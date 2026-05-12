<script setup lang="ts">
defineProps<{
  label: string;
  value: string;
  helper?: string;
}>();
</script>

<template>
  <article class="stat-card page-card">
    <div class="page-card__body">
      <small>{{ label }}</small>
      <strong>{{ value }}</strong>
      <p>{{ helper }}</p>
    </div>
  </article>
</template>

<style scoped>
.stat-card small {
  color: var(--muted);
}

.stat-card strong {
  display: block;
  margin-top: 10px;
  font-size: 28px;
}

.stat-card p {
  margin: 12px 0 0;
  color: var(--muted);
}
</style>
