<script lang="ts">
  import { onMount } from 'svelte';
  let devices = [];
  let templates = ['Sales_Floor', 'Exec_Suite', 'Support_Desk', 'Warehouse_DECT', 'Default_Basic'];

  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/extensions`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        devices = (data.extensions || []).map(ext => {
          const padded = (ext.ExtensionNumber || '0').padStart(6, '0');
          const lastOctets = padded.slice(-6).match(/.{2}/g).join(':');
          return {
            mac: `00:15:65:${lastOctets}`.toUpperCase(),
            vendor: 'Yealink',
            model: 'SIP-T46S',
            template: 'Default_Basic',
            ext: ext.ExtensionNumber,
            status: ext.IsActive ? 'Provisioned' : 'Offline',
            ip: '192.168.1.1' + ext.ExtensionNumber.slice(-1)
          };
        });
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Endpoint Provisioning</h1>
      <p class="text-slate-500 font-medium text-sm">Zero-touch configuration and MAC address mapping for SIP hardware.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-slate-100 hover:bg-slate-200 text-slate-900 px-4 py-2 rounded-xl font-bold transition-colors border border-slate-200 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path></svg>
        Import CSV
      </button>
      <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Add Device
      </button>
    </div>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
    <div class="lg:col-span-2 bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl flex flex-col">
      <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
        <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m14-6h2m-2 6h2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"></path></svg>
          Mapped Devices
        </h3>
        <div class="relative">
          <svg class="w-4 h-4 text-slate-500 absolute left-3 top-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <input type="text" placeholder="Search MAC..." class="bg-white border border-slate-200 text-slate-900 text-xs rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block w-48 pl-9 p-1.5 transition-colors">
        </div>
      </div>
      
      <div class="flex-1 overflow-x-auto p-2">
        <table class="w-full text-left">
          <thead>
            <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
              <th class="p-4">MAC Address</th>
              <th class="p-4">Hardware</th>
              <th class="p-4">Template</th>
              <th class="p-4 text-center">Ext</th>
              <th class="p-4">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-800/50">
            {#each devices as dev}
              <tr class="hover:bg-slate-100/30 transition-colors group">
                <td class="p-4 font-mono font-bold text-slate-900 tracking-widest">{dev.mac}</td>
                <td class="p-4">
                  <p class="font-bold text-slate-700">{dev.vendor}</p>
                  <p class="text-[10px] text-slate-500 font-mono">{dev.model}</p>
                </td>
                <td class="p-4">
                  <span class="px-2 py-1 bg-slate-100 text-slate-700 rounded text-xs font-medium border border-slate-200">{dev.template}</span>
                </td>
                <td class="p-4 text-center font-mono font-bold text-blue-600">{dev.ext}</td>
                <td class="p-4">
                  <div class="flex items-center space-x-2">
                    <div class="w-2 h-2 rounded-full {dev.status === 'Provisioned' ? 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.8)]' : dev.status === 'Rebooting' ? 'bg-amber-500 animate-pulse' : 'bg-rose-500'}"></div>
                    <span class="text-xs font-bold {dev.status === 'Provisioned' ? 'text-emerald-400' : dev.status === 'Rebooting' ? 'text-amber-400' : 'text-rose-400'}">{dev.status}</span>
                  </div>
                  <p class="text-[10px] text-gray-600 font-mono mt-0.5">{dev.ip}</p>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- Provisioning Server Info -->
    <div class="space-y-6">
      <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-xl p-6">
        <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest mb-6 flex items-center">
          <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
          HTTP Provisioning URL
        </h3>
        
        <div class="bg-slate-50 border border-slate-200 rounded-xl p-4 flex items-center justify-between group mb-4">
          <code class="text-emerald-400 text-xs font-mono">http://{window.location.hostname}/cfg/</code>
          <button aria-label="Copy URL" class="text-gray-600 hover:text-slate-900 transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg></button>
        </div>
        
        <p class="text-xs text-slate-500 leading-relaxed">
          Configure DHCP Option 66 to point to this URL. Devices will automatically pull their XML configuration based on their MAC address during boot.
        </p>
      </div>

      <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-xl">
        <div class="p-5 border-b border-slate-200 bg-slate-50/30">
          <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
            <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            Base Templates
          </h3>
        </div>
        <div class="p-2">
          <ul class="divide-y divide-gray-800/50">
            {#each templates as tpl}
              <li class="p-3 flex justify-between items-center hover:bg-slate-100/30 rounded-lg transition-colors cursor-pointer group">
                <span class="text-sm font-medium text-slate-700 group-hover:text-slate-900 transition-colors">{tpl}</span>
                <svg class="w-4 h-4 text-gray-600 group-hover:text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
              </li>
            {/each}
          </ul>
          <button class="w-full mt-2 py-3 text-xs font-bold text-blue-600 hover:text-indigo-300 uppercase tracking-widest border border-dashed border-slate-200 hover:border-indigo-500/50 rounded-xl transition-colors">
            + Create Template
          </button>
        </div>
      </div>
    </div>
  </div>
</div>


