# Builder brief — REQ-644 ([impact-rule-change] Capture distinguishes a session earmark from work that needs the user as operator)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-644-capture-earmark-vs-operator-blocked
- Branch: worktree-agent-REQ-644-capture-earmark-vs-operator-blocked, created from main HEAD with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-644-capture-earmark-vs-operator-blocked.md. Read it fully: What, Why, Verified Facts, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-140/input.md (the maintainer's own words for the distinction are in "Two changes", item 2).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-07-232637/REQ-644-handback.md
- Route A, tdd: false, impact-rule-change, effort-mechanical. Prose-only, two shipped files.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md (same directory; domain: general, so no extra domain crew). Prime: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (read before editing). Lessons: the `[family: alternate-writer-contract-drift]` bullets in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md. Rule: a changed rule is not done until every writer and reader that restates it is swept; grep the restated phrase, not only the cited line.

## The change (decided at capture, do not reopen)
From the REQ's `## Detailed Requirements`, one sentence per site, wording latitude only:
1. `skills/do-work/actions/capture.md`, the **Earmark assessment** bullet (today line 107): add one sentence. Content to carry: `assigned_to` routes work between sessions; when the user must be present as operator, capture `status: blocked` with `blocked_by` naming the person (the user's words, Frontmatter Quoting) and `blocked_at`, which lands in Needs Input · Blocked and is released by `do-work clarify`. Point at the External-condition assessment for the mechanics rather than restating them.
2. `capture.md`, the **External-condition assessment** bullet (today line 108): the sentence "Keep this distinct from the three look-alikes: …" lists depends_on, pending-answers, and Answerer. Add the earmark as a fourth: work reserved for another session is `assigned_to` (not blocked). Because the count "three" is written into that sentence, change it to a condition-keyed phrasing or the correct count; do not leave a stale number (CLAUDE.md: closed enumerations go stale).
3. `skills/do-work/docs/work-guide.md`, the paragraph **Earmarking with `assigned_to`** (today lines 113-119, which opens with "To say 'leave this one for me'"): one sentence saying the field is for another session or checkout, and that work waiting on you as operator is captured `blocked` with `blocked_by` naming you, so it shows under Needs Input · Blocked. Adjust the opening phrase if it now contradicts the sentence ("leave this one for me" reads as operator-present).
4. Sweep every other restatement of the earmark rule in `skills/` (grep: earmark, assigned_to, "leave this", "for the laptop"). Today's grep found only the three sites above plus the schema line in `actions/work-reference.md:114` and the Composed Exit Summary row 6 (`work-reference.md:511`); both describe the field's mechanics, not when to use it, and are outside your boundary. Account for every hit in the hand-back; change none outside the two files.

Constraints: prose only; no new field, status, flag, Go change, or SKILL.md row. Keep the never-invent rule intact: no earmark in the user's words means no `assigned_to`; no person named means no `blocked_by`. Keep every existing heading and bold label unchanged (`_dev/tests/shipped-package-reference-contract.sh` pins citation strings). The Board REQ-643 (a sibling in this wave) will move earmarked REQs from Pending → Ready to a Pending → Earmarked group; do not describe the board's Pending placement of an earmark in your sentences (say only where the blocked shape lands: Needs Input · Blocked), so your text is true before and after that sibling ships.

## Write boundary
Only skills/do-work/actions/capture.md and skills/do-work/docs/work-guide.md (the REQ's `write_set`; Route A has no `## Scope`). Never touch VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/version.md, work-reference.md, clarify.md, crew-members/, the lessons satellite, do-work/lessons-index.md, any Go source or test, SKILL.md, or anything under do-work/ (except the hand-back file). If the sweep shows another file must change, stop and say so in the hand-back; do not write it.

## Verify before hand-back
- From the worktree root, with wall times: `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`. Both must pass.
- Red-Green from the REQ: read your edited Earmark assessment as a cold capture agent receiving "I need to be at the keyboard for this one, hold it for me" and confirm the text now routes it to `blocked` + `blocked_by` + `blocked_at`, and that "leave this for cloud-alpha" still routes to `assigned_to`.
- The restatement grep above, each command pasted, every hit accounted for.
- `git diff --stat` and `git diff --check` on the branch.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash.
- File manifest: each file with (modified) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY] (the orchestrator ticks the REQ's boxes from it); [UNIFY] lists git diff --stat and each file checked.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds; do not fix them; each line ends `→ report only` unless impact-critical).
- Lessons read.
- A proposed lesson bullet for _dev/primes/lessons-action-files.md with a `[family: <slug>]` marker, or "none" with a reason, and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words, one or two sentences then bullets). The orchestrator writes both.
- Integration seams (none expected).
- The two test wall times and the restatement-grep accounting.
