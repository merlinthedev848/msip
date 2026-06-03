<script lang="ts">
  import { slide, fade } from 'svelte/transition';
  import { onMount, onDestroy } from 'svelte';
  import { UserAgent, Registerer, Inviter, SessionState } from 'sip.js';

  let isOpen = false;
  let currentNumber = "";
  let activeCall = null; // null | 'calling' | 'connected' | 'ringing'
  let callDuration = 0;
  let callTimer: ReturnType<typeof setInterval>;

  let remoteAudio: HTMLAudioElement;
  let userAgent: UserAgent;
  let registerer: Registerer;
  let session: any = null;

  // Configure WebRTC details
  let sipExt = "1001";
  let sipPass = "password"; // Default for demo
  let sipDomain = "localhost";
  let wsServer = "ws://localhost:5066"; // FreeSWITCH native WS port
  let isRegistered = false;

  const toggle = () => {
    isOpen = !isOpen;
    if (isOpen && !isRegistered) {
      connectSip();
    }
  };

  const appendDigit = (digit) => {
    if (currentNumber.length < 15) currentNumber += digit;
  };

  const backspace = () => {
    currentNumber = currentNumber.slice(0, -1);
  };

  const connectSip = async () => {
    const uri = UserAgent.makeURI(`sip:${sipExt}@${sipDomain}`);
    if (!uri) return;

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

    try {
      await userAgent.start();
      registerer = new Registerer(userAgent);
      await registerer.register();
      isRegistered = true;
    } catch (error) {
      console.error("SIP Registration Failed:", error);
    }
  };

  onDestroy(() => {
    if (userAgent) userAgent.stop();
  });

  const dial = async () => {
    if (!currentNumber || !userAgent) return;
    activeCall = 'calling';
    
    const target = UserAgent.makeURI(`sip:${currentNumber}@${sipDomain}`);
    if (!target) return;

    const sdhOptions = {
      constraints: { audio: true, video: false },
      peerConnectionConfiguration: {
        iceServers: [
          { urls: 'stun:localhost:3478' },
          { urls: 'turn:localhost:3478', username: 'pbxuser', credential: 'pbxpass' }
        ]
      }
    };

    session = new Inviter(userAgent, target, {
      sessionDescriptionHandlerOptions: sdhOptions
    });

    session.stateChange.addListener((state: SessionState) => {
      if (state === SessionState.Established) {
        activeCall = 'connected';
        callDuration = 0;
        callTimer = setInterval(() => callDuration++, 1000);
        
        const sdh = session.sessionDescriptionHandler;
        const stream = new MediaStream();
        sdh.peerConnection.getReceivers().forEach((receiver) => {
          if (receiver.track) stream.addTrack(receiver.track);
        });
        if (remoteAudio) {
          remoteAudio.srcObject = stream;
          remoteAudio.play();
        }
      } else if (state === SessionState.Terminated) {
        hangup();
      }
    });

    try {
      await session.invite();
    } catch (e) {
      console.error("Invite failed", e);
      hangup();
    }
  };

  const answer = async () => {
    if (!session) return;
    try {
      const sdhOptions = {
        constraints: { audio: true, video: false },
        peerConnectionConfiguration: {
          iceServers: [
            { urls: 'stun:localhost:3478' },
            { urls: 'turn:localhost:3478', username: 'pbxuser', credential: 'pbxpass' }
          ]
        }
      };

      await session.accept({
        sessionDescriptionHandlerOptions: sdhOptions
      });
    } catch (e) {
      console.error(e);
    }
  };

  const hangup = () => {
    if (session) {
      if (session.state === SessionState.Established) {
        session.bye();
      } else if (session.state === SessionState.Initial || session.state === SessionState.Establishing) {
        if (session instanceof Inviter) session.cancel();
        else session.reject();
      }
    }
    session = null;
    activeCall = null;
    currentNumber = "";
    clearInterval(callTimer);
  };

  $: formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };
</script>

<audio bind:this={remoteAudio} autoplay style="display:none"></audio>

