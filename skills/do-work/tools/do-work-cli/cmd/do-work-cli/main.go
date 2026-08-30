package main

import (
	"os"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/cleanup"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/doctor"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/suiteinstall"
)

func main() {
	// stdout carries only the rendered CommandResult; narration and the single install
	// confirmation use stderr and stdin. Registering a command is one import and one map
	// entry, which is the shape every later command family copies.
	handlers := map[string]commandruntime.CommandHandler{}
	for name, handler := range suiteinstall.Handlers(os.Stderr, os.Stdin) {
		handlers[name] = handler
	}
	for name, handler := range cleanup.Handlers() {
		handlers[name] = handler
	}
	for name, handler := range doctor.Handlers() {
		handlers[name] = handler
	}
	runtime := commandruntime.NewRuntime(os.Stdout, handlers)
	os.Exit(runtime.Run(os.Args[1:]))
}
