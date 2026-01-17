<template>
  <div class="home-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1>智能改卷助手</h1>
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
        <div class="welcome-section">
          <el-card class="welcome-card">
            <h2>欢迎使用智能改卷助手</h2>
            <p>使用AI技术，为您提供快速、准确的试卷批改服务</p>
          </el-card>
        </div>

        <div class="actions-section">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-card class="action-card" @click="goToUpload">
                <el-icon class="action-icon"><Upload /></el-icon>
                <h3>上传试卷</h3>
                <p>拍摄或上传试卷图片进行智能批改</p>
              </el-card>
            </el-col>

            <el-col :span="8">
              <el-card class="action-card" @click="goToHistory">
                <el-icon class="action-icon"><Clock /></el-icon>
                <h3>历史记录</h3>
                <p>查看历史批改记录和评分结果</p>
              </el-card>
            </el-col>

            <el-col :span="8">
              <el-card class="action-card">
                <el-icon class="action-icon"><DataAnalysis /></el-icon>
                <h3>统计分析</h3>
                <p>查看学习进度和错题分析</p>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div class="recent-section">
          <el-card>
            <template #header>
              <div class="recent-header">
                <span>最近批改记录</span>
                <el-button type="text" @click="goToHistory">查看全部</el-button>
              </div>
            </template>

            <el-table :data="recentRecords" v-loading="loading">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="scope">
                  <el-tag :type="getStatusType(scope.row.status)">
                    {{ getStatusText(scope.row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="aiScore" label="得分" width="100">
                <template #default="scope">
                  <span v-if="scope.row.aiScore"
                    >{{ scope.row.aiScore }}分</span
                  >
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="时间" width="180">
                <template #default="scope">
                  {{ formatDate(scope.row.createdAt) }}
                </template>
              </el-table-column>
              <el-table-column label="操作">
                <template #default="scope">
                  <el-button
                    type="text"
                    @click="viewDetail(scope.row.id)"
                    :disabled="scope.row.status !== 'completed'"
                  >
                    查看详情
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowDown,
  SwitchButton,
  Upload,
  Clock,
  DataAnalysis,
} from "@element-plus/icons-vue";

export default {
  name: "Home",
  components: {
    ArrowDown,
    SwitchButton,
    Upload,
    Clock,
    DataAnalysis,
  },
  setup() {
    const router = useRouter();
    const loading = ref(false);
    const recentRecords = ref([]);
    const userName = ref("用户");
    const userAvatar = ref("");

    const loadRecentRecords = async () => {
      loading.value = true;
      try {
        const token = localStorage.getItem("token");
        const response = await fetch("/api/grading?limit=5", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          recentRecords.value = data.records || [];
        }
      } catch (error) {
        console.error("加载记录失败:", error);
      } finally {
        loading.value = false;
      }
    };

    const loadUserInfo = () => {
      const userStr = localStorage.getItem("user");
      if (userStr) {
        const user = JSON.parse(userStr);
        userName.value = user.name || user.openId || "用户";
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

    const goToUpload = () => {
      router.push("/upload");
    };

    const goToHistory = () => {
      router.push("/history");
    };

    const viewDetail = (id) => {
      router.push(`/result/${id}`);
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

    onMounted(() => {
      loadUserInfo();
      loadRecentRecords();
    });

    return {
      loading,
      recentRecords,
      userName,
      userAvatar,
      getStatusType,
      getStatusText,
      formatDate,
      goToUpload,
      goToHistory,
      viewDetail,
      logout,
    };
  },
};
</script>

<style scoped>
.home-container {
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
  background: linear-gradient(135deg, #409eff 0%, #36cfc9 100%);
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
  color: #409eff;
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

.recent-section {
  background: white;
  border-radius: 8px;
}

.recent-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style>
