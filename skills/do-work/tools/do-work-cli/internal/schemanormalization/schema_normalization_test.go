package schemanormalization

import (
	"reflect"
	"testing"
)

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
		{"gate deferred alias", "gate_deferred", "yes", "true", true, false, ""},
		{"repair default", "repository_gate_repair", "", "false", true, true, ""},
		{"deferred base is verbatim", "deferred_implementation_base", " abc123 ", "abc123", true, false, ""},
		{"absent default is quiet", "impact", "", "impact-user-visible", true, true, ""},
		{"priority absent", "priority", "", "next", true, true, ""},
		{"priority now", "priority", " now ", "now", true, false, ""},
		{"priority next", "priority", "NEXT", "next", true, false, ""},
		{"priority later", "priority", "later", "later", true, false, ""},
		{
			"priority invalid", "priority", "urgent", "next", false, true,
			"⚠ priority: 'urgent' not recognized — expected one of [now, next, later]. Treating as 'next'.",
		},
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

func TestSchemaAliasesAreCentralizedAndDefensivelyCopied(t *testing.T) {
	wantAliases := map[string][]string{
		"addendum_to":    {"amends", "parent", "amendment_to"},
		"depends_on":     {"dependencies"},
		"batch":          {"batch_name"},
		"related":        {"related_reqs"},
		"suggested_spec": {"spec_hint", "suggested-spec"},
	}
	for canonicalKey, want := range wantAliases {
		aliases := SchemaFieldAliases(canonicalKey)
		if !reflect.DeepEqual(aliases, want) {
			t.Fatalf("SchemaFieldAliases(%q) = %v, want %v", canonicalKey, aliases, want)
		}
		aliases[0] = "mutated"
		if got := SchemaFieldAliases(canonicalKey); !reflect.DeepEqual(got, want) {
			t.Fatalf("caller mutation changed aliases for %q: %v", canonicalKey, got)
		}
		for _, aliasKey := range want {
			if gotCanonical, found := CanonicalFieldForAlias(aliasKey); !found || gotCanonical != canonicalKey {
				t.Fatalf("CanonicalFieldForAlias(%q) = %q, %v", aliasKey, gotCanonical, found)
			}
		}
	}
	if aliases := SchemaFieldAliases("unknown"); aliases != nil {
		t.Fatalf("unknown aliases = %v, want nil", aliases)
	}
	if canonicalKey, found := CanonicalFieldForAlias("depends_on"); found || canonicalKey != "" {
		t.Fatalf("canonical key reported as alias: %q, %v", canonicalKey, found)
	}
}

func TestNormalizeFieldDistinguishesExactCanonicalAuthoringValues(t *testing.T) {
	testCases := []struct {
		name      string
		fieldName string
		rawValue  string
		canonical bool
	}{
		{"canonical enum", "status", "pending", true},
		{"case-normalized enum", "status", "PENDING", false},
		{"read alias", "status", "done", false},
		{"defaulted empty", "impact", "", false},
		{"invalid default", "impact", "visible", false},
		{"canonical boolean", "tdd", "true", true},
		{"boolean alias", "tdd", "yes", false},
		{"verbatim", "deferred_implementation_base", "abc123", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := NormalizeField(testCase.fieldName, testCase.rawValue)
			if result.IsCanonicalAuthoringValue != testCase.canonical {
				t.Fatalf("NormalizeField(%q, %q).IsCanonicalAuthoringValue = %v", testCase.fieldName, testCase.rawValue, result.IsCanonicalAuthoringValue)
			}
		})
	}
}

func TestEverySchemaContractPreservesReadAliasesAndExactWriteEvidence(t *testing.T) {
	for fieldName, contract := range fieldContracts {
		for _, canonicalValue := range contract.canonicalValues {
			result := NormalizeField(fieldName, canonicalValue)
			if !result.IsCanonicalAuthoringValue || result.ResolvedValue != canonicalValue || !result.IsRecognized || result.IsDefaulted || result.WarningMessage != "" {
				t.Errorf("canonical %s=%q: %#v", fieldName, canonicalValue, result)
			}
		}
		for aliasValue, canonicalValue := range contract.aliasValues {
			result := NormalizeField(fieldName, aliasValue)
			if result.IsCanonicalAuthoringValue || result.ResolvedValue != canonicalValue || !result.IsRecognized || result.IsDefaulted || result.WarningMessage != "" {
				t.Errorf("legacy %s=%q: %#v", fieldName, aliasValue, result)
			}
		}
		result := NormalizeField(fieldName, "")
		if result.IsCanonicalAuthoringValue || result.ResolvedValue != contract.defaultValue || !result.IsRecognized || !result.IsDefaulted || result.WarningMessage != "" {
			t.Errorf("absent %s: %#v", fieldName, result)
		}
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

// TestDependencySourceReadyAcceptsClaimedWithCommitOnly pins REQ-570: the heavy
// hold is a phase of a claimed request, so a claimed dependency that already
// landed its implementation commit is source-ready. Every other status is not,
// commit or no commit — the rule is the condition, not a second status list, so
// no retired value is named here. The repository sweep in the REQ proves the
// retired hold status is gone from the tree.
//
// The withdrawn cases are the fail-closed half. This is the rule both the CLI
// and the board read, and callers reach it with whatever the request carried,
// so a commit that is only whitespace must not grant authority any more than an
// absent one does.
func TestDependencySourceReadyAcceptsClaimedWithCommitOnly(t *testing.T) {
	if !DependencySourceReady("claimed", "abc123") {
		t.Fatal("a claimed request holding its landed implementation commit must be source-ready")
	}
	for _, withdrawnCommit := range []string{"", "   ", "\n\t"} {
		if DependencySourceReady("claimed", withdrawnCommit) {
			t.Fatalf("a claimed request whose commit is %q must fail closed", withdrawnCommit)
		}
	}
	if DependencySourceReady("pending", "abc123") {
		t.Fatal("pending remediation must block dependencies even with a stale commit")
	}
	for _, unrecognizedStatus := range []string{"pending-answers", "blocked", "almost-done"} {
		if DependencySourceReady(unrecognizedStatus, "abc123") {
			t.Fatalf("%s must not grant source readiness", unrecognizedStatus)
		}
	}
}
