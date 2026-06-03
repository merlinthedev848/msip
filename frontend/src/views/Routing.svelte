<script lang="ts">
  import { onMount } from 'svelte';
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Button from '../components/ui/Button.svelte';

  let routes = [];
  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/inbound-routes`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        routes = data.inbound_routes || [];
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-white mb-1">Inbound Routing</h1>
      <p class="text-gray-400 text-sm">Map phone numbers (DIDs) to system destinations.</p>
    </div>
    <Button variant="primary">Add Route</Button>
  </header>

  <div class="grid grid-cols-1 gap-4">
    {#if routes.length === 0}
      <div class="p-8 text-center text-gray-400">No inbound routes configured.</div>
    {/if}
    {#each routes as route}
    <GlassCard hoverEffect={true} className="flex flex-col md:flex-row items-center justify-between p-5 gap-4">
      <div class="flex items-center space-x-6 w-full md:w-auto">
        <div class="flex-shrink-0 bg-indigo-500/10 border border-indigo-500/30 p-3 rounded-xl">
          <svg class="w-6 h-6 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
        </div>
        <div>
          <h3 class="text-lg font-mono font-bold text-white">{route.DIDNumber}</h3>
          <p class="text-sm text-gray-400">Inbound Route</p>
        </div>
      </div>
      
      <div class="hidden md:flex flex-1 items-center justify-center">
        <div class="h-px bg-gray-700 w-full max-w-[100px] relative">
          <div class="absolute right-0 top-1/2 -translate-y-1/2 w-2 h-2 border-t border-r border-gray-700 transform rotate-45"></div>
        </div>
      </div>

      <div class="flex items-center w-full md:w-auto bg-gray-900/50 border border-gray-700 rounded-lg p-3">
        <div class="w-2 h-2 rounded-full bg-emerald-500 mr-3"></div>
        <div>
          <p class="text-xs text-gray-500 font-medium">DESTINATION ({route.DestinationType})</p>
          <p class="text-sm font-semibold text-gray-200">{route.DestinationTarget}</p>
        </div>
      </div>
      
      <div class="flex ml-4">
        <Button variant="secondary" className="px-3 py-1.5 text-xs">Edit</Button>
      </div>
    </GlassCard>
    {/each}
  </div>

</div>

