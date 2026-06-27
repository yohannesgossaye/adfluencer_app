# Adfluence Business Portal API Specification

Backend design for **business onboarding**, **profile CRUD**, **logo upload**, **campaigns**, **creator recommendations**, **invites**, and the path to **contracts**. Based on:

- **Frontend UI (mock today):** `/business/dashboard`, `/business/campaigns`, `/business/campaigns/create`, `/business/recommendations`, `/business/influencers/[id]`, `/business/contracts`, `/business/messages`
- **PostgreSQL schema:** `auth.sql` (`business_profiles`), `campaigns.sql`, `marketplace.sql` (contracts — later phase)
- **Auth:** completed — see `AUTH-API.md`
- **Media:** MinIO/S3 — see `MEDIA-STORAGE.md`
- **Influencer side:** profile + admin approval done — see `INFLUENCER-PROFILE-API.md`

**Frontend integration:** out of scope for this backend phase — APIs must work via **Postman/curl** first.

---

## 0. Critical product rules (do not forget)

### Platform revenue

| Rule | Detail |
|------|--------|
| **Fees apply on contracts only** | **5% from business** + **5% from influencer** on each `contracts.amount_etb` |
| **No subscriptions** | No SaaS tiers, no monthly platform plans |
| **No pricing packages API** | Frontend “Basic / Standard / Premium” on `/business/influencers/[id]` is **UI mock only** — do **not** implement |
| **Creator rates** | Show `influencer_profiles.min_price_etb` / `max_price_etb` as collaboration range — not platform fees |
| **Campaign budget** | `campaigns.budget_min_etb` / `budget_max_etb` = what business is willing to spend on creators — separate from platform fee |

Fee math (stored on `contracts` at creation time):

```
amount_etb                    = agreed collaboration amount (ETB)
business_fee_amount_etb       = amount × 5%
influencer_fee_amount_etb     = amount × 5%
total_business_payable_etb    = amount + business_fee_amount_etb
influencer_receivable_etb     = amount − influencer_fee_amount_etb
platform_total_fee_etb        = business_fee + influencer_fee
```

### Schema rule

Use **existing tables** unless explicitly approved. **36 tables today** (includes `influencer_portfolio_items`). Business Phase 1–3 uses tables already in `auth.sql` + `campaigns.sql` — **no new tables required** for profile + campaigns + invites.

### Business vs influencer onboarding

| | Business | Influencer |
|---|----------|------------|
| Signup creates | `users` + `business_profiles` | `users` + `influencer_profiles` |
| Profile complete | `users.status → active` | `users.status → pending_approval` |
| Admin verification | **Required** for marketplace visibility & contracts | **Required** for marketplace visibility |
| `is_verified` default | `false` until admin approves | `false` until admin approves application |
| Hidden until verified? | **Yes** — not visible to creators in marketplace | **Yes** — not in recommendations or public creator view |

### Marketplace visibility rules (both sides — comfort & trust)

| Action | Business requirement | Influencer requirement |
|--------|---------------------|------------------------|
| Appear in recommendations | `status=active` + `is_verified=true` | `status=active` + `is_verified=true` |
| Public creator profile (business view) | — | `status=active` + `is_verified=true` |
| Receive campaign invite | `status=active` (can invite while unverified V1) | `status=active` + `is_verified=true` |
| Accept invite | — | `status=active` + `is_verified=true` |
| Create contract | `status=active` + **`is_verified=true`** | Must be verified (already invited) |
| Sign contract | `is_verified=true` | `is_verified=true` |

**Phase 1 complete profile** sets `users.status = active` but leaves `business_profiles.is_verified = false`. Admin business verification API is a **future admin phase** (mirror influencer approval).

GET profile returns `marketplace_access` object so frontend can show “pending verification” banner.

---

## 1. Overview

### Purpose

After business signup (`users.status = pending_profile`), the owner completes **business profile**, optionally uploads a **logo**, then the account becomes **`active`**. Active businesses can **create campaigns**, **browse matched creators**, **send invites**, and later **create contracts** (separate phase).

### End-to-end lifecycle (target state)

```
SIGNUP (auth — role=business)
  users.status = pending_profile
  business_profiles row created (owner_full_name + business_name from signup name)

PHASE 1 — BUSINESS PROFILE
  GET  /business/profile
  PUT  /business/profile
  POST /business/profile/logo
  POST /business/profile/complete     → users.status = active (no admin gate V1)

PHASE 2 — CAMPAIGNS
  POST   /business/campaigns          → draft
  GET    /business/campaigns
  GET    /business/campaigns/{id}
  PUT    /business/campaigns/{id}
  POST   /business/campaigns/{id}/publish   → draft → active (or pending)
  POST   /business/campaigns/{id}/cancel

PHASE 3 — MATCHING & INVITES
  POST   /business/campaigns/{id}/recommendations/generate
  GET    /business/campaigns/{id}/recommendations
  GET    /business/influencers/{user_id}     → public creator card (active only)
  POST   /business/campaigns/{id}/invites
  GET    /business/campaigns/{id}/invites
  DELETE /business/campaigns/{id}/invites/{invite_id}   → cancel pending invite

PHASE 4 — INFLUENCER SIDE (separate doc section)
  GET  /influencer/invites
  POST /influencer/invites/{id}/accept
  POST /influencer/invites/{id}/reject

PHASE 5 — CONTRACTS (marketplace.sql — later)
  POST /business/contracts  (from accepted invite)
  … both sides accept → active → payments

PHASE 6 — DASHBOARD KPIs (last)
  GET /business/dashboard/summary
```

