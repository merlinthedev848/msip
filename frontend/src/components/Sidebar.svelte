<script lang="ts">
  const menuSections = [
    {
      title: "Dashboard & Telemetry",
      items: ["Live Status Dashboard", "Call History (CDRs)"]
    },
    {
      title: "Call Routing Engine",
      items: ["Inbound Routes", "Outbound Routes", "SIP Trunks"]
    },
    {
      title: "Core Telephony Features",
      items: ["Extensions Management", "IVR / Receptionist", "Ring Groups", "Visual Voicemail"]
    },
    {
      title: "Video & Collaboration",
      items: ["Video Rooms", "Camera Integration"]
    },
    {
      title: "Messaging & Chat",
      items: ["SMS Trunks", "Team Chat Channels"]
    },
    {
      title: "Developer & API Gateway",
      items: ["API Credentials", "Webhook Subscriptions", "Throttling & IP Control"]
    },
    {
      title: "Security & Operations",
      items: ["Firewall & IDS", "Anti-Fraud Protections"]
    },
    {
      title: "System Updates",
      items: ["Version Control"]
    }
  ];
</script>

<aside class="w-72 bg-gray-950 h-screen overflow-y-auto border-r border-gray-800 p-4 shadow-xl flex flex-col hidden md:flex sticky top-0">
  <div class="mb-8 px-2">
    <h1 class="text-2xl font-black bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-indigo-500 tracking-tight">MSIP PBX</h1>
    <p class="text-xs text-gray-500 font-mono mt-1">LIGHTWEIGHT UC ENGINE</p>
  </div>
  
  <nav class="space-y-6 flex-1">
    {#each menuSections as section}
      <div>
        <h2 class="text-xs font-bold text-gray-400 uppercase tracking-wider mb-3 px-2">{section.title}</h2>
        <ul class="space-y-1">
          {#each section.items as item}
            <li>
              <a href="#" class="block px-2 py-2 rounded-lg text-sm text-gray-300 hover:text-white hover:bg-gray-800 transition-all duration-200 ease-in-out group flex items-center">
                <span class="w-1.5 h-1.5 rounded-full bg-gray-600 mr-2 group-hover:bg-indigo-400 transition-colors duration-200"></span>
                {item}
              </a>
            </li>
          {/each}
        </ul>
      </div>
    {/each}
  </nav>

  <div class="mt-8 px-2 pt-4 border-t border-gray-800">
    <div class="flex items-center space-x-3 text-sm text-gray-400">
      <div class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
      <span class="font-mono text-xs">System Online</span>
    </div>
  </div>
</aside>
