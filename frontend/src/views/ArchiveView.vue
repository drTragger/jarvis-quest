<script setup>
import { ref, onMounted, computed, defineAsyncComponent, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import JarvisLayout from '../layouts/JarvisLayout.vue'
import { getPassword } from '../lib/auth'
import { stepIds, collectibles, archiveRecords } from '../data/steps'

const router = useRouter()
const unlocked = ref(null)
const openId = ref(null)

const completedSteps = computed(() => {
  if (unlocked.value === null) return []
  return stepIds.filter((id) => id < unlocked.value)
})

onMounted(async () => {
  const res = await fetch(`/api/progress?password=${encodeURIComponent(getPassword())}`)
  if (res.status === 401) {
    router.push({ name: 'login' })
    return
  }
  const data = await res.json()
  unlocked.value = data.unlocked
})

function toggle(id) {
  openId.value = openId.value === id ? null : id
}

const recordComponentCache = {}
function getRecordComponent(id) {
  if (!archiveRecords[id]) return null
  if (!recordComponentCache[id]) {
    recordComponentCache[id] = markRaw(defineAsyncComponent(archiveRecords[id]))
  }
  return recordComponentCache[id]
}
</script>

<template>
  <JarvisLayout eyebrow="АРХІВ ДОСЬЄ" wide>
    <button class="archive-back" @click="router.push({ name: 'hub' })">&lt; НАЗАД ДО МЕНЮ</button>

    <div class="archive-empty" v-if="unlocked !== null && completedSteps.length === 0">
      Архів порожній. Розблокуйте перший запис, пройшовши місію.
    </div>

    <div class="archive-list" v-else>
      <div v-for="id in completedSteps" :key="id" class="archive-entry" :class="{ 'is-open': openId === id }">
        <button class="archive-entry-head" @click="toggle(id)">
          <span class="archive-entry-date">{{ collectibles[id]?.date || `ЗАПИС №${id}` }}</span>
          <span class="archive-entry-title">{{ collectibles[id]?.title || `Сегмент ${id}` }}</span>
        </button>

        <div v-if="openId === id" class="archive-entry-body">
          <component :is="getRecordComponent(id)" v-if="getRecordComponent(id)" />
          <p v-else class="archive-entry-text">{{ collectibles[id]?.excerpt || 'Опис відсутній.' }}</p>
        </div>
      </div>
    </div>
  </JarvisLayout>
</template>

<style scoped>
.archive-back {
  align-self: flex-start;
  background: rgba(5, 7, 10, 0.75);
  border: 1px solid rgba(126, 20, 255, 0.4);
  border-radius: 4px;
  padding: 6px 10px;
  color: var(--jarvis-cyan);
  font-family: var(--jarvis-mono);
  font-size: clamp(11px, 2.6vw, 13px);
  letter-spacing: 1px;
  cursor: pointer;
  margin-bottom: 20px;
}
.archive-back:hover,
.archive-back:focus-visible {
  border-color: var(--jarvis-cyan);
  background: rgba(5, 7, 10, 0.9);
}
.archive-empty {
  color: var(--jarvis-text-dim);
  font-size: clamp(12px, 3vw, 14px);
}
.archive-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}
.archive-entry {
  border: 1px solid rgba(126, 20, 255, 0.3);
  border-radius: 4px;
  overflow: hidden;
  transition: border-color 0.2s;
}
.archive-entry:hover,
.archive-entry.is-open {
  border-color: var(--jarvis-cyan);
}
.archive-entry-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  width: 100%;
  padding: 14px 16px;
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
}
.archive-entry-date {
  font-family: var(--jarvis-mono);
  font-size: clamp(10px, 2.4vw, 12px);
  color: var(--jarvis-text-dim);
}
.archive-entry-title {
  font-family: var(--jarvis-display);
  font-weight: 700;
  font-size: clamp(13px, 3vw, 15px);
  color: var(--jarvis-text);
}
.archive-entry-body {
  padding: 0 16px 20px;
}
.archive-entry-text {
  margin: 0;
  font-size: clamp(12px, 2.8vw, 14px);
  color: var(--jarvis-text-dim);
  line-height: 1.6;
}
</style>