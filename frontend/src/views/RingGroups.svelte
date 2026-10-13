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

  interface RingGroup {
    ID: string;
    Name: string;
    Extension: string;
    Strategy: string;
    Members: string[];
  }

  interface ExtensionItem {
    ID: string;
    ExtensionNumber: string;
  }

  let groups: RingGroup[] = [];
  let availableExts: ExtensionItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    Extension: '',
    Strategy: 'simultaneous',
    MembersText: ''
  };

  $: filtered = groups.filter(g => 
    !search.trim() || `${g.Name} ${g.Extension}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const [ivrsData, extsData] = await Promise.all([
        api.get<{ ivrs: any[] }>('/ivrs'),
        api.get<{ extensions: ExtensionItem[] }>('/extensions')
      ]);
      availableExts = extsData.extensions || [];
      // Ring groups use IVRMenu model with prefix "RG:"
      const allIvrs = ivrsData.ivrs || [];
      groups = allIvrs
        .filter(i => i.Name.startsWith('Ring Group:') || i.Greeting === 'ring_group')
        .map(i => ({
          ID: i.ID,
          Name: i.Name.replace('Ring Group:', '').trim() || i.Name,
          Extension: i.Extension,
          Strategy: 'Simultaneous Ring',
          Members: (i.Choices || []).map((c: any) => c.Destination)
        }));
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      Name: '',
      Extension: '600',
      Strategy: 'simultaneous',
      MembersText: availableExts.map(e => e.ExtensionNumber).slice(0, 3).join(', ')
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Ring group name is required');
    if (!form.Extension.trim()) return toast.error('Extension is required');

    const members = form.MembersText.split(',').map(m => m.trim()).filter(Boolean);
    if (members.length === 0) return toast.error('Please enter at least one member extension');

    saving = true;
    try {
      const payload = {
        Name: `Ring Group: ${form.Name.trim()}`,
        Extension: form.Extension.trim(),
        Greeting: 'ring_group',
        Timeout: 30000,
        Choices: members.map((m, idx) => ({
          Digits: String(idx + 1),
          Action: 'bridge',
          Destination: m
        }))
      };

      if (editingId) {
        await api.put(`/ivrs/${editingId}`, payload);
        toast.success('Ring Group updated');
      } else {
        await api.post('/ivrs', payload);
        toast.success('Ring Group created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(g: RingGroup) {
    if (!confirm(`Delete ring group "${g.Name}"?`)) return;
    try {
      await api.del(`/ivrs/${g.ID}`);
      toast.success('Ring group deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Ring Groups" description="Ring multiple extensions simultaneously or sequentially when a virtual extension is dialed.">
  <Button on:click={openCreate}>+ New Ring Group</Button>
</PageHeader>

<Card padded={false} title="Active Ring Groups" subtitle="{groups.length} groups configured">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search ring group…" bind:value={search} aria-label="Search ring groups" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No ring groups match your search' : 'No ring groups configured'}
      message={search ? '' : 'Create a ring group (e.g. Sales or Support) to ring multiple extensions when called.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Create first ring group</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Group Name</th>
            <th>Extension</th>
            <th>Strategy</th>
            <th>Destination Members</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as g (g.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{g.Name}</td>
              <td class="mono font-semibold text-sky-700 text-xs">{g.Extension}</td>
              <td><Badge tone="info">{g.Strategy}</Badge></td>
              <td>
                <div class="flex flex-wrap gap-1">
                  {#each g.Members as m}
                    <span class="mono px-2 py-0.5 bg-slate-100 rounded text-xs text-slate-700 border border-slate-200">
                      Ext {m}
                    </span>
                  {/each}
                </div>
              </td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-danger" on:click={() => remove(g)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Ring Group' : 'Create Ring Group'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="rg-name" class="label">Group Name</label>
      <input id="rg-name" class="input" bind:value={form.Name} placeholder="e.g. Sales Team or Tech Support" required />
    </div>

    <div>
      <label for="rg-ext" class="label">Group Virtual Extension</label>
      <input id="rg-ext" class="input mono" bind:value={form.Extension} placeholder="e.g. 600" required />
    </div>

    <div>
      <label for="rg-members" class="label">Member Extensions (comma-separated)</label>
      <input id="rg-members" class="input mono" bind:value={form.MembersText} placeholder="e.g. 1001, 1002, 1003" required />
      <p class="help">Enter extension numbers to ring when this group is called.</p>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
