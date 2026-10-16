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

  interface AudioPrompt {
    id: string;
    name: string;
    type: 'moh' | 'system' | 'custom';
    description: string;
    file_path: string;
    is_default: boolean;
  }

  interface RecordingRule {
    id: string;
    name: string;
    target_type: string;
    target_id: string;
    direction: string;
    storage_days: number;
    is_active: boolean;
  }

  let activeTab: 'moh' | 'custom' | 'system' | 'recording' = 'moh';

  let prompts: AudioPrompt[] = [];
  let recordingRules: RecordingRule[] = [];
  let availableExtensions: { id: string; number: string; label: string }[] = [];
  let availableQueues: { id: string; number: string; label: string }[] = [];
  let availableTrunks: { id: string; number: string; label: string }[] = [];

  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  // Prompt Modal
  let promptModalOpen = false;
  let editingPromptId: string | null = null;
  let promptForm = {
    name: '',
    type: 'moh' as 'moh' | 'custom' | 'system',
    description: '',
    file_path: 'local_stream://moh',
    is_default: false
  };

  // Recording Rule Modal
  let ruleModalOpen = false;
  let editingRuleId: string | null = null;
  let ruleForm = {
    name: '',
    target_type: 'all',
    target_id: '',
    direction: 'both',
    storage_days: 90,
    is_active: true
  };

  $: filteredPrompts = prompts
    .filter(p => activeTab === 'recording' || p.type === activeTab)
    .filter(p => !search.trim() || `${p.name} ${p.description}`.toLowerCase().includes(search.trim().toLowerCase()));

  $: filteredRules = recordingRules.filter(r =>
    !search.trim() || `${r.name} ${r.target_type} ${r.target_id}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [pRes, rRes, extRes, qRes, trunkRes] = await Promise.all([
        api.get<{ prompts: AudioPrompt[] }>('/audio-prompts').catch(() => ({ prompts: [] })),
        api.get<{ rules: RecordingRule[] }>('/recording-rules').catch(() => ({ rules: [] })),
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] })),
        api.get<{ queues: any[] }>('/queues').catch(() => ({ queues: [] })),
        api.get<{ trunks: any[] }>('/trunks').catch(() => ({ trunks: [] }))
      ]);
      prompts = pRes.prompts || [];
      recordingRules = rRes.rules || [];

      availableExtensions = (extRes.extensions || []).map((e: any) => ({
        id: e.ID,
        number: e.ExtensionNumber,
        label: `Ext ${e.ExtensionNumber} (${e.DisplayName || 'User'})`
      }));

      availableQueues = (qRes.queues || []).map((q: any) => ({
        id: q.ID,
        number: q.Extension || q.Name,
        label: `Queue: ${q.Name} (${q.Extension || ''})`
      }));

      availableTrunks = (trunkRes.trunks || []).map((t: any) => ({
        id: t.ID,
        number: t.ProviderName,
        label: `Trunk: ${t.ProviderName} (${t.SIPServer})`
      }));
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreatePrompt(type: 'moh' | 'custom' | 'system' | 'recording') {
    editingPromptId = null;
    const resolvedType = type === 'recording' ? 'custom' : type;
    promptForm = {
      name: '',
      type: resolvedType,
      description: '',
      file_path: '',
      is_default: false
    };
    promptModalOpen = true;
  }

  function openEditPrompt(p: AudioPrompt) {
    editingPromptId = p.id;
    promptForm = {
      name: p.name,
      type: p.type,
      description: p.description || '',
      file_path: p.file_path,
      is_default: p.is_default
    };
    promptModalOpen = true;
  }

  async function savePrompt(e: Event) {
    e.preventDefault();
    if (!promptForm.name.trim()) return toast.error('Prompt title is required');
    if (!promptForm.file_path.trim()) return toast.error('Audio file path is required');

    saving = true;
    try {
      if (editingPromptId) {
        await api.put(`/audio-prompts/${editingPromptId}`, promptForm);
        toast.success('Audio prompt updated');
      } else {
        await api.post('/audio-prompts', promptForm);
        toast.success('Audio prompt created');
      }
      promptModalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deletePrompt(id: string, name: string) {
    if (!confirm(`Delete audio file "${name}"?`)) return;
    try {
      await api.del(`/audio-prompts/${id}`);
      toast.success('Prompt deleted');
      prompts = prompts.filter(p => p.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  // RECORDING RULES
  function openCreateRule() {
    editingRuleId = null;
    ruleForm = {
      name: '',
      target_type: 'all',
      target_id: '',
      direction: 'both',
      storage_days: 90,
      is_active: true
    };
    ruleModalOpen = true;
  }

  function openEditRule(r: RecordingRule) {
    editingRuleId = r.id;
    ruleForm = {
      name: r.name,
      target_type: r.target_type,
      target_id: r.target_id || '',
      direction: r.direction,
      storage_days: r.storage_days || 90,
      is_active: r.is_active
    };
    ruleModalOpen = true;
  }

  async function saveRule(e: Event) {
    e.preventDefault();
    if (!ruleForm.name.trim()) return toast.error('Rule name is required');

    saving = true;
    try {
      if (editingRuleId) {
        await api.put(`/recording-rules/${editingRuleId}`, ruleForm);
        toast.success('Recording policy updated');
      } else {
        await api.post('/recording-rules', ruleForm);
        toast.success('Recording policy created');
      }
      ruleModalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteRule(id: string, name: string) {
    if (!confirm(`Delete recording rule "${name}"?`)) return;
    try {
      await api.del(`/recording-rules/${id}`);
      toast.success('Rule removed');
      recordingRules = recordingRules.filter(r => r.id !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Voice Prompts & Call Recording"
    subtitle="Configure Music on Hold (MoH), interactive voice announcements, and compliance call recording policies."
  >
    {#if activeTab === 'recording'}
      <Button variant="primary" on:click={openCreateRule}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Recording Policy
      </Button>
    {:else}
      <Button variant="primary" on:click={() => openCreatePrompt(activeTab)}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add Audio File
      </Button>
    {/if}
  </PageHeader>

  <!-- TAB SWITCHER -->
  <div class="flex border-b border-slate-200 space-x-6 text-sm font-medium">
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'moh' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'moh'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
      </svg>
      <span>Music on Hold (MoH) ({prompts.filter(p => p.type === 'moh').length})</span>
    </button>
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'custom' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'custom'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 100-6 3 3 0 000 6z" />
      </svg>
      <span>Custom Prompts ({prompts.filter(p => p.type === 'custom').length})</span>
    </button>
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'system' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'system'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
      <span>System Prompts ({prompts.filter(p => p.type === 'system').length})</span>
    </button>
    <button
      type="button"
      class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'recording' ? 'border-sky-600 text-sky-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
      on:click={() => activeTab = 'recording'}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
      </svg>
      <span>Call Recording Policies ({recordingRules.length})</span>
    </button>
  </div>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder={activeTab === 'recording' ? 'Filter recording rules...' : 'Filter audio files...'}
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
    <div class="p-12 text-center text-slate-400 text-sm">Loading audio assets and policies...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if activeTab !== 'recording'}
    <!-- AUDIO PROMPTS LIST -->
    {#if filteredPrompts.length === 0}
      <EmptyState
        title="No Audio Prompts in this Category"
        description="Add high-fidelity WAV or MP3 audio files for call queues, parked callers, and IVR attendants."
      >
        <Button variant="primary" on:click={() => openCreatePrompt(activeTab)}>Upload New Audio</Button>
      </EmptyState>
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Audio Prompt Title</th>
                <th class="py-3 px-4">Description</th>
                <th class="py-3 px-4">Media Stream / Path</th>
                <th class="py-3 px-4">Role</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredPrompts as p}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-semibold text-slate-900 flex items-center space-x-2">
                    <span>{p.name}</span>
                    {#if p.is_default}
                      <span class="px-1.5 py-0.5 bg-sky-100 text-sky-800 rounded text-[10px] font-bold uppercase tracking-wider">
                        Default
                      </span>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-slate-600">{p.description || '—'}</td>
                  <td class="py-3 px-4 font-mono text-slate-500 text-[11px]">{p.file_path}</td>
                  <td class="py-3 px-4 uppercase text-[11px] font-semibold text-slate-500">
                    <Badge variant={p.is_default ? 'info' : 'neutral'}>{p.type}</Badge>
                  </td>
                  <td class="py-3 px-4 text-right space-x-2">
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 font-medium"
                      on:click={() => openEditPrompt(p)}
                    >
                      Edit
                    </button>
                    {#if !p.is_default}
                      <button
                        type="button"
                        class="text-rose-600 hover:text-rose-800 font-medium"
                        on:click={() => deletePrompt(p.id, p.name)}
                      >
                        Delete
                      </button>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </Card>
    {/if}
  {:else}
    <!-- RECORDING POLICIES -->
    {#if filteredRules.length === 0}
      <EmptyState
        title="No Automatic Recording Rules Configured"
        description="Create enterprise compliance recording policies to automatically capture audio on inbound, outbound, or queue calls."
      >
        <Button variant="primary" on:click={openCreateRule}>Add Recording Policy</Button>
      </EmptyState>
    {:else}
      <Card padding="none">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
                <th class="py-3 px-4">Policy Name</th>
                <th class="py-3 px-4">Target Scope</th>
                <th class="py-3 px-4">Call Direction</th>
                <th class="py-3 px-4">Retention Window</th>
                <th class="py-3 px-4">Status</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each filteredRules as r}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-4 font-semibold text-slate-900">{r.name}</td>
                  <td class="py-3 px-4 capitalize">
                    <span class="font-medium text-slate-800">{r.target_type}</span>
                    {#if r.target_id}
                      <span class="ml-1 text-slate-500 font-mono">({r.target_id})</span>
                    {/if}
                  </td>
                  <td class="py-3 px-4 capitalize">
                    <span class="px-2 py-0.5 bg-slate-100 rounded text-slate-700 font-medium text-[11px]">
                      {r.direction} Calls
                    </span>
                  </td>
                  <td class="py-3 px-4 font-mono text-slate-700 font-medium">
                    {r.storage_days} Days
                  </td>
                  <td class="py-3 px-4">
                    {#if r.is_active}
                      <Badge variant="success">Recording Active</Badge>
                    {:else}
                      <Badge variant="neutral">Paused</Badge>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-right space-x-2">
                    <button
                      type="button"
                      class="text-sky-600 hover:text-sky-800 font-medium"
                      on:click={() => openEditRule(r)}
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
  {/if}
</div>

<!-- AUDIO PROMPT MODAL -->
<Modal
  open={promptModalOpen}
  title={editingPromptId ? 'Edit Audio Prompt' : 'Add Audio Prompt'}
  on:close={() => promptModalOpen = false}
>
  <form on:submit={savePrompt} class="space-y-4 text-xs">
    <div>
      <label for="ap-name" class="block font-medium text-slate-700 mb-1">Prompt Title</label>
      <input
        id="ap-name"
        type="text"
        bind:value={promptForm.name}
        placeholder="e.g. Corporate Chime or Holiday Announcement"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div>
      <label for="ap-type" class="block font-medium text-slate-700 mb-1">Category</label>
      <select
        id="ap-type"
        bind:value={promptForm.type}
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
      >
        <option value="moh">Music on Hold (MoH)</option>
        <option value="custom">Custom Auto-Attendant Prompt</option>
        <option value="system">System Notification</option>
      </select>
    </div>

    <div>
      <label for="ap-path" class="block font-medium text-slate-700 mb-1">File Path / Stream URI</label>
      <input
        id="ap-path"
        type="text"
        bind:value={promptForm.file_path}
        placeholder="e.g. sounds/moh/chime.wav or local_stream://moh"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
      />
    </div>

    <div>
      <label for="ap-desc" class="block font-medium text-slate-700 mb-1">Description (Optional)</label>
      <input
        id="ap-desc"
        type="text"
        bind:value={promptForm.description}
        placeholder="e.g. Played when callers are placed on hold or in queue"
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="flex items-center space-x-2 pt-2">
      <input
        type="checkbox"
        id="ap_default"
        bind:checked={promptForm.is_default}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="ap_default" class="text-slate-700 font-medium">Set as default playlist for this category</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => promptModalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingPromptId ? 'Save Changes' : 'Save Audio'}
      </Button>
    </div>
  </form>
</Modal>

<!-- RECORDING RULE MODAL -->
<Modal
  open={ruleModalOpen}
  title={editingRuleId ? 'Edit Recording Policy' : 'Create Recording Policy'}
  on:close={() => ruleModalOpen = false}
>
  <form on:submit={saveRule} class="space-y-4 text-xs">
    <div>
      <label for="rr-name" class="block font-medium text-slate-700 mb-1">Policy Title</label>
      <input
        id="rr-name"
        type="text"
        bind:value={ruleForm.name}
        placeholder="e.g. Record All Support Desk Calls"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
      />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="rr-ttype" class="block font-medium text-slate-700 mb-1">Target Scope</label>
        <select
          id="rr-ttype"
          bind:value={ruleForm.target_type}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
        >
          <option value="all">Company-Wide (All Calls)</option>
          <option value="extension">Specific Extension</option>
          <option value="queue">Call Center Queue</option>
          <option value="trunk">SIP Trunk Gateway</option>
        </select>
      </div>
      <div>
        <label for="rr-tid" class="block font-medium text-slate-700 mb-1">Target Selection</label>
        {#if ruleForm.target_type === 'extension'}
          <select
            id="rr-tid"
            bind:value={ruleForm.target_id}
            class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
          >
            <option value="">-- Select Extension --</option>
            {#each availableExtensions as ext}
              <option value={ext.number}>{ext.label}</option>
            {/each}
          </select>
        {:else if ruleForm.target_type === 'queue'}
          <select
            id="rr-tid"
            bind:value={ruleForm.target_id}
            class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
          >
            <option value="">-- Select Queue --</option>
            {#each availableQueues as q}
              <option value={q.number}>{q.label}</option>
            {/each}
          </select>
        {:else if ruleForm.target_type === 'trunk'}
          <select
            id="rr-tid"
            bind:value={ruleForm.target_id}
            class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
          >
            <option value="">-- Select SIP Trunk --</option>
            {#each availableTrunks as t}
              <option value={t.number}>{t.label}</option>
            {/each}
          </select>
        {:else}
          <input
            id="rr-tid"
            type="text"
            disabled
            value="All Endpoints (Tenant-Wide Policy)"
            class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono disabled:bg-slate-100 italic text-slate-500"
          />
        {/if}
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="rr-dir" class="block font-medium text-slate-700 mb-1">Direction</label>
        <select
          id="rr-dir"
          bind:value={ruleForm.direction}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
        >
          <option value="both">Both Inbound & Outbound</option>
          <option value="inbound">Inbound Calls Only</option>
          <option value="outbound">Outbound Calls Only</option>
        </select>
      </div>
      <div>
        <label for="rr-days" class="block font-medium text-slate-700 mb-1">Retention Window (Days)</label>
        <input
          id="rr-days"
          type="number"
          min="1"
          max="3650"
          bind:value={ruleForm.storage_days}
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
    </div>

    <div class="flex items-center space-x-2 pt-2">
      <input
        type="checkbox"
        id="rule_active"
        bind:checked={ruleForm.is_active}
        class="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500"
      />
      <label for="rule_active" class="text-slate-700 font-medium">Activate this recording rule immediately</label>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => ruleModalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingRuleId ? 'Save Changes' : 'Create Policy'}
      </Button>
    </div>
  </form>
</Modal>
