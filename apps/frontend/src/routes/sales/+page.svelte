<script lang="ts">
  import AppLayout from '$lib/components/layout/AppLayout.svelte';
  import { 
    Card, 
    CardContent,
    CardHeader,
    CardTitle,
    Button,
    Table,
    TableHeader,
    TableBody,
    TableRow,
    TableHead,
    TableCell,
    Badge
  } from '$lib/components/ui';
  import { Plus, ShoppingCart, DollarSign, FileText, TrendingUp } from 'lucide-svelte';

  let orders = $state([
    { id: 1, order_no: 'SO202401150001', customer_name: '张三', total_amount: 268, status: 'paid', payment_method: '微信支付', created_at: '2024-01-15 14:30' },
    { id: 2, order_no: 'SO202401150002', customer_name: '李四', total_amount: 156, status: 'paid', payment_method: '支付宝', created_at: '2024-01-15 15:20' },
    { id: 3, order_no: 'SO202401160001', customer_name: '王五', total_amount: 89, status: 'completed', payment_method: '现金', created_at: '2024-01-16 10:15' }
  ]);

  const stats = {
    todaySales: 3580,
    todayOrders: 28,
    monthSales: 86500
  };

  const statusMap: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' }> = {
    pending: { label: '待支付', variant: 'secondary' },
    paid: { label: '已支付', variant: 'default' },
    completed: { label: '已完成', variant: 'default' },
    cancelled: { label: '已取消', variant: 'destructive' }
  };
</script>

<svelte:head>
  <title>销售管理 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">销售管理</h1>
        <p class="text-gray-500 mt-1">管理销售订单和统计数据</p>
      </div>
      <Button>
        <Plus class="w-4 h-4 mr-2" />
        新建销售单
      </Button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">今日销售额</p>
              <p class="text-2xl font-bold text-gray-900">¥{stats.todaySales.toLocaleString()}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
              <DollarSign class="w-6 h-6 text-green-600" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">今日订单</p>
              <p class="text-2xl font-bold text-gray-900">{stats.todayOrders}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
              <ShoppingCart class="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">本月销售额</p>
              <p class="text-2xl font-bold text-gray-900">¥{stats.monthSales.toLocaleString()}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-purple-100 flex items-center justify-center">
              <TrendingUp class="w-6 h-6 text-purple-600" />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Orders Table -->
    <Card>
      <CardHeader>
        <CardTitle>销售订单</CardTitle>
      </CardHeader>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>订单号</TableHead>
              <TableHead>客户</TableHead>
              <TableHead class="text-right">金额</TableHead>
              <TableHead>支付方式</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {#each orders as order}
              <TableRow>
                <TableCell class="font-mono">{order.order_no}</TableCell>
                <TableCell>{order.customer_name}</TableCell>
                <TableCell class="text-right">¥{order.total_amount}</TableCell>
                <TableCell>{order.payment_method}</TableCell>
                <TableCell>
                  <Badge variant={statusMap[order.status].variant}>
                    {statusMap[order.status].label}
                  </Badge>
                </TableCell>
                <TableCell>{order.created_at}</TableCell>
              </TableRow>
            {/each}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</AppLayout>
