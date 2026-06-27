"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { GradientButton } from "@/components/common/gradient-button";
import { useAuth } from "@/components/providers/auth-provider";
import { mockUsers } from "@/data/mock-users";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
      <Card className="w-full border-white/15 bg-[#06111F]/85 text-white shadow-[0_22px_60px_rgba(6,17,31,0.55)] backdrop-blur-xl">
        <CardHeader className="space-y-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#FF7A00] to-[#7C3AED]">
              <span className="text-xl font-bold text-white">A</span>
            </div>
            <div>
              <p className="text-lg font-semibold text-white">Adfluence</p>
              <p className="text-xs text-[#CBD5E1]">Admin access</p>
            </div>
          </div>
          <div>
            <CardTitle className="text-2xl text-white">Admin Portal</CardTitle>
            <p className="mt-2 text-sm text-[#CBD5E1]">
              Secure access for Adfluence platform management.
            </p>
          </div>
        </CardHeader>
        <CardContent>
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
                className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
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
                className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
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
        </CardContent>
      </Card>
    </div>
  );
}

