package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BaiduOCRService struct {
	apiKey    string
	secretKey string
	token     string
	tokenExp  time.Time
}

type BaiduTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type BaiduOCRResponse struct {
	WordsResult []struct {
		Words string `json:"words"`
	} `json:"words_result"`
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

func NewBaiduOCRService(apiKey, secretKey string) *BaiduOCRService {
	return &BaiduOCRService{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (s *BaiduOCRService) GetAccessToken() error {
	if s.apiKey == "" || s.secretKey == "" {
		return errors.New("Baidu OCR credentials not configured")
	}

	if s.token != "" && time.Now().Before(s.tokenExp) {
		return nil
	}

	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		s.apiKey, s.secretKey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp BaiduTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	s.token = tokenResp.AccessToken
	s.tokenExp = time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second) // 提前5分钟过期

	return nil
}

func (s *BaiduOCRService) RecognizeText(imageData []byte) (string, error) {
	if err := s.GetAccessToken(); err != nil {
		return "", err
	}

	apiURL := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=%s", s.token)

	form := url.Values{}
	form.Set("image", base64.StdEncoding.EncodeToString(imageData))

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ocrResp BaiduOCRResponse
	if err := json.Unmarshal(body, &ocrResp); err != nil {
		return "", fmt.Errorf("failed to parse OCR response: %w", err)
	}

	if ocrResp.ErrorCode != 0 {
		return "", fmt.Errorf("OCR API error: %d - %s", ocrResp.ErrorCode, ocrResp.ErrorMsg)
	}

	var result string
	for _, word := range ocrResp.WordsResult {
		result += word.Words + "\n"
	}

	return result, nil
}
