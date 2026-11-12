--
-- PostgreSQL database dump
--

\restrict F09d1iD4rin60UlBEf8T1cluR9yfUxtwjSl0g4eaGeJthmHdbcgBZpuoxBNqKYf

-- Dumped from database version 17.11 (Debian 17.11-1.pgdg13+2)
-- Dumped by pg_dump version 17.11 (Debian 17.11-1.pgdg13+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: audit_action; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.audit_action AS ENUM (
    'login',
    'logout',
    'signup',
    'password_reset',
    'profile_update',
    'status_change',
    'approve',
    'reject',
    'block',
    'unblock',
    'suspend',
    'unsuspend',
    'delete',
    'contract_status_change',
    'payment_status_change',
    'dispute_resolved',
    'permission_update',
    'admin_action'
);


ALTER TYPE public.audit_action OWNER TO adfluence;

--
-- Name: campaign_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.campaign_status AS ENUM (
    'draft',
    'pending',
    'active',
    'completed',
    'cancelled',
    'paused'
);


ALTER TYPE public.campaign_status OWNER TO adfluence;

--
-- Name: contract_deliverable_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.contract_deliverable_status AS ENUM (
    'pending',
    'submitted',
    'approved',
    'revision_requested',
    'completed'
);


ALTER TYPE public.contract_deliverable_status OWNER TO adfluence;

--
-- Name: contract_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.contract_status AS ENUM (
    'draft',
    'pending_influencer',
    'pending_business',
    'accepted',
    'active',
    'completed',
    'cancelled',
    'disputed',
    'rejected'
);


ALTER TYPE public.contract_status OWNER TO adfluence;

--
-- Name: conversation_participant_role; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.conversation_participant_role AS ENUM (
    'business',
    'influencer'
);


ALTER TYPE public.conversation_participant_role OWNER TO adfluence;

--
-- Name: conversation_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.conversation_status AS ENUM (
    'active',
    'archived',
    'closed'
);


ALTER TYPE public.conversation_status OWNER TO adfluence;

--
-- Name: creator_size; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.creator_size AS ENUM (
    'nano',
    'micro',
    'mid',
    'macro',
    'mega',
    'any'
);


ALTER TYPE public.creator_size OWNER TO adfluence;

--
-- Name: deliverable_type; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.deliverable_type AS ENUM (
    'tiktok_video',
    'instagram_reel',
    'instagram_story',
    'youtube_short',
    'product_review',
    'post',
    'other'
);


ALTER TYPE public.deliverable_type OWNER TO adfluence;

--
-- Name: dispute_severity; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.dispute_severity AS ENUM (
    'low',
    'medium',
    'high'
);


ALTER TYPE public.dispute_severity OWNER TO adfluence;

--
-- Name: dispute_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.dispute_status AS ENUM (
    'open',
    'under_review',
    'resolved',
    'escalated',
    'closed'
);


ALTER TYPE public.dispute_status OWNER TO adfluence;

--
-- Name: document_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.document_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);


ALTER TYPE public.document_status OWNER TO adfluence;

--
-- Name: document_type; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.document_type AS ENUM (
    'fayda_id',
    'passport',
    'business_license',
    'trade_license',
    'tax_certificate',
    'portfolio_sample',
    'other'
);


ALTER TYPE public.document_type OWNER TO adfluence;

--
-- Name: invite_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.invite_status AS ENUM (
    'pending',
    'accepted',
    'rejected',
    'cancelled'
);


ALTER TYPE public.invite_status OWNER TO adfluence;

--
-- Name: message_report_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.message_report_status AS ENUM (
    'open',
    'under_review',
    'resolved',
    'dismissed'
);


ALTER TYPE public.message_report_status OWNER TO adfluence;

--
-- Name: notification_type; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.notification_type AS ENUM (
    'campaign_invite',
    'invite_accepted',
    'invite_rejected',
    'contract_created',
    'contract_accepted',
    'contract_completed',
    'payment_paid',
    'payout_paid',
    'support_reply',
    'verification_approved',
    'verification_rejected',
    'profile_approved',
    'profile_rejected',
    'dispute_opened',
    'dispute_resolved',
    'message_report',
    'general'
);


ALTER TYPE public.notification_type OWNER TO adfluence;

--
-- Name: payment_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.payment_status AS ENUM (
    'pending',
    'paid',
    'failed',
    'refunded',
    'cancelled'
);


ALTER TYPE public.payment_status OWNER TO adfluence;

--
-- Name: payout_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.payout_status AS ENUM (
    'pending',
    'processing',
    'paid',
    'failed',
    'cancelled'
);


ALTER TYPE public.payout_status OWNER TO adfluence;

--
-- Name: recommendation_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.recommendation_status AS ENUM (
    'recommended',
    'invited',
    'dismissed'
);


ALTER TYPE public.recommendation_status OWNER TO adfluence;

--
-- Name: social_platform; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.social_platform AS ENUM (
    'tiktok',
    'instagram',
    'youtube',
    'telegram',
    'facebook'
);


ALTER TYPE public.social_platform OWNER TO adfluence;

--
-- Name: support_issue_type; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.support_issue_type AS ENUM (
    'campaign_issue',
    'contract_issue',
    'payment_question',
    'creator_communication',
    'profile_approval',
    'technical_issue',
    'general_feedback',
    'complaint'
);


ALTER TYPE public.support_issue_type OWNER TO adfluence;

--
-- Name: support_priority; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.support_priority AS ENUM (
    'low',
    'normal',
    'high',
    'urgent'
);


ALTER TYPE public.support_priority OWNER TO adfluence;

--
-- Name: support_ticket_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.support_ticket_status AS ENUM (
    'open',
    'in_review',
    'resolved',
    'escalated',
    'closed'
);


ALTER TYPE public.support_ticket_status OWNER TO adfluence;

--
-- Name: user_role; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.user_role AS ENUM (
    'business',
    'influencer',
    'admin'
);


ALTER TYPE public.user_role OWNER TO adfluence;

--
-- Name: user_status; Type: TYPE; Schema: public; Owner: adfluence
--

CREATE TYPE public.user_status AS ENUM (
    'pending_profile',
    'pending_approval',
    'active',
    'rejected',
    'suspended',
    'blocked',
    'removed'
);


ALTER TYPE public.user_status OWNER TO adfluence;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: admin_profiles; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.admin_profiles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    full_name character varying(255) NOT NULL,
    permission_level character varying(50) DEFAULT 'super_admin'::character varying NOT NULL,
    permissions jsonb DEFAULT '{"manage_users": true, "view_reports": true, "manage_support": true, "manage_disputes": true, "manage_payments": true, "manage_campaigns": true, "manage_categories": true}'::jsonb NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.admin_profiles OWNER TO adfluence;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.audit_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    actor_user_id uuid,
    target_user_id uuid,
    action public.audit_action NOT NULL,
    entity_type character varying(100),
    entity_id uuid,
    description text,
    metadata jsonb,
    ip_address inet,
    user_agent text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.audit_logs OWNER TO adfluence;

