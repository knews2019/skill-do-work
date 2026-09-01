---
id: REQ-269
title: Draw the cross-package citation class by what a citation is, not by the punctuation around it
status: completed
created_at: 2026-08-18T21:44:57Z
status_changed_at: 2026-08-18T22:20:09Z
claimed_at: 2026-08-20T23:21:54Z
completed_at: 2026-08-20T23:36:44Z
kb_status: promoted
kb_entry: REQ-269-draw-the-cross-package-citation-class-by.md
commit: f71dfee
user_request: UR-055
addendum_to: REQ-259
domain: general
route: C
review_generated: true
sweep: true
sweep_key: citation-class-drawn-at-the-marker
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-08-20T23:22:22Z
  basis:
    - Route C
    - 6-file write set
    - 4 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
depends_on: []
maintenance: true
write_set:
- _dev/tests/shipped-package-reference-contract.sh
- _dev/primes/prime-action-files.md
- skills/do-work/crew-members/prompt-injection.md
- skills/do-work-toolbox/crew-members/prompt-injection.md
- skills/do-work-knowledge/crew-members/prompt-injection.md
- skills/do-work/actions/work-reference.md
- skills/do-work-toolbox/actions/code-review.md
- skills/do-work-toolbox/actions/validate-feedback.md
---

# Draw the Cross-Package Citation Class by What a Citation Is, Not by the Punctuation Around It

## What

Eight consecutive REQs in this area have bounded the cross-package citation rule by a **marker** — a leading `../`, a backtick, a fence character — instead of by the **thing** the rule governs: a path in shipped text that a reader is expected to follow from the citing file's own directory. Each fix closes one spelling and leaves the class open, so the next review finds the same defect wearing different punctuation. This REQ ends that by making the checker's condition *be* the rule's condition.

Concretely, three markers are currently doing the bounding, and each has live escapees:

1. **The backtick.** `_dev/tests/shipped-package-reference-contract.sh` checks backticked spans only. Two of REQ-259's three sites existed solely because they were bare text; REQ-259 fixed them by backticking, which moved those two sites into coverage without closing the hole.
2. **The leading `../`.** REQ-259's own sweep found the same skills-folder base spelled with no `../` at all and reasoned it out of the class on the grounds that "it never claimed to be a relative path" — a spelling test standing in for a semantic one.
3. **The fence.** REQ-249 exempted fenced blocks with a stated rationale: their text lands in *some other file*, so it is content rather than a citation from here. The exemption is keyed on the fence character instead of on that rationale, so it also shields in-file annotations that are citations from the file in every sense that matters.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-259's independent review, Important finding 1 (gate: user-visible), consolidating the builder's own D-T1, D-T2 and D-T3 — one root cause, so one sweep REQ rather than three follow-ups. Created `pending-answers` per the generation-≥2 cascade depth stop, since REQ-259 is itself `review_generated: true`.

The review verified the hole by mutation rather than by argument: restoring the exact pre-REQ-259 state (bare text at the wrong depth) makes `shipped-package-reference-contract.sh` **PASS**. That is the class, still open, demonstrated.

## Instances

