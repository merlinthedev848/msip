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

  interface FeatureCode {
    id: string;
    code: string;
    name: string;
    description: string;
    is_active: boolean;
  }

  let codes: FeatureCode[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    code: '',
    name: '',
    description: '',
    is_active: true
  };

  $: filteredCodes = codes.filter(c =>
    !search.trim() || `${c.code} ${c.name} ${c.description}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const res = await api.get<{ feature_codes: FeatureCode[] }>('/feature-codes');
      codes = res.feature_codes || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      code: '',
      name: '',
      description: '',
      is_active: true
    };
    modalOpen = true;
  }

  function openEdit(c: FeatureCode) {
    editingId = c.id;
    form = {
      code: c.code,
      name: c.name,
      description: c.description || '',
      is_active: c.is_active
    };
    modalOpen = true;
  }

  async function toggleStatus(c: FeatureCode) {
    try {
      const updated = !c.is_active;
      await api.put(`/feature-codes/${c.id}`, { ...c, is_active: updated });
      c.is_active = updated;
      codes = [...codes];
      toast.success(`${c.name} ${updated ? 'enabled' : 'disabled'}`);
    } catch (err: any) {
      toastError(err);
    }
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.code.trim()) return toast.error('Dial code is required');
    if (!form.name.trim()) return toast.error('Feature name is required');

    saving = true;
    try {
      if (editingId) {
        await api.put(`/feature-codes/${editingId}`, form);
        toast.success('Feature code updated');
      } else {
        await api.post('/feature-codes', form);
        toast.success('Feature code added');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteCode(id: string, name: string) {
    if (!confirm(`Delete feature code "${name}"?`)) return;
    try {
      await api.del(`/feature-codes/${id}`);
      toast.success('Feature code deleted');
      codes = codes.filter(c => c.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="PBX Feature Codes"
    subtitle="Standard dial codes used from physical desk phones and softphones for voicemail, transfers, and AI tools."
  >
    <Button variant="primary" on:click={openCreate}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Feature Code
    </Button>
  </PageHeader>

  <!-- SEARCH & STATS BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder="Search code or description..."
        class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
      />
      <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
    <div class="text-xs text-slate-500">
      Total Codes: <span class="font-bold text-slate-700">{codes.length}</span>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading feature codes catalog...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if filteredCodes.length === 0}
    <EmptyState
      title="No Feature Codes Found"
      description="Feature codes give phone users access to instant dial shortcuts like *97 for voicemail or *8 for call pickup."
    >
      <Button variant="primary" on:click={openCreate}>Create Feature Code</Button>
    </EmptyState>
  {:else}
    <Card padding="none">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
              <th class="py-3 px-4">Dial Code</th>
              <th class="py-3 px-4">Feature Function</th>
              <th class="py-3 px-4">Description</th>
              <th class="py-3 px-4">Status</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            {#each filteredCodes as c}
              <tr class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3 px-4">
                  <span class="px-2.5 py-1 bg-sky-50 text-sky-700 border border-sky-200 rounded font-mono font-bold text-xs tracking-wide">
                    {c.code}
                  </span>
                </td>
                <td class="py-3 px-4 font-semibold text-slate-900">{c.name}</td>
                <td class="py-3 px-4 text-slate-600 max-w-md">{c.description || '—'}</td>
                <td class="py-3 px-4">
                  <button
                    type="button"
                    class="cursor-pointer"
                    on:click={() => toggleStatus(c)}
                    title="Click to toggle status"
                  >
                    {#if c.is_active}
                      <Badge variant="success">Active</Badge>
                    {:else}
                      <Badge variant="neutral">Disabled</Badge>
                    {/if}
                  </button>
                </td>
                <td class="py-3 px-4 text-right space-x-2">
                  <button
                    type="button"
                    class="text-sky-600 hover:text-sky-800 font-medium"
                    on:click={() => openEdit(c)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    class="text-rose-600 hover:text-rose-800 font-medium"
                    on:click={() => deleteCode(c.id, c.name)}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </Card>
  {/if}
</div>

<!-- FEATURE CODE MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? 'Edit Feature Code' : 'New Feature Code'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="fc-code" class="block font-medium text-slate-700 mb-1">Dial String</label>
        <input
          id="fc-code"
          type="text"
          bind:value={form.code}
          placeholder="e.g. *97 or *8"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono font-bold text-sm"
        />
      </div>
      <div>
        <label for="fc-name" class="block font-medium text-slate-700 mb-1">Feature Name</label>
        <input
          id="fc-name"
          type="text"
          bind:value={form.name}
          placeholder="e.g. Personal Voicemail"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
    </div>

    <div>
      <label for="fc-desc" class="block font-medium text-slate-700 mb-1">Description</label>
      <input
        id="fc-desc"
        type="text"
        bind:value={form.description}
        placeholder="e.g. Dials user's personal mailbox without asking for mailbox number"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="flex items-center space-x-2 pt-2">
      <input
        type="checkbox"
        id="fc_active"
        bind:checked={form.is_active}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="fc_active" class="text-slate-700 font-medium">Feature code is active and dialable</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Create Feature Code'}
      </Button>
    </div>
  </form>
</Modal>
