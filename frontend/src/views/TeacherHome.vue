<template>
  <div class="teacher-home-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
          <h1>教师控制台</h1>
          <div class="user-info">
            <el-dropdown>
              <span class="el-dropdown-link">
                <el-avatar :size="32" :src="userAvatar" />
                <span class="username">{{ userName }}</span>
                <el-icon class="el-icon--right">
                  <arrow-down />
                </el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="logout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>

      <el-main class="main-content">
        <!-- 欢迎卡片 -->
        <div class="welcome-section">
          <el-card class="welcome-card">
            <h2>欢迎使用智能改卷教师端</h2>
            <p>自动化批改试卷，实时监控进度，智能分析数据</p>
          </el-card>
        </div>

        <!-- 快速操作 -->
        <div class="actions-section">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="action-card" @click="goToTaskConfig">
                <el-icon class="action-icon"><Setting /></el-icon>
                <h3>创建任务</h3>
                <p>配置新的自动化批改任务</p>
              </el-card>
            </el-col>

            <el-col :span="6">
              <el-card class="action-card" @click="goToTaskHistory">
                <el-icon class="action-icon"><Clock /></el-icon>
                <h3>任务历史</h3>
                <p>查看历史任务记录</p>
              </el-card>
            </el-col>

            <el-col :span="6">
              <el-card class="action-card" @click="goToTaskMonitor">
                <el-icon class="action-icon"><Monitor /></el-icon>
                <h3>实时监控</h3>
                <p>监控任务执行状态</p>
              </el-card>
            </el-col>

            <el-col :span="6">
              <el-card class="action-card" @click="goToTaskAnalysis">
                <el-icon class="action-icon"><DataAnalysis /></el-icon>
                <h3>数据分析</h3>
                <p>查看批改数据分析</p>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <!-- 最近任务 -->
        <div class="recent-tasks-section">
          <el-card>
            <template #header>
              <div class="recent-header">
                <span>最近任务</span>
                <el-button type="text" @click="goToTaskHistory"
                  >查看全部</el-button
                >
              </div>
            </template>

            <el-table
              :data="recentTasks"
              v-loading="loading"
              empty-text="暂无任务"
            >
              <el-table-column prop="id" label="任务ID" width="120" />
              <el-table-column
                prop="targetUrl"
                label="目标网站"
                show-overflow-tooltip
              />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="scope">
                  <el-tag :type="getStatusType(scope.row.status)">
                    {{ getStatusText(scope.row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column
                prop="totalPapers"
                label="试卷总数"
                width="100"
              />
              <el-table-column
                prop="completedPapers"
                label="已完成"
                width="100"
              />
              <el-table-column prop="averageScore" label="平均分" width="100">
                <template #default="scope">
                  <span v-if="scope.row.averageScore">{{
                    scope.row.averageScore.toFixed(1)
                  }}</span>
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="创建时间" width="180">
                <template #default="scope">
                  {{ formatDate(scope.row.createdAt) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200" fixed="right">
                <template #default="scope">
                  <el-button
                    type="primary"
                    size="small"
                    @click="goToMonitor(scope.row.id)"
                    :disabled="scope.row.status === 'cancelled'"
                  >
                    监控
                  </el-button>
                  <el-button
                    type="success"
                    size="small"
                    @click="goToAnalysis(scope.row.id)"
                    :disabled="scope.row.status !== 'completed'"
                  >
                    分析
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </div>

        <!-- 统计卡片 -->
        <div class="statistics-section">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.totalTasks }}</div>
                  <div class="stat-label">总任务数</div>
                </div>
                <el-icon class="stat-icon"><Document /></el-icon>
              </el-card>
            </el-col>

            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.runningTasks }}</div>
                  <div class="stat-label">运行中</div>
                </div>
                <el-icon class="stat-icon"><Loading /></el-icon>
              </el-card>
            </el-col>

            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.completedTasks }}</div>
                  <div class="stat-label">已完成</div>
                </div>
                <el-icon class="stat-icon"><CircleCheck /></el-icon>
              </el-card>
            </el-col>

            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.totalPapers }}</div>
                  <div class="stat-label">总批改数</div>
                </div>
                <el-icon class="stat-icon"><Files /></el-icon>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {
  ArrowLeft,
  SwitchButton,
  Setting,
  Clock,
  Monitor,
  DataAnalysis,
  Document,
  Loading,
  CircleCheck,
  Files,
} from "@element-plus/icons-vue";

