---
id: REQ-093
title: "Confirm: six shipped Go-source sites cite the export-ignored CLAUDE.md, and the suite's guard catches none of them"
status: completed
status_changed_at: 2026-08-04T09:00:52Z
claimed_at: 2026-08-04T09:03:16Z
completed_at: 2026-08-04T09:09:43Z
kb_status: pending
route: B
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
  - CLAUDE.md
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

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified by `## Answer` and the six comment sites are named by
line. But the guard half is a rewrite, not an edit, and needed discovery first: whether the suite
already had an allowlist idiom to copy (it did not — this is the first), and how its checks report
and accumulate failures.

**Planning:** Not required

## Exploration

- `_dev/tests/contract-regressions.sh` has **no** existing allowlist/exemption idiom — greps for
  `allowlist|exempt|whitelist|skip_files` return only the `shipped_citation_paths` array. So the
  per-file allowlist is a new pattern in this suite, designed here rather than copied.
- The existing check is a four-line shape: build a pattern, `grep -rIEn` over
  `shipped_citation_paths`, `printf` a `FAIL:` line to stderr on any hit, `fail_count++`. The
  inverted version keeps that shape and contract (stderr, `fail_count`) and adds a filter loop.
- The suite is bash (arrays already in use), so `<<<` and `${var%%:*}` are available.
- `_dev/` is itself export-ignored, so the suite's *own* comment may cite `CLAUDE.md` freely — that
  is not a site to fix, and the check must not be run against `_dev/`.
- All six Go sites are comments only; none is in a string literal or a test assertion, so no
  behavior or assertion meaning can shift.

## Scope

**Files I will touch:**
- `tools/queue-kanban/verify.go` (modify) — 3 comment sites: 123, 156, 186
- `tools/queue-kanban/verify_test.go` (modify) — 3 comment sites: 89, 141, 171
- `prompts/prompt-kit-step6-constraint-architecture.md` (modify) — drop the unsourceable
  "CLAUDE.md standard" attribution (declared extension, see `## Answer`)
- `_dev/tests/contract-regressions.sh` (modify) — invert the citation check; add the 14-file
  per-file allowlist
- `CLAUDE.md` (modify) — **declared scope drift taken mid-build**, see D-01

**Files I will NOT touch:**
- `actions/memory-reference.md` — REQ-088's authorized wording; its redundancy finding is a
  separate one-line call the user has not made
- The 14 allowlisted files — their mentions are correct and must not be "fixed"
- Any `verify.go`/`verify_test.go` code, only comments

**Acceptance criteria (restated from REQ):**
- [ ] All six citations in `verify.go`/`verify_test.go` state their rule inline; neither file
      mentions `CLAUDE.md`/`AGENTS.md` at all
- [ ] `go build` and `go test ./...` pass in `tools/queue-kanban`
- [ ] The suite's check flags **any** `CLAUDE.md`/`AGENTS.md` mention in a shipped path
- [ ] A per-file (not per-directory) allowlist exempts the 14 legitimate files
- [ ] The check passes on the post-change tree
- [ ] The check **demonstrably fails** on a newly introduced unallowed mention — proven, not assumed
- [ ] `prompt-kit-step6:78` no longer attributes its rule to a "CLAUDE.md standard"
- [ ] Full suite exits 0

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ✓ passing — `bash _dev/tests/contract-regressions.sh` exit 0; `go test ./...` ok
**Dependencies:** ✓ Go toolchain present

*Checked by work action* — the stale `baseline.json` left by a prior session was deleted and
re-recorded for this REQ rather than inherited.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/verify.go` (modified)
- `tools/queue-kanban/verify_test.go` (modified)
- `prompts/prompt-kit-step6-constraint-architecture.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `CLAUDE.md` (modified)

**What was done:** Rewrote six comments in `verify.go` (123, 156, 186) and `verify_test.go`
(89, 141, 171) to state the release invariant they describe instead of citing
`CLAUDE.md § Before Every Commit`; neither file now mentions the maintainer doc. Dropped the
unsourceable "CLAUDE.md standard:" attribution from `prompt-kit-step6-constraint-architecture.md:78`,
keeping the rule. Replaced the suite's idiom-matching `self_citation_pattern` with an inverted
check: it greps every `CLAUDE.md`/`AGENTS.md` mention across the shipped paths and filters out a
14-entry per-file `maintainer_doc_mention_allowlist`, keeping the original check's stderr/`fail_count`
contract and adding a `FAIL:` message that names the allowlist as the third remedy. Updated
`CLAUDE.md:122`, which described the guard as grepping "the common citation idioms" — stale the
moment the guard stopped doing that.

