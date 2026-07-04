<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { getPassword } from '../lib/auth'
import { stepIds } from '../data/steps'
import OnboardingOverlay from '../components/OnboardingOverlay.vue'
import { hasSeenOnboarding, markOnboardingSeen } from '../lib/onboarding'
import onboarding1 from '../assets/audio/onboarding-1.mp3'
import onboarding2 from '../assets/audio/onboarding-2.mp3'
import onboarding3 from '../assets/audio/onboarding-3.mp3'
import onboarding4 from '../assets/audio/onboarding-4.mp3'
import onboarding5 from '../assets/audio/onboarding-5.mp3'

const router = useRouter()
const unlocked = ref(null)
const loading = ref(true)
const uptimeSeconds = ref(0)
let uptimeTimer = null
const showOnboarding = ref(false)

const onboardingSteps = [
	{ selector: null, title: 'Ласкаво прошу', text: 'Вітаю, Владе. Це головна панель керування. Дозволь коротко показати, що тут є.', audioSrc: onboarding1 },
	{ selector: '[data-onboard="mission"]', title: 'Продовжити місію', text: 'Тут ти продовжуєш проходження місій. Кожна розкрита загадка наближає тебе до фінальної нагороди.', audioSrc: onboarding2 },
	{ selector: '[data-onboard="archive"]', title: 'Архів', text: 'У архіві зберігаються записи з уже пройдених місій. Можеш переглянути їх у будь-який момент.', audioSrc: onboarding3 },
	{ selector: '[data-onboard="terminal"]', title: 'Термінал', text: 'А це прямий канал зв\'язку зі мною. Він не обов\'язковий, але можеш звернутись, якщо стане нудно.', audioSrc: onboarding4 },
	{ selector: null, title: 'Все готово', text: 'Тепер усе в твоїх руках. Успіхів, Владе.', audioSrc: onboarding5 }
]

function onOnboardingFinish() {
	showOnboarding.value = false
	markOnboardingSeen()
}

const completedCount = computed(() => {
  if (unlocked.value === null) return 0
  return stepIds.filter((id) => id < unlocked.value).length
})

const isFullyCompleted = computed(() => unlocked.value !== null && unlocked.value > stepIds.length)

const uptimeLabel = computed(() => {
  const m = Math.floor(uptimeSeconds.value / 60).toString().padStart(2, '0')
  const s = (uptimeSeconds.value % 60).toString().padStart(2, '0')
  return `${m}:${s}`
})

const logLines = [
  'СКАНУВАННЯ ЗАВЕРШЕНО. АНОМАЛІЙ НЕ ВИЯВЛЕНО.',
  "З'ЄДНАННЯ З АРХІВОМ СТАБІЛЬНЕ.",
  'РЕЗЕРВНЕ ЖИВЛЕННЯ В НОРМІ.',
  'ОЧІКУЮ НАСТУПНОЇ КОМАНДИ.',
  'ПРОТОКОЛ ДОСТУПУ АКТИВНИЙ.',
  'ЦІЛІСНІСТЬ ДАНИХ ПІДТВЕРДЖЕНО.'
]

onMounted(async () => {
  const res = await fetch(`/api/progress?password=${encodeURIComponent(getPassword())}`)
  if (res.status === 401) {
    router.push({ name: 'login' })
    return
  }
  const data = await res.json()
  unlocked.value = data.unlocked
  loading.value = false

  if (!hasSeenOnboarding()) {
	showOnboarding.value = true
  }

  uptimeTimer = setInterval(() => {
    uptimeSeconds.value++
  }, 1000)

  gsap.from('.hub-panel', { opacity: 0, y: 16, duration: 0.5, stagger: 0.12, ease: 'power2.out' })
})

onUnmounted(() => {
  if (uptimeTimer) clearInterval(uptimeTimer)
})

function goToNextMission() {
  router.push(isFullyCompleted.value ? { name: 'finish' } : { name: 'quest-step', params: { stepId: unlocked.value } })
}
function goToArchive() {
  router.push({ name: 'archive' })
}
function goToTerminal() {
  router.push({ name: 'terminal' })
}
</script>

