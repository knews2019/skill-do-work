package main

import (
	"os"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func main() {
	runtime := commandruntime.NewRuntime(os.Stdout, nil)
	os.Exit(runtime.Run(os.Args[1:]))
}
