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

  interface ChatChannel {
    ID: string;
    Name: string;
    Type: string;
  }

  let chats: ChatChannel[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    Type: 'public'
  };

  $: filtered = chats.filter(c => 
    !search.trim() || c.Name.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ chats: ChatChannel[] }>('/chats');
      chats = data.chats || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { Name: '', Type: 'public' };
    modalOpen = true;
  }

  function openEdit(c: ChatChannel) {
    editingId = c.ID;
    form = { Name: c.Name, Type: c.Type };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Channel name is required');

    saving = true;
    try {
      const payload = {
        Name: form.Name.trim(),
        Type: form.Type
      };

      if (editingId) {
        await api.put(`/chats/${editingId}`, payload);
        toast.success('Channel updated');
      } else {
        await api.post('/chats', payload);
        toast.success('Channel created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(c: ChatChannel) {
    if (!confirm(`Delete channel #${c.Name}?`)) return;
    try {
      await api.del(`/chats/${c.ID}`);
      toast.success('Channel deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Team Chat Channels" description="Internal messaging channels and group broadcast rooms.">
  <Button on:click={openCreate}>+ New Channel</Button>
</PageHeader>

<Card padded={false} title="Configured Channels" subtitle="{chats.length} channels">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search channel…" bind:value={search} aria-label="Search channels" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No channels match your search' : 'No team chat channels'}
      message={search ? '' : 'Create chat channels for internal team collaboration.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Create first channel</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Channel Name</th>
            <th>Access</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as c (c.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">#{c.Name}</td>
              <td><Badge tone={c.Type === 'public' ? 'success' : 'neutral'}>{c.Type}</Badge></td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(c)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(c)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Channel' : 'New Chat Channel'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="chat-name" class="label">Channel Name</label>
      <input id="chat-name" class="input" bind:value={form.Name} placeholder="e.g. general or support" required />
    </div>

    <div>
      <label for="chat-type" class="label">Channel Type</label>
      <select id="chat-type" class="select" bind:value={form.Type}>
        <option value="public">Public</option>
        <option value="private">Private</option>
      </select>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