--
-- Name: auth_login_lockouts; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.auth_login_lockouts (
    email character varying(255) NOT NULL,
    failed_attempts integer DEFAULT 0 NOT NULL,
    lockout_tier integer DEFAULT 0 NOT NULL,
    locked_until timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.auth_login_lockouts OWNER TO adfluence;

--
-- Name: auth_rate_limit_events; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.auth_rate_limit_events (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    scope character varying(50) NOT NULL,
    rate_key character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.auth_rate_limit_events OWNER TO adfluence;

--
-- Name: bank_accounts; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.bank_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    bank_name character varying(150),
    account_holder_name character varying(255),
    account_number character varying(100),
    mobile_money_provider character varying(100),
    mobile_money_number character varying(50),
    is_primary boolean DEFAULT true NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.bank_accounts OWNER TO adfluence;

--
-- Name: business_profiles; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.business_profiles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    owner_full_name character varying(255) NOT NULL,
    business_name character varying(255) NOT NULL,
    legal_business_name character varying(255),
    logo_url text,
    description text,
    industry character varying(100),
    business_type character varying(100),
    business_size character varying(100),
    city character varying(100),
    region character varying(100),
    address text,
    tin_number character varying(100),
    website_url text,
    instagram_url text,
    tiktok_url text,
    facebook_url text,
    linkedin_url text,
    preferred_platforms text[],
    monthly_marketing_budget numeric(14,2),
    is_verified boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.business_profiles OWNER TO adfluence;

--
-- Name: campaign_deliverables; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.campaign_deliverables (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    campaign_id uuid NOT NULL,
    deliverable_type public.deliverable_type NOT NULL,
    quantity integer DEFAULT 1 NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.campaign_deliverables OWNER TO adfluence;

--
-- Name: campaign_invites; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.campaign_invites (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    campaign_id uuid NOT NULL,
    business_user_id uuid NOT NULL,
    influencer_user_id uuid NOT NULL,
    message text,
    status public.invite_status DEFAULT 'pending'::public.invite_status NOT NULL,
    responded_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.campaign_invites OWNER TO adfluence;

--
-- Name: campaign_platforms; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.campaign_platforms (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    campaign_id uuid NOT NULL,
    platform public.social_platform NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.campaign_platforms OWNER TO adfluence;

--
-- Name: campaign_recommendations; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.campaign_recommendations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    campaign_id uuid NOT NULL,
    influencer_user_id uuid NOT NULL,
    match_score numeric(5,2) DEFAULT 0 NOT NULL,
    score_breakdown jsonb DEFAULT '{}'::jsonb NOT NULL,
    rank_position integer,
    status public.recommendation_status DEFAULT 'recommended'::public.recommendation_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.campaign_recommendations OWNER TO adfluence;

--
-- Name: campaigns; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.campaigns (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    business_user_id uuid NOT NULL,
    title character varying(255) NOT NULL,
    goal text,
    description text NOT NULL,
    category_id uuid,
    target_audience text,
    city character varying(100),
    region character varying(100),
    budget_min_etb numeric(14,2),
    budget_max_etb numeric(14,2),
    creator_size_preference public.creator_size DEFAULT 'any'::public.creator_size NOT NULL,
    number_of_creators_needed integer DEFAULT 1 NOT NULL,
    start_date date,
    end_date date,
    status public.campaign_status DEFAULT 'draft'::public.campaign_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.campaigns OWNER TO adfluence;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(100) NOT NULL,
    slug character varying(120) NOT NULL,
    description text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.categories OWNER TO adfluence;

--
-- Name: chat_admin_access_logs; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.chat_admin_access_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    admin_user_id uuid NOT NULL,
    conversation_id uuid NOT NULL,
    reason text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.chat_admin_access_logs OWNER TO adfluence;

--
-- Name: contract_deliverables; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.contract_deliverables (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    contract_id uuid NOT NULL,
    deliverable_type public.deliverable_type NOT NULL,
    quantity integer DEFAULT 1 NOT NULL,
    description text,
    status public.contract_deliverable_status DEFAULT 'pending'::public.contract_deliverable_status NOT NULL,
    due_date date,
    completed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.contract_deliverables OWNER TO adfluence;

--
-- Name: contract_status_history; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.contract_status_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    contract_id uuid NOT NULL,
    old_status public.contract_status,
    new_status public.contract_status NOT NULL,
    changed_by uuid,
    reason text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.contract_status_history OWNER TO adfluence;

--
-- Name: contracts; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.contracts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    campaign_id uuid NOT NULL,
    business_user_id uuid NOT NULL,
    influencer_user_id uuid NOT NULL,
    conversation_id uuid,
    invite_id uuid,
    title character varying(255) NOT NULL,
    description text,
    amount_etb numeric(14,2) NOT NULL,
    business_fee_percent numeric(5,2) DEFAULT 5.00 NOT NULL,
    influencer_fee_percent numeric(5,2) DEFAULT 5.00 NOT NULL,
    business_fee_amount_etb numeric(14,2) NOT NULL,
    influencer_fee_amount_etb numeric(14,2) NOT NULL,
    total_business_payable_etb numeric(14,2) NOT NULL,
    influencer_receivable_etb numeric(14,2) NOT NULL,
    platform_total_fee_etb numeric(14,2) NOT NULL,
    start_date date,
    end_date date,
    status public.contract_status DEFAULT 'draft'::public.contract_status NOT NULL,
    accepted_by_business_at timestamp with time zone,
    accepted_by_influencer_at timestamp with time zone,
    completed_at timestamp with time zone,
    cancelled_at timestamp with time zone,
    cancellation_reason text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.contracts OWNER TO adfluence;

--
-- Name: conversation_participants; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.conversation_participants (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    conversation_id uuid NOT NULL,
    user_id uuid NOT NULL,
    role public.conversation_participant_role NOT NULL,
    last_read_at timestamp with time zone,
    is_muted boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.conversation_participants OWNER TO adfluence;

--
-- Name: conversations; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.conversations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    campaign_id uuid,
    status public.conversation_status DEFAULT 'active'::public.conversation_status NOT NULL,
    last_message_preview text,
    last_message_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.conversations OWNER TO adfluence;

--
-- Name: disputes; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.disputes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    contract_id uuid NOT NULL,
    opened_by_user_id uuid NOT NULL,
    against_user_id uuid NOT NULL,
    assigned_admin_user_id uuid,
    issue_type public.support_issue_type NOT NULL,
    description text NOT NULL,
    severity public.dispute_severity DEFAULT 'medium'::public.dispute_severity NOT NULL,
    status public.dispute_status DEFAULT 'open'::public.dispute_status NOT NULL,
    resolution_note text,
    resolved_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    resolved_at timestamp with time zone
);


ALTER TABLE public.disputes OWNER TO adfluence;

--
-- Name: email_verification_tokens; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.email_verification_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    used_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.email_verification_tokens OWNER TO adfluence;

--
-- Name: influencer_portfolio_items; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.influencer_portfolio_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    file_url text NOT NULL,
    object_key text NOT NULL,
    media_type character varying(10) NOT NULL,
    mime_type character varying(100) NOT NULL,
    file_size_bytes bigint NOT NULL,
    sort_order smallint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT influencer_portfolio_items_file_size_bytes_check CHECK ((file_size_bytes > 0)),
    CONSTRAINT influencer_portfolio_items_media_type_check CHECK (((media_type)::text = ANY ((ARRAY['image'::character varying, 'video'::character varying])::text[]))),
    CONSTRAINT influencer_portfolio_items_sort_order_check CHECK (((sort_order >= 0) AND (sort_order <= 2)))
);


ALTER TABLE public.influencer_portfolio_items OWNER TO adfluence;

--
-- Name: influencer_profile_categories; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.influencer_profile_categories (
    influencer_user_id uuid NOT NULL,
    category_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.influencer_profile_categories OWNER TO adfluence;

--
-- Name: influencer_profile_platforms; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.influencer_profile_platforms (
    influencer_user_id uuid NOT NULL,
    platform public.social_platform NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.influencer_profile_platforms OWNER TO adfluence;

--
-- Name: influencer_profiles; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.influencer_profiles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    full_name character varying(255) NOT NULL,
    creator_name character varying(255) NOT NULL,
    username character varying(100),
    profile_image_url text,
    bio text,
    city character varying(100),
    region character varying(100),
    languages text[],
    social_links jsonb DEFAULT '{}'::jsonb NOT NULL,
    instagram_followers integer DEFAULT 0 NOT NULL,
    tiktok_followers integer DEFAULT 0 NOT NULL,
    youtube_subscribers integer DEFAULT 0 NOT NULL,
    total_followers integer DEFAULT 0 NOT NULL,
    engagement_rate numeric(5,2),
    audience_description text,
    content_style text,
    portfolio_description text,
    min_price_etb numeric(12,2),
    max_price_etb numeric(12,2),
    approval_note text,
    is_verified boolean DEFAULT false NOT NULL,
    approved_at timestamp with time zone,
    approved_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.influencer_profiles OWNER TO adfluence;

--
-- Name: message_reports; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.message_reports (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    conversation_id uuid NOT NULL,
    firebase_message_id character varying(255) NOT NULL,
    reporter_user_id uuid NOT NULL,
    reported_user_id uuid NOT NULL,
    reason text NOT NULL,
    status public.message_report_status DEFAULT 'open'::public.message_report_status NOT NULL,
    resolved_by uuid,
    resolution_note text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    resolved_at timestamp with time zone
);


ALTER TABLE public.message_reports OWNER TO adfluence;

--
-- Name: messages; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    conversation_id uuid NOT NULL,
    sender_user_id uuid NOT NULL,
    receiver_user_id uuid NOT NULL,
    body text NOT NULL,
    media_url text,
    media_type character varying(50) DEFAULT 'text'::character varying,
    is_read boolean DEFAULT false NOT NULL,
    read_at timestamp with time zone,
    deleted_for_sender boolean DEFAULT false NOT NULL,
    deleted_for_receiver boolean DEFAULT false NOT NULL,
    is_deleted_for_everyone boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.messages OWNER TO adfluence;

--
-- Name: notifications; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.notifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    type public.notification_type NOT NULL,
    entity_type character varying(100),
    entity_id uuid,
    read_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.notifications OWNER TO adfluence;

--
-- Name: password_reset_tokens; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.password_reset_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    used_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.password_reset_tokens OWNER TO adfluence;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.payments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    contract_id uuid NOT NULL,
    business_user_id uuid NOT NULL,
    influencer_user_id uuid NOT NULL,
    amount_etb numeric(14,2) NOT NULL,
    business_fee_etb numeric(14,2) NOT NULL,
    influencer_fee_etb numeric(14,2) NOT NULL,
    platform_fee_total_etb numeric(14,2) NOT NULL,
    total_paid_by_business_etb numeric(14,2) NOT NULL,
    influencer_payout_etb numeric(14,2) NOT NULL,
    payment_method character varying(50) NOT NULL,
    provider character varying(100),
    provider_reference character varying(255),
    status public.payment_status DEFAULT 'pending'::public.payment_status NOT NULL,
    paid_at timestamp with time zone,
    failed_at timestamp with time zone,
    failure_reason text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.payments OWNER TO adfluence;

--
-- Name: payouts; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.payouts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    payment_id uuid NOT NULL,
    influencer_user_id uuid NOT NULL,
    bank_account_id uuid,
    amount_etb numeric(14,2) NOT NULL,
    status public.payout_status DEFAULT 'pending'::public.payout_status NOT NULL,
    provider character varying(100),
    provider_reference character varying(255),
    paid_at timestamp with time zone,
    failed_at timestamp with time zone,
    failure_reason text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.payouts OWNER TO adfluence;

--
-- Name: phone_verification_tokens; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.phone_verification_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    used_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.phone_verification_tokens OWNER TO adfluence;

--
-- Name: reviews; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.reviews (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    contract_id uuid NOT NULL,
    reviewer_user_id uuid NOT NULL,
    reviewed_user_id uuid NOT NULL,
    rating integer NOT NULL,
    comment text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT reviews_rating_check CHECK (((rating >= 1) AND (rating <= 5)))
);


ALTER TABLE public.reviews OWNER TO adfluence;

--
-- Name: support_messages; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.support_messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ticket_id uuid NOT NULL,
    sender_user_id uuid NOT NULL,
    message text NOT NULL,
    attachment_url text,
    is_internal_note boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.support_messages OWNER TO adfluence;

--
-- Name: support_tickets; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.support_tickets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_by_user_id uuid NOT NULL,
    assigned_admin_user_id uuid,
    related_campaign_id uuid,
    related_contract_id uuid,
    subject character varying(255) NOT NULL,
    message text NOT NULL,
    issue_type public.support_issue_type NOT NULL,
    priority public.support_priority DEFAULT 'normal'::public.support_priority NOT NULL,
    status public.support_ticket_status DEFAULT 'open'::public.support_ticket_status NOT NULL,
    resolution_note text,
    resolved_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    resolved_at timestamp with time zone
);


ALTER TABLE public.support_tickets OWNER TO adfluence;

--
-- Name: user_sessions; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.user_sessions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    refresh_token_hash text NOT NULL,
    ip_address inet,
    user_agent text,
    expires_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    refresh_ttl_seconds integer,
    rotated_at timestamp with time zone
);


ALTER TABLE public.user_sessions OWNER TO adfluence;

--
-- Name: users; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    phone_number character varying(50),
    password_hash text NOT NULL,
    role public.user_role NOT NULL,
    status public.user_status DEFAULT 'pending_profile'::public.user_status NOT NULL,
    email_verified boolean DEFAULT false NOT NULL,
    phone_verified boolean DEFAULT false NOT NULL,
    last_login_at timestamp with time zone,
    suspended_until timestamp with time zone,
    status_reason text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO adfluence;

--
-- Name: verification_documents; Type: TABLE; Schema: public; Owner: adfluence
--

CREATE TABLE public.verification_documents (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    document_type public.document_type NOT NULL,
    file_url text NOT NULL,
    file_name character varying(255),
    mime_type character varying(100),
    file_size_bytes bigint,
    status public.document_status DEFAULT 'pending'::public.document_status NOT NULL,
    rejection_reason text,
    reviewed_by uuid,
    reviewed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.verification_documents OWNER TO adfluence;

--
-- Data for Name: admin_profiles; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.admin_profiles (id, user_id, full_name, permission_level, permissions, created_at, updated_at) FROM stdin;
a2000000-0000-4000-8000-000000000001	a1000000-0000-4000-8000-000000000001	Adfluence Admin	super_admin	{"manage_users": true, "view_reports": true, "manage_support": true, "manage_disputes": true, "manage_payments": true, "manage_campaigns": true, "manage_categories": true}	2026-10-20 07:26:27.132952+00	2026-10-20 07:26:27.132952+00
\.


--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.audit_logs (id, actor_user_id, target_user_id, action, entity_type, entity_id, description, metadata, ip_address, user_agent, created_at) FROM stdin;
f5099ce0-4850-431b-864c-2f03c4502028	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	signup	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user signed up	\N	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 07:55:56.295599+00
11933fb8-cbe9-43a6-8b6e-c1ae6cc0d987	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-20 07:58:42.385678+00
4d1310e3-d8f6-41d9-ae1e-36f183bf9ae1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-20 07:58:52.77238+00
faa2d9e1-22a8-443e-a19a-88aac5cc034d	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-20 07:58:57.82543+00
dad59b4b-94d9-496c-af5a-6bcaecd9be24	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-20 08:13:44.831194+00
e1294f47-039f-4907-926a-e5e853f8162d	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-20 08:13:47.503186+00
b19c0888-f2e8-4618-a337-b4e3c24061b8	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-20 08:13:50.147756+00
29876f11-5903-49e2-9d36-0f7ab17d0e75	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-20 08:32:47.615538+00
376c59a6-6569-4ef4-ba20-68616efbb7bb	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_profile	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	influencer avatar uploaded	\N	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 08:49:57.113874+00
71436ac4-1932-44de-abc9-b2b30f43ca9f	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_profile	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	influencer avatar deleted	\N	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 08:50:11.632715+00
142b75b8-4aee-49a0-9adf-d5a657e8310e	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_profile	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	influencer avatar uploaded	\N	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 08:50:29.555489+00
71089153-bd80-42ad-aa84-d230c8e2cd05	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_profile	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	influencer avatar deleted	\N	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 08:50:39.83745+00
505d8e22-683e-4807-96e9-bc0081f2add7	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 08:51:46.139265+00
7dbe9650-55b4-4b17-a1d6-08c3c031d577	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 08:51:55.921334+00
8706b357-7dc0-4a7f-9726-db1a543511d4	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Mobile Safari/537.36	2026-10-20 08:54:41.93305+00
270b03d8-07b1-4a21-b2a7-480b8fb238f5	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	102.213.71.203	PostmanRuntime/2.12.0	2026-10-20 13:39:20.77948+00
4569815f-4941-4086-989d-195d4e02a68f	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:27:22.095847+00
2b7510bc-c235-46d4-894e-33232726bcb8	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_profile	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	influencer avatar uploaded	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:27:32.118851+00
ecad4c31-395b-4255-8d7e-ad0bccd6e20f	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_profile	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	influencer profile updated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:27:40.610367+00
f37a19c1-64d3-4416-8477-c16a2ea54f96	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_portfolio_item	4502c1d7-1c62-4972-be68-5531f7ca75cc	influencer portfolio sample uploaded	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:27:56.806228+00
48db26d9-07e3-432e-b42f-491a90dc4905	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	profile_update	influencer_portfolio_item	4502c1d7-1c62-4972-be68-5531f7ca75cc	influencer portfolio sample deleted	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:28:05.398381+00
7c0bd89d-e561-4533-bef8-6be605f865f8	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	signup	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	user signed up	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:29:24.409315+00
13fd4b0a-1034-4d0d-94cf-a8ec2f53d69e	a1000000-0000-4000-8000-000000000001	a1000000-0000-4000-8000-000000000001	login	user	a1000000-0000-4000-8000-000000000001	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:30:07.429867+00
2dc34f46-ba47-4272-9d52-8ff0e56ae374	a1000000-0000-4000-8000-000000000001	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	approve	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	influencer application approved	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:30:39.208854+00
3fb59e69-a33e-43a0-9fd1-7e6aa6eaad8b	a1000000-0000-4000-8000-000000000001	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	reject	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	influencer application rejected	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:30:57.553841+00
d43eee4e-ab3a-453c-a227-7a5e3a19b0d9	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	login	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:31:31.274707+00
c6fa2b5a-06c1-49d8-a4a8-7354ac79c400	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	profile_update	influencer_profile	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	influencer avatar uploaded	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:31:42.314466+00
1109c97d-6330-43b7-ba88-c61eef3607f1	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	profile_update	influencer_profile	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	influencer profile updated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:31:54.668459+00
fec58e07-7967-4ed4-9bea-cc03c5ade2eb	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	signup	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user signed up	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:33:27.399538+00
304a7e51-1d8f-4ddd-b7a1-f9dde67763c4	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	profile_update	campaign	2499c495-01d1-4824-b064-dfd70a2755b6	business campaign created	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:35:33.034326+00
08b3d2a5-47bc-4968-9e02-f3043de34d8d	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	profile_update	campaign_recommendations	2499c495-01d1-4824-b064-dfd70a2755b6	campaign recommendations generated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:35:33.516899+00
d5e106e3-0cac-472d-acd2-19e631ab6519	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	profile_update	campaign_recommendations	2499c495-01d1-4824-b064-dfd70a2755b6	campaign recommendations generated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:35:34.888982+00
13102db8-3554-46a2-9dba-51c3c46ecee1	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	profile_update	campaign_recommendations	2499c495-01d1-4824-b064-dfd70a2755b6	campaign recommendations generated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:35:40.998593+00
6fe29b93-2277-4048-b158-39877bde9b8a	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	profile_update	campaign_recommendations	2499c495-01d1-4824-b064-dfd70a2755b6	campaign recommendations generated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:35:52.848491+00
5036eae4-0b70-4a6e-99a6-cd03442aad0a	a1000000-0000-4000-8000-000000000001	a1000000-0000-4000-8000-000000000001	login	user	a1000000-0000-4000-8000-000000000001	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:36:15.044045+00
22212fd6-3082-4059-a919-08d73f462437	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	profile_update	influencer_profile	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	influencer profile updated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:38:15.941019+00
48301ac7-f323-49f5-bbb1-64594a1745bf	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	status_change	campaign_invite	3690599e-cfda-4467-a1bd-1b229e6da28d	business sent campaign invite	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:39:07.634829+00
f7190f38-617b-4401-8bc6-fabf8f0f8426	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	login	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:39:55.420367+00
05a87f2f-8173-44ac-8898-e7b5d04b4139	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	status_change	campaign_invite	3690599e-cfda-4467-a1bd-1b229e6da28d	influencer accepted campaign invite	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:40:05.592071+00
49247f09-5128-468e-8bec-1416a7046de3	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	login	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:46:03.426761+00
ebce856b-124e-4d49-b3e6-532d0081bdd1	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	login	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:46:08.66948+00
ef4f587d-76e3-45cf-9c7f-6ff44ea449bb	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	login	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:46:43.494406+00
5d2f9de2-e541-43dc-9c82-54f047585dfb	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	login	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:46:49.204975+00
b71f8861-edcd-47cf-a323-dd08925f7b7d	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	login	user	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	user logged in	\N	196.188.179.241	PostmanRuntime/2.12.0	2026-10-21 06:50:12.871141+00
66871a41-188a-4cc0-8471-2e4450d76176	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	login	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged in	\N	102.213.71.192	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-01 09:06:07.376819+00
4b84c942-978d-41b2-93b0-16c49dd9b13f	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	login	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged in	\N	102.213.71.192	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-01 09:07:53.933904+00
53cbe7c6-5b86-42a6-89ef-42d30e7ff233	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	login	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:37:36.632122+00
11e60cbc-5f71-455f-b816-1031f75ccf16	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	login	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:38:32.240838+00
3f3f0ea4-636e-436f-8cd7-35b96467deee	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	profile_update	campaign_recommendations	2499c495-01d1-4824-b064-dfd70a2755b6	campaign recommendations generated	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:38:45.441102+00
ced7a907-8f6b-472e-b5fa-993ab2e677e1	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	login	user	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:44:28.467297+00
7925e046-752b-481b-a50f-25ffc76bc94e	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	login	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged in	\N	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-20 17:45:15.242641+00
91ddd60c-823f-4395-bd99-bbf9a9192ee3	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	login	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged in	\N	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-01 11:59:18.497533+00
54935ad2-cb98-445b-a4b9-e67414b6f2c4	1429fd1b-6412-4e13-ba93-1887a7a74872	1429fd1b-6412-4e13-ba93-1887a7a74872	logout	user	1429fd1b-6412-4e13-ba93-1887a7a74872	user logged out	\N	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-01 12:31:03.663629+00
2069d85e-0f7c-4c05-b508-009c3752e049	a1000000-0000-4000-8000-000000000001	a1000000-0000-4000-8000-000000000001	login	user	a1000000-0000-4000-8000-000000000001	user logged in	\N	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-01 12:31:11.599513+00
6b1a3cf4-df18-43cf-8ed2-2a9d2d1062b5	a1000000-0000-4000-8000-000000000001	a1000000-0000-4000-8000-000000000001	logout	user	a1000000-0000-4000-8000-000000000001	user logged out	\N	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-01 12:32:32.824244+00
ccbdb4cb-f86c-46a7-a972-2ef06f779836	a1000000-0000-4000-8000-000000000001	a1000000-0000-4000-8000-000000000001	login	user	a1000000-0000-4000-8000-000000000001	user logged in	\N	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-01 13:26:25.186216+00
\.


--
-- Data for Name: auth_login_lockouts; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.auth_login_lockouts (email, failed_attempts, lockout_tier, locked_until, updated_at) FROM stdin;
business@adfluence.com	1	0	\N	2026-10-20 08:51:25.816329+00
\.


--
-- Data for Name: auth_rate_limit_events; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.auth_rate_limit_events (id, scope, rate_key, created_at) FROM stdin;
\.


--
-- Data for Name: bank_accounts; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.bank_accounts (id, user_id, bank_name, account_holder_name, account_number, mobile_money_provider, mobile_money_number, is_primary, is_verified, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: business_profiles; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.business_profiles (id, user_id, owner_full_name, business_name, legal_business_name, logo_url, description, industry, business_type, business_size, city, region, address, tin_number, website_url, instagram_url, tiktok_url, facebook_url, linkedin_url, preferred_platforms, monthly_marketing_budget, is_verified, created_at, updated_at) FROM stdin;
ae03d442-98b8-4ff3-b4fa-bec6c28630be	1429fd1b-6412-4e13-ba93-1887a7a74872	Gossaye Zewde	Gossaye Zewde	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	f	2026-10-20 17:33:27.395307+00	2026-10-20 17:33:27.395307+00
\.


--
-- Data for Name: campaign_deliverables; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.campaign_deliverables (id, campaign_id, deliverable_type, quantity, description, created_at, updated_at) FROM stdin;
1061cb15-824f-4a27-b68c-b6fe6fd3fa18	2499c495-01d1-4824-b064-dfd70a2755b6	tiktok_video	1	TikTok video	2026-10-20 17:35:33.025758+00	2026-10-20 17:35:33.025758+00
f651026b-015f-4836-a900-f1b979827eb8	2499c495-01d1-4824-b064-dfd70a2755b6	instagram_reel	1	Instagram reel	2026-10-20 17:35:33.025758+00	2026-10-20 17:35:33.025758+00
\.


--
-- Data for Name: campaign_invites; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.campaign_invites (id, campaign_id, business_user_id, influencer_user_id, message, status, responded_at, created_at, updated_at) FROM stdin;
3690599e-cfda-4467-a1bd-1b229e6da28d	2499c495-01d1-4824-b064-dfd70a2755b6	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	Hello Yohannes Gossaye, we'd love to invite you to collaborate with us!	accepted	2026-10-20 17:40:05.589549+00	2026-10-20 17:39:07.631407+00	2026-10-20 17:40:05.589549+00
\.


--
-- Data for Name: campaign_platforms; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.campaign_platforms (id, campaign_id, platform, created_at) FROM stdin;
0260d50e-f641-4c41-9698-d3c302a0da00	2499c495-01d1-4824-b064-dfd70a2755b6	tiktok	2026-10-20 17:35:33.025758+00
75dcb730-19dc-43d4-bfd7-d1da6f3fe225	2499c495-01d1-4824-b064-dfd70a2755b6	instagram	2026-10-20 17:35:33.025758+00
c946b40f-572b-4562-ac1b-433dc16180f6	2499c495-01d1-4824-b064-dfd70a2755b6	youtube	2026-10-20 17:35:33.025758+00
\.


--
-- Data for Name: campaign_recommendations; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.campaign_recommendations (id, campaign_id, influencer_user_id, match_score, score_breakdown, rank_position, status, created_at) FROM stdin;
7b4f81e0-41fa-4a43-8f99-c5a3b41ebf54	2499c495-01d1-4824-b064-dfd70a2755b6	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	90.00	{"niche": 40, "budget": 20, "platform": 25, "verified": 5, "creator_size": 0}	1	invited	2026-10-20 17:38:45.438219+00
\.


--
-- Data for Name: campaigns; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.campaigns (id, business_user_id, title, goal, description, category_id, target_audience, city, region, budget_min_etb, budget_max_etb, creator_size_preference, number_of_creators_needed, start_date, end_date, status, created_at, updated_at) FROM stdin;
2499c495-01d1-4824-b064-dfd70a2755b6	1429fd1b-6412-4e13-ba93-1887a7a74872	bole butic	Brand awareness	sell style and fashion wear and cloth	c1000000-0000-4000-8000-000000000002	man 20+	Ethiopia wide	\N	10000.00	25000.00	micro	1	2026-10-20	2026-11-20	active	2026-10-20 17:35:33.025758+00	2026-10-20 17:35:33.025758+00
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.categories (id, name, slug, description, is_active, created_at, updated_at) FROM stdin;
c1000000-0000-4000-8000-000000000001	Beauty	beauty	Beauty and skincare creators	t	2026-10-20 07:26:27.129252+00	2026-10-20 07:26:27.129252+00
c1000000-0000-4000-8000-000000000002	Fashion	fashion	Fashion and styling creators	t	2026-10-20 07:26:27.129252+00	2026-10-20 07:26:27.129252+00
c1000000-0000-4000-8000-000000000003	Food	food	Food and restaurant creators	t	2026-10-20 07:26:27.129252+00	2026-10-20 07:26:27.129252+00
c1000000-0000-4000-8000-000000000004	Tech	tech	Technology and gadget reviewers	t	2026-10-20 07:26:27.129252+00	2026-10-20 07:26:27.129252+00
c1000000-0000-4000-8000-000000000005	Fitness	fitness	Fitness and wellness creators	t	2026-10-20 07:26:27.129252+00	2026-10-20 07:26:27.129252+00
c1000000-0000-4000-8000-000000000006	Lifestyle	lifestyle	General lifestyle creators	t	2026-10-20 07:26:27.129252+00	2026-10-20 07:26:27.129252+00
c1000000-0000-4000-8000-000000000007	Travel	travel	Travel and tourism creators	t	2026-10-20 07:26:27.129252+00	2026-10-20 07:26:27.129252+00
\.


--
-- Data for Name: chat_admin_access_logs; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.chat_admin_access_logs (id, admin_user_id, conversation_id, reason, created_at) FROM stdin;
\.


--
-- Data for Name: contract_deliverables; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.contract_deliverables (id, contract_id, deliverable_type, quantity, description, status, due_date, completed_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: contract_status_history; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.contract_status_history (id, contract_id, old_status, new_status, changed_by, reason, created_at) FROM stdin;
\.


--
-- Data for Name: contracts; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.contracts (id, campaign_id, business_user_id, influencer_user_id, conversation_id, invite_id, title, description, amount_etb, business_fee_percent, influencer_fee_percent, business_fee_amount_etb, influencer_fee_amount_etb, total_business_payable_etb, influencer_receivable_etb, platform_total_fee_etb, start_date, end_date, status, accepted_by_business_at, accepted_by_influencer_at, completed_at, cancelled_at, cancellation_reason, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: conversation_participants; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.conversation_participants (id, conversation_id, user_id, role, last_read_at, is_muted, created_at) FROM stdin;
041b6780-d032-4b31-8efd-17bfa710e7b9	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	influencer	2026-10-20 17:46:32.768423+00	f	2026-10-20 17:39:15.383425+00
80de05a1-0813-41df-bf5b-b94d14f35885	7c730843-fab6-4a41-a0e3-fa47344582ac	1429fd1b-6412-4e13-ba93-1887a7a74872	business	2026-11-01 11:59:36.84529+00	f	2026-10-20 17:39:15.383425+00
\.


--
-- Data for Name: conversations; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.conversations (id, campaign_id, status, last_message_preview, last_message_at, created_at, updated_at) FROM stdin;
7c730843-fab6-4a41-a0e3-fa47344582ac	2499c495-01d1-4824-b064-dfd70a2755b6	active	hey	2026-11-01 12:30:36.852569+00	2026-10-20 17:39:15.383425+00	2026-11-01 12:30:36.852569+00
\.


--
-- Data for Name: disputes; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.disputes (id, contract_id, opened_by_user_id, against_user_id, assigned_admin_user_id, issue_type, description, severity, status, resolution_note, resolved_by, created_at, updated_at, resolved_at) FROM stdin;
\.


--
-- Data for Name: email_verification_tokens; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.email_verification_tokens (id, user_id, token_hash, expires_at, used_at, created_at) FROM stdin;
\.


--
-- Data for Name: influencer_portfolio_items; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.influencer_portfolio_items (id, user_id, file_url, object_key, media_type, mime_type, file_size_bytes, sort_order, created_at) FROM stdin;
\.


--
-- Data for Name: influencer_profile_categories; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.influencer_profile_categories (influencer_user_id, category_id, created_at) FROM stdin;
3a709749-dbc8-4918-bfb6-3a0b0ada40b1	c1000000-0000-4000-8000-000000000004	2026-10-20 17:27:40.605862+00
08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	c1000000-0000-4000-8000-000000000002	2026-10-20 17:38:15.935514+00
\.


--
-- Data for Name: influencer_profile_platforms; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.influencer_profile_platforms (influencer_user_id, platform, created_at) FROM stdin;
3a709749-dbc8-4918-bfb6-3a0b0ada40b1	youtube	2026-10-20 17:27:40.605862+00
08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	youtube	2026-10-20 17:38:15.935514+00
08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	tiktok	2026-10-20 17:38:15.935514+00
08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	instagram	2026-10-20 17:38:15.935514+00
\.


--
-- Data for Name: influencer_profiles; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.influencer_profiles (id, user_id, full_name, creator_name, username, profile_image_url, bio, city, region, languages, social_links, instagram_followers, tiktok_followers, youtube_subscribers, total_followers, engagement_rate, audience_description, content_style, portfolio_description, min_price_etb, max_price_etb, approval_note, is_verified, approved_at, approved_by, created_at, updated_at) FROM stdin;
fa243c5a-24ee-4988-9c9e-4748c8cf4d5e	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	wini late	wini late	\N	https://13-49-175-98.sslip.io/uploads/media/influencer/3a709749-dbc8-4918-bfb6-3a0b0ada40b1/avatar/9d918ce8-1680-4dcc-9947-8f93baf7cc3a.jpg		Addis Ababa	\N	{Amharic,English}	{}	0	0	84000	84000	7.20	\N	\N	\N	15000.00	35000.00	add tiktoke and instagram follower	f	\N	\N	2026-10-20 07:55:56.289013+00	2026-10-20 17:31:02.453025+00
c44e5fab-7d0c-4e9a-ba51-72fa69a5bd65	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	Yohannes Gossaye	Yohannes Gossaye	\N	https://13-49-175-98.sslip.io/uploads/media/influencer/08cd0a98-4a57-47c5-93f5-8e4ea2f0392b/avatar/a682fccb-f50c-4830-beda-03eebe3913f1.jpg	fashion	Addis Ababa	\N	{Amharic,English}	{}	0	0	0	84000	7.20	\N	\N	\N	10000.00	35000.00	\N	t	2026-10-20 17:30:39.205777+00	a1000000-0000-4000-8000-000000000001	2026-10-20 17:29:24.405853+00	2026-10-20 17:38:15.935514+00
\.


--
-- Data for Name: message_reports; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.message_reports (id, conversation_id, firebase_message_id, reporter_user_id, reported_user_id, reason, status, resolved_by, resolution_note, created_at, resolved_at) FROM stdin;
\.


--
-- Data for Name: messages; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.messages (id, conversation_id, sender_user_id, receiver_user_id, body, media_url, media_type, is_read, read_at, deleted_for_sender, deleted_for_receiver, is_deleted_for_everyone, created_at, updated_at) FROM stdin;
eb1ef3d2-17a1-4b47-b747-8e97bc298368	7c730843-fab6-4a41-a0e3-fa47344582ac	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	hey john how are you	\N	text	t	2026-10-20 17:39:56.861198+00	f	f	f	2026-10-20 17:39:28.759554+00	2026-10-20 17:39:56.861198+00
53e60127-d799-4df8-a79a-275ac4917723	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	how are u	\N	text	t	2026-10-20 17:45:19.75354+00	f	f	t	2026-10-20 17:40:25.894821+00	2026-10-20 17:45:19.75354+00
aa826be8-e932-4671-adfa-9ca57a7d4fe0	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	yes i am fine did u see my offer	\N	text	t	2026-10-20 17:45:19.75354+00	f	f	t	2026-10-20 17:40:41.800016+00	2026-10-20 17:45:19.75354+00
7f232620-13c6-40cd-a06e-4c44efe9ff1f	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	hey	\N	text	t	2026-10-20 17:45:19.75354+00	f	f	t	2026-10-20 17:41:39.145983+00	2026-10-20 17:45:19.75354+00
a22368b3-8ac6-4995-b25a-fb4ba28defcc	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	kklkl	\N	text	t	2026-10-20 17:45:19.75354+00	f	f	t	2026-10-20 17:41:57.922775+00	2026-10-20 17:45:19.75354+00
82c8b45d-e62a-418c-b052-cc99134e4389	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	hey	\N	text	t	2026-10-20 17:45:19.75354+00	f	f	f	2026-10-20 17:44:46.269687+00	2026-10-20 17:45:19.75354+00
74da0817-7c67-49d0-8f76-1789d91d8c01	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	hey how are u	\N	text	t	2026-10-20 17:45:19.75354+00	f	f	f	2026-10-20 17:45:01.808421+00	2026-10-20 17:45:19.75354+00
f0413810-aa0a-433d-af78-323297a76599	7c730843-fab6-4a41-a0e3-fa47344582ac	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	i am fine	\N	text	t	2026-10-20 17:46:12.501787+00	f	f	f	2026-10-20 17:45:32.099683+00	2026-10-20 17:46:12.501787+00
10ac58b5-1826-44ca-88d6-0678b0668a68	7c730843-fab6-4a41-a0e3-fa47344582ac	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	30k did u agree	\N	text	t	2026-10-20 17:46:12.501787+00	f	f	f	2026-10-20 17:45:48.006507+00	2026-10-20 17:46:12.501787+00
0fd772fd-3b84-4f00-8608-45900a8e42df	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	Please send the contract	\N	text	t	2026-10-20 17:46:53.897872+00	f	f	f	2026-10-20 17:46:25.205127+00	2026-10-20 17:46:53.897872+00
43f19fc7-42f4-4804-be5f-6febe2dcbbe0	7c730843-fab6-4a41-a0e3-fa47344582ac	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	1429fd1b-6412-4e13-ba93-1887a7a74872	for 40k	\N	text	t	2026-10-20 17:46:53.897872+00	f	f	f	2026-10-20 17:46:33.523721+00	2026-10-20 17:46:53.897872+00
19923da6-1aec-4d32-b3e8-f01a5f727fd2	7c730843-fab6-4a41-a0e3-fa47344582ac	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	not that is too much maxmum 35k	\N	text	f	\N	f	f	f	2026-10-20 17:47:16.147955+00	2026-10-20 17:47:16.147955+00
3a4964b3-d6de-4bfb-895d-342862dd9fc4	7c730843-fab6-4a41-a0e3-fa47344582ac	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	hey	\N	text	f	\N	f	f	f	2026-11-01 12:30:31.727967+00	2026-11-01 12:30:31.727967+00
496beffe-104c-4bd5-9bc2-43c6a105c794	7c730843-fab6-4a41-a0e3-fa47344582ac	1429fd1b-6412-4e13-ba93-1887a7a74872	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	hey	\N	text	f	\N	f	f	f	2026-11-01 12:30:36.852569+00	2026-11-01 12:30:36.852569+00
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.notifications (id, user_id, title, message, type, entity_type, entity_id, read_at, created_at) FROM stdin;
\.


--
-- Data for Name: password_reset_tokens; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.password_reset_tokens (id, user_id, token_hash, expires_at, used_at, created_at) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.payments (id, contract_id, business_user_id, influencer_user_id, amount_etb, business_fee_etb, influencer_fee_etb, platform_fee_total_etb, total_paid_by_business_etb, influencer_payout_etb, payment_method, provider, provider_reference, status, paid_at, failed_at, failure_reason, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: payouts; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.payouts (id, payment_id, influencer_user_id, bank_account_id, amount_etb, status, provider, provider_reference, paid_at, failed_at, failure_reason, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: phone_verification_tokens; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.phone_verification_tokens (id, user_id, token_hash, expires_at, used_at, created_at) FROM stdin;
\.


--
-- Data for Name: reviews; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.reviews (id, contract_id, reviewer_user_id, reviewed_user_id, rating, comment, created_at) FROM stdin;
\.


--
-- Data for Name: support_messages; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.support_messages (id, ticket_id, sender_user_id, message, attachment_url, is_internal_note, created_at) FROM stdin;
\.


--
-- Data for Name: support_tickets; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.support_tickets (id, created_by_user_id, assigned_admin_user_id, related_campaign_id, related_contract_id, subject, message, issue_type, priority, status, resolution_note, resolved_by, created_at, updated_at, resolved_at) FROM stdin;
\.


--
-- Data for Name: user_sessions; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.user_sessions (id, user_id, refresh_token_hash, ip_address, user_agent, expires_at, revoked_at, created_at, refresh_ttl_seconds, rotated_at) FROM stdin;
193ede08-3e7c-498f-8076-79b9c58baa00	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	619c52e7bbd3fa16e32f5c685a3644862c7299ed2f1c6ad785741fdc1d0ed24b	196.188.179.241	PostmanRuntime/2.12.0	2026-11-19 07:58:42.384138+00	\N	2026-10-20 07:58:42.38414+00	\N	\N
a381e820-1093-490e-a014-b03129d5096c	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	9b463ae0d480d73a285417341a2032b2efe31a5cbaf74054e1b6827f8d1e80af	196.188.179.241	PostmanRuntime/2.12.0	2026-11-19 07:58:52.770989+00	\N	2026-10-20 07:58:52.770991+00	\N	\N
9fc36600-a9b8-457b-937e-4af14ad2bf7b	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	d53c8045fe3b2146d44a7126b767ca2ef73892bb68749cdcda05b6f7b6877f21	196.188.179.241	PostmanRuntime/2.12.0	2026-11-19 07:58:57.823831+00	\N	2026-10-20 07:58:57.823833+00	\N	\N
96205929-e641-44a4-92b6-938f432ef928	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	9463457df420cd07ede3bb48129ba720541fa0cb8968477a39c9ebf9c5449cb0	196.188.179.241	PostmanRuntime/2.12.0	2026-11-19 08:13:44.829803+00	\N	2026-10-20 08:13:44.829805+00	\N	\N
e9b0c650-cf0c-40e7-923c-2c17ae327405	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	55b68865728a04bb9ac35a7ec8dfe1d58a0518bd55505e802fcde3243e3edd48	196.188.179.241	PostmanRuntime/2.12.0	2026-11-19 08:13:47.50174+00	\N	2026-10-20 08:13:47.501742+00	\N	\N
b020b4a8-ec4c-45f7-bd00-62508e7cf069	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	46f213e4a1a740f3e97886cb2d19f6e7196a70754349c948125da7b07bd5289a	196.188.179.241	PostmanRuntime/2.12.0	2026-11-19 08:13:50.146396+00	\N	2026-10-20 08:13:50.146398+00	\N	\N
0d6bca16-6121-478f-b313-82fe542844c8	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	ff983e5573ece908433dae0436a5f75ee348beac12f590296909da5cbe18291c	196.188.179.241	PostmanRuntime/2.12.0	2026-11-19 08:32:47.614112+00	\N	2026-10-20 08:32:47.614114+00	\N	\N
2a41e103-5aa4-4182-90d4-514f76aed97b	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	cbb7849e6d72295b2f0e894f040de446f9f1678168d7661b1717545eb4bf5cd7	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-27 07:55:56.293748+00	2026-10-20 08:49:32.504189+00	2026-10-20 07:55:56.29375+00	\N	\N
1b67138f-cb97-4632-b55f-e532d05e5853	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	3cf85ba267be23fdca2aff788dd9703ca59a84323b60cfbe6a53a4286eedc5cb	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-27 08:49:32.505718+00	\N	2026-10-20 08:49:32.50572+00	\N	\N
68e579e6-f362-4f59-95db-a8dcf18f49ad	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	88d8e12a98029996a6481f5c1e83f9db0b006a0b6cf5118ffd5aea2c43e13b8e	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 08:51:46.137801+00	\N	2026-10-20 08:51:46.137802+00	\N	\N
066ea2dc-b7d4-4bd8-a64c-547c2b07a020	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	2e31e663e31d89534ea774fdc51fc9cd36ca80ce06872bf8358d5c0fd84e2b00	196.188.179.241	Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Mobile Safari/537.36	2026-11-19 08:54:41.930975+00	\N	2026-10-20 08:54:41.930977+00	\N	\N
97977234-c9b4-4070-a170-ebfacdbbd751	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	89d68986a1702e7e658621edced5b13dfbb3f659b69460654c53b5e2e892c1ec	102.213.71.203	PostmanRuntime/2.12.0	2026-11-19 13:39:20.778116+00	\N	2026-10-20 13:39:20.778118+00	\N	\N
879cbdd1-e142-4691-9abd-2de9e61b3a24	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	f390e9ef09e5c4eccac47eff5b8e0dc8dfa9ce04a9bbd5592504e7edffcc8eeb	196.188.179.241	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 08:51:55.919997+00	2026-10-20 17:27:17.011378+00	2026-10-20 08:51:55.919999+00	\N	\N
194fc34b-d312-44bb-a438-31b2c20d353e	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	264f59b75a8cbc269d006001490347b889d3d1436c35e1cb560701acfdda5382	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-27 17:27:17.012243+00	\N	2026-10-20 17:27:17.012246+00	\N	\N
ce063242-5edf-4238-8bc5-a04ed12f424a	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	98a24830a31fe0b0686a013541aafdb9725568de720181d999e0974d23111067	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:27:22.094263+00	\N	2026-10-20 17:27:22.094265+00	\N	\N
3eac6941-29e9-4b4f-9fda-fdf54f236b78	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	9266b168ad76022709c5fddc4d831887c258d0215171b92e1ee9c76b243f98ce	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-27 17:29:24.408034+00	\N	2026-10-20 17:29:24.408037+00	\N	\N
279e94bd-90b5-497e-8d7f-4971c76ad934	a1000000-0000-4000-8000-000000000001	eb587255c65abee03e5cd4d63236e1c49d2196af19b6a2b1103c17cc6f62bc87	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:30:07.428411+00	\N	2026-10-20 17:30:07.428413+00	\N	\N
61e533b6-09f8-4e91-ab32-af395d4f2736	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	375d98a5544ce055c44069376238e365a216491cc22a4663cdb95d41b3902a7a	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:31:31.273031+00	\N	2026-10-20 17:31:31.273033+00	\N	\N
65751dac-2798-4bd4-a71b-5df3583cabf2	1429fd1b-6412-4e13-ba93-1887a7a74872	1a848345afe286d43ad92d74907390d4027c94df0133f9421c8fb8b68c72a128	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-27 17:33:27.398103+00	\N	2026-10-20 17:33:27.398105+00	\N	\N
df850def-3426-46b0-af1f-995578bfdc09	a1000000-0000-4000-8000-000000000001	b5f1a014aa02ba8ba2cfc9606cb19ca0c74621838fba926f626b4cf6f4bbce0f	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:36:15.042621+00	\N	2026-10-20 17:36:15.042624+00	\N	\N
bea57051-947c-4734-be4a-a2f358dfb3d2	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	c9b8ebadfb946b93a582142fec356fcf552f44fa709c40d832b657875bbd9620	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:37:36.630261+00	\N	2026-10-20 17:37:36.630263+00	\N	\N
399b3628-7fbd-4049-b28f-071be3fafe70	1429fd1b-6412-4e13-ba93-1887a7a74872	1db968e75e348d98b872ef0089f3e840a73fe0e2462a663faf4771cd57686d22	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:38:32.239309+00	\N	2026-10-20 17:38:32.239311+00	\N	\N
2e2a42e3-318a-4262-be40-829aad58a779	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	e472f666ed7bc98c95b33cd895c64cd9812073137e39537ce43c31ffb74b05e0	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:39:55.41878+00	\N	2026-10-20 17:39:55.418783+00	\N	\N
9a5411b7-f8bf-4a54-883c-26d64e28c795	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	c4877098dd504f6aaa1af30493a1f7f659c3370a4141d73bebb16c6a2bf6f7dc	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:44:28.465641+00	\N	2026-10-20 17:44:28.465643+00	\N	\N
85d89046-94cd-452a-aa48-4d2b526e84ca	1429fd1b-6412-4e13-ba93-1887a7a74872	8938665afc49150b08ddf8156204f8bde427544e28dc2314e49533a165cb230c	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:45:15.241192+00	\N	2026-10-20 17:45:15.241194+00	\N	\N
4cfe9ccd-e7ce-42be-902f-6946618757f2	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	ab9cf41cf836159214fe735da032288498a0a00ce3ebda413611ab6018fd093e	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:46:03.42503+00	\N	2026-10-20 17:46:03.425032+00	\N	\N
c6da35df-35b7-45c9-90cd-756004c72303	08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	b961d39af9ae2a8eb63b41f9150814b7e41bd9feb2d52de57f96dd98a6634aa2	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:46:08.668125+00	\N	2026-10-20 17:46:08.668127+00	\N	\N
39607e90-ff40-4c0c-87b7-9720eb4a0777	1429fd1b-6412-4e13-ba93-1887a7a74872	1e120f350b0aef9dd995699ed4c4b88e6f90619b7960cc6342de1f4e2e137a8e	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:46:43.492902+00	\N	2026-10-20 17:46:43.492904+00	\N	\N
b38c97d4-09eb-4b01-b173-fccba083e14a	3a709749-dbc8-4918-bfb6-3a0b0ada40b1	16c6ec73137d7099029ce64c9ff48cd03f109324ca8d8b5333ca4dc863329323	196.188.179.241	PostmanRuntime/2.12.0	2026-11-20 06:50:12.869316+00	\N	2026-10-21 06:50:12.869318+00	\N	\N
a74d7a4e-fc6b-4ce7-b063-ed4b3c8024dc	1429fd1b-6412-4e13-ba93-1887a7a74872	6515e572b4e543a94fcf7ce7b7e701b6c66ebc1667ccbd14d8b8817f7beae411	102.213.71.199	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-19 17:46:49.203633+00	2026-11-01 09:06:04.081877+00	2026-10-20 17:46:49.203635+00	\N	\N
a1eadc41-e4bd-42fa-8811-f3ef90ad25cf	1429fd1b-6412-4e13-ba93-1887a7a74872	df46b7387feb98637d32912d3e4c233338b5cecd9cb1fef4fbfac16ae4d1e7d2	102.213.71.192	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-08 09:06:04.083503+00	\N	2026-11-01 09:06:04.083507+00	\N	\N
e7a4df18-e6e3-4f7a-9d81-be61b139f1d0	1429fd1b-6412-4e13-ba93-1887a7a74872	568be8303d7318e9dbfd5d972404e7f8cd832fedf572797e6822cda26e691786	102.213.71.192	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-12-01 09:06:07.375324+00	\N	2026-11-01 09:06:07.375327+00	\N	\N
87b304e3-bac7-4665-bf8d-52387cebc6e5	1429fd1b-6412-4e13-ba93-1887a7a74872	f121409bd119b2648650c26c36a0b38613b14b80f32fa17b0c1266b0787b4e70	102.213.71.192	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-12-01 09:07:53.932+00	2026-11-01 10:32:29.845302+00	2026-11-01 09:07:53.932002+00	\N	\N
33161ea7-a170-4a31-aa95-ce4929293537	1429fd1b-6412-4e13-ba93-1887a7a74872	f72762790d95f6f61e803195864316b869b5823a6ff6e51c20239afe8b5da460	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-08 10:32:29.847202+00	2026-11-01 10:32:30.240322+00	2026-11-01 10:32:29.847205+00	\N	\N
b4155bc0-2cfa-4abf-a754-28cf3f8d2f27	1429fd1b-6412-4e13-ba93-1887a7a74872	97c8e10357ab5b448a829505f9614f32e6953219c911883a928d0b2cc4c5810c	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-08 10:32:30.241707+00	\N	2026-11-01 10:32:30.24171+00	\N	\N
43fda16f-a768-4616-90f2-31dd0d8d8ee5	1429fd1b-6412-4e13-ba93-1887a7a74872	6577d8da185e802be4a27c480fad3026f1d4b3bf340fc80d60c3b9260aa1cdf8	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-12-01 11:59:18.494862+00	2026-11-01 12:29:53.798435+00	2026-11-01 11:59:18.494866+00	\N	\N
db860d8b-9df1-436c-97a4-0c231a73cc98	1429fd1b-6412-4e13-ba93-1887a7a74872	dfb34056e2e766552c43c21afd5589e640574dbbb9e0c01d60c80e90a87b0388	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-08 12:29:53.799924+00	2026-11-01 12:31:03.61033+00	2026-11-01 12:29:53.799926+00	\N	\N
f81cf730-6ec3-471c-be05-0e063b99c773	a1000000-0000-4000-8000-000000000001	397cef8af86e090bb902178baa0282094e7b0762f1d0a5fcebb7578902b7b45d	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-12-01 12:31:11.59781+00	2026-11-01 12:32:32.822643+00	2026-11-01 12:31:11.597812+00	\N	\N
ddc47cb5-82ee-4bf0-9782-37f8e4f6fdf4	a1000000-0000-4000-8000-000000000001	da85c6c9ce6b9f749b1dd4cc2acb77d16749b99ec823fadfeaf8635932f7b907	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-12-01 13:26:25.184595+00	2026-11-01 14:38:18.836035+00	2026-11-01 13:26:25.184597+00	\N	\N
26e97b45-6f12-4671-b8fb-1f53e7a2f5f3	a1000000-0000-4000-8000-000000000001	68f6dea0332a45f3fedb7ac4ffe9686886a197d757e9345e7f6faabeb6720f8d	102.213.71.237	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36	2026-11-08 14:38:18.837438+00	\N	2026-11-01 14:38:18.83744+00	\N	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.users (id, email, phone_number, password_hash, role, status, email_verified, phone_verified, last_login_at, suspended_until, status_reason, created_at, updated_at) FROM stdin;
1429fd1b-6412-4e13-ba93-1887a7a74872	gossaye@gmail.com	251911675623	$argon2id$v=19$m=65536,t=4,p=2$/1mN4xEVwL0TVQWEKcPPWQ$u4kralekg5rvsZTfPe/XsgRi4xYo403cQOtC1rgbOzA	business	pending_profile	f	f	2026-11-01 11:59:18.49335+00	\N	\N	2026-10-20 17:33:27.395307+00	2026-11-01 11:59:18.49335+00
a1000000-0000-4000-8000-000000000001	admin@adfluence.com	\N	$argon2id$v=19$m=65536,t=4,p=2$SDDP2Vq45jYwGPUG3auV3g$B/d7aBWuRin0+R9uU0E4TmTxvlY4oY2Xqa+/03fyltA	admin	active	t	f	2026-11-01 13:26:25.183229+00	\N	\N	2026-10-20 07:26:27.130998+00	2026-11-01 13:26:25.183229+00
08cd0a98-4a57-47c5-93f5-8e4ea2f0392b	yohannesgossaye73@gmail.com	251953422606	$argon2id$v=19$m=65536,t=4,p=2$EATjVk7U2peDYH526As2+A$e04ba1KOR6CMxrnOBIThni5p0X8gNSsHNQ5mxI7fTRU	influencer	active	f	f	2026-10-20 17:46:08.666852+00	\N	\N	2026-10-20 17:29:24.405853+00	2026-10-20 17:46:08.666852+00
3a709749-dbc8-4918-bfb6-3a0b0ada40b1	wini23@gmail.com	251946627616	$argon2id$v=19$m=65536,t=4,p=2$IWs/BKYdhuLTikkVTAAkag$k/mH6J/g8aTVnlkSFvYhgibOTZo+d+KzPTcZx934XNU	influencer	rejected	f	f	2026-10-21 06:50:12.868021+00	\N	Application rejected by admin.	2026-10-20 07:55:56.289013+00	2026-10-21 06:50:12.868021+00
\.


--
-- Data for Name: verification_documents; Type: TABLE DATA; Schema: public; Owner: adfluence
--

COPY public.verification_documents (id, user_id, document_type, file_url, file_name, mime_type, file_size_bytes, status, rejection_reason, reviewed_by, reviewed_at, created_at, updated_at) FROM stdin;
\.


--
-- Name: admin_profiles admin_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.admin_profiles
    ADD CONSTRAINT admin_profiles_pkey PRIMARY KEY (id);


--
-- Name: admin_profiles admin_profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.admin_profiles
    ADD CONSTRAINT admin_profiles_user_id_key UNIQUE (user_id);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: auth_login_lockouts auth_login_lockouts_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.auth_login_lockouts
    ADD CONSTRAINT auth_login_lockouts_pkey PRIMARY KEY (email);


--
-- Name: auth_rate_limit_events auth_rate_limit_events_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.auth_rate_limit_events
    ADD CONSTRAINT auth_rate_limit_events_pkey PRIMARY KEY (id);


--
-- Name: bank_accounts bank_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.bank_accounts
    ADD CONSTRAINT bank_accounts_pkey PRIMARY KEY (id);


--
-- Name: business_profiles business_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.business_profiles
    ADD CONSTRAINT business_profiles_pkey PRIMARY KEY (id);


--
-- Name: business_profiles business_profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.business_profiles
    ADD CONSTRAINT business_profiles_user_id_key UNIQUE (user_id);


--
-- Name: campaign_deliverables campaign_deliverables_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_deliverables
    ADD CONSTRAINT campaign_deliverables_pkey PRIMARY KEY (id);


--
-- Name: campaign_invites campaign_invites_campaign_id_influencer_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_invites
    ADD CONSTRAINT campaign_invites_campaign_id_influencer_user_id_key UNIQUE (campaign_id, influencer_user_id);


--
-- Name: campaign_invites campaign_invites_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_invites
    ADD CONSTRAINT campaign_invites_pkey PRIMARY KEY (id);


--
-- Name: campaign_platforms campaign_platforms_campaign_id_platform_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_platforms
    ADD CONSTRAINT campaign_platforms_campaign_id_platform_key UNIQUE (campaign_id, platform);


--
-- Name: campaign_platforms campaign_platforms_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_platforms
    ADD CONSTRAINT campaign_platforms_pkey PRIMARY KEY (id);


--
-- Name: campaign_recommendations campaign_recommendations_campaign_id_influencer_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_recommendations
    ADD CONSTRAINT campaign_recommendations_campaign_id_influencer_user_id_key UNIQUE (campaign_id, influencer_user_id);


--
-- Name: campaign_recommendations campaign_recommendations_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_recommendations
    ADD CONSTRAINT campaign_recommendations_pkey PRIMARY KEY (id);


--
-- Name: campaigns campaigns_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: categories categories_slug_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);


--
-- Name: chat_admin_access_logs chat_admin_access_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.chat_admin_access_logs
    ADD CONSTRAINT chat_admin_access_logs_pkey PRIMARY KEY (id);


--
-- Name: contract_deliverables contract_deliverables_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contract_deliverables
    ADD CONSTRAINT contract_deliverables_pkey PRIMARY KEY (id);


--
-- Name: contract_status_history contract_status_history_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contract_status_history
    ADD CONSTRAINT contract_status_history_pkey PRIMARY KEY (id);


--
-- Name: contracts contracts_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contracts
    ADD CONSTRAINT contracts_pkey PRIMARY KEY (id);


--
-- Name: conversation_participants conversation_participants_conversation_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.conversation_participants
    ADD CONSTRAINT conversation_participants_conversation_id_user_id_key UNIQUE (conversation_id, user_id);


--
-- Name: conversation_participants conversation_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.conversation_participants
    ADD CONSTRAINT conversation_participants_pkey PRIMARY KEY (id);


--
-- Name: conversations conversations_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.conversations
    ADD CONSTRAINT conversations_pkey PRIMARY KEY (id);


--
-- Name: disputes disputes_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.disputes
    ADD CONSTRAINT disputes_pkey PRIMARY KEY (id);


--
-- Name: email_verification_tokens email_verification_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT email_verification_tokens_pkey PRIMARY KEY (id);


--
-- Name: email_verification_tokens email_verification_tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT email_verification_tokens_token_hash_key UNIQUE (token_hash);


--
-- Name: influencer_portfolio_items influencer_portfolio_items_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_portfolio_items
    ADD CONSTRAINT influencer_portfolio_items_pkey PRIMARY KEY (id);


--
-- Name: influencer_portfolio_items influencer_portfolio_items_user_id_sort_order_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_portfolio_items
    ADD CONSTRAINT influencer_portfolio_items_user_id_sort_order_key UNIQUE (user_id, sort_order);


--
-- Name: influencer_profile_categories influencer_profile_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profile_categories
    ADD CONSTRAINT influencer_profile_categories_pkey PRIMARY KEY (influencer_user_id, category_id);


--
-- Name: influencer_profile_platforms influencer_profile_platforms_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profile_platforms
    ADD CONSTRAINT influencer_profile_platforms_pkey PRIMARY KEY (influencer_user_id, platform);


--
-- Name: influencer_profiles influencer_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profiles
    ADD CONSTRAINT influencer_profiles_pkey PRIMARY KEY (id);


--
-- Name: influencer_profiles influencer_profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profiles
    ADD CONSTRAINT influencer_profiles_user_id_key UNIQUE (user_id);


--
-- Name: influencer_profiles influencer_profiles_username_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profiles
    ADD CONSTRAINT influencer_profiles_username_key UNIQUE (username);


--
-- Name: message_reports message_reports_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.message_reports
    ADD CONSTRAINT message_reports_pkey PRIMARY KEY (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: password_reset_tokens password_reset_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (id);


--
-- Name: password_reset_tokens password_reset_tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_token_hash_key UNIQUE (token_hash);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: payouts payouts_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_pkey PRIMARY KEY (id);


--
-- Name: phone_verification_tokens phone_verification_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.phone_verification_tokens
    ADD CONSTRAINT phone_verification_tokens_pkey PRIMARY KEY (id);


--
-- Name: phone_verification_tokens phone_verification_tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.phone_verification_tokens
    ADD CONSTRAINT phone_verification_tokens_token_hash_key UNIQUE (token_hash);


--
-- Name: reviews reviews_contract_id_reviewer_user_id_reviewed_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_contract_id_reviewer_user_id_reviewed_user_id_key UNIQUE (contract_id, reviewer_user_id, reviewed_user_id);


--
-- Name: reviews reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id);


--
-- Name: support_messages support_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_messages
    ADD CONSTRAINT support_messages_pkey PRIMARY KEY (id);


--
-- Name: support_tickets support_tickets_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_tickets
    ADD CONSTRAINT support_tickets_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: verification_documents verification_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.verification_documents
    ADD CONSTRAINT verification_documents_pkey PRIMARY KEY (id);


--
-- Name: verification_documents verification_documents_user_id_key; Type: CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.verification_documents
    ADD CONSTRAINT verification_documents_user_id_key UNIQUE (user_id);


--
-- Name: idx_audit_logs_action; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_audit_logs_action ON public.audit_logs USING btree (action);


--
-- Name: idx_audit_logs_actor; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_audit_logs_actor ON public.audit_logs USING btree (actor_user_id);


--
-- Name: idx_audit_logs_created_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_audit_logs_created_at ON public.audit_logs USING btree (created_at);


--
-- Name: idx_audit_logs_target; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_audit_logs_target ON public.audit_logs USING btree (target_user_id);


--
-- Name: idx_auth_rate_limit_scope_key_created; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_auth_rate_limit_scope_key_created ON public.auth_rate_limit_events USING btree (scope, rate_key, created_at DESC);


--
-- Name: idx_bank_accounts_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_bank_accounts_user_id ON public.bank_accounts USING btree (user_id);


--
-- Name: idx_business_profiles_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_business_profiles_user_id ON public.business_profiles USING btree (user_id);


--
-- Name: idx_campaign_deliverables_campaign_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_deliverables_campaign_id ON public.campaign_deliverables USING btree (campaign_id);


--
-- Name: idx_campaign_invites_business_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_invites_business_user_id ON public.campaign_invites USING btree (business_user_id);


--
-- Name: idx_campaign_invites_campaign_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_invites_campaign_id ON public.campaign_invites USING btree (campaign_id);


--
-- Name: idx_campaign_invites_influencer_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_invites_influencer_user_id ON public.campaign_invites USING btree (influencer_user_id);


--
-- Name: idx_campaign_invites_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_invites_status ON public.campaign_invites USING btree (status);


--
-- Name: idx_campaign_platforms_campaign_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_platforms_campaign_id ON public.campaign_platforms USING btree (campaign_id);


--
-- Name: idx_campaign_recommendations_campaign_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_recommendations_campaign_id ON public.campaign_recommendations USING btree (campaign_id);


--
-- Name: idx_campaign_recommendations_influencer; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_recommendations_influencer ON public.campaign_recommendations USING btree (influencer_user_id);


--
-- Name: idx_campaign_recommendations_rank; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_recommendations_rank ON public.campaign_recommendations USING btree (campaign_id, rank_position);


--
-- Name: idx_campaign_recommendations_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaign_recommendations_status ON public.campaign_recommendations USING btree (status);


--
-- Name: idx_campaigns_business_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaigns_business_user_id ON public.campaigns USING btree (business_user_id);


--
-- Name: idx_campaigns_category_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaigns_category_id ON public.campaigns USING btree (category_id);


--
-- Name: idx_campaigns_created_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaigns_created_at ON public.campaigns USING btree (created_at);


--
-- Name: idx_campaigns_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_campaigns_status ON public.campaigns USING btree (status);


--
-- Name: idx_categories_is_active; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_categories_is_active ON public.categories USING btree (is_active);


--
-- Name: idx_categories_slug; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_categories_slug ON public.categories USING btree (slug);


--
-- Name: idx_chat_admin_access_logs_admin; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_chat_admin_access_logs_admin ON public.chat_admin_access_logs USING btree (admin_user_id);


--
-- Name: idx_chat_admin_access_logs_conversation; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_chat_admin_access_logs_conversation ON public.chat_admin_access_logs USING btree (conversation_id);


--
-- Name: idx_contract_deliverables_contract_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contract_deliverables_contract_id ON public.contract_deliverables USING btree (contract_id);


--
-- Name: idx_contract_deliverables_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contract_deliverables_status ON public.contract_deliverables USING btree (status);


--
-- Name: idx_contract_status_history_contract_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contract_status_history_contract_id ON public.contract_status_history USING btree (contract_id);


--
-- Name: idx_contract_status_history_created_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contract_status_history_created_at ON public.contract_status_history USING btree (created_at);


--
-- Name: idx_contracts_business_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contracts_business_user_id ON public.contracts USING btree (business_user_id);


--
-- Name: idx_contracts_campaign_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contracts_campaign_id ON public.contracts USING btree (campaign_id);


--
-- Name: idx_contracts_conversation_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contracts_conversation_id ON public.contracts USING btree (conversation_id);


--
-- Name: idx_contracts_created_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contracts_created_at ON public.contracts USING btree (created_at);


--
-- Name: idx_contracts_influencer_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contracts_influencer_user_id ON public.contracts USING btree (influencer_user_id);


--
-- Name: idx_contracts_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_contracts_status ON public.contracts USING btree (status);


--
-- Name: idx_conversation_participants_conversation; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_conversation_participants_conversation ON public.conversation_participants USING btree (conversation_id);


--
-- Name: idx_conversation_participants_user; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_conversation_participants_user ON public.conversation_participants USING btree (user_id);


--
-- Name: idx_conversations_campaign_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_conversations_campaign_id ON public.conversations USING btree (campaign_id);


--
-- Name: idx_conversations_last_message_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_conversations_last_message_at ON public.conversations USING btree (last_message_at DESC);


--
-- Name: idx_disputes_against_user; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_disputes_against_user ON public.disputes USING btree (against_user_id);


--
-- Name: idx_disputes_contract_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_disputes_contract_id ON public.disputes USING btree (contract_id);


--
-- Name: idx_disputes_opened_by; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_disputes_opened_by ON public.disputes USING btree (opened_by_user_id);


--
-- Name: idx_disputes_severity; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_disputes_severity ON public.disputes USING btree (severity);


--
-- Name: idx_disputes_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_disputes_status ON public.disputes USING btree (status);


--
-- Name: idx_email_verification_tokens_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_email_verification_tokens_user_id ON public.email_verification_tokens USING btree (user_id);


--
-- Name: idx_influencer_portfolio_items_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_influencer_portfolio_items_user_id ON public.influencer_portfolio_items USING btree (user_id);


--
-- Name: idx_influencer_profile_categories_category; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_influencer_profile_categories_category ON public.influencer_profile_categories USING btree (category_id);


--
-- Name: idx_influencer_profile_platforms_platform; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_influencer_profile_platforms_platform ON public.influencer_profile_platforms USING btree (platform);


--
-- Name: idx_influencer_profiles_is_verified; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_influencer_profiles_is_verified ON public.influencer_profiles USING btree (is_verified);


--
-- Name: idx_influencer_profiles_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_influencer_profiles_user_id ON public.influencer_profiles USING btree (user_id);


--
-- Name: idx_influencer_profiles_username; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_influencer_profiles_username ON public.influencer_profiles USING btree (username);


--
-- Name: idx_message_reports_conversation; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_message_reports_conversation ON public.message_reports USING btree (conversation_id);


--
-- Name: idx_message_reports_reporter; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_message_reports_reporter ON public.message_reports USING btree (reporter_user_id);


--
-- Name: idx_message_reports_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_message_reports_status ON public.message_reports USING btree (status);


--
-- Name: idx_messages_conversation_created; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_messages_conversation_created ON public.messages USING btree (conversation_id, created_at DESC);


--
-- Name: idx_messages_is_read; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_messages_is_read ON public.messages USING btree (conversation_id, receiver_user_id, is_read);


--
-- Name: idx_messages_receiver; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_messages_receiver ON public.messages USING btree (receiver_user_id);


--
-- Name: idx_messages_sender; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_messages_sender ON public.messages USING btree (sender_user_id);


--
-- Name: idx_notifications_created_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_notifications_created_at ON public.notifications USING btree (created_at DESC);


--
-- Name: idx_notifications_read_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_notifications_read_at ON public.notifications USING btree (user_id, read_at);


--
-- Name: idx_notifications_type; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_notifications_type ON public.notifications USING btree (type);


--
-- Name: idx_notifications_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_notifications_user_id ON public.notifications USING btree (user_id);


--
-- Name: idx_password_reset_tokens_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_password_reset_tokens_user_id ON public.password_reset_tokens USING btree (user_id);


--
-- Name: idx_payments_business_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payments_business_user_id ON public.payments USING btree (business_user_id);


--
-- Name: idx_payments_contract_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payments_contract_id ON public.payments USING btree (contract_id);


--
-- Name: idx_payments_influencer_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payments_influencer_user_id ON public.payments USING btree (influencer_user_id);


--
-- Name: idx_payments_paid_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payments_paid_at ON public.payments USING btree (paid_at);


--
-- Name: idx_payments_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payments_status ON public.payments USING btree (status);


--
-- Name: idx_payouts_influencer_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payouts_influencer_user_id ON public.payouts USING btree (influencer_user_id);


--
-- Name: idx_payouts_paid_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payouts_paid_at ON public.payouts USING btree (paid_at);


--
-- Name: idx_payouts_payment_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payouts_payment_id ON public.payouts USING btree (payment_id);


--
-- Name: idx_payouts_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_payouts_status ON public.payouts USING btree (status);


--
-- Name: idx_phone_verification_tokens_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_phone_verification_tokens_user_id ON public.phone_verification_tokens USING btree (user_id);


--
-- Name: idx_reviews_contract_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_reviews_contract_id ON public.reviews USING btree (contract_id);


--
-- Name: idx_reviews_rating; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_reviews_rating ON public.reviews USING btree (rating);


--
-- Name: idx_reviews_reviewed_user; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_reviews_reviewed_user ON public.reviews USING btree (reviewed_user_id);


--
-- Name: idx_reviews_reviewer; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_reviews_reviewer ON public.reviews USING btree (reviewer_user_id);


--
-- Name: idx_support_messages_sender; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_messages_sender ON public.support_messages USING btree (sender_user_id);


--
-- Name: idx_support_messages_ticket_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_messages_ticket_id ON public.support_messages USING btree (ticket_id);


--
-- Name: idx_support_tickets_assigned_admin; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_tickets_assigned_admin ON public.support_tickets USING btree (assigned_admin_user_id);


--
-- Name: idx_support_tickets_created_by; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_tickets_created_by ON public.support_tickets USING btree (created_by_user_id);


--
-- Name: idx_support_tickets_issue_type; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_tickets_issue_type ON public.support_tickets USING btree (issue_type);


--
-- Name: idx_support_tickets_related_campaign; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_tickets_related_campaign ON public.support_tickets USING btree (related_campaign_id);


--
-- Name: idx_support_tickets_related_contract; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_tickets_related_contract ON public.support_tickets USING btree (related_contract_id);


--
-- Name: idx_support_tickets_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_support_tickets_status ON public.support_tickets USING btree (status);


--
-- Name: idx_user_sessions_expires_at; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_user_sessions_expires_at ON public.user_sessions USING btree (expires_at);


--
-- Name: idx_user_sessions_refresh_token_hash; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE UNIQUE INDEX idx_user_sessions_refresh_token_hash ON public.user_sessions USING btree (refresh_token_hash);


--
-- Name: idx_user_sessions_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_user_sessions_user_id ON public.user_sessions USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: idx_users_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_users_status ON public.users USING btree (status);


--
-- Name: idx_verification_documents_status; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_verification_documents_status ON public.verification_documents USING btree (status);


--
-- Name: idx_verification_documents_user_id; Type: INDEX; Schema: public; Owner: adfluence
--

CREATE INDEX idx_verification_documents_user_id ON public.verification_documents USING btree (user_id);


--
-- Name: admin_profiles admin_profiles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.admin_profiles
    ADD CONSTRAINT admin_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: audit_logs audit_logs_actor_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_actor_user_id_fkey FOREIGN KEY (actor_user_id) REFERENCES public.users(id);


--
-- Name: audit_logs audit_logs_target_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_target_user_id_fkey FOREIGN KEY (target_user_id) REFERENCES public.users(id);


--
-- Name: bank_accounts bank_accounts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.bank_accounts
    ADD CONSTRAINT bank_accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: business_profiles business_profiles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.business_profiles
    ADD CONSTRAINT business_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: campaign_deliverables campaign_deliverables_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_deliverables
    ADD CONSTRAINT campaign_deliverables_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE;


--
-- Name: campaign_invites campaign_invites_business_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_invites
    ADD CONSTRAINT campaign_invites_business_user_id_fkey FOREIGN KEY (business_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: campaign_invites campaign_invites_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_invites
    ADD CONSTRAINT campaign_invites_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE;


--
-- Name: campaign_invites campaign_invites_influencer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_invites
    ADD CONSTRAINT campaign_invites_influencer_user_id_fkey FOREIGN KEY (influencer_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: campaign_platforms campaign_platforms_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_platforms
    ADD CONSTRAINT campaign_platforms_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE;


--
-- Name: campaign_recommendations campaign_recommendations_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_recommendations
    ADD CONSTRAINT campaign_recommendations_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE CASCADE;


--
-- Name: campaign_recommendations campaign_recommendations_influencer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaign_recommendations
    ADD CONSTRAINT campaign_recommendations_influencer_user_id_fkey FOREIGN KEY (influencer_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: campaigns campaigns_business_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_business_user_id_fkey FOREIGN KEY (business_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: campaigns campaigns_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: chat_admin_access_logs chat_admin_access_logs_admin_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.chat_admin_access_logs
    ADD CONSTRAINT chat_admin_access_logs_admin_user_id_fkey FOREIGN KEY (admin_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: chat_admin_access_logs chat_admin_access_logs_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.chat_admin_access_logs
    ADD CONSTRAINT chat_admin_access_logs_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;


--
-- Name: contract_deliverables contract_deliverables_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contract_deliverables
    ADD CONSTRAINT contract_deliverables_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.contracts(id) ON DELETE CASCADE;


--
-- Name: contract_status_history contract_status_history_changed_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contract_status_history
    ADD CONSTRAINT contract_status_history_changed_by_fkey FOREIGN KEY (changed_by) REFERENCES public.users(id);


--
-- Name: contract_status_history contract_status_history_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contract_status_history
    ADD CONSTRAINT contract_status_history_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.contracts(id) ON DELETE CASCADE;


--
-- Name: contracts contracts_business_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contracts
    ADD CONSTRAINT contracts_business_user_id_fkey FOREIGN KEY (business_user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: contracts contracts_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contracts
    ADD CONSTRAINT contracts_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE RESTRICT;


--
-- Name: contracts contracts_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contracts
    ADD CONSTRAINT contracts_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE SET NULL;


--
-- Name: contracts contracts_influencer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contracts
    ADD CONSTRAINT contracts_influencer_user_id_fkey FOREIGN KEY (influencer_user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: contracts contracts_invite_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.contracts
    ADD CONSTRAINT contracts_invite_id_fkey FOREIGN KEY (invite_id) REFERENCES public.campaign_invites(id) ON DELETE SET NULL;


--
-- Name: conversation_participants conversation_participants_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.conversation_participants
    ADD CONSTRAINT conversation_participants_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;


--
-- Name: conversation_participants conversation_participants_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.conversation_participants
    ADD CONSTRAINT conversation_participants_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: conversations conversations_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.conversations
    ADD CONSTRAINT conversations_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id) ON DELETE SET NULL;


--
-- Name: disputes disputes_against_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.disputes
    ADD CONSTRAINT disputes_against_user_id_fkey FOREIGN KEY (against_user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: disputes disputes_assigned_admin_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.disputes
    ADD CONSTRAINT disputes_assigned_admin_user_id_fkey FOREIGN KEY (assigned_admin_user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: disputes disputes_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.disputes
    ADD CONSTRAINT disputes_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.contracts(id) ON DELETE RESTRICT;


--
-- Name: disputes disputes_opened_by_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.disputes
    ADD CONSTRAINT disputes_opened_by_user_id_fkey FOREIGN KEY (opened_by_user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: disputes disputes_resolved_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.disputes
    ADD CONSTRAINT disputes_resolved_by_fkey FOREIGN KEY (resolved_by) REFERENCES public.users(id);


--
-- Name: email_verification_tokens email_verification_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT email_verification_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: influencer_portfolio_items influencer_portfolio_items_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_portfolio_items
    ADD CONSTRAINT influencer_portfolio_items_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: influencer_profile_categories influencer_profile_categories_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profile_categories
    ADD CONSTRAINT influencer_profile_categories_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: influencer_profile_categories influencer_profile_categories_influencer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profile_categories
    ADD CONSTRAINT influencer_profile_categories_influencer_user_id_fkey FOREIGN KEY (influencer_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: influencer_profile_platforms influencer_profile_platforms_influencer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profile_platforms
    ADD CONSTRAINT influencer_profile_platforms_influencer_user_id_fkey FOREIGN KEY (influencer_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: influencer_profiles influencer_profiles_approved_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profiles
    ADD CONSTRAINT influencer_profiles_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES public.users(id);


--
-- Name: influencer_profiles influencer_profiles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.influencer_profiles
    ADD CONSTRAINT influencer_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: message_reports message_reports_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.message_reports
    ADD CONSTRAINT message_reports_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;


--
-- Name: message_reports message_reports_reported_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.message_reports
    ADD CONSTRAINT message_reports_reported_user_id_fkey FOREIGN KEY (reported_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: message_reports message_reports_reporter_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.message_reports
    ADD CONSTRAINT message_reports_reporter_user_id_fkey FOREIGN KEY (reporter_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: message_reports message_reports_resolved_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.message_reports
    ADD CONSTRAINT message_reports_resolved_by_fkey FOREIGN KEY (resolved_by) REFERENCES public.users(id);


--
-- Name: messages messages_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_conversation_id_fkey FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;


--
-- Name: messages messages_receiver_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_receiver_user_id_fkey FOREIGN KEY (receiver_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: messages messages_sender_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_sender_user_id_fkey FOREIGN KEY (sender_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: notifications notifications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: password_reset_tokens password_reset_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: payments payments_business_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_business_user_id_fkey FOREIGN KEY (business_user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: payments payments_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.contracts(id) ON DELETE RESTRICT;


--
-- Name: payments payments_influencer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_influencer_user_id_fkey FOREIGN KEY (influencer_user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: payouts payouts_bank_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_bank_account_id_fkey FOREIGN KEY (bank_account_id) REFERENCES public.bank_accounts(id) ON DELETE SET NULL;


--
-- Name: payouts payouts_influencer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_influencer_user_id_fkey FOREIGN KEY (influencer_user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: payouts payouts_payment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payments(id) ON DELETE RESTRICT;


--
-- Name: phone_verification_tokens phone_verification_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.phone_verification_tokens
    ADD CONSTRAINT phone_verification_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: reviews reviews_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.contracts(id) ON DELETE CASCADE;


--
-- Name: reviews reviews_reviewed_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_reviewed_user_id_fkey FOREIGN KEY (reviewed_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: reviews reviews_reviewer_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_reviewer_user_id_fkey FOREIGN KEY (reviewer_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: support_messages support_messages_sender_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_messages
    ADD CONSTRAINT support_messages_sender_user_id_fkey FOREIGN KEY (sender_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: support_messages support_messages_ticket_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_messages
    ADD CONSTRAINT support_messages_ticket_id_fkey FOREIGN KEY (ticket_id) REFERENCES public.support_tickets(id) ON DELETE CASCADE;


--
-- Name: support_tickets support_tickets_assigned_admin_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_tickets
    ADD CONSTRAINT support_tickets_assigned_admin_user_id_fkey FOREIGN KEY (assigned_admin_user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: support_tickets support_tickets_created_by_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_tickets
    ADD CONSTRAINT support_tickets_created_by_user_id_fkey FOREIGN KEY (created_by_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: support_tickets support_tickets_related_campaign_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_tickets
    ADD CONSTRAINT support_tickets_related_campaign_id_fkey FOREIGN KEY (related_campaign_id) REFERENCES public.campaigns(id) ON DELETE SET NULL;


--
-- Name: support_tickets support_tickets_related_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_tickets
    ADD CONSTRAINT support_tickets_related_contract_id_fkey FOREIGN KEY (related_contract_id) REFERENCES public.contracts(id) ON DELETE SET NULL;


--
-- Name: support_tickets support_tickets_resolved_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.support_tickets
    ADD CONSTRAINT support_tickets_resolved_by_fkey FOREIGN KEY (resolved_by) REFERENCES public.users(id);


--
-- Name: user_sessions user_sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: verification_documents verification_documents_reviewed_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.verification_documents
    ADD CONSTRAINT verification_documents_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES public.users(id);


--
-- Name: verification_documents verification_documents_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: adfluence
--

ALTER TABLE ONLY public.verification_documents
    ADD CONSTRAINT verification_documents_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict F09d1iD4rin60UlBEf8T1cluR9yfUxtwjSl0g4eaGeJthmHdbcgBZpuoxBNqKYf

