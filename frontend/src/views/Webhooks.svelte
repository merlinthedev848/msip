<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface WebhookItem {
    ID: string;
    Event: string;
    TargetUrl: string;
  }

  let webhooks: WebhookItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Event: 'call.completed',
    TargetUrl: ''
  };

  $: filtered = webhooks.filter(w => 
    !search.trim() || `${w.Event} ${w.TargetUrl}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ webhooks: WebhookItem[] }>('/webhooks');
      webhooks = data.webhooks || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { Event: 'call.completed', TargetUrl: 'https://' };
    modalOpen = true;
  }

  function openEdit(w: WebhookItem) {
    editingId = w.ID;
    form = {
      Event: w.Event,
      TargetUrl: w.TargetUrl
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.TargetUrl.trim()) return toast.error('Target URL is required');

    saving = true;
    try {
      const payload = {
        Event: form.Event,
        TargetUrl: form.TargetUrl.trim()
      };

      if (editingId) {
        await api.put(`/webhooks/${editingId}`, payload);
        toast.success('Webhook updated');
      } else {
        await api.post('/webhooks', payload);
        toast.success('Webhook registered');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(w: WebhookItem) {
    if (!confirm(`Delete webhook for ${w.Event}?`)) return;
    try {
      await api.del(`/webhooks/${w.ID}`);
      toast.success('Webhook deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Event Webhooks" description="Real-time HTTP POST notifications dispatched on call state changes.">
  <Button on:click={openCreate}>+ Register Webhook</Button>
</PageHeader>

<Card padded={false} title="Active Subscriptions" subtitle="{webhooks.length} webhooks">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search event or url…" bind:value={search} aria-label="Search webhooks" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No webhooks match your search' : 'No webhooks registered'}
      message={search ? '' : 'Register webhook URLs to receive instant HTTP callbacks when calls start or complete.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Register first webhook</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Event Topic</th>
            <th>Target Endpoint URL</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as w (w.ID)}
            <tr>
              <td><Badge tone="info">{w.Event}</Badge></td>
              <td class="mono text-xs text-slate-700 max-w-md truncate">{w.TargetUrl}</td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(w)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(w)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Webhook' : 'Register Webhook'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="hook-event" class="label">Event Trigger</label>
      <select id="hook-event" class="select" bind:value={form.Event}>
        <option value="call.completed">call.completed (FreeSWITCH CDR saved)</option>
        <option value="call.started">call.started (Channel answered)</option>
        <option value="voicemail.new">voicemail.new (New audio message)</option>
        <option value="sms.received">sms.received (Inbound SMS)</option>
      </select>
    </div>

    <div>
      <label for="hook-url" class="label">Target Endpoint URL</label>
      <input id="hook-url" class="input mono" bind:value={form.TargetUrl} placeholder="https://api.yourdomain.com/webhooks/pbx" required />
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
