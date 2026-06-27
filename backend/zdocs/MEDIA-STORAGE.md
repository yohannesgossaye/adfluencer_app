# Adfluence Media Storage Specification

Object storage for profile images, portfolio samples (image/video), business logos, and future attachments. Uses **MinIO locally** (S3-compatible) and the same client pattern in production.

**Related docs:** `INFLUENCER-PROFILE-API.md` (upload endpoints), `AUTH-API.md` (auth only).

**Frontend integration:** out of scope for this backend phase — APIs must work via **Postman/curl** with `multipart/form-data` first.

---

## 1. Overview

### Why MinIO / S3

| Need | Storage |
|------|---------|
| Influencer avatar | `influencer_profiles.profile_image_url` + object in bucket |
| Content samples (up to 3) | `influencer_portfolio_items` + objects in bucket |
| Business logo (later) | `business_profiles.logo_url` |
| ID / license docs (later) | `verification_documents.file_url` |

PostgreSQL stores **metadata and URLs/keys only**. Binary files live in object storage.

### Design principles

1. **Backend receives multipart** in V1 (simple Postman testing, matches frontend plan).
2. **S3-compatible SDK** — same code for local MinIO and production MinIO/AWS.
3. **Object keys are opaque** — never expose raw bucket paths; use stored `object_key` for delete/replace.
4. **Public read for media** in dev via bucket policy; production may use CDN or presigned URLs later.

---

## 2. Local MinIO setup (Docker)

Add to `backend/docker-compose.yml` (not yet in repo — implement in Phase B):

```yaml
services:
  minio:
    image: minio/minio:latest
    container_name: adfluence_minio
    restart: unless-stopped
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin123
    ports:
      - "9000:9000"   # S3 API
      - "9001:9001"   # Web console → http://localhost:9001
    volumes:
      - adfluence_minio_data:/data
    healthcheck:
      test: ["CMD", "mc", "ready", "local"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  adfluence_minio_data:
```

### First-time bucket setup

1. `make up` (postgres + minio)
2. Open **http://localhost:9001** → login `minioadmin` / `minioadmin123`
3. Create bucket: **`adfluence-media`**
4. Bucket policy → **public read** for dev (or use anonymous download on prefix `media/`)

Optional init script (`docker/init/minio-init.sh`) can create bucket + policy automatically — implement when wiring compose.

---

## 3. Environment variables

Add to `backend/.env`:

```env
# MinIO / S3 (local defaults)
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY_ID=minioadmin
S3_SECRET_ACCESS_KEY=minioadmin123
S3_BUCKET_NAME=adfluence-media
S3_PUBLIC_URL=http://localhost:9000/adfluence-media
S3_REGION=us-east-1
S3_USE_PATH_STYLE=true
```

| Variable | Required | Description |
|----------|----------|-------------|
| `S3_ENDPOINT` | yes | MinIO or AWS endpoint URL |
| `S3_ACCESS_KEY_ID` | yes | Access key |
| `S3_SECRET_ACCESS_KEY` | yes | Secret key |
| `S3_BUCKET_NAME` | yes | Bucket name (`adfluence-media`) |
| `S3_PUBLIC_URL` | yes | Base URL for browser-accessible files (no trailing slash) |
| `S3_REGION` | no | Default `us-east-1` (MinIO ignores but SDK needs it) |
| `S3_USE_PATH_STYLE` | no | `true` for MinIO path-style URLs |

### Production

Same variables; point `S3_ENDPOINT` and `S3_PUBLIC_URL` at your hosted MinIO or AWS S3 + CloudFront.

When S3 is **not configured**, upload endpoints return **503** `storage_not_configured`. Profile JSON endpoints still work without storage.

---

## 4. Go storage service (reference pattern)

Adapt from existing project code. Target path: `backend/platform/storage/storage.go`.

### Interface

```go
type StorageService interface {
	UploadBytes(ctx context.Context, objectName string, data io.Reader, size int64, contentType string) (string, error)
	Check(ctx context.Context, objectName string) (bool, error)
	GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, objectName string) error
	Enabled() bool
}
```

