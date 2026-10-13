<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let blockedIps = [];
  let rules = [];
  let isModalOpen = false;
  let newIp = '';
  let newAction = 'accept';
  let newNotes = '';

  async function fetchRules() {
    try {
      const res = await fetch(`/api/v1/firewall`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        rules = (data.firewall || []).map((r, i) => ({
          id: r.ID, name: r.Notes || 'Firewall Rule', action: r.Action.toUpperCase(), proto: 'ALL', port: 'ALL', src: r.IPAddress, order: i + 1
        }));
      }
    } catch (e) {}
  }

  onMount(fetchRules);

  async function handleCreateRule(e) {
    e.preventDefault();
    try {
      const res = await fetch(`/api/v1/firewall`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ IPAddress: newIp, Action: newAction, Notes: newNotes })
      });
      if (res.ok) {
        isModalOpen = false;
        newIp = ''; newNotes = '';
        fetchRules();
      }
    } catch (e) {}
  }

  async function handleDeleteRule(id) {
    if(!confirm("Delete this rule?")) return;
    try {
      const res = await fetch(`/api/v1/firewall/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchRules();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Firewall & Fail2Ban</h1>
      <p class="text-slate-500 font-medium text-sm">Manage iptables rules, ACLs, and automatic intrusion prevention.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-slate-100 hover:bg-slate-200 text-white px-4 py-2 rounded-xl font-bold transition-colors border border-slate-800 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        Flush Bans
      </button>
      <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm" on:click={() => isModalOpen = true}>
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Add Rule
      </button>
    </div>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
    <!-- Fail2Ban Jail -->
    <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl flex flex-col">
      <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
        <h3 class="text-xs font-bold text-white uppercase tracking-wider flex items-center font-mono">
          <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
          Active Jails (Banned IPs)
        </h3>
        <span class="bg-rose-950/60 text-rose-400 px-2.5 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider border border-rose-500/30 font-mono">{blockedIps.length} Banned</span>
      </div>
      <div class="flex-1 overflow-auto max-h-96">
        <table class="w-full text-left">
          <thead>
            <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] sticky top-0 bg-[#0E121B] font-mono">
              <th class="p-4">IP Address</th>
              <th class="p-4">Reason</th>
              <th class="p-4">Hits</th>
              <th class="p-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#1E2538]">
            {#if blockedIps.length === 0}
              <tr>
                <td colspan="4" class="p-8 text-center text-slate-400 bg-[#121622] font-mono text-xs">No active banned IPs.</td>
              </tr>
            {/if}
            {#each blockedIps as ip}
              <tr class="hover:bg-[#182030]/50 transition-colors group">
                <td class="p-4">
                  <div class="flex items-center space-x-2">
                    <span class="px-1.5 py-0.5 bg-[#182030] rounded text-[9px] font-bold border border-[#26334A] text-cyan-400 font-mono">{ip.country}</span>
                    <span class="font-mono font-bold text-rose-400 text-xs">{ip.ip}</span>
                  </div>
                  <span class="text-[10px] text-slate-500 block mt-1 font-mono">{ip.time}</span>
                </td>
                <td class="p-4">
                  <span class="text-xs font-mono text-slate-300">{ip.reason}</span>
                </td>
                <td class="p-4 font-mono font-bold text-slate-400 text-xs">{ip.hits}</td>
                <td class="p-4 text-right">
                  <button title="Unban" aria-label="Unban IP" class="text-slate-400 hover:text-emerald-400 p-1.5 transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"></path></svg></button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- Active ACLs -->
    <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl flex flex-col">
      <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
        <h3 class="text-xs font-bold text-white uppercase tracking-wider flex items-center font-mono">
          <svg class="w-4 h-4 mr-2 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
          Active ACL Rules
        </h3>
      </div>
      <div class="flex-1 overflow-auto max-h-96 p-3">
        <div class="space-y-2">
          {#if rules.length === 0}
            <div class="p-8 text-center text-slate-400 bg-[#121622] border border-[#202738] rounded-xl font-mono text-xs">No rules configured.</div>
          {/if}
          {#each rules as rule}
            <div class="p-3 bg-[#0B0E17] rounded-lg border border-[#1E2538] hover:border-cyan-500/30 transition-colors flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <div class="w-7 h-7 rounded bg-[#182030] flex items-center justify-center border border-[#26334A] font-mono text-xs font-bold text-cyan-400">
                  {rule.order}
                </div>
                <div>
                  <h4 class="text-xs font-bold text-white flex items-center space-x-2 font-mono">
                    <span>{rule.name}</span>
                    <span class="px-2 py-0.5 rounded text-[9px] font-bold uppercase tracking-wider font-mono {rule.action === 'ALLOW' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' : 'bg-rose-950/60 text-rose-400 border border-rose-500/30'}">{rule.action}</span>
                  </h4>
                  <p class="text-[11px] text-slate-400 mt-0.5 font-mono">
                    <span class="text-cyan-400">{rule.proto}</span> {rule.port} <span class="text-slate-500">from</span> {rule.src}
                  </p>
                </div>
              </div>
              <div class="flex space-x-1">
                <button title="Delete" aria-label="Delete" class="text-rose-400 hover:text-white hover:bg-rose-600 p-1.5 rounded transition-colors" on:click={() => handleDeleteRule(rule.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Add Firewall Rule">
  <form on:submit={handleCreateRule} class="space-y-4">
    <div>
      <label for="fw_ip" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">IP Address / CIDR</label>
      <input id="fw_ip" type="text" bind:value={newIp} required placeholder="192.168.1.0/24" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div>
      <label for="fw_action" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Action</label>
      <select id="fw_action" bind:value={newAction} class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none">
        <option value="accept">Accept</option>
        <option value="drop">Drop</option>
      </select>
    </div>
    <div>
      <label for="fw_notes" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Notes</label>
      <input id="fw_notes" type="text" bind:value={newNotes} placeholder="e.g. Office Network" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Save Rule</button>
    </div>
  </form>
</Modal>
