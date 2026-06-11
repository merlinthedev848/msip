<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let rates = [];
  let isModalOpen = false;
  let newName = '';
  let newRate = 100;

  async function fetchRules() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/throttling`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        rates = (data.throttling || []).map(r => ({
          id: r.ID, ip: r.Name, desc: 'API Rule', type: 'ALL', current: 0, max: r.Rate, status: 'Normal'
        }));
      }
    } catch (e) {}
  }

  onMount(fetchRules);

  async function handleCreateRule(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/throttling`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: newName, Rate: parseInt(newRate) })
      });
      if (res.ok) {
        isModalOpen = false;
        newName = '';
        fetchRules();
      }
    } catch (e) {}
  }

  async function handleDeleteRule(id) {
    if(!confirm("Delete this rule?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/throttling/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchRules();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Rate Throttling</h1>
      <p class="text-slate-500 font-medium text-sm">Monitor and limit SIP packet rates per IP to prevent DoS attacks.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Limit
    </button>
  </header>

  <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl">
    <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
      <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
        <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
        Active Limits (Packets/sec)
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
          <th class="p-6">Source IP</th>
          <th class="p-6">SIP Method</th>
          <th class="p-6 w-1/3">Traffic Load</th>
          <th class="p-6 text-right">Status</th>
          <th class="p-6 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#if rates.length === 0}
          <tr>
            <td colspan="5" class="p-8 text-center text-slate-500 bg-white/40">No rate rules configured.</td>
          </tr>
        {/if}
        {#each rates as rate}
          <tr class="hover:bg-slate-100/30 transition-colors group">
            <td class="p-6">
              <span class="font-bold text-slate-900 block">{rate.ip}</span>
              <span class="text-xs text-slate-500">{rate.desc}</span>
            </td>
            <td class="p-6">
              <span class="px-2 py-1 bg-slate-100 border border-slate-200 rounded text-xs font-mono font-bold text-blue-600">{rate.type}</span>
            </td>
            <td class="p-6">
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-bold text-slate-900">{rate.current} <span class="text-slate-500">pps</span></span>
                <span class="text-xs font-bold text-slate-500">Limit: {rate.max}</span>
              </div>
              <div class="w-full h-2 bg-slate-100 rounded-full overflow-hidden relative">
                <div class="absolute h-full transition-all duration-1000 ease-out
                  {rate.status === 'Normal' ? 'bg-emerald-500' : 
                   rate.status === 'Warning' ? 'bg-amber-500' : 
                   'bg-rose-500'}" 
                  style="width: {Math.min((rate.current / rate.max) * 100, 100)}%">
                </div>
              </div>
            </td>
            <td class="p-6 text-right">
              <span class="px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-widest
                {rate.status === 'Normal' ? 'bg-emerald-50 text-emerald-400 border border-emerald-500/20' : 
                 rate.status === 'Warning' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 
                 'bg-rose-500/10 text-rose-400 border border-rose-500/20 animate-pulse'}">
                {rate.status}
              </span>
            </td>
            <td class="p-6 text-right">
              <button title="Settings" aria-label="Settings" class="text-slate-500 hover:text-slate-900 p-2 transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
              <button title="Delete Rule" aria-label="Delete Rule" class="text-rose-500 hover:text-slate-900 hover:bg-rose-500 p-2 rounded-lg transition-colors" on:click={() => handleDeleteRule(rate.id)}><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Add Throttling Rule">
  <form on:submit={handleCreateRule} class="space-y-4">
    <div>
      <label for="throttle_ip" class="block text-sm font-bold text-slate-500 mb-1">Target IP / CIDR</label>
      <input id="throttle_ip" type="text" bind:value={newName} required placeholder="192.168.1.50" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label for="throttle_rate" class="block text-sm font-bold text-slate-500 mb-1">Max Packets/sec</label>
      <input id="throttle_rate" type="number" bind:value={newRate} required min="10" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Save Rule</button>
    </div>
  </form>
</Modal>
