<template>
  <div class="task-analysis">
    <el-page-header @back="goBack" title="返回任务列表">
      <template #content>
        <span class="page-title">任务数据分析</span>
      </template>
    </el-page-header>

    <div class="analysis-content" v-if="task && analysis">
      <!-- 任务概览 -->
      <el-card class="overview-card" header="任务概览">
        <div class="overview-grid">
          <div class="overview-item">
            <div class="label">任务ID</div>
            <div class="value">{{ task.id }}</div>
          </div>
          <div class="overview-item">
            <div class="label">目标网站</div>
            <div class="value">{{ task.targetUrl }}</div>
          </div>
          <div class="overview-item">
            <div class="label">账号</div>
            <div class="value">{{ task.account }}</div>
          </div>
          <div class="overview-item">
            <div class="label">状态</div>
            <el-tag :type="getStatusType(task.status)">{{
              getStatusText(task.status)
            }}</el-tag>
          </div>
        </div>
      </el-card>

      <!-- 统计图表 -->
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card header="分数分布">
            <div class="chart-container">
              <div class="chart-placeholder">
                <el-icon size="64"><PieChart /></el-icon>
                <p>分数分布图表</p>
                <div class="score-data">
                  <div
                    v-for="item in analysis.scoreDistribution"
                    :key="item.range"
                    class="score-item"
                  >
                    <span class="range">{{ item.range }}</span>
                    <span class="count"
                      >{{ item.count }}人 ({{ item.percentage }}%)</span
                    >
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="12">
          <el-card header="错误分析">
            <div class="error-analysis">
              <div
                v-if="
                  analysis.errorAnalysis && analysis.errorAnalysis.length > 0
                "
              >
                <div
                  v-for="(error, index) in analysis.errorAnalysis"
                  :key="index"
                  class="error-item"
                >
                  <el-tag type="danger" size="small">低分</el-tag>
                  <span class="error-info">
                    {{ error.studentName }} - {{ error.score }}分 -
                    {{ error.reason }}
                  </span>
                </div>
              </div>
              <div v-else class="no-data">
                <el-icon size="32"><SuccessFilled /></el-icon>
                <p>暂无错误数据</p>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 性能指标 -->
      <el-card class="metrics-card" header="性能指标">
        <el-row :gutter="20">
          <el-col
            :span="6"
            v-for="(value, key) in analysis.performanceMetrics"
            :key="key"
          >
            <div class="metric-item">
              <div class="metric-label">{{ getMetricLabel(key) }}</div>
              <div class="metric-value">{{ formatMetric(value, key) }}</div>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 导出按钮 -->
      <div class="actions">
        <el-button type="primary" @click="exportReport">
          <el-icon><Download /></el-icon>
          导出分析报告
        </el-button>
        <el-button @click="refreshAnalysis">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { ElMessage } from "element-plus";
import {
  PieChart,
  Download,
  Refresh,
  SuccessFilled,
} from "@element-plus/icons-vue";

const router = useRouter();
const route = useRoute();

const task = ref(null);
const analysis = ref(null);
const loading = ref(false);

const goBack = () => {
  router.push("/teacher-home");
};

const getStatusType = (status) => {
  const statusMap = {
    completed: "success",
    running: "warning",
    failed: "danger",
    pending: "info",
  };
  return statusMap[status] || "info";
};

const getStatusText = (status) => {
  const statusMap = {
    completed: "已完成",
    running: "运行中",
    failed: "失败",
    pending: "待执行",
  };
  return statusMap[status] || "未知";
};

const getMetricLabel = (key) => {
  const labelMap = {
    averageScore: "平均分",
    maxScore: "最高分",
    minScore: "最低分",
    passRate: "及格率",
  };
  return labelMap[key] || key;
};

const formatMetric = (value, key) => {
  if (key.includes("Rate")) {
    return `${value}%`;
  }
  return value;
};

const exportReport = () => {
  ElMessage.info("导出功能开发中...");
};

const refreshAnalysis = async () => {
  await loadTaskAnalysis();
  ElMessage.success("数据已刷新");
};

const loadTaskAnalysis = async () => {
  try {
    loading.value = true;
    const taskId = route.params.id;

    // 模拟API调用
    await new Promise((resolve) => setTimeout(resolve, 1000));

    task.value = {
      id: taskId,
      targetUrl: "https://www.7net.cc/build/home/index.html",
      account: "123123",
      status: "completed",
    };

    analysis.value = {
      scoreDistribution: [
        { range: "优秀(90-100)", count: 5, percentage: 10 },
        { range: "良好(80-89)", count: 15, percentage: 30 },
        { range: "及格(60-79)", count: 20, percentage: 40 },
        { range: "不及格(0-59)", count: 10, percentage: 20 },
      ],
      errorAnalysis: [
        { studentName: "学生1", score: 45, reason: "计算错误" },
        { studentName: "学生2", score: 52, reason: "概念理解错误" },
      ],
      performanceMetrics: {
        averageScore: 75.5,
        maxScore: 98,
        minScore: 45,
        passRate: 80.0,
      },
    };
  } catch (error) {
    ElMessage.error("加载分析数据失败");
    console.error("Error loading analysis:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadTaskAnalysis();
});
</script>

<style scoped>
.task-analysis {
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.analysis-content {
  margin-top: 20px;
}

.overview-card {
  margin-bottom: 20px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.overview-item {
  text-align: center;
}

.overview-item .label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.overview-item .value {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.chart-container {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #909399;
}

.score-data {
  margin-top: 20px;
}

.score-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  padding: 4px 0;
}

.range {
  font-weight: 500;
}

.count {
  color: #409eff;
}

.error-analysis {
  max-height: 300px;
  overflow-y: auto;
}

.error-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px;
  background-color: #fef0f0;
  border-radius: 4px;
}

.error-info {
  margin-left: 8px;
  font-size: 14px;
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}

.metrics-card {
  margin-top: 20px;
}

.metric-item {
  text-align: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.metric-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.metric-value {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
}

.actions {
  margin-top: 20px;
  text-align: center;
}

.loading-container {
  padding: 20px;
}
</style>
