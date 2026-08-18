# REQ-251 builder brief

**Route A — direct implementation.** Estimated 5 active minutes (P50, high confidence, trivial). Two text-only edits, no behaviour change: two fixture literals in `verify_test.go` (lines ~1186/~1230) hold the retired reversed-span wording as input, and `timestamp.go:35`'s comment overstates its own test's claim — bring it in line with `timestamp_test.go:43-46`'s twin. Verify with the REQ's grep requirement: after the change, the retired wording appears nowhere in shipped source. Line numbers may have drifted slightly — REQ-248 just landed in this package; locate by content, not by line. Run the package tests (`go test -count=1 ./...` from the tool dir) and the full gate.

## How this build runs

You are a **worktree builder** dispatched by the do-work work pipeline. Everything binding is in this brief.

**Your tree, your branch.** Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-251-retire-the-stale-copies-of-the-future-stamp-message` — a full checkout on branch `worktree-agent-REQ-251-retire-the-stale-copies-of-the-future-stamp-message`, cut from integration tip `ad69e56`.

- Never write anything under `/home/user/skill-do-work` — the one exception is your hand-back file, named below.
- Never read or write `do-work/` in your own worktree (stale snapshot; your REQ body is inlined below).
- Commit on your own branch, in small increments so an interruption costs one step. Do not touch `VERSION`, `CHANGELOG.md`, or `skills/do-work/actions/version.md` — serial-only, integrator-owned.
- A needed one-line edit to a file outside your write set is an *integration seam*: hand back the exact line and where it goes; do not edit the file. A larger need: stop and report in your hand-back.
- Out-of-scope finds go in `## Discovered Tasks` in your hand-back — never fixed inline.

**Crew rules** (read from your own worktree before writing code): `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`. Read every `prime_files` path too.

**P-A-U phasing is mandatory** — work the [PLAN]/[APPLY]/[UNIFY] block in your REQ body; record the evidence in your hand-back (the orchestrator transcribes it and audits it against the diff). Log significant choices as D-XX with reasoning (DECIDE & STATE vs ESCALATE with Value/Risk).

## Environment notes

- `bash _dev/tests/maintainer-verify.sh` exits 0 at your branch point — your baseline and your gate. Exit code is the only proof; never pipe it through `tail`.
- Toolchain present: Go 1.26.1, ShellCheck 0.11.0, `just`, Node 22, Chromium (Playwright, `/opt/pw-browsers/chromium`).
- **Never run bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — build to scratch (`go build -o /tmp/<name> .`).
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp; never carry or compute one.
- Build test fixtures in scratch space, never in this repo's own `do-work/`.

## Hand-back

Write your report to exactly this absolute path (the one main-tree write you may make; never stage or commit it):

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-182646/REQ-251-handback.md
```

Structure: `# REQ-NNN hand-back` with **Branch**, **Commits** (oldest first), then `## What I built`, `## File manifest` (one full path per line, `(new|modified|deleted)` + one factual line), `## P-A-U evidence`, `## Testing evidence` (real RED and GREEN output — never from a prototype or memory; the observed maintainer-verify exit code), `## Decisions (D-XX)`, `## Integration seams` (exact lines or "none"), `## Discovered Tasks`, `## Pushback`.

**Standing warning, now seven-for-seven in this repo:** every recent REQ shipped a mechanism that looked like it closed a class and closed only the instance — reviews keep finding the hole exactly where the real data lives. Assume your first fix has that shape and hunt the hole before the reviewer does.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

---
id: REQ-251
title: Retire the stale copies of the future-stamp message
status: claimed
created_at: 2026-08-18T13:55:32Z
claimed_at: 2026-08-18T18:25:40Z
route: A
status_changed_at: 2026-08-18T13:55:32Z
user_request: UR-055
addendum_to: REQ-245
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/verify_test.go
- skills/do-work-board/tools/queue-kanban/timestamp.go
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T18:26:21Z
  basis:
    - trivial short-circuit
---

# Retire the Stale Copies of the Future-Stamp Message

## What

REQ-245 rewrote the future-stamp diagnosis in five renderers. Two copies of the old wording survive outside its write set, both harmless to behaviour and both misleading to the next person who greps for the message.

## Instances

- [ ] **`verify_test.go:1186` and `:1230` hold hand-typed copies of the retired reversed-span reason as fixture literals.** They are input, never asserted against `detectCompletionAnomaly` — REQ-245's review read the surrounding asserts and confirmed they check `"REQ-9330"` and `"is earlier than claimed_at"` only. So nothing is broken. But the message they copy has now moved twice, and the next person greping for it will find two copies of a sentence that no longer exists anywhere.
- [ ] **`timestamp.go:35` — "the exact corruption the Timestamp rule warns about" — is now strictly stronger than its own test's claim.** It is largely defensible, since `formatCanonicalTimestamp` can only prevent the timezone cause; REQ-245's twin comment at `timestamp_test.go:43-46` already says so explicitly ("one of the two corruptions… This is the corruption a correct writer rules out"). Bring the source comment in line with the test that describes it.

## Requirements

- No behaviour change and no new checks. This is text.
- After the change, a grep for the retired wording across shipped source returns nothing but `do-work/` history, `CHANGELOG.md` release notes and `ai-reports/` narrative — all correctly frozen.

## Context

Both were found by REQ-245's builder and its independent reviewer, and both were deliberately left out of REQ-245 rather than widening that REQ a third time.

---

## Triage

**Route: A** - Simple

**Reasoning:** Two named stale text sites with line numbers and the replacement direction stated; no behaviour change, no new checks.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
