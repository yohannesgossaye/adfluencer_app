-- Phase 2A business logic. Safe to re-run. Runs in one transaction.
--
-- 1. Business verification becomes an explicit workflow the business can see:
--    unsubmitted -> pending -> verified | rejected (with a reason) -> pending ...
--    is_verified stays in sync (TRUE only when verification_status = 'verified').
-- 2. users.previous_status remembers the status before an admin block, so
--    unblocking restores it instead of making everyone "active".
-- 3. Campaigns paused by the old admin "pause" action were stored as
--    'pending'; they are moved to 'paused', the status pause now uses.

BEGIN;

ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verification_status TEXT NOT NULL DEFAULT 'unsubmitted';
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verification_note TEXT;
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verification_submitted_at TIMESTAMPTZ;
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ;
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS verified_by UUID REFERENCES users(id) ON DELETE SET NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'business_profiles_verification_status_check'
    ) THEN
        ALTER TABLE business_profiles
            ADD CONSTRAINT business_profiles_verification_status_check
            CHECK (verification_status IN ('unsubmitted', 'pending', 'verified', 'rejected'));
    END IF;
END $$;

-- Backfill (only rows still at the default, so re-runs change nothing).
UPDATE business_profiles
SET verification_status = 'verified',
    verified_at = COALESCE(verified_at, updated_at)
WHERE is_verified = TRUE
  AND verification_status = 'unsubmitted';

UPDATE business_profiles bp
SET verification_status = 'pending',
    verification_submitted_at = COALESCE(bp.verification_submitted_at, u.updated_at)
FROM users u
WHERE u.id = bp.user_id
  AND u.status = 'active'
  AND bp.is_verified = FALSE
  AND bp.verification_status = 'unsubmitted';

CREATE INDEX IF NOT EXISTS idx_business_profiles_verification_queue
    ON business_profiles (verification_status, verification_submitted_at);

ALTER TABLE users ADD COLUMN IF NOT EXISTS previous_status user_status;

UPDATE campaigns SET status = 'paused', updated_at = NOW() WHERE status = 'pending';

COMMIT;
