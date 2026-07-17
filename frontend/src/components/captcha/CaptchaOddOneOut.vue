<script setup>
import { ref } from 'vue'

const emit = defineEmits(['solved', 'failed'])
const GRID_SIZE = 9

const oddIndex = ref(0)
const hue = ref(200)
const lightnessDelta = ref(6)

function generate() {
  oddIndex.value = Math.floor(Math.random() * GRID_SIZE)
  hue.value = Math.floor(Math.random() * 360)
  lightnessDelta.value = 4 + Math.random() * 2
  if (Math.random() > 0.5) lightnessDelta.value *= -1
}
generate()

function colorFor(i) {
  const lightness = i === oddIndex.value ? 52 + lightnessDelta.value : 52
  return `hsl(${hue.value}, 55%, ${lightness}%)`
}

function select(i) {
  if (i === oddIndex.value) {
    emit('solved')
  } else {
    emit('failed')
    generate()
  }
}
</script>

<template>
  <div class="captcha-odd">
    <p class="captcha-label">ЗНАЙДИ ВІДМІННИЙ СЕГМЕНТ</p>
    <div class="odd-grid">
      <button
        v-for="i in GRID_SIZE"
        :key="i"
        class="odd-cell"
        @click="select(i - 1)"
        aria-label="Сегмент"
      >
        <span class="diamond" :style="{ background: colorFor(i - 1) }"></span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.captcha-odd {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}
.captcha-label {
  font-family: var(--jarvis-mono);
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--jarvis-text-dim);
  margin: 0;
}
.odd-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  width: 100%;
  max-width: 260px;
}
.odd-cell {
  aspect-ratio: 1;
  border: 1px solid rgba(126, 20, 255, 0.35);
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  position: relative;
  transition: border-color 0.15s;
}
.diamond {
  position: absolute;
  inset: 26%;
  clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%);
}
.odd-cell:hover,
.odd-cell:focus-visible {
  border-color: var(--jarvis-cyan);
  outline: none;
}
</style>
