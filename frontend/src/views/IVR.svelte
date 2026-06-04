<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Button from '../components/ui/Button.svelte';
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-4">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-slate-900 mb-1">IVR Builder</h1>
      <p class="text-slate-500 text-sm">Visual node designer for interactive voice response menus.</p>
    </div>
    <div class="flex space-x-3">
      <Button variant="secondary">Upload Audio</Button>
      <Button variant="primary">Save Tree</Button>
    </div>
  </header>

  <GlassCard className="flex-1 min-h-[500px] relative overflow-hidden bg-white/80">
    <!-- Grid Background -->
    <div class="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px]"></div>
    
    <!-- Dummy Nodes -->
    <div class="absolute inset-0 p-8 flex flex-col items-center">
      
      <!-- Start Node -->
      <div class="bg-indigo-600 border border-indigo-400 rounded-xl px-6 py-3 shadow-lg shadow-indigo-500/20 z-10 cursor-move hover:-translate-y-1 transition-transform">
        <p class="text-slate-900 font-bold">Main Welcome (Start)</p>
        <p class="text-indigo-200 text-xs">Audio: greeting.wav</p>
      </div>
      
      <!-- Connector -->
      <div class="w-px h-12 bg-indigo-500/50"></div>
      
      <!-- Split Node -->
      <div class="bg-slate-100 border border-slate-300 rounded-xl px-6 py-3 shadow-lg z-10 text-center">
        <p class="text-slate-900 font-semibold mb-2">Input Menu</p>
        <div class="flex space-x-8">
          <div class="flex flex-col items-center">
            <span class="w-6 h-6 bg-slate-200 rounded text-xs flex items-center justify-center font-bold text-slate-700">1</span>
            <div class="w-px h-8 bg-gray-600 mt-2"></div>
            <div class="bg-emerald-500/20 border border-emerald-500/50 text-emerald-400 text-xs px-3 py-1 rounded mt-2">Sales Queue</div>
          </div>
          <div class="flex flex-col items-center">
            <span class="w-6 h-6 bg-slate-200 rounded text-xs flex items-center justify-center font-bold text-slate-700">2</span>
            <div class="w-px h-8 bg-gray-600 mt-2"></div>
            <div class="bg-blue-500/20 border border-blue-500/50 text-blue-400 text-xs px-3 py-1 rounded mt-2">Support</div>
          </div>
          <div class="flex flex-col items-center">
            <span class="w-6 h-6 bg-slate-200 rounded text-xs flex items-center justify-center font-bold text-slate-700">t</span>
            <div class="w-px h-8 bg-gray-600 mt-2"></div>
            <div class="bg-red-500/20 border border-red-500/50 text-red-400 text-xs px-3 py-1 rounded mt-2">Hangup</div>
          </div>
        </div>
      </div>
      
    </div>
  </GlassCard>
</div>

