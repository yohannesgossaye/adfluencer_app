import type { MockUserRole } from "@/data/mock-users";

export type MockSessionUser = {
  email: string;
  role: MockUserRole;
  name: string;
  redirect: string;
};

const MOCK_AUTH_KEY = "adfluence_mock_user";

export function getMockUser(): MockSessionUser | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem(MOCK_AUTH_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as MockSessionUser;
  } catch {
    return null;
  }
}

export function setMockUser(user: MockSessionUser) {
  if (typeof window === "undefined") return;
  localStorage.setItem(MOCK_AUTH_KEY, JSON.stringify(user));
}

export function clearMockUser() {
  if (typeof window === "undefined") return;
  localStorage.removeItem(MOCK_AUTH_KEY);
}

export function isRole(role: MockUserRole): boolean {
  return getMockUser()?.role === role;
}

