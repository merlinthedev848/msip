<script lang="ts">
  import { Router, Route, links } from "svelte-routing";
  import Sidebar from './components/Sidebar.svelte';
  
  import Dashboard from './views/Dashboard.svelte';
  import Extensions from './views/Extensions.svelte';
  import Routing from './views/Routing.svelte';

  export let url = "";
</script>

<!-- The use:links action automatically intercepts anchor tags and uses the router instead of hard reloads -->
<div use:links>
  <Router {url}>
    <main class="flex min-h-screen bg-gray-900 font-sans selection:bg-indigo-500 selection:text-white text-gray-100">
      <Sidebar />
      
      <div class="flex-1 flex flex-col h-screen overflow-hidden">
        <!-- Header -->
        <header class="h-16 border-b border-gray-800 bg-gray-900/80 backdrop-blur-md flex items-center justify-between px-8 sticky top-0 z-20">
          <div class="flex items-center">
            <h2 class="text-sm font-semibold text-gray-400 uppercase tracking-widest">Workspace</h2>
          </div>
          <div class="flex items-center space-x-4">
            <button class="bg-gray-800 hover:bg-gray-700 text-sm py-1.5 px-4 rounded-full transition-colors border border-gray-700 text-gray-300 shadow-sm flex items-center space-x-2">
              <div class="w-5 h-5 rounded-full bg-indigo-500 flex items-center justify-center text-xs text-white font-bold">A</div>
              <span>Admin</span>
            </button>
          </div>
        </header>

        <!-- Main Content Area -->
        <div class="flex-1 overflow-y-auto p-8 relative">
          <!-- Ambient Background Glow -->
          <div class="absolute top-0 left-1/4 w-96 h-96 bg-indigo-600/10 rounded-full blur-3xl pointer-events-none -z-10"></div>
          <div class="absolute bottom-0 right-1/4 w-96 h-96 bg-emerald-600/5 rounded-full blur-3xl pointer-events-none -z-10"></div>
          
          <div class="max-w-7xl mx-auto">
            <Route path="/" component={Dashboard} />
            <Route path="/extensions" component={Extensions} />
            <Route path="/inbound-routes" component={Routing} />
            
            <!-- Fallback Catch All for Unimplemented Routes -->
            <Route path="/outbound-routes">
              <div class="p-8 text-center text-gray-500">Outbound Routes module coming soon.</div>
            </Route>
            <Route path="/trunks">
              <div class="p-8 text-center text-gray-500">SIP Trunks module coming soon.</div>
            </Route>
          </div>
        </div>
      </div>
    </main>
  </Router>
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background-color: #111827; /* gray-900 */
  }
</style>
