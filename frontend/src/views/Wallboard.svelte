<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { api, wsUrl } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import DestinationPicker from '../components/DestinationPicker.svelte';

  interface QueueMemberItem {
    ID?: string;
    Member: string;
    Extension: string;
    Penalty: number;
    Disabled: boolean;
  }

  interface QueueItem {
    ID: string;
    Name: string;
    Extension: string;
    MaxCapacity: number;
    RingStrategy: string;
    MemberRingTimeout: number;
    MemberRetryDelay: number;
    MemberWrapUpDelay: number;
    InternalNumber: string;
    PlayMusicWhileWaiting: boolean;
    Music: string;
    CallerNamePrefix: string;
    RingCadence: string;
    AnnounceWaitTime: boolean;
    AnnouncePosition: boolean;
    AnnouncementInterval: number;
    NoAnswerTimeout: number;
    NoAnswerDestination: string;
    ParkEnabled: boolean;
    ParkTelephoneNumber: string;
    ParkResumeWhenAtPosition: number;
    ParkAllowedMinPosition: number;
    ParkResumeRetryIntervalMin: number;
    ParkResumeMaxAttempts: number;
    ParkResumeMenuTimeoutSec: number;
    ParkResumeDialTimeoutSec: number;
    ParkConfirmation: string;
    Members: QueueMemberItem[];
    CreatedAt: string;
  }

  interface WallboardMetrics {
    active_calls: number;
    agents_online: number;
    agents_busy: number;
    queue_length: number;
    avg_wait_time: number;
    sla_percentage: number;
  }

  interface QueueStat {
    name: string;
    waiting: number;
    agents: number;
    sla: number;
  }

  // Top Tabs: 'queues' or 'wallboard'
  let activeTab: 'queues' | 'wallboard' = 'queues';

  // Queue View mode: 'list' or 'edit'
  let currentView: 'list' | 'edit' = 'list';
  let editingQueue: QueueItem | null = null;

  let queuesList: QueueItem[] = [];
  let availableExts: any[] = [];
  let telephoneNumbers: any[] = [];
  let audioPrompts: any[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // Member Modal State
  let memberModalOpen = false;
  let editingMemberIndex = -1;
  let memberForm = {
    Member: '',
    Extension: '',
    Penalty: 0,
    Disabled: false
  };

  // Main Form Data matching Screenshot 3
  let form = {
    ID: '',
    Name: '',
    Extension: '8001',
    MaxCapacity: 10,
    RingStrategy: 'Ring All',
    MemberRingTimeout: 20,
    MemberRetryDelay: 10,
    MemberWrapUpDelay: 0,
    InternalNumber: '991',
    PlayMusicWhileWaiting: true,
    Music: 'Default',
    CallerNamePrefix: '',
    RingCadence: 'Default',
    AnnounceWaitTime: false,
    AnnouncePosition: true,
    AnnouncementInterval: 60,
    NoAnswerTimeout: 20,
    NoAnswerDestination: 'Voicemail',
    ParkEnabled: false,
    ParkTelephoneNumber: '',
    ParkResumeWhenAtPosition: 2,
    ParkAllowedMinPosition: 2,
    ParkResumeRetryIntervalMin: 1,
    ParkResumeMaxAttempts: 3,
    ParkResumeMenuTimeoutSec: 15,
    ParkResumeDialTimeoutSec: 20,
    ParkConfirmation: '',
    Members: [] as QueueMemberItem[]
  };

  // Live Wallboard Telemetry
  let metrics: WallboardMetrics = {
    active_calls: 0,
    agents_online: 5,
    agents_busy: 1,
    queue_length: 0,
    avg_wait_time: 14,
    sla_percentage: 98.4
  };
  let liveQueueStats: QueueStat[] = [];
  let wsConnected = false;
  let ws: WebSocket | null = null;
  let reconnectTimer: any;

  function connectWS() {
    try {
      const token = localStorage.getItem('pbx_token');
      const streamUrl = wsUrl('/wallboard/stream') + (token ? `?token=${encodeURIComponent(token)}` : '');
      ws = new WebSocket(streamUrl);
      ws.onopen = () => { wsConnected = true; };
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.metrics) metrics = data.metrics;
          if (data.queues) liveQueueStats = data.queues;
        } catch (e) {}
      };
      ws.onclose = () => {
        wsConnected = false;
        reconnectTimer = setTimeout(connectWS, 4000);
      };
      ws.onerror = () => { wsConnected = false; };
    } catch (e) {
      wsConnected = false;
    }
  }

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [qRes, extRes, telRes, promptRes] = await Promise.all([
        api.get<{ queues: QueueItem[] }>('/queues'),
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] })),
        api.get<{ telephone_numbers: any[] }>('/telephone-numbers').catch(() => ({ telephone_numbers: [] })),
        api.get<{ prompts: any[]; audio_prompts: any[] }>('/audio-prompts').catch(() => ({ prompts: [], audio_prompts: [] }))
      ]);
      queuesList = qRes.queues || [];
      availableExts = extRes.extensions || [];
      telephoneNumbers = telRes.telephone_numbers || [];
      audioPrompts = (promptRes.prompts || promptRes.audio_prompts || []);
    } catch (e: any) {
      loadError = e.message || 'Failed to load queues';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    loadData();
    connectWS();
  });

  onDestroy(() => {
    if (ws) ws.close();
    if (reconnectTimer) clearTimeout(reconnectTimer);
  });

  function openCreate() {
    editingQueue = null;
    form = {
      ID: '',
      Name: '',
      Extension: (8000 + queuesList.length + 1).toString(),
      MaxCapacity: 10,
      RingStrategy: 'Ring All',
      MemberRingTimeout: 20,
      MemberRetryDelay: 10,
      MemberWrapUpDelay: 0,
      InternalNumber: (990 + queuesList.length + 1).toString(),
      PlayMusicWhileWaiting: true,
      Music: audioPrompts.length > 0 ? audioPrompts[0].Name : 'Default',
      CallerNamePrefix: '',
      RingCadence: 'Default',
      AnnounceWaitTime: false,
      AnnouncePosition: true,
      AnnouncementInterval: 60,
      NoAnswerTimeout: 20,
      NoAnswerDestination: 'Not Set',
      ParkEnabled: false,
      ParkTelephoneNumber: telephoneNumbers.length > 0 ? telephoneNumbers[0].Number : '',
      ParkResumeWhenAtPosition: 2,
      ParkAllowedMinPosition: 2,
      ParkResumeRetryIntervalMin: 1,
      ParkResumeMaxAttempts: 3,
      ParkResumeMenuTimeoutSec: 15,
      ParkResumeDialTimeoutSec: 20,
      ParkConfirmation: '',
      Members: []
    };
    currentView = 'edit';
  }

  function openEdit(q: QueueItem) {
    editingQueue = q;
    form = {
      ID: q.ID,
      Name: q.Name || '',
      Extension: q.Extension || '8001',
      MaxCapacity: q.MaxCapacity ?? 10,
      RingStrategy: q.RingStrategy || 'Ring All',
      MemberRingTimeout: q.MemberRingTimeout ?? 20,
      MemberRetryDelay: q.MemberRetryDelay ?? 10,
      MemberWrapUpDelay: q.MemberWrapUpDelay ?? 0,
      InternalNumber: q.InternalNumber || '991',
      PlayMusicWhileWaiting: q.PlayMusicWhileWaiting !== false,
      Music: q.Music || 'Default',
      CallerNamePrefix: q.CallerNamePrefix || '',
      RingCadence: q.RingCadence || 'Default',
      AnnounceWaitTime: !!q.AnnounceWaitTime,
      AnnouncePosition: q.AnnouncePosition !== false,
      AnnouncementInterval: q.AnnouncementInterval ?? 60,
      NoAnswerTimeout: q.NoAnswerTimeout ?? 20,
      NoAnswerDestination: q.NoAnswerDestination || 'Voicemail',
      ParkEnabled: !!q.ParkEnabled,
      ParkTelephoneNumber: q.ParkTelephoneNumber || '',
      ParkResumeWhenAtPosition: q.ParkResumeWhenAtPosition ?? 2,
      ParkAllowedMinPosition: q.ParkAllowedMinPosition ?? 2,
      ParkResumeRetryIntervalMin: q.ParkResumeRetryIntervalMin ?? 1,
      ParkResumeMaxAttempts: q.ParkResumeMaxAttempts ?? 3,
      ParkResumeMenuTimeoutSec: q.ParkResumeMenuTimeoutSec ?? 15,
      ParkResumeDialTimeoutSec: q.ParkResumeDialTimeoutSec ?? 20,
      ParkConfirmation: q.ParkConfirmation || '',
      Members: (q.Members || []).map(m => ({
        Member: m.Member,
        Extension: m.Extension || '',
        Penalty: m.Penalty ?? 0,
        Disabled: !!m.Disabled
      }))
    };
    currentView = 'edit';
  }

  function backToList() {
    currentView = 'list';
    editingQueue = null;
  }

  async function handleSave() {
    if (!form.Name.trim()) return toast.error('Queue Name is required');

    saving = true;
    try {
      const payload: any = {
        Name: form.Name.trim(),
        Extension: form.Extension.trim(),
        MaxCapacity: Number(form.MaxCapacity),
        RingStrategy: form.RingStrategy,
        MemberRingTimeout: Number(form.MemberRingTimeout),
        MemberRetryDelay: Number(form.MemberRetryDelay),
        MemberWrapUpDelay: Number(form.MemberWrapUpDelay),
        InternalNumber: form.InternalNumber.trim(),
        PlayMusicWhileWaiting: form.PlayMusicWhileWaiting,
        Music: form.Music,
        CallerNamePrefix: form.CallerNamePrefix.trim(),
        RingCadence: form.RingCadence,
        AnnounceWaitTime: form.AnnounceWaitTime,
        AnnouncePosition: form.AnnouncePosition,
        AnnouncementInterval: Number(form.AnnouncementInterval),
        NoAnswerTimeout: Number(form.NoAnswerTimeout),
        NoAnswerDestination: form.NoAnswerDestination,
        ParkEnabled: form.ParkEnabled,
        ParkTelephoneNumber: form.ParkTelephoneNumber,
        ParkResumeWhenAtPosition: Number(form.ParkResumeWhenAtPosition),
        ParkAllowedMinPosition: Number(form.ParkAllowedMinPosition),
        ParkResumeRetryIntervalMin: Number(form.ParkResumeRetryIntervalMin),
        ParkResumeMaxAttempts: Number(form.ParkResumeMaxAttempts),
        ParkResumeMenuTimeoutSec: Number(form.ParkResumeMenuTimeoutSec),
        ParkResumeDialTimeoutSec: Number(form.ParkResumeDialTimeoutSec),
        ParkConfirmation: form.ParkConfirmation,
        Members: form.Members.map(m => ({
          Member: m.Member,
          Extension: m.Extension,
          Penalty: Number(m.Penalty),
          Disabled: m.Disabled
        }))
      };

      if (editingQueue) {
        await api.put(`/queues/${editingQueue.ID}`, payload);
        toast.success(`Queue '${form.Name}' updated`);
      } else {
        await api.post('/queues', payload);
        toast.success(`Queue '${form.Name}' created`);
      }

      await loadData();
      currentView = 'list';
    } catch (e: any) {
      toastError(e, 'Failed to save call queue');
    } finally {
      saving = false;
    }
  }

  async function deleteQueue(q: QueueItem) {
    if (!confirm(`Are you sure you want to delete call queue '${q.Name}'?`)) return;
    try {
      await api.delete(`/queues/${q.ID}`);
      toast.success(`Queue '${q.Name}' deleted`);
      loadData();
    } catch (e: any) {
      toastError(e, 'Failed to delete queue');
    }
  }

  // Member CRUD & toggle
  function openAddMember() {
    editingMemberIndex = -1;
    const defaultExt = availableExts.length > 0 ? availableExts[0] : null;
    memberForm = {
      Member: defaultExt ? (defaultExt.DisplayName || `${defaultExt.FirstName || ''} ${defaultExt.LastName || ''}`.trim() || `User ${defaultExt.ExtensionNumber}`) : '',
      Extension: defaultExt ? defaultExt.ExtensionNumber : '',
      Penalty: 0,
      Disabled: false
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

  function toggleDisabled(index: number) {
    form.Members[index].Disabled = !form.Members[index].Disabled;
  }

  $: filteredList = queuesList.filter(q =>
    !search.trim() ||
    `${q.Name} ${q.InternalNumber} ${q.Extension} ${q.RingStrategy}`.toLowerCase().includes(search.trim().toLowerCase())
  );
</script>

<div class="p-6 bg-[#EEF2F6] min-h-screen text-slate-800">
  <!-- Top Navigation Tabs: Call Queues (ACD) vs Live Wallboard -->
  <div class="flex items-center space-x-2 border-b border-slate-300 pb-3 mb-5">
    <button
      type="button"
      on:click={() => { activeTab = 'queues'; currentView = 'list'; }}
      class="px-4 py-2 font-bold text-xs rounded-md transition-colors {activeTab === 'queues' ? 'bg-[#0284C7] text-white shadow-sm' : 'bg-white text-slate-700 hover:bg-slate-100 border border-slate-200'}"
    >
      Call Queues (ACD Management)
    </button>
    <button
      type="button"
      on:click={() => activeTab = 'wallboard'}
      class="px-4 py-2 font-bold text-xs rounded-md transition-colors flex items-center space-x-1.5 {activeTab === 'wallboard' ? 'bg-[#0284C7] text-white shadow-sm' : 'bg-white text-slate-700 hover:bg-slate-100 border border-slate-200'}"
    >
      <span class="w-2 h-2 rounded-full {wsConnected ? 'bg-emerald-400 animate-pulse' : 'bg-slate-300'}"></span>
      <span>Live Realtime Wallboard</span>
    </button>
  </div>

  {#if activeTab === 'queues'}
    {#if currentView === 'list'}
      <!-- QUEUES LIST VIEW -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-4 mb-4 border-b border-slate-200 gap-4">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-lg bg-sky-100 text-[#0284C7] flex items-center justify-center font-bold shadow-sm">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 100-6 3 3 0 000 6z" />
            </svg>
          </div>
          <div>
            <h1 class="text-xl font-bold text-slate-900 tracking-tight">Call Queues (ACD)</h1>
            <p class="text-xs text-slate-500">Configure hunt strategies, wrap up delays, wait announcements, and park failover</p>
          </div>
        </div>

        <div class="flex items-center space-x-3">
          <div class="relative">
            <input
              type="text"
              bind:value={search}
              placeholder="Search queues..."
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
            <span>Add Queue</span>
          </button>
        </div>
      </div>

      <!-- Table -->
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <table class="w-full text-left text-xs text-slate-700">
          <thead class="bg-[#2B527A] text-white font-semibold">
            <tr>
              <th class="px-4 py-3">Queue Name</th>
              <th class="px-4 py-3">Internal Ext</th>
              <th class="px-4 py-3">Strategy</th>
              <th class="px-4 py-3">Capacity</th>
              <th class="px-4 py-3">Members</th>
              <th class="px-4 py-3">Failover Target</th>
              <th class="px-4 py-3 text-right">Options</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            {#each filteredList as item}
              <tr class="hover:bg-slate-50 transition-colors">
                <td class="px-4 py-3 font-bold text-slate-900">{item.Name}</td>
                <td class="px-4 py-3 text-slate-600 font-mono font-bold">{item.InternalNumber || item.Extension}</td>
                <td class="px-4 py-3">
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-sky-50 text-[#0284C7] border border-sky-200">
                    {item.RingStrategy}
                  </span>
                </td>
                <td class="px-4 py-3 text-slate-600">{item.MaxCapacity} callers</td>
                <td class="px-4 py-3">
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-50 text-emerald-800">
                    {(item.Members || []).length} Agents
                  </span>
                </td>
                <td class="px-4 py-3 text-slate-600 truncate max-w-xs">{item.NoAnswerDestination || 'Voicemail'}</td>
                <td class="px-4 py-3 text-right">
                  <div class="inline-flex items-center space-x-2">
                    <button
                      type="button"
                      on:click={() => openEdit(item)}
                      title="Edit Queue"
                      class="p-1 text-slate-600 hover:text-[#0284C7] hover:bg-sky-50 rounded"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                      </svg>
                    </button>
                    <button
                      type="button"
                      on:click={() => deleteQueue(item)}
                      title="Delete Queue"
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
                  No call queues configured. Click "Add Queue" to create one.
                </td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>

    {:else}
      <!-- EDIT VIEW MATCHING SCREENSHOT 3 -->
      <!-- Top Bar with queue title -->
      <div class="flex items-center space-x-2.5 pb-4 mb-4 border-b border-slate-300">
        <div class="text-[#0284C7] flex items-center space-x-1 font-bold text-xs">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 100-6 3 3 0 000 6z" />
          </svg>
        </div>
        <h1 class="text-sm font-bold text-slate-900 tracking-tight">
          {editingQueue ? `Editing queue '${form.Name}'` : 'Create New Queue'}
        </h1>
      </div>

      <form on:submit|preventDefault={handleSave} class="space-y-4">
        <!-- CARD 1: Queue Details (3-Column Layout matching Screenshot 3) -->
        <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
          <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
            Queue Details
          </div>
          <div class="p-5 grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Column 1 -->
            <div class="space-y-3.5">
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Queue Name</label>
                <input
                  type="text"
                  bind:value={form.Name}
                  placeholder="e.g. Customer Support Queue"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-medium"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Maximum Capacity</label>
                <input
                  type="number"
                  bind:value={form.MaxCapacity}
                  min="1"
                  max="500"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Ring Strategy</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg>
                  </span>
                  <select
                    bind:value={form.RingStrategy}
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  >
                    <option value="Ring All">Ring All</option>
                    <option value="Round Robin">Round Robin</option>
                    <option value="Least Recent">Least Recent</option>
                    <option value="Fewest Calls">Fewest Calls</option>
                    <option value="Random">Random</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Member Ring Timeout (s)</label>
                <input
                  type="number"
                  bind:value={form.MemberRingTimeout}
                  min="5"
                  max="120"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Member Retry Delay (s)</label>
                <input
                  type="number"
                  bind:value={form.MemberRetryDelay}
                  min="0"
                  max="60"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Member Wrap Up Delay (s)</label>
                <input
                  type="number"
                  bind:value={form.MemberWrapUpDelay}
                  min="0"
                  max="300"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
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
                    placeholder="e.g. 991"
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none font-bold"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Play music while waiting</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"/></svg>
                  </span>
                  <select
                    bind:value={form.PlayMusicWhileWaiting}
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  >
                    <option value={true}>Yes</option>
                    <option value={false}>No</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Music</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"/></svg>
                  </span>
                  <select
                    bind:value={form.Music}
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  >
                    <option value="Default">Default</option>
                    {#each audioPrompts as prompt}
                      <option value={prompt.Name}>{prompt.Name}</option>
                    {/each}
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
                    placeholder="e.g. Accounts"
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
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
                    <option value="Please select">Please select</option>
                    <option value="Default">Default</option>
                    <option value="Standard">Standard</option>
                    <option value="Ring-Ring">Ring-Ring</option>
                    <option value="Short-Short">Short-Short</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Announce estimated wait time</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400 font-bold text-xs">⏸</span>
                  <select
                    bind:value={form.AnnounceWaitTime}
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  >
                    <option value={false}>No</option>
                    <option value={true}>Yes</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Announce queue position</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400 font-bold text-xs">📍</span>
                  <select
                    bind:value={form.AnnouncePosition}
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  >
                    <option value={true}>Yes</option>
                    <option value={false}>No</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Announcement interval</label>
                <input
                  type="number"
                  bind:value={form.AnnouncementInterval}
                  min="10"
                  max="300"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
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
                    max="600"
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

            <!-- Column 3 (Call Park / Failover Settings) -->
            <div class="space-y-3.5">
              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Enabled</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400 font-bold text-xs">P</span>
                  <select
                    bind:value={form.ParkEnabled}
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  >
                    <option value={false}>No</option>
                    <option value={true}>Yes</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Telephone Number</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg>
                  </span>
                  <select
                    bind:value={form.ParkTelephoneNumber}
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
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Resume When at Position</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400 font-bold text-xs">⏸</span>
                  <input
                    type="number"
                    bind:value={form.ParkResumeWhenAtPosition}
                    min="1"
                    max="10"
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Allowed Minimum Position</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400 font-bold text-xs">⏸</span>
                  <input
                    type="number"
                    bind:value={form.ParkAllowedMinPosition}
                    min="1"
                    max="10"
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Resume Retry Interval (minutes)</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                  </span>
                  <input
                    type="number"
                    bind:value={form.ParkResumeRetryIntervalMin}
                    min="1"
                    max="60"
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Resume Max Attempts</label>
                <input
                  type="number"
                  bind:value={form.ParkResumeMaxAttempts}
                  min="1"
                  max="10"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Resume Menu Timeout (seconds)</label>
                <input
                  type="number"
                  bind:value={form.ParkResumeMenuTimeoutSec}
                  min="1"
                  max="60"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Resume Dial Timeout (seconds)</label>
                <input
                  type="number"
                  bind:value={form.ParkResumeDialTimeoutSec}
                  min="1"
                  max="120"
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-700 mb-1">Park Confirmation</label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-slate-400 font-bold text-xs">▶</span>
                  <select
                    bind:value={form.ParkConfirmation}
                    class="w-full pl-9 pr-3 py-2 bg-white border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
                  >
                    <option value="">Please select</option>
                    {#each audioPrompts as prompt}
                      <option value={prompt.Name}>{prompt.Name}</option>
                    {/each}
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

        <!-- CARD 2: Members (Matching Screenshot 3 Bottom Card) -->
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

            <!-- Table with Navy Header matching Screenshot 3 -->
            <div class="border border-slate-200 rounded overflow-hidden">
              <table class="w-full text-left text-xs">
                <thead class="bg-[#2B527A] text-white font-semibold">
                  <tr>
                    <th class="px-4 py-2.5">Member</th>
                    <th class="px-4 py-2.5">Penalty</th>
                    <th class="px-4 py-2.5">Disabled</th>
                    <th class="px-4 py-2.5 text-right">Options</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 bg-white">
                  {#each form.Members as mem, idx}
                    <tr class="hover:bg-slate-50 transition-colors">
                      <td class="px-4 py-2.5 font-medium text-slate-800">
                        {mem.Member} {mem.Extension ? `(${mem.Extension})` : ''}
                      </td>
                      <td class="px-4 py-2.5 text-slate-600 font-mono">{mem.Penalty}</td>
                      <td class="px-4 py-2.5">
                        <button
                          type="button"
                          on:click={() => toggleDisabled(idx)}
                          class="px-3 py-0.5 rounded text-[11px] font-bold cursor-pointer transition-colors {mem.Disabled ? 'bg-rose-100 text-rose-700 hover:bg-rose-200' : 'bg-emerald-600 text-white hover:bg-emerald-700'}"
                        >
                          {mem.Disabled ? 'Yes' : 'No'}
                        </button>
                      </td>
                      <td class="px-4 py-2.5 text-right">
                        <div class="inline-flex items-center space-x-2">
                          <button
                            type="button"
                            on:click={() => editMember(idx)}
                            class="p-1 text-slate-500 hover:text-[#0284C7]"
                            title="Edit Agent Penalty"
                          >
                            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                          </button>
                          <button
                            type="button"
                            on:click={() => removeMember(idx)}
                            class="p-1 text-slate-500 hover:text-rose-600"
                            title="Remove Agent"
                          >
                            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                          </button>
                        </div>
                      </td>
                    </tr>
                  {/each}
                  {#if form.Members.length === 0}
                    <tr>
                      <td colspan="4" class="px-4 py-6 text-center text-slate-400 bg-slate-50/50">
                        No agents assigned to this queue. Click "+ Add" to add members.
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

  {:else}
    <!-- LIVE WALLBOARD METRICS VIEW -->
    <div class="space-y-6">
      <div class="flex items-center justify-between pb-3 border-b border-slate-200">
        <div>
          <h2 class="text-lg font-bold text-slate-900 tracking-tight">Call Center Realtime Wallboard</h2>
          <p class="text-xs text-slate-500">Live operational telemetry streamed from FreeSWITCH event socket</p>
        </div>
        <div class="flex items-center space-x-2 text-xs">
          <span class="w-2.5 h-2.5 rounded-full {wsConnected ? 'bg-emerald-500 animate-pulse' : 'bg-rose-500'}"></span>
          <span class="font-medium text-slate-600">{wsConnected ? 'Realtime Stream Active' : 'Connecting to Core...'}</span>
        </div>
      </div>

      <!-- KPI Metrics Cards -->
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        <div class="bg-white p-4 rounded-lg border border-slate-200 shadow-sm text-center">
          <div class="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-1">Active Calls</div>
          <div class="text-2xl font-black text-[#0284C7]">{metrics.active_calls}</div>
        </div>
        <div class="bg-white p-4 rounded-lg border border-slate-200 shadow-sm text-center">
          <div class="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-1">Agents Online</div>
          <div class="text-2xl font-black text-emerald-600">{metrics.agents_online}</div>
        </div>
        <div class="bg-white p-4 rounded-lg border border-slate-200 shadow-sm text-center">
          <div class="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-1">Agents Busy</div>
          <div class="text-2xl font-black text-amber-600">{metrics.agents_busy}</div>
        </div>
        <div class="bg-white p-4 rounded-lg border border-slate-200 shadow-sm text-center">
          <div class="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-1">Queue Depth</div>
          <div class="text-2xl font-black text-purple-600">{metrics.queue_length}</div>
        </div>
        <div class="bg-white p-4 rounded-lg border border-slate-200 shadow-sm text-center">
          <div class="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-1">Avg Wait Time</div>
          <div class="text-2xl font-black text-slate-800">{metrics.avg_wait_time}s</div>
        </div>
        <div class="bg-white p-4 rounded-lg border border-slate-200 shadow-sm text-center">
          <div class="text-[11px] font-bold text-slate-500 uppercase tracking-wider mb-1">SLA Target</div>
          <div class="text-2xl font-black text-emerald-600">{metrics.sla_percentage}%</div>
        </div>
      </div>

      <!-- Live Queue Wallboard Cards -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-sm overflow-hidden">
        <div class="px-5 py-3 border-b border-slate-100 font-bold text-xs text-slate-800">
          Live Queues Monitor
        </div>
        <div class="p-5">
          {#if liveQueueStats.length > 0}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each liveQueueStats as q}
                <div class="p-4 border border-slate-200 rounded-lg bg-slate-50/50 flex flex-col justify-between">
                  <div class="flex items-center justify-between mb-3">
                    <div class="font-bold text-xs text-slate-900">{q.name}</div>
                    <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-100 text-emerald-800">SLA {q.sla}%</span>
                  </div>
                  <div class="grid grid-cols-2 gap-2 text-center">
                    <div class="bg-white p-2 rounded border border-slate-200">
                      <div class="text-[10px] text-slate-400 font-semibold uppercase">Callers Waiting</div>
                      <div class="text-lg font-bold text-slate-800">{q.waiting}</div>
                    </div>
                    <div class="bg-white p-2 rounded border border-slate-200">
                      <div class="text-[10px] text-slate-400 font-semibold uppercase">Logged Agents</div>
                      <div class="text-lg font-bold text-[#0284C7]">{q.agents}</div>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {:else}
            <div class="p-8 text-center text-slate-400 text-xs">
              No active queues currently reporting live traffic. Queues and agent metrics will appear here in real time when calls arrive.
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<!-- Modal: Queue Member Editor -->
{#if memberModalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-[2px]">
    <div class="bg-white rounded-lg shadow-xl border border-slate-200 w-full max-w-md overflow-hidden">
      <div class="px-5 py-3.5 bg-[#2B527A] text-white font-bold text-xs flex justify-between items-center">
        <span>{editingMemberIndex >= 0 ? 'Edit Queue Agent' : 'Add Queue Agent'}</span>
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
          <label class="block text-xs font-semibold text-slate-700 mb-1">Agent Display Name</label>
          <input
            type="text"
            bind:value={memberForm.Member}
            placeholder="e.g. Agent Name"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Skill Penalty (0 = Highest Priority)</label>
          <input
            type="number"
            bind:value={memberForm.Penalty}
            min="0"
            max="10"
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Disabled</label>
          <select
            bind:value={memberForm.Disabled}
            class="w-full px-3 py-2 border border-slate-300 rounded text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
          >
            <option value={false}>No (Active)</option>
            <option value={true}>Yes (Paused)</option>
          </select>
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
          Save Agent
        </button>
      </div>
    </div>
  </div>
{/if}
