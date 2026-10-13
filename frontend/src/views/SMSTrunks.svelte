<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let trunks = [];
  let messages = [];
  let isModalOpen = false;
  let newSmsTo = '';
  let newSmsBody = '';

  async function fetchSMS() {
    try {
      const res = await fetch(`/api/v1/sms`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        messages = (data.sms || []).map(m => ({
          id: m.ID, dir: m.Direction, from: m.Sender, to: m.Receiver, status: 'Sent', text: m.Body, time: m.CreatedAt
        }));
      }
    } catch (e) {}
  }

  onMount(fetchSMS);

  async function handleSendSMS(e) {
    e.preventDefault();
    try {
      const res = await fetch(`/api/v1/sms`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Direction: 'outbound', Sender: 'System', Receiver: newSmsTo, Body: newSmsBody })
      });
      if (res.ok) {
        isModalOpen = false;
        newSmsTo = ''; newSmsBody = '';
        fetchSMS();
      }
    } catch (e) {}
  }

  async function handleDeleteSMS(id) {
    if(!confirm("Delete this message?")) return;
    try {
      const res = await fetch(`/api/v1/sms/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchSMS();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">SMS Trunks & Routing</h1>
      <p class="text-slate-500 font-medium text-sm">Manage SMPP binds, message routing, and A2P compliance.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Send Message
    </button>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <!-- Active Binds -->
    <div class="lg:col-span-2 bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl">
      <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
        <h3 class="text-xs font-bold text-white uppercase tracking-wider flex items-center font-mono">
          <svg class="w-4 h-4 mr-2 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"></path></svg>
          Active SMPP Gateway Binds
        </h3>
      </div>
      <table class="w-full text-left">
        <thead>
          <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] font-mono bg-[#0E121B]">
            <th class="p-4">Provider</th>
            <th class="p-4">Host</th>
            <th class="p-4 text-center">Binds</th>
            <th class="p-4 text-center">TPS</th>
            <th class="p-4">Latency</th>
            <th class="p-4 text-right">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#1E2538]">
          {#each trunks as trunk}
            <tr class="hover:bg-[#182030]/50 transition-colors">
              <td class="p-4 font-bold text-white text-sm font-mono">{trunk.name}</td>
              <td class="p-4 text-xs font-mono text-cyan-400">{trunk.host}</td>
              <td class="p-4 text-center font-mono font-bold text-cyan-300">{trunk.binds}</td>
              <td class="p-4 text-center font-mono font-bold text-emerald-400">{trunk.tps}</td>
              <td class="p-4 text-xs font-mono text-slate-400">{trunk.latency}</td>
              <td class="p-4 text-right">
                <span class="px-2.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider font-mono {trunk.status === 'Connected' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' : 'bg-rose-950/60 text-rose-400 border border-rose-500/30'}">
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
      <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex flex-col justify-center items-center shadow-2xl h-full">
        <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-2 font-mono">Outbound Messages (24h)</p>
        <h2 class="text-4xl font-black text-white font-mono">14,208</h2>
        <div class="mt-3 flex items-center text-emerald-400 font-mono text-xs">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
          <span class="font-bold">+12.5% throughput</span>
        </div>
      </div>
    </div>
  </div>

<Modal bind:isOpen={isModalOpen} title="Send SMS Message">
  <form on:submit={handleSendSMS} class="space-y-4">
    <div>
      <label for="sms_recipient" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Recipient Number</label>
      <input id="sms_recipient" type="text" bind:value={newSmsTo} required placeholder="+1234567890" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="sms_body" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Message Body</label>
      <textarea id="sms_body" bind:value={newSmsBody} required rows="3" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none"></textarea>
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Send</button>
    </div>
  </form>
</Modal>

  <!-- Message Log -->
  <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl">
    <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
      <h3 class="text-xs font-bold text-white uppercase tracking-wider font-mono">Recent Message Log</h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] font-mono bg-[#0E121B]">
          <th class="p-4">Time</th>
          <th class="p-4">Type</th>
          <th class="p-4">From</th>
          <th class="p-4">To</th>
          <th class="p-4">Message Preview</th>
          <th class="p-4 text-right">Status</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[#1E2538]">
        {#each messages as msg}
          <tr class="hover:bg-[#182030]/50 transition-colors group">
            <td class="p-4 text-xs font-mono text-slate-400">{msg.time}</td>
            <td class="p-4">
              <span class="text-xs font-bold font-mono {msg.dir === 'Outbound' ? 'text-cyan-400' : 'text-emerald-400'}">{msg.dir}</span>
            </td>
            <td class="p-4 font-mono font-bold text-white text-xs">{msg.from}</td>
            <td class="p-4 font-mono font-bold text-white text-xs">{msg.to}</td>
            <td class="p-4 text-xs text-slate-300 font-mono truncate max-w-xs">{msg.text}</td>
            <td class="p-4 text-right">
              <span class="px-2.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider font-mono
                {msg.status === 'Delivered' || msg.status === 'Received' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' : 
                 msg.status === 'Failed' ? 'bg-rose-950/60 text-rose-400 border border-rose-500/30' : 
                 'bg-amber-950/60 text-amber-400 border border-amber-500/30'}">
                {msg.status}
              </span>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
