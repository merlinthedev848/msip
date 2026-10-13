<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { UserAgent, Registerer } from 'sip.js';

  let isOpen = false;
  let showConfig = false;
  let currentNumber = "";
  let activeCall: 'calling' | 'connected' | 'ringing' | null = null;
  let callDuration = 0;
  let callTimer: any;

  let remoteAudio: HTMLAudioElement;
  let userAgent: UserAgent | null = null;
  let registerer: Registerer | null = null;
  let session: any = null;

  let sipExt = localStorage.getItem('pbx_softphone_ext') || "";
  let sipPass = localStorage.getItem('pbx_softphone_pass') || "";
  let isRegistered = false;
  let registrationStatus = "Not configured";

  function saveConfig() {
    localStorage.setItem('pbx_softphone_ext', sipExt.trim());
    localStorage.setItem('pbx_softphone_pass', sipPass.trim());
    showConfig = false;
    connectSip();
  }

  const toggle = () => {
    isOpen = !isOpen;
    if (isOpen && !isRegistered && sipExt && sipPass) {
      connectSip();
    }
  };

  const appendDigit = (digit: string) => {
    if (currentNumber.length < 20) currentNumber += digit;
  };

  const backspace = () => {
    currentNumber = currentNumber.slice(0, -1);
  };

  const connectSip = async () => {
    if (!sipExt || !sipPass) {
      registrationStatus = "Configure extension credentials";
      return;
    }

    try {
      registrationStatus = "Registering...";
      const sipDomain = window.location.hostname;
      const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsPort = window.location.protocol === 'https:' ? '7443' : '5066';
      const wsServer = `${wsProtocol}//${window.location.hostname}:${wsPort}`;

      const uri = UserAgent.makeURI(`sip:${sipExt}@${sipDomain}`);
      if (!uri) {
        registrationStatus = "Invalid SIP URI";
        return;
      }

      userAgent = new UserAgent({
        uri: uri,
        transportOptions: { server: wsServer },
        authorizationUsername: sipExt,
        authorizationPassword: sipPass,
        delegate: {
          onInvite(invitation) {
            session = invitation;
            activeCall = 'ringing';
            currentNumber = invitation.remoteIdentity.uri.user || 'Unknown';
          }
        }
      });

      await userAgent.start();
      registerer = new Registerer(userAgent);
      registerer.stateChange.addListener((newState) => {
        if (newState === 'Registered') {
          isRegistered = true;
          registrationStatus = `Registered (Ext ${sipExt})`;
        } else if (newState === 'Unregistered') {
          isRegistered = false;
          registrationStatus = 'Unregistered';
        } else if (newState === 'Terminated') {
          isRegistered = false;
          registrationStatus = 'Registration failed';
        }
      });
      await registerer.register();
    } catch (e: any) {
      isRegistered = false;
      registrationStatus = e.message || 'SIP connection error';
    }
  };

  const hangupCall = () => {
    if (session) {
      session.bye();
      session = null;
    }
    activeCall = null;
    clearInterval(callTimer);
    callDuration = 0;
  };

  onMount(() => {
    if (sipExt && sipPass) {
      registrationStatus = "Ready to connect";
    }
  });

  onDestroy(() => {
    if (registerer) registerer.unregister();
    if (userAgent) userAgent.stop();
    if (callTimer) clearInterval(callTimer);
  });
</script>

<audio bind:this={remoteAudio} autoplay></audio>

