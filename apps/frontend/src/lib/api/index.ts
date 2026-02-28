import { writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// 响应类型
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
}

// 分页数据类型
export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// 用户状态
export interface User {
  id: number;
  username: string;
  email: string;
  phone?: string;
  real_name?: string;
  avatar?: string;
  status: number;
  roles: Role[];
  created_at: string;
}

export interface Role {
  id: number;
  name: string;
  code: string;
  permissions: Permission[];
}

export interface Permission {
  id: number;
  name: string;
  code: string;
  resource: string;
  action: string;
}

// 认证状态
export interface AuthState {
  user: User | null;
  token: string | null;
  permissions: string[];
  isAuthenticated: boolean;
}

// 创建认证状态存储
export const authStore = writable<AuthState>({
  user: null,
  token: null,
  permissions: [],
  isAuthenticated: false
});

// 从 localStorage 恢复认证状态
export function initAuth() {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    const permissionsStr = localStorage.getItem('permissions');

    if (token && userStr) {
      try {
        const user = JSON.parse(userStr);
        const permissions = permissionsStr ? JSON.parse(permissionsStr) : [];
        authStore.set({ user, token, permissions, isAuthenticated: true });
      } catch {
        logout();
      }
    }
  }
}

// 保存认证状态
export function saveAuth(token: string, user: User, permissions: string[]) {
  if (typeof window !== 'undefined') {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    localStorage.setItem('permissions', JSON.stringify(permissions));
  }
  authStore.set({ user, token, permissions, isAuthenticated: true });
}

// 清除认证状态
export function clearAuth() {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('permissions');
  }
  authStore.set({ user: null, token: null, permissions: [], isAuthenticated: false });
}

// 登出
export async function logout() {
  try {
    await api.post('/logout');
  } catch {
    // ignore
  }
  clearAuth();
  goto('/auth/login');
}

// API 请求类
class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    const data: ApiResponse<T> = await response.json();

    // 处理认证错误
    if (data.code === 401) {
      clearAuth();
      goto('/auth/login');
      throw new Error('登录已过期，请重新登录');
    }

    // 处理其他错误
    if (data.code !== 200) {
      throw new Error(data.message || '请求失败');
    }

    return data;
  }

  async post<T = unknown>(path: string, body?: unknown): Promise<T> {
    try {
      const response = await fetch(`${this.baseUrl}${path}`, {
        method: 'POST',
        headers: this.getHeaders(),
        body: body ? JSON.stringify(body) : undefined,
      });

      const result = await this.handleResponse<T>(response);
      return result.data as T;
    } catch (error) {
      const message = error instanceof Error ? error.message : '网络错误';
      toast.error(message);
      throw error;
    }
  }

  // 带进度的文件上传
  async upload<T = unknown>(path: string, file: File, onProgress?: (percent: number) => void): Promise<T> {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();
      const formData = new FormData();
      formData.append('file', file);

      xhr.upload.addEventListener('progress', (e) => {
        if (e.lengthComputable && onProgress) {
          const percent = Math.round((e.loaded / e.total) * 100);
          onProgress(percent);
        }
      });

      xhr.addEventListener('load', async () => {
        if (xhr.status === 200) {
          try {
            const result = JSON.parse(xhr.responseText);
            if (result.code === 200) {
              resolve(result.data);
            } else {
              reject(new Error(result.message));
            }
          } catch {
            reject(new Error('解析响应失败'));
          }
        } else {
          reject(new Error('上传失败'));
        }
      });

      xhr.addEventListener('error', () => reject(new Error('网络错误')));
      xhr.addEventListener('abort', () => reject(new Error('上传已取消')));

      xhr.open('POST', `${this.baseUrl}${path}`);

      // 添加认证头
      const token = localStorage.getItem('token');
      if (token) {
        xhr.setRequestHeader('Authorization', `Bearer ${token}`);
      }

      xhr.send(formData);
    });
  }
}

// 导出 API 客户端实例
export const api = new ApiClient(API_BASE_URL);

// 便捷方法
export const authApi = {
  login: (username: string, password: string) =>
    api.post<{ token: string; user: User; permissions: string[] }>('/login', { username, password }),

  register: (data: { username: string; password: string; email: string; real_name?: string }) =>
    api.post<User>('/register', data),

  logout: () => api.post('/logout'),

  profile: () => api.post<{ user: User; permissions: string[] }>('/profile'),

  changePassword: (old_password: string, new_password: string) =>
    api.post('/users/change-password', { old_password, new_password }),
};

export const bookApi = {
  list: (params: { page?: number; page_size?: number; title?: string; author?: string; category?: string }) =>
    api.post<PageData<Book>>('/books/list', params),

  detail: (id: number) => api.post<Book>('/books/detail', { id }),

  create: (data: Partial<Book>) => api.post<Book>('/books/create', data),

  update: (data: Partial<Book>) => api.post<Book>('/books/update', data),

  delete: (id: number) => api.post('/books/delete', { id }),

  categories: () => api.post<string[]>('/books/categories'),
};

