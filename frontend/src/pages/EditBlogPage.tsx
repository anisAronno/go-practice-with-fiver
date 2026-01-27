import { useEffect, useState, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { blogService } from '../services'
import { useAuthStore } from '../store/authStore'
import type { Blog } from '../types'
import LoadingSpinner from '../components/LoadingSpinner'
import toast from 'react-hot-toast'
import { Save, ArrowLeft, Upload, Trash2 } from 'lucide-react'

export default function EditBlogPage() {
  const { id } = useParams<{ id: string }>()
  const [blog, setBlog] = useState<Blog | null>(null)
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [image, setImage] = useState<File | null>(null)
  const [preview, setPreview] = useState<string | null>(null)
  const [existingImage, setExistingImage] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const fileRef = useRef<HTMLInputElement>(null)
  const { user } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    fetchBlog()
  }, [id])

  const fetchBlog = async () => {
    if (!id) return
    setLoading(true)
    try {
      const response = await blogService.getById(Number(id))
      const blogData = response.data.data
      
      if (blogData.user_id !== user?.id && user?.role !== 'admin') {
        toast.error('You can only edit your own blogs')
        navigate('/blogs')
        return
      }
      
      setBlog(blogData)
      setTitle(blogData.title)
      setContent(blogData.content)
      setExistingImage(blogData.image || null)
    } catch {
      toast.error('Blog not found')
      navigate('/blogs')
    } finally {
      setLoading(false)
    }
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      if (file.size > 5 * 1024 * 1024) {
        toast.error('Image must be less than 5MB')
        return
      }
      setImage(file)
      setPreview(URL.createObjectURL(file))
      setExistingImage(null)
    }
  }

  const removeImage = async () => {
    if (existingImage && blog) {
      try {
        await blogService.deleteImage(blog.id)
        toast.success('Image removed')
      } catch {
        toast.error('Failed to remove image')
        return
      }
    }
    setImage(null)
    setPreview(null)
    setExistingImage(null)
    if (fileRef.current) fileRef.current.value = ''
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!blog) return
    
    setSaving(true)
    try {
      await blogService.update(blog.id, { title, content })

      if (image) {
        await blogService.uploadImage(blog.id, image)
      }

      toast.success('Blog updated successfully!')
      navigate(`/blogs/${blog.id}`)
    } catch (error: unknown) {
      const err = error as { response?: { data?: { message?: string } } }
      toast.error(err.response?.data?.message || 'Failed to update blog')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <LoadingSpinner className="py-20" size={40} />
  if (!blog) return null

  const currentImage = preview || existingImage

  return (
    <div>
      <button
        onClick={() => navigate(-1)}
        className="inline-flex items-center text-gray-600 hover:text-indigo-600 mb-6"
      >
        <ArrowLeft size={18} className="mr-1" />
        Back
      </button>

      <div className="bg-white rounded-xl shadow-sm border p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Edit Blog</h1>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              minLength={3}
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Image</label>
            {currentImage ? (
              <div className="relative inline-block">
                <img src={currentImage} alt="Preview" className="w-full max-w-md h-48 object-cover rounded-lg" />
                <button
                  type="button"
                  onClick={removeImage}
                  className="absolute top-2 right-2 p-1 bg-red-500 text-white rounded-full hover:bg-red-600"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            ) : (
              <div
                onClick={() => fileRef.current?.click()}
                className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center cursor-pointer hover:border-indigo-400"
              >
                <Upload className="mx-auto h-12 w-12 text-gray-400" />
                <p className="mt-2 text-sm text-gray-500">Click to upload image</p>
                <p className="text-xs text-gray-400">PNG, JPG up to 5MB</p>
              </div>
            )}
            <input ref={fileRef} type="file" accept="image/*" onChange={handleImageChange} className="hidden" />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Content</label>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              required
              minLength={10}
              rows={12}
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none"
            />
          </div>

          <div className="flex justify-end">
            <button
              type="submit"
              disabled={saving}
              className="bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 transition flex items-center space-x-2 disabled:opacity-50"
            >
              <Save size={18} />
              <span>{saving ? 'Saving...' : 'Save Changes'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
