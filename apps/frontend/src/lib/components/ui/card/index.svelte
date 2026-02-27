<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('rounded-xl border bg-card text-card-foreground shadow', className);
</script>

<div class={classes} {...restProps} />
