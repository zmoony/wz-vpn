<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage } from "../lib/element-plus";

import {
  createDDNSConfig,
  fetchDDNSConfigs,
  syncDDNSConfigs,
  updateDDNSConfig,
  type DDNSConfig,
  type DDNSPayload,
  type DDNSRuntimeStatus,
} from "../api/ddns";
import PageHeader from "../components/PageHeader.vue";

const loading = ref(false);
const saving = ref(false);
const syncing = ref(false);
const editingId = ref<number | null>(null);
const items = ref<DDNSConfig[]>([]);
const runtimeStatus = ref<DDNSRuntimeStatus | null>(null);

const form = reactive<DDNSPayload>({
  provider: "aliyun",
  accessKeyId: "",
  accessKeySecret: "",
  domain: "",
  subdomain: "home",
  enabled: true,
  checkIntervalSeconds: 300,
});

const actionLabel = computed(() => (editingId.value ? "更新 DDNS 配置" : "新增 DDNS 配置"));
const enabledCount = computed(() => items.value.filter((item) => item.enabled).length);

function resetForm() {
  editingId.value = null;
  form.provider = "aliyun";
  form.accessKeyId = "";
  form.accessKeySecret = "";
  form.domain = "";
  form.subdomain = "home";
  form.enabled = true;
  form.checkIntervalSeconds = 300;
}

async function load() {
  loading.value = true;
  try {
    const response = await fetchDDNSConfigs();
    items.value = response.items;
    runtimeStatus.value = response.runtimeStatus;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载 DDNS 配置失败");
  } finally {
    loading.value = false;
  }
}

function startEdit(item: DDNSConfig) {
  editingId.value = item.id;
  form.provider = item.provider;
  form.accessKeyId = item.accessKeyId;
  form.accessKeySecret = "";
  form.domain = item.domain;
  form.subdomain = item.subdomain;
  form.enabled = item.enabled;
  form.checkIntervalSeconds = item.checkIntervalSeconds;
}

async function submit() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateDDNSConfig(editingId.value, { ...form });
      ElMessage.success("DDNS 配置已更新");
    } else {
      await createDDNSConfig({ ...form });
      ElMessage.success("DDNS 配置已创建");
    }
    resetForm();
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存 DDNS 配置失败");
  } finally {
    saving.value = false;
  }
}

async function syncNow() {
  syncing.value = true;
  try {
    await syncDDNSConfigs();
    ElMessage.success("已触发 DDNS 同步");
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "触发同步失败");
  } finally {
    syncing.value = false;
  }
}

onMounted(load);
</script>

<template>
  <div class="page-shell ddns-page">
    <section class="page-card ddns-hero">
      <div class="page-card__body ddns-hero__body">
        <PageHeader
          title="DDNS"
          description="这轮已经把配置落库、受管配置文件写入、状态读取和 reload / sync 入口补上。这里适合统一维护公网域名到家庭网络的动态解析。"
        >
          <el-space wrap>
            <el-button :loading="loading" @click="load">刷新状态</el-button>
            <el-button type="primary" :loading="syncing" @click="syncNow">立即同步</el-button>
          </el-space>
        </PageHeader>

        <div class="ddns-summary">
          <div class="surface-muted ddns-summary__item">
            <span>记录总数</span>
            <strong>{{ items.length }}</strong>
          </div>
          <div class="surface-muted ddns-summary__item">
            <span>启用中</span>
            <strong>{{ enabledCount }}</strong>
          </div>
          <div class="surface-muted ddns-summary__item">
            <span>最近状态</span>
            <strong>{{ runtimeStatus?.lastStatus || "unknown" }}</strong>
          </div>
        </div>
      </div>
    </section>

    <section class="ddns-grid">
      <article class="page-card">
        <div class="page-card__body">
          <PageHeader
            :title="actionLabel"
            description="当前主路径仍然是阿里云风格的动态解析。更新已有配置时，Secret 不会回显，需要重新输入或保持后端原值。"
          />
          <el-form label-position="top" class="ddns-form">
            <el-form-item label="Provider">
              <el-input v-model="form.provider" placeholder="aliyun" />
            </el-form-item>
            <el-form-item label="AccessKey ID">
              <el-input v-model="form.accessKeyId" />
            </el-form-item>
            <el-form-item label="AccessKey Secret">
              <el-input v-model="form.accessKeySecret" type="password" show-password />
            </el-form-item>
            <el-form-item label="根域名">
              <el-input v-model="form.domain" placeholder="example.com" />
            </el-form-item>
            <el-form-item label="记录名">
              <el-input v-model="form.subdomain" placeholder="home" />
            </el-form-item>
            <el-form-item label="检查间隔（秒）">
              <el-input-number v-model="form.checkIntervalSeconds" :min="60" :step="60" />
            </el-form-item>
            <el-form-item>
              <el-switch v-model="form.enabled" active-text="启用该 DDNS 记录" />
            </el-form-item>
            <el-space wrap>
              <el-button type="primary" :loading="saving" @click="submit">{{ actionLabel }}</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-space>
          </el-form>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <PageHeader
            title="运行状态"
            description="重点看最近状态、当前 IP 和最后同步时间；如果状态异常，优先结合错误信息判断是凭据、网络还是外部 provider 问题。"
          />
          <div class="status-grid">
            <div class="surface-muted status-item">
              <small>最近状态</small>
              <strong>{{ runtimeStatus?.lastStatus || "unknown" }}</strong>
            </div>
            <div class="surface-muted status-item">
              <small>当前 IPv6</small>
              <strong>{{ runtimeStatus?.lastKnownIpv6 || "-" }}</strong>
            </div>
            <div class="surface-muted status-item">
              <small>最后同步</small>
              <strong>{{ runtimeStatus?.lastSyncedAt || "-" }}</strong>
            </div>
            <div class="surface-muted status-item">
              <small>错误</small>
              <strong>{{ runtimeStatus?.lastError || "-" }}</strong>
            </div>
          </div>

          <el-table :data="items" class="ddns-table">
            <el-table-column prop="provider" label="Provider" width="100" />
            <el-table-column label="记录" min-width="180">
              <template #default="{ row }">
                {{ row.subdomain }}.{{ row.domain }}
              </template>
            </el-table-column>
            <el-table-column prop="checkIntervalSeconds" label="间隔" width="100" />
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? "启用" : "停用" }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="lastStatus" label="最近结果" width="120" />
            <el-table-column prop="lastError" label="错误信息" min-width="220" />
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="startEdit(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.ddns-hero {
  background:
    radial-gradient(circle at 88% 10%, rgba(245, 184, 65, 0.18), transparent 18%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(249, 250, 246, 0.97));
}

.ddns-hero__body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ddns-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.ddns-summary__item {
  padding: 16px 18px;
}

.ddns-summary__item span {
  display: block;
  color: var(--text-muted);
  font-size: 13px;
}

.ddns-summary__item strong {
  display: block;
  margin-top: 8px;
  font-size: 26px;
  word-break: break-word;
}

.ddns-grid {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 20px;
}

.ddns-form {
  margin-top: 18px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin-top: 18px;
}

.status-item {
  padding: 16px;
}

.status-item small {
  color: var(--text-muted);
}

.status-item strong {
  display: block;
  margin-top: 8px;
  word-break: break-all;
}

.ddns-table {
  margin-top: 18px;
}

@media (max-width: 1200px) {
  .ddns-grid,
  .ddns-summary,
  .status-grid {
    grid-template-columns: 1fr;
  }
}
</style>
