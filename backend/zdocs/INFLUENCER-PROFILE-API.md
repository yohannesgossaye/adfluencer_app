# Adfluence Influencer Profile API Specification

Backend design for influencer onboarding, profile CRUD, media uploads, profile submission, and admin approval. Based on:

- **Frontend UI:** `/influencer/profile`, `/influencer/dashboard`, `/admin/influencers`
- **PostgreSQL schema:** `auth.sql`, `admin.sql`, new `influencer_portfolio.sql`
- **Auth:** completed — see `AUTH-API.md`
- **Media:** MinIO/S3 — see `MEDIA-STORAGE.md`

**Out of scope:** campaigns, invites, contracts, dashboard KPIs, frontend integration (separate phase).

**Platform revenue (do not confuse with profile pricing):** Adfluence earns **5% from the business** and **5% from the influencer** on each **contract** only. There are **no subscriptions** and **no platform pricing tiers**. The frontend “Basic / Standard / Premium” package cards on `/influencer/profile` are **UI mock only** — **do not implement** as backend entities.

**Schema rule:** Use **existing tables only**. No new tables without explicit approval. Media metadata for portfolio samples in Phase B may use an **optional JSONB column** on `influencer_profiles` (ALTER, not a new table) — discuss before adding.

**Implementation strategy:** Backend first, test with Postman. Frontend mock removal + API wiring happens later.

---

## 1. Overview

### Purpose

After influencer signup (`users.status = pending_profile`), the creator completes their profile, uploads samples, and **submits for admin review**. Admin approves or rejects. Approved creators become marketplace-visible (`users.status = active`).

### Lifecycle diagram

```
SIGNUP (auth)
  users.status = pending_profile
  influencer_profiles row created (creator_name from signup only)

SAVE PROFILE (draft, repeatable)
  PUT /influencer/profile
  users.status stays pending_profile (or rejected if resubmitting)

UPLOAD MEDIA (optional but required for submit)
  POST /influencer/profile/avatar
  POST /influencer/profile/portfolio (max 3)

SUBMIT FOR APPROVAL
  POST /influencer/profile/submit
  users.status → pending_approval

ADMIN REVIEW
  GET  /admin/influencers/applications
  POST /admin/influencers/applications/{user_id}/approve
  POST /admin/influencers/applications/{user_id}/reject

APPROVE → users.status = active, influencer_profiles.approved_at, is_verified
REJECT  → users.status = rejected, users.status_reason set
          Creator can edit profile and submit again
```

### Frontend routes (reference — integration later)

| Route | Current behavior | Backend replaces |
|-------|------------------|------------------|
| `/influencer/profile` | localStorage `adfluence_influencer_profile` | GET/PUT profile + uploads |
| `/influencer/dashboard` | Mock 86% completion | GET profile `completeness_percent` (dashboard API later) |
| `/admin/influencers` | localStorage applications | Admin application APIs |

---

## 2. User status rules (influencer)

Maps to `user_status` enum. Enforced on every influencer profile endpoint.

| Status | Can view profile? | Can edit (PUT)? | Can upload media? | Can submit? |
|--------|-------------------|-----------------|-------------------|-------------|
| `pending_profile` | yes | yes | yes | yes (if complete) |
| `pending_approval` | yes | yes (updates draft; does **not** auto-resubmit) | yes | no (already submitted) |
| `active` | yes | yes | yes | no |
| `rejected` | yes | yes | yes | yes (resubmit after edits) |
| `suspended` | yes (read) | no | no | no |
| `blocked` | no | no | no | no |
| `removed` | no | no | no | no |

**Note:** While `pending_approval`, allow profile edits so creators can fix typos; changing status back to `pending_profile` on edit is **optional V1** — keep `pending_approval` until admin acts.

**Submit guard:** Only `pending_profile` or `rejected` may call submit. Return **409** `invalid_status_for_submit` otherwise.

**Role guard:** All `/influencer/*` routes require JWT `role = influencer`. Return **403** `forbidden` for business/admin tokens.

---

## 3. Database

### 3.1 Existing tables (already migrated)

#### `influencer_profiles` (`auth.sql`)

