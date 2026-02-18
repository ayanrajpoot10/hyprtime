import { writable, derived } from "svelte/store";
import type { ScreenTimeOverview, DailyData, ViewMode } from "../types";
import { fetchOverview, fetchTodayStats, fetchDailyStats } from "../services/screentime";

// --- State stores ---
export const overview = writable<ScreenTimeOverview | null>(null);
export const dailyData = writable<DailyData | null>(null);
export const loading = writable<boolean>(true);
export const error = writable<string>("");
export const viewMode = writable<ViewMode>("overview");
export const selectedDate = writable<string>(
    new Date().toISOString().split("T")[0]
);

// --- Derived stores ---
export const hasOverviewData = derived(overview, ($o) => $o !== null);
export const hasDailyData = derived(dailyData, ($d) => $d !== null);

// --- Actions ---
export async function loadOverview(): Promise<void> {
    loading.set(true);
    error.set("");
    try {
        const [overviewData, todayData] = await Promise.all([
            fetchOverview(),
            fetchTodayStats(),
        ]);
        overview.set(overviewData);
        dailyData.set(todayData);
    } catch (err: any) {
        error.set(err.message || "Failed to load data");
        console.error(err);
    } finally {
        loading.set(false);
    }
}

export async function loadDailyStats(date: string): Promise<void> {
    loading.set(true);
    error.set("");
    try {
        const data = await fetchDailyStats(date);
        dailyData.set(data);
    } catch (err: any) {
        error.set(err.message || "Failed to load daily stats");
        console.error(err);
    } finally {
        loading.set(false);
    }
}

export function switchView(view: ViewMode): void {
    viewMode.set(view);
    if (view === "overview") {
        loadOverview();
    } else {
        selectedDate.update((d) => {
            loadDailyStats(d);
            return d;
        });
    }
}

export function refresh(): void {
    viewMode.update((v) => {
        if (v === "overview") {
            loadOverview();
        } else {
            selectedDate.update((d) => {
                loadDailyStats(d);
                return d;
            });
        }
        return v;
    });
}
