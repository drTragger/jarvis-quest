<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { useQuestStep } from '../composables/useQuestStep'
import { useHint } from '../composables/useHint'
import StepSuccessOverlay from '../components/StepSuccessOverlay.vue'
import { playSound } from '../composables/useSound'
import step4Voice from '../assets/audio/step4-voice.mp3'
import CaptchaMath from '../components/captcha/CaptchaMath.vue'
import CaptchaOddOneOut from '../components/captcha/CaptchaOddOneOut.vue'
import CaptchaGlyphs from '../components/captcha/CaptchaGlyphs.vue'
import CaptchaMemory from '../components/captcha/CaptchaMemory.vue'
import CaptchaSlider from '../components/captcha/CaptchaSlider.vue'

const STEP_ID = 4
const SECRET_LETTERS = ['R', 'O', 'B', 'O', 'T']
const ROUND_COMPONENTS = [CaptchaMath, CaptchaOddOneOut, CaptchaGlyphs, CaptchaMemory, CaptchaSlider]
const TOTAL_ROUNDS = SECRET_LETTERS.length

const { loading, error, success, submitAnswer, proceedAfterSuccess } = useQuestStep(STEP_ID)
const { scanning, revealed, hintText, remaining, cooldown, checkHint } = useHint(STEP_ID)

const round = ref(1)
const revealedLetters = ref([])
const gauntletDone = computed(() => revealedLetters.value.length === TOTAL_ROUNDS)
const currentComponent = computed(() => ROUND_COMPONENTS[round.value - 1])

const captchaRejected = ref(false)
let audioCtx = null

const answer = ref('')
const rejected = ref(false)

function beep(freq, duration = 0.12) {
  try {
    audioCtx = audioCtx || new (window.AudioContext || window.webkitAudioContext)()
    const osc = audioCtx.createOscillator()
    const gain = audioCtx.createGain()
    osc.frequency.value = freq
    osc.type = 'square'
    gain.gain.setValueAtTime(0.05, audioCtx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.0001, audioCtx.currentTime + duration)
    osc.connect(gain)
    gain.connect(audioCtx.destination)
    osc.start()
    osc.stop(audioCtx.currentTime + duration)
  } catch {
    // audio feedback is a nice-to-have, ignore if unsupported
  }
}

function onRoundSolved() {
  beep(880)
  revealedLetters.value = [...revealedLetters.value, SECRET_LETTERS[round.value - 1]]
  gsap.fromTo(
    `.letter-slot-${round.value - 1}`,
    { rotateX: -90, opacity: 0 },
    { rotateX: 0, opacity: 1, duration: 0.5, ease: 'back.out(1.7)' }
  )
  if (!gauntletDone.value) round.value += 1
}

function onRoundFailed() {
  beep(140, 0.2)
  captchaRejected.value = true
  gsap.timeline()
    .fromTo('.captcha-stage', { x: -6 }, { x: 0, duration: 0.4, ease: 'elastic.out(1, 0.4)' })
    .eventCallback('onComplete', () => {
      setTimeout(() => { captchaRejected.value = false }, 800)
    })
}

async function submit() {
  if (!answer.value.trim() || loading.value) return
  const ok = await submitAnswer(answer.value)
  if (!ok) {
    rejected.value = true
    gsap.timeline()
      .fromTo('.frame-content', { x: -5 }, { x: 0, duration: 0.4, ease: 'elastic.out(1, 0.4)' })
      .eventCallback('onComplete', () => {
        setTimeout(() => { rejected.value = false }, 900)
      })
  }
}

onMounted(() => {
  gsap.from('.frame-content', { opacity: 0, y: 16, duration: 0.5 })
  playSound('step4-voice', step4Voice, { volume: 0.9 })
})

onBeforeUnmount(() => {
  if (audioCtx) audioCtx.close()
})
</script>

