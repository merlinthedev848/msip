<script lang="ts">
  import Sidebar from './components/Sidebar.svelte';
  import { onMount } from 'svelte';
  
  import Dashboard from './views/Dashboard.svelte';
  import CDRs from './views/CDRs.svelte';
  import Routing from './views/Routing.svelte';
  import OutboundRoutes from './views/OutboundRoutes.svelte';
  import Trunks from './views/Trunks.svelte';
  import Extensions from './views/Extensions.svelte';
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
  
  import Wallboard from './views/Wallboard.svelte';
  import Domains from './views/Domains.svelte';
  import Provisioning from './views/Provisioning.svelte';
  import RingGroups from './views/RingGroups.svelte';
  import ReportsAnalytics from './views/ReportsAnalytics.svelte';

  import TimeConditions from './views/TimeConditions.svelte';
  import OutboundRestrictions from './views/OutboundRestrictions.svelte';
  import PagingPickup from './views/PagingPickup.svelte';
  import AudioPrompts from './views/AudioPrompts.svelte';
  import FeatureCodes from './views/FeatureCodes.svelte';
  import Diagnostics from './views/Diagnostics.svelte';

  import ButtonProfiles from './views/ButtonProfiles.svelte';
  import Contacts from './views/Contacts.svelte';
  import TelephoneNumbers from './views/TelephoneNumbers.svelte';
  import Sites from './views/Sites.svelte';
  import Diverts from './views/Diverts.svelte';
  import Devices from './views/Devices.svelte';
  import SIPSettings from './views/SIPSettings.svelte';
  import AIAssistants from './views/AIAssistants.svelte';
  import AIBookings from './views/AIBookings.svelte';
  
  import Login from './views/Login.svelte';
  import Softphone from './components/Softphone.svelte';
  import Toaster from './components/ui/Toaster.svelte';

  let currentPath = window.location.pathname;
  let token = localStorage.getItem('pbx_token');
  let user = JSON.parse(localStorage.getItem('pbx_user') || '{}');

  let searchQuery = '';
  let showSearchResults = false;
  let sidebarCollapsed = false;

  const searchablePages = [
    { name: 'Dashboard & Metrics', path: '/', category: 'Overview' },
    { name: 'AI Virtual Receptionists & Voice AI', path: '/ai-assistants', category: 'Voice' },
    { name: 'AI Appointment Bookings & Calendar', path: '/ai-bookings', category: 'Voice' },
    { name: 'Extensions (SIP Endpoints)', path: '/extensions', category: 'Telephony' },
    { name: 'Call Detail Records (CDRs / Logs)', path: '/call-logs', category: 'Telephony' },
    { name: 'Digital Receptionists (IVR)', path: '/digital-receptionists', category: 'Voice' },
    { name: 'Ring Groups', path: '/ring-groups', category: 'Voice' },
    { name: 'Call Queues & Wallboard', path: '/call-queues', category: 'Voice' },
    { name: 'Voicemail Server', path: '/voicemail', category: 'Voice' },
    { name: 'Time Conditions & Holidays', path: '/time-conditions', category: 'Routing' },
    { name: 'Outbound Restrictions & AutoCLIP', path: '/outbound-restrictions', category: 'Routing' },
    { name: 'Paging, Intercom & Call Pickup', path: '/paging-pickup', category: 'Voice' },
    { name: 'Voice Prompts & Call Recording', path: '/audio-prompts', category: 'Media' },
    { name: 'PBX Feature Codes (*97, *8, *0)', path: '/feature-codes', category: 'Telephony' },
    { name: 'SIP, IAX & Jitter Diagnostics', path: '/diagnostics', category: 'System' },
    { name: 'System Settings & Admins', path: '/admins', category: 'Settings' },
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
    { name: 'Button Profiles (BLF / DSS Keys)', path: '/button-profiles', category: 'Telephony' },
    { name: 'Contacts & Directory', path: '/contacts', category: 'Telephony' },
    { name: 'Telephone Numbers (DIDs)', path: '/telephone-numbers', category: 'Carrier' },
    { name: 'Sites & Branch Offices', path: '/sites', category: 'Tenancy' },
    { name: 'Call Diverts & Forwarding', path: '/diverts', category: 'Routing' },
    { name: 'Devices & Hardware Handsets', path: '/devices', category: 'Hardware' },
    { name: 'Multi-Tenant Domains', path: '/domains', category: 'Tenancy' },
    { name: 'Git Version History', path: '/versioning', category: 'System' },
    { name: 'Advanced SIP & RTP Settings', path: '/sip-settings', category: 'System' }
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
    '/menus': IVRs,
    '/ring-groups': RingGroups,
    '/call-distribution-groups': RingGroups,
    '/call-queues': Wallboard,
    '/queues': Wallboard,
    '/voicemail': Voicemail,
    '/voicemails': Voicemail,
    '/system-settings': Admins,
    '/reports': ReportsAnalytics,
    '/reports-analytics': ReportsAnalytics,
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
    '/provisioning': Provisioning,
    '/time-conditions': TimeConditions,
    '/holidays': TimeConditions,
    '/outbound-restrictions': OutboundRestrictions,
    '/autoclip': OutboundRestrictions,
    '/paging-pickup': PagingPickup,
    '/paging': PagingPickup,
    '/intercom': PagingPickup,
    '/pickup': PagingPickup,
    '/audio-prompts': AudioPrompts,
    '/moh': AudioPrompts,
    '/recordings': AudioPrompts,
    '/voice-prompts': AudioPrompts,
    '/feature-codes': FeatureCodes,
    '/codes': FeatureCodes,
    '/diagnostics': Diagnostics,
    '/jitter': Diagnostics,
    '/button-profiles': ButtonProfiles,
    '/buttons': ButtonProfiles,
    '/contacts': Contacts,
    '/directory': Contacts,
    '/telephone-numbers': TelephoneNumbers,
    '/numbers': TelephoneNumbers,
    '/dids': TelephoneNumbers,
    '/sites': Sites,
    '/locations': Sites,
    '/diverts': Diverts,
    '/forwarding': Diverts,
    '/devices': Devices,
    '/handsets': Devices,
    '/sip-settings': SIPSettings,
    '/sip': SIPSettings,
    '/rtp': SIPSettings,
    '/pbx-settings': SIPSettings,
    '/ai-assistants': AIAssistants,
    '/ai-assistant': AIAssistants,
    '/ai-receptionist': AIAssistants,
    '/ai': AIAssistants,
    '/ai-bookings': AIBookings,
    '/bookings': AIBookings,
    '/appointments': AIBookings
  };

  $: normalizedPath = (currentPath.endsWith('/') && currentPath.length > 1) ? currentPath.slice(0, -1) : currentPath;
  $: ActiveComponent = currentPath === '/wallboard-fullscreen' ? Wallboard : (routes[normalizedPath] || routes[currentPath]);
