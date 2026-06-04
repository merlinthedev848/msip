<script lang="ts">
  import { onMount } from 'svelte';
  let commits = [];

  onMount(async () => {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/commits`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        commits = data.commits || [];
      }
    } catch (e) {}
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-2">Config Version Control</h1>
      <p class="text-slate-500 font-medium text-sm">Audit trail and Git-backed rollback for all PBX dialplan changes.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-indigo-600 hover:bg-indigo-500 text-slate-900 px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"></path></svg>
        Backup Now
      </button>
    </div>
  </header>

  <div class="bg-white/60 rounded-[2rem] border border-slate-200 overflow-hidden shadow-2xl p-8">
    <div class="space-y-6">
      {#each commits as commit, i}
        <div class="relative pl-8 group">
          <!-- Timeline Dot -->
          <div class="absolute -left-[9px] top-1.5 w-4 h-4 rounded-full border-4 border-gray-900 
            {i === 0 ? 'bg-indigo-500 shadow-[0_0_10px_rgba(99,102,241,0.8)]' : 'bg-gray-600'}"></div>
          
          <div class="bg-slate-50/50 rounded-2xl border border-slate-200 p-5 hover:border-indigo-500/30 transition-colors">
            <div class="flex justify-between items-start mb-2">
              <div class="flex items-center space-x-3">
                <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-widest font-mono bg-slate-100 text-slate-500">{commit.hash}</span>
                <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-widest
                  {commit.type === 'Security' ? 'bg-rose-500/10 text-rose-400' :
                   commit.type === 'Provisioning' ? 'bg-emerald-50 text-emerald-400' :
                   'bg-blue-50 text-blue-600'}">
                  {commit.type}
                </span>
              </div>
              <span class="text-xs font-bold text-slate-500">{commit.time}</span>
            </div>
            
            <h3 class="text-lg font-bold text-slate-900 mb-3">{commit.msg}</h3>
            
            <div class="flex justify-between items-center mt-4 pt-4 border-t border-slate-200/50">
              <div class="flex items-center space-x-2">
                <div class="w-6 h-6 rounded-full bg-slate-100 flex items-center justify-center text-[10px] font-bold text-slate-900 border border-slate-200">
                  {commit.author.charAt(0)}
                </div>
                <span class="text-xs font-bold text-slate-500">{commit.author}</span>
              </div>
              
              <div class="flex space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button aria-label="View Diff" class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 rounded-lg text-xs font-bold text-slate-900 transition-colors border border-slate-200 flex items-center">
                  <svg class="w-3 h-3 mr-1.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path></svg>
                  View Diff
                </button>
                <button aria-label="Rollback" class="px-3 py-1.5 bg-rose-500/10 hover:bg-rose-500 rounded-lg text-xs font-bold text-rose-400 hover:text-slate-900 transition-colors border border-rose-500/20 flex items-center">
                  <svg class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"></path></svg>
                  Rollback to Here
                </button>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>


