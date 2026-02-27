<script lang="ts">
  import { Toaster as Sonner, toast } from 'svelte-sonner';
  import { cn } from '$lib/utils.js';

  let {
    position = 'bottom-right',
    class: className,
    ...restProps
  }: Parameters<typeof Sonner>[0] = $props();
</script>

<Sonner
  {position}
  class={cn('toaster group', className)}
  toastOptions={{
    classes: {
      toast: 'group toast group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[.toaster]:shadow-lg',
      description: 'group-[.toast]:text-muted-foreground',
      actionButton: 'group-[.toast]:bg-primary group-[.toast]:text-primary-foreground',
      cancelButton: 'group-[.toast]:bg-muted group-[.toast]:text-muted-foreground',
      error: 'group-[.toaster]:border-red-500/50 group-[.toaster]:bg-red-50 group-[.toaster]:text-red-900',
      success: 'group-[.toaster]:border-green-500/50 group-[.toaster]:bg-green-50 group-[.toaster]:text-green-900',
      warning: 'group-[.toaster]:border-yellow-500/50 group-[.toaster]:bg-yellow-50 group-[.toaster]:text-yellow-900',
      info: 'group-[.toaster]:border-blue-500/50 group-[.toaster]:bg-blue-50 group-[.toaster]:text-blue-900'
    }
  }}
  {...restProps}
/>
