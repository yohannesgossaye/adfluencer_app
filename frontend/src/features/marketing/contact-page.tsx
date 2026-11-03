"use client";

import { useState } from "react";
import { CheckCircle2, Loader2, Mail } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { authApi } from "@/lib/api/auth";
import { authErrorMessage } from "@/lib/api-errors";
import { SITE_CONTACT_EMAIL } from "@/lib/site-contact";
import { cn } from "@/lib/utils";

type Role = "business" | "creator" | "other";
const ROLES: Array<{ id: Role; label: string }> = [
  { id: "business", label: "I'm a business" },
  { id: "creator", label: "I'm a creator" },
  { id: "other", label: "Something else" },
];

const field =
  "w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] dark:bg-[color:var(--adfluence-surface)] px-4 text-[15px] text-white placeholder:text-[#64748B] outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";

export function ContactPage() {
  const [form, setForm] = useState({ name: "", email: "", role: "business" as Role, subject: "", message: "", website: "" });
  const [state, setState] = useState<"idle" | "sending" | "sent">("idle");
  const [error, setError] = useState("");
  const set = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    setForm((f) => ({ ...f, [k]: e.target.value }));

  const problem = !form.name.trim()
    ? "Please enter your name."
    : !/^\S+@\S+\.\S+$/.test(form.email.trim())
      ? "Please enter a valid email address, so we can reply."
      : form.message.trim().length < 10
        ? "Please write a message of at least 10 characters."
        : "";

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (problem) {
      setError(problem);
      return;
    }
    setError("");
    setState("sending");
    try {
      await authApi.contact(form);
      setState("sent");
    } catch (err) {
      setError(authErrorMessage(err, "Your message could not be sent. Please try again, or email us directly."));
      setState("idle");
    }
  };

  return (
    <section className="mx-auto grid max-w-[1200px] gap-10 px-5 py-12 sm:px-8 lg:grid-cols-[minmax(0,1fr)_minmax(0,560px)] lg:gap-20 lg:py-20">
      <div>
        <h1 className="font-[family-name:var(--font-display)] text-[40px] leading-[1.05] font-extrabold tracking-[-0.03em] text-white sm:text-6xl">
          Contact us
        </h1>
        <p className="mt-5 max-w-[460px] text-[17px] leading-relaxed text-[#CBD5E1]">
          Questions about campaigns, creators, contracts or fees? Send us a message and we&apos;ll reply by email within
          one working day, in Amharic or English.
        </p>
        <div className="mt-8 rounded-2xl border border-white/10 bg-[color:var(--adfluence-card)] p-5">
          <p className="text-sm text-[#94A3B8]">Prefer your own email app?</p>
          <a
            href={`mailto:${SITE_CONTACT_EMAIL}`}
            className="mt-1 inline-flex items-center gap-2 text-[15px] font-semibold text-white hover:text-[color:var(--adfluence-orange)]"
          >
            <Mail className="size-4 shrink-0" aria-hidden /> Email us directly
          </a>
        </div>
      </div>

      <div className="rounded-[24px] border border-white/10 bg-[color:var(--adfluence-card)] p-6 sm:p-8">
        {state === "sent" ? (
          <div role="status" className="flex flex-col items-start gap-4 py-6">
            <CheckCircle2 className="size-10 text-emerald-300" aria-hidden />
            <h2 className="font-[family-name:var(--font-display)] text-2xl font-bold text-white">Thanks, your message is on its way</h2>
            <p className="text-[15px] leading-relaxed text-[#CBD5E1]">
              We&apos;ll reply to <span className="font-medium text-white">{form.email.trim()}</span> within one working day.
            </p>
            <Button
              type="button"
              variant="secondary"
              className="h-11 rounded-xl border border-white/15 px-5"
              onClick={() => {
                setForm((f) => ({ ...f, subject: "", message: "" }));
                setState("idle");
              }}
            >
              Send another message
            </Button>
          </div>
        ) : (
          <form onSubmit={submit} noValidate className="space-y-5">
            <fieldset>
              <legend className="mb-2 text-sm text-[#CBD5E1]">Who are you?</legend>
              <div className="flex flex-wrap gap-2">
                {ROLES.map((r) => (
                  <label
                    key={r.id}
                    className={cn(
                      "inline-flex min-h-11 cursor-pointer items-center rounded-full border px-4 text-sm transition-colors has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-[#A78BFA]",
                      form.role === r.id
                        ? "border-transparent bg-[color:var(--adfluence-action)] font-semibold text-white"
                        : "border-white/15 text-[#CBD5E1] hover:text-white",
                    )}
                  >
                    <input
                      type="radio"
                      name="role"
                      value={r.id}
                      checked={form.role === r.id}
                      onChange={() => setForm((f) => ({ ...f, role: r.id }))}
                      className="sr-only"
                    />
                    {r.label}
                  </label>
                ))}
              </div>
            </fieldset>
            <div className="grid gap-5 sm:grid-cols-2">
              <label className="block space-y-2">
                <span className="text-sm text-[#CBD5E1]">Your name</span>
                <Input value={form.name} onChange={set("name")} autoComplete="name" maxLength={100} className={cn(field, "h-12")} />
              </label>
              <label className="block space-y-2">
                <span className="text-sm text-[#CBD5E1]">Your email</span>
                <Input type="email" value={form.email} onChange={set("email")} autoComplete="email" inputMode="email" className={cn(field, "h-12")} />
              </label>
            </div>
            <label className="block space-y-2">
              <span className="text-sm text-[#CBD5E1]">
                Subject <span className="text-[#64748B]">(optional)</span>
              </span>
              <Input value={form.subject} onChange={set("subject")} maxLength={150} className={cn(field, "h-12")} />
            </label>
            <label className="block space-y-2">
              <span className="text-sm text-[#CBD5E1]">Message</span>
              <textarea value={form.message} onChange={set("message")} rows={6} maxLength={3000} className={cn(field, "resize-y py-3")} />
            </label>
            {/* Hidden from people; bots that fill it in are ignored. */}
            <input type="text" name="website" value={form.website} onChange={set("website")} tabIndex={-1} autoComplete="off" aria-hidden className="hidden" />
            {error ? (
              <p role="alert" className="rounded-xl border border-rose-400/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-100">
                {error}
              </p>
            ) : null}
            <Button type="submit" disabled={state === "sending"} className="h-12 w-full rounded-xl text-[15px] font-semibold">
              {state === "sending" ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
              Send message
            </Button>
          </form>
        )}
      </div>
    </section>
  );
}
