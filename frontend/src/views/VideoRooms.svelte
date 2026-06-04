<script lang="ts">
  import { onMount } from 'svelte';
  let rooms = [];
  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/video-rooms`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        rooms = (data.video_rooms || []).map(r => ({
          id: r.RoomCode, name: r.Name, users: 0, max: r.MaxUsers, duration: '0m', host: 'System', active: r.IsActive
        }));
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">WebRTC Video Rooms</h1>
      <p class="text-slate-500 font-medium text-sm">Manage multi-party video conferencing bridges.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center">
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
      Create Room
    </button>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    {#each rooms as room}
      <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl flex flex-col group hover:border-indigo-500/30 transition-colors">
        <!-- Room Header / Preview -->
        <div class="h-40 bg-slate-50 relative flex items-center justify-center overflow-hidden border-b border-slate-200">
          {#if room.active}
            <div class="absolute inset-0 grid grid-cols-2 gap-1 p-1 opacity-20">
              <div class="bg-indigo-500 rounded-lg"></div>
              <div class="bg-emerald-500 rounded-lg"></div>
              <div class="bg-rose-500 rounded-lg"></div>
              <div class="bg-amber-500 rounded-lg"></div>
            </div>
            <div class="absolute inset-0 bg-gradient-to-t from-gray-950 to-transparent"></div>
            
            <button class="relative z-10 bg-indigo-600 hover:bg-indigo-500 text-slate-900 font-bold py-3 px-6 rounded-xl shadow-xl shadow-indigo-600/20 flex items-center transform transition-transform hover:scale-105 active:scale-95">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
              Join Call
            </button>
            
            <div class="absolute top-4 right-4 flex items-center space-x-1.5 bg-white/80 backdrop-blur px-2.5 py-1 rounded-full border border-slate-200">
              <div class="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.8)] animate-pulse"></div>
              <span class="text-[10px] font-bold text-slate-900 tracking-widest uppercase">Live</span>
            </div>
          {:else}
            <div class="flex flex-col items-center justify-center opacity-50">
              <svg class="w-12 h-12 text-gray-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              <span class="text-xs font-bold text-slate-500 uppercase tracking-widest">Room Empty</span>
            </div>
          {/if}
          
          <div class="absolute bottom-4 left-4 bg-white/80 backdrop-blur px-3 py-1.5 rounded-xl border border-slate-200 font-mono text-xs font-bold text-slate-700">
            {room.duration}
          </div>
        </div>

        <!-- Room Details -->
        <div class="p-6 flex-1 flex flex-col justify-between">
          <div class="mb-6">
            <div class="flex items-center space-x-3 mb-2">
              <span class="text-sm font-black text-indigo-500 font-mono bg-blue-50 px-2 py-0.5 rounded border border-blue-100">{room.id}</span>
              <h3 class="text-xl font-bold text-slate-900">{room.name}</h3>
            </div>
            <p class="text-sm text-slate-500 flex items-center">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
              Host: <span class="font-bold text-slate-700 ml-1">{room.host}</span>
            </p>
          </div>

          <div>
            <div class="flex justify-between items-end mb-2">
              <span class="text-xs font-bold text-slate-500 uppercase tracking-widest">Participants</span>
              <span class="text-sm font-mono font-bold {room.users > 0 ? 'text-blue-600' : 'text-gray-600'}">{room.users} <span class="text-gray-600">/ {room.max}</span></span>
            </div>
            <div class="w-full h-2 bg-slate-100 rounded-full overflow-hidden">
              <div class="h-full {room.users/room.max > 0.8 ? 'bg-rose-500' : room.users > 0 ? 'bg-indigo-500' : 'bg-transparent'}" style="width: {(room.users/room.max)*100}%"></div>
            </div>
          </div>
        </div>
        
        <!-- Room Actions -->
        <div class="p-4 border-t border-slate-200 bg-slate-50/30 flex justify-end space-x-2">
          <button aria-label="Copy Invite Link" class="p-2 text-slate-500 hover:text-slate-900 transition-colors rounded-lg hover:bg-slate-100"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path></svg></button>
          <button aria-label="Room Settings" class="p-2 text-slate-500 hover:text-slate-900 transition-colors rounded-lg hover:bg-slate-100"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
          {#if room.active}
            <button aria-label="End Meeting" class="p-2 text-rose-500 hover:text-slate-900 hover:bg-rose-500 transition-colors rounded-lg"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg></button>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

