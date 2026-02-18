<script context="module">
    import { fade } from "svelte/transition";
</script>

<script lang="ts">
    import StatCard from "../components/StatCard.svelte";
    import AppList from "../components/AppList.svelte";
    import DatePicker from "../components/DatePicker.svelte";
    import type { DailyData } from "../lib/types";
    import { selectedDate, loadDailyStats } from "../lib/stores/screentime";

    export let dailyData: DailyData;

    let currentDate: string;
    selectedDate.subscribe((d) => (currentDate = d));

    const today = new Date().toISOString().split("T")[0];

    const handleDateChange = (date: string) => {
        selectedDate.set(date);
        loadDailyStats(date);
    };
</script>

<div class="daily-view" in:fade={{ duration: 200 }}>
    <div class="daily-header">
        <DatePicker
            value={currentDate}
            max={today}
            onChange={handleDateChange}
        />
    </div>

    <div class="stats-row">
        <StatCard
            title="Screen Time"
            value={dailyData.total_time_formatted}
            icon="calendar"
            variant="accent"
        />
        <StatCard
            title="Apps Used"
            value={dailyData.apps ? String(dailyData.apps.length) : "0"}
            icon="apps"
            variant="primary"
        />
    </div>

    <AppList apps={dailyData.apps || []} title="Applications" />
</div>

<style>
    .daily-view {
        display: flex;
        flex-direction: column;
        gap: 24px;
    }

    .daily-header {
        display: flex;
        align-items: center;
        justify-content: flex-start;
    }

    .stats-row {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
        gap: 12px;
    }
</style>
