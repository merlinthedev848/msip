<script lang="ts">
  import Sidebar from './components/Sidebar.svelte';
  import { onMount } from 'svelte';
  
  import Dashboard from './views/Dashboard.svelte';
  import CDRs from './views/CDRs.svelte';
  import IVRs from './views/IVRs.svelte';
  import Routing from './views/Routing.svelte';
  import OutboundRoutes from './views/OutboundRoutes.svelte';
  import Trunks from './views/Trunks.svelte';
  import Extensions from './views/Extensions.svelte';
  import IVR from './views/IVR.svelte';
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
    '/cdrs': CDRs,
    '/inbound-routes': Routing,
    '/outbound-routes': OutboundRoutes,
    '/trunks': Trunks,
    '/extensions': Extensions,
    '/ivr': IVRs,
    '/voicemail': Voicemail,
    '/video-rooms': VideoRooms,
    '/cameras': Camera,
    '/sms': SMSTrunks,
    '/chat': TeamChat,
    '/api-credentials': APICredentials,
    '/webhooks': Webhooks,
    '/throttling': Throttling,
    '/firewall': Firewall,
    '/fraud': AntiFraud,
    '/versioning': VersionControl,
    '/admins': Admins,
    '/wallboard': Wallboard,
    '/domains': Domains,
    '/provisioning': Provisioning
  };

  $: ActiveComponent = currentPath === '/wallboard-fullscreen' ? Wallboard : routes[currentPath];
</script>

<div>
  {#if !token}
    <Login />
  {:else if currentPath === '/wallboard-fullscreen'}
    <div class="h-screen w-screen bg-gray-950 overflow-hidden">
      <Wallboard />
    </div>
  {:else}
    <main class="flex min-h-screen bg-slate-50 font-sans selection:bg-blue-500 selection:text-white text-slate-900">
      <Sidebar />
      
      <div class="flex-1 flex flex-col h-screen overflow-hidden">
        <!-- Header -->
        <header class="h-16 border-b border-slate-200 bg-white flex items-center justify-between px-8 sticky top-0 z-20 shadow-sm">
          <div class="flex items-center">
            <h2 class="text-sm font-semibold text-slate-500 uppercase tracking-widest">Workspace</h2>
          </div>
          <div class="flex items-center space-x-4">
            <button on:click={handleLogout} class="bg-white hover:bg-slate-50 hover:text-slate-900 text-sm py-1.5 px-4 rounded-full transition-colors border border-slate-200 text-slate-600 shadow-sm flex items-center space-x-2">
              <div class="w-5 h-5 rounded-full bg-blue-600 flex items-center justify-center text-xs text-white font-bold">
                {user.email ? user.email[0].toUpperCase() : 'U'}
              </div>
              <span>{user.email || 'Admin'} (Logout)</span>
            </button>
          </div>
        </header>

        <!-- Main Content Area -->
        <div class="flex-1 overflow-y-auto p-8 relative">
          
          <div class="w-full pb-20">
            {#if ActiveComponent}
              <svelte:component this={ActiveComponent} />
            {:else}
              <div class="p-12 text-center flex flex-col items-center justify-center h-full">
                <div class="w-16 h-16 bg-gray-800 rounded-2xl flex items-center justify-center mb-4">
                  <svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
                </div>
                <h2 class="text-xl font-semibold text-slate-900 mb-2">404 Not Found</h2>
                <p class="text-slate-500 max-w-md">The route "{currentPath}" does not exist.</p>
              </div>
            {/if}
          </div>
        </div>
      </div>
      
    </main>
  {/if}
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background-color: #f8fafc;
  }
</style>
