"use client";

import Image from "next/image";
import { useEffect, useRef, useState } from "react";
import { AlertCircle, Trash2, UploadCloud } from "lucide-react";

import { ConfirmActionDialog, type ConfirmActionRequest } from "@/components/common/confirm-action-dialog";
import { useInfluencerAccount } from "@/components/influencer/account-status";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  influencerApi,
  type InfluencerPortfolioItem,
  type InfluencerProfileResponse,
} from "@/lib/api/influencer";
import { CREATOR_STATUS, REQUIRED_FOR_SUBMIT } from "@/lib/influencer-profile";
import { parseAmount } from "@/lib/numbers";
import { cn, getInitials } from "@/lib/utils";

// Values match what the server stores (lowercase); labels are for display.
const platformOptions = [
  { value: "tiktok", label: "TikTok" },
  { value: "instagram", label: "Instagram" },
  { value: "youtube", label: "YouTube" },
] as const;
const platformLabel = (v: string) => platformOptions.find((o) => o.value === v.toLowerCase())?.label ?? v;

type NumberField = "followers" | "engagement" | "priceMin" | "priceMax";


const STATUS_TONE = {
  neutral: "border-white/15 bg-white/5 text-[#CBD5E1]",
  amber: "border-amber-400/30 bg-amber-500/15 text-amber-200",
  green: "border-emerald-400/30 bg-emerald-500/15 text-emerald-200",
  rose: "border-rose-400/30 bg-rose-500/15 text-rose-200",
} as const;

