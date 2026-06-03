<script lang="ts">
  import { onMount } from 'svelte';
  let cameras = [];
  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/cameras`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        cameras = (data.cameras || []).map(c => ({
          id: c.ID.substring(0,8), name: c.Name, sipCode: '*9XX', status: 'Online', recording: false, motion: false
        }));
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">IP Camera Feeds</h1>
      <p class="text-gray-400 font-medium text-sm">Monitor H.264 SIP-enabled security cameras and door phones.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-gray-800 hover:bg-gray-700 text-white px-4 py-2 rounded-xl font-bold transition-colors border border-gray-700 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
        Grid View
      </button>
      <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Add Camera
      </button>
    </div>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
    {#each cameras as cam}
      <div class="bg-gray-900/60 rounded-[2rem] border border-gray-800 overflow-hidden shadow-2xl group hover:border-indigo-500/30 transition-colors flex flex-col">
        <!-- Video Feed Wrapper -->
        <div class="relative bg-gray-950 aspect-video flex items-center justify-center overflow-hidden border-b border-gray-800">
          {#if cam.status === 'Online'}
            <!-- Simulated Noise/Static Background for "Feed" -->
            <div class="absolute inset-0 opacity-10 mix-blend-screen bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0IiBoZWlnaHQ9IjQiPgo8cmVjdCB3aWR0aD0iNCIgaGVpZ2h0PSI0IiBmaWxsPSIjMDAwIj48L3JlY3Q+CjxyZWN0IHdpZHRoPSIxIiBoZWlnaHQ9IjEiIGZpbGw9IiNmZmYiPjwvcmVjdD4KPHJlY3QgeD0iMiIgeT0iMiIgd2lkdGg9IjEiIGhlaWdodD0iMSIgZmlsbD0iI2ZmZiI+PC9yZWN0Pgo8L3N2Zz4=')]"></div>
            
            <!-- Timestamp Overlay -->
            <div class="absolute top-4 right-4 bg-black/60 backdrop-blur font-mono text-[10px] text-white px-2 py-1 rounded">
              2026-06-03 14:02:12
            </div>

            <!-- Recording/Motion Indicators -->
            <div class="absolute top-4 left-4 flex space-x-2">
              {#if cam.recording}
                <div class="flex items-center space-x-1 bg-black/60 backdrop-blur px-2 py-1 rounded">
                  <div class="w-1.5 h-1.5 rounded-full bg-rose-500 animate-pulse"></div>
                  <span class="text-[9px] font-bold text-white uppercase tracking-wider">REC</span>
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
            <svg class="w-8 h-8 text-white/10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 4v16m8-8H4"></path></svg>

          {:else}
            <!-- Offline State -->
            <div class="flex flex-col items-center">
              <svg class="w-12 h-12 text-rose-500/50 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 010 12.728m0 0l-2.829-2.829m2.829 2.829L21 21M15.536 8.464a5 5 0 010 7.072m0 0l-2.829-2.829m-4.243 2.829a4.978 4.978 0 01-1.414-2.83m-1.414 5.658a9 9 0 01-2.167-9.238m7.824 2.167a1 1 0 111.414 1.414m-1.414-1.414L3 3m8.293 8.293l1.414 1.414"></path></svg>
              <span class="text-xs font-bold text-rose-500 uppercase tracking-widest">Feed Offline</span>
            </div>
          {/if}
        </div>

        <!-- Camera Details -->
        <div class="p-5 flex-1 flex justify-between items-center">
          <div>
            <div class="flex items-center space-x-2 mb-1">
              <h3 class="text-lg font-bold text-white">{cam.name}</h3>
              <span class="px-2 py-0.5 rounded bg-gray-800 text-[10px] text-gray-400 font-mono border border-gray-700">{cam.id}</span>
            </div>
            <p class="text-sm text-gray-500 flex items-center">
              Dial Code: <span class="font-bold font-mono text-indigo-400 ml-1.5">{cam.sipCode}</span>
            </p>
          </div>
          
          <button aria-label="Open Door / Unlock" class="w-12 h-12 rounded-full bg-gray-800 hover:bg-gray-700 border border-gray-700 text-gray-400 hover:text-emerald-400 transition-colors flex items-center justify-center">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"></path></svg>
          </button>
        </div>
      </div>
    {/each}
  </div>
</div>

