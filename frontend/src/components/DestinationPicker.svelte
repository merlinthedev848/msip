<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { api } from '../lib/api';

  export let value: string = 'Not Set';
  export let label: string = '';

  const dispatch = createEventDispatcher();

  let modalOpen = false;
  let activeCategory: 'extension' | 'menu' | 'queue' | 'group' | 'voicemail' | 'external' | 'hangup' = 'extension';

  let extensions: any[] = [];
  let menus: any[] = [];
  let queues: any[] = [];
  let groups: any[] = [];
  let externalNumber: string = '';

  onMount(async () => {
    try {
      const [extRes, ivrRes, qRes, rgRes] = await Promise.all([
        api.get<{ extensions: any[] }>('/extensions').catch(() => ({ extensions: [] })),
        api.get<{ ivrs: any[] }>('/ivrs').catch(() => ({ ivrs: [] })),
        api.get<{ queues: any[] }>('/queues').catch(() => ({ queues: [] })),
        api.get<{ ring_groups: any[] }>('/ring-groups').catch(() => ({ ring_groups: [] }))
      ]);
      extensions = extRes.extensions || [];
      menus = ivrRes.ivrs || [];
      queues = qRes.queues || [];
      groups = rgRes.ring_groups || [];
    } catch (e) {}
  });

  function choose(target: string) {
    value = target;
    dispatch('change', target);
    modalOpen = false;
  }

  function getIcon(val: string): string {
    if (!val || val === 'Not Set' || val === 'Unknown') return '?';
    if (val.toLowerCase().includes('end call') || val.toLowerCase().includes('hangup')) return '🛑';
    if (val.toLowerCase().includes('menu') || val.toLowerCase().includes('ivr')) return '≡';
    if (val.toLowerCase().includes('group') || val.toLowerCase().includes('door') || val.toLowerCase().includes('hunt')) return '👥';
    if (val.toLowerCase().includes('queue') || val.toLowerCase().includes('accounts')) return '⏳';
    if (val.toLowerCase().includes('voicemail')) return '✉';
    return '📞';
  }
</script>

