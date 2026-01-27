import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { blogService } from '../services'
import type { Blog } from '../types'
import { Zap, Shield, Layers, ArrowRight, Clock } from 'lucide-react'

export default function HomePage() {
  const [recentBlogs, setRecentBlogs] = useState<Blog[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadRecentBlogs()
  }, [])

  const loadRecentBlogs = async () => {
    try {
      const res = await blogService.getAll(1, 3)
      setRecentBlogs(res.data.data || [])
    } catch (err) {
      console.error('Failed to load blogs:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="-mt-8 -mx-4 sm:-mx-6 lg:-mx-8">
      <section className="bg-gradient-to-br from-indigo-600 via-indigo-700 to-purple-800 text-white py-24 px-4">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-5xl md:text-6xl font-bold mb-6">
            Build Fast with
            <span className="block text-indigo-200">Go & Fiber</span>
          </h1>
          <p className="text-xl text-indigo-100 mb-8 max-w-2xl mx-auto">
            A modern full-stack boilerplate featuring Go, Fiber framework, React, TypeScript, and Docker. Production-ready from day one.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/blogs"
              className="inline-flex items-center justify-center gap-2 bg-white text-indigo-700 px-8 py-3 rounded-xl font-semibold hover:bg-indigo-50 transition"
            >
              Explore Blog
              <ArrowRight size={20} />
            </Link>
            <Link
              to="/register"
              className="inline-flex items-center justify-center gap-2 bg-indigo-500 text-white px-8 py-3 rounded-xl font-semibold hover:bg-indigo-400 transition border border-indigo-400"
            >
              Get Started
            </Link>
          </div>
        </div>
      </section>

      <section className="py-20 px-4 bg-gray-50">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">Why GoFiver?</h2>
            <p className="text-gray-600 max-w-2xl mx-auto">
              Everything you need to build modern web applications with the power of Go.
            </p>
          </div>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-white p-8 rounded-2xl shadow-sm border hover:shadow-md transition">
              <div className="w-14 h-14 bg-indigo-100 rounded-xl flex items-center justify-center mb-6">
                <Zap className="text-indigo-600" size={28} />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Blazing Fast</h3>
              <p className="text-gray-600">
                Built on Fiber, the fastest HTTP framework for Go. Handle millions of requests with minimal resources.
              </p>
            </div>
            <div className="bg-white p-8 rounded-2xl shadow-sm border hover:shadow-md transition">
              <div className="w-14 h-14 bg-green-100 rounded-xl flex items-center justify-center mb-6">
                <Shield className="text-green-600" size={28} />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Secure by Default</h3>
              <p className="text-gray-600">
                JWT authentication, role-based access control, and industry best practices baked in.
              </p>
            </div>
            <div className="bg-white p-8 rounded-2xl shadow-sm border hover:shadow-md transition">
              <div className="w-14 h-14 bg-purple-100 rounded-xl flex items-center justify-center mb-6">
                <Layers className="text-purple-600" size={28} />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">Production Ready</h3>
              <p className="text-gray-600">
                Docker, Nginx reverse proxy, and scalable architecture. Deploy with confidence.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section className="py-20 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-between mb-12">
            <div>
              <h2 className="text-3xl font-bold text-gray-900 mb-2">Recent Posts</h2>
              <p className="text-gray-600">Discover the latest articles from our community</p>
            </div>
            <Link
              to="/blogs"
              className="hidden sm:inline-flex items-center gap-2 text-indigo-600 font-medium hover:text-indigo-700"
            >
              View All <ArrowRight size={18} />
            </Link>
          </div>

          {loading ? (
            <div className="text-center py-12">
              <div className="animate-spin w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full mx-auto"></div>
            </div>
          ) : recentBlogs.length > 0 ? (
            <div className="grid md:grid-cols-3 gap-8">
              {recentBlogs.map(blog => (
                <article key={blog.id} className="bg-white rounded-2xl shadow-sm border overflow-hidden hover:shadow-md transition group">
                  <div className="h-40 bg-gradient-to-br from-indigo-100 to-purple-100 flex items-center justify-center">
                    <span className="text-6xl font-bold text-indigo-200 group-hover:scale-110 transition">
                      {blog.title.charAt(0)}
                    </span>
                  </div>
                  <div className="p-6">
                    <Link to={`/blogs/${blog.id}`}>
                      <h3 className="text-lg font-semibold text-gray-900 mb-2 group-hover:text-indigo-600 transition line-clamp-2">
                        {blog.title}
                      </h3>
                    </Link>
                    <p className="text-gray-600 text-sm mb-4 line-clamp-2">
                      {blog.content}
                    </p>
                    <div className="flex items-center justify-between text-sm text-gray-500">
                      <span>{blog.user?.name || 'Anonymous'}</span>
                      <span className="flex items-center gap-1">
                        <Clock size={14} />
                        {new Date(blog.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  </div>
                </article>
              ))}
            </div>
          ) : (
            <div className="text-center py-12 text-gray-500">
              No posts yet. Be the first to write one!
            </div>
          )}

          <div className="text-center mt-8 sm:hidden">
            <Link
              to="/blogs"
              className="inline-flex items-center gap-2 text-indigo-600 font-medium"
            >
              View All Posts <ArrowRight size={18} />
            </Link>
          </div>
        </div>
      </section>

      <section className="py-20 px-4 bg-indigo-600">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl font-bold text-white mb-4">
            Ready to Get Started?
          </h2>
          <p className="text-indigo-100 mb-8 text-lg">
            Join our community and start building amazing applications today.
          </p>
          <Link
            to="/register"
            className="inline-flex items-center gap-2 bg-white text-indigo-700 px-8 py-3 rounded-xl font-semibold hover:bg-indigo-50 transition"
          >
            Create Account
            <ArrowRight size={20} />
          </Link>
        </div>
      </section>
    </div>
  )
}
