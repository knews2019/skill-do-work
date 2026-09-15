package archivefetch

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

type credentialTransport func(*http.Request) (*http.Response, error)

func (transport credentialTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func TestDownloadCredentialsStayOnTrustedHTTPSDestinations(t *testing.T) {
	for _, destination := range []string{
		"https://github.com/file", "https://api.github.com/file",
		"https://codeload.github.com/file", "https://raw.githubusercontent.com/file",
		"https://github.com:443/file",
		"http://github.com/file", "http://127.0.0.1/file",
		"https://mirror.example/file", "https://github.com.evil.example/file",
		"https://uploads.github.com/file", "https://github.com:8443/file",
		"https://github.com./file", "https://github.com@evil.example/file",
		"https://user@github.com/file",
	} {
		for _, primary := range []string{"preferred", ""} {
			t.Run(destination+"/primary="+primary, func(t *testing.T) {
				t.Setenv("GH_TOKEN", primary)
				t.Setenv("GITHUB_TOKEN", "fallback")
				want := ""
				switch destination {
				case "https://github.com/file", "https://api.github.com/file", "https://codeload.github.com/file", "https://raw.githubusercontent.com/file", "https://github.com:443/file":
					want = "Bearer fallback"
					if primary != "" {
						want = "Bearer " + primary
					}
				case "https://github.com@evil.example/file":
					want = "Basic Z2l0aHViLmNvbTo=" // Explicit URL credentials, not the environment token.
				case "https://user@github.com/file":
					want = "Basic dXNlcjo="
				}
				previous := atomicHTTPClient.Transport
				t.Cleanup(func() { atomicHTTPClient.Transport = previous })
				calls := 0
				atomicHTTPClient.Transport = credentialTransport(func(request *http.Request) (*http.Response, error) {
					calls++
					if got := request.Header.Get("Authorization"); got != want {
						t.Errorf("authorization = %q, want %q", got, want)
					}
					return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("payload")), Request: request}, nil
				})
				result := DownloadAtomic(context.Background(), destination, filepath.Join(t.TempDir(), "download"))
				if result.Err != nil || calls != 1 {
					t.Fatalf("result=%+v calls=%d", result, calls)
				}
			})
		}
	}
}

func TestRedirectsCannotCarryCredentialsOutsideTrustBoundary(t *testing.T) {
	for _, destination := range []string{
		"http://github.com/file", "https://uploads.github.com/file",
		"https://github.com:8443/file", "https://mirror.example/file",
		"https://github.com/other",
	} {
		t.Run(destination, func(t *testing.T) {
			t.Setenv("GH_TOKEN", "redirect-token")
			previous := atomicHTTPClient.Transport
			t.Cleanup(func() { atomicHTTPClient.Transport = previous })
			calls := 0
			atomicHTTPClient.Transport = credentialTransport(func(request *http.Request) (*http.Response, error) {
				calls++
				want := ""
				if calls == 1 || destination == "https://github.com/other" {
					want = "Bearer redirect-token"
				}
				if got := request.Header.Get("Authorization"); got != want {
					t.Errorf("hop %d authorization=%q want=%q", calls, got, want)
				}
				response := &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("payload")), Request: request}
				if calls == 1 {
					response.StatusCode = 302
					response.Header.Set("Location", destination)
				}
				return response, nil
			})
			result := DownloadAtomic(context.Background(), "https://github.com/start", filepath.Join(t.TempDir(), "download"))
			if result.Err != nil || calls != 2 {
				t.Fatalf("result=%+v calls=%d", result, calls)
			}
		})
	}
}

func TestGitHubCredentialTrustRejectsMissingURL(t *testing.T) {
	if TrustedGitHubURL(nil) || TrustedGitHubURL(&url.URL{}) {
		t.Fatal("missing destination was trusted")
	}
}
