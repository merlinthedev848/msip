<script lang="ts">
  export let status: 'success' | 'warning' | 'error' | 'neutral' = 'neutral';
  export let pulse = false;
  
  const styles = {
    success: 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20',
    warning: 'bg-amber-500/10 text-amber-400 border border-amber-500/20',
    error: 'bg-red-500/10 text-red-400 border border-red-500/20',
    neutral: 'bg-gray-500/10 text-gray-400 border border-gray-500/20',
  };
  
  const dots = {
    success: 'bg-emerald-400',
    warning: 'bg-amber-400',
    error: 'bg-red-400',
    neutral: 'bg-gray-400',
  };
</script>

<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {styles[status]} shadow-sm backdrop-blur-sm">
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
