import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { User } from '$lib/types';

interface AuthState {
  token: string | null;
  user: User | null;
  permissions: string[];
  isAuthenticated: boolean;
}

function createAuthStore() {
  const initialState: AuthState = {
    token: null,
    user: null,
    permissions: [],
    isAuthenticated: false
  };

  // 从 localStorage 恢复状态
  if (browser) {
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    const permissionsStr = localStorage.getItem('permissions');
    
    if (token && userStr) {
      initialState.token = token;
      try {
        initialState.user = JSON.parse(userStr);
      } catch {
        initialState.user = null;
      }
      initialState.permissions = permissionsStr ? JSON.parse(permissionsStr) : [];
      initialState.isAuthenticated = true;
    }
  }

  const { subscribe, set, update } = writable<AuthState>(initialState);

  return {
    subscribe,
    
    login: (token: string, user: User, permissions: string[]) => {
      if (browser) {
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('permissions', JSON.stringify(permissions));
      }
      set({
        token,
        user,
        permissions,
        isAuthenticated: true
      });
    },
    
    logout: () => {
      if (browser) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        localStorage.removeItem('permissions');
      }
      set({
        token: null,
        user: null,
        permissions: [],
        isAuthenticated: false
      });
    },
    
    updatePermissions: (permissions: string[]) => {
      if (browser) {
        localStorage.setItem('permissions', JSON.stringify(permissions));
      }
      update(state => ({
        ...state,
        permissions
      }));
    },
    
    hasPermission: (permission: string): boolean => {
      let hasPerm = false;
      const unsub = subscribe(state => {
        hasPerm = state.permissions.includes(permission) || state.permissions.includes('*');
      });
      unsub();
      return hasPerm;
    }
  };
}

export const auth = createAuthStore();

// UI 状态管理
interface UIState {
  sidebarOpen: boolean;
  theme: 'light' | 'dark';
  loading: boolean;
}

function createUIStore() {
  const { subscribe, set, update } = writable<UIState>({
    sidebarOpen: true,
    theme: 'light',
    loading: false
  });

  return {
    subscribe,
    
    toggleSidebar: () => {
      update(state => ({
        ...state,
        sidebarOpen: !state.sidebarOpen
      }));
    },
    
    setSidebarOpen: (open: boolean) => {
      update(state => ({
        ...state,
        sidebarOpen: open
      }));
    },
    
    setTheme: (theme: 'light' | 'dark') => {
      if (browser) {
        document.documentElement.classList.toggle('dark', theme === 'dark');
      }
      update(state => ({
        ...state,
        theme
      }));
    },
    
    setLoading: (loading: boolean) => {
      update(state => ({
        ...state,
        loading
      }));
    }
  };
}

export const ui = createUIStore();
