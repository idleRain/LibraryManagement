<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    inset,
    ...restProps
  }: DropdownMenuPrimitive.LabelProps & HTMLAttributes<HTMLDivElement> & { inset?: boolean } = $props();

  $: classes = cn('px-2 py-1.5 text-sm font-semibold', inset && 'pl-8', className);
</script>

<DropdownMenuPrimitive.Label class={classes} {...restProps} />
