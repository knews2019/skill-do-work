package toolboxcommands

import "testing"

func TestAuditRenameAndCopyResolution(t *testing.T) {
	aliases := map[string]string{"old": "new"}
	if got := auditResolveAlias(aliases, "old"); got != "new" {
		t.Fatal(got)
	}
	copyMap := map[string]string{"old": "middle", "middle": "live"}
	if got, ok := auditSurvivingCopy(copyMap, map[string]bool{"live": true}, "old"); !ok || got != "live" {
		t.Fatalf("%q %v", got, ok)
	}
}
