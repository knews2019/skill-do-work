# REQ-260 builder brief

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-260-gofmt-the-durations-day-domain-truncation` — this is your working directory; everything you do happens here.
**Branch:** `worktree-agent-REQ-260-gofmt-the-durations-day-domain-truncation` (already checked out in that worktree).
**Hand-back file (write when done):** `/home/user/skill-do-work/do-work/runs/work-2026-08-18-211613/REQ-260-handback.md` — absolute main-tree path, the ONE main-tree path you may write. Never stage or commit it.

## The REQ (verbatim)

```markdown
---
id: REQ-260
title: Run the Go formatter as part of the canonical verify
status: claimed
claimed_at: 2026-08-18T21:16:24Z
route: A
created_at: 2026-08-18T18:41:26Z
status_changed_at: 2026-08-18T20:59:31Z
user_request: UR-051
addendum_to: REQ-251
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- _dev/tests/maintainer-verify.sh
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route A
    - 1-file write set
    - 3 acceptance criteria
    - full-suite verification
---

# Run the Go Formatter as Part of the Canonical Verify

## What

`bash _dev/tests/maintainer-verify.sh` runs `go vet` but never the Go formatter, so a formatting slip lands and survives silently — one did, in the board tool's day-domain truncation expression, and was only caught by a builder reading adjacent code. That specific slip is already fixed (REQ-252 corrected it in passing), so this REQ is the rule rather than the instance: the canonical gate should run the formatter over tracked Go files and fail on a non-empty result, exactly as it already runs ShellCheck over tracked shell files.

## Requirements

- The gate runs the Go formatter over tracked Go files and fails when any is unformatted, selecting files the way the ShellCheck lane does (`git ls-files`, never a hand-maintained path list — Closed Enumerations Go Stale).
- The failure names the offending files.
- `bash _dev/tests/maintainer-verify.sh` exits 0 on the current tree — the package is formatted clean today, so a red result means the check found something real.

## Context

Discovered by REQ-251's builder ([low]); the one-character instance was fixed in passing by REQ-252. The user widened the scope at clarify: the instance is closed, the gap that let it through is not.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-251: a gofmt formatting miss in `durations.go` from REQ-248. Should I process this as a new task? → Yes, and widened: add a formatter check to the canonical verify so this class cannot recur silently
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — or fold it into REQ-252, which already owns the file.

**Answered [2026-08-18]:** User approved via `do-work clarify` and **widened the scope**: the original one-character fix is already satisfied by REQ-252's in-passing correction, so this REQ now covers making the gate run the formatter over tracked Go files. Title and requirements updated to match; `effort_estimate` raised from trivial to normal.

---

## Triage

**Route: A** - Simple

**Reasoning:** One named file, and the change mirrors a pattern that already exists inside it (the ShellCheck lane's `git ls-files` selection and fail-on-non-empty shape). Scope is obvious and bounded.

**Planning:** Not required
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
