<script lang="ts">
  import '../app.css';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores';
  import { onMount } from 'svelte';

  let { children, data } = $props();

  onMount(() => {
    // 检查登录状态
    if (!$auth.isAuthenticated) {
      goto('/auth/login');
    }
  });
</script>

{@render children?.()}
