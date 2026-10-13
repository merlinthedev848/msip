<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface ThrottleItem {
    ID: string;
    Name: string;
    Rate: number;
  }

  let rules: ThrottleItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    Rate: 60
  };

  $: filtered = rules.filter(r => 
    !search.trim() || r.Name.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ throttling: ThrottleItem[] }>('/throttling');
      rules = data.throttling || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { Name: '', Rate: 60 };
    modalOpen = true;
  }

  function openEdit(r: ThrottleItem) {
    editingId = r.ID;
    form = {
      Name: r.Name,
      Rate: r.Rate || 60
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Policy name is required');

    saving = true;
    try {
      const payload = {
        Name: form.Name.trim(),
        Rate: Number(form.Rate) || 60
      };

      if (editingId) {
        await api.put(`/throttling/${editingId}`, payload);
        toast.success('Throttling policy updated');
      } else {
        await api.post('/throttling', payload);
        toast.success('Throttling policy created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(r: ThrottleItem) {
    if (!confirm(`Delete rate throttling policy "${r.Name}"?`)) return;
    try {
      await api.del(`/throttling/${r.ID}`);
      toast.success('Policy deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Rate Throttling" description="Configure SIP INVITE and REST API request rate limits.">
  <Button on:click={openCreate}>+ New Policy</Button>
</PageHeader>

<Card padded={false} title="Active Throttling Policies" subtitle="{rules.length} policies">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search policy…" bind:value={search} aria-label="Search policies" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No policies match your search' : 'No rate limit policies configured'}
      message={search ? '' : 'Add a rate limit policy to prevent SIP INVITE flood attacks and API abuse.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first policy</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Policy Name</th>
            <th>Rate Limit</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{r.Name}</td>
              <td class="mono text-xs text-sky-700 font-semibold">{r.Rate} requests / min</td>
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Policy' : 'New Rate Limit Policy'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="throttle-name" class="label">Policy Name</label>
      <input id="throttle-name" class="input" bind:value={form.Name} placeholder="e.g. Inbound SIP INVITE Filter" required />
    </div>

    <div>
      <label for="throttle-rate" class="label">Allowed Request Rate (requests / minute)</label>
      <input id="throttle-rate" type="number" min="1" max="10000" class="input mono" bind:value={form.Rate} required />
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
