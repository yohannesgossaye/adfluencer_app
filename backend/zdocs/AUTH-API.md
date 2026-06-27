# Adfluence Auth API Specification

Backend auth design for Golang implementation. Based on the **frontend UI** (`/auth/signup`, `/auth/login`, `/admin/login`) and the **PostgreSQL schema** (`auth.sql`, `admin.sql`).

**Out of scope for this document:** pricing/subscriptions (replaced by 5% business fee + 5% influencer fee on contracts). Campaign, chat, contract APIs are separate.

---

## 1. Overview

### Roles

| Role | Public signup? | Default `users.status` after signup | Portal redirect after login |
|------|----------------|-------------------------------------|-----------------------------|
| `business` | Yes (`/auth/signup?role=business`) | `pending_profile` | `/business/dashboard` |
| `influencer` | Yes (`/auth/signup?role=influencer`) | `pending_profile` | `/influencer/dashboard` |
| `admin` | **No** — seed only | `active` (seeded) | `/admin/dashboard` |

### Auth model (recommended)

- **Access token:** short-lived JWT (Bearer), stored in frontend `localStorage` key `access_token`
- **Refresh token:** opaque random string, hashed in `user_sessions.refresh_token_hash`
- **Password:** bcrypt/argon2 hash in `users.password_hash` — never return to client

### Frontend today (mock)

| Route | Current behavior | Backend must replace |
|-------|------------------|-------------------|
| `/auth/signup` | Demo only — no API, redirects to login | `POST /api/v1/auth/signup` |
| `/auth/login` | Matches `mock-users.ts` | `POST /api/v1/auth/login` |
| `/admin/login` | Admin mock only | Same login endpoint OR `POST /api/v1/auth/admin/login` |
| `/auth/forgot-password` | Placeholder | `POST /api/v1/auth/forgot-password` |
| `/auth/reset-password` | Placeholder | `POST /api/v1/auth/reset-password` |
| `/auth/verify` | Placeholder | `POST /api/v1/auth/verify-email` |

---

## 2. User status → login behavior

Maps to `user_status` enum in schema. Backend **must** enforce this on login and on every protected request.

| Status | Login allowed? | Frontend redirect / screen |
|--------|----------------|----------------------------|
| `pending_profile` | Yes | Portal dashboard → prompt profile completion (`/business/...` or `/influencer/profile`) |
| `pending_approval` | Yes | Waiting-approval page (influencer after profile submit) |
| `active` | Yes | Normal dashboard |
| `rejected` | Yes | Rejection reason page (`users.status_reason`) |
| `suspended` | No (if `now() < suspended_until`) | Error: account suspended until date |
| `blocked` | No | Error: account blocked — contact support |
| `removed` | No | Error: account removed |

### Status transitions (auth-related only)

```
SIGNUP (business)
  → pending_profile
  → (profile completed) → active

SIGNUP (influencer)
  → pending_profile
  → (profile completed + submitted) → pending_approval
  → (admin approves) → active
  → (admin rejects) → rejected

ADMIN SEED
  → active
```

Profile completion APIs are separate documents, but auth responses must always return current `status` so the frontend can redirect correctly.

---

## 3. Database tables touched by auth

| Action | Tables written |
|--------|----------------|
| Business signup | `users`, `business_profiles`, `audit_logs` (signup) |
| Influencer signup | `users`, `influencer_profiles`, `audit_logs` (signup) |
| Admin seed | `users`, `admin_profiles` |
| Login | `users` (last_login_at), `user_sessions`, `audit_logs` (login, optional) |
| Logout | `user_sessions` (revoked_at) |
| Refresh token | `user_sessions` (new row or rotate) |
| Forgot password | `password_reset_tokens` |
| Reset password | `password_reset_tokens`, `users` (password_hash) |
| Verify email | `email_verification_tokens`, `users` (email_verified) |

**Not created at signup:** `verification_documents`, `bank_accounts` — user adds later.

---

## 4. API endpoints

Base path assumed: `/api/v1/auth`

All JSON. Timestamps ISO 8601 UTC. UUIDs as strings.

---

### 4.1 Signup — Business

**`POST /api/v1/auth/signup`**

Frontend: `/auth/signup` with role **Business**, fields from `signup-page.tsx`.

