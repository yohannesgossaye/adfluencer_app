-- Removes the demo creators added by seed_demo_creators.sql (emails ending
-- in .demo@example.com) and everything that belongs to them. Contracts block
-- the delete on purpose: cancel or remove those first.
BEGIN;
DELETE FROM campaign_recommendations WHERE influencer_user_id IN (SELECT id FROM users WHERE email LIKE '%.demo@example.com' AND role = 'influencer');
DELETE FROM users WHERE email LIKE '%.demo@example.com' AND role = 'influencer';
COMMIT;
