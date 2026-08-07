package main

import (
	"log"
	"os/exec"
	"runtime"
)

// selectBrowserOpenCommand builds the *exec.Cmd that opens targetUrl in the
// platform's default browser, selecting the opener by goos ("darwin" |
// "windows" | anything else, treated as a generic xdg-open-capable unix like
// linux). Accepting goos as a parameter — rather than reading runtime.GOOS
// directly — is the injectable seam: tests can exercise every branch of the
// selection table on any host and assert the constructed argv, without
// starting a process (exec.Command only builds the *exec.Cmd; it launches
// nothing until Start()/Run() is called).
func selectBrowserOpenCommand(goos string, targetUrl string) *exec.Cmd {
	switch goos {
	case "darwin":
		return exec.Command("open", targetUrl)
	case "windows":
		// rundll32's URL protocol handler avoids the quoting/parsing quirks of
		// routing through `cmd /c start`.
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", targetUrl)
	default:
		return exec.Command("xdg-open", targetUrl)
	}
}

// openBrowser launches the platform's default browser at targetUrl in the
// background: Start(), never Run() — a slow or missing browser opener must
// never block the server. targetUrl is always self-constructed from the
// resolved listen address (see bindServeListenerAndAnnounce in serve.go),
// never raw user input, so no shell is invoked and there is no
// argv-injection surface — mirroring the argv-safety instinct behind the git
// command construction in model.go's lookupGitCommitDate. A failed launch is
// logged as a warning; it is never fatal, so the board keeps serving either
// way.
func openBrowser(targetUrl string) {
	command := selectBrowserOpenCommand(runtime.GOOS, targetUrl)
	if startErr := command.Start(); startErr != nil {
		log.Printf("queue-kanban serve: could not open browser at %s: %v", targetUrl, startErr)
	}
}
