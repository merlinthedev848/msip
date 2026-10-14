<script lang="ts">
  import { onMount } from 'svelte';
  
  export let collapsed = false;
  let currentPath = window.location.pathname;
  let user = JSON.parse(localStorage.getItem('pbx_user') || '{}');

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
    { name: 'Extensions & Users', path: '/extensions', icon: 'users' },
    { name: 'Phone Numbers (DIDs)', path: '/telephone-numbers', icon: 'phone-numbers' },
    { name: 'Call Logs & Traces', path: '/call-logs', icon: 'call-logs' },
    { name: 'Reports & Analytics', path: '/reports-analytics', icon: 'dashboards' },
    { 
      name: 'Call Management', 
      icon: 'call-groups', 
      expanded: false,
      children: [
        { name: 'Ring Groups', path: '/ring-groups' },
        { name: 'Auto Attendant (IVR)', path: '/digital-receptionists' },
        { name: 'Call Queues & Wallboard', path: '/call-queues' },
        { name: 'Voicemail Server', path: '/voicemail' },
        { name: 'Call Diverts & Forwarding', path: '/diverts' },
        { name: 'Time Conditions & Holidays', path: '/time-conditions' },
        { name: 'Hold Music & Recordings', path: '/audio-prompts' },
        { name: 'Paging & Call Pickup', path: '/paging-pickup' }
      ]
    },
    { 
      name: 'AI & Voice Intelligence', 
      icon: 'ai', 
      expanded: false,
      children: [
        { name: 'AI Virtual Receptionists', path: '/ai-assistants' },
        { name: 'AI Appointment Bookings', path: '/ai-bookings' }
      ]
    },
    { 
      name: 'Carriers & Routing', 
      icon: 'call-control', 
      expanded: false,
      children: [
        { name: 'SIP Trunks & Gateways', path: '/trunks' },
        { name: 'Inbound Call Rules', path: '/inbound-routes' },
        { name: 'Outbound Dial Rules', path: '/outbound-routes' },
        { name: 'Outbound Barring & AutoCLIP', path: '/outbound-restrictions' }
      ]
    },
    { 
      name: 'Directory & Devices', 
      icon: 'contacts', 
      expanded: false,
      children: [
        { name: 'Company Directory', path: '/contacts' },
        { name: 'IP Handsets & MACs', path: '/devices' },
        { name: 'Button & BLF Profiles', path: '/button-profiles' },
        { name: 'Auto-Provisioning', path: '/provisioning' },
        { name: 'Office Sites & Branches', path: '/sites' }
      ]
    },
    { 
      name: 'Collaboration', 
      icon: 'integration', 
      expanded: false,
      children: [
        { name: 'Team Chat', path: '/chat' },
        { name: 'SMS & MMS Messaging', path: '/sms' },
        { name: 'Video Conference Rooms', path: '/video-rooms' },
        { name: 'Security Cameras (RTSP)', path: '/cameras' }
      ]
    },
    { 
      name: 'Security & System', 
      icon: 'security', 
      expanded: false,
      children: [
        { name: 'System Administrators', path: '/admins' },
        { name: 'Tenant Domains', path: '/domains' },
        { name: 'SIP & RTP Settings', path: '/sip-settings' },
        { name: 'Feature Codes (*XX)', path: '/feature-codes' },
        { name: 'SIP & Jitter Diagnostics', path: '/diagnostics' },
        { name: 'Firewall & Access Control', path: '/firewall' },
        { name: 'Anti-Fraud Shield', path: '/fraud' },
        { name: 'API Keys & Webhooks', path: '/api-credentials' },
        { name: 'Change History', path: '/versioning' }
      ]
    }
  ];

  // Auto-expand accordion when a child route is active
  $: if (currentPath) {
    let changed = false;
    navItems.forEach(item => {
      if (item.children && item.children.some(c => c.path === currentPath)) {
        if (!item.expanded) {
          item.expanded = true;
          changed = true;
        }
      }
    });
    if (changed) {
      navItems = [...navItems];
    }
  }

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
    user = JSON.parse(localStorage.getItem('pbx_user') || '{}');
    const updatePath = () => { currentPath = window.location.pathname; };
    window.addEventListener('popstate', updatePath);
    window.addEventListener('app-navigate', (e: any) => { currentPath = e.detail.path; });

    return () => {
      window.removeEventListener('popstate', updatePath);
      window.removeEventListener('app-navigate', () => {});
    };
  });
</script>

<aside class="{collapsed ? 'w-16' : 'w-60'} transition-all duration-200 bg-white h-[calc(100vh-50px)] overflow-y-auto border-r border-slate-200 py-3 flex flex-col flex-shrink-0 select-none shadow-[1px_0_4px_rgba(0,0,0,0.02)]">
  <!-- User Profile Header matching customer portal -->
  <div class="px-3 pb-3 mb-1 border-b border-slate-100 flex items-center space-x-2.5 {collapsed ? 'justify-center' : ''}">
    <div class="w-8 h-8 rounded-full bg-slate-200 text-slate-600 flex items-center justify-center font-bold text-xs shadow-sm flex-shrink-0">
      <svg class="w-4 h-4 text-slate-500" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
      </svg>
    </div>
    {#if !collapsed}
      <div class="min-w-0 flex-1">
        <div class="text-xs font-bold text-slate-800 truncate">{user.email ? user.email.split('@')[0] : 'Portal Administrator'}</div>
        <div class="text-[10px] text-slate-400 capitalize truncate">{user.role || 'Supervisor'}</div>
      </div>
    {/if}
  </div>

  <!-- Navigation List -->
  <nav class="space-y-0.5 px-2">
    {#each navItems as item, idx}
      {@const active = isActive(item)}
      <div>
        <button 
          type="button"
          on:click={() => toggleItem(idx)}
          title={item.name}
          class="w-full flex items-center {collapsed ? 'justify-center' : 'justify-between'} px-3 py-2 rounded-md text-xs transition-colors {active ? 'text-[#0284C7] font-bold bg-sky-50/70' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50 font-medium'}"
        >
          <div class="flex items-center space-x-2.5 truncate">
            <span class="w-4 h-4 flex-shrink-0 flex items-center justify-center {active ? 'text-[#0284C7]' : 'text-slate-500'}">
              {#if item.icon === 'dashboard'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-width="2" stroke-linecap="round" d="M12 7v5l3 3"/></svg>
              {:else if item.icon === 'buttons'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/></svg>
              {:else if item.icon === 'contacts'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/></svg>
              {:else if item.icon === 'call-groups'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/></svg>
              {:else if item.icon === 'menus'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"/></svg>
              {:else if item.icon === 'queues'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 100-6 3 3 0 000 6z"/></svg>
              {:else if item.icon === 'ai'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>
              {:else if item.icon === 'phone-numbers'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg>
              {:else if item.icon === 'sites'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>
              {:else if item.icon === 'users'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>
              {:else if item.icon === 'voicemail'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="6" cy="12" r="3" stroke-width="2"/><circle cx="18" cy="12" r="3" stroke-width="2"/><path stroke-width="2" d="M6 15h12"/></svg>
              {:else if item.icon === 'dashboards'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>
              {:else if item.icon === 'call-logs'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
              {:else if item.icon === 'advanced'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
              {:else if item.icon === 'call-control'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/></svg>
              {:else if item.icon === 'integration'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z"/></svg>
              {:else if item.icon === 'security'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/></svg>
              {:else if item.icon === 'system'}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/></svg>
              {/if}
            </span>
            {#if !collapsed}
              <span class="truncate">{item.name}</span>
            {/if}
          </div>

          {#if !collapsed && item.children}
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
