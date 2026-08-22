import { apiClient } from "@/lib/api/api-client";
import type { ApiResponse } from "@/lib/api/auth";

// --- Dashboard Types ---
export type AdminKPIs = {
  total_influencers: number;
  total_businesses: number;
  pending_influencer_approvals: number;
  active_campaigns: number;
  total_campaigns: number;
  total_users: number;
  open_disputes: number;
  open_support_tickets: number;
  blocked_users: number;
  completed_contracts: number;
  monthly_platform_revenue_etb: number;
};

export type AdminRevenueSummary = {
  total_campaign_value_etb: number;
  total_platform_fees_etb: number;
  platform_fee_percent_business: number;
  platform_fee_percent_influencer: number;
  monthly_trend: Array<{
    month: string;
    platform_revenue_etb: number;
    completed_contracts: number;
  }>;
};

export type AdminPendingApproval = {
  user_id: string;
  name: string;
  niche?: string;
  platform?: string;
  followers: number;
  submitted_at: string;
};

export type AdminActivity = {
  id: string;
  action: string;
  label: string;
  detail: string;
  created_at: string;
};

export type AdminAlert = {
  id: string;
  title: string;
  detail: string;
  severity: string;
  entity_type?: string;
};

export type AdminDashboardSummaryResponse = {
  kpis: AdminKPIs;
  pending_approvals: AdminPendingApproval[];
  recent_activity: AdminActivity[];
  revenue_summary: AdminRevenueSummary;
  alerts: AdminAlert[];
};

// --- Influencer Approval Types ---
export type InfluencerApplicationItem = {
  user_id: string;
  creator_name: string;
  profile_image_url?: string;
  primary_niche?: string;
  primary_platform?: string;
  total_followers: number;
  engagement_rate?: number;
  completeness_percent: number;
  pricing_summary?: string;
  status: string;
  submitted_at: string;
  bio_preview?: string;
};

export type InfluencerApplicationDetail = {
  user_id: string;
  status: string;
  status_reason?: string;
  email: string;
  phone_number?: string;
  profile: {
    creator_name: string;
    bio?: string;
    profile_image_url?: string;
    city?: string;
    niches: Array<{ id: string; name: string; slug: string }>;
    platforms: string[];
    social_links: Record<string, string>;
    total_followers: number;
    engagement_rate?: number;
    min_price_etb?: number;
    max_price_etb?: number;
    audience_description?: string;
    languages: string[];
    portfolio_items: Array<{
      id: string;
      file_url: string;
      media_type: string;
      mime_type: string;
      file_size_bytes: number;
      sort_order: number;
      created_at: string;
    }>;
  };
  completeness_percent: number;
  admin_note?: string;
  submitted_at: string;
};

// --- Business Verification Types ---
export type BusinessVerificationItem = {
  user_id: string;
  business_name: string;
  owner_full_name: string;
  industry?: string;
  city?: string;
  logo_url?: string;
  email: string;
  is_verified: boolean;
  registered_at: string;
};

export type BusinessVerificationDetail = {
  user_id: string;
  email: string;
  phone_number?: string;
  status: string;
  business_name: string;
  owner_full_name: string;
  industry?: string;
  city?: string;
  description?: string;
  logo_url?: string;
  is_verified: boolean;
  registered_at: string;
};

// --- User Governance Types ---
export type AdminUserListItem = {
  user_id: string;
  name: string;
  email: string;
  role: string;
  status: string;
  joined_at: string;
  last_active_at?: string;
  campaigns_count: number;
  contracts_count: number;
};

export type AdminUserStats = {
  total: number;
  active: number;
  blocked: number;
  pending: number;
};

// --- Category Types ---
export type AdminCategoryItem = {
  id: string;
  name: string;
  slug: string;
  description?: string;
  is_active: boolean;
  influencer_count: number;
  campaign_count: number;
  created_at: string;
  updated_at: string;
};

// --- Campaign Supervision Types ---
export type AdminCampaignListItem = {
  id: string;
  title: string;
  business_name: string;
  niche?: string;
  budget_summary?: string;
  creators_invited: number;
  status: string;
  created_at: string;
  contract_summary: string;
  payment_summary: string;
  contracts_active: number;
  contracts_pending: number;
};

export type AdminCampaignStats = {
  active: number;
  pending: number;
  flagged: number;
};

