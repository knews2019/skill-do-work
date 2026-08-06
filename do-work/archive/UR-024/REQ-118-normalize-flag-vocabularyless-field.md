---
id: REQ-118
title: The normalize flag must stop calling vocabulary-less field values unrecognized
status: completed
created_at: 2026-08-06T10:53:03Z
claimed_at: 2026-08-06T11:10:18Z
completed_at: 2026-08-06T11:14:00Z
commit: 8d1a9f2
route: A
kb_status: pending
user_request: UR-024
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/frontmatter_cli.go, tools/queue-kanban/frontmatter_cli_test.go, tools/queue-kanban/model.go]
maintenance: false
related: [REQ-116, REQ-117]
batch: schema-contract-board-fixes
---

# The Normalize Flag Must Stop Calling Vocabulary-Less Field Values Unrecognized

## What

`queue-kanban frontmatter get <file> <field> --normalize` emits a `⚠ … not recognized — no canonical vocabulary is defined for this field.` line on **every** call when the field has no Schema Read Contract row — a timestamp, a title, a path list. Gate the warning on the field actually having a row, so the flag becomes a clean no-op for fields the contract itself places outside its scope.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** The root cause is that `isKnownSchemaFieldValue` returns `false` for two different facts — a bad value, and a field the contract does not govern — so the caller cannot tell them apart. Add the missing predicate (`hasSchemaFieldContract`, a map lookup, never a hand-listed exemption set) and gate `runFrontmatterCommand` on it: no contract row means print the value and return, with `--in-set` split off as a usage error. Reword `schemaFieldWarningText`'s contract-less branch to match the contract's classification rather than deleting it.
- [x] **[APPLY]:** Four test cases first, all RED for the right reasons (three showing the actual warning text, one showing `exit = 1, want 2`). Then the predicate in `model.go`, the gate in `frontmatter_cli.go`, the reworded branch. Three files, all in `write_set`.
- [x] **[UNIFY]:** `git diff --stat` → 3 implementation files. `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` passes (4.3s), which includes REQ-115's `TestRunFrontmatterCommandWarnsOnUnrecognizedStatus` and `…InSetWarnsOnUnrecognizedStatus` — the hole this must not re-open. No debug artifacts. Prose-caller sweep: `grep -rn 'frontmatter get' actions/ docs/ CLAUDE.md` → exactly one call site, `actions/commit.md`'s `status --in-set terminal-success`, whose field has a contract row and is untouched. Built the binary and exercised all four paths for real: `created_at --normalize` → value, silent, exit 0; `status --normalize` → `completed`, silent, exit 0; `status --in-set terminal-success` → exit 0; `created_at --in-set terminal-success` → usage error, exit 2.

## Why (if provided)

Stderr-only and harmless to `value=$(…)` capture, so this starts as noise. It is slightly more than noise, though: `actions/work-reference.md`'s Schema Read Contract states that **fields with no canonical vocabulary are outside this contract and are read verbatim** — "no alias map, no case folding, no path canonicalization, no warning." Printing "not recognized" for such a field tells the reader the opposite of the contract's own classification of it. The cause is that `isKnownSchemaFieldValue` returns `false` for a field with no contract row — the same answer it gives for a genuinely bad value — so the CLI's warn branch cannot distinguish "bad value" from "field the contract doesn't govern."

## Detailed Requirements

- In `runFrontmatterCommand`, gate the normalize/warn branch on the field having a contract row. For a row-less field: print the trimmed value on stdout, print **nothing** on stderr, exit 0 — identical to the same call without `--normalize`.
- The gate must be a lookup of the contract table, not a hardcoded list of exempt field names. A hand-enumerated exemption list is exactly the pattern `CLAUDE.md` → **Closed Enumerations Go Stale** forbids: it would silently go wrong the next time the contract grows a row.
- `--in-set` gets the same treatment for a row-less field, and must not answer a membership question it cannot answer meaningfully — both set names (`terminal-success`, `terminal-resolved`) are `status` sets, so asking for either on a row-less field is a caller error. Reporting it as a usage error (exit 2) is preferred over silently exiting 1; the builder decides which, and states the choice.
- **Behavior for fields that DO have a row must not change.** A typo'd `status` must still warn — that is REQ-115's fix (0.175.1) and this REQ must not re-open the hole it closed. Keep or extend `frontmatter_cli_test.go`'s existing coverage of it rather than editing those assertions to fit.
- Resolve the now-unreachable branch rather than leaving it. `schemaFieldWarningText`'s no-contract branch has exactly one caller (`frontmatter_cli.go`), so gating the caller strands it. Either delete it and let the formatter require a row, or keep it as a guarded invariant with a comment saying why it is unreachable — but do not leave dead code with a stale comment (`crew-members/coding-guardrails.md`, surgical changes).

