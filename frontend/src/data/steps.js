export const stepComponents = {
  1: () => import('../quests/Step1.vue'),
  2: () => import('../quests/Step2.vue'),
  3: () => import('../quests/Step3.vue'),
  4: () => import('../quests/Step4.vue')
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
  },
  3: {
    date: 'ЗАПИС №3',
    title: 'Протокол передачі сигналу',
    excerpt: 'ASCII-послідовність перехопленої передачі декодовано. Черговий рівень захисту архіву знято.'
  },
  4: {
    date: 'ЗАПИС №4',
    title: 'Протокол перевірки людяності',
    excerpt: 'П\'ять спотворених перевірок розпізнано вручну. Система визнала оператора тим, чим і підозрювала.'
  },
}

export const archiveRecords = {
  1: () => import('../quests/archive/Step1Record.vue'),
  2: () => import('../quests/archive/Step2Record.vue'),
  3: () => import('../quests/archive/Step3Record.vue'),
  4: () => import('../quests/archive/Step4Record.vue')
}