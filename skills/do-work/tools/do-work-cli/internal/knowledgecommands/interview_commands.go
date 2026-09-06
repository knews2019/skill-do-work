package knowledgecommands

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/sharedprimitives"
)

type interviewOptions struct {
	knowledgeRoot string
	template      string
	kb            string
	dryRun        bool
	commit        bool
	confirm       bool
}

type interviewTemplate struct {
	Name         string
	Description  string
	Version      string
	TopicCluster string
	Layers       []interviewLayer
	Exports      []string
	Raw          string
}

type interviewLayer struct{ ID, Title string }

// rootedCreateTestHook is nil in production and lets tests deterministically
// replace a validated parent immediately before rooted publication.
var rootedCreateTestHook func(repositoryRoot, relative string)

func parseInterviewOptions(arguments []string, mutable bool) (interviewOptions, error) {
	options := interviewOptions{kb: "kb"}
	seen := map[string]bool{}
	singleton := func(name string) error {
		if seen[name] {
			return fmt.Errorf("%s may be specified only once", name)
		}
		seen[name] = true
		return nil
	}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		value := func(name string) (string, error) {
			index++
			if index >= len(arguments) {
				return "", fmt.Errorf("%s requires a value", name)
			}
			return arguments[index], nil
		}
		var err error
		switch {
		case argument == "--knowledge-root":
			err = singleton("--knowledge-root")
			if err == nil {
				options.knowledgeRoot, err = value(argument)
			}
		case strings.HasPrefix(argument, "--knowledge-root="):
			err = singleton("--knowledge-root")
			if err == nil {
				options.knowledgeRoot = strings.TrimPrefix(argument, "--knowledge-root=")
			}
		case argument == "--template":
			err = singleton("--template")
			if err == nil {
				options.template, err = value(argument)
			}
		case strings.HasPrefix(argument, "--template="):
			err = singleton("--template")
			if err == nil {
				options.template = strings.TrimPrefix(argument, "--template=")
			}
		case argument == "--kb":
			err = singleton("--kb")
			if err == nil {
				options.kb, err = value(argument)
			}
		case strings.HasPrefix(argument, "--kb="):
			err = singleton("--kb")
			if err == nil {
				options.kb = strings.TrimPrefix(argument, "--kb=")
			}
		case argument == "--dry-run" && mutable:
			err = singleton("--dry-run")
			options.dryRun = err == nil
		case argument == "--commit" && mutable:
			err = singleton("--commit")
			options.commit = err == nil
		case argument == "--confirm" && mutable:
			err = singleton("--confirm")
			options.confirm = err == nil
		case !strings.HasPrefix(argument, "-") && options.template == "":
			err = singleton("--template")
			if err == nil {
				options.template = argument
			}
		default:
			return options, fmt.Errorf("unknown option %q", argument)
		}
		if err != nil {
			return options, err
		}
	}
	if options.dryRun && options.commit {
		return options, errors.New("--dry-run and --commit cannot be combined")
	}
	if options.template != "" && (filepath.Base(options.template) != options.template || options.template == "." || options.template == "..") {
		return options, fmt.Errorf("template %q must be a single slug", options.template)
	}
	return options, nil
}

func handleInterviewList(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseInterviewOptions(arguments, false)
	if err != nil {
		return usageResult(CommandInterviewList, err)
	}
	if options.template != "" {
		return usageResult(CommandInterviewList, errors.New("interview-list does not take a template"))
	}
	knowledgeRoot, err := locateKnowledgeRoot(executionContext.RepositoryRoot, options.knowledgeRoot)
	if err != nil {
		return interviewFailure(CommandInterviewList, "INTERVIEW-TEMPLATE-ROOT-MISSING", options.knowledgeRoot, err)
	}
	paths, _ := filepath.Glob(filepath.Join(knowledgeRoot, "interviews", "*.md"))
	sort.Strings(paths)
	if len(paths) == 0 {
		return interviewFindingResult(CommandInterviewList, "INTERVIEW-NO-TEMPLATES", resultmodel.SeverityWarning, filepath.ToSlash(filepath.Join(options.knowledgeRoot, "interviews")), "no interview templates were found", resultmodel.OutcomeFindings)
	}
	findings := make([]resultmodel.CommandFinding, 0, len(paths))
	for _, path := range paths {
		template, parseError := loadInterviewTemplate(path)
		display := filepath.ToSlash(path)
		if relative, relError := filepath.Rel(executionContext.RepositoryRoot, path); relError == nil {
			display = filepath.ToSlash(relative)
		}
		if parseError != nil {
			findings = append(findings, interviewFinding(CommandInterviewList, "INTERVIEW-TEMPLATE-MALFORMED", resultmodel.SeverityError, display, parseError.Error(), resultmodel.FixabilityManual, "repair the template frontmatter"))
			continue
		}
		if strings.TrimSuffix(filepath.Base(path), ".md") != template.Name {
			findings = append(findings, interviewFinding(CommandInterviewList, "INTERVIEW-TEMPLATE-SLUG-MISMATCH", resultmodel.SeverityError, display, fmt.Sprintf("frontmatter name=%q does not match filename", template.Name), resultmodel.FixabilityManual, "rename the template or correct its declared name"))
			continue
		}
		findings = append(findings, interviewFinding(CommandInterviewList, "INTERVIEW-TEMPLATE", resultmodel.SeverityInfo, display, fmt.Sprintf("%s version=%s layers=%d exports=%d description=%s", template.Name, template.Version, len(template.Layers), len(template.Exports), strings.Join(strings.Fields(template.Description), " ")), resultmodel.FixabilityManual, ""))
	}
	outcome := resultmodel.OutcomeSuccess
	for _, finding := range findings {
		if finding.Severity == resultmodel.SeverityError {
			outcome = resultmodel.OutcomeFindings
		}
	}
	return retargetInterviewResult(resultmodel.CommandResult{Outcome: outcome, Findings: findings}, options)
}

func handleInterviewStatus(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseInterviewOptions(arguments, false)
	if err != nil {
		return usageResult(CommandInterviewStatus, err)
	}
	template, session, _, failure := loadInterviewState(CommandInterviewStatus, executionContext.RepositoryRoot, options)
	if failure != nil {
		return *failure
	}
	oldVersion := stringValue(session["template_version"])
	migrationPending, migrationError := migrateInterviewSession(template, session)
	if migrationError != nil {
		return interviewFailure(CommandInterviewStatus, "INTERVIEW-MIGRATION-UNAVAILABLE", sessionPath(options.template), migrationError)
	}
	approved, entries := sessionCounts(session)
	evidence := fmt.Sprintf("template=%s status=%s approved=%d/%d entries=%d pending_layer=%s started=%s last_activity=%s review_runs=%d last_exported_at=%s", options.template, stringValue(session["status"]), approved, len(template.Layers), entries, stringValue(session["pending_layer"]), stringValue(session["started_at"]), stringValue(session["last_activity_at"]), intValue(session["review_runs"]), stringValue(session["last_exported_at"]))
	findings := []resultmodel.CommandFinding{interviewFinding(CommandInterviewStatus, "INTERVIEW-STATUS", resultmodel.SeverityInfo, sessionPath(options.template), evidence, resultmodel.FixabilityManual, "")}
	if migrationPending {
		findings = append(findings, interviewFinding(CommandInterviewStatus, "INTERVIEW-MIGRATION-PENDING", resultmodel.SeverityWarning, sessionPath(options.template), fmt.Sprintf("session %s is rendered in memory at template %s; next mutating command persists it", oldVersion, template.Version), resultmodel.FixabilityAutomatic, "status is read-only"))
	}
	return retargetInterviewResult(resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: findings}, options)
}

