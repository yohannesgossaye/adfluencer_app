"use client";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { SignupForm, SignupFormHeader } from "@/components/auth/signup-form";
import { parseMarketingAuthRole } from "@/lib/marketing-auth-url";

export function SignupPageFeature({ initialRole }: { initialRole?: string }) {
  const role = parseMarketingAuthRole(initialRole) ?? "business";

  return (
    <AuthPageCard size="lg">
      <SignupFormHeader />
      <SignupForm idPrefix="page-signup" initialRole={role} />
    </AuthPageCard>
  );
}
