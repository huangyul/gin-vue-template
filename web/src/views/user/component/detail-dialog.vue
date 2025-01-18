<template>
  <el-dialog v-model="dialogVisible" :title width="500">
    <el-form ref="formRef" :model="form" :rules>
      <el-form-item label="用户名" label-width="140px" prop="username">
        <el-input
          v-model="form.username"
          :disabled="type == 'edit'"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item label="密码" label-width="140px" prop="password">
        <el-input
          v-model="form.password"
          :disabled="type == 'edit'"
          type="password"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        v-if="type == 'create'"
        label="确认密码"
        label-width="140px"
        prop="confirmPassword"
      >
        <el-input
          v-model="form.confirmPassword"
          type="password"
          autocomplete="off"
        />
      </el-form-item>

      <el-form-item label="昵称" label-width="140px">
        <el-input v-model="form.nickname" autocomplete="off" />
      </el-form-item>
      <el-form-item label="头像" label-width="140px">
        <el-upload
          action="/api/file/upload"
          :show-file-list="false"
          :on-success="handleAvatarSuccess"
          :before-upload="beforeAvatarUpload"
          :headers="{
            Authorization: `Bearer ${getToken().accessToken}`
          }"
        >
          <img v-if="form.avatar" :src="form.avatar" class="avatar" />
          <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
        </el-upload>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave"> 保存 </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { createUser, getUserInfo, updateUser } from "@/api/user";
import { ElMessage, type FormInstance, type UploadProps } from "element-plus";
import { ref } from "vue";
import { getToken } from "@/utils/auth";
import { Plus } from "@element-plus/icons-vue";
import type { ApiResponse } from "@/types/api";

const emits = defineEmits(["refresh"]);
const dialogVisible = ref(false);
const title = ref("");
const form = ref({
  username: "",
  password: "",
  confirmPassword: "",
  nickname: "",
  avatar: "",
  id: 0
});
const type = ref<"create" | "edit">("create");
const rules = {
  username: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 3, max: 15, message: "长度在 3 到 15 个字符", trigger: "blur" }
  ],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, max: 10, message: "长度在 6 到 10 个字符", trigger: "blur" }
  ],
  confirmPassword: [
    { required: true, message: "请再次输入密码", trigger: "blur" },
    { min: 6, max: 10, message: "长度在 6 到 10 个字符", trigger: "blur" }
  ]
};
const formRef = ref<FormInstance>();

const handleClose = () => {
  dialogVisible.value = false;
};
const handleOpen = async (id?: number) => {
  if (id) {
    type.value = "edit";
    // 编辑
    title.value = "编辑用户";
    const res = await getUserInfo(id);
    form.value.id = res.id;
    form.value.username = res.username;
    form.value.nickname = res.nickname;
    form.value.password = "******";
    form.value.avatar = res.avatar;
  } else {
    type.value = "create";
    // 新增
    form.value = {
      username: "",
      password: "",
      nickname: "",
      confirmPassword: "",
      avatar: "",
      id: 0
    };
  }
  dialogVisible.value = true;
};

const handleSave = async () => {
  if (form.value.id) {
    await updateUser({
      id: form.value.id,
      nickname: form.value.nickname
    });
    ElMessage.success("更新成功");
  } else {
    await formRef.value?.validate();
    if (form.value.password !== form.value.confirmPassword) {
      ElMessage.error("两次密码不一致");
      return;
    }
    await createUser({
      username: form.value.username,
      password: form.value.password,
      nickname: form.value.nickname,
      avatar: form.value.avatar
    });
    ElMessage.success("创建成功");
  }
  emits("refresh");
  handleClose();
};

const imageUrl = ref("");

const handleAvatarSuccess: UploadProps["onSuccess"] = (
  response: ApiResponse<{ url: string }>
) => {
  form.value.avatar = response.data.url;
};

const beforeAvatarUpload: UploadProps["beforeUpload"] = rawFile => {
  const validTypes = ["image/jpeg", "image/png"];
  if (!validTypes.includes(rawFile.type)) {
    ElMessage.error("头像图片必须是 JPG 或 PNG 格式！");
    return false;
  } else if (rawFile.size / 1024 / 1024 > 2) {
    ElMessage.error("头像图片大小不能超过 2MB！");
    return false;
  }
  return true;
};

defineExpose({
  handleOpen,
  handleClose
});
</script>

<style scoped lang="scss">
.avatar-uploader .el-upload {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}

.avatar-uploader .avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style>