| Column | Type | Profile API usage |
|--------|------|-------------------|
| `id` | UUID PK | Profile id |
| `user_id` | UUID UNIQUE | Owner |
| `full_name` | VARCHAR(255) | Legal/display full name |
| `creator_name` | VARCHAR(255) | Public creator name (UI: "Creator name") |
| `username` | VARCHAR(100) UNIQUE | Optional slug; generate in Phase C |
| `profile_image_url` | TEXT | Avatar URL after MinIO upload |
| `bio` | TEXT | Bio |
| `city` | VARCHAR(100) | Location (UI: "Addis Ababa") |
| `region` | VARCHAR(100) | Optional; split from location later |
| `languages` | TEXT[] | e.g. `{Amharic,English}` |
| `social_links` | JSONB | `{tiktok, instagram, youtube}` URLs |
| `instagram_followers` | INTEGER | Per-platform or 0 |
| `tiktok_followers` | INTEGER | Per-platform or 0 |
| `youtube_subscribers` | INTEGER | Per-platform or 0 |
| `total_followers` | INTEGER | Sum or UI single "followers" field |
| `engagement_rate` | NUMERIC(5,2) | e.g. 7.20 |
| `audience_description` | TEXT | Optional extended notes |
| `content_style` | TEXT | Optional |
| `portfolio_description` | TEXT | Optional |
| `min_price_etb` | NUMERIC(12,2) | Creator collaboration rate **min** (ETB) — not a platform fee |
| `max_price_etb` | NUMERIC(12,2) | Creator collaboration rate **max** (ETB) |
| `approval_note` | TEXT | Admin internal note (sync with admin UI) |
| `is_verified` | BOOLEAN | Marketplace badge |
| `approved_at` | TIMESTAMPTZ | Set on admin approve |
| `approved_by` | UUID | Admin user id |

**At signup** only `full_name`, `creator_name`, `user_id` are populated (from signup `name`).

#### `influencer_profile_categories` (`admin.sql`)

Many-to-many: influencer ↔ `categories` (niches).

#### `influencer_profile_platforms` (`admin.sql`)

Many-to-many: influencer ↔ `social_platform` enum (`tiktok`, `instagram`, `youtube`, …).

#### `categories` (`admin.sql`)

Admin-managed niches. **Seed** for MVP (match frontend):

| name | slug |
|------|------|
| Beauty | beauty |
| Fashion | fashion |
| Food | food |
| Tech | tech |
| Fitness | fitness |
| Lifestyle | lifestyle |
| Travel | travel |

Migration: `migration/seed_categories.sql` (new, idempotent).

#### `users` (status fields used here)

| Column | Usage |
|--------|-------|
| `status` | Lifecycle |
| `status_reason` | Rejection message shown to creator |
| `status_updated_at` | Optional; set on approve/reject |

### 3.2 Media / portfolio (Phase B — no new table)

**Phase A:** JSON profile only. No upload endpoints.

**Phase B:** Avatar → existing `profile_image_url`. Content samples (up to 3) → MinIO files; store URLs in **`influencer_profiles.social_links` sub-key** or a **JSONB column added only with approval** (e.g. `portfolio_samples`). Do **not** add `influencer_portfolio_items` or other new tables unless explicitly approved.

See `MEDIA-STORAGE.md` for MinIO upload flow.

| Action | Tables |
|--------|--------|
| GET profile | `influencer_profiles`, `influencer_profile_categories`, `influencer_profile_platforms`, `categories`, `users` |
| PUT profile | `influencer_profiles`, category/platform join tables, `audit_logs` |
| Upload avatar (Phase B) | `influencer_profiles`, MinIO |
| Submit | `users` (status → pending_approval), `audit_logs` |
| Admin approve | `users`, `influencer_profiles`, `audit_logs` |
| Admin reject | `users`, `influencer_profiles`, `audit_logs` |

---

## 4. API conventions

### Base paths

| Portal | Base |
|--------|------|
| Influencer | `/api/v1/influencer` |
| Admin | `/api/v1/admin` |

### Auth

All endpoints require `Authorization: Bearer {access_token}` unless noted.

### Response envelope (same as auth)

```json
{
  "status": 200,
  "message": "profile retrieved successfully",
  "data": { }
}
```

Errors:

