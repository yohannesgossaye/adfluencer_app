export type Organization = {
  id: string;
  name: string;
  slug: string;
  createdAt: string;
};

export type DashboardStats = {
  activeUsers: number;
  revenue: number;
  series: { label: string; value: number }[];
};
