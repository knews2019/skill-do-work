---
id: REQ-460
title: '[impact-rule-change] Make outside-text delimiter containment condition-complete'
status: claimed
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
status_changed_at: 2026-09-01T18:54:11Z
claimed_at: 2026-09-03T02:13:10Z
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

## Triage

**Route: B** - Medium

**Reasoning:** The defect and the single function are named exactly, but "every Markdown structural class the request format supports" is not enumerated anywhere in the REQ — the set of shapes the classifier must cover had to be established from the format itself before any predicate could be written.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`summaryRequiresContainment` (`internal/publication/answer.go:22`) is ten lines: trim leading spaces and tabs, test ten literal prefixes — `#`, `>`, `` ` ``, `~`, `- `, `* `, `+ `, `[`, `<`, `---` — then fall back to `orderedListDelimiterPattern` (`^[0-9]+[.)][ \t]`). Its one caller is `BuildAnswerPlan` at `:117`, which uses the verdict to decide whether a one-line answer summary may be written inline or must be carried as a file-backed raw payload with canonical containment.

The gaps are structural, not a missing entry or two:

- **Thematic breaks.** CommonMark accepts three or more `*`, `_`, or `-` with any internal spaces. `---` is listed; `***`, `___`, `* * *`, `- - -` and every longer run are not. The REQ names `***` and `___` explicitly.
- **Bare list markers.** The `- `, `* `, `+ ` prefixes all require a trailing space, so a line that is exactly `-`, `*`, or `+` — a valid empty list item — passes as prose.
- **Ordered lists.** The pattern requires a trailing space or tab, so `1.` alone passes.
- **Setext heading underlines.** A run of `=` is a heading underline in context and appears nowhere in the list.
- **Indentation.** `TrimLeft` discards leading whitespace before testing, which is what an indented code block is made of — four leading spaces is itself structure, and trimming it away is how the current shape loses that class entirely.

This is the pattern `_dev/primes/prime-shell-commands.md` § **Closed Enumerations Go Stale** names, and CLAUDE.md restates as "state conditions, not lists": the action contract asks whether a line *could be read as the file's own delimiter or structure*, and the code answers with a membership test against ten examples. The fix is to express that condition once and demote the examples to fixtures.

The safe direction is unambiguous here: a false positive costs a file-backed payload for a line that did not need one and loses nothing, while a false negative writes an unescaped delimiter into the document it is being embedded in. The classifier should resolve doubt toward containment.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modify) — replace the prefix enumeration with one structural predicate expressed as a condition
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modify) — table-driven adversarial cases spanning every structural class, with the class named per row

**Files I will NOT touch:** the C0/DEL and multiline containment checks (the REQ forbids weakening either), and every other file in `internal/publication`.

**Acceptance criteria (restated from REQ):**
- [ ] The containment boundary is defined by the semantic Markdown/document condition, not a hand-maintained prefix enumeration
- [ ] A one-line answer summary that can be read as document structure or a delimiter requires a file-backed raw payload and canonical containment
- [ ] Bytes stay lossless for safe plain summaries and for contained outside text
- [ ] Adversarial cases are table-driven and span every supported structural class, including `***` and `___`
- [ ] C0/DEL and multiline containment are not weakened
- [ ] Examples in tests are fixtures, not the definition of the accepted set
