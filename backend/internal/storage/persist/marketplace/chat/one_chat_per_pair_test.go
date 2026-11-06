package chat

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"adinfluence_service_app/internal/testutil"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

type pairFixture struct{ business, creator, campaignA, campaignB string }

func newPair(t *testing.T, db *sql.DB) pairFixture {
	t.Helper()
	var f pairFixture
	q := func(dst *string, query string, args ...any) {
		t.Helper()
		if err := db.QueryRow(query, args...).Scan(dst); err != nil {
			t.Fatal(err)
		}
	}
	q(&f.business, `INSERT INTO users (email, password_hash, role, status) VALUES ('biz@example.com', 'x', 'business', 'active') RETURNING id`)
	q(&f.creator, `INSERT INTO users (email, password_hash, role, status) VALUES ('creator@example.com', 'x', 'influencer', 'active') RETURNING id`)
	q(&f.campaignA, `INSERT INTO campaigns (business_user_id, title, description) VALUES ($1, 'Bole boutique', 'x') RETURNING id`, f.business)
	q(&f.campaignB, `INSERT INTO campaigns (business_user_id, title, description) VALUES ($1, 'Holiday sale', 'x') RETURNING id`, f.business)
	return f
}

// Every way of opening a chat between the same business and creator lands
// in the same conversation, with or without a campaign.
func TestOneConversationPerPair(t *testing.T) {
	db := testutil.MigratedDB(t)
	p := &chatPersistence{db: db, log: utils.New(zap.NewNop())}
	ctx := context.Background()
	f := newPair(t, db)

	first, created, err := p.GetOrCreateConversation(ctx, f.business, f.creator, nil)
	if err != nil || !created {
		t.Fatalf("first chat: created=%v err=%v", created, err)
	}
	for _, c := range []struct {
		name     string
		from, to string
		campaign *string
	}{
		{"again without a campaign", f.business, f.creator, nil},
		{"from a campaign", f.business, f.creator, &f.campaignA},
		{"from another campaign, by the creator", f.creator, f.business, &f.campaignB},
	} {
		got, created, err := p.GetOrCreateConversation(ctx, c.from, c.to, c.campaign)
		if err != nil || created || got.ID != first.ID {
			t.Fatalf("%s: id=%v created=%v err=%v, want the first chat %s", c.name, got, created, err, first.ID)
		}
	}
	// A chat opened without a campaign takes the first campaign it is opened from, and keeps it.
	var campaign sql.NullString
	db.QueryRow(`SELECT campaign_id FROM conversations WHERE id = $1`, first.ID).Scan(&campaign)
	if campaign.String != f.campaignA {
		t.Fatalf("campaign label = %q, want campaign A", campaign.String)
	}
	var n int
	db.QueryRow(`SELECT count(*) FROM conversations`).Scan(&n)
	if n != 1 {
		t.Fatalf("%d conversations, want 1", n)
	}
}

