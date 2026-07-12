import type { ReactNode } from "react";

import { DashboardHeader } from "@/components/dashboard/dashboard-header";
import { DashboardSidebar } from "@/components/dashboard/dashboard-sidebar";
import type { NavItem } from "@/lib/navigation";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu } from "lucide-react";

export function DashboardShell({
  sidebarTitle,
  navItems,
  headerTitle,
  headerSubtitle,
  children,
}: {
  sidebarTitle: string;
  navItems: NavItem[];
  headerTitle: string;
  headerSubtitle: string;
  children: ReactNode;
}) {
  return (
    <div className="min-h-screen overflow-x-hidden bg-[color:var(--adfluence-bg)] text-white">
      {/* Desktop: fixed sidebar — does not scroll with page */}
      <div className="fixed inset-y-0 left-0 z-40 hidden h-[100dvh] w-[260px] md:block">
        <DashboardSidebar title={sidebarTitle} items={navItems} variant="desktop" />
      </div>

      <div
        className="flex min-h-screen flex-col md:ml-[260px] md:h-[100dvh] md:min-h-0 md:overflow-hidden"
      >
        <DashboardHeader
          title={headerTitle}
          subtitle={headerSubtitle}
          leading={
            <Sheet>
              <SheetTrigger asChild>
                <Button
                  variant="outline"
                  size="icon"
                  className="size-11 border-white/15 bg-white/5 text-white hover:bg-white/10 md:size-9"
                  aria-label="Open sidebar"
                >
                  <Menu className="size-4" />
                </Button>
              </SheetTrigger>
              <SheetContent
                side="left"
                className="flex h-full w-[min(100%,260px)] max-w-[260px] flex-col border-white/10 bg-[color:var(--adfluence-surface)] p-0 sm:max-w-[260px]"
              >
                <DashboardSidebar title={sidebarTitle} items={navItems} variant="drawer" />
              </SheetContent>
            </Sheet>
          }
        />
        <main className="min-w-0 flex-1 overflow-y-auto bg-[color:var(--adfluence-bg)] p-4 md:p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
