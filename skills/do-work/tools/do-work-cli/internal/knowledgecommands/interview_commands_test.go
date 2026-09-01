package knowledgecommands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func writeInterviewFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const tinyInterviewTemplate = `---
name: tiny
description: Tiny deterministic interview.
version: 1.0.0
topic_cluster: tiny-topic
layers:
  - id: facts
    title: Facts
    order: 1
exports:
  - path: USER.md
    kind: narrative
---

# Tiny

## Export Templates

### ` + "`USER.md`" + ` — profile

` + "```markdown" + `
# {{session.session_id}}
{{#each facts.entries}}
- {{title}}: {{summary}}
{{/each}}
` + "```" + `
`

const tinySession = `{
  "template":"tiny","template_version":"1.0.0","session_id":"session-1",
  "started_at":"2026-08-01T10:00:00Z","last_activity_at":"2026-08-02T10:00:00Z",
  "status":"complete","pending_layer":null,"previous_version":null,
  "review_completed_at":"2026-08-02T11:00:00Z","review_runs":1,"last_exported_at":null,
  "layers":{"facts":{"approved":true,"entries":[{"title":"One","summary":"First fact","source_confidence":"confirmed"}]}}
}`

func TestInterviewListStatusAndVersionsAreDeterministicAndReadOnly(t *testing.T) {
	root := t.TempDir()
	knowledgeRoot := filepath.Join(root, "knowledge")
	writeInterviewFixture(t, knowledgeRoot, "interviews/tiny.md", tinyInterviewTemplate)
	writeInterviewFixture(t, root, "do-work/interview/tiny/session.json", tinySession)
	writeInterviewFixture(t, root, "do-work/interview/tiny/versions/v10-2026-08-20/session.json", tinySession)
	writeInterviewFixture(t, root, "do-work/interview/tiny/versions/v2-2026-08-10/session.json", tinySession)

	listed := handleInterviewList(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--knowledge-root", knowledgeRoot})
	if listed.Outcome != resultmodel.OutcomeSuccess || len(listed.Findings) != 1 || !strings.Contains(strings.Join(listed.Findings[0].Evidence, " "), "tiny") {
		t.Fatalf("list = %#v", listed)
	}
	if !strings.Contains(strings.Join(listed.Findings[0].NextArgv, "\x00"), knowledgeRoot) {
		t.Fatalf("list lost absolute knowledge-root spelling: %#v", listed.Findings[0].NextArgv)
	}
	before, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	status := handleInterviewStatus(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--knowledge-root", knowledgeRoot, "--template", "tiny"})
	after, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	if status.Outcome != resultmodel.OutcomeSuccess || string(before) != string(after) || !strings.Contains(strings.Join(status.Findings[0].Evidence, " "), "approved=1/1") {
		t.Fatalf("status = %#v changed=%v", status, string(before) != string(after))
	}
	if !strings.Contains(strings.Join(status.Findings[0].NextArgv, "\x00"), "--template\x00tiny") {
		t.Fatalf("status lost template in next argv: %#v", status.Findings[0].NextArgv)
	}
	versions := handleInterviewVersions(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--template", "tiny"})
	if len(versions.Findings) != 2 || !strings.Contains(versions.Findings[0].Evidence[0], "v2-") || !strings.Contains(versions.Findings[1].Evidence[0], "v10-") {
		t.Fatalf("versions = %#v", versions)
	}
}

func TestInterviewExportRendersDeclaredTemplateAndStampsAfterPublication(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init")
	runGitFixture(t, root, "config", "user.email", "fixture@example.com")
	runGitFixture(t, root, "config", "user.name", "Fixture")
	knowledgeRoot := filepath.Join(root, "knowledge")
	writeInterviewFixture(t, knowledgeRoot, "interviews/tiny.md", tinyInterviewTemplate)
	writeInterviewFixture(t, root, "do-work/interview/tiny/session.json", tinySession)
	writeInterviewFixture(t, root, "do-work/interview/tiny/CHANGELOG.md", "# Changes\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-m", "fixture")

	result := handleInterviewExport(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--knowledge-root", knowledgeRoot, "--template", "tiny"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("export = %#v", result)
	}
	exportBytes, err := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/exports/USER.md"))
	if err != nil || !strings.Contains(string(exportBytes), "# session-1") || !strings.Contains(string(exportBytes), "- One: First fact") {
		t.Fatalf("export bytes=%q err=%v", exportBytes, err)
	}
	var session map[string]any
	data, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	if err := json.Unmarshal(data, &session); err != nil || session["last_exported_at"] == nil || session["last_exported_at"] == "" {
		t.Fatalf("session stamp = %#v err=%v", session, err)
	}
}

func TestInterviewResetRequiresConfirmationWithoutWriting(t *testing.T) {
	root := t.TempDir()
	writeInterviewFixture(t, root, "do-work/interview/tiny/session.json", tinySession)
	before, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	result := handleInterviewReset(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--template", "tiny"})
	after, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	if result.Outcome != resultmodel.OutcomeRefused || string(before) != string(after) {
		t.Fatalf("reset = %#v changed=%v", result, string(before) != string(after))
	}
}

func TestShippedInterviewTemplateDeclaresAndRendersFiveExports(t *testing.T) {
	templatePath, err := filepath.Abs("../../../../../do-work-knowledge/interviews/work-operating-model.md")
	if err != nil {
		t.Fatal(err)
	}
	template, err := loadInterviewTemplate(templatePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(template.Exports) != 5 || len(template.Layers) != 5 {
		t.Fatalf("shipped template layers=%d exports=%d", len(template.Layers), len(template.Exports))
	}
	session := map[string]any{"session_id": "fixture", "role_or_name_or_repo": "fixture-repo", "last_exported_at": "2026-09-01T12:00:00Z", "previous_version": "none", "layers": map[string]any{}}
	data := map[string]any{
		"session":    session,
		"template":   map[string]any{"name": template.Name, "version": template.Version, "topic_cluster": template.TopicCluster},
		"derived":    deriveInterviewValues(session),
		"all_layers": map[string]any{"entries": []any{}},
	}
	for _, layer := range template.Layers {
		data[layer.ID] = map[string]any{"entries": []any{}}
	}
	renders, err := renderInterviewExports(template, data)
	if err != nil {
		t.Fatal(err)
	}
	if len(renders) != 5 {
		t.Fatalf("rendered exports = %#v", mapKeysBytes(renders))
	}
}

func TestShippedRendererAcceptsNullAndJSONEscapesEveryStringSubstitution(t *testing.T) {
	templatePath, err := filepath.Abs("../../../../../do-work-knowledge/interviews/work-operating-model.md")
	if err != nil {
		t.Fatal(err)
	}
	template, err := loadInterviewTemplate(templatePath)
	if err != nil {
		t.Fatal(err)
	}
	session := map[string]any{"session_id": "fixture\"\\\n", "role_or_name_or_repo": "fixture", "last_exported_at": "2026-09-01T12:00:00Z", "previous_version": nil, "layers": map[string]any{}}
	derived := deriveInterviewValues(session)
	derived["time_blocks"] = []any{map[string]any{"label": "focus \"quoted\" \\ block\nnext", "days": []any{"Mon"}, "start": "09:00", "end": "10:00", "type": "deep_work", "source_entries": []any{"rhythm.one"}}}
	data := map[string]any{"session": session, "template": map[string]any{"name": template.Name, "version": template.Version, "topic_cluster": template.TopicCluster}, "derived": derived, "all_layers": map[string]any{"entries": []any{}}}
	for _, layer := range template.Layers {
		data[layer.ID] = map[string]any{"entries": []any{}}
	}
	renders, err := renderInterviewExports(template, data)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"operating-model.json", "schedule-recommendations.json"} {
		var decoded any
		if err := json.Unmarshal(renders[name], &decoded); err != nil {
			t.Fatalf("%s is invalid JSON: %v\n%s", name, err, renders[name])
		}
	}
	var operating map[string]any
	_ = json.Unmarshal(renders["operating-model.json"], &operating)
	if value, exists := operating["previous_version"]; !exists || value != nil {
		t.Fatalf("previous_version = %#v, want JSON null", value)
	}
}

func TestInterviewResetArchivesCompleteStateAndSeedsVersionBacklink(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init")
	runGitFixture(t, root, "config", "user.email", "fixture@example.com")
	runGitFixture(t, root, "config", "user.name", "Fixture")
	knowledgeRoot := filepath.Join(root, "knowledge")
	writeInterviewFixture(t, knowledgeRoot, "interviews/tiny.md", tinyInterviewTemplate)
	writeInterviewFixture(t, root, "do-work/interview/tiny/session.json", tinySession)
	writeInterviewFixture(t, root, "do-work/interview/tiny/checkpoints/nested/one.md", "checkpoint\n")
	writeInterviewFixture(t, root, "do-work/interview/tiny/exports/USER.md", "export\n")
	writeInterviewFixture(t, root, "do-work/interview/tiny/CHANGELOG.md", "# Changes\n")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-m", "fixture")
	result := handleInterviewReset(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--knowledge-root", knowledgeRoot, "--template", "tiny", "--confirm"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("reset = %#v", result)
	}
	for _, archived := range []string{"session.json", "checkpoints/nested/one.md", "exports/USER.md"} {
		if _, err := os.Stat(filepath.Join(root, "do-work/interview/tiny/versions/v1-"+nowUTC().Format("2006-01-02"), filepath.FromSlash(archived))); err != nil {
			t.Fatalf("archive %s: %v", archived, err)
		}
	}
	for _, removed := range []string{"do-work/interview/tiny/checkpoints/nested/one.md", "do-work/interview/tiny/exports/USER.md"} {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(removed))); !os.IsNotExist(err) {
			t.Fatalf("working state not cleared %s: %v", removed, err)
		}
	}
	data, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	var fresh map[string]any
	_ = json.Unmarshal(data, &fresh)
	if fresh["previous_version"] != "v1" {
		t.Fatalf("previous_version = %#v", fresh["previous_version"])
	}
}

func TestInterviewResetPropagatesBrokenEntryAndDoesNotPartiallyReset(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init")
	knowledgeRoot := filepath.Join(root, "knowledge")
	writeInterviewFixture(t, knowledgeRoot, "interviews/tiny.md", tinyInterviewTemplate)
	writeInterviewFixture(t, root, "do-work/interview/tiny/session.json", tinySession)
	if err := os.MkdirAll(filepath.Join(root, "do-work/interview/tiny/checkpoints"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "missing-outside"), filepath.Join(root, "do-work/interview/tiny/checkpoints/broken.md")); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	result := handleInterviewReset(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--knowledge-root", knowledgeRoot, "--template", "tiny", "--confirm"})
	after, _ := os.ReadFile(filepath.Join(root, "do-work/interview/tiny/session.json"))
	if result.Outcome != resultmodel.OutcomeFailure || string(before) != string(after) {
		t.Fatalf("broken reset = %#v changed=%v", result, string(before) != string(after))
	}
}

func TestInterviewVersionsRejectDuplicateNumericIDs(t *testing.T) {
	root := t.TempDir()
	writeInterviewFixture(t, root, "do-work/interview/tiny/versions/v2-2026-08-01/session.json", tinySession)
	writeInterviewFixture(t, root, "do-work/interview/tiny/versions/v2-2026-08-02/session.json", tinySession)
	result := handleInterviewVersions(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--template", "tiny"})
	if result.Outcome != resultmodel.OutcomeFindings || !resultHasFindingCode(result, "INTERVIEW-VERSION-DUPLICATE-ID") {
		t.Fatalf("duplicate versions = %#v", result)
	}
}

func TestInterviewIngestRefusesMalformedQueueWithoutWriting(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init")
	knowledgeRoot := filepath.Join(root, "knowledge")
	writeInterviewFixture(t, knowledgeRoot, "interviews/tiny.md", tinyInterviewTemplate)
	writeInterviewFixture(t, root, "do-work/interview/tiny/session.json", strings.Replace(tinySession, `"last_exported_at":null`, `"last_exported_at":"2026-08-02T12:00:00Z"`, 1))
	writeInterviewFixture(t, root, "do-work/interview/tiny/exports/USER.md", "# Export\n")
	writeInterviewFixture(t, root, "kb/raw/_inbox_queue.md", "not a queue\n")
	writeInterviewFixture(t, root, "kb/wiki/.gitkeep", "")
	result := handleInterviewIngest(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--knowledge-root", knowledgeRoot, "--template", "tiny", "--kb", "kb"})
	if result.Outcome != resultmodel.OutcomeRefused || !resultHasFindingCode(result, "INTERVIEW-INGEST-QUEUE-MALFORMED") {
		t.Fatalf("malformed queue = %#v", result)
	}
	entries, _ := os.ReadDir(filepath.Join(root, "kb/raw/inbox"))
	if len(entries) != 0 {
		t.Fatalf("ingest wrote files: %#v", entries)
	}
}

func TestCollisionNameRetriesOccupiedSameSecondPrefix(t *testing.T) {
	root := t.TempDir()
	inbox, capture := filepath.Join(root, "inbox"), filepath.Join(root, "capture")
	writeInterviewFixture(t, inbox, "tiny-USER.md", "one")
	writeInterviewFixture(t, inbox, "120000-tiny-USER.md", "two")
	writeInterviewFixture(t, capture, "nested/120000-2-tiny-USER.md", "three")
	name, err := collisionName(inbox, capture, "tiny-USER.md", time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC), nil)
	if err != nil || name != "120000-3-tiny-USER.md" {
		t.Fatalf("collision name=%q err=%v", name, err)
	}
}

func TestInterviewSingletonOptionsRejectRepetition(t *testing.T) {
	for _, arguments := range [][]string{
		{"--knowledge-root", "one", "--knowledge-root", "two"},
		{"--template", "tiny", "--template", "tiny"},
		{"--kb", "one", "--kb", "two"},
		{"--dry-run", "--dry-run"},
		{"--commit", "--commit"},
		{"--confirm", "--confirm"},
	} {
		if _, err := parseInterviewOptions(arguments, true); err == nil {
			t.Fatalf("repeated singleton accepted: %#v", arguments)
		}
	}
}

func resultHasFindingCode(result resultmodel.CommandResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}
