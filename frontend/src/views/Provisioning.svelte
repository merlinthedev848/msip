<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface ProvisionedDevice {
    mac: string;
    ext: string;
    code: string;
    status: string;
  }

  let devices: ProvisionedDevice[] = [];
  let loading = true;
  let loadError = '';
  let copyFeedback = false;

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ extensions: any[] }>('/extensions');
      devices = (data.extensions || []).map(ext => {
        const padded = (ext.ExtensionNumber || '0').padStart(6, '0');
        const lastOctets = padded.slice(-6).match(/.{2}/g)?.join(':') || '00:00:00';
        return {
          mac: `00:15:65:${lastOctets}`.toUpperCase(),
          ext: ext.ExtensionNumber,
          code: ext.ProvisioningCode || 'N/A',
          status: ext.IsActive ? 'Active' : 'Disabled'
        };
      });
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function copyUrl() {
    const url = `http://${window.location.hostname}:8081/api/v1/provisioning/`;
    navigator.clipboard.writeText(url);
    copyFeedback = true;
    setTimeout(() => copyFeedback = false, 2000);
  }

  onMount(load);
</script>

<PageHeader title="Auto Provisioning" description="Zero-touch SIP endpoint provisioning URL and mapped extension codes.">
  <Button variant="secondary" on:click={load}>Refresh</Button>
</PageHeader>

<div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
  <div class="lg:col-span-2">
    <Card padded={false} title="Configured Endpoints" subtitle="{devices.length} endpoints available for auto-provisioning">
      {#if loading || loadError || devices.length === 0}
        <EmptyState
          {loading}
          error={loadError}
          title="No provisionable extensions found"
          message="Create extensions in the Extensions view to generate provisioning codes."
        >
          <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
        </EmptyState>
      {:else}
        <div class="overflow-x-auto">
          <table class="tbl">
            <thead>
              <tr>
                <th>Extension</th>
                <th>Provisioning Code</th>
                <th>Device Identifier</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {#each devices as dev}
                <tr>
                  <td class="font-bold text-slate-800 mono">{dev.ext}</td>
                  <td>
                    <span class="mono px-2 py-0.5 bg-sky-50 text-sky-800 rounded font-semibold text-xs border border-sky-200">
                      {dev.code}
                    </span>
                  </td>
                  <td class="mono text-xs text-slate-600">{dev.mac}</td>
                  <td>
                    <Badge tone={dev.status === 'Active' ? 'success' : 'neutral'}>{dev.status}</Badge>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </Card>
  </div>

  <div class="space-y-5">
    <Card title="Server Provisioning URL">
      <div class="space-y-3">
        <p class="text-xs text-slate-600 leading-relaxed">
          SIP devices query the HTTP provisioning endpoint with their 8-digit provisioning code to retrieve registration credentials automatically:
        </p>

        <div class="bg-slate-50 border border-slate-200 rounded-md p-2.5 flex items-center justify-between">
          <code class="text-xs text-sky-700 font-mono truncate">
            http://{window.location.hostname}:8081/api/v1/provisioning/&lt;code&gt;
          </code>
          <button 
            type="button" 
            on:click={copyUrl} 
            class="text-xs text-slate-500 hover:text-slate-800 ml-2 font-medium"
          >
            {copyFeedback ? 'Copied!' : 'Copy'}
          </button>
        </div>

        <p class="text-[11px] text-slate-500">
          Point DHCP Option 66 to this URL for plug-and-play phone setup on your local voice VLAN.
        </p>
      </div>
    </Card>
  </div>
</div>
