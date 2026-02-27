<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Tabs as TabsPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: TabsPrimitive.ListProps & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn(
    'inline-flex h-10 items-center justify-center rounded-md bg-muted p-1 text-muted-foreground',
    className
  );
</script>

<TabsPrimitive.List class={classes} {...restProps} />
