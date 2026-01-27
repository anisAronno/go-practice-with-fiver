import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import { blogService, userService } from '../services'
import type { Blog, User } from '../types'
import ConfirmModal from '../components/ConfirmModal'
import { FileText, Users, PlusCircle, Trash2, Edit, RotateCcw, XCircle, ChevronLeft, ChevronRight } from 'lucide-react'

type Tab = 'blogs' | 'users'
type View = 'active' | 'trashed'

interface Modal {
  isOpen: boolean
  type: 'delete' | 'restore' | 'force-delete'
  id: number
  name: string
  entity: 'blog' | 'user'
}

const API_URL = import.meta.env.VITE_API_URL || ''

export default function DashboardPage() {
  const { user } = useAuthStore()
  const isAdmin = user?.role === 'admin'
  
  const [tab, setTab] = useState<Tab>('blogs')
  const [view, setView] = useState<View>('active')
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [modal, setModal] = useState<Modal>({ isOpen: false, type: 'delete', id: 0, name: '', entity: 'blog' })

  useEffect(() => {
    setPage(1)
  }, [tab, view])

  useEffect(() => {
    loadData()
  }, [tab, view, page, isAdmin])

  const loadData = async () => {
    setLoading(true)
    try {
      if (tab === 'blogs') {
        if (view === 'active') {
          const res = isAdmin ? await blogService.getAll(page) : await blogService.getMyBlogs(page)
          setBlogs(res.data.data || [])
          setTotalPages(res.data.meta?.total_pages || 1)
        } else if (isAdmin) {
          const res = await blogService.getTrashed(page)
          setBlogs(res.data.data || [])
          setTotalPages(res.data.meta?.total_pages || 1)
        }
      } else if (isAdmin) {
        if (view === 'active') {
          const res = await userService.getAll(page)
          setUsers(res.data.data || [])
          setTotalPages(res.data.meta?.total_pages || 1)
        } else {
          const res = await userService.getTrashed(page)
          setUsers(res.data.data || [])
          setTotalPages(res.data.meta?.total_pages || 1)
        }
      }
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const openModal = (type: Modal['type'], id: number, name: string, entity: Modal['entity']) => {
    setModal({ isOpen: true, type, id, name, entity })
  }

  const closeModal = () => setModal(prev => ({ ...prev, isOpen: false }))

  const handleConfirm = async () => {
    try {
      if (modal.entity === 'blog') {
        if (modal.type === 'delete') await blogService.delete(modal.id)
        else if (modal.type === 'restore') await blogService.restore(modal.id)
        else await blogService.forceDelete(modal.id)
        setBlogs(blogs.filter(b => b.id !== modal.id))
      } else {
        if (modal.type === 'delete') await userService.delete(modal.id)
        else if (modal.type === 'restore') await userService.restore(modal.id)
        else await userService.forceDelete(modal.id)
        setUsers(users.filter(u => u.id !== modal.id))
      }
    } catch (err) {
      console.error(err)
    } finally {
      closeModal()
    }
  }

  const getImageUrl = (path: string) => {
    if (!path) return ''
    if (path.startsWith('http')) return path
    return `${API_URL}${path}`
  }

  return (
    <div className="flex gap-6 min-h-[600px]">
      {/* Sidebar */}
      <div className="w-56 shrink-0">
        <div className="bg-white rounded-xl border p-4 sticky top-4">
          <h2 className="font-semibold text-gray-900 mb-4">Dashboard</h2>
          <nav className="space-y-1">
            <button
              onClick={() => { setTab('blogs'); setView('active') }}
              className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-left ${tab === 'blogs' ? 'bg-indigo-50 text-indigo-700' : 'text-gray-600 hover:bg-gray-50'}`}
            >
              <FileText size={18} />
              <span>{isAdmin ? 'All Blogs' : 'My Blogs'}</span>
            </button>
            {isAdmin && (
              <button
                onClick={() => { setTab('users'); setView('active') }}
                className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-left ${tab === 'users' ? 'bg-indigo-50 text-indigo-700' : 'text-gray-600 hover:bg-gray-50'}`}
              >
                <Users size={18} />
                <span>Users</span>
              </button>
            )}
          </nav>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{tab === 'blogs' ? 'Blogs' : 'Users'}</h1>
            <p className="text-gray-500 text-sm">{isAdmin ? (tab === 'blogs' ? 'Manage all blogs' : 'Manage users') : 'My blogs'}</p>
          </div>
          <div className="flex items-center gap-3">
            {isAdmin && (
              <select value={view} onChange={e => setView(e.target.value as View)} className="px-3 py-2 border rounded-lg text-sm">
                <option value="active">Active</option>
                <option value="trashed">Trashed</option>
              </select>
            )}
            {tab === 'blogs' && (
              <Link to="/blogs/create" className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm hover:bg-indigo-700">
                <PlusCircle size={16} />
                New Blog
              </Link>
            )}
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center py-12">
            <div className="animate-spin w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full" />
          </div>
        ) : (
          <>
            <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    {tab === 'blogs' ? (
                      <>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Title</th>
                        {isAdmin && <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Author</th>}
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                        <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
                      </>
                    ) : (
                      <>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Role</th>
                        <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
                      </>
                    )}
                  </tr>
                </thead>
                <tbody className="divide-y">
                  {tab === 'blogs' ? blogs.map(blog => (
                    <tr key={blog.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-3">
                          {blog.image && <img src={getImageUrl(blog.image)} alt="" className="w-10 h-10 rounded object-cover" />}
                          <span className="font-medium text-gray-900">{blog.title}</span>
                        </div>
                      </td>
                      {isAdmin && <td className="px-4 py-3 text-gray-500 text-sm">{blog.user?.name || '-'}</td>}
                      <td className="px-4 py-3 text-gray-500 text-sm">{new Date(blog.created_at).toLocaleDateString()}</td>
                      <td className="px-4 py-3 text-right">
                        <div className="flex justify-end gap-1">
                          {view === 'active' ? (
                            <>
                              <Link to={`/blogs/${blog.id}/edit`} className="p-2 text-gray-500 hover:text-indigo-600"><Edit size={16} /></Link>
                              <button onClick={() => openModal('delete', blog.id, blog.title, 'blog')} className="p-2 text-gray-500 hover:text-red-600"><Trash2 size={16} /></button>
                            </>
                          ) : (
                            <>
                              <button onClick={() => openModal('restore', blog.id, blog.title, 'blog')} className="p-2 text-gray-500 hover:text-green-600"><RotateCcw size={16} /></button>
                              <button onClick={() => openModal('force-delete', blog.id, blog.title, 'blog')} className="p-2 text-gray-500 hover:text-red-600"><XCircle size={16} /></button>
                            </>
                          )}
                        </div>
                      </td>
                    </tr>
                  )) : users.map(u => (
                    <tr key={u.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3 font-medium text-gray-900">{u.name}</td>
                      <td className="px-4 py-3 text-gray-500 text-sm">{u.email}</td>
                      <td className="px-4 py-3">
                        <span className={`px-2 py-1 text-xs rounded-full ${u.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'}`}>{u.role}</span>
                      </td>
                      <td className="px-4 py-3 text-right">
                        <div className="flex justify-end gap-1">
                          {view === 'active' ? (
                            <button onClick={() => openModal('delete', u.id, u.name, 'user')} className="p-2 text-gray-500 hover:text-red-600"><Trash2 size={16} /></button>
                          ) : (
                            <>
                              <button onClick={() => openModal('restore', u.id, u.name, 'user')} className="p-2 text-gray-500 hover:text-green-600"><RotateCcw size={16} /></button>
                              <button onClick={() => openModal('force-delete', u.id, u.name, 'user')} className="p-2 text-gray-500 hover:text-red-600"><XCircle size={16} /></button>
                            </>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                  {((tab === 'blogs' && blogs.length === 0) || (tab === 'users' && users.length === 0)) && (
                    <tr><td colSpan={4} className="px-4 py-12 text-center text-gray-500">No {tab} found</td></tr>
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center gap-4 mt-6">
                <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1} className="flex items-center gap-1 px-3 py-2 border rounded-lg disabled:opacity-50 hover:bg-gray-50">
                  <ChevronLeft size={16} /> Prev
                </button>
                <span className="text-sm text-gray-600">Page {page} of {totalPages}</span>
                <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="flex items-center gap-1 px-3 py-2 border rounded-lg disabled:opacity-50 hover:bg-gray-50">
                  Next <ChevronRight size={16} />
                </button>
              </div>
            )}
          </>
        )}
      </div>

      <ConfirmModal
        isOpen={modal.isOpen}
        title={modal.type === 'delete' ? 'Delete' : modal.type === 'restore' ? 'Restore' : 'Permanently Delete'}
        message={`${modal.type === 'restore' ? 'Restore' : 'Delete'} "${modal.name}"?`}
        type={modal.type}
        onConfirm={handleConfirm}
        onCancel={closeModal}
      />
    </div>
  )
}
