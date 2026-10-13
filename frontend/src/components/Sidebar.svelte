<script lang="ts">
  import { onMount } from 'svelte';
  
  let currentPath = window.location.pathname;

  const handleNav = (path: string) => {
    window.history.pushState({}, '', path);
    currentPath = path;
    window.dispatchEvent(new CustomEvent('app-navigate', { detail: { path } }));
  };

  interface NavItem {
    name: string;
    path?: string;
    icon: string;
    expanded?: boolean;
    children?: { name: string; path: string }[];
  }

  let navItems: NavItem[] = [
    { name: 'Dashboard', path: '/', icon: 'dashboard' },
    { 
      name: 'Extension and Trunk', 
      icon: 'ext-trunk', 
      expanded: true,
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
      expanded: true,
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
        { name: 'Admins & Roles', path: '/admins' },
        { name: 'Change History', path: '/versioning' },
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
        { name: 'API Credentials', path: '/api-credentials' }
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

  function toggleItem(index: number) {
    if (navItems[index].children) {
      navItems[index].expanded = !navItems[index].expanded;
    } else if (navItems[index].path) {
      handleNav(navItems[index].path!);
    }
  }

  function isActive(item: NavItem) {
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

<aside class="w-60 bg-white h-[calc(100vh-50px)] overflow-y-auto border-r border-slate-200 py-3 flex flex-col flex-shrink-0 select-none shadow-[1px_0_4px_rgba(0,0,0,0.02)]">
  <!-- Brand Header -->
  <div class="px-4 pb-3 flex items-center space-x-2.5">
    <div class="w-7 h-7 flex items-center justify-center text-[#0284C7]">
      <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none">
        <path d="M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12" stroke="#0284C7" stroke-width="2.5" stroke-linecap="round"/>
        <path d="M12 6V12L16 14" stroke="#0EA5E9" stroke-width="2.5" stroke-linecap="round"/>
      </svg>
    </div>
    <span class="text-base font-bold text-slate-800 tracking-tight">MSIP PBX</span>
  </div>

  <!-- Navigation List -->
  <nav class="space-y-0.5 px-2">
    {#each navItems as item, idx}
      {@const active = isActive(item)}
      <div>
        <button 
          type="button"
          on:click={() => toggleItem(idx)}
          class="w-full flex items-center justify-between px-3 py-2 rounded-md text-xs transition-colors {active ? 'text-[#0284C7] font-bold bg-sky-50/70' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50 font-medium'}"
        >
          <div class="flex items-center space-x-2.5 truncate">
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

          {#if item.children}
            <svg class="w-3.5 h-3.5 text-slate-400 transition-transform duration-150 {item.expanded ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          {/if}
        </button>

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
