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

  interface Site {
    ID: string;
    Name: string;
    SiteCode: string;
    Timezone: string;
    CallerIDNumber: string;
    EmergencyAddress: string;
    IsDefault: boolean;
  }

  let sites: Site[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    SiteCode: '',
    Timezone: 'UTC',
    CallerIDNumber: '',
    EmergencyAddress: '',
    IsDefault: false
  };

  $: filteredSites = sites.filter(s =>
    !search.trim() || `${s.Name} ${s.SiteCode} ${s.CallerIDNumber} ${s.EmergencyAddress}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const res = await api.get<{ sites: Site[] }>('/sites');
      sites = res.sites || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      Name: '',
      SiteCode: String((sites.length + 1) * 10),
      Timezone: 'Europe/London',
      CallerIDNumber: '',
      EmergencyAddress: '',
      IsDefault: sites.length === 0
    };
    modalOpen = true;
  }

  function openEdit(s: Site) {
    editingId = s.ID;
    form = {
      Name: s.Name,
      SiteCode: s.SiteCode,
      Timezone: s.Timezone || 'UTC',
      CallerIDNumber: s.CallerIDNumber || '',
      EmergencyAddress: s.EmergencyAddress || '',
      IsDefault: s.IsDefault
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Site name is required');
    if (!form.SiteCode.trim()) return toast.error('Site code prefix is required');

    saving = true;
    try {
      if (editingId) {
        await api.put(`/sites/${editingId}`, form);
        toast.success('Site location updated');
      } else {
        await api.post('/sites', form);
        toast.success('Site location created');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteSite(id: string, name: string) {
    if (!confirm(`Delete branch office site "${name}"?`)) return;
    try {
      await api.del(`/sites/${id}`);
      toast.success('Site removed');
      sites = sites.filter(s => s.ID !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Sites & Branch Locations"
    subtitle="Configure multi-site company topology, branch dial prefixes, and emergency 999/911 dispatch addresses."
  >
    <Button variant="primary" on:click={openCreate}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Site Location
    </Button>
  </PageHeader>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder="Search sites or emergency addresses..."
        class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
      />
      <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
    <div class="text-xs text-slate-500">
      Configured Sites: <span class="font-bold text-slate-700">{sites.length}</span>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading multi-site locations...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if filteredSites.length === 0}
    <EmptyState
      title="No Sites Configured"
      description="Define office locations to route emergency calls with the correct postal address and organize branch dial plans."
    >
      <Button variant="primary" on:click={openCreate}>Create Primary Site</Button>
    </EmptyState>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each filteredSites as s}
        <Card>
          <div class="flex justify-between items-start mb-3">
            <div>
              <div class="flex items-center space-x-2">
                <h3 class="font-bold text-slate-900 text-sm">{s.Name}</h3>
                {#if s.IsDefault}
                  <span class="px-1.5 py-0.5 bg-sky-100 text-sky-800 rounded text-[10px] font-bold uppercase tracking-wider">
                    HQ
                  </span>
                {/if}
              </div>
              <p class="text-xs text-slate-500 font-mono mt-0.5">Site Prefix: {s.SiteCode}</p>
            </div>
            <div class="flex space-x-2">
              <button
                type="button"
                class="text-sky-600 hover:text-sky-800 text-xs font-medium"
                on:click={() => openEdit(s)}
              >
                Edit
              </button>
              {#if !s.IsDefault}
                <button
                  type="button"
                  class="text-rose-600 hover:text-rose-800 text-xs font-medium"
                  on:click={() => deleteSite(s.ID, s.Name)}
                >
                  Delete
                </button>
              {/if}
            </div>
          </div>

          <div class="space-y-2 text-xs text-slate-600 pt-2 border-t border-slate-100">
            <div class="flex justify-between">
              <span class="text-slate-400">Local Timezone:</span>
              <span class="font-medium text-slate-800">{s.Timezone}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">Default Site CLI:</span>
              <span class="font-mono font-medium text-sky-700">{s.CallerIDNumber || 'Default Trunk'}</span>
            </div>
            <div class="pt-2">
              <span class="text-slate-400 block text-[11px] mb-0.5">Emergency Dispatch Address:</span>
              <p class="text-slate-800 bg-slate-50 p-2 rounded border border-slate-200 text-[11px] leading-relaxed">
                {s.EmergencyAddress || 'Not set (will inherit default organization billing address)'}
              </p>
            </div>
          </div>
        </Card>
      {/each}
    </div>
  {/if}
</div>

<!-- SITE MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? 'Edit Site Location' : 'New Branch Office Site'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="site-name" class="block font-medium text-slate-700 mb-1">Site / Office Name</label>
        <input
          id="site-name"
          type="text"
          bind:value={form.Name}
          placeholder="e.g. Manchester Branch or London HQ"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
      <div>
        <label for="site-code" class="block font-medium text-slate-700 mb-1">Site Code Prefix</label>
        <input
          id="site-code"
          type="text"
          bind:value={form.SiteCode}
          placeholder="e.g. 10 or 20"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono font-bold"
        />
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="site-tz" class="block font-medium text-slate-700 mb-1">Office Timezone</label>
        <select
          id="site-tz"
          bind:value={form.Timezone}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
        >
          <option value="Europe/London">Europe/London (GMT/BST)</option>
          <option value="America/New_York">America/New_York (EST/EDT)</option>
          <option value="America/Chicago">America/Chicago (CST/CDT)</option>
          <option value="America/Los_Angeles">America/Los_Angeles (PST/PDT)</option>
          <option value="Europe/Paris">Europe/Paris (CET)</option>
          <option value="Asia/Tokyo">Asia/Tokyo (JST)</option>
          <option value="UTC">UTC Universal</option>
        </select>
      </div>
      <div>
        <label for="site-cli" class="block font-medium text-slate-700 mb-1">Outbound Emergency CLI</label>
        <input
          id="site-cli"
          type="text"
          bind:value={form.CallerIDNumber}
          placeholder="e.g. +442079460912"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
    </div>

    <div>
      <label for="site-addr" class="block font-medium text-slate-700 mb-1">Emergency Dispatch Address (E911 / 999)</label>
      <textarea
        id="site-addr"
        bind:value={form.EmergencyAddress}
        rows="3"
        placeholder="Full street address, building floor, and postal code for emergency service dispatchers"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 leading-normal"
      ></textarea>
    </div>

    <div class="flex items-center space-x-2 pt-1">
      <input
        type="checkbox"
        id="site_default"
        bind:checked={form.IsDefault}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="site_default" class="text-slate-700 font-medium">Set as primary Headquarters location</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Create Site'}
      </Button>
    </div>
  </form>
</Modal>
