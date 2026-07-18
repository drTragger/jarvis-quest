import { ref } from 'vue'

export const MORSE_MAP = {
  'А': '.-', 'Б': '-...', 'В': '.--', 'Г': '--.', 'Ґ': '--.-', 'Д': '-..', 'Е': '.',
  'Є': '..-..', 'Ж': '...-', 'З': '--..', 'И': '-.--', 'І': '..', 'Ї': '.---',
  'Й': '.---', 'К': '-.-', 'Л': '.-..', 'М': '--', 'Н': '-.', 'О': '---', 'П': '.--.',
  'Р': '.-.', 'С': '...', 'Т': '-', 'У': '..-', 'Ф': '..-.', 'Х': '....', 'Ц': '-.-.',
  'Ч': '---.', 'Ш': '----', 'Щ': '--.-', 'Ь': '-..-', 'Ю': '..--', 'Я': '.-.-'
}

const MORSE_TO_LETTER = Object.fromEntries(
  Object.entries(MORSE_MAP).map(([letter, code]) => [code, letter])
)

const DOT_MAX_MS = 280
const LETTER_GAP_MS = 900

export function useMorseInput() {
  const word = ref('')
  const currentSymbol = ref('')
  const pressing = ref(false)
  let pressStart = 0
  let letterGapTimer = null

  function resolveLetter() {
    if (!currentSymbol.value) return
    word.value += MORSE_TO_LETTER[currentSymbol.value] || '?'
    currentSymbol.value = ''
  }

  function armLetterGap() {
    clearTimeout(letterGapTimer)
    letterGapTimer = setTimeout(resolveLetter, LETTER_GAP_MS)
  }

  function keyDown() {
    if (pressing.value) return
    pressing.value = true
    clearTimeout(letterGapTimer)
    pressStart = performance.now()
  }

  function keyUp() {
    if (!pressing.value) return
    pressing.value = false
    const duration = performance.now() - pressStart
    currentSymbol.value += duration < DOT_MAX_MS ? '.' : '-'
    armLetterGap()
  }

  function clear() {
    clearTimeout(letterGapTimer)
    word.value = ''
    currentSymbol.value = ''
    pressing.value = false
  }

  return { word, currentSymbol, pressing, keyDown, keyUp, clear }
}
