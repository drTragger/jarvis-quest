<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { getPassword } from '../lib/auth'
import { invalidateProgressCache, resolveDestination } from '../router'
import { usePlatformDetect } from '../composables/usePlatformDetect'
import { playSound } from '../composables/useSound'
import soundTestFile from '../assets/audio/rules-sound-test.mp3'

const router = useRouter()
const loading = ref(false)

const { platform } = usePlatformDetect()
const isMobile = computed(() => platform.value === 'mobile')

const testingSound = ref(false)

function testSound() {
  if (testingSound.value) return
  testingSound.value = true
  const sound = playSound('rules-sound-test', soundTestFile, { volume: 0.9 })
  const stop = () => { testingSound.value = false }
  sound.once('end', stop)
  sound.once('loaderror', stop)
  sound.once('playerror', stop)
}

async function start() {
  if (loading.value) return
  loading.value = true

  await fetch('/api/rules/seen', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password: getPassword() })
  }).catch(() => {})

  invalidateProgressCache()
  const dest = await resolveDestination()
  router.push(dest)
}

onMounted(() => {
  gsap.from('.frame-content', { opacity: 0, y: 16, duration: 0.5, clearProps: 'transform' })
  gsap.from('.rule-item', { opacity: 0, y: 12, duration: 0.4, stagger: 0.1, delay: 0.2 })
})
</script>

<template>
  <JarvisLayout eyebrow="ПРОТОКОЛ ДОСТУПУ" wide>
    <h1 class="title">ПЕРШ НІЖ РОЗПОЧАТИ</h1>
    <p class="subtitle">
      Владе, перш ніж я відкрию доступ до системи — кілька правил користування системою.
    </p>

    <div class="rules-list">
      <div class="rule-item">
        <span class="rule-index">01</span>
        <div class="rule-body">
          <p class="rule-title">Гугли, рахуй, питай ШІ</p>
          <p class="rule-text">
            Пошук в інтернеті, калькулятор, навіть штучний інтелект — усе дозволено.
            Я нічого не перевіряю і не караю. Але чим більше думаєш сам, тим цікавішим буде фінал.
          </p>
        </div>
      </div>

      <div class="rule-item">
        <span class="rule-index">02</span>
        <div class="rule-body">
          <p class="rule-title">Телефон замість комп'ютера</p>
          <p class="rule-text">
            Протокол розрахований насамперед на мобільний інтерфейс. З телефону все
            виглядатиме так, як задумано.
          </p>
          <span class="device-badge" :class="{ 'is-good': isMobile, 'is-bad': !isMobile }">
            {{ isMobile ? '✓ ТИ НА ТЕЛЕФОНІ. ВСЕ ДОБРЕ' : '⨯ ЦЕ НЕ ТЕЛЕФОН' }}
          </span>
        </div>
      </div>

      <div class="rule-item">
        <span class="rule-index">03</span>
        <div class="rule-body">
          <p class="rule-title">Увімкни звук</p>
          <p class="rule-text">
            Частина інформації передається голосом — без звуку ти пропустиш половину суті.
            Якщо звук не грає, спробуй відкрити цю сторінку в новій вкладці браузера: Safari
            з цим протоколом ладнає гірше за інших, тож краще обрати щось інше.
          </p>
          <button class="sound-test-btn" :disabled="testingSound" @click="testSound">
            {{ testingSound ? '▶ ЗВУЧИТЬ...' : '▶ ПЕРЕВІРИТИ ЗВУК' }}
          </button>
        </div>
      </div>

      <div class="rule-item">
        <span class="rule-index">04</span>
        <div class="rule-body">
          <p class="rule-title">Застряг — не панікуй</p>
          <p class="rule-text">
            У кожному розділі є кнопка «ЗАПУСТИТИ ПОГЛИБЛЕНИЙ АНАЛІЗ» — це і є підказка.
            Але вона не з'явиться за одну секунду: доступ до неї відкривається лише через
            24 години після того, як ти вперше зайдеш у розділ. Це не збій і не жарт — просто
            дай собі час подумати самостійно. А коли час настане, підказка вже чекатиме.
          </p>
        </div>
      </div>

      <div class="rule-item">
        <span class="rule-index">05</span>
        <div class="rule-body">
          <p class="rule-title">Протокол самодостатній</p>
          <p class="rule-text">
            Усе необхідне для проходження вже сказано в самому інтерфейсі. Писати мені з
            питаннями "що робити далі" не потрібно. Звертайся лише якщо знайшов справжню
            технічну помилку — щось реально не працює.
          </p>
        </div>
      </div>

      <div class="rule-item">
        <span class="rule-index">06</span>
        <div class="rule-body">
          <p class="rule-title">У системі є великодки</p>
          <p class="rule-text">
            Десь у інтерфейсі заховано кілька непотрібних для проходження дрібниць.
            Вони ні на що не впливають і не дають жодної переваги — просто прикол.
            Знайдеш — вважай бонусом.
          </p>
        </div>
      </div>
    </div>

    <button class="jarvis-btn" :disabled="loading" @click="start">
      <span>{{ loading ? 'ЗАВАНТАЖЕННЯ...' : 'ЗРОЗУМІВ. РОЗПОЧАТИ' }}</span>
    </button>

    <p class="ps-note">P.S. тобі глина😁 Але я в тебе вірю.</p>
  </JarvisLayout>
