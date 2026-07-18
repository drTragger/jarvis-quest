<script setup>
import { ref, onBeforeUnmount } from 'vue'

const emit = defineEmits(['solved', 'failed'])
const TARGET_WIDTH = 1.4

const trackEl = ref(null)
const handlePos = ref(0)
const targetStart = ref(0)
const dragging = ref(false)

function generate() {
  targetStart.value = 12 + Math.random() * (84 - TARGET_WIDTH)
  handlePos.value = 0
}
generate()

function clamp(v) {
  return Math.max(0, Math.min(100, v))
}

function updateFromClientX(clientX) {
  const rect = trackEl.value.getBoundingClientRect()
  const pct = ((clientX - rect.left) / rect.width) * 100
  handlePos.value = clamp(pct)
}

function onPointerDown(e) {
  dragging.value = true
  updateFromClientX(e.clientX)
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp)
}
function onPointerMove(e) {
  if (!dragging.value) return
  updateFromClientX(e.clientX)
}
function onPointerUp() {
  dragging.value = false
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
  const inZone = handlePos.value >= targetStart.value && handlePos.value <= targetStart.value + TARGET_WIDTH
  if (inZone) {
    emit('solved')
  } else {
    emit('failed')
    generate()
  }
}

onBeforeUnmount(() => {
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
})
</script>

<template>
  <div class="captcha-slider">
    <p class="captcha-label">СУМІСТИ МАРКЕР ІЗ ВУЗЬКОЮ ЗОНОЮ КАЛІБРУВАННЯ</p>
    <div ref="trackEl" class="slider-track">
      <div class="slider-zone" :style="{ left: targetStart + '%', width: TARGET_WIDTH + '%' }"></div>
      <div
        class="slider-handle"
        :class="{ 'is-dragging': dragging }"
        :style="{ left: handlePos + '%' }"
        @pointerdown="onPointerDown"
      ></div>
    </div>
  </div>
</template>

<style scoped>
.captcha-slider {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}
.captcha-label {
  font-family: var(--jarvis-mono);
  font-size: 11px;
  letter-spacing: 1.5px;
  color: var(--jarvis-text-dim);
  margin: 0;
  text-align: center;
}
.slider-track {
  position: relative;
  width: 100%;
  max-width: 300px;
  height: 8px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 4px;
  margin: 26px 0 20px;
}
.slider-zone {
  position: absolute;
  top: -6px;
  bottom: -6px;
  min-width: 2px;
  background: var(--jarvis-cyan);
  box-shadow: 0 0 8px rgba(71, 191, 255, 0.7);
}
.slider-handle {
  position: absolute;
  top: 50%;
  width: 22px;
  height: 22px;
  background: var(--jarvis-purple);
  border: 2px solid var(--jarvis-text);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  cursor: grab;
  box-shadow: 0 0 10px rgba(126, 20, 255, 0.6);
  touch-action: none;
}
.slider-handle.is-dragging {
  cursor: grabbing;
}
</style>
