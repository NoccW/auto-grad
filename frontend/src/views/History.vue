<template>
  <div class="history-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
          <h1>历史记录</h1>
          <div></div>
        </div>
      </el-header>

      <el-main class="main-content">
        <!-- 搜索和筛选 -->
        <el-card class="filter-card">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-input
                v-model="searchForm.keyword"
                placeholder="搜索关键词"
                :prefix-icon="Search"
                @input="handleSearch"
              />
            </el-col>
            <el-col :span="4">
              <el-select
                v-model="searchForm.status"
                placeholder="状态筛选"
                clearable
                @change="handleFilter"
              >
                <el-option label="全部" value="" />
                <el-option label="待处理" value="pending" />
                <el-option label="处理中" value="processing" />
                <el-option label="已完成" value="completed" />
                <el-option label="失败" value="failed" />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-date-picker
                v-model="searchForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                @change="handleFilter"
              />
            </el-col>
            <el-col :span="4">
              <el-button type="primary" @click="refreshData">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </el-col>
          </el-row>
        </el-card>

        <!-- 数据表格 -->
        <el-card class="table-card">
          <el-table
            :data="gradingRecords"
            v-loading="loading"
            stripe
            style="width: 100%"
          >
            <el-table-column type="expand">
              <template #default="props">
                <div class="expand-content">
                  <el-descriptions :column="2" border>
                    <el-descriptions-item label="记录ID">
                      {{ props.row.id }}
                    </el-descriptions-item>
                    <el-descriptions-item label="状态">
                      <el-tag :type="getStatusType(props.row.status)">
                        {{ getStatusText(props.row.status) }}
                      </el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="得分">
                      <span v-if="props.row.aiScore"
                        >{{ props.row.aiScore }}分</span
                      >
                      <span v-else>-</span>
                    </el-descriptions-item>
                    <el-descriptions-item label="创建时间">
                      {{ formatDate(props.row.createdAt) }}
                    </el-descriptions-item>
                    <el-descriptions-item
                      label="更新时间"
                      v-if="props.row.updatedAt !== props.row.createdAt"
                    >
                      {{ formatDate(props.row.updatedAt) }}
                    </el-descriptions-item>
                    <el-descriptions-item
                      label="错题数量"
                      v-if="props.row.wrongQuestions"
                    >
                      {{
                        parsedWrongQuestions(props.row.wrongQuestions).length
                      }}道
                    </el-descriptions-item>
                  </el-descriptions>

                  <!-- OCR识别结果预览 -->
                  <div v-if="props.row.ocrResult" class="ocr-preview">
                    <h4>OCR识别结果预览</h4>
                    <el-input
                      :model-value="props.row.ocrResult"
                      type="textarea"
                      :rows="3"
                      readonly
                    />
                  </div>
                </div>
              </template>
            </el-table-column>

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
                <span v-if="scope.row.aiScore">{{ scope.row.aiScore }}分</span>
                <span v-else>-</span>
              </template>
            </el-table-column>

            <el-table-column prop="createdAt" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.createdAt) }}
              </template>
            </el-table-column>

            <el-table-column label="错题数量" width="100">
              <template #default="scope">
                <span v-if="scope.row.wrongQuestions">
                  {{ parsedWrongQuestions(scope.row.wrongQuestions).length }}道
                </span>
                <span v-else>-</span>
              </template>
            </el-table-column>

            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button
                  type="primary"
                  size="small"
                  @click="viewDetail(scope.row)"
                  :disabled="scope.row.status !== 'completed'"
                >
                  查看详情
                </el-button>
                <el-button
                  type="warning"
                  size="small"
                  @click="regrade(scope.row)"
                  :disabled="scope.row.status === 'processing'"
                >
                  重新批改
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="deleteRecord(scope.row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { ArrowLeft, Search, Refresh } from "@element-plus/icons-vue";

