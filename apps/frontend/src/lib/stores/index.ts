import { writable, derived } from 'svelte/store';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';

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
}

// 认证状态
export interface AuthState {
  user: User | null;
  token: string | null;
  permissions: string[];
  isAuthenticated: boolean;
  loading: boolean;
}

// 初始状态
const initialState: AuthState = {
  user: null,
  token: null,
  permissions: [],
  isAuthenticated: false,
  loading: true
};

// 创建认证状态存储
function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>(initialState);

  return {
    subscribe,

    // 初始化（从 localStorage 恢复）
    init: () => {
      if (typeof window !== 'undefined') {
        const token = localStorage.getItem('token');
        const userStr = localStorage.getItem('user');
        const permissionsStr = localStorage.getItem('permissions');

        if (token && userStr) {
          try {
            const user = JSON.parse(userStr);
            const permissions = permissionsStr ? JSON.parse(permissionsStr) : [];
            set({ user, token, permissions, isAuthenticated: true, loading: false });
          } catch {
            authStore.clear();
          }
        } else {
          update(s => ({ ...s, loading: false }));
        }
      }
    },

    // 保存认证状态
    save: (token: string, user: User, permissions: string[]) => {
      if (typeof window !== 'undefined') {
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('permissions', JSON.stringify(permissions));
      }
      set({ user, token, permissions, isAuthenticated: true, loading: false });
    },

    // 清除认证状态
    clear: () => {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        localStorage.removeItem('permissions');
      }
      set({ ...initialState, loading: false });
    },

    // 更新用户信息
    updateUser: (user: Partial<User>) => {
      update(state => {
        if (state.user) {
          const newUser = { ...state.user, ...user };
          if (typeof window !== 'undefined') {
            localStorage.setItem('user', JSON.stringify(newUser));
          }
          return { ...state, user: newUser };
        }
        return state;
      });
    },

    // 更新权限
    updatePermissions: (permissions: string[]) => {
      update(state => {
        if (typeof window !== 'undefined') {
          localStorage.setItem('permissions', JSON.stringify(permissions));
        }
        return { ...state, permissions };
      });
    }
  };
}

export const authStore = createAuthStore();

// 派生状态：用户显示名称
export const userDisplayName = derived(authStore, $auth => {
  if (!$auth.user) return '';
  return $auth.user.real_name || $auth.user.username;
});

// 派生状态：是否是管理员
export const isAdmin = derived(authStore, $auth => {
  if (!$auth.user) return false;
  return $auth.user.roles?.some(r => r.code === 'admin' || r.code === 'super_admin') || false;
});

// 派生状态：角色名称列表
export const roleNames = derived(authStore, $auth => {
  if (!$auth.user) return [];
  return $auth.user.roles?.map(r => r.name) || [];
});

// 检查权限
export function hasPermission(permission: string): boolean {
  let has = false;
  authStore.subscribe(state => {
    has = state.permissions.includes(permission) || state.permissions.includes(permission.split(':')[0]);
  })();
  return has;
}

// 检查多个权限（任一满足）
export function hasAnyPermission(permissions: string[]): boolean {
  let has = false;
  authStore.subscribe(state => {
    has = permissions.some(p =>
      state.permissions.includes(p) ||
      state.permissions.includes(p.split(':')[0])
    );
  })();
  return has;
}

// 检查多个权限（全部满足）
export function hasAllPermissions(permissions: string[]): boolean {
  let has = false;
  authStore.subscribe(state => {
    has = permissions.every(p =>
      state.permissions.includes(p) ||
      state.permissions.includes(p.split(':')[0])
    );
  })();
  return has;
}

// Toast 状态管理
function createToastStore() {
  const { subscribe, update } = writable<{ count: number }>({ count: 0 });

  return {
    subscribe,
    success: (message: string) => {
      toast.success(message);
      update(s => ({ count: s.count + 1 }));
    },
    error: (message: string) => {
      toast.error(message);
      update(s => ({ count: s.count + 1 }));
    },
    warning: (message: string) => {
      toast.warning(message);
      update(s => ({ count: s.count + 1 }));
    },
    info: (message: string) => {
      toast.info(message);
      update(s => ({ count: s.count + 1 }));
    }
  };
}

export const toastStore = createToastStore();

// 侧边栏状态
export const sidebarCollapsed = writable<boolean>(false);

// 主题状态
export const theme = writable<'light' | 'dark'>('light');

// 初始化主题
export function initTheme() {
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('theme') as 'light' | 'dark' | null;
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    theme.set(saved || (prefersDark ? 'dark' : 'light'));
  }
}

// 切换主题
export function toggleTheme() {
  theme.update(t => {
    const newTheme = t === 'light' ? 'dark' : 'light';
    if (typeof window !== 'undefined') {
      localStorage.setItem('theme', newTheme);
      document.documentElement.classList.toggle('dark', newTheme === 'dark');
    }
    return newTheme;
  });
}

// 加载状态
export const loading = writable<boolean>(false);

// 全局搜索
export const searchQuery = writable<string>('');

// 分页状态
export interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

export function createPaginationStore(initialPage = 1, initialPageSize = 10) {
  const { subscribe, set, update } = writable<PaginationState>({
    page: initialPage,
    pageSize: initialPageSize,
    total: 0
  });

  return {
    subscribe,
    setPage: (page: number) => update(s => ({ ...s, page })),
    setPageSize: (pageSize: number) => update(s => ({ ...s, pageSize, page: 1 })),
    setTotal: (total: number) => update(s => ({ ...s, total })),
    reset: () => set({ page: initialPage, pageSize: initialPageSize, total: 0 })
  };
}
