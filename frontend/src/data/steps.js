export const stepComponents = {
  1: () => import('../quests/Step1.vue'),
  2: () => import('../quests/Step2.vue')
}

export const stepIds = Object.keys(stepComponents).map(Number).sort((a, b) => a - b)

export const collectibles = {
  1: {
    date: 'ЗАПИС №1',
    title: 'Досьє: Протокол J.A.R.V.I.S.',
    excerpt: 'Акронім розшифровано. Система підтвердила особу оператора та відкрила доступ до наступного архівного сегмента.'
  },
  2: {
    date: 'ЗАПИС №2',
    title: 'Перехоплений сигнал',
    excerpt: 'Приховане повідомлення в частотному спектрі розкодовано. Ще один рівень захисту архіву знято.'
  }
}

export const archiveRecords = {
  1: () => import('../quests/archive/Step1Record.vue'),
  2: () => import('../quests/archive/Step2Record.vue')
}