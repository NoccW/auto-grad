<template>
  <div class="result-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
          <h1>批改结果</h1>
          <div></div>
        </div>
      </el-header>

      <el-main class="main-content" v-loading="loading">
        <div v-if="gradingRecord">
          <!-- 结果概览 -->
          <el-card class="overview-card">
            <el-row :gutter="20">
              <el-col :span="8">
                <div class="score-display">
                  <div class="score-circle">
                    <el-progress
                      type="circle"
                      :percentage="scorePercentage"
                      :color="scoreColor"
                      :width="120"
                    >
                      <template #default="{ percentage }">
                        <span class="score-text"
                          >{{ gradingRecord.aiScore || 0 }}分</span
                        >
                      </template>
                    </el-progress>
                  </div>
                  <h3>得分</h3>
                </div>
              </el-col>

              <el-col :span="8">
                <div class="status-display">
                  <el-tag
                    :type="getStatusType(gradingRecord.status)"
                    size="large"
                  >
                    {{ getStatusText(gradingRecord.status) }}
                  </el-tag>
                  <h3>批改状态</h3>
                </div>
              </el-col>

              <el-col :span="8">
                <div class="time-display">
                  <div class="time-info">
                    <el-icon><Clock /></el-icon>
                    <span>{{ formatDate(gradingRecord.createdAt) }}</span>
                  </div>
                  <h3>批改时间</h3>
                </div>
              </el-col>
            </el-row>
          </el-card>

          <!-- 图片展示 -->
          <el-row :gutter="20" class="image-row">
            <el-col :span="12">
              <el-card>
                <template #header>
                  <h3>试卷图片</h3>
                </template>
                <el-image
                  :src="getImageUrl(gradingRecord.paperImageUrl)"
                  fit="contain"
                  style="width: 100%; height: 400px"
                  :preview-src-list="[getImageUrl(gradingRecord.paperImageUrl)]"
                />
              </el-card>
            </el-col>

            <el-col :span="12" v-if="gradingRecord.answerImageUrl">
              <el-card>
                <template #header>
                  <h3>参考答案</h3>
                </template>
                <el-image
                  :src="getImageUrl(gradingRecord.answerImageUrl)"
                  fit="contain"
                  style="width: 100%; height: 400px"
                  :preview-src-list="[
                    getImageUrl(gradingRecord.answerImageUrl),
                  ]"
                />
              </el-card>
            </el-col>
          </el-row>

          <!-- OCR识别结果 -->
          <el-card class="ocr-card" v-if="gradingRecord.ocrResult">
            <template #header>
              <h3>OCR识别结果</h3>
            </template>
            <el-input
              v-model="gradingRecord.ocrResult"
              type="textarea"
              :rows="6"
              readonly
              placeholder="OCR识别的文字内容"
            />
          </el-card>

          <!-- AI 反馈 -->
          <el-card class="feedback-card" v-if="gradingRecord.feedback">
            <template #header>
              <h3>AI 评分反馈</h3>
            </template>
            <el-input
              v-model="gradingRecord.feedback"
              type="textarea"
              :rows="4"
              readonly
              placeholder="AI 评分反馈"
            />
          </el-card>

          <!-- 错题分析 -->
          <el-card class="analysis-card" v-if="gradingRecord.wrongQuestions">
            <template #header>
              <h3>错题分析</h3>
            </template>
            <div class="wrong-questions" v-if="parsedWrongQuestions.length > 0">
              <el-tag
                v-for="(question, index) in parsedWrongQuestions"
                :key="index"
                type="danger"
                class="question-tag"
              >
                {{ question }}
              </el-tag>
            </div>
            <div v-else>
              <p>暂无错题记录</p>
            </div>
          </el-card>

          <!-- 正确答案 -->
          <el-card class="answers-card" v-if="gradingRecord.correctAnswers">
            <template #header>
              <h3>正确答案</h3>
            </template>
            <div class="correct-answers" v-if="parsedCorrectAnswers.length > 0">
              <div
                v-for="(answer, index) in parsedCorrectAnswers"
                :key="index"
                class="answer-item"
              >
                <el-tag type="success">{{ answer }}</el-tag>
              </div>
            </div>
            <div v-else>
              <p>暂无答案记录</p>
            </div>
          </el-card>

          <!-- 操作按钮 -->
          <div class="action-buttons">
            <el-button type="primary" @click="downloadReport">
              <el-icon><Download /></el-icon>
              下载报告
            </el-button>
            <el-button @click="shareResult">
              <el-icon><Share /></el-icon>
              分享结果
            </el-button>
            <el-button @click="printResult">
              <el-icon><Printer /></el-icon>
              打印结果
            </el-button>
            <el-button type="success" @click="regrade">
              <el-icon><Refresh /></el-icon>
              重新批改
            </el-button>
          </div>
        </div>

        <el-empty v-else description="未找到批改结果" />
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {
  ArrowLeft,
  Clock,
  Download,
  Share,
  Printer,
  Refresh,
} from "@element-plus/icons-vue";

