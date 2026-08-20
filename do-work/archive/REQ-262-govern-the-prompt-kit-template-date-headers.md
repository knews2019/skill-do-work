---
id: REQ-262
title: Govern the prompt-kit templates' date headers
status: completed
created_at: 2026-08-18T19:30:47Z
status_changed_at: 2026-08-18T20:59:31Z
claimed_at: 2026-08-20T22:51:58Z
completed_at: 2026-08-20T23:18:39Z
kb_status: pending
commit:
user_request: UR-055
addendum_to: REQ-253
domain: general
route: B
effort_estimate: trivial
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-20T22:52:38Z
  basis:
    - trivial short-circuit
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work-reference.md
---

# Govern the Prompt-Kit Templates' Date Headers

## What

Three prompt-kit templates in do-work-knowledge carry `Date: [today]` headers that no paragraph of the Timestamp rule governs and that sit outside the citation checker's reach (they are template content, not action prose). Decide whether they join the date-only paragraph's consumer list (UTC, cited) or are declared template-content-out-of-scope like the fenced-block exemption.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Discovered by REQ-253's builder's `[today]` class grep ([low]). The exact template paths are re-derived at claim time by the same grep — the count, not the list, is the contract.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-253: three prompt-kit templates carry ungoverned `Date: [today]` headers. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — prompt templates are consumer-facing content and can stay ungoverned.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

