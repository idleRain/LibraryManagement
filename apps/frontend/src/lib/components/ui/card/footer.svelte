<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('flex items-center p-6 pt-0', className);
</script>

<div class={classes} {...restProps} />