export default {
  name: "TeacherHome",
  components: {
    ArrowLeft,
    SwitchButton,
    Setting,
    Clock,
    Monitor,
    DataAnalysis,
    Document,
    Loading,
    CircleCheck,
    Files,
  },
  setup() {
    const router = useRouter();
    const loading = ref(false);
    const recentTasks = ref([]);
    const statistics = ref({
      totalTasks: 0,
      runningTasks: 0,
      completedTasks: 0,
      totalPapers: 0,
    });
    const userName = ref("教师用户");
    const userAvatar = ref("");

    const loadRecentTasks = async () => {
      loading.value = true;
      try {
        const token = localStorage.getItem("token");
        const response = await fetch("/api/teacher/tasks?limit=5", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          recentTasks.value = data.tasks || [];
          calculateStatistics(data.tasks || []);
        } else {
          const errorData = await response.json();
          ElMessage.error(errorData.error || "加载任务失败");
        }
      } catch (error) {
        console.error("加载任务失败:", error);
        ElMessage.error("加载任务失败");
      } finally {
        loading.value = false;
      }
    };

    const calculateStatistics = (tasks) => {
      statistics.value.totalTasks = tasks.length;
      statistics.value.runningTasks = tasks.filter(
        (t) => t.status === "running",
      ).length;
      statistics.value.completedTasks = tasks.filter(
        (t) => t.status === "completed",
      ).length;
      statistics.value.totalPapers = tasks.reduce(
        (sum, t) => sum + (t.totalPapers || 0),
        0,
      );
    };

    const getStatusType = (status) => {
      const statusMap = {
        pending: "info",
        running: "warning",
        completed: "success",
        failed: "danger",
        cancelled: "info",
      };
      return statusMap[status] || "info";
    };

    const getStatusText = (status) => {
      const statusMap = {
        pending: "待执行",
        running: "运行中",
        completed: "已完成",
        failed: "失败",
        cancelled: "已取消",
      };
      return statusMap[status] || status;
    };

    const formatDate = (dateStr) => {
      return new Date(dateStr).toLocaleString("zh-CN");
    };

    const goBack = () => {
      router.push("/");
    };

    const goToTaskConfig = () => {
      router.push("/task-config");
    };

    const goToTaskHistory = () => {
      router.push("/task-history");
    };

    const goToTaskMonitor = (taskId) => {
      router.push(`/task-monitor/${taskId}`);
    };

    const goToTaskAnalysis = (taskId) => {
      router.push(`/task-analysis/${taskId}`);
    };

    const logout = async () => {
      try {
        await ElMessageBox.confirm("确定要退出登录吗？", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        });

        localStorage.removeItem("token");
        localStorage.removeItem("user");
        ElMessage.success("退出成功");
        router.push("/login");
      } catch (error) {
        // 用户取消
      }
    };

    const loadUserInfo = () => {
      const userStr = localStorage.getItem("user");
      if (userStr) {
        const user = JSON.parse(userStr);
        userName.value = user.name || user.openId || "教师用户";
      }
    };

    onMounted(() => {
      loadUserInfo();
      loadRecentTasks();
    });

    return {
      loading,
      recentTasks,
      statistics,
      userName,
      userAvatar,
      getStatusType,
      getStatusText,
      formatDate,
      goBack,
      goToTaskConfig,
      goToTaskHistory,
      goToTaskMonitor,
      goToTaskAnalysis,
      logout,
    };
  },
};
</script>

<style scoped>
.teacher-home-container {
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
  font-size: 24px;
}

.user-info {
  display: flex;
  align-items: center;
}

.el-dropdown-link {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.username {
  margin: 0 8px;
}

.main-content {
  padding: 20px;
}

.welcome-section {
  margin-bottom: 30px;
}

.welcome-card {
  text-align: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  color: white;
}

.welcome-card h2 {
  margin: 0 0 10px 0;
  font-size: 28px;
}

.welcome-card p {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
}

.actions-section {
  margin-bottom: 30px;
}

.action-card {
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.action-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.action-icon {
  font-size: 48px;
  color: #67c23a;
  margin-bottom: 15px;
}

.action-card h3 {
  margin: 0 0 10px 0;
  color: #303133;
}

.action-card p {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.recent-tasks-section {
  margin-bottom: 30px;
}

.recent-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.statistics-section {
  margin-bottom: 30px;
}

.stat-card {
  position: relative;
  overflow: hidden;
}

.stat-content {
  padding: 20px 20px 20px 80px;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.stat-label {
  color: #606266;
  font-size: 14px;
}

.stat-icon {
  position: absolute;
  top: 50%;
  right: 20px;
  transform: translateY(-50%);
  font-size: 40px;
  color: #67c23a;
  opacity: 0.2;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style>