export function InfluencerProfile() {
  const [profileData, setProfileData] = useState<InfluencerProfileResponse | null>(null);
  const [availableNiches, setAvailableNiches] = useState<Array<{ id: string; name: string; slug: string }>>([]);
  const account = useInfluencerAccount();
  const [loadState, setLoadState] = useState<"loading" | "ready" | "error">("loading");
  // The form as last loaded from the server, to detect unsaved edits.
  const [savedSnapshot, setSavedSnapshot] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<NumberField, string>>>({});
  const [message, setMessage] = useState("");
  const [confirm, setConfirm] = useState<ConfirmActionRequest | null>(null);
  const [error, setError] = useState("");
  const [isSaving, setIsSaving] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [portfolioItems, setPortfolioItems] = useState<InfluencerPortfolioItem[]>([]);

  const avatarInputRef = useRef<HTMLInputElement | null>(null);
  const portfolioInputRef = useRef<HTMLInputElement | null>(null);

  // Starts empty and is filled from the saved profile; never pre-filled with
  // example values that could be saved to a real account.
  const [form, setForm] = useState({
    name: "",
    bio: "",
    niches: [] as string[],
    platforms: [] as string[],
    tiktok: "",
    instagram: "",
    youtube: "",
    followers: "",
    engagement: "",
    priceMin: "",
    priceMax: "",
    languages: "",
    location: "",
  });

  // Returns the loaded profile, or null when it could not be loaded.
  const loadProfile = async (): Promise<InfluencerProfileResponse | null> => {
    const [profRes, catRes] = await Promise.allSettled([
      influencerApi.getProfile(),
      influencerApi.listCategories(),
    ]);

    if (catRes.status === "fulfilled" && catRes.value.data?.items) {
      setAvailableNiches(catRes.value.data.items);
    }

    if (profRes.status !== "fulfilled" || !profRes.value.data) {
      setLoadState("error");
      return null;
    }
    const data = profRes.value.data;
    setProfileData(data);
    const p = data.profile;
    setPortfolioItems(p?.portfolio_items || []);
    if (p) {
      const loaded = {
        name: p.creator_name || p.full_name || "",
        bio: p.bio || "",
        niches: p.niches?.map((n) => n.name) || [],
        platforms: (p.platforms || []).map((x) => x.toLowerCase()),
        tiktok: p.social_links?.tiktok || "",
        instagram: p.social_links?.instagram || "",
        youtube: p.social_links?.youtube || "",
        followers: p.total_followers ? String(p.total_followers) : "",
        engagement: p.engagement_rate != null ? String(p.engagement_rate) : "",
        priceMin: p.min_price_etb != null ? String(p.min_price_etb) : "",
        priceMax: p.max_price_etb != null ? String(p.max_price_etb) : "",
        languages: p.languages?.join(", ") || "",
        location: p.city || "",
      };
      setForm(loaded);
      setSavedSnapshot(JSON.stringify(loaded));
    }
    setLoadState("ready");
    return data;
  };

  useEffect(() => {
    // Deferred like the other portal screens, so the effect itself sets no state.
    void Promise.resolve().then(loadProfile);
  }, []);

  const toggleMulti = (key: "niches" | "platforms", value: string) => {
    setForm((f) => {
      const arr = f[key];
      const next = arr.includes(value) ? arr.filter((x) => x !== value) : [...arr, value];
      return { ...f, [key]: next };
    });
  };

  // Checks the number fields; returns the parsed values, or null after
  // showing an error next to each field that can't be read.
  const readNumbers = () => {
    const errors: Partial<Record<NumberField, string>> = {};
    const followers = parseAmount(form.followers);
    const engagement = parseAmount(form.engagement);
    const priceMin = parseAmount(form.priceMin);
    const priceMax = parseAmount(form.priceMax);
    if (Number.isNaN(followers) || (followers !== undefined && !Number.isInteger(followers)))
      errors.followers = "Enter a whole number, e.g. 12000 or 12k.";
    if (Number.isNaN(engagement) || (engagement !== undefined && engagement > 100))
      errors.engagement = "Enter a percentage between 0 and 100, e.g. 4.5.";
    if (Number.isNaN(priceMin)) errors.priceMin = "Enter an amount in ETB, e.g. 5000 or 5k.";
    if (Number.isNaN(priceMax)) errors.priceMax = "Enter an amount in ETB, e.g. 20000 or 20k.";
    if (!errors.priceMin && !errors.priceMax && (priceMin === undefined) !== (priceMax === undefined))
      errors[priceMin === undefined ? "priceMin" : "priceMax"] = "Enter both a minimum and a maximum price.";
    if (!errors.priceMin && !errors.priceMax && priceMin !== undefined && priceMax !== undefined && priceMax < priceMin)
      errors.priceMax = "Maximum must be at least the minimum.";
    setFieldErrors(errors);
    if (Object.keys(errors).length > 0) return null;
    return { followers, engagement, priceMin, priceMax };
  };

  // Saves the form; returns true on success.
  const onSave = async (): Promise<boolean> => {
    setError("");
    setMessage("");
    if (loadState !== "ready") return false;
    const nums = readNumbers();
    if (!nums) {
      setError("Please fix the highlighted fields.");
      return false;
    }
    setIsSaving(true);
    try {
      await influencerApi.updateProfile({
        creator_name: form.name.trim(),
        bio: form.bio.trim(),
        niches: form.niches,
        platforms: form.platforms,
        social_links: {
          tiktok: form.tiktok,
          instagram: form.instagram,
          youtube: form.youtube,
        },
        followers: nums.followers,
        engagement_rate: nums.engagement,
        min_price_etb: nums.priceMin,
        max_price_etb: nums.priceMax,
        languages: form.languages.split(",").map((l) => l.trim()).filter(Boolean),
        city: form.location.trim(),
      });
      setMessage("Profile saved.");
      await loadProfile();
      return true;
    } catch (err: unknown) {
      const apiMsg = (err as { response?: { data?: { message?: string; error?: string } } })?.response?.data?.message;
      setError(apiMsg || (err instanceof Error ? err.message : "Failed to update profile."));
      return false;
    } finally {
      setIsSaving(false);
    }
  };

  const onAvatarSelected = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setIsUploading(true);
    setError("");
    try {
      const res = await influencerApi.uploadAvatar(file);
      if (res.data?.profile_image_url) {
        await loadProfile();
        setMessage("Avatar updated successfully.");
        setTimeout(() => setMessage(""), 4000);
      }
    } catch (err: unknown) {
      const apiMsg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message;
      setError(apiMsg || (err instanceof Error ? err.message : "Failed to upload avatar."));
    } finally {
      setIsUploading(false);
    }
  };

  const avatarUrl = profileData?.profile?.profile_image_url;

  const onDeleteAvatar = async () => {
    if (!avatarUrl) return;
    setIsUploading(true);
    setError("");
    try {
      await influencerApi.deleteAvatar();
      await loadProfile();
      setMessage("Profile avatar removed successfully.");
      setTimeout(() => setMessage(""), 4000);
    } catch (err: unknown) {
      const apiMsg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message;
      setError(apiMsg || (err instanceof Error ? err.message : "Failed to remove avatar."));
    } finally {
      setIsUploading(false);
    }
  };

  const onPortfolioSelected = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setIsUploading(true);
    setError("");
    try {
      const res = await influencerApi.uploadPortfolio(file);
      if (res.data?.portfolio_items) {
        setPortfolioItems(res.data.portfolio_items);
        setMessage("Portfolio item added.");
        setTimeout(() => setMessage(""), 4000);
      }
    } catch (err: unknown) {
      const apiMsg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message;
      setError(apiMsg || (err instanceof Error ? err.message : "Failed to upload portfolio item."));
    } finally {
      setIsUploading(false);
    }
  };

  const onDeletePortfolio = async (id: string) => {
    setError("");
    try {
      const res = await influencerApi.deletePortfolio(id);
      if (res.data?.portfolio_items) {
        setPortfolioItems(res.data.portfolio_items);
        setMessage("Portfolio item removed.");
        setTimeout(() => setMessage(""), 4000);
      }
    } catch (err: unknown) {
      const apiMsg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message;
      setError(apiMsg || (err instanceof Error ? err.message : "Failed to delete item."));
    }
  };

  // Saves unsaved edits first, so the reviewer sees what the creator sees.
  const onSubmitForReview = async () => {
    setError("");
    let completeness = profileData?.completeness;
    if (isDirty) {
      if (!(await onSave())) return;
      completeness = (await loadProfile())?.completeness;
    }
    if (!completeness?.can_submit) {
      setError("Your profile is saved, but some required details are still missing (see the list below).");
      return;
    }
    setIsSubmitting(true);
    try {
      const res = await influencerApi.submitProfile();
      if (res.data) {
        setMessage(res.data.message || "Profile submitted for admin verification!");
        await Promise.all([loadProfile(), account.refresh()]);
        setTimeout(() => setMessage(""), 5000);
      }
    } catch (err: unknown) {
      const apiMsg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message;
      setError(apiMsg || (err instanceof Error ? err.message : "Failed to submit profile for review."));
    } finally {
      setIsSubmitting(false);
    }
  };


  const completeness = profileData?.completeness?.percent ?? 0;
  const isRejected = profileData?.user?.status === "rejected";
  // The server refuses submission until every required field is filled.
  const canSubmit = profileData?.completeness?.can_submit ?? false;
  const userStatus = profileData?.user?.status;
  const awaitingSubmit = userStatus === "pending_profile" || isRejected;
  const isDirty = loadState === "ready" && savedSnapshot !== "" && JSON.stringify(form) !== savedSnapshot;
  const missingRequired = (profileData?.completeness?.missing ?? []).filter((f) => f in REQUIRED_FOR_SUBMIT);

  return (
    <div className="space-y-6 text-white">
      {loadState === "loading" ? (
        <p role="status" className="text-sm text-[#94A3B8]">Loading your profile…</p>
      ) : null}
      {loadState === "error" ? (
        <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-4 text-sm text-rose-100">
          Your profile could not be loaded, so editing is paused to avoid overwriting it.{" "}
          <button type="button" className="underline" onClick={() => { setLoadState("loading"); void loadProfile(); }}>
            Try again
          </button>
        </div>
      ) : null}
      <input
        ref={avatarInputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={onAvatarSelected}
      />
      <input
        ref={portfolioInputRef}
        type="file"
        accept="image/*,video/*"
        className="hidden"
        onChange={onPortfolioSelected}
      />

      <section className="flex flex-col gap-5 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:flex-row md:items-center">
        <div
          onClick={() => avatarInputRef.current?.click()}
          className="group relative mx-auto flex h-24 w-24 cursor-pointer shrink-0 items-center justify-center overflow-hidden rounded-full border-2 border-[color:var(--adfluence-orange)]/40 bg-gradient-to-br from-[color:var(--adfluence-orange)]/30 to-white/10 md:mx-0"
          title="Click to change profile avatar"
        >
          {avatarUrl ? (
            <Image
              src={avatarUrl}
              alt={form.name}
              fill
              unoptimized
              className="object-cover transition-opacity group-hover:opacity-75"
              sizes="96px"
            />
          ) : (
            <span className="text-2xl font-bold text-white uppercase tracking-wider">
              {getInitials(form.name, "CR")}
            </span>
          )}
          <div className="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 transition-opacity group-hover:opacity-100">
            <UploadCloud className="size-6 text-white" />
          </div>
        </div>
        <div className="flex-1 text-center md:text-left">
          <h1 className="text-2xl font-semibold text-white">{form.name}</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            {form.niches.join(", ") || "General"} · {form.platforms.map(platformLabel).join(", ")} · {form.location}
          </p>
          <div className="mt-2 flex flex-wrap justify-center gap-2 md:justify-start">
            {userStatus ? (
              <span className={cn("rounded-full border px-2 py-0.5 text-xs", STATUS_TONE[CREATOR_STATUS[userStatus]?.tone ?? "neutral"])}>
                {CREATOR_STATUS[userStatus]?.label ?? userStatus}
              </span>
            ) : null}
            <span className="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-xs text-[#CBD5E1]">
              Profile {completeness}% complete
            </span>
          </div>

          <div className="mt-3 flex flex-wrap items-center justify-center gap-2 md:justify-start">
            <Button
              type="button"
              size="sm"
              variant="outline"
              disabled={isUploading}
              onClick={() => avatarInputRef.current?.click()}
              className="h-8 rounded-lg border-white/15 bg-transparent text-xs text-white hover:bg-white/10"
            >
              <UploadCloud className="mr-1.5 size-3.5" />
              {avatarUrl ? "Change Photo" : "Upload Photo"}
            </Button>
            {avatarUrl ? (
              <Button
                type="button"
                size="sm"
                variant="outline"
                disabled={isUploading}
                onClick={() =>
                  setConfirm({
                    title: "Remove your profile photo?",
                    description: "Businesses will see your initials instead until you upload a new photo.",
                    confirmLabel: "Remove photo",
                    tone: "danger",
                    onConfirm: onDeleteAvatar,
                  })
                }
                className="h-11 rounded-lg border-white/15 bg-transparent text-xs text-rose-200 hover:bg-rose-500/10 sm:h-8"
              >
                <Trash2 className="mr-1.5 size-3.5" />
                Remove Photo
              </Button>
            ) : null}
          </div>
        </div>

        {awaitingSubmit && (canSubmit || isDirty) ? (
          <Button
            type="button"
            onClick={() => void onSubmitForReview()}
            disabled={isSubmitting || isSaving}
            className="h-11 rounded-xl bg-emerald-600 px-5 text-white hover:bg-emerald-500"
          >
            {isDirty ? "Save & " : ""}
            {isRejected ? "Re-apply for Verification" : "Submit for Verification"}
          </Button>
        ) : null}
      </section>

      {userStatus === "pending_approval" ? (
        <section role="status" className="rounded-2xl border border-amber-400/30 bg-amber-500/10 p-4 text-sm text-amber-50">
          <span className="font-semibold">Your profile is under review.</span> An admin is checking it. Changes you save now are visible to the reviewer.
        </section>
      ) : null}
      {userStatus === "active" ? (
        <p className="text-xs text-[#94A3B8]">You&apos;re approved. Changes you save are shown to businesses right away.</p>
      ) : null}

      {awaitingSubmit && !canSubmit && missingRequired.length > 0 ? (
        <section role="status" className="rounded-2xl border border-amber-400/30 bg-amber-500/10 p-5">
          <h2 className="text-base font-semibold text-amber-100">Complete these to submit for verification</h2>
          <p className="mt-1 text-sm text-amber-100/80">An admin can only review your profile after you submit it.</p>
          <ul className="mt-3 grid gap-1.5 text-sm text-amber-50 sm:grid-cols-2">
            {missingRequired.map((f) => (
              <li key={f} className="flex items-start gap-2">
                <AlertCircle className="mt-0.5 size-4 shrink-0 text-amber-300" aria-hidden />
                {REQUIRED_FOR_SUBMIT[f]}
              </li>
            ))}
          </ul>
        </section>
      ) : null}

      {isRejected ? (
        <section className="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-5">
          <div className="flex items-start gap-3">
            <AlertCircle className="size-6 text-rose-400 shrink-0 mt-0.5" />
            <div className="space-y-1">
              <h2 className="text-base font-semibold text-rose-200">Changes Requested by Admin</h2>
              <p className="text-sm text-rose-300">
                Reason: {profileData?.user?.status_reason || "Your profile was reviewed and requires updates before approval."}
              </p>
              <p className="text-xs text-[#CBD5E1] pt-1">
                You can make adjustments to your bio, niches, pricing, or media samples below and click <strong>&quot;Re-apply for Verification&quot;</strong> to submit your application again.
              </p>
            </div>
          </div>
        </section>
      ) : null}

      {error ? (
        <p className="rounded-xl border border-rose-400/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-200">
          {error}
        </p>
      ) : null}

      {message ? (
        <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-200">
          {message}
        </p>
      ) : null}

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
        <h2 className="text-lg font-semibold text-white">Profile details</h2>
        <div className="mt-4 grid gap-4 md:grid-cols-2">
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">Creator name</span>
            <Input
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">Bio</span>
            <textarea
              value={form.bio}
              onChange={(e) => setForm((f) => ({ ...f, bio: e.target.value }))}
              rows={4}
              maxLength={2000}
              className="w-full rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm text-white"
            />
            <span className={`block text-xs ${form.bio.trim().length >= 20 ? "text-[#94A3B8]" : "text-amber-200"}`}>
              {form.bio.trim().length} characters{form.bio.trim().length < 20 ? " (at least 20 needed to submit)" : ""}
            </span>
          </label>
          <div className="md:col-span-2">
            <p className="text-sm text-[#CBD5E1]">Niche (multi-select)</p>
            <div className="mt-2 flex flex-wrap gap-2">
              {availableNiches.map((n) => (
                <button
                  key={n.id}
                  type="button"
                  onClick={() => toggleMulti("niches", n.name)}
                  className={`rounded-xl border px-3 py-2 text-sm transition-colors ${
                    form.niches.includes(n.name) || form.niches.includes(n.slug)
                      ? "border-[color:var(--adfluence-orange)]/50 bg-[color:var(--adfluence-orange)]/15 text-white"
                      : "border-white/10 bg-[#06111F] text-[#CBD5E1] hover:border-white/20"
                  }`}
                >
                  {n.name}
                </button>
              ))}
            </div>
          </div>
          <div className="md:col-span-2">
            <p className="text-sm text-[#CBD5E1]">Platforms</p>
            <div className="mt-2 flex flex-wrap gap-2">
              {platformOptions.map((p) => (
                <button
                  key={p.value}
                  type="button"
                  aria-pressed={form.platforms.includes(p.value)}
                  onClick={() => toggleMulti("platforms", p.value)}
                  className={`rounded-xl border px-3 py-2 text-sm transition-colors ${
                    form.platforms.includes(p.value)
                      ? "border-[color:var(--adfluence-orange)]/50 bg-[color:var(--adfluence-orange)]/15 text-white"
                      : "border-white/10 bg-[#06111F] text-[#CBD5E1] hover:border-white/20"
                  }`}
                >
                  {p.label}
                </button>
              ))}
            </div>
          </div>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">TikTok link</span>
            <Input
              value={form.tiktok}
              placeholder="https://tiktok.com/@yourhandle"
              onChange={(e) => setForm((f) => ({ ...f, tiktok: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">Instagram link</span>
            <Input
              value={form.instagram}
              placeholder="https://instagram.com/yourhandle"
              onChange={(e) => setForm((f) => ({ ...f, instagram: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">YouTube link</span>
            <Input
              value={form.youtube}
              placeholder="https://youtube.com/@yourchannel"
              onChange={(e) => setForm((f) => ({ ...f, youtube: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2">
            <span className="text-sm text-[#CBD5E1]">Followers</span>
            <Input
              value={form.followers}
              inputMode="numeric"
              placeholder="e.g. 12000 or 12k"
              aria-invalid={fieldErrors.followers ? true : undefined}
              onChange={(e) => setForm((f) => ({ ...f, followers: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
            {fieldErrors.followers ? <span className="block text-xs text-rose-300">{fieldErrors.followers}</span> : null}
          </label>
          <label className="space-y-2">
            <span className="text-sm text-[#CBD5E1]">Engagement rate (%)</span>
            <Input
              value={form.engagement}
              inputMode="decimal"
              placeholder="e.g. 4.5"
              aria-invalid={fieldErrors.engagement ? true : undefined}
              onChange={(e) => setForm((f) => ({ ...f, engagement: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
            {fieldErrors.engagement ? <span className="block text-xs text-rose-300">{fieldErrors.engagement}</span> : null}
          </label>
          <label className="space-y-2">
            <span className="text-sm text-[#CBD5E1]">Price min (ETB)</span>
            <Input
              value={form.priceMin}
              inputMode="numeric"
              placeholder="e.g. 5000"
              aria-invalid={fieldErrors.priceMin ? true : undefined}
              onChange={(e) => setForm((f) => ({ ...f, priceMin: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
            {fieldErrors.priceMin ? <span className="block text-xs text-rose-300">{fieldErrors.priceMin}</span> : null}
          </label>
          <label className="space-y-2">
            <span className="text-sm text-[#CBD5E1]">Price max (ETB)</span>
            <Input
              value={form.priceMax}
              inputMode="numeric"
              placeholder="e.g. 20000"
              aria-invalid={fieldErrors.priceMax ? true : undefined}
              onChange={(e) => setForm((f) => ({ ...f, priceMax: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
            {fieldErrors.priceMax ? <span className="block text-xs text-rose-300">{fieldErrors.priceMax}</span> : null}
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">Languages (separate with commas)</span>
            <Input
              value={form.languages}
              onChange={(e) => setForm((f) => ({ ...f, languages: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">City</span>
            <Input
              value={form.location}
              onChange={(e) => setForm((f) => ({ ...f, location: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
        </div>
      </section>

      <section>
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-white">Portfolio & Content samples</h2>
          <Button
            type="button"
            onClick={() => portfolioInputRef.current?.click()}
            disabled={isUploading}
            className="h-9 rounded-xl bg-white/10 text-white hover:bg-white/15"
          >
            <UploadCloud className="mr-2 size-4" />
            {isUploading ? "Uploading..." : "Upload Media"}
          </Button>
        </div>

        {portfolioItems.filter((item) => Boolean(item.file_url)).length === 0 ? (
          <div
            onClick={() => portfolioInputRef.current?.click()}
            className="flex cursor-pointer flex-col items-center justify-center rounded-2xl border border-dashed border-white/20 bg-[rgba(255,255,255,0.05)] p-8 text-center transition-colors hover:bg-white/10"
          >
            <UploadCloud className="size-10 text-[#CBD5E1]/60" />
            <p className="mt-2 text-sm font-medium text-white">No portfolio media uploaded yet.</p>
            <p className="text-xs text-[#CBD5E1]">Click to upload photos or short video clips (MP4, WebM, PNG, JPG)</p>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {portfolioItems
              .filter((item) => Boolean(item.file_url))
              .map((item) => (
                <div
                  key={item.id}
                  className="group relative flex flex-col rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.05)] p-3"
                >
                  <div className="relative aspect-video overflow-hidden rounded-xl border border-white/10 bg-black/40">
                    {item.media_type === "video" ? (
                      <video
                        src={item.file_url}
                        controls
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <Image
                        src={item.file_url}
                        alt="Portfolio item"
                        fill
                        unoptimized
                        className="object-cover"
                        sizes="(max-width:768px) 100vw, 33vw"
                      />
                    )}
                  </div>
                  {/* Always visible: phones have no hover. */}
                  <button
                    type="button"
                    onClick={() =>
                      setConfirm({
                        title: "Delete this sample?",
                        description: "It is removed from your profile and businesses won't see it any more.",
                        confirmLabel: "Delete sample",
                        tone: "danger",
                        onConfirm: () => onDeletePortfolio(item.id),
                      })
                    }
                    className="absolute right-4 top-4 inline-flex size-11 items-center justify-center rounded-lg bg-black/60 text-white transition-colors hover:bg-red-600 sm:size-9"
                    aria-label="Delete this sample"
                  >
                    <Trash2 className="size-4" />
                  </button>
                </div>
              ))}
          </div>
        )}
      </section>


      <Button
        type="button"
        onClick={() => void onSave()}
        disabled={isSaving || loadState !== "ready"}
        className="h-12 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] md:w-auto md:px-10"
      >
        {isSaving ? "Saving..." : isDirty ? "Save changes" : "Save Profile"}
      </Button>
      <ConfirmActionDialog request={confirm} onClose={() => setConfirm(null)} />
    </div>
  );
}
