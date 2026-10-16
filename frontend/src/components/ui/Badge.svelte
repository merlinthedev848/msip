<script context="module" lang="ts">
  export type Tone = 'success' | 'warning' | 'error' | 'danger' | 'neutral' | 'info';
</script>

<script lang="ts">
  export let tone: Tone = 'neutral';
  export let variant: Tone | undefined = undefined;
  /** @deprecated use tone */
  export let status: Tone | undefined = undefined;

  $: t = variant ?? status ?? tone;
  const styles: Record<Tone, string> = {
    success: 'bg-emerald-50 text-emerald-700 border-emerald-200',
    warning: 'bg-amber-50 text-amber-700 border-amber-200',
    error: 'bg-rose-50 text-rose-700 border-rose-200',
    danger: 'bg-rose-50 text-rose-700 border-rose-200',
    neutral: 'bg-slate-100 text-slate-600 border-slate-200',
    info: 'bg-sky-50 text-sky-700 border-sky-200',
  };
</script>

<span class="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-medium border {styles[t] || styles.neutral}">
  <slot />
</span>
