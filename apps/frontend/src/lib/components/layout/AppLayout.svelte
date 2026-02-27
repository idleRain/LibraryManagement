<script lang="ts">
  import { page } from '$app/stores';
  import { auth, ui } from '$lib/stores';
  import { goto } from '$app/navigation';
  import { 
    Sidebar, 
    SidebarHeader, 
    SidebarContent, 
    SidebarMenu, 
    SidebarMenuItem,
    Avatar,
    AvatarFallback,
    Button,
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    Sonner
  } from '$lib/components/ui';
  import { 
    Book, 
    Package, 
    ShoppingCart, 
    FileText, 
    BookOpen, 
    Users, 
    Settings,
    LogOut,
    Menu,
    Home,
    LayoutDashboard
  } from 'lucide-svelte';

  const menuItems = [
    { path: '/', label: '仪表盘', icon: LayoutDashboard },
    { path: '/books', label: '图书管理', icon: Book },
    { path: '/stocks', label: '库存管理', icon: Package },
    { path: '/purchases', label: '采购管理', icon: ShoppingCart },
    { path: '/sales', label: '销售管理', icon: FileText },
    { path: '/borrows', label: '借阅管理', icon: BookOpen },
    { path: '/users', label: '用户管理', icon: Users },
    { path: '/roles', label: '角色权限', icon: Settings }
  ];

  function handleLogout() {
    auth.logout();
    goto('/auth/login');
  }

  $: currentPath = $page.url.pathname;
  $: user = $auth.user;
  $: sidebarOpen = $ui.sidebarOpen;
</script>

<div class="min-h-screen flex bg-gray-50">
  <!-- Sidebar -->
  <aside class="fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-gray-200 transform transition-transform duration-300 lg:translate-x-0" class:hidden={!sidebarOpen}>
    <Sidebar>
      <!-- Logo -->
      <SidebarHeader>
        <div class="flex items-center gap-3">
          <div class="flex items-center justify-center w-10 h-10 rounded-lg bg-primary text-white">
            <Book class="w-6 h-6" />
          </div>
          <div>
            <h1 class="text-lg font-bold text-gray-900">图书管理</h1>
            <p class="text-xs text-gray-500">Library System</p>
          </div>
        </div>
      </SidebarHeader>

      <!-- Navigation -->
      <SidebarContent>
        <SidebarMenu>
          {#each menuItems as item}
            <a href={item.path}>
              <SidebarMenuItem active={currentPath === item.path || (item.path !== '/' && currentPath.startsWith(item.path))}>
                <item.icon class="w-5 h-5" />
                <span>{item.label}</span>
              </SidebarMenuItem>
            </a>
          {/each}
        </SidebarMenu>
      </SidebarContent>
    </Sidebar>
  </aside>

  <!-- Main Content -->
  <div class="flex-1 lg:ml-64">
    <!-- Header -->
    <header class="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6 sticky top-0 z-40">
      <button
        on:click={() => ui.setSidebarOpen(!sidebarOpen)}
        class="lg:hidden p-2 rounded-lg hover:bg-gray-100"
      >
        <Menu class="w-6 h-6" />
      </button>

      <div class="flex-1"></div>

      <!-- User Menu -->
      <div class="flex items-center gap-4">
        <DropdownMenu>
          <DropdownMenuTrigger class="flex items-center gap-3 cursor-pointer">
            <div class="text-right hidden sm:block">
              <div class="text-sm font-medium text-gray-700">
                {user?.real_name || user?.username || '用户'}
              </div>
              <div class="text-xs text-gray-500">
                {user?.roles?.map(r => r.name).join(', ') || '普通用户'}
              </div>
            </div>
            <Avatar>
              <AvatarFallback class="bg-primary text-white">
                {(user?.real_name || user?.username || 'U').charAt(0).toUpperCase()}
              </AvatarFallback>
            </Avatar>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem>
              <a href="/profile" class="flex items-center gap-2 w-full">
                <Settings class="w-4 h-4" />
                个人设置
              </a>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem on:click={handleLogout} class="text-red-600">
              <LogOut class="w-4 h-4" />
              退出登录
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>

    <!-- Page Content -->
    <main class="p-6">
      <slot />
    </main>
  </div>

  <!-- Overlay for mobile -->
  {#if sidebarOpen}
    <div
      class="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
      on:click={() => ui.setSidebarOpen(false)}
      on:keydown={(e) => e.key === 'Escape' && ui.setSidebarOpen(false)}
      role="button"
      tabindex="0"
    ></div>
  {/if}
</div>

<Sonner position="top-center" />
