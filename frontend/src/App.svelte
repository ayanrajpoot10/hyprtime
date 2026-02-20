<script lang="ts">
  import { onDestroy } from "svelte";
  import { Window } from "@wailsio/runtime";

  import Header from "./components/Header.svelte";
  import LoadingSpinner from "./components/LoadingSpinner.svelte";
  import ErrorMessage from "./components/ErrorMessage.svelte";
  import DailyView from "./views/DailyView.svelte";

  import {
    dailyData,
    loading,
    error,
    selectedDate,
    loadDailyStats,
    refresh,
  } from "./lib/stores/screentime";

  // Initial load with today's date
  loadDailyStats(new Date().toISOString().split("T")[0]);

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
    {:else if $dailyData}
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
