<template>
  <div class="upload-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
          <h1>上传试卷</h1>
          <div></div>
        </div>
      </el-header>

      <el-main class="main-content">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-card class="upload-card">
              <template #header>
                <h3>上传试卷图片</h3>
              </template>

              <div class="upload-area" v-loading="uploading">
                <el-upload
                  ref="uploadRef"
                  class="upload-dragger"
                  drag
                  :action="uploadUrl"
                  :headers="uploadHeaders"
                  :on-success="handleUploadSuccess"
                  :on-error="handleUploadError"
                  :before-upload="beforeUpload"
                  accept="image/*"
                >
                  <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                  <div class="el-upload__text">
                    将图片拖拽到此处，或<em>点击上传</em>
                  </div>
                  <template #tip>
                    <div class="el-upload__tip">
                      支持 jpg/png/gif 格式，文件大小不超过 10MB
                    </div>
                  </template>
                </el-upload>
              </div>

              <div class="preview-area" v-if="paperImageUrl">
                <h4>试卷预览</h4>
                <el-image
                  :src="getImageUrl(paperImageUrl)"
                  fit="contain"
                  style="width: 100%; max-height: 300px"
                />
              </div>
            </el-card>
          </el-col>

          <el-col :span="12">
            <el-card class="reference-card">
              <template #header>
                <h3>参考答案（可选）</h3>
              </template>

              <div class="upload-area" v-loading="uploadingAnswer">
                <el-upload
                  ref="answerUploadRef"
                  class="upload-dragger"
                  drag
                  :action="uploadUrl"
                  :headers="uploadHeaders"
                  :on-success="handleAnswerUploadSuccess"
                  :on-error="handleUploadError"
                  :before-upload="beforeUpload"
                  accept="image/*"
                >
                  <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                  <div class="el-upload__text">
                    将参考答案图片拖拽到此处，或<em>点击上传</em>
                  </div>
                  <template #tip>
                    <div class="el-upload__tip">
                      上传标准答案图片，有助于提高批改准确性
                    </div>
                  </template>
                </el-upload>
              </div>

              <div class="preview-area" v-if="answerImageUrl">
                <h4>参考答案预览</h4>
                <el-image
                  :src="getImageUrl(answerImageUrl)"
                  fit="contain"
                  style="width: 100%; max-height: 300px"
                />
              </div>
            </el-card>
          </el-col>
        </el-row>

        <div class="action-section" v-if="paperImageUrl">
          <el-card>
            <template #header>
              <h3>批改设置</h3>
            </template>

            <el-form :model="gradingForm" label-width="120px">
              <el-form-item label="学科类型">
                <el-select
                  v-model="gradingForm.subject"
                  placeholder="请选择学科"
                >
                  <el-option label="数学" value="math" />
                  <el-option label="语文" value="chinese" />
                  <el-option label="英语" value="english" />
                  <el-option label="物理" value="physics" />
                  <el-option label="化学" value="chemistry" />
                  <el-option label="生物" value="biology" />
                  <el-option label="其他" value="other" />
                </el-select>
              </el-form-item>

              <el-form-item label="年级">
                <el-select v-model="gradingForm.grade" placeholder="请选择年级">
                  <el-option label="小学一年级" value="grade1" />
                  <el-option label="小学二年级" value="grade2" />
                  <el-option label="小学三年级" value="grade3" />
                  <el-option label="小学四年级" value="grade4" />
                  <el-option label="小学五年级" value="grade5" />
                  <el-option label="小学六年级" value="grade6" />
                  <el-option label="初中一年级" value="grade7" />
                  <el-option label="初中二年级" value="grade8" />
                  <el-option label="初中三年级" value="grade9" />
                  <el-option label="高中一年级" value="grade10" />
                  <el-option label="高中二年级" value="grade11" />
                  <el-option label="高中三年级" value="grade12" />
                </el-select>
              </el-form-item>

              <el-form-item>
                <el-button
                  type="primary"
                  size="large"
                  :loading="processing"
                  @click="startGrading"
                >
                  开始智能批改
                </el-button>
                <el-button size="large" @click="resetForm">
                  重新上传
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { ArrowLeft, UploadFilled } from "@element-plus/icons-vue";

