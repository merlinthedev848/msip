// Layout style store — permanent single layout
import { writable } from 'svelte/store';

export type LayoutStyle = 'standard';
export const layoutStyle = writable<LayoutStyle>('standard');
export function setLayoutStyle(_style: any) {
  layoutStyle.set('standard');
}
