<script lang="ts">
  import { onDestroy } from "svelte";
  import { Window } from "@wailsio/runtime";

  import Header from "./components/Header.svelte";
  import LoadingSpinner from "./components/LoadingSpinner.svelte";
  import ErrorMessage from "./components/ErrorMessage.svelte";
  import OverviewView from "./views/OverviewView.svelte";
  import DailyView from "./views/DailyView.svelte";

  import {
    overview,
    dailyData,
    loading,
    error,
    viewMode,
    loadOverview,
    refresh,
  } from "./lib/stores/screentime";

  // Initial data load
  loadOverview();

  // Auto-refresh every 30 seconds
  const interval = setInterval(refresh, 30000);

  // Show window after frontend is ready
  Window.Show();

  onDestroy(() => {
    clearInterval(interval);
  });
</script>

<div class="app-shell">
  <Header />

  <main class="main-content">
    {#if $loading}
      <LoadingSpinner />
    {:else if $error}
      <ErrorMessage message={$error} onRetry={refresh} />
    {:else if $viewMode === "overview" && $overview}
      <OverviewView overview={$overview} />
    {:else if $viewMode === "daily" && $dailyData}
      <DailyView dailyData={$dailyData} />
    {/if}
  </main>
</div>

<style>
  .app-shell {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    background: var(--bg);
  }

  .main-content {
    flex: 1;
    padding: 24px;
    max-width: 860px;
    width: 100%;
    margin: 0 auto;
    box-sizing: border-box;
  }

  @media (max-width: 600px) {
    .main-content {
      padding: 16px;
    }
  }
</style>
