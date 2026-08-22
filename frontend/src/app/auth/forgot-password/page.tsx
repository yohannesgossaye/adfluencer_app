"use client";

import Link from "next/link";
import { useState } from "react";
import axios from "axios";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { authApi } from "@/lib/api/auth";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) {
      setError("Please enter your email address.");
      return;
    }

    setError(null);
    setMessage(null);
    setIsLoading(true);

    try {
      const res = await authApi.forgotPassword(email.trim());
      setMessage(res.message || "Password reset instructions have been sent to your email.");
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response?.data?.message) {
        setError(err.response.data.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Failed to process request. Please try again.");
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
          <h2 className="text-2xl font-semibold text-white">Forgot password</h2>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            Enter your email and we&apos;ll send you a link to reset your password.
          </p>
        </div>
      </div>

      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-2">
          <label htmlFor="forgot-email" className="text-sm text-[#CBD5E1]">
            Email
          </label>
          <Input
            id="forgot-email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
            placeholder="you@company.com"
          />
        </div>

        <GradientButton className="h-12 w-full rounded-xl px-6" type="submit" disabled={isLoading}>
          {isLoading ? "Sending..." : "Send Reset Link"}
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