</script>

<Toaster />

<div class="min-h-screen bg-[#EEF2F6] flex flex-col font-sans select-none">
  {#if !token}
    <Login />
  {:else if currentPath === '/wallboard-fullscreen'}
    <div class="h-screen w-screen bg-[#EEF2F6] overflow-hidden">
      <Wallboard />
    </div>
  {:else}
    <!-- TOP NAVIGATION BAR -->
    <header class="h-[50px] bg-white border-b border-slate-200 flex items-center justify-between px-4 z-30 flex-shrink-0 shadow-[0_1px_2px_rgba(0,0,0,0.03)]">
      <!-- Left: Hamburger / Brand Label -->
      <div class="flex items-center space-x-3">
        <button 
          type="button" 
          on:click={() => sidebarCollapsed = !sidebarCollapsed} 
          class="p-1.5 text-slate-500 hover:text-slate-800 hover:bg-slate-100 rounded transition-colors" 
          title="Toggle Navigation"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <span class="text-xs font-bold text-slate-700 uppercase tracking-wider hidden sm:inline">MSIP Telephony Control</span>
      </div>

      <!-- Center: Quick Global Search -->
      <div class="flex-1 max-w-sm mx-4 relative hidden md:block">
        <input 
          type="text" 
          bind:value={searchQuery}
          on:focus={() => showSearchResults = true}
          on:blur={() => setTimeout(() => showSearchResults = false, 200)}
          placeholder="Search extensions, routes, settings..." 
          class="w-full bg-slate-50 border border-slate-200 text-slate-800 placeholder-slate-400 text-xs rounded-md pl-3 pr-8 py-1.5 focus:outline-none focus:border-sky-500 focus:bg-white transition-colors"
        />
        <div class="absolute inset-y-0 right-0 pr-2.5 flex items-center pointer-events-none text-slate-400">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>

        {#if showSearchResults && filteredSearch.length > 0}
          <div class="absolute top-full left-0 right-0 mt-1 bg-white border border-slate-200 rounded-lg shadow-xl z-50 overflow-hidden divide-y divide-slate-100 max-h-72 overflow-y-auto">
            {#each filteredSearch as item}
              <button 
                type="button"
                on:click={() => handleNavigate(item.path)}
                class="w-full px-3 py-2 text-left text-xs flex items-center justify-between hover:bg-slate-50 transition-colors"
              >
                <span class="font-semibold text-slate-700">{item.name}</span>
                <span class="text-[10px] uppercase font-bold text-slate-400">{item.category}</span>
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Right: Utility Icons & User Menu -->
      <div class="flex items-center space-x-3 text-slate-600">
        <button on:click={() => handleNavigate('/call-queues')} class="p-1.5 hover:text-slate-900 hover:bg-slate-100 rounded text-slate-500" title="Live Wallboard">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>
        </button>

        <button on:click={() => handleNavigate('/versioning')} class="p-1.5 hover:text-slate-900 hover:bg-slate-100 rounded text-slate-500" title="Change History">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="1.8"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M12 7v5l3 3"/></svg>
        </button>

        <div class="h-4 w-px bg-slate-200"></div>

        <!-- User Profile -->
        <div class="flex items-center space-x-1.5 cursor-pointer">
          <div class="w-6 h-6 rounded-full bg-[#0284C7] text-white flex items-center justify-center text-[11px] font-bold">
            <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"/></svg>
          </div>
          <span class="text-xs font-semibold text-slate-700">{user.email ? user.email.split('@')[0] : 'admin'}</span>
        </div>

        <button 
          type="button" 
          on:click={handleLogout} 
          class="text-[11px] font-medium text-slate-400 hover:text-rose-600 px-1.5 py-0.5 rounded transition-colors"
        >
          Logout
        </button>
      </div>
    </header>

    <!-- Main Workspace Layout -->
    <div class="flex-1 flex overflow-hidden">
      <Sidebar collapsed={sidebarCollapsed} />

      <!-- Main Canvas (Always permanent clean light background) -->
      <main class="flex-1 overflow-y-auto bg-[#EEF2F6] text-slate-800 p-5 md:p-6 transition-colors">
        <div class="max-w-[1440px] mx-auto pb-20">
          {#if ActiveComponent}
            <svelte:component this={ActiveComponent} />
          {:else}
            <div class="p-12 text-center bg-white border border-slate-200 rounded-xl max-w-md mx-auto my-12 shadow-sm">
              <h2 class="text-base font-bold text-slate-800 mb-1">Route Not Found</h2>
              <p class="text-slate-500 text-xs mb-4">The route "{currentPath}" does not have an active view registered.</p>
              <button on:click={() => handleNavigate('/')} class="px-4 py-2 bg-[#0284C7] text-white text-xs font-semibold rounded-lg shadow-sm">Return to Dashboard</button>
            </div>
          {/if}
        </div>
      </main>

      <!-- WebRTC Softphone Drawer -->
      <Softphone />
    </div>
  {/if}
</div>
