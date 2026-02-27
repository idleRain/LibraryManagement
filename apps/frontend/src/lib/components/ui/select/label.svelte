<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Select as SelectPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: SelectPrimitive.LabelProps & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn('py-1.5 pl-8 pr-2 text-sm font-semibold', className);
</script>

<SelectPrimitive.Label class={classes} {...restProps} />
