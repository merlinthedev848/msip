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

  interface CallDivert {
    ID: string;
    ExtensionNumber: string;
    AlwaysEnabled: boolean;
    AlwaysTarget: string;
    BusyEnabled: boolean;
    BusyTarget: string;
    NoAnswerEnabled: boolean;
    NoAnswerTarget: string;
    NoAnswerTimeout: number;
    OfflineEnabled: boolean;
    OfflineTarget: string;
  }

  interface ExtensionOption {
    ExtensionNumber: string;
    DisplayName?: string;
  }

  let diverts: CallDivert[] = [];
  let availableExts: ExtensionOption[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    ExtensionNumber: '',
    AlwaysEnabled: false,
    AlwaysTarget: '',
    BusyEnabled: false,
    BusyTarget: '',
    NoAnswerEnabled: true,
    NoAnswerTarget: '',
    NoAnswerTimeout: 20,
    OfflineEnabled: false,
    OfflineTarget: ''
  };

  $: filteredDiverts = diverts.filter(d =>
    !search.trim() || d.ExtensionNumber.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [divRes, extRes] = await Promise.all([
        api.get<{ diverts: CallDivert[] }>('/diverts'),
        api.get<{ extensions: ExtensionOption[] }>('/extensions')
      ]);
      diverts = divRes.diverts || [];
      availableExts = extRes.extensions || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    const existingExts = new Set(diverts.map(d => d.ExtensionNumber));
    const nextUnsetExt = availableExts.find(e => !existingExts.has(e.ExtensionNumber));

    form = {
      ExtensionNumber: nextUnsetExt ? nextUnsetExt.ExtensionNumber : availableExts[0]?.ExtensionNumber || '1001',
      AlwaysEnabled: false,
      AlwaysTarget: '',
      BusyEnabled: false,
      BusyTarget: '',
      NoAnswerEnabled: true,
      NoAnswerTarget: '1001',
      NoAnswerTimeout: 20,
      OfflineEnabled: false,
      OfflineTarget: ''
    };
    modalOpen = true;
  }

  function openEdit(d: CallDivert) {
    editingId = d.ID;
    form = {
      ExtensionNumber: d.ExtensionNumber,
      AlwaysEnabled: d.AlwaysEnabled,
      AlwaysTarget: d.AlwaysTarget || '',
      BusyEnabled: d.BusyEnabled,
      BusyTarget: d.BusyTarget || '',
      NoAnswerEnabled: d.NoAnswerEnabled,
      NoAnswerTarget: d.NoAnswerTarget || '',
      NoAnswerTimeout: d.NoAnswerTimeout || 20,
      OfflineEnabled: d.OfflineEnabled,
      OfflineTarget: d.OfflineTarget || ''
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.ExtensionNumber.trim()) return toast.error('Extension number is required');

    saving = true;
    try {
      if (editingId) {
        await api.put(`/diverts/${editingId}`, form);
        toast.success('Call divert policy updated');
      } else {
        await api.post('/diverts', form);
        toast.success('Call divert policy created');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteDivert(id: string, ext: string) {
    if (!confirm(`Reset call diverts for extension ${ext}?`)) return;
    try {
      await api.del(`/diverts/${id}`);
      toast.success('Divert policy reset');
      diverts = diverts.filter(d => d.ID !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Call Diverts & Forwarding"
    subtitle="Configure unconditional, busy, no-answer timeout, and handset offline forwarding rules."
  >
    <Button variant="primary" on:click={openCreate}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Divert Rule
    </Button>
  </PageHeader>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder="Filter by extension..."
        class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
      />
      <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
    <div class="text-xs text-slate-500">
      Configured Diverts: <span class="font-bold text-slate-700">{diverts.length}</span>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading extension call divert rules...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if filteredDiverts.length === 0}
    <EmptyState
      title="No Call Diverts Configured"
      description="Create call divert rules to automatically forward calls when staff are away from desk, on another call, or after hours."
    >
      <Button variant="primary" on:click={openCreate}>Configure First Extension</Button>
    </EmptyState>
  {:else}
    <Card padding="none">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
              <th class="py-3 px-4">Extension</th>
              <th class="py-3 px-4">Divert Always (Unconditional)</th>
              <th class="py-3 px-4">Divert on Busy</th>
              <th class="py-3 px-4">Divert on No Answer</th>
              <th class="py-3 px-4">Divert on Offline</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            {#each filteredDiverts as d}
              <tr class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3 px-4">
                  <span class="px-2.5 py-1 bg-sky-50 text-sky-700 border border-sky-200 rounded font-mono font-bold text-xs">
                    Ext {d.ExtensionNumber}
                  </span>
                </td>
                <td class="py-3 px-4">
                  {#if d.AlwaysEnabled}
                    <span class="inline-flex items-center text-rose-700 font-medium">
                      <span class="w-2 h-2 rounded-full bg-rose-500 mr-1.5"></span>
                      To: {d.AlwaysTarget}
                    </span>
                  {:else}
                    <span class="text-slate-400">Disabled</span>
                  {/if}
                </td>
                <td class="py-3 px-4">
                  {#if d.BusyEnabled}
                    <span class="inline-flex items-center text-amber-700 font-medium">
                      <span class="w-2 h-2 rounded-full bg-amber-500 mr-1.5"></span>
                      To: {d.BusyTarget}
                    </span>
                  {:else}
                    <span class="text-slate-400">Disabled</span>
                  {/if}
                </td>
                <td class="py-3 px-4">
                  {#if d.NoAnswerEnabled}
                    <span class="inline-flex items-center text-emerald-700 font-medium">
                      <span class="w-2 h-2 rounded-full bg-emerald-500 mr-1.5"></span>
                      To: {d.NoAnswerTarget} ({d.NoAnswerTimeout}s)
                    </span>
                  {:else}
                    <span class="text-slate-400">Disabled</span>
                  {/if}
                </td>
                <td class="py-3 px-4">
                  {#if d.OfflineEnabled}
                    <span class="inline-flex items-center text-purple-700 font-medium">
                      <span class="w-2 h-2 rounded-full bg-purple-500 mr-1.5"></span>
                      To: {d.OfflineTarget}
                    </span>
                  {:else}
                    <span class="text-slate-400">Disabled</span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-right space-x-2">
                  <button
                    type="button"
                    class="text-sky-600 hover:text-sky-800 font-medium"
                    on:click={() => openEdit(d)}
                  >
                    Configure
                  </button>
                  <button
                    type="button"
                    class="text-rose-600 hover:text-rose-800 font-medium"
                    on:click={() => deleteDivert(d.ID, d.ExtensionNumber)}
                  >
                    Reset
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

<!-- DIVERT MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? `Configure Diverts: Ext ${form.ExtensionNumber}` : 'New Extension Divert Rule'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <div>
      <label for="div-ext" class="block font-medium text-slate-700 mb-1">Target Extension</label>
      <select
        id="div-ext"
        bind:value={form.ExtensionNumber}
        disabled={!!editingId}
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white disabled:bg-slate-100"
      >
        {#each availableExts as ext}
          <option value={ext.ExtensionNumber}>Extension {ext.ExtensionNumber} ({ext.DisplayName || 'User'})</option>
        {/each}
      </select>
    </div>

    <!-- 1. Always Divert -->
    <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg space-y-2">
      <div class="flex items-center justify-between">
        <label for="div_always_en" class="font-bold text-slate-800 cursor-pointer">Divert Always (Unconditional Forward)</label>
        <input
          type="checkbox"
          id="div_always_en"
          bind:checked={form.AlwaysEnabled}
          class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
        />
      </div>
      {#if form.AlwaysEnabled}
        <input
          type="text"
          bind:value={form.AlwaysTarget}
          placeholder="Destination number or extension (e.g. 1002 or +15550199)"
          required={form.AlwaysEnabled}
          class="w-full px-2.5 py-1.5 bg-white border border-slate-200 rounded text-xs font-mono focus:outline-none focus:border-sky-500"
        />
      {/if}
    </div>

    <!-- 2. Busy Divert -->
    <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg space-y-2">
      <div class="flex items-center justify-between">
        <label for="div_busy_en" class="font-bold text-slate-800 cursor-pointer">Divert on Busy (Engaged on another call)</label>
        <input
          type="checkbox"
          id="div_busy_en"
          bind:checked={form.BusyEnabled}
          class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
        />
      </div>
      {#if form.BusyEnabled}
        <input
          type="text"
          bind:value={form.BusyTarget}
          placeholder="Destination number or colleague extension"
          required={form.BusyEnabled}
          class="w-full px-2.5 py-1.5 bg-white border border-slate-200 rounded text-xs font-mono focus:outline-none focus:border-sky-500"
        />
      {/if}
    </div>

    <!-- 3. No Answer Divert -->
    <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg space-y-2">
      <div class="flex items-center justify-between">
        <label for="div_noans_en" class="font-bold text-slate-800 cursor-pointer">Divert on No Answer (Ring Timeout)</label>
        <input
          type="checkbox"
          id="div_noans_en"
          bind:checked={form.NoAnswerEnabled}
          class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
        />
      </div>
      {#if form.NoAnswerEnabled}
        <div class="grid grid-cols-2 gap-2">
          <input
            type="text"
            bind:value={form.NoAnswerTarget}
            placeholder="Target voicemail or number"
            required={form.NoAnswerEnabled}
            class="w-full px-2.5 py-1.5 bg-white border border-slate-200 rounded text-xs font-mono focus:outline-none focus:border-sky-500"
          />
          <div class="flex items-center space-x-1.5">
            <input
              type="number"
              min="5"
              max="120"
              bind:value={form.NoAnswerTimeout}
              class="w-20 px-2 py-1.5 bg-white border border-slate-200 rounded text-xs font-mono focus:outline-none focus:border-sky-500"
            />
            <span class="text-slate-500">sec timeout</span>
          </div>
        </div>
      {/if}
    </div>

    <!-- 4. Offline Divert -->
    <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg space-y-2">
      <div class="flex items-center justify-between">
        <label for="div_off_en" class="font-bold text-slate-800 cursor-pointer">Divert on Offline (Unregistered Handset Failover)</label>
        <input
          type="checkbox"
          id="div_off_en"
          bind:checked={form.OfflineEnabled}
          class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
        />
      </div>
      {#if form.OfflineEnabled}
        <input
          type="text"
          bind:value={form.OfflineTarget}
          placeholder="Emergency cell phone or backup branch extension"
          required={form.OfflineEnabled}
          class="w-full px-2.5 py-1.5 bg-white border border-slate-200 rounded text-xs font-mono focus:outline-none focus:border-sky-500"
        />
      {/if}
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Create Divert Rule'}
      </Button>
    </div>
  </form>
</Modal>