#### Request DTO — `SignupRequest`

| Field | Type | Required | Validation | Maps to DB |
|-------|------|----------|------------|------------|
| `role` | `string` | yes | `"business"` or `"influencer"` only | `users.role` |
| `name` | `string` | yes | 2–255 chars | `business_profiles.business_name` OR `influencer_profiles.creator_name` |
| `email` | `string` | yes | valid email, unique | `users.email` |
| `phone` | `string` | yes | Ethiopian format `+251...`, unique | `users.phone_number` |
| `password` | `string` | yes | min 8 chars, 1 letter + 1 number (recommend) | `users.password_hash` (hashed) |
| `confirm_password` | `string` | yes | must equal `password` | not stored |

**Business example:**

```json
{
  "role": "business",
  "name": "Acme Coffee",
  "email": "owner@acmecoffee.et",
  "phone": "+251911223344",
  "password": "SecurePass1",
  "confirm_password": "SecurePass1"
}
```

#### Backend writes (business)

**`users`**

| Column | Value |
|--------|-------|
| `email` | request.email (lowercase) |
| `phone_number` | normalized phone |
| `password_hash` | hashed password |
| `role` | `business` |
| `status` | `pending_profile` |
| `email_verified` | `false` |
| `phone_verified` | `false` |

**`business_profiles`**

| Column | Value |
|--------|-------|
| `user_id` | new user id |
| `owner_full_name` | same as `name` for now (or split later) |
| `business_name` | request.name |
| `is_verified` | `false` |

#### Response DTO — `SignupResponse` (201)

```json
{
  "user": {
    "id": "uuid",
    "email": "owner@acmecoffee.et",
    "phone_number": "+251911223344",
    "role": "business",
    "status": "pending_profile",
    "email_verified": false,
    "phone_verified": false,
    "created_at": "2026-06-26T10:00:00Z"
  },
  "profile": {
    "id": "uuid",
    "business_name": "Acme Coffee",
    "owner_full_name": "Acme Coffee",
    "is_verified": false
  },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "opaque-random-string",
    "expires_in": 900
  },
  "redirect_to": "/business/dashboard"
}
```

**Notes:**
- Option A: return tokens and auto-login (better UX — frontend can skip login page).
- Option B: no tokens — frontend redirects to `/auth/login` (matches current mock).
- Recommend **Option A** with `redirect_to` based on role + status.

#### Errors

| HTTP | Code | When |
|------|------|------|
| 400 | `VALIDATION_ERROR` | missing fields, password mismatch, weak password |
| 409 | `EMAIL_ALREADY_EXISTS` | email taken |
| 409 | `PHONE_ALREADY_EXISTS` | phone taken |
| 403 | `SIGNUP_DISABLED` | role admin attempted |

---

### 4.2 Signup — Influencer

Same endpoint `POST /api/v1/auth/signup` with `role: "influencer"`.

#### Request DTO

Same `SignupRequest` shape. Label on frontend is **Creator name** instead of Business name.

**Influencer example:**

```json
{
  "role": "influencer",
  "name": "Selam Beauty Studio",
  "email": "selam@creator.et",
  "phone": "+251922334455",
  "password": "SecurePass1",
  "confirm_password": "SecurePass1"
}
```

#### Backend writes (influencer)

**`users`** — same as business but `role = influencer`, `status = pending_profile`

**`influencer_profiles`**

| Column | Value |
|--------|-------|
| `user_id` | new user id |
| `full_name` | request.name |
| `creator_name` | request.name |
| `is_verified` | `false` |

Username can be generated later during profile completion (e.g. slug from creator_name).

#### Response DTO — `SignupResponse` (201)

```json
{
  "user": {
    "id": "uuid",
    "email": "selam@creator.et",
    "phone_number": "+251922334455",
    "role": "influencer",
    "status": "pending_profile",
    "email_verified": false,
    "phone_verified": false,
    "created_at": "2026-06-26T10:00:00Z"
  },
  "profile": {
    "id": "uuid",
    "full_name": "Selam Beauty Studio",
    "creator_name": "Selam Beauty Studio",
    "is_verified": false
  },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "opaque-random-string",
    "expires_in": 900
  },
  "redirect_to": "/influencer/dashboard"
}
```

