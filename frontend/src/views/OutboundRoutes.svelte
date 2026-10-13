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

  interface OutboundRouteItem {
    ID: string;
    Name: string;
    DialPattern: string;
    PrefixStrip: number;
    PrependDigits: string;
    TrunkID: string;
    FallbackTrunkID: string;
    Priority: number;
    OutboundCIDNumber: string;
    OutboundCIDName: string;
    EmergencyRoute: boolean;
    CreatedAt: string;
  }

  interface TrunkItem {
    ID: string;
    ProviderName: string;
  }

  let routes: OutboundRouteItem[] = [];
  let trunks: TrunkItem[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    Name: '',
    DialPattern: '',
    PrefixStrip: 0,
    PrependDigits: '',
    TrunkID: '',
    FallbackTrunkID: '',
    Priority: 1,
    OutboundCIDNumber: '',
    OutboundCIDName: '',
    EmergencyRoute: false
  };

  $: filtered = routes.filter(r => 
    !search.trim() || `${r.Name || ''} ${r.DialPattern} ${r.OutboundCIDNumber || ''}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  function getTrunkName(id: string) {
    const found = trunks.find(t => t.ID === id);
    return found ? found.ProviderName : (id ? id.slice(0, 8) : '—');
  }

  async function load() {
    loading = true;
    loadError = '';
    try {
      const [rData, tData] = await Promise.all([
        api.get<{ outbound_routes: OutboundRouteItem[] }>('/outbound-routes'),
        api.get<{ trunks: TrunkItem[] }>('/trunks')
      ]);
      routes = rData.outbound_routes || [];
      trunks = tData.trunks || [];
      if (trunks.length > 0 && !form.TrunkID) {
        form.TrunkID = trunks[0].ID;
      }
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
      DialPattern: '9|.',
      PrefixStrip: 1,
      PrependDigits: '',
      TrunkID: trunks.length > 0 ? trunks[0].ID : '',
      FallbackTrunkID: '',
      Priority: routes.length + 1,
      OutboundCIDNumber: '',
      OutboundCIDName: '',
      EmergencyRoute: false
    };
    modalOpen = true;
  }

  function openEdit(r: OutboundRouteItem) {
    editingId = r.ID;
    form = {
      Name: r.Name || '',
      DialPattern: r.DialPattern,
      PrefixStrip: r.PrefixStrip || 0,
      PrependDigits: r.PrependDigits || '',
      TrunkID: r.TrunkID,
      FallbackTrunkID: r.FallbackTrunkID || '',
      Priority: r.Priority || 1,
      OutboundCIDNumber: r.OutboundCIDNumber || '',
      OutboundCIDName: r.OutboundCIDName || '',
      EmergencyRoute: r.EmergencyRoute || false
    };
    modalOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.DialPattern.trim()) return toast.error('Dial pattern is required');
    if (!form.TrunkID) return toast.error('Primary Trunk selection is required');

    saving = true;
    try {
      const payload: any = {
        Name: form.Name.trim() || `Outbound ${form.DialPattern.trim()}`,
        DialPattern: form.DialPattern.trim(),
        PrefixStrip: Number(form.PrefixStrip) || 0,
        PrependDigits: form.PrependDigits.trim(),
        TrunkID: form.TrunkID,
        FallbackTrunkID: form.FallbackTrunkID,
        Priority: Number(form.Priority) || 1,
        OutboundCIDNumber: form.OutboundCIDNumber.trim(),
        OutboundCIDName: form.OutboundCIDName.trim(),
        EmergencyRoute: form.EmergencyRoute
      };

      if (editingId) {
        await api.put(`/outbound-routes/${editingId}`, payload);
        toast.success('Outbound route updated');
      } else {
        await api.post('/outbound-routes', payload);
        toast.success('Outbound route created');
      }
      modalOpen = false;
      await load();
    } catch (err) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function remove(r: OutboundRouteItem) {
    if (!confirm(`Delete outbound dial rule "${r.Name || r.DialPattern}"?`)) return;
    try {
      await api.del(`/outbound-routes/${r.ID}`);
      toast.success('Outbound route deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  onMount(load);
</script>

<PageHeader title="Outbound Dial Rules" description="Configure dial patterns (e.g. 9 for outside line, 00 for international) and carrier trunk failovers.">
  <Button on:click={openCreate}>+ New Dial Rule</Button>
</PageHeader>

<Card padded={false} title="Configured Dial Rules" subtitle="{routes.length} outbound patterns">
  <svelte:fragment slot="actions">
    <input class="input !py-1.5 !text-xs w-60" placeholder="Search rule, pattern, CID…" bind:value={search} aria-label="Search routes" />
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search ? 'No dial rules match your search' : 'No outbound routes configured'}
      message={search ? '' : 'Add an outbound dial rule so internal extensions can route calls to external phone carriers.'}
    >
      {#if !search}<Button size="sm" on:click={openCreate}>Create First Dial Rule</Button>{/if}
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Priority</th>
            <th>Rule Name</th>
            <th>Dial Pattern Match</th>
            <th>Digit Manipulation</th>
            <th>Primary Gateway</th>
            <th>Failover Trunk</th>
            <th>Outbound Caller ID</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as r (r.ID)}
            <tr class="hover:bg-slate-50 transition-colors">
              <td class="mono font-bold text-slate-500 text-xs">#{r.Priority}</td>
              <td class="font-bold text-slate-800 text-sm">
                <div class="flex items-center space-x-2">
                  <span>{r.Name || 'Dial Rule'}</span>
                  {#if r.EmergencyRoute}
                    <span class="px-1.5 py-0.5 rounded bg-rose-100 text-rose-800 text-[10px] font-bold">911 / 999</span>
                  {/if}
                </div>
              </td>
              <td>
                <span class="mono px-2 py-0.5 bg-slate-100 text-slate-800 rounded text-xs border border-slate-200 font-bold">
                  {r.DialPattern}
                </span>
              </td>
              <td class="text-xs text-slate-600 mono">
                {#if r.PrefixStrip > 0}
                  <span>Strip {r.PrefixStrip}</span>
                {/if}
                {#if r.PrependDigits}
                  <span> + Prepend {r.PrependDigits}</span>
                {/if}
                {#if !r.PrefixStrip && !r.PrependDigits}
                  <span class="text-slate-400">Exact</span>
                {/if}
              </td>
              <td class="text-xs font-semibold text-[#0284C7]">
                {getTrunkName(r.TrunkID)}
              </td>
              <td class="text-xs text-slate-600">
                {r.FallbackTrunkID ? getTrunkName(r.FallbackTrunkID) : '—'}
              </td>
              <td class="mono text-xs text-slate-800">
                {r.OutboundCIDNumber || 'Extension CID'}
              </td>
              <td class="text-right space-x-3 whitespace-nowrap">
                <button type="button" class="link-action" on:click={() => openEdit(r)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(r)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<Modal bind:isOpen={modalOpen} title={editingId ? 'Edit Dial Rule' : 'Create Outbound Dial Rule'} maxWidth="max-w-2xl">
  <form on:submit={save} class="space-y-4">
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div class="sm:col-span-2">
        <label for="out-name" class="label">Rule Name *</label>
        <input id="out-name" class="input font-bold" bind:value={form.Name} placeholder="e.g. Standard Domestic Calls, International Dialing, 911 Emergency" required />
      </div>

      <div>
        <label for="out-pattern" class="label">Dial Pattern *</label>
        <input id="out-pattern" class="input mono font-bold" bind:value={form.DialPattern} placeholder="e.g. 9|XXXXXXXXXX or 00X." required />
        <p class="text-[11px] text-slate-500 mt-1">FreeSWITCH regex or standard prefix pattern (e.g. 9 for outside line, 00 for international).</p>
      </div>

      <div>
        <label for="out-priority" class="label">Rule Evaluation Priority</label>
        <input id="out-priority" type="number" min="1" max="100" class="input mono font-bold" bind:value={form.Priority} />
        <p class="text-[11px] text-slate-500 mt-1">Lower numbers are evaluated first.</p>
      </div>

      <!-- Digit Manipulation -->
      <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg sm:col-span-2 grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label for="out-strip" class="label">Digits to Strip from Front</label>
          <input id="out-strip" type="number" min="0" max="10" class="input mono" bind:value={form.PrefixStrip} />
          <p class="text-[11px] text-slate-500 mt-1">E.g. Set to 1 if callers dial 9 to reach an outside line.</p>
        </div>

        <div>
          <label for="out-prepend" class="label">Digits to Prepend (Prefix)</label>
          <input id="out-prepend" class="input mono" bind:value={form.PrependDigits} placeholder="e.g. +1 or 011" />
          <p class="text-[11px] text-slate-500 mt-1">Added in front of the number sent to the carrier.</p>
        </div>
      </div>

      <!-- Gateways -->
      <div>
        <label for="out-trunk" class="label">Primary Carrier Trunk *</label>
        <select id="out-trunk" class="input font-semibold" bind:value={form.TrunkID} required>
          <option value="">Select Carrier Trunk...</option>
          {#each trunks as t}
            <option value={t.ID}>{t.ProviderName}</option>
          {/each}
        </select>
      </div>

      <div>
        <label for="out-fb-trunk" class="label">Failover Backup Trunk (Optional)</label>
        <select id="out-fb-trunk" class="input" bind:value={form.FallbackTrunkID}>
          <option value="">None (Fail if primary unreachable)</option>
          {#each trunks.filter(t => t.ID !== form.TrunkID) as t}
            <option value={t.ID}>{t.ProviderName}</option>
          {/each}
        </select>
      </div>

      <!-- Override Caller ID -->
      <div>
        <label for="out-cid-num" class="label">Override Outbound Caller ID Number</label>
        <input id="out-cid-num" class="input mono" bind:value={form.OutboundCIDNumber} placeholder="Leave empty to use user's extension DID" />
      </div>

      <div>
        <label for="out-cid-name" class="label">Override Outbound Caller ID Name</label>
        <input id="out-cid-name" class="input" bind:value={form.OutboundCIDName} placeholder="Leave empty to use user's display name" />
      </div>

      <div class="sm:col-span-2 flex items-center space-x-2 pt-2">
        <input id="out-emergency" type="checkbox" bind:checked={form.EmergencyRoute} class="rounded border-slate-300 text-rose-600 focus:ring-rose-500" />
        <label for="out-emergency" class="text-xs font-bold text-rose-700">Mark as Emergency Route (Bypasses all outbound barring and restrictions)</label>
      </div>
    </div>

    <div class="flex justify-end gap-2 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : (editingId ? 'Update Rule' : 'Save Rule')}</Button>
    </div>
  </form>
</Modal>
