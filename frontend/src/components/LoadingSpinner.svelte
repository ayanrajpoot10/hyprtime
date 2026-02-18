<script lang="ts">
</script>

<div class="loading-container">
    <div class="spinner">
        <div class="spinner-ring" />
    </div>
    <p class="loading-text">Loading screen time data...</p>
</div>

<style>
    .loading-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 80px 24px;
        gap: 20px;
    }

    .spinner {
        position: relative;
        width: 40px;
        height: 40px;
    }

    .spinner-ring {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        border: 3px solid var(--border);
        border-top-color: var(--accent);
        animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .loading-text {
        margin: 0;
        font-size: 0.9rem;
        color: var(--text-tertiary);
        font-weight: 500;
    }
</style>
