# REQ-561 builder handback

REQ-561 — Add a three-value priority field the selector orders by and the board shows

## Identity

- Branch: `worktree-agent-REQ-561-add-a-three-value-priority-field-the-selector-orders-by-and-the-board-shows`
- Base: `649296d2fe3232fa4ee3bc694ee994832989ceee`
- Commit: `f3d92379c5b9f9740e0cbdcdd58a6552ba30457a`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-561-add-a-three-value-priority-field-the-selector-orders-by-and-the-board-shows`

## Outcome

Added the closed `priority: now | next | later` schema contract with effective default `next`, typed selected/excluded JSON projection distinct from `selection_priority`, ordinary-class stable ordering before fan-out, and UR-expansion ordering without changing explicit REQ order. The board independently parses the same contract, stable-sorts Pending Ready and Pending Waiting, projects raw/invalid evidence, and renders only `now`/`later` badges. Capture and board contracts document user-words-only authorship and queued-addendum changes.

## RED → GREEN evidence

- Core RED command: `go test -count=1 ./internal/schemanormalization ./internal/requestmodel ./internal/resultmodel ./internal/nextselection`. `priority_absent` resolved to `""`, `priority_next` retained `"NEXT"`, and `priority_invalid` retained `"urgent"`; `TestRequestPriorityOrdersOrdinaryReadyWorkBeforeFanOut` selected `[REQ-821 REQ-822]`, wanted `[REQ-823 REQ-822]`.
- Board RED command: `go test -count=1 -run 'TestPendingGroupsSortByRequestPriority|TestParseRequestTicketPriorityContract' .`. `TestPendingGroupsSortByRequestPriority` returned Ready `[REQ-101 REQ-102 REQ-103]`, wanted `[REQ-102 REQ-101 REQ-103]`.
- GREEN: `go test -count=1 ./internal/schemanormalization ./internal/requestmodel ./internal/resultmodel ./internal/nextselection` passed; `DO_WORK_HEAVY_TESTS=1 go test -count=1 ./internal/nextselection` passed; full do-work-cli module tests and `go vet ./...` passed.
- GREEN: queue-kanban `go test -count=1 .` and `go vet ./...` passed; static and live payload regression passed.
- Browser GREEN: `QUEUE_KANBAN_BROWSER='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR=1 go test -count=1 -run '^TestBrowserBehaviorPriorityOrderAndBadges$' -v .` passed in HeadlessChrome 152.0.0.0. Exact inspected URL for both theme evaluations: `file:///var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/TestBrowserBehaviorPriorityOrderAndBadges2124794665/002/probe.html`. Light and dark both returned Ready `[REQ-403 REQ-401 REQ-402]`, Waiting `[REQ-405 REQ-404]`, four contained non-overlapping now/later badges, and no next badge.
- An initial attempt to reuse the file-page harness for live navigation correctly refused its mismatched URL (`measured on "http://127.0.0.1:57767/", not the probe page`). Live/static payload parity was then verified at the supported server seam by `go test -count=1 -run TestStaticAndLivePriorityPayloadsAgree .`; no failed live-DOM claim is retained.
- Contract GREEN: `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/shipped-package-reference-contract.sh` passed.
- Hygiene: `git diff --check`, `gofmt`, and exact manifest review passed.

## Manifest

- `skills/do-work-board/actions/board.md`
- `skills/do-work-board/tools/queue-kanban/generate.go`
- `skills/do-work-board/tools/queue-kanban/model.go`
- `skills/do-work-board/tools/queue-kanban/model_test.go`
- `skills/do-work-board/tools/queue-kanban/priority_browser_probe_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/serve_test.go`
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work/actions/capture-reference.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go`
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go`
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go`
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go`

No tracked `do-work/`, release metadata, queue record, generated site, or screenshot changed on the builder branch.

## Decisions and deviations

- Kept authored `priority` separate from scheduling-class `selection_priority` in every typed result.
- Kept explicit REQ token order caller-authored; only default scans and each UR expansion receive authored-priority ordering.
- Used `serve_test.go` to decode and compare both generated static and live data instead of adding a separate `generate_test.go` change; this exercises the shared real-fixture parse/projection seam once and avoids duplicate fixture machinery.
- Used the dedicated browser test from exploration; the clipboard contract and browser harness remained unchanged.
- Release metadata, queue stamps, and REQ-530 terminal provenance remain orchestrator-only integration work.

## Guidance read

The builder lane used `CLAUDE.md`, the REQ brief/plan/exploration, action and board primes, their full lesson satellites, do-work-cli and queue-kanban shipped primes/lessons, release guidance, and the general, coding-guardrails, backend, testing, and communication crew files. The implementation applies the alternate-writer lesson by changing the executable core schema and the independent board parser in the same increment.

## Discovered tasks

No critical follow-up was found. The browser probe covers static exact-URL evidence; the live HTTP payload agreement is pinned in `serve_test.go`, while an additional live DOM probe would duplicate the same client path and stays report-only.
