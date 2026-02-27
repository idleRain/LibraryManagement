<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: DropdownMenuPrimitive.SeparatorProps & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('-mx-1 my-1 h-px bg-muted', className);
</script>

<DropdownMenuPrimitive.Separator class={classes} {...restProps} />
