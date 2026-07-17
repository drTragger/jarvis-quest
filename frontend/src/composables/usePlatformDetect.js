import { ref, onMounted, onUnmounted } from 'vue'

function detectPlatform() {
  const ua = navigator.userAgent || ''
  const uaMobile = /Android|iPhone|iPad|iPod|Mobile|Windows Phone/i.test(ua)
  const coarsePointer = window.matchMedia('(pointer: coarse)').matches
  const noHover = window.matchMedia('(hover: none)').matches
  const hasTouch = navigator.maxTouchPoints > 0
  const narrowViewport = window.innerWidth <= 900

  const mobileVotes = [uaMobile, coarsePointer, noHover, hasTouch, narrowViewport].filter(Boolean).length
  return mobileVotes >= 3 ? 'mobile' : 'desktop'
}

export function usePlatformDetect() {
  const platform = ref(detectPlatform())
  const update = () => { platform.value = detectPlatform() }

  let pointerMq
  onMounted(() => {
    window.addEventListener('resize', update)
    pointerMq = window.matchMedia('(pointer: coarse)')
    pointerMq.addEventListener('change', update)
  })
  onUnmounted(() => {
    window.removeEventListener('resize', update)
    pointerMq?.removeEventListener('change', update)
  })

  return { platform }
}