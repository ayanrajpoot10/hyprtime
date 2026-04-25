export interface AppData {
    class: string;
    total_time: number;
    total_time_formatted: string;
    open_count: number;
    last_seen: string;
    percentage: number;
}

export interface DailyData {
    date: string;
    total_time: number;
    total_time_formatted: string;
    apps: AppData[];
}
