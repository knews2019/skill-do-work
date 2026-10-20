package toolboxcommands

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/ownedprocess"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const secondOpinionUsage = "ai-report-second-opinion [--backend codex|gemini|agy|claude [--model <slug>] | --all --exclude <own-cli> | --prepare --backend <cli> --model <slug>] [--preset <name>] [--editor-direction <text>] [<bundle-dir>|<page.html>|latest] | --finalize <staging-page>"

// secondOpinionPresetsFolder is the folder name for presets, both in a report bundle and beside the shipped prompt.
const secondOpinionPresetsFolder = "second-opinion-presets"

var secondOpinionModelPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// secondOpinionStagePattern is the shape of a launcher-reserved staging page name, so --finalize only ever publishes a page prepare created.
var secondOpinionStagePattern = regexp.MustCompile(`^\.[A-Za-z0-9][A-Za-z0-9._-]*-staging-[0-9]+\.html$`)

// secondOpinionTreePagePattern is the whole basename of an iteration-tree page; capture 1 is its tree id. Parent, depth
// and order come from these names only. `c` followed by a letter (codex, claude) cannot extend the id.
var secondOpinionTreePagePattern = regexp.MustCompile(`^index-(v1(?:c[0-9]{2})+)-[A-Za-z0-9][A-Za-z0-9._-]*\.html$`)

// errSecondOpinionModelUnverified marks native evidence that names no single backend model; a finished page is then published under an unverified name instead of being discarded.
var errSecondOpinionModelUnverified = errors.New("actual model is absent or ambiguous")
var secondOpinionJudge = secondOpinionRender

// secondOpinionProgressSink receives the launcher's live progress while a backend or the judge runs: heartbeats and the
// backend's mirrored output lines. stdout stays the typed result only; tests replace the sink and shorten the interval.
var secondOpinionProgressSink io.Writer = os.Stderr
var secondOpinionHeartbeatInterval = 30 * time.Second

// secondOpinionProgressMutex keeps whole lines whole: the heartbeat and both output copiers share the sink.
var secondOpinionProgressMutex sync.Mutex

// secondOpinionProgress names what a heartbeat reports; pagePath stays empty until a staging page exists.
type secondOpinionProgress struct {
	phase, backend, pagePath string
}

func secondOpinionProgressLine(line string) {
	secondOpinionProgressMutex.Lock()
	defer secondOpinionProgressMutex.Unlock()
	_, _ = io.WriteString(secondOpinionProgressSink, line)
}

// secondOpinionHeartbeat reports elapsed time, phase, backend and page size every interval until the returned stop
// runs. A stalled backend shows as heartbeats whose page size no longer changes. stop returns after the ticker exits.
func secondOpinionHeartbeat(progress secondOpinionProgress) (stop func()) {
	started := time.Now()
	stopped, finished := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(finished)
		ticker := time.NewTicker(secondOpinionHeartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stopped:
				return
			case <-ticker.C:
				page := "not written yet"
				if info, err := os.Stat(progress.pagePath); progress.pagePath != "" && err == nil && info.Size() > 0 {
					page = fmt.Sprintf("%d bytes", info.Size())
				}
				secondOpinionProgressLine(fmt.Sprintf("[second-opinion] %s %s running %s, page %s\n", progress.backend, progress.phase, time.Since(started).Round(time.Second), page))
			}
		}
	}()
	return func() { close(stopped); <-finished }
}

// secondOpinionBacklinkRequirement fills {{BACKLINK_REQUIREMENT}} for both the exec path and the session prepare.
const secondOpinionBacklinkRequirement = `Open the body with <p>Original report: <a href="index.html">index.html</a></p>, before the main content. Never modify index.html.`

// secondOpinionChangeListRequirement fills {{CHANGE_LIST_REQUIREMENT}} on a tree run; a plain run gets "None for this pass.".
const secondOpinionChangeListRequirement = `Right after the backlink, add a list titled "Changed in this iteration". Each item names one change you made to the base page and links to #rev-<n> (rev-1, rev-2 and so on). Give each changed section the matching id="rev-<n>" and a visible "changed" badge. Every link must resolve to a section on this page, because the launcher checks it.`

// secondOpinionChangeListKind is the judge finding kind for a change list whose #rev-* links do not resolve.
const secondOpinionChangeListKind = "AI-REPORT-SECOND-OPINION-CHANGE-LIST"

type secondOpinionOptions struct {
	backend, model, exclude, bundle, finalize string
	preset, editorDirection                   string
	all, prepare                              bool
}

// secondOpinionDirection is what the person directing this pass asked for: a named preset, free text, or both, and the
// base page to rewrite. It travels in the prepare record to finalize; the zero value is a plain run.
type secondOpinionDirection struct {
	Preset          string `json:"preset,omitempty"`
	PresetFile      string `json:"preset_file,omitempty"`
	PresetText      string `json:"preset_text,omitempty"`
	EditorDirection string `json:"editor_direction,omitempty"`
	// Base is the bare filename of a tree page to build on; empty means index.html.
	Base string `json:"base,omitempty"`
}

// parentTreeID is the tree id the new page hangs under: empty on a plain run, v1 for index.html, else the base's id.
func (direction secondOpinionDirection) parentTreeID() string {
	if match := secondOpinionTreePagePattern.FindStringSubmatch(direction.Base); match != nil {
		return match[1]
	}
	if direction.filenameSlot() != "" {
		return "v1"
	}
	return ""
}

// combinedText is the direction the editor receives: the preset text, a blank line, then the free text.
func (direction secondOpinionDirection) combinedText() string {
	parts := []string{}
	for _, part := range []string{direction.PresetText, direction.EditorDirection} {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, "\n\n")
}

// filenameSlot names the variation in the published filename: the preset, or custom for free text only.
func (direction secondOpinionDirection) filenameSlot() string {
	if direction.Preset != "" {
		return direction.Preset
	}
	if direction.EditorDirection != "" {
		return "custom"
	}
	return ""
}

