"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";

import { GradientButton } from "@/components/common/gradient-button";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

type Role = "business" | "influencer";

function normalizeInitialRole(value: string | undefined): Role {
  return value === "influencer" ? "influencer" : "business";
}

export function SignupPageFeature({ initialRole }: { initialRole?: string }) {
  const router = useRouter();
  const [role, setRole] = useState<Role>(() => normalizeInitialRole(initialRole));
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [message, setMessage] = useState<string | null>(null);

  const roleLabel = useMemo(
    () => (role === "business" ? "Business name" : "Creator name"),
    [role],
  );

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !email || !phone || !password || !confirmPassword) return;
    setMessage("Demo signup complete. Please use a mock login account to continue.");
    setTimeout(() => {
      router.push("/auth/login");
    }, 1200);
  };

  return (
    <Card className="w-full max-w-xl border-white/15 bg-[#0B1A2B]/70 text-white shadow-[0_22px_60px_rgba(6,17,31,0.45)] backdrop-blur-xl">
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
          <CardTitle className="text-2xl text-white">Create your Adfluence account</CardTitle>
          <p className="mt-2 text-sm text-[#CBD5E1]">
            Join Ethiopia&apos;s creator marketing platform.
          </p>
        </div>
      </CardHeader>

      <CardContent>
        <form onSubmit={onSubmit} className="space-y-4">
          <div className="grid gap-3 sm:grid-cols-2">
            <button
              type="button"
              onClick={() => setRole("business")}
              className={cn(
                "rounded-xl border p-4 text-left transition-colors",
                role === "business"
                  ? "border-[#7C3AED] bg-white/10"
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
                  ? "border-[#FF7A00] bg-white/10"
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
            <label htmlFor="signup-name" className="text-sm text-[#CBD5E1]">
              {roleLabel}
            </label>
            <Input
              id="signup-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
              placeholder={role === "business" ? "Acme Coffee" : "Selam Beauty Studio"}
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <label htmlFor="signup-email" className="text-sm text-[#CBD5E1]">
                Email
              </label>
              <Input
                id="signup-email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
                placeholder="you@example.com"
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="signup-phone" className="text-sm text-[#CBD5E1]">
                Phone number
              </label>
              <Input
                id="signup-phone"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
                placeholder="+251 9xx xxx xxx"
              />
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <label htmlFor="signup-password" className="text-sm text-[#CBD5E1]">
                Password
              </label>
              <Input
                id="signup-password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="signup-confirm-password" className="text-sm text-[#CBD5E1]">
                Confirm password
              </label>
              <Input
                id="signup-confirm-password"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className="h-12 rounded-xl border-white/15 bg-white/5 text-white placeholder:text-[#CBD5E1]/70"
              />
            </div>
          </div>

          <GradientButton className="h-12 w-full rounded-xl px-6" type="submit">
            Create Account
          </GradientButton>

          {message ? (
            <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-200">
              {message}
            </p>
          ) : null}

          <div className="text-center text-sm text-[#CBD5E1]">
            Already have an account?{" "}
            <Link href="/auth/login" className="text-white hover:text-[#FF7A00]">
              Login
            </Link>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