- [x] **`skills/do-work/crew-members/prompt-injection.md:3`** — `through do-work-toolbox/actions/completed-work-presentation-reference.md`, bare, no `../`; resolves from nowhere but `skills/`. — **fixed**: now `` `../../do-work-toolbox/actions/completed-work-presentation-reference.md` ``. Re-verified: the checker reported this exact line before the fix and passes after.
- [x] **`skills/do-work-toolbox/crew-members/prompt-injection.md:3`** — same text; wrong twice over, since from this file the correct citation is the same-package `../actions/completed-work-presentation-reference.md`. — **fixed**: corrected to the same-package `` `../actions/completed-work-presentation-reference.md` ``, as the REQ predicted. No contract requires the three crew copies to be byte-identical (checked: nothing in `_dev/tests/` compares them), so each now carries the citation correct from where it sits.
- [x] **`skills/do-work-knowledge/crew-members/prompt-injection.md:3`** — same text. — **fixed**: now `` `../../do-work-toolbox/actions/completed-work-presentation-reference.md` ``.
- [x] **`skills/do-work/actions/work-reference.md` lines 130, 132, 137, 204** — four `../do-work-board/...` citations at the retired one-`../` depth (correct is `../../` from `actions/`), shielded by the fence around the Schema Read Contract yaml block whose `#` annotations are documentation for the agent reading that very file. — **fixed, and the real set was larger than the four named**: the widened checker reported lines 105, 130 (×4), 132, 135, 137, 138 and 205. Line 105 and one on 205 were bare rather than backticked. Eight backticked occurrences were re-depthed to `../../do-work-board/…` and the bare one on 205 was backticked and re-depthed; line 105 turned out to be `<skill-root>/../do-work-board/…`, correctly *not* a citation, and is untouched.
- [x] **The fence exemption itself** — split it by its stated rationale (fenced text that lands in another file) rather than by the fence character, in `_dev/primes/prime-action-files.md` § Cross-Referencing and in the checker. — **split by its rationale**, in both homes: `_dev/primes/prime-action-files.md` § Cross-Referencing now states that a fenced payload is exempt because its text lands in another file while an annotation beside it is documentation for this file's reader, and `citation_candidate_tokens` implements exactly that.
- [x] **The checker's span condition** — widen `_dev/tests/shipped-package-reference-contract.sh` past backticked spans so a bare-text citation is checked like any other. — **widened**: the scan surface is now every path-shaped token in prose, in inline code, in HTML comment interiors, and in in-fence annotations. `run_citation_surface_fixtures` pins it, and both halves were mutation-tested to confirm they bite.

## Requirements

- The checker's condition is the rule's condition: a cross-package citation in shipped text is checked whether or not it is backticked, and whether or not it leads with `../`.
- The fence exemption is decided by whether the text lands in another file, not by the presence of a fence.
- **Done means the class cannot recur, not that six more spots were patched.** The Red-Green Proof is the mutation the review ran: restoring a bare, wrong-depth citation anywhere in shipped text must make the checker FAIL, naming the file and line. Today it passes.
- Every instance above is either fixed or reported with a reason.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED:** in the current tree, revert `skills/do-work/actions/commit.md:17` to the bare, wrong-depth form `../do-work-toolbox/actions/inspect.md` (no backticks) and run `bash _dev/tests/shipped-package-reference-contract.sh`. Observed today: `shipped package reference contract: PASS`, exit 0 — the blind spot.

**GREEN:** the same mutation FAILs, naming `skills/do-work/actions/commit.md:17`, and the unmutated tree still passes.

## Open Questions

- [x] REQ-259's review found that the cross-package citation class is bounded by punctuation rather than by what a citation is, with six live escapees across three packages, and demonstrated by mutation that the pre-REQ-259 defect still passes the checker. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — accept that bare-text citations are checked by review rather than mechanically.

**Answered [2026-08-18]:** User approved via `do-work clarify`, presented with the reviewer's mutation evidence — restoring the pre-REQ-259 defect still passes the checker today. Approved as a class fix, not a six-site patch: the Requirements' "done means the class cannot recur" clause is the acceptance bar, and the Red-Green Proof is that mutation.

---

## Triage

**Route: C** - Complex

**Reasoning:** The acceptance bar is explicitly "the class cannot recur, not that six more spots were patched" — that means changing the condition inside `shipped-package-reference-contract.sh`, which is a 900-line checker other REQs depend on, and splitting a documented exemption by its rationale across a prime file and the checker. Multiple packages, a contract test's core predicate, and a mutation-based proof to satisfy.

**Planning:** Required

## Plan

**The class, stated once:** a cross-package citation is a token in shipped text whose first path segment — after any `../` lead — names a sibling package directory, and which a reader is expected to resolve from the citing file's own directory. Nothing about backticks, leading `../`, or fences appears in that sentence, and after this REQ nothing about them appears in the checker's condition either.

