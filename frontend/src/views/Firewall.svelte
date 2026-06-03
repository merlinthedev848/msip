<script lang="ts">
  import { onMount } from 'svelte';
  let blockedIps = [];
  let rules = [];

  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/firewall-rules`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        rules = (data.firewall_rules || []).map((r, i) => ({
          id: r.ID, name: r.Notes || 'Firewall Rule', action: r.Action.toUpperCase(), proto: 'ALL', port: 'ALL', src: r.IPAddress, order: i + 1
        }));
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Firewall & Fail2Ban</h1>
      <p class="text-gray-400 font-medium text-sm">Manage iptables rules, ACLs, and automatic intrusion prevention.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-gray-800 hover:bg-gray-700 text-white px-4 py-2 rounded-xl font-bold transition-colors border border-gray-700 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        Flush Bans
      </button>
      <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Add Rule
      </button>
    </div>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
    <!-- Fail2Ban Jail -->
    <div class="bg-gray-900/60 rounded-[2rem] border border-gray-800 overflow-hidden shadow-2xl flex flex-col">
      <div class="p-5 border-b border-gray-800 flex justify-between items-center bg-gray-950/30">
        <h3 class="text-sm font-bold text-white uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
          Active Jails (Banned IPs)
        </h3>
        <span class="bg-rose-500/10 text-rose-400 px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-widest border border-rose-500/20">{blockedIps.length} Banned</span>
      </div>
      <div class="flex-1 overflow-auto max-h-96">
        <table class="w-full text-left">
          <thead>
            <tr class="text-[10px] font-bold text-gray-500 uppercase tracking-[0.2em] border-b border-gray-800 sticky top-0 bg-gray-900/95 backdrop-blur">
              <th class="p-4">IP Address</th>
              <th class="p-4">Reason</th>
              <th class="p-4">Hits</th>
              <th class="p-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-800/50">
            {#each blockedIps as ip}
              <tr class="hover:bg-gray-800/30 transition-colors group">
                <td class="p-4">
                  <div class="flex items-center space-x-2">
                    <span class="w-5 h-4 bg-gray-800 rounded flex items-center justify-center text-[8px] font-bold border border-gray-700">{ip.country}</span>
                    <span class="font-mono font-bold text-rose-400">{ip.ip}</span>
                  </div>
                  <span class="text-[10px] text-gray-500 block mt-1">{ip.time}</span>
                </td>
                <td class="p-4">
                  <span class="text-sm font-bold text-gray-300">{ip.reason}</span>
                </td>
                <td class="p-4 font-mono font-bold text-gray-400">{ip.hits}</td>
                <td class="p-4 text-right">
                  <button aria-label="Unban IP" class="text-gray-500 hover:text-emerald-400 p-2 transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"></path></svg></button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- Active ACLs -->
    <div class="bg-gray-900/60 rounded-[2rem] border border-gray-800 overflow-hidden shadow-2xl flex flex-col">
      <div class="p-5 border-b border-gray-800 flex justify-between items-center bg-gray-950/30">
        <h3 class="text-sm font-bold text-white uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
          Active ACL Rules
        </h3>
      </div>
      <div class="flex-1 overflow-auto max-h-96 p-2">
        <div class="space-y-2">
          {#each rules as rule}
            <div class="p-4 bg-gray-950/50 rounded-xl border border-gray-800 hover:border-gray-700 transition-colors flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div class="w-8 h-8 rounded-lg bg-gray-900 flex items-center justify-center border border-gray-800 font-mono text-xs font-bold text-gray-500">
                  {rule.order}
                </div>
                <div>
                  <h4 class="text-sm font-bold text-white flex items-center space-x-2">
                    <span>{rule.name}</span>
                    <span class="px-2 py-0.5 rounded text-[9px] font-bold uppercase tracking-widest {rule.action === 'ALLOW' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">{rule.action}</span>
                  </h4>
                  <p class="text-xs text-gray-500 mt-1 font-mono">
                    <span class="text-indigo-400">{rule.proto}</span> {rule.port} <span class="text-gray-600">from</span> {rule.src}
                  </p>
                </div>
              </div>
              <div class="flex space-x-1">
                <button aria-label="Move Up" class="p-1.5 text-gray-600 hover:text-white transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path></svg></button>
                <button aria-label="Move Down" class="p-1.5 text-gray-600 hover:text-white transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg></button>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>

