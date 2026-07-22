export function compactParams<T extends object>(params: T) {
  return Object.fromEntries(
    Object.entries(params).filter(([, value]) => {
      if (value === undefined || value === null || value === '') {
        return false
      }
      return !Array.isArray(value) || value.length > 0
    })
  )
}
