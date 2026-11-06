-- Test data for the browser tests (e2e/*.spec.ts). Loaded into the
-- throwaway CI database only, after the migrations and the demo creators.
-- Every account's password is E2e-pass-2026.
BEGIN;

INSERT INTO users (id, email, password_hash, role, status, email_verified) VALUES
  ('e2e00000-0000-4000-8000-000000000001', 'e2e.business@example.com', '$argon2id$v=19$m=65536,t=4,p=2$Qen3CrgBgs3SNZkftDfoVQ$rqNAyBkZ1hPuU5U6SK7UC5tAfkhObVHyacyRsFxZND4', 'business', 'active', TRUE),
  ('e2e00000-0000-4000-8000-000000000002', 'e2e.creator@example.com', '$argon2id$v=19$m=65536,t=4,p=2$Qen3CrgBgs3SNZkftDfoVQ$rqNAyBkZ1hPuU5U6SK7UC5tAfkhObVHyacyRsFxZND4', 'influencer', 'active', TRUE),
  ('e2e00000-0000-4000-8000-000000000003', 'e2e.admin@example.com', '$argon2id$v=19$m=65536,t=4,p=2$Qen3CrgBgs3SNZkftDfoVQ$rqNAyBkZ1hPuU5U6SK7UC5tAfkhObVHyacyRsFxZND4', 'admin', 'active', TRUE);

INSERT INTO business_profiles (user_id, owner_full_name, business_name, is_verified, verification_status)
VALUES ('e2e00000-0000-4000-8000-000000000001', 'Test Owner', 'E2E Coffee', TRUE, 'verified');

INSERT INTO influencer_profiles (user_id, full_name, creator_name, username, bio, city, total_followers, tiktok_followers,
  engagement_rate, min_price_etb, max_price_etb, is_verified, approved_at)
VALUES ('e2e00000-0000-4000-8000-000000000002', 'Test Creator', 'E2E Creator', 'e2e_creator', 'Test creator account.',
  'Addis Ababa', 20000, 20000, 5.0, 3000, 8000, TRUE, NOW());
INSERT INTO influencer_profile_categories (influencer_user_id, category_id)
SELECT 'e2e00000-0000-4000-8000-000000000002', id FROM categories WHERE slug = 'food';
INSERT INTO influencer_profile_platforms (influencer_user_id, platform)
VALUES ('e2e00000-0000-4000-8000-000000000002', 'tiktok');

INSERT INTO admin_profiles (user_id, full_name) VALUES ('e2e00000-0000-4000-8000-000000000003', 'E2E Admin');

-- A live sports campaign to rank and invite for (find-creators.spec.ts).
INSERT INTO campaigns (id, business_user_id, title, description, category_id, budget_min_etb, budget_max_etb,
  creator_size_preference, number_of_creators_needed, status)
SELECT 'e2e00000-0000-4000-8000-0000000000c1', 'e2e00000-0000-4000-8000-000000000001', 'E2E campaign',
  'Football season launch', id, 20000, 90000, 'any', 2, 'active' FROM categories WHERE slug = 'sports';
INSERT INTO campaign_platforms (campaign_id, platform) VALUES ('e2e00000-0000-4000-8000-0000000000c1', 'youtube');

-- A live food campaign with a pending invite to the test creator (invite.spec.ts).
INSERT INTO campaigns (id, business_user_id, title, description, category_id, budget_min_etb, budget_max_etb,
  creator_size_preference, number_of_creators_needed, status)
SELECT 'e2e00000-0000-4000-8000-0000000000c2', 'e2e00000-0000-4000-8000-000000000001', 'E2E food tasting',
  'Cafe opening videos', id, 3000, 8000, 'any', 1, 'active' FROM categories WHERE slug = 'food';
INSERT INTO campaign_platforms (campaign_id, platform) VALUES ('e2e00000-0000-4000-8000-0000000000c2', 'tiktok');
INSERT INTO campaign_invites (id, campaign_id, business_user_id, influencer_user_id, status)
VALUES ('e2e00000-0000-4000-8000-0000000000a1', 'e2e00000-0000-4000-8000-0000000000c2',
  'e2e00000-0000-4000-8000-000000000001', 'e2e00000-0000-4000-8000-000000000002', 'pending');

-- A paid, in-progress contract to report a problem on (dispute.spec.ts).
INSERT INTO contracts (id, campaign_id, business_user_id, influencer_user_id, title, amount_etb, business_fee_percent,
  influencer_fee_percent, business_fee_amount_etb, influencer_fee_amount_etb, total_business_payable_etb,
  influencer_receivable_etb, platform_total_fee_etb, status, accepted_by_influencer_at)
VALUES ('e2e00000-0000-4000-8000-0000000000d1', 'e2e00000-0000-4000-8000-0000000000c2',
  'e2e00000-0000-4000-8000-000000000001', 'e2e00000-0000-4000-8000-000000000002', 'E2E tasting video',
  5000, 5, 5, 250, 250, 5250, 4750, 500, 'active', NOW());
INSERT INTO payments (contract_id, business_user_id, influencer_user_id, amount_etb, business_fee_etb, influencer_fee_etb,
  platform_fee_total_etb, total_paid_by_business_etb, influencer_payout_etb, payment_method, status, paid_at)
VALUES ('e2e00000-0000-4000-8000-0000000000d1', 'e2e00000-0000-4000-8000-000000000001',
  'e2e00000-0000-4000-8000-000000000002', 5000, 250, 250, 500, 5250, 4750, 'admin_confirmed', 'paid', NOW());

COMMIT;