---

### 4.3 Login

**`POST /api/v1/auth/login`**

Frontend: `/auth/login` — email + password. Also used for business, influencer, and admin (same form).

#### Request DTO — `LoginRequest`

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `email` | `string` | yes | case-insensitive match |
| `password` | `string` | yes | |
| `remember_me` | `boolean` | no | if true, longer refresh token TTL |
| `portal` | `string` | no | `"admin"` when logging in from `/admin/login` |
| `next` | `string` | no | safe internal path only, e.g. `/business/campaigns` — overrides default `redirect_to` |

```json
{
  "email": "yohannes@gmail.com",
  "password": "yoya1111",
  "remember_me": true,
  "next": "/business/campaigns"
}
```

#### Backend logic (order matters)

1. Check login lockout (progressive rate limit — see §14)
2. Find user by email
3. Verify password hash (same error if email not found)
4. Check `users.status`:
   - `blocked` / `removed` → 403
   - `suspended` + `suspended_until > now()` → 403
   - `suspended` + expired → optionally auto-set `active` (product decision)
5. Clear login lockout on success
6. Update `users.last_login_at`
7. Insert `user_sessions` row
8. Write `audit_logs` action `login`
9. Load profile snippet by role
10. Compute `redirect_to` from role + status, or use `next` if valid internal path

#### Response DTO — `LoginResponse` (200)

```json
{
  "user": {
    "id": "uuid",
    "email": "yohannes@gmail.com",
    "phone_number": "+251911000000",
    "role": "business",
    "status": "active",
    "email_verified": true,
    "phone_verified": false,
    "status_reason": null,
    "suspended_until": null,
    "last_login_at": "2026-06-26T12:00:00Z"
  },
  "profile": {
    "type": "business",
    "id": "uuid",
    "business_name": "Yohannes Trading",
    "owner_full_name": "Yohannes",
    "is_verified": true
  },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "opaque-random-string",
    "expires_in": 900
  },
  "redirect_to": "/business/dashboard"
}
```

**Profile union in response** — `profile.type` discriminates:

| `profile.type` | Fields |
|----------------|--------|
| `business` | `id`, `business_name`, `owner_full_name`, `is_verified` |
| `influencer` | `id`, `full_name`, `creator_name`, `username`, `is_verified` |
| `admin` | `id`, `full_name`, `permission_level` |

#### Redirect rules — `redirect_to`

| role | status | redirect_to |
|------|--------|-------------|
| business | `pending_profile` | `/business/dashboard` (show complete-profile banner) |
| business | `active` | `/business/dashboard` |
| business | `rejected` | `/auth/account-status` or `/business/account-rejected` |
| influencer | `pending_profile` | `/influencer/profile` |
| influencer | `pending_approval` | `/influencer/dashboard` (waiting approval UI) |
| influencer | `active` | `/influencer/dashboard` |
| influencer | `rejected` | `/auth/account-status` |
| admin | `active` | `/admin/dashboard` |

Frontend `api-client.ts` supports `?next=` on 401 — login should accept optional `next` query/body to override redirect after success.

#### Optional on login request

| Field | Type | Notes |
|-------|------|-------|
| `next` | `string` | safe internal path only, e.g. `/business/campaigns` |

#### Errors

| HTTP | Code | Message (frontend shows) |
|------|------|--------------------------|
| 401 | `INVALID_CREDENTIALS` | Incorrect email or password. |
| 403 | `ACCOUNT_BLOCKED` | This account has been blocked. Please contact Adfluence support. |
| 403 | `ACCOUNT_SUSPENDED` | Account suspended until {date}. |
| 403 | `ACCOUNT_REMOVED` | Account no longer available. |
| 403 | `NOT_ADMIN` | Admin access required (when `portal: admin`). |
| 429 | `LOGIN_LOCKED` | Too many failed login attempts — try again later. |

---

### 4.4 Admin login

**Option A (recommended):** same `POST /api/v1/auth/login` — admin uses `/auth/login` or `/admin/login`.

**Option B:** separate `POST /api/v1/auth/admin/login` — identical DTO/response but rejects non-admin roles with 403.

Frontend `/admin/login` only accepts admin credentials today.

