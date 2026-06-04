<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import { onMount } from 'svelte';

  let cdrs = [];
  let isLoading = true;

  async function fetchCDRs() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/cdr`, {
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
    const interval = setInterval(fetchCDRs, 10000); // Auto-refresh
    return () => clearInterval(interval);
  });

  const formatDuration = (seconds) => {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}m ${s}s`;
  };
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">Call History (CDRs)</h1>
      <p class="text-slate-500 text-sm">Real-time billing and call detail records.</p>
    </div>
  </header>

  <GlassCard className="p-0 overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-white/40 text-slate-500 text-xs uppercase tracking-wider">
            <th class="p-4 font-medium">Time</th>
            <th class="p-4 font-medium">Caller ID</th>
            <th class="p-4 font-medium">Destination</th>
            <th class="p-4 font-medium">Duration</th>
            <th class="p-4 font-medium">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-800/50 text-sm text-slate-700">
          {#if isLoading}
            <tr><td colspan="5" class="p-8 text-center text-slate-500">Loading CDRs...</td></tr>
          {:else if cdrs.length === 0}
            <tr><td colspan="5" class="p-8 text-center text-slate-500">No calls recorded yet. Make a call to see records here!</td></tr>
          {:else}
            {#each cdrs as cdr}
              <tr class="hover:bg-slate-100/30 transition-colors">
                <td class="p-4 text-slate-500">{new Date(cdr.CreatedAt).toLocaleString()}</td>
                <td class="p-4 font-mono text-blue-600 font-medium">
                  {cdr.CallerIDName ? cdr.CallerIDName + ' ' : ''}&lt;{cdr.CallerIDNumber}&gt;
                </td>
                <td class="p-4 text-slate-900 font-mono">{cdr.Destination}</td>
                <td class="p-4 text-slate-700">
                  {formatDuration(cdr.Duration)}
                  {#if cdr.Billsec > 0}
                    <span class="text-xs text-emerald-400 ml-2">({formatDuration(cdr.Billsec)} billed)</span>
                  {/if}
                </td>
                <td class="p-4">
                  {#if cdr.HangupCause === 'NORMAL_CLEARING'}
                    <Badge status="success">Completed</Badge>
                  {:else if cdr.HangupCause === 'ORIGINATOR_CANCEL'}
                    <Badge status="neutral">Canceled</Badge>
                  {:else}
                    <Badge status="error">{cdr.HangupCause}</Badge>
                  {/if}
                </td>
              </tr>
              {#if cdr.Transcription}
              <tr class="bg-indigo-900/10 border-b border-slate-200/50">
                <td colspan="5" class="p-4 pl-8">
                  <div class="flex items-start space-x-3">
                    <svg class="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
                    <div>
                      <p class="text-xs font-bold text-blue-600 uppercase tracking-wider mb-1">AI Dictation / Transcription</p>
                      <p class="text-sm text-slate-700 italic">"{cdr.Transcription}"</p>
                    </div>
                  </div>
                </td>
              </tr>
              {/if}
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </GlassCard>
</div>