<div class="fixed bottom-6 right-6 z-50 flex flex-col items-end">
  
  {#if isOpen}
    <div 
      transition:slide={{duration: 300}}
      class="bg-gray-900/95 backdrop-blur-xl border border-gray-700/50 rounded-3xl shadow-2xl overflow-hidden mb-4 w-72 flex flex-col"
    >
      <!-- Header -->
      <div class="bg-gray-950 px-4 py-3 border-b border-gray-800 flex justify-between items-center">
        <div class="flex items-center space-x-2">
          <div class="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.8)]"></div>
          <span class="text-xs font-bold text-gray-300 uppercase tracking-widest">Ext 1001</span>
        </div>
        <button on:click={toggle} aria-label="Close Softphone" class="text-gray-500 hover:text-white transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      {#if !activeCall}
        <!-- Display -->
        <div class="p-4 bg-gray-900 flex justify-between items-center h-20">
          <h2 class="text-3xl font-mono text-white tracking-widest truncate">{currentNumber || '...'}</h2>
          {#if currentNumber}
            <button on:click={backspace} aria-label="Backspace" class="text-gray-500 hover:text-rose-400 p-2 transition-colors">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M3 12l6.414 6.414a2 2 0 001.414.586H19a2 2 0 002-2V7a2 2 0 00-2-2h-8.172a2 2 0 00-1.414.586L3 12z"></path></svg>
            </button>
          {/if}
        </div>

        <!-- Dialpad -->
        <div class="p-4 grid grid-cols-3 gap-3 bg-gray-900/50">
          {#each [
            {n:'1', l:''}, {n:'2', l:'ABC'}, {n:'3', l:'DEF'},
            {n:'4', l:'GHI'}, {n:'5', l:'JKL'}, {n:'6', l:'MNO'},
            {n:'7', l:'PQRS'}, {n:'8', l:'TUV'}, {n:'9', l:'WXYZ'},
            {n:'*', l:''}, {n:'0', l:'+'}, {n:'#', l:''}
          ] as key}
            <button 
              on:click={() => appendDigit(key.n)}
              class="h-14 rounded-2xl bg-gray-800 hover:bg-gray-700 active:bg-indigo-600 transition-colors flex flex-col items-center justify-center border border-gray-700/50"
            >
              <span class="text-xl font-bold text-white leading-none">{key.n}</span>
              <span class="text-[9px] font-bold text-gray-500 mt-1 uppercase">{key.l}</span>
            </button>
          {/each}
        </div>

        <!-- Action Bar -->
        <div class="p-4 pb-6 flex justify-center space-x-6 border-t border-gray-800 bg-gray-900/80">
          <button aria-label="Call History" class="w-12 h-12 rounded-full bg-gray-800 text-gray-400 hover:text-white flex items-center justify-center transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
          </button>
          <button 
            on:click={dial}
            aria-label="Dial"
            class="w-16 h-16 rounded-full bg-emerald-500 hover:bg-emerald-400 text-white flex items-center justify-center shadow-lg shadow-emerald-500/30 transform -translate-y-2 transition-transform active:scale-95"
          >
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
          </button>
          <button aria-label="Contacts" class="w-12 h-12 rounded-full bg-gray-800 text-gray-400 hover:text-white flex items-center justify-center transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
          </button>
        </div>
      {:else}
        <!-- Active Call View -->
        <div in:fade class="flex-1 flex flex-col p-6 items-center justify-center bg-gray-900/80 min-h-[300px]">
          {#if activeCall === 'calling'}
            <div class="w-16 h-16 rounded-full border-4 border-indigo-500/20 border-t-indigo-500 animate-spin mb-6"></div>
            <h3 class="text-lg font-bold text-white mb-1">Calling...</h3>
            <p class="text-sm font-mono text-gray-400">{currentNumber}</p>
          {:else}
            <div class="w-20 h-20 rounded-full bg-gradient-to-br from-indigo-500 to-violet-600 flex items-center justify-center mb-6 shadow-lg shadow-indigo-500/30">
              <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-1">Unknown Contact</h3>
            <p class="text-sm font-mono text-indigo-400 mb-6">{currentNumber}</p>
            <div class="px-4 py-1.5 bg-gray-950 rounded-full text-white font-mono text-xl tracking-widest border border-gray-800 shadow-inner">
              {formatTime(callDuration)}
            </div>
          {/if}
          
          <div class="mt-8 flex space-x-4">
            <button aria-label="Mute" class="w-12 h-12 rounded-full bg-gray-800 text-white flex items-center justify-center hover:bg-gray-700">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"></path></svg>
            </button>
            <button on:click={hangup} aria-label="Hangup" class="w-16 h-16 rounded-full bg-rose-500 hover:bg-rose-400 text-white flex items-center justify-center shadow-lg shadow-rose-500/30 transform transition-transform active:scale-95">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 8l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M5 3a2 2 0 00-2 2v1c0 8.284 6.716 15 15 15h1a2 2 0 002-2v-3.28a1 1 0 00-.684-.948l-4.493-1.498a1 1 0 00-1.21.502l-1.13 2.257a11.042 11.042 0 01-5.516-5.516l2.257-1.13a1 1 0 00.502-1.21L9.22 3.684A1 1 0 008.27 3H5z"></path></svg>
            </button>
            <button aria-label="Transfer" class="w-12 h-12 rounded-full bg-gray-800 text-white flex items-center justify-center hover:bg-gray-700">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9V5a2 2 0 012-2h2a2 2 0 012 2v4M14 11h-4m4 0h4m-4 0v4m-4-4v4m0 4h4m-4 0h-4m0 0v-4"></path></svg>
            </button>
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Toggle Button (Always visible) -->
  <button 
    on:click={toggle}
    class="w-16 h-16 rounded-full bg-indigo-600 hover:bg-indigo-500 text-white shadow-xl shadow-indigo-600/30 flex items-center justify-center transition-all transform hover:scale-105 active:scale-95"
  >
    {#if isOpen}
      <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
    {:else}
      <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
    {/if}
  </button>
</div>
