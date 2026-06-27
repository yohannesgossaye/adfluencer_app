"use client";

import Image from "next/image";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const CREATOR_IMAGE = "/images/influencers/influencer-2.jpg";

const nicheOptions = ["Beauty", "Fashion", "Food", "Tech", "Fitness", "Lifestyle", "Travel"] as const;
const platformOptions = ["TikTok", "Instagram", "YouTube"] as const;

const PROFILE_KEY = "adfluence_influencer_profile";

export function InfluencerProfile() {
  const [saved, setSaved] = useState(false);
  const [form, setForm] = useState({
    name: "Biruk Alemu",
    bio: "Tech creator focused on practical reviews, app walkthroughs, and honest recommendations for Ethiopian audiences.",
    niches: ["Tech"] as string[],
    platforms: ["YouTube"] as string[],
    tiktok: "https://tiktok.com/@birukcreates",
    instagram: "https://instagram.com/birukcreates",
    youtube: "https://youtube.com/@birukcreates",
    followers: "84000",
    engagement: "7.2",
    priceMin: "15000",
    priceMax: "35000",
    languages: "Amharic, English",
    location: "Addis Ababa",
  });

  useEffect(() => {
    try {
      const raw = localStorage.getItem(PROFILE_KEY);
      if (!raw) return;
      const parsed = JSON.parse(raw) as typeof form;
      if (parsed && typeof parsed === "object") setForm((f) => ({ ...f, ...parsed }));
    } catch {
      /* ignore */
    }
  }, []);

  const toggleMulti = (key: "niches" | "platforms", value: string) => {
    setForm((f) => {
      const arr = f[key];
      const next = arr.includes(value) ? arr.filter((x) => x !== value) : [...arr, value];
      return { ...f, [key]: next };
    });
  };

  const onSave = () => {
    if (typeof window !== "undefined") {
      localStorage.setItem(PROFILE_KEY, JSON.stringify(form));
    }
    setSaved(true);
    setTimeout(() => setSaved(false), 4000);
  };

  return (
    <div className="space-y-6 rounded-2xl bg-[#0B1A2B] p-4 text-white md:p-6">
      <section className="flex flex-col gap-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 md:flex-row md:items-center">
        <div className="relative mx-auto h-24 w-24 shrink-0 overflow-hidden rounded-full border-2 border-[#FF7A00]/40 md:mx-0">
          <Image src={CREATOR_IMAGE} alt="Biruk Alemu" fill className="object-cover" sizes="96px" />
        </div>
        <div className="flex-1 text-center md:text-left">
          <h1 className="text-2xl font-semibold text-white">Biruk Alemu</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">Tech · YouTube · {form.location}</p>
          <div className="mt-2 flex flex-wrap justify-center gap-2 md:justify-start">
            <span className="rounded-full border border-emerald-400/30 bg-emerald-500/15 px-2 py-0.5 text-xs text-emerald-200">
              Verified creator
            </span>
            <span className="rounded-full border border-white/15 bg-white/5 px-2 py-0.5 text-xs text-[#CBD5E1]">
              Profile {86}% complete
            </span>
          </div>
        </div>
      </section>

      {saved ? (
        <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-200">
          Profile saved locally.
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
                      ? "border-[#FF7A00]/50 bg-gradient-to-r from-[#FF7A00]/20 to-[#7C3AED]/20 text-white"
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
                      ? "border-[#FF7A00]/50 bg-gradient-to-r from-[#FF7A00]/20 to-[#7C3AED]/20 text-white"
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
              onChange={(e) => setForm((f) => ({ ...f, tiktok: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">Instagram link</span>
            <Input
              value={form.instagram}
              onChange={(e) => setForm((f) => ({ ...f, instagram: e.target.value }))}
              className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
            />
          </label>
          <label className="space-y-2 md:col-span-2">
            <span className="text-sm text-[#CBD5E1]">YouTube link</span>
            <Input
              value={form.youtube}
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
        <h2 className="mb-3 text-lg font-semibold text-white">Content samples</h2>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="flex flex-col rounded-2xl border border-dashed border-white/20 bg-[rgba(255,255,255,0.05)] p-4"
            >
              <div className="relative aspect-video overflow-hidden rounded-xl border border-white/10">
                <Image
                  src={CREATOR_IMAGE}
                  alt={`Sample ${i}`}
                  fill
                  className="object-cover opacity-80"
                  sizes="(max-width:768px) 100vw, 33vw"
                />
              </div>
              <p className="mt-2 text-center text-xs text-[#CBD5E1]">Upload coming soon</p>
            </div>
          ))}
        </div>
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
                  ? "rounded-2xl border border-[#FF7A00]/40 bg-[linear-gradient(135deg,rgba(255,122,0,0.12),rgba(124,58,237,0.12))] p-4"
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
        className="h-12 w-full rounded-xl bg-gradient-to-r from-[#FF7A00] to-[#7C3AED] text-white hover:opacity-95 md:w-auto md:px-10"
      >
        Save Profile
      </Button>
    </div>
  );
}
