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

  interface FraudRuleItem {
    ID: string;
    Prefix: string;
    MaxCost: number;
    BlockAction: string;
  }

  let rules: FraudRuleItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Prefix: '',
    MaxCost: 0,
    BlockAction: 'BLOCK'
  };

  $: filtered = rules.filter(r => 
    !search.trim() || r.Prefix.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ fraud: FraudRuleItem[] }>('/fraud');
      rules = data.fraud || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { Prefix: '', MaxCost: 0, BlockAction: 'BLOCK' };
    modalOpen = true;
  }

  function openEdit(r: FraudRuleItem) {
    editingId = r.ID;
    form = {
      Prefix: r.Prefix,
      MaxCost: r.MaxCost || 0,
      BlockAction: r.BlockAction || 'BLOCK'
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Prefix.trim()) return toast.error('Destination prefix is required');

    saving = true;
    try {
      const payload = {
        Prefix: form.Prefix.trim(),
        MaxCost: Number(form.MaxCost) || 0,
        BlockAction: form.BlockAction
      };

      if (editingId) {
        await api.put(`/fraud/${editingId}`, payload);
        toast.success('Fraud rule updated');
      } else {
        await api.post('/fraud', payload);
        toast.success('Fraud rule created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(r: FraudRuleItem) {
    if (!confirm(`Delete anti-fraud rule for prefix ${r.Prefix}?`)) return;
    try {
      await api.del(`/fraud/${r.ID}`);
      toast.success('Fraud rule deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Anti-Fraud Protection" description="Block high-rate toll destinations and enforce international call cost limits.">
  <Button on:click={openCreate}>+ New Rule</Button>
</PageHeader>

<Card padded={false} title="Protected Prefixes" subtitle="{rules.length} rules active">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search prefix…" bind:value={search} aria-label="Search prefixes" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No fraud rules match your search' : 'No fraud rules configured'}
      message={search ? '' : 'Add destination prefixes (e.g. international country codes) to block unauthorized toll fraud.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first rule</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Country / Dial Prefix</th>
            <th>Max Rate / Minute</th>
            <th>Enforcement Action</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr>
              <td class="mono font-bold text-slate-800 text-sm">+{r.Prefix}</td>
              <td class="mono text-xs text-slate-700">${r.MaxCost.toFixed(2)}</td>
              <td><Badge tone={r.BlockAction === 'BLOCK' ? 'error' : 'warning'}>{r.BlockAction}</Badge></td>
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Fraud Rule' : 'New Anti-Fraud Rule'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="fraud-prefix" class="label">Destination Dial Prefix (Country Code)</label>
      <input id="fraud-prefix" class="input mono" bind:value={form.Prefix} placeholder="e.g. 234 or 93" required />
    </div>

    <div>
      <label for="fraud-cost" class="label">Maximum Allowable Rate ($/min)</label>
      <input id="fraud-cost" type="number" step="0.01" class="input mono" bind:value={form.MaxCost} required />
    </div>

    <div>
      <label for="fraud-action" class="label">Enforcement Action</label>
      <select id="fraud-action" class="select" bind:value={form.BlockAction}>
        <option value="BLOCK">Block Immediately (403 Forbidden)</option>
        <option value="PIN_CHALLENGE">Require Supervisor PIN</option>
      </select>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
