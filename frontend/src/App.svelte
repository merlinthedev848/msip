<script lang="ts">
  import Sidebar from './components/Sidebar.svelte';
  import { onMount } from 'svelte';
  
  import Dashboard from './views/Dashboard.svelte';
  import CDRs from './views/CDRs.svelte';
  import Routing from './views/Routing.svelte';
  import OutboundRoutes from './views/OutboundRoutes.svelte';
  import Trunks from './views/Trunks.svelte';
  import Extensions from './views/Extensions.svelte';
  import IVR from './views/IVR.svelte';
  import IVRs from './views/IVRs.svelte';
  import Voicemail from './views/Voicemail.svelte';
  import VideoRooms from './views/VideoRooms.svelte';
  import Camera from './views/Camera.svelte';
  import SMSTrunks from './views/SMSTrunks.svelte';
  import TeamChat from './views/TeamChat.svelte';
  import APICredentials from './views/APICredentials.svelte';
  import Webhooks from './views/Webhooks.svelte';
  import Throttling from './views/Throttling.svelte';
  import Firewall from './views/Firewall.svelte';
  import AntiFraud from './views/AntiFraud.svelte';
  import VersionControl from './views/VersionControl.svelte';
  import Admins from './views/Admins.svelte';
  
  // Hybrid Features
  import Wallboard from './views/Wallboard.svelte';
  import Domains from './views/Domains.svelte';
  import Provisioning from './views/Provisioning.svelte';
  
  import Login from './views/Login.svelte';
  import Softphone from './components/Softphone.svelte';

  let currentPath = window.location.pathname;
  let token = localStorage.getItem('pbx_token');
  let user = JSON.parse(localStorage.getItem('pbx_user') || '{}');

  let searchQuery = '';
  let showSearchResults = false;

  const searchablePages = [
    { name: 'Dashboard & Metrics', path: '/', category: 'Overview' },
    { name: 'Extensions (SIP Endpoints)', path: '/extensions', category: 'Telephony' },
    { name: 'Call Detail Records (CDRs / Logs)', path: '/call-logs', category: 'Telephony' },
    { name: 'Digital Receptionists (IVR)', path: '/digital-receptionists', category: 'Voice' },
    { name: 'Ring Groups', path: '/ring-groups', category: 'Voice' },
    { name: 'Call Queues & Wallboard', path: '/call-queues', category: 'Voice' },
    { name: 'Voicemail Server', path: '/voicemail', category: 'Voice' },
    { name: 'System Settings & Admins', path: '/system-settings', category: 'Settings' },
    { name: 'Reports & Analytics', path: '/reports-analytics', category: 'Analytics' },
    { name: 'Inbound Routing Rules', path: '/inbound-routes', category: 'Routing' },
    { name: 'Outbound Dial Routes', path: '/outbound-routes', category: 'Routing' },
    { name: 'SIP Trunks & Gateways', path: '/trunks', category: 'Carrier' },
    { name: 'Auto Provisioning Templates', path: '/provisioning', category: 'Hardware' },
    { name: 'Video Conference Rooms', path: '/video-rooms', category: 'Media' },
    { name: 'IP Camera RTSP Feeds', path: '/cameras', category: 'Media' },
    { name: 'Team Chat Channels', path: '/chat', category: 'Messaging' },
    { name: 'SMS & MMS Trunks', path: '/sms', category: 'Messaging' },
    { name: 'Firewall & IP ACLs', path: '/firewall', category: 'Security' },
    { name: 'Anti-Fraud Protection', path: '/fraud', category: 'Security' },
    { name: 'API Key Management', path: '/api-credentials', category: 'Developers' },
    { name: 'Webhooks Dispatcher', path: '/webhooks', category: 'Developers' },
    { name: 'Rate Throttling Policies', path: '/throttling', category: 'Security' },
    { name: 'Multi-Tenant Domains', path: '/domains', category: 'Tenancy' },
    { name: 'Git Version History', path: '/versioning', category: 'System' }
  ];

  $: filteredSearch = searchQuery.trim() 
    ? searchablePages.filter(p => p.name.toLowerCase().includes(searchQuery.toLowerCase()) || p.category.toLowerCase().includes(searchQuery.toLowerCase()))
    : [];

  function handleNavigate(path: string) {
    window.history.pushState({}, '', path);
    currentPath = path;
    searchQuery = '';
    showSearchResults = false;
    window.dispatchEvent(new CustomEvent('app-navigate', { detail: { path } }));
  }

  function handleLogout() {
    localStorage.removeItem('pbx_token');
    localStorage.removeItem('pbx_user');
    window.location.href = '/';
  }

  onMount(() => {
    const handlePopState = () => {
      currentPath = window.location.pathname;
    };
    window.addEventListener('popstate', handlePopState);
    
    window.addEventListener('app-navigate', (e: any) => {
      currentPath = e.detail.path;
    });

    return () => {
      window.removeEventListener('popstate', handlePopState);
      window.removeEventListener('app-navigate', () => {});
    };
  });

  const routes: Record<string, any> = {
    '/': Dashboard,
    '/dashboard': Dashboard,
    '/cdrs': CDRs,
    '/cdr': CDRs,
    '/call-logs': CDRs,
    '/inbound-routes': Routing,
    '/inbound': Routing,
    '/routing': Routing,
    '/outbound-routes': OutboundRoutes,
    '/outbound': OutboundRoutes,
    '/trunks': Trunks,
    '/trunk': Trunks,
    '/extensions': Extensions,
    '/extension': Extensions,
    '/ivr': IVRs,
    '/ivrs': IVRs,
    '/digital-receptionists': IVRs,
    '/ivr-builder': IVR,
    '/ring-groups': Extensions,
    '/call-queues': Wallboard,
    '/voicemail': Voicemail,
    '/voicemails': Voicemail,
    '/system-settings': Admins,
    '/reports': CDRs,
    '/reports-analytics': CDRs,
    '/video-rooms': VideoRooms,
    '/video': VideoRooms,
    '/cameras': Camera,
    '/camera': Camera,
    '/sms': SMSTrunks,
    '/chat': TeamChat,
    '/api-credentials': APICredentials,
    '/credentials': APICredentials,
    '/webhooks': Webhooks,
    '/throttling': Throttling,
    '/firewall': Firewall,
    '/fraud': AntiFraud,
    '/anti-fraud': AntiFraud,
    '/versioning': VersionControl,
    '/version-control': VersionControl,
    '/commits': VersionControl,
    '/admins': Admins,
    '/admin': Admins,
    '/users': Admins,
    '/wallboard': Wallboard,
    '/domains': Domains,
    '/domain': Domains,
    '/tenants': Domains,
    '/provisioning': Provisioning
  };

  $: normalizedPath = (currentPath.endsWith('/') && currentPath.length > 1) ? currentPath.slice(0, -1) : currentPath;
  $: ActiveComponent = currentPath === '/wallboard-fullscreen' ? Wallboard : (routes[normalizedPath] || routes[currentPath]);
  $: isDashboard = currentPath === '/' || currentPath === '/dashboard';
