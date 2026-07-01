package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// kanbanServePortEnvVar is the environment variable that overrides --port for
// the serve subcommand. Named distinctly from SA_PORT (the main goserver) so
// the two servers can coexist in the same shell session without a collision.
const kanbanServePortEnvVar = "QUEUE_KANBAN_PORT"

// kanbanServeDefaultListenAddress is the default listen address for
// queue-kanban serve — loopback only, because the board data carries every
// rendered REQ/UR body and must not be LAN-readable unless the user asks for
// it with an explicit host. Port 8090 avoids the main goserver :8080.
const kanbanServeDefaultListenAddress = "127.0.0.1:8090"

// liveBoardServer is the http.Handler for queue-kanban serve. It holds the
// mtime cache that keeps repeated requests cheap: each request stats the
// do-work tree and only rebuilds the board model when at least one file's
// mtime has changed since the last build.
type liveBoardServer struct {
	repoRoot     string
	recentWindow time.Duration

	cacheMu          sync.Mutex
	cachedFileMtimes map[string]time.Time // absPath → last-seen mtime
	cachedBoardData  *generatedBoardData  // nil until the first request
}

// newLiveBoardServer creates a liveBoardServer for the given repoRoot and
// recentWindow. The mtime cache starts empty; the first request triggers a
// full tree walk and board build.
func newLiveBoardServer(repoRoot string, recentWindow time.Duration) *liveBoardServer {
	return &liveBoardServer{
		repoRoot:         repoRoot,
		recentWindow:     recentWindow,
		cachedFileMtimes: map[string]time.Time{},
	}
}

// ServeHTTP dispatches HTTP requests for the live board. Routes:
//
//	GET /              → board HTML shell (same template as generate's index.html)
//	GET /board-data.js → fresh board data as a JS global assignment
//
// Every other path returns 404. Security headers are set on all responses,
// mirroring the goserver pattern (REQ-782).
func (liveServer *liveBoardServer) ServeHTTP(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	setKanbanSecurityHeaders(responseWriter.Header())

	requestPath := httpRequest.URL.Path
	// Strip trailing slashes except for the root "/".
	if requestPath != "/" && strings.HasSuffix(requestPath, "/") {
		requestPath = strings.TrimRight(requestPath, "/")
	}

	switch requestPath {
	case "/":
		liveServer.serveBoardHtml(responseWriter, httpRequest)
	case "/board-data.js":
		liveServer.serveLiveBoardDataJs(responseWriter, httpRequest)
	default:
		http.NotFound(responseWriter, httpRequest)
	}
}

