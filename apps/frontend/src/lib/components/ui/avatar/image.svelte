<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Avatar as AvatarPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: AvatarPrimitive.ImageProps & HTMLAttributes<HTMLImageElement> = $props();

  $: classes = cn('aspect-square h-full w-full', className);
</script>

<AvatarPrimitive.Image class={classes} {...restProps} />
