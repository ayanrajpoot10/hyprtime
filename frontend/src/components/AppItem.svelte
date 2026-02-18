<script lang="ts">
    import type { AppData } from "../lib/types";

    export let app: AppData;
    export let rank: number = 0;

    const APP_COLORS = [
        "#6366f1",
        "#8b5cf6",
        "#ec4899",
        "#f43f5e",
        "#f97316",
        "#eab308",
        "#22c55e",
        "#14b8a6",
        "#06b6d4",
        "#3b82f6",
    ];

    $: barColor = APP_COLORS[rank % APP_COLORS.length];
    $: initials = app.class
        .split(/[-_.\s]/)
        .map((w) => w.charAt(0).toUpperCase())
        .slice(0, 2)
        .join("");
</script>

<div class="app-item">
    <div class="app-row">
        <div
            class="app-icon"
            style="background: {barColor}20; color: {barColor};"
        >
            {initials}
        </div>
        <div class="app-details">
            <div class="app-name-row">
                <span class="app-name">{app.class}</span>
                <span class="app-time" style="color: {barColor};"
                    >{app.total_time_formatted}</span
                >
            </div>
            <div class="app-meta-row">
                <span class="app-meta">
                    <svg
                        width="12"
                        height="12"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        ><path
                            d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4M10 17l5-5-5-5M13.8 12H3"
                        /></svg
                    >
                    {app.open_count} opens
                </span>
                {#if app.last_seen}
                    <span class="app-meta">
                        <svg
                            width="12"
                            height="12"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            ><circle cx="12" cy="12" r="10" /><polyline
                                points="12 6 12 12 16 14"
                            /></svg
                        >
                        {new Date(app.last_seen).toLocaleTimeString([], {
                            hour: "2-digit",
                            minute: "2-digit",
                        })}
                    </span>
                {/if}
                <span class="app-percentage">{app.percentage.toFixed(1)}%</span>
            </div>
        </div>
    </div>
    <div class="progress-track">
        <div
            class="progress-fill"
            style="width: {app.percentage}%; background: {barColor};"
        />
    </div>
</div>

<style>
    .app-item {
        padding: 16px 20px;
        background: var(--surface-elevated);
        border: 1px solid var(--border);
        border-radius: 14px;
        transition: all 0.2s ease;
    }

    .app-item:hover {
        border-color: var(--border-hover);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
        transform: translateY(-1px);
    }

    .app-row {
        display: flex;
        align-items: center;
        gap: 14px;
        margin-bottom: 12px;
    }

    .app-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 40px;
        height: 40px;
        border-radius: 10px;
        font-size: 0.8rem;
        font-weight: 700;
        flex-shrink: 0;
        letter-spacing: 0.02em;
    }

    .app-details {
        flex: 1;
        min-width: 0;
    }

    .app-name-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 12px;
        margin-bottom: 6px;
    }

    .app-name {
        font-size: 0.95rem;
        font-weight: 600;
        color: var(--text-primary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .app-time {
        font-size: 0.95rem;
        font-weight: 700;
        flex-shrink: 0;
    }

    .app-meta-row {
        display: flex;
        align-items: center;
        gap: 12px;
        flex-wrap: wrap;
    }

    .app-meta {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 0.75rem;
        color: var(--text-tertiary);
    }

    .app-percentage {
        font-size: 0.75rem;
        font-weight: 600;
        color: var(--text-secondary);
        margin-left: auto;
    }

    .progress-track {
        width: 100%;
        height: 4px;
        background: var(--surface);
        border-radius: 2px;
        overflow: hidden;
    }

    .progress-fill {
        height: 100%;
        border-radius: 2px;
        transition: width 0.5s cubic-bezier(0.4, 0, 0.2, 1);
    }

    @media (max-width: 600px) {
        .app-item {
            padding: 12px 14px;
        }

        .app-icon {
            width: 34px;
            height: 34px;
            font-size: 0.7rem;
        }
    }
</style>
