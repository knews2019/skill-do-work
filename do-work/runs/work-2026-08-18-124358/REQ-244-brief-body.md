---
id: REQ-244
title: Cite the Timestamp rule at every timestamp write site
status: claimed
created_at: 2026-08-18T12:28:33Z
user_request: UR-055
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-245]
batch: timestamp-stamping-integrity
estimate:
  p50_active_minutes: 55
  confidence: low
  calculated_at: 2026-08-18T13:05:12Z
  basis:
    - Route C
    - 12-file write set
    - 1 new files
    - 4 subsystems involved
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
claimed_at: 2026-08-18T13:05:12Z
route: C
---

# Cite the Timestamp Rule at Every Timestamp Write Site

## What

Sweep all four skills for every timestamp write site — templates and action steps carrying `[timestamp]`, `<timestamp>`, `<now>`, `[UTC timestamp]`, or any `*_at:`/date-shaped placeholder — normalize each to the spellings the Timestamp rule recognizes (`<timestamp>` / `<now>`), and add an inline citation of the rule (`Timestamp rule, actions/work-reference.md`) at each site that lacks one.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

An agent filling a template from context never re-reads the rule when nothing at the site points to it; a fabricated `created_at` on two review-generated REQs was reported as the resulting incident. The rule's own design already mandates "Every other site cites the rule and stops" — uncited bare placeholders are drift from that architecture.

## Detailed Requirements

Sites confirmed uncited at capture (starting set — the sweep is the requirement, this list is not the extent):

- `skills/do-work/actions/review-work.md:365` — "Review Fix" follow-up template `created_at: [timestamp]` (the site that produced the reported incident)
- `skills/do-work/actions/review-work.md:425` — report footer `**Overall: [X]%** | [timestamp]`
- `skills/do-work/actions/work-reference.md:627` — Builder-Decided Follow-up Template `created_at: [timestamp]`
- `skills/do-work/actions/work-reference.md:893` — Session Checkpoint Template `session_ended: [timestamp]`
- `skills/do-work-toolbox/actions/code-review.md:301` — follow-up template `created_at: [timestamp]`
- `skills/do-work/actions/forensics.md:216,257` and `skills/do-work/actions/roadmap.md:135,244` — `**Scan date:** [timestamp]`
- `skills/do-work-toolbox/actions/present-work.md:86` — `**Generated:** [UTC timestamp]`
- `skills/do-work-toolbox/actions/deep-explore.md:250` — `completed_at: <timestamp>` (recognized spelling, no citation in a skill that never loads the rule)
- `skills/do-work-knowledge/actions/interview.md` / `interview-reference.md` — the `<now>` cluster (`started_at`, `last_activity_at`, `approved_at`, `last_validated_at`, `review_completed_at`, `last_exported_at`); recognized spelling, cross-skill, uncited

Grep-verified at capture: `grep -c "Timestamp rule"` returns 0 for review-work.md, code-review.md, roadmap.md, present-work.md, and interview.md.

## Constraints

- **Citations only, never command copies.** `skills/do-work/actions/work-reference.md` ~line 101 states the Timestamp rule's paragraph "is the only place in `actions/` that spells a command for obtaining one" and documents why per-site copies failed (Windows agents). The sweep must not recreate that.
- Cross-skill citations from do-work-toolbox / do-work-knowledge follow the existing precedent in `skills/do-work-knowledge/actions/memory.md`.
- Distinguish instants from date-only stamps: the rule's own "Date-only stamps" paragraph governs `YYYY-MM-DD` sites (log filenames, headings) — do not convert those to instant placeholders. Path slugs like `work-<timestamp>` in run-directory names are names, not stamps, and are out of scope.
- Finding provenance (validate-feedback triage, this session): verdict Accept; Surface-cost N/A — aligning sites to an existing documented rule, no new defensive surface.

## Red-Green Proof

