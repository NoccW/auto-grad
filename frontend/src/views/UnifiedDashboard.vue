<template>
  <div class="unified-dashboard">
    <!-- 顶部导航栏 -->
    <el-header class="header">
      <div class="header-left">
        <h1 class="logo">智能改卷系统</h1>
      </div>
      <div class="header-center">
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          @select="handleMenuSelect"
          class="nav-menu"
        >
          <el-menu-item index="dashboard">首页</el-menu-item>
          <el-menu-item v-if="userRole === 'parent'" index="upload"
            >上传作业</el-menu-item
          >
          <el-menu-item index="results">成绩查看</el-menu-item>
          <el-menu-item index="history">历史记录</el-menu-item>
          <el-menu-item v-if="userRole === 'teacher'" index="tasks"
            >任务管理</el-menu-item
          >
        </el-menu>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleUserMenu">
          <span class="user-info">
            <el-avatar>{{ username.charAt(0) }}</el-avatar>
            <span class="username">{{ username }}</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人资料</el-dropdown-item>
              <el-dropdown-item command="settings">系统设置</el-dropdown-item>
              <el-dropdown-item
                v-if="userRole === 'teacher'"
                command="switch-parent"
                >切换到家长端</el-dropdown-item
              >
              <el-dropdown-item
                v-if="userRole === 'parent'"
                command="switch-teacher"
                >切换到教师端</el-dropdown-item
              >
              <el-dropdown-item divided command="logout"
                >退出登录</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <el-dialog v-model="showProfileDialog" title="个人资料" width="420px">
      <el-form label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="profileForm.name" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="profileForm.email" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showProfileDialog = false">取消</el-button>
          <el-button type="primary" @click="saveProfile">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 主要内容区域 -->
    <el-main class="main-content">
      <!-- 根据角色显示不同的仪表板内容 -->
      <div v-if="userRole === 'teacher'" class="teacher-dashboard">
        <TeacherDashboard />
      </div>
      <div v-else-if="userRole === 'parent'" class="parent-dashboard">
        <ParentDashboard />
      </div>
      <div v-else class="role-selection">
        <el-card class="role-card">
          <template #header>
            <h2>选择用户角色</h2>
          </template>
          <div class="role-options">
            <el-card
              class="role-option"
              :class="{ active: selectedRole === 'teacher' }"
              @click="selectRole('teacher')"
            >
              <div class="role-icon">
                <el-icon size="64"><UserFilled /></el-icon>
              </div>
              <h3>教师端</h3>
              <p>创建和管理自动化改卷任务，查看统计数据和班级分析</p>
              <el-button type="primary" size="large">进入教师端</el-button>
            </el-card>

            <el-card
              class="role-option"
              :class="{ active: selectedRole === 'parent' }"
              @click="selectRole('parent')"
            >
              <div class="role-icon">
                <el-icon size="64"><Avatar /></el-icon>
              </div>
              <h3>家长端</h3>
              <p>上传学生作业，查看改卷结果和学习进度分析</p>
              <el-button type="success" size="large">进入家长端</el-button>
            </el-card>
          </div>
        </el-card>
      </div>
    </el-main>
  </div>
</template>