| Method | Returns | Notes |
|--------|---------|-------|
| `UploadBytes` | `objectName` (key) | Upload to bucket; caller builds public URL via `GetFileURL` |
| `Check` | exists bool | HeadObject |
| `GetFileURL` | public URL string | Dev: `{S3_PUBLIC_URL}/{objectName}`; prod: same or presigned if expiry > 0 |
| `Delete` | error | Remove object on profile replace / portfolio delete |
| `Enabled` | bool | All required env vars set |

### S3 client config (MinIO)

```go
awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
	awsconfig.WithRegion(cfg.Region),
	awsconfig.WithBaseEndpoint(cfg.Endpoint),
	awsconfig.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	),
)

client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
	o.UsePathStyle = cfg.UsePathStyle // true for MinIO
})
```

Use `manager.NewUploader(client)` for `PutObject` via `UploadBytes`.

### Initiator wiring

Same pattern as SMTP/SMS in `initiator/initiator.go`:

- If `storage.Enabled()` → log `object storage ready`
- Else → log warn `upload endpoints disabled until S3_* env vars are set`

Inject into influencer module (and later business module).

---

## 5. Object key naming convention

Structured keys — easy to debug, enforce ownership in path:

```
media/{owner_type}/{user_id}/{purpose}/{uuid}.{ext}
```

| Segment | Values |
|---------|--------|
| `owner_type` | `influencer`, `business`, `admin` |
| `user_id` | UUID from JWT |
| `purpose` | `avatar`, `portfolio`, `logo`, `document` |
| `uuid` | New UUID per upload |
| `ext` | From validated MIME (`.jpg`, `.png`, `.webp`, `.mp4`, `.webm`) |

**Examples:**

```
media/influencer/a1b2c3d4-.../avatar/f7e8d9c0-....jpg
media/influencer/a1b2c3d4-.../portfolio/b2c3d4e5-....mp4
media/business/c4d5e6f7-.../logo/a8b9c0d1-....png
```

**Rules:**

- Never use original filename as key (collision + security).
- Store `object_key` in DB for delete/replace.
- Store `file_url` (public URL) in DB for API responses.

---

## 6. Upload validation rules

Apply in handler/module **before** calling `UploadBytes`.

### Avatar (influencer)

| Rule | Value |
|------|-------|
| Max size | **5 MB** |
| Allowed MIME | `image/jpeg`, `image/png`, `image/webp` |
| Form field name | `file` |
| Replaces | Previous avatar object deleted after successful DB update |

### Portfolio sample (influencer)

| Rule | Value |
|------|-------|
| Max count | **3** per user (enforced in DB + app) |
| Max size image | **10 MB** |
| Max size video | **100 MB** |
| Allowed image MIME | `image/jpeg`, `image/png`, `image/webp` |
| Allowed video MIME | `video/mp4`, `video/webm` |
| Form field name | `file` |
| Optional form field | `sort_order` (0–2) |

### Business logo (later)

| Rule | Value |
|------|-------|
| Max size | **5 MB** |
| Allowed MIME | `image/jpeg`, `image/png`, `image/webp` |

### Security

- Validate `Content-Type` header **and** sniff first 512 bytes (magic bytes) where practical.
- Reject executable extensions even if MIME spoofed.
- Max body size on route: use `http.MaxBytesReader` (e.g. 105 MB for portfolio video endpoint).

---

## 7. Upload flows

### 7.1 Avatar upload (multipart)

**`POST /api/v1/influencer/profile/avatar`**

```
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

file: (binary)
```

**Steps:**

1. Auth middleware → user_id, role=`influencer`
2. Parse multipart, read `file`
3. Validate size + MIME
4. Build object key → `UploadBytes`
5. Build public URL → `GetFileURL`
6. Transaction:
   - If old `profile_image_url` / old object_key exists → `Delete` old object (best effort)
   - Update `influencer_profiles.profile_image_url`
7. Return updated profile snippet

**Response `data`:**

```json
{
  "profile_image_url": "http://localhost:9000/adfluence-media/media/influencer/.../avatar/....jpg"
}
```

### 7.2 Portfolio upload (multipart)

**`POST /api/v1/influencer/profile/portfolio`**

```
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

file: (binary)
sort_order: 0
```

**Steps:**

