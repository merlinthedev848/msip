<script lang="ts">
  export let className = '';
  export let hoverEffect = false;
</script>

<div class="bg-slate-900/70 backdrop-blur-xl border border-slate-800/90 rounded-2xl p-6 shadow-2xl shadow-black/40 relative overflow-hidden transition-all duration-300 {hoverEffect ? 'hover:border-indigo-500/40 hover:shadow-indigo-500/10 hover:shadow-2xl' : ''} {className}">
  {#if hoverEffect}
    <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 via-purple-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"></div>
  {/if}
  <div class="relative z-10 h-full">
    <slot />
  </div>
</div>
