# REQ-246 builder brief

**Route C.** Estimated 45 active minutes (P50, medium confidence). Plan and scope are in the REQ body below — they are the orchestrator's, and you may argue with them in your hand-back, but do not silently depart from them.

**Why this one matters:** detection of a wrong timestamp already exists on the read side; repair depended on the same agent that wrote the bad stamp. This closes that loop mechanically, so a fabricated future stamp gets corrected regardless of agent compliance.

## How this build runs

You are a **worktree builder** dispatched by the do-work work pipeline (`skills/do-work/actions/work.md`). Read that file's Step 6 expectations if you need them; everything binding is repeated here.

**Your tree, your branch.** Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-246-repair-wrong-queue-and-working-timestamps-from-the-session-hook`. That is a full checkout of this repository on branch `worktree-agent-REQ-246-repair-wrong-queue-and-working-timestamps-from-the-session-hook`, cut from the integration tip `67dae6b`.

- Never write anything under `/home/user/skill-do-work` — the main tree belongs to the orchestrator. The one exception is your hand-back file, named below.
- Never read or write anything under `do-work/` in your own worktree. Your worktree carries a **stale snapshot** of the queue; the live queue lives in the main tree and is the orchestrator's alone. Your REQ body is inlined in this brief — that is your copy of it.
- Commit your work on your own branch, as many commits as the work naturally splits into. Do not bump `VERSION`, do not touch `CHANGELOG.md`, do not touch `skills/do-work/actions/version.md` — those are serial-only files the integrator owns, and a builder bumping them races every sibling.
- If you need one line of wiring in a file outside your write set — a shared registry entry, a pointer in someone else's doc — **do not edit it**. Hand back the exact line and where it goes as an *integration seam*; the orchestrator applies it inside the merge commit.
- If you discover you need a file outside your declared scope, stop and report it in your hand-back rather than writing it silently.

**Crew rules load from your own worktree** (they ship, so they are there at the same paths): read `skills/do-work/crew-members/general.md`, `skills/do-work/crew-members/coding-guardrails.md`, and `skills/do-work/crew-members/communication-style.md` before you write code. This REQ is `tdd: true`, so also read `skills/do-work/crew-members/testing.md` and follow RED → GREEN → REFACTOR. Read every path in the REQ's `prime_files` too — those primes encode prior mistakes in this exact area.

**P-A-U phasing is mandatory.** Your REQ body carries an `AI Execution State (P-A-U Loop)` block. Work it: [PLAN] a brief technical approach, [APPLY] inside declared scope only, [UNIFY] run `git diff --stat`, run the linters, verify no debug artifacts, and list each file you checked. The orchestrator audits the checked boxes against the diff — a checked [UNIFY] over a diff containing stray instrumentation is a qualification failure. Record the P-A-U evidence in your hand-back, not in a `do-work/` file.

**Log significant decisions as D-XX** in your hand-back with reasoning. A reversible, low-reach choice is DECIDE & STATE (reasoning only); an irreversible, taste-dependent, or genuinely contestable one is ESCALATE — add `Value:` and `Risk:` lines.

**Out-of-scope finds** go in a `## Discovered Tasks` list in your hand-back. Do not fix them inline.

## Environment notes for this checkout

This is a Linux container, not the maintainer's machine. Before you start, know:

- `bash _dev/tests/maintainer-verify.sh` is the canonical pass/fail gate and it **exits 0 at your branch point** — that is your baseline. Exit code zero is the only proof; never pipe it through `| tail` or judge it from a summary. It takes a few minutes.
- The toolchain was installed for this session: Go 1.26.1, ShellCheck 0.11.0, `just` 1.21.0. They are on `PATH`.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB `queue-kanban` binary into the source tree, which is gitignored (so nothing warns you) and multiplies through the installer probe's copies. Build with `go build -o /tmp/<name> .`.
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry a timestamp forward and never compute one.
- Browser tooling may drop a `.playwright-mcp/` directory into the main tree without you issuing a write. It is gitignored and holds sibling agents' evidence — if you clean it, remove only your own files.

## Hand-back

When you are done, write your report to exactly this absolute path:

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-162355/REQ-246-handback.md
```

That is the one main-tree path you may write. Never stage it, never commit it — it is an orchestrator working file, not branch content.

Structure it as:

```markdown
# REQ-NNN hand-back

**Branch:** worktree-agent-REQ-246-repair-wrong-queue-and-working-timestamps-from-the-session-hook
**Commits:** <short hashes, oldest first>

