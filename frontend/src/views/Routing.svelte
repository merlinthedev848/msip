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

  interface InboundRouteItem {
    ID: string;
    DIDNumber: string;
    DestinationType: string;
    DestinationTarget: string;
    CreatedAt: string;
  }

  let routes: InboundRouteItem[] = [];
  let availableExts: any[] = [];
  let availableIvrs: any[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    DIDNumber: '',
    DestinationType: 'EXTENSION',
    DestinationTarget: ''
  };

  $: filtered = routes.filter(r => 
    !search.trim() || `${r.DIDNumber} ${r.DestinationTarget} ${r.DestinationType}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const [routeData, extData, ivrData] = await Promise.all([
        api.get<{ routes: InboundRouteItem[] }>('/inbound-routes'),
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] })),
        api.get<{ ivrs: any[] }>('/ivrs').catch(() => ({ ivrs: [] }))
      ]);
      routes = routeData.routes || [];
      availableExts = extData.extensions || [];
      availableIvrs = ivrData.ivrs || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    const defaultTarget = availableExts.length > 0 ? availableExts[0].ExtensionNumber : '';
    form = { DIDNumber: '', DestinationType: 'EXTENSION', DestinationTarget: defaultTarget };
    modalOpen = true;
  }

  function openEdit(r: InboundRouteItem) {
    editingId = r.ID;
    form = {
      DIDNumber: r.DIDNumber,
      DestinationType: r.DestinationType,
      DestinationTarget: r.DestinationTarget
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.DIDNumber.trim()) return toast.error('DID Number or pattern is required');
    if (!form.DestinationTarget.trim()) return toast.error('Destination Target is required');

    saving = true;
    try {
      const payload = {
        DIDNumber: form.DIDNumber.trim(),
        DestinationType: form.DestinationType,
        DestinationTarget: form.DestinationTarget.trim()
      };

      if (editingId) {
        await api.put(`/inbound-routes/${editingId}`, payload);
        toast.success('Inbound route updated');
      } else {
        await api.post('/inbound-routes', payload);
        toast.success('Inbound route created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(r: InboundRouteItem) {
    if (!confirm(`Delete inbound route for ${r.DIDNumber}?`)) return;
    try {
      await api.del(`/inbound-routes/${r.ID}`);
      toast.success('Inbound route deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Inbound Routes" description="Map incoming DID telephone numbers directly to PBX extensions, IVRs, or queues.">
  <Button on:click={openCreate}>+ New Inbound Route</Button>
</PageHeader>

<Card padded={false} title="Inbound Dialplan" subtitle="{routes.length} routing rules">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search DID or target…" bind:value={search} aria-label="Search routes" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No routes match your search' : 'No inbound routes configured'}
      message={search ? '' : 'Add a DID number mapping to route incoming carrier calls into FreeSWITCH.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first route</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>DID Pattern</th>
            <th>Type</th>
            <th>Target Destination</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr>
              <td class="mono font-bold text-slate-800 text-sm">{r.DIDNumber}</td>
              <td><Badge tone="info">{r.DestinationType}</Badge></td>
              <td class="mono text-xs text-slate-700 font-semibold">{r.DestinationTarget}</td>
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Incoming Call Rule' : 'New Incoming Call Rule'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="route-did" class="label">Incoming Phone Number (DID)</label>
      <input id="route-did" class="input mono" bind:value={form.DIDNumber} placeholder="e.g. +14155550199 or 4155550199" required />
      <span class="help">The telephone number assigned to you by your SIP Trunk carrier.</span>
    </div>

    <div>
      <label for="route-type" class="label">Forward Call To</label>
      <select id="route-type" class="select" bind:value={form.DestinationType}>
        <option value="EXTENSION">Extension (Specific User Handset)</option>
        <option value="IVR">Auto Attendant (Press 1 for Sales IVR)</option>
        <option value="QUEUE">Call Queue / Wallboard</option>
        <option value="VOICEMAIL">Voicemail Server</option>
      </select>
    </div>

    <div>
      <label for="route-target" class="label">Destination Target</label>
      {#if form.DestinationType === 'EXTENSION' && availableExts.length > 0}
        <select id="route-target" class="select" bind:value={form.DestinationTarget} required>
          <option value="" disabled>Choose an extension...</option>
          {#each availableExts as ext}
            <option value={ext.ExtensionNumber}>Extension {ext.ExtensionNumber}</option>
          {/each}
        </select>
      {:else if form.DestinationType === 'IVR' && availableIvrs.length > 0}
        <select id="route-target" class="select" bind:value={form.DestinationTarget} required>
          <option value="" disabled>Choose an Auto Attendant...</option>
          {#each availableIvrs as ivr}
            <option value={ivr.Name}>{ivr.Name} (Ext {ivr.Extension})</option>
          {/each}
        </select>
      {:else}
        <input id="route-target" class="input mono" bind:value={form.DestinationTarget} placeholder="e.g. 1001" required />
      {/if}
      <span class="help">Where the phone rings when an outside caller dials this number.</span>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save Rule'}</Button>
    </div>
  </form>
</Modal>
