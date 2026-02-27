<script lang="ts">
  import { goto } from '$app/navigation';
  import { authApi } from '$lib/api';
  import { auth } from '$lib/stores';
  import { Button, Input, Label, Card, CardHeader, CardContent, CardTitle, CardDescription } from '$lib/components/ui';
  import { toast } from 'svelte-sonner';
  import { Book } from 'lucide-svelte';

  let username = $state('');
  let password = $state('');
  let loading = $state(false);
  let errors = $state<{ username?: string; password?: string }>({});

  async function handleLogin() {
    errors = {};
    
    if (!username) {
      errors.username = '请输入用户名';
    }
    if (!password) {
      errors.password = '请输入密码';
    }
    
    if (Object.keys(errors).length > 0) return;

    loading = true;
    
    try {
      const response = await authApi.login(username, password);
      
      if (response.code === 200 && response.data) {
        auth.login(
          response.data.token,
          response.data.user as any,
          response.data.permissions
        );
        toast.success('登录成功');
        goto('/');
      } else {
        toast.error(response.message || '登录失败');
      }
    } catch (error) {
      toast.error('网络错误，请稍后重试');
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<svelte:head>
  <title>登录 - 图书管理系统</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
  <div class="w-full max-w-md">
    <!-- Logo -->
    <div class="text-center mb-8">
      <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary text-white mb-4">
        <Book class="w-8 h-8" />
      </div>
      <h1 class="text-2xl font-bold text-gray-900">图书管理系统</h1>
      <p class="text-gray-500 mt-2">Library Management System</p>
    </div>

    <!-- Login Card -->
    <Card>
      <CardHeader>
        <CardTitle>登录账户</CardTitle>
        <CardDescription>请输入您的账户信息登录系统</CardDescription>
      </CardHeader>
      <CardContent>
        <form on:submit|preventDefault={handleLogin} class="space-y-4">
          <div class="space-y-2">
            <Label for="username">用户名</Label>
            <Input
              id="username"
              type="text"
              placeholder="请输入用户名"
              bind:value={username}
              error={errors.username}
              on:keypress={handleKeyPress}
            />
          </div>
          
          <div class="space-y-2">
            <Label for="password">密码</Label>
            <Input
              id="password"
              type="password"
              placeholder="请输入密码"
              bind:value={password}
              error={errors.password}
              on:keypress={handleKeyPress}
            />
          </div>

          <Button type="submit" class="w-full" {loading}>
            登录
          </Button>
        </form>

        <div class="mt-6 text-center text-sm text-gray-500">
          还没有账户？
          <a href="/auth/register" class="text-primary hover:underline">
            立即注册
          </a>
        </div>
      </CardContent>
    </Card>

    <!-- Demo Account -->
    <div class="mt-4 text-center text-sm text-gray-500">
      <p>演示账户：admin / admin123</p>
    </div>
  </div>
</div>
