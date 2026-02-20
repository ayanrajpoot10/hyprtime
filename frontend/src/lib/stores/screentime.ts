import { writable } from "svelte/store";
import type { DailyData } from "../types";
import { fetchDailyStats } from "../services/screentime";

export const dailyData = writable<DailyData | null>(null);
export const loading = writable<boolean>(true);
export const error = writable<string>("");
export const selectedDate = writable<string>(
    new Date().toISOString().split("T")[0]
);

export async function loadDailyStats(date: string): Promise<void> {
    loading.set(true);
    error.set("");
    try {
        const data = await fetchDailyStats(date);
        dailyData.set(data);
    } catch (err: any) {
        error.set(err.message || "Failed to load daily stats");
    } finally {
        loading.set(false);
    }
}
