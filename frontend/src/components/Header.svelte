<script lang="ts">
    import type { ViewMode } from "../lib/types";
    import { viewMode } from "../lib/stores/screentime";
    import { switchView, refresh } from "../lib/stores/screentime";

    let currentView: ViewMode;
    viewMode.subscribe((v) => (currentView = v));
</script>

<header class="header">
    <div class="header-left">
        <div class="brand">
            <div class="brand-icon">
                <svg
                    width="22"
                    height="22"
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
            </div>
            <h1>HyprTime</h1>
        </div>
    </div>

    <nav class="nav-tabs">
        <button
            class="nav-tab"
            class:active={currentView === "overview"}
            on:click={() => switchView("overview")}
        >
            <svg
                width="16"
                height="16"
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
            Overview
        </button>
        <button
            class="nav-tab"
            class:active={currentView === "daily"}
            on:click={() => switchView("daily")}
        >
            <svg
                width="16"
                height="16"
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
            Daily
        </button>
    </nav>

    <button class="refresh-btn" on:click={refresh} title="Refresh data">
        <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
        >
            <polyline points="23 4 23 10 17 10" />
            <polyline points="1 20 1 14 7 14" />
            <path
                d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"
            />
        </svg>
    </button>
</header>

<style>
    .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 24px;
        background: var(--surface-elevated);
        border-bottom: 1px solid var(--border);
        backdrop-filter: blur(20px);
        -webkit-backdrop-filter: blur(20px);
        position: sticky;
        top: 0;
        z-index: 100;
    }

    .header-left {
        display: flex;
        align-items: center;
        gap: 16px;
    }

    .brand {
        display: flex;
        align-items: center;
        gap: 10px;
    }

    .brand-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 36px;
        height: 36px;
        background: var(--accent-gradient);
        border-radius: 10px;
        color: white;
    }

    h1 {
        margin: 0;
        font-size: 1.25rem;
        font-weight: 700;
        color: var(--text-primary);
        letter-spacing: -0.02em;
    }

    .nav-tabs {
        display: flex;
        gap: 4px;
        background: var(--surface);
        padding: 4px;
        border-radius: 10px;
        border: 1px solid var(--border);
    }

    .nav-tab {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 16px;
        border: none;
        background: transparent;
        color: var(--text-secondary);
        font-size: 0.85rem;
        font-weight: 500;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s ease;
        white-space: nowrap;
        width: auto;
        height: auto;
        line-height: normal;
        margin: 0;
    }

    .nav-tab:hover {
        color: var(--text-primary);
        background: var(--surface-hover);
    }

    .nav-tab.active {
        background: var(--accent);
        color: white;
        box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
    }

    .refresh-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 36px;
        height: 36px;
        padding: 0;
        margin: 0;
        border: 1px solid var(--border);
        background: var(--surface);
        color: var(--text-secondary);
        border-radius: 10px;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .refresh-btn:hover {
        color: var(--accent);
        border-color: var(--accent);
        background: var(--accent-subtle);
    }

    .refresh-btn:active {
        transform: rotate(180deg);
    }

    @media (max-width: 600px) {
        .header {
            padding: 12px 16px;
            flex-wrap: wrap;
            gap: 12px;
        }

        h1 {
            font-size: 1.1rem;
        }

        .nav-tabs {
            order: 3;
            width: 100%;
            justify-content: center;
        }
    }
</style>
