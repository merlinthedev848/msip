import { writable } from 'svelte/store';

export type LayoutStyle = 'communigate' | 'yeastar';

const stored = typeof localStorage !== 'undefined' ? localStorage.getItem('pbx_layout_style') as LayoutStyle : null;
const initialLayout: LayoutStyle = (stored === 'yeastar' || stored === 'communigate') ? stored : 'yeastar';

export const layoutStyle = writable<LayoutStyle>(initialLayout);

export function setLayoutStyle(style: LayoutStyle) {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('pbx_layout_style', style);
  }
  layoutStyle.set(style);
  window.dispatchEvent(new CustomEvent('layout-changed', { detail: { style } }));
}
