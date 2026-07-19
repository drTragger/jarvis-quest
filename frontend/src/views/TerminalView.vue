<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { usePlatformDetect } from '../composables/usePlatformDetect'
import { getPassword } from '../lib/auth'
import { invalidateProgressCache } from '../router'

const { platform } = usePlatformDetect()

const router = useRouter()
const input = ref('')
const lines = ref(['JARVIS ONLINE. Введіть команду або "help".'])
const logEl = ref(null)
const inputEl = ref(null)
const currentStep = ref(null)

const responses = {
  help: 'Доступні команди: status, joke, stark, rules, exit',
  status: 'Всі системи функціонують в межах норми. Реактор стабільний.',
  joke: 'Сер, я обчислив 4 мільярди можливих результатів. Жоден не смішний.',
  stark: 'Останній відомий контакт: майстерня, Малібу. Дані застарілі.',
  rules: 'Відкриваю протокол доступу...',
  exit: 'Термінал не можна закрити звідси. Скористайтесь кнопкою "Назад".'
}

async function loadCurrentStep() {
  const password = getPassword()
  if (!password) return
  try {
    const res = await fetch(`/api/progress?password=${encodeURIComponent(password)}`)
    if (res.status === 401) return
    const data = await res.json()
    currentStep.value = data.unlocked ?? null
  } catch {
    currentStep.value = null
  }
}

async function tryAsQuestAnswer(cmd) {
  if (!currentStep.value) return false
  const password = getPassword()
  const res = await fetch('/api/answer', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      password,
      step_id: `step${currentStep.value}`,
      answer: cmd
    })
  })
  if (res.status === 401) {
    router.push({ name: 'login' })
    return true
  }
  const data = await res.json()
  if (data.correct) {
    invalidateProgressCache()
    loadCurrentStep()
    return true
  }
  return false
}

async function run() {
  const cmd = input.value.trim().toLowerCase()
  if (!cmd) return
  lines.value.push(`marmelade@jarvis:~$ ${cmd}`)
  input.value = ''

  if (cmd in responses) {
    lines.value.push(responses[cmd])
  } else {
    const accepted = await tryAsQuestAnswer(cmd)
    lines.value.push(accepted ? 'КОМАНДУ ВИКОНАНО. ПРОТОКОЛ МІСІЇ ОНОВЛЕНО.' : `Команду "${cmd}" не розпізнано.`)
  }

  await nextTick()
  logEl.value?.scrollTo({ top: logEl.value.scrollHeight })

  if (cmd === 'rules') {
    setTimeout(() => router.push({ name: 'rules' }), 500)
  }
}

function focusInput() {
  inputEl.value?.focus()
}

onMounted(() => {
  // Programmatic focus doesn't open the virtual keyboard on mobile anyway
  // (and it made the caret render oddly before the user's first tap), so
  // only autofocus on devices where it's actually helpful.
  if (platform.value !== 'mobile') focusInput()
  loadCurrentStep()
})
</script>

<template>
  <JarvisLayout eyebrow="ТЕРМІНАЛ JARVIS" :show-scan="false">
    <button class="term-back" @click="router.push({ name: 'hub' })">&lt; НАЗАД ДО МЕНЮ</button>

    <div class="terminal-window">
      <div class="terminal-titlebar">
        <span class="term-dot term-dot-red"></span>
        <span class="term-dot term-dot-yellow"></span>
        <span class="term-dot term-dot-green"></span>
        <span class="terminal-title">marmelade@jarvis: ~</span>
      </div>

      <div class="term-log" ref="logEl" @click="focusInput" @touchstart="focusInput">
        <p v-for="(line, i) in lines" :key="i" class="term-line">{{ line }}</p>
        <div class="term-prompt-row">
          <span class="term-prompt">marmelade@jarvis:~$</span>
          <input
            ref="inputEl"
            v-model="input"
            class="term-input-inline"
            autocomplete="off"
            autocapitalize="off"
            autocorrect="off"
            spellcheck="false"
            @keyup.enter="run"
          />
        </div>
      </div>
    </div>
  </JarvisLayout>
</template>

<style scoped>
.term-back {
  align-self: flex-start;
  background: none;
  border: none;
  color: var(--jarvis-cyan);
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  cursor: pointer;
  margin-bottom: 16px;
  padding: 0;
}

.terminal-window {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
  background: #0a0c0a;
  border: 1px solid rgba(51, 255, 102, 0.25);
  box-shadow: 0 0 24px rgba(51, 255, 102, 0.08), 0 8px 24px rgba(0, 0, 0, 0.5);
}

.terminal-titlebar {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 9px 12px;
  background: linear-gradient(180deg, #2a2a2a, #1c1c1c);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.term-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}
.term-dot-red { background: #ff5f56; }
.term-dot-yellow { background: #ffbd2e; }
.term-dot-green { background: #27c93f; }
.terminal-title {
  flex: 1;
  text-align: center;
  font-family: var(--jarvis-mono);
  font-size: 11px;
  color: #8a8a8a;
  margin-right: 46px;
}

.term-log {
  position: relative;
  width: 100%;
  height: 280px;
  overflow-y: auto;
  padding: 14px 16px;
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  line-height: 1.7;
  color: #3ee06a;
  text-align: left;
  cursor: text;
  background-image: repeating-linear-gradient(
    0deg,
    rgba(255, 255, 255, 0.025) 0px,
    rgba(255, 255, 255, 0.025) 1px,
    transparent 1px,
    transparent 3px
  );
}
.term-line {
  margin: 0 0 6px;
  word-break: break-word;
  text-shadow: 0 0 6px rgba(51, 255, 102, 0.35);
}
.term-prompt-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.term-prompt {
  flex-shrink: 0;
  color: #3ee06a;
  text-shadow: 0 0 6px rgba(51, 255, 102, 0.35);
}
.term-input-inline {
  flex: 1;
  min-width: 0;
  background: none;
  border: none;
  outline: none;
  margin: 0;
  padding: 0;
  -webkit-appearance: none;
  appearance: none;
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  line-height: 1.7;
  color: #eafff0;
  caret-color: #3ee06a;
}

@media (max-width: 560px) {
  /* iOS Safari auto-zooms the page on focus if an input's font-size is
     under 16px, which is what made the caret look misplaced/jumpy. */
  .term-input-inline {
    font-size: 16px;
  }
}
</style>
