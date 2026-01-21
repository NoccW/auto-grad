<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <div class="card-header">
          <h2>智能改卷助手</h2>
          <p>欢迎注册</p>
        </div>
      </template>

      <el-form
        :model="registerForm"
        :rules="rules"
        ref="registerFormRef"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="openId">
          <el-input v-model="registerForm.openId" placeholder="请输入用户名" />
        </el-form-item>

        <el-form-item label="姓名" prop="name">
          <el-input v-model="registerForm.name" placeholder="请输入真实姓名" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="registerForm.email" placeholder="请输入邮箱" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="用户类型" prop="userRole">
          <el-radio-group v-model="registerForm.userRole">
            <el-radio label="parent">
              <User style="margin-right: 8px" />
              家长用户
            </el-radio>
            <el-radio label="teacher">
              <School style="margin-right: 8px" />
              教师用户
            </el-radio>
          </el-radio-group>
          <div class="role-description">
            <span v-if="registerForm.userRole === 'parent'">
              家长用户可以上传试卷进行AI智能批改
            </span>
            <span v-else-if="registerForm.userRole === 'teacher'">
              教师用户可以创建自动化批改任务，批量处理试卷
            </span>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            @click="handleRegister"
            style="width: 100%"
          >
            注册
          </el-button>
        </el-form-item>

        <el-form-item>
          <el-button type="text" @click="goToLogin" style="width: 100%">
            已有账号？立即登录
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
import { User, School } from "@element-plus/icons-vue";

export default {
  name: "Register",
  components: {
    User,
    School,
  },
  setup() {
    const router = useRouter();
    const registerFormRef = ref();
    const loading = ref(false);

    const registerForm = reactive({
      openId: "",
      name: "",
      email: "",
      password: "",
      confirmPassword: "",
      userRole: "parent",
    });

    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== registerForm.password) {
        callback(new Error("两次输入的密码不一致"));
      } else {
        callback();
      }
    };

    const rules = {
      openId: [
        { required: true, message: "请输入用户名", trigger: "blur" },
        {
          min: 3,
          max: 20,
          message: "用户名长度在 3 到 20 个字符",
          trigger: "blur",
        },
      ],
      name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
      email: [
        { required: true, message: "请输入邮箱", trigger: "blur" },
        { type: "email", message: "请输入正确的邮箱地址", trigger: "blur" },
      ],
      password: [
        { required: true, message: "请输入密码", trigger: "blur" },
        { min: 6, message: "密码长度不能少于6位", trigger: "blur" },
      ],
      confirmPassword: [
        { required: true, message: "请确认密码", trigger: "blur" },
        { validator: validateConfirmPassword, trigger: "blur" },
      ],
      userRole: [
        { required: true, message: "请选择用户类型", trigger: "change" },
      ],
    };

    const handleRegister = async () => {
      if (!registerFormRef.value) return;

      try {
        await registerFormRef.value.validate();
        loading.value = true;

        const response = await fetch("/api/auth/register", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(registerForm),
        });

        const data = await response.json();

        if (response.ok && data?.user && data?.token) {
          const userRole = data.user.userRole || data.user.role || registerForm.userRole;

          localStorage.setItem("token", data.token);
          localStorage.setItem("user", JSON.stringify(data.user));
          localStorage.setItem("userRole", userRole);
          localStorage.setItem("username", data.user.username || registerForm.name);

          ElMessage.success("注册成功");

          // 根据角色跳转到对应页面
          if (userRole === "teacher") {
            router.push("/teacher");
          } else {
            router.push("/parent");
          }
        } else {
          ElMessage.error(data.error || "注册失败");
        }
      } catch (error) {
        console.error("注册错误:", error);
        ElMessage.error("注册失败，请重试");
      } finally {
        loading.value = false;
      }
    };

    const goToLogin = () => {
      router.push("/login");
    };

    return {
      registerForm,
      rules,
      registerFormRef,
      loading,
      handleRegister,
      goToLogin,
    };
  },
};
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px 0;
}

.register-card {
  width: 450px;
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

.role-description {
  margin-top: 10px;
  padding: 10px;
  background: #f0f9ff;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
  line-height: 1.4;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}
</style>