export const stockApi = {
  list: (params: { page?: number; page_size?: number; book_title?: string }) =>
    api.post<PageData<Stock>>('/stocks/list', params),

  detail: (id: number) => api.post<Stock>('/stocks/detail', { id }),

  in: (book_id: number, quantity: number, reason?: string) =>
    api.post<Stock>('/stocks/in', { book_id, quantity, reason }),

  out: (book_id: number, quantity: number, reason: string) =>
    api.post<Stock>('/stocks/out', { book_id, quantity, reason }),

  records: (params: { page?: number; page_size?: number; book_id?: number; type?: string }) =>
    api.post<PageData<StockRecord>>('/stocks/records', params),

  low: (threshold?: number) => api.post<Stock[]>('/stocks/low', { threshold }),
};

export const borrowApi = {
  list: (params: { page?: number; page_size?: number; status?: string; user_id?: number }) =>
    api.post<PageData<BorrowRecord>>('/borrows/list', params),

  borrow: (book_id: number, days?: number) => api.post('/borrows/borrow', { book_id, days }),

  return: (id: string) => api.post('/borrows/return', { id }),

  renew: (id: string) => api.post<{ new_due_date: string }>('/borrows/renew', { id }),

  overdue: (params: { page?: number; page_size?: number }) =>
    api.post<PageData<BorrowRecord>>('/borrows/overdue', params),

  stats: () => api.post<{ current_borrowed: number; total_borrowed: number; overdue_count: number }>('/borrows/stats'),

  payFine: (id: string) => api.post('/borrows/pay-fine', { id }),
};

export const saleApi = {
  list: (params: { page?: number; page_size?: number; status?: string; start_date?: string; end_date?: string }) =>
    api.post<PageData<SaleOrder>>('/sales/list', params),

  create: (data: { customer_name?: string; customer_phone?: string; items: SaleItem[] }) =>
    api.post<SaleOrder>('/sales/create', data),

  cancel: (id: number) => api.post('/sales/cancel', { id }),

  stats: (start_date?: string, end_date?: string) =>
    api.post<{ total_orders: number; total_amount: number; total_quantity: number }>('/sales/stats', { start_date, end_date }),
};

export const cartApi = {
  list: () => api.post<Cart[]>('/cart/list'),

  add: (book_id: number, quantity: number) => api.post('/cart/add', { book_id, quantity }),

  remove: (id: number) => api.post('/cart/remove', { id }),

  clear: () => api.post('/cart/clear'),

  update: (id: number, quantity: number) => api.post('/cart/update', { id, quantity }),
};

export const logApi = {
  list: (params: { page?: number; page_size?: number; user_id?: number; module?: string; action?: string }) =>
    api.post<PageData<OperationLog>>('/logs/list', params),

  modules: () => api.post<string[]>('/logs/modules'),

  actions: () => api.post<string[]>('/logs/actions'),

  dashboardStats: () => api.post<{
    today_count: number;
    week_count: number;
    month_count: number;
    active_users: number;
    action_stats: { action: string; count: number }[];
  }>('/logs/dashboard-stats'),
};

// 类型定义
export interface Book {
  id: number;
  isbn: string;
  title: string;
  author: string;
  publisher: string;
  category: string;
  price: number;
  cover_image?: string;
  description?: string;
  status: number;
  stock?: Stock;
}

export interface Stock {
  id: number;
  book_id: number;
  book?: Book;
  total_quantity: number;
  available_quantity: number;
  borrowed_quantity: number;
  sold_quantity: number;
  location?: string;
}

export interface StockRecord {
  id: number;
  book_id: number;
  book?: Book;
  type: string;
  quantity: number;
  reason?: string;
  created_at: string;
}

export interface BorrowRecord {
  id: string;
  book_id: number;
  book_title: string;
  user_id: number;
  user_name: string;
  borrow_date: string;
  due_date: string;
  return_date?: string;
  status: string;
  renew_count: number;
  fine: number;
  fine_paid: boolean;
}

export interface SaleOrder {
  id: number;
  order_no: string;
  customer_name?: string;
  customer_phone?: string;
  total_amount: number;
  pay_amount: number;
  status: string;
  created_at: string;
  items?: SaleOrderItem[];
}

export interface SaleOrderItem {
  id: number;
  book_id: number;
  book?: Book;
  quantity: number;
  unit_price: number;
  total_price: number;
}

export interface SaleItem {
  book_id: number;
  quantity: number;
  unit_price: number;
  discount?: number;
}

export interface Cart {
  id: number;
  book_id: number;
  book?: Book;
  quantity: number;
}

export interface OperationLog {
  id: number;
  user_id?: number;
  username?: string;
  module: string;
  action: string;
  ip?: string;
  content?: string;
  created_at: string;
}
