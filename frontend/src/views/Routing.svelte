<script lang="ts">
  import { onMount } from 'svelte';
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';

  let routes = [];
  let isModalOpen = false;
  let newDID = '';
  let newTarget = '';
  let newType = 'EXTENSION';

  async function fetchRoutes() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/inbound-routes`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        routes = data.inbound_routes || [];
      }
    } catch (e) {}
  }

  onMount(fetchRoutes);

  async function handleCreateRoute(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/inbound-routes`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ DIDNumber: newDID, DestinationType: newType, DestinationTarget: newTarget })
      });
      if (res.ok) {
        isModalOpen = false;
        newDID = ''; newTarget = '';
        fetchRoutes();
      }
    } catch (e) {}
  }

  async function handleDeleteRoute(id) {
    if(!confirm("Delete this inbound route?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/inbound-routes/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchRoutes();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">Inbound Routing</h1>
      <p class="text-slate-500 text-sm">Map phone numbers (DIDs) to system destinations.</p>
    </div>
    <Button variant="primary" on:click={() => isModalOpen = true}>Add Route</Button>
  </header>

  <div class="grid grid-cols-1 gap-4">
    {#if routes.length === 0}
      <div class="p-8 text-center text-slate-500">No inbound routes configured.</div>
    {/if}
    {#each routes as route}
    <GlassCard hoverEffect={true} className="flex flex-col md:flex-row items-center justify-between p-5 gap-4">
      <div class="flex items-center space-x-6 w-full md:w-auto">
        <div class="flex-shrink-0 bg-blue-50 border border-indigo-500/30 p-3 rounded-xl">
          <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
        </div>
        <div>
          <h3 class="text-lg font-mono font-bold text-slate-900">{route.DIDNumber}</h3>
          <p class="text-sm text-slate-500">Inbound Route</p>
        </div>
      </div>
      
      <div class="hidden md:flex flex-1 items-center justify-center">
        <div class="h-px bg-slate-200 w-full max-w-[100px] relative">
          <div class="absolute right-0 top-1/2 -translate-y-1/2 w-2 h-2 border-t border-r border-slate-200 transform rotate-45"></div>
        </div>
      </div>

      <div class="flex items-center w-full md:w-auto bg-white border border-slate-200 rounded-lg p-3">
        <div class="w-2 h-2 rounded-full bg-emerald-500 mr-3"></div>
        <div>
          <p class="text-xs text-slate-500 font-medium">DESTINATION ({route.DestinationType})</p>
          <p class="text-sm font-semibold text-slate-800">{route.DestinationTarget}</p>
        </div>
      </div>
      
      <div class="flex ml-4">
        <button class="text-slate-500 hover:text-slate-900 p-2 transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg></button>
        <button class="text-slate-500 hover:text-red-400 p-2 transition-colors" on:click={() => handleDeleteRoute(route.ID)}><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
      </div>
    </GlassCard>
    {/each}
  </div>

</div>

<Modal bind:isOpen={isModalOpen} title="Create Inbound Route">
  <form on:submit={handleCreateRoute} class="space-y-4">
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">DID Number / Pattern</label>
      <input type="text" bind:value={newDID} required placeholder="+1234567890" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Destination Type</label>
      <select bind:value={newType} class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500">
        <option value="EXTENSION">Extension</option>
        <option value="IVR">IVR Menu</option>
        <option value="VOICEMAIL">Voicemail</option>
      </select>
    </div>
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Target ID / Extension</label>
      <input type="text" bind:value={newTarget} required placeholder="1000" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Add Route</button>
    </div>
  </form>
</Modal>