```json
{
  "status": 400,
  "message": "validation failed",
  "error": {
    "code": "invalid_bio",
    "message": "bio is required"
  }
}
```

### Timestamps

ISO 8601 UTC strings in JSON.

---

## 5. Influencer endpoints

### 5.1 Get my profile

**`GET /api/v1/influencer/profile`**

Frontend: `/influencer/profile` (load on mount).

**Auth:** influencer role.

**Response `data` — `InfluencerProfileResponse`:**

```json
{
  "user": {
    "id": "uuid",
    "email": "yohannesgossaye@gmail.com",
    "phone_number": "+251953422606",
    "status": "pending_profile",
    "status_reason": null,
    "email_verified": false,
    "phone_verified": false
  },
  "profile": {
    "id": "uuid",
    "creator_name": "Biruk Alemu",
    "full_name": "Biruk Alemu",
    "username": null,
    "bio": "Tech creator focused on...",
    "profile_image_url": null,
    "city": "Addis Ababa",
    "region": null,
    "languages": ["Amharic", "English"],
    "niches": [
      { "id": "uuid", "name": "Tech", "slug": "tech" }
    ],
    "platforms": ["youtube"],
    "social_links": {
      "tiktok": "https://tiktok.com/@birukcreates",
      "instagram": "https://instagram.com/birukcreates",
      "youtube": "https://youtube.com/@birukcreates"
    },
    "instagram_followers": 0,
    "tiktok_followers": 0,
    "youtube_subscribers": 84000,
    "total_followers": 84000,
    "engagement_rate": 7.2,
    "min_price_etb": 15000,
    "max_price_etb": 35000,
    "audience_description": null,
    "is_verified": false,
    "approved_at": null
  },
  "portfolio_items": [
    {
      "id": "uuid",
      "file_url": "http://localhost:9000/adfluence-media/media/...",
      "media_type": "video",
      "mime_type": "video/mp4",
      "sort_order": 0,
      "created_at": "2026-06-26T12:00:00Z"
    }
  ],
  "completeness": {
    "percent": 72,
    "missing": ["profile_image", "portfolio_sample", "min_price_etb"],
    "can_submit": false
  }
}
```

**Errors:**

| HTTP | code | When |
|------|------|------|
| 401 | unauthorized | Missing/invalid token |
| 403 | forbidden | Not influencer role |
| 404 | profile_not_found | No profile row (should not happen post-signup) |

---

### 5.2 Update profile (draft save)

**`PUT /api/v1/influencer/profile`**

Frontend: "Save Profile" button.

**Content-Type:** `application/json`

**Request — `UpdateInfluencerProfileRequest`:**

| Field | Type | Required | Validation | DB mapping |
|-------|------|----------|------------|------------|
| `creator_name` | string | yes | 2–255 chars | `creator_name`, also update `full_name` if not split |
| `bio` | string | yes | 20–2000 chars | `bio` |
| `niches` | string[] | yes | 1–5 items; each must match active category **name** or **slug** | `influencer_profile_categories` |
| `platforms` | string[] | yes | 1–5; values: `tiktok`, `instagram`, `youtube` (lowercase) | `influencer_profile_platforms` |
| `social_links` | object | no | Valid URLs if present | `social_links` JSONB |
| `social_links.tiktok` | string | no | URL or empty | |
| `social_links.instagram` | string | no | URL or empty | |
| `social_links.youtube` | string | no | URL or empty | |
| `followers` | integer | no | ≥ 0 | See follower mapping §5.2.1 |
| `instagram_followers` | integer | no | ≥ 0 | Direct column |
| `tiktok_followers` | integer | no | ≥ 0 | Direct column |
| `youtube_subscribers` | integer | no | ≥ 0 | Direct column |
| `engagement_rate` | number | no | 0–100 | `engagement_rate` |
| `min_price_etb` | number | no | ≥ 0 | `min_price_etb` |
| `max_price_etb` | number | no | ≥ min | `max_price_etb` |
| `languages` | string[] | no | 1–10 items | `languages` TEXT[] |
| `city` | string | no | 2–100 chars | `city` (UI "Location") |
| `region` | string | no | 2–100 chars | `region` |
| `audience_description` | string | no | max 2000 | `audience_description` |

