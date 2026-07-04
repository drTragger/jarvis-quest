const KEY = 'jarvis-quest-password'

export function getPassword() {
  return localStorage.getItem(KEY)
}

export function setPassword(password) {
  localStorage.setItem(KEY, password)
}

export function clearPassword() {
  localStorage.removeItem(KEY)
}