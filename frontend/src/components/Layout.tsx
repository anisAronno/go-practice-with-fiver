import { useState } from 'react'
import { Outlet, Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import { Home, FileText, LogIn, LogOut, LayoutDashboard } from 'lucide-react'
import LogoutModal from './LogoutModal'

export default function Layout() {
  const { isAuthenticated, user, logout } = useAuthStore()
  const navigate = useNavigate()
  const [showLogoutModal, setShowLogoutModal] = useState(false)

  const handleLogout = () => {
    logout()
    setShowLogoutModal(false)
    navigate('/')
  }

  return (
    <div className="min-h-screen flex flex-col">
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link to="/" className="flex items-center space-x-2">
              <span className="text-2xl font-bold text-indigo-600">GoFiver</span>
            </Link>
            <nav className="flex items-center space-x-4">
              <Link
                to="/"
                className="flex items-center space-x-1 text-gray-600 hover:text-indigo-600 transition"
              >
                <Home size={18} />
                <span>Home</span>
              </Link>
              <Link
                to="/blogs"
                className="flex items-center space-x-1 text-gray-600 hover:text-indigo-600 transition"
              >
                <FileText size={18} />
                <span>Blog</span>
              </Link>

              {isAuthenticated ? (
                <>
                  <Link
                    to="/dashboard"
                    className="flex items-center space-x-1 text-gray-600 hover:text-indigo-600 transition"
                  >
                    <LayoutDashboard size={18} />
                    <span>Dashboard</span>
                  </Link>
                  <div className="flex items-center space-x-3 ml-4 pl-4 border-l">
                    <span className="text-sm text-gray-500">
                      {user?.name}
                      {user?.role === 'admin' && (
                        <span className="ml-1 px-1.5 py-0.5 text-xs bg-indigo-100 text-indigo-700 rounded">
                          Admin
                        </span>
                      )}
                    </span>
                    <button
                      onClick={() => setShowLogoutModal(true)}
                      className="flex items-center space-x-1 text-red-600 hover:text-red-700 transition"
                    >
                      <LogOut size={18} />
                      <span>Logout</span>
                    </button>
                  </div>
                </>
              ) : (
                <Link
                  to="/login"
                  className="flex items-center space-x-1 bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition"
                >
                  <LogIn size={18} />
                  <span>Login</span>
                </Link>
              )}
            </nav>
          </div>
        </div>
      </header>

      <main className="flex-1">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Outlet />
        </div>
      </main>

      <footer className="bg-white border-t py-4">
        <div className="max-w-7xl mx-auto px-4 text-center text-gray-500 text-sm">
          GoFiver - Built with Go, Fiber & React
        </div>
      </footer>

      <LogoutModal
        isOpen={showLogoutModal}
        onConfirm={handleLogout}
        onCancel={() => setShowLogoutModal(false)}
      />
    </div>
  )
}
