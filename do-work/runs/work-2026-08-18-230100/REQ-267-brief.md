# REQ-267 builder brief

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-267-close-the-two-remaining-repairer-shape-divergences` — your working directory.
**Branch:** `worktree-agent-REQ-267-close-the-two-remaining-repairer-shape-divergences` (already checked out there).
**Hand-back (write when done):** `/home/user/skill-do-work/do-work/runs/work-2026-08-18-230100/REQ-267-handback.md` — absolute main-tree path, the ONE main-tree path you may write. Never stage or commit it.

## The REQ (verbatim)

```markdown
---
id: REQ-267
title: Close the two remaining repairer shape divergences
status: claimed
claimed_at: 2026-08-18T22:59:48Z
route: B
created_at: 2026-08-18T21:03:15Z
status_changed_at: 2026-08-18T22:20:09Z
user_request: UR-056
addendum_to: REQ-255
domain: general
review_generated: true
sweep: true
sweep_key: repairer-detector-shape-parity
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T22:59:48Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Close the Two Remaining Repairer Shape Divergences

## What

REQ-255 closed six shape divergences between the timestamp repairer and the board's readers. Its independent review, fuzzing the whole value space, found two more — both pre-existing at REQ-255's branch point, neither among its declared instances, and neither corrupting.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] **An unterminated frontmatter block is repaired where the board sees only body text (reproduced by execution).** A file whose opening fence never closes is scanned to end-of-file by the extractor and its stamp lines are rewritten by the unattended hook, while the board's `splitFrontmatter` returns *no frontmatter* for exactly that shape. The script's own scope comment states the fence-bounded contract the code does not honour here. **Second face of the same root cause:** when such a file also ends with the defective stamp on its final line with no trailing newline, the changed-line guard expects four diff lines and sees two, so the repair is refused and the SessionStart hook exits 1 **every session, permanently**, with no self-heal. Refusing the fence-broken shape the way the read side does closes both at once.
- [ ] **Quoted stamps with padding inside the quotes are board-parseable but refused, and the refusal is undocumented (reproduced by execution).** A Go probe replicating the board's pipeline (YAML unquote, then trim, then parse) accepts `"2093-01-01 00:00:00 "` and would flag it future; the repairer refuses it byte-identical and the refusal falls into the header's catch-all "anything else unparseable", which is false for this shape. The header's parity rule claims the opposite family-wide, so the documentation is wrong rather than merely silent.

- [ ] **The clean-fixture comment still justifies the offset refusal with the reason REQ-257 repudiated (added from REQ-257's review, Important 3).** `_dev/tests/prescribed-shell-scripts-behavior.sh:1150` says the repairer must not touch "a numeric-offset value **it cannot compare without timezone arithmetic**" — verbatim the justification REQ-257's header now states is wrong ("the arithmetic is the risk, not the obstacle"). Same family as instance 2: a statement about what the repairer refuses that is wrong rather than merely silent. This file is already in this REQ's write set.

## Requirements

- Each instance is repaired to canonical form or refused byte-identical **with the refusal documented beside the existing refusal entries** — the never-half-rewrite rule from REQ-255 still governs.
- The permanent hook-failure loop is gone: no input shape may leave the SessionStart hook exiting nonzero every session with no path to self-heal.
- Lock-in cases per instance, through both scan scopes.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-255's independent review, findings I-1 and I-2 (gate: trivial each — never-corrupt holds, and both shapes are already malformed or exotic). Created `pending-answers` per the generation-≥2 cascade stop. The review also noted a consequence for a queued sibling: REQ-257's description called offset and fractional seconds "the last" board-detectable-but-unrepaired class, which instance 2 makes inaccurate — corrected in that REQ directly.

## Open Questions

- [x] REQ-255's review found two more shapes where the repairer and the board disagree — one repairs a file the board reads as having no frontmatter (and, in one variant, makes the session hook fail forever), the other refuses a shape the board accepts while the docs claim it is handled. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`, presented as the only known live defect in the queue that can wedge a session permanently. Nothing was put out of scope — both instances stand, and the never-half-rewrite rule from REQ-255 still governs the refusal path.
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
