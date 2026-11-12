package createdadmin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

const testKey = "k3y-for-tests-0123456789abcdef-0123456789"

type fakeService struct{ got *adminDto.CreateAdminRequest }

func (f *fakeService) CreateAdmin(_ context.Context, req adminDto.CreateAdminRequest, _ authDto.SessionMeta) (string, error) {
	f.got = &req
	return "new-id", nil
}

func TestCreateAdminGuards(t *testing.T) {
	good := `{"admin_name":"Hana","admin_email":"  Hana@Adfluence.ET ","password":"Strong-pass-1","admin_role":"super_admin"}`
	cases := []struct {
		name, key, remote, header, body string
		want                            int
		wantType                        string
	}{
		{"off without a key", "", "127.0.0.1:1", testKey, good, 404, ""},
		{"off with a short key", "short-key", "127.0.0.1:1", "short-key", good, 404, ""},
		{"refused from the internet", testKey, "102.213.71.236:1", testKey, good, 404, ""},
		{"missing key", testKey, "127.0.0.1:1", "", good, 401, "invalid_api_key"},
		{"wrong key", testKey, "172.17.0.1:1", testKey + "x", good, 401, "invalid_api_key"},
		{"bad JSON", testKey, "127.0.0.1:1", testKey, `{`, 400, "invalid_request"},
		{"short password", testKey, "127.0.0.1:1", testKey, strings.Replace(good, "Strong-pass-1", "short", 1), 400, ""},
		{"bad email", testKey, "127.0.0.1:1", testKey, strings.Replace(good, "  Hana@Adfluence.ET ", "not-an-email", 1), 400, ""},
		{"wrong role", testKey, "127.0.0.1:1", testKey, strings.Replace(good, "super_admin", "god", 1), 400, ""},
		{"created from the server (Docker network)", testKey, "172.17.0.1:1", testKey, good, 201, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &fakeService{}
			h := InitCreateAdminHandler(svc, c.key, utils.New(zap.NewNop()))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/create_admin", strings.NewReader(c.body))
			req.RemoteAddr = c.remote
			if c.header != "" {
				req.Header.Set("X-Api-Key", c.header)
			}
			rec := httptest.NewRecorder()
			h.CreateAdmin(rec, req)
			if rec.Code != c.want {
				t.Fatalf("status %d, want %d (%s)", rec.Code, c.want, rec.Body.String())
			}
			if c.wantType != "" {
				var body struct {
					Error struct{ Type string } `json:"error"`
				}
				_ = json.Unmarshal(rec.Body.Bytes(), &body)
				if body.Error.Type != c.wantType {
					t.Fatalf("error type %q, want %q", body.Error.Type, c.wantType)
				}
			}
			if c.want != 201 && svc.got != nil {
				t.Fatal("the service ran although the request was refused")
			}
			if c.want == 201 && (svc.got == nil || svc.got.AdminEmail != "hana@adfluence.et") {
				t.Fatalf("service got %+v, want the email normalised", svc.got)
			}
		})
	}
}
