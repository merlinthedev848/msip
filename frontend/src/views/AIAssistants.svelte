<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime, formatDuration } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface AIAssistant {
    ID: string;
    Name: string;
    ExtensionNumber: string;
    RoleDescription: string;
    SystemPrompt: string;
    GreetingMessage: string;
    LLMProvider: string;
    LLMModel: string;
    TTSProvider: string;
    VoiceID: string;
    EnableBooking: boolean;
    EnableTransfer: boolean;
    TransferTarget: string;
    EnableVoicemail: boolean;
    MaxCallDurationSec: number;
    IsActive: boolean;
    CreatedAt: string;
  }

  interface AICallSession {
    ID: string;
    AssistantID: string;
    AssistantName: string;
    CallUUID: string;
    CallerNumber: string;
    Summary: string;
    Sentiment: string;
    ActionTaken: string;
    MessagesJSON: string;
    DurationSec: number;
    CreatedAt: string;
  }

  interface AITranscriptionSettings {
    AutoTranscribeVoicemail: boolean;
    AutoTranscribeRecordings: boolean;
    STTProvider: string;
    STTModel: string;
    STTAPIKey: string;
    TTSProvider: string;
    TTSAPIKey: string;
    TTSDefaultVoice: string;
    LLMProvider: string;
    LLMAPIKey: string;
    LLMDefaultModel: string;
    LLMBaseURL: string;
    SentimentAnalysis: boolean;
    EmailTranscripts: boolean;
    NotificationEmail: string;
  }

  interface SimMessage {
    role: string;
    content: string;
    time: string;
  }

  interface ToolCall {
    name: string;
    arguments: Record<string, any>;
    result: string;
  }

  let activeTab: 'assistants' | 'simulator' | 'settings' | 'sessions' = 'assistants';
  let assistants: AIAssistant[] = [];
  let sessions: AICallSession[] = [];
  let settings: AITranscriptionSettings = {
    AutoTranscribeVoicemail: true,
    AutoTranscribeRecordings: true,
    STTProvider: 'whisper',
    STTModel: 'whisper-1',
    STTAPIKey: '',
    TTSProvider: 'openai',
    TTSAPIKey: '',
    TTSDefaultVoice: 'alloy',
    LLMProvider: 'openai',
    LLMAPIKey: '',
    LLMDefaultModel: 'gpt-4o',
    LLMBaseURL: 'https://api.openai.com/v1',
    SentimentAnalysis: true,
    EmailTranscripts: false,
    NotificationEmail: '',
  };

  let loading = true;
  let loadError = '';

  // Modal State
  let showModal = false;
  let isEditing = false;
  let formAssistant: AIAssistant = getBlankAssistant();

  // Simulator State
  let simAssistantId = '';
  let simCallerNumber = '+1 (555) 234-5678';
  let simMessages: SimMessage[] = [];
  let simInput = '';
  let simSending = false;
  let simLastTools: ToolCall[] = [];

  // Session Detail Modal
  let selectedSession: AICallSession | null = null;
  let selectedSessionMessages: SimMessage[] = [];

  function getBlankAssistant(): AIAssistant {
    return {
      ID: '',
      Name: 'Sarah - AI Virtual Receptionist',
      ExtensionNumber: '890',
      RoleDescription: 'Virtual Front Desk Receptionist & Booking Agent',
      SystemPrompt: `You are Sarah, a professional, warm, and highly efficient AI virtual receptionist for MSIP Enterprise PBX. 
You answer incoming calls, book appointments, answers questions about our services, and transfer callers to extensions when requested.
Be concise, conversational, and courteous. Keep spoken responses to 1-3 sentences.`,
      GreetingMessage: 'Hello, thank you for calling MSIP PBX! My name is Sarah. How can I help you today?',
      LLMProvider: 'openai',
      LLMModel: 'gpt-4o',
      TTSProvider: 'openai',
      VoiceID: 'alloy',
      EnableBooking: true,
      EnableTransfer: true,
      TransferTarget: '1001',
      EnableVoicemail: true,
      MaxCallDurationSec: 600,
      IsActive: true,
      CreatedAt: '',
    };
  }

  async function loadAll() {
    loading = true;
    loadError = '';
    try {
      const [asstRes, sessRes, settRes] = await Promise.all([
        api.get<{ assistants: AIAssistant[] }>('/ai/assistants'),
        api.get<{ sessions: AICallSession[] }>('/ai/sessions'),
        api.get<AITranscriptionSettings>('/ai/settings'),
      ]);
      assistants = asstRes.assistants || [];
      sessions = sessRes.sessions || [];
      if (settRes) settings = settRes;
      if (assistants.length > 0 && !simAssistantId) {
        simAssistantId = assistants[0].ID;
        resetSimulator();
      }
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreateModal() {
    isEditing = false;
    formAssistant = getBlankAssistant();
    showModal = true;
  }

  function openEditModal(asst: AIAssistant) {
    isEditing = true;
    formAssistant = { ...asst };
    showModal = true;
  }

  async function saveAssistant() {
    if (!formAssistant.Name.trim()) {
      toast.error('Assistant name is required');
      return;
    }
    try {
      if (isEditing) {
        await api.put(`/ai/assistants/${formAssistant.ID}`, formAssistant);
        toast.success('AI Assistant updated successfully');
      } else {
        await api.post('/ai/assistants', formAssistant);
        toast.success('AI Assistant created successfully');
      }
      showModal = false;
      await loadAll();
    } catch (e) {
      toastError(e);
    }
  }

  async function deleteAssistant(asst: AIAssistant) {
    if (!confirm(`Delete AI assistant "${asst.Name}"?`)) return;
    try {
      await api.del(`/ai/assistants/${asst.ID}`);
      toast.success('AI Assistant removed');
      await loadAll();
    } catch (e) {
      toastError(e);
    }
  }

  async function saveSettings() {
    try {
      await api.put('/ai/settings', settings);
      toast.success('AI & Transcription settings saved');
    } catch (e) {
      toastError(e);
    }
  }

  function resetSimulator() {
    const sel = assistants.find(a => a.ID === simAssistantId) || assistants[0];
    simMessages = [
      {
        role: 'assistant',
        content: sel ? sel.GreetingMessage : 'Hello, how can I assist you today?',
        time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      },
    ];
    simLastTools = [];
    simInput = '';
  }

  async function sendSimulatorMessage() {
    if (!simInput.trim() || simSending) return;
    const userText = simInput.trim();
    simInput = '';
    simMessages = [
      ...simMessages,
      {
        role: 'user',
        content: userText,
        time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      },
    ];

    simSending = true;
    try {
      const res = await api.post<{
        reply: string;
        tool_calls: ToolCall[];
        assistant: string;
        timestamp: string;
      }>('/ai/simulate', {
        assistant_id: simAssistantId,
        caller_number: simCallerNumber,
        messages: simMessages,
      });

      simLastTools = res.tool_calls || [];
      simMessages = [
        ...simMessages,
        {
          role: 'assistant',
          content: res.reply,
          time: res.timestamp || new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        },
      ];
    } catch (e) {
      toastError(e);
      simMessages = [
        ...simMessages,
        {
          role: 'assistant',
          content: 'Sorry, I encountered an internal processing error.',
          time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        },
      ];
    } finally {
      simSending = false;
    }
  }

  function viewSessionDetails(sess: AICallSession) {
    selectedSession = sess;
    try {
      selectedSessionMessages = JSON.parse(sess.MessagesJSON || '[]');
    } catch {
      selectedSessionMessages = [];
    }
  }

  function getSentimentTone(sentiment?: string): 'success' | 'warning' | 'danger' | 'neutral' {
    if (!sentiment) return 'neutral';
    const s = sentiment.toLowerCase();
    if (s.includes('pos')) return 'success';
    if (s.includes('neg')) return 'danger';
    return 'neutral';
  }

  onMount(loadAll);
</script>

<PageHeader 
  title="AI Assistants & Virtual Reception" 
  description="Multi-engine conversational AI receptionist, automated appointment booking, real-time IVR, and Whisper transcription."
>
  <div class="flex items-center space-x-2">
    <Button variant="secondary" on:click={loadAll}>Refresh</Button>
    <Button variant="primary" on:click={openCreateModal}>+ New AI Assistant</Button>
  </div>
</PageHeader>

<!-- Navigation Tabs -->
<div class="flex border-b border-slate-200 mb-6 space-x-6 text-sm font-medium">
  <button
    type="button"
    class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'assistants' ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    on:click={() => (activeTab = 'assistants')}
  >
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 100-6 3 3 0 000 6z"/></svg>
    <span>AI Assistants ({assistants.length})</span>
  </button>

  <button
    type="button"
    class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'simulator' ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    on:click={() => (activeTab = 'simulator')}
  >
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>
    <span>Interactive Simulator</span>
  </button>

  <button
    type="button"
    class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'sessions' ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    on:click={() => (activeTab = 'sessions')}
  >
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
    <span>Call Sessions & Transcripts ({sessions.length})</span>
  </button>

  <button
    type="button"
    class="pb-3 border-b-2 transition-colors flex items-center space-x-2 {activeTab === 'settings' ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    on:click={() => (activeTab = 'settings')}
  >
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
    <span>AI & STT Settings</span>
  </button>
</div>

<!-- TAB 1: ASSISTANTS TABLE -->
{#if activeTab === 'assistants'}
  <Card padded={false} title="Configured AI Assistants" subtitle="Virtual receptionists listening on internal extensions and public DIDs">
    {#if loading || loadError || assistants.length === 0}
      <EmptyState
        {loading}
        error={loadError}
        title="No AI Assistants configured"
        message="Create your first conversational assistant to start handling calls, booking appointments, and answering caller inquiries."
      >
        <svelte:fragment slot="retry">
          <Button size="sm" variant="primary" on:click={openCreateModal}>+ Create AI Assistant</Button>
        </svelte:fragment>
      </EmptyState>
    {:else}
      <div class="overflow-x-auto">
        <table class="tbl">
          <thead>
            <tr>
              <th>Assistant Name</th>
              <th>Ext</th>
              <th>Role & Persona</th>
              <th>Engine / Model</th>
              <th>Voice</th>
              <th>Capabilities</th>
              <th>Status</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each assistants as asst (asst.ID)}
              <tr class="hover:bg-slate-50/70 transition-colors">
                <td class="font-semibold text-slate-900">
                  <div class="flex items-center space-x-2">
                    <span class="w-2.5 h-2.5 rounded-full {asst.IsActive ? 'bg-emerald-500' : 'bg-slate-300'}"></span>
                    <span>{asst.Name}</span>
                  </div>
                </td>
                <td class="mono font-bold text-indigo-700">{asst.ExtensionNumber || '—'}</td>
                <td class="text-xs text-slate-600 max-w-xs truncate">{asst.RoleDescription}</td>
                <td>
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-slate-100 text-slate-800 border border-slate-200">
                    {asst.LLMProvider.toUpperCase()} / {asst.LLMModel}
                  </span>
                </td>
                <td class="text-xs text-slate-600 capitalize">{asst.VoiceID || 'alloy'}</td>
                <td>
                  <div class="flex items-center space-x-1.5">
                    {#if asst.EnableBooking}
                      <span class="px-1.5 py-0.5 rounded text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200" title="Appointment Booking Enabled">
                        📅 Booking
                      </span>
                    {/if}
                    {#if asst.EnableTransfer}
                      <span class="px-1.5 py-0.5 rounded text-[10px] font-semibold bg-blue-50 text-blue-700 border border-blue-200" title="Call Transfer to {asst.TransferTarget}">
                        📞 Ext {asst.TransferTarget}
                      </span>
                    {/if}
                  </div>
                </td>
                <td>
                  <Badge tone={asst.IsActive ? 'success' : 'neutral'}>
                    {asst.IsActive ? 'Active' : 'Disabled'}
                  </Badge>
                </td>
                <td class="text-right space-x-3 whitespace-nowrap">
                  <button
                    type="button"
                    class="link-action text-indigo-600 font-medium"
                    on:click={() => {
                      simAssistantId = asst.ID;
                      resetSimulator();
                      activeTab = 'simulator';
                    }}
                  >
                    Test Chat
                  </button>
                  <button type="button" class="link-action" on:click={() => openEditModal(asst)}>Edit</button>
                  <button type="button" class="link-danger" on:click={() => deleteAssistant(asst)}>Delete</button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </Card>
{/if}

<!-- TAB 2: INTERACTIVE SIMULATOR -->
{#if activeTab === 'simulator'}
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Simulator Chat Frame -->
    <div class="lg:col-span-2">
      <Card title="Live Conversational AI Simulator" subtitle="Interact directly with your assistant as a caller to test voice prompts and booking tools">
        <div class="flex flex-wrap items-center gap-3 mb-4 pb-4 border-b border-slate-100">
          <div class="flex-1 min-w-[200px]">
            <label class="block text-xs font-semibold text-slate-600 mb-1" for="sim-select-assistant">Testing Assistant:</label>
            <select id="sim-select-assistant" class="input !py-1.5 text-xs w-full" bind:value={simAssistantId} on:change={resetSimulator}>
              {#each assistants as asst}
                <option value={asst.ID}>{asst.Name} (Ext {asst.ExtensionNumber})</option>
              {/each}
            </select>
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-600 mb-1" for="sim-caller-number">Caller Number:</label>
            <input id="sim-caller-number" class="input !py-1.5 text-xs w-44" bind:value={simCallerNumber} />
          </div>
          <div class="self-end">
            <Button size="sm" variant="secondary" on:click={resetSimulator}>Reset Call</Button>
          </div>
        </div>

        <!-- Chat Bubble History -->
        <div class="h-96 overflow-y-auto space-y-3 p-4 bg-slate-50/70 rounded-xl border border-slate-200/60 mb-4 flex flex-col">
          {#each simMessages as msg}
            <div class="flex flex-col {msg.role === 'user' ? 'items-end' : 'items-start'}">
              <div class="flex items-center space-x-1.5 mb-1 text-[11px] text-slate-400">
                <span class="font-medium text-slate-600 capitalize">{msg.role === 'user' ? 'Caller' : 'AI Assistant'}</span>
                <span>•</span>
                <span>{msg.time}</span>
              </div>
              <div
                class="max-w-md p-3 rounded-2xl text-xs leading-relaxed shadow-sm {msg.role === 'user' ? 'bg-indigo-600 text-white rounded-br-none' : 'bg-white text-slate-800 border border-slate-200/80 rounded-bl-none'}"
              >
                {msg.content}
              </div>
            </div>
          {/each}

          {#if simSending}
            <div class="flex items-center space-x-2 text-slate-500 text-xs italic p-2">
              <span class="animate-pulse">Assistant is thinking & analyzing tools…</span>
            </div>
          {/if}
        </div>

        <!-- Simulator Input -->
        <form on:submit|preventDefault={sendSimulatorMessage} class="flex items-center space-x-2">
          <input
            class="input flex-1 !py-2.5 text-sm"
            placeholder="Type caller utterance (e.g. 'Can I schedule an appointment tomorrow at 10 AM?')..."
            bind:value={simInput}
            disabled={simSending}
            aria-label="Caller utterance"
          />
          <Button variant="primary" type="submit" disabled={simSending || !simInput.trim()}>
            Send Message
          </Button>
        </form>

        <!-- Quick Test Prompts -->
        <div class="mt-4 pt-4 border-t border-slate-100 flex flex-wrap gap-2 text-xs">
          <span class="text-slate-400 font-medium py-1">Sample utterances:</span>
          <button
            type="button"
            class="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 rounded-full text-slate-700 transition-colors"
            on:click={() => {
              simInput = 'I would like to book an appointment tomorrow at 2:00 PM for Alex.';
              sendSimulatorMessage();
            }}
          >
            "Book appointment tomorrow at 2 PM"
          </button>
          <button
            type="button"
            class="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 rounded-full text-slate-700 transition-colors"
            on:click={() => {
              simInput = 'Can you transfer me to extension 1001?';
              sendSimulatorMessage();
            }}
          >
            "Transfer me to extension 1001"
          </button>
          <button
            type="button"
            class="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 rounded-full text-slate-700 transition-colors"
            on:click={() => {
              simInput = 'What are your operating hours and office location?';
              sendSimulatorMessage();
            }}
          >
            "What are your business hours?"
          </button>
        </div>
      </Card>
    </div>

    <!-- Active Tool Execution Inspector -->
    <div class="space-y-6">
      <Card title="Executed Tool Calls" subtitle="Autonomous actions triggered during this turn">
        {#if simLastTools.length === 0}
          <div class="text-xs text-slate-400 italic py-6 text-center">
            No tools executed yet. Ask the assistant to book an appointment or transfer to see execution payloads.
          </div>
        {:else}
          <div class="space-y-3">
            {#each simLastTools as tool}
              <div class="p-3 bg-emerald-50/60 border border-emerald-200 rounded-lg text-xs space-y-1.5">
                <div class="flex items-center justify-between">
                  <span class="font-bold text-emerald-800 font-mono text-[11px]">{tool.name}</span>
                  <Badge tone="success">Triggered</Badge>
                </div>
                <div class="text-slate-700 text-[11px]">
                  <strong>Arguments:</strong>
                  <pre class="bg-white p-1.5 rounded mt-1 overflow-x-auto text-[10px] text-slate-800 font-mono border border-emerald-100">{JSON.stringify(tool.arguments, null, 2)}</pre>
                </div>
                <div class="text-emerald-700 text-[11px] font-medium pt-1">
                  ✓ {tool.result}
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </Card>

      <Card title="Dialplan Integration" subtitle="How calls reach this assistant">
        <div class="space-y-2 text-xs text-slate-600">
          <p>
            This assistant is accessible via FreeSWITCH audio fork at:
          </p>
          <div class="p-2 bg-slate-900 text-emerald-400 font-mono rounded text-[11px] leading-relaxed">
            Dial Ext: {assistants.find(a => a.ID === simAssistantId)?.ExtensionNumber || '890'}<br/>
            Feature Code: *99 (Global AI)<br/>
            WebSocket: ws://pbx_backend:8081/api/ai/stream
          </div>
          <p class="text-[11px] text-slate-500 pt-1">
            Incoming DIDs can route directly to this AI assistant under Inbound Routes.
          </p>
        </div>
      </Card>
    </div>
  </div>
{/if}

<!-- TAB 3: CALL SESSIONS & TRANSCRIPTS -->
{#if activeTab === 'sessions'}
  <Card padded={false} title="AI Call Sessions & Transcripts" subtitle="Real-time call transcripts, automated summaries, and sentiment analytics">
    {#if loading || loadError || sessions.length === 0}
      <EmptyState
        {loading}
        error={loadError}
        title="No AI call sessions recorded yet"
        message="Calls routed to your AI assistants will log complete multi-turn transcripts, sentiment, and executed tools here."
      />
    {:else}
      <div class="overflow-x-auto">
        <table class="tbl">
          <thead>
            <tr>
              <th>Call UUID / Date</th>
              <th>Assistant</th>
              <th>Caller</th>
              <th>Duration</th>
              <th>Action Taken</th>
              <th>Sentiment</th>
              <th>Summary</th>
              <th class="text-right">Transcript</th>
            </tr>
          </thead>
          <tbody>
            {#each sessions as sess (sess.ID)}
              <tr class="hover:bg-slate-50/70 transition-colors">
                <td class="text-xs">
                  <div class="mono font-semibold text-slate-800">{sess.CallUUID.substring(0, 8)}…</div>
                  <div class="text-slate-400 text-[11px]">{formatDateTime(sess.CreatedAt)}</div>
                </td>
                <td class="font-medium text-slate-800 text-xs">{sess.AssistantName}</td>
                <td class="mono text-xs text-slate-700">{sess.CallerNumber || 'Unknown'}</td>
                <td class="mono text-xs text-slate-600">{formatDuration(sess.DurationSec)}</td>
                <td>
                  <Badge tone={sess.ActionTaken === 'booking_created' ? 'success' : sess.ActionTaken === 'transferred' ? 'warning' : 'neutral'}>
                    {sess.ActionTaken.replace('_', ' ').toUpperCase()}
                  </Badge>
                </td>
                <td>
                  <Badge tone={getSentimentTone(sess.Sentiment)}>
                    {sess.Sentiment.toUpperCase()}
                  </Badge>
                </td>
                <td class="text-xs text-slate-600 max-w-sm truncate" title={sess.Summary}>
                  {sess.Summary || 'No summary generated'}
                </td>
                <td class="text-right whitespace-nowrap">
                  <button type="button" class="link-action font-medium" on:click={() => viewSessionDetails(sess)}>
                    View Transcript
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </Card>
{/if}

<!-- TAB 4: GLOBAL AI & TRANSCRIPTION SETTINGS -->
{#if activeTab === 'settings'}
  <Card title="Speech-to-Text & AI Synthesis Configuration" subtitle="Configure OpenAI Whisper, ElevenLabs, Deepgram, and local LLM endpoints">
    <form on:submit|preventDefault={saveSettings} class="space-y-6 max-w-3xl">
      <!-- Automation Toggles -->
      <div class="p-4 bg-slate-50 rounded-xl border border-slate-200/80 space-y-3">
        <div class="font-semibold text-slate-900 text-sm">Automated PBX Transcriptions</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <label class="flex items-center space-x-2 text-xs font-medium text-slate-700 cursor-pointer">
            <input type="checkbox" class="rounded border-slate-300 text-indigo-600" bind:checked={settings.AutoTranscribeVoicemail} />
            <span>Automatically transcribe all new Voicemails</span>
          </label>
          <label class="flex items-center space-x-2 text-xs font-medium text-slate-700 cursor-pointer">
            <input type="checkbox" class="rounded border-slate-300 text-indigo-600" bind:checked={settings.AutoTranscribeRecordings} />
            <span>Automatically transcribe recorded calls (CDRs)</span>
          </label>
          <label class="flex items-center space-x-2 text-xs font-medium text-slate-700 cursor-pointer">
            <input type="checkbox" class="rounded border-slate-300 text-indigo-600" bind:checked={settings.SentimentAnalysis} />
            <span>Enable Sentiment Analysis & Summaries</span>
          </label>
          <label class="flex items-center space-x-2 text-xs font-medium text-slate-700 cursor-pointer">
            <input type="checkbox" class="rounded border-slate-300 text-indigo-600" bind:checked={settings.EmailTranscripts} />
            <span>Email transcripts immediately to mailbox owners</span>
          </label>
        </div>
      </div>

      <!-- STT Configuration -->
      <div class="border-t border-slate-100 pt-4 space-y-4">
        <div class="font-semibold text-slate-900 text-sm">Speech-to-Text (STT) Engine</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-stt-provider">STT Provider</label>
            <select id="settings-stt-provider" class="input text-xs w-full" bind:value={settings.STTProvider}>
              <option value="whisper">OpenAI Whisper API (Standard)</option>
              <option value="deepgram">Deepgram Nova-2</option>
              <option value="google">Google Cloud Speech-to-Text</option>
              <option value="vosk">Vosk / Local Offline Model</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-stt-model">STT Model Identifier</label>
            <input id="settings-stt-model" class="input text-xs w-full" bind:value={settings.STTModel} placeholder="whisper-1" />
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-stt-api-key">STT API Key</label>
            <input id="settings-stt-api-key" type="password" class="input text-xs w-full font-mono" bind:value={settings.STTAPIKey} placeholder="sk-..." />
            <p class="text-[11px] text-slate-400 mt-1">If no API key is specified, the system utilizes intelligent built-in local transcription fallback.</p>
          </div>
        </div>
      </div>

      <!-- LLM Configuration -->
      <div class="border-t border-slate-100 pt-4 space-y-4">
        <div class="font-semibold text-slate-900 text-sm">Large Language Model (LLM) Backend</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-llm-provider">LLM Provider</label>
            <select id="settings-llm-provider" class="input text-xs w-full" bind:value={settings.LLMProvider}>
              <option value="openai">OpenAI (GPT-4o, GPT-4o-mini)</option>
              <option value="claude">Anthropic Claude</option>
              <option value="gemini">Google Gemini 1.5</option>
              <option value="ollama">Local Ollama / vLLM (Self-hosted)</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-llm-model">Default LLM Model</label>
            <input id="settings-llm-model" class="input text-xs w-full" bind:value={settings.LLMDefaultModel} placeholder="gpt-4o" />
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-llm-base-url">LLM Base URL (for Ollama or custom proxies)</label>
            <input id="settings-llm-base-url" class="input text-xs w-full font-mono" bind:value={settings.LLMBaseURL} placeholder="https://api.openai.com/v1" />
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-llm-api-key">LLM API Key</label>
            <input id="settings-llm-api-key" type="password" class="input text-xs w-full font-mono" bind:value={settings.LLMAPIKey} placeholder="sk-..." />
          </div>
        </div>
      </div>

      <!-- TTS Configuration -->
      <div class="border-t border-slate-100 pt-4 space-y-4">
        <div class="font-semibold text-slate-900 text-sm">Text-to-Speech (TTS) Voice Synthesis</div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-tts-provider">TTS Provider</label>
            <select id="settings-tts-provider" class="input text-xs w-full" bind:value={settings.TTSProvider}>
              <option value="openai">OpenAI TTS</option>
              <option value="elevenlabs">ElevenLabs Neural Voice</option>
              <option value="deepgram">Deepgram Aura</option>
              <option value="system">FreeSWITCH Local Flite / Native</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-tts-voice">Default Voice Persona</label>
            <select id="settings-tts-voice" class="input text-xs w-full" bind:value={settings.TTSDefaultVoice}>
              <option value="alloy">Alloy (Balanced, Neutral)</option>
              <option value="echo">Echo (Warm, Crisp)</option>
              <option value="fable">Fable (Expressive)</option>
              <option value="onyx">Onyx (Deep, Authoritative)</option>
              <option value="nova">Nova (Energetic, Friendly)</option>
              <option value="shimmer">Shimmer (Clear, Professional)</option>
            </select>
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="settings-tts-api-key">TTS API Key</label>
            <input id="settings-tts-api-key" type="password" class="input text-xs w-full font-mono" bind:value={settings.TTSAPIKey} placeholder="sk-..." />
          </div>
        </div>
      </div>

      <div class="flex justify-end pt-4 border-t border-slate-100">
        <Button variant="primary" type="submit">Save Settings</Button>
      </div>
    </form>
  </Card>
{/if}

<!-- CREATE / EDIT ASSISTANT MODAL -->
{#if showModal}
  <div class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full p-6 space-y-6">
      <div class="flex items-center justify-between border-b border-slate-100 pb-3">
        <h3 class="text-base font-bold text-slate-900">
          {isEditing ? 'Edit AI Assistant Persona' : 'Configure New AI Assistant'}
        </h3>
        <button type="button" class="text-slate-400 hover:text-slate-600" on:click={() => (showModal = false)}>✕</button>
      </div>

      <form on:submit|preventDefault={saveAssistant} class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-asst-name">Assistant Name *</label>
            <input id="modal-asst-name" class="input text-xs w-full" bind:value={formAssistant.Name} placeholder="Sarah - Front Desk" required />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-asst-ext">Extension Number *</label>
            <input id="modal-asst-ext" class="input text-xs w-full mono font-bold" bind:value={formAssistant.ExtensionNumber} placeholder="890" required />
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-asst-role">Role Description</label>
            <input id="modal-asst-role" class="input text-xs w-full" bind:value={formAssistant.RoleDescription} placeholder="Front Desk Receptionist & Booking Agent" />
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-asst-greeting">Greeting Spoken Message</label>
            <input id="modal-asst-greeting" class="input text-xs w-full" bind:value={formAssistant.GreetingMessage} placeholder="Hello, thank you for calling MSIP PBX. How can I help you today?" />
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-asst-prompt">System Prompt & Instructions *</label>
            <textarea
              id="modal-asst-prompt"
              class="input text-xs w-full h-28 leading-relaxed font-sans"
              bind:value={formAssistant.SystemPrompt}
              placeholder="You are an AI virtual receptionist..."
              required
            ></textarea>
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-asst-llm-provider">LLM Provider</label>
            <select id="modal-asst-llm-provider" class="input text-xs w-full" bind:value={formAssistant.LLMProvider}>
              <option value="openai">OpenAI</option>
              <option value="claude">Anthropic Claude</option>
              <option value="gemini">Google Gemini</option>
              <option value="ollama">Local Ollama</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-asst-voice">Voice Persona</label>
            <select id="modal-asst-voice" class="input text-xs w-full" bind:value={formAssistant.VoiceID}>
              <option value="alloy">Alloy (Neutral)</option>
              <option value="echo">Echo (Warm)</option>
              <option value="fable">Fable (Expressive)</option>
              <option value="onyx">Onyx (Deep)</option>
              <option value="nova">Nova (Friendly)</option>
              <option value="shimmer">Shimmer (Professional)</option>
            </select>
          </div>

          <!-- Feature Toggles -->
          <div class="sm:col-span-2 p-3 bg-slate-50 rounded-lg border border-slate-200/80 space-y-2">
            <div class="text-xs font-semibold text-slate-800">Agent Tool Execution</div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
              <label class="flex items-center space-x-2 cursor-pointer">
                <input type="checkbox" class="rounded border-slate-300 text-indigo-600" bind:checked={formAssistant.EnableBooking} />
                <span>Autonomous Appointment Booking</span>
              </label>
              <label class="flex items-center space-x-2 cursor-pointer">
                <input type="checkbox" class="rounded border-slate-300 text-indigo-600" bind:checked={formAssistant.EnableTransfer} />
                <span>Autonomous Call Transfer</span>
              </label>
            </div>
            {#if formAssistant.EnableTransfer}
              <div class="pt-2">
                <label class="block text-[11px] font-semibold text-slate-600 mb-1" for="modal-asst-transfer">Transfer Target Extension</label>
                <input id="modal-asst-transfer" class="input !py-1 text-xs w-48 mono" bind:value={formAssistant.TransferTarget} placeholder="1001" />
              </div>
            {/if}
          </div>

          <div class="flex items-center space-x-2 sm:col-span-2 pt-2">
            <input type="checkbox" id="modal-asst-active" class="rounded border-slate-300 text-indigo-600" bind:checked={formAssistant.IsActive} />
            <label for="modal-asst-active" class="text-xs font-semibold text-slate-800">Assistant is Active & Accepting Calls</label>
          </div>
        </div>

        <div class="flex justify-end space-x-3 pt-4 border-t border-slate-100">
          <Button variant="secondary" type="button" on:click={() => (showModal = false)}>Cancel</Button>
          <Button variant="primary" type="submit">Save Assistant</Button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- SESSION DETAIL MODAL -->
{#if selectedSession}
  <div class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full p-6 space-y-4">
      <div class="flex items-center justify-between border-b border-slate-100 pb-3">
        <div>
          <h3 class="text-base font-bold text-slate-900">Call Transcript & Analytics</h3>
          <p class="text-xs text-slate-500">Call UUID: {selectedSession.CallUUID}</p>
        </div>
        <button type="button" class="text-slate-400 hover:text-slate-600" on:click={() => (selectedSession = null)}>✕</button>
      </div>

      <div class="grid grid-cols-3 gap-3 p-3 bg-slate-50 rounded-xl text-xs">
        <div>
          <span class="text-slate-400 block text-[10px] uppercase font-bold">Caller</span>
          <span class="font-mono font-medium text-slate-800">{selectedSession.CallerNumber || 'Unknown'}</span>
        </div>
        <div>
          <span class="text-slate-400 block text-[10px] uppercase font-bold">Duration</span>
          <span class="font-mono text-slate-800">{formatDuration(selectedSession.DurationSec)}</span>
        </div>
        <div>
          <span class="text-slate-400 block text-[10px] uppercase font-bold">Sentiment</span>
          <Badge tone={getSentimentTone(selectedSession.Sentiment)}>{selectedSession.Sentiment.toUpperCase()}</Badge>
        </div>
      </div>

      {#if selectedSession.Summary}
        <div class="p-3 bg-indigo-50/60 border border-indigo-100 rounded-xl text-xs">
          <span class="font-bold text-indigo-900 block mb-1">AI Executive Summary:</span>
          <p class="text-indigo-950 leading-relaxed">{selectedSession.Summary}</p>
        </div>
      {/if}

      <div class="h-64 overflow-y-auto p-4 bg-slate-50 rounded-xl border border-slate-200/80 space-y-3">
        {#each selectedSessionMessages as msg}
          <div class="text-xs">
            <span class="font-bold uppercase text-[10px] {msg.role === 'user' ? 'text-indigo-600' : 'text-emerald-700'}">
              {msg.role}:
            </span>
            <span class="text-slate-800 ml-1.5 leading-relaxed">{msg.content}</span>
          </div>
        {/each}
        {#if selectedSessionMessages.length === 0}
          <div class="text-xs text-slate-400 italic">No turns captured for this session.</div>
        {/if}
      </div>

      <div class="flex justify-end pt-2">
        <Button variant="secondary" on:click={() => (selectedSession = null)}>Close</Button>
      </div>
    </div>
  </div>
{/if}
