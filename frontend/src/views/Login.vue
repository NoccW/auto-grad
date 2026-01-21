<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>智能改卷系统</h2>
          <p>统一登录入口</p>
        </div>
      </template>

      <el-form
        :model="loginForm"
        :rules="rules"
        ref="loginFormRef"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item label="角色" prop="role">
          <el-radio-group v-model="loginForm.role">
            <el-radio label="parent">家长端</el-radio>
            <el-radio label="teacher">教师端</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>

        <el-form-item>
          <el-button type="text" @click="goToRegister" style="width: 100%">
            没有账号？立即注册
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";

export default {
  name: "Login",
  setup() {
    const router = useRouter();
    const loginFormRef = ref();
    const loading = ref(false);

    const loginForm = reactive({
      username: "",
      password: "",
      role: "parent",
    });

    const rules = {
      username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
      password: [{ required: true, message: "请输入密码", trigger: "blur" }],
      role: [{ required: true, message: "请选择角色", trigger: "change" }],
    };

    const handleLogin = async () => {
      if (!loginFormRef.value) return;

      try {
        await loginFormRef.value.validate();
        loading.value = true;

        const response = await fetch("/api/auth/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(loginForm),
        });

        const data = await response.json();

        if (response.ok) {
          localStorage.setItem("token", data.token);
          localStorage.setItem("user", JSON.stringify(data.user));
          localStorage.setItem("userRole", data.user.role);
          localStorage.setItem("username", data.user.username);

          ElMessage.success("登录成功");

          // 根据角色跳转到对应页面
          if (data.user.role === "teacher") {
            router.push("/teacher");
          } else {
            router.push("/parent");
          }
        } else {
          ElMessage.error(data.error || "登录失败");
        }
      } catch (error) {
        console.error("登录错误:", error);
        ElMessage.error("登录失败，请重试");
      } finally {
        loading.value = false;
      }
    };

    const goToRegister = () => {
      router.push("/register");
    };

    return {
      loginForm,
      rules,
      loginFormRef,
      loading,
      handleLogin,
      goToRegister,
    };
  },
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  border-radius: 15px;
  overflow: hidden;
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0;
  color: #409eff;
  font-size: 24px;
}

.card-header p {
  margin: 10px 0 0 0;
  color: #666;
  font-size: 14px;
}

:deep(.el-card__header) {
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}
</style>
