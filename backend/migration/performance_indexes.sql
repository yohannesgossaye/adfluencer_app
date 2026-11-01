-- Phase 3B: indexes for dashboards and lists. Indexes only; no data changes.
-- Safe to re-run.
BEGIN;

-- Dashboards and contract lists filter by owner and status together.
-- Each composite index also serves owner-only lookups, so it replaces the
-- single-column index on the same leading column.
CREATE INDEX IF NOT EXISTS idx_contracts_business_status
    ON contracts (business_user_id, status);
CREATE INDEX IF NOT EXISTS idx_contracts_influencer_status
    ON contracts (influencer_user_id, status);
DROP INDEX IF EXISTS idx_contracts_business_user_id;
DROP INDEX IF EXISTS idx_contracts_influencer_user_id;

-- Influencer invite list: newest pending invites for one creator.
CREATE INDEX IF NOT EXISTS idx_campaign_invites_influencer_status_created
    ON campaign_invites (influencer_user_id, status, created_at DESC);
DROP INDEX IF EXISTS idx_campaign_invites_influencer_user_id;

-- Business "recent campaigns": one owner, newest first.
CREATE INDEX IF NOT EXISTS idx_campaigns_business_updated
    ON campaigns (business_user_id, updated_at DESC);
DROP INDEX IF EXISTS idx_campaigns_business_user_id;

-- Payout totals on the influencer dashboard.
CREATE INDEX IF NOT EXISTS idx_payouts_influencer_status
    ON payouts (influencer_user_id, status);
DROP INDEX IF EXISTS idx_payouts_influencer_user_id;

-- These duplicate an index Postgres already keeps for a UNIQUE constraint:
-- they cost disk and slow every write, and are never chosen over the unique one.
DROP INDEX IF EXISTS idx_users_email;                    -- users.email UNIQUE
DROP INDEX IF EXISTS idx_business_profiles_user_id;      -- business_profiles.user_id UNIQUE
DROP INDEX IF EXISTS idx_influencer_profiles_user_id;    -- influencer_profiles.user_id UNIQUE
DROP INDEX IF EXISTS idx_categories_slug;                -- categories.slug UNIQUE
DROP INDEX IF EXISTS idx_campaign_platforms_campaign_id; -- UNIQUE (campaign_id, platform)

COMMIT;

ANALYZE contracts, campaign_invites, campaigns, payouts;