export default {
  name: "History",
  components: {
    ArrowLeft,
    Search,
    Refresh,
  },
  setup() {
    const router = useRouter();
    const loading = ref(false);
    const gradingRecords = ref([]);
    const total = ref(0);
    const currentPage = ref(1);
    const pageSize = ref(10);

    const searchForm = reactive({
      keyword: "",
      status: "",
      dateRange: null,
    });

    const loadRecords = async () => {
      loading.value = true;
      try {
        const token = localStorage.getItem("token");

        // 构建查询参数
        const params = new URLSearchParams({
          page: currentPage.value,
          limit: pageSize.value,
        });

        if (searchForm.status) {
          params.append("status", searchForm.status);
        }

        if (searchForm.keyword) {
          params.append("keyword", searchForm.keyword);
        }

        if (searchForm.dateRange && searchForm.dateRange.length === 2) {
          params.append("startDate", searchForm.dateRange[0]);
          params.append("endDate", searchForm.dateRange[1]);
        }

        const response = await fetch(`/api/grading?${params}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          gradingRecords.value = data.records || [];
          total.value = data.total || 0;
        } else {
          const errorData = await response.json();
          ElMessage.error(errorData.error || "加载记录失败");
        }
      } catch (error) {
        console.error("加载记录错误:", error);
        ElMessage.error("加载记录失败");
      } finally {
        loading.value = false;
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

    const parsedWrongQuestions = (wrongQuestionsStr) => {
      if (!wrongQuestionsStr) return [];
      try {
        const parsed = JSON.parse(wrongQuestionsStr);
        return parsed.questions || [];
      } catch {
        return [];
      }
    };

    const handleSearch = () => {
      currentPage.value = 1;
      loadRecords();
    };

    const handleFilter = () => {
      currentPage.value = 1;
      loadRecords();
    };

    const refreshData = () => {
      loadRecords();
    };

    const handleSizeChange = (newSize) => {
      pageSize.value = newSize;
      currentPage.value = 1;
      loadRecords();
    };

    const handleCurrentChange = (newPage) => {
      currentPage.value = newPage;
      loadRecords();
    };

    const viewDetail = (record) => {
      router.push(`/result/${record.id}`);
    };

    const regrade = async (record) => {
      try {
        loading.value = true;
        const token = localStorage.getItem("token");

        const response = await fetch(`/api/grading/${record.id}/process`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          ElMessage.success("重新批改完成");
          await loadRecords();
        } else {
          const errorData = await response.json();
          ElMessage.error(errorData.error || "重新批改失败");
        }
      } catch (error) {
        console.error("重新批改错误:", error);
        ElMessage.error("重新批改失败");
      } finally {
        loading.value = false;
      }
    };

    const deleteRecord = async (record) => {
      try {
        await ElMessageBox.confirm(
          `确定要删除这条批改记录吗？（ID: ${record.id}）`,
          "删除确认",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          },
        );

        loading.value = true;
        const token = localStorage.getItem("token");

        const response = await fetch(`/api/grading/${record.id}`, {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          ElMessage.success("删除成功");
          await loadRecords();
        } else {
          const errorData = await response.json();
          ElMessage.error(errorData.error || "删除失败");
        }
      } catch (error) {
        if (error !== "cancel") {
          console.error("删除错误:", error);
          ElMessage.error("删除失败");
        }
      } finally {
        loading.value = false;
      }
    };

    const goBack = () => {
      router.push("/");
    };

    onMounted(() => {
      loadRecords();
    });

    return {
      loading,
      gradingRecords,
      total,
      currentPage,
      pageSize,
      searchForm,
      getStatusType,
      getStatusText,
      formatDate,
      parsedWrongQuestions,
      handleSearch,
      handleFilter,
      refreshData,
      handleSizeChange,
      handleCurrentChange,
      viewDetail,
      regrade,
      deleteRecord,
      goBack,
    };
  },
};
</script>

<style scoped>
.history-container {
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

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  min-height: 400px;
}

.expand-content {
  padding: 20px;
  background: #fafafa;
  border-radius: 4px;
}

.ocr-preview {
  margin-top: 20px;
}

.ocr-preview h4 {
  margin: 0 0 10px 0;
  color: #303133;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-table__expanded-cell) {
  padding: 0;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
}
</style>
