<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import PageHeader from '../components/ui/PageHeader.svelte';
  import Card from '../components/ui/Card.svelte';
  import Button from '../components/ui/Button.svelte';
  import Badge from '../components/ui/Badge.svelte';
  import EmptyState from '../components/ui/EmptyState.svelte';

  interface CommitItem {
    hash: string;
    type: string;
    time: string;
    msg: string;
    author: string;
  }

  let commits: CommitItem[] = [];
  let loading = true;
  let loadError = '';

  async function load() {
    loading = true;
    loadError = '';
    try {
      const data = await api.get<{ commits: CommitItem[] }>('/commits');
      commits = data.commits || [];
    } catch (e: any) {
      loadError = e.message;
    } finally {
      loading = false;
    }
  }

  onMount(load);
</script>

<PageHeader title="Change History" description="Audit log of system configuration updates and version history.">
  <Button variant="secondary" on:click={load}>Refresh</Button>
</PageHeader>

<Card padded={false} title="Audit Trail" subtitle="{commits.length} change entries">
  {#if loading || loadError || commits.length === 0}
    <EmptyState
      {loading}
      error={loadError}
      title="No change history entries found"
      message="Audit logs and configuration revisions will appear here."
    >
      <svelte:fragment slot="retry"><Button size="sm" variant="secondary" on:click={load}>Retry</Button></svelte:fragment>
    </EmptyState>
  {:else}
    <div class="overflow-x-auto">
      <table class="tbl">
        <thead>
          <tr>
            <th>Commit / Hash</th>
            <th>Category</th>
            <th>Change Description</th>
            <th>Author</th>
            <th>Timestamp</th>
          </tr>
        </thead>
        <tbody>
          {#each commits as c}
            <tr>
              <td>
                <span class="mono px-2 py-0.5 bg-sky-50 text-sky-800 rounded text-xs font-bold border border-sky-200">
                  {c.hash}
                </span>
              </td>
              <td><Badge tone="info">{c.type || 'System'}</Badge></td>
              <td class="text-xs text-slate-800 font-medium">{c.msg}</td>
              <td class="text-xs text-slate-600">{c.author || 'System'}</td>
              <td class="text-xs text-slate-500">{c.time}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>
