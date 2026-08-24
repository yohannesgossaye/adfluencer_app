import { apiClient } from "@/lib/api/api-client";
import type { ApiResponse } from "@/lib/api/auth";

// --- Dashboard Types ---
export type DashboardGreeting = {
  business_name: string;
  owner_full_name: string;
};

export type DashboardKPIs = {
  active_campaigns: number;
  total_campaigns: number;
  recommended_creators_total: number;
  recommended_creators_high_match: number;
  pending_contracts: number;
  contracts_awaiting_influencer: number;
  contracts_awaiting_payment: number;
  pending_invites: number;
};

export type DashboardContractsSummary = {
  pending: number;
  active: number;
  completed: number;
};

export type DashboardCampaignRow = {
  id: string;
  title: string;
  niche?: string;
  budget_summary?: string;
  status: string;
  progress_percent?: number;
  accepted_creators: number;
  creators_needed: number;
  updated_at: string;
};

export type DashboardRecommendationRow = {
  influencer_user_id: string;
  creator_name: string;
  profile_image_url?: string;
  match_score: number;
  campaign_id: string;
  campaign_title: string;
};

export type MarketplaceAccessResponse = {
  visible_in_marketplace: boolean;
  can_create_contracts: boolean;
  message?: string;
};

export type DashboardSummaryResponse = {
  greeting: DashboardGreeting;
  marketplace_access: MarketplaceAccessResponse;
  kpis: DashboardKPIs;
  contracts_summary: DashboardContractsSummary;
  recent_campaigns: DashboardCampaignRow[];
  top_recommendations: DashboardRecommendationRow[];
};

// --- Profile Types ---
export type BusinessSocialLinks = {
  instagram?: string;
  tiktok?: string;
  facebook?: string;
  linkedin?: string;
};

export type BusinessProfileDetail = {
  id: string;
  owner_full_name: string;
  business_name: string;
  legal_business_name?: string;
  logo_url?: string;
  description?: string;
  industry?: string;
  business_type?: string;
  business_size?: string;
  city?: string;
  region?: string;
  address?: string;
  tin_number?: string;
  website_url?: string;
  social_links: BusinessSocialLinks;
  preferred_platforms: string[];
  monthly_marketing_budget?: number;
  is_verified: boolean;
};

export type BusinessCompleteness = {
  percent: number;
  missing: string[];
  can_complete: boolean;
};

export type BusinessProfileResponse = {
  user: {
    id: string;
    email: string;
    phone_number?: string;
    status: string;
    email_verified: boolean;
    phone_verified: boolean;
  };
  profile: BusinessProfileDetail;
  completeness: BusinessCompleteness;
  marketplace_access: MarketplaceAccessResponse;
};

export type UpdateBusinessProfileRequest = {
  owner_full_name?: string;
  business_name?: string;
  legal_business_name?: string;
  description?: string;
  industry?: string;
  business_type?: string;
  business_size?: string;
  city?: string;
  region?: string;
  address?: string;
  tin_number?: string;
  website_url?: string;
  social_links?: BusinessSocialLinks;
  preferred_platforms?: string[];
  monthly_marketing_budget?: number;
};

// --- Campaign Types ---
export type CampaignDeliverable = {
  id?: string;
  type: string;
  quantity: number;
  description?: string;
  status?: string;
};

export type CreateCampaignRequest = {
  title: string;
  goal: string;
  description: string;
  niche: string;
  target_audience?: string;
  city: string;
  region?: string;
  platforms: string[];
  budget_min_etb?: number;
  budget_max_etb?: number;
  creator_size_preference: string;
  number_of_creators_needed: number;
  deliverables: CampaignDeliverable[];
  start_date: string;
  end_date: string;
  additional_notes?: string;
};

