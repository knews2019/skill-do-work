# Hand-back: REQ-694 (req append-section refuses a body that hides later sections; frontmatter set refuses status and id)

- Branch: `worktree-agent-REQ-694-append-section-and-set-guards`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-694-append-section-and-set-guards`
- Base commit: `fd9a2378`
- Commits: `cb507a44` `[REQ-694] refuse hidden-section bodies and owned frontmatter fields` (one commit, two paths, staged by exact path)

## File manifest

- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go` (modified): `handleRequest` counts column-0 visible sections in its existing loop and, after `ReplaceBodySpan`, refuses `SECTION-WRITE-FAILED` unless the count grew by exactly one (replaces the single-copy re-check, comment rewritten). `handleFrontmatterSet` refuses `status` and `id` with `FRONTMATTER-FIELD-OWNED` right after argument parsing; doc comment gains one sentence.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go` (modified): adds `TestFrontmatterSetRefusesLifecycleOwnedFields` (2 cases) and `TestRequestAppendSectionRefusesBodyThatHidesLaterSections` (3 rows).

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read general, backend, coding-guardrails, shared-principles, communication-style, testing; prime-do-work-cli, prime-releases; lessons-releases (whole) and lessons-do-work-cli lines 13, 34, 40, 93. Approach exactly as the brief: N1 = column-0 section count before (in the existing loop at `:121-130`) vs after the insert, must be +1, replacing the single-copy re-check (D-02); F3 = `switch field` on `status`/`id` right after parsing, before stamp checks and before `resolveActiveRequest` (D-03). Tests first (tdd: true). No lesson contradicts the plan: the count keeps the `HeadingIndent == 0` rule from the `rule-direction-checked-against-callers` lesson, and the count is over visible sections only, as `lifecycle-section-evidence` requires.
- [x] **[APPLY]:** Tests written and run RED first, then the two guards. No new package-level identifier (REQ-693 shares the package). One commit.
- [x] **[UNIFY]:** `git diff fd9a2378 --stat`:
  ```
   .../internal/corehelpers/request_writers.go        | 27 +++++++++----
   .../internal/corehelpers/request_writers_test.go   | 45 ++++++++++++++++++++++
   2 files changed, 64 insertions(+), 8 deletions(-)
  ```
  Checks (run from the worktree root; wall times on a loaded machine, build cache warm):
  - `bash .../REQ-694-probe.sh`: exit 0, about 1 s, prints `REQ-694 GREEN probe: pass`.
  - `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^(TestFrontmatterSet|TestRequestAppendSection|TestRequestIDResolves)' ./internal/corehelpers/`: exit 0, about 1 s (`ok ... 0.272s`).
  - `gofmt -l skills/do-work/tools/do-work-cli/internal/corehelpers`: no output, exit 0.
  - `go vet -C skills/do-work/tools/do-work-cli ./internal/corehelpers/`: clean, exit 0.
  - `git diff --check`: clean, exit 0. `git status --short` after commit: clean.
  Files checked: both files above, read in full diff; no debug artifacts. A temporary probe test (`zz_parity_probe_test.go`, see Discovered Tasks) was created and deleted before the hand-back and never staged.

## Proof record

**RED** (tests only, production untouched), exit 1:
```
--- FAIL: TestFrontmatterSetRefusesLifecycleOwnedFields/status   request_writers_test.go:148: outcome=success findings=[], want refused (exit 1)
--- FAIL: TestFrontmatterSetRefusesLifecycleOwnedFields/id       request_writers_test.go:148: outcome=success findings=[], want refused (exit 1)
--- FAIL: TestRequestAppendSectionRefusesBodyThatHidesLaterSections/an_open_fence_in_the_body    request_writers_test.go:254: outcome=success findings=[], want refused (exit 1)
--- FAIL: TestRequestAppendSectionRefusesBodyThatHidesLaterSections/an_open_comment_in_the_body  request_writers_test.go:254: outcome=success findings=[], want refused (exit 1)
--- PASS: TestRequestAppendSectionRefusesBodyThatHidesLaterSections/an_open_fence_already_at_the_end_of_the_file
```
Row 3 passes at base, as expected: the old single-copy re-check refused it.

**GREEN**, same command, exit 0: `--- PASS` for both tests and all 5 subtests; `grep -c -- '--- SKIP'` = 0.

**Mutation 1** (count check deleted, `_ = sectionsAfter` left so it compiles), N1 test exit 1: all three rows fail with `request_writers_test.go:254: outcome=success findings=[], want refused (exit 1)`, including row 3 (the old check is gone too).
**Mutation 2** (F3 `switch` deleted), F3 test exit 1: `status` and `id` both fail with `request_writers_test.go:148: outcome=success findings=[], want refused (exit 1)`.
Restored from a saved copy each time (`cmp` equal) before the commit.

## Decisions

