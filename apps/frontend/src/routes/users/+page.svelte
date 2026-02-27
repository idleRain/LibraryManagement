<script lang="ts">
  import AppLayout from '$lib/components/layout/AppLayout.svelte';
  import { 
    Card, 
    CardContent,
    Button,
    Input,
    Table,
    TableHeader,
    TableBody,
    TableRow,
    TableHead,
    TableCell,
    Badge,
    Avatar,
    AvatarFallback
  } from '$lib/components/ui';
  import { Plus, Search, Edit, Key } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  let users = $state([
    { id: 1, username: 'admin', email: 'admin@library.com', real_name: '系统管理员', status: 1, roles: [{ name: '管理员' }] },
    { id: 2, username: 'zhangsan', email: 'zhangsan@example.com', real_name: '张三', status: 1, roles: [{ name: '普通用户' }] },
    { id: 3, username: 'lisi', email: 'lisi@example.com', real_name: '李四', status: 0, roles: [{ name: '普通用户' }] }
  ]);

  let searchUsername = $state('');

  function handleEdit(user: any) {
    toast.info('编辑用户: ' + user.username);
  }

  function handleResetPassword(user: any) {
    toast.success('已重置用户 ' + user.username + ' 的密码');
  }
</script>

<svelte:head>
  <title>用户管理 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">用户管理</h1>
        <p class="text-gray-500 mt-1">管理系统用户账户</p>
      </div>
      <Button>
        <Plus class="w-4 h-4 mr-2" />
        添加用户
      </Button>
    </div>

    <!-- Search -->
    <Card>
      <CardContent class="p-4">
        <div class="flex gap-4">
          <div class="flex-1">
            <Input placeholder="搜索用户名..." bind:value={searchUsername}>
              <Search class="w-4 h-4 text-gray-400" slot="icon" />
            </Input>
          </div>
          <Button variant="outline">搜索</Button>
        </div>
      </CardContent>
    </Card>

    <!-- Users Table -->
    <Card>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>用户</TableHead>
              <TableHead>邮箱</TableHead>
              <TableHead>角色</TableHead>
              <TableHead>状态</TableHead>
              <TableHead class="text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {#each users as user}
              <TableRow>
                <TableCell>
                  <div class="flex items-center gap-3">
                    <Avatar>
                      <AvatarFallback class="bg-primary text-white">
                        {(user.real_name || user.username).charAt(0).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <p class="font-medium">{user.real_name || user.username}</p>
                      <p class="text-sm text-gray-500">@{user.username}</p>
                    </div>
                  </div>
                </TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>
                  {#each user.roles as role}
                    <Badge variant="secondary">{role.name}</Badge>
                  {/each}
                </TableCell>
                <TableCell>
                  <Badge variant={user.status === 1 ? 'default' : 'destructive'}>
                    {user.status === 1 ? '启用' : '禁用'}
                  </Badge>
                </TableCell>
                <TableCell class="text-right">
                  <div class="flex items-center justify-end gap-2">
                    <Button variant="ghost" size="sm" on:click={() => handleEdit(user)}>
                      <Edit class="w-4 h-4" />
                    </Button>
                    <Button variant="ghost" size="sm" on:click={() => handleResetPassword(user)}>
                      <Key class="w-4 h-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            {/each}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</AppLayout>
