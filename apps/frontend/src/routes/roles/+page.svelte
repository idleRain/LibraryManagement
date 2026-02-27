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
  import { Plus, Shield, Key, Edit } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  let roles = $state([
    { id: 1, name: '管理员', code: 'admin', description: '系统管理员，拥有所有权限', permissions: 16 },
    { id: 2, name: '普通用户', code: 'user', description: '普通用户，拥有基本查看权限', permissions: 4 },
    { id: 3, name: '图书管理员', code: 'librarian', description: '图书管理员，管理图书和借阅', permissions: 8 }
  ]);

  let permissions = $state([
    { id: 1, name: '用户查看', code: 'user:read', resource: 'user', action: 'read' },
    { id: 2, name: '用户创建', code: 'user:create', resource: 'user', action: 'create' },
    { id: 3, name: '图书查看', code: 'book:read', resource: 'book', action: 'read' },
    { id: 4, name: '图书管理', code: 'book:manage', resource: 'book', action: 'manage' },
    { id: 5, name: '借阅查看', code: 'borrow:read', resource: 'borrow', action: 'read' },
    { id: 6, name: '借阅管理', code: 'borrow:manage', resource: 'borrow', action: 'manage' }
  ]);

  function handleEditRole(role: any) {
    toast.info('编辑角色: ' + role.name);
  }
</script>

<svelte:head>
  <title>角色权限 - 图书管理系统</title>
</svelte:head>

<AppLayout>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">角色权限</h1>
        <p class="text-gray-500 mt-1">管理系统角色和权限配置</p>
      </div>
      <Button>
        <Plus class="w-4 h-4 mr-2" />
        添加角色
      </Button>
    </div>

    <Tabs value="roles">
      <TabsList>
        <TabsTrigger value="roles">
          <Shield class="w-4 h-4 mr-2" />
          角色管理
        </TabsTrigger>
        <TabsTrigger value="permissions">
          <Key class="w-4 h-4 mr-2" />
          权限列表
        </TabsTrigger>
      </TabsList>

      <TabsContent value="roles">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>角色名称</TableHead>
                  <TableHead>角色编码</TableHead>
                  <TableHead>描述</TableHead>
                  <TableHead class="text-center">权限数</TableHead>
                  <TableHead class="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {#each roles as role}
                  <TableRow>
                    <TableCell class="font-medium">{role.name}</TableCell>
                    <TableCell>
                      <Badge variant="secondary">{role.code}</Badge>
                    </TableCell>
                    <TableCell>{role.description}</TableCell>
                    <TableCell class="text-center">{role.permissions}</TableCell>
                    <TableCell class="text-right">
                      <Button variant="ghost" size="sm" on:click={() => handleEditRole(role)}>
                        <Edit class="w-4 h-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                {/each}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="permissions">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>权限名称</TableHead>
                  <TableHead>权限编码</TableHead>
                  <TableHead>资源</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {#each permissions as perm}
                  <TableRow>
                    <TableCell class="font-medium">{perm.name}</TableCell>
                    <TableCell>
                      <Badge variant="secondary">{perm.code}</Badge>
                    </TableCell>
                    <TableCell>{perm.resource}</TableCell>
                    <TableCell>{perm.action}</TableCell>
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
