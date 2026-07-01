package main

import "testing"

// TestSelectBrowserOpenCommandPerPlatform asserts the opener-selection table:
// darwin → open, windows → rundll32's URL protocol handler, and any other
// GOOS (linux and other unix-likes) → xdg-open. exec.Command only builds the
// *exec.Cmd — it never launches a process — so this test never opens a real
// browser regardless of which binaries exist on the machine running it.
func TestSelectBrowserOpenCommandPerPlatform(t *testing.T) {
	const targetUrl = "http://localhost:8090"

	platformCases := []struct {
		goos        string
		wantProgram string
	}{
		{"darwin", "open"},
		{"linux", "xdg-open"},
		{"windows", "rundll32"},
		{"freebsd", "xdg-open"}, // unlisted GOOS falls to the generic unix opener
	}
	for _, platformCase := range platformCases {
		command := selectBrowserOpenCommand(platformCase.goos, targetUrl)
		if len(command.Args) == 0 || command.Args[0] != platformCase.wantProgram {
			t.Errorf("selectBrowserOpenCommand(%q, ...) program = %v, want %q", platformCase.goos, command.Args, platformCase.wantProgram)
		}
		if lastArg := command.Args[len(command.Args)-1]; lastArg != targetUrl {
			t.Errorf("selectBrowserOpenCommand(%q, ...) last arg = %q, want target URL %q", platformCase.goos, lastArg, targetUrl)
		}
	}
}

// TestSelectBrowserOpenCommandWindowsUsesUrlDllHandler pins the exact windows
// argv shape (rundll32 url.dll,FileProtocolHandler <url>) since it is the
// only platform in the table with more than one argument before the URL.
func TestSelectBrowserOpenCommandWindowsUsesUrlDllHandler(t *testing.T) {
	const targetUrl = "http://localhost:9000"
	command := selectBrowserOpenCommand("windows", targetUrl)
	wantArgs := []string{"rundll32", "url.dll,FileProtocolHandler", targetUrl}
	if len(command.Args) != len(wantArgs) {
		t.Fatalf("selectBrowserOpenCommand(\"windows\", ...) args = %v, want %v", command.Args, wantArgs)
	}
	for i, wantArg := range wantArgs {
		if command.Args[i] != wantArg {
			t.Errorf("selectBrowserOpenCommand(\"windows\", ...) args[%d] = %q, want %q", i, command.Args[i], wantArg)
		}
	}
}
