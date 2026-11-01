/** Where a signed-in user belongs, matching the backend's redirect rules. */
export function homeFor(user?: { role?: string; status?: string } | null): string {
  switch (user?.role) {
    case "admin":
      return "/admin/dashboard";
    case "business":
      return "/business/dashboard";
    case "influencer":
      // The creator dashboard opens after approval; until then they work on their profile.
      return user.status === "active" ? "/influencer/dashboard" : "/influencer/profile";
    default:
      return "/auth/login";
  }
}
