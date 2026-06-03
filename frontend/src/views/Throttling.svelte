<script lang="ts">
  let rates = [
    { ip: '192.168.1.100', desc: 'Main SIP Trunk', type: 'INVITE', current: 45, max: 100, status: 'Normal' },
    { ip: '10.0.0.50', desc: 'SBC Gateway', type: 'OPTIONS', current: 180, max: 200, status: 'Warning' },
    { ip: '45.33.12.99', desc: 'Unknown Scanner', type: 'REGISTER', current: 500, max: 50, status: 'Throttled' },
    { ip: '192.168.1.105', desc: 'Exec Phone', type: 'SUBSCRIBE', current: 5, max: 20, status: 'Normal' }
  ];
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Rate Throttling</h1>
      <p class="text-gray-400 font-medium text-sm">Monitor and limit SIP packet rates per IP to prevent DoS attacks.</p>
    </div>
  </header>

  <div class="bg-gray-900/60 rounded-[2rem] border border-gray-800 overflow-hidden shadow-2xl">
    <div class="p-5 border-b border-gray-800 flex justify-between items-center bg-gray-950/30">
      <h3 class="text-sm font-bold text-white uppercase tracking-widest flex items-center">
        <svg class="w-4 h-4 mr-2 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
        Active Limits (Packets/sec)
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-gray-500 uppercase tracking-[0.2em] border-b border-gray-800">
          <th class="p-6">Source IP</th>
          <th class="p-6">SIP Method</th>
          <th class="p-6 w-1/3">Traffic Load</th>
          <th class="p-6 text-right">Status</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each rates as rate}
          <tr class="hover:bg-gray-800/30 transition-colors group">
            <td class="p-6">
              <span class="font-bold text-white block">{rate.ip}</span>
              <span class="text-xs text-gray-500">{rate.desc}</span>
            </td>
            <td class="p-6">
              <span class="px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs font-mono font-bold text-indigo-400">{rate.type}</span>
            </td>
            <td class="p-6">
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-bold text-white">{rate.current} <span class="text-gray-500">pps</span></span>
                <span class="text-xs font-bold text-gray-500">Limit: {rate.max}</span>
              </div>
              <div class="w-full h-2 bg-gray-800 rounded-full overflow-hidden relative">
                <div class="absolute h-full transition-all duration-1000 ease-out
                  {rate.status === 'Normal' ? 'bg-emerald-500' : 
                   rate.status === 'Warning' ? 'bg-amber-500' : 
                   'bg-rose-500'}" 
                  style="width: {Math.min((rate.current / rate.max) * 100, 100)}%">
                </div>
              </div>
            </td>
            <td class="p-6 text-right">
              <span class="px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-widest
                {rate.status === 'Normal' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 
                 rate.status === 'Warning' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 
                 'bg-rose-500/10 text-rose-400 border border-rose-500/20 animate-pulse'}">
                {rate.status}
              </span>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

