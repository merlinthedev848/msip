<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  
  let domains = [];
  let isModalOpen = false;
  let newName = '';
  let newDescription = '';
  
  async function fetchDomains() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/domains`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        domains = (data.domains || []).map(d => ({
          id: d.ID, name: d.Name, description: d.Domain, extCount: 0, maxExt: 50, cc: 0, maxCc: 20, status: 'Active'
        }));
      }
    } catch (e) {}
  }

  onMount(fetchDomains);

  async function handleAddDomain(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/domains`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: newName, Domain: newDescription })
      });
      if (res.ok) {
        isModalOpen = false;
        newName = '';
        newDescription = '';
        fetchDomains();
      }
    } catch (e) {}
  }

  async function handleDeleteDomain(id) {
    if(!confirm("Delete this domain?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/domains/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchDomains();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Multi-Tenant Domains</h1>
      <p class="text-slate-500 font-medium text-sm">Manage isolated tenant environments and resource limits.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Domain
    </button>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
    <div class="bg-white/40 rounded-2xl border border-slate-200 p-6">
      <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-2">Total Tenants</p>
      <h2 class="text-4xl font-black text-slate-900 font-mono">{domains.length}</h2>
    </div>
    <div class="bg-white/40 rounded-2xl border border-slate-200 p-6">
      <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-2">Global Extensions</p>
      <h2 class="text-4xl font-black text-slate-900 font-mono">{domains.reduce((a,b)=>a+b.extCount, 0)}</h2>
    </div>
    <div class="bg-white/40 rounded-2xl border border-slate-200 p-6">
      <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-2">Concurrent Calls</p>
      <div class="flex items-baseline space-x-2">
        <h2 class="text-4xl font-black text-slate-900 font-mono">{domains.reduce((a,b)=>a+b.cc, 0)}</h2>
        <span class="text-sm font-bold text-emerald-500 bg-emerald-50 px-2 py-0.5 rounded">Live</span>
      </div>
    </div>
  </div>

  <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl">
    <div class="p-6 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
      <h3 class="text-lg font-bold text-slate-900">Configured Domains</h3>
      <div class="relative">
        <svg class="w-4 h-4 text-slate-500 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        <input type="text" placeholder="Search tenants..." class="bg-white border border-slate-200 text-slate-900 text-sm rounded-xl focus:ring-indigo-500 focus:border-indigo-500 block w-64 pl-9 p-2 placeholder-gray-500 transition-colors">
      </div>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
          <th class="p-6">Domain Name</th>
          <th class="p-6 text-center">Extensions</th>
          <th class="p-6 text-center">Concurrent Limits</th>
          <th class="p-6">Status</th>
          <th class="p-6 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each domains as domain}
          <tr class="hover:bg-slate-100/30 transition-colors group">
            <td class="p-6">
              <div class="flex items-center space-x-4">
                <div class="w-10 h-10 rounded-xl bg-slate-100 flex items-center justify-center text-blue-600 font-bold border border-slate-200">
                  {domain.name.charAt(0).toUpperCase()}
                </div>
                <div>
                  <p class="font-bold text-slate-900 text-base group-hover:text-indigo-300 transition-colors">{domain.name}</p>
                  <p class="text-xs text-slate-500">{domain.description}</p>
                </div>
              </div>
            </td>
            <td class="p-6">
              <div class="flex flex-col items-center">
                <span class="font-mono font-bold text-slate-900 text-lg">{domain.extCount} <span class="text-gray-600 text-sm">/ {domain.maxExt}</span></span>
                <div class="w-24 h-1.5 bg-slate-100 rounded-full mt-2 overflow-hidden">
                  <div class="h-full {domain.extCount/domain.maxExt > 0.8 ? 'bg-amber-500' : 'bg-indigo-500'}" style="width: {(domain.extCount/domain.maxExt)*100}%"></div>
                </div>
              </div>
            </td>
            <td class="p-6">
              <div class="flex flex-col items-center">
                <span class="font-mono font-bold text-emerald-400 text-lg">{domain.cc} <span class="text-gray-600 text-sm">/ {domain.maxCc}</span></span>
                <div class="w-24 h-1.5 bg-slate-100 rounded-full mt-2 overflow-hidden">
                  <div class="h-full {domain.cc/domain.maxCc > 0.8 ? 'bg-rose-500' : 'bg-emerald-500'}" style="width: {(domain.cc/domain.maxCc)*100}%"></div>
                </div>
              </div>
            </td>
            <td class="p-6">
              <span class="px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-widest
                {domain.status === 'Active' ? 'bg-emerald-50 text-emerald-400 border border-emerald-500/20' : 
                 domain.status === 'Warning' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 
                 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">
                {domain.status}
              </span>
            </td>
            <td class="p-6 text-right space-x-2">
              <button aria-label="Edit Domain" class="text-slate-500 hover:text-blue-600 transition-colors p-2"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg></button>
              <button aria-label="Delete Domain" class="text-slate-500 hover:text-red-400 transition-colors p-2" on:click={() => handleDeleteDomain(domain.id)}><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create New Domain">
  <form on:submit={handleAddDomain} class="space-y-4">
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Domain Identity</label>
      <input type="text" bind:value={newName} required placeholder="e.g. Acme Corp" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">FQDN / Hostname</label>
      <input type="text" bind:value={newDescription} required placeholder="pbx.acmecorp.com" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Create Domain</button>
    </div>
  </form>
</Modal>


