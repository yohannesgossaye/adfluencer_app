# Adfluence Database Quick Reference

Beekeeper connection: **localhost:5433** · user `adfluence` · db `adfluence_db`

**API specs:** `AUTH-API.md` · `INFLUENCER-PROFILE-API.md` · `BUSINESS-API.md` · `MEDIA-STORAGE.md`

---

## Short commands (Makefile)

Run from `backend/` folder:

```bash
make up        # start postgres
make down      # stop postgres
make reset     # wipe + recreate + migrate (all data lost)
make migrate   # run SQL files again (manual)
make psql      # open terminal SQL shell
make status    # check container + table count
make tables    # list all tables (36 with influencer_portfolio_items)
make help      # show all commands
```

---

## What you see in Beekeeper sidebar

### Tables (32) — grouped by purpose

#### Auth — login & accounts (`auth.sql`)
| Table | What it stores |
|-------|----------------|
| `users` | Email, phone, password hash, role, status |
| `business_profiles` | Business name, owner, city, logo, is_verified |
| `influencer_profiles` | Creator name, bio, followers, prices, approval |
| `admin_profiles` | Admin name, permissions |
| `verification_documents` | Uploaded ID / license (1 per user) |
| `bank_accounts` | Bank or Telebirr payout details |
| `user_sessions` | Refresh tokens (login sessions) |
| `password_reset_tokens` | Forgot-password links |
| `email_verification_tokens` | Email verify links |
| `audit_logs` | Who did what (admin actions) |

#### Admin reference (`admin.sql`)
| Table | What it stores |
|-------|----------------|
| `categories` | Niches: beauty, fashion, tech… |
| `influencer_profile_categories` | Which categories each influencer belongs to |
| `influencer_profile_platforms` | TikTok / Instagram / YouTube per influencer |

#### Campaigns & matching (`campaigns.sql`)
| Table | What it stores |
|-------|----------------|
| `campaigns` | Business campaign brief, budget, dates |
| `campaign_platforms` | Platforms targeted (tiktok, instagram…) |
| `campaign_deliverables` | Content needed (reel, story, post…) |
| `campaign_recommendations` | Match scores for recommendation page |
| `campaign_invites` | Business invited influencer |

#### Chat, contracts, money (`marketplace.sql`)
| Table | What it stores |
|-------|----------------|
| `conversations` | Chat thread metadata (messages in Firebase) |
| `conversation_participants` | Business + influencer in a chat |
| `message_reports` | Reported chat messages |
| `chat_admin_access_logs` | When admin viewed a chat |
| `contracts` | Deal between business & influencer + 5%+5% fees |
| `contract_deliverables` | Agreed work items on contract |
| `contract_status_history` | Contract status changes |
| `payments` | Business paid for contract |
| `payouts` | Money sent to influencer |

#### Support & feedback (`support.sql`)
| Table | What it stores |
|-------|----------------|
| `support_tickets` | Help requests |
| `support_messages` | Messages inside a ticket |
| `disputes` | Serious contract conflicts |
| `reviews` | Star ratings after completed work |
| `notifications` | In-app alerts |

---

### Functions below tables (pgcrypto)

These are **not tables** — they come from the `pgcrypto` extension (enabled in `auth.sql`).

| Function | Simple meaning |
|----------|----------------|
| `gen_random_uuid()` | Create random UUID for new row ids |
| `gen_random_bytes()` | Random bytes (tokens, secrets) |
| `crypt` / `gen_salt` | Hash passwords (backend may use bcrypt in Go instead) |
| `digest` / `hmac` | Hash data for security |
| `encrypt` / `decrypt` | Encrypt/decrypt data |
| `armor` / `dearmor` | Encode binary as text |

You can ignore these in Beekeeper day to day — the Go backend uses them indirectly.

---

## If Beekeeper shows fewer than 32 tables

Scroll the table list — all 32 exist. Or run:

```bash
make tables
# or
make status
```

Missing from a short screenshot are usually: `admin_profiles`, `audit_logs`, `bank_accounts`, `business_profiles`, `campaign_deliverables`, `campaign_invites`, `campaign_platforms`, `campaign_recommendations`.

---

## Docker details

| Item | Value |
|------|--------|
| Container | `adfluence_postgres` |
| Port | **5433** (host) → 5432 (container) |
| User / password | `adfluence` / `adfluence123` |
| Database | `adfluence_db` |

Migrations run **automatically** on first `make up` (empty volume).  
After that, schema changes need `make reset` or run updated SQL in Beekeeper.

---

## Run SQL in Beekeeper

1. Connect to `adfluence_db`
2. New query tab
3. Paste SQL from `backend/migration/*.sql`
4. **Ctrl+Enter** or **F5** to run

Order if running by hand: `auth` → `admin` → `campaigns` → `marketplace` → `support`
