package schemanormalization

import "testing"

func TestNormalizeFieldAppliesAliasesDefaultsAndExactWarnings(t *testing.T) {
	testCases := []struct {
		name           string
		fieldName      string
		rawValue       string
		wantValue      string
		wantRecognized bool
		wantDefault    bool
		wantWarning    string
	}{
		{"domain alias", "domain", " back_end ", "backend", true, false, ""},
		{"status alias", "status", "DONE", "completed", true, false, ""},
		{"route case", "route", "b", "B", true, false, ""},
		{"testing alias", "testing_status", "returned with feedback", "returned", true, false, ""},
		{"legacy effort", "effort_estimate", "trivial", "effort-mechanical", true, false, ""},
		{"absent default is quiet", "impact", "", "impact-user-visible", true, true, ""},
		{
			"unknown warns and defaults", "domain", "server", "general", false, true,
			"⚠ domain: 'server' not recognized — expected one of [frontend, backend, ui-design, general, security, testing, cms]. Treating as 'general'.",
		},
		{
			"unknown without default retains evidence", "status", "almost-done", "almost-done", false, false,
			"⚠ status: 'almost-done' not recognized — expected one of [pending, claimed, completed, completed-with-issues, failed, cancelled, pending-answers, blocked, blocked-archive-collision, blocked-dependency-cycle]. No default is defined; reporting it unchanged.",
		},
		{"verbatim field", "assigned_to", " Cloud-Alpha ", "Cloud-Alpha", true, false, ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := NormalizeField(testCase.fieldName, testCase.rawValue)
			if result.ResolvedValue != testCase.wantValue || result.IsRecognized != testCase.wantRecognized || result.IsDefaulted != testCase.wantDefault || result.WarningMessage != testCase.wantWarning {
				t.Fatalf("NormalizeField(%q, %q) = %#v", testCase.fieldName, testCase.rawValue, result)
			}
		})
	}
}

func TestPreferredFieldGivesCanonicalKeyPrecedence(t *testing.T) {
	fields := map[string][]string{
		"depends_on":   {"REQ-10"},
		"dependencies": {"REQ-99"},
		"parent":       {"REQ-4"},
	}
	values, sourceKey, found := PreferredField(fields, "depends_on", "dependencies")
	if !found || sourceKey != "depends_on" || len(values) != 1 || values[0] != "REQ-10" {
		t.Fatalf("preferred dependency = %v, %q, %v", values, sourceKey, found)
	}
	values, sourceKey, found = PreferredField(fields, "addendum_to", "amends", "parent", "amendment_to")
	if !found || sourceKey != "parent" || len(values) != 1 || values[0] != "REQ-4" {
		t.Fatalf("preferred addendum = %v, %q, %v", values, sourceKey, found)
	}
}

func TestTerminalPredicatesKeepFailureAndCancellationDistinct(t *testing.T) {
	if !IsTerminalSuccess("done") || !DependencySatisfied("completed-with-issues") {
		t.Fatal("completion aliases and completed-with-issues must satisfy dependencies")
	}
	if IsTerminalSuccess("cancelled") || DependencySatisfied("failed") {
		t.Fatal("cancelled and failed must not satisfy dependencies")
	}
	if !IsTerminalResolved("cancelled") || IsTerminalResolved("failed") {
		t.Fatal("cancelled resolves a REQ; failed does not")
	}
	if !IsStopped("failed") || IsStopped("blocked") {
		t.Fatal("failed is stopped; blocked remains unfinished")
	}
}
