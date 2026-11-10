package toolboxcommands

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const sqliteFixture = "SQLite format 3\x00 rest of a database page"

func runFolderStructureCheck(t *testing.T, repository string, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	handler := Handlers()[CommandFolderStructureCheck]
	if handler == nil {
		t.Fatal("folder-structure-check is not registered")
	}
	return handler(commandruntime.ExecutionContext{RepositoryRoot: repository}, arguments)
}

func folderStructureGapKeys(result resultmodel.CommandResult) []string {
	keys := []string{}
	for _, finding := range result.Findings {
		if finding.Severity == resultmodel.SeverityWarning {
			keys = append(keys, finding.Code+" "+finding.AffectedPaths[0])
		}
	}
	sort.Strings(keys)
	return keys
}

// A layout that follows the target passes with exit 0, so a project can wire the check into CI
// without a false red on day one. A tracked fixture database is not a gap, and neither is an
// ignored database in node_modules/ or in do-work's own evidence (seen in g1w).
func TestFolderStructureCheckPassesTheTargetLayout(t *testing.T) {
	repository := toolboxTestRepository(t)
	writeOrganizeFiles(t, repository, map[string]string{
		".gitignore":           "data/*\n!data/README.md\nnode_modules/\ndo-work/deliverables/\n",
		"README.md":            "readme",
		"CHANGELOG.md":         "changes",
		"justfile":             "port := \"8080\"\nurl:=\"http://localhost\"\nset shell := [\"bash\", \"-c\"]\n\n# deploy to the server\ndeploy:\n    bash cmd/deploy.sh\n\n[positional-arguments]\nserve $port=\"8090\" *args:\n    bash cmd/serve.sh \"$@\"\n",
		"cmd/deploy.sh":        "echo deploy",
		"src/be-golang/go.mod": "module example",
		"src/be-golang/internal/store/fixture.db": sqliteFixture,
		"src/fe-node/package.json":                "{}",
		"src/mobile-react-native/App.tsx":         "app",
		"data/README.md":                          "be-golang/app.db",
		"data/be-golang/app.db":                   sqliteFixture,
		"node_modules/cache.db":                   sqliteFixture,
		"do-work/deliverables/evidence/cms.db":    sqliteFixture,
		"docs/spec/home.md":                       "spec",
	})
	toolboxTestGit(t, repository, "add", "-A")
	toolboxTestGit(t, repository, "commit", "-qm", "fixture")

	result := runFolderStructureCheck(t, repository)
	if result.Outcome != resultmodel.OutcomeSuccess || len(result.Findings) != 0 {
		t.Fatalf("target layout should pass, got %s %v", result.Outcome, folderStructureGapKeys(result))
	}
}

// Each gap class the action relies on: g1w-style root clutter (node-source/), a src/ folder that
// misses the <role>-<tech> form, a runtime database written inside the backend (gitignored, so
// only a content walk finds it), data/ committed to Git, and a justfile recipe holding logic.
// The run-kanban recipe pins REQ-713: its parameters carry default values, and the check used to
// stop at the "=" and skip the recipe, so its three body lines were never counted.
func TestFolderStructureCheckReportsEachGapClass(t *testing.T) {
	repository := toolboxTestRepository(t)
	writeOrganizeFiles(t, repository, map[string]string{
		".gitignore":                    "*.sqlite\n",
		"README.md":                     "readme",
		"skills-lock.json":              "{}",
		"node-source/package.json":      "{}",
		"src/backend-cms-golang/go.mod": "module example",
		"src/main.go":                   "package main",
		"src/be-golang/go.mod":          "module example",
		"src/be-golang/player.sqlite":   sqliteFixture,
		"data/seed.json":                "{}",
		"justfile":                      "test:\n    cd src/be-golang\n    go test ./...\n\nlint:\n    bash cmd/lint.sh\n\n# default values, one holding a colon, must not hide the recipe (run-kanban in the suite's own justfile)\nrun-kanban $port=\"8090\" url='http://localhost' +rest: lint\n    echo \"$port\"\n    echo \"$url\"\n    bash cmd/serve.sh\n",
		"docs/inbox/idea.md":            "idea",
	})
	toolboxTestGit(t, repository, "add", "-A")
	toolboxTestGit(t, repository, "commit", "-qm", "fixture")

	result := runFolderStructureCheck(t, repository)
	want := []string{
		"FOLDER-STRUCTURE-DATA-NOT-IGNORED data/seed.json",
		"FOLDER-STRUCTURE-DATA-OUTSIDE src/be-golang/player.sqlite",
		"FOLDER-STRUCTURE-JUST-LOGIC justfile: run-kanban",
		"FOLDER-STRUCTURE-JUST-LOGIC justfile: test",
		"FOLDER-STRUCTURE-ROOT-ENTRY node-source/",
		"FOLDER-STRUCTURE-ROOT-ENTRY skills-lock.json",
		"FOLDER-STRUCTURE-SRC-ENTRY src/backend-cms-golang",
		"FOLDER-STRUCTURE-SRC-ENTRY src/main.go",
	}
	if got := folderStructureGapKeys(result); !reflect.DeepEqual(got, want) {
		t.Fatalf("gaps\n got %v\nwant %v", got, want)
	}
	if result.Outcome != resultmodel.OutcomeFindings || resultmodel.ExitCode(result.Outcome) != 1 {
		t.Fatalf("gaps must exit 1 so CI can block, got %s", result.Outcome)
	}
	inboxReported := false
	for _, finding := range result.Findings {
		inboxReported = inboxReported || (finding.Code == "FOLDER-STRUCTURE-INBOX" && finding.Severity == resultmodel.SeverityInfo)
	}
	if !inboxReported {
		t.Fatal("docs/inbox/ files must be listed as an info finding")
	}

	// --allow lets a project keep a tool-mandated root file without editing the shipped layout.
	allowedResult := runFolderStructureCheck(t, repository, "--allow", "skills-lock.json", "--allow=node-source/")
	for _, key := range folderStructureGapKeys(allowedResult) {
		if key == "FOLDER-STRUCTURE-ROOT-ENTRY skills-lock.json" || key == "FOLDER-STRUCTURE-ROOT-ENTRY node-source/" {
			t.Fatalf("--allow did not clear %s", key)
		}
	}
}

// The suite installs this template as every consumer's managed justfile block, and the check only
// ever reads a repository's root justfile, so nothing else would see a consumer-side gap. Pins the
// failure where every consumer's folder-structure-check reported FOLDER-STRUCTURE-JUST-LOGIC for
// the suite's own run-kanban and kanban-static recipes.
func TestShippedBoardTemplateRecipesHoldNoLogic(t *testing.T) {
	shippedTemplatePath := filepath.Join("..", "..", "..", "..", "..", "do-work-board", "justfile.template")
	shippedTemplate, err := os.ReadFile(shippedTemplatePath)
	if err != nil {
		t.Fatalf("read shipped board template: %v", err)
	}
	consumerRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(consumerRoot, "justfile"), shippedTemplate, 0o644); err != nil {
		t.Fatalf("write consumer justfile: %v", err)
	}
	if gaps := justfileLogicGaps(consumerRoot); len(gaps) != 0 {
		t.Fatalf("shipped board template recipes hold logic: %v", gaps)
	}
}
