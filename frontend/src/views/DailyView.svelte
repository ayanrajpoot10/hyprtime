<script lang="ts">
    import AppList from "../components/AppList.svelte";
    import DatePicker from "../components/DatePicker.svelte";
    import type { DailyData } from "../lib/types";
    import { selectedDate, loadDailyStats } from "../lib/stores/screentime";

    export let dailyData: DailyData;

    const today = new Date().toISOString().split("T")[0];

    const handleDateChange = (date: string) => {
        selectedDate.set(date);
        loadDailyStats(date);
    };
</script>

<div class="view">
    <div class="top">
        <DatePicker
            value={$selectedDate}
            max={today}
            onChange={handleDateChange}
        />
        <div class="summary">
            <span class="total">{dailyData.total_time_formatted}</span>
            <span class="meta">{dailyData.apps?.length ?? 0} apps</span>
        </div>
    </div>

    <AppList apps={dailyData.apps ?? []} />
</div>

<style>
    .view {
        display: flex;
        flex-direction: column;
        gap: 24px;
    }

    .top {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
    }

    .summary {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: 2px;
    }

    .total {
        font-size: 1.5rem;
        font-weight: 700;
        color: var(--text-primary);
        letter-spacing: -0.02em;
    }

    .meta {
        font-size: 0.8rem;
        color: var(--text-tertiary);
    }
</style>
