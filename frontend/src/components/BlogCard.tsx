import { Link } from 'react-router-dom'
import type { Blog } from '../types'
import { Calendar, User } from 'lucide-react'

interface BlogCardProps {
  blog: Blog
}

export default function BlogCard({ blog }: BlogCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  return (
    <article className="bg-white rounded-xl shadow-sm border hover:shadow-md transition p-6">
      <Link to={`/blogs/${blog.id}`}>
        <h2 className="text-xl font-semibold text-gray-900 hover:text-indigo-600 transition mb-3">
          {blog.title}
        </h2>
      </Link>
      <p className="text-gray-600 line-clamp-3 mb-4">
        {blog.content}
      </p>
      <div className="flex items-center justify-between text-sm text-gray-500">
        <div className="flex items-center space-x-4">
          {blog.user && (
            <span className="flex items-center space-x-1">
              <User size={14} />
              <span>{blog.user.name}</span>
            </span>
          )}
          <span className="flex items-center space-x-1">
            <Calendar size={14} />
            <span>{formatDate(blog.created_at)}</span>
          </span>
        </div>
        <Link
          to={`/blogs/${blog.id}`}
          className="text-indigo-600 hover:text-indigo-700 font-medium"
        >
          Read more →
        </Link>
      </div>
    </article>
  )
}