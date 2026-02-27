<script lang="ts" generics="T extends Record<string, any>">
  export let columns: { key: string; label: string; width?: string }[] = [];
  export let data: T[] = [];
  export let loading: boolean = false;
  export let emptyText: string = '暂无数据';

  $: hasData = data && data.length > 0;
</script>

<div class="overflow-x-auto">
  <table class="table">
    <thead>
      <tr>
        {#each columns as col}
          <th style={col.width ? `width: ${col.width}` : ''}>
            {col.label}
          </th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#if loading}
        <tr>
          <td colspan={columns.length} class="text-center py-8">
            <div class="flex items-center justify-center gap-2">
              <svg class="animate-spin h-5 w-5 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-gray-500">加载中...</span>
            </div>
          </td>
        </tr>
      {:else if hasData}
        {#each data as row, i}
          <tr>
            {#each columns as col}
              <td>
                <slot name={`cell-${col.key}`} {row} {value: row[col.key]}>
                  {row[col.key] ?? '-'}
                </slot>
              </td>
            {/each}
          </tr>
        {/each}
      {:else}
        <tr>
          <td colspan={columns.length} class="text-center py-8 text-gray-500">
            {emptyText}
          </td>
        </tr>
      {/if}
    </tbody>
  </table>
</div>
