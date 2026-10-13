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
      const res = await fetch(`/api/v1/video-rooms`, {
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
      const res = await fetch(`/api/v1/video-rooms`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: newName, RoomCode: newCode, MaxUsers: Number(newMaxUsers), IsActive: true })
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
      const res = await fetch(`/api/v1/video-rooms/${id}`, {
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
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">WebRTC Video Rooms</h1>
      <p class="text-slate-500 font-medium text-sm">Manage multi-party video conferencing bridges.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
      Create Room
    </button>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    {#each rooms as room}
      <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl flex flex-col group hover:border-cyan-500/30 transition-colors">
        <!-- Room Header / Preview -->
        <div class="h-36 bg-[#080A0F] relative flex items-center justify-center overflow-hidden border-b border-[#1E2538]">
          {#if room.active}
            <div class="absolute inset-0 bg-[#0B0E17]"></div>
            
            <button class="relative z-10 bg-cyan-600 hover:bg-cyan-500 text-slate-950 font-bold py-2.5 px-5 rounded-lg flex items-center font-mono text-xs uppercase tracking-wider transition-all">
              <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
              Join Conference
            </button>
            
            <div class="absolute top-3 right-3 flex items-center space-x-1.5 bg-[#0E121B] text-emerald-400 px-2 py-0.5 rounded border border-emerald-500/30 font-mono text-[10px]">
              <div class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></div>
              <span class="font-bold tracking-wider uppercase">Active</span>
            </div>
          {:else}
            <div class="flex flex-col items-center justify-center opacity-60">
              <svg class="w-8 h-8 text-slate-500 mb-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              <span class="text-[10px] font-bold text-slate-500 uppercase tracking-widest font-mono">Conference Idle</span>
            </div>
          {/if}
          
          <div class="absolute bottom-3 left-3 bg-[#0E121B] text-slate-300 px-2 py-0.5 rounded border border-[#1E2538] font-mono text-[10px]">
            {room.duration}
          </div>
        </div>

        <!-- Room Details -->
        <div class="p-5 flex-1 flex flex-col justify-between bg-[#121622]">
          <div class="mb-4">
            <div class="flex items-center space-x-2.5 mb-1.5">
              <span class="text-xs font-bold text-cyan-400 font-mono bg-[#182030] px-2 py-0.5 rounded border border-[#26334A]">{room.id}</span>
              <h3 class="text-base font-bold text-white font-mono">{room.name}</h3>
            </div>
            <p class="text-xs text-slate-400 flex items-center font-mono">
              <svg class="w-3.5 h-3.5 mr-1.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
              Host: <span class="font-bold text-slate-200 ml-1">{room.host}</span>
            </p>
          </div>

          <div>
            <div class="flex justify-between items-end mb-1.5 font-mono">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Participants</span>
              <span class="text-xs font-bold {room.users > 0 ? 'text-cyan-400' : 'text-slate-500'}">{room.users} <span class="text-slate-500">/ {room.max}</span></span>
            </div>
            <div class="w-full h-1.5 bg-[#182030] rounded-full overflow-hidden border border-[#26334A]">
              <div class="h-full {room.users/room.max > 0.8 ? 'bg-rose-500' : room.users > 0 ? 'bg-cyan-500' : 'bg-transparent'}" style="width: {(room.users/room.max)*100}%"></div>
            </div>
          </div>
        </div>
        
        <!-- Room Actions -->
        <div class="p-3 border-t border-[#1E2538] bg-[#0E121B] flex justify-end space-x-1">
          <button aria-label="Copy Invite Link" class="p-1.5 text-slate-400 hover:text-white transition-colors rounded hover:bg-[#182030]"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path></svg></button>
          <button aria-label="Delete Room" class="p-1.5 text-rose-400 hover:text-white transition-colors rounded hover:bg-rose-600" on:click={() => handleDeleteRoom(room.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
        </div>
      </div>
    {/each}
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create Video Room">
  <form on:submit={handleCreateRoom} class="space-y-4">
    <div>
      <label for="room_name" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Room Name</label>
      <input id="room_name" type="text" bind:value={newName} required placeholder="e.g. Daily Standup" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="room_code" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Room Code / PIN</label>
      <input id="room_code" type="text" bind:value={newCode} required placeholder="e.g. 555222" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="room_max_users" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Max Users</label>
      <input id="room_max_users" type="number" bind:value={newMaxUsers} required min="2" max="100" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Create Room</button>
    </div>
  </form>
</Modal>
