<script lang="ts">
  import { navigate } from 'svelte-routing';
  import { apiFetch } from '../lib/api';

  let email = '';
  let password = '';
  let error = '';
  let isLoading = false;

  async function handleLogin(e) {
    e.preventDefault();
    isLoading = true;
    error = '';

    try {
      const res = await apiFetch('/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      if (res.ok) {
        const data = await res.json();
        localStorage.setItem('pbx_token', data.token);
        localStorage.setItem('pbx_user', JSON.stringify(data.user));
        // Redirect to dashboard
        window.location.href = '/';
      } else {
        const errData = await res.json();
        error = errData.error || 'Invalid credentials';
      }
    } catch (err) {
      error = 'Failed to connect to authentication server.';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-[#080A0F] flex items-center justify-center p-4">
  <div class="max-w-md w-full bg-[#121622] border border-[#202738] rounded-xl shadow-2xl p-8 relative">
    <div>
      <div class="text-center mb-8">
        <div class="w-14 h-14 bg-[#182030] border border-[#26334A] rounded-xl mx-auto mb-5 flex items-center justify-center shadow-inner">
          <svg class="w-7 h-7 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
        </div>
        <div class="inline-flex items-center space-x-2 px-2.5 py-0.5 rounded border border-cyan-500/30 bg-cyan-950/40 text-cyan-400 font-mono text-[11px] uppercase tracking-wider mb-2">
          <span class="w-1.5 h-1.5 rounded-full bg-cyan-400 animate-pulse"></span>
          <span>Core Telephony Gateway</span>
        </div>
        <h1 class="text-2xl font-black tracking-tight text-white font-mono mt-1">MSIP PBX CONTROL</h1>
        <p class="text-slate-400 text-xs font-mono mt-1">Operator & Administration Console</p>
      </div>

      <form on:submit={handleLogin} class="space-y-5">
        {#if error}
          <div class="bg-rose-950/40 border border-rose-500/40 rounded-lg p-3 text-xs text-rose-300 font-mono text-center">
            {error}
          </div>
        {/if}

        <div>
          <label for="email" class="block text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2 font-mono">Operator ID / Email</label>
          <input id="email" type="email" bind:value={email} required placeholder="admin@pbx.local" class="w-full bg-[#0B0E17] border border-[#1E2538] rounded-lg px-4 py-3 text-white focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-all outline-none font-mono text-sm">
        </div>
        
        <div>
          <label for="password" class="block text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2 font-mono">Access Token / Password</label>
          <input id="password" type="password" bind:value={password} required placeholder="••••••••" class="w-full bg-[#0B0E17] border border-[#1E2538] rounded-lg px-4 py-3 text-white focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-all outline-none font-mono text-sm">
        </div>
        
        <button type="submit" disabled={isLoading} class="w-full bg-cyan-500 hover:bg-cyan-400 text-slate-950 font-bold py-3.5 rounded-lg transition-all disabled:opacity-50 mt-4 uppercase text-xs tracking-wider font-mono">
          {isLoading ? 'Verifying Credentials...' : 'Authenticate & Connect'}
        </button>
      </form>
    </div>
  </div>
</div>

