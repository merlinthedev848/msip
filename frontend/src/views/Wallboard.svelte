<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import GlassCard from '../components/ui/GlassCard.svelte';

  // Live WebSocket Metrics
  let metrics = {
    active_calls: 0,
    agents_online: 0,
    agents_busy: 0,
    queue_length: 0,
    avg_wait_time: 0,
    sla_percentage: 0,
    timestamp: 0
  };

  let ws: WebSocket;
  let isConnected = false;

  let queueStats = [];

  onMount(() => {
    connectWebSocket();
  });

  function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    ws = new WebSocket(`${protocol}://${window.location.host}/api/v1/wallboard/stream`);
    ws.onopen = () => { isConnected = true; };
    ws.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data);
        metrics = data.metrics || metrics;
        queueStats = data.queues || queueStats;
      } catch (err) {}
    };
    ws.onclose = () => { isConnected = false; setTimeout(connectWebSocket, 5000); };
  }

  onDestroy(() => {
    if (ws) ws.close();
  });

  $: formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  
  <header class="flex justify-between items-center mb-10">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2 flex items-center">
        Live Call Centre
        <span class="ml-3 px-2 py-0.5 rounded-full bg-blue-50 text-blue-600 border border-blue-100 text-xs font-bold uppercase tracking-wider flex items-center">
          <span class="w-1.5 h-1.5 rounded-full {isConnected ? 'bg-indigo-400 animate-pulse' : 'bg-rose-500'} mr-2"></span>
          {isConnected ? 'Stream Active' : 'Connecting...'}
        </span>
      </h1>
      <p class="text-slate-500 font-medium text-sm">Real-time queue and agent performance metrics.</p>
    </div>
    
    <button on:click={() => window.open('/wallboard-fullscreen', '_blank')} class="bg-slate-100 hover:bg-slate-200 text-white px-4 py-2 rounded-xl font-bold text-sm transition-colors flex items-center border border-slate-800">
      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
      Pop Out Display
    </button>
  </header>

  <!-- Massive KPIs -->
  <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-4 gap-6">
    
    <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex flex-col justify-center items-center shadow-2xl relative">
      <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 font-mono">Calls Waiting</p>
      <h2 class="text-7xl font-black {metrics.queue_length > 5 ? 'text-rose-400 animate-pulse' : 'text-white'} font-mono">
        {metrics.queue_length}
      </h2>
    </div>

    <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex flex-col justify-center items-center shadow-2xl relative">
      <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 font-mono">Active Agents</p>
      <h2 class="text-7xl font-black text-white font-mono">{metrics.agents_online}</h2>
    </div>

    <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex flex-col justify-center items-center shadow-2xl relative">
      <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 font-mono">Avg Wait Time</p>
      <h2 class="text-6xl font-black text-cyan-400 font-mono">{formatTime(metrics.avg_wait_time)}</h2>
    </div>

    <div class="bg-[#121622] rounded-xl border border-[#202738] p-6 flex flex-col justify-center items-center shadow-2xl relative">
      <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 font-mono">Global SLA</p>
      <div class="flex items-baseline font-mono">
        <h2 class="text-6xl font-black {metrics.sla_percentage < 90 ? 'text-amber-400' : 'text-emerald-400'}">{(metrics.sla_percentage).toFixed(1)}</h2>
        <span class="text-2xl font-bold text-slate-400 ml-1.5">%</span>
      </div>
    </div>
  </div>

  <!-- Detailed Queues Table -->
  <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden mt-8 shadow-2xl">
    <div class="p-4 border-b border-[#1E2538] bg-[#0E121B]">
      <h3 class="text-xs font-bold text-white uppercase tracking-wider font-mono">Queue Status & Live Performance</h3>
    </div>
    <table class="w-full text-left">
      <thead>
        <tr class="bg-[#0E121B] border-b border-[#1E2538] text-[10px] font-bold text-slate-400 uppercase tracking-wider font-mono">
          <th class="p-4">Queue Name</th>
          <th class="p-4 text-center">Waiting</th>
          <th class="p-4 text-center">Agents Ready</th>
          <th class="p-4 text-right">SLA %</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[#1E2538]">
        {#each queueStats as queue}
          <tr class="hover:bg-[#182030]/50 transition-colors font-mono">
            <td class="p-4 font-bold text-white text-base">{queue.name}</td>
            <td class="p-4 text-center">
              <span class="text-2xl font-black {queue.waiting > 0 ? 'text-rose-400' : 'text-slate-500'}">{queue.waiting}</span>
            </td>
            <td class="p-4 text-center">
              <span class="text-xl font-bold text-cyan-400">{queue.agents}</span>
            </td>
            <td class="p-4 text-right">
              <div class="flex flex-col items-end">
                <span class="text-lg font-bold {queue.sla > 95 ? 'text-emerald-400' : 'text-amber-400'}">{queue.sla.toFixed(1)}%</span>
                <div class="w-32 h-1.5 bg-[#182030] rounded-full mt-1.5 overflow-hidden border border-[#26334A]">
                  <div class="h-full bg-emerald-400 transition-all duration-500" style="width: {queue.sla}%"></div>
                </div>
              </div>
            </td>
          </tr>
        {:else}
          <tr>
            <td colspan="4" class="p-6 text-center text-slate-400 font-mono text-xs">
              {isConnected ? 'No active queues configured.' : 'Connecting to live statistics stream...'}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>


