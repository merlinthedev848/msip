<script lang="ts">
  import { onMount } from 'svelte';
  import { slide } from 'svelte/transition';
  
  let currentPath = window.location.pathname;

  const handleNav = (path) => {
    window.history.pushState({}, '', path);
    currentPath = path;
    window.dispatchEvent(new CustomEvent('app-navigate', { detail: { path } }));
  };
  
  let menuSections = [
    {
      title: 'Platform',
      items: [
        { name: 'Global Dashboard', path: '/' },
        { name: 'Domains & Tenants', path: '/domains' },
        { name: 'Admins & Roles', path: '/admins' },
        { name: 'Call Center Wallboard', path: '/wallboard' }
      ]
    },
    {
      title: 'Call Routing',
      items: [
        { name: 'Inbound Routes', path: '/inbound-routes' },
        { name: 'Outbound Routes', path: '/outbound-routes' },
        { name: 'SIP Trunks', path: '/trunks' }
      ]
    },
    {
      title: 'Endpoints & Hardware',
      items: [
        { name: 'SIP Extensions', path: '/extensions' },
        { name: 'Auto Provisioning', path: '/provisioning' }
      ]
    },
    {
      title: 'Voice Applications',
      items: [
        { name: 'IVR / Auto Attendant', path: '/ivr' },
        { name: 'Voicemail Server', path: '/voicemail' },
        { name: 'Call Detail Records', path: '/cdrs' }
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
      title: 'Security & API',
      items: [
        { name: 'API Credentials', path: '/api-credentials' },
        { name: 'Event Webhooks', path: '/webhooks' },
        { name: 'Rate Throttling', path: '/throttling' },
        { name: 'Firewall & ACLs', path: '/firewall' },
        { name: 'Anti-Fraud System', path: '/fraud' },
        { name: 'Version Control', path: '/versioning' }
      ]
    }
  ].map(section => ({
    ...section,
    expanded: section.items.some(item => item.path === currentPath)
  }));
  
  onMount(() => {
    const updatePath = () => { currentPath = window.location.pathname; };
    window.addEventListener('popstate', updatePath);

    return () => {
      window.removeEventListener('popstate', updatePath);
    };
  });

  const toggleSection = (index) => {
    menuSections[index].expanded = !menuSections[index].expanded;
  };
</script>

<aside class="w-72 bg-white h-screen overflow-y-auto scrollbar-hide border-r border-slate-200 p-4 shadow-sm flex flex-col hidden md:flex sticky top-0 z-30">
  <div class="mb-8 px-2 pt-2">
    <div class="flex items-center space-x-3 mb-2">
      <div class="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center shadow-sm shadow-blue-500/20">
        <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
      </div>
      <div>
        <h1 class="text-xl font-black text-slate-900 tracking-tight">MSIP PBX</h1>
      </div>
    </div>
  </div>
  
  <nav class="space-y-4 flex-1 pb-10">
    {#each menuSections as section, i}
      <div class="bg-slate-50 rounded-xl overflow-hidden border border-slate-200">
        <button 
          class="w-full flex items-center justify-between p-3 focus:outline-none group hover:bg-slate-100 transition-colors"
          on:click={() => toggleSection(i)}
        >
          <h2 class="text-[11px] font-bold text-slate-500 group-hover:text-slate-900 uppercase tracking-wider transition-colors">{section.title}</h2>
          <svg 
            class="w-4 h-4 text-slate-400 transition-transform duration-200 {section.expanded ? 'rotate-180 text-slate-700' : ''}" 
            fill="none" stroke="currentColor" viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
          </svg>
        </button>
        
        {#if section.expanded}
          <ul class="px-2 pb-2 space-y-1" transition:slide={{duration: 200}}>
            {#each section.items as item}
              <li>
                <a 
                  href={item.path} 
                  on:click|preventDefault={() => handleNav(item.path)}
                  class="block px-3 py-2 rounded-xl text-sm transition-all duration-200 ease-in-out group flex items-center font-medium {currentPath === item.path ? 'bg-blue-50 text-blue-700 shadow-sm border border-blue-100' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100 border border-transparent'}"
                >
                  <span class="w-1.5 h-1.5 rounded-full mr-3 transition-colors duration-200 flex-shrink-0 {currentPath === item.path ? 'bg-blue-600' : 'bg-slate-300 group-hover:bg-slate-400'}"></span>
                  {item.name}
                </a>
              </li>
            {/each}
          </ul>
        {/if}
      </div>
    {/each}
  </nav>

  <div class="mt-4 px-4 pt-4 border-t border-slate-200 sticky bottom-0 bg-white pb-4 z-10">
    <div class="flex items-center space-x-3 text-sm text-slate-500 bg-slate-50 p-3 rounded-xl border border-slate-200">
      <div class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shadow-sm"></div>
      <span class="font-mono text-[11px] uppercase tracking-wider">System Online</span>
    </div>
  </div>
</aside>
