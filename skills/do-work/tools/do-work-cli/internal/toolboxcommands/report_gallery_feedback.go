package toolboxcommands

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/rootedfs"
)

const (
	reportGalleryFeedbackFile = "do-work/report-gallery-feedback.json"
	reportGalleryFeedbackLock = "do-work/.report-gallery-feedback.lock"
	reportGalleryFeedbackAPI  = "/api/report-gallery-feedback"
	reportGalleryNoteLimit    = 500
)

type reportGalleryTarget struct {
	Kind       string `json:"target_kind"`
	ReportPath string `json:"report_path"`
	FigureID   string `json:"figure_id,omitempty"`
}

type reportGalleryEntry struct {
	reportGalleryTarget
	Mark string `json:"mark"`
	Note string `json:"note,omitempty"`
}

type reportGalleryFeedback struct {
	Version int                  `json:"version"`
	Entries []reportGalleryEntry `json:"entries"`
}

type reportGalleryReview struct {
	Title      string
	ReportPath string
	Href       string
	Figures    []reportGalleryFigure
}

type reportGalleryFigure struct {
	ID, Label, Href string
}

type reportGalleryLivePage struct {
	Reviews   map[string]reportGalleryReview // keyed by the renderer's report link
	Generated []reportGalleryReview
	Targets   map[reportGalleryTarget]bool
}

func (live *reportGalleryLivePage) discover(repositoryRoot string, worked []workedReportCard, bundles []reportBundle) error {
	live.Reviews = make(map[string]reportGalleryReview)
	live.Targets = make(map[reportGalleryTarget]bool)
	add := func(filename, title string) (reportGalleryReview, error) {
		review := reportGalleryReview{Title: title, ReportPath: filename, Href: reportGalleryHref(filename)}
		if err := validateNoLinkedAncestors(repositoryRoot, filename, true); err != nil {
			return review, err
		}
		if strings.EqualFold(path.Ext(filename), ".html") {
			content, err := rootedReadFile(repositoryRoot, filename)
			if err != nil {
				return review, err
			}
			review.Figures = reportGalleryFigures(content, review.Href)
		}
		live.Targets[reportGalleryTarget{Kind: "report", ReportPath: filename}] = true
		for _, figure := range review.Figures {
			live.Targets[reportGalleryTarget{Kind: "figure", ReportPath: filename, FigureID: figure.ID}] = true
		}
		live.Reviews[review.Href] = review
		return review, nil
	}
	for _, card := range worked {
		href, err := url.Parse(card.Report)
		if err != nil {
			return err
		}
		if _, err := add(path.Join(reportGalleryFolder, href.Path), card.Subject); err != nil {
			return err
		}
	}
	for _, bundle := range bundles {
		if bundle.Entry == "" {
			continue
		}
		if _, exists := live.Reviews[reportGalleryHref(bundle.Entry)]; !exists {
			review, err := add(bundle.Entry, bundle.Title)
			if err != nil {
				return err
			}
			live.Generated = append(live.Generated, review)
		}
	}
	return nil
}

const reportGalleryTagAttrs = `(?:[^'"<>]|"[^"]*"|'[^']*')*`

var (
	reportGalleryHidden = regexp.MustCompile(`(?is)<!--.*?(?:-->|$)` +
		galleryHiddenSections("script", "style", "textarea", "xmp", "iframe", "template", "noscript"))
	reportGalleryVisual = regexp.MustCompile(`(?is)<(figure)\b(` + reportGalleryTagAttrs + `)>(.*?)</figure\s*>|<(svg)\b(` + reportGalleryTagAttrs +
		`)>(.*?)</svg\s*>|<(canvas)\b(` + reportGalleryTagAttrs + `)>(.*?)</canvas\s*>|<(img)\b(` + reportGalleryTagAttrs + `)>`)
	reportGalleryVisualTag = regexp.MustCompile(`(?is)<(?:img|svg|canvas)\b(` + reportGalleryTagAttrs + `)>`)
	reportGalleryCaption   = regexp.MustCompile(`(?is)<figcaption\b` + reportGalleryTagAttrs + `>(.*?)</figcaption\s*>`)
	reportGalleryTitle     = regexp.MustCompile(`(?is)<title\b(` + reportGalleryTagAttrs + `)>(.*?)</title\s*>`)
	reportGalleryAttribute = regexp.MustCompile(`^\s+([^\s=/>]+)(?:\s*=\s*("[^"]*"|'[^']*'|[^\s"'<>=` + "`" + `]+))?`)
)

func galleryHiddenSections(tags ...string) string {
	var patterns strings.Builder
	for _, tag := range tags {
		patterns.WriteString(`|<` + tag + `\b` + reportGalleryTagAttrs + `>.*?(?:</` + tag + `\s*>|$)`)
	}
	return patterns.String()
}

