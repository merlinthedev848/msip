<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import DestinationPicker from '../components/DestinationPicker.svelte';

  interface HandsetKey {
    key: number;
    label: string;
    type: string;
    description: string;
  }

  interface ExtensionItem {
    ID: string;
    ExtensionNumber: string;
    PasswordHash?: string;
    IsActive: boolean;
    ProvisioningCode: string;
    FirstName: string;
    LastName: string;
    DisplayName: string;
    Email: string;
    Department: string;
    HomeNumber: string;
    CallerIDName: string;
    CallerIDNumber: string;
    OutboundDID: string;
    DirectNumber: string;
    PresentationNumber: string;
    SiteID?: string;
    ButtonProfileID?: string;
    VoicemailNotification: string;
    MobileNumber: string;
    CallWaiting: boolean;
    DoNotDisturb: boolean;
    Forwarding: boolean;
    ForwardInternalTarget: string;
    ForwardExternalTarget: string;
    MobileTwinning: boolean;
    InternalNoAnswerTimeout: number;
    InternalNoAnswerAction: string;
    InternalOnBusyAction: string;
    ExternalNoAnswerTimeout: number;
    ExternalNoAnswerAction: string;
    ExternalOnBusyAction: string;
    HandsetButtons: string;
    VoicemailEnabled: boolean;
    VoicemailPIN: string;
    VoicemailToEmail: boolean;
    SIPTransport?: string;
    SRTPMode?: string;
    Codecs?: string;
    DTMFMode?: string;
    NATMode?: string;
    CallRecording?: string;
    MaxConcurrentCalls?: number;
    RingTimeoutSeconds?: number;
    CreatedAt: string;
  }

  let extensions: ExtensionItem[] = [];
  let sites: any[] = [];
  let buttonProfiles: any[] = [];
  let telephoneNumbers: any[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // View state: 'list' or 'edit'
  let currentView: 'list' | 'edit' = 'list';
  let editingExt: ExtensionItem | null = null;
  let generatedPasswordModal: string = '';

  // Handset Key Search & Pagination
  let keySearch = '';
  let keyPageSize = 10;
  let keyCurrentPage = 1;

  // Handset Key Modal state
  let keyModalOpen = false;
  let editingKeyIndex: number = -1;
  let keyForm: HandsetKey = {
    key: 1,
    label: '',
    type: 'BLF',
    description: ''
  };

  // Main Form Data matching Screenshot 1 + Advanced SIP
  let form = {
    ID: '',
    FirstName: '',
    LastName: '',
    DisplayName: '',
    Email: '',
    MobileNumber: '',
    HomeNumber: '',
    InternalNumber: '',
    DirectNumber: '',
    PresentationNumber: '',
    SiteID: '',
    ButtonProfileID: '',
    VoicemailNotification: 'Off',
    CallWaiting: false,
    DoNotDisturb: false,
    Forwarding: false,
    ForwardInternalTarget: '',
    ForwardExternalTarget: '',
    MobileTwinning: false,
    InternalNoAnswerTimeout: 3,
    InternalNoAnswerAction: 'Not Set',
    InternalOnBusyAction: 'Unknown',
    ExternalNoAnswerTimeout: 0,
    ExternalNoAnswerAction: 'Not Set',
    ExternalOnBusyAction: 'Unknown',
    HandsetKeys: [] as HandsetKey[],
    Password: '',
    // Advanced Settings
    VoicemailEnabled: true,
    VoicemailPIN: '1234',
    VoicemailToEmail: true,
    SIPTransport: 'UDP',
    SRTPMode: 'optional',
    Codecs: 'OPUS, G.722, PCMA, PCMU',
    DTMFMode: 'rfc2833',
    NATMode: 'auto-nat',
    CallRecording: 'disabled',
    MaxConcurrentCalls: 4,
    RingTimeoutSeconds: 20
  };

  $: filteredHandsetKeys = form.HandsetKeys.filter(k =>
    !keySearch.trim() || `${k.key} ${k.label} ${k.type} ${k.description}`.toLowerCase().includes(keySearch.trim().toLowerCase())
  );
  $: totalKeyPages = Math.ceil(filteredHandsetKeys.length / keyPageSize) || 1;
  $: paginatedHandsetKeys = filteredHandsetKeys.slice((keyCurrentPage - 1) * keyPageSize, keyCurrentPage * keyPageSize);

  function generatePassword(length = 14): string {
    const charset = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$';
    let res = '';
    const array = new Uint32Array(length);
    window.crypto.getRandomValues(array);
    for (let i = 0; i < length; i++) {
      res += charset[array[i] % charset.length];
    }
    return res;
  }

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [extRes, sitesRes, bpRes, didsRes] = await Promise.all([
        api.get<{ extensions: ExtensionItem[] }>('/extensions'),
        api.get<{ sites: any[] }>('/sites').catch(() => ({ sites: [] })),
        api.get<{ button_profiles: any[] }>('/button-profiles').catch(() => ({ button_profiles: [] })),
        api.get<{ telephone_numbers: any[] }>('/telephone-numbers').catch(() => ({ telephone_numbers: [] }))
      ]);
      extensions = extRes.extensions || [];
      sites = sitesRes.sites || [];
      buttonProfiles = bpRes.button_profiles || [];
      telephoneNumbers = didsRes.telephone_numbers || [];
    } catch (e: any) {
      loadError = e.message || 'Failed to load user extensions';
    } finally {
      loading = false;
    }
  }

  onMount(loadData);

  function openCreate() {
    const autoPass = generatePassword(14);
    const nextExtNum = (1000 + extensions.length + 1).toString();
    editingExt = null;
    form = {
      ID: '',
      FirstName: '',
      LastName: '',
      DisplayName: '',
      Email: '',
      MobileNumber: '',
      HomeNumber: '',
      InternalNumber: nextExtNum,
      DirectNumber: telephoneNumbers.length > 0 ? telephoneNumbers[0].Number : '',
      PresentationNumber: telephoneNumbers.length > 0 ? telephoneNumbers[0].Number : '',
      SiteID: sites.length > 0 ? sites[0].ID : '',
      ButtonProfileID: buttonProfiles.length > 0 ? buttonProfiles[0].ID : '',
      VoicemailNotification: 'Off',
      CallWaiting: false,
      DoNotDisturb: false,
      Forwarding: false,
      ForwardInternalTarget: '',
      ForwardExternalTarget: '',
      MobileTwinning: false,
      InternalNoAnswerTimeout: 3,
      InternalNoAnswerAction: 'Not Set',
      InternalOnBusyAction: 'Unknown',
      ExternalNoAnswerTimeout: 0,
      ExternalNoAnswerAction: 'Not Set',
      ExternalOnBusyAction: 'Unknown',
      HandsetKeys: [],
      Password: autoPass,
      VoicemailEnabled: true,
      VoicemailPIN: '1234',
      VoicemailToEmail: true,
      SIPTransport: 'UDP',
      SRTPMode: 'optional',
      Codecs: 'OPUS, G.722, PCMA, PCMU',
      DTMFMode: 'rfc2833',
      NATMode: 'auto-nat',
      CallRecording: 'disabled',
      MaxConcurrentCalls: 4,
      RingTimeoutSeconds: 20
    };
    currentView = 'edit';
  }

  function openEdit(ext: ExtensionItem) {
    editingExt = ext;
    let parsedKeys: HandsetKey[] = [];
    try {
      if (ext.HandsetButtons) {
        parsedKeys = JSON.parse(ext.HandsetButtons);
      }
    } catch (e) {
      parsedKeys = [];
    }

    form = {
      ID: ext.ID,
      FirstName: ext.FirstName || '',
      LastName: ext.LastName || '',
      DisplayName: ext.DisplayName || `${ext.FirstName || ''} ${ext.LastName || ''}`.trim() || ext.ExtensionNumber,
      Email: ext.Email || '',
      MobileNumber: ext.MobileNumber || '',
      HomeNumber: ext.HomeNumber || '',
      InternalNumber: ext.ExtensionNumber,
      DirectNumber: ext.DirectNumber || ext.OutboundDID || '',
      PresentationNumber: ext.PresentationNumber || ext.CallerIDNumber || '',
      SiteID: ext.SiteID || (sites.length > 0 ? sites[0].ID : ''),
      ButtonProfileID: ext.ButtonProfileID || '',
      VoicemailNotification: ext.VoicemailNotification || 'Off',
      CallWaiting: !!ext.CallWaiting,
      DoNotDisturb: !!ext.DoNotDisturb,
      Forwarding: !!ext.Forwarding,
      ForwardInternalTarget: ext.ForwardInternalTarget || '',
      ForwardExternalTarget: ext.ForwardExternalTarget || '',
      MobileTwinning: !!ext.MobileTwinning,
      InternalNoAnswerTimeout: ext.InternalNoAnswerTimeout ?? 3,
      InternalNoAnswerAction: ext.InternalNoAnswerAction || 'Not Set',
      InternalOnBusyAction: ext.InternalOnBusyAction || 'Unknown',
      ExternalNoAnswerTimeout: ext.ExternalNoAnswerTimeout ?? 0,
      ExternalNoAnswerAction: ext.ExternalNoAnswerAction || 'Not Set',
      ExternalOnBusyAction: ext.ExternalOnBusyAction || 'Unknown',
      HandsetKeys: parsedKeys,
      Password: '',
      VoicemailEnabled: ext.VoicemailEnabled ?? true,
      VoicemailPIN: ext.VoicemailPIN || '1234',
      VoicemailToEmail: ext.VoicemailToEmail ?? true,
      SIPTransport: ext.SIPTransport || 'UDP',
      SRTPMode: ext.SRTPMode || 'optional',
      Codecs: ext.Codecs || 'OPUS, G.722, PCMA, PCMU',
      DTMFMode: ext.DTMFMode || 'rfc2833',
      NATMode: ext.NATMode || 'auto-nat',
      CallRecording: ext.CallRecording || 'disabled',
      MaxConcurrentCalls: ext.MaxConcurrentCalls || 4,
      RingTimeoutSeconds: ext.RingTimeoutSeconds || 20
    };
    currentView = 'edit';
  }

  function backToList() {
    currentView = 'list';
    editingExt = null;
  }

  // Update DisplayName automatically when First or Last Name change
  $: if (!editingExt && (form.FirstName || form.LastName)) {
    form.DisplayName = `${form.FirstName} ${form.LastName}`.trim();
  }

  async function handleSave() {
    if (!form.InternalNumber.trim()) {
      return toast.error('Internal Number is required');
    }

    saving = true;
    try {
      const payload: any = {
        ExtensionNumber: form.InternalNumber.trim(),
        FirstName: form.FirstName.trim(),
        LastName: form.LastName.trim(),
        DisplayName: form.DisplayName.trim(),
        Email: form.Email.trim(),
        MobileNumber: form.MobileNumber.trim(),
        HomeNumber: form.HomeNumber.trim(),
        DirectNumber: form.DirectNumber.trim(),
        PresentationNumber: form.PresentationNumber.trim(),
        CallerIDNumber: form.PresentationNumber.trim(),
        CallerIDName: form.DisplayName.trim(),
        OutboundDID: form.DirectNumber.trim(),
        SiteID: form.SiteID || null,
        ButtonProfileID: form.ButtonProfileID || null,
        VoicemailNotification: form.VoicemailNotification,
        CallWaiting: form.CallWaiting,
        DoNotDisturb: form.DoNotDisturb,
        Forwarding: form.Forwarding,
        ForwardInternalTarget: form.ForwardInternalTarget.trim(),
        ForwardExternalTarget: form.ForwardExternalTarget.trim(),
        MobileTwinning: form.MobileTwinning,
        InternalNoAnswerTimeout: Number(form.InternalNoAnswerTimeout),
        InternalNoAnswerAction: form.InternalNoAnswerAction,
        InternalOnBusyAction: form.InternalOnBusyAction,
        ExternalNoAnswerTimeout: Number(form.ExternalNoAnswerTimeout),
        ExternalNoAnswerAction: form.ExternalNoAnswerAction,
        ExternalOnBusyAction: form.ExternalOnBusyAction,
        HandsetButtons: JSON.stringify(form.HandsetKeys),
        VoicemailEnabled: form.VoicemailEnabled,
        VoicemailPIN: form.VoicemailPIN.trim(),
        VoicemailToEmail: form.VoicemailToEmail,
        SIPTransport: form.SIPTransport,
        SRTPMode: form.SRTPMode,
        Codecs: form.Codecs.trim(),
        DTMFMode: form.DTMFMode,
        NATMode: form.NATMode,
        CallRecording: form.CallRecording,
        MaxConcurrentCalls: Number(form.MaxConcurrentCalls),
        RingTimeoutSeconds: Number(form.RingTimeoutSeconds)
      };

      if (form.Password) {
        payload.PasswordHash = form.Password;
      }

      if (editingExt) {
        await api.put(`/extensions/${editingExt.ID}`, payload);
        toast.success(`User '${form.InternalNumber} ${form.DisplayName}' updated successfully`);
      } else {
        const created = await api.post<{ extension: ExtensionItem; generated_password?: string }>('/extensions', payload);
        toast.success(`User created with extension ${form.InternalNumber}`);
        if (created.generated_password) {
          generatedPasswordModal = created.generated_password;
        }
      }

      await loadData();
      if (!generatedPasswordModal) {
        currentView = 'list';
      }
    } catch (e: any) {
      toastError(e, 'Failed to save user extension');
    } finally {
      saving = false;
    }
  }

  async function deleteUser(ext: ExtensionItem) {
    if (!confirm(`Are you sure you want to delete user ${ext.ExtensionNumber} (${ext.DisplayName || ext.FirstName})?`)) return;
    try {
      await api.delete(`/extensions/${ext.ID}`);
      toast.success(`User ${ext.ExtensionNumber} deleted`);
      loadData();
    } catch (e: any) {
      toastError(e, 'Failed to delete user');
    }
  }

  // Handset Buttons Management
  function openAddKey() {
    editingKeyIndex = -1;
    const nextKeyNum = form.HandsetKeys.length > 0 ? Math.max(...form.HandsetKeys.map(k => k.key)) + 1 : 1;
    keyForm = {
      key: nextKeyNum,
      label: `Line ${nextKeyNum}`,
      type: 'BLF',
      description: 'Ext 1001'
    };
    keyModalOpen = true;
  }

  function editKey(index: number) {
    editingKeyIndex = index;
    keyForm = { ...form.HandsetKeys[index] };
    keyModalOpen = true;
  }

  function saveKey() {
    if (editingKeyIndex >= 0) {
      form.HandsetKeys[editingKeyIndex] = { ...keyForm };
    } else {
      form.HandsetKeys = [...form.HandsetKeys, { ...keyForm }];
    }
    keyModalOpen = false;
  }

  function removeKey(index: number) {
    form.HandsetKeys = form.HandsetKeys.filter((_, i) => i !== index);
  }

  function deleteAllKeys() {
    if (confirm('Are you sure you want to delete all programmable handset keys?')) {
      form.HandsetKeys = [];
    }
  }

  $: filteredList = extensions.filter(e =>
    !search.trim() ||
    `${e.ExtensionNumber} ${e.DisplayName} ${e.FirstName} ${e.LastName} ${e.Email} ${e.DirectNumber}`.toLowerCase().includes(search.trim().toLowerCase())
  );
