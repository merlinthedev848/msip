<script lang="ts">
  import { onMount } from 'svelte';
  let channels = [];
  let messages = [];
  import Modal from '../components/ui/Modal.svelte';
  let isModalOpen = false;
  let newName = '';

  async function fetchChats() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/chats`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        channels = (data.chats || []).map(c => ({
          id: c.ID, name: c.Name, unread: 0, active: false
        }));
        if (channels.length > 0) channels[0].active = true;
      }
    } catch (e) {}
  }

  onMount(fetchChats);

  async function handleCreateChat(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/chats`, {
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
        fetchChats();
      }
    } catch (e) {}
  }

  async function handleDeleteChat(id) {
    if(!confirm("Delete this channel?")) return;
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/chats/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchChats();
    } catch (e) {}
  }
</script>

<div class="h-[calc(100vh-120px)] flex flex-col animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-6 shrink-0">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Team Chat</h1>
      <p class="text-slate-500 font-medium text-sm">Persistent encrypted XMPP messaging channels.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      New Channel
    </button>
  </header>

  <div class="flex-1 flex space-x-6 min-h-0">
    <!-- Channels Sidebar -->
    <div class="w-72 bg-white/60 rounded-[2rem] border border-slate-200 shadow-2xl flex flex-col overflow-hidden shrink-0">
      <div class="p-5 border-b border-slate-200 bg-slate-50/30">
        <h3 class="text-xs font-bold text-slate-500 uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
          Active Channels
        </h3>
      </div>
      <div class="space-y-2 p-3">
        {#each channels as chan}
          <div class="w-full flex items-center justify-between p-3 rounded-xl transition-colors group {chan.active ? 'bg-indigo-600/10 border border-blue-100' : 'hover:bg-slate-50 border border-transparent'}">
            <div class="flex items-center">
              <span class="text-lg font-bold text-slate-500 mr-2">#</span>
              <span class="font-bold text-sm {chan.active ? 'text-blue-600' : 'text-slate-700 group-hover:text-slate-900'}">{chan.name}</span>
            </div>
            <div class="flex items-center space-x-2">
              {#if chan.unread > 0}
                <span class="bg-rose-500 text-slate-900 text-[10px] font-bold px-2 py-0.5 rounded-full shadow-[0_0_10px_rgba(244,63,94,0.5)]">{chan.unread}</span>
              {/if}
              <button title="Delete Channel" aria-label="Delete Channel" class="text-rose-500 hover:text-slate-900 p-1 rounded transition-colors opacity-0 group-hover:opacity-100" on:click|stopPropagation={() => handleDeleteChat(chan.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
            </div>
          </div>
        {/each}
      </div>
    </div>

    <!-- Chat Area -->
    <div class="flex-1 bg-white/60 rounded-[2rem] border border-slate-200 shadow-2xl flex flex-col overflow-hidden">
      <!-- Chat Header -->
      <div class="p-5 border-b border-slate-200 bg-slate-50/30 flex justify-between items-center shrink-0">
        <div class="flex items-center space-x-3">
          <span class="text-2xl font-black text-indigo-500">#</span>
          <h2 class="text-lg font-bold text-slate-900">general</h2>
          <span class="px-2 py-0.5 rounded-md bg-slate-100 border border-slate-200 text-[10px] text-slate-500 font-bold uppercase">14 Online</span>
        </div>
        <div class="flex space-x-2">
          <button aria-label="Search Chat" class="p-2 text-slate-500 hover:text-slate-900 transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg></button>
          <button aria-label="Channel Settings" class="p-2 text-slate-500 hover:text-slate-900 transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg></button>
        </div>
      </div>

      <!-- Messages -->
      <div class="flex-1 p-5 overflow-y-auto space-y-6">
        {#each messages as msg}
          {#if msg.isSys}
            <div class="flex justify-center">
              <span class="text-[10px] font-bold text-slate-500 uppercase tracking-widest bg-slate-50 px-3 py-1 rounded-full">{msg.text}</span>
            </div>
          {:else}
            <div class="flex items-start space-x-4">
              <div class="w-10 h-10 rounded-xl bg-slate-100 flex items-center justify-center font-black {msg.color} border border-slate-200 shrink-0">
                {msg.user.charAt(0)}
              </div>
              <div class="flex-1">
                <div class="flex items-baseline space-x-2 mb-1">
                  <span class="font-bold text-slate-900 text-sm">{msg.user}</span>
                  <span class="text-[10px] font-bold text-slate-500">{msg.time}</span>
                </div>
                <p class="text-slate-700 text-sm leading-relaxed">{msg.text}</p>
              </div>
            </div>
          {/if}
        {/each}
      </div>

      <!-- Input Box -->
      <div class="p-4 border-t border-slate-200 bg-slate-50/50 shrink-0">
        <div class="relative flex items-center">
          <button aria-label="Attach File" class="absolute left-3 text-slate-500 hover:text-blue-600 transition-colors p-1"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"></path></svg></button>
          <input type="text" placeholder="Message #general..." class="w-full bg-white border border-slate-200 text-slate-900 rounded-xl pl-12 pr-12 py-4 focus:ring-indigo-500 focus:border-indigo-500 transition-colors placeholder-gray-600 font-medium">
          <button aria-label="Send Message" class="absolute right-3 bg-indigo-600 hover:bg-indigo-500 text-slate-900 p-2 rounded-lg shadow-lg shadow-indigo-500/30 transition-all"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path></svg></button>
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create Channel">
  <form on:submit={handleCreateChat} class="space-y-4">
    <div>
      <label for="chat_name" class="block text-sm font-bold text-slate-500 mb-1">Channel Name</label>
      <input id="chat_name" type="text" bind:value={newName} required placeholder="e.g. general" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Create</button>
    </div>
  </form>
</Modal>
