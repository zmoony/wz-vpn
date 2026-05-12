<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage } from "element-plus";

import PageHeader from "../components/PageHeader.vue";
import StatCard from "../components/StatCard.vue";
import { fetchSystemStats, type SystemStats } from "../api/sys";

const loading = ref(false);
const stats = ref<SystemStats | null>(null);

async function load() {
  loading.value = true;
  try {
    stats.value = await fetchSystemStats();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载系统状态失败");
  } finally {
    loading.value = false;
  }
}

onMounted(load);
</script>

<template>
  <div class="page-shell">
    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="系统总览"
          description="这里先接入系统状态和 WireGuard 主链路，其它模块保留为可扩展入口。"
        >
          <el-button :loading="loading" @click="load">刷新状态</el-button>
        </PageHeader>
      </div>
    </section>

    <section class="stats-grid" v-if="stats">
      <StatCard label="CPU" :value="`${stats.cpuPercent.toFixed(1)}%`" :helper="`负载 ${stats.load1} / ${stats.load5} / ${stats.load15}`" />
      <StatCard label="内存" :value="`${stats.memoryPercent.toFixed(1)}%`" :helper="`${stats.memoryUsedMb}MB / ${stats.memoryTotalMb}MB`" />
      <StatCard label="温度" :value="`${stats.temperatureC.toFixed(1)}°C`" helper="树莓派核心温度" />
      <StatCard label="磁盘" :value="`${stats.diskPercent.toFixed(1)}%`" :helper="`${stats.diskUsedGb.toFixed(1)}GB / ${stats.diskTotalGb.toFixed(1)}GB`" />
    </section>

    <section class="page-card" v-if="stats">
      <div class="page-card__body">
        <h3>网卡流量</h3>
        <el-table :data="stats.network">
          <el-table-column prop="name" label="网卡" />
          <el-table-column prop="rxBytes" label="接收字节" />
          <el-table-column prop="txBytes" label="发送字节" />
        </el-table>
      </div>
    </section>
  </div>
</template>