export default {
  name: "Result",
  components: {
    ArrowLeft,
    Clock,
    Download,
    Share,
    Printer,
    Refresh,
  },
  setup() {
    const route = useRoute();
    const router = useRouter();
    const loading = ref(false);
    const gradingRecord = ref(null);
    let pollTimer = null;

    const scorePercentage = computed(() => {
      return gradingRecord.value?.aiScore || 0;
    });

    const scoreColor = computed(() => {
      const score = scorePercentage.value;
      if (score >= 90) return "#67c23a";
      if (score >= 80) return "#409eff";
      if (score >= 70) return "#e6a23c";
      if (score >= 60) return "#f56c6c";
      return "#f56c6c";
    });

    const parsedWrongQuestions = computed(() => {
      if (!gradingRecord.value?.wrongQuestions) return [];
      try {
        const parsed = JSON.parse(gradingRecord.value.wrongQuestions);
        return parsed.questions || [];
      } catch {
        return [];
      }
    });

    const parsedCorrectAnswers = computed(() => {
      if (!gradingRecord.value?.correctAnswers) return [];
      try {
        const parsed = JSON.parse(gradingRecord.value.correctAnswers);
        return parsed.answers || [];
      } catch {
        return [];
      }
    });

    const loadGradingResult = async () => {
      loading.value = true;
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(`/api/grading/${route.params.id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          gradingRecord.value = data;
          if (data.status === "processing" || data.status === "pending") {
            schedulePoll();
          } else {
            clearPoll();
          }
        } else {
          const errorData = await response.json();
          ElMessage.error(errorData.error || "加载结果失败");
        }
      } catch (error) {
        console.error("加载结果错误:", error);
        ElMessage.error("加载结果失败");
      } finally {
        loading.value = false;
      }
    };

    const schedulePoll = () => {
      clearPoll();
      pollTimer = setTimeout(() => {
        loadGradingResult();
      }, 3000);
    };

    const clearPoll = () => {
      if (pollTimer) {
        clearTimeout(pollTimer);
        pollTimer = null;
      }
    };

    const getStatusType = (status) => {
      const statusMap = {
        pending: "info",
        processing: "warning",
        completed: "success",
        failed: "danger",
      };
      return statusMap[status] || "info";
    };

    const getStatusText = (status) => {
      const statusMap = {
        pending: "待处理",
        processing: "处理中",
        completed: "已完成",
        failed: "失败",
      };
      return statusMap[status] || status;
    };

    const formatDate = (dateStr) => {
      return new Date(dateStr).toLocaleString("zh-CN");
    };

    const getImageUrl = (relativePath) => {
      return `/uploads/${relativePath}`;
    };

    const goBack = () => {
      router.push("/");
    };

    const downloadReport = () => {
      ElMessage.info("下载功能开发中...");
    };

    const shareResult = () => {
      ElMessage.info("分享功能开发中...");
    };

    const printResult = () => {
      window.print();
    };

    const regrade = async () => {
      try {
        loading.value = true;
        const token = localStorage.getItem("token");

        const response = await fetch(
          `/api/grading/${route.params.id}/process`,
          {
            method: "POST",
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );

        if (response.ok) {
          ElMessage.success("已重新提交批改任务");
          gradingRecord.value = {
            ...(gradingRecord.value || {}),
            status: "processing",
            aiScore: 0,
          };
          schedulePoll();
        } else {
          const errorData = await response.json();
          ElMessage.error(errorData.error || "重新批改失败");
        }
      } catch (error) {
        console.error("重新批改错误:", error);
        ElMessage.error("重新批改失败");
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      loadGradingResult();
    });

    onUnmounted(() => {
      clearPoll();
    });

    return {
      loading,
      gradingRecord,
      scorePercentage,
      scoreColor,
      parsedWrongQuestions,
      parsedCorrectAnswers,
      getStatusType,
      getStatusText,
      formatDate,
      getImageUrl,
      goBack,
      downloadReport,
      shareResult,
      printResult,
      regrade,
    };
  },
};
</script>

<style scoped>
.result-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  padding: 0 20px;
}

.header-content h1 {
  margin: 0;
  color: #409eff;
  font-size: 20px;
}

.main-content {
  padding: 20px;
}

.overview-card {
  margin-bottom: 20px;
  text-align: center;
}

.score-display,
.status-display,
.time-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 180px;
}

.score-circle {
  margin-bottom: 15px;
}

.score-text {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.status-display h3,
.time-display h3 {
  margin-top: 15px;
  color: #606266;
}

.time-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;
}

.image-row {
  margin-bottom: 20px;
}

.ocr-card,
.feedback-card,
.analysis-card,
.answers-card {
  margin-bottom: 20px;
}

.wrong-questions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.question-tag {
  margin: 5px;
}

.correct-answers {
  space-y: 10px;
}

.answer-item {
  margin-bottom: 10px;
}

.action-buttons {
  text-align: center;
  margin-top: 30px;
}

.action-buttons .el-button {
  margin: 0 10px;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-tag) {
  margin: 2px;
}

@media print {
  .header {
    display: none;
  }

  .action-buttons {
    display: none;
  }
}
</style>