<template>
  <JarvisLayout eyebrow="ГОЛОВНИЙ ІНТЕРФЕЙС" wide>
    <div class="hub-head">
      <div class="hub-head-text">
        <h1 class="hub-title">Панель керування JARVIS</h1>
        <p class="hub-operator">ОПЕРАТОР: ВЛАД&nbsp;&nbsp;·&nbsp;&nbsp;РІВЕНЬ ДОСТУПУ: ALPHA</p>
      </div>

      <svg class="hub-radar" viewBox="0 0 100 100" aria-hidden="true">
        <circle class="radar-ring" cx="50" cy="50" r="46" />
        <circle class="radar-ring" cx="50" cy="50" r="32" />
        <circle class="radar-ring" cx="50" cy="50" r="18" />
        <line class="radar-cross" x1="4" y1="50" x2="96" y2="50" />
        <line class="radar-cross" x1="50" y1="4" x2="50" y2="96" />
        <g class="radar-sweep">
          <path d="M 50 50 L 50 4 A 46 46 0 0 1 90 26 Z" />
        </g>
        <circle class="radar-dot" cx="50" cy="50" r="3" />
      </svg>
    </div>

    <div class="hub-telemetry" v-if="!loading">
      <div class="telemetry-item">
        <span class="telemetry-label">СТАБІЛЬНІСТЬ РЕАКТОРА</span>
        <span class="telemetry-value">98.6%</span>
      </div>
      <div class="telemetry-item">
        <span class="telemetry-label">ЧАС СЕАНСУ</span>
        <span class="telemetry-value">{{ uptimeLabel }}</span>
      </div>
      <div class="telemetry-item">
        <span class="telemetry-label">АРХІВНИХ ЗАПИСІВ</span>
        <span class="telemetry-value">{{ completedCount }}</span>
      </div>
    </div>

    <div class="hub-grid" v-if="!loading">
      <button class="hub-panel hub-tile" data-onboard="mission" @click="goToNextMission">
        <svg class="hub-tile-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path d="M12 3 L20 10 L20 21 L14 21 L14 15 L10 15 L10 21 L4 21 L4 10 Z" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
        </svg>
        <span class="hub-tile-text">
          <span class="hub-tile-title">{{ isFullyCompleted ? 'Фінальний протокол' : 'Продовжити місію' }}</span>
          <span class="hub-tile-desc">{{ isFullyCompleted ? 'Всі сегменти розшифровано' : 'Перейти до наступного завдання' }}</span>
        </span>
      </button>

      <button class="hub-panel hub-tile" data-onboard="archive" @click="goToArchive" :disabled="completedCount === 0">
        <svg class="hub-tile-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <rect x="3.5" y="6" width="17" height="14" rx="1.5" stroke="currentColor" stroke-width="1.6" />
          <path d="M3.5 10 H20.5" stroke="currentColor" stroke-width="1.6" />
          <path d="M8 6 V4.5 H16 V6" stroke="currentColor" stroke-width="1.6" />
        </svg>
        <span class="hub-tile-text">
          <span class="hub-tile-title">Архів</span>
          <span class="hub-tile-desc">{{ completedCount === 0 ? 'Записів ще немає' : 'Переглянути зібрані досьє' }}</span>
        </span>
      </button>

      <button class="hub-panel hub-tile" data-onboard="terminal" @click="goToTerminal">
        <svg class="hub-tile-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <rect x="3" y="4" width="18" height="16" rx="1.5" stroke="currentColor" stroke-width="1.6" />
          <path d="M6.5 9 L10 12 L6.5 15" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M12.5 15 H17.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
        </svg>
        <span class="hub-tile-text">
          <span class="hub-tile-title">Термінал</span>
          <span class="hub-tile-desc">Прямий канал зв'язку з системою</span>
        </span>
      </button>
    </div>

    <div class="hub-panel hub-log" v-if="!loading" aria-hidden="true">
      <div class="hub-log-track">
        <span v-for="(line, i) in [...logLines, ...logLines]" :key="i" class="hub-log-item">&gt; {{ line }}</span>
      </div>
    </div>
  </JarvisLayout>

  <OnboardingOverlay
	v-if="showOnboarding"
	:steps="onboardingSteps"
	@finish="onOnboardingFinish"
/>
</template>

<style scoped>
.hub-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  width: 100%;
  margin-bottom: 18px;
}
.hub-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(20px, 5vw, 28px);
  margin: 0 0 6px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.hub-operator {
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  color: var(--jarvis-text-dim);
  letter-spacing: 1px;
  margin: 0;
}

