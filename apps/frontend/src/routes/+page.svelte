<script lang="ts">
  import AppLayout from '$lib/components/layout/AppLayout.svelte';
  import { Card, CardHeader, CardContent, CardTitle } from '$lib/components/ui';
  import { Book, Users, BookOpen, DollarSign, TrendingUp, AlertTriangle } from 'lucide-svelte';

  // 模拟数据
  const stats = {
    totalBooks: 12580,
    totalUsers: 3420,
    totalBorrowed: 856,
    totalSales: 125680.50,
    lowStockCount: 12,
    overdueCount: 5
  };

  const recentBorrows = [
    { id: 1, book_title: '深入理解计算机系统', user_name: '张三', borrow_date: '2024-01-15', status: 'borrowed' },
    { id: 2, book_title: 'JavaScript高级程序设计', user_name: '李四', borrow_date: '2024-01-14', status: 'returned' },
    { id: 3, book_title: 'Python编程从入门到实践', user_name: '王五', borrow_date: '2024-01-13', status: 'overdue' },
    { id: 4, book_title: '算法导论', user_name: '赵六', borrow_date: '2024-01-12', status: 'borrowed' },
    { id: 5, book_title: '设计模式', user_name: '钱七', borrow_date: '2024-01-11', status: 'borrowed' }
  ];

  const lowStockBooks = [
    { id: 1, title: '深入理解计算机系统', available: 2 },
    { id: 2, title: 'JavaScript高级程序设计', available: 1 },
    { id: 3, title: 'Python编程从入门到实践', available: 0 },
    { id: 4, title: '算法导论', available: 3 }
  ];
</script>

<svelte:head>
  <title>仪表盘 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900">仪表盘</h1>
      <p class="text-gray-500 mt-1">欢迎回来，查看系统概览</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">图书总数</p>
              <p class="text-2xl font-bold text-gray-900">{stats.totalBooks.toLocaleString()}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
              <Book class="w-6 h-6 text-blue-600" />
            </div>
          </div>
          <div class="mt-4 flex items-center text-sm text-green-600">
            <TrendingUp class="w-4 h-4 mr-1" />
            <span>较上月 +12%</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">注册用户</p>
              <p class="text-2xl font-bold text-gray-900">{stats.totalUsers.toLocaleString()}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
              <Users class="w-6 h-6 text-green-600" />
            </div>
          </div>
          <div class="mt-4 flex items-center text-sm text-green-600">
            <TrendingUp class="w-4 h-4 mr-1" />
            <span>较上月 +8%</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">借阅中</p>
              <p class="text-2xl font-bold text-gray-900">{stats.totalBorrowed}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-orange-100 flex items-center justify-center">
              <BookOpen class="w-6 h-6 text-orange-600" />
            </div>
          </div>
          <div class="mt-4 flex items-center text-sm text-red-600">
            <AlertTriangle class="w-4 h-4 mr-1" />
            <span>{stats.overdueCount} 本逾期</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">销售总额</p>
              <p class="text-2xl font-bold text-gray-900">¥{stats.totalSales.toLocaleString()}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-purple-100 flex items-center justify-center">
              <DollarSign class="w-6 h-6 text-purple-600" />
            </div>
          </div>
          <div class="mt-4 flex items-center text-sm text-green-600">
            <TrendingUp class="w-4 h-4 mr-1" />
            <span>较上月 +15%</span>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent Borrows -->
      <Card>
        <CardHeader>
          <CardTitle>最近借阅</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            {#each recentBorrows as borrow}
              <div class="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <div>
                  <p class="font-medium text-gray-900">{borrow.book_title}</p>
                  <p class="text-sm text-gray-500">{borrow.user_name} · {borrow.borrow_date}</p>
                </div>
                <span class="px-2 py-1 text-xs rounded-full {borrow.status === 'borrowed' ? 'bg-blue-100 text-blue-700' : borrow.status === 'returned' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}">
                  {borrow.status === 'borrowed' ? '借阅中' : borrow.status === 'returned' ? '已归还' : '已逾期'}
                </span>
              </div>
            {/each}
          </div>
        </CardContent>
      </Card>

      <!-- Low Stock Alert -->
      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <AlertTriangle class="w-5 h-5 text-red-500" />
            库存预警
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            {#each lowStockBooks as book}
              <div class="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <p class="font-medium text-gray-900">{book.title}</p>
                <span class="px-2 py-1 text-xs rounded-full {book.available === 0 ? 'bg-red-100 text-red-700' : 'bg-yellow-100 text-yellow-700'}">
                  {book.available === 0 ? '已售罄' : `仅剩 ${book.available} 本`}
                </span>
              </div>
            {/each}
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Quick Actions -->
    <Card>
      <CardHeader>
        <CardTitle>快捷操作</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <a href="/books" class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-primary hover:bg-primary-50 transition-colors">
            <Book class="w-8 h-8 text-primary mb-2" />
            <span class="text-sm font-medium">图书管理</span>
          </a>
          <a href="/borrows" class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-primary hover:bg-primary-50 transition-colors">
            <BookOpen class="w-8 h-8 text-primary mb-2" />
            <span class="text-sm font-medium">借阅管理</span>
          </a>
          <a href="/sales" class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-primary hover:bg-primary-50 transition-colors">
            <DollarSign class="w-8 h-8 text-primary mb-2" />
            <span class="text-sm font-medium">销售管理</span>
          </a>
          <a href="/stocks" class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-primary hover:bg-primary-50 transition-colors">
            <Package class="w-8 h-8 text-primary mb-2" />
            <span class="text-sm font-medium">库存管理</span>
          </a>
        </div>
      </CardContent>
    </Card>
  </div>
</AppLayout>

<script lang="ts">
  import { Package } from 'lucide-svelte';
</script>
