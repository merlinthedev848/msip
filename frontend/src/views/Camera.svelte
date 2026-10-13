<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let cameras = [];
  let isModalOpen = false;
  let newName = '';
  let newSipCode = '*9XX';

  async function fetchCameras() {
    try {
      const res = await fetch(`/api/v1/cameras`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        cameras = (data.cameras || []).map(c => ({
          id: c.ID, name: c.Name, sipCode: newSipCode, status: 'Online', recording: false, motion: false
        }));
      }
    } catch (e) {}
  }

  onMount(fetchCameras);

  async function handleCreateCamera(e) {
    e.preventDefault();
    try {
      const res = await fetch(`/api/v1/cameras`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: newName })
      });
      if (res.ok) {
        isModalOpen = false;
        newName = '';
        fetchCameras();
      }
    } catch (e) {}
  }

  async function handleDeleteCamera(id) {
    if(!confirm("Delete this camera?")) return;
    try {
      const res = await fetch(`/api/v1/cameras/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchCameras();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">IP Camera Feeds</h1>
      <p class="text-slate-500 font-medium text-sm">Monitor H.264 SIP-enabled security cameras and door phones.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-slate-100 hover:bg-slate-200 text-white px-4 py-2 rounded-xl font-bold transition-colors border border-slate-800 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
        Grid View
      </button>
      <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm" on:click={() => isModalOpen = true}>
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Add Camera
      </button>
    </div>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
    {#if cameras.length === 0}
      <div class="col-span-1 md:col-span-3 p-8 text-center text-slate-400 bg-[#121622] border border-[#202738] rounded-xl font-mono text-xs">No cameras configured.</div>
    {/if}
    {#each cameras as cam}
      <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl flex flex-col">
        <!-- Video Feed Wrapper -->
        <div class="relative bg-[#080A0F] aspect-video flex items-center justify-center overflow-hidden border-b border-[#1E2538]">
          {#if cam.status === 'Online'}
            <!-- Timestamp Overlay -->
            <div class="absolute top-3 right-3 bg-[#0B0E17] border border-[#1E2538] font-mono text-[10px] text-cyan-400 px-2 py-1 rounded">
              CAM STREAM ACTIVE
            </div>

            <!-- Recording/Motion Indicators -->
            <div class="absolute top-3 left-3 flex space-x-2">
              {#if cam.recording}
                <div class="flex items-center space-x-1.5 bg-[#0B0E17] border border-rose-500/40 px-2 py-1 rounded">
                  <div class="w-1.5 h-1.5 rounded-full bg-rose-500 animate-pulse"></div>
                  <span class="text-[9px] font-bold text-rose-400 uppercase tracking-wider font-mono">REC</span>
                </div>
              {/if}
              {#if cam.motion}
                <div class="flex items-center space-x-1 bg-[#0B0E17] border border-amber-500/50 px-2 py-1 rounded">
                  <svg class="w-3 h-3 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
                  <span class="text-[9px] font-bold text-amber-400 uppercase tracking-wider font-mono">Motion</span>
                </div>
              {/if}
            </div>

            <!-- Crosshair -->
            <svg class="w-8 h-8 text-cyan-500/20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 4v16m8-8H4"></path></svg>

          {:else}
            <!-- Offline State -->
            <div class="flex flex-col items-center justify-center p-6 text-slate-500">
              <span class="text-xs font-mono uppercase tracking-wider mb-2">Feed Offline</span>
              <button title="Delete Camera" aria-label="Delete Camera" class="p-1.5 text-rose-400 hover:text-white hover:bg-rose-600 rounded transition-colors" on:click={() => handleDeleteCamera(cam.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
            </div>
          {/if}
        </div>

        <!-- Camera Details -->
        <div class="p-4 flex-1 flex justify-between items-center bg-[#0E121B]">
          <div>
            <div class="flex items-center space-x-2 mb-1">
              <h3 class="text-sm font-bold text-white font-mono">{cam.name}</h3>
              <span class="px-1.5 py-0.5 rounded bg-[#182030] text-[10px] text-cyan-400 font-mono border border-[#26334A]">{cam.id ? cam.id.substring(0,8) : 'N/A'}</span>
            </div>
            <p class="text-xs text-slate-400 flex items-center font-mono">
              Dial Code: <span class="font-bold text-cyan-400 ml-1.5">{cam.sipCode}</span>
            </p>
          </div>
          
          <button aria-label="Open Door / Unlock" class="w-10 h-10 rounded-lg bg-[#182030] hover:bg-[#202B40] border border-[#26334A] text-slate-300 hover:text-emerald-400 transition-colors flex items-center justify-center">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"></path></svg>
          </button>
        </div>
      </div>
    {/each}
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create Camera Feed">
  <form on:submit={handleCreateCamera} class="space-y-4">
    <div>
      <label for="cam_name" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Camera Name</label>
      <input id="cam_name" type="text" bind:value={newName} required placeholder="e.g. Front Door Camera" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Create Camera</button>
    </div>
  </form>
</Modal>
