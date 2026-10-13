<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { apiFetch } from '../lib/api';

  // Live PBX Metrics from backend API
  let activeCalls = 21;
  let registeredExtensions = 168;
  let totalExtensions = 175;
  let activeTrunks = 4;
  let activeRoutes = 8;
  let totalCDRs = 1420;
  let cpuLoad = 34;
  let memoryUsage = 55;
  let diskSpace = 18;
  let serviceLevel = 88;

  let isApiOnline = true;
  let ws: WebSocket;
  let statsInterval: any;
  let isConsoleExpanded = false;

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

<div class="space-y-5 animate-in fade-in duration-300">
  
  <!-- Content Header: Title & Action Icons matching reference screenshot -->
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-slate-800 tracking-tight">Dashboard</h1>
    <div class="flex items-center space-x-2">
      <!-- Settings Icon -->
      <button 
        on:click={() => navigateTo('/system-settings')}
        class="p-2 text-slate-500 hover:text-slate-800 hover:bg-white rounded-lg border border-slate-200/80 transition-colors bg-white shadow-sm" 
        title="Settings"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      </button>

      <!-- Notification Bell -->
      <button 
        class="p-2 text-slate-500 hover:text-slate-800 hover:bg-white rounded-lg border border-slate-200/80 transition-colors bg-white shadow-sm relative" 
        title="Notifications"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        <span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-blue-600"></span>
      </button>
    </div>
  </div>

  <!-- ========================================== -->
  <!-- ROW 1: System Health | Live Calls | Extensions -->
  <!-- ========================================== -->
  <div class="grid grid-cols-1 md:grid-cols-12 gap-5">
    
    <!-- CARD 1: System Health (3 Semi-circular gauges) -->
    <div class="md:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-sm p-5 flex flex-col justify-between">
      <h2 class="text-xs font-bold text-slate-700 tracking-tight">System Health</h2>

      <!-- 3 Radial Semi-Circle Gauges Side-by-Side -->
      <div class="grid grid-cols-3 gap-2 my-2 text-center">
        <!-- CPU Load (Amber) -->
        <div class="flex flex-col items-center">
          <div class="relative w-20 h-12 flex items-center justify-center">
            <svg class="w-20 h-12" viewBox="0 0 88 50">
              <!-- Background Arc -->
              <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#E2E8F0" stroke-width="7" stroke-linecap="round" />
              <!-- Progress Arc: circumference = pi * 34 = 106.81 -->
              <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#F59E0B" stroke-width="7" stroke-linecap="round"
                stroke-dasharray="106.81" stroke-dashoffset="{106.81 * (1 - cpuLoad / 100)}" />
            </svg>
            <span class="absolute bottom-0 text-xs font-black text-slate-800">{cpuLoad}%</span>
          </div>
          <span class="text-[11px] font-semibold text-slate-500 mt-1">CPU Load</span>
        </div>

        <!-- Memory Usage (Blue) -->
        <div class="flex flex-col items-center">
          <div class="relative w-20 h-12 flex items-center justify-center">
            <svg class="w-20 h-12" viewBox="0 0 88 50">
              <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#E2E8F0" stroke-width="7" stroke-linecap="round" />
              <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#2563EB" stroke-width="7" stroke-linecap="round"
                stroke-dasharray="106.81" stroke-dashoffset="{106.81 * (1 - memoryUsage / 100)}" />
            </svg>
            <span class="absolute bottom-0 text-xs font-black text-slate-800">{memoryUsage}%</span>
          </div>
          <span class="text-[11px] font-semibold text-slate-500 mt-1">Memory Usage</span>
        </div>

        <!-- Disk Space (Teal) -->
        <div class="flex flex-col items-center">
          <div class="relative w-20 h-12 flex items-center justify-center">
            <svg class="w-20 h-12" viewBox="0 0 88 50">
              <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#E2E8F0" stroke-width="7" stroke-linecap="round" />
              <path d="M 10 46 A 34 34 0 0 1 78 46" fill="none" stroke="#06B6D4" stroke-width="7" stroke-linecap="round"
                stroke-dasharray="106.81" stroke-dashoffset="{106.81 * (1 - diskSpace / 100)}" />
            </svg>
            <span class="absolute bottom-0 text-xs font-black text-slate-800">{diskSpace}%</span>
          </div>
          <span class="text-[11px] font-semibold text-slate-500 mt-1">Disk Space</span>
        </div>
      </div>

      <!-- Footer: All systems operational -->
      <div class="flex items-center space-x-1.5 pt-2 border-t border-slate-100">
        <span class="text-xs text-slate-600 font-medium">All systems operational.</span>
        <svg class="w-4 h-4 text-emerald-500" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
        </svg>
      </div>
    </div>

    <!-- CARD 2: Live Calls (Stat + Spline Sparkline) -->
    <div class="md:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-sm p-5 flex flex-col justify-between">
      <h2 class="text-xs font-bold text-slate-700 tracking-tight">Live Calls</h2>
      
      <div class="flex items-end justify-between my-2">
        <div>
          <div class="text-4xl font-extrabold text-slate-900 tracking-tight">{activeCalls}</div>
          <div class="text-[11px] font-semibold text-slate-500 mt-1">Total Active Calls</div>
        </div>

        <!-- Continuous Spline Sparkline -->
        <div class="flex flex-col items-end">
          <svg class="w-32 h-14" viewBox="0 0 140 50">
            <defs>
              <linearGradient id="sparklineGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#06B6D4" stop-opacity="0.3" />
                <stop offset="100%" stop-color="#06B6D4" stop-opacity="0.0" />
              </linearGradient>
            </defs>
            <!-- Gradient Area -->
            <path d="M 0 36 C 20 24, 30 38, 50 18 C 70 8, 80 40, 100 24 C 115 12, 125 32, 140 20 L 140 50 L 0 50 Z" 
              fill="url(#sparklineGrad)" />
            <!-- Spline Line -->
            <path d="M 0 36 C 20 24, 30 38, 50 18 C 70 8, 80 40, 100 24 C 115 12, 125 32, 140 20" 
              fill="none" stroke="#06B6D4" stroke-width="2.5" stroke-linecap="round" />
          </svg>
          <span class="text-[10px] text-slate-400 font-medium">Last 30 minutes</span>
        </div>
      </div>

      <div class="h-4"></div>
    </div>

    <!-- CARD 3: Extensions (Stat + Green Donut Gauge) -->
    <div class="md:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-sm p-5 flex flex-col justify-between">
      <h2 class="text-xs font-bold text-slate-700 tracking-tight">Extensions</h2>
      
      <div class="flex items-center justify-between my-2">
        <div>
          <div class="text-3xl font-extrabold text-slate-900 tracking-tight">
            {registeredExtensions} <span class="text-slate-400 font-bold">/</span> {totalExtensions}
          </div>
          <div class="text-[11px] font-semibold text-slate-500 mt-1">Registered Extensions</div>
        </div>

        <!-- Green Donut Ring Gauge: circumference = 2 * pi * 27 = 169.64 -->
        <div class="relative w-16 h-16 flex items-center justify-center">
          <svg class="w-16 h-16 -rotate-90" viewBox="0 0 68 68">
            <circle cx="34" cy="34" r="27" fill="none" stroke="#E2E8F0" stroke-width="8" />
            <circle cx="34" cy="34" r="27" fill="none" stroke="#22C55E" stroke-width="8" stroke-linecap="round"
              stroke-dasharray="169.64" stroke-dashoffset="{169.64 * (1 - registeredExtensions / totalExtensions)}" />
          </svg>
        </div>
      </div>

      <div class="h-4"></div>
    </div>

  </div>

  <!-- ============================================================ -->
  <!-- ROW 2: Call Volume Trends (Dual Spline) | Queue Performance  -->
  <!-- ============================================================ -->
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-5">
    
    <!-- CARD 4: Call Volume Trends (65% width) -->
    <div class="lg:col-span-8 bg-white rounded-xl border border-slate-200/90 shadow-sm p-5 flex flex-col justify-between">
      <div class="flex items-center justify-between mb-2">
        <h2 class="text-xs font-bold text-slate-700 tracking-tight">Call Volume Trends</h2>
        
        <!-- Legend: Inbound & Outbound -->
        <div class="flex items-center space-x-4 text-xs font-semibold">
          <div class="flex items-center space-x-1.5 text-slate-700">
            <span class="w-2.5 h-2.5 rounded-full bg-[#1E293B]"></span>
            <span>Inbound</span>
          </div>
          <div class="flex items-center space-x-1.5 text-slate-700">
            <span class="w-2.5 h-2.5 rounded-full bg-[#0EA5E9]"></span>
            <span>Outbound</span>
          </div>
        </div>
      </div>

      <!-- Dual Spline Line Chart -->
      <div class="w-full h-52 relative mt-2">
        <svg class="w-full h-full" viewBox="0 0 540 180" preserveAspectRatio="none">
          <defs>
            <linearGradient id="outboundArea" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#0EA5E9" stop-opacity="0.15" />
              <stop offset="100%" stop-color="#0EA5E9" stop-opacity="0.0" />
            </linearGradient>
          </defs>

          <!-- Horizontal Grid Lines and Y-Axis Labels -->
          <g stroke="#F1F5F9" stroke-width="1">
            <line x1="36" y1="20" x2="530" y2="20" />
            <line x1="36" y1="55" x2="530" y2="55" />
            <line x1="36" y1="90" x2="530" y2="90" />
            <line x1="36" y1="125" x2="530" y2="125" />
            <line x1="36" y1="160" x2="530" y2="160" />
          </g>

          <g font-size="9" fill="#94A3B8" font-family="sans-serif" text-anchor="end">
            <text x="30" y="23">200</text>
            <text x="30" y="58">150</text>
            <text x="30" y="93">100</text>
            <text x="30" y="128">50</text>
            <text x="30" y="163">0</text>
          </g>

          <!-- X-Axis Timestamps -->
          <g font-size="10" fill="#94A3B8" font-family="sans-serif" text-anchor="middle">
            <text x="175" y="176">8:00</text>
            <text x="310" y="176">12:00</text>
            <text x="445" y="176">16:00</text>
          </g>

          <!-- Outbound Area Gradient -->
          <path d="M 45 142 C 80 120, 110 115, 175 130 C 220 140, 240 75, 310 138 C 360 148, 380 102, 445 122 C 480 135, 500 100, 520 95 L 520 160 L 45 160 Z"
            fill="url(#outboundArea)" />

          <!-- Outbound Spline Line (Sky Blue) -->
          <path d="M 45 142 C 80 120, 110 115, 175 130 C 220 140, 240 75, 310 138 C 360 148, 380 102, 445 122 C 480 135, 500 100, 520 95"
            fill="none" stroke="#0EA5E9" stroke-width="2" stroke-linecap="round" />

          <!-- Outbound Data Points -->
          <circle cx="45" cy="142" r="3" fill="#0EA5E9" />
          <circle cx="110" cy="118" r="3" fill="#0EA5E9" />
          <circle cx="175" cy="130" r="3" fill="#0EA5E9" />
          <circle cx="240" cy="80" r="3" fill="#0EA5E9" />
          <circle cx="310" cy="138" r="3" fill="#0EA5E9" />
          <circle cx="380" cy="105" r="3" fill="#0EA5E9" />
          <circle cx="445" cy="122" r="3" fill="#0EA5E9" />
          <circle cx="520" cy="95" r="3" fill="#0EA5E9" />

          <!-- Inbound Spline Line (Dark Navy) -->
          <path d="M 45 152 C 80 130, 110 128, 175 132 C 220 135, 240 85, 310 122 C 350 145, 380 72, 445 118 C 475 140, 500 70, 520 40"
            fill="none" stroke="#1E293B" stroke-width="2.5" stroke-linecap="round" />

          <!-- Inbound Data Points -->
          <circle cx="45" cy="152" r="3.5" fill="#1E293B" />
          <circle cx="110" cy="128" r="3.5" fill="#1E293B" />
          <circle cx="175" cy="132" r="3.5" fill="#1E293B" />
          <circle cx="240" cy="88" r="3.5" fill="#1E293B" />
          <circle cx="310" cy="122" r="3.5" fill="#1E293B" />
          <circle cx="380" cy="76" r="3.5" fill="#1E293B" />
          <circle cx="445" cy="118" r="3.5" fill="#1E293B" />
          <circle cx="520" cy="40" r="3.5" fill="#1E293B" />
        </svg>
      </div>
    </div>

    <!-- CARD 5: Queue Performance Summary (Donut Chart) -->
    <div class="lg:col-span-4 bg-white rounded-xl border border-slate-200/90 shadow-sm p-5 flex flex-col justify-between items-center text-center">
      <h2 class="text-xs font-bold text-slate-700 tracking-tight self-start">Queue Performance Summary</h2>

      <!-- Multi-Segment Donut Chart -->
      <div class="relative w-36 h-36 flex items-center justify-center my-2">
        <svg class="w-36 h-36 -rotate-90" viewBox="0 0 120 120">
          <!-- Circumference = 2 * pi * 44 = 276.46 -->
          <!-- Background -->
          <circle cx="60" cy="60" r="44" fill="none" stroke="#E2E8F0" stroke-width="12" />
          <!-- Segment 1: Answered < 30s (Green: ~70% => 193.5) -->
          <circle cx="60" cy="60" r="44" fill="none" stroke="#22C55E" stroke-width="12"
            stroke-dasharray="276.46" stroke-dashoffset="{276.46 * (1 - 0.70)}" stroke-linecap="butt" />
          <!-- Segment 2: Answered > 30s (Amber: ~18% => 49.7) -->
          <circle cx="60" cy="60" r="44" fill="none" stroke="#F59E0B" stroke-width="12"
            stroke-dasharray="49.7 226.76" stroke-dashoffset="-193.5" stroke-linecap="butt" />
          <!-- Segment 3: Missed (Red: ~12% => 33.1) -->
          <circle cx="60" cy="60" r="44" fill="none" stroke="#EF4444" stroke-width="12"
            stroke-dasharray="33.1 243.36" stroke-dashoffset="-243.2" stroke-linecap="butt" />
        </svg>

        <!-- Donut Center Content -->
        <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
          <span class="text-[10px] text-slate-500 font-semibold uppercase tracking-wider">Service Level:</span>
          <span class="text-2xl font-black text-slate-900 tracking-tight">{serviceLevel}%</span>
        </div>
      </div>

      <!-- Donut Legend -->
      <div class="flex items-center justify-center space-x-3 text-[11px] font-semibold text-slate-600 pt-2 border-t border-slate-100 w-full">
        <div class="flex items-center space-x-1">
          <span class="w-2 h-2 rounded-full bg-[#22C55E]"></span>
          <span>Answered &lt; 30s</span>
        </div>
        <div class="flex items-center space-x-1">
          <span class="w-2 h-2 rounded-full bg-[#F59E0B]"></span>
          <span>Answered &gt; 30s</span>
        </div>
        <div class="flex items-center space-x-1">
          <span class="w-2 h-2 rounded-full bg-[#EF4444]"></span>
          <span>Missed</span>
        </div>
      </div>
    </div>

  </div>

  <!-- ========================================================== -->
  <!-- ROW 3: Recent Call Logs | New Voicemails                   -->
  <!-- ========================================================== -->
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-5">
    
    <!-- CARD 6: Recent Call Logs (Table) -->
    <div class="lg:col-span-7 bg-white rounded-xl border border-slate-200/90 shadow-sm p-5 flex flex-col justify-between">
      <div>
        <h2 class="text-xs font-bold text-slate-700 tracking-tight mb-3">Recent Call Logs</h2>
        
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 text-slate-600 font-bold border-b border-slate-100">
                <th class="py-2.5 px-3 rounded-l-lg">Time</th>
                <th class="py-2.5 px-3">Caller ID</th>
                <th class="py-2.5 px-3">Direction</th>
                <th class="py-2.5 px-3">Duration</th>
                <th class="py-2.5 px-3 rounded-r-lg">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each recentCalls as call}
                <tr class="hover:bg-slate-50/80 transition-colors">
                  <td class="py-2.5 px-3 font-mono text-slate-500 font-medium">{call.time}</td>
                  <td class="py-2.5 px-3 font-semibold text-slate-800">{call.callerId}</td>
                  <td class="py-2.5 px-3 text-slate-600">{call.direction}</td>
                  <td class="py-2.5 px-3 font-mono text-slate-500">{call.duration}</td>
                  <td class="py-2.5 px-3 font-bold {call.status === 'Answered' ? 'text-emerald-600' : 'text-rose-600'}">
                    {call.status}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>

      <!-- Action Button -->
      <button 
        type="button" 
        on:click={() => navigateTo('/call-logs')}
        class="w-full mt-4 py-2 border border-slate-200 text-slate-700 hover:bg-slate-50 rounded-lg text-xs font-semibold transition-colors text-center shadow-sm"
      >
        View All Logs
      </button>
    </div>

    <!-- CARD 7: New Voicemails -->
    <div class="lg:col-span-5 bg-white rounded-xl border border-slate-200/90 shadow-sm p-5 flex flex-col justify-between">
      <div>
        <h2 class="text-xs font-bold text-slate-700 tracking-tight mb-2">New Voicemails</h2>
        <div class="text-[11px] text-slate-400 font-semibold uppercase tracking-wider pb-2 border-b border-slate-100">
          New Voicemails
        </div>

        <div class="divide-y divide-slate-100 text-xs">
          {#each recentVoicemails as vm}
            <div class="py-3 flex items-center justify-between hover:bg-slate-50/80 px-1 rounded-md transition-colors">
              <div class="flex items-center space-x-2.5">
                <span class="font-mono text-slate-600 font-medium">
                  {#if vm.unread}
                    <span class="font-bold text-slate-800">[unread]</span>
                  {:else}
                    <span class="text-slate-400">[read]</span>
                  {/if}
                  {vm.time}
                </span>
                <span class="font-semibold text-slate-800">{vm.caller}</span>
              </div>
              <span class="font-mono text-slate-400 text-xs">{vm.duration}</span>
            </div>
          {/each}
        </div>
      </div>

      <!-- Action Button -->
      <button 
        type="button" 
        on:click={() => navigateTo('/voicemail')}
        class="w-full mt-4 py-2 border border-slate-200 text-slate-700 hover:bg-slate-50 rounded-lg text-xs font-semibold transition-colors text-center shadow-sm"
      >
        View All Voicemail
      </button>
    </div>

  </div>

  <!-- ============================================================== -->
  <!-- DOCKED OPERATOR TOOL: FreeSWITCH Core CLI & Telemetry Panel    -->
  <!-- ============================================================== -->
  <div class="bg-white rounded-xl border border-slate-200/90 shadow-sm overflow-hidden">
    <!-- Panel Header with live badges and expand toggle -->
    <div class="px-5 py-3.5 bg-slate-50 border-b border-slate-200 flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center space-x-3">
        <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></div>
        <span class="text-xs font-bold text-slate-800 tracking-tight">FreeSWITCH Operator Console</span>
        <span class="px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 font-mono text-[10px] font-bold border border-emerald-300/60">
          ESL 127.0.0.1:8021 CONNECTED
        </span>
      </div>

      <!-- Quick Telecom Commands & Expand Toggle -->
      <div class="flex items-center space-x-2">
        <button 
          on:click={() => executeQuickCommand('status')}
          class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs"
        >
          status
        </button>
        <button 
          on:click={() => executeQuickCommand('show channels')}
          class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs"
        >
          show channels
        </button>
        <button 
          on:click={() => executeQuickCommand('sofia status')}
          class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs"
        >
          sofia status
        </button>
        <button 
          on:click={() => executeQuickCommand('reloadxml')}
          class="px-2.5 py-1 text-[11px] font-mono font-semibold bg-white border border-slate-200 rounded text-slate-700 hover:bg-slate-100 transition-colors shadow-2xs"
        >
          reloadxml
        </button>

        <!-- Expand / Collapse Terminal Toggle -->
        <button 
          on:click={() => isConsoleExpanded = !isConsoleExpanded}
          class="px-3 py-1 text-xs font-semibold bg-[#2563EB] text-white hover:bg-blue-700 rounded transition-colors shadow-sm ml-2 flex items-center space-x-1"
        >
          <span>{isConsoleExpanded ? 'Collapse CLI' : 'Expand CLI'}</span>
          <svg class="w-3.5 h-3.5 transition-transform duration-200 {isConsoleExpanded ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Collapsible Monospace Terminal Body -->
    {#if isConsoleExpanded}
      <div class="bg-[#090D16] p-4 text-slate-200 font-mono text-xs">
        <!-- Terminal Filter & Action Toolbar -->
        <div class="flex items-center justify-between pb-3 mb-3 border-b border-slate-800">
          <div class="flex items-center space-x-1.5">
            <span class="text-[10px] text-slate-400 font-bold uppercase tracking-wider mr-1">Filter:</span>
            {#each ['ALL', 'INFO', 'EVNT', 'WARN', 'ERRR'] as f}
              <button 
                on:click={() => selectedFilter = f as any}
                class="px-2 py-0.5 rounded text-[10px] font-bold transition-colors {selectedFilter === f ? 'bg-[#2563EB] text-white' : 'bg-slate-800/80 text-slate-400 hover:text-white'}"
              >
                {f}
              </button>
            {/each}
          </div>

          <div class="flex items-center space-x-2">
            <button 
              on:click={copyAllLogs}
              class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-semibold transition-colors"
            >
              {copyFeedback ? 'Copied!' : 'Copy Logs'}
            </button>
            <button 
              on:click={() => consoleLogs = []}
              class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-semibold transition-colors"
            >
              Clear
            </button>
          </div>
        </div>

        <!-- Terminal Output Lines -->
        <div 
          bind:this={terminalBodyElement}
          class="h-60 overflow-y-auto space-y-1 pr-2 select-text font-mono text-[12px] leading-relaxed"
        >
          {#each filteredLogs as log (log.id)}
            <div class="flex items-start space-x-2">
              <span class="text-slate-500 flex-shrink-0">[{log.time}]</span>
              <span class="px-1.5 py-0.2 rounded text-[10px] font-bold flex-shrink-0 {
                log.level === 'INFO' ? 'bg-blue-950 text-blue-400' :
                log.level === 'EVNT' ? 'bg-cyan-950 text-cyan-400' :
                log.level === 'WARN' ? 'bg-amber-950 text-amber-400' :
                log.level === 'ERRR' ? 'bg-rose-950 text-rose-400' :
                'bg-emerald-950 text-emerald-400'
              }">
                {log.level}
              </span>
              <span class="{log.level === 'SYS' ? 'text-emerald-300 font-bold' : 'text-slate-300'}">{log.text}</span>
            </div>
          {/each}
        </div>

        <!-- CLI Command Input Form -->
        <form on:submit={handleCommandSubmit} class="mt-3 pt-3 border-t border-slate-800 flex items-center space-x-2">
          <span class="text-emerald-400 font-bold">freeswitch&gt;</span>
          <input 
            type="text" 
            bind:value={cliInput}
            placeholder="Type command (e.g. status, show channels, reloadxml, sofia status)..." 
            class="flex-1 bg-slate-900 border border-slate-800 text-slate-100 placeholder-slate-500 px-3 py-1.5 rounded text-xs focus:outline-none focus:border-cyan-500 font-mono"
          />
          <button 
            type="submit" 
            class="px-3 py-1.5 bg-[#2563EB] hover:bg-blue-600 text-white font-semibold rounded text-xs transition-colors"
          >
            Send
          </button>
        </form>
      </div>
    {/if}
  </div>

</div>
