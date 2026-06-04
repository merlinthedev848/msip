<script lang="ts">
  import { onMount } from 'svelte';
  let trunks = [];
  let messages = [];

  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/sms`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        messages = (data.sms || []).map(m => ({
          id: m.ID, dir: m.Direction, from: m.Sender, to: m.Receiver, status: 'Sent', text: m.Body, time: m.CreatedAt
        }));
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">SMS Trunks & Routing</h1>
      <p class="text-slate-500 font-medium text-sm">Manage SMPP binds, message routing, and A2P compliance.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center">
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add SMPP Trunk
    </button>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <!-- Active Binds -->
    <div class="lg:col-span-2 bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl">
      <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
        <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"></path></svg>
          Active SMPP Binds
        </h3>
      </div>
      <table class="w-full text-left">
        <thead>
          <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
            <th class="p-4">Provider</th>
            <th class="p-4">Host</th>
            <th class="p-4 text-center">Binds</th>
            <th class="p-4 text-center">TPS</th>
            <th class="p-4">Latency</th>
            <th class="p-4 text-right">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-800/50">
          {#each trunks as trunk}
            <tr class="hover:bg-slate-100/30 transition-colors">
              <td class="p-4 font-bold text-slate-900">{trunk.name}</td>
              <td class="p-4 text-xs font-mono text-slate-500">{trunk.host}</td>
              <td class="p-4 text-center font-mono font-bold text-blue-600">{trunk.binds}</td>
              <td class="p-4 text-center font-mono font-bold text-emerald-400">{trunk.tps}</td>
              <td class="p-4 text-xs font-mono text-slate-500">{trunk.latency}</td>
              <td class="p-4 text-right">
                <span class="px-2 py-1 rounded text-[10px] font-bold uppercase tracking-widest {trunk.status === 'Connected' ? 'bg-emerald-50 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">
                  {trunk.status}
                </span>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Quick Stats -->
    <div class="space-y-6">
      <div class="bg-white/60 rounded-[2rem] border border-slate-200 p-6 flex flex-col justify-center items-center shadow-2xl relative overflow-hidden h-full">
        <div class="absolute -right-10 -top-10 w-40 h-40 bg-blue-50 rounded-full blur-3xl"></div>
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-2 z-10">Outbound Messages (24h)</p>
        <h2 class="text-5xl font-black text-slate-900 font-mono z-10">14,208</h2>
        <div class="mt-4 flex items-center text-emerald-400 z-10">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
          <span class="text-sm font-bold">+12.5%</span>
        </div>
      </div>
    </div>
  </div>

  <!-- Message Log -->
  <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl">
    <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
      <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest">Recent Message Log</h3>
      <div class="flex space-x-2">
        <button aria-label="Filter" class="p-2 bg-slate-100 rounded-lg text-slate-500 hover:text-slate-900 transition-colors border border-slate-200"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"></path></svg></button>
        <button aria-label="Refresh Logs" class="p-2 bg-slate-100 rounded-lg text-slate-500 hover:text-slate-900 transition-colors border border-slate-200"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg></button>
      </div>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
          <th class="p-4">Time</th>
          <th class="p-4">Type</th>
          <th class="p-4">From</th>
          <th class="p-4">To</th>
          <th class="p-4">Message Preview</th>
          <th class="p-4 text-right">Status</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each messages as msg}
          <tr class="hover:bg-slate-100/30 transition-colors group cursor-pointer">
            <td class="p-4 text-xs font-mono text-slate-500">{msg.time}</td>
            <td class="p-4">
              <span class="text-xs font-bold {msg.dir === 'Outbound' ? 'text-blue-600' : 'text-emerald-400'}">{msg.dir}</span>
            </td>
            <td class="p-4 font-mono font-bold text-slate-700 text-sm">{msg.from}</td>
            <td class="p-4 font-mono font-bold text-slate-700 text-sm">{msg.to}</td>
            <td class="p-4 text-sm text-slate-500 truncate max-w-xs">{msg.text}</td>
            <td class="p-4 text-right">
              <span class="px-2 py-1 rounded text-[10px] font-bold uppercase tracking-widest
                {msg.status === 'Delivered' || msg.status === 'Received' ? 'bg-emerald-50 text-emerald-400' : 
                 msg.status === 'Failed' ? 'bg-rose-500/10 text-rose-400' : 
                 'bg-amber-500/10 text-amber-400'}">
                {msg.status}
              </span>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

