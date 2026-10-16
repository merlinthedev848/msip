<script lang="ts">
  export let title = '';
  export let description = '';
  export let subtitle = '';
</script>

<div class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-3 mb-5">
  <div>
    <h1 class="text-xl font-bold text-slate-800 tracking-tight">{title}</h1>
    {#if description || subtitle}<p class="text-xs text-slate-500 mt-1">{description || subtitle}</p>{/if}
  </div>
  <div class="flex items-center gap-2"><slot /></div>
</div>
