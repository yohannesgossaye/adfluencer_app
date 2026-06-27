import {
  IconBuilding,
  IconHome,
  IconLogin,
} from "@tabler/icons-react";

export type NavItem = {
  title: string;
  href: string;
  icon: typeof IconHome;
};

export const mainNav: NavItem[] = [
  { title: "Dashboard", href: "/", icon: IconHome },
  { title: "Organizations", href: "/organizations", icon: IconBuilding },
];

export const secondaryNav: NavItem[] = [
  { title: "Sign in", href: "/auth/login", icon: IconLogin },
];