<!-- Floating Launcher Button Bottom-Right -->
<div class="fixed bottom-6 right-6 z-40">
  <button 
    on:click={toggle}
    aria-label="Toggle Softphone"
    class="w-14 h-14 rounded-full bg-[#0284C7] hover:bg-sky-700 text-white shadow-lg flex items-center justify-center transition-all duration-200 focus:outline-none focus:ring-4 focus:ring-sky-200"
  >
    {#if isOpen}
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
    {:else}
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
    {/if}
  </button>
</div>

<!-- Softphone Panel -->
{#if isOpen}
  <div class="fixed bottom-24 right-6 z-40 w-80 bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans">
    
    <!-- Top Header -->
    <div class="px-4 py-3 bg-slate-50 border-b border-slate-200 flex items-center justify-between">
      <div class="flex items-center space-x-2">
        <span class="w-2.5 h-2.5 rounded-full {isRegistered ? 'bg-emerald-500' : 'bg-slate-400'}"></span>
        <span class="text-xs font-bold text-slate-800">WebRTC Softphone</span>
      </div>
      <button 
        type="button"
        on:click={() => showConfig = !showConfig}
        class="text-xs text-sky-700 font-semibold hover:underline"
      >
        {showConfig ? 'Keypad' : 'Settings'}
      </button>
    </div>

    <!-- Status Bar -->
    <div class="px-4 py-1.5 bg-slate-100 text-[11px] text-slate-600 border-b border-slate-200 flex items-center justify-between">
      <span class="truncate">{registrationStatus}</span>
      {#if !isRegistered && !showConfig}
        <button on:click={connectSip} class="text-sky-700 font-bold hover:underline">Connect</button>
      {/if}
    </div>

    {#if showConfig}
      <!-- Settings Panel -->
      <form on:submit|preventDefault={saveConfig} class="p-4 space-y-3 text-xs">
        <div>
          <label for="sp-ext" class="block font-medium text-slate-700 mb-1">SIP Extension</label>
          <input id="sp-ext" class="w-full bg-white border border-slate-300 rounded px-2.5 py-1.5 mono" bind:value={sipExt} placeholder="e.g. 1001" required />
        </div>
        <div>
          <label for="sp-pass" class="block font-medium text-slate-700 mb-1">SIP Password</label>
          <input id="sp-pass" type="password" class="w-full bg-white border border-slate-300 rounded px-2.5 py-1.5 mono" bind:value={sipPass} placeholder="Secret" required />
        </div>
        <button type="submit" class="w-full py-2 bg-[#0284C7] hover:bg-sky-700 text-white font-bold rounded mt-2">
          Save & Register
        </button>
      </form>
    {:else}
      <!-- Dialpad Area -->
      <div class="p-4 space-y-3">
        <!-- Number Input Display -->
        <div class="flex items-center bg-slate-50 border border-slate-200 rounded-lg px-3 py-2">
          <input 
            type="text" 
            bind:value={currentNumber} 
            placeholder="Dial number or ext..." 
            class="bg-transparent border-none text-slate-900 font-mono text-base font-bold w-full focus:outline-none"
          />
          {#if currentNumber}
            <button on:click={backspace} class="text-slate-400 hover:text-slate-700 p-1" title="Backspace">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M3 12l7-7h10a2 2 0 012 2v10a2 2 0 01-2 2H10l-7-7z" /></svg>
            </button>
          {/if}
        </div>

        <!-- 12-button Keypad -->
        <div class="grid grid-cols-3 gap-2">
          {#each ['1', '2', '3', '4', '5', '6', '7', '8', '9', '*', '0', '#'] as d}
            <button 
              type="button"
              on:click={() => appendDigit(d)}
              class="h-10 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold rounded-lg transition-colors flex items-center justify-center text-sm font-mono"
            >
              {d}
            </button>
          {/each}
        </div>

        <!-- Call Actions -->
        <div class="pt-2">
          {#if activeCall}
            <button 
              type="button"
              on:click={hangupCall}
              class="w-full py-2.5 bg-rose-600 hover:bg-rose-700 text-white font-bold rounded-xl shadow-md flex items-center justify-center space-x-2 text-xs"
            >
              <span>Hang Up</span>
            </button>
          {:else}
            <button 
              type="button"
              class="w-full py-2.5 bg-[#0284C7] hover:bg-sky-700 text-white font-bold rounded-xl shadow-md flex items-center justify-center space-x-2 text-xs"
            >
              <span>Call</span>
            </button>
          {/if}
        </div>
      </div>
    {/if}

  </div>
{/if}
