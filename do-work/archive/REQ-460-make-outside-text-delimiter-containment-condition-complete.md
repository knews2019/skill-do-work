---
id: REQ-460
title: '[impact-rule-change] Make outside-text delimiter containment condition-complete'
status: completed
created_at: 2026-09-01T00:12:38Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-420]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-413]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-413
sweep: true
sweep_key: markdown-delimiter-containment-prefix-gaps
status_changed_at: 2026-09-03T11:10:03Z
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-03T02:13:57Z
  basis:
    - Route B
    - 2-file write set
    - 1 subsystem involved
    - 4 acceptance criteria
    - cross-route regression gates
claimed_at: 2026-09-03T11:10:55Z
route: B
dispatch_at: 2026-09-03T11:17:09Z
builder_handback_at: 2026-09-03T11:19:48Z
integration_at: 2026-09-03T11:19:48Z
review_at: 2026-09-03T11:24:06Z
kb_status: pending
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
completed_at: 2026-09-03T11:31:24Z
commit: 7e16f05c4e95ebf50fcf2d065e4f0145246d46ad
release_at: 2026-09-03T11:31:24Z
---

# Make Outside-Text Delimiter Containment Condition-Complete

> **Deferred (maintainer, 2026-09-01):** substantive-effort edge-case hardening; hold
> until REQ-420's parity work ships rather than spending a pipeline slot on it while
> the UR-081 spine is unfinished. The `depends_on: [REQ-420]` gate above is this
> decision made mechanical — remove it to un-defer.

## What

Replace the finite prefix list used for answer-summary containment with a condition-complete classifier for lines that Markdown can interpret as document-owned delimiters or structure.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-action-files.md` and its satellite, `prime-shell-commands.md` § Closed Enumerations Go Stale, `prime-kanban-board.md` (for the goldmark question), and CLAUDE.md. Approach: state the condition once — Markdown opens every block construct from one of three line-start ingredients — and test those three instead of enumerating constructs.
- [x] **[APPLY]:** Two files, both inside the declared Scope. Nothing else touched.
- [x] **[UNIFY]:** Orchestrator independently re-ran `go build ./...`, `go vet ./...` (clean) and `gofmt -l .` (silent), read the full `answer.go` diff, and confirmed `git status --porcelain` lists only the two Scope files. No debug artifacts; the added comments state the condition and the doubt-resolution direction rather than narrating steps.

## Finding Provenance

REQ-413's fresh re-review found that `summaryRequiresContainment` recognizes a fixed set of ten prefixes but accepts other delimiter-shaped lines such as `***` and `___`. The action contract states a condition—whether the line could be read as the file's own delimiter—not a closed list of examples.

## Instances

- Markdown thematic breaks formed from asterisks, underscores, or hyphens.
- Heading, block, fence, list, quote, HTML-like, indentation-sensitive, and other structural shapes supported by the request format.
- Future syntax representatives that satisfy the same document-boundary condition without appearing in an example list.

## Detailed Requirements

- Define the containment boundary from the semantic Markdown/document condition rather than a hand-maintained prefix enumeration.
- Require file-backed raw payload and canonical containment whenever a one-line answer summary can be interpreted as document structure or delimiter.
- Preserve lossless bytes for safe plain summaries and for contained outside text.
- Use table-driven or generated adversarial cases spanning every supported structural class, including `***` and `___`.

## Constraints

- Do not weaken C0/DEL or multiline containment.
- Examples are fixtures, not the definition of the accepted set.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm. Express one structural predicate and make the tests demonstrate classes rather than mirror another finite production list.

## Red-Green Proof

**RED prompt/case:** Submit `***`, `___`, and representatives from every Markdown structural class as inline answer summaries without a raw payload.
**Why RED now:** The current ten-prefix test accepts delimiter shapes outside its enumeration.
**GREEN when:** Every structural/delimiter-shaped representative requires canonical containment while ordinary one-line prose remains inline and lossless.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-413-rereview.md` for source context and independent evidence.

---
*Source: REQ-413 fresh re-review finding F5.*

## Triage

**Route: B** — Medium

**Reasoning:** The invariant is firm and localized to answer-summary containment plus its direct fixtures, but the implementation must replace a finite Markdown-prefix enumeration with one semantic predicate and prove every structural class without weakening plain-summary byte preservation.

**Planning:** Not required for Route B; exploration will re-establish the exact parser and test seams after recovery.

**Estimate:** 20 active minutes (P50, medium confidence; frozen from the recovered estimate).

## Plan

**Planning not required** — Route B: exploration-guided implementation.

*Skipped by work action*

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3968 tokens; relevant to shipped action-contract restatements, but the partial-slug satellite has no narrower family matching Markdown containment within the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 3006 tokens; relevant to typed publication containment, but the partial-slug satellite has no narrower family matching this delimiter-classification bug within the 2000-token budget.

## Exploration

`summaryRequiresContainment` and its caller-seam tests are already present in current `HEAD` through implementation commit `7e16f05c`, which is an ancestor of this claim. The current tree is byte-equivalent to that commit for the two implementation paths, and neither path belongs to the preserved REQ-461 dirty diff. The unfinished state is therefore lifecycle/release provenance, not missing code.

The implementation replaces ten literal prefixes with three line-start ingredients: block-significant whitespace, ASCII punctuation, or a digit run followed by `.` or `)`. Its table covers headings, setext underlines, thematic breaks including `***` and `___`, bare and ordered lists, fences, quotes, HTML/link forms, tables, indentation, and dialect fences; caller-seam tests prove raw-payload refusal, lossless containment, and safe prose inlining. C0/DEL and multiline handling remain in the unchanged shared containment helpers.

