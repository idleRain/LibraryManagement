<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Separator as SeparatorPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    orientation = 'horizontal',
    decorative = true,
    class: className,
    ...restProps
  }: SeparatorPrimitive.Props & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('shrink-0 bg-border', orientation === 'horizontal' ? 'h-[1px] w-full' : 'h-full w-[1px]', className);
</script>

<SeparatorPrimitive.Root
  {orientation}
  {decorative}
  class={classes}
  {...restProps}
/>
