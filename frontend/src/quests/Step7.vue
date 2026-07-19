<script setup>
import { computed, onMounted } from 'vue'
import gsap from 'gsap'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { useGeoLocator } from '../composables/useGeoLocator'
import { playSound } from '../composables/useSound'
import step7Voice from '../assets/audio/step7-voice.mp3'

const TARGET_LAT = 48.9157593
const TARGET_LON = 24.6971441
const LOCK_RADIUS_M = 40
const MAX_RANGE_M = 5000
const FAR_THRESHOLD_M = 15000

const {
  distance,
  bearingToTarget,
  heading,
  headingAvailable,
  locked,
  geoError,
  compassPermissionNeeded,
  requestCompass
} = useGeoLocator(TARGET_LAT, TARGET_LON, LOCK_RADIUS_M)

const tooFar = computed(() => distance.value !== null && distance.value > FAR_THRESHOLD_M)

const needleRotation = computed(() => {
  if (bearingToTarget.value === null) return 0
  if (headingAvailable.value && heading.value !== null) {
    return (bearingToTarget.value - heading.value + 360) % 360
  }
  return bearingToTarget.value
})

const signalStrength = computed(() => {
  if (distance.value === null) return 0
  if (distance.value <= LOCK_RADIUS_M) return 100
  if (distance.value >= MAX_RANGE_M) return 0
  const t = 1 - (distance.value - LOCK_RADIUS_M) / (MAX_RANGE_M - LOCK_RADIUS_M)
  return Math.round(Math.pow(t, 0.6) * 100)
})

const signalBars = computed(() => Math.round(signalStrength.value / 20))

const distanceLabel = computed(() => {
  if (distance.value === null) return '—'
  if (distance.value < 1000) return `${Math.round(distance.value)} м`
  return `${(distance.value / 1000).toFixed(1)} км`
})

const COMPASS_LABELS = ['ПН', 'ПН-СХ', 'СХ', 'ПД-СХ', 'ПД', 'ПД-ЗХ', 'ЗХ', 'ПН-ЗХ']
const bearingLabel = computed(() => {
  if (bearingToTarget.value === null) return '—'
  const idx = Math.round(bearingToTarget.value / 45) % 8
  return `${COMPASS_LABELS[idx]}, ${Math.round(bearingToTarget.value)}°`
})

onMounted(() => {
  const tl = gsap.timeline({ defaults: { ease: 'power2.out' } })
  tl.from('.locator-panel', { opacity: 0, y: 20, duration: 0.6 })
    .eventCallback('onComplete', () => {
      playSound('step7-voice', step7Voice, { volume: 0.9 })
    })
})
</script>

<template>
  <JarvisLayout eyebrow="ЛОКАТОР СИГНАЛУ" :show-scan="false">
    <div class="locator-panel">
      <p class="locator-id">ДЖЕРЕЛО СИГНАЛУ &middot; ТОЧНІ КООРДИНАТИ ЗАСЕКРЕЧЕНО &middot; ТРІАНГУЛЯЦІЯ АКТИВНА</p>
      <h1 class="locator-title">Знайди джерело сигналу</h1>
      <p class="locator-text">
        Архів впіймав слабкий, стабільний сигнал, але не може вирахувати
        координати — тільки напрямок і силу. Рухайся туди, куди вказує
        стрілка. Що ближче ти підходиш, то сильніший сигнал.
      </p>

      <template v-if="geoError">
        <div class="error-card">
          <p>ГЕОЛОКАЦІЯ НЕДОСТУПНА</p>
          <p class="error-detail">{{ geoError }}</p>
          <p class="error-hint">Дозволь доступ до геолокації в браузері й онови сторінку.</p>
        </div>
      </template>

      <template v-else-if="tooFar">
        <div class="error-card">
          <p>СИГНАЛ ЗАНАДТО СЛАБКИЙ</p>
          <p class="error-detail">Джерело сигналу знаходиться в місті Івано-Франківськ.</p>
          <p class="error-hint">Наблизся до міста, щоб тріангуляція запрацювала.</p>
        </div>
      </template>

      <template v-else-if="!locked">
        <div class="compass-wrap">
          <div class="compass-dial">
            <span class="compass-tick tick-n">ПН</span>
            <span class="compass-tick tick-e">СХ</span>
            <span class="compass-tick tick-s">ПД</span>
            <span class="compass-tick tick-w">ЗХ</span>
            <div class="compass-needle" :style="{ transform: `rotate(${needleRotation}deg)` }">
              <span class="needle-tip"></span>
            </div>
          </div>
        </div>

        <p class="bearing-label">{{ bearingLabel }}</p>

        <div class="signal-meter">
          <div class="signal-bars">
            <span v-for="n in 5" :key="n" class="signal-bar" :class="{ 'is-filled': n <= signalBars }"></span>
          </div>
          <p class="signal-pct">{{ signalStrength }}%</p>
        </div>
        <p class="distance-label">Відстань до джерела: {{ distanceLabel }}</p>

        <button v-if="compassPermissionNeeded" class="jarvis-btn" type="button" @click="requestCompass">
          УВІМКНУТИ КОМПАС
        </button>
      </template>

      <template v-else>
        <div class="lock-card">
          <p class="lock-title">СИГНАЛ ЗАФІКСОВАНО</p>
          <p class="lock-text">Джерело сигналу поряд. Знайди його.</p>
          <p class="lock-text">Шукай щось своє. Пароль "navigation"</p>
        </div>
      </template>
    </div>
  </JarvisLayout>
