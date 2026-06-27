# Adfluence Postman

## Import

1. Open Postman → **Import**
2. Select both files:
   - `Adfluence.postman_collection.json`
   - `Adfluence.postman_environment.json` (optional)
3. Enable environment **Adfluence — Local**
4. Start backend: `go run cmd/main.go` (from `backend/`)

## Test flow

1. **Auth → Public → Login** (or Signup) — tokens auto-save
2. **Auth → Protected → Get Me**
3. **Auth → Protected → Send Verify Email** → set `verifyOtp` → **Verify Email**
4. **Admin → Login — Admin** (after `make seed-admin`)

## Folders

| Folder | Status |
|--------|--------|
| Auth | All 12 endpoints implemented |
| Admin | Login + Get Me via auth |
| Business | Stub (not wired in router yet) |
| Influencer | Stub (not wired in router yet) |
| Marketplace | Stub (not wired in router yet) |

## Variables

| Variable | Set by |
|----------|--------|
| `baseUrl` | Default `http://localhost:8080` |
| `accessToken` | Login / Signup / Refresh tests |
| `refreshToken` | Login / Signup / Refresh tests |
| `resetToken` | Manual — from reset email link |
| `verifyOtp` | Manual — from verification email |