func handleInterviewVersions(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseInterviewOptions(arguments, false)
	if err != nil {
		return usageResult(CommandInterviewVersions, err)
	}
	if options.template == "" {
		return usageResult(CommandInterviewVersions, errors.New("--template is required"))
	}
	versionsRoot := filepath.Join(executionContext.RepositoryRoot, "do-work", "interview", options.template, "versions")
	entries, readError := os.ReadDir(versionsRoot)
	if os.IsNotExist(readError) {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{interviewFinding(CommandInterviewVersions, "INTERVIEW-NO-VERSIONS", resultmodel.SeverityInfo, filepath.ToSlash(filepath.Join("do-work/interview", options.template, "versions")), "no archived versions yet", resultmodel.FixabilityManual, "")}}
	}
	if readError != nil {
		return interviewFailure(CommandInterviewVersions, "INTERVIEW-VERSIONS-READ-FAILED", versionsRoot, readError)
	}
	type versionRecord struct {
		number  int
		name    string
		finding resultmodel.CommandFinding
	}
	records := []versionRecord{}
	malformed := []resultmodel.CommandFinding{}
	versionNames := map[int][]string{}
	pattern := regexp.MustCompile(`^v([0-9]+)-([0-9]{4}-[0-9]{2}-[0-9]{2})$`)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		match := pattern.FindStringSubmatch(entry.Name())
		display := filepath.ToSlash(filepath.Join("do-work/interview", options.template, "versions", entry.Name()))
		if match == nil {
			malformed = append(malformed, interviewFinding(CommandInterviewVersions, "INTERVIEW-VERSION-MALFORMED", resultmodel.SeverityWarning, display, "archive directory does not match v<N>-YYYY-MM-DD", resultmodel.FixabilityManual, "archives are immutable; inspect manually"))
			continue
		}
		number, _ := strconv.Atoi(match[1])
		versionNames[number] = append(versionNames[number], entry.Name())
		data, readErr := os.ReadFile(filepath.Join(versionsRoot, entry.Name(), "session.json"))
		if readErr != nil {
			malformed = append(malformed, interviewFinding(CommandInterviewVersions, "INTERVIEW-VERSION-MALFORMED", resultmodel.SeverityWarning, display, readErr.Error(), resultmodel.FixabilityManual, "archives are immutable; inspect manually"))
			continue
		}
		var session map[string]any
		if json.Unmarshal(data, &session) != nil {
			malformed = append(malformed, interviewFinding(CommandInterviewVersions, "INTERVIEW-VERSION-MALFORMED", resultmodel.SeverityWarning, display, "archived session.json is malformed", resultmodel.FixabilityManual, "archives are immutable; inspect manually"))
			continue
		}
		layers, entriesCount := sessionCounts(session)
		records = append(records, versionRecord{number: number, name: entry.Name(), finding: interviewFinding(CommandInterviewVersions, "INTERVIEW-VERSION", resultmodel.SeverityInfo, display, fmt.Sprintf("%s layers=%d entries=%d", entry.Name(), layers, entriesCount), resultmodel.FixabilityManual, "")})
	}
	for number, names := range versionNames {
		if len(names) < 2 {
			continue
		}
		sort.Strings(names)
		malformed = append(malformed, interviewFinding(CommandInterviewVersions, "INTERVIEW-VERSION-DUPLICATE-ID", resultmodel.SeverityError, filepath.ToSlash(filepath.Join("do-work/interview", options.template, "versions")), fmt.Sprintf("numeric version v%d is used by %s", number, strings.Join(names, ", ")), resultmodel.FixabilityManual, "archives are immutable; resolve the duplicate IDs manually"))
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].number != records[j].number {
			return records[i].number < records[j].number
		}
		return records[i].name < records[j].name
	})
	findings := append([]resultmodel.CommandFinding{}, malformed...)
	for _, record := range records {
		findings = append(findings, record.finding)
	}
	if len(findings) == 0 {
		findings = append(findings, interviewFinding(CommandInterviewVersions, "INTERVIEW-NO-VERSIONS", resultmodel.SeverityInfo, filepath.ToSlash(filepath.Join("do-work/interview", options.template, "versions")), "no archived versions yet", resultmodel.FixabilityManual, ""))
	}
	outcome := resultmodel.OutcomeSuccess
	if len(malformed) > 0 {
		outcome = resultmodel.OutcomeFindings
	}
	return retargetInterviewResult(resultmodel.CommandResult{Outcome: outcome, Findings: findings}, options)
}

func handleInterviewExport(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseInterviewOptions(arguments, true)
	if err != nil {
		return usageResult(CommandInterviewExport, err)
	}
	template, session, templatePath, failure := loadInterviewState(CommandInterviewExport, executionContext.RepositoryRoot, options)
	if failure != nil {
		return *failure
	}
	oldVersion := stringValue(session["template_version"])
	migrationPending, migrationError := migrateInterviewSession(template, session)
	if migrationError != nil {
		return interviewFailure(CommandInterviewExport, "INTERVIEW-MIGRATION-UNAVAILABLE", sessionPath(options.template), migrationError)
	}
	approved, _ := sessionCounts(session)
	if approved != len(template.Layers) || stringValue(session["status"]) != "complete" {
		return interviewFindingResult(CommandInterviewExport, "INTERVIEW-EXPORT-INCOMPLETE", resultmodel.SeverityWarning, sessionPath(options.template), fmt.Sprintf("approved=%d/%d; every layer must be approved", approved, len(template.Layers)), resultmodel.OutcomeRefused)
	}
	if intValue(session["review_runs"]) < 1 || session["review_completed_at"] == nil {
		return interviewFindingResult(CommandInterviewExport, "INTERVIEW-EXPORT-REVIEW-REQUIRED", resultmodel.SeverityWarning, sessionPath(options.template), "run the contradiction review before exporting", resultmodel.OutcomeRefused)
	}
	operationTime := nowUTC()
	now := operationTime.Format(time.RFC3339)
	previousExport := stringValue(session["last_exported_at"])
	lastActivity := stringValue(session["last_activity_at"])
	session["last_exported_at"] = now
	if stringValue(session["role_or_name_or_repo"]) == "" {
		session["role_or_name_or_repo"] = filepath.Base(executionContext.RepositoryRoot)
	}
	derived := deriveInterviewValues(session)
	rootData := map[string]any{"session": session, "template": map[string]any{"name": template.Name, "version": template.Version, "topic_cluster": template.TopicCluster}, "derived": derived, "all_layers": map[string]any{"entries": allInterviewEntries(session)}}
	for key, value := range mapValue(session["layers"]) {
		rootData[key] = value
	}
	renders, renderError := renderInterviewExports(template, rootData)
	if renderError != nil {
		return interviewFailure(CommandInterviewExport, "INTERVIEW-EXPORT-RENDER-FAILED", filepath.ToSlash(templatePath), renderError)
	}
	sessionBytes, _ := json.MarshalIndent(session, "", "  ")
	sessionBytes = append(sessionBytes, '\n')
	base := filepath.ToSlash(filepath.Join("do-work/interview", options.template))
	changelogPath := filepath.ToSlash(filepath.Join(base, "CHANGELOG.md"))
	changelogBytes, readError := os.ReadFile(filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(changelogPath)))
	if readError != nil && !os.IsNotExist(readError) {
		return interviewFailure(CommandInterviewExport, "INTERVIEW-CHANGELOG-READ-FAILED", changelogPath, readError)
	}
	if os.IsNotExist(readError) {
		changelogBytes = []byte("# Interview CHANGELOG — " + options.template + "\n\n")
	}
	filenames := make([]string, 0, len(renders))
	for name := range renders {
		filenames = append(filenames, name)
	}
	sort.Strings(filenames)
	if migrationPending {
		changelogBytes = append(changelogBytes, []byte(fmt.Sprintf("\n## %s — migrated %s to %s\nSession migration persisted before export.\n", operationTime.Format("2006-01-02 15:04"), oldVersion, template.Version))...)
	}
	changelogBytes = append(changelogBytes, []byte(fmt.Sprintf("\n## %s — exports written\n%s\n", operationTime.Format("2006-01-02 15:04"), strings.Join(filenames, ", ")))...)
	targets := []string{sessionPath(options.template), changelogPath}
	for _, name := range filenames {
		targets = append(targets, filepath.ToSlash(filepath.Join(base, "exports", name)))
	}
	createdDirs := absentDirectories(executionContext.RepositoryRoot, []string{filepath.ToSlash(filepath.Join(base, "exports"))})
	transaction := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{RepositoryRoot: executionContext.RepositoryRoot, TargetPaths: targets, CreatedDirectoryPaths: createdDirs, DryRun: options.dryRun, Commit: options.commit, CommitMessage: "Export " + options.template + " interview"}, func(recorder *gittransaction.MutationRecorder) error {
		if err := createTransactionDirectories(executionContext.RepositoryRoot, recorder, createdDirs); err != nil {
			return err
		}
		for _, name := range filenames {
			if err := publishTransactionFile(executionContext.RepositoryRoot, recorder, filepath.ToSlash(filepath.Join(base, "exports", name)), renders[name], 0o644, false); err != nil {
				return err
			}
		}
		if err := publishTransactionFile(executionContext.RepositoryRoot, recorder, sessionPath(options.template), sessionBytes, 0o644, false); err != nil {
			return err
		}
		return publishTransactionFile(executionContext.RepositoryRoot, recorder, changelogPath, changelogBytes, 0o644, false)
	})
	result := gittransaction.BuildCommandResult(CommandInterviewExport, transaction)
	if result.Outcome == resultmodel.OutcomeSuccess {
		freshness := "first export"
		if previousExport != "" {
			freshness = "fresh re-export"
			if exportedAt, exportErr := time.Parse(time.RFC3339, previousExport); exportErr == nil {
				if activityAt, activityErr := time.Parse(time.RFC3339, lastActivity); activityErr == nil && activityAt.After(exportedAt) {
					freshness = "session changed after the previous export"
				}
			}
		}
		result.Findings = append(result.Findings, interviewFinding(CommandInterviewExport, "INTERVIEW-EXPORT-FRESHNESS", resultmodel.SeverityInfo, sessionPath(options.template), freshness, resultmodel.FixabilityManual, ""))
	}
	if options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		result.Changes = plannedChanges(targets)
	}
	return retargetInterviewResult(result, options)
}

