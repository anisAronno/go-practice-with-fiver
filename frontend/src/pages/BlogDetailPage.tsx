import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { blogService } from '../services'
import { useAuthStore } from '../store/authStore'
import type { Blog } from '../types'
import LoadingSpinner from '../components/LoadingSpinner'
import toast from 'react-hot-toast'
import { Calendar, User, Edit, Trash2, ArrowLeft } from 'lucide-react'

export default function BlogDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [blog, setBlog] = useState<Blog | null>(null)
  const [loading, setLoading] = useState(true)
  const [deleting, setDeleting] = useState(false)
  const { user, isAuthenticated } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    fetchBlog()
  }, [id])

  const fetchBlog = async () => {
    if (!id) return
    setLoading(true)
    try {
      const response = await blogService.getById(Number(id))
      setBlog(response.data.data)
    } catch (error) {
      console.error('Failed to fetch blog:', error)
      toast.error('Blog not found')
      navigate('/blogs')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!blog || !window.confirm('Are you sure you want to delete this blog?')) return
    
    setDeleting(true)
    try {
      await blogService.delete(blog.id)
      toast.success('Blog deleted successfully')
      navigate('/blogs')
    } catch (error) {
      toast.error('Failed to delete blog')
    } finally {
      setDeleting(false)
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  if (loading) {
    return <LoadingSpinner className="py-20" size={40} />
  }

  if (!blog) {
    return null
  }

  const isOwner = isAuthenticated && user?.id === blog.user_id

  return (
    <div className="max-w-3xl mx-auto">
      <Link
        to="/blogs"
        className="inline-flex items-center text-gray-600 hover:text-indigo-600 mb-6"
      >
        <ArrowLeft size={18} className="mr-1" />
        Back to Blogs
      </Link>

      <article className="bg-white rounded-xl shadow-sm border p-8">
        <header className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">{blog.title}</h1>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4 text-gray-500">
              {blog.user && (
                <span className="flex items-center space-x-1">
                  <User size={16} />
                  <span>{blog.user.name}</span>
                </span>
              )}
              <span className="flex items-center space-x-1">
                <Calendar size={16} />
                <span>{formatDate(blog.created_at)}</span>
              </span>
            </div>

            {isOwner && (
              <div className="flex items-center space-x-2">
                <Link
                  to={`/blogs/${blog.id}/edit`}
                  className="flex items-center space-x-1 text-indigo-600 hover:text-indigo-700 px-3 py-1 rounded-lg border border-indigo-200 hover:bg-indigo-50"
                >
                  <Edit size={16} />
                  <span>Edit</span>
                </Link>
                <button
                  onClick={handleDelete}
                  disabled={deleting}
                  className="flex items-center space-x-1 text-red-600 hover:text-red-700 px-3 py-1 rounded-lg border border-red-200 hover:bg-red-50 disabled:opacity-50"
                >
                  <Trash2 size={16} />
                  <span>{deleting ? 'Deleting...' : 'Delete'}</span>
                </button>
              </div>
            )}
          </div>
        </header>

        <div className="prose max-w-none">
          <p className="text-gray-700 whitespace-pre-wrap leading-relaxed">
            {blog.content}
          </p>
        </div>
      </article>
    </div>
  )
}
