<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import Modal from '../components/ui/Modal.svelte';

  interface SIPSettingsData {
    ID?: string;
    TenantID?: string;

    // Signaling & Listening Ports
    user_agent: string;
    sip_port: number;
    sip_tls_port: number;
    external_sip_port: number;
    external_sip_tls_port: number;
    wss_port: number;
    tcp_keepalive_sec: number;
    compact_headers: boolean;
    allow_contact_rewrite: boolean;

    // RTP Media & Jitter Buffer
    rtp_port_start: number;
    rtp_port_end: number;
    rtp_timeout_sec: number;
    rtp_hold_timeout_sec: number;
    jitter_buffer_mode: string;
    jitter_buffer_initial_ms: number;
    jitter_buffer_max_ms: number;
    dtmf_mode: string;
    dtmf_duration_ms: number;
    audio_codecs: string;
    video_codecs: string;
    enable_cng: boolean;

    // NAT, STUN, TURN & ICE
    nat_mode: string;
    external_ipv4: string;
    external_ipv6: string;
    local_networks: string;
    stun_server: string;
    turn_server: string;
    turn_user: string;
    turn_password?: string;
    enable_ice: boolean;
    aggressive_ice: boolean;
    nat_ping_interval_sec: number;

    // SIP Timers & Registration TTL
    session_expires_sec: number;
    min_se_sec: number;
    session_refresher: string;
    register_min_sec: number;
    register_max_sec: number;
    register_default_sec: number;
    register_retry_sec: number;
    max_channels: number;

    // Security & SRTP
    srtp_mode: string;
    srtp_suites: string;
    tls_version_min: string;
    tls_verify_peer: boolean;
    auth_realm: string;
    challenge_all_invites: boolean;

    // QoS / DSCP
    dscp_audio: number;
    dscp_video: number;
    dscp_signaling: number;
  }

  let settings: SIPSettingsData = {
    user_agent: 'MSIP-PBX-Core/2.0',
    sip_port: 5060,
    sip_tls_port: 5061,
    external_sip_port: 5080,
    external_sip_tls_port: 5081,
    wss_port: 7443,
    tcp_keepalive_sec: 30,
    compact_headers: false,
    allow_contact_rewrite: true,

    rtp_port_start: 16384,
    rtp_port_end: 32768,
    rtp_timeout_sec: 30,
    rtp_hold_timeout_sec: 1800,
    jitter_buffer_mode: 'adaptive',
    jitter_buffer_initial_ms: 60,
    jitter_buffer_max_ms: 200,
    dtmf_mode: 'rfc2833',
    dtmf_duration_ms: 100,
    audio_codecs: 'OPUS,G722,PCMA,PCMU',
    video_codecs: 'VP8,H264',
    enable_cng: true,

    nat_mode: 'auto',
    external_ipv4: '',
    external_ipv6: '',
    local_networks: '127.0.0.1/32, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16',
    stun_server: 'stun:stun.l.google.com:19302',
    turn_server: '',
    turn_user: '',
    turn_password: '',
    enable_ice: true,
    aggressive_ice: true,
    nat_ping_interval_sec: 20,

    session_expires_sec: 1800,
    min_se_sec: 120,
    session_refresher: 'auto',
    register_min_sec: 60,
    register_max_sec: 3600,
    register_default_sec: 3600,
    register_retry_sec: 30,
    max_channels: 250,

    srtp_mode: 'optional',
    srtp_suites: 'AEAD_AES_256_GCM,AES_CM_128_HMAC_SHA1_80,AES_CM_128_HMAC_SHA1_32',
    tls_version_min: '1.2',
    tls_verify_peer: false,
    auth_realm: 'msip.local',
    challenge_all_invites: true,

    dscp_audio: 46,
    dscp_video: 34,
    dscp_signaling: 26,
  };

  let loading = true;
  let saving = false;
  let applying = false;
  let resetting = false;

  let activeTab: 'signaling' | 'rtp' | 'nat' | 'timers' | 'security' | 'qos' = 'signaling';
  let message = '';
  let messageType: 'success' | 'error' = 'success';
  let showResetModal = false;
  let applyResult: any = null;

  async function loadSettings() {
    loading = true;
    message = '';
    try {
      const data = await api.get<SIPSettingsData>('/sip-settings');
      if (data) {
        settings = { ...settings, ...data };
      }
    } catch (e: any) {
      message = e.message || 'Failed to load SIP settings';
      messageType = 'error';
    } finally {
      loading = false;
    }
  }

  function validate(): string | null {
    if (!settings.user_agent.trim()) return 'User Agent identifier cannot be empty';
    if (settings.sip_port < 1024 && settings.sip_port !== 5060) return 'Internal SIP port must be 5060 or between 1024 and 65535';
    if (settings.sip_tls_port < 1024 && settings.sip_tls_port !== 5061) return 'Internal SIP TLS port must be 5061 or between 1024 and 65535';
    if (settings.external_sip_port < 1024 || settings.external_sip_port > 65535) return 'External SIP port must be between 1024 and 65535';
    if (settings.wss_port < 1024 || settings.wss_port > 65535) return 'WebRTC WSS port must be between 1024 and 65535';

    if (settings.rtp_port_start >= settings.rtp_port_end) return 'RTP Start Port must be strictly less than RTP End Port';
    if (settings.rtp_port_end - settings.rtp_port_start < 100) return 'RTP Port span must be at least 100 ports for audio stream capacity';

    if (
      (settings.sip_port >= settings.rtp_port_start && settings.sip_port <= settings.rtp_port_end) ||
      (settings.sip_tls_port >= settings.rtp_port_start && settings.sip_tls_port <= settings.rtp_port_end) ||
      (settings.external_sip_port >= settings.rtp_port_start && settings.external_sip_port <= settings.rtp_port_end) ||
      (settings.wss_port >= settings.rtp_port_start && settings.wss_port <= settings.rtp_port_end)
    ) {
      return 'Signaling ports (SIP/TLS/WSS) cannot overlap with RTP media port range';
    }

    if (settings.jitter_buffer_initial_ms < 0 || settings.jitter_buffer_max_ms < settings.jitter_buffer_initial_ms) {
      return 'Jitter Buffer max delay must be greater than or equal to initial delay';
    }

    if (settings.min_se_sec < 90) return 'Min-SE (Session Expires) cannot be less than 90 seconds (RFC 4028)';
    if (settings.session_expires_sec < settings.min_se_sec) return 'Session Expires must be greater than or equal to Min-SE';

    if (settings.register_min_sec <= 0 || settings.register_max_sec < settings.register_min_sec) {
      return 'Registration Max TTL must be greater than or equal to Min TTL';
    }

    if (
      settings.dscp_audio < 0 || settings.dscp_audio > 63 ||
      settings.dscp_video < 0 || settings.dscp_video > 63 ||
      settings.dscp_signaling < 0 || settings.dscp_signaling > 63
    ) {
      return 'DSCP values must be in the valid DiffServ range [0, 63]';
    }

    return null;
  }

  async function handleSave() {
    const err = validate();
    if (err) {
      message = err;
      messageType = 'error';
      return;
    }

    saving = true;
    message = '';
    applyResult = null;
    try {
      const updated = await api.put<SIPSettingsData>('/sip-settings', settings);
      settings = { ...settings, ...updated };
      message = 'SIP and RTP settings saved successfully. Remember to apply changes to reload the active stack.';
      messageType = 'success';
    } catch (e: any) {
      message = e.message || 'Failed to save SIP settings';
      messageType = 'error';
    } finally {
      saving = false;
    }
  }

  async function handleApply() {
    applying = true;
    message = '';
    try {
      const res = await api.post<any>('/sip-settings/apply', {});
      applyResult = res;
      message = res.message || 'Core SIP & RTP stack reloaded successfully';
      messageType = 'success';
    } catch (e: any) {
      message = e.message || 'Failed to apply configuration';
      messageType = 'error';
    } finally {
      applying = false;
    }
  }

  async function handleResetConfirm() {
    resetting = true;
    message = '';
    try {
      const res = await api.post<any>('/sip-settings/reset', {});
      if (res && res.settings) {
        settings = { ...res.settings };
      }
      showResetModal = false;
      message = 'Settings restored to factory defaults.';
      messageType = 'success';
    } catch (e: any) {
      message = e.message || 'Failed to reset settings';
      messageType = 'error';
    } finally {
      resetting = false;
    }
  }

  onMount(loadSettings);

  $: rtpPortSpan = Math.max(0, settings.rtp_port_end - settings.rtp_port_start);
  $: rtpSimultaneousCalls = Math.floor(rtpPortSpan / 4);
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
  <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
    <PageHeader
      title="Advanced SIP & RTP Settings"
      subtitle="Complete low-level administration of PBX signaling ports, RTP media port allocation, adaptive jitter buffer, NAT/ICE traversal, session timers, and QoS DSCP markings."
    />
    <div class="flex flex-wrap items-center gap-3">
      <Button variant="secondary" on:click={() => showResetModal = true} disabled={loading || saving || applying}>
        Reset to Defaults
      </Button>
      <Button variant="secondary" on:click={handleApply} disabled={loading || saving || applying}>
        {#if applying}
          <span class="inline-block animate-spin mr-1.5">&#x21BB;</span> Reloading...
        {:else}
          Apply & Reload Stack
        {/if}
      </Button>
      <Button variant="primary" on:click={handleSave} disabled={loading || saving || applying}>
        {#if saving}
          <span class="inline-block animate-spin mr-1.5">&#x21BB;</span> Saving...
        {:else}
          Save Settings
        {/if}
      </Button>
    </div>
  </div>

  {#if message}
    <div class="p-4 rounded-xl text-sm flex items-start gap-3 transition-all {messageType === 'success' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">
      <span class="text-lg leading-none">{messageType === 'success' ? '✓' : '⚠️'}</span>
      <div class="flex-1 font-medium">{message}</div>
    </div>
  {/if}

  {#if applyResult && applyResult.modules}
    <div class="p-4 rounded-xl bg-slate-900 border border-emerald-500/30 text-xs text-slate-300 space-y-2">
      <div class="flex items-center justify-between text-emerald-400 font-semibold uppercase tracking-wider text-[11px]">
        <span>Telephony Stack Reload Telemetry</span>
        <span>{applyResult.applied_at || 'Now'}</span>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-2 pt-1">
        {#each applyResult.modules as mod}
          <div class="bg-slate-800/80 rounded-lg p-2.5 border border-slate-700/60">
            <div class="text-slate-400 font-mono text-[10px] truncate">{mod.module}</div>
            <div class="text-emerald-400 font-semibold mt-0.5 capitalize">{mod.status}</div>
            <div class="text-slate-500 text-[10px] mt-0.5 truncate">{mod.port || mod.port_range || mod.mode || mod.audio_ef}</div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Tab Navigation -->
  <div class="flex flex-wrap gap-2 border-b border-slate-800 pb-3">
    <button
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'signaling' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-slate-800'}"
      on:click={() => activeTab = 'signaling'}
    >
      Signaling & Ports
    </button>
    <button
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'rtp' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-slate-800'}"
      on:click={() => activeTab = 'rtp'}
    >
      RTP Media & Jitter Buffer
    </button>
    <button
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'nat' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-slate-800'}"
      on:click={() => activeTab = 'nat'}
    >
      NAT, STUN, TURN & ICE
    </button>
    <button
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'timers' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-slate-800'}"
      on:click={() => activeTab = 'timers'}
    >
      SIP Timers & TTL
    </button>
    <button
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'security' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-slate-800'}"
      on:click={() => activeTab = 'security'}
    >
      Security & SRTP
    </button>
    <button
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'qos' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-slate-800'}"
      on:click={() => activeTab = 'qos'}
    >
      QoS & DSCP
    </button>
  </div>

  {#if loading}
    <div class="p-12 text-center text-slate-500 flex flex-col items-center justify-center gap-3">
      <div class="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
      <p class="text-sm">Reading core PBX and SIP parameters...</p>
    </div>
  {:else}
    <!-- SECTION 1: SIGNALING & PORTS -->
    {#if activeTab === 'signaling'}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">SIP Engine & Identification</h3>
              <p class="text-xs text-slate-400">Core SIP stack identification and protocol header treatment.</p>
            </div>

            <div>
              <label for="f-user-agent" class="block text-xs font-medium text-slate-300 mb-1">User Agent Header</label>
              <input
                id="f-user-agent"
                type="text"
                bind:value={settings.user_agent}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                placeholder="e.g. MSIP-PBX-Core/2.0"
              />
              <span class="text-[11px] text-slate-500">Advertised in the SIP <code>User-Agent</code> header on dialogs.</span>
            </div>

            <div>
              <label for="f-tcp-keepalive" class="block text-xs font-medium text-slate-300 mb-1">TCP Keep-Alive Interval (Seconds)</label>
              <input
                id="f-tcp-keepalive"
                type="number"
                min="5"
                max="3600"
                bind:value={settings.tcp_keepalive_sec}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
              <span class="text-[11px] text-slate-500">Probing interval to prevent firewall state table eviction on persistent TCP trunks.</span>
            </div>

            <div class="pt-2 space-y-3">
              <label class="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={settings.compact_headers}
                  class="rounded bg-slate-900 border-slate-700 text-indigo-600 focus:ring-0 focus:ring-offset-0 w-4 h-4"
                />
                <div>
                  <div class="text-xs font-medium text-slate-200">Compact SIP Headers</div>
                  <div class="text-[11px] text-slate-500">Use single-letter abbreviations (e.g. <code>v</code> instead of <code>Via</code>, <code>f</code> for <code>From</code>) for bandwidth reduction.</div>
                </div>
              </label>

              <label class="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={settings.allow_contact_rewrite}
                  class="rounded bg-slate-900 border-slate-700 text-indigo-600 focus:ring-0 focus:ring-offset-0 w-4 h-4"
                />
                <div>
                  <div class="text-xs font-medium text-slate-200">Auto-Rewrite Contact Header (NAT Fix)</div>
                  <div class="text-[11px] text-slate-500">Replaces private RFC1918 IPs in the <code>Contact</code> header with actual source socket IP/port.</div>
                </div>
              </label>
            </div>
          </div>
        </Card>

        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">Signaling & Protocol Ports</h3>
              <p class="text-xs text-slate-400">Listening interfaces for local extensions, carriers, and WebRTC clients.</p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-sip-port" class="block text-xs font-medium text-slate-300 mb-1">Internal SIP Port (UDP/TCP)</label>
                <input
                  id="f-sip-port"
                  type="number"
                  bind:value={settings.sip_port}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">Default: <code>5060</code></span>
              </div>

              <div>
                <label for="f-sip-tls-port" class="block text-xs font-medium text-slate-300 mb-1">Internal SIP TLS Port</label>
                <input
                  id="f-sip-tls-port"
                  type="number"
                  bind:value={settings.sip_tls_port}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">Default: <code>5061</code></span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-ext-sip-port" class="block text-xs font-medium text-slate-300 mb-1">External SIP Port (Trunks)</label>
                <input
                  id="f-ext-sip-port"
                  type="number"
                  bind:value={settings.external_sip_port}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">Default: <code>5080</code></span>
              </div>

              <div>
                <label for="f-ext-tls-port" class="block text-xs font-medium text-slate-300 mb-1">External TLS Port</label>
                <input
                  id="f-ext-tls-port"
                  type="number"
                  bind:value={settings.external_sip_tls_port}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">Default: <code>5081</code></span>
              </div>
            </div>

            <div>
              <label for="f-wss-port" class="block text-xs font-medium text-slate-300 mb-1">WebRTC Secure WebSocket (WSS) Port</label>
              <input
                id="f-wss-port"
                type="number"
                bind:value={settings.wss_port}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
              <span class="text-[11px] text-slate-500">Listening port for browser softphones & WebRTC web clients. Default: <code>7443</code>.</span>
            </div>
          </div>
        </Card>
      </div>
    {/if}

    <!-- SECTION 2: RTP MEDIA & JITTER BUFFER -->
    {#if activeTab === 'rtp'}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">RTP Media Port Range</h3>
              <p class="text-xs text-slate-400">Dynamic UDP port pool reserved for audio and video media streams.</p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-rtp-start" class="block text-xs font-medium text-slate-300 mb-1">RTP Port Start</label>
                <input
                  id="f-rtp-start"
                  type="number"
                  bind:value={settings.rtp_port_start}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
              </div>

              <div>
                <label for="f-rtp-end" class="block text-xs font-medium text-slate-300 mb-1">RTP Port End</label>
                <input
                  id="f-rtp-end"
                  type="number"
                  bind:value={settings.rtp_port_end}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
              </div>
            </div>

            <div class="p-3 bg-slate-900/80 border border-slate-700 rounded-lg flex items-center justify-between text-xs">
              <div>
                <span class="text-slate-400">Total Allocated Ports:</span>
                <span class="font-mono text-indigo-400 font-semibold ml-1.5">{rtpPortSpan.toLocaleString()}</span>
              </div>
              <div>
                <span class="text-slate-400">Max Concurrent Calls:</span>
                <span class="font-mono text-emerald-400 font-semibold ml-1.5">~{rtpSimultaneousCalls.toLocaleString()}</span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-rtp-timeout" class="block text-xs font-medium text-slate-300 mb-1">RTP Media Timeout (Sec)</label>
                <input
                  id="f-rtp-timeout"
                  type="number"
                  bind:value={settings.rtp_timeout_sec}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">Hang up call if no RTP packets received.</span>
              </div>

              <div>
                <label for="f-rtp-hold-timeout" class="block text-xs font-medium text-slate-300 mb-1">RTP Hold Timeout (Sec)</label>
                <input
                  id="f-rtp-hold-timeout"
                  type="number"
                  bind:value={settings.rtp_hold_timeout_sec}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">Maximum allowed duration for call on hold.</span>
              </div>
            </div>

            <label class="flex items-center gap-3 cursor-pointer pt-2">
              <input
                type="checkbox"
                bind:checked={settings.enable_cng}
                class="rounded bg-slate-900 border-slate-700 text-indigo-600 focus:ring-0 focus:ring-offset-0 w-4 h-4"
              />
              <div>
                <div class="text-xs font-medium text-slate-200">Comfort Noise Generation (CNG) & Silence Suppression</div>
                <div class="text-[11px] text-slate-500">Suppresses voice bandwidth transmission during speech pauses.</div>
              </div>
            </label>
          </div>
        </Card>

        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">Adaptive Jitter Buffer & Codecs</h3>
              <p class="text-xs text-slate-400">Audio packet smoothing algorithm, DTMF relay, and priority codec negotiation.</p>
            </div>

            <div>
              <label for="f-jitter-mode" class="block text-xs font-medium text-slate-300 mb-1">Jitter Buffer Engine</label>
              <select
                id="f-jitter-mode"
                bind:value={settings.jitter_buffer_mode}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
              >
                <option value="adaptive">Adaptive Dynamic Buffer (Recommended)</option>
                <option value="fixed">Fixed Millisecond Delay</option>
                <option value="disabled">Disabled (Direct RTP Pass-Through)</option>
              </select>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-jitter-initial" class="block text-xs font-medium text-slate-300 mb-1">Initial Delay (ms)</label>
                <input
                  id="f-jitter-initial"
                  type="number"
                  bind:value={settings.jitter_buffer_initial_ms}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
              </div>

              <div>
                <label for="f-jitter-max" class="block text-xs font-medium text-slate-300 mb-1">Maximum Delay (ms)</label>
                <input
                  id="f-jitter-max"
                  type="number"
                  bind:value={settings.jitter_buffer_max_ms}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-dtmf-mode" class="block text-xs font-medium text-slate-300 mb-1">DTMF Signaling Mode</label>
                <select
                  id="f-dtmf-mode"
                  bind:value={settings.dtmf_mode}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
                >
                  <option value="rfc2833">RFC 2833 / RFC 4733 (Default)</option>
                  <option value="info">SIP INFO Messages</option>
                  <option value="inband">In-Band Audio Tones</option>
                  <option value="auto">Auto-Detect</option>
                </select>
              </div>

              <div>
                <label for="f-dtmf-duration" class="block text-xs font-medium text-slate-300 mb-1">DTMF Tone Duration (ms)</label>
                <input
                  id="f-dtmf-duration"
                  type="number"
                  bind:value={settings.dtmf_duration_ms}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
              </div>
            </div>

            <div>
              <label for="f-audio-codecs" class="block text-xs font-medium text-slate-300 mb-1">Preferred Audio Codecs (CSV)</label>
              <input
                id="f-audio-codecs"
                type="text"
                bind:value={settings.audio_codecs}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
              <span class="text-[11px] text-slate-500">Ordered by negotiation preference. E.g. <code>OPUS,G722,PCMA,PCMU</code></span>
            </div>

            <div>
              <label for="f-video-codecs" class="block text-xs font-medium text-slate-300 mb-1">Preferred Video Codecs (CSV)</label>
              <input
                id="f-video-codecs"
                type="text"
                bind:value={settings.video_codecs}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
              <span class="text-[11px] text-slate-500">E.g. <code>VP8,H264</code></span>
            </div>
          </div>
        </Card>
      </div>
    {/if}

    <!-- SECTION 3: NAT, STUN, TURN & ICE -->
    {#if activeTab === 'nat'}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">NAT Mode & Public Interfaces</h3>
              <p class="text-xs text-slate-400">Configure how the PBX identifies external IPs behind cloud NAT gateways.</p>
            </div>

            <div>
              <label for="f-nat-mode" class="block text-xs font-medium text-slate-300 mb-1">NAT Traversal Mode</label>
              <select
                id="f-nat-mode"
                bind:value={settings.nat_mode}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
              >
                <option value="auto">Auto-Detect External IP (Via STUN / AWS / GCP Metadata)</option>
                <option value="manual">Manual Explicit IP Specification</option>
                <option value="stun">Dynamic STUN Server Discovery</option>
                <option value="disabled">Disabled (Direct Public IP / Non-NATed)</option>
              </select>
            </div>

            <div>
              <label for="f-ext-ipv4" class="block text-xs font-medium text-slate-300 mb-1">External Public IPv4 Address</label>
              <input
                id="f-ext-ipv4"
                type="text"
                bind:value={settings.external_ipv4}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                placeholder="Leave blank for auto-detection"
              />
              <span class="text-[11px] text-slate-500">Public IP injected into SIP <code>c=</code> SDP line and <code>Contact</code> header.</span>
            </div>

            <div>
              <label for="f-ext-ipv6" class="block text-xs font-medium text-slate-300 mb-1">External Public IPv6 Address</label>
              <input
                id="f-ext-ipv6"
                type="text"
                bind:value={settings.external_ipv6}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                placeholder="Optional IPv6"
              />
            </div>

            <div>
              <label for="f-local-nets" class="block text-xs font-medium text-slate-300 mb-1">Local Networks CIDR (ACL Bypass)</label>
              <textarea
                id="f-local-nets"
                rows="3"
                bind:value={settings.local_networks}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-indigo-500 font-mono"
              ></textarea>
              <span class="text-[11px] text-slate-500">Comma-separated local subnets treated as LAN without NAT traversal.</span>
            </div>
          </div>
        </Card>

        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">STUN, TURN & ICE Negotiation</h3>
              <p class="text-xs text-slate-400">WebRTC and remote SIP handset firewall punching infrastructure.</p>
            </div>

            <div>
              <label for="f-stun-server" class="block text-xs font-medium text-slate-300 mb-1">STUN Server URI</label>
              <input
                id="f-stun-server"
                type="text"
                bind:value={settings.stun_server}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                placeholder="stun:stun.l.google.com:19302"
              />
            </div>

            <div>
              <label for="f-turn-server" class="block text-xs font-medium text-slate-300 mb-1">TURN Relay Server URI</label>
              <input
                id="f-turn-server"
                type="text"
                bind:value={settings.turn_server}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                placeholder="turn:turn.domain.com:3478"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-turn-user" class="block text-xs font-medium text-slate-300 mb-1">TURN Username</label>
                <input
                  id="f-turn-user"
                  type="text"
                  bind:value={settings.turn_user}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
                />
              </div>

              <div>
                <label for="f-turn-pass" class="block text-xs font-medium text-slate-300 mb-1">TURN Secret / Password</label>
                <input
                  id="f-turn-pass"
                  type="password"
                  bind:value={settings.turn_password}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
                  placeholder="••••••••"
                />
              </div>
            </div>

            <div>
              <label for="f-nat-ping" class="block text-xs font-medium text-slate-300 mb-1">NAT Keep-Alive Ping (Sec)</label>
              <input
                id="f-nat-ping"
                type="number"
                bind:value={settings.nat_ping_interval_sec}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
              <span class="text-[11px] text-slate-500">SIP OPTIONS ping rate sent to registered handsets to maintain UDP pinhole state.</span>
            </div>

            <div class="pt-2 space-y-3">
              <label class="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={settings.enable_ice}
                  class="rounded bg-slate-900 border-slate-700 text-indigo-600 focus:ring-0 focus:ring-offset-0 w-4 h-4"
                />
                <div>
                  <div class="text-xs font-medium text-slate-200">Enable Interactive Connectivity Establishment (ICE)</div>
                  <div class="text-[11px] text-slate-500">Collects host, srflx, and relay candidates for WebRTC & SIP sessions.</div>
                </div>
              </label>

              <label class="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={settings.aggressive_ice}
                  class="rounded bg-slate-900 border-slate-700 text-indigo-600 focus:ring-0 focus:ring-offset-0 w-4 h-4"
                />
                <div>
                  <div class="text-xs font-medium text-slate-200">Aggressive ICE Nomination</div>
                  <div class="text-[11px] text-slate-500">Fast-tracks media candidate pair selection for faster call connect times.</div>
                </div>
              </label>
            </div>
          </div>
        </Card>
      </div>
    {/if}

    <!-- SECTION 4: SIP TIMERS & TTL -->
    {#if activeTab === 'timers'}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">RFC 4028 Session Timers</h3>
              <p class="text-xs text-slate-400">Dead dialog detection and automatic call cleanup via re-INVITE / UPDATE refreshes.</p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="f-sess-expires" class="block text-xs font-medium text-slate-300 mb-1">Session-Expires TTL (Sec)</label>
                <input
                  id="f-sess-expires"
                  type="number"
                  min="90"
                  bind:value={settings.session_expires_sec}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">Default: <code>1800</code> (30 mins)</span>
              </div>

              <div>
                <label for="f-min-se" class="block text-xs font-medium text-slate-300 mb-1">Min-SE Floor (Sec)</label>
                <input
                  id="f-min-se"
                  type="number"
                  min="90"
                  bind:value={settings.min_se_sec}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[11px] text-slate-500">RFC 4028 minimum: <code>90</code></span>
              </div>
            </div>

            <div>
              <label for="f-refresher" class="block text-xs font-medium text-slate-300 mb-1">Session Refresher Preference</label>
              <select
                id="f-refresher"
                bind:value={settings.session_refresher}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
              >
                <option value="auto">Auto-Negotiate (Prefer Caller / UAC)</option>
                <option value="uac">Force User Agent Client (UAC)</option>
                <option value="uas">Force User Agent Server (UAS / PBX)</option>
              </select>
            </div>

            <div>
              <label for="f-max-channels" class="block text-xs font-medium text-slate-300 mb-1">Maximum Concurrent Channels / Dialogs</label>
              <input
                id="f-max-channels"
                type="number"
                min="10"
                max="5000"
                bind:value={settings.max_channels}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
              <span class="text-[11px] text-slate-500">System-wide tenant capacity threshold. Returns 503 Service Unavailable when exceeded.</span>
            </div>
          </div>
        </Card>

        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">Endpoint Registration TTL & Retries</h3>
              <p class="text-xs text-slate-400">Controls registration expiry boundaries for SIP clients and upstream trunks.</p>
            </div>

            <div class="grid grid-cols-3 gap-3">
              <div>
                <label for="f-reg-min" class="block text-xs font-medium text-slate-300 mb-1">Min TTL (Sec)</label>
                <input
                  id="f-reg-min"
                  type="number"
                  min="10"
                  bind:value={settings.register_min_sec}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[10px] text-slate-500">Default: 60</span>
              </div>

              <div>
                <label for="f-reg-def" class="block text-xs font-medium text-slate-300 mb-1">Default TTL</label>
                <input
                  id="f-reg-def"
                  type="number"
                  min="60"
                  bind:value={settings.register_default_sec}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[10px] text-slate-500">Default: 3600</span>
              </div>

              <div>
                <label for="f-reg-max" class="block text-xs font-medium text-slate-300 mb-1">Max TTL</label>
                <input
                  id="f-reg-max"
                  type="number"
                  min="60"
                  bind:value={settings.register_max_sec}
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <span class="text-[10px] text-slate-500">Default: 3600</span>
              </div>
            </div>

            <div>
              <label for="f-reg-retry" class="block text-xs font-medium text-slate-300 mb-1">Trunk Registration Retry Interval (Sec)</label>
              <input
                id="f-reg-retry"
                type="number"
                min="5"
                bind:value={settings.register_retry_sec}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
              <span class="text-[11px] text-slate-500">Backoff interval when an outbound carrier trunk registration fails or returns 408/500.</span>
            </div>
          </div>
        </Card>
      </div>
    {/if}

    <!-- SECTION 5: SECURITY, TLS & SRTP -->
    {#if activeTab === 'security'}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">SRTP Media Encryption</h3>
              <p class="text-xs text-slate-400">Encrypt real-time audio and video packets to prevent eavesdropping.</p>
            </div>

            <div>
              <label for="f-srtp-mode" class="block text-xs font-medium text-slate-300 mb-1">SRTP Encryption Policy</label>
              <select
                id="f-srtp-mode"
                bind:value={settings.srtp_mode}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
              >
                <option value="optional">Optional / Opportunistic (Encrypt when peer supports SRTP)</option>
                <option value="mandatory">Mandatory / Enforced (Reject unencrypted calls)</option>
                <option value="disabled">Disabled (Plain unencrypted RTP only)</option>
              </select>
            </div>

            <div>
              <label for="f-srtp-suites" class="block text-xs font-medium text-slate-300 mb-1">Allowed SRTP Cipher Suites (CSV)</label>
              <textarea
                id="f-srtp-suites"
                rows="3"
                bind:value={settings.srtp_suites}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-indigo-500 font-mono"
              ></textarea>
              <span class="text-[11px] text-slate-500">Prioritized crypto suites. Includes Galois/Counter Mode (GCM) and HMAC-SHA1.</span>
            </div>
          </div>
        </Card>

        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">TLS & SIP Authentication</h3>
              <p class="text-xs text-slate-400">Cryptographic trust requirements and anti-scanning defenses.</p>
            </div>

            <div>
              <label for="f-tls-ver" class="block text-xs font-medium text-slate-300 mb-1">Minimum TLS Protocol Version</label>
              <select
                id="f-tls-ver"
                bind:value={settings.tls_version_min}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500"
              >
                <option value="1.2">TLS 1.2 (Standard compatibility)</option>
                <option value="1.3">TLS 1.3 (Maximum security)</option>
              </select>
            </div>

            <div>
              <label for="f-auth-realm" class="block text-xs font-medium text-slate-300 mb-1">SIP Digest Authentication Realm</label>
              <input
                id="f-auth-realm"
                type="text"
                bind:value={settings.auth_realm}
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
              />
            </div>

            <div class="pt-2 space-y-3">
              <label class="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={settings.tls_verify_peer}
                  class="rounded bg-slate-900 border-slate-700 text-indigo-600 focus:ring-0 focus:ring-offset-0 w-4 h-4"
                />
                <div>
                  <div class="text-xs font-medium text-slate-200">Verify Remote TLS Peer Certificate</div>
                  <div class="text-[11px] text-slate-500">Rejects SIP TLS trunks with self-signed or unverified certificates.</div>
                </div>
              </label>

              <label class="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={settings.challenge_all_invites}
                  class="rounded bg-slate-900 border-slate-700 text-indigo-600 focus:ring-0 focus:ring-offset-0 w-4 h-4"
                />
                <div>
                  <div class="text-xs font-medium text-slate-200">Challenge All Outbound INVITE Dialogs (Anti-Scanner)</div>
                  <div class="text-[11px] text-slate-500">Issues 407 Proxy Authentication challenge to every unauthenticated call attempt.</div>
                </div>
              </label>
            </div>
          </div>
        </Card>
      </div>
    {/if}

    <!-- SECTION 6: QOS & DSCP -->
    {#if activeTab === 'qos'}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">DiffServ DSCP Header Tagging</h3>
              <p class="text-xs text-slate-400">Sets IPv4 ToS / IPv6 Traffic Class bytes so enterprise routers prioritize telephony packets.</p>
            </div>

            <div>
              <label for="f-dscp-audio" class="block text-xs font-medium text-slate-300 mb-1">Audio RTP Stream DSCP</label>
              <div class="flex items-center gap-3">
                <input
                  id="f-dscp-audio"
                  type="number"
                  min="0"
                  max="63"
                  bind:value={settings.dscp_audio}
                  class="w-32 bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <Badge variant={settings.dscp_audio === 46 ? 'success' : 'info'}>
                  {settings.dscp_audio === 46 ? 'Expedited Forwarding (EF - 46)' : `DSCP ${settings.dscp_audio}`}
                </Badge>
              </div>
              <span class="text-[11px] text-slate-500">Standard for VoIP voice packets is <code>46</code> (Expedited Forwarding, lowest latency queue).</span>
            </div>

            <div>
              <label for="f-dscp-video" class="block text-xs font-medium text-slate-300 mb-1">Video Stream DSCP</label>
              <div class="flex items-center gap-3">
                <input
                  id="f-dscp-video"
                  type="number"
                  min="0"
                  max="63"
                  bind:value={settings.dscp_video}
                  class="w-32 bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <Badge variant={settings.dscp_video === 34 ? 'success' : 'info'}>
                  {settings.dscp_video === 34 ? 'Assured Forwarding 41 (AF41 - 34)' : `DSCP ${settings.dscp_video}`}
                </Badge>
              </div>
              <span class="text-[11px] text-slate-500">Standard for video conference packets is <code>34</code> (AF41).</span>
            </div>

            <div>
              <label for="f-dscp-sig" class="block text-xs font-medium text-slate-300 mb-1">SIP Signaling Protocol DSCP</label>
              <div class="flex items-center gap-3">
                <input
                  id="f-dscp-sig"
                  type="number"
                  min="0"
                  max="63"
                  bind:value={settings.dscp_signaling}
                  class="w-32 bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-indigo-500 font-mono"
                />
                <Badge variant={settings.dscp_signaling === 26 ? 'success' : 'info'}>
                  {settings.dscp_signaling === 26 ? 'Assured Forwarding 31 (AF31 - 26)' : `DSCP ${settings.dscp_signaling}`}
                </Badge>
              </div>
              <span class="text-[11px] text-slate-500">Standard for SIP control signaling is <code>26</code> (AF31) or <code>24</code> (CS3).</span>
            </div>
          </div>
        </Card>

        <Card>
          <div class="space-y-4">
            <div class="border-b border-slate-800 pb-3">
              <h3 class="text-base font-semibold text-white">QoS Quick Reference Guide</h3>
              <p class="text-xs text-slate-400">Industry-standard mapping of DiffServ values to router queue policies.</p>
            </div>

            <div class="overflow-x-auto">
              <table class="w-full text-left text-xs">
                <thead>
                  <tr class="border-b border-slate-700/80 text-slate-400">
                    <th class="py-2">Traffic Class</th>
                    <th class="py-2">DSCP Val</th>
                    <th class="py-2">Hex / TOS</th>
                    <th class="py-2">Router Treatment</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-800 text-slate-300 font-mono text-[11px]">
                  <tr>
                    <td class="py-2 text-emerald-400 font-sans font-semibold">Voice (RTP)</td>
                    <td class="py-2">46 (EF)</td>
                    <td class="py-2">0xB8</td>
                    <td class="py-2 font-sans text-slate-400">Strict Priority Queue (LLQ)</td>
                  </tr>
                  <tr>
                    <td class="py-2 text-sky-400 font-sans font-semibold">Interactive Video</td>
                    <td class="py-2">34 (AF41)</td>
                    <td class="py-2">0x88</td>
                    <td class="py-2 font-sans text-slate-400">High-Bandwidth Queue</td>
                  </tr>
                  <tr>
                    <td class="py-2 text-indigo-400 font-sans font-semibold">SIP Signaling</td>
                    <td class="py-2">26 (AF31)</td>
                    <td class="py-2">0x68</td>
                    <td class="py-2 font-sans text-slate-400">Control Plane Queue</td>
                  </tr>
                  <tr>
                    <td class="py-2 text-slate-400 font-sans font-semibold">Best Effort</td>
                    <td class="py-2">0 (BE)</td>
                    <td class="py-2">0x00</td>
                    <td class="py-2 font-sans text-slate-400">Standard Data Queue</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </Card>
      </div>
    {/if}
  {/if}
</div>

<!-- Modal confirmation for Reset to Defaults -->
<Modal
  isOpen={showResetModal}
  title="Reset SIP & RTP Settings"
  onClose={() => showResetModal = false}
>
  <div class="space-y-4">
    <p class="text-sm text-slate-300">
      Are you sure you want to restore all PBX, SIP, RTP, and QoS settings to system factory defaults?
    </p>
    <div class="p-3 bg-amber-500/10 border border-amber-500/20 rounded-lg text-xs text-amber-300">
      ⚠️ Any custom SIP ports, RTP ranges, or STUN/TURN credentials will be reverted to standard defaults (5060, 16384–32768, adaptive jitter buffer).
    </div>
    <div class="flex justify-end gap-3 pt-2">
      <Button variant="secondary" on:click={() => showResetModal = false} disabled={resetting}>
        Cancel
      </Button>
      <Button variant="danger" on:click={handleResetConfirm} disabled={resetting}>
        {#if resetting}
          Restoring...
        {:else}
          Yes, Restore Defaults
        {/if}
      </Button>
    </div>
  </div>
</Modal>
