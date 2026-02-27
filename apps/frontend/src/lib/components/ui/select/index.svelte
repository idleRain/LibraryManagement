<script lang="ts">
  import { Select as SelectPrimitive } from 'bits-ui';

  let { value = $bindable(''), ...restProps }: SelectPrimitive.Props = $props();
</script>

<SelectPrimitive.Root {value} {...restProps}>
  <slot />
</SelectPrimitive.Root>
