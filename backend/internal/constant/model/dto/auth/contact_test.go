package auth

import "testing"

func TestContactRequestValidate(t *testing.T) {
	ok := ContactRequest{Name: "Selam", Email: " Selam@Example.com ", Role: "business", Message: "We want to run a TikTok campaign."}.Normalize()
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid message rejected: %v", err)
	}
	if ok.Email != "selam@example.com" {
		t.Fatalf("email not normalised: %q", ok.Email)
	}
	bad := map[string]ContactRequest{
		"no name":       {Email: "a@b.co", Role: "other", Message: "hello there friends"},
		"bad email":     {Name: "Abel", Email: "not-an-email", Role: "other", Message: "hello there friends"},
		"bad role":      {Name: "Abel", Email: "a@b.co", Role: "admin", Message: "hello there friends"},
		"short message": {Name: "Abel", Email: "a@b.co", Role: "creator", Message: "hi"},
	}
	for name, r := range bad {
		if err := r.Normalize().Validate(); err == nil {
			t.Errorf("%s: expected a validation error", name)
		}
	}
}
