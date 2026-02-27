<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  let {
    class: className,
    active = false,
    ...restProps
  }: HTMLAttributes<HTMLAnchorElement> & { active?: boolean } = $props();

  $: classes = cn(
    'flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors',
    active
      ? 'bg-sidebar-accent text-sidebar-accent-foreground font-medium'
      : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
    className
  );
</script>

<a class={classes} {...restProps}>
  <slot />
</a>
