import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { blogService } from '../services'
import type { Blog } from '../types'
import BlogCard from '../components/BlogCard'
import LoadingSpinner from '../components/LoadingSpinner'
import { PlusCircle, ChevronLeft, ChevronRight } from 'lucide-react'

export default function MyBlogsPage() {
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)

  useEffect(() => {
    fetchMyBlogs()
  }, [page])

  const fetchMyBlogs = async () => {
    setLoading(true)
    try {
      const response = await blogService.getMyBlogs(page)
      setBlogs(response.data.data)
      setTotalPages(response.data.meta.total_pages)
    } catch (error) {
      console.error('Failed to fetch blogs:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <LoadingSpinner className="py-20" size={40} />
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">My Blogs</h1>
        <Link
          to="/blogs/create"
          className="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition flex items-center space-x-2"
        >
          <PlusCircle size={18} />
          <span>New Blog</span>
        </Link>
      </div>

      {blogs.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border">
          <p className="text-gray-500 mb-4">You haven't written any blogs yet</p>
          <Link
            to="/blogs/create"
            className="text-indigo-600 hover:underline font-medium"
          >
            Create your first blog →
          </Link>
        </div>
      ) : (
        <>
          <div className="grid gap-6">
            {blogs.map((blog) => (
              <BlogCard key={blog.id} blog={blog} />
            ))}
          </div>

          {totalPages > 1 && (
            <div className="flex justify-center items-center space-x-4 mt-8">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="flex items-center space-x-1 px-4 py-2 border rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                <ChevronLeft size={18} />
                <span>Previous</span>
              </button>
              <span className="text-gray-600">
                Page {page} of {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="flex items-center space-x-1 px-4 py-2 border rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
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
