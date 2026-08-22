"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";
import axios from "axios";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { authApi } from "@/lib/api/auth";

export default function ResetPasswordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get("token") || "";

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      setError("Reset token is missing or invalid. Please request a new link.");
      return;
    }
    if (!password || !confirmPassword) {
      setError("Please fill in both password fields.");
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
      const res = await authApi.resetPassword(password, token);
      setMessage(res.message || "Password has been successfully reset! Redirecting to login...");
      setTimeout(() => {
        router.push("/auth/login");
      }, 1500);
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response?.data?.message) {
        setError(err.response.data.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Failed to reset password. The link may be expired.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthPageCard>
      <div className="space-y-4">
        <BrandLogo href="/" />
        <div>
          <h2 className="text-2xl font-semibold text-white">Reset password</h2>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            Choose a new secure password for your account.
          </p>
        </div>
      </div>

      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-2">
          <label htmlFor="reset-password" className="text-sm text-[#CBD5E1]">
            New Password
          </label>
          <Input
            id="reset-password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
            placeholder="••••••••"
          />
        </div>

        <div className="space-y-2">
          <label htmlFor="reset-confirm-password" className="text-sm text-[#CBD5E1]">
            Confirm New Password
          </label>
          <Input
            id="reset-confirm-password"
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
            placeholder="••••••••"
          />
        </div>

        <GradientButton className="h-12 w-full rounded-xl px-6" type="submit" disabled={isLoading}>
          {isLoading ? "Resetting..." : "Reset Password"}
        </GradientButton>

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

        <Button
          asChild
          variant="outline"
          className="h-11 w-full rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
        >
          <Link href="/auth/login">Back to sign in</Link>
        </Button>
      </form>
    </AuthPageCard>
  );
}
