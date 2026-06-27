export type AdminCampaignStatus = "draft" | "pending" | "active" | "completed" | "flagged";

export type AdminCampaignRecord = {
  id: string;
  title: string;
  business: string;
  niche: string;
  budget: string;
  creatorsInvited: number;
  status: AdminCampaignStatus;
  createdAt: string;
  contractStatus: string;
  paymentStatus: string;
};

export const initialAdminCampaigns: AdminCampaignRecord[] = [
  {
    id: "adc-1",
    title: "Betu Cosmetics Launch",
    business: "Betu Cosmetics",
    niche: "Beauty",
    budget: "ETB 280K",
    creatorsInvited: 6,
    status: "active",
    createdAt: "2026-04-28",
    contractStatus: "2 active, 1 pending",
    paymentStatus: "Escrow funded",
  },
  {
    id: "adc-2",
    title: "Ride App Creator Push",
    business: "Ride Ethiopia",
    niche: "Tech",
    budget: "ETB 340K",
    creatorsInvited: 9,
    status: "pending",
    createdAt: "2026-05-02",
    contractStatus: "Awaiting signatures",
    paymentStatus: "Pending deposit",
  },
  {
    id: "adc-3",
    title: "Coffee House Awareness",
    business: "Habesha Coffee Roasters",
    niche: "Food",
    budget: "ETB 190K",
    creatorsInvited: 4,
    status: "completed",
    createdAt: "2026-03-10",
    contractStatus: "Closed",
    paymentStatus: "Released",
  },
  {
    id: "adc-4",
    title: "Sheba Tech Week",
    business: "Sheba Tech Hub",
    niche: "Tech",
    budget: "ETB 120K",
    creatorsInvited: 5,
    status: "flagged",
    createdAt: "2026-05-09",
    contractStatus: "Policy review",
    paymentStatus: "On hold",
  },
  {
    id: "adc-5",
    title: "Spring Fashion Edit",
    business: "Abay Fashion House",
    niche: "Fashion",
    budget: "ETB 210K",
    creatorsInvited: 3,
    status: "draft",
    createdAt: "2026-05-11",
    contractStatus: "None",
    paymentStatus: "Not started",
  },
];

export const CAMPAIGNS_STORAGE_KEY = "adfluence_admin_campaigns";

export function loadAdminCampaigns(): AdminCampaignRecord[] {
  if (typeof window === "undefined") return initialAdminCampaigns;
  try {
    const raw = localStorage.getItem(CAMPAIGNS_STORAGE_KEY);
    if (!raw) return initialAdminCampaigns;
    const parsed = JSON.parse(raw) as AdminCampaignRecord[];
    return Array.isArray(parsed) && parsed.length ? parsed : initialAdminCampaigns;
  } catch {
    return initialAdminCampaigns;
  }
}

export function saveAdminCampaigns(rows: AdminCampaignRecord[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(CAMPAIGNS_STORAGE_KEY, JSON.stringify(rows));
}
