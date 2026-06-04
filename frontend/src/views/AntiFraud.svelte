<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let blockedDestinations = [];
  let isModalOpen = false;
  let newPrefix = '';
  let newAction = 'BLOCK';

  async function fetchRules() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/fraud`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        blockedDestinations = (data.fraud || []).map(r => ({
          id: r.ID, country: r.BlockAction, code: r.Prefix, callsBlocked: 0, costSaved: '$0.00'
        }));
      }
    } catch (e) {}
  }

  onMount(fetchRules);

  async function handleCreateRule(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/fraud`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/fraud/${id}`, {
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
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Toll Fraud Prevention</h1>
      <p class="text-slate-500 font-medium text-sm">Block high-cost international destinations and rate centers.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Rule
    </button>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <div class="bg-white/40 rounded-2xl border border-slate-200 p-6 flex items-center">
      <div class="w-16 h-16 rounded-2xl bg-emerald-50 text-emerald-500 flex items-center justify-center mr-6 border border-emerald-500/20 shadow-sm">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-1">Est. Savings (30d)</p>
        <h2 class="text-3xl font-black text-slate-900 font-mono">$2,360.90</h2>
      </div>
    </div>
    <div class="bg-white/40 rounded-2xl border border-slate-200 p-6 flex items-center">
      <div class="w-16 h-16 rounded-2xl bg-rose-500/10 text-rose-500 flex items-center justify-center mr-6 border border-rose-500/20 shadow-sm">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path></svg>
      </div>
      <div>
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-1">Calls Blocked</p>
        <h2 class="text-3xl font-black text-slate-900 font-mono">1,407</h2>
      </div>
    </div>
    <div class="bg-white/40 rounded-2xl border border-slate-200 p-6 flex items-center">
      <div class="w-16 h-16 rounded-2xl bg-blue-50 text-indigo-500 flex items-center justify-center mr-6 border border-blue-100 shadow-sm">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-1">Strict Mode</p>
        <h2 class="text-xl font-black text-slate-900">Enabled</h2>
        <p class="text-[10px] text-slate-500">Whitelisting US/CAN/UK only</p>
      </div>
    </div>
  </div>

  <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl">
    <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
      <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
        <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        Blocked High-Cost Destinations
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
          <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
            <th class="p-4">Prefix / Code</th>
            <th class="p-4">Destination</th>
            <th class="p-4 text-center">Calls Blocked</th>
            <th class="p-4 text-right">Est. Savings</th>
            <th class="p-4 text-right">Actions</th>
          </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each blockedDestinations as dest}
          <tr class="hover:bg-slate-100/30 transition-colors group">
            <td class="p-6">
              <span class="font-bold text-slate-900 text-base">{dest.country}</span>
            </td>
            <td class="p-6 font-mono font-bold text-blue-600">{dest.code}</td>
            <td class="p-6 text-center font-mono font-bold text-slate-500">{dest.callsBlocked}</td>
              <td class="p-4 text-right text-emerald-400 font-bold font-mono text-sm">{dest.costSaved}</td>
              <td class="p-4 text-right space-x-2">
                <button class="text-slate-500 hover:text-slate-900 p-1 transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
                <button class="text-rose-500 hover:text-slate-900 hover:bg-rose-500 p-1 rounded transition-colors" on:click={() => handleDeleteRule(dest.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
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
      <label class="block text-sm font-bold text-slate-500 mb-1">Prefix / Pattern</label>
      <input type="text" bind:value={newPrefix} required placeholder="e.g. +900" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Action</label>
      <select bind:value={newAction} class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500">
        <option value="BLOCK">Block Calls</option>
        <option value="MONITOR">Monitor Only</option>
      </select>
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Save Rule</button>
    </div>
  </form>
</Modal>