func handleAIReportSecondOpinion(execution commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	options, err := secondOpinionParse(args)
	if err != nil {
		return usageResult(CommandAIReportSecondOpinion, err.Error()+"; "+secondOpinionUsage)
	}
	var bundle, stageName string
	var template []byte
	direction := secondOpinionDirection{Preset: options.preset, EditorDirection: options.editorDirection}
	if options.finalize != "" {
		// The staging page names its bundle; the current-revision chain must not redirect finalize to a newer bundle.
		stagePath := judgeAbsolutePath(execution.RepositoryRoot, options.finalize)
		bundle, err = secondOpinionBundle(execution.RepositoryRoot, filepath.Dir(stagePath))
		if err == nil && bundle != filepath.Dir(stagePath) {
			err = errors.New("staging page must sit in a current report bundle")
		}
		if err != nil {
			return usageResult(CommandAIReportSecondOpinion, err.Error())
		}
		stageName = filepath.Base(stagePath)
	} else {
		if strings.HasSuffix(options.bundle, ".html") {
			// A page target names both the bundle and the base page; like --finalize, it is never moved to a successor bundle.
			pagePath := judgeAbsolutePath(execution.RepositoryRoot, options.bundle)
			bundle, err = secondOpinionBundle(execution.RepositoryRoot, filepath.Dir(pagePath))
			if err == nil && bundle != filepath.Dir(pagePath) {
				err = errors.New("base page must sit in a current report bundle")
			}
			if info, statErr := os.Lstat(pagePath); err == nil && (statErr != nil || !info.Mode().IsRegular() || filepath.Base(pagePath) != "index.html" && !secondOpinionTreePagePattern.MatchString(filepath.Base(pagePath))) {
				err = errors.New("base page must be index.html or a tree page (index-v1c...html) in the bundle")
			}
			if filepath.Base(pagePath) != "index.html" {
				direction.Base = filepath.Base(pagePath)
			}
		} else {
			bundle, err = secondOpinionBundle(execution.RepositoryRoot, options.bundle)
		}
		if err != nil {
			return usageResult(CommandAIReportSecondOpinion, err.Error())
		}
		templatePath := os.Getenv("DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT")
		template, err = os.ReadFile(templatePath)
		if err != nil {
			return usageResult(CommandAIReportSecondOpinion, "use the suite's do-work-cli.sh launcher, or set DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT to its trusted shipped template: "+err.Error())
		}
		if options.preset != "" {
			// A preset in the report bundle wins over the shipped one of the same name, which sits beside the trusted prompt.
			folders := []string{filepath.Join(bundle, secondOpinionPresetsFolder), filepath.Join(filepath.Dir(templatePath), secondOpinionPresetsFolder)}
			for _, folder := range folders {
				candidate := filepath.Join(folder, options.preset+".md")
				if info, err := os.Lstat(candidate); err == nil {
					if !info.Mode().IsRegular() {
						return usageResult(CommandAIReportSecondOpinion, "preset must be a regular file: "+candidate)
					}
					direction.PresetFile = candidate
					break
				}
			}
			if direction.PresetFile == "" {
				found := []string{}
				for _, folder := range folders {
					matches, _ := filepath.Glob(filepath.Join(folder, "*.md"))
					for _, match := range matches {
						found = append(found, strings.TrimSuffix(filepath.Base(match), ".md"))
					}
				}
				sort.Strings(found)
				if found = slices.Compact(found); len(found) == 0 {
					found = []string{"none"}
				}
				return usageResult(CommandAIReportSecondOpinion, fmt.Sprintf("unknown preset %q; presets found: %s", options.preset, strings.Join(found, ", ")))
			}
			data, err := os.ReadFile(direction.PresetFile)
			direction.PresetText = strings.TrimSpace(string(data))
			if err != nil || direction.PresetText == "" || strings.Contains(direction.PresetText, "{{") {
				return usageResult(CommandAIReportSecondOpinion, "preset must be a readable, non-empty file without {{: "+direction.PresetFile)
			}
			if relative, err := filepath.Rel(execution.RepositoryRoot, direction.PresetFile); err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				direction.PresetFile = filepath.ToSlash(relative)
			}
		}
	}
	ctx, stop, signalStatus := mutationSignalContext()
	defer stop()
	payload := &resultmodel.AIReportSecondOpinionResult{BundlePath: bundle, InvocationMode: "targeted", BackendRuns: []resultmodel.AIReportSecondOpinionRun{}}
	backends := []string{options.backend}
	if options.prepare || options.finalize != "" {
		payload.InvocationMode = "session"
	} else if options.all {
		payload.InvocationMode = "batch"
		backends = []string{"codex", "gemini", "agy"}
	} else if options.backend == "" {
		payload.InvocationMode = "terminal"
		backends = nil
		for _, name := range []string{"codex", "gemini", "agy"} {
			if _, err := exec.LookPath(name); err == nil {
				backends = []string{name}
				break
			}
		}
		if len(backends) == 0 {
			backends = []string{"codex"}
		}
	}
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, AIReportSecondOpinion: payload}
	basePage := ""
	if direction.parentTreeID() != "" {
		relative, _ := filepath.Rel(execution.RepositoryRoot, bundle)
		basePage = filepath.ToSlash(filepath.Join(relative, cmp.Or(direction.Base, "index.html")))
	}
	var summary strings.Builder
	for _, backend := range backends {
		row := resultmodel.AIReportSecondOpinionRun{BackendName: backend, RequestedModel: options.model, RunState: "failed", JudgeVerdict: "not-run", BannerState: "not-inserted", JudgeFindings: []resultmodel.AIReportSecondOpinionFinding{}, OutsidePaths: []string{}, SeenModels: []string{}, Preset: direction.Preset, PresetFile: direction.PresetFile, EditorDirection: direction.EditorDirection, BasePage: basePage}
		switch {
		case options.finalize != "":
			secondOpinionFinalize(execution.RepositoryRoot, bundle, stageName, &row, &result)
		case options.prepare:
			// The session is the backend: its declared model is recorded as such, never probed through a subprocess.
			row.ActualModel = options.model
			row.ModelEvidence = "session-declared"
			secondOpinionPrepare(execution.RepositoryRoot, bundle, string(template), direction, &row)
		case backend == options.exclude:
			row.RunState = "excluded"
			row.FailureReason = "excluded calling backend"
		default:
			if executable, lookupErr := exec.LookPath(backend); lookupErr != nil {
				row.RunState = "skipped"
				row.FailureReason = "backend is not installed"
			} else {
				secondOpinionRun(ctx, execution.RepositoryRoot, bundle, executable, string(template), direction, &row, &result)
			}
		}
		payload.BackendRuns = append(payload.BackendRuns, row)
		fmt.Fprintf(&summary, "%s: %s; model=%s; page=%s; reply=%s; judge=%s; banner=%s", row.BackendName, row.RunState, row.ActualModel, row.PagePath, row.ReplyPath, row.JudgeVerdict, row.BannerState)
		if row.FailureReason != "" {
			fmt.Fprintf(&summary, "; %s", row.FailureReason)
		}
		summary.WriteByte('\n')
		if row.RunState != "ran" && row.RunState != "excluded" && row.RunState != "prepared" || row.RunState == "ran" && row.JudgeVerdict != "pass" {
			severity := resultmodel.SeverityWarning
			if !options.all {
				severity = resultmodel.SeverityError
				result.Outcome = resultmodel.OutcomeFailure
			} else {
				result.Outcome = resultmodel.OutcomeFindings
			}
			result.Findings = append(result.Findings, toolboxFinding(CommandAIReportSecondOpinion, "AI-REPORT-SECOND-OPINION-"+strings.ToUpper(row.RunState), severity, row.OutsidePaths, row.BackendName+": "+row.FailureReason, resultmodel.FixabilityManual, "this backend did not produce a render-verified linked second opinion"))
		}
		if ctx.Err() != nil {
			break
		}
	}
	output := summary.String()
	result.ExactTextOutput = &output
	if status := signalStatus(); status != 0 {
		result.ExitCodeOverride = status
		result.Outcome = resultmodel.OutcomeFailure
	}
	return result
}

