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
  UpdateUserRequest,
  UpdateRoleRequest,
} from '../types'

export const authService = {
  login: (data: LoginRequest) =>
    api.post<ApiResponse<AuthResponse>>('/auth/login', data),

  register: (data: RegisterRequest) =>
    api.post<ApiResponse<AuthResponse>>('/auth/register', data),

  me: () => api.get<ApiResponse<User>>('/auth/me'),
}

export const userService = {
  getAll: (page = 1, perPage = 15) =>
    api.get<PaginatedResponse<User>>(`/users?page=${page}&per_page=${perPage}`),

  getById: (id: number) => api.get<ApiResponse<User>>(`/users/${id}`),

  update: (id: number, data: UpdateUserRequest) =>
    api.put<ApiResponse<User>>(`/users/${id}`, data),

  updateRole: (id: number, data: UpdateRoleRequest) =>
    api.patch<ApiResponse<User>>(`/users/${id}/role`, data),

  delete: (id: number) => api.delete<ApiResponse<null>>(`/users/${id}`),

  getTrashed: (page = 1, perPage = 15) =>
    api.get<PaginatedResponse<User>>(`/admin/users/trashed?page=${page}&per_page=${perPage}`),

  restore: (id: number) =>
    api.post<ApiResponse<null>>(`/admin/users/${id}/restore`),

  forceDelete: (id: number) =>
    api.delete<ApiResponse<null>>(`/admin/users/${id}/force`),
}

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

  getTrashed: (page = 1, perPage = 15) =>
    api.get<PaginatedResponse<Blog>>(`/admin/blogs/trashed?page=${page}&per_page=${perPage}`),

  restore: (id: number) =>
    api.post<ApiResponse<null>>(`/admin/blogs/${id}/restore`),

  forceDelete: (id: number) =>
    api.delete<ApiResponse<null>>(`/admin/blogs/${id}/force`),

  uploadImage: (id: number, file: File) => {
    const formData = new FormData()
    formData.append('image', file)
    return api.post<ApiResponse<{ image: string }>>(`/blogs/${id}/image`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  deleteImage: (id: number) =>
    api.delete<ApiResponse<null>>(`/blogs/${id}/image`),
}
