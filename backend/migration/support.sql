-- =============================================================================
-- SUPPORT.SQL — Support, disputes, reviews, notifications (5 tables)
-- =============================================================================
-- support_tickets
-- support_messages
-- disputes
-- reviews
-- notifications
-- =============================================================================

CREATE TABLE support_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    created_by_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    assigned_admin_user_id UUID
        REFERENCES users(id) ON DELETE SET NULL,

    related_campaign_id UUID
        REFERENCES campaigns(id) ON DELETE SET NULL,
    related_contract_id UUID
        REFERENCES contracts(id) ON DELETE SET NULL,

    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,

    issue_type support_issue_type NOT NULL,
    priority support_priority NOT NULL DEFAULT 'normal',
    status support_ticket_status NOT NULL DEFAULT 'open',

    resolution_note TEXT,
    resolved_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE support_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    ticket_id UUID NOT NULL
        REFERENCES support_tickets(id) ON DELETE CASCADE,
    sender_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    message TEXT NOT NULL,
    attachment_url TEXT,

    is_internal_note BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE RESTRICT,

    opened_by_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    against_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,
    assigned_admin_user_id UUID
        REFERENCES users(id) ON DELETE SET NULL,

    issue_type support_issue_type NOT NULL,
    description TEXT NOT NULL,

    severity dispute_severity NOT NULL DEFAULT 'medium',
    status dispute_status NOT NULL DEFAULT 'open',

    resolution_note TEXT,
    resolved_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    contract_id UUID NOT NULL
        REFERENCES contracts(id) ON DELETE CASCADE,

    reviewer_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    reviewed_user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (contract_id, reviewer_user_id, reviewed_user_id)
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,

    type notification_type NOT NULL,

    entity_type VARCHAR(100),
    entity_id UUID,

    read_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_support_tickets_created_by ON support_tickets(created_by_user_id);
CREATE INDEX idx_support_tickets_assigned_admin ON support_tickets(assigned_admin_user_id);
CREATE INDEX idx_support_tickets_status ON support_tickets(status);
CREATE INDEX idx_support_tickets_issue_type ON support_tickets(issue_type);
CREATE INDEX idx_support_tickets_related_campaign ON support_tickets(related_campaign_id);
CREATE INDEX idx_support_tickets_related_contract ON support_tickets(related_contract_id);

CREATE INDEX idx_support_messages_ticket_id ON support_messages(ticket_id);
CREATE INDEX idx_support_messages_sender ON support_messages(sender_user_id);

CREATE INDEX idx_disputes_contract_id ON disputes(contract_id);
CREATE INDEX idx_disputes_opened_by ON disputes(opened_by_user_id);
CREATE INDEX idx_disputes_against_user ON disputes(against_user_id);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_disputes_severity ON disputes(severity);

CREATE INDEX idx_reviews_contract_id ON reviews(contract_id);
CREATE INDEX idx_reviews_reviewer ON reviews(reviewer_user_id);
CREATE INDEX idx_reviews_reviewed_user ON reviews(reviewed_user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_read_at ON notifications(user_id, read_at);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