## What I built
## File manifest
- `path` (new|modified|deleted) — one factual line each
## P-A-U evidence
## Testing evidence
Red-green: the test name, the failure text before, the pass after. Quote real output — never a transcript from a prototype or from memory.
Full gate: the `maintainer-verify.sh` exit code you actually observed.
## Decisions (D-XX)
## Integration seams
Exact lines and where they go, or "none".
## Discovered Tasks
## Pushback
Anything in this brief you think is wrong, with your evidence.
```

**One standing warning, from the previous session's own record.** Five consecutive REQs shipped a mechanism that looked like it closed a class and closed only the instance — and in three of five the remaining hole was exactly where the real data lives. Assume your first fix has that shape and go looking for the hole before a reviewer does. A passing assertion is not evidence about the thing you did not sample.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

---
id: REQ-246
title: Repair detectably wrong queue and working timestamps from the session hook
status: claimed
created_at: 2026-08-18T12:38:26Z
claimed_at: 2026-08-18T16:09:27Z
route: C
user_request: UR-056
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-247, REQ-244, REQ-245]
batch: timestamp-stamping-integrity
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-08-18T16:10:30Z
  basis:
    - Route C
    - 4-file write set
    - 2 new files
    - 2 subsystems involved
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- skills/do-work/hooks/session-start.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
- _dev/tests/session-start-hook-behavior.sh
---

# Repair Detectably Wrong Queue and Working Timestamps From the Session Hook

## What

A core script that scans REQ files in `do-work/queue/` and `do-work/working/` for detectably wrong `*_at` stamps and rewrites them with a mechanically derived correct value — no agent judgment anywhere in the path. Wired into the SessionStart hook the way `scripts/cleanup-req-reservations.sh` already is, so repair happens before any agent or board render sees the file.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Detection already exists read-side (board generate/verify, forensics Check 11) but repair depended on the same agent that wrote the bad stamp. This closes that loop mechanically: a fabricated future stamp gets corrected regardless of agent compliance.

## Detailed Requirements

- **Detect only provable wrongness:** stamps later than now plus the 2-minute skew allowance (matching `actions/forensics.md` Check 11 and the board's `futureTimestampSkewAllowance` in `model.go`), and impossible orderings (`completed_at` earlier than `claimed_at`, `claimed_at` earlier than `created_at`).
- **Replacement value, by file state:** file dirty against HEAD → file mtime (the actual write instant); committed content → author time of the commit that introduced the stamp line. Clamp so the repaired set satisfies `created_at ≤ claimed_at ≤ completed_at ≤ now`.
- **Guard style of `tools/checks/record-commit-hash.sh`:** verify before replace, atomic write, a tripped guard leaves the file byte-identical, nonzero exit on any failure.
- **Audit trail:** print each correction (file, field, old value, new value, replacement source) to the hook output. No new frontmatter fields. (Log-only was the capture recommendation; the user did not override it — builder may surface it for confirmation if the hand-back warrants.)
- **Hook wiring:** invoked from the SessionStart hook (`skills/do-work/hooks/`) alongside `cleanup-req-reservations.sh`; also invocable directly.
- **Scope is queue + working only** (user-confirmed). Archived files are REQ-247's territory.

## Constraints

- Not in the board tool: its frontmatter surface is documented read-only (`frontmatter_cli.go` ~line 38) and its write-surface count is pinned in CLAUDE.md § Kanban Board Write Surfaces.
- No LLM, no heuristic text parsing beyond frontmatter field extraction.
- Platform floor applies: POSIX shell, no Go build in the repair path.
- Provenance: replaces the Discuss verdict on validate-feedback Finding 2 (UR-055 triage); Surface-cost — earned by the fabricated-stamp incident report, and cheaper than the alternative (per-site guard prose scattered across write sites, which the Timestamp rule's centralization clause forbids).

## Red-Green Proof

**RED prompt/case:** A `_dev/tests/` lock-in test builds a fixture REQ with `created_at` two hours in the future in a scratch queue, runs the repair script, and asserts the stamp is rewritten to the fixture file's mtime and the correction is logged. Fails today because the script does not exist.
**Why RED now:** Nothing repairs a detected-wrong stamp; the board only warns.
**GREEN when:** The test passes, an ordering-violation fixture is likewise repaired, a clean fixture passes through byte-identical, and `bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** User confirmed (scope via ask tool; repair-window and mtime sources are the user's own design)

## Full Context

See `do-work/user-requests/UR-056/input.md` for complete verbatim input.

---
*Source: "basically if a wrong timestamp is detected, it should be automatically mechanically no-llm corrected by the tool" + ask-tool answer "Queue + working (Recommended)"*

---

## Triage

**Route: C** - Complex

**Reasoning:** A new shipped script with its own detection, derivation and clamping logic, plus hook wiring and a lock-in test suite — multiple components, many precise requirements, and a guard-style contract to match.

**Planning:** Required

## Plan

**Approach**

1. **One script, two entry conditions.** `skills/do-work/scripts/repair-req-timestamps.sh [project-root]` scans `do-work/queue/` and `do-work/working/` only. It runs from the SessionStart hook and stands alone as a direct invocation; the hook guards it exactly the way it guards `cleanup-req-reservations.sh`, so a partial install or a repairer failure can never break the status banner.

2. **Detection is arithmetic, never judgment.** Two predicates, both mechanical: a stamp later than `now + 120s` (the same 2-minute allowance `actions/forensics.md` Check 11 and the board's `futureTimestampSkewAllowance` use), and an impossible ordering among `created_at`, `claimed_at`, `completed_at`. Nothing else is a defect — a stamp that is merely surprising is left alone.

3. **Replacement source is decided by file state, not by preference.** Dirty against HEAD → the file's mtime, which is the actual write instant. Committed content → the author time of the commit that introduced that stamp line. Then clamp the repaired set so `created_at ≤ claimed_at ≤ completed_at ≤ now`.

4. **Guard style is `tools/checks/record-commit-hash.sh`'s, and that is a hard requirement, not a stylistic nod.** Verify before replace; write atomically; a tripped guard leaves the file byte-identical; any failure exits nonzero. Free-form frontmatter edits truncated six archived REQs to 0 bytes in a consumer repo — that is the incident this style exists to prevent.

5. **Tests are lock-in, not decoration.** The `## Red-Green Proof` cases go into `_dev/tests/prescribed-shell-scripts-behavior.sh` (future stamp repaired to mtime, ordering violation repaired and clamped, clean fixture byte-identical, tripped guard byte-identical and nonzero) and the wiring probe into `_dev/tests/session-start-hook-behavior.sh`. Each case names the failure it pins.

**Ordering:** script first with its cases RED, then the hook wiring, then the probe.

**Watch for:** the previous session's recurring finding — a mechanism that looks like it closes a class and closes only the instance. Here the instance is *a future `created_at` in the queue*; the class is *every detectably wrong `*_at` in either directory, under both replacement sources*. Test the derivation you did not reach for first.

*Written inline by the orchestrator (no separate Plan agent) — Route C.*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (new) — the mechanical repairer: detection, replacement derivation, ordering clamp, atomic guarded write, audit trail
- `skills/do-work/hooks/session-start.sh` (modify) — invoke the repairer alongside `cleanup-req-reservations.sh`, guarded the same way
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — the lock-in cases named in `## Red-Green Proof`
- `_dev/tests/session-start-hook-behavior.sh` (modify) — the hook-wiring probe

**Files I will NOT touch:**
- `skills/do-work-board/tools/queue-kanban/**` — the board's frontmatter surface is read-only and its write-surface count is pinned (REQ Constraints).
- Any shipped action markdown under `skills/*/actions/` — REQ-249 is sweeping that tree in this same wave. If wiring needs documenting there, hand the exact lines back as an integration seam.
- `do-work/archive/**` — REQ-247's territory.

**Acceptance criteria (restated from REQ):**
- [ ] A fixture REQ with `created_at` two hours in the future is rewritten to the fixture file's mtime, and the correction is logged.
- [ ] An ordering-violation fixture (`completed_at` < `claimed_at`, or `claimed_at` < `created_at`) is likewise repaired, clamped so `created_at ≤ claimed_at ≤ completed_at ≤ now`.
- [ ] A clean fixture passes through byte-identical.
- [ ] A tripped guard leaves the file byte-identical and exits nonzero.
- [ ] The repairer runs from the SessionStart hook and can also be invoked directly.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 (recorded in `do-work/working/baseline.json`)
**Dependencies:** ⚠ this checkout needed Go 1.26.1, ShellCheck 0.11.0 and `just` installed before the baseline could run at all, and one pre-existing Linux-only test failure had to be fixed first (0.212.8) — see the REQ brief.

*Checked by work action*