export type CampaignDetail = {
  id: string;
  title: string;
  goal?: string;
  description: string;
  category?: { id: string; name: string; slug: string };
  target_audience?: string;
  city?: string;
  region?: string;
  platforms: string[];
  budget_min_etb?: number;
  budget_max_etb?: number;
  creator_size_preference: string;
  number_of_creators_needed: number;
  deliverables: CampaignDeliverable[];
  start_date?: string;
  end_date?: string;
  status: string;
  business_name: string;
  invites_sent: number;
  invites_accepted: number;
  created_at: string;
  updated_at: string;
};

export type CampaignListItem = {
  id: string;
  title: string;
  business_name: string;
  primary_platform?: string;
  budget_summary?: string;
  status: string;
  start_date?: string;
  end_date?: string;
  creators_needed: number;
  invites_sent: number;
  invites_accepted: number;
  created_at: string;
};

// --- Recommendation Types ---
export type ScoreBreakdown = {
  niche: number;
  platform: number;
  budget: number;
  creator_size: number;
  verified: number;
};

export type RecommendationItem = {
  influencer_user_id: string;
  creator_name: string;
  profile_image_url?: string;
  primary_niche?: string;
  platforms: string[];
  total_followers: number;
  engagement_rate?: number;
  min_price_etb?: number;
  max_price_etb?: number;
  match_score: number;
  score_breakdown: ScoreBreakdown;
  is_verified: boolean;
  invite_status?: string;
  rank_position?: number;
};

export type InfluencerViewDetail = {
  influencer_user_id: string;
  creator_name: string;
  full_name: string;
  bio?: string;
  profile_image_url?: string;
  city?: string;
  region?: string;
  audience_description?: string;
  languages: string[];
  niches: Array<{ id: string; name: string; slug: string }>;
  platforms: string[];
  social_links: {
    tiktok?: string;
    instagram?: string;
    youtube?: string;
    telegram?: string;
  };
  total_followers: number;
  engagement_rate?: number;
  min_price_etb?: number;
  max_price_etb?: number;
  is_verified: boolean;
  portfolio_items: Array<{
    id: string;
    file_url?: string;
    media_url?: string;
    media_type: string;
    thumbnail_url?: string;
    title?: string;
    caption?: string;
  }>;
  match_score?: number;
  score_breakdown?: ScoreBreakdown;
};

// --- Invite Types ---
export type InviteResponse = {
  id: string;
  campaign_id: string;
  influencer_user_id: string;
  creator_name?: string;
  status: string;
  message?: string;
  created_at: string;
};

// --- Contract Types ---
export type FeeBreakdown = {
  amount_etb: number;
  business_fee_percent: number;
  influencer_fee_percent: number;
  business_fee_amount_etb: number;
  influencer_fee_amount_etb: number;
  total_business_payable_etb: number;
  influencer_receivable_etb: number;
  platform_total_fee_etb: number;
};

export type ContractDetail = {
  id: string;
  title: string;
  description?: string;
  status: string;
  fees: FeeBreakdown;
  deliverables: CampaignDeliverable[];
  campaign: { id: string; title: string };
  business: { user_id: string; business_name?: string; profile_image_url?: string };
  influencer: { user_id: string; creator_name?: string; profile_image_url?: string };
  start_date?: string;
  end_date?: string;
  payment?: {
    id: string;
    status: string;
    total_paid_by_business_etb: number;
    payment_method: string;
    paid_at?: string;
  };
  payout?: {
    id: string;
    status: string;
    amount_etb: number;
    paid_at?: string;
  };
  accepted_by_influencer_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
};

export type CreateContractRequest = {
  invite_id: string;
  title: string;
  description?: string;
  amount_etb: number;
  deliverables: CampaignDeliverable[];
  start_date?: string;
  end_date?: string;
};

