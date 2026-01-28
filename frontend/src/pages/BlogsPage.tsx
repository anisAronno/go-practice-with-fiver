import { useEffect, useState, useCallback } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { blogService } from '../services'
import type { Blog } from '../types'
import { ChevronLeft, ChevronRight, Calendar, User, Search, X } from 'lucide-react'

const API_URL = import.meta.env.VITE_API_URL || ''

const getImageUrl = (path: string) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${API_URL}${path}`
}

export default function BlogsPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(parseInt(searchParams.get('page') || '1'))
  const [totalPages, setTotalPages] = useState(1)
  const [search, setSearch] = useState(searchParams.get('search') || '')
  const [searchInput, setSearchInput] = useState(searchParams.get('search') || '')

  const updateUrl = useCallback((newPage: number, newSearch: string) => {
    const params = new URLSearchParams()
    if (newPage > 1) params.set('page', String(newPage))
    if (newSearch) params.set('search', newSearch)
    setSearchParams(params)
  }, [setSearchParams])

  const fetchBlogs = useCallback(async () => {
    setLoading(true)
    try {
      const response = await blogService.getAll(page, 20, search)
      setBlogs(response.data.data || [])
      setTotalPages(response.data.meta?.total_pages || 1)
    } catch (error) {
      console.error('Failed to fetch blogs:', error)
    } finally {
      setLoading(false)
    }
  }, [page, search])

  useEffect(() => {
    fetchBlogs()
  }, [fetchBlogs])

  useEffect(() => {
    const timer = setTimeout(() => {
      if (searchInput !== search) {
        setSearch(searchInput)
        setPage(1)
        updateUrl(1, searchInput)
      }
    }, 300)
    return () => clearTimeout(timer)
  }, [searchInput, search, updateUrl])

  useEffect(() => {
    updateUrl(page, search)
  }, [page, search, updateUrl])

  const clearSearch = () => {
    setSearchInput('')
    setSearch('')
    setPage(1)
    updateUrl(1, '')
  }

  const formatDate = (date: string) => new Date(date).toLocaleDateString()

  return (
    <div>
      <div className="flex flex-wrap gap-4 mb-8">
        <h1 className="text-3xl font-bold text-gray-900">All Blogs</h1>
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3 top-4 text-gray-400" size={20} />
          <input
            type="text"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            placeholder="Search blogs..."
            className="w-full pl-10 pr-10 py-2.5 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
          {searchInput && (
            <button
              onClick={clearSearch}
              className="absolute right-3 top-4 text-gray-400 hover:text-red-600"
            >
              <X size={18} />
            </button>
          )}
        </div>
      </div>

      {loading ? (
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
      ) : blogs.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border">
          <p className="text-gray-500">
            {search ? `No blogs found for "${search}"` : 'No blogs found'}
          </p>
          {search && (
            <button onClick={clearSearch} className="mt-4 text-indigo-600 hover:underline">
              Clear search
            </button>
          )}
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
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page === 1}
                className="flex items-center space-x-1 px-4 py-2 border rounded-lg disabled:opacity-50 hover:bg-gray-50"
              >
                <ChevronLeft size={18} />
                <span>Previous</span>
              </button>
              <span className="text-gray-600">Page {page} of {totalPages}</span>
              <button
                onClick={() => setPage(Math.min(totalPages, page + 1))}
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
