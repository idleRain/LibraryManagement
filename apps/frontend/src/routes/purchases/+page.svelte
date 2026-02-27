<script lang="ts">
  import AppLayout from '$lib/components/layout/AppLayout.svelte';
  import { 
    Card, 
    CardContent,
    Button,
    Table,
    TableHeader,
    TableBody,
    TableRow,
    TableHead,
    TableCell,
    Badge,
    Tabs,
    TabsList,
    TabsTrigger,
    TabsContent
  } from '$lib/components/ui';
  import { Plus, Building2, FileText } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  let suppliers = $state([
    { id: 1, name: '北京图书批发中心', contact: '张经理', phone: '010-12345678', status: 1 },
    { id: 2, name: '上海书城', contact: '李经理', phone: '021-87654321', status: 1 },
    { id: 3, name: '广州图书市场', contact: '王经理', phone: '020-11112222', status: 0 }
  ]);

  let orders = $state([
    { id: 1, order_no: 'PO202401150001', supplier: { name: '北京图书批发中心' }, total_amount: 15000, status: 'completed', created_at: '2024-01-15' },
    { id: 2, order_no: 'PO202401160001', supplier: { name: '上海书城' }, total_amount: 8500, status: 'pending', created_at: '2024-01-16' },
    { id: 3, order_no: 'PO202401170001', supplier: { name: '广州图书市场' }, total_amount: 12000, status: 'approved', created_at: '2024-01-17' }
  ]);

  const statusMap: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' }> = {
    pending: { label: '待审核', variant: 'secondary' },
    approved: { label: '已审核', variant: 'default' },
    completed: { label: '已完成', variant: 'default' },
    cancelled: { label: '已取消', variant: 'destructive' }
  };
</script>

<svelte:head>
  <title>采购管理 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">采购管理</h1>
        <p class="text-gray-500 mt-1">管理供应商和采购订单</p>
      </div>
      <Button>
        <Plus class="w-4 h-4 mr-2" />
        新建采购单
      </Button>
    </div>

    <Tabs value="orders">
      <TabsList>
        <TabsTrigger value="orders">
          <FileText class="w-4 h-4 mr-2" />
          采购订单
        </TabsTrigger>
        <TabsTrigger value="suppliers">
          <Building2 class="w-4 h-4 mr-2" />
          供应商
        </TabsTrigger>
      </TabsList>

      <TabsContent value="orders">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>订单号</TableHead>
                  <TableHead>供应商</TableHead>
                  <TableHead class="text-right">金额</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>创建时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {#each orders as order}
                  <TableRow>
                    <TableCell class="font-mono">{order.order_no}</TableCell>
                    <TableCell>{order.supplier.name}</TableCell>
                    <TableCell class="text-right">¥{order.total_amount.toLocaleString()}</TableCell>
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
      </TabsContent>

      <TabsContent value="suppliers">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>供应商名称</TableHead>
                  <TableHead>联系人</TableHead>
                  <TableHead>电话</TableHead>
                  <TableHead>状态</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {#each suppliers as supplier}
                  <TableRow>
                    <TableCell class="font-medium">{supplier.name}</TableCell>
                    <TableCell>{supplier.contact}</TableCell>
                    <TableCell>{supplier.phone}</TableCell>
                    <TableCell>
                      <Badge variant={supplier.status === 1 ? 'default' : 'secondary'}>
                        {supplier.status === 1 ? '启用' : '禁用'}
                      </Badge>
                    </TableCell>
                  </TableRow>
                {/each}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</AppLayout>
