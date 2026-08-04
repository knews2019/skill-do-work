---
id: REQ-093
title: "Confirm: six shipped Go-source sites cite the export-ignored CLAUDE.md, and the suite's guard catches none of them"
status: pending
status_changed_at: 2026-08-04T09:00:52Z
created_at: 2026-08-04T04:28:03Z
user_request: UR-015
addendum_to: REQ-088
discovered_during: REQ-088
domain: general
maintenance: true
write_set:
  - tools/queue-kanban/verify.go
  - tools/queue-kanban/verify_test.go
  - prompts/prompt-kit-step6-constraint-architecture.md
  - _dev/tests/contract-regressions.sh
---

# Confirm: Six Shipped Go-Source Sites Cite the Export-Ignored CLAUDE.md

## What

REQ-088 fixed one line of `actions/memory-reference.md` that pointed a consumer at this repo's
`CLAUDE.md` — a file that is `export-ignore`d, so it never lands in a consumer install and the
pointer dangles. While verifying that REQ-088's one-site inventory was complete, a sweep of every
shipped path found **six more sites of the same defect**, all in shipped Go source:

| File | Line | The citation |
|---|---|---|
| `tools/queue-kanban/verify.go` | 123 | `covers CLAUDE.md § Before Every Commit items 1 and 2` |
| `tools/queue-kanban/verify.go` | 156 | `(CLAUDE.md § Before Every Commit item 2)` |
| `tools/queue-kanban/verify.go` | 186 | `the duplicate version numbers CLAUDE.md records as having already happened` |
| `tools/queue-kanban/verify_test.go` | 89 | `each direction has its own cause — CLAUDE.md § Before Every Commit` |
| `tools/queue-kanban/verify_test.go` | 141 | `CLAUDE.md records as having already happened` |
| `tools/queue-kanban/verify_test.go` | 171 | `CLAUDE.md § Before Every Commit, item 2` |

`tools/queue-kanban/` is deliberately **not** export-ignored — it ships in the tarball so
`do-work update` carries the board into every consumer install. `CLAUDE.md` is export-ignored. So
these six comments ship to consumers who cannot open what they cite.

A seventh, lower-confidence site: `prompts/prompt-kit-step6-constraint-architecture.md:78`
attributes a rule to a "CLAUDE.md standard" — *"if removing a line wouldn't cause mistakes, cut
it"* — and that sentence is not in this repo's `CLAUDE.md` at all. It may be a generic idiom
meaning "the sort of rule you'd put in a CLAUDE.md" rather than a citation, which is why it is
listed separately rather than counted among the six.

## The guard that was supposed to catch this catches none of it

`_dev/tests/contract-regressions.sh` already has a check for exactly this defect. Its pattern is:

```
self_citation_pattern='(see|per|→) `?CLAUDE\.md|CLAUDE\.md`? *→|(see|per) `?AGENTS\.md'
```

Run verbatim over the shipped paths it scans, it returns **zero hits** — it matched neither the six
sites above nor the line REQ-088 was filed to fix. It only recognizes three idioms (`see CLAUDE.md`,
`per CLAUDE.md`, `CLAUDE.md →`). The two idioms actually in use here — a `§`-section reference and a
`CLAUDE.md: <rule>` parenthetical — sail past it. The guard's own comment says its patterns are
"illustrative, not exhaustive," which is honest but means the check currently provides no coverage
of the two forms that occur in practice.

## Why this is your call, not the builder's

Two reasons, and they are different in kind:

1. **You already declined the guard-widening half once.** REQ-088 offered widening this grep as
   option (c) and your answer put it explicitly out of scope. That decision was made when the
   only known instance was the single line REQ-088 fixed — widening a test pattern to catch one
   already-fixed line is hard to justify. There are now six live instances the pattern misses, so
   the trade-off has changed. Reversing your own out-of-scope call on new evidence is yours to do,
   not a builder's to assume.

2. **Nothing is broken.** As with REQ-088, a dangling citation costs a consumer one confused
   lookup and corrupts nothing. These are code comments — no behavior changes either way. So this
   is worth queuing and not worth a builder sweeping in unasked.

## What Would Change

If both halves are accepted: six comment lines in two Go files get their `CLAUDE.md §` references
replaced with the rule stated inline (the house pattern — the same rule is already restated inline
at nine shipped sites and cites `CLAUDE.md` at none of them), and
`_dev/tests/contract-regressions.sh` gains a broader citation pattern. No Go behavior changes; the
board and `verify` compile and run exactly as now. If only the first half is accepted, the same six
lines change and the guard stays as it is, so the next occurrence is again found by hand.

## Open Questions

- [x] Should the six dangling `CLAUDE.md` citations in `tools/queue-kanban/verify.go` and
  `verify_test.go` be fixed — replacing each `CLAUDE.md § Before Every Commit` reference with the
  rule stated inline? → Yes, fix all six
  Recommended: **Yes, add to queue** (this REQ flips to `pending` and a builder makes the edits).
  Value: `tools/queue-kanban/` ships to every consumer install while `CLAUDE.md` does not, so six
  comments currently point readers at a file they do not have. Fixing them matches what the other
  nine shipped statements of these rules already do — state the rule where it is used.
  Risk: very low and fully reversible — comment text only, no Go behavior changes, and the board's
  test suite proves nothing shifted.
  Also: **(a)** No, discard it — a comment pointing at a missing file is harmless enough to leave;
  **(b)** Fix only the three in `verify.go` (shipped production source) and leave `verify_test.go`
  alone, on the grounds that a consumer is far less likely to read the test file.

- [x] Separately: should `_dev/tests/contract-regressions.sh`'s citation check be widened so this
  idiom class is caught mechanically in future, rather than by someone happening to grep for it?
  → Yes — but **invert the check** rather than adding idioms. See `## Answer` for the full decision;
  the recommendation below (extend the pattern with `§` and `:`) was **rejected** as the wrong shape.
  Recommended: **Yes, widen it as part of the same change.** Concretely, extend the pattern to also
  match a `§`/section reference (`CLAUDE.md §`) and a rule-attribution colon (`CLAUDE.md:`),
  keeping the existing three idioms.
  Value: this is what stops an eighth site appearing. The current pattern scored zero against
  seven real instances, including the one it was written to prevent — so today the check reads as
  coverage without being coverage.
  Risk: a wider pattern can flag legitimate references to a *consumer project's* `CLAUDE.md`,
  which are explicitly allowed. There are several in `actions/prime.md`, `actions/capture.md`,
  `actions/tidy-repo.md` and `actions/bkb-reference.md`, so widening probably needs an allowance
  for those — that is the real cost of this option, and it is why it is a bigger change than the
  first question.
  Also: **(a)** No, leave the guard alone — you have chosen this once already and the six sites
  above can simply be fixed by hand; **(b)** Widen it but only over `tools/` and `prompts/`, the
  two shipped trees where a `CLAUDE.md` mention is least likely to mean a consumer's own file.

