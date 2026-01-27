import { useEffect } from 'react'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import Layout from './components/Layout'
import HomePage from './pages/HomePage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import BlogsPage from './pages/BlogsPage'
import BlogDetailPage from './pages/BlogDetailPage'
import CreateBlogPage from './pages/CreateBlogPage'
import EditBlogPage from './pages/EditBlogPage'
import DashboardPage from './pages/DashboardPage'
import ProtectedRoute from './components/ProtectedRoute'
import GuestRoute from './components/GuestRoute'
import { useAuthStore } from './store/authStore'
import { authService } from './services'

function AuthLoader({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, token, _hasHydrated, setAuth, logout } = useAuthStore()

  useEffect(() => {
    const refreshUser = async () => {
      // Only verify token after hydration completes and we have auth
      if (!_hasHydrated) return
      
      if (token && isAuthenticated) {
        try {
          const response = await authService.me()
          if (response.data.data) {
            setAuth(response.data.data, token)
          }
        } catch {
          // Token invalid, clear auth
          logout()
        }
      }
    }
    refreshUser()
  }, [_hasHydrated, isAuthenticated, token])

  return <>{children}</>
}

const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      { index: true, element: <HomePage /> },
      { path: 'blogs', element: <BlogsPage /> },
      { path: 'blogs/:id', element: <BlogDetailPage /> },
      {
        element: <GuestRoute />,
        children: [
          { path: 'login', element: <LoginPage /> },
          { path: 'register', element: <RegisterPage /> },
        ],
      },
      {
        element: <ProtectedRoute />,
        children: [
          { path: 'dashboard', element: <DashboardPage /> },
          { path: 'blogs/create', element: <CreateBlogPage /> },
          { path: 'blogs/:id/edit', element: <EditBlogPage /> },
        ],
      },
    ],
  },
], {
  future: {
    v7_relativeSplatPath: true,
  },
})

function App() {
  return (
    <AuthLoader>
      <Toaster position="top-right" />
      <RouterProvider router={router} future={{ v7_startTransition: true }} />
    </AuthLoader>
  )
}

export default App
