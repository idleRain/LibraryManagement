<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLTableSectionElement> = $props();

  $: classes = cn('[&_tr]:border-b', className);
</script>

<thead class={classes} {...restProps} />