// --- Dispute Types ---
export type AdminDisputeListItem = {
  id: string;
  title: string;
  business_name: string;
  influencer_name: string;
  campaign_title: string;
  issue_type: string;
  issue_type_label: string;
  severity: string;
  status: string;
  preview: string;
  created_at: string;
};

export type AdminDisputeDetail = {
  id: string;
  title: string;
  business_name: string;
  influencer_name: string;
  influencer_email: string;
  influencer_user_id: string;
  campaign_title: string;
  campaign_id: string;
  contract_id: string;
  issue_type: string;
  issue_type_label: string;
  severity: string;
  status: string;
  description: string;
  preview: string;
  admin_note?: string;
  resolution_note?: string;
  created_at: string;
  updated_at: string;
  resolved_at?: string;
};

// --- Support Desk Types ---
export type AdminSupportTicketItem = {
  id: string;
  subject: string;
  message_preview: string;
  issue_type: string;
  issue_type_label: string;
  priority: string;
  status: string;
  sender_user_id: string;
  sender_name: string;
  sender_email: string;
  sender_role: string;
  related_campaign_title?: string;
  admin_reply?: string;
  created_at: string;
};

export type AdminSupportTicketStats = {
  open_count: number;
  high_priority_count: number;
  business_count: number;
  influencer_count: number;
  resolved_count: number;
};

export type AdminSupportTicketDetail = {
  id: string;
  subject: string;
  message: string;
  issue_type: string;
  issue_type_label: string;
  priority: string;
  status: string;
  sender_user_id: string;
  sender_name: string;
  sender_email: string;
  sender_role: string;
  related_campaign_id?: string;
  related_campaign_title?: string;
  admin_reply?: string;
  internal_note?: string;
  created_at: string;
  updated_at: string;
};

