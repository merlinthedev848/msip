<script lang="ts">
  export let title = '';
  export let subtitle = '';
  export let padded = true;
  export let className = '';
</script>

<section class="bg-white border border-slate-200 rounded-lg shadow-sm {className}">
  {#if title || $$slots.actions}
    <header class="px-5 py-3 border-b border-slate-200 flex items-center justify-between gap-3">
      <div class="min-w-0">
        {#if title}<h2 class="text-sm font-semibold text-slate-800 truncate">{title}</h2>{/if}
        {#if subtitle}<p class="text-xs text-slate-500 mt-0.5">{subtitle}</p>{/if}
      </div>
      <div class="flex items-center gap-2"><slot name="actions" /></div>
    </header>
  {/if}
  <div class={padded ? 'p-5' : ''}>
    <slot />
  </div>
</section>
