import Link from "next/link";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { BrandLogo } from "@/components/common/brand-logo";
import { Button } from "@/components/ui/button";

export default function ResetPasswordPage() {
  return (
    <AuthPageCard>
      <div className="space-y-4">
        <BrandLogo href="/" />
        <div>
          <h2 className="text-2xl font-semibold text-white">Reset password</h2>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            Choose a new password using the secure link from your email.
          </p>
        </div>
      </div>
      <Button
        asChild
        variant="outline"
        className="h-11 w-full rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
      >
        <Link href="/auth/login">Back to sign in</Link>
      </Button>
    </AuthPageCard>
  );
}
