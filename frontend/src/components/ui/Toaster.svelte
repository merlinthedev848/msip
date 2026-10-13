<script lang="ts">
  import { toasts, dismiss } from '../../lib/toast';
  const tones = {
    success: 'border-emerald-200 bg-emerald-50 text-emerald-800',
    error: 'border-rose-200 bg-rose-50 text-rose-800',
    info: 'border-sky-200 bg-sky-50 text-sky-800',
  };
</script>

<div class="fixed top-14 right-4 z-[60] space-y-2 w-80">
  {#each $toasts as t (t.id)}
    <div class="border rounded-md shadow-md px-3 py-2 text-xs flex items-start justify-between gap-2 {tones[t.tone]}" role="status">
      <span>{t.message}</span>
      <button type="button" class="opacity-60 hover:opacity-100" aria-label="Dismiss" title="Dismiss" on:click={() => dismiss(t.id)}>✕</button>
    </div>
  {/each}
</div>
