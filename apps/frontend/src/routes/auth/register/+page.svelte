<script lang="ts">
  import { goto } from '$app/navigation';
  import { authApi } from '$lib/api';
  import { Button, Input, Label, Card, CardHeader, CardContent, CardTitle, CardDescription } from '$lib/components/ui';
  import { toast } from 'svelte-sonner';
  import { Book, ArrowLeft } from 'lucide-svelte';

  let username = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let email = $state('');
  let realName = $state('');
  let loading = $state(false);
  let errors = $state<Record<string, string>>({});

  function validate(): boolean {
    errors = {};
    
    if (!username || username.length < 3) {
      errors.username = '用户名至少3个字符';
    }
    if (!password || password.length < 6) {
      errors.password = '密码至少6个字符';
    }
    if (password !== confirmPassword) {
      errors.confirmPassword = '两次密码不一致';
    }
    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = '请输入有效的邮箱地址';
    }
    
    return Object.keys(errors).length === 0;
  }

  async function handleRegister() {
    if (!validate()) return;

    loading = true;
    
    try {
      const response = await authApi.register({
        username,
        password,
        email,
        real_name: realName
      });
      
      if (response.code === 200) {
        toast.success('注册成功，请登录');
        goto('/auth/login');
      } else {
        toast.error(response.message || '注册失败');
      }
    } catch (error) {
      toast.error('网络错误，请稍后重试');
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      handleRegister();
    }
  }
</script>

<svelte:head>
  <title>注册 - 图书管理系统</title>
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

    <!-- Register Card -->
    <Card>
      <CardHeader>
        <CardTitle>注册账户</CardTitle>
        <CardDescription>创建新账户以使用系统</CardDescription>
      </CardHeader>
      <CardContent>
        <form on:submit|preventDefault={handleRegister} class="space-y-4">
          <div class="space-y-2">
            <Label for="username">用户名 *</Label>
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
            <Label for="email">邮箱 *</Label>
            <Input
              id="email"
              type="email"
              placeholder="请输入邮箱"
              bind:value={email}
              error={errors.email}
              on:keypress={handleKeyPress}
            />
          </div>
          
          <div class="space-y-2">
            <Label for="realName">真实姓名</Label>
            <Input
              id="realName"
              type="text"
              placeholder="请输入真实姓名"
              bind:value={realName}
              on:keypress={handleKeyPress}
            />
          </div>
          
          <div class="space-y-2">
            <Label for="password">密码 *</Label>
            <Input
              id="password"
              type="password"
              placeholder="请输入密码"
              bind:value={password}
              error={errors.password}
              on:keypress={handleKeyPress}
            />
          </div>
          
          <div class="space-y-2">
            <Label for="confirmPassword">确认密码 *</Label>
            <Input
              id="confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              bind:value={confirmPassword}
              error={errors.confirmPassword}
              on:keypress={handleKeyPress}
            />
          </div>

          <Button type="submit" class="w-full" {loading}>
            注册
          </Button>
        </form>

        <div class="mt-6 text-center">
          <a href="/auth/login" class="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-primary">
            <ArrowLeft class="w-4 h-4" />
            返回登录
          </a>
        </div>
      </CardContent>
    </Card>
  </div>
</div>
