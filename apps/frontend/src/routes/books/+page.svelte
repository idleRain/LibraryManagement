<script lang="ts">
  import AppLayout from '$lib/components/layout/AppLayout.svelte';
  import { 
    Card, 
    CardHeader, 
    CardContent, 
    CardTitle,
    Button,
    Input,
    Label,
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
    Select,
    SelectTrigger,
    SelectContent,
    SelectItem,
    Textarea,
    Pagination
  } from '$lib/components/ui';
  import { Plus, Search, Edit, Trash2, BookOpen, Package } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  // 模拟数据
  let books = $state([
    { id: 1, isbn: '978-7-111-54784-2', title: '深入理解计算机系统', author: 'Randal E. Bryant', publisher: '机械工业出版社', category: '计算机', price: 139.00, stock: { available_quantity: 15, total_quantity: 20 } },
    { id: 2, isbn: '978-7-115-42533-4', title: 'JavaScript高级程序设计', author: 'Matt Frisbie', publisher: '人民邮电出版社', category: '计算机', price: 129.00, stock: { available_quantity: 8, total_quantity: 10 } },
    { id: 3, isbn: '978-7-115-42602-7', title: 'Python编程从入门到实践', author: 'Eric Matthes', publisher: '人民邮电出版社', category: '计算机', price: 89.00, stock: { available_quantity: 0, total_quantity: 5 } },
    { id: 4, isbn: '978-7-111-40701-0', title: '算法导论', author: 'Thomas H. Cormen', publisher: '机械工业出版社', category: '计算机', price: 128.00, stock: { available_quantity: 12, total_quantity: 15 } },
    { id: 5, isbn: '978-7-115-42813-7', title: '设计模式', author: 'Erich Gamma', publisher: '人民邮电出版社', category: '计算机', price: 99.00, stock: { available_quantity: 5, total_quantity: 8 } }
  ]);

  let searchTitle = $state('');
  let searchCategory = $state('');
  let currentPage = $state(1);
  let totalPages = $state(10);
  let dialogOpen = $state(false);
  let editingBook = $state<any>(null);

  // 表单字段
  let formData = $state({
    isbn: '',
    title: '',
    author: '',
    publisher: '',
    category: '',
    price: 0,
    description: '',
    pages: 0,
    language: '中文'
  });

  const categories = ['计算机', '文学', '历史', '科学', '艺术', '经济', '教育', '其他'];

  function handleSearch() {
    toast.info('搜索功能需要连接后端API');
  }

  function openAddDialog() {
    editingBook = null;
    formData = {
      isbn: '',
      title: '',
      author: '',
      publisher: '',
      category: '',
      price: 0,
      description: '',
      pages: 0,
      language: '中文'
    };
    dialogOpen = true;
  }

  function openEditDialog(book: any) {
    editingBook = book;
    formData = {
      isbn: book.isbn,
      title: book.title,
      author: book.author,
      publisher: book.publisher,
      category: book.category,
      price: book.price,
      description: book.description || '',
      pages: book.pages || 0,
      language: book.language || '中文'
    };
    dialogOpen = true;
  }

  function handleSave() {
    if (!formData.title || !formData.isbn) {
      toast.error('请填写必填字段');
      return;
    }

    if (editingBook) {
      // 更新
      const index = books.findIndex(b => b.id === editingBook.id);
      if (index !== -1) {
        books[index] = { ...books[index], ...formData };
      }
      toast.success('更新成功');
    } else {
      // 新增
      books.push({
        id: Date.now(),
        ...formData,
        stock: { available_quantity: 0, total_quantity: 0 }
      });
      toast.success('添加成功');
    }
    dialogOpen = false;
  }

  function handleDelete(book: any) {
    if (confirm(`确定要删除《${book.title}》吗？`)) {
      books = books.filter(b => b.id !== book.id);
      toast.success('删除成功');
    }
  }

  function handlePageChange(page: number) {
    currentPage = page;
  }
</script>

