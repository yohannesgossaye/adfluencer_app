-- =============================================================================
-- MARKETPLACE.SQL — Chat, contracts, payments (9 tables)
-- =============================================================================
-- conversations
-- conversation_participants
-- message_reports
-- chat_admin_access_logs
-- contracts
-- contract_deliverables
-- contract_status_history
-- payments
-- payouts
-- =============================================================================
-- Firebase stores message bodies: conversations/{id}/messages/{messageId}
-- Link chat to contract: contracts.conversation_id -> conversations.id
-- =============================================================================

CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID REFERENCES campaigns(id) ON DELETE SET NULL,

    status conversation_status NOT NULL DEFAULT 'active',

    last_message_preview TEXT,
    last_message_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE conversation_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    conversation_id UUID NOT NULL
        REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    role conversation_participant_role NOT NULL,

    last_read_at TIMESTAMPTZ,
    is_muted BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (conversation_id, user_id)
);

CREATE TABLE message_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    conversation_id UUID NOT NULL
        REFERENCES conversations(id) ON DELETE CASCADE,
    firebase_message_id VARCHAR(255) NOT NULL,

    reporter_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    reported_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    reason TEXT NOT NULL,

    status message_report_status NOT NULL DEFAULT 'open',

    resolved_by UUID REFERENCES users(id),
    resolution_note TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE chat_admin_access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    admin_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    conversation_id UUID NOT NULL
        REFERENCES conversations(id) ON DELETE CASCADE,

    reason TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    campaign_id UUID NOT NULL
        REFERENCES campaigns(id) ON DELETE RESTRICT,
    business_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,

    conversation_id UUID REFERENCES conversations(id) ON DELETE SET NULL,
    invite_id UUID REFERENCES campaign_invites(id) ON DELETE SET NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,

    amount_etb NUMERIC(14,2) NOT NULL,

    business_fee_percent NUMERIC(5,2) NOT NULL DEFAULT 5.00,
    influencer_fee_percent NUMERIC(5,2) NOT NULL DEFAULT 5.00,

    business_fee_amount_etb NUMERIC(14,2) NOT NULL,
    influencer_fee_amount_etb NUMERIC(14,2) NOT NULL,
    total_business_payable_etb NUMERIC(14,2) NOT NULL,
    influencer_receivable_etb NUMERIC(14,2) NOT NULL,
    platform_total_fee_etb NUMERIC(14,2) NOT NULL,

    start_date DATE,
    end_date DATE,

    status contract_status NOT NULL DEFAULT 'draft',

    accepted_by_business_at TIMESTAMPTZ,
    accepted_by_influencer_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contract_deliverables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE CASCADE,

    deliverable_type deliverable_type NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    description TEXT,

    status contract_deliverable_status NOT NULL DEFAULT 'pending',

    due_date DATE,
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contract_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE CASCADE,

    old_status contract_status,
    new_status contract_status NOT NULL,

    changed_by UUID REFERENCES users(id),
    reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE RESTRICT,
    business_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,

    amount_etb NUMERIC(14,2) NOT NULL,
    business_fee_etb NUMERIC(14,2) NOT NULL,
    influencer_fee_etb NUMERIC(14,2) NOT NULL,
    platform_fee_total_etb NUMERIC(14,2) NOT NULL,
    total_paid_by_business_etb NUMERIC(14,2) NOT NULL,
    influencer_payout_etb NUMERIC(14,2) NOT NULL,

    payment_method VARCHAR(50) NOT NULL,
    provider VARCHAR(100),
    provider_reference VARCHAR(255),

    status payment_status NOT NULL DEFAULT 'pending',

    paid_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID NOT NULL
        REFERENCES payments(id) ON DELETE RESTRICT,
    influencer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    bank_account_id UUID
        REFERENCES bank_accounts(id) ON DELETE SET NULL,

    amount_etb NUMERIC(14,2) NOT NULL,

    status payout_status NOT NULL DEFAULT 'pending',

    provider VARCHAR(100),
    provider_reference VARCHAR(255),

    paid_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_conversations_campaign_id ON conversations(campaign_id);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC);

CREATE INDEX idx_conversation_participants_conversation ON conversation_participants(conversation_id);
CREATE INDEX idx_conversation_participants_user ON conversation_participants(user_id);

CREATE INDEX idx_message_reports_conversation ON message_reports(conversation_id);
CREATE INDEX idx_message_reports_status ON message_reports(status);
CREATE INDEX idx_message_reports_reporter ON message_reports(reporter_user_id);

CREATE INDEX idx_chat_admin_access_logs_admin ON chat_admin_access_logs(admin_user_id);
CREATE INDEX idx_chat_admin_access_logs_conversation ON chat_admin_access_logs(conversation_id);

CREATE INDEX idx_contracts_campaign_id ON contracts(campaign_id);
CREATE INDEX idx_contracts_business_user_id ON contracts(business_user_id);
CREATE INDEX idx_contracts_influencer_user_id ON contracts(influencer_user_id);
CREATE INDEX idx_contracts_conversation_id ON contracts(conversation_id);
CREATE INDEX idx_contracts_status ON contracts(status);
CREATE INDEX idx_contracts_created_at ON contracts(created_at);

CREATE INDEX idx_contract_deliverables_contract_id ON contract_deliverables(contract_id);
CREATE INDEX idx_contract_deliverables_status ON contract_deliverables(status);

CREATE INDEX idx_contract_status_history_contract_id ON contract_status_history(contract_id);
CREATE INDEX idx_contract_status_history_created_at ON contract_status_history(created_at);

CREATE INDEX idx_payments_contract_id ON payments(contract_id);
CREATE INDEX idx_payments_business_user_id ON payments(business_user_id);
CREATE INDEX idx_payments_influencer_user_id ON payments(influencer_user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_paid_at ON payments(paid_at);

CREATE INDEX idx_payouts_payment_id ON payouts(payment_id);
CREATE INDEX idx_payouts_influencer_user_id ON payouts(influencer_user_id);
CREATE INDEX idx_payouts_status ON payouts(status);
CREATE INDEX idx_payouts_paid_at ON payouts(paid_at);
