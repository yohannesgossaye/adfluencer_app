"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

import { GradientButton } from "@/components/common/gradient-button";
import { useAuth } from "@/components/providers/auth-provider";
import { mockUsers } from "@/data/mock-users";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { isEmailBlocked } from "@/lib/blocked-users";
import { setMockUser } from "@/lib/mock-auth";

export function LoginPageFeature() {
  const router = useRouter();
  const { setAccessToken } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError("Incorrect email or password.");
      return;
    }
    const matchedUser = mockUsers.find(
      (u) =>
        u.email.toLowerCase() === email.trim().toLowerCase() &&
        u.password === password,
    );
    if (!matchedUser) {
      setError("Incorrect email or password.");
      return;
    }
    if (isEmailBlocked(matchedUser.email)) {
      setError("This account has been blocked. Please contact Adfluence support.");
      return;
    }
    setError(null);
    setAccessToken("demo-session-token");
    setMockUser({
      email: matchedUser.email,
      role: matchedUser.role,
      name: matchedUser.name,
      redirect: matchedUser.redirect,
    });
    router.push(matchedUser.redirect);
  };

  return (
    <Card className="w-full max-w-md border-white/15 bg-[#0B1A2B]/70 text-white shadow-[0_22px_60px_rgba(6,17,31,0.45)] backdrop-blur-xl">
      <CardHeader className="space-y-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED]">
            <span className="text-xl font-bold text-white">A</span>
          </div>
          <div>
            <p className="text-lg font-semibold text-white">Adfluence</p>
            <p className="text-xs text-[#CBD5E1]">Connect. Collaborate. Grow.</p>
          </div>
        </div>
        <div>
          <CardTitle className="text-2xl text-white">Welcome back</CardTitle>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            Log in to manage campaigns, creators, and collaborations.
          </p>
        </div>
      </CardHeader>

      <CardContent>
        <form onSubmit={onSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="login-email" className="text-sm text-[#CBD5E1]">
              Email
            </label>
            <Input
              id="login-email"
              type="email"
              autoComplete="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
              placeholder="you@company.com"
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="login-password" className="text-sm text-[#CBD5E1]">
              Password
            </label>
            <Input
              id="login-password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
              placeholder="••••••••"
            />
          </div>

          <div className="flex flex-wrap items-center justify-between gap-3 text-sm">
            <label className="flex items-center gap-2 text-[#CBD5E1]">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                className="h-4 w-4 rounded border-white/25 bg-transparent accent-[#7C3AED]"
              />
              Remember me
            </label>
            <Link href="/auth/forgot-password" className="text-[#CBD5E1] hover:text-white">
              Forgot password?
            </Link>
          </div>

          <GradientButton className="h-12 w-full rounded-xl px-6" type="submit">
            Login
          </GradientButton>

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
            <Link href="/auth/signup" className="text-white hover:text-[#FF7A00]">
              Sign up
            </Link>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