### Frontend routes (reference — integration later)

| Route | Current behavior | Backend replaces |
|-------|------------------|------------------|
| `/business/dashboard` | Mock KPIs + hardcoded name | Dashboard API (Phase 6); name from GET profile |
| `/business/campaigns` | Static campaign cards | GET `/business/campaigns` |
| `/business/campaigns/create` | localStorage draft | POST/PUT campaign |
| `/business/recommendations` | Mock scoring + localStorage invites | Recommendations + invite APIs |
| `/business/influencers/[id]` | Mock influencer data + package cards | GET public creator profile (**no packages**) |
| `/business/contracts` | Mock contracts | Contracts API (Phase 5) |
| `/business/messages` | Mock threads | Conversations API + Firebase (later) |

**Note:** There is **no `/business/profile` page yet** in the frontend. Backend still implements profile APIs so dashboard can show completion banner and create-campaign can pre-fill `business_name`.

### Implementation strategy

Same clean Go pattern as influencer:

```
DTO → persistence → module → handler → route
```

HTTP messages in **handler**, business rules in **module**, SQL in **persist**.

---

## 2. User status rules (business)

Maps to `user_status` enum. Enforced on every `/business/*` route.

| Status | Can view profile? | Can edit profile? | Can upload logo? | Can complete profile? | Can manage campaigns? |
|--------|-------------------|-------------------|------------------|----------------------|------------------------|
| `pending_profile` | yes | yes | yes | yes (if complete) | no |
| `pending_approval` | yes | yes | yes | no | no (unused for business V1) |
| `active` | yes | yes | yes | no | yes |
| `rejected` | yes | yes | yes | yes (resubmit) | no until active |
| `suspended` | yes (read) | no | no | no | no |
| `blocked` | no | no | no | no | no |
| `removed` | no | no | no | no | no |

**Role guard:** All `/business/*` routes require JWT `role = business`. Return **403** `forbidden` for influencer/admin tokens.

**Campaign guard:** All campaign/recommendation/invite routes require `users.status = active` and `business_profiles` row exists. Return **403** `profile_incomplete` or **409** `account_not_active` as appropriate.

---

## 3. Database

### 3.1 `business_profiles` (`auth.sql`) — Phase 1

| Column | Type | API usage |
|--------|------|-----------|
| `id` | UUID PK | Profile id |
| `user_id` | UUID UNIQUE | Owner |
| `owner_full_name` | VARCHAR(255) | Signup name; editable |
| `business_name` | VARCHAR(255) | Brand name (UI: create-campaign “Brand / business name”) |
| `legal_business_name` | VARCHAR(255) | Optional registered name |
| `logo_url` | TEXT | MinIO public URL after logo upload |
| `description` | TEXT | About the business |
| `industry` | VARCHAR(100) | e.g. Beauty, Tech — maps to niche/industry |
| `business_type` | VARCHAR(100) | e.g. Retail, SaaS, Restaurant |
| `business_size` | VARCHAR(100) | e.g. Solo, SME, Enterprise |
| `city` | VARCHAR(100) | Location |
| `region` | VARCHAR(100) | Optional |
| `address` | TEXT | Optional street address |
| `tin_number` | VARCHAR(100) | Tax ID (optional V1) |
| `website_url` | TEXT | |
| `instagram_url` | TEXT | |
| `tiktok_url` | TEXT | |
| `facebook_url` | TEXT | |
| `linkedin_url` | TEXT | |
| `preferred_platforms` | TEXT[] | e.g. `{tiktok,instagram}` |
| `monthly_marketing_budget` | NUMERIC(14,2) | Optional budget hint |
| `is_verified` | BOOLEAN | Admin/badge later; default false |
| `created_at` / `updated_at` | TIMESTAMPTZ | |

**Signup seed:** On business signup, auth already inserts `owner_full_name = name`, `business_name = name` (same as signup `name` field).

### 3.2 Campaigns (`campaigns.sql`) — Phase 2–3

#### `campaigns`

| Column | Type | API usage |
|--------|------|-----------|
| `id` | UUID PK | |
| `business_user_id` | UUID FK → users | Owner; must match JWT |
| `title` | VARCHAR(255) | Campaign title * |
| `goal` | TEXT | e.g. “Brand awareness” |
| `description` | TEXT | Long brief * |
| `category_id` | UUID FK → categories | Niche/industry match |
| `target_audience` | TEXT | e.g. “Women 18–30” |
| `city` | VARCHAR(100) | Location |
| `region` | VARCHAR(100) | Optional |
| `budget_min_etb` | NUMERIC(14,2) | Parsed from UI budget range |
| `budget_max_etb` | NUMERIC(14,2) | Parsed from UI budget range |
| `creator_size_preference` | creator_size enum | nano/micro/mid/macro/mega/any |
| `number_of_creators_needed` | INTEGER | Default 1 |
| `start_date` / `end_date` | DATE | Timeline |
| `status` | campaign_status enum | draft → active → completed/cancelled |
| `created_at` / `updated_at` | TIMESTAMPTZ | |

