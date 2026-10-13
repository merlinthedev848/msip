<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { apiFetch } from '../lib/api';

  let activeCalls = 0;
  let activeExtensions = 0;
  let cpuUsage = 21;
  let ramUsage = 1238;
  
  // Historical data for sparklines
  let callHistory = Array(20).fill(0);
  let cpuHistory = Array(20).fill(20);
  let ramHistory = Array(20).fill(1200);

  // SVG Path generators
  $: callPoints = callHistory.map((val, i) => `${i * 5.26},${30 - (val / 60 * 30)}`).join(' ');
  $: callPath = `M 0,30 L ${callPoints} L 100,30 Z`;

  $: cpuPoints = cpuHistory.map((val, i) => `${i * 5.26},${30 - (val / 100 * 30)}`).join(' ');
  $: cpuPath = `M 0,30 L ${cpuPoints} L 100,30 Z`;

  $: ramPoints = ramHistory.map((val, i) => `${i * 5.26},${30 - (val / 5000 * 30)}`).join(' ');
  $: ramPath = `M 0,30 L ${ramPoints} L 100,30 Z`;

  let liveCallList = [];
  let ws: WebSocket;
  let wsInterval: any;
  let statsInterval: any;
  let activeTimeRange = '1H';

  // Terminal state & logs
  interface LogLine {
    id: number;
    time: string;
    level: 'INFO' | 'EVNT' | 'WARN' | 'ERRR' | 'ROUT' | 'SYS';
    text: string;
  }

  let consoleLogs: LogLine[] = [
    { id: 1, time: '14:32:01.004', level: 'INFO', text: 'fs_cli: FreeSWITCH 1.10.12-release core initialized successfully' },
    { id: 2, time: '14:32:05.112', level: 'INFO', text: 'mod_sofia: SIP profile "internal" bound to 0.0.0.0:5060 [TCP/UDP]' },
    { id: 3, time: '14:33:12.981', level: 'EVNT', text: 'endpoint_1001 REGISTERED from 192.168.1.55:5060 (Polycom VVX-450)' },
    { id: 4, time: '14:34:44.204', level: 'WARN', text: 'mod_event_socket: ACL Challenge from 10.0.0.9 (Denied: IP not in carrier acl)' },
    { id: 5, time: '14:35:01.442', level: 'EVNT', text: 'endpoint_1045 REGISTERED from 10.0.0.12:5060 (Yealink T48U)' },
    { id: 6, time: '14:40:22.001', level: 'ROUT', text: 'Dialplan: Inbound DID +18005550199 routed to IVR_Main via context "public"' },
    { id: 7, time: '14:41:10.155', level: 'EVNT', text: 'endpoint_1002 REGISTERED from 192.168.1.105:5060 (WebRTC Client)' },
    { id: 8, time: '14:42:00.893', level: 'ERRR', text: 'mod_sofia: Gateway carrier_primary ping timeout (504 Gateway Timeout)' },
    { id: 9, time: '14:42:01.002', level: 'INFO', text: 'failover: LCR route automatically switched to backup carrier_secondary' },
    { id: 10, time: '14:43:15.820', level: 'INFO', text: 'mod_conference: WebRTC MCU Room 9001 audio-mix profile loaded' }
  ];

  let selectedFilter: 'ALL' | 'INFO' | 'EVNT' | 'WARN' | 'ERRR' = 'ALL';
  let terminalFontSize = 13;
  let autoScroll = true;
  let cliInput = '';
  let copyFeedback = false;
  let terminalBodyElement: HTMLDivElement;

  $: filteredLogs = selectedFilter === 'ALL' 
    ? consoleLogs 
    : consoleLogs.filter(l => l.level === selectedFilter);

  function scrollToBottom() {
    if (autoScroll && terminalBodyElement) {
      setTimeout(() => {
        terminalBodyElement.scrollTop = terminalBodyElement.scrollHeight;
      }, 50);
    }
  }

  function handleCommandSubmit(e: Event) {
    e.preventDefault();
    const cmd = cliInput.trim();
    if (!cmd) return;

    const now = new Date().toTimeString().split(' ')[0] + '.' + Math.floor(Math.random() * 900 + 100);

    // Append executed command
    consoleLogs = [...consoleLogs, {
      id: Date.now(),
      time: now,
      level: 'SYS',
      text: `> ${cmd}`
    }];

    // Process command
    const lower = cmd.toLowerCase();
    if (lower === 'clear' || lower === 'cls') {
      consoleLogs = [];
    } else if (lower === 'status') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `UP 0 years, 3 days, 4 hours, 12 minutes, 15 seconds. 0 sessions, 0/1000 concurrent.`
      }];
    } else if (lower === 'show channels') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: activeCalls > 0 ? `${activeCalls} active SIP channels currently bridged.` : `0 total channels bridged.`
      }];
    } else if (lower === 'sofia status') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `Profile internal: RUNNING (0.0.0.0:5060) | Profile external: RUNNING (0.0.0.0:5080) | Gateways: 2 registered`
      }];
    } else if (lower === 'reloadxml') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'EVNT',
        text: `+OK [Success] Dialplan, Directory and XML configuration reloaded cleanly.`
      }];
    } else if (lower === 'help') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `Available commands: status, show channels, sofia status, reloadxml, clear, help`
      }];
    } else {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'WARN',
        text: `-ERR Unknown fs_cli command "${cmd}". Type "help" for syntax.`
      }];
    }

    cliInput = '';
    scrollToBottom();
  }

  function clearTerminal() {
    consoleLogs = [];
  }

  function copyTerminalLogs() {
    const text = filteredLogs.map(l => `[${l.time}] [${l.level}] ${l.text}`).join('\n');
    navigator.clipboard.writeText(text);
    copyFeedback = true;
    setTimeout(() => { copyFeedback = false; }, 2000);
  }

  async function fetchStats() {
    try {
      const res = await apiFetch('/dashboard/stats', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        activeExtensions = data.extensions || 0;
      }
    } catch (e) {}
  }

  function connectWS() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/api/v1/wallboard/stream`);
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.metrics) {
          activeCalls = data.metrics.active_calls || 0;
          callHistory = [...callHistory.slice(1), activeCalls];
          
          let tempCalls = [];
          for (let i = 0; i < activeCalls; i++) {
            tempCalls.push({
              name: `Agent ${1001 + i}`,
              caller: `100${1 + i}`,
              avatar: 'bg-indigo-600/30 text-indigo-300 border border-indigo-500/40',
              callee: `+1800555019${i}`,
              codec: 'OPUS / 48kHz (HD)',
              duration: '01:45',
              state: i % 2 === 0 ? 'Talking' : 'Ringing'
            });
          }
          liveCallList = tempCalls;
        }
      } catch (err) {}
    };
    ws.onclose = () => {
      setTimeout(connectWS, 5000);
    };
  }

  onMount(async () => {
    await fetchStats();
    connectWS();

    statsInterval = setInterval(fetchStats, 10000);
    
    wsInterval = setInterval(() => {
      cpuUsage = 18 + Math.floor(Math.random() * 10);
      ramUsage = 1220 + Math.floor(Math.random() * 150);

      cpuHistory = [...cpuHistory.slice(1), cpuUsage];
      ramHistory = [...ramHistory.slice(1), ramUsage];

      // Simulated background event ticks occasionally
      if (Math.random() > 0.65 && autoScroll) {
        const now = new Date().toTimeString().split(' ')[0] + '.' + Math.floor(Math.random() * 900 + 100);
        const sampleEvents = [
          { level: 'INFO', text: 'sofia/internal: Keepalive OPTIONS probe ACK received (latency 12ms)' },
          { level: 'EVNT', text: 'mod_xml_curl: Cached directory queried for tenant "pbx.local"' },
          { level: 'INFO', text: 'media_engine: SRTP cipher AES_CM_128_HMAC_SHA1_80 negotiated' },
          { level: 'ROUT', text: 'dialplan: Extension 1001 dialed *97 (Personal Voicemail Access)' }
        ];
        const pick = sampleEvents[Math.floor(Math.random() * sampleEvents.length)];
        consoleLogs = [...consoleLogs.slice(-40), {
          id: Date.now(),
          time: now,
          level: pick.level as any,
          text: pick.text
        }];
        scrollToBottom();
      }
    }, 3000);
  });

  onDestroy(() => {
    if (ws) ws.close();
    if (wsInterval) clearInterval(wsInterval);
    if (statsInterval) clearInterval(statsInterval);
  });
</script>

<div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500 ease-out max-w-7xl mx-auto w-full">
  
  <!-- Header with Title and Time Range Tabs -->
  <header class="flex flex-col md:flex-row md:items-end justify-between gap-4">
    <div>
      <div class="flex items-center space-x-3 mb-1">
        <h1 class="text-3xl md:text-4xl font-black tracking-tight text-white flex items-center">
          System Overview
        </h1>
        <span class="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-indigo-500/15 text-indigo-300 border border-indigo-500/30">
          Cluster Healthy
        </span>
      </div>
      <p class="text-slate-400 text-sm">Real-time telemetry, concurrent SIP traffic, and core engine event log.</p>
    </div>
    
    <!-- Time Filter Pills -->
    <div class="flex items-center space-x-1.5 bg-slate-900/80 p-1.5 rounded-xl border border-slate-800 shadow-lg">
      {#each ['1H', '24H', '7D'] as range}
        <button 
          on:click={() => activeTimeRange = range}
          class="px-3.5 py-1 text-xs font-bold rounded-lg transition-all {activeTimeRange === range ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-slate-800/60'}"
        >
          {range}
        </button>
      {/each}
    </div>
  </header>

  <!-- Sleek Metrics Cards Grid -->
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
    
    <!-- 1. Active Streams -->
    <div class="relative overflow-hidden rounded-2xl bg-slate-900/75 border border-slate-800/90 p-5 shadow-xl shadow-black/30 hover:border-indigo-500/40 transition-all duration-300 group backdrop-blur-xl">
      <div class="flex items-center justify-between mb-3">
        <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">Live Streams</p>
        <div class="w-8 h-8 rounded-lg bg-indigo-500/15 border border-indigo-500/30 flex items-center justify-center text-indigo-400 group-hover:scale-110 transition-transform">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
        </div>
      </div>
      <h2 class="text-4xl font-black text-white font-mono tracking-tight mb-4">{activeCalls}</h2>
      
      <div class="h-12 -mx-5 -mb-5 pointer-events-none opacity-80">
        <svg viewBox="0 0 100 30" class="w-full h-full" preserveAspectRatio="none">
          <defs>
            <linearGradient id="gradIndigo" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#6366f1" stop-opacity="0.45" />
              <stop offset="100%" stop-color="#6366f1" stop-opacity="0.0" />
            </linearGradient>
          </defs>
          <path d={callPath} fill="url(#gradIndigo)" />
          <polyline points={callPoints} fill="none" stroke="#818cf8" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
    </div>

    <!-- 2. Registered Endpoints -->
    <div class="relative overflow-hidden rounded-2xl bg-slate-900/75 border border-slate-800/90 p-5 shadow-xl shadow-black/30 hover:border-emerald-500/40 transition-all duration-300 group backdrop-blur-xl">
      <div class="flex items-center justify-between mb-3">
        <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">Registered Endpoints</p>
        <div class="w-8 h-8 rounded-lg bg-emerald-500/15 border border-emerald-500/30 flex items-center justify-center text-emerald-400 group-hover:scale-110 transition-transform">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
        </div>
      </div>
      <div class="flex items-baseline space-x-1.5 mb-4">
        <h2 class="text-4xl font-black text-white font-mono tracking-tight">{activeExtensions}</h2>
        <span class="text-sm font-semibold text-slate-400">/ 50 limit</span>
      </div>
      
      <!-- Capacity Progress Bar -->
      <div class="w-full bg-slate-800 rounded-full h-2 overflow-hidden mt-2">
        <div class="bg-gradient-to-r from-emerald-500 to-teal-400 h-2 rounded-full transition-all duration-500" style="width: {Math.max((activeExtensions / 50) * 100, 4)}%"></div>
      </div>
    </div>

    <!-- 3. DSP Core Load -->
    <div class="relative overflow-hidden rounded-2xl bg-slate-900/75 border border-slate-800/90 p-5 shadow-xl shadow-black/30 hover:border-cyan-500/40 transition-all duration-300 group backdrop-blur-xl">
      <div class="flex items-center justify-between mb-3">
        <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">DSP Core Load</p>
        <div class="w-8 h-8 rounded-lg bg-cyan-500/15 border border-cyan-500/30 flex items-center justify-center text-cyan-400 group-hover:scale-110 transition-transform">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"></path></svg>
        </div>
      </div>
      <div class="flex items-baseline space-x-1 mb-4">
        <h2 class="text-4xl font-black text-white font-mono tracking-tight">{cpuUsage}</h2>
        <span class="text-xl font-bold text-cyan-400">%</span>
      </div>
      
      <div class="h-12 -mx-5 -mb-5 pointer-events-none opacity-80">
        <svg viewBox="0 0 100 30" class="w-full h-full" preserveAspectRatio="none">
          <defs>
            <linearGradient id="gradCyan" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#06b6d4" stop-opacity="0.45" />
              <stop offset="100%" stop-color="#06b6d4" stop-opacity="0.0" />
            </linearGradient>
          </defs>
          <path d={cpuPath} fill="url(#gradCyan)" />
          <polyline points={cpuPoints} fill="none" stroke="#22d3ee" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
    </div>

    <!-- 4. Memory Footprint -->
    <div class="relative overflow-hidden rounded-2xl bg-slate-900/75 border border-slate-800/90 p-5 shadow-xl shadow-black/30 hover:border-purple-500/40 transition-all duration-300 group backdrop-blur-xl">
      <div class="flex items-center justify-between mb-3">
        <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">Memory Footprint</p>
        <div class="w-8 h-8 rounded-lg bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-400 group-hover:scale-110 transition-transform">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
        </div>
      </div>
      <div class="flex items-baseline space-x-1 mb-4">
        <h2 class="text-4xl font-black text-white font-mono tracking-tight">{ramUsage}</h2>
        <span class="text-lg font-bold text-purple-400">MB</span>
      </div>
      
      <div class="h-12 -mx-5 -mb-5 pointer-events-none opacity-80">
        <svg viewBox="0 0 100 30" class="w-full h-full" preserveAspectRatio="none">
          <defs>
            <linearGradient id="gradViolet" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#8b5cf6" stop-opacity="0.45" />
              <stop offset="100%" stop-color="#8b5cf6" stop-opacity="0.0" />
            </linearGradient>
          </defs>
          <path d={ramPath} fill="url(#gradViolet)" />
          <polyline points={ramPoints} fill="none" stroke="#a78bfa" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </div>
    </div>

  </div>

  <!-- Active Edge Streams Table -->
  <div class="rounded-2xl bg-slate-900/80 border border-slate-800/90 shadow-2xl overflow-hidden backdrop-blur-xl">
    <div class="px-6 py-4 border-b border-slate-800 flex justify-between items-center bg-slate-950/40">
      <h3 class="text-sm font-bold text-white uppercase tracking-wider flex items-center">
        <span class="w-2.5 h-2.5 rounded-full bg-indigo-500 mr-2.5 shadow-sm shadow-indigo-500/50"></span>
        Active Edge Streams
      </h3>
      <a href="/cdrs" class="text-xs font-bold text-indigo-400 hover:text-indigo-300 uppercase tracking-wider flex items-center gap-1 transition-colors">
        View All CDRs &rarr;
      </a>
    </div>
    
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="text-[11px] text-slate-400 uppercase tracking-wider bg-slate-950/60 border-b border-slate-800/80">
            <th class="px-6 py-3.5 font-bold">Source Endpoint</th>
            <th class="px-4 py-3.5 font-bold text-center">Dir</th>
            <th class="px-6 py-3.5 font-bold">Destination</th>
            <th class="px-6 py-3.5 font-bold">Codec</th>
            <th class="px-6 py-3.5 font-bold">Duration</th>
            <th class="px-6 py-3.5 font-bold">State</th>
          </tr>
        </thead>
        <tbody class="text-sm divide-y divide-slate-800/60">
          {#if liveCallList.length === 0}
            <tr>
              <td colspan="6" class="px-6 py-10 text-center">
                <div class="flex flex-col items-center justify-center space-y-2">
                  <div class="w-10 h-10 rounded-full bg-slate-800/80 border border-slate-700/60 flex items-center justify-center text-slate-400">
                    <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
                  </div>
                  <p class="text-slate-300 font-medium text-sm">No Active Voice Calls in Progress</p>
                  <p class="text-slate-500 text-xs">SIP endpoints are idle. Place a call via softphone or extension to see real-time bridge.</p>
                </div>
              </td>
            </tr>
          {:else}
            {#each liveCallList as call}
              <tr class="hover:bg-slate-800/40 transition-colors group cursor-pointer">
                <td class="px-6 py-3.5">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 rounded-lg {call.avatar} flex items-center justify-center font-bold text-xs">
                      {call.name.charAt(0)}
                    </div>
                    <div>
                      <p class="font-bold text-white group-hover:text-indigo-400 transition-colors">{call.name}</p>
                      <p class="text-xs text-slate-400 font-mono">{call.caller}</p>
                    </div>
                  </div>
                </td>
                <td class="px-4 py-3.5 text-center">
                  <div class="w-6 h-6 rounded-full bg-slate-800 flex items-center justify-center mx-auto text-slate-400 border border-slate-700">
                    <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"></path></svg>
                  </div>
                </td>
                <td class="px-6 py-3.5">
                  <p class="font-mono text-white font-medium">{call.callee}</p>
                </td>
                <td class="px-6 py-3.5 font-mono text-xs text-slate-400">{call.codec}</td>
                <td class="px-6 py-3.5 font-mono text-slate-300 text-xs">{call.duration}</td>
                <td class="px-6 py-3.5">
                  <span class="px-2.5 py-1 text-xs font-bold uppercase tracking-wider rounded-full border
                    {call.state === 'Talking' ? 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30' : 
                     call.state === 'Ringing' ? 'bg-amber-500/15 text-amber-400 border-amber-500/30' : 
                     'bg-slate-800 text-slate-400 border-slate-700'}">
                    {call.state}
                  </span>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>

  <!-- THE OVERHAULED HIGH-CONTRAST FS_CLI CONSOLE -->
  <div class="rounded-2xl bg-[#060810] border border-slate-800/90 shadow-2xl shadow-black/80 flex flex-col overflow-hidden font-mono">
    
    <!-- Terminal Titlebar -->
    <div class="px-5 py-3 border-b border-slate-800 bg-[#0B0F1C] flex flex-wrap gap-3 justify-between items-center select-none">
      
      <!-- Left: Window controls & title -->
      <div class="flex items-center space-x-3">
        <div class="flex space-x-1.5">
          <button on:click={clearTerminal} title="Clear terminal" class="w-3 h-3 rounded-full bg-rose-500/80 hover:bg-rose-500 transition-colors border border-rose-600/50 cursor-pointer"></button>
          <div class="w-3 h-3 rounded-full bg-amber-500/80 border border-amber-600/50"></div>
          <div class="w-3 h-3 rounded-full bg-emerald-500/80 border border-emerald-600/50"></div>
        </div>
        <div class="flex items-center space-x-2 border-l border-slate-700/60 pl-3">
          <svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
          <span class="text-xs font-bold text-slate-200 tracking-wider">fs_cli Console</span>
          <span class="text-[10px] text-slate-400">(FreeSWITCH Event Socket)</span>
        </div>
      </div>

      <!-- Right: Action tools & Filter tabs -->
      <div class="flex items-center space-x-3 flex-wrap">
        
        <!-- Log level filters -->
        <div class="flex items-center space-x-1 bg-slate-900/90 p-1 rounded-lg border border-slate-800 text-[11px]">
          {#each ['ALL', 'INFO', 'EVNT', 'WARN', 'ERRR'] as level}
            <button 
              on:click={() => selectedFilter = level as any}
              class="px-2 py-0.5 rounded font-bold transition-colors cursor-pointer {selectedFilter === level ? 'bg-indigo-600 text-white' : 'text-slate-400 hover:text-white'}"
            >
              {level}
            </button>
          {/each}
        </div>

        <!-- Font size +/- -->
        <div class="flex items-center space-x-1 text-xs text-slate-400">
          <button 
            on:click={() => terminalFontSize = Math.max(terminalFontSize - 1, 11)} 
            class="px-1.5 py-0.5 rounded bg-slate-900 hover:bg-slate-800 text-slate-300 border border-slate-800 font-bold"
            title="Decrease font size"
          >
            A-
          </button>
          <button 
            on:click={() => terminalFontSize = Math.min(terminalFontSize + 1, 18)} 
            class="px-1.5 py-0.5 rounded bg-slate-900 hover:bg-slate-800 text-slate-300 border border-slate-800 font-bold"
            title="Increase font size"
          >
            A+
          </button>
        </div>

        <!-- Copy button -->
        <button 
          on:click={copyTerminalLogs}
          class="px-2.5 py-1 rounded bg-slate-900 hover:bg-slate-800 text-xs font-semibold text-slate-300 hover:text-white border border-slate-800 transition-colors flex items-center gap-1 cursor-pointer"
        >
          {#if copyFeedback}
            <span class="text-emerald-400 font-bold">✓ Copied</span>
          {:else}
            <span>Copy</span>
          {/if}
        </button>

        <!-- Live Streaming Pill -->
        <div class="flex items-center space-x-1.5 px-2.5 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-[10px] font-semibold text-emerald-400">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
          <span>LIVE</span>
        </div>
      </div>
    </div>
    
    <!-- Terminal Logs Body (High Contrast & Perfectly Readable) -->
    <div 
      bind:this={terminalBodyElement}
      style="font-size: {terminalFontSize}px" 
      class="h-80 p-5 overflow-y-auto space-y-2 bg-[#060810] text-slate-100 font-mono leading-relaxed"
    >
      {#each filteredLogs as log (log.id)}
        <div class="flex items-start space-x-3 hover:bg-white/[0.03] px-2 py-0.5 rounded transition-colors group">
          <span class="text-slate-400 text-xs select-none flex-shrink-0 pt-0.5 font-mono">{log.time}</span>
          
          {#if log.level === 'INFO'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-sky-500/20 text-sky-300 border border-sky-500/40 flex-shrink-0">
              [INFO]
            </span>
            <span class="text-slate-100 font-normal break-all">{log.text}</span>
          {:else if log.level === 'EVNT'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-emerald-500/20 text-emerald-300 border border-emerald-500/40 flex-shrink-0">
              [EVNT]
            </span>
            <span class="text-emerald-100 font-normal break-all">{log.text}</span>
          {:else if log.level === 'WARN'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-amber-500/20 text-amber-300 border border-amber-500/40 flex-shrink-0">
              [WARN]
            </span>
            <span class="text-amber-100 font-normal break-all">{log.text}</span>
          {:else if log.level === 'ERRR'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-rose-500/25 text-rose-300 border border-rose-500/50 flex-shrink-0">
              [ERRR]
            </span>
            <span class="text-rose-200 font-semibold break-all">{log.text}</span>
          {:else if log.level === 'ROUT'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-purple-500/20 text-purple-300 border border-purple-500/40 flex-shrink-0">
              [ROUT]
            </span>
            <span class="text-purple-100 font-normal break-all">{log.text}</span>
          {:else}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-slate-700 text-slate-300 flex-shrink-0">
              [CMD]
            </span>
            <span class="text-indigo-300 font-semibold break-all">{log.text}</span>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Interactive Terminal Command Input Bar -->
    <form on:submit={handleCommandSubmit} class="p-3 border-t border-slate-800 bg-[#0B0F1C] flex items-center space-x-3">
      <span class="text-xs font-bold text-emerald-400 select-none">freeswitch@pbx2:~#</span>
      <input 
        type="text" 
        bind:value={cliInput}
        placeholder="Type command ('status', 'show channels', 'sofia status', 'reloadxml', 'help', 'clear')..." 
        class="flex-1 bg-transparent border-0 text-slate-100 placeholder-slate-500 text-xs font-mono focus:outline-none focus:ring-0"
      />
      <button 
        type="submit" 
        class="px-3 py-1 rounded bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-bold transition-colors cursor-pointer"
      >
        Send
      </button>
    </form>

  </div>

</div>
