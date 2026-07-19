<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { useQuestStep } from '../composables/useQuestStep'
import { useHint } from '../composables/useHint'
import { playSound } from '../composables/useSound'
import StepSuccessOverlay from '../components/StepSuccessOverlay.vue'
import transmissionFile from '../assets/audio/step2-transmission.wav'
import step2Voice from '../assets/audio/step2-voice.mp3'
import errorSound from '../assets/audio/access-denied.wav'

const { loading, error, success, submitAnswer, proceedAfterSuccess } = useQuestStep(2)
const { scanning, revealed, hintText, remainingLabel, cooldown, checkHint } = useHint(2)
const answer = ref('')
const rejected = ref(false)
const audioEl = ref(null)
const isPlaying = ref(false)
const canvasEl = ref(null)
let audioCtx = null
let analyser = null
let sourceNode = null
let rafId = null

function drawWaveform() {
  if (!analyser || !canvasEl.value) return
  const canvas = canvasEl.value
  const ctx = canvas.getContext('2d')
  const bufferLength = analyser.fftSize
  const dataArray = new Uint8Array(bufferLength)

  function loop() {
    rafId = requestAnimationFrame(loop)
    analyser.getByteTimeDomainData(dataArray)
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    ctx.lineWidth = 2
    ctx.strokeStyle = '#47bfff'
    ctx.beginPath()
    const sliceWidth = canvas.width / bufferLength
    let x = 0
    for (let i = 0; i < bufferLength; i++) {
      const v = dataArray[i] / 128.0
      const y = (v * canvas.height) / 2
      if (i === 0) ctx.moveTo(x, y)
      else ctx.lineTo(x, y)
      x += sliceWidth
    }
    ctx.stroke()
  }
  loop()
}

function setupAnalyser() {
  if (audioCtx) return
  audioCtx = new (window.AudioContext || window.webkitAudioContext)()
  sourceNode = audioCtx.createMediaElementSource(audioEl.value)
  analyser = audioCtx.createAnalyser()
  analyser.fftSize = 2048
  sourceNode.connect(analyser)
  analyser.connect(audioCtx.destination)
}

function togglePlay() {
  setupAnalyser()
  if (audioCtx.state === 'suspended') audioCtx.resume()

  if (isPlaying.value) {
    audioEl.value.pause()
  } else {
    audioEl.value.play()
    drawWaveform()
  }
}

function onPlay() {
  isPlaying.value = true
}
function onPause() {
  isPlaying.value = false
  if (rafId) cancelAnimationFrame(rafId)
}
function onEnded() {
  isPlaying.value = false
  if (rafId) cancelAnimationFrame(rafId)
}

onMounted(() => {
  const tl = gsap.timeline({ defaults: { ease: 'power2.out' } })
  tl.from('.signal-panel', { opacity: 0, y: 20, duration: 0.6 })
    .eventCallback('onComplete', () => {
      playSound('step2-voice', step2Voice, { volume: 0.9 })
    })
})

onUnmounted(() => {
  if (rafId) cancelAnimationFrame(rafId)
  if (audioCtx) audioCtx.close()
})

async function submit() {
  if (!answer.value.trim() || loading.value) return
  const ok = await submitAnswer(answer.value)
  if (!ok) {
    rejected.value = true
    playSound('access-denied', errorSound, { volume: 0.6 })
    gsap.timeline()
      .fromTo('.signal-panel', { x: -5 }, { x: 0, duration: 0.4, ease: 'elastic.out(1, 0.4)' })
      .fromTo('.rejected-tag', { opacity: 0, scale: 1.4 }, { opacity: 1, scale: 1, duration: 0.25 }, 0)
      .to('.rejected-tag', { opacity: 0, duration: 0.4, delay: 1.2 })
  }
}
</script>

