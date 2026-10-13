<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let blockedDestinations = [];
  let isModalOpen = false;
  let newPrefix = '';
  let newAction = 'BLOCK';

  async function fetchRules() {
    try {
      const res = await fetch(`/api/v1/fraud`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        blockedDestinations = (data.fraud || []).map(r => ({
          id: r.ID, prefix: r.Prefix, action: r.BlockAction, callsBlocked: 0, costSaved: '$0.00'
        }));
      }
    } catch (e) {}
  }

  onMount(fetchRules);

  async function handleCreateRule(e) {
    e.preventDefault();
    try {
      const res = await fetch(`/api/v1/fraud`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Prefix: newPrefix, BlockAction: newAction })
      });
      if (res.ok) {
        isModalOpen = false;
        newPrefix = '';
        fetchRules();
      }
    } catch (e) {}
  }

  async function handleDeleteRule(id) {
    if(!confirm("Delete this fraud prevention rule?")) return;
    try {
      const res = await fetch(`/api/v1/fraud/${id}`, {
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
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Toll Fraud Prevention</h1>
      <p class="text-slate-500 font-medium text-sm">Block high-cost international destinations and rate centers.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Rule
    </button>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex items-center shadow-lg">
      <div class="w-14 h-14 rounded-xl bg-emerald-950/50 text-emerald-400 flex items-center justify-center mr-5 border border-emerald-500/30">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1 font-mono">Est. Savings (30d)</p>
        <h2 class="text-3xl font-black text-white font-mono">$0.00</h2>
      </div>
    </div>
    <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex items-center shadow-lg">
      <div class="w-14 h-14 rounded-xl bg-rose-950/50 text-rose-400 flex items-center justify-center mr-5 border border-rose-500/30">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path></svg>
      </div>
      <div>
        <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1 font-mono">Calls Blocked</p>
        <h2 class="text-3xl font-black text-white font-mono">0</h2>
      </div>
    </div>
    <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex items-center shadow-lg">
      <div class="w-14 h-14 rounded-xl bg-cyan-950/50 text-cyan-400 flex items-center justify-center mr-5 border border-cyan-500/30">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1 font-mono">Strict Mode</p>
        <h2 class="text-xl font-bold text-white font-mono">Disabled</h2>
        <p class="text-[10px] text-slate-500 font-mono">Standard monitoring only</p>
      </div>
    </div>
  </div>

  <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl">
    <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
      <h3 class="text-xs font-bold text-white uppercase tracking-wider flex items-center font-mono">
        <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        Blocked High-Cost Destinations
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
          <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] font-mono bg-[#0E121B]">
            <th class="p-4">Prefix / Pattern</th>
            <th class="p-4">Action</th>
            <th class="p-4 text-center">Calls Blocked</th>
            <th class="p-4 text-right">Est. Savings</th>
            <th class="p-4 text-right">Actions</th>
          </tr>
      </thead>
      <tbody class="divide-y divide-[#1E2538]">
        {#each blockedDestinations as dest}
          <tr class="hover:bg-[#182030]/50 transition-colors group">
            <td class="p-4">
              <span class="font-bold text-white text-sm font-mono">{dest.prefix}</span>
            </td>
            <td class="p-4 font-mono font-bold text-cyan-400 text-xs">{dest.action}</td>
            <td class="p-4 text-center font-mono font-bold text-slate-300 text-sm">{dest.callsBlocked}</td>
            <td class="p-4 text-right text-emerald-400 font-bold font-mono text-sm">{dest.costSaved}</td>
            <td class="p-4 text-right space-x-2">
              <button title="Delete Rule" aria-label="Delete Rule" class="text-rose-400 hover:text-white hover:bg-rose-600 p-1 rounded transition-colors" on:click={() => handleDeleteRule(dest.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Add Anti-Fraud Rule">
  <form on:submit={handleCreateRule} class="space-y-4">
    <div>
      <label for="prefix" class="block text-sm font-bold text-slate-500 mb-1">Prefix / Pattern</label>
      <input id="prefix" type="text" bind:value={newPrefix} required placeholder="e.g. +900" class="w-full bg-slate-950/80 border border-slate-800 text-slate-100 rounded-xl px-4 py-2 text-white focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label for="action" class="block text-sm font-bold text-slate-500 mb-1">Action</label>
      <select id="action" bind:value={newAction} class="w-full bg-slate-950/80 border border-slate-800 text-slate-100 rounded-xl px-4 py-2 text-white focus:ring-indigo-500 focus:border-indigo-500">
        <option value="BLOCK">Block Calls</option>
        <option value="MONITOR">Monitor Only</option>
      </select>
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-white" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Save Rule</button>
    </div>
  </form>
</Modal>



