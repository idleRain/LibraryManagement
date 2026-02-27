<script lang="ts">
  export let type: 'success' | 'warning' | 'danger' | 'info' | 'default' = 'default';
  export let size: 'sm' | 'md' = 'sm';

  $: classes = [
    'badge',
    `badge-${type}`,
    size === 'md' ? 'px-3 py-1 text-sm' : ''
  ].filter(Boolean).join(' ');
</script>

<span class={classes}>
  <slot />
</span>

<style>
  .badge-default {
    @apply bg-gray-100 text-gray-800;
  }
</style>
