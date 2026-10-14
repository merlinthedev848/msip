<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import DestinationPicker from '../components/DestinationPicker.svelte';

  interface Choice {
    ID?: string;
    Digits: string;
    Name: string;
    Action: string;
    Destination: string;
  }

  interface IVRItem {
    ID: string;
    Name: string;
    Extension: string;
    Greeting: string;
    Timeout: number;
    MaxFailures: number;
    TimeoutDestType: string;
    TimeoutDestTarget: string;
    InvalidDestType: string;
    InvalidDestTarget: string;
    AllowDirectDial: boolean;
    PinCode: string;
    SiteID?: string;
    AnnounceOnly: boolean;
    TimeoutDestination: string;
    SoundFile: string;
    Choices: Choice[];
    CreatedAt: string;
  }

  let ivrs: IVRItem[] = [];
  let sites: any[] = [];
  let audioPrompts: any[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // View mode: 'list' or 'edit'
  let currentView: 'list' | 'edit' = 'list';
  let editingMenu: IVRItem | null = null;

  // Option Modal
  let optionModalOpen = false;
  let editingOptionIndex: number = -1;
  let optionForm: Choice = {
    Digits: '1',
    Name: 'Customer Support',
    Action: 'transfer',
    Destination: 'End Call'
  };

  // Main Form Data matching Screenshot 2
  let form = {
    ID: '',
    Name: '',
    Extension: '7000',
    PinCode: '',
    SiteID: '',
    AnnounceOnly: false,
    MenuTimeout: 10,
    TimeoutDestination: 'End Call',
    SoundFile: '',
    Choices: [] as Choice[]
  };

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [ivrRes, siteRes, audioRes] = await Promise.all([
        api.get<{ ivrs: IVRItem[] }>('/ivrs'),
        api.get<{ sites: any[] }>('/sites').catch(() => ({ sites: [] })),
        api.get<{ prompts: any[]; audio_prompts: any[] }>('/audio-prompts').catch(() => ({ prompts: [], audio_prompts: [] }))
      ]);
      ivrs = ivrRes.ivrs || [];
      sites = siteRes.sites || [];
      audioPrompts = (audioRes.prompts || audioRes.audio_prompts || []);
    } catch (e: any) {
      loadError = e.message || 'Failed to load IVR menus';
    } finally {
      loading = false;
    }
  }

  onMount(loadData);

  function openCreate() {
    editingMenu = null;
    form = {
      ID: '',
      Name: '',
      Extension: (7000 + ivrs.length + 1).toString(),
      PinCode: Math.floor(1000 + Math.random() * 9000).toString(),
      SiteID: sites.length > 0 ? sites[0].ID : '',
      AnnounceOnly: false,
      MenuTimeout: 10,
      TimeoutDestination: 'End Call',
      SoundFile: audioPrompts.length > 0 ? audioPrompts[0].Name : '',
      Choices: []
    };
    currentView = 'edit';
  }

  function openEdit(menu: IVRItem) {
    editingMenu = menu;
    form = {
      ID: menu.ID,
      Name: menu.Name || '',
      Extension: menu.Extension || '7000',
      PinCode: menu.PinCode || '',
      SiteID: menu.SiteID || (sites.length > 0 ? sites[0].ID : ''),
      AnnounceOnly: !!menu.AnnounceOnly,
      MenuTimeout: menu.Timeout > 0 ? Math.round(menu.Timeout / 1000) : 10,
      TimeoutDestination: menu.TimeoutDestination || 'End Call',
      SoundFile: menu.SoundFile || menu.Greeting || (audioPrompts.length > 0 ? audioPrompts[0].Name : 'Default Prompt'),
      Choices: (menu.Choices || []).map(c => ({
        Digits: c.Digits,
        Name: c.Name || `Option ${c.Digits}`,
        Action: c.Action || 'transfer',
        Destination: c.Destination
      }))
    };
    currentView = 'edit';
  }

  function backToList() {
    currentView = 'list';
    editingMenu = null;
  }

  async function handleSave() {
    if (!form.Name.trim()) {
      return toast.error('Menu Name is required');
    }

    saving = true;
    try {
      const payload: any = {
        Name: form.Name.trim(),
        Extension: form.Extension.trim() || '7000',
        Greeting: form.SoundFile || 'ivr/welcome.wav',
        SoundFile: form.SoundFile,
        Timeout: form.MenuTimeout * 1000,
        MaxFailures: 3,
        TimeoutDestination: form.TimeoutDestination,
        TimeoutDestType: 'DESTINATION',
        TimeoutDestTarget: form.TimeoutDestination,
        InvalidDestType: 'DESTINATION',
        InvalidDestTarget: form.TimeoutDestination,
        AllowDirectDial: true,
        PinCode: form.PinCode.trim(),
        SiteID: form.SiteID || null,
        AnnounceOnly: form.AnnounceOnly,
        Choices: form.Choices.map(c => ({
          Digits: c.Digits,
          Name: c.Name,
          Action: c.Action || 'transfer',
          Destination: c.Destination
        }))
      };

      if (editingMenu) {
        await api.put(`/ivrs/${editingMenu.ID}`, payload);
        toast.success(`Menu '${form.Name}' updated`);
      } else {
        await api.post('/ivrs', payload);
        toast.success(`Menu '${form.Name}' created`);
      }

      await loadData();
      currentView = 'list';
    } catch (e: any) {
      toastError(e, 'Failed to save IVR menu');
    } finally {
      saving = false;
    }
  }

  async function deleteMenu(menu: IVRItem) {
    if (!confirm(`Are you sure you want to delete menu '${menu.Name}'?`)) return;
    try {
      await api.delete(`/ivrs/${menu.ID}`);
      toast.success(`Menu '${menu.Name}' deleted`);
      loadData();
    } catch (e: any) {
      toastError(e, 'Failed to delete menu');
    }
  }

  // Options CRUD
  function openAddOption() {
    editingOptionIndex = -1;
    const existingDigits = form.Choices.map(c => parseInt(c.Digits)).filter(n => !isNaN(n));
    const nextDigit = existingDigits.length > 0 ? Math.max(...existingDigits) + 1 : 1;
    optionForm = {
      Digits: nextDigit <= 9 ? nextDigit.toString() : '*',
      Name: `Option ${nextDigit}`,
      Action: 'transfer',
      Destination: 'End Call'
    };
    optionModalOpen = true;
  }

  function editOption(index: number) {
    editingOptionIndex = index;
    optionForm = { ...form.Choices[index] };
    optionModalOpen = true;
  }

  function saveOption() {
    if (!optionForm.Digits.trim()) return toast.error('Digit is required');
    if (!optionForm.Name.trim()) optionForm.Name = `Option ${optionForm.Digits}`;

    if (editingOptionIndex >= 0) {
      form.Choices[editingOptionIndex] = { ...optionForm };
    } else {
      form.Choices = [...form.Choices, { ...optionForm }];
    }
    optionModalOpen = false;
  }

  function removeOption(index: number) {
    form.Choices = form.Choices.filter((_, i) => i !== index);
  }

  $: filteredList = ivrs.filter(i =>
    !search.trim() ||
    `${i.Name} ${i.Extension} ${i.PinCode}`.toLowerCase().includes(search.trim().toLowerCase())
  );
