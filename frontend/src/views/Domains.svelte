<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  
  let domains = [];
  let isModalOpen = false;
  let isEditModalOpen = false;
  
  let newName = '';
  let newDescription = '';
  
  let editDomainId = '';
  let editName = '';
  let editDescription = '';
  
  async function fetchDomains() {
    try {
      const res = await fetch(`/api/v1/domains`, {
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
      const res = await fetch(`/api/v1/domains`, {
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

  function openEditModal(domain) {
    editDomainId = domain.id;
    editName = domain.name;
    editDescription = domain.description;
    isEditModalOpen = true;
  }

  async function handleEditDomain(e) {
    e.preventDefault();
    try {
      const res = await fetch(`/api/v1/domains/${editDomainId}`, {
        method: 'PUT',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: editName, Domain: editDescription })
      });
      if (res.ok) {
        isEditModalOpen = false;
        fetchDomains();
      } else {
        alert("Failed to update domain.");
      }
    } catch (e) {
      alert("Error updating domain.");
    }
  }

  async function handleDeleteDomain(id) {
    if(!confirm("Delete this domain?")) return;
    try {
      const res = await fetch(`/api/v1/domains/${id}`, {
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
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Multi-Tenant Domains</h1>
      <p class="text-slate-500 font-medium text-sm">Manage isolated tenant environments and resource limits.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Domain
    </button>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
    <div class="bg-[#121622] rounded-xl border border-[#202738] p-5 shadow-lg">
      <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 font-mono">Total Tenants</p>
      <h2 class="text-3xl font-black text-white font-mono">{domains.length}</h2>
    </div>
    <div class="bg-[#121622] rounded-xl border border-[#202738] p-5 shadow-lg">
      <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 font-mono">Global Extensions</p>
      <h2 class="text-3xl font-black text-white font-mono">{domains.reduce((a,b)=>a+b.extCount, 0)}</h2>
    </div>
    <div class="bg-[#121622] rounded-xl border border-[#202738] p-5 shadow-lg">
      <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 font-mono">Concurrent Calls</p>
      <div class="flex items-baseline space-x-2">
        <h2 class="text-3xl font-black text-white font-mono">{domains.reduce((a,b)=>a+b.cc, 0)}</h2>
        <span class="text-xs font-bold text-emerald-400 bg-emerald-950/60 border border-emerald-500/30 px-2 py-0.5 rounded font-mono">LIVE</span>
      </div>
    </div>
  </div>

  <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl">
    <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
      <h3 class="text-xs font-bold text-white uppercase tracking-wider font-mono">Configured Tenant Domains</h3>
      <div class="relative">
        <svg class="w-4 h-4 text-slate-500 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        <input type="text" placeholder="Search tenants..." class="bg-[#0B0E17] border border-[#1E2538] text-white text-xs rounded-lg focus:border-cyan-500 block w-64 pl-9 p-2 placeholder-slate-500 outline-none font-mono">
      </div>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] font-mono bg-[#0E121B]">
          <th class="p-4">Domain Name</th>
          <th class="p-4 text-center">Extensions</th>
          <th class="p-4 text-center">Concurrent Limits</th>
          <th class="p-4">Status</th>
          <th class="p-4 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[#1E2538]">
        {#each domains as domain}
          <tr class="hover:bg-[#182030]/50 transition-colors group">
            <td class="p-4">
              <div class="flex items-center space-x-3">
                <div class="w-9 h-9 rounded-lg bg-[#182030] flex items-center justify-center text-cyan-400 font-bold border border-[#26334A] font-mono text-sm">
                  {domain.name.charAt(0).toUpperCase()}
                </div>
                <div>
                  <p class="font-bold text-white text-sm group-hover:text-cyan-300 transition-colors font-mono">{domain.name}</p>
                  <p class="text-xs text-slate-500 font-mono">{domain.description}</p>
                </div>
              </div>
            </td>
            <td class="p-4">
              <div class="flex flex-col items-center">
                <span class="font-mono font-bold text-white text-sm">{domain.extCount} <span class="text-slate-500 text-xs">/ {domain.maxExt}</span></span>
                <div class="w-24 h-1.5 bg-[#182030] rounded-full mt-1.5 overflow-hidden border border-[#26334A]">
                  <div class="h-full {domain.extCount/domain.maxExt > 0.8 ? 'bg-amber-400' : 'bg-cyan-500'}" style="width: {(domain.extCount/domain.maxExt)*100}%"></div>
                </div>
              </div>
            </td>
            <td class="p-4">
              <div class="flex flex-col items-center">
                <span class="font-mono font-bold text-emerald-400 text-sm">{domain.cc} <span class="text-slate-500 text-xs">/ {domain.maxCc}</span></span>
                <div class="w-24 h-1.5 bg-[#182030] rounded-full mt-1.5 overflow-hidden border border-[#26334A]">
                  <div class="h-full {domain.cc/domain.maxCc > 0.8 ? 'bg-rose-500' : 'bg-emerald-500'}" style="width: {(domain.cc/domain.maxCc)*100}%"></div>
                </div>
              </div>
            </td>
            <td class="p-4">
              <span class="px-2.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider font-mono
                {domain.status === 'Active' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' : 
                 domain.status === 'Warning' ? 'bg-amber-950/60 text-amber-400 border border-amber-500/30' : 
                 'bg-rose-950/60 text-rose-400 border border-rose-500/30'}">
                {domain.status}
              </span>
            </td>
            <td class="p-4 text-right space-x-2">
              <button aria-label="Edit Domain" class="text-slate-400 hover:text-cyan-400 transition-colors p-1.5" on:click={() => openEditModal(domain)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg></button>
              <button aria-label="Delete Domain" class="text-slate-400 hover:text-rose-400 transition-colors p-1.5" on:click={() => handleDeleteDomain(domain.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
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
      <label for="new_domain_name" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Domain Identity</label>
      <input id="new_domain_name" type="text" bind:value={newName} required placeholder="e.g. Acme Corp" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="new_domain_desc" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">FQDN / Hostname</label>
      <input id="new_domain_desc" type="text" bind:value={newDescription} required placeholder="pbx.acmecorp.com" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Create Domain</button>
    </div>
  </form>
</Modal>

<Modal bind:isOpen={isEditModalOpen} title="Edit Domain Settings">
  <form on:submit={handleEditDomain} class="space-y-4">
    <div>
      <label for="edit_domain_name" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Domain Identity</label>
      <input id="edit_domain_name" type="text" bind:value={editName} required placeholder="e.g. Acme Corp" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="edit_domain_desc" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">FQDN / Hostname</label>
      <input id="edit_domain_desc" type="text" bind:value={editDescription} required placeholder="pbx.acmecorp.com" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isEditModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Save Changes</button>
    </div>
  </form>
</Modal>



