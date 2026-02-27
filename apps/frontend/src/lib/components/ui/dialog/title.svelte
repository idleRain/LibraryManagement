<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLHeadingElement> = $props();

  $: classes = cn('text-lg font-semibold leading-none tracking-tight', className);
</script>

<h2 class={classes} {...restProps} />
