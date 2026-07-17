<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['solved', 'failed'])

const SUPERSCRIPTS = { 2: '²', 3: '³', 4: '⁴' }

const a = ref(3)
const b = ref(0)
const c = ref(0)
const n = ref(2)
const answer = ref('')

function randomNonZeroMultiple(step, maxMultiplier) {
  const magnitude = step * (1 + Math.floor(Math.random() * maxMultiplier))
  return Math.random() > 0.5 ? magnitude : -magnitude
}

function generate() {
  a.value = randomNonZeroMultiple(3, 4)
  b.value = Math.random() > 0.15 ? randomNonZeroMultiple(2, 4) : 0
  c.value = Math.floor(Math.random() * 13) - 6
  n.value = 2 + Math.floor(Math.random() * 3)
  answer.value = ''
}
generate()

const expression = computed(() => {
  const terms = []
  if (a.value !== 0) terms.push({ coef: a.value, deg: 2 })
  if (b.value !== 0) terms.push({ coef: b.value, deg: 1 })
  if (c.value !== 0) terms.push({ coef: c.value, deg: 0 })

  return terms
    .map((term, i) => {
      const abs = Math.abs(term.coef)
      const coefStr = abs === 1 && term.deg !== 0 ? '' : String(abs)
      const varStr = term.deg === 2 ? 'x²' : term.deg === 1 ? 'x' : ''
      const body = `${coefStr}${varStr}`
      if (i === 0) return term.coef < 0 ? `-${body}` : body
      return term.coef < 0 ? ` - ${body}` : ` + ${body}`
    })
    .join('')
})

const upperBound = computed(() => SUPERSCRIPTS[n.value] || n.value)

function targetValue() {
  return (a.value / 3) * n.value ** 3 + (b.value / 2) * n.value ** 2 + c.value * n.value
}

function submit() {
  if (answer.value.trim() === '') return
  if (Number(answer.value.trim()) === targetValue()) {
    emit('solved')
  } else {
    emit('failed')
    generate()
  }
}
</script>

<template>
  <div class="captcha-math">
    <p class="captcha-label">ОБЧИСЛИ ВИЗНАЧЕНИЙ ІНТЕГРАЛ</p>
    <p class="captcha-expr">
      <span class="integral-sign">&int;<sub>0</sub><sup>{{ upperBound }}</sup></span>
      ({{ expression }}) dx = ?
    </p>
    <input
      v-model="answer"
      class="jarvis-input"
      type="text"
      inputmode="numeric"
      pattern="-?[0-9]*"
      placeholder="ВВЕДИ ЧИСЛОВЕ ЗНАЧЕННЯ"
      autocomplete="off"
      @keyup.enter="submit"
    />
    <button class="jarvis-btn" @click="submit"><span>ПЕРЕВІРИТИ</span></button>
  </div>
</template>

<style scoped>
.captcha-math {
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
.captcha-expr {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(20px, 5.4vw, 27px);
  color: var(--jarvis-text);
  margin: 4px 0 6px;
  letter-spacing: 1px;
  flex-wrap: wrap;
  text-align: center;
}
.integral-sign {
  font-size: clamp(30px, 8vw, 40px);
  line-height: 1;
  color: var(--jarvis-cyan);
  text-shadow: 0 0 14px rgba(71, 191, 255, 0.5);
}
.integral-sign sub {
  font-size: 0.4em;
}
.integral-sign sup {
  font-size: 0.4em;
}
</style>
