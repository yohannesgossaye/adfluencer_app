export type AdminCategory = {
  id: string;
  name: string;
  enabled: boolean;
  influencerCount: number;
  campaignCount: number;
};

export const initialAdminCategories: AdminCategory[] = [
  { id: "cat-1", name: "Beauty", enabled: true, influencerCount: 42, campaignCount: 28 },
  { id: "cat-2", name: "Fashion", enabled: true, influencerCount: 36, campaignCount: 22 },
  { id: "cat-3", name: "Food", enabled: true, influencerCount: 29, campaignCount: 19 },
  { id: "cat-4", name: "Tech", enabled: true, influencerCount: 24, campaignCount: 31 },
  { id: "cat-5", name: "Fitness", enabled: true, influencerCount: 18, campaignCount: 11 },
  { id: "cat-6", name: "Lifestyle", enabled: true, influencerCount: 51, campaignCount: 44 },
  { id: "cat-7", name: "Travel", enabled: true, influencerCount: 15, campaignCount: 12 },
  { id: "cat-8", name: "Education", enabled: true, influencerCount: 12, campaignCount: 9 },
  { id: "cat-9", name: "Music", enabled: true, influencerCount: 21, campaignCount: 14 },
  { id: "cat-10", name: "Gaming", enabled: false, influencerCount: 8, campaignCount: 6 },
];

export const CATEGORIES_KEY = "adfluence_admin_categories";

export function loadCategories(): AdminCategory[] {
  if (typeof window === "undefined") return initialAdminCategories;
  try {
    const raw = localStorage.getItem(CATEGORIES_KEY);
    if (!raw) return initialAdminCategories;
    const parsed = JSON.parse(raw) as AdminCategory[];
    return Array.isArray(parsed) && parsed.length ? parsed : initialAdminCategories;
  } catch {
    return initialAdminCategories;
  }
}

export function saveCategories(rows: AdminCategory[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(CATEGORIES_KEY, JSON.stringify(rows));
}
