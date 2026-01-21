<template>
  <div class="task-config-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
          <h1>创建自动化任务</h1>
          <div></div>
        </div>
      </el-header>

      <el-main class="main-content">
        <el-card class="config-card">
          <template #header>
            <h3>任务配置</h3>
          </template>

          <el-form
            :model="taskForm"
            :rules="rules"
            ref="taskFormRef"
            label-width="120px"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="目标网站" prop="targetUrl">
                  <el-input
                    v-model="taskForm.targetUrl"
                    placeholder="https://www.7net.cc/build/home/index/index.html#"
                  />
                  <div class="form-tip">输入要自动批改的目标网站地址</div>
                </el-form-item>
              </el-col>

              <el-col :span="6">
                <el-form-item label="账号" prop="account">
                  <el-input
                    v-model="taskForm.account"
                    placeholder="请输入登录账号"
                  />
                </el-form-item>
              </el-col>

              <el-col :span="6">
                <el-form-item label="密码" prop="password">
                  <el-input
                    v-model="taskForm.password"
                    type="password"
                    placeholder="请输入登录密码"
                    show-password
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="8">
                <el-form-item label="试卷数量限制" prop="paperLimit">
                  <el-input-number
                    v-model="taskForm.paperLimit"
                    :min="1"
                    :max="10000"
                    placeholder="批改试卷数量上限"
                    style="width: 100%"
                  />
                  <div class="form-tip">最大支持10,000张试卷</div>
                </el-form-item>
              </el-col>

              <el-col :span="8">
                <el-form-item label="并发数量">
                  <el-input-number
                    v-model="taskForm.concurrent"
                    :min="1"
                    :max="5"
                    placeholder="同时处理的试卷数"
                    style="width: 100%"
                  />
                  <div class="form-tip">建议不超过3个并发</div>
                </el-form-item>
              </el-col>

              <el-col :span="8">
                <el-form-item label="超时时间">
                  <el-input-number
                    v-model="taskForm.timeout"
                    :min="30"
                    :max="300"
                    placeholder="单张试卷超时时间"
                    style="width: 100%"
                  />
                  <div class="form-tip">单位：秒，建议60-120秒</div>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="批改规则">
              <el-input
                v-model="taskForm.gradingRules"
                type="textarea"
                :rows="4"
                placeholder="请输入批改规则，例如：按题目给分，每题5分..."
              />
              <div class="form-tip">定义AI评分的标准和规则</div>
            </el-form-item>

            <el-form-item label="任务描述">
              <el-input
                v-model="taskForm.description"
                type="textarea"
                :rows="2"
                placeholder="任务描述（可选）"
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :loading="creating"
                @click="createTask"
              >
                创建并执行任务
              </el-button>
              <el-button size="large" @click="resetForm"> 重置表单 </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 高级配置 -->
        <el-card class="advanced-card">
          <template #header>
            <h3>高级配置</h3>
          </template>

          <el-collapse v-model="advancedExpanded">
            <el-collapse-item title="浏览器配置" name="browser">
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="浏览器类型">
                    <el-select v-model="taskForm.browserType">
                      <el-option label="Chrome" value="chrome" />
                      <el-option label="Firefox" value="firefox" />
                    </el-select>
                  </el-form-item>
                </el-col>

                <el-col :span="8">
                  <el-form-item label="无头模式">
                    <el-switch
                      v-model="taskForm.headless"
                      active-text="启用"
                      inactive-text="禁用"
                    />
                  </el-form-item>
                </el-col>

                <el-col :span="8">
                  <el-form-item label="调试模式">
                    <el-switch
                      v-model="taskForm.debug"
                      active-text="启用"
                      inactive-text="禁用"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-collapse-item>

            <el-collapse-item title="重试配置" name="retry">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="最大重试次数">
                    <el-input-number
                      v-model="taskForm.maxRetries"
                      :min="0"
                      :max="5"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>

                <el-col :span="12">
                  <el-form-item label="重试间隔">
                    <el-input-number
                      v-model="taskForm.retryDelay"
                      :min="5"
                      :max="60"
                      style="width: 100%"
                    />
                    <span style="margin-left: 10px">秒</span>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-collapse-item>

            <el-collapse-item title="通知配置" name="notification">
              <el-form-item label="任务完成通知">
                <el-switch
                  v-model="taskForm.notifyOnComplete"
                  active-text="启用"
                  inactive-text="禁用"
                />
              </el-form-item>

              <el-form-item label="失败通知">
                <el-switch
                  v-model="taskForm.notifyOnFailure"
                  active-text="启用"
                  inactive-text="禁用"
                />
              </el-form-item>
            </el-collapse-item>
          </el-collapse>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { ArrowLeft } from "@element-plus/icons-vue";

