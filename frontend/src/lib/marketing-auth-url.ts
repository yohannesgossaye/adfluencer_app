export type MarketingAuthMode = "login" | "signup";

export type MarketingAuthRole = "business" | "influencer";

export function buildMarketingAuthHref(
  pathname: string,
  mode: MarketingAuthMode,
  role?: MarketingAuthRole,
) {
  const params = new URLSearchParams();
  params.set("auth", mode);
  if (role) {
    params.set("role", role);
  }
  return `${pathname}?${params.toString()}`;
}

export function parseMarketingAuthRole(
  value: string | null | undefined,
): MarketingAuthRole | undefined {
  if (value === "business" || value === "influencer") {
    return value;
  }
  return undefined;
}
