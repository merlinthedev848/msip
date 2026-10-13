import { writable } from 'svelte/store';

function normalize(path: string): string {
  if (!path) return '/';
  return path.length > 1 && path.endsWith('/') ? path.slice(0, -1) : path;
}

export const currentPath = writable<string>(normalize(window.location.pathname));

window.addEventListener('popstate', () => currentPath.set(normalize(window.location.pathname)));

export function navigate(path: string) {
  const target = normalize(path);
  if (target !== normalize(window.location.pathname)) {
    window.history.pushState({}, '', target);
  }
  currentPath.set(target);
}
