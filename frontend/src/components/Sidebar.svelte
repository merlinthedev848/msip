<script lang="ts">
  import { onMount } from 'svelte';
  import { layoutStyle, setLayoutStyle, type LayoutStyle } from '../lib/theme';
  
  let currentPath = window.location.pathname;
  let currentLayout: LayoutStyle = 'yeastar';

  layoutStyle.subscribe(val => {
    currentLayout = val;
  });

  const handleNav = (path: string) => {
    window.history.pushState({}, '', path);
    currentPath = path;
    window.dispatchEvent(new CustomEvent('app-navigate', { detail: { path } }));
  };

  // -------------------------------------------------------------
  // COMMUNIGATE PBX NAVIGATION DEFINITION
  // -------------------------------------------------------------
  const communigatePrimaryNav = [
    { name: 'Dashboard', path: '/', alias: ['/', '/dashboard'], icon: 'dashboard' },
    { name: 'Extensions', path: '/extensions', alias: ['/extensions', '/extension'], icon: 'phone' },
    { name: 'Call Logs', path: '/call-logs', alias: ['/call-logs', '/cdrs', '/cdr'], icon: 'logs' },
    { name: 'Digital Receptionists', path: '/digital-receptionists', alias: ['/digital-receptionists', '/ivr', '/ivrs', '/ivr-builder'], icon: 'ivr' },
    { name: 'Ring Groups', path: '/ring-groups', alias: ['/ring-groups'], icon: 'users' },
    { name: 'Call Queues', path: '/call-queues', alias: ['/call-queues', '/wallboard'], icon: 'queue' },
    { name: 'Voicemail', path: '/voicemail', alias: ['/voicemail', '/voicemails'], icon: 'voicemail' },
    { name: 'System Settings', path: '/system-settings', alias: ['/system-settings', '/admins', '/admin'], icon: 'settings' },
    { name: 'Reports & Analytics', path: '/reports-analytics', alias: ['/reports-analytics', '/reports'], icon: 'reports' }
  ];

  let showCommunigateAdvanced = false;

  const communigateAdvancedSections = [
    {
      title: 'Call Routing & Carrier',
      items: [
        { name: 'Inbound Routes', path: '/inbound-routes' },
        { name: 'Outbound Routes', path: '/outbound-routes' },
        { name: 'SIP Trunks', path: '/trunks' },
        { name: 'Auto Provisioning', path: '/provisioning' }
      ]
    },
    {
      title: 'Unified Communications',
      items: [
        { name: 'Video Rooms', path: '/video-rooms' },
        { name: 'IP Cameras', path: '/cameras' },
        { name: 'Team Chat', path: '/chat' },
        { name: 'SMS & Messaging', path: '/sms' }
      ]
    },
    {
      title: 'Security & Infrastructure',
      items: [
        { name: 'Multi-Tenant Domains', path: '/domains' },
        { name: 'API Credentials', path: '/api-credentials' },
        { name: 'Event Webhooks', path: '/webhooks' },
        { name: 'Rate Throttling', path: '/throttling' },
        { name: 'Firewall & ACLs', path: '/firewall' },
        { name: 'Anti-Fraud Shield', path: '/fraud' },
        { name: 'Version Control', path: '/versioning' }
      ]
    }
  ];

  // -------------------------------------------------------------
  // YEASTAR PBX NAVIGATION DEFINITION (Matching Reference Image)
  // -------------------------------------------------------------
  interface YeastarNavItem {
    name: string;
    path?: string;
    icon: string;
    expanded?: boolean;
    children?: { name: string; path: string }[];
  }

  let yeastarNavItems: YeastarNavItem[] = [
    { name: 'Dashboard', path: '/', icon: 'dashboard' },
    { 
      name: 'Extension and Trunk', 
      icon: 'ext-trunk', 
      expanded: false,
      children: [
        { name: 'Extensions', path: '/extensions' },
        { name: 'SIP Trunks', path: '/trunks' },
        { name: 'Auto Provisioning', path: '/provisioning' }
      ]
    },
    { 
      name: 'Contacts', 
      icon: 'contacts', 
      expanded: false,
      children: [
        { name: 'Company Directory', path: '/extensions' },
        { name: 'Speed Dial List', path: '/extensions' }
      ]
    },
    { 
      name: 'Auto Provisioning', 
      path: '/provisioning',
      icon: 'provisioning'
    },
    { 
      name: 'Call Control', 
      icon: 'call-control', 
      expanded: false,
      children: [
        { name: 'Inbound Routes', path: '/inbound-routes' },
        { name: 'Outbound Routes', path: '/outbound-routes' },
        { name: 'Call Queues & Wallboard', path: '/call-queues' }
      ]
    },
    { 
      name: 'Call Features', 
      icon: 'call-features', 
      expanded: false,
      children: [
        { name: 'IVR / Auto Attendant', path: '/digital-receptionists' },
        { name: 'Ring Groups', path: '/ring-groups' },
        { name: 'Voicemail Server', path: '/voicemail' },
        { name: 'Video Rooms', path: '/video-rooms' }
      ]
    },
    { 
      name: 'PBX Settings', 
      icon: 'pbx-settings', 
      expanded: false,
      children: [
        { name: 'Tenant Domains', path: '/domains' },
        { name: 'Audio Prompts & MoH', path: '/digital-receptionists' }
      ]
    },
    { 
      name: 'System', 
      icon: 'system', 
      expanded: false,
      children: [
        { name: 'Cluster & Admins', path: '/system-settings' },
        { name: 'Version Control', path: '/versioning' },
        { name: 'IP Camera Feeds', path: '/cameras' }
      ]
    },
    { 
      name: 'Security', 
      icon: 'security', 
      expanded: false,
      children: [
        { name: 'Firewall & IP ACLs', path: '/firewall' },
        { name: 'Anti-Fraud Shield', path: '/fraud' },
        { name: 'Rate Throttling', path: '/throttling' },
        { name: 'API Key Management', path: '/api-credentials' }
      ]
    },
    { 
      name: 'Maintenance', 
      path: '/system-settings', 
      icon: 'maintenance'
    },
    { 
      name: 'Integration', 
      icon: 'integration', 
      expanded: false,
      children: [
        { name: 'Event Webhooks', path: '/webhooks' },
        { name: 'Team Chat', path: '/chat' },
        { name: 'SMS Messaging', path: '/sms' }
      ]
    },
    { 
      name: 'Reports and Recordings', 
      icon: 'reports', 
      expanded: false,
      children: [
        { name: 'Call Detail Records', path: '/call-logs' },
        { name: 'CDR Analytics', path: '/reports-analytics' }
      ]
    },
    { 
      name: 'Plan', 
      path: '/system-settings', 
      icon: 'plan'
    }
  ];

  function toggleYeastarItem(index: number) {
    if (yeastarNavItems[index].children) {
      yeastarNavItems[index].expanded = !yeastarNavItems[index].expanded;
    } else if (yeastarNavItems[index].path) {
      handleNav(yeastarNavItems[index].path!);
    }
  }

  function isCommunigateActive(item: { path: string; alias?: string[] }) {
    if (item.alias && item.alias.includes(currentPath)) return true;
    return currentPath === item.path;
  }

  function isYeastarActive(item: YeastarNavItem) {
    if (item.path === currentPath) return true;
    if (item.path === '/' && (currentPath === '/' || currentPath === '/dashboard')) return true;
    if (item.children && item.children.some(c => c.path === currentPath)) return true;
    return false;
  }

  onMount(() => {
    const updatePath = () => { currentPath = window.location.pathname; };
    window.addEventListener('popstate', updatePath);
    window.addEventListener('app-navigate', (e: any) => { currentPath = e.detail.path; });

    return () => {
      window.removeEventListener('popstate', updatePath);
      window.removeEventListener('app-navigate', () => {});
    };
  });
