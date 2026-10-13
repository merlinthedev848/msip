<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { wsUrl } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Badge from '../components/ui/Badge.svelte';

  interface WallboardMetrics {
    active_calls: number;
    agents_online: number;
    agents_busy: number;
    queue_length: number;
    avg_wait_time: number;
    sla_percentage: number;
  }

  interface QueueStat {
    name: string;
    waiting: number;
    agents: number;
    sla: number;
  }

  let metrics: WallboardMetrics = {
    active_calls: 0,
    agents_online: 0,
    agents_busy: 0,
    queue_length: 0,
    avg_wait_time: 0,
    sla_percentage: 100
  };

  let queues: QueueStat[] = [];
  let isConnected = false;
  let ws: WebSocket | null = null;
  let reconnectTimer: any;

  function connectWS() {
    try {
      ws = new WebSocket(wsUrl('/wallboard/stream'));
      ws.onopen = () => {
        isConnected = true;
      };
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.metrics) {
            metrics = data.metrics;
          }
          if (data.queues) {
            queues = data.queues;
          }
        } catch (e) {}
      };
      ws.onclose = () => {
        isConnected = false;
        reconnectTimer = setTimeout(connectWS, 3000);
      };
      ws.onerror = () => {
        isConnected = false;
      };
    } catch (e) {
      isConnected = false;
      reconnectTimer = setTimeout(connectWS, 3000);
    }
  }

  function formatTime(seconds: number) {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  }

  onMount(() => {
    connectWS();
  });

  onDestroy(() => {
    if (reconnectTimer) clearTimeout(reconnectTimer);
    if (ws) ws.close();
  });
</script>

<PageHeader title="Call Queues & Live Wallboard" description="Real-time call center SLA metrics and agent queue performance.">
  <Badge tone={isConnected ? 'success' : 'error'}>
    {isConnected ? 'Live Stream Active' : 'Connecting to Stream...'}
  </Badge>
</PageHeader>

<div class="space-y-6">
  <!-- Key Performance Indicators -->
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Calls Waiting</div>
      <div class="text-4xl font-extrabold {metrics.queue_length > 0 ? 'text-amber-600' : 'text-slate-800'} mono">
        {metrics.queue_length}
      </div>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Active Agents</div>
      <div class="text-4xl font-extrabold text-slate-800 mono">
        {metrics.agents_online}
      </div>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Avg Wait Time</div>
      <div class="text-4xl font-extrabold text-sky-700 mono">
        {formatTime(metrics.avg_wait_time)}
      </div>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Global SLA</div>
      <div class="text-4xl font-extrabold {metrics.sla_percentage < 90 ? 'text-amber-600' : 'text-emerald-600'} mono">
        {metrics.sla_percentage}%
      </div>
    </div>
  </div>

  <!-- Detailed Queues Table -->
  <Card padded={false} title="Queue Performance Breakdown">
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Queue Name</th>
            <th class="text-center">Callers Waiting</th>
            <th class="text-center">Agents Ready</th>
            <th class="text-right">SLA %</th>
          </tr>
        </thead>
        <tbody>
          {#each queues as q}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{q.name}</td>
              <td class="text-center mono font-bold text-sm {q.waiting > 0 ? 'text-amber-600' : 'text-slate-500'}">
                {q.waiting}
              </td>
              <td class="text-center mono text-sky-700 font-bold text-sm">{q.agents}</td>
              <td class="text-right">
                <span class="mono font-bold {q.sla < 90 ? 'text-amber-600' : 'text-emerald-600'}">
                  {q.sla}%
                </span>
              </td>
            </tr>
          {:else}
            <tr>
              <td colspan="4" class="p-8 text-center text-slate-500 text-xs">
                {isConnected ? 'No queues reporting metrics.' : 'Connecting to live statistics stream…'}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </Card>
</div>