</script>

<div class="p-6 bg-[#EEF2F6] min-h-screen text-slate-800">
  {#if currentView === 'list'}
    <!-- LIST VIEW -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-5 mb-5 border-b border-slate-200 gap-4">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 rounded-lg bg-sky-100 text-[#0284C7] flex items-center justify-center font-bold shadow-sm">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
          </svg>
        </div>
        <div>
          <h1 class="text-xl font-bold text-slate-900 tracking-tight">Auto Attendants (IVR Menus)</h1>
          <p class="text-xs text-slate-500">Configure interactive digital receptionists, audio greetings, and DTMF key actions</p>
        </div>
      </div>

      <div class="flex items-center space-x-3">
        <div class="relative">
          <input
            type="text"
            bind:value={search}
            placeholder="Search menus..."
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
          <span>Add Menu</span>
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
      <table class="w-full text-left text-xs text-slate-700">
        <thead class="bg-[#2B527A] text-white font-semibold">
          <tr>
            <th class="px-4 py-3">Menu Name</th>
            <th class="px-4 py-3">Internal Ext</th>
            <th class="px-4 py-3">Pin Code</th>
            <th class="px-4 py-3">Sound Greeting</th>
            <th class="px-4 py-3">Timeout (s)</th>
            <th class="px-4 py-3">Options Configured</th>
            <th class="px-4 py-3 text-right">Options</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          {#each filteredList as item}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="px-4 py-3 font-bold text-slate-900">{item.Name}</td>
              <td class="px-4 py-3 text-slate-600 font-mono font-medium">{item.Extension}</td>
              <td class="px-4 py-3 text-slate-600 font-mono">{item.PinCode || '—'}</td>
              <td class="px-4 py-3 text-slate-600 flex items-center space-x-1.5">
                <span class="text-[#0284C7]">▶</span>
                <span class="truncate max-w-xs">{item.SoundFile || item.Greeting || 'Default'}</span>
              </td>
              <td class="px-4 py-3 text-slate-600">{Math.round((item.Timeout || 10000) / 1000)}s</td>
              <td class="px-4 py-3">
                <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-sky-100 text-sky-800">
                  {(item.Choices || []).length} Options
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <div class="inline-flex items-center space-x-2">
                  <button
                    type="button"
                    on:click={() => openEdit(item)}
                    title="Edit Menu"
                    class="p-1 text-slate-600 hover:text-[#0284C7] hover:bg-sky-50 rounded"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                  </button>
                  <button
                    type="button"
                    on:click={() => deleteMenu(item)}
                    title="Delete Menu"
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
                No IVR menus found. Click "Add Menu" to create one.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>

  {:else}
    <!-- EDIT VIEW MATCHING SCREENSHOT 2 -->
    <!-- Top Bar with menu title -->
    <div class="flex items-center space-x-2.5 pb-4 mb-4 border-b border-slate-300">
      <div class="text-[#0284C7] flex items-center space-x-1 font-bold text-xs">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
        </svg>
      </div>
      <h1 class="text-sm font-bold text-slate-900 tracking-tight">
        {editingMenu ? `Editing menu '${form.Name}'` : 'Create New Menu'}
      </h1>
    </div>

    <form on:submit|preventDefault={handleSave} class="space-y-4">
      <!-- CARD 1: Menu Details (3-Column Layout) -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
          Menu Details
        </div>
        <div class="p-5 grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Column 1 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Name</label>
              <div class="relative">
                <input
                  type="text"
                  bind:value={form.Name}
                  placeholder="e.g. Main Auto Attendant"
                  class="w-full px-3 py-2 pr-9 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-medium"
                />
                <button
                  type="button"
                  on:click={() => { if (!form.Name) form.Name = 'Main Auto Attendant'; else form.Name += ' (Copy)'; }}
                  title="Generate or Auto-complete"
                  class="absolute right-2 top-2 text-slate-400 hover:text-slate-600 text-xs px-1 cursor-pointer transition-colors"
                >
                  •••
                </button>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Pin Code</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/></svg>
                </span>
                <input
                  type="text"
                  bind:value={form.PinCode}
                  placeholder="e.g. 9753"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-mono"
                />
              </div>
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
          </div>

          <!-- Column 2 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Announce Only</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"/></svg>
                </span>
                <select
                  bind:value={form.AnnounceOnly}
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                >
                  <option value={false}>No</option>
                  <option value={true}>Yes</option>
                </select>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Menu Timeout</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-slate-400">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                </span>
                <input
                  type="number"
                  bind:value={form.MenuTimeout}
                  min="1"
                  max="60"
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>
            </div>

            <div>
              <DestinationPicker
                label="Timeout Destination"
                bind:value={form.TimeoutDestination}
                on:change={(e) => form.TimeoutDestination = e.detail}
              />
            </div>
          </div>

          <!-- Column 3 -->
          <div class="space-y-3.5">
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Sound File</label>
              <div class="relative">
                <span class="absolute left-3 top-2.5 text-[#0284C7]">
                  <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd"/></svg>
                </span>
                <select
                  bind:value={form.SoundFile}
                  class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                >
                  <option value="">Please select</option>
                  {#each audioPrompts as prompt}
                    <option value={prompt.Name}>{prompt.Name}</option>
                  {/each}
                  {#if audioPrompts.length === 0}
                    <option value="" disabled>No audio prompts uploaded yet</option>
                  {/if}
                </select>
              </div>
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

      <!-- CARD 2: Options (Matching Screenshot 2 Bottom Card) -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
          Options
        </div>
        <div class="p-5">
          <!-- Top + Add Button -->
          <div class="mb-4">
            <button
              type="button"
              on:click={openAddOption}
              class="px-3 py-1.5 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 text-xs font-semibold rounded shadow-sm transition-colors flex items-center space-x-1"
            >
              <span>+ Add</span>
            </button>
          </div>

          <!-- Table with Navy Header matching Screenshot 2 -->
          <div class="border border-slate-200 rounded overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-[#2B527A] text-white font-semibold">
                <tr>
                  <th class="px-4 py-2.5">Option</th>
                  <th class="px-4 py-2.5">Name</th>
                  <th class="px-4 py-2.5">Destination</th>
                  <th class="px-4 py-2.5 text-right">Options</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 bg-white">
                {#each form.Choices as choice, idx}
                  <tr class="hover:bg-slate-50 transition-colors">
                    <td class="px-4 py-2.5 font-bold text-slate-900 font-mono text-xs">{choice.Digits}</td>
                    <td class="px-4 py-2.5 text-slate-800 font-medium">{choice.Name}</td>
                    <td class="px-4 py-2.5 text-slate-600 flex items-center space-x-1.5">
                      <span class="text-slate-400">↳</span>
                      <span>{choice.Destination}</span>
                    </td>
                    <td class="px-4 py-2.5 text-right">
                      <div class="inline-flex items-center space-x-2">
                        <button
                          type="button"
                          on:click={() => editOption(idx)}
                          class="p-1 text-slate-500 hover:text-[#0284C7]"
                          title="Edit Option"
                        >
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                        </button>
                        <button
                          type="button"
                          on:click={() => removeOption(idx)}
                          class="p-1 text-slate-500 hover:text-rose-600"
                          title="Remove Option"
                        >
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
                {#if form.Choices.length === 0}
                  <tr>
                    <td colspan="4" class="px-4 py-6 text-center text-slate-400 bg-slate-50/50">
                      No options configured. Click "+ Add" to add DTMF choices (e.g. Press 1 for Sales).
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

<!-- Modal: IVR Option Editor -->
{#if optionModalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-[2px]">
    <div class="bg-white rounded-lg shadow-xl border border-slate-200 w-full max-w-md overflow-hidden">
      <div class="px-5 py-3.5 bg-[#2B527A] text-white font-bold text-xs flex justify-between items-center">
        <span>{editingOptionIndex >= 0 ? 'Edit Menu Option' : 'Add Menu Option'}</span>
        <button type="button" on:click={() => optionModalOpen = false} class="text-white/80 hover:text-white">✕</button>
      </div>
      <div class="p-5 space-y-4">
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Key / Digit (0-9, *, #)</label>
          <input
            type="text"
            bind:value={optionForm.Digits}
            placeholder="e.g. 1"
            maxlength="2"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-bold"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Option Name / Department</label>
          <input
            type="text"
            bind:value={optionForm.Name}
            placeholder="e.g. Accounts & Billing"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
        </div>
        <div>
          <DestinationPicker
            label="Destination"
            bind:value={optionForm.Destination}
            on:change={(e) => optionForm.Destination = e.detail}
          />
        </div>
      </div>
      <div class="px-5 py-3 bg-slate-50 border-t border-slate-200 flex justify-end space-x-2">
        <button
          type="button"
          on:click={() => optionModalOpen = false}
          class="px-3 py-1.5 bg-white border border-slate-300 text-slate-700 font-semibold text-xs rounded hover:bg-slate-50"
        >
          Cancel
        </button>
        <button
          type="button"
          on:click={saveOption}
          class="px-4 py-1.5 bg-[#0284C7] hover:bg-[#0369A1] text-white font-bold text-xs rounded shadow-sm"
        >
          Save Option
        </button>
      </div>
    </div>
  </div>
{/if}