## Constraints

- Read-only command; no write surface (see the UR's Batch Constraints).
- The stdout/stderr split is load-bearing and stays as documented in `frontmatter_cli.go`: the value goes to stdout and **nothing else does**, so a caller can capture stdout cleanly even when a warning fires.
- No change to `actions/commit.md` or any other prose caller — the single existing call site (`frontmatter get … status --in-set terminal-success`) uses a field that has a row, so its behavior is untouched.

## Dependencies

None. Independent of REQ-116 and REQ-117; touches different files apart from a possible one-function edit in `model.go`.

## Builder Guidance

**Firm on the silence, open on the mechanics.** Silent no-op is the decided behavior (see the Open Question below for what was considered and rejected). How the gate is expressed, and what happens to the stranded formatter branch, are the builder's calls.

## Open Questions

- [~] Should a row-less field with `--normalize` print nothing at all, or one reworded line saying the flag has no effect because the field is outside the contract? → **Deferred: proceeding with print nothing.** The user was asked and did not answer, so the reading most faithful to their original complaint wins — they called the per-call line noisy, and a reworded per-call line is still a per-call line. The accepted cost is that a caller who passes `--normalize` expecting normalization gets no hint it did nothing.
  Recommended: silent no-op (implemented).
  Also: one reworded note per call; silent on `--normalize` but a note on `--in-set`.

## Red-Green Proof

**RED prompt/case:** `queue-kanban frontmatter get <any REQ> created_at --normalize` — assert stderr is empty. It fails today, printing `⚠ created_at: '2026-08-05T15:53:39Z' not recognized — no canonical vocabulary is defined for this field.` while still exiting 0 with the correct value on stdout.
**Why RED now:** `isKnownSchemaFieldValue` returns `false` for a field with no contract row, which is indistinguishable at the call site from an unrecognized value, so the warn branch always fires.
**GREEN when:** `created_at --normalize` prints the value on stdout, nothing on stderr, exit 0; `status --normalize` on a typo'd status still warns (REQ-115's behavior intact); `route --normalize` still uppercases `a` to `A` silently; and `--in-set` on a row-less field reports a usage error rather than a silent false.
**Validation:** User adjusted — the user's capture text asked for both a gate and a rewording; the gate is implemented, and the rewording question was put to them and left unanswered, so the Open Question above records the assumption the build proceeds under.

## Assets

None.

## Triage

**Route: A** - Simple

**Reasoning:** One named file, one named gate, and the decided behaviour recorded in the Open Question. The only judgment is what happens to the stranded formatter branch.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

