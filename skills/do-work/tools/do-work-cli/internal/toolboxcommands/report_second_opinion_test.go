package toolboxcommands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func secondOpinionFixture(t *testing.T) (string, string) {
	t.Helper()
	repository := toolboxTestRepository(t)
	bundle := filepath.Join(repository, "ai-reports", "2026-10-15_1200-example")
	if err := os.MkdirAll(bundle, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "index.html"), []byte("<!doctype html><html><body><h1>Original</h1></body></html>"), 0644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", ".")
	toolboxTestGit(t, repository, "commit", "-qm", "fixture")
	template := filepath.Join(t.TempDir(), "prompt.md")
	os.WriteFile(template, []byte("Output: {{OUTPUT_FILE}}\nPrefix: {{ASSET_PREFIX}}\nRepository: {{REPOSITORY_ROOT}}\nBundle: {{BUNDLE_PATH}}\n{{BACKLINK_REQUIREMENT}}\nDirection: {{EDITOR_DIRECTION}}\n"), 0644)
	t.Setenv("DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT", template)
	bin := t.TempDir()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(gitPath, filepath.Join(bin, "git")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	names, apps := aiReportJudgeBrowserNames, aiReportJudgeBrowserAppPaths
	aiReportJudgeBrowserNames = nil
	aiReportJudgeBrowserAppPaths = nil
	t.Cleanup(func() { aiReportJudgeBrowserNames = names; aiReportJudgeBrowserAppPaths = apps })
	return repository, bundle
}

func secondOpinionFakeCodex(t *testing.T, extra string) {
	t.Helper()
	script := `#!/bin/sh
printf 'OpenAI Codex test\n--------\nworkdir: %s\nmodel: actual-model-1\nprovider: openai\nsandbox: workspace-write\nsession id: fixture\n--------\n' "$PWD" >&2
for prompt do :; done
case "$prompt" in
 *'Output: '*)
 output=${prompt#*'Output: '}; output=${output%%'
'*}
 printf '<!doctype html><html><body><p>Original report: <a href="index.html">index.html</a></p><h1>Clear report</h1></body></html>' > "$output"
 ` + extra + `
 ;;
esac
printf 'File: %s/%s\n' "$PWD" "${output:-}"
printf 'Clarity changes\n- Clear flow\nCorrections: none\nOpen questions: none\n'
`
	if err := os.WriteFile(filepath.Join(os.Getenv("PATH"), "codex"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
}

func secondOpinionInvoke(t *testing.T, repository string, args ...string) resultmodel.CommandResult {
	t.Helper()
	handler := Handlers()["ai-report-second-opinion"]
	if handler == nil {
		t.Fatal("ai-report-second-opinion command is absent")
	}
	return handler(commandruntime.ExecutionContext{RepositoryRoot: repository}, args)
}

func secondOpinionRows(t *testing.T, result resultmodel.CommandResult) []map[string]any {
	t.Helper()
	raw, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	var payload struct {
		SecondOpinion struct {
			Runs []map[string]any `json:"runs"`
		} `json:"ai_report_second_opinion"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.SecondOpinion.Runs) == 0 {
		t.Fatalf("missing typed backend summary: %s", raw)
	}
	return payload.SecondOpinion.Runs
}

func TestSecondOpinionNativeModelCollisionAndSkippedJudge(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionFakeCodex(t, "")
	original, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	for _, suffix := range []string{"", "-2"} {
		result := secondOpinionInvoke(t, repo, "--backend", "codex")
		row := secondOpinionRows(t, result)[0]
		if row["actual_model"] != "actual-model-1" || row["judge_verdict"] != "skipped" {
			t.Fatalf("row=%v", row)
		}
		reply, _ := os.ReadFile(filepath.Join(bundle, "codex-actual-model-1"+suffix+"-reply.md"))
		if !strings.Contains(string(reply), "index-codex-actual-model-1"+suffix+".html") || strings.Contains(string(reply), "staging-") || strings.Contains(string(reply), bundle) {
			t.Fatalf("saved reply points to deleted stage: %s", reply)
		}
		for _, name := range []string{"index-codex-actual-model-1" + suffix + ".html", "codex-actual-model-1" + suffix + "-reply.md"} {
			if _, err := os.Stat(filepath.Join(bundle, name)); err != nil {
				t.Fatal(err)
			}
		}
	}
	after, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	if string(after) != string(original) {
		t.Fatal("skipped judge changed original")
	}
}

func TestSecondOpinionTargetNeverSubstitutes(t *testing.T) {
	repo, _ := secondOpinionFixture(t)
	secondOpinionFakeCodex(t, "")
	result := secondOpinionInvoke(t, repo, "--backend", "gemini")
	row := secondOpinionRows(t, result)[0]
	if row["backend"] != "gemini" || row["state"] != "skipped" || result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("result=%+v row=%v", result, row)
	}
	result = secondOpinionInvoke(t, repo, "--backend", "codex", "--model", "requested-model")
	row = secondOpinionRows(t, result)[0]
	if row["state"] != "failed" {
		t.Fatalf("model substitution accepted: %v", row)
	}
}

func TestSecondOpinionOutsideDirtAndOriginalMutationRefuse(t *testing.T) {
	for _, scenario := range []string{"preexisting", "created", "original"} {
		t.Run(scenario, func(t *testing.T) {
			repo, bundle := secondOpinionFixture(t)
			extra := ""
			if scenario == "preexisting" {
				os.WriteFile(filepath.Join(repo, "outside.txt"), []byte("dirty"), 0644)
			}
			if scenario == "created" {
				extra = "printf dirty > ../../outside.txt"
			}
			if scenario == "original" {
				extra = "printf changed > index.html"
			}
			secondOpinionFakeCodex(t, extra)
			result := secondOpinionInvoke(t, repo, "--backend", "codex")
			row := secondOpinionRows(t, result)[0]
			if row["state"] != "failed" {
				t.Fatalf("unsafe run accepted: %v", row)
			}
			original, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
			if strings.Contains(string(original), "second-opinion-banner:") {
				t.Fatal("failed run inserted banner")
			}
		})
	}
}

func TestSecondOpinionBatchIncludesMissingAndExcluded(t *testing.T) {
	repo, _ := secondOpinionFixture(t)
	secondOpinionFakeCodex(t, "")
	result := secondOpinionInvoke(t, repo, "--all", "--exclude", "agy")
	rows := secondOpinionRows(t, result)
	if len(rows) != 3 || rows[0]["backend"] != "codex" || rows[1]["backend"] != "gemini" || rows[1]["state"] != "skipped" || rows[2]["state"] != "excluded" {
		t.Fatalf("rows=%v", rows)
	}
}

func TestSecondOpinionBannerOnlyRealBodyAndPerFileIdempotency(t *testing.T) {
	page := "index-codex-model.html"
	original := []byte(`<!doctype html><!-- <body> --><html><head><title>Example <body> tag</title><script>const fake="<body><!-- second-opinion-banner:index-codex-model.html -->";</script></head><body class="real"><textarea><body></textarea><h1>Original</h1></body></html>`)
	updated, inserted, err := secondOpinionBanner(original, page)
	if err != nil || !inserted {
		t.Fatalf("inserted=%v err=%v", inserted, err)
	}
	expected := `<body class="real">` + "\n<!-- second-opinion-banner:" + page + " -->"
	if !strings.Contains(string(updated), expected) {
		t.Fatalf("banner missed actual body: %s", updated)
	}
	again, inserted, err := secondOpinionBanner(updated, page)
	if err != nil || inserted || string(again) != string(updated) {
		t.Fatal("same-filename banner not idempotent")
	}
	other, inserted, err := secondOpinionBanner(updated, "index-codex-model-2.html")
	if err != nil || !inserted || len(other) <= len(updated) {
		t.Fatal("different filename lost its banner")
	}
}

func TestSecondOpinionJudgeVerdictsGateOriginal(t *testing.T) {
	for _, verdict := range []string{"pass", "fail", "skipped", "error"} {
		t.Run(verdict, func(t *testing.T) {
			repo, bundle := secondOpinionFixture(t)
			secondOpinionFakeCodex(t, "")
			before, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
			previous := secondOpinionJudge
			secondOpinionJudge = func(string, []byte) (string, []resultmodel.AIReportSecondOpinionFinding) {
				return verdict, []resultmodel.AIReportSecondOpinionFinding{{FindingKind: "test-finding", FindingDetail: "observable evidence"}}
			}
			t.Cleanup(func() { secondOpinionJudge = previous })
			result := secondOpinionInvoke(t, repo, "--backend", "codex")
			row := secondOpinionRows(t, result)[0]
			after, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
			if row["judge_verdict"] != verdict {
				t.Fatalf("row=%v", row)
			}
			if verdict == "pass" {
				if row["banner_state"] != "inserted" || string(before) == string(after) || result.Outcome != resultmodel.OutcomeSuccess {
					t.Fatalf("pass did not link: %v", row)
				}
			} else if string(before) != string(after) || result.Outcome == resultmodel.OutcomeSuccess {
				t.Fatalf("non-pass changed original or claimed success: %v", row)
			}
		})
	}
}

func secondOpinionFakeJSON(t *testing.T, backend string, fail bool) {
	t.Helper()
	response := `printf '{"stats":{"models":{"native-1":{}}},"response":"Clarity: simpler. Corrections: none. Open questions: none."}\n'`
	if backend == "claude" {
		response = `printf '{"type":"result","subtype":"success","is_error":false,"modelUsage":{"native-1":{}},"result":"Clarity: simpler. Corrections: none. Open questions: none."}\n'`
	}
	if backend == "agy" {
		response = `printf '{"event":"init","conversation_id":"fake","init":{"model":"gemini-3.8-flash-low","cwd":"%s","permission_mode":"auto"}}\n' "$PWD"
printf '{"event":"result","result":{"conversation_id":"fake","status":"SUCCESS","response":"Clarity: simpler. Corrections: none. Open questions: none."}}\n'`
	}
	script := `#!/bin/sh
case "$1" in models) printf 'gemini-3.8-flash-low Gemini 3.8 Flash (Low)\n'; exit 0;; esac
prompt=
while [ "$#" -gt 0 ]; do
 case "$1" in -p) prompt=$2; shift 2;; *) shift;; esac
done
if [ "$prompt" = /model ]; then printf '{"status":"SUCCESS","command":{"name":"model","data":{"id":"gemini-3.8-flash-low","effort":"low"}}}\n'; exit 0; fi
case "$prompt" in *'Output: '*)
 output=${prompt#*'Output: '}; output=${output%%'
'*}
 printf '<html><body><p>Original report: <a href="index.html">index.html</a></p><h1>Clear</h1></body></html>' > "$output"
 ;; esac
` + response + "\n"
	if fail {
		script = "#!/bin/sh\nprintf '{\"error\":{\"message\":\"unavailable\"}}\\n'\nexit 3\n"
	}
	if err := os.WriteFile(filepath.Join(os.Getenv("PATH"), backend), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
}

func TestSecondOpinionAllNativeAdaptersAndSerialFailureContinuation(t *testing.T) {
	for _, backend := range []string{"gemini", "claude", "agy"} {
		t.Run(backend, func(t *testing.T) {
			repo, _ := secondOpinionFixture(t)
			secondOpinionFakeJSON(t, backend, false)
			result := secondOpinionInvoke(t, repo, "--backend", backend)
			row := secondOpinionRows(t, result)[0]
			if row["state"] != "ran" || row["actual_model"] == "" || !strings.Contains(row["page"].(string), "index-"+backend+"-") {
				t.Fatalf("row=%v", row)
			}
		})
	}
	t.Run("batch continues after failure", func(t *testing.T) {
		repo, _ := secondOpinionFixture(t)
		secondOpinionFakeCodex(t, "")
		secondOpinionFakeJSON(t, "gemini", true)
		secondOpinionFakeJSON(t, "agy", false)
		rows := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--all", "--exclude", "claude"))
		if len(rows) != 3 || rows[0]["state"] != "ran" || rows[1]["state"] != "failed" || rows[2]["state"] != "ran" {
			t.Fatalf("rows=%v", rows)
		}
	})
	t.Run("terminal selects once", func(t *testing.T) {
		repo, bundle := secondOpinionFixture(t)
		secondOpinionFakeCodex(t, "")
		secondOpinionFakeJSON(t, "gemini", false)
		os.WriteFile(filepath.Join(os.Getenv("PATH"), "codex"), []byte("#!/bin/sh\nexit 4\n"), 0755)
		rows := secondOpinionRows(t, secondOpinionInvoke(t, repo))
		files, _ := os.ReadDir(bundle)
		if len(rows) != 1 || rows[0]["backend"] != "codex" || rows[0]["state"] != "failed" || len(files) != 1 {
			t.Fatalf("fallback substituted after failure: %v files=%v", rows, files)
		}
	})
}

func TestSecondOpinionRejectsUnverifiedIdentityAndInvalidAGYCompletion(t *testing.T) {
	cases := []struct{ backend, stdout, stderr string }{
		{"codex", "model: invented", "user\nmodel: invented\nprovider: fake\nsession id: fake"},
		{"gemini", `{"response":"model: invented","stats":{"models":{}}}`, ""},
		{"claude", `{"type":"result","subtype":"success","modelUsage":{"one":{},"two":{}}}`, ""},
		{"agy", `{"event":"init","conversation_id":"x","init":{"model":"native","cwd":"/tmp","permission_mode":"auto"}}`, ""},
		{"agy", "{\"event\":\"result\",\"result\":{\"status\":\"SUCCESS\"}}", ""},
		{"agy", "{\"event\":\"init\",\"conversation_id\":\"x\",\"init\":{\"model\":\"native\",\"cwd\":\"/tmp\",\"permission_mode\":\"auto\"}}\n{\"event\":\"result\",\"result\":{\"conversation_id\":\"x\",\"status\":\"ERROR\"}}", ""},
	}
	for _, test := range cases {
		if _, _, _, err := secondOpinionMetadata(test.backend, "", secondOpinionCaptured{stdout: test.stdout, stderr: test.stderr}); err == nil {
			t.Fatalf("unverified identity accepted: %+v", test)
		}
	}
}

func TestSecondOpinionPreservesAssetsAndRejectsUnprefixedSymlinkAndMissingBacklink(t *testing.T) {
	for _, extra := range []string{"printf changed > asset.css", "printf new > unprefixed.css", "/bin/ln -s index.html codex-actual-model-1-link", "printf '<body><h1>Missing original link</h1></body>' > \"$output\""} {
		t.Run(extra, func(t *testing.T) {
			repo, bundle := secondOpinionFixture(t)
			os.WriteFile(filepath.Join(bundle, "asset.css"), []byte("old"), 0644)
			secondOpinionFakeCodex(t, extra)
			rows := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))
			if rows[0]["state"] != "failed" {
				t.Fatalf("unsafe backend accepted: %v", rows)
			}
			if strings.Contains(extra, "Missing original link") {
				// The backend finished the page; a backlink refusal keeps the bytes and names them.
				staged, _ := filepath.Glob(filepath.Join(bundle, ".codex-actual-model-1-staging-*.html"))
				if len(staged) != 1 {
					t.Fatalf("staged page not kept: %v", staged)
				}
				if info, err := os.Stat(staged[0]); err != nil || info.Size() == 0 {
					t.Fatalf("kept staged page is empty: %v %v", info, err)
				}
				if !strings.Contains(rows[0]["reason"].(string), "staged page kept at") {
					t.Fatalf("reason does not name the kept page: %v", rows[0]["reason"])
				}
			}
		})
	}
}

func TestSecondOpinionLatestAndLinkedBundleRefusal(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	newer := filepath.Join(repo, "ai-reports", "2026-10-16_1200-new")
	os.MkdirAll(newer, 0755)
	os.WriteFile(filepath.Join(newer, "index.html"), []byte("<body>new</body>"), 0644)
	chosen, err := secondOpinionBundle(repo, "latest")
	if err != nil || chosen != newer {
		t.Fatalf("latest=%s err=%v", chosen, err)
	}
	for _, selection := range []string{"../escape", repo, filepath.Join(bundle, "nested")} {
		if _, err := secondOpinionBundle(repo, selection); err == nil {
			t.Fatalf("invalid bundle accepted: %s", selection)
		}
	}
	linked := filepath.Join(repo, "ai-reports", "linked")
	os.Symlink(bundle, linked)
	if _, err := secondOpinionBundle(repo, linked); err == nil {
		t.Fatal("symlink bundle accepted")
	}
}

func TestSecondOpinionOutsideRenameIncludesBothSides(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	os.WriteFile(filepath.Join(repo, "outside.txt"), []byte("tracked"), 0644)
	toolboxTestGit(t, repo, "add", ".")
	toolboxTestGit(t, repo, "commit", "-qm", "outside fixture")
	toolboxTestGit(t, repo, "mv", "outside.txt", filepath.Join(bundle, "moved.txt"))
	outside, err := secondOpinionOutside(repo, bundle)
	if err != nil || len(outside) != 1 || outside[0] != "outside.txt" {
		t.Fatalf("rename source lost: %v %v", outside, err)
	}
}

func TestSecondOpinionCapturedRunnerCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	directory := t.TempDir()
	oldGrace := reportImageGracePeriod
	reportImageGracePeriod = 20 * time.Millisecond
	t.Cleanup(func() { reportImageGracePeriod = oldGrace })
	done := make(chan secondOpinionCaptured, 1)
	go func() {
		done <- secondOpinionCapture(ctx, secondOpinionProgress{}, directory, "/bin/sh", "-c", "trap '' TERM; (printf ready > ready; /bin/sleep 0.3; printf late > cancel-late; /bin/sleep 30) & wait")
	}()
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(directory, "ready")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			cancel()
			t.Fatal("backend did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case result := <-done:
		if !errors.Is(result.err, context.Canceled) {
			t.Fatalf("result=%+v", result)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("capture cancellation did not reap process group")
	}
	// Killing a waited-for sleep can resume its parent's next statement during
	// teardown. The boundary is no further writes after capture returns.
	before, beforeErr := os.ReadFile(filepath.Join(directory, "cancel-late"))
	t.Logf("at capture return: data=%q error=%v", before, beforeErr)
	time.Sleep(650 * time.Millisecond)
	after, afterErr := os.ReadFile(filepath.Join(directory, "cancel-late"))
	if string(before) != string(after) || os.IsNotExist(beforeErr) != os.IsNotExist(afterErr) {
		t.Fatal("cancelled descendant wrote after capture returned")
	}
}

func TestSecondOpinionBackendCommandSafetyContract(t *testing.T) {
	for _, backend := range []string{"codex", "gemini", "agy", "claude"} {
		args := secondOpinionArgv(backend, "native-model", "literal; $(touch unsafe)", false)
		joined := " " + strings.Join(args, " ") + " "
		if !strings.Contains(joined, "literal; $(touch unsafe)") {
			t.Fatal("prompt lost its literal argv value")
		}
		required := map[string][]string{"codex": {" exec -s workspace-write ", " -m native-model "}, "gemini": {" -p ", " --approval-mode auto_edit ", " -m native-model "}, "agy": {" --model native-model ", " --print-timeout 15m ", " --output-format stream-json ", " -p "}, "claude": {" -p ", " --model native-model "}}
		for _, flag := range required[backend] {
			if !strings.Contains(joined, flag) {
				t.Fatalf("%s lost required invocation %q: %q", backend, flag, args)
			}
		}
		if backend == "codex" && strings.Contains(joined, "--dangerously") {
			t.Fatal("Codex sandbox bypassed")
		}
	}
}

func TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionFakeCodex(t, "")
	backendPath := filepath.Join(os.Getenv("PATH"), "codex")
	script, _ := os.ReadFile(backendPath)
	// A native generation header differs from the native discovery header.
	changed := strings.Replace(string(script), "printf 'OpenAI Codex", "case \"$*\" in *'Output: '*) actual=different-model;; *) actual=actual-model-1;; esac\nprintf 'OpenAI Codex", 1)
	changed = strings.Replace(changed, "model: actual-model-1\\n", "model: %s\\n", 1)
	changed = strings.Replace(changed, "\"$PWD\" >&2", "\"$PWD\" \"$actual\" >&2", 1)
	os.WriteFile(backendPath, []byte(changed), 0755)
	rows := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))
	if rows[0]["state"] != "unverified" || rows[0]["actual_model"] != "" || fmt.Sprint(rows[0]["seen_models"]) != "[different-model]" {
		t.Fatalf("model switch not published as unverified: %v", rows)
	}
	files, _ := os.ReadDir(bundle)
	names := []string{}
	for _, file := range files {
		names = append(names, file.Name())
	}
	if fmt.Sprint(names) != "[codex-actual-model-1-unverified-reply.md index-codex-actual-model-1-unverified.html index.html]" {
		t.Fatalf("unverified pair missing or staging left behind: %v", names)
	}
	original, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	if strings.Contains(string(original), "second-opinion-banner:") {
		t.Fatal("unverified run inserted banner")
	}
}

// secondOpinionFakeClaude prints one envelope for the probe turn and another after writing the page.
func secondOpinionFakeClaude(t *testing.T, probeEnvelope, generationEnvelope string) {
	t.Helper()
	script := `#!/bin/sh
prompt=
while [ "$#" -gt 0 ]; do
 case "$1" in -p) prompt=$2; shift 2;; *) shift;; esac
