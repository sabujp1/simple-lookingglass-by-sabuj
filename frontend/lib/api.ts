import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/store/auth';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000,
});

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = useAuthStore.getState().token;
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
  user: User;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  role?: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface Router {
  id: string;
  name: string;
  hostname: string;
  ip_address: string;
  vendor: string;
  model: string;
  port: number;
  username: string;
  is_active: boolean;
  zone_id: string;
  zone_name?: string;
  last_check: string;
  created_at: string;
  updated_at: string;
}

export interface Zone {
  id: string;
  name: string;
  description: string;
  is_active: boolean;
  router_count: number;
  created_at: string;
  updated_at: string;
}

export interface QueryRequest {
  router_id: string;
  query_type: 'ping' | 'traceroute' | 'bgp' | 'bgp_as_path' | 'bgp_communities' | 'bgp_prefix';
  target: string;
  options?: Record<string, any>;
}

export interface QueryResult {
  id: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  router_id: string;
  router_name: string;
  query_type: string;
  target: string;
  output: string;
  error?: string;
  started_at: string;
  completed_at?: string;
  user_id: string;
}

export interface AuditLog {
  id: string;
  user_id: string;
  username: string;
  action: string;
  resource: string;
  resource_id?: string;
  details?: Record<string, any>;
  ip_address: string;
  created_at: string;
}

export interface PaginationParams {
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export const authApi = {
  login: (data: LoginRequest) => api.post<LoginResponse>('/auth/login', data),
  register: (data: RegisterRequest) => api.post<User>('/auth/register', data),
  refresh: (refreshToken: string) => 
    api.post<{ token: string }>('/auth/refresh', { refresh_token: refreshToken }),
  me: () => api.get<User>('/auth/me'),
};

export const routerApi = {
  list: (params?: PaginationParams & { zone_id?: string }) => 
    api.get<PaginatedResponse<Router>>('/routers', { params }),
  get: (id: string) => api.get<Router>(`/routers/${id}`),
  create: (data: Omit<Router, 'id' | 'created_at' | 'updated_at'>) => 
    api.post<Router>('/routers', data),
  update: (id: string, data: Partial<Router>) => 
    api.put<Router>(`/routers/${id}`, data),
  delete: (id: string) => api.delete(`/routers/${id}`),
  test: (id: string) => api.post<{ success: boolean; message: string }>(`/routers/${id}/test`),
  status: (id: string) => api.get<{ online: boolean; latency?: number }>(`/routers/${id}/status`),
};

export const zoneApi = {
  list: (params?: PaginationParams) => 
    api.get<PaginatedResponse<Zone>>('/zones', { params }),
  get: (id: string) => api.get<Zone>(`/zones/${id}`),
  create: (data: Omit<Zone, 'id' | 'created_at' | 'updated_at'>) => 
    api.post<Zone>('/zones', data),
  update: (id: string, data: Partial<Zone>) => 
    api.put<Zone>(`/zones/${id}`, data),
  delete: (id: string) => api.delete(`/zones/${id}`),
};

export const queryApi = {
  ping: (routerId: string, target: string, options?: { count?: number; size?: number }) =>
    api.post<QueryResult>('/queries/ping', { router_id: routerId, target, options }),
  traceroute: (routerId: string, target: string, options?: { max_hops?: number }) =>
    api.post<QueryResult>('/queries/traceroute', { router_id: routerId, target, options }),
  bgp: (routerId: string, query: string) =>
    api.post<QueryResult>('/queries/bgp', { router_id: routerId, target: query }),
  execute: (data: QueryRequest) => api.post<QueryResult>('/queries/execute', data),
  get: (id: string) => api.get<QueryResult>(`/queries/${id}`),
  list: (params?: PaginationParams & { router_id?: string; status?: string }) =>
    api.get<PaginatedResponse<QueryResult>>('/queries', { params }),
};

export const adminApi = {
  users: {
    list: (params?: PaginationParams) => 
      api.get<PaginatedResponse<User>>('/admin/users', { params }),
    get: (id: string) => api.get<User>(`/admin/users/${id}`),
    update: (id: string, data: Partial<User>) => 
      api.put<User>(`/admin/users/${id}`, data),
    delete: (id: string) => api.delete(`/admin/users/${id}`),
  },
  logs: (params?: PaginationParams & { user_id?: string; action?: string }) =>
    api.get<PaginatedResponse<AuditLog>>('/admin/logs', { params }),
  stats: () => api.get<{
    total_routers: number;
    active_routers: number;
    total_zones: number;
    total_users: number;
    queries_today: number;
    queries_by_type: Record<string, number>;
  }>('/admin/stats'),
};

export default api;