func secondOpinionParse(args []string) (secondOpinionOptions, error) {
	options := secondOpinionOptions{bundle: "latest"}
	bundleSeen := false
	for i := 0; i < len(args); i++ {
		var err error
		argument := args[i]
		switch {
		case argument == "--backend" || strings.HasPrefix(argument, "--backend="):
			options.backend, err = optionValue(args, &i, "--backend")
		case argument == "--model" || strings.HasPrefix(argument, "--model="):
			options.model, err = optionValue(args, &i, "--model")
		case argument == "--exclude" || strings.HasPrefix(argument, "--exclude="):
			options.exclude, err = optionValue(args, &i, "--exclude")
		case argument == "--all":
			options.all = true
		case argument == "--prepare":
			options.prepare = true
		case argument == "--finalize" || strings.HasPrefix(argument, "--finalize="):
			options.finalize, err = optionValue(args, &i, "--finalize")
		case argument == "--preset" || strings.HasPrefix(argument, "--preset="):
			options.preset, err = optionValue(args, &i, "--preset")
			if err == nil && !regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`).MatchString(options.preset) {
				err = errors.New("--preset must be a lowercase name of letters, digits and hyphens")
			}
		case argument == "--editor-direction" || strings.HasPrefix(argument, "--editor-direction="):
			// optionValue returns "" for --editor-direction '', and {{ would trip the unresolved-placeholder check.
			options.editorDirection, err = optionValue(args, &i, "--editor-direction")
			if err == nil && strings.TrimSpace(options.editorDirection) == "" {
				err = errors.New("--editor-direction must not be empty")
			} else if err == nil && strings.Contains(options.editorDirection, "{{") {
				err = errors.New("--editor-direction must not contain {{")
			}
		case strings.HasPrefix(argument, "-") || bundleSeen:
			err = fmt.Errorf("unexpected argument %q", argument)
		default:
			options.bundle = argument
			bundleSeen = true
		}
		if err != nil {
			return options, err
		}
	}
	validBackend := func(name string) bool {
		return name == "codex" || name == "gemini" || name == "agy" || name == "claude"
	}
	if options.backend != "" && !validBackend(options.backend) || options.exclude != "" && !validBackend(options.exclude) {
		return options, errors.New("unknown backend")
	}
	if options.all && (options.backend != "" || options.model != "" || options.exclude == "") || !options.all && options.exclude != "" || options.model != "" && options.backend == "" {
		return options, errors.New("--all requires --exclude and cannot name a backend/model; --model requires --backend")
	}
	if options.model != "" && !secondOpinionModelPattern.MatchString(options.model) {
		return options, errors.New("model must be a filename-safe native slug")
	}
	if options.prepare && (options.backend == "" || options.model == "" || options.all || options.exclude != "") {
		return options, errors.New("--prepare requires --backend and --model and cannot batch")
	}
	if options.finalize != "" && (bundleSeen || options.backend != "" || options.model != "" || options.all || options.exclude != "" || options.prepare || options.preset != "" || options.editorDirection != "") {
		return options, errors.New("--finalize takes only the staging page path")
	}
	return options, nil
}

func secondOpinionBundle(repository, selection string) (string, error) {
	catalog, _, err := reportIndexWalk(repository)
	if err != nil {
		return "", err
	}
	if selection == "latest" {
		if len(catalog) == 0 {
			return "", errors.New("no report bundles found")
		}
		selection = catalog[0].Path
	}
	bundle := judgeAbsolutePath(repository, selection)
	relative, err := filepath.Rel(repository, bundle)
	if err != nil || filepath.Dir(relative) != reportsFolder || strings.HasPrefix(filepath.Base(relative), ".") {
		return "", errors.New("bundle must be a physical direct child of ai-reports")
	}
	if err := validateNoLinkedAncestors(repository, filepath.ToSlash(relative)+"/index.html", false); err != nil {
		return "", err
	}
	for _, entry := range catalog {
		if entry.Path == filepath.ToSlash(relative) && entry.SupersededBy != "" { // Follow the same current-revision chain as revise.
			seen := map[string]bool{}
			current := entry
			for current.SupersededBy != "" {
				if seen[current.Path] {
					return "", errors.New("report supersession cycle")
				}
				seen[current.Path] = true
				next := current.SupersededBy
				found := false
				for _, candidate := range catalog {
					if candidate.Path == next {
						current = candidate
						found = true
						break
					}
				}
				if !found {
					return "", errors.New("missing report successor")
				}
			}
			relative = filepath.FromSlash(current.Path)
			bundle = filepath.Join(repository, relative)
			break
		}
	}
	if err := validateNoLinkedAncestors(repository, filepath.ToSlash(relative)+"/index.html", false); err != nil {
		return "", err
	}
	info, err := os.Lstat(filepath.Join(bundle, "index.html"))
	if err != nil || !info.Mode().IsRegular() {
		return "", errors.New("bundle must hold a regular index.html")
	}
	return bundle, nil
}

// secondOpinionFile is one bundle entry. info and data are live when the snapshot was taken in this process; a snapshot rebuilt from a prepare record carries only mode and digest.
type secondOpinionFile struct {
	info   os.FileInfo
	data   []byte
	mode   fs.FileMode
	digest string
}

// secondOpinionRecord is the prepare sidecar: the before-snapshot as digests, never bytes, so finalize can prove preservation in a later process.
type secondOpinionRecord struct {
	Backend string                         `json:"backend"`
	Model   string                         `json:"model"`
	Prefix  string                         `json:"prefix"`
	Files   map[string]secondOpinionDigest `json:"files"`
	// A record written before directions existed has no direction and finalizes as a plain run.
	Direction secondOpinionDirection `json:"direction,omitzero"`
}
type secondOpinionDigest struct {
	Mode   uint32 `json:"mode"`
	SHA256 string `json:"sha256"`
}

func secondOpinionDigests(before map[string]secondOpinionFile) map[string]secondOpinionDigest {
	digests := map[string]secondOpinionDigest{}
	for name, item := range before {
		digests[name] = secondOpinionDigest{Mode: uint32(item.mode), SHA256: item.digest}
	}
	return digests
}

func secondOpinionSnapshot(bundle string) (map[string]secondOpinionFile, error) {
	files := map[string]secondOpinionFile{}
	err := filepath.WalkDir(bundle, func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if name == bundle {
			return nil
		}
		relative, _ := filepath.Rel(bundle, name)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() && !info.IsDir() {
			return fmt.Errorf("bundle contains a link or special file: %s", relative)
		}
		item := secondOpinionFile{info: info, mode: info.Mode()}
		if info.Mode().IsRegular() {
			item.data, err = os.ReadFile(name)
			if err != nil {
				return err
			}
			item.digest = fmt.Sprintf("%x", sha256.Sum256(item.data))
		}
		files[relative] = item
		return nil
	})
	return files, err
}

// secondOpinionPreserved proves every pre-existing bundle path is untouched and every new one is an allowed staging
// name or a prefixed regular root file. A live before-entry is held to mode, bytes and inode identity; an entry rebuilt
// from a prepare record (no info) is held to mode and digest, which is all that survives the process boundary.
func secondOpinionPreserved(bundle string, before map[string]secondOpinionFile, allowed []string, prefix string) error {
	after, err := secondOpinionSnapshot(bundle)
	if err != nil {
		return err
	}
	for name, old := range before {
		current, ok := after[name]
		changed := !ok
		if ok && old.info != nil {
			changed = old.info.Mode() != current.info.Mode() || !bytes.Equal(old.data, current.data) || !os.SameFile(old.info, current.info)
		} else if ok {
			changed = old.mode != current.mode || old.digest != current.digest
		}
		if changed {
			return fmt.Errorf("backend changed existing bundle path: %s", name)
		}
	}
	for name, current := range after {
		if _, ok := before[name]; ok || slices.Contains(allowed, name) {
			continue
		}
		if prefix == "" || filepath.Dir(name) != "." || !strings.HasPrefix(name, prefix) || !current.info.Mode().IsRegular() {
			return fmt.Errorf("backend created an unprefixed or non-file asset: %s", name)
		}
	}
	return nil
}

func secondOpinionOutside(repository, bundle string) ([]string, error) {
	command := exec.Command("git", "-C", repository, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git status failed: %w", err)
	}
	relative, err := filepath.Rel(repository, bundle)
	if err != nil {
		return nil, err
	}
	prefix := filepath.ToSlash(relative) + "/"
	outside := []string{}
	tokens := bytes.Split(output, []byte{0})
	for i := 0; i < len(tokens); i++ {
		record := string(tokens[i])
		if record == "" {
			continue
		}
		if len(record) < 4 {
			return nil, errors.New("malformed git status")
		}
		paths := []string{record[3:]}
		if strings.ContainsAny(record[:2], "RC") {
			i++
			if i >= len(tokens) || len(tokens[i]) == 0 {
				return nil, errors.New("incomplete git rename status")
			}
			paths = append(paths, string(tokens[i]))
		}
		for _, name := range paths {
			if !strings.HasPrefix(name, prefix) {
				outside = append(outside, name)
			}
		}
	}
	sort.Strings(outside)
	return outside, nil
}

// secondOpinionOutsideCheck records the Git status outside the bundle on the row and refuses any dirt there.
func secondOpinionOutsideCheck(repository, bundle string, row *resultmodel.AIReportSecondOpinionRun) error {
	paths, err := secondOpinionOutside(repository, bundle)
	row.OutsidePaths = paths
	if err != nil {
		return err
	}
	if len(paths) > 0 {
		return fmt.Errorf("outside-bundle changes: %s", strings.Join(paths, ", "))
	}
	return nil
}

// secondOpinionPrepare is the session form's first phase: it reserves the staging page, records the bundle's
// digests beside it, and hands the filled prompt back so the editorial pass runs in the agent's own session
// instead of a hidden subprocess of its CLI. Nothing here executes a backend.
func secondOpinionPrepare(repository, bundle, template string, direction secondOpinionDirection, row *resultmodel.AIReportSecondOpinionRun) {
	fail := func(err error) { row.FailureReason = err.Error() }
	if err := secondOpinionOutsideCheck(repository, bundle, row); err != nil {
		fail(err)
		return
	}
	before, err := secondOpinionSnapshot(bundle)
	if err != nil {
		fail(err)
		return
	}
	prefix := row.BackendName + "-" + row.RequestedModel + "-"
	stage, err := os.CreateTemp(bundle, "."+prefix+"staging-*.html")
	if err != nil {
		fail(err)
		return
	}
	stage.Close()
	stagePath := stage.Name()
	stem := strings.TrimSuffix(stagePath, ".html")
	defer func() {
		// A prepare that failed after reserving its files leaves nothing behind: no agent has written to them yet.
		if row.RunState != "prepared" {
			_ = os.Remove(stagePath)
			_ = os.Remove(stem + ".json")
		}
	}()
	record, err := json.Marshal(secondOpinionRecord{Backend: row.BackendName, Model: row.RequestedModel, Prefix: prefix, Files: secondOpinionDigests(before), Direction: direction})
	if err == nil {
		err = os.WriteFile(stem+".json", record, 0600)
	}
	if err != nil {
		fail(err)
		return
	}
	// The agent's cwd is the repository, not the bundle, so the output path is absolute here.
	prompt, err := secondOpinionPrompt(template, repository, bundle, stagePath, prefix, direction)
	if err != nil {
		fail(err)
		return
	}
	row.StagingPage = stagePath
	row.StagingReply = stem + "-reply.md"
	row.AssetPrefix = prefix
	row.PromptText = prompt
	row.RunState = "prepared"
}

// secondOpinionFinalize is the session form's second phase: it re-checks the bundle against the prepare record and
// publishes the page the agent wrote through the shared launcher tail. Only a never-written staging page is removed;
// every other refusal keeps the page and reply stage on disk and names them.
func secondOpinionFinalize(repository, bundle, stageName string, row *resultmodel.AIReportSecondOpinionRun, result *resultmodel.CommandResult) {
	fail := func(err error) { row.FailureReason = err.Error() }
	if !secondOpinionStagePattern.MatchString(stageName) {
		fail(errors.New("not a launcher-reserved staging page"))
		return
	}
	stem := strings.TrimSuffix(stageName, ".html")
	sidecar, replyStage := stem+".json", stem+"-reply.md"
	relative, _ := filepath.Rel(repository, bundle)
	defer func() {
		if info, err := os.Lstat(filepath.Join(bundle, stageName)); err == nil && info.Size() > 0 {
			row.FailureReason += "; staged page kept at " + filepath.ToSlash(filepath.Join(relative, stageName))
		}
	}()
	data, err := os.ReadFile(filepath.Join(bundle, sidecar))
	if err != nil {
		fail(errors.New("no prepare record for this staging page"))
		return
	}
	var record secondOpinionRecord
	if err := json.Unmarshal(data, &record); err != nil {
		fail(fmt.Errorf("prepare record unreadable: %w", err))
		return
	}
	row.BackendName, row.RequestedModel, row.ActualModel, row.ModelEvidence = record.Backend, record.Model, record.Model, "session-declared"
	row.Preset, row.PresetFile, row.EditorDirection = record.Direction.Preset, record.Direction.PresetFile, record.Direction.EditorDirection
	if record.Direction.parentTreeID() != "" {
		row.BasePage = filepath.ToSlash(filepath.Join(relative, cmp.Or(record.Direction.Base, "index.html")))
	}
	if err := secondOpinionOutsideCheck(repository, bundle, row); err != nil {
		fail(err)
		return
	}
	if info, err := os.Lstat(filepath.Join(bundle, stageName)); err != nil || info.Size() == 0 {
		// Nothing was written: the reserved names go back, and nothing is published.
		if bundleRoot, err := os.OpenRoot(bundle); err == nil {
			for _, leaf := range []string{stageName, sidecar, replyStage} {
				_ = bundleRoot.Remove(leaf)
			}
			bundleRoot.Close()
		}
		fail(errors.New("staging page was not written"))
		return
	}
	recorded := map[string]secondOpinionFile{}
	for name, digest := range record.Files {
		recorded[name] = secondOpinionFile{mode: fs.FileMode(digest.Mode), digest: digest.SHA256}
	}
	if err := secondOpinionPreserved(bundle, recorded, []string{stageName, sidecar, replyStage}, record.Prefix); err != nil {
		fail(err)
		return
	}
	reply, err := os.ReadFile(filepath.Join(bundle, replyStage))
	if err != nil || strings.TrimSpace(string(reply)) == "" {
		fail(errors.New("no final editorial reply"))
		return
	}
	after, err := secondOpinionSnapshot(bundle)
	if err != nil {
		fail(err)
		return
	}
	secondOpinionPublication(repository, bundle, stageName, after[stageName].data, string(reply), after["index.html"], []string{stageName, sidecar, replyStage}, "", record.Direction, row, result)
}

func secondOpinionRun(ctx context.Context, repository, bundle, executable, template string, direction secondOpinionDirection, row *resultmodel.AIReportSecondOpinionRun, result *resultmodel.CommandResult) {
	fail := func(err error) { row.FailureReason = err.Error() }
	checkOutside := func() error { return secondOpinionOutsideCheck(repository, bundle, row) }
	if err := checkOutside(); err != nil {
		fail(err)
		return
	}
	before, err := secondOpinionSnapshot(bundle)
	if err != nil {
		fail(err)
		return
	}
	// The tiny read-only probe resolves defaults before the asset prefix is sent.
	probeContext, stopProbe := context.WithTimeout(ctx, 90*time.Second)
	defer stopProbe()
	model := row.RequestedModel
	if row.BackendName == "agy" && model == "" {
		capture := secondOpinionCapture(probeContext, secondOpinionProgress{phase: "probe", backend: row.BackendName}, bundle, executable, "models")
		if err := checkOutside(); err != nil {
			fail(err)
			return
		}
		if err := secondOpinionPreserved(bundle, before, nil, ""); err != nil {
			fail(err)
			return
		}
		if capture.err != nil {
			fail(capture.err)
			return
		}
		model, err = secondOpinionFlashModel(capture.stdout)
		if err != nil {
			fail(err)
			return
		}
	}
	probePrompt := "Reply OK only. Do not use tools or modify any files."
	argv := secondOpinionArgv(row.BackendName, model, probePrompt, true)
	capture := secondOpinionCapture(probeContext, secondOpinionProgress{phase: "probe", backend: row.BackendName}, bundle, executable, argv...)
	if err := checkOutside(); err != nil {
		fail(err)
		return
	}
	if err := secondOpinionPreserved(bundle, before, nil, ""); err != nil {
		fail(err)
		return
	}
	resolved, _, _, err := secondOpinionMetadata(row.BackendName, model, capture)
	if err != nil {
		fail(err)
		return
	}
	row.ActualModel = resolved
	if row.RequestedModel != "" && resolved != row.RequestedModel {
		fail(fmt.Errorf("requested model %s resolved to %s; substitution refused", row.RequestedModel, resolved))
		return
	}
	if model != "" && resolved != model {
		fail(fmt.Errorf("selected model %s resolved to %s; substitution refused", model, resolved))
		return
	}
	row.ActualModel = resolved
	stopProbe()
	prefix := row.BackendName + "-" + resolved + "-"
	bundleRoot, err := os.OpenRoot(bundle)
	if err != nil {
		fail(err)
		return
	}
	defer bundleRoot.Close()
	stage, err := os.CreateTemp(bundle, "."+prefix+"staging-*.html")
	if err != nil {
		fail(err)
		return
	}
	stageName := filepath.Base(stage.Name())
	stage.Close()
	relative, _ := filepath.Rel(repository, bundle)
	// Publication removes the staging leaf itself; what is left here is a page the backend never wrote, which goes, or a finished page a refusal kept, which is named.
	defer func() {
		info, statErr := bundleRoot.Stat(stageName)
		if statErr != nil {
			return
		}
		if info.Size() == 0 {
			_ = bundleRoot.Remove(stageName)
			return
		}
		row.FailureReason += "; staged page kept at " + filepath.ToSlash(filepath.Join(relative, stageName))
	}()
	// The backend's cwd is the bundle, so the output file is the bare stage name here.
	prompt, err := secondOpinionPrompt(template, repository, bundle, stageName, prefix, direction)
	if err != nil {
		fail(err)
		return
	}
	capture = secondOpinionCapture(ctx, secondOpinionProgress{phase: "generation", backend: row.BackendName, pagePath: filepath.Join(bundle, stageName)}, bundle, executable, secondOpinionArgv(row.BackendName, resolved, prompt, false)...)
	if err := checkOutside(); err != nil {
		fail(err)
		return
	}
	if err := secondOpinionPreserved(bundle, before, []string{stageName}, prefix); err != nil {
		fail(err)
		return
	}
	actual, reply, seen, err := secondOpinionMetadata(row.BackendName, resolved, capture)
	row.SeenModels = append([]string{}, seen...)
	if err != nil && !errors.Is(err, errSecondOpinionModelUnverified) {
		fail(err)
		return
	}
	unverified := err != nil || actual != resolved
	row.ActualModel = actual
	row.ModelEvidence = "native"
	if !unverified && strings.TrimSpace(reply) == "" {
		fail(errors.New("backend returned no final editorial reply"))
		return
	}
	page, err := os.ReadFile(filepath.Join(bundle, stageName))
	if err != nil {
		fail(err)
		return
	}
	if unverified && len(page) == 0 {
		// Nothing finished to keep: an empty staging page is a plain failure, and the deferred cleanup removes it.
		fail(fmt.Errorf("backend wrote no page; actual model unverified: requested %s; seen %s", resolved, strings.Join(row.SeenModels, ", ")))
		return
	}
	unverifiedModel := ""
	if unverified {
		unverifiedModel = resolved
	}
	secondOpinionPublication(repository, bundle, stageName, page, reply, before["index.html"], []string{stageName}, unverifiedModel, direction, row, result)
}

// secondOpinionPublication is the launcher tail shared by the exec path and the session finalize: backlink, sibling
// names, direction and tree marks, page and reply publication, removal of the staging leaves, the render judge, the
// change-list check on a tree page, and the pass-gated banner (plain run) or iteration list (tree run).
// unverifiedModel is the requested model whose identity the native evidence did not confirm (empty when verified):
// the finished page is then kept under an -unverified name with no judge and no banner, never deleted.
func secondOpinionPublication(repository, bundle, stageName string, page []byte, reply string, original secondOpinionFile, stageLeaves []string, unverifiedModel string, direction secondOpinionDirection, row *resultmodel.AIReportSecondOpinionRun, result *resultmodel.CommandResult) {
	fail := func(err error) { row.FailureReason = err.Error() }
	relative, _ := filepath.Rel(repository, bundle)
	// Strip launcher marks an editor copied from another page first: a header block ahead of the backlink fails its check.
	page = secondOpinionDecorate(page, secondOpinionDirection{}, secondOpinionTreeMarks{})
	publishedModel := row.ActualModel
	if unverifiedModel != "" {
		publishedModel = unverifiedModel + "-unverified"
		row.ActualModel = ""
		row.FailureReason = fmt.Sprintf("actual model unverified: requested %s; seen %s", unverifiedModel, strings.Join(row.SeenModels, ", "))
	} else if err := secondOpinionBacklink(page); err != nil {
		fail(err)
		return
	}
	// The editor's slug is read only on a directed run; a plain run keeps the plain names.
	variation := direction.filenameSlot()
	if variation != "" {
		if slug := secondOpinionVariationSlug(reply); slug != "" {
			variation += "-" + slug
		}
	}
	pageName, replyName, treeID, err := secondOpinionNames(bundle, row.BackendName, publishedModel, variation, direction.parentTreeID())
	if err != nil {
		fail(err)
		return
	}
	// The judge and the reader see the same decorated bytes; the header needs the id, so decoration follows naming.
	page = secondOpinionDecorate(page, direction, secondOpinionTreeMarks{treeID: treeID, modelLabel: row.BackendName + " " + publishedModel, publishDate: time.Now().UTC().Format(time.DateOnly)})
	publish := func(name string, data []byte) error {
		path := filepath.ToSlash(filepath.Join(relative, name))
		if err := rootedPublishFile(repository, path, data, 0644, false); err != nil {
			return err
		}
		result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: path, Kind: "created", Detail: "second opinion"})
		return nil
	}
	if err := publish(pageName, page); err != nil {
		fail(err)
		return
	}
	row.PagePath = filepath.ToSlash(filepath.Join(relative, pageName))
	row.TreeID = treeID
	// The editor knows the reserved stage name; the saved reply names its published page.
	reply = strings.NewReplacer("file://"+filepath.ToSlash(filepath.Join(bundle, stageName)), pageName, filepath.Join(bundle, stageName), pageName, filepath.ToSlash(filepath.Join(bundle, stageName)), pageName, stageName, pageName).Replace(reply)
	if err := publish(replyName, []byte(reply)); err != nil {
		fail(err)
		return
	}
	row.ReplyPath = filepath.ToSlash(filepath.Join(relative, replyName))
	// Both sibling files exist now, so the staging leaves have been published and can go.
	if bundleRoot, err := os.OpenRoot(bundle); err == nil {
		for _, leaf := range stageLeaves {
			_ = bundleRoot.Remove(leaf)
		}
		bundleRoot.Close()
	}
	if unverifiedModel != "" {
		row.RunState = "unverified"
		return
	}
	stopHeartbeat := secondOpinionHeartbeat(secondOpinionProgress{phase: "judge", backend: row.BackendName, pagePath: filepath.Join(bundle, pageName)})
	verdict, findings := secondOpinionJudge(bundle, page)
	stopHeartbeat()
	if treeID != "" {
		// A change-list link with no target is a judge finding: the page stays published, index.html stays untouched.
		if changeFindings := secondOpinionChangeList(page); len(changeFindings) > 0 {
			findings = append(findings, changeFindings...)
			if verdict == "pass" || verdict == "skipped" {
				verdict = "fail"
			}
		}
	}
	row.JudgeVerdict = verdict
	row.JudgeFindings = findings
	row.RunState = "ran"
	if verdict != "pass" {
		row.FailureReason = "render judge " + verdict + "; original left unchanged"
		return
	}
	if err := secondOpinionOutsideCheck(repository, bundle, row); err != nil {
		row.RunState = "failed"
		fail(err)
		return
	}
	current, err := os.Lstat(filepath.Join(bundle, "index.html"))
	if err != nil || !os.SameFile(original.info, current) || current.Mode() != original.info.Mode() {
		row.RunState = "failed"
		fail(errors.New("original identity changed before banner publication"))
		return
	}
	originalNow, err := os.ReadFile(filepath.Join(bundle, "index.html"))
	if err != nil || !bytes.Equal(originalNow, original.data) {
		row.RunState = "failed"
		fail(errors.New("original bytes changed before banner publication"))
		return
	}
	var updated []byte
	var inserted bool
	detail := "second-opinion banner only"
	if treeID != "" {
		// A tree page is linked from the one iteration list, never from a per-filename banner.
		updated, err = secondOpinionTreeList(bundle, original.data)
		inserted, detail = !bytes.Equal(updated, original.data), "second-opinion iteration list only"
	} else {
		updated, inserted, err = secondOpinionBanner(original.data, pageName)
	}
	if err != nil {
		row.RunState = "failed"
		fail(err)
		return
	}
	if !inserted {
		row.BannerState = "already-present"
		return
	}
	originalPath := filepath.ToSlash(filepath.Join(relative, "index.html"))
	if err := rootedPublishFile(repository, originalPath, updated, original.info.Mode(), true); err != nil {
		row.RunState = "failed"
		fail(err)
		return
	}
	row.BannerState = "inserted"
	result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: originalPath, Kind: "modified", Detail: detail})
}

// secondOpinionPrompt fills the shipped template for one run. A directed run refuses a template without
// {{EDITOR_DIRECTION}}, which would drop the direction silently, and a tree run one without {{BASE_PAGE}} or
// {{CHANGE_LIST_REQUIREMENT}}; a plain run with such a template still works.
func secondOpinionPrompt(template, repository, bundle, outputFile, prefix string, direction secondOpinionDirection) (string, error) {
	text := direction.combinedText()
	if text != "" && !strings.Contains(template, "{{EDITOR_DIRECTION}}") {
		return "", errors.New("prompt template has no {{EDITOR_DIRECTION}} placeholder for the direction")
	}
	if text == "" {
		text = "None."
	}
	basePage, changeList := "index.html", "None for this pass."
	if direction.parentTreeID() != "" {
		if !strings.Contains(template, "{{BASE_PAGE}}") || !strings.Contains(template, "{{CHANGE_LIST_REQUIREMENT}}") {
			return "", errors.New("prompt template has no {{BASE_PAGE}} or {{CHANGE_LIST_REQUIREMENT}} placeholder for the iteration")
		}
		basePage, changeList = cmp.Or(direction.Base, "index.html"), secondOpinionChangeListRequirement
	}
	prompt := strings.NewReplacer("{{REPOSITORY_ROOT}}", repository, "{{BUNDLE_PATH}}", bundle, "{{OUTPUT_FILE}}", outputFile, "{{ASSET_PREFIX}}", prefix, "{{BACKLINK_REQUIREMENT}}", secondOpinionBacklinkRequirement, "{{EDITOR_DIRECTION}}", text, "{{BASE_PAGE}}", basePage, "{{CHANGE_LIST_REQUIREMENT}}", changeList).Replace(template)
	if strings.Contains(prompt, "{{") {
		return "", errors.New("unresolved prompt placeholder")
	}
	return prompt, nil
}

// secondOpinionVariationSlug cleans the editor's first "Variation slug:" reply line into at most 4 kebab words and
// 40 characters. It is editor text, so it only ever becomes [a-z0-9-]; empty means no slug segment.
func secondOpinionVariationSlug(reply string) string {
	match := regexp.MustCompile(`(?im)^[ \t]*[-*]?[ \t]*variation slug:(.*)$`).FindStringSubmatch(reply)
	if match == nil {
		return ""
	}
	words := strings.FieldsFunc(strings.ToLower(match[1]), func(character rune) bool {
		return (character < 'a' || character > 'z') && (character < '0' || character > '9')
	})
	if len(words) > 4 {
		words = words[:4]
	}
	slug := strings.Join(words, "-")
	if len(slug) > 40 {
		slug = strings.TrimRight(slug[:40], "-")
	}
	return slug
}

func secondOpinionArgv(backend, model, prompt string, probe bool) []string {
	switch backend {
	case "codex":
		sandbox := "workspace-write"
		if probe {
			sandbox = "read-only"
		}
		args := []string{"exec", "-s", sandbox, "--ephemeral"}
		if model != "" {
			args = append(args, "-m", model)
		}
		return append(args, prompt)
	case "gemini":
		approval := "auto_edit"
		if probe {
			approval = "plan"
		}
		args := []string{"-p", prompt, "--output-format", "json", "--approval-mode", approval}
		if model != "" {
			args = append(args, "-m", model)
		}
		return args
	case "agy":
		if probe {
			return []string{"--model", model, "--print-timeout", "30s", "--output-format", "json", "-p", "/model"}
		}
		return []string{"--model", model, "--print-timeout", "15m", "--output-format", "stream-json", "--dangerously-skip-permissions", "-p", prompt}
	case "claude":
		args := []string{"-p", prompt, "--output-format", "json"}
		if model != "" {
			args = append(args, "--model", model)
		}
		if probe {
			args = append(args, "--tools", "")
		} else {
			args = append(args, "--permission-mode", "acceptEdits")
		}
		return args
	}
	return nil
}

type secondOpinionCaptured struct {
	stdout, stderr string
	err            error
}

// secondOpinionBuffer captures up to 16 MB of one backend stream and mirrors every line to the progress sink as it
// arrives, past the cap too: pending holds at most one unfinished line of up to 64 KiB, so the mirror cannot grow memory.
// captured is a named field, never embedded: an embedded bytes.Buffer promotes ReadFrom, io.Copy prefers it, and
// neither the cap nor the mirror in Write would ever run.
type secondOpinionBuffer struct {
	captured     bytes.Buffer
	truncated    bool
	mirrorPrefix string
	pending      []byte
}

func (buffer *secondOpinionBuffer) Write(data []byte) (int, error) {
	length := len(data)
	buffer.pending = append(buffer.pending, data...)
	for {
		index := bytes.IndexByte(buffer.pending, '\n')
		if index < 0 && len(buffer.pending) < 64*1024 {
			break
		}
		if index < 0 {
			index = len(buffer.pending) - 1
		}
		buffer.mirrorLine(buffer.pending[:index+1])
		buffer.pending = buffer.pending[:copy(buffer.pending, buffer.pending[index+1:])]
	}
	remaining := 16*1024*1024 - buffer.captured.Len()
	if len(data) > remaining {
		data = data[:remaining]
		buffer.truncated = true
	}
	_, _ = buffer.captured.Write(data)
	return length, nil
}

func (buffer *secondOpinionBuffer) mirrorLine(line []byte) {
	text := buffer.mirrorPrefix + string(line)
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	secondOpinionProgressLine(text)
}

// flushMirror prints a last line the backend left without a newline; call it only after Wait joined the copiers.
func (buffer *secondOpinionBuffer) flushMirror() {
	if len(buffer.pending) > 0 {
		buffer.mirrorLine(buffer.pending)
		buffer.pending = nil
	}
}

// The supervisor owns the isolated group until cleanup, even if the backend
// exits before its children. Reaping the backend as the group leader would lose
// Getpgid ownership proof (ESRCH does not mean the descendants have stopped).
// Keep arguments positional; neither prompts nor executable paths are shell code.
const secondOpinionSupervisor = `"$@" </dev/null 3>&- &
backend_pid=$!
wait "$backend_pid"
backend_status=$?
printf '%s\n' "$backend_status" >&3
exec 3>&-
IFS= read -r release
exit "$backend_status"
`

func secondOpinionCapture(ctx context.Context, progress secondOpinionProgress, directory, executable string, args ...string) secondOpinionCaptured {
	bounded, cancel := context.WithTimeout(ctx, 16*time.Minute)
	defer cancel()
	if err := bounded.Err(); err != nil {
		return secondOpinionCaptured{err: err}
	}
	command := exec.Command("/bin/sh", append([]string{"-c", secondOpinionSupervisor, "second-opinion-supervisor", executable}, args...)...)
	command.Dir = directory
	if err := configureOwnedProcess(command); err != nil {
		return secondOpinionCaptured{err: err}
	}
	statusRead, statusWrite, err := os.Pipe()
	if err != nil {
		return secondOpinionCaptured{err: err}
	}
	defer statusRead.Close()
	defer statusWrite.Close()
	holdRead, holdWrite, err := os.Pipe()
	if err != nil {
		return secondOpinionCaptured{err: err}
	}
	defer holdRead.Close()
	defer holdWrite.Close()
	// Only the supervisor sees these channels. The backend gets closed stdin
	// and closes fd 3 before exec, so its output cannot forge completion status.
	command.Stdin = holdRead
	command.ExtraFiles = []*os.File{statusWrite}
	// Non-file writers make exec copy both pipes through goroutines that Wait joins, so each buffer has one writer.
	stdout := secondOpinionBuffer{mirrorPrefix: "[" + progress.backend + " stdout] "}
	stderr := secondOpinionBuffer{mirrorPrefix: "[" + progress.backend + " stderr] "}
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		return secondOpinionCaptured{err: err}
	}
	stopHeartbeat := secondOpinionHeartbeat(progress)
	defer stopHeartbeat()
	_ = statusWrite.Close()
	_ = holdRead.Close()
	done := make(chan error, 1)
	go func() {
		data, readErr := io.ReadAll(io.LimitReader(statusRead, 32))
		if readErr != nil {
			done <- fmt.Errorf("backend completion channel: %w", readErr)
			return
		}
		text := string(data)
		status, parseErr := strconv.Atoi(strings.TrimSuffix(text, "\n"))
		if parseErr != nil || status < 0 || status > 255 || text != strconv.Itoa(status)+"\n" {
			done <- errors.New("invalid backend completion status")
			return
		}
		if status != 0 {
			done <- fmt.Errorf("exit status %d", status)
			return
		}
		done <- nil
	}()
	select {
	case err = <-done:
	case <-bounded.Done():
		err = bounded.Err()
	}
	// Do not Wait before this proof-bound cleanup: the supervisor must retain
	// its PID/group identity while descendants are terminated bottom-up.
	cleanupErr := ownedprocess.TerminateGroup(command.Process.Pid, reportImageGracePeriod)
	_ = statusRead.Close()
	_ = holdWrite.Close()
	waitErr := command.Wait()
	stdout.flushMirror()
	stderr.flushMirror()
	if cleanupErr != nil && !errors.Is(cleanupErr, os.ErrProcessDone) {
		err = errors.Join(err, cleanupErr)
	}
	// The supervisor is intentionally signalled; its exit status is not the
	// backend's. Still retain a pipe-copy or other non-exit Wait failure.
	var exitErr *exec.ExitError
	if waitErr != nil && !errors.As(waitErr, &exitErr) {
		err = errors.Join(err, waitErr)
	}
	if stdout.truncated || stderr.truncated {
		err = errors.New("backend output exceeded capture limit")
	}
	return secondOpinionCaptured{stdout: stdout.captured.String(), stderr: stderr.captured.String(), err: err}
}

func secondOpinionMetadata(backend, expected string, capture secondOpinionCaptured) (model, reply string, seen []string, err error) {
	if capture.err != nil {
		return "", "", nil, fmt.Errorf("backend failed: %w; %s", capture.err, secondOpinionErrorDetail(backend, capture))
	}
	switch backend {
	case "codex":
		// Only the native startup envelope, before the prompt/tools/answer, is authority.
		lines := strings.Split(capture.stderr, "\n")
		started := false
		header := false
		fields := map[string]string{}
		for _, line := range lines {
			if !started {
				if strings.HasPrefix(line, "OpenAI Codex ") {
					started = true
				}
				continue
			}
			if line == "--------" {
				if header {
					break
				}
				header = true
				continue
			}
			if header {
				key, value, ok := strings.Cut(line, ": ")
				if ok {
					fields[key] = value
				}
			}
		}
		if fields["workdir"] == "" || fields["provider"] == "" || fields["session id"] == "" {
			return "", "", nil, errors.New("Codex native startup model metadata unavailable")
		}
		model = fields["model"]
		reply = capture.stdout
	case "gemini", "claude":
		var envelope struct {
			Response string          `json:"response"`
			Result   string          `json:"result"`
			IsError  bool            `json:"is_error"`
			Type     string          `json:"type"`
			Subtype  string          `json:"subtype"`
			Error    json.RawMessage `json:"error"`
			Stats    struct {
				Models map[string]json.RawMessage `json:"models"`
			} `json:"stats"`
			ModelUsage map[string]json.RawMessage `json:"modelUsage"`
		}
		if err := json.Unmarshal([]byte(capture.stdout), &envelope); err != nil {
			return "", "", nil, fmt.Errorf("invalid %s native JSON: %w", backend, err)
		}
		if envelope.IsError || len(envelope.Error) > 0 && string(envelope.Error) != "null" {
			return "", "", nil, errors.New("backend native result reports an error")
		}
		models := envelope.Stats.Models
		reply = envelope.Response
		if backend == "claude" {
			if envelope.Type != "result" || envelope.Subtype != "success" {
				return "", "", nil, errors.New("Claude terminal success result missing")
			}
			models = envelope.ModelUsage
			reply = envelope.Result
		}
		for name := range models {
			seen = append(seen, name)
		}
		sort.Strings(seen)
		// The requested model wins when the evidence lists it: a subagent's model beside it is recorded, not ambiguity.
		if _, ok := models[expected]; expected != "" && ok {
			model = expected
		} else if len(seen) == 1 {
			model = seen[0]
		} else {
			// The reply and the seen set travel with this error on purpose: the caller publishes a finished page unverified.
			return "", reply, seen, fmt.Errorf("%s %w (%d models)", backend, errSecondOpinionModelUnverified, len(models))
		}
	case "agy":
		stderrLower := strings.ToLower(capture.stderr)
		if strings.Contains(stderrLower, "print timeout") || strings.Contains(stderrLower, "timed out") || strings.Contains(stderrLower, "truncated") {
			return "", "", nil, errors.New("AGY reported incomplete output on stderr")
		}
		// /model is a zero-turn native command, not model-authored text.
		var selection struct {
			Status  string `json:"status"`
			Command struct {
				Name string `json:"name"`
				Data struct {
					ID string `json:"id"`
				} `json:"data"`
			} `json:"command"`
		}
		if json.Unmarshal([]byte(capture.stdout), &selection) == nil && selection.Command.Name == "model" && selection.Status == "SUCCESS" {
			model = selection.Command.Data.ID
			break
		}
		seenInit, seenResult := false, false
		conversation := ""
		for _, line := range strings.Split(capture.stdout, "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var event struct {
				Event          string `json:"event"`
				ConversationID string `json:"conversation_id"`
				Init           struct {
					Model          string `json:"model"`
					CWD            string `json:"cwd"`
					PermissionMode string `json:"permission_mode"`
				} `json:"init"`
				Result struct {
					ConversationID string            `json:"conversation_id"`
					Status         string            `json:"status"`
					Response       string            `json:"response"`
					Error          string            `json:"error"`
					DeniedActions  []json.RawMessage `json:"denied_actions"`
				} `json:"result"`
			}
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				return "", "", nil, errors.New("invalid AGY event stream")
			}
			switch event.Event {
			case "init":
				if seenInit || event.Init.CWD == "" || event.Init.PermissionMode == "" {
					return "", "", nil, errors.New("invalid AGY init metadata")
				}
				seenInit = true
				model = event.Init.Model
				conversation = event.ConversationID
			case "result":
				if seenResult || !seenInit || event.Result.Status != "SUCCESS" || event.Result.ConversationID != conversation || len(event.Result.DeniedActions) > 0 {
					return "", "", nil, fmt.Errorf("AGY terminal result incomplete or unsuccessful: %s %s", event.Result.Status, event.Result.Error)
				}
				seenResult = true
				reply = event.Result.Response
			}
		}
		if !seenInit || !seenResult {
			return "", "", nil, errors.New("AGY requires native init and terminal SUCCESS")
		}
	}
	if seen == nil && model != "" {
		seen = []string{model}
	}
	if !secondOpinionModelPattern.MatchString(model) {
		return "", "", nil, errors.New("actual model is absent or is not a filename-safe native slug")
	}
	return model, reply, seen, nil
}

func secondOpinionErrorDetail(backend string, capture secondOpinionCaptured) string {
	// Native error fields are useful; avoid dumping tool transcripts or account diagnostics.
	var envelope struct {
		Error  json.RawMessage `json:"error"`
		Result struct {
			Error string `json:"error"`
		} `json:"result"`
	}
	for _, line := range strings.Split(capture.stdout, "\n") {
		if json.Unmarshal([]byte(line), &envelope) == nil {
			if envelope.Result.Error != "" {
				return envelope.Result.Error
			}
			if len(envelope.Error) > 0 {
				return string(envelope.Error)
			}
		}
	}
	return backend + " returned no successful native result; inspect backend authentication/availability"
}

func secondOpinionFlashModel(catalog string) (string, error) {
	pattern := regexp.MustCompile(`(?m)^\s*(gemini-([0-9]+(?:\.[0-9]+)*)-flash-(low|medium))\s`)
	matches := pattern.FindAllStringSubmatch(catalog, -1)
	if len(matches) == 0 {
		return "", errors.New("AGY catalog has no Gemini Flash Low/Medium model")
	}
	newer := func(left, right string) int {
		a, b := strings.Split(left, "."), strings.Split(right, ".")
		for i := 0; i < len(a) || i < len(b); i++ {
			x, y := 0, 0
			if i < len(a) {
				x, _ = strconv.Atoi(a[i])
			}
			if i < len(b) {
				y, _ = strconv.Atoi(b[i])
			}
			if x > y {
				return 1
			}
			if x < y {
				return -1
			}
		}
		return 0
	}
	best := matches[0]
	for _, match := range matches[1:] {
		comparison := newer(match[2], best[2])
		if comparison > 0 || comparison == 0 && match[3] == "low" {
			best = match
		}
	}
	// A newer Flash release that lacks Low/Medium must not silently select an older one.
	all := regexp.MustCompile(`(?m)^\s*gemini-([0-9]+(?:\.[0-9]+)*)-flash-[A-Za-z]+\s`).FindAllStringSubmatch(catalog, -1)
	for _, match := range all {
		if newer(match[1], best[2]) > 0 {
			return "", errors.New("newest AGY Flash release has no Low/Medium variant")
		}
	}
	return best[1], nil
}

// secondOpinionNames is the one builder of published second-opinion names, where the stem is
// <backend>-<model>[-<variation>] and variation is empty on a plain run. A plain run (empty parentID) gets
// index-<stem>[-N].html and <stem>[-N]-reply.md. A tree run gets index-<id>-<stem>.html and <id>-<stem>-reply.md, where
// id is parentID + cNN and NN is one above the highest number any tree page name uses directly under the parent, so a
// deleted page's number is never reused while its children remain. It also returns the id (empty on a plain run).
func secondOpinionNames(bundle, backend, model, variation, parentID string) (string, string, string, error) {
	stem := backend + "-" + model
	if variation != "" {
		stem += "-" + variation
	}
	number := 1
	if parentID != "" {
		entries, err := os.ReadDir(bundle)
		if err != nil {
			return "", "", "", err
		}
		for _, entry := range entries {
			if match := secondOpinionTreePagePattern.FindStringSubmatch(entry.Name()); match != nil && strings.HasPrefix(match[1], parentID+"c") {
				used, _ := strconv.Atoi(match[1][len(parentID)+1 : len(parentID)+3])
				number = max(number, used+1)
			}
		}
	}
	for ; ; number++ {
		treeID, suffix := "", ""
		if number > 1 {
			suffix = "-" + strconv.Itoa(number)
		}
		page, reply := "index-"+stem+suffix+".html", stem+suffix+"-reply.md"
		if parentID != "" {
			if number > 99 {
				return "", "", "", fmt.Errorf("no free two-digit child id under %s", parentID)
			}
			treeID = fmt.Sprintf("%sc%02d", parentID, number)
			page, reply = "index-"+treeID+"-"+stem+".html", treeID+"-"+stem+"-reply.md"
		}
		free := true
		for _, name := range []string{page, reply} {
			_, err := os.Lstat(filepath.Join(bundle, name))
			if err == nil {
				free = false
				break
			}
			if !os.IsNotExist(err) {
				return "", "", "", err
			}
		}
		if free {
			return page, reply, treeID, nil
		}
	}
}

// Mask comments and raw-text elements without moving offsets. This prevents a
// body or marker written in a script/example from authorizing an insertion.
var secondOpinionRawHTML = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script\s*>|<style\b[^>]*>.*?</style\s*>|<title\b[^>]*>.*?</title\s*>|<textarea\b[^>]*>.*?</textarea\s*>|<xmp\b[^>]*>.*?</xmp\s*>|<iframe\b[^>]*>.*?</iframe\s*>|<noembed\b[^>]*>.*?</noembed\s*>|<noframes\b[^>]*>.*?</noframes\s*>|<plaintext\b[^>]*>.*$`)
var secondOpinionHiddenHTML = regexp.MustCompile(`(?is)<!--.*?-->|` + secondOpinionRawHTML.String())
var secondOpinionBodyTag = regexp.MustCompile(`(?is)<body\b(?:[^'"<>]|"[^"]*"|'[^']*')*>`)

func secondOpinionVisibleHTML(data []byte) []byte {
	visible := append([]byte(nil), data...)
	for _, span := range secondOpinionHiddenHTML.FindAllIndex(data, -1) {
		for i := span[0]; i < span[1]; i++ {
			visible[i] = ' '
		}
	}
	return visible
}
func secondOpinionBacklink(page []byte) error {
	visible := secondOpinionVisibleHTML(page)
	body := secondOpinionBodyTag.FindIndex(visible)
	if body == nil {
		return errors.New("new page has no body")
	}
	opening := strings.TrimSpace(string(visible[body[1]:]))
	// Harmless inline formatting and attribute order do not change the backlink.
	attributes := `(?:[^'"<>]|"[^"]*"|'[^']*')*`
	inline := regexp.MustCompile(`(?is)^</?(?:span|strong|b|em|i|small|code)\b`)
	tags := regexp.MustCompile(`(?is)</?[a-z][a-z0-9]*\b` + attributes + `>`)
	opening = tags.ReplaceAllStringFunc(opening, func(tag string) string {
		if inline.MatchString(tag) {
			return ""
		}
		return tag
	})
	link := regexp.MustCompile(`(?is)^(?:<(?:p|div|aside|nav|header)\b` + attributes + `>\s*)*Original report:\s*<a\b(` + attributes + `)>(.*?)</a\s*>`)
	match := link.FindStringSubmatch(opening)
	invalid := errors.New("new page must open its body with Original report: index.html linking to index.html")
	if match == nil || strings.TrimSpace(html.UnescapeString(match[2])) != "index.html" {
		return invalid
	}
	// Parse whole attributes so an href-looking string inside title/data is not a link.
	attribute := regexp.MustCompile(`^\s+([^\s=/>]+)(?:\s*=\s*("[^"]*"|'[^']*'|[^\s"'<>=` + "`" + `]+))?`)
	remainder := match[1]
	hrefCount := 0
	for strings.TrimSpace(remainder) != "" {
		parts := attribute.FindStringSubmatch(remainder)
		if parts == nil {
			return invalid
		}
		remainder = remainder[len(parts[0]):]
		if strings.EqualFold(parts[1], "href") {
			hrefCount++
			value := strings.Trim(parts[2], "\"'")
			if html.UnescapeString(value) != "index.html" {
				return invalid
			}
		}
	}
	if hrefCount != 1 {
		return invalid
	}
	return nil
}
func secondOpinionBanner(original []byte, pageName string) ([]byte, bool, error) {
	if !regexp.MustCompile(`^index-[A-Za-z0-9][A-Za-z0-9._-]*\.html$`).MatchString(pageName) {
		return nil, false, errors.New("unsafe second-opinion page name")
	}
	visible := secondOpinionVisibleHTML(original)
	body := secondOpinionBodyTag.FindIndex(visible)
	if body == nil {
		return nil, false, errors.New("original has no opening body tag")
	}
	marker := "<!-- second-opinion-banner:" + pageName + " -->"
	// Only a real HTML comment following the body counts, never a JS string.
	rawText := secondOpinionRawHTML
	commentSurface := append([]byte(nil), original...)
	for _, span := range rawText.FindAllIndex(original, -1) {
		for i := span[0]; i < span[1]; i++ {
			commentSurface[i] = ' '
		}
	}
	comments := regexp.MustCompile(`(?s)<!--.*?-->`).FindAllIndex(commentSurface, -1)
	for _, span := range comments {
		if span[0] >= body[1] && string(commentSurface[span[0]:span[1]]) == marker {
			return original, false, nil
		}
	}
	banner := "\n" + marker + "\n<aside aria-label=\"Second opinion\">A second-opinion version is available: <a href=\"" + html.EscapeString(pageName) + "\">Read the second opinion</a>.</aside>\n<!-- /second-opinion-banner -->\n"
	updated := append([]byte(nil), original[:body[1]]...)
	updated = append(updated, []byte(banner)...)
	updated = append(updated, original[body[1]:]...)
	return updated, true, nil
}

// secondOpinionStripBlocks removes every launcher block between the start and end marker comments. Markers count only
// as real comments on the raw-text-masked surface, so a marker inside a script, style or textarea string never matches.
func secondOpinionStripBlocks(page []byte, startMarker, endMarker string) []byte {
	comment := regexp.MustCompile(`(?s)<!--.*?-->`)
	for {
		commentSurface := append([]byte(nil), page...)
		for _, span := range secondOpinionRawHTML.FindAllIndex(page, -1) {
			for i := span[0]; i < span[1]; i++ {
				commentSurface[i] = ' '
			}
		}
		start, end := -1, -1
		for _, span := range comment.FindAllIndex(commentSurface, -1) {
			text := string(commentSurface[span[0]:span[1]])
			if start < 0 && text == startMarker {
				start = span[0]
			} else if start >= 0 && text == endMarker {
				end = span[1]
				break
			}
		}
		if end < 0 {
			return page
		}
		page = append(append([]byte(nil), page[:start]...), page[end:]...)
	}
}

// secondOpinionTreeMarks are the tree fields of one published page: its id plus two display values. An empty id means
// a plain run with no tree marks.
type secondOpinionTreeMarks struct {
	treeID, modelLabel, publishDate string
}

// secondOpinionDecorate first removes every launcher header block, copied iteration list and ai-report-second-opinion-*
// meta already on the page, then adds the direction and tree metas before the first real </head> and a header block
// right after the opening body tag: the tree id with a link to the parent page, then the slot and the full direction.
// Decorating twice gives the same bytes, and a plain direction with no tree marks only strips. Tags and markers count
// only on the masked surface, so text inside a script, style or comment never matches.
func secondOpinionDecorate(page []byte, direction secondOpinionDirection, tree secondOpinionTreeMarks) []byte {
	const headerStart, headerEnd = "<!-- second-opinion-header -->", "<!-- /second-opinion-header -->"
	insert := func(data []byte, at int, text string) []byte {
		return append(append(append([]byte(nil), data[:at]...), text...), data[at:]...)
	}
	page = secondOpinionStripBlocks(page, headerStart, headerEnd)
	page = secondOpinionStripBlocks(page, secondOpinionTreeStart, secondOpinionTreeEnd)
	// Either attribute order counts: the iteration list reads metas first-wins, so a copied one must never survive.
	spans := reportMetaTag.FindAllIndex(secondOpinionVisibleHTML(page), -1)
	for i := len(spans) - 1; i >= 0; i-- {
		for _, attribute := range reportMetaAttribute.FindAllSubmatch(page[spans[i][0]:spans[i][1]], -1) {
			if strings.EqualFold(string(attribute[1]), "name") && strings.HasPrefix(strings.ToLower(strings.TrimSpace(string(attribute[2])+string(attribute[3]))), "ai-report-second-opinion-") {
				page = append(append([]byte(nil), page[:spans[i][0]]...), page[spans[i][1]:]...)
				break
			}
		}
	}
	slot := direction.filenameSlot()
	if slot == "" && tree.treeID == "" {
		return page
	}
	text := html.EscapeString(direction.combinedText())
	metas, paragraphs := "", ""
	if slot != "" {
		metas = `<meta name="ai-report-second-opinion-preset" content="` + html.EscapeString(slot) + `"><meta name="ai-report-second-opinion-direction" content="` + text + `">`
	}
	if tree.treeID != "" {
		parent := html.EscapeString(cmp.Or(direction.Base, "index.html"))
		metas += `<meta name="ai-report-second-opinion-id" content="` + html.EscapeString(tree.treeID) + `"><meta name="ai-report-second-opinion-base" content="` + parent + `"><meta name="ai-report-second-opinion-model" content="` + html.EscapeString(tree.modelLabel) + `"><meta name="ai-report-second-opinion-date" content="` + html.EscapeString(tree.publishDate) + `">`
		paragraphs = `<p><strong>Iteration ` + html.EscapeString(tree.treeID) + `</strong>, built on <a href="` + parent + `">` + parent + `</a></p>`
	}
	if slot != "" {
		paragraphs += `<p><strong>Second-opinion preset:</strong> ` + html.EscapeString(slot) + `</p><p style="white-space:pre-line"><strong>Direction:</strong> ` + text + `</p>`
	}
	if head := regexp.MustCompile(`(?i)</head\s*>`).FindIndex(secondOpinionVisibleHTML(page)); head != nil {
		page = insert(page, head[0], metas)
		metas = ""
	}
	// overflow-wrap keeps one long direction token from failing the judge's horizontal-overflow check.
	block := headerStart + metas + `<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere">` + paragraphs + `</aside>` + headerEnd
	at := 0
	if body := secondOpinionBodyTag.FindIndex(secondOpinionVisibleHTML(page)); body != nil {
		at = body[1]
	}
	return insert(page, at, block)
}

const secondOpinionTreeStart, secondOpinionTreeEnd = "<!-- second-opinion-tree -->", "<!-- /second-opinion-tree -->"

// secondOpinionTreeList rebuilds the one iteration list on index.html from the bundle's tree page names, right after
// the opening body tag. A byte sort of the names is pre-order, because "-" sorts before "c". Every tree page on disk is
// listed, judge-failed and unverified ones too, so no child is orphaned; preset, model and date are display metas.
func secondOpinionTreeList(bundle string, original []byte) ([]byte, error) {
	updated := secondOpinionStripBlocks(original, secondOpinionTreeStart, secondOpinionTreeEnd)
	body := secondOpinionBodyTag.FindIndex(secondOpinionVisibleHTML(updated))
	if body == nil {
		return nil, errors.New("original has no opening body tag")
	}
	entries, err := os.ReadDir(bundle)
	if err != nil {
		return nil, err
	}
	var items strings.Builder
	for _, entry := range entries {
		match := secondOpinionTreePagePattern.FindStringSubmatch(entry.Name())
		if match == nil || !entry.Type().IsRegular() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(bundle, entry.Name()))
		if err != nil {
			return nil, err
		}
		metas := readReportMetas(secondOpinionVisibleHTML(data))
		preset := cmp.Or(metas["ai-report-second-opinion-preset"], "plain")
		depth := (len(match[1]) - len("v1")) / len("c01")
		fmt.Fprintf(&items, "<li style=\"margin-left:%gem\"><a href=\"%s\">%s</a> · %s · %s · %s</li>\n", float64(depth-1)*1.5, html.EscapeString(entry.Name()), html.EscapeString(match[1]), html.EscapeString(preset), html.EscapeString(metas["ai-report-second-opinion-model"]), html.EscapeString(metas["ai-report-second-opinion-date"]))
	}
	// No whitespace outside the markers, so removing the old block and inserting a fresh one never accumulates bytes.
	block := secondOpinionTreeStart + "\n<nav aria-label=\"Second-opinion iterations\"><p><strong>Second-opinion iterations</strong> (tree id · preset · model · date)</p>\n<ul>\n" + items.String() + "</ul></nav>\n" + secondOpinionTreeEnd
	return append(append(append([]byte(nil), updated[:body[1]]...), block...), updated[body[1]:]...), nil
}

// secondOpinionChangeList checks a tree page's "Changed in this iteration" list on the visible surface: it needs at
// least one #rev-* link, and every such link needs a real id="rev-*" attribute (whitespace before id, so data-id does
// not count). Each failure is one judge finding.
func secondOpinionChangeList(page []byte) []resultmodel.AIReportSecondOpinionFinding {
	visible := secondOpinionVisibleHTML(page)
	targets := map[string]bool{}
	for _, match := range regexp.MustCompile(`(?i)\sid\s*=\s*["']?(rev-[^"'\s>]*)`).FindAllSubmatch(visible, -1) {
		targets[string(match[1])] = true
	}
	links := regexp.MustCompile(`(?i)\shref\s*=\s*["']?#(rev-[^"'\s>]*)`).FindAllSubmatch(visible, -1)
	if len(links) == 0 {
		return []resultmodel.AIReportSecondOpinionFinding{{FindingKind: secondOpinionChangeListKind, FindingDetail: "no \"Changed in this iteration\" link to a #rev-* section"}}
	}
	findings := []resultmodel.AIReportSecondOpinionFinding{}
	for _, link := range links {
		if !targets[string(link[1])] {
			findings = append(findings, resultmodel.AIReportSecondOpinionFinding{FindingKind: secondOpinionChangeListKind, TargetPath: "#" + string(link[1]), FindingDetail: "change-list link has no matching id on the page"})
		}
	}
	return findings
}

func secondOpinionRender(bundle string, page []byte) (string, []resultmodel.AIReportSecondOpinionFinding) {
	findings := []resultmodel.AIReportSecondOpinionFinding{}
	fail := func(err error) (string, []resultmodel.AIReportSecondOpinionFinding) {
		return "error", append(findings, resultmodel.AIReportSecondOpinionFinding{FindingKind: judgeCodeError, FindingDetail: err.Error()})
	}
	browser := discoverJudgeBrowser()
	if browser == "" {
		return "skipped", append(findings, resultmodel.AIReportSecondOpinionFinding{FindingKind: judgeCodeSkipped, FindingDetail: "no browser engine found"})
	}
	scratch, err := os.MkdirTemp("", "ai-report-second-opinion-judge-")
	if err != nil {
		return fail(err)
	}
	defer os.RemoveAll(scratch)
	renderBundle := filepath.Join(scratch, "bundle")
	captures := filepath.Join(scratch, "captures")
	if err := os.MkdirAll(renderBundle, 0700); err != nil {
		return fail(err)
	}
	if err := os.MkdirAll(captures, 0700); err != nil {
		return fail(err)
	}
	files, err := secondOpinionSnapshot(bundle)
	if err != nil {
		return fail(err)
	}
	for name, item := range files {
		if name == "index.html" {
			continue
		}
		destination := filepath.Join(renderBundle, name)
		if item.info.IsDir() {
			err = os.MkdirAll(destination, 0700)
		} else {
			if err = os.MkdirAll(filepath.Dir(destination), 0700); err == nil {
				err = os.WriteFile(destination, item.data, 0600)
			}
		}
		if err != nil {
			return fail(err)
		}
	}
	if err := os.WriteFile(filepath.Join(renderBundle, "index.html"), page, 0600); err != nil {
		return fail(err)
	}
	report := judgeReport{SchemaVersion: 1, Bundle: renderBundle, Captures: []judgeCapture{}, Findings: []judgeFinding{}}
	err = judgeRender(&report, renderBundle, captures, browser)
	for _, finding := range report.Findings {
		detail := finding.Detail
		if finding.Kind == judgeCodeOverflow {
			detail = fmt.Sprintf("scroll width %d exceeds viewport %d", finding.ScrollWidth, finding.InnerWidth)
		}
		if finding.Status != 0 {
			detail = fmt.Sprintf("HTTP %d: %s", finding.Status, detail)
		}
		findings = append(findings, resultmodel.AIReportSecondOpinionFinding{FindingKind: finding.Kind, ViewportName: finding.Viewport, ColorScheme: finding.Scheme, TargetPath: finding.Target, FindingDetail: detail})
	}
	if err != nil {
		return fail(err)
	}
	if len(findings) > 0 {
		return "fail", findings
	}
	return "pass", findings
}
