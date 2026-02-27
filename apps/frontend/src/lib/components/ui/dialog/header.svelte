<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('flex flex-col space-y-1.5 text-center sm:text-left', className);
</script>

<div class={classes} {...restProps} />
