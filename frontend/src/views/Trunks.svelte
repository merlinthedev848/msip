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

  interface TrunkItem {
    ID: string;
    ProviderName: string;
    SIPServer: string;
    SIPPort: number;
    Transport: string;
    AuthMethod: string;
    Username: string;
    Register: boolean;
    MaxChannels: number;
    Codecs: string;
    FromDomain: string;
    OutboundProxy: string;
    CallerIDNumber: string;
    DTMFMode?: string;
    QualifyFrequency?: number;
    RegistrationExpiry?: number;
    DirectMedia?: boolean;
    IsActive: boolean;
  }

  let trunks: TrunkItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let activeTab: 'general' | 'auth' | 'advanced' = 'general';

  let form = {
    ProviderName: '',
    SIPServer: '',
    SIPPort: 5060,
    Transport: 'UDP',
    AuthMethod: 'USERPASS',
    Username: '',
    Password: '',
    Register: true,
    MaxChannels: 10,
    Codecs: 'G.711u, G.711a, G.722, Opus',
    FromDomain: '',
    OutboundProxy: '',
    CallerIDNumber: '',
    DTMFMode: 'rfc2833',
    QualifyFrequency: 60,
    RegistrationExpiry: 3600,
    DirectMedia: false,
    IsActive: true
  };

  $: filtered = trunks.filter(t => 
    !search.trim() || `${t.ProviderName} ${t.SIPServer} ${t.Username} ${t.CallerIDNumber}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ trunks: TrunkItem[] }>('/trunks');
      trunks = data.trunks || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    activeTab = 'general';
    form = {
      ProviderName: '',
      SIPServer: '',
      SIPPort: 5060,
      Transport: 'UDP',
      AuthMethod: 'USERPASS',
      Username: '',
      Password: '',
      Register: true,
      MaxChannels: 10,
      Codecs: 'G.711u, G.711a, G.722, Opus',
      FromDomain: '',
      OutboundProxy: '',
      CallerIDNumber: '',
      DTMFMode: 'rfc2833',
      QualifyFrequency: 60,
      RegistrationExpiry: 3600,
      DirectMedia: false,
      IsActive: true
    };
    modalOpen = true;
  }

  function openEdit(t: TrunkItem) {
    editingId = t.ID;
    activeTab = 'general';
    form = {
      ProviderName: t.ProviderName,
      SIPServer: t.SIPServer,
      SIPPort: t.SIPPort || 5060,
      Transport: t.Transport || 'UDP',
      AuthMethod: t.AuthMethod || 'USERPASS',
      Username: t.Username || '',
      Password: '',
      Register: t.Register ?? false,
      MaxChannels: t.MaxChannels || 10,
      Codecs: t.Codecs || 'G.711u, G.711a, G.722, Opus',
      FromDomain: t.FromDomain || '',
      OutboundProxy: t.OutboundProxy || '',
      CallerIDNumber: t.CallerIDNumber || '',
      DTMFMode: t.DTMFMode || 'rfc2833',
      QualifyFrequency: t.QualifyFrequency ?? 60,
      RegistrationExpiry: t.RegistrationExpiry ?? 3600,
      DirectMedia: t.DirectMedia ?? false,
      IsActive: t.IsActive ?? true
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.ProviderName.trim()) return toast.error('Trunk provider name is required');
    if (!form.SIPServer.trim()) return toast.error('SIP server address is required');

    saving = true;
    try {
      const payload: any = {
        ProviderName: form.ProviderName.trim(),
        SIPServer: form.SIPServer.trim(),
        SIPPort: Number(form.SIPPort) || 5060,
        Transport: form.Transport,
        AuthMethod: form.AuthMethod,
        Username: form.Username.trim(),
        Register: form.Register,
        MaxChannels: Number(form.MaxChannels) || 10,
        Codecs: form.Codecs.trim(),
        FromDomain: form.FromDomain.trim() || form.SIPServer.trim(),
        OutboundProxy: form.OutboundProxy.trim(),
        CallerIDNumber: form.CallerIDNumber.trim(),
        DTMFMode: form.DTMFMode,
        QualifyFrequency: Number(form.QualifyFrequency) || 0,
        RegistrationExpiry: Number(form.RegistrationExpiry) || 3600,
        DirectMedia: form.DirectMedia,
        IsActive: form.IsActive
      };

      if (form.Password) {
        payload.Password = form.Password;
      }

      if (editingId) {
        await api.put(`/trunks/${editingId}`, payload);
        toast.success(`Trunk ${form.ProviderName} updated`);
      } else {
        await api.post('/trunks', payload);
        toast.success(`Trunk ${form.ProviderName} configured`);
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(t: TrunkItem) {
    if (!confirm(`Permanently remove SIP Trunk "${t.ProviderName}"? Active calls will drop.`)) return;
    try {
      await api.del(`/trunks/${t.ID}`);
      toast.success('SIP trunk removed');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="SIP Trunks & Gateways" description="Connect telecom carriers (Twilio, Telnyx, VoIP.ms, Bandwidth) for outbound and inbound calls.">
  <Button on:click={openCreate}>+ New SIP Trunk</Button>
</PageHeader>

<Card padded={false} title="Configured Carrier Trunks" subtitle="{trunks.length} SIP gateways">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-60" placeholder="Search provider, host, caller ID…" bind:value={search} aria-label="Search trunks" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No trunks match your search' : 'No SIP trunks configured'}
      message={search ? '' : 'Connect your telecom carrier gateway to allow extensions to dial external telephone lines.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add First SIP Trunk</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Carrier / Provider</th>
            <th>SIP Host & Port</th>
            <th>Auth Mode</th>
            <th>Registered Username</th>
            <th>Trunk CLI (Default Caller ID)</th>
            <th>Channels</th>
            <th>Status</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as t (t.ID)}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="font-bold text-slate-800 text-sm">
                <div>{t.ProviderName}</div>
                <div class="text-[10px] text-slate-400 font-mono font-normal">{t.Transport || 'UDP'}</div>
              </td>
              <td class="mono text-xs font-semibold text-slate-700">
                {t.SIPServer}:{t.SIPPort || 5060}
              </td>
              <td>
                <span class="mono px-2 py-0.5 rounded text-xs border {t.AuthMethod === 'IP_ACL' ? 'bg-amber-50 text-amber-800 border-amber-200' : 'bg-sky-50 text-sky-800 border-sky-200'}">
                  {t.AuthMethod === 'IP_ACL' ? 'IP Whitelist' : (t.Register ? 'SIP Registration' : 'Credentials')}
                </span>
              </td>
              <td class="mono text-xs text-slate-600">{t.Username || '—'}</td>
              <td class="mono text-xs font-bold text-slate-800">
                {t.CallerIDNumber || 'Pass-Through'}
              </td>
              <td class="text-xs text-slate-600">{t.MaxChannels || 10} Lines</td>
              <td>
                <Badge tone={t.IsActive ? 'success' : 'neutral'}>
                  {t.IsActive ? 'Online' : 'Disabled'}
                </Badge>
              </td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(t)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(t)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? `Edit Trunk: ${form.ProviderName}` : 'Add New SIP Trunk'} maxWidth="max-w-2xl">
  <!-- Tabs -->
  <div class="flex items-center space-x-1 border-b border-slate-200 mb-4 pb-1">
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'general' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'general'}
    >
      General & Host
    </button>
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'auth' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'auth'}
    >
      Authentication & Credentials
    </button>
    <button
      type="button"
      class="px-3 py-1.5 text-xs font-semibold rounded-md transition-colors {activeTab === 'advanced' ? 'bg-[#0284C7] text-white' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
      on:click={() => activeTab = 'advanced'}
    >
      Codecs & Capacity
    </button>
  </div>

  <form on:submit={save} class="space-y-4">
    {#if activeTab === 'general'}
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="sm:col-span-2">
          <label for="trunk-name" class="label">Provider / Gateway Name *</label>
          <input id="trunk-name" class="input font-bold" bind:value={form.ProviderName} placeholder="e.g. Twilio Primary, Telnyx East, VoIP.ms" required />
        </div>

        <div>
          <label for="trunk-host" class="label">SIP Server Host / FQDN *</label>
          <input id="trunk-host" class="input mono font-semibold" bind:value={form.SIPServer} placeholder="e.g. sip.twilio.com or 192.168.1.10" required />
        </div>

        <div>
          <label for="trunk-port" class="label">SIP Port</label>
          <input id="trunk-port" type="number" class="input mono font-semibold" bind:value={form.SIPPort} placeholder="5060" />
        </div>

        <div>
          <label for="trunk-transport" class="label">Transport Protocol</label>
          <select id="trunk-transport" class="input" bind:value={form.Transport}>
            <option value="UDP">UDP (Standard Telephony)</option>
            <option value="TCP">TCP</option>
            <option value="TLS">TLS (Encrypted SIP)</option>
          </select>
        </div>

        <div>
          <label for="trunk-cli" class="label">Default Outbound Caller ID (CLI)</label>
          <input id="trunk-cli" class="input mono" bind:value={form.CallerIDNumber} placeholder="e.g. +14155550199 or 02086884243" />
        </div>

        <div class="sm:col-span-2 flex items-center space-x-2 pt-2">
          <input id="trunk-active" type="checkbox" bind:checked={form.IsActive} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
          <label for="trunk-active" class="text-xs font-semibold text-slate-800">Trunk is Active and Ready for Dialing</label>
        </div>
      </div>
    {/if}

    {#if activeTab === 'auth'}
      <div class="space-y-4">
        <div>
          <label for="trunk-auth" class="label">Authentication Mode</label>
          <select id="trunk-auth" class="input" bind:value={form.AuthMethod}>
            <option value="USERPASS">Registration & Digest Credentials (Username & Password)</option>
            <option value="IP_ACL">IP-Based Authentication (Carrier Whitelist / No Password)</option>
          </select>
        </div>

        {#if form.AuthMethod === 'USERPASS'}
          <div class="p-4 bg-slate-50 border border-slate-200 rounded-lg space-y-3">
            <div>
              <label for="trunk-user" class="label">Authentication Username / Account ID</label>
              <input id="trunk-user" class="input mono" bind:value={form.Username} placeholder="SIP username or trunk ID" />
            </div>

            <div>
              <label for="trunk-pass" class="label">SIP Secret / Password {editingId ? '(leave blank to keep current)' : ''}</label>
              <input id="trunk-pass" type="password" class="input mono" bind:value={form.Password} placeholder="Carrier trunk password" />
            </div>

            <div class="flex items-center space-x-2 pt-1">
              <input id="trunk-register" type="checkbox" bind:checked={form.Register} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
              <label for="trunk-register" class="text-xs font-semibold text-slate-700">Send Periodic SIP REGISTER requests to provider</label>
            </div>
          </div>
        {:else}
          <div class="p-4 bg-sky-50 border border-sky-200 rounded-lg text-xs text-sky-900 space-y-1">
            <span class="font-bold block">IP Whitelist / Peer Mode</span>
            <p>Calls will be transmitted directly to the SIP host with no digest challenge. Ensure your PBX public IP is whitelisted in your carrier portal.</p>
          </div>
        {/if}
      </div>
    {/if}

    {#if activeTab === 'advanced'}
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label for="trunk-max" class="label">Concurrent Channel Limit</label>
          <input id="trunk-max" type="number" min="1" max="1000" class="input mono font-bold" bind:value={form.MaxChannels} />
          <p class="text-[11px] text-slate-500 mt-1">Maximum simultaneous calls allowed through this trunk.</p>
        </div>

        <div>
          <label for="trunk-codecs" class="label">Permitted Audio Codecs</label>
          <input id="trunk-codecs" class="input mono" bind:value={form.Codecs} placeholder="G.711u, G.711a, G.722, Opus" />
          <p class="text-[11px] text-slate-500 mt-1">Comma-separated RTP codecs in order of priority.</p>
        </div>

        <div>
          <label for="trunk-dtmf" class="label">DTMF Signaling Mode</label>
          <select id="trunk-dtmf" class="input" bind:value={form.DTMFMode}>
            <option value="rfc2833">RFC 2833 / RFC 4733 (Out-of-band RTP Events - Recommended)</option>
            <option value="info">SIP INFO (Application/DTMF-Relay)</option>
            <option value="inband">Inband Audio (Tone within audio stream)</option>
          </select>
          <p class="text-[11px] text-slate-500 mt-1">Tone delivery protocol for IVRs and carrier touch-tones.</p>
        </div>

        <div>
          <label for="trunk-qualify" class="label">Qualify Keepalive Interval (seconds)</label>
          <input id="trunk-qualify" type="number" min="0" max="600" class="input mono" bind:value={form.QualifyFrequency} placeholder="60" />
          <p class="text-[11px] text-slate-500 mt-1">Frequency of SIP OPTIONS keepalive pings (0 to disable).</p>
        </div>

        <div>
          <label for="trunk-expiry" class="label">Registration Expiry (seconds)</label>
          <input id="trunk-expiry" type="number" min="60" max="86400" class="input mono" bind:value={form.RegistrationExpiry} placeholder="3600" />
          <p class="text-[11px] text-slate-500 mt-1">SIP Registration TTL requested from upstream carrier.</p>
        </div>

        <div class="sm:col-span-2 flex items-center space-x-2 pt-2 border-t border-slate-200">
          <input id="trunk-direct-media" type="checkbox" bind:checked={form.DirectMedia} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
          <div>
            <label for="trunk-direct-media" class="text-xs font-semibold text-slate-800 cursor-pointer">Direct Media / RTP Audio Bypass (Re-INVITE)</label>
            <p class="text-[11px] text-slate-500">Allows audio stream to travel directly between carrier and endpoint, reducing PBX server load.</p>
          </div>
        </div>
      </div>
    {/if}

    <div class="flex justify-between items-center pt-3 border-t border-slate-200">
      <div class="text-xs text-slate-400">
        {#if activeTab !== 'general'}
          <button type="button" class="text-sky-600 hover:underline" on:click={() => activeTab = 'general'}>&larr; General Settings</button>
        {/if}
      </div>
      <div class="flex gap-2">
        <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
        <Button type="submit" disabled={saving}>{saving ? 'Saving…' : (editingId ? 'Update Trunk' : 'Save Trunk')}</Button>
      </div>
    </div>
  </form>
</Modal>