Historical review findings already landed as REQ-528 and REQ-529, so this recovery must not mint duplicates. The prior completion/release commit `a4a887ed` is not an ancestor of this branch; a current-branch lifecycle/release transaction is still required and will record `7e16f05c` as supplied implementation provenance.

*Generated by fresh Explore agent during run-with-recovery*

## Scope

**Files I will touch (supplied implementation scope; no new edit required):**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` — condition-complete summary classifier from supplied commit 7e16f05c.
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` — structural-class and caller-seam RED/GREEN coverage from supplied commit 7e16f05c.

**Files I will NOT change during recovery:** the two implementation files above, shared C0/DEL and multiline containment helpers, request-state cancellation containment, and the preserved REQ-461 dirty paths. This pass restores lifecycle/release evidence for code already in `HEAD`.

**Acceptance criteria (restated from REQ):**
- [x] Structural classification is expressed as a semantic line-start condition, not a literal prefix list.
- [x] `***`, `___`, setext, bare list markers, tables, indentation, and dialect fences require raw-payload containment.
- [x] Safe letter-led and non-marker digit-led summaries remain inline byte-for-byte; contained summaries retain their exact source bytes.
- [x] C0/DEL rejection and multiline containment remain unchanged.
- [x] Test rows name structural classes so examples remain fixtures rather than the production definition.

## Decisions

### D-05: Reuse the merged implementation commit and rebuild only the missing tail

**Decision:** DECIDE & STATE — treat `7e16f05c` as supplied implementation provenance, verify its exact two-file diff against current `HEAD`, and do not rewrite already-landed code.

**Reasoning:** The earlier implementation branch was merged through `18920627`, but its later release commit `a4a887ed` was not. Reapplying identical source would manufacture a hollow diff and obscure provenance; the current branch needs only reconstructed evidence plus a new lifecycle/release transaction.

## Implementation Summary

**Supplied implementation commit:** `7e16f05c4e95ebf50fcf2d065e4f0145246d46ad`

**Files changed by that commit and adopted here:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go`
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go`

`summaryRequiresContainment` now tests the structural condition directly: block-significant leading whitespace, ASCII punctuation, or a digit run followed by `.` or `)`. The caller requires a byte-matching raw payload for these summaries and emits canonical containment; safe prose remains inline. Table-driven tests cover all requested structural classes and exercise both refusal and successful containment through `BuildAnswerPlan`.

Fresh builder validation found both current blobs exactly equal to the supplied commit (`answer.go` blob `4100e6ee127cffbc7bb0247d54f1fcaae88edd98`; `answer_test.go` blob `74491508f161dc0fbe41940d7fd5dc0dd4115b82`). No source edit or new implementation commit was needed during recovery.

## Qualification

**Passed** — the two-file manifest matches Scope and `write_set`; all six restated acceptance criteria trace to the supplied diff; P-A-U is complete.

- `git merge-base --is-ancestor 7e16f05c HEAD` confirms the supplied implementation is in the current branch.
- Reverse-application and blob checks confirm the exact implementation is present without a hollow rewrite.
- `go test -count=1 ./internal/publication` passed in 21.029s.
- `go vet ./...` passed.
- `git diff --check` passed.
- The shared C0/DEL and multiline containment helpers are unchanged by `7e16f05c`.

## Testing

**RED evidence:** against the parent predicate, `***`, `___`, setext underlines, bare list markers, bare ordered markers, indentation, tables, and dialect fences miss the ten-prefix enumeration for the documented reasons.

**GREEN evidence:** the current publication package passes the structural-class matrix and `BuildAnswerPlan` caller-seam tests. Structural summaries refuse a missing raw payload and preserve matching raw bytes inside canonical containment, while ordinary letter-led and non-marker digit-led prose remains inline.

*Validated by a fresh builder during run-with-recovery; no source edits made.*

## Review

**Approve — Overall: 89%** | Acceptance: Pass | Risk: Low

| Dimension | Score |
|---|---:|
| Requirements | 88% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 100% |

The supplied implementation closes the prefix-enumeration defect with the declared two-file scope. Fresh review reran the full publication package, vet, diff hygiene, and the named structural/caller-seam tests at exit 0. The result is a strict containment widening: required structural forms are carried losslessly while ordinary prose remains inline, and the shared C0/DEL and multiline helpers are untouched.

No new Important findings were found. The prior adjacent findings remain routed without duplication: disposition-marker spoofing is REQ-528, cancellation-reason parity is REQ-529, and the inaccurate prose restatements are already in `do-work/prose-backlog.md`. Minor observations remain report-only: the ordered-list comment overstates Markdown universality, conservative punctuation/number over-catch is not pinned, and whitespace-only input is refused before this predicate at the caller seam.

*Reviewed independently against `7e16f05c^..7e16f05c` by the review-work action.*

## Lessons Learned

**What worked:**
- State the structural condition first and make examples test fixtures. Leading whitespace, ASCII punctuation, and an ordered-list digit marker cover the requested class without another brittle prefix catalog.
- Compare old and new predicates as a superset, then exercise the real `BuildAnswerPlan` seam; this proves both no weakening and byte-preserving containment.

**What didn't:**
- The old implementation trimmed leading whitespace before classification and thereby erased the evidence for indented structure.
- Reasoning only from a few passing delimiters made a partially correct enumeration look complete; review also found that the original hazard explanation ignored the summary's actual mid-line write position.

**Worth knowing:** punctuation-led prose and number-like text are intentionally over-contained as the safer reversible tradeoff. Adjacent disposition and cancellation consumers are owned by REQ-528 and REQ-529.

## Orientation

Answer publication now decides containment from a structural line-start condition instead of ten remembered prefixes. The behavior lives in the do-work CLI publication subsystem; no module or data-flow map changed, and both referenced primes remain structurally current.

---
