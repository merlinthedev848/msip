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
    CreatedAt: string;
  }

  let extensions: ExtensionItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    ExtensionNumber: '',
    Password: '',
    ConfirmPassword: '',
    IsActive: true
  };

  $: filtered = extensions.filter(e => 
    !search.trim() || `${e.ExtensionNumber} ${e.ProvisioningCode}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ extensions: ExtensionItem[] }>('/extensions');
      extensions = data.extensions || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { ExtensionNumber: '', Password: '', ConfirmPassword: '', IsActive: true };
    modalOpen = true;
  }

  function openEdit(ext: ExtensionItem) {
    editingId = ext.ID;
    form = { ExtensionNumber: ext.ExtensionNumber, Password: '', ConfirmPassword: '', IsActive: ext.IsActive };
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
        IsActive: form.IsActive
      };
      if (form.Password) {
        payload.PasswordHash = form.Password;
      }

      if (editingId) {
        await api.put(`/extensions/${editingId}`, payload);
        toast.success('Extension updated');
      } else {
        await api.post('/extensions', payload);
        toast.success('Extension created');
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
    if (!confirm(`Delete extension ${ext.ExtensionNumber}?`)) return;
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

<PageHeader title="Extensions" description="Manage SIP user endpoints and provision codes.">
  <Button on:click={openCreate}>+ New Extension</Button>
</PageHeader>

<Card padded={false} title="Configured Extensions" subtitle="{extensions.length} extensions">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search extension…" bind:value={search} aria-label="Search extensions" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No extensions match your search' : 'No extensions found'}
      message={search ? '' : 'Create your first extension to connect a SIP hardware device or softphone.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first extension</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Extension</th>
            <th>Status</th>
            <th>Provisioning Code</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as ext (ext.ID)}
            <tr>
              <td class="mono font-bold text-slate-800 text-sm">{ext.ExtensionNumber}</td>
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Extension' : 'Create Extension'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="ext-number" class="label">Extension Number</label>
      <input id="ext-number" class="input mono" bind:value={form.ExtensionNumber} placeholder="e.g. 1001" required />
    </div>

    <div>
      <label for="ext-pass" class="label">SIP Password {editingId ? '(leave blank to keep current)' : ''}</label>
      <input id="ext-pass" type="password" class="input mono" bind:value={form.Password} placeholder="SIP registration secret" required={!editingId} />
    </div>

    {#if form.Password}
      <div>
        <label for="ext-confirm" class="label">Confirm Password</label>
        <input id="ext-confirm" type="password" class="input mono" bind:value={form.ConfirmPassword} placeholder="Confirm password" required />
      </div>
    {/if}

    <div class="flex items-center space-x-2 pt-1">
      <input id="ext-active" type="checkbox" bind:checked={form.IsActive} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
      <label for="ext-active" class="text-xs font-medium text-slate-700">Active</label>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
