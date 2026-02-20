import { ScreenTimeService } from "../../../bindings/hyprtime";
import type { DailyData } from "../types";

export async function fetchDailyStats(date: string): Promise<DailyData> {
    const data = await ScreenTimeService.GetDailyStats(date);
    if (!data) throw new Error("No daily stats returned");
    return data as unknown as DailyData;
}
