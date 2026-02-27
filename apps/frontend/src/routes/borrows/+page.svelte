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
    Badge,
    Tabs,
    TabsList,
    TabsTrigger,
    TabsContent
  } from '$lib/components/ui';
  import { BookOpen, Clock, AlertTriangle, RotateCcw } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  let records = $state([
    { id: '1', book_title: '深入理解计算机系统', user_name: '张三', borrow_date: '2024-01-10', due_date: '2024-02-10', status: 'borrowed', renew_count: 0 },
    { id: '2', book_title: 'JavaScript高级程序设计', user_name: '李四', borrow_date: '2024-01-05', due_date: '2024-01-20', status: 'overdue', renew_count: 1 },
    { id: '3', book_title: 'Python编程从入门到实践', user_name: '王五', borrow_date: '2024-01-08', due_date: '2024-02-08', return_date: '2024-01-25', status: 'returned', renew_count: 0 }
  ]);

  const stats = {
    borrowed: 156,
    overdue: 8,
    returned: 342
  };

  const statusMap: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' }> = {
    borrowed: { label: '借阅中', variant: 'default' },
    returned: { label: '已归还', variant: 'secondary' },
    overdue: { label: '已逾期', variant: 'destructive' },
    lost: { label: '已丢失', variant: 'destructive' }
  };

  function handleReturn(id: string) {
    toast.success('归还成功');
  }

  function handleRenew(id: string) {
    toast.success('续借成功');
  }
</script>

<svelte:head>
  <title>借阅管理 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">借阅管理</h1>
        <p class="text-gray-500 mt-1">管理图书借阅、归还和逾期</p>
      </div>
      <Button>
        <BookOpen class="w-4 h-4 mr-2" />
        新建借阅
      </Button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">借阅中</p>
              <p class="text-2xl font-bold text-gray-900">{stats.borrowed}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
              <BookOpen class="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">已逾期</p>
              <p class="text-2xl font-bold text-red-600">{stats.overdue}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center">
              <AlertTriangle class="w-6 h-6 text-red-600" />
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500">已归还</p>
              <p class="text-2xl font-bold text-gray-900">{stats.returned}</p>
            </div>
            <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
              <RotateCcw class="w-6 h-6 text-green-600" />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Tabs -->
    <Tabs value="all">
      <TabsList>
        <TabsTrigger value="all">全部</TabsTrigger>
        <TabsTrigger value="borrowed">借阅中</TabsTrigger>
        <TabsTrigger value="overdue">已逾期</TabsTrigger>
        <TabsTrigger value="returned">已归还</TabsTrigger>
      </TabsList>

      <TabsContent value="all">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>图书名称</TableHead>
                  <TableHead>借阅人</TableHead>
                  <TableHead>借阅日期</TableHead>
                  <TableHead>应还日期</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead class="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {#each records as record}
                  <TableRow>
                    <TableCell class="font-medium">{record.book_title}</TableCell>
                    <TableCell>{record.user_name}</TableCell>
                    <TableCell>{record.borrow_date}</TableCell>
                    <TableCell>{record.due_date}</TableCell>
                    <TableCell>
                      <Badge variant={statusMap[record.status].variant}>
                        {statusMap[record.status].label}
                      </Badge>
                    </TableCell>
                    <TableCell class="text-right">
                      {#if record.status === 'borrowed' || record.status === 'overdue'}
                        <div class="flex items-center justify-end gap-2">
                          <Button variant="outline" size="sm" on:click={() => handleReturn(record.id)}>
                            归还
                          </Button>
                          {#if record.renew_count < 2}
                            <Button variant="outline" size="sm" on:click={() => handleRenew(record.id)}>
                              续借
                            </Button>
                          {/if}
                        </div>
                      {/if}
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