**Status enum:** `draft`, `pending`, `active`, `completed`, `cancelled`, `paused`

**V1 transitions:**
- Create → `draft`
- Publish → `active` (or `pending` if you add admin campaign review later — **not in V1**)
- Cancel → `cancelled`
- Auto `completed` when `end_date` passed — optional cron later

#### `campaign_platforms`

| Column | Usage |
|--------|--------|
| `campaign_id` | Parent |
| `platform` | `social_platform` enum: tiktok, instagram, youtube, telegram, facebook |

Replace on each campaign PUT (same pattern as influencer platforms).

#### `campaign_deliverables`

| Column | Usage |
|--------|--------|
| `campaign_id` | Parent |
| `deliverable_type` | See enum mapping below |
| `quantity` | Default 1 |
| `description` | Optional note |

#### `campaign_recommendations`

| Column | Usage |
|--------|--------|
| `campaign_id` | |
| `influencer_user_id` | Matched creator |
| `match_score` | 0–100 |
| `score_breakdown` | JSONB `{niche, platform, budget, audience, verified}` |
| `rank_position` | 1, 2, 3… |
| `status` | recommended, invited, dismissed |

Regenerated on `POST .../recommendations/generate` (delete old `recommended` rows or upsert).

#### `campaign_invites`

| Column | Usage |
|--------|--------|
| `campaign_id` | |
| `business_user_id` | Sender |
| `influencer_user_id` | Recipient |
| `message` | Optional invite note |
| `status` | pending, accepted, rejected, cancelled |
| `responded_at` | Set when influencer responds |

**Unique:** `(campaign_id, influencer_user_id)` — one invite per creator per campaign.

### 3.3 Tables used later (not Phase 1–3)

| Table | Phase |
|-------|-------|
| `contracts`, `contract_deliverables`, `contract_status_history` | Phase 5 |
| `conversations`, `conversation_participants` | Messages |
| `payments`, `payouts` | Chapa / payout integration |
| `reviews` | Post-contract ratings |
| `notifications` | Trigger on invite/contract events |
| `verification_documents` | Business license upload (optional later) |

### 3.4 Dependency on influencer data

Recommendations and public creator view only include influencers where:

```sql
users.role = 'influencer'
AND users.status = 'active'
AND influencer_profiles.is_verified = true
```

Portfolio samples from `influencer_portfolio_items`. Avatar from `profile_image_url`.

---

## 4. Enum & UI mapping

### 4.1 Campaign goal (free text V1)

Frontend dropdown values — store as `campaigns.goal` text:

- Brand awareness
- Product launch
- App downloads
- Sales conversion
- Event promotion

### 4.2 Category / niche

Frontend niches map to `categories` table (seed: `make seed-categories`):

| UI label | Resolve by |
|----------|------------|
| Beauty | `categories.name` or `slug` |
| Fashion | … |
| Food, Tech, Fitness, Lifestyle, Travel | … |

Same resolver pattern as influencer `ResolveCategoryIDs`.

### 4.3 Location (create-campaign UI)

| UI option | Backend mapping |
|-----------|-----------------|
| Addis Ababa | `city = 'Addis Ababa'` |
| Ethiopia wide | `city = 'Addis Ababa'`, `region = 'Ethiopia wide'` (or dedicated region value) |
| Regional campaign | `region = 'Regional campaign'`, city optional |

**V1:** Accept `city` + optional `region` strings from client; do not over-normalize.

### 4.4 Platforms

| UI | DB enum |
|----|---------|
| TikTok | `tiktok` |
| Instagram | `instagram` |
| YouTube | `youtube` |

### 4.5 Budget range (UI) → numeric

| UI `budgetRange` | `budget_min_etb` | `budget_max_etb` |
|------------------|------------------|------------------|
| Under 10,000 ETB | 0 | 10000 |
| 10,000–25,000 ETB | 10000 | 25000 |
| 25,000–50,000 ETB | 25000 | 50000 |
| 50,000+ ETB | 50000 | 150000 |

API accepts either:
- `budget_range` string (server parses), **or**
- explicit `budget_min_etb` + `budget_max_etb` (preferred for API clarity)

### 4.6 Creator size (UI) → `creator_size` enum

| UI | DB enum | Follower heuristic (matching) |
|----|---------|--------------------------------|
| Micro creator | `micro` | &lt; 70,000 total_followers |
| Mid-tier creator | `mid` | 70,000 – 109,999 |
| Large creator | `macro` | ≥ 110,000 |
| (any) | `any` | no filter |

Schema also has `nano`, `mega` — available for future; UI V1 uses micro/mid/macro only.

