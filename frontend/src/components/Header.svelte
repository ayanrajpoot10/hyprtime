<script lang="ts">
    import { refresh } from "../lib/stores/screentime";
</script>

<header class="header">
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
        }

        h1 {
            font-size: 1.1rem;
        }
    }
</style>
