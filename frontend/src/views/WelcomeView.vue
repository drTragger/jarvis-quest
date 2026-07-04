<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { resolveDestination } from '../router'
import { unlockAudio, playSound } from '../composables/useSound'
import typeSound from '../assets/audio/type-tick.wav'
import greetingSound from '../assets/audio/greeting-vlad.mp3'

const router = useRouter()
const activated = ref(false)
const bootLines = ref([])
const progress = ref(0)
const showButton = ref(false)

const log = [
  '> ЗАПУСК ПРОТОКОЛУ J.A.R.V.I.S. ... OK',
  '> ПЕРЕВІРКА СИСТЕМ ... OK',
  '> ІДЕНТИФІКАЦІЯ КОРИСТУВАЧА: ВЛАД ... OK',
  '> З ДНЕМ НАРОДЖЕННЯ. ГОТОВИЙ ДО ІНІЦІАЛІЗАЦІЇ.'
]

function activateSystem() {
  activated.value = true
  unlockAudio(typeSound).then(() => {
    runIntroSequence()
  })
}

function runIntroSequence() {
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

  if (reduceMotion) {
    bootLines.value = [...log]
    progress.value = 100
    showButton.value = true
    return
  }

  const tl = gsap.timeline({ defaults: { ease: 'power2.out' } })
  tl.from('.core', { opacity: 0, scale: 0.6, duration: 0.8 })
    .from('.title', { opacity: 0, y: 16, duration: 0.6 }, '-=0.3')
    .from('.subtitle', { opacity: 0, y: 10, duration: 0.5 }, '-=0.3')
    .eventCallback('onComplete', runBootSequence)
}

async function start() {
  const dest = await resolveDestination()
  router.push(dest)
}

function typeLine(text) {
  return new Promise((resolve) => {
    bootLines.value.push('')
    const idx = bootLines.value.length - 1
    let i = 0
    function step() {
      i++
      bootLines.value[idx] = text.slice(0, i)
      if (i % 2 === 0) {
        playSound('type-tick', typeSound, { volume: 0.15 })
      }
      if (i < text.length) {
        setTimeout(step, 14 + Math.random() * 18)
      } else {
        resolve()
      }
    }
    step()
  })
}

async function runBootLines() {
  for (const line of log) {
    await typeLine(line)
    await new Promise((r) => setTimeout(r, 220))
  }
}

function runFakeProgress() {
  return new Promise((resolve) => {
    function step() {
      const remaining = 100 - progress.value
      const increment = Math.max(0.6, Math.random() * remaining * 0.18)
      progress.value = Math.min(99.4, progress.value + increment)
      if (progress.value < 99.3) {
        setTimeout(step, 90 + Math.random() * 170)
      } else {
        setTimeout(() => {
          progress.value = 100
          resolve()
        }, 350)
      }
    }
    step()
  })
}

async function runBootSequence() {
  await Promise.all([runBootLines(), runFakeProgress()])
  showButton.value = true
  gsap.from('.jarvis-btn', { opacity: 0, y: 10, duration: 0.5 })
  playSound('greeting-vlad', greetingSound, { volume: 0.9 })
}
</script>

