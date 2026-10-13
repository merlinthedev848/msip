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
  <div class="fixed inset-0 z-50 flex items-center justify-center animate-in fade-in duration-200">
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="absolute inset-0 bg-black/75" on:click={close}></div>
    
    <!-- Modal Panel -->
    <div class="relative w-full max-w-lg bg-[#121622] border border-[#202738] rounded-xl shadow-2xl overflow-hidden animate-in zoom-in-95 duration-200">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#1E2538] bg-[#0E121B] flex justify-between items-start">
        <div>
          <h2 class="text-lg font-bold text-white font-mono">{title}</h2>
          {#if description}
            <p class="text-xs text-slate-400 mt-0.5 font-mono">{description}</p>
          {/if}
        </div>
        <button aria-label="Close" on:click={close} class="p-1.5 text-slate-400 hover:text-white transition-colors rounded hover:bg-[#1E2538]">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <!-- Body Slot -->
      <div class="px-6 py-5 text-slate-200">
        <slot></slot>
      </div>
    </div>
  </div>
{/if}