<script setup>
  import { ref, reactive, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { ArrowDown, UserFilled, Avatar } from "@element-plus/icons-vue";
import TeacherDashboard from "./TeacherHome.vue";
import ParentDashboard from "./ParentDashboard.vue";

const router = useRouter();
const activeMenu = ref("dashboard");
const selectedRole = ref("");

// 用户信息
const userRole = computed(() => localStorage.getItem("userRole") || "");
  const username = computed(() => {
    const cached = localStorage.getItem("username");
    if (cached) return cached;
    if (userRole.value === "teacher") {
      return "张老师";
    } else if (userRole.value === "parent") {
      return "李家长";
    }
    return "用户";
  });

  const showProfileDialog = ref(false);
  const profileForm = reactive({
    name: "",
    email: "",
  });

const handleMenuSelect = (index) => {
  activeMenu.value = index;

  switch (index) {
    case "dashboard":
      router.push("/");
      break;
    case "upload":
      router.push("/upload");
      break;
    case "results":
      router.push("/result/latest");
      break;
    case "history":
      router.push("/history");
      break;
    case "tasks":
      router.push("/teacher");
      break;
  }
};

const handleUserMenu = async (command) => {
  switch (command) {
      case "profile":
        await openProfile();
        break;
    case "settings":
      ElMessage.info("系统设置功能开发中...");
      break;
    case "switch-parent":
      await switchRole("parent");
      break;
    case "switch-teacher":
      await switchRole("teacher");
      break;
    case "logout":
      await logout();
      break;
  }
};

const selectRole = (role) => {
  selectedRole.value = role;
};

  const switchRole = async (role) => {
  try {
    await ElMessageBox.confirm(
      `确定要切换到${role === "teacher" ? "教师" : "家长"}端吗？`,
      "角色切换",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "info",
      },
    );

    // 更新用户角色
    localStorage.setItem("userRole", role);

    // 重新加载页面
    window.location.reload();
  } catch {
    // 用户取消操作
  }
  };

  const openProfile = async () => {
    try {
      const token = localStorage.getItem("token");
      const resp = await fetch("/api/auth/me", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (resp.ok) {
        const data = await resp.json();
        profileForm.name = data.username || "";
        profileForm.email = data.email || "";
        showProfileDialog.value = true;
      } else {
        ElMessage.error("获取个人资料失败");
      }
    } catch (e) {
      console.error(e);
      ElMessage.error("获取个人资料失败");
    }
  };

  const saveProfile = async () => {
    try {
      const token = localStorage.getItem("token");
      const resp = await fetch("/api/auth/profile", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(profileForm),
      });
      if (!resp.ok) {
        const err = await resp.json();
        throw new Error(err.error || "保存失败");
      }
      localStorage.setItem("username", profileForm.name);
      ElMessage.success("资料已更新");
      showProfileDialog.value = false;
    } catch (e) {
      console.error(e);
      ElMessage.error(e.message || "保存失败");
    }
  };

const logout = async () => {
  try {
    await ElMessageBox.confirm("确定要退出登录吗？", "退出确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    // 清除用户信息
    localStorage.removeItem("token");
    localStorage.removeItem("userRole");
    localStorage.removeItem("username");

    // 跳转到登录页
    router.push("/login");
    ElMessage.success("已退出登录");
  } catch {
    // 用户取消操作
  }
};

onMounted(() => {
  // 如果用户没有角色，显示角色选择页面
  if (!userRole.value) {
    activeMenu.value = "dashboard";
  } else {
    // 根据角色设置默认菜单
    activeMenu.value = "dashboard";
  }
});
</script>

<style scoped>
.unified-dashboard {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-left .logo {
  margin: 0;
  color: #409eff;
  font-size: 20px;
  font-weight: 600;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.nav-menu {
  border-bottom: none;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.username {
  margin: 0 8px;
  font-weight: 500;
}

.main-content {
  flex: 1;
  padding: 20px;
  background-color: #f5f7fa;
}

.role-selection {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
}

.role-card {
  max-width: 800px;
  width: 100%;
}

.role-card h2 {
  text-align: center;
  margin: 0;
  color: #303133;
}

.role-options {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
  margin-top: 30px;
}

.role-option {
  text-align: center;
  padding: 30px;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.role-option:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.role-option.active {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.role-icon {
  margin-bottom: 20px;
  color: #409eff;
}

.role-option h3 {
  margin: 0 0 15px 0;
  color: #303133;
  font-size: 20px;
}

.role-option p {
  margin: 0 0 25px 0;
  color: #606266;
  line-height: 1.6;
}

.teacher-dashboard,
.parent-dashboard {
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .header {
    padding: 0 10px;
  }

  .header-left .logo {
    font-size: 16px;
  }

  .nav-menu {
    display: none;
  }

  .role-options {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .role-option {
    padding: 20px;
  }
}
</style>
