<script lang="ts">
  let tokens = [
    { id: 'tok_live_8f9x...', name: 'Billing Integration', scopes: ['read:cdrs', 'read:accounts'], created: '2025-11-01', lastUsed: '2 mins ago', status: 'Active' },
    { id: 'tok_live_2b4p...', name: 'Mobile App Push', scopes: ['write:messages', 'read:extensions'], created: '2026-01-15', lastUsed: 'Just now', status: 'Active' },
    { id: 'tok_test_9c1m...', name: 'Staging Environment', scopes: ['*'], created: '2026-05-20', lastUsed: 'Never', status: 'Active' },
    { id: 'tok_live_old_...', name: 'Legacy CRM Sync', scopes: ['read:calls'], created: '2023-08-10', lastUsed: '6 months ago', status: 'Revoked' }
  ];
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">API Credentials</h1>
      <p class="text-gray-400 font-medium text-sm">Manage JWTs and API keys for external integrations.</p>
    </div>
    <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center">
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path></svg>
      Generate Token
    </button>
  </header>

  <div class="bg-gray-900/60 rounded-[2rem] border border-gray-800 overflow-hidden shadow-2xl">
    <div class="p-5 border-b border-gray-800 flex justify-between items-center bg-gray-950/30">
      <h3 class="text-sm font-bold text-white uppercase tracking-widest flex items-center">
        <svg class="w-4 h-4 mr-2 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
        Active Access Tokens
      </h3>
    </div>
    
    <table class="w-full text-left">
      <thead>
        <tr class="text-[10px] font-bold text-gray-500 uppercase tracking-[0.2em] border-b border-gray-800">
          <th class="p-6">Token Name</th>
          <th class="p-6">Key Prefix</th>
          <th class="p-6">Permissions (Scopes)</th>
          <th class="p-6">Last Used</th>
          <th class="p-6 text-right">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-800/50">
        {#each tokens as token}
          <tr class="hover:bg-gray-800/30 transition-colors group {token.status === 'Revoked' ? 'opacity-50' : ''}">
            <td class="p-6">
              <span class="font-bold text-white block">{token.name}</span>
              <span class="text-xs text-gray-500">Created: {token.created}</span>
            </td>
            <td class="p-6">
              <div class="flex items-center space-x-2">
                <span class="font-mono text-sm text-gray-400">{token.id}</span>
                <button aria-label="Copy Key" class="text-gray-600 hover:text-white transition-colors" disabled={token.status==='Revoked'}><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg></button>
              </div>
            </td>
            <td class="p-6">
              <div class="flex flex-wrap gap-2">
                {#each token.scopes as scope}
                  <span class="px-2 py-0.5 bg-indigo-500/10 border border-indigo-500/20 rounded text-[10px] font-bold font-mono text-indigo-400">{scope}</span>
                {/each}
              </div>
            </td>
            <td class="p-6 text-sm text-gray-400 font-medium">
              {token.lastUsed}
            </td>
            <td class="p-6 text-right">
              {#if token.status === 'Active'}
                <button class="bg-rose-500/10 text-rose-400 hover:bg-rose-500 hover:text-white px-3 py-1 rounded text-xs font-bold transition-colors">Revoke</button>
              {:else}
                <span class="px-2 py-1 bg-gray-800 text-gray-500 rounded text-[10px] font-bold uppercase tracking-widest">Revoked</span>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

