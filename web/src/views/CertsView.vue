<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";

import PageHeader from "../components/PageHeader.vue";
import {
  createCertificate,
  fetchCertificates,
  renewCertificate,
  updateCertificate,
  type CertificateConfig,
  type CertificatePayload,
} from "../api/certs";

const loading = ref(false);
const saving = ref(false);
const renewingId = ref<number | null>(null);
const editingId = ref<number | null>(null);
const items = ref<CertificateConfig[]>([]);

const form = reactive<CertificatePayload>({
  rootDomain: "",
  provider: "aliyun",
  accessKeyId: "",
  accessKeySecret: "",
  installDir: "",
  enabledAutoRenew: true,
});

const actionLabel = computed(() => (editingId.value ? "更新证书配置" : "申请证书"));

function resetForm() {
  editingId.value = null;
  form.rootDomain = "";
  form.provider = "aliyun";
  form.accessKeyId = "";
  form.accessKeySecret = "";
  form.installDir = "";
  form.enabledAutoRenew = true;
}

async function load() {
  loading.value = true;
  try {
    const response = await fetchCertificates();
    items.value = response.items;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "加载证书列表失败");
  } finally {
    loading.value = false;
  }
}

function startEdit(item: CertificateConfig) {
  editingId.value = item.id;
  form.rootDomain = item.rootDomain;
  form.provider = item.provider;
  form.accessKeyId = item.accessKeyId;
  form.accessKeySecret = "";
  form.installDir = item.installDir;
  form.enabledAutoRenew = item.enabledAutoRenew;
}

async function submit() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateCertificate(editingId.value, { ...form });
      ElMessage.success("证书配置已更新");
    } else {
      await createCertificate({ ...form });
      ElMessage.success("证书申请已触发");
    }
    resetForm();
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存证书配置失败");
  } finally {
    saving.value = false;
  }
}

async function renew(item: CertificateConfig) {
  renewingId.value = item.id;
  try {
    await renewCertificate(item.id);
    ElMessage.success(`已触发 ${item.rootDomain} 的证书续期`);
    await load();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "触发续期失败");
  } finally {
    renewingId.value = null;
  }
}

onMounted(load);
</script>

<template>
  <div class="page-shell certs-page">
    <section class="page-card">
      <div class="page-card__body">
        <PageHeader
          title="证书管理"
          description="这轮先支持单个根域名的泛域名证书，也就是 *.example.com + example.com 的主路径。"
        >
          <el-button :loading="loading" @click="load">刷新列表</el-button>
        </PageHeader>
      </div>
    </section>

    <section class="certs-grid">
      <article class="page-card">
        <div class="page-card__body">
          <h3>{{ actionLabel }}</h3>
          <el-form label-position="top">
            <el-form-item label="根域名">
              <el-input v-model="form.rootDomain" placeholder="example.com" />
            </el-form-item>
            <el-form-item label="Provider">
              <el-input v-model="form.provider" placeholder="aliyun" />
            </el-form-item>
            <el-form-item label="AccessKey ID">
              <el-input v-model="form.accessKeyId" />
            </el-form-item>
            <el-form-item label="AccessKey Secret">
              <el-input v-model="form.accessKeySecret" type="password" show-password />
            </el-form-item>
            <el-form-item label="安装目录">
              <el-input v-model="form.installDir" placeholder="/opt/pi-gateway/data/certs/example.com" />
            </el-form-item>
            <el-form-item>
              <el-switch v-model="form.enabledAutoRenew" active-text="启用自动续期标记" />
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
            <el-table-column prop="rootDomain" label="根域名" min-width="160" />
            <el-table-column prop="provider" label="Provider" width="100" />
            <el-table-column label="剩余天数" width="120">
              <template #default="{ row }">
                {{ row.notAfter ? row.daysRemaining : "-" }}
              </template>
            </el-table-column>
            <el-table-column label="到期时间" min-width="180">
              <template #default="{ row }">
                {{ row.notAfter || "-" }}
              </template>
            </el-table-column>
            <el-table-column label="最近结果" width="120">
              <template #default="{ row }">
                {{ row.lastIssueStatus || "-" }}
              </template>
            </el-table-column>
            <el-table-column label="证书路径" min-width="220">
              <template #default="{ row }">
                {{ row.fullchainPath || row.installDir }}
              </template>
            </el-table-column>
            <el-table-column label="错误信息" min-width="220">
              <template #default="{ row }">
                {{ row.lastIssueError || "-" }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-space>
                  <el-button link type="primary" @click="startEdit(row)">编辑</el-button>
                  <el-button
                    link
                    type="success"
                    :loading="renewingId === row.id"
                    @click="renew(row)"
                  >
                    续期
                  </el-button>
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
.certs-grid {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 20px;
}

@media (max-width: 1200px) {
  .certs-grid {
    grid-template-columns: 1fr;
  }
}
</style>