func galleryVisualAttrs(attributes string) map[string]string {
	values := make(map[string]string)
	for {
		match := reportGalleryAttribute.FindStringSubmatch(attributes)
		if match == nil {
			return values
		}
		name := strings.ToLower(match[1])
		if _, exists := values[name]; !exists {
			values[name] = strings.TrimSpace(html.UnescapeString(strings.Trim(match[2], "\"'")))
		}
		attributes = attributes[len(match[0]):]
	}
}

func reportGalleryFigures(content []byte, href string) []reportGalleryFigure {
	type visualCandidate struct {
		figure reportGalleryFigure
		digest string
	}
	var candidates []visualCandidate
	seen := map[string]int{}
	visible := reportGalleryHidden.ReplaceAllString(string(content), " ")
	for _, match := range reportGalleryVisual.FindAllStringSubmatch(visible, -1) {
		// The same matcher consumes an entire semantic figure before looking for
		// standalone visuals. A figure and the SVG or image inside it are one target.
		kind, opening, body := match[1], match[2], match[3]
		if kind == "" {
			switch {
			case match[4] != "":
				kind, opening, body = match[4], match[5], match[6]
			case match[7] != "":
				kind, opening, body = match[7], match[8], match[9]
			default:
				kind, opening = match[10], match[11]
			}
		}
		kind = strings.ToLower(kind)
		outer := galleryVisualAttrs(opening)
		visual := outer
		if kind == "figure" {
			if embedded := reportGalleryVisualTag.FindStringSubmatch(body); embedded != nil {
				visual = galleryVisualAttrs(embedded[1])
			}
		}
		if outer["aria-hidden"] == "true" || (kind == "img" && outer["alt"] == "") {
			continue
		}
		label, titleID := visual["aria-label"], ""
		if title := reportGalleryTitle.FindStringSubmatch(body); title != nil {
			titleID = galleryVisualAttrs(title[1])["id"]
			if label == "" {
				label = cleanReportText(title[2])
			}
		}
		if label == "" {
			label = visual["alt"]
		}
		if outer["aria-label"] != "" {
			label = outer["aria-label"]
		}
		if label == "" {
			if caption := reportGalleryCaption.FindStringSubmatch(body); caption != nil {
				label = cleanReportText(caption[1])
			}
		}
		if kind != "figure" && label == "" && outer["role"] != "img" {
			continue
		}
		anchor := outer["id"]
		if anchor == "" {
			anchor = visual["id"]
		}
		id := anchor
		if id == "" {
			id = titleID
		}
		digest := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(strings.Fields(match[0]), " "))))[:24]
		if id == "" {
			identity := visual["src"]
			if identity == "" {
				identity = label
			}
			if identity != "" {
				id = fmt.Sprintf("figure-%x", sha256.Sum256([]byte(identity)))[:31]
			} else {
				id = "figure-" + digest
			}
		}
		link := ""
		if anchor != "" {
			link = href + (&url.URL{Fragment: anchor}).String()
		}
		if label == "" {
			label = fmt.Sprintf("Figure %d", len(candidates)+1)
		}
		seen[id]++
		candidates = append(candidates, visualCandidate{reportGalleryFigure{id, label, link}, digest})
	}
	figures := make([]reportGalleryFigure, 0, len(candidates))
	assigned := map[string]bool{}
	for _, candidate := range candidates {
		figure := candidate.figure
		if seen[figure.ID] > 1 {
			base := figure.ID + "-" + candidate.digest
			figure.ID = base
			// Reserve authored IDs even when they look like generated disambiguators.
			for suffix := 2; seen[figure.ID] != 0 || assigned[figure.ID]; suffix++ {
				figure.ID = fmt.Sprintf("%s-%d", base, suffix)
			}
		}
		assigned[figure.ID] = true
		figures = append(figures, figure)
	}
	return figures
}

func validReportGalleryEntry(entry reportGalleryEntry) bool {
	return (entry.Kind == "report" && entry.FigureID == "" || entry.Kind == "figure" && entry.FigureID != "") &&
		fs.ValidPath(entry.ReportPath) && strings.HasPrefix(entry.ReportPath, reportsFolder+"/") &&
		(entry.Mark == "like" || entry.Mark == "dislike") && utf8.RuneCountInString(entry.Note) <= reportGalleryNoteLimit
}

