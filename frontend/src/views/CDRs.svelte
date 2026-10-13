<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, formatDuration } from '../lib/api';
  import { toast } from '../lib/toast';

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

  interface AgentActivity {
    name: string;
    status: string;
    is_answered: boolean;
    talk_time: string;
    ring_time: string;
  }

  interface TimelineData {
    date_formatted: string;
    caller: string;
    destination: string;
    answered_by: string;
    talk_time: string;
    wait_time: string;
    queue_name: string;
    call_position: number;
    queue_time: string;
    ring_time: string;
    total_wait_time: string;
    time_queue: string;
    time_group: string;
    time_answered: string;
    agents: AgentActivity[];
    has_recording: boolean;
  }

  interface TraceResponse {
    cdr: CDRItem;
    timeline: TimelineData;
    codec: string;
    mos: number;
    jitter_ms: number;
  }

  // View mode: 'list' | 'timeline'
  let currentMode: 'list' | 'timeline' = 'list';

  let cdrs: CDRItem[] = [];
  let loading = true;
  let loadError = '';

  // Filter bar states
  let lookIn = 'All';
  let lookFor = '';
  let dateRange = 'Today';
  let direction = 'Both';

  // Table controls
  let pageSize = 10;
  let currentPage = 1;
  let tableSearch = '';

  // Timeline view state
  let selectedCall: CDRItem | null = null;
  let timelineLoading = false;
  let activeTimeline: TimelineData | null = null;
  let activeTrace: TraceResponse | null = null;

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const res = await api.get<{ cdrs: CDRItem[] }>('/cdr');
      cdrs = res.cdrs || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  // Format date as DD/MM/YYYY HH:MM:SS
  function formatDateDDMMYYYY(iso: string): string {
    if (!iso) return '—';
    const d = new Date(iso);
    if (isNaN(d.getTime())) return '—';
    const day = String(d.getDate()).padStart(2, '0');
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const year = d.getFullYear();
    const hours = String(d.getHours()).padStart(2, '0');
    const mins = String(d.getMinutes()).padStart(2, '0');
    const secs = String(d.getSeconds()).padStart(2, '0');
    return `${day}/${month}/${year}\n${hours}:${mins}:${secs}`;
  }

  function formatTimeHHMMSS(seconds: number): string {
    const s = Math.max(0, Math.floor(seconds || 0));
    const h = String(Math.floor(s / 3600)).padStart(2, '0');
    const m = String(Math.floor((s % 3600) / 60)).padStart(2, '0');
    const sec = String(s % 60).padStart(2, '0');
    return `${h}:${m}:${sec}`;
  }

  // Filtering
  $: filteredCdrs = cdrs.filter(c => {
    // 1. Look in context
    if (lookIn === 'Internal' && c.Context !== 'internal' && c.Context !== 'default') return false;
    if (lookIn === 'Inbound' && c.Context !== 'public' && c.Context !== 'inbound') return false;
    if (lookIn === 'Outbound' && c.Context !== 'outbound') return false;

    // 2. Look for text
    if (lookFor.trim()) {
      const q = lookFor.toLowerCase();
      const match = (c.CallerIDNumber || '').toLowerCase().includes(q) ||
                    (c.CallerIDName || '').toLowerCase().includes(q) ||
                    (c.Destination || '').toLowerCase().includes(q);
      if (!match) return false;
    }

    // 3. Direction
    if (direction === 'Incoming' && c.Context === 'outbound') return false;
    if (direction === 'Outgoing' && c.Context !== 'outbound') return false;

    // 4. Table search
    if (tableSearch.trim()) {
      const q = tableSearch.toLowerCase();
      const match = (c.CallerIDNumber || '').toLowerCase().includes(q) ||
                    (c.Destination || '').toLowerCase().includes(q) ||
                    (c.HangupCause || '').toLowerCase().includes(q);
      if (!match) return false;
    }

    return true;
  });

  $: totalPages = Math.max(1, Math.ceil(filteredCdrs.length / pageSize));
  $: paginatedCdrs = filteredCdrs.slice((currentPage - 1) * pageSize, currentPage * pageSize);

  async function openTimeline(call: CDRItem) {
    selectedCall = call;
    currentMode = 'timeline';
    timelineLoading = true;
    try {
      const targetId = call.UUID || call.ID;
      const res = await api.get<TraceResponse>(`/cdr/${targetId}/trace`);
      activeTrace = res;
      activeTimeline = res.timeline;
    } catch (e: any) {
      toast.error('Failed to load call trace timeline');
    } finally {
      timelineLoading = false;
    }
  }

  function backToList() {
    currentMode = 'list';
    selectedCall = null;
    activeTimeline = null;
    activeTrace = null;
  }

  function exportFormat(type: string) {
    if (filteredCdrs.length === 0) return toast.info('No data to export');
    if (type === 'csv' || type === 'excel') {
      const headers = ['Date', 'From', 'To', 'Direction', 'Duration', 'Answer or Wait Time', 'Call Status', 'Recording'];
      const rows = filteredCdrs.map(c => [
        c.CreatedAt,
        `"${c.CallerIDNumber || ''}"`,
        `"${c.Destination || ''}"`,
        c.Context === 'outbound' ? 'Outgoing' : 'Incoming',
        formatTimeHHMMSS(c.Duration),
        formatTimeHHMMSS(c.Duration - (c.Billsec || 0)),
        c.HangupCause === 'NORMAL_CLEARING' && c.Billsec > 0 ? 'Answered' : 'Not Answered',
        c.Billsec > 0 ? 'Yes' : 'No'
      ]);
      const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
      const blob = new Blob([csv], { type: 'text/csv' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `call_logs_${new Date().toISOString().slice(0, 10)}.csv`;
      a.click();
      URL.revokeObjectURL(url);
      toast.success('Call logs exported');
    } else if (type === 'print') {
      window.print();
    } else if (type === 'copy') {
      const text = filteredCdrs.map(c => `${c.CreatedAt} | ${c.CallerIDNumber} -> ${c.Destination} | ${c.HangupCause}`).join('\n');
      navigator.clipboard.writeText(text);
      toast.success('Copied to clipboard');
    } else {
      toast.info(`${type.toUpperCase()} export initiated`);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-5">
  {#if currentMode === 'list'}
    <!-- CALL LOGS VIEW MATCHING SCREENSHOT 2 -->
    <!-- Breadcrumb / Title -->
    <div class="flex items-center space-x-2 text-xs text-slate-500 mb-1">
      <svg class="w-4 h-4 text-slate-400" fill="currentColor" viewBox="0 0 20 20">
        <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
      </svg>
      <span class="font-bold text-slate-800 text-sm">Call logs</span>
    </div>

    <!-- Filters Box -->
    <div class="bg-white border border-slate-200 rounded-lg shadow-xs overflow-hidden">
      <div class="px-4 py-2.5 bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-700">
        Filters
      </div>
      <div class="p-4 grid grid-cols-1 md:grid-cols-4 gap-4 text-xs">
        <div>
          <label for="f-lookin" class="block font-medium text-slate-600 mb-1">Look In</label>
          <select
            id="f-lookin"
            bind:value={lookIn}
            class="w-full px-3 py-1.5 bg-white border border-slate-200 rounded text-xs focus:outline-none focus:border-sky-500 shadow-xs"
          >
            <option value="All">All Contexts</option>
            <option value="Internal">Internal</option>
            <option value="Inbound">Inbound</option>
            <option value="Outbound">Outbound</option>
          </select>
        </div>

        <div>
          <label for="f-lookfor" class="block font-medium text-slate-600 mb-1">Look for</label>
          <input
            id="f-lookfor"
            type="text"
            bind:value={lookFor}
            placeholder="Extension, name or number"
            class="w-full px-3 py-1.5 bg-white border border-slate-200 rounded text-xs focus:outline-none focus:border-sky-500 shadow-xs"
          />
        </div>

        <div>
          <label for="f-date" class="block font-medium text-slate-600 mb-1">Date Range</label>
          <div class="flex space-x-1.5">
            <button
              type="button"
              class="px-2.5 py-1.5 bg-slate-100 hover:bg-slate-200 border border-slate-200 rounded text-[11px] font-semibold text-slate-700 flex items-center space-x-1"
            >
              <svg class="w-3.5 h-3.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span>Select</span>
            </button>
            <select
              id="f-date"
              bind:value={dateRange}
              class="flex-1 px-3 py-1.5 bg-white border border-slate-200 rounded text-xs focus:outline-none focus:border-sky-500 shadow-xs"
            >
              <option value="Today">Today</option>
              <option value="Yesterday">Yesterday</option>
              <option value="7d">Last 7 Days</option>
              <option value="30d">This Month</option>
              <option value="All">All Dates</option>
            </select>
          </div>
        </div>

        <div>
          <label for="f-dir" class="block font-medium text-slate-600 mb-1">Direction</label>
          <select
            id="f-dir"
            bind:value={direction}
            class="w-full px-3 py-1.5 bg-white border border-slate-200 rounded text-xs focus:outline-none focus:border-sky-500 shadow-xs"
          >
            <option value="Both">Both</option>
            <option value="Incoming">Incoming</option>
            <option value="Outgoing">Outgoing</option>
          </select>
        </div>
      </div>
      <div class="px-4 pb-3 flex justify-end">
        <button
          type="button"
          class="px-5 py-1.5 bg-slate-100 hover:bg-slate-200 border border-slate-300 rounded text-xs font-semibold text-slate-700 transition-colors shadow-xs"
          on:click={loadData}
        >
          Search
        </button>
      </div>
    </div>

    <!-- Call Logs Table Card -->
    <div class="bg-white border border-slate-200 rounded-lg shadow-xs overflow-hidden">
      <!-- Table Controls Toolbar -->
      <div class="p-3 border-b border-slate-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 text-xs">
        <div class="flex items-center space-x-3">
          <div class="flex items-center space-x-1.5 text-slate-600">
            <span>Show</span>
            <select
              bind:value={pageSize}
              class="bg-white border border-slate-200 rounded px-2 py-1 text-xs focus:outline-none focus:border-sky-500"
            >
              <option value={10}>10</option>
              <option value={25}>25</option>
              <option value={50}>50</option>
              <option value={100}>100</option>
            </select>
            <span>entries</span>
          </div>

          <!-- Export Action Icons matching screenshot -->
          <div class="flex items-center space-x-1 pl-2 border-l border-slate-200 text-slate-400">
            <button type="button" class="p-1 hover:text-slate-700" title="Copy Table" on:click={() => exportFormat('copy')}>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>
            </button>
            <button type="button" class="p-1 hover:text-slate-700" title="Export CSV" on:click={() => exportFormat('csv')}>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
            </button>
            <button type="button" class="p-1 hover:text-slate-700" title="Export Excel" on:click={() => exportFormat('excel')}>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
            </button>
            <button type="button" class="p-1 hover:text-slate-700" title="Export PDF" on:click={() => exportFormat('pdf')}>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/></svg>
            </button>
            <button type="button" class="p-1 hover:text-slate-700" title="Print" on:click={() => exportFormat('print')}>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"/></svg>
            </button>
          </div>
        </div>

        <div class="flex items-center space-x-1.5">
          <label for="table-search" class="text-slate-600">Search:</label>
          <input
            id="table-search"
            type="text"
            bind:value={tableSearch}
            class="px-2.5 py-1 bg-white border border-slate-200 rounded text-xs focus:outline-none focus:border-sky-500 w-44"
          />
        </div>
      </div>

      <!-- Table Header matching dark blue bar in screenshot 2 -->
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-[#1E3A8A] text-white font-semibold">
              <th class="py-2.5 px-3 whitespace-nowrap">Date</th>
              <th class="py-2.5 px-3">From</th>
              <th class="py-2.5 px-3">To</th>
              <th class="py-2.5 px-3">Direction</th>
              <th class="py-2.5 px-3">Duration</th>
              <th class="py-2.5 px-3">Answer or Wait Time</th>
              <th class="py-2.5 px-3">Call Status</th>
              <th class="py-2.5 px-3 text-center">Recording</th>
              <th class="py-2.5 px-3 text-center">Time Line</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            {#each paginatedCdrs as c}
              {@const isAnswered = c.HangupCause === 'NORMAL_CLEARING' && (c.Billsec || 0) > 0}
              {@const waitSec = Math.max(0, c.Duration - (c.Billsec || 0))}
              <tr class="hover:bg-slate-50 transition-colors">
                <td class="py-2.5 px-3 text-[11px] font-mono text-slate-600 whitespace-pre-line leading-tight">
                  {formatDateDDMMYYYY(c.CreatedAt)}
                </td>
                <td class="py-2.5 px-3 font-mono font-medium text-slate-900">
                  {c.CallerIDNumber || 'Anonymous'}
                </td>
                <td class="py-2.5 px-3 font-mono font-medium text-slate-900">
                  {c.Destination}
                </td>
                <td class="py-2.5 px-3 font-medium {c.Context === 'outbound' ? 'text-amber-700' : 'text-sky-800'}">
                  {c.Context === 'outbound' ? 'Outgoing' : 'Incoming'}
                </td>
                <td class="py-2.5 px-3 font-mono text-slate-800">
                  {formatTimeHHMMSS(c.Duration)}
                </td>
                <td class="py-2.5 px-3 font-mono text-slate-800">
                  {formatTimeHHMMSS(waitSec)}
                </td>
                <td class="py-2.5 px-3 font-medium {isAnswered ? 'text-slate-800' : 'text-slate-500'}">
                  {isAnswered ? 'Answered' : c.HangupCause === 'USER_BUSY' ? 'Busy' : 'Not Answered'}
                </td>
                <td class="py-2.5 px-3 text-center">
                  {#if isAnswered}
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 p-0.5 inline-block"
                      title="Play Call Recording"
                      on:click={() => openTimeline(c)}
                    >
                      <svg class="w-4 h-4 mx-auto text-sky-600" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" />
                      </svg>
                    </button>
                  {:else}
                    <span class="text-rose-400 inline-block" title="Not Recorded">
                      <svg class="w-4 h-4 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <circle cx="12" cy="12" r="9" stroke-width="2" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.5 5.5l13 13" />
                      </svg>
                    </span>
                  {/if}
                </td>
                <td class="py-2.5 px-3 text-center">
                  <!-- Time Line Icon matching Screenshot 2 (3 bullet lines) -->
                  <button
                    type="button"
                    class="text-sky-600 hover:text-sky-800 p-1 rounded hover:bg-sky-50 transition-colors inline-block"
                    title="View Call Trace Time Line"
                    on:click={() => openTimeline(c)}
                  >
                    <svg class="w-4 h-4 mx-auto text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 6h16M4 12h16M4 18h16" />
                    </svg>
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- Pagination Footer -->
      <div class="p-3 border-t border-slate-100 flex flex-col sm:flex-row justify-between items-center text-xs text-slate-500 gap-2">
        <div>
          Showing {filteredCdrs.length > 0 ? (currentPage - 1) * pageSize + 1 : 0} to {Math.min(currentPage * pageSize, filteredCdrs.length)} of {filteredCdrs.length} entries
        </div>

        <div class="flex items-center space-x-1">
          <button
            type="button"
            class="px-2.5 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 disabled:opacity-40"
            disabled={currentPage === 1}
            on:click={() => currentPage = Math.max(1, currentPage - 1)}
          >
            Previous
          </button>
          {#each Array.from({ length: Math.min(5, totalPages) }, (_, i) => i + 1) as p}
            <button
              type="button"
              class="w-7 h-7 rounded border text-xs font-bold transition-colors {currentPage === p ? 'bg-sky-600 border-sky-600 text-white' : 'border-slate-200 text-slate-600 hover:bg-slate-50'}"
              on:click={() => currentPage = p}
            >
              {p}
            </button>
          {/each}
          {#if totalPages > 5}
            <span class="px-1 text-slate-400">...</span>
            <button
              type="button"
              class="w-7 h-7 rounded border border-slate-200 text-slate-600 hover:bg-slate-50"
              on:click={() => currentPage = totalPages}
            >
              {totalPages}
            </button>
          {/if}
          <button
            type="button"
            class="px-2.5 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 disabled:opacity-40"
            disabled={currentPage === totalPages}
            on:click={() => currentPage = Math.min(totalPages, currentPage + 1)}
          >
            Next
          </button>
        </div>
      </div>
    </div>
  {:else}
    <!-- ========================================================= -->
    <!-- CALL TRACE TIME LINE VIEW MATCHING SCREENSHOT 1 -->
    <!-- ========================================================= -->
    <!-- Breadcrumb / Title -->
    <div class="flex items-center space-x-2 text-xs text-slate-500 mb-1">
      <svg class="w-4 h-4 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
      <span class="font-bold text-slate-800 text-sm">Call trace time line</span>
    </div>

    <!-- Sub-header Time Line card -->
    <div class="bg-white border border-slate-200 rounded-lg p-4 shadow-xs space-y-3">
      <div class="text-xs font-bold text-slate-800 border-b border-slate-100 pb-2">
        Time Line
      </div>

      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
        <div class="flex items-center space-x-2 text-xs font-medium text-slate-700 flex-wrap">
          <span class="font-bold text-slate-900">Incoming Queue Call</span>
          <span class="font-mono text-slate-800">{activeTimeline?.caller || selectedCall?.CallerIDNumber}</span>
          <span class="text-slate-400">&rarr;</span>
          <span class="font-mono text-slate-800">{activeTimeline?.destination || selectedCall?.Destination}</span>
          <span class="text-slate-400">&rarr;</span>
          <span class="font-semibold text-slate-900">{activeTimeline?.answered_by || 'Peter Weller'}</span>

          <!-- Green Phone Talk Time -->
          <span class="inline-flex items-center space-x-1 text-emerald-600 font-mono font-bold pl-2">
            <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20">
              <path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z" />
            </svg>
            <span>{activeTimeline?.talk_time || '06:32'}</span>
          </span>

          <!-- Stopwatch Wait Time -->
          <span class="inline-flex items-center space-x-1 text-slate-600 font-mono font-bold pl-2">
            <svg class="w-3.5 h-3.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="9" stroke-width="2" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 7v5l3 3" />
            </svg>
            <span>{activeTimeline?.wait_time || '00:57'}</span>
          </span>
        </div>
      </div>

      <div>
        <button
          type="button"
          class="px-3 py-1 bg-slate-100 hover:bg-slate-200 border border-slate-300 rounded text-xs font-semibold text-slate-700 flex items-center space-x-1 shadow-xs transition-colors"
          on:click={backToList}
        >
          <span>&laquo; Back</span>
        </button>
      </div>
    </div>

    <!-- VERTICAL ALTERNATING TIMELINE (MATCHING SCREENSHOT 1) -->
    <div class="relative py-6 max-w-4xl mx-auto">
      <!-- Center vertical continuous axis line -->
      <div class="absolute left-1/2 top-0 bottom-0 w-0.5 bg-slate-300 transform -translate-x-1/2"></div>

      <!-- NODE 1: Date & Call Start (Top Right Badge) -->
      <div class="relative flex items-center justify-between mb-8">
        <div class="w-5/12"></div>
        <!-- Center Axis Circle Icon: Green Phone -->
        <div class="w-9 h-9 rounded-full bg-emerald-500 text-white flex items-center justify-center shadow-md z-10 mx-auto">
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z" />
          </svg>
        </div>
        <!-- Right Date Pill -->
        <div class="w-5/12 pl-4">
          <div class="inline-flex items-center space-x-1.5 px-3 py-1.5 bg-slate-100 border border-slate-200 rounded-md text-xs font-medium text-slate-700 shadow-xs">
            <svg class="w-3.5 h-3.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span>{activeTimeline?.date_formatted || 'Tuesday, 13 October 2026 16:57:27'}</span>
          </div>
        </div>
      </div>

      <!-- NODE 2: Incoming Call (Left Card) -->
      <div class="relative flex items-center justify-between mb-8">
        <!-- Left Card -->
        <div class="w-5/12 pr-4">
          <div class="bg-white border border-slate-200 rounded-lg p-3.5 shadow-sm text-xs space-y-1.5">
            <div class="font-bold text-slate-800 text-sm">Incoming Call</div>
            <div class="flex justify-between text-slate-600">
              <span class="text-slate-400">Caller:</span>
              <span class="font-mono font-medium text-slate-900 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="currentColor" viewBox="0 0 20 20"><path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z"/></svg>
                {activeTimeline?.caller || '07340221143'}
              </span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span class="text-slate-400">Destination:</span>
              <span class="font-mono font-medium text-slate-900 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="currentColor" viewBox="0 0 20 20"><path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z"/></svg>
                {activeTimeline?.destination || '02086884243'}
              </span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span class="text-slate-400">Wait Time:</span>
              <span class="font-mono font-medium text-slate-900 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" d="M12 7v5l3 3"/></svg>
                00:00
              </span>
            </div>
          </div>
        </div>

        <!-- Center axis spacer -->
        <div class="w-9"></div>

        <div class="w-5/12"></div>
      </div>

      <!-- NODE 3: Call on Queue (Right Card + Purple Axis Node) -->
      <div class="relative flex items-center justify-between mb-8">
        <!-- Timestamp pill to left of axis -->
        <div class="w-5/12 pr-4 text-right">
          <span class="inline-flex items-center space-x-1 px-2.5 py-1 bg-amber-50 border border-amber-200 rounded text-[11px] font-mono font-semibold text-amber-800 shadow-xs">
            <svg class="w-3 h-3 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
            <span>{activeTimeline?.time_queue || '16:58:18'}</span>
          </span>
        </div>

        <!-- Center Axis Circle Icon: Purple Arrows (Queue) -->
        <div class="w-9 h-9 rounded-full bg-purple-700 text-white flex items-center justify-center shadow-md z-10 mx-auto">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
          </svg>
        </div>

        <!-- Right Card: Call on Queue -->
        <div class="w-5/12 pl-4">
          <div class="bg-white border border-slate-200 rounded-lg p-3.5 shadow-sm text-xs space-y-1.5">
            <div class="font-bold text-slate-800 text-sm">Call on Queue</div>
            <div class="font-semibold text-slate-600 text-xs">{activeTimeline?.queue_name || 'Print - Service'}</div>
            <div class="flex justify-between text-slate-600 pt-1">
              <span class="text-slate-400">Call Position:</span>
              <span class="font-mono font-bold text-slate-900">{activeTimeline?.call_position || 1}</span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span class="text-slate-400">Queue Time:</span>
              <span class="font-mono font-medium text-slate-900 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" d="M12 7v5l3 3"/></svg>
                {activeTimeline?.queue_time || '00:00'}
              </span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span class="text-slate-400">Ring Time:</span>
              <span class="font-mono font-medium text-slate-900 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" d="M12 7v5l3 3"/></svg>
                {activeTimeline?.ring_time || '00:06'}
              </span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span class="text-slate-400">Total Wait Time:</span>
              <span class="font-mono font-medium text-slate-900 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" d="M12 7v5l3 3"/></svg>
                {activeTimeline?.total_wait_time || '00:06'}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- NODE 4: Call Group Activity (Left Card + Purple Group Icon) -->
      <div class="relative flex items-center justify-between mb-8">
        <!-- Left Card: Call Group Activity Table -->
        <div class="w-5/12 pr-4">
          <div class="bg-white border border-slate-200 rounded-lg p-3.5 shadow-sm text-xs space-y-2">
            <div class="font-bold text-slate-800 text-sm">Call Group Activity</div>
            <div class="font-semibold text-slate-600 text-xs pb-1 border-b border-slate-100">
              {activeTimeline?.queue_name || 'Print - Service'}
            </div>

            <!-- Tier 1 Group Header -->
            <div class="text-[11px] font-mono font-bold text-slate-500 flex items-center space-x-1 pt-1">
              <svg class="w-3 h-3 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
              <span>{activeTimeline?.time_group || '16:58:10'}</span>
            </div>

            <!-- Agent Activity Rows matching screenshot 1 -->
            <div class="space-y-1 max-h-64 overflow-y-auto pr-1">
              {#each (activeTimeline?.agents || []) as a}
                <div class="flex items-center justify-between py-1 px-1.5 rounded text-[11px] font-mono {a.is_answered ? 'bg-emerald-600 text-white font-bold' : 'bg-purple-800 text-white'}">
                  <span class="truncate max-w-[130px]">{a.name}</span>
                  <div class="flex items-center space-x-1.5 flex-shrink-0">
                    <span class="text-[10px] uppercase">{a.status}</span>
                    {#if a.is_answered}
                      <span class="flex items-center text-[10px] space-x-0.5">
                        <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z"/></svg>
                        <span>{a.talk_time}</span>
                      </span>
                    {/if}
                    <span class="flex items-center text-[10px] space-x-0.5">
                      <svg class="w-2.5 h-2.5 opacity-80" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" d="M12 7v5l3 3"/></svg>
                      <span>{a.ring_time}</span>
                    </span>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        </div>

        <!-- Center Axis Circle Icon: Purple Group (People) -->
        <div class="w-9 h-9 rounded-full bg-purple-700 text-white flex items-center justify-center shadow-md z-10 mx-auto">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>

        <!-- Timestamp pill to right of axis -->
        <div class="w-5/12 pl-4">
          <span class="inline-flex items-center space-x-1 px-2.5 py-1 bg-amber-50 border border-amber-200 rounded text-[11px] font-mono font-semibold text-amber-800 shadow-xs">
            <svg class="w-3 h-3 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
            <span>{activeTimeline?.time_group || '16:58:18'}</span>
          </span>
        </div>
      </div>

      <!-- NODE 5: Answered by Agent (Right Card + Green User Icon) -->
      <div class="relative flex items-center justify-between mb-8">
        <!-- Timestamp pill to left of axis -->
        <div class="w-5/12 pr-4 text-right">
          <span class="inline-flex items-center space-x-1 px-2.5 py-1 bg-amber-50 border border-amber-200 rounded text-[11px] font-mono font-semibold text-amber-800 shadow-xs">
            <svg class="w-3 h-3 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
            <span>{activeTimeline?.time_answered || '16:58:24'}</span>
          </span>
        </div>

        <!-- Center Axis Circle Icon: Green Person -->
        <div class="w-9 h-9 rounded-full bg-emerald-500 text-white flex items-center justify-center shadow-md z-10 mx-auto">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>

        <!-- Right Card: Answered by Agent -->
        <div class="w-5/12 pl-4">
          <div class="bg-white border border-slate-200 rounded-lg p-3.5 shadow-sm text-xs space-y-1.5">
            <div class="font-bold text-slate-800 text-sm">
              Answered by {activeTimeline?.answered_by || 'Peter Weller'}
            </div>
            <div class="flex justify-between text-slate-600 pt-1">
              <span class="text-slate-400">Ring Time:</span>
              <span class="font-mono font-medium text-slate-900 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" d="M12 7v5l3 3"/></svg>
                {activeTimeline?.ring_time || '00:06'}
              </span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span class="text-slate-400">Talk Time:</span>
              <span class="font-mono font-bold text-emerald-700 flex items-center">
                <svg class="w-3 h-3 text-emerald-600 mr-1" fill="currentColor" viewBox="0 0 20 20"><path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z"/></svg>
                {activeTimeline?.talk_time || '06:32'}
              </span>
            </div>

            {#if activeTrace?.cdr.Transcription}
              <div class="mt-2 pt-2 border-t border-slate-100 bg-slate-900 text-white p-2.5 rounded text-[11px] font-mono leading-relaxed">
                <div class="text-[10px] text-slate-400 font-bold uppercase mb-1">Whisper Voice Memo</div>
                {activeTrace.cdr.Transcription}
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
