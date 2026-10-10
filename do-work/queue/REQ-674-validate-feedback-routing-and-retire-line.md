---
id: REQ-674
title: 'Core do-work routes validate-feedback to the toolbox, and the suite installer narrates a retire line for the old standalone do-validate-feedback skill'
status: pending
created_at: 2026-10-09T21:26:55Z
user_request: UR-149
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-670, REQ-671, REQ-672, REQ-673]
batch: validate-feedback-capture
required_lessons: ["skills/do-work/tools/lessons-do-work-update.md"]
claimed_at: 2026-10-10T12:52:22Z
status_changed_at: 2026-10-10T12:57:19Z
---
# Core do-work Routes validate-feedback, and the Installer Narrates a Retire Line for do-validate-feedback
## What
Two parts. C5a: add a `validate-feedback, triage feedback` row to the core `do-work` routing table that forwards to `../do-work-toolbox/actions/validate-feedback.md`, so "do-work validate-feedback: …" no longer falls through to capture. C5b: the suite installer gains a read-only check that finds an older standalone `do-validate-feedback` skill folder under `.claude/skills/` or `.agents/skills/` and narrates a line naming it and the `git rm -r <path>` command to remove it. The installer never deletes it.
## Why
Report item C5 (UR-149 input): 21 prompts invoked `do-work validate-feedback`, and core `do-work` has no route for it, so the phrase lands in the capture fallback. Separately, 18 identical copies of a 62-line standalone `do-validate-feedback` skill sit in `.claude/skills/` (and one `.agents/skills/`) across 9 consumer repos. It predates the toolbox action and lacks its prompt-injection guard, Surface-cost check, and provenance-preserving capture handoff, and the harness advertises it under the phrase "validate feedback".
## Verified Facts (checked at 0.305.87)
- `skills/do-work/SKILL.md:24-46`: the core routing table has no `validate-feedback` row; `:46` routes "`capture-request:` / `capture request:` or unmatched descriptive multi-word input" to capture. The table header says "first match wins".
- `skills/do-work-toolbox/SKILL.md:18`: `| validate-feedback, triage feedback, review feedback | ./actions/validate-feedback.md |`, the only forward today.
- No core routing row forwards into a sibling package today; this row would be the first.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go:743-791` (`reviewAndConfirm`) narrates module diffs and calls `warnAboutDirtyManagedPaths` (`:786`, defined at `:809`) before the `[y/N]` prompt. Nothing looks at sibling skill folders the suite does not own.
- `skills/do-work/actions/version.md:60`: the update "does not mutate … any other project configuration". So retiring an old skill must be an offer, not part of the install transaction.
- Tests for the installer live in `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` and `update_transaction_test.go`.
## Detailed Requirements
C5a, core routing:
1. Add a row to `skills/do-work/SKILL.md`'s routing table above the capture fallback at `:46`: `| validate-feedback, triage feedback | ../do-work-toolbox/actions/validate-feedback.md |`.
2. "do-work validate-feedback: …" routes to the toolbox action, not to capture.
C5b, installer:
3. Add a read-only check beside `warnAboutDirtyManagedPaths`: if `.claude/skills/` or `.agents/skills/` contains a folder named `do-validate-feedback`, narrate "`<path>` is an older standalone copy that answers the same phrase as `do-work-toolbox validate-feedback`; remove it with `git rm -r <path>`", with the real path filled in.
4. Keep a short list of known predecessor names, starting with `do-validate-feedback`.
5. Never delete or modify the folder inside the transaction, per `version.md:60`. An install or update into a repo with `.claude/skills/do-validate-feedback/` prints the retire line and leaves the folder in place.
6. Release per `_dev/primes/prime-releases.md`.
## Constraints
- The installer check is read-only and never changes the transaction's outcome or its `[y/N]` prompt.
- No new REQ fields, no new statuses.
- Related queued work, no edge: REQ-668 (the `do-work status` action) and REQ-669 (the `do-work trace` action) each add a row to the same routing table; whichever lands later rebases its row.
- Out of scope: any change to verdict rules, capture templates or run selection.
## Assumptions (recorded at capture, no questions asked)
- C5 is one numbered item with two parts (C5a routing, C5b installer), so it is one REQ as the capture instruction asked. The builder may ship the two parts as separate commits within the REQ.
- The predecessor-name list is a deliberate closed list, an exception to the maintainer's "state conditions, not lists" rule: the condition "a skill that answers the same phrase" would need parsing other skills' descriptions, which no current case needs. One name today; a second known predecessor is the time to revisit.
- "Offer to retire it" (Request C5) means the narrated line with the `git rm -r` command, as the Proposed direction says; no separate prompt is added.
- The row is placed so first-match-wins resolves "validate-feedback" and "triage feedback" to it; "review feedback" is left out of the core row (the report's row omits it and the core `review` row would match first).
- The check also runs on the update path, since `version.md:60` says the updater delegates to the installed full-suite installer.
- `tdd: true` covers C5b (a Go test in `suiteinstall` asserting the retire line and the untouched folder); C5a is prose and is proven by the routing acceptance check.
## Dependencies
None. Independent of REQ-670 to REQ-673; the route works with today's read-only action.
## Builder Guidance
Certainty is high; the report names the row text, the check's location, the message and the no-delete rule. Latitude: the helper's name and the exact narration wording.
## Red-Green Proof
**RED prompt/case:** an install into a fixture repo containing `.claude/skills/do-validate-feedback/SKILL.md`; and the input "do-work validate-feedback: <a sample finding>".
**Why RED now:** `install_transaction.go:743-791` inspects only managed paths, so no line names the old folder; `SKILL.md:24-46` has no `validate-feedback` row, so the input falls to the capture fallback at `:46`.
**GREEN when:** the install output contains the retire line naming `.claude/skills/do-validate-feedback` and `git rm -r .claude/skills/do-validate-feedback`, the folder still exists afterwards, and the same holds for `.agents/skills/`. A repo without the folder prints no such line. "do-work validate-feedback: …" routes to `../do-work-toolbox/actions/validate-feedback.md`.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers do-work-cli internals, which the installer check lives in.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the budget left after the selected entry; `slugged: partial`). Matching reason: its index row covers changing action routing.
## Full Context
See `do-work/user-requests/UR-149/input.md` for complete verbatim input (sections Request C5, What happened, Where the behaviour lives today, Proposed direction C5a and C5b, Acceptance check). No queued candidate shares this root cause.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-validate-feedback-capture.md`, Request item C5: "Give `validate-feedback` a forward in the core `do-work` routing table, and have the suite installer detect the older standalone `do-validate-feedback` skill that still answers the same phrase in consumer repos, and offer to retire it."*
