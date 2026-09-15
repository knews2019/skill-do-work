package knowledgecommands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func populatedInterviewSession() map[string]any {
	entry := func(id, status string, stakeholders []any, details map[string]any) map[string]any {
		return map[string]any{"entry_id": id, "title": id, "summary": id, "cadence": "weekly", "trigger": "Monday", "inputs": []any{}, "stakeholders": stakeholders, "constraints": []any{}, "details": details, "source_confidence": "confirmed", "status": status, "last_validated_at": "2026-09-15"}
	}
	window := func(label string) []any {
		return []any{map[string]any{"label": label, "days": []any{"Mon"}, "start": "09:00", "end": "10:00"}}
	}
	rhythm := func(id, status string) map[string]any {
		return entry(id, status, []any{}, map[string]any{"energy_pattern": "focus", "non_calendar_reality": "Daily interruptions: " + id, "time_windows": window(id), "interruptions": []any{}})
	}
	decision := func(id, status string, inputs []any) map[string]any {
		return entry(id, status, []any{"Alex"}, map[string]any{"decision_name": id, "decision_inputs": inputs, "thresholds": []any{}, "escalation_rule": "never", "reversible": true})
	}
	dependency := func(id, status, needed string) map[string]any {
		return entry(id, status, []any{"Alex", "Sam"}, map[string]any{"dependency_owner": "Sam", "deliverable": id, "needed_by": needed, "failure_impact": "delay", "fallback": "ask"})
	}
	knowledge := func(id, status string) map[string]any {
		return entry(id, status, []any{"Lee"}, map[string]any{"knowledge_area": id, "why_it_matters": "context", "where_it_lives": "in my head", "who_else_knows": []any{}, "risk_if_missing": "high"})
	}
	friction := entry("sync", "active", []any{"Pat"}, map[string]any{"priority": "high", "frequency": "daily", "time_cost": "20 minutes", "systems_involved": "Daily CRM sync interruptions", "current_workaround": "wait", "automation_candidate": true, "time_windows": window("sync")})
	layers := map[string]any{}
	for name, entries := range map[string][]any{
		"operating_rhythms":       {rhythm("focus", "active"), rhythm("old-rhythm", "stale")},
		"recurring_decisions":     {decision("weekly", "active", []any{"Dashboard", "Dashboard"}), decision("old-decision", "stale", []any{"Dashboard", "Retired"})},
		"dependencies":            {dependency("handoff", "active", "weekly Monday at 09:00"), dependency("old-handoff", "stale", "daily at 12:00"), dependency("one-off", "active", "2026-10-15")},
		"institutional_knowledge": {knowledge("context", "active"), knowledge("old-context", "stale")},
		"friction":                {friction},
	} {
		layers[name] = map[string]any{"approved": true, "entries": entries}
	}
	return map[string]any{"template": "work-operating-model", "template_version": "1.1.0", "session_id": "populated", "role_or_name_or_repo": "fixture", "status": "complete", "review_runs": 1, "review_completed_at": "2026-09-15T09:00:00Z", "last_activity_at": "2026-09-15T09:00:00Z", "previous_version": nil, "layers": layers}
}

