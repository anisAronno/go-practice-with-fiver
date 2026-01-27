export interface User {
  id: number
  name: string
  email: string
  created_at: string
  updated_at: string
}

export interface Blog {
  id: number
  title: string
  content: string
  user_id: number
  user?: User
  created_at: string
  updated_at: string
}

export interface ApiResponse<T> {
  success: boolean
  message?: string
  data: T
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  meta: {
    current_page: number
    per_page: number
    total: number
    total_pages: number
  }
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface AuthResponse {
  user: User
  token: string
}

export interface CreateBlogRequest {
  title: string
  content: string
}

export interface UpdateBlogRequest {
  title?: string
  content?: string
}
