<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn(
    'flex h-full w-64 flex-col bg-sidebar text-sidebar-foreground',
    className
  );
</script>

<div class={classes} {...restProps}>
  <slot />
</div>
