<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, formatDuration } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';

  interface CDRItem {
    ID: string;
    Duration: number;
    Billsec: number;
    HangupCause: string;
    Context: string;
    CreatedAt: string;
  }

  let cdrs: CDRItem[] = [];
  let loading = true;
  let totalCalls = 0;
  let answeredCalls = 0;
  let totalBillableSeconds = 0;
  let avgDuration = 0;
  let answerRate = 0;

  async function load() {
    loading = true;
    try {
      const data = await api.get<{ cdrs: CDRItem[] }>('/cdr');
      cdrs = data.cdrs || [];
      totalCalls = cdrs.length;
      answeredCalls = cdrs.filter(c => c.HangupCause === 'NORMAL_CLEARING').length;
      totalBillableSeconds = cdrs.reduce((acc, c) => acc + (c.Billsec || 0), 0);
      avgDuration = totalCalls > 0 ? Math.round(totalBillableSeconds / totalCalls) : 0;
      answerRate = totalCalls > 0 ? Math.round((answeredCalls / totalCalls) * 100) : 0;
    } catch (e) {
      //
    } finally {
      loading = false;
    }
  }

  onMount(load);
</script>

<PageHeader title="CDR Analytics & Telephony Reports" description="Aggregated call statistics, answer rates, and billable duration calculated from call logs.">
  <Button variant="secondary" on:click={load}>Refresh Data</Button>
</PageHeader>

<div class="space-y-6">
  <!-- Key Metrics Row -->
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Total Processed Calls</div>
      <div class="text-4xl font-extrabold text-slate-800 mono">{totalCalls}</div>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Answered Rate</div>
      <div class="text-4xl font-extrabold text-emerald-600 mono">{answerRate}%</div>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Total Billable Time</div>
      <div class="text-4xl font-extrabold text-sky-700 mono">{formatDuration(totalBillableSeconds)}</div>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs text-center">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1">Avg Call Duration</div>
      <div class="text-4xl font-extrabold text-slate-700 mono">{formatDuration(avgDuration)}</div>
    </div>
  </div>

  <!-- Traffic Summary by Context -->
  <Card title="Call Volume by Telephony Context">
    <div class="divide-y divide-slate-100 text-sm">
      <div class="py-3 flex items-center justify-between">
        <span class="font-medium text-slate-700">Internal Station-to-Station</span>
        <Badge tone="info">{cdrs.filter(c => c.Context === 'internal' || c.Context === 'default').length} calls</Badge>
      </div>
      <div class="py-3 flex items-center justify-between">
        <span class="font-medium text-slate-700">Inbound PSTN Calls</span>
        <Badge tone="success">{cdrs.filter(c => c.Context === 'public' || c.Context === 'inbound').length} calls</Badge>
      </div>
      <div class="py-3 flex items-center justify-between">
        <span class="font-medium text-slate-700">Outbound Gateway Calls</span>
        <Badge tone="neutral">{cdrs.filter(c => c.Context === 'outbound').length} calls</Badge>
      </div>
    </div>
  </Card>
</div>
