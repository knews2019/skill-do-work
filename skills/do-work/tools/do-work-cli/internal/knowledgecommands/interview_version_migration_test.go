package knowledgecommands

import "testing"

// REQ-557 / D-03: migrateInterviewSession's same-major branch swapped a lenient
// comparator for the strict one plus an explicit parsed flag. The lenient
// comparator scored an unparseable component as 0 and could never report that it
// had failed to read a version, so a template version like "1.09.0" or "1.0.x"
// was ordered numerically and the session was stamped with it. The strict parser
// refuses instead.
//
// template.Version is read verbatim from a project-authored template's YAML
// frontmatter, so both refusing rows below are reachable from project data, not
// only from tool-written bytes. semverMajor still gates the branch on three
// dot-separated parts with an Atoi-able first part, which a leading zero, an
// empty middle part, a trailing dot and a non-numeric third part all pass.
//
// The refusal is deliberate, and the last row pins its bound: a pair both sides
// of the strict parser accept still stamps exactly as before.
func TestMigrateInterviewSessionRefusesAVersionOnlyTheLenientComparatorCouldOrder(t *testing.T) {
	const bareSemverRefusal = "template and session versions must be bare semver"
	for _, versionCase := range []struct {
		name            string
		sessionVersion  string
		templateVersion string
		wantError       string
		wantStamped     string
	}{
		{
			name:            "session patch part is not a number",
			sessionVersion:  "1.0.x",
			templateVersion: "1.0.1",
			wantError:       bareSemverRefusal,
			wantStamped:     "1.0.x",
		},
		{
			name:            "template minor part carries a leading zero",
			sessionVersion:  "1.0.0",
			templateVersion: "1.09.0",
			wantError:       bareSemverRefusal,
			wantStamped:     "1.0.0",
		},
		{
			name:            "template patch part is empty",
			sessionVersion:  "1.0.0",
			templateVersion: "1.0.",
			wantError:       bareSemverRefusal,
			wantStamped:     "1.0.0",
		},
		{
			name:            "both versions are bare semver and the template is newer",
			sessionVersion:  "1.0.0",
			templateVersion: "1.1.0",
			wantStamped:     "1.1.0",
		},
		{
			name:            "both versions are bare semver and equal",
			sessionVersion:  "1.0.0",
			templateVersion: "1.0.0",
			wantStamped:     "1.0.0",
		},
	} {
		t.Run(versionCase.name, func(t *testing.T) {
			session := map[string]any{"template_version": versionCase.sessionVersion}
			template := interviewTemplate{Version: versionCase.templateVersion}

			migrationPending, migrationError := migrateInterviewSession(template, session)

			if migrationPending {
				t.Fatalf("migrateInterviewSession(%q -> %q) reported a pending major migration",
					versionCase.sessionVersion, versionCase.templateVersion)
			}
			switch {
			case versionCase.wantError == "" && migrationError != nil:
				t.Fatalf("migrateInterviewSession(%q -> %q) refused a bare-semver pair: %v",
					versionCase.sessionVersion, versionCase.templateVersion, migrationError)
			case versionCase.wantError != "" && migrationError == nil:
				t.Fatalf("migrateInterviewSession(%q -> %q) accepted a version the strict parser cannot read; "+
					"the lenient comparator's numeric scoring is back",
					versionCase.sessionVersion, versionCase.templateVersion)
			case versionCase.wantError != "" && migrationError.Error() != versionCase.wantError:
				t.Fatalf("migrateInterviewSession(%q -> %q) error = %q, want %q",
					versionCase.sessionVersion, versionCase.templateVersion, migrationError, versionCase.wantError)
			}
			if stamped := stringValue(session["template_version"]); stamped != versionCase.wantStamped {
				t.Fatalf("migrateInterviewSession(%q -> %q) left template_version %q, want %q",
					versionCase.sessionVersion, versionCase.templateVersion, stamped, versionCase.wantStamped)
			}
		})
	}
}
