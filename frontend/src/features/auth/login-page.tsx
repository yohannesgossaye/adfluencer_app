"use client";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { LoginForm, LoginFormHeader } from "@/components/auth/login-form";

export function LoginPageFeature() {
  return (
    <AuthPageCard>
      <LoginFormHeader />
      <LoginForm idPrefix="page-login" />
    </AuthPageCard>
  );
}
