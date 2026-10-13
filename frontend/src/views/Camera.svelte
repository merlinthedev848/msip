<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface CameraItem {
    ID: string;
    Name: string;
    RTSPUrl: string;
    Location: string;
  }

  let cameras: CameraItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    RTSPUrl: '',
    Location: ''
  };

  $: filtered = cameras.filter(c => 
    !search.trim() || `${c.Name} ${c.Location} ${c.RTSPUrl}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ cameras: CameraItem[] }>('/cameras');
      cameras = data.cameras || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { Name: '', RTSPUrl: 'rtsp://', Location: '' };
    modalOpen = true;
  }

  function openEdit(c: CameraItem) {
    editingId = c.ID;
    form = {
      Name: c.Name,
      RTSPUrl: c.RTSPUrl,
      Location: c.Location || ''
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Camera name is required');
    if (!form.RTSPUrl.trim()) return toast.error('RTSP URL is required');

    saving = true;
    try {
      const payload = {
        Name: form.Name.trim(),
        RTSPUrl: form.RTSPUrl.trim(),
        Location: form.Location.trim()
      };

      if (editingId) {
        await api.put(`/cameras/${editingId}`, payload);
        toast.success('Camera updated');
      } else {
        await api.post('/cameras', payload);
        toast.success('Camera registered');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(c: CameraItem) {
    if (!confirm(`Delete IP camera "${c.Name}"?`)) return;
    try {
      await api.del(`/cameras/${c.ID}`);
      toast.success('Camera deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="IP Camera RTSP Feeds" description="Surveillance cameras and video intercom streaming endpoints.">
  <Button on:click={openCreate}>+ Register Camera</Button>
</PageHeader>

<Card padded={false} title="Surveillance Endpoints" subtitle="{cameras.length} camera streams">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search camera or location…" bind:value={search} aria-label="Search cameras" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No cameras match your search' : 'No IP cameras registered'}
      message={search ? '' : 'Register RTSP camera streams for security monitoring or door phone video.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Register first camera</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Camera Name</th>
            <th>Location</th>
            <th>RTSP URL</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as c (c.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{c.Name}</td>
              <td class="text-xs text-slate-700 font-medium">{c.Location || '—'}</td>
              <td class="mono text-xs text-slate-500 max-w-xs truncate">{c.RTSPUrl}</td>
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Camera' : 'Register IP Camera'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="cam-name" class="label">Camera Name</label>
      <input id="cam-name" class="input" bind:value={form.Name} placeholder="e.g. Front Gate Intercom" required />
    </div>

    <div>
      <label for="cam-loc" class="label">Physical Location</label>
      <input id="cam-loc" class="input" bind:value={form.Location} placeholder="e.g. Building A Entrance" />
    </div>

    <div>
      <label for="cam-url" class="label">RTSP Stream URL</label>
      <input id="cam-url" class="input mono" bind:value={form.RTSPUrl} placeholder="rtsp://admin:pass@192.168.1.50:554/stream1" required />
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
