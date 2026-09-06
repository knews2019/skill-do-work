## 0.305.18 — Build Integration Test CLI Once per Test Binary in Suiteinstall (2026-09-06)

Integration tests in `suiteinstall` now compile the `do-work-cli` executable once per test binary using `sync.Once` and clean up the temporary directory in `TestMain`, eliminating redundant `go build` invocations across heavy signal handling tests.

- `internal/suiteinstall/suite_commands_test.go` builds the test CLI binary on demand in an isolated temporary directory via `sync.Once` and deletes it upon process completion in `TestMain`.
- Retains signal interruption, recovery, and post-verification exit status assertions while reducing wall time by ~1.15s and total CPU by ~1.00s.
