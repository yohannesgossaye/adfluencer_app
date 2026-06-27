export type NavIconKey =
  | "dashboard"
  | "briefcase"
  | "sparkles"
  | "recommendations"
  | "messages"
  | "contracts"
  | "profile"
  | "users"
  | "building"
  | "disputes"
  | "categories"
  | "support";

export type NavItem = {
  label: string;
  href: string;
  icon: NavIconKey;
};

export const businessNav: NavItem[] = [
  { label: "Dashboard", href: "/business/dashboard", icon: "dashboard" },
  { label: "Campaigns", href: "/business/campaigns", icon: "briefcase" },
  {
    label: "Create Campaign",
    href: "/business/campaigns/create",
    icon: "sparkles",
  },
  {
    label: "Recommendations",
    href: "/business/recommendations",
    icon: "recommendations",
  },
  { label: "Messages", href: "/business/messages", icon: "messages" },
  { label: "Contracts", href: "/business/contracts", icon: "contracts" },
  { label: "Support", href: "/business/support", icon: "support" },
];

export const influencerNav: NavItem[] = [
  {
    label: "Dashboard",
    href: "/influencer/dashboard",
    icon: "dashboard",
  },
  { label: "Profile", href: "/influencer/profile", icon: "profile" },
  { label: "Invites", href: "/influencer/invites", icon: "briefcase" },
  { label: "Messages", href: "/influencer/messages", icon: "messages" },
  { label: "Contracts", href: "/influencer/contracts", icon: "contracts" },
  { label: "Support", href: "/influencer/support", icon: "support" },
];

export const adminNav: NavItem[] = [
  { label: "Dashboard", href: "/admin/dashboard", icon: "dashboard" },
  { label: "Influencers", href: "/admin/influencers", icon: "users" },
  { label: "Campaigns", href: "/admin/campaigns", icon: "briefcase" },
  { label: "Users", href: "/admin/users", icon: "building" },
  { label: "Disputes", href: "/admin/disputes", icon: "disputes" },
  { label: "Categories", href: "/admin/categories", icon: "categories" },
  { label: "Support", href: "/admin/support", icon: "support" },
];
