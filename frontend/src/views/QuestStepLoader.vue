<script setup>
import { computed, defineAsyncComponent } from 'vue'
import { useRouter } from 'vue-router'
import { stepComponents } from '../data/steps'

const props = defineProps({
  stepId: { type: [String, Number], required: true }
})

const router = useRouter()
const stepNum = computed(() => Number(props.stepId))
const loader = stepComponents[stepNum.value]

if (!loader) {
  router.replace({ name: 'not-found' })
}

const StepComponent = loader ? defineAsyncComponent(loader) : null
</script>

<template>
  <component :is="StepComponent" v-if="StepComponent" />
</template>