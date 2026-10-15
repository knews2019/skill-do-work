package toolboxcommands

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	os.WriteFile(template, []byte("Output: {{OUTPUT_FILE}}\nPrefix: {{ASSET_PREFIX}}\nRepository: {{REPOSITORY_ROOT}}\nBundle: {{BUNDLE_PATH}}\n{{BACKLINK_REQUIREMENT}}\n"), 0644)
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
		if _, _, err := secondOpinionMetadata(test.backend, secondOpinionCaptured{stdout: test.stdout, stderr: test.stderr}); err == nil {
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
		done <- secondOpinionCapture(ctx, directory, "/bin/sh", "-c", "trap '' TERM; (printf ready > ready; /bin/sleep 0.3; printf late > cancel-late; /bin/sleep 30) & wait")
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

func TestSecondOpinionFreshGenerationModelMismatchRefusesPublication(t *testing.T) {
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
	files, _ := os.ReadDir(bundle)
	if rows[0]["state"] != "failed" || rows[0]["actual_model"] != "different-model" || len(files) != 1 {
		t.Fatalf("model switch published: %v files=%v", rows, files)
	}
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
			captured := secondOpinionCapture(ctx, directory, "/bin/sh", "-c", `(/bin/sleep 0.3; printf late > late-write; /bin/sleep 0.6)`+redirect+` &`)
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

func TestSecondOpinionCapturePreservesBackendStatusAndChannels(t *testing.T) {
	for _, status := range []string{"0", "7"} {
		captured := secondOpinionCapture(context.Background(), t.TempDir(), "/bin/sh", "-c", `printf '%s' "$1"; printf warning >&2; if read value; then exit 99; fi; if (printf forged >&3) 2>/dev/null; then exit 98; fi; exit "$2"`, "backend", `literal $() ' argument`, status)
		if captured.stdout != `literal $() ' argument` || captured.stderr != "warning" {
			t.Fatalf("capture=%+v", captured)
		}
		if status == "0" && captured.err != nil || status == "7" && (captured.err == nil || !strings.Contains(captured.err.Error(), "7")) {
			t.Fatalf("status=%s capture=%+v", status, captured)
		}
	}
}
