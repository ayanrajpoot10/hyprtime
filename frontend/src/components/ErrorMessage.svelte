<script lang="ts">
    export let message: string = "Something went wrong";
    export let onRetry: (() => void) | null = null;
</script>

<div class="error-container">
    <div class="error-icon">
        <svg
            width="28"
            height="28"
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
    </div>
    <p class="error-text">{message}</p>
    {#if onRetry}
        <button class="retry-btn" on:click={onRetry}>
            <svg
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <polyline points="23 4 23 10 17 10" />
                <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
            </svg>
            Try Again
        </button>
    {/if}
</div>

<style>
    .error-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 60px 24px;
        gap: 12px;
    }

    .error-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 52px;
        height: 52px;
        border-radius: 14px;
        background: rgba(239, 68, 68, 0.1);
        color: #ef4444;
    }

    .error-text {
        margin: 0;
        font-size: 0.9rem;
        color: var(--text-secondary);
        text-align: center;
        max-width: 300px;
    }

    .retry-btn {
        display: flex;
        align-items: center;
        gap: 6px;
        margin-top: 8px;
        padding: 8px 18px;
        width: auto;
        height: auto;
        line-height: normal;
        background: var(--accent);
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 0.85rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .retry-btn:hover {
        background: var(--accent-hover);
        transform: translateY(-1px);
    }
</style>
