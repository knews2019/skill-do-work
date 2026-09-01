---
id: REQ-312
title: "[impact-rule-change] Resolve same-package citations in the shipped reference contract"
status: completed
created_at: 2026-08-21T03:20:32Z
status_changed_at: 2026-08-21T15:38:59Z
claimed_at: 2026-08-21T16:54:45Z
completed_at: 2026-08-21T17:54:32Z
commit: 99ea028
user_request: UR-055
addendum_to: REQ-299
domain: general
route: C
review_generated: true
impact: impact-rule-change
kb_status: promoted
kb_entry: REQ-312-resolve-same-package-citations-in-the-sh.md
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-08-21T17:20:50Z
  basis:
    - Route C
    - 11-file write set
    - 4 subsystems involved
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
depends_on: []
maintenance: true
write_set:
- _dev/tests/contract-regressions.sh
- _dev/tests/shipped-package-reference-contract.sh
- skills/do-work/actions/cleanup.md
- skills/do-work/actions/work.md
- skills/do-work/crew-members/prompt-injection.md
- skills/do-work/crew-members/ui-design.md
- skills/do-work-knowledge/actions/prompts.md
- skills/do-work-knowledge/crew-members/prompt-injection.md
- skills/do-work-toolbox/actions/install.md
- skills/do-work-toolbox/actions/tidy-repo.md
- skills/do-work-toolbox/crew-members/prompt-injection.md
- skills/do-work-toolbox/crew-members/ui-design.md
---

# Resolve Same-Package Citations in the Shipped Reference Contract

## What

`_dev/tests/shipped-package-reference-contract.sh` is the guard that keeps a shipped
instruction from pointing at a file that is not there once the suite is installed under
`.claude/skills/`. It resolves **cross-package** citations only. A dangling citation to a
file in the *same* package ships silently.

Measured, not inferred. Four dangling citations were planted one at a time in
`skills/do-work/actions/work-reference.md` and the contract run after each:

| Planted citation | Contract verdict |
|---|---|
| `actions/no-such-file.md` | PASS |
| `../docs/no-such-file.md` | PASS |
| `crew-members/no-such-file.md` | PASS |
| `../../do-work-board/actions/no-such-file.md` | **FAIL**, correctly |

Same-package citations are the large majority of citations in the suite, so the guard
covers the smaller half of the class it exists for.

## Context

Found during REQ-299, whose Consumer-Install Constraint required confirming this contract
"actually covers the new text". It does not — it passed on planted breakage inside the exact
paragraph REQ-299 added.

**Nothing is broken today.** A sweep of the installed suite found no live dangling
same-package citation. The seven candidates a naive sweep surfaced are all legitimate: five
are the deliberately non-existent `prompts/init.md` used as the hostile example in
`crew-members/prompt-injection.md` and `do-work-knowledge/actions/prompts.md`, and two are
consumer-project paths (`docs/design/REQ-NNN-wireframe.md`, `docs/worklog.md`) that the
suite does not own. This is a coverage gap, not an outage — which is exactly why it needs a
check rather than a one-time sweep.

## Requirements

- The contract resolves same-package citations at their **installed** destination, the same
  way it already resolves cross-package ones.
- The two false-positive classes above do not break the build: a path the suite does not own
  (a consumer-project path) and a path a file deliberately cites as non-existent.
- Mutation-proven: each of the three PASS rows above becomes a FAIL, and the suite stays
  green on the untouched tree.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** plant `actions/no-such-file.md` in a shipped action file and run
