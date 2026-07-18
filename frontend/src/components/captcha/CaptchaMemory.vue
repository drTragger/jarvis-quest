<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const emit = defineEmits(['solved', 'failed'])
const COLORS = ['#47bfff', '#7e14ff', '#ff5470', '#ffb648']
const SEQUENCE_LENGTH = 20

const sequence = ref([])
const userProgress = ref(0)
const activeIndex = ref(-1)
const clickIndex = ref(-1)
const playing = ref(false)
let timers = []
let clickTimer = null

function clearTimers() {
  timers.forEach(clearTimeout)
  timers = []
}

function generate() {
  sequence.value = Array.from({ length: SEQUENCE_LENGTH }, () => Math.floor(Math.random() * COLORS.length))
  userProgress.value = 0
}
generate()

function play() {
  if (playing.value) return
  clearTimers()
  playing.value = true
  sequence.value.forEach((tileIndex, step) => {
    timers.push(setTimeout(() => {
      activeIndex.value = tileIndex
      timers.push(setTimeout(() => { activeIndex.value = -1 }, 320))
    }, step * 520))
  })
  timers.push(setTimeout(() => {
    playing.value = false
  }, sequence.value.length * 520))
}

function pick(tileIndex) {
  if (playing.value) return

  clickIndex.value = tileIndex
  clearTimeout(clickTimer)
  clickTimer = setTimeout(() => { clickIndex.value = -1 }, 180)

  const expected = sequence.value[userProgress.value]
  if (tileIndex !== expected) {
    emit('failed')
    generate()
    return
  }
  userProgress.value += 1
  if (userProgress.value === sequence.value.length) {
    emit('solved')
  }
}

onMounted(() => {
  timers.push(setTimeout(play, 500))
})
onBeforeUnmount(() => {
  clearTimers()
  clearTimeout(clickTimer)
})
</script>

<template>
  <div class="captcha-memory">
    <p class="captcha-label">ЗАПАМ'ЯТАЙ ТА ПОВТОРИ ПОСЛІДОВНІСТЬ ({{ sequence.length }})</p>
    <div class="memory-grid">
      <button
        v-for="(color, i) in COLORS"
        :key="i"
        class="memory-tile"
        :style="{ background: color, opacity: (activeIndex === i || clickIndex === i) ? 1 : 0.3 }"
        :disabled="playing"
        @click="pick(i)"
        aria-label="Сегмент пам'яті"
      ></button>
    </div>
    <button class="replay-btn" :disabled="playing" @click="play">ВІДТВОРИТИ ЩЕ РАЗ</button>
  </div>
</template>

<style scoped>
.captcha-memory {
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
  text-align: center;
}
.memory-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  width: 100%;
  max-width: 190px;
}
.memory-tile {
  aspect-ratio: 1;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  cursor: pointer;
  transition: opacity 0.12s;
}
.memory-tile:disabled {
  cursor: not-allowed;
}
.replay-btn {
  background: none;
  border: 1px solid rgba(232, 230, 240, 0.25);
  border-radius: 4px;
  padding: 8px 12px;
  color: var(--jarvis-text-dim);
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  cursor: pointer;
}
.replay-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