#### Admin — no signup endpoint

Admin accounts are created only via:

- **Database seed** (migration/SQL or Go seed command)
- Never exposed on `/auth/signup`

#### Admin seed data (development)

Seed one super admin:

**`users`**

| Column | Example |
|--------|---------|
| email | `admin@adfluence.com` |
| password_hash | hash of `admin1111` (change in production) |
| role | `admin` |
| status | `active` |
| email_verified | `true` |

**`admin_profiles`**

| Column | Example |
|--------|---------|
| user_id | seeded user id |
| full_name | `Adfluence Admin` |
| permission_level | `super_admin` |

Login response same `LoginResponse` with `role: admin`, `redirect_to: /admin/dashboard`.

---

### 4.5 Refresh token

**`POST /api/v1/auth/refresh`**

#### Request DTO — `RefreshTokenRequest`

```json
{
  "refresh_token": "opaque-random-string"
}
```

#### Response DTO — `RefreshTokenResponse` (200)

```json
{
  "access_token": "eyJ...",
  "refresh_token": "new-opaque-string",
  "expires_in": 900
}
```

Rotate refresh token on each use. Set old `user_sessions.revoked_at`.

---

### 4.6 Logout

**`POST /api/v1/auth/logout`**

Requires `Authorization: Bearer {access_token}`.

#### Request DTO — `LogoutRequest`

```json
{
  "refresh_token": "opaque-random-string"
}
```

#### Response (200)

```json
{
  "message": "Logged out successfully"
}
```

Revoke matching `user_sessions` row. Frontend should clear `access_token` and mock session.

---

### 4.7 Get current user (session check)

**`GET /api/v1/auth/me`**

Requires Bearer token.

#### Response DTO — `MeResponse` (200)

Same shape as login `user` + `profile` without tokens.

```json
{
  "user": { "...": "..." },
  "profile": { "type": "business", "...": "..." }
}
```

Used on app load to restore session and enforce status-based UI.

---

### 4.8 Forgot password

**`POST /api/v1/auth/forgot-password`**

Frontend: `/auth/forgot-password` (placeholder today).

#### Request DTO — `ForgotPasswordRequest`

```json
{
  "email": "owner@acmecoffee.et"
}
```

#### Response (200) — always generic (prevent email enumeration)

```json
{
  "message": "If an account exists for this email, a reset link has been sent."
}
```

#### Backend writes

- Rate limit: max **3 requests per hour** per email → 429 `too_many_forgot_password_requests`
- Insert `password_reset_tokens` (hash token, `expires_at` e.g. +1 hour)
- Send email with link: `/auth/reset-password?token=...`

---

### 4.9 Reset password

**`POST /api/v1/auth/reset-password`**

#### Request DTO — `ResetPasswordRequest`

| Field | Type | Required |
|-------|------|----------|
| `token` | string | yes |
| `password` | string | yes |
| `confirm_password` | string | yes |

```json
{
  "token": "raw-token-from-email",
  "password": "NewSecurePass1",
  "confirm_password": "NewSecurePass1"
}
```

#### Response (200)

```json
{
  "message": "Password updated successfully",
  "redirect_to": "/auth/login"
}
```

#### Backend writes

- Validate token hash in `password_reset_tokens`, not expired, `used_at` null
- Update `users.password_hash`
- Set `password_reset_tokens.used_at`
- Revoke all `user_sessions` for that user (force re-login)

---

### 4.10 Verify email (OTP)

Two-step flow — user must be logged in (Bearer token).

#### 4.10.1 Send OTP

**`POST /api/v1/auth/send-verify-email`**

Requires `Authorization: Bearer {access_token}`. No request body.

#### Response (200)

Handler message: `"verification code sent"`.

Email body: `Your Adfluence email verification code is: 482913` (6 digits, **5 minute** expiry).

Rate limit: **3 sends per hour** per user → 429 `too_many_verify_email_requests`.

If already verified → 409 `email_already_verified`.

---

#### 4.10.2 Verify OTP

**`POST /api/v1/auth/verify-email`**

Frontend: `/auth/verify`.

#### Request DTO — `VerifyEmailRequest`

```json
{
  "otp": "482913"
}
```

#### Response (200)

Handler message: `"email verified successfully"`.

