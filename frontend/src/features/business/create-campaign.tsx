"use client";

import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

type CampaignFormState = {
  title: string;
  businessName: string;
  goal: string;
  description: string;
  niche: string;
  targetAudience: string;
  location: string;
  platforms: string[];
  budgetRange: string;
  creatorSize: string;
  creatorsNeeded: string;
  deliverables: string[];
  startDate: string;
  endDate: string;
  additionalNotes: string;
};

const goals = [
  "Brand awareness",
  "Product launch",
  "App downloads",
  "Sales conversion",
  "Event promotion",
] as const;

const niches = ["Beauty", "Fashion", "Food", "Tech", "Fitness", "Lifestyle", "Travel"] as const;

const locations = ["Addis Ababa", "Ethiopia wide", "Regional campaign"] as const;

const platforms = ["TikTok", "Instagram", "YouTube"] as const;

const budgetRanges = [
  "Under 10,000 ETB",
  "10,000–25,000 ETB",
  "25,000–50,000 ETB",
  "50,000+ ETB",
] as const;

const creatorSizes = ["Micro creator", "Mid-tier creator", "Large creator"] as const;

const deliverables = [
  "TikTok video",
  "Instagram reel",
  "Instagram story",
  "YouTube short",
  "Product review",
] as const;

const initialForm: CampaignFormState = {
  title: "",
  businessName: "",
  goal: "",
  description: "",
  niche: "",
  targetAudience: "",
  location: "",
  platforms: [],
  budgetRange: "",
  creatorSize: "",
  creatorsNeeded: "",
  deliverables: [],
  startDate: "",
  endDate: "",
  additionalNotes: "",
};

function toggleArrayValue(arr: string[], value: string) {
  return arr.includes(value) ? arr.filter((entry) => entry !== value) : [...arr, value];
}

