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

  interface Device {
    ID: string;
    MACAddress: string;
    BrandModel: string;
    ExtensionNumber: string;
    IPAddress: string;
    FirmwareVersion: string;
    Status: string;
  }

  interface ExtensionOption {
    ExtensionNumber: string;
    DisplayName?: string;
  }

  let devices: Device[] = [];
  let availableExts: ExtensionOption[] = [];
  let loading = true;
  let loadError = '';
  let saving = false;
  let search = '';

  let modalOpen = false;
  let editingId: string | null = null;
  let form = {
    MACAddress: '',
    BrandModel: 'Generic HD SIP Phone',
    ExtensionNumber: '',
    IPAddress: '',
    FirmwareVersion: 'v1.0.0',
    Status: 'online'
  };

  $: filteredDevices = devices.filter(d =>
    !search.trim() || `${d.MACAddress} ${d.BrandModel} ${d.ExtensionNumber} ${d.IPAddress}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  async function loadData() {
    loading = true;
    loadError = '';
    try {
      const [devRes, extRes] = await Promise.all([
        api.get<{ devices: Device[] }>('/devices'),
        api.get<{ extensions: ExtensionOption[] }>('/extensions')
      ]);
      devices = devRes.devices || [];
      availableExts = extRes.extensions || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    editingId = null;
    form = {
      MACAddress: '',
      BrandModel: '',
      ExtensionNumber: availableExts[0]?.ExtensionNumber || '',
      IPAddress: '',
      FirmwareVersion: '',
      Status: 'offline'
    };
    modalOpen = true;
  }

  function openEdit(d: Device) {
    editingId = d.ID;
    form = {
      MACAddress: d.MACAddress,
      BrandModel: d.BrandModel,
      ExtensionNumber: d.ExtensionNumber,
      IPAddress: d.IPAddress || '',
      FirmwareVersion: d.FirmwareVersion || '',
      Status: d.Status || 'online'
    };
    modalOpen = true;
  }

  const BRAND_PRESETS = [
    'Yealink T46U (16-Line Gigabit IP Phone)',
    'Yealink T48U (7-inch Touchscreen Executive)',
    'Yealink T54W (Built-in Bluetooth & Wi-Fi)',
    'Yealink T33G (Entry-Level 4-Line Phone)',
    'Yealink W73P (Cordless DECT Handset & Base)',
    'Grandstream GRP2614 (Carrier-Grade Dual-Screen)',
    'Grandstream GRP2615 (10-Line Professional)',
    'Grandstream GXP2170 (12-Line Enterprise)',
    'Grandstream HT802 (2-Port Analog Telephone Adapter)',
    'Fanvil X6U (High-End 20-Line Phone)',
    'Fanvil V64 (Prime Business Phone with Wi-Fi)',
    'Fanvil i10V (SIP Video Intercom / Door Station)',
    'Poly VVX 450 (12-Line IP Desk Phone)',
    'Poly Edge E450 (8-Key Color IP Phone)',
    'Cisco 8841 (Gigabit Multiplatform Firmware)',
    'Cisco 8851 (Multiplatform with Bluetooth/USB)',
    'Snom D785 (Dual Display Desk Phone)',
    'Generic HD SIP Phone'
  ];

  async function rebootDevice(d: Device) {
    if (!confirm(`Send remote SIP check-sync reboot signal to handset ${d.BrandModel} (${d.MACAddress})?`)) return;
    try {
      await api.post(`/devices/${d.ID}/reboot`, {});
      toast.success(`Reboot command dispatched to MAC ${d.MACAddress}`);
    } catch (err: any) {
      toastError(err);
    }
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!form.MACAddress.trim()) return toast.error('MAC address is required');

    saving = true;
    try {
      if (editingId) {
        await api.put(`/devices/${editingId}`, form);
        toast.success('Device updated');
      } else {
        await api.post('/devices', form);
        toast.success('Device registered');
      }
      modalOpen = false;
      await loadData();
    } catch (err: any) {
      toastError(err);
    } finally {
      saving = false;
    }
  }

  async function deleteDevice(id: string, mac: string) {
    if (!confirm(`Delete device "${mac}" from inventory?`)) return;
    try {
      await api.del(`/devices/${id}`);
      toast.success('Device deleted');
      devices = devices.filter(d => d.ID !== id);
    } catch (err: any) {
      toastError(err);
    }
  }

  onMount(loadData);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Device Inventory (Handsets & ATAs)"
    subtitle="Registered physical SIP desk phones, MAC addresses, LAN IPs, and remote reboot control."
  >
    <Button variant="primary" on:click={openCreate}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Hardware Device
    </Button>
  </PageHeader>

  <!-- SEARCH BAR -->
  <div class="flex items-center justify-between">
    <div class="relative w-72">
      <input
        type="text"
        bind:value={search}
        placeholder="Search MAC, model, extension..."
        class="w-full pl-9 pr-3 py-1.5 text-xs bg-white border border-slate-200 rounded-md shadow-sm focus:outline-none focus:border-sky-500 focus:ring-1 focus:ring-sky-500"
      />
      <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
    <div class="text-xs text-slate-500">
      Provisioned Devices: <span class="font-bold text-slate-700">{devices.length}</span>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Loading hardware phone inventory...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if filteredDevices.length === 0}
    <EmptyState
      title="No Hardware Devices Found"
      description="Register desk phones by MAC address to bind them to user extensions and enable zero-touch auto provisioning."
    >
      <Button variant="primary" on:click={openCreate}>Register First Handset</Button>
    </EmptyState>
  {:else}
    <Card padding="none">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase tracking-wider">
              <th class="py-3 px-4">MAC Address</th>
              <th class="py-3 px-4">Brand & Model</th>
              <th class="py-3 px-4">Assigned Extension</th>
              <th class="py-3 px-4">LAN IP Address</th>
              <th class="py-3 px-4">Firmware</th>
              <th class="py-3 px-4">Provisioning Status</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            {#each filteredDevices as d}
              <tr class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3 px-4 font-mono font-bold text-slate-900 uppercase">
                  {d.MACAddress}
                </td>
                <td class="py-3 px-4 font-medium text-slate-800">{d.BrandModel}</td>
                <td class="py-3 px-4">
                  <span class="px-2 py-0.5 bg-sky-50 text-sky-700 border border-sky-200 rounded font-mono font-bold text-[11px]">
                    Ext {d.ExtensionNumber}
                  </span>
                </td>
                <td class="py-3 px-4 font-mono text-slate-600">{d.IPAddress || 'DHCP (Dynamic)'}</td>
                <td class="py-3 px-4 font-mono text-slate-500 text-[11px]">{d.FirmwareVersion || '—'}</td>
                <td class="py-3 px-4">
                  {#if d.Status === 'online'}
                    <Badge variant="success">Online / Provisioned</Badge>
                  {:else if d.Status === 'provisioning'}
                    <Badge variant="info">Syncing</Badge>
                  {:else}
                    <Badge variant="neutral">Offline</Badge>
                  {/if}
                </td>
                <td class="py-3 px-4 text-right space-x-2">
                  <button
                    type="button"
                    class="text-amber-600 hover:text-amber-800 font-medium"
                    on:click={() => rebootDevice(d)}
                    title="Send SIP check-sync reboot"
                  >
                    Reboot
                  </button>
                  <button
                    type="button"
                    class="text-sky-600 hover:text-sky-800 font-medium"
                    on:click={() => openEdit(d)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    class="text-rose-600 hover:text-rose-800 font-medium"
                    on:click={() => deleteDevice(d.ID, d.MACAddress)}
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

<!-- DEVICE MODAL -->
<Modal
  open={modalOpen}
  title={editingId ? 'Edit Hardware Device' : 'Register Hardware Device'}
  on:close={() => modalOpen = false}
>
  <form on:submit={save} class="space-y-4 text-xs">
    <div>
      <label for="dev-mac" class="block font-medium text-slate-700 mb-1">Physical MAC Address</label>
      <input
        id="dev-mac"
        type="text"
        bind:value={form.MACAddress}
        placeholder="e.g. 00:15:65:8A:2B:11 or 0015658A2B11"
        required
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono font-bold uppercase"
      />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="dev-brand" class="block font-medium text-slate-700 mb-1">Brand & Model</label>
        <input
          id="dev-brand"
          type="text"
          list="brand-presets-list"
          bind:value={form.BrandModel}
          placeholder="Select preset or type model"
          required
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500"
        />
        <datalist id="brand-presets-list">
          {#each BRAND_PRESETS as preset}
            <option value={preset}></option>
          {/each}
        </datalist>
      </div>
      <div>
        <label for="dev-ext" class="block font-medium text-slate-700 mb-1">Assigned Extension</label>
        <select
          id="dev-ext"
          bind:value={form.ExtensionNumber}
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
        >
          <option value="">-- Unassigned (Hot-Desking Pool) --</option>
          {#each availableExts as ext}
            <option value={ext.ExtensionNumber}>Ext {ext.ExtensionNumber} ({ext.DisplayName || 'User'})</option>
          {/each}
        </select>
      </div>
    </div>

    <!-- Provisioning Guide Pill -->
    <div class="p-3 bg-sky-50/70 border border-sky-200 rounded-lg text-xs space-y-1">
      <div class="font-semibold text-sky-900 flex items-center">
        <svg class="w-3.5 h-3.5 mr-1 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Auto-Provisioning URL (HTTP/HTTPS)
      </div>
      <p class="text-sky-700">Enter this configuration URL in your handset's Auto-Provision / Upgrade menu:</p>
      <code class="block bg-white px-2 py-1 border border-sky-300 rounded text-sky-900 font-mono text-[11px] select-all">
        {window.location.origin}/provision/
      </code>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="dev-ip" class="block font-medium text-slate-700 mb-1">Local IP Address</label>
        <input
          id="dev-ip"
          type="text"
          bind:value={form.IPAddress}
          placeholder="e.g. 192.168.1.105"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
      <div>
        <label for="dev-fw" class="block font-medium text-slate-700 mb-1">Firmware Version</label>
        <input
          id="dev-fw"
          type="text"
          bind:value={form.FirmwareVersion}
          placeholder="e.g. v66.86.0.5"
          class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 font-mono"
        />
      </div>
    </div>

    <div>
      <label for="dev-status" class="block font-medium text-slate-700 mb-1">Status</label>
      <select
        id="dev-status"
        bind:value={form.Status}
        class="w-full px-3 py-2 border border-slate-200 rounded-md focus:outline-none focus:border-sky-500 bg-white"
      >
        <option value="online">Online / Registered</option>
        <option value="provisioning">Provisioning in Progress</option>
        <option value="offline">Offline / Disconnected</option>
      </select>
    </div>

    <div class="flex justify-end space-x-3 pt-3 border-t border-slate-200">
      <Button variant="secondary" type="button" on:click={() => modalOpen = false}>Cancel</Button>
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : editingId ? 'Save Changes' : 'Register Device'}
      </Button>
    </div>
  </form>
</Modal>