func readReportGalleryFeedback(repositoryRoot string, repository *os.Root) (reportGalleryFeedback, []byte, error) {
	empty := reportGalleryFeedback{Version: 1, Entries: []reportGalleryEntry{}}
	if err := validateNoLinkedAncestors(repositoryRoot, reportGalleryFeedbackFile, true); err != nil {
		return empty, nil, err
	}
	info, err := repository.Lstat(reportGalleryFeedbackFile)
	if errors.Is(err, fs.ErrNotExist) {
		return empty, nil, nil
	}
	if err != nil {
		return empty, nil, err
	}
	if !info.Mode().IsRegular() {
		return empty, nil, errors.New("feedback record must be a regular file")
	}
	original, err := rootedfs.ReadFile(repository, reportGalleryFeedbackFile)
	if err != nil {
		return empty, nil, err
	}
	var feedback reportGalleryFeedback
	if err := json.Unmarshal(original, &feedback); err != nil {
		return empty, nil, fmt.Errorf("invalid feedback JSON: %w", err)
	}
	if feedback.Version != 1 || feedback.Entries == nil {
		return empty, nil, errors.New("feedback record must have version 1 and an entries array")
	}
	seen := make(map[reportGalleryTarget]bool, len(feedback.Entries))
	for _, entry := range feedback.Entries {
		if !validReportGalleryEntry(entry) || seen[entry.reportGalleryTarget] {
			return empty, nil, errors.New("feedback record contains an invalid or duplicate target")
		}
		seen[entry.reportGalleryTarget] = true
	}
	return feedback, original, nil
}

func saveReportGalleryFeedback(repositoryRoot string, repository *os.Root, entry reportGalleryEntry) error {
	if err := rootedMkdirAll(repositoryRoot, path.Dir(reportGalleryFeedbackFile), 0o755); err != nil {
		return err
	}
	if err := validateNoLinkedAncestors(repositoryRoot, reportGalleryFeedbackLock, false); err != nil {
		return err
	}
	// The exclusive claim serializes writers from separate gallery processes.
	// The final publication also checks the snapshot for edits by other tools.
	lock, err := repository.OpenFile(reportGalleryFeedbackLock, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return errors.New("another gallery may be saving. If no gallery is saving, remove " + reportGalleryFeedbackLock)
	}
	if err != nil {
		return err
	}
	defer repository.Remove(reportGalleryFeedbackLock)
	if err := lock.Close(); err != nil {
		return err
	}

	feedback, original, err := readReportGalleryFeedback(repositoryRoot, repository)
	if err != nil {
		return err
	}
	index := len(feedback.Entries)
	for i, saved := range feedback.Entries {
		if saved.reportGalleryTarget == entry.reportGalleryTarget {
			index = i
			break
		}
	}
	if index == len(feedback.Entries) {
		feedback.Entries = append(feedback.Entries, entry)
	} else {
		feedback.Entries[index] = entry
	}
	content, err := json.MarshalIndent(feedback, "", "  ")
	if err != nil {
		return err
	}
	return rootedPublishFileChecked(repositoryRoot, reportGalleryFeedbackFile, append(content, '\n'), 0o600, original != nil, func(current []byte) error {
		if !bytes.Equal(current, original) {
			return errors.New("the feedback record changed during the save")
		}
		return nil
	})
}

func serveReportGallery(repositoryRoot string, repository *os.Root, toolboxDocs string, samples []reportGallerySample, openBrowser bool) resultmodel.CommandResult {
	bundles, _, err := reportIndexWalk(repositoryRoot)
	if err != nil {
		return reportGalleryFailure(reportsFolder, err)
	}
	live := &reportGalleryLivePage{}
	page, err := renderReportGalleryPage(repositoryRoot, repository, toolboxDocs, samples, bundles, live)
	if err != nil {
		return reportGalleryFailure(ReportGalleryPageFile, err)
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return reportGalleryFailure(ReportGalleryPageFile, err)
	}
	host := listener.Addr().String()
	address := "http://" + host + "/" + ReportGalleryPageFile
	server := &http.Server{
		Handler:           reportGalleryHTTP(repositoryRoot, repository, toolboxDocs, host, page, live.Targets),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	fmt.Println("Report gallery: " + address)
	if openBrowser {
		openReportGalleryBrowser(address)
	}
	select {
	case err := <-done:
		return reportGalleryFailure(ReportGalleryPageFile, err)
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdown); err != nil {
			_ = server.Close()
			return reportGalleryFailure(ReportGalleryPageFile, err)
		}
		<-done
		return exactOutputResult("", nil)
	}
}

func openReportGalleryBrowser(address string) {
	program := "xdg-open"
	switch runtime.GOOS {
	case "darwin":
		program = "open"
	case "windows":
		program = "rundll32"
	}
	args := []string{address}
	if program == "rundll32" {
		args = append([]string{"url.dll,FileProtocolHandler"}, args...)
	}
	command := exec.Command(program, args...)
	if err := command.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Could not open a browser automatically. Open "+address+" in your browser.")
		return
	}
	go func() {
		if err := command.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, "Could not open a browser automatically. Open "+address+" in your browser.")
		}
	}()
}

