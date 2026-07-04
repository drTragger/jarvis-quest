<script setup>
import { ref, onMounted } from 'vue'
import gsap from 'gsap'
import { useQuestStep } from '../composables/useQuestStep'
import { useHint } from '../composables/useHint'

const { loading, error, submitAnswer } = useQuestStep(1)
const { scanning, revealed, hintText, remaining, cooldown, checkHint } = useHint(1)
const answer = ref('')
const rejected = ref(false)

onMounted(() => {
  const tl = gsap.timeline({ defaults: { ease: 'power2.out' } })
  tl.from('.dossier', { opacity: 0, y: 24, rotate: -4, duration: 0.6 })
    .from('.stamp', { opacity: 0, scale: 2, rotate: -25, duration: 0.4, ease: 'back.out(3)' }, '-=0.2')
})

async function submit() {
  if (!answer.value.trim() || loading.value) return
  const ok = await submitAnswer(answer.value)
  if (!ok) {
    rejected.value = true
    gsap.timeline()
      .fromTo('.dossier', { x: -5 }, { x: 0, duration: 0.4, ease: 'elastic.out(1, 0.4)' })
      .fromTo('.rejected-stamp', { opacity: 0, scale: 1.6 }, { opacity: 1, scale: 1, duration: 0.25 }, 0)
      .to('.rejected-stamp', { opacity: 0, duration: 0.4, delay: 1.2 })
  }
}
</script>

<template>
  <section class="dossier-page">
    <div class="dossier">
      <div class="stamp" aria-hidden="true">
        <span>ЗАСЕКРЕЧЕНО</span>
      </div>

      <p class="file-id">АРХІВНИЙ ЗАПИС № 001&nbsp;&nbsp;·&nbsp;&nbsp;РІВЕНЬ ДОСТУПУ: ALPHA</p>

      <h1 class="dossier-title">Розшифровка протоколу</h1>

      <p class="dossier-text">
        Ти вже бачив це ім'я на екрані входу. Кожна літера абревіатури —
        окреме слово розшифровки. Архів зберіг майже всю послідовність,
        окрім одного фрагмента — його вирізали під час чистки даних.
      </p>

      <p class="acronym-line">
        J.A.R.V.I.S. = Just&nbsp;A&nbsp;Rather&nbsp;Very&nbsp;<span class="redacted" aria-hidden="true">██████████</span>&nbsp;System
      </p>

      <p class="dossier-hint">Відновити витерте слово. Англійською, малими літерами, без пробілів.</p>

      <div class="answer-row">
        <label class="answer-label" for="step1-answer">ВІДНОВЛЕНЕ СЛОВО:</label>
        <input
          id="step1-answer"
          v-model="answer"
          class="dossier-input"
          :class="{ 'is-error': error }"
          autocomplete="off"
          @keyup.enter="submit"
        />
      </div>

      <button class="stamp-btn" :disabled="loading" @click="submit">
        {{ loading ? 'ПЕРЕВІРКА...' : 'ПІДТВЕРДИТИ' }}
      </button>

      <div v-if="rejected" class="rejected-stamp" aria-hidden="true">ВІДХИЛЕНО</div>

      <div class="hint-block">
        <template v-if="!revealed">
          <button class="hint-btn" :disabled="scanning || cooldown" @click="checkHint">
            {{ scanning ? 'СКАНУВАННЯ АРХІВУ...' : 'ЗАПУСТИТИ ДОДАТКОВЕ СКАНУВАННЯ' }}
          </button>
          <p v-if="remaining !== null" class="hint-wait">
            Архів ще не розшифрував фрагмент. Спробуй за {{ remaining }} с., або дай ще одну відповідь.
          </p>
        </template>
        <template v-else>
          <p class="hint-label">ЧАСТКОВО ВІДНОВЛЕНІ ДАНІ:</p>
          <p class="hint-text">{{ hintText }}</p>
        </template>
      </div>
    </div>
  </section>
</template>

