<script setup>
import { ref, nextTick, onMounted } from 'vue'

const emit = defineEmits(['solved', 'failed'])
const CAPTCHA_CHARS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'

const canvasEl = ref(null)
const input = ref('')
const refreshCooldown = ref(false)
let captchaText = ''

function randomCaptchaText(length = 6) {
  let out = ''
  for (let i = 0; i < length; i++) {
    out += CAPTCHA_CHARS[Math.floor(Math.random() * CAPTCHA_CHARS.length)]
  }
  return out
}

function drawCaptcha() {
  const canvas = canvasEl.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  const w = canvas.width
  const h = canvas.height

  ctx.clearRect(0, 0, w, h)
  ctx.fillStyle = 'rgba(13, 18, 32, 0.92)'
  ctx.fillRect(0, 0, w, h)

  for (let i = 0; i < 5; i++) {
    ctx.strokeStyle = `rgba(126, 20, 255, ${0.2 + Math.random() * 0.3})`
    ctx.lineWidth = 1 + Math.random() * 1.5
    ctx.beginPath()
    ctx.moveTo(Math.random() * w, Math.random() * h)
    ctx.bezierCurveTo(Math.random() * w, Math.random() * h, Math.random() * w, Math.random() * h, Math.random() * w, Math.random() * h)
    ctx.stroke()
  }

  const charWidth = w / (captchaText.length + 1)
  ctx.textBaseline = 'middle'
  ctx.textAlign = 'center'
  for (let i = 0; i < captchaText.length; i++) {
    const ch = captchaText[i]
    const x = charWidth * (i + 0.85)
    const y = h / 2 + (Math.random() - 0.5) * 16
    const angle = (Math.random() - 0.5) * 0.7
    const size = 22 + Math.random() * 10
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(angle)
    ctx.font = `${size}px 'Share Tech Mono', monospace`
    ctx.fillStyle = Math.random() > 0.5 ? '#47bfff' : '#e8e6f0'
    ctx.fillText(ch, 0, 0)
    ctx.restore()
  }

  for (let i = 0; i < 45; i++) {
    ctx.fillStyle = `rgba(255, 255, 255, ${Math.random() * 0.18})`
    ctx.fillRect(Math.random() * w, Math.random() * h, 1.4, 1.4)
  }

  for (let i = 0; i < 3; i++) {
    ctx.strokeStyle = `rgba(255, 255, 255, ${0.08 + Math.random() * 0.14})`
    ctx.lineWidth = 1
    ctx.beginPath()
    ctx.moveTo(0, Math.random() * h)
    ctx.lineTo(w, Math.random() * h)
    ctx.stroke()
  }
}

function newCaptcha() {
  captchaText = randomCaptchaText()
  input.value = ''
  nextTick(drawCaptcha)
}

function submit() {
  if (!input.value.trim()) return
  if (input.value.trim() === captchaText) {
    emit('solved')
  } else {
    emit('failed')
    newCaptcha()
  }
}

function refresh() {
  if (refreshCooldown.value) return
  refreshCooldown.value = true
  newCaptcha()
  setTimeout(() => { refreshCooldown.value = false }, 1500)
}

onMounted(newCaptcha)
</script>

<template>
  <div class="captcha-glyphs">
    <p class="captcha-label">РОЗПІЗНАЙ СИМВОЛИ</p>
    <div class="glyphs-panel">
      <canvas ref="canvasEl" class="glyphs-canvas" width="320" height="90"></canvas>
      <button class="glyphs-refresh" :disabled="refreshCooldown" title="Згенерувати нову послідовність" @click="refresh">&#8635;</button>
    </div>
    <input
      v-model="input"
      class="jarvis-input"
      type="text"
      placeholder="ВВЕДІТЬ СИМВОЛИ ІЗ ЗОБРАЖЕННЯ"
      autocomplete="off"
      autocapitalize="off"
      autocorrect="off"
      spellcheck="false"
      @keyup.enter="submit"
    />
    <button class="jarvis-btn" @click="submit"><span>ПЕРЕВІРИТИ</span></button>
  </div>
</template>

<style scoped>
.captcha-glyphs {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.captcha-label {
  font-family: var(--jarvis-mono);
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--jarvis-text-dim);
  margin: 0;
}
.glyphs-panel {
  position: relative;
  width: 100%;
  max-width: 320px;
}
.glyphs-canvas {
  width: 100%;
  height: 90px;
  border: 1px solid rgba(126, 20, 255, 0.4);
  border-radius: 6px;
  display: block;
}
.glyphs-refresh {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 26px;
  height: 26px;
  background: rgba(5, 7, 10, 0.7);
  border: 1px solid rgba(71, 191, 255, 0.5);
  border-radius: 4px;
  color: var(--jarvis-cyan);
  font-size: 14px;
  line-height: 1;
  cursor: pointer;
}
.glyphs-refresh:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
