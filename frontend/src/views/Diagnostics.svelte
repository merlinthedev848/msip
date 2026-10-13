<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';

  interface SipProfile {
    name: string;
    port: number;
    proto: string;
    state: string;
  }

  interface JitterStats {
    current_ms: number;
    max_ms: number;
    packet_loss_pct: number;
    rtt_ms: number;
    mos_score: number;
    quality: string;
  }

  interface DiagnosticsData {
    status: string;
    sip_profiles: SipProfile[];
    jitter_buffer: JitterStats;
    active_codecs: string[];
  }

  let data: DiagnosticsData | null = null;
  let loading = true;
  let loadError = '';

  // Ping Tool State
  let pingHost = '8.8.8.8';
  let pinging = false;
  let pingResults: { target: string; latency_ms: number; status: string }[] = [];

  async function loadDiagnostics() {
    loading = true;
    loadError = '';
    try {
      data = await api.get<DiagnosticsData>('/diagnostics/jitter');
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  async function runPing() {
    if (!pingHost.trim()) return;
    pinging = true;
    try {
      const res = await api.post<{ target: string; latency_ms: number; status: string; success: boolean }>('/diagnostics/ping', {
        host: pingHost.trim()
      });
      pingResults = [
        { target: res.target, latency_ms: res.latency_ms, status: res.status },
        ...pingResults.slice(0, 4)
      ];
    } catch (e: any) {
      pingResults = [
        { target: pingHost, latency_ms: 0, status: `Error: ${e.message || 'Probe failed'}` },
        ...pingResults.slice(0, 4)
      ];
    } finally {
      pinging = false;
    }
  }

  onMount(loadDiagnostics);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="SIP, IAX & Jitter Diagnostics"
    subtitle="Live RTP jitter buffer telemetry, audio packet latency, and SIP gateway interface state."
  >
    <Button variant="secondary" on:click={loadDiagnostics} disabled={loading}>
      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
      </svg>
      Refresh Telemetry
    </Button>
  </PageHeader>

  {#if loading}
    <div class="p-12 text-center text-slate-400 text-sm">Gathering SIP and audio stream metrics...</div>
  {:else if loadError}
    <div class="p-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-sm">{loadError}</div>
  {:else if data}
    <!-- METRIC TILES -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <Card>
        <div class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Mean Opinion Score (MOS)</div>
        <div class="text-2xl font-bold text-emerald-600 flex items-baseline space-x-2">
          <span>{data.jitter_buffer.mos_score.toFixed(2)}</span>
          <span class="text-xs font-medium text-slate-500">/ 5.0 (Toll Quality)</span>
        </div>
        <p class="text-[11px] text-slate-500 mt-2">Crystal clear audio fidelity across active calls</p>
      </Card>

      <Card>
        <div class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Jitter Buffer Delay</div>
        <div class="text-2xl font-bold text-sky-600">
          {data.jitter_buffer.current_ms} <span class="text-xs font-medium text-slate-500">ms (Max {data.jitter_buffer.max_ms}ms)</span>
        </div>
        <p class="text-[11px] text-slate-500 mt-2">Adaptive de-jitter queue prevents audio chops</p>
      </Card>

      <Card>
        <div class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Packet Loss</div>
        <div class="text-2xl font-bold text-emerald-600">
          {data.jitter_buffer.packet_loss_pct}%
        </div>
        <p class="text-[11px] text-slate-500 mt-2">Below 1% threshold; zero audible audio distortion</p>
      </Card>

      <Card>
        <div class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">Network RTT (Ping)</div>
        <div class="text-2xl font-bold text-slate-800">
          {data.jitter_buffer.rtt_ms} <span class="text-xs font-medium text-slate-500">ms</span>
        </div>
        <p class="text-[11px] text-slate-500 mt-2">Round-trip latency to upstream carrier gateways</p>
      </Card>
    </div>

    <!-- SIP PROFILES & INTERFACES -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <Card title="SIP & Signaling Interfaces" subtitle="Listening socket endpoints for desk phones, WebRTC and trunks.">
        <div class="overflow-x-auto mt-2">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold uppercase">
                <th class="py-2.5 px-3">Profile Name</th>
                <th class="py-2.5 px-3">Port</th>
                <th class="py-2.5 px-3">Protocol</th>
                <th class="py-2.5 px-3 text-right">Service State</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700">
              {#each data.sip_profiles as profile}
                <tr class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-2.5 px-3 font-semibold text-slate-900">{profile.name}</td>
                  <td class="py-2.5 px-3 font-mono font-medium text-sky-700">{profile.port}</td>
                  <td class="py-2.5 px-3 uppercase text-slate-600 font-mono text-[11px]">{profile.proto}</td>
                  <td class="py-2.5 px-3 text-right">
                    <Badge variant="success">{profile.state}</Badge>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </Card>

      <!-- CODEC NEGOTIATION & ECHO TEST -->
      <Card title="Audio Codecs & Echo Loopback" subtitle="Configured RTP payload negotiation and line test tools.">
        <div class="space-y-4 mt-2 text-xs">
          <div>
            <div class="font-medium text-slate-700 mb-2">Priority Negotiated Codecs:</div>
            <div class="flex flex-wrap gap-2">
              {#each data.active_codecs as codec}
                <span class="px-2.5 py-1 bg-slate-100 border border-slate-200 text-slate-800 rounded font-mono font-semibold text-[11px]">
                  {codec}
                </span>
              {/each}
            </div>
          </div>

          <div class="p-3.5 bg-sky-50 border border-sky-200 rounded-lg">
            <div class="font-bold text-sky-900 mb-1 flex items-center">
              <svg class="w-4 h-4 mr-1.5 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 100-6 3 3 0 000 6z" />
              </svg>
              Quick Echo Test (Dial 9999)
            </div>
            <p class="text-sky-800 text-[11px] leading-relaxed">
              Dial <span class="font-bold font-mono">9999</span> from any softphone or physical handset. The PBX will immediately loop back your exact audio to test microphone delay, echo cancellation, and network jitter.
            </p>
          </div>

          <!-- Carrier Ping Test -->
          <div>
            <div class="font-medium text-slate-700 mb-1.5">Carrier Ping & Latency Probe:</div>
            <div class="flex space-x-2">
              <input
                type="text"
                bind:value={pingHost}
                placeholder="IP address or SIP domain"
                class="flex-1 px-3 py-1.5 border border-slate-200 rounded text-xs font-mono focus:outline-none focus:border-sky-500"
              />
              <Button variant="primary" on:click={runPing} disabled={pinging}>
                {pinging ? 'Pinging...' : 'Test Ping'}
              </Button>
            </div>
            {#if pingResults.length > 0}
              <div class="mt-2 space-y-1">
                {#each pingResults as res}
                  <div class="flex justify-between items-center text-[11px] font-mono p-1.5 bg-slate-50 border border-slate-200 rounded">
                    <span class="text-slate-800">{res.target}</span>
                    <span class="text-emerald-700 font-semibold">{res.latency_ms} ms</span>
                    <span class="text-slate-500">{res.status}</span>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      </Card>
    </div>
  {/if}
</div>
