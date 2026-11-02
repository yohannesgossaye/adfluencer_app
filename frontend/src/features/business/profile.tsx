"use client";

import { useEffect, useMemo, useRef, useState, type ReactNode } from "react";
import axios from "axios";
import { Loader2, Upload } from "lucide-react";
import { toast } from "sonner";

import { VerificationStatusCard } from "@/components/business/verification-status-card";
import { UserAvatar } from "@/components/common/user-avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  businessApi,
  type BusinessProfileResponse,
  type UpdateBusinessProfileRequest,
} from "@/lib/api/business";
import { cn } from "@/lib/utils";

const PLATFORMS = [
  { id: "tiktok", label: "TikTok" },
  { id: "instagram", label: "Instagram" },
  { id: "youtube", label: "YouTube" },
  { id: "facebook", label: "Facebook" },
  { id: "telegram", label: "Telegram" },
] as const;

const BUSINESS_TYPES = ["Sole proprietorship", "PLC", "Share company", "Partnership", "NGO", "Other"];
const BUSINESS_SIZES = ["1–10 employees", "11–50 employees", "51–200 employees", "200+ employees"];

const MISSING_LABELS: Record<string, string> = {
  owner_full_name: "Owner full name",
  business_name: "Business name",
  description: "Description (at least 20 characters)",
  industry: "Industry",
  city: "City",
  preferred_platforms: "Preferred platforms",
  logo_url: "Logo",
  website_or_social: "Website or a social link",
  business_type: "Business type",
  monthly_marketing_budget: "Monthly marketing budget",
};

const MAX_LOGO_BYTES = 5 * 1024 * 1024;

type Form = {
  owner_full_name: string;
  business_name: string;
  legal_business_name: string;
  tin_number: string;
  business_type: string;
  business_size: string;
  industry: string;
  description: string;
  monthly_marketing_budget: string;
  city: string;
  region: string;
  address: string;
  website_url: string;
  instagram: string;
  tiktok: string;
  facebook: string;
  linkedin: string;
  preferred_platforms: string[];
};

function toForm(p: BusinessProfileResponse["profile"]): Form {
  return {
    owner_full_name: p.owner_full_name ?? "",
    business_name: p.business_name ?? "",
    legal_business_name: p.legal_business_name ?? "",
    tin_number: p.tin_number ?? "",
    business_type: p.business_type ?? "",
    business_size: p.business_size ?? "",
    industry: p.industry ?? "",
    description: p.description ?? "",
    monthly_marketing_budget: p.monthly_marketing_budget != null ? String(p.monthly_marketing_budget) : "",
    city: p.city ?? "",
    region: p.region ?? "",
    address: p.address ?? "",
    website_url: p.website_url ?? "",
    instagram: p.social_links?.instagram ?? "",
    tiktok: p.social_links?.tiktok ?? "",
    facebook: p.social_links?.facebook ?? "",
    linkedin: p.social_links?.linkedin ?? "",
    preferred_platforms: p.preferred_platforms ?? [],
  };
}

function toRequest(f: Form): UpdateBusinessProfileRequest {
  const opt = (v: string) => (v.trim() ? v.trim() : undefined);
  const budget = f.monthly_marketing_budget.trim();
  return {
    owner_full_name: f.owner_full_name.trim(),
    business_name: f.business_name.trim(),
    legal_business_name: opt(f.legal_business_name),
    tin_number: opt(f.tin_number),
    business_type: opt(f.business_type),
    business_size: opt(f.business_size),
    industry: f.industry.trim(),
    description: f.description.trim(),
    monthly_marketing_budget: budget ? Number(budget) : undefined,
    city: f.city.trim(),
    region: opt(f.region),
    address: opt(f.address),
    website_url: opt(f.website_url),
    social_links: {
      instagram: opt(f.instagram),
      tiktok: opt(f.tiktok),
      facebook: opt(f.facebook),
      linkedin: opt(f.linkedin),
    },
    preferred_platforms: f.preferred_platforms,
  };
}

function validate(f: Form): string | null {
  if (f.owner_full_name.trim().length < 2) return "Enter the owner’s full name.";
  if (f.business_name.trim().length < 2) return "Enter your business name.";
  if (f.industry.trim().length < 2) return "Enter your industry.";
  if (f.description.trim().length < 20) return "Describe your business in at least 20 characters.";
  if (f.city.trim().length < 2) return "Enter your city.";
  if (f.preferred_platforms.length === 0) return "Choose at least one platform you want to advertise on.";
  const budget = f.monthly_marketing_budget.trim();
  if (budget && (!Number.isFinite(Number(budget)) || Number(budget) < 0)) return "Monthly budget must be a positive number.";
  for (const [label, value] of [
    ["Website", f.website_url],
    ["Instagram", f.instagram],
    ["TikTok", f.tiktok],
    ["Facebook", f.facebook],
    ["LinkedIn", f.linkedin],
  ] as const) {
    if (value.trim() && !/^https?:\/\/\S+\.\S+/.test(value.trim())) return `${label} must be a full link starting with https://`;
  }
  return null;
}

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { message?: string; error?: { details?: Array<{ description?: string }> } } | undefined;
    const detail = data?.error?.details?.[0]?.description;
    if (detail && detail !== data?.message) return detail;
    if (data?.message) return data.message;
    if (!err.response) return "You appear to be offline. Check your connection.";
  }
  return fallback;
}

