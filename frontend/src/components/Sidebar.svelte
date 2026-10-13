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

<aside class="w-72 bg-[#090D16] h-screen overflow-y-auto scrollbar-hide border-r border-slate-800/90 p-4 shadow-2xl flex flex-col hidden md:flex sticky top-0 z-30 select-none">
  <!-- Brand Header -->
  <div class="mb-6 px-2 pt-2">
    <div class="flex items-center space-x-3 mb-1">
      <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-indigo-600 via-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/30 border border-indigo-400/30">
        <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path>
        </svg>
      </div>
      <div>
        <h1 class="text-xl font-black text-white tracking-tight flex items-center gap-1.5">
          MSIP <span class="text-xs font-bold px-1.5 py-0.5 rounded bg-indigo-500/20 text-indigo-400 border border-indigo-500/30 font-mono">PBX</span>
        </h1>
        <p class="text-[10px] text-slate-400 uppercase tracking-widest font-semibold">Telecom Cloud</p>
      </div>
    </div>
  </div>
  
  <!-- Navigation Sections -->
  <nav class="space-y-3 flex-1 pb-8">
    {#each menuSections as section, i}
      <div class="bg-slate-900/50 rounded-xl overflow-hidden border border-slate-800/80 transition-all duration-200">
        <button 
          class="w-full flex items-center justify-between p-3 focus:outline-none group hover:bg-slate-800/50 transition-colors"
          on:click={() => toggleSection(i)}
        >
          <h2 class="text-[11px] font-bold text-slate-400 group-hover:text-slate-200 uppercase tracking-wider transition-colors">{section.title}</h2>
          <svg 
            class="w-4 h-4 text-slate-500 transition-transform duration-200 {section.expanded ? 'rotate-180 text-indigo-400' : ''}" 
            fill="none" stroke="currentColor" viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
          </svg>
        </button>
        
        {#if section.expanded}
          <ul class="px-2 pb-2 space-y-0.5" transition:slide={{duration: 200}}>
            {#each section.items as item}
              <li>
                <a 
                  href={item.path} 
                  on:click|preventDefault={() => handleNav(item.path)}
                  class="block px-3 py-2 rounded-lg text-sm transition-all duration-200 ease-in-out group flex items-center font-medium {currentPath === item.path ? 'bg-gradient-to-r from-indigo-500/20 via-purple-500/15 to-transparent text-indigo-300 border border-indigo-500/40 shadow-sm shadow-indigo-500/10' : 'text-slate-400 hover:text-white hover:bg-slate-800/60 border border-transparent'}"
                >
                  <span class="w-1.5 h-1.5 rounded-full mr-2.5 transition-colors duration-200 flex-shrink-0 {currentPath === item.path ? 'bg-indigo-400 shadow-sm shadow-indigo-400' : 'bg-slate-600 group-hover:bg-slate-400'}"></span>
                  {item.name}
                </a>
              </li>
            {/each}
          </ul>
        {/if}
      </div>
    {/each}
  </nav>

  <!-- Sticky Footer System Status -->
  <div class="mt-auto pt-3 border-t border-slate-800/90 bg-[#090D16] sticky bottom-0 pb-2 z-10">
    <div class="flex items-center justify-between text-xs text-slate-400 bg-slate-900/90 p-2.5 rounded-xl border border-slate-800/80">
      <div class="flex items-center space-x-2.5">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
        </span>
        <span class="font-mono text-[11px] font-semibold text-emerald-400 tracking-wider">SYSTEM ONLINE</span>
      </div>
      <span class="text-[10px] font-mono text-slate-500">FS 1.10</span>
    </div>
  </div>
</aside>
