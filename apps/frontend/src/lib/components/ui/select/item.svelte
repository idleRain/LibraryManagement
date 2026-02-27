<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Select as SelectPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';
  import { Check } from 'lucide-svelte';

  let {
    class: className,
    children,
    ...restProps
  }: SelectPrimitive.ItemProps & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn(
    'relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    className
  );
</script>

<SelectPrimitive.Item class={classes} {...restProps}>
  <span class="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
    <SelectPrimitive.ItemIndicator>
      <Check class="h-4 w-4" />
    </SelectPrimitive.ItemIndicator>
  </span>
  <SelectPrimitive.ItemText>
    {@render children?.()}
  </SelectPrimitive.ItemText>
</SelectPrimitive.Item>
