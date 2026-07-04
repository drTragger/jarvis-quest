<script setup>
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import JarvisLayout from '../layouts/JarvisLayout.vue'

const router = useRouter()
const input = ref('')
const lines = ref(['JARVIS ONLINE. Введіть команду або "help".'])
const logEl = ref(null)

const responses = {
  help: 'Доступні команди: status, joke, stark, exit',
  status: 'Всі системи функціонують в межах норми. Реактор стабільний.',
  joke: 'Сер, я обчислив 4 мільярди можливих результатів. Жоден не смішний.',
  stark: 'Останній відомий контакт: майстерня, Малібу. Дані застарілі.',
  exit: 'Термінал не можна закрити звідси. Скористайтесь кнопкою "Назад".'
}

async function run() {
  const cmd = input.value.trim().toLowerCase()
  if (!cmd) return
  lines.value.push(`> ${cmd}`)
  lines.value.push(responses[cmd] || `Команду "${cmd}" не розпізнано.`)
  input.value = ''
  await nextTick()
  logEl.value?.scrollTo({ top: logEl.value.scrollHeight })
}
</script>

<template>
  <JarvisLayout eyebrow="ТЕРМІНАЛ JARVIS" :show-scan="false">
    <button class="term-back" @click="router.push({ name: 'hub' })">&lt; НАЗАД ДО МЕНЮ</button>

    <div class="term-log" ref="logEl">
      <p v-for="(line, i) in lines" :key="i" class="term-line">{{ line }}</p>
    </div>

    <input
      v-model="input"
      class="jarvis-input term-input"
      placeholder="ВВЕДІТЬ КОМАНДУ..."
      autocomplete="off"
      @keyup.enter="run"
    />
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
.term-log {
  width: 100%;
  max-height: 240px;
  overflow-y: auto;
  margin-bottom: 16px;
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  color: var(--jarvis-text-dim);
}
.term-line { margin: 0 0 6px; }
.term-input { width: 100%; }
</style>