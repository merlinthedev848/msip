<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let webhooks = [];
  let isModalOpen = false;
  let newUrl = '';
  let newEvent = 'call.answered';

  async function fetchWebhooks() {
    try {
      const res = await fetch(`/api/v1/webhooks`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        webhooks = (data.webhooks || []).map(w => ({
          id: w.ID, url: w.TargetUrl, trigger: w.Event, status: 'Healthy', events: 0, lastCode: 200
        }));
      }
    } catch (e) {}
  }

  onMount(fetchWebhooks);

  async function handleCreateWebhook(e) {
    e.preventDefault();
    try {
      const res = await fetch(`/api/v1/webhooks`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ TargetUrl: newUrl, Event: newEvent })
      });
      if (res.ok) {
        isModalOpen = false;
        newUrl = '';
        fetchWebhooks();
      }
    } catch (e) {}
  }

  async function handleDeleteWebhook(id) {
    if(!confirm("Delete this webhook?")) return;
    try {
      const res = await fetch(`/api/v1/webhooks/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchWebhooks();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Event Webhooks</h1>
      <p class="text-slate-500 font-medium text-sm">Push real-time PBX events to your external web services.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Endpoint
    </button>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <!-- Webhooks List -->
    <div class="lg:col-span-2 bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl">
      <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
        <h3 class="text-xs font-bold text-white uppercase tracking-wider font-mono">Configured Endpoints</h3>
      </div>
      <table class="w-full text-left">
        <thead>
          <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] font-mono bg-[#0E121B]">
            <th class="p-4">Endpoint URL</th>
            <th class="p-4">Event Trigger</th>
            <th class="p-4 text-center">Last Response</th>
            <th class="p-4 text-right">Status</th>
            <th class="p-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#1E2538]">
          {#if webhooks.length === 0}
            <tr>
              <td colspan="5" class="p-8 text-center text-slate-400 bg-[#121622] font-mono text-xs">No endpoints configured.</td>
            </tr>
          {/if}
          {#each webhooks as hook}
            <tr class="hover:bg-[#182030]/50 transition-colors">
              <td class="p-4 font-mono font-bold text-cyan-400 text-xs">{hook.url}</td>
              <td class="p-4 text-xs font-mono font-bold text-slate-300">{hook.trigger}</td>
              <td class="p-4 text-center">
                <span class="px-2 py-0.5 rounded text-xs font-mono font-bold {hook.lastCode === 200 ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' : 'bg-rose-950/60 text-rose-400 border border-rose-500/30'}">
                  HTTP {hook.lastCode}
                </span>
              </td>
              <td class="p-4 text-right">
                <span class="px-2.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider font-mono {hook.status === 'Healthy' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' : 'bg-rose-950/60 text-rose-400 border border-rose-500/30'}">
                  {hook.status}
                </span>
              </td>
              <td class="p-4 text-right">
                <button title="Delete Webhook" aria-label="Delete Webhook" class="text-rose-400 hover:text-white hover:bg-rose-600 p-1 rounded transition-colors" on:click={() => handleDeleteWebhook(hook.id)}>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Stats -->
    <div class="space-y-6">
      <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex flex-col justify-center items-center shadow-2xl h-full">
        <p class="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-2 font-mono">Events Delivered (24h)</p>
        <h2 class="text-4xl font-black text-white font-mono">2,810</h2>
        <div class="mt-4 flex space-x-4 w-full px-4 border-t border-[#1E2538] pt-4 font-mono">
          <div class="flex-1 text-center">
            <p class="text-[10px] font-bold text-slate-400 uppercase">Success</p>
            <p class="text-base font-mono font-bold text-emerald-400">99.8%</p>
          </div>
          <div class="flex-1 text-center border-l border-[#1E2538]">
            <p class="text-[10px] font-bold text-slate-400 uppercase">Failed</p>
            <p class="text-base font-mono font-bold text-rose-400">12</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create Webhook">
  <form on:submit={handleCreateWebhook} class="space-y-4">
    <div>
      <label for="hook_url" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Target URL</label>
      <input id="hook_url" type="url" bind:value={newUrl} required placeholder="https://api.yourdomain.com/webhook" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="hook_event" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Event Type</label>
      <select id="hook_event" bind:value={newEvent} class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none">
        <option value="call.answered">call.answered</option>
        <option value="call.ended">call.ended</option>
        <option value="sms.received">sms.received</option>
        <option value="extension.registered">extension.registered</option>
      </select>
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Create</button>
    </div>
  </form>
</Modal>