func handleInterviewIngest(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseInterviewOptions(arguments, true)
	if err != nil {
		return usageResult(CommandInterviewIngest, err)
	}
	template, session, _, failure := loadInterviewState(CommandInterviewIngest, executionContext.RepositoryRoot, options)
	if failure != nil {
		return *failure
	}
	oldVersion := stringValue(session["template_version"])
	migrationPending, migrationError := migrateInterviewSession(template, session)
	if migrationError != nil {
		return interviewFailure(CommandInterviewIngest, "INTERVIEW-MIGRATION-UNAVAILABLE", sessionPath(options.template), migrationError)
	}
	lastExported := stringValue(session["last_exported_at"])
	if lastExported == "" {
		return interviewFindingResult(CommandInterviewIngest, "INTERVIEW-INGEST-EXPORTS-STALE", resultmodel.SeverityWarning, sessionPath(options.template), "session has no successful export stamp; run interview-export first", resultmodel.OutcomeRefused)
	}
	if exportedAt, exportErr := time.Parse(time.RFC3339, lastExported); exportErr == nil {
		if activityAt, activityErr := time.Parse(time.RFC3339, stringValue(session["last_activity_at"])); activityErr == nil && activityAt.After(exportedAt) {
			return interviewFindingResult(CommandInterviewIngest, "INTERVIEW-INGEST-EXPORTS-STALE", resultmodel.SeverityWarning, sessionPath(options.template), "session changed after its exported artifacts; run interview-export again", resultmodel.OutcomeRefused)
		}
	}
	kbAbsolute, kbRelative, locateError := locateKnowledgeBase(executionContext.RepositoryRoot, options.kb)
	if locateError != nil {
		return interviewFindingResult(CommandInterviewIngest, "INTERVIEW-INGEST-KB-MISSING", resultmodel.SeverityWarning, options.kb, locateError.Error(), resultmodel.OutcomeRefused)
	}
	physicalRepositoryRoot, physicalRootError := physicalPath(filepath.Clean(executionContext.RepositoryRoot))
	if physicalRootError != nil || !pathInside(physicalRepositoryRoot, kbAbsolute) {
		return interviewFailure(CommandInterviewIngest, "INTERVIEW-INGEST-KB-OUTSIDE-REPOSITORY", options.kb, errors.New("mutating ingest target must be inside the repository"))
	}
	exportsRoot := filepath.Join(executionContext.RepositoryRoot, "do-work", "interview", options.template, "exports")
	exportEntries, readError := os.ReadDir(exportsRoot)
	if readError != nil || len(exportEntries) == 0 {
		return interviewFindingResult(CommandInterviewIngest, "INTERVIEW-INGEST-EXPORTS-MISSING", resultmodel.SeverityWarning, filepath.ToSlash(filepath.Join("do-work/interview", options.template, "exports")), "run interview-export first", resultmodel.OutcomeRefused)
	}
	stamp := nowUTC()
	inbox := filepath.Join(kbAbsolute, "raw", "inbox")
	capture := filepath.Join(kbAbsolute, "raw", "capture")
	writes := map[string][]byte{}
	reservedNames := map[string]bool{}
	for _, entry := range exportEntries {
		if entry.IsDir() {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 || !entry.Type().IsRegular() {
			return interviewFailure(CommandInterviewIngest, "INTERVIEW-INGEST-READ-FAILED", entry.Name(), errors.New("export is not a regular file"))
		}
		data, readErr := os.ReadFile(filepath.Join(exportsRoot, entry.Name()))
		if readErr != nil {
			return interviewFailure(CommandInterviewIngest, "INTERVIEW-INGEST-READ-FAILED", entry.Name(), readErr)
		}
		exportStem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		name := options.template + "-" + exportStem + ".md"
		name, collisionError := collisionName(inbox, capture, name, stamp, reservedNames)
		if collisionError != nil {
			return interviewFailure(CommandInterviewIngest, "INTERVIEW-INGEST-COLLISION-CHECK-FAILED", name, collisionError)
		}
		reservedNames[name] = true
		body := string(data)
		if strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			body = "```json\n" + strings.TrimSpace(body) + "\n```\n"
		}
		front := fmt.Sprintf("---\ntitle: %q\ntype: source-summary\ntopic_cluster: %q\nsources:\n  - %q\nconfidence: high\ncreated: %q\n---\n\n", template.Name+" — "+strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())), template.TopicCluster, "interview/"+options.template+"/exports/"+entry.Name(), stringValue(session["last_exported_at"]))
		writes[filepath.ToSlash(filepath.Join(kbRelative, "raw/inbox", name))] = []byte(front + body)
	}
	for _, layer := range template.Layers {
		layerValue := mapValue(mapValue(session["layers"])[layer.ID])
		entries := sliceValue(layerValue["entries"])
		name, collisionError := collisionName(inbox, capture, options.template+"-"+layer.ID+".md", stamp, reservedNames)
		if collisionError != nil {
			return interviewFailure(CommandInterviewIngest, "INTERVIEW-INGEST-COLLISION-CHECK-FAILED", layer.ID, collisionError)
		}
		reservedNames[name] = true
		var body strings.Builder
		fmt.Fprintf(&body, "# %s\n\n", layer.Title)
		confirmed, synthesized := 0, 0
		for _, raw := range entries {
			item := mapValue(raw)
			fmt.Fprintf(&body, "- **%s** — %s\n", stringValue(item["title"]), stringValue(item["summary"]))
			if stringValue(item["source_confidence"]) == "confirmed" {
				confirmed++
			} else {
				synthesized++
			}
		}
		confidence := "medium"
		if confirmed >= synthesized {
			confidence = "high"
		}
		front := fmt.Sprintf("---\ntitle: %q\ntype: concept\ntopic_cluster: %q\nsources:\n  - %q\nrelated:\n  - page: %q\n    rel: evidence-for\nconfidence: %q\ncreated: %q\n---\n\n", template.Name+" — "+layer.Title, template.TopicCluster, "interview/"+options.template+"/session.json", options.template+"-user-md", confidence, stringValue(session["last_exported_at"]))
		writes[filepath.ToSlash(filepath.Join(kbRelative, "raw/inbox", name))] = []byte(front + body.String())
	}
	queuePath := filepath.ToSlash(filepath.Join(kbRelative, "raw/_inbox_queue.md"))
	queueBytes, queueError := os.ReadFile(filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(queuePath)))
	if queueError != nil {
		return interviewFailure(CommandInterviewIngest, "INTERVIEW-INGEST-QUEUE-MISSING", queuePath, queueError)
	}
	if !validInboxQueue(queueBytes) {
		return interviewFindingResult(CommandInterviewIngest, "INTERVIEW-INGEST-QUEUE-MALFORMED", resultmodel.SeverityError, queuePath, "queue must contain the canonical Inbox Queue heading and Markdown table separator", resultmodel.OutcomeRefused)
	}
	if len(queueBytes) > 0 && queueBytes[len(queueBytes)-1] != '\n' {
		queueBytes = append(queueBytes, '\n')
	}
	paths := make([]string, 0, len(writes))
	for path := range writes {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		queueBytes = append(queueBytes, []byte(fmt.Sprintf("| %s | interview | ready | %s | normal |\n", filepath.Base(path), template.TopicCluster))...)
	}
	writes[queuePath] = queueBytes
	if stringValue(session["template_version"]) != oldVersion {
		sessionBytes, _ := json.MarshalIndent(session, "", "  ")
		writes[sessionPath(options.template)] = append(sessionBytes, '\n')
		if migrationPending {
			changelogPath := filepath.ToSlash(filepath.Join("do-work/interview", options.template, "CHANGELOG.md"))
			changelog, _ := os.ReadFile(filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(changelogPath)))
			if len(changelog) == 0 {
				changelog = []byte("# Interview CHANGELOG — " + options.template + "\n\n")
			}
			changelog = append(changelog, []byte(fmt.Sprintf("\n## %s — migrated %s to %s\nSession migration persisted before ingest.\n", stamp.Format("2006-01-02 15:04"), oldVersion, template.Version))...)
			writes[changelogPath] = changelog
		}
	}
	targets := sharedprimitives.UniqueSortedStrings(mapKeysBytes(writes))
	createdDirs := absentDirectories(executionContext.RepositoryRoot, []string{filepath.ToSlash(filepath.Join(kbRelative, "raw/inbox"))})
	transaction := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{RepositoryRoot: executionContext.RepositoryRoot, TargetPaths: targets, CreatedDirectoryPaths: createdDirs, DryRun: options.dryRun, Commit: options.commit, CommitMessage: "Ingest " + options.template + " interview"}, func(recorder *gittransaction.MutationRecorder) error {
		if err := createTransactionDirectories(executionContext.RepositoryRoot, recorder, createdDirs); err != nil {
			return err
		}
		for _, path := range targets {
			if err := publishTransactionFile(executionContext.RepositoryRoot, recorder, path, writes[path], 0o644, false); err != nil {
				return err
			}
		}
		return nil
	})
	result := gittransaction.BuildCommandResult(CommandInterviewIngest, transaction)
	if options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		result.Changes = plannedChanges(targets)
	}
	return retargetInterviewResult(result, options)
}

