import Link from "next/link";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { BrandLogo } from "@/components/common/brand-logo";
import { Button } from "@/components/ui/button";

export default function ForgotPasswordPage() {
  return (
    <AuthPageCard>
      <div className="space-y-4">
        <BrandLogo href="/" />
        <div>
          <h2 className="text-2xl font-semibold text-white">Forgot password</h2>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            Enter your email and we&apos;ll send you a link to reset your password.
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
