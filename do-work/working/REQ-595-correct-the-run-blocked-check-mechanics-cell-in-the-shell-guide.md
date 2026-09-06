---
id: REQ-595
status: claimed
domain: general
created_at: 2026-09-06T03:56:38Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
route: B
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-06T04:39:27Z
  basis:
    - Route B
    - 1-file write set
    - 3 acceptance criteria
maintenance: true
depends_on: [REQ-555]
related: [REQ-555]
write_set: [skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Correct the run-blocked-check mechanics cell, which describes shell the Go command does not run'
claimed_at: 2026-09-06T04:38:49Z
---

# Correct the run-blocked-check Mechanics Cell in the Shell Guide

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Row 13 of the "Shipped executable homes" table in
`skills/do-work/docs/prescribed-shell-primitives.md` gives `run-blocked-check`'s owned mechanics as
"GNU timeout selection and isolated stock-Bash process-group timeout/cleanup". The mechanics moved into
Go, and the Go implementation selects no external `timeout` binary and builds no stock-Bash process
group. The cell describes shell that no longer runs.

## Why

Same finding class as REQ-555: the guide is the pointer target from sixteen shipped files and it
currently describes an implementation that does not exist. REQ-555 corrected the route column; this
cell is the same defect one column over, and was equally stale before that change, so it was recorded
rather than folded in.

## Context

Found during the independent three-lens review of REQ-555, which confirmed the cell is false against
`blocked_probe_unix.go` and scored it as an adjacent finding outside that request's declared class.

## Detailed Requirements

- Read the Go implementation of the `run-blocked-check` subcommand and write the mechanics it actually
  owns. Do not paraphrase the old cell.
- Check the other thirteen Mechanics cells against their implementations in the same pass, and correct
  any that are stale. The row that was checked is not evidence about the rows that were not.
- Scope is exactly this finding class: no behaviour change, no code change, no route-column change.

## Constraints

- Guide prose only. `_dev/tests/audit-lockins.sh` already pins the route column and the orchestration
  claim from REQ-555; do not weaken either.
- If a Mechanics cell cannot be verified against code, say so in the request rather than guessing.

## Dependencies

Depends on REQ-555, which rewrote the route column of the same table.

## Open Questions

None.

## Triage

**Route: B** — Explore then build.

**Reasoning:** The edit itself is prose in one table, Route A shaped. What makes exploration real is the
request's second requirement: check the other thirteen Mechanics cells in the same pass, because the row
that was checked says nothing about the rows that were not. Each cell is a claim about what a Go
subcommand owns, and settling it means finding that subcommand's implementation and reading it. Fourteen
of those is discovery, not typing.

**Planning:** Skipped. There is nothing to sequence — the write set is one file and the work is
whatever the audit finds.

**This request exists because REQ-555 declined to widen.** The stale cell was found during REQ-555's
review and was equally stale before that change, so folding it in would have broken that request's
"scope is exactly this finding class" constraint. Same finding class, its own request.

## Plan

**Planning not required** — Route B: one file, and the work is whatever the fourteen-cell audit finds.

*Skipped by work action*

## Exploration

Three read-only agents split the fourteen rows and checked each Mechanics cell clause by clause against
the Go implementation, starting from the command dispatch and following into the package that does the
work. Full report in the run directory as `REQ-595-mechanics-audit.json`.

**Eleven cells hold. Three do not**, and the two the request did not know about fail the same way as
each other: the cell credits the subcommand with a step someone else performs.

- **`run-blocked-check` — stale, and false twice over.** "GNU timeout selection" is false: no
  `timeout` or `gtimeout` lookup exists anywhere in the module. "stock-Bash process-group" is false: the
  child is `sh -c`, and the group comes from Go's `syscall.SysProcAttr{Setpgid: true}` at
  `internal/nextselection/blocked_probe_unix.go:37-41`, with a leader check at `:47-48` and an
  in-process `time.Timer` at `:62-73`. What it does own is the group's whole lifetime — timer, TERM, a
  500ms grace, KILL, and a post-exit sweep of survivors — plus first-hand launched/timed-out facts
  rather than facts re-derived from an exit status, a bounded diagnostic identity, and the
  green/matching-red/new-red baseline comparison.
- **`install-memory-hooks` — partly stale.** "verification, and rollback" credits the command with two
  steps that are still hand-written in the knowledge actions (`memory-reference.md:163`,
  `setup-memory.md:65-68`). The Go code stops at a pre-mutation backup plus a rename. What it does own
  is per-event gating by script filename, an order-preserving settings merge that neither reorders
  consumer keys nor duplicates entries, and removal of the retired Stop hook.
- **`record-timing-event` — partly stale, and the document disagrees with itself.** "the folded
  per-request summary" is a different subcommand's work: `fold-timing-summary` does the folding, and the
  prose lower on the same page already says so. `record-timing-event` appends one bounded JSON record to
  a Git-private per-request stream and never touches the request file. It also owns something the cell
  misses: elapsed seconds are derived either from an explicit start or by chaining to the previous
  event's end.

**Two stale prose claims outside the table, in the same document and the same class**, verified against
`internal/toolboxcommands/report_image.go`:

- Line 138: "Publication happens once, as a single same-filesystem rename of the complete verified
  batch", and the two later clauses that depend on the rename. There is no rename. The command claims
  `generated/` with an exclusive `mkdir` (`:193`), writes each verified image into that claimed
  directory one at a time (`:212-224`), and removes the whole directory on any failure through a
  deferred identity-checked `RemoveAll` (`:197-208`). The all-or-nothing outcome survives; the mechanism
  named does not.
- Line 136: the staging directory is described as "adjacent to `generated/`". It is
  `os.MkdirTemp("", "do-work-generated-staging.*")` at `:164` — the system temporary directory. The
  bare-filename requirement the sentence exists to justify still holds; its stated reason does not.

**Two accurate cells understate what they own**, recorded and not changed because the request's class is
stale prose, not incomplete prose: `atomic-download` also validates the URL scheme and injects a Bearer
header from `GH_TOKEN`/`GITHUB_TOKEN`, and `capture-screenshot`'s publication is a hard link with an
inode-identity confirmation.

*Generated by three Explore agents*

## Scope

**Files I will touch:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — three Mechanics cells rewritten to what the Go code owns, and two stale mechanism claims in the report-image-batch section corrected

**One widening inside the declared file, and it is a judgment call stated rather than hidden.** The
request's Detailed Requirements name the Mechanics cells. The audit also found two false mechanism
claims in the prose of the same document — a publication described as a same-filesystem rename that is
an exclusive `mkdir` plus per-file writes, and a staging directory described as adjacent to `generated/`
that is in the system temporary directory. They are the same finding class as the cell (the guide
describes an implementation that does not exist), they are in the file already declared, and they were
found by the same pass. Leaving a known-false mechanism claim in place to protect a scope boundary
would ship the defect this batch exists to remove. No new file is touched.

**Files I will NOT touch:** any Go source — this is prose only, no behaviour change and no code change,
as the request requires. The route column, which REQ-555 settled and which
`_dev/tests/audit-lockins.sh` pins. The two accurate-but-understated cells, because incomplete is not
the class this request removes.

**Acceptance criteria:**
- [ ] The `run-blocked-check` cell names the mechanics the Go implementation owns, with no external
  timeout binary and no Bash process group
- [ ] The `install-memory-hooks` and `record-timing-event` cells stop crediting the subcommand with a
  step another subcommand or a hand-written action step performs
- [ ] The report-image-batch prose describes the publication mechanism that exists
- [ ] Every remaining cell was checked against its implementation and the eleven that hold are named
- [ ] `_dev/tests/audit-lockins.sh` and `_dev/tests/prescribed-shell-canonicalization.sh` both still exit 0
