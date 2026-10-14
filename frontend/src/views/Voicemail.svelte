<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, formatDuration } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Modal from '../components/ui/Modal.svelte';
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
  let playingAudioId: string | null = null;

  // Voicemail Server & AI Settings Modal
  let settingsModalOpen = false;
  let settingsSaving = false;
  let settingsForm = {
    auto_transcribe_voicemail: true,
    stt_provider: 'whisper',
    sentiment_analysis: true,
    email_transcripts: false,
    notification_email: ''
  };

  $: filtered = voicemails.filter(v => 
    !search.trim() || 
    `${v.Extension} ${v.CallerID} ${v.Transcription || ''} ${v.Summary || ''}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function openSettings() {
    try {
      const res = await api.get<{ settings: any }>('/ai/settings');
      if (res.settings) {
        settingsForm = {
          auto_transcribe_voicemail: res.settings.auto_transcribe_voicemail ?? true,
          stt_provider: res.settings.stt_provider || 'whisper',
          sentiment_analysis: res.settings.sentiment_analysis ?? true,
          email_transcripts: res.settings.email_transcripts ?? false,
          notification_email: res.settings.notification_email || ''
        };
      }
    } catch {
      // Use defaults if not found
    }
    settingsModalOpen = true;
  }

  async function saveSettings(e: Event) {
    e.preventDefault();
    settingsSaving = true;
    try {
      await api.put('/ai/settings', settingsForm);
      toast.success('Voicemail server & AI settings saved');
      settingsModalOpen = false;
    } catch (err: any) {
      toastError(err);
    } finally {
      settingsSaving = false;
    }
  }

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
  <div class="flex gap-2">
    <Button variant="secondary" on:click={openSettings}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
      Server & AI Settings
    </Button>
    <Button variant="secondary" on:click={load}>Refresh</Button>
  </div>
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

                {#if playingAudioId === v.ID}
                  <div class="mt-2 p-2 bg-slate-100 rounded flex items-center gap-2">
                    <audio controls class="h-8 w-full max-w-xs" src="/api/v1/voicemails/{v.ID}/audio" preload="none">
                      Your browser does not support the audio player.
                    </audio>
                  </div>
                {/if}
              </td>
              <td class="text-xs text-slate-500 whitespace-nowrap align-top pt-3">{formatDateTime(v.CreatedAt)}</td>
              <td class="text-right space-x-3 whitespace-nowrap align-top pt-3">
                <button 
                  type="button" 
                  class="link-action text-sky-600 hover:text-sky-800"
                  on:click={() => playingAudioId = playingAudioId === v.ID ? null : v.ID}
                >
                  {playingAudioId === v.ID ? 'Hide Audio' : 'Play Audio'}
                </button>
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

<!-- Voicemail Server & AI Settings Modal -->
<Modal
  open={settingsModalOpen}
  title="Voicemail Server & AI Transcription Settings"
  on:close={() => settingsModalOpen = false}
>
  <form on:submit={saveSettings} class="space-y-4 text-xs">
    <div class="space-y-3 p-4 bg-slate-50 border border-slate-200 rounded-lg">
      <div class="flex items-start space-x-3">
        <input
          id="vm-auto-transcribe"
          type="checkbox"
          bind:checked={settingsForm.auto_transcribe_voicemail}
          class="mt-1 h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
        />
        <div>
          <label for="vm-auto-transcribe" class="font-semibold text-slate-800 cursor-pointer">
            Auto-Transcribe New Voicemails with AI (Whisper)
          </label>
          <p class="text-slate-500 text-[11px] mt-0.5">
            Immediately run speech-to-text upon voicemail completion and generate executive summaries.
          </p>
        </div>
      </div>

      <div class="flex items-start space-x-3">
        <input
          id="vm-sentiment"
          type="checkbox"
          bind:checked={settingsForm.sentiment_analysis}
          class="mt-1 h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
        />
        <div>
          <label for="vm-sentiment" class="font-semibold text-slate-800 cursor-pointer">
            Perform Sentiment Analysis
          </label>
          <p class="text-slate-500 text-[11px] mt-0.5">
            Classify voicemail tone as positive, neutral, or urgent/negative for supervisor review.
          </p>
        </div>
      </div>

      <div class="flex items-start space-x-3">
        <input
          id="vm-email-transcripts"
          type="checkbox"
          bind:checked={settingsForm.email_transcripts}
          class="mt-1 h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
        />
        <div>
          <label for="vm-email-transcripts" class="font-semibold text-slate-800 cursor-pointer">
            Email Audio & Transcript to Recipient
          </label>
          <p class="text-slate-500 text-[11px] mt-0.5">
            Dispatch MP3 attachment and full text transcript to the extension user's email address.
          </p>
        </div>
      </div>
    </div>

    <div>
      <label for="vm-stt" class="block font-medium text-slate-700 mb-1">Speech Recognition (STT) Provider</label>
      <select
        id="vm-stt"
        bind:value={settingsForm.stt_provider}
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
      >
        <option value="whisper">OpenAI Whisper (Cloud / Local)</option>
        <option value="deepgram">Deepgram Nova-2 (Ultra-fast)</option>
        <option value="google">Google Cloud Speech-to-Text</option>
        <option value="vosk">Vosk (On-Premise / Offline)</option>
      </select>
    </div>

    <div>
      <label for="vm-admin-email" class="block font-medium text-slate-700 mb-1">Central Compliance Notification Email (Optional)</label>
      <input
        id="vm-admin-email"
        type="email"
        bind:value={settingsForm.notification_email}
        placeholder="voicemail-archive@company.com"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
      />
      <p class="text-[11px] text-slate-500 mt-1">Receive carbon-copy (BCC) transcripts for compliance audits.</p>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => settingsModalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={settingsSaving}>
        {settingsSaving ? 'Saving...' : 'Save Voicemail Settings'}
      </Button>
    </div>
  </form>
</Modal>
