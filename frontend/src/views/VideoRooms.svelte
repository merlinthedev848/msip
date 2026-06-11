<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';

  let rooms = [];
  let isModalOpen = false;
  let newName = '';
  let newMaxUsers = 10;
  let newCode = '';

  async function fetchRooms() {
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
  }

  onMount(fetchRooms);

  async function handleCreateRoom(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/video-rooms`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: newName, RoomCode: newCode, MaxUsers: parseInt(newMaxUsers), IsActive: true })
      });
      if (res.ok) {
        isModalOpen = false;
        newName = ''; newCode = '';
        fetchRooms();
      }
    } catch (e) {}
  }

  async function handleDeleteRoom(id) {
    if(!confirm("Delete this video room?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/video-rooms/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchRooms();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">WebRTC Video Rooms</h1>
      <p class="text-slate-500 font-medium text-sm">Manage multi-party video conferencing bridges.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
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
          <button aria-label="Delete Room" class="p-2 text-red-500 hover:text-slate-900 transition-colors rounded-lg hover:bg-red-800" on:click={() => handleDeleteRoom(room.id)}><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
          {#if room.active}
            <button aria-label="End Meeting" class="p-2 text-rose-500 hover:text-slate-900 hover:bg-rose-500 transition-colors rounded-lg"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg></button>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create Video Room">
  <form on:submit={handleCreateRoom} class="space-y-4">
    <div>
      <label for="room_name" class="block text-sm font-bold text-slate-500 mb-1">Room Name</label>
      <input id="room_name" type="text" bind:value={newName} required placeholder="e.g. Daily Standup" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label for="room_code" class="block text-sm font-bold text-slate-500 mb-1">Room Code / PIN</label>
      <input id="room_code" type="text" bind:value={newCode} required placeholder="e.g. 555222" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label for="room_max_users" class="block text-sm font-bold text-slate-500 mb-1">Max Users</label>
      <input id="room_max_users" type="number" bind:value={newMaxUsers} required min="2" max="100" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Create Room</button>
    </div>
  </form>
</Modal>
