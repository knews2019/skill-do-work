---
id: REQ-342
title: "[impact-critical] Neutralize user text written into a REQ body"
status: completed
completed_at: 2026-08-24T10:52:00Z
claimed_at: 2026-08-24T08:55:00Z
created_at: 2026-08-23T22:35:07Z
user_request: UR-068
domain: security
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-343, REQ-344]
maintenance: false
route: C
impact: impact-critical
effort_estimate: effort-substantive
write_set:
  - skills/do-work/actions/clarify.md
  - _dev/tests/contract-regressions.sh
---

# Neutralize User Text Written Into a REQ Body

## What

The Canonical answered-question format writes a user's typed answer verbatim into a REQ body, and
nothing states how that text is neutralized first. Arbitrary typed text can therefore forge the
body's own delimiters. State the neutralization rule at the named entry point — `actions/clarify.md`
Step 4 — so every caller that records a user answer obeys it, not just clarify.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-action-files.md` and the prime's § Closed Enumerations Go Stale. Compared four mechanisms against the two reader classes the file's own history names — the board (goldmark) and line-based prose greps — and verified by experiment that fence-only defeats one but not the other.
- [x] **[APPLY]:** Two files, both inside the write set. The rule is stated once inside the named entry point's blockquote; the two other touches in `clarify.md` are pointers, and the lock-in's exactly-once count check enforces that.
- [x] **[UNIFY]:** Audited by the orchestrator against the merged range `ec24585..2407f27`: mutated the shipped rule by replacing "illustrative" with "the complete set to check" — the narrowing the REQ forbids — and confirmed the lock-in fails with exit 1 naming that exact property; restored byte-identical. Separately confirmed the mechanism defeats a line-based scan while the file's own headings survive.

## Why

Three consequences, all proven on the user's fixture, and none of them announce themselves:

- A typed answer containing a line-leading `- [ ]` becomes a real unchecked open question. Clarify
  Step 5's "any remaining unchecked item wins" then pins that REQ in `pending-answers` forever and
  re-presents the pasted fragment at every session.
- An unclosed code fence swallows every following section into a code block on the board, while the
  actions' prose greps still see those sections — so the board and the pipeline disagree about what
  the REQ says.
- A paste landing above the opening frontmatter fence loses the entire frontmatter. Status, title and
  `user_request` all read empty, so the REQ silently drops out of its UR.

Pasting a code snippet as an answer is ordinary use, not an attack, and an unbalanced fence is the
common accident. That is what makes this `impact-critical` rather than a hardening nicety: the queue
loses a REQ's identity and its UR membership without erroring.

## Context

`actions/clarify.md` Step 4 is the **Canonical answered-question format** — a named entry point that
`actions/work.md` Step 3.5 already cites for a mid-run answer, so a rule stated there is inherited
rather than copied. Verified: the file today contains no escaping, neutralization or fence-balancing
language at all. Its only "verbatim" mentions are about rewriting stored *question* wording for a
cold reader (`:38`, `:216`, `:231`), which is a different concern in the opposite direction.

This REQ was captured immediately after a live instance of the same class in this repo: a clarify
write keyed its regex on the file's first line-leading `- [ ]`, which is the P-A-U `[PLAN]` checkbox
rather than the question, and falsely ticked `[PLAN]` on three REQs. That was an agent's own text,
not the user's, and it was caught and reverted — but it is the same failure mode from the other side:
this format's delimiters are ambiguous, and neither the writer nor the reader is told so.

## Detailed Requirements

- The Canonical answered-question format states how user text is neutralized before it is written.
- The rule covers at minimum: a leading `- [ ]` or `- [x]`, a bare `---` line, a `## ` heading, and an
  unbalanced code fence.
- **State the condition, not a closed list of characters.** The four above are the shapes the fixture
  proved; the rule must be keyed on "text that could be read as one of this file's own delimiters"
  so a fifth shape is covered without an edit (`_dev/primes/prime-shell-commands.md` § Closed
  Enumerations Go Stale).
- Stated **once**, at the named entry point, and cited by callers rather than restated. Grep every
  current citation of the format and confirm each inherits rather than duplicating.
- The neutralization must preserve what the user actually said — a reader must still be able to see
  the answer's real content. A rule that silently drops characters is not acceptable.

## Constraints

- `_dev/primes/prime-action-files.md` governs any action-file change. Read it first.
- Do not weaken clarify Step 5's "any remaining unchecked item wins" rule — it is correct, and the
  defect is that forged text can create an item it then acts on.
