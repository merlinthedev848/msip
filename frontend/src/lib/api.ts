/**
 * MSIP API client.
 *
 * All browser traffic goes through the same origin under /api/v1, which nginx
 * (production) or the Vite dev proxy forwards to the KrakenD gateway.
 * Every request carries the session JWT, disables HTTP caching, parses JSON
 * safely and surfaces backend error messages as ApiError instances.
 */

export const API_BASE = '/api/v1';

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export function getToken(): string | null {
  return localStorage.getItem('pbx_token');
}

export function getCurrentUser(): { id?: string; email?: string; role?: string; tenantId?: string } {
  try {
    return JSON.parse(localStorage.getItem('pbx_user') || '{}');
  } catch {
    return {};
  }
}

export function logout() {
  localStorage.removeItem('pbx_token');
  localStorage.removeItem('pbx_user');
  window.location.href = '/';
}

async function request<T = any>(method: string, path: string, body?: unknown): Promise<T> {
  const clean = path.startsWith('/') ? path : `/${path}`;
  const headers: Record<string, string> = { Accept: 'application/json' };
  const token = getToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;
  if (body !== undefined) headers['Content-Type'] = 'application/json';

  let res: Response;
  try {
    res = await fetch(`${API_BASE}${clean}`, {
      method,
      headers,
      cache: 'no-store',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError(0, 'Unable to reach the PBX API. Check that the backend and gateway containers are running.');
  }

  if (res.status === 401 && clean !== '/auth/login') {
    logout();
    throw new ApiError(401, 'Session expired. Please sign in again.');
  }

  const text = await res.text();
  let data: any = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = null;
    }
  }

  if (!res.ok) {
    const msg = (data && (data.error || data.message)) || `Request failed with status ${res.status}`;
    throw new ApiError(res.status, msg);
  }
  return data as T;
}

export const api = {
  get: <T = any>(path: string) => request<T>('GET', path),
  post: <T = any>(path: string, body?: unknown) => request<T>('POST', path, body ?? {}),
  put: <T = any>(path: string, body?: unknown) => request<T>('PUT', path, body ?? {}),
  del: <T = any>(path: string) => request<T>('DELETE', path),
};

/** Builds an absolute ws:// or wss:// URL for a same-origin API path. */
export function wsUrl(path: string): string {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const clean = path.startsWith('/') ? path : `/${path}`;
  return `${proto}//${window.location.host}${API_BASE}${clean}`;
}

/** Formats an ISO timestamp for tables. */
export function formatDateTime(value?: string): string {
  if (!value) return '—';
  const d = new Date(value);
  if (isNaN(d.getTime()) || d.getFullYear() < 2000) return '—';
  return d.toLocaleString(undefined, { year: 'numeric', month: 'short', day: '2-digit', hour: '2-digit', minute: '2-digit' });
}

export function formatDate(value?: string): string {
  if (!value) return '—';
  const d = new Date(value);
  if (isNaN(d.getTime()) || d.getFullYear() < 2000) return '—';
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: '2-digit' });
}

/** Formats seconds as m:ss or h:mm:ss. */
export function formatDuration(totalSeconds?: number): string {
  const s = Math.max(0, Math.floor(totalSeconds || 0));
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = (s % 60).toString().padStart(2, '0');
  return h > 0 ? `${h}:${m.toString().padStart(2, '0')}:${sec}` : `${m}:${sec}`;
}
