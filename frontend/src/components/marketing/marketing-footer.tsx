import Link from "next/link";

import { BrandLogo } from "@/components/common/brand-logo";
import { Container } from "@/components/common/container";

export function MarketingFooter() {
  return (
    <footer className="border-t border-white/10 bg-[color:var(--adfluence-surface)] py-16 text-white">
      <Container>
        <div className="mb-12 grid grid-cols-2 gap-8 text-white md:grid-cols-4">
          <div>
            <h4 className="mb-4 font-semibold">Platform</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>
                <Link href="/for-businesses" className="hover:text-white transition-colors">
                  For Businesses
                </Link>
              </li>
              <li>
                <Link href="/for-influencers" className="hover:text-white transition-colors">
                  For Influencers
                </Link>
              </li>
              <li>
                <Link href="/#how-it-works" className="hover:text-white transition-colors">
                  How It Works
                </Link>
              </li>
              <li>
                <Link href="/pricing" className="hover:text-white transition-colors">
                  How Fees Work
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h4 className="mb-4 font-semibold">Company</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>About Us</li>
              <li>Careers</li>
              <li>Blog</li>
              <li>Press Kit</li>
            </ul>
          </div>
          <div>
            <h4 className="mb-4 font-semibold">Resources</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>Help Center</li>
              <li>Community</li>
              <li>Case Studies</li>
              <li>API Docs</li>
            </ul>
          </div>
          <div>
            <h4 className="mb-4 font-semibold">Contact</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>hello@adfluence.et</li>
              <li>Addis Ababa, Ethiopia</li>
              <li>adfluence.et</li>
            </ul>
          </div>
        </div>

        <div className="flex flex-col items-center justify-between gap-4 border-t border-white/10 pt-8 text-[#CBD5E1] md:flex-row">
          <BrandLogo href="/" showName />
          <div className="text-sm text-[#CBD5E1]">
            © 2026 Adfluence. All rights reserved.
          </div>
          <div className="flex gap-6 text-sm text-[#CBD5E1]">
            <span>Privacy</span>
            <span>Terms</span>
            <span>Cookies</span>
          </div>
        </div>
      </Container>
    </footer>
  );
}
