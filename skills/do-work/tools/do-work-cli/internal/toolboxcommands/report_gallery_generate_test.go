package toolboxcommands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestAIReportGalleryGenerationLineage(t *testing.T) {
	root := t.TempDir()
	const evidence = "evidence/measurement.tsv"
	const manifest = "evidence/manifest.json"
	const parent = "ai-reports/2026-09_1-report/2026-09-06_1606_first/index.html"
	const child = "ai-reports/2026-11-19_0900_child/index.html"
	const sidecar = "ai-reports/2026-09_1-report/2026-09-06_1606_first/gallery-generation.json"
	const feedback = "do-work/report-gallery-feedback.json"

	for file, content := range map[string]string{
		evidence: "elapsed_seconds\n48.606\n",
		manifest: `{"source_id":"fixture-baseline-seconds","files":["evidence/measurement.tsv"]}`,
		parent:   `<title>First</title><figure id="timing"><figcaption>48.606 seconds</figcaption></figure>`,
		child:    `<title>Child</title><figure id="comparison"><figcaption>Wall elapsed 48.606 s</figcaption></figure>`,
		"ai-reports/2026-09_1-report/2026-09-06_1606_first/gallery-design.md": "---\ntitle: Fixture timing\nuse_when: Show elapsed time.\n---\n\n## Borrow\n- Lead with the median.\n- Spell out the units.\n",
		"ai-reports/2026-11-19_0900_child/gallery-design.md":                  "---\ntitle: Fixture comparison\nuse_when: Compare elapsed times.\n---\n\n## Borrow\n- Compare the median.\n- Spell out the units.\n",
	} {
		writeGalleryFixture(t, root, file, content)
	}
	run := func(args ...string) string {
		t.Helper()
		result := runGalleryFixture(t, root, args...)
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil {
			t.Fatalf("%v failed: %+v", args, result.Findings)
		}
		return strings.TrimSpace(*result.ExactTextOutput)
	}
	reject := func(reason string, args ...string) {
		t.Helper()
		result := runGalleryFixture(t, root, args...)
		if result.Outcome == resultmodel.OutcomeSuccess || len(result.Findings) != 1 ||
			!strings.Contains(strings.Join(result.Findings[0].Evidence, " "), reason) {
			t.Fatalf("%v should reject with %q, got %s %+v", args, reason, result.Outcome, result.Findings)
		}
	}
	decode := func(output string) map[string]json.RawMessage {
		t.Helper()
		var fields map[string]json.RawMessage
		if err := json.Unmarshal([]byte(output), &fields); err != nil {
			t.Fatal(err)
		}
		return fields
	}
	writeRecord := func(plan, before, after, executedModel string) {
		t.Helper()
		record, err := json.Marshal(map[string]any{
			"plan": json.RawMessage(plan),
			"execution": map[string]string{
				"provider_instance_id": "fixture-codex", "executed_model": executedModel,
				"status": "completed", "task_id": "fixture-task", "thread_id": "fixture-thread", "run_id": "fixture-run",
			},
			"changes": []map[string]string{{"before": before, "after": after, "description": "Make the elapsed unit visible."}},
		})
		if err != nil {
			t.Fatal(err)
		}
		writeGalleryFixture(t, root, "record.json", string(record))
	}
	writeReview := func(digest, renderVerdict, verdict, reason string) {
		t.Helper()
		content, err := json.Marshal(map[string]any{
			"bundle_sha256": digest, "render_verdict": renderVerdict, "verdict": verdict,
			"findings": []map[string]string{{"anchor": "timing", "reason": reason}},
		})
		if err != nil {
			t.Fatal(err)
		}
		writeGalleryFixture(t, root, "review.json", string(content))
	}

	initial := run("--prepare-generation", manifest)
	fields := decode(initial)
	if string(fields["parent_report"]) != `""` || string(fields["marks"]) != `[]` ||
		!strings.Contains(string(fields["evidence"]), `"fixture-baseline-seconds"`) {
		t.Fatalf("initial plan should pin real evidence and explicitly omit parent and marks: %s", initial)
	}
	if _, err := os.Stat(filepath.Join(root, feedback)); !os.IsNotExist(err) {
		t.Fatalf("initial preparation must not create feedback: %v", err)
	}

	writeRecord(initial, "", parent+"#timing", "gpt-6.2-sol")
	reject("completed T3 execution", "--record-generation", parent, "record.json")
	writeRecord(initial, "", parent+"#timing", "gpt-6.2-astra")
	run("--record-generation", parent, "record.json")
	stored, err := os.ReadFile(filepath.Join(root, sidecar))
	if err != nil || !strings.Contains(string(stored), `"executed_model": "gpt-6.2-astra"`) {
		t.Fatalf("generation receipt not saved: %v %s", err, stored)
	}
	reject("already exists", "--record-generation", parent, "record.json")
	reject("no factual review", "--prepare-generation", manifest, parent)

	digest := run("--generation-digest", parent)
	writeReview(digest, "skipped", "pass", "Browser skipped.")
	reject("requires a verified render", "--record-factual-review", parent, "review.json")
	writeReview(digest, "skipped", "unverified", "Browser skipped.")
	run("--record-factual-review", parent, "review.json")
	reject("factual review unverified", "--prepare-generation", manifest, parent)
	writeReview(digest, "pass", "pass", "The rendered figure says 48.606 seconds, matching measurement.tsv.")
	run("--record-factual-review", parent, "review.json")
	reject("no current feedback", "--prepare-generation", manifest, parent)

	const reportMark = `{"target_kind":"report","report_path":"` + parent + `","mark":"like"}`
	const figureMark = `{"target_kind":"figure","report_path":"` + parent + `","figure_id":"timing","mark":"dislike","note":"Show median in seconds."}`
	writeGalleryFixture(t, root, feedback, `{"version":1,"entries":[`+reportMark+`,`+figureMark+`,{"target_kind":"figure","report_path":"`+child+`","figure_id":"timing","mark":"like"},{"target_kind":"figure","report_path":"`+parent+`","figure_id":"removed","mark":"like"}]}`)
	followup := run("--prepare-generation", manifest, parent)
	fields = decode(followup)
	if string(fields["parent_report"]) != `"`+parent+`"` || string(fields["marks"]) != `[`+reportMark+`,`+figureMark+`]` ||
		string(fields["evidence"]) != string(decode(initial)["evidence"]) {
		t.Fatalf("child plan must preserve pin and copy exact current marks for this report: %s", followup)
	}
	writeRecord(followup, parent+"#timing", child+"#missing", "gpt-6.2-astra")
	reject("missing", "--record-generation", child, "record.json")
	writeRecord(followup, parent+"#timing", child+"#comparison", "gpt-6.2-astra")
	// Later feedback edits cannot replace the marks the generator already used.
	writeGalleryFixture(t, root, feedback, `{"version":1,"entries":[`+reportMark+`]}`)
	run("--record-generation", child, "record.json")
	childRecord, err := os.ReadFile(filepath.Join(root, "ai-reports/2026-11-19_0900_child/gallery-generation.json"))
	if err != nil || !strings.Contains(string(childRecord), "Show median in seconds.") {
		t.Fatalf("child lost the snapshotted figure note: %v %s", err, childRecord)
	}

	writeGalleryFixture(t, root, evidence, "elapsed_seconds\n49\n")
	reject("evidence pin", "--prepare-generation", manifest, parent)
	writeGalleryFixture(t, root, evidence, "elapsed_seconds\n48.606\n")

	writeGalleryFixture(t, root, filepath.ToSlash(filepath.Join(filepath.Dir(parent), "presentation.css")), "figure { color: red; }")
	reject("stale factual review", "--prepare-generation", manifest, parent)
	writeGalleryFixture(t, root, parent, `<title>First</title><figure id="timing"><figcaption>48.606 milliseconds</figcaption></figure>`)
	writeReview(digest, "pass", "pass", "Old review attempted after an edit.")
	reject("changed", "--record-factual-review", parent, "review.json")
	wrongDigest := run("--generation-digest", parent)
	writeReview(wrongDigest, "pass", "fail", "Rendered milliseconds contradict elapsed seconds in measurement.tsv.")
	run("--record-factual-review", parent, "review.json")
	writeReview(wrongDigest, "pass", "pass", "A like is not a factual receipt.")
	reject("verdict already exists", "--record-factual-review", parent, "review.json")
	reject("milliseconds contradict", "--prepare-generation", manifest, parent)
	resolved := run("--resolve", "report:ai-reports/2026-09_1-report/2026-09-06_1606_first/")
	if !strings.Contains(resolved, "gallery-design.md") {
		t.Fatalf("factually disqualified candidate lost its reusable design: %s", resolved)
	}
}