```json
{
  "user": {
    "id": "uuid",
    "email_verified": true
  }
}
```

#### Errors

| HTTP | Code | When |
|------|------|------|
| 400 | `invalid_verify_otp` | Wrong or expired OTP |
| 409 | `email_already_verified` | Already verified |

Login still allowed with `email_verified: false` for MVP.

---

## 5. JWT access token claims (recommended)

Payload inside `access_token` (do not trust alone — always verify signature):

| Claim | Value |
|-------|-------|
| `sub` | `users.id` |
| `email` | `users.email` |
| `role` | `users.role` |
| `status` | `users.status` |
| `exp` | expiry unix |
| `iat` | issued at |

Middleware on every protected route:

1. Validate JWT
2. Load user from DB (status may have changed)
3. Re-check status rules
4. Attach `user_id`, `role` to request context

---

## 6. Shared types (Go DTO package reference)

### Enums (match PostgreSQL)

```
UserRole:        business | influencer | admin
UserStatus:      pending_profile | pending_approval | active | rejected | suspended | blocked | removed
```

### `UserDTO` (never include password_hash)

| Field | Type | JSON |
|-------|------|------|
| ID | UUID | `id` |
| Email | string | `email` |
| PhoneNumber | *string | `phone_number` |
| Role | UserRole | `role` |
| Status | UserStatus | `status` |
| EmailVerified | bool | `email_verified` |
| PhoneVerified | bool | `phone_verified` |
| StatusReason | *string | `status_reason` |
| SuspendedUntil | *time | `suspended_until` |
| LastLoginAt | *time | `last_login_at` |
| CreatedAt | time | `created_at` |

### `TokenPairDTO`

| Field | JSON |
|-------|------|
| AccessToken | `access_token` |
| RefreshToken | `refresh_token` |
| ExpiresIn | `expires_in` (seconds) |

### `BusinessProfileSnippetDTO`

| Field | JSON |
|-------|------|
| ID | `id` |
| BusinessName | `business_name` |
| OwnerFullName | `owner_full_name` |
| IsVerified | `is_verified` |

### `InfluencerProfileSnippetDTO`

| Field | JSON |
|-------|------|
| ID | `id` |
| FullName | `full_name` |
| CreatorName | `creator_name` |
| Username | `username` (nullable) |
| IsVerified | `is_verified` |

### `AdminProfileSnippetDTO`

| Field | JSON |
|-------|------|
| ID | `id` |
| FullName | `full_name` |
| PermissionLevel | `permission_level` |

### Standard error response — `ErrorResponse`

```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Incorrect email or password.",
    "details": []
  }
}
```

`details` optional for field validation:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      { "field": "password", "message": "must be at least 8 characters" },
      { "field": "confirm_password", "message": "passwords do not match" }
    ]
  }
}
```

---

## 7. Influencer profile completion (post-auth, not signup)

Frontend `/influencer/profile` collects data **after** login. Separate API (Profile module), but auth must expose status transitions.

When influencer saves complete profile → backend sets `users.status = pending_approval`.

| Frontend field | DB destination |
|----------------|----------------|
| Creator name | `influencer_profiles.creator_name`, `full_name` |
| Bio | `influencer_profiles.bio` |
| Niches (multi) | `influencer_profile_categories` → `categories` |
| Platforms (multi) | `influencer_profile_platforms` |
| TikTok / Instagram / YouTube links | `influencer_profiles.social_links` JSONB |
| Followers | `instagram_followers`, `tiktok_followers`, `youtube_subscribers`, `total_followers` |
| Engagement rate | `influencer_profiles.engagement_rate` |
| Price min/max | `min_price_etb`, `max_price_etb` |
| Languages | `languages` TEXT[] |
| Location | `city`, `region` |

Admin approves via admin portal → `users.status = active`, `influencer_profiles.approved_at`, `approved_by`.

---

## 8. Business profile completion (post-auth)

No dedicated onboarding page yet — business lands on dashboard. Minimal profile completion API later.

Suggested fields from `business_profiles` schema:

| Field | Column |
|-------|--------|
| Business name | `business_name` |
| Owner name | `owner_full_name` |
| Description | `description` |
| City / region | `city`, `region` |
| Industry | `industry` |
| Logo | `logo_url` |

When minimum fields filled → `users.status = active` (business does not need admin approval for MVP unless you add verification gate).

---

## 9. Audit log entries (auth events)

Write to `audit_logs` for traceability:

| Event | action | entity_type | entity_id |
|-------|--------|-------------|-----------|
| Signup | `signup` | `user` | user id |
| Login | `login` | `user` | user id |
| Logout | `logout` | `user` | user id |
| Password reset | `password_reset` | `user` | user id |
| Admin blocks user | `block` | `user` | target user id |

---

## 10. Admin seed checklist

Run once in dev/staging (Go seed command or SQL):

1. Insert `users` — role `admin`, status `active`, email verified
2. Insert `admin_profiles` — full_name, permission_level `super_admin`
3. Do **not** create business/influencer profiles for admin user

**Never:**

- Allow `role: admin` on public signup
- Return admin credentials in API responses

---

## 11. Endpoint summary

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api/v1/auth/signup` | Business or influencer registration |
| POST | `/api/v1/auth/login` | All roles login |
| POST | `/api/v1/auth/refresh` | New access token |
| POST | `/api/v1/auth/logout` | Revoke session |
| GET | `/api/v1/auth/me` | Current user + profile |
| POST | `/api/v1/auth/forgot-password` | Send reset email |
| POST | `/api/v1/auth/reset-password` | Set new password |
| POST | `/api/v1/auth/send-verify-email` | Send email OTP (Bearer) |
| POST | `/api/v1/auth/verify-email` | Confirm email with OTP (Bearer) |

