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

  interface Contact {
    ID: string;
    FirstName: string;
    LastName: string;
    Company: string;
    PhoneNumber: string;
    MobileNumber: string;
    Email: string;
    SpeedDial: string;
    Group: string;
  }

  let contacts: Contact[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';
  let selectedGroup = 'All';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    FirstName: '',
    LastName: '',
    Company: '',
    PhoneNumber: '',
    MobileNumber: '',
    Email: '',
    SpeedDial: '',
    Group: 'General'
  };

  $: groups = ['All', ...Array.from(new Set(contacts.map(c => c.Group || 'General')))];

  $: filteredContacts = contacts.filter(c => {
    const matchesSearch = !search.trim() || 
      `${c.FirstName} ${c.LastName} ${c.Company} ${c.PhoneNumber} ${c.MobileNumber} ${c.SpeedDial}`
        .toLowerCase()
        .includes(search.trim().toLowerCase());
    const matchesGroup = selectedGroup === 'All' || (c.Group || 'General') === selectedGroup;
    return matchesSearch && matchesGroup;
  });

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const res = await api.get<{ contacts: Contact[] }>('/contacts');
      contacts = res.contacts || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      FirstName: '',
      LastName: '',
      Company: '',
      PhoneNumber: '',
      MobileNumber: '',
      Email: '',
      SpeedDial: '',
      Group: 'General'
    };
    modalOpen = true;
  }

  function openEdit(c: Contact) {
    editingId = c.ID;
    form = {
      FirstName: c.FirstName,
      LastName: c.LastName || '',
      Company: c.Company || '',
      PhoneNumber: c.PhoneNumber,
      MobileNumber: c.MobileNumber || '',
      Email: c.Email || '',
      SpeedDial: c.SpeedDial || '',
      Group: c.Group || 'General'
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.FirstName.trim()) return toast.error('First name is required');
    if (!form.PhoneNumber.trim()) return toast.error('Phone number is required');

    saving = true;
    try {
      if (editingId) {
        await api.put(`/contacts/${editingId}`, form);
        toast.success('Contact updated');
      } else {
        await api.post('/contacts', form);
        toast.success('Contact added');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteContact(id: string, name: string) {
    if (!confirm(`Delete contact "${name}"?`)) return;
    try {
      await api.del(`/contacts/${id}`);
      toast.success('Contact removed');
      contacts = contacts.filter(c => c.ID !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  function copyXmlUrl() {
    const url = `${window.location.origin}/api/v1/contacts/phonebook.xml`;
    navigator.clipboard.writeText(url);
    toast.success('Phonebook XML URL copied to clipboard');
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Company Contacts & Directory"
    subtitle="Global organization address book synchronized directly to physical VoIP handsets and softphones."
  >
    <div class="flex items-center space-x-3">
      <Button variant="secondary" on:click={copyXmlUrl}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
        </svg>
        Copy XML Phonebook URL
      </Button>
      <Button variant="primary" on:click={openCreate}>
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add Contact
      </Button>
    </div>
  </PageHeader>

  <!-- SEARCH & GROUP FILTERS -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
    <div class="flex items-center space-x-3">
      <div class="relative w-72">
        <input
          type="text"
          bind:value={search}
          placeholder="Search name, company, number..."
          class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
        />
        <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
      </div>

      <select
        bind:value={selectedGroup}
        class="text-xs bg-white border border-slate-200 rounded-md px-3 py-1.5 shadow-sm focus:outline-none focus:border-sky-500"
      >
        {#each groups as g}
          <option value={g}>Group: {g}</option>
        {/each}
      </select>
    </div>

    <div class="text-xs text-slate-500">
      Total Directory Contacts: <span class="font-bold text-slate-700">{contacts.length}</span>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading company directory...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if filteredContacts.length === 0}
    <EmptyState
      title="No Directory Contacts Found"
      description="Add clients, vendors, and partner departments so staff can quickly dial them from any desk phone."
    >
      <Button variant="primary" on:click={openCreate}>Create First Contact</Button>
    </EmptyState>
  {:else}
    <Card padding="none">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
              <th class="py-3 px-4">Contact Name</th>
              <th class="py-3 px-4">Organization</th>
              <th class="py-3 px-4">Phone Number</th>
              <th class="py-3 px-4">Mobile</th>
              <th class="py-3 px-4">Speed Dial</th>
              <th class="py-3 px-4">Group</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            {#each filteredContacts as c}
              <tr class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3 px-4 font-semibold text-slate-900 flex items-center space-x-2.5">
                  <div class="w-7 h-7 rounded-full bg-sky-100 text-sky-700 font-bold flex items-center justify-center text-xs">
                    {c.FirstName[0] || 'C'}
                  </div>
                  <div>
                    <span>{c.FirstName} {c.LastName || ''}</span>
                    {#if c.Email}
                      <span class="block text-[11px] text-slate-400 font-normal">{c.Email}</span>
                    {/if}
                  </div>
                </td>
                <td class="py-3 px-4 text-slate-600">{c.Company || '—'}</td>
                <td class="py-3 px-4 font-mono font-medium text-sky-700">
                  {c.PhoneNumber}
                </td>
                <td class="py-3 px-4 font-mono text-slate-600">
                  {c.MobileNumber || '—'}
                </td>
                <td class="py-3 px-4">
                  {#if c.SpeedDial}
                    <span class="px-2 py-0.5 bg-slate-100 rounded text-slate-800 font-mono font-bold text-[11px]">
                      {c.SpeedDial}
                    </span>
                  {:else}
                    <span class="text-slate-400">—</span>
                  {/if}
                </td>
                <td class="py-3 px-4">
                  <Badge variant="neutral">{c.Group || 'General'}</Badge>
                </td>
                <td class="py-3 px-4 text-right space-x-2">
                  <button
                    type="button"
                    class="text-sky-600 hover:text-sky-800 font-medium"
                    on:click={() => openEdit(c)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    class="text-rose-600 hover:text-rose-800 font-medium"
                    on:click={() => deleteContact(c.ID, `${c.FirstName} ${c.LastName || ''}`)}
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
</div>

<!-- CONTACT MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? 'Edit Directory Contact' : 'New Directory Contact'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="ct-fname" class="block font-medium text-slate-700 mb-1">First Name</label>
        <input
          id="ct-fname"
          type="text"
          bind:value={form.FirstName}
          placeholder="e.g. Alice"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
      <div>
        <label for="ct-lname" class="block font-medium text-slate-700 mb-1">Last Name</label>
        <input
          id="ct-lname"
          type="text"
          bind:value={form.LastName}
          placeholder="e.g. Smith"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="ct-comp" class="block font-medium text-slate-700 mb-1">Company / Department</label>
        <input
          id="ct-comp"
          type="text"
          bind:value={form.Company}
          placeholder="e.g. Acme Corp or Legal Dept"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
      <div>
        <label for="ct-grp" class="block font-medium text-slate-700 mb-1">Directory Group</label>
        <input
          id="ct-grp"
          type="text"
          bind:value={form.Group}
          placeholder="e.g. Clients, Suppliers, VIP"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="ct-phone" class="block font-medium text-slate-700 mb-1">Primary Phone Number</label>
        <input
          id="ct-phone"
          type="text"
          bind:value={form.PhoneNumber}
          placeholder="e.g. +15550199 or 1002"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
      <div>
        <label for="ct-mob" class="block font-medium text-slate-700 mb-1">Mobile Number (Optional)</label>
        <input
          id="ct-mob"
          type="text"
          bind:value={form.MobileNumber}
          placeholder="e.g. +15550199123"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="ct-email" class="block font-medium text-slate-700 mb-1">Email Address (Optional)</label>
        <input
          id="ct-email"
          type="email"
          bind:value={form.Email}
          placeholder="alice@example.com"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
      </div>
      <div>
        <label for="ct-spd" class="block font-medium text-slate-700 mb-1">Speed Dial Code (Optional)</label>
        <input
          id="ct-spd"
          type="text"
          bind:value={form.SpeedDial}
          placeholder="e.g. 01, 02"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono font-bold"
        />
      </div>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Add Contact'}
      </Button>
    </div>
  </form>
</Modal>
