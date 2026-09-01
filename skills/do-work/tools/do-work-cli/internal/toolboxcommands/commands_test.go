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