<style scoped>
.dossier-page {
  min-height: 100vh;
  min-height: 100dvh;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  box-sizing: border-box;
  background: radial-gradient(ellipse at 50% 20%, #2a2620 0%, #0f0d0a 75%);
}

.dossier {
  position: relative;
  width: 100%;
  max-width: 560px;
  background: #ece7d8;
  background-image:
    repeating-linear-gradient(0deg, rgba(0,0,0,0.015) 0px, transparent 1px, transparent 2px),
    radial-gradient(ellipse at 30% 20%, rgba(0,0,0,0.06), transparent 60%),
    radial-gradient(ellipse at 80% 90%, rgba(0,0,0,0.08), transparent 55%);
  color: #1c1a14;
  padding: clamp(24px, 5vw, 44px);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
  clip-path: polygon(
    0% 1%, 2% 0%, 98% 0%, 100% 1%,
    100% 98%, 99% 100%, 97% 99%, 3% 100%, 1% 99%, 0% 98%
  );
}

.stamp {
  position: absolute;
  top: clamp(14px, 4vw, 28px);
  right: clamp(14px, 4vw, 28px);
  border: 3px solid #9c1f1f;
  color: #9c1f1f;
  font-family: 'Oswald', sans-serif;
  font-weight: 700;
  font-size: clamp(10px, 2.4vw, 13px);
  letter-spacing: 2px;
  padding: 6px 10px;
  transform: rotate(-8deg);
  opacity: 0.85;
  mix-blend-mode: multiply;
}

.file-id {
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  color: #6b6252;
  margin: 0 0 18px;
}

.dossier-title {
  font-family: 'Special Elite', monospace;
  font-size: clamp(22px, 5.5vw, 30px);
  margin: 0 0 18px;
  line-height: 1.2;
}

.dossier-text {
  font-family: 'Special Elite', monospace;
  font-size: clamp(14px, 3.4vw, 16px);
  line-height: 1.7;
  margin: 0 0 20px;
}

.acronym-line {
  font-family: 'Special Elite', monospace;
  font-size: clamp(14px, 3.6vw, 17px);
  line-height: 1.8;
  margin: 0 0 20px;
  word-break: break-word;
}

.redacted {
  display: inline-block;
  background: #0b0b0b;
  color: transparent;
  border-radius: 2px;
  user-select: none;
}

.dossier-hint {
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(11px, 2.8vw, 13px);
  color: #6b6252;
  margin: 0 0 26px;
  font-style: italic;
}

.answer-row {
  margin-bottom: 20px;
}

.answer-label {
  display: block;
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1.5px;
  color: #6b6252;
  margin-bottom: 8px;
}

.dossier-input {
  width: 100%;
  box-sizing: border-box;
  background: transparent;
  border: none;
  border-bottom: 2px solid #1c1a14;
  padding: 8px 4px;
  font-family: 'Special Elite', monospace;
  font-size: clamp(15px, 3.6vw, 18px);
  color: #1c1a14;
  outline: none;
}
.dossier-input.is-error {
  border-bottom-color: #9c1f1f;
}

.stamp-btn {
  width: 100%;
  padding: 14px;
  background: transparent;
  border: 2px solid #1c1a14;
  color: #1c1a14;
  font-family: 'Oswald', sans-serif;
  font-weight: 600;
  font-size: clamp(13px, 3vw, 15px);
  letter-spacing: 3px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}
.stamp-btn:hover,
.stamp-btn:focus-visible {
  background: #1c1a14;
  color: #ece7d8;
  outline: none;
}
.stamp-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.rejected-stamp {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) rotate(-10deg);
  border: 4px solid #9c1f1f;
  color: #9c1f1f;
  font-family: 'Oswald', sans-serif;
  font-weight: 700;
  font-size: clamp(20px, 6vw, 32px);
  letter-spacing: 4px;
  padding: 10px 24px;
  background: rgba(236, 231, 216, 0.85);
  pointer-events: none;
  mix-blend-mode: multiply;
}

.hint-block {
  margin-top: 24px;
  padding-top: 18px;
  border-top: 1px dashed #a89f88;
}

.hint-btn {
  width: 100%;
  padding: 10px;
  background: transparent;
  border: 1px dashed #6b6252;
  color: #6b6252;
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
  0%, 100% { opacity: 0.7; }
  50% { opacity: 1; }
}

.hint-label {
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1.5px;
  color: #6b6252;
  margin: 0 0 6px;
}

.hint-text {
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(12px, 3vw, 14px);
  color: #1c1a14;
  margin: 0;
  min-height: 1.4em;
}

.hint-wait {
  margin-top: 10px;
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.6vw, 12px);
  color: #6b6252;
}

@media (prefers-reduced-motion: reduce) {
  .dossier, .stamp, .rejected-stamp, .hint-btn { animation: none !important; }
}
</style>