</script>

<div class="min-h-screen bg-[#EEF2F6] flex flex-col font-sans select-none">
  {#if !token}
    <Login />
  {:else if currentPath === '/wallboard-fullscreen'}
    <div class="h-screen w-screen bg-[#06080F] overflow-hidden">
      <Wallboard />
    </div>
  {:else}
    <!-- Top Bar: Deep Navy #1B2A4A, Spans Full Width -->
    <header class="h-[52px] bg-[#1B2A4A] border-b border-[#24355A] flex items-center justify-between px-4 z-30 flex-shrink-0">
      <!-- Left: Logo & Brand Name -->
      <button type="button" class="flex items-center space-x-2.5 w-60 flex-shrink-0 cursor-pointer text-left focus:outline-none" on:click={() => handleNavigate('/')}>
        <div class="w-7 h-7 rounded-full bg-cyan-600/30 border border-cyan-400/40 flex items-center justify-center text-cyan-400">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path>
          </svg>
        </div>
        <span class="text-white font-semibold text-sm tracking-tight">CommuniGate PBX</span>
      </button>

      <!-- Center: Global Search Input -->
      <div class="flex-1 max-w-md mx-4 relative hidden sm:block">
        <div class="relative">
          <input 
            type="text" 
            bind:value={searchQuery}
            on:focus={() => showSearchResults = true}
            on:blur={() => setTimeout(() => showSearchResults = false, 200)}
            placeholder="Global search" 
            class="w-full bg-[#273B58] border border-[#384C6B] text-slate-100 placeholder-slate-400 text-xs rounded-md pl-3 pr-9 py-1.5 focus:outline-none focus:border-cyan-400 focus:bg-[#2F4466] transition-colors"
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-slate-400">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>

        <!-- Global Search Autocomplete Popup -->
        {#if showSearchResults && filteredSearch.length > 0}
          <div class="absolute top-full left-0 right-0 mt-1.5 bg-white border border-slate-200 rounded-lg shadow-xl z-50 overflow-hidden divide-y divide-slate-100 max-h-80 overflow-y-auto">
            {#each filteredSearch as item}
              <button 
                type="button"
                on:click={() => handleNavigate(item.path)}
                class="w-full px-3.5 py-2.5 text-left text-xs flex items-center justify-between hover:bg-slate-50 transition-colors group"
              >
                <div class="flex items-center space-x-2">
                  <span class="w-1.5 h-1.5 rounded-full bg-blue-600"></span>
                  <span class="font-semibold text-slate-800 group-hover:text-blue-600">{item.name}</span>
                </div>
                <span class="text-[10px] uppercase font-bold text-slate-400 tracking-wider">{item.category}</span>
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Right: User Avatar, Name, Dot, Gear, Logout -->
      <div class="flex items-center space-x-3">
        <!-- User Profile Pill -->
        <div class="flex items-center space-x-2 text-xs">
          <div class="w-6 h-6 rounded-full bg-[#2563EB] text-white flex items-center justify-center text-[11px] font-bold shadow-sm">
            {user.email ? user.email[0].toUpperCase() : 'D'}
          </div>
          <span class="text-xs font-semibold text-slate-100 hidden md:inline">
            {user.email ? (user.email.split('@')[0].charAt(0).toUpperCase() + user.email.split('@')[0].slice(1)) : 'David K.'}
          </span>
          <span class="w-2 h-2 rounded-full bg-emerald-400 flex-shrink-0" title="Online"></span>
          <svg class="w-3.5 h-3.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </div>

        <!-- Settings Gear Icon -->
        <button 
          on:click={() => handleNavigate('/system-settings')}
          class="p-1.5 text-slate-300 hover:text-white hover:bg-[#253655] rounded-md transition-colors" 
          title="System Settings"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </button>

        <!-- Sign Out Button -->
        <button 
          on:click={handleLogout} 
          class="text-[11px] font-semibold text-slate-400 hover:text-rose-300 px-2 py-1 hover:bg-[#253655] rounded-md transition-colors"
        >
          Logout
        </button>
      </div>
    </header>

    <!-- Main Workspace Layout -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Sidebar (White, matching reference) -->
      <Sidebar />

      <!-- Main Canvas -->
      <main class="flex-1 overflow-y-auto {isDashboard ? 'bg-[#EEF2F6] text-slate-800' : 'bg-[#0E1320] text-slate-100'} p-5 md:p-6 transition-colors">
        <div class="max-w-[1400px] mx-auto pb-20">
          {#if ActiveComponent}
            <svelte:component this={ActiveComponent} />
          {:else}
            <div class="p-12 text-center bg-white border border-slate-200 rounded-xl max-w-md mx-auto my-12 shadow-sm">
              <h2 class="text-base font-bold text-slate-800 mb-1">Route Not Found</h2>
              <p class="text-slate-500 text-xs mb-4">The route "{currentPath}" does not have an active view registered.</p>
              <button on:click={() => handleNavigate('/')} class="px-4 py-2 bg-[#2563EB] text-white text-xs font-semibold rounded-lg shadow-sm">Return to Dashboard</button>
            </div>
          {/if}
        </div>
      </main>

      <!-- WebRTC Softphone Drawer -->
      <Softphone />
    </div>
  {/if}
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background-color: #EEF2F6;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }
</style>
