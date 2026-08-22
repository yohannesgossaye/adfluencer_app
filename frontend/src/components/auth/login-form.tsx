"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";
import axios from "axios";

import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { useAuth } from "@/components/providers/auth-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { authApi } from "@/lib/api/auth";
import { setMockUser } from "@/lib/mock-auth";

type LoginFormProps = {
  idPrefix?: string;
  onSuccess?: () => void;
  onSwitchToSignup?: () => void;
  showCancel?: boolean;
  onCancel?: () => void;
};

export function LoginForm({
  idPrefix = "login",
  onSuccess,
  onSwitchToSignup,
  showCancel = false,
  onCancel,
}: LoginFormProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const nextParam = searchParams.get("next") || undefined;
  const { setSession } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError("Please enter your email and password.");
      return;
    }

    setError(null);
    setIsLoading(true);

    try {
      const res = await authApi.login({
        email: email.trim(),
        password,
        remember_me: rememberMe,
        next: nextParam,
      });

      if (res.data) {
        const { user, profile, tokens, redirect_to } = res.data;
        setSession(tokens, user, profile);

        const displayName =
          profile.business_name ||
          profile.creator_name ||
          profile.full_name ||
          user.email;

        setMockUser({
          email: user.email,
          role: user.role,
          name: displayName,
          redirect: redirect_to,
        });

        onSuccess?.();
        const targetRoute = nextParam || redirect_to || "/dashboard";
        router.push(targetRoute);
      }
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response?.data?.message) {
        setError(err.response.data.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Login failed. Please check your credentials.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="space-y-2">
        <label htmlFor={`${idPrefix}-email`} className="text-sm text-[#CBD5E1]">
          Email
        </label>
        <Input
          id={`${idPrefix}-email`}
          type="email"
          autoComplete="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
          placeholder="you@company.com"
        />
      </div>
      <div className="space-y-2">
        <label htmlFor={`${idPrefix}-password`} className="text-sm text-[#CBD5E1]">
          Password
        </label>
        <Input
          id={`${idPrefix}-password`}
          type="password"
          autoComplete="current-password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
          placeholder="••••••••"
        />
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3 text-sm">
        <label className="flex items-center gap-2 text-[#CBD5E1]">
          <input
            type="checkbox"
            checked={rememberMe}
            onChange={(e) => setRememberMe(e.target.checked)}
            className="h-4 w-4 rounded border-white/25 bg-transparent accent-[color:var(--adfluence-orange)]"
          />
          Remember me
        </label>
        <Link href="/auth/forgot-password" className="text-[#CBD5E1] hover:text-white">
          Forgot password?
        </Link>
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
            {isLoading ? "Logging in..." : "Login"}
          </GradientButton>
        </div>
      ) : (
        <GradientButton className="h-12 w-full rounded-xl px-6" type="submit" disabled={isLoading}>
          {isLoading ? "Logging in..." : "Login"}
        </GradientButton>
      )}

      {error ? (
        <p className="rounded-xl border border-red-400/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
          {error}
        </p>
      ) : null}

      <p className="text-center text-xs text-[#CBD5E1]">
        Business, Influencer, and Admin accounts use the same login.
      </p>

      <div className="pt-1 text-center text-sm text-[#CBD5E1]">
        Don&apos;t have an account?{" "}
        {onSwitchToSignup ? (
          <button
            type="button"
            onClick={onSwitchToSignup}
            className="text-white hover:text-[color:var(--adfluence-orange)]"
          >
            Sign up
          </button>
        ) : (
          <Link
            href="/auth/signup"
            className="text-white hover:text-[color:var(--adfluence-orange)]"
          >
            Sign up
          </Link>
        )}
      </div>
    </form>
  );
}

export function LoginFormHeader() {
  return (
    <div className="space-y-4">
      <BrandLogo href="/" />
      <div>
        <h2 className="text-2xl font-semibold text-white">Welcome back</h2>
        <p className="mt-2 text-sm text-[#CBD5E1]">
          Log in to manage campaigns, creators, and collaborations.
        </p>
      </div>
    </div>
  );
}
