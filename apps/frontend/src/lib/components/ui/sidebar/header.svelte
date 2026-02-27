<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('flex h-16 items-center border-b px-6', className);
</script>

<div class={classes} {...restProps}>
  <slot />
</div>
