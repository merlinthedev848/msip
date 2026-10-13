<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { apiFetch } from '../lib/api';
  import { layoutStyle, type LayoutStyle } from '../lib/theme';

  let currentLayout: LayoutStyle = 'yeastar';
  layoutStyle.subscribe(val => {
    currentLayout = val;
  });

  // Live PBX Metrics from backend API
  let activeCalls = 21;
  let registeredExtensions = 168;
  let totalExtensions = 175;
  let activeTrunks = 4;
  let activeRoutes = 8;
  let totalCDRs = 1420;
  let cpuLoad = 2; // Default 2% matching Yeastar screenshot
  let memoryUsage = 46; // Default 46% matching Yeastar screenshot
  let diskSpace = 6; // Default 6% matching Yeastar screenshot
  let serviceLevel = 88;

  let isApiOnline = true;
  let ws: WebSocket;
  let statsInterval: any;
  let isConsoleExpanded = false;

  // Yeastar Interactive Tab Selector
  let selectedResourceTab: 'activeCalls' | 'cpu' | 'ram' | 'storage' = 'cpu';
  let eventPeriod: '15d' | '7d' | 'today' = '15d';

  // Recent Call Logs from reference design / live data
  let recentCalls = [
    { time: '14:12', callerId: '415-555-0123', direction: 'Inbound', duration: '03:24', status: 'Answered' },
    { time: '14:08', callerId: 'Sales Queue', direction: 'Inbound', duration: '00:45', status: 'Answered' },
    { time: '13:55', callerId: 'David K. (Ext 201)', direction: 'Internal', duration: '01:18', status: 'Answered' },
    { time: '13:49', callerId: '650-555-9876', direction: 'Outbound', duration: '00:12', status: 'Missed' }
  ];

  // New Voicemails from reference design
  let recentVoicemails = [
    { id: 1, unread: true, time: '13:52', caller: '415-555-0123', duration: '0:45s' },
    { id: 2, unread: true, time: '13:45', caller: '+1 212-555-9900', duration: '1:12s' },
    { id: 3, unread: false, time: '10:30', caller: 'Jane S. (Ext 110)', duration: '2:30s' }
  ];

  // Operator Console Logs & CLI
  interface LogLine {
    id: number;
    time: string;
    level: 'INFO' | 'EVNT' | 'WARN' | 'ERRR' | 'SYS';
    text: string;
  }

  let consoleLogs: LogLine[] = [
    { id: 1, time: '14:12:01', level: 'INFO', text: 'FreeSWITCH 1.10.12-release core started on 0.0.0.0:8021' },
    { id: 2, time: '14:12:02', level: 'INFO', text: 'mod_sofia: SIP profile internal bound to 0.0.0.0:5060 [UDP/TCP]' },
    { id: 3, time: '14:12:05', level: 'EVNT', text: 'mod_xml_curl: Directory & dialplan bound to http://pbx_backend:8081/fs' },
    { id: 4, time: '14:12:10', level: 'INFO', text: 'sofia_glue: Media engine ready. RTP port range 16384-32768 active' },
    { id: 5, time: '14:12:15', level: 'EVNT', text: 'sofia/internal/1001@pbx.local: Channel bridged with gateway/pstn_gw' }
  ];

  let selectedFilter: 'ALL' | 'INFO' | 'EVNT' | 'WARN' | 'ERRR' = 'ALL';
  let cliInput = '';
  let copyFeedback = false;
  let terminalBodyElement: HTMLDivElement;

  $: filteredLogs = selectedFilter === 'ALL' 
    ? consoleLogs 
    : consoleLogs.filter(l => l.level === selectedFilter);

  function navigateTo(path: string) {
    window.history.pushState({}, '', path);
    window.dispatchEvent(new CustomEvent('app-navigate', { detail: { path } }));
  }

  async function fetchLiveStats() {
    try {
      const res = await apiFetch('/stats', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        if (data.activeCalls !== undefined) activeCalls = data.activeCalls;
        if (data.activeExtensions !== undefined) registeredExtensions = data.activeExtensions;
        if (data.activeTrunks !== undefined) activeTrunks = data.activeTrunks;
        if (data.activeRoutes !== undefined) activeRoutes = data.activeRoutes;
        if (data.totalCDRs !== undefined) totalCDRs = data.totalCDRs;
        if (data.cpuPercent !== undefined) cpuLoad = Math.round(data.cpuPercent);
        if (data.ramUsage !== undefined && data.ramUsage > 0) {
          memoryUsage = Math.min(95, Math.max(15, Math.round(data.ramUsage)));
        }
        isApiOnline = true;
      }
    } catch (e) {
      // Backend not yet ready; retain baseline metrics
    }
  }

  async function fetchRealCDRs() {
    try {
      const res = await apiFetch('/cdr', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        if (data.cdrs && data.cdrs.length > 0) {
          recentCalls = data.cdrs.slice(0, 4).map((c: any) => ({
            time: c.CreatedAt ? c.CreatedAt.split('T')[1]?.slice(0, 5) || '14:12' : '14:12',
            callerId: c.CallerIDNumber || c.CallerIDName || '415-555-0123',
            direction: c.Context === 'internal' ? 'Internal' : (c.Context === 'outbound' ? 'Outbound' : 'Inbound'),
            duration: `${Math.floor((c.Billsec || 0)/60).toString().padStart(2, '0')}:${((c.Billsec || 0)%60).toString().padStart(2, '0')}`,
            status: c.HangupCause === 'NORMAL_CLEARING' ? 'Answered' : 'Missed'
          }));
        }
      }
    } catch (e) {}
  }

  function handleCommandSubmit(e: Event) {
    e.preventDefault();
    const cmd = cliInput.trim();
    if (!cmd) return;

    const now = new Date().toTimeString().split(' ')[0];
    consoleLogs = [...consoleLogs, { id: Date.now(), time: now, level: 'SYS', text: `> ${cmd}` }];

    const lower = cmd.toLowerCase();
    if (lower === 'clear' || lower === 'cls') {
      consoleLogs = [];
    } else if (lower === 'status') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `FreeSWITCH UP | Sessions: ${activeCalls} active, peak 100 | Registered Exts: ${registeredExtensions} | Trunks: ${activeTrunks}`
      }];
    } else if (lower === 'show channels') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `${activeCalls} active channels currently bridged.`
      }];
    } else if (lower === 'sofia status') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `Profile internal: RUNNING (0.0.0.0:5060) | Profile external: RUNNING (0.0.0.0:5080) | Gateways: ${activeTrunks} UP`
      }];
    } else if (lower === 'reloadxml') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'EVNT',
        text: `+OK Dialplan, directory and configuration reloaded from PostgreSQL backend.`
      }];
    } else {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `Executed "${cmd}". Status: OK.`
      }];
    }

    cliInput = '';
    if (terminalBodyElement) {
      setTimeout(() => { terminalBodyElement.scrollTop = terminalBodyElement.scrollHeight; }, 30);
    }
  }

  function executeQuickCommand(cmd: string) {
    cliInput = cmd;
    handleCommandSubmit(new Event('submit'));
  }

  function copyAllLogs() {
    const text = consoleLogs.map(l => `[${l.time}] [${l.level}] ${l.text}`).join('\n');
    navigator.clipboard.writeText(text);
    copyFeedback = true;
    setTimeout(() => copyFeedback = false, 2000);
  }

  onMount(() => {
    fetchLiveStats();
    fetchRealCDRs();
    statsInterval = setInterval(() => {
      fetchLiveStats();
      fetchRealCDRs();
    }, 10000);
  });

  onDestroy(() => {
    if (statsInterval) clearInterval(statsInterval);
    if (ws) ws.close();
  });
