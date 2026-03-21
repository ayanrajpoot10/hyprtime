import * as ScreenTimeService from "../../../bindings/hyprtime/internal/gui/service/screentimeservice";
import type { DailyData } from "../types";

export async function fetchDailyStats(date: string): Promise<DailyData> {
    const data = await ScreenTimeService.GetDailyStats(date);
    if (!data) throw new Error("No daily stats returned");
    return data as unknown as DailyData;
}
