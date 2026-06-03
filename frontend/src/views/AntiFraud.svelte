<script lang="ts">
  let blockedDestinations = [
    { country: 'Cuba', code: '+53', callsBlocked: 405, costSaved: '$850.40' },
    { country: 'Somalia', code: '+252', callsBlocked: 112, costSaved: '$310.00' },
    { country: 'North Korea', code: '+850', callsBlocked: 0, costSaved: '$0.00' },
    { country: 'Latvia (Premium)', code: '+371', callsBlocked: 890, costSaved: '$1,200.50' }
  ];
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Toll Fraud Prevention</h1>
      <p class="text-gray-400 font-medium text-sm">Block high-cost international destinations and rate centers.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center">
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Add Rule
    </button>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <div class="bg-gray-900/40 rounded-2xl border border-gray-800 p-6 flex items-center">
      <div class="w-16 h-16 rounded-2xl bg-emerald-500/10 text-emerald-500 flex items-center justify-center mr-6 border border-emerald-500/20 shadow-[0_0_15px_rgba(16,185,129,0.2)]">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <p class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-1">Est. Savings (30d)</p>
        <h2 class="text-3xl font-black text-white font-mono">$2,360.90</h2>
      </div>
    </div>
    <div class="bg-gray-900/40 rounded-2xl border border-gray-800 p-6 flex items-center">
      <div class="w-16 h-16 rounded-2xl bg-rose-500/10 text-rose-500 flex items-center justify-center mr-6 border border-rose-500/20 shadow-[0_0_15px_rgba(244,63,94,0.2)]">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path></svg>
      </div>
      <div>
        <p class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-1">Calls Blocked</p>
        <h2 class="text-3xl font-black text-white font-mono">1,407</h2>
      </div>
    </div>
    <div class="bg-gray-900/40 rounded-2xl border border-gray-800 p-6 flex items-center">
      <div class="w-16 h-16 rounded-2xl bg-indigo-500/10 text-indigo-500 flex items-center justify-center mr-6 border border-indigo-500/20 shadow-[0_0_15px_rgba(99,102,241,0.2)]">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <p class="text-xs font-bold text-gray-500 uppercase tracking-widest mb-1">Strict Mode</p>
        <h2 class="text-xl font-black text-white">Enabled</h2>
        <p class="text-[10px] text-gray-400">Whitelisting US/CAN/UK only</p>
      </div>
    </div>
  </div>

  <div class="bg-gray-900/60 rounded-[2rem] border border-gray-800 overflow-hidden shadow-2xl">
    <div class="p-5 border-b border-gray-800 flex justify-between items-center bg-gray-950/30">
      <h3 class="text-sm font-bold text-white uppercase tracking-widest flex items-center">
        <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        Blocked High-Cost Destinations
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-gray-500 uppercase tracking-[0.2em] border-b border-gray-800">
          <th class="p-6">Destination</th>
          <th class="p-6">Dial Code</th>
          <th class="p-6 text-center">Calls Dropped</th>
          <th class="p-6 text-right">Est. Savings</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each blockedDestinations as dest}
          <tr class="hover:bg-gray-800/30 transition-colors group">
            <td class="p-6">
              <span class="font-bold text-white text-base">{dest.country}</span>
            </td>
            <td class="p-6 font-mono font-bold text-indigo-400">{dest.code}</td>
            <td class="p-6 text-center font-mono font-bold text-gray-400">{dest.callsBlocked}</td>
            <td class="p-6 text-right font-mono font-bold text-emerald-400">{dest.costSaved}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

