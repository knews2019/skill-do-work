# Builder brief — REQ-647 (Clarify tells the user when a released REQ stays earmarked and clears it on request)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-647-clarify-reports-released-req-still-earmarked
- Branch: worktree-agent-REQ-647-clarify-reports-released-req-still-earmarked, created from main HEAD dd54712a with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-647-clarify-reports-released-req-still-earmarked.md. Read it fully: What, Why, Verified Facts, Detailed Requirements 1-7, Constraints, Builder Guidance, Red-Green Proof. Requirement 7 (release) belongs to the integrator; you PROPOSE the changelog entry and the lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-142/input.md (holds the consumer incident timeline and the full triage).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-08-155307/REQ-647-handback.md
- Route A, tdd: false, impact-user-visible, effort-mechanical, domain general. Prose only.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md — same directory. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (read before editing any action file: template, earned sections, cross-reference spelling) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (read; the release step is the integrator's). Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md read every bullet carrying `[family: alternate-writer-contract-drift]` (REQ-644's entry is this rule's own history).

## The change (decided at capture, do not reopen)
Files: `skills/do-work/actions/clarify.md` and `skills/do-work/docs/work-guide.md` in your worktree.
1. `clarify.md` Step 5.5, the "Yes → unblock" bullet (line ~188): keep the transaction sentence and the "do not reproduce the mutation free-form" rule. After the transaction, add the condition: when the released REQ carries `assigned_to`, tell the user on the spot that the REQ is still earmarked for `<session name, verbatim>` and that the default run will skip it, then ask one question with two concrete options (clear-questions.md is already loaded by Step 3): clear the earmark now (remove the `assigned_to` line by hand edit, the documented clear path: `actions/work-reference.md` → Request File Schema `assigned_to`, and the exit summary's assigned-elsewhere row), or keep it and run the REQ by name with `do-work run REQ-NNN`. Replace the closing sentence with: "The REQ re-enters the queue for the next `do-work run` unless it carries `assigned_to`, in which case the default run skips it until it is named explicitly or the field is cleared."
2. Say how that hand edit is committed: exactly the way clarify's other hand edits in the same invocation are committed today. Read clarify.md to find that (Step 4's `answer --manifest … --commit` is the question path; look for how reclaim lines and the Builder Was Right / Discarded fast paths reach a commit). Do not invent a new commit step or transaction. If clarify has no stated commit discipline for hand edits, say so in the hand-back and reference `actions/commit.md` in one clause.
3. `clarify.md` Step 6 (Report): one sentence keyed on the condition: the summary lists every released REQ that stays earmarked, with the verbatim session name and the run-by-name command.
4. `work-guide.md`, the "Earmarking with `assigned_to`" paragraph (line ~113-120): one sentence saying that when clarify releases a blocked REQ that is also earmarked, it says so and offers to clear the earmark. Keep every existing sentence.
5. `clarify.md` Verification Checklist: one line that every unblocked REQ still carrying `assigned_to` was named to the user and in the report, only if the checklist's existing shape takes it without ceremony; otherwise skip and say so.
- Never clear the field silently on unblock; never clear it on "keep". `assigned_to` stays verbatim (no normalization). No Go, no flag, no transaction, no `verify` probe, no edit to `work-reference.md` or `capture.md`.
- Cross-references use the prime's spelling (`actions/work-reference.md` → **section name**). Keep every existing bold label and heading unchanged; `_dev/tests/contract-regressions.sh` pins several of them.

## Write boundary
Exactly `skills/do-work/actions/clarify.md` and `skills/do-work/docs/work-guide.md`. Anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `bash _dev/tests/shipped-package-reference-contract.sh` (≈1 s) and `bash _dev/tests/contract-regressions.sh` (≈20-60 s); both exit 0.
- `grep -n assigned_to skills/do-work/actions/clarify.md` shows the Step 5.5 condition and the Step 6 line.
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit preferred). Base: dd54712a.
- File manifest: each file with (modified) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, the two test scripts' exit codes, and each file checked.
- Red-green evidence in prose terms: quote the RED read (clarify.md before: no `assigned_to` mention) and the GREEN read (the new condition, the new closing sentence, the Step 6 line), plus the grep output.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk). Expected: the commit-discipline choice for requirement 2, the checklist choice for requirement 5.
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for `_dev/primes/lessons-action-files.md` in the file's bullet shape (`- [family: <slug>] [REQ-647: <one-line lesson>](<relative archive link, the integrator fixes the path>)`) and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words, one or two sentences on why it matters, then specific bullets). The integrator writes both.
- Integration seams (none expected). Test wall times.
