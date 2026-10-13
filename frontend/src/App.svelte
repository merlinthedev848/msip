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
    <main class="flex min-h-screen bg-[#080A0F] font-sans selection:bg-[#2563EB] selection:text-white text-slate-100">
      <Sidebar />
      
      <div class="flex-1 flex flex-col h-screen overflow-hidden bg-[#0A0D14]">
        <!-- Solid Tactile Header -->
        <header class="h-14 border-b border-[#1E2536] bg-[#0E121B] flex items-center justify-between px-6 sticky top-0 z-20 shadow-md">
          <div class="flex items-center space-x-3">
            <span class="text-[11px] font-bold font-mono text-slate-400 uppercase tracking-widest">CONSOLE</span>
            <span class="text-slate-600 font-mono">/</span>
            <span class="text-xs font-mono font-bold text-[#38BDF8]">{pageTitle}</span>
          </div>

          <div class="flex items-center space-x-3">
            <!-- Telecom Cluster Status -->
            <div class="hidden sm:flex items-center space-x-2 px-2.5 py-1 rounded bg-[#0D281E] border border-[#10B981]/40 text-xs font-mono font-bold text-[#10B981]">
              <span class="w-2 h-2 rounded-full bg-[#10B981]"></span>
              <span>FreeSWITCH ESL: CONNECTED</span>
            </div>

            <!-- User Logout Button -->
            <button 
              on:click={handleLogout} 
              class="bg-[#141A26] hover:bg-[#1E2538] text-xs font-mono font-semibold py-1.5 px-3 rounded-lg transition-colors border border-[#242D40] text-slate-300 hover:text-white flex items-center space-x-2 cursor-pointer"
            >
              <span class="w-5 h-5 rounded bg-[#2563EB] flex items-center justify-center text-[10px] text-white font-bold">
                {user.email ? user.email[0].toUpperCase() : 'A'}
              </span>
              <span>{user.email || 'Admin'}</span>
              <span class="text-slate-500 hover:text-rose-400">[Logout]</span>
            </button>
          </div>
        </header>

        <!-- Main Content Area -->
        <div class="flex-1 overflow-y-auto p-5 md:p-6 relative">
          <div class="w-full pb-20">
            {#if ActiveComponent}
              <svelte:component this={ActiveComponent} />
            {:else}
              <div class="p-12 text-center flex flex-col items-center justify-center h-full">
                <div class="w-14 h-14 bg-[#121622] border border-[#202738] rounded-xl flex items-center justify-center mb-3 text-[#38BDF8]">
                  <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
                </div>
                <h2 class="text-lg font-bold text-white mb-1 font-mono">404 ROUTE NOT REGISTERED</h2>
                <p class="text-slate-400 max-w-md text-xs font-mono">Path "{currentPath}" is not found in the router table.</p>
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
