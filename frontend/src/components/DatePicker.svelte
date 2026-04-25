<script lang="ts">
    export let value: string;
    export let max: string;
    export let onChange: (date: string) => void;

    $: isToday = value === max;
    $: displayDate = new Date(value + "T00:00:00").toLocaleDateString("en-US", {
        weekday: "short",
        month: "short",
        day: "numeric",
        year: "numeric",
    });

    const prev = () => {
        const d = new Date(value);
        d.setDate(d.getDate() - 1);
        onChange(d.toISOString().split("T")[0]);
    };

    const next = () => {
        const d = new Date(value);
        d.setDate(d.getDate() + 1);
        const nd = d.toISOString().split("T")[0];
        if (nd <= max) onChange(nd);
    };

    const handleChange = (e: Event) => {
        onChange((e.target as HTMLInputElement).value);
    };
</script>

<div class="picker">
    <button class="arrow" on:click={prev}>‹</button>
    <label class="date">
        <span>{displayDate}</span>
        <input type="date" {value} {max} on:change={handleChange} />
    </label>
    <button class="arrow" on:click={next} disabled={isToday}>›</button>
</div>

<style>
    .picker {
        display: flex;
        align-items: center;
        gap: 4px;
    }

    .arrow {
        width: 28px;
        height: 28px;
        border: none;
        background: transparent;
        color: var(--text-secondary);
        cursor: pointer;
        font-size: 1.25rem;
        line-height: 1;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition:
            background 0.15s,
            color 0.15s;
    }

    .arrow:hover:not(:disabled) {
        background: var(--surface-elevated);
        color: var(--text-primary);
    }

    .arrow:disabled {
        opacity: 0.3;
        cursor: not-allowed;
    }

    .date {
        position: relative;
        display: flex;
        align-items: center;
        font-size: 0.9rem;
        font-weight: 500;
        color: var(--text-primary);
        cursor: pointer;
        padding: 4px 8px;
        border-radius: 6px;
        transition: background 0.15s;
    }

    .date:hover {
        background: var(--surface-elevated);
    }

    .date input {
        position: absolute;
        inset: 0;
        opacity: 0;
        width: 100%;
        cursor: pointer;
        border: none;
    }
</style>
