<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, formatDuration } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface CDRItem {
    ID: string;
    UUID?: string;
    CallerIDNumber: string;
    CallerIDName: string;
    Destination: string;
    Context: string;
    Duration: number;
    Billsec: number;
    HangupCause: string;
    Transcription: string;
    CreatedAt: string;
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
    cdr: CDRItem;
    hops: CallHop[];
    codec: string;
    mos: number;
    jitter_ms: number;
  }

  let cdrs: CDRItem[] = [];
  let loading = true;
  let loadError = '';
  let search = '';
  let filterContext = 'ALL';

  // Trace Modal
  let traceModalOpen = false;
  let traceLoading = false;
  let traceData: TraceData | null = null;

  $: filtered = cdrs.filter(c => {
    const matchesSearch = !search.trim() || 
      `${c.CallerIDNumber} ${c.CallerIDName} ${c.Destination} ${c.HangupCause}`.toLowerCase().includes(search.trim().toLowerCase());
    const matchesContext = filterContext === 'ALL' || c.Context === filterContext;
    return matchesSearch && matchesContext;
  });

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ cdrs: CDRItem[] }>('/cdr');
      cdrs = data.cdrs || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  async function openTrace(c: CDRItem) {
    traceModalOpen = true;
    traceLoading = true;
    try {
      const targetId = c.UUID || c.ID;
      traceData = await api.get<TraceData>(`/cdr/${targetId}/trace`);
    } catch (e: any) {
      traceData = null;
    } finally {
      traceLoading = false;
    }
  }

  function exportCSV() {
    if (cdrs.length === 0) return;
    const headers = ['Time', 'Caller ID', 'Destination', 'Context', 'Duration (s)', 'Billable (s)', 'Disposition', 'Transcription'];
    const rows = cdrs.map(c => [
      c.CreatedAt,
      `"${c.CallerIDNumber || ''}"`,
      `"${c.Destination || ''}"`,
      c.Context || '',
      c.Duration || 0,
      c.Billsec || 0,
      c.HangupCause || '',
      `"${(c.Transcription || '').replace(/"/g, '""')}"`
    ]);
    const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `msip_cdrs_${new Date().toISOString().slice(0, 10)}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  }

  onMount(load);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Call Detail Records (CDRs) & Tracking"
    subtitle="Individual call session logs recorded by FreeSWITCH and Kamailio with click-to-trace lifecycle inspection."
  >
    <div class="flex items-center space-x-2">
      <Button variant="secondary" on:click={exportCSV} disabled={cdrs.length === 0}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        Export CSV
      </Button>
      <Button variant="secondary" on:click={load}>Refresh</Button>
    </div>
  </PageHeader>

  <Card padding="none">
    <div class="p-4 border-b border-slate-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
      <div class="text-xs text-slate-500">
        Total Logged Sessions: <span class="font-bold text-slate-700">{cdrs.length}</span>
      </div>
      <div class="flex items-center space-x-2">
        <select
          class="text-xs bg-white border border-slate-200 rounded-md px-2.5 py-1.5 focus:outline-none focus:border-sky-500 shadow-xs"
          bind:value={filterContext}
        >
          <option value="ALL">All Contexts</option>
          <option value="internal">Internal / Station</option>
          <option value="public">Inbound PSTN</option>
          <option value="outbound">Outbound Gateway</option>
        </select>
        <div class="relative w-56">
          <input
            type="text"
            placeholder="Search caller, destination..."
            bind:value={search}
            class="w-full pl-8 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 shadow-xs"
          />
          <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    {#if loading}
      <div class="p-12 text-center text-slate-400 text-sm">Loading call logs...</div>
    {:else if loadError}
      <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm m-4">{loadError}</div>
    {:else if filtered.length === 0}
      <EmptyState
        title={search ? 'No Calls Match Your Filter' : 'No Call Records Logged Yet'}
        description={search ? 'Try clearing your search query or switching context filter.' : 'Completed SIP calls processed through the system will appear here.'}
      />
    {:else}
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
              <th class="py-3 px-4">Date & Time</th>
              <th class="py-3 px-4">Caller ID</th>
              <th class="py-3 px-4">Destination</th>
              <th class="py-3 px-4">Context</th>
              <th class="py-3 px-4">Talk Time</th>
              <th class="py-3 px-4">Disposition</th>
              <th class="py-3 px-4 text-right">Drill-Down Trace</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            {#each filtered as c}
              <tr class="hover:bg-slate-50/70 transition-colors cursor-pointer" on:click={() => openTrace(c)}>
                <td class="py-3 px-4 text-slate-500 font-mono text-[11px] whitespace-nowrap">
                  {formatDateTime(c.CreatedAt)}
                </td>
                <td class="py-3 px-4 font-mono font-bold text-slate-900">
                  {c.CallerIDNumber || 'Anonymous'}
                  {#if c.CallerIDName && c.CallerIDName !== c.CallerIDNumber}
                    <span class="text-slate-400 font-normal ml-1">({c.CallerIDName})</span>
                  {/if}
                </td>
                <td class="py-3 px-4 font-mono font-semibold text-sky-700">
                  {c.Destination}
                </td>
                <td class="py-3 px-4">
                  <span class="px-2 py-0.5 bg-slate-100 rounded text-slate-700 font-medium text-[11px] capitalize">
                    {c.Context || 'default'}
                  </span>
                </td>
                <td class="py-3 px-4 font-mono text-slate-800 font-medium">
                  {formatDuration(c.Billsec || c.Duration)}
                </td>
                <td class="py-3 px-4">
                  {#if c.HangupCause === 'NORMAL_CLEARING' && (c.Billsec || 0) > 0}
                    <Badge variant="success">Answered</Badge>
                  {:else if c.HangupCause === 'USER_BUSY'}
                    <Badge variant="warning">Busy</Badge>
                  {:else}
                    <Badge variant="neutral">Missed ({c.HangupCause || 'No Answer'})</Badge>
                  {/if}
                  {#if c.Transcription}
                    <span class="ml-1.5 px-1.5 py-0.5 bg-purple-50 text-purple-700 border border-purple-200 rounded text-[10px] font-bold">
                      AI Audio
                    </span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-right">
                  <button
                    type="button"
                    class="text-sky-600 hover:text-sky-800 font-semibold inline-flex items-center space-x-1"
                    on:click|stopPropagation={() => openTrace(c)}
                  >
                    <span>Trace Call</span>
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </Card>
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
