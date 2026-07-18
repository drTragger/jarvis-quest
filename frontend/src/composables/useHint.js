import { ref, computed } from 'vue'
import { getPassword } from '../lib/auth'

const SCRAMBLE_CHARS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ01#$%&'

function formatRemaining(seconds) {
  if (seconds >= 3600) {
    const h = Math.floor(seconds / 3600)
    const m = Math.floor((seconds % 3600) / 60)
    return m > 0 ? `${h} год ${m} хв` : `${h} год`
  }
  if (seconds >= 60) {
    return `${Math.floor(seconds / 60)} хв`
  }
  return `${seconds} с`
}

function scrambleReveal(target, text, duration = 900) {
  const start = performance.now()
  function frame(now) {
    const progress = Math.min(1, (now - start) / duration)
    const revealCount = Math.floor(progress * text.length)
    let out = ''
    for (let i = 0; i < text.length; i++) {
      if (i < revealCount || text[i] === ' ') out += text[i]
      else out += SCRAMBLE_CHARS[Math.floor(Math.random() * SCRAMBLE_CHARS.length)]
    }
    target.value = out
    if (progress < 1) requestAnimationFrame(frame)
    else target.value = text
  }
  requestAnimationFrame(frame)
}

export function useHint(stepId) {
  const scanning = ref(false)
  const revealed = ref(false)
  const hintText = ref('')
  const remaining = ref(null)
  const remainingLabel = computed(() => remaining.value === null ? null : formatRemaining(remaining.value))
  const cooldown = ref(false)

  async function checkHint() {
    if (scanning.value || cooldown.value || revealed.value) return
    scanning.value = true
    remaining.value = null

    try {
      const res = await fetch(
        `/api/hint?password=${encodeURIComponent(getPassword())}&step_id=step${stepId}`
      )
      const data = await res.json()

      if (data.eligible) {
        revealed.value = true
        scrambleReveal(hintText, data.hint)
      } else {
        remaining.value = data.seconds_remaining ?? null
        cooldown.value = true
        setTimeout(() => { cooldown.value = false }, 5000)
      }
    } finally {
      scanning.value = false
    }
  }

  return { scanning, revealed, hintText, remaining, remainingLabel, cooldown, checkHint }
}