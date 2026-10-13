<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface DomainItem {
    ID: string;
    Name: string;
    Domain: string;
    CreatedAt: string;
  }

  let domains: DomainItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = { Name: '', Domain: '' };

  $: filtered = domains.filter(d => 
    !search.trim() || `${d.Name} ${d.Domain}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ domains: DomainItem[] }>('/domains');
      domains = data.domains || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { Name: '', Domain: '' };
    modalOpen = true;
  }

  function openEdit(d: DomainItem) {
    editingId = d.ID;
    form = { Name: d.Name, Domain: d.Domain };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Tenant name is required');
    if (!form.Domain.trim()) return toast.error('Domain is required');

    saving = true;
    try {
      const payload = {
        Name: form.Name.trim(),
        Domain: form.Domain.trim().toLowerCase()
      };

      if (editingId) {
        await api.put(`/domains/${editingId}`, payload);
        toast.success('Tenant domain updated');
      } else {
        await api.post('/domains', payload);
        toast.success('Tenant domain created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(d: DomainItem) {
    if (!confirm(`Delete tenant domain "${d.Name}" (${d.Domain})? All associated extensions and routes will be deleted.`)) return;
    try {
      await api.del(`/domains/${d.ID}`);
      toast.success('Tenant domain deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Tenant Domains" description="Multi-tenant SIP isolation, domains, and organization partitions.">
  <Button on:click={openCreate}>+ New Tenant Domain</Button>
</PageHeader>

<Card padded={false} title="Configured Tenancies" subtitle="{domains.length} domains">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search tenant or domain…" bind:value={search} aria-label="Search domains" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No tenants match your search' : 'No tenant domains configured'}
      message={search ? '' : 'Add multi-tenant domains to partition extensions and routing per customer organization.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first domain</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Organization Name</th>
            <th>SIP Domain</th>
            <th>Created</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as d (d.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{d.Name}</td>
              <td class="mono font-semibold text-xs text-sky-700">{d.Domain}</td>
              <td class="text-xs text-slate-500">{formatDateTime(d.CreatedAt)}</td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(d)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(d)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Tenant Domain' : 'New Tenant Domain'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="dom-name" class="label">Organization Name</label>
      <input id="dom-name" class="input" bind:value={form.Name} placeholder="e.g. Acme Corporation" required />
    </div>

    <div>
      <label for="dom-domain" class="label">SIP Realm / Domain</label>
      <input id="dom-domain" class="input mono" bind:value={form.Domain} placeholder="e.g. acme.pbx.com" required />
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
