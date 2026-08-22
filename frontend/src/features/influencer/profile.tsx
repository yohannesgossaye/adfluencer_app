"use client";

import Image from "next/image";
import { useEffect, useRef, useState } from "react";
import { Trash2, UploadCloud } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  influencerApi,
  type InfluencerPortfolioItem,
  type InfluencerProfileResponse,
} from "@/lib/api/influencer";

const CREATOR_IMAGE = "/images/influencers/influencer-2.jpg";

const nicheOptions = ["Beauty", "Fashion", "Food", "Tech", "Fitness", "Lifestyle", "Travel"] as const;
const platformOptions = ["TikTok", "Instagram", "YouTube"] as const;

export function InfluencerProfile() {
  const [profileData, setProfileData] = useState<InfluencerProfileResponse | null>(null);
  const [saved, setSaved] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [isSaving, setIsSaving] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [portfolioItems, setPortfolioItems] = useState<InfluencerPortfolioItem[]>([]);

  const avatarInputRef = useRef<HTMLInputElement | null>(null);
  const portfolioInputRef = useRef<HTMLInputElement | null>(null);

  const [form, setForm] = useState({
    name: "Biruk Alemu",
    bio: "Tech creator focused on practical reviews, app walkthroughs, and honest recommendations for Ethiopian audiences.",
    niches: ["Tech"] as string[],
    platforms: ["YouTube"] as string[],
    tiktok: "",
    instagram: "",
    youtube: "",
    followers: "84000",
    engagement: "7.2",
    priceMin: "15000",
    priceMax: "35000",
    languages: "Amharic, English",
    location: "Addis Ababa",
  });

  const loadProfile = async () => {
    try {
      const res = await influencerApi.getProfile();
      if (res.data) {
        setProfileData(res.data);
        const p = res.data.profile;
        setPortfolioItems(p.portfolio_items || []);
        setForm({
          name: p.creator_name || p.full_name || "Creator",
          bio: p.bio || "",
          niches: p.niches?.map((n) => n.name) || ["Tech"],
          platforms: p.platforms || ["YouTube"],
          tiktok: p.social_links?.tiktok || "",
          instagram: p.social_links?.instagram || "",
          youtube: p.social_links?.youtube || "",
          followers: String(p.total_followers || 84000),
          engagement: String(p.engagement_rate || 7.2),
          priceMin: String(p.min_price_etb || 15000),
          priceMax: String(p.max_price_etb || 35000),
          languages: p.languages?.join(", ") || "Amharic, English",
          location: p.city || "Addis Ababa",
        });
      }
    } catch {
      // fallback
    }
  };

  useEffect(() => {
    loadProfile();
  }, []);

  const toggleMulti = (key: "niches" | "platforms", value: string) => {
    setForm((f) => {
      const arr = f[key];
      const next = arr.includes(value) ? arr.filter((x) => x !== value) : [...arr, value];
      return { ...f, [key]: next };
    });
  };

  const onSave = async () => {
    setError("");
    setMessage("");
    setIsSaving(true);
    try {
      await influencerApi.updateProfile({
        creator_name: form.name.trim(),
        bio: form.bio.trim(),
        niches: form.niches,
        platforms: form.platforms.map((p) => p.toLowerCase()),
        social_links: {
          tiktok: form.tiktok,
          instagram: form.instagram,
          youtube: form.youtube,
        },
        followers: Number(form.followers) || undefined,
        engagement_rate: Number(form.engagement) || undefined,
        min_price_etb: Number(form.priceMin) || undefined,
        max_price_etb: Number(form.priceMax) || undefined,
        languages: form.languages.split(",").map((l) => l.trim()).filter(Boolean),
        city: form.location.trim(),
      });
      setSaved(true);
      setMessage("Profile saved successfully.");
      await loadProfile();
      setTimeout(() => setSaved(false), 4000);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to update profile.");
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
      setError(err instanceof Error ? err.message : "Failed to upload avatar.");
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
      setError(err instanceof Error ? err.message : "Failed to upload portfolio item.");
    } finally {
      setIsUploading(false);
    }
  };

  const onDeletePortfolio = async (id: string) => {
    try {
      const res = await influencerApi.deletePortfolio(id);
      if (res.data?.portfolio_items) {
        setPortfolioItems(res.data.portfolio_items);
      }
    } catch {
      // fallback
    }
  };

  const onSubmitForReview = async () => {
    setError("");
    try {
      const res = await influencerApi.submitProfile();
      setMessage(res.data?.message || "Profile submitted for review!");
      await loadProfile();
      setTimeout(() => setMessage(""), 4000);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to submit profile for review.");
    }
  };

  const completenessPercent = profileData?.completeness?.percent ?? 86;
  const avatarUrl = profileData?.profile?.profile_image_url || CREATOR_IMAGE;
  const isVerified = profileData?.profile?.is_verified ?? false;
  const canSubmit = profileData?.completeness?.can_submit ?? false;

  return (
    <div className="space-y-6 text-white">
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

      <section className="flex flex-col gap-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:flex-row md:items-center">
        <div
          onClick={() => avatarInputRef.current?.click()}
          className="group relative mx-auto h-24 w-24 cursor-pointer shrink-0 overflow-hidden rounded-full border-2 border-[color:var(--adfluence-orange)]/40 md:mx-0"
          title="Click to upload profile avatar"
        >
          <Image
            src={avatarUrl}
            alt={form.name}
            fill
            className="object-cover transition-opacity group-hover:opacity-75"
            sizes="96px"
          />
          <div className="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition-opacity group-hover:opacity-100">
            <UploadCloud className="size-6 text-white" />
          </div>
        </div>
        <div className="flex-1 text-center md:text-left">
          <h1 className="text-2xl font-semibold text-white">{form.name}</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            {form.niches.join(", ") || "General"} · {form.platforms.join(", ")} · {form.location}
          </p>
          <div className="mt-2 flex flex-wrap justify-center gap-2 md:justify-start">
            {isVerified ? (
              <span className="rounded-full border border-emerald-400/30 bg-emerald-500/15 px-2 py-0.5 text-xs text-emerald-200">
                Verified creator
              </span>
            ) : (
              <span className="rounded-full border border-amber-400/30 bg-amber-500/15 px-2 py-0.5 text-xs text-amber-200">
                Pending Verification
              </span>
            )}
            <span className="rounded-full border border-white/15 bg-white/5 px-2 py-0.5 text-xs text-[#CBD5E1]">
              Profile {completenessPercent}% complete
            </span>
          </div>
        </div>

        {canSubmit ? (
          <Button
            type="button"
            onClick={onSubmitForReview}
            className="h-11 rounded-xl bg-emerald-600 px-5 text-white hover:bg-emerald-500"
          >
            Submit for Verification
          </Button>
        ) : null}
      </section>

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
              className="w-full rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm text-white"
            />
          </label>
          <div className="md:col-span-2">
            <p className="text-sm text-[#CBD5E1]">Niche (multi-select)</p>
            <div className="mt-2 flex flex-wrap gap-2">
              {nicheOptions.map((n) => (
                <button
                  key={n}
                  type="button"
                  onClick={() => toggleMulti("niches", n)}
                  className={`rounded-xl border px-3 py-2 text-sm transition-colors ${
                    form.niches.includes(n)
                      ? "border-[color:var(--adfluence-orange)]/50 bg-[color:var(--adfluence-orange)]/15 text-white"
                      : "border-white/10 bg-[#06111F] text-[#CBD5E1] hover:border-white/20"
                  }`}
                >
                  {n}
                </button>
              ))}
            </div>
          </div>
          <div className="md:col-span-2">
            <p className="text-sm text-[#CBD5E1]">Platforms</p>
            <div className="mt-2 flex flex-wrap gap-2">
              {platformOptions.map((p) => (
                <button
                  key={p}
                  type="button"
                  onClick={() => toggleMulti("platforms", p)}
                  className={`rounded-xl border px-3 py-2 text-sm transition-colors ${
                    form.platforms.includes(p)
                      ? "border-[color:var(--adfluence-orange)]/50 bg-[color:var(--adfluence-orange)]/15 text-white"
                      : "border-white/10 bg-[#06111F] text-[#CBD5E1] hover:border-white/20"
                  }`}
                >
                  {p}
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
              onChange={(e) => setForm((f) => ({ ...f, followers: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2">
            <span className="text-sm text-[#CBD5E1]">Engagement rate (%)</span>
            <Input
              value={form.engagement}
              onChange={(e) => setForm((f) => ({ ...f, engagement: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2">
            <span className="text-sm text-[#CBD5E1]">Price min (ETB)</span>
            <Input
              value={form.priceMin}
              onChange={(e) => setForm((f) => ({ ...f, priceMin: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2">
            <span className="text-sm text-[#CBD5E1]">Price max (ETB)</span>
            <Input
              value={form.priceMax}
              onChange={(e) => setForm((f) => ({ ...f, priceMax: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">Languages</span>
            <Input
              value={form.languages}
              onChange={(e) => setForm((f) => ({ ...f, languages: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">Location</span>
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

        {portfolioItems.length === 0 ? (
          <div
            onClick={() => portfolioInputRef.current?.click()}
            className="flex cursor-pointer flex-col items-center justify-center rounded-2xl border border-dashed border-white/20 bg-[rgba(255,255,255,0.05)] p-8 text-center"
          >
            <UploadCloud className="size-10 text-[#CBD5E1]/60" />
            <p className="mt-2 text-sm text-white">No portfolio media uploaded yet.</p>
            <p className="text-xs text-[#CBD5E1]">Click to upload photos or short video clips</p>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {portfolioItems.map((item) => (
              <div
                key={item.id}
                className="group relative flex flex-col rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.05)] p-3"
              >
                <div className="relative aspect-video overflow-hidden rounded-xl border border-white/10">
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
                      className="object-cover"
                      sizes="(max-width:768px) 100vw, 33vw"
                    />
                  )}
                </div>
                <button
                  type="button"
                  onClick={() => onDeletePortfolio(item.id)}
                  className="absolute right-4 top-4 rounded-lg bg-red-600/80 p-1.5 text-white opacity-0 transition-opacity hover:bg-red-600 group-hover:opacity-100"
                  title="Delete media"
                >
                  <Trash2 className="size-4" />
                </button>
              </div>
            ))}
          </div>
        )}
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-white">Pricing packages</h2>
        <div className="grid gap-4 md:grid-cols-3">
          {(
            [
              { title: "Basic", items: ["1 short video", "1 story", "Posting window: 7 days"] },
              { title: "Standard", items: ["2 videos", "3 stories", "30-day usage rights"] },
              { title: "Premium", items: ["Campaign bundle", "Multi-platform", "Priority support"] },
            ] as const
          ).map((pkg, index) => (
            <article
              key={pkg.title}
              className={
                index === 1
                  ? "rounded-2xl border border-[color:var(--adfluence-orange)]/40 bg-[color:var(--adfluence-orange)]/10 p-4"
                  : "rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
              }
            >
              <div className="flex items-center justify-between">
                <h3 className="font-semibold text-white">{pkg.title}</h3>
                {index === 1 ? (
                  <span className="text-[10px] uppercase text-orange-200">Recommended</span>
                ) : null}
              </div>
              <ul className="mt-2 space-y-1 text-sm text-[#CBD5E1]">
                {pkg.items.map((line) => (
                  <li key={line}>• {line}</li>
                ))}
              </ul>
            </article>
          ))}
        </div>
      </section>

      <Button
        type="button"
        onClick={onSave}
        disabled={isSaving}
        className="h-12 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] md:w-auto md:px-10"
      >
        {isSaving ? "Saving..." : "Save Profile"}
      </Button>
    </div>
  );
}
