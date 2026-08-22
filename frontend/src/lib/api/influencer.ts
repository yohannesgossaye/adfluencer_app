import { apiClient } from "@/lib/api/api-client";
import type { ApiResponse } from "@/lib/api/auth";
import type { ContractDetail } from "@/lib/api/business";

export type { ContractDetail };

// --- Dashboard Types ---
export type InfluencerGreeting = {
  creator_name: string;
};

export type InfluencerKPIs = {
  pending_invites: number;
  active_collaborations: number;
  pending_contracts: number;
  completed_contracts: number;
  profile_score_percent: number;
  total_completed_earnings_etb: number;
};

export type InfluencerContractsSummary = {
  pending: number;
  active: number;
  completed: number;
};

export type InfluencerEarningsSummary = {
  pending_payout_etb: number;
  completed_earnings_etb: number;
  next_payout_at?: string;
};

export type InfluencerDashboardInvite = {
  id: string;
  business_name: string;
  campaign_title: string;
  status: string;
  budget_summary?: string;
  end_date?: string;
};

export type InfluencerDashboardContract = {
  id: string;
  business_name: string;
  status: string;
};

export type InfluencerCompleteness = {
  percent: number;
  missing: string[];
  can_submit: boolean;
};

export type InfluencerDashboardSummaryResponse = {
  greeting: InfluencerGreeting;
  completeness: InfluencerCompleteness;
  kpis: InfluencerKPIs;
  contracts_summary: InfluencerContractsSummary;
  earnings_summary: InfluencerEarningsSummary;
  pending_invites: InfluencerDashboardInvite[];
  active_contracts: InfluencerDashboardContract[];
};

// --- Profile Types ---
export type InfluencerSocialLinks = {
  tiktok?: string;
  instagram?: string;
  youtube?: string;
  telegram?: string;
};

export type InfluencerCategory = {
  id: string;
  name: string;
  slug: string;
};

export type InfluencerPortfolioItem = {
  id: string;
  file_url: string;
  media_type: string;
  mime_type: string;
  file_size_bytes: number;
  sort_order: number;
  created_at: string;
};

export type InfluencerProfileDetail = {
  id: string;
  creator_name: string;
  full_name: string;
  username?: string;
  bio?: string;
  profile_image_url?: string;
  city?: string;
  region?: string;
  languages: string[];
  niches: InfluencerCategory[];
  platforms: string[];
  social_links: InfluencerSocialLinks;
  instagram_followers: number;
  tiktok_followers: number;
  youtube_subscribers: number;
  total_followers: number;
  engagement_rate?: number;
  min_price_etb?: number;
  max_price_etb?: number;
  audience_description?: string;
  is_verified: boolean;
  approved_at?: string;
  portfolio_items: InfluencerPortfolioItem[];
};

export type InfluencerProfileResponse = {
  user: {
    id: string;
    email: string;
    phone_number?: string;
    status: string;
    status_reason?: string;
    email_verified: boolean;
    phone_verified: boolean;
  };
  profile: InfluencerProfileDetail;
  completeness: InfluencerCompleteness;
};

export type UpdateInfluencerProfileRequest = {
  creator_name?: string;
  bio?: string;
  niches?: string[];
  platforms?: string[];
  social_links?: InfluencerSocialLinks;
  followers?: number;
  instagram_followers?: number;
  tiktok_followers?: number;
  youtube_subscribers?: number;
  engagement_rate?: number;
  min_price_etb?: number;
  max_price_etb?: number;
  languages?: string[];
  city?: string;
  region?: string;
  audience_description?: string;
};

// --- Invite Types ---
export type InfluencerInviteCampaign = {
  id: string;
  title: string;
  business_name: string;
  status: string;
  budget_min_etb?: number;
  budget_max_etb?: number;
  platforms: string[];
  start_date?: string;
  end_date?: string;
};

export type InfluencerInviteItem = {
  id: string;
  status: string;
  message?: string;
  campaign: InfluencerInviteCampaign;
  created_at: string;
  responded_at?: string;
};

// --- Support Types ---
export type InfluencerSupportTicket = {
  id: string;
  subject: string;
  message?: string;
  issue_type: string;
  issue_type_label?: string;
  priority: string;
  status: string;
  admin_reply?: string;
  related_campaign_id?: string;
  related_campaign_title?: string;
  created_at: string;
};

export type CreateInfluencerSupportTicketRequest = {
  issue_type: string;
  subject: string;
  message: string;
  priority: string;
  related_invite_id?: string;
};

