import { Container } from "@/components/common/container";

export function MarketingFooter() {
  return (
    <footer className="bg-[color:var(--adfluence-navy-dark)] border-t border-white/10 py-16">
      <Container>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12 text-white">
          <div>
            <h4 className="font-semibold mb-4">Platform</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>For Businesses</li>
              <li>For Influencers</li>
              <li>How It Works</li>
              <li>Pricing</li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold mb-4">Company</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>About Us</li>
              <li>Careers</li>
              <li>Blog</li>
              <li>Press Kit</li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold mb-4">Resources</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>Help Center</li>
              <li>Community</li>
              <li>Case Studies</li>
              <li>API Docs</li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold mb-4">Contact</h4>
            <ul className="space-y-2 text-sm text-[#CBD5E1]">
              <li>hello@adfluence.et</li>
              <li>Addis Ababa, Ethiopia</li>
              <li>adfluence.et</li>
            </ul>
          </div>
        </div>

        <div className="pt-8 border-t border-white/10 flex flex-col md:flex-row justify-between items-center gap-4 text-[#CBD5E1]">
          <div className="flex items-center gap-3 text-white">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED] flex items-center justify-center">
              <span className="text-white font-bold text-xl">A</span>
            </div>
            <span className="font-bold text-xl">Adfluence</span>
          </div>
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
