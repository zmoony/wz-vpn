<script setup lang="ts">
import PageHeader from "../components/PageHeader.vue";

defineProps<{
  module: string;
}>();
</script>

<template>
  <div class="page-shell">
    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          :title="module"
          description="当前页面骨架和路由已接好，后续会继续补真实配置表单、状态读取和系统适配层。"
        />
      </div>
    </section>

    <section class="page-card">
      <div class="page-card__body">
        <el-empty :description="`${module} 模块骨架已生成，等待下一轮实现。`" />
      </div>
    </section>
  </div>
</template>