`bash _dev/tests/shipped-package-reference-contract.sh`. It prints
`shipped package reference contract: PASS`.
**GREEN when:** the same plant fails the contract naming the file and line, and the
untouched tree still passes.
**Validation:** Discovered task from REQ-299; apply `actions/work-reference.md` →
**Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-299: the shipped reference
  contract — the guard that stops an installed instruction from pointing at a file that
  isn't there — checks citations that cross package boundaries but never checks citations
  inside a package, which is most of them. I planted three broken same-package pointers in a
  live action file and the guard passed all three. Nothing is broken right now; the risk is
  that the next broken pointer ships unnoticed. Closing it means teaching the guard to
  resolve same-package citations too, and teaching it to ignore two kinds of path that
  legitimately do not exist: paths in the consumer's own project, and the deliberately fake
  `prompts/init.md` the prompt-injection rules cite as an attack example. Should I process
  this as a new task? → Confirmed: Yes, add to queue

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  **Answered 2026-08-21** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the builder's recommendation via
  `do-work clarify`: extend the guard to same-package citations with the REQ's existing
  exclusions for consumer-owned paths and deliberately nonexistent examples. No additional
  scope was requested.

  Why this is yours rather than mine: the fix is small, but the false-positive policy is a
  judgment call about what the suite owns. Getting it wrong turns a useful guard into one
  that fails on legitimate text, and the cheapest wrong answer — skip anything that does not
  resolve — restores the gap it is meant to close.

---

## Triage

**Route: C** - Complex

**Reasoning:** Reproduction confirmed the three blind spots, but the requested false-positive
policy cannot be derived from a missing path: a real dangling citation and an intentionally absent
example have the same syntax. Closing the class without allowlists requires a package-wide resolver,
explicit consumer roots in live prose across three skill packages, and mutation coverage for source
and installed topology. That is an eleven-file, multi-subsystem contract change.

**Planning:** Required

## Plan

1. Capture untouched RED evidence for the three same-package spellings, then add fixtures that pin
   source and installed resolution, containment, cross-package skip non-fallthrough, explicit
   `<project-root>` exclusions, changelog narrowing, and exact-first sentence-period handling.
2. Make `citation_candidate_tokens` the sole extraction path by deleting the obsolete self-only
   backticked-span helpers and fixtures. Extend `citation_messages` with an explicit cross-package-
   first branch and a same-package branch derived from content-directory names present across the
   manifest modules, never a closed hand list.
3. Resolve leadless paths from the owning package root and `../`-led paths from the citing directory;
   require normalized source and installed targets to remain inside their owning mapped roots. If
   the exact target is missing and ends in one period, retry exactly one period-stripped candidate
   in both topologies while preserving the original token in diagnostics.
4. Root consumer-owned and deliberately absent live examples explicitly at `<project-root>` across
   the ten measured instruction files. Exclude only package-root `CHANGELOG.md` from the new
   same-package arm; retain its existing link, cross-package, and mirror checks.
5. Run the focused suite GREEN, replay all three captured mutations against the final checker and
   require nonzero file/line diagnostics, restore each mutation, then run syntax, shell lint,
   exact-scope/diff checks, and the canonical repository gate directly and unpiped.

**Policy decisions:** The resolver never treats “missing” as an exemption; prose declares external
ownership with `<project-root>`. Cross-package classification always terminates its branch, even
when an existing semantic skip returns no error. Same-package grammar is derived from real module
content directories, and the period fallback is exact-first and one-terminal-period-only.

**Plan validation:** All four requirements map to the tasks above; the original eleven-file scope
contained no release, prime, or historical-record edits. The canonical gate later exposed one stale
existing assertion, recorded in D-02 and added to the final scope.

*Generated by Plan agent*

## Exploration

The existing `citation_candidate_tokens` surface is already correct: it sees prose, inline code,
HTML comments, and annotations inside fenced blocks while excluding payload content that lands in
another file. The gap is entirely in `citation_messages`, which only classifies tokens whose first
post-lead segment names a sibling package.

The package-wide grammar needs two same-package forms:

- leadless `actions/...`, `crew-members/...`, and similar paths resolve from the owning package root;
- `../`-led paths resolve from the citing file's directory only when both normalized source and
  installed targets remain inside that package's mapped roots.

Consumer-owned and deliberately absent examples cannot be distinguished from broken citations by
syntax or non-existence. The condition-based solution is to mark them with an explicit
`<project-root>/...` root in live prose. The final scan currently identifies ten likely instruction
files for that clarification. Historical changelog references must remain outside the new
same-package arm, and sentence-final punctuation needs a narrowly verified fallback rather than
being silently treated as part of the filename.