// --- API Service Methods ---
export const adminApi = {
  // 1. Dashboard Summary
  getDashboardSummary: async (): Promise<ApiResponse<AdminDashboardSummaryResponse>> => {
    const res = await apiClient.get<ApiResponse<AdminDashboardSummaryResponse>>(
      "/api/v1/admin/dashboard/summary"
    );
    return res.data;
  },

  // 2. Influencer Applications & Verifications
  listInfluencerApplications: async (params?: {
    status?: string;
    search?: string;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<{ items: InfluencerApplicationItem[]; pagination: { total: number; page: number; limit: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: InfluencerApplicationItem[]; pagination: { total: number; page: number; limit: number } }>>(
      "/api/v1/admin/influencers/applications",
      { params }
    );
    return res.data;
  },

  getInfluencerApplication: async (userId: string): Promise<ApiResponse<InfluencerApplicationDetail>> => {
    const res = await apiClient.get<ApiResponse<InfluencerApplicationDetail>>(
      `/api/v1/admin/influencers/applications/${userId}`
    );
    return res.data;
  },

  approveInfluencerApplication: async (
    userId: string,
    data?: { admin_note?: string }
  ): Promise<ApiResponse<{ user_id: string; status: string; approved_at: string }>> => {
    const res = await apiClient.post<ApiResponse<{ user_id: string; status: string; approved_at: string }>>(
      `/api/v1/admin/influencers/applications/${userId}/approve`,
      data
    );
    return res.data;
  },

  rejectInfluencerApplication: async (
    userId: string,
    data: { reason: string; admin_note?: string }
  ): Promise<ApiResponse<{ user_id: string; status: string; reason: string }>> => {
    const res = await apiClient.post<ApiResponse<{ user_id: string; status: string; reason: string }>>(
      `/api/v1/admin/influencers/applications/${userId}/reject`,
      data
    );
    return res.data;
  },

  updateInfluencerApplicationNotes: async (
    userId: string,
    data: { admin_note: string }
  ): Promise<ApiResponse<{ user_id: string; admin_note: string }>> => {
    const res = await apiClient.patch<ApiResponse<{ user_id: string; admin_note: string }>>(
      `/api/v1/admin/influencers/applications/${userId}/notes`,
      data
    );
    return res.data;
  },

  // 3. Business Verifications
  listBusinessVerifications: async (params?: {
    status?: string;
    search?: string;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<{ items: BusinessVerificationItem[]; pagination: { total: number; page: number; limit: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: BusinessVerificationItem[]; pagination: { total: number; page: number; limit: number } }>>(
      "/api/v1/admin/businesses/verifications",
      { params }
    );
    return res.data;
  },

  getBusinessVerification: async (userId: string): Promise<ApiResponse<BusinessVerificationDetail>> => {
    const res = await apiClient.get<ApiResponse<BusinessVerificationDetail>>(
      `/api/v1/admin/businesses/verifications/${userId}`
    );
    return res.data;
  },

  verifyBusiness: async (
    userId: string,
    data?: { admin_note?: string }
  ): Promise<ApiResponse<{ user_id: string; is_verified: boolean; verified_at: string }>> => {
    const res = await apiClient.post<ApiResponse<{ user_id: string; is_verified: boolean; verified_at: string }>>(
      `/api/v1/admin/businesses/verifications/${userId}/verify`,
      data
    );
    return res.data;
  },

  revokeBusinessVerification: async (
    userId: string,
    data?: { reason?: string }
  ): Promise<ApiResponse<{ user_id: string; is_verified: boolean }>> => {
    const res = await apiClient.post<ApiResponse<{ user_id: string; is_verified: boolean }>>(
      `/api/v1/admin/businesses/verifications/${userId}/revoke`,
      data
    );
    return res.data;
  },

  // 4. User Governance
  listUsers: async (params?: {
    role?: string;
    status?: string;
    search?: string;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<{ items: AdminUserListItem[]; stats: AdminUserStats; pagination: { total: number; page: number; limit: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: AdminUserListItem[]; stats: AdminUserStats; pagination: { total: number; page: number; limit: number } }>>(
      "/api/v1/admin/users",
      { params }
    );
    return res.data;
  },

  getUser: async (userId: string): Promise<ApiResponse<AdminUserListItem & { status_reason?: string }>> => {
    const res = await apiClient.get<ApiResponse<AdminUserListItem & { status_reason?: string }>>(
      `/api/v1/admin/users/${userId}`
    );
    return res.data;
  },

  blockUser: async (userId: string, data?: { reason?: string }): Promise<ApiResponse<{ user_id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ user_id: string; status: string }>>(
      `/api/v1/admin/users/${userId}/block`,
      data
    );
    return res.data;
  },

  unblockUser: async (userId: string): Promise<ApiResponse<{ user_id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ user_id: string; status: string }>>(
      `/api/v1/admin/users/${userId}/unblock`
    );
    return res.data;
  },

  removeUser: async (userId: string, data?: { reason?: string }): Promise<ApiResponse<{ user_id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ user_id: string; status: string }>>(
      `/api/v1/admin/users/${userId}/remove`,
      data
    );
    return res.data;
  },

  // 5. Categories
  listCategories: async (): Promise<ApiResponse<{ items: AdminCategoryItem[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: AdminCategoryItem[] }>>(
      "/api/v1/admin/categories"
    );
    return res.data;
  },

  createCategory: async (data: { name: string; description?: string }): Promise<ApiResponse<AdminCategoryItem>> => {
    const res = await apiClient.post<ApiResponse<AdminCategoryItem>>(
      "/api/v1/admin/categories",
      data
    );
    return res.data;
  },

  updateCategory: async (
    categoryId: string,
    data: { name?: string; description?: string }
  ): Promise<ApiResponse<AdminCategoryItem>> => {
    const res = await apiClient.patch<ApiResponse<AdminCategoryItem>>(
      `/api/v1/admin/categories/${categoryId}`,
      data
    );
    return res.data;
  },

  disableCategory: async (categoryId: string): Promise<ApiResponse<{ id: string; is_active: boolean }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; is_active: boolean }>>(
      `/api/v1/admin/categories/${categoryId}/disable`
    );
    return res.data;
  },

  enableCategory: async (categoryId: string): Promise<ApiResponse<{ id: string; is_active: boolean }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; is_active: boolean }>>(
      `/api/v1/admin/categories/${categoryId}/enable`
    );
    return res.data;
  },

  // 6. Campaign Supervision
  listCampaigns: async (params?: {
    status?: string;
    niche?: string;
    business?: string;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<{ items: AdminCampaignListItem[]; stats: AdminCampaignStats; pagination: { total: number; page: number; limit: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: AdminCampaignListItem[]; stats: AdminCampaignStats; pagination: { total: number; page: number; limit: number } }>>(
      "/api/v1/admin/campaigns",
      { params }
    );
    return res.data;
  },

  getCampaign: async (campaignId: string): Promise<ApiResponse<AdminCampaignListItem & { goal?: string; description: string }>> => {
    const res = await apiClient.get<ApiResponse<AdminCampaignListItem & { goal?: string; description: string }>>(
      `/api/v1/admin/campaigns/${campaignId}`
    );
    return res.data;
  },

  pauseCampaign: async (campaignId: string): Promise<ApiResponse<{ id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; status: string }>>(
      `/api/v1/admin/campaigns/${campaignId}/pause`
    );
    return res.data;
  },

  completeCampaign: async (campaignId: string): Promise<ApiResponse<{ id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; status: string }>>(
      `/api/v1/admin/campaigns/${campaignId}/complete`
    );
    return res.data;
  },

  flagCampaign: async (campaignId: string): Promise<ApiResponse<{ id: string; status: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; status: string }>>(
      `/api/v1/admin/campaigns/${campaignId}/flag`
    );
    return res.data;
  },

  // 7. Dispute Resolution
  listDisputes: async (params?: {
    tab?: string;
    search?: string;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<{ items: AdminDisputeListItem[]; pagination: { total: number; page: number; limit: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: AdminDisputeListItem[]; pagination: { total: number; page: number; limit: number } }>>(
      "/api/v1/admin/disputes",
      { params }
    );
    return res.data;
  },

  getDispute: async (disputeId: string): Promise<ApiResponse<AdminDisputeDetail>> => {
    const res = await apiClient.get<ApiResponse<AdminDisputeDetail>>(
      `/api/v1/admin/disputes/${disputeId}`
    );
    return res.data;
  },

  updateDisputeStatus: async (
    disputeId: string,
    data: { status: string }
  ): Promise<ApiResponse<{ id: string; status: string }>> => {
    const res = await apiClient.patch<ApiResponse<{ id: string; status: string }>>(
      `/api/v1/admin/disputes/${disputeId}/status`,
      data
    );
    return res.data;
  },

  updateDisputeNotes: async (
    disputeId: string,
    data: { admin_note?: string; resolution_note?: string }
  ): Promise<ApiResponse<{ id: string; admin_note: string; resolution_note?: string }>> => {
    const res = await apiClient.patch<ApiResponse<{ id: string; admin_note: string; resolution_note?: string }>>(
      `/api/v1/admin/disputes/${disputeId}/notes`,
      data
    );
    return res.data;
  },

  // 8. Support Desk
  listSupportTickets: async (params?: {
    role?: string;
    status?: string;
    priority?: string;
    search?: string;
    page?: number;
    limit?: number;
  }): Promise<ApiResponse<{ items: AdminSupportTicketItem[]; stats: AdminSupportTicketStats; pagination: { total: number; page: number; limit: number } }>> => {
    const res = await apiClient.get<ApiResponse<{ items: AdminSupportTicketItem[]; stats: AdminSupportTicketStats; pagination: { total: number; page: number; limit: number } }>>(
      "/api/v1/admin/support/tickets",
      { params }
    );
    return res.data;
  },

  getSupportTicket: async (ticketId: string): Promise<ApiResponse<AdminSupportTicketDetail>> => {
    const res = await apiClient.get<ApiResponse<AdminSupportTicketDetail>>(
      `/api/v1/admin/support/tickets/${ticketId}`
    );
    return res.data;
  },

  updateSupportTicketStatus: async (
    ticketId: string,
    data: { status: string }
  ): Promise<ApiResponse<{ id: string; status: string }>> => {
    const res = await apiClient.patch<ApiResponse<{ id: string; status: string }>>(
      `/api/v1/admin/support/tickets/${ticketId}/status`,
      data
    );
    return res.data;
  },

  replySupportTicket: async (
    ticketId: string,
    data: { admin_reply?: string; internal_note?: string }
  ): Promise<ApiResponse<{ id: string; admin_reply?: string }>> => {
    const res = await apiClient.post<ApiResponse<{ id: string; admin_reply?: string }>>(
      `/api/v1/admin/support/tickets/${ticketId}/reply`,
      data
    );
    return res.data;
  },
};