done
case "$prompt" in *'Output: '*)
 output=${prompt#*'Output: '}; output=${output%%'
'*}
 printf '<html><body><p>Original report: <a href="index.html">index.html</a></p><h1>Clear</h1></body></html>' > "$output"
 printf '%s\n' '` + generationEnvelope + `'
 ;;
*) printf '%s\n' '` + probeEnvelope + `' ;;
esac
`
	if err := os.WriteFile(filepath.Join(os.Getenv("PATH"), "claude"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
}

func TestSecondOpinionSubagentModelIsNotAmbiguous(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	envelope := `{"type":"result","subtype":"success","is_error":false,"modelUsage":{"claude-fable-5-1":{},"claude-opus-5-5":{}},"result":"Clarity: simpler. Corrections: none. Open questions: none."}`
	secondOpinionFakeClaude(t, envelope, envelope)
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "claude", "--model", "claude-fable-5-1"))[0]
	if row["state"] != "ran" || row["actual_model"] != "claude-fable-5-1" || row["page"] != "ai-reports/2026-10-15_1200-example/index-claude-claude-fable-5-1.html" || fmt.Sprint(row["seen_models"]) != "[claude-fable-5-1 claude-opus-5-5]" {
		t.Fatalf("subagent model treated as ambiguity: %v", row)
	}
	if _, err := os.Stat(filepath.Join(bundle, "index-claude-claude-fable-5-1.html")); err != nil {
		t.Fatal(err)
	}
	files, _ := os.ReadDir(bundle)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			t.Fatalf("staging page left after publication: %s", file.Name())
		}
	}
}

func TestSecondOpinionUnverifiedModelWithoutPageFails(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	probe := `{"type":"result","subtype":"success","is_error":false,"modelUsage":{"claude-fable-5-1":{}},"result":"OK"}`
	generation := `{"type":"result","subtype":"success","is_error":false,"modelUsage":{"claude-opus-5-5":{},"claude-haiku-5-5":{}},"result":"No page was written."}`
	secondOpinionFakeClaude(t, probe, generation)
	// The backend answers with an unverifiable model and never writes the page: nothing finished exists to keep.
	backendPath := filepath.Join(os.Getenv("PATH"), "claude")
	script, _ := os.ReadFile(backendPath)
	changed := strings.Replace(string(script), ` printf '<html><body><p>Original report: <a href="index.html">index.html</a></p><h1>Clear</h1></body></html>' > "$output"`+"\n", "", 1)
	if changed == string(script) {
		t.Fatal("fake claude script did not drop its page write")
	}
	os.WriteFile(backendPath, []byte(changed), 0755)
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "claude", "--model", "claude-fable-5-1"))[0]
	if row["state"] != "failed" || !strings.Contains(row["reason"].(string), "backend wrote no page") {
		t.Fatalf("empty page not refused as a plain failure: %v", row)
	}
	files, _ := os.ReadDir(bundle)
	if len(files) != 1 || files[0].Name() != "index.html" {
		t.Fatalf("empty page published or staging left behind: %v", files)
	}
}

func TestSecondOpinionUnverifiedModelKeepsFinishedPage(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	probe := `{"type":"result","subtype":"success","is_error":false,"modelUsage":{"claude-fable-5-1":{}},"result":"OK"}`
	generation := `{"type":"result","subtype":"success","is_error":false,"modelUsage":{"claude-opus-5-5":{},"claude-haiku-5-5":{}},"result":"Clarity: simpler. Corrections: none. Open questions: none."}`
	secondOpinionFakeClaude(t, probe, generation)
	original, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	result := secondOpinionInvoke(t, repo, "--backend", "claude", "--model", "claude-fable-5-1")
	row := secondOpinionRows(t, result)[0]
	if row["state"] != "unverified" || row["actual_model"] != "" || fmt.Sprint(row["seen_models"]) != "[claude-haiku-5-5 claude-opus-5-5]" || row["banner_state"] != "not-inserted" || row["judge_verdict"] != "not-run" {
		t.Fatalf("unverified model did not keep the finished page: %v", row)
	}
	page, err := os.ReadFile(filepath.Join(bundle, "index-claude-claude-fable-5-1-unverified.html"))
	if err != nil || !strings.Contains(string(page), "<h1>Clear</h1>") {
		t.Fatalf("finished page bytes not kept: %q %v", page, err)
	}
	reply, err := os.ReadFile(filepath.Join(bundle, "claude-claude-fable-5-1-unverified-reply.md"))
	if err != nil || !strings.Contains(string(reply), "Clarity: simpler.") || strings.Contains(string(reply), "staging-") {
		t.Fatalf("reply not kept beside the page: %q %v", reply, err)
	}
	after, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	if string(after) != string(original) {
		t.Fatal("unverified run changed the original")
	}
	files, _ := os.ReadDir(bundle)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			t.Fatalf("staging page left after publication: %s", file.Name())
		}
	}
	if result.Outcome == resultmodel.OutcomeSuccess || len(result.Findings) != 1 || result.Findings[0].Code != "AI-REPORT-SECOND-OPINION-UNVERIFIED" {
		t.Fatalf("typed finding missing: outcome=%v findings=%+v", result.Outcome, result.Findings)
	}
	evidence := strings.Join(result.Findings[0].Evidence, "\n")
	for _, model := range []string{"claude-fable-5-1", "claude-opus-5-5", "claude-haiku-5-5"} {
		if !strings.Contains(evidence, model) {
			t.Fatalf("finding evidence omits %s: %s", model, evidence)
		}
	}
}

// secondOpinionFakeLogger installs a backend on the fixture PATH that only records its argv; the log file proves whether the launcher ever ran it.
func secondOpinionFakeLogger(t *testing.T, backend string) string {
	t.Helper()
	logPath := filepath.Join(t.TempDir(), backend+"-invocations.log")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" >> '" + logPath + "'\n"
	if err := os.WriteFile(filepath.Join(os.Getenv("PATH"), backend), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	return logPath
}

func TestSecondOpinionSessionPrepareNeverExecsOwnBackend(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	logPath := secondOpinionFakeLogger(t, "claude")
	result := secondOpinionInvoke(t, repo, "--prepare", "--backend", "claude", "--model", "claude-fable-5-1")
	row := secondOpinionRows(t, result)[0]
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Fatalf("the session's own backend was executed: %v", err)
	}
	if row["state"] != "prepared" || row["model_evidence"] != "session-declared" || row["asset_prefix"] != "claude-claude-fable-5-1-" || result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("row=%v outcome=%v", row, result.Outcome)
	}
	stage, _ := row["staging_page"].(string)
	if filepath.Dir(stage) != bundle || !strings.HasPrefix(filepath.Base(stage), ".") {
		t.Fatalf("staging page is not a dotfile inside the bundle: %q", stage)
	}
	if info, err := os.Stat(stage); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("staging page not reserved: %v", err)
	}
	stem := strings.TrimSuffix(stage, ".html")
	if _, err := os.Stat(stem + ".json"); err != nil {
		t.Fatalf("prepare record missing: %v", err)
	}
	if row["staging_reply"] != stem+"-reply.md" {
		t.Fatalf("reply stage not derived from the page stage: %v", row["staging_reply"])
	}
	prompt, _ := row["prompt"].(string)
	if !strings.Contains(prompt, "Output: "+stage+"\n") || !strings.Contains(prompt, "Original report:") || strings.Contains(prompt, "{{") {
		t.Fatalf("prompt not filled for the session: %q", prompt)
	}
}

func TestSecondOpinionSessionFinalizePublishesAgentWrittenPage(t *testing.T) {
	// prepare reserves the staging page for the session; the test then plays the agent and writes the page and reply.
	prepare := func(t *testing.T) (repo, bundle, stage, logPath string) {
		repo, bundle = secondOpinionFixture(t)
		logPath = secondOpinionFakeLogger(t, "claude")
		row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--prepare", "--backend", "claude", "--model", "claude-fable-5-1"))[0]
		if row["state"] != "prepared" {
			t.Fatalf("prepare failed: %v", row)
		}
		stage, _ = row["staging_page"].(string)
		return repo, bundle, stage, logPath
	}
	writePage := func(t *testing.T, bundle, stage string) {
		if err := os.WriteFile(stage, []byte(`<html><body><p>Original report: <a href="index.html">index.html</a></p><h1>Clear</h1><img src="claude-claude-fable-5-1-chart.svg"></body></html>`), 0644); err != nil {
			t.Fatal(err)
		}
		os.WriteFile(filepath.Join(bundle, "claude-claude-fable-5-1-chart.svg"), []byte("<svg/>"), 0644)
		os.WriteFile(strings.TrimSuffix(stage, ".html")+"-reply.md", []byte("Page: "+stage+" ("+filepath.Base(stage)+")\nClarity: simpler. Corrections: none. Open questions: none.\n"), 0644)
	}
	t.Run("publishes", func(t *testing.T) {
		repo, bundle, stage, logPath := prepare(t)
		writePage(t, bundle, stage)
		result := secondOpinionInvoke(t, repo, "--finalize", stage)
		row := secondOpinionRows(t, result)[0]
		if row["state"] != "ran" || row["backend"] != "claude" || row["model_evidence"] != "session-declared" || row["actual_model"] != "claude-fable-5-1" || row["judge_verdict"] != "skipped" || row["banner_state"] != "not-inserted" {
			t.Fatalf("agent-written page not published: %v", row)
		}
		if row["page"] != "ai-reports/2026-10-15_1200-example/index-claude-claude-fable-5-1.html" || row["reply"] != "ai-reports/2026-10-15_1200-example/claude-claude-fable-5-1-reply.md" {
			t.Fatalf("published names: %v", row)
		}
		reply, err := os.ReadFile(filepath.Join(bundle, "claude-claude-fable-5-1-reply.md"))
		if err != nil || !strings.Contains(string(reply), "index-claude-claude-fable-5-1.html") || strings.Contains(string(reply), "staging-") {
			t.Fatalf("reply not rewritten to the published page: %q %v", reply, err)
		}
		for _, name := range []string{stage, strings.TrimSuffix(stage, ".html") + ".json", strings.TrimSuffix(stage, ".html") + "-reply.md"} {
			if _, err := os.Stat(name); !os.IsNotExist(err) {
				t.Fatalf("staging leaf left behind: %s", name)
			}
		}
		if _, err := os.Stat(logPath); !os.IsNotExist(err) {
			t.Fatal("finalize executed the session's own backend")
		}
	})
	t.Run("judge pass inserts banner", func(t *testing.T) {
		repo, bundle, stage, _ := prepare(t)
		writePage(t, bundle, stage)
		before, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
		previous := secondOpinionJudge
		secondOpinionJudge = func(string, []byte) (string, []resultmodel.AIReportSecondOpinionFinding) {
			return "pass", []resultmodel.AIReportSecondOpinionFinding{}
		}
		t.Cleanup(func() { secondOpinionJudge = previous })
		result := secondOpinionInvoke(t, repo, "--finalize", stage)
		row := secondOpinionRows(t, result)[0]
		if row["state"] != "ran" || row["banner_state"] != "inserted" || result.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("pass did not link: %v", row)
		}
		expected, _, _ := secondOpinionBanner(before, "index-claude-claude-fable-5-1.html")
		after, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
		if string(after) != string(expected) {
			t.Fatalf("original changed beyond the banner: %s", after)
		}
	})
	t.Run("original modified keeps stage", func(t *testing.T) {
		repo, bundle, stage, _ := prepare(t)
		writePage(t, bundle, stage)
		os.WriteFile(filepath.Join(bundle, "index.html"), []byte("<html><body><h1>Edited</h1></body></html>"), 0644)
		row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--finalize", stage))[0]
		if row["state"] != "failed" || !strings.Contains(fmt.Sprint(row["reason"]), "staged page kept at") {
			t.Fatalf("modified original accepted: %v", row)
		}
		if info, err := os.Stat(stage); err != nil || info.Size() == 0 {
			t.Fatalf("finished staging page deleted: %v", err)
		}
		original, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
		if strings.Contains(string(original), "second-opinion-banner:") {
			t.Fatal("failed finalize inserted banner")
		}
	})
	t.Run("never written cleans up", func(t *testing.T) {
		repo, _, stage, _ := prepare(t)
		row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--finalize", stage))[0]
		if row["state"] != "failed" || !strings.Contains(fmt.Sprint(row["reason"]), "staging page was not written") {
			t.Fatalf("empty stage accepted: %v", row)
		}
		for _, name := range []string{stage, strings.TrimSuffix(stage, ".html") + ".json"} {
			if _, err := os.Stat(name); !os.IsNotExist(err) {
				t.Fatalf("never-written stage not cleaned up: %s", name)
			}
		}
	})
}

func TestSecondOpinionBacklinkAllowsSemanticAttributeAndInlineVariants(t *testing.T) {
	for _, markup := range []string{
		`<p>Original report: <a class="original" href="index.html">index.html</a></p>`,
		`<div><strong>Original report:</strong> <a title="Original" href='index.html'><code>index.html</code></a></div>`,
		`<p>Original report: <a data-note="href='other.html'" href=index.html>index.html</a></p>`,
	} {
		if err := secondOpinionBacklink([]byte(`<html><body>` + markup + `<h1>Report</h1></body></html>`)); err != nil {
			t.Errorf("valid opening %s: %v", markup, err)
		}
	}
	for _, markup := range []string{
		`<p>Original report: <a title="href='index.html'" href="other.html">index.html</a></p>`,
		`<p>Original report: <a href="other.html" href="index.html">index.html</a></p>`,
		`<h1>Report</h1><p>Original report: <a href="index.html">index.html</a></p>`,
		`<p>Original report: <a href="index.html">wrong.html</a></p>`,
		`<p>Original report: <a href="<span>index.html</span>">index.html</a></p>`,
	} {
		if err := secondOpinionBacklink([]byte(`<html><body>` + markup + `</body></html>`)); err == nil {
			t.Errorf("accepted invalid opening %s", markup)
		}
	}
}

func TestSecondOpinionCaptureDrainsExitedLeaderChildren(t *testing.T) {
	oldGrace := reportImageGracePeriod
	reportImageGracePeriod = 20 * time.Millisecond
	t.Cleanup(func() { reportImageGracePeriod = oldGrace })
	for _, redirect := range []string{` >/dev/null 2>&1`, ``} {
		t.Run(redirect, func(t *testing.T) {
			directory := t.TempDir()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			started := time.Now()
			captured := secondOpinionCapture(ctx, secondOpinionProgress{}, directory, "/bin/sh", "-c", `(/bin/sleep 0.3; printf late > late-write; /bin/sleep 0.6)`+redirect+` &`)
			if captured.err != nil {
				t.Fatalf("capture: %v", captured.err)
			}
			if elapsed := time.Since(started); elapsed > 600*time.Millisecond {
				t.Errorf("inherited pipes delayed capture: %v", elapsed)
			}
			time.Sleep(650 * time.Millisecond)
			if _, err := os.Stat(filepath.Join(directory, "late-write")); err == nil {
				t.Fatal("capture returned while a descendant could still write")
			}
		})
	}
}

// progressRecorder is a test progress sink that writes markerPath once it has seen every wanted fragment.
type progressRecorder struct {
	mutex      sync.Mutex
	text       strings.Builder
	wanted     []string
	markerPath string
}

func (recorder *progressRecorder) Write(data []byte) (int, error) {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	recorder.text.Write(data)
	for _, fragment := range recorder.wanted {
		if !strings.Contains(recorder.text.String(), fragment) {
			return len(data), nil
		}
	}
	return len(data), os.WriteFile(recorder.markerPath, nil, 0644)
}

func (recorder *progressRecorder) String() string {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	return recorder.text.String()
}

// Pins the silent pane: the launcher buffered all backend output until exit, so a working backend and a stalled
// one looked the same. The fake backend holds until the test has seen its line and a heartbeat on the sink.
func TestSecondOpinionCaptureStreamsProgressWhileBackendRuns(t *testing.T) {
	repo, _ := secondOpinionFixture(t)
	markerPath := filepath.Join(t.TempDir(), "progress-seen")
	recorder := &progressRecorder{wanted: []string{"[codex stderr] drafting\n", "codex generation", " bytes"}, markerPath: markerPath}
	previousSink, previousInterval := secondOpinionProgressSink, secondOpinionHeartbeatInterval
	secondOpinionProgressSink, secondOpinionHeartbeatInterval = recorder, 20*time.Millisecond
	t.Cleanup(func() { secondOpinionProgressSink, secondOpinionHeartbeatInterval = previousSink, previousInterval })
	// The page is written before this runs, so a generation heartbeat can report its size; the wait is bounded at 4s.
	secondOpinionFakeCodex(t, `printf 'drafting\n' >&2
 waited=0
 while [ ! -e '`+markerPath+`' ] && [ "$waited" -lt 200 ]; do /bin/sleep 0.02; waited=$((waited+1)); done
 printf 'tail without newline' >&2`)
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))[0]
	if row["actual_model"] != "actual-model-1" || row["page"] == nil {
		t.Fatalf("row=%v", row)
	}
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("backend exited before its line and a heartbeat reached the sink: %q", recorder.String())
	}
	progress := recorder.String()
	for _, fragment := range []string{"[codex stdout] Clarity changes\n", "[codex stderr] tail without newline\n"} {
		if !strings.Contains(progress, fragment) {
			t.Fatalf("progress lacks %q: %q", fragment, progress)
		}
	}
}

// Pins the dead cap: with bytes.Buffer embedded, io.Copy used its promoted ReadFrom, so a runaway backend was
// buffered without limit and without the truncation error (and the live mirror in Write never ran either).
func TestSecondOpinionCaptureCapsBackendOutput(t *testing.T) {
	previousSink := secondOpinionProgressSink
	secondOpinionProgressSink = io.Discard
	t.Cleanup(func() { secondOpinionProgressSink = previousSink })
	captured := secondOpinionCapture(context.Background(), secondOpinionProgress{}, t.TempDir(), "/usr/bin/head", "-c", "16777300", "/dev/zero")
	if len(captured.stdout) != 16*1024*1024 || captured.err == nil {
		t.Fatalf("captured %d bytes, err=%v", len(captured.stdout), captured.err)
	}
}

func TestSecondOpinionCapturePreservesBackendStatusAndChannels(t *testing.T) {
	for _, status := range []string{"0", "7"} {
		captured := secondOpinionCapture(context.Background(), secondOpinionProgress{}, t.TempDir(), "/bin/sh", "-c", `printf '%s' "$1"; printf warning >&2; if read value; then exit 99; fi; if (printf forged >&3) 2>/dev/null; then exit 98; fi; exit "$2"`, "backend", `literal $() ' argument`, status)
		if captured.stdout != `literal $() ' argument` || captured.stderr != "warning" {
			t.Fatalf("capture=%+v", captured)
		}
		if status == "0" && captured.err != nil || status == "7" && (captured.err == nil || !strings.Contains(captured.err.Error(), "7")) {
			t.Fatalf("status=%s capture=%+v", status, captured)
		}
	}
}

