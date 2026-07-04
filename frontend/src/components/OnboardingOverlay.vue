<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { playSound } from '../composables/useSound'

const props = defineProps({
	steps: { type: Array, required: true }
})
const emit = defineEmits(['finish'])

const stepIndex = ref(0)
const rect = ref(null)

const currentStep = computed(() => props.steps[stepIndex.value])
const isLastStep = computed(() => stepIndex.value === props.steps.length - 1)

function measureTarget() {
	const step = currentStep.value
	if (!step?.selector) {
		rect.value = null
		return
	}
	const el = document.querySelector(step.selector)
	if (!el) {
		rect.value = null
		return
	}
	const r = el.getBoundingClientRect()
	rect.value = { top: r.top, left: r.left, width: r.width, height: r.height }
}

const spotlightStyle = computed(() => {
	if (!rect.value) return null
	const pad = 8
	return {
		top: `${rect.value.top - pad}px`,
		left: `${rect.value.left - pad}px`,
		width: `${rect.value.width + pad * 2}px`,
		height: `${rect.value.height + pad * 2}px`
	}
})

const cardStyle = computed(() => {
	if (!rect.value) return null
	const spaceBelow = window.innerHeight - rect.value.top - rect.value.height
	const placeBelow = spaceBelow > 180
	const maxLeft = window.innerWidth - 300
	const left = Math.max(12, Math.min(rect.value.left, maxLeft))
	return placeBelow
		? { top: `${rect.value.top + rect.value.height + 20}px`, left: `${left}px` }
		: { top: `${rect.value.top - 20}px`, left: `${left}px`, transform: 'translateY(-100%)' }
})

let currentHowl = null

function playStepAudio() {
	if (currentHowl) currentHowl.stop()
	const step = currentStep.value
	currentHowl = step?.audioSrc
		? playSound(`onboarding-${stepIndex.value}`, step.audioSrc, { volume: 0.9 })
		: null
}

function next() {
	if (isLastStep.value) {
		finish()
		return
	}
	stepIndex.value++
}

function finish() {
	if (currentHowl) currentHowl.stop()
	emit('finish')
}

watch(stepIndex, async () => {
	await nextTick()
	measureTarget()
	playStepAudio()
})

function handleResize() {
	measureTarget()
}

onMounted(async () => {
	await nextTick()
	measureTarget()
	playStepAudio()
	window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
	window.removeEventListener('resize', handleResize)
})
</script>

<template>
	<div class="onboarding-root">
		<div
			class="onboarding-backdrop"
			:class="{ 'has-spotlight': !!spotlightStyle }"
			:style="spotlightStyle || {}"
		></div>

		<div class="onboarding-card" :class="{ 'is-centered': !rect }" :style="rect ? cardStyle : {}">
			<p class="onboarding-eyebrow">JARVIS &middot; {{ stepIndex + 1 }}/{{ steps.length }}</p>
			<h3 class="onboarding-title">{{ currentStep.title }}</h3>
			<p class="onboarding-text">{{ currentStep.text }}</p>
			<div class="onboarding-actions">
				<button class="onboarding-skip" @click="finish">ПРОПУСТИТИ</button>
				<button class="onboarding-next" @click="next">{{ isLastStep ? 'ЗРОЗУМІЛО' : 'ДАЛІ' }}</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
.onboarding-root {
	position: fixed;
	inset: 0;
	z-index: 100;
}
.onboarding-backdrop {
	position: fixed;
	inset: 0;
	background: rgba(5, 7, 10, 0.8);
	transition: all 0.25s ease;
}
.onboarding-backdrop.has-spotlight {
	background: transparent;
	border-radius: 8px;
	box-shadow: 0 0 0 9999px rgba(5, 7, 10, 0.8);
	border: 1px solid var(--jarvis-cyan);
}
.onboarding-card {
	position: fixed;
	max-width: 280px;
	padding: 16px 18px;
	background: rgba(13, 18, 32, 0.95);
	border: 1px solid rgba(126, 20, 255, 0.5);
	border-radius: 6px;
	box-shadow: 0 0 24px rgba(71, 191, 255, 0.2);
}
.onboarding-card.is-centered {
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	max-width: 320px;
	text-align: center;
}
.onboarding-eyebrow {
	font-family: var(--jarvis-mono);
	font-size: 10px;
	letter-spacing: 1.5px;
	color: var(--jarvis-cyan);
	margin: 0 0 8px;
}
.onboarding-title {
	font-family: var(--jarvis-display);
	font-weight: 700;
	font-size: clamp(14px, 3.4vw, 16px);
	margin: 0 0 8px;
}
.onboarding-text {
	font-size: clamp(12px, 2.8vw, 13px);
	line-height: 1.5;
	color: var(--jarvis-text-dim);
	margin: 0 0 14px;
}
.onboarding-actions {
	display: flex;
	justify-content: space-between;
	gap: 10px;
}
.onboarding-skip,
.onboarding-next {
	background: none;
	border: 1px solid rgba(232, 230, 240, 0.25);
	border-radius: 4px;
	padding: 6px 10px;
	color: var(--jarvis-text-dim);
	font-family: var(--jarvis-mono);
	font-size: 11px;
	letter-spacing: 1px;
	cursor: pointer;
}
.onboarding-next {
	border-color: var(--jarvis-cyan);
	color: var(--jarvis-cyan);
}
.onboarding-skip:hover,
.onboarding-next:hover {
	background: rgba(255, 255, 255, 0.05);
}
</style>