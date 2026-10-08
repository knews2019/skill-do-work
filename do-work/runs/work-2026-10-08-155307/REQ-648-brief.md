# Builder brief — REQ-648 ([impact-rule-change] The assigned_to schema line and the work action say operator-gated work is blocked, never earmarked)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-648-schema-and-work-action-state-operator-rule
- Branch: worktree-agent-REQ-648-schema-and-work-action-state-operator-rule, created from main HEAD dd54712a with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-648-schema-and-work-action-state-operator-rule.md. Read it fully: What, Why, Verified Facts, Detailed Requirements 1-7, Constraints, Builder Guidance, Red-Green Proof. Requirement 7 (release) belongs to the integrator; you PROPOSE the changelog entry and the lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-142/input.md.
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-08-155307/REQ-648-handback.md
- Route A, tdd: false, impact-rule-change, effort-mechanical, domain general. Prose only.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md — same directory. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (read before editing any action file) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md. Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md read every bullet carrying `[family: alternate-writer-contract-drift]`.

## The change (decided at capture, do not reopen)
Files: `skills/do-work/actions/work-reference.md` and `skills/do-work/actions/work.md` in your worktree.
1. `work-reference.md` line ~114, the `assigned_to:` schema comment: append (where it reads best, before or after the board-placement sentence) "**For another session or checkout only** — work that waits on the user as operator is `status: blocked` with `blocked_by` naming them (`actions/capture.md` Step 1's Earmark assessment), never an invented earmark." Keep every existing sentence of that line, including the parser lock-step sentence (no parser change: the field's read semantics and board placement do not change).
2. `work.md` Error Handling table (line ~519-528), beside the "Builder reports a missing external precondition" row, add: `| Orchestrator wants to park a queued REQ on the operator (production access, credentials, a person at the keyboard) | Flip it to `status: blocked` with `blocked_by` naming the operator and `blocked_at`, the same mid-run flip as a missing precondition. Never write `assigned_to` for this: the earmark is for another session, it is invisible as "waiting on you", and the default run hides it from the user by skipping it. |`
3. Pointer fix: the existing row says "Apply the blocked-flip test (Step 8 → **Mid-run blocked flip**)" and no such heading exists in work.md. Point both rows (existing and new) at the procedure that exists: `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row (the blocked-flip test). Leave the prose mentions at `work-reference.md:165` and `:788` ("mid-run blocked flip", "mid-run blocked-flip procedure") as they are; they describe, they do not point.
4. Do not touch `actions/capture.md`, `actions/capture-reference.md`, `docs/work-guide.md`, or `model.go`.
- Keep every existing bold label, heading, and table shape unchanged; `_dev/tests/contract-regressions.sh` pins several. The new sentences say when NOT to write `assigned_to`; they add no case for writing it.

## Write boundary
Exactly `skills/do-work/actions/work-reference.md` and `skills/do-work/actions/work.md`. Anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `grep -n "operator" skills/do-work/actions/work-reference.md skills/do-work/actions/work.md` finds the schema sentence and the table row.
- `grep -n "Mid-run blocked flip" skills/do-work/actions/work.md` prints nothing.
- `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`; both exit 0.
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit preferred). Base: dd54712a.
- File manifest: each file with (modified) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, the test scripts' exit codes, and each file checked.
- Red-green evidence in prose terms: the RED read (both greps empty before) and the GREEN read (both greps hit; the dangling pointer gone).
- `## Decisions` (D-01 onwards). Expected: where the sentence sits in the schema line; the pointer wording.
- `## Discovered Tasks` (each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for `_dev/primes/lessons-action-files.md` in the file's bullet shape (`- [family: <slug>] [REQ-648: <one-line lesson>](<relative archive link, the integrator fixes the path>)`) and a proposed CHANGELOG entry (descriptive title, plain words, why it matters, specific bullets). The integrator writes both.
- Integration seams (none expected). Test wall times.
