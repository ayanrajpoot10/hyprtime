<script context="module">
    import { fade } from "svelte/transition";
</script>

<script lang="ts">
    import StatCard from "../components/StatCard.svelte";
    import AppList from "../components/AppList.svelte";
    import type { ScreenTimeOverview } from "../lib/types";

    export let overview: ScreenTimeOverview;
</script>

<div class="overview-view" in:fade={{ duration: 200 }}>
    <div class="stats-row">
        <StatCard
            title="Total Screen Time"
            value={overview.total_time_formatted}
            icon="clock"
            variant="primary"
        />
        <StatCard
            title="Today's Screen Time"
            value={overview.today_time_formatted}
            icon="today"
            variant="secondary"
        />
    </div>

    <AppList apps={overview.top_apps} title="Top Applications" />
</div>

<style>
    .overview-view {
        display: flex;
        flex-direction: column;
        gap: 24px;
    }

    .stats-row {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
        gap: 12px;
    }
</style>
