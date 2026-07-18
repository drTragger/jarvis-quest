<script setup>
import { ref, computed, onMounted } from 'vue'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { usePlatformDetect } from '../composables/usePlatformDetect'
import { useQuestStep } from '../composables/useQuestStep'
import { useHint } from '../composables/useHint'
import { playSound } from '../composables/useSound'
import StepSuccessOverlay from '../components/StepSuccessOverlay.vue'
import lockedVoice from '../assets/audio/step3-locked-voice.mp3'
import unlockedVoice from '../assets/audio/step3-unlocked-voice.mp3'
import errorSound from '../assets/audio/access-denied.wav'

const STEP_ID = 3
const REQUIRED_PLATFORM = 'desktop'

const { platform } = usePlatformDetect()
const isLocked = computed(() => platform.value !== REQUIRED_PLATFORM)

const { loading, error, success, submitAnswer, proceedAfterSuccess } = useQuestStep(STEP_ID)
const { scanning, revealed, hintText, remainingLabel, cooldown, checkHint } = useHint(STEP_ID)
const answer = ref('')
const rejected = ref(false)

async function submit() {
  if (!answer.value.trim() || loading.value) return
  const ok = await submitAnswer(answer.value)
  if (!ok) {
    rejected.value = true
    playSound('access-denied', errorSound, { volume: 0.6 })
    gsap.timeline()
      .fromTo('.frame-content', { x: -5 }, { x: 0, duration: 0.4, ease: 'elastic.out(1, 0.4)', clearProps: 'transform' })
      .eventCallback('onComplete', () => {
        setTimeout(() => { rejected.value = false }, 900)
      })
  }
}

onMounted(() => {
  gsap.from('.frame-content', { opacity: 0, y: 16, duration: 0.5, clearProps: 'transform' })
  playSound(
    isLocked.value ? 'step3-locked' : 'step3-unlocked',
    isLocked.value ? lockedVoice : unlockedVoice,
    { volume: 0.9 }
  )
})
</script>

<template>
  <JarvisLayout eyebrow="ПРОТОКОЛ ПЕРЕДАЧІ СИГНАЛУ" :alert="isLocked" :warn="rejected">
    <StepSuccessOverlay v-if="success" @done="proceedAfterSuccess" />
    <template v-if="isLocked">
      <div class="lock-icon" aria-hidden="true">⨯</div>
      <h1 class="title">СИГНАЛ НЕ РОЗПІЗНАНО</h1>
      <p class="subtitle">
        Приймач цього типу не призначений для розшифрування.
      </p>
      <p class="cryptic-text">
        Джерело — {{ REQUIRED_PLATFORM === 'desktop' ? 'стаціонарне' : 'мобільне' }}.
      </p>
      <div class="static-line" aria-hidden="true"></div>
    </template>

    <template v-else>
      <h1 class="title">СИГНАЛ ПРИЙНЯТО</h1>
      <p class="subtitle">Канал стабільний. Розшифруйте послідовність, щоб продовжити.</p>

      <div class="cipher-block">
        <p class="cipher-label">ПЕРЕХОПЛЕНО ASCII-ПОСЛІДОВНІСТЬ</p>
        <p class="cipher-hex">30 76 65 72 72 69 64 65</p>
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
    </template>
  </JarvisLayout>
</template>

<style scoped>
.title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(19px, 4.8vw, 24px);
  margin: 0 0 14px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.subtitle {
  font-size: clamp(13px, 3vw, 15px);
  line-height: 1.6;
  color: var(--jarvis-text-dim);
  margin: 0 0 20px;
}
.lock-icon {
  font-size: 48px;
  color: var(--jarvis-danger);
  margin-bottom: 12px;
}
.cryptic-text {
  font-family: var(--jarvis-mono);
  font-size: 13px;
  letter-spacing: 2px;
  color: var(--jarvis-text-dim);
  margin-top: 4px;
}
.static-line {
  width: 60%;
  height: 1px;
  margin-top: 20px;
  background: repeating-linear-gradient(
    90deg,
    var(--jarvis-danger) 0px,
    var(--jarvis-danger) 4px,
    transparent 4px,
    transparent 9px
  );
  opacity: 0.5;
}
.cipher-block {
  width: 100%;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(126, 20, 255, 0.35);
  padding: 16px;
  text-align: left;
  margin-bottom: 20px;
  box-sizing: border-box;
}
.cipher-label {
  font-size: 10px;
  letter-spacing: 2px;
  color: var(--jarvis-cyan);
  margin: 0 0 10px;
}
.cipher-hex {
  font-family: var(--jarvis-mono);
  font-size: 16px;
  letter-spacing: 4px;
  color: var(--jarvis-text);
  margin: 0;
}
.jarvis-btn {
  margin-top: 16px;
}
.rejected-text {
  color: var(--jarvis-danger);
  font-size: 12px;
  letter-spacing: 3px;
  margin-top: 12px;
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