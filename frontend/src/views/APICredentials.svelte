<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from '../components/ui/Modal.svelte';
  let tokens = [];
  let isModalOpen = false;
  let newTokenName = '';

  async function fetchTokens() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/credentials`, {
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/credentials`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({ Name: newTokenName, TokenKey: 'key_' + Math.random().toString(36).substr(2, 9) })
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
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/credentials/${id}`, {
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
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">API Credentials</h1>
      <p class="text-slate-500 font-medium text-sm">Manage JWTs and API keys for external integrations.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center" on:click={() => isModalOpen = true}>
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path></svg>
      Generate Token
    </button>
  </header>

  <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl">
    <div class="p-5 border-b border-slate-200 flex justify-between items-center bg-slate-50/30">
      <h3 class="text-sm font-bold text-slate-900 uppercase tracking-widest flex items-center">
        <svg class="w-4 h-4 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
        Active Access Tokens
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] border-b border-slate-200">
          <th class="p-6">Token Name</th>
          <th class="p-6">Key Prefix</th>
          <th class="p-6">Permissions (Scopes)</th>
          <th class="p-6">Last Used</th>
          <th class="p-6 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each tokens as token}
          <tr class="hover:bg-slate-100/30 transition-colors group {token.status === 'Revoked' ? 'opacity-50' : ''}">
            <td class="p-6">
              <span class="font-bold text-slate-900 block">{token.name}</span>
              <span class="text-xs text-slate-500">Created: {token.created}</span>
            </td>
            <td class="p-6">
              <div class="flex items-center space-x-2">
                <span class="font-mono text-sm text-slate-500">{token.id}</span>
                <button aria-label="Copy Key" class="text-gray-600 hover:text-slate-900 transition-colors" disabled={token.status==='Revoked'}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg></button>
              </div>
            </td>
            <td class="p-6">
              <div class="flex flex-wrap gap-2">
                {#each token.scopes as scope}
                  <span class="px-2 py-0.5 bg-blue-50 border border-blue-100 rounded text-[10px] font-bold font-mono text-blue-600">{scope}</span>
                {/each}
              </div>
            </td>
            <td class="p-6 text-sm text-slate-500 font-medium">
              {token.lastUsed}
            </td>
            <td class="p-6 text-right space-x-2">
              <button title="Settings" aria-label="Settings" class="text-slate-500 hover:text-slate-900 p-2 transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
              <button title="Revoke Token" aria-label="Revoke Token" class="text-rose-500 hover:text-slate-900 hover:bg-rose-500 p-2 rounded-lg transition-colors" on:click={() => handleDeleteToken(token.rawId)}><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button>
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
      <label for="token_name" class="block text-sm font-bold text-slate-500 mb-1">Token Name</label>
      <input id="token_name" type="text" bind:value={newTokenName} required placeholder="e.g. My Integration" class="w-full bg-white border border-slate-200 rounded-xl px-4 py-2 text-slate-900 focus:ring-indigo-500 focus:border-indigo-500" />
    </div>
    <div class="pt-4 flex justify-end space-x-3">
      <button type="button" class="px-4 py-2 text-slate-500 hover:text-slate-900" on:click={() => isModalOpen = false}>Cancel</button>
      <button type="submit" class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold shadow-lg shadow-indigo-500/20">Generate</button>
    </div>
  </form>
</Modal>
