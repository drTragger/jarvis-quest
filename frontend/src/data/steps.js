export const stepComponents = {
  1: () => import('../quests/Step1.vue')
}

export const stepIds = Object.keys(stepComponents).map(Number).sort((a, b) => a - b)