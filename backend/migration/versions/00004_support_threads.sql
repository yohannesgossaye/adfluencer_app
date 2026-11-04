-- +goose Up
-- Support tickets become conversations: every public reply is a row in
-- support_messages. Until now an admin's single reply lived in
-- support_tickets.resolution_note (and was overwritten by the next one);
-- copy those into the thread so no reply is lost.
INSERT INTO support_messages (ticket_id, sender_user_id, message, is_internal_note, created_at)
SELECT st.id, sender.id, st.resolution_note, FALSE, st.updated_at
FROM support_tickets st
CROSS JOIN LATERAL (
    SELECT COALESCE(st.assigned_admin_user_id, st.resolved_by,
                    (SELECT u.id FROM users u WHERE u.role = 'admin' ORDER BY u.created_at LIMIT 1)) AS id
) sender
WHERE btrim(COALESCE(st.resolution_note, '')) <> ''
  AND sender.id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM support_messages sm WHERE sm.ticket_id = st.id AND sm.is_internal_note = FALSE
  );

CREATE INDEX IF NOT EXISTS idx_support_messages_ticket_created
    ON support_messages (ticket_id, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_support_messages_ticket_created;
-- The copied replies stay: resolution_note still holds the same text.
