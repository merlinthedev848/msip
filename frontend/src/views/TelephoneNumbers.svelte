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

  interface TelephoneNumber {
    ID: string;
    Number: string;
    Description: string;
    DestinationType: string;
    DestinationTarget: string;
    OutboundCLI: boolean;
    Provider: string;
  }

  interface ExtensionOption {
    ExtensionNumber: string;
    DisplayName?: string;
  }

  interface IVROption {
    Name: string;
    Extension?: string;
  }

  let numbers: TelephoneNumber[] = [];
  let availableExts: ExtensionOption[] = [];
  let availableIvrs: IVROption[] = [];

  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Number: '',
    Description: '',
    DestinationType: 'extension',
    DestinationTarget: '',
    OutboundCLI: true,
    Provider: 'Default SIP Trunk'
  };

  $: filteredNumbers = numbers.filter(n =>
    !search.trim() || `${n.Number} ${n.Description} ${n.DestinationTarget} ${n.Provider}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [numRes, extRes, ivrRes] = await Promise.all([
        api.get<{ telephone_numbers: TelephoneNumber[] }>('/telephone-numbers'),
        api.get<{ extensions: ExtensionOption[] }>('/extensions'),
        api.get<{ ivrs: IVROption[] }>('/ivrs')
      ]);
      numbers = numRes.telephone_numbers || [];
      availableExts = extRes.extensions || [];
      availableIvrs = ivrRes.ivrs || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      Number: '',
      Description: '',
      DestinationType: 'extension',
      DestinationTarget: availableExts[0]?.ExtensionNumber || '',
      OutboundCLI: true,
      Provider: ''
    };
    modalOpen = true;
  }

  function openEdit(n: TelephoneNumber) {
    editingId = n.ID;
    form = {
      Number: n.Number,
      Description: n.Description || '',
      DestinationType: n.DestinationType || 'extension',
      DestinationTarget: n.DestinationTarget,
      OutboundCLI: n.OutboundCLI,
      Provider: n.Provider || 'Default Carrier Gateway'
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Number.trim()) return toast.error('Telephone number is required');
    if (!form.DestinationTarget.trim()) return toast.error('Destination target is required');

    saving = true;
    try {
      if (editingId) {
        await api.put(`/telephone-numbers/${editingId}`, form);
        toast.success('Telephone number updated');
      } else {
        await api.post('/telephone-numbers', form);
        toast.success('Telephone number added');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteNumber(id: string, num: string) {
    if (!confirm(`Remove telephone number "${num}" from PBX inventory?`)) return;
    try {
      await api.del(`/telephone-numbers/${id}`);
      toast.success('Number removed');
      numbers = numbers.filter(n => n.ID !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Telephone Numbers (DIDs)"
    subtitle="Manage carrier phone numbers, inbound destination routing, and outbound Caller ID (CLI) presentation."
  >
    <Button variant="primary" on:click={openCreate}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Telephone Number
    </Button>
  </PageHeader>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder="Search number, description..."
        class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
      />
      <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
    <div class="text-xs text-slate-500">
      Assigned DIDs: <span class="font-bold text-slate-700">{numbers.length}</span>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading telephone numbers inventory...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if filteredNumbers.length === 0}
    <EmptyState
      title="No Telephone Numbers Allocated"
      description="Add external phone lines and DIDs provided by your SIP trunk carriers."
    >
      <Button variant="primary" on:click={openCreate}>Add First DID</Button>
    </EmptyState>
  {:else}
    <Card padding="none">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
              <th class="py-3 px-4">Telephone Number</th>
              <th class="py-3 px-4">Description / Label</th>
              <th class="py-3 px-4">Inbound Route Target</th>
              <th class="py-3 px-4">Outbound CLI Presentation</th>
              <th class="py-3 px-4">Carrier Gateway</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            {#each filteredNumbers as n}
              <tr class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3 px-4 font-mono font-bold text-sky-700 text-sm">
                  {n.Number}
                </td>
                <td class="py-3 px-4 font-medium text-slate-900">{n.Description || 'Main Company Number'}</td>
                <td class="py-3 px-4">
                  <span class="inline-flex items-center text-slate-800 font-medium">
                    <span class="w-2 h-2 rounded-full bg-emerald-500 mr-1.5"></span>
                    <span class="capitalize text-slate-500 mr-1">{n.DestinationType}:</span>
                    {n.DestinationTarget}
                  </span>
                </td>
                <td class="py-3 px-4">
                  {#if n.OutboundCLI}
                    <Badge variant="success">Enabled (Present as CLI)</Badge>
                  {:else}
                    <Badge variant="neutral">Inbound Only</Badge>
                  {/if}
                </td>
                <td class="py-3 px-4 text-slate-600 font-medium">{n.Provider || 'Default Trunk'}</td>
                <td class="py-3 px-4 text-right space-x-2">
                  <button
                    type="button"
                    class="text-sky-600 hover:text-sky-800 font-medium"
                    on:click={() => openEdit(n)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    class="text-rose-600 hover:text-rose-800 font-medium"
                    on:click={() => deleteNumber(n.ID, n.Number)}
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

<!-- NUMBER MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? 'Edit Telephone Number' : 'Add Telephone Number (DID)'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <div>
      <label for="tn-num" class="block font-medium text-slate-700 mb-1">Phone Number (DID / E.164)</label>
      <input
        id="tn-num"
        type="text"
        bind:value={form.Number}
        placeholder="e.g. +18005550100 or 02079460912"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono font-bold text-sm"
      />
    </div>

    <div>
      <label for="tn-desc" class="block font-medium text-slate-700 mb-1">Description / Department</label>
      <input
        id="tn-desc"
        type="text"
        bind:value={form.Description}
        placeholder="e.g. Main London Line or Customer Support Toll-Free"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="tn-desttype" class="block font-medium text-slate-700 mb-1">Inbound Destination</label>
        <select
          id="tn-desttype"
          bind:value={form.DestinationType}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
        >
          <option value="extension">User Extension</option>
          <option value="ivr">Auto Attendant Menu (IVR)</option>
          <option value="group">Ring Group</option>
          <option value="queue">Call Center Queue</option>
          <option value="voicemail">Voicemail Server</option>
        </select>
      </div>
      <div>
        <label for="tn-desttarget" class="block font-medium text-slate-700 mb-1">Target Number / Choice</label>
        {#if form.DestinationType === 'extension'}
          <select
            id="tn-desttarget"
            bind:value={form.DestinationTarget}
            class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
          >
            {#each availableExts as ext}
              <option value={ext.ExtensionNumber}>Ext {ext.ExtensionNumber} ({ext.DisplayName || 'User'})</option>
            {/each}
          </select>
        {:else if form.DestinationType === 'ivr'}
          <select
            id="tn-desttarget"
            bind:value={form.DestinationTarget}
            class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
          >
            {#each availableIvrs as ivr}
              <option value={ivr.Extension || ivr.Name}>{ivr.Name}</option>
            {/each}
          </select>
        {:else}
          <input
            id="tn-desttarget"
            type="text"
            bind:value={form.DestinationTarget}
            placeholder="e.g. 600 or 8000"
            required
            class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
          />
        {/if}
      </div>
    </div>

    <div>
      <label for="tn-prov" class="block font-medium text-slate-700 mb-1">SIP Trunk Provider</label>
      <input
        id="tn-prov"
        type="text"
        bind:value={form.Provider}
        placeholder="e.g. Primary Carrier Gateway"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="flex items-center space-x-2 pt-2">
      <input
        type="checkbox"
        id="tn_cli"
        bind:checked={form.OutboundCLI}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="tn_cli" class="text-slate-700 font-medium">Allow staff to present this number as Outbound Caller ID (CLI)</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Save Number'}
      </Button>
    </div>
  </form>
</Modal>