const inputClass =
  "h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25";
const selectClass =
  "h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";

function Field({ label, hint, children, wide }: { label: string; hint?: string; children: ReactNode; wide?: boolean }) {
  return (
    <label className={cn("space-y-2", wide && "md:col-span-2")}>
      <span className="text-sm text-[#CBD5E1]">{label}</span>
      {children}
      {hint ? <span className="block text-xs text-[#94A3B8]">{hint}</span> : null}
    </label>
  );
}

function Section({ title, description, children }: { title: string; description?: string; children: ReactNode }) {
  return (
    <section className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
      <div>
        <h2 className="text-lg font-semibold">{title}</h2>
        {description ? <p className="mt-0.5 text-xs text-[#94A3B8]">{description}</p> : null}
      </div>
      <div className="grid gap-4 md:grid-cols-2">{children}</div>
    </section>
  );
}

export function BusinessProfile() {
  const [data, setData] = useState<BusinessProfileResponse | null>(null);
  const [form, setForm] = useState<Form | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [confirmIdentityChange, setConfirmIdentityChange] = useState(false);
  const fileRef = useRef<HTMLInputElement | null>(null);

  const load = async () => {
    try {
      const res = await businessApi.getProfile();
      if (res.data) {
        setData(res.data);
        setForm(toForm(res.data.profile));
        setLoadError(null);
      }
    } catch (err) {
      setLoadError(apiError(err, "Could not load your business profile."));
    }
  };

  useEffect(() => {
    void Promise.resolve().then(load);
  }, []);

  const saved = useMemo(() => (data ? toForm(data.profile) : null), [data]);
  const dirty = useMemo(() => !!form && !!saved && JSON.stringify(form) !== JSON.stringify(saved), [form, saved]);
  const verification = data?.profile.verification;
  const identityChanged =
    !!form &&
    !!saved &&
    (form.business_name.trim() !== saved.business_name ||
      form.legal_business_name.trim() !== saved.legal_business_name ||
      form.tin_number.trim() !== saved.tin_number);

  const set = <K extends keyof Form>(key: K, value: Form[K]) => setForm((f) => (f ? { ...f, [key]: value } : f));

  const save = async () => {
    if (!form) return;
    const problem = validate(form);
    if (problem) {
      toast.error(problem);
      return;
    }
    if (verification?.status === "verified" && identityChanged && !confirmIdentityChange) {
      setConfirmIdentityChange(true);
      return;
    }
    setSaving(true);
    try {
      const res = await businessApi.updateProfile(toRequest(form));
      if (res.data) {
        setData(res.data);
        setForm(toForm(res.data.profile));
      }
      setConfirmIdentityChange(false);
      toast.success(
        res.data?.profile.verification.status === "pending" && verification?.status === "verified"
          ? "Profile saved. Your business is being re-verified."
          : "Profile saved.",
      );
    } catch (err) {
      toast.error(apiError(err, "Could not save your profile."));
    } finally {
      setSaving(false);
    }
  };

  const submitForVerification = async () => {
    setSubmitting(true);
    try {
      const res = await businessApi.completeProfile();
      toast.success(res.data?.message ?? "Submitted for verification.");
      await load();
    } catch (err) {
      toast.error(apiError(err, "Could not submit your profile."));
    } finally {
      setSubmitting(false);
    }
  };

  const uploadLogo = async (file: File) => {
    if (!["image/jpeg", "image/png", "image/webp"].includes(file.type)) {
      toast.error("Logo must be a JPEG, PNG or WebP image.");
      return;
    }
    if (file.size > MAX_LOGO_BYTES) {
      toast.error("Logo must be 5 MB or smaller.");
      return;
    }
    setUploading(true);
    try {
      await businessApi.uploadLogo(file);
      toast.success("Logo updated.");
      await load();
    } catch (err) {
      toast.error(apiError(err, "Could not upload the logo."));
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = "";
    }
  };

  if (loadError) {
    return (
      <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
        <p>{loadError}</p>
        <Button type="button" variant="link" className="mt-2 h-auto p-0 text-white" onClick={() => void load()}>
          Try again
        </Button>
      </div>
    );
  }

  if (!data || !form || !verification) {
    return (
      <p className="flex items-center gap-2 py-10 text-sm text-[#94A3B8]">
        <Loader2 className="size-4 animate-spin" aria-hidden /> Loading your business profile…
      </p>
    );
  }

  const canSubmit = verification.status === "unsubmitted" || verification.status === "rejected";
  const missing = data.completeness.missing.map((m) => MISSING_LABELS[m] ?? m);

  return (
    <div className="space-y-6 pb-24 text-white">
      <div>
        <h1 className="text-2xl font-semibold">Business profile</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Creators see this information. Our team uses it to verify your business.
        </p>
      </div>

      <VerificationStatusCard
        status={verification.status}
        note={verification.note}
        submittedAt={verification.submitted_at}
        verifiedAt={verification.verified_at}
        action={
          canSubmit ? (
            <Button
              type="button"
              onClick={() => void submitForVerification()}
              disabled={submitting || dirty || !data.completeness.can_complete}
              className="h-11 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]"
            >
              {submitting ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
              {verification.status === "rejected" ? "Resubmit for verification" : "Submit for verification"}
            </Button>
          ) : null
        }
      />
      {canSubmit && (dirty || !data.completeness.can_complete) ? (
        <p className="-mt-3 text-xs text-[#94A3B8]">
          {dirty ? "Save your changes before submitting." : "Fill in the required fields (marked *) before submitting."}
        </p>
      ) : null}

      <section className="flex flex-col gap-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 sm:flex-row sm:items-center md:p-5">
        <UserAvatar src={data.profile.logo_url} name={data.profile.business_name || "Business"} size="xl" />
        <div className="min-w-0 flex-1">
          <p className="font-semibold">{data.profile.business_name || "Your business"}</p>
          <div className="mt-2 h-2 w-full max-w-sm overflow-hidden rounded-full bg-white/10" aria-hidden>
            <div className="h-full rounded-full bg-[color:var(--adfluence-action)]" style={{ width: `${data.completeness.percent}%` }} />
          </div>
          <p className="mt-1 text-xs text-[#CBD5E1]">Profile {data.completeness.percent}% complete</p>
          {missing.length > 0 ? <p className="mt-1 text-xs text-[#94A3B8]">Still missing: {missing.join(", ")}</p> : null}
        </div>
        <input
          ref={fileRef}
          id="business-logo"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          className="sr-only"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) void uploadLogo(file);
          }}
        />
        <Button
          type="button"
          variant="outline"
          disabled={uploading}
          onClick={() => fileRef.current?.click()}
          className="h-11 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
        >
          {uploading ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : <Upload className="mr-2 size-4" aria-hidden />}
          {data.profile.logo_url ? "Change logo" : "Upload logo"}
        </Button>
      </section>

      <form
        className="space-y-6"
        onSubmit={(e) => {
          e.preventDefault();
          void save();
        }}
      >
        <Section title="Business identity" description="Used to verify your business. Changing these after verification starts a new review.">
          <Field label="Business name *">
            <Input id="business_name" value={form.business_name} maxLength={255} onChange={(e) => set("business_name", e.target.value)} className={inputClass} placeholder="e.g. Selam Beauty" />
          </Field>
          <Field label="Legal business name" hint="As on your trade license. Recommended for faster verification.">
            <Input id="legal_business_name" value={form.legal_business_name} maxLength={255} onChange={(e) => set("legal_business_name", e.target.value)} className={inputClass} placeholder="e.g. Selam Beauty PLC" />
          </Field>
          <Field label="TIN number" hint="Taxpayer Identification Number. Recommended for faster verification.">
            <Input id="tin_number" value={form.tin_number} maxLength={50} inputMode="numeric" onChange={(e) => set("tin_number", e.target.value)} className={inputClass} placeholder="e.g. 0012345678" />
          </Field>
          <Field label="Business type">
            <select id="business_type" value={form.business_type} onChange={(e) => set("business_type", e.target.value)} className={selectClass}>
              <option value="">Select type</option>
              {BUSINESS_TYPES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </Field>
          <Field label="Business size">
            <select id="business_size" value={form.business_size} onChange={(e) => set("business_size", e.target.value)} className={selectClass}>
              <option value="">Select size</option>
              {BUSINESS_SIZES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </Field>
          <Field label="Owner full name *">
            <Input id="owner_full_name" value={form.owner_full_name} maxLength={255} onChange={(e) => set("owner_full_name", e.target.value)} className={inputClass} />
          </Field>
        </Section>

        <Section title="About your business">
          <Field label="Industry *">
            <Input id="industry" value={form.industry} maxLength={100} onChange={(e) => set("industry", e.target.value)} className={inputClass} placeholder="e.g. Beauty & cosmetics" />
          </Field>
          <Field label="Monthly marketing budget (ETB)">
            <Input id="monthly_marketing_budget" value={form.monthly_marketing_budget} inputMode="decimal" onChange={(e) => set("monthly_marketing_budget", e.target.value)} className={inputClass} placeholder="e.g. 50000" />
          </Field>
          <Field label="Description *" wide hint={`${form.description.trim().length}/5000 · at least 20 characters`}>
            <textarea
              id="description"
              value={form.description}
              maxLength={5000}
              rows={4}
              onChange={(e) => set("description", e.target.value)}
              placeholder="What do you sell or offer, and who are your customers?"
              className="w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white outline-none placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
            />
          </Field>
          <fieldset className="space-y-2 md:col-span-2">
            <legend className="text-sm text-[#CBD5E1]">Platforms you want to advertise on *</legend>
            <div className="flex flex-wrap gap-2">
              {PLATFORMS.map((p) => {
                const on = form.preferred_platforms.includes(p.id);
                return (
                  <button
                    key={p.id}
                    type="button"
                    aria-pressed={on}
                    onClick={() =>
                      set("preferred_platforms", on ? form.preferred_platforms.filter((x) => x !== p.id) : [...form.preferred_platforms, p.id])
                    }
                    className={cn(
                      "min-h-11 rounded-full border px-4 text-sm transition-colors md:min-h-9",
                      on
                        ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
                        : "border-white/10 bg-[color:var(--adfluence-surface)] text-[#CBD5E1] hover:text-white",
                    )}
                  >
                    {p.label}
                  </button>
                );
              })}
            </div>
          </fieldset>
        </Section>

        <Section title="Location">
          <Field label="City *">
            <Input id="city" value={form.city} maxLength={100} onChange={(e) => set("city", e.target.value)} className={inputClass} placeholder="e.g. Addis Ababa" />
          </Field>
          <Field label="Region">
            <Input id="region" value={form.region} maxLength={100} onChange={(e) => set("region", e.target.value)} className={inputClass} />
          </Field>
          <Field label="Address" wide>
            <Input id="address" value={form.address} maxLength={255} onChange={(e) => set("address", e.target.value)} className={inputClass} placeholder="e.g. Bole, Atlas area, building name" />
          </Field>
        </Section>

        <Section title="Online presence" description="Add your website or at least one social profile. Links must start with https://">
          <Field label="Website">
            <Input id="website_url" value={form.website_url} onChange={(e) => set("website_url", e.target.value)} className={inputClass} placeholder="https://" />
          </Field>
          <Field label="Instagram">
            <Input id="instagram" value={form.instagram} onChange={(e) => set("instagram", e.target.value)} className={inputClass} placeholder="https://instagram.com/…" />
          </Field>
          <Field label="TikTok">
            <Input id="tiktok" value={form.tiktok} onChange={(e) => set("tiktok", e.target.value)} className={inputClass} placeholder="https://tiktok.com/@…" />
          </Field>
          <Field label="Facebook">
            <Input id="facebook" value={form.facebook} onChange={(e) => set("facebook", e.target.value)} className={inputClass} placeholder="https://facebook.com/…" />
          </Field>
          <Field label="LinkedIn">
            <Input id="linkedin" value={form.linkedin} onChange={(e) => set("linkedin", e.target.value)} className={inputClass} placeholder="https://linkedin.com/company/…" />
          </Field>
        </Section>

        <div className="sticky bottom-0 z-10 -mx-4 flex flex-wrap items-center justify-end gap-3 border-t border-white/10 bg-[color:var(--adfluence-bg)]/95 px-4 py-3 backdrop-blur md:-mx-6 md:px-6">
          {confirmIdentityChange ? (
            <p role="alert" className="mr-auto max-w-xl text-xs text-amber-100">
              You changed your business name, legal name or TIN. Saving will send your business back to verification, and publishing,
              invites and contracts pause until it is approved again.
            </p>
          ) : dirty ? (
            <p className="mr-auto text-xs text-[#94A3B8]">You have unsaved changes.</p>
          ) : null}
          {confirmIdentityChange ? (
            <Button type="button" variant="ghost" onClick={() => setConfirmIdentityChange(false)} className="h-11 rounded-xl text-white hover:bg-white/10">
              Cancel
            </Button>
          ) : null}
          <Button
            type="submit"
            disabled={saving || !dirty}
            className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-6 text-white hover:bg-[color:var(--adfluence-action-hover)]"
          >
            {saving ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            {confirmIdentityChange ? "Save and re-verify" : "Save changes"}
          </Button>
        </div>
      </form>
    </div>
  );
}
