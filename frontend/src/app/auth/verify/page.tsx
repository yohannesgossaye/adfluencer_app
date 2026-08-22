"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import axios from "axios";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { authApi } from "@/lib/api/auth";
import { useAuth } from "@/components/providers/auth-provider";

export default function VerifyPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const initialCode = searchParams.get("code") || searchParams.get("token") || "";
  const type = searchParams.get("type") || "email"; // 'email' or 'phone'

  const { refreshMe } = useAuth();
  const [code, setCode] = useState(initialCode);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isResending, setIsResending] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!code) {
      setError("Please enter the verification code.");
      return;
    }

    setError(null);
    setMessage(null);
    setIsLoading(true);

    try {
      if (type === "phone") {
        await authApi.verifyPhone(code.trim());
      } else {
        await authApi.verifyEmail(code.trim());
      }

      setMessage("Verification successful! Redirecting to dashboard...");
      await refreshMe();
      setTimeout(() => {
        router.push("/dashboard");
      }, 1200);
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response?.data?.message) {
        setError(err.response.data.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Verification failed. Please check the code or request a new one.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleResend = async () => {
    setError(null);
    setMessage(null);
    setIsResending(true);

    try {
      if (type === "phone") {
        await authApi.sendVerifyPhone();
      } else {
        await authApi.sendVerifyEmail();
      }
      setMessage("A new verification code has been sent!");
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response?.data?.message) {
        setError(err.response.data.message);
      } else {
        setError("Failed to resend verification code.");
      }
    } finally {
      setIsResending(false);
    }
  };

  useEffect(() => {
    if (initialCode) {
      // Auto-submit if code passed in URL
      setCode(initialCode);
    }
  }, [initialCode]);

  return (
    <AuthPageCard>
      <div className="space-y-4">
        <BrandLogo href="/" />
        <div>
          <h2 className="text-2xl font-semibold text-white">
            {type === "phone" ? "Verify your phone" : "Verify your email"}
          </h2>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            {type === "phone"
              ? "Enter the 6-digit SMS code sent to your phone."
              : "Enter the verification code sent to your email."}
          </p>
        </div>
      </div>

      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-2">
          <label htmlFor="verify-code" className="text-sm text-[#CBD5E1]">
            Verification Code
          </label>
          <Input
            id="verify-code"
            type="text"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
            placeholder="123456"
          />
        </div>

        <GradientButton className="h-12 w-full rounded-xl px-6" type="submit" disabled={isLoading}>
          {isLoading ? "Verifying..." : "Verify Code"}
        </GradientButton>

        <Button
          type="button"
          variant="outline"
          onClick={handleResend}
          disabled={isResending}
          className="h-11 w-full rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
        >
          {isResending ? "Sending..." : "Resend Code"}
        </Button>

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

        <div className="pt-2 text-center text-sm text-[#CBD5E1]">
          <Link href="/auth/login" className="text-white hover:text-[color:var(--adfluence-orange)]">
            Back to Login
          </Link>
        </div>
      </form>
    </AuthPageCard>
  );
}
