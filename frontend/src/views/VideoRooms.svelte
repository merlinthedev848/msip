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

  interface VideoRoomItem {
    ID: string;
    Name: string;
    RoomCode: string;
    MaxUsers: number;
    IsActive: boolean;
  }

  let rooms: VideoRoomItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    RoomCode: '',
    MaxUsers: 10,
    IsActive: true
  };

  $: filtered = rooms.filter(r => 
    !search.trim() || `${r.Name} ${r.RoomCode}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ 'video-rooms': VideoRoomItem[] }>('/video-rooms');
      rooms = data['video-rooms'] || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      Name: '',
      RoomCode: 'conf-' + Math.floor(1000 + Math.random() * 9000),
      MaxUsers: 10,
      IsActive: true
    };
    modalOpen = true;
  }

  function openEdit(r: VideoRoomItem) {
    editingId = r.ID;
    form = {
      Name: r.Name,
      RoomCode: r.RoomCode,
      MaxUsers: r.MaxUsers || 10,
      IsActive: r.IsActive ?? true
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Room name is required');
    if (!form.RoomCode.trim()) return toast.error('Room code is required');

    saving = true;
    try {
      const payload = {
        Name: form.Name.trim(),
        RoomCode: form.RoomCode.trim(),
        MaxUsers: Number(form.MaxUsers) || 10,
        IsActive: form.IsActive
      };

      if (editingId) {
        await api.put(`/video-rooms/${editingId}`, payload);
        toast.success('Video room updated');
      } else {
        await api.post('/video-rooms', payload);
        toast.success('Video room created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(r: VideoRoomItem) {
    if (!confirm(`Delete video room "${r.Name}"?`)) return;
    try {
      await api.del(`/video-rooms/${r.ID}`);
      toast.success('Video room deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Video Conference Rooms" description="WebRTC multi-party video bridge endpoints.">
  <Button on:click={openCreate}>+ New Video Room</Button>
</PageHeader>

<Card padded={false} title="Configured Conference Bridges" subtitle="{rooms.length} video rooms">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search room…" bind:value={search} aria-label="Search rooms" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No video rooms match your search' : 'No video rooms configured'}
      message={search ? '' : 'Create a WebRTC video room for peer-to-peer or MCU conferencing.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Create first room</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Room Name</th>
            <th>Room Code</th>
            <th>Capacity</th>
            <th>Status</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{r.Name}</td>
              <td class="mono font-semibold text-sky-700 text-xs">{r.RoomCode}</td>
              <td class="text-xs text-slate-600">{r.MaxUsers} participants</td>
              <td><Badge tone={r.IsActive ? 'success' : 'neutral'}>{r.IsActive ? 'Active' : 'Disabled'}</Badge></td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(r)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(r)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Video Room' : 'New Video Room'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="room-name" class="label">Room Name</label>
      <input id="room-name" class="input" bind:value={form.Name} placeholder="e.g. Executive Boardroom" required />
    </div>

    <div>
      <label for="room-code" class="label">Room Code</label>
      <input id="room-code" class="input mono" bind:value={form.RoomCode} placeholder="e.g. conf-101" required />
    </div>

    <div>
      <label for="room-users" class="label">Max Participants</label>
      <input id="room-users" type="number" min="2" max="100" class="input mono" bind:value={form.MaxUsers} required />
    </div>

    <div class="flex items-center space-x-2 pt-1">
      <input id="room-active" type="checkbox" bind:checked={form.IsActive} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
      <label for="room-active" class="text-xs font-medium text-slate-700">Active</label>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