<!-- D-XX counter: none used. Next decision: D-01. -->

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is a single rule decision (templates join the date-only paragraph's consumer list, or are declared out of scope), but the exact template paths and the current wording of the Timestamp rule's consumer list both have to be discovered before the decision can be written anywhere.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The three sites (re-derived by grep, per the REQ's Context — the count is the contract, not the list):**

- `skills/do-work-knowledge/prompts/prompt-kit-step0-pen-and-paper-exercises-to-prepare-prompt.md:87`
- `skills/do-work-knowledge/prompts/prompt-kit-step3-spec-engineer.md:67`
- `skills/do-work-knowledge/prompts/prompt-kit-step4-intent-and-delegation-framework.md:62`

Three, matching the REQ. All three sit **inside a fenced block** that is the artifact the prompt tells the model to emit for the user (`=== PRE-FLIGHT BRIEF ===`, `=== PROJECT SPECIFICATION ===`, `=== DELEGATION & INTENT FRAMEWORK ===`). `[today]` is a fill-in placeholder addressed to the model, alongside `[name]`, `[who maintains this]`, and `[one-line descriptor the user gave]` — not a stamp any skill step writes.

**The rule's home:** `skills/do-work/actions/work-reference.md` line 113 — the "Date-only stamps are a different shape" paragraph. It already carries two carve-outs shaped exactly like this decision:

- *"A local-time date is a different thing again and is correct where it is used deliberately (changelog entry headings, run-directory names, report slugs); those sites are not governed here."*
- The `## HH:MM UTC` daily-log headings: *"this rule deliberately does not govern them; their write sites are marked, and a sweep walks past every `## HH:MM UTC` heading instead of converting it."*

The second is the precedent that matters: the fix for an ungoverned date-shaped token is to **state the exclusion at the rule's home so a sweep walks past it**, not to enumerate the sites.

**Citation-checker reach confirmed:** `_dev/tests/shipped-package-reference-contract.sh` masks fenced, list-fenced, and indented code before checking (its `strip_markdown_code` phase 1, line ~214). Template content inside a fenced block is already unreachable by it, so nothing there needs changing.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work-reference.md` (modify) — extend the date-only paragraph's carve-out sentence so it governs template placeholders by condition

**Files I will NOT touch:**
- The three prompt-kit templates. The decision is that their `Date: [today]` lines are correct as written; editing them would be the opposite of the outcome.
- `_dev/tests/shipped-package-reference-contract.sh`. Its fenced-code masking already puts these sites out of reach; no checker change is implied.

**Acceptance criteria (restated from the REQ):**
1. The three `Date: [today]` headers are governed by a decision, one way or the other — either they join the date-only paragraph's consumer list, or they are declared template-content-out-of-scope.
2. The decision lives where the rule lives, so a reader of the Timestamp rule finds it without re-deriving the question.
3. The statement is keyed on the condition, not on a three-path list that goes stale the moment a fourth template is written.

## Decisions

- **D-01** (ESCALATE): The canonical gate `bash _dev/tests/maintainer-verify.sh` was red before this REQ touched anything, so no hand-back in this run could show the exit 0 CLAUDE.md requires as proof. Root cause: REQ-283 (`2308afd`) routed board verification through the board skill but never carried the change into `_dev/tests/staged-skills-contract.sh` or `skills/do-work/actions/help.md`. A separate escaped-line-continuation defect in `_dev/tests/record-commit-hash-guards.sh` was creating a fixture as a 0-byte file. Builder chose: repair the gate before implementing, and land the repair as **its own commit (`8e9cc46`, version 0.222.5) outside this REQ's boundary** rather than folding it into REQ-262. Reasoning: the repair is a precondition of verifying REQ-262, not part of its change, so a shared commit would misattribute four files to a prose decision about date headers. Because it stayed outside the boundary, `## Scope` and `write_set` are unchanged and there is no scope drift to report. Value: every REQ in this run and afterwards can produce the exit-0 proof the project's own rule demands. Risk: the fix updates a lock-in test's expected pass-through string — if REQ-283's routing decision was itself wrong, this cements it. Reversible in one commit; REQ-283's archived `## What` was read first and states the intent explicitly.
- **D-02** (DECIDE & STATE): The two candidate outcomes the REQ named were "join the date-only paragraph's consumer list" or "declared template-content-out-of-scope". Builder chose out-of-scope. Reasoning: all three sites sit inside the fenced block the prompt tells the model to emit *for the user*, alongside `[name]` and `[who maintains this]`. Making them UTC would stamp a machine instant into a human's specification document, which is the wrong shape for the artifact and would still not be a value any pipeline step reads.
- **D-03** (DECIDE & STATE): Written as a condition, not as three paths. Reasoning: `crew-members/maintenance.md` names the stale-enumeration case directly — generalize a hand-maintained list to a trigger condition so there is nothing to keep in sync — and the REQ's own Context says the count is the contract, not the list.

## Implementation Summary

**What was done:** Added one condition-keyed sentence to the date-only-stamp paragraph in `skills/do-work/actions/work-reference.md`, declaring that a date placeholder inside a template artifact's fenced block is out of scope for the Timestamp rule. The sentence keys on what the site *is* (a fill-in token addressed to the model that emits a document for the user) rather than on which files currently contain one, and names neither the paths nor how many there are. The three prompt-kit templates were deliberately not edited: the decision is that their `Date: [today]` lines are correct as written.

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified) — one sentence added to the date-only paragraph placing template-artifact placeholders out of scope, keyed on the condition.

**Tests touched:** none. The change is prose in an action reference; `_dev/tests/shipped-package-reference-contract.sh` already covers this file's citation and mirror contracts and passes.

## Qualification

Passed — 1 file verified, 3 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `git diff --stat` shows one project file changed (`skills/do-work/actions/work-reference.md`, 1 insertion / 1 deletion — the sentence is inserted into an existing long line). Grepped the diff for `console.log`, `debugger`, `TODO`, `FIXME` — none. No linter applies to Markdown prose; `maintainer-verify.sh` exits 0, which includes this file's citation and shipped-reference contracts.
- **Substantive:** the modified file's diff carries the decision text itself, not whitespace.
- **Requirements traced:** AC1 → the sentence states the verdict; AC2 → it sits in the date-only paragraph at the rule's home, three sentences from the `## HH:MM UTC` carve-out it is modeled on; AC3 → no path list, condition + count only.
- **Flowing:** not applicable — no data path; this is instruction prose.

## Review

**Overall: 92%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| The three `Date: [today]` headers are governed by a decision | ✅ Declared out of scope at the rule's home |
| The decision lives where the rule lives | ✅ `work-reference.md` line 113, in the date-only paragraph, adjacent to the `## HH:MM UTC` carve-out it is modeled on |
| Keyed on the condition, not a three-path list | ✅ No paths, and after the sweep finding below, no count either |

### Findings

**Important — none.**

**Minor:**

- **M1 (fixed during review):** The first draft of the sentence carried "`../../do-work-knowledge/prompts/` holds three today". The Restatement Sweep matched it against `_dev/primes/prime-kanban-board.md:51` — REQ-261's recorded lesson that *a conditional keyed on a count that does not bear on the argument is clutter* — and against `CHANGELOG.md` 0.213.x, which removed a count-keyed tripwire from **this same paragraph** with the note "the rule never needed the count anyway". Shipping a fresh count into the paragraph a prior REQ had just de-counted would have re-introduced the exact drift. Removed; the sentence now reads on the condition alone.

**Nit:**

- **N1:** The paragraph is now five sentences of carve-out on one line. Not this REQ's to fix, and each carve-out earns its place, but if a sixth class of date-shaped token turns up, the paragraph wants restructuring into a short list rather than a seventh clause. Recorded as an observation, not queued — a REQ minted for it today would be speculative.

### Restatement Sweep

Redefined element: the scope boundary of the date-only stamp rule — which date-shaped tokens it governs.

Swept every other statement and consumer of that boundary:

- `skills/do-work-toolbox/actions/ui-review.md:210` (report header `**Date**` cites the rule) — still accurate; the new exclusion is a different class and does not touch it.
- `skills/do-work-knowledge/actions/memory.md:50` and `actions/memory-reference.md:55` (`memory/logs/` mirror; `## HH:MM UTC` heading) — still accurate, unchanged by this diff.
- `skills/do-work/CHANGELOG.md:285-288, 358, 421, 2244` — dated history describing past states, not live contract. Not stale.
- `_dev/primes/prime-kanban-board.md:51` — produced finding M1 above, fixed in this diff.
- `_dev/tests/shipped-package-reference-contract.sh` — the citation-count contract still passes, so the new citation did not move a pinned figure.

No stale restatement remains.

### Acceptance Testing

`bash _dev/tests/maintainer-verify.sh` exits 0 with the change in place, which is the acceptance test that applies: this file is covered by the shipped-package reference contract (citation resolution, changelog mirror byte-identity) and by the action-shell-block probes. Non-behavioral prose change, so no red-green pair applies; the exit-0 run is the regression evidence.

Re-read the rendered paragraph end to end as the acceptance check a reader would perform: the new sentence sits between the local-time carve-out and the `## HH:MM UTC` carve-out, reads in the same voice, and answers the REQ's question without sending the reader anywhere else.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | N/A (prose change; contract suite is the proof) |
| Scope Discipline | 95% |
| Risk | None |
| Acceptance | Pass |

Scope Discipline is 95% rather than 100% only because the run repaired a red canonical gate before implementing. That repair was deliberately kept out of this REQ's commit and out of its Scope declaration (D-01), so `scope-drift.sh` reports clean and the deduction is a note, not drift.

### Follow-up REQs Created

None. M1 was fixed in-diff; N1 is speculative and deliberately not queued.

*Reviewed in orchestrated mode by the work orchestrator*

## Discovered Tasks

- **impact-rule-change** The work pipeline's Step 6.5 test resolution is area-scoped (prime test map, else per-file generic detection), so a change can pass testing and review while leaving the repository's own canonical whole-repo gate red. REQ-283 did exactly that and `maintainer-verify.sh` stayed red across three REQs. Fold-first scan found no `pending`/`pending-answers` REQ sharing this root cause and it is not a prose restatement, so the standing sweep does not take it — queued as **REQ-309** (`pending-answers`).
- **fixed in-run, not queued** — `_dev/tests/record-commit-hash-guards.sh` line 428 carried `\\` where a line continuation was meant, so `printf` never ran and a bare redirect created the "healthy neighbour" fixture as a 0-byte file. Repaired in commit `8e9cc46` because the gate could not go green otherwise. Recorded here so the discovery is not invisible, not as an open task.

## Lessons Learned

**What worked:** Reading the target paragraph's *existing* carve-outs before writing a new one. The `## HH:MM UTC` exclusion three sentences away was already the exact shape this decision needed — an out-of-scope class stated at the rule's home so a sweep walks past it — so the work was matching an established pattern rather than inventing a policy. Also: running the Restatement Sweep against the REQ's own diff, not only against other files. It is what caught M1.

**What didn't:** The first draft stated "holds three today". It felt like useful precision and it was the opposite — `_dev/primes/prime-kanban-board.md:51` records REQ-261 deciding that a count which does not bear on the argument is clutter, and the changelog shows a count-keyed tripwire being removed from *this same paragraph* for the same reason. Writing a count into a paragraph a prior REQ had just de-counted would have re-opened settled drift. The sweep caught it; reading the prime first would have prevented it.

**Worth knowing:** The three prompt-kit `Date: [today]` lines were never reachable by the citation checker in the first place — `_dev/tests/shipped-package-reference-contract.sh` masks fenced, list-fenced, and indented code before it reads. So the REQ's premise ("outside the citation checker's reach") was correct but understated: they are outside it *structurally*, by a property of every fenced block, not by an accident this decision needed to preserve. That is why the exclusion could be stated as a condition with nothing to keep in sync.

## Orientation

The Timestamp rule's date-only paragraph now answers a question it previously left open: a date placeholder inside a template artifact is out of scope, decided by what the site is rather than by which file it sits in. Lives in the core action reference (`skills/do-work/actions/work-reference.md`), the schema-and-rules subsystem that `_dev/primes/prime-action-files.md` indexes. Leaf change — no new module, data flow, or contract, and no renamed concept.

Prime staleness spot-check: `_dev/primes/prime-action-files.md` — its referenced paths all still exist; this change makes nothing in it stale.
