<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { setPassword } from '../lib/auth'
import { invalidateProgressCache } from '../router'
import { stepIds } from '../data/steps'

const router = useRouter()
const password = ref('')
const error = ref(false)
const loading = ref(false)
const granted = ref(false)

async function submit() {
  if (!password.value.trim() || loading.value) return
  loading.value = true
  error.value = false

  const res = await fetch(`/api/progress?password=${encodeURIComponent(password.value)}`)
  loading.value = false

  if (res.status === 401) {
    error.value = true
    gsap.fromTo('.login-input', { x: -6 }, { x: 0, duration: 0.4, ease: 'elastic.out(1, 0.4)' })
    return
  }

  const data = await res.json()
  setPassword(password.value)
  invalidateProgressCache()

  const dest = stepIds.includes(data.unlocked)
  ? { name: 'quest-step', params: { stepId: data.unlocked } }
  : { name: 'finish' }

  granted.value = true

  await gsap.timeline()
    .from('.granted-core', { opacity: 0, scale: 0.5, duration: 0.5, ease: 'back.out(2)' })
    .from('.granted-text', { opacity: 0, y: 10, duration: 0.4 }, '-=0.2')
    .to({}, { duration: 0.7 })

  router.push(dest)
}
</script>

<template>
  <JarvisLayout eyebrow="ІДЕНТИФІКАЦІЯ КОРИСТУВАЧА" :alert="error">
    <template v-if="!granted">
      <h2 class="login-title">Введіть код доступу</h2>
      <p class="login-hint">Система впізнає лише одну людину.</p>

      <input
        v-model="password"
        type="password"
        class="jarvis-input login-input"
        :class="{ 'is-error': error }"
        placeholder="КОД ДОСТУПУ"
        autocomplete="off"
        @keyup.enter="submit"
      />

      <button class="jarvis-btn login-btn" :disabled="loading" @click="submit">
        {{ loading ? 'ПЕРЕВІРКА...' : 'АВТОРИЗУВАТИ' }}
      </button>

      <p v-if="error" class="login-error">&gt; ДОСТУП ЗАБОРОНЕНО. НЕВІРНИЙ КОД.</p>
    </template>

    <template v-else>
      <div class="granted-core" aria-hidden="true">
        <div class="granted-ring granted-ring-outer"></div>
        <div class="granted-ring granted-ring-inner"></div>
        <div class="granted-dot"></div>
      </div>
      <p class="granted-text">ДОСТУП НАДАНО</p>
    </template>
  </JarvisLayout>
</template>

<style scoped>
.login-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(20px, 5vw, 26px);
  margin: 0 0 10px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.login-hint {
  font-size: clamp(12px, 3vw, 14px);
  color: var(--jarvis-text-dim);
  margin: 0 0 24px;
}
.login-input {
  margin-bottom: 16px;
  text-align: center;
}
.login-btn {
  margin-top: 4px;
}
.login-error {
  margin-top: 14px;
  font-size: clamp(11px, 2.8vw, 13px);
  color: var(--jarvis-danger);
}

.granted-core {
  position: relative;
  width: clamp(72px, 22vw, 104px);
  height: clamp(72px, 22vw, 104px);
  margin-bottom: 20px;
}
.granted-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid transparent;
}
.granted-ring-outer {
  border-top-color: var(--jarvis-cyan);
  border-bottom-color: var(--jarvis-cyan);
  animation: spin 1.4s linear infinite;
}
.granted-ring-inner {
  inset: 20%;
  border-left-color: var(--jarvis-purple);
  border-right-color: var(--jarvis-purple);
  animation: spin 0.9s linear infinite reverse;
}
.granted-dot {
  position: absolute;
  inset: 40%;
  border-radius: 50%;
  background: var(--jarvis-cyan);
  box-shadow: 0 0 20px 6px rgba(71, 191, 255, 0.8);
}

.granted-text {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(18px, 4.5vw, 24px);
  letter-spacing: 4px;
  color: var(--jarvis-text);
  text-shadow: 0 0 18px rgba(71, 191, 255, 0.7);
  margin: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (prefers-reduced-motion: reduce) {
  .granted-ring-outer,
  .granted-ring-inner {
    animation: none;
  }
}
</style>