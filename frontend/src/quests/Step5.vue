<script setup>
import { ref, onMounted } from 'vue'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { useQuestStep } from '../composables/useQuestStep'
import { useHint } from '../composables/useHint'
import StepSuccessOverlay from '../components/StepSuccessOverlay.vue'
import { playSound } from '../composables/useSound'
import step5Voice from '../assets/audio/step5-voice.mp3'
import errorSound from '../assets/audio/access-denied.wav'
import ironmanImage from '../assets/images/step5-ironman.png'

const STEP_ID = 5

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
  playSound('step5-voice', step5Voice, { volume: 0.9 })
})
</script>

<template>
  <JarvisLayout eyebrow="ПЕРЕХОПЛЕНЕ ЗОБРАЖЕННЯ" :warn="rejected">
    <StepSuccessOverlay v-if="success" @done="proceedAfterSuccess" />

    <h1 class="title">ЗОБРАЖЕННЯ БЕЗ ВИДИМИХ СЛІДІВ</h1>
    <p class="subtitle">
      Архів перехопив це зображення. На вигляд — нічого особливого, самі лише
      тіні та звичні кольори. Але дані рідко зникають повністю. Іноді достатньо
      трохи інакше подивитись на світло й тінь, щоб побачити те, що ховається в темряві.
    </p>

    <div class="image-panel">
      <img :src="ironmanImage" alt="Перехоплене зображення" class="scan-image" />
    </div>

    <a :href="ironmanImage" download="intercepted-image.png" class="image-download">
      &darr; ЗАВАНТАЖИТИ ЗОБРАЖЕННЯ ДЛЯ АНАЛІЗУ
    </a>

    <input
      v-model="answer"
      class="jarvis-input answer-input"
      type="text"
      placeholder="ВВЕДІТЬ РОЗШИФРОВАНЕ СЛОВО"
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
  </JarvisLayout>
</template>

<style scoped>
.title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(18px, 4.4vw, 23px);
  margin: 0 0 14px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.subtitle {
  font-size: clamp(12px, 2.8vw, 14px);
  line-height: 1.6;
  color: var(--jarvis-text-dim);
  margin: 0 0 20px;
}
.image-panel {
  width: 100%;
  max-width: 320px;
  margin-bottom: 10px;
  border: 1px solid rgba(126, 20, 255, 0.4);
  border-radius: 6px;
  overflow: hidden;
}
.scan-image {
  width: 100%;
  height: auto;
  display: block;
}
.image-download {
  display: inline-block;
  margin-bottom: 20px;
  color: var(--jarvis-cyan);
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  letter-spacing: 1px;
  text-decoration: none;
  border-bottom: 1px dashed var(--jarvis-cyan);
}
.image-download:hover,
.image-download:focus-visible {
  border-bottom-style: solid;
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
.answer-input {
  margin-top: 6px;
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
