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
    CreatedAt: string;
  }

  let voicemails: VoicemailItem[] = [];
  let loading = true;
  let loadError = '';
  let search = '';

  $: filtered = voicemails.filter(v => 
    !search.trim() || `${v.Extension} ${v.CallerID}`.toLowerCase().includes(search.trim().toLowerCase())
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

  onMount(load);
</script>

<PageHeader title="Voicemail Server" description="Stored audio messages, duration, and mailbox states.">
  <Button variant="secondary" on:click={load}>Refresh</Button>
</PageHeader>

<Card padded={false} title="Mailbox Messages" subtitle="{voicemails.length} recordings">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search extension or caller…" bind:value={search} aria-label="Search messages" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No voicemails match your search' : 'No voicemails recorded'}
      message={search ? '' : 'Voicemail recordings left on PBX extensions will appear here.'}
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
            <th>Received</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as v (v.ID)}
            <tr>
              <td class="mono font-bold text-slate-800 text-sm">{v.Extension}</td>
              <td class="mono text-xs text-slate-700">{v.CallerID || 'Unknown'}</td>
              <td class="mono text-xs text-slate-600">{formatDuration(v.Duration)}</td>
              <td>
                <Badge tone={v.IsRead ? 'neutral' : 'warning'}>
                  {v.IsRead ? 'Read' : 'New'}
                </Badge>
              </td>
              <td class="text-xs text-slate-500">{formatDateTime(v.CreatedAt)}</td>
              <td class="text-right space-x-3 whitespace-nowrap">
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