## Qualification

**Passed** — 5 files verified, 8 acceptance criteria traced, P-A-U confirmed against the diff.

- **Files exist / show in diff:** all five appear in `git diff` with the expected hunks.
- **Changes are substantive:** the six Go comments are rewritten prose, not whitespace; the suite
  hunk replaces a 4-line check with a filtered loop plus allowlist; `CLAUDE.md`'s clause is
  rewritten.
- **No behavior shift in Go:** all six sites are comments. `go build` succeeds and
  `go test ./...` reports `ok` — the same result as the pre-change baseline, not merely "passing."
- **Requirements traced:** each of the eight acceptance criteria has a named piece of evidence in
  `## Testing` below, including the negative proof the criteria demanded rather than assumed.
- **Scope:** four of five files were declared at Step 5.5. The fifth (`CLAUDE.md`) is declared
  drift — `## Scope` and `write_set` were both amended at the moment it was taken, and it is
  raised as a review finding below rather than folded in silently.
- **No debug artifacts:** the diff contains no `console.log`, `fmt.Println`, `TODO`, or commented-out
  code. The two temporary probe lines used for negative testing were reverted with
  `git checkout --` and verified gone by a clean `git status`.

## Testing

**Tests run:**
- `bash _dev/tests/contract-regressions.sh` → **exit 0**
- `cd tools/queue-kanban && go build -o queue-kanban . && go test ./...` → **ok** (2.8s)

**Baseline comparison:** pre-change baseline recorded by `tools/checks/preflight.sh` was suite
exit 0 / `go test` ok. Post-change is identical. No regression, no newly-excluded failure.

**Red-green validation — the guard's new behavior, proven in both directions:**

| Probe | Expectation | Result |
|---|---|---|
| Idiom-marked mention (`see CLAUDE.md`) added to `crew-members/general.md` (not allowlisted) | FAIL, naming file:line | ✓ exit 1, printed `crew-members/general.md:43` |
| **Bare-prose mention** (`CLAUDE.md records this as having happened`) added to `actions/board.md` | FAIL | ✓ exit 1 |
| Clean tree, 14 allowlisted files still mentioning the doc | pass, no false positive | ✓ exit 0 |

The middle probe is the decisive one: a bare-prose mention carries **neither** a `§` nor a colon, so
the idiom-widening this REQ originally recommended would have let it through. The inverted check
catches it.

**Scored against the real occurrences (the claim, measured rather than asserted):** with the three
content fixes stashed and REQ-088's original line restored from `976e815`, the new check flags
**8 of 8** mention-lines — `actions/memory-reference.md:88`,
`prompts/prompt-kit-step6-constraint-architecture.md:78`, `verify.go:123/156/186`,
`verify_test.go:89/141/171`. The pattern it replaced flagged **0 of 8**.

**Correction to this REQ's own figures:** the body above says "seven" occurrences and the decision
record cites "0/7". The true total is **8 mention-lines across 4 files** (6 in Go source, 1 in
`prompt-kit-step6`, 1 the already-fixed `memory-reference.md` line). The "0/N" scoring holds at
every N — the old pattern caught none of them — and the user's "4/6" figure for the proposed
widening was specifically about the six Go sites and is correct as stated.

*Verified by work action*

## Decisions

- **D-01**: Took declared scope drift into `CLAUDE.md` to fix a stale restatement this change
  created. `CLAUDE.md:122` described the guard as one that "greps for the common citation idioms" —
  true before this REQ, false after it. **ESCALATE**-class (it edits the repo's canonical rule
  statement, and it is outside the approved file list), so: `## Scope` and `write_set` were amended
  at the moment the drift was taken, and it is raised as a review finding for the user to judge.
  **Value:** the canonical statement of the rule now describes the check that actually runs; left
  alone, the next maintainer reads that idiom-matching is the mechanism and may "restore" it.
  **Risk:** low and fully reversible — one clause in one sentence, no behavior attached. The
  alternative (report-only, fix later) ships a maintainer doc that contradicts its own test suite.
- **D-02**: Did **not** add a staleness check for the allowlist itself — a probe that reports
  allowlist entries which no longer mention the maintainer doc, or that name deleted files. It is
  the same closed-enumeration rot the user's own reasoning identified, one level up. Not built
  because it was not part of the decision and `crew-members/maintenance.md` § 3 asks any addition to
  earn itself against a concrete failing case; there is no such case yet, since the allowlist was
  verified exactly complete today. **DECIDE & STATE**, and recorded in `## Discovered Tasks` so the
  idea survives without being smuggled in.