- Scope is the writing rule and its lock-in. Repairing REQ files already damaged by this is REQ-343's
  detection half plus a separate repair decision, not this REQ.

## Builder Guidance

**Certainty: firm on the defect and the entry point, open on the neutralization mechanism.** The user
proved all three consequences and named the minimum shapes; how to neutralize them (indent the
pasted block, fence it with a longer fence, prefix a zero-width guard, quote it as a blockquote) is
the builder's call. Prefer whichever keeps the answer legible to a human reading the archived REQ,
because that trail is the reason the answer is written down at all.

Note the tension this sits next to: `actions/capture-reference.md` § REQ Title Convention says the
parser's line-based recovery "is a salvage path with a narrower contract than the parser proper",
while `frontmatter.go:104` says "recovery is the contract here". REQ-344 covers the frontmatter half;
do not resolve that tension here, but do not lean on recovery either.

## Open Questions

None — the user stated the defect, the entry point, and the minimum shapes.

## Red-Green Proof

**RED prompt/case:** Record an answer whose text contains a line-leading `- [ ]` in a REQ that has no
other open questions, then run `do-work clarify`: the REQ is still `pending-answers` and the pasted
fragment is presented as a question. Separately, record an answer containing one unbalanced code
fence and open the board: every section below it renders inside a code block.

**Why RED now:** Nothing in the Canonical answered-question format says the text is neutralized, so
it is written through as-is.

**GREEN when:** Both pastes are recorded with their content intact and neither forges a delimiter —
the REQ resolves out of `pending-answers`, the board renders the following sections normally, and a
semantic instruction test fails if the neutralization rule is removed or narrowed from the named
entry point.

**Validation:** User confirmed — the defect, the entry point and the minimum shapes are stated
verbatim in the input, from a fixture the user ran.

## Assets

None. Reproduced on the user's fixture; the three consequences are recorded in the UR input.

---
*Source: UR-068 — see `do-work/user-requests/UR-068/input.md` for complete verbatim input.*

---

## Triage

**Route: C** - Complex

**Reasoning:** A rule stated in prose at a named entry point, inherited by six callers, whose lock-in has to detect narrowing rather than absence. The mechanism was explicitly the builder's to choose against two different reader classes, and the wrong choice defeats one of them silently.

**Planning:** Required — the plan was the mechanism comparison, carried out by experiment rather than argued.

## Scope

**Files I will touch:**
- `skills/do-work/actions/clarify.md` (modify) — the neutralization rule inside the Step 4 blockquote, one Red Flag, one checklist clause
- `_dev/tests/contract-regressions.sh` (modify) — the semantic lock-in and the exactly-once count check

**Files I will NOT touch:** the six citing callers — they cite the format by name and inherit the rule; editing them would create the second copy the exactly-once check exists to prevent.

**Acceptance criteria (restated from REQ):**
- [x] The format states how user text is neutralized before it is written
- [x] Covers the four proven shapes: `- [ ]`/`- [x]`, a bare `---`, a `## ` heading, an unbalanced fence
- [x] Keyed on the condition, not a closed character list — shapes marked illustrative
- [x] Stated once at the named entry point; every caller inherits rather than restates
- [x] Neutralization preserves what the user actually said

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/clarify.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** The Canonical answered-question format now carries the neutralization rule, keyed on the condition "could this line be read as one of this file's own delimiters?" with the four proven shapes present as illustrative examples explicitly marked "not a checklist". A one-line answer stays inline after the arrow, because a delimiter must start a line and inline text never does. Anything longer is summarized on the answer line and placed in the dated note inside a blockquote whose lines open a code fence longer than the longest backtick run in the text: the `> ` prefix removes every line start from a line-based scan, the longer fence keeps those lines literal for a Markdown reader and cannot be closed from inside, and nothing but the prefix is added — so the answer's content survives verbatim. Placement is bound by the same condition: never the file's first line, which is the one delimiter no container can guard.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`, `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Gate exit 0

**Red-green validation (all three consequences, board-rendered):**
- Forged open question: RED board render `{"total":3,"unchecked":1}` → GREEN `{"total":1,"unchecked":0}` — the forged unchecked item that would pin the REQ in `pending-answers` is gone, the real `Q-01` remains
- Unbalanced fence: RED `bodyHtml contains <h2 ...Plan>? false` → GREEN `true`; the user's unbalanced fence is preserved exactly and rendered inert
- Paste above the frontmatter fence: RED `status=[... has no frontmatter block]` → GREEN `status=[pending] title=[...] user_request=[UR-999]`
- Content preserved: the pasted `- [ ] retry once` and `- [x] then give up` are still legible in the archived REQ, quote-prefixed

