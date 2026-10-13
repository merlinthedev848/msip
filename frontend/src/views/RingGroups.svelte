<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import DestinationPicker from '../components/DestinationPicker.svelte';

  interface GroupMemberItem {
    ID?: string;
    Member: string;
    Extension: string;
    DelaySec: number;
  }

  interface CallGroupItem {
    ID: string;
    Name: string;
    SiteID?: string;
    SiteName: string;
    InternalNumber: string;
    RingCadence: string;
    CallerNamePrefix: string;
    NoAnswerTimeout: number;
    NoAnswerDestination: string;
    Members: GroupMemberItem[];
    CreatedAt: string;
  }

  let groups: CallGroupItem[] = [];
  let availableExts: any[] = [];
  let sites: any[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // View state: 'list' or 'edit'
  let currentView: 'list' | 'edit' = 'list';
  let editingGroup: CallGroupItem | null = null;

  // Member Modal State
  let memberModalOpen = false;
  let editingMemberIndex = -1;
  let memberForm = {
    Member: '',
    Extension: '',
    DelaySec: 0
  };

  // Main Form Data matching Screenshot 4
  let form = {
    ID: '',
    Name: '',
    SiteID: '',
    SiteName: '',
    InternalNumber: '990',
    RingCadence: 'Default',
    CallerNamePrefix: '',
    NoAnswerTimeout: 15,
    NoAnswerDestination: 'Not Set',
    Members: [] as GroupMemberItem[]
  };

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [rgRes, extRes, siteRes] = await Promise.all([
        api.get<{ ring_groups: CallGroupItem[] }>('/ring-groups'),
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] })),
        api.get<{ sites: any[] }>('/sites').catch(() => ({ sites: [] }))
      ]);
      groups = rgRes.ring_groups || [];
      availableExts = extRes.extensions || [];
      sites = siteRes.sites || [];
    } catch (e: any) {
      loadError = e.message || 'Failed to load call distribution groups';
    } finally {
      loading = false;
    }
  }

  onMount(loadData);

  function openCreate() {
    editingGroup = null;
    const defaultSite = sites.length > 0 ? sites[0] : null;
    form = {
      ID: '',
      Name: '',
      SiteID: defaultSite ? defaultSite.ID : '',
      SiteName: defaultSite ? defaultSite.Name : '',
      InternalNumber: (990 + groups.length).toString(),
      RingCadence: 'Default',
      CallerNamePrefix: '',
      NoAnswerTimeout: 15,
      NoAnswerDestination: 'Not Set',
      Members: []
    };
    currentView = 'edit';
  }

  function openEdit(group: CallGroupItem) {
    editingGroup = group;
    form = {
      ID: group.ID,
      Name: group.Name || '',
      SiteID: group.SiteID || '',
      SiteName: group.SiteName || (sites.length > 0 ? sites[0].Name : ''),
      InternalNumber: group.InternalNumber || '990',
      RingCadence: group.RingCadence || 'Default',
      CallerNamePrefix: group.CallerNamePrefix || '',
      NoAnswerTimeout: group.NoAnswerTimeout || 15,
      NoAnswerDestination: group.NoAnswerDestination || 'Not Set',
      Members: (group.Members || []).map(m => ({
        Member: m.Member,
        Extension: m.Extension || '',
        DelaySec: m.DelaySec ?? 0
      }))
    };
    currentView = 'edit';
  }

  function backToList() {
    currentView = 'list';
    editingGroup = null;
  }

  async function handleSave() {
    if (!form.Name.trim()) return toast.error('Group Name is required');
    if (!form.InternalNumber.trim()) return toast.error('Internal Number is required');

    // Find site name if site ID chosen
    if (form.SiteID) {
      const matched = sites.find(s => s.ID === form.SiteID);
      if (matched) form.SiteName = matched.Name;
    }

    saving = true;
    try {
      const payload: any = {
        Name: form.Name.trim(),
        InternalNumber: form.InternalNumber.trim(),
        SiteID: form.SiteID || null,
        SiteName: form.SiteName || '',
        RingCadence: form.RingCadence,
        CallerNamePrefix: form.CallerNamePrefix.trim(),
        NoAnswerTimeout: Number(form.NoAnswerTimeout),
        NoAnswerDestination: form.NoAnswerDestination,
        Members: form.Members.map(m => ({
          Member: m.Member,
          Extension: m.Extension,
          DelaySec: Number(m.DelaySec)
        }))
      };

      if (editingGroup) {
        await api.put(`/ring-groups/${editingGroup.ID}`, payload);
        toast.success(`Group '${form.Name}' updated`);
      } else {
        await api.post('/ring-groups', payload);
        toast.success(`Group '${form.Name}' created`);
      }

      await loadData();
      currentView = 'list';
    } catch (e: any) {
      toastError(e, 'Failed to save call distribution group');
    } finally {
      saving = false;
    }
  }

  async function deleteGroup(group: CallGroupItem) {
    if (!confirm(`Are you sure you want to delete call distribution group '${group.Name}'?`)) return;
    try {
      await api.delete(`/ring-groups/${group.ID}`);
      toast.success(`Group '${group.Name}' deleted`);
      loadData();
    } catch (e: any) {
      toastError(e, 'Failed to delete group');
    }
  }

  // Member CRUD
  function openAddMember() {
    editingMemberIndex = -1;
    const defaultExt = availableExts.length > 0 ? availableExts[0] : null;
    memberForm = {
      Member: defaultExt ? (defaultExt.DisplayName || `${defaultExt.FirstName || ''} ${defaultExt.LastName || ''}`.trim() || `User ${defaultExt.ExtensionNumber}`) : '',
      Extension: defaultExt ? defaultExt.ExtensionNumber : '',
      DelaySec: 0
    };
    memberModalOpen = true;
  }

  function editMember(index: number) {
    editingMemberIndex = index;
    memberForm = { ...form.Members[index] };
    memberModalOpen = true;
  }

  function saveMember() {
    if (!memberForm.Member.trim()) return toast.error('Member name is required');
    if (editingMemberIndex >= 0) {
      form.Members[editingMemberIndex] = { ...memberForm };
    } else {
      form.Members = [...form.Members, { ...memberForm }];
    }
    memberModalOpen = false;
  }

  function removeMember(index: number) {
    form.Members = form.Members.filter((_, i) => i !== index);
  }

  $: filteredList = groups.filter(g =>
    !search.trim() ||
    `${g.Name} ${g.InternalNumber} ${g.SiteName}`.toLowerCase().includes(search.trim().toLowerCase())
  );
