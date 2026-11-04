-- +goose Up
-- Sports creators (athletes, sports journalists, football commentary) had no niche.
INSERT INTO categories (id, name, slug, description, is_active, created_at, updated_at)
VALUES ('c1000000-0000-4000-8000-000000000008', 'Sports', 'sports', 'Sports, athletes and sports journalism', TRUE, NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;

-- +goose Down
DELETE FROM categories c
WHERE c.slug = 'sports'
  AND NOT EXISTS (SELECT 1 FROM influencer_profile_categories ipc WHERE ipc.category_id = c.id)
  AND NOT EXISTS (SELECT 1 FROM campaigns cp WHERE cp.category_id = c.id);