**RED prompt/case:** A new lock-in check in `_dev/tests/` (wired into `maintainer-verify.sh`) greps shipped `skills/*/actions/` for bare timestamp placeholders (`[timestamp]`, `[UTC timestamp]`) and for stamp write sites in files that never cite the Timestamp rule — it fails on the current tree, naming the sites listed above.
**Why RED now:** Those sites exist today with no citation; an agent filling them has nothing pointing at the rule or a clock command.
**GREEN when:** Every stamp write site uses a recognized spelling with an inline Timestamp-rule citation, the lock-in check passes, and `bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** Inferred during capture

## Builder Guidance

Certainty: Firm on the sweep and citations; the exact lock-in check pattern (how "uncited site" is detected mechanically) is the builder's call — keep it condition-keyed, not a hand-maintained site list, per CLAUDE.md's Closed Enumerations rule.

## Full Context

See `do-work/user-requests/UR-055/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding 1 — "AUDIT: sweep all four skills for every timestamp write site … and bring each under the Timestamp rule with an inline citation, normalizing placeholder spelling to the forms the rule recognizes. The list above is a starting set, not the full extent."*

---

# Builder Guardrails (orchestrator-issued — binding)

## Your tree

- Work **only** inside your worktree (path below). It is a full checkout on your own branch.
- **Never write anywhere in the main tree** except the single hand-back file named below. That is the only main-tree path you may touch.
- **Never touch `do-work/`** — not the queue, not `working/`, not `CHECKPOINT.md`, not `archive/`. Queue state is the orchestrator's alone. Your branch must contain **zero** commits touching `do-work/`; the orchestrator runs `git diff --name-only <pre>...<your-branch> -- do-work/` and a single path there stops your hand-back.
- **Never touch `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, or `CHANGELOG.md`.** Those are serial-only integrator-owned files. A bump on your branch races every sibling.
- **Scratch files go in `/tmp` or inside your worktree — never the main tree root.** A previous builder left a PNG in the repo root; that is a write-set violation. Screenshots, fixtures, generated boards: `/tmp`.

## Commit on your branch

Commit your implementation on your own branch before handing back. Message body only — no version bump, no changelog entry. End the message with:

```
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
```

## The P-A-U loop is yours to fill

The REQ body contains an `## AI Execution State (P-A-U Loop)` section with three checkboxes, or the orchestrator will add one. **You must tick all three and write the required content into each**, in your worktree's copy of nothing — instead, put the filled P-A-U block **into your hand-back file** verbatim under a `## P-A-U` heading, since you may not write `do-work/`. `qualify.sh` FAILs on unticked boxes and the orchestrator will otherwise have to fill them from your evidence.

- **[PLAN]** — read the listed `prime_files` and agent rules, then write the technical approach. No code yet.
- **[APPLY]** — code written exactly as planned, scope strictly limited to planned files.
- **[UNIFY]** — run `git diff --stat`, review every changed file, run the project's linters/tests, confirm no debug artifacts. List each file you verified and what you checked.

## Evidence rules — every one of these was learned by getting it wrong