func handleInterviewReset(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseInterviewOptions(arguments, true)
	if err != nil {
		return usageResult(CommandInterviewReset, err)
	}
	if options.template == "" {
		return usageResult(CommandInterviewReset, errors.New("--template is required"))
	}
	if !options.confirm {
		return interviewFindingResult(CommandInterviewReset, "INTERVIEW-RESET-CONFIRMATION-REQUIRED", resultmodel.SeverityWarning, sessionPath(options.template), "reset requires --confirm after the archive plan is reviewed", resultmodel.OutcomeRefused)
	}
	knowledgeRoot, locateError := locateKnowledgeRoot(executionContext.RepositoryRoot, options.knowledgeRoot)
	if locateError != nil {
		return interviewFailure(CommandInterviewReset, "INTERVIEW-TEMPLATE-ROOT-MISSING", options.knowledgeRoot, locateError)
	}
	template, templateError := loadInterviewTemplate(filepath.Join(knowledgeRoot, "interviews", options.template+".md"))
	if templateError != nil {
		return interviewFailure(CommandInterviewReset, "INTERVIEW-TEMPLATE-MALFORMED", options.template, templateError)
	}
	base := filepath.Join(executionContext.RepositoryRoot, "do-work", "interview", options.template)
	sessionBytes, readError := readRegularInterviewFile(filepath.Join(base, "session.json"))
	if readError != nil {
		return interviewFailure(CommandInterviewReset, "INTERVIEW-SESSION-MISSING", sessionPath(options.template), readError)
	}
	next, versionError := nextInterviewVersion(filepath.Join(base, "versions"))
	if versionError != nil {
		return interviewFailure(CommandInterviewReset, "INTERVIEW-RESET-VERSIONS-INVALID", filepath.ToSlash(filepath.Join("do-work/interview", options.template, "versions")), versionError)
	}
	operationTime := nowUTC()
	versionName := fmt.Sprintf("v%d-%s", next, operationTime.Format("2006-01-02"))
	archiveRelative := filepath.ToSlash(filepath.Join("do-work/interview", options.template, "versions", versionName))
	archiveAbsolute := filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(archiveRelative))
	if _, statError := os.Lstat(archiveAbsolute); statError == nil {
		return interviewFindingResult(CommandInterviewReset, "INTERVIEW-RESET-ARCHIVE-COLLISION", resultmodel.SeverityWarning, archiveRelative, "the exact archive already exists and is immutable", resultmodel.OutcomeRefused)
	} else if !os.IsNotExist(statError) {
		return interviewFailure(CommandInterviewReset, "INTERVIEW-RESET-ARCHIVE-INSPECT-FAILED", archiveRelative, statError)
	}
	writes := map[string][]byte{filepath.ToSlash(filepath.Join(archiveRelative, "session.json")): sessionBytes}
	oldPaths := []string{}
	for _, directory := range []string{"checkpoints", "exports"} {
		workingDirectory := filepath.Join(base, directory)
		files, collectError := collectInterviewStateFiles(executionContext.RepositoryRoot, workingDirectory, filepath.ToSlash(filepath.Join(archiveRelative, directory)))
		if collectError != nil {
			return interviewFailure(CommandInterviewReset, "INTERVIEW-RESET-STATE-READ-FAILED", filepath.ToSlash(filepath.Join("do-work/interview", options.template, directory)), collectError)
		}
		for _, file := range files {
			writes[file.archivePath] = file.data
			oldPaths = append(oldPaths, file.workingPath)
		}
	}
	now := operationTime.Format(time.RFC3339)
	fresh := map[string]any{"template": options.template, "template_version": template.Version, "session_id": randomID(), "started_at": now, "last_activity_at": now, "status": "in_progress", "pending_layer": firstLayer(template), "previous_version": fmt.Sprintf("v%d", next), "review_completed_at": nil, "review_runs": 0, "last_exported_at": nil, "layers": map[string]any{}}
	freshBytes, _ := json.MarshalIndent(fresh, "", "  ")
	freshBytes = append(freshBytes, '\n')
	writes[sessionPath(options.template)] = freshBytes
	changelogPath := filepath.ToSlash(filepath.Join("do-work/interview", options.template, "CHANGELOG.md"))
	changelogAbsolute := filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(changelogPath))
	changelog, changelogError := readRegularInterviewFile(changelogAbsolute)
	if os.IsNotExist(changelogError) {
		changelog = []byte("# Interview CHANGELOG — " + options.template + "\n\n")
	} else if changelogError != nil {
		return interviewFailure(CommandInterviewReset, "INTERVIEW-CHANGELOG-READ-FAILED", changelogPath, changelogError)
	}
	changelog = append(changelog, []byte(fmt.Sprintf("\n## %s — reset (archived as v%d)\nFresh session started; v%d retained for reference.\n", operationTime.Format("2006-01-02 15:04"), next, next))...)
	writes[changelogPath] = changelog
	targets := append([]string{}, oldPaths...)
	for path := range writes {
		targets = append(targets, path)
	}
	targets = sharedprimitives.UniqueSortedStrings(targets)
	dirs := []string{filepath.ToSlash(filepath.Dir(archiveRelative)), archiveRelative}
	for path := range writes {
		if strings.HasPrefix(path, archiveRelative+"/") {
			for directory := filepath.ToSlash(filepath.Dir(path)); strings.HasPrefix(directory, archiveRelative); directory = filepath.ToSlash(filepath.Dir(directory)) {
				dirs = append(dirs, directory)
				if directory == archiveRelative {
					break
				}
			}
		}
	}
	createdDirs := absentDirectories(executionContext.RepositoryRoot, sharedprimitives.UniqueSortedStrings(dirs))
	transaction := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{RepositoryRoot: executionContext.RepositoryRoot, TargetPaths: targets, CreatedDirectoryPaths: createdDirs, DryRun: options.dryRun, Commit: options.commit, CommitMessage: "Reset " + options.template + " interview"}, func(recorder *gittransaction.MutationRecorder) error {
		if err := createTransactionDirectories(executionContext.RepositoryRoot, recorder, createdDirs); err != nil {
			return err
		}
		for _, path := range sharedprimitives.UniqueSortedStrings(mapKeysBytes(writes)) {
			data := writes[path]
			if err := publishTransactionFile(executionContext.RepositoryRoot, recorder, path, data, 0o644, false); err != nil {
				return err
			}
		}
		root, err := os.OpenRoot(executionContext.RepositoryRoot)
		if err != nil {
			return err
		}
		defer root.Close()
		for _, path := range oldPaths {
			if err := recorder.RecordTouched(path); err != nil {
				return err
			}
			if err := root.Remove(filepath.FromSlash(path)); err != nil {
				return err
			}
		}
		return nil
	})
	result := gittransaction.BuildCommandResult(CommandInterviewReset, transaction)
	if options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		result.Changes = plannedChanges(targets)
	}
	return retargetInterviewResult(result, options)
}