// secondOpinionFixturePresets writes bold.md and visual.md beside the fixture template, where the launcher looks for
// shipped presets, so tests never read the real shipped files.
func secondOpinionFixturePresets(t *testing.T) string {
	t.Helper()
	folder := filepath.Join(filepath.Dir(os.Getenv("DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT")), "second-opinion-presets")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}
	for name, text := range map[string]string{"bold": "Shipped bold: add motion.\n", "visual": "Shipped visual: diagrams first.\n"} {
		if err := os.WriteFile(filepath.Join(folder, name+".md"), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return folder
}

// Pins: a preset that never reaches the prompt, an unknown name that hides what exists, and direction text that would
// trip the placeholder check or ride along on --finalize.
func TestSecondOpinionPresetResolutionAndArgumentRules(t *testing.T) {
	repo, _ := secondOpinionFixture(t)
	shipped := secondOpinionFixturePresets(t)
	session := []string{"--prepare", "--backend", "claude", "--model", "claude-fable-5-1"}
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, append(session, "--preset", "bold")...))[0]
	prompt, _ := row["prompt"].(string)
	if row["state"] != "prepared" || row["preset"] != "bold" || row["preset_file"] != filepath.Join(shipped, "bold.md") || !strings.Contains(prompt, "Direction: Shipped bold: add motion.\n") || strings.Contains(prompt, "{{") {
		t.Fatalf("shipped preset not applied: %v", row)
	}
	refused := map[string][]string{
		`unknown preset "missing"; presets found: bold, visual`: append(session, "--preset", "missing"),
		"--editor-direction must not be empty":                  append(session, "--editor-direction", " "),
		"--editor-direction must not contain {{":                append(session, "--editor-direction", "keep {{OUTPUT_FILE}}"),
		"--finalize takes only the staging page path":           {"--finalize", "ai-reports/2026-10-15_1200-example/.x-staging-1.html", "--preset", "bold"},
	}
	for want, args := range refused {
		result := secondOpinionInvoke(t, repo, args...)
		if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) == 0 || !strings.Contains(strings.Join(result.Findings[0].Evidence, " "), want) {
			t.Errorf("%v: want refusal %q, got %+v", args, want, result.Findings)
		}
	}
}

