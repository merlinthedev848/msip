<script lang="ts">
  import { onMount } from 'svelte';
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';

  let voicemails = [];
  async function fetchVoicemails() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/voicemails`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        voicemails = data.voicemails || [];
      }
    } catch (e) {}
  }

  onMount(fetchVoicemails);

  async function handleDelete(id) {
    if(!confirm("Delete this voicemail?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/voicemails/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchVoicemails();
    } catch (e) {}
  }

  async function handleSimulate() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/voicemails`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ CallerID: 'Test Caller', Extension: '1001', Duration: 15, IsRead: false })
      });
      if (res.ok) fetchVoicemails();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">Visual Voicemail</h1>
      <p class="text-slate-500 text-sm">Listen and read transcriptions of incoming messages.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-4 py-2 rounded-xl font-bold text-sm transition-colors" on:click={handleSimulate}>
      Simulate Voicemail
    </button>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    {#if voicemails.length === 0}
      <div class="col-span-1 md:col-span-2 p-8 text-center text-slate-500 bg-white/40 border border-slate-200 rounded-3xl">No voicemails found.</div>
    {/if}
    {#each voicemails as vm}
    <GlassCard hoverEffect={true}>
      <div class="flex justify-between items-start mb-4">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-full bg-indigo-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
          </div>
          <div>
            <p class="font-semibold text-slate-900">{vm.CallerID}</p>
            <p class="text-xs text-slate-500">Extension: {vm.Extension} • Duration: {vm.Duration}s</p>
          </div>
        </div>
        {#if !vm.IsRead}
          <Badge status="error">New</Badge>
        {/if}
      </div>
      
      <div class="bg-white p-4 rounded-lg border border-slate-200 mb-4">
        <p class="text-sm text-slate-700 italic font-serif leading-relaxed">[Voicemail transcription unavailable]</p>
      </div>
      
      <div class="flex items-center space-x-4">
        <button title="Play Voicemail" aria-label="Play Voicemail" class="w-10 h-10 rounded-full bg-indigo-600 hover:bg-indigo-500 flex items-center justify-center transition-colors shadow-lg shadow-indigo-500/30">
          <svg class="w-4 h-4 text-slate-900 ml-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd"></path></svg>
        </button>
        <div class="flex-1 h-2 bg-slate-100 rounded-full overflow-hidden">
          <div class="h-full bg-indigo-500 w-1/3"></div>
        </div>
        <span class="text-xs text-slate-500 font-mono">0:00 / {vm.Duration}</span>
        <button title="Delete Voicemail" aria-label="Delete Voicemail" class="text-rose-500 hover:text-slate-900 p-2 rounded hover:bg-rose-500 ml-4 transition-colors" on:click={() => handleDelete(vm.ID)}>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
        </button>
      </div>
    </GlassCard>
    {/each}
  </div>
</div>
