<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let disabled = false;
  export let className = '';
  
  const baseClasses = "relative inline-flex items-center justify-center px-6 py-2.5 text-sm font-semibold rounded-xl transition-all duration-300 ease-out overflow-hidden focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-900 shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed transform hover:-translate-y-0.5 active:translate-y-0";
  
  const variants = {
    primary: "bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500 shadow-sm border border-transparent",
    secondary: "bg-white text-slate-700 hover:bg-slate-50 focus:ring-slate-500 border border-slate-300 shadow-sm",
    danger: "bg-red-50 text-red-600 hover:bg-red-100 focus:ring-red-500 border border-red-200"
  };
</script>

<button {type} {disabled} class="{baseClasses} {variants[variant]} {className}" on:click>
  {#if variant === 'primary'}
    <div class="absolute inset-0 opacity-0 hover:opacity-10 transition-opacity duration-300 bg-black"></div>
  {/if}
  <span class="relative z-10 flex items-center space-x-2">
    <slot />
  </span>
</button>
