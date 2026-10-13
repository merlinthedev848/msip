<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import { onMount } from 'svelte';
  import { apiFetch } from '../lib/api';

  interface CDRItem {
    ID: string;
    CreatedAt: string;
    CallerIDName: string;
    CallerIDNumber: string;
    Destination: string;
    Duration: number;
    Billsec: number;
    HangupCause: string;
    Context?: string;
  }

  let cdrs: CDRItem[] = [];
  let isLoading = true;
  let searchTerm = '';
  let statusFilter: 'all' | 'completed' | 'cancelled' | 'failed' = 'all';

  async function fetchCDRs() {
    try {
      const res = await apiFetch('/cdr', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        cdrs = data.cdrs || [];
      }
    } catch (e) {
      console.error("Failed to fetch CDRs", e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchCDRs();
    const interval = setInterval(fetchCDRs, 10000);
    return () => clearInterval(interval);
  });

  const formatDuration = (seconds: number) => {
    const m = Math.floor((seconds || 0) / 60);
    const s = (seconds || 0) % 60;
    return `${m}m ${s}s`;
  };

  $: filteredCDRs = cdrs.filter(cdr => {
    const matchesSearch = 
      (cdr.CallerIDNumber || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      (cdr.CallerIDName || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      (cdr.Destination || '').toLowerCase().includes(searchTerm.toLowerCase());
    
    if (!matchesSearch) return false;

    if (statusFilter === 'completed') return cdr.HangupCause === 'NORMAL_CLEARING';
    if (statusFilter === 'cancelled') return cdr.HangupCause === 'ORIGINATOR_CANCEL';
    if (statusFilter === 'failed') return cdr.HangupCause && cdr.HangupCause !== 'NORMAL_CLEARING' && cdr.HangupCause !== 'ORIGINATOR_CANCEL';

    return true;
  });

  function exportCSV() {
    if (cdrs.length === 0) return;
    const headers = ['Time', 'Caller ID', 'Destination', 'Duration (s)', 'Billed (s)', 'Status'];
    const rows = cdrs.map(c => [
      c.CreatedAt,
      `"${c.CallerIDName || ''} <${c.CallerIDNumber || ''}>"`,
      `"${c.Destination || ''}"`,
      c.Duration || 0,
      c.Billsec || 0,
      c.HangupCause || 'UNKNOWN'
    ]);
    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map(e => e.join(','))].join('\n');
    const encodedUri = encodeURI(csvContent);
    const link = document.createElement('a');
    link.setAttribute('href', encodedUri);
    link.setAttribute('download', `pbx_cdr_export_${Date.now()}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500 ease-out max-w-7xl mx-auto w-full">
  
  <!-- Header -->
  <header class="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
    <div>
      <div class="flex items-center space-x-3 mb-1">
        <h1 class="text-3xl font-black tracking-tight text-white">Call Detail Records (CDRs)</h1>
        <span class="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-indigo-500/15 text-indigo-300 border border-indigo-500/30">
          Billing Ledger
        </span>
      </div>
      <p class="text-slate-400 text-sm">Real-time call accounting, carrier termination billing, and hangup diagnostics.</p>
    </div>

    <!-- Export & Refresh tools -->
    <div class="flex items-center space-x-3">
      <button 
        on:click={exportCSV}
        disabled={cdrs.length === 0}
        class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-800 transition-all flex items-center space-x-2 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer shadow-md"
      >
        <svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
        <span>Export CSV</span>
      </button>

      <button 
        on:click={fetchCDRs}
        class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-indigo-600 hover:bg-indigo-500 text-white transition-all flex items-center space-x-2 shadow-lg shadow-indigo-600/30 cursor-pointer"
      >
        <svg class="w-4 h-4 {isLoading ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
        <span>Refresh</span>
      </button>
    </div>
  </header>

  <!-- Filter Controls Card -->
  <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex flex-col md:flex-row gap-4 justify-between items-center backdrop-blur-xl shadow-xl">
    <!-- Search Bar -->
    <div class="relative w-full md:w-96">
      <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-500">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
      </div>
      <input 
        type="text" 
        bind:value={searchTerm}
        placeholder="Filter by caller, extension, or destination..." 
        class="w-full bg-slate-950/70 border border-slate-800 text-slate-100 placeholder-slate-500 text-xs rounded-xl pl-10 pr-4 py-2 focus:outline-none focus:border-indigo-500 font-mono"
      />
    </div>

    <!-- Status Filter Tabs -->
    <div class="flex items-center space-x-1.5 w-full md:w-auto overflow-x-auto">
      {#each [
        { id: 'all', label: 'All Calls' },
        { id: 'completed', label: 'Answered' },
        { id: 'cancelled', label: 'Cancelled' },
        { id: 'failed', label: 'Failed' }
      ] as tab}
        <button 
          on:click={() => statusFilter = tab.id as any}
          class="px-3 py-1.5 rounded-lg text-xs font-semibold whitespace-nowrap transition-colors {statusFilter === tab.id ? 'bg-indigo-600/30 text-indigo-300 border border-indigo-500/40' : 'text-slate-400 hover:text-white hover:bg-slate-800/60'}"
        >
          {tab.label}
        </button>
      {/each}
    </div>
  </div>

  <!-- CDRs Table Container -->
  <div class="bg-slate-900/80 border border-slate-800 rounded-2xl overflow-hidden shadow-2xl backdrop-blur-xl">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-950/80 text-slate-400 text-xs uppercase tracking-wider border-b border-slate-800">
            <th class="p-4 font-bold">Call Time</th>
            <th class="p-4 font-bold">Caller ID</th>
            <th class="p-4 font-bold">Destination</th>
            <th class="p-4 font-bold">Duration</th>
            <th class="p-4 font-bold">Hangup Cause</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 text-sm">
          {#if isLoading}
            <tr>
              <td colspan="5" class="p-12 text-center text-slate-400">
                <div class="flex flex-col items-center justify-center space-y-2">
                  <div class="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
                  <p class="text-sm font-medium">Loading call records from PostgreSQL...</p>
                </div>
              </td>
            </tr>
          {:else if filteredCDRs.length === 0}
            <tr>
              <td colspan="5" class="p-16 text-center">
                <div class="flex flex-col items-center justify-center space-y-3">
                  <div class="w-12 h-12 rounded-2xl bg-slate-800/80 border border-slate-700/80 flex items-center justify-center text-indigo-400 shadow-xl">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
                  </div>
                  <h3 class="text-base font-bold text-white">No Call Detail Records Found</h3>
                  <p class="text-slate-400 text-xs max-w-sm">No incoming or outgoing calls match your query. Place a test call to extension 1001 or dial *97 for voicemail to generate your first CDR.</p>
                </div>
              </td>
            </tr>
          {:else}
            {#each filteredCDRs as cdr}
              <tr class="hover:bg-slate-800/40 transition-colors">
                <td class="p-4 text-slate-400 font-mono text-xs">
                  {new Date(cdr.CreatedAt).toLocaleString()}
                </td>
                <td class="p-4 font-mono text-indigo-300 font-semibold">
                  {cdr.CallerIDName ? cdr.CallerIDName + ' ' : ''}&lt;{cdr.CallerIDNumber}&gt;
                </td>
                <td class="p-4 text-white font-mono font-medium">
                  {cdr.Destination}
                </td>
                <td class="p-4 text-slate-300">
                  <div class="flex items-center space-x-2">
                    <span>{formatDuration(cdr.Duration)}</span>
                    {#if cdr.Billsec > 0}
                      <span class="text-xs text-emerald-400 font-mono">({formatDuration(cdr.Billsec)} billed)</span>
                    {/if}
                  </div>
                </td>
                <td class="p-4">
                  {#if cdr.HangupCause === 'NORMAL_CLEARING'}
                    <Badge status="success">Completed</Badge>
                  {:else if cdr.HangupCause === 'ORIGINATOR_CANCEL'}
                    <Badge status="warning">Cancelled</Badge>
                  {:else if cdr.HangupCause === 'USER_BUSY'}
                    <Badge status="error">Busy</Badge>
                  {:else if cdr.HangupCause === 'NO_ANSWER'}
                    <Badge status="warning">No Answer</Badge>
                  {:else}
                    <Badge status="neutral">{cdr.HangupCause || 'Ended'}</Badge>
                  {/if}
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>

</div>
