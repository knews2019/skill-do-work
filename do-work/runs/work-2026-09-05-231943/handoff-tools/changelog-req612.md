## 0.305.23 — Remove Redundant Go Test Discovery using Native -skip (2026-09-06)

Replaced pre-execution Go test discovery in `run-go-tests-with-budget.sh` with Go 1.20+ native `-skip`, eliminating subprocess startup while strictly preserving test selection and safety guards.

- Replaces `go test -list` and Python regex assembly with pure-bash `-skip "^($skip_pattern)"`, escaping regex metacharacters with zero child processes.
- Preserves exact 100% test selection equivalence across fast test suites (verified on `queue-kanban` 402 fast tests).
- Retains empty-selection refusal (`no fast Go tests remain after applying the heavy prefixes`) in the post-test results processor with zero discovery overhead.
- Added comprehensive behavior probes in `_dev/tests/run-go-tests-with-budget-behavior.sh` for heavy exclusion, empty-selection refusal, and metacharacter escaping.