### 4.7 Deliverables (UI) → `deliverable_type` enum

| UI checkbox | DB enum |
|-------------|---------|
| TikTok video | `tiktok_video` |
| Instagram reel | `instagram_reel` |
| Instagram story | `instagram_story` |
| YouTube short | `youtube_short` |
| Product review | `product_review` |

### 4.8 Campaign list UI status → DB

| UI card status | `campaigns.status` |
|----------------|-------------------|
| Pending | `draft` or `pending` |
| Active | `active` |
| Completed | `completed` |

**Progress bar (72%, etc.):** Derived in dashboard phase from contracts/deliverables — **not stored on campaign V1**. List API may return `progress_percent: null` until Phase 6.

---

## 5. Phase 1 — Business profile API

### 5.1 GET profile

**`GET /api/v1/business/profile`**

```
Authorization: Bearer {access_token}
```

**Response `data`:**

```json
{
  "user": {
    "id": "uuid",
    "email": "brand@example.com",
    "phone_number": "+251...",
    "status": "pending_profile",
    "status_reason": null,
    "email_verified": true,
    "phone_verified": false
  },
  "profile": {
    "id": "uuid",
    "owner_full_name": "Yohannes Bekele",
    "business_name": "Betu Cosmetics",
    "legal_business_name": null,
    "logo_url": null,
    "description": null,
    "industry": null,
    "business_type": null,
    "business_size": null,
    "city": null,
    "region": null,
    "address": null,
    "tin_number": null,
    "website_url": null,
    "social_links": {
      "instagram": "",
      "tiktok": "",
      "facebook": "",
      "linkedin": ""
    },
    "preferred_platforms": [],
    "monthly_marketing_budget": null,
    "is_verified": false
  },
  "completeness": {
    "percent": 45,
    "missing": ["description", "city", "industry"],
    "can_complete": false
  },
  "marketplace_access": {
    "visible_in_marketplace": false,
    "can_create_contracts": false,
    "message": "Complete your business profile to activate your account."
  }
}
```

### 5.2 PUT profile

**`PUT /api/v1/business/profile`**

```json
{
  "owner_full_name": "Yohannes Bekele",
  "business_name": "Betu Cosmetics",
  "legal_business_name": "",
  "description": "Ethiopian beauty brand focused on ...",
  "industry": "Beauty",
  "business_type": "Retail",
  "business_size": "SME",
  "city": "Addis Ababa",
  "region": "",
  "address": "",
  "tin_number": "",
  "website_url": "https://...",
  "social_links": {
    "instagram": "https://instagram.com/...",
    "tiktok": "",
    "facebook": "",
    "linkedin": ""
  },
  "preferred_platforms": ["instagram", "tiktok"],
  "monthly_marketing_budget": 50000
}
```

**Validation (V1 minimum for `can_complete`):**

| Field | Rule |
|-------|------|
| `owner_full_name` | required, 2–255 chars |
| `business_name` | required, 2–255 chars |
| `description` | min 20 chars |
| `industry` | required (must resolve to category name or free text in `industry` column) |
| `city` | required |
| `preferred_platforms` | at least 1 |

Returns same shape as GET.

### 5.3 POST logo upload

**`POST /api/v1/business/profile/logo`**

```
Content-Type: multipart/form-data
file: (binary)
```

Same rules as influencer avatar:
- JPEG, PNG, WebP — max **5 MB**
- MinIO key: `media/business/{user_id}/logo/{uuid}.{ext}`
- Updates `business_profiles.logo_url`
- Deletes previous object (best effort)

**Response:**

```json
{
  "logo_url": "http://localhost:9000/adfluence-media/media/business/.../logo/....png"
}
```

### 5.4 POST complete profile (activate account)

**`POST /api/v1/business/profile/complete`**

Unlike influencer, **no admin gate on status** — but **`is_verified` stays false** until admin verifies the business.

**Preconditions:**
- `completeness.can_complete === true`
- `users.status` in (`pending_profile`, `rejected`)

**Effect:**
- `users.status → active`
- `business_profiles.is_verified` unchanged (still `false`)
- Audit log: profile completed, pending admin verification

**Response:**

```json
{
  "status": "active",
  "message": "Business profile completed. Admin verification is required before marketplace visibility and contracts.",
  "is_verified": false,
  "completed_at": "2026-06-27T12:00:00Z"
}
```

**Note:** Logo optional for complete in V1 (recommended in completeness scoring but not blocking).

### 5.5 Profile completeness scoring

| Field | Points |
|-------|--------|
| owner_full_name | 10 |
| business_name | 10 |
| description (≥20 chars) | 15 |
| industry | 10 |
| city | 10 |
| preferred_platforms (≥1) | 10 |
| logo_url | 10 |
| website or any social link | 5 |
| business_type | 5 |
| monthly_marketing_budget set | 5 |

Max 100. `can_complete` when required fields pass (not necessarily 100%).

---

## 6. Phase 2 — Campaign APIs

Base path: `/api/v1/business/campaigns`  
Requires: `role=business`, `status=active`.

