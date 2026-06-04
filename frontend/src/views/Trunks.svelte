<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import { onMount } from 'svelte';

  let trunks = [];
  let isApiConnected = true;
  
  let isModalOpen = false;
  let newTrunkName = '';
  let newTrunkServer = '';
  let newTrunkAuth = 'IP_ACL';
  let isSubmitting = false;

  async function fetchTrunks() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/trunks`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        trunks = data.trunks || [];
        isApiConnected = true;
      } else { throw new Error("API Offline"); }
    } catch (e) {
      isApiConnected = false;
      trunks = [];
    }
  }

  onMount(() => {
    fetchTrunks();
  });

  async function handleCreateTrunk(e) {
    e.preventDefault();
    isSubmitting = true;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/trunks`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({
          ProviderName: newTrunkName,
          SIPServer: newTrunkServer,
          AuthMethod: newTrunkAuth
        })
      });
      
      if (res.ok) {
        isModalOpen = false;
        newTrunkName = '';
        newTrunkServer = '';
        await fetchTrunks();
      } else {
        alert("Failed to create trunk.");
      }
    } catch (err) {
      alert("API Error: Make sure backend is running.");
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDeleteTrunk(id) {
    if(!confirm("Are you sure you want to delete this trunk?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/trunks/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchTrunks();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">SIP Trunks</h1>
      <p class="text-slate-500 text-sm">Manage external carrier connections.</p>
    </div>
    <Button variant="primary" on:click={() => isModalOpen = true}>Add Trunk</Button>
  </header>

  {#if !isApiConnected}
  <div class="mb-6 p-4 bg-yellow-500/10 border border-yellow-500/30 rounded-lg flex items-start space-x-3 text-yellow-500 text-sm">
    <svg class="w-5 h-5 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
    <div>
      <p class="font-bold">Backend API Offline</p>
      <p class="opacity-80 mt-1">The Go backend (port 8081) is not responding.</p>
    </div>
  </div>
  {/if}

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Static Fallback / Examples if DB is empty -->
    {#if trunks.length === 0}
      <GlassCard>
        <div class="flex justify-between items-start mb-6">
          <div>
            <h3 class="text-xl font-bold text-slate-900">Twilio Primary</h3>
            <p class="text-sm text-slate-500">sip.twilio.com</p>
          </div>
          <Badge status="success" pulse={true}>Registered</Badge>
        </div>
        <div class="space-y-4">
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Auth Method</span>
            <span class="text-slate-700">IP ACL</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Channels</span>
            <span class="text-slate-700">Unlimited</span>
          </div>
        </div>
      </GlassCard>
    {/if}

    <!-- Dynamic Trunks from DB -->
    {#each trunks as trunk}
      <GlassCard>
        <div class="flex justify-between items-start mb-6">
          <div>
            <h3 class="text-xl font-bold text-slate-900">{trunk.ProviderName}</h3>
            <p class="text-sm text-slate-500">{trunk.SIPServer}</p>
          </div>
          <Badge status="success" pulse={true}>Registered</Badge>
        </div>
        <div class="space-y-4">
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Auth Method</span>
            <span class="text-slate-700">{trunk.AuthMethod}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Database ID</span>
            <span class="text-slate-700 font-mono text-xs">{trunk.ID.substring(0, 8)}...</span>
          </div>
        </div>
        <div class="mt-6 pt-4 border-t border-slate-200 flex justify-end space-x-3">
          <Button variant="secondary" className="px-4 text-xs">Configure</Button>
          <Button variant="danger" className="px-4 text-xs" on:click={() => handleDeleteTrunk(trunk.ID)}>Delete</Button>
        </div>
      </GlassCard>
    {/each}
  </div>
</div>

<!-- Add Trunk Modal -->
<Modal bind:isOpen={isModalOpen} title="Add SIP Trunk" description="Connect to an external carrier or upstream PBX.">
  <form on:submit={handleCreateTrunk} class="space-y-5">
    <div>
      <label for="trunkName" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Provider Name</label>
      <input id="trunkName" type="text" bind:value={newTrunkName} required placeholder="e.g. Twilio, Flowroute" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none">
    </div>
    <div>
      <label for="trunkServer" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">SIP Server (IP/Domain)</label>
      <input id="trunkServer" type="text" bind:value={newTrunkServer} required placeholder="e.g. sip.provider.com" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none">
    </div>
    <div>
      <label for="trunkAuth" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Authentication Method</label>
      <select id="trunkAuth" bind:value={newTrunkAuth} class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none appearance-none">
        <option value="IP_ACL">IP Authentication (ACL)</option>
        <option value="USERPASS">Username / Password</option>
      </select>
    </div>
    
    <div class="pt-4 flex justify-end space-x-3">
      <Button variant="secondary" on:click={() => isModalOpen = false} type="button">Cancel</Button>
      <button type="submit" disabled={isSubmitting} class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2.5 rounded-xl font-bold transition-all shadow-lg shadow-indigo-600/20 disabled:opacity-50">
        {isSubmitting ? 'Saving...' : 'Add Trunk'}
      </button>
    </div>
  </form>
</Modal>


