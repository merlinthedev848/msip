<script lang="ts">
  import { onMount } from 'svelte';
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';

  let routes = [];
  let trunks = [];
  let isModalOpen = false;
  let isEditModalOpen = false;

  let newPattern = '';
  let newPriority = 10;
  let newTrunkId = '';

  let editRouteId = '';
  let editPattern = '';
  let editPriority = 10;
  let editTrunkId = '';

  async function fetchRoutes() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/outbound-routes`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        routes = data.outbound_routes || [];
      }
    } catch (e) {}
  }

  async function fetchTrunks() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/trunks`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        trunks = data.trunks || [];
        if (trunks.length > 0) {
          newTrunkId = trunks[0].ID;
          editTrunkId = trunks[0].ID;
        }
      }
    } catch (e) {}
  }

  onMount(() => {
    fetchRoutes();
    fetchTrunks();
  });

  async function handleCreateRoute(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/outbound-routes`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ DialPattern: newPattern, TrunkID: newTrunkId, Priority: parseInt(newPriority) })
      });
      if (res.ok) {
        isModalOpen = false;
        newPattern = '';
        fetchRoutes();
      }
    } catch (e) {}
  }

  function openEditModal(route) {
    editRouteId = route.ID;
    editPattern = route.DialPattern;
    editPriority = route.Priority;
    editTrunkId = route.TrunkID;
    isEditModalOpen = true;
  }

  async function handleEditRoute(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/outbound-routes/${editRouteId}`, {
        method: 'PUT',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ DialPattern: editPattern, TrunkID: editTrunkId, Priority: parseInt(editPriority) })
      });
      if (res.ok) {
        isEditModalOpen = false;
        fetchRoutes();
      } else {
        alert("Failed to update route.");
      }
    } catch (e) {
      alert("Error updating route.");
    }
  }

  async function handleDeleteRoute(id) {
    if(!confirm("Delete this outbound route?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/outbound-routes/${id}`, {
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
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">Outbound Routes</h1>
      <p class="text-slate-500 text-sm">Define dial patterns and prioritize carrier trunks.</p>
    </div>
    <Button variant="primary" on:click={() => isModalOpen = true}>New Rule</Button>
  </header>

  {#if routes.length === 0}
    <div class="p-8 text-center text-slate-500 bg-white/40 border border-slate-200 rounded-3xl">No outbound routes configured.</div>
  {/if}
  {#each routes as route}
  <GlassCard hoverEffect={true} className="border-l-4 border-l-indigo-500">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
      <div>
        <h3 class="text-lg font-bold text-slate-900 mb-1">Route (Priority {route.Priority})</h3>
        <p class="text-sm font-mono text-slate-500 bg-white inline-block px-2 py-1 rounded">Pattern: {route.DialPattern}</p>
      </div>
      <div class="mt-4 md:mt-0 flex items-center space-x-4">
        <div class="text-right">
          <p class="text-xs text-slate-500 uppercase">Trunk ID</p>
          <p class="text-sm font-semibold text-emerald-400 font-mono">{route.TrunkID ? route.TrunkID.substring(0,8) : 'N/A'}</p>
        </div>
        <Button variant="secondary" className="px-3 py-1" on:click={() => openEditModal(route)}>Edit</Button>
        <button title="Delete Route" aria-label="Delete Route" class="text-slate-500 hover:text-red-400 p-2 transition-colors ml-2" on:click={() => handleDeleteRoute(route.ID)}><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
      </div>
    </div>
  </GlassCard>
  {/each}
</div>

<Modal bind:isOpen={isModalOpen} title="Create Outbound Route">
  <form on:submit={handleCreateRoute} class="space-y-4">
    <div>
      <label for="new_pattern" class="block text-sm font-bold text-slate-500 mb-1">Dial Pattern</label>
      <input id="new_pattern" type="text" bind:value={newPattern} required placeholder="e.g. ^0[1-9]\d{8,9}$" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label for="new_trunk" class="block text-sm font-bold text-slate-500 mb-1">Select Trunk</label>
      <select id="new_trunk" bind:value={newTrunkId} required class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500">
        {#each trunks as trunk}
          <option value={trunk.ID}>{trunk.ProviderName} ({trunk.SIPServer})</option>
        {/each}
      </select>
    </div>
    <div>
      <label for="new_priority" class="block text-sm font-bold text-slate-500 mb-1">Priority</label>
      <input id="new_priority" type="number" bind:value={newPriority} required min="1" max="100" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Add Rule</button>
    </div>
  </form>
</Modal>

<Modal bind:isOpen={isEditModalOpen} title="Edit Outbound Route">
  <form on:submit={handleEditRoute} class="space-y-4">
    <div>
      <label for="edit_pattern" class="block text-sm font-bold text-slate-500 mb-1">Dial Pattern</label>
      <input id="edit_pattern" type="text" bind:value={editPattern} required placeholder="e.g. ^0[1-9]\d{8,9}$" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label for="edit_trunk" class="block text-sm font-bold text-slate-500 mb-1">Select Trunk</label>
      <select id="edit_trunk" bind:value={editTrunkId} required class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500">
        {#each trunks as trunk}
          <option value={trunk.ID}>{trunk.ProviderName} ({trunk.SIPServer})</option>
        {/each}
      </select>
    </div>
    <div>
      <label for="edit_priority" class="block text-sm font-bold text-slate-500 mb-1">Priority</label>
      <input id="edit_priority" type="number" bind:value={editPriority} required min="1" max="100" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isEditModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Save Changes</button>
    </div>
  </form>
</Modal>