### 6.1 Create campaign (draft)

**`POST /api/v1/business/campaigns`**

```json
{
  "title": "Betu Cosmetics Launch",
  "goal": "Product launch",
  "description": "We are launching a new skincare line ...",
  "niche": "Beauty",
  "target_audience": "Women 18-30 in urban areas",
  "city": "Addis Ababa",
  "region": "",
  "platforms": ["instagram", "tiktok"],
  "budget_min_etb": 25000,
  "budget_max_etb": 50000,
  "creator_size_preference": "mid",
  "number_of_creators_needed": 3,
  "deliverables": [
    { "type": "instagram_reel", "quantity": 2 },
    { "type": "instagram_story", "quantity": 3 }
  ],
  "start_date": "2026-07-01",
  "end_date": "2026-08-01",
  "additional_notes": "Optional — append to description or store in goal/description"
}
```

**Notes:**
- `business_name` on create-campaign UI comes from **business profile**, not campaign payload (pre-fill on frontend).
- `additional_notes` → concatenate into `description` or ignore in V1.
- Creates `campaigns` row with `status = draft` + platforms + deliverables in transaction.

### 6.2 List campaigns

**`GET /api/v1/business/campaigns?status=&page=1&limit=20`**

**Response item:**

```json
{
  "id": "uuid",
  "title": "Betu Cosmetics Launch",
  "business_name": "Betu Cosmetics",
  "primary_platform": "instagram",
  "budget_summary": "ETB 25,000 – 50,000",
  "status": "active",
  "start_date": "2026-07-01",
  "end_date": "2026-08-01",
  "creators_needed": 3,
  "invites_sent": 2,
  "invites_accepted": 1,
  "created_at": "..."
}
```

### 6.3 Get campaign detail

**`GET /api/v1/business/campaigns/{campaign_id}`**

Full campaign + platforms + deliverables + invite counts.

### 6.4 Update campaign

**`PUT /api/v1/business/campaigns/{campaign_id}`**

Only allowed when `status` in (`draft`, `active`, `paused`).  
Cannot change `business_user_id`.  
Replacing platforms/deliverables same as influencer category pattern.

### 6.5 Publish campaign

**`POST /api/v1/business/campaigns/{campaign_id}/publish`**

Validates required fields (same as create-campaign UI validation).  
`draft → active`.  
Returns updated campaign.

### 6.6 Cancel campaign

**`POST /api/v1/business/campaigns/{campaign_id}/cancel`**

`active|draft|paused → cancelled`.  
Pending invites → `cancelled`.

---

## 7. Phase 3 — Recommendations & invites

### 7.1 Generate recommendations

**`POST /api/v1/business/campaigns/{campaign_id}/recommendations/generate`**

**Algorithm V1 (mirror frontend mock logic):**

For each eligible influencer, compute `match_score` 0–100:

| Factor | Max points | Logic |
|--------|------------|--------|
| Niche fit | 40 | `category_id` matches any influencer niche |
| Platform fit | 25 | overlap campaign platforms ∩ influencer platforms |
| Budget fit | 20 | influencer `min_price_etb`/`max_price_etb` overlaps campaign budget band |
| Creator size | 10 | `total_followers` matches `creator_size_preference` |
| Verified | 5 | `is_verified = true` |

Store in `campaign_recommendations` with `score_breakdown` JSON.  
Sort desc, assign `rank_position`.  
Return top **50** (paginate in GET).

**Minimum score to list:** 78 (same as frontend fallback).

### 7.2 List recommendations

**`GET /api/v1/business/campaigns/{campaign_id}/recommendations?search=&niche=&platform=&min_score=&sort=match&page=1&limit=20`**

**Item shape (for recommendation cards):**

```json
{
  "influencer_user_id": "uuid",
  "creator_name": "Biruk Alemu",
  "profile_image_url": "...",
  "primary_niche": "Tech",
  "platforms": ["youtube"],
  "total_followers": 84000,
  "engagement_rate": 7.2,
  "min_price_etb": 15000,
  "max_price_etb": 35000,
  "match_score": 91,
  "score_breakdown": {
    "niche": 40,
    "platform": 25,
    "budget": 20,
    "creator_size": 10,
    "verified": 5
  },
  "is_verified": true,
  "invite_status": null
}
```

### 7.3 Public creator profile (business view)

**`GET /api/v1/business/influencers/{user_id}?campaign_id={optional}`**

For `/business/influencers/[id]` page.

**Include:**
- Creator profile fields (name, bio, city, followers, engagement, prices, niches, platforms, social links)
- `portfolio_items[]` (from `influencer_portfolio_items`)
- If `campaign_id` provided: campaign-specific `match_score` + breakdown

**Exclude / do not implement:**
- Pricing package cards (Basic/Standard/Premium)
- Mock reviews (use `reviews` table in later phase)
- Fake performance metrics (avg reach/views) unless computed from real data later

**Access:** Only `active` + verified influencers. Others → **404** `creator_not_found`.

### 7.4 Send invite

**`POST /api/v1/business/campaigns/{campaign_id}/invites`**

