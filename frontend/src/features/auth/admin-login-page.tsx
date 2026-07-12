"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { AuthPageCard } from "@/components/auth/auth-page-card";
import { BrandLogo } from "@/components/common/brand-logo";
import { GradientButton } from "@/components/common/gradient-button";
import { useAuth } from "@/components/providers/auth-provider";
import { mockUsers } from "@/data/mock-users";
import { Input } from "@/components/ui/input";
import { isEmailBlocked } from "@/lib/blocked-users";
import { setMockUser } from "@/lib/mock-auth";

export function AdminLoginPageFeature() {
  const router = useRouter();
  const { setAccessToken } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const admin = mockUsers.find((u) => u.role === "admin");
    if (
      !admin ||
      email.trim().toLowerCase() !== admin.email.toLowerCase() ||
      password !== admin.password
    ) {
      setError("Invalid admin credentials.");
      return;
    }
    if (isEmailBlocked(admin.email)) {
      setError("This account has been blocked. Please contact Adfluence support.");
      return;
    }
    setError(null);
    setAccessToken("demo-admin-token");
    setMockUser({
      email: admin.email,
      role: admin.role,
      name: admin.name,
      redirect: admin.redirect,
    });
    router.push(admin.redirect);
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
          <GradientButton className="h-12 w-full rounded-xl px-6" type="submit">
            Login
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