### The three markers and what replaces each

1. **The backtick.** `backticked_citation_messages` is called only over `backticked_span_texts(...)`. Replace the scan surface with *every whitespace-delimited token in the file*, backticked or not. The backtick stops being a gate and becomes what it always was — punctuation the tokenizer strips.
2. **The leading `../`.** `backticked_citation_lead = re.compile(r"^(?:\.\./)+")` currently *requires* the lead. Make it optional. The condition becomes "first post-lead segment names a package", which is the semantic test.
3. **The fence.** Split by REQ-249's stated rationale — fenced text that lands in *another* file is content; a `#` or `<!-- -->` annotation inside a fenced block is documentation addressed to the reader of *this* file, so its paths are citations from here. The code already has this instinct: `backticked_span_texts` deliberately reaches inside HTML comments because "JIT_CONTEXT headers carry real cross-package citations". This generalizes that one exception into the rule.

### The ambiguity that survives, narrowed rather than widened

Dropping the `../` requirement collides with the consumer queue root: `do-work/queue/`, `do-work/testers.md` and friends are the consuming project's own state, not citations. Measured across the tracked shipped markdown, a naive widening produces 635 non-resolving hits, of which ~600 are exactly that.

The existing code already documents this as "the core package's directory name equals the consumer queue root". The narrowing that resolves it: **only `do-work` collides.** `do-work-board`, `do-work-knowledge` and `do-work-toolbox` have no consumer-state meaning at all, so a bare token leading with one of those three is unambiguously a citation. A bare token leading with `do-work` is consumer state. A `../`-led token leading with `do-work` is a citation and is checked exactly as today.

That single rule takes the widened result from 635 hits to a two-dozen list that is entirely real.

### Ordered tasks

1. Rename the marker out of the identifiers: `backticked_citation_lead` → `citation_lead`, `backticked_citation_messages` → `citation_messages`, `run_backticked_citation_fixtures` → `run_citation_fixtures`, and the failure message `backticked citation does not resolve in …` → `cross-package citation does not resolve in …`. A name that says "backticked" is the marker-bounded thinking in the most durable place it lives.
2. Make the lead optional and add the bare-`do-work`-is-consumer-state skip, with its reason written at the condition.
3. Add a token scan over prose and over in-fence annotation lines, keeping the existing backticked-span scan as one of its inputs so no current coverage is lost.
4. Split the fence exemption in the checker and restate it by rationale in `_dev/primes/prime-action-files.md` § Cross-Referencing.
5. Extend the fixture table with the cases that pin the new condition: bare non-core citation must fail; bare `do-work/queue/` must pass; wrong-depth in-fence annotation must fail; fenced payload must still pass.
6. Fix every site the widened checker reports; tick each `## Instances` entry with its evidence.
7. Run the Red-Green Proof mutation and confirm it now FAILs.

### Plan validation

- **Requirement coverage:** every Requirement maps to a task — condition-is-the-rule → 2+3, fence-by-rationale → 4, class-cannot-recur → 5+7, instances → 6, verify green → 7.
- **No orphan tasks:** task 1 is the only one not named by a Requirement; it is the naming half of the same defect and is cheap to drop if reviewed as scope creep.
- **Scope sanity:** 7 tasks, above the 3-task quality threshold. Flagged as the action requires. Not split, because the acceptance bar is explicitly a class fix and a partial one fails it — but this is the REQ most likely in this queue to have wanted splitting at capture.

## Exploration

