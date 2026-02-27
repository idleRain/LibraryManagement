<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Avatar as AvatarPrimitive } from 'bits-ui';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: AvatarPrimitive.FallbackProps & HTMLAttributes<HTMLSpanElement> = $props();

  $: classes = cn('flex h-full w-full items-center justify-center rounded-full bg-muted', className);
</script>

<AvatarPrimitive.Fallback class={classes} {...restProps} />
