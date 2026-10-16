<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let isOpen = false;
  export let open: boolean | undefined = undefined;
  export let title = '';
  export let description = '';
  export let width: 'md' | 'lg' = 'md';
  export let maxWidth: string = '';
  export let onClose: (() => void) | undefined = undefined;

  const dispatch = createEventDispatcher<{ close: void }>();

  $: active = open !== undefined ? open : isOpen;

  function close() {
    isOpen = false;
    if (open !== undefined) {
      open = false;
    }
    dispatch('close');
    onClose?.();
  }

  function onKey(e: KeyboardEvent) {
    if (active && e.key === 'Escape') close();
  }
</script>

<svelte:window on:keydown={onKey} />

{#if active}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="absolute inset-0 bg-slate-900/40" on:click={close}></div>
    <div
      role="dialog"
      aria-modal="true"
      aria-label={title}
      class="relative w-full {maxWidth || (width === 'lg' ? 'max-w-2xl' : 'max-w-lg')} bg-white border border-slate-200 rounded-lg shadow-xl max-h-[90vh] flex flex-col"
    >
      <div class="px-5 py-4 border-b border-slate-200 flex justify-between items-start">
        <div>
          <h2 class="text-base font-semibold text-slate-800">{title}</h2>
          {#if description}
            <p class="text-xs text-slate-500 mt-0.5">{description}</p>
          {/if}
        </div>
        <button type="button" aria-label="Close" title="Close" on:click={close} class="p-1 text-slate-400 hover:text-slate-700 rounded hover:bg-slate-100">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>
      <div class="px-5 py-4 overflow-y-auto text-sm text-slate-700">
        <slot />
      </div>
    </div>
  </div>
{/if}