func TestInterviewExportPreservesPopulatedCanonicalEntriesAndScalarSources(t *testing.T) {
	root := t.TempDir()
	template, err := os.ReadFile("../../../../../do-work-knowledge/interviews/work-operating-model.md")
	if err != nil {
		t.Fatal(err)
	}
	knowledgeRoot := filepath.Join(root, "knowledge")
	writeInterviewFixture(t, knowledgeRoot, "interviews/work-operating-model.md", string(template))
	session := populatedInterviewSession()
	parsed, err := loadInterviewTemplate(filepath.Join(knowledgeRoot, "interviews/work-operating-model.md"))
	if err != nil {
		t.Fatal(err)
	}
	session["template_version"] = parsed.Version
	// An approved empty layer must serialize as [], never null.
	mapValue(session["layers"])["institutional_knowledge"] = map[string]any{"approved": true, "entries": []any{}}
	original, _ := json.Marshal(session)
	writeInterviewFixture(t, root, sessionPath("work-operating-model"), string(original))
	runGitFixture(t, root, "init")
	runGitFixture(t, root, "config", "user.email", "fixture@example.com")
	runGitFixture(t, root, "config", "user.name", "Fixture")
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-qm", "fixture")
	result := handleInterviewExport(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--knowledge-root", knowledgeRoot, "--template", "work-operating-model"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("export: %+v", result)
	}
	exportRoot := filepath.Join(root, "do-work/interview/work-operating-model/exports")
	data, err := os.ReadFile(filepath.Join(exportRoot, "operating-model.json"))
	if err != nil {
		t.Fatal(err)
	}
	var dump map[string]any
	if err := json.Unmarshal(data, &dump); err != nil {
		t.Fatal(err)
	}
	for layer, raw := range mapValue(session["layers"]) {
		got := mapValue(mapValue(dump["layers"])[layer])["entries"]
		want := mapValue(raw)["entries"]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s entries = %#v, want %#v", layer, got, want)
		}
	}
	data, err = os.ReadFile(filepath.Join(exportRoot, "schedule-recommendations.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schedule map[string]any
	if err := json.Unmarshal(data, &schedule); err != nil {
		t.Fatal(err)
	}
	if len(sliceValue(schedule["standing_slots"])) != 1 || len(sliceValue(schedule["avoid_windows"])) != 2 || len(sliceValue(schedule["time_blocks"])) != 1 {
		t.Errorf("published schedule lost derived fields: %s", data)
	}

	for _, name := range []string{"SOUL.md", "HEARTBEAT.md"} {
		data, err := os.ReadFile(filepath.Join(exportRoot, name))
		if err != nil || !strings.Contains(string(data), "- Dashboard") || strings.Contains(string(data), "- {}") || strings.Contains(string(data), "Retired") {
			t.Errorf("%s lost active scalar sources: %s (%v)", name, data, err)
		}
	}
}

func TestInterviewDerivationsUseActiveEntriesAndPopulateDeclaredFields(t *testing.T) {
	session := populatedInterviewSession()
	before, _ := json.Marshal(session)
	derived := deriveInterviewValues(session)
	for key, want := range map[string]int{"time_blocks": 1, "avoid_windows": 2, "standing_slots": 1, "tacit_knowledge": 1} {
		if got := len(sliceValue(derived[key])); got != want {
			t.Errorf("%s count = %d, want %d: %#v", key, got, want, derived[key])
		}
	}
	if !reflect.DeepEqual(derived["authoritative_inputs"], []any{}) || !reflect.DeepEqual(derived["advisory_inputs"], []any{"Dashboard"}) {
		t.Errorf("stale or duplicate decisions affect trust: %#v", derived)
	}
	if !reflect.DeepEqual(derived["overnight_scan_sources"], []any{"Daily interruptions: focus", "Dashboard"}) {
		t.Errorf("scan sources: %#v", derived["overnight_scan_sources"])
	}
	if strings.Contains(stringValue(derived["rhythm_synthesis"]), "old-rhythm") {
		t.Error("stale rhythm in synthesis")
	}
	slots := sliceValue(derived["standing_slots"])
	if len(slots) == 1 {
		want := map[string]any{"label": "handoff", "cadence": "weekly", "day": "Mon", "time": "09:00", "counterparty": "Sam", "source_entries": []any{"dependencies.handoff"}}
		if !reflect.DeepEqual(slots[0], want) {
			t.Errorf("standing slot: %#v, want %#v", slots[0], want)
		}
	}
	wantTones := []any{map[string]any{"stakeholder": "Alex", "tone": "formal"}, map[string]any{"stakeholder": "Lee", "tone": "informal"}, map[string]any{"stakeholder": "Pat", "tone": "terse"}, map[string]any{"stakeholder": "Sam", "tone": "formal"}}
	if !reflect.DeepEqual(derived["stakeholder_tones"], wantTones) {
		t.Errorf("tones: %#v", derived["stakeholder_tones"])
	}
	after, _ := json.Marshal(session)
	if string(before) != string(after) {
		t.Fatal("derivation mutated historical entries")
	}
	// Without supplied windows, even qualifying recurring friction cannot invent times.
	for _, raw := range interviewLayerEntries(session, "friction") {
		delete(mapValue(mapValue(raw)["details"]), "time_windows")
	}
	if got := len(sliceValue(deriveInterviewValues(session)["avoid_windows"])); got != 1 {
		t.Errorf("avoid windows without a supplied friction window: %d", got)
	}
}

func TestInterviewEachRetainsScalarTypesAndObjectContext(t *testing.T) {
	data := map[string]any{"values": []any{"Dashboard \"quoted\"", 9, true, nil}, "objects": []any{map[string]any{"title": "Object"}}}
	rendered, err := renderTemplateBlock(`{{#each values}}"{{this}}"{{#unless @last}},{{/unless}}{{/each}}`, data, data, 0, 1, true)
	if err != nil || rendered != `"Dashboard \"quoted\"",9,true,null` {
		t.Fatalf("scalar JSON iteration: %q, %v", rendered, err)
	}
	rendered, err = renderTemplateBlock(`{{#each objects}}{{title}}{{/each}}`, data, data, 0, 1, false)
	if err != nil || rendered != "Object" {
		t.Fatalf("object iteration: %q, %v", rendered, err)
	}
}

func TestInterviewStandingSlotsParseOnlyRecurringTiming(t *testing.T) {
	for _, tc := range []struct {
		input, cadence, day, clock string
		ok                         bool
	}{
		{"daily", "daily", "rolling", "", true},
		{"every day at 9:05", "daily", "rolling", "09:05", true},
		{"monthly at 14:30", "monthly", "rolling", "14:30", true},
		{"Friday", "weekly", "Fri", "", true},
		{"weekly Monday at 09:00", "weekly", "Mon", "09:00", true},
		{"2026-10-15", "", "", "", false},
		{"when ready", "", "", "", false},
		{"daily at 29:00", "", "", "", false},
	} {
		t.Run(tc.input, func(t *testing.T) {
			cadence, day, clock, ok := parseInterviewCadence(tc.input)
			if cadence != tc.cadence || day != tc.day || clock != tc.clock || ok != tc.ok {
				t.Errorf("got %q %q %q %v", cadence, day, clock, ok)
			}
		})
	}
}

func TestInterviewAvoidanceUsesAnUnambiguousSuppliedWindow(t *testing.T) {
	first := map[string]any{"label": "morning sync", "days": []any{"Mon"}, "start": "09:00", "end": "10:00"}
	second := map[string]any{"label": "afternoon sync", "days": []any{"Tue"}, "start": "14:00", "end": "15:00"}
	details := map[string]any{"time_windows": []any{first, second}}
	if got := suppliedAvoidWindow(details, "Daily afternoon sync delays"); !reflect.DeepEqual(got, second) {
		t.Errorf("wrong window: %#v", got)
	}
	if got := suppliedAvoidWindow(details, "Daily interruptions"); got != nil {
		t.Errorf("invented a choice between windows: %#v", got)
	}
}
