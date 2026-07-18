<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import gsap from 'gsap'
import { useQuestStep } from '../composables/useQuestStep'
import { useHint } from '../composables/useHint'
import { useMorseInput, MORSE_MAP } from '../composables/useMorseInput'
import { playSound, stopSound } from '../composables/useSound'
import StepSuccessOverlay from '../components/StepSuccessOverlay.vue'
import step6Voice from '../assets/audio/step6-voice.mp3'
import errorSound from '../assets/audio/access-denied.wav'
import morseBeep from '../assets/audio/morse-beep.wav'

const { loading, error, success, submitAnswer, proceedAfterSuccess } = useQuestStep(6)
const { scanning, revealed, hintText, remainingLabel, cooldown, checkHint } = useHint(6)
const { word, currentSymbol, pressing, keyDown, keyUp, clear } = useMorseInput()
const rejected = ref(false)
const beaconOn = ref(false)

const BEACON_WORD = 'ВЛАДИСЛАВ'
const DOT_MS = 200
const DASH_MS = 600
const SYMBOL_GAP_MS = 200
const LETTER_GAP_MS = 600
const LOOP_GAP_MS = 1600

let beaconTimer = null
let beaconRunning = true

async function wait(ms) {
  return new Promise((resolve) => { beaconTimer = setTimeout(resolve, ms) })
}

async function runBeacon() {
  while (beaconRunning) {
    for (const letter of BEACON_WORD) {
      const code = MORSE_MAP[letter]
      for (const symbol of code) {
        beaconOn.value = true
        await wait(symbol === '.' ? DOT_MS : DASH_MS)
        if (!beaconRunning) return
        beaconOn.value = false
        await wait(SYMBOL_GAP_MS)
        if (!beaconRunning) return
      }
      await wait(LETTER_GAP_MS)
      if (!beaconRunning) return
    }
    await wait(LOOP_GAP_MS)
  }
}

onMounted(() => {
  const tl = gsap.timeline({ defaults: { ease: 'power2.out' } })
  tl.from('.telegraph-panel', { opacity: 0, y: 24, duration: 0.6 })
    .eventCallback('onComplete', () => {
      playSound('step6-voice', step6Voice, { volume: 0.9 })
    })
  runBeacon()
})

onUnmounted(() => {
  beaconRunning = false
  clearTimeout(beaconTimer)
})

async function submit() {
  if (!word.value.trim() || loading.value) return
  const ok = await submitAnswer(word.value.toLowerCase())
  if (!ok) {
    rejected.value = true
    playSound('access-denied', errorSound, { volume: 0.6 })
    clear()
    gsap.timeline()
      .fromTo('.telegraph-panel', { x: -5 }, { x: 0, duration: 0.4, ease: 'elastic.out(1, 0.4)' })
      .fromTo('.rejected-plate', { opacity: 0, scale: 1.4 }, { opacity: 1, scale: 1, duration: 0.25 }, 0)
      .to('.rejected-plate', { opacity: 0, duration: 0.4, delay: 1.2 })
  }
}

function pressKey() {
  keyDown()
  playSound('morse-beep', morseBeep, { volume: 0.7, loop: true })
}
function releaseKey() {
  keyUp()
  stopSound('morse-beep')
}

function onKeyDown(e) {
  if (e.code !== 'Space' || e.repeat) return
  e.preventDefault()
  pressKey()
}
function onKeyUp(e) {
  if (e.code !== 'Space') return
  e.preventDefault()
  releaseKey()
}
</script>

