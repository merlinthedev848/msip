<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let disabled = false;
  export let className = '';
  
  const baseClasses = "relative inline-flex items-center justify-center px-6 py-2.5 text-sm font-semibold rounded-xl transition-all duration-300 ease-out overflow-hidden focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-900 shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed transform hover:-translate-y-0.5 active:translate-y-0";
  
  const variants = {
    primary: "bg-indigo-600 text-white hover:bg-indigo-500 focus:ring-indigo-500 shadow-indigo-500/25 hover:shadow-indigo-500/40 border border-indigo-500/50",
    secondary: "bg-gray-800 text-gray-200 hover:bg-gray-700 hover:text-white focus:ring-gray-600 border border-gray-700 shadow-gray-900/50 hover:border-gray-600",
    danger: "bg-red-500/10 text-red-400 hover:bg-red-500/20 focus:ring-red-500 border border-red-500/30 hover:border-red-500/50"
  };
</script>

<button {type} {disabled} class="{baseClasses} {variants[variant]} {className}" on:click>
  {#if variant === 'primary'}
    <div class="absolute inset-0 bg-gradient-to-r from-indigo-600 to-violet-600 opacity-0 hover:opacity-100 transition-opacity duration-300"></div>
  {/if}
  <span class="relative z-10 flex items-center space-x-2">
    <slot />
  </span>
</button>