// Pins: a preset kept in the report bundle that is not found, and a leftover lookup that still reads the shared
// ai-reports/second-opinion-presets/ folder outside the bundle.
func TestSecondOpinionPresetResolvesFromBundleOnly(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionFixturePresets(t)
	own := filepath.Join(bundle, "second-opinion-presets")
	os.MkdirAll(own, 0755)
	session := []string{"--prepare", "--backend", "claude", "--model", "claude-fable-5-1"}
	for _, preset := range [][2]string{{"mine", "Bundle mine: tell it as a story.\n"}, {"bold", "Bundle bold: go further.\n"}} {
		name, text := preset[0], preset[1]
		os.WriteFile(filepath.Join(own, name+".md"), []byte(text), 0644)
		row := secondOpinionRows(t, secondOpinionInvoke(t, repo, append(session, "--preset", name)...))[0]
		prompt, _ := row["prompt"].(string)
		if row["state"] != "prepared" || row["preset_file"] != "ai-reports/2026-10-15_1200-example/second-opinion-presets/"+name+".md" || !strings.Contains(prompt, "Direction: "+text) {
			t.Fatalf("uncommitted bundle preset %s not applied: %v", name, row)
		}
	}
	legacy := filepath.Join(repo, "ai-reports", "second-opinion-presets")
	os.MkdirAll(legacy, 0755)
	os.WriteFile(filepath.Join(legacy, "legacy.md"), []byte("Legacy: shared folder.\n"), 0644)
	toolboxTestGit(t, repo, "add", "ai-reports/second-opinion-presets")
	toolboxTestGit(t, repo, "commit", "-qm", "legacy preset")
	result := secondOpinionInvoke(t, repo, append(session, "--preset", "legacy")...)
	if want := `unknown preset "legacy"; presets found: bold, mine, visual`; result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) == 0 || !strings.Contains(strings.Join(result.Findings[0].Evidence, " "), want) {
		t.Fatalf("shared folder preset still read: want refusal %q, got %+v", want, result.Findings)
	}
}

