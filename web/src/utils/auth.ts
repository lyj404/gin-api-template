export interface UserInfo {
  id: string
  name: string
}

export const getUserInfo = (): UserInfo | null => {
  const token = localStorage.getItem('accessToken')
  if (!token) return null
  try {
    const payload = token.split('.')[1]
    const base64 = payload.replace(/-/g, '+').replace(/_/g, '/')
    const padded = base64 + '='.repeat((4 - base64.length % 4) % 4)
    const binary = atob(padded)
    const bytes = Uint8Array.from(binary, c => c.charCodeAt(0))
    const decoded = JSON.parse(new TextDecoder('utf-8').decode(bytes))
    return { id: decoded.id, name: decoded.name }
  } catch {
    return null
  }
}
