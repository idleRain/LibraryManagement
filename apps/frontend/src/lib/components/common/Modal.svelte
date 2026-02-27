<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { X } from 'lucide-svelte';
  import { fly, fade } from 'svelte/transition';

  export let open: boolean = false;
  export let title: string = '';
  export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
  export let closable: boolean = true;

  const dispatch = createEventDispatcher();

  function close() {
    if (closable) {
      open = false;
      dispatch('close');
    }
  }

  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      close();
    }
  }

  $: sizeClass = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl'
  }[size];
</script>

{#if open}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4"
    transition:fade={{ duration: 150 }}
    on:click={handleBackdropClick}
    on:keydown={(e) => e.key === 'Escape' && close()}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black bg-opacity-50"></div>
    
    <!-- Modal -->
    <div
      class="relative bg-white rounded-xl shadow-xl w-full {sizeClass} max-h-[90vh] overflow-hidden"
      transition:fly={{ y: -20, duration: 200 }}
    >
      <!-- Header -->
      {#if title || closable}
        <div class="flex items-center justify-between px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-800">{title}</h3>
          {#if closable}
            <button
              class="p-1 rounded-lg hover:bg-gray-100 text-gray-400 hover:text-gray-600"
              on:click={close}
            >
              <X class="w-5 h-5" />
            </button>
          {/if}
        </div>
      {/if}
      
      <!-- Body -->
      <div class="px-6 py-4 overflow-y-auto max-h-[calc(90vh-120px)]">
        <slot />
      </div>
      
      <!-- Footer -->
      <div class="px-6 py-4 border-t border-gray-200 bg-gray-50">
        <slot name="footer" />
      </div>
    </div>
  </div>
{/if}
