<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "../lib/element-plus";

import PageHeader from "../components/PageHeader.vue";
import StatCard from "../components/StatCard.vue";
import { fetchSystemStats, type SystemStats } from "../api/sys";

const loading = ref(false);
const stats = ref<SystemStats | null>(null);

const highlights = computed(() => {
  if (!stats.value) {
    return [];
  }
  return [
    {
      label: "运行温度",
      value: `${stats.value.temperatureC.toFixed(1)}°C`,
      helper: stats.value.temperatureC >= 70 ? "建议观察散热状态" : "温度处于可接受区间",
    },
    {
      label: "内存压力",
      value: `${stats.value.memoryPercent.toFixed(1)}%`,
      helper: `${stats.value.memoryUsedMb}MB / ${stats.value.memoryTotalMb}MB`,
    },
    {
      label: "磁盘占用",
      value: `${stats.value.diskPercent.toFixed(1)}%`,
      helper: `${stats.value.diskUsedGb.toFixed(1)}GB / ${stats.value.diskTotalGb.toFixed(1)}GB`,
    },
    {
      label: "系统负载",
      value: `${stats.value.load1}`,
      helper: `1 / 5 / 15 分钟：${stats.value.load1} / ${stats.value.load5} / ${stats.value.load15}`,
    },
  ];
});

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
  <div class="page-shell dashboard-page">
    <section class="page-card dashboard-hero">
      <div class="page-card__body dashboard-hero__body">
        <PageHeader
          title="系统总览"
          description="把当前网关主机的运行状态、资源占用和网络接口放到一块首页画布里，方便你在进入 WireGuard、反代和防火墙前先判断整体健康度。"
        >
          <el-button :loading="loading" @click="load">刷新状态</el-button>
        </PageHeader>

        <div v-if="stats" class="dashboard-hero__summary">
          <div class="dashboard-hero__badge">
            <strong>{{ stats.cpuPercent.toFixed(1) }}%</strong>
            <span>当前 CPU 占用</span>
          </div>
          <p>
            这台网关主机正在稳定运行，建议优先关注温度、磁盘和负载的持续变化，再决定是否继续调整 VPN、
            防火墙或反向代理配置。
          </p>
        </div>
      </div>
    </section>

    <section v-if="stats" class="stats-grid">
      <StatCard
        v-for="item in highlights"
        :key="item.label"
        :label="item.label"
        :value="item.value"
        :helper="item.helper"
      />
    </section>

    <section v-else class="page-card">
      <div class="page-card__body">
        <div class="empty-panel">
          <strong>还没有加载到系统状态</strong>
          <p>你可以点击右上角“刷新状态”重新拉取。如果后端暂时不可用，这里会保持为空而不是展示过期数据。</p>
        </div>
      </div>
    </section>

    <section class="dashboard-grid" v-if="stats">
      <article class="page-card">
        <div class="page-card__body">
          <h3 class="section-title">运行摘要</h3>
          <div class="dashboard-summary">
            <div class="surface-muted dashboard-summary__item">
              <span>CPU</span>
              <strong>{{ stats.cpuPercent.toFixed(1) }}%</strong>
            </div>
            <div class="surface-muted dashboard-summary__item">
              <span>温度</span>
              <strong>{{ stats.temperatureC.toFixed(1) }}°C</strong>
            </div>
            <div class="surface-muted dashboard-summary__item">
              <span>磁盘</span>
              <strong>{{ stats.diskPercent.toFixed(1) }}%</strong>
            </div>
          </div>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <h3 class="section-title">接口流量</h3>
          <p class="section-subtitle">持续观察收发字节变化，有助于判断 WireGuard、反代和 DDNS 是否正在活跃工作。</p>
          <el-table :data="stats.network">
            <el-table-column prop="name" label="网卡" />
            <el-table-column prop="rxBytes" label="接收字节" />
            <el-table-column prop="txBytes" label="发送字节" />
          </el-table>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.dashboard-hero__body {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.dashboard-hero {
  background:
    radial-gradient(circle at 92% 12%, rgba(245, 184, 65, 0.22), transparent 22%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(245, 249, 248, 0.96));
}

.dashboard-hero__summary {
  display: grid;
  grid-template-columns: 190px 1fr;
  gap: 18px;
  align-items: center;
}

.dashboard-hero__summary p {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.8;
}

.dashboard-hero__badge {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 20px;
  border-radius: 18px;
  color: #ffffff;
  background: linear-gradient(160deg, var(--brand) 0%, var(--brand-deep) 100%);
  box-shadow: 0 16px 28px rgba(22, 152, 142, 0.22);
}

.dashboard-hero__badge strong {
  font-size: 34px;
  line-height: 1.1;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 20px;
}

.dashboard-summary {
  display: grid;
  gap: 14px;
}

.dashboard-summary__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 18px;
}

.dashboard-summary__item span {
  color: var(--text-secondary);
}

.dashboard-summary__item strong {
  color: var(--text);
  font-size: 24px;
}

@media (max-width: 1080px) {
  .dashboard-hero__summary,
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}
</style>
