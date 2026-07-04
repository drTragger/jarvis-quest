const KEY = 'jarvis_hub_onboarding_seen'

export function hasSeenOnboarding() {
	return localStorage.getItem(KEY) === '1'
}

export function markOnboardingSeen() {
	localStorage.setItem(KEY, '1')
}

export function resetOnboarding() {
	localStorage.removeItem(KEY)
}