</template>

<style scoped>
.title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(20px, 5vw, 26px);
  margin: 0 0 12px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.subtitle {
  font-size: clamp(13px, 3vw, 15px);
  line-height: 1.6;
  color: var(--jarvis-text-dim);
  margin: 0 0 24px;
}
.rules-list {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 28px;
}
.rule-item {
  display: flex;
  gap: 14px;
  align-items: flex-start;
  padding: 14px 16px;
  background: rgba(13, 18, 32, 0.6);
  border: 1px solid rgba(126, 20, 255, 0.3);
  border-left: 3px solid var(--jarvis-cyan);
  border-radius: 4px;
  text-align: left;
}
.rule-index {
  flex-shrink: 0;
  font-family: var(--jarvis-mono);
  font-weight: 700;
  font-size: clamp(13px, 3vw, 15px);
  color: var(--jarvis-cyan);
  text-shadow: 0 0 10px rgba(71, 191, 255, 0.5);
  margin-top: 1px;
}
.rule-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}
.rule-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(13px, 3.2vw, 15px);
  color: var(--jarvis-text);
  margin: 0;
}
.rule-text {
  font-size: clamp(12px, 2.8vw, 13px);
  line-height: 1.55;
  color: var(--jarvis-text-dim);
  margin: 0;
}
.device-badge {
  display: inline-block;
  margin-top: 8px;
  padding: 5px 10px;
  border-radius: 3px;
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 11px);
  letter-spacing: 0.5px;
  width: fit-content;
}
.device-badge.is-good {
  color: var(--jarvis-cyan);
  border: 1px solid rgba(71, 191, 255, 0.5);
  background: rgba(71, 191, 255, 0.08);
}
.device-badge.is-bad {
  color: var(--jarvis-danger);
  border: 1px solid rgba(255, 84, 112, 0.5);
  background: rgba(255, 84, 112, 0.08);
}
.sound-test-btn {
  align-self: flex-start;
  margin-top: 8px;
  background: none;
  border: 1px solid rgba(71, 191, 255, 0.5);
  border-radius: 4px;
  padding: 7px 12px;
  color: var(--jarvis-cyan);
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 11px);
  letter-spacing: 0.5px;
  cursor: pointer;
  transition: background 0.15s;
}
.sound-test-btn:hover:not(:disabled) {
  background: rgba(71, 191, 255, 0.12);
}
.sound-test-btn:disabled {
  opacity: 0.6;
  cursor: default;
}
.ps-note {
  margin-top: 18px;
  font-size: clamp(12px, 2.8vw, 13px);
  color: var(--jarvis-text-dim);
  font-style: italic;
}
</style>
