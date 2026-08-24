---
id: REQ-360
title: "[impact-rule-change] Close the neutralization contract's reach"
status: claimed
claimed_at: 2026-08-24T17:42:55Z
status_changed_at: 2026-08-24T17:42:55Z
route: C
created_at: 2026-08-24T10:50:00Z
user_request: UR-068
addendum_to: REQ-342
domain: security
review_generated: true
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
sweep: true
sweep_key: neutralization-contract-reach
impact: impact-rule-change
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-08-24T17:42:55Z
  basis:
    - Route C
    - 10-file seeded write set
    - 3 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
  - _dev/tests/contract-regressions.sh
  - skills/do-work-board/tools/queue-kanban/frontmatter_test.go
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/verify-requests.md
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/capture-reference.md
  - skills/do-work/actions/stakeholder-answers.md
  - skills/do-work/actions/abandon.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/sample-archived-req.md
  - skills/do-work/docs/capture-guide.md
  - skills/do-work/docs/work-guide.md
---

# Close the Neutralization Contract's Reach

## What

Two sibling contracts — REQ-342's body-side neutralization rule and REQ-344's Frontmatter Quoting —
each landed at a named entry point without closing its reach. Five ways, sharing one root cause: the condition quantifies over **lines**, so damage that is not
line-shaped escapes it; and one live caller writes a user's answer into a REQ body without citing the
format at all, so it inherits nothing.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

One root cause: the contract was stated at the entry point without closing either end of its reach —
neither every damaging shape, nor every writer.

## Instances

- [ ] **The condition is line-scoped, so non-line-shaped damage escapes it.** The rule reads "Judge
  every line of it against one condition: could this line be read as one of this file's own
  delimiters?" A byte that changes how the **file** is read rather than how a **line** is read is
  outside it, and neither the `> ` prefix nor the over-long fence neutralizes it. Measured on an
  answer containing a NUL, written exactly per branch 2:

  | scan | control | with NUL |
  |---|---|---|
  | `grep -cE '^[[:space:]]*- \[ \]'` | 3 | 3 |
  | `grep -nE '^[[:space:]]*- \[ \]'` lines printed | 3 | **0** |
  | `grep -nE '^## '` lines printed | 3 | **0** |
  | board renders `## Plan` | yes | yes |

  grep reports "binary file matches" on stderr and prints nothing on stdout. Any reader consuming
  matched *lines* — Step 5's "any remaining unchecked item wins", a `while read` loop — sees an empty
  result on a REQ that has two real open questions, while the board shows them. That is the REQ's own
  second consequence with a different first cause. `_dev/primes/prime-shell-commands.md` §
  *Prescribed Shell Commands Must Surface What the Steps Consume* already records the trap.
  Judged `impact-rule-change`.

- [ ] **`actions/verify-requests.md:168` writes a user answer and cites nothing.** Inside Step 7's
  "the user is here right now — resolve them on the spot": `2. If the user answers → add the resolved
  question to the REQ's `## Open Questions` section as `- [x] [question] → [user's answer]``. It
  obtains typed text interactively and writes it into a REQ body in exactly the shape the format
  governs, and the string "Canonical answered-question format" appears nowhere in that file. Line 169
  is the same site's deferred branch. Judged `impact-rule-change`.

  REQ-342's builder did what its REQ asked — "grep every current citation of the format and confirm
  each inherits" — and a citation grep is structurally blind to a caller that never cites. The UR's
  wording is the stronger one: "so every caller that records an answer obeys it, not just clarify."

- [ ] **The UR `input.md` template prescribes an unquoted title.** `actions/capture-reference.md:194`
  shows `title: Add keyboard shortcuts`. A UR title is the user's own phrasing, so the Frontmatter
  Quoting condition governs it. Measured through the shipped CLI: `title: Fix the #1 board complaint`
  reads back as `Fix the`, exit 0, no warning. Note the reach gap behind it — the contract's home is
  § **Request** File Schema while UR files run through the same parser. Judged `impact-user-visible`.

- [ ] **The block-scalar branch has an unstated precondition.** Written exactly as prescribed
  (`key: |-`, text indented beneath), a value whose **first line begins with a space** fails the
  strict parse, drops the whole block to salvage, and reads back as the literal `|-`:

  ```
  blocked_by: |-
     staging creds provisioned
    by the platform team
  → blocked_by = "|-"   exit=0   (estimate: dropped)
  ```

  YAML needs an explicit indentation indicator (`|-2`) there, and pasted text with an indented first
  line is the ordinary case for a multi-line `blocked_check`. Also measured: a lone CR folds to a
  space under the single-quoted form and fails the strict parse under the block form, and a control
  character (0x1b, NUL) fails under both. So the contract's opening claim — "no character in it can
  be read as YAML syntax" — over-promises on what the two forms deliver. Judged `impact-rule-change`.

- [ ] **Four shipped files still demonstrate the forbidden form, and one write site was missed.**
  `actions/review-work.md:361` (core's own follow-up minting template, the exact sibling of the
  `code-review.md` line REQ-344 corrected); `docs/capture-guide.md:34`, five lines above the sentence
  REQ-344 rewrote, so that file now contradicts itself; `docs/work-guide.md:123`, the snippet the
  guide tells a user to copy; `actions/sample-archived-req.md:3`, unquoted. And
  `actions/stakeholder-answers.md` Step 5 rewrites `blocked_by:` with a person's name and received no
  citation, where `work.md` and `clarify.md` did for the identical rewrite. Judged
  `impact-rule-change`.

## Detailed Requirements

