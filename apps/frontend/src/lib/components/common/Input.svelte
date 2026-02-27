<script lang="ts">
  export let type: string = 'text';
  export let value: string = '';
  export let placeholder: string = '';
  export let disabled: boolean = false;
  export let error: string = '';
  export let label: string = '';
  export let id: string = `input-${Math.random().toString(36).slice(2)}`;

  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    value = target.value;
  }
</script>

<div class="w-full">
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700 mb-1">
      {label}
    </label>
  {/if}
  <input
    {id}
    {type}
    {value}
    {placeholder}
    {disabled}
    class="input {error ? 'input-error' : ''}"
    on:input={handleInput}
    {...$$restProps}
  />
  {#if error}
    <p class="mt-1 text-sm text-red-600">{error}</p>
  {/if}
</div>
