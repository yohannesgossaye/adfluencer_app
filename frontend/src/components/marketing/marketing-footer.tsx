import Link from "next/link";

import { BrandLogo } from "@/components/common/brand-logo";
import { TibebBand } from "@/components/marketing/tibeb-band";
import { contactLinks } from "@/lib/site-contact";

// Only pages that exist. Add Privacy and Terms here once those pages are written.
const LINKS = [
  { href: "/for-businesses", label: "For businesses" },
  { href: "/for-influencers", label: "For creators" },
  { href: "/pricing", label: "Fees" },
  { href: "/auth/login", label: "Log in" },
];

export function MarketingFooter() {
  return (
    <footer className="bg-[color:var(--adfluence-surface)] text-white">
      <TibebBand />
      <div className="mx-auto flex max-w-[1200px] flex-col gap-8 px-5 py-10 sm:px-8 lg:flex-row lg:items-center lg:justify-between">
        <BrandLogo href="/" showName />
        <nav aria-label="Footer" className="flex flex-wrap gap-x-7 gap-y-3 text-sm text-[#94A3B8]">
          {LINKS.map((l) => (
            <Link key={l.href} href={l.href} className="transition-colors hover:text-white">
              {l.label}
            </Link>
          ))}
        </nav>
        <div className="flex flex-wrap gap-x-6 gap-y-3 text-sm text-[#94A3B8]">
          <a href={contactLinks.email} className="transition-colors hover:text-white">
            Email
          </a>
          <a href={contactLinks.phone} className="transition-colors hover:text-white">
            Phone
          </a>
          <a href={contactLinks.telegram} target="_blank" rel="noopener noreferrer" className="transition-colors hover:text-white">
            Telegram
          </a>
          <span>© {new Date().getFullYear()} Adfluence · Addis Ababa</span>
        </div>
      </div>
    </footer>
  );
}
