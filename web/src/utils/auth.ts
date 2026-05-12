export interface UserInfo {
  id: string
  name: string
}

export const getUserInfo = (): UserInfo | null => {
  const token = localStorage.getItem('accessToken')
  if (!token) return null
  try {
    const payload = token.split('.')[1]
    const decoded = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
    return { id: decoded.id, name: decoded.name }
  } catch {
    return null
  }
}