// serveBoardHtml serves the HTML shell assembled from the embedded template.
// The board data is NOT inlined here — the template.html already carries a
// <script src="board-data.js"> tag that the browser follows to /board-data.js.
// Passing time.Now() to assembleStaticPage makes the "Generated …" timestamp
// reflect when the page was requested, not when the binary started.
func (liveServer *liveBoardServer) serveBoardHtml(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet && httpRequest.Method != http.MethodHead {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pageHtml, assembleErr := assembleStaticPage(time.Now(), deriveProjectName(liveServer.repoRoot))
	if assembleErr != nil {
		log.Printf("queue-kanban serve: assembling board HTML: %v", assembleErr)
		http.Error(responseWriter, "Internal error assembling board HTML", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(pageHtml))
}

// serveLiveBoardDataJs re-walks the do-work tree (using the mtime cache to
// skip re-parsing when no file changed) and emits the fresh board model as:
//
//	window.queueKanbanBoardData = <JSON>;
//
// This is the same format generateStaticSite writes to board-data.js so the
// frontend JavaScript is unchanged across the generate and serve modes.
func (liveServer *liveBoardServer) serveLiveBoardDataJs(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet && httpRequest.Method != http.MethodHead {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	boardData, refreshErr := liveServer.refreshBoardData()
	if refreshErr != nil {
		log.Printf("queue-kanban serve: refreshing board data: %v", refreshErr)
		http.Error(responseWriter, "Internal error building board data", http.StatusInternalServerError)
		return
	}

	jsText, encodeErr := encodeBoardDataForJsAssignment(*boardData)
	if encodeErr != nil {
		log.Printf("queue-kanban serve: encoding board data: %v", encodeErr)
		http.Error(responseWriter, "Internal error encoding board data", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(jsText))
}

// refreshBoardData checks whether any file in the do-work tree has changed
// since the last build (by comparing mtime fingerprints). Returns the cached
// board data unchanged if the tree is clean; otherwise rebuilds via the shared
// buildBoard + buildGeneratedBoardData path that generate uses.
func (liveServer *liveBoardServer) refreshBoardData() (*generatedBoardData, error) {
	liveServer.cacheMu.Lock()
	defer liveServer.cacheMu.Unlock()

	discovered, enumerateErr := enumerateDoWorkTree(liveServer.repoRoot)
	if enumerateErr != nil {
		return nil, fmt.Errorf("enumerating do-work tree: %w", enumerateErr)
	}

	currentFileMtimes := buildTreeMtimeFingerprint(discovered)

	if liveServer.cachedBoardData != nil && treeMtimeFingerprintsEqual(liveServer.cachedFileMtimes, currentFileMtimes) {
		// The tree is unchanged but time has moved on: serve a copy with a fresh
		// GeneratedAt, because the client computes the Recently-done cutoff from
		// generatedAt as "now" — returning the frozen build instant would stop
		// completed items from ever aging out of the 24h/48h/7d window while the
		// page header (stamped per request) claims the data is current.
		refreshedBoardData := *liveServer.cachedBoardData
		refreshedBoardData.GeneratedAt = formatTimestamp(time.Now())
		return &refreshedBoardData, nil
	}

	// Cache miss (or first request): rebuild using the same path as generate.
	// The real git lookup is used (best-effort; a failed lookup leaves the
	// completion undated), matching generate's behavior exactly.
	board, buildErr := buildBoard(liveServer.repoRoot, time.Now(), liveServer.recentWindow, lookupGitCommitDate)
	if buildErr != nil {
		return nil, fmt.Errorf("building board model: %w", buildErr)
	}

	boardData, projectErr := buildGeneratedBoardData(board)
	if projectErr != nil {
		return nil, fmt.Errorf("projecting board data: %w", projectErr)
	}

	liveServer.cachedFileMtimes = currentFileMtimes
	liveServer.cachedBoardData = &boardData
	return liveServer.cachedBoardData, nil
}

// buildTreeMtimeFingerprint stats every file discovered by enumerateDoWorkTree
// and returns a map of absPath → mtime. Files that cannot be stat'd are omitted
// (consistent with the best-effort walk contract in walk.go).
func buildTreeMtimeFingerprint(discovered discoveredTreeFiles) map[string]time.Time {
	fingerprint := make(map[string]time.Time, len(discovered.RequestFiles)+len(discovered.UserRequestFiles))
	for _, ref := range discovered.RequestFiles {
		if fileInfo, statErr := os.Stat(ref.AbsolutePath); statErr == nil {
			fingerprint[ref.AbsolutePath] = fileInfo.ModTime()
		}
	}
	for _, urFilePath := range discovered.UserRequestFiles {
		if fileInfo, statErr := os.Stat(urFilePath); statErr == nil {
			fingerprint[urFilePath] = fileInfo.ModTime()
		}
	}
	return fingerprint
}

// treeMtimeFingerprintsEqual reports whether two mtime fingerprints represent
// the same set of files with the same modification times. It returns false
// whenever a file has been added, removed, or modified.
func treeMtimeFingerprintsEqual(previous map[string]time.Time, current map[string]time.Time) bool {
	if len(previous) != len(current) {
		return false
	}
	for filePath, previousMtime := range previous {
		currentMtime, fileExists := current[filePath]
		if !fileExists || !previousMtime.Equal(currentMtime) {
			return false
		}
	}
	return true
}

// setKanbanSecurityHeaders sets defense-in-depth security response headers,
// mirroring the goserver pattern from backend/goserver/server.go (REQ-782).
func setKanbanSecurityHeaders(responseHeader http.Header) {
	responseHeader.Set(
		"Content-Security-Policy",
		"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data:; "+
			"connect-src *; "+
			"frame-ancestors 'none'",
	)
	responseHeader.Set("X-Content-Type-Options", "nosniff")
	responseHeader.Set("X-Frame-Options", "DENY")
	responseHeader.Set("Referrer-Policy", "no-referrer")
}

// resolveServeListenAddress determines the listen address from the --port flag,
// the QUEUE_KANBAN_PORT env var, or the hardcoded default 127.0.0.1:8090. The
// flag takes priority over the env var. Bare port numbers (no colon) are
// accepted from either source and bind loopback only; a value containing a
// colon (host:port, or a host-less ":port" meaning all interfaces) passes
// through verbatim — LAN exposure is always an explicit choice.
func resolveServeListenAddress(portFlagValue string) string {
	if portFlagValue != "" {
		if !strings.Contains(portFlagValue, ":") {
			return "127.0.0.1:" + portFlagValue
		}
		return portFlagValue
	}
	envValue := os.Getenv(kanbanServePortEnvVar)
	if envValue != "" {
		if !strings.Contains(envValue, ":") {
			return "127.0.0.1:" + envValue
		}
		return envValue
	}
	return kanbanServeDefaultListenAddress
}

// bindServeListenerAndAnnounce binds listenAddress via TCP and, ONLY on a
// successful bind, prints the startup banner (and the all-interfaces warning,
// when applicable) and — if openAfterBind is true — invokes browserOpener
// with the printed board URL. This ordering is deliberate: binding before
// announcing is the fix for the false-positive "live board at …" line a port
// collision used to print (immediately followed by a fatal bind error) —
// nothing is announced and no browser opens for a server that never came up.
// browserOpener is injected so tests can assert it fires (or doesn't) without
// launching a real browser.
func bindServeListenerAndAnnounce(listenAddress string, repoRoot string, openAfterBind bool, browserOpener func(string)) (net.Listener, error) {
	listener, listenErr := net.Listen("tcp", listenAddress)
	if listenErr != nil {
		return nil, listenErr
	}

	displayAddress := listenAddress
	if strings.HasPrefix(displayAddress, ":") {
		displayAddress = "localhost" + displayAddress
		fmt.Printf("queue-kanban: warning: %s binds ALL interfaces — the board (including REQ bodies) is reachable from the local network\n", listenAddress)
	}
	boardUrl := fmt.Sprintf("http://%s", displayAddress)
	fmt.Printf("queue-kanban: live board at %s  (re-walks do-work/ per /board-data.js request)\n", boardUrl)
	fmt.Printf("queue-kanban: serving do-work/ from %s\n", repoRoot)

	if openAfterBind {
		browserOpener(boardUrl)
	}

	return listener, nil
}

// runServeCommand starts the live local board server, resolving the repo root
// and listen address from flags/env, then serving until SIGINT or SIGTERM.
// The shutdown sequence mirrors the goserver bootstrap in
// backend/goserver/server.go.
func runServeCommand(args []string) {
	flagSet := flag.NewFlagSet("serve", flag.ExitOnError)
	portFlag := flagSet.String("port", "",
		fmt.Sprintf("listen port (default %s; override: %s env var; bare ports bind loopback, use host:port for LAN exposure)", kanbanServeDefaultListenAddress, kanbanServePortEnvVar))
	repoRootFlag := flagSet.String("repo-root", "",
		"repo root containing do-work/ (default: walk up from the working directory)")
	openFlag := flagSet.Bool("open", false,
		"open the default browser at the board URL after a successful bind")
	_ = flagSet.Parse(args)

	repoRoot, resolveErr := resolveRepoRootOrDefault(*repoRootFlag)
	if resolveErr != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban serve:", resolveErr)
		os.Exit(1)
	}

	listenAddress := resolveServeListenAddress(*portFlag)
	liveServer := newLiveBoardServer(repoRoot, defaultRecentWindow)

	listener, listenErr := bindServeListenerAndAnnounce(listenAddress, repoRoot, *openFlag, openBrowser)
	if listenErr != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban serve:", listenErr)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:         listenAddress,
		Handler:      liveServer,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	shutdownComplete := make(chan struct{})
	go func() {
		receivedSignal := <-shutdownChannel
		log.Printf("Received %s, shutting down gracefully...", receivedSignal)
		shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelShutdown()
		if shutdownErr := httpServer.Shutdown(shutdownContext); shutdownErr != nil {
			log.Printf("queue-kanban serve: graceful shutdown error: %v", shutdownErr)
		}
		close(shutdownComplete)
	}()

	serverErr := httpServer.Serve(listener)
	if serverErr != nil && serverErr != http.ErrServerClosed {
		log.Fatalf("queue-kanban serve: server failed: %v", serverErr)
	}
	<-shutdownComplete
}
