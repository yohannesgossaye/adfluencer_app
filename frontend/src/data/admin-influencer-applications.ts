export type AdminApplicationStatus = "pending" | "approved" | "rejected";

export type AdminInfluencerApplication = {
  id: string;
  name: string;
  imageUrl: string;
  niche: string;
  platform: string;
  followers: number;
  engagementRate: number;
  submittedAt: string;
  status: AdminApplicationStatus;
  completeness: string;
  bio: string;
  tiktok: string;
  instagram: string;
  youtube: string;
  audienceNotes: string;
  pricing: string;
  adminNotes: string;
};

export const APPLICATIONS_KEY = "adfluence_admin_influencer_applications";

export const initialAdminApplications: AdminInfluencerApplication[] = [
  {
    id: "app-1",
    name: "Selam Tadesse",
    imageUrl: "/images/influencers/influencer-1.jpg",
    niche: "Beauty",
    platform: "TikTok",
    followers: 62000,
    engagementRate: 6.4,
    submittedAt: "2026-05-10",
    status: "pending",
    completeness: "72% — missing pricing",
    bio: "Beauty creator focused on honest reviews and tutorials for Ethiopian skincare brands.",
    tiktok: "https://tiktok.com/@selambeauty",
    instagram: "https://instagram.com/selambeauty",
    youtube: "",
    audienceNotes: "Primarily women 18–34 in Addis Ababa and regional cities.",
    pricing: "ETB 12K–22K per package",
    adminNotes: "",
  },
  {
    id: "app-2",
    name: "Dawit Mulugeta",
    imageUrl: "/images/influencers/influencer-2.jpg",
    niche: "Tech",
    platform: "YouTube",
    followers: 54000,
    engagementRate: 5.8,
    submittedAt: "2026-05-09",
    status: "pending",
    completeness: "88%",
    bio: "Gadget reviews and app walkthroughs for local startups.",
    tiktok: "",
    instagram: "https://instagram.com/dawittech",
    youtube: "https://youtube.com/@dawittech",
    audienceNotes: "Tech-curious professionals and students.",
    pricing: "ETB 18K–32K",
    adminNotes: "",
  },
  {
    id: "app-3",
    name: "Marta Haile",
    imageUrl: "/images/influencers/influencer-3.jpg",
    niche: "Fashion",
    platform: "Instagram",
    followers: 88000,
    engagementRate: 6.9,
    submittedAt: "2026-05-08",
    status: "pending",
    completeness: "91%",
    bio: "Editorial fashion and styling content with premium brand fit.",
    tiktok: "https://tiktok.com/@martastyle",
    instagram: "https://instagram.com/martastyle",
    youtube: "",
    audienceNotes: "Urban fashion-forward audience.",
    pricing: "ETB 25K–40K",
    adminNotes: "",
  },
  {
    id: "app-4",
    name: "Kidus Alemayehu",
    imageUrl: "/images/influencers/influencer-4.jpg",
    niche: "Lifestyle",
    platform: "TikTok",
    followers: 41000,
    engagementRate: 5.5,
    submittedAt: "2026-05-07",
    status: "pending",
    completeness: "65% — add samples",
    bio: "Day-in-the-life and café culture content.",
    tiktok: "https://tiktok.com/@kiduslife",
    instagram: "https://instagram.com/kiduslife",
    youtube: "",
    audienceNotes: "Young adults 20–28.",
    pricing: "ETB 10K–18K",
    adminNotes: "",
  },
  {
    id: "app-5",
    name: "Yared Bekele",
    imageUrl: "/images/influencers/influencer-5.jpg",
    niche: "Food",
    platform: "Instagram",
    followers: 72000,
    engagementRate: 7.1,
    submittedAt: "2026-05-06",
    status: "pending",
    completeness: "84%",
    bio: "Restaurant reviews and recipe shorts.",
    tiktok: "",
    instagram: "https://instagram.com/yaredfood",
    youtube: "https://youtube.com/@yaredfood",
    audienceNotes: "Food explorers across Addis.",
    pricing: "ETB 15K–28K",
    adminNotes: "",
  },
  {
    id: "app-6",
    name: "Helina Solomon",
    imageUrl: "/images/influencers/influencer-1.jpg",
    niche: "Fitness",
    platform: "YouTube",
    followers: 35000,
    engagementRate: 6.0,
    submittedAt: "2026-05-05",
    status: "pending",
    completeness: "79%",
    bio: "Home workouts and wellness tips.",
    tiktok: "https://tiktok.com/@helinafit",
    instagram: "https://instagram.com/helinafit",
    youtube: "https://youtube.com/@helinafit",
    audienceNotes: "Health-conscious mixed gender 22–40.",
    pricing: "ETB 9K–16K",
    adminNotes: "",
  },
  {
    id: "app-7",
    name: "Nahom Girma",
    imageUrl: "/images/influencers/influencer-2.jpg",
    niche: "Travel",
    platform: "Instagram",
    followers: 95000,
    engagementRate: 6.3,
    submittedAt: "2026-05-04",
    status: "pending",
    completeness: "93%",
    bio: "Domestic travel guides and eco-tourism features.",
    tiktok: "",
    instagram: "https://instagram.com/nahomtravel",
    youtube: "https://youtube.com/@nahomtravel",
    audienceNotes: "Adventure seekers nationwide.",
    pricing: "ETB 22K–38K",
    adminNotes: "",
  },
  {
    id: "app-8",
    name: "Liya Worku",
    imageUrl: "/images/influencers/influencer-3.jpg",
    niche: "Beauty",
    platform: "TikTok",
    followers: 48000,
    engagementRate: 5.9,
    submittedAt: "2026-05-03",
    status: "pending",
    completeness: "70%",
    bio: "Makeup tutorials and GRWM content.",
    tiktok: "https://tiktok.com/@liyaglam",
    instagram: "https://instagram.com/liyaglam",
    youtube: "",
    audienceNotes: "Gen Z beauty audience.",
    pricing: "ETB 11K–20K",
    adminNotes: "",
  },
  {
    id: "app-9",
    name: "Biniam Fikru",
    imageUrl: "/images/influencers/influencer-4.jpg",
    niche: "Music",
    platform: "YouTube",
    followers: 120000,
    engagementRate: 5.2,
    submittedAt: "2026-04-28",
    status: "approved",
    completeness: "95%",
    bio: "Live sessions and artist collaborations.",
    tiktok: "",
    instagram: "https://instagram.com/biniammusic",
    youtube: "https://youtube.com/@biniammusic",
    audienceNotes: "Music fans 16–35.",
    pricing: "ETB 30K–55K",
    adminNotes: "",
  },
  {
    id: "app-10",
    name: "Rahel Desta",
    imageUrl: "/images/influencers/influencer-5.jpg",
    niche: "Fashion",
    platform: "Instagram",
    followers: 22000,
    engagementRate: 4.8,
    submittedAt: "2026-04-20",
    status: "rejected",
    completeness: "55%",
    bio: "Street style content.",
    tiktok: "",
    instagram: "https://instagram.com/raheld",
    youtube: "",
    audienceNotes: "Limited reach vs policy minimum.",
    pricing: "ETB 6K–10K",
    adminNotes: "",
  },
];

export function loadApplications(): AdminInfluencerApplication[] {
  if (typeof window === "undefined") return initialAdminApplications;
  try {
    const raw = localStorage.getItem(APPLICATIONS_KEY);
    if (!raw) return initialAdminApplications;
    const parsed = JSON.parse(raw) as AdminInfluencerApplication[];
    if (!Array.isArray(parsed) || !parsed.length) return initialAdminApplications;
    return parsed.map((p) => ({ ...p, adminNotes: p.adminNotes ?? "" }));
  } catch {
    return initialAdminApplications;
  }
}

export function saveApplications(list: AdminInfluencerApplication[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(APPLICATIONS_KEY, JSON.stringify(list));
}