**Semantic lock-in (the REQ's hardest clause):** the test detects the rule being *narrowed*, not merely deleted. Independently re-run by the orchestrator — replacing "illustrative" with "the complete set to check" produces `FAIL: ... keyed on the delimiter condition rather than a character list (REQ-342)` and exit 1. The in-suite matrix replays 12 further narrowings, each required to trip its own predicate.

**Caller sweep:** six live citations, none restating the rule — `actions/stakeholder-answers.md:41,102`, `actions/work.md:253`, `actions/work-reference.md:755`, `crew-members/clear-questions.md:41` in all three packages. Enforced mechanically by the exactly-once count check.

*Verified by work action*

## Decisions

<!-- D-XX counter: last used D-05. Next decision: D-06. -->

- **D-01 — Blockquote plus an over-long code fence, not fence-only, indent-only or blockquote-only. DECIDE.** Decided by experiment, not preference: with the answer fenced but unquoted, `grep '^[[:space:]]*- \[ \]'` still matched the pasted checkbox — fence-only defeats a Markdown reader and not a line-based scan. Blockquote-only defeats the scan but still renders a task checkbox and a heading inside the quote. The file's own history names both reader classes: the board is goldmark, and the sibling incident the REQ cites was a regex keying on the first line-leading `- [ ]`.
- **D-02 — The answer line carries a one-line summary; the verbatim text lives in the dated note. DECIDE.** The checkbox line stays scannable and the exact words stay one line away, rather than a bare "see below".
- **D-03 — Placement is part of the same condition, not a separate rule. DECIDE.** The opening frontmatter fence is a delimiter like the others, and no container can guard it, so placement *is* its neutralization. This deliberately does not lean on the parser's line-based recovery: `splitFrontmatter` returns `hasFrontmatter=false` outright when the fence is not the first line, so there is nothing there to lean on — which is what the REQ's Constraints required.
- **D-04 — One Red Flag and one checklist clause rather than a new checklist entry. DECIDE.** Both are pointers to the named format, carrying neither the condition nor the mechanism; the exactly-once count check would fail on a restatement.
- **D-05 — The lock-in is an extractor plus mutation matrix, not `assert_block_contains` calls. DECIDE.** `require(file, token)` tests vocabulary, and the requirement was to detect narrowing. Isolating the blockquote also stops nearby prose from lending vocabulary to a weakened property.

## Discovered Tasks

- `skills/do-work/actions/work.md:253` cites the format by name and also restates the `- [x] [question] → [the user's answer]` shape inline. Harmless today — the neutralization rule is inherited, not copied — but it is a second copy of the answer-line form that can drift from the entry point.
- `_dev/tests/contract-regressions.sh` has an existing `grep … | wc -l` count under `set -euo pipefail` (`three_attempt_count`, near the Step 6.5 block) with no `|| true`. If that phrase ever reaches zero occurrences the whole suite aborts at that line and reports nothing about the checks below it. The builder hit exactly this while writing its own count check and guarded its own; the older one is unguarded.
- `queue-kanban frontmatter get` prints its "has no frontmatter block" message to stdout in a way that reads like a value to a shell caller. Adjacent to REQ-344.

## Open Questions

- [~] Do the Red Flag line and the checklist clause count as restating the rule rather than citing it? → **D-06**: Builder judged them pointers — neither carries the condition or the mechanism, and the exactly-once check confirms the condition sentence exists in one file only. Orchestrator concurs: the check is mechanical and it passes. Value: the symptom is named where an operator meets it. Risk: if the maintainer reads them as restatement, both are one-line deletions. Reversible; not carried to a follow-up REQ, because the mechanical check already enforces the property the question is about.

## Review

**Overall: 86%** | 2026-08-24T10:36:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 88% |
| Code Quality | 85% |
| Test Adequacy | 70% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:**
- **The condition quantifies over lines, so damage that is not line-shaped escapes it.** A NUL inside a correctly-quoted, correctly-fenced answer makes every `grep -n` on that REQ print zero lines while `grep -c` still counts and the board renders normally — the same board/pipeline disagreement this REQ calls impact-critical, reached through a shape the rule does not cover. — `impact-rule-change` → **REQ-360**
- **`actions/verify-requests.md:168` writes a user's typed answer into a REQ body and never cites the format**, so it inherits nothing. The builder's sweep grepped citations, which is structurally blind to a caller that does not cite. — `impact-rule-change` → **REQ-360** (same sweep; one root cause — the contract landed without closing either end of its reach)
- **The lock-in detects deletion and phrase substitution, not narrowing by added qualifier.** Six narrowings that keep all 12 phrases byte-intact pass unchanged, including "Judge its first line" and making branch 2 conditional on length — two that reintroduce the original defect outright. — `impact-rule-change` → **REQ-361**

**Minor findings:** 3 (report only) — the lock-in's own comments claim it "grades by meaning" when it is phrase matching; `clarify.md:160`'s Step 5.5 reclaim writes a user answer citing the format only indirectly; `stakeholder-answers.md` Step 4.1 preserves a third-party reply verbatim into a UR `input.md`, the same class of write in a different file type.

**What the reviewer could not defeat, stated plainly:** the mechanism held on every line-shaped input tried — 15 adversarial cases, each built as a real REQ file, scanned three ways and rendered through a real board build. Text already carrying `> ` prefixes; a line that is exactly `> ` plus five backticks; a 12-backtick run (the fence scales to 13, and runs inside quoted lines count, so no length wins); the dated-note marker followed by a forged checkbox; CRLF throughout; an empty first line; mixed tilde and backtick fences; Unicode look-alikes; and the plain nasty paste with all four shapes at once. Branch 1's inline justification also survived: a literal newline routes to branch 2 by the rule's own trigger, and a lone CR does not split under goldmark, grep, or Go's `strings.Split`.

**Verified claims:** D-03 holds — `splitFrontmatter` (`frontmatter.go:36-42`) accepts a fence only at offset 0 after an optional BOM and otherwise returns `hasFrontmatter=false` immediately, so there is no recovery path to lean on. The condition is genuinely condition-keyed rather than a disguised enumeration. The exactly-once claim holds for what it measured: six citing files, each a pointer carrying neither condition nor mechanism, and the count check works (adding a second copy moved it 1 → 2).

**Acceptance:** Pass — gate exit 0 with browser and JS lanes both run; all three consequences independently reproduced GREEN on real board renders.

**Follow-ups created:** REQ-360 (sweep), REQ-361

*Reviewed by review-work action*

## Correction to This REQ's Own Testing Section

This REQ's Testing section above claims the lock-in "detects the rule being *narrowed*, not merely
deleted". **That claim is too strong and the orchestrator repeated it before checking.** It detects
deletion and phrase substitution. Narrowing by added qualifier escapes: replacing "Judge every line
of it" with "Judge its first line" — which lets a delimiter on line five of a paste straight through —
ships with the suite at exit 0 and no complaint. Independently reproduced 2026-08-24.

The orchestrator's earlier verification tested one mutation, of the substitution kind the matrix was
built around, and generalized from it. REQ-361 carries the fix; the escape table is recorded there as
its RED set.

## Lessons Learned

**What worked:** Choosing the mechanism by experiment rather than by argument. Fence-only *looks*
sufficient until you run `grep '^[[:space:]]*- \[ \]'` against a fenced-but-unquoted answer and watch
it match — the builder ran exactly that and it decided the design. Isolating the named entry point's
own blockquote before grading it is the one part of the lock-in that is genuinely more than
`require(file, token)`: it stops nearby prose lending vocabulary to a weakened property.

**What didn't:** The lock-in's mutation matrix replaces the exact phrase each predicate matches, so it
proves each predicate catches its own phrase's deletion — close to a tautology. The axis it holds
constant, insertion, is the one an author actually narrows a rule along: a hedge, a scope limit, a
length threshold. Six such mutations ship green. A mutation matrix built by inverting your own
predicates tests your predicates, not the property.

**Worth knowing:** The condition's subject is a *line*, and that is load-bearing in a way the prose
does not admit — a NUL byte defeats every line-based reader while leaving the board correct, which is
this REQ's own second consequence arriving through a door the rule does not watch. Also:
`insertQuestionOptionHardBreaks` (`render.go:50`) tracks fences with a prefix check that never fires
on a `> `-prefixed fence, so it walks a quoted answer believing it is outside a code block. Harmless
today only because the same `> ` prefix stops those lines matching its other triggers — safe by one
guard twice, not by two independent guards.

## Orientation

An answer typed into `do-work clarify` can no longer forge the delimiters of the REQ file it is
written into — no invented open question that pins the REQ in `pending-answers`, no unbalanced fence
that swallows the sections below it on the board, no paste above the frontmatter that costs the REQ
its identity and UR membership. The rule lives at the Canonical answered-question format in
`actions/clarify.md` Step 4, and the six actions that record an answer inherit it by citation.
