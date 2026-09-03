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
- [x] The containment boundary is defined by the semantic Markdown/document condition, not a hand-maintained prefix enumeration
- [x] A one-line answer summary that can be read as document structure or a delimiter requires a file-backed raw payload and canonical containment
- [x] Bytes stay lossless for safe plain summaries and for contained outside text
- [x] Adversarial cases are table-driven and span every supported structural class, including `***` and `___`
- [x] C0/DEL and multiline containment are not weakened
- [x] Examples in tests are fixtures, not the definition of the accepted set

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modified)

**What was done:** `summaryRequiresContainment` no longer tests membership in a ten-prefix list. Its doc comment states the condition once — Markdown builds every block construct at a line start from exactly three ingredients: block-significant leading whitespace, an ASCII punctuation mark, or a digit run forming an ordered-list marker — and the body tests those three. `isMarkdownBlockPunctuation` is four range comparisons covering the whole ASCII punctuation block, which is CommonMark's own definition of ASCII punctuation, so a construct a dialect adds from any of those characters is caught with no edit. `orderedListDelimiterPattern` became `orderedListMarkerPattern` and lost its trailing-space requirement, so bare `1.`, `1)` and `1234.` are caught. The leading-whitespace trim was removed from the front of the predicate: it used to discard the very bytes that make an indented code block, which is how that class was lost entirely.

## Decisions

- **D-01** — Doubt resolves toward containment at the **character** level, not the construct level: any leading ASCII punctuation requires containment, with no carve-out for marks no dialect currently uses. **ESCALATE.** (The builder recorded this as DECIDE & STATE; reclassified here, because it widens behavior for summaries a user would not think of as structural and is genuinely contestable.) **Value:** the predicate has no inert-punctuation list to revisit, so a dialect claiming `:`, `$`, `+` or `|` is caught the day it ships rather than the day someone notices — `+++`, `:::note` and `$$` are all caught today for exactly this reason, and the REQ asks for future syntax representatives by name. **Risk:** a plain summary like `(none)` or `"blue"` now takes a file-backed payload it does not need. No bytes are lost — the text still lands verbatim inside canonical containment and the checkbox reads `See contained answer note` — but it is extra ceremony for an innocuous answer, and a user who expected their two-word reply inline will see it moved. Reversible by narrowing the ranges; irreversible in no sense.
- **D-02** — Any leading space or tab requires containment, without measuring indent depth. **DECIDE & STATE.** Four spaces opens an indented code block and one to three do not, but a summary that begins with whitespace is degenerate either way, and the depth rule would put indentation arithmetic inside a predicate whose whole point is stating one condition. The old code was strictly worse: it trimmed the whitespace away before testing.
- **D-03** — One irreducible enumeration remains: the digit case. **DECIDE & STATE.** An ordered-list marker is the only Markdown block opener that can begin with a character prose also begins with, so it cannot be decided from the first byte. It stays a regex, documented as the single exception. The punctuation set is data (four ranges), not ten code paths.
- **D-04** — No Markdown parser. **DECIDE & STATE.** `do-work-cli`'s `go.mod` declares zero dependencies and `go.mod`/`go.sum` are outside this write set. Beyond that, goldmark would answer the narrower question: it parses `+++`, `:::note` and `$$` as ordinary paragraphs, so a parser-backed predicate would inline all three — the wrong direction under D-01.

## Discovered Tasks