<template>
  <JarvisLayout eyebrow="ПЕРЕХОПЛЕНИЙ СИГНАЛ" :show-scan="false">
    <StepSuccessOverlay v-if="success" @done="proceedAfterSuccess" />
    <div class="signal-panel">
      <p class="signal-id">ФРАГМЕНТ ПЕРЕДАЧІ &middot; ДЖЕРЕЛО НЕВІДОМЕ &middot; ЦІЛІСНІСТЬ: 41%</p>
      <h1 class="signal-title">Розшифрувати перехоплення</h1>
      <p class="signal-text">
        Архів перехопив уривок передачі. Аудіоряд пошкоджено — на слух це
        суцільний шум. Але дані рідко втрачаються повністю, вони просто
        змінюють форму.
      </p>

      <audio
        ref="audioEl"
        :src="transmissionFile"
        @play="onPlay"
        @pause="onPause"
        @ended="onEnded"
        preload="auto"
      ></audio>

      <div class="signal-player">
        <button class="signal-play-btn" @click="togglePlay">
          {{ isPlaying ? 'ПАУЗА' : 'ВІДТВОРИТИ' }}
        </button>
        <canvas ref="canvasEl" class="signal-canvas" width="400" height="80"></canvas>
      </div>

      <a :href="transmissionFile" download="signal-fragment.wav" class="signal-download">
        &darr; ЗАВАНТАЖИТИ ФАЙЛ ДЛЯ АНАЛІЗУ
      </a>

      <div class="answer-row">
        <label class="answer-label" for="step2-answer">РОЗШИФРОВАНЕ СЛОВО:</label>
        <input
          id="step2-answer"
          v-model="answer"
          class="jarvis-input"
          :class="{ 'is-error': error }"
          autocomplete="off"
          @keyup.enter="submit"
        />
      </div>

      <button class="jarvis-btn" :disabled="loading" @click="submit">
        {{ loading ? 'ПЕРЕВІРКА...' : 'ПІДТВЕРДИТИ' }}
      </button>

      <div v-if="rejected" class="rejected-tag">СИГНАЛ НЕ РОЗПІЗНАНО</div>

      <div class="hint-block">
        <template v-if="!revealed">
          <button class="hint-btn" :disabled="scanning || cooldown" @click="checkHint">
            {{ scanning ? 'АНАЛІЗ...' : 'ЗАПУСТИТИ ПОГЛИБЛЕНИЙ АНАЛІЗ' }}
          </button>
          <p v-if="remainingLabel !== null" class="hint-wait">
            Аналіз ще не завершено. Підказка стане доступною за {{ remainingLabel }}.
          </p>
        </template>
        <template v-else>
          <p class="hint-label">РЕЗУЛЬТАТ АНАЛІЗУ:</p>
          <p class="hint-text">{{ hintText }}</p>
        </template>
      </div>
    </div>
  </JarvisLayout>
</template>

<style scoped>
.signal-panel {
  width: 100%;
  max-width: 480px;
}
.signal-id {
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  color: var(--jarvis-text-dim);
  margin: 0 0 14px;
}
.signal-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(19px, 4.8vw, 24px);
  margin: 0 0 14px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.signal-text {
  font-size: clamp(13px, 3vw, 15px);
  line-height: 1.6;
  color: var(--jarvis-text-dim);
  margin: 0 0 20px;
}
.signal-player {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px;
  background: rgba(13, 18, 32, 0.7);
  border: 1px solid rgba(126, 20, 255, 0.35);
  border-radius: 6px;
  margin-bottom: 10px;
}
.signal-play-btn {
  flex-shrink: 0;
  padding: 10px 14px;
  background: none;
  border: 1px solid var(--jarvis-cyan);
  border-radius: 4px;
  color: var(--jarvis-cyan);
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  cursor: pointer;
}
.signal-play-btn:hover {
  background: rgba(71, 191, 255, 0.1);
}
.signal-canvas {
  flex: 1;
  min-width: 0;
  width: 100%;
  aspect-ratio: 400 / 80;
  max-height: 80px;
  background: rgba(5, 7, 10, 0.6);
  border-radius: 4px;
}
.signal-download {
  display: inline-block;
  margin-bottom: 8px;
  color: var(--jarvis-cyan);
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  letter-spacing: 1px;
  text-decoration: none;
  border-bottom: 1px dashed var(--jarvis-cyan);
}
.signal-download:hover,
.signal-download:focus-visible {
  border-bottom-style: solid;
}
.signal-note {
  font-size: clamp(10px, 2.4vw, 12px);
  color: var(--jarvis-text-dim);
  font-style: italic;
  margin: 0 0 22px;
}
.answer-row {
  margin-bottom: 16px;
}
.answer-label {
  display: block;
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1.5px;
  color: var(--jarvis-text-dim);
  margin-bottom: 8px;
}
.rejected-tag {
  margin-top: 12px;
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  color: var(--jarvis-danger);
  letter-spacing: 1px;
}
.hint-block {
  margin-top: 24px;
}
.hint-btn {
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
.hint-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.hint-wait,
.hint-label {
  font-size: clamp(11px, 2.6vw, 13px);
  color: var(--jarvis-text-dim);
  margin: 10px 0 0;
}
.hint-text {
  font-family: var(--jarvis-mono);
  font-size: clamp(12px, 2.8vw, 14px);
  color: var(--jarvis-cyan);
  margin: 6px 0 0;
}
</style>