<template>
  <JarvisLayout eyebrow="ПРОТОКОЛ ПЕРЕВІРКИ ЛЮДЯНОСТІ" :warn="rejected || captchaRejected">
    <StepSuccessOverlay v-if="success" @done="proceedAfterSuccess" />

    <h1 class="title">ПІДТВЕРДЬ, ЩО ТИ НЕ РОБОТ</h1>
    <p class="subtitle">
      Перш ніж продовжити, система вимагає пройти {{ TOTAL_ROUNDS }} перевірок розпізнавання.
      Кожна наступна перевірка складніша за попередню. Кожна пройдена перевірка відновлює один символ
      прихованої послідовності.
    </p>

    <div class="progress-line">СЕГМЕНТ {{ Math.min(round, TOTAL_ROUNDS) }} / {{ TOTAL_ROUNDS }}</div>

    <div class="letter-slots" aria-hidden="true">
      <span
        v-for="(letter, i) in SECRET_LETTERS"
        :key="i"
        class="letter-slot"
        :class="[`letter-slot-${i}`, { filled: i < revealedLetters.length }]"
      >{{ i < revealedLetters.length ? revealedLetters[i] : '_' }}</span>
    </div>

    <template v-if="!gauntletDone">
      <div class="captcha-stage">
        <component :is="currentComponent" :key="round" @solved="onRoundSolved" @failed="onRoundFailed" />
      </div>
      <p v-if="captchaRejected" class="rejected-text">НЕ РОЗПІЗНАНО</p>
    </template>

    <template v-else>
      <div class="unlocked-block">
        <p class="unlocked-label">ПОСЛІДОВНІСТЬ ВІДНОВЛЕНА</p>
        <p class="unlocked-hint">Введи відновлене слово, щоб підтвердити перевірку.</p>
      </div>

      <input
        v-model="answer"
        class="jarvis-input"
        type="text"
        placeholder="ВВЕДІТЬ ВІДПОВІДЬ"
        :class="{ 'is-error': error }"
        autocomplete="off"
        @keyup.enter="submit"
      />
      <button class="jarvis-btn" :disabled="loading" @click="submit">
        <span>{{ loading ? 'ПЕРЕВІРКА...' : 'ПІДТВЕРДИТИ' }}</span>
      </button>

      <p v-if="rejected" class="rejected-text">ВІДХИЛЕНО</p>
    </template>

    <div class="hint-block">
      <template v-if="!revealed">
        <button class="hint-btn" :disabled="scanning || cooldown" @click="checkHint">
          {{ scanning ? 'АНАЛІЗ...' : 'ЗАПУСТИТИ ПОГЛИБЛЕНИЙ АНАЛІЗ' }}
        </button>
        <p v-if="remaining !== null" class="hint-wait">
          Аналіз ще не завершено. Спробуй за {{ remaining }} с., або дай ще одну відповідь.
        </p>
      </template>
      <template v-else>
        <p class="hint-label">РЕЗУЛЬТАТ АНАЛІЗУ:</p>
        <p class="hint-text">{{ hintText }}</p>
      </template>
    </div>
  </JarvisLayout>
</template>

<style scoped>
.title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(18px, 4.4vw, 23px);
  margin: 0 0 12px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.subtitle {
  font-size: clamp(12px, 2.8vw, 14px);
  line-height: 1.6;
  color: var(--jarvis-text-dim);
  margin: 0 0 18px;
}
.progress-line {
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 2px;
  color: var(--jarvis-cyan);
  margin-bottom: 14px;
}
.letter-slots {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}
.letter-slot {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(126, 20, 255, 0.4);
  border-radius: 4px;
  font-family: var(--jarvis-mono);
  font-weight: 700;
  font-size: 16px;
  color: var(--jarvis-text-dim);
  background: rgba(255, 255, 255, 0.03);
}
.letter-slot.filled {
  color: var(--jarvis-cyan);
  border-color: var(--jarvis-cyan);
  box-shadow: 0 0 12px rgba(71, 191, 255, 0.4);
  text-shadow: 0 0 10px rgba(71, 191, 255, 0.6);
}
.captcha-stage {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-bottom: 6px;
}
.jarvis-btn {
  margin-top: 4px;
}
.rejected-text {
  color: var(--jarvis-danger);
  font-size: 12px;
  letter-spacing: 3px;
  margin-top: 12px;
}
.unlocked-block {
  width: 100%;
  padding: 14px 16px;
  margin-bottom: 16px;
  background: rgba(71, 191, 255, 0.06);
  border: 1px solid rgba(71, 191, 255, 0.4);
  border-radius: 6px;
}
.unlocked-label {
  font-family: var(--jarvis-mono);
  font-size: 12px;
  letter-spacing: 2px;
  color: var(--jarvis-cyan);
  margin: 0 0 6px;
}
.unlocked-hint {
  font-size: 12px;
  color: var(--jarvis-text-dim);
  margin: 0;
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
