import { apiClient } from "@/lib/api/api-client";

export type UserRole = "business" | "influencer" | "admin";
export type UserStatus =
  | "pending_profile"
  | "pending_approval"
  | "active"
  | "rejected"
  | "suspended"
  | "blocked"
  | "removed";

export type AuthUser = {
  id: string;
  email: string;
  phone_number?: string;
  role: UserRole;
  status: UserStatus;
  email_verified: boolean;
  phone_verified: boolean;
  status_reason?: string | null;
  suspended_until?: string | null;
  last_login_at?: string | null;
  created_at?: string;
};

export type AuthProfile = {
  id: string;
  type?: string;
  business_name?: string;
  owner_full_name?: string;
  full_name?: string;
  creator_name?: string;
  username?: string;
  is_verified?: boolean;
  permission_level?: string;
};

export type AuthTokens = {
  access_token: string;
  refresh_token: string;
  expires_in: number;
};

export type SignupRequest = {
  role: "business" | "influencer";
  name: string;
  email: string;
  phone_number: string;
  password: string;
  confirm_password: string;
};

export type SignupResponseData = {
  user: AuthUser;
  profile: AuthProfile;
  tokens: AuthTokens;
  redirect_to: string;
};

export type LoginRequest = {
  email: string;
  password: string;
  remember_me?: boolean;
  portal?: string;
  next?: string;
};

export type LoginResponseData = {
  user: AuthUser;
  profile: AuthProfile;
  tokens: AuthTokens;
  redirect_to: string;
};

export type MeResponseData = {
  user: AuthUser;
  profile: AuthProfile;
};

export type ApiResponse<T> = {
  status: number;
  message: string;
  data?: T;
  error?: {
    type: string;
    message: string;
    details?: Array<{ title: string; description: string }>;
  };
};

export const authApi = {
  signup: async (data: SignupRequest): Promise<ApiResponse<SignupResponseData>> => {
    const res = await apiClient.post<ApiResponse<SignupResponseData>>("/api/v1/auth/signup", data);
    return res.data;
  },

  login: async (data: LoginRequest): Promise<ApiResponse<LoginResponseData>> => {
    const res = await apiClient.post<ApiResponse<LoginResponseData>>("/api/v1/auth/login", data);
    return res.data;
  },

  refreshToken: async (refreshToken: string): Promise<ApiResponse<AuthTokens>> => {
    const res = await apiClient.post<ApiResponse<AuthTokens>>("/api/v1/auth/refresh", {
      refresh_token: refreshToken,
    });
    return res.data;
  },

  forgotPassword: async (email: string): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.post<ApiResponse<{ message: string }>>("/api/v1/auth/forgot-password", {
      email,
    });
    return res.data;
  },

  resetPassword: async (password: string, token: string): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.post<ApiResponse<{ message: string }>>("/api/v1/auth/reset-password", {
      password,
      token,
    });
    return res.data;
  },

  getMe: async (): Promise<ApiResponse<MeResponseData>> => {
    const res = await apiClient.get<ApiResponse<MeResponseData>>("/api/v1/auth/me");
    return res.data;
  },

  logout: async (refreshToken?: string): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.post<ApiResponse<{ message: string }>>("/api/v1/auth/logout", {
      refresh_token: refreshToken,
    });
    return res.data;
  },

  sendVerifyEmail: async (): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.post<ApiResponse<{ message: string }>>("/api/v1/auth/send-verify-email");
    return res.data;
  },

  verifyEmail: async (code: string): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.post<ApiResponse<{ message: string }>>("/api/v1/auth/verify-email", {
      code,
    });
    return res.data;
  },

  sendVerifyPhone: async (): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.post<ApiResponse<{ message: string }>>("/api/v1/auth/send-verify-phone");
    return res.data;
  },

  verifyPhone: async (code: string): Promise<ApiResponse<{ message: string }>> => {
    const res = await apiClient.post<ApiResponse<{ message: string }>>("/api/v1/auth/verify-phone", {
      code,
    });
    return res.data;
  },
};
