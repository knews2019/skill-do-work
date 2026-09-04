package heavyverification

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/ownedprocess"
)

// runFingerprintProbe bounds discovery as well as execution. An unavailable or
// hung runtime is uncertain evidence, never permission to reuse a recent pass.
func runFingerprintProbe(repositoryRoot string, argv []string, timeout time.Duration) ([]byte, error) {
	command := exec.Command(argv[0], argv[1:]...)
	command.Dir = repositoryRoot
	if !ownedprocess.ConfigureGroup(command) {
		return nil, fmt.Errorf("toolchain probe process ownership is unavailable")
	}
	var output bytes.Buffer
	command.Stdout, command.Stderr = &output, &output
	if err := command.Start(); err != nil {
		return nil, err
	}
	waited := make(chan error, 1)
	go func() { waited <- command.Wait() }()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err := <-waited:
		return output.Bytes(), err
	case <-timer.C:
		_ = ownedprocess.TerminateGroup(command.Process.Pid, laneTerminationGracePeriod)
		<-waited
		return nil, fmt.Errorf("toolchain probe timed out")
	}
}
