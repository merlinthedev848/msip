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

<aside class="w-72 bg-gray-950 h-screen overflow-y-auto scrollbar-hide border-r border-gray-800 p-4 shadow-2xl flex flex-col hidden md:flex sticky top-0 z-30">
  <div class="mb-8 px-2 pt-2">
    <div class="flex items-center space-x-3 mb-2">
      <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-violet-600 flex items-center justify-center shadow-lg shadow-indigo-500/20">
        <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
      </div>
      <div>
        <h1 class="text-xl font-black text-white tracking-tight">MSIP PBX</h1>
      </div>
    </div>
    <p class="text-[10px] text-gray-500 font-mono tracking-widest uppercase">Lightweight Engine</p>
  </div>
  
  <nav class="space-y-4 flex-1 pb-10">
    {#each menuSections as section, i}
      <div class="bg-gray-900/30 rounded-xl overflow-hidden border border-gray-800/50">
        <button 
          class="w-full flex items-center justify-between p-3 focus:outline-none group hover:bg-gray-800/40 transition-colors"
          on:click={() => toggleSection(i)}
        >
          <h2 class="text-[11px] font-bold text-gray-400 group-hover:text-white uppercase tracking-wider transition-colors">{section.title}</h2>
          <svg 
            class="w-4 h-4 text-gray-500 transition-transform duration-200 {section.expanded ? 'rotate-180 text-white' : ''}" 
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
                  class="block px-3 py-2 rounded-xl text-sm transition-all duration-200 ease-in-out group flex items-center font-medium {currentPath === item.path ? 'bg-gray-800/80 text-white shadow-sm border border-gray-700/50' : 'text-gray-400 hover:text-white hover:bg-gray-800/50 border border-transparent'}"
                >
                  <span class="w-1.5 h-1.5 rounded-full mr-3 transition-colors duration-200 flex-shrink-0 {currentPath === item.path ? 'bg-indigo-500 shadow-[0_0_8px_rgba(99,102,241,0.8)]' : 'bg-gray-700 group-hover:bg-gray-500'}"></span>
                  {item.name}
                </a>
              </li>
            {/each}
          </ul>
        {/if}
      </div>
    {/each}
  </nav>

  <div class="mt-4 px-4 pt-4 border-t border-gray-800/50 sticky bottom-0 bg-gray-950 pb-4 z-10">
    <div class="flex items-center space-x-3 text-sm text-gray-400 bg-gray-900/50 p-3 rounded-xl border border-gray-800">
      <div class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.8)]"></div>
      <span class="font-mono text-[11px] uppercase tracking-wider">System Online</span>
    </div>
  </div>
</aside>
