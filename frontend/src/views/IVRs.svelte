<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';
  import Modal from '../components/ui/Modal.svelte';
  import { onMount } from 'svelte';

  let ivrs = [];
  let isApiConnected = true;
  
  let isModalOpen = false;
  let newIvrName = '';
  let newIvrExtension = '';
  let newIvrGreeting = 'phrase:greeting';
  let isSubmitting = false;

  async function fetchIVRs() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/ivrs`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        ivrs = data.ivrs || [];
        isApiConnected = true;
      } else { throw new Error("API Offline"); }
    } catch (e) {
      isApiConnected = false;
      ivrs = [];
    }
  }

  onMount(() => {
    fetchIVRs();
  });

  async function handleCreateIVR(e) {
    e.preventDefault();
    isSubmitting = true;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/ivrs`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({
          Name: newIvrName,
          Extension: newIvrExtension,
          Greeting: newIvrGreeting
        })
      });
      
      if (res.ok) {
        isModalOpen = false;
        newIvrName = '';
        newIvrExtension = '';
        await fetchIVRs();
      } else {
        alert("Failed to create IVR.");
      }
    } catch (err) {
      alert("API Error");
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">Auto Attendants (IVR)</h1>
      <p class="text-slate-500 text-sm">Configure voice menus and digit routing.</p>
    </div>
    <Button variant="primary" on:click={() => isModalOpen = true}>New Menu</Button>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    {#each ivrs as ivr}
      <GlassCard>
        <div class="flex justify-between items-start mb-6">
          <div>
            <h3 class="text-xl font-bold text-slate-900">{ivr.Name}</h3>
            <p class="text-sm font-mono text-blue-600">Ext: {ivr.Extension}</p>
          </div>
          <Badge status="success">Active</Badge>
        </div>
        <div class="space-y-4">
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Greeting</span>
            <span class="text-slate-700 font-mono text-xs">{ivr.Greeting}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Configured Choices</span>
            <span class="text-slate-700">{ivr.Choices ? ivr.Choices.length : 0} options</span>
          </div>
        </div>
        <div class="mt-6 pt-4 border-t border-slate-200 flex justify-end space-x-3">
          <Button variant="secondary" className="px-4 text-xs">Edit Menu</Button>
        </div>
      </GlassCard>
    {/each}
    {#if ivrs.length === 0}
      <div class="col-span-full p-8 text-center text-slate-500 border border-dashed border-slate-200 rounded-2xl">
        No IVR Menus configured. Click "New Menu" to get started.
      </div>
    {/if}
  </div>
</div>

<!-- Add IVR Modal -->
<Modal bind:isOpen={isModalOpen} title="Create Auto Attendant" description="Create a new interactive voice menu.">
  <form on:submit={handleCreateIVR} class="space-y-5">
    <div>
      <label for="ivrName" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Menu Name</label>
      <input id="ivrName" type="text" bind:value={newIvrName} required placeholder="e.g. Main Menu" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 outline-none">
    </div>
    <div>
      <label for="ivrExt" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Virtual Extension</label>
      <input id="ivrExt" type="text" bind:value={newIvrExtension} required placeholder="e.g. 8000" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 outline-none">
    </div>
    
    <div class="pt-4 flex justify-end space-x-3">
      <Button variant="secondary" on:click={() => isModalOpen = false} type="button">Cancel</Button>
      <button type="submit" disabled={isSubmitting} class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2.5 rounded-xl font-bold transition-all shadow-lg disabled:opacity-50">
        {isSubmitting ? 'Saving...' : 'Create Menu'}
      </button>
    </div>
  </form>
</Modal>
