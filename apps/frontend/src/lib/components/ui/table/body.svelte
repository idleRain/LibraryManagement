<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLTableSectionElement> = $props();

  $: classes = cn('[&_tr:last-child]:border-0', className);
</script>

<tbody class={classes} {...restProps} />