## Discovered Tasks

- **[low]** The new `maintainer_doc_mention_allowlist` is itself a hand-maintained enumeration and
  can rot in the two directions the inverted check cannot see: an entry whose file no longer
  mentions the maintainer doc (the exemption outlived its reason, so a genuinely bad citation could
  later hide behind it), and an entry naming a file that no longer exists. A cheap probe — for each
  allowlist entry, assert the file exists and still matches the grep — would close both. Deliberately
  not built here (D-02); it is a candidate for the next maintenance pass on the suite.

## Review

**Overall: 95%** | 2026-08-04T09:03:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 97% |
| Scope | 88% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 1 important, 1 minor
**Acceptance:** Pass — all eight acceptance criteria met, each with named evidence; the guard's new
behavior is proven in both directions rather than assumed.
**Suggested testing:** 0 items
**Follow-ups created:** None

**Important finding — declared scope drift into `CLAUDE.md`, for the user to ratify.** This REQ's
approved file list did not include `CLAUDE.md`; the Restatement Sweep found that `CLAUDE.md:122`
described the guard as grepping citation idioms, which this change makes false. It was fixed rather
than only reported, because shipping a maintainer doc that misdescribes its own test suite is the
failure mode this REQ exists to prevent — one level up. `## Scope` and `write_set` were amended when
the drift was taken (D-01). **If you would rather the doc edit had been a separate REQ, it is one
clause and trivially revertible.**

**Minor finding — this REQ's own body overstates the guard's blind spot as "seven" occurrences.**
The real figure is eight mention-lines across four files. The body text and the `## Answer` were
written before the full scoring run and were not retro-edited (the intent trail records what was
believed at decision time); the correction is recorded in `## Testing`. No decision depended on the
difference — the old pattern scored zero at either count.

**Restatement Sweep (run — mandatory):** One stale restatement found and fixed (`CLAUDE.md:122`,
above). Also checked and found **not** stale: `actions/forensics.md` Check 14 and the
`tools/queue-kanban/` co-location rule in `CLAUDE.md` both describe `verify`'s *release* checks,
which this REQ did not touch (only their comments), so neither needed amending. The archived
REQ-088's own "the suite greps for the common citation idioms" sentence was deliberately left
alone — an archived REQ records what was true when it was written, and rewriting it would falsify
the trail.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Scoring a proposed guard against the real occurrences before shipping it, instead
of reasoning about which idioms "should" appear. The measurement (old pattern 0/8, proposed widening
4/6 on the Go sites, inverted check 8/8) settled the design argument in one command, and the
bare-prose negative probe demonstrated the failure the idiom approach would have shipped.

**What didn't:** `git add` with a list containing one stale pathspec — the REQ-090 queue path that
`git mv` had already moved. The whole invocation aborted, so the commit captured only the rename and
left both the cancellation body and REQ-093's answers unstaged. Both halves of this trap are written
down in `do-work/HANDDOWN-UR-015-016.md` (`git mv` stages content at move time; `git add` aborts on
one bad pathspec) and it was still walked into, because the two combine into a single silent failure:
the commit *succeeded*, reporting nothing wrong. `git status --porcelain -uall` after staging and
before committing is the check that catches it — reading the commit's own `--name-status` afterwards
is what caught it here, one step too late.

**Worth knowing:** `_dev/` is export-ignored, so the contract suite may cite `CLAUDE.md` in its own
comments — the inverted check must never be pointed at `_dev/`, or it flags its own explanation.
Relatedly, the per-file allowlist shape was chosen over per-directory for a concrete reason worth
keeping: allowlisting `actions/` would have exempted `actions/memory-reference.md`, the file the
whole thread started from.

## Orientation

The contract suite now catches dangling maintainer-doc references by their *existence* rather than
their phrasing, so a shipped file that mentions `CLAUDE.md` fails until someone classifies it —
and the six comments in the board's `verify` subsystem that shipped such a reference for months
state their rule inline instead. [MAP CHANGED] — this replaces the mechanism of an existing
contract guard, not just its coverage: the enforcement point moved from a pattern list to a file
allowlist, which is where the judgement now lives and where future maintainers must go to record an
exemption. `tools/queue-kanban/prime-do-kanban.md` was checked for staleness; its referenced paths
all still exist and none of its claims concern the citation guard, so it needs no amendment.
