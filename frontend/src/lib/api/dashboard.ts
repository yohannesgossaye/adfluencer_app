import { apiClient } from "@/lib/api/api-client";
import { apiRoutes } from "@/lib/api/api-routes";
import type { DashboardStats } from "@/lib/api/types";

export async function getDashboardStats(): Promise<DashboardStats> {
  if (!apiClient.defaults.baseURL) {
    return {
      activeUsers: 1284,
      revenue: 48200,
      series: [
        { label: "Mon", value: 120 },
        { label: "Tue", value: 200 },
        { label: "Wed", value: 150 },
        { label: "Thu", value: 280 },
        { label: "Fri", value: 220 },
        { label: "Sat", value: 90 },
        { label: "Sun", value: 60 },
      ],
    };
  }

  const { data } = await apiClient.get<DashboardStats>(apiRoutes.dashboardStats);
  return data;
}
