<script lang="ts">
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Button from '../components/ui/Button.svelte';

  const extensions = [
    { id: '1001', name: 'Alice Smith', type: 'SIP', status: 'Online', ip: '192.168.1.45' },
    { id: '1002', name: 'Bob Jones', type: 'WebRTC', status: 'Offline', ip: '-' },
    { id: '1003', name: 'Sales Queue', type: 'Queue', status: 'Online', ip: '-' },
    { id: '1004', name: 'Support', type: 'SIP', status: 'Online', ip: '192.168.1.88' },
  ];
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out">
  
  <header class="flex justify-between items-end mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-white mb-1">Extensions</h1>
      <p class="text-gray-400 text-sm">Manage SIP devices and user endpoints.</p>
    </div>
    <Button variant="primary">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      <span>New Extension</span>
    </Button>
  </header>

  <GlassCard className="p-0 overflow-hidden">
    <!-- Search Bar -->
    <div class="p-4 border-b border-gray-800/50 bg-gray-900/20 flex items-center">
      <svg class="w-5 h-5 text-gray-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
      <input type="text" placeholder="Search extensions..." class="bg-transparent border-none text-gray-300 focus:outline-none w-full placeholder-gray-600 text-sm" />
    </div>

    <!-- Data Table -->
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-gray-900/40 text-gray-400 text-xs uppercase tracking-wider">
            <th class="p-4 font-medium">Ext</th>
            <th class="p-4 font-medium">Name</th>
            <th class="p-4 font-medium">Type</th>
            <th class="p-4 font-medium">Status</th>
            <th class="p-4 font-medium">IP Address</th>
            <th class="p-4 font-medium text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-800/50 text-sm text-gray-300">
          {#each extensions as ext}
            <tr class="hover:bg-gray-800/30 transition-colors group cursor-pointer">
              <td class="p-4 font-mono text-indigo-400 font-medium">{ext.id}</td>
              <td class="p-4 text-white font-medium">{ext.name}</td>
              <td class="p-4">
                <span class="px-2 py-1 bg-gray-800 rounded-md text-xs text-gray-400 border border-gray-700">{ext.type}</span>
              </td>
              <td class="p-4">
                {#if ext.status === 'Online'}
                  <Badge status="success">{ext.status}</Badge>
                {:else}
                  <Badge status="neutral">{ext.status}</Badge>
                {/if}
              </td>
              <td class="p-4 font-mono text-xs text-gray-500">{ext.ip}</td>
              <td class="p-4 text-right">
                <button class="text-gray-500 hover:text-white transition-colors opacity-0 group-hover:opacity-100 p-1">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </GlassCard>

</div>
