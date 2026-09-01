package hookcommands

import (
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func TestHandlersRegisterAllHookCommands(t *testing.T) {
	handlers := Handlers(strings.NewReader(""))
	for _, commandName := range []string{CommandSessionStart, CommandMemorySessionStart, CommandMemoryStopCapture} {
		if handlers[commandName] == nil {
			t.Errorf("missing handler %q", commandName)
		}
	}
	result := handlers[CommandMemoryStopCapture](commandruntime.ExecutionContext{RepositoryRoot: t.TempDir()}, []string{"unexpected"})
	if result.Outcome != "failure" {
		t.Fatalf("unexpected argument result=%+v", result)
	}
}
