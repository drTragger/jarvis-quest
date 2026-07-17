<script setup>
</script>

<template>
  <div class="record-signal">
    <p class="record-id">ПРОТОКОЛ ПЕРЕДАЧІ СИГНАЛУ &middot; КАНАЛ СТАБІЛЬНИЙ: 100%</p>
    <h2 class="record-title">Послідовність розшифровано</h2>
    <p class="record-text">
      ASCII-байти перехопленої передачі зіставлено з символами. Протокол
      підтвердив цілісність даних і відкрив доступ до наступного сегмента.
    </p>
    <p class="record-hex">30 76 65 72 72 69 64 65</p>
    <p class="record-word">0verride</p>
  </div>
</template>

<style scoped>
.record-signal {
  width: 100%;
  padding: clamp(18px, 4vw, 26px);
  background: rgba(13, 18, 32, 0.7);
  border: 1px solid rgba(126, 20, 255, 0.4);
  border-radius: 6px;
}
.record-id {
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  letter-spacing: 1px;
  color: var(--jarvis-text-dim);
  margin: 0 0 12px;
}
.record-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(16px, 4vw, 19px);
  margin: 0 0 12px;
  color: var(--jarvis-text);
}
.record-text {
  font-size: clamp(12px, 2.8vw, 14px);
  line-height: 1.6;
  color: var(--jarvis-text-dim);
  margin: 0 0 16px;
}
.record-hex {
  font-family: var(--jarvis-mono);
  font-size: clamp(12px, 2.8vw, 14px);
  letter-spacing: 3px;
  color: var(--jarvis-text-dim);
  margin: 0 0 10px;
}
.record-word {
  font-family: var(--jarvis-mono);
  font-weight: 700;
  font-size: clamp(20px, 5vw, 26px);
  letter-spacing: 4px;
  color: var(--jarvis-cyan);
  text-shadow: 0 0 14px rgba(126, 20, 255, 0.6);
  margin: 0;
}
</style>
