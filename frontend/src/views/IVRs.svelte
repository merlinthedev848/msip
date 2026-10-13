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
    MaxFailures: number;
    TimeoutDestType: string;
    TimeoutDestTarget: string;
    InvalidDestType: string;
    InvalidDestTarget: string;
    AllowDirectDial: boolean;
    Choices: Choice[];
    CreatedAt: string;
  }

  let ivrs: IVRItem[] = [];
  let availableExts: any[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    Extension: '',
    Greeting: 'ivr/welcome.wav',
    Timeout: 5,
    MaxFailures: 3,
    TimeoutDestType: 'EXTENSION',
    TimeoutDestTarget: '',
    InvalidDestType: 'EXTENSION',
    InvalidDestTarget: '',
    AllowDirectDial: true,
    Choices: [] as Choice[]
  };

  $: filtered = ivrs.filter(i => 
    !search.trim() || `${i.Name} ${i.Extension}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function load() {
    loading = true;
    loadError = '';
    try {
      const [ivrData, extData] = await Promise.all([
        api.get<{ ivrs: IVRItem[] }>('/ivrs'),
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] }))
      ]);
      ivrs = ivrData.ivrs || [];
      availableExts = extData.extensions || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    const defaultExt = availableExts.length > 0 ? availableExts[0].ExtensionNumber : '1001';
    form = {
      Name: '',
      Extension: '7000',
      Greeting: 'ivr/welcome.wav',
      Timeout: 5,
      MaxFailures: 3,
      TimeoutDestType: 'EXTENSION',
      TimeoutDestTarget: defaultExt,
      InvalidDestType: 'EXTENSION',
      InvalidDestTarget: defaultExt,
      AllowDirectDial: true,
      Choices: [
        { Digits: '1', Action: 'transfer', Destination: defaultExt },
        { Digits: '2', Action: 'transfer', Destination: defaultExt }
      ]
    };
    modalOpen = true;
  }

  function openEdit(i: IVRItem) {
    editingId = i.ID;
    form = {
      Name: i.Name,
      Extension: i.Extension,
      Greeting: i.Greeting,
      Timeout: Math.round((i.Timeout || 5000) / 1000) || 5,
      MaxFailures: i.MaxFailures || 3,
      TimeoutDestType: i.TimeoutDestType || 'EXTENSION',
      TimeoutDestTarget: i.TimeoutDestTarget || '',
      InvalidDestType: i.InvalidDestType || 'EXTENSION',
      InvalidDestTarget: i.InvalidDestTarget || '',
      AllowDirectDial: i.AllowDirectDial ?? true,
      Choices: (i.Choices || []).map(c => ({ Digits: c.Digits, Action: c.Action || 'transfer', Destination: c.Destination }))
    };
    modalOpen = true;
  }

  function addChoice() {
    const nextDigit = String(form.Choices.length + 1);
    form.Choices = [...form.Choices, { Digits: nextDigit, Action: 'transfer', Destination: '' }];
  }

  function removeChoice(idx: number) {
    form.Choices = form.Choices.filter((_, i) => i !== idx);
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.Name.trim()) return toast.error('Menu name is required');
    if (!form.Extension.trim()) return toast.error('Extension number is required');

    saving = true;
    try {
      const payload: any = {
        Name: form.Name.trim(),
        Extension: form.Extension.trim(),
        Greeting: form.Greeting.trim(),
        Timeout: (Number(form.Timeout) || 5) * 1000,
        MaxFailures: Number(form.MaxFailures) || 3,
        TimeoutDestType: form.TimeoutDestType,
        TimeoutDestTarget: form.TimeoutDestTarget.trim(),
        InvalidDestType: form.InvalidDestType,
        InvalidDestTarget: form.InvalidDestTarget.trim(),
        AllowDirectDial: form.AllowDirectDial,
        Choices: form.Choices
      };

      if (editingId) {
        await api.put(`/ivrs/${editingId}`, payload);
        toast.success('Auto Attendant (IVR) updated');
      } else {
        await api.post('/ivrs', payload);
        toast.success('Auto Attendant (IVR) created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(i: IVRItem) {
    if (!confirm(`Delete Auto Attendant menu "${i.Name}"?`)) return;
    try {
      await api.del(`/ivrs/${i.ID}`);
      toast.success('IVR menu deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Digital Receptionists (Auto Attendant IVR)" description="Create 'Press 1 for Sales, 2 for Support' interactive voice response menus.">
  <Button on:click={openCreate}>+ New Auto Attendant</Button>
</PageHeader>

<Card padded={false} title="Configured Voice Menus" subtitle="{ivrs.length} IVR menus">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-56" placeholder="Search menu…" bind:value={search} aria-label="Search menus" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No menus match your search' : 'No IVR menus configured'}
      message={search ? '' : 'Create an Auto Attendant to greet callers with an audio prompt and route by keypad digits.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Create First IVR Menu</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Menu Name</th>
            <th>Extension</th>
            <th>Greeting Audio Prompt</th>
            <th>Keypad Options</th>
            <th>Timeout Fallback</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as ivr (ivr.ID)}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="font-bold text-slate-800 text-sm">
                <div>{ivr.Name}</div>
                {#if ivr.AllowDirectDial}
                  <span class="text-[10px] text-emerald-700 font-medium">Direct Extension Dialing Allowed</span>
                {/if}
              </td>
              <td class="mono font-bold text-sky-700 text-xs">{ivr.Extension}</td>
              <td class="mono text-xs text-slate-600">{ivr.Greeting}</td>
              <td>
                <div class="flex flex-wrap gap-1">
                  {#each (ivr.Choices || []) as c}
                    <span class="mono px-2 py-0.5 bg-slate-100 rounded text-xs text-slate-800 border border-slate-200 font-semibold">
                      Press {c.Digits} &rarr; {c.Destination}
                    </span>
                  {:else}
                    <span class="text-xs text-slate-400">None</span>
                  {/each}
                </div>
              </td>
              <td class="mono text-xs text-slate-700">
                {ivr.TimeoutDestTarget ? `${ivr.TimeoutDestType}: ${ivr.TimeoutDestTarget}` : 'Default (Operator)'}
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

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Auto Attendant (IVR)' : 'Create Auto Attendant (IVR)'} maxWidth="max-w-2xl">
  <form on:submit={save} class="space-y-4">
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div class="sm:col-span-2">
        <label for="ivr-name" class="label">Menu Name *</label>
        <input id="ivr-name" class="input font-bold" bind:value={form.Name} placeholder="e.g. Main Company Daytime Reception" required />
      </div>

      <div>
        <label for="ivr-ext" class="label">Menu Extension *</label>
        <input id="ivr-ext" class="input mono font-bold" bind:value={form.Extension} placeholder="e.g. 7000" required />
        <p class="text-[11px] text-slate-500 mt-1">Internal number mapped to this automated voice attendant.</p>
      </div>

      <div>
        <label for="ivr-timeout" class="label">Keypress Wait Timeout (Seconds)</label>
        <input id="ivr-timeout" type="number" min="1" max="60" class="input mono font-bold" bind:value={form.Timeout} required />
        <p class="text-[11px] text-slate-500 mt-1">Seconds to wait for caller input after audio prompt finishes.</p>
      </div>

      <div class="sm:col-span-2">
        <label for="ivr-greeting" class="label">Greeting Audio Prompt File *</label>
        <input id="ivr-greeting" class="input mono" bind:value={form.Greeting} placeholder="e.g. ivr/welcome.wav or prompts/office_menu.wav" required />
        <p class="text-[11px] text-slate-500 mt-1">Audio recording played when callers enter this menu (stored in audio prompts directory).</p>
      </div>

      <!-- Fallbacks -->
      <div>
        <label for="ivr-to-dest" class="label">Destination on Timeout (No Key Pressed)</label>
        <div class="flex gap-2">
          <select class="input !w-32 text-xs" bind:value={form.TimeoutDestType}>
            <option value="EXTENSION">Extension</option>
            <option value="VOICEMAIL">Voicemail</option>
            <option value="REPEAT">Replay Menu</option>
            <option value="HANGUP">Hang Up</option>
          </select>
          <input class="input mono text-xs flex-1" bind:value={form.TimeoutDestTarget} placeholder="Target (e.g. 1001)" />
        </div>
      </div>

      <div>
        <label for="ivr-inv-dest" class="label">Destination on Invalid Keypad Digit</label>
        <div class="flex gap-2">
          <select class="input !w-32 text-xs" bind:value={form.InvalidDestType}>
            <option value="EXTENSION">Extension</option>
            <option value="VOICEMAIL">Voicemail</option>
            <option value="REPEAT">Replay Menu</option>
            <option value="HANGUP">Hang Up</option>
          </select>
          <input class="input mono text-xs flex-1" bind:value={form.InvalidDestTarget} placeholder="Target (e.g. 1001)" />
        </div>
      </div>

      <div class="sm:col-span-2 flex items-center space-x-2 pt-1">
        <input id="ivr-direct" type="checkbox" bind:checked={form.AllowDirectDial} class="rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
        <label for="ivr-direct" class="text-xs font-semibold text-slate-800">Allow Callers to Directly Dial Known Extensions (e.g. 1001, 1002)</label>
      </div>
    </div>

    <!-- Keypad Digit Mapping -->
    <div class="border-t border-slate-200 pt-3">
      <div class="flex items-center justify-between mb-2">
        <div>
          <span class="text-xs font-bold text-slate-800 uppercase tracking-wider">Keypad Digit Choices (DTMF)</span>
          <p class="text-[11px] text-slate-500">Define what happens when callers press numbers 0-9, * or #.</p>
        </div>
        <button type="button" on:click={addChoice} class="text-xs text-sky-700 font-bold hover:underline">+ Add Key Choice</button>
      </div>

      <div class="space-y-2">
        {#each form.Choices as choice, idx}
          <div class="flex items-center gap-2 p-2 bg-slate-50 border border-slate-200 rounded-lg">
            <div class="flex items-center space-x-1">
              <span class="text-xs text-slate-500 font-medium">Press</span>
              <input class="input mono w-14 text-center font-bold" bind:value={choice.Digits} placeholder="1" aria-label="Digits" required />
            </div>
            <span class="text-xs text-slate-400 font-bold">&rarr;</span>
            <input class="input mono flex-1" bind:value={choice.Destination} placeholder="Target Extension / Queue (e.g. 1001, 600)" aria-label="Destination" required />
            <button type="button" on:click={() => removeChoice(idx)} class="text-rose-600 hover:text-rose-800 p-1 text-sm font-bold" title="Remove choice">✕</button>
          </div>
        {/each}
      </div>
    </div>

    <div class="flex justify-end gap-2 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : (editingId ? 'Update Auto Attendant' : 'Save Auto Attendant')}</Button>
    </div>
  </form>
</Modal>
