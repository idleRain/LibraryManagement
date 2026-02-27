<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('flex flex-1 flex-col overflow-y-auto p-4', className);
</script>

<div class={classes} {...restProps}>
  <slot />
</div>
