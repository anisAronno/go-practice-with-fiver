import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { blogService } from '../services'
import type { Blog } from '../types'
import { ChevronLeft, ChevronRight, Calendar, User } from 'lucide-react'

const API_URL = import.meta.env.VITE_API_URL || ''

const getImageUrl = (path: string) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${API_URL}${path}`
}

export default function BlogsPage() {
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)

  useEffect(() => {
    fetchBlogs()
  }, [page])

  const fetchBlogs = async () => {
    setLoading(true)
    try {
      const response = await blogService.getAll(page)
      setBlogs(response.data.data || [])
      setTotalPages(response.data.meta?.total_pages || 1)
    } catch (error) {
      console.error('Failed to fetch blogs:', error)
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (date: string) => new Date(date).toLocaleDateString()

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {[...Array(6)].map((_, i) => (
          <div key={i} className="bg-white rounded-xl border overflow-hidden animate-pulse">
            <div className="h-48 bg-gray-200" />
            <div className="p-5 space-y-3">
              <div className="h-6 bg-gray-200 rounded w-3/4" />
              <div className="h-4 bg-gray-200 rounded w-1/2" />
            </div>
          </div>
        ))}
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-8">All Blogs</h1>

      {blogs.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border">
          <p className="text-gray-500">No blogs found</p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {blogs.map((blog) => (
              <Link
                key={blog.id}
                to={`/blogs/${blog.id}`}
                className="bg-white rounded-xl border overflow-hidden hover:shadow-lg transition group"
              >
                <div className="h-48 bg-gray-100 overflow-hidden">
                  {blog.image ? (
                    <img
                      src={getImageUrl(blog.image)}
                      alt={blog.title}
                      className="w-full h-full object-cover group-hover:scale-105 transition"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-indigo-100 to-purple-100">
                      <span className="text-4xl font-bold text-indigo-300">
                        {blog.title.charAt(0)}
                      </span>
                    </div>
                  )}
                </div>
                <div className="p-5">
                  <h2 className="text-lg font-semibold text-gray-900 group-hover:text-indigo-600 transition line-clamp-2">
                    {blog.title}
                  </h2>
                  <div className="flex items-center gap-4 mt-3 text-sm text-gray-500">
                    {blog.user && (
                      <span className="flex items-center gap-1">
                        <User size={14} />
                        {blog.user.name}
                      </span>
                    )}
                    <span className="flex items-center gap-1">
                      <Calendar size={14} />
                      {formatDate(blog.created_at)}
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>

          {totalPages > 1 && (
            <div className="flex justify-center items-center space-x-4 mt-8">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="flex items-center space-x-1 px-4 py-2 border rounded-lg disabled:opacity-50 hover:bg-gray-50"
              >
                <ChevronLeft size={18} />
                <span>Previous</span>
              </button>
              <span className="text-gray-600">Page {page} of {totalPages}</span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="flex items-center space-x-1 px-4 py-2 border rounded-lg disabled:opacity-50 hover:bg-gray-50"
              >
                <span>Next</span>
                <ChevronRight size={18} />
              </button>
            </div>
          )}
        </>
      )}
    </div>
  )
}
