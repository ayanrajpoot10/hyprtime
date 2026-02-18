<script lang="ts">
    export let value: string;
    export let max: string;
    export let onChange: (date: string) => void;

    const handleChange = (e: Event) => {
        const input = e.target as HTMLInputElement;
        onChange(input.value);
    };

    const goToPrevDay = () => {
        const d = new Date(value);
        d.setDate(d.getDate() - 1);
        const newDate = d.toISOString().split("T")[0];
        onChange(newDate);
    };

    const goToNextDay = () => {
        const d = new Date(value);
        d.setDate(d.getDate() + 1);
        const newDate = d.toISOString().split("T")[0];
        if (newDate <= max) {
            onChange(newDate);
        }
    };

    $: isToday = value === max;

    $: displayDate = new Date(value + "T00:00:00").toLocaleDateString("en-US", {
        weekday: "short",
        month: "short",
        day: "numeric",
        year: "numeric",
    });
</script>

<div class="date-picker">
    <button class="nav-arrow" on:click={goToPrevDay} title="Previous day">
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
            <polyline points="15 18 9 12 15 6" />
        </svg>
    </button>

    <label class="date-label">
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
        <span class="date-display">{displayDate}</span>
        <input
            type="date"
            {value}
            {max}
            on:change={handleChange}
            class="date-input"
        />
    </label>

    <button
        class="nav-arrow"
        class:disabled={isToday}
        on:click={goToNextDay}
        disabled={isToday}
        title="Next day"
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
            <polyline points="9 18 15 12 9 6" />
        </svg>
    </button>
</div>

<style>
    .date-picker {
        display: flex;
        align-items: center;
        gap: 8px;
        background: var(--surface-elevated);
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 6px;
    }

    .nav-arrow {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 34px;
        height: 34px;
        padding: 0;
        margin: 0;
        border: none;
        background: var(--surface);
        color: var(--text-secondary);
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .nav-arrow:hover:not(.disabled) {
        background: var(--accent-subtle);
        color: var(--accent);
    }

    .nav-arrow.disabled {
        opacity: 0.3;
        cursor: not-allowed;
    }

    .date-label {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 6px 12px;
        color: var(--text-primary);
        cursor: pointer;
        position: relative;
        font-weight: 500;
    }

    .date-label svg {
        color: var(--accent);
        flex-shrink: 0;
    }

    .date-display {
        font-size: 0.9rem;
        white-space: nowrap;
    }

    .date-input {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        opacity: 0;
        cursor: pointer;
        border: none;
        padding: 0;
    }
</style>
