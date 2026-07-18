<script setup>
</script>

<template>
  <div class="record-signal">
    <p class="record-id">СКАНУВАННЯ ЗОБРАЖЕННЯ &middot; ПРИХОВАНИЙ ТЕКСТ ПРОЯВЛЕНО</p>
    <h2 class="record-title">Слід, прихований у тінях</h2>
    <p class="record-text">
      Підняття яскравості найтемнішої ділянки зображення виявило сирий рядок "ivoenavhz".
      Другий прохід тим самим зсувом Цезаря розкрив справжнє слово.
    </p>
    <p class="record-word">VIBRANIUM</p>
  </div>
</template>

<style scoped>
.record-signal {
  width: 100%;
  padding: clamp(18px, 4vw, 26px);
  background: rgba(13, 18, 32, 0.7);
  border: 1px solid rgba(71, 191, 255, 0.4);
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
.record-word {
  font-family: var(--jarvis-mono);
  font-weight: 700;
  font-size: clamp(20px, 5vw, 26px);
  letter-spacing: 4px;
  color: var(--jarvis-cyan);
  text-shadow: 0 0 14px rgba(71, 191, 255, 0.6);
  margin: 0;
}
</style>
