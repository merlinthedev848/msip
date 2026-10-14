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

  interface InboundRouteItem {
    ID: string;
    DIDNumber: string;
    CallerIDPattern: string;
    Name: string;
    DestinationType: string;
    DestinationTarget: string;
    FallbackDestType: string;
    FallbackDestTarget: string;
    TimeConditionID: string;
    PrependCID: string;
    CallRecording: string;
    MaxRingTimeSeconds: number;
    CreatedAt: string;
  }

  let routes: InboundRouteItem[] = [];
  let availableExts: any[] = [];
  let availableIvrs: any[] = [];
  let availableGroups: any[] = [];
  let availableQueues: any[] = [];
  let availableAssistants: any[] = [];
  let availableTimeConds: any[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let activeTab: 'routing' | 'advanced' = 'routing';

  let form = {
    Name: '',
    DIDNumber: '',
    CallerIDPattern: '',
    DestinationType: 'EXTENSION',
    DestinationTarget: '',
    FallbackDestType: 'VOICEMAIL',
    FallbackDestTarget: '',
    TimeConditionID: '',
    PrependCID: '',
    CallRecording: 'off',
    MaxRingTimeSeconds: 30
  };

  $: filtered = routes.filter(r => 
    !search.trim() || `${r.Name || ''} ${r.DIDNumber} ${r.DestinationTarget} ${r.DestinationType}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const [routeData, extData, ivrData, timeData, groupData, queueData, aiData] = await Promise.all([
        api.get<{ routes: InboundRouteItem[] }>('/inbound-routes'),
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] })),
        api.get<{ ivrs: any[] }>('/ivrs').catch(() => ({ ivrs: [] })),
        api.get<{ conditions: any[] }>('/time-conditions').catch(() => ({ conditions: [] })),
        api.get<{ ring_groups: any[] }>('/ring-groups').catch(() => ({ ring_groups: [] })),
        api.get<{ queues: any[] }>('/queues').catch(() => ({ queues: [] })),
        api.get<{ assistants: any[] }>('/ai/assistants').catch(() => ({ assistants: [] }))
      ]);
      routes = routeData.routes || [];
      availableExts = extData.extensions || [];
      availableIvrs = ivrData.ivrs || [];
      availableTimeConds = timeData.conditions || [];
      availableGroups = groupData.ring_groups || [];
      availableQueues = queueData.queues || [];
      availableAssistants = aiData.assistants || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    activeTab = 'routing';
    const defaultTarget = availableExts.length > 0 ? availableExts[0].ExtensionNumber : '';
    form = {
      Name: '',
      DIDNumber: '',
      CallerIDPattern: '',
      DestinationType: 'EXTENSION',
      DestinationTarget: defaultTarget,
      FallbackDestType: 'VOICEMAIL',
      FallbackDestTarget: defaultTarget,
      TimeConditionID: '',
      PrependCID: '',
      CallRecording: 'off',
      MaxRingTimeSeconds: 30
    };
    modalOpen = true;
  }

  function openEdit(r: InboundRouteItem) {
    editingId = r.ID;
    activeTab = 'routing';
    form = {
      Name: r.Name || '',
      DIDNumber: r.DIDNumber,
      CallerIDPattern: r.CallerIDPattern || '',
      DestinationType: r.DestinationType || 'EXTENSION',
      DestinationTarget: r.DestinationTarget,
      FallbackDestType: r.FallbackDestType || 'VOICEMAIL',
      FallbackDestTarget: r.FallbackDestTarget || '',
      TimeConditionID: r.TimeConditionID || '',
      PrependCID: r.PrependCID || '',
      CallRecording: r.CallRecording || 'off',
      MaxRingTimeSeconds: r.MaxRingTimeSeconds || 30
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.DIDNumber.trim()) return toast.error('Inbound DID number or pattern is required');
    if (!form.DestinationTarget.trim()) return toast.error('Destination target is required');

    saving = true;
    try {
      const payload: any = {
        Name: form.Name.trim() || `Route for ${form.DIDNumber.trim()}`,
        DIDNumber: form.DIDNumber.trim(),
        CallerIDPattern: form.CallerIDPattern.trim(),
        DestinationType: form.DestinationType,
        DestinationTarget: form.DestinationTarget.trim(),
        FallbackDestType: form.FallbackDestType,
        FallbackDestTarget: form.FallbackDestTarget.trim(),
        TimeConditionID: form.TimeConditionID,
        PrependCID: form.PrependCID.trim(),
        CallRecording: form.CallRecording,
        MaxRingTimeSeconds: Number(form.MaxRingTimeSeconds) || 30
      };

      if (editingId) {
        await api.put(`/inbound-routes/${editingId}`, payload);
        toast.success('Inbound route updated');
      } else {
        await api.post('/inbound-routes', payload);
        toast.success('Inbound route created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(r: InboundRouteItem) {
    if (!confirm(`Delete inbound route for DID ${r.DIDNumber}?`)) return;
    try {
      await api.del(`/inbound-routes/${r.ID}`);
      toast.success('Inbound route deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Inbound Call Routing" description="Map incoming telephone numbers (DIDs) to extensions, Auto Attendants (IVRs), ring groups, or time conditions.">
  <Button on:click={openCreate}>+ New Inbound Rule</Button>
</PageHeader>

<Card padded={false} title="Configured Inbound Rules" subtitle="{routes.length} incoming dial patterns">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-60" placeholder="Search route, DID, destination…" bind:value={search} aria-label="Search routes" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No routes match your search' : 'No inbound routes configured'}
      message={search ? '' : 'Create your first inbound rule to decide where callers go when they dial your company numbers.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Create Inbound Route</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Rule Name</th>
            <th>Incoming DID Number</th>
            <th>Primary Destination</th>
            <th>Failover / Timeout</th>
            <th>Call Recording</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="font-bold text-slate-800 text-sm">
                <div>{r.Name || 'Standard Rule'}</div>
                {#if r.PrependCID}
                  <div class="text-[10px] text-amber-600 font-mono">CID Tag: [{r.PrependCID}]</div>
                {/if}
              </td>
              <td>
                <span class="mono px-2 py-0.5 bg-sky-50 text-sky-800 rounded text-xs border border-sky-200 font-bold">
                  {r.DIDNumber}
                </span>
              </td>
              <td>
                <div class="flex items-center space-x-1.5">
                  <Badge tone={r.DestinationType === 'IVR' ? 'warning' : 'neutral'}>
                    {r.DestinationType}
                  </Badge>
                  <span class="mono text-xs font-semibold text-slate-800">{r.DestinationTarget}</span>
                </div>
              </td>
              <td>
                {#if r.FallbackDestTarget}
                  <div class="text-xs text-slate-600">
                    <span class="text-slate-400">{r.FallbackDestType}:</span>
                    <span class="mono font-semibold">{r.FallbackDestTarget}</span>
                  </div>
                {:else}
                  <span class="text-xs text-slate-400">Hangup</span>
                {/if}
              </td>
              <td>
                <span class="text-xs font-medium capitalize {r.CallRecording === 'always' ? 'text-emerald-700 font-bold' : 'text-slate-500'}">
                  {r.CallRecording || 'Off'}
                </span>
              </td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(r)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(r)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Inbound Rule' : 'Create Inbound Rule'} maxWidth="max-w-2xl">
  <!-- Tabs -->
  <div class="flex items-center space-x-1 border-b border-slate-200 mb-4 pb-1">
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'routing' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'routing'}
    >
      Destination & DID
    </button>
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'advanced' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'advanced'}
    >
      Failover, Time & Recording
    </button>
  </div>

  <form on:submit={save} class="space-y-4">
    {#if activeTab === 'routing'}
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="sm:col-span-2">
          <label for="route-name" class="label">Rule Description Name</label>
          <input id="route-name" class="input" bind:value={form.Name} placeholder="e.g. Main Company Number, Support Line, UK DID" />
        </div>

        <div>
          <label for="route-did" class="label">Inbound DID Number (Match Pattern) *</label>
          <input id="route-did" class="input mono font-bold" bind:value={form.DIDNumber} placeholder="e.g. 02086884243 or _0208XXXXXXX" required />
          <p class="text-[11px] text-slate-500 mt-1">Exact phone number or FreeSWITCH dial pattern matched from carrier SIP INVITE.</p>
        </div>

        <div>
          <label for="route-cid" class="label">Caller ID Filter (Optional)</label>
          <input id="route-cid" class="input mono" bind:value={form.CallerIDPattern} placeholder="Leave blank to match all callers" />
        </div>

        <div>
          <label for="route-dtype" class="label">Primary Destination Type</label>
          <select id="route-dtype" class="input font-semibold" bind:value={form.DestinationType}>
            <option value="EXTENSION">Extension (Desk Phone / User)</option>
            <option value="AI_ASSISTANT">AI Virtual Receptionist (Conversational AI)</option>
            <option value="IVR">Digital Receptionist (Auto Attendant IVR)</option>
            <option value="RING_GROUP">Ring Group (Hunt Group)</option>
            <option value="QUEUE">Call Queue / Wallboard</option>
            <option value="VOICEMAIL">Voicemail Mailbox</option>
          </select>
        </div>

        <div>
          <label for="route-target" class="label">Primary Target *</label>
          {#if form.DestinationType === 'EXTENSION'}
            <select id="route-target" class="input font-mono" bind:value={form.DestinationTarget}>
              <option value="">Select Extension...</option>
              {#each availableExts as ext}
                <option value={ext.ExtensionNumber}>{ext.ExtensionNumber} - {ext.DisplayName || `${ext.FirstName || ''} ${ext.LastName || ''}`.trim() || 'User'}</option>
              {/each}
            </select>
          {:else if form.DestinationType === 'AI_ASSISTANT'}
            <select id="route-target" class="input font-mono" bind:value={form.DestinationTarget}>
              <option value="">Select AI Receptionist...</option>
              {#each availableAssistants as asst}
                <option value={asst.ExtensionNumber || asst.Extension || '890'}>🤖 {asst.Name} (Ext {asst.ExtensionNumber || asst.Extension || '890'})</option>
              {/each}
            </select>
          {:else if form.DestinationType === 'IVR'}
            <select id="route-target" class="input font-mono" bind:value={form.DestinationTarget}>
              <option value="">Select IVR Menu...</option>
              {#each availableIvrs as ivr}
                <option value={ivr.Extension || ivr.Name}>{ivr.Name} (Ext {ivr.Extension})</option>
              {/each}
            </select>
          {:else if form.DestinationType === 'RING_GROUP'}
            <select id="route-target" class="input font-mono" bind:value={form.DestinationTarget}>
              <option value="">Select Ring Group...</option>
              {#each availableGroups as grp}
                <option value={grp.InternalNumber || grp.Name}>👥 {grp.Name} (Ext {grp.InternalNumber})</option>
              {/each}
            </select>
          {:else if form.DestinationType === 'QUEUE'}
            <select id="route-target" class="input font-mono" bind:value={form.DestinationTarget}>
              <option value="">Select Call Queue...</option>
              {#each availableQueues as q}
                <option value={q.InternalNumber || q.Extension || q.Name}>⏳ {q.Name} (Ext {q.InternalNumber || q.Extension})</option>
              {/each}
            </select>
          {:else if form.DestinationType === 'VOICEMAIL'}
            <select id="route-target" class="input font-mono" bind:value={form.DestinationTarget}>
              <option value="*98">General Company Voicemail (*98)</option>
              {#each availableExts as ext}
                <option value={ext.ExtensionNumber}>Voicemail Box {ext.ExtensionNumber} ({ext.DisplayName || 'User'})</option>
              {/each}
            </select>
          {:else}
            <input id="route-target" class="input mono font-bold" bind:value={form.DestinationTarget} placeholder="Target extension (e.g. 1001)" required />
          {/if}
        </div>
      </div>
    {/if}

    {#if activeTab === 'advanced'}
      <div class="space-y-4">
        <!-- Failover -->
        <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label for="route-fb-type" class="label">No-Answer Failover Type</label>
            <select id="route-fb-type" class="input" bind:value={form.FallbackDestType}>
              <option value="VOICEMAIL">Voicemail</option>
              <option value="EXTENSION">Colleague Extension</option>
              <option value="AI_ASSISTANT">AI Virtual Receptionist</option>
              <option value="IVR">Auto Attendant (IVR)</option>
              <option value="HANGUP">Terminate Call</option>
            </select>
          </div>
          <div>
            <label for="route-fb-target" class="label">Failover Target</label>
            {#if form.FallbackDestType === 'EXTENSION' || form.FallbackDestType === 'VOICEMAIL'}
              <select id="route-fb-target" class="input font-mono" bind:value={form.FallbackDestTarget}>
                {#each availableExts as ext}
                  <option value={ext.ExtensionNumber}>Ext {ext.ExtensionNumber} ({ext.DisplayName || 'User'})</option>
                {/each}
              </select>
            {:else if form.FallbackDestType === 'AI_ASSISTANT'}
              <select id="route-fb-target" class="input font-mono" bind:value={form.FallbackDestTarget}>
                {#each availableAssistants as asst}
                  <option value={asst.ExtensionNumber || asst.Extension || '890'}>🤖 {asst.Name} (Ext {asst.ExtensionNumber || asst.Extension || '890'})</option>
                {/each}
              </select>
            {:else if form.FallbackDestType === 'IVR'}
              <select id="route-fb-target" class="input font-mono" bind:value={form.FallbackDestTarget}>
                {#each availableIvrs as ivr}
                  <option value={ivr.Extension || ivr.Name}>{ivr.Name} (Ext {ivr.Extension})</option>
                {/each}
              </select>
            {:else if form.FallbackDestType === 'HANGUP'}
              <input id="route-fb-target" class="input mono text-slate-400" value="Hang Up (Normal Disconnect)" disabled />
            {:else}
              <input id="route-fb-target" class="input mono" bind:value={form.FallbackDestTarget} placeholder="e.g. 1001" />
            {/if}
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label for="route-ringtime" class="label">Max Ring Timeout (Seconds)</label>
            <input id="route-ringtime" type="number" min="5" max="120" class="input mono font-bold" bind:value={form.MaxRingTimeSeconds} />
            <p class="text-[11px] text-slate-500 mt-1">Seconds to ring destination before triggering failover.</p>
          </div>

          <div>
            <label for="route-prepend" class="label">Prepend to Caller ID (Tagging)</label>
            <input id="route-prepend" class="input font-mono" bind:value={form.PrependCID} placeholder="e.g. Sales: or VIP:" />
            <p class="text-[11px] text-slate-500 mt-1">Prepend prefix to caller ID display so operators know which DID was dialed.</p>
          </div>

          <div>
            <label for="route-rec" class="label">Call Recording Policy</label>
            <select id="route-rec" class="input" bind:value={form.CallRecording}>
              <option value="off">Disabled (Do Not Record)</option>
              <option value="always">Always Record Inbound Calls</option>
            </select>
          </div>

          {#if availableTimeConds.length > 0}
            <div>
              <label for="route-tc" class="label">Apply Business Hours / Time Condition</label>
              <select id="route-tc" class="input" bind:value={form.TimeConditionID}>
                <option value="">No Time Condition (Always Route)</option>
                {#each availableTimeConds as tc}
                  <option value={tc.id}>{tc.name} ({tc.time_start} - {tc.time_end})</option>
                {/each}
              </select>
            </div>
          {/if}
        </div>
      </div>
    {/if}

    <div class="flex justify-between items-center pt-3 border-t border-slate-200">
      <div class="text-xs text-slate-400">
        {#if activeTab !== 'routing'}
          <button type="button" class="text-sky-600 hover:underline" on:click={() => activeTab = 'routing'}>&larr; Primary Route</button>
        {/if}
      </div>
      <div class="flex gap-2">
        <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
        <Button type="submit" disabled={saving}>{saving ? 'Saving…' : (editingId ? 'Update Route' : 'Create Route')}</Button>
      </div>
    </div>
  </form>
</Modal>