**Example (matches frontend mock):**

```json
{
  "creator_name": "Biruk Alemu",
  "bio": "Tech creator focused on practical reviews, app walkthroughs, and honest recommendations for Ethiopian audiences.",
  "niches": ["Tech"],
  "platforms": ["youtube"],
  "social_links": {
    "tiktok": "https://tiktok.com/@birukcreates",
    "instagram": "https://instagram.com/birukcreates",
    "youtube": "https://youtube.com/@birukcreates"
  },
  "followers": 84000,
  "engagement_rate": 7.2,
  "min_price_etb": 15000,
  "max_price_etb": 35000,
  "languages": ["Amharic", "English"],
  "city": "Addis Ababa"
}
```

#### 5.2.1 Follower field mapping

Frontend sends single `followers`. Backend logic:

1. If `instagram_followers`, `tiktok_followers`, `youtube_subscribers` provided explicitly → use them.
2. Else if `followers` provided → distribute to selected platforms:
   - If **one** platform in `platforms` → set that platform's column = `followers`.
   - If **multiple** platforms → set `total_followers = followers`, leave per-platform 0 unless URLs imply primary platform (optional: set first platform only).
3. Always recompute `total_followers` = sum of three columns when any per-platform value set.

#### 5.2.2 Transaction steps

1. Validate request + status allows edit
2. Update `influencer_profiles` row
3. Replace `influencer_profile_categories` (delete all for user, re-insert)
4. Replace `influencer_profile_platforms` (delete all, re-insert)
5. Write `audit_logs` (`profile_update`)
6. Return full GET profile response

**Does not change `users.status`.**

**Errors:**

| HTTP | code |
|------|------|
| 400 | validation errors (`invalid_bio`, `invalid_niche`, …) |
| 403 | forbidden / status blocks edit |
| 409 | `profile_locked` (blocked/suspended) |

---

### 5.3 Upload avatar

**`POST /api/v1/influencer/profile/avatar`**

**Content-Type:** `multipart/form-data`

| Field | Type | Required |
|-------|------|----------|
| `file` | file | yes |

See `MEDIA-STORAGE.md` §7.1 for validation and MinIO flow.

**Response `data`:**

```json
{
  "profile_image_url": "http://localhost:9000/adfluence-media/media/influencer/.../avatar/....jpg"
}
```

---

### 5.4 Upload portfolio sample

**`POST /api/v1/influencer/profile/portfolio`**

**Content-Type:** `multipart/form-data`

| Field | Type | Required |
|-------|------|----------|
| `file` | file | yes |
| `sort_order` | integer | no (default next 0,1,2) |

Max **3** items. See `MEDIA-STORAGE.md` §7.2.

**Response `data`:**

```json
{
  "item": {
    "id": "uuid",
    "file_url": "...",
    "media_type": "image",
    "mime_type": "image/jpeg",
    "sort_order": 0,
    "created_at": "..."
  },
  "portfolio_count": 1
}
```

---

### 5.5 Delete portfolio sample

**`DELETE /api/v1/influencer/profile/portfolio/{item_id}`**

Deletes MinIO object + DB row. Returns updated `portfolio_items` array.

---

### 5.6 Submit profile for admin approval

**`POST /api/v1/influencer/profile/submit`**

No body.

**Preconditions (all required — else 400 `profile_incomplete`):**

| Check | Rule |
|-------|------|
| Status | `pending_profile` or `rejected` |
| `creator_name` | non-empty |
| `bio` | ≥ 20 chars |
| `niches` | ≥ 1 category linked |
| `platforms` | ≥ 1 platform linked |
| `city` | non-empty |
| `min_price_etb` & `max_price_etb` | both set, max ≥ min |
| `engagement_rate` | set |
| `total_followers` | > 0 |
| At least one social link | matching a selected platform |
| `profile_image_url` | set |
| Portfolio items | ≥ 1 (recommend 2+; MVP require ≥ 1) |