<div class="w-full">
  {#if label}
    <label class="block text-xs font-semibold text-slate-700 mb-1">{label}</label>
  {/if}
  
  <div class="inline-flex w-full rounded-md shadow-sm border border-slate-300 overflow-hidden bg-white text-xs">
    <button
      type="button"
      on:click={() => modalOpen = true}
      class="px-3 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold border-r border-slate-300 flex items-center space-x-1.5 transition-colors focus:outline-none"
    >
      <svg class="w-3.5 h-3.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
      </svg>
      <span>Select</span>
      <svg class="w-2.5 h-2.5 text-slate-400 ml-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>
    <div class="px-3 py-2 flex-1 bg-white text-slate-800 flex items-center space-x-2 truncate font-medium">
      <span class="text-slate-500 font-bold">{getIcon(value)}</span>
      <span class="truncate">{value || 'Not Set'}</span>
    </div>
  </div>
</div>

<!-- Modal Destination Selector -->
{#if modalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-[2px] animate-fade-in">
    <div class="bg-white rounded-lg shadow-2xl border border-slate-200 w-full max-w-xl overflow-hidden flex flex-col max-h-[85vh]">
      <!-- Header -->
      <div class="px-5 py-3.5 bg-gradient-to-r from-[#1E3A8A] to-[#0284C7] text-white flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <svg class="w-4 h-4 text-sky-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7" />
          </svg>
          <span class="font-bold text-sm tracking-wide">Choose Call Destination</span>
        </div>
        <button
          type="button"
          on:click={() => modalOpen = false}
          class="text-white/80 hover:text-white p-1 rounded hover:bg-white/10"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Category Tabs -->
      <div class="flex border-b border-slate-200 bg-slate-50 overflow-x-auto text-xs px-2 pt-2">
        <button
          type="button"
          on:click={() => activeCategory = 'extension'}
          class="px-3 py-2 font-semibold border-b-2 whitespace-nowrap {activeCategory === 'extension' ? 'border-[#0284C7] text-[#0284C7] bg-white rounded-t' : 'border-transparent text-slate-600 hover:text-slate-900'}"
        >
          Extension
        </button>
        <button
          type="button"
          on:click={() => activeCategory = 'group'}
          class="px-3 py-2 font-semibold border-b-2 whitespace-nowrap {activeCategory === 'group' ? 'border-[#0284C7] text-[#0284C7] bg-white rounded-t' : 'border-transparent text-slate-600 hover:text-slate-900'}"
        >
          Ring Group
        </button>
        <button
          type="button"
          on:click={() => activeCategory = 'queue'}
          class="px-3 py-2 font-semibold border-b-2 whitespace-nowrap {activeCategory === 'queue' ? 'border-[#0284C7] text-[#0284C7] bg-white rounded-t' : 'border-transparent text-slate-600 hover:text-slate-900'}"
        >
          Call Queue
        </button>
        <button
          type="button"
          on:click={() => activeCategory = 'menu'}
          class="px-3 py-2 font-semibold border-b-2 whitespace-nowrap {activeCategory === 'menu' ? 'border-[#0284C7] text-[#0284C7] bg-white rounded-t' : 'border-transparent text-slate-600 hover:text-slate-900'}"
        >
          IVR Menu
        </button>
        <button
          type="button"
          on:click={() => activeCategory = 'voicemail'}
          class="px-3 py-2 font-semibold border-b-2 whitespace-nowrap {activeCategory === 'voicemail' ? 'border-[#0284C7] text-[#0284C7] bg-white rounded-t' : 'border-transparent text-slate-600 hover:text-slate-900'}"
        >
          Voicemail
        </button>
        <button
          type="button"
          on:click={() => activeCategory = 'external'}
          class="px-3 py-2 font-semibold border-b-2 whitespace-nowrap {activeCategory === 'external' ? 'border-[#0284C7] text-[#0284C7] bg-white rounded-t' : 'border-transparent text-slate-600 hover:text-slate-900'}"
        >
          External
        </button>
        <button
          type="button"
          on:click={() => activeCategory = 'hangup'}
          class="px-3 py-2 font-semibold border-b-2 whitespace-nowrap {activeCategory === 'hangup' ? 'border-[#0284C7] text-[#0284C7] bg-white rounded-t' : 'border-transparent text-slate-600 hover:text-slate-900'}"
        >
          System Action
        </button>
      </div>

      <!-- Selection Body -->
      <div class="p-4 overflow-y-auto flex-1 max-h-96">
        {#if activeCategory === 'extension'}
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {#each extensions as ext}
              <button
                type="button"
                on:click={() => choose(`Ext ${ext.ExtensionNumber} (${ext.DisplayName || ext.FirstName || 'User'})`)}
                class="p-2.5 text-left border rounded-lg hover:border-[#0284C7] hover:bg-sky-50 transition-colors flex items-center justify-between group"
              >
                <div>
                  <div class="font-bold text-xs text-slate-800 group-hover:text-[#0284C7]">{ext.ExtensionNumber}</div>
                  <div class="text-[11px] text-slate-500">{ext.DisplayName || `${ext.FirstName || ''} ${ext.LastName || ''}`.trim() || 'Internal User'}</div>
                </div>
                <span class="text-xs text-slate-400 group-hover:text-[#0284C7]">→</span>
              </button>
            {/each}
            {#if extensions.length === 0}
              <div class="col-span-2 text-center py-6 text-slate-400 text-xs">No extensions configured</div>
            {/if}
          </div>

        {:else if activeCategory === 'group'}
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {#each groups as grp}
              <button
                type="button"
                on:click={() => choose(`${grp.Name}`)}
                class="p-2.5 text-left border rounded-lg hover:border-[#0284C7] hover:bg-sky-50 transition-colors flex items-center justify-between group"
              >
                <div>
                  <div class="font-bold text-xs text-slate-800 group-hover:text-[#0284C7]">{grp.Name}</div>
                  <div class="text-[11px] text-slate-500">Ext {grp.InternalNumber} • {grp.SiteName || 'Main'}</div>
                </div>
                <span class="text-xs text-slate-400 group-hover:text-[#0284C7]">→</span>
              </button>
            {/each}
            {#if groups.length === 0}
              <div class="col-span-2 text-center py-6 text-slate-400 text-xs">No ring groups configured</div>
            {/if}
          </div>

        {:else if activeCategory === 'queue'}
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {#each queues as q}
              <button
                type="button"
                on:click={() => choose(`${q.Name}`)}
                class="p-2.5 text-left border rounded-lg hover:border-[#0284C7] hover:bg-sky-50 transition-colors flex items-center justify-between group"
              >
                <div>
                  <div class="font-bold text-xs text-slate-800 group-hover:text-[#0284C7]">{q.Name}</div>
                  <div class="text-[11px] text-slate-500">Strategy: {q.RingStrategy} • Ext {q.InternalNumber || q.Extension}</div>
                </div>
                <span class="text-xs text-slate-400 group-hover:text-[#0284C7]">→</span>
              </button>
            {/each}
            {#if queues.length === 0}
              <div class="col-span-2 text-center py-6 text-slate-400 text-xs">No queues configured</div>
            {/if}
          </div>

        {:else if activeCategory === 'menu'}
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {#each menus as m}
              <button
                type="button"
                on:click={() => choose(`${m.Name}`)}
                class="p-2.5 text-left border rounded-lg hover:border-[#0284C7] hover:bg-sky-50 transition-colors flex items-center justify-between group"
              >
                <div>
                  <div class="font-bold text-xs text-slate-800 group-hover:text-[#0284C7]">{m.Name}</div>
                  <div class="text-[11px] text-slate-500">Ext {m.Extension} • Timeout: {m.Timeout / 1000}s</div>
                </div>
                <span class="text-xs text-slate-400 group-hover:text-[#0284C7]">→</span>
              </button>
            {/each}
            {#if menus.length === 0}
              <div class="col-span-2 text-center py-6 text-slate-400 text-xs">No IVR menus configured</div>
            {/if}
          </div>

        {:else if activeCategory === 'voicemail'}
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            <button
              type="button"
              on:click={() => choose('General Voicemail (*98)')}
              class="p-2.5 text-left border rounded-lg hover:border-[#0284C7] hover:bg-sky-50 transition-colors flex items-center justify-between group"
            >
              <div>
                <div class="font-bold text-xs text-slate-800 group-hover:text-[#0284C7]">General Company Voicemail</div>
                <div class="text-[11px] text-slate-500">Box *98</div>
              </div>
              <span class="text-xs text-slate-400 group-hover:text-[#0284C7]">→</span>
            </button>
            {#each extensions as ext}
              <button
                type="button"
                on:click={() => choose(`Voicemail Ext ${ext.ExtensionNumber}`)}
                class="p-2.5 text-left border rounded-lg hover:border-[#0284C7] hover:bg-sky-50 transition-colors flex items-center justify-between group"
              >
                <div>
                  <div class="font-bold text-xs text-slate-800 group-hover:text-[#0284C7]">Voicemail Box {ext.ExtensionNumber}</div>
                  <div class="text-[11px] text-slate-500">{ext.DisplayName || ext.FirstName || 'User'}</div>
                </div>
                <span class="text-xs text-slate-400 group-hover:text-[#0284C7]">→</span>
              </button>
            {/each}
          </div>

        {:else if activeCategory === 'external'}
          <div class="space-y-3">
            <label class="block text-xs font-semibold text-slate-700">Forward to External Telephone / Mobile Number</label>
            <div class="flex space-x-2">
              <input
                type="text"
                bind:value={externalNumber}
                placeholder="e.g. +447123456789 or 01233220950"
                class="flex-1 px-3 py-2 border border-slate-300 rounded-md text-xs focus:ring-1 focus:ring-[#0284C7] focus:outline-none"
              />
              <button
                type="button"
                on:click={() => { if (externalNumber.trim()) choose(`External: ${externalNumber.trim()}`); }}
                class="px-4 py-2 bg-[#0284C7] hover:bg-sky-700 text-white font-bold text-xs rounded-md shadow-sm transition-colors"
              >
                Set
              </button>
            </div>
            <p class="text-[11px] text-slate-500">Calls routed to external numbers will consume carrier trunk channels.</p>
          </div>

        {:else if activeCategory === 'hangup'}
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            <button
              type="button"
              on:click={() => choose('End Call')}
              class="p-3 text-left border border-rose-200 bg-rose-50/50 rounded-lg hover:bg-rose-100 transition-colors flex items-center justify-between group"
            >
              <div>
                <div class="font-bold text-xs text-rose-800">End Call (Hang Up)</div>
                <div class="text-[11px] text-rose-600">Disconnects the call with normal clearing</div>
              </div>
              <span class="text-xs text-rose-400">🛑</span>
            </button>
            <button
              type="button"
              on:click={() => choose('Play Busy Tone')}
              class="p-3 text-left border border-amber-200 bg-amber-50/50 rounded-lg hover:bg-amber-100 transition-colors flex items-center justify-between group"
            >
              <div>
                <div class="font-bold text-xs text-amber-800">Play Busy Tone</div>
                <div class="text-[11px] text-amber-600">Returns standard SIP 486 Busy tone</div>
              </div>
              <span class="text-xs text-amber-400">⚠️</span>
            </button>
          </div>
        {/if}
      </div>

      <!-- Footer -->
      <div class="px-5 py-2.5 bg-slate-50 border-t border-slate-200 flex justify-end">
        <button
          type="button"
          on:click={() => modalOpen = false}
          class="px-4 py-1.5 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 font-semibold text-xs rounded shadow-sm transition-colors"
        >
          Cancel
        </button>
      </div>
    </div>
  </div>
{/if}
