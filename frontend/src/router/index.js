import { createRouter, createWebHashHistory } from 'vue-router'
import WelcomeView from '../views/WelcomeView.vue'
import LoginView from '../views/LoginView.vue'
import QuestStepLoader from '../views/QuestStepLoader.vue'
import FinishView from '../views/FinishView.vue'
import NotFoundView from '../views/NotFoundView.vue'
import { getPassword } from '../lib/auth'
import { stepIds } from '../data/steps'

const routes = [
  { path: '/', name: 'welcome', component: WelcomeView },
  { path: '/login', name: 'login', component: LoginView },
  { path: '/quest/:stepId(\\d+)', name: 'quest-step', component: QuestStepLoader, props: true },
  { path: '/finish', name: 'finish', component: FinishView },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundView }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

let cachedUnlocked = null

async function fetchUnlocked() {
  if (cachedUnlocked !== null) return cachedUnlocked
  const password = getPassword()
  if (!password) return null

  try {
    const res = await fetch(`/api/progress?password=${encodeURIComponent(password)}`)
    if (res.status === 401) return null
    const data = await res.json()
    cachedUnlocked = data.unlocked || 1
  } catch {
    return null
  }
  return cachedUnlocked
}

export function invalidateProgressCache() {
  cachedUnlocked = null
}

export async function resolveDestination() {
  const unlocked = await fetchUnlocked()
  if (unlocked === null) return { name: 'login' }
  if (unlocked > stepIds.length) return { name: 'finish' }
  return { name: 'quest-step', params: { stepId: unlocked } }
}

router.beforeEach(async (to) => {
  if (to.name === 'welcome' || to.name === 'not-found') return

  const unlocked = await fetchUnlocked()

  if (unlocked === null) {
    if (to.name !== 'login') return { name: 'login' }
    return
  }

  if (to.name === 'login') return resolveDestination()

  const completed = unlocked > stepIds.length
  if (completed && to.name !== 'finish') return { name: 'finish' }
  if (!completed && to.name === 'finish') return { name: 'quest-step', params: { stepId: unlocked } }
  if (to.name === 'quest-step') {
    const requested = Number(to.params.stepId)
    if (requested > unlocked) return { name: 'quest-step', params: { stepId: unlocked } }
  }
})

export default router