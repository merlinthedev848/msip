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

  interface PagingGroup {
    id: string;
    name: string;
    number: string;
    duplex_mode: string;
    members: string[];
  }

  interface PickupGroup {
    id: string;
    name: string;
    code: string;
    members: string[];
  }

  interface ExtensionOption {
    ExtensionNumber: string;
    DisplayName?: string;
  }

  let activeTab: 'paging' | 'pickup' = 'paging';

  let pagingGroups: PagingGroup[] = [];
  let pickupGroups: PickupGroup[] = [];
  let availableExts: ExtensionOption[] = [];

  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // Paging Modal
  let pageModalOpen = false;
  let editingPageId: string | null = null;
  let pageForm = {
    name: '',
    number: '630',
    duplex_mode: 'one-way',
    members_text: ''
  };

  // Pickup Modal
  let pickupModalOpen = false;
  let editingPickupId: string | null = null;
  let pickupForm = {
    name: '',
    code: '*8',
    members_text: ''
  };

  $: filteredPaging = pagingGroups.filter(p =>
    !search.trim() || `${p.name} ${p.number}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  $: filteredPickup = pickupGroups.filter(p =>
    !search.trim() || `${p.name} ${p.code}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [pRes, pickRes, extRes] = await Promise.all([
        api.get<{ groups: PagingGroup[] }>('/paging-groups'),
        api.get<{ groups: PickupGroup[] }>('/pickup-groups'),
        api.get<{ extensions: ExtensionOption[] }>('/extensions')
      ]);
      pagingGroups = (pRes.groups || []).map((p: any) => ({
        ...p,
        id: p.id || p.ID,
        members: Array.isArray(p.members) ? p.members : (typeof p.members === 'string' ? p.members.split(',').map((s: string) => s.trim()).filter(Boolean) : [])
      }));
      pickupGroups = (pickRes.groups || []).map((p: any) => ({
        ...p,
        id: p.id || p.ID,
        members: Array.isArray(p.members) ? p.members : (typeof p.members === 'string' ? p.members.split(',').map((s: string) => s.trim()).filter(Boolean) : [])
      }));
      availableExts = extRes.extensions || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreatePaging() {
    editingPageId = null;
    pageForm = {
      name: '',
      number: '630',
      duplex_mode: 'one-way',
      members_text: availableExts.map(e => e.ExtensionNumber).join(', ')
    };
    pageModalOpen = true;
  }

  function openEditPaging(p: PagingGroup) {
    editingPageId = p.id;
    pageForm = {
      name: p.name,
      number: p.number,
      duplex_mode: p.duplex_mode || 'one-way',
      members_text: (p.members || []).join(', ')
    };
    pageModalOpen = true;
  }

  async function savePaging(e: Event) {
    e.preventDefault();
    if (!pageForm.name.trim()) return toast.error('Paging group name is required');
    if (!pageForm.number.trim()) return toast.error('Dial number is required');

    const members = pageForm.members_text.split(',').map(m => m.trim()).filter(Boolean);
    if (members.length === 0) return toast.error('Add at least one extension to the group');

    saving = true;
    try {
      const payload = {
        name: pageForm.name,
        number: pageForm.number,
        duplex_mode: pageForm.duplex_mode,
        members
      };
      if (editingPageId) {
        await api.put(`/paging-groups/${editingPageId}`, payload);
        toast.success('Paging group updated');
      } else {
        await api.post('/paging-groups', payload);
        toast.success('Paging group created');
      }
      pageModalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deletePaging(id: string, name: string) {
    if (!confirm(`Delete paging group "${name}"?`)) return;
    try {
      await api.del(`/paging-groups/${id}`);
      toast.success('Paging group removed');
      pagingGroups = pagingGroups.filter(p => p.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  // PICKUP
  function openCreatePickup() {
    editingPickupId = null;
    pickupForm = {
      name: '',
      code: '*8',
      members_text: availableExts.slice(0, 3).map(e => e.ExtensionNumber).join(', ')
    };
    pickupModalOpen = true;
  }

  function openEditPickup(p: PickupGroup) {
    editingPickupId = p.id;
    pickupForm = {
      name: p.name,
      code: p.code,
      members_text: (p.members || []).join(', ')
    };
    pickupModalOpen = true;
  }

  async function savePickup(e: Event) {
    e.preventDefault();
    if (!pickupForm.name.trim()) return toast.error('Pickup group name is required');
    if (!pickupForm.code.trim()) return toast.error('Pickup code is required');

    const members = pickupForm.members_text.split(',').map(m => m.trim()).filter(Boolean);
    if (members.length === 0) return toast.error('Add at least one extension to the group');

    saving = true;
    try {
      const payload = {
        name: pickupForm.name,
        code: pickupForm.code,
        members
      };
      if (editingPickupId) {
        await api.put(`/pickup-groups/${editingPickupId}`, payload);
        toast.success('Pickup group updated');
      } else {
        await api.post('/pickup-groups', payload);
        toast.success('Pickup group created');
      }
      pickupModalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deletePickup(id: string, name: string) {
    if (!confirm(`Delete pickup group "${name}"?`)) return;
    try {
      await api.del(`/pickup-groups/${id}`);
      toast.success('Pickup group removed');
      pickupGroups = pickupGroups.filter(p => p.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Paging, Intercom & Call Pickup"
    subtitle="Broadcast emergency announcements across speakerphones and intercept ringing desk phones."
  >
    {#if activeTab === 'paging'}
      <Button variant="primary" on:click={openCreatePaging}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Paging Group
      </Button>
    {:else}
      <Button variant="primary" on:click={openCreatePickup}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Pickup Group
      </Button>
    {/if}
  </PageHeader>

  <!-- TAB SWITCHER -->
  <div class="flex border-b border-slate-200 space-x-6 text-sm font-medium">
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'paging' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'paging'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
      </svg>
      <span>Paging & Intercom Groups ({pagingGroups.length})</span>
    </button>
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'pickup' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'pickup'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
      </svg>
      <span>Call Pickup Groups (*8) ({pickupGroups.length})</span>
    </button>
  </div>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder={activeTab === 'paging' ? 'Search paging groups...' : 'Search pickup groups...'}
        class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
      />
      <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading paging and pickup groups...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if activeTab === 'paging'}
    {#if filteredPaging.length === 0}
      <EmptyState
        title="No Paging Groups Created"
        description="Paging allows auto-answering desk phone speakerphones for warehouse alerts or public announcements."
      >
        <Button variant="primary" on:click={openCreatePaging}>Create Paging Zone</Button>
      </EmptyState>
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Group Name</th>
                <th class="py-3 px-4">Dial Code</th>
                <th class="py-3 px-4">Duplex Mode</th>
                <th class="py-3 px-4">Target Extensions</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredPaging as p}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-semibold text-slate-900">{p.name}</td>
                  <td class="py-3 px-4 font-mono font-medium text-sky-700">{p.number}</td>
                  <td class="py-3 px-4 capitalize">
                    {#if p.duplex_mode === 'two-way'}
                      <span class="px-2 py-0.5 bg-purple-50 border border-purple-200 text-purple-700 rounded text-[11px] font-medium">
                        Two-Way Intercom
                      </span>
                    {:else}
                      <span class="px-2 py-0.5 bg-blue-50 border border-blue-200 text-blue-700 rounded text-[11px] font-medium">
                        One-Way Broadcast
                      </span>
                    {/if}
                  </td>
                  <td class="py-3 px-4">
                    <div class="flex flex-wrap gap-1">
                      {#each (p.members || []) as m}
                        <span class="px-1.5 py-0.5 bg-slate-100 rounded text-slate-700 font-mono text-[11px]">
                          {m}
                        </span>
                      {/each}
                    </div>
                  </td>
                  <td class="py-3 px-4 text-right space-x-2">
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 font-medium"
                      on:click={() => openEditPaging(p)}
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="text-rose-600 hover:text-rose-800 font-medium"
                      on:click={() => deletePaging(p.id, p.name)}
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
  {:else}
    <!-- PICKUP TAB -->
    <div class="p-4 bg-sky-50 border border-sky-200 rounded-lg text-xs text-sky-900 leading-relaxed flex items-start space-x-3">
      <svg class="w-5 h-5 text-sky-600 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <div>
        <span class="font-bold">What is Call Pickup?</span>
        When an extension in your department is ringing while the owner is away from their desk, picking up your phone and dialing <span class="font-mono font-bold bg-sky-100 px-1 rounded">*8</span> intercepts the incoming call and answers it on your own handset immediately.
      </div>
    </div>

    {#if filteredPickup.length === 0}
      <EmptyState
        title="No Pickup Groups Defined"
        description="Create department pickup groups so team members can answer each other's ringing phones with *8."
      >
        <Button variant="primary" on:click={openCreatePickup}>Create Pickup Group</Button>
      </EmptyState>
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Department / Group Name</th>
                <th class="py-3 px-4">Pickup Code</th>
                <th class="py-3 px-4">Group Extension Members</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredPickup as p}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-semibold text-slate-900">{p.name}</td>
                  <td class="py-3 px-4 font-mono font-bold text-sky-700">{p.code}</td>
                  <td class="py-3 px-4">
                    <div class="flex flex-wrap gap-1">
                      {#each (p.members || []) as m}
                        <span class="px-2 py-0.5 bg-slate-100 rounded text-slate-700 font-mono text-[11px]">
                          Ext {m}
                        </span>
                      {/each}
                    </div>
                  </td>
                  <td class="py-3 px-4 text-right space-x-2">
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 font-medium"
                      on:click={() => openEditPickup(p)}
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="text-rose-600 hover:text-rose-800 font-medium"
                      on:click={() => deletePickup(p.id, p.name)}
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
  {/if}
</div>

<!-- PAGING MODAL -->
<Modal
  open={pageModalOpen}
  title={editingPageId ? 'Edit Paging Group' : 'Create Paging Group'}
  on:close={() => pageModalOpen = false}
>
  <form on:submit={savePaging} class="space-y-4 text-xs">
    <div>
      <label for="p-name" class="block font-medium text-slate-700 mb-1">Group Name</label>
      <input
        id="p-name"
        type="text"
        bind:value={pageForm.name}
        placeholder="e.g. Warehouse Speaker Alert"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="p-num" class="block font-medium text-slate-700 mb-1">Dial Number</label>
        <input
          id="p-num"
          type="text"
          bind:value={pageForm.number}
          placeholder="e.g. 630"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
      <div>
        <label for="p-duplex" class="block font-medium text-slate-700 mb-1">Mode</label>
        <select
          id="p-duplex"
          bind:value={pageForm.duplex_mode}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
        >
          <option value="one-way">One-Way Broadcast (Muted Listeners)</option>
          <option value="two-way">Two-Way Intercom (Hands-free Talkback)</option>
        </select>
      </div>
    </div>

    <div>
      <label for="p-members" class="block font-medium text-slate-700 mb-1">Member Extensions (Comma Separated)</label>
      <input
        id="p-members"
        type="text"
        bind:value={pageForm.members_text}
        placeholder="e.g. 1001, 1002, 1003"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
      />
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => pageModalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingPageId ? 'Save Changes' : 'Create Group'}
      </Button>
    </div>
  </form>
</Modal>

<!-- PICKUP MODAL -->
<Modal
  open={pickupModalOpen}
  title={editingPickupId ? 'Edit Pickup Group' : 'Create Pickup Group'}
  on:close={() => pickupModalOpen = false}
>
  <form on:submit={savePickup} class="space-y-4 text-xs">
    <div>
      <label for="pick-name" class="block font-medium text-slate-700 mb-1">Group Name</label>
      <input
        id="pick-name"
        type="text"
        bind:value={pickupForm.name}
        placeholder="e.g. Sales Team Pickup"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div>
      <label for="pick-code" class="block font-medium text-slate-700 mb-1">Feature Dial Code</label>
      <input
        id="pick-code"
        type="text"
        bind:value={pickupForm.code}
        placeholder="e.g. *8"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
      />
    </div>

    <div>
      <label for="pick-members" class="block font-medium text-slate-700 mb-1">Member Extensions (Comma Separated)</label>
      <input
        id="pick-members"
        type="text"
        bind:value={pickupForm.members_text}
        placeholder="e.g. 1001, 1002, 1003"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
      />
      <span class="text-[11px] text-slate-400 mt-0.5 block">Any phone in this list can dial *8 to grab a call ringing on another phone in the list.</span>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => pickupModalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingPickupId ? 'Save Changes' : 'Create Group'}
      </Button>
    </div>
  </form>
</Modal>
