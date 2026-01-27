import { Link } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import { ArrowRight, Code, Zap, Shield } from 'lucide-react'

export default function HomePage() {
  const { isAuthenticated } = useAuthStore()

  return (
    <div className="space-y-16">
      <section className="text-center py-16">
        <h1 className="text-5xl font-bold text-gray-900 mb-6">
          Welcome to <span className="text-indigo-600">GoFiver</span>
        </h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto mb-8">
          A modern blog platform built with Go, Fiber Framework, and React.
          Fast, secure, and enterprise-ready.
        </p>
        <div className="flex justify-center space-x-4">
          <Link
            to="/blogs"
            className="bg-indigo-600 text-white px-6 py-3 rounded-lg hover:bg-indigo-700 transition flex items-center space-x-2"
          >
            <span>Browse Blogs</span>
            <ArrowRight size={18} />
          </Link>
          {!isAuthenticated && (
            <Link
              to="/register"
              className="border border-indigo-600 text-indigo-600 px-6 py-3 rounded-lg hover:bg-indigo-50 transition"
            >
              Get Started
            </Link>
          )}
        </div>
      </section>

      <section className="grid md:grid-cols-3 gap-8">
        <div className="bg-white p-8 rounded-xl shadow-sm border text-center">
          <div className="w-14 h-14 bg-indigo-100 rounded-lg flex items-center justify-center mx-auto mb-4">
            <Zap className="text-indigo-600" size={28} />
          </div>
          <h3 className="text-lg font-semibold mb-2">Lightning Fast</h3>
          <p className="text-gray-600">
            Built on Fiber, one of the fastest web frameworks for Go.
          </p>
        </div>

        <div className="bg-white p-8 rounded-xl shadow-sm border text-center">
          <div className="w-14 h-14 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-4">
            <Shield className="text-green-600" size={28} />
          </div>
          <h3 className="text-lg font-semibold mb-2">Secure by Design</h3>
          <p className="text-gray-600">
            JWT authentication and enterprise-grade security practices.
          </p>
        </div>

        <div className="bg-white p-8 rounded-xl shadow-sm border text-center">
          <div className="w-14 h-14 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-4">
            <Code className="text-purple-600" size={28} />
          </div>
          <h3 className="text-lg font-semibold mb-2">Modern Stack</h3>
          <p className="text-gray-600">
            Go backend with React frontend for the best developer experience.
          </p>
        </div>
      </section>

      {/* CTA Section */}
        <h2 className="text-3xl font-bold mb-4">Ready to Start Blogging?</h2>
        <p className="text-indigo-100 mb-6 max-w-xl mx-auto">
          Join our community and share your thoughts with the world.
        </p>
        {isAuthenticated ? (
          <Link
            to="/blogs/create"
            className="bg-white text-indigo-600 px-6 py-3 rounded-lg hover:bg-indigo-50 transition inline-block font-medium"
          >
            Create Your First Blog
          </Link>
        ) : (
          <Link
            to="/register"
            className="bg-white text-indigo-600 px-6 py-3 rounded-lg hover:bg-indigo-50 transition inline-block font-medium"
          >
            Sign Up Now
          </Link>
        )}
    </div>
  )
}
