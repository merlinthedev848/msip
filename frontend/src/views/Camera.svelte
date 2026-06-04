<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let cameras = [];
  let isModalOpen = false;
  let newName = '';
  let newSipCode = '*9XX';

  async function fetchCameras() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/cameras`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/cameras`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/cameras/${id}`, {
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
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">IP Camera Feeds</h1>
      <p class="text-slate-500 font-medium text-sm">Monitor H.264 SIP-enabled security cameras and door phones.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-slate-100 hover:bg-slate-200 text-slate-900 px-4 py-2 rounded-xl font-bold transition-colors border border-slate-200 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
        Grid View
      </button>
      <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm" on:click={() => isModalOpen = true}>
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Add Camera
      </button>
    </div>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
    {#each cameras as cam}
      <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl group hover:border-indigo-500/30 transition-colors flex flex-col">
        <!-- Video Feed Wrapper -->
        <div class="relative bg-slate-50 aspect-video flex items-center justify-center overflow-hidden border-b border-slate-200">
          {#if cam.status === 'Online'}
            <!-- Simulated Noise/Static Background for "Feed" -->
            <div class="absolute inset-0 opacity-10 mix-blend-screen bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0IiBoZWlnaHQ9IjQiPgo8cmVjdCB3aWR0aD0iNCIgaGVpZ2h0PSI0IiBmaWxsPSIjMDAwIj48L3JlY3Q+CjxyZWN0IHdpZHRoPSIxIiBoZWlnaHQ9IjEiIGZpbGw9IiNmZmYiPjwvcmVjdD4KPHJlY3QgeD0iMiIgeT0iMiIgd2lkdGg9IjEiIGhlaWdodD0iMSIgZmlsbD0iI2ZmZiI+PC9yZWN0Pgo8L3N2Zz4=')]"></div>
            
            <!-- Timestamp Overlay -->
            <div class="absolute top-4 right-4 bg-black/60 backdrop-blur font-mono text-[10px] text-slate-900 px-2 py-1 rounded">
              2026-06-03 14:02:12
            </div>

            <!-- Recording/Motion Indicators -->
            <div class="absolute top-4 left-4 flex space-x-2">
              {#if cam.recording}
                <div class="flex items-center space-x-1 bg-black/60 backdrop-blur px-2 py-1 rounded">
                  <div class="w-1.5 h-1.5 rounded-full bg-rose-500 animate-pulse"></div>
                  <span class="text-[9px] font-bold text-slate-900 uppercase tracking-wider">REC</span>
                </div>
              {/if}
              {#if cam.motion}
                <div class="flex items-center space-x-1 bg-black/60 backdrop-blur px-2 py-1 rounded border border-amber-500/50">
                  <svg class="w-3 h-3 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
                  <span class="text-[9px] font-bold text-amber-500 uppercase tracking-wider">Motion</span>
                </div>
              {/if}
            </div>

            <!-- Crosshair -->
            <svg class="w-8 h-8 text-slate-900/10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 4v16m8-8H4"></path></svg>

          {:else}
            <!-- Offline State -->
            <button class="p-2 text-slate-500 hover:text-slate-900 hover:bg-slate-100 rounded-lg transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
          <button class="p-2 text-rose-500 hover:text-slate-900 hover:bg-rose-500 rounded-lg transition-colors" on:click={() => handleDeleteCamera(cam.id)}><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
          {/if}
        </div>

        <!-- Camera Details -->
        <div class="p-5 flex-1 flex justify-between items-center">
          <div>
            <div class="flex items-center space-x-2 mb-1">
              <h3 class="text-lg font-bold text-slate-900">{cam.name}</h3>
              <span class="px-2 py-0.5 rounded bg-slate-100 text-[10px] text-slate-500 font-mono border border-slate-200">{cam.id}</span>
            </div>
            <p class="text-sm text-slate-500 flex items-center">
              Dial Code: <span class="font-bold font-mono text-blue-600 ml-1.5">{cam.sipCode}</span>
            </p>
          </div>
          
          <button aria-label="Open Door / Unlock" class="w-12 h-12 rounded-full bg-slate-100 hover:bg-slate-200 border border-slate-200 text-slate-500 hover:text-emerald-400 transition-colors flex items-center justify-center">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"></path></svg>
          </button>
        </div>
      </div>
    {/each}
  </div>
</div>