export default {
  name: "TaskConfig",
  components: {
    ArrowLeft,
  },
  setup() {
    const router = useRouter();
    const taskFormRef = ref();
    const creating = ref(false);
    const advancedExpanded = ref([]);

    const taskForm = reactive({
      targetUrl: "https://www.7net.cc/build/home/index/index.html#",
      account: "123123",
      password: "123123",
      paperLimit: 100,
      concurrent: 1,
      timeout: 60,
      gradingRules: "按题目分数给分，每题5分，总分100分",
      description: "",
      browserType: "chrome",
      headless: true,
      debug: false,
      maxRetries: 3,
      retryDelay: 10,
      notifyOnComplete: true,
      notifyOnFailure: true,
    });

    const rules = {
      targetUrl: [
        { required: true, message: "请输入目标网站地址", trigger: "blur" },
        { type: "url", message: "请输入有效的URL地址", trigger: "blur" },
      ],
      account: [{ required: true, message: "请输入账号", trigger: "blur" }],
      password: [
        { required: true, message: "请输入密码", trigger: "blur" },
        { min: 6, message: "密码长度不能少于6位", trigger: "blur" },
      ],
      paperLimit: [
        { required: true, message: "请输入试卷数量限制", trigger: "blur" },
        {
          type: "number",
          min: 1,
          max: 10000,
          message: "试卷数量限制在1-10000之间",
          trigger: "blur",
        },
      ],
      concurrent: [
        { required: true, message: "请输入并发数量", trigger: "blur" },
        {
          type: "number",
          min: 1,
          max: 5,
          message: "并发数量在1-5之间",
          trigger: "blur",
        },
      ],
      timeout: [
        { required: true, message: "请输入超时时间", trigger: "blur" },
        {
          type: "number",
          min: 30,
          max: 300,
          message: "超时时间在30-300秒之间",
          trigger: "blur",
        },
      ],
      gradingRules: [
        { required: true, message: "请输入批改规则", trigger: "blur" },
      ],
    };

    const createTask = async () => {
      if (!taskFormRef.value) return;

      try {
        await taskFormRef.value.validate();
        creating.value = true;

        const token = localStorage.getItem("token");

        // 创建任务
        const response = await fetch("/api/teacher/tasks", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(taskForm),
        });

        const data = await response.json();

        if (response.ok) {
          ElMessage.success("任务创建成功");

          // 直接执行任务
          const executeResponse = await fetch(
            `/api/teacher/tasks/${data.id}/execute`,
            {
              method: "POST",
              headers: {
                Authorization: `Bearer ${token}`,
              },
            },
          );

          if (executeResponse.ok) {
            ElMessage.success("任务已开始执行");
            router.push(`/task-monitor/${data.id}`);
          } else {
            const executeData = await executeResponse.json();
            ElMessage.error(`任务执行失败: ${executeData.error}`);
            router.push(`/task-monitor/${data.id}`);
          }
        } else {
          ElMessage.error(data.error || "创建任务失败");
        }
      } catch (error) {
        console.error("创建任务错误:", error);
        ElMessage.error("创建任务失败");
      } finally {
        creating.value = false;
      }
    };

    const resetForm = () => {
      if (taskFormRef.value) {
        taskFormRef.value.resetFields();
      }
      // 重置为默认值
      taskForm.targetUrl = "https://www.7net.cc/build/home/index/index.html#";
      taskForm.account = "123123";
      taskForm.password = "123123";
      taskForm.paperLimit = 100;
      taskForm.concurrent = 1;
      taskForm.timeout = 60;
      taskForm.gradingRules = "按题目分数给分，每题5分，总分100分";
    };

    const goBack = () => {
      router.push("/teacher-home");
    };

    return {
      taskForm,
      rules,
      taskFormRef,
      creating,
      advancedExpanded,
      createTask,
      resetForm,
      goBack,
    };
  },
};
</script>

<style scoped>
.task-config-container {
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

.config-card {
  margin-bottom: 20px;
}

.advanced-card {
  margin-bottom: 20px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  line-height: 1.4;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-collapse-item__header) {
  font-weight: 500;
}

:deep(.el-collapse-item__content) {
  padding-bottom: 0;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style>
