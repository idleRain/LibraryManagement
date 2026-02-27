<script lang="ts">
  import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    open = $bindable(false),
    children,
    ...restProps
  }: DropdownMenuPrimitive.Props & HTMLAttributes<HTMLDivElement> = $props();
</script>

<DropdownMenuPrimitive.Root {open} {...restProps}>
  {@render children?.()}
</DropdownMenuPrimitive.Root>
