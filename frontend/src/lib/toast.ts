import { writable } from 'svelte/store';

export type ToastTone = 'success' | 'error' | 'info';
export interface Toast {
  id: number;
  tone: ToastTone;
  message: string;
}

export const toasts = writable<Toast[]>([]);

let seq = 0;
function push(tone: ToastTone, message: string, ms = 4000) {
  const id = ++seq;
  toasts.update((t) => [...t, { id, tone, message }]);
  setTimeout(() => dismiss(id), ms);
}

export function dismiss(id: number) {
  toasts.update((t) => t.filter((x) => x.id !== id));
}

export const toast = {
  success: (m: string) => push('success', m),
  error: (m: string) => push('error', m, 6000),
  info: (m: string) => push('info', m),
};

/** Convenience: show an error toast from any thrown value. */
export function toastError(e: unknown, fallback = 'Something went wrong') {
  const msg = e instanceof Error && e.message ? e.message : fallback;
  toast.error(msg);
}
