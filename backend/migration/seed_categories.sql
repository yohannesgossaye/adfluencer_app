-- Seed categories for influencer niches (matches frontend profile.tsx).
-- Safe to re-run. For existing DB: make seed-categories

INSERT INTO categories (id, name, slug, description, is_active, created_at, updated_at)
VALUES
    ('c1000000-0000-4000-8000-000000000001', 'Beauty', 'beauty', 'Beauty and skincare creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000002', 'Fashion', 'fashion', 'Fashion and styling creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000003', 'Food', 'food', 'Food and restaurant creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000004', 'Tech', 'tech', 'Technology and gadget reviewers', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000005', 'Fitness', 'fitness', 'Fitness and wellness creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000006', 'Lifestyle', 'lifestyle', 'General lifestyle creators', TRUE, NOW(), NOW()),
    ('c1000000-0000-4000-8000-000000000007', 'Travel', 'travel', 'Travel and tourism creators', TRUE, NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;
