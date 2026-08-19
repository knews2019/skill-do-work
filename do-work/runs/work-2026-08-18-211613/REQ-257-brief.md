# REQ-257 builder brief

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-257-decide-whether-the-repairer-learns-offset-and-fractional-stamps` — this is your working directory; everything you do happens here.
**Branch:** `worktree-agent-REQ-257-decide-whether-the-repairer-learns-offset-and-fractional-stamps` (already checked out in that worktree).
**Hand-back file (write when done):** `/home/user/skill-do-work/do-work/runs/work-2026-08-18-211613/REQ-257-handback.md` — absolute main-tree path, the ONE main-tree path you may write. Never stage or commit it.

## The REQ (verbatim)

```markdown
---
id: REQ-257
title: Decide whether the timestamp repairer learns offset and fractional stamps
status: claimed
claimed_at: 2026-08-18T21:16:24Z
route: B
created_at: 2026-08-18T17:49:24Z
status_changed_at: 2026-08-18T20:55:14Z
user_request: UR-056
addendum_to: REQ-246
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-255]
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route B
    - 2-file write set
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Decide Whether the Timestamp Repairer Learns Offset and Fractional Stamps

## What

REQ-246's repairer deliberately refuses stamps with numeric UTC offsets (`2093-01-01T00:00:00+02:00`) or fractional seconds — repairing them needs timezone arithmetic, and a wrong guess would rewrite a correct stamp (REQ-246 D-04, documented in the script header). The board and forensics still detect and warn on those shapes, so they remain a detection-without-repair residual. **Correction (REQ-255 review, I-2):** these are no longer the *only* such residual — a quoted stamp with padding inside the quotes is also board-parseable and refused here; that one is tracked in REQ-267. This asks whether that residual matters enough to implement offset arithmetic in `comparison_key_for`, or whether the documented refusal is the permanent answer.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Builder-discovered on REQ-246 (Discovered Tasks, first item), classified [normal]. Gated behind REQ-255 so the repairer's shape handling settles once, on top of the parity sweep, rather than twice.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-246: the repairer refuses offset/fractional stamps that the board still warns about — implementing offset arithmetic in `comparison_key_for` would close the last board-detectable-but-unrepaired timestamp class. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the documented refusal (D-04) is the permanent answer and the board's warning is disclosure enough.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (decide, then either implement offset/fractional arithmetic or make the refusal permanent and provable) but grounding that decision needs discovery: how `comparison_key_for` and `extract_timestamp_fields` currently recognize shapes after REQ-255's parity sweep, and what the read-side detectors accept.

**Planning:** Not required

---

## Exploration

`skills/do-work/scripts/repair-req-timestamps.sh` holds the whole shape surface in two functions the REQ names:

- `comparison_key_for` (line ~206) — turns a value token into a sortable key; ends by gating on `calendar_components_valid` (line ~179). This is the single place a shape is recognized or refused, and REQ-247's auditor sources this file, so widening it widens the auditor in the same edit.
- `extract_timestamp_fields` (line ~259) — whitespace-token extraction; REQ-255 taught it the quoted and unquoted space-separated spellings.
- The script header (line ~51) states the current refusal: *a numeric UTC offset or fractional seconds is not provably wrong without timezone arithmetic* (REQ-246 D-04).

Lock-ins live in `_dev/tests/prescribed-shell-scripts-behavior.sh` as a `# repair-req-timestamps:` comment-headed case group (lines ~1098–1500), including the two REQ-255 space-separated cases and the skew-constant lock-step case at ~1474 — the pattern any new case follows.

*Explored inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — the decision's implementation: either offset/fractional recognition in `comparison_key_for`, or the refusal made permanent and provable
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-ins for whichever answer wins

**Files I will NOT touch:** `skills/do-work/scripts/audit-archive-timestamps.sh` (it sources the repairer, so it inherits the change without an edit), the SessionStart hook, the board tool, `CHANGELOG.md`, `VERSION`, `skills/do-work/actions/version.md`.