<template>
  <section class="telegraph-page">
    <StepSuccessOverlay v-if="success" @done="proceedAfterSuccess" />
    <div class="telegraph-panel" tabindex="0" @keydown="onKeyDown" @keyup="onKeyUp">
      <div class="rivet rivet-tl" aria-hidden="true"></div>
      <div class="rivet rivet-tr" aria-hidden="true"></div>
      <div class="rivet rivet-bl" aria-hidden="true"></div>
      <div class="rivet rivet-br" aria-hidden="true"></div>

      <p class="plate-id">СТАРИЙ ПРОТОКОЛ &middot; ТІЛЬКИ РУЧНА ПЕРЕДАЧА &middot; АВТОДЕКОДЕР ВІДСУТНІЙ</p>
      <h1 class="plate-title">Отримай сигнал і дай відповідь</h1>
      <p class="plate-text">
        Нам вдалося дізнатися, що відповідь — це позивний. Маяк моргає
        повільно й монотонно, знову і знову одне й те саме слово. Порахуй
        короткі й довгі спалахи сам і вистукай відповідь на ключі.
      </p>

      <div class="lamp-row">
        <div class="lamp-housing">
          <div class="lamp-bulb" :class="{ 'is-on': beaconOn }" aria-hidden="true"></div>
        </div>
        <p class="lamp-label">СИГНАЛ МАЯКА</p>
      </div>

      <div class="ticker-tape">
        <p class="ticker-word">{{ word || '·' }}</p>
        <p class="ticker-buffer">{{ currentSymbol }}</p>
      </div>

      <div class="key-assembly">
        <div class="key-base" aria-hidden="true">
          <div class="key-screw key-screw-l"></div>
          <div class="key-screw key-screw-r"></div>
        </div>
        <button
          class="key-lever"
          :class="{ 'is-pressed': pressing }"
          @mousedown.prevent="pressKey"
          @mouseup.prevent="releaseKey"
          @mouseleave="pressing && releaseKey()"
          @touchstart.prevent="pressKey"
          @touchend.prevent="releaseKey"
        >
          <span class="key-arm"></span>
          <span class="key-knob"></span>
        </button>
      </div>
      <p class="key-hint">Тримай важіль (або пробіл): коротко &mdash; крапка, довго &mdash; тире.</p>

      <div class="brass-actions">
        <button class="brass-btn ghost" type="button" @click="clear">ОЧИСТИТИ</button>
        <button class="brass-btn primary" :disabled="loading" @click="submit">
          {{ loading ? 'ПЕРЕВІРКА...' : 'ПЕРЕДАТИ' }}
        </button>
      </div>

      <div v-if="rejected" class="rejected-plate" aria-hidden="true">СИГНАЛ НЕ РОЗПІЗНАНО</div>

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
  </section>
</template>

