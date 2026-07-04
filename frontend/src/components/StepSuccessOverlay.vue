<script setup>
import { onMounted } from 'vue'
import gsap from 'gsap'

const emit = defineEmits(['done'])

onMounted(() => {
  gsap.timeline({ defaults: { ease: 'power2.out' } })
    .from('.success-core', { opacity: 0, scale: 0.4, duration: 0.5, ease: 'back.out(2)' })
    .from('.success-text', { opacity: 0, y: 10, duration: 0.4 }, '-=0.2')
    .to('.success-panel', { opacity: 0, duration: 0.4, delay: 1.1 })
    .eventCallback('onComplete', () => emit('done'))
})
</script>

<template>
  <div class="success-overlay">
    <div class="success-panel">
      <div class="success-core" aria-hidden="true">
        <div class="success-ring success-ring-outer"></div>
        <div class="success-ring success-ring-inner"></div>
        <div class="success-dot"></div>
      </div>
      <p class="success-text">СЕГМЕНТ РОЗШИФРОВАНО</p>
    </div>
  </div>
</template>

<style scoped>
.success-overlay {
  position: fixed;
  inset: 0;
  z-index: 90;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(5, 7, 10, 0.92);
}
.success-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 18px;
}
.success-core {
  position: relative;
  width: clamp(72px, 22vw, 104px);
  height: clamp(72px, 22vw, 104px);
}
.success-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid transparent;
}
.success-ring-outer {
  border-top-color: var(--jarvis-cyan);
  border-bottom-color: var(--jarvis-cyan);
  animation: success-spin 1.4s linear infinite;
}
.success-ring-inner {
  inset: 20%;
  border-left-color: var(--jarvis-purple);
  border-right-color: var(--jarvis-purple);
  animation: success-spin 0.9s linear infinite reverse;
}
.success-dot {
  position: absolute;
  inset: 40%;
  border-radius: 50%;
  background: var(--jarvis-cyan);
  box-shadow: 0 0 20px 6px rgba(71, 191, 255, 0.8);
}
.success-text {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(16px, 4vw, 20px);
  letter-spacing: 3px;
  color: var(--jarvis-text);
  text-shadow: 0 0 18px rgba(71, 191, 255, 0.7);
  margin: 0;
}
@keyframes success-spin {
  to { transform: rotate(360deg); }
}
@media (prefers-reduced-motion: reduce) {
  .success-ring-outer,
  .success-ring-inner {
    animation: none;
  }
}
</style>