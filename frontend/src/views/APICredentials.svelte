<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface CredentialItem {
    ID: string;
    Name: string;
    TokenKey: string;
    CreatedAt: string;
  }

  let credentials: CredentialItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let form = { Name: '' };
  let createdSecret = '';
  let createdTokenKey = '';

  $: filtered = credentials.filter(c => 
    !search.trim() || `${c.Name} ${c.TokenKey}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ credentials: CredentialItem[] }>('/credentials');
      credentials = data.credentials || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    form = { Name: '' };
    createdSecret = '';
    createdTokenKey = '';
    modalOpen = true;
  }

  async function createKey(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Application name is required');

    saving = true;
    try {
      const res = await api.post<any>('/credentials', { Name: form.Name.trim() });
      createdSecret = res.Secret || '';
      createdTokenKey = res.TokenKey || '';
      toast.success('API Key generated');
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(c: CredentialItem) {
    if (!confirm(`Revoke API key "${c.Name}"? Applications using this token will stop working immediately.`)) return;
    try {
      await api.del(`/credentials/${c.ID}`);
      toast.success('API key revoked');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="API Credentials" description="Bearer tokens and webhook authorization keys for external telecom integrations.">
  <Button on:click={openCreate}>+ Generate Key</Button>
</PageHeader>

<Card padded={false} title="Active Application Keys" subtitle="{credentials.length} keys issued">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search application…" bind:value={search} aria-label="Search keys" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No API keys match your search' : 'No API keys issued'}
      message={search ? '' : 'Issue an API key to allow external software to access the PBX REST API.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Generate first key</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Application</th>
            <th>API Token Key</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as c (c.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{c.Name}</td>
              <td class="mono font-semibold text-xs text-sky-700">{c.TokenKey}</td>
              <td class="text-right whitespace-nowrap">
                <button type="button" class="link-danger" on:click={() => remove(c)}>Revoke</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={createdSecret ? 'API Key Generated' : 'Generate API Key'}>
  {#if createdSecret}
    <div class="space-y-4">
      <div class="p-3 bg-amber-50 border border-amber-200 rounded-md text-amber-900 text-xs leading-relaxed">
        <strong>Important:</strong> Copy your API secret key now. For security purposes, this secret will not be displayed again.
      </div>

      <div>
        <label for="created-token" class="label">API Key</label>
        <input id="created-token" class="input mono" readonly value={createdTokenKey} />
      </div>

      <div>
        <label for="created-secret" class="label">API Secret</label>
        <input id="created-secret" class="input mono font-bold text-sky-700" readonly value={createdSecret} />
      </div>

      <div class="flex justify-end pt-2">
        <Button on:click={() => modalOpen = false}>Done</Button>
      </div>
    </div>
  {:else}
    <form on:submit={createKey} class="space-y-4">
      <div>
        <label for="cred-name" class="label">Application Name</label>
        <input id="cred-name" class="input" bind:value={form.Name} placeholder="e.g. CRM Integration or Billing Sync" required />
      </div>

      <div class="flex justify-end gap-2 pt-2">
        <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
        <Button type="submit" disabled={saving}>{saving ? 'Generating…' : 'Generate'}</Button>
      </div>
    </form>
  {/if}
</Modal>
