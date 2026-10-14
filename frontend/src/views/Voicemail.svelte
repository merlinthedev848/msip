<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, formatDuration } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface VoicemailItem {
    ID: string;
    Extension: string;
    CallerID: string;
    Duration: number;
    FilePath: string;
    IsRead: boolean;
    Transcription?: string;
    Summary?: string;
    Sentiment?: string;
    CreatedAt: string;
  }

  let voicemails: VoicemailItem[] = [];
  let loading = true;
  let loadError = '';
  let search = '';
  let transcribing: Record<string, boolean> = {};
  let expandedVoicemail: string | null = null;

  $: filtered = voicemails.filter(v => 
    !search.trim() || 
    `${v.Extension} ${v.CallerID} ${v.Transcription || ''} ${v.Summary || ''}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ voicemails: VoicemailItem[] }>('/voicemails');
      voicemails = data.voicemails || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  async function toggleRead(v: VoicemailItem) {
    try {
      await api.put(`/voicemails/${v.ID}`, { IsRead: !v.IsRead });
      toast.success(v.IsRead ? 'Marked unread' : 'Marked read');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  async function remove(v: VoicemailItem) {
    if (!confirm(`Delete voicemail from ${v.CallerID}?`)) return;
    try {
      await api.del(`/voicemails/${v.ID}`);
      toast.success('Voicemail deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  async function transcribe(v: VoicemailItem) {
    transcribing[v.ID] = true;
    try {
      const res = await api.post<{ transcription: string; summary: string; sentiment: string }>(`/voicemails/${v.ID}/transcribe`, {});
      v.Transcription = res.transcription;
      v.Summary = res.summary;
      v.Sentiment = res.sentiment;
      voicemails = [...voicemails];
      expandedVoicemail = v.ID;
      toast.success('Voicemail transcribed successfully via AI engine');
    } catch (e) {
      toastError(e);
    } finally {
      transcribing[v.ID] = false;
    }
  }

  function getSentimentTone(sentiment?: string): 'success' | 'warning' | 'danger' | 'neutral' {
    if (!sentiment) return 'neutral';
    const s = sentiment.toLowerCase();
    if (s.includes('pos')) return 'success';
    if (s.includes('neg')) return 'danger';
    return 'neutral';
  }

  onMount(load);
</script>

<PageHeader title="Voicemail Server" description="Stored audio messages, Whisper AI transcription, and caller sentiment analysis.">
  <Button variant="secondary" on:click={load}>Refresh</Button>
</PageHeader>

<Card padded={false} title="Mailbox Messages" subtitle="{voicemails.length} recordings with AI transcriptions">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-64" placeholder="Search extension, caller, transcript…" bind:value={search} aria-label="Search messages" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No voicemails match your search' : 'No voicemails recorded'}
      message={search ? '' : 'Voicemail recordings left on PBX extensions will appear here with automatic transcription.'}
    >
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Extension</th>
            <th>Caller ID</th>
            <th>Duration</th>
            <th>Status</th>
            <th>AI Transcript & Sentiment</th>
            <th>Received</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as v (v.ID)}
            <tr class="hover:bg-slate-50/70 transition-colors">
              <td class="mono font-bold text-slate-800 text-sm align-top pt-3">{v.Extension}</td>
              <td class="mono text-xs text-slate-700 align-top pt-3">{v.CallerID || 'Unknown'}</td>
              <td class="mono text-xs text-slate-600 align-top pt-3">{formatDuration(v.Duration)}</td>
              <td class="align-top pt-3">
                <Badge tone={v.IsRead ? 'neutral' : 'warning'}>
                  {v.IsRead ? 'Read' : 'New'}
                </Badge>
              </td>
              <td class="max-w-md align-top pt-3">
                {#if v.Transcription}
                  <div class="space-y-1.5">
                    <div class="flex items-center space-x-2">
                      <span class="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-semibold bg-indigo-50 text-indigo-700 border border-indigo-200">
                        <svg class="w-3 h-3 mr-1 text-indigo-500" fill="currentColor" viewBox="0 0 20 20"><path d="M7 4a3 3 0 016 0v4a3 3 0 11-6 0V4zm4 10.93A7.001 7.001 0 0017 8a1 1 0 10-2 0A5 5 0 015 8a1 1 0 00-2 0 7.001 7.001 0 006 6.93V17H6a1 1 0 100 2h8a1 1 0 100-2h-3v-2.07z"/></svg>
                        Whisper AI
                      </span>
                      {#if v.Sentiment}
                        <Badge tone={getSentimentTone(v.Sentiment)}>
                          {v.Sentiment.toUpperCase()}
                        </Badge>
                      {/if}
                    </div>
                    {#if v.Summary}
                      <p class="text-xs text-slate-800 font-medium line-clamp-1">{v.Summary}</p>
                    {/if}
                    <p class="text-xs text-slate-600 line-clamp-2 leading-relaxed bg-slate-50 p-2 rounded border border-slate-100 font-sans">
                      {v.Transcription}
                    </p>
                  </div>
                {:else}
                  <div class="flex items-center space-x-2 py-1">
                    <span class="text-xs text-slate-400 italic">Not transcribed yet</span>
                    <button
                      type="button"
                      class="text-[11px] font-medium text-indigo-600 hover:text-indigo-800 hover:underline flex items-center"
                      disabled={transcribing[v.ID]}
                      on:click={() => transcribe(v)}
                    >
                      {#if transcribing[v.ID]}
                        <span class="animate-spin inline-block w-3 h-3 border border-indigo-600 border-t-transparent rounded-full mr-1"></span>
                        Transcribing…
                      {:else}
                        Transcribe with AI →
                      {/if}
                    </button>
                  </div>
                {/if}
              </td>
              <td class="text-xs text-slate-500 whitespace-nowrap align-top pt-3">{formatDateTime(v.CreatedAt)}</td>
              <td class="text-right space-x-3 whitespace-nowrap align-top pt-3">
                {#if !v.Transcription}
                  <button 
                    type="button" 
                    class="link-action text-indigo-600 hover:text-indigo-800 font-medium" 
                    disabled={transcribing[v.ID]}
                    on:click={() => transcribe(v)}
                  >
                    {transcribing[v.ID] ? 'Processing…' : 'Transcribe'}
                  </button>
                {/if}
                <button type="button" class="link-action" on:click={() => toggleRead(v)}>
                  {v.IsRead ? 'Mark Unread' : 'Mark Read'}
                </button>
                <button type="button" class="link-danger" on:click={() => remove(v)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>
