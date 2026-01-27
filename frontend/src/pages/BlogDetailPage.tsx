import { useEffect, useState, useRef } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { blogService } from '../services'
import { useAuthStore } from '../store/authStore'
import type { Blog } from '../types'
import toast from 'react-hot-toast'
import { Calendar, User, Edit, Trash2, ArrowLeft } from 'lucide-react'

const API_URL = import.meta.env.VITE_API_URL || ''

export default function BlogDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [blog, setBlog] = useState<Blog | null>(null)
  const [loading, setLoading] = useState(true)
  const [notFound, setNotFound] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const { user, isAuthenticated } = useAuthStore()
  const navigate = useNavigate()
  const toastShown = useRef(false)

  useEffect(() => {
    fetchBlog()
  }, [id])

  const fetchBlog = async () => {
    if (!id) return
    setLoading(true)
    setNotFound(false)
    try {
      const response = await blogService.getById(Number(id))
      setBlog(response.data.data)
    } catch {
      setNotFound(true)
      if (!toastShown.current) {
        toast.error('Blog not found')
        toastShown.current = true
      }
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!blog || !window.confirm('Delete this blog?')) return
    setDeleting(true)
    try {
      await blogService.delete(blog.id)
      toast.success('Blog deleted')
      navigate('/blogs')
    } catch {
      toast.error('Failed to delete')
    } finally {
      setDeleting(false)
    }
  }

  const formatDate = (date: string) =>
    new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })

  const getImageUrl = (path: string) => {
    if (!path) return ''
    if (path.startsWith('http')) return path
    return `${API_URL}${path}`
  }

  if (loading) {
    return (
      <div className="animate-pulse">
        <div className="h-8 bg-gray-200 rounded w-24 mb-6" />
        <div className="bg-white rounded-xl border overflow-hidden">
          <div className="h-80 bg-gray-200" />
          <div className="p-8 space-y-4">
            <div className="h-8 bg-gray-200 rounded w-3/4" />
            <div className="h-4 bg-gray-200 rounded w-1/4" />
          </div>
        </div>
      </div>
    )
  }

  if (notFound) {
    return (
      <div className="text-center py-20">
        <h2 className="text-2xl font-bold text-gray-900 mb-4">Blog Not Found</h2>
        <p className="text-gray-500 mb-6">The blog you're looking for doesn't exist.</p>
        <Link to="/blogs" className="inline-flex items-center gap-2 text-indigo-600 hover:text-indigo-700">
          <ArrowLeft size={18} />
          Back to Blogs
        </Link>
      </div>
    )
  }

  if (!blog) return null

  const isOwner = isAuthenticated && user?.id === blog.user_id

  return (
    <div>
      <Link to="/blogs" className="inline-flex items-center text-gray-600 hover:text-indigo-600 mb-6">
        <ArrowLeft size={18} className="mr-1" />
        Back to Blogs
      </Link>

      <article className="bg-white rounded-xl shadow-sm border overflow-hidden">
        {blog.image && (
          <div className="h-80 md:h-96 overflow-hidden">
            <img src={getImageUrl(blog.image)} alt={blog.title} className="w-full h-full object-cover" />
          </div>
        )}

        <div className="p-8">
          <header className="mb-6">
            <h1 className="text-3xl font-bold text-gray-900 mb-4">{blog.title}</h1>
            <div className="flex items-center justify-between flex-wrap gap-4">
              <div className="flex items-center gap-4 text-gray-500">
                {blog.user && (
                  <span className="flex items-center gap-1">
                    <User size={16} />
                    {blog.user.name}
                  </span>
                )}
                <span className="flex items-center gap-1">
                  <Calendar size={16} />
                  {formatDate(blog.created_at)}
                </span>
              </div>

              {isOwner && (
                <div className="flex items-center gap-2">
                  <Link
                    to={`/blogs/${blog.id}/edit`}
                    className="flex items-center gap-1 text-indigo-600 hover:text-indigo-700 px-3 py-1 rounded-lg border border-indigo-200 hover:bg-indigo-50"
                  >
                    <Edit size={16} />
                    Edit
                  </Link>
                  <button
                    onClick={handleDelete}
                    disabled={deleting}
                    className="flex items-center gap-1 text-red-600 hover:text-red-700 px-3 py-1 rounded-lg border border-red-200 hover:bg-red-50 disabled:opacity-50"
                  >
                    <Trash2 size={16} />
                    {deleting ? 'Deleting...' : 'Delete'}
                  </button>
                </div>
              )}
            </div>
          </header>

          <div className="prose max-w-none">
            <p className="text-gray-700 whitespace-pre-wrap leading-relaxed">{blog.content}</p>
          </div>
        </div>
      </article>
    </div>
  )
}