// --- API Service Methods ---
export const influencerApi = {
  // 1. Dashboard Summary
  getDashboardSummary: async (): Promise<ApiResponse<InfluencerDashboardSummaryResponse>> => {
    const res = await apiClient.get<ApiResponse<InfluencerDashboardSummaryResponse>>(
      "/api/v1/influencer/dashboard/summary"
    );
    return res.data;
  },

  // 2. Profile CRUD & Media
  getProfile: async (): Promise<ApiResponse<InfluencerProfileResponse>> => {
    const res = await apiClient.get<ApiResponse<InfluencerProfileResponse>>(
      "/api/v1/influencer/profile"
    );
    return res.data;
  },

  updateProfile: async (
    data: UpdateInfluencerProfileRequest
  ): Promise<ApiResponse<InfluencerProfileResponse>> => {
    const res = await apiClient.put<ApiResponse<InfluencerProfileResponse>>(
      "/api/v1/influencer/profile",
      data
    );
    return res.data;
  },

  submitProfile: async (): Promise<ApiResponse<{ status: string; message: string; submitted_at: string }>> => {
    const res = await apiClient.post<ApiResponse<{ status: string; message: string; submitted_at: string }>>(
      "/api/v1/influencer/profile/submit"
    );
    return res.data;
  },

  uploadAvatar: async (file: File): Promise<ApiResponse<{ profile_image_url: string }>> => {
    const formData = new FormData();
    formData.append("file", file);
    const res = await apiClient.post<ApiResponse<{ profile_image_url: string }>>(
      "/api/v1/influencer/profile/avatar",
      formData,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    return res.data;
  },

  deleteAvatar: async (): Promise<ApiResponse<null>> => {
    const res = await apiClient.delete<ApiResponse<null>>(
      "/api/v1/influencer/profile/avatar"
    );
    return res.data;
  },

  uploadPortfolio: async (
    file: File
  ): Promise<ApiResponse<{ item: InfluencerPortfolioItem; portfolio_items: InfluencerPortfolioItem[]; portfolio_count: number }>> => {
    const formData = new FormData();
    formData.append("file", file);
    const res = await apiClient.post<ApiResponse<{ item: InfluencerPortfolioItem; portfolio_items: InfluencerPortfolioItem[]; portfolio_count: number }>>(
      "/api/v1/influencer/profile/portfolio",
      formData,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    return res.data;
  },

  deletePortfolio: async (
    itemId: string
  ): Promise<ApiResponse<{ portfolio_items: InfluencerPortfolioItem[]; portfolio_count: number }>> => {
    const res = await apiClient.delete<ApiResponse<{ portfolio_items: InfluencerPortfolioItem[]; portfolio_count: number }>>(
      `/api/v1/influencer/profile/portfolio/${itemId}`
    );
    return res.data;
  },

  // 3. Invites
  listInvites: async (params?: {
    status?: string;
  }): Promise<ApiResponse<{ items: InfluencerInviteItem[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: InfluencerInviteItem[] }>>(
      "/api/v1/influencer/invites",
      { params }
    );
    return res.data;
  },

  getInvite: async (inviteId: string): Promise<ApiResponse<InfluencerInviteItem>> => {
    const res = await apiClient.get<ApiResponse<InfluencerInviteItem>>(
      `/api/v1/influencer/invites/${inviteId}`
    );
    return res.data;
  },

  acceptInvite: async (
    inviteId: string
  ): Promise<ApiResponse<{ id: string; status: string; responded_at: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; status: string; responded_at: string }>>(
      `/api/v1/influencer/invites/${inviteId}/accept`
    );
    return res.data;
  },

  rejectInvite: async (
    inviteId: string
  ): Promise<ApiResponse<{ id: string; status: string; responded_at: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; status: string; responded_at: string }>>(
      `/api/v1/influencer/invites/${inviteId}/reject`
    );
    return res.data;
  },

  // 4. Contracts
  listContracts: async (): Promise<ApiResponse<{ items: ContractDetail[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: ContractDetail[] }>>(
      "/api/v1/influencer/contracts"
    );
    return res.data;
  },

  getContract: async (contractId: string): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.get<ApiResponse<{ contract: ContractDetail }>>(
      `/api/v1/influencer/contracts/${contractId}`
    );
    return res.data;
  },

  acceptContract: async (
    contractId: string
  ): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ contract: ContractDetail }>>(
      `/api/v1/influencer/contracts/${contractId}/accept`
    );
    return res.data;
  },

  rejectContract: async (
    contractId: string
  ): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ contract: ContractDetail }>>(
      `/api/v1/influencer/contracts/${contractId}/reject`
    );
    return res.data;
  },

  // 5. Support
  listSupportInviteOptions: async (): Promise<ApiResponse<{ items: Array<{ id: string; campaign_title: string; business_name: string }> }>> => {
    const res = await apiClient.get<ApiResponse<{ items: Array<{ id: string; campaign_title: string; business_name: string }> }>>(
      "/api/v1/influencer/support/invite-options"
    );
    return res.data;
  },

  createSupportTicket: async (
    data: CreateInfluencerSupportTicketRequest
  ): Promise<ApiResponse<{ id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; status: string }>>(
      "/api/v1/influencer/support/tickets",
      data
    );
    return res.data;
  },

  listSupportTickets: async (): Promise<ApiResponse<{ items: InfluencerSupportTicket[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: InfluencerSupportTicket[] }>>(
      "/api/v1/influencer/support/tickets"
    );
    return res.data;
  },

  getSupportTicket: async (ticketId: string): Promise<ApiResponse<InfluencerSupportTicket>> => {
    const res = await apiClient.get<ApiResponse<InfluencerSupportTicket>>(
      `/api/v1/influencer/support/tickets/${ticketId}`
    );
    return res.data;
  },

  // 6. Categories / Niches
  listCategories: async (): Promise<ApiResponse<{ items: Array<{ id: string; name: string; slug: string }> }>> => {
    const res = await apiClient.get<ApiResponse<{ items: Array<{ id: string; name: string; slug: string }> }>>(
      "/api/v1/influencer/categories"
    );
    return res.data;
  },
};
