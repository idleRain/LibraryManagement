<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLHeadingElement> = $props();

  $: classes = cn('font-semibold leading-none tracking-tight', className);
</script>

<h3 class={classes} {...restProps} />
