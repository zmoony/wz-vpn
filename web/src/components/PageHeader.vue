<script setup lang="ts">
defineProps<{
  title: string;
  description: string;
}>();
</script>

<template>
  <div class="page-header">
    <div class="page-header__content">
      <div class="page-header__eyebrow">山水控制台</div>
      <h2>{{ title }}</h2>
      <p>{{ description }}</p>
    </div>
    <slot />
  </div>
</template>
