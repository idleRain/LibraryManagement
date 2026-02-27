<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { type HTMLInputAttributes } from 'svelte/elements';

  interface Props extends HTMLInputAttributes {
    error?: string;
  }

  let {
    class: className,
    type = 'text',
    error,
    value = $bindable(''),
    ...restProps
  }: Props = $props();

  $: classes = cn(
    'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
    error && 'border-red-500 focus-visible:ring-red-500',
    className
  );
</script>

<div class="w-full">
  <input
    {type}
    {value}
    class={classes}
    {...restProps}
  />
  {#if error}
    <p class="mt-1 text-sm text-red-500">{error}</p>
  {/if}
</div>
