<template>
  <div class="flex flex-col h-full">
    <el-card class="mb-2">
      <el-upload
        ref="uploadRef"
        action="/api/file/upload"
        :show-file-list="false"
        :on-success="handleUpload"
        :headers="{
          Authorization: `Bearer ${getToken().accessToken}`
        }"
      >
        <template #trigger>
          <el-button type="primary">select file</el-button>
        </template>
        <template #tip>
          <div class="el-upload__tip">
            jpg/png files with a size less than 500kb
          </div>
        </template>
      </el-upload>
      <el-form :inline="true" class="mt-2">
        <el-form-item label="文件名">
          <el-input v-model="searchForm.filename" />
        </el-form-item>
        <el-form-item label="上传用户">
          <el-select v-model="searchForm.userId" style="width: 192px" clearable>
            <el-option
              v-for="item in searchOptions.userOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <div class="flex-1 mb-2 overflow-hidden">
      <el-table :data="tableData" border style="height: 100%; overflow: auto">
        <el-table-column prop="file_name" label="文件名称" />
        <el-table-column prop="upload_user" label="上传用户" />
        <el-table-column prop="upload_time" label="创建时间" />
        <el-table-column label="操作">
          <template #default="{ row }">
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
  </div>
</template>

<script setup lang="ts">
import { deleteUser, getUserList } from "@/api/user";
import { deleteFileById, getList, getOptions, uploadFile } from "@/api/file";
import type { File } from "@/types/file";
import { onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { SelectOptions } from "@/types/common";
import type { UploadInstance, UploadRawFile } from "element-plus";
import { getToken } from "@/utils/auth";
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(100);

const tableData = ref<
  {
    file_name: string;
    link: string;
    upload_user: number;
    upload_time: string;
    id: number;
  }[]
>([]);

const searchForm = ref({
  filename: "",
  userId: ""
});

const searchOptions = ref<{
  userOptions: SelectOptions;
}>({
  userOptions: []
});

const handleReset = () => {
  searchForm.value = {
    filename: "",
    userId: ""
  };
  handleSearch();
};

const handleSearch = async () => {
  const res = await getList({
    page: currentPage.value,
    page_size: pageSize.value,
    file_name: searchForm.value.filename,
    user_id: searchForm.value.userId
  });
  tableData.value = res.data;
  total.value = res.total;
};

const handleCreate = () => {};
const handleDelete = async (id: number) => {
  await deleteFileById(id);
  ElMessage.success("删除成功");
  handleSearch();
};

const uploadRef = ref<UploadInstance>();

const handleUpload = () => {
  ElMessage.success("上传成功");
};

const init = async () => {
  const res = await getOptions();
  searchOptions.value.userOptions = res.user;
};

onMounted(() => {
  init();
  handleSearch();
});
</script>

<style lang="scss" scoped>
html,
body,
#app {
  height: 100%;
  margin: 0;
}
</style>
