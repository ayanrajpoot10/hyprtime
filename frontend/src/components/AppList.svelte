<script lang="ts">
    import type { AppData } from "../lib/types";
    import AppItem from "./AppItem.svelte";

    export let apps: AppData[] = [];
</script>

<div class="list">
    {#if apps.length === 0}
        <p class="empty">No application data for this day.</p>
    {:else}
        {#each apps as app (app.class)}
            <AppItem {app} />
        {/each}
    {/if}
</div>

<style>
    .list {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    .empty {
        padding: 48px 0;
        text-align: center;
        color: var(--text-tertiary);
        font-size: 0.9rem;
    }
</style>
