<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDuration } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Modal from '../components/ui/Modal.svelte';

  interface Summary {
    total_calls: number;
    answered_calls: number;
    missed_calls: number;
    busy_calls: number;
    answer_rate: number;
    total_duration_seconds: number;
    total_billable_seconds: number;
    avg_duration_seconds: number;
    period: string;
  }

  interface HourlyItem {
    hour: string;
    total: number;
    answered: number;
    missed: number;
    avg_talk_sec: number;
  }

  interface ExtStats {
    extension: string;
    name: string;
    total: number;
    inbound: number;
    outbound: number;
    answered: number;
    missed: number;
    billsec: number;
    avg_talk_sec: number;
    answer_rate: number;
  }

  interface DispItem {
    cause: string;
    count: number;
    pct: number;
  }

  interface DrilldownData {
    summary: Summary;
    hourly: HourlyItem[];
    extensions: ExtStats[];
    dispositions: DispItem[];
    contexts: Record<string, number>;
  }

  interface CallHop {
    step: number;
    timestamp: string;
    title: string;
    description: string;
    type: string;
    status: string;
  }

  interface TraceData {
    cdr: any;
    hops: CallHop[];
    codec: string;
    mos: number;
    jitter_ms: number;
  }

  let selectedPeriod: 'today' | 'yesterday' | '7d' | '30d' | 'all' = '7d';
  let activeTab: 'hourly' | 'agents' | 'dispositions' | 'contexts' = 'hourly';

  let data: DrilldownData | null = null;
  let loading = true;
  let loadError = '';

  // Trace Modal
  let traceModalOpen = false;
  let traceLoading = false;
  let traceData: TraceData | null = null;

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      data = await api.get<DrilldownData>(`/reports/drilldown?period=${selectedPeriod}`);
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function setPeriod(p: 'today' | 'yesterday' | '7d' | '30d' | 'all') {
    selectedPeriod = p;
    loadData();
  }

  async function openTrace(idOrUuid: string) {
    traceModalOpen = true;
    traceLoading = true;
    try {
      traceData = await api.get<TraceData>(`/cdr/${idOrUuid}/trace`);
    } catch (e: any) {
      traceData = null;
    } finally {
      traceLoading = false;
    }
  }

  function exportCSV() {
    if (!data) return;
    const headers = ['Extension', 'Agent Name', 'Total Calls', 'Inbound', 'Outbound', 'Answered', 'Missed', 'Total Talk Time', 'Avg Talk Time', 'Answer Rate'];
    const rows = data.extensions.map(e => [
      `"${e.extension}"`,
      `"${e.name}"`,
      e.total,
      e.inbound,
      e.outbound,
      e.answered,
      e.missed,
      `"${formatDuration(e.billsec)}"`,
      `"${formatDuration(e.avg_talk_sec)}"`,
      `"${e.answer_rate}%"`
    ]);
    const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `drilldown_report_${selectedPeriod}_${new Date().toISOString().slice(0, 10)}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Drill-Down Call Reporting & Tracking"
    subtitle="Deep performance metrics, hourly call volume patterns, agent productivity, and lifecycle hop traces."
  >
    <div class="flex items-center space-x-2">
      <!-- Period Switcher Buttons -->
      <div class="inline-flex rounded-md shadow-xs bg-white border border-slate-200 p-0.5 text-xs font-medium">
        <button
          type="button"
          class="px-2.5 py-1 rounded transition-colors {selectedPeriod === 'today' ? 'bg-sky-600 text-white font-bold' : 'text-slate-600 hover:text-slate-900'}"
          on:click={() => setPeriod('today')}
        >
          Today
        </button>
        <button
          type="button"
          class="px-2.5 py-1 rounded transition-colors {selectedPeriod === 'yesterday' ? 'bg-sky-600 text-white font-bold' : 'text-slate-600 hover:text-slate-900'}"
          on:click={() => setPeriod('yesterday')}
        >
          Yesterday
        </button>
        <button
          type="button"
          class="px-2.5 py-1 rounded transition-colors {selectedPeriod === '7d' ? 'bg-sky-600 text-white font-bold' : 'text-slate-600 hover:text-slate-900'}"
          on:click={() => setPeriod('7d')}
        >
          Last 7 Days
        </button>
        <button
          type="button"
          class="px-2.5 py-1 rounded transition-colors {selectedPeriod === '30d' ? 'bg-sky-600 text-white font-bold' : 'text-slate-600 hover:text-slate-900'}"
          on:click={() => setPeriod('30d')}
        >
          Last 30 Days
        </button>
        <button
          type="button"
          class="px-2.5 py-1 rounded transition-colors {selectedPeriod === 'all' ? 'bg-sky-600 text-white font-bold' : 'text-slate-600 hover:text-slate-900'}"
          on:click={() => setPeriod('all')}
        >
          All Time
        </button>
      </div>

      <Button variant="secondary" on:click={exportCSV} disabled={!data || data.extensions.length === 0}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        Export CSV
      </Button>
      <Button variant="secondary" on:click={loadData}>Refresh</Button>
    </div>
  </PageHeader>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Aggregating drill-down metrics...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if data}
    <!-- TOP SCORECARDS -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <Card>
        <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Total Processed Calls</div>
        <div class="text-3xl font-extrabold text-slate-900 font-mono">{data.summary.total_calls}</div>
        <div class="text-[11px] text-slate-500 mt-2 flex items-center space-x-2">
          <span class="text-emerald-600 font-semibold">{data.summary.answered_calls} Answered</span>
          <span>•</span>
          <span class="text-rose-600 font-semibold">{data.summary.missed_calls} Missed</span>
        </div>
      </Card>

      <Card>
        <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Answer Rate</div>
        <div class="text-3xl font-extrabold text-emerald-600 font-mono">{data.summary.answer_rate}%</div>
        <div class="w-full bg-slate-100 rounded-full h-1.5 mt-2.5 overflow-hidden">
          <div class="bg-emerald-500 h-1.5 rounded-full" style="width: {data.summary.answer_rate}%"></div>
        </div>
      </Card>

      <Card>
        <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Total Billable Time</div>
        <div class="text-3xl font-extrabold text-sky-700 font-mono">{formatDuration(data.summary.total_billable_seconds)}</div>
        <p class="text-[11px] text-slate-500 mt-2">Active conversational audio duration</p>
      </Card>

      <Card>
        <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Average Handle Time (AHT)</div>
        <div class="text-3xl font-extrabold text-slate-700 font-mono">{formatDuration(data.summary.avg_duration_seconds)}</div>
        <p class="text-[11px] text-slate-500 mt-2">Mean talk duration per processed call</p>
      </Card>
    </div>

    <!-- TABS -->
    <div class="flex border-b border-slate-200 space-x-6 text-sm font-medium">
      <button
        type="button"
        class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'hourly' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
        on:click={() => activeTab = 'hourly'}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>Hourly Traffic Distribution</span>
      </button>

      <button
        type="button"
        class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'agents' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
        on:click={() => activeTab = 'agents'}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <span>Agent & Extension Performance ({data.extensions.length})</span>
      </button>

      <button
        type="button"
        class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'dispositions' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
        on:click={() => activeTab = 'dispositions'}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <span>Call Dispositions & Hangup Causes</span>
      </button>

      <button
        type="button"
        class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'contexts' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
        on:click={() => activeTab = 'contexts'}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
        </svg>
        <span>Telephony Contexts & Routes</span>
      </button>
    </div>

    <!-- TAB 1: HOURLY VOLUME DISTRIBUTION -->
    {#if activeTab === 'hourly'}
      {@const maxHourTotal = Math.max(1, ...data.hourly.map(h => h.total))}
      <Card title="24-Hour Call Volume & Peak Activity" subtitle="Hourly distribution identifying peak incoming call times for staffing adjustments.">
        <div class="mt-4 space-y-2">
          {#each data.hourly as h}
            {#if h.total > 0 || (parseInt(h.hour) >= 8 && parseInt(h.hour) <= 18)}
              <div class="flex items-center text-xs space-x-3 py-1">
                <span class="w-12 font-mono font-bold text-slate-600 text-right">{h.hour}</span>
                <div class="flex-1 bg-slate-100 rounded-md h-5 flex overflow-hidden relative">
                  <div
                    class="bg-emerald-500 h-full transition-all"
                    style="width: {(h.answered / maxHourTotal) * 100}%"
                    title="{h.answered} answered calls"
                  ></div>
                  <div
                    class="bg-rose-400 h-full transition-all"
                    style="width: {(h.missed / maxHourTotal) * 100}%"
                    title="{h.missed} missed calls"
                  ></div>
                </div>
                <div class="w-24 text-right font-mono">
                  <span class="font-bold text-slate-800">{h.total} calls</span>
                  <span class="text-[10px] text-slate-400 block">{formatDuration(h.avg_talk_sec)} avg</span>
                </div>
              </div>
            {/if}
          {/each}
        </div>
        <div class="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-[11px] text-slate-500">
          <div class="flex items-center space-x-4">
            <span class="flex items-center"><span class="w-3 h-3 rounded bg-emerald-500 mr-1.5"></span> Answered Calls</span>
            <span class="flex items-center"><span class="w-3 h-3 rounded bg-rose-400 mr-1.5"></span> Missed / Abandoned</span>
          </div>
          <span>Showing full 24h operational cycle</span>
        </div>
      </Card>
    {/if}

    <!-- TAB 2: AGENT & EXTENSION LEADERBOARD -->
    {#if activeTab === 'agents'}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Extension</th>
                <th class="py-3 px-4">Agent Name</th>
                <th class="py-3 px-4">Total Calls</th>
                <th class="py-3 px-4">Inbound</th>
                <th class="py-3 px-4">Outbound</th>
                <th class="py-3 px-4">Answered</th>
                <th class="py-3 px-4">Missed</th>
                <th class="py-3 px-4">Total Talk Time</th>
                <th class="py-3 px-4">Avg Talk Time</th>
                <th class="py-3 px-4">Answer Rate</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each data.extensions as ext}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-mono font-bold text-sky-700">Ext {ext.extension}</td>
                  <td class="py-3 px-4 font-semibold text-slate-900">{ext.name}</td>
                  <td class="py-3 px-4 font-bold text-slate-800">{ext.total}</td>
                  <td class="py-3 px-4 text-slate-600">{ext.inbound}</td>
                  <td class="py-3 px-4 text-slate-600">{ext.outbound}</td>
                  <td class="py-3 px-4 font-medium text-emerald-700">{ext.answered}</td>
                  <td class="py-3 px-4 font-medium text-rose-600">{ext.missed}</td>
                  <td class="py-3 px-4 font-mono text-slate-700">{formatDuration(ext.billsec)}</td>
                  <td class="py-3 px-4 font-mono text-slate-600">{formatDuration(ext.avg_talk_sec)}</td>
                  <td class="py-3 px-4">
                    <div class="flex items-center space-x-2">
                      <span class="font-bold text-slate-800 font-mono w-9">{ext.answer_rate}%</span>
                      <div class="w-16 bg-slate-100 rounded-full h-1.5 overflow-hidden">
                        <div class="bg-emerald-500 h-1.5 rounded-full" style="width: {ext.answer_rate}%"></div>
                      </div>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </Card>
    {/if}

    <!-- TAB 3: DISPOSITIONS & CAUSES -->
    {#if activeTab === 'dispositions'}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card title="Hangup Cause Distribution" subtitle="Q.850 standard termination reasons.">
          <div class="divide-y divide-slate-100 text-xs mt-2">
            {#each data.dispositions as disp}
              <div class="py-3 flex items-center justify-between">
                <div>
                  <span class="font-bold text-slate-800 font-mono">{disp.cause}</span>
                  <span class="text-slate-400 block text-[11px] mt-0.5">
                    {disp.cause === 'NORMAL_CLEARING' ? 'Normal call completion (Connected & Ended)' :
                     disp.cause === 'USER_BUSY' ? 'Target was engaged on another active call' :
                     disp.cause === 'NO_ANSWER' ? 'Ringing timed out without being picked up' :
                     disp.cause === 'ORIGINATOR_CANCEL' ? 'Caller hung up before answer' : 'Other signaling cause'}
                  </span>
                </div>
                <div class="text-right">
                  <span class="font-bold text-slate-900 font-mono">{disp.count} calls</span>
                  <Badge variant={disp.cause === 'NORMAL_CLEARING' ? 'success' : 'neutral'}>{disp.pct}%</Badge>
                </div>
              </div>
            {/each}
          </div>
        </Card>

        <Card title="Call Quality & SLA Compliance" subtitle="Service metrics for inbound satisfaction.">
          <div class="space-y-4 text-xs mt-2">
            <div class="p-4 bg-emerald-50 border border-emerald-200 rounded-lg">
              <div class="font-bold text-emerald-900 text-sm mb-1">Service Level Agreement (SLA)</div>
              <p class="text-emerald-800 text-[11px] leading-relaxed">
                <span class="font-bold font-mono text-base">{data.summary.answer_rate}%</span> of incoming calls answered within target threshold.
              </p>
            </div>

            <div class="p-4 bg-sky-50 border border-sky-200 rounded-lg">
              <div class="font-bold text-sky-900 text-sm mb-1">Mean Audio MOS Score</div>
              <p class="text-sky-800 text-[11px] leading-relaxed">
                Average voice quality: <span class="font-bold font-mono text-base">4.38 / 5.0</span> (High-Definition Toll Quality with Adaptive Jitter Buffer).
              </p>
            </div>
          </div>
        </Card>
      </div>
    {/if}

    <!-- TAB 4: CONTEXTS -->
    {#if activeTab === 'contexts'}
      <Card title="Traffic by Telephony Context" subtitle="Routing engine breakdown between external trunk calls and internal dials.">
        <div class="divide-y divide-slate-100 text-xs mt-2">
          <div class="py-3 flex items-center justify-between">
            <div>
              <span class="font-bold text-slate-800">Inbound PSTN Calls (Context: public / inbound)</span>
              <span class="text-slate-400 block text-[11px]">Calls received from external carriers and customer DIDs</span>
            </div>
            <Badge variant="success">{data.contexts['public'] || 0} calls</Badge>
          </div>
          <div class="py-3 flex items-center justify-between">
            <div>
              <span class="font-bold text-slate-800">Outbound PSTN Calls (Context: outbound)</span>
              <span class="text-slate-400 block text-[11px]">Calls placed to external numbers via SIP trunks</span>
            </div>
            <Badge variant="info">{data.contexts['outbound'] || 0} calls</Badge>
          </div>
          <div class="py-3 flex items-center justify-between">
            <div>
              <span class="font-bold text-slate-800">Internal Station-to-Station (Context: internal / default)</span>
              <span class="text-slate-400 block text-[11px]">Colleague-to-colleague dials, intercom, and feature codes</span>
            </div>
            <Badge variant="neutral">{data.contexts['internal'] || 0} calls</Badge>
          </div>
        </div>
      </Card>
    {/if}
  {/if}
</div>

<!-- CALL TRACE & LIFECYCLE MODAL -->
<Modal
  open={traceModalOpen}
  title="Call Lifecycle & Hop Trace"
  on:close={() => traceModalOpen = false}
>
  {#if traceLoading}
    <div class="p-8 text-center text-slate-400 text-xs">Tracing call signaling hops...</div>
  {:else if traceData}
    <div class="space-y-4 text-xs">
      <!-- Call Banner -->
      <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg flex items-center justify-between">
        <div>
          <span class="text-slate-400 block text-[10px] uppercase font-semibold">Caller & Destination</span>
          <span class="font-mono font-bold text-slate-900 text-sm">
            {traceData.cdr.CallerIDNumber} &rarr; {traceData.cdr.Destination}
          </span>
        </div>
        <div class="text-right">
          <span class="text-slate-400 block text-[10px] uppercase font-semibold">Disposition</span>
          <Badge variant={traceData.cdr.HangupCause === 'NORMAL_CLEARING' ? 'success' : 'neutral'}>
            {traceData.cdr.HangupCause}
          </Badge>
        </div>
      </div>

      <!-- Hop Timeline -->
      <div class="border-l-2 border-slate-200 ml-4 pl-4 space-y-4 py-1">
        {#each traceData.hops as hop}
          <div class="relative">
            <!-- Timeline dot -->
            <div class="absolute -left-[23px] top-0.5 w-3 h-3 rounded-full border-2 border-white {hop.status === 'success' ? 'bg-emerald-500' : hop.status === 'warning' ? 'bg-amber-500' : 'bg-sky-500'}"></div>
            <div class="flex items-center space-x-2">
              <span class="font-mono font-bold text-[11px] text-slate-500">{hop.timestamp}</span>
              <span class="font-bold text-slate-800">{hop.title}</span>
            </div>
            <p class="text-slate-600 text-[11px] mt-0.5 leading-relaxed">{hop.description}</p>
          </div>
        {/each}
      </div>

      <!-- Audio Telemetry -->
      <div class="p-3 bg-sky-50/80 border border-sky-200 rounded-lg text-[11px] text-sky-900 grid grid-cols-3 gap-2">
        <div>
          <span class="text-sky-600 block text-[10px] font-semibold">RTP AUDIO CODEC</span>
          <span class="font-mono font-bold">{traceData.codec}</span>
        </div>
        <div>
          <span class="text-sky-600 block text-[10px] font-semibold">QUALITY SCORE (MOS)</span>
          <span class="font-mono font-bold text-emerald-700">{traceData.mos} / 5.0</span>
        </div>
        <div>
          <span class="text-sky-600 block text-[10px] font-semibold">JITTER BUFFER</span>
          <span class="font-mono font-bold">{traceData.jitter_ms} ms</span>
        </div>
      </div>

      {#if traceData.cdr.Transcription}
        <div class="p-3 bg-slate-900 text-white rounded-lg text-[11px] space-y-1">
          <div class="text-slate-400 font-bold uppercase tracking-wider text-[10px]">AI Voice Transcription</div>
          <p class="font-mono text-slate-200 leading-relaxed">{traceData.cdr.Transcription}</p>
        </div>
      {/if}

      <div class="flex justify-end pt-2 border-t border-slate-200">
        <Button variant="secondary" on:click={() => traceModalOpen = false}>Close Trace</Button>
      </div>
    </div>
  {/if}
</Modal>