**Admin only:** seed + login (`make seed-admin`). Login with `"portal": "admin"`. No signup route.

---

## 12. Frontend integration checklist

When wiring Golang backend to existing frontend:

1. Replace mock signup submit → `POST /api/v1/auth/signup` with `role` from toggle
2. Replace mock login → `POST /api/v1/auth/login`
3. Store `tokens.access_token` in `access_token` localStorage key
4. Use `redirect_to` from response instead of hardcoded mock redirects
5. Handle 401 → clear token, redirect `/auth/login?next=`
6. Implement `?next=` on login page — pass as `"next"` in login JSON body
7. Map `user.status` to banners / blocking screens
8. Admin `/admin/login` → same login API, verify `role === admin`
9. Logout → call logout API + clear both `access_token` and mock user storage

---

## 13. Security notes (implementation)

- Hash passwords with argon2id (`platform/jwt/password.go`)
- Hash refresh tokens, reset tokens, and OTP before storing
- **Login lockout:** after **5 failed attempts** per email → progressive lock:
  - 1st lock: **5 minutes**
  - 2nd lock: **30 minutes**
  - 3rd+ lock: **2 hours**
  - Cleared on successful login
- **Forgot password:** max **3 requests per hour** per email (still returns 200 if email unknown)
- **Verify email OTP:** max **3 sends per hour** per user
- Normalize email to lowercase before save/lookup
- Validate `next` redirect is relative path only (no open redirect)
- Admin seed: `make seed-admin` — change password in production
- Return same error for wrong email vs wrong password on login
- `audit_logs` written on signup, login, logout

### Environment variables (backend `.env`)

| Variable | Purpose |
|----------|---------|
| `DATABASE_URL` | PostgreSQL connection |
| `JWT_SECRET`, `JWT_*_TTL` | Token signing and expiry |
| `SMTP_*` | Gmail / SMTP for emails |
| `FRONTEND_URL` | Password reset links (e.g. `http://localhost:3000`) |
| `LOG_FORMAT`, `LOG_LEVEL` | Console vs JSON logs |

### Migrations

Run order: `auth.sql` → **`auth_security.sql`** → `admin.sql` → … (`make migrate`)

Admin seed (once): `make seed-admin`

---

## 14. Audit logs (implemented)

| Event | `audit_logs.action` | When |
|-------|---------------------|------|
| Signup | `signup` | After successful signup |
| Login | `login` | After successful login |
| Logout | `logout` | After session revoked |

Stores `actor_user_id`, `target_user_id`, `ip_address`, `user_agent` where available.

---

*Document version: Adfluence V1 — aligned with implemented backend auth (OTP verify, rate limits, audit logs).*
