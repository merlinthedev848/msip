<script lang="ts">
  import { onMount } from 'svelte';
  let commits = [];

  onMount(async () => {
    try {
      const res = await fetch(`/api/v1/commits`, {
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
      <h1 class="text-3xl font-black text-white tracking-tight mb-2">Config Version Control</h1>
      <p class="text-slate-500 font-medium text-sm">Audit trail and Git-backed rollback for all PBX dialplan changes.</p>
    </div>
    <div class="flex space-x-3">
      <button class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-xl font-bold transition-colors shadow-lg shadow-indigo-600/20 flex items-center text-sm">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"></path></svg>
        Backup Now
      </button>
    </div>
  </header>

  <div class="bg-[#121622] rounded-xl border border-[#202738] overflow-hidden shadow-2xl p-6">
    <div class="space-y-5">
      {#each commits as commit, i}
        <div class="relative pl-6 group">
          <!-- Timeline Dot -->
          <div class="absolute -left-[7px] top-2 w-3.5 h-3.5 rounded-full border-2 border-[#121622] 
            {i === 0 ? 'bg-cyan-400 shadow-[0_0_8px_rgba(34,211,238,0.8)]' : 'bg-slate-600'}"></div>
          
          <div class="bg-[#0B0E17] rounded-xl border border-[#1E2538] p-4 hover:border-cyan-500/30 transition-colors">
            <div class="flex justify-between items-start mb-2">
              <div class="flex items-center space-x-2.5">
                <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-widest font-mono bg-[#182030] text-cyan-400 border border-[#26334A]">{commit.hash}</span>
                <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider font-mono
                  {commit.type === 'Security' ? 'bg-rose-950/60 text-rose-400 border border-rose-500/30' :
                   commit.type === 'Provisioning' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-500/30' :
                   'bg-cyan-950/60 text-cyan-400 border border-cyan-500/30'}">
                  {commit.type}
                </span>
              </div>
              <span class="text-xs font-mono text-slate-500">{commit.time}</span>
            </div>
            
            <h3 class="text-sm font-bold text-white mb-2 font-mono">{commit.msg}</h3>
            
            <div class="flex justify-between items-center mt-3 pt-3 border-t border-[#1E2538]">
              <div class="flex items-center space-x-2">
                <div class="w-5 h-5 rounded bg-[#182030] flex items-center justify-center text-[10px] font-bold text-cyan-400 border border-[#26334A] font-mono">
                  {commit.author.charAt(0)}
                </div>
                <span class="text-xs font-mono text-slate-400">{commit.author}</span>
              </div>
              
              <div class="flex space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button aria-label="View Diff" class="px-2.5 py-1 bg-[#182030] hover:bg-[#202B40] rounded text-xs font-mono font-bold text-slate-300 hover:text-white transition-colors border border-[#26334A] flex items-center">
                  <svg class="w-3 h-3 mr-1 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path></svg>
                  View Diff
                </button>
                <button aria-label="Rollback" class="px-2.5 py-1 bg-rose-950/50 hover:bg-rose-900 rounded text-xs font-mono font-bold text-rose-400 hover:text-white transition-colors border border-rose-500/30 flex items-center">
                  <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"></path></svg>
                  Rollback
                </button>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>


