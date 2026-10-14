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

  interface TimeCondition {
    id: string;
    name: string;
    time_start: string;
    time_end: string;
    days_of_week: string;
    matched_dest_type: string;
    matched_dest_val: string;
    unmatched_dest_type: string;
    unmatched_dest_val: string;
    is_active: boolean;
  }

  interface Holiday {
    id: string;
    name: string;
    start_date: string;
    end_date: string;
    dest_type: string;
    dest_val: string;
    is_active: boolean;
  }

  interface OptionItem {
    id: string;
    number: string;
    label: string;
  }

  let activeTab: 'conditions' | 'holidays' = 'conditions';

  let conditions: TimeCondition[] = [];
  let holidays: Holiday[] = [];
  let availableExts: OptionItem[] = [];
  let availableIvrs: OptionItem[] = [];
  let availableRingGroups: OptionItem[] = [];
  let availableQueues: OptionItem[] = [];
  let availableAi: OptionItem[] = [];

  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // Time Condition Modal
  let condModalOpen = false;
  let editingCondId: string | null = null;
  let condForm = {
    name: '',
    time_start: '09:00',
    time_end: '17:00',
    days_of_week: 'mon-fri',
    matched_dest_type: 'ivr',
    matched_dest_val: '',
    unmatched_dest_type: 'voicemail',
    unmatched_dest_val: '',
    is_active: true
  };

  // Holiday Modal
  let holidayModalOpen = false;
  let editingHolidayId: string | null = null;
  let holidayForm = {
    name: '',
    start_date: '',
    end_date: '',
    dest_type: 'voicemail',
    dest_val: '',
    is_active: true
  };

  $: filteredConditions = conditions.filter(c =>
    !search.trim() || c.name.toLowerCase().includes(search.trim().toLowerCase())
  );

  $: filteredHolidays = holidays.filter(h =>
    !search.trim() || h.name.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [tcRes, holRes, extRes, ivrRes, rgRes, qRes, aiRes] = await Promise.all([
        api.get<{ time_conditions: TimeCondition[] }>('/time-conditions').catch(() => ({ time_conditions: [] })),
        api.get<{ holidays: Holiday[] }>('/holidays').catch(() => ({ holidays: [] })),
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] })),
        api.get<{ ivrs: any[] }>('/ivrs').catch(() => ({ ivrs: [] })),
        api.get<{ ring_groups: any[] }>('/ring-groups').catch(() => ({ ring_groups: [] })),
        api.get<{ queues: any[] }>('/queues').catch(() => ({ queues: [] })),
        api.get<{ assistants: any[] }>('/ai/assistants').catch(() => ({ assistants: [] }))
      ]);

      conditions = tcRes.time_conditions || [];
      holidays = holRes.holidays || [];

      availableExts = (extRes.extensions || []).map((e: any) => ({
        id: e.ID,
        number: e.ExtensionNumber,
        label: `Ext ${e.ExtensionNumber} (${e.DisplayName || 'User'})`
      }));

      availableIvrs = (ivrRes.ivrs || []).map((i: any) => ({
        id: i.ID,
        number: i.Extension || i.Name,
        label: `IVR: ${i.Name} (${i.Extension || 'No Ext'})`
      }));

      availableRingGroups = (rgRes.ring_groups || []).map((rg: any) => ({
        id: rg.ID,
        number: rg.ExtensionNumber || rg.Name,
        label: `Ring Group: ${rg.Name} (${rg.ExtensionNumber})`
      }));

      availableQueues = (qRes.queues || []).map((q: any) => ({
        id: q.ID,
        number: q.Extension || q.Name,
        label: `Queue: ${q.Name} (${q.Extension || ''})`
      }));

      availableAi = (aiRes.assistants || []).map((a: any) => ({
        id: a.ID,
        number: a.Extension || a.Name,
        label: `🤖 AI: ${a.Name} (${a.Extension || ''})`
      }));
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreateCond() {
    editingCondId = null;
    condForm = {
      name: '',
      time_start: '09:00',
      time_end: '17:00',
      days_of_week: 'mon-fri',
      matched_dest_type: availableIvrs.length > 0 ? 'ivr' : (availableExts.length > 0 ? 'extension' : 'hangup'),
      matched_dest_val: availableIvrs[0]?.number || availableExts[0]?.number || '',
      unmatched_dest_type: 'voicemail',
      unmatched_dest_val: availableExts[0]?.number || '',
      is_active: true
    };
    condModalOpen = true;
  }

  function openEditCond(c: TimeCondition) {
    editingCondId = c.id;
    condForm = {
      name: c.name,
      time_start: c.time_start || '09:00',
      time_end: c.time_end || '17:00',
      days_of_week: c.days_of_week || 'mon-fri',
      matched_dest_type: c.matched_dest_type || 'ivr',
      matched_dest_val: c.matched_dest_val || '',
      unmatched_dest_type: c.unmatched_dest_type || 'voicemail',
      unmatched_dest_val: c.unmatched_dest_val || '',
      is_active: c.is_active
    };
    condModalOpen = true;
  }

  async function saveCond(e: Event) {
    e.preventDefault();
    if (!condForm.name.trim()) return toast.error('Condition name is required');
    saving = true;
    try {
      if (editingCondId) {
        await api.put(`/time-conditions/${editingCondId}`, condForm);
        toast.success('Time condition updated');
      } else {
        await api.post('/time-conditions', condForm);
        toast.success('Time condition created');
      }
      condModalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteCond(id: string, name: string) {
    if (!confirm(`Delete time condition "${name}"?`)) return;
    try {
      await api.del(`/time-conditions/${id}`);
      toast.success('Condition deleted');
      conditions = conditions.filter(c => c.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  // HOLIDAYS
  function openCreateHoliday() {
    editingHolidayId = null;
    const today = new Date().toISOString().split('T')[0];
    holidayForm = {
      name: '',
      start_date: today,
      end_date: today,
      dest_type: 'voicemail',
      dest_val: availableExts[0]?.number || '',
      is_active: true
    };
    holidayModalOpen = true;
  }

  function openEditHoliday(h: Holiday) {
    editingHolidayId = h.id;
    holidayForm = {
      name: h.name,
      start_date: h.start_date || '',
      end_date: h.end_date || '',
      dest_type: h.dest_type || 'voicemail',
      dest_val: h.dest_val || '',
      is_active: h.is_active
    };
    holidayModalOpen = true;
  }

  async function saveHoliday(e: Event) {
    e.preventDefault();
    if (!holidayForm.name.trim()) return toast.error('Holiday name is required');
    if (!holidayForm.start_date) return toast.error('Start date is required');
    saving = true;
    try {
      if (editingHolidayId) {
        await api.put(`/holidays/${editingHolidayId}`, holidayForm);
        toast.success('Holiday updated');
      } else {
        await api.post('/holidays', holidayForm);
        toast.success('Holiday created');
      }
      holidayModalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteHoliday(id: string, name: string) {
    if (!confirm(`Delete holiday "${name}"?`)) return;
    try {
      await api.del(`/holidays/${id}`);
      toast.success('Holiday deleted');
      holidays = holidays.filter(h => h.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Time Conditions & Holidays"
    subtitle="Define business office hours, lunch schedules, and public holiday call routing paths."
  >
    <div class="flex items-center space-x-3">
      {#if activeTab === 'conditions'}
        <Button variant="primary" on:click={openCreateCond}>
          <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Time Condition
        </Button>
      {:else}
        <Button variant="primary" on:click={openCreateHoliday}>
          <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          Add Holiday
        </Button>
      {/if}
    </div>
  </PageHeader>

  <!-- TAB NAVIGATION -->
  <div class="flex border-b border-slate-200 space-x-6 text-sm font-medium">
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'conditions' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'conditions'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span>Office Hours / Time Conditions ({conditions.length})</span>
    </button>
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'holidays' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'holidays'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <span>Holiday Calendar ({holidays.length})</span>
    </button>
  </div>

  <!-- SEARCH & FILTER BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder={activeTab === 'conditions' ? 'Search time conditions...' : 'Search holiday schedules...'}
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
    <div class="p-12 text-center text-slate-400 text-sm">Loading schedules and rules...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if activeTab === 'conditions'}
    <!-- TIME CONDITIONS TAB -->
    {#if filteredConditions.length === 0}
      <EmptyState
        title="No Time Conditions Defined"
        description="Create office hour windows so inbound callers get routed to reception during work hours and voicemail after hours."
      >
        <Button variant="primary" on:click={openCreateCond}>Create Business Hours</Button>
      </EmptyState>
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Condition Name</th>
                <th class="py-3 px-4">Office Hours</th>
                <th class="py-3 px-4">Active Days</th>
                <th class="py-3 px-4">Open Destination</th>
                <th class="py-3 px-4">Closed Destination</th>
                <th class="py-3 px-4">Status</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredConditions as c}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-semibold text-slate-900">{c.name}</td>
                  <td class="py-3 px-4 font-mono font-medium text-sky-700">
                    {c.time_start} – {c.time_end}
                  </td>
                  <td class="py-3 px-4 capitalize">
                    <span class="px-2 py-0.5 bg-slate-100 rounded text-slate-700 font-mono text-[11px]">{c.days_of_week}</span>
                  </td>
                  <td class="py-3 px-4">
                    <span class="inline-flex items-center text-emerald-700 font-medium">
                      <span class="w-2 h-2 rounded-full bg-emerald-500 mr-1.5"></span>
                      {c.matched_dest_type.toUpperCase()}: {c.matched_dest_val}
                    </span>
                  </td>
                  <td class="py-3 px-4">
                    <span class="inline-flex items-center text-slate-600">
                      <span class="w-2 h-2 rounded-full bg-amber-500 mr-1.5"></span>
                      {c.unmatched_dest_type.toUpperCase()}: {c.unmatched_dest_val}
                    </span>
                  </td>
                  <td class="py-3 px-4">
                    {#if c.is_active}
                      <Badge variant="success">Active</Badge>
                    {:else}
                      <Badge variant="neutral">Disabled</Badge>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-right space-x-2">
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 font-medium"
                      on:click={() => openEditCond(c)}
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="text-rose-600 hover:text-rose-800 font-medium"
                      on:click={() => deleteCond(c.id, c.name)}
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
    <!-- HOLIDAYS TAB -->
    {#if filteredHolidays.length === 0}
      <EmptyState
        title="No Holidays Configured"
        description="Add statutory bank holidays and corporate shutdowns to automatically override standard office hours."
      >
        <Button variant="primary" on:click={openCreateHoliday}>Add Public Holiday</Button>
      </EmptyState>
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Holiday Title</th>
                <th class="py-3 px-4">Date Window</th>
                <th class="py-3 px-4">Special Closed Routing</th>
                <th class="py-3 px-4">Status</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredHolidays as h}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-semibold text-slate-900">{h.name}</td>
                  <td class="py-3 px-4 font-mono font-medium text-slate-700">
                    {h.start_date} {h.end_date && h.end_date !== h.start_date ? `to ${h.end_date}` : '(1 Day)'}
                  </td>
                  <td class="py-3 px-4">
                    <span class="inline-flex items-center text-purple-700 font-medium">
                      <span class="w-2 h-2 rounded-full bg-purple-500 mr-1.5"></span>
                      {h.dest_type.toUpperCase()}: {h.dest_val}
                    </span>
                  </td>
                  <td class="py-3 px-4">
                    {#if h.is_active}
                      <Badge variant="success">Active</Badge>
                    {:else}
                      <Badge variant="neutral">Disabled</Badge>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-right space-x-2">
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 font-medium"
                      on:click={() => openEditHoliday(h)}
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="text-rose-600 hover:text-rose-800 font-medium"
                      on:click={() => deleteHoliday(h.id, h.name)}
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

<!-- TIME CONDITION MODAL -->
<Modal
  open={condModalOpen}
  title={editingCondId ? 'Edit Time Condition' : 'Create Time Condition'}
  on:close={() => condModalOpen = false}
>
  <form on:submit={saveCond} class="space-y-4 text-xs">
    <div>
      <label for="tc-name" class="block font-medium text-slate-700 mb-1">Condition Name</label>
      <input
        id="tc-name"
        type="text"
        bind:value={condForm.name}
        placeholder="e.g. Standard Office Hours"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="tc-start" class="block font-medium text-slate-700 mb-1">Open Time (24h)</label>
        <input
          id="tc-start"
          type="time"
          bind:value={condForm.time_start}
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
      <div>
        <label for="tc-end" class="block font-medium text-slate-700 mb-1">Close Time (24h)</label>
        <input
          id="tc-end"
          type="time"
          bind:value={condForm.time_end}
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
    </div>

    <div>
      <label for="tc-days" class="block font-medium text-slate-700 mb-1">Days of Week</label>
      <select
        id="tc-days"
        bind:value={condForm.days_of_week}
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
      >
        <option value="mon-fri">Monday through Friday (Business Days)</option>
        <option value="mon-sat">Monday through Saturday</option>
        <option value="all">Every Day (7 Days / Week)</option>
        <option value="sat-sun">Weekends Only (Saturday & Sunday)</option>
      </select>
    </div>

    <!-- During Hours Routing -->
    <div class="p-3 bg-emerald-50/70 border border-emerald-200 rounded-lg space-y-3">
      <div class="font-semibold text-emerald-900 flex items-center">
        <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 mr-2"></span>
        During Office Hours (Matched)
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="tc-mtype" class="block text-emerald-800 mb-1 font-semibold">Destination Type</label>
          <select
            id="tc-mtype"
            bind:value={condForm.matched_dest_type}
            class="w-full px-2.5 py-1.5 border border-emerald-300 rounded bg-white text-xs"
          >
            <option value="ivr">Digital Receptionist (IVR)</option>
            <option value="ai">🤖 AI Virtual Receptionist</option>
            <option value="extension">Extension Handset</option>
            <option value="ringgroup">Ring Group</option>
            <option value="queue">Call Center Queue</option>
            <option value="voicemail">Voicemail Box</option>
            <option value="hangup">Disconnect Call</option>
          </select>
        </div>
        <div>
          <label for="tc-mval" class="block text-emerald-800 mb-1 font-semibold">Target Choice</label>
          {#if condForm.matched_dest_type === 'extension' || condForm.matched_dest_type === 'voicemail'}
            <select
              id="tc-mval"
              bind:value={condForm.matched_dest_val}
              class="w-full px-2.5 py-1.5 border border-emerald-300 rounded bg-white text-xs"
            >
              {#each availableExts as ext}
                <option value={ext.number}>{ext.label}</option>
              {/each}
            </select>
          {:else if condForm.matched_dest_type === 'ivr'}
            <select
              id="tc-mval"
              bind:value={condForm.matched_dest_val}
              class="w-full px-2.5 py-1.5 border border-emerald-300 rounded bg-white text-xs"
            >
              {#each availableIvrs as ivr}
                <option value={ivr.number}>{ivr.label}</option>
              {/each}
            </select>
          {:else if condForm.matched_dest_type === 'ai'}
            <select
              id="tc-mval"
              bind:value={condForm.matched_dest_val}
              class="w-full px-2.5 py-1.5 border border-emerald-300 rounded bg-white text-xs"
            >
              {#each availableAi as ai}
                <option value={ai.number}>{ai.label}</option>
              {/each}
            </select>
          {:else if condForm.matched_dest_type === 'ringgroup'}
            <select
              id="tc-mval"
              bind:value={condForm.matched_dest_val}
              class="w-full px-2.5 py-1.5 border border-emerald-300 rounded bg-white text-xs"
            >
              {#each availableRingGroups as rg}
                <option value={rg.number}>{rg.label}</option>
              {/each}
            </select>
          {:else if condForm.matched_dest_type === 'queue'}
            <select
              id="tc-mval"
              bind:value={condForm.matched_dest_val}
              class="w-full px-2.5 py-1.5 border border-emerald-300 rounded bg-white text-xs"
            >
              {#each availableQueues as q}
                <option value={q.number}>{q.label}</option>
              {/each}
            </select>
          {:else if condForm.matched_dest_type === 'hangup'}
            <input
              id="tc-mval"
              type="text"
              disabled
              value="Terminates Call Immediately"
              class="w-full px-2.5 py-1.5 border border-slate-200 rounded bg-slate-100 text-slate-500 text-xs italic"
            />
          {:else}
            <input
              id="tc-mval"
              type="text"
              bind:value={condForm.matched_dest_val}
              placeholder="Destination value"
              class="w-full px-2.5 py-1.5 border border-emerald-300 rounded bg-white text-xs"
            />
          {/if}
        </div>
      </div>
    </div>

    <!-- Outside Hours Routing -->
    <div class="p-3 bg-amber-50/70 border border-amber-200 rounded-lg space-y-3">
      <div class="font-semibold text-amber-900 flex items-center">
        <span class="w-2.5 h-2.5 rounded-full bg-amber-500 mr-2"></span>
        After Hours / Closed (Unmatched)
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="tc-umtype" class="block text-amber-800 mb-1 font-semibold">Destination Type</label>
          <select
            id="tc-umtype"
            bind:value={condForm.unmatched_dest_type}
            class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white text-xs"
          >
            <option value="voicemail">General Voicemail Box</option>
            <option value="ai">🤖 AI Virtual Receptionist (Night Service)</option>
            <option value="extension">Night Service Extension</option>
            <option value="ivr">After-Hours IVR Announcement</option>
            <option value="ringgroup">After-Hours Ring Group</option>
            <option value="queue">After-Hours Queue</option>
            <option value="hangup">Disconnect Call</option>
          </select>
        </div>
        <div>
          <label for="tc-umval" class="block text-amber-800 mb-1 font-semibold">Target Box / Extension</label>
          {#if condForm.unmatched_dest_type === 'voicemail' || condForm.unmatched_dest_type === 'extension'}
            <select
              id="tc-umval"
              bind:value={condForm.unmatched_dest_val}
              class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white text-xs"
            >
              {#each availableExts as ext}
                <option value={ext.number}>{ext.label}</option>
              {/each}
            </select>
          {:else if condForm.unmatched_dest_type === 'ai'}
            <select
              id="tc-umval"
              bind:value={condForm.unmatched_dest_val}
              class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white text-xs"
            >
              {#each availableAi as ai}
                <option value={ai.number}>{ai.label}</option>
              {/each}
            </select>
          {:else if condForm.unmatched_dest_type === 'ivr'}
            <select
              id="tc-umval"
              bind:value={condForm.unmatched_dest_val}
              class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white text-xs"
            >
              {#each availableIvrs as ivr}
                <option value={ivr.number}>{ivr.label}</option>
              {/each}
            </select>
          {:else if condForm.unmatched_dest_type === 'ringgroup'}
            <select
              id="tc-umval"
              bind:value={condForm.unmatched_dest_val}
              class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white text-xs"
            >
              {#each availableRingGroups as rg}
                <option value={rg.number}>{rg.label}</option>
              {/each}
            </select>
          {:else if condForm.unmatched_dest_type === 'queue'}
            <select
              id="tc-umval"
              bind:value={condForm.unmatched_dest_val}
              class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white text-xs"
            >
              {#each availableQueues as q}
                <option value={q.number}>{q.label}</option>
              {/each}
            </select>
          {:else if condForm.unmatched_dest_type === 'hangup'}
            <input
              id="tc-umval"
              type="text"
              disabled
              value="Terminates Call Immediately"
              class="w-full px-2.5 py-1.5 border border-slate-200 rounded bg-slate-100 text-slate-500 text-xs italic"
            />
          {:else}
            <input
              id="tc-umval"
              type="text"
              bind:value={condForm.unmatched_dest_val}
              placeholder="e.g. 1001"
              class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white text-xs"
            />
          {/if}
        </div>
      </div>
    </div>

    <div class="flex items-center space-x-2 pt-2">
      <input
        type="checkbox"
        id="cond_active"
        bind:checked={condForm.is_active}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="cond_active" class="text-slate-700 font-medium">Enable this condition rule</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => condModalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingCondId ? 'Save Changes' : 'Create Condition'}
      </Button>
    </div>
  </form>
</Modal>

<!-- HOLIDAY MODAL -->
<Modal
  open={holidayModalOpen}
  title={editingHolidayId ? 'Edit Holiday' : 'Add Holiday Schedule'}
  on:close={() => holidayModalOpen = false}
>
  <form on:submit={saveHoliday} class="space-y-4 text-xs">
    <div>
      <label for="h-name" class="block font-medium text-slate-700 mb-1">Holiday Name</label>
      <input
        id="h-name"
        type="text"
        bind:value={holidayForm.name}
        placeholder="e.g. Christmas Day or New Year's Day"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="h-sdate" class="block font-medium text-slate-700 mb-1">Start Date</label>
        <input
          id="h-sdate"
          type="date"
          bind:value={holidayForm.start_date}
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
      <div>
        <label for="h-edate" class="block font-medium text-slate-700 mb-1">End Date</label>
        <input
          id="h-edate"
          type="date"
          bind:value={holidayForm.end_date}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
    </div>

    <div class="p-3 bg-purple-50/70 border border-purple-200 rounded-lg space-y-3">
      <div class="font-semibold text-purple-900 flex items-center">
        <span class="w-2.5 h-2.5 rounded-full bg-purple-500 mr-2"></span>
        Holiday Caller Routing
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="h-dtype" class="block text-purple-800 mb-1 font-semibold">Routing Destination</label>
          <select
            id="h-dtype"
            bind:value={holidayForm.dest_type}
            class="w-full px-2.5 py-1.5 border border-purple-300 rounded bg-white text-xs"
          >
            <option value="voicemail">Holiday Voicemail Greeting</option>
            <option value="ai">🤖 AI Virtual Receptionist (Holiday Desk)</option>
            <option value="ivr">Holiday IVR Announcement</option>
            <option value="extension">On-Call Duty Extension</option>
            <option value="ringgroup">Emergency Ring Group</option>
            <option value="queue">Emergency Queue</option>
            <option value="hangup">Disconnect Call</option>
          </select>
        </div>
        <div>
          <label for="h-dval" class="block text-purple-800 mb-1 font-semibold">Target Choice</label>
          {#if holidayForm.dest_type === 'extension' || holidayForm.dest_type === 'voicemail'}
            <select
              id="h-dval"
              bind:value={holidayForm.dest_val}
              class="w-full px-2.5 py-1.5 border border-purple-300 rounded bg-white text-xs"
            >
              {#each availableExts as ext}
                <option value={ext.number}>{ext.label}</option>
              {/each}
            </select>
          {:else if holidayForm.dest_type === 'ai'}
            <select
              id="h-dval"
              bind:value={holidayForm.dest_val}
              class="w-full px-2.5 py-1.5 border border-purple-300 rounded bg-white text-xs"
            >
              {#each availableAi as ai}
                <option value={ai.number}>{ai.label}</option>
              {/each}
            </select>
          {:else if holidayForm.dest_type === 'ivr'}
            <select
              id="h-dval"
              bind:value={holidayForm.dest_val}
              class="w-full px-2.5 py-1.5 border border-purple-300 rounded bg-white text-xs"
            >
              {#each availableIvrs as ivr}
                <option value={ivr.number}>{ivr.label}</option>
              {/each}
            </select>
          {:else if holidayForm.dest_type === 'ringgroup'}
            <select
              id="h-dval"
              bind:value={holidayForm.dest_val}
              class="w-full px-2.5 py-1.5 border border-purple-300 rounded bg-white text-xs"
            >
              {#each availableRingGroups as rg}
                <option value={rg.number}>{rg.label}</option>
              {/each}
            </select>
          {:else if holidayForm.dest_type === 'queue'}
            <select
              id="h-dval"
              bind:value={holidayForm.dest_val}
              class="w-full px-2.5 py-1.5 border border-purple-300 rounded bg-white text-xs"
            >
              {#each availableQueues as q}
                <option value={q.number}>{q.label}</option>
              {/each}
            </select>
          {:else if holidayForm.dest_type === 'hangup'}
            <input
              id="h-dval"
              type="text"
              disabled
              value="Terminates Call Immediately"
              class="w-full px-2.5 py-1.5 border border-slate-200 rounded bg-slate-100 text-slate-500 text-xs italic"
            />
          {:else}
            <input
              id="h-dval"
              type="text"
              bind:value={holidayForm.dest_val}
              placeholder="e.g. 1001 or holiday_prompt"
              class="w-full px-2.5 py-1.5 border border-purple-300 rounded bg-white text-xs"
            />
          {/if}
        </div>
      </div>
    </div>

    <div class="flex items-center space-x-2 pt-2">
      <input
        type="checkbox"
        id="h_active"
        bind:checked={holidayForm.is_active}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="h_active" class="text-slate-700 font-medium">Enable this holiday rule</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => holidayModalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingHolidayId ? 'Save Changes' : 'Add Holiday'}
      </Button>
    </div>
  </form>
</Modal>
