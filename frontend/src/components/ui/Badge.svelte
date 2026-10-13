<script lang="ts">
  export let status: 'success' | 'warning' | 'error' | 'neutral' = 'neutral';
  export let pulse = false;
  
  const styles = {
    success: 'bg-emerald-500/15 text-emerald-300 border border-emerald-500/30',
    warning: 'bg-amber-500/15 text-amber-300 border border-amber-500/30',
    error: 'bg-rose-500/15 text-rose-300 border border-rose-500/30',
    neutral: 'bg-slate-800/80 text-slate-300 border border-slate-700/80',
  };
  
  const dots = {
    success: 'bg-emerald-400',
    warning: 'bg-amber-400',
    error: 'bg-rose-400',
    neutral: 'bg-slate-400',
  };
</script>

<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {styles[status]} shadow-sm">
  {#if pulse}
    <span class="mr-1.5 flex h-2 w-2 relative">
      <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 {dots[status]}"></span>
      <span class="relative inline-flex rounded-full h-2 w-2 {dots[status]}"></span>
    </span>
  {:else}
    <span class="mr-1.5 h-1.5 w-1.5 rounded-full {dots[status]}"></span>
  {/if}
  <slot />
</span>
