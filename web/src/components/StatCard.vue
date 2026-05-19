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
      <small class="stat-card__label">{{ label }}</small>
      <strong class="stat-card__value">{{ value }}</strong>
      <p class="stat-card__helper">{{ helper }}</p>
    </div>
  </article>
</template>

<style scoped>
.stat-card {
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(249, 246, 239, 0.96));
}

.stat-card__label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--earth);
  font-size: 13px;
  letter-spacing: 0.04em;
}

.stat-card__label::before {
  content: "";
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: linear-gradient(180deg, var(--accent), var(--brand));
}

.stat-card__value {
  display: block;
  margin-top: 14px;
  color: var(--text);
  font-size: 30px;
  line-height: 1.15;
}

.stat-card__helper {
  margin: 14px 0 0;
  color: var(--text-secondary);
  line-height: 1.6;
}
</style>
