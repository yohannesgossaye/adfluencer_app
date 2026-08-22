"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import axios from "axios";

import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { useAuth } from "@/components/providers/auth-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { authApi } from "@/lib/api/auth";
import { setMockUser } from "@/lib/mock-auth";
import type { MarketingAuthRole } from "@/lib/marketing-auth-url";

type SignupFormProps = {
  idPrefix?: string;
  initialRole?: MarketingAuthRole;
  onSuccess?: () => void;
  onSwitchToLogin?: () => void;
  showCancel?: boolean;
  onCancel?: () => void;
};

export function SignupForm({
  idPrefix = "signup",
  initialRole = "business",
  onSuccess,
  onSwitchToLogin,
  showCancel = false,
  onCancel,
}: SignupFormProps) {
  const router = useRouter();
  const { setSession } = useAuth();
  const [role, setRole] = useState<MarketingAuthRole>(initialRole);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setRole(initialRole);
  }, [initialRole]);

  const roleLabel = useMemo(
    () => (role === "business" ? "Business name" : "Creator name"),
    [role],
  );

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !email || !phone || !password || !confirmPassword) {
      setError("Please fill in all required fields.");
      return;
    }
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }

    setError(null);
    setMessage(null);
    setIsLoading(true);

    try {
      const res = await authApi.signup({
        role,
        name: name.trim(),
        email: email.trim(),
        phone_number: phone.trim(),
        password,
        confirm_password: confirmPassword,
      });

      if (res.data) {
        const { user, profile, tokens, redirect_to } = res.data;
        if (tokens) {
          setSession(tokens, user, profile);
        }

        const displayName =
          profile?.business_name ||
          profile?.creator_name ||
          profile?.full_name ||
          user.email;

        setMockUser({
          email: user.email,
          role: user.role,
          name: displayName,
          redirect: redirect_to,
        });

        setMessage("Account created successfully! Redirecting...");
        onSuccess?.();
        setTimeout(() => {
          router.push(redirect_to || (role === "business" ? "/business/dashboard" : "/influencer/profile"));
        }, 1000);
      }
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response?.data?.message) {
        setError(err.response.data.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Failed to create account. Please try again.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="grid gap-3 sm:grid-cols-2">
        <button
          type="button"
          onClick={() => setRole("business")}
          className={cn(
            "rounded-xl border p-4 text-left transition-colors",
            role === "business"
              ? "border-[color:var(--adfluence-orange)] bg-white/10"
              : "border-white/10 bg-white/5 hover:bg-white/8",
          )}
        >
          <p className="text-sm font-semibold text-white">Business</p>
          <p className="mt-1 text-xs text-[#CBD5E1]">I want to launch campaigns</p>
        </button>
        <button
          type="button"
          onClick={() => setRole("influencer")}
          className={cn(
            "rounded-xl border p-4 text-left transition-colors",
            role === "influencer"
              ? "border-[color:var(--adfluence-orange)] bg-white/10"
              : "border-white/10 bg-white/5 hover:bg-white/8",
          )}
        >
          <p className="text-sm font-semibold text-white">Influencer</p>
          <p className="mt-1 text-xs text-[#CBD5E1]">
            I want to get brand collaborations
          </p>
        </button>
      </div>

      <div className="space-y-2">
        <label htmlFor={`${idPrefix}-name`} className="text-sm text-[#CBD5E1]">
          {roleLabel}
        </label>
        <Input
          id={`${idPrefix}-name`}
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
          placeholder={role === "business" ? "Acme Coffee" : "Selam Beauty Studio"}
        />
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <label htmlFor={`${idPrefix}-email`} className="text-sm text-[#CBD5E1]">
            Email
          </label>
          <Input
            id={`${idPrefix}-email`}
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
            placeholder="you@example.com"
          />
        </div>
        <div className="space-y-2">
          <label htmlFor={`${idPrefix}-phone`} className="text-sm text-[#CBD5E1]">
            Phone number
          </label>
          <Input
            id={`${idPrefix}-phone`}
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
            placeholder="+251 9xx xxx xxx"
          />
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <label htmlFor={`${idPrefix}-password`} className="text-sm text-[#CBD5E1]">
            Password
          </label>
          <Input
            id={`${idPrefix}-password`}
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
          />
        </div>
        <div className="space-y-2">
          <label
            htmlFor={`${idPrefix}-confirm-password`}
            className="text-sm text-[#CBD5E1]"
          >
            Confirm password
          </label>
          <Input
            id={`${idPrefix}-confirm-password`}
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
          />
        </div>
      </div>

      {showCancel ? (
        <div className="flex flex-col-reverse gap-3 sm:flex-row">
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={isLoading}
            className="h-12 flex-1 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
          >
            Cancel
          </Button>
          <GradientButton className="h-12 flex-1 rounded-xl px-6" type="submit" disabled={isLoading}>
            {isLoading ? "Creating Account..." : "Create Account"}
          </GradientButton>
        </div>
      ) : (
        <GradientButton className="h-12 w-full rounded-xl px-6" type="submit" disabled={isLoading}>
          {isLoading ? "Creating Account..." : "Create Account"}
        </GradientButton>
      )}

      {error ? (
        <p className="rounded-xl border border-red-400/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
          {error}
        </p>
      ) : null}

      {message ? (
        <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-200">
          {message}
        </p>
      ) : null}

      <div className="text-center text-sm text-[#CBD5E1]">
        Already have an account?{" "}
        {onSwitchToLogin ? (
          <button
            type="button"
            onClick={onSwitchToLogin}
            className="text-white hover:text-[color:var(--adfluence-orange)]"
          >
            Login
          </button>
        ) : (
          <Link
            href="/auth/login"
            className="text-white hover:text-[color:var(--adfluence-orange)]"
          >
            Login
          </Link>
        )}
      </div>
    </form>
  );
}

export function SignupFormHeader() {
  return (
    <div className="space-y-4">
      <BrandLogo href="/" />
      <div>
        <h2 className="text-2xl font-semibold text-white">
          Create your Adfluence account
        </h2>
        <p className="mt-2 text-sm text-[#CBD5E1]">
          Join Ethiopia&apos;s creator marketing platform.
        </p>
      </div>
    </div>
  );
}
