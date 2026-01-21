<template>
  <div class="task-history">
    <el-page-header @back="goBack" title="返回教师主页">
      <template #content>
        <span class="page-title">任务历史记录</span>
      </template>
    </el-page-header>

    <div class="history-content">
      <!-- 搜索和筛选 -->
      <el-card class="filter-card">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-input
              v-model="searchQuery"
              placeholder="搜索任务ID或账号"
              prefix-icon="Search"
              clearable
            />
          </el-col>
          <el-col :span="4">
            <el-select v-model="statusFilter" placeholder="状态筛选" clearable>
              <el-option label="全部" value="" />
              <el-option label="已完成" value="completed" />
              <el-option label="运行中" value="running" />
              <el-option label="失败" value="failed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
            />
          </el-col>
          <el-col :span="4">
            <el-button type="primary" @click="loadTasks">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
          </el-col>
        </el-row>
      </el-card>

      <!-- 任务列表 -->
      <el-card class="tasks-card">
        <template #header>
          <div class="card-header">
            <span>任务列表 ({{ filteredTasks.length }})</span>
            <div class="header-actions">
              <el-button size="small" @click="exportHistory">
                <el-icon><Download /></el-icon>
                导出记录
              </el-button>
              <el-button size="small" @click="refreshTasks">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </div>
        </template>

        <el-table :data="paginatedTasks" stripe style="width: 100%">
          <el-table-column prop="id" label="任务ID" width="180" />
          <el-table-column prop="targetUrl" label="目标网站" min-width="200">
            <template #default="{ row }">
              <el-link :href="row.targetUrl" target="_blank" type="primary">
                {{ formatUrl(row.targetUrl) }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column prop="account" label="账号" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="totalPapers" label="总试卷数" width="100" />
          <el-table-column prop="completedPapers" label="已完成" width="100" />
          <el-table-column prop="averageScore" label="平均分" width="100">
            <template #default="{ row }">
              <span v-if="row.averageScore > 0">{{
                row.averageScore.toFixed(1)
              }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="160">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="viewTask(row)">查看</el-button>
              <el-button
                v-if="row.status === 'running'"
                size="small"
                type="warning"
                @click="monitorTask(row)"
              >
                监控
              </el-button>
              <el-button
                v-if="row.status === 'completed'"
                size="small"
                type="success"
                @click="analyzeTask(row)"
              >
                分析
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="filteredTasks.length"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Search, Download, Refresh } from "@element-plus/icons-vue";

const router = useRouter();

const searchQuery = ref("");
const statusFilter = ref("");
const dateRange = ref([]);
const tasks = ref([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(20);

const goBack = () => {
  router.push("/teacher-home");
};

const getStatusType = (status) => {
  const statusMap = {
    completed: "success",
    running: "warning",
    failed: "danger",
    cancelled: "info",
    pending: "info",
  };
  return statusMap[status] || "info";
};

const getStatusText = (status) => {
  const statusMap = {
    completed: "已完成",
    running: "运行中",
    failed: "失败",
    cancelled: "已取消",
    pending: "待执行",
  };
  return statusMap[status] || "未知";
};

const formatUrl = (url) => {
  if (!url) return "-";
  try {
    const domain = new URL(url).hostname;
    return domain;
  } catch {
    return url.length > 30 ? url.substring(0, 30) + "..." : url;
  }
};

const formatDate = (dateString) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleString("zh-CN");
};

const filteredTasks = computed(() => {
  let filtered = tasks.value;

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(
      (task) =>
        task.id.toLowerCase().includes(query) ||
        task.account.toLowerCase().includes(query),
    );
  }

  // 状态过滤
  if (statusFilter.value) {
    filtered = filtered.filter((task) => task.status === statusFilter.value);
  }

  // 日期过滤
  if (dateRange.value && dateRange.value.length === 2) {
    const [startDate, endDate] = dateRange.value;
    filtered = filtered.filter((task) => {
      const taskDate = new Date(task.createdAt).toISOString().split("T")[0];
      return taskDate >= startDate && taskDate <= endDate;
    });
  }

  return filtered;
});

const paginatedTasks = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredTasks.value.slice(start, end);
});

const viewTask = (task) => {
  router.push(`/task-monitor/${task.id}`);
};

const monitorTask = (task) => {
  router.push(`/task-monitor/${task.id}`);
};

const analyzeTask = (task) => {
  router.push(`/task-analysis/${task.id}`);
};

const handleSizeChange = (size) => {
  pageSize.value = size;
  currentPage.value = 1;
};

const handleCurrentChange = (page) => {
  currentPage.value = page;
};

const exportHistory = () => {
  ElMessage.info("导出功能开发中...");
};

const refreshTasks = async () => {
  await loadTasks();
  ElMessage.success("数据已刷新");
};

const loadTasks = async () => {
  try {
    loading.value = true;

    // 模拟API调用
    await new Promise((resolve) => setTimeout(resolve, 1000));

    tasks.value = [
      {
        id: "task_1640996805",
        targetUrl: "https://www.7net.cc/build/home/index/index.html",
        account: "123123",
        status: "completed",
        totalPapers: 100,
        completedPapers: 95,
        failedPapers: 5,
        averageScore: 78.5,
        createdAt: "2024-01-17T10:30:00Z",
        updatedAt: "2024-01-17T11:30:00Z",
      },
      {
        id: "task_1640996806",
        targetUrl: "https://www.7net.cc/build/home/index.html",
        account: "test123",
        status: "running",
        totalPapers: 200,
        completedPapers: 120,
        failedPapers: 5,
        averageScore: 82.3,
        createdAt: "2024-01-17T09:00:00Z",
        updatedAt: "2024-01-17T09:00:00Z",
      },
      {
        id: "task_1640996807",
        targetUrl: "https://www.example.com",
        account: "demo456",
        status: "failed",
        totalPapers: 50,
        completedPapers: 15,
        failedPapers: 35,
        averageScore: 65.2,
        createdAt: "2024-01-16T15:20:00Z",
        updatedAt: "2024-01-16T16:45:00Z",
      },
      {
        id: "task_1640996808",
        targetUrl: "https://www.test-site.com",
        account: "user789",
        status: "cancelled",
        totalPapers: 80,
        completedPapers: 30,
        failedPapers: 0,
        averageScore: 0,
        createdAt: "2024-01-15T14:10:00Z",
        updatedAt: "2024-01-15T14:30:00Z",
      },
    ];
  } catch (error) {
    ElMessage.error("加载任务历史失败");
    console.error("Error loading tasks:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadTasks();
});
</script>

<style scoped>
.task-history {
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.history-content {
  margin-top: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.tasks-card {
  min-height: 500px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>
