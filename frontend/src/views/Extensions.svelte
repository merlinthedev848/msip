<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';

  import { onMount } from 'svelte';
  
  let extensions = [];
  let isApiConnected = true;
  
  let isModalOpen = false;
  let newExtNumber = '';
  let newExtPassword = '';
  let newExtConfirm = '';
  let isSubmitting = false;

  async function fetchExtensions() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/extensions`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        extensions = (data.extensions || []).map(e => ({
          id: e.ID,
          ext: e.ExtensionNumber,
          name: 'DB Endpoint',
          status: e.IsActive ? 'Online' : 'Offline',
          ip: 'Dynamic',
          code: e.ProvisioningCode || 'N/A'
        }));
        isApiConnected = true;
      } else { throw new Error("API Offline"); }
    } catch (e) {
      isApiConnected = false;
      extensions = [];
    }
  }

  onMount(() => {
    fetchExtensions();
  });

  async function handleCreateExtension(e) {
    e.preventDefault();
    if (newExtPassword !== newExtConfirm) {
      alert("Passwords do not match!");
      return;
    }
    isSubmitting = true;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/extensions`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({
          ExtensionNumber: newExtNumber,
          PasswordHash: newExtPassword,
          IsActive: true
        })
      });
      
      if (res.ok) {
        isModalOpen = false;
        newExtNumber = '';
        newExtPassword = '';
        newExtConfirm = '';
        await fetchExtensions();
      } else {
        alert("Failed to create extension.");
      }
    } catch (err) {
      alert("API Error: Make sure backend is running.");
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDeleteExtension(id) {
    if(!confirm("Are you sure you want to delete this extension?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/extensions/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchExtensions();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">Extensions</h1>
      <p class="text-slate-500 text-sm">Manage SIP devices and user endpoints.</p>
    </div>
    <Button variant="primary" on:click={() => isModalOpen = true}>
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      <span>New Extension</span>
    </Button>
  </header>

  {#if !isApiConnected}
  <div class="mb-6 p-4 bg-yellow-500/10 border border-yellow-500/30 rounded-lg flex items-start space-x-3 text-yellow-500 text-sm">
    <svg class="w-5 h-5 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
    <div>
      <p class="font-bold">Backend API Offline</p>
      <p class="opacity-80 mt-1">The Go backend (port 8081) is not responding. Showing simulated data instead of PostgreSQL records.</p>
    </div>
  </div>
  {/if}

  <GlassCard className="p-0 overflow-hidden">
    <div class="p-4 border-b border-slate-200/50 bg-white/20 flex items-center">
      <svg class="w-5 h-5 text-slate-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
      <input type="text" placeholder="Search extensions..." class="bg-transparent border-none text-slate-700 focus:outline-none w-full placeholder-gray-600 text-sm" />
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-white border-b border-slate-200">
            <th class="py-4 px-6 text-xs font-semibold text-slate-500 uppercase tracking-wider">Ext</th>
            <th class="py-4 px-6 text-xs font-semibold text-slate-500 uppercase tracking-wider">Status</th>
            <th class="py-4 px-6 text-xs font-semibold text-slate-500 uppercase tracking-wider">Provisioning Code</th>
            <th class="py-4 px-6 text-xs font-semibold text-slate-500 uppercase tracking-wider text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-800/50 text-sm text-slate-700">
          {#each extensions as ext}
            <tr class="hover:bg-slate-100/30 transition-colors group cursor-pointer">
              <td class="p-4 font-mono text-blue-600 font-medium">{ext.ext}</td>
              <td class="p-4">
                <span class="px-2.5 py-1 rounded-md text-xs font-medium border {ext.status === 'Online' ? 'bg-emerald-50 text-emerald-400 border-emerald-500/20' : 'bg-gray-500/10 text-slate-500 border-gray-500/20'}">
                  {ext.status}
                </span>
              </td>
              <td class="p-4">
                <span class="font-mono text-sm text-blue-600 tracking-widest bg-blue-50 px-2 py-1 rounded border border-blue-100">
                  {ext.code}
                </span>
              </td>
              <td class="p-4 text-right flex justify-end space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button title="Edit Extension" aria-label="Edit Extension" class="text-slate-500 hover:text-slate-900 transition-colors p-1">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
                </button>
                <button title="Delete Extension" aria-label="Delete Extension" class="text-slate-500 hover:text-red-400 transition-colors p-1" on:click={() => handleDeleteExtension(ext.id)}>
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </GlassCard>
</div>

<!-- Add Extension Modal -->
<Modal bind:isOpen={isModalOpen} title="Provision New Extension" description="Create a new SIP extension in the PostgreSQL database. It will immediately be available to register on Kamailio/FreeSWITCH.">
  <form on:submit={handleCreateExtension} class="space-y-5">
    <div>
      <label for="extNumber" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Extension Number</label>
      <input id="extNumber" type="text" bind:value={newExtNumber} required placeholder="e.g. 1005" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none">
    </div>
    <div>
      <label for="extPassword" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">SIP Password</label>
      <input id="extPassword" type="password" bind:value={newExtPassword} required placeholder="Strong password" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none">
    </div>
    <div>
      <label for="extConfirm" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Confirm Password</label>
      <input id="extConfirm" type="password" bind:value={newExtConfirm} required placeholder="Confirm password" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none">
    </div>
    
    <div class="pt-4 flex justify-end space-x-3">
      <Button variant="secondary" on:click={() => isModalOpen = false} type="button">Cancel</Button>
      <button type="submit" disabled={isSubmitting} class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2.5 rounded-xl font-bold transition-all shadow-lg shadow-indigo-600/20 disabled:opacity-50">
        {isSubmitting ? 'Provisioning...' : 'Create Extension'}
      </button>
    </div>
  </form>
</Modal>



