<template>
  <div class="task-monitor-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
          <h1>任务实时监控</h1>
          <div></div>
        </div>
      </el-header>

      <el-main class="main-content">
        <!-- 任务状态概览 -->
        <el-row :gutter="20" class="status-overview">
          <el-col :span="6">
            <el-card class="status-card">
              <div class="status-content">
                <div class="status-icon" :class="getStatusClass()">
                  <el-icon v-if="taskStatus.status === 'running'">
                    <Loading />
                  </el-icon>
                  <el-icon v-else-if="taskStatus.status === 'completed'">
                    <CircleCheck />
                  </el-icon>
                  <el-icon v-else-if="taskStatus.status === 'failed'">
                    <CircleClose />
                  </el-icon>
                  <el-icon v-else>
                    <Clock />
                  </el-icon>
                </div>
                <div class="status-info">
                  <div class="status-title">任务状态</div>
                  <div class="status-value">
                    {{ getStatusText(taskStatus.status) }}
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="6">
            <el-card class="progress-card">
              <div class="progress-content">
                <el-progress
                  :percentage="taskProgress"
                  :status="getProgressStatus()"
                  :stroke-width="8"
                />
                <div class="progress-text">{{ taskProgress.toFixed(1) }}%</div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-content">
                <div class="stats-number">
                  {{ taskStatus.totalPapers || 0 }}
                </div>
                <div class="stats-label">总试卷数</div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-content">
                <div class="stats-number">
                  {{ taskStatus.averageScore || 0 }}
                </div>
                <div class="stats-label">平均分</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 详细进度 -->
        <el-row :gutter="20" class="detail-stats">
          <el-col :span="8">
            <el-card>
              <template #header>
                <h4>批改进度</h4>
              </template>
              <div class="progress-detail">
                <div class="progress-item">
                  <span class="label">已完成:</span>
                  <span class="value">{{
                    taskStatus.completedPapers || 0
                  }}</span>
                </div>
                <div class="progress-item">
                  <span class="label">失败:</span>
                  <span class="value error">{{
                    taskStatus.failedPapers || 0
                  }}</span>
                </div>
                <div class="progress-item">
                  <span class="label">剩余:</span>
                  <span class="value">{{
                    (taskStatus.totalPapers || 0) -
                    (taskStatus.completedPapers || 0) -
                    (taskStatus.failedPapers || 0)
                  }}</span>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="8">
            <el-card>
              <template #header>
                <h4>执行时间</h4>
              </template>
              <div class="time-detail">
                <div class="time-item">
                  <span class="label">开始时间:</span>
                  <span class="value">{{ formatTime(task.createdAt) }}</span>
                </div>
                <div class="time-item">
                  <span class="label">预计完成:</span>
                  <span class="value">{{ getEstimatedCompletion() }}</span>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="8">
            <el-card>
              <template #header>
                <h4>操作</h4>
              </template>
              <div class="action-buttons">
                <el-button
                  type="danger"
                  :disabled="taskStatus.status !== 'running'"
                  @click="cancelTask"
                >
                  取消任务
                </el-button>
                <el-button
                  type="primary"
                  :disabled="taskStatus.status !== 'completed'"
                  @click="viewAnalysis"
                >
                  查看分析
                </el-button>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 实时日志 -->
        <el-card class="logs-card">
          <template #header>
            <div class="logs-header">
              <h4>实时日志</h4>
              <div class="logs-controls">
                <el-button size="small" @click="clearLogs">
                  清空日志
                </el-button>
                <el-button size="small" @click="exportLogs">
                  导出日志
                </el-button>
              </div>
            </div>
          </template>

          <div class="logs-container" ref="logsContainer">
            <div
              v-for="log in logs"
              :key="log.id"
              :class="['log-item', `log-${log.level}`]"
            >
              <span class="log-time">{{ formatLogTime(log.timestamp) }}</span>
              <span class="log-level">{{ log.level }}</span>
              <span class="log-message">{{ log.message }}</span>
            </div>

            <div v-if="logs.length === 0" class="no-logs">暂无日志信息</div>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, reactive, onMounted, onUnmounted, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  Loading,
  CircleCheck,
  CircleClose,
  Clock,
} from "@element-plus/icons-vue";

