import { useEffect, useState } from 'react'
import { blogService } from '../services'
import type { Blog } from '../types'
import BlogCard from '../components/BlogCard'
import LoadingSpinner from '../components/LoadingSpinner'
import { ChevronLeft, ChevronRight } from 'lucide-react'

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
      <h1 className="text-3xl font-bold text-gray-900 mb-8">All Blogs</h1>

      {blogs.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border">
          <p className="text-gray-500">No blogs found</p>
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
