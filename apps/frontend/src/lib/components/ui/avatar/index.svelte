<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Avatar as AvatarPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: AvatarPrimitive.Props & HTMLAttributes<HTMLSpanElement> = $props();

  $: classes = cn('relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full', className);
</script>

<AvatarPrimitive.Root class={classes} {...restProps} />
