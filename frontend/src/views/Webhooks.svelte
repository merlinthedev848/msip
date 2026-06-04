<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let webhooks = [];
  let isModalOpen = false;
  let newUrl = '';
  let newEvent = 'call.answered';

  async function fetchWebhooks() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/webhooks`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/webhooks`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/webhooks/${id}`, {
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
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Event Webhooks</h1>
      <p class="text-slate-500 font-medium text-sm">Push real-time PBX events to your external web services.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Endpoint
    </button>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <!-- Webhooks List -->
    <div class="lg:col-span-2 bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl">
      <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
        <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest">Configured Endpoints</h3>
      </div>
      <table class="w-full text-left">
        <thead>
          <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
            <th class="p-4">Endpoint URL</th>
            <th class="p-4">Event Trigger</th>
            <th class="p-4 text-center">Last Response</th>
            <th class="p-4 text-right">Status</th>
            <th class="p-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-800/50">
          {#each webhooks as hook}
            <tr class="hover:bg-slate-100/30 transition-colors">
              <td class="p-4 font-mono font-bold text-blue-600 text-sm">{hook.url}</td>
              <td class="p-4 text-sm font-bold text-slate-700">{hook.trigger}</td>
              <td class="p-4 text-center">
                <span class="px-2 py-1 rounded text-xs font-mono font-bold {hook.lastCode === 200 ? 'bg-emerald-50 text-emerald-400' : 'bg-rose-500/10 text-rose-400'}">
                  HTTP {hook.lastCode}
                </span>
              </td>
              <td class="p-4 text-right">
                <span class="px-2 py-1 rounded text-[10px] font-bold uppercase tracking-widest {hook.status === 'Healthy' ? 'bg-emerald-50 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">
                  {hook.status}
                </span>
              </td>
              <td class="p-4 text-right">
                <button class="text-slate-500 hover:text-red-400 p-1" on:click={() => handleDeleteWebhook(hook.id)}>
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
      <div class="bg-white/60 rounded-[2rem] border border-slate-200 p-6 flex flex-col justify-center items-center shadow-2xl relative overflow-hidden h-full">
        <div class="absolute -left-10 -top-10 w-40 h-40 bg-emerald-50 rounded-full blur-3xl"></div>
        <p class="text-xs font-bold text-slate-500 uppercase tracking-widest mb-2 z-10">Events Delivered (24h)</p>
        <h2 class="text-5xl font-black text-slate-900 font-mono z-10">2,810</h2>
        <div class="mt-4 flex space-x-4 w-full px-4">
          <div class="flex-1 text-center">
            <p class="text-[10px] font-bold text-slate-500 uppercase">Success</p>
            <p class="text-lg font-mono font-bold text-emerald-400">99.8%</p>
          </div>
          <div class="flex-1 text-center border-l border-slate-200">
            <p class="text-[10px] font-bold text-slate-500 uppercase">Failed</p>
            <p class="text-lg font-mono font-bold text-rose-400">12</p>
          </div>
        </div>
      </div>
    </div>
  </div>
    </div>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create Webhook">
  <form on:submit={handleCreateWebhook} class="space-y-4">
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Target URL</label>
      <input type="url" bind:value={newUrl} required placeholder="https://api.yourdomain.com/webhook" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Event Type</label>
      <select bind:value={newEvent} class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500">
        <option value="call.answered">call.answered</option>
        <option value="call.ended">call.ended</option>
        <option value="sms.received">sms.received</option>
        <option value="extension.registered">extension.registered</option>
      </select>
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Create</button>
    </div>
  </form>
</Modal>



