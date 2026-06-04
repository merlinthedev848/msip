<script lang="ts">
  import { onMount, onDestroy } from 'svelte';

  let activeCalls = 0;
  let activeExtensions = 0;
  let cpuUsage = 0;
  let ramUsage = 0;
  
  // Historical data for sparklines
  let callHistory = Array(20).fill(0);
  let cpuHistory = Array(20).fill(0);
  let ramHistory = Array(20).fill(0);

  // SVG Path generators
  $: callPoints = callHistory.map((val, i) => `${i * 5.26},${30 - (val / 60 * 30)}`).join(' ');
  $: callPath = `M 0,30 L ${callPoints} L 100,30 Z`;

  $: cpuPoints = cpuHistory.map((val, i) => `${i * 5.26},${30 - (val / 100 * 30)}`).join(' ');
  $: cpuPath = `M 0,30 L ${cpuPoints} L 100,30 Z`;

  $: ramPoints = ramHistory.map((val, i) => `${i * 5.26},${30 - (val / 500 * 30)}`).join(' ');
  $: ramPath = `M 0,30 L ${ramPoints} L 100,30 Z`;

  let liveCallList = [];
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  
  <header class="flex flex-col md:flex-row md:items-end justify-between gap-4">
    <div>
      <h1 class="text-4xl font-black tracking-tighter text-slate-900 mb-2 flex items-center">
        System Overview
      </h1>
    </div>
    
    <div class="flex items-center space-x-3 bg-white p-1.5 rounded-xl border border-slate-200/50 ">
      <button class="px-4 py-1.5 text-xs font-bold bg-slate-100 text-slate-900 rounded-lg shadow-sm">1H</button>
      <button class="px-4 py-1.5 text-xs font-bold text-slate-500 hover:text-slate-900 transition-colors">24H</button>
      <button class="px-4 py-1.5 text-xs font-bold text-slate-500 hover:text-slate-900 transition-colors">7D</button>
    </div>
  </header>

  <!-- Sleek Top Stats -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    
    <!-- Active Calls -->
    <div class="relative overflow-hidden rounded-3xl bg-white/40 border border-slate-200/60 p-6  shadow-2xl group hover:border-indigo-500/30 transition-all duration-300">
      <div class="absolute top-0 right-0 p-6 opacity-5 group-hover:opacity-10 transition-opacity transform group-hover:scale-110 duration-500">
        <svg class="w-24 h-24 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
      </div>
      <p class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] mb-2">Live Streams</p>
      <h2 class="text-5xl font-black text-slate-900 tracking-tighter mb-6 font-mono">{activeCalls}</h2>
      
      <div class="absolute bottom-0 left-0 right-0 h-16 pointer-events-none">
        <svg viewBox="0 0 100 30" class="w-full h-full" preserveAspectRatio="none">
          <defs>
            <linearGradient id="gradIndigo" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#6366f1" stop-opacity="0.3" />
              <stop offset="100%" stop-color="#6366f1" stop-opacity="0" />
            </linearGradient>
          </defs>
          <path d={callPath} fill="url(#gradIndigo)" />
          <polyline points={callPoints} fill="none" stroke="#818cf8" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
    </div>

    <!-- Active Endpoints -->
    <div class="relative overflow-hidden rounded-3xl bg-white/40 border border-slate-200/60 p-6  shadow-2xl group hover:border-emerald-500/30 transition-all duration-300">
      <div class="absolute top-0 right-0 p-6 opacity-5 group-hover:opacity-10 transition-opacity transform group-hover:scale-110 duration-500">
        <svg class="w-24 h-24 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
      </div>
      <p class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] mb-2">Registered Endpoints</p>
      <div class="flex items-baseline space-x-2 mb-6">
        <h2 class="text-5xl font-black text-slate-900 tracking-tighter font-mono">{activeExtensions}</h2>
        <span class="text-sm font-bold text-gray-600">/ 50</span>
      </div>
      
      <!-- Fake mini progress bar -->
      <div class="absolute bottom-6 left-6 right-6 h-1.5 bg-slate-50 rounded-full overflow-hidden">
        <div class="h-full bg-emerald-500 rounded-full transition-all duration-500 ease-out" style="width: {(activeExtensions / 50) * 100}%"></div>
      </div>
    </div>

    <!-- CPU Usage -->
    <div class="relative overflow-hidden rounded-3xl bg-white/40 border border-slate-200/60 p-6  shadow-2xl group hover:border-cyan-500/30 transition-all duration-300">
      <p class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] mb-2">DSP Core Load</p>
      <div class="flex items-baseline space-x-1 mb-6">
        <h2 class="text-5xl font-black text-slate-900 tracking-tighter font-mono">{cpuUsage}</h2>
        <span class="text-lg font-bold text-slate-500">%</span>
      </div>
      
      <div class="absolute bottom-0 left-0 right-0 h-16 pointer-events-none">
        <svg viewBox="0 0 100 30" class="w-full h-full" preserveAspectRatio="none">
          <defs>
            <linearGradient id="gradCyan" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#06b6d4" stop-opacity="0.3" />
              <stop offset="100%" stop-color="#06b6d4" stop-opacity="0" />
            </linearGradient>
          </defs>
          <path d={cpuPath} fill="url(#gradCyan)" />
          <polyline points={cpuPoints} fill="none" stroke="#22d3ee" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
    </div>

    <!-- RAM Usage -->
    <div class="relative overflow-hidden rounded-3xl bg-white/40 border border-slate-200/60 p-6  shadow-2xl group hover:border-violet-500/30 transition-all duration-300">
      <p class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] mb-2">Memory Footprint</p>
      <div class="flex items-baseline space-x-1 mb-6">
        <h2 class="text-5xl font-black text-slate-900 tracking-tighter font-mono">{ramUsage}</h2>
        <span class="text-lg font-bold text-slate-500">MB</span>
      </div>
      
      <div class="absolute bottom-0 left-0 right-0 h-16 pointer-events-none">
        <svg viewBox="0 0 100 30" class="w-full h-full" preserveAspectRatio="none">
          <defs>
            <linearGradient id="gradViolet" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#8b5cf6" stop-opacity="0.3" />
              <stop offset="100%" stop-color="#8b5cf6" stop-opacity="0" />
            </linearGradient>
          </defs>
          <path d={ramPath} fill="url(#gradViolet)" />
          <polyline points={ramPoints} fill="none" stroke="#a78bfa" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
    </div>

  </div>

  <!-- Bottom Data Area -->
  <div class="space-y-6 mt-8">
    
    <!-- Deep Table -->
    <div class="rounded-3xl bg-white/40 border border-slate-200/60  shadow-2xl flex flex-col overflow-hidden">
      <div class="px-6 py-5 border-b border-slate-200/50 flex justify-between items-center bg-white/20">
        <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
          <div class="w-2 h-2 rounded-full bg-indigo-500 mr-3"></div>
          Active Edge Streams
        </h3>
        <button class="text-[10px] font-bold text-blue-600 hover:text-indigo-300 uppercase tracking-widest">View All CDRs &rarr;</button>
      </div>
      
      <div class="flex-1 overflow-x-auto p-2">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="text-[10px] text-slate-500 uppercase tracking-[0.2em]">
              <th class="px-6 py-4 font-bold">Source Endpoint</th>
              <th class="px-4 py-4 font-bold text-center">Dir</th>
              <th class="px-6 py-4 font-bold">Destination</th>
              <th class="px-6 py-4 font-bold">Time</th>
              <th class="px-6 py-4 font-bold">State</th>
            </tr>
          </thead>
          <tbody class="text-sm text-slate-700">
            {#each liveCallList as call}
              <tr class="hover:bg-slate-50 transition-colors group cursor-pointer border-b border-slate-200/30 last:border-0">
                <td class="px-6 py-4">
                  <div class="flex items-center space-x-4">
                    <div class="w-10 h-10 rounded-full {call.avatar} flex items-center justify-center text-slate-900 font-bold text-xs shadow-inner">
                      {call.name.charAt(0)}
                    </div>
                    <div>
                      <p class="font-bold text-slate-900 group-hover:text-indigo-300 transition-colors">{call.name}</p>
                      <p class="text-xs text-slate-500 font-mono">{call.caller}</p>
                    </div>
                  </div>
                </td>
                <td class="px-4 py-4 text-center">
                  <div class="w-8 h-8 rounded-full bg-slate-50 flex items-center justify-center mx-auto border border-slate-200">
                    <svg class="w-4 h-4 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"></path></svg>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <p class="font-bold text-slate-700">{call.callee}</p>
                  <p class="text-[10px] font-mono text-slate-500 uppercase mt-0.5 tracking-wider">{call.codec}</p>
                </td>
                <td class="px-6 py-4 font-mono text-slate-500 text-xs">{call.duration}</td>
                <td class="px-6 py-4">
                  <span class="px-3 py-1 text-[10px] font-bold uppercase tracking-widest rounded-full 
                    {call.state === 'Talking' ? 'bg-emerald-50 text-emerald-400 border border-emerald-500/20' : 
                     call.state === 'Waiting' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 
                     'bg-slate-100 text-slate-500 border border-slate-200'}">
                    {call.state}
                  </span>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- Live Terminal -->
    <div class="rounded-3xl bg-[#0a0a0f] border border-slate-200/80 shadow-2xl flex flex-col overflow-hidden relative">
      <!-- Glossy top reflection -->
      <div class="absolute top-0 left-0 right-0 h-32 bg-gradient-to-b from-white/[0.03] to-transparent pointer-events-none"></div>
      
      <div class="px-5 py-4 border-b border-gray-900 flex justify-between items-center bg-[#0d0d12] z-10">
        <h3 class="text-xs font-bold text-slate-500 uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
          fs_cli Console
        </h3>
        <div class="flex space-x-1.5">
          <div class="w-2.5 h-2.5 rounded-full bg-red-500/20 border border-red-500/50"></div>
          <div class="w-2.5 h-2.5 rounded-full bg-yellow-500/20 border border-yellow-500/50"></div>
          <div class="w-2.5 h-2.5 rounded-full bg-green-500/20 border border-green-500/50"></div>
        </div>
      </div>
      
      <div class="flex-1 p-5 font-mono text-[11px] leading-relaxed overflow-y-auto z-10 space-y-2">
        <div><span class="text-gray-600">14:32:01.004</span> <span class="text-blue-600">[INFO]</span> <span class="text-slate-700">fs_cli: FreeSWITCH core initialized</span></div>
        <div><span class="text-gray-600">14:32:05.112</span> <span class="text-blue-600">[INFO]</span> <span class="text-slate-700">mod_sofia: SIP profile 'internal' started</span></div>
        <div><span class="text-gray-600">14:33:12.981</span> <span class="text-emerald-400">[EVNT]</span> <span class="text-slate-700">endpoint_1001 REGISTERED from 192.168.1.55</span></div>
        <div><span class="text-gray-600">14:34:44.204</span> <span class="text-yellow-400">[WARN]</span> <span class="text-slate-700">mod_event_socket: Auth failed 10.0.0.9</span></div>
        <div><span class="text-gray-600">14:35:01.442</span> <span class="text-emerald-400">[EVNT]</span> <span class="text-slate-700">endpoint_1045 REGISTERED from 10.0.0.12</span></div>
        <div><span class="text-gray-600">14:40:22.001</span> <span class="text-cyan-400">[ROUT]</span> <span class="text-slate-700">Dialplan: Call routed to IVR_Main</span></div>
        <div><span class="text-gray-600">14:41:10.155</span> <span class="text-emerald-400">[EVNT]</span> <span class="text-slate-700">endpoint_1002 REGISTERED from 192.168.1.105</span></div>
        <div class="text-red-400"><span class="text-red-500/50">14:42:00.893</span> [ERRR] mod_sofia: Timeout contacting sip_provider_primary</div>
        <div><span class="text-gray-600">14:42:01.002</span> <span class="text-blue-600">[INFO]</span> <span class="text-slate-700">failover: Switching to secondary trunk</span></div>
        <div class="mt-4 text-emerald-400 animate-pulse">root@pbx-core:~# _</div>
      </div>
    </div>
    
  </div>

</div>

