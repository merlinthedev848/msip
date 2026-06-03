<script lang="ts">
  import { onMount } from 'svelte';
  let domains = [];
  
  onMount(async () => {
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
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Multi-Tenant Domains</h1>
      <p class="text-gray-400 font-medium text-sm">Manage isolated tenant environments and resource limits.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center">
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Domain
    </button>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
    <div class="bg-gray-900/40 rounded-2xl border border-gray-800 p-6">
      <p class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-2">Total Tenants</p>
      <h2 class="text-4xl font-black text-white font-mono">{domains.length}</h2>
    </div>
    <div class="bg-gray-900/40 rounded-2xl border border-gray-800 p-6">
      <p class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-2">Global Extensions</p>
      <h2 class="text-4xl font-black text-white font-mono">{domains.reduce((a,b)=>a+b.extCount, 0)}</h2>
    </div>
    <div class="bg-gray-900/40 rounded-2xl border border-gray-800 p-6">
      <p class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-2">Concurrent Calls</p>
      <div class="flex items-baseline space-x-2">
        <h2 class="text-4xl font-black text-white font-mono">{domains.reduce((a,b)=>a+b.cc, 0)}</h2>
        <span class="text-sm font-bold text-emerald-500 bg-emerald-500/10 px-2 py-0.5 rounded">Live</span>
      </div>
    </div>
  </div>

  <div class="bg-gray-900/60 rounded-[2rem] border border-gray-800 overflow-hidden shadow-2xl">
    <div class="p-6 border-b border-gray-800 flex justify-between items-center bg-gray-950/30">
      <h3 class="text-lg font-bold text-white">Configured Domains</h3>
      <div class="relative">
        <svg class="w-4 h-4 text-gray-500 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        <input type="text" placeholder="Search tenants..." class="bg-gray-900 border border-gray-700 text-white text-sm rounded-xl focus:ring-indigo-500 focus:border-indigo-500 block w-64 pl-9 p-2 placeholder-gray-500 transition-colors">
      </div>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-gray-500 uppercase tracking-[0.2em] border-b border-gray-800">
          <th class="p-6">Domain Name</th>
          <th class="p-6 text-center">Extensions</th>
          <th class="p-6 text-center">Concurrent Limits</th>
          <th class="p-6">Status</th>
          <th class="p-6 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each domains as domain}
          <tr class="hover:bg-gray-800/30 transition-colors group">
            <td class="p-6">
              <div class="flex items-center space-x-4">
                <div class="w-10 h-10 rounded-xl bg-gray-800 flex items-center justify-center text-indigo-400 font-bold border border-gray-700">
                  {domain.name.charAt(0).toUpperCase()}
                </div>
                <div>
                  <p class="font-bold text-white text-base group-hover:text-indigo-300 transition-colors">{domain.name}</p>
                  <p class="text-xs text-gray-500">{domain.description}</p>
                </div>
              </div>
            </td>
            <td class="p-6">
              <div class="flex flex-col items-center">
                <span class="font-mono font-bold text-white text-lg">{domain.extCount} <span class="text-gray-600 text-sm">/ {domain.maxExt}</span></span>
                <div class="w-24 h-1.5 bg-gray-800 rounded-full mt-2 overflow-hidden">
                  <div class="h-full {domain.extCount/domain.maxExt > 0.8 ? 'bg-amber-500' : 'bg-indigo-500'}" style="width: {(domain.extCount/domain.maxExt)*100}%"></div>
                </div>
              </div>
            </td>
            <td class="p-6">
              <div class="flex flex-col items-center">
                <span class="font-mono font-bold text-emerald-400 text-lg">{domain.cc} <span class="text-gray-600 text-sm">/ {domain.maxCc}</span></span>
                <div class="w-24 h-1.5 bg-gray-800 rounded-full mt-2 overflow-hidden">
                  <div class="h-full {domain.cc/domain.maxCc > 0.8 ? 'bg-rose-500' : 'bg-emerald-500'}" style="width: {(domain.cc/domain.maxCc)*100}%"></div>
                </div>
              </div>
            </td>
            <td class="p-6">
              <span class="px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-widest
                {domain.status === 'Active' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 
                 domain.status === 'Warning' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 
                 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">
                {domain.status}
              </span>
            </td>
            <td class="p-6 text-right space-x-2">
              <button aria-label="Edit Domain" class="text-gray-500 hover:text-indigo-400 transition-colors p-2"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg></button>
              <button aria-label="Domain Settings" class="text-gray-500 hover:text-white transition-colors p-2"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

