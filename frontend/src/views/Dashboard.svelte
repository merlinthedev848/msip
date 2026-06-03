<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  
  // Dummy data for charts/telemetry
  const activeCalls = 24;
  const maxCapacity = 100;
  
  // Simulated chart data points
  const sparklineData = [30, 45, 25, 60, 80, 55, 40, 65, 75, 45, 35, 60, 40, 24];
  const maxPoint = Math.max(...sparklineData);
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out">
  
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-white mb-1">Live Telemetry</h1>
      <p class="text-gray-400 text-sm">Real-time system health and call density.</p>
    </div>
    <Badge status="success" pulse={true}>All Systems Operational</Badge>
  </header>

  <!-- Top Metrics -->
  <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
    <GlassCard hoverEffect={true} className="flex flex-col justify-between">
      <div class="flex justify-between items-start mb-4">
        <h3 class="text-sm font-medium text-gray-400">Active Calls</h3>
        <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
      </div>
      <div class="flex items-baseline">
        <p class="text-4xl font-light text-white">{activeCalls}</p>
        <p class="ml-2 text-sm text-gray-500">/ {maxCapacity}</p>
      </div>
    </GlassCard>

    <GlassCard hoverEffect={true} className="flex flex-col justify-between">
      <div class="flex justify-between items-start mb-4">
        <h3 class="text-sm font-medium text-gray-400">Extensions</h3>
        <svg class="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
      </div>
      <p class="text-4xl font-light text-white">412</p>
    </GlassCard>
    
    <GlassCard hoverEffect={true} className="flex flex-col justify-between">
      <div class="flex justify-between items-start mb-4">
        <h3 class="text-sm font-medium text-gray-400">Trunk Status</h3>
        <svg class="w-5 h-5 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
      </div>
      <p class="text-4xl font-light text-white">3 <span class="text-lg text-gray-500">Up</span></p>
    </GlassCard>

    <GlassCard hoverEffect={true} className="flex flex-col justify-between">
      <div class="flex justify-between items-start mb-4">
        <h3 class="text-sm font-medium text-gray-400">RAM</h3>
        <svg class="w-5 h-5 text-pink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"></path></svg>
      </div>
      <div class="flex items-baseline">
        <p class="text-4xl font-light text-white">182</p>
        <p class="ml-1 text-sm text-gray-500">MB</p>
      </div>
    </GlassCard>
  </div>

  <!-- Telemetry Canvas Chart -->
  <GlassCard className="h-80 flex flex-col relative group">
    <div class="flex justify-between items-center mb-6 z-10">
      <h3 class="font-semibold text-gray-200">Call Volume Stream (Live)</h3>
      <div class="flex space-x-2">
        <span class="w-2 h-2 rounded-full bg-indigo-500 animate-ping"></span>
      </div>
    </div>
    
    <div class="flex-1 w-full relative z-0 flex items-end opacity-80 group-hover:opacity-100 transition-opacity">
      {#each sparklineData as point, i}
        <div class="flex-1 flex flex-col justify-end items-center px-1 group/bar h-full">
          <!-- Tooltip -->
          <div class="opacity-0 group-hover/bar:opacity-100 transition-opacity absolute top-10 bg-gray-900 border border-gray-700 text-xs px-2 py-1 rounded text-gray-300 pointer-events-none transform -translate-y-2 group-hover/bar:-translate-y-4">
            {point} calls
          </div>
          <div 
            class="w-full bg-gradient-to-t from-indigo-900/50 to-indigo-500 rounded-t-sm transition-all duration-500 ease-in-out hover:to-indigo-400"
            style="height: {(point / maxPoint) * 100}%"
          ></div>
        </div>
      {/each}
    </div>
    <!-- Background Grid -->
    <div class="absolute inset-0 top-16 border-t border-gray-800/50 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px] pointer-events-none z-[-1]"></div>
  </GlassCard>

</div>
