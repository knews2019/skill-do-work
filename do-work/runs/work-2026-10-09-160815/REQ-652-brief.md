# Builder brief — REQ-652 (Routine operator work is never a REQ; a tracked operator configuration is a blocked, dependency-gated REQ that clarify completes)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-652-operator-work-routine-not-captured-tracked-config-blocked-dependency-gated
- Branch: worktree-agent-REQ-652-operator-work-routine-not-captured-tracked-config-blocked-dependency-gated, created from main HEAD BASE_HASH with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-652-operator-work-routine-not-captured-tracked-config-blocked-dependency-gated.md. Read it fully: What, Why (the maintainer's three quoted statements are the spec), Verified Facts, Detailed Requirements 1-7, Constraints, Builder Guidance, Red-Green Proof. Requirement 7 (release) belongs to the integrator; you PROPOSE the changelog entry and the lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-143/input.md (the full upstream proposal, the maintainer's answers and refinements, the F11 verdict).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-09-160815/REQ-652-handback.md
- Route A, tdd: false, impact-rule-change, effort-mechanical, domain general. Prose only.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md — same directory. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (read before editing any action file: template, earned sections, cross-reference spelling, conditions not lists) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (read; the release step is the integrator's). Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md read every bullet carrying `[family: alternate-writer-contract-drift]` (REQ-644 and REQ-648 are this rule's own history).

## The change (decided by the maintainer, do not reopen; keep it SHORT, the maintainer said the orchestrator skill "should not concern itself too much" with operator work)
Files in your worktree: `skills/do-work/actions/capture.md`, `skills/do-work/actions/clarify.md`, `skills/do-work/actions/work.md`, `skills/do-work/actions/work-reference.md`, `skills/do-work/docs/work-guide.md`.
1. `capture.md` line ~108 (External-condition assessment): two sentences keyed on the condition. (a) A routine operator act after the code ships (deploy, publish, approve, verify on live hosts are examples, not a list) is never captured as a REQ: the user does it on their own time, only the AI-buildable parts become REQs, and the capture summary names the operator step left to the user. (b) When the request itself asks to track a special operator configuration, capture that one REQ as `status: blocked`, `blocked_by: '<the operator, in the user's words>'`, `blocked_at`, and `depends_on` listing every AI-buildable REQ minted from the same request, with the checklist in the body; it waits under Pending → Waiting until those complete, then surfaces under Needs input · Blocked, and `do-work clarify` completes it when the operator says it is done.
2. `capture.md` line ~107 (Earmark assessment): the "at the keyboard" sentence points at the new rule instead of flatly capturing `blocked`. Delete words before adding; if the existing sentence can simply gain a clause, do that.
3. `clarify.md` Step 5.5, plain-blocked prompt (~179-190): add `4. Done — I did it myself` to the example prompt and one bullet after "Yes → unblock": on Done, clarify appends a `## Operator receipts` section holding the user's words verbatim (Step 4's Outside-text containment), sets `status: completed` and `status_changed_at: <now>` in place, and does not run `unblock`; the board then shows it under Done and `do-work cleanup` archives it. Cite the precedent in one clause: the confirmed `builder_decided: true` path already flips to `completed` here, and `actions/cleanup.md` Pass 0 sweeps terminal REQs. Step 6 (report) lists each REQ completed this way. No new CLI verb. Add one Verification Checklist line only if the checklist takes it without ceremony.
4. `work.md` line ~523 (the "park on the operator" row) and `work-reference.md` line ~788 (Failure Classification, Environment row): one sentence each: when a run finds that the remaining work is a routine operator act, the REQ completes with what the AI built and the hand-back names the operator step; it is not flipped to `blocked`, not failed, and abandon is not recommended. The blocked flip stays for a precondition to AI work and for a tracked special configuration.
5. `docs/work-guide.md`, the Earmarking or blocked paragraph: one sentence on the shape and where it shows on the board (Pending → Waiting, then Needs input · Blocked, then Done).
- Keep every existing bold label and heading unchanged; `_dev/tests/contract-regressions.sh` pins several. Cross-references use the prime's spelling (`actions/<file>.md` → **Section**). Conditions, not lists: mark example verbs as examples. No Go, no field, no status, no flag, no transaction, no probe, no board change.
- Record in `## Decisions` the challenge the REQ asks for: the proposal's "no REQ, report-only line, abandon" shape was rejected by the maintainer in favour of the blocked + `depends_on` placement the board already has and a clarify completion path.

## Write boundary
Exactly the five files above. Anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `bash _dev/tests/shipped-package-reference-contract.sh` (≈1 s) and `bash _dev/tests/contract-regressions.sh` (≈20-60 s); both exit 0.
- `grep -n "own time" skills/do-work/actions/capture.md`; `grep -n "Done — I did it myself" skills/do-work/actions/clarify.md`; `grep -n operator skills/do-work/actions/work.md skills/do-work/actions/work-reference.md` show the new sentences.
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit preferred). Base: BASE_HASH.
- File manifest: each file with (modified) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, the two test scripts' exit codes, and each file checked.
- Red-green evidence in prose terms: the RED read (the three places before) and the GREEN read (the new sentences, the clarify option), plus the grep output.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk), including the recorded challenge.
- `## Discovered Tasks` (report only unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for `_dev/primes/lessons-action-files.md` in the file's bullet shape (`- [family: <slug>] [REQ-652: <one-line lesson>](<relative archive link, the integrator fixes the path>)`) and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words, one or two sentences on why it matters, then specific bullets). The integrator writes both.
- Integration seams (none expected). Test wall times.
