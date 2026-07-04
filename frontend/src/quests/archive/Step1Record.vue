<script setup>
</script>

<template>
  <div class="record-page">
    <div class="dossier-card">
      <div class="stamp" aria-hidden="true">
        <span>РОЗСЕКРЕЧЕНО</span>
      </div>

      <p class="file-id">АРХІВНИЙ ЗАПИС № 001&nbsp;&nbsp;·&nbsp;&nbsp;РІВЕНЬ ДОСТУПУ: ALPHA</p>

      <h1 class="dossier-title">Розшифровка протоколу</h1>

      <p class="dossier-text">
        Ти вже бачив це ім'я на екрані входу. Кожна літера абревіатури —
        окреме слово розшифровки. Фрагмент, який раніше приховувала чистка
        даних, тепер повністю відновлено.
      </p>

      <p class="acronym-line">
        J.A.R.V.I.S. = Just&nbsp;A&nbsp;Rather&nbsp;Very&nbsp;<span class="revealed-word">Intelligent</span>&nbsp;System
      </p>
    </div>
  </div>
</template>

<style scoped>
.record-page {
  width: 100%;
  display: flex;
  justify-content: center;
}
.dossier-card {
  position: relative;
  width: 100%;
  max-width: 560px;
  background: #ece7d8;
  background-image:
    repeating-linear-gradient(0deg, rgba(0,0,0,0.015) 0px, transparent 1px, transparent 2px),
    radial-gradient(ellipse at 30% 20%, rgba(0,0,0,0.06), transparent 60%),
    radial-gradient(ellipse at 80% 90%, rgba(0,0,0,0.08), transparent 55%);
  color: #1c1a14;
  padding: clamp(20px, 4vw, 32px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.4);
  clip-path: polygon(
    0% 1%, 2% 0%, 98% 0%, 100% 1%,
    100% 98%, 99% 100%, 97% 99%, 3% 100%, 1% 99%, 0% 98%
  );
}
.stamp {
  position: absolute;
  top: clamp(12px, 3.5vw, 22px);
  right: clamp(12px, 3.5vw, 22px);
  border: 3px solid #1a7a5e;
  color: #1a7a5e;
  font-family: 'Oswald', sans-serif;
  font-weight: 700;
  font-size: clamp(10px, 2.2vw, 12px);
  letter-spacing: 2px;
  padding: 6px 10px;
  transform: rotate(-8deg);
  opacity: 0.85;
  mix-blend-mode: multiply;
}
.file-id {
  font-family: 'Share Tech Mono', monospace;
  font-size: clamp(10px, 2.2vw, 12px);
  letter-spacing: 1px;
  color: #6b6252;
  margin: 0 0 16px;
}
.dossier-title {
  font-family: 'Special Elite', monospace;
  font-size: clamp(18px, 4.6vw, 24px);
  margin: 0 0 14px;
  line-height: 1.2;
}
.dossier-text {
  font-family: 'Special Elite', monospace;
  font-size: clamp(13px, 3vw, 15px);
  line-height: 1.7;
  margin: 0 0 16px;
}
.acronym-line {
  font-family: 'Special Elite', monospace;
  font-size: clamp(13px, 3.2vw, 15px);
  line-height: 1.8;
  margin: 0;
  word-break: break-word;
}
.revealed-word {
  font-weight: 700;
  text-decoration: underline;
  text-decoration-color: #1a7a5e;
  text-underline-offset: 3px;
}
</style>