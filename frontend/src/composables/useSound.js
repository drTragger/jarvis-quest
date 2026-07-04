import { Howl, Howler } from 'howler'

const cache = new Map()
let unlocked = false

export function isUnlocked() {
	return unlocked
}

export function unlockAudio(unlockSrc) {
	return new Promise((resolve) => {
		if (unlocked) {
			resolve()
			return
		}
		const primer = new Howl({
			src: [unlockSrc],
			volume: 0.01
		})
		primer.once('play', () => {
			unlocked = true
			resolve()
		})
		primer.once('loaderror', () => {
			unlocked = true
			resolve()
		})
		primer.once('playerror', () => {
			unlocked = true
			resolve()
		})
		primer.play()
	})
}

export function playSound(key, src, options = {}) {
	if (Howler.ctx && Howler.ctx.state === 'suspended') {
		Howler.ctx.resume()
	}

	let sound = cache.get(key)
	if (!sound) {
		sound = new Howl({
			src: [src],
			volume: options.volume ?? 1,
			loop: options.loop ?? false,
			rate: options.rate ?? 1
		})
		cache.set(key, sound)
	}
	sound.play()
	return sound
}

export function stopSound(key) {
	const sound = cache.get(key)
	if (sound) sound.stop()
}

if (typeof window !== 'undefined') {
	const resumeCtx = () => {
		if (Howler.ctx && Howler.ctx.state === 'suspended') {
			Howler.ctx.resume()
		}
	}
	window.addEventListener('click', resumeCtx)
	window.addEventListener('touchend', resumeCtx)
	document.addEventListener('visibilitychange', () => {
		if (document.visibilityState === 'visible') resumeCtx()
	})
}