```json
{
  "influencer_user_id": "uuid",
  "message": "We'd love to collaborate on our summer launch."
}
```

**Rules:**
- Campaign must be `active`
- Influencer must be active + verified
- Unique per (campaign, influencer)
- Sets recommendation status → `invited` if row exists
- Creates notification row (optional V1 — log only)

**Response:**

```json
{
  "id": "invite_uuid",
  "status": "pending",
  "created_at": "..."
}
```

### 7.5 List invites for campaign

**`GET /api/v1/business/campaigns/{campaign_id}/invites`**

### 7.6 Cancel invite

**`DELETE /api/v1/business/campaigns/{campaign_id}/invites/{invite_id}`**

Only `pending` invites. Status → `cancelled`.

---

## 8. Phase 4 — Influencer invite APIs (cross-reference)

Implemented under `/api/v1/influencer/invites` (separate handler package, same DB).

| Method | Path | Action |
|--------|------|--------|
| GET | `/influencer/invites?status=` | List invites for logged-in influencer |
| GET | `/influencer/invites/{id}` | Detail + campaign summary |
| POST | `/influencer/invites/{id}/accept` | `pending → accepted` |
| POST | `/influencer/invites/{id}/reject` | `pending → rejected` |

After accept → business can create contract (Phase 5).

---

## 9. Phase 5 — Contracts preview (spec only — implement later)

See `marketplace.sql`. Entry point after **accepted invite**:

**`POST /api/v1/business/contracts`**

```json
{
  "invite_id": "uuid",
  "title": "Betu Cosmetics Launch — Biruk Alemu",
  "amount_etb": 30000,
  "description": "...",
  "deliverables": [...],
  "start_date": "2026-07-01",
  "end_date": "2026-07-31"
}
```

Server computes fee columns from defaults 5%/5%.  
Creates `contracts` + optional `conversation` + `contract_deliverables`.

Influencer accepts contract → both sides → `active`.

**Do not start contracts until Phases 1–3 work in Postman.**

---

## 10. Code structure (Go)

### 10.1 New packages (mirror influencer)

```
internal/constant/model/dto/business/
internal/storage/persist/business/
internal/storage/business_storage.go
internal/module/business/
internal/module/business/helper/          # completeness, budget parse, matching
internal/handler/rest/http/business/
internal/glue/routing/business/
```

Campaign + invite logic can live in:
- `internal/module/campaign/` **or**
- nested under `business` module as `business/campaign.go`

**Recommendation:** Single `business` module for profile + campaigns + invites in V1; split later if files grow.

### 10.2 Initiator wiring

```go
// initiator/persistence.go — add business persist
// initiator/module.go     — InitBusinessModule(...)
// initiator/handler.go  — businessHTTP.InitBusinessHandler(...)
// initiator/route.go    — businessRoutes.BusinessRoutes(...)
```

Pass `platformStorage.Service` for logo upload (same as influencer).

### 10.3 Storage keys (MinIO)

Add to `platform/storage/keys.go`:

```go
func BusinessLogoObjectKey(userID, ext string) string
// media/business/{user_id}/logo/{uuid}.{ext}
```

---

## 11. Errors

| HTTP | Code | When |
|------|------|------|
| 400 | `invalid_request` | Malformed JSON |
| 400 | `validation_error` | Field validation |
| 400 | `invalid_file` | Bad upload |
| 401 | `unauthorized` | Missing/invalid JWT |
| 403 | `forbidden` | Wrong role |
| 403 | `profile_locked` | Suspended/blocked |
| 403 | `account_not_active` | Campaign API while not active |
| 404 | `profile_not_found` | No business_profiles row |
| 404 | `campaign_not_found` | |
| 404 | `creator_not_found` | Influencer not visible |
| 404 | `invite_not_found` | |
| 409 | `profile_incomplete` | Complete profile failed |
| 409 | `invalid_status_for_complete` | Wrong user status |
| 409 | `invalid_campaign_status` | e.g. publish when cancelled |
| 409 | `invite_already_exists` | Duplicate invite |
| 409 | `invite_not_pending` | Cancel/respond wrong state |
| 503 | `storage_not_configured` | MinIO missing |

Add to `internal/constant/errors/errors.go` + `ErrorMap` + `ClientErrorCode`.

---

## 12. Security & authorization

| Rule | Enforcement |
|------|-------------|
| Business A cannot read Business B campaigns | All queries filter `business_user_id = JWT user_id` |
| Business cannot invite self | Role check on influencer id |
| Business cannot view draft influencer profiles | Only active + verified |
| Campaign id in URL must belong to caller | Module-level ownership check |
| Admin portal | Separate `/admin/*` routes (campaign moderation optional later) |

---

## 13. Postman collection (planned)

Folder: **Business → Profile**

| Request | Method | URL |
|---------|--------|-----|
| Get Profile | GET | `/api/v1/business/profile` |
| Update Profile | PUT | `/api/v1/business/profile` |
| Upload Logo | POST | `/api/v1/business/profile/logo` |
| Complete Profile | POST | `/api/v1/business/profile/complete` |

