import { ref, onMounted, onUnmounted } from 'vue'

const EARTH_RADIUS_M = 6371000

function toRad(deg) {
  return (deg * Math.PI) / 180
}
function toDeg(rad) {
  return (rad * 180) / Math.PI
}

function haversineDistance(lat1, lon1, lat2, lon2) {
  const dLat = toRad(lat2 - lat1)
  const dLon = toRad(lon2 - lon1)
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) ** 2
  return EARTH_RADIUS_M * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
}

function bearingTo(lat1, lon1, lat2, lon2) {
  const phi1 = toRad(lat1)
  const phi2 = toRad(lat2)
  const dLon = toRad(lon2 - lon1)
  const y = Math.sin(dLon) * Math.cos(phi2)
  const x = Math.cos(phi1) * Math.sin(phi2) - Math.sin(phi1) * Math.cos(phi2) * Math.cos(dLon)
  return (toDeg(Math.atan2(y, x)) + 360) % 360
}

export function useGeoLocator(targetLat, targetLon, lockRadiusM = 40) {
  const distance = ref(null)
  const bearingToTarget = ref(null)
  const heading = ref(null)
  const headingAvailable = ref(false)
  const locked = ref(false)
  const geoError = ref(null)
  const compassPermissionNeeded = ref(false)

  let watchId = null

  function onPosition(pos) {
    const { latitude, longitude } = pos.coords
    distance.value = haversineDistance(latitude, longitude, targetLat, targetLon)
    bearingToTarget.value = bearingTo(latitude, longitude, targetLat, targetLon)
    if (distance.value <= lockRadiusM) locked.value = true
  }

  function onGeoError(err) {
    geoError.value = err.message || 'Геолокація недоступна'
  }

  function onOrientation(e) {
    let h = null
    if (typeof e.webkitCompassHeading === 'number') {
      h = e.webkitCompassHeading
    } else if (e.alpha !== null && e.alpha !== undefined) {
      h = (360 - e.alpha) % 360
    }
    if (h !== null) {
      heading.value = h
      headingAvailable.value = true
    }
  }

  function attachOrientationListeners() {
    window.addEventListener('deviceorientationabsolute', onOrientation, true)
    window.addEventListener('deviceorientation', onOrientation, true)
  }

  async function requestCompass() {
    if (
      typeof DeviceOrientationEvent !== 'undefined' &&
      typeof DeviceOrientationEvent.requestPermission === 'function'
    ) {
      try {
        const result = await DeviceOrientationEvent.requestPermission()
        if (result !== 'granted') return
      } catch {
        return
      }
    }
    compassPermissionNeeded.value = false
    attachOrientationListeners()
  }

  onMounted(() => {
    if (!navigator.geolocation) {
      geoError.value = 'Геолокація не підтримується цим пристроєм'
      return
    }
    watchId = navigator.geolocation.watchPosition(onPosition, onGeoError, {
      enableHighAccuracy: true,
      maximumAge: 2000,
      timeout: 15000
    })

    if (
      typeof DeviceOrientationEvent !== 'undefined' &&
      typeof DeviceOrientationEvent.requestPermission === 'function'
    ) {
      compassPermissionNeeded.value = true
    } else {
      attachOrientationListeners()
    }
  })

  onUnmounted(() => {
    if (watchId !== null) navigator.geolocation.clearWatch(watchId)
    window.removeEventListener('deviceorientationabsolute', onOrientation, true)
    window.removeEventListener('deviceorientation', onOrientation, true)
  })

  return {
    distance,
    bearingToTarget,
    heading,
    headingAvailable,
    locked,
    geoError,
    compassPermissionNeeded,
    requestCompass
  }
}