</script>

<div class="p-6 bg-[#EEF2F6] min-h-screen text-slate-800">
  {#if currentView === 'list'}
    <!-- LIST VIEW -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-5 mb-5 border-b border-slate-200 gap-4">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 rounded-lg bg-sky-100 text-[#0284C7] flex items-center justify-center font-bold shadow-sm">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>
        <div>
          <h1 class="text-xl font-bold text-slate-900 tracking-tight">Call Distribution Groups (Ring Groups)</h1>
          <p class="text-xs text-slate-500">Configure hunt groups, staggered delay member ringing, and failover routing</p>
        </div>
      </div>

      <div class="flex items-center space-x-3">
        <div class="relative">
          <input
            type="text"
            bind:value={search}
            placeholder="Search groups..."
            class="pl-9 pr-3 py-1.5 bg-white border border-slate-300 rounded-md text-xs text-slate-700 focus:outline-none focus:ring-1 focus:ring-[#0284C7] w-52 shadow-sm"
          />
          <svg class="w-4 h-4 text-slate-400 absolute left-2.5 top-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <button
          type="button"
          on:click={openCreate}
          class="px-4 py-2 bg-[#0284C7] hover:bg-[#0369A1] text-white text-xs font-bold rounded-md shadow-sm transition-colors flex items-center space-x-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H8" />
          </svg>
          <span>Add Group</span>
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
      <table class="w-full text-left text-xs text-slate-700">
        <thead class="bg-[#2B527A] text-white font-semibold">
          <tr>
            <th class="px-4 py-3">Group Name</th>
            <th class="px-4 py-3">Internal Ext</th>
            <th class="px-4 py-3">Office Site</th>
            <th class="px-4 py-3">Members</th>
            <th class="px-4 py-3">Timeout (s)</th>
            <th class="px-4 py-3">No Answer Destination</th>
            <th class="px-4 py-3 text-right">Options</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          {#each filteredList as item}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="px-4 py-3 font-bold text-slate-900">{item.Name}</td>
              <td class="px-4 py-3 text-slate-600 font-mono font-bold">{item.InternalNumber}</td>
              <td class="px-4 py-3 text-slate-600">{item.SiteName || '—'}</td>
              <td class="px-4 py-3">
                <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-sky-100 text-sky-800">
                  {(item.Members || []).length} Members
                </span>
              </td>
              <td class="px-4 py-3 text-slate-600">{item.NoAnswerTimeout}s</td>
              <td class="px-4 py-3 text-slate-600 truncate max-w-xs">{item.NoAnswerDestination || 'Not Set'}</td>
              <td class="px-4 py-3 text-right">
                <div class="inline-flex items-center space-x-2">
                  <button
                    type="button"
                    on:click={() => openEdit(item)}
                    title="Edit Group"
                    class="p-1 text-slate-600 hover:text-[#0284C7] hover:bg-sky-50 rounded"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                  </button>
                  <button
                    type="button"
                    on:click={() => deleteGroup(item)}
                    title="Delete Group"
                    class="p-1 text-slate-600 hover:text-rose-600 hover:bg-rose-50 rounded"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          {/each}
          {#if filteredList.length === 0}
            <tr>
              <td colspan="7" class="px-4 py-8 text-center text-slate-400">
                No call distribution groups found. Click "Add Group" to create one.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>

  {:else}
    <!-- EDIT VIEW MATCHING SCREENSHOT 4 -->
    <!-- Top Bar with group title -->
    <div class="flex items-center space-x-2.5 pb-4 mb-4 border-b border-slate-300">
      <div class="text-[#0284C7] flex items-center space-x-1 font-bold text-xs">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </div>
      <h1 class="text-sm font-bold text-slate-900 tracking-tight">
        {editingGroup ? `Editing call distribution group '${form.Name}'` : 'Create Call Distribution Group'}
      </h1>
    </div>

    <form on:submit|preventDefault={handleSave} class="space-y-4">
      <!-- CARD 1: Group Details (2-Column Layout matching Screenshot 4) -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
          Group Details
        </div>
        <div class="p-5 grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Column 1 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Group Name</label>
              <input
                type="text"
                bind:value={form.Name}
                placeholder="e.g. Sales Ring Group"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-medium"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Site</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 14v3m4-3v3m4-3v3M3 21h18M3 10h18M3 7l9-4 9 4M4 10h16v11H4V10z"/></svg>
                </span>
                <select
                  bind:value={form.SiteID}
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                >
                  <option value="">Please select</option>
                  {#each sites as s}
                    <option value={s.ID}>{s.Name}</option>
                  {/each}
                  {#if sites.length === 0}
                    <option value="" disabled>No sites configured</option>
                  {/if}
                </select>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Internal Number</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg>
                </span>
                <input
                  type="text"
                  bind:value={form.InternalNumber}
                  placeholder="e.g. 990"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-bold"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Ring Cadence</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/></svg>
                </span>
                <select
                  bind:value={form.RingCadence}
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                >
                  <option value="Default">Default</option>
                  <option value="Standard">Standard</option>
                  <option value="Ring-Ring">Ring-Ring</option>
                  <option value="Short-Short">Short-Short</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Column 2 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Caller Name Prefix</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400 font-bold text-xs">“</span>
                <input
                  type="text"
                  bind:value={form.CallerNamePrefix}
                  placeholder="e.g. Front Door (9)"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">No Answer Timeout</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" stroke-linecap="round" d="M12 7v5l3 3"/></svg>
                </span>
                <input
                  type="number"
                  bind:value={form.NoAnswerTimeout}
                  min="1"
                  max="300"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>
            </div>

            <div>
              <DestinationPicker
                label="No Answer Destination"
                bind:value={form.NoAnswerDestination}
                on:change={(e) => form.NoAnswerDestination = e.detail}
              />
            </div>
          </div>
        </div>

        <!-- Action Buttons (Save disk icon + Back) -->
        <div class="px-5 py-3 bg-slate-50 border-t border-slate-200 flex items-center space-x-2">
          <button
            type="submit"
            disabled={saving}
            class="px-3.5 py-1.5 bg-[#0284C7] hover:bg-[#0369A1] text-white text-xs font-bold rounded shadow-sm transition-colors flex items-center space-x-1.5"
          >
            <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20">
              <path d="M7.707 10.293a1 1 0 10-1.414 1.414l3 3a1 1 0 001.414 0l6-6a1 1 0 00-1.414-1.414H11V4a1 1 0 10-2 0v6.586l-1.293-1.293z" />
            </svg>
            <span>{saving ? 'Saving...' : 'Save'}</span>
          </button>
          <button
            type="button"
            on:click={backToList}
            class="px-3 py-1.5 bg-white border border-slate-300 hover:bg-slate-100 text-slate-700 text-xs font-semibold rounded shadow-sm transition-colors flex items-center space-x-1"
          >
            <span>↶</span>
            <span>Back</span>
          </button>
        </div>
      </div>

      <!-- CARD 2: Members (Matching Screenshot 4 Bottom Card) -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
          Members
        </div>
        <div class="p-5">
          <!-- Top + Add Button -->
          <div class="mb-4">
            <button
              type="button"
              on:click={openAddMember}
              class="px-3 py-1.5 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 text-xs font-semibold rounded shadow-sm transition-colors flex items-center space-x-1"
            >
              <span>+ Add</span>
            </button>
          </div>

          <!-- Table with Navy Header matching Screenshot 4 -->
          <div class="border border-slate-200 rounded overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-[#2B527A] text-white font-semibold">
                <tr>
                  <th class="px-4 py-2.5">Member</th>
                  <th class="px-4 py-2.5">Delay (s)</th>
                  <th class="px-4 py-2.5 text-right">Options</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 bg-white">
                {#each form.Members as mem, idx}
                  <tr class="hover:bg-slate-50 transition-colors">
                    <td class="px-4 py-2.5 font-medium text-slate-800">
                      {mem.Member} {mem.Extension ? `(${mem.Extension})` : ''}
                    </td>
                    <td class="px-4 py-2.5 text-slate-600 font-mono">{mem.DelaySec}</td>
                    <td class="px-4 py-2.5 text-right">
                      <div class="inline-flex items-center space-x-2">
                        <button
                          type="button"
                          on:click={() => editMember(idx)}
                          class="p-1 text-slate-500 hover:text-[#0284C7]"
                          title="Edit Member"
                        >
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                        </button>
                        <button
                          type="button"
                          on:click={() => removeMember(idx)}
                          class="p-1 text-slate-500 hover:text-rose-600"
                          title="Remove Member"
                        >
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
                {#if form.Members.length === 0}
                  <tr>
                    <td colspan="3" class="px-4 py-6 text-center text-slate-400 bg-slate-50/50">
                      No members added to this distribution group. Click "+ Add" to add extensions.
                    </td>
                  </tr>
                {/if}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </form>
  {/if}
</div>

<!-- Modal: Ring Group Member Editor -->
{#if memberModalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-[2px]">
    <div class="bg-white rounded-lg shadow-xl border border-slate-200 w-full max-w-md overflow-hidden">
      <div class="px-5 py-3.5 bg-[#2B527A] text-white font-bold text-xs flex justify-between items-center">
        <span>{editingMemberIndex >= 0 ? 'Edit Group Member' : 'Add Group Member'}</span>
        <button type="button" on:click={() => memberModalOpen = false} class="text-white/80 hover:text-white">✕</button>
      </div>
      <div class="p-5 space-y-4">
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Select User / Extension</label>
          <select
            bind:value={memberForm.Extension}
            on:change={(e) => {
              const ext = availableExts.find(x => x.ExtensionNumber === e.currentTarget.value);
              if (ext) {
                memberForm.Member = ext.DisplayName || `${ext.FirstName || ''} ${ext.LastName || ''}`.trim() || `User ${ext.ExtensionNumber}`;
              }
            }}
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          >
            {#each availableExts as ext}
              <option value={ext.ExtensionNumber}>
                {ext.ExtensionNumber} - {ext.DisplayName || `${ext.FirstName || ''} ${ext.LastName || ''}`.trim() || 'Internal User'}
              </option>
            {/each}
          </select>
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Member Display Name</label>
          <input
            type="text"
            bind:value={memberForm.Member}
            placeholder="e.g. Agent Name"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Delay Before Ringing (Seconds)</label>
          <input
            type="number"
            bind:value={memberForm.DelaySec}
            min="0"
            max="120"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
          <p class="text-[11px] text-slate-400 mt-1">Set to 0 for immediate ringing alongside the group.</p>
        </div>
      </div>
      <div class="px-5 py-3 bg-slate-50 border-t border-slate-200 flex justify-end space-x-2">
        <button
          type="button"
          on:click={() => memberModalOpen = false}
          class="px-3 py-1.5 bg-white border border-slate-300 text-slate-700 font-semibold text-xs rounded hover:bg-slate-50"
        >
          Cancel
        </button>
        <button
          type="button"
          on:click={saveMember}
          class="px-4 py-1.5 bg-[#0284C7] hover:bg-[#0369A1] text-white font-bold text-xs rounded shadow-sm"
        >
          Save Member
        </button>
      </div>
    </div>
  </div>
{/if}
