<template>
  <el-config-provider :locale="zhCn">
    <router-view v-if="!$route.meta.requiresAuth || isAuthenticated" />
    <div v-else-if="!isRoleValid($route.meta.role)" class="access-denied">
      <div class="access-denied-content">
        <el-result
          icon="warning"
          title="访问被拒绝"
          sub-title="您没有权限访问此页面，请使用正确的角色登录。"
        >
          <template #extra>
            <el-button type="primary" @click="goToHome">返回首页</el-button>
            <el-button @click="logout">重新登录</el-button>
          </template>
        </el-result>
      </div>
    </div>
    <router-view v-else />
  </el-config-provider>
</template>

<script>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import zhCn from "element-plus/dist/locale/zh-cn.mjs";

export default {
  name: "App",
  setup() {
    const router = useRouter();
    const isAuthenticated = ref(false);
    const userRole = ref("");

    const checkAuth = async () => {
      const token = localStorage.getItem("token");
      if (!token) {
        if (router.currentRoute.value.path !== "/login") {
          router.push("/login");
        }
        return;
      }

      try {
        const response = await fetch("/api/auth/me", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          isAuthenticated.value = true;
          userRole.value = data.userRole || "";

          // 检查路由权限
          if (router.currentRoute.value.path === "/login") {
            // 根据角色跳转到相应主页
            if (data.userRole === "teacher") {
              router.push("/teacher-home");
            } else {
              router.push("/");
            }
          } else if (!isRoleValid(router.currentRoute.value.meta.role)) {
            // 如果权限不匹配，跳转到对应的主页
            if (data.userRole === "teacher") {
              router.push("/teacher-home");
            } else {
              router.push("/");
            }
          }
        } else {
          localStorage.removeItem("token");
          localStorage.removeItem("user");
          if (router.currentRoute.value.path !== "/login") {
            router.push("/login");
          }
        }
      } catch (error) {
        console.error("认证检查失败:", error);
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        if (router.currentRoute.value.path !== "/login") {
          router.push("/login");
        }
      }
    };

    const isRoleValid = (requiredRole) => {
      if (!requiredRole || !userRole.value) {
        return true;
      }
      return userRole.value === requiredRole;
    };

    const goToHome = () => {
      if (userRole.value === "teacher") {
        router.push("/teacher-home");
      } else {
        router.push("/");
      }
    };

    const logout = () => {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      isAuthenticated.value = false;
      userRole.value = "";
      router.push("/login");
    };

    onMounted(async () => {
      await checkAuth();
    });

    return {
      isAuthenticated,
      userRole,
      isRoleValid,
      goToHome,
      logout,
      zhCn,
    };
  },
};
</script>

<style>
#app {
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue",
    Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.loading-container {
  width: 100vw;
  height: 100vh;
}

.access-denied {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}

.access-denied-content {
  text-align: center;
  max-width: 500px;
  padding: 40px;
}

body {
  margin: 0;
  padding: 0;
}
</style>
