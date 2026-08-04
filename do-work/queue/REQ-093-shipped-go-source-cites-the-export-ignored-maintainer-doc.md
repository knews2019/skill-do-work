---
id: REQ-093
title: "Confirm: six shipped Go-source sites cite the export-ignored CLAUDE.md, and the suite's guard catches none of them"
status: pending-answers
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

- [ ] Should the six dangling `CLAUDE.md` citations in `tools/queue-kanban/verify.go` and
  `verify_test.go` be fixed — replacing each `CLAUDE.md § Before Every Commit` reference with the
  rule stated inline?
  Recommended: **Yes, add to queue** (this REQ flips to `pending` and a builder makes the edits).
  Value: `tools/queue-kanban/` ships to every consumer install while `CLAUDE.md` does not, so six
  comments currently point readers at a file they do not have. Fixing them matches what the other
  nine shipped statements of these rules already do — state the rule where it is used.
  Risk: very low and fully reversible — comment text only, no Go behavior changes, and the board's
  test suite proves nothing shifted.
  Also: **(a)** No, discard it — a comment pointing at a missing file is harmless enough to leave;
  **(b)** Fix only the three in `verify.go` (shipped production source) and leave `verify_test.go`
  alone, on the grounds that a consumer is far less likely to read the test file.

- [ ] Separately: should `_dev/tests/contract-regressions.sh`'s citation check be widened so this
  idiom class is caught mechanically in future, rather than by someone happening to grep for it?
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