export default {
  name: "Upload",
  components: {
    ArrowLeft,
    UploadFilled,
  },
  setup() {
    const router = useRouter();
    const uploadRef = ref();
    const answerUploadRef = ref();
    const uploading = ref(false);
    const uploadingAnswer = ref(false);
    const processing = ref(false);

    const paperImageUrl = ref("");
    const answerImageUrl = ref("");

    const gradingForm = reactive({
      subject: "",
      grade: "",
    });

    const uploadUrl = "/api/upload";
    const uploadHeaders = {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    };

    const beforeUpload = (file) => {
      const isImage = file.type.startsWith("image/");
      const isLt10M = file.size / 1024 / 1024 < 10;

      if (!isImage) {
        ElMessage.error("只能上传图片文件!");
        return false;
      }
      if (!isLt10M) {
        ElMessage.error("图片大小不能超过 10MB!");
        return false;
      }
      return true;
    };

    const handleUploadSuccess = (response) => {
      uploading.value = false;
      paperImageUrl.value = response.relativePath || response.url || response.filename;
      ElMessage.success("试卷图片上传成功");
    };

    const handleAnswerUploadSuccess = (response) => {
      uploadingAnswer.value = false;
      answerImageUrl.value =
        response.relativePath || response.url || response.filename;
      ElMessage.success("参考答案上传成功");
    };

    const handleUploadError = (error) => {
      uploading.value = false;
      uploadingAnswer.value = false;
      ElMessage.error("图片上传失败");
      console.error("上传错误:", error);
    };

    const getImageUrl = (relativePath) => {
      if (!relativePath) return "";
      // 兼容后端返回的 url 或文件名
      if (relativePath.startsWith("http")) return relativePath;
      if (relativePath.startsWith("/")) return relativePath;
      return `/uploads/${relativePath}`;
    };

    const startGrading = async () => {
      if (!paperImageUrl.value) {
        ElMessage.error("请先上传试卷图片");
        return;
      }

      processing.value = true;

      try {
        const token = localStorage.getItem("token");

        // 创建批改记录（后端可异步接入真实评分流程）
        const createResponse = await fetch("/api/grading", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            paperImageUrl: paperImageUrl.value,
            answerImageUrl: answerImageUrl.value || undefined,
          }),
        });

        const createData = await createResponse.json();

        if (!createResponse.ok) {
          throw new Error(createData.error || "创建批改记录失败");
        }

        ElMessage.success("批改任务已提交，等待处理");
        router.push(`/result/${createData.id}`);
      } catch (error) {
        console.error("批改错误:", error);
        ElMessage.error("批改失败，请重试");
      } finally {
        processing.value = false;
      }
    };

    const resetForm = () => {
      paperImageUrl.value = "";
      answerImageUrl.value = "";
      gradingForm.subject = "";
      gradingForm.grade = "";
      if (uploadRef.value) {
        uploadRef.value.clearFiles();
      }
      if (answerUploadRef.value) {
        answerUploadRef.value.clearFiles();
      }
    };

    const goBack = () => {
      router.push("/");
    };

    return {
      uploadRef,
      answerUploadRef,
      uploading,
      uploadingAnswer,
      processing,
      paperImageUrl,
      answerImageUrl,
      gradingForm,
      uploadUrl,
      uploadHeaders,
      beforeUpload,
      handleUploadSuccess,
      handleAnswerUploadSuccess,
      handleUploadError,
      getImageUrl,
      startGrading,
      resetForm,
      goBack,
    };
  },
};
</script>

<style scoped>
.upload-container {
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
  font-size: 20px;
}

.main-content {
  padding: 20px;
}

.upload-card,
.reference-card {
  margin-bottom: 20px;
  min-height: 500px;
}

.upload-area {
  margin-bottom: 20px;
}

.upload-dragger {
  width: 100%;
}

.preview-area {
  text-align: center;
}

.preview-area h4 {
  margin: 20px 0 10px 0;
  color: #303133;
}

.action-section {
  margin-top: 20px;
}

:deep(.el-upload-dragger) {
  width: 100%;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}
</style>
