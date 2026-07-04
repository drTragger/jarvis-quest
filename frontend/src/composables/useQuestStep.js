import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getPassword } from '../lib/auth'
import { invalidateProgressCache } from '../router'
import { stepIds } from '../data/steps'

export function useQuestStep(stepId) {
  const router = useRouter()
  const loading = ref(false)
  const error = ref(false)
  const success = ref(false)
  let pendingDestination = null

  fetch(`/api/step/start?password=${encodeURIComponent(getPassword())}&step_id=step${stepId}`).catch(() => {})

  async function submitAnswer(answer) {
    if (loading.value) return false
    loading.value = true
    error.value = false

    const res = await fetch('/api/answer', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        password: getPassword(),
        step_id: `step${stepId}`,
        answer
      })
    })

    if (res.status === 401) {
      loading.value = false
      router.push({ name: 'login' })
      return false
    }

    const data = await res.json()
    loading.value = false

    if (!data.correct) {
      error.value = true
      return false
    }

    invalidateProgressCache()
    const hasNext = stepIds.includes(data.unlocked)
    pendingDestination = hasNext ? { name: 'hub' } : { name: 'finish' }
    success.value = true
    return true
  }

  function proceedAfterSuccess() {
    if (pendingDestination) router.push(pendingDestination)
  }

  return { loading, error, success, submitAnswer, proceedAfterSuccess }
}