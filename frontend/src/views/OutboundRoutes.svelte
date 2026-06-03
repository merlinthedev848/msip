<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Button from '../components/ui/Button.svelte';
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-white mb-1">Outbound Routes</h1>
      <p class="text-gray-400 text-sm">Define dial patterns and prioritize carrier trunks.</p>
    </div>
    <Button variant="primary">New Rule</Button>
  </header>

  <GlassCard hoverEffect={true} className="border-l-4 border-l-indigo-500">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
      <div>
        <h3 class="text-lg font-bold text-white mb-1">International Standard</h3>
        <p class="text-sm font-mono text-gray-400 bg-gray-900/50 inline-block px-2 py-1 rounded">Pattern: ^011([1-9][0-9]{6,14})$</p>
      </div>
      <div class="mt-4 md:mt-0 flex items-center space-x-4">
        <div class="text-right">
          <p class="text-xs text-gray-500 uppercase">Primary Trunk</p>
          <p class="text-sm font-semibold text-emerald-400">Telnyx SIP TLS</p>
        </div>
        <Button variant="secondary" className="px-3 py-1">Edit</Button>
      </div>
    </div>
  </GlassCard>

  <GlassCard hoverEffect={true} className="border-l-4 border-l-emerald-500">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
      <div>
        <h3 class="text-lg font-bold text-white mb-1">US Domestic (Toll-Free)</h3>
        <p class="text-sm font-mono text-gray-400 bg-gray-900/50 inline-block px-2 py-1 rounded">Pattern: ^1?(800|888|877|866|855|844|833)[0-9]{7}$</p>
      </div>
      <div class="mt-4 md:mt-0 flex items-center space-x-4">
        <div class="text-right">
          <p class="text-xs text-gray-500 uppercase">Primary Trunk</p>
          <p class="text-sm font-semibold text-emerald-400">Twilio Elastic</p>
        </div>
        <Button variant="secondary" className="px-3 py-1">Edit</Button>
      </div>
    </div>
  </GlassCard>
</div>

