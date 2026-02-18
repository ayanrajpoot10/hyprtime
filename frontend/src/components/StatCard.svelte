<script lang="ts">
    export let title: string;
    export let value: string;
    export let icon: string = "clock";
    export let variant: "primary" | "secondary" | "accent" = "primary";
</script>

<div class="stat-card {variant}">
    <div class="stat-icon">
        {#if icon === "clock"}
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <circle cx="12" cy="12" r="10" />
                <polyline points="12 6 12 12 16 14" />
            </svg>
        {:else if icon === "today"}
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
        {:else if icon === "calendar"}
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                <line x1="16" y1="2" x2="16" y2="6" />
                <line x1="8" y1="2" x2="8" y2="6" />
                <line x1="3" y1="10" x2="21" y2="10" />
            </svg>
        {:else if icon === "apps"}
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <rect x="3" y="3" width="7" height="7" />
                <rect x="14" y="3" width="7" height="7" />
                <rect x="14" y="14" width="7" height="7" />
                <rect x="3" y="14" width="7" height="7" />
            </svg>
        {/if}
    </div>
    <div class="stat-content">
        <span class="stat-label">{title}</span>
        <span class="stat-value">{value}</span>
    </div>
</div>

<style>
    .stat-card {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px 24px;
        background: var(--surface-elevated);
        border: 1px solid var(--border);
        border-radius: 16px;
        transition: all 0.25s ease;
    }

    .stat-card:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    }

    .stat-card.primary .stat-icon {
        background: var(--accent-gradient);
        color: white;
    }

    .stat-card.secondary .stat-icon {
        background: linear-gradient(135deg, #10b981, #059669);
        color: white;
    }

    .stat-card.accent .stat-icon {
        background: linear-gradient(135deg, #f59e0b, #d97706);
        color: white;
    }

    .stat-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 44px;
        height: 44px;
        border-radius: 12px;
        flex-shrink: 0;
    }

    .stat-content {
        display: flex;
        flex-direction: column;
        gap: 4px;
        min-width: 0;
    }

    .stat-label {
        font-size: 0.8rem;
        font-weight: 500;
        color: var(--text-tertiary);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .stat-value {
        font-size: 1.5rem;
        font-weight: 700;
        color: var(--text-primary);
        letter-spacing: -0.02em;
        line-height: 1.2;
    }

    @media (max-width: 600px) {
        .stat-card {
            padding: 16px 18px;
        }

        .stat-value {
            font-size: 1.25rem;
        }
    }
</style>
