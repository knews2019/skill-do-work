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
	"path/filepath"
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

	cacheMu             sync.Mutex
	cachedFileMtimes    map[string]time.Time        // absPath → last-seen mtime
	cachedBoardData     *generatedBoardData         // nil until the first request
	cachedBoardMarkdown *generatedBoardMarkdownData // rebuilt with cachedBoardData, served lazily
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
//	GET /board-data.js     → fresh board data as a JS global assignment
//	GET /board-markdown.js → raw Markdown for Copy, loaded only on demand
//	GET /file?path=…       → read-only view of one repo file (loopback-only)
//	POST /api/testing/profile → add a tester profile to do-work/testers.md
//	POST /api/testing/status  → write one REQ's testing placeholders (testing_api.go)
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
	case "/board-markdown.js":
		liveServer.serveLiveBoardMarkdownJs(responseWriter, httpRequest)
	case "/file":
		liveServer.serveRepoFileView(responseWriter, httpRequest)
	case "/api/testing/profile":
		liveServer.serveTestingProfileApi(responseWriter, httpRequest)
	case "/api/testing/status":
		liveServer.serveTestingStatusApi(responseWriter, httpRequest)
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

	// Serve mode is the only mode whose page can reach the /api/testing/*
	// write endpoints, so flag the payload (a copy — never the cached struct,
	// which generate-mode snapshots must not inherit the flag from).
	liveBoardData := *boardData
	liveBoardData.LiveTestingApi = true
	liveBoardData.LiveFileApi = true

	jsText, encodeErr := encodeBoardDataForJsAssignment(liveBoardData)
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

