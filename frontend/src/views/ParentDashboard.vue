<template>
  <div class="parent-dashboard">
    <!-- 学生信息卡片 -->
    <el-row :gutter="20" class="info-row">
      <el-col :span="24">
        <el-card class="student-info-card">
          <template #header>
            <div class="card-header">
              <span>学生信息</span>
              <el-button type="text" @click="openStudentDialog">
                <el-icon><Edit /></el-icon>
              </el-button>
            </div>
          </template>
          <div class="student-info">
            <div class="info-item">
              <span class="label">姓名：</span>
              <span class="value">{{ studentInfo.name }}</span>
            </div>
            <div class="info-item">
              <span class="label">班级：</span>
              <span class="value">{{ studentInfo.class }}</span>
            </div>
            <div class="info-item">
              <span class="label">学校：</span>
              <span class="value">{{ studentInfo.school }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快速操作卡片 -->
    <el-row :gutter="20" class="actions-row">
      <el-col :span="8">
        <el-card class="action-card upload-card" @click="goToUpload">
          <div class="action-content">
            <el-icon size="48" class="action-icon"><UploadFilled /></el-icon>
            <h3>上传作业</h3>
            <p>上传学生作业进行智能改卷</p>
            <el-button type="primary">立即上传</el-button>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="action-card results-card" @click="goToResults">
          <div class="action-content">
            <el-icon size="48" class="action-icon"><Document /></el-icon>
            <h3>查看成绩</h3>
            <p>查看最新改卷结果和反馈</p>
            <el-button type="success">查看成绩</el-button>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="action-card history-card" @click="goToHistory">
          <div class="action-content">
            <el-icon size="48" class="action-icon"><Clock /></el-icon>
            <h3>历史记录</h3>
            <p>查看历史改卷记录和统计</p>
            <el-button type="warning">历史记录</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近结果和统计 -->
    <el-row :gutter="20" class="data-row">
      <el-col :span="16">
        <el-card class="recent-results-card">
          <template #header>
            <div class="card-header">
              <span>最近改卷结果</span>
              <el-button type="text" @click="viewAllResults"
                >查看全部</el-button
              >
            </div>
          </template>
          <div class="recent-results">
            <div
              v-for="result in recentResults"
              :key="result.id"
              class="result-item"
              @click="viewResult(result.id)"
            >
              <div class="result-info">
                <div class="subject">{{ result.subject }}</div>
                <div class="score-info">
                  <span class="score">{{ result.score }}</span>
                  <span class="total-score">/ {{ result.totalScore }}</span>
                </div>
                <div class="date">{{ formatDate(result.submitTime) }}</div>
              </div>
              <div class="status">
                <el-tag :type="getStatusType(result.status)" size="small">
                  {{ getStatusText(result.status) }}
                </el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="statistics-card">
          <template #header>
            <span>学习统计</span>
          </template>
          <div class="statistics">
            <div class="stat-item">
              <div class="stat-value">{{ statistics.totalSubmissions }}</div>
              <div class="stat-label">总提交次数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">
                {{ statistics.averageScore.toFixed(1) }}
              </div>
              <div class="stat-label">平均分数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ statistics.completedTasks }}</div>
              <div class="stat-label">已完成任务</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ statistics.pendingTasks }}</div>
              <div class="stat-label">待处理任务</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 学习进度图表 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="24">
        <el-card class="progress-card">
          <template #header>
            <span>学习进度</span>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder">
              <el-icon size="64"><TrendCharts /></el-icon>
              <p>学习进度图表</p>
              <div class="progress-data">
                <div class="subject-progress">
                  <span>数学</span>
                  <el-progress :percentage="85" status="success" />
                </div>
                <div class="subject-progress">
                  <span>语文</span>
                  <el-progress :percentage="92" status="success" />
                </div>
                <div class="subject-progress">
                  <span>英语</span>
                  <el-progress :percentage="78" status="warning" />
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-model="showStudentDialog" title="编辑学生信息" width="400px">
      <el-form label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="studentForm.name" />
        </el-form-item>
        <el-form-item label="班级">
          <el-input v-model="studentForm.class" />
        </el-form-item>
        <el-form-item label="学校">
          <el-input v-model="studentForm.school" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showStudentDialog = false">取消</el-button>
          <el-button type="primary" @click="saveStudentInfo">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {
  Edit,
  UploadFilled,
  Document,
  Clock,
  TrendCharts,
} from "@element-plus/icons-vue";

const router = useRouter();

// 学生信息
const studentInfo = ref({
  name: "李小明",
  class: "三年级一班",
  school: "示例小学",
});
const studentForm = reactive({
  name: "",
  class: "",
  school: "",
});
const showStudentDialog = ref(false);

