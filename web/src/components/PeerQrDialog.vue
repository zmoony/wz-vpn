<script setup lang="ts">
defineProps<{
  modelValue: boolean;
  conf: string;
  qrCode: string;
}>();

const emit = defineEmits<{
  "update:modelValue": [boolean];
}>();

function close() {
  emit("update:modelValue", false);
}
</script>

<template>
  <el-dialog :model-value="modelValue" title="客户端配置" width="760px" @close="close">
    <div class="qr-dialog">
      <section>
        <h4>二维码输出</h4>
        <pre>{{ qrCode }}</pre>
      </section>
      <section>
        <h4>配置文件</h4>
        <pre>{{ conf }}</pre>
      </section>
    </div>
  </el-dialog>
</template>

<style scoped>
.qr-dialog {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

pre {
  min-height: 240px;
  margin: 0;
  padding: 16px;
  overflow: auto;
  white-space: pre-wrap;
  border-radius: 12px;
  background: #0f172a;
  color: #ecf3ff;
}

@media (max-width: 900px) {
  .qr-dialog {
    grid-template-columns: 1fr;
  }
}
</style>
