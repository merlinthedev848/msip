<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
  export let isOpen = false;
  export let title = "Modal Title";
  export let description = "";

  const dispatch = createEventDispatcher();

  function close() {
    isOpen = false;
    dispatch('close');
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center animate-in fade-in duration-300">
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" on:click={close}></div>
    
    <!-- Modal Panel -->
    <div class="relative w-full max-w-lg bg-gray-900 border border-gray-700/50 rounded-2xl shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300 slide-in-from-bottom-8">
      <!-- Header -->
      <div class="px-6 py-5 border-b border-gray-800 bg-gray-950/50 flex justify-between items-start">
        <div>
          <h2 class="text-xl font-bold text-white">{title}</h2>
          {#if description}
            <p class="text-sm text-gray-400 mt-1">{description}</p>
          {/if}
        </div>
        <button aria-label="Close" on:click={close} class="p-2 text-gray-500 hover:text-white transition-colors rounded-lg hover:bg-gray-800">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <!-- Body Slot -->
      <div class="px-6 py-6 text-gray-300">
        <slot></slot>
      </div>
    </div>
  </div>
{/if}
