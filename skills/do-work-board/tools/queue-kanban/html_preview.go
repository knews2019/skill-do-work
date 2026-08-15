package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const htmlPreviewListenAddress = "127.0.0.1:0"

// htmlFolderPreviewManager lazily gives each folder containing a linked HTML
// file its own ephemeral loopback HTTP origin. The folder boundary makes local
// relative and root-relative assets work while the separate origin keeps active
// repository HTML away from the live board's testing-write APIs.
type htmlFolderPreviewManager struct {
	previewMutex       sync.Mutex
	previewByDirectory map[string]*htmlFolderPreviewServer
	managerClosed      bool
}

type htmlFolderPreviewServer struct {
	directoryRoot string
	baseUrl       string
	httpServer    *http.Server
	listener      net.Listener
}

func newHtmlFolderPreviewManager() *htmlFolderPreviewManager {
	return &htmlFolderPreviewManager{
		previewByDirectory: map[string]*htmlFolderPreviewServer{},
	}
}

// previewUrlForHtmlFile starts (or reuses) a read-only preview server rooted at
// the HTML file's real containing directory and returns an absolute URL for the
// file. Each directory gets a distinct port/origin; the basename is URL-escaped
// so ordinary spaces and punctuation remain valid links.
func (previewManager *htmlFolderPreviewManager) previewUrlForHtmlFile(resolvedHtmlPath string) (string, error) {
	resolvedDirectory, resolveError := filepath.EvalSymlinks(filepath.Dir(resolvedHtmlPath))
	if resolveError != nil {
		return "", fmt.Errorf("resolving HTML preview directory: %w", resolveError)
	}
	fileInfo, statError := os.Stat(resolvedHtmlPath)
	if statError != nil || !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("HTML preview target is not a regular file")
	}

	previewManager.previewMutex.Lock()
	defer previewManager.previewMutex.Unlock()
	if previewManager.managerClosed {
		return "", fmt.Errorf("HTML preview manager is closed")
	}
	if existingPreview := previewManager.previewByDirectory[resolvedDirectory]; existingPreview != nil {
		return existingPreview.baseUrl + "/" + url.PathEscape(filepath.Base(resolvedHtmlPath)), nil
	}

	previewListener, listenError := net.Listen("tcp", htmlPreviewListenAddress)
	if listenError != nil {
		return "", fmt.Errorf("binding HTML preview listener: %w", listenError)
	}
	previewHandler := &htmlFolderPreviewHandler{directoryRoot: resolvedDirectory}
	productionHandler, handlerError := newLiveBoardProductionHandler(
		htmlPreviewListenAddress, previewListener.Addr(), previewHandler)
	if handlerError != nil {
		_ = previewListener.Close()
		return "", fmt.Errorf("constructing HTML preview authority guard: %w", handlerError)
	}
	previewServer := &htmlFolderPreviewServer{
		directoryRoot: resolvedDirectory,
		baseUrl:       "http://" + previewListener.Addr().String(),
		listener:      previewListener,
		httpServer: &http.Server{
			Addr:         htmlPreviewListenAddress,
			Handler:      productionHandler,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
	previewManager.previewByDirectory[resolvedDirectory] = previewServer
	go previewManager.serveFolderPreview(previewServer)
	return previewServer.baseUrl + "/" + url.PathEscape(filepath.Base(resolvedHtmlPath)), nil
}

func (previewManager *htmlFolderPreviewManager) serveFolderPreview(previewServer *htmlFolderPreviewServer) {
	serveError := previewServer.httpServer.Serve(previewServer.listener)
	if serveError == nil || errors.Is(serveError, http.ErrServerClosed) {
		return
	}
	log.Printf("queue-kanban serve: HTML preview for %s stopped: %v", previewServer.directoryRoot, serveError)

	previewManager.previewMutex.Lock()
	if previewManager.previewByDirectory[previewServer.directoryRoot] == previewServer {
		delete(previewManager.previewByDirectory, previewServer.directoryRoot)
	}
	previewManager.previewMutex.Unlock()
}

// shutdown stops every lazily-created preview listener. The live board stops
// accepting new requests first, so no new preview can race this close path.
func (previewManager *htmlFolderPreviewManager) shutdown(shutdownContext context.Context) error {
	previewManager.previewMutex.Lock()
	previewManager.managerClosed = true
	previewServers := make([]*htmlFolderPreviewServer, 0, len(previewManager.previewByDirectory))
	for _, previewServer := range previewManager.previewByDirectory {
		previewServers = append(previewServers, previewServer)
	}
	previewManager.previewByDirectory = map[string]*htmlFolderPreviewServer{}
	previewManager.previewMutex.Unlock()

	var shutdownErrors []error
	for _, previewServer := range previewServers {
		if shutdownError := previewServer.httpServer.Shutdown(shutdownContext); shutdownError != nil {
			shutdownErrors = append(shutdownErrors,
				fmt.Errorf("shutting down HTML preview for %s: %w", previewServer.directoryRoot, shutdownError))
		}
	}
	return errors.Join(shutdownErrors...)
}

// htmlFolderPreviewHandler serves regular files from one canonical directory.
// It intentionally has no directory-listing path: a directory resolves only
// through its index.html, and every file/symlink is re-checked against the
// folder root before it is opened.
type htmlFolderPreviewHandler struct {
	directoryRoot string
}

func (previewHandler *htmlFolderPreviewHandler) ServeHTTP(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	setHtmlPreviewSecurityHeaders(responseWriter.Header())
	if httpRequest.Method != http.MethodGet && httpRequest.Method != http.MethodHead {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLoopbackRemoteAddr(httpRequest.RemoteAddr) {
		http.Error(responseWriter, "HTML previews are loopback-only", http.StatusForbidden)
		return
	}

	requestedRelativePath := filepath.FromSlash(strings.TrimPrefix(httpRequest.URL.Path, "/"))
	if requestedRelativePath == "" {
		requestedRelativePath = "index.html"
	}
	resolvedFilePath, resolveError := resolveRepoFilePath(previewHandler.directoryRoot, requestedRelativePath)
	if resolveError != nil {
		if os.IsNotExist(resolveError) {
			http.NotFound(responseWriter, httpRequest)
			return
		}
		http.Error(responseWriter, "Path is outside this HTML preview folder", http.StatusBadRequest)
		return
	}
	fileInfo, statError := os.Stat(resolvedFilePath)
	if statError == nil && fileInfo.IsDir() {
		requestedRelativePath = filepath.Join(requestedRelativePath, "index.html")
		resolvedFilePath, resolveError = resolveRepoFilePath(previewHandler.directoryRoot, requestedRelativePath)
		if resolveError != nil {
			if os.IsNotExist(resolveError) {
				http.NotFound(responseWriter, httpRequest)
				return
			}
			http.Error(responseWriter, "Path is outside this HTML preview folder", http.StatusBadRequest)
			return
		}
		fileInfo, statError = os.Stat(resolvedFilePath)
	}
	if statError != nil || !fileInfo.Mode().IsRegular() {
		http.NotFound(responseWriter, httpRequest)
		return
	}

	previewFile, openError := os.Open(resolvedFilePath)
	if openError != nil {
		log.Printf("queue-kanban serve: opening %s for HTML preview: %v", resolvedFilePath, openError)
		http.Error(responseWriter, "Internal error reading preview file", http.StatusInternalServerError)
		return
	}
	defer previewFile.Close()
	if contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(resolvedFilePath))); contentType != "" {
		responseWriter.Header().Set("Content-Type", contentType)
	}
	responseWriter.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(responseWriter, httpRequest, filepath.Base(resolvedFilePath), fileInfo.ModTime(), previewFile)
}

// setHtmlPreviewSecurityHeaders isolates active repository HTML from the board
// without restricting the resources authored into the page. The distinct
// loopback origin is the primary boundary; COOP/no-opener defense and no CORS
// keep the preview from inheriting a browser relationship to the board.
func setHtmlPreviewSecurityHeaders(responseHeader http.Header) {
	responseHeader.Set("Content-Security-Policy", "frame-ancestors 'none'")
	responseHeader.Set("Cross-Origin-Opener-Policy", "same-origin")
	responseHeader.Set("X-Content-Type-Options", "nosniff")
	responseHeader.Set("X-Frame-Options", "DENY")
	responseHeader.Set("Referrer-Policy", "no-referrer")
}