1. **Two REDs when the first is a reference error.** A test that fails because a constant or function does not exist yet proves nothing. Put the code in place, break exactly one rule, and let the assertion fail *for the reason it exists*. Report both RED outputs.
2. **`git stash push` on a clean file stashes nothing** — and the resulting green run reads as proof when it is vacuous. To reproduce RED against pre-change code, check out the pre-change blob by hash (`git show <hash>:<path>`) instead.
3. **Assert page identity inside the same call that reads the DOM.** If you drive a browser, return `location.href` (and, where relevant, the page's own rule text) from the *same* `evaluate` call as every measurement. A shared browser instance can be navigated by a sibling between your navigate and your evaluate, and the numbers come back confident, well-formed, and about somebody else's page. A URL checked *before* navigating is not the same claim. Prefer an isolated browser context.
4. **A programmatic `.focus()` does not trigger `:focus-visible` in Chrome.** Use a real `Tab` keypress if focus styling is in question.
5. **Generate the artifact and look at it.** For anything that changes what appears on screen, a passing assertion is not evidence about two glyphs sharing a coordinate. Measure `getBoundingClientRect()` intersections in the live DOM when the question is "do two things overlap"; read the rendered text when the question is "what does this say".
6. **Push back if the brief is wrong.** If a requirement contradicts an existing test, or a piece of code you wrote turns out unneeded, say so in the hand-back rather than quietly editing the test or keeping dead code. Two builders pushed back last session and both were right.

## Verification bar

`bash _dev/tests/maintainer-verify.sh` from your worktree root. **Exit code 0 is the only proof.** Never pipe it through `tail`/`head` — the pipeline's exit status hides the failure. Run it, then `echo $?` on its own line, and paste that.

## Hand-back

Write **one** file, at the absolute path given below, containing:

1. `## Branch` — your branch name.
2. `## P-A-U` — the three filled, ticked checkboxes with their content.
3. `## Files Changed` — `git diff --stat` against your branch's merge base, plus one line per file saying what changed and why.
4. `## Red-Green Evidence` — the RED output(s) and the GREEN output, quoted.
5. `## Verification` — the `maintainer-verify.sh` tail and its `echo $?` line.
6. `## Integration Seams` — anything the orchestrator must apply by hand in the merge commit (shared registries, cross-REQ text). Say "none" if none.
7. `## Decisions` — numbered D-01, D-02… for choices with reach beyond this REQ.
8. `## Lessons Learned` — what a future session should know. Omit if genuinely nothing.
9. `## Pushback` — anything in this brief you think is wrong. Omit if none.

Your final message back should be a short summary; the hand-back file is the real deliverable.

## Your Assignment

- **Worktree path (your working directory):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-244-cite-the-timestamp-rule-at-every-stamp-write-site`
- **Branch name:** `worktree-agent-REQ-244-cite-the-timestamp-rule-at-every-stamp-write-site`
- **Hand-back file (absolute, main tree — the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-18-124358/REQ-244-handback.md`
- **Repo root of the MAIN tree (read-only for you):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`

## Orchestrator Notes for This REQ

**You have no `write_set`, and that is deliberate — it is not licence for free scope.**

The REQ's requirement *is* the sweep, and its listed sites are explicitly "a starting set, not the extent". So:

1. **Run the sweep FIRST, before writing anything.**
2. **Then declare your own scope** in a `## Scope` section in your hand-back — "Files I will touch", each with (new|modify), and "Files I will NOT touch". That declaration is what the orchestrator mirrors into `write_set`, and it is what your diff will be judged against.
3. Only then write code.

A file that appears in your diff but not in your declared Scope is scope drift and will be treated as such.

**Sibling-builder collision guard.** REQ-242 is building concurrently and owns `skills/do-work-board/tools/queue-kanban/web/board-durations.js` and `generate_test.go`. REQ-245 has just landed changes to `skills/do-work-board/tools/queue-kanban/{model.go,verify.go,web/board-cards.js,web/board-core.js,web/board.css,prime-do-kanban.md}`. Stay out of all of those. Your corpus is action markdown under `skills/*/actions/`.

**REQ-243's link checker is already live in your worktree.** It resolves every relative `.md` link in shipped markdown from the citing file's own directory AND checks that each `#anchor` names a real heading in the target. You are about to add a lot of citations of the form `Timestamp rule, actions/work-reference.md`. If `maintainer-verify.sh` fails on a link or anchor after your edits, **that is the checker working, not a regression** — fix your link. Take advantage of it: it means a mistyped citation cannot ship.

**Constraints worth repeating because they are easy to violate at scale:**

- **Citations only, never command copies.** `skills/do-work/actions/work-reference.md` ~line 101 says the Timestamp rule's paragraph "is the only place in `actions/` that spells a command for obtaining one", and records why per-site copies failed (Windows agents). Do not recreate that. A site gets a *pointer*, never a `date -u ...` line.
- **Distinguish instants from date-only stamps.** The rule's own "Date-only stamps" paragraph governs `YYYY-MM-DD` sites (log filenames, headings) — do not convert those to instant placeholders. Path slugs like `work-<timestamp>` in run-directory names are **names, not stamps**, and are out of scope. Getting this wrong is the most likely way to make this REQ worse than the problem.
- **Your lock-in check must be condition-keyed, not a hand-maintained site list** (CLAUDE.md → Closed Enumerations Go Stale; the detail is in `_dev/primes/prime-shell-commands.md`). The exact mechanical definition of "uncited site" is your call — that is stated in the REQ's Builder Guidance — but a list of filenames will rot and will be rejected.
- **`maintenance: false` here, but CLAUDE.md's "Delete before you add" still applies.** Before adding a check, look at whether an existing `_dev/tests/` suite is its natural home. REQ-243 just proved the value of that sweep: it found half of its own REQ already implemented.

**Context for why this REQ exists.** Last session three `completed_at` stamps were written by extrapolating the clock forward instead of reading it. The board's own future-stamp check caught them (corrected in `818ea17`). The diagnosis is that an agent filling a template from context never re-reads the rule when nothing at the site points to it. You are building the fix for that. Do not commit the same error while doing so — read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything, and never carry a stamp forward from earlier in your session.