- **D-06 (DECIDE & STATE):** The before-count is `sectionsBefore++` inside the existing loop (after its `HeadingIndent != 0` skip), and the after-count is a short local loop. No closure: an earlier draft used a local `visibleCount` closure, but it re-parsed the body that the loop already walks, so I removed it.
- **D-07 (DECIDE & STATE):** The `FRONTMATTER-FIELD-OWNED` finding names `target` (the argument as typed, for example `REQ-701`) as its path, not the resolved file. The refusal fires before resolution on purpose (D-03), so the file is never read. The existing `FRONTMATTER-AT-INVALID` refusal uses the same shape.
- **D-08 (DECIDE & STATE):** The field match is exact and case-sensitive (`status`, `id`). Frontmatter keys are lowercase by schema. Matching `Status` would be a guess at a failure nobody has seen (YAGNI).
- **D-09 (DECIDE & STATE):** New `SECTION-WRITE-FAILED` evidence text: `the inserted ## <name> would hide itself or a later section; close every fence and comment in the body`. It tells the writer what to fix in both cases (a hidden heading, or hidden later sections).
- **D-10 (DECIDE & STATE):** Test row 3 uses a REQ whose body ends `## What` then an unclosed fence, with no later canonical section, and an ordinary `--from` body (D-04).

## Discovered Tasks

- The count check does not catch a body whose open fence flips the parity of a fenced example further down the file. Reproduced on `cb507a44` with a temporary test (deleted): REQ body `## What`, `## Review`, then a fenced example holding `## Hidden`; `--from` body ```` ``` ```` + `open`. Result: exit 0. `## Review` is hidden and `## Hidden` becomes a visible section, so the count still grows by one. A fix that does not depend on what comes after the insert point: refuse when the `--from` body leaves a fence or comment open, for example by checking that a sentinel heading appended after the body is visible (`VisibleSections(body + "\n\n## Sentinel\n")`). The design in the brief was marked as decided, so I did not change it. impact-user-visible → report only
- No shipped prose names the new `FRONTMATTER-FIELD-OWNED` code, and no action tells an agent to run `frontmatter set` on `status` or `id` (pre-dispatch checked this). If `work-reference.md` → Stamps are append-only ever lists the writer's refusal codes, add this one there. impact-negligible → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file): families `canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`. None applies to a code-only change; the integrator owns the release.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` lines 13 and 34 (`rule-direction-checked-against-callers`): the count keeps the column-0 rule, and no shipped caller relies on `frontmatter set status|id` (checked at pre-dispatch). Lines 40 and 93 (`lifecycle-section-evidence`): hidden text (comment) and quoted text (fence) are both covered by the visible-section count. The Discovered Task above is a new form of the same family.
- Missing files: none.

## Anti-bloat check

`git diff fd9a2378 --stat` is pasted under [UNIFY] (2 files, +64 -8).
Added beyond what the REQ named:
- The constant string `FRONTMATTER-FIELD-OWNED` (named by D-03; the only new code).
- The locals `sectionsBefore` and `sectionsAfter` (they replace the removed `visibleCopies`).
Nothing else: no package-level function, flag, option, config, file or extra test. The two tests are the ones D-05 names.

## Proposed CHANGELOG entry (the integrator writes it with the version)

**Request Writers Refuse Hidden Sections and Lifecycle-Owned Fields**

`req append-section` no longer accepts a section body that hides the sections after it from `advance`. `frontmatter set` no longer edits the two fields that only the lifecycle owns.

- `req append-section` refuses `SECTION-WRITE-FAILED`, and writes nothing, unless the number of visible column-0 sections grows by exactly one. A body that leaves a code fence or HTML comment open could hide the later `## Review` from `advance`. The same check still refuses an append into a file that already ends inside an open fence.
- `frontmatter set` refuses `status` and `id` with the new code `FRONTMATTER-FIELD-OWNED`. The refusal names the owner: the lifecycle commands (`advance`, `unblock`, `finalize`) for `status`, and "never rewritten" for `id`.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`
```
- [family: lifecycle-section-evidence] [REQ-694: a "section count grew by exactly one" check after an insert proves the net count, not which sections stay visible — a body that opens a fence can hide a later real section and reveal a fenced example heading below it at the same time; check the inserted body closes every fence and comment itself instead of judging the whole file after the insert](../../../../do-work/archive/UR-154/REQ-694-append-section-hidden-sections-and-set-status-guard.md)
```

## Integration seams

- Package-level identifiers added: none (REQ-693 edits `handoff.go` in the same package). Only `request_writers.go` and `request_writers_test.go` changed; no other builder touches them.
- REQ-695 (`internal/runstatus/`) and REQ-696 (prose and `_dev/tests/staged-skills-contract.sh`): no overlap.
- Test wall times: RED 1 s, GREEN 1 s, each mutation run about 1 s, probe about 1 s, writer suite 1 s. No timing budget failed, so nothing was rerun.
