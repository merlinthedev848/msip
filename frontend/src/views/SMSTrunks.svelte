<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface SMSItem {
    ID: string;
    Direction: string;
    Sender: string;
    Receiver: string;
    Body: string;
    CreatedAt: string;
  }

  let messages: SMSItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let form = {
    Receiver: '',
    Sender: '',
    Body: ''
  };

  $: filtered = messages.filter(m => 
    !search.trim() || `${m.Sender} ${m.Receiver} ${m.Body}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ sms: SMSItem[] }>('/sms');
      messages = data.sms || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    form = { Receiver: '', Sender: '', Body: '' };
    modalOpen = true;
  }

  async function sendSMS(e: Event) {
    e.preventDefault();
    if (!form.Receiver.trim()) return toast.error('Receiver number is required');
    if (!form.Body.trim()) return toast.error('Message body is required');

    saving = true;
    try {
      await api.post('/sms', {
        Direction: 'outbound',
        Sender: form.Sender.trim() || 'PBX',
        Receiver: form.Receiver.trim(),
        Body: form.Body.trim()
      });
      toast.success('SMS dispatched');
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(m: SMSItem) {
    if (!confirm('Delete this SMS log record?')) return;
    try {
      await api.del(`/sms/${m.ID}`);
      toast.success('SMS record deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="SMS & Messaging" description="Inbound and outbound text message logs and dispatch.">
  <Button on:click={openCreate}>+ Send SMS</Button>
</PageHeader>

<Card padded={false} title="Message Log" subtitle="{messages.length} messages">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search message…" bind:value={search} aria-label="Search SMS" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No messages match your search' : 'No SMS messages logged'}
      message={search ? '' : 'Inbound and outbound text messages will appear here.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Send first SMS</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Direction</th>
            <th>Sender</th>
            <th>Recipient</th>
            <th>Message Body</th>
            <th>Date</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as m (m.ID)}
            <tr>
              <td>
                <Badge tone={m.Direction === 'inbound' ? 'info' : 'success'}>
                  {m.Direction}
                </Badge>
              </td>
              <td class="mono font-semibold text-xs text-slate-800">{m.Sender}</td>
              <td class="mono font-semibold text-xs text-sky-700">{m.Receiver}</td>
              <td class="text-xs text-slate-700 max-w-sm truncate">{m.Body}</td>
              <td class="text-xs text-slate-500">{formatDateTime(m.CreatedAt)}</td>
              <td class="text-right whitespace-nowrap">
                <button type="button" class="link-danger" on:click={() => remove(m)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title="Send Text Message">
  <form on:submit={sendSMS} class="space-y-4">
    <div>
      <label for="sms-sender" class="label">Sender DID / Caller ID</label>
      <input id="sms-sender" class="input mono" bind:value={form.Sender} placeholder="e.g. +14155550199" />
    </div>

    <div>
      <label for="sms-receiver" class="label">Recipient Phone Number</label>
      <input id="sms-receiver" class="input mono" bind:value={form.Receiver} placeholder="e.g. +14155550122" required />
    </div>

    <div>
      <label for="sms-body" class="label">Message Body</label>
      <textarea id="sms-body" rows="4" class="textarea" bind:value={form.Body} placeholder="Type message text…" required></textarea>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Dispatching…' : 'Send'}</Button>
    </div>
  </form>
</Modal>
