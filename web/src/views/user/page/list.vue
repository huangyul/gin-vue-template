<template>
  <div class="flex flex-col h-full">
    <el-card class="mb-2">
      <el-button type="success" @click="handleCreate">新建</el-button>
      <el-form :inline="true" class="mt-2">
        <el-form-item label="昵称">
          <el-input v-model="searchForm.nickname" />
        </el-form-item>
        <el-form-item label="账号">
          <el-input v-model="searchForm.username" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <div class="flex-1 mb-2 overflow-hidden">
      <el-table :data="tableData" border style="height: 100%; overflow: auto">
        <el-table-column prop="username" label="账号" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="created_at" label="创建时间" />
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row.id)">
              Edit
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)">
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="flex justify-end">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 30, 40]"
        background
        size="small"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSearch"
        @current-change="handleSearch"
      />
    </div>

    <DetailDialog ref="detailRef" @refresh="handleSearch" />
  </div>
</template>

<script setup lang="ts">
import { deleteUser, getUserList } from "@/api/user";
import { User } from "@/types/user";
import { onMounted, ref } from "vue";
import DetailDialog from "../component/detail-dialog.vue";
import { ElMessage } from "element-plus";

const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(100);

const tableData = ref<User[]>([]);

const searchForm = ref({
  nickname: "",
  username: ""
});

const handleReset = () => {
  searchForm.value = {
    nickname: "",
    username: ""
  };
  handleSearch();
};

const handleSearch = async () => {
  const res = await getUserList({
    page: currentPage.value,
    page_size: pageSize.value,
    nickname: searchForm.value.nickname,
    username: searchForm.value.username
  });
  tableData.value = res.data;
  total.value = res.total;
};

const detailRef = ref<InstanceType<typeof DetailDialog | null>>(null);

const handleEdit = (id: number) => {
  detailRef.value?.handleOpen(id);
};

const handleCreate = () => {
  detailRef.value?.handleOpen();
};
const handleDelete = async (id: number) => {
  await deleteUser(id);
  ElMessage.success("删除成功");
  handleSearch();
};

onMounted(() => {
  handleSearch();
});
</script>

<style lang="scss" scoped>
/* 确保页面容器高度生效 */
html,
body,
#app {
  height: 100%;
  margin: 0;
}
</style>
