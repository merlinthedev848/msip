<script lang="ts">
  export let title = 'Nothing here yet';
  export let message = '';
  export let description = '';
  export let loading = false;
  export let error = '';
</script>

<div class="py-12 px-6 text-center">
  {#if loading}
    <div class="inline-block w-6 h-6 border-2 border-slate-200 border-t-[#0284C7] rounded-full animate-spin"></div>
    <p class="text-xs text-slate-500 mt-3">Loading…</p>
  {:else if error}
    <p class="text-sm font-medium text-rose-700">Could not load data</p>
    <p class="text-xs text-slate-500 mt-1">{error}</p>
    <div class="mt-4"><slot name="retry" /></div>
  {:else}
    <p class="text-sm font-medium text-slate-700">{title}</p>
    {#if message || description}<p class="text-xs text-slate-500 mt-1">{message || description}</p>{/if}
    <div class="mt-4"><slot /></div>
  {/if}
</div>
