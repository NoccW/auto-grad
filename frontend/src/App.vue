<template>
  <el-config-provider :locale="zhCn">
    <router-view v-if="$route.meta.requiresAuth === false || isAuthenticated" />
    <div v-else class="loading-container">
      <el-loading :fullscreen="true" />
    </div>
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

    onMounted(async () => {
      // 检查本地存储的token
      const token = localStorage.getItem("token");
      if (token) {
        try {
          // 验证token有效性
          const response = await fetch("/api/auth/me", {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          if (response.ok) {
            isAuthenticated.value = true;
            // 如果已经在登录页面，跳转到主页
            if (router.currentRoute.value.path === "/login") {
              router.push("/");
            }
          } else {
            localStorage.removeItem("token");
            if (router.currentRoute.value.path !== "/login") {
              router.push("/login");
            }
          }
        } catch (error) {
          localStorage.removeItem("token");
          if (router.currentRoute.value.path !== "/login") {
            router.push("/login");
          }
        }
      } else if (router.currentRoute.value.path !== "/login") {
        router.push("/login");
      } else {
        isAuthenticated.value = true;
      }
    });

    return {
      isAuthenticated,
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

body {
  margin: 0;
  padding: 0;
}
</style>