<style scoped>
.telegraph-page {
  min-height: 100vh;
  min-height: 100dvh;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  box-sizing: border-box;
  background:
    radial-gradient(ellipse at 50% 15%, #4a3520 0%, #1a1108 70%),
    repeating-linear-gradient(90deg, rgba(0,0,0,0.08) 0px, transparent 2px, transparent 40px);
}

.telegraph-panel {
  position: relative;
  width: 100%;
  max-width: 520px;
  background:
    linear-gradient(160deg, #caa969 0%, #a3813f 45%, #8a6a30 100%);
  border: 6px solid #2b1c0f;
  border-radius: 8px;
  box-shadow:
    0 20px 50px rgba(0, 0, 0, 0.6),
    inset 0 0 40px rgba(43, 28, 15, 0.35);
  padding: clamp(24px, 5vw, 42px);
  outline: none;
}

.rivet {
  position: absolute;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, #f0dca0, #6b4f22 75%);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.6);
}
.rivet-tl { top: 10px; left: 10px; }
.rivet-tr { top: 10px; right: 10px; }
.rivet-bl { bottom: 10px; left: 10px; }
.rivet-br { bottom: 10px; right: 10px; }

.plate-id {
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  color: #3a2a12;
  opacity: 0.75;
  margin: 0 0 14px;
}
.plate-title {
  font-family: 'Oswald', sans-serif;
  font-weight: 700;
  font-size: clamp(20px, 5vw, 26px);
  letter-spacing: 1px;
  color: #2b1c0f;
  text-shadow: 0 1px 0 rgba(255, 234, 180, 0.35);
  margin: 0 0 14px;
}
.plate-text {
  font-family: 'Special Elite', monospace;
  font-size: clamp(13px, 3vw, 15px);
  line-height: 1.7;
  color: #3a2a12;
  margin: 0 0 24px;
}

.lamp-row {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  margin-bottom: 22px;
}
.lamp-housing {
  width: 68px;
  height: 68px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, #e8cf90, #6b4f22 75%);
  border: 3px solid #2b1c0f;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: inset 0 2px 6px rgba(0, 0, 0, 0.5);
}
.lamp-bulb {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: radial-gradient(circle at 40% 35%, #6e2c14, #2c0f06 80%);
  border: 2px solid #1a0a04;
  transition: background 0.05s linear, box-shadow 0.05s linear;
}
.lamp-bulb.is-on {
  background: radial-gradient(circle at 40% 35%, #ffe0a0, #ff7a1a 65%, #b73a0a 100%);
  box-shadow: 0 0 26px 10px rgba(255, 138, 40, 0.85);
}
.lamp-label {
  font-family: 'Share Tech Mono', monospace;
  font-size: 11px;
  letter-spacing: 2px;
  color: #2b1c0f;
  opacity: 0.75;
  margin: 0;
}

.ticker-tape {
  text-align: center;
  margin: 0 auto 22px;
  max-width: 340px;
  background: #ece7d8;
  background-image:
    repeating-linear-gradient(0deg, rgba(0,0,0,0.02) 0px, transparent 1px, transparent 2px);
  border: 1px solid #6b4f22;
  border-radius: 2px;
  padding: 14px 10px;
  min-height: 56px;
  box-shadow: inset 0 2px 6px rgba(0, 0, 0, 0.25);
}
.ticker-word {
  font-family: 'Special Elite', monospace;
  font-weight: 700;
  font-size: clamp(20px, 5vw, 26px);
  letter-spacing: 4px;
  color: #1c1a14;
  margin: 0 0 4px;
}
.ticker-buffer {
  font-family: 'Share Tech Mono', monospace;
  font-size: 16px;
  letter-spacing: 3px;
  color: #6b6252;
  min-height: 1.2em;
  margin: 0;
}

.key-assembly {
  position: relative;
  width: 180px;
  height: 120px;
  margin: 0 auto 10px;
}
.key-base {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 170px;
  height: 42px;
  background: linear-gradient(180deg, #3a2a16, #1c130a);
  border: 2px solid #14100a;
  border-radius: 4px;
  box-shadow: 0 6px 10px rgba(0, 0, 0, 0.5);
}
.key-screw {
  position: absolute;
  top: 50%;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  transform: translateY(-50%);
  background: radial-gradient(circle at 35% 30%, #e8cf90, #6b4f22 75%);
}
.key-screw-l { left: 12px; }
.key-screw-r { right: 12px; }

.key-lever {
  position: absolute;
  bottom: 30px;
  left: 50%;
  transform-origin: 20% 90%;
  transform: translateX(-50%) rotate(0deg);
  width: 150px;
  height: 24px;
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  user-select: none;
  touch-action: none;
  transition: transform 0.08s ease;
}
.key-arm {
  position: absolute;
  left: 8px;
  top: 8px;
  width: 100px;
  height: 8px;
  border-radius: 4px;
  background: linear-gradient(180deg, #f0dca0, #a3813f 55%, #6b4f22);
  box-shadow: 0 2px 3px rgba(0, 0, 0, 0.5);
}
.key-knob {
  position: absolute;
  right: 0;
  top: -6px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, #2a2a2a, #0a0a0a 75%);
  border: 3px solid #4a2c14;
  box-shadow: 0 3px 5px rgba(0, 0, 0, 0.6);
}
.key-lever.is-pressed {
  transform: translateX(-50%) rotate(8deg);
}
.key-lever.is-pressed .key-knob {
  box-shadow: 0 3px 5px rgba(0, 0, 0, 0.6), 0 0 16px 4px rgba(255, 138, 40, 0.7);
}

.key-hint {
  text-align: center;
  font-family: 'Share Tech Mono', monospace;
  font-size: 12px;
  color: #3a2a12;
  opacity: 0.8;
  margin: 0 0 22px;
}

.brass-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}
.brass-btn {
  flex: 1;
  padding: 12px;
  font-family: 'Oswald', sans-serif;
  font-weight: 600;
  letter-spacing: 1.5px;
  font-size: clamp(12px, 2.8vw, 14px);
  cursor: pointer;
  border-radius: 3px;
  transition: transform 0.1s ease;
}
.brass-btn.ghost {
  background: transparent;
  border: 2px solid #2b1c0f;
  color: #2b1c0f;
}
.brass-btn.primary {
  flex: 2;
  background: linear-gradient(180deg, #3a2a16, #1c130a);
  border: 2px solid #14100a;
  color: #f0dca0;
}
.brass-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.brass-btn:active:not(:disabled) {
  transform: scale(0.97);
}

.rejected-plate {
  text-align: center;
  font-family: 'Oswald', sans-serif;
  font-weight: 700;
  letter-spacing: 2px;
  color: #7a1f10;
  background: rgba(236, 231, 216, 0.6);
  border: 1px dashed #7a1f10;
  padding: 6px;
  margin-bottom: 10px;
}

.hint-block {
  margin-top: 20px;
  padding-top: 18px;
  border-top: 1px dashed rgba(43, 28, 15, 0.4);
}
.hint-btn {
  width: 100%;
  padding: 10px;
  background: transparent;
  border: 1px dashed #3a2a12;
  color: #3a2a12;
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(11px, 2.8vw, 13px);
  letter-spacing: 1.5px;
  cursor: pointer;
  animation: hint-pulse 2s ease-in-out infinite;
}
.hint-btn:disabled {
  cursor: not-allowed;
  animation: none;
}
@keyframes hint-pulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}
.hint-label {
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1.5px;
  color: #3a2a12;
  opacity: 0.8;
  margin: 0 0 6px;
}
.hint-text {
  font-family: 'Special Elite', monospace;
  font-size: clamp(12px, 3vw, 14px);
  color: #1c1a14;
  margin: 0;
  min-height: 1.4em;
}
.hint-wait {
  margin-top: 10px;
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.6vw, 12px);
  color: #3a2a12;
  opacity: 0.8;
}

@media (prefers-reduced-motion: reduce) {
  .lamp-bulb, .key-lever { transition: none; }
}
</style>
