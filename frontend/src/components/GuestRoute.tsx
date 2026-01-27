import { Navigate, Outlet } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import LoadingSpinner from './LoadingSpinner'

export default function GuestRoute() {
  const { isAuthenticated, _hasHydrated } = useAuthStore()

  // Wait for Zustand to rehydrate from localStorage
  if (!_hasHydrated) {
    return <LoadingSpinner />
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />
  }

  return <Outlet />
}