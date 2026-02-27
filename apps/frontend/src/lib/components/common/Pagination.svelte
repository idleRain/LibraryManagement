<script lang="ts">
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';

  export let total: number = 0;
  export let page: number = 1;
  export let pageSize: number = 10;
  export let showSizeChanger: boolean = true;
  export let pageSizeOptions: number[] = [10, 20, 50, 100];

  const dispatch = createEventDispatcher();

  $: totalPages = Math.ceil(total / pageSize);
  $: start = (page - 1) * pageSize + 1;
  $: end = Math.min(page * pageSize, total);

  function goToPage(p: number) {
    if (p >= 1 && p <= totalPages) {
      dispatch('change', { page: p, pageSize });
    }
  }

  function handlePageSizeChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    dispatch('change', { page: 1, pageSize: parseInt(target.value) });
  }

  function getPageNumbers(): (number | string)[] {
    const pages: (number | string)[] = [];
    
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);
      
      if (page > 3) {
        pages.push('...');
      }
      
      const start = Math.max(2, page - 1);
      const end = Math.min(totalPages - 1, page + 1);
      
      for (let i = start; i <= end; i++) {
        pages.push(i);
      }
      
      if (page < totalPages - 2) {
        pages.push('...');
      }
      
      pages.push(totalPages);
    }
    
    return pages;
  }
</script>

{#if total > 0}
  <div class="flex items-center justify-between px-4 py-3 bg-white border-t border-gray-100">
    <div class="flex items-center gap-4 text-sm text-gray-600">
      <span>显示 {start}-{end} 条，共 {total} 条</span>
      
      {#if showSizeChanger}
        <div class="flex items-center gap-2">
          <span>每页</span>
          <select
            class="border border-gray-300 rounded px-2 py-1 text-sm"
            value={pageSize}
            on:change={handlePageSizeChange}
          >
            {#each pageSizeOptions as size}
              <option value={size}>{size}</option>
            {/each}
          </select>
          <span>条</span>
        </div>
      {/if}
    </div>
    
    <div class="flex items-center gap-1">
      <button
        class="p-2 rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={page <= 1}
        on:click={() => goToPage(page - 1)}
      >
        <ChevronLeft class="w-4 h-4" />
      </button>
      
      {#each getPageNumbers() as p}
        {#if p === '...'}
          <span class="px-3 py-1 text-gray-400">...</span>
        {:else}
          <button
            class="min-w-[32px] h-8 px-3 rounded text-sm {page === p
              ? 'bg-primary-600 text-white'
              : 'hover:bg-gray-100 text-gray-700'}"
            on:click={() => goToPage(p as number)}
          >
            {p}
          </button>
        {/if}
      {/each}
      
      <button
        class="p-2 rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={page >= totalPages}
        on:click={() => goToPage(page + 1)}
      >
        <ChevronRight class="w-4 h-4" />
      </button>
    </div>
  </div>
{/if}
