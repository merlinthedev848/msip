<script lang="ts">
  import { onMount } from 'svelte';
  let webhooks = [];

  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/webhooks`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        webhooks = (data.webhooks || []).map(w => ({
          url: w.TargetUrl, trigger: w.Event, status: 'Healthy', events: 0, lastCode: 200
        }));
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Event Webhooks</h1>
      <p class="text-slate-500 font-medium text-sm">Push real-time PBX events to your external web services.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center">
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

