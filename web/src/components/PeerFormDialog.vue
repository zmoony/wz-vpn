<script setup lang="ts">
import { reactive, watch } from "vue";

const props = defineProps<{
  modelValue: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [boolean];
  submit: [{ name: string; description: string }];
}>();

const form = reactive({
  name: "",
  description: "",
});

watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      form.name = "";
      form.description = "";
    }
  },
);

function close() {
  emit("update:modelValue", false);
}

function submit() {
  emit("submit", { name: form.name, description: form.description });
}
</script>

<template>
  <el-dialog :model-value="modelValue" title="新增 Peer" width="520px" @close="close">
    <el-form label-position="top">
      <el-form-item label="名称">
        <el-input v-model="form.name" placeholder="如 iphone-15" />
      </el-form-item>
      <el-form-item label="说明">
        <el-input v-model="form.description" type="textarea" :rows="3" placeholder="可选备注" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button type="primary" :loading="loading" @click="submit">创建</el-button>
    </template>
  </el-dialog>
</template>