.hub-radar {
  width: clamp(56px, 14vw, 84px);
  height: clamp(56px, 14vw, 84px);
  flex-shrink: 0;
}
.radar-ring {
  fill: none;
  stroke: rgba(71, 191, 255, 0.3);
  stroke-width: 1;
}
.radar-cross {
  stroke: rgba(71, 191, 255, 0.18);
  stroke-width: 1;
}
.radar-dot {
  fill: var(--jarvis-cyan);
}
.radar-sweep {
  transform-origin: 50px 50px;
  animation: radar-spin 3.2s linear infinite;
}
.radar-sweep path {
  fill: rgba(71, 191, 255, 0.22);
}
@keyframes radar-spin {
  to { transform: rotate(360deg); }
}

.hub-telemetry {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 1px;
  width: 100%;
  background: rgba(126, 20, 255, 0.25);
  border: 1px solid rgba(126, 20, 255, 0.25);
  margin-bottom: 24px;
}
.telemetry-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 10px 14px;
  background: rgba(13, 18, 32, 0.75);
}
.telemetry-label {
  font-family: var(--jarvis-mono);
  font-size: clamp(9px, 2.2vw, 10.5px);
  color: var(--jarvis-text-dim);
  letter-spacing: 1px;
}
.telemetry-value {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(15px, 3.6vw, 18px);
  color: var(--jarvis-cyan);
}

.hub-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  width: 100%;
  margin-bottom: 24px;
}
.hub-tile {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  padding: 20px;
  background: rgba(13, 18, 32, 0.6);
  border: 1px solid rgba(126, 20, 255, 0.35);
  border-radius: 6px;
  color: var(--jarvis-text);
  text-align: left;
  cursor: pointer;
  transition: border-color 0.2s, transform 0.2s, background 0.2s;
}
.hub-tile::before,
.hub-tile::after {
  content: '';
  position: absolute;
  width: 10px;
  height: 10px;
  border: 1px solid transparent;
  opacity: 0;
  transition: opacity 0.2s;
}
.hub-tile::before {
  top: 4px;
  left: 4px;
  border-top-color: var(--jarvis-cyan);
  border-left-color: var(--jarvis-cyan);
}
.hub-tile::after {
  bottom: 4px;
  right: 4px;
  border-bottom-color: var(--jarvis-cyan);
  border-right-color: var(--jarvis-cyan);
}
.hub-tile:hover:not(:disabled),
.hub-tile:focus-visible:not(:disabled) {
  border-color: var(--jarvis-cyan);
  background: rgba(13, 18, 32, 0.9);
  transform: translateY(-2px);
}
.hub-tile:hover:not(:disabled)::before,
.hub-tile:hover:not(:disabled)::after {
  opacity: 1;
}
.hub-tile:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.hub-tile-icon {
  width: 26px;
  height: 26px;
  color: var(--jarvis-cyan);
}
.hub-tile-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(14px, 3.4vw, 16px);
  letter-spacing: 1px;
}
.hub-tile-desc {
  font-size: clamp(11px, 2.6vw, 13px);
  color: var(--jarvis-text-dim);
}
.hub-tile-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.hub-log {
  width: 100%;
  overflow: hidden;
  border-top: 1px solid rgba(232, 230, 240, 0.12);
  border-bottom: 1px solid rgba(232, 230, 240, 0.12);
  padding: 10px 0;
  mask-image: linear-gradient(90deg, transparent, black 10%, black 90%, transparent);
}
.hub-log-track {
  display: flex;
  gap: 40px;
  white-space: nowrap;
  animation: log-scroll 22s linear infinite;
}
.hub-log-item {
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  color: var(--jarvis-text-dim);
  letter-spacing: 0.5px;
}
@keyframes log-scroll {
  from { transform: translateX(0); }
  to { transform: translateX(-50%); }
}

@media (prefers-reduced-motion: reduce) {
  .radar-sweep,
  .hub-log-track {
    animation: none;
  }
}

@media (max-width: 560px) {
  .hub-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 14px;
  }
  .hub-radar {
    width: 44px;
    height: 44px;
  }
  .hub-telemetry {
    grid-template-columns: 1fr;
    margin-bottom: 18px;
  }
  .telemetry-item {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
  }
  .hub-grid {
    gap: 10px;
    margin-bottom: 18px;
  }
  .hub-tile {
    flex-direction: row;
    align-items: center;
    gap: 14px;
    padding: 14px 16px;
  }
  .hub-tile-icon {
    width: 22px;
    height: 22px;
    flex-shrink: 0;
  }
  .hub-tile-text {
    gap: 2px;
  }
  .hub-log {
    margin-top: 6px;
    padding: 12px 0;
  }
}
</style>