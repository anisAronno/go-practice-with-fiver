import { useEffect } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import Layout from './components/Layout'
import HomePage from './pages/HomePage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import BlogsPage from './pages/BlogsPage'
import BlogDetailPage from './pages/BlogDetailPage'
import CreateBlogPage from './pages/CreateBlogPage'
import EditBlogPage from './pages/EditBlogPage'
import MyBlogsPage from './pages/MyBlogsPage'
import DashboardPage from './pages/DashboardPage'
import ProtectedRoute from './components/ProtectedRoute'
import GuestRoute from './components/GuestRoute'
import { useAuthStore } from './store/authStore'
import { authService } from './services'

function AppContent() {
  const { isAuthenticated, setAuth, logout } = useAuthStore()

  useEffect(() => {
    const refreshUser = async () => {
      if (isAuthenticated) {
        try {
          const response = await authService.me()
          const token = localStorage.getItem('token')
          if (token && response.data.data) {
            setAuth(response.data.data, token)
          }
        } catch {
          logout()
        }
      }
    }
    refreshUser()
  }, [])

  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />

        <Route element={<GuestRoute />}>
          <Route path="login" element={<LoginPage />} />
          <Route path="register" element={<RegisterPage />} />
        </Route>

        <Route path="blogs" element={<BlogsPage />} />
        <Route path="blogs/:id" element={<BlogDetailPage />} />

        <Route element={<ProtectedRoute />}>
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="blogs/create" element={<CreateBlogPage />} />
          <Route path="blogs/:id/edit" element={<EditBlogPage />} />
          <Route path="my-blogs" element={<MyBlogsPage />} />
        </Route>
      </Route>
    </Routes>
  )
}

function App() {
  return (
    <BrowserRouter>
      <Toaster position="top-right" />
      <AppContent />
    </BrowserRouter>
  )
}

export default App