"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, X } from "lucide-react";

import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { useMarketingAuth } from "@/components/marketing/marketing-auth-provider";
import { Button } from "@/components/ui/button";
import { Sheet, SheetClose, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { cn } from "@/lib/utils";

export function MarketingNavbar() {
  const pathname = usePathname();
  const { openLogin, openSignup } = useMarketingAuth();

  const links = [
    { href: "/for-businesses", label: "For businesses" },
    { href: "/for-influencers", label: "For creators" },
    { href: "/pricing", label: "Fees" },
    { href: "/contact", label: "Contact us" },
  ] as const;

  return (
    <header className="sticky top-0 z-50 overflow-x-hidden border-b border-white/10 bg-[color:var(--adfluence-bg)] text-white">
      <div className="relative flex h-[72px] min-w-0 items-center gap-2 px-4 sm:gap-3 sm:px-6 md:px-10 lg:px-16 xl:px-20">
        <BrandLogo showNameFrom="sm" className="min-w-0" />
        <nav aria-label="Main" className="absolute left-1/2 hidden -translate-x-1/2 items-center gap-9 text-[15px] font-medium text-[#CBD5E1] lg:flex">
          {links.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              aria-current={pathname === link.href ? "page" : undefined}
              className={cn("transition-colors hover:text-white", pathname === link.href && "text-white")}
            >
              {link.label}
            </Link>
          ))}
        </nav>
        <div className="ml-auto flex shrink-0 items-center gap-2 sm:gap-3">
          <Sheet>
            <SheetTrigger asChild>
              <Button
                variant="outline"
                size="icon"
                className="h-11 w-11 rounded-xl border-white/20 bg-transparent text-white hover:bg-white/10 lg:hidden"
                aria-label="Open menu"
              >
                <Menu className="size-5" />
              </Button>
            </SheetTrigger>
            <SheetContent
              side="right"
              showCloseButton={false}
              className="w-[88vw] max-w-[360px] border-l border-white/10 bg-[color:var(--adfluence-surface)] text-white backdrop-blur-xl"
            >
              <div className="flex h-full flex-col">
                <div className="flex items-start justify-between border-b border-white/10 p-4">
                  <BrandLogo showName />
                  <SheetClose asChild>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-11 w-11 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                      aria-label="Close menu"
                    >
                      <X className="size-5" />
                    </Button>
                  </SheetClose>
                </div>

                <div className="flex-1 space-y-2 p-4">
                  {links.map((link) => {
                    const active =
                      link.href.startsWith("/#")
                          ? false
                          : pathname === link.href || pathname.startsWith(`${link.href}/`);
                    return (
                      <SheetClose asChild key={link.href}>
                        <Link
                          href={link.href}
                          className={cn(
                            "flex h-11 items-center gap-3 rounded-xl px-3 text-sm text-[#CBD5E1] transition-colors",
                            "hover:bg-white/8 hover:text-white",
                            active ? "bg-white/8 text-white" : null,
                          )}
                        >
                          <span className="relative flex size-2.5 items-center justify-center">
                            <span
                              className={cn(
                                "size-2.5 rounded-full bg-white/20",
                                active ? "bg-[color:var(--adfluence-action)]" : null,
                              )}
                            />
                          </span>
                          {link.label}
                        </Link>
                      </SheetClose>
                    );
                  })}
                </div>

                <div className="border-t border-white/10 p-4">
                  <div className="grid gap-3">
                    <SheetClose asChild>
                      <Button
                        type="button"
                        variant="outline"
                        onClick={openLogin}
                        className="h-12 w-full rounded-xl border-white/15 bg-transparent px-6 text-white hover:bg-white/10"
                      >
                        Log in
                      </Button>
                    </SheetClose>
                    <SheetClose asChild>
                      <GradientButton
                        type="button"
                        onClick={() => openSignup()}
                        className="h-12 w-full rounded-xl px-6"
                      >
                        Join free
                      </GradientButton>
                    </SheetClose>
                  </div>
                </div>
              </div>
            </SheetContent>
          </Sheet>
          <Button
            type="button"
            variant="ghost"
            onClick={openLogin}
            className="hidden h-10 rounded-xl px-5 font-semibold sm:inline-flex"
          >
            Log in
          </Button>
          <GradientButton
            type="button"
            onClick={() => openSignup()}
            className="hidden h-10 rounded-xl px-5 font-semibold sm:inline-flex"
          >
            Join free
          </GradientButton>
        </div>
      </div>
    </header>
  );
}