- `internal/publication/capture_files.go:59-66` validates and contains raw input bytes but never asks the structural question about any one-line field it inlines. Whether capture inlines a user-supplied line anywhere was not audited.
- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/update-suite/INT` failed once under full-module load and passed on re-run and `-count=3`. Same flake already tracked as REQ-525; this is a second sighting, on the `update-suite/INT` subtest rather than `install-suite`.

## Qualification

**Passed** — 2 files verified, 6 requirements traced, P-A-U confirmed.

Mechanical: `tools/checks/qualify.sh` → `OK: mechanical qualification passed`; `tools/checks/scope-drift.sh` → `OK: Implementation Summary matches the Scope declaration`.

Independent (orchestrator-run, not the builder's report):
- `go build ./... && go vet ./...` clean; `gofmt -l .` printed nothing; `go test -count=1 ./internal/publication/` → `ok … 18.6s`.
- Read the full `answer.go` diff. The predicate is three tests behind a doc comment that states the condition and the doubt-resolution direction; `isMarkdownBlockPunctuation` is four range comparisons, so the punctuation set is data rather than ten branches. No stub, no placeholder, no debug artifact.
- Requirement trace: the boundary is now the condition (block-significant leading whitespace, ASCII punctuation, or an ordered-list digit run) rather than a prefix list; structural summaries route to the file-backed payload through the unchanged `BuildAnswerPlan` seam at `answer.go:117`; bytes are lossless on both paths because neither the inline nor the contained writer changed; the tests are table-driven with each row named by the class it represents; and `validateOutsideBytes` (C0/DEL) and `containedOutsideBytes` (multiline) in `publication_manifest.go` are untouched, which I confirmed against the diff rather than the report.
- The removed `strings.TrimLeft(summary, " \t")` is the fix for the indented-code class, not a regression: trimming ran *before* the test, so the bytes that constitute the structure were discarded before anything could see them.

## Testing

**Tests run:** `go build ./... && go vet ./... && gofmt -l .`; `go test -count=1 ./internal/publication/`; `go test ./...` for the module; canonical repository gate `bash _dev/tests/maintainer-verify.sh`.
**Result:** ✓ `internal/publication` green (`ok … 18.6s`). Gate exits 1 on the recorded baseline failure only — `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback`, tracked as REQ-524 and in a package this diff does not touch.

**Red-green validation:** traces the REQ's Red-Green Proof, which asked for `***`, `___`, and a representative from every Markdown structural class submitted as inline summaries without a raw payload.

Under the old ten-prefix predicate, 19 subtests fail. By class:
- Thematic breaks the list missed: `***`, `___`. (`---`, `* * *`, `- - -`, `-----` were already caught by the old `---`/`* `/`- ` prefixes and stayed green throughout — worth recording, since it means only part of that class was ever broken.)
- Setext heading underlines: `===`, `=` — absent from the old list entirely.
- Bare list markers: `-`, `*`, `+` — the old prefixes each required a trailing space.
- Bare ordered markers: `1.`, `1)`, `1234.` — the old pattern required a trailing space or tab.
- Indentation as structure: four leading spaces, one leading tab — the old code trimmed these away before testing, so the class was unreachable.
- Tables: `| option | cost |`, `|---|---|`.
- Dialect fences: `+++` (TOML frontmatter), `:::note` (directive/admonition), `$$` (math block).

All 8 prose rows pass at both ends, so the negative side is not an artifact of the new predicate.

**Caller-seam evidence:** `***` as an inline summary was accepted by `BuildAnswerPlan` with no refusal and written into the document as `- [x] Choice? → ***`. It now refuses with `ANSWER-RAW-PAYLOAD-REQUIRED`, and with a matching payload lands as `> ***` inside the canonical fenced blockquote with the checkbox reading `See contained answer note`. This is the finding's actual harm — an unescaped thematic break written into the REQ file it was being embedded in — reproduced and then closed.

**New tests added:**
- `TestSummaryContainmentDecidesByMarkdownStructureCondition` — 31 structural rows plus 8 prose rows, each row named by the class it stands for rather than by its input, per the REQ's requirement that examples are fixtures and not the definition of the accepted set.
- `TestBuildAnswerPlanCarriesStructuralSummaryContainedAndProseInline` — three subtests through the `BuildAnswerPlan` seam, so the predicate is pinned at the caller and not only in isolation.

**Existing tests updated:** none. No prior test asserted the ten-prefix behavior, which is itself worth noting — the enumeration was never pinned as a contract, only as an implementation.

*Verified by work action*

## Review

**Overall: 89%** | 2026-09-03T03:05:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 88% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 — `allResolvedQuestionsMatch` (`internal/publication/answer.go:410`) tests `bytes.Contains(line, marker)` across the whole `- [x] ` line, so a plain answer summary containing `→ Discarded:` makes an *answered* question read as discarded and the REQ is silently cancelled and archived; the `→ Confirmed:` variant reaches `completed` the same way. Reproduced in a scratch fixture. **Pre-existing — the old ten-prefix predicate inlined the same text — so not attributed to this diff**, but squarely on the contract this REQ redefined — `impact-critical` → REQ-528 created
- F2 — `cancellationReasonBlock` (`internal/requeststate/state_apply.go:895-903`) makes the same inline-or-contain decision on a newline test alone, so a one-line reason of `***` or `## Notes` is written in as document structure; that package already has its own containment writer and validator, only the structural judgment is missing — `impact-rule-change` → REQ-529 created (`pending-answers`: whether the two seams share one predicate is a package-direction call)
- F3 — restatement drift against `skills/do-work/actions/clarify.md`: `:105` conditions inlining on "cannot be a delimiter **where it lands**" while the code is position-blind, and `:106` says the containment branch keeps the summary after the `→` when the code writes `See contained answer note`. Also `_dev/tests/contract-regressions.sh:6605` requires the condition sentence exactly once across `skills/`, and this diff added a second wording of it in a Go doc comment — `impact-rule-change` → `do-work/prose-backlog.md` (prose-only, no root-cause match)

