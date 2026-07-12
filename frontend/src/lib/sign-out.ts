import { AUTH_TOKEN_STORAGE_KEY } from "@/lib/constants";
import { clearMockUser } from "@/lib/mock-auth";

export function signOut() {
  clearMockUser();
  if (typeof window !== "undefined") {
    window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
  }
}
