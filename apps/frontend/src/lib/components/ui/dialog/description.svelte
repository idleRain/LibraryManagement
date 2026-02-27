<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: HTMLAttributes<HTMLParagraphElement> = $props();

  $: classes = cn('text-sm text-muted-foreground', className);
</script>

<p class={classes} {...restProps} />
