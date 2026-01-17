package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DeepSeekService struct {
	apiKey  string
	baseURL string
}

type DeepSeekRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
}

type DeepSeekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

type GradingResult struct {
	Score          int      `json:"score"`
	TotalScore     int      `json:"totalScore"`
	WrongQuestions []string `json:"wrongQuestions"`
	CorrectAnswers []string `json:"correctAnswers"`
	Feedback       string   `json:"feedback"`
}

func NewDeepSeekService(apiKey string) *DeepSeekService {
	return &DeepSeekService{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com",
	}
}

func (s *DeepSeekService) GradePaper(ocrText, referenceAnswer string) (*GradingResult, error) {
	if s.apiKey == "" {
		return nil, errors.New("DeepSeek API key not configured")
	}

	prompt := s.buildGradingPrompt(ocrText, referenceAnswer)

	request := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: prompt},
		},
		Stream:      false,
		Temperature: 0.1,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := s.baseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response DeepSeekResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("DeepSeek API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from DeepSeek API")
	}

	return s.parseGradingResult(response.Choices[0].Message.Content)
}

func (s *DeepSeekService) buildGradingPrompt(ocrText, referenceAnswer string) string {
	prompt := `你是一个专业的试卷批改助手。请根据以下OCR识别的学生答案和参考答案，进行评分和分析。

学生答案：
` + ocrText + `

参考答案：
` + referenceAnswer + `

请按照以下JSON格式返回评分结果：
{
  "score": 85,
  "totalScore": 100,
  "wrongQuestions": ["第3题", "第7题"],
  "correctAnswers": ["第3题答案：C", "第7题答案：A"],
  "feedback": "整体表现良好，建议加强数学计算练习"
}

评分标准：
1. 内容准确性 (40分)
2. 逻辑完整性 (30分)
3. 表达清晰度 (20分)
4. 格式规范性 (10分)

请严格按照JSON格式返回，不要包含其他文字。`
	return prompt
}

func (s *DeepSeekService) parseGradingResult(content string) (*GradingResult, error) {
	var result GradingResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 如果JSON解析失败，返回默认结果
		return &GradingResult{
			Score:          75,
			TotalScore:     100,
			WrongQuestions: []string{"解析失败"},
			CorrectAnswers: []string{"请手动检查"},
			Feedback:       "AI解析失败，请手动批改",
		}, nil
	}
	return &result, nil
}
