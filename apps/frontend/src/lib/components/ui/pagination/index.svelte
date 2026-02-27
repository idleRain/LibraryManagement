<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Pagination as PaginationPrimitive } from 'bits-ui';
  import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-svelte';
  import Button from '$lib/components/ui/button/index.svelte';
  import type { HTMLAttributes } from 'svelte/elements';

  interface Props extends HTMLAttributes<HTMLDivElement> {
    page: number;
    totalPages: number;
    onPageChange: (page: number) => void;
    showEdges?: boolean;
  }

  let {
    page,
    totalPages,
    onPageChange,
    showEdges = true,
    class: className,
    ...restProps
  }: Props = $props();

  function getPageNumbers(): (number | 'ellipsis')[] {
    const pages: (number | 'ellipsis')[] = [];
    
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);
      
      if (page > 3) {
        pages.push('ellipsis');
      }
      
      const start = Math.max(2, page - 1);
      const end = Math.min(totalPages - 1, page + 1);
      
      for (let i = start; i <= end; i++) {
        pages.push(i);
      }
      
      if (page < totalPages - 2) {
        pages.push('ellipsis');
      }
      
      pages.push(totalPages);
    }
    
    return pages;
  }

  $: pages = getPageNumbers();
</script>

<nav
  role="navigation"
  aria-label="pagination"
  class={cn('mx-auto flex w-full justify-center', className)}
  {...restProps}
>
  <ul class="flex flex-row items-center gap-1">
    {#if showEdges}
      <li>
        <Button
          variant="ghost"
          size="icon"
          disabled={page <= 1}
          on:click={() => onPageChange(page - 1)}
          aria-label="Go to previous page"
        >
          <ChevronLeft class="h-4 w-4" />
        </Button>
      </li>
    {/if}

    {#each pages as p}
      <li>
        {#if p === 'ellipsis'}
          <span class="flex h-9 w-9 items-center justify-center">
            <MoreHorizontal class="h-4 w-4" />
          </span>
        {:else}
          <Button
            variant={page === p ? 'default' : 'ghost'}
            size="icon"
            on:click={() => onPageChange(p)}
            aria-label="Go to page {p}"
            aria-current={page === p ? 'page' : undefined}
          >
            {p}
          </Button>
        {/if}
      </li>
    {/each}

    {#if showEdges}
      <li>
        <Button
          variant="ghost"
          size="icon"
          disabled={page >= totalPages}
          on:click={() => onPageChange(page + 1)}
          aria-label="Go to next page"
        >
          <ChevronRight class="h-4 w-4" />
        </Button>
      </li>
    {/if}
  </ul>
</nav>