type interviewStateFile struct {
	workingPath string
	archivePath string
	data        []byte
}

func collectInterviewStateFiles(repositoryRoot, workingDirectory, archiveDirectory string) ([]interviewStateFile, error) {
	files := []interviewStateFile{}
	err := filepath.WalkDir(workingDirectory, func(path string, entry os.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 || !entry.Type().IsRegular() {
			return fmt.Errorf("%s is not a regular file", path)
		}
		data, readError := readRegularInterviewFile(path)
		if readError != nil {
			return readError
		}
		workingRelative, relativeError := filepath.Rel(repositoryRoot, path)
		if relativeError != nil {
			return relativeError
		}
		insideDirectory, relativeError := filepath.Rel(workingDirectory, path)
		if relativeError != nil || insideDirectory == ".." || strings.HasPrefix(insideDirectory, ".."+string(filepath.Separator)) {
			return fmt.Errorf("state path %s escaped its declared directory", path)
		}
		files = append(files, interviewStateFile{workingPath: filepath.ToSlash(workingRelative), archivePath: filepath.ToSlash(filepath.Join(archiveDirectory, insideDirectory)), data: data})
		return nil
	})
	if os.IsNotExist(err) {
		return files, nil
	}
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].workingPath < files[j].workingPath })
	return files, nil
}

func locateKnowledgeRoot(repositoryRoot, supplied string) (string, error) {
	candidates := []string{}
	if supplied != "" {
		candidates = append(candidates, supplied)
	} else {
		candidates = append(candidates, filepath.Join(repositoryRoot, ".claude/skills/do-work-knowledge"), filepath.Join(repositoryRoot, "skills/do-work-knowledge"))
	}
	for _, candidate := range candidates {
		if !filepath.IsAbs(candidate) {
			candidate = filepath.Join(repositoryRoot, candidate)
		}
		physical, err := physicalPath(filepath.Clean(candidate))
		if err != nil {
			continue
		}
		info, err := os.Stat(physical)
		if err == nil && info.IsDir() {
			return physical, nil
		}
	}
	return "", fmt.Errorf("knowledge root not found; supply --knowledge-root or install do-work-knowledge")
}

func loadInterviewTemplate(path string) (interviewTemplate, error) {
	data, err := readRegularInterviewFile(path)
	if err != nil {
		return interviewTemplate{}, err
	}
	raw := strings.TrimPrefix(strings.ReplaceAll(string(data), "\r\n", "\n"), "\ufeff")
	if !strings.HasPrefix(raw, "---\n") {
		return interviewTemplate{}, errors.New("template has no YAML frontmatter")
	}
	end := strings.Index(raw[4:], "\n---\n")
	if end < 0 {
		return interviewTemplate{}, errors.New("template frontmatter is unterminated")
	}
	front := raw[4 : 4+end]
	template := interviewTemplate{Raw: raw}
	section := ""
	descriptionLines := []string{}
	for _, line := range strings.Split(front, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(line, " ") && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key, value := strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			section = key
			switch key {
			case "name":
				template.Name = value
			case "version":
				template.Version = value
			case "topic_cluster":
				template.TopicCluster = value
			case "description":
				if value != "|" {
					template.Description = value
				}
			}
			continue
		}
		if section == "description" && strings.HasPrefix(line, "  ") && !strings.HasPrefix(trimmed, "-") {
			descriptionLines = append(descriptionLines, trimmed)
			continue
		}
		if section == "layers" && strings.HasPrefix(trimmed, "- id:") {
			template.Layers = append(template.Layers, interviewLayer{ID: strings.TrimSpace(strings.TrimPrefix(trimmed, "- id:"))})
			continue
		}
		if section == "layers" && strings.HasPrefix(trimmed, "title:") && len(template.Layers) > 0 {
			template.Layers[len(template.Layers)-1].Title = strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "title:")), "\"'")
			continue
		}
		if section == "exports" && strings.HasPrefix(trimmed, "- path:") {
			template.Exports = append(template.Exports, strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "- path:")), "`\"'"))
		}
	}
	if template.Description == "" {
		template.Description = strings.Join(descriptionLines, " ")
	}
	if template.Name == "" || template.Version == "" || len(template.Layers) == 0 {
		return interviewTemplate{}, errors.New("template requires name, version, and at least one layer")
	}
	for _, export := range template.Exports {
		if filepath.Base(export) != export || export == "." || export == ".." {
			return interviewTemplate{}, fmt.Errorf("unsafe export path %q", export)
		}
	}
	return template, nil
}

func loadInterviewState(command, repositoryRoot string, options interviewOptions) (interviewTemplate, map[string]any, string, *resultmodel.CommandResult) {
	if options.template == "" {
		failure := usageResult(command, errors.New("--template is required"))
		return interviewTemplate{}, nil, "", &failure
	}
	knowledgeRoot, err := locateKnowledgeRoot(repositoryRoot, options.knowledgeRoot)
	if err != nil {
		failure := interviewFailure(command, "INTERVIEW-TEMPLATE-ROOT-MISSING", options.knowledgeRoot, err)
		return interviewTemplate{}, nil, "", &failure
	}
	templatePath := filepath.Join(knowledgeRoot, "interviews", options.template+".md")
	template, err := loadInterviewTemplate(templatePath)
	if err != nil {
		failure := interviewFailure(command, "INTERVIEW-TEMPLATE-MALFORMED", filepath.ToSlash(templatePath), err)
		return interviewTemplate{}, nil, templatePath, &failure
	}
	data, err := readRegularInterviewFile(filepath.Join(repositoryRoot, filepath.FromSlash(sessionPath(options.template))))
	if err != nil {
		failure := interviewFailure(command, "INTERVIEW-SESSION-MISSING", sessionPath(options.template), err)
		return interviewTemplate{}, nil, templatePath, &failure
	}
	var session map[string]any
	if err := json.Unmarshal(data, &session); err != nil {
		failure := interviewFailure(command, "INTERVIEW-SESSION-MALFORMED", sessionPath(options.template), err)
		return interviewTemplate{}, nil, templatePath, &failure
	}
	if stringValue(session["template"]) != options.template {
		failure := interviewFailure(command, "INTERVIEW-SESSION-TEMPLATE-MISMATCH", sessionPath(options.template), fmt.Errorf("session template=%q", stringValue(session["template"])))
		return interviewTemplate{}, nil, templatePath, &failure
	}
	return template, session, templatePath, nil
}

func readRegularInterviewFile(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", path)
	}
	return os.ReadFile(path)
}

