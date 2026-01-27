import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Only clear auth and redirect for explicit 401s on protected endpoints
    // Don't redirect during initial page load or for public endpoints
    if (error.response?.status === 401) {
      const url = error.config?.url || ''
      // Only redirect if accessing protected endpoints (not login/register/public)
      const isProtectedEndpoint = url.includes('/auth/me') || 
        url.includes('/blogs/my') ||
        url.includes('/dashboard') ||
        (url.includes('/users') && !url.includes('/blogs'))
      
      if (isProtectedEndpoint) {
        localStorage.removeItem('token')
        localStorage.removeItem('auth-storage')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