func reportGalleryHTTP(repositoryRoot string, repository *os.Root, toolboxDocs, host string, page []byte, targets map[reportGalleryTarget]bool) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Host != host {
			http.Error(writer, "Open the loopback address printed by just run-report-gallery.", http.StatusForbidden)
			return
		}
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		switch request.URL.Path {
		case reportGalleryFeedbackAPI:
			reportGalleryFeedbackRequest(writer, request, repositoryRoot, repository, "http://"+host, targets)
		case "/", "/" + ReportGalleryPageFile:
			if request.Method != http.MethodGet && request.Method != http.MethodHead {
				writer.Header().Set("Allow", "GET, HEAD")
				http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if request.URL.Path == "/" {
				http.Redirect(writer, request, "/"+ReportGalleryPageFile, http.StatusTemporaryRedirect)
				return
			}
			writer.Header().Set("Cache-Control", "no-store")
			writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			writer.Header().Set("X-Frame-Options", "DENY")
			http.ServeContent(writer, request, "index.html", time.Time{}, bytes.NewReader(page))
		default:
			serveReportGalleryAsset(writer, request, repositoryRoot, repository, toolboxDocs)
		}
	})
}

func reportGalleryFeedbackRequest(writer http.ResponseWriter, request *http.Request, repositoryRoot string, repository *os.Root, origin string, targets map[reportGalleryTarget]bool) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Cache-Control", "no-store")
	fail := func(code int, message string) {
		writer.WriteHeader(code)
		_ = json.NewEncoder(writer).Encode(map[string]string{"error": message})
	}
	if value := request.Header.Get("Origin"); value != "" && value != origin {
		fail(http.StatusForbidden, "Open the gallery from just run-report-gallery to read or save feedback.")
		return
	}
	switch request.Method {
	case http.MethodGet:
		feedback, _, err := readReportGalleryFeedback(repositoryRoot, repository)
		if err != nil {
			fail(http.StatusInternalServerError, fmt.Sprintf("Check %s: %v", reportGalleryFeedbackFile, err))
			return
		}
		_ = json.NewEncoder(writer).Encode(feedback)
	case http.MethodPut:
		if request.Header.Get("Origin") != origin {
			fail(http.StatusForbidden, "Save feedback from the opened gallery.")
			return
		}
		mediaType, _, _ := mime.ParseMediaType(request.Header.Get("Content-Type"))
		if mediaType != "application/json" {
			fail(http.StatusUnsupportedMediaType, "Feedback must be JSON.")
			return
		}
		decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 16<<10))
		decoder.DisallowUnknownFields()
		var entry reportGalleryEntry
		if err := decoder.Decode(&entry); err != nil || decoder.Decode(new(any)) != io.EOF ||
			!validReportGalleryEntry(entry) || !targets[entry.reportGalleryTarget] {
			fail(http.StatusBadRequest, "Choose a current report or figure and a mark. Keep the note to 500 characters.")
			return
		}
		if err := saveReportGalleryFeedback(repositoryRoot, repository, entry); err != nil {
			fail(http.StatusConflict, fmt.Sprintf("Check %s. %v.", reportGalleryFeedbackFile, err))
			return
		}
		_ = json.NewEncoder(writer).Encode(entry)
	default:
		writer.Header().Set("Allow", "GET, PUT")
		fail(http.StatusMethodNotAllowed, "Method not allowed.")
	}
}

func serveReportGalleryAsset(writer http.ResponseWriter, request *http.Request, repositoryRoot string, repository *os.Root, toolboxDocs string) {
	filename := strings.TrimPrefix(request.URL.Path, "/")
	inAllowedFolder := false
	for _, folder := range []string{reportsFolder, reportGalleryFolder, toolboxDocs} {
		if folder == "" || !strings.HasPrefix(filename, folder+"/") {
			continue
		}
		tail := strings.TrimPrefix(filename, folder+"/")
		inAllowedFolder = !strings.HasPrefix(tail, ".") && !strings.Contains(tail, "/.")
		break
	}
	if !inAllowedFolder || !fs.ValidPath(filename) || strings.Contains(filename, `\`) ||
		validateNoLinkedAncestors(repositoryRoot, filename, true) != nil {
		http.NotFound(writer, request)
		return
	}
	if request.Method != http.MethodGet && request.Method != http.MethodHead {
		writer.Header().Set("Allow", "GET, HEAD")
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	info, err := repository.Lstat(filename)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(writer, request)
		return
	}
	file, err := repository.Open(filename)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	defer file.Close()
	// Reports may use JavaScript for charts. Give active assets a separate,
	// opaque origin so their scripts cannot read or change gallery feedback.
	writer.Header().Set("Content-Security-Policy", "sandbox allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-downloads")
	http.ServeContent(writer, request, path.Base(filename), info.ModTime(), file)
}
