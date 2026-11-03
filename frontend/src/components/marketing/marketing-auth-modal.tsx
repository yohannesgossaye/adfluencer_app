"use client";

import { LoginForm, LoginFormHeader } from "@/components/auth/login-form";
import { SignupForm, SignupFormHeader } from "@/components/auth/signup-form";
import { Dialog, DialogContent, DialogDescription, DialogTitle } from "@/components/ui/dialog";

import { useSearchParams } from "next/navigation";

import { parseMarketingAuthRole, type MarketingAuthMode } from "@/lib/marketing-auth-url";

import { useMarketingAuth } from "./marketing-auth-provider";

/** Render inside <Suspense>: it reads the URL (?auth=login|signup&role=…). */
export function MarketingAuthModal() {
  const { openLogin, openSignup, closeAuth } = useMarketingAuth();
  const searchParams = useSearchParams();
  const authParam = searchParams.get("auth");
  const mode: MarketingAuthMode | null = authParam === "login" || authParam === "signup" ? authParam : null;
  const signupRole = parseMarketingAuthRole(searchParams.get("role")) ?? "business";
  const isOpen = mode !== null;
  const switchToLogin = openLogin;
  const switchToSignup = openSignup;

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(open) => {
        if (!open) {
          closeAuth();
        }
      }}
    >
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-xl" showCloseButton>
        <DialogTitle className="sr-only">Authentication</DialogTitle>
        <DialogDescription className="sr-only">Log in or create an account on Adfluence</DialogDescription>
        <div className="space-y-6 p-6 pr-14">
          {mode === "login" ? (
            <>
              <LoginFormHeader />
              <LoginForm
                idPrefix="modal-login"
                showCancel
                onCancel={closeAuth}
                onSwitchToSignup={() => switchToSignup(signupRole)}
              />
            </>
          ) : null}
          {mode === "signup" ? (
            <>
              <SignupFormHeader />
              <SignupForm
                idPrefix="modal-signup"
                initialRole={signupRole}
                showCancel
                onCancel={closeAuth}
                onSwitchToLogin={switchToLogin}
              />
            </>
          ) : null}
        </div>
      </DialogContent>
    </Dialog>
  );
}
