<script lang="ts">
  import { onMount } from 'svelte';
  let channels = [];
  let messages = [];

  onMount(async () => {
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
  });
</script>

<div class="h-[calc(100vh-120px)] flex flex-col animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-6 shrink-0">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Team Chat</h1>
      <p class="text-gray-400 font-medium text-sm">Persistent encrypted XMPP messaging channels.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center">
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      New Channel
    </button>
  </header>

  <div class="flex-1 flex space-x-6 min-h-0">
    <!-- Channels Sidebar -->
    <div class="w-72 bg-gray-900/60 rounded-[2rem] border border-gray-800 shadow-2xl flex flex-col overflow-hidden shrink-0">
      <div class="p-5 border-b border-gray-800 bg-gray-950/30">
        <h3 class="text-xs font-bold text-gray-400 uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
          Active Channels
        </h3>
      </div>
      <div class="space-y-2 p-3">
        {#each channels as chan}
          <button class="w-full flex items-center justify-between p-3 rounded-xl transition-colors group {chan.active ? 'bg-indigo-600/10 border border-indigo-500/20' : 'hover:bg-gray-800/40 border border-transparent'}">
            <div class="flex items-center">
              <span class="text-lg font-bold text-gray-500 mr-2">#</span>
              <span class="font-bold text-sm {chan.active ? 'text-indigo-400' : 'text-gray-300 group-hover:text-white'}">{chan.name}</span>
            </div>
            {#if chan.unread > 0}
              <span class="bg-rose-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-full shadow-[0_0_10px_rgba(244,63,94,0.5)]">{chan.unread}</span>
            {/if}
          </button>
        {/each}
      </div>
    </div>

    <!-- Chat Area -->
    <div class="flex-1 bg-gray-900/60 rounded-[2rem] border border-gray-800 shadow-2xl flex flex-col overflow-hidden">
      <!-- Chat Header -->
      <div class="p-5 border-b border-gray-800 bg-gray-950/30 flex justify-between items-center shrink-0">
        <div class="flex items-center space-x-3">
          <span class="text-2xl font-black text-indigo-500">#</span>
          <h2 class="text-lg font-bold text-white">general</h2>
          <span class="px-2 py-0.5 rounded-md bg-gray-800 border border-gray-700 text-[10px] text-gray-400 font-bold uppercase">14 Online</span>
        </div>
        <div class="flex space-x-2">
          <button aria-label="Search Chat" class="p-2 text-gray-500 hover:text-white transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg></button>
          <button aria-label="Channel Settings" class="p-2 text-gray-500 hover:text-white transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg></button>
        </div>
      </div>

      <!-- Messages -->
      <div class="flex-1 p-5 overflow-y-auto space-y-6">
        {#each messages as msg}
          {#if msg.isSys}
            <div class="flex justify-center">
              <span class="text-[10px] font-bold text-gray-500 uppercase tracking-widest bg-gray-800/50 px-3 py-1 rounded-full">{msg.text}</span>
            </div>
          {:else}
            <div class="flex items-start space-x-4">
              <div class="w-10 h-10 rounded-xl bg-gray-800 flex items-center justify-center font-black {msg.color} border border-gray-700 shrink-0">
                {msg.user.charAt(0)}
              </div>
              <div class="flex-1">
                <div class="flex items-baseline space-x-2 mb-1">
                  <span class="font-bold text-white text-sm">{msg.user}</span>
                  <span class="text-[10px] font-bold text-gray-500">{msg.time}</span>
                </div>
                <p class="text-gray-300 text-sm leading-relaxed">{msg.text}</p>
              </div>
            </div>
          {/if}
        {/each}
      </div>

      <!-- Input Box -->
      <div class="p-4 border-t border-gray-800 bg-gray-950/50 shrink-0">
        <div class="relative flex items-center">
          <button aria-label="Attach File" class="absolute left-3 text-gray-500 hover:text-indigo-400 transition-colors p-1"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"></path></svg></button>
          <input type="text" placeholder="Message #general..." class="w-full bg-gray-900 border border-gray-700 text-white rounded-xl pl-12 pr-12 py-4 focus:ring-indigo-500 focus:border-indigo-500 transition-colors placeholder-gray-600 font-medium">
          <button aria-label="Send Message" class="absolute right-3 bg-indigo-600 hover:bg-indigo-500 text-white p-2 rounded-lg shadow-lg shadow-indigo-500/30 transition-all"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path></svg></button>
        </div>
      </div>
    </div>
  </div>
</div>

