<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  eyebrow: { type: String, default: '' },
  wide: { type: Boolean, default: false },
  showScan: { type: Boolean, default: true },
  alert: { type: Boolean, default: false },
  warn: { type: Boolean, default: false }
})

const flashing = ref(false)
let flashTimer = null

const warnActive = ref(false)
let warnTimer = null

watch(() => props.alert, (val) => {
  if (!val) return
  flashing.value = true
  clearTimeout(flashTimer)
  flashTimer = setTimeout(() => { flashing.value = false }, 600)
})

watch(() => props.warn, (val) => {
  if (!val) return
  warnActive.value = true
  clearTimeout(warnTimer)
  warnTimer = setTimeout(() => { warnActive.value = false }, 900)
})
</script>

<template>
  <section class="jarvis-frame" :class="{ 'is-alert': alert, 'is-flashing': flashing, 'is-warn': warnActive }">
  <div class="hud-grid" aria-hidden="true"></div>
  <div class="alert-vignette" aria-hidden="true"></div>
  <div v-if="showScan" class="scan-line" aria-hidden="true"></div>
  <div v-if="alert" class="alert-scan" aria-hidden="true"></div>
  <div v-if="warnActive" class="warn-scan" aria-hidden="true"></div>

  <span class="corner corner-tl" aria-hidden="true"></span>
  <span class="corner corner-tr" aria-hidden="true"></span>
  <span class="corner corner-bl" aria-hidden="true"></span>
  <span class="corner corner-br" aria-hidden="true"></span>

  <div class="frame-content" :class="{ 'frame-content-wide': wide }">
    <p v-if="eyebrow" class="frame-eyebrow">{{ eyebrow }}</p>
    <slot name="header" />
    <slot />
  </div>
</section>
</template>

<style scoped>
.jarvis-frame {
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  width: 100%;
  background: radial-gradient(ellipse at 50% 30%, var(--jarvis-bg-soft) 0%, var(--jarvis-bg) 70%);
  color: var(--jarvis-text);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--jarvis-mono);
  padding: 16px;
  box-sizing: border-box;
}

.hud-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(126, 20, 255, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(126, 20, 255, 0.08) 1px, transparent 1px);
  background-size: 32px 32px;
  mask-image: radial-gradient(ellipse at 50% 40%, black 20%, transparent 75%);
}
@media (min-width: 640px) {
  .hud-grid { background-size: 48px 48px; }
}

.scan-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--jarvis-cyan), transparent);
  box-shadow: 0 0 12px 2px rgba(71, 191, 255, 0.6);
  animation: scan 5s linear infinite;
}
@keyframes scan {
  0% { top: -2%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 102%; opacity: 0; }
}

.corner {
  position: absolute;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(126, 20, 255, 0.5);
}
.corner-tl { top: 12px; left: 12px; border-right: none; border-bottom: none; }
.corner-tr { top: 12px; right: 12px; border-left: none; border-bottom: none; }
.corner-bl { bottom: 12px; left: 12px; border-right: none; border-top: none; }
.corner-br { bottom: 12px; right: 12px; border-left: none; border-top: none; }
@media (min-width: 640px) {
  .corner { width: 28px; height: 28px; }
  .corner-tl, .corner-tr { top: 20px; }
  .corner-bl, .corner-br { bottom: 20px; }
  .corner-tl, .corner-bl { left: 20px; }
  .corner-tr, .corner-br { right: 20px; }
}

.frame-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  width: 100%;
  max-width: 480px;
  padding: 12px;
}
.frame-content-wide {
  max-width: 720px;
}

.frame-eyebrow {
  font-size: clamp(10px, 2.6vw, 12px);
  letter-spacing: 3px;
  color: var(--jarvis-cyan);
  margin: 0 0 8px;
}

@media (prefers-reduced-motion: reduce) {
  .scan-line { animation: none; }
}

.alert-vignette {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: radial-gradient(ellipse at 50% 50%, transparent 30%, rgba(255, 40, 70, 0.35) 100%);
  opacity: 0;
  transition: opacity 0.25s;
}
.jarvis-frame.is-alert .alert-vignette {
  opacity: 1;
  animation: alert-pulse 1.3s ease-in-out infinite;
}
@keyframes alert-pulse {
  0%, 100% { opacity: 0.45; }
  50% { opacity: 1; }
}

.alert-scan {
  position: absolute;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, transparent, #ff3355, transparent);
  box-shadow: 0 0 16px 3px rgba(255, 51, 85, 0.8);
  animation: scan 1.1s linear infinite;
}

.jarvis-frame.is-alert .scan-line {
  animation: none;
  opacity: 0;
}

.jarvis-frame.is-alert .corner {
  border-color: rgba(255, 51, 85, 0.85);
  transition: border-color 0.2s;
}

.jarvis-frame.is-alert .hud-grid {
  background-image:
    linear-gradient(rgba(255, 51, 85, 0.12) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 51, 85, 0.12) 1px, transparent 1px);
}

.jarvis-frame.is-flashing {
  animation: alert-shake 0.5s ease-in-out;
}
@keyframes alert-shake {
  0%, 100% { transform: translateX(0); }
  20% { transform: translateX(-6px); }
  40% { transform: translateX(5px); }
  60% { transform: translateX(-4px); }
  80% { transform: translateX(3px); }
}

@media (prefers-reduced-motion: reduce) {
  .jarvis-frame.is-flashing { animation: none; }
  .jarvis-frame.is-alert .alert-vignette { animation: none; opacity: 0.6; }
  .alert-scan { animation: none; }
}

.warn-scan {
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #ffb648, transparent);
  box-shadow: 0 0 12px 2px rgba(255, 182, 72, 0.6);
  animation: scan 0.9s linear 1;
}

.jarvis-frame.is-warn .scan-line {
  animation: none;
  opacity: 0;
}

.jarvis-frame.is-warn .corner {
  border-color: rgba(255, 182, 72, 0.7);
  transition: border-color 0.3s;
}

.jarvis-frame.is-warn .alert-vignette {
  opacity: 0.22;
  background: radial-gradient(ellipse at 50% 50%, transparent 40%, rgba(255, 182, 72, 0.25) 100%);
  animation: warn-fade 0.9s ease-out 1;
}
@keyframes warn-fade {
  0% { opacity: 0.35; }
  100% { opacity: 0; }
}

@media (prefers-reduced-motion: reduce) {
  .warn-scan { animation: none; }
  .jarvis-frame.is-warn .alert-vignette { animation: none; opacity: 0; }
}
</style>