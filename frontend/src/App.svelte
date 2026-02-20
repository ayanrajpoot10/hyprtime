<script lang="ts">
  import { Window } from "@wailsio/runtime";
  import DailyView from "./views/DailyView.svelte";
  import {
    dailyData,
    loading,
    error,
    loadDailyStats,
  } from "./lib/stores/screentime";

  loadDailyStats(new Date().toISOString().split("T")[0]);
  Window.Show();
</script>

<div class="app">
  <main>
    {#if $loading}
      <div class="state"><div class="spinner" /></div>
    {:else if $error}
      <div class="state"><p class="error-msg">{$error}</p></div>
    {:else if $dailyData}
      <DailyView dailyData={$dailyData} />
    {/if}
  </main>
</div>

<style>
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  main {
    flex: 1;
    padding: 24px;
    max-width: 800px;
    width: 100%;
    margin: 0 auto;
    box-sizing: border-box;
  }

  .state {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 80px 24px;
  }

  .spinner {
    width: 28px;
    height: 28px;
    border: 2px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .error-msg {
    color: var(--text-secondary);
    font-size: 0.9rem;
  }

  @media (max-width: 600px) {
    main {
      padding: 16px;
    }
  }
</style>
