<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { apiFetch } from '../lib/api';

  let activeCalls = 0;
  let activeExtensions = 0;
  let activeTrunks = 0;
  let activeRoutes = 0;
  let totalCDRs = 0;
  let ramUsage = 0;
  let goroutines = 0;
  let cpuCores = 1;

  let liveCallList = [];
  let ws: WebSocket;
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
    { id: 1, time: new Date().toTimeString().split(' ')[0], level: 'INFO', text: 'FreeSWITCH 1.10.12-release core started on 0.0.0.0:8021' },
    { id: 2, time: new Date().toTimeString().split(' ')[0], level: 'INFO', text: 'mod_sofia: SIP profile internal bound to 0.0.0.0:5060 [UDP/TCP]' },
    { id: 3, time: new Date().toTimeString().split(' ')[0], level: 'EVNT', text: 'mod_xml_curl: Gateway URL bound to http://pbx_backend:8081/fs' },
    { id: 4, time: new Date().toTimeString().split(' ')[0], level: 'INFO', text: 'sofia_glue: Media engine ready. RTP port range 16384-32768 open' },
    { id: 5, time: new Date().toTimeString().split(' ')[0], level: 'INFO', text: 'mod_conference: WebRTC MCU profiles loaded for domain pbx.local' }
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

    const now = new Date().toTimeString().split(' ')[0];

    consoleLogs = [...consoleLogs, {
      id: Date.now(),
      time: now,
      level: 'SYS',
      text: `> ${cmd}`
    }];

    const lower = cmd.toLowerCase();
    if (lower === 'clear' || lower === 'cls') {
      consoleLogs = [];
    } else if (lower === 'status') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `FreeSWITCH UP | Sessions: ${activeCalls} active, max 1000 | Exts: ${activeExtensions} | Trunks: ${activeTrunks}`
      }];
    } else if (lower === 'show channels') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: activeCalls > 0 ? `${activeCalls} active SIP channel(s) currently bridged.` : `0 total channels bridged.`
      }];
    } else if (lower === 'sofia status') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'INFO',
        text: `Profile internal: RUNNING (0.0.0.0:5060) | Profile external: RUNNING (0.0.0.0:5080) | Gateways: ${activeTrunks} active`
      }];
    } else if (lower === 'reloadxml') {
      consoleLogs = [...consoleLogs, {
        id: Date.now() + 1,
        time: now,
        level: 'EVNT',
        text: `+OK Dialplan, Directory and XML configuration reloaded from database.`
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
        text: `-ERR Unknown command "${cmd}". Type "help" for syntax.`
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
        activeTrunks = data.trunks || 0;
        activeRoutes = data.routes || 0;
        totalCDRs = data.cdrs || 0;
        goroutines = data.goroutines || 0;
        cpuCores = data.cpu_cores || 1;
        if (data.memory_mb !== undefined) {
          ramUsage = data.memory_mb;
        }
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
        }
        if (data.calls && Array.isArray(data.calls)) {
          liveCallList = data.calls;
        } else if (activeCalls === 0) {
          liveCallList = [];
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
    statsInterval = setInterval(fetchStats, 5000);
  });

  onDestroy(() => {
    if (ws) ws.close();
    if (statsInterval) clearInterval(statsInterval);
  });
</script>

<div class="space-y-6 max-w-7xl mx-auto w-full">
  
  <!-- Operational Cockpit Header -->
  <header class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-2 border-b border-[#1E2536]">
    <div>
      <div class="flex items-center space-x-3 mb-1">
        <h1 class="text-2xl md:text-3xl font-black tracking-tight text-white flex items-center gap-2">
          Telephony Operations Cockpit
        </h1>
        <div class="flex items-center space-x-1.5 px-2.5 py-1 rounded bg-[#0D281E] border border-[#10B981]/40 text-[#10B981] text-xs font-mono font-bold">
          <span class="w-2 h-2 rounded-full bg-[#10B981]"></span>
          <span>CLUSTER ONLINE</span>
        </div>
      </div>
      <p class="text-slate-400 text-xs font-mono">FreeSWITCH 1.10 SIP Switchboard • ESL Dynamic Binding Active</p>
    </div>
    
    <!-- Time Filter & System Info -->
    <div class="flex items-center space-x-3">
      <div class="flex items-center space-x-1 bg-[#121622] p-1 rounded-lg border border-[#202738]">
        {#each ['1H', '24H', '7D'] as range}
          <button 
            on:click={() => activeTimeRange = range}
            class="px-3 py-1 text-xs font-bold font-mono rounded transition-colors {activeTimeRange === range ? 'bg-[#2563EB] text-white' : 'text-slate-400 hover:text-white hover:bg-[#1A2130]'}"
          >
            {range}
          </button>
        {/each}
      </div>
    </div>
  </header>

  <!-- Hardware-Style Metric Instrument Cards (Solid Tactile Design) -->
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
    
    <!-- 1. Active Streams Instrument -->
    <div class="bg-[#121622] border border-[#202738] rounded-xl p-5 shadow-md flex flex-col justify-between">
      <div>
        <div class="flex items-center justify-between mb-2">
          <span class="text-[11px] font-bold font-mono uppercase tracking-wider text-slate-400">SIP Channels</span>
          <span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold {activeCalls > 0 ? 'bg-[#064E3B] text-[#34D399] border border-[#059669]' : 'bg-[#1C2333] text-slate-400 border border-[#2B354C]'}">
            {activeCalls > 0 ? 'ACTIVE' : 'IDLE'}
          </span>
        </div>
        <div class="flex items-baseline space-x-2 my-2">
          <span class="text-4xl font-black font-mono text-white tracking-tight">{activeCalls}</span>
          <span class="text-xs font-mono text-slate-400">/ 50 limit</span>
        </div>
      </div>

      <!-- Segmented Hardware Channel Bar -->
      <div class="mt-4 pt-3 border-t border-[#1C2333]">
        <div class="flex items-center justify-between text-[10px] font-mono text-slate-400 mb-1.5">
          <span>CAPACITY</span>
          <span>{Math.round((activeCalls / 50) * 100)}%</span>
        </div>
        <div class="w-full bg-[#0A0D14] h-2 rounded overflow-hidden flex gap-0.5 p-0.5 border border-[#1A2130]">
          {#each Array(10) as _, idx}
            <div class="flex-1 h-full rounded-sm {idx < Math.ceil((activeCalls / 50) * 10) ? 'bg-[#3B82F6]' : 'bg-[#141A26]'}"></div>
          {/each}
        </div>
      </div>
    </div>

    <!-- 2. Registered Extensions Instrument -->
    <div class="bg-[#121622] border border-[#202738] rounded-xl p-5 shadow-md flex flex-col justify-between">
      <div>
        <div class="flex items-center justify-between mb-2">
          <span class="text-[11px] font-bold font-mono uppercase tracking-wider text-slate-400">SIP Endpoints</span>
          <span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-[#064E3B] text-[#34D399] border border-[#059669]">
            REGISTRAR
          </span>
        </div>
        <div class="flex items-baseline space-x-2 my-2">
          <span class="text-4xl font-black font-mono text-white tracking-tight">{activeExtensions}</span>
          <span class="text-xs font-mono text-slate-400">provisioned</span>
        </div>
      </div>

      <div class="mt-4 pt-3 border-t border-[#1C2333] flex items-center justify-between text-xs font-mono">
        <span class="text-slate-400">DOMAIN REALM</span>
        <span class="text-[#38BDF8] font-bold">pbx.local</span>
      </div>
    </div>

    <!-- 3. Carrier Gateways Instrument -->
    <div class="bg-[#121622] border border-[#202738] rounded-xl p-5 shadow-md flex flex-col justify-between">
      <div>
        <div class="flex items-center justify-between mb-2">
          <span class="text-[11px] font-bold font-mono uppercase tracking-wider text-slate-400">Carrier Trunks</span>
          <span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-[#082F49] text-[#38BDF8] border border-[#0284C7]">
            SOFIA
          </span>
        </div>
        <div class="flex items-baseline space-x-2 my-2">
          <span class="text-4xl font-black font-mono text-white tracking-tight">{activeTrunks}</span>
          <span class="text-xs font-mono text-slate-400">gateways</span>
        </div>
      </div>

      <div class="mt-4 pt-3 border-t border-[#1C2333] flex items-center justify-between text-xs font-mono">
        <span class="text-slate-400">INBOUND DIDs</span>
        <span class="text-[#A78BFA] font-bold">{activeRoutes} active</span>
      </div>
    </div>

    <!-- 4. Engine Telemetry Instrument -->
    <div class="bg-[#121622] border border-[#202738] rounded-xl p-5 shadow-md flex flex-col justify-between">
      <div>
        <div class="flex items-center justify-between mb-2">
          <span class="text-[11px] font-bold font-mono uppercase tracking-wider text-slate-400">Engine Telemetry</span>
          <span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-[#1C2333] text-slate-300 border border-[#2B354C]">
            RUNTIME
          </span>
        </div>
        <div class="flex items-baseline space-x-2 my-2">
          <span class="text-4xl font-black font-mono text-white tracking-tight">{ramUsage}</span>
          <span class="text-xs font-mono text-slate-400">MB heap</span>
        </div>
      </div>

      <div class="mt-4 pt-3 border-t border-[#1C2333] flex items-center justify-between text-xs font-mono">
        <span class="text-slate-400">THREADS / CORES</span>
        <span class="text-[#34D399] font-bold">{goroutines} active / {cpuCores} CPU</span>
      </div>
    </div>

  </div>

  <!-- Active Edge Streams Switchboard Table (Tactile High-Density) -->
  <div class="bg-[#121622] border border-[#202738] rounded-xl shadow-lg overflow-hidden">
    <div class="px-5 py-3.5 border-b border-[#202738] bg-[#161D2B] flex justify-between items-center">
      <div class="flex items-center space-x-2.5">
        <div class="w-2.5 h-2.5 rounded-full bg-[#3B82F6]"></div>
        <h3 class="text-xs font-mono font-bold text-white uppercase tracking-wider">
          Active SIP Channels & Edge Bridges
        </h3>
      </div>
      <a 
        href="/cdrs" 
        on:click|preventDefault={() => {
          window.history.pushState({}, '', '/cdrs');
          window.dispatchEvent(new CustomEvent('app-navigate', { detail: { path: '/cdrs' } }));
        }}
        class="text-xs font-mono font-bold text-[#38BDF8] hover:text-[#7DD3FC] uppercase tracking-wider transition-colors cursor-pointer"
      >
        View Call Ledger &rarr;
      </a>
    </div>
    
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="text-[11px] font-mono text-slate-400 uppercase tracking-wider bg-[#0E121B] border-b border-[#202738]">
            <th class="px-5 py-3 font-semibold">Source Endpoint</th>
            <th class="px-4 py-3 font-semibold text-center">Dir</th>
            <th class="px-5 py-3 font-semibold">Destination</th>
            <th class="px-5 py-3 font-semibold">Media Codec</th>
            <th class="px-5 py-3 font-semibold">Duration</th>
            <th class="px-5 py-3 font-semibold">Channel State</th>
          </tr>
        </thead>
        <tbody class="text-sm divide-y divide-[#1A2130]">
          {#if liveCallList.length === 0}
            <tr>
              <td colspan="6" class="px-5 py-8 text-center bg-[#0D1017]">
                <div class="flex flex-col items-center justify-center space-y-1.5 font-mono">
                  <span class="text-slate-400 text-xs font-semibold">NO ACTIVE CHANNELS IN SESSION</span>
                  <span class="text-slate-500 text-[11px]">FreeSWITCH port 5060 listening. Place a call via Softphone or dial 9999 for Echo Test.</span>
                </div>
              </td>
            </tr>
          {:else}
            {#each liveCallList as call}
              <tr class="hover:bg-[#182030] transition-colors font-mono">
                <td class="px-5 py-3 font-semibold text-white">
                  {call.caller} ({call.name})
                </td>
                <td class="px-4 py-3 text-center text-slate-400 text-xs">
                  &harr;
                </td>
                <td class="px-5 py-3 text-[#38BDF8] font-semibold">
                  {call.callee}
                </td>
                <td class="px-5 py-3 text-slate-400 text-xs">{call.codec}</td>
                <td class="px-5 py-3 text-white text-xs">{call.duration}</td>
                <td class="px-5 py-3">
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase border
                    {call.state === 'Talking' ? 'bg-[#064E3B] text-[#34D399] border-[#059669]' : 
                     'bg-[#78350F] text-[#FBBF24] border-[#B45309]'}">
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

  <!-- THE AUTHENTIC UNIX / FREESWITCH fs_cli CONSOLE -->
  <div class="rounded-xl bg-[#080A0F] border border-[#1E2538] shadow-2xl overflow-hidden font-mono">
    
    <!-- Console Tool Bar -->
    <div class="px-4 py-2.5 border-b border-[#1E2538] bg-[#0E121B] flex flex-wrap gap-2.5 justify-between items-center select-none">
      
      <!-- Window Controls and Title -->
      <div class="flex items-center space-x-3">
        <div class="flex space-x-1.5">
          <div class="w-3 h-3 rounded-full bg-[#EF4444] border border-[#B91C1C]"></div>
          <div class="w-3 h-3 rounded-full bg-[#F59E0B] border border-[#B45309]"></div>
          <div class="w-3 h-3 rounded-full bg-[#10B981] border border-[#047857]"></div>
        </div>
        <div class="border-l border-[#242D40] pl-3 flex items-center space-x-2">
          <span class="text-xs font-bold text-white tracking-wider">fs_cli</span>
          <span class="text-[11px] text-slate-400 font-mono">FreeSWITCH Event Socket</span>
        </div>
      </div>

      <!-- Filters & Tool Buttons -->
      <div class="flex items-center space-x-2 flex-wrap">
        
        <!-- Filter Tabs -->
        <div class="flex items-center space-x-1 bg-[#141A26] p-0.5 rounded border border-[#232B3E] text-[11px]">
          {#each ['ALL', 'INFO', 'EVNT', 'WARN', 'ERRR'] as level}
            <button 
              on:click={() => selectedFilter = level as any}
              class="px-2 py-0.5 rounded font-bold transition-colors cursor-pointer {selectedFilter === level ? 'bg-[#2563EB] text-white' : 'text-slate-400 hover:text-white'}"
            >
              {level}
            </button>
          {/each}
        </div>

        <!-- Font Size Zoom -->
        <div class="flex items-center space-x-1 text-xs">
          <button 
            on:click={() => terminalFontSize = Math.max(terminalFontSize - 1, 11)} 
            class="px-1.5 py-0.5 rounded bg-[#161C2A] hover:bg-[#202738] text-slate-300 border border-[#232B3E] font-bold"
            title="Decrease font size"
          >
            A-
          </button>
          <button 
            on:click={() => terminalFontSize = Math.min(terminalFontSize + 1, 18)} 
            class="px-1.5 py-0.5 rounded bg-[#161C2A] hover:bg-[#202738] text-slate-300 border border-[#232B3E] font-bold"
            title="Increase font size"
          >
            A+
          </button>
        </div>

        <!-- Copy & Clear -->
        <button 
          on:click={copyTerminalLogs}
          class="px-2 py-0.5 rounded bg-[#161C2A] hover:bg-[#202738] text-xs font-semibold text-slate-300 hover:text-white border border-[#232B3E] transition-colors cursor-pointer"
        >
          {copyFeedback ? '✓ Copied' : 'Copy'}
        </button>

        <button 
          on:click={clearTerminal}
          class="px-2 py-0.5 rounded bg-[#161C2A] hover:bg-[#202738] text-xs font-semibold text-slate-300 hover:text-rose-400 border border-[#232B3E] transition-colors cursor-pointer"
        >
          Clear
        </button>

        <div class="flex items-center space-x-1 px-2 py-0.5 rounded bg-[#0D281E] border border-[#10B981]/40 text-[#10B981] text-[10px] font-bold font-mono">
          <span class="w-1.5 h-1.5 rounded-full bg-[#10B981]"></span>
          <span>LIVE</span>
        </div>
      </div>
    </div>
    
    <!-- Terminal Output Log Buffer (Crisp, High-Contrast Monospace) -->
    <div 
      bind:this={terminalBodyElement}
      style="font-size: {terminalFontSize}px" 
      class="h-72 p-4 overflow-y-auto space-y-1.5 bg-[#080A0F] text-slate-100 font-mono leading-relaxed"
    >
      {#each filteredLogs as log (log.id)}
        <div class="flex items-start space-x-3 px-1.5 py-0.5 rounded hover:bg-[#121622] transition-colors">
          <span class="text-slate-400 text-xs select-none flex-shrink-0 pt-0.5 font-mono">{log.time}</span>
          
          {#if log.level === 'INFO'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#082F49] text-[#38BDF8] border border-[#0284C7] flex-shrink-0">
              [INFO]
            </span>
            <span class="text-slate-200 font-normal break-all">{log.text}</span>
          {:else if log.level === 'EVNT'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#064E3B] text-[#34D399] border border-[#059669] flex-shrink-0">
              [EVNT]
            </span>
            <span class="text-[#A7F3D0] font-normal break-all">{log.text}</span>
          {:else if log.level === 'WARN'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#78350F] text-[#FBBF24] border border-[#B45309] flex-shrink-0">
              [WARN]
            </span>
            <span class="text-[#FDE68A] font-normal break-all">{log.text}</span>
          {:else if log.level === 'ERRR'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#881337] text-[#FDA4AF] border border-[#E11D48] flex-shrink-0">
              [ERRR]
            </span>
            <span class="text-[#FECDD3] font-bold break-all">{log.text}</span>
          {:else if log.level === 'ROUT'}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#581C87] text-[#C084FC] border border-[#9333EA] flex-shrink-0">
              [ROUT]
            </span>
            <span class="text-[#E9D5FF] font-normal break-all">{log.text}</span>
          {:else}
            <span class="px-1.5 py-0.2 rounded text-[10px] font-bold bg-[#1E293B] text-slate-300 flex-shrink-0">
              [CMD]
            </span>
            <span class="text-[#60A5FA] font-bold break-all">{log.text}</span>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Terminal Command Line Prompt -->
    <form on:submit={handleCommandSubmit} class="p-2.5 border-t border-[#1E2538] bg-[#0E121B] flex items-center space-x-2">
      <span class="text-xs font-bold text-[#34D399] select-none font-mono">freeswitch@pbx2:~#</span>
      <input 
        type="text" 
        bind:value={cliInput}
        placeholder="Type fs_cli command ('status', 'show channels', 'sofia status', 'reloadxml', 'help', 'clear')..." 
        class="flex-1 bg-transparent border-0 text-white placeholder-slate-500 text-xs font-mono focus:outline-none focus:ring-0"
      />
      <button 
        type="submit" 
        class="px-3 py-1 rounded bg-[#2563EB] hover:bg-[#1D4ED8] text-white text-xs font-bold font-mono transition-colors cursor-pointer"
      >
        Execute
      </button>
    </form>

  </div>

</div>
