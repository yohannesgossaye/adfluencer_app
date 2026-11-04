-- +goose Up
-- Disputes opened from a contract ("Report a problem"): the reason the
-- person picked, and how the admin settled it.
ALTER TABLE disputes ADD COLUMN IF NOT EXISTS reason TEXT;
ALTER TABLE disputes ADD COLUMN IF NOT EXISTS outcome TEXT;
ALTER TABLE disputes DROP CONSTRAINT IF EXISTS disputes_reason_check;
ALTER TABLE disputes ADD CONSTRAINT disputes_reason_check CHECK (reason IS NULL OR reason IN (
    'not_delivered', 'delivered_late', 'not_as_agreed', 'creator_unresponsive',
    'business_unresponsive', 'extra_work_requested', 'approval_refused', 'other'
));
ALTER TABLE disputes DROP CONSTRAINT IF EXISTS disputes_outcome_check;
ALTER TABLE disputes ADD CONSTRAINT disputes_outcome_check CHECK (outcome IS NULL OR outcome IN ('continue', 'complete', 'terminate'));
-- At most one unresolved dispute per contract.
CREATE UNIQUE INDEX IF NOT EXISTS uq_disputes_one_open_per_contract
    ON disputes (contract_id) WHERE status IN ('open', 'under_review', 'escalated');

-- +goose Down
DROP INDEX IF EXISTS uq_disputes_one_open_per_contract;
ALTER TABLE disputes DROP CONSTRAINT IF EXISTS disputes_outcome_check;
ALTER TABLE disputes DROP CONSTRAINT IF EXISTS disputes_reason_check;
ALTER TABLE disputes DROP COLUMN IF EXISTS outcome;
ALTER TABLE disputes DROP COLUMN IF EXISTS reason;