// Pins: a direction that never reaches the editor, a directed page published under the plain name or without its
// visible marks, a session finalize that forgets the prepared direction, a plain run that gains marks, and a second
// decoration that stacks a second header.
func TestSecondOpinionDirectedRunNamesAndMarksPage(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	shipped := secondOpinionFixturePresets(t)
	promptLog := filepath.Join(t.TempDir(), "prompt.txt")
	// The fake editor proposes a slug only when its prompt carries a direction, so the plain case still proves the plain name.
	secondOpinionFakeCodex(t, `printf '%s' "$prompt" > '`+promptLog+`'
 case "$prompt" in *'Direction: None.'*) ;; *) printf 'Variation slug: Add A Timeline!\n' ;; esac`)
	page := func(name string) string {
		data, err := os.ReadFile(filepath.Join(bundle, name))
		if err != nil {
			t.Fatalf("%s not published: %v", name, err)
		}
		return string(data)
	}
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex", "--preset", "bold", "--editor-direction", "add a timeline"))[0]
	prompt, _ := os.ReadFile(promptLog)
	if !strings.Contains(string(prompt), "Direction: Shipped bold: add motion.\n\nadd a timeline\n") {
		t.Fatalf("prompt lacks preset-then-free-text direction: %s", prompt)
	}
	if row["page"] != "ai-reports/2026-10-15_1200-example/index-codex-actual-model-1-bold-add-a-timeline.html" || row["reply"] != "ai-reports/2026-10-15_1200-example/codex-actual-model-1-bold-add-a-timeline-reply.md" {
		t.Fatalf("directed names: %v", row)
	}
	if row["preset"] != "bold" || row["preset_file"] != filepath.Join(shipped, "bold.md") || row["editor_direction"] != "add a timeline" {
		t.Fatalf("row lacks direction provenance: %v", row)
	}
	directed := page("index-codex-actual-model-1-bold-add-a-timeline.html")
	for _, mark := range []string{`<meta name="ai-report-second-opinion-preset" content="bold">`, `<meta name="ai-report-second-opinion-direction" content="Shipped bold: add motion.`, "<!-- second-opinion-header -->", "<!-- /second-opinion-header -->", "add a timeline"} {
		if !strings.Contains(directed, mark) {
			t.Fatalf("directed page lacks %q: %s", mark, directed)
		}
	}
	row = secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex", "--editor-direction", "add a timeline"))[0]
	if row["page"] != "ai-reports/2026-10-15_1200-example/index-codex-actual-model-1-custom-add-a-timeline.html" || row["preset"] != nil {
		t.Fatalf("free-text-only run: %v", row)
	}
	row = secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))[0]
	plain := page("index-codex-actual-model-1.html")
	if row["page"] != "ai-reports/2026-10-15_1200-example/index-codex-actual-model-1.html" || strings.Contains(plain, "second-opinion-header") || strings.Contains(plain, "ai-report-second-opinion-preset") || row["preset"] != nil || row["editor_direction"] != nil {
		t.Fatalf("plain run changed: %v %s", row, plain)
	}
	row = secondOpinionRows(t, secondOpinionInvoke(t, repo, "--prepare", "--backend", "claude", "--model", "claude-fable-5-1", "--preset", "bold"))[0]
	stage, _ := row["staging_page"].(string)
	os.WriteFile(stage, []byte(`<html><body><p>Original report: <a href="index.html">index.html</a></p><h1>Bold</h1></body></html>`), 0644)
	os.WriteFile(strings.TrimSuffix(stage, ".html")+"-reply.md", []byte("Clarity: bolder.\n- Variation slug: Motion Pass\n"), 0644)
	row = secondOpinionRows(t, secondOpinionInvoke(t, repo, "--finalize", stage))[0]
	if row["state"] != "ran" || row["page"] != "ai-reports/2026-10-15_1200-example/index-claude-claude-fable-5-1-bold-motion-pass.html" || row["preset"] != "bold" || row["preset_file"] != filepath.Join(shipped, "bold.md") {
		t.Fatalf("finalize lost the prepared direction: %v", row)
	}
	direction := secondOpinionDirection{Preset: "bold", PresetText: "Go <further> & \"faster\"", EditorDirection: "add a timeline"}
	withHead := []byte(`<html><head><title>T</title><script>const fake = "<!-- second-opinion-header --></head>";</script></head><body class="x"><p>Original report: <a href="index.html">index.html</a></p></body></html>`)
	once := secondOpinionDecorate(withHead, direction)
	if twice := secondOpinionDecorate(once, direction); string(twice) != string(once) {
		t.Fatalf("decoration not idempotent:\n%s\n%s", once, twice)
	}
	if !strings.Contains(string(once), `const fake = "<!-- second-opinion-header --></head>";</script><meta name="ai-report-second-opinion-preset" content="bold"><meta name="ai-report-second-opinion-direction" content="Go &lt;further&gt; &amp; &#34;faster&#34;`) {
		t.Fatalf("metas not before the real </head>, or not escaped: %s", once)
	}
	if string(secondOpinionDecorate(once, secondOpinionDirection{})) != string(withHead) {
		t.Fatalf("plain decoration did not strip the launcher marks: %s", secondOpinionDecorate(once, secondOpinionDirection{}))
	}
}
