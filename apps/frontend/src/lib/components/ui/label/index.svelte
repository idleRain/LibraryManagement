<script lang="ts">
  import { cn } from '$lib/utils.js';
  import { Label as LabelPrimitive } from 'bits-ui';
  import { type HTMLLabelAttributes } from 'svelte/elements';

  let {
    class: className,
    ...restProps
  }: LabelPrimitive.Props & HTMLLabelAttributes = $props();

  $: classes = cn(
    'text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70',
    className
  );
</script>

<LabelPrimitive.Root class={classes} {...restProps} />
