package initiator

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFilesOnlyFSServesFilesButNotDirectories(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "media", "influencer")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.png"), []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	srv := http.StripPrefix("/uploads/", http.FileServer(filesOnlyFS{http.Dir(root)}))

	cases := []struct {
		path string
		want int
	}{
		{"/uploads/media/influencer/a.png", http.StatusOK},
		{"/uploads/media/influencer/", http.StatusNotFound},
		{"/uploads/media/", http.StatusNotFound},
		{"/uploads/", http.StatusNotFound},
		{"/uploads/media/influencer/missing.png", http.StatusNotFound},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.want {
			t.Errorf("GET %s = %d, want %d", tc.path, rec.Code, tc.want)
		}
	}
}

func TestParseList(t *testing.T) {
	got := parseList(" http://localhost:3000/ , ,https://app.example.com")
	want := []string{"http://localhost:3000", "https://app.example.com"}
	if len(got) != len(want) {
		t.Fatalf("parseList = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("parseList = %v, want %v", got, want)
		}
	}
}