func migrateInterviewSession(template interviewTemplate, session map[string]any) (bool, error) {
	old := stringValue(session["template_version"])
	if old == "" {
		old = "1.0.0"
	}
	oldMajor, oldErr := semverMajor(old)
	newMajor, newErr := semverMajor(template.Version)
	if oldErr != nil || newErr != nil {
		return false, errors.New("template and session versions must be bare semver")
	}
	if oldMajor > newMajor {
		return false, nil
	}
	if oldMajor == newMajor {
		// The strict comparator reports parsing separately, so a version the old
		// lenient parser scored as "equal" is now an explicit refusal rather than a
		// silently skipped version stamp.
		versionOrdering, versionsParsed := sharedprimitives.CompareSemanticVersions(old, template.Version)
		if !versionsParsed {
			return false, errors.New("template and session versions must be bare semver")
		}
		if versionOrdering < 0 {
			session["template_version"] = template.Version
		}
		return false, nil
	}
	for major := oldMajor; major < newMajor; major++ {
		if !strings.Contains(template.Raw, fmt.Sprintf("## Migration from v%d.x", major)) {
			return false, fmt.Errorf("session %s needs undocumented migration to %s", old, template.Version)
		}
	}
	for _, layerRaw := range mapValue(session["layers"]) {
		layer := mapValue(layerRaw)
		for _, entryRaw := range sliceValue(layer["entries"]) {
			entry := mapValue(entryRaw)
			details := mapValue(entry["details"])
			if interruptions, ok := details["interruptions"].([]any); ok {
				converted := make([]any, 0, len(interruptions))
				for _, item := range interruptions {
					if text, stringItem := item.(string); stringItem {
						converted = append(converted, map[string]any{"source": text, "priority": "medium"})
					} else {
						converted = append(converted, item)
					}
				}
				details["interruptions"] = converted
			}
			for _, windowRaw := range sliceValue(details["time_windows"]) {
				window := mapValue(windowRaw)
				if _, exists := window["days"]; !exists {
					window["days"] = []any{"Mon", "Tue", "Wed", "Thu", "Fri"}
				}
			}
		}
	}
	session["template_version"] = template.Version
	return true, nil
}

func renderInterviewExports(template interviewTemplate, root map[string]any) (map[string][]byte, error) {
	renders := map[string][]byte{}
	for _, name := range template.Exports {
		marker := "### `" + name + "`"
		start := strings.Index(template.Raw, marker)
		if start < 0 {
			return nil, fmt.Errorf("export %s has no declared render block", name)
		}
		fence := strings.Index(template.Raw[start:], "```")
		if fence < 0 {
			return nil, fmt.Errorf("export %s render block is not fenced", name)
		}
		bodyStart := start + fence + 3
		if newline := strings.Index(template.Raw[bodyStart:], "\n"); newline >= 0 {
			bodyStart += newline + 1
		}
		end := strings.Index(template.Raw[bodyStart:], "\n```")
		if end < 0 {
			return nil, fmt.Errorf("export %s render block is unterminated", name)
		}
		rendered, err := renderTemplateBlock(template.Raw[bodyStart:bodyStart+end], root, root, 0, 1, strings.EqualFold(filepath.Ext(name), ".json"))
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", name, err)
		}
		renders[name] = []byte(strings.TrimSpace(rendered) + "\n")
	}
	return renders, nil
}

func renderTemplateBlock(template string, root, current map[string]any, itemIndex, itemCount int, jsonOutput bool) (string, error) {
	for {
		start := strings.Index(template, "{{#each ")
		if start < 0 {
			break
		}
		exprEnd := strings.Index(template[start:], "}}")
		if exprEnd < 0 {
			return "", errors.New("unterminated each expression")
		}
		exprEnd += start
		closeIndex := matchingTemplateClose(template, exprEnd+2, "each")
		if closeIndex < 0 {
			return "", errors.New("unclosed each block")
		}
		expr := strings.TrimSpace(template[start+8 : exprEnd])
		body := template[exprEnd+2 : closeIndex]
		items := evaluateEach(expr, root, current)
		var replacement strings.Builder
		for index, raw := range items {
			item := mapValue(raw)
			rendered, err := renderTemplateBlock(body, root, item, index, len(items), jsonOutput)
			if err != nil {
				return "", err
			}
			replacement.WriteString(rendered)
		}
		template = template[:start] + replacement.String() + template[closeIndex+len("{{/each}}"):]
	}
	for {
		start := strings.Index(template, "{{#unless @last}}")
		if start < 0 {
			break
		}
		end := strings.Index(template[start:], "{{/unless}}")
		if end < 0 {
			return "", errors.New("unclosed unless block")
		}
		end += start
		body := template[start+len("{{#unless @last}}") : end]
		if itemIndex == itemCount-1 {
			body = ""
		}
		template = template[:start] + body + template[end+len("{{/unless}}"):]
	}
	var renderError error
	if jsonOutput {
		quotedPattern := regexp.MustCompile(`"\{\{\s*([^{}]+?)\s*\}\}"`)
		template = quotedPattern.ReplaceAllStringFunc(template, func(token string) string {
			expr := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(token, `"{{`), `}}"`))
			if strings.HasPrefix(expr, "json ") || strings.HasPrefix(expr, "json_entries ") {
				return token
			}
			value, found := resolveScopedFound(root, current, expr)
			if !found {
				renderError = fmt.Errorf("unresolved template field %q", expr)
				return ""
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				renderError = fmt.Errorf("encode template field %q: %w", expr, err)
				return ""
			}
			return string(encoded)
		})
	}
	pattern := regexp.MustCompile(`\{\{\s*([^{}]+?)\s*\}\}`)
	template = pattern.ReplaceAllStringFunc(template, func(token string) string {
		expr := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(token, "{{"), "}}"))
		var value any
		switch {
		case strings.HasPrefix(expr, "json_entries "):
			layer := strings.TrimSpace(strings.TrimPrefix(expr, "json_entries "))
			value = mapValue(resolveValue(root, "layers."+layer))["entries"]
			bytes, _ := json.Marshal(value)
			return string(bytes)
		case strings.HasPrefix(expr, "json "):
			value = resolveScoped(root, current, strings.TrimSpace(strings.TrimPrefix(expr, "json ")))
			bytes, _ := json.Marshal(value)
			return string(bytes)
		default:
			var found bool
			value, found = resolveScopedFound(root, current, expr)
			if !found {
				renderError = fmt.Errorf("unresolved template field %q", expr)
				return ""
			}
		}
		if jsonOutput {
			encoded, err := json.Marshal(value)
			if err != nil {
				renderError = fmt.Errorf("encode template field %q: %w", expr, err)
				return ""
			}
			return string(encoded)
		}
		return displayTemplateValue(value)
	})
	return template, renderError
}

func matchingTemplateClose(template string, from int, kind string) int {
	open, close, depth := "{{#"+kind+" ", "{{/"+kind+"}}", 1
	position := from
	for position < len(template) {
		nextOpen := strings.Index(template[position:], open)
		nextClose := strings.Index(template[position:], close)
		if nextClose < 0 {
			return -1
		}
		if nextOpen >= 0 && nextOpen < nextClose {
			depth++
			position += nextOpen + len(open)
			continue
		}
		position += nextClose
		depth--
		if depth == 0 {
			return position
		}
		position += len(close)
	}
	return -1
}

