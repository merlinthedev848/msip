<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, getCurrentUser } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface UserItem {
    ID: string;
    Email: string;
    Role: string;
    CreatedAt: string;
  }

  let users: UserItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';
  let currentUser = getCurrentUser();

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Email: '',
    Password: '',
    Role: 'admin'
  };

  $: filtered = users.filter(u => 
    !search.trim() || `${u.Email} ${u.Role}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ users: UserItem[] }>('/users');
      users = data.users || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = { Email: '', Password: '', Role: 'admin' };
    modalOpen = true;
  }

  function openEdit(u: UserItem) {
    editingId = u.ID;
    form = { Email: u.Email, Password: '', Role: u.Role || 'admin' };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Email.trim()) return toast.error('Email address is required');
    if (!editingId && !form.Password) return toast.error('Password is required');

    saving = true;
    try {
      const payload: any = {
        Email: form.Email.trim().toLowerCase(),
        Role: form.Role
      };
      if (form.Password) {
        payload.PasswordHash = form.Password;
      }

      if (editingId) {
        await api.put(`/users/${editingId}`, payload);
        toast.success('Admin updated');
      } else {
        await api.post('/users', payload);
        toast.success('Admin invited');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(u: UserItem) {
    if (u.Email.toLowerCase() === (currentUser.email || '').toLowerCase()) {
      return toast.error('You cannot delete your own account while logged in');
    }
    if (!confirm(`Revoke dashboard access for ${u.Email}?`)) return;
    try {
      await api.del(`/users/${u.ID}`);
      toast.success('User revoked');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Admins & Roles" description="Manage dashboard administrators and RBAC permissions.">
  <Button on:click={openCreate}>+ Invite User</Button>
</PageHeader>

<Card padded={false} title="Active Administrators" subtitle="{users.length} accounts">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search email or role…" bind:value={search} aria-label="Search users" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No users match your search' : 'No administrator users'}
      message={search ? '' : 'Invite administrative team members to manage PBX extensions and routing.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Invite first user</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Email Address</th>
            <th>Role</th>
            <th>Created</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as u (u.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">
                {u.Email}
                {#if u.Email.toLowerCase() === (currentUser.email || '').toLowerCase()}
                  <span class="ml-1.5 text-xs text-sky-600 font-semibold">(You)</span>
                {/if}
              </td>
              <td>
                <Badge tone={u.Role === 'superadmin' ? 'warning' : (u.Role === 'admin' ? 'info' : 'neutral')}>
                  {u.Role.toUpperCase()}
                </Badge>
              </td>
              <td class="text-xs text-slate-500">{formatDateTime(u.CreatedAt)}</td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(u)}>Edit</button>
                {#if u.Email.toLowerCase() !== (currentUser.email || '').toLowerCase()}
                  <button type="button" class="link-danger" on:click={() => remove(u)}>Revoke</button>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Admin' : 'Invite New Admin'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="adm-email" class="label">Email Address</label>
      <input id="adm-email" type="email" class="input" bind:value={form.Email} placeholder="name@company.com" required />
    </div>

    <div>
      <label for="adm-pass" class="label">Password {editingId ? '(leave blank to keep current)' : ''}</label>
      <input id="adm-pass" type="password" class="input mono" bind:value={form.Password} placeholder="Dashboard password" required={!editingId} />
    </div>

    <div>
      <label for="adm-role" class="label">Permission Role</label>
      <select id="adm-role" class="select" bind:value={form.Role}>
        <option value="superadmin">SUPERADMIN (Full Tenant & Cluster Control)</option>
        <option value="admin">ADMIN (PBX Configuration & User Management)</option>
        <option value="viewer">VIEWER (Read-Only Access)</option>
      </select>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
