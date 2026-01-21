require('dotenv').config();
const axios = require('axios');
const qs = require('qs');
const fs = require('fs');
const path = require('path');

// --- 配置区域 ---
const DEEPSEEK_API_KEY = process.env.DEEPSEEK_API_KEY;
const BAIDU_API_KEY = process.env.BAIDU_API_KEY;
const BAIDU_SECRET_KEY = process.env.BAIDU_SECRET_KEY;

// 输入图片文件夹路径
const INPUT_DIR = './papers';
// 输出结果文件路径
const OUTPUT_FILE = './results.json';

// 评分规则
const GRADING_RULES = `
按题目分数每题给分
`;

let BAIDU_ACCESS_TOKEN = null;
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// --- 核心逻辑 ---
async function run() {
    console.log(`=== 本地智能阅卷脚本 (OCR + DeepSeek) ===`);
    
    // 0. 检查环境
    if (!fs.existsSync(INPUT_DIR)) {
        console.error(`❌ 错误：找不到图片文件夹 "${INPUT_DIR}"，请先创建并放入图片。`);
        return;
    }
    if (!BAIDU_API_KEY || !BAIDU_SECRET_KEY || !DEEPSEEK_API_KEY) {
        console.error("❌ 错误：请在 .env 文件中配置所有 API Key");
        return;
    }

    // 1. 获取百度 Token
    console.log("🔄 正在获取百度 OCR 授权...");
    await refreshBaiduToken();
    if (!BAIDU_ACCESS_TOKEN) return;
    console.log("✅ 百度授权成功！");

    // 2. 读取图片列表
    const files = fs.readdirSync(INPUT_DIR).filter(file => {
        const ext = path.extname(file).toLowerCase();
        return ['.jpg', '.jpeg', '.png', '.bmp'].includes(ext);
    });

    console.log(`📂 发现 ${files.length} 张图片，开始处理...`);
    console.log("------------------------------------------------");

    const results = [];

    // 3. 循环处理
    for (let i = 0; i < files.length; i++) {
        const fileName = files[i];
        const filePath = path.join(INPUT_DIR, fileName);
        
        console.log(`\n[${i + 1}/${files.length}] 正在处理: ${fileName}`);

        try {
            // 读取文件并转 Base64
            const fileBuffer = fs.readFileSync(filePath);
            const imageBase64 = fileBuffer.toString('base64');

            // Step A: 百度 OCR 识别
            process.stdout.write("   👀 OCR 识别中... ");
            const studentAnswer = await recognizeHandwriting(imageBase64);
            
            if (!studentAnswer) {
                console.log("❌ 识别失败或内容为空");
                results.push({ file: fileName, score: 0, answer: "", reason: "OCR失败" });
                continue;
            }
            console.log("✅ 完成");
            console.log(`   📝 内容预览: ${studentAnswer.substring(0, 30)}...`);

            // Step B: DeepSeek 评分
            process.stdout.write("   🤖 DeepSeek 评分中... ");
            const score = await getScoreFromDeepSeek(studentAnswer, GRADING_RULES);
            console.log(`✅ 得分: ${score}`);

            // 保存结果
            results.push({
                file: fileName,
                score: score,
                answer: studentAnswer
            });

            // 避免触发 API 速率限制 (QPS)
            await delay(1000); 

        } catch (err) {
            console.error(`   ❌ 处理异常: ${err.message}`);
        }
    }

    // 4. 输出最终报告
    console.log("\n================================================");
    console.log("🎉 处理完成！");
    fs.writeFileSync(OUTPUT_FILE, JSON.stringify(results, null, 2), 'utf8');
    console.log(`💾 结果已保存至: ${OUTPUT_FILE}`);
    
    // 简单打印结果表
    console.table(results.map(r => ({ 文件名: r.file, 分数: r.score, 字数: r.answer.length })));
}

// --- 辅助函数 ---

async function refreshBaiduToken() {
    try {
        const url = `https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=${BAIDU_API_KEY}&client_secret=${BAIDU_SECRET_KEY}`;
        const res = await axios.post(url);
        if (res.data.access_token) BAIDU_ACCESS_TOKEN = res.data.access_token;
    } catch (e) { 
        console.error("\n❌ 百度Token获取失败:", e.response ? e.response.data : e.message); 
    }
}

async function recognizeHandwriting(imageBase64) {
    if (!BAIDU_ACCESS_TOKEN) return "";
    
    // 使用高精度含位置版 (accurate_basic) 或 手写文字识别 (handwriting)
    // 注意：百度通用文字识别（高精度版）免费额度较少，手写版可能更适合
    const url = `https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic?access_token=${BAIDU_ACCESS_TOKEN}`;
    
    try {
        // 百度OCR限制 body 大小，如果图片过大可能需要压缩，这里直接传
        const data = qs.stringify({ 
            image: imageBase64, 
            language_type: 'CHN_ENG' 
        });
        const res = await axios.post(url, data, { 
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' } 
        });
        
        if (res.data.words_result) {
            return res.data.words_result.map(item => item.words).join("，");
        } else {
            return "";
        }
    } catch (e) { 
        console.error("OCR API 错误:", e.message);
        return ""; 
    }
}

async function getScoreFromDeepSeek(studentAnswer, gradingRules) {
    if (!DEEPSEEK_API_KEY) return 0;
    
    const prompt = `
你是一名语文阅卷老师。
【评分规则】：
${gradingRules}

【学生回答】：
${studentAnswer}

请根据【评分规则】对【学生回答】进行打分。
要求：
1. 忽略OCR识别产生的明显错别字，关注语义是否符合得分点。
2. 只要意思对即可给分。
    `;
    
    try {
        const response = await axios.post('https://api.deepseek.com/chat/completions', {
            model: "deepseek-chat", // 或者 deepseek-reasoner
            messages: [{ role: "user", content: prompt }],
            temperature: 0.1 // 降低随机性，让分数更稳定
        }, { 
            headers: { 
                'Authorization': `Bearer ${DEEPSEEK_API_KEY}`, 
                'Content-Type': 'application/json' 
            } 
        });
        
        const content = response.data.choices[0].message.content;
        // 提取内容中的第一个数字
        const match = content.match(/\d+/);
        return match ? parseInt(match[0]) : 0;
    } catch (e) { 
        console.error("DeepSeek API 错误:", e.message);
        return 0; 
    }
}

run();