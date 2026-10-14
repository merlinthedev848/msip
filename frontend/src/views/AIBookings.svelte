<script lang="ts">
  import { onMount } from 'svelte';
  import { api, formatDateTime } from '../lib/api';
  import { toast, toastError } from '../lib/toast';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface AIBooking {
    ID: string;
    AssistantID: string;
    AssistantName: string;
    CallerNumber: string;
    CallerName: string;
    ServiceType: string;
    BookingDate: string; // YYYY-MM-DD
    BookingTime: string; // HH:MM
    DurationMinutes: number;
    Status: string; // confirmed, pending, completed, cancelled
    Notes: string;
    CallUUID: string;
    CreatedAt: string;
  }

  let bookings: AIBooking[] = [];
  let loading = true;
  let loadError = '';
  let search = '';
  let statusFilter = 'all';

  // Modal State
  let showModal = false;
  let isEditing = false;
  let formBooking: AIBooking = getBlankBooking();

  function getBlankBooking(): AIBooking {
    const today = new Date();
    today.setDate(today.getDate() + 1);
    const dateStr = today.toISOString().split('T')[0];

    return {
      ID: '',
      AssistantID: '',
      AssistantName: 'Sarah - AI Virtual Receptionist',
      CallerNumber: '+1 (555) 000-0000',
      CallerName: '',
      ServiceType: 'General Consultation',
      BookingDate: dateStr,
      BookingTime: '10:00',
      DurationMinutes: 30,
      Status: 'confirmed',
      Notes: 'Booked via AI conversational assistant',
      CallUUID: '',
      CreatedAt: '',
    };
  }

  $: filtered = bookings.filter(b => {
    const matchesSearch = !search.trim() ||
      `${b.CallerName} ${b.CallerNumber} ${b.ServiceType} ${b.Notes}`.toLowerCase().includes(search.trim().toLowerCase());
    const matchesStatus = statusFilter === 'all' || b.Status.toLowerCase() === statusFilter.toLowerCase();
    return matchesSearch && matchesStatus;
  });

  $: stats = {
    total: bookings.length,
    confirmed: bookings.filter(b => b.Status === 'confirmed').length,
    pending: bookings.filter(b => b.Status === 'pending').length,
    completed: bookings.filter(b => b.Status === 'completed').length,
    cancelled: bookings.filter(b => b.Status === 'cancelled').length,
  };

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ bookings: AIBooking[] }>('/ai/bookings');
      bookings = data.bookings || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  function openCreateModal() {
    isEditing = false;
    formBooking = getBlankBooking();
    showModal = true;
  }

  function openEditModal(b: AIBooking) {
    isEditing = true;
    formBooking = { ...b };
    showModal = true;
  }

  async function saveBooking() {
    if (!formBooking.CallerName.trim() || !formBooking.CallerNumber.trim()) {
      toast.error('Caller name and number are required');
      return;
    }
    try {
      if (isEditing) {
        await api.put(`/ai/bookings/${formBooking.ID}`, formBooking);
        toast.success('Appointment updated');
      } else {
        await api.post('/ai/bookings', formBooking);
        toast.success('Appointment booked successfully');
      }
      showModal = false;
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  async function updateStatus(b: AIBooking, newStatus: string) {
    try {
      await api.put(`/ai/bookings/${b.ID}`, { ...b, Status: newStatus });
      toast.success(`Booking marked as ${newStatus}`);
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  async function remove(b: AIBooking) {
    if (!confirm(`Delete appointment for ${b.CallerName}?`)) return;
    try {
      await api.del(`/ai/bookings/${b.ID}`);
      toast.success('Appointment deleted');
      await load();
    } catch (e) {
      toastError(e);
    }
  }

  function getStatusTone(status: string): 'success' | 'warning' | 'danger' | 'neutral' {
    switch (status.toLowerCase()) {
      case 'confirmed': return 'success';
      case 'completed': return 'neutral';
      case 'pending': return 'warning';
      case 'cancelled': return 'danger';
      default: return 'neutral';
    }
  }

  onMount(load);
</script>

<PageHeader 
  title="AI Appointment Bookings" 
  description="Autonomous calendar scheduling handled by AI Virtual Receptionists during incoming phone calls."
>
  <div class="flex items-center space-x-2">
    <Button variant="secondary" on:click={load}>Refresh</Button>
    <Button variant="primary" on:click={openCreateModal}>+ Schedule Appointment</Button>
  </div>
</PageHeader>

<!-- Summary Metric Cards -->
<div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
  <div class="p-4 bg-white rounded-xl border border-slate-200/80 shadow-sm">
    <span class="text-xs font-semibold text-slate-500 uppercase tracking-wide">Total Bookings</span>
    <div class="text-2xl font-bold text-slate-900 mt-1">{stats.total}</div>
  </div>
  <div class="p-4 bg-white rounded-xl border border-emerald-100 bg-emerald-50/20 shadow-sm">
    <span class="text-xs font-semibold text-emerald-700 uppercase tracking-wide">Confirmed</span>
    <div class="text-2xl font-bold text-emerald-800 mt-1">{stats.confirmed}</div>
  </div>
  <div class="p-4 bg-white rounded-xl border border-amber-100 bg-amber-50/20 shadow-sm">
    <span class="text-xs font-semibold text-amber-700 uppercase tracking-wide">Pending</span>
    <div class="text-2xl font-bold text-amber-800 mt-1">{stats.pending}</div>
  </div>
  <div class="p-4 bg-white rounded-xl border border-slate-200/80 shadow-sm">
    <span class="text-xs font-semibold text-slate-500 uppercase tracking-wide">Completed</span>
    <div class="text-2xl font-bold text-slate-700 mt-1">{stats.completed}</div>
  </div>
</div>

<Card padded={false} title="Scheduled Appointments" subtitle="{filtered.length} appointments listed">
  <svelte:fragment slot="actions">
    <div class="flex items-center space-x-3">
      <select class="input !py-1.5 !text-xs w-36" bind:value={statusFilter} aria-label="Filter status">
        <option value="all">All Statuses</option>
        <option value="confirmed">Confirmed</option>
        <option value="pending">Pending</option>
        <option value="completed">Completed</option>
        <option value="cancelled">Cancelled</option>
      </select>
      <input class="input !py-1.5 !text-xs w-56" placeholder="Search caller, service…" bind:value={search} aria-label="Search bookings" />
    </div>
  </svelte:fragment>

  {#if loading || loadError || filtered.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title={search || statusFilter !== 'all' ? 'No appointments match your filters' : 'No appointments scheduled'}
      message={search ? '' : 'When callers request an appointment, the AI assistant automatically books it and logs it here.'}
    >
      <svelte:fragment slot="retry">
        <Button size="sm" variant="primary" on:click={openCreateModal}>+ Schedule Appointment</Button>
      </svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Caller / Client</th>
            <th>Service</th>
            <th>Appointment Date & Time</th>
            <th>Duration</th>
            <th>Booked By</th>
            <th>Status</th>
            <th>Notes</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filtered as b (b.ID)}
            <tr class="hover:bg-slate-50/70 transition-colors">
              <td>
                <div class="font-bold text-slate-900 text-xs">{b.CallerName}</div>
                <div class="mono text-[11px] text-slate-500">{b.CallerNumber}</div>
              </td>
              <td class="text-xs font-medium text-slate-800">{b.ServiceType}</td>
              <td>
                <div class="font-semibold text-indigo-700 text-xs flex items-center space-x-1.5">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
                  <span>{b.BookingDate}</span>
                </div>
                <div class="text-[11px] text-slate-500 mono font-medium">{b.BookingTime}</div>
              </td>
              <td class="mono text-xs text-slate-600">{b.DurationMinutes}m</td>
              <td class="text-xs text-slate-600">
                <span class="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] bg-slate-100 text-slate-700 font-medium border border-slate-200">
                  🤖 {b.AssistantName || 'AI Receptionist'}
                </span>
              </td>
              <td>
                <Badge tone={getStatusTone(b.Status)}>
                  {b.Status.toUpperCase()}
                </Badge>
              </td>
              <td class="text-xs text-slate-500 max-w-xs truncate" title={b.Notes}>
                {b.Notes || '—'}
              </td>
              <td class="text-right space-x-2 whitespace-nowrap">
                {#if b.Status === 'pending'}
                  <button type="button" class="link-action text-emerald-600 font-medium" on:click={() => updateStatus(b, 'confirmed')}>
                    Confirm
                  </button>
                {/if}
                {#if b.Status === 'confirmed'}
                  <button type="button" class="link-action text-slate-600 font-medium" on:click={() => updateStatus(b, 'completed')}>
                    Complete
                  </button>
                {/if}
                <button type="button" class="link-action" on:click={() => openEditModal(b)}>Edit</button>
                <button type="button" class="link-danger" on:click={() => remove(b)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<!-- CREATE / EDIT BOOKING MODAL -->
{#if showModal}
  <div class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-6 space-y-6">
      <div class="flex items-center justify-between border-b border-slate-100 pb-3">
        <h3 class="text-base font-bold text-slate-900">
          {isEditing ? 'Edit Appointment' : 'Schedule New Appointment'}
        </h3>
        <button type="button" class="text-slate-400 hover:text-slate-600" on:click={() => (showModal = false)}>✕</button>
      </div>

      <form on:submit|preventDefault={saveBooking} class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-name">Caller / Client Name *</label>
            <input id="modal-booking-name" class="input text-xs w-full" bind:value={formBooking.CallerName} placeholder="Alex Smith" required />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-number">Caller Phone Number *</label>
            <input id="modal-booking-number" class="input text-xs w-full mono font-medium" bind:value={formBooking.CallerNumber} placeholder="+15551234567" required />
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-service">Service Type</label>
            <input id="modal-booking-service" class="input text-xs w-full" bind:value={formBooking.ServiceType} placeholder="Dental Checkup / Consultation" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-date">Date (YYYY-MM-DD) *</label>
            <input id="modal-booking-date" type="date" class="input text-xs w-full" bind:value={formBooking.BookingDate} required />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-time">Time (HH:MM) *</label>
            <input id="modal-booking-time" type="time" class="input text-xs w-full" bind:value={formBooking.BookingTime} required />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-duration">Duration (Minutes)</label>
            <input id="modal-booking-duration" type="number" class="input text-xs w-full" bind:value={formBooking.DurationMinutes} min="5" max="240" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-status">Booking Status</label>
            <select id="modal-booking-status" class="input text-xs w-full" bind:value={formBooking.Status}>
              <option value="confirmed">Confirmed</option>
              <option value="pending">Pending</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>
          <div class="sm:col-span-2">
            <label class="block text-xs font-semibold text-slate-700 mb-1" for="modal-booking-notes">Notes</label>
            <textarea id="modal-booking-notes" class="input text-xs w-full h-20" bind:value={formBooking.Notes} placeholder="Reason for visit, special requirements..."></textarea>
          </div>
        </div>

        <div class="flex justify-end space-x-3 pt-4 border-t border-slate-100">
          <Button variant="secondary" type="button" on:click={() => (showModal = false)}>Cancel</Button>
          <Button variant="primary" type="submit">Save Appointment</Button>
        </div>
      </form>
    </div>
  </div>
{/if}