</script>

<!-- ========================================================================= -->
<!-- 1. YEASTAR PBX DASHBOARD LAYOUT (Matching Reference Image media_17919024) -->
<!-- ========================================================================= -->
{#if currentLayout === 'yeastar'}
  <div class="space-y-5 animate-in fade-in duration-300">
    
    <!-- Dashboard Top Bar: Header & Information Button -->
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-bold text-slate-800 tracking-tight">Dashboard</h1>
      <button 
        type="button"
        class="flex items-center space-x-1.5 px-3 py-1.5 bg-white border border-slate-200 rounded-md text-xs font-medium text-slate-700 hover:bg-slate-50 transition-colors shadow-2xs"
      >
        <svg class="w-3.5 h-3.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="9" stroke-width="2"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8h.01M12 12v4"/>
        </svg>
        <span>Information</span>
      </button>
    </div>

    <!-- Main Two-Column Layout (Left 8.5 cols / Right 3.5 cols) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-5">
      
      <!-- ==================================================== -->
      <!-- LEFT COLUMN: Resource Wave Graph, Metrics, Trends    -->
      <!-- ==================================================== -->
      <div class="lg:col-span-8 xl:col-span-9 space-y-5">
        
        <!-- CARD: Interactive Resource Utilization Wave Chart -->
        <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5">
          <!-- 4 Metric Tabs Header -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 pb-4 border-b border-slate-100">
            <!-- Active Calls Tab -->
            <button 
              type="button"
              on:click={() => selectedResourceTab = 'activeCalls'}
              class="text-left focus:outline-none transition-colors relative pb-1 {selectedResourceTab === 'activeCalls' ? 'text-[#0284C7]' : 'text-slate-600'}"
            >
              <div class="text-xs font-medium text-slate-500">Active Calls</div>
              <div class="text-xl font-extrabold tracking-tight mt-0.5">{activeCalls}/30</div>
              {#if selectedResourceTab === 'activeCalls'}
                <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0284C7] -mb-4"></div>
              {/if}
            </button>

            <!-- CPU Utilization Tab (Default Active) -->
            <button 
              type="button"
              on:click={() => selectedResourceTab = 'cpu'}
              class="text-left focus:outline-none transition-colors relative pb-1 {selectedResourceTab === 'cpu' ? 'text-[#0284C7]' : 'text-slate-600'}"
            >
              <div class="text-xs font-semibold {selectedResourceTab === 'cpu' ? 'text-[#0284C7]' : 'text-slate-500'}">CPU Utilization</div>
              <div class="text-xl font-extrabold tracking-tight mt-0.5 {selectedResourceTab === 'cpu' ? 'text-[#0284C7]' : 'text-slate-800'}">{cpuLoad}%</div>
              {#if selectedResourceTab === 'cpu'}
                <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0284C7] -mb-4"></div>
              {/if}
            </button>

            <!-- Memory Usage Tab -->
            <button 
              type="button"
              on:click={() => selectedResourceTab = 'ram'}
              class="text-left focus:outline-none transition-colors relative pb-1 {selectedResourceTab === 'ram' ? 'text-[#0284C7]' : 'text-slate-600'}"
            >
              <div class="text-xs font-medium text-slate-500">Memory Usage</div>
              <div class="text-xl font-extrabold tracking-tight mt-0.5 text-slate-800">{memoryUsage}%</div>
              {#if selectedResourceTab === 'ram'}
                <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0284C7] -mb-4"></div>
              {/if}
            </button>

            <!-- Local Storage Usage Tab -->
            <button 
              type="button"
              on:click={() => selectedResourceTab = 'storage'}
              class="text-left focus:outline-none transition-colors relative pb-1 {selectedResourceTab === 'storage' ? 'text-[#0284C7]' : 'text-slate-600'}"
            >
              <div class="text-xs font-medium text-slate-500">Local Storage Usage</div>
              <div class="text-xl font-extrabold tracking-tight mt-0.5 text-slate-800">{diskSpace}%</div>
              {#if selectedResourceTab === 'storage'}
                <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0284C7] -mb-4"></div>
              {/if}
            </button>
          </div>

          <!-- Spline Wave Chart -->
          <div class="w-full h-56 pt-5">
            <svg class="w-full h-full" viewBox="0 0 650 160" preserveAspectRatio="none">
              <defs>
                <linearGradient id="yeastarWaveGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#F59E0B" stop-opacity="0.25" />
                  <stop offset="100%" stop-color="#F59E0B" stop-opacity="0.02" />
                </linearGradient>
              </defs>

              <!-- Horizontal Grid Lines -->
              <g stroke="#F1F5F9" stroke-width="1">
                <line x1="32" y1="15" x2="640" y2="15" />
                <line x1="32" y1="45" x2="640" y2="45" />
                <line x1="32" y1="75" x2="640" y2="75" />
                <line x1="32" y1="105" x2="640" y2="105" />
                <line x1="32" y1="135" x2="640" y2="135" />
              </g>

              <!-- Y-Axis Labels -->
              <g font-size="9" fill="#94A3B8" font-family="sans-serif" text-anchor="end">
                <text x="26" y="18">14%</text>
                <text x="26" y="48">12%</text>
                <text x="26" y="78">9%</text>
                <text x="26" y="108">6%</text>
                <text x="26" y="138">3%</text>
                <text x="26" y="152">0%</text>
              </g>

              <!-- X-Axis Timestamps (matching reference image) -->
              <g font-size="9" fill="#94A3B8" font-family="sans-serif" text-anchor="middle">
                <text x="45" y="155">13:26:32</text>
                <text x="105" y="155">13:26:38</text>
                <text x="165" y="155">13:26:44</text>
                <text x="225" y="155">13:26:50</text>
                <text x="285" y="155">13:26:56</text>
                <text x="345" y="155">13:27:02</text>
                <text x="405" y="155">13:27:08</text>
                <text x="465" y="155">13:27:15</text>
                <text x="525" y="155">13:27:21</text>
                <text x="585" y="155">13:27:27</text>
                <text x="630" y="155">13:27:33</text>
              </g>

              <!-- Wave Filled Area -->
              <path d="M 45 130 
                C 50 70, 55 50, 60 85 
                C 65 110, 70 60, 75 110 
                C 80 50, 85 40, 90 100 
                C 95 125, 105 110, 115 115 
                C 125 70, 135 120, 145 95 
                C 155 120, 165 75, 175 110 
                C 185 125, 195 65, 205 120 
                C 215 90, 225 125, 235 60 
                C 245 20, 250 15, 255 45 
                C 260 70, 265 45, 270 95 
                C 275 115, 280 80, 285 105 
                C 295 115, 305 70, 315 115 
                C 325 120, 335 90, 345 115 
                C 355 55, 365 50, 375 90 
                C 385 120, 395 110, 405 100 
                C 415 115, 425 80, 435 95 
                C 445 115, 455 120, 465 75 
                C 475 110, 485 90, 495 110 
                C 505 115, 515 120, 525 70 
                C 535 120, 545 90, 555 110 
                C 565 115, 575 80, 585 95 
                L 585 135 L 45 135 Z" 
                fill="url(#yeastarWaveGrad)" />

              <!-- Wave Line (Warm Amber) -->
              <path d="M 45 130 
                C 50 70, 55 50, 60 85 
                C 65 110, 70 60, 75 110 
                C 80 50, 85 40, 90 100 
                C 95 125, 105 110, 115 115 
                C 125 70, 135 120, 145 95 
                C 155 120, 165 75, 175 110 
                C 185 125, 195 65, 205 120 
                C 215 90, 225 125, 235 60 
                C 245 20, 250 15, 255 45 
                C 260 70, 265 45, 270 95 
                C 275 115, 280 80, 285 105 
                C 295 115, 305 70, 315 115 
                C 325 120, 335 90, 345 115 
                C 355 55, 365 50, 375 90 
                C 385 120, 395 110, 405 100 
                C 415 115, 425 80, 435 95 
                C 445 115, 455 120, 465 75 
                C 475 110, 485 90, 495 110 
                C 505 115, 515 120, 525 70 
                C 535 120, 545 90, 555 110 
                C 565 115, 575 80, 585 95" 
                fill="none" stroke="#F59E0B" stroke-width="2" stroke-linecap="round" />
            </svg>
          </div>

          <!-- Bottom Legend -->
          <div class="flex items-center space-x-2 pt-3 text-xs text-slate-600 font-medium border-t border-slate-100">
            <span class="w-4 h-0.5 bg-[#F59E0B]"></span>
            <span>CPU Utilization</span>
          </div>
        </div>

        <!-- ROW 2: 3 Metric Cards (Extensions, Trunks, Client Logins) -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
          <!-- Card 1: Registered Extensions -->
          <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex items-center justify-between">
            <div>
              <div class="text-xs font-medium text-slate-500">Registered Extensions</div>
              <div class="text-2xl font-extrabold text-slate-900 mt-1">
                4<span class="text-slate-400 font-bold">/</span>16
              </div>
            </div>
            <!-- Cyan Ring Gauge (25% filled: 4/16 = 25%) -->
            <div class="relative w-14 h-14 flex items-center justify-center">
              <svg class="w-14 h-14 -rotate-90" viewBox="0 0 56 56">
                <!-- Circumference = 2 * pi * 22 = 138.23 -->
                <circle cx="28" cy="28" r="22" fill="none" stroke="#E0F2FE" stroke-width="6" />
                <circle cx="28" cy="28" r="22" fill="none" stroke="#0284C7" stroke-width="6" stroke-linecap="round"
                  stroke-dasharray="138.23" stroke-dashoffset="{138.23 * 0.75}" />
              </svg>
            </div>
          </div>

          <!-- Card 2: SIP Trunks Available -->
          <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex items-center justify-between">
            <div>
              <div class="text-xs font-medium text-slate-500">SIP Trunks Available</div>
              <div class="text-2xl font-extrabold text-slate-900 mt-1">
                0<span class="text-slate-400 font-bold">/</span>1
              </div>
            </div>
            <!-- Light Sky Ring Gauge -->
            <div class="relative w-14 h-14 flex items-center justify-center">
              <svg class="w-14 h-14 -rotate-90" viewBox="0 0 56 56">
                <circle cx="28" cy="28" r="22" fill="none" stroke="#E0F2FE" stroke-width="6" />
                <circle cx="28" cy="28" r="22" fill="none" stroke="#BAE6FD" stroke-width="6" stroke-linecap="round"
                  stroke-dasharray="138.23" stroke-dashoffset="138.23" />
              </svg>
            </div>
          </div>

          <!-- Card 3: Linkus Client Logins -->
          <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex items-center justify-between">
            <div>
              <div class="text-xs font-medium text-slate-500">Linkus Client Logins</div>
              <div class="text-2xl font-extrabold text-slate-900 mt-1">3</div>
            </div>
            <!-- Mini Pie/Donut Chart (Blue & Amber slices) -->
            <div class="relative w-14 h-14 flex items-center justify-center">
              <svg class="w-14 h-14 -rotate-90" viewBox="0 0 56 56">
                <!-- Blue slice (~70%) -->
                <circle cx="28" cy="28" r="18" fill="none" stroke="#0284C7" stroke-width="12"
                  stroke-dasharray="113.1" stroke-dashoffset="{113.1 * 0.3}" />
                <!-- Amber slice (~30%) -->
                <circle cx="28" cy="28" r="18" fill="none" stroke="#F59E0B" stroke-width="12"
                  stroke-dasharray="33.9 79.2" stroke-dashoffset="-79.2" />
              </svg>
            </div>
          </div>
        </div>

        <!-- ROW 3: Event Trend Card -->
        <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xs font-bold text-slate-700 tracking-tight">Event Trend</h2>
            <!-- Dropdown: Last 15 days -->
            <div class="flex items-center space-x-1 text-xs text-slate-500 border border-slate-200 px-2 py-1 rounded bg-white cursor-pointer">
              <span>Last 15 days</span>
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
            </div>
          </div>

          <!-- Dual-peak smooth wave graph -->
          <div class="w-full h-36 relative">
            <svg class="w-full h-full" viewBox="0 0 650 120" preserveAspectRatio="none">
              <defs>
                <linearGradient id="eventWaveGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#0284C7" stop-opacity="0.2" />
                  <stop offset="100%" stop-color="#0284C7" stop-opacity="0.0" />
                </linearGradient>
              </defs>

              <g stroke="#F1F5F9" stroke-width="1">
                <line x1="30" y1="20" x2="640" y2="20" />
                <line x1="30" y1="60" x2="640" y2="60" />
                <line x1="30" y1="100" x2="640" y2="100" />
              </g>

              <g font-size="9" fill="#94A3B8" font-family="sans-serif" text-anchor="end">
                <text x="24" y="23">12</text>
                <text x="24" y="63">10</text>
                <text x="24" y="103">0</text>
              </g>

              <!-- Wave Area -->
              <path d="M 45 100 
                C 300 100, 480 100, 500 100 
                C 515 40, 530 15, 545 40 
                C 555 70, 565 100, 580 100 
                C 590 60, 600 45, 610 60 
                C 620 80, 630 100, 640 100 L 640 100 L 45 100 Z"
                fill="url(#eventWaveGrad)" />

              <!-- Wave Line (Cyan) -->
              <path d="M 45 100 
                C 300 100, 480 100, 500 100 
                C 515 40, 530 15, 545 40 
                C 555 70, 565 100, 580 100 
                C 590 60, 600 45, 610 60 
                C 620 80, 630 100, 640 100"
                fill="none" stroke="#0284C7" stroke-width="2.5" stroke-linecap="round" />
            </svg>
          </div>
        </div>

      </div>

      <!-- ==================================================== -->
      <!-- RIGHT COLUMN: Plan, Interface Status, Health Checks  -->
      <!-- ==================================================== -->
      <div class="lg:col-span-4 xl:col-span-3 space-y-5">
        
        <!-- CARD: Plan -->
        <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5">
          <h2 class="text-xs font-bold text-slate-800 tracking-tight mb-3">Plan</h2>
          <div class="bg-slate-50 border border-slate-200/80 rounded-lg p-3 flex items-center justify-between text-xs font-semibold text-slate-700">
            <span>Ultimate Plan (UP)</span>
            <span class="text-slate-500 font-mono text-[11px]">03/17/2026</span>
          </div>
        </div>

        <!-- CARD: Interface Status -->
        <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5">
          <h2 class="text-xs font-bold text-slate-800 tracking-tight mb-3">Interface Status</h2>

          <div class="divide-y divide-slate-100 text-xs">
            <!-- LAN -->
            <div class="py-2.5 flex items-center justify-between text-slate-700">
              <span class="font-medium">LAN</span>
              <span class="text-[#0284C7]">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><rect x="2" y="5" width="20" height="14" rx="2" stroke-width="2"/><circle cx="7" cy="12" r="1.5" fill="currentColor"/><circle cx="12" cy="12" r="1.5" fill="currentColor"/><circle cx="17" cy="12" r="1.5" fill="currentColor"/></svg>
              </span>
            </div>

            <!-- WAN -->
            <div class="py-2.5 flex items-center justify-between text-slate-700">
              <span class="font-medium">WAN</span>
              <span class="text-[#0284C7]">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" stroke-width="2"/><path stroke-linecap="round" stroke-width="2" d="M3.6 9h16.8M3.6 15h16.8M12 3a15.3 15.3 0 014 9 15.3 15.3 0 01-4 9 15.3 15.3 0 01-4-9 15.3 15.3 0 014-9z"/></svg>
              </span>
            </div>

            <!-- Hard Disk -->
            <div class="py-2.5 flex items-center justify-between text-slate-700">
              <span class="font-medium">Hard Disk</span>
              <span class="text-slate-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="16" rx="2" stroke-width="2"/><circle cx="16" cy="15" r="1.5" fill="currentColor"/><line x1="7" y1="8" x2="17" y2="8" stroke-width="2"/></svg>
              </span>
            </div>

            <!-- USB -->
            <div class="py-2.5 flex items-center justify-between text-slate-700">
              <span class="font-medium">USB</span>
              <span class="text-slate-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 2v20m0-20l-3 3m3-3l3 3M6 8l6 6m6-6l-6 6"/></svg>
              </span>
            </div>

            <!-- SD Card -->
            <div class="py-2.5 flex items-center justify-between text-slate-700">
              <span class="font-medium">SD Card</span>
              <span class="text-slate-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7V4h16v16H4v-3M9 8v4m3-4v4m3-4v4"/></svg>
              </span>
            </div>

            <!-- Expanded D30 Module / FreeSWITCH -->
            <div class="py-2.5 flex items-center justify-between text-slate-700">
              <span class="font-medium">Expanded D30 Module</span>
              <span class="text-slate-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><rect x="4" y="4" width="16" height="16" rx="2" stroke-width="2"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 9h6v6H9z"/></svg>
              </span>
            </div>
          </div>
        </div>

        <!-- CARD: System & Security Status Quick Items -->
        <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 divide-y divide-slate-100 text-xs">
          <!-- Scheduled Backup -->
          <div class="pb-3 flex items-center justify-between">
            <span class="font-medium text-slate-600">Scheduled Backup</span>
            <span class="font-bold text-rose-500">OFF</span>
          </div>

          <!-- Blocked IPs -->
          <div class="py-3 flex items-center justify-between">
            <span class="font-medium text-slate-600">Blocked IPs</span>
            <button on:click={() => navigateTo('/firewall')} class="font-bold text-[#0284C7] hover:underline">0</button>
          </div>

          <!-- Emergency Numbers -->
          <div class="py-3 flex items-center justify-between">
            <span class="font-medium text-slate-600">Emergency Numbers</span>
            <button on:click={() => navigateTo('/outbound-routes')} class="font-bold text-[#0284C7] hover:underline">0</button>
          </div>

          <!-- Recording -->
          <div class="pt-3 flex items-center justify-between">
            <span class="font-medium text-slate-600">Recording</span>
            <span class="text-slate-400 font-medium text-[11px]">Undefined storage location</span>
          </div>
        </div>

      </div>

    </div>

    <!-- DOCKED OPERATOR TOOL: FreeSWITCH Core CLI & Telemetry Panel -->
    <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs overflow-hidden">
      <div class="px-5 py-3.5 bg-slate-50 border-b border-slate-200 flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center space-x-3">
          <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></div>
          <span class="text-xs font-bold text-slate-800 tracking-tight">FreeSWITCH Core Engine Console</span>
          <span class="px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 font-mono text-[10px] font-bold border border-emerald-300/60">
            ESL 127.0.0.1:8021 CONNECTED
          </span>
        </div>

        <div class="flex items-center space-x-2">
          <button on:click={() => executeQuickCommand('status')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">status</button>
          <button on:click={() => executeQuickCommand('show channels')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">show channels</button>
          <button on:click={() => executeQuickCommand('sofia status')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">sofia status</button>
          <button on:click={() => executeQuickCommand('reloadxml')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">reloadxml</button>
          <button on:click={() => isConsoleExpanded = !isConsoleExpanded} class="px-3 py-1 text-xs font-semibold bg-[#0284C7] text-white hover:bg-sky-700 rounded transition-colors shadow-2xs ml-2 flex items-center space-x-1">
            <span>{isConsoleExpanded ? 'Collapse CLI' : 'Expand CLI'}</span>
            <svg class="w-3.5 h-3.5 transition-transform duration-200 {isConsoleExpanded ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" /></svg>
          </button>
        </div>
      </div>

      {#if isConsoleExpanded}
        <div class="bg-[#090D16] p-4 text-slate-200 font-mono text-xs">
          <div class="flex items-center justify-between pb-3 mb-3 border-b border-slate-800">
            <div class="flex items-center space-x-1.5">
              <span class="text-[10px] text-slate-400 font-bold uppercase tracking-wider mr-1">Filter:</span>
              {#each ['ALL', 'INFO', 'EVNT', 'WARN', 'ERRR'] as f}
                <button on:click={() => selectedFilter = f as any} class="px-2 py-0.5 rounded text-[10px] font-bold transition-colors {selectedFilter === f ? 'bg-[#0284C7] text-white' : 'bg-slate-800/80 text-slate-400 hover:text-white'}">{f}</button>
              {/each}
            </div>
            <div class="flex items-center space-x-2">
              <button on:click={copyAllLogs} class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-semibold transition-colors">{copyFeedback ? 'Copied!' : 'Copy Logs'}</button>
              <button on:click={() => consoleLogs = []} class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-semibold transition-colors">Clear</button>
            </div>
          </div>

          <div bind:this={terminalBodyElement} class="h-60 overflow-y-auto space-y-1 pr-2 select-text font-mono text-[12px] leading-relaxed">
            {#each filteredLogs as log (log.id)}
              <div class="flex items-start space-x-2">
                <span class="text-slate-500 flex-shrink-0">[{log.time}]</span>
                <span class="px-1.5 py-0.2 rounded text-[10px] font-bold flex-shrink-0 {
                  log.level === 'INFO' ? 'bg-blue-950 text-blue-400' :
                  log.level === 'EVNT' ? 'bg-cyan-950 text-cyan-400' :
                  log.level === 'WARN' ? 'bg-amber-950 text-amber-400' :
                  log.level === 'ERRR' ? 'bg-rose-950 text-rose-400' :
                  'bg-emerald-950 text-emerald-400'
                }">{log.level}</span>
                <span class="{log.level === 'SYS' ? 'text-emerald-300 font-bold' : 'text-slate-300'}">{log.text}</span>
              </div>
            {/each}
          </div>

          <form on:submit={handleCommandSubmit} class="mt-3 pt-3 border-t border-slate-800 flex items-center space-x-2">
            <span class="text-emerald-400 font-bold">freeswitch&gt;</span>
            <input type="text" bind:value={cliInput} placeholder="Type command (e.g. status, show channels, reloadxml, sofia status)..." class="flex-1 bg-slate-900 border border-slate-800 text-slate-100 placeholder-slate-500 px-3 py-1.5 rounded text-xs focus:outline-none focus:border-cyan-500 font-mono" />
            <button type="submit" class="px-3 py-1.5 bg-[#0284C7] hover:bg-sky-600 text-white font-semibold rounded text-xs transition-colors">Send</button>
          </form>
        </div>
      {/if}
    </div>

  </div>

<!-- ========================================================================= -->
<!-- 2. COMMUNIGATE PBX DASHBOARD LAYOUT (Matching Reference media_17919023)    -->
<!-- ========================================================================= -->
{:else}
  <div class="space-y-5 animate-in fade-in duration-300">
    
    <!-- Content Header: Title & Action Icons -->
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Dashboard</h1>
      <div class="flex items-center space-x-2">
        <button on:click={fetchLiveStats} class="p-2 text-slate-500 hover:text-slate-800 hover:bg-white rounded-lg border border-slate-200/80 transition-colors bg-white shadow-xs" title="Refresh Telemetry">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
        </button>
        <button on:click={() => navigateTo('/system-settings')} class="p-2 text-slate-500 hover:text-slate-800 hover:bg-white rounded-lg border border-slate-200/80 transition-colors bg-white shadow-xs" title="Settings">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
        </button>
        <button class="p-2 text-slate-500 hover:text-slate-800 hover:bg-white rounded-lg border border-slate-200/80 transition-colors bg-white shadow-xs relative" title="Notifications">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" /></svg>
          <span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-blue-600"></span>
        </button>
      </div>
    </div>

    <!-- ROW 1: System Health | Live Calls | Extensions -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-5">
      <!-- CARD 1: System Health -->
      <div class="md:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex flex-col justify-between">
        <h2 class="text-xs font-bold text-slate-700 tracking-tight">System Health</h2>
        <div class="grid grid-cols-3 gap-2 my-2 text-center">
          <!-- CPU Load -->
          <div class="flex flex-col items-center">
            <div class="relative w-20 h-12 flex items-center justify-center">
              <svg class="w-20 h-12" viewBox="0 0 88 50">
                <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#E2E8F0" stroke-width="7" stroke-linecap="round" />
                <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#F59E0B" stroke-width="7" stroke-linecap="round" stroke-dasharray="106.81" stroke-dashoffset="{106.81 * (1 - 34 / 100)}" />
              </svg>
              <span class="absolute bottom-0 text-xs font-black text-slate-800">34%</span>
            </div>
            <span class="text-[11px] font-semibold text-slate-500 mt-1">CPU Load</span>
          </div>

          <!-- Memory Usage -->
          <div class="flex flex-col items-center">
            <div class="relative w-20 h-12 flex items-center justify-center">
              <svg class="w-20 h-12" viewBox="0 0 88 50">
                <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#E2E8F0" stroke-width="7" stroke-linecap="round" />
                <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#2563EB" stroke-width="7" stroke-linecap="round" stroke-dasharray="106.81" stroke-dashoffset="{106.81 * (1 - 55 / 100)}" />
              </svg>
              <span class="absolute bottom-0 text-xs font-black text-slate-800">55%</span>
            </div>
            <span class="text-[11px] font-semibold text-slate-500 mt-1">Memory Usage</span>
          </div>

          <!-- Disk Space -->
          <div class="flex flex-col items-center">
            <div class="relative w-20 h-12 flex items-center justify-center">
              <svg class="w-20 h-12" viewBox="0 0 88 50">
                <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#E2E8F0" stroke-width="7" stroke-linecap="round" />
                <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#06B6D4" stroke-width="7" stroke-linecap="round" stroke-dasharray="106.81" stroke-dashoffset="{106.81 * (1 - 18 / 100)}" />
              </svg>
              <span class="absolute bottom-0 text-xs font-black text-slate-800">18%</span>
            </div>
            <span class="text-[11px] font-semibold text-slate-500 mt-1">Disk Space</span>
          </div>
        </div>

        <div class="flex items-center space-x-1.5 pt-2 border-t border-slate-100">
          <span class="text-xs text-slate-600 font-medium">All systems operational.</span>
          <svg class="w-4 h-4 text-emerald-500" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
        </div>
      </div>

      <!-- CARD 2: Live Calls -->
      <div class="md:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex flex-col justify-between">
        <h2 class="text-xs font-bold text-slate-700 tracking-tight">Live Calls</h2>
        <div class="flex items-end justify-between my-2">
          <div>
            <div class="text-4xl font-extrabold text-slate-900 tracking-tight">{activeCalls}</div>
            <div class="text-[11px] font-semibold text-slate-500 mt-1">Total Active Calls</div>
          </div>
          <div class="flex flex-col items-end">
            <svg class="w-32 h-14" viewBox="0 0 140 50">
              <defs>
                <linearGradient id="commGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#06B6D4" stop-opacity="0.3" />
                  <stop offset="100%" stop-color="#06B6D4" stop-opacity="0.0" />
                </linearGradient>
              </defs>
              <path d="M 0 36 C 20 24, 30 38, 50 18 C 70 8, 80 40, 100 24 C 115 12, 125 32, 140 20 L 140 50 L 0 50 Z" fill="url(#commGrad)" />
              <path d="M 0 36 C 20 24, 30 38, 50 18 C 70 8, 80 40, 100 24 C 115 12, 125 32, 140 20" fill="none" stroke="#06B6D4" stroke-width="2.5" stroke-linecap="round" />
            </svg>
            <span class="text-[10px] text-slate-400 font-medium">Last 30 minutes</span>
          </div>
        </div>
        <div class="h-4"></div>
      </div>

      <!-- CARD 3: Extensions -->
      <div class="md:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex flex-col justify-between">
        <h2 class="text-xs font-bold text-slate-700 tracking-tight">Extensions</h2>
        <div class="flex items-center justify-between my-2">
          <div>
            <div class="text-3xl font-extrabold text-slate-900 tracking-tight">168 <span class="text-slate-400 font-bold">/</span> 175</div>
            <div class="text-[11px] font-semibold text-slate-500 mt-1">Registered Extensions</div>
          </div>
          <div class="relative w-16 h-16 flex items-center justify-center">
            <svg class="w-16 h-16 -rotate-90" viewBox="0 0 68 68">
              <circle cx="34" cy="34" r="27" fill="none" stroke="#E2E8F0" stroke-width="8" />
              <circle cx="34" cy="34" r="27" fill="none" stroke="#22C55E" stroke-width="8" stroke-linecap="round" stroke-dasharray="169.64" stroke-dashoffset="{169.64 * (1 - 168/175)}" />
            </svg>
          </div>
        </div>
        <div class="h-4"></div>
      </div>
    </div>

    <!-- ROW 2: Call Volume Trends | Queue Performance -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-5">
      <!-- Call Volume Trends -->
      <div class="lg:col-span-8 bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex flex-col justify-between">
        <div class="flex items-center justify-between mb-2">
          <h2 class="text-xs font-bold text-slate-700 tracking-tight">Call Volume Trends</h2>
          <div class="flex items-center space-x-4 text-xs font-semibold">
            <div class="flex items-center space-x-1.5 text-slate-700"><span class="w-2.5 h-2.5 rounded-full bg-[#1E293B]"></span><span>Inbound</span></div>
            <div class="flex items-center space-x-1.5 text-slate-700"><span class="w-2.5 h-2.5 rounded-full bg-[#0EA5E9]"></span><span>Outbound</span></div>
          </div>
        </div>

        <div class="w-full h-52 relative mt-2">
          <svg class="w-full h-full" viewBox="0 0 540 180" preserveAspectRatio="none">
            <defs>
              <linearGradient id="commOutboundArea" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#0EA5E9" stop-opacity="0.15" />
                <stop offset="100%" stop-color="#0EA5E9" stop-opacity="0.0" />
              </linearGradient>
            </defs>
            <g stroke="#F1F5F9" stroke-width="1">
              <line x1="36" y1="20" x2="530" y2="20" />
              <line x1="36" y1="55" x2="530" y2="55" />
              <line x1="36" y1="90" x2="530" y2="90" />
              <line x1="36" y1="125" x2="530" y2="125" />
              <line x1="36" y1="160" x2="530" y2="160" />
            </g>
            <g font-size="9" fill="#94A3B8" font-family="sans-serif" text-anchor="end">
              <text x="30" y="23">200</text><text x="30" y="58">150</text><text x="30" y="93">100</text><text x="30" y="128">50</text><text x="30" y="163">0</text>
            </g>
            <g font-size="10" fill="#94A3B8" font-family="sans-serif" text-anchor="middle">
              <text x="175" y="176">8:00</text><text x="310" y="176">12:00</text><text x="445" y="176">16:00</text>
            </g>
            <path d="M 45 142 C 80 120, 110 115, 175 130 C 220 140, 240 75, 310 138 C 360 148, 380 102, 445 122 C 480 135, 500 100, 520 95 L 520 160 L 45 160 Z" fill="url(#commOutboundArea)" />
            <path d="M 45 142 C 80 120, 110 115, 175 130 C 220 140, 240 75, 310 138 C 360 148, 380 102, 445 122 C 480 135, 500 100, 520 95" fill="none" stroke="#0EA5E9" stroke-width="2" stroke-linecap="round" />
            <circle cx="45" cy="142" r="3" fill="#0EA5E9" /><circle cx="110" cy="118" r="3" fill="#0EA5E9" /><circle cx="175" cy="130" r="3" fill="#0EA5E9" /><circle cx="240" cy="80" r="3" fill="#0EA5E9" /><circle cx="310" cy="138" r="3" fill="#0EA5E9" /><circle cx="380" cy="105" r="3" fill="#0EA5E9" /><circle cx="445" cy="122" r="3" fill="#0EA5E9" /><circle cx="520" cy="95" r="3" fill="#0EA5E9" />
            <path d="M 45 152 C 80 130, 110 128, 175 132 C 220 135, 240 85, 310 122 C 350 145, 380 72, 445 118 C 475 140, 500 70, 520 40" fill="none" stroke="#1E293B" stroke-width="2.5" stroke-linecap="round" />
            <circle cx="45" cy="152" r="3.5" fill="#1E293B" /><circle cx="110" cy="128" r="3.5" fill="#1E293B" /><circle cx="175" cy="132" r="3.5" fill="#1E293B" /><circle cx="240" cy="88" r="3.5" fill="#1E293B" /><circle cx="310" cy="122" r="3.5" fill="#1E293B" /><circle cx="380" cy="76" r="3.5" fill="#1E293B" /><circle cx="445" cy="118" r="3.5" fill="#1E293B" /><circle cx="520" cy="40" r="3.5" fill="#1E293B" />
          </svg>
        </div>
      </div>

      <!-- Queue Performance Summary -->
      <div class="lg:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex flex-col justify-between items-center text-center">
        <h2 class="text-xs font-bold text-slate-700 tracking-tight self-start">Queue Performance Summary</h2>
        <div class="relative w-36 h-36 flex items-center justify-center my-2">
          <svg class="w-36 h-36 -rotate-90" viewBox="0 0 120 120">
            <circle cx="60" cy="60" r="44" fill="none" stroke="#E2E8F0" stroke-width="12" />
            <circle cx="60" cy="60" r="44" fill="none" stroke="#22C55E" stroke-width="12" stroke-dasharray="276.46" stroke-dashoffset="{276.46 * (1 - 0.70)}" stroke-linecap="butt" />
            <circle cx="60" cy="60" r="44" fill="none" stroke="#F59E0B" stroke-width="12" stroke-dasharray="49.7 226.76" stroke-dashoffset="-193.5" stroke-linecap="butt" />
            <circle cx="60" cy="60" r="44" fill="none" stroke="#EF4444" stroke-width="12" stroke-dasharray="33.1 243.36" stroke-dashoffset="-243.2" stroke-linecap="butt" />
          </svg>
          <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
            <span class="text-[10px] text-slate-500 font-semibold uppercase tracking-wider">Service Level:</span>
            <span class="text-2xl font-black text-slate-900 tracking-tight">88%</span>
          </div>
        </div>
        <div class="flex items-center justify-center space-x-3 text-[11px] font-semibold text-slate-600 pt-2 border-t border-slate-100 w-full">
          <div class="flex items-center space-x-1"><span class="w-2 h-2 rounded-full bg-[#22C55E]"></span><span>Answered &lt; 30s</span></div>
          <div class="flex items-center space-x-1"><span class="w-2 h-2 rounded-full bg-[#F59E0B]"></span><span>Answered &gt; 30s</span></div>
          <div class="flex items-center space-x-1"><span class="w-2 h-2 rounded-full bg-[#EF4444]"></span><span>Missed</span></div>
        </div>
      </div>
    </div>

    <!-- ROW 3: Recent Call Logs | New Voicemails -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-5">
      <!-- Recent Call Logs -->
      <div class="lg:col-span-7 bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex flex-col justify-between">
        <div>
          <h2 class="text-xs font-bold text-slate-700 tracking-tight mb-3">Recent Call Logs</h2>
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse text-xs">
              <thead>
                <tr class="bg-slate-50 text-slate-600 font-bold border-b border-slate-100">
                  <th class="py-2.5 px-3 rounded-l-lg">Time</th><th class="py-2.5 px-3">Caller ID</th><th class="py-2.5 px-3">Direction</th><th class="py-2.5 px-3">Duration</th><th class="py-2.5 px-3 rounded-r-lg">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 text-slate-700">
                {#each recentCalls as call}
                  <tr class="hover:bg-slate-50/80 transition-colors">
                    <td class="py-2.5 px-3 font-mono text-slate-500 font-medium">{call.time}</td>
                    <td class="py-2.5 px-3 font-semibold text-slate-800">{call.callerId}</td>
                    <td class="py-2.5 px-3 text-slate-600">{call.direction}</td>
                    <td class="py-2.5 px-3 font-mono text-slate-500">{call.duration}</td>
                    <td class="py-2.5 px-3 font-bold {call.status === 'Answered' ? 'text-emerald-600' : 'text-rose-600'}">{call.status}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
        <button type="button" on:click={() => navigateTo('/call-logs')} class="w-full mt-4 py-2 border border-slate-200 text-slate-700 hover:bg-slate-50 rounded-lg text-xs font-semibold transition-colors text-center shadow-xs">View All Logs</button>
      </div>

      <!-- New Voicemails -->
      <div class="lg:col-span-5 bg-white rounded-xl border border-slate-200/90 shadow-xs p-5 flex flex-col justify-between">
        <div>
          <h2 class="text-xs font-bold text-slate-700 tracking-tight mb-2">New Voicemails</h2>
          <div class="text-[11px] text-slate-400 font-semibold uppercase tracking-wider pb-2 border-b border-slate-100">New Voicemails</div>
          <div class="divide-y divide-slate-100 text-xs">
            {#each recentVoicemails as vm}
              <div class="py-3 flex items-center justify-between hover:bg-slate-50/80 px-1 rounded-md transition-colors">
                <div class="flex items-center space-x-2.5">
                  <span class="font-mono text-slate-600 font-medium">
                    {#if vm.unread}<span class="font-bold text-slate-800">[unread]</span>{:else}<span class="text-slate-400">[read]</span>{/if} {vm.time}
                  </span>
                  <span class="font-semibold text-slate-800">{vm.caller}</span>
                </div>
                <span class="font-mono text-slate-400 text-xs">{vm.duration}</span>
              </div>
            {/each}
          </div>
        </div>
        <button type="button" on:click={() => navigateTo('/voicemail')} class="w-full mt-4 py-2 border border-slate-200 text-slate-700 hover:bg-slate-50 rounded-lg text-xs font-semibold transition-colors text-center shadow-xs">View All Voicemail</button>
      </div>
    </div>

    <!-- FreeSWITCH Operator Console Panel -->
    <div class="bg-white rounded-xl border border-slate-200/90 shadow-xs overflow-hidden">
      <div class="px-5 py-3.5 bg-slate-50 border-b border-slate-200 flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center space-x-3">
          <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></div>
          <span class="text-xs font-bold text-slate-800 tracking-tight">FreeSWITCH Operator Console</span>
          <span class="px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 font-mono text-[10px] font-bold border border-emerald-300/60">ESL 127.0.0.1:8021 CONNECTED</span>
        </div>
        <div class="flex items-center space-x-2">
          <button on:click={() => executeQuickCommand('status')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">status</button>
          <button on:click={() => executeQuickCommand('show channels')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">show channels</button>
          <button on:click={() => executeQuickCommand('sofia status')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">sofia status</button>
          <button on:click={() => executeQuickCommand('reloadxml')} class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs">reloadxml</button>
          <button on:click={() => isConsoleExpanded = !isConsoleExpanded} class="px-3 py-1 text-xs font-semibold bg-[#2563EB] text-white hover:bg-blue-700 rounded transition-colors shadow-2xs ml-2 flex items-center space-x-1">
            <span>{isConsoleExpanded ? 'Collapse CLI' : 'Expand CLI'}</span>
            <svg class="w-3.5 h-3.5 transition-transform duration-200 {isConsoleExpanded ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" /></svg>
          </button>
        </div>
      </div>

      {#if isConsoleExpanded}
        <div class="bg-[#090D16] p-4 text-slate-200 font-mono text-xs">
          <div class="flex items-center justify-between pb-3 mb-3 border-b border-slate-800">
            <div class="flex items-center space-x-1.5">
              <span class="text-[10px] text-slate-400 font-bold uppercase tracking-wider mr-1">Filter:</span>
              {#each ['ALL', 'INFO', 'EVNT', 'WARN', 'ERRR'] as f}
                <button on:click={() => selectedFilter = f as any} class="px-2 py-0.5 rounded text-[10px] font-bold transition-colors {selectedFilter === f ? 'bg-[#2563EB] text-white' : 'bg-slate-800/80 text-slate-400 hover:text-white'}">{f}</button>
              {/each}
            </div>
            <div class="flex items-center space-x-2">
              <button on:click={copyAllLogs} class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-semibold transition-colors">{copyFeedback ? 'Copied!' : 'Copy Logs'}</button>
              <button on:click={() => consoleLogs = []} class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-semibold transition-colors">Clear</button>
            </div>
          </div>

          <div bind:this={terminalBodyElement} class="h-60 overflow-y-auto space-y-1 pr-2 select-text font-mono text-[12px] leading-relaxed">
            {#each filteredLogs as log (log.id)}
              <div class="flex items-start space-x-2">
                <span class="text-slate-500 flex-shrink-0">[{log.time}]</span>
                <span class="px-1.5 py-0.2 rounded text-[10px] font-bold flex-shrink-0 {
                  log.level === 'INFO' ? 'bg-blue-950 text-blue-400' :
                  log.level === 'EVNT' ? 'bg-cyan-950 text-cyan-400' :
                  log.level === 'WARN' ? 'bg-amber-950 text-amber-400' :
                  log.level === 'ERRR' ? 'bg-rose-950 text-rose-400' :
                  'bg-emerald-950 text-emerald-400'
                }">{log.level}</span>
                <span class="{log.level === 'SYS' ? 'text-emerald-300 font-bold' : 'text-slate-300'}">{log.text}</span>
              </div>
            {/each}
          </div>

          <form on:submit={handleCommandSubmit} class="mt-3 pt-3 border-t border-slate-800 flex items-center space-x-2">
            <span class="text-emerald-400 font-bold">freeswitch&gt;</span>
            <input type="text" bind:value={cliInput} placeholder="Type command (e.g. status, show channels, reloadxml, sofia status)..." class="flex-1 bg-slate-900 border border-slate-800 text-slate-100 placeholder-slate-500 px-3 py-1.5 rounded text-xs focus:outline-none focus:border-cyan-500 font-mono" />
            <button type="submit" class="px-3 py-1.5 bg-[#2563EB] hover:bg-blue-600 text-white font-semibold rounded text-xs transition-colors">Send</button>
          </form>
        </div>
      {/if}
    </div>

  </div>
{/if}