// Migration 00006 merges a pair's duplicate chats into the oldest one
// without losing any message, and contracts follow the merged chat.
func TestDuplicateChatsAreMerged(t *testing.T) {
	db := testutil.MigratedDB(t)
	f := newPair(t, db)
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatal(err)
		}
	}
	// Recreate production before 00006 exactly: the old rule (one chat per
	// pair and campaign) in place, and a pair with a campaign chat plus a
	// campaign-less "Direct message" chat.
	exec(`DROP INDEX uq_conversations_pair`)
	exec(`CREATE UNIQUE INDEX uq_conversations_pair_campaign ON conversations (pair_low, pair_high,
		(COALESCE(campaign_id, '00000000-0000-0000-0000-000000000000'::uuid))) WHERE pair_low IS NOT NULL`)
	newChat := func(campaign any, age string) string {
		var id string
		if err := db.QueryRow(`INSERT INTO conversations (campaign_id, status, pair_low, pair_high, created_at)
			VALUES ($1, 'active', LEAST($2::uuid, $3::uuid), GREATEST($2::uuid, $3::uuid), NOW() - $4::interval) RETURNING id`,
			campaign, f.business, f.creator, age).Scan(&id); err != nil {
			t.Fatal(err)
		}
		exec(`INSERT INTO conversation_participants (conversation_id, user_id, role) VALUES ($1, $2, 'business'), ($1, $3, 'influencer')`, id, f.business, f.creator)
		return id
	}
	campaignChat := newChat(f.campaignA, "3 days")
	directChat := newChat(nil, "1 hour")
	msg := func(conv, from, to, body, age string) {
		exec(`INSERT INTO messages (conversation_id, sender_user_id, receiver_user_id, body, created_at) VALUES ($1, $2, $3, $4, NOW() - $5::interval)`, conv, from, to, body, age)
	}
	msg(campaignChat, f.business, f.creator, "Yes, I'm interested", "2 days")
	msg(campaignChat, f.creator, f.business, "Can we start this week?", "1 day")
	msg(directChat, f.business, f.creator, "how are u", "10 minutes")
	var contract string
	if err := db.QueryRow(`INSERT INTO contracts (campaign_id, business_user_id, influencer_user_id, title, amount_etb, business_fee_amount_etb,
		influencer_fee_amount_etb, total_business_payable_etb, influencer_receivable_etb, platform_total_fee_etb, conversation_id)
		VALUES ($1, $2, $3, 'Video', 1000, 50, 50, 1050, 950, 100, $4) RETURNING id`, f.campaignA, f.business, f.creator, directChat).Scan(&contract); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile("../../../../../migration/versions/00006_one_chat_per_pair.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(strings.Split(string(raw), "-- +goose Down")[0]); err != nil {
		t.Fatal(err)
	}

	var convs, msgs int
	var keeper, preview string
	db.QueryRow(`SELECT count(*) FROM conversations`).Scan(&convs)
	db.QueryRow(`SELECT id, last_message_preview FROM conversations`).Scan(&keeper, &preview)
	db.QueryRow(`SELECT count(*) FROM messages WHERE conversation_id = $1`, keeper).Scan(&msgs)
	if convs != 1 || keeper != campaignChat || msgs != 3 || preview != "how are u" {
		t.Fatalf("after merge: %d chats, keeper %s (want %s), %d messages, preview %q", convs, keeper, campaignChat, msgs, preview)
	}
	var contractChat string
	db.QueryRow(`SELECT conversation_id FROM contracts WHERE id = $1`, contract).Scan(&contractChat)
	if contractChat != campaignChat {
		t.Fatalf("contract chat = %s, want the merged chat", contractChat)
	}
	var participants int
	db.QueryRow(`SELECT count(*) FROM conversation_participants WHERE conversation_id = $1`, keeper).Scan(&participants)
	if participants != 2 {
		t.Fatalf("participants = %d, want 2", participants)
	}
	// The new rule is in place: a second chat for the pair is refused.
	if _, err := db.Exec(`INSERT INTO conversations (status, pair_low, pair_high) VALUES ('active', LEAST($1::uuid, $2::uuid), GREATEST($1::uuid, $2::uuid))`, f.business, f.creator); err == nil {
		t.Fatal("a second chat for the same pair was allowed")
	}
}

// The case that stopped production: the oldest chat has no campaign and a
// newer duplicate has one, with the old per-campaign rule still in place.
func TestMergeWhenOldestChatHasNoCampaign(t *testing.T) {
	db := testutil.MigratedDB(t)
	f := newPair(t, db)
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`DROP INDEX uq_conversations_pair`)
	exec(`CREATE UNIQUE INDEX uq_conversations_pair_campaign ON conversations (pair_low, pair_high,
		(COALESCE(campaign_id, '00000000-0000-0000-0000-000000000000'::uuid))) WHERE pair_low IS NOT NULL`)
	var direct, campaignChat string
	db.QueryRow(`INSERT INTO conversations (campaign_id, status, pair_low, pair_high, created_at)
		VALUES (NULL, 'active', LEAST($1::uuid, $2::uuid), GREATEST($1::uuid, $2::uuid), NOW() - interval '5 days') RETURNING id`, f.business, f.creator).Scan(&direct)
	db.QueryRow(`INSERT INTO conversations (campaign_id, status, pair_low, pair_high, created_at)
		VALUES ($3, 'active', LEAST($1::uuid, $2::uuid), GREATEST($1::uuid, $2::uuid), NOW() - interval '1 day') RETURNING id`, f.business, f.creator, f.campaignA).Scan(&campaignChat)
	exec(`INSERT INTO messages (conversation_id, sender_user_id, receiver_user_id, body) VALUES ($1, $2, $3, 'hello'), ($4, $3, $2, 'hi back')`,
		direct, f.business, f.creator, campaignChat)

	raw, err := os.ReadFile("../../../../../migration/versions/00006_one_chat_per_pair.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(strings.Split(string(raw), "-- +goose Down")[0]); err != nil {
		t.Fatalf("migration failed: %v", err)
	}
	var convs, msgs int
	var keeper string
	var campaign sql.NullString
	db.QueryRow(`SELECT count(*) FROM conversations`).Scan(&convs)
	db.QueryRow(`SELECT id, campaign_id FROM conversations`).Scan(&keeper, &campaign)
	db.QueryRow(`SELECT count(*) FROM messages WHERE conversation_id = $1`, keeper).Scan(&msgs)
	if convs != 1 || keeper != direct || msgs != 2 || campaign.String != f.campaignA {
		t.Fatalf("after merge: %d chats, keeper %s (want %s), %d messages, campaign %q (want A)", convs, keeper, direct, msgs, campaign.String)
	}
}
