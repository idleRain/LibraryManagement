<script lang="ts">
  import { cn } from '$lib/utils.js';
  import type { HTMLAttributes } from 'svelte/elements';

  type Variant = 'default' | 'secondary' | 'destructive' | 'outline';

  interface Props extends HTMLAttributes<HTMLDivElement> {
    variant?: Variant;
  }

  let {
    variant = 'default',
    class: className,
    ...restProps
  }: Props = $props();

  const variantClasses: Record<Variant, string> = {
    default: 'border-transparent bg-primary text-primary-foreground hover:bg-primary/80',
    secondary: 'border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80',
    destructive: 'border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80',
    outline: 'text-foreground'
  };

  $: classes = cn(
    'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
    variantClasses[variant],
    className
  );
</script>

<div class={classes} {...restProps}>
  <slot />
</div>