export function CreateCampaign() {
  const router = useRouter();
  const [form, setForm] = useState<CampaignFormState>(initialForm);
  const [error, setError] = useState<string>("");

  const estimatedMatches = useMemo(() => {
    let score = 8;
    if (form.goal) score += 4;
    if (form.niche) score += 4;
    score += form.platforms.length * 6;
    if (form.budgetRange) score += 5;
    if (form.creatorSize) score += 4;
    if (Number(form.creatorsNeeded) >= 5) score += 3;
    return score;
  }, [form]);

  const onSubmit = () => {
    const requiredMissing =
      !form.title ||
      !form.businessName ||
      !form.goal ||
      !form.description ||
      !form.niche ||
      !form.location ||
      form.platforms.length === 0 ||
      !form.budgetRange ||
      !form.creatorSize ||
      !form.startDate ||
      !form.endDate;

    if (requiredMissing) {
      setError("Please complete the required campaign details.");
      return;
    }

    if (form.endDate < form.startDate) {
      setError("Campaign end date must be after the start date.");
      return;
    }

    const draft = {
      ...form,
      estimatedMatches,
      createdAt: new Date().toISOString(),
    };

    localStorage.setItem("adfluence_campaign_draft", JSON.stringify(draft));
    setError("");
    router.push("/business/recommendations");
  };

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-white">Create Campaign</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            Tell us what you want to promote and we&apos;ll recommend the best creators.
          </p>
        </div>
        <span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs font-medium text-[#CBD5E1]">
          Step 1 of 4
        </span>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
        <p className="mb-3 text-xs font-medium uppercase tracking-[0.16em] text-[#CBD5E1]">
          Campaign setup progress
        </p>
        <div className="grid grid-cols-4 gap-2">
          {[1, 2, 3, 4].map((step) => (
            <div
              key={step}
              className={`h-2 rounded-full ${
                step === 1
                  ? "bg-[color:var(--adfluence-orange)]"
                  : "bg-white/10"
              }`}
            />
          ))}
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_340px]">
        <div className="space-y-6">
          <section className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
            <h2 className="text-lg font-semibold">Campaign Basics</h2>
            <div className="grid gap-4 md:grid-cols-2">
              <label className="space-y-2 md:col-span-2">
                <span className="text-sm text-[#CBD5E1]">Campaign title *</span>
                <Input
                  value={form.title}
                  onChange={(e) => setForm((prev) => ({ ...prev, title: e.target.value }))}
                  placeholder="e.g. Summer Product Launch"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Brand / business name *</span>
                <Input
                  value={form.businessName}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, businessName: e.target.value }))
                  }
                  placeholder="e.g. Betu Cosmetics"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Campaign goal *</span>
                <select
                  value={form.goal}
                  onChange={(e) => setForm((prev) => ({ ...prev, goal: e.target.value }))}
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                >
                  <option value="" disabled>
                    Select campaign goal
                  </option>
                  {goals.map((goal) => (
                    <option key={goal} value={goal} className="bg-[color:var(--adfluence-surface)]">
                      {goal}
                    </option>
                  ))}
                </select>
              </label>
              <label className="space-y-2 md:col-span-2">
                <span className="text-sm text-[#CBD5E1]">Campaign description *</span>
                <textarea
                  value={form.description}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, description: e.target.value }))
                  }
                  rows={4}
                  placeholder="Describe your campaign, key message, product details, and outcomes you want."
                  className="w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white outline-none placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
            </div>
          </section>

          <section className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
            <h2 className="text-lg font-semibold">Targeting</h2>
            <div className="grid gap-4 md:grid-cols-2">
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Industry / niche *</span>
                <select
                  value={form.niche}
                  onChange={(e) => setForm((prev) => ({ ...prev, niche: e.target.value }))}
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                >
                  <option value="" disabled>
                    Select industry
                  </option>
                  {niches.map((niche) => (
                    <option key={niche} value={niche} className="bg-[color:var(--adfluence-surface)]">
                      {niche}
                    </option>
                  ))}
                </select>
              </label>
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Target audience</span>
                <Input
                  value={form.targetAudience}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, targetAudience: e.target.value }))
                  }
                  placeholder="e.g. Women 18-30 in urban areas"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Location *</span>
                <select
                  value={form.location}
                  onChange={(e) => setForm((prev) => ({ ...prev, location: e.target.value }))}
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                >
                  <option value="" disabled>
                    Select location
                  </option>
                  {locations.map((location) => (
                    <option key={location} value={location} className="bg-[color:var(--adfluence-surface)]">
                      {location}
                    </option>
                  ))}
                </select>
              </label>
              <div className="space-y-2">
                <p className="text-sm text-[#CBD5E1]">Platforms *</p>
                <div className="flex flex-wrap gap-2">
                  {platforms.map((platform) => (
                    <label
                      key={platform}
                      className="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-[#E2E8F0]"
                    >
                      <input
                        type="checkbox"
                        checked={form.platforms.includes(platform)}
                        onChange={() =>
                          setForm((prev) => ({
                            ...prev,
                            platforms: toggleArrayValue(prev.platforms, platform),
                          }))
                        }
                        className="h-4 w-4 rounded border-white/20 bg-[color:var(--adfluence-surface)] accent-[color:var(--adfluence-orange)]"
                      />
                      {platform}
                    </label>
                  ))}
                </div>
              </div>
            </div>
          </section>

          <section className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
            <h2 className="text-lg font-semibold">Budget &amp; Creator Type</h2>
            <div className="grid gap-4 md:grid-cols-2">
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Budget range *</span>
                <select
                  value={form.budgetRange}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, budgetRange: e.target.value }))
                  }
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                >
                  <option value="" disabled>
                    Select budget
                  </option>
                  {budgetRanges.map((range) => (
                    <option key={range} value={range} className="bg-[color:var(--adfluence-surface)]">
                      {range}
                    </option>
                  ))}
                </select>
              </label>
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Preferred creator size *</span>
                <select
                  value={form.creatorSize}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, creatorSize: e.target.value }))
                  }
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                >
                  <option value="" disabled>
                    Select creator tier
                  </option>
                  {creatorSizes.map((size) => (
                    <option key={size} value={size} className="bg-[color:var(--adfluence-surface)]">
                      {size}
                    </option>
                  ))}
                </select>
              </label>
              <label className="space-y-2 md:col-span-2">
                <span className="text-sm text-[#CBD5E1]">Number of creators needed</span>
                <Input
                  type="number"
                  min={1}
                  value={form.creatorsNeeded}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, creatorsNeeded: e.target.value }))
                  }
                  placeholder="e.g. 5"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
            </div>
          </section>

          <section className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
            <h2 className="text-lg font-semibold">Deliverables &amp; Timeline</h2>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2 md:col-span-2">
                <p className="text-sm text-[#CBD5E1]">Deliverables</p>
                <div className="flex flex-wrap gap-2">
                  {deliverables.map((deliverable) => (
                    <label
                      key={deliverable}
                      className="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-[#E2E8F0]"
                    >
                      <input
                        type="checkbox"
                        checked={form.deliverables.includes(deliverable)}
                        onChange={() =>
                          setForm((prev) => ({
                            ...prev,
                            deliverables: toggleArrayValue(prev.deliverables, deliverable),
                          }))
                        }
                        className="h-4 w-4 rounded border-white/20 bg-[color:var(--adfluence-surface)] accent-[color:var(--adfluence-orange)]"
                      />
                      {deliverable}
                    </label>
                  ))}
                </div>
              </div>
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Campaign start date *</span>
                <Input
                  type="date"
                  value={form.startDate}
                  onChange={(e) => setForm((prev) => ({ ...prev, startDate: e.target.value }))}
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm text-[#CBD5E1]">Campaign end date *</span>
                <Input
                  type="date"
                  value={form.endDate}
                  onChange={(e) => setForm((prev) => ({ ...prev, endDate: e.target.value }))}
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
              <label className="space-y-2 md:col-span-2">
                <span className="text-sm text-[#CBD5E1]">Additional notes</span>
                <textarea
                  value={form.additionalNotes}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, additionalNotes: e.target.value }))
                  }
                  rows={4}
                  placeholder="Anything else our team should know before recommending creators?"
                  className="w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white outline-none placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
            </div>
          </section>

          {error ? (
            <p className="rounded-xl border border-rose-400/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-200">
              {error}
            </p>
          ) : null}
        </div>

        <aside className="xl:sticky xl:top-6 xl:h-fit">
          <div className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
            <h3 className="text-lg font-semibold text-white">Campaign Summary</h3>
            <div className="space-y-3 text-sm">
              <p className="text-[#CBD5E1]">
                Goal: <span className="text-white">{form.goal || "Not selected"}</span>
              </p>
              <p className="text-[#CBD5E1]">
                Niche: <span className="text-white">{form.niche || "Not selected"}</span>
              </p>
              <p className="text-[#CBD5E1]">
                Platforms:{" "}
                <span className="text-white">
                  {form.platforms.length ? form.platforms.join(", ") : "Not selected"}
                </span>
              </p>
              <p className="text-[#CBD5E1]">
                Budget range:{" "}
                <span className="text-white">{form.budgetRange || "Not selected"}</span>
              </p>
              <p className="text-[#CBD5E1]">
                Creator size:{" "}
                <span className="text-white">{form.creatorSize || "Not selected"}</span>
              </p>
              <p className="text-[#CBD5E1]">
                Estimated matches: <span className="text-white">{estimatedMatches}</span>
              </p>
            </div>
            <p className="text-xs text-[#CBD5E1]">
              Recommendations are based on niche, platform, budget, and creator fit.
            </p>
            <Button
              onClick={onSubmit}
              className="h-11 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
            >
              Find Influencers
            </Button>
          </div>
        </aside>
      </div>
    </div>
  );
}