---

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/frontmatter_cli.go` (modified)
- `tools/queue-kanban/frontmatter_cli_test.go` (modified)
- `tools/queue-kanban/model.go` (modified)

**What was done:** Added `hasSchemaFieldContract(fieldName)` — a lookup of the contract table, the predicate that was missing — and gated `runFrontmatterCommand`'s normalize/resolve branch on it. A field with no row now prints its value and returns, observably identical to the same `get` without `--normalize`. `--in-set` on such a field returns a usage error (exit 2) naming the field, rather than a silent exit 1 that would read as a real negative. `schemaFieldWarningText`'s contract-less branch is reworded to say the field is outside the contract and read verbatim, instead of calling the value unrecognized (see D-01).

## Decisions

- **D-01**: Kept `schemaFieldWarningText`'s contract-less branch rather than deleting it, reworded, with a comment stating it is unreachable from this package and why it survives. **DECIDE & STATE** — four lines, no behaviour change on any reachable path.

  Reasoning: the REQ named both options and required that neither leaves dead code under a stale comment. Deleting it is tempting under YAGNI, but the function would then fall through to its `defaultValue == ""` branch on an ungated call and print `expected one of []` — a worse failure than the one being fixed, for the next caller who forgets the gate. Keeping it also satisfies the second half of the user's ask literally: they asked for the warning to be *gated* **and** for the wording to match the contract's outside-the-contract classification, which only makes sense if the text still exists. The comment states the unreachability as fact, so it is a guard, not orphaned code.

## Testing

**Tests run:** `go test ./...` (in `tools/queue-kanban/`), plus `go test -run Frontmatter .`
**Result:** ✓ All passing

**Red-green validation:**
- `TestRunFrontmatterCommandNormalizeIsSilentOnFieldOutsideTheContract` (3 sub-cases: `created_at`, `title`, `id`): ✗ before — each printed `⚠ <field>: '<value>' not recognized — no canonical vocabulary is defined for this field.` → ✓ after, stderr empty
- `TestRunFrontmatterCommandNormalizeMatchesPlainGetOutsideTheContract`: ✗ before (`--normalize` changed observable output) → ✓ after
- `TestRunFrontmatterCommandInSetRejectsFieldOutsideTheContract`: ✗ before (`exit = 1, want 2`) → ✓ after
- Regression guard, green before and after: REQ-115's `TestRunFrontmatterCommandWarnsOnUnrecognizedStatus` and `TestRunFrontmatterCommandInSetWarnsOnUnrecognizedStatus`. A typo'd `status` must still warn — re-opening that hole is the one way this change could do real damage.
- End-to-end with the built binary: all four paths, exit codes and stderr as documented in UNIFY.

**New tests added:**
- `TestRunFrontmatterCommandNormalizeIsSilentOnFieldOutsideTheContract`
- `TestRunFrontmatterCommandNormalizeMatchesPlainGetOutsideTheContract`
- `TestRunFrontmatterCommandInSetRejectsFieldOutsideTheContract`

*Verified by work action*

## Lessons Learned

**What worked:** Treating the noise as a symptom and asking what the code could not express. The warning was not a wording bug — `isKnownSchemaFieldValue` was being asked a question it structurally could not answer, returning one `false` for "this value is wrong" and "this field isn't mine". Once that was named, the fix was a missing predicate rather than a condition bolted onto the call site, and the same predicate is what any future caller needs.

**What didn't:** The instinct to delete the now-unreachable warning branch. Following it would have left `schemaFieldWarningText` falling through to `expected one of []` for an ungated caller — a worse message than the one being removed. Zero-value structs make "unreachable, so delete it" more dangerous in Go than it looks: the fallthrough path still runs, it just runs on empty data.

**Worth knowing:** The gate splits `--normalize` and `--in-set` deliberately. Silence is right for `--normalize` (nothing to normalize against, so no-op) and wrong for `--in-set`, because both set names are `status` sets — answering "not a member" for a timestamp would look like a real negative at a call site written as `if …; then`. One narrower looseness is left alone on purpose: `--in-set` on a field that *has* a row but is not `status` (e.g. `domain`) still exits 1 rather than erroring. Pre-existing, out of this REQ's scope, and harmless today since the only prose call site passes `status`.

## Orientation

`queue-kanban frontmatter get … --normalize` is now quiet on fields the Schema Read Contract doesn't govern — a timestamp or a title prints its value and nothing else, while a typo'd `status` still warns exactly as before. Lives in the frontmatter CLI surface REQ-112 added (`tools/queue-kanban/prime-do-kanban.md`). Not `[MAP CHANGED]` — one gate and one new predicate inside an existing command; no new subcommand, no write surface, no change to the single prose call site. Prime staleness spot-check: `prime-do-kanban.md`'s referenced paths all still exist. Its subcommand list does not name `frontmatter` — pre-existing since REQ-112 and left alone here rather than fixed as a drive-by; noted in the review as a Minor finding.

## Review

**Overall: 96%** | 2026-08-06T11:13:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — all four paths exercised with the built binary; the one prose call site (`actions/commit.md`) passes a field with a contract row and is unaffected.
**Suggested testing:** 2 items
**Follow-up REQs created:** None

**Minor:** `prime-do-kanban.md`'s opening subcommand list (`summary | generate | serve | next-req | next-version | verify | now`) still omits `frontmatter`, which REQ-112 added. Not this REQ's doing and deliberately not fixed here — a drive-by edit to a prime's header is the adjacent-code "improvement" `coding-guardrails.md` § Surgical Changes rules out. Worth a one-line fix in whatever REQ next touches that file.

**Suggested additional testing:** (1) a `--normalize` on a field that is *present but empty* (`title:` with no value) — the absent-field path returns exit 1 before the new gate is reached, which is correct but untested for a row-less field; (2) if a future contract row is added for a field currently row-less, confirm the gate flips it back into the warn path with no code change — that is the property the map lookup buys over an exemption list.

**Restatement sweep:** ran. The diff changes what `--normalize` *does* for a class of fields, so the consumers are the documented callers rather than code: `grep -rn 'frontmatter get' actions/ docs/ CLAUDE.md` → one site, `actions/commit.md`, passing `status --in-set terminal-success` (has a row, unaffected). `grep 'no canonical vocabulary is defined'` across all Markdown → only this REQ's own bug description, which is history and correct to keep. `actions/work-reference.md`'s outside-the-contract paragraph needed no change — the code now agrees with it, which is the point of the REQ.

*Reviewed by review-work action*

---
*Source: `do-work/user-requests/UR-024/input.md` — Finding 3 of the 0.174.15-series feedback triage: "frontmatter get … --normalize on a field with no contract row always warns 'not recognized' — gate the warn branch on the field having a row, and reword to match work-reference.md:214's outside-the-contract classification"*
