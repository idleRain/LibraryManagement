<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    inset,
    ...restProps
  }: DropdownMenuPrimitive.ItemProps & HTMLAttributes<HTMLDivElement> & { inset?: boolean } = $props();

  $: classes = cn(
    'relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    inset && 'pl-8',
    className
  );
</script>

<DropdownMenuPrimitive.Item class={classes} {...restProps} />
