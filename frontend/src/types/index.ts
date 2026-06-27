export type UserRole = "business" | "influencer" | "admin";

export type InfluencerStatus = "pending" | "approved" | "rejected";
export type CampaignStatus =
  | "draft"
  | "pending"
  | "active"
  | "completed"
  | "cancelled";
export type ContractStatus =
  | "pending"
  | "accepted"
  | "rejected"
  | "active"
  | "completed";

export type Influencer = {
  id: string;
  name: string;
  niche:
    | "beauty"
    | "fashion"
    | "food"
    | "fitness"
    | "tech"
    | "lifestyle"
    | "travel";
  city: string;
  platforms: Array<"TikTok" | "Instagram" | "YouTube">;
  followers: number;
  engagementRate: number;
  status: InfluencerStatus;
  imageUrl?: string;
  matchScore?: number;
  priceRangeEtb?: string;
};

export type Campaign = {
  id: string;
  businessName: string;
  title: string;
  budgetEtb: number;
  status: CampaignStatus;
  platform: "TikTok" | "Instagram" | "YouTube";
};

export type Message = {
  id: string;
  from: string;
  to: string;
  subject: string;
  preview: string;
  unread: boolean;
  timestamp: string;
};

export type Contract = {
  id: string;
  campaignTitle: string;
  businessName: string;
  influencerName: string;
  amountEtb: number;
  status: ContractStatus;
};

export type StatItem = {
  id: string;
  label: string;
  value: string;
  hint?: string;
};