Maintenance opportunity: `comment_backtick_span`, `backticked_span_texts`, and their self-only
fixture are obsolete after `citation_candidate_tokens` became the live extraction path; delete them
instead of extending two extractors.

The three captured mutations were reproduced in an isolated clean clone: each passed with exit 0
before implementation.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` (modify) — align the existing last30days health assertion with
  the explicit project-root spelling required by this REQ
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — package-wide resolver, fixtures,
  mutation contract, changelog narrowing, punctuation handling, and obsolete-helper deletion
- `skills/do-work/actions/cleanup.md` (modify) — explicitly root consumer-owned docs examples
- `skills/do-work/actions/work.md` (modify) — explicitly root consumer design-document examples
- `skills/do-work/crew-members/prompt-injection.md` (modify) — explicitly root hostile/project-local
  prompt examples
- `skills/do-work/crew-members/ui-design.md` (modify) — explicitly root consumer wireframe path
- `skills/do-work-knowledge/actions/prompts.md` (modify) — explicitly root deliberately absent and
  project-local prompt examples
- `skills/do-work-knowledge/crew-members/prompt-injection.md` (modify) — mirror the consumer-root
  prompt examples
- `skills/do-work-toolbox/actions/install.md` (modify) — explicitly root the external last30days
  install path while preserving toolbox-owned script citations
- `skills/do-work-toolbox/actions/tidy-repo.md` (modify) — explicitly root emitted consumer worklog
  paths
- `skills/do-work-toolbox/crew-members/prompt-injection.md` (modify) — explicitly root consumer
  prompt examples while preserving its distinct valid cross-package citation
- `skills/do-work-toolbox/crew-members/ui-design.md` (modify) — mirror the consumer-root wireframe
  path

**Files I will NOT touch:** root or installed changelogs, version files, primes, REQ/UR history, and
the parent request. Historical changelog text stays intact; only the new same-package arm ignores it.

**Acceptance criteria (restated from REQ):**
- [x] Same-package citations resolve in both source and installed topology with containment enforced.
- [x] Consumer-owned and deliberately absent paths remain legitimate through explicit semantic roots,
  never missing-file allowlists.
- [x] Each of the three captured PASS mutations becomes a file-and-line-naming FAIL, and the untouched
  focused suite remains green.
- [x] The direct, unpiped canonical repository gate exits 0.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the action and shell primes and their linked lessons; reproduced all three
  blind spots in an isolated clean clone; declared an eleven-file implementation scope covering the resolver,
  regression fixtures, and only the live prose that needed an explicit consumer root.
- [x] **[APPLY]:** Added manifest-derived same-package classification, paired source/installed
  resolution and containment, exact-first terminal-period handling, and explicit `<project-root>`
  ownership markers. Removed the obsolete secondary backtick extractor and kept cross-package
  semantics isolated from the new branch.
- [x] **[UNIFY]:** Reviewed the final twelve-file diff and `git diff --stat`; preserved the two byte-mirror
  pairs, found no debug or scratch artifacts, passed Bash syntax, ShellCheck 0.11.0, focused contract
  tests, seven targeted semantic mutations, and final replay of all three captured mutations.

## Decisions

<!-- D-XX counter: D-02 -->

### D-01: Rephrase generic prompt-body prose instead of falsely rooting it

**Type:** DECIDE & STATE

The prompt-injection crew's generic references to prompt bodies describe a class of input, not a
filesystem path. Those references now say “prompt-file body” or “shipped prompt-library bodies”; only
the hostile and project-local examples use `<project-root>`. This keeps the new ownership marker
semantically truthful while avoiding a missing-path allowlist.

### D-02: Expand scope for the stale last30days contract pin

**Type:** DECIDE & STATE

The first canonical gate failed because `_dev/tests/contract-regressions.sh` still required the
unrooted phrase this REQ intentionally replaced in `actions/install.md`. The assertion belongs to
the same changed contract surface, so the scope expands by that one test file and the assertion now
requires the explicit `<project-root>/.claude/skills/last30days/scripts/last30days.py` spelling.
This preserves the original runnable-payload guarantee instead of weakening or bypassing it.

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/shipped-package-reference-contract.sh` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/crew-members/prompt-injection.md` (modified)
- `skills/do-work/crew-members/ui-design.md` (modified)
- `skills/do-work-knowledge/actions/prompts.md` (modified)
- `skills/do-work-knowledge/crew-members/prompt-injection.md` (modified)
- `skills/do-work-toolbox/actions/install.md` (modified)
- `skills/do-work-toolbox/actions/tidy-repo.md` (modified)
- `skills/do-work-toolbox/crew-members/prompt-injection.md` (modified)
- `skills/do-work-toolbox/crew-members/ui-design.md` (modified)

**What was done:** The shipped-reference contract now recognizes same-package citations from the
suite's live manifest topology, verifies paired source and installed targets stay inside the owning
module and exist, and reports dangling citations with their original token. Regression fixtures pin
the new classifier, containment, cross-package isolation, changelog exception, and punctuation
behavior; consumer-owned and deliberately absent paths are explicitly rooted in shipped prose.

## Qualification

Passed — 12 implementation files verified, 4 requirements traced, P-A-U confirmed. Mechanical
qualification and Scope/Implementation Summary parity both pass after the recorded D-02 expansion;
the changes are substantive, contain no placeholder data flow, and stay inside the declared write set.

## Testing

**Tests run:** `bash -n _dev/tests/shipped-package-reference-contract.sh`;
`shellcheck _dev/tests/shipped-package-reference-contract.sh`;
`bash _dev/tests/shipped-package-reference-contract.sh`;
`bash _dev/tests/contract-regressions.sh`;
`QUEUE_KANBAN_BROWSER='…/chromium_headless_shell-1212/chrome-headless-shell-mac-arm64/chrome-headless-shell' bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing; the canonical repository gate exited 0 directly and unpiped.

