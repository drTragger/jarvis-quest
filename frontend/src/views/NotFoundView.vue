<script setup>
import { useRouter } from 'vue-router'
import Button from 'primevue/button'

const router = useRouter()
</script>

<template>
  <section class="not-found">
    <h1>404</h1>
    <Button label="На головну" @click="router.push({ name: 'welcome' })" />
  </section>
</template>