- `_dev/tests/shipped-package-reference-contract.sh` — `citation` predicate at 1183-1241; called at 1485-1495 over `backticked_span_texts` only; fixtures at 1244-1330. `strip_markdown_code` (391) masks fenced/list-fenced/indented code in phase 1 (214) and backtick spans in phase 2 (322).
- `backticked_span_texts` (368) already reaches into HTML comments deliberately — the precedent for task 3.
- Measured escapee list under the refined predicate (deduped, probe artifacts removed): the three `crew-members/prompt-injection.md:3` bare citations; `skills/do-work-toolbox/actions/code-review.md:328` and `validate-feedback.md:114` (two) shielded by a fence; `skills/do-work/actions/work-reference.md` lines 130, 132, 135, 137, 138, 205 at the retired one-`../` depth inside the Schema Read Contract yaml block's `#` annotations.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — the condition, the scan surface, the fence split, the fixtures
- `_dev/primes/prime-action-files.md` (modify) — § Cross-Referencing restated by rationale
- `skills/do-work/crew-members/prompt-injection.md` (modify) — instance 1
- `skills/do-work-toolbox/crew-members/prompt-injection.md` (modify) — instance 2
- `skills/do-work-knowledge/crew-members/prompt-injection.md` (modify) — instance 3
- `skills/do-work/actions/work-reference.md` (modify) — instance 4
- `skills/do-work-toolbox/actions/code-review.md` (modify) — fence-shielded citation the widened checker reports
- `skills/do-work-toolbox/actions/validate-feedback.md` (modify) — same

**Files I will NOT touch:** `skills/do-work/CHANGELOG.md` beyond the release entry — its historical `do-work/` mentions are consumer-state and dated history.

**Acceptance criteria (restated from the REQ):**
1. The checker's condition is the rule's condition — checked whether or not backticked, whether or not `../`-led.
2. The fence exemption is decided by whether the text lands in another file.
3. The Red-Green mutation FAILs, naming file and line. It passes today.
4. Every `## Instances` entry is fixed or reported with a reason.
5. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Implementation Summary

**What was done:** Redrew the cross-package citation class in the contract checker so its condition is the rule's condition. A token is a citation when its first path segment — after an optional `../` lead — names a sibling package directory. Backticks, the `../` lead, and the fence stopped being part of that test. The scan surface widened from backticked spans only to every path-shaped token in prose, inline code, HTML comment interiors, and in-fence annotations; the fence exemption was split by REQ-249's stated rationale so a fenced payload stays exempt while an annotation beside it does not. Two conditions take a token back out of the class, both about meaning: a path rooted somewhere else (`<skill-root>/…`, `.claude/skills/…`) is not a path from here, and a bare `do-work/…` token is the consumer's queue root — a collision only the core package's name has.

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified) — `citation_lead` made optional; `citation_shape` and `citation_candidate_tokens`/`citation_tokens_in` added as the new scan surface; `mask_block_code` grew an optional `masked_line_ranges` out-parameter so the fence split reuses the one fence state machine; `citation_messages` rewritten with the narrowed consumer-state ambiguity; `run_citation_surface_fixtures` added; three cases added to the citation fixture table; `backticked_*` renamed out of the citation identifiers and the failure message.
- `_dev/primes/prime-action-files.md` (modified) — § Cross-Referencing restates the class by condition and the fence exemption by rationale.
- `skills/do-work/crew-members/prompt-injection.md` (modified) — bare citation → `../../do-work-toolbox/…`.
- `skills/do-work-knowledge/crew-members/prompt-injection.md` (modified) — same.
- `skills/do-work-toolbox/crew-members/prompt-injection.md` (modified) — corrected to the same-package `../actions/…`.
- `skills/do-work/actions/work-reference.md` (modified) — nine citations re-depthed from `../do-work-board/…` to `../../do-work-board/…`, one of them bare and now backticked.

**Declared but not touched** (see D-04): `skills/do-work-knowledge/actions/setup-memory.md` and `skills/do-work-toolbox/actions/install.md`. The first widened pass reported them; anchoring the shape match at the token start showed their `<skill-root>/../do-work/scripts/…` tokens were never citations from those files, so editing them would have introduced the error rather than fixed one.

**Tests touched:** `run_citation_surface_fixtures` (new, 5 assertions) and three new cases in `run_citation_fixtures`. Both new fixture sets were mutation-tested.

## Decisions