1. Count existing portfolio rows for user; if ≥ 3 → **409** `portfolio_limit_reached`
2. Upload to MinIO
3. Insert `influencer_portfolio_items` row
4. Recalculate profile completeness (optional in response)

### 7.3 Portfolio delete

**`DELETE /api/v1/influencer/profile/portfolio/{item_id}`**

1. Verify item belongs to authenticated user
2. Delete MinIO object
3. Delete DB row
4. Return updated portfolio list

### 7.4 Future: presigned upload (V2)

For large videos in production, optional flow:

1. `POST /api/v1/influencer/profile/portfolio/presign` → returns presigned PUT URL + pending row id
2. Client uploads direct to MinIO
3. `POST /api/v1/influencer/profile/portfolio/{id}/confirm` → backend verifies object exists via `Check`

**Not in V1** — backend multipart only for now.

---

## 8. URL construction

### Dev (public bucket)

```
{S3_PUBLIC_URL}/{object_key}
→ http://localhost:9000/adfluence-media/media/influencer/{user_id}/avatar/{uuid}.jpg
```

Requires bucket policy allowing `GetObject` on `adfluence-media/*` for anonymous read in dev.

### GetFileURL implementation

Reference pattern:

```go
func (s *S3Storage) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	if expiry > 0 {
		// optional: presign for private buckets
	}
	escaped := strings.ReplaceAll(objectName, " ", "%20")
	return url.JoinPath(s.publicURL, escaped)
}
```

---

## 9. Database columns using storage

| Table | Column | Stores |
|-------|--------|--------|
| `influencer_profiles` | `profile_image_url` | Public URL |
| `influencer_portfolio_items` | `file_url` | Public URL |
| `influencer_portfolio_items` | `object_key` | MinIO key for delete |
| `influencer_portfolio_items` | `file_name`, `mime_type`, `file_size_bytes`, `media_type` | Metadata |
| `business_profiles` | `logo_url` | Public URL (later) |
| `verification_documents` | `file_url` | Public or private URL (later) |

**New migration:** `migration/influencer_portfolio.sql` — see `INFLUENCER-PROFILE-API.md` §3.

---

## 10. Error codes (storage)

| HTTP | Code | When |
|------|------|------|
| 400 | `invalid_file` | Missing file, wrong MIME, too large |
| 400 | `invalid_multipart` | Cannot parse form |
| 409 | `portfolio_limit_reached` | Already 3 samples |
| 404 | `portfolio_item_not_found` | Delete wrong id |
| 503 | `storage_not_configured` | S3 env not set |
| 500 | `upload_failed` | MinIO/S3 error |

---

## 11. Makefile commands (planned)

```bash
make up              # postgres + minio
make minio-console   # open http://localhost:9001 (document only)
make migrate-portfolio  # influencer_portfolio_items table
```

---

## 12. Postman testing (multipart)

1. Login as influencer → save `accessToken`
2. **POST** `{{baseUrl}}/api/v1/influencer/profile/avatar`
   - Body → **form-data**
   - Key `file` → type **File** → pick `.jpg`
   - Header: `Authorization: Bearer {{accessToken}}`
3. **POST** portfolio with `.mp4` or `.jpg` (repeat up to 3)
4. **GET** profile → verify `profile_image_url` and `portfolio_items[]`

---

## 13. Dependencies (Go)

Add to `go.mod` when implementing:

```
github.com/aws/aws-sdk-go-v2/config
github.com/aws/aws-sdk-go-v2/credentials
github.com/aws/aws-sdk-go-v2/service/s3
github.com/aws/aws-sdk-go-v2/feature/s3/manager
```

---

## 14. Implementation checklist

- [ ] Add MinIO to `docker-compose.yml`
- [ ] Create bucket `adfluence-media` + dev public read policy
- [ ] `platform/storage` — `StorageService` + `S3Storage` (from reference code)
- [ ] Config in `initiator/config.go` — `S3Config` struct
- [ ] Wire storage into influencer module
- [ ] Avatar + portfolio handlers with multipart parsing
- [ ] Migration `influencer_portfolio_items`
- [ ] Postman requests with file upload
- [ ] Delete old object on avatar replace / portfolio delete

---

*Document version: Adfluence V1 — backend-only phase; frontend integration deferred.*
