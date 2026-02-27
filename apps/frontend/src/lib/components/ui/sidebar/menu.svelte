<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLUListElement> = $props();

  $: classes = cn('flex w-full flex-col gap-1', className);
</script>

<ul class={classes} {...restProps}>
  <slot />
</ul>
