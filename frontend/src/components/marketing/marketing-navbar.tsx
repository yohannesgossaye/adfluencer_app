/* eslint-disable react/jsx-newline */
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, X } from "lucide-react";

import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { Sheet, SheetClose, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { cn } from "@/lib/utils";

export function MarketingNavbar() {
  const pathname = usePathname();

  const links = [
    { href: "/", label: "Home" },
    { href: "/for-businesses", label: "For Businesses" },
    { href: "/for-influencers", label: "For Influencers" },
    { href: "/#how-it-works", label: "How It Works" },
    { href: "/pricing", label: "Pricing" },
  ] as const;

  return (
    <header className="sticky top-0 z-50 border-b border-white/10 bg-[#0B1A2B] text-white">
      <div className="relative flex h-[72px] items-center px-4 sm:px-6 md:px-10 lg:px-16 xl:px-20">
        <Link href="/" className="flex shrink-0 items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED]">
            <span className="text-xl font-bold text-white">A</span>
          </div>
          <div className="leading-tight">
            <p className="text-xl font-bold">Adfluence</p>
            <p className="text-xs text-slate-300">Connect. Collaborate. Grow.</p>
          </div>
        </Link>
        <nav className="hidden items-center gap-8 text-sm text-slate-300 lg:flex absolute left-1/2 -translate-x-1/2">
          <Link href="/" className="hover:text-white transition-colors">
            Home
          </Link>
          <Link href="/for-businesses" className="hover:text-white transition-colors">
            For Businesses
          </Link>
          <Link href="/for-influencers" className="hover:text-white transition-colors">
            For Influencers
          </Link>
          <Link href="/#how-it-works" className="hover:text-white transition-colors">
            How It Works
          </Link>
          <Link href="/pricing" className="hover:text-white transition-colors">
            Pricing
          </Link>
        </nav>
        <div className="ml-auto flex shrink-0 items-center gap-3">
          <Sheet>
            <SheetTrigger asChild>
              <Button
                variant="outline"
                size="icon"
                className="h-10 w-10 rounded-xl border-white/20 bg-transparent text-white hover:bg-white/10 lg:hidden"
                aria-label="Open menu"
              >
                <Menu className="size-5" />
              </Button>
            </SheetTrigger>
            <SheetContent
              side="right"
              showCloseButton={false}
              className="w-[88vw] max-w-[360px] border-l border-white/10 bg-[#06111F] text-white backdrop-blur-xl"
            >
              <div className="flex h-full flex-col">
                <div className="flex items-start justify-between border-b border-white/10 p-4">
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-[#FF7A00] to-[#7C3AED]">
                      <span className="text-lg font-bold text-white">A</span>
                    </div>
                    <div className="leading-tight">
                      <p className="text-base font-semibold text-white">Adfluence</p>
                      <p className="text-xs text-[#CBD5E1]">
                        Connect. Collaborate. Grow.
                      </p>
                    </div>
                  </div>
                  <SheetClose asChild>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-10 w-10 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                      aria-label="Close menu"
                    >
                      <X className="size-5" />
                    </Button>
                  </SheetClose>
                </div>

                <div className="flex-1 space-y-2 p-4">
                  {links.map((link) => {
                    const active =
                      link.href === "/"
                        ? pathname === "/"
                        : link.href.startsWith("/#")
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
                                active ? "adfluence-gradient" : null,
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
                        asChild
                        variant="outline"
                        className="h-12 w-full rounded-xl border-white/15 bg-transparent px-6 text-white hover:bg-white/10"
                      >
                        <Link href="/auth/login">Login</Link>
                      </Button>
                    </SheetClose>
                    <SheetClose asChild>
                      <GradientButton asChild className="h-12 w-full rounded-xl px-6">
                        <Link href="/auth/signup">Sign Up</Link>
                      </GradientButton>
                    </SheetClose>
                  </div>
                </div>
              </div>
            </SheetContent>
          </Sheet>
          <Button
            asChild
            variant="outline"
            className="hidden h-10 rounded-xl border-white/20 bg-transparent px-5 text-white hover:bg-white/10 sm:inline-flex"
          >
            <Link href="/auth/login">Login</Link>
          </Button>
          <GradientButton asChild className="h-10 rounded-xl px-5">
            <Link href="/auth/signup">Sign Up</Link>
          </GradientButton>
        </div>
      </div>
    </header>
  );
}
