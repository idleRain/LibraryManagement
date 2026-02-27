<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: DropdownMenuPrimitive.TriggerProps & HTMLAttributes<HTMLButtonElement> = $props();

  $: classes = cn('outline-none', className);
</script>

<DropdownMenuPrimitive.Trigger class={classes} {...restProps} />