<template>
  <JarvisLayout eyebrow="ОСОБИСТИЙ ІНТЕРФЕЙС">
    <template v-if="!activated">
      <div class="core" aria-hidden="true">
        <div class="ring ring-outer"></div>
        <div class="ring ring-mid"></div>
        <div class="ring ring-inner"></div>
        <div class="core-dot"></div>
      </div>
      <h1 class="title">J&nbsp;·&nbsp;A&nbsp;·&nbsp;R&nbsp;·&nbsp;V&nbsp;·&nbsp;I&nbsp;·&nbsp;S</h1>
      <p class="subtitle">Система в режимі очікування.</p>
      <button class="jarvis-btn" @click="activateSystem">
        <span>АКТИВУВАТИ СИСТЕМУ</span>
      </button>
    </template>

    <template v-else>
      <div class="core" aria-hidden="true">
        <div class="ring ring-outer"></div>
        <div class="ring ring-mid"></div>
        <div class="ring ring-inner"></div>
        <div class="core-dot"></div>
      </div>

      <h1 class="title">J&nbsp;·&nbsp;A&nbsp;·&nbsp;R&nbsp;·&nbsp;V&nbsp;·&nbsp;I&nbsp;·&nbsp;S</h1>
      <p class="subtitle">Система активна. Очікую команди.</p>

      <div class="boot-log">
        <p v-for="(line, i) in bootLines" :key="i" class="boot-line">
          {{ line }}<span v-if="i === bootLines.length - 1 && !showButton" class="cursor">▋</span>
        </p>
      </div>

      <div class="progress-wrap">
        <div class="progress-label">
          <span>ЗАВАНТАЖЕННЯ СИСТЕМИ</span>
          <span>{{ Math.floor(progress) }}%</span>
        </div>
        <div class="progress-track">
          <div class="progress-fill" :style="{ width: progress + '%' }"></div>
        </div>
      </div>

      <button v-if="showButton" class="jarvis-btn" @click="start">
        <span>ІНІЦІЮВАТИ</span>
      </button>
    </template>
  </JarvisLayout>
</template>

<style scoped>
.core {
  position: relative;
  width: clamp(64px, 20vw, 96px);
  height: clamp(64px, 20vw, 96px);
  margin-bottom: clamp(18px, 4vw, 28px);
}

.ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 1.5px solid transparent;
}
.ring-outer {
  border-top-color: #7e14ff;
  border-bottom-color: #7e14ff;
  animation: spin 6s linear infinite;
}
.ring-mid {
  inset: 15%;
  border-left-color: #47bfff;
  border-right-color: #47bfff;
  animation: spin 4s linear infinite reverse;
}
.ring-inner {
  inset: 32%;
  border-top-color: #e8e6f0;
  animation: spin 2.5s linear infinite;
}
.core-dot {
  position: absolute;
  inset: 44%;
  border-radius: 50%;
  background: #47bfff;
  box-shadow: 0 0 16px 4px rgba(71, 191, 255, 0.7);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }
@keyframes pulse {
  0%, 100% { opacity: 0.6; transform: scale(0.9); }
  50% { opacity: 1; transform: scale(1.1); }
}

.title {
  font-family: 'Orbitron', sans-serif;
  font-weight: 700;
  font-size: clamp(22px, 7vw, 44px);
  letter-spacing: clamp(2px, 1vw, 4px);
  margin: 0 0 10px;
  text-shadow: 0 0 18px rgba(126, 20, 255, 0.6);
}

.subtitle {
  font-size: clamp(12px, 3vw, 14px);
  color: #9ca3af;
  margin: 0 0 24px;
}

.boot-log {
  width: 100%;
  min-height: 96px;
  text-align: left;
  font-size: clamp(11px, 3vw, 13px);
  line-height: 1.9;
  color: #7fd4ff;
  word-break: break-word;
}

.boot-line {
  margin: 0;
}

.cursor {
  display: inline-block;
  margin-left: 2px;
  animation: blink 1s step-start infinite;
  color: #47bfff;
}
@keyframes blink { 50% { opacity: 0; } }

.progress-wrap {
  width: 100%;
  margin-top: 18px;
}

.progress-label {
  display: flex;
  justify-content: space-between;
  font-size: clamp(9px, 2.4vw, 11px);
  letter-spacing: 1.5px;
  color: #9ca3af;
  margin-bottom: 6px;
}

.progress-track {
  width: 100%;
  height: 6px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(126, 20, 255, 0.3);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #7e14ff, #47bfff);
  box-shadow: 0 0 8px rgba(71, 191, 255, 0.6);
  transition: width 0.15s linear;
}

.jarvis-btn {
  margin-top: 24px;
}

@media (prefers-reduced-motion: reduce) {
  .ring-outer, .ring-mid, .ring-inner, .core-dot, .cursor {
    animation: none;
  }
}
</style>