**Acceptance criteria (restated from REQ):**
- [ ] The offset/fractional residual is decided either way, and the decision is stated where a reader of the script meets it
- [ ] Whichever way it goes, a lock-in pins it — a refusal that is only a comment is not pinned
- [ ] The one code body stays shared: nothing is duplicated into the auditor
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0
```

## Never touch

- Anything under `do-work/` — the queue, `working/`, `CHECKPOINT.md`, `runs/` — is the orchestrator's. Your worktree carries a stale committed snapshot of it; treat it as absent. The single exception is your own hand-back file at the absolute path above, which you write but never stage or commit.
- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` — serial-only files owned by the integrator. Bumping any of them races every sibling builder. Skip the "Before Every Commit" ritual in `CLAUDE.md` entirely; it belongs to the integrating commit, not to yours.
- Any file outside your REQ's `## Scope` "Files I will touch" list. Discovering mid-build that you need one is a **stop-and-report to the orchestrator, never a silent write** — put it in the hand-back and say why.
- The main tree at `/home/user/skill-do-work` — you are not in it and must not write to it.

## Environment

- Go 1.26.1 at `/usr/local/go`, ShellCheck 0.11.0, `just` 1.21.0, Node 22 are all installed and are the versions the gate pins.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB gitignored binary into the source tree. Build to a scratch directory with `-o`.
- The canonical gate is `bash _dev/tests/maintainer-verify.sh`, run from your worktree root. **Exit code zero is the only proof it passed** — never pipe it through `tail` or `head`, because the pipeline's exit status hides the failure. Redirect to a file and echo `$?` if you want to read the output.
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry a timestamp forward and never compute one.

## How to work

Follow the do-work crew rules that always load during implementation — read them from your worktree:

- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- plus `skills/do-work/crew-members/testing.md` if this REQ has `tdd: true`, and `skills/do-work/crew-members/maintenance.md` if it has `maintenance: true`.

Read every path in the REQ's `prime_files` before touching code, including its `## Lessons` links.

**P-A-U is mandatory and it is your own note-taking** — the REQ file itself lives in the main tree and is not yours to edit, so write the [PLAN] / [APPLY] / [UNIFY] evidence into your hand-back instead. [PLAN] before any code. [APPLY] stays inside the declared scope. [UNIFY] runs `git diff --stat` against your branch point, runs the native linters, verifies no debug artifacts, and lists each file you checked.

**The failure mode that beat nine of twelve REQs last session:** a mechanism that looks like it closes a class and closes only the instance. The two REQs that beat it grepped or fuzzed the **primitive** before declaring the class closed, and both found the real hole where no instance list pointed. Assume your first fix has that shape. An instance list is a sample, not the class.

**Commit in small, individually-green increments on your branch.** Builders were repeatedly killed mid-run by server-side errors last session and nothing was lost, because each increment was already committed. Do not wait until the end to commit.

## Hand-back format

Write `REQ-NNN-handback.md` at the absolute path given above, with these sections:

- `**Branch:**` the operative name, and the commit it was cut from
- `**Commits (oldest first):**` short hash + subject for each
- `## What I built` — factual, what you actually built, not what the REQ asked for
- `## File manifest` — every file created/modified/deleted with `(new)` / `(modified)` / `(deleted)` and a phrase on what changed in it
- `## P-A-U evidence` — [PLAN] / [APPLY] / [UNIFY], with what you actually ran
- `## Testing evidence` — real observed output, never a prototype or a paraphrase. If the REQ is `tdd: true`, the RED must be a real run against pre-change code, quoted. State the gate's exit code as a number you observed.
- `## Decisions (D-XX)` — each marked DECIDE & STATE or ESCALATE; ESCALATE entries carry `Value:` and `Risk:` lines
- `## Integration seams` — exact lines and where they go, for any shared file you must not edit yourself. `None.` if there are none.
- `## Discovered Tasks` — out-of-scope finds, not fixed inline
- `## Pushback` — where the REQ was wrong, or where you disagree. `None.` if you have none.

Your report is a claim, not evidence: the orchestrator judges from git state. Say plainly what you did not do.