</template>

<style scoped>
.locator-panel {
  width: 100%;
  max-width: 480px;
}
.locator-id {
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  color: var(--jarvis-text-dim);
  margin: 0 0 14px;
}
.locator-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(19px, 4.8vw, 24px);
  margin: 0 0 14px;
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.5);
}
.locator-text {
  font-size: clamp(13px, 3vw, 15px);
  line-height: 1.6;
  color: var(--jarvis-text-dim);
  margin: 0 0 26px;
}

.compass-wrap {
  display: flex;
  justify-content: center;
  margin-bottom: 12px;
}
.compass-dial {
  position: relative;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  border: 2px solid rgba(71, 191, 255, 0.4);
  background: radial-gradient(circle, rgba(71, 191, 255, 0.05), transparent 70%);
}
.compass-tick {
  position: absolute;
  font-family: var(--jarvis-mono);
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--jarvis-text-dim);
}
.tick-n { top: 6px; left: 50%; transform: translateX(-50%); }
.tick-s { bottom: 6px; left: 50%; transform: translateX(-50%); }
.tick-e { right: 8px; top: 50%; transform: translateY(-50%); }
.tick-w { left: 8px; top: 50%; transform: translateY(-50%); }

.compass-needle {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 4px;
  height: 78px;
  margin-left: -2px;
  margin-top: -78px;
  transform-origin: 50% 100%;
  background: linear-gradient(180deg, #47bfff, rgba(71, 191, 255, 0.15));
  transition: transform 0.2s ease-out;
}
.needle-tip {
  position: absolute;
  top: -6px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 7px solid transparent;
  border-right: 7px solid transparent;
  border-bottom: 12px solid #47bfff;
}

.bearing-label {
  text-align: center;
  font-family: var(--jarvis-mono);
  font-size: 14px;
  letter-spacing: 2px;
  color: var(--jarvis-cyan);
  margin: 0 0 22px;
}

.signal-meter {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 8px;
}
.signal-bars {
  display: flex;
  align-items: flex-end;
  gap: 4px;
}
.signal-bar {
  width: 8px;
  background: rgba(71, 191, 255, 0.15);
  border: 1px solid rgba(71, 191, 255, 0.3);
}
.signal-bar:nth-child(1) { height: 8px; }
.signal-bar:nth-child(2) { height: 14px; }
.signal-bar:nth-child(3) { height: 20px; }
.signal-bar:nth-child(4) { height: 26px; }
.signal-bar:nth-child(5) { height: 32px; }
.signal-bar.is-filled {
  background: #47bfff;
  box-shadow: 0 0 8px rgba(71, 191, 255, 0.7);
}
.signal-pct {
  font-family: var(--jarvis-mono);
  font-weight: 700;
  font-size: 16px;
  color: var(--jarvis-cyan);
  margin: 0;
}

.distance-label {
  text-align: center;
  font-family: var(--jarvis-mono);
  font-size: 12px;
  color: var(--jarvis-text-dim);
  margin: 0 0 22px;
}

.error-card {
  text-align: center;
  padding: 20px;
  border: 1px dashed rgba(255, 71, 71, 0.5);
  color: #ff8080;
  font-family: var(--jarvis-mono);
  font-size: 13px;
  line-height: 1.6;
}
.error-detail {
  color: var(--jarvis-text-dim);
  font-size: 12px;
}
.error-hint {
  color: var(--jarvis-text-dim);
  font-size: 12px;
}

.lock-card {
  text-align: center;
  padding: 22px;
  border: 1px solid rgba(71, 191, 255, 0.5);
  border-radius: 6px;
  background: rgba(71, 191, 255, 0.06);
}
.lock-title {
  font-family: var(--jarvis-mono);
  font-weight: 700;
  letter-spacing: 3px;
  color: var(--jarvis-cyan);
  text-shadow: 0 0 14px rgba(71, 191, 255, 0.7);
  margin: 0 0 12px;
}
.lock-text {
  font-size: 13px;
  color: var(--jarvis-text-dim);
  margin: 0 0 14px;
}
.lock-text:last-child {
  margin-bottom: 0;
  font-weight: 700;
  color: var(--jarvis-cyan);
}

@media (prefers-reduced-motion: reduce) {
  .compass-needle { transition: none; }
}
</style>
