<script lang="ts">
    import type { AppData } from "../lib/types";

    export let app: AppData;

    $: initials = app.class
        .split(/[-_.\s]/)
        .map((w) => w[0]?.toUpperCase() ?? "")
        .slice(0, 2)
        .join("");
</script>

<div class="item">
    <div class="avatar">
        <span class="initials">{initials}</span>
    </div>
    <div class="info">
        <div class="row">
            <span class="name">{app.class}</span>
            <span class="time">{app.total_time_formatted}</span>
        </div>
        <div class="track">
            <div class="fill" style="width:{app.percentage}%;" />
        </div>
    </div>
</div>

<style>
    .item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 10px 12px;
        border-radius: 8px;
        transition: background 0.15s;
    }

    .item:hover {
        background: var(--surface-elevated);
    }

    .avatar {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 36px;
        height: 36px;
        border-radius: 8px;
        flex-shrink: 0;
        background: var(--surface-elevated);
    }

    .initials {
        font-size: 0.75rem;
        font-weight: 700;
        color: var(--text-secondary);
    }

    .info {
        flex: 1;
        min-width: 0;
    }

    .row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 6px;
    }

    .name {
        font-size: 0.9rem;
        font-weight: 500;
        color: var(--text-primary);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .time {
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--text-secondary);
        flex-shrink: 0;
        margin-left: 8px;
    }

    .track {
        height: 3px;
        background: var(--border);
        border-radius: 2px;
        overflow: hidden;
    }

    .fill {
        height: 100%;
        border-radius: 2px;
        background: var(--accent);
        transition: width 0.4s ease;
        opacity: 0.6;
    }
</style>
