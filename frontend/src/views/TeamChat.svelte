<script lang="ts">
  import { onMount } from 'svelte';
  let channels = [];
  let messages = [];
  import Modal from '../components/ui/Modal.svelte';
  let isModalOpen = false;
  let newName = '';

  async function fetchChats() {
    try {
      const res = await fetch(`/api/v1/chats`, {
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
      const res = await fetch(`/api/v1/chats`, {
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
      const res = await fetch(`/api/v1/chats/${id}`, {
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
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Team Chat</h1>
      <p class="text-slate-500 font-medium text-sm">Persistent encrypted XMPP messaging channels.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      New Channel
    </button>
  </header>

  <div class="flex-1 flex space-x-6 min-h-0">
    <!-- Channels Sidebar -->
    <div class="w-72 bg-[#121622] rounded-xl border border-[#202738] shadow-2xl flex flex-col overflow-hidden shrink-0">
      <div class="p-4 border-b border-[#1E2538] bg-[#0E121B]">
        <h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider flex items-center font-mono">
          <svg class="w-4 h-4 mr-2 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
          Active Channels
        </h3>
      </div>
      <div class="space-y-1.5 p-3">
        {#each channels as chan}
          <div class="w-full flex items-center justify-between p-2.5 rounded-lg transition-colors group {chan.active ? 'bg-cyan-500/10 border border-cyan-500/30' : 'hover:bg-[#182030] border border-transparent'}">
            <div class="flex items-center">
              <span class="text-base font-bold text-slate-500 mr-2 font-mono">#</span>
              <span class="font-bold text-xs font-mono {chan.active ? 'text-cyan-400' : 'text-slate-300 group-hover:text-white'}">{chan.name}</span>
            </div>
            <div class="flex items-center space-x-2">
              {#if chan.unread > 0}
                <span class="bg-rose-600 text-white text-[10px] font-bold px-1.5 py-0.2 rounded font-mono">{chan.unread}</span>
              {/if}
              <button title="Delete Channel" aria-label="Delete Channel" class="text-rose-400 hover:text-white p-1 rounded transition-colors opacity-0 group-hover:opacity-100" on:click|stopPropagation={() => handleDeleteChat(chan.id)}><svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
            </div>
          </div>
        {/each}
      </div>
    </div>

    <!-- Chat Area -->
    <div class="flex-1 bg-[#121622] rounded-xl border border-[#202738] shadow-2xl flex flex-col overflow-hidden">
      <!-- Chat Header -->
      <div class="p-4 border-b border-[#1E2538] bg-[#0E121B] flex justify-between items-center shrink-0">
        <div class="flex items-center space-x-3">
          <span class="text-xl font-black text-cyan-400 font-mono">#</span>
          <h2 class="text-sm font-bold text-white font-mono">general</h2>
          <span class="px-2 py-0.5 rounded bg-[#182030] border border-[#26334A] text-[10px] text-emerald-400 font-mono font-bold uppercase">14 Online</span>
        </div>
      </div>

      <!-- Messages -->
      <div class="flex-1 p-5 overflow-y-auto space-y-5">
        {#each messages as msg}
          {#if msg.isSys}
            <div class="flex justify-center">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest bg-[#0B0E17] border border-[#1E2538] px-3 py-1 rounded-full font-mono">{msg.text}</span>
            </div>
          {:else}
            <div class="flex items-start space-x-3">
              <div class="w-8 h-8 rounded-lg bg-[#182030] flex items-center justify-center font-bold text-cyan-400 border border-[#26334A] shrink-0 font-mono text-xs">
                {msg.user.charAt(0)}
              </div>
              <div class="flex-1">
                <div class="flex items-baseline space-x-2 mb-1">
                  <span class="font-bold text-white text-xs font-mono">{msg.user}</span>
                  <span class="text-[10px] font-mono text-slate-500">{msg.time}</span>
                </div>
                <p class="text-slate-300 text-xs leading-relaxed font-mono">{msg.text}</p>
              </div>
            </div>
          {/if}
        {/each}
      </div>

      <!-- Input Box -->
      <div class="p-3 border-t border-[#1E2538] bg-[#0E121B] shrink-0">
        <div class="relative flex items-center">
          <input type="text" placeholder="Message #general..." class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-4 py-2.5 focus:border-cyan-500 transition-colors placeholder-slate-500 font-mono text-xs outline-none">
          <button aria-label="Send Message" class="absolute right-2 bg-cyan-600 hover:bg-cyan-500 text-slate-950 p-1.5 rounded transition-all"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path></svg></button>
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Create Channel">
  <form on:submit={handleCreateChat} class="space-y-4">
    <div>
      <label for="chat_name" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Channel Name</label>
      <input id="chat_name" type="text" bind:value={newName} required placeholder="e.g. general" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Create</button>
    </div>
  </form>
</Modal>
