<script lang="ts">
  export let className = '';
  export let hoverEffect = false;
</script>

<div class="bg-gray-800/40 border border-gray-700/50 rounded-2xl p-6 backdrop-blur-md shadow-lg relative overflow-hidden {hoverEffect ? 'hover:border-gray-600 hover:shadow-xl transition-all duration-300 group' : ''} {className}">
  {#if hoverEffect}
    <div class="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"></div>
  {/if}
  <div class="relative z-10 h-full">
    <slot />
  </div>
</div>
