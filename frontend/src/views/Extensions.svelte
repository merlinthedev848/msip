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

  interface ExtensionItem {
    ID: string;
    ExtensionNumber: string;
    IsActive: boolean;
    ProvisioningCode: string;
    FirstName: string;
    LastName: string;
    Email: string;
    Department: string;
    CallerIDName: string;
    CallerIDNumber: string;
    OutboundDID: string;
    MobileNumber: string;
    ForwardAlwaysEnabled: boolean;
    ForwardAlwaysTarget: string;
    ForwardBusyEnabled: boolean;
    ForwardBusyTarget: string;
    ForwardNoAnswerEnabled: boolean;
    ForwardNoAnswerTarget: string;
    RingTimeoutSeconds: number;
    VoicemailEnabled: boolean;
    VoicemailPIN: string;
    VoicemailToEmail: boolean;
    CreatedAt: string;
  }

  interface DIDOption {
    Number: string;
    Description: string;
  }

  let extensions: ExtensionItem[] = [];
  let availableDIDs: DIDOption[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let activeTab: 'basic' | 'cid' | 'forwarding' | 'voicemail' = 'basic';
  let showPassword = false;

  let form = {
    ExtensionNumber: '',
    Password: '',
    ConfirmPassword: '',
    IsActive: true,
    FirstName: '',
    LastName: '',
    Email: '',
    Department: 'General',
    CallerIDName: '',
    CallerIDNumber: '',
    OutboundDID: '',
    MobileNumber: '',
    ForwardAlwaysEnabled: false,
    ForwardAlwaysTarget: '',
    ForwardBusyEnabled: false,
    ForwardBusyTarget: '',
    ForwardNoAnswerEnabled: true,
    ForwardNoAnswerTarget: '',
    RingTimeoutSeconds: 20,
    VoicemailEnabled: true,
    VoicemailPIN: '',
    VoicemailToEmail: true
  };

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

  function regeneratePassword() {
    const p = generatePassword(14);
    form.Password = p;
    form.ConfirmPassword = p;
  }

  async function copyPassword() {
    if (!form.Password) return;
    try {
      await navigator.clipboard.writeText(form.Password);
      toast.success('SIP secret copied to clipboard');
    } catch {
      toast.error('Failed to copy password');
    }
  }

  $: filtered = extensions.filter(e => {
    const q = search.trim().toLowerCase();
    if (!q) return true;
    const fullName = `${e.FirstName || ''} ${e.LastName || ''}`.trim();
    const str = `${e.ExtensionNumber} ${fullName} ${e.Email || ''} ${e.CallerIDNumber || ''} ${e.OutboundDID || ''} ${e.Department || ''} ${e.ProvisioningCode || ''}`.toLowerCase();
    return str.includes(q);
  });

  async function load() {
    loading = true;
    loadError = '';
    try {
      const [extData, numData] = await Promise.all([
        api.get<{ extensions: ExtensionItem[] }>('/extensions'),
        api.get<{ numbers: any[] }>('/telephone-numbers').catch(() => ({ numbers: [] }))
      ]);
      extensions = extData.extensions || [];
      availableDIDs = (numData.numbers || []).map(n => ({
        Number: n.Number,
        Description: n.Description || n.Number
      }));
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    activeTab = 'basic';
    showPassword = true;
    const initialPass = generatePassword(14);
    form = {
      ExtensionNumber: '',
      Password: initialPass,
      ConfirmPassword: initialPass,
      IsActive: true,
      FirstName: '',
      LastName: '',
      Email: '',
      Department: 'General',
      CallerIDName: '',
      CallerIDNumber: '',
      OutboundDID: '',
      MobileNumber: '',
      ForwardAlwaysEnabled: false,
      ForwardAlwaysTarget: '',
      ForwardBusyEnabled: false,
      ForwardBusyTarget: '',
      ForwardNoAnswerEnabled: true,
      ForwardNoAnswerTarget: '',
      RingTimeoutSeconds: 20,
      VoicemailEnabled: true,
      VoicemailPIN: '',
      VoicemailToEmail: true
    };
    modalOpen = true;
  }

  function openEdit(ext: ExtensionItem) {
    editingId = ext.ID;
    activeTab = 'basic';
    showPassword = false;
    form = {
      ExtensionNumber: ext.ExtensionNumber,
      Password: '',
      ConfirmPassword: '',
      IsActive: ext.IsActive ?? true,
      FirstName: ext.FirstName || '',
      LastName: ext.LastName || '',
      Email: ext.Email || '',
      Department: ext.Department || 'General',
      CallerIDName: ext.CallerIDName || '',
      CallerIDNumber: ext.CallerIDNumber || '',
      OutboundDID: ext.OutboundDID || '',
      MobileNumber: ext.MobileNumber || '',
      ForwardAlwaysEnabled: ext.ForwardAlwaysEnabled || false,
      ForwardAlwaysTarget: ext.ForwardAlwaysTarget || '',
      ForwardBusyEnabled: ext.ForwardBusyEnabled || false,
      ForwardBusyTarget: ext.ForwardBusyTarget || '',
      ForwardNoAnswerEnabled: ext.ForwardNoAnswerEnabled ?? true,
      ForwardNoAnswerTarget: ext.ForwardNoAnswerTarget || '',
      RingTimeoutSeconds: ext.RingTimeoutSeconds || 20,
      VoicemailEnabled: ext.VoicemailEnabled ?? true,
      VoicemailPIN: ext.VoicemailPIN || '',
      VoicemailToEmail: ext.VoicemailToEmail ?? true
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.ExtensionNumber.trim()) return toast.error('Extension number is required');
    if (!editingId && !form.Password) return toast.error('SIP Password is required');
    if (form.Password && form.Password !== form.ConfirmPassword) {
      return toast.error('Passwords do not match');
    }

    saving = true;
    try {
      const payload: any = {
        ExtensionNumber: form.ExtensionNumber.trim(),
        IsActive: form.IsActive,
        FirstName: form.FirstName.trim(),
        LastName: form.LastName.trim(),
        Email: form.Email.trim(),
        Department: form.Department.trim(),
        CallerIDName: form.CallerIDName.trim() || `${form.FirstName} ${form.LastName}`.trim(),
        CallerIDNumber: form.CallerIDNumber.trim() || form.ExtensionNumber.trim(),
        OutboundDID: form.OutboundDID.trim(),
        MobileNumber: form.MobileNumber.trim(),
        ForwardAlwaysEnabled: form.ForwardAlwaysEnabled,
        ForwardAlwaysTarget: form.ForwardAlwaysTarget.trim(),
        ForwardBusyEnabled: form.ForwardBusyEnabled,
        ForwardBusyTarget: form.ForwardBusyTarget.trim(),
        ForwardNoAnswerEnabled: form.ForwardNoAnswerEnabled,
        ForwardNoAnswerTarget: form.ForwardNoAnswerTarget.trim() || (form.MobileNumber.trim() || ''),
        RingTimeoutSeconds: Number(form.RingTimeoutSeconds) || 20,
        VoicemailEnabled: form.VoicemailEnabled,
        VoicemailPIN: form.VoicemailPIN.trim() || form.ExtensionNumber.trim(),
        VoicemailToEmail: form.VoicemailToEmail
      };

      if (form.Password) {
        payload.PasswordHash = form.Password;
      }

      if (editingId) {
        await api.put(`/extensions/${editingId}`, payload);
        toast.success(`Extension ${form.ExtensionNumber} updated`);
      } else {
        await api.post('/extensions', payload);
        toast.success(`Extension ${form.ExtensionNumber} created`);
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(ext: ExtensionItem) {
    if (!confirm(`Permanently delete extension ${ext.ExtensionNumber} (${ext.FirstName} ${ext.LastName})?`)) return;
    try {
      await api.del(`/extensions/${ext.ID}`);
      toast.success('Extension deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Extensions & Users" description="Manage user profiles, SIP registrations, caller IDs, assigned DIDs, and mobile call forwarding.">
  <Button on:click={openCreate}>+ New Extension</Button>
</PageHeader>

<Card padded={false} title="User Extensions Directory" subtitle="{extensions.length} configured extensions">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-64" placeholder="Search by name, ext, DID, email…" bind:value={search} aria-label="Search extensions" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No extensions match your search' : 'No extensions found'}
      message={search ? '' : 'Create your first extension to connect a SIP user endpoint, configure caller ID, and assign a DID.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add First Extension</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Ext</th>
            <th>User / Display Name</th>
            <th>Email</th>
            <th>Assigned DID / Outbound CLI</th>
            <th>Forwarding / Mobile</th>
            <th>Status</th>
            <th>Provisioning</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as ext (ext.ID)}
            {@const fullName = `${ext.FirstName || ''} ${ext.LastName || ''}`.trim()}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="mono font-bold text-slate-900 text-sm">{ext.ExtensionNumber}</td>
              <td>
                <div class="font-semibold text-slate-800 text-xs">
                  {fullName || `Extension ${ext.ExtensionNumber}`}
                </div>
                {#if ext.Department}
                  <div class="text-[10px] text-slate-400 capitalize">{ext.Department}</div>
                {/if}
              </td>
              <td class="text-xs text-slate-600">
                {ext.Email || '—'}
              </td>
              <td>
                {#if ext.OutboundDID}
                  <span class="mono px-2 py-0.5 bg-emerald-50 text-emerald-800 rounded text-xs border border-emerald-200 font-semibold">
                    {ext.OutboundDID}
                  </span>
                {:else if ext.CallerIDNumber}
                  <span class="mono text-xs text-slate-600">
                    {ext.CallerIDNumber}
                  </span>
                {:else}
                  <span class="text-slate-400 text-xs">Internal Only</span>
                {/if}
              </td>
              <td>
                {#if ext.ForwardAlwaysEnabled && ext.ForwardAlwaysTarget}
                  <span class="inline-flex items-center space-x-1 text-[11px] text-amber-700 font-medium">
                    <span>Forward:</span>
                    <span class="font-mono font-bold">{ext.ForwardAlwaysTarget}</span>
                  </span>
                {:else if ext.MobileNumber}
                  <span class="inline-flex items-center space-x-1 text-[11px] text-sky-700">
                    <svg class="w-3 h-3 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/></svg>
                    <span class="font-mono">{ext.MobileNumber}</span>
                  </span>
                {:else}
                  <span class="text-slate-400 text-xs">Desk Only</span>
                {/if}
              </td>
              <td>
                <Badge tone={ext.IsActive ? 'success' : 'neutral'}>
                  {ext.IsActive ? 'Active' : 'Disabled'}
                </Badge>
              </td>
              <td>
                <span class="mono px-2 py-0.5 bg-sky-50 text-sky-800 rounded text-xs border border-sky-200">
                  {ext.ProvisioningCode || '—'}
                </span>
              </td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(ext)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(ext)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<!-- Comprehensive Extension Configuration Modal -->
<Modal bind:isOpen={modalOpen} title={editingId ? `Edit Extension ${form.ExtensionNumber}` : 'Create New Extension'} maxWidth="max-w-2xl">
  <!-- Tab Strip -->
  <div class="flex items-center space-x-1 border-b border-slate-200 mb-4 pb-1">
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'basic' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'basic'}
    >
      User & SIP Account
    </button>
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'cid' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'cid'}
    >
      Caller ID & DID Mapping
    </button>
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'forwarding' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'forwarding'}
    >
      Forwarding & Mobile
    </button>
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'voicemail' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'voicemail'}
    >
      Voicemail
    </button>
  </div>

  <form on:submit={save} class="space-y-4">

    <!-- TAB 1: BASIC USER & SIP ACCOUNT -->
    {#if activeTab === 'basic'}
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label for="ext-number" class="label">Extension Number *</label>
          <input id="ext-number" class="input mono font-bold" bind:value={form.ExtensionNumber} placeholder="e.g. 1001" required />
          <p class="text-[11px] text-slate-500 mt-1">2 to 5 digit extension assigned to this user.</p>
        </div>

        <div>
          <label for="ext-dept" class="label">Department / Group</label>
          <input id="ext-dept" class="input" bind:value={form.Department} placeholder="Sales, Support, Admin..." />
        </div>

        <div>
          <label for="ext-first" class="label">First Name</label>
          <input id="ext-first" class="input" bind:value={form.FirstName} placeholder="John" />
        </div>

        <div>
          <label for="ext-last" class="label">Last Name</label>
          <input id="ext-last" class="input" bind:value={form.LastName} placeholder="Smith" />
        </div>

        <div class="sm:col-span-2">
          <label for="ext-email" class="label">User Email Address</label>
          <input id="ext-email" type="email" class="input" bind:value={form.Email} placeholder="john.smith@company.com" />
          <p class="text-[11px] text-slate-500 mt-1">Used for voicemail-to-email notifications and operator invitations.</p>
        </div>
      </div>

      <!-- SIP Password Generator -->
      <div class="border-t border-slate-100 pt-3">
        <div class="flex items-center justify-between mb-1">
          <label for="ext-pass" class="label mb-0">SIP Registration Password {editingId ? '(leave blank to keep current)' : '*'}</label>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="text-xs text-sky-600 hover:text-sky-800 font-medium"
              on:click={regeneratePassword}
            >
              ↻ Regenerate
            </button>
            {#if form.Password}
              <button
                type="button"
                class="text-xs text-slate-600 hover:text-slate-800 font-medium"
                on:click={copyPassword}
              >
                📋 Copy
              </button>
              <button
                type="button"
                class="text-xs text-slate-600 hover:text-slate-800 font-medium"
                on:click={() => showPassword = !showPassword}
              >
                {showPassword ? 'Hide' : 'Show'}
              </button>
            {/if}
          </div>
        </div>

        <div class="relative mt-1">
          <input
            id="ext-pass"
            type={showPassword ? 'text' : 'password'}
            class="input mono pr-10"
            bind:value={form.Password}
            placeholder="SIP registration secret"
            required={!editingId}
          />
        </div>
        <p class="text-[11px] text-slate-500 mt-1">
          Strong cryptographic password automatically pre-generated for hardware desk phone or softphone registration.
        </p>
      </div>

      {#if form.Password}
        <div>
          <label for="ext-confirm" class="label">Confirm SIP Password *</label>
          <input
            id="ext-confirm"
            type={showPassword ? 'text' : 'password'}
            class="input mono"
            bind:value={form.ConfirmPassword}
            placeholder="Re-enter SIP secret"
            required
          />
        </div>
      {/if}

      <div class="flex items-center space-x-2 pt-2">
        <input id="ext-active" type="checkbox" bind:checked={form.IsActive} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
        <label for="ext-active" class="text-xs font-semibold text-slate-800">Extension is Active and Allowed to Register</label>
      </div>
    {/if}

    <!-- TAB 2: CALLER ID & DID MAPPING -->
    {#if activeTab === 'cid'}
      <div class="space-y-4">
        <div>
          <label for="ext-cid-name" class="label">Outbound Caller ID Name (CNAM)</label>
          <input
            id="ext-cid-name"
            class="input"
            bind:value={form.CallerIDName}
            placeholder="e.g. John Smith or Acme Sales"
          />
          <p class="text-[11px] text-slate-500 mt-1">Display name shown when this extension calls internal peers or external lines.</p>
        </div>

        <div>
          <label for="ext-cid-num" class="label">Internal Caller ID Number</label>
          <input
            id="ext-cid-num"
            class="input mono"
            bind:value={form.CallerIDNumber}
            placeholder={form.ExtensionNumber || '1001'}
          />
          <p class="text-[11px] text-slate-500 mt-1">Number shown on internal phones when dialing colleagues (defaults to extension number).</p>
        </div>

        <div class="p-3 bg-sky-50 border border-sky-100 rounded-lg space-y-2">
          <label for="ext-outbound-did" class="label font-bold text-sky-900">Assigned Outbound Direct Inward Dial (DID / CLI)</label>
          <div class="flex gap-2">
            <input
              id="ext-outbound-did"
              class="input mono font-semibold bg-white"
              bind:value={form.OutboundDID}
              placeholder="e.g. +14155550199 or 02086884243"
            />
            {#if availableDIDs.length > 0}
              <select
                class="input !w-48 text-xs bg-white"
                on:change={(e) => {
                  const target = e.target as HTMLSelectElement;
                  if (target.value) form.OutboundDID = target.value;
                }}
              >
                <option value="">Select DID...</option>
                {#each availableDIDs as d}
                  <option value={d.Number}>{d.Number} ({d.Description})</option>
                {/each}
              </select>
            {/if}
          </div>
          <p class="text-[11px] text-sky-700 leading-relaxed">
            The external telephone number presented when this extension calls out over public SIP trunks. Callee can ring this number directly back to reach this user.
          </p>
        </div>
      </div>
    {/if}

    <!-- TAB 3: FORWARDING & MOBILE -->
    {#if activeTab === 'forwarding'}
      <div class="space-y-4">
        <div>
          <label for="ext-mobile" class="label">Mobile Telephone Number</label>
          <input
            id="ext-mobile"
            type="tel"
            class="input mono"
            bind:value={form.MobileNumber}
            placeholder="e.g. +14155550188 or 07340221143"
          />
          <p class="text-[11px] text-slate-500 mt-1">User's mobile cell phone for twin-ringing or diversion rules.</p>
        </div>

        <div>
          <label for="ext-timeout" class="label">Handset Ring Timeout (Seconds)</label>
          <input
            id="ext-timeout"
            type="number"
            min="5"
            max="120"
            class="input w-32 mono font-bold"
            bind:value={form.RingTimeoutSeconds}
          />
          <p class="text-[11px] text-slate-500 mt-1">Seconds to ring this extension's physical desk phone before triggering No-Answer forwarding or Voicemail.</p>
        </div>

        <div class="border-t border-slate-200 pt-3 space-y-3">
          <div class="text-xs font-bold text-slate-800 uppercase tracking-wider">Automated Call Diversion Rules</div>

          <!-- Always Forward -->
          <div class="p-3 bg-white border border-slate-200 rounded-lg space-y-2">
            <div class="flex items-center justify-between">
              <label for="fwd-always" class="text-xs font-semibold text-slate-800 flex items-center space-x-2">
                <input id="fwd-always" type="checkbox" bind:checked={form.ForwardAlwaysEnabled} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
                <span>Call Forward Always (Unconditional)</span>
              </label>
              {#if form.ForwardAlwaysEnabled}
                <span class="text-[10px] uppercase font-bold text-amber-600 bg-amber-50 px-2 py-0.5 rounded border border-amber-200">Active</span>
              {/if}
            </div>
            {#if form.ForwardAlwaysEnabled}
              <input class="input mono text-xs mt-1" bind:value={form.ForwardAlwaysTarget} placeholder="Destination extension or mobile number..." />
            {/if}
          </div>

          <!-- Busy Forward -->
          <div class="p-3 bg-white border border-slate-200 rounded-lg space-y-2">
            <label for="fwd-busy" class="text-xs font-semibold text-slate-800 flex items-center space-x-2">
              <input id="fwd-busy" type="checkbox" bind:checked={form.ForwardBusyEnabled} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
              <span>Forward When Busy (On another call)</span>
            </label>
            {#if form.ForwardBusyEnabled}
              <input class="input mono text-xs mt-1" bind:value={form.ForwardBusyTarget} placeholder="Destination extension or mobile number..." />
            {/if}
          </div>

          <!-- No Answer Forward -->
          <div class="p-3 bg-white border border-slate-200 rounded-lg space-y-2">
            <label for="fwd-noans" class="text-xs font-semibold text-slate-800 flex items-center space-x-2">
              <input id="fwd-noans" type="checkbox" bind:checked={form.ForwardNoAnswerEnabled} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
              <span>Forward When No Answer (Timeout)</span>
            </label>
            {#if form.ForwardNoAnswerEnabled}
              <input class="input mono text-xs mt-1" bind:value={form.ForwardNoAnswerTarget} placeholder="Destination (e.g. mobile number, colleague extension, or leave empty for voicemail)..." />
            {/if}
          </div>
        </div>
      </div>
    {/if}

    <!-- TAB 4: VOICEMAIL PREFERENCES -->
    {#if activeTab === 'voicemail'}
      <div class="space-y-4">
        <div class="flex items-center space-x-2">
          <input id="ext-vm-enable" type="checkbox" bind:checked={form.VoicemailEnabled} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
          <label for="ext-vm-enable" class="text-xs font-semibold text-slate-800">Enable Personal Voicemail Box for this Extension</label>
        </div>

        {#if form.VoicemailEnabled}
          <div class="p-4 bg-slate-50 border border-slate-200 rounded-lg space-y-3">
            <div>
              <label for="ext-vm-pin" class="label">Voicemail Access PIN *</label>
              <input
                id="ext-vm-pin"
                type="password"
                maxlength="10"
                class="input mono w-40 font-bold"
                bind:value={form.VoicemailPIN}
                placeholder={form.ExtensionNumber || '1001'}
              />
              <p class="text-[11px] text-slate-500 mt-1">4-digit PIN dialed after dialing <span class="font-bold">*97</span> to listen to voice messages.</p>
            </div>

            <div class="flex items-center space-x-2 pt-2 border-t border-slate-200">
              <input id="ext-vm-email" type="checkbox" bind:checked={form.VoicemailToEmail} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
              <label for="ext-vm-email" class="text-xs font-medium text-slate-700">
                Email audio recordings (.wav) to user's email address ({form.Email || 'No email specified'})
              </label>
            </div>
          </div>
        {/if}
      </div>
    {/if}

    <!-- Modal Footer -->
    <div class="flex justify-between items-center pt-3 border-t border-slate-200">
      <div class="text-xs text-slate-400">
        {#if activeTab !== 'basic'}
          <button type="button" class="text-sky-600 hover:underline" on:click={() => activeTab = 'basic'}>&larr; Back to Basic Settings</button>
        {/if}
      </div>
      <div class="flex gap-2">
        <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
        <Button type="submit" disabled={saving}>{saving ? 'Saving…' : (editingId ? 'Update Extension' : 'Create Extension')}</Button>
      </div>
    </div>
  </form>
</Modal>
