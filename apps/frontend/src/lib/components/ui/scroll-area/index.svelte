<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { ScrollArea as ScrollAreaPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: ScrollAreaPrimitive.Props & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('relative overflow-hidden', className);
</script>

<ScrollAreaPrimitive.Root class={classes} {...restProps}>
  <ScrollAreaPrimitive.Viewport class="h-full w-full rounded-[inherit]">
    <slot />
  </ScrollAreaPrimitive.Viewport>
  <ScrollAreaPrimitive.Scrollbar
    orientation="vertical"
    class="flex touch-none select-none transition-colors data-[disabled=true]:opacity-0 data-[orientation=vertical]:h-full data-[orientation=vertical]:w-2.5 data-[orientation=horizontal]:h-2.5 data-[orientation=horizontal]:flex-col data-[orientation=horizontal]:w-full"
  >
    <ScrollAreaPrimitive.Thumb
      class="relative flex-1 rounded-full bg-border"
    />
  </ScrollAreaPrimitive.Scrollbar>
  <ScrollAreaPrimitive.Scrollbar
    orientation="horizontal"
    class="flex touch-none select-none transition-colors data-[disabled=true]:opacity-0 data-[orientation=vertical]:h-full data-[orientation=vertical]:w-2.5 data-[orientation=horizontal]:h-2.5 data-[orientation=horizontal]:flex-col data-[orientation=horizontal]:w-full"
  >
    <ScrollAreaPrimitive.Thumb
      class="relative flex-1 rounded-full bg-border"
    />
  </ScrollAreaPrimitive.Scrollbar>
  <ScrollAreaPrimitive.Corner />
</ScrollAreaPrimitive.Root>
