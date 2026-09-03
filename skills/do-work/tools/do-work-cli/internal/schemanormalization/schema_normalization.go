// Package schemanormalization owns the request schema's read-time aliases,
// defaults, warnings, and status predicates.
package schemanormalization

import (
	"fmt"
	"sort"
	"strings"
)

type fieldContract struct {
	canonicalValues []string
	aliasValues     map[string]string
	defaultValue    string
	upperCase       bool
	verbatim        bool
}

var fieldContracts = map[string]fieldContract{
	"domain": {
		canonicalValues: []string{"frontend", "backend", "ui-design", "general", "security", "testing", "cms"},
		aliasValues: map[string]string{
			"back-end": "backend", "back_end": "backend", "front-end": "frontend", "front_end": "frontend",
			"ui_design": "ui-design", "sec": "security", "test": "testing",
			"content-management": "cms", "content_management": "cms",
		},
		defaultValue: "general",
	},
	"status": {
		canonicalValues: []string{"pending", "claimed", "completed", "completed-with-issues", "failed", "cancelled", "pending-answers", "pending-heavy-testing", "blocked", "blocked-archive-collision", "blocked-dependency-cycle"},
		aliasValues: map[string]string{
			"complete": "completed", "done": "completed", "finished": "completed", "closed": "completed",
			"canceled": "cancelled", "abandoned": "cancelled", "wont-do": "cancelled", "wontfix": "cancelled",
		},
	},
	"route": {
		canonicalValues: []string{"A", "B", "C"},
		aliasValues:     map[string]string{},
		upperCase:       true,
	},
	"caveman": {
		canonicalValues: []string{"false", "true", "lite", "full", "ultra"},
		aliasValues:     map[string]string{"yes": "true", "on": "true", "light": "lite"},
		defaultValue:    "false",
	},
	"maintenance": {
		canonicalValues: []string{"true", "false"},
		aliasValues:     map[string]string{"yes": "true", "on": "true", "t": "true", "no": "false", "off": "false", "f": "false"},
		defaultValue:    "false",
	},
	"tdd": {
		canonicalValues: []string{"true", "false"},
		aliasValues:     map[string]string{"test_first": "true", "yes": "true", "on": "true", "t": "true", "no": "false", "off": "false", "f": "false"},
		defaultValue:    "false",
	},
	"error_type": {
		canonicalValues: []string{"intent", "spec", "code", "environment"},
		aliasValues:     map[string]string{},
		defaultValue:    "code",
	},
	"kb_status": {
		canonicalValues: []string{"promoted", "pending", "declined", "skipped"},
		aliasValues:     map[string]string{"skip": "skipped", "rejected": "declined"},
		defaultValue:    "pending",
	},
	"impact": {
		canonicalValues: []string{"impact-critical", "impact-user-visible", "impact-rule-change", "impact-negligible"},
		aliasValues:     map[string]string{},
		defaultValue:    "impact-user-visible",
	},
	"priority": {
		canonicalValues: []string{"now", "next", "later"},
		aliasValues:     map[string]string{},
		defaultValue:    "next",
	},
	"effort_estimate": {
		canonicalValues: []string{"effort-mechanical", "effort-substantive"},
		aliasValues:     map[string]string{"trivial": "effort-mechanical", "normal": "effort-substantive"},
		defaultValue:    "effort-substantive",
	},
	"testing_status": {
		canonicalValues: []string{"in-testing", "tested", "returned"},
		aliasValues: map[string]string{
			"in_testing": "in-testing", "in testing": "in-testing", "testing": "in-testing",
			"selected-for-testing": "in-testing", "selected for testing": "in-testing",
			"returned-with-feedback": "returned", "returned_with_feedback": "returned", "returned with feedback": "returned",
		},
	},
	"builder_decided": {
		canonicalValues: []string{"true"},
		aliasValues:     map[string]string{},
		defaultValue:    "false",
	},
	"gate_deferred": {
		canonicalValues: []string{"true", "false"},
		aliasValues:     map[string]string{"yes": "true", "on": "true", "t": "true", "no": "false", "off": "false", "f": "false"},
		defaultValue:    "false",
	},
	"repository_gate_repair": {
		canonicalValues: []string{"true", "false"},
		aliasValues:     map[string]string{"yes": "true", "on": "true", "t": "true", "no": "false", "off": "false", "f": "false"},
		defaultValue:    "false",
	},
	"deferred_implementation_base":  {verbatim: true},
	"deferred_implementation_merge": {verbatim: true},
}