## Answer

- **[2026-08-04] Decided by:** user, via the run's batch question review
- **Decision, part 1 — fix all six sites.** Replace each `CLAUDE.md § Before Every Commit`
  reference in `tools/queue-kanban/verify.go` (lines 123, 156, 186) and
  `tools/queue-kanban/verify_test.go` (lines 89, 141, 171) with the rule stated inline. Comment
  text only; no Go behavior changes.
- **Decision, part 2 — widen the guard by INVERTING it, not by adding idioms.** The builder's
  recommendation (extend `self_citation_pattern` with `§` and `:`) was **rejected**. Instead: the
  check flags **any** mention of `CLAUDE.md`/`AGENTS.md` across the shipped paths, with a short
  **per-file allowlist** for the files whose job genuinely is a consumer's own `CLAUDE.md`.
  Both parts land in the **same change**.

### Why the recommendation was rejected (the user's reasoning, recorded)

- **Idiom enumeration is the prohibited shape.** A hand-maintained list of citation idioms is
  exactly the anti-pattern the maintainer doc's "Closed Enumerations Go Stale" section forbids —
  and that section sits roughly 200 lines above this very guard. The rule's own remedy applies:
  state the trigger *condition* and treat any value list as illustrative. Inverting the check makes
  the condition ("a shipped file mentions the maintainer doc") the thing tested, and moves the
  enumeration to a **file** list, which fails loudly when a new file appears rather than silently
  when a new idiom does.
- **The proposed widening measurably underperforms.** Scored against the seven real instances: the
  current pattern catches **0/7**. The proposed `§`-plus-colon widening catches **4/6** of the Go
  sites — it misses `verify.go:186` and `verify_test.go:141`, which are bare prose
  (`CLAUDE.md records as having already happened`) carrying neither marker. It also
  **false-positives** on `actions/prime.md:146`, where the colon is ordinary sentence punctuation
  (`… should be registered in CLAUDE.md:`).
- **The colon idiom was over-generalized from a single instance** — the one line REQ-088 already
  fixed. Building a permanent guard clause around a defect that no longer exists is the weakest
  possible justification for a new pattern.

### Allowlist starting set (per-file, not per-directory)

Fourteen files, surveyed as legitimately referring to a *consumer project's* `CLAUDE.md`/`AGENTS.md`:

`actions/prime.md`, `docs/prime-guide.md`, `actions/version.md`, `actions/tidy-repo.md`,
`actions/bkb.md`, `actions/bkb-reference.md`, `docs/bkb-guide.md`, `actions/capture.md`,
`actions/validate-feedback.md`, `actions/prompts.md`, `README.md`, `prompts/README.md`,
`prompts/prompt-kit-step2-personal-context-doc.md`, `tools/do-work-update.sh`.

**Per-file, deliberately not per-directory** — the point is to keep the residual (a genuinely bad
citation hiding inside an allowlisted file) as small as possible. Allowlisting `actions/` wholesale
would have exempted `actions/memory-reference.md`, the file this whole thread started with.

**Verified against the tree, not taken on trust.** A per-file mention count across the shipped
paths returns **17** files. Fourteen are the allowlist above; the other three are
`tools/queue-kanban/verify.go`, `tools/queue-kanban/verify_test.go` (the six sites being fixed) and
`prompts/prompt-kit-step6-constraint-architecture.md` (see the scope note below). The allowlist is
therefore exactly complete as of this REQ — every shipped mention is either allowlisted or fixed by
this change.

### Scope extension taken at decision time (flagged, not folded in quietly)

`prompts/prompt-kit-step6-constraint-architecture.md:78` is **neither** allowlisted **nor** among
the six sites, so an inverted check would fail the suite the moment it lands. It reads *"Keep it
concise — CLAUDE.md standard: if removing a line wouldn't cause mistakes, cut it."* That sentence
is **not in this repo's `CLAUDE.md` at all**, so it is an unsourceable attribution rather than a
citation of a real rule. Resolution: **drop the attribution**, keeping the rule itself. This
extends the REQ beyond the six sites the user approved, so it is declared here and in `## Scope`,
mirrored into `write_set`, and raised as a review finding for the user to judge rather than
presented as inevitable.

- **Explicitly out of scope:** re-litigating REQ-088's authorized wording at
  `actions/memory-reference.md:88`, and the redundancy review finding recorded against it. That is
  a separate one-line call the user has not made.
