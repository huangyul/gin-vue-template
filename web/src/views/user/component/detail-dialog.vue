<template>
  <el-dialog v-model="dialogVisible" :title width="500">
    <el-form ref="formRef" :model="form" :rules>
      <el-form-item label="用户名" label-width="140px" prop="username">
        <el-input v-model="form.username" autocomplete="off" />
      </el-form-item>
      <el-form-item label="密码" label-width="140px" prop="password">
        <el-input v-model="form.password" type="password" autocomplete="off" />
      </el-form-item>
      <el-form-item label="昵称" label-width="140px">
        <el-input v-model="form.nickname" autocomplete="off" />
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
import { ElMessage, type FormInstance } from "element-plus";
import { ref } from "vue";

const emits = defineEmits(["refresh"]);
const dialogVisible = ref(false);
const title = ref("");
const form = ref({
  username: "",
  password: "",
  nickname: "",
  id: 0
});
const rules = {
  username: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 3, max: 15, message: "长度在 3 到 15 个字符", trigger: "blur" }
  ],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, max: 10, message: "长度在 6 到 10 个字符", trigger: "blur" }
  ]
};
const formRef = ref<FormInstance>();

const handleClose = () => {
  dialogVisible.value = false;
};
const handleOpen = async (id?: number) => {
  if (id) {
    // 编辑
    title.value = "编辑用户";
    const res = await getUserInfo(id);
    form.value.id = res.id;
    form.value.username = res.username;
    form.value.nickname = res.nickname;
    form.value.password = "******";
  } else {
    // 新增
    form.value = {
      username: "",
      password: "",
      nickname: "",
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
    await createUser({
      username: form.value.username,
      password: form.value.password,
      nickname: form.value.nickname
    });
    ElMessage.success("创建成功");
  }
  emits("refresh");
  handleClose();
};

defineExpose({
  handleOpen,
  handleClose
});
</script>
