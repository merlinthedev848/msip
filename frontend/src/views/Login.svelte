<script lang="ts">
  import { navigate } from 'svelte-routing';

  let email = '';
  let password = '';
  let error = '';
  let isLoading = false;

  async function handleLogin(e) {
    e.preventDefault();
    isLoading = true;
    error = '';

    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/auth/login`, {
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

<div class="min-h-screen bg-[#0B0F19] flex items-center justify-center p-4">
  <div class="max-w-md w-full bg-white  border border-slate-200 rounded-2xl shadow-2xl p-8 relative overflow-hidden">
    <!-- Decorative glowing orb -->
    <div class="absolute -top-24 -right-24 w-48 h-48 bg-indigo-600 rounded-full mix-blend-multiply filter blur-3xl opacity-20"></div>
    <div class="absolute -bottom-24 -left-24 w-48 h-48 bg-purple-600 rounded-full mix-blend-multiply filter blur-3xl opacity-20"></div>

    <div class="relative z-10">
      <div class="text-center mb-10">
        <div class="w-16 h-16 bg-gradient-to-tr from-indigo-600 to-purple-600 rounded-2xl mx-auto mb-6 flex items-center justify-center shadow-lg shadow-indigo-500/30 transform -rotate-6 hover:rotate-0 transition-transform">
          <svg class="w-8 h-8 text-slate-900" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
        </div>
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-2">MSIP PBX</h1>
        <p class="text-slate-500 text-sm">Enterprise Unified Communications</p>
      </div>

      <form on:submit={handleLogin} class="space-y-5">
        {#if error}
          <div class="bg-rose-500/10 border border-rose-500/20 rounded-lg p-3 text-sm text-rose-400 text-center animate-in fade-in slide-in-from-top-2">
            {error}
          </div>
        {/if}

        <div>
          <label for="email" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Email Address</label>
          <input id="email" type="email" bind:value={email} required placeholder="admin@pbx.local" class="w-full bg-slate-50/50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none">
        </div>
        
        <div>
          <label for="password" class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Password</label>
          <input id="password" type="password" bind:value={password} required placeholder="••••••••" class="w-full bg-slate-50/50 border border-slate-200 rounded-xl px-4 py-3 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none">
        </div>
        
        <button type="submit" disabled={isLoading} class="w-full bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-slate-900 py-3.5 rounded-xl font-bold transition-all shadow-lg shadow-indigo-600/20 disabled:opacity-50 mt-4">
          {isLoading ? 'Authenticating...' : 'Sign In'}
        </button>
      </form>
    </div>
  </div>
</div>
