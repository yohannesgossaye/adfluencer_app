import { mockUsers } from "@/data/mock-users";

export type AdminUserRole = "business" | "influencer" | "admin";
export type AdminUserStatus = "active" | "pending" | "blocked" | "removed";

export type AdminUserRecord = {
  id: string;
  name: string;
  email: string;
  role: AdminUserRole;
  status: AdminUserStatus;
  joinedAt: string;
  lastActive: string;
  campaignsCount: number;
  contractsCount: number;
};

export const initialAdminUsers: AdminUserRecord[] = [
  {
    id: "u-1",
    name: mockUsers[0]?.name ?? "Yohannes",
    email: mockUsers[0]?.email ?? "yohannes@gmail.com",
    role: "business",
    status: "active",
    joinedAt: "2025-11-02",
    lastActive: "2026-05-11",
    campaignsCount: 6,
    contractsCount: 4,
  },
  {
    id: "u-2",
    name: mockUsers[1]?.name ?? "Biruk",
    email: mockUsers[1]?.email ?? "biruk@gmail.com",
    role: "influencer",
    status: "active",
    joinedAt: "2025-10-18",
    lastActive: "2026-05-10",
    campaignsCount: 0,
    contractsCount: 3,
  },
  {
    id: "u-3",
    name: mockUsers[2]?.name ?? "Admin",
    email: mockUsers[2]?.email ?? "admin@adfluence.com",
    role: "admin",
    status: "active",
    joinedAt: "2025-01-05",
    lastActive: "2026-05-11",
    campaignsCount: 0,
    contractsCount: 0,
  },
  {
    id: "u-4",
    name: "Meron Gebre",
    email: "meron@addisbrands.et",
    role: "business",
    status: "pending",
    joinedAt: "2026-05-08",
    lastActive: "2026-05-08",
    campaignsCount: 0,
    contractsCount: 0,
  },
  {
    id: "u-5",
    name: "Samuel Tesfaye",
    email: "samuel.creates@gmail.com",
    role: "influencer",
    status: "blocked",
    joinedAt: "2026-02-14",
    lastActive: "2026-04-02",
    campaignsCount: 0,
    contractsCount: 1,
  },
  {
    id: "u-6",
    name: "Helina Imports",
    email: "ops@helinimports.et",
    role: "business",
    status: "removed",
    joinedAt: "2025-08-01",
    lastActive: "2026-01-10",
    campaignsCount: 2,
    contractsCount: 2,
  },
];

export const USERS_STORAGE_KEY = "adfluence_admin_users";

export function loadAdminUsers(): AdminUserRecord[] {
  if (typeof window === "undefined") return initialAdminUsers;
  try {
    const raw = localStorage.getItem(USERS_STORAGE_KEY);
    if (!raw) return initialAdminUsers;
    const parsed = JSON.parse(raw) as AdminUserRecord[];
    return Array.isArray(parsed) && parsed.length ? parsed : initialAdminUsers;
  } catch {
    return initialAdminUsers;
  }
}

export function saveAdminUsers(users: AdminUserRecord[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(USERS_STORAGE_KEY, JSON.stringify(users));
}