// 最近结果
const recentResults = ref([]);

// 统计数据
const statistics = ref({
  totalSubmissions: 0,
  averageScore: 0,
  completedTasks: 0,
  pendingTasks: 0,
});

const openStudentDialog = () => {
  studentForm.name = studentInfo.value.name;
  studentForm.class = studentInfo.value.class;
  studentForm.school = studentInfo.value.school;
  showStudentDialog.value = true;
};

const saveStudentInfo = async () => {
  try {
    const response = await fetch("/api/parent/student", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(studentForm),
    });
    if (!response.ok) {
      const err = await response.json();
      throw new Error(err.error || "更新学生信息失败");
    }
    studentInfo.value = { ...studentForm };
    ElMessage.success("学生信息已更新");
    showStudentDialog.value = false;
  } catch (error) {
    console.error("学生信息更新错误:", error);
    ElMessage.error(error.message || "更新学生信息失败");
  }
};

const goToUpload = () => {
  router.push("/upload");
};

const goToResults = () => {
  router.push("/result/latest");
};

const goToHistory = () => {
  router.push("/history");
};

const viewAllResults = () => {
  router.push("/result/latest");
};

const viewResult = (id) => {
  router.push(`/result/${id}`);
};

const formatDate = (dateString) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleDateString("zh-CN");
};

const getStatusType = (status) => {
  const statusMap = {
    completed: "success",
    processing: "warning",
    failed: "danger",
  };
  return statusMap[status] || "info";
};

const getStatusText = (status) => {
  const statusMap = {
    completed: "已完成",
    processing: "处理中",
    failed: "失败",
  };
  return statusMap[status] || "未知";
};

const loadDashboardData = async () => {
  try {
    const response = await fetch("/api/parent/dashboard");
    if (!response.ok) {
      throw new Error("Failed to load dashboard data");
    }

    const data = await response.json();

    if (data.studentInfo) {
      studentInfo.value = data.studentInfo;
    }

    recentResults.value = data.recentResults || [];

    statistics.value = data.statistics || {
      totalSubmissions: 0,
      averageScore: 0,
      completedTasks: 0,
      pendingTasks: 0,
    };
  } catch (error) {
    ElMessage.error("加载仪表板数据失败");
    console.error("Error loading dashboard data:", error);
  }
};

onMounted(() => {
  loadDashboardData();
});
</script>

<style scoped>
.parent-dashboard {
  padding: 20px;
}

.info-row {
  margin-bottom: 20px;
}

.student-info-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.student-info {
  display: flex;
  gap: 40px;
}

.info-item {
  display: flex;
  align-items: center;
}

.info-item .label {
  font-weight: 500;
  color: #606266;
  margin-right: 8px;
}

.info-item .value {
  font-weight: 600;
  color: #303133;
}

.actions-row {
  margin-bottom: 20px;
}

.action-card {
  cursor: pointer;
  transition: all 0.3s;
  height: 100%;
}

.action-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.action-content {
  text-align: center;
  padding: 20px;
}

.action-icon {
  margin-bottom: 15px;
}

.action-content h3 {
  margin: 0 0 10px 0;
  color: #303133;
}

.action-content p {
  margin: 0 0 20px 0;
  color: #606266;
}

.data-row {
  margin-bottom: 20px;
}

.recent-results-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.recent-results {
  max-height: 400px;
  overflow-y: auto;
}

.result-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.3s;
}

.result-item:hover {
  background-color: #f8f9fa;
}

.result-item:last-child {
  border-bottom: none;
}

.result-info {
  flex: 1;
}

.result-info .subject {
  font-weight: 600;
  color: #303133;
  margin-bottom: 5px;
}

.score-info {
  color: #606266;
  margin-bottom: 5px;
}

.score {
  font-weight: 600;
  color: #409eff;
  font-size: 18px;
}

.total-score {
  color: #909399;
  font-size: 14px;
}

.result-info .date {
  font-size: 12px;
  color: #c0c4cc;
}

.statistics {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
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

.progress-data {
  margin-top: 30px;
  width: 100%;
  max-width: 400px;
}

.subject-progress {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.subject-progress span {
  width: 60px;
  font-size: 14px;
  color: #606266;
}

.subject-progress .el-progress {
  flex: 1;
  margin-left: 15px;
}

@media (max-width: 768px) {
  .student-info {
    flex-direction: column;
    gap: 10px;
  }

  .actions-row .el-col {
    margin-bottom: 15px;
  }

  .result-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .result-info .date {
    align-self: flex-end;
  }
}
</style>
