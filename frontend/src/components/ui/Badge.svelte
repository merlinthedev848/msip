<script lang="ts">
  export let status: 'success' | 'warning' | 'error' | 'neutral' = 'neutral';
  export let pulse = false;
  
  const styles = {
    success: 'bg-emerald-100 text-emerald-800 border border-emerald-200',
    warning: 'bg-amber-100 text-amber-800 border border-amber-200',
    error: 'bg-red-100 text-red-800 border border-red-200',
    neutral: 'bg-slate-100 text-slate-800 border border-slate-200',
  };
  
  const dots = {
    success: 'bg-emerald-500',
    warning: 'bg-amber-500',
    error: 'bg-red-500',
    neutral: 'bg-slate-500',
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
