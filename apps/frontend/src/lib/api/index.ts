import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import type { ApiResponse, PageData } from '$lib/types';

const API_BASE = '/api';

class ApiClient {
  private getToken(): string | null {
    if (browser) {
      return localStorage.getItem('token');
    }
    return null;
  }

  private async request<T>(
    endpoint: string,
    body?: unknown
  ): Promise<ApiResponse<T>> {
    const token = this.getToken();
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json'
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      method: 'POST',
      headers,
      body: body ? JSON.stringify(body) : undefined
    });

    const data = await response.json();

    // 处理未授权情况
    if (data.code === 401) {
      if (browser) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        localStorage.removeItem('permissions');
        goto('/auth/login');
      }
    }

    return data;
  }

  async post<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, body);
  }
}

export const api = new ApiClient();

// 认证 API
export const authApi = {
  login: (username: string, password: string) =>
    api.post<{ token: string; user: unknown; permissions: string[] }>('/login', { username, password }),
  
  register: (data: { username: string; password: string; email: string; real_name?: string }) =>
    api.post('/register', data)
};

// 用户 API
export const userApi = {
  getList: (params?: { page?: number; page_size?: number; username?: string }) =>
    api.post<PageData<unknown>>('/users/list', params),
  
  get: (id: number) =>
    api.post<unknown>(`/users/detail`, { id }),
  
  create: (data: unknown) =>
    api.post('/users/create', data),
  
  update: (id: number, data: unknown) =>
    api.post(`/users/update`, { id, ...data }),
  
  delete: (id: number) =>
    api.post('/users/delete', { id }),
  
  assignRoles: (id: number, roleIds: number[]) =>
    api.post(`/users/assign-roles`, { id, role_ids: roleIds }),
  
  changePassword: (oldPassword: string, newPassword: string) =>
    api.post('/users/change-password', { old_password: oldPassword, new_password: newPassword })
};

// 角色 API
export const roleApi = {
  getList: () =>
    api.post<unknown[]>('/roles/list'),
  
  create: (data: { name: string; code: string; description?: string }) =>
    api.post('/roles/create', data),
  
  update: (id: number, data: { name?: string; description?: string }) =>
    api.post('/roles/update', { id, ...data }),
  
  delete: (id: number) =>
    api.post('/roles/delete', { id }),
  
  assignPermissions: (id: number, permissionIds: number[]) =>
    api.post('/roles/assign-permissions', { id, permission_ids: permissionIds })
};

// 权限 API
export const permissionApi = {
  getList: () =>
    api.post<unknown[]>('/permissions/list'),
  
  create: (data: { name: string; code: string; resource: string; action: string; description?: string }) =>
    api.post('/permissions/create', data),
  
  update: (id: number, data: { name?: string; description?: string }) =>
    api.post('/permissions/update', { id, ...data }),
  
  delete: (id: number) =>
    api.post('/permissions/delete', { id })
};

// 图书 API
export const bookApi = {
  getList: (params?: { page?: number; page_size?: number; title?: string; author?: string; category?: string }) =>
    api.post<PageData<unknown>>('/books/list', params),
  
  get: (id: number) =>
    api.post<unknown>('/books/detail', { id }),
  
  create: (data: unknown) =>
    api.post('/books/create', data),
  
  update: (id: number, data: unknown) =>
    api.post('/books/update', { id, ...data }),
  
  delete: (id: number) =>
    api.post('/books/delete', { id }),
  
  getCategories: () =>
    api.post<string[]>('/books/categories')
};

// 库存 API
export const stockApi = {
  getList: (params?: { page?: number; page_size?: number; book_title?: string }) =>
    api.post<PageData<unknown>>('/stocks/list', params),
  
  get: (id: number) =>
    api.post<unknown>('/stocks/detail', { id }),
  
  stockIn: (bookId: number, quantity: number, reason?: string) =>
    api.post('/stocks/in', { book_id: bookId, quantity, reason }),
  
  stockOut: (bookId: number, quantity: number, reason: string) =>
    api.post('/stocks/out', { book_id: bookId, quantity, reason }),
  
  getRecords: (params?: { page?: number; page_size?: number; book_id?: number; type?: string }) =>
    api.post<PageData<unknown>>('/stocks/records', params),
  
  getLowStock: (threshold?: number) =>
    api.post<unknown[]>('/stocks/low', { threshold })
};

// 采购 API
export const purchaseApi = {
  getSuppliers: () =>
    api.post<unknown[]>('/suppliers/list'),
  
  createSupplier: (data: unknown) =>
    api.post('/suppliers/create', data),
  
  updateSupplier: (id: number, data: unknown) =>
    api.post('/suppliers/update', { id, ...data }),
  
  deleteSupplier: (id: number) =>
    api.post('/suppliers/delete', { id }),
  
  getOrders: (params?: { page?: number; page_size?: number; status?: string }) =>
    api.post<PageData<unknown>>('/purchases/list', params),
  
  getOrder: (id: number) =>
    api.post<unknown>('/purchases/detail', { id }),
  
  createOrder: (data: unknown) =>
    api.post('/purchases/create', data),
  
  updateOrderStatus: (id: number, status: string) =>
    api.post('/purchases/update-status', { id, status }),
  
  deleteOrder: (id: number) =>
    api.post('/purchases/delete', { id })
};

// 销售 API
export const saleApi = {
  getOrders: (params?: { page?: number; page_size?: number; status?: string; start_date?: string; end_date?: string }) =>
    api.post<PageData<unknown>>('/sales/list', params),
  
  getOrder: (id: number) =>
    api.post<unknown>('/sales/detail', { id }),
  
  createOrder: (data: unknown) =>
    api.post('/sales/create', data),
  
  cancelOrder: (id: number) =>
    api.post('/sales/cancel', { id }),
  
  getStats: (startDate?: string, endDate?: string) =>
    api.post('/sales/stats', { start_date: startDate, end_date: endDate }),
  
  getCart: () =>
    api.post<unknown[]>('/cart/list'),
  
  addToCart: (bookId: number, quantity: number) =>
    api.post('/cart/add', { book_id: bookId, quantity }),
  
  removeFromCart: (id: number) =>
    api.post('/cart/remove', { id })
};

// 借阅 API
export const borrowApi = {
  getRecords: (params?: { page?: number; page_size?: number; status?: string; user_id?: number }) =>
    api.post<PageData<unknown>>('/borrows/list', params),
  
  getRecord: (id: string) =>
    api.post<unknown>('/borrows/detail', { id }),
  
  borrowBook: (bookId: number, days?: number) =>
    api.post('/borrows/borrow', { book_id: bookId, days }),
  
  returnBook: (id: string) =>
    api.post('/borrows/return', { id }),
  
  renewBook: (id: string) =>
    api.post('/borrows/renew', { id }),
  
  getOverdueRecords: (params?: { page?: number; page_size?: number }) =>
    api.post<PageData<unknown>>('/borrows/overdue', params),
  
  getStats: () =>
    api.post('/borrows/stats')
};

// 仪表盘 API
export const dashboardApi = {
  getStats: () =>
    api.post<unknown>('/dashboard/stats')
};
