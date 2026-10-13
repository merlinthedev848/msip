<script lang="ts">
  export let className = '';
  export let hoverEffect = false;
</script>

<div class="bg-[#121622] border border-[#202738] rounded-xl p-5 shadow-md shadow-black/20 relative overflow-hidden transition-all duration-150 {hoverEffect ? 'hover:border-[#344057] hover:shadow-lg' : ''} {className}">
  <div class="relative z-10 h-full">
    <slot />
  </div>
</div>
