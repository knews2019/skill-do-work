package toolboxcommands

import (
	"bytes"
	"context"
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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/ownedprocess"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const secondOpinionUsage = "ai-report-second-opinion [--backend codex|gemini|agy|claude [--model <slug>] | --all --exclude <own-cli>] [<bundle-dir>|latest]"

var secondOpinionModelPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// errSecondOpinionModelUnverified marks native evidence that names no single backend model; a finished page is then published under an unverified name instead of being discarded.
var errSecondOpinionModelUnverified = errors.New("actual model is absent or ambiguous")
var secondOpinionJudge = secondOpinionRender

type secondOpinionOptions struct {
	backend, model, exclude, bundle string
	all                             bool
}

func handleAIReportSecondOpinion(execution commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	options, err := secondOpinionParse(args)
	if err != nil {
		return usageResult(CommandAIReportSecondOpinion, err.Error()+"; "+secondOpinionUsage)
	}
	bundle, err := secondOpinionBundle(execution.RepositoryRoot, options.bundle)
	if err != nil {
		return usageResult(CommandAIReportSecondOpinion, err.Error())
	}
	templatePath := os.Getenv("DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT")
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return usageResult(CommandAIReportSecondOpinion, "use the suite's do-work-cli.sh launcher, or set DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT to its trusted shipped template: "+err.Error())
	}
	ctx, stop, signalStatus := mutationSignalContext()
	defer stop()
	payload := &resultmodel.AIReportSecondOpinionResult{BundlePath: bundle, InvocationMode: "targeted", BackendRuns: []resultmodel.AIReportSecondOpinionRun{}}
	backends := []string{options.backend}
	if options.all {
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
	var summary strings.Builder
	for _, backend := range backends {
		row := resultmodel.AIReportSecondOpinionRun{BackendName: backend, RequestedModel: options.model, RunState: "failed", JudgeVerdict: "not-run", BannerState: "not-inserted", JudgeFindings: []resultmodel.AIReportSecondOpinionFinding{}, OutsidePaths: []string{}, SeenModels: []string{}}
		if backend == options.exclude {
			row.RunState = "excluded"
			row.FailureReason = "excluded calling backend"
		} else if executable, lookupErr := exec.LookPath(backend); lookupErr != nil {
			row.RunState = "skipped"
			row.FailureReason = "backend is not installed"
		} else {
			secondOpinionRun(ctx, execution.RepositoryRoot, bundle, executable, string(template), &row, &result)
		}
		payload.BackendRuns = append(payload.BackendRuns, row)
		fmt.Fprintf(&summary, "%s: %s; model=%s; page=%s; reply=%s; judge=%s; banner=%s", row.BackendName, row.RunState, row.ActualModel, row.PagePath, row.ReplyPath, row.JudgeVerdict, row.BannerState)
		if row.FailureReason != "" {
			fmt.Fprintf(&summary, "; %s", row.FailureReason)
		}
		summary.WriteByte('\n')
		if row.RunState != "ran" && row.RunState != "excluded" || row.RunState == "ran" && row.JudgeVerdict != "pass" {
			severity := resultmodel.SeverityWarning
			if !options.all {
				severity = resultmodel.SeverityError
				result.Outcome = resultmodel.OutcomeFailure
			} else {
				result.Outcome = resultmodel.OutcomeFindings
			}
			result.Findings = append(result.Findings, toolboxFinding(CommandAIReportSecondOpinion, "AI-REPORT-SECOND-OPINION-"+strings.ToUpper(row.RunState), severity, row.OutsidePaths, backend+": "+row.FailureReason, resultmodel.FixabilityManual, "this backend did not produce a render-verified linked second opinion"))
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

type secondOpinionFile struct {
	info os.FileInfo
	data []byte
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
		item := secondOpinionFile{info: info}
		if info.Mode().IsRegular() {
			item.data, err = os.ReadFile(name)
			if err != nil {
				return err
			}
		}
		files[relative] = item
		return nil
	})
	return files, err
}
func secondOpinionPreserved(bundle string, before map[string]secondOpinionFile, stage, prefix string) error {
	after, err := secondOpinionSnapshot(bundle)
	if err != nil {
		return err
	}
	for name, old := range before {
		current, ok := after[name]
		if !ok || old.info.Mode() != current.info.Mode() || !bytes.Equal(old.data, current.data) || !os.SameFile(old.info, current.info) {
			return fmt.Errorf("backend changed existing bundle path: %s", name)
		}
	}
	for name, current := range after {
		if _, ok := before[name]; ok {
			continue
		}
		if name == stage {
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

func secondOpinionRun(ctx context.Context, repository, bundle, executable, template string, row *resultmodel.AIReportSecondOpinionRun, result *resultmodel.CommandResult) {
	fail := func(err error) { row.FailureReason = err.Error() }
	checkOutside := func() error {
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
		capture := secondOpinionCapture(probeContext, bundle, executable, "models")
		if err := checkOutside(); err != nil {
			fail(err)
			return
		}
		if err := secondOpinionPreserved(bundle, before, "", ""); err != nil {
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
	capture := secondOpinionCapture(probeContext, bundle, executable, argv...)
	if err := checkOutside(); err != nil {
		fail(err)
		return
	}
	if err := secondOpinionPreserved(bundle, before, "", ""); err != nil {
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
	stagePublished := false
	// The staging leaf is removed only once its bytes were published to a sibling page, or when the backend never wrote it. A finished page is never deleted on refusal; the failure reason names where it was kept.
	defer func() {
		info, statErr := bundleRoot.Stat(stageName)
		if statErr != nil {
			return
		}
		if stagePublished || info.Size() == 0 {
			_ = bundleRoot.Remove(stageName)
			return
		}
		row.FailureReason += "; staged page kept at " + filepath.ToSlash(filepath.Join(relative, stageName))
	}()
	prompt := strings.NewReplacer("{{REPOSITORY_ROOT}}", repository, "{{BUNDLE_PATH}}", bundle, "{{OUTPUT_FILE}}", stageName, "{{ASSET_PREFIX}}", prefix, "{{BACKLINK_REQUIREMENT}}", `Open the body with <p>Original report: <a href="index.html">index.html</a></p>, before the main content. Never modify index.html.`).Replace(template)
	if strings.Contains(prompt, "{{") {
		fail(errors.New("unresolved prompt placeholder"))
		return
	}
	capture = secondOpinionCapture(ctx, bundle, executable, secondOpinionArgv(row.BackendName, resolved, prompt, false)...)
	if err := checkOutside(); err != nil {
		fail(err)
		return
	}
	if err := secondOpinionPreserved(bundle, before, stageName, prefix); err != nil {
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
	if !unverified {
		if err := secondOpinionBacklink(page); err != nil {
			fail(err)
			return
		}
	}
	publishedModel := actual
	if unverified {
		// The finished page is kept under a name that says the model was not confirmed; no banner, no judge.
		publishedModel = resolved + "-unverified"
		row.ActualModel = ""
		row.FailureReason = fmt.Sprintf("actual model unverified: requested %s; seen %s", resolved, strings.Join(row.SeenModels, ", "))
	}
	pageName, replyName, err := secondOpinionNames(bundle, row.BackendName, publishedModel)
	if err != nil {
		fail(err)
		return
	}
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
	stagePublished = true
	// The backend knows the reserved stage name; the saved reply names its published page.
	reply = strings.NewReplacer("file://"+filepath.ToSlash(filepath.Join(bundle, stageName)), pageName, filepath.Join(bundle, stageName), pageName, filepath.ToSlash(filepath.Join(bundle, stageName)), pageName, stageName, pageName).Replace(reply)
	if err := publish(replyName, []byte(reply)); err != nil {
		fail(err)
		return
	}
	row.ReplyPath = filepath.ToSlash(filepath.Join(relative, replyName))
	if unverified {
		row.RunState = "unverified"
		return
	}
	verdict, findings := secondOpinionJudge(bundle, page)
	row.JudgeVerdict = verdict
	row.JudgeFindings = findings
	row.RunState = "ran"
	if verdict != "pass" {
		row.FailureReason = "render judge " + verdict + "; original left unchanged"
		return
	}
	if err := checkOutside(); err != nil {
		row.RunState = "failed"
		fail(err)
		return
	}
	original := before["index.html"]
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
	updated, inserted, err := secondOpinionBanner(original.data, pageName)
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
	result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: originalPath, Kind: "modified", Detail: "second-opinion banner only"})
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
type secondOpinionBuffer struct {
	bytes.Buffer
	truncated bool
}

func (buffer *secondOpinionBuffer) Write(data []byte) (int, error) {
	length := len(data)
	remaining := 16*1024*1024 - buffer.Len()
	if len(data) > remaining {
		data = data[:remaining]
		buffer.truncated = true
	}
	_, _ = buffer.Buffer.Write(data)
	return length, nil
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

func secondOpinionCapture(ctx context.Context, directory, executable string, args ...string) secondOpinionCaptured {
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
	var stdout, stderr secondOpinionBuffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		return secondOpinionCaptured{err: err}
	}
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
	return secondOpinionCaptured{stdout: stdout.String(), stderr: stderr.String(), err: err}
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

func secondOpinionNames(bundle, backend, model string) (string, string, error) {
	stem := backend + "-" + model
	for number := 1; ; number++ {
		suffix := ""
		if number > 1 {
			suffix = "-" + strconv.Itoa(number)
		}
		page, reply := "index-"+stem+suffix+".html", stem+suffix+"-reply.md"
		free := true
		for _, name := range []string{page, reply} {
			_, err := os.Lstat(filepath.Join(bundle, name))
			if err == nil {
				free = false
				break
			}
			if !os.IsNotExist(err) {
				return "", "", err
			}
		}
		if free {
			return page, reply, nil
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
