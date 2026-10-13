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

  interface KeyConfig {
    key: number;
    type: string;
    value: string;
    label: string;
  }

  interface ButtonProfile {
    ID: string;
    Name: string;
    DeviceModel: string;
    Description: string;
    Keys: string;
  }

  let profiles: ButtonProfile[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    DeviceModel: 'Standard IP Handset (10 Keys)',
    Description: '',
    keysList: [] as KeyConfig[]
  };

  $: filteredProfiles = profiles.filter(p =>
    !search.trim() || `${p.Name} ${p.DeviceModel} ${p.Description}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const res = await api.get<{ button_profiles: ButtonProfile[] }>('/button-profiles');
      profiles = res.button_profiles || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function parseKeys(keysJson: string): KeyConfig[] {
    try {
      return JSON.parse(keysJson || '[]');
    } catch {
      return [];
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      Name: '',
      DeviceModel: 'Standard IP Handset (10 Keys)',
      Description: '',
      keysList: [
        { key: 1, type: 'Line', value: '1', label: 'Line 1' }
      ]
    };
    modalOpen = true;
  }

  function openEdit(p: ButtonProfile) {
    editingId = p.ID;
    form = {
      Name: p.Name,
      DeviceModel: p.DeviceModel,
      Description: p.Description || '',
      keysList: parseKeys(p.Keys)
    };
    if (form.keysList.length === 0) {
      form.keysList = [
        { key: 1, type: 'Line', value: '1', label: 'Line 1' }
      ];
    }
    modalOpen = true;
  }

  function addKey() {
    const nextKey = form.keysList.length + 1;
    form.keysList = [
      ...form.keysList,
      { key: nextKey, type: 'BLF', value: '', label: `Key ${nextKey}` }
    ];
  }

  function removeKey(idx: number) {
    form.keysList = form.keysList.filter((_, i) => i !== idx);
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Profile name is required');

    saving = true;
    try {
      const payload = {
        Name: form.Name.trim(),
        DeviceModel: form.DeviceModel.trim(),
        Description: form.Description.trim(),
        Keys: JSON.stringify(form.keysList)
      };

      if (editingId) {
        await api.put(`/button-profiles/${editingId}`, payload);
        toast.success('Button profile updated');
      } else {
        await api.post('/button-profiles', payload);
        toast.success('Button profile created');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteProfile(id: string, name: string) {
    if (!confirm(`Delete button profile "${name}"?`)) return;
    try {
      await api.del(`/button-profiles/${id}`);
      toast.success('Profile deleted');
      profiles = profiles.filter(p => p.ID !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Button Profiles (BLF & DSS Keys)"
    subtitle="Configure programmable line keys, busy lamp field (BLF) status lights, and speed dials for phone screens."
  >
    <Button variant="primary" on:click={openCreate}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Button Profile
    </Button>
  </PageHeader>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder="Filter button profiles..."
        class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
      />
      <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
    <div class="text-xs text-slate-500">
      Total Templates: <span class="font-bold text-slate-700">{profiles.length}</span>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading programmable button profiles...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if filteredProfiles.length === 0}
    <EmptyState
      title="No Button Profiles Defined"
      description="Create visual line key templates for phone models to control what each physical button dials."
    >
      <Button variant="primary" on:click={openCreate}>Create First Profile</Button>
    </EmptyState>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each filteredProfiles as p}
        {@const keys = parseKeys(p.Keys)}
        <Card>
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-bold text-slate-900 text-sm">{p.Name}</h3>
              <p class="text-xs text-sky-700 font-medium">{p.DeviceModel}</p>
            </div>
            <div class="flex space-x-2">
              <button
                type="button"
                class="text-sky-600 hover:text-sky-800 text-xs font-medium"
                on:click={() => openEdit(p)}
              >
                Edit
              </button>
              <button
                type="button"
                class="text-rose-600 hover:text-rose-800 text-xs font-medium"
                on:click={() => deleteProfile(p.ID, p.Name)}
              >
                Delete
              </button>
            </div>
          </div>

          <p class="text-xs text-slate-600 mb-4 line-clamp-2">{p.Description || 'Standard key layout for extension handsets.'}</p>

          <!-- Handset LCD Keys Preview -->
          <div class="bg-slate-900 text-white p-3 rounded-lg border border-slate-800 font-mono text-[11px] space-y-1.5 shadow-inner">
            <div class="text-[10px] text-slate-400 border-b border-slate-800 pb-1 mb-1.5 flex justify-between">
              <span>LCD SCREEN PREVIEW</span>
              <span>{keys.length} KEYS</span>
            </div>
            {#each keys.slice(0, 6) as k}
              <div class="flex items-center justify-between py-0.5 px-1.5 bg-slate-800/80 rounded border border-slate-700/50">
                <span class="flex items-center space-x-1.5">
                  <span class="w-1.5 h-1.5 rounded-full {k.type === 'BLF' ? 'bg-emerald-400 animate-pulse' : 'bg-sky-400'}"></span>
                  <span class="truncate max-w-[110px]">{k.label || `Key ${k.key}`}</span>
                </span>
                <span class="text-slate-400 text-[10px]">{k.type} ({k.value})</span>
              </div>
            {/each}
            {#if keys.length > 6}
              <div class="text-[10px] text-slate-500 text-center pt-1">+ {keys.length - 6} more keys configured</div>
            {/if}
          </div>
        </Card>
      {/each}
    </div>
  {/if}
</div>

<!-- BUTTON PROFILE MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? 'Edit Button Profile' : 'New Button Profile'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="bp-name" class="block font-medium text-slate-700 mb-1">Profile Name</label>
        <input
          id="bp-name"
          type="text"
          bind:value={form.Name}
          placeholder="e.g. Sales Desk Phone (10 Keys)"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
      <div>
        <label for="bp-model" class="block font-medium text-slate-700 mb-1">Target Device Model</label>
        <input
          id="bp-model"
          type="text"
          bind:value={form.DeviceModel}
          placeholder="e.g. Standard IP Handset"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
    </div>

    <div>
      <label for="bp-desc" class="block font-medium text-slate-700 mb-1">Description (Optional)</label>
      <input
        id="bp-desc"
        type="text"
        bind:value={form.Description}
        placeholder="e.g. Assigned to Receptionists and Frontline Support"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <!-- Programmable Keys List -->
    <div class="border-t border-slate-200 pt-3">
      <div class="flex justify-between items-center mb-2">
        <span class="font-bold text-slate-800">Programmable LCD Screen Keys ({form.keysList.length})</span>
        <Button variant="secondary" type="button" on:click={addKey}>
          + Add Key
        </Button>
      </div>

      <div class="space-y-2 max-h-60 overflow-y-auto pr-1">
        {#each form.keysList as key, idx}
          <div class="p-2.5 bg-slate-50 border border-slate-200 rounded-md flex items-center space-x-2">
            <span class="font-mono text-slate-500 text-xs w-6 text-center">#{idx + 1}</span>
            <input
              type="text"
              bind:value={key.label}
              placeholder="Display Label"
              required
              class="w-32 px-2 py-1 bg-white border border-slate-200 rounded text-xs focus:outline-none focus:border-sky-500"
            />
            <select
              bind:value={key.type}
              class="w-28 px-2 py-1 bg-white border border-slate-200 rounded text-xs focus:outline-none focus:border-sky-500"
            >
              <option value="Line">Line Key</option>
              <option value="BLF">BLF (Busy Lamp)</option>
              <option value="SpeedDial">Speed Dial</option>
              <option value="Intercom">Intercom</option>
              <option value="DND">DND</option>
              <option value="Park">Call Park</option>
            </select>
            <input
              type="text"
              bind:value={key.value}
              placeholder="Value / Ext"
              required
              class="flex-1 px-2 py-1 bg-white border border-slate-200 rounded text-xs font-mono focus:outline-none focus:border-sky-500"
            />
            <button
              type="button"
              class="text-rose-500 hover:text-rose-700 p-1"
              on:click={() => removeKey(idx)}
              title="Remove key"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        {/each}
      </div>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Create Profile'}
      </Button>
    </div>
  </form>
</Modal>
