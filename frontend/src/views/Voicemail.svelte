<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-white mb-1">Visual Voicemail</h1>
      <p class="text-gray-400 text-sm">Listen and read transcriptions of incoming messages.</p>
    </div>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <GlassCard hoverEffect={true}>
      <div class="flex justify-between items-start mb-4">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-full bg-indigo-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
          </div>
          <div>
            <p class="font-semibold text-white">+1 (555) 098-7654</p>
            <p class="text-xs text-gray-500">Today, 10:23 AM • 0:45</p>
          </div>
        </div>
        <Badge status="error">New</Badge>
      </div>
      
      <div class="bg-gray-900/50 p-4 rounded-lg border border-gray-800 mb-4">
        <p class="text-sm text-gray-300 italic font-serif leading-relaxed">"Hi, this is Sarah from accounting. I had a question about the recent invoice. Please call me back when you get a chance."</p>
      </div>
      
      <div class="flex items-center space-x-4">
        <button title="Play Voicemail" aria-label="Play Voicemail" class="w-10 h-10 rounded-full bg-indigo-600 hover:bg-indigo-500 flex items-center justify-center transition-colors shadow-lg shadow-indigo-500/30">
          <svg class="w-4 h-4 text-white ml-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd"></path></svg>
        </button>
        <div class="flex-1 h-2 bg-gray-800 rounded-full overflow-hidden">
          <div class="h-full bg-indigo-500 w-1/3"></div>
        </div>
        <span class="text-xs text-gray-500 font-mono">0:15 / 0:45</span>
      </div>
    </GlassCard>
  </div>
</div>

