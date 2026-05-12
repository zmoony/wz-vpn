<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

import {
  createProxyHost,
  deleteProxyHost,
  fetchProxyHosts,
  updateProxyHost,
  type ProxyHost,
  type ProxyPayload,
} from "../api/proxy";
import PageHeader from "../components/PageHeader.vue";

const loading = ref(false);
const saving = ref(false);
const editingId = ref<number | null>(null);
const items = ref<ProxyHost[]>([]);

const form = reactive<ProxyPayload>({
  name: "",
  serverName: "",
  upstreamUrl: "http://127.0.0.1:3000",
  certificateCertPath: "/opt/pi-gateway/data/certs/fullchain.pem",
  certificateKeyPath: "/opt/pi-gateway/data/certs/privkey.pem",
  enabled: true,
  description: "",
});

const actionLabel = computed(() => (editingId.value ? "更新反代" : "新增反代"));

function resetForm() {
  editingId.value = null;
  form.name = "";
  form.serverName = "";
  form.upstreamUrl = "http://127.0.0.1:3000";
  form.certificateCertPath = "/opt/pi-gateway/data/certs/fullchain.pem";
  form.certificateKeyPath = "/opt/pi-gateway/data/certs/privkey.pem";
  form.enabled = true;
  form.description = "";
}

async function load() {
  loading.value = true;
  try {
    const response = await fetchProxyHosts();
    items.value = response.items;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载反代列表失败");
  } finally {
    loading.value = false;
  }
}

function startEdit(item: ProxyHost) {
  editingId.value = item.id;
  form.name = item.name;
  form.serverName = item.serverName;
  form.upstreamUrl = item.upstreamUrl;
  form.certificateCertPath = item.certificateCertPath;
  form.certificateKeyPath = item.certificateKeyPath;
  form.enabled = item.enabled;
  form.description = item.description;
}

async function submit() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateProxyHost(editingId.value, { ...form });
      ElMessage.success("反向代理已更新");
    } else {
      await createProxyHost({ ...form });
      ElMessage.success("反向代理已创建");
    }
    resetForm();
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存反代配置失败");
  } finally {
    saving.value = false;
  }
}

async function remove(item: ProxyHost) {
  await ElMessageBox.confirm(`确定删除反代 "${item.name}" 吗？`, "删除确认", { type: "warning" });
  try {
    await deleteProxyHost(item.id);
    ElMessage.success("反代已删除");
    if (editingId.value === item.id) {
      resetForm();
    }
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "删除反代失败");
  }
}

onMounted(load);
</script>

<template>
  <div class="page-shell proxy-page">
    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="反向代理"
          description="这一轮已经把配置落库、nginx vhost 生成、nginx -t 校验和 reload 链路接起来了。"
        >
          <el-button :loading="loading" @click="load">刷新列表</el-button>
        </PageHeader>
      </div>
    </section>

    <section class="proxy-grid">
      <article class="page-card">
        <div class="page-card__body">
          <h3>{{ actionLabel }}</h3>
          <el-form label-position="top">
            <el-form-item label="名称">
              <el-input v-model="form.name" placeholder="如 blog" />
            </el-form-item>
            <el-form-item label="子域名">
              <el-input v-model="form.serverName" placeholder="如 blog.example.com" />
            </el-form-item>
            <el-form-item label="后端地址">
              <el-input v-model="form.upstreamUrl" placeholder="http://127.0.0.1:3000" />
            </el-form-item>
            <el-form-item label="证书路径">
              <el-input v-model="form.certificateCertPath" placeholder="/path/to/fullchain.pem" />
            </el-form-item>
            <el-form-item label="私钥路径">
              <el-input v-model="form.certificateKeyPath" placeholder="/path/to/privkey.pem" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="form.description" type="textarea" :rows="3" placeholder="可选备注" />
            </el-form-item>
            <el-form-item>
              <el-switch v-model="form.enabled" active-text="启用后立即写入并 reload" />
            </el-form-item>
            <el-space>
              <el-button type="primary" :loading="saving" @click="submit">{{ actionLabel }}</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-space>
          </el-form>
        </div>
      </article>

      <article class="page-card">
        <div class="page-card__body">
          <el-table :data="items">
            <el-table-column prop="name" label="名称" width="140" />
            <el-table-column prop="serverName" label="域名" min-width="180" />
            <el-table-column prop="upstreamUrl" label="后端" min-width="180" />
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? "启用" : "停用" }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="最近应用" min-width="160">
              <template #default="{ row }">
                <span v-if="row.lastApplyStatus === 'applied'">已应用</span>
                <span v-else-if="row.lastApplyStatus === 'error'">失败</span>
                <span v-else>未应用</span>
              </template>
            </el-table-column>
            <el-table-column label="错误信息" min-width="220">
              <template #default="{ row }">
                {{ row.lastApplyError || "-" }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-space>
                  <el-button link type="primary" @click="startEdit(row)">编辑</el-button>
                  <el-button link type="danger" @click="remove(row)">删除</el-button>
                </el-space>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.proxy-grid {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 20px;
}

@media (max-width: 1200px) {
  .proxy-grid {
    grid-template-columns: 1fr;
  }
}
</style>
