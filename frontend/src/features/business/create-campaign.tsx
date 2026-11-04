"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { businessApi, type CampaignDetail } from "@/lib/api/business";

type CampaignFormState = {
  title: string;
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

const locations = ["Addis Ababa", "Ethiopia wide", "Regional campaign"] as const;

const platforms = ["TikTok", "Instagram", "YouTube"] as const;

const budgetRanges = [
  "Under 10,000 ETB",
  "10,000–25,000 ETB",
  "25,000–50,000 ETB",
  "50,000+ ETB",
] as const;

const creatorSizes = [
  "Any size",
  "Nano creator",
  "Micro creator",
  "Mid-tier creator",
  "Large creator",
  "Very large creator",
] as const;

const deliverables = [
  "TikTok video",
  "Instagram reel",
  "Instagram story",
  "YouTube short",
  "Product review",
] as const;

const initialForm: CampaignFormState = {
  title: "",
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

const BUDGETS: Record<(typeof budgetRanges)[number], [number, number]> = {
  "Under 10,000 ETB": [0, 10000],
  "10,000–25,000 ETB": [10000, 25000],
  "25,000–50,000 ETB": [25000, 50000],
  "50,000+ ETB": [50000, 200000],
};

const SIZE_TO_API: Record<(typeof creatorSizes)[number], string> = {
  "Any size": "any",
  "Nano creator": "nano",
  "Micro creator": "micro",
  "Mid-tier creator": "mid",
  "Large creator": "macro",
  "Very large creator": "mega",
};

// A saved budget that matches none of the presets stays selectable as-is.
const CUSTOM_BUDGET = "custom";

function etbRange(min: number, max: number) {
  return `ETB ${min.toLocaleString("en-US")} – ${max.toLocaleString("en-US")}`;
}

const DELIVERABLE_TO_API: Record<(typeof deliverables)[number], string> = {
  "TikTok video": "tiktok_video",
  "Instagram reel": "instagram_reel",
  "Instagram story": "instagram_story",
  "YouTube short": "youtube_short",
  "Product review": "product_review",
};

/** Maps a saved campaign back into the form (edit mode). */
function campaignToForm(c: CampaignDetail): CampaignFormState {
  const preset = (Object.entries(BUDGETS) as Array<[(typeof budgetRanges)[number], [number, number]]>).find(
    ([, [min, max]]) => min === (c.budget_min_etb ?? -1) && max === (c.budget_max_etb ?? -1),
  )?.[0];
  const budget = preset ?? (c.budget_max_etb ? CUSTOM_BUDGET : "");
  const size = (Object.entries(SIZE_TO_API).find(([, api]) => api === c.creator_size_preference)?.[0] ?? "") as string;
  const chosen = (Object.entries(DELIVERABLE_TO_API) as Array<[string, string]>)
    .filter(([, api]) => c.deliverables.some((d) => d.type === api))
    .map(([label]) => label);
  // The server stores notes inside the description (MergeDescription); split them back out.
  const [description, notes = ""] = (c.description ?? "").split("\n\nAdditional notes:\n");
  return {
    title: c.title,
    goal: c.goal ?? "",
    description,
    niche: c.category?.name ?? "",
    targetAudience: c.target_audience ?? "",
    location: c.city ?? "",
    platforms: c.platforms.map((p) => platforms.find((x) => x.toLowerCase() === p) ?? p),
    budgetRange: budget,
    creatorSize: size,
    creatorsNeeded: String(c.number_of_creators_needed ?? ""),
    deliverables: chosen,
    startDate: c.start_date?.slice(0, 10) ?? "",
    endDate: c.end_date?.slice(0, 10) ?? "",
    additionalNotes: notes,
  };
}

type FieldErrors = Partial<Record<keyof CampaignFormState, string>>;

function toggleArrayValue(arr: string[], value: string) {
  return arr.includes(value) ? arr.filter((entry) => entry !== value) : [...arr, value];
}

/** Create a campaign draft, or edit an existing campaign when campaignId is given. */
export function CreateCampaign({ campaignId }: { campaignId?: string } = {}) {
  const router = useRouter();
  const isEdit = Boolean(campaignId);
  const [form, setForm] = useState<CampaignFormState>(initialForm);
  // What the campaign had when opened for editing, so saving never silently changes it.
  const [saved, setSaved] = useState<{ budget: [number, number] | null; quantities: Record<string, number> }>({
    budget: null,
    quantities: {},
  });
  const [loadingCampaign, setLoadingCampaign] = useState(isEdit);
  const [categories, setCategories] = useState<Array<{ id: string; name: string; slug: string }>>([]);
  const [error, setError] = useState<string>("");
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const fieldError = (name: keyof CampaignFormState) =>
    fieldErrors[name] ? (
      <span role="alert" className="block text-xs font-medium text-rose-300">
        {fieldErrors[name]}
      </span>
    ) : null;
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    async function loadCategories() {
      try {
        const res = await businessApi.listCategories();
        if (res.data?.items) {
          setCategories(res.data.items);
        }
      } catch {
        setError("Niches could not be loaded. Refresh the page to try again.");
      }
    }
    loadCategories();
  }, []);

  useEffect(() => {
    if (!campaignId) return;
    businessApi
      .getCampaign(campaignId)
      .then((res) => {
        if (res.data) {
          const c = res.data;
          setForm(campaignToForm(c));
          setSaved({
            budget: c.budget_max_etb ? [c.budget_min_etb ?? 0, c.budget_max_etb] : null,
            quantities: Object.fromEntries(c.deliverables.map((d) => [d.type, d.quantity ?? 1])),
          });
        }
      })
      .catch(() => setError("Could not load this campaign."))
      .finally(() => setLoadingCampaign(false));
  }, [campaignId]);


  const onSubmit = async () => {
    // Check everything at once and show each problem next to its field.
    const problems: FieldErrors = {};
    if (form.title.trim().length < 3) problems.title = "Enter a title of at least 3 characters.";
    if (!form.goal) problems.goal = "Choose a goal.";
    if (form.description.trim().length < 20) problems.description = "Describe the campaign in at least 20 characters.";
    if (!form.niche) problems.niche = "Choose an industry or niche.";
    if (!form.location) problems.location = "Choose where the campaign runs.";
    if (form.platforms.length === 0) problems.platforms = "Choose at least one platform.";
    if (!form.budgetRange) problems.budgetRange = "Choose a budget range.";
    if (!form.creatorSize) problems.creatorSize = "Choose a creator size.";
    if (!form.startDate) problems.startDate = "Choose a start date.";
    if (!form.endDate) problems.endDate = "Choose an end date.";
    else if (form.startDate && form.endDate < form.startDate) problems.endDate = "The end date must be on or after the start date.";
    const creators = Number(form.creatorsNeeded || "1");
    if (!Number.isInteger(creators) || creators < 1 || creators > 50) problems.creatorsNeeded = "Enter a number from 1 to 50.";
    if (form.deliverables.length === 0) problems.deliverables = "Choose at least one deliverable, e.g. TikTok video.";
    setFieldErrors(problems);
    const first = Object.keys(problems)[0];
    if (first) {
      setError(`Please fix ${Object.keys(problems).length === 1 ? "the highlighted field" : `the ${Object.keys(problems).length} highlighted fields`}.`);
      document.getElementById(`field-${first}`)?.scrollIntoView({ behavior: "smooth", block: "center" });
      return;
    }

    setError("");
    setIsLoading(true);

    const [budgetMin, budgetMax] =
      form.budgetRange === CUSTOM_BUDGET && saved.budget
        ? saved.budget
        : (BUDGETS[form.budgetRange as (typeof budgetRanges)[number]] ?? [0, 0]);

    const campaignDeliverables = form.deliverables.map((d) => {
      const type = DELIVERABLE_TO_API[d as (typeof deliverables)[number]] ?? "other";
      return { type, quantity: saved.quantities[type] ?? 1, description: d };
    });

    const payload = {
      title: form.title.trim(),
      goal: form.goal,
      description: form.description.trim(),
      niche: form.niche,
      target_audience: form.targetAudience,
      city: form.location,
      platforms: form.platforms.map((p) => p.toLowerCase()),
      budget_min_etb: budgetMin,
      budget_max_etb: budgetMax,
      creator_size_preference: SIZE_TO_API[form.creatorSize as (typeof creatorSizes)[number]] ?? "any",
      number_of_creators_needed: Number(form.creatorsNeeded) || 1,
      deliverables: campaignDeliverables,
      start_date: form.startDate.slice(0, 10),
      end_date: form.endDate.slice(0, 10),
      additional_notes: form.additionalNotes,
    };

    try {
      if (campaignId) {
        await businessApi.updateCampaign(campaignId, payload);
        router.push(`/business/campaigns/${campaignId}`);
      } else {
        const res = await businessApi.createCampaign(payload);
        // New campaigns are drafts; the detail page offers Publish.
        router.push(res.data?.id ? `/business/campaigns/${res.data.id}?created=1` : "/business/campaigns");
      }
    } catch (err: unknown) {
      const data = axios.isAxiosError(err)
        ? (err.response?.data as { message?: string; error?: { type?: string; details?: Array<{ description?: string }> } } | undefined)
        : undefined;
      if (data?.error?.type === "invalid_campaign_status") {
        setError("This campaign can no longer be edited (it is completed or cancelled).");
      } else if (data?.error?.type === "account_not_active" || data?.error?.type === "profile_locked") {
        setError("Your account cannot create campaigns right now. Check your business profile.");
      } else {
        setError(data?.error?.details?.[0]?.description || data?.message || "Could not save the campaign. Please try again.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const step1Complete = useMemo(() => {
    return form.title.trim().length >= 3 && Boolean(form.goal) && form.description.trim().length >= 20;
  }, [form.title, form.goal, form.description]);

  const step2Complete = useMemo(() => {
    return Boolean(form.niche) && Boolean(form.location) && form.platforms.length > 0;
  }, [form.niche, form.location, form.platforms]);

  const step3Complete = useMemo(() => {
    return Boolean(form.budgetRange) && Boolean(form.creatorSize);
  }, [form.budgetRange, form.creatorSize]);

  const step4Complete = useMemo(() => {
    return Boolean(form.startDate) && Boolean(form.endDate) && form.endDate >= form.startDate;
  }, [form.startDate, form.endDate]);

  const completedStepsCount = useMemo(() => {
    let count = 0;
    if (step1Complete) count++;
    if (step2Complete) count++;
    if (step3Complete) count++;
    if (step4Complete) count++;
    return count;
  }, [step1Complete, step2Complete, step3Complete, step4Complete]);

  const stepsList = [
    { label: "1. Basics", complete: step1Complete },
    { label: "2. Targeting", complete: step2Complete },
    { label: "3. Budget", complete: step3Complete },
    { label: "4. Timeline", complete: step4Complete },
  ];

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-white">{isEdit ? "Edit campaign" : "Create campaign"}</h1>
          <p className="mt-1 text-sm text-[#CBD5E1]">
            {isEdit
              ? "Update the campaign details. Creators you already invited keep their invites."
              : "Saved as a draft first. Publish it from the campaign page when you’re ready to invite creators."}
          </p>
        </div>
        <span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs font-medium text-[#CBD5E1]">
          {completedStepsCount === 4 ? "Ready to save (4/4)" : `${completedStepsCount} of 4 sections complete`}
        </span>
      </section>

      <section className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
        <div className="mb-3 flex items-center justify-between">
          <p className="text-xs font-medium uppercase tracking-[0.16em] text-[#CBD5E1]">
            Campaign setup progress
          </p>
          <span className="text-xs font-medium text-[color:var(--adfluence-orange)]">
            {completedStepsCount === 4 ? "100% Ready" : `${completedStepsCount * 25}% Complete`}
          </span>
        </div>
        <div className="grid grid-cols-4 gap-2">
          {stepsList.map((step, idx) => (
            <div key={idx} className="space-y-1.5">
              <div
                className={`h-2 rounded-full transition-all duration-300 ${
                  step.complete
                    ? "bg-[color:var(--adfluence-action)] shadow-[0_0_8px_rgba(249,115,22,0.4)]"
                    : "bg-white/10"
                }`}
              />
              <p
                className={`hidden text-[11px] font-medium transition-colors sm:block ${
                  step.complete ? "text-white" : "text-[#94A3B8]"
                }`}
              >
                {step.label}
              </p>
            </div>
          ))}
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_340px]">
        <div className="space-y-6">
          <section className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:p-5">
            <h2 className="text-lg font-semibold">Campaign Basics</h2>
            <div className="grid gap-4 md:grid-cols-2">
              <label className="space-y-2 md:col-span-2">
                <span id="field-title" className="text-sm text-[#CBD5E1]">Campaign title *</span>
                {fieldError("title")}
                <Input
                  value={form.title}
                  onChange={(e) => setForm((prev) => ({ ...prev, title: e.target.value }))}
                  placeholder="e.g. Summer Product Launch"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
              <label className="space-y-2">
                <span id="field-goal" className="text-sm text-[#CBD5E1]">Campaign goal *</span>
                {fieldError("goal")}
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
                <span id="field-description" className="text-sm text-[#CBD5E1]">Campaign description *</span>
                {fieldError("description")}
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
                <span id="field-niche" className="text-sm text-[#CBD5E1]">Industry / niche *</span>
                {fieldError("niche")}
                <select
                  value={form.niche}
                  onChange={(e) => setForm((prev) => ({ ...prev, niche: e.target.value }))}
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                >
                  <option value="" disabled>
                    Select industry
                  </option>
                  {categories.map((c) => (
                    <option key={c.id} value={c.name} className="bg-[color:var(--adfluence-surface)]">
                      {c.name}
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
                <span id="field-location" className="text-sm text-[#CBD5E1]">Location *</span>
                {fieldError("location")}
                <select
                  value={form.location}
                  onChange={(e) => setForm((prev) => ({ ...prev, location: e.target.value }))}
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                >
                  <option value="" disabled>
                    Select location
                  </option>
                  {/* Keep a saved city that isn't in the preset list selectable when editing. */}
                  {[...locations, ...(form.location && !(locations as readonly string[]).includes(form.location) ? [form.location] : [])].map((location) => (
                    <option key={location} value={location} className="bg-[color:var(--adfluence-surface)]">
                      {location}
                    </option>
                  ))}
                </select>
              </label>
              <div className="space-y-2">
                <p id="field-platforms" className="text-sm text-[#CBD5E1]">Platforms *</p>
                {fieldError("platforms")}
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
                <span id="field-budgetRange" className="text-sm text-[#CBD5E1]">Budget range *</span>
                {fieldError("budgetRange")}
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
                  {form.budgetRange === CUSTOM_BUDGET && saved.budget ? (
                    <option value={CUSTOM_BUDGET} className="bg-[color:var(--adfluence-surface)]">
                      {etbRange(saved.budget[0], saved.budget[1])} (current)
                    </option>
                  ) : null}
                  {budgetRanges.map((range) => (
                    <option key={range} value={range} className="bg-[color:var(--adfluence-surface)]">
                      {range}
                    </option>
                  ))}
                </select>
              </label>
              <label className="space-y-2">
                <span id="field-creatorSize" className="text-sm text-[#CBD5E1]">Preferred creator size *</span>
                {fieldError("creatorSize")}
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
                <span id="field-creatorsNeeded" className="text-sm text-[#CBD5E1]">Number of creators needed</span>
                {fieldError("creatorsNeeded")}
                <Input
                  type="number"
                  min={1}
                  max={50}
                  inputMode="numeric"
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
                <p id="field-deliverables" className="text-sm text-[#CBD5E1]">Deliverables *</p>
                {fieldError("deliverables")}
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
                <span id="field-startDate" className="text-sm text-[#CBD5E1]">Campaign start date *</span>
                {fieldError("startDate")}
                <Input
                  type="date"
                  value={form.startDate}
                  onChange={(e) => setForm((prev) => ({ ...prev, startDate: e.target.value }))}
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                />
              </label>
              <label className="space-y-2">
                <span id="field-endDate" className="text-sm text-[#CBD5E1]">Campaign end date *</span>
                {fieldError("endDate")}
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
            <p className="rounded-xl border border-rose-400/30 bg-rose-500/10 p-4 text-sm text-rose-200">
              {error}
            </p>
          ) : null}

          <div className="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center">
            <Button
              type="button"
              onClick={onSubmit}
              disabled={isLoading || loadingCampaign}
              className="h-12 rounded-xl bg-[color:var(--adfluence-action)] px-8 text-base font-semibold text-white shadow-lg shadow-[color:var(--adfluence-orange)]/20 hover:bg-[color:var(--adfluence-action-hover)]"
            >
              {isLoading ? "Saving…" : isEdit ? "Save changes" : "Save as draft"}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => router.push(campaignId ? `/business/campaigns/${campaignId}` : "/business/campaigns")}
              className="h-12 rounded-xl border-white/10 text-sm text-[#CBD5E1] hover:bg-white/5"
            >
              Cancel
            </Button>
          </div>
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
            </div>
            <p className="text-xs text-[#CBD5E1]">
              Recommendations are based on niche, platform, budget, and creator fit.
            </p>
            {error ? (
              <p className="rounded-xl border border-rose-500/30 bg-rose-500/10 p-3 text-xs text-rose-300">
                {error}
              </p>
            ) : null}
            <Button
              type="button"
              onClick={onSubmit}
              disabled={isLoading || loadingCampaign}
              className="h-12 w-full rounded-xl bg-[color:var(--adfluence-action)] text-base font-semibold text-white shadow-lg shadow-[color:var(--adfluence-orange)]/20 hover:bg-[color:var(--adfluence-action-hover)]"
            >
              {isLoading ? "Saving…" : isEdit ? "Save changes" : "Save as draft"}
            </Button>
          </div>
        </aside>
      </div>
    </div>
  );
}
