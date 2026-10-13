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

  interface FirewallRule {
    ID: string;
    IPAddress: string;
    Action: string;
    Notes: string;
    CreatedAt: string;
  }

  let rules: FirewallRule[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = { IPAddress: '', Action: 'accept', Notes: '' };

  $: filtered = rules.filter((r) =>
    !search.trim() || `${r.IPAddress} ${r.Notes} ${r.Action}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  const isBlock = (a: string) => ['drop', 'block', 'reject', 'deny'].includes((a || '').toLowerCase());

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ firewall: FirewallRule[] }>('/firewall');
      rules = data.firewall || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { IPAddress: '', Action: 'accept', Notes: '' };
    modalOpen = true;
  }

  function openEdit(r: FirewallRule) {
    editingId = r.ID;
    form = { IPAddress: r.IPAddress, Action: isBlock(r.Action) ? 'drop' : 'accept', Notes: r.Notes || '' };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.IPAddress.trim()) return toast.error('IP address or CIDR is required');
    saving = true;
    try {
      const payload = { IPAddress: form.IPAddress.trim(), Action: form.Action, Notes: form.Notes.trim() };
      if (editingId) {
        await api.put(`/firewall/${editingId}`, payload);
        toast.success('Firewall rule updated');
      } else {
        await api.post('/firewall', payload);
        toast.success('Firewall rule created');
      }
      modalOpen = false;
      await load();
    } catch (e) {
      toastError(e);
    } finally {
      saving = false;
    }
  }

  async function remove(r: FirewallRule) {
    if (!confirm(`Delete firewall rule for ${r.IPAddress}?`)) return;
    try {
      await api.del(`/firewall/${r.ID}`);
      toast.success('Firewall rule deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Firewall & IP ACLs" description="Allow or block SIP and API traffic by source IP address or CIDR range.">
  <Button on:click={openCreate}>+ Add Rule</Button>
</PageHeader>

<Card padded={false} title="Rules" subtitle="{rules.length} configured">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search IP or notes…" bind:value={search} aria-label="Search rules" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No rules match your search' : 'No firewall rules'}
      message={search ? '' : 'All traffic is currently permitted. Add a rule to allow or block a source address.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first rule</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr><th>Source IP / CIDR</th><th>Action</th><th>Notes</th><th>Created</th><th class="text-right">Actions</th></tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr>
              <td class="mono font-medium text-slate-800">{r.IPAddress}</td>
              <td><Badge tone={isBlock(r.Action) ? 'error' : 'success'}>{isBlock(r.Action) ? 'Block' : 'Allow'}</Badge></td>
              <td class="text-slate-600">{r.Notes || '—'}</td>
              <td class="text-slate-500 text-xs">{formatDateTime(r.CreatedAt)}</td>
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Firewall Rule' : 'Add Firewall Rule'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="fw-ip" class="label">Source IP address or CIDR</label>
      <input id="fw-ip" class="input mono" bind:value={form.IPAddress} placeholder="e.g. 203.0.113.0/24" required />
    </div>
    <div>
      <label for="fw-action" class="label">Action</label>
      <select id="fw-action" class="select" bind:value={form.Action}>
        <option value="accept">Allow</option>
        <option value="drop">Block</option>
      </select>
    </div>
    <div>
      <label for="fw-notes" class="label">Notes</label>
      <input id="fw-notes" class="input" bind:value={form.Notes} placeholder="Optional description" />
    </div>
    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => (modalOpen = false)}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
