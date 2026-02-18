import { ScreenTimeService } from "../../../bindings/hyprtime";
import type { ScreenTimeOverview, DailyData } from "../types";

export async function fetchOverview(): Promise<ScreenTimeOverview> {
    const data = await ScreenTimeService.GetOverview();
    if (!data) throw new Error("No overview data returned");
    return data as unknown as ScreenTimeOverview;
}

export async function fetchTodayStats(): Promise<DailyData> {
    const data = await ScreenTimeService.GetTodayStats();
    if (!data) throw new Error("No today stats returned");
    return data as unknown as DailyData;
}

export async function fetchDailyStats(date: string): Promise<DailyData> {
    const data = await ScreenTimeService.GetDailyStats(date);
    if (!data) throw new Error("No daily stats returned");
    return data as unknown as DailyData;
}
