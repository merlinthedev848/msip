<script lang="ts">
  import { api } from '../lib/api';

  let email = '';
  let password = '';
  let error = '';
  let isLoading = false;

  async function handleLogin(e: Event) {
    e.preventDefault();
    isLoading = true;
    error = '';

    try {
      const data = await api.post('/auth/login', {
        email: email.trim(),
        password: password
      });

      if (data && data.token) {
        localStorage.setItem('pbx_token', data.token);
        localStorage.setItem('pbx_user', JSON.stringify(data.user || {}));
        window.location.href = '/';
      } else {
        error = 'Invalid authentication response';
      }
    } catch (err: any) {
      error = err.message || 'Invalid email or password.';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-[#EEF2F6] flex items-center justify-center p-4">
  <div class="max-w-md w-full bg-white border border-slate-200 rounded-2xl shadow-xl p-8 relative">
    <div class="text-center mb-6">
      <div class="w-12 h-12 bg-sky-50 border border-sky-200 rounded-xl mx-auto mb-4 flex items-center justify-center text-[#0284C7]">
        <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none">
          <path d="M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12" stroke="#0284C7" stroke-width="2.5" stroke-linecap="round"/>
          <path d="M12 6V12L16 14" stroke="#0EA5E9" stroke-width="2.5" stroke-linecap="round"/>
        </svg>
      </div>
      <h1 class="text-xl font-bold text-slate-800 tracking-tight">MSIP PBX CONTROL</h1>
      <p class="text-slate-500 text-xs mt-1">Telephony Administration & Operator Portal</p>
    </div>

    <form on:submit={handleLogin} class="space-y-4">
      {#if error}
        <div class="bg-rose-50 border border-rose-200 rounded-md p-3 text-xs text-rose-700 text-center font-medium">
          {error}
        </div>
      {/if}

      <div>
        <label for="login-email" class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1.5">Operator Email / Account</label>
        <input id="login-email" type="email" bind:value={email} required placeholder="admin@pbx.local" class="input" />
      </div>
      
      <div>
        <label for="login-password" class="block text-xs font-semibold text-slate-600 uppercase tracking-wider mb-1.5">Password</label>
        <input id="login-password" type="password" bind:value={password} required placeholder="••••••••" class="input mono" />
      </div>
      
      <button 
        type="submit" 
        disabled={isLoading} 
        class="w-full bg-[#0284C7] hover:bg-sky-700 text-white font-bold py-2.5 rounded-lg transition-colors shadow-sm disabled:opacity-50 text-sm mt-2"
      >
        {isLoading ? 'Authenticating...' : 'Sign In to PBX'}
      </button>

      <p class="text-center text-[11px] text-slate-400 mt-4">
        Default administrator credentials: <span class="mono font-semibold text-slate-600">admin@pbx.local / admin123</span>
      </p>
    </form>
  </div>
</div>
