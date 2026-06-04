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

  let queueStats = [
    { name: 'Support Tier 1', waiting: 4, agents: 8, sla: 92 },
    { name: 'Sales Global', waiting: 1, agents: 4, sla: 98 },
    { name: 'Billing', waiting: 0, agents: 2, sla: 100 }
  ];

  let interval;
  onMount(() => {
    connectWebSocket();
    interval = setInterval(() => {
      // simulate real-time updates for now if websocket is down
      if (!isConnected) {
        metrics.queue_length = Math.floor(Math.random() * 10);
      }
    }, 5000);
  });

  function connectWebSocket() {
    ws = new WebSocket(`ws://${window.location.hostname}:8081/api/wallboard/stream`);
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
    if (interval) clearInterval(interval);
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
      <h1 class="text-4xl font-black text-slate-900 tracking-tighter">Call Center Wallboard</h1>
      <p class="text-slate-500 font-medium">Real-time queue analytics and SLA tracking.</p>
    </div>
    <div class="text-right">
      <div class="text-3xl font-black font-mono text-slate-900 tracking-widest">{new Date().toLocaleTimeString('en-US', { hour12: false })}</div>
      <p class="text-emerald-400 text-xs font-bold uppercase tracking-[0.2em]">Live Connection</p>
    </div>
  </header>

  <!-- Massive KPIs -->
  <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-4 gap-6">
    
    <div class="bg-white/60 rounded-[2rem] border border-slate-200 p-8 flex flex-col justify-center items-center relative overflow-hidden group">
      <div class="absolute inset-0 bg-gradient-to-br from-rose-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
      <p class="text-sm font-bold text-slate-500 uppercase tracking-[0.3em] mb-4">Calls Waiting</p>
      <h2 class="text-8xl font-black {metrics.queue_length > 5 ? 'text-rose-500 animate-pulse' : 'text-slate-900'} font-mono transition-colors duration-500">
        {metrics.queue_length}
      </h2>
    </div>

    <div class="bg-white/60 rounded-[2rem] border border-slate-200 p-8 flex flex-col justify-center items-center relative overflow-hidden group">
      <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
      <p class="text-sm font-bold text-slate-500 uppercase tracking-[0.3em] mb-4">Active Agents</p>
      <h2 class="text-8xl font-black text-slate-900 font-mono">{metrics.agents_online}</h2>
    </div>

    <div class="bg-white/60 rounded-[2rem] border border-slate-200 p-8 flex flex-col justify-center items-center relative overflow-hidden group">
      <div class="absolute inset-0 bg-gradient-to-br from-amber-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
      <p class="text-sm font-bold text-slate-500 uppercase tracking-[0.3em] mb-4">Avg Wait Time</p>
      <h2 class="text-7xl font-black text-slate-900 font-mono">{formatTime(metrics.avg_wait_time)}</h2>
    </div>

    <div class="bg-white/60 rounded-[2rem] border border-slate-200 p-8 flex flex-col justify-center items-center relative overflow-hidden group">
      <div class="absolute inset-0 bg-gradient-to-br from-emerald-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
      <p class="text-sm font-bold text-slate-500 uppercase tracking-[0.3em] mb-4">Global SLA</p>
      <div class="flex items-baseline">
        <h2 class="text-7xl font-black {metrics.sla_percentage < 90 ? 'text-amber-500' : 'text-emerald-400'} font-mono">{(metrics.sla_percentage).toFixed(1)}</h2>
        <span class="text-3xl font-bold text-slate-500 ml-2">%</span>
      </div>
    </div>
  </div>

  <!-- Detailed Queues Table -->
  <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden mt-8">
    <table class="w-full text-left">
      <thead>
        <tr class="bg-slate-50/50 border-b border-slate-200">
          <th class="p-6 text-xs font-bold text-slate-500 uppercase tracking-[0.2em]">Queue Name</th>
          <th class="p-6 text-xs font-bold text-slate-500 uppercase tracking-[0.2em] text-center">Waiting</th>
          <th class="p-6 text-xs font-bold text-slate-500 uppercase tracking-[0.2em] text-center">Agents Ready</th>
          <th class="p-6 text-xs font-bold text-slate-500 uppercase tracking-[0.2em] text-right">SLA %</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each queueStats as queue}
          <tr class="hover:bg-slate-100/30 transition-colors">
            <td class="p-6 font-bold text-slate-900 text-xl">{queue.name}</td>
            <td class="p-6 text-center">
              <span class="text-3xl font-black font-mono {queue.waiting > 0 ? 'text-rose-400' : 'text-gray-600'}">{queue.waiting}</span>
            </td>
            <td class="p-6 text-center">
              <span class="text-2xl font-bold font-mono text-blue-600">{queue.agents}</span>
            </td>
            <td class="p-6 text-right">
              <div class="flex flex-col items-end">
                <span class="text-2xl font-black font-mono {queue.sla > 95 ? 'text-emerald-400' : 'text-amber-400'}">{queue.sla.toFixed(1)}%</span>
                <div class="w-32 h-1.5 bg-slate-100 rounded-full mt-2 overflow-hidden">
                  <div class="h-full bg-emerald-500 transition-all duration-500" style="width: {queue.sla}%"></div>
                </div>
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