<svelte:head>
  <title>图书管理 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">图书管理</h1>
        <p class="text-gray-500 mt-1">管理图书馆所有图书信息</p>
      </div>
      <Button on:click={openAddDialog}>
        <Plus class="w-4 h-4 mr-2" />
        添加图书
      </Button>
    </div>

    <!-- Search Card -->
    <Card>
      <CardContent class="p-4">
        <div class="flex flex-wrap gap-4">
          <div class="flex-1 min-w-[200px]">
            <Input
              placeholder="搜索图书名称..."
              bind:value={searchTitle}
              on:keypress={(e) => e.key === 'Enter' && handleSearch()}
            >
              <Search class="w-4 h-4 text-gray-400" slot="icon" />
            </Input>
          </div>
          <div class="w-48">
            <Select bind:value={searchCategory}>
              <SelectTrigger>
                <span>{searchCategory || '选择分类'}</span>
              </SelectTrigger>
              <SelectContent>
                {#each categories as cat}
                  <SelectItem value={cat}>{cat}</SelectItem>
                {/each}
              </SelectContent>
            </Select>
          </div>
          <Button variant="outline" on:click={handleSearch}>
            <Search class="w-4 h-4 mr-2" />
            搜索
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Table Card -->
    <Card>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ISBN</TableHead>
              <TableHead>书名</TableHead>
              <TableHead>作者</TableHead>
              <TableHead>出版社</TableHead>
              <TableHead>分类</TableHead>
              <TableHead>价格</TableHead>
              <TableHead>库存</TableHead>
              <TableHead class="text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {#each books as book}
              <TableRow>
                <TableCell class="font-mono text-sm">{book.isbn}</TableCell>
                <TableCell class="font-medium">{book.title}</TableCell>
                <TableCell>{book.author}</TableCell>
                <TableCell>{book.publisher}</TableCell>
                <TableCell>
                  <Badge variant="secondary">{book.category}</Badge>
                </TableCell>
                <TableCell>¥{book.price.toFixed(2)}</TableCell>
                <TableCell>
                  <span class="{book.stock.available_quantity === 0 ? 'text-red-600' : book.stock.available_quantity < 5 ? 'text-yellow-600' : 'text-green-600'}">
                    {book.stock.available_quantity} / {book.stock.total_quantity}
                  </span>
                </TableCell>
                <TableCell class="text-right">
                  <div class="flex items-center justify-end gap-2">
                    <Button variant="ghost" size="sm" on:click={() => openEditDialog(book)}>
                      <Edit class="w-4 h-4" />
                    </Button>
                    <Button variant="ghost" size="sm" on:click={() => handleDelete(book)}>
                      <Trash2 class="w-4 h-4 text-red-500" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            {/each}
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Pagination -->
    <div class="flex justify-center">
      <Pagination {currentPage} {totalPages} onPageChange={handlePageChange} />
    </div>
  </div>
</AppLayout>

<!-- Add/Edit Dialog -->
<Dialog bind:open={dialogOpen}>
  <DialogHeader>
    <DialogTitle>{editingBook ? '编辑图书' : '添加图书'}</DialogTitle>
    <DialogDescription>
      {editingBook ? '修改图书信息' : '填写新图书信息'}
    </DialogDescription>
  </DialogHeader>
  
  <div class="space-y-4 py-4">
    <div class="grid grid-cols-2 gap-4">
      <div class="space-y-2">
        <Label>ISBN *</Label>
        <Input bind:value={formData.isbn} placeholder="978-7-XXX-XXXXX-X" />
      </div>
      <div class="space-y-2">
        <Label>书名 *</Label>
        <Input bind:value={formData.title} placeholder="请输入书名" />
      </div>
    </div>
    
    <div class="grid grid-cols-2 gap-4">
      <div class="space-y-2">
        <Label>作者</Label>
        <Input bind:value={formData.author} placeholder="请输入作者" />
      </div>
      <div class="space-y-2">
        <Label>出版社</Label>
        <Input bind:value={formData.publisher} placeholder="请输入出版社" />
      </div>
    </div>
    
    <div class="grid grid-cols-3 gap-4">
      <div class="space-y-2">
        <Label>分类</Label>
        <Select bind:value={formData.category}>
          <SelectTrigger>
            <span>{formData.category || '选择分类'}</span>
          </SelectTrigger>
          <SelectContent>
            {#each categories as cat}
              <SelectItem value={cat}>{cat}</SelectItem>
            {/each}
          </SelectContent>
        </Select>
      </div>
      <div class="space-y-2">
        <Label>价格</Label>
        <Input type="number" bind:value={formData.price} placeholder="0.00" />
      </div>
      <div class="space-y-2">
        <Label>页数</Label>
        <Input type="number" bind:value={formData.pages} placeholder="0" />
      </div>
    </div>
    
    <div class="space-y-2">
      <Label>简介</Label>
      <Textarea bind:value={formData.description} placeholder="请输入图书简介" rows={3} />
    </div>
  </div>
  
  <DialogFooter>
    <Button variant="outline" on:click={() => dialogOpen = false}>取消</Button>
    <Button on:click={handleSave}>保存</Button>
  </DialogFooter>
</Dialog>
