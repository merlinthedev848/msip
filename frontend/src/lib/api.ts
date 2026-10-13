/**
 * Robust API Client for MSIP PBX
 * Automatically detects environment:
 * - Production: uses relative '/api/v1' (proxied by Nginx on port 80/443 with no firewall issues or CORS)
 * - Development: uses port 8080 (KrakenD)
 * - Automatic multi-tier fallback to port 8080 and 8081 if connection fails
 */

export function getApiBase(): string {
  if (typeof window === 'undefined') return '/api/v1';
  if (window.location.port === '5173') {
    return `/api/v1`;
  }
  return '/api/v1';
}

export function getWsBase(): string {
  if (typeof window === 'undefined') return '';
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  if (window.location.port === '5173') {
    return `${proto}//${window.location.hostname}:8080/api/v1`;
  }
  return `${proto}//${window.location.host}/api/v1`;
}

export const API_BASE = getApiBase();

export async function apiFetch(endpoint: string, options?: RequestInit): Promise<Response> {
  const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const path = cleanEndpoint.startsWith('/api/v1') ? cleanEndpoint.replace('/api/v1', '') : cleanEndpoint;

  // 1. Primary attempt: relative /api/v1 (or API_BASE)
  try {
    const res = await fetch(`${API_BASE}${path}`, options);
    return res;
  } catch (err1) {
    if (typeof window === 'undefined') throw err1;

    // 2. Fallback to Gateway on port 8080
    try {
      return await fetch(`/api/v1${path}`, options);
    } catch (err2) {
      // 3. Fallback to Go Backend directly on port 8081
      return await fetch(`http://${window.location.hostname}:8081/api/v1${path}`, options);
    }
  }
}
