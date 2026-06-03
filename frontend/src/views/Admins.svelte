<script lang="ts">
  import { onMount } from 'svelte';
  import GlassCard from '../components/ui/GlassCard.svelte';
  import Button from '../components/ui/Button.svelte';

  let users = [];
  let isLoading = true;

  let showModal = false;
  let newEmail = '';
  let newPassword = '';
  let newRole = 'viewer';

  async function fetchUsers() {
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/users`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('pbx_token')}` }
      });
      if (res.ok) {
        const data = await res.json();
        users = data.users || [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      isLoading = false;
    }
  }

  async function handleAddUser(e) {
    e.preventDefault();
    try {
      const res = await fetch(`http://${window.location.hostname}:8080/api/v1/users`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('pbx_token')}`
        },
        body: JSON.stringify({
          Email: newEmail,
          PasswordHash: newPassword, // Note: Handled by backend hashing
          Role: newRole
        })
      });
      if (res.ok) {
        showModal = false;
        newEmail = '';
        newPassword = '';
        newRole = 'viewer';
        fetchUsers();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(() => {
    fetchUsers();
  });
</script>

<div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full">
  <header class="flex justify-between items-end mb-4">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-white mb-1">Admins & Roles</h1>
      <p class="text-gray-400 text-sm">Manage dashboard access and custom role permissions for this tenant.</p>
    </div>
    <Button variant="primary" on:click={() => showModal = true}>
      <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Invite User
    </Button>
  </header>

  <GlassCard className="p-0 overflow-hidden">
    {#if isLoading}
      <div class="p-8 text-center text-gray-400 flex flex-col items-center">
        <svg class="w-8 h-8 animate-spin text-indigo-500 mb-4" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
        Loading users...
      </div>
    {:else if users.length === 0}
      <div class="p-12 text-center flex flex-col items-center justify-center">
        <div class="w-16 h-16 bg-gray-800/50 rounded-2xl flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
        </div>
        <h3 class="text-lg font-medium text-white mb-2">No users found</h3>
        <p class="text-gray-400 text-sm max-w-sm">There are no additional users configured for this tenant. Click "Invite User" to add one.</p>
      </div>
    {:else}
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-gray-900/50 border-b border-gray-800">
              <th class="py-4 px-6 text-xs font-semibold text-gray-400 uppercase tracking-wider">Email Address</th>
              <th class="py-4 px-6 text-xs font-semibold text-gray-400 uppercase tracking-wider">Role</th>
              <th class="py-4 px-6 text-xs font-semibold text-gray-400 uppercase tracking-wider">Created</th>
              <th class="py-4 px-6 text-xs font-semibold text-gray-400 uppercase tracking-wider text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-800/50">
            {#each users as user}
              <tr class="hover:bg-gray-800/30 transition-colors group">
                <td class="py-4 px-6">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 rounded-full bg-indigo-500/20 text-indigo-400 flex items-center justify-center font-bold text-sm">
                      {user.Email[0].toUpperCase()}
                    </div>
                    <span class="text-white font-medium">{user.Email}</span>
                  </div>
                </td>
                <td class="py-4 px-6">
                  <span class="px-2.5 py-1 rounded-md text-xs font-medium border
                    {user.Role === 'superadmin' ? 'bg-purple-500/10 text-purple-400 border-purple-500/20' : 
                     user.Role === 'admin' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 
                     'bg-gray-500/10 text-gray-400 border-gray-500/20'}">
                    {user.Role.toUpperCase()}
                  </span>
                </td>
                <td class="py-4 px-6 text-sm text-gray-400">
                  {new Date(user.CreatedAt).toLocaleDateString()}
                </td>
                <td class="py-4 px-6 text-right">
                  <button class="text-gray-500 hover:text-indigo-400 p-2 transition-colors">Edit</button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </GlassCard>
</div>

{#if showModal}
<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
  <div class="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl overflow-hidden animate-in zoom-in-95 duration-200">
    <div class="px-6 py-4 border-b border-gray-800 flex justify-between items-center bg-gray-900/50">
      <h3 class="text-lg font-semibold text-white">Invite User</h3>
      <button class="text-gray-500 hover:text-white transition-colors" aria-label="Close" on:click={() => showModal = false}>
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
      </button>
    </div>
    
    <form on:submit={handleAddUser} class="p-6 space-y-4">
      <div>
        <label for="email" class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Email Address</label>
        <input id="email" type="email" bind:value={newEmail} required class="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
      </div>
      <div>
        <label for="password" class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Password</label>
        <input id="password" type="password" bind:value={newPassword} required class="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
      </div>
      <div>
        <label for="role" class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Custom Role</label>
        <select id="role" bind:value={newRole} class="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
          <option value="admin">Administrator (Full Access)</option>
          <option value="viewer">Viewer (Read-Only)</option>
        </select>
      </div>
      
      <div class="pt-4 flex space-x-3">
        <Button variant="secondary" type="button" className="flex-1" on:click={() => showModal = false}>Cancel</Button>
        <Button variant="primary" type="submit" className="flex-1">Send Invite</Button>
      </div>
    </form>
  </div>
</div>
{/if}
