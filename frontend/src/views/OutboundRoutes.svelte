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

  interface OutboundRouteItem {
    ID: string;
    DialPattern: string;
    TrunkID: string;
    Priority: number;
    CreatedAt: string;
  }

  interface TrunkItem {
    ID: string;
    ProviderName: string;
  }

  let routes: OutboundRouteItem[] = [];
  let trunks: TrunkItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    DialPattern: '',
    TrunkID: '',
    Priority: 1
  };

  $: filtered = routes.filter(r => 
    !search.trim() || r.DialPattern.toLowerCase().includes(search.trim().toLowerCase())
  );

  function getTrunkName(id: string) {
    const found = trunks.find(t => t.ID === id);
    return found ? found.ProviderName : id;
  }

  async function load() {
    loading = true;
    loadError = '';
    try {
      const [rData, tData] = await Promise.all([
        api.get<{ routes: OutboundRouteItem[] }>('/outbound-routes'),
        api.get<{ trunks: TrunkItem[] }>('/trunks')
      ]);
      routes = rData.routes || [];
      trunks = tData.trunks || [];
      if (trunks.length > 0 && !form.TrunkID) {
        form.TrunkID = trunks[0].ID;
      }
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      DialPattern: '',
      TrunkID: trunks[0]?.ID || '',
      Priority: 1
    };
    modalOpen = true;
  }

  function openEdit(r: OutboundRouteItem) {
    editingId = r.ID;
    form = {
      DialPattern: r.DialPattern,
      TrunkID: r.TrunkID,
      Priority: r.Priority
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.DialPattern.trim()) return toast.error('Dial pattern is required');
    if (!form.TrunkID) return toast.error('Please select an outbound trunk');

    saving = true;
    try {
      const payload = {
        DialPattern: form.DialPattern.trim(),
        TrunkID: form.TrunkID,
        Priority: Number(form.Priority) || 1
      };

      if (editingId) {
        await api.put(`/outbound-routes/${editingId}`, payload);
        toast.success('Outbound route updated');
      } else {
        await api.post('/outbound-routes', payload);
        toast.success('Outbound route created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(r: OutboundRouteItem) {
    if (!confirm(`Delete outbound route for ${r.DialPattern}?`)) return;
    try {
      await api.del(`/outbound-routes/${r.ID}`);
      toast.success('Outbound route deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Outbound Routes" description="Rules for matching dialed digits and routing calls out through carrier SIP trunks.">
  <Button on:click={openCreate}>+ New Outbound Route</Button>
</PageHeader>

<Card padded={false} title="Dial Routing Rules" subtitle="{routes.length} dial rules">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search pattern…" bind:value={search} aria-label="Search patterns" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No routes match your search' : 'No outbound routes configured'}
      message={search ? '' : 'Add a dial pattern matching external phone numbers to bridge through a SIP trunk.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first rule</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Dial Pattern</th>
            <th>Carrier Trunk</th>
            <th>Priority</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr>
              <td class="mono font-bold text-slate-800 text-sm">{r.DialPattern}</td>
              <td class="font-semibold text-slate-700">{getTrunkName(r.TrunkID)}</td>
              <td><Badge tone="info">Priority {r.Priority}</Badge></td>
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Outbound Route' : 'Add Outbound Route'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="out-pattern" class="label">Dial Pattern (e.g. _9NXXXXXX or _1XXXXXXXXXX)</label>
      <input id="out-pattern" class="input mono" bind:value={form.DialPattern} placeholder="e.g. _NXXXXXXXXX" required />
    </div>

    <div>
      <label for="out-trunk" class="label">Outbound SIP Trunk</label>
      <select id="out-trunk" class="select" bind:value={form.TrunkID} required>
        {#each trunks as t}
          <option value={t.ID}>{t.ProviderName}</option>
        {:else}
          <option value="" disabled>No trunks available — add a trunk first</option>
        {/each}
      </select>
    </div>

    <div>
      <label for="out-priority" class="label">Priority (Lower number = evaluated first)</label>
      <input id="out-priority" type="number" min="1" max="100" class="input mono" bind:value={form.Priority} required />
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving || trunks.length === 0}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
