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

  interface OutboundRestriction {
    id: string;
    name: string;
    pattern: string;
    barred_types: string;
    allowed_exts: string;
    pin_code: string;
    is_active: boolean;
  }

  interface AutoCLIPRecord {
    id: string;
    caller_number: string;
    original_callee_ext: string;
    trunk_name: string;
    created_at: string;
    expires_at: string;
  }

  let activeTab: 'restrictions' | 'autoclip' = 'restrictions';

  let restrictions: OutboundRestriction[] = [];
  let autoclipRecords: AutoCLIPRecord[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // Barring Modal
  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    name: '',
    pattern: '00X.|011X.',
    barred_types: 'international',
    allowed_exts: '',
    pin_code: '',
    is_active: true
  };

  $: filteredRestrictions = restrictions.filter(r =>
    !search.trim() || `${r.name} ${r.pattern} ${r.barred_types}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  $: filteredAutoCLIP = autoclipRecords.filter(a =>
    !search.trim() || `${a.caller_number} ${a.original_callee_ext} ${a.trunk_name}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [resData, clipData] = await Promise.all([
        api.get<{ restrictions: OutboundRestriction[] }>('/outbound-restrictions'),
        api.get<{ records: AutoCLIPRecord[] }>('/autoclip')
      ]);
      restrictions = resData.restrictions || [];
      autoclipRecords = clipData.records || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      name: '',
      pattern: '00X.|011X.',
      barred_types: 'international',
      allowed_exts: '',
      pin_code: '',
      is_active: true
    };
    modalOpen = true;
  }

  function applyPreset(type: 'intl' | 'premium' | 'all') {
    if (type === 'intl') {
      form.name = 'Block International Dialing';
      form.pattern = '00X.|011X.|+X.';
      form.barred_types = 'international';
    } else if (type === 'premium') {
      form.name = 'Block Premium Numbers';
      form.pattern = '090X.|098X.|1900X.';
      form.barred_types = 'premium';
    } else if (type === 'all') {
      form.name = 'Lock All Outbound Calls';
      form.pattern = 'X.';
      form.barred_types = 'all';
      form.pin_code = '1234';
    }
  }

  function openEdit(r: OutboundRestriction) {
    editingId = r.id;
    form = {
      name: r.name,
      pattern: r.pattern,
      barred_types: r.barred_types,
      allowed_exts: r.allowed_exts || '',
      pin_code: r.pin_code || '',
      is_active: r.is_active
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.name.trim()) return toast.error('Rule name is required');
    if (!form.pattern.trim()) return toast.error('Number dial pattern is required');

    saving = true;
    try {
      if (editingId) {
        await api.put(`/outbound-restrictions/${editingId}`, form);
        toast.success('Restriction rule updated');
      } else {
        await api.post('/outbound-restrictions', form);
        toast.success('Restriction rule created');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteRule(id: string, name: string) {
    if (!confirm(`Delete rule "${name}"?`)) return;
    try {
      await api.del(`/outbound-restrictions/${id}`);
      toast.success('Rule removed');
      restrictions = restrictions.filter(r => r.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  async function deleteClip(id: string) {
    try {
      await api.del(`/autoclip/${id}`);
      toast.success('AutoCLIP mapping cleared');
      autoclipRecords = autoclipRecords.filter(r => r.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Outbound Restriction & AutoCLIP"
    subtitle="Prevent toll fraud by barring international and premium numbers, and manage smart caller callbacks."
  >
    {#if activeTab === 'restrictions'}
      <Button variant="primary" on:click={openCreate}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Barring Rule
      </Button>
    {/if}
  </PageHeader>

  <!-- TAB SWITCHER -->
  <div class="flex border-b border-slate-200 space-x-6 text-sm font-medium">
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'restrictions' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'restrictions'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
      </svg>
      <span>Outbound Barring Rules ({restrictions.length})</span>
    </button>
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'autoclip' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'autoclip'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
      </svg>
      <span>AutoCLIP Smart Callbacks ({autoclipRecords.length})</span>
    </button>
  </div>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder={activeTab === 'restrictions' ? 'Filter barring rules...' : 'Filter callback numbers...'}
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
    <div class="p-12 text-center text-slate-400 text-sm">Loading security and callback rules...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if activeTab === 'restrictions'}
    {#if filteredRestrictions.length === 0}
      <EmptyState
        title="No Outbound Restrictions Configured"
        description="Prevent unauthorized dialing of costly international destinations, satellite numbers, or premium tariff codes."
      >
        <Button variant="primary" on:click={openCreate}>Create Barring Rule</Button>
      </EmptyState>
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Rule Name</th>
                <th class="py-3 px-4">Barred Prefix / Pattern</th>
                <th class="py-3 px-4">Type</th>
                <th class="py-3 px-4">Exempt Extensions</th>
                <th class="py-3 px-4">PIN Code Bypass</th>
                <th class="py-3 px-4">Status</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredRestrictions as r}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-semibold text-slate-900">{r.name}</td>
                  <td class="py-3 px-4 font-mono font-medium text-rose-700">{r.pattern}</td>
                  <td class="py-3 px-4 capitalize">
                    <span class="px-2 py-0.5 bg-rose-50 border border-rose-200 text-rose-700 rounded text-[11px] font-medium">
                      {r.barred_types}
                    </span>
                  </td>
                  <td class="py-3 px-4">
                    {#if r.allowed_exts}
                      <span class="font-mono text-slate-600">{r.allowed_exts}</span>
                    {:else}
                      <span class="text-slate-400 italic">None (All Blocked)</span>
                    {/if}
                  </td>
                  <td class="py-3 px-4">
                    {#if r.pin_code}
                      <span class="inline-flex items-center text-amber-700 font-mono">
                        <svg class="w-3.5 h-3.5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                        PIN Required
                      </span>
                    {:else}
                      <span class="text-slate-400">Strict Block</span>
                    {/if}
                  </td>
                  <td class="py-3 px-4">
                    {#if r.is_active}
                      <Badge variant="success">Enforced</Badge>
                    {:else}
                      <Badge variant="neutral">Disabled</Badge>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-right space-x-2">
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 font-medium"
                      on:click={() => openEdit(r)}
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="text-rose-600 hover:text-rose-800 font-medium"
                      on:click={() => deleteRule(r.id, r.name)}
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
    <!-- AUTOCLIP TAB -->
    <div class="p-4 bg-sky-50 border border-sky-200 rounded-lg text-xs text-sky-900 leading-relaxed flex items-start space-x-3">
      <svg class="w-5 h-5 text-sky-600 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <div>
        <span class="font-bold">What is AutoCLIP?</span>
        AutoCLIP remembers outgoing calls made by your extensions. When an outside customer calls back on your main number, the PBX automatically bypasses the IVR receptionist and rings the exact extension who dialed them.
      </div>
    </div>

    {#if filteredAutoCLIP.length === 0}
      <EmptyState
        title="No Active AutoCLIP Call Mappings"
        description="AutoCLIP records are automatically generated whenever extensions make outbound calls to external customers."
      />
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Customer Phone Number</th>
                <th class="py-3 px-4">Original Caller (Extension)</th>
                <th class="py-3 px-4">Trunk Gateway</th>
                <th class="py-3 px-4">Call Placed</th>
                <th class="py-3 px-4 text-right">Clear Memory</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredAutoCLIP as c}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-mono font-semibold text-slate-900">{c.caller_number}</td>
                  <td class="py-3 px-4">
                    <span class="px-2 py-0.5 bg-sky-50 text-sky-700 border border-sky-200 rounded font-mono font-medium">
                      Ext {c.original_callee_ext}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-slate-600">{c.trunk_name || 'Default SIP Trunk'}</td>
                  <td class="py-3 px-4 text-slate-500">{formatDateTime(c.created_at)}</td>
                  <td class="py-3 px-4 text-right">
                    <button
                      type="button"
                      class="text-rose-600 hover:text-rose-800 font-medium"
                      on:click={() => deleteClip(c.id)}
                    >
                      Clear
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

<!-- BARRING RULE MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? 'Edit Outbound Barring Rule' : 'New Outbound Barring Rule'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <!-- Quick Presets -->
    <div>
      <span class="block text-slate-600 font-medium mb-1.5">Apply Standard Preset:</span>
      <div class="flex space-x-2">
        <button
          type="button"
          class="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded text-[11px] font-medium transition-colors"
          on:click={() => applyPreset('intl')}
        >
          International (00/011)
        </button>
        <button
          type="button"
          class="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded text-[11px] font-medium transition-colors"
          on:click={() => applyPreset('premium')}
        >
          Premium Rate (090/1900)
        </button>
        <button
          type="button"
          class="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded text-[11px] font-medium transition-colors"
          on:click={() => applyPreset('all')}
        >
          Lock All Calls with PIN
        </button>
      </div>
    </div>

    <div>
      <label for="bar-name" class="block font-medium text-slate-700 mb-1">Rule Name</label>
      <input
        id="bar-name"
        type="text"
        bind:value={form.name}
        placeholder="e.g. Block High-Cost International"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div>
      <label for="bar-pattern" class="block font-medium text-slate-700 mb-1">Number Pattern / Prefix</label>
      <input
        id="bar-pattern"
        type="text"
        bind:value={form.pattern}
        placeholder="e.g. 00X.|011X. (pipe separated)"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
      />
      <span class="text-[11px] text-slate-400 mt-0.5 block">Use X for single digit, . for one or more digits.</span>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="bar-types" class="block font-medium text-slate-700 mb-1">Category</label>
        <select
          id="bar-types"
          bind:value={form.barred_types}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
        >
          <option value="international">International Numbers</option>
          <option value="premium">Premium Rate Services</option>
          <option value="mobile">Cellular / Mobile Numbers</option>
          <option value="all">All Outbound Destinations</option>
        </select>
      </div>
      <div>
        <label for="bar-pin" class="block font-medium text-slate-700 mb-1">PIN Bypass Code (Optional)</label>
        <input
          id="bar-pin"
          type="text"
          bind:value={form.pin_code}
          placeholder="Leave blank for hard block"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
    </div>

    <div>
      <label for="bar-exts" class="block font-medium text-slate-700 mb-1">Exempt Extensions (Whitelist)</label>
      <input
        id="bar-exts"
        type="text"
        bind:value={form.allowed_exts}
        placeholder="e.g. 1000, 1001, 1002 (comma separated)"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
      />
      <span class="text-[11px] text-slate-400 mt-0.5 block">These extensions will never be restricted by this rule.</span>
    </div>

    <div class="flex items-center space-x-2 pt-2">
      <input
        type="checkbox"
        id="bar_active"
        bind:checked={form.is_active}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="bar_active" class="text-slate-700 font-medium">Enforce this restriction rule</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Create Barring Rule'}
      </Button>
    </div>
  </form>
</Modal>
