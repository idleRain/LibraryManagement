<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLTableElement> = $props();

  $: classes = cn('w-full caption-bottom text-sm', className);
</script>

<table class={classes} {...restProps} />
