---
id: REQ-555
title: '[impact-negligible] Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-06T03:17:52Z
  basis:
    - trivial short-circuit
    - Route A
    - 3-file write set
depends_on: [REQ-554]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
route: A
write_set: [skills/do-work/docs/prescribed-shell-primitives.md, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-06T03:16:18Z
---

# Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form

## What
The "Shipped executable homes" table in `skills/do-work/docs/prescribed-shell-primitives.md` assigns owned mechanics to nine `*.sh` paths that are each a 6-to-11-line `exec` shim over `do-work-cli.sh` (the mechanics moved to Go at 0.260.1), and one sentence below it says `scripts/protected-inventory.sh` "orchestrates" two check scripts, which a six-line shim cannot do. Reword the route column to the `tools/do-work-cli.sh … <subcommand>` form the toolbox rows already use and delete the orchestration sentence.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` (§ Unchecked Exit Status Reads as
  Content, § Closed Enumerations Go Stale) and re-verified every claim the request makes against
  HEAD before choosing Route A. Approach: take each route from the shim's own `exec` line, delete
  the false orchestration clause, and pin both shapes with one assertion block that fails loudly
  when its target is renamed.
- [x] **[APPLY]:** Two files changed, both in the declared write set:
  `skills/do-work/docs/prescribed-shell-primitives.md` and `_dev/tests/audit-lockins.sh`. Nothing
  else was touched. `_dev/tests/prescribed-shell-canonicalization.sh` was declared and turned out
  not to need a change; that is stated in the summary rather than papered over with an edit.
- [x] **[UNIFY]:** `git diff --stat` — `_dev/tests/audit-lockins.sh` +53, the guide +12/-10. Both
  files read in full. The guide: fourteen table rows now read the same way, each subcommand
  checked against its shim's `exec` line, and the added sentence states a condition rather than a
  count. The lock-in: every command's exit status is read, the awk pass exits 3 when the heading
  is gone and that status is checked, and `rg`'s 0/1/>1 are told apart. No debug artifacts, no
  commented-out code, no `TODO`. `bash _dev/tests/audit-lockins.sh` exits 0 and
  `bash _dev/tests/prescribed-shell-canonicalization.sh` exits 0; four ablations each print the
  intended FAIL line and exit 1.

## Why
The guide is the pointer target from 16 shipped files (ratchet-enforced) and it currently misroutes readers to shims and describes an orchestration that no longer exists.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 7, sweep_key `stale-shell-ownership-prose`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -5. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- Seven rows naming `scripts/show-commit-diff.sh`, `add-local-git-exclude.sh`, `atomic-download.sh`, `capture-screenshot.sh`, `run-blocked-check.sh`, `protected-inventory.sh`, `stage-exact-deletion.sh` (9-11 lines each, all `exec … do-work-cli.sh`).
- Two rows naming `../../do-work-knowledge/scripts/lexical-memory-recall.sh` and `install-memory-hooks.sh` (6 lines each).
- The sentence `which orchestrates \`tools/checks/uncommitted-inventory.sh\` and \`tools/checks/associate-files.sh\`` is false at HEAD: delete it.
- The shims themselves are not touched (retained by the 0.260.1 decision); only the guide's description of them changes.
- Reproduce at dc8a64e3: `awk 'NR>=9 && NR<=22 && /\.sh`/' skills/do-work/docs/prescribed-shell-primitives.md | grep -oE '[^`]*\.sh' | while read -r p; do case "$p" in ../../*) fp="skills/${p#../../}";; *) fp="skills/do-work/$p";; esac; echo "$(wc -l < "$fp" | tr -d ' ') lines $fp $(grep -c do-work-cli "$fp") cli-exec"; done; rg -n 'which orchestrates' skills/do-work/docs/prescribed-shell-primitives.md`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Same guide and ratchet as REQ-554: land after it so the heading and pointer counts are re-baselined once.
- Lock-in limit: shim rows in the executable-homes table: 0 after this REQ (today 9).

## Dependencies
Depends on REQ-554, which already edits this guide and re-baselines the ratchet, so this REQ is a table rewrite only.

## Builder Guidance
Firm: the route column names the CLI subcommand, not a shim.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints nine shim rows (6-11 lines each, all `do-work-cli` execs) and the orchestration sentence.
**GREEN when:** The table has zero `.sh` rows for mechanics owned by Go, the sentence is gone; the lock-in pins shim rows in that table at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for stale-shell-ownership-prose.*

## Triage

**Route: A** — Build directly.

**Reasoning:** Every claim the request makes was re-checked against HEAD and all of them hold, so there
is nothing left to discover. The nine paths in the "Shipped executable homes" table are 6-to-11-line
`exec` shims over `tools/do-work-cli.sh`, each naming its own subcommand on that exec line, so the
route column can be rewritten mechanically from the shims themselves. The orchestration sentence is
still present, and it is still false: `scripts/protected-inventory.sh` is six lines, and the two files
it is said to orchestrate — `tools/checks/uncommitted-inventory.sh` and
`tools/checks/associate-files.sh` — are themselves 9-line and 19-line compatibility launchers over the
same Go command, which no code in `internal/` ever launches. The work is one table rewrite, one
sentence deletion and one counting assertion.

**Planning:** Skipped.

**REQ-554 landed first, as this request requires.** It rewrote the paragraph the orchestration sentence
sits in and left that sentence untouched by an explicit decision, so this request edits settled text
rather than racing it.

## Plan

**Planning not required** — Route A: the write set, the target table and the sentence to delete are all
named by the request, and every claim was re-checked against HEAD before the route was chosen.

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `_dev/tests/audit-lockins.sh` (modified)

**What was done:** The nine `.sh` rows in the "Shipped executable homes" table now name the
`tools/do-work-cli.sh … <subcommand>` route, taken from each shim's own `exec` line rather than
guessed, so all fourteen rows read the same way. The false orchestration clause is deleted. One
assertion block in `_dev/tests/audit-lockins.sh` pins both shapes.

**One sentence was added, and it is earned.** Rewriting the route column alone leaves the guide
contradicting itself: the table would name only CLI subcommands while the paragraph below it still
says the inventory ships behind `scripts/protected-inventory.sh`, which is true — `actions/commit.md`
and `../../do-work-toolbox/actions/inspect.md` invoke that path today. A reader following the table
would conclude the launchers are gone and could delete one. The added sentence states the condition
(where a route also ships a retained launcher of the same name) rather than listing which nine, so it
cannot go stale as rows move, and it says plainly that a behaviour change is made in the command and
never in the launcher.

**The expected net line delta was -5; the actual guide delta is +2** (12 insertions, 10 deletions).
The nine rewritten rows are net zero, the deleted clause was part of a longer line, and the added
sentence is +2 with its blank line. The estimate assumed a pure deletion; the correctness of the guide
was preferred to the number.

`_dev/tests/prescribed-shell-canonicalization.sh` is in the declared write set and was not changed: it
checks only that the nine shims exist and are executable, which this request does not touch. Declaring
a file and then not needing it is recorded rather than papered over.

## Decisions — implementation

- **D-01 — the route column names the subcommand each shim actually execs. DECIDE & STATE.** Every one
  of the nine was read: `show-commit-diff`, `add-local-git-exclude`, `atomic-download`,
  `capture-screenshot`, `run-blocked-check`, `protected-inventory`, `stage-exact-deletion`,
  `lexical-memory-recall`, `install-memory-hooks`. None was inferred from the file name.
- **D-02 — the retained launchers are stated once, keyed on a condition, not counted. DECIDE &
  STATE.** "Where a route above also ships a retained `scripts/*.sh` launcher of the same name" holds
  as rows are added or removed. A sentence saying "the first nine rows" would be a hand-maintained
  list, which `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale bans.
- **D-03 — the shims themselves are untouched, and so is the sentence that routes readers to one.**
  The 0.260.1 decision retained them and shipped actions still invoke them by path. Rewriting those
  invocations is a different change with a different blast radius, and this request does not name it.
- **D-04 — the lock-in pins two shapes in one block because they are one defect. DECIDE & STATE.** A
  shim in the route column and the claim that a six-line launcher orchestrates two scripts are both
  "the prose describes a shell ownership that no longer exists". The request asks for one assertion,
  and splitting them would have produced two.
- **D-05 — a missing guide or a renamed heading fails rather than reads clean. DECIDE & STATE.** The
  awk pass exits 3 when it never sees the heading, and that status is checked. A ratchet that goes
  quiet when its target is renamed is the failure mode REQ-552 and REQ-554 both hit in this same file.

## Qualification

**Passed.** Read from the range `a49a542f..ad8a8050`, two files, 65 insertions and 10 deletions.
Canonical `qualify` is satisfied.

- **Every claim in the request holds at HEAD, and each was checked rather than assumed.** The nine
  paths are 9, 10, 9, 9, 11, 6, 9, 6 and 6 lines, and every one of them `exec`s
  `do-work-cli.sh --format text <subcommand> "$@"`. The orchestration clause was still present. This is
  the opposite of the sibling REQ-556, whose stated baseline did not survive contact with HEAD.
- **The orchestration clause is false in two independent ways, both verified.**
  `scripts/protected-inventory.sh` is six lines, and `tools/checks/uncommitted-inventory.sh` and
  `tools/checks/associate-files.sh` are themselves 9-line and 19-line compatibility launchers over the
  same command. Nothing under `internal/` launches either of them: the only matches in the Go tree are
  in `inventory_test.go`, which drives the launcher chain deliberately to hold the launcher contract.
- **One sentence was added, and the addition is declared rather than smuggled.** `maintenance.md` asks
  for a concrete case that fails without it, and that case is the paragraph immediately below the table,
  which still routes readers to `scripts/protected-inventory.sh` because that is genuinely what
  `commit.md` and `inspect.md` invoke. Without the sentence the table and the paragraph contradict each
  other. It is stated in the Implementation Summary as an addition and is the reason the net line delta
  is +2 rather than the expected -5.
- **The declared write set is a ceiling that was not filled.**
  `_dev/tests/prescribed-shell-canonicalization.sh` was declared and not changed, because it only
  asserts the nine shims exist and are executable, which this request does not touch. It still exits 0.

## Testing

**The lock-in was proven red four ways before it was accepted**, and the guide was restored to the
green state after each.

- One shim row returned: `prescribed-shell-primitives.md:9 routes owned mechanics to a shim
  (\`scripts/show-commit-diff.sh\`)`, exit 1.
- All nine returned: nine `FAIL:` lines, one per row, exit 1. The count is what the request pins at 0.
- The orchestration clause returned, reworded to a shorter sentence than the deleted one: caught, exit
  1. The check matches the claim, not the original wording.
- The table heading renamed from `## Shipped executable homes` to `## Executable homes`: `the "##
  Shipped executable homes" heading is gone ... (awk exit 3); the route ratchet cannot run`, exit 1.
  This is the case that matters most — a rename is how a ratchet goes quiet instead of red.

Green at HEAD: `bash _dev/tests/audit-lockins.sh` prints `Audit lock-in regressions passed.` and exits
0; `bash _dev/tests/prescribed-shell-canonicalization.sh` prints its passed line and exits 0. The full
fast gate is run once for the batch rather than per request.
