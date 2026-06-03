<script lang="ts">
  import { onMount } from 'svelte';
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Button from '../components/ui/Button.svelte';

  let routes = [];
  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/outbound-routes`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        routes = data.outbound_routes || [];
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-white mb-1">Outbound Routes</h1>
      <p class="text-gray-400 text-sm">Define dial patterns and prioritize carrier trunks.</p>
    </div>
    <Button variant="primary">New Rule</Button>
  </header>

  {#if routes.length === 0}
    <div class="p-8 text-center text-gray-400">No outbound routes configured.</div>
  {/if}
  {#each routes as route}
  <GlassCard hoverEffect={true} className="border-l-4 border-l-indigo-500">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
      <div>
        <h3 class="text-lg font-bold text-white mb-1">Route (Priority {route.Priority})</h3>
        <p class="text-sm font-mono text-gray-400 bg-gray-900/50 inline-block px-2 py-1 rounded">Pattern: {route.DialPattern}</p>
      </div>
      <div class="mt-4 md:mt-0 flex items-center space-x-4">
        <div class="text-right">
          <p class="text-xs text-gray-500 uppercase">Trunk ID</p>
          <p class="text-sm font-semibold text-emerald-400 font-mono">{route.TrunkID.substring(0,8)}</p>
        </div>
        <Button variant="secondary" className="px-3 py-1">Edit</Button>
      </div>
    </div>
  </GlassCard>
  {/each}
</div>