**Optional recommended (warn in completeness, don't block V1):**

- Email verified
- Phone verified
- 3 portfolio samples

**On success:**

1. `users.status` → `pending_approval`
2. Clear `users.status_reason` if resubmitting from `rejected`
3. `audit_logs` → `status_change` or custom `profile_submitted`
4. Response includes updated user status

```json
{
  "status": "pending_approval",
  "message": "Profile submitted for review",
  "submitted_at": "2026-06-26T12:00:00Z"
}
```

**Errors:**

| HTTP | code |
|------|------|
| 400 | `profile_incomplete` + `missing` array |
| 409 | `invalid_status_for_submit` |
| 409 | `already_pending_approval` |

---

### 5.7 Profile completeness score

Computed on GET profile (not stored). Matches dashboard "Profile 86%" UI.

| Field / check | Weight |
|---------------|--------|
| `creator_name` | 10 |
| `bio` (≥ 20 chars) | 15 |
| ≥ 1 niche | 10 |
| ≥ 1 platform | 10 |
| ≥ 1 social link for selected platform | 10 |
| `total_followers` > 0 | 10 |
| `engagement_rate` set | 5 |
| `min_price_etb` & `max_price_etb` | 10 |
| `city` | 5 |
| `languages` ≥ 1 | 5 |
| `profile_image_url` | 10 |
| Each portfolio item (max 3) | 10 each (30 max) |

**Total:** 100 points max.

`can_submit` = all submit preconditions met (§5.6).

---

## 6. Admin endpoints

Base: `/api/v1/admin/influencers`

**Auth:** JWT `role = admin`. Use login with `"portal": "admin"`.

### 6.1 List applications

**`GET /api/v1/admin/influencers/applications`**

Frontend: `/admin/influencers`

**Query params:**

| Param | Type | Description |
|-------|------|-------------|
| `status` | string | `pending`, `approved`, `rejected`, or omit for all |
| `search` | string | Match creator_name, niche, platform |
| `page` | int | Default 1 |
| `limit` | int | Default 20, max 50 |

**Data source:** Users with `role = influencer` and relevant statuses. Primary queue = `pending_approval`.

**Response `data`:**

```json
{
  "items": [
    {
      "user_id": "uuid",
      "creator_name": "Selam Tadesse",
      "profile_image_url": "...",
      "primary_niche": "Beauty",
      "primary_platform": "tiktok",
      "total_followers": 62000,
      "engagement_rate": 6.4,
      "completeness_percent": 72,
      "pricing_summary": "ETB 12000 – 22000",
      "status": "pending_approval",
      "submitted_at": "2026-05-10T00:00:00Z",
      "bio_preview": "Beauty creator focused on..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5
  }
}
```

Map frontend mock fields:

| Mock field | API field |
|------------|-----------|
| `name` | `creator_name` |
| `imageUrl` | `profile_image_url` |
| `niche` | `primary_niche` |
| `platform` | `primary_platform` |
| `completeness` | `completeness_percent` + human `missing` optional |
| `adminNotes` | `influencer_profiles.approval_note` |

---

### 6.2 Get application detail

**`GET /api/v1/admin/influencers/applications/{user_id}`**

Full profile + portfolio + user status for admin side panel.

Includes: bio, all social links, audience_description, pricing, all portfolio URLs, admin notes, audit snippet.

---

### 6.3 Approve application

**`POST /api/v1/admin/influencers/applications/{user_id}/approve`**

**Request (optional):**

```json
{
  "admin_note": "Strong portfolio, approved for marketplace."
}
```

**Steps:**

1. Verify target user is influencer with status `pending_approval`
2. `users.status` → `active`
3. `influencer_profiles.is_verified` → `true`
4. `influencer_profiles.approved_at` → now
5. `influencer_profiles.approved_by` → admin user id
6. Store `approval_note` if provided
7. `audit_logs` → `approve`

**Response:**

```json
{
  "user_id": "uuid",
  "status": "active",
  "approved_at": "2026-06-26T12:00:00Z"
}
```

---

### 6.4 Reject application

**`POST /api/v1/admin/influencers/applications/{user_id}/reject`**

**Request:**

```json
{
  "reason": "Portfolio samples do not meet quality guidelines. Please upload clearer content.",
  "admin_note": "Internal: resubmit after guideline update."
}
```

| Field | Required |
|-------|----------|
| `reason` | yes (shown to creator → `users.status_reason`) |
| `admin_note` | no → `approval_note` |

**Steps:**

1. Verify status `pending_approval`
2. `users.status` → `rejected`
3. `users.status_reason` → `reason`
4. `influencer_profiles.is_verified` → `false`
5. Clear `approved_at`, `approved_by`
6. `audit_logs` → `reject`

Creator can edit profile and call submit again.

---

### 6.5 Update admin notes (without status change)

**`PATCH /api/v1/admin/influencers/applications/{user_id}/notes`**

```json
{
  "admin_note": "Waiting for ID verification."
}
```

Updates `influencer_profiles.approval_note` only.

---

## 7. Public / business read (future — not V1)

When business views creator on `/business/influencers/[id]`:

**`GET /api/v1/influencers/{user_id}/public`**

- Only if `users.status = active` and `is_verified = true`
- Returns public-safe fields (no email, phone, admin notes)

Document now; implement with campaigns/recommendations phase.

---

## 8. Backend code structure (follow auth pattern)

```
internal/constant/model/dto/influencer/   — DTOs + validation
internal/storage/persist/influencer/      — SQL
internal/module/influencer/               — business logic
internal/handler/rest/http/influencer/    — HTTP messages
internal/glue/routing/influencer/         — routes
internal/module/admin/                    — approval logic (or admin/influencer subpackage)
internal/glue/routing/admin/              — admin routes
platform/storage/                         — MinIO (MEDIA-STORAGE.md)
migration/influencer_portfolio.sql
migration/seed_categories.sql
```

### Route registration

Wire in `initiator/route.go` (currently only auth registered):

```go
router.Route("/api/v1", func(r chi.Router) {
    auth.AuthRoutes(r, rest.auth, authMW)
    influencer.InfluencerRoutes(r, rest.influencer, authMW)
    admin.AdminRoutes(r, rest.admin, authMW)
})
```

### Replace stub route

Current stub `POST /influencer/get_profile` → replace with REST routes above.

---

## 9. Error codes (full list)

| HTTP | code | Description |
|------|------|-------------|
| 400 | `invalid_bio` | Bio validation |
| 400 | `invalid_niche` | Unknown or inactive category |
| 400 | `invalid_platform` | Not in social_platform enum |
| 400 | `invalid_social_link` | Malformed URL |
| 400 | `invalid_price_range` | max < min |
| 400 | `profile_incomplete` | Submit preconditions fail |
| 400 | `invalid_file` | Upload validation |
| 401 | unauthorized | |
| 403 | forbidden | Wrong role |
| 403 | profile_locked | Suspended/blocked |
| 404 | profile_not_found | |
| 404 | portfolio_item_not_found | |
| 409 | invalid_status_for_submit | |
| 409 | portfolio_limit_reached | 3 samples max |
| 503 | storage_not_configured | MinIO not set |
| 500 | internal | |

---

## 10. Audit log events

| Event | action | entity_type | entity_id |
|-------|--------|-------------|-----------|
| Profile saved | `profile_update` | `influencer_profile` | profile id |
| Submitted | `status_change` | `user` | user id |
| Admin approve | `approve` | `user` | target user id |
| Admin reject | `reject` | `user` | target user id |

Include `actor_user_id`, `ip_address`, `user_agent` from request meta (same as auth).

---

## 11. Implementation phases

### Phase A — JSON profile (no MinIO)

- [ ] Migration `seed_categories.sql`
- [ ] GET/PUT `/influencer/profile`
- [ ] Completeness calculation
- [ ] POST `/influencer/profile/submit` (without media requirements relaxed: no avatar/portfolio required temporarily OR mock URLs)
- [ ] Wire influencer routes in initiator
- [ ] Postman folder **Influencer → Profile**

**Testing:** Postman only. Submit may use relaxed rules until Phase B.

### Phase B — MinIO + media

- [ ] MinIO docker-compose
- [ ] `platform/storage` (S3Storage from reference code)
- [ ] Migration `influencer_portfolio.sql`
- [ ] Avatar + portfolio upload/delete
- [ ] Enforce full submit preconditions (avatar + ≥1 portfolio)
- [ ] Postman multipart requests

### Phase C — Admin approval

- [ ] GET list/detail applications
- [ ] Approve / reject / notes
- [ ] Admin routes + middleware (admin role)
- [ ] Postman folder **Admin → Influencer Applications**

### Phase D — Polish (optional before frontend)

- [ ] Username generation/slug
- [ ] Pricing packages JSONB
- [ ] Public profile read endpoint
- [ ] Resubmit flow edge cases

### Phase E — Frontend integration (later)

- [ ] Replace localStorage in `profile.tsx`
- [ ] Wire `/admin/influencers`
- [ ] Status-based redirects (pending_approval screen)

---

## 12. Postman collection (planned additions)

Under **Influencer → Profile**:

| Request | Method | Path |
|---------|--------|------|
| Get Profile | GET | `/api/v1/influencer/profile` |
| Update Profile | PUT | `/api/v1/influencer/profile` |
| Upload Avatar | POST | `/api/v1/influencer/profile/avatar` |
| Upload Portfolio | POST | `/api/v1/influencer/profile/portfolio` |
| Delete Portfolio Item | DELETE | `/api/v1/influencer/profile/portfolio/{id}` |
| Submit for Approval | POST | `/api/v1/influencer/profile/submit` |

Under **Admin → Influencers**:

| Request | Method | Path |
|---------|--------|------|
| List Applications | GET | `/api/v1/admin/influencers/applications` |
| Get Application | GET | `/api/v1/admin/influencers/applications/{user_id}` |
| Approve | POST | `/api/v1/admin/influencers/applications/{user_id}/approve` |
| Reject | POST | `/api/v1/admin/influencers/applications/{user_id}/reject` |
| Save Admin Notes | PATCH | `/api/v1/admin/influencers/applications/{user_id}/notes` |

---

## 13. Test flow (end-to-end, Postman)

1. **Signup — Influencer** (or use existing user)
2. **Login** → save token
3. **GET profile** → verify defaults from signup
4. **PUT profile** → fill all JSON fields
5. **Upload avatar** (Phase B)
6. **Upload 1–3 portfolio files** (Phase B)
7. **Submit profile** → status `pending_approval`
8. **Admin login** (`portal: admin`)
9. **List applications** → find user
10. **Approve** or **Reject**
11. **Login as influencer** → GET `/auth/me` → verify new status

---

## 14. Environment variables (summary)

From auth + new storage:

```env
# existing auth, db, jwt, smtp, sms ...

# MinIO (Phase B) — see MEDIA-STORAGE.md
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY_ID=minioadmin
S3_SECRET_ACCESS_KEY=minioadmin123
S3_BUCKET_NAME=adfluence-media
S3_PUBLIC_URL=http://localhost:9000/adfluence-media
```

---

## 15. Frontend field mapping (complete)

| UI label (`profile.tsx`) | JSON field | DB |
|--------------------------|------------|-----|
| Creator name | `creator_name` | `influencer_profiles.creator_name` |
| Bio | `bio` | `bio` |
| Niche (multi) | `niches[]` | `influencer_profile_categories` |
| Platforms (multi) | `platforms[]` | `influencer_profile_platforms` |
| TikTok link | `social_links.tiktok` | JSONB |
| Instagram link | `social_links.instagram` | JSONB |
| YouTube link | `social_links.youtube` | JSONB |
| Followers | `followers` | follower columns §5.2.1 |
| Engagement rate (%) | `engagement_rate` | `engagement_rate` |
| Price min (ETB) | `min_price_etb` | `min_price_etb` |
| Price max (ETB) | `max_price_etb` | `max_price_etb` |
| Languages | `languages[]` | `languages` |
| Location | `city` | `city` |
| Avatar image | multipart upload | `profile_image_url` |
| Content samples (×3) | multipart upload | `influencer_portfolio_items` |
| Pricing packages | Phase C | JSONB optional |
| Profile % complete | `completeness.percent` | computed |
| Verified creator badge | `is_verified` | after admin approve |

---

## 16. Related documents

| Doc | Purpose |
|-----|---------|
| `AUTH-API.md` | Signup, login, tokens, verify email/phone |
| `MEDIA-STORAGE.md` | MinIO, S3 client, multipart, object keys |
| `DATABASE-QUICKREF.md` | Table overview |

---

*Document version: Adfluence V1 — influencer profile + admin approval; backend-first; frontend integration deferred.*
