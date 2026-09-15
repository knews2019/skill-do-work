package corehelpers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type downloadTransport func(*http.Request) (*http.Response, error)

func (transport downloadTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func TestAtomicDownloadDoesNotSendGitHubTokenToLocalhost(t *testing.T) {
	t.Setenv("GH_TOKEN", "dummy-review-token")
	t.Setenv("GITHUB_TOKEN", "fallback-token")
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if header := request.Header.Get("Authorization"); header != "" {
			t.Errorf("localhost received credential: %q", header)
		}
		_, _ = response.Write([]byte("payload"))
	}))
	defer server.Close()
	directory := t.TempDir()
	result := handleAtomicDownload(testContext(directory), []string{"--source-url", server.URL, "--target-path", filepath.Join(directory, "download")})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("download failed: %+v", result)
	}
}

func TestAtomicDownloadUntrustedURLsNeverPassCredentialsToCurl(t *testing.T) {
	t.Setenv("GH_TOKEN", "dummy-review-token")
	t.Setenv("GITHUB_TOKEN", "fallback-token")
	binDirectory := t.TempDir()
	// A real invocation still has to produce bytes, but any bearer header fails it.
	script := "#!/bin/sh\nwhile [ \"$#\" -gt 0 ]; do case \"$1\" in -H|--header) exit 99 ;; -o) output_path=\"$2\"; shift 2 ;; *) shift ;; esac; done\nprintf payload > \"$output_path\"\n"
	if err := os.WriteFile(filepath.Join(binDirectory, "curl"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDirectory+string(os.PathListSeparator)+os.Getenv("PATH"))
	for _, destination := range []string{"https://mirror.example/file", "http://github.com/file", "https://github.com.evil.example/file", "https://uploads.github.com/file", "https://github.com:8443/file"} {
		t.Run(destination, func(t *testing.T) {
			directory := t.TempDir()
			result := handleAtomicDownload(testContext(directory), []string{"--source-url", destination, "--target-path", filepath.Join(directory, "download")})
			if result.Outcome != resultmodel.OutcomeSuccess {
				t.Fatalf("anonymous download failed: %+v", result)
			}
		})
	}
}

func TestAtomicDownloadTrustedCredentialsUseRedirectSafeTransport(t *testing.T) {
	for _, primary := range []string{"preferred", ""} {
		t.Run("primary="+primary, func(t *testing.T) {
			t.Setenv("GH_TOKEN", primary)
			t.Setenv("GITHUB_TOKEN", "fallback")
			installFakeCurl(t, "exit 99")
			previous := http.DefaultTransport
			t.Cleanup(func() { http.DefaultTransport = previous })
			calls := 0
			http.DefaultTransport = downloadTransport(func(request *http.Request) (*http.Response, error) {
				calls++
				want := ""
				if calls == 1 {
					want = "Bearer fallback"
					if primary != "" {
						want = "Bearer " + primary
					}
				}
				if got := request.Header.Get("Authorization"); got != want {
					t.Errorf("hop %d authorization=%q want=%q", calls, got, want)
				}
				response := &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("payload")), Request: request}
				if calls == 1 {
					response.StatusCode = 302
					response.Header.Set("Location", "https://untrusted.github.com/file")
				}
				return response, nil
			})
			directory := t.TempDir()
			target := filepath.Join(directory, "download")
			if err := os.WriteFile(target, nil, 0o600); err != nil {
				t.Fatal(err)
			}
			result := handleAtomicDownload(testContext(directory), []string{"--source-url", "https://github.com/start", "--target-path", target})
			if result.Outcome != resultmodel.OutcomeSuccess || calls != 2 {
				t.Fatalf("result=%+v calls=%d", result, calls)
			}
			contents, err := os.ReadFile(target)
			if err != nil || string(contents) != "payload" {
				t.Fatalf("published bytes=%q error=%v", contents, err)
			}
		})
	}
}