export default {
  name: "TaskMonitor",
  components: {
    ArrowLeft,
    Loading,
    CircleCheck,
    CircleClose,
    Clock,
  },
  setup() {
    const route = useRoute();
    const router = useRouter();
    const logsContainer = ref(null);

    const task = reactive({});
    const taskStatus = reactive({
      status: "pending",
      totalPapers: 0,
      completedPapers: 0,
      failedPapers: 0,
      averageScore: 0,
      message: "",
    });

    const logs = ref([]);
    const ws = ref(null);
    const pollingInterval = ref(null);

    const taskId = route.params.id;

    const taskProgress = ref(0);

    // 计算进度
    const calculateProgress = () => {
      const total = taskStatus.totalPapers || 0;
      const completed = taskStatus.completedPapers || 0;
      const failed = taskStatus.failedPapers || 0;

      if (total === 0) {
        taskProgress.value = 0;
        return;
      }

      taskProgress.value = ((completed + failed) / total) * 100;
    };

    const getStatusClass = () => {
      return `status-${taskStatus.status}`;
    };

    const getProgressStatus = () => {
      if (taskStatus.status === "failed") return "exception";
      if (taskStatus.status === "completed") return "success";
      return "";
    };

    const getStatusText = (status) => {
      const statusMap = {
        pending: "等待执行",
        running: "执行中",
        completed: "已完成",
        failed: "执行失败",
        cancelled: "已取消",
      };
      return statusMap[status] || status;
    };

    const formatTime = (timeStr) => {
      if (!timeStr) return "-";
      return new Date(timeStr).toLocaleString("zh-CN");
    };

    const formatLogTime = (timestamp) => {
      return new Date(timestamp).toLocaleTimeString("zh-CN");
    };

    const getEstimatedCompletion = () => {
      if (taskStatus.status !== "running" || taskProgress.value === 0) {
        return "-";
      }

      const rate =
        (taskProgress.value /
          (Date.now() - new Date(task.createdAt).getTime())) *
        1000;
      const remaining = (100 - taskProgress.value) / rate;

      if (remaining > 0 && remaining < 86400) {
        const hours = Math.floor(remaining / 3600);
        const minutes = Math.floor((remaining % 3600) / 60);
        return `${hours}小时${minutes}分钟后`;
      }

      return "计算中...";
    };

    const loadTaskStatus = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(`/api/teacher/tasks/${taskId}/status`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          Object.assign(taskStatus, data);
          calculateProgress();

          // 添加状态更新日志
          if (data.message) {
            addLog("info", data.message);
          }
        }
      } catch (error) {
        console.error("加载任务状态失败:", error);
      }
    };

    const loadTaskDetail = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch(`/api/teacher/tasks/${taskId}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          Object.assign(task, data.task);
          addLog("info", "任务详情加载完成");
        }
      } catch (error) {
        console.error("加载任务详情失败:", error);
      }
    };

    const cancelTask = async () => {
      try {
        await ElMessageBox.confirm("确定要取消这个任务吗？", "取消任务", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        });

        const token = localStorage.getItem("token");
        const response = await fetch(`/api/teacher/tasks/${taskId}/cancel`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          ElMessage.success("任务已取消");
          addLog("warning", "用户手动取消任务");
        } else {
          const errorData = await response.json();
          ElMessage.error(errorData.error || "取消任务失败");
        }
      } catch (error) {
        if (error !== "cancel") {
          ElMessage.error("取消任务失败");
        }
      }
    };

    const viewAnalysis = () => {
      router.push(`/task-analysis/${taskId}`);
    };

    const goBack = () => {
      router.push("/teacher-home");
    };

    const addLog = (level, message) => {
      const log = {
        id: Date.now(),
        level: level,
        message: message,
        timestamp: new Date().toISOString(),
      };

      logs.value.unshift(log);

      // 限制日志数量
      if (logs.value.length > 100) {
        logs.value = logs.value.slice(0, 100);
      }

      // 滚动到底部
      nextTick(() => {
        if (logsContainer.value) {
          logsContainer.value.scrollTop = logsContainer.value.scrollHeight;
        }
      });
    };

    const clearLogs = () => {
      logs.value = [];
    };

    const exportLogs = () => {
      const logText = logs.value
        .map(
          (log) =>
            `${formatLogTime(log.timestamp)} [${log.level}] ${log.message}`,
        )
        .join("\n");

      const blob = new Blob([logText], { type: "text/plain" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `task-${taskId}-logs.txt`;
      a.click();
      URL.revokeObjectURL(url);
    };

    // 启动WebSocket连接
    const startWebSocket = () => {
      try {
        ws.value = new WebSocket(`ws://localhost:3000/ws/task/${taskId}`);

        ws.value.onopen = () => {
          addLog("info", "WebSocket连接已建立");
        };

        ws.value.onmessage = (event) => {
          try {
            const status = JSON.parse(event.data);
            Object.assign(taskStatus, status);
            calculateProgress();
          } catch (error) {
            addLog("error", `WebSocket消息解析错误: ${error.message}`);
          }
        };

        ws.value.onclose = () => {
          addLog("warning", "WebSocket连接已断开，切换到轮询模式");
          // 降级到轮询模式
          startPolling();
        };

        ws.value.onerror = (error) => {
          addLog("error", `WebSocket连接错误: ${error}`);
        };
      } catch (error) {
        addLog("error", `WebSocket初始化失败: ${error.message}`);
        startPolling();
      }
    };

    // 启动轮询模式
    const startPolling = () => {
      if (pollingInterval.value) {
        clearInterval(pollingInterval.value);
      }

      pollingInterval.value = setInterval(() => {
        loadTaskStatus();
      }, 30000); // 30秒轮询一次
    };

    // 清理资源
    const cleanup = () => {
      if (ws.value) {
        ws.value.close();
        ws.value = null;
      }

      if (pollingInterval.value) {
        clearInterval(pollingInterval.value);
        pollingInterval.value = null;
      }
    };

    onMounted(async () => {
      await loadTaskDetail();
      await loadTaskStatus();

      addLog("info", "任务监控页面已加载");
      addLog("info", `开始监控任务 ${taskId}`);

      // 尝试启动WebSocket
      startWebSocket();
    });

    onUnmounted(() => {
      cleanup();
    });

    return {
      task,
      taskStatus,
      taskProgress,
      logs,
      logsContainer,
      getStatusClass,
      getProgressStatus,
      getStatusText,
      formatTime,
      formatLogTime,
      getEstimatedCompletion,
      cancelTask,
      viewAnalysis,
      goBack,
      clearLogs,
      exportLogs,
    };
  },
};
</script>

