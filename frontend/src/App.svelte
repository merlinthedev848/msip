<script lang="ts">
  import Sidebar from './components/Sidebar.svelte';
</script>

<main class="flex min-h-screen bg-gray-900 font-sans selection:bg-indigo-500 selection:text-white text-gray-100">
  <Sidebar />
  
  <div class="flex-1 flex flex-col h-screen overflow-hidden">
    <!-- Header -->
    <header class="h-16 border-b border-gray-800 bg-gray-900/50 backdrop-blur-md flex items-center justify-between px-8 sticky top-0 z-10">
      <div class="flex items-center">
        <h2 class="text-lg font-semibold text-gray-100">Live Status Dashboard</h2>
      </div>
      <div class="flex items-center space-x-4">
        <button class="bg-gray-800 hover:bg-gray-700 text-sm py-1.5 px-4 rounded-full transition-colors border border-gray-700 text-gray-300">
          Admin
        </button>
      </div>
    </header>

    <!-- Main Content Area -->
    <div class="flex-1 overflow-y-auto p-8">
      <div class="max-w-7xl mx-auto space-y-6">
        
        <!-- Metrics Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-gray-800/40 border border-gray-700/50 p-6 rounded-2xl shadow-sm hover:border-gray-600 transition-colors backdrop-blur-sm relative overflow-hidden group">
            <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <h3 class="text-sm font-medium text-gray-400 mb-1">Active Calls</h3>
            <p class="text-4xl font-light text-white">0</p>
          </div>
          
          <div class="bg-gray-800/40 border border-gray-700/50 p-6 rounded-2xl shadow-sm hover:border-gray-600 transition-colors backdrop-blur-sm relative overflow-hidden group">
            <div class="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <h3 class="text-sm font-medium text-gray-400 mb-1">Registered Extensions</h3>
            <p class="text-4xl font-light text-white">12</p>
          </div>
          
          <div class="bg-gray-800/40 border border-gray-700/50 p-6 rounded-2xl shadow-sm hover:border-gray-600 transition-colors backdrop-blur-sm relative overflow-hidden group">
            <div class="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <h3 class="text-sm font-medium text-gray-400 mb-1">RAM Footprint</h3>
            <p class="text-4xl font-light text-blue-400">142<span class="text-lg text-gray-500 ml-1">MB</span></p>
          </div>
        </div>

        <!-- Telemetry Canvas -->
        <div class="bg-gray-800/40 border border-gray-700/50 rounded-2xl p-6 backdrop-blur-sm shadow-sm h-96 flex flex-col justify-center items-center relative overflow-hidden">
             <!-- Decorative elements to make it feel premium -->
             <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-indigo-500 to-transparent opacity-20"></div>
             
             <svg class="w-16 h-16 text-gray-700 mb-4 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
               <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
             </svg>
             <p class="text-gray-400 text-sm tracking-wide">Waiting for system telemetry streams...</p>
        </div>

      </div>
    </div>
  </div>
</main>

<style>
  /* Base structural reset inside the component to ensure styles apply strictly */
  :global(body) {
    margin: 0;
    padding: 0;
  }
</style>