- The neutralization condition covers damage that is not line-shaped, or states explicitly and in one
  place what it does not cover and why that residue is acceptable. Either is a legitimate answer;
  silence is not.
- Every action that writes user-obtained text into a REQ body cites the format. **Find them by their
  behaviour, not by their citations** — the citation grep is what missed this one. Grep for the write
  itself: `- [x]`, `→`, "the user's answer", "resolve them on the spot".
- The lock-in gains a case for whichever instances land, so the reach is pinned rather than argued.

## Constraints

- Do not restate the rule in `verify-requests.md`. It cites; the entry point states. REQ-342's
  exactly-once count check enforces that and must keep passing.
- `_dev/primes/prime-action-files.md` governs both action files. Read it first.

## Red-Green Proof

**RED prompt/case:** A REQ carrying a NUL inside a correctly-quoted, correctly-fenced answer returns
zero lines from `grep -nE '^[[:space:]]*- \[ \]'` while the board renders its two real open questions.
Separately, `grep -c "Canonical answered-question format" skills/do-work/actions/verify-requests.md`
returns 0 while that file writes a user answer into a REQ body at line 168.

**GREEN when:** the non-line-shaped case is either neutralized or documented as out of reach with its
reason, `verify-requests.md` cites the format, and a lock-in covers both.

**Validation:** Inferred during REQ-342's review; the NUL measurement and the uncited caller were both
reproduced there.

---
*Source: REQ-342 review findings F1 and F2 (UR-068), folded into one sweep on the shared root cause.*

## Triage

**Route: C** — This five-instance security/rule sweep spans body neutralization, YAML frontmatter, writer reach, examples, and contract lock-ins across thirteen files. It needs an explicit plan and exploration before the shared contract is revised.

## Plan

1. Add RED semantic and mutation lock-ins for non-line-shaped text damage, each uncited body writer, every live forbidden frontmatter example, and byte-preserving multiline scalar cases.
2. Broaden the single canonical body-neutralization condition so every do-work Markdown writer inherits delimiter containment plus an accepted-text precondition.
3. Wire every behaviorally discovered body writer and correct all copyable UR/REQ/frontmatter examples without restating either canonical contract.
4. Make the Frontmatter Quoting block-scalar branch byte-honest for indentation, terminal LF count, and invalid controls; pin strict parsing and nested frontmatter survival in Go.
5. Run focused mutation/contract/parser tests and the canonical repository gate, then audit every sweep instance and the exactly-once canonical-condition count.

**Plan validation:** All five sweep instances map to a closure and at least one falsifiable lock-in; no task is orphaned. Five tasks is the upper safe bound, but the shared contract/writer/test structure makes further splitting more likely to recreate reach drift than reduce it.

## Exploration

- The canonical answered-question format in `clarify.md` is the only body-neutralization entry point. Generalize its condition to externally supplied text written into do-work Markdown records while keeping the answer-line plus dated-reasoning record shape answer-specific.
- `verify-requests.md` is the named missed writer. Behavior-shaped search also found queued-addendum/UR input writes in `capture.md`, a new-UR verbatim reply in `stakeholder-answers.md`, and a cancellation reason in `abandon.md`; each must actively cite the same condition.
- Five live hand-authored frontmatter examples need single quotes and inheritance: the UR title in `capture-reference.md`, the guide title in `capture-guide.md`, the review follow-up title, `sample-archived-req.md`'s title, and `work-guide.md`'s `assigned_to`. Encoder-emitted board YAML remains the explicit exception.
- `work-reference.md` must require accepted text to contain only LF/TAB among C0 controls, indent every physical block line, and select `|-`, `|`, or `|+` for zero, one, or multiple terminal LF bytes. Strict-parser fixtures belong in `frontmatter_test.go` so salvage cannot masquerade as success.
- Contract regressions must isolate each writer block, mutation-test canonical wording and reach, inventory forbidden hand-authored forms, and keep one canonical-condition occurrence so callers inherit rather than restate.

## Scope

**Files I will touch:**

- `_dev/tests/contract-regressions.sh`
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go`
- `skills/do-work/actions/clarify.md`
- `skills/do-work/actions/verify-requests.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work/actions/capture-reference.md`
- `skills/do-work/actions/stakeholder-answers.md`
- `skills/do-work/actions/abandon.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/sample-archived-req.md`
- `skills/do-work/docs/capture-guide.md`
- `skills/do-work/docs/work-guide.md`

**Acceptance criteria:**

- The one canonical body-neutralization condition covers delimiter-shaped and non-line-shaped damage; invalid C0/DEL controls other than LF/TAB are refused and reported before any write.
- Every behaviorally discovered action that records external text into a do-work body actively cites the canonical condition and preserves its existing state semantics.
- Every live hand-authored frontmatter example uses the single-quoted form and cites Frontmatter Quoting; complete/verbatim input remains byte-identical apart from containment bytes.
- Multiline frontmatter preserves leading indentation and exact terminal LF count through strict YAML, while invalid controls are rejected before hand composition and nested `estimate:` remains intact.
- Semantic/mutation lock-ins fail on narrowing, passive/missing citations, forbidden examples, bad chomp/indentation, or duplicate/restated canonical conditions; the full gate stays green.

## Decisions

- **D-01 — Refuse invalid control bytes rather than normalize them.** To keep the existing COMPLETE/UNEDITED and byte-identical promises honest, externally supplied text containing C0/DEL controls other than LF/TAB is not written; the caller reports the offending code point. CR-to-LF normalization and a hand-authored escape table are rejected because both silently change the supplied bytes.
