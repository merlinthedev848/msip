<script lang="ts">
  import { onMount } from 'svelte';
  
  let currentPath = window.location.pathname;

  const handleNav = (path: string) => {
    window.history.pushState({}, '', path);
    currentPath = path;
    window.dispatchEvent(new CustomEvent('app-navigate', { detail: { path } }));
  };
  
  // Primary Navigation matching reference screenshot
  const primaryNav = [
    { 
      name: 'Dashboard', 
      path: '/', 
      alias: ['/', '/dashboard'],
      icon: 'dashboard' 
    },
    { 
      name: 'Extensions', 
      path: '/extensions', 
      alias: ['/extensions', '/extension'],
      icon: 'phone' 
    },
    { 
      name: 'Call Logs', 
      path: '/call-logs', 
      alias: ['/call-logs', '/cdrs', '/cdr'],
      icon: 'logs' 
    },
    { 
      name: 'Digital Receptionists', 
      path: '/digital-receptionists', 
      alias: ['/digital-receptionists', '/ivr', '/ivrs', '/ivr-builder'],
      icon: 'ivr' 
    },
    { 
      name: 'Ring Groups', 
      path: '/ring-groups', 
      alias: ['/ring-groups'],
      icon: 'users' 
    },
    { 
      name: 'Call Queues', 
      path: '/call-queues', 
      alias: ['/call-queues', '/wallboard'],
      icon: 'queue' 
    },
    { 
      name: 'Voicemail', 
      path: '/voicemail', 
      alias: ['/voicemail', '/voicemails'],
      icon: 'voicemail' 
    },
    { 
      name: 'System Settings', 
      path: '/system-settings', 
      alias: ['/system-settings', '/admins', '/admin'],
      icon: 'settings' 
    },
    { 
      name: 'Reports & Analytics', 
      path: '/reports-analytics', 
      alias: ['/reports-analytics', '/reports'],
      icon: 'reports' 
    }
  ];

  let showAdvanced = false;

  const advancedSections = [
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

  onMount(() => {
    const updatePath = () => { currentPath = window.location.pathname; };
    window.addEventListener('popstate', updatePath);
    window.addEventListener('app-navigate', (e: any) => { currentPath = e.detail.path; });

    return () => {
      window.removeEventListener('popstate', updatePath);
      window.removeEventListener('app-navigate', () => {});
    };
  });

  function isActive(item: { path: string; alias?: string[] }) {
    if (item.alias && item.alias.includes(currentPath)) return true;
    return currentPath === item.path;
  }
</script>

<aside class="w-64 bg-white h-[calc(100vh-52px)] overflow-y-auto border-r border-slate-200 p-3 flex flex-col flex-shrink-0 select-none shadow-[1px_0_4px_rgba(0,0,0,0.02)]">
  <!-- Primary Navigation Links -->
  <nav class="space-y-1">
    {#each primaryNav as item}
      {@const active = isActive(item)}
      <button 
        type="button"
        on:click={() => handleNav(item.path)}
        class="w-full flex items-center space-x-3 px-3.5 py-2.5 rounded-lg text-xs font-semibold transition-all duration-150 text-left {active ? 'bg-[#2563EB] text-white shadow-sm' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80'}"
      >
        <!-- Icon Rendering -->
        <span class="w-4 h-4 flex-shrink-0 flex items-center justify-center {active ? 'text-white' : 'text-slate-500'}">
          {#if item.icon === 'dashboard'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
          {:else if item.icon === 'phone'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
            </svg>
          {:else if item.icon === 'logs'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          {:else if item.icon === 'ivr'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 11.5V14m0-2.5v-6a1.5 1.5 0 113 0m-3 6a1.5 1.5 0 00-3 0v2a7.5 7.5 0 0015 0v-5a1.5 1.5 0 00-3 0m-6-3V11m0-5.5v-1a1.5 1.5 0 013 0v1m0 0V11m0-5.5a1.5 1.5 0 013 0v3m0 0V11" />
            </svg>
          {:else if item.icon === 'users'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          {:else if item.icon === 'queue'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
            </svg>
          {:else if item.icon === 'voicemail'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <circle cx="5.5" cy="11.5" r="3.5" stroke="currentColor" stroke-width="2" fill="none" />
              <circle cx="18.5" cy="11.5" r="3.5" stroke="currentColor" stroke-width="2" fill="none" />
              <line x1="5.5" y1="15" x2="18.5" y2="15" stroke="currentColor" stroke-width="2" />
            </svg>
          {:else if item.icon === 'settings'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          {:else if item.icon === 'reports'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          {/if}
        </span>
        <span class="truncate">{item.name}</span>
      </button>
    {/each}
  </nav>

  <!-- Divider -->
  <div class="my-4 border-t border-slate-200"></div>

  <!-- Advanced Telecom & Security Modules Toggle -->
  <div class="mb-2">
    <button 
      type="button" 
      on:click={() => showAdvanced = !showAdvanced} 
      class="w-full flex items-center justify-between px-3 py-2 text-[11px] font-bold text-slate-500 uppercase tracking-wider hover:text-slate-800 transition-colors"
    >
      <span>Advanced Modules</span>
      <svg class="w-3.5 h-3.5 transition-transform duration-200 {showAdvanced ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    {#if showAdvanced}
      <div class="space-y-4 pt-1 pl-1 pr-1 animate-in fade-in duration-150">
        {#each advancedSections as sec}
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

  <!-- Sticky Footer: Telecom Node Status -->
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
