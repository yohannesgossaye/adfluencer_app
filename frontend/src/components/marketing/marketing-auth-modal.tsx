"use client";

import { LoginForm, LoginFormHeader } from "@/components/auth/login-form";
import { SignupForm, SignupFormHeader } from "@/components/auth/signup-form";
import { Dialog, DialogContent, DialogDescription, DialogTitle } from "@/components/ui/dialog";

import { useMarketingAuth } from "./marketing-auth-provider";

export function MarketingAuthModal() {
  const {
    mode,
    signupRole,
    isOpen,
    closeAuth,
    switchToLogin,
    switchToSignup,
  } = useMarketingAuth();

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