func evaluateEach(expr string, root, current map[string]any) []any {
	sortPart := ""
	if index := strings.Index(expr, " sorted by "); index >= 0 {
		sortPart = expr[index+11:]
		expr = expr[:index]
	}
	wherePart := ""
	if index := strings.Index(expr, " where "); index >= 0 {
		wherePart = expr[index+7:]
		expr = expr[:index]
	}
	items := append([]any(nil), sliceValue(resolveScoped(root, current, strings.TrimSpace(expr)))...)
	if wherePart != "" {
		filtered := []any{}
		for _, item := range items {
			if evaluatePredicate(wherePart, mapValue(item)) {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}
	if sortPart != "" {
		terms := strings.Split(sortPart, ",")
		sort.SliceStable(items, func(i, j int) bool {
			left, right := mapValue(items[i]), mapValue(items[j])
			for _, term := range terms {
				fields := strings.Fields(strings.TrimSpace(term))
				if len(fields) == 0 {
					continue
				}
				a, b := displayTemplateValue(resolveValue(left, fields[0])), displayTemplateValue(resolveValue(right, fields[0]))
				if a == b {
					continue
				}
				if len(fields) > 1 && fields[1] == "desc" {
					return a > b
				}
				return a < b
			}
			return false
		})
	}
	return items
}

func evaluatePredicate(predicate string, item map[string]any) bool {
	for _, orPart := range strings.Split(predicate, " or ") {
		all := true
		for _, andPart := range strings.Split(orPart, " and ") {
			condition := strings.TrimSpace(andPart)
			fields := strings.Fields(condition)
			if len(fields) < 2 {
				all = false
				break
			}
			value := resolveValue(item, fields[0])
			ok := false
			switch fields[1] {
			case "exists":
				ok = value != nil && displayTemplateValue(value) != ""
			case "==", "!=", "contains", "mentions":
				if len(fields) < 3 {
					break
				}
				expected := strings.Trim(strings.Join(fields[2:], " "), "\"'")
				actual := displayTemplateValue(value)
				switch fields[1] {
				case "==":
					ok = strings.EqualFold(actual, expected)
				case "!=":
					ok = !strings.EqualFold(actual, expected)
				default:
					ok = strings.Contains(strings.ToLower(actual), strings.ToLower(expected))
				}
			}
			if !ok {
				all = false
				break
			}
		}
		if all {
			return true
		}
	}
	return false
}

func deriveInterviewValues(session map[string]any) map[string]any {
	derived := map[string]any{}
	rhythms := interviewLayerEntries(session, "operating_rhythms")
	parts, overnight := []string{}, []any{}
	for _, raw := range rhythms {
		item := mapValue(raw)
		details := mapValue(item["details"])
		for _, key := range []string{"energy_pattern", "non_calendar_reality"} {
			if value := stringValue(details[key]); value != "" {
				parts = append(parts, value)
				if key == "non_calendar_reality" {
					overnight = append(overnight, value)
				}
			}
		}
	}
	derived["rhythm_synthesis"] = strings.Join(parts, " ")
	inputCounts := map[string]int{}
	for _, raw := range interviewLayerEntries(session, "recurring_decisions") {
		details := mapValue(mapValue(raw)["details"])
		for _, input := range sliceValue(details["decision_inputs"]) {
			inputCounts[stringValue(input)]++
		}
	}
	authoritative, advisory := []any{}, []any{}
	for input, count := range inputCounts {
		if count >= 2 {
			authoritative = append(authoritative, input)
		} else {
			advisory = append(advisory, input)
		}
		overnight = append(overnight, input)
	}
	sort.Slice(authoritative, func(i, j int) bool { return stringValue(authoritative[i]) < stringValue(authoritative[j]) })
	sort.Slice(advisory, func(i, j int) bool { return stringValue(advisory[i]) < stringValue(advisory[j]) })
	derived["authoritative_inputs"], derived["advisory_inputs"], derived["overnight_scan_sources"] = authoritative, advisory, uniqueAnyStrings(overnight)
	tacit := []any{}
	for _, raw := range interviewLayerEntries(session, "institutional_knowledge") {
		item := mapValue(raw)
		location := strings.ToLower(stringValue(mapValue(item["details"])["where_it_lives"]))
		if strings.Contains(location, "head") || strings.Contains(location, "undocumented") {
			tacit = append(tacit, item)
		}
	}
	derived["tacit_knowledge"] = tacit
	derived["stakeholder_tones"] = []any{}
	timeBlocks := []any{}
	for _, raw := range rhythms {
		item := mapValue(raw)
		details := mapValue(item["details"])
		pattern := strings.ToLower(stringValue(details["energy_pattern"]))
		for _, windowRaw := range sliceValue(details["time_windows"]) {
			window := mapValue(windowRaw)
			kind := "reactive"
			joined := pattern + " " + strings.ToLower(stringValue(window["label"]))
			if strings.Contains(joined, "deep") || strings.Contains(joined, "focus") {
				kind = "deep_work"
			} else if strings.Contains(joined, "admin") || strings.Contains(joined, "email") {
				kind = "admin"
			}
			timeBlocks = append(timeBlocks, map[string]any{"label": window["label"], "days": window["days"], "start": window["start"], "end": window["end"], "type": kind, "source_entries": []any{"operating_rhythms." + stringValue(item["entry_id"])}})
		}
	}
	derived["time_blocks"] = timeBlocks
	derived["avoid_windows"], derived["standing_slots"] = []any{}, []any{}
	return derived
}

func publishTransactionFile(repositoryRoot string, recorder *gittransaction.MutationRecorder, relative string, contents []byte, mode os.FileMode, private bool) error {
	relative = filepath.ToSlash(filepath.Clean(relative))
	absolute := filepath.Join(repositoryRoot, filepath.FromSlash(relative))
	info, err := os.Lstat(absolute)
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("target %s is not a regular file", relative)
		}
		if err := recorder.RecordTouched(relative); err != nil {
			return err
		}
		if err := replaceRootedFile(repositoryRoot, relative, contents, info.Mode()); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		if rootedCreateTestHook != nil {
			rootedCreateTestHook(repositoryRoot, relative)
		}
		// Ownership is the successful rooted create, so the transaction learns about the
		// destination only once this invocation has actually published it.
		if err := createRootedFile(repositoryRoot, relative, contents, mode); err != nil {
			return err
		}
		if err := recorder.RecordCreated(relative); err != nil {
			return err
		}
	} else {
		return err
	}
	if private {
		return recorder.RecordPublished(relative)
	}
	return nil
}

func createRootedFile(repositoryRoot, relative string, contents []byte, mode os.FileMode) error {
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	rel := filepath.FromSlash(relative)
	file, err := root.OpenFile(rel, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode.Perm())
	if err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = root.Remove(rel)
		}
	}()
	if _, err = file.Write(contents); err != nil {
		file.Close()
		return err
	}
	if err = file.Chmod(mode); err != nil {
		file.Close()
		return err
	}
	if err = file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	keep = true
	return nil
}

func replaceRootedFile(repositoryRoot, relative string, contents []byte, mode os.FileMode) error {
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	rel := filepath.FromSlash(relative)
	before, err := root.Lstat(rel)
	if err != nil || !before.Mode().IsRegular() {
		return fmt.Errorf("unsafe replacement target %s", relative)
	}
	old, err := root.ReadFile(rel)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(old)
	suffix := make([]byte, 8)
	if _, err = rand.Read(suffix); err != nil {
		return err
	}
	temp := filepath.Join(filepath.Dir(rel), "."+filepath.Base(rel)+".tmp-"+hex.EncodeToString(suffix))
	if err = createRootedFileWithRoot(root, temp, contents, mode); err != nil {
		return err
	}
	defer root.Remove(temp)
	current, err := root.Lstat(rel)
	if err != nil || !os.SameFile(before, current) {
		return fmt.Errorf("target %s changed before publication", relative)
	}
	currentBytes, err := root.ReadFile(rel)
	if err != nil || sha256.Sum256(currentBytes) != digest {
		return fmt.Errorf("target %s contents changed before publication", relative)
	}
	return root.Rename(temp, rel)
}

func createRootedFileWithRoot(root *os.Root, relative string, contents []byte, mode os.FileMode) error {
	file, err := root.OpenFile(relative, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode.Perm())
	if err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = root.Remove(relative)
		}
	}()
	if _, err = file.Write(contents); err != nil {
		file.Close()
		return err
	}
	if err = file.Chmod(mode); err != nil {
		file.Close()
		return err
	}
	if err = file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	keep = true
	return nil
}

func createTransactionDirectories(repositoryRoot string, recorder *gittransaction.MutationRecorder, directories []string) error {
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	sort.Slice(directories, func(i, j int) bool { return strings.Count(directories[i], "/") < strings.Count(directories[j], "/") })
	for _, directory := range directories {
		if err = root.Mkdir(filepath.FromSlash(directory), 0o755); err != nil {
			return err
		}
		if err = recorder.RecordCreatedDirectory(directory); err != nil {
			return err
		}
	}
	return nil
}
func absentDirectories(root string, directories []string) []string {
	result := []string{}
	for _, directory := range sharedprimitives.UniqueSortedStrings(directories) {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(directory))); os.IsNotExist(err) {
			result = append(result, directory)
		}
	}
	return result
}
func plannedChanges(paths []string) []resultmodel.RecordedChange {
	changes := make([]resultmodel.RecordedChange, 0, len(paths))
	for _, path := range sharedprimitives.UniqueSortedStrings(paths) {
		changes = append(changes, resultmodel.RecordedChange{Path: path, Kind: "planned", Detail: "would publish deterministic interview bytes"})
	}
	return changes
}

