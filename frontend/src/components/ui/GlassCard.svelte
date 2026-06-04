<script lang="ts">
  export let className = '';
  export let hoverEffect = false;
</script>

<div class="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm relative overflow-hidden {hoverEffect ? 'hover:border-blue-300 hover:shadow-md transition-all duration-300 group' : ''} {className}">
  {#if hoverEffect}
    <div class="absolute inset-0 bg-gradient-to-br from-blue-50 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"></div>
  {/if}
  <div class="relative z-10 h-full">
    <slot />
  </div>
</div>
