---
id: REQ-460
title: '[impact-rule-change] Make outside-text delimiter containment condition-complete'
status: pending
created_at: 2026-09-01T00:12:38Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
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

---

