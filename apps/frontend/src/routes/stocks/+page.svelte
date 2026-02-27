<script lang="ts">
  import AppLayout from '$lib/components/layout/AppLayout.svelte';
  import { 
    Card, 
    CardHeader, 
    CardContent, 
    CardTitle,
    Button,
    Input,
    Table,
    TableHeader,
    TableBody,
    TableRow,
    TableHead,
    TableCell,
    Dialog,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
    Badge,
    Tabs,
    TabsList,
    TabsTrigger,
    TabsContent
  } from '$lib/components/ui';
  import { Plus, ArrowUp, ArrowDown, Package, History } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  // 模拟数据
  let stocks = $state([
    { id: 1, book_id: 1, book: { title: '深入理解计算机系统' }, total_quantity: 20, available_quantity: 15, borrowed_quantity: 3, sold_quantity: 2, location: 'A区-01-03' },
    { id: 2, book_id: 2, book: { title: 'JavaScript高级程序设计' }, total_quantity: 10, available_quantity: 8, borrowed_quantity: 2, sold_quantity: 0, location: 'A区-02-01' },
    { id: 3, book_id: 3, book: { title: 'Python编程从入门到实践' }, total_quantity: 5, available_quantity: 0, borrowed_quantity: 3, sold_quantity: 2, location: 'B区-01-05' },
    { id: 4, book_id: 4, book: { title: '算法导论' }, total_quantity: 15, available_quantity: 12, borrowed_quantity: 2, sold_quantity: 1, location: 'A区-03-02' }
  ]);

  let records = $state([
    { id: 1, book: { title: '深入理解计算机系统' }, type: 'in', quantity: 10, reason: '采购入库', operator: { real_name: '管理员' }, created_at: '2024-01-15 10:30:00' },
    { id: 2, book: { title: 'JavaScript高级程序设计' }, type: 'out', quantity: 2, reason: '销售出库', operator: { real_name: '管理员' }, created_at: '2024-01-15 09:20:00' },
    { id: 3, book: { title: 'Python编程从入门到实践' }, type: 'out', quantity: 1, reason: '报损', operator: { real_name: '管理员' }, created_at: '2024-01-14 16:45:00' }
  ]);

  let dialogOpen = $state(false);
  let dialogType = $state<'in' | 'out'>('in');
  let formData = $state({
    book_id: 0,
    quantity: 1,
    reason: ''
  });

  function openDialog(type: 'in' | 'out') {
    dialogType = type;
    formData = { book_id: 0, quantity: 1, reason: '' };
    dialogOpen = true;
  }

  function handleSubmit() {
    if (!formData.book_id || !formData.quantity || !formData.reason) {
      toast.error('请填写完整信息');
      return;
    }
    toast.success(dialogType === 'in' ? '入库成功' : '出库成功');
    dialogOpen = false;
  }
</script>

<svelte:head>
  <title>库存管理 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">库存管理</h1>
        <p class="text-gray-500 mt-1">管理图书库存和出入库记录</p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" on:click={() => openDialog('in')}>
          <ArrowUp class="w-4 h-4 mr-2" />
          入库
        </Button>
        <Button variant="outline" on:click={() => openDialog('out')}>
          <ArrowDown class="w-4 h-4 mr-2" />
          出库
        </Button>
      </div>
    </div>

    <!-- Tabs -->
    <Tabs value="stock">
      <TabsList>
        <TabsTrigger value="stock">库存列表</TabsTrigger>
        <TabsTrigger value="records">出入库记录</TabsTrigger>
      </TabsList>

      <TabsContent value="stock">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>图书名称</TableHead>
                  <TableHead class="text-center">总库存</TableHead>
                  <TableHead class="text-center">可借</TableHead>
                  <TableHead class="text-center">借出</TableHead>
                  <TableHead class="text-center">已售</TableHead>
                  <TableHead>存放位置</TableHead>
                  <TableHead>状态</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {#each stocks as stock}
                  <TableRow>
                    <TableCell class="font-medium">{stock.book.title}</TableCell>
                    <TableCell class="text-center">{stock.total_quantity}</TableCell>
                    <TableCell class="text-center">{stock.available_quantity}</TableCell>
                    <TableCell class="text-center">{stock.borrowed_quantity}</TableCell>
                    <TableCell class="text-center">{stock.sold_quantity}</TableCell>
                    <TableCell>{stock.location}</TableCell>
                    <TableCell>
                      <Badge variant={stock.available_quantity === 0 ? 'destructive' : stock.available_quantity < 5 ? 'secondary' : 'default'}>
                        {stock.available_quantity === 0 ? '缺货' : stock.available_quantity < 5 ? '低库存' : '正常'}
                      </Badge>
                    </TableCell>
                  </TableRow>
                {/each}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="records">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>图书名称</TableHead>
                  <TableHead>类型</TableHead>
                  <TableHead class="text-center">数量</TableHead>
                  <TableHead>原因</TableHead>
                  <TableHead>操作人</TableHead>
                  <TableHead>时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {#each records as record}
                  <TableRow>
                    <TableCell class="font-medium">{record.book.title}</TableCell>
                    <TableCell>
                      <Badge variant={record.type === 'in' ? 'default' : 'destructive'}>
                        {record.type === 'in' ? '入库' : '出库'}
                      </Badge>
                    </TableCell>
                    <TableCell class="text-center">{record.quantity}</TableCell>
                    <TableCell>{record.reason}</TableCell>
                    <TableCell>{record.operator.real_name}</TableCell>
                    <TableCell>{record.created_at}</TableCell>
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

<!-- Dialog -->
<Dialog bind:open={dialogOpen}>
  <DialogHeader>
    <DialogTitle>{dialogType === 'in' ? '图书入库' : '图书出库'}</DialogTitle>
    <DialogDescription>
      {dialogType === 'in' ? '填写入库信息' : '填写出库信息'}
    </DialogDescription>
  </DialogHeader>
  
  <div class="space-y-4 py-4">
    <div class="space-y-2">
      <Label>数量</Label>
      <Input type="number" bind:value={formData.quantity} placeholder="请输入数量" />
    </div>
    <div class="space-y-2">
      <Label>原因</Label>
      <Input bind:value={formData.reason} placeholder="请输入原因" />
    </div>
  </div>
  
  <DialogFooter>
    <Button variant="outline" on:click={() => dialogOpen = false}>取消</Button>
    <Button on:click={handleSubmit}>确认</Button>
  </DialogFooter>
</Dialog>

<script lang="ts">
  import { Label } from '$lib/components/ui';
</script>
