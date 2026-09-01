package toolboxcommands

import "testing"

func TestHandlersRegisterCanonicalToolboxCommands(t *testing.T) {
	handlers := Handlers()
	for _, name := range []string{CommandNote, CommandArchitecture, CommandReportImage, CommandReportImageBatch, CommandPortfolio, CommandLast30Days, CommandAuditMetrics} {
		if handlers[name] == nil {
			t.Errorf("canonical command %q is not registered", name)
		}
	}
	if len(handlers) != 7 {
		t.Fatalf("registered %d toolbox commands, want 7", len(handlers))
	}
}

func TestRemediationMutationOptionsHaveLeadingRegionAndDoubleDash(t *testing.T) {
	rest, dry, commit, err := parseMutationFlags([]string{"--dry-run", "literal", "--commit"})
	if err != nil || !dry || commit || len(rest) != 2 || rest[1] != "--commit" {
		t.Fatalf("leading region parse=%q dry=%v commit=%v err=%v", rest, dry, commit, err)
	}
	rest, dry, commit, err = parseMutationFlags([]string{"--", "--dry-run", "--commit"})
	if err != nil || dry || commit || len(rest) != 2 {
		t.Fatalf("double-dash parse=%q dry=%v commit=%v err=%v", rest, dry, commit, err)
	}
}
