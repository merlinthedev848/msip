<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, formatDuration } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface CDRItem {
    ID: string;
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

  let cdrs: CDRItem[] = [];
  let loading = true;
  let loadError = '';
  let search = '';
  let filterContext = 'ALL';

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

  function exportCSV() {
    if (cdrs.length === 0) return;
    const headers = ['Time', 'Caller ID', 'Destination', 'Context', 'Duration (s)', 'Billable (s)', 'Disposition'];
    const rows = cdrs.map(c => [
      c.CreatedAt,
      `"${c.CallerIDNumber || ''}"`,
      `"${c.Destination || ''}"`,
      c.Context || '',
      c.Duration || 0,
      c.Billsec || 0,
      c.HangupCause || ''
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

<PageHeader title="Call Detail Records (CDRs)" description="Call logs recorded directly by the FreeSWITCH engine and PostgreSQL.">
  <div class="flex items-center gap-2">
    <Button variant="secondary" on:click={exportCSV} disabled={cdrs.length === 0}>Export CSV</Button>
    <Button variant="secondary" on:click={load}>Refresh</Button>
  </div>
</PageHeader>

<Card padded={false} title="Call History" subtitle="{cdrs.length} records">
  <svelte:fragment slot="actions">
    <div class="flex items-center gap-2">
      <select class="select !py-1 !text-xs w-32" bind:value={filterContext} aria-label="Filter context">
        <option value="ALL">All Contexts</option>
        <option value="internal">Internal</option>
        <option value="inbound">Inbound</option>
        <option value="outbound">Outbound</option>
      </select>
      <input class="input !py-1.5 !text-xs w-48" placeholder="Search caller, dest…" bind:value={search} aria-label="Search CDRs" />
    </div>
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No call records match your search' : 'No call records logged yet'}
      message={search ? '' : 'Completed SIP calls processed through Kamailio and FreeSWITCH will appear here.'}
    >
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Date & Time</th>
            <th>Caller ID</th>
            <th>Destination</th>
            <th>Context</th>
            <th>Duration</th>
            <th>Disposition</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as c (c.ID)}
            <tr>
              <td class="text-xs text-slate-500">{formatDateTime(c.CreatedAt)}</td>
              <td class="mono font-semibold text-slate-800 text-xs">
                {c.CallerIDNumber || 'Anonymous'}
                {#if c.CallerIDName && c.CallerIDName !== c.CallerIDNumber}
                  <span class="text-slate-400 font-normal">({c.CallerIDName})</span>
                {/if}
              </td>
              <td class="mono text-xs font-semibold text-sky-700">{c.Destination}</td>
              <td><Badge tone="info">{c.Context || 'default'}</Badge></td>
              <td class="mono text-xs text-slate-700">{formatDuration(c.Billsec || c.Duration)}</td>
              <td>
                <Badge tone={c.HangupCause === 'NORMAL_CLEARING' ? 'success' : 'error'}>
                  {c.HangupCause === 'NORMAL_CLEARING' ? 'Answered' : (c.HangupCause || 'Failed')}
                </Badge>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>
