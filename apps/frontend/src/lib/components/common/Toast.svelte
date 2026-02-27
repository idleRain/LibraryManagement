<script lang="ts">
  import { toast } from '$lib/stores';
  import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-svelte';
  import { fly } from 'svelte/transition';

  const icons = {
    success: CheckCircle,
    error: AlertCircle,
    warning: AlertTriangle,
    info: Info
  };

  const colors = {
    success: 'bg-green-50 border-green-200 text-green-800',
    error: 'bg-red-50 border-red-200 text-red-800',
    warning: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    info: 'bg-blue-50 border-blue-200 text-blue-800'
  };

  const iconColors = {
    success: 'text-green-500',
    error: 'text-red-500',
    warning: 'text-yellow-500',
    info: 'text-blue-500'
  };
</script>

<div class="fixed top-4 right-4 z-[100] space-y-2">
  {#each $toast as t}
    <div
      class="flex items-center gap-3 px-4 py-3 rounded-lg border shadow-lg {colors[t.type]}"
      transition:fly={{ x: 100, duration: 300 }}
    >
      <svelte:component this={icons[t.type]} class="w-5 h-5 {iconColors[t.type]}" />
      <span class="text-sm font-medium">{t.message}</span>
      <button
        class="ml-2 p-1 rounded hover:bg-black hover:bg-opacity-10"
        on:click={() => toast.remove(t.id)}
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  {/each}
</div>
