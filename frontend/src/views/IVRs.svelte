<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface Choice {
    Digits: string;
    Action: string;
    Destination: string;
  }

  interface IVRItem {
    ID: string;
    Name: string;
    Extension: string;
    Greeting: string;
    Timeout: number;
    Choices: Choice[];
    CreatedAt: string;
  }

  let ivrs: IVRItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    Extension: '',
    Greeting: '',
    Timeout: 10000,
    Choices: [] as Choice[]
  };

  $: filtered = ivrs.filter(i => 
    !search.trim() || `${i.Name} ${i.Extension}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ ivrs: IVRItem[] }>('/ivrs');
      ivrs = data.ivrs || [];
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
      Extension: '',
      Greeting: 'welcome.wav',
      Timeout: 10000,
      Choices: [{ Digits: '1', Action: 'transfer', Destination: '1001' }]
    };
    modalOpen = true;
  }

  function openEdit(ivr: IVRItem) {
    editingId = ivr.ID;
    form = {
      Name: ivr.Name,
      Extension: ivr.Extension,
      Greeting: ivr.Greeting,
      Timeout: ivr.Timeout || 10000,
      Choices: (ivr.Choices && ivr.Choices.length > 0) ? JSON.parse(JSON.stringify(ivr.Choices)) : [{ Digits: '1', Action: 'transfer', Destination: '1001' }]
    };
    modalOpen = true;
  }

  function addChoice() {
    form.Choices = [...form.Choices, { Digits: '', Action: 'transfer', Destination: '' }];
  }

  function removeChoice(idx: number) {
    form.Choices = form.Choices.filter((_, i) => i !== idx);
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Name is required');
    if (!form.Extension.trim()) return toast.error('Extension is required');

    saving = true;
    try {
      const payload = {
        Name: form.Name.trim(),
        Extension: form.Extension.trim(),
        Greeting: form.Greeting.trim(),
        Timeout: Number(form.Timeout) || 10000,
        Choices: form.Choices.filter(c => c.Digits.trim() && c.Destination.trim())
      };

      if (editingId) {
        await api.put(`/ivrs/${editingId}`, payload);
        toast.success('IVR Menu updated');
      } else {
        await api.post('/ivrs', payload);
        toast.success('IVR Menu created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(ivr: IVRItem) {
    if (!confirm(`Delete IVR menu "${ivr.Name}"?`)) return;
    try {
      await api.del(`/ivrs/${ivr.ID}`);
      toast.success('IVR menu deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Digital Receptionists (IVR)" description="Interactive auto attendants and keypress menu trees.">
  <Button on:click={openCreate}>+ New IVR Menu</Button>
</PageHeader>

<Card padded={false} title="Configured Auto Attendants" subtitle="{ivrs.length} menus">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search menu…" bind:value={search} aria-label="Search menus" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No IVR menus match your search' : 'No IVR menus configured'}
      message={search ? '' : 'Create an IVR menu to greet incoming callers and route by DTMF keypresses.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Add first IVR</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Name</th>
            <th>Extension</th>
            <th>Audio Prompt</th>
            <th>Keypress Options</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as ivr (ivr.ID)}
            <tr>
              <td class="font-bold text-slate-800 text-sm">{ivr.Name}</td>
              <td class="mono font-semibold text-sky-700 text-xs">{ivr.Extension}</td>
              <td class="mono text-xs text-slate-600">{ivr.Greeting}</td>
              <td>
                <div class="flex flex-wrap gap-1">
                  {#each (ivr.Choices || []) as c}
                    <span class="mono px-1.5 py-0.5 bg-slate-100 rounded text-[11px] text-slate-700 border border-slate-200">
                      Press {c.Digits} &rarr; {c.Destination}
                    </span>
                  {:else}
                    <span class="text-xs text-slate-400">None</span>
                  {/each}
                </div>
              </td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(ivr)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(ivr)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit IVR Menu' : 'New IVR Menu'}>
  <form on:submit={save} class="space-y-4">
    <div>
      <label for="ivr-name" class="label">Menu Name</label>
      <input id="ivr-name" class="input" bind:value={form.Name} placeholder="e.g. Main Company Reception" required />
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div>
        <label for="ivr-ext" class="label">Menu Extension</label>
        <input id="ivr-ext" class="input mono" bind:value={form.Extension} placeholder="e.g. 7000" required />
      </div>
      <div>
        <label for="ivr-timeout" class="label">Timeout (ms)</label>
        <input id="ivr-timeout" type="number" class="input mono" bind:value={form.Timeout} required />
      </div>
    </div>

    <div>
      <label for="ivr-greeting" class="label">Audio Prompt / WAV File</label>
      <input id="ivr-greeting" class="input mono" bind:value={form.Greeting} placeholder="e.g. ivr/main_welcome.wav" required />
    </div>

    <div>
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs font-semibold text-slate-700">DTMF Keypress Choices</span>
        <button type="button" on:click={addChoice} class="text-xs text-sky-700 font-semibold hover:underline">+ Add Key</button>
      </div>

      <div class="space-y-2">
        {#each form.Choices as choice, idx}
          <div class="flex items-center gap-2">
            <input class="input mono w-20 text-center" bind:value={choice.Digits} placeholder="Key" aria-label="Digits" required />
            <span class="text-xs text-slate-400">&rarr;</span>
            <input class="input mono flex-1" bind:value={choice.Destination} placeholder="Target Ext/Queue" aria-label="Destination" required />
            <button type="button" on:click={() => removeChoice(idx)} class="text-rose-600 hover:text-rose-800 p-1" title="Remove choice">✕</button>
          </div>
        {/each}
      </div>
    </div>

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="secondary" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
    </div>
  </form>
</Modal>
