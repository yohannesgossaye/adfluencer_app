# Browser tests (Playwright)

A robot uses the real site in Chrome to check the main flows on every pull request:

| Test | Flow |
|---|---|
| `login.spec.ts` | Business, creator and admin each log in |
| `find-creators.spec.ts` | Search "football", rank for a campaign, invite |
| `invite.spec.ts` | Creator accepts an invite |
| `dispute.spec.ts` | Business reports a problem, admin resolves it |
| `support.spec.ts` | Creator writes, admin replies, creator answers |

CI (`.github/workflows/ci.yml`, job **E2E**) does the following:
1. Starts an empty Postgres and the API, which creates the tables.
2. Loads `backend/scripts/seed_demo_creators.sql` and `e2e/fixtures.sql`.
3. Builds the site and runs `pnpm exec playwright test`.

The database is thrown away after the run, and the live server is never used.

The test accounts are in `fixtures.sql` (password `E2e-pass-2026`).

## Run on your machine
1. Start the API on a fresh database and load the two SQL files.
2. Build the site with `NEXT_PUBLIC_API_URL` pointing to that API, and start it on port 3000.
3. Run `E2E_CHROME=chrome pnpm test:e2e`. `E2E_CHROME=chrome` uses your installed Chrome; without it, first run `pnpm exec playwright install chromium`.

The tests change the data (they accept the invite and resolve the dispute), so reload the fixtures into a fresh database before running again.
