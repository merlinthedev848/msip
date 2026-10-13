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

  const routes = {
    '/': Dashboard,
    '/dashboard': Dashboard,
    '/cdrs': CDRs,
    '/cdr': CDRs,
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
    '/ivr-builder': IVR,
    '/voicemail': Voicemail,
    '/voicemails': Voicemail,
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
  $: pageTitle = currentPath === '/' ? 'System Overview' : normalizedPath.replace('/', '').replace(/-/g, ' ').toUpperCase();
</script>

<div class="min-h-screen bg-[#0B0F19] text-slate-100">
  {#if !token}
    <Login />
  {:else if currentPath === '/wallboard-fullscreen'}
    <div class="h-screen w-screen bg-[#06080F] overflow-hidden">
      <Wallboard />
    </div>
  {:else}
    <main class="flex min-h-screen bg-[#0B0F19] font-sans selection:bg-indigo-500 selection:text-white text-slate-100">
      <Sidebar />
      
      <div class="flex-1 flex flex-col h-screen overflow-hidden bg-gradient-to-br from-[#0B0F19] via-[#0E1424] to-[#0A0D18]">
        <!-- Modern Executive Header -->
        <header class="h-16 border-b border-slate-800/80 bg-[#0D1322]/80 backdrop-blur-xl flex items-center justify-between px-8 sticky top-0 z-20 shadow-lg shadow-black/20">
          <div class="flex items-center space-x-3">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-widest">WORKSPACE</span>
            <span class="text-slate-600">/</span>
            <span class="text-xs font-mono font-semibold text-indigo-400">{pageTitle}</span>
          </div>

          <div class="flex items-center space-x-4">
            <!-- Telecom Cluster Status -->
            <div class="hidden sm:flex items-center space-x-2 px-3 py-1 rounded-full bg-slate-900/80 border border-slate-800 text-xs">
              <span class="relative flex h-2 w-2">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
              </span>
              <span class="text-slate-300 font-medium">FreeSWITCH Connected</span>
            </div>

            <!-- User Logout Pill -->
            <button 
              on:click={handleLogout} 
              class="bg-slate-900/90 hover:bg-slate-800 text-sm py-1.5 px-4 rounded-full transition-all border border-slate-800 hover:border-slate-700 text-slate-300 hover:text-white shadow-md flex items-center space-x-2.5 group cursor-pointer"
            >
              <div class="w-6 h-6 rounded-full bg-gradient-to-tr from-indigo-600 to-purple-600 flex items-center justify-center text-xs text-white font-bold shadow-sm shadow-indigo-500/30">
                {user.email ? user.email[0].toUpperCase() : 'A'}
              </div>
              <span class="font-medium text-xs text-slate-300 group-hover:text-white">{user.email || 'Admin'}</span>
              <span class="text-[11px] text-slate-500 group-hover:text-rose-400 transition-colors">(Logout)</span>
            </button>
          </div>
        </header>

        <!-- Main Content Area -->
        <div class="flex-1 overflow-y-auto p-6 md:p-8 relative">
          <!-- Subtle ambient background glowing orbs -->
          <div class="absolute top-0 right-1/4 w-96 h-96 bg-indigo-600/5 rounded-full filter blur-3xl pointer-events-none"></div>
          <div class="absolute bottom-1/3 left-10 w-96 h-96 bg-purple-600/5 rounded-full filter blur-3xl pointer-events-none"></div>

          <div class="w-full pb-20 relative z-10">
            {#if ActiveComponent}
              <svelte:component this={ActiveComponent} />
            {:else}
              <div class="p-12 text-center flex flex-col items-center justify-center h-full">
                <div class="w-16 h-16 bg-slate-900 border border-slate-800 rounded-2xl flex items-center justify-center mb-4 text-indigo-400 shadow-xl">
                  <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
                </div>
                <h2 class="text-xl font-bold text-white mb-2">404 View Not Found</h2>
                <p class="text-slate-400 max-w-md text-sm">The route "{currentPath}" is not registered.</p>
              </div>
            {/if}
          </div>
        </div>

        <!-- Embedded WebRTC Softphone Bar -->
        <Softphone />
      </div>
    </main>
  {/if}
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background-color: #0B0F19 !important;
    color: #f1f5f9;
  }
</style>