**Red-green validation:**
- Captured `actions/no-such-file.md`: pre-change contract exited 0 → final checker exits 1 and names
  `skills/do-work/actions/work-reference.md:3` plus the original token.
- Captured `../docs/no-such-file.md`: pre-change contract exited 0 → final checker exits 1 and names
  the citing file/line plus the original token.
- Captured `crew-members/no-such-file.md`: pre-change contract exited 0 → final checker exits 1 and
  names the citing file/line plus the original token.
- Seven semantic mutation groups covering leadless resolution, paired containment, installed mapping,
  explicit-root anchoring, changelog narrowing, punctuation fallback, and cross-skip fallthrough all
  made the focused contract fail and were restored before the final green run.

**New tests added:**
- Same-package source/installed resolution and containment fixtures in
  `_dev/tests/shipped-package-reference-contract.sh`.
- Explicit-root, changelog, exact-first punctuation, cross-skip, and live-mutation regression probes.

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh`: REQ-312 D-02 preserves the last30days complete-payload pin
  while requiring its now-explicit consumer project root.

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-21T17:53:37Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `_dev/primes/prime-action-files.md` still describes the checker as enforcing exactly the
  cross-package class even though it now also enforces same-package citations — impact-rule-change
  → appended to `do-work/prose-backlog.md` under the prose-only Fold-First route.

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the same-package gap is closed in both topologies and mutation-proven.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Deriving citation grammar from the manifest's live content directories closed the class without a
  closed enum, while paired mutation probes proved both source and installed behavior.
- Replaying the three original silent passes in a clean clone gave the Finding-Closure Ratchet an
  exact before/after oracle instead of relying on the untouched corpus alone.

**What didn't:**
- The original eleven-file scope missed an existing contract pin for the last30days prose. The
  canonical gate caught the stale expectation; D-02 expanded the scope without weakening the pin.

**Worth knowing:** A newly added immediate content directory in any manifest module intentionally
broadens the same-package citation grammar. Consumer-owned or deliberately absent path examples must
declare their root explicitly; missing targets are never an exemption.

## Orientation

**[MAP CHANGED]** Shipped Markdown reference validation now covers both sibling-package and
same-package citations. The contract lives in the action-file and shell-command verification
subsystems: it classifies from the suite manifest, checks source and installed topology together,
and treats `<project-root>` as the explicit boundary for consumer-owned examples. This matters
because the majority same-package class can no longer ship a dangling pointer silently.
