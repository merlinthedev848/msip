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
    AuthMethod: string;
    Username: string;
    Register: boolean;
  }

  let trunks: TrunkItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    ProviderName: '',
    SIPServer: '',
    AuthMethod: 'IP_ACL',
    Username: '',
    Password: '',
    Register: false
  };

  $: filtered = trunks.filter(t => 
    !search.trim() || `${t.ProviderName} ${t.SIPServer} ${t.Username}`.toLowerCase().includes(search.trim().toLowerCase())
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
    form = { ProviderName: '', SIPServer: '', AuthMethod: 'IP_ACL', Username: '', Password: '', Register: false };
    modalOpen = true;
  }

  function openEdit(t: TrunkItem) {
    editingId = t.ID;
    form = {
      ProviderName: t.ProviderName,
      SIPServer: t.SIPServer,
      AuthMethod: t.AuthMethod || 'IP_ACL',
      Username: t.Username || '',
      Password: '',
      Register: t.Register ?? false
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.ProviderName.trim()) return toast.error('Provider name is required');
    if (!form.SIPServer.trim()) return toast.error('SIP Server is required');

    saving = true;
    try {
      const payload: any = {
        ProviderName: form.ProviderName.trim(),
        SIPServer: form.SIPServer.trim(),
        AuthMethod: form.AuthMethod,
        Username: form.Username.trim(),
        Register: form.Register
      };
      if (form.Password) {
        payload.Password = form.Password;
      }

      if (editingId) {
        await api.put(`/trunks/${editingId}`, payload);
        toast.success('Trunk updated');
      } else {
        await api.post('/trunks', payload);
        toast.success('Trunk created');
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
    if (!confirm(`Delete SIP trunk ${t.ProviderName}?`)) return;
    try {
      await api.del(`/trunks/${t.ID}`);
      toast.success('Trunk deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="SIP Trunks" description="Manage carrier gateways and external VoIP interconnects.">
  <Button on:click={openCreate}>+ New Trunk</Button>
</PageHeader>

<Card padded={false} title="Configured Gateways" subtitle="{trunks.length} trunks">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search trunk…" bind:value={search} aria-label="Search trunks" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No trunks match your search' : 'No SIP trunks configured'}
      message={search ? '' : 'Add an outbound/inbound SIP trunk to route PSTN calls through your carrier.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first trunk</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Provider</th>
            <th>SIP Server</th>
            <th>Auth Method</th>
            <th>Username</th>
            <th>Registration</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as t (t.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{t.ProviderName}</td>
              <td class="mono text-xs text-slate-700">{t.SIPServer}</td>
              <td><Badge tone="info">{t.AuthMethod}</Badge></td>
              <td class="mono text-xs text-slate-600">{t.Username || '—'}</td>
              <td>
                <Badge tone={t.Register ? 'success' : 'neutral'}>
                  {t.Register ? 'Registered' : 'IP-Auth'}
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Trunk' : 'Add SIP Trunk'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="trunk-name" class="label">Provider Name</label>
      <input id="trunk-name" class="input" bind:value={form.ProviderName} placeholder="e.g. Twilio Elastic SIP" required />
    </div>

    <div>
      <label for="trunk-server" class="label">SIP Server / Host</label>
      <input id="trunk-server" class="input mono" bind:value={form.SIPServer} placeholder="e.g. sip.carrier.com" required />
    </div>

    <div>
      <label for="trunk-auth" class="label">Authentication Method</label>
      <select id="trunk-auth" class="select" bind:value={form.AuthMethod}>
        <option value="IP_ACL">IP ACL / Whitelist</option>
        <option value="USERPASS">Username & Password</option>
      </select>
    </div>

    {#if form.AuthMethod === 'USERPASS'}
      <div>
        <label for="trunk-user" class="label">Username</label>
        <input id="trunk-user" class="input mono" bind:value={form.Username} placeholder="SIP Username" />
      </div>

      <div>
        <label for="trunk-pass" class="label">Password {editingId ? '(leave blank to keep current)' : ''}</label>
        <input id="trunk-pass" type="password" class="input mono" bind:value={form.Password} placeholder="SIP Password" />
      </div>

      <div class="flex items-center space-x-2 pt-1">
        <input id="trunk-reg" type="checkbox" bind:checked={form.Register} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
        <label for="trunk-reg" class="text-xs font-medium text-slate-700">Send SIP REGISTER packets to gateway</label>
      </div>
    {/if}

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