</script>

<!-- ================================================================= -->
<!-- YEASTAR PBX SIDEBAR                                               -->
<!-- ================================================================= -->
{#if currentLayout === 'yeastar'}
  <aside class="w-60 bg-white h-[calc(100vh-50px)] overflow-y-auto border-r border-slate-200 py-3 flex flex-col flex-shrink-0 select-none shadow-[1px_0_4px_rgba(0,0,0,0.02)]">
    <!-- Yeastar Brand Header -->
    <div class="px-4 pb-3 flex items-center space-x-2.5">
      <div class="w-7 h-7 flex items-center justify-center text-[#0284C7]">
        <!-- Yeastar Pinwheel / Cross Logo -->
        <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none">
          <path d="M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12" stroke="#0284C7" stroke-width="2.5" stroke-linecap="round"/>
          <path d="M12 6V12L16 14" stroke="#0EA5E9" stroke-width="2.5" stroke-linecap="round"/>
        </svg>
      </div>
      <span class="text-base font-bold text-slate-800 tracking-tight">Yeastar</span>
    </div>

    <!-- Navigation List -->
    <nav class="space-y-0.5 px-2">
      {#each yeastarNavItems as item, idx}
        {@const active = isYeastarActive(item)}
        <div>
          <button 
            type="button"
            on:click={() => toggleYeastarItem(idx)}
            class="w-full flex items-center justify-between px-3 py-2 rounded-md text-xs transition-colors {active ? 'text-[#0284C7] font-bold bg-sky-50/70' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50 font-medium'}"
          >
            <div class="flex items-center space-x-2.5 truncate">
              <!-- Icon -->
              <span class="w-4 h-4 flex-shrink-0 flex items-center justify-center {active ? 'text-[#0284C7]' : 'text-slate-500'}">
                {#if item.icon === 'dashboard'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" stroke-linecap="round" d="M12 7v5l3 3"/></svg>
                {:else if item.icon === 'ext-trunk'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>
                {:else if item.icon === 'contacts'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/></svg>
                {:else if item.icon === 'provisioning'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>
                {:else if item.icon === 'call-control'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/></svg>
                {:else if item.icon === 'call-features'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg>
                {:else if item.icon === 'pbx-settings'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/></svg>
                {:else if item.icon === 'system'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/></svg>
                {:else if item.icon === 'security'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/></svg>
                {:else if item.icon === 'maintenance'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                {:else if item.icon === 'integration'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z"/></svg>
                {:else if item.icon === 'reports'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>
                {:else if item.icon === 'plan'}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/></svg>
                {/if}
              </span>
              <span class="truncate">{item.name}</span>
            </div>

            <!-- Chevron if has children -->
            {#if item.children}
              <svg class="w-3.5 h-3.5 text-slate-400 transition-transform duration-150 {item.expanded ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            {/if}
          </button>

          <!-- Sub-items dropdown -->
          {#if item.children && item.expanded}
            <div class="ml-6 pl-2 border-l border-slate-200 space-y-0.5 py-1">
              {#each item.children as sub}
                <button 
                  type="button"
                  on:click={() => handleNav(sub.path)}
                  class="w-full text-left px-2 py-1 text-xs rounded transition-colors {currentPath === sub.path ? 'text-[#0284C7] font-semibold bg-sky-50' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'}"
                >
                  {sub.name}
                </button>
              {/each}
            </div>
          {/if}
        </div>
      {/each}
    </nav>
  </aside>

<!-- ================================================================= -->
<!-- COMMUNIGATE PBX SIDEBAR                                           -->
<!-- ================================================================= -->
{:else}
  <aside class="w-64 bg-white h-[calc(100vh-52px)] overflow-y-auto border-r border-slate-200 p-3 flex flex-col flex-shrink-0 select-none shadow-[1px_0_4px_rgba(0,0,0,0.02)]">
    <!-- Primary Navigation Links -->
    <nav class="space-y-1">
      {#each communigatePrimaryNav as item}
        {@const active = isCommunigateActive(item)}
        <button 
          type="button"
          on:click={() => handleNav(item.path)}
          class="w-full flex items-center space-x-3 px-3.5 py-2.5 rounded-lg text-xs font-semibold transition-all duration-150 text-left {active ? 'bg-[#2563EB] text-white shadow-sm' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80'}"
        >
          <span class="w-4 h-4 flex-shrink-0 flex items-center justify-center {active ? 'text-white' : 'text-slate-500'}">
            {#if item.icon === 'dashboard'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>
            {:else if item.icon === 'phone'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
            {:else if item.icon === 'logs'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            {:else if item.icon === 'ivr'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 11.5V14m0-2.5v-6a1.5 1.5 0 113 0m-3 6a1.5 1.5 0 00-3 0v2a7.5 7.5 0 0015 0v-5a1.5 1.5 0 00-3 0m-6-3V11m0-5.5v-1a1.5 1.5 0 013 0v1m0 0V11m0-5.5a1.5 1.5 0 013 0v3m0 0V11" /></svg>
            {:else if item.icon === 'users'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
            {:else if item.icon === 'queue'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" /></svg>
            {:else if item.icon === 'voicemail'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="5.5" cy="11.5" r="3.5" stroke="currentColor" stroke-width="2" fill="none" /><circle cx="18.5" cy="11.5" r="3.5" stroke="currentColor" stroke-width="2" fill="none" /><line x1="5.5" y1="15" x2="18.5" y2="15" stroke="currentColor" stroke-width="2" /></svg>
            {:else if item.icon === 'settings'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            {:else if item.icon === 'reports'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>
            {/if}
          </span>
          <span class="truncate">{item.name}</span>
        </button>
      {/each}
    </nav>

    <div class="my-4 border-t border-slate-200"></div>

    <!-- Advanced Telecom & Security Modules Toggle -->
    <div class="mb-2">
      <button 
        type="button" 
        on:click={() => showCommunigateAdvanced = !showCommunigateAdvanced} 
        class="w-full flex items-center justify-between px-3 py-2 text-[11px] font-bold text-slate-500 uppercase tracking-wider hover:text-slate-800 transition-colors"
      >
        <span>Advanced Modules</span>
        <svg class="w-3.5 h-3.5 transition-transform duration-200 {showCommunigateAdvanced ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {#if showCommunigateAdvanced}
        <div class="space-y-4 pt-1 pl-1 pr-1 animate-in fade-in duration-150">
          {#each communigateAdvancedSections as sec}
            <div>
              <div class="px-2 py-1 text-[10px] font-bold uppercase tracking-wider text-slate-400">{sec.title}</div>
              <div class="space-y-0.5">
                {#each sec.items as subItem}
                  {@const isSubActive = currentPath === subItem.path}
                  <button 
                    type="button"
                    on:click={() => handleNav(subItem.path)}
                    class="w-full flex items-center px-3 py-1.5 rounded-md text-xs font-medium text-left transition-colors {isSubActive ? 'bg-blue-50 text-[#2563EB] font-semibold' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}"
                  >
                    <span class="w-1.5 h-1.5 rounded-full mr-2 {isSubActive ? 'bg-[#2563EB]' : 'bg-slate-300'}"></span>
                    <span class="truncate">{subItem.name}</span>
                  </button>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Sticky Footer -->
    <div class="mt-auto pt-3 border-t border-slate-200 text-xs">
      <div class="flex items-center justify-between px-2 py-1.5 rounded-lg bg-slate-50 border border-slate-200/80">
        <div class="flex items-center space-x-2">
          <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
          <span class="text-[11px] font-semibold text-slate-700">Cluster 1</span>
        </div>
        <span class="text-[10px] font-mono text-slate-500">FreeSWITCH 1.10</span>
      </div>
    </div>
  </aside>
{/if}