// FieldResult is the complete evidence for one schema-backed field read.
type FieldResult struct {
	FieldName      string
	OriginalValue  string
	ResolvedValue  string
	IsRecognized   bool
	IsDefaulted    bool
	WarningMessage string
}

// SchemaFieldNames returns every field governed by the Schema Read Contract.
func SchemaFieldNames() []string {
	fieldNames := make([]string, 0, len(fieldContracts))
	for fieldName := range fieldContracts {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)
	return fieldNames
}

// NormalizeField resolves one value under the Schema Read Contract.
func NormalizeField(fieldName string, rawValue string) FieldResult {
	trimmedValue := strings.TrimSpace(rawValue)
	contract, contracted := fieldContracts[fieldName]
	if !contracted {
		return FieldResult{FieldName: fieldName, OriginalValue: rawValue, ResolvedValue: trimmedValue, IsRecognized: true}
	}
	if trimmedValue == "" {
		return FieldResult{FieldName: fieldName, OriginalValue: rawValue, ResolvedValue: contract.defaultValue, IsRecognized: true, IsDefaulted: true}
	}
	if contract.verbatim {
		return FieldResult{FieldName: fieldName, OriginalValue: rawValue, ResolvedValue: trimmedValue, IsRecognized: true}
	}
	normalizedValue := strings.ToLower(trimmedValue)
	if contract.upperCase {
		normalizedValue = strings.ToUpper(trimmedValue)
	}
	if aliasValue, isAlias := contract.aliasValues[normalizedValue]; isAlias {
		normalizedValue = aliasValue
	}
	for _, canonicalValue := range contract.canonicalValues {
		if normalizedValue == canonicalValue {
			return FieldResult{FieldName: fieldName, OriginalValue: rawValue, ResolvedValue: normalizedValue, IsRecognized: true}
		}
	}
	warning := ""
	usedDefault := true
	resolvedValue := contract.defaultValue
	if contract.defaultValue == "" {
		warning = fmt.Sprintf("⚠ %s: '%s' not recognized — expected one of [%s]. No default is defined; reporting it unchanged.",
			fieldName, trimmedValue, strings.Join(contract.canonicalValues, ", "))
		usedDefault = false
		resolvedValue = normalizedValue
	} else {
		warning = fmt.Sprintf("⚠ %s: '%s' not recognized — expected one of [%s]. Treating as '%s'.",
			fieldName, trimmedValue, strings.Join(contract.canonicalValues, ", "), contract.defaultValue)
	}
	return FieldResult{
		FieldName: fieldName, OriginalValue: rawValue, ResolvedValue: resolvedValue,
		IsRecognized: false, IsDefaulted: usedDefault, WarningMessage: warning,
	}
}

// PreferredField returns the canonical key when present, otherwise the first
// present alias. Presence wins even when the selected value is empty.
func PreferredField(fields map[string][]string, canonicalKey string, aliasKeys ...string) (values []string, sourceKey string, found bool) {
	if values, found := fields[canonicalKey]; found {
		return append([]string(nil), values...), canonicalKey, true
	}
	for _, aliasKey := range aliasKeys {
		if values, found := fields[aliasKey]; found {
			return append([]string(nil), values...), aliasKey, true
		}
	}
	return nil, "", false
}

// IsTerminalSuccess reports completed work that satisfies dependencies.
func IsTerminalSuccess(status string) bool {
	normalizedStatus := NormalizeField("status", status).ResolvedValue
	return normalizedStatus == "completed" || normalizedStatus == "completed-with-issues"
}

// IsTerminalResolved reports work that no longer holds a user request open.
func IsTerminalResolved(status string) bool {
	normalizedStatus := NormalizeField("status", status).ResolvedValue
	return IsTerminalSuccess(normalizedStatus) || normalizedStatus == "cancelled"
}

// IsStopped reports work that ended, including a failed attempt.
func IsStopped(status string) bool {
	normalizedStatus := NormalizeField("status", status).ResolvedValue
	return IsTerminalResolved(normalizedStatus) || normalizedStatus == "failed"
}

// DependencySatisfied reports whether a dependency status unblocks a dependent.
func DependencySatisfied(status string) bool { return IsTerminalSuccess(status) }