- **D-01** (DECIDE & STATE): Renamed `backticked_citation_*` → `citation_*` and the failure message to `cross-package citation does not resolve…`. Reasoning: the marker-bounded thinking the REQ describes was literally encoded in the identifiers, and a future reader grepping `backticked` would have found the citation logic and concluded the backtick was part of the contract.
- **D-02** (ESCALATE): The widened predicate needed a rule for bare tokens, and a naive widening produced 635 non-resolving hits against the shipped corpus — unusable. Builder chose: skip a **bare** `do-work/…` token as consumer state, and treat a bare token under any other package name as an unambiguous citation. Reasoning: the existing code already documents the core package's name colliding with the consumer queue root; this narrows that documented ambiguity to the one package it actually affects rather than inventing a new exemption. Measured effect: 635 hits → 17, all real. Value: the bare-citation half of the class closes without a false-positive flood that would have got the check disabled. Risk: if a consumer ever names its queue root after a non-core package, a real consumer-state path would be reported as a broken citation. Reversible in one condition; no such consumer exists, and the queue root is `do-work/` by contract.
- **D-03** (DECIDE & STATE): The fence split is implemented as "annotations inside a masked block are scanned, payload is not", with the annotation found by `#` or `<!--`. Reasoning: it is the mechanical form of REQ-249's own stated reason. The code already had the instinct — `backticked_span_texts` reaches into HTML comments because "JIT_CONTEXT headers carry real cross-package citations" — so this generalizes an existing one-off rather than adding a new concept.
- **D-04** (ESCALATE): The first widened pass reported `setup-memory.md:42` and `install.md:39,42,183` as broken. They are not: the tokens are `<skill-root>/../do-work/scripts/…`, rooted at a placeholder, and a mid-token regex search had matched the tail. Builder chose to anchor the shape match at the token start rather than add an exemption list for placeholder roots. Reasoning: "a citation from here starts the token" is the same semantic test the rest of this REQ applies, and an exemption list would have been the fourth marker. Value: `.claude/skills/…`, `<suite-root>/…` and every future rooted spelling are handled by the condition with nothing to maintain. Risk: a citation preceded by an unstripped punctuation character would be missed; the opening-punctuation strip set covers the markdown cases and `run_citation_surface_fixtures` pins two of them. **The two files stayed in the write set but needed no edit** — recorded here rather than silently dropped.