// serveLiveBoardMarkdownJs emits the raw REQ/UR bodies used by the Copy button.
// The main page does not request this route; board.js loads it on the first copy.
func (liveServer *liveBoardServer) serveLiveBoardMarkdownJs(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet && httpRequest.Method != http.MethodHead {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	markdownData, refreshError := liveServer.refreshBoardMarkdownData()
	if refreshError != nil {
		log.Printf("queue-kanban serve: refreshing raw Markdown data: %v", refreshError)
		http.Error(responseWriter, "Internal error building Markdown data", http.StatusInternalServerError)
		return
	}
	jsText, encodeError := encodeBoardMarkdownForJsAssignment(*markdownData)
	if encodeError != nil {
		log.Printf("queue-kanban serve: encoding raw Markdown data: %v", encodeError)
		http.Error(responseWriter, "Internal error encoding Markdown data", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(jsText))
}

// repoFileViewMaxBytes caps what GET /file will serve. The endpoint exists so
// the drawer's file-path links can open the referenced doc/spec/prime file —
// those are small text files; anything larger is not what the link is for.
const repoFileViewMaxBytes = 2 << 20

// serveRepoFileView serves one repo file as read-only plain text so file-path
// mentions in REQ/UR bodies can be real links. Guards, in order: loopback
// callers only (a LAN-exposed board must not turn into a whole-repo file
// reader — the exposure warning promises REQ bodies, nothing more), then repo
// containment via resolveRepoFilePath, then a regular-file + size check.
// Always text/plain (with the global nosniff header), never the file's own
// content type, so a crafted HTML/SVG file cannot execute in the board origin.
func (liveServer *liveBoardServer) serveRepoFileView(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet && httpRequest.Method != http.MethodHead {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLoopbackRemoteAddr(httpRequest.RemoteAddr) {
		http.Error(responseWriter, "The file view is loopback-only", http.StatusForbidden)
		return
	}
	requestedPath := httpRequest.URL.Query().Get("path")
	if requestedPath == "" {
		http.Error(responseWriter, "Missing ?path=<repo-relative file>", http.StatusBadRequest)
		return
	}
	resolvedFilePath, resolveErr := resolveRepoFilePath(liveServer.repoRoot, requestedPath)
	if resolveErr != nil {
		if os.IsNotExist(resolveErr) {
			http.Error(responseWriter, "File not found in this repository: "+requestedPath, http.StatusNotFound)
			return
		}
		http.Error(responseWriter, "Path is outside this repository", http.StatusBadRequest)
		return
	}
	fileInfo, statErr := os.Stat(resolvedFilePath)
	if statErr != nil || !fileInfo.Mode().IsRegular() {
		http.Error(responseWriter, "File not found in this repository: "+requestedPath, http.StatusNotFound)
		return
	}
	if fileInfo.Size() > repoFileViewMaxBytes {
		http.Error(responseWriter, "File exceeds the board file-view size cap", http.StatusRequestEntityTooLarge)
		return
	}
	fileBytes, readErr := os.ReadFile(resolvedFilePath)
	if readErr != nil {
		log.Printf("queue-kanban serve: reading %s for /file: %v", resolvedFilePath, readErr)
		http.Error(responseWriter, "Internal error reading file", http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write(fileBytes)
}

// resolveRepoFilePath maps a repo-relative path from a drawer file link to an
// absolute path, refusing anything that escapes the repo root. Absolute paths
// and ".." traversal are rejected before touching the filesystem; symlinks are
// then resolved and the REAL location re-checked, so a symlink inside the repo
// cannot be used as a hop to files outside it. A missing file surfaces as an
// os.IsNotExist error so the caller can 404 instead of 400.
func resolveRepoFilePath(repoRoot string, requestedPath string) (string, error) {
	if filepath.IsAbs(requestedPath) {
		return "", fmt.Errorf("absolute path rejected")
	}
	cleanedRelativePath := filepath.Clean(requestedPath)
	if cleanedRelativePath == ".." || strings.HasPrefix(cleanedRelativePath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes the repo root")
	}
	resolvedRepoRoot, rootResolveErr := filepath.EvalSymlinks(repoRoot)
	if rootResolveErr != nil {
		return "", rootResolveErr
	}
	resolvedFilePath, fileResolveErr := filepath.EvalSymlinks(filepath.Join(resolvedRepoRoot, cleanedRelativePath))
	if fileResolveErr != nil {
		return "", fileResolveErr
	}
	containmentPath, relErr := filepath.Rel(resolvedRepoRoot, resolvedFilePath)
	if relErr != nil || containmentPath == ".." || strings.HasPrefix(containmentPath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes the repo root")
	}
	return resolvedFilePath, nil
}

// isLoopbackRemoteAddr reports whether an http.Request.RemoteAddr is a
// loopback peer. An unparseable address counts as non-loopback — the guard
// fails closed.
func isLoopbackRemoteAddr(remoteAddr string) bool {
	remoteHost, _, splitErr := net.SplitHostPort(remoteAddr)
	if splitErr != nil {
		return false
	}
	parsedRemoteIp := net.ParseIP(remoteHost)
	return parsedRemoteIp != nil && parsedRemoteIp.IsLoopback()
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
	boardMarkdownData := buildGeneratedBoardMarkdownData(board)
	liveServer.cachedBoardMarkdown = &boardMarkdownData
	return liveServer.cachedBoardData, nil
}

// refreshBoardMarkdownData shares refreshBoardData's tree walk and mtime cache,
// then returns the matching raw-source payload. The payload is immutable after a
// rebuild, so returning its cached pointer after releasing the lock is safe.
func (liveServer *liveBoardServer) refreshBoardMarkdownData() (*generatedBoardMarkdownData, error) {
	if _, refreshError := liveServer.refreshBoardData(); refreshError != nil {
		return nil, refreshError
	}
	liveServer.cacheMu.Lock()
	defer liveServer.cacheMu.Unlock()
	if liveServer.cachedBoardMarkdown == nil {
		return nil, fmt.Errorf("raw Markdown cache was not built with board data")
	}
	return liveServer.cachedBoardMarkdown, nil
}

// buildTreeMtimeFingerprint stats every file discovered by enumerateDoWorkTree
// — REQ files, UR input.md files, notes.md, and testers.md — and returns a map
// of absPath → mtime. Files that cannot be stat'd are omitted (consistent with the
// best-effort walk contract in walk.go). notes.md must be fingerprinted like
// any other input: it feeds the rendered board, so appending a note has to
// invalidate the cache or the live server would keep serving the old strip.
func buildTreeMtimeFingerprint(discovered discoveredTreeFiles) map[string]time.Time {
	fingerprint := make(map[string]time.Time, len(discovered.RequestFiles)+len(discovered.UserRequestFiles)+1)
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
	if discovered.NotesFilePath != "" {
		if fileInfo, statErr := os.Stat(discovered.NotesFilePath); statErr == nil {
			fingerprint[discovered.NotesFilePath] = fileInfo.ModTime()
		}
	}
	// testers.md feeds the testing view's profile picker, so adding a profile
	// must invalidate the cached board data just like a REQ edit does.
	if discovered.TestersFilePath != "" {
		if fileInfo, statErr := os.Stat(discovered.TestersFilePath); statErr == nil {
			fingerprint[discovered.TestersFilePath] = fileInfo.ModTime()
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

// describeListenExposure classifies listenAddress for the startup banner. It
// returns the address to print in the board URL plus, when the bind is
// reachable beyond loopback, a warning line. The board data carries every
// rendered REQ/UR body, so *every* non-loopback bind must warn — not just the
// host-less ":port" spelling: the wildcard hosts ("0.0.0.0", "[::]") bind all
// interfaces exactly like ":port" does, and an explicit LAN IP or machine
// hostname is likewise network-reachable. Loopback binds (127.0.0.1 and
// friends, ::1, "localhost") stay silent. An address SplitHostPort cannot
// parse gets no warning — net.Listen rejects it before anything is announced.
func describeListenExposure(listenAddress string) (displayAddress string, exposureWarning string) {
	allInterfacesWarning := fmt.Sprintf(
		"queue-kanban: warning: %s binds ALL interfaces — the board (including REQ bodies) is reachable from the local network",
		listenAddress)
	if strings.HasPrefix(listenAddress, ":") {
		return "localhost" + listenAddress, allInterfacesWarning
	}
	bindHost, _, splitError := net.SplitHostPort(listenAddress)
	if splitError != nil {
		return listenAddress, ""
	}
	switch {
	case bindHost == "" || bindHost == "0.0.0.0" || bindHost == "::":
		return listenAddress, allInterfacesWarning
	case isLoopbackBindHost(bindHost):
		return listenAddress, ""
	default:
		return listenAddress, fmt.Sprintf(
			"queue-kanban: warning: %s binds a non-loopback address — the board (including REQ bodies) is reachable from the network",
			listenAddress)
	}
}

// isLoopbackBindHost reports whether a bind host is loopback: the literal
// "localhost" or any IP the net package classifies as loopback. A hostname
// that is not "localhost" returns false — it may resolve anywhere, so the
// caller treats it as exposed.
func isLoopbackBindHost(bindHost string) bool {
	if strings.EqualFold(bindHost, "localhost") {
		return true
	}
	parsedBindIp := net.ParseIP(bindHost)
	return parsedBindIp != nil && parsedBindIp.IsLoopback()
}

// bindServeListenerAndAnnounce binds listenAddress via TCP and, ONLY on a
// successful bind, prints the startup banner (and the non-loopback exposure warning,
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

	displayAddress, exposureWarning := describeListenExposure(listenAddress)
	if exposureWarning != "" {
		fmt.Println(exposureWarning)
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
