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

// 用户类型
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

    if (data.code === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      localStorage.removeItem('permissions');
      goto('/auth/login');
      throw new Error('登录已过期，请重新登录');
    }

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
}

// 导出 API 客户端实例
export const api = new ApiClient(API_BASE_URL);

// ========== 认证 API ==========
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

// ========== 仪表盘 API ==========
export const dashboardApi = {
  overview: () => api.post<{
    books: { total: number; active: number; inactive: number };
    users: { total: number; active: number; new_today: number };
    stocks: { total: number; low_stock: number; out_of_stock: number };
    borrows: { borrowed: number; overdue: number };
    sales: { today_orders: number; today_amount: number };
  }>('/dashboard/overview'),

  borrowTrend: (days = 7) => api.post<{ date: string; count: number }[]>('/dashboard/borrow-trend', { days }),

  salesTrend: (days = 7) => api.post<{ date: string; orders: number; amount: number }[]>('/dashboard/sales-trend', { days }),

  categoryStats: () => api.post<{ category: string; count: number }[]>('/dashboard/category-stats'),

  topBorrowed: (limit = 10) => api.post<{ book_id: number; title: string; count: number }[]>('/dashboard/top-borrowed', { limit }),

  topSold: (limit = 10) => api.post<{ book_id: number; title: string; quantity: number; amount: number }[]>('/dashboard/top-sold', { limit }),

  activities: (limit = 20) => api.post<{ borrows: OperationLog[]; sales: OperationLog[] }>('/dashboard/activities', { limit }),

  alerts: () => api.post<{ type: string; module: string; title: string; message: string; detail: number }[]>('/dashboard/alerts'),

  search: (keyword: string, limit = 5) => api.post<{
    books: Book[];
    users: User[];
    suppliers: Supplier[];
    borrows: BorrowRecord[];
  }>('/search', { keyword, limit }),
};

// ========== 图书 API ==========
export const bookApi = {
  list: (params: { page?: number; page_size?: number; title?: string; author?: string; category?: string }) =>
    api.post<PageData<Book>>('/books/list', params),

  detail: (id: number) => api.post<Book>('/books/detail', { id }),

  create: (data: Partial<Book>) => api.post<Book>('/books/create', data),

  update: (data: Partial<Book>) => api.post<Book>('/books/update', data),

  delete: (id: number) => api.post('/books/delete', { id }),

  categories: () => api.post<string[]>('/books/categories'),
};

// ========== 库存 API ==========
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

// ========== 借阅 API ==========
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

// ========== 销售 API ==========
export const saleApi = {
  list: (params: { page?: number; page_size?: number; status?: string; start_date?: string; end_date?: string }) =>
    api.post<PageData<SaleOrder>>('/sales/list', params),

  create: (data: { customer_name?: string; customer_phone?: string; items: SaleItem[] }) =>
    api.post<SaleOrder>('/sales/create', data),

  cancel: (id: number) => api.post('/sales/cancel', { id }),

  stats: (start_date?: string, end_date?: string) =>
    api.post<{ total_orders: number; total_amount: number; total_quantity: number }>('/sales/stats', { start_date, end_date }),
};

// ========== 购物车 API ==========
export const cartApi = {
  list: () => api.post<Cart[]>('/cart/list'),

  add: (book_id: number, quantity: number) => api.post('/cart/add', { book_id, quantity }),

  remove: (id: number) => api.post('/cart/remove', { id }),

  clear: () => api.post('/cart/clear'),

  update: (id: number, quantity: number) => api.post('/cart/update', { id, quantity }),
};

// ========== 日志 API ==========
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

// ========== 批量操作 API ==========
export const batchApi = {
  importBooks: async (file: File): Promise<{ imported: number; errors: string[]; total: number }> => {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch(`${API_BASE_URL}/batch/import-books`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
      body: formData,
    });
    
    const result = await response.json();
    if (result.code !== 200) {
      throw new Error(result.message);
    }
    return result.data;
  },

  updateBooks: (ids: number[], update: Record<string, unknown>) =>
    api.post<{ updated: number }>('/batch/update-books', { ids, update }),

  deleteBooks: (ids: number[]) =>
    api.post<{ deleted: number }>('/batch/delete-books', { ids }),

  importUsers: async (file: File): Promise<{ imported: number; errors: string[]; total: number }> => {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch(`${API_BASE_URL}/batch/import-users`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
      body: formData,
    });
    
    const result = await response.json();
    if (result.code !== 200) {
      throw new Error(result.message);
    }
    return result.data;
  },

  updateUsers: (ids: number[], status: number) =>
    api.post<{ updated: number }>('/batch/update-users', { ids, status }),

  stockIn: (items: { book_id: number; quantity: number }[], reason?: string) =>
    api.post<{ success: number; errors: string[]; total: number }>('/batch/stock-in', { items, reason }),

  getTemplateUrl: (type: 'books' | 'users') =>
    `${API_BASE_URL}/batch/template?type=${type}`,
};

// ========== 系统配置 API ==========
export const configApi = {
  list: () => api.post<Record<string, string>>('/config/list'),

  get: (key: string) => api.post<SystemConfig>('/config/get', { key }),

  update: (key: string, value: string) =>
    api.post<SystemConfig>('/config/update', { key, value }),

  batchUpdate: (configs: Record<string, string>) =>
    api.post('/config/batch-update', configs),

  getBorrowRules: () => api.post<{
    borrow_days: number;
    max_borrow_books: number;
    max_renew_times: number;
    fine_per_day: number;
  }>('/config/borrow-rules'),

  updateBorrowRules: (rules: {
    borrow_days?: number;
    max_borrow_books?: number;
    max_renew_times?: number;
    fine_per_day?: number;
  }) => api.post('/config/update-borrow-rules', rules),

  getSiteInfo: () => api.post<{
    site_name: string;
    site_logo: string;
    site_footer: string;
    contact_email: string;
    contact_phone: string;
  }>('/config/site-info'),

  updateSiteInfo: (info: Record<string, string>) =>
    api.post('/config/update-site-info', info),

  reset: () => api.post<Record<string, string>>('/config/reset'),
};

// ========== 类型定义 ==========
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

export interface Supplier {
  id: number;
  name: string;
  code: string;
  contact: string;
  phone: string;
  email: string;
  address: string;
  status: number;
}

export interface SystemConfig {
  id: number;
  config_key: string;
  config_value: string;
  config_type: string;
  description: string;
}