<style scoped>
.task-monitor-container {
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
  color: #67c23a;
  font-size: 24px;
}

.main-content {
  padding: 20px;
}

.status-overview {
  margin-bottom: 20px;
}

.status-card,
.progress-card,
.stats-card {
  height: 120px;
}

.status-content {
  display: flex;
  align-items: center;
  height: 100%;
}

.status-icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
}

.status-icon.status-running {
  color: #e6a23c;
}

.status-icon.status-completed {
  color: #67c23a;
}

.status-icon.status-failed {
  color: #f56c6c;
}

.status-icon.status-pending {
  color: #909399;
}

.status-info {
  flex: 1;
}

.status-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 5px;
}

.status-value {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.progress-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
  padding: 0 20px;
}

.progress-text {
  text-align: center;
  margin-top: 10px;
  font-size: 16px;
  font-weight: bold;
  color: #67c23a;
}

.stats-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.stats-number {
  font-size: 32px;
  font-weight: bold;
  color: #67c23a;
  margin-bottom: 10px;
}

.stats-label {
  color: #606266;
  font-size: 14px;
}

.detail-stats {
  margin-bottom: 20px;
}

.progress-detail,
.time-detail {
  padding: 20px 0;
}

.progress-item,
.time-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.label {
  color: #606266;
}

.value {
  font-weight: 500;
  color: #303133;
}

.value.error {
  color: #f56c6c;
}

.action-buttons {
  display: flex;
  gap: 10px;
  justify-content: center;
}

.logs-card {
  min-height: 400px;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logs-controls {
  display: flex;
  gap: 10px;
}

.logs-container {
  height: 300px;
  overflow-y: auto;
  background: #f8f9fa;
  border-radius: 4px;
  padding: 10px;
}

.log-item {
  display: flex;
  align-items: center;
  padding: 5px 10px;
  margin-bottom: 5px;
  background: white;
  border-radius: 3px;
  font-family: "Courier New", monospace;
  font-size: 12px;
}

.log-info {
  border-left: 3px solid #409eff;
}

.log-warning {
  border-left: 3px solid #e6a23c;
}

.log-error {
  border-left: 3px solid #f56c6c;
}

.log-time {
  color: #909399;
  margin-right: 10px;
  min-width: 80px;
}

.log-level {
  margin-right: 10px;
  min-width: 60px;
  font-weight: bold;
}

.log-message {
  flex: 1;
}

.no-logs {
  text-align: center;
  color: #909399;
  padding: 50px;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style>
