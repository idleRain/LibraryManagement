<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Select as SelectPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: SelectPrimitive.SeparatorProps & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('-mx-1 my-1 h-px bg-muted', className);
</script>

<SelectPrimitive.Separator class={classes} {...restProps} />
