import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import { userService, blogService } from '../services'
import type { User, Blog } from '../types'
import ConfirmModal from '../components/ConfirmModal'
import {
  PlusCircle,
  FileText,
  Users,
  Trash2,
  Edit,
  Shield,
  UserIcon,
  RotateCcw,
  XCircle,
} from 'lucide-react'

type TabType = 'blogs' | 'users' | 'trashed-blogs' | 'trashed-users'

interface ModalState {
  isOpen: boolean
  type: 'delete' | 'restore' | 'force-delete'
  targetType: 'blog' | 'user'
  targetId: number
  targetName: string
}

export default function DashboardPage() {
  const { user: currentUser } = useAuthStore()
  const [activeTab, setActiveTab] = useState<TabType>('blogs')
  const [blogs, setBlogs] = useState<Blog[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [trashedBlogs, setTrashedBlogs] = useState<Blog[]>([])
  const [trashedUsers, setTrashedUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [modal, setModal] = useState<ModalState>({
    isOpen: false,
    type: 'delete',
    targetType: 'blog',
    targetId: 0,
    targetName: '',
  })

  const isAdmin = currentUser?.role === 'admin'

  useEffect(() => {
    loadData()
  }, [activeTab, isAdmin])

  const loadData = async () => {
    setLoading(true)
    try {
      if (activeTab === 'blogs') {
        if (isAdmin) {
          const res = await blogService.getAll(1, 50)
          setBlogs(res.data.data || [])
        } else {
          const res = await blogService.getMyBlogs(1, 50)
          setBlogs(res.data.data || [])
        }
      } else if (activeTab === 'users' && isAdmin) {
        const res = await userService.getAll(1, 50)
        setUsers(res.data.data || [])
      } else if (activeTab === 'trashed-blogs' && isAdmin) {
        const res = await blogService.getTrashed(1, 50)
        setTrashedBlogs(res.data.data || [])
      } else if (activeTab === 'trashed-users' && isAdmin) {
        const res = await userService.getTrashed(1, 50)
        setTrashedUsers(res.data.data || [])
      }
    } catch (err) {
      console.error('Failed to load data:', err)
    } finally {
      setLoading(false)
    }
  }

  const openModal = (
    type: ModalState['type'],
    targetType: ModalState['targetType'],
    targetId: number,
    targetName: string
  ) => {
    setModal({ isOpen: true, type, targetType, targetId, targetName })
  }

  const closeModal = () => {
    setModal(prev => ({ ...prev, isOpen: false }))
  }

  const handleConfirm = async () => {
    const { type, targetType, targetId } = modal
    try {
      if (targetType === 'blog') {
        if (type === 'delete') {
          await blogService.delete(targetId)
          setBlogs(blogs.filter(b => b.id !== targetId))
        } else if (type === 'restore') {
          await blogService.restore(targetId)
          setTrashedBlogs(trashedBlogs.filter(b => b.id !== targetId))
        } else if (type === 'force-delete') {
          await blogService.forceDelete(targetId)
          setTrashedBlogs(trashedBlogs.filter(b => b.id !== targetId))
        }
      } else {
        if (type === 'delete') {
          await userService.delete(targetId)
          setUsers(users.filter(u => u.id !== targetId))
        } else if (type === 'restore') {
          await userService.restore(targetId)
          setTrashedUsers(trashedUsers.filter(u => u.id !== targetId))
        } else if (type === 'force-delete') {
          await userService.forceDelete(targetId)
          setTrashedUsers(trashedUsers.filter(u => u.id !== targetId))
        }
      }
    } catch (err) {
      console.error('Action failed:', err)
    } finally {
      closeModal()
    }
  }

  const handleToggleRole = async (id: number, currentRole: string) => {
    const newRole = currentRole === 'admin' ? 'author' : 'admin'
    try {
      await userService.updateRole(id, { role: newRole as 'admin' | 'author' })
      setUsers(
        users.map(u =>
          u.id === id ? { ...u, role: newRole as 'admin' | 'author' } : u
        )
      )
    } catch (err) {
      console.error('Failed to update role:', err)
    }
  }

  const getModalMessage = () => {
    const { type, targetName } = modal
    if (type === 'delete') {
      return `Are you sure you want to delete "${targetName}"? It will be moved to trash.`
    } else if (type === 'restore') {
      return `Are you sure you want to restore "${targetName}"?`
    } else {
      return `Are you sure you want to permanently delete "${targetName}"? This action cannot be undone.`
    }
  }

  const getModalTitle = () => {
    const { type, targetType } = modal
    const itemType = targetType === 'blog' ? 'Blog' : 'User'
    if (type === 'delete') return `Delete ${itemType}`
    if (type === 'restore') return `Restore ${itemType}`
    return `Permanently Delete ${itemType}`
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-500 mt-1">
          Welcome back, {currentUser?.name}
          {isAdmin && (
            <span className="ml-2 px-2 py-0.5 bg-indigo-100 text-indigo-700 text-sm rounded-full">
              Admin
            </span>
          )}
        </p>
      </div>

      <div className="flex flex-wrap gap-3 mb-6">
        <button
          onClick={() => setActiveTab('blogs')}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition ${
            activeTab === 'blogs'
              ? 'bg-indigo-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
        >
          <FileText size={18} />
          {isAdmin ? 'All Blogs' : 'My Blogs'}
        </button>
        {isAdmin && (
          <>
            <button
              onClick={() => setActiveTab('users')}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition ${
                activeTab === 'users'
                  ? 'bg-indigo-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              <Users size={18} />
              Users
            </button>
            <button
              onClick={() => setActiveTab('trashed-blogs')}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition ${
                activeTab === 'trashed-blogs'
                  ? 'bg-red-600 text-white'
                  : 'bg-red-50 text-red-700 hover:bg-red-100'
              }`}
            >
              <Trash2 size={18} />
              Trashed Blogs
            </button>
            <button
              onClick={() => setActiveTab('trashed-users')}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition ${
                activeTab === 'trashed-users'
                  ? 'bg-red-600 text-white'
                  : 'bg-red-50 text-red-700 hover:bg-red-100'
              }`}
            >
              <Trash2 size={18} />
              Trashed Users
            </button>
          </>
        )}
        <Link
          to="/blogs/create"
          className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg font-medium hover:bg-green-700 transition ml-auto"
        >
          <PlusCircle size={18} />
          New Blog
        </Link>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : (
        <>
          {activeTab === 'blogs' && (
            <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Title
                    </th>
                    {isAdmin && (
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Author
                      </th>
                    )}
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Date
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {blogs.map(blog => (
                    <tr key={blog.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-3">
                          {blog.image && (
                            <img
                              src={blog.image}
                              alt=""
                              className="w-10 h-10 rounded object-cover"
                            />
                          )}
                          <Link
                            to={`/blogs/${blog.id}`}
                            className="text-indigo-600 hover:underline font-medium"
                          >
                            {blog.title}
                          </Link>
                        </div>
                      </td>
                      {isAdmin && (
                        <td className="px-6 py-4 text-gray-500 text-sm">
                          {blog.user?.name || 'Unknown'}
                        </td>
                      )}
                      <td className="px-6 py-4 text-gray-500 text-sm">
                        {new Date(blog.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Link
                            to={`/blogs/${blog.id}/edit`}
                            className="p-2 text-gray-500 hover:text-indigo-600 transition"
                          >
                            <Edit size={16} />
                          </Link>
                          <button
                            onClick={() =>
                              openModal('delete', 'blog', blog.id, blog.title)
                            }
                            className="p-2 text-gray-500 hover:text-red-600 transition"
                          >
                            <Trash2 size={16} />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                  {blogs.length === 0 && (
                    <tr>
                      <td
                        colSpan={isAdmin ? 4 : 3}
                        className="px-6 py-12 text-center text-gray-500"
                      >
                        No blogs found
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}

          {activeTab === 'users' && isAdmin && (
            <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Email
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Role
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Joined
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {users.map(user => (
                    <tr key={user.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 font-medium text-gray-900">
                        {user.name}
                      </td>
                      <td className="px-6 py-4 text-gray-500">{user.email}</td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-full ${
                            user.role === 'admin'
                              ? 'bg-indigo-100 text-indigo-700'
                              : 'bg-gray-100 text-gray-700'
                          }`}
                        >
                          {user.role === 'admin' ? (
                            <Shield size={12} />
                          ) : (
                            <UserIcon size={12} />
                          )}
                          {user.role}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-gray-500 text-sm">
                        {new Date(user.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 text-right">
                        {user.id !== currentUser?.id && (
                          <div className="flex items-center justify-end gap-2">
                            <button
                              onClick={() =>
                                handleToggleRole(user.id, user.role)
                              }
                              className="p-2 text-gray-500 hover:text-indigo-600 transition"
                              title={`Make ${user.role === 'admin' ? 'Author' : 'Admin'}`}
                            >
                              <Shield size={16} />
                            </button>
                            <button
                              onClick={() =>
                                openModal('delete', 'user', user.id, user.name)
                              }
                              className="p-2 text-gray-500 hover:text-red-600 transition"
                            >
                              <Trash2 size={16} />
                            </button>
                          </div>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {activeTab === 'trashed-blogs' && isAdmin && (
            <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Title
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Author
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Deleted At
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {trashedBlogs.map(blog => (
                    <tr key={blog.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 font-medium text-gray-900">
                        {blog.title}
                      </td>
                      <td className="px-6 py-4 text-gray-500 text-sm">
                        {blog.user?.name || 'Unknown'}
                      </td>
                      <td className="px-6 py-4 text-gray-500 text-sm">
                        {blog.deleted_at
                          ? new Date(blog.deleted_at).toLocaleDateString()
                          : '-'}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            onClick={() =>
                              openModal('restore', 'blog', blog.id, blog.title)
                            }
                            className="p-2 text-gray-500 hover:text-green-600 transition"
                            title="Restore"
                          >
                            <RotateCcw size={16} />
                          </button>
                          <button
                            onClick={() =>
                              openModal(
                                'force-delete',
                                'blog',
                                blog.id,
                                blog.title
                              )
                            }
                            className="p-2 text-gray-500 hover:text-red-600 transition"
                            title="Delete Permanently"
                          >
                            <XCircle size={16} />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                  {trashedBlogs.length === 0 && (
                    <tr>
                      <td
                        colSpan={4}
                        className="px-6 py-12 text-center text-gray-500"
                      >
                        No trashed blogs
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}

          {activeTab === 'trashed-users' && isAdmin && (
            <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Email
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Role
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Deleted At
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {trashedUsers.map(user => (
                    <tr key={user.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 font-medium text-gray-900">
                        {user.name}
                      </td>
                      <td className="px-6 py-4 text-gray-500">{user.email}</td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-full ${
                            user.role === 'admin'
                              ? 'bg-indigo-100 text-indigo-700'
                              : 'bg-gray-100 text-gray-700'
                          }`}
                        >
                          {user.role === 'admin' ? (
                            <Shield size={12} />
                          ) : (
                            <UserIcon size={12} />
                          )}
                          {user.role}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-gray-500 text-sm">
                        {user.deleted_at
                          ? new Date(user.deleted_at).toLocaleDateString()
                          : '-'}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            onClick={() =>
                              openModal('restore', 'user', user.id, user.name)
                            }
                            className="p-2 text-gray-500 hover:text-green-600 transition"
                            title="Restore"
                          >
                            <RotateCcw size={16} />
                          </button>
                          <button
                            onClick={() =>
                              openModal(
                                'force-delete',
                                'user',
                                user.id,
                                user.name
                              )
                            }
                            className="p-2 text-gray-500 hover:text-red-600 transition"
                            title="Delete Permanently"
                          >
                            <XCircle size={16} />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                  {trashedUsers.length === 0 && (
                    <tr>
                      <td
                        colSpan={5}
                        className="px-6 py-12 text-center text-gray-500"
                      >
                        No trashed users
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}

      <ConfirmModal
        isOpen={modal.isOpen}
        title={getModalTitle()}
        message={getModalMessage()}
        type={modal.type}
        onConfirm={handleConfirm}
        onCancel={closeModal}
      />
    </div>
  )
}
