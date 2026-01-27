import api from './api'
import type {
  ApiResponse,
  PaginatedResponse,
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
  Blog,
  CreateBlogRequest,
  UpdateBlogRequest,
} from '../types'

// Auth services
export const authService = {
  login: (data: LoginRequest) =>
    api.post<ApiResponse<AuthResponse>>('/auth/login', data),

  register: (data: RegisterRequest) =>
    api.post<ApiResponse<AuthResponse>>('/auth/register', data),

  me: () => api.get<ApiResponse<User>>('/auth/me'),
}

// User services
export const userService = {
  getAll: (page = 1, perPage = 15) =>
    api.get<PaginatedResponse<User>>(`/users?page=${page}&per_page=${perPage}`),

  getById: (id: number) => api.get<ApiResponse<User>>(`/users/${id}`),

  update: (id: number, data: Partial<User>) =>
    api.put<ApiResponse<User>>(`/users/${id}`, data),

  delete: (id: number) => api.delete<ApiResponse<null>>(`/users/${id}`),
}

// Blog services
export const blogService = {
  getAll: (page = 1, perPage = 15) =>
    api.get<PaginatedResponse<Blog>>(`/blogs?page=${page}&per_page=${perPage}`),

  getById: (id: number) => api.get<ApiResponse<Blog>>(`/blogs/${id}`),

  create: (data: CreateBlogRequest) =>
    api.post<ApiResponse<Blog>>('/blogs', data),

  update: (id: number, data: UpdateBlogRequest) =>
    api.put<ApiResponse<Blog>>(`/blogs/${id}`, data),

  delete: (id: number) => api.delete<ApiResponse<null>>(`/blogs/${id}`),

  getMyBlogs: (page = 1, perPage = 15) =>
    api.get<PaginatedResponse<Blog>>(`/blogs/my?page=${page}&per_page=${perPage}`),

  getUserBlogs: (userId: number, page = 1, perPage = 15) =>
    api.get<PaginatedResponse<Blog>>(`/users/${userId}/blogs?page=${page}&per_page=${perPage}`),
}