func interviewFinding(command, code string, severity resultmodel.FindingSeverity, path, evidence string, fixability resultmodel.FindingFixability, stop string) resultmodel.CommandFinding {
	finding := knowledgeFinding(command, code, severity, nil, evidence, fixability, stop)
	if path != "" {
		finding.AffectedPaths = []string{path}
		finding.NextArgv = []string{"do-work-cli", command}
		finding.VerificationArgv = []string{"do-work-cli", "--format", "json", command}
	}
	return finding
}

func retargetInterviewResult(result resultmodel.CommandResult, options interviewOptions) resultmodel.CommandResult {
	for index := range result.Findings {
		finding := &result.Findings[index]
		command := ""
		if len(finding.NextArgv) > 1 {
			command = finding.NextArgv[1]
			finding.NextJustRecipe = command
			if options.template != "" {
				finding.NextJustRecipe += " " + quoteRecipeArgument(options.template)
			}
		}
		if options.knowledgeRoot != "" {
			finding.NextArgv = appendFindingOption(finding.NextArgv, "--knowledge-root", options.knowledgeRoot)
			finding.VerificationArgv = appendFindingOption(finding.VerificationArgv, "--knowledge-root", options.knowledgeRoot)
			finding.NextJustRecipe += " --knowledge-root " + quoteRecipeArgument(options.knowledgeRoot)
		}
		if options.template != "" {
			finding.NextArgv = appendFindingOption(finding.NextArgv, "--template", options.template)
			finding.VerificationArgv = appendFindingOption(finding.VerificationArgv, "--template", options.template)
		}
		if options.kb != "" && options.kb != "kb" && len(finding.NextArgv) > 1 && finding.NextArgv[1] == CommandInterviewIngest {
			finding.NextArgv = appendFindingOption(finding.NextArgv, "--kb", options.kb)
			finding.VerificationArgv = appendFindingOption(finding.VerificationArgv, "--kb", options.kb)
			finding.NextJustRecipe += " --kb " + quoteRecipeArgument(options.kb)
		}
	}
	return result
}

func appendFindingOption(argv []string, name, value string) []string {
	for index, item := range argv {
		if item == name || strings.HasPrefix(item, name+"=") {
			if item == name && index+1 < len(argv) {
				argv[index+1] = value
			}
			return argv
		}
	}
	return append(argv, name, value)
}
func interviewFailure(command, code, path string, err error) resultmodel.CommandResult {
	return interviewFindingResult(command, code, resultmodel.SeverityError, path, err.Error(), resultmodel.OutcomeFailure)
}
func interviewFindingResult(command, code string, severity resultmodel.FindingSeverity, path, evidence string, outcome resultmodel.CommandOutcome) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: outcome, Findings: []resultmodel.CommandFinding{interviewFinding(command, code, severity, path, evidence, resultmodel.FixabilityManual, "the deterministic interview operation cannot continue")}}
}
func sessionPath(template string) string {
	return filepath.ToSlash(filepath.Join("do-work/interview", template, "session.json"))
}
func sessionCounts(session map[string]any) (int, int) {
	approved, entries := 0, 0
	for _, raw := range mapValue(session["layers"]) {
		layer := mapValue(raw)
		if value, ok := layer["approved"].(bool); ok && value {
			approved++
		}
		entries += len(sliceValue(layer["entries"]))
	}
	return approved, entries
}
func interviewLayerEntries(session map[string]any, layer string) []any {
	return sliceValue(mapValue(mapValue(session["layers"])[layer])["entries"])
}
func allInterviewEntries(session map[string]any) []any {
	result := []any{}
	for layer, raw := range mapValue(session["layers"]) {
		for _, entryRaw := range sliceValue(mapValue(raw)["entries"]) {
			entry := mapValue(entryRaw)
			entry["layer_id"] = layer
			result = append(result, entry)
		}
	}
	return result
}
func mapValue(value any) map[string]any {
	if result, ok := value.(map[string]any); ok {
		return result
	}
	return map[string]any{}
}
func sliceValue(value any) []any {
	if result, ok := value.([]any); ok {
		return result
	}
	return []any{}
}
func stringValue(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	if number, ok := value.(json.Number); ok {
		return number.String()
	}
	return fmt.Sprint(value)
}
func intValue(value any) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case json.Number:
		n, _ := typed.Int64()
		return int(n)
	default:
		n, _ := strconv.Atoi(stringValue(value))
		return n
	}
}
func resolveScoped(root, current map[string]any, path string) any {
	if value, found := resolveValueFound(current, path); found {
		return value
	}
	return resolveValue(root, path)
}
func resolveScopedFound(root, current map[string]any, path string) (any, bool) {
	if value, found := resolveValueFound(current, path); found {
		return value, true
	}
	return resolveValueFound(root, path)
}
func resolveValue(value any, path string) any {
	resolved, _ := resolveValueFound(value, path)
	return resolved
}
func resolveValueFound(value any, path string) (any, bool) {
	current := value
	for _, part := range strings.Split(path, ".") {
		if part == "this" {
			continue
		}
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		var exists bool
		current, exists = object[part]
		if !exists {
			return nil, false
		}
	}
	return current, true
}
func displayTemplateValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case bool:
		return strconv.FormatBool(typed)
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			parts = append(parts, displayTemplateValue(item))
		}
		return strings.Join(parts, ", ")
	case map[string]any:
		bytes, _ := json.Marshal(typed)
		return string(bytes)
	default:
		return fmt.Sprint(typed)
	}
}
func semverMajor(value string) (int, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return 0, errors.New("invalid semver")
	}
	major, err := strconv.Atoi(parts[0])
	return major, err
}
func firstLayer(template interviewTemplate) any {
	if len(template.Layers) == 0 {
		return nil
	}
	return template.Layers[0].ID
}
func randomID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("session-%d", nowUTC().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
func nextInterviewVersion(root string) (int, error) {
	entries, readError := os.ReadDir(root)
	if os.IsNotExist(readError) {
		return 1, nil
	}
	if readError != nil {
		return 0, readError
	}
	maximum := 0
	pattern := regexp.MustCompile(`^v([0-9]+)-`)
	seen := map[int]string{}
	for _, entry := range entries {
		match := pattern.FindStringSubmatch(entry.Name())
		if match != nil {
			number, _ := strconv.Atoi(match[1])
			if previous, duplicate := seen[number]; duplicate {
				return 0, fmt.Errorf("duplicate numeric version v%d: %s and %s", number, previous, entry.Name())
			}
			seen[number] = entry.Name()
			if number > maximum {
				maximum = number
			}
		}
	}
	return maximum + 1, nil
}
func collisionName(inbox, capture, name string, stamp time.Time, reserved map[string]bool) (string, error) {
	prefix := stamp.Format("150405")
	for attempt := 0; ; attempt++ {
		candidate := name
		if attempt == 1 {
			candidate = prefix + "-" + name
		} else if attempt > 1 {
			candidate = fmt.Sprintf("%s-%d-%s", prefix, attempt, name)
		}
		if reserved[candidate] {
			continue
		}
		taken, err := collisionTargetExists(inbox, capture, candidate)
		if err != nil {
			return "", err
		}
		if !taken {
			return candidate, nil
		}
	}
}
func collisionTargetExists(inbox, capture, name string) (bool, error) {
	if _, err := os.Lstat(filepath.Join(inbox, name)); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	found := false
	err := filepath.WalkDir(capture, func(path string, entry os.DirEntry, walkError error) error {
		if os.IsNotExist(walkError) && path == capture {
			return filepath.SkipDir
		}
		if walkError != nil {
			return walkError
		}
		if !entry.IsDir() && entry.Name() == name {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found, err
}
func validInboxQueue(data []byte) bool {
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.Contains(normalized, "# Inbox Queue") {
		return false
	}
	separator := regexp.MustCompile(`(?m)^\|(?:\s*:?-{3,}:?\s*\|){3,}\s*$`)
	return separator.MatchString(normalized)
}
func pathInside(root, target string) bool {
	relative, err := filepath.Rel(root, target)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}
func uniqueAnyStrings(values []any) []any {
	seen := map[string]bool{}
	stringsList := []string{}
	for _, value := range values {
		text := stringValue(value)
		if text != "" && !seen[text] {
			seen[text] = true
			stringsList = append(stringsList, text)
		}
	}
	sort.Strings(stringsList)
	result := make([]any, len(stringsList))
	for index, value := range stringsList {
		result[index] = value
	}
	return result
}