Folder: **Business → Campaigns**

| Request | Method | URL |
|---------|--------|-----|
| Create Campaign | POST | `/api/v1/business/campaigns` |
| List Campaigns | GET | `/api/v1/business/campaigns` |
| Get Campaign | GET | `/api/v1/business/campaigns/{{campaignId}}` |
| Update Campaign | PUT | `/api/v1/business/campaigns/{{campaignId}}` |
| Publish Campaign | POST | `/api/v1/business/campaigns/{{campaignId}}/publish` |
| Cancel Campaign | POST | `/api/v1/business/campaigns/{{campaignId}}/cancel` |

Folder: **Business → Recommendations & Invites**

| Request | Method | URL |
|---------|--------|-----|
| Generate Recommendations | POST | `/api/v1/business/campaigns/{{campaignId}}/recommendations/generate` |
| List Recommendations | GET | `/api/v1/business/campaigns/{{campaignId}}/recommendations` |
| View Creator | GET | `/api/v1/business/influencers/{{influencerUserId}}` |
| Send Invite | POST | `/api/v1/business/campaigns/{{campaignId}}/invites` |
| List Invites | GET | `/api/v1/business/campaigns/{{campaignId}}/invites` |
| Cancel Invite | DELETE | `/api/v1/business/campaigns/{{campaignId}}/invites/{{inviteId}}` |

Collection variables: `campaignId`, `inviteId`, `influencerUserId`.

---

## 14. Implementation checklist

### Phase 1 — Business profile (start here)

- [x] DTOs: `GetProfileResponse`, `UpdateProfileRequest`, `UploadLogoResponse`, `CompleteProfileResponse`
- [x] Completeness helper (`module/business/helper/completeness.go`)
- [x] Persistence: get/update profile, update logo URL, complete profile (update users.status)
- [x] Module + handler + routes under `/business/profile`
- [x] `BusinessLogoObjectKey` + logo upload (reuse storage validate)
- [x] Wire initiator (persistence, module, handler, routes)
- [x] Postman folder **Business → Profile**
- [x] Errors registered
- [x] `marketplace_access` in GET profile (verification gating preview)
- [ ] Admin business verification API (sets `is_verified=true`) — future admin phase

### Phase 2 — Campaigns

- [x] Campaign DTOs + validation
- [x] CRUD persistence (transaction for platforms/deliverables)
- [x] Publish / cancel status transitions
- [x] Module + handler + routes under `/business/campaigns`
- [x] Postman **Business → Campaigns**
- [x] Errors registered (`campaign_not_found`, `invalid_campaign_status`, `account_not_active`)

### Phase 3 — Recommendations & invites

- [x] Matching helper (score + breakdown)
- [x] Generate + list recommendations (`campaign_recommendations` table)
- [x] Public creator GET for business
- [x] Invite create/list/cancel (`campaign_invites` table)
- [x] Module + handler + routes
- [x] Postman **Recommendations & Invites**
- [x] Errors registered (`creator_not_found`, `invite_not_found`, `invite_already_exists`, `invalid_invite_status`)

### Phase 4 — Influencer invites

- [x] `/influencer/invites` list/detail/accept/reject
- [x] Persistence on `campaign_invites` (influencer-scoped queries)
- [x] Module + handler + routes
- [x] Postman **Influencer → Invites**
- [x] Accept requires `active` + `is_verified`; campaign must be `active`

### Phase 5 — Contracts (Phase 5A — no Chapa yet)

- [x] Fee calculation helper (5% business + 5% influencer)
- [x] Create contract from accepted invite (`contracts`, `contract_deliverables`, `contract_status_history`)
- [x] Business list/get/pay (escrow stub)/approve-work/cancel
- [x] Influencer list/get/accept/reject
- [x] Payments + payouts records on pay/approve
- [x] Postman **Business → Contracts** + **Influencer → Contracts**
- [ ] Chapa payment gateway + webhooks (Phase 5B)
- [ ] Real bank/Telebirr payout transfer (Phase 5B)

### Phase 6 — Dashboard

- [x] `GET /business/dashboard/summary` — KPIs from campaigns/contracts/recommendations/invites
- [x] Recent campaigns table (5 rows, progress from accepted invites)
- [x] Top recommendations strip (5 rows, match_score ≥ 78)
- [x] Admin business verification — list/get/verify/revoke (`/admin/businesses/verifications`)
- [x] Postman **Business → Dashboard** + **Admin → Business Verifications**

### Phase 7 — Support (business + admin)

- [x] Business: create/list/get tickets, campaign dropdown
- [x] Admin: list (filters + stats), get detail, update status, public reply + internal note
- [x] Uses existing `support_tickets` + `support_messages` tables (no new tables)
- [x] Postman **Business → Support** + **Admin → Support**

---

## 15. Dashboard API

### `GET /api/v1/business/dashboard/summary`

**Auth:** Business bearer token, account `active`.

Returns real counts only — no reach, engagement, or performance chart data (frontend keeps those mocked).

**Response highlights:**

