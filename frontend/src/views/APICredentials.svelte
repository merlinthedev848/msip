<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let tokens = [];
  let isModalOpen = false;
  let newTokenName = '';

  async function fetchTokens() {
    try {
      const res = await fetch(`/api/v1/credentials`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        tokens = (data.credentials || []).map(t => ({
          id: t.TokenKey, rawId: t.ID, name: t.Name, scopes: ['*'], created: t.CreatedAt, lastUsed: 'Never', status: 'Active'
        }));
      }
    } catch (e) {}
  }

  onMount(fetchTokens);

  async function handleCreateToken(e) {
    e.preventDefault();
    try {
      const res = await fetch(`/api/v1/credentials`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ 
          Name: newTokenName, 
          TokenKey: 'pbx_live_' + (window.crypto?.randomUUID ? window.crypto.randomUUID().replace(/-/g, '') : Date.now().toString(36))
        })
      });
      if (res.ok) {
        isModalOpen = false;
        newTokenName = '';
        fetchTokens();
      }
    } catch (e) {}
  }

  async function handleDeleteToken(id) {
    if(!confirm("Revoke this token?")) return;
    try {
      const res = await fetch(`/api/v1/credentials/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) fetchTokens();
    } catch (e) {}
  }
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">API Credentials</h1>
      <p class="text-slate-500 font-medium text-sm">Manage JWTs and API keys for external integrations.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path></svg>
      Generate Token
    </button>
  </header>

  <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl">
    <div class="p-4 border-b border-[#1E2538] flex justify-between items-center bg-[#0E121B]">
      <h3 class="text-xs font-bold text-white uppercase tracking-wider flex items-center font-mono">
        <svg class="w-4 h-4 mr-2 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
        Active Access Tokens
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-[#1E2538] font-mono bg-[#0E121B]">
          <th class="p-4">Token Name</th>
          <th class="p-4">Key Identifier</th>
          <th class="p-4">Permissions (Scopes)</th>
          <th class="p-4">Last Used</th>
          <th class="p-4 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[#1E2538]">
        {#each tokens as token}
          <tr class="hover:bg-[#182030]/50 transition-colors group {token.status === 'Revoked' ? 'opacity-50' : ''}">
            <td class="p-4">
              <span class="font-bold text-white block text-sm font-mono">{token.name}</span>
              <span class="text-xs text-slate-500 font-mono">Created: {token.created}</span>
            </td>
            <td class="p-4">
              <div class="flex items-center space-x-2">
                <span class="font-mono text-xs text-cyan-400 bg-cyan-950/40 px-2 py-0.5 rounded border border-cyan-500/20">{token.id}</span>
                <button aria-label="Copy Key" class="text-slate-400 hover:text-white transition-colors" disabled={token.status==='Revoked'}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg></button>
              </div>
            </td>
            <td class="p-4">
              <div class="flex flex-wrap gap-1.5">
                {#each token.scopes as scope}
                  <span class="px-2 py-0.5 bg-[#182030] border border-[#26334A] rounded text-[10px] font-bold font-mono text-cyan-400">{scope}</span>
                {/each}
              </div>
            </td>
            <td class="p-4 text-xs text-slate-400 font-mono">
              {token.lastUsed}
            </td>
            <td class="p-4 text-right space-x-2">
              <button title="Revoke Token" aria-label="Revoke Token" class="text-rose-400 hover:text-white hover:bg-rose-600 p-1.5 rounded transition-colors" on:click={() => handleDeleteToken(token.rawId)}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<Modal bind:isOpen={isModalOpen} title="Generate API Token">
  <form on:submit={handleCreateToken} class="space-y-4">
    <div>
      <label for="token_name" class="block text-xs font-mono text-slate-400 uppercase tracking-wider mb-1.5">Token Name</label>
      <input id="token_name" type="text" bind:value={newTokenName} required placeholder="e.g. My Integration" class="w-full bg-[#0B0E17] border border-[#1E2538] text-white rounded-lg px-3 py-2 font-mono text-sm focus:border-cyan-500 outline-none" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-400 hover:text-white font-mono text-xs" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-cyan-600 hover:bg-cyan-500 text-slate-950 px-5 py-2 rounded-lg font-bold font-mono text-xs uppercase tracking-wider">Generate</button>
    </div>
  </form>
</Modal>
