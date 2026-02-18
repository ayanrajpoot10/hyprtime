export interface AppData {
    class: string;
    total_time: number;
    total_time_formatted: string;
    open_count: number;
    last_seen: string;
    percentage: number;
}

export interface ScreenTimeOverview {
    total_time: number;
    total_time_formatted: string;
    today_time: number;
    today_time_formatted: string;
    top_apps: AppData[];
}

export interface DailyData {
    date: string;
    total_time: number;
    total_time_formatted: string;
    apps: AppData[];
}

export type ViewMode = "overview" | "daily";
