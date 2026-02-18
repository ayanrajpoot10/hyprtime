<script lang="ts">
    import type { AppData } from "../lib/types";
    import AppItem from "./AppItem.svelte";
    import EmptyState from "./EmptyState.svelte";

    export let apps: AppData[] = [];
    export let title: string = "Applications";
</script>

<section class="app-list-section">
    <div class="section-header">
        <h2>{title}</h2>
        {#if apps.length > 0}
            <span class="badge">{apps.length}</span>
        {/if}
    </div>

    {#if apps.length > 0}
        <div class="app-list">
            {#each apps as app, i (app.class)}
                <AppItem {app} rank={i} />
            {/each}
        </div>
    {:else}
        <EmptyState message="No application data available" icon="apps" />
    {/if}
</section>

<style>
    .app-list-section {
        margin-top: 8px;
    }

    .section-header {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-bottom: 16px;
    }

    h2 {
        margin: 0;
        font-size: 1.1rem;
        font-weight: 600;
        color: var(--text-primary);
        letter-spacing: -0.01em;
    }

    .badge {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 22px;
        height: 22px;
        padding: 0 7px;
        background: var(--accent-subtle);
        color: var(--accent);
        font-size: 0.75rem;
        font-weight: 600;
        border-radius: 6px;
    }

    .app-list {
        display: flex;
        flex-direction: column;
        gap: 8px;
    }
</style>
