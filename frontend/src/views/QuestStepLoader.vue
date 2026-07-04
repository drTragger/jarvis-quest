<script setup>
import { computed, defineAsyncComponent } from 'vue'
import { useRouter } from 'vue-router'
import { stepComponents } from '../data/steps'

const props = defineProps({
  stepId: { type: [String, Number], required: true }
})

const router = useRouter()
const stepNum = computed(() => Number(props.stepId))
const loader = stepComponents[stepNum.value]

if (!loader) {
  router.replace({ name: 'not-found' })
}

const StepComponent = loader ? defineAsyncComponent(loader) : null

function goToHub() {
  router.push({ name: 'hub' })
}
</script>

<template>
  <button class="quest-back" @click="goToHub">&lt; МЕНЮ</button>
  <component :is="StepComponent" v-if="StepComponent" />
</template>

<style scoped>
.quest-back {
  position: fixed;
  top: 14px;
  left: 14px;
  z-index: 50;
  background: rgba(5, 7, 10, 0.75);
  border: 1px solid rgba(126, 20, 255, 0.4);
  border-radius: 4px;
  padding: 6px 10px;
  color: var(--jarvis-cyan);
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  cursor: pointer;
  backdrop-filter: blur(2px);
}
.quest-back:hover,
.quest-back:focus-visible {
  border-color: var(--jarvis-cyan);
  background: rgba(5, 7, 10, 0.9);
}
</style>