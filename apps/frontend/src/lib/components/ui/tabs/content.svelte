<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Tabs as TabsPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: TabsPrimitive.ContentProps & HTMLAttributes<HTMLDivElement> = $props();

  $: classes = cn(
    'mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
    className
  );
</script>

<TabsPrimitive.Content class={classes} {...restProps} />