// --- Support Types ---
export type SupportTicket = {
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

export type CreateSupportTicketRequest = {
  issue_type: string;
  subject: string;
  message: string;
  priority: string;
  related_campaign_id?: string;
};

// --- API Service Methods ---
export const businessApi = {
  // 1. Dashboard Summary
  getDashboardSummary: async (): Promise<ApiResponse<DashboardSummaryResponse>> => {
    const res = await apiClient.get<ApiResponse<DashboardSummaryResponse>>(
      "/api/v1/business/dashboard/summary"
    );
    return res.data;
  },

  // 2. Profile CRUD & Logo
  getProfile: async (): Promise<ApiResponse<BusinessProfileResponse>> => {
    const res = await apiClient.get<ApiResponse<BusinessProfileResponse>>(
      "/api/v1/business/profile"
    );
    return res.data;
  },

  updateProfile: async (
    data: UpdateBusinessProfileRequest
  ): Promise<ApiResponse<BusinessProfileResponse>> => {
    const res = await apiClient.put<ApiResponse<BusinessProfileResponse>>(
      "/api/v1/business/profile",
      data
    );
    return res.data;
  },

  uploadLogo: async (file: File): Promise<ApiResponse<{ logo_url: string }>> => {
    const formData = new FormData();
    formData.append("file", file);
    const res = await apiClient.post<ApiResponse<{ logo_url: string }>>(
      "/api/v1/business/profile/logo",
      formData,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    return res.data;
  },

  completeProfile: async (): Promise<ApiResponse<{ status: string; message: string; is_verified: boolean }>> => {
    const res = await apiClient.post<ApiResponse<{ status: string; message: string; is_verified: boolean }>>(
      "/api/v1/business/profile/complete"
    );
    return res.data;
  },

  // 3. Campaigns
  createCampaign: async (data: CreateCampaignRequest): Promise<ApiResponse<{ campaign: CampaignDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ campaign: CampaignDetail }>>(
      "/api/v1/business/campaigns",
      data
    );
    return res.data;
  },

  listCampaigns: async (params?: {
    status?: string;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<{ items: CampaignListItem[]; pagination: { total: number; page: number; limit: number; total_pages: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: CampaignListItem[]; pagination: { total: number; page: number; limit: number; total_pages: number } }>>(
      "/api/v1/business/campaigns",
      { params }
    );
    return res.data;
  },

  getCampaign: async (campaignId: string): Promise<ApiResponse<{ campaign: CampaignDetail }>> => {
    const res = await apiClient.get<ApiResponse<{ campaign: CampaignDetail }>>(
      `/api/v1/business/campaigns/${campaignId}`
    );
    return res.data;
  },

  updateCampaign: async (
    campaignId: string,
    data: CreateCampaignRequest
  ): Promise<ApiResponse<{ campaign: CampaignDetail }>> => {
    const res = await apiClient.put<ApiResponse<{ campaign: CampaignDetail }>>(
      `/api/v1/business/campaigns/${campaignId}`,
      data
    );
    return res.data;
  },

  publishCampaign: async (campaignId: string): Promise<ApiResponse<{ campaign: CampaignDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ campaign: CampaignDetail }>>(
      `/api/v1/business/campaigns/${campaignId}/publish`
    );
    return res.data;
  },

  cancelCampaign: async (campaignId: string): Promise<ApiResponse<{ campaign: CampaignDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ campaign: CampaignDetail }>>(
      `/api/v1/business/campaigns/${campaignId}/cancel`
    );
    return res.data;
  },

  // 4. Smart Recommendations & Discovery
  generateRecommendations: async (
    campaignId: string
  ): Promise<ApiResponse<{ generated: number; total: number }>> => {
    const res = await apiClient.post<ApiResponse<{ generated: number; total: number }>>(
      `/api/v1/business/campaigns/${campaignId}/recommendations/generate`
    );
    return res.data;
  },

  listRecommendations: async (
    campaignId: string,
    params?: {
      search?: string;
      niche?: string;
      platform?: string;
      min_score?: number;
      sort?: string;
      page?: number;
      limit?: number;
    }
  ): Promise<ApiResponse<{ items: RecommendationItem[]; pagination: { total: number; page: number; limit: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: RecommendationItem[]; pagination: { total: number; page: number; limit: number } }>>(
      `/api/v1/business/campaigns/${campaignId}/recommendations`,
      { params }
    );
    return res.data;
  },

  getInfluencerForBusiness: async (
    userId: string,
    campaignId?: string
  ): Promise<ApiResponse<InfluencerViewDetail>> => {
    const res = await apiClient.get<ApiResponse<InfluencerViewDetail>>(
      `/api/v1/business/influencers/${userId}`,
      { params: campaignId ? { campaign_id: campaignId } : undefined }
    );
    return res.data;
  },

  // 5. Invites
  createInvite: async (
    campaignId: string,
    data: { influencer_user_id: string; message?: string }
  ): Promise<ApiResponse<InviteResponse>> => {
    const res = await apiClient.post<ApiResponse<InviteResponse>>(
      `/api/v1/business/campaigns/${campaignId}/invites`,
      data
    );
    return res.data;
  },

  listInvites: async (campaignId: string): Promise<ApiResponse<{ items: InviteResponse[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: InviteResponse[] }>>(
      `/api/v1/business/campaigns/${campaignId}/invites`
    );
    return res.data;
  },

  cancelInvite: async (
    campaignId: string,
    inviteId: string
  ): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.delete<ApiResponse<{ message: string }>>(
      `/api/v1/business/campaigns/${campaignId}/invites/${inviteId}`
    );
    return res.data;
  },

  // 6. Contracts
  createContract: async (
    data: CreateContractRequest
  ): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ contract: ContractDetail }>>(
      "/api/v1/business/contracts",
      data
    );
    return res.data;
  },

  listContracts: async (): Promise<ApiResponse<{ items: ContractDetail[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: ContractDetail[] }>>(
      "/api/v1/business/contracts"
    );
    return res.data;
  },

  getContract: async (contractId: string): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.get<ApiResponse<{ contract: ContractDetail }>>(
      `/api/v1/business/contracts/${contractId}`
    );
    return res.data;
  },

  payContract: async (
    contractId: string,
    paymentMethod: string = "manual_stub"
  ): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ contract: ContractDetail }>>(
      `/api/v1/business/contracts/${contractId}/pay`,
      { payment_method: paymentMethod }
    );
    return res.data;
  },

  approveWork: async (contractId: string): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ contract: ContractDetail }>>(
      `/api/v1/business/contracts/${contractId}/approve-work`
    );
    return res.data;
  },

  cancelContract: async (contractId: string): Promise<ApiResponse<{ contract: ContractDetail }>> => {
    const res = await apiClient.post<ApiResponse<{ contract: ContractDetail }>>(
      `/api/v1/business/contracts/${contractId}/cancel`
    );
    return res.data;
  },

  // 7. Support
  listSupportCampaignOptions: async (): Promise<ApiResponse<{ items: Array<{ id: string; title: string }> }>> => {
    const res = await apiClient.get<ApiResponse<{ items: Array<{ id: string; title: string }> }>>(
      "/api/v1/business/support/campaign-options"
    );
    return res.data;
  },

  createSupportTicket: async (
    data: CreateSupportTicketRequest
  ): Promise<ApiResponse<{ id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; status: string }>>(
      "/api/v1/business/support/tickets",
      data
    );
    return res.data;
  },

  listSupportTickets: async (): Promise<ApiResponse<{ items: SupportTicket[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: SupportTicket[] }>>(
      "/api/v1/business/support/tickets"
    );
    return res.data;
  },

  getSupportTicket: async (ticketId: string): Promise<ApiResponse<SupportTicket>> => {
    const res = await apiClient.get<ApiResponse<SupportTicket>>(
      `/api/v1/business/support/tickets/${ticketId}`
    );
    return res.data;
  },

  listCategories: async (): Promise<ApiResponse<{ items: Array<{ id: string; name: string; slug: string }> }>> => {
    const res = await apiClient.get<ApiResponse<{ items: Array<{ id: string; name: string; slug: string }> }>>(
      "/api/v1/business/categories"
    );
    return res.data;
  },
};
