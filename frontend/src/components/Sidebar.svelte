<script lang="ts">
  import { globalHistory } from "svelte-routing/src/history";
  
  let currentPath = window.location.pathname;
  
  // Listen to router changes to highlight the active menu item
  globalHistory.listen(({ location }) => {
    currentPath = location.pathname;
  });

  const menuSections = [
    {
      title: "Dashboard & Telemetry",
      items: [
        { name: "Live Status Dashboard", path: "/" },
        { name: "Call History (CDRs)", path: "/cdrs" }
      ]
    },
    {
      title: "Call Routing Engine",
      items: [
        { name: "Inbound Routes", path: "/inbound-routes" },
        { name: "Outbound Routes", path: "/outbound-routes" },
        { name: "SIP Trunks", path: "/trunks" }
      ]
    },
    {
      title: "Core Telephony Features",
      items: [
        { name: "Extensions Management", path: "/extensions" },
        { name: "IVR / Receptionist", path: "/ivr" },
        { name: "Visual Voicemail", path: "/voicemail" }
      ]
    }
  ];
</script>

<aside class="w-72 bg-gray-950 h-screen overflow-y-auto border-r border-gray-800 p-4 shadow-2xl flex flex-col hidden md:flex sticky top-0 z-30">
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
  
  <nav class="space-y-8 flex-1">
    {#each menuSections as section}
      <div>
        <h2 class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-3 px-2">{section.title}</h2>
        <ul class="space-y-1">
          {#each section.items as item}
            <li>
              <a 
                href={item.path} 
                class="block px-3 py-2 rounded-xl text-sm transition-all duration-200 ease-in-out group flex items-center {currentPath === item.path ? 'bg-gray-800/80 text-white font-medium shadow-sm border border-gray-700/50' : 'text-gray-400 hover:text-white hover:bg-gray-800/50'}"
              >
                <span class="w-1.5 h-1.5 rounded-full mr-3 transition-colors duration-200 {currentPath === item.path ? 'bg-indigo-500 shadow-[0_0_8px_rgba(99,102,241,0.8)]' : 'bg-gray-700 group-hover:bg-gray-500'}"></span>
                {item.name}
              </a>
            </li>
          {/each}
        </ul>
      </div>
    {/each}
  </nav>

  <div class="mt-8 px-4 pt-4 border-t border-gray-800/50">
    <div class="flex items-center space-x-3 text-sm text-gray-400 bg-gray-900/50 p-3 rounded-xl border border-gray-800">
      <div class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.8)]"></div>
      <span class="font-mono text-[11px] uppercase tracking-wider">System Online</span>
    </div>
  </div>
</aside>