</script>

<div class="p-6 bg-[#EEF2F6] min-h-screen text-slate-800">
  {#if currentView === 'list'}
    <!-- LIST VIEW -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-5 mb-5 border-b border-slate-200 gap-4">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 rounded-lg bg-sky-100 text-[#0284C7] flex items-center justify-center font-bold shadow-sm">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        </div>
        <div>
          <h1 class="text-xl font-bold text-slate-900 tracking-tight">Users & Extensions</h1>
          <p class="text-xs text-slate-500">Configure SIP endpoints, user routing policies, DIDs, and handset keys</p>
        </div>
      </div>

      <div class="flex items-center space-x-3">
        <div class="relative">
          <input
            type="text"
            bind:value={search}
            placeholder="Search users..."
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
          <span>Add User</span>
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
      <table class="w-full text-left text-xs text-slate-700">
        <thead class="bg-[#2B527A] text-white font-semibold">
          <tr>
            <th class="px-4 py-3">Internal Ext</th>
            <th class="px-4 py-3">Display Name</th>
            <th class="px-4 py-3">Direct Number (DID)</th>
            <th class="px-4 py-3">Mobile Number</th>
            <th class="px-4 py-3">Email Address</th>
            <th class="px-4 py-3">DND / Forward</th>
            <th class="px-4 py-3 text-right">Options</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          {#each filteredList as ext}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="px-4 py-3 font-bold text-slate-900 flex items-center space-x-2">
                <span class="w-2 h-2 rounded-full {ext.IsActive ? 'bg-emerald-500' : 'bg-slate-300'}"></span>
                <span>{ext.ExtensionNumber}</span>
              </td>
              <td class="px-4 py-3 font-medium text-slate-800">
                {ext.DisplayName || `${ext.FirstName || ''} ${ext.LastName || ''}`.trim() || 'Internal User'}
              </td>
              <td class="px-4 py-3 text-slate-600">
                {ext.DirectNumber || ext.OutboundDID || '—'}
              </td>
              <td class="px-4 py-3 text-slate-600">
                {ext.MobileNumber || '—'}
              </td>
              <td class="px-4 py-3 text-slate-600">
                {ext.Email || '—'}
              </td>
              <td class="px-4 py-3">
                <div class="flex items-center space-x-1.5">
                  {#if ext.DoNotDisturb}
                    <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-rose-100 text-rose-700">DND ON</span>
                  {/if}
                  {#if ext.Forwarding}
                    <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-amber-100 text-amber-700">FWD ON</span>
                  {/if}
                  {#if !ext.DoNotDisturb && !ext.Forwarding}
                    <span class="text-slate-400 text-[11px]">Normal</span>
                  {/if}
                </div>
              </td>
              <td class="px-4 py-3 text-right">
                <div class="inline-flex items-center space-x-2">
                  <button
                    type="button"
                    on:click={() => openEdit(ext)}
                    title="Edit User"
                    class="p-1 text-slate-600 hover:text-[#0284C7] hover:bg-sky-50 rounded"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                  </button>
                  <button
                    type="button"
                    on:click={() => deleteUser(ext)}
                    title="Delete User"
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
                No users found. Click "Add User" to create one.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>

  {:else}
    <!-- EDIT VIEW MATCHING SCREENSHOT 1 -->
    <!-- Top Bar with user title -->
    <div class="flex items-center space-x-2.5 pb-4 mb-4 border-b border-slate-300">
      <div class="text-[#0284C7]">
        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
        </svg>
      </div>
      <h1 class="text-sm font-bold text-slate-900 tracking-tight">
        {editingExt ? `Editing user '${form.InternalNumber} ${form.DisplayName}'` : 'Create New User'}
      </h1>
    </div>

    <form on:submit|preventDefault={handleSave} class="space-y-4">
      <!-- CARD 1: User Details (3-Column Layout) -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
          User Details
        </div>
        <div class="p-5 grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Column 1 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">First Name</label>
              <div class="relative">
                <input
                  type="text"
                  bind:value={form.FirstName}
                  placeholder="e.g. John"
                  class="w-full px-3 py-2 pr-9 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
                <button
                  type="button"
                  title="Generate ID"
                  on:click={() => { if (!form.FirstName) form.FirstName = form.InternalNumber; }}
                  class="absolute right-2 top-2 text-slate-400 hover:text-slate-600 text-xs px-1"
                >
                  •••
                </button>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Last Name</label>
              <input
                type="text"
                bind:value={form.LastName}
                placeholder="e.g. Smith"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Display Name</label>
              <input
                type="text"
                bind:value={form.DisplayName}
                placeholder="e.g. John Smith"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-medium"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Email Address</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>
                </span>
                <input
                  type="email"
                  bind:value={form.Email}
                  placeholder="user@company.com"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Mobile Number</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/></svg>
                </span>
                <input
                  type="text"
                  bind:value={form.MobileNumber}
                  placeholder="e.g. +447123456789"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Home Number</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/></svg>
                </span>
                <input
                  type="text"
                  bind:value={form.HomeNumber}
                  placeholder="Optional landline"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>
            </div>
          </div>

          <!-- Column 2 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Internal Number</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg>
                </span>
                <input
                  type="text"
                  bind:value={form.InternalNumber}
                  placeholder="e.g. 1001 or 991"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-bold"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Direct Number</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>
                </span>
                <select
                  bind:value={form.DirectNumber}
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                >
                  <option value="">Please select</option>
                  {#each telephoneNumbers as num}
                    <option value={num.Number}>{num.Number} ({num.Description || 'DID'})</option>
                  {/each}
                </select>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Presentation Number</label>
              <input
                type="text"
                bind:value={form.PresentationNumber}
                placeholder="e.g. +18005550100"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
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
                </select>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Button Profile</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/></svg>
                </span>
                <select
                  bind:value={form.ButtonProfileID}
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                >
                  <option value="">Please select</option>
                  {#each buttonProfiles as bp}
                    <option value={bp.ID}>{bp.Name}</option>
                  {/each}
                </select>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Voicemail Notification</label>
              <select
                bind:value={form.VoicemailNotification}
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              >
                <option value="Off">Off</option>
                <option value="Email">Email</option>
                <option value="SMS">SMS</option>
                <option value="Email and SMS">Email and SMS</option>
              </select>
            </div>
          </div>

          <!-- Column 3 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Call Waiting</label>
              <select
                bind:value={form.CallWaiting}
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              >
                <option value={false}>No</option>
                <option value={true}>Yes</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Do Not Disturb</label>
              <select
                bind:value={form.DoNotDisturb}
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              >
                <option value={false}>Off</option>
                <option value={true}>On</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Forwarding</label>
              <select
                bind:value={form.Forwarding}
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              >
                <option value={false}>Off</option>
                <option value={true}>On</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Forward internal calls to</label>
              <input
                type="text"
                bind:value={form.ForwardInternalTarget}
                placeholder="Extension or Target"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Forward external calls to</label>
              <input
                type="text"
                bind:value={form.ForwardExternalTarget}
                placeholder="Mobile or Landline"
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Mobile Twinning</label>
              <select
                bind:value={form.MobileTwinning}
                class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              >
                <option value={false}>Off</option>
                <option value={true}>On</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Action Buttons matching screenshot (Blue save disk icon + Back) -->
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

      <!-- MIDDLE ROW: Internal Call Settings & External Call Settings -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Internal Call Settings Card -->
        <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden flex flex-col justify-between">
          <div>
            <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
              Internal Call Settings
            </div>
            <div class="p-5 space-y-4">
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">No Answer Timeout</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                  </span>
                  <input
                    type="number"
                    bind:value={form.InternalNoAnswerTimeout}
                    min="0"
                    max="300"
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <DestinationPicker
                  label="No answer action"
                  bind:value={form.InternalNoAnswerAction}
                  on:change={(e) => form.InternalNoAnswerAction = e.detail}
                />
              </div>

              <div>
                <DestinationPicker
                  label="On Busy Action"
                  bind:value={form.InternalOnBusyAction}
                  on:change={(e) => form.InternalOnBusyAction = e.detail}
                />
              </div>
            </div>
          </div>

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

        <!-- External Call Settings Card -->
        <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden flex flex-col justify-between">
          <div>
            <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
              External Call Settings
            </div>
            <div class="p-5 space-y-4">
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">No Answer Timeout</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                  </span>
                  <input
                    type="number"
                    bind:value={form.ExternalNoAnswerTimeout}
                    min="0"
                    max="300"
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <DestinationPicker
                  label="No answer action"
                  bind:value={form.ExternalNoAnswerAction}
                  on:change={(e) => form.ExternalNoAnswerAction = e.detail}
                />
              </div>

              <div>
                <DestinationPicker
                  label="On Busy Action"
                  bind:value={form.ExternalOnBusyAction}
                  on:change={(e) => form.ExternalOnBusyAction = e.detail}
                />
              </div>
            </div>
          </div>

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
      </div>

      <!-- CARD: Advanced Telephony, SIP & Media Settings -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800 flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <svg class="w-4 h-4 text-[#0284C7]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span>Advanced SIP, Codecs & Voicemail Settings</span>
          </div>
          <span class="text-[11px] text-slate-400 font-normal">Industry-standard FreeSWITCH & 3CX compatible parameters</span>
        </div>
        <div class="p-5 grid grid-cols-1 md:grid-cols-3 gap-5">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">SIP Transport Protocol</label>
            <select bind:value={form.SIPTransport} class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-medium">
              <option value="UDP">UDP (Standard Telephony)</option>
              <option value="TCP">TCP (Reliable Stream)</option>
              <option value="TLS">TLS (Encrypted SIP Signaling)</option>
              <option value="WSS">WSS (WebRTC Browser / Softphone)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">SRTP Media Encryption</label>
            <select bind:value={form.SRTPMode} class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none">
              <option value="disabled">Disabled (Plain RTP)</option>
              <option value="optional">Optional (Negotiate SRTP if supported)</option>
              <option value="mandatory">Mandatory (Require SRTP Crypto)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">DTMF Signaling Mode</label>
            <select bind:value={form.DTMFMode} class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none">
              <option value="rfc2833">RFC 2833 (RTP Event / Standard)</option>
              <option value="info">SIP INFO</option>
              <option value="inband">Inband Audio Tones</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">NAT Traversal Mode</label>
            <select bind:value={form.NATMode} class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none">
              <option value="auto-nat">Auto-NAT (Comedia & Force rport)</option>
              <option value="force-rport">Force rport Only</option>
              <option value="never">Never (Local LAN only)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Call Recording Policy</label>
            <select bind:value={form.CallRecording} class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none">
              <option value="disabled">Disabled (Do Not Record)</option>
              <option value="always">Always Record All Calls</option>
              <option value="external_only">Record External Calls Only</option>
              <option value="inbound_only">Record Inbound Calls Only</option>
              <option value="on_demand">On-Demand via Star Code (*1)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Concurrent Channel Limit</label>
            <select bind:value={form.MaxConcurrentCalls} class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-bold">
              <option value={1}>1 Line (Call Waiting Off)</option>
              <option value={2}>2 Concurrent Calls</option>
              <option value={4}>4 Concurrent Calls (Standard)</option>
              <option value={8}>8 Concurrent Calls</option>
              <option value={16}>16 Concurrent Calls (Operator Desk)</option>
            </select>
          </div>

          <div class="md:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1">Permitted Audio Codecs (In Priority Order)</label>
            <input type="text" bind:value={form.Codecs} placeholder="OPUS, G.722, PCMA, PCMU" class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-mono" />
            <p class="text-[11px] text-slate-400 mt-1">FreeSWITCH prioritizes codecs left-to-right. OPUS provides HD wideband audio with adaptive bitrates.</p>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Ring Timeout (Seconds)</label>
            <input type="number" min="5" max="120" bind:value={form.RingTimeoutSeconds} class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-bold" />
          </div>

          <!-- Voicemail Sub-section -->
          <div class="md:col-span-3 pt-3 border-t border-slate-100 grid grid-cols-1 sm:grid-cols-3 gap-4 bg-slate-50 p-3 rounded-lg">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Voicemail PIN</label>
              <input type="password" bind:value={form.VoicemailPIN} placeholder="1234" class="w-full px-3 py-1.5 bg-white border border-slate-300 rounded text-xs font-mono" />
            </div>
            <div class="flex items-center space-x-2 pt-4">
              <input type="checkbox" id="ext-vm-en" bind:checked={form.VoicemailEnabled} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
              <label for="ext-vm-en" class="text-xs font-semibold text-slate-800">Voicemail Mailbox Active</label>
            </div>
            <div class="flex items-center space-x-2 pt-4">
              <input type="checkbox" id="ext-vm-email" bind:checked={form.VoicemailToEmail} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
              <label for="ext-vm-email" class="text-xs font-semibold text-slate-800">Email WAV Audio to User</label>
            </div>
          </div>
        </div>
      </div>

      <!-- CARD 3: Handset Buttons (Matching Screenshot 1 Bottom Card) -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
          Handset Buttons
        </div>
        <div class="p-5">
          <!-- Top Button Actions -->
          <div class="flex items-center space-x-2 mb-4">
            <button
              type="button"
              on:click={openAddKey}
              class="px-3 py-1.5 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 text-xs font-semibold rounded shadow-sm transition-colors flex items-center space-x-1"
            >
              <span>+ Add</span>
            </button>
            <button
              type="button"
              on:click={deleteAllKeys}
              class="px-3 py-1.5 bg-[#2B527A] hover:bg-[#1E3A8A] text-white text-xs font-semibold rounded shadow-sm transition-colors"
            >
              Delete All
            </button>
          </div>

          <!-- Table Controls -->
          <div class="flex items-center justify-between text-xs text-slate-500 mb-3">
            <div class="flex items-center space-x-1.5">
              <span>Show</span>
              <select bind:value={keyPageSize} on:change={() => keyCurrentPage = 1} class="border border-slate-300 rounded px-2 py-1 text-xs bg-white">
                <option value={10}>10</option>
                <option value={25}>25</option>
                <option value={50}>50</option>
              </select>
              <span>entries</span>
            </div>
            <div class="flex items-center space-x-1.5">
              <span>Search:</span>
              <input type="text" bind:value={keySearch} on:input={() => keyCurrentPage = 1} placeholder="Filter keys..." class="border border-slate-300 rounded px-2 py-1 text-xs bg-white w-40" />
            </div>
          </div>

          <!-- Table with Navy Header matching screenshot -->
          <div class="border border-slate-200 rounded overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-[#2B527A] text-white font-semibold">
                <tr>
                  <th class="px-4 py-2.5">Key</th>
                  <th class="px-4 py-2.5">Label</th>
                  <th class="px-4 py-2.5">Key Type</th>
                  <th class="px-4 py-2.5">Key Description</th>
                  <th class="px-4 py-2.5 text-right">Options</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 bg-white">
                {#each paginatedHandsetKeys as k, idx}
                  <tr class="hover:bg-slate-50 transition-colors">
                    <td class="px-4 py-2.5 font-bold text-slate-900">{k.key}</td>
                    <td class="px-4 py-2.5 text-slate-800">{k.label}</td>
                    <td class="px-4 py-2.5">
                      <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-sky-100 text-sky-800">{k.type}</span>
                    </td>
                    <td class="px-4 py-2.5 text-slate-600">{k.description || '—'}</td>
                    <td class="px-4 py-2.5 text-right">
                      <div class="inline-flex items-center space-x-2">
                        <button
                          type="button"
                          on:click={() => editKey((keyCurrentPage - 1) * keyPageSize + idx)}
                          class="p-1 text-slate-500 hover:text-[#0284C7]"
                          title="Edit Key"
                        >
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                        </button>
                        <button
                          type="button"
                          on:click={() => removeKey((keyCurrentPage - 1) * keyPageSize + idx)}
                          class="p-1 text-slate-500 hover:text-rose-600"
                          title="Remove Key"
                        >
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
                {#if filteredHandsetKeys.length === 0}
                  <tr>
                    <td colspan="5" class="px-4 py-6 text-center text-slate-400 bg-slate-50/50">
                      {keySearch ? 'No matching keys found' : 'No data available in table'}
                    </td>
                  </tr>
                {/if}
              </tbody>
            </table>
          </div>

          <div class="flex items-center justify-between text-xs text-slate-500 mt-3">
            <div>Showing {filteredHandsetKeys.length > 0 ? (keyCurrentPage - 1) * keyPageSize + 1 : 0} to {Math.min(keyCurrentPage * keyPageSize, filteredHandsetKeys.length)} of {filteredHandsetKeys.length} entries</div>
            <div class="flex space-x-1">
              <button
                type="button"
                on:click={() => { if (keyCurrentPage > 1) keyCurrentPage--; }}
                class="px-2 py-1 border border-slate-300 rounded {keyCurrentPage > 1 ? 'text-slate-700 hover:bg-slate-50 cursor-pointer' : 'text-slate-300 cursor-not-allowed'}"
                disabled={keyCurrentPage <= 1}
              >
                Previous
              </button>
              <button
                type="button"
                on:click={() => { if (keyCurrentPage < totalKeyPages) keyCurrentPage++; }}
                class="px-2 py-1 border border-slate-300 rounded {keyCurrentPage < totalKeyPages ? 'text-slate-700 hover:bg-slate-50 cursor-pointer' : 'text-slate-300 cursor-not-allowed'}"
                disabled={keyCurrentPage >= totalKeyPages}
              >
                Next
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  {/if}
</div>

<!-- Modal: Handset Key Editor -->
{#if keyModalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-[2px]">
    <div class="bg-white rounded-lg shadow-xl border border-slate-200 w-full max-w-md overflow-hidden">
      <div class="px-5 py-3.5 bg-[#2B527A] text-white font-bold text-xs flex justify-between items-center">
        <span>{editingKeyIndex >= 0 ? 'Edit Handset Key' : 'Add Handset Key'}</span>
        <button type="button" on:click={() => keyModalOpen = false} class="text-white/80 hover:text-white">✕</button>
      </div>
      <div class="p-5 space-y-3.5">
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Key Number (1-30)</label>
          <input
            type="number"
            bind:value={keyForm.key}
            min="1"
            max="30"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Key Label</label>
          <input
            type="text"
            bind:value={keyForm.label}
            placeholder="e.g. Line 1, Front Desk, Voicemail"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Key Type</label>
          <select
            bind:value={keyForm.type}
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          >
            <option value="Line">Line Key</option>
            <option value="BLF">BLF (Busy Lamp Field)</option>
            <option value="SpeedDial">Speed Dial</option>
            <option value="Park">Call Park</option>
            <option value="Intercom">Intercom</option>
            <option value="Voicemail">Voicemail</option>
            <option value="DND">Do Not Disturb</option>
          </select>
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Key Description / Target Extension</label>
          <input
            type="text"
            bind:value={keyForm.description}
            placeholder="e.g. 1001, *97, *8"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
        </div>
      </div>
      <div class="px-5 py-3 bg-slate-50 border-t border-slate-200 flex justify-end space-x-2">
        <button
          type="button"
          on:click={() => keyModalOpen = false}
          class="px-3 py-1.5 bg-white border border-slate-300 text-slate-700 font-semibold text-xs rounded hover:bg-slate-50"
        >
          Cancel
        </button>
        <button
          type="button"
          on:click={saveKey}
          class="px-4 py-1.5 bg-[#0284C7] hover:bg-[#0369A1] text-white font-bold text-xs rounded shadow-sm"
        >
          Save Key
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal: Generated Password Display on User Creation -->
{#if generatedPasswordModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-[2px]">
    <div class="bg-white rounded-lg shadow-2xl border border-slate-200 w-full max-w-md p-6 text-center">
      <div class="w-12 h-12 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto mb-3">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <h3 class="text-sm font-bold text-slate-900 mb-1">User Extension Created!</h3>
      <p class="text-xs text-slate-500 mb-4">A high-entropy SIP password was automatically generated for this handset endpoint:</p>
      <div class="p-3 bg-slate-100 border border-slate-300 rounded-lg text-slate-900 font-mono text-sm tracking-wider font-bold mb-4 select-all flex items-center justify-between">
        <span>{generatedPasswordModal}</span>
        <button
          type="button"
          on:click={async () => {
            await navigator.clipboard.writeText(generatedPasswordModal);
            toast.success('Password copied to clipboard!');
          }}
          class="text-xs text-[#0284C7] hover:underline font-sans ml-2"
        >
          Copy
        </button>
      </div>
      <button
        type="button"
        on:click={() => { generatedPasswordModal = ''; currentView = 'list'; }}
        class="w-full py-2 bg-[#0284C7] hover:bg-[#0369A1] text-white font-bold text-xs rounded shadow-sm transition-colors"
      >
        Done & Return to Users
      </button>
    </div>
  </div>
{/if}
