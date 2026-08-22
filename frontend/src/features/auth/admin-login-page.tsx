"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import axios from "axios";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { useAuth } from "@/components/providers/auth-provider";
import { Input } from "@/components/ui/input";
import { authApi } from "@/lib/api/auth";
import { setMockUser } from "@/lib/mock-auth";

export function AdminLoginPageFeature() {
  const router = useRouter();
  const { setSession } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError("Please enter your admin email and password.");
      return;
    }

    setError(null);
    setIsLoading(true);

    try {
      const res = await authApi.login({
        email: email.trim(),
        password,
        portal: "admin",
      });

      if (res.data) {
        const { user, profile, tokens, redirect_to } = res.data;
        if (user.role !== "admin") {
          setError("Admin access required.");
          return;
        }

        setSession(tokens, user, profile);
        setMockUser({
          email: user.email,
          role: "admin",
          name: profile.full_name || "Adfluence Admin",
          redirect: redirect_to || "/admin/dashboard",
        });

        router.push(redirect_to || "/admin/dashboard");
      }
    } catch (err: unknown) {
      if (axios.isAxiosError(err) && err.response?.data?.message) {
        setError(err.response.data.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Invalid admin credentials.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="mx-auto flex min-h-[80vh] w-full max-w-lg items-center px-4">
      <AuthPageCard>
        <div className="space-y-4">
          <BrandLogo href="/" tagline="Admin access" />
          <div>
            <h2 className="text-2xl font-semibold text-white">Admin Portal</h2>
            <p className="mt-2 text-sm text-[#CBD5E1]">
              Secure access for Adfluence platform management.
            </p>
          </div>
        </div>
        <form onSubmit={onSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="admin-email" className="text-sm text-[#CBD5E1]">
              Email
            </label>
            <Input
              id="admin-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
              placeholder="admin@adfluence.et"
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="admin-password" className="text-sm text-[#CBD5E1]">
              Password
            </label>
            <Input
              id="admin-password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="h-12 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#CBD5E1]/70"
              placeholder="••••••••"
            />
          </div>
          <GradientButton className="h-12 w-full rounded-xl px-6" type="submit" disabled={isLoading}>
            {isLoading ? "Logging in..." : "Login"}
          </GradientButton>
          {error ? (
            <p className="rounded-xl border border-red-400/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
              {error}
            </p>
          ) : null}
        </form>
      </AuthPageCard>
    </div>
  );
}
