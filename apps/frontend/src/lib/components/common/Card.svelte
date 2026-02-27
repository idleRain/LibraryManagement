<script lang="ts">
  export let title: string = '';
  export let subtitle: string = '';
  export let actions: boolean = false;
</script>

<div class="card">
  {#if title || actions}
    <div class="card-header flex items-center justify-between">
      <div>
        {#if title}
          <h3 class="text-lg font-semibold text-gray-800">{title}</h3>
        {/if}
        {#if subtitle}
          <p class="text-sm text-gray-500 mt-1">{subtitle}</p>
        {/if}
      </div>
      {#if actions}
        <div class="flex items-center gap-2">
          <slot name="actions" />
        </div>
      {/if}
    </div>
  {/if}
  <div class="card-body">
    <slot />
  </div>
</div>
