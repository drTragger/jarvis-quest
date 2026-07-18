import { createRouter, createWebHistory } from 'vue-router'
import WelcomeView from '../views/WelcomeView.vue'
import LoginView from '../views/LoginView.vue'
import RulesView from '../views/RulesView.vue'
import HubView from '../views/HubView.vue'
import ArchiveView from '../views/ArchiveView.vue'
import TerminalView from '../views/TerminalView.vue'
import QuestStepLoader from '../views/QuestStepLoader.vue'
import FinishView from '../views/FinishView.vue'
import NotFoundView from '../views/NotFoundView.vue'
import { getPassword } from '../lib/auth'
import { stepIds } from '../data/steps'
import { stopAllSounds } from '../composables/useSound'

const routes = [
  { path: '/', name: 'welcome', component: WelcomeView },
  { path: '/login', name: 'login', component: LoginView },
  { path: '/rules', name: 'rules', component: RulesView },
  { path: '/hub', name: 'hub', component: HubView },
  { path: '/archive', name: 'archive', component: ArchiveView },
  { path: '/terminal', name: 'terminal', component: TerminalView },
  { path: '/quest/:stepId(\\d+)', name: 'quest-step', component: QuestStepLoader, props: true },
  { path: '/finish', name: 'finish', component: FinishView },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

let cachedProgress = null

async function fetchProgress() {
  if (cachedProgress !== null) return cachedProgress
  const password = getPassword()
  if (!password) return null

  try {
    const res = await fetch(`/api/progress?password=${encodeURIComponent(password)}`)
    if (res.status === 401) return null
    const data = await res.json()
    cachedProgress = { unlocked: data.unlocked || 1, rulesSeen: !!data.rules_seen }
  } catch {
    return null
  }
  return cachedProgress
}

export function invalidateProgressCache() {
  cachedProgress = null
}

export async function resolveDestination() {
  const progress = await fetchProgress()
  if (progress === null) return { name: 'login' }
  if (!progress.rulesSeen) return { name: 'rules' }
  if (progress.unlocked > stepIds.length) return { name: 'finish' }
  return { name: 'hub' }
}

const POST_COMPLETION_ALLOWED = ['finish', 'hub', 'archive', 'terminal', 'rules']

router.beforeEach(async (to) => {
  if (to.name === 'welcome' || to.name === 'not-found') return

  const progress = await fetchProgress()

  if (progress === null) {
    if (to.name !== 'login') return { name: 'login' }
    return
  }

  if (to.name === 'login') return resolveDestination()

  if (!progress.rulesSeen && to.name !== 'rules') return { name: 'rules' }

  const completed = progress.unlocked > stepIds.length
  if (completed && !POST_COMPLETION_ALLOWED.includes(to.name)) return { name: 'finish' }
  if (!completed && to.name === 'finish') return { name: 'hub' }
  if (to.name === 'quest-step') {
    const requested = Number(to.params.stepId)
    if (requested !== progress.unlocked) return { name: 'quest-step', params: { stepId: progress.unlocked } }
  }
})

router.afterEach(() => {
  stopAllSounds()
})

export default router