| Field | Source |
|-------|--------|
| `greeting.business_name`, `greeting.owner_full_name` | `business_profiles` |
| `marketplace_access` | Same rules as GET profile |
| `kpis.active_campaigns` | `campaigns.status = active` |
| `kpis.recommended_creators_total` | `campaign_recommendations.match_score >= 78` |
| `kpis.recommended_creators_high_match` | `match_score >= 85` |
| `kpis.pending_contracts` | `pending_influencer` + `accepted` |
| `kpis.contracts_awaiting_influencer` | `pending_influencer` |
| `kpis.contracts_awaiting_payment` | `accepted` (awaiting pay) |
| `contracts_summary` | pending / active / completed contract counts |
| `recent_campaigns` | Last 5 campaigns with niche, budget, progress |
| `top_recommendations` | Top 5 by match_score across business campaigns |

**SQL strategy:** One scalar-subquery row for all KPIs (no cross joins), separate lightweight queries for recent campaigns and top recommendations.

---

## 18. Support API (business + admin)

Uses `support_tickets` and `support_messages` (internal notes only in V1).

### Business

| Method | Path | Description |
|--------|------|-------------|
| GET | `/business/support/campaign-options` | Campaign dropdown for related campaign |
| POST | `/business/support/tickets` | Submit ticket |
| GET | `/business/support/tickets` | My requests (includes `admin_reply`) |
| GET | `/business/support/tickets/{ticket_id}` | Ticket detail |

**Create body:**

```json
{
  "issue_type": "payment_question",
  "subject": "Contract payment question",
  "message": "Describe the issue…",
  "priority": "medium",
  "related_campaign_id": "uuid-or-null"
}
```

**Issue types (business):** `campaign_issue`, `contract_issue`, `payment_question`, `creator_communication`, `technical_issue`, `general_feedback`

**Priority:** `low`, `medium` (stored as `normal`), `high`

### Admin

| Method | Path | Description |
|--------|------|-------------|
| GET | `/admin/support/tickets?role=&status=&priority=&search=` | List + stats |
| GET | `/admin/support/tickets/{ticket_id}` | Detail + internal note |
| PATCH | `/admin/support/tickets/{ticket_id}/status` | `open`, `in_review`, `resolved`, `escalated`, `closed` |
| POST | `/admin/support/tickets/{ticket_id}/reply` | Public `admin_reply` + optional `internal_note` |

Public reply is stored on `support_tickets.resolution_note` (shown to business as `admin_reply`). Internal notes use `support_messages.is_internal_note=true`.

**Messages/chat deferred** — Firebase vs Supabase TBD; not part of support tickets.

---

## 19. Admin business verification

Business accounts become `active` on profile complete but **`is_verified` stays false** until admin approval (required for marketplace/contracts).

| Method | Path | Description |
|--------|------|-------------|
| GET | `/admin/businesses/verifications?status=pending\|verified` | List businesses awaiting or completed verification |
| GET | `/admin/businesses/verifications/{user_id}` | Detail |
| POST | `/admin/businesses/verifications/{user_id}/verify` | Set `is_verified=true` |
| POST | `/admin/businesses/verifications/{user_id}/revoke` | Set `is_verified=false` |

---

## 20. Testing flow (Postman)

1. **Signup/login** as business (`role=business`)
2. **GET profile** → `pending_profile`
3. **PUT profile** → fill required fields
4. **POST logo** (optional)
5. **POST complete** → `active` (still `is_verified=false`)
6. **Admin → Verify Business** (or manual DB update) before contracts/marketplace
7. **POST campaign** → draft
8. **POST publish** → active
9. **POST generate recommendations** (requires at least one active verified influencer in DB)
10. **POST invite** to influencer user id
11. Login as **influencer** → accept invite (Phase 4)
12. **POST contract** (Phase 5)
13. **GET dashboard/summary** — verify KPIs match your data
14. **POST support ticket** → **Admin → Reply** → business sees `admin_reply`

---

## 21. Explicitly out of scope (V1)

| Item | Reason |
|------|--------|
| Pricing packages on creator page | UI mock only |
| Platform subscription tiers | Not in product model |
| Admin approval for business signup | V1 auto-activate on complete; **admin verifies marketplace** via `/admin/businesses/verifications` |
| Campaign admin moderation | Optional later |
| Firebase chat message bodies | Metadata only in Postgres V1 — **messages API deferred** (Firebase vs Supabase TBD) |
| Chapa / payment webhooks | Payments phase |
| Business verification documents | `verification_documents` later |
| Email notifications for invites | `notifications` table insert optional |
| Frontend wiring | Deferred |

---

## 22. Related documents

| Doc | Purpose |
|-----|---------|
| `AUTH-API.md` | Signup, login, JWT, business `pending_profile` |
| `INFLUENCER-PROFILE-API.md` | Creator profile, portfolio, admin approval |
| `MEDIA-STORAGE.md` | MinIO, logo upload, S3 env vars |
| `DATABASE-QUICKREF.md` | Table list, Makefile |

---

*Document version: Adfluence V1 — business backend specification. Implement Phase 1 first, then campaigns, then invites, then contracts.*
