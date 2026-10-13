<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let rates = [];
  let isModalOpen = false;
  let newName = '';
  let newRate = 100;

  async function fetchRules() {
    try {
      const res = await fetch(`/api/v1/throttling`, {
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
      const res = await fetch(`/api/v1/throttling`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: newName, Rate: Number(newRate) })
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
      const res = await fetch(`/api/v1/throttling/${id}`, {
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
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Rate Throttling</h1>
      <p class="text-slate-500 font-medium text-sm">Monitor and limit SIP packet rates per IP to prevent DoS attacks.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Limit
    </button>
  </header>

  <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl">
    <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
      <h3 class="text-xs font-bold text-white uppercase tracking-wider flex items-center font-mono">
        <svg class="w-4 h-4 mr-2 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
        Active Rate Limiters (Packets / Second)
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] font-mono bg-[#0E121B]">
          <th class="p-4">Source IP</th>
          <th class="p-4">SIP Method</th>
          <th class="p-4 w-1/3">Traffic Load</th>
          <th class="p-4 text-right">Status</th>
          <th class="p-4 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[#1E2538]">
        {#if rates.length === 0}
          <tr>
            <td colspan="5" class="p-8 text-center text-slate-400 bg-[#121622] font-mono text-xs">No rate rules configured.</td>
          </tr>
        {/if}
        {#each rates as rate}
          <tr class="hover:bg-[#182030]/50 transition-colors group">
            <td class="p-4">
              <span class="font-bold text-white block font-mono text-sm">{rate.ip}</span>
              <span class="text-xs text-slate-500 font-mono">{rate.desc}</span>
            </td>
            <td class="p-4">
              <span class="px-2 py-0.5 bg-[#182030] border border-[#26334A] rounded text-xs font-mono font-bold text-cyan-400">{rate.type}</span>
            </td>
            <td class="p-4">
              <div class="flex items-center justify-between mb-1.5 font-mono text-xs">
                <span class="font-bold text-white">{rate.current} <span class="text-slate-500">pps</span></span>
                <span class="text-slate-400">Limit: {rate.max}</span>
              </div>
              <div class="w-full h-1.5 bg-[#182030] rounded-full overflow-hidden relative border border-[#26334A]">
                <div class="absolute h-full transition-all duration-1000 ease-out
                  {rate.status === 'Normal' ? 'bg-emerald-500' : 
                   rate.status === 'Warning' ? 'bg-amber-400' : 
                   'bg-rose-500'}" 
                  style="width: {Math.min((rate.current / rate.max) * 100, 100)}%">
                </div>
              </div>
            </td>
            <td class="p-4 text-right">
              <span class="px-2.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider font-mono
                {rate.status === 'Normal' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' : 
                 rate.status === 'Warning' ? 'bg-amber-950/60 text-amber-400 border border-amber-500/30' : 
                 'bg-rose-950/60 text-rose-400 border border-rose-500/30 animate-pulse'}">
                {rate.status}
              </span>
            </td>
            <td class="p-4 text-right">
              <button title="Delete Rule" aria-label="Delete Rule" class="text-rose-400 hover:text-white hover:bg-rose-600 p-1.5 rounded transition-colors" on:click={() => handleDeleteRule(rate.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
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
      <label for="throttle_ip" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Target IP / CIDR</label>
      <input id="throttle_ip" type="text" bind:value={newName} required placeholder="192.168.1.50" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="throttle_rate" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Max Packets/sec</label>
      <input id="throttle_rate" type="number" bind:value={newRate} required min="10" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Save Rule</button>
    </div>
  </form>
</Modal>