- **D-05** (DECIDE & STATE): `skills/do-work-toolbox/actions/code-review.md:328` and `validate-feedback.md:114` were declared in `## Scope` from the exploration probe and then correctly **not** touched. Both sit inside a ```` ```markdown ```` REQ template whose text is copied into a new REQ file, so the fence exemption's surviving half applies exactly as designed and the widened checker never reports them. Builder chose to leave them and record the reason rather than edit them to satisfy a declaration. Reasoning: editing them would have introduced an error to make a scope list tidy. This is the source of the `scope-drift.sh` DRIFT reported in Qualification, and it is intended.

## Qualification

Passed — 6 files verified, 5 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `git diff --stat` reviewed file by file; the only executable change is the checker, which is exercised by its own fixtures plus the whole-repo run. Grepped the diff for `console.log`, `debugger`, `print(`, `TODO`, `FIXME` — the only `print(` matches are the checker's pre-existing reporting, which is contract output. `shellcheck` clean (the file is in `maintainer-verify`'s 73-file lint set) and the run exits 0.
- **Substantive:** the checker diff changes the predicate and the scan surface, not formatting.
- **Requirements traced:** AC1 → `citation_messages` + `citation_candidate_tokens`; AC2 → the fence split in both homes; AC3 → mutation re-run below; AC4 → all six `## Instances` ticked with evidence; AC5 → verify exits 0.
- **Flowing:** the checker actually probes the filesystem — proven by the mutation runs, which change its verdict.

## Testing

- `bash _dev/tests/shipped-package-reference-contract.sh` — PASS, exit 0, on the unmutated tree.
- `bash _dev/tests/maintainer-verify.sh` — exit 0. Baseline was green before implementation (`do-work/working/baseline.json`, `launched: true`, `exit_status: 0`), so this is a clean no-regression comparison rather than a comparison against a red record.
- `shellcheck` clean: the checker is inside `maintainer-verify`'s 73-file warning-level lint set.

**Red-green validation** — traced to the REQ's `## Red-Green Proof`, which named the mutation rather than a test:

| Mutation of `skills/do-work/actions/commit.md:17` | Before this REQ | After |
|---|---|---|
| bare, wrong depth (`../do-work-toolbox/actions/inspect.md`) — **the captured RED** | PASS (the blind spot) | **FAIL**, naming `skills/do-work/actions/commit.md:17` |
| backticked, wrong depth | FAIL | FAIL |
| bare, no `../` lead at all (`do-work-toolbox/actions/inspect.md`) | PASS | **FAIL** |
| unmutated | PASS | PASS |

The captured RED was reproduced first, on the untouched tree, before any implementation — confirming the blind spot rather than assuming it. The GREEN row is the same mutation after.

**Fixture mutation testing** — a fixture that passes on its first run proves nothing, so both new fixture sets were checked for bite by reverting the code they pin:

| Reverted behavior | `run_citation_surface_fixtures` |
|---|---|
| anchoring (`match` → `search`) | FAILs: leaked `../do-work/scripts/tool.sh` |
| fence split (annotation scan dropped) | FAILs: missed `../do-work-board/actions/board.md` |

**Existing tests updated:** none. The three added cases in `run_citation_fixtures` and the new `run_citation_surface_fixtures` are additions; no existing assertion changed meaning.

## Review

**Overall: 89%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Checker's condition is the rule's condition — backtick and `../` no longer gate | ✅ |
| Fence exemption decided by whether the text lands in another file | ✅ Payload exempt, annotations checked; stated in prime and implemented in `citation_candidate_tokens` |
| Class cannot recur — the captured mutation must FAIL | ✅ Reproduced RED first, then GREEN; two further spellings also caught |
| Every instance fixed or reported with a reason | ✅ All six ticked with evidence; two extra sites found beyond the four named |
| `maintainer-verify.sh` exits 0 | ✅ |

### Findings

**Important — none.**

**Minor:**

- **M1:** `scope-drift.sh` reports two declared-but-untouched files. Investigated and intended — D-05. Recording it rather than silencing it, because the correct resolution of a drift report is an explanation, not an edited declaration.
- **M2:** The REQ named four wrong-depth lines in `work-reference.md`; the widened checker found nine sites there, two of them bare rather than backticked. The capture-time list undercounted, which is itself the REQ's thesis — a list written by eye is bounded by what the eye was scanning for. No action; noted because it is evidence the class fix was the right shape.

**Nit:**

- **N1:** `citation_candidate_tokens` yields duplicates when the same path appears twice on one line, so `work-reference.md:130` produced four FAIL lines for three distinct paths. Harmless for a checker whose output is a fix list, and de-duplicating would hide a genuinely repeated citation. Not queued.

### Restatement Sweep

Redefined element: what counts as a cross-package citation, and the fence exemption's boundary.

- `_dev/primes/prime-action-files.md` § Cross-Referencing — the canonical prose home. **Updated in this diff**; it was the one stale restatement and it is the file the REQ's write set named first.
- `_dev/tests/contract-regressions.sh` — enforces the separate CLAUDE.md/AGENTS.md citation ban; unrelated condition, unaffected. Verified by the green run.
- `skills/do-work/actions/work-reference.md` Schema Read Contract block — a consumer of the citation form, not a statement of the rule. Its citations were corrected as instance 4.
- `_dev/primes/prime-shell-commands.md` and `prime-kanban-board.md` — grepped for citation-form statements; neither states this rule.
- The identifiers and the failure message inside the checker — a restatement in the most durable form there is. Renamed (D-01).

No stale restatement remains.

### Acceptance Testing

The checker is the deliverable, so acceptance is behavioral and is the mutation table in `## Testing` above: the exact captured RED now fails, two further spellings of the same class fail, and the clean tree passes. Ran against the whole live corpus (every tracked `.md` under `skills/`), not a fixture repo.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 85% |
| Test Adequacy | 95% |
| Scope Discipline | 85% |
| Risk | Low |
| Acceptance | Pass |

Code Quality 85%: `citation_candidate_tokens` is doing three jobs (prose, regions, in-fence annotations) in one generator. It is readable and commented at each branch, but it is the function most likely to want splitting the next time this area is touched. Scope Discipline 85% for M1's declared-but-untouched pair, intended and explained.

**Risk: Low, not None.** This widens a gate that runs on every commit. The false-positive mode is a rooted or consumer-state path misread as a citation; measured against the full corpus it is currently zero, and the two conditions that produce it are pinned by fixtures. The failure mode is loud and immediate, never silent.

### Follow-up REQs Created

One — see `## Discovered Tasks`.

*Reviewed in orchestrated mode by the work orchestrator*

## Discovered Tasks

- **impact-rule-change** A fenced template's payload is exempt from citation checking because it lands in another file, and nothing verifies the payload resolves at that destination. `skills/do-work-toolbox/actions/code-review.md:328` and `validate-feedback.md:114` both carry a citation that is wrong where it lands. Fold-first scan: no `pending`/`pending-answers` REQ shares this root cause, and it is not a prose restatement so the standing sweep does not take it. Queued as **REQ-310** (`pending-answers`), with the count-the-population-first instruction that may close it as a two-line fix.


## Lessons Learned

**What worked:** Reproducing the captured RED on the untouched tree *before* writing any code. The REQ handed over a mutation and a claim that it passes today; running it first turned that claim into an observation and made every later "it fails now" mean something. The second thing that worked was measuring before designing — a naive widening of the predicate produces 635 non-resolving hits, and knowing that number is what forced the narrowed consumer-state rule instead of a plausible-looking widening that would have been reverted within a day.

**What didn't:** The exploration probe was written twice before it was right, and both failures were the same mistake in miniature as the REQ itself. First `str.strip(".")` ate the leading `..` of every `../` token, so instance 4 looked clean when it was not. Then a mid-token regex `search` matched the tail of `<skill-root>/../do-work/scripts/…` and reported four rooted paths as broken citations. Both are punctuation standing in for meaning — the exact defect being fixed, committed by the tool built to find it. The second one survived into the implementation and was only caught because the real checker reported the four sites; the fix for it (anchor the match at the token start) turned out to be the correct semantic rule and is now a fixture.

**Worth knowing:** The three `crew-members/prompt-injection.md` copies are byte-identical across the packages but nothing enforces that, and they must now differ — from `do-work-toolbox` the correct citation is the same-package `../actions/…`, not a cross-package hop. Anyone tempted to re-mirror those three files will silently reintroduce instance 2. Related: `mask_block_code` now takes an optional `masked_line_ranges` out-parameter so the fence split reuses one fence state machine rather than growing a second walker — if a third caller needs the ranges, that is the seam to extend, not to copy.

## Orientation

The cross-package citation check now enforces the class it always meant to: a path is a citation when it names a sibling package from this file's own directory, whether or not it is backticked and whether or not it is spelled with `../`. A fenced block's payload stays exempt for its real reason while the annotations beside it are checked. Lives in the shipped-reference contract (`_dev/tests/shipped-package-reference-contract.sh`) with its prose home in `_dev/primes/prime-action-files.md` § Cross-Referencing.

**[MAP CHANGED]** — this renames a concept and widens a contract. `backticked_citation_*` is gone from the checker; the citation identifiers and the failure message no longer name a marker, so anyone grepping `backticked` for the citation logic will not find it. The check's reach also grew from backticked spans to every path-shaped token in prose, inline code, comment interiors, and in-fence annotations, which means shipped text that was previously invisible to it is now covered.

Prime staleness spot-check: `_dev/primes/prime-action-files.md` — updated by this REQ, and its other referenced paths still resolve (the green contract run is the evidence).
