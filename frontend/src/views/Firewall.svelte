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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/firewall`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/firewall`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/firewall/${id}`, {
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
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Firewall & Fail2Ban</h1>
      <p class="text-slate-500 font-medium text-sm">Manage iptables rules, ACLs, and automatic intrusion prevention.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-slate-100 hover:bg-slate-200 text-slate-900 px-4 py-2 rounded-xl font-bold transition-colors border border-slate-200 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        Flush Bans
      </button>
      <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm" on:click={() => isModalOpen = true}>
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Add Rule
      </button>
    </div>
  </header>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
    <!-- Fail2Ban Jail -->
    <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl flex flex-col">
      <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
        <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
          Active Jails (Banned IPs)
        </h3>
        <span class="bg-rose-500/10 text-rose-400 px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-widest border border-rose-500/20">{blockedIps.length} Banned</span>
      </div>
      <div class="flex-1 overflow-auto max-h-96">
        <table class="w-full text-left">
          <thead>
            <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200 sticky top-0 bg-white/95 backdrop-blur">
              <th class="p-4">IP Address</th>
              <th class="p-4">Reason</th>
              <th class="p-4">Hits</th>
              <th class="p-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-800/50">
            {#each blockedIps as ip}
              <tr class="hover:bg-slate-100/30 transition-colors group">
                <td class="p-4">
                  <div class="flex items-center space-x-2">
                    <span class="w-5 h-4 bg-slate-100 rounded flex items-center justify-center text-[8px] font-bold border border-slate-200">{ip.country}</span>
                    <span class="font-mono font-bold text-rose-400">{ip.ip}</span>
                  </div>
                  <span class="text-[10px] text-slate-500 block mt-1">{ip.time}</span>
                </td>
                <td class="p-4">
                  <span class="text-sm font-bold text-slate-700">{ip.reason}</span>
                </td>
                <td class="p-4 font-mono font-bold text-slate-500">{ip.hits}</td>
                <td class="p-4 text-right">
                  <button aria-label="Unban IP" class="text-slate-500 hover:text-emerald-400 p-2 transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"></path></svg></button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- Active ACLs -->
    <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl flex flex-col">
      <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
        <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
          <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
          Active ACL Rules
        </h3>
      </div>
      <div class="flex-1 overflow-auto max-h-96 p-2">
        <div class="space-y-2">
          {#each rules as rule}
            <div class="p-4 bg-slate-50/50 rounded-xl border border-slate-200 hover:border-slate-200 transition-colors flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div class="w-8 h-8 rounded-lg bg-white flex items-center justify-center border border-slate-200 font-mono text-xs font-bold text-slate-500">
                  {rule.order}
                </div>
                <div>
                  <h4 class="text-sm font-bold text-slate-900 flex items-center space-x-2">
                    <span>{rule.name}</span>
                    <span class="px-2 py-0.5 rounded text-[9px] font-bold uppercase tracking-widest {rule.action === 'ALLOW' ? 'bg-emerald-50 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">{rule.action}</span>
                  </h4>
                  <p class="text-xs text-slate-500 mt-1 font-mono">
                    <span class="text-blue-600">{rule.proto}</span> {rule.port} <span class="text-gray-600">from</span> {rule.src}
                  </p>
                </div>
              </div>
              <div class="flex space-x-1">
                <button aria-label="Edit" class="text-slate-500 hover:text-slate-900 p-1 transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
                <button aria-label="Delete" class="text-rose-500 hover:text-slate-900 hover:bg-rose-500 p-1 rounded transition-colors" on:click={() => handleDeleteRule(rule.id)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
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
      <label class="block text-sm font-bold text-slate-500 mb-1">IP Address / CIDR</label>
      <input type="text" bind:value={newIp} required placeholder="192.168.1.0/24" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Action</label>
      <select bind:value={newAction} class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500">
        <option value="accept">Accept</option>
        <option value="drop">Drop</option>
      </select>
    </div>
    <div>
      <label class="block text-sm font-bold text-slate-500 mb-1">Notes</label>
      <input type="text" bind:value={newNotes} placeholder="e.g. Office Network" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Save Rule</button>
    </div>
  </form>
</Modal>