**Minor findings:** 4 (report only) — M1 the `orderedListMarkerPattern` comment asserts a false universal (Pandoc `fancy_lists` opens ordered lists from letters and roman numerals, all still inlined, correctly); M2 dropping the trailing-space requirement over-catches `1.2.3`, `3.50`, `2026.` — consistent with D-01 but unrecorded by any test row; M3 no row exercises the empty/whitespace branch or the intentional over-catch; M4 the two stated rules contradict on paper for an all-whitespace summary, resolved silently by ordering and unreachable at the seam. M1 routed to the prose backlog with F3.

**Acceptance:** Pass — package green at exit 0, revert-and-confirm independently reproduced at exactly 19 red subtests, both document shapes exercised end-to-end, and a 160,906-input differential of old against new predicate found **0 weakenings**, proving the new predicate is a strict superset that cannot bypass anything the old one caught.
**Suggested testing:** 4 items — a `→ Discarded:` spoof row at the `BuildAnswerPlan` seam once REQ-528 lands, a cancel/abandon fixture for REQ-529's RED, a row recording the intentional over-catch so a future narrowing is visible, and a human read of one real clarify round trip on an answer starting with a quote or paren.
**Follow-ups created:** REQ-528, REQ-529; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Stating the condition in the doc comment *before* writing the body. The comment names the three line-start ingredients Markdown builds every block construct from, and the body is three tests against them — so the code reads as the contract rather than as a coincidence.
- Taking the punctuation set wholesale from CommonMark's own ASCII-punctuation definition instead of the marks today's syntax uses. That is what catches `+++`, `:::note` and `$$` with nothing enumerating them, which is exactly what the REQ asked for by name.
- A differential run of the old predicate against the new one over 160,906 inputs. It turned "this should be a superset" into "0 weakenings, proven", which is a much stronger claim than any fixture set gives.

**What didn't:**
- My stated justification for the widening was wrong, and the review caught it. I argued containment prevents "an unescaped delimiter written into the document"; a summary is only ever written mid-line, after `→ ` or after `: `, so it never starts a line and that hazard is unreachable at both write sites. `clarify.md:105` already said so. The decision survives as defence-in-depth for a future line-start writer, but it was defended on a hazard that does not exist here.
- The sweep stopped at the prose that looked authoritative. It found `clarify.md` and concluded the prose already stated the condition correctly — but `clarify.md` is position-qualified and the code is position-blind, so the two now disagree in the opposite direction, and the sibling seam in `internal/requeststate` was missed entirely.
- Only part of the thematic-break class was ever broken: `---` and `* * *` were caught, `***` and `___` were not. A partially-correct enumeration is worse than an obviously-empty one, because the passing cases make it look like the class is handled.

**Worth knowing:**
- A summary is written at exactly two sites (`answer.go:196`, `:200`), both mid-line. Any future reasoning about delimiter hazards here has to start from that, not from "a delimiter is dangerous".
- `answer.go:136` already refuses an empty summary using the identical `TrimSpace(summary) == ""` expression the predicate's guard uses, so that guard is unreachable in production. Worth knowing before anyone "fixes" the apparent contradiction between it and the leading-whitespace rule.
- Pandoc's `fancy_lists` opens ordered lists from letters and roman numerals. `A) yes, use the first option` stays inline, and that is the right call — but it means the digit case is the opener this predicate *chooses* to cover, not the only one that exists.

## Orientation

An answer summary is now carried as a file-backed payload whenever Markdown could read it as document structure, decided by one stated condition rather than a list of ten prefixes. Lives in the CLI's `publication` package, on the `answer` path that `do-work clarify` drives.

`prime_files`: `_dev/primes/prime-action-files.md` — spot-checked, its referenced paths all still exist and it needed no change; the contract this REQ implements is stated in `skills/do-work/actions/clarify.md`, which the review found now disagrees with the code in two places (both on the prose backlog, not silently left).

**[MAP CHANGED]** — no new module, but the same inline-or-contain decision is now condition-complete in one package and absent in its sibling (`internal/requeststate`), which is a split REQ-529 exists to close.
