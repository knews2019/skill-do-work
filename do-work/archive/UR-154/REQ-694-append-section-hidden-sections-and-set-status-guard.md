---
id: REQ-694
title: 'Addendum: req append-section refuses a body that hides later sections, and frontmatter set refuses status and id'
status: completed
created_at: 2026-10-10T19:19:53Z
user_request: UR-154
addendum_to: REQ-659
domain: backend
prime_files: ["_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-693, REQ-695, REQ-696]
batch: review-followups-ur154
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go", "skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go"]
claimed_at: 2026-10-10T19:22:00Z
route: A
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-10-10T19:23:55Z
  basis:
    - Route A
    - 2-file write set
    - 2 acceptance criteria
builder_handback_at: 2026-10-10T19:28:12Z
dispatch_at: 2026-10-10T19:26:34Z
integration_at: 2026-10-10T19:44:13Z
review_at: 2026-10-10T19:58:51Z
kb_status: pending
commit: d7cddc738955b004eef7288567ceea6d5a2e0fec
heavy_verified_at: 2026-10-10T20:04:08Z
heavy_verified_revision: d7cddc738955b004eef7288567ceea6d5a2e0fec
completed_at: 2026-10-10T20:04:57Z
release_at: 2026-10-10T20:04:57Z
---
# Addendum: req append-section Refuses a Body That Hides Later Sections, and frontmatter set Refuses status and id

## What

Two schema gaps the REQ-659 (frontmatter set and req append-section writers) reviews left as report only. `req append-section` accepts a body that leaves a code fence or HTML comment open, which hides every later section from `advance` (re-review N1). `frontmatter set` writes any value into `status` or `id`, although those fields have owning commands (review F3).

## Prior Implementation

REQ-659 shipped both writers in `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go`, archived at commit `8a33190c`. `handleRequest` inserts the section in `requestmodel.CanonicalSectionOrder` and then re-checks that exactly one visible copy of the new section exists. It does not check that the sections after the insertion point stay visible. `handleFrontmatterSet` guards `*_at` stamps (append-only, canonical form) and structured fields, and nothing else.

## Detailed Requirements

- N1: after the insertion, count the column-0 visible sections (`requestmodel.VisibleSections`) and refuse with `SECTION-WRITE-FAILED`, writing nothing, unless the count grew by exactly one. This replaces or extends the current single-copy re-check (the reviewer notes it also covers the end-of-file case).
- F3: `frontmatter set` refuses the fields `status` and `id`, and the refusal names the owning command (status: the lifecycle transitions such as `advance`, `unblock`, `finalize`; id: never rewritten). Other scalar fields keep today's behaviour.

## Red-Green Proof

**RED prompt/case:** (1) A REQ with a `## Review` section; `req append-section REQ-N --section Testing --from body.md` where `body.md` is a line of three backticks followed by `## x` (or `text <!-- open`). (2) `frontmatter set REQ-N status bogus-status`.
**Why RED now:** (1) exits 0 and the following `## Review` disappears inside the open fence or comment, so `advance` reads the phase as missing. (2) exits 0 and writes `status: bogus-status`.
**GREEN when:** (1) refuses `SECTION-WRITE-FAILED` and the file is byte-identical. (2) refuses, names the owning command, and the file is byte-identical.
**Validation:** Inferred during capture

## Constraints

- One test per failure in `request_writers_test.go`; each must fail with the new guard removed.
- No new allowed-values table for other fields (YAGNI): only `status` and `id` are refused.
- No `depends_on` edges in this batch; no other REQ touches these files.

## Builder Guidance

Certainty: high (both reproduced by the reviewer, re-checked at HEAD `0e8e0ef9`). The reviewer sized N1 at about four lines.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 20802 tokens, bare only (`slugged: partial`); matches the do-work-cli corehelpers path (families `writer-anchored-marker-position`, `lifecycle-section-evidence`).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/archive/UR-145/REQ-659-frontmatter-set-and-req-append-section.md` → `## Review`, `do-work/runs/work-2026-10-10-131527/REQ-659-review.md` (F3), `do-work/runs/work-2026-10-10-131527/REQ-659-rereview.md` (N1).

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
*Source: UR-154 R2 and R3 — "`req append-section` accepts a body that leaves a code fence or comment open, which hides every later section from `advance` (review N1; suggested fix: refuse unless the visible section count grows by exactly one)." and "`frontmatter set` accepts any value for `status` (review F3, impact-rule-change)."*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the one production file and its test file, and both changes are decided: a visible-section count check after the insert in `handleRequest`, and a refusal for `status` and `id` at the top of `handleFrontmatterSet`. The reviewer reproduced both failures and sized the first fix at about four lines.

**Planning:** Not required

No `## Open Questions` section exists, so Step 3.5 records nothing there. Pre-dispatch decisions:

- **D-01 (DECIDE & STATE):** Route A, not B. The location is known (`request_writers.go:161-171` and `:25-41`), so exploration would only re-read what this triage already read. The builder brief carries the file:line pointers.
- **D-02 (DECIDE & STATE):** N1 replaces the single-copy re-check at `request_writers.go:161-171` with the count check (column-0 visible sections after the insert must equal the count before plus one). The REQ allows "replaces or extends"; the count check also covers the end-of-file case the old check guarded (a hidden new heading does not raise the count), so keeping both would be dead code.
- **D-03 (DECIDE & STATE):** F3 adds one new refusal code, `FRONTMATTER-FIELD-OWNED`, checked right after argument parsing in `handleFrontmatterSet` and before the target is resolved or read, so the file is never touched. The evidence text names the owner: for `status`, the lifecycle commands (`advance`, `unblock`, `finalize`) or the hand write at the transition's defining site (clarify's `pending-answers` to `pending` flip is a hand write, `actions/clarify.md:136`); for `id`, that an id is never rewritten. A shipped-prose search found no action that tells an agent to run `frontmatter set` on `status` or `id`, so no caller breaks.
- **D-04 (DECIDE & STATE):** Because D-02 removes the only guard for an unclosed fence or comment at the end of the file (a hidden new heading), the new N1 test carries that case as a third table row beside the open-fence and open-comment bodies. It pins that the replacement keeps the old guarantee; no test pinned it before. Each row fails with the count check removed.
- **D-05 (DECIDE & STATE):** Test names are fixed so the GREEN probe can name them: `TestRequestAppendSectionRefusesBodyThatHidesLaterSections` and `TestFrontmatterSetRefusesLifecycleOwnedFields`, both in `request_writers_test.go`. The builder's first implementation decision is D-06.

<!-- D-XX counter: last used D-13. Next decision: D-14. -->

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go` (modified)

**What was done:** `req append-section` now refuses `SECTION-WRITE-FAILED` before reading the request when the `--from` body leaves a fence or comment open: a `## Sentinel` heading placed after the body must stay visible. That catches a body that would hide later sections mid-file and a body appended at the end that would hide the next section written after it. The base single-copy re-check (the request file's own open fence at the end hides the new heading) is unchanged. `frontmatter set` refuses `status` and `id` with the new code `FRONTMATTER-FIELD-OWNED` before the target is resolved, naming the owner. Branch history: builder `cb507a44` (net-count check), integrator `87ffbb4e` (coordinator ruling: offset comparison, D-11), integrator `571b256c` (review F1 fix: body sentinel, D-13). Tests: `TestRequestAppendSectionRefusesBodyThatHidesLaterSections` (5 rows) and `TestFrontmatterSetRefusesLifecycleOwnedFields` (2 rows).

## Decisions

(from the builder hand-back, verbatim; D-11 to D-13 added by the integrator)

- **D-06 (DECIDE & STATE):** The before-count is `sectionsBefore++` inside the existing loop (after its `HeadingIndent != 0` skip), and the after-count is a short local loop. No closure: an earlier draft used a local `visibleCount` closure, but it re-parsed the body that the loop already walks, so I removed it.
- **D-07 (DECIDE & STATE):** The `FRONTMATTER-FIELD-OWNED` finding names `target` (the argument as typed, for example `REQ-701`) as its path, not the resolved file. The refusal fires before resolution on purpose (D-03), so the file is never read. The existing `FRONTMATTER-AT-INVALID` refusal uses the same shape.
- **D-08 (DECIDE & STATE):** The field match is exact and case-sensitive (`status`, `id`). Frontmatter keys are lowercase by schema. Matching `Status` would be a guess at a failure nobody has seen (YAGNI).
- **D-09 (DECIDE & STATE):** New `SECTION-WRITE-FAILED` evidence text: `the inserted ## <name> would hide itself or a later section; close every fence and comment in the body`. It tells the writer what to fix in both cases (a hidden heading, or hidden later sections).
- **D-10 (DECIDE & STATE):** Test row 3 uses a REQ whose body ends `## What` then an unclosed fence, with no later canonical section, and an ordinary `--from` body (D-04).
- **D-11 (DECIDE & STATE, integrator, coordinator ruling):** Replaced D-02's decided count check (and D-06's `sectionsBefore++` count) with: every column-0 section visible before the write is still visible after it, by name and shifted offset, and the only new visible heading is the appended section. This overrides the brief's decided design because the decided design does not fix the named failure: the builder reproduced a body that opens a fence, hides `## Review`, and reveals a fenced example heading lower in the file, so the count still grows by one and the write succeeds (hand-back Discovered Task 1). Committed on the builder branch as `87ffbb4e`.
- **D-12 (DECIDE & STATE, integrator; superseded by D-13, the row is back to the builder's exact `## Hidden` repro):** The new N1 table row is the builder's repro with the fenced example heading named `## Review` instead of `## Hidden`. With that name the count still grows by one and the section names still match, so the row fails under both the count-only check (verified at `cb507a44`) and a names-only comparison (verified by mutation), and passes only with the offset comparison. One row pins both.
- **D-13 (DECIDE & STATE, integrator, after review F1):** Replaced the D-11 offset-comparison loop with a check on the `--from` body itself: `VisibleSections(body + "\n\n## Sentinel\n")` must return exactly one section, else `SECTION-WRITE-FAILED` naming the `--from` file; the base single-copy re-check stays as it was. Challenge to the coordinator ruling, recorded here because the run is unattended: the ruling's check observes only sections that already exist, so a body appended at the end (the usual Testing-before-Review order) could still leave a fence open and hide the `## Review` written next (review F1, reproduced on the live CLI). The insert point is always a visible section start or the end of a file whose state the re-check covers, so a body that closes everything it opens cannot change how the rest of the file parses; the body check therefore covers every case the ruling's loop covered, and with it in place the loop had no failing row left (it would have been untested code). Net: 7 added lines instead of 20, every row fails under a mutation (sentinel deleted: rows 1, 2, 4, 5; re-check deleted: row 3). The refusal text now says a line starting with three or more of one punctuation character opens a fence (review F2).

## Discovered Tasks

(from the builder hand-back, verbatim; the first item was resolved at integration by D-11 and then D-13)

- The count check does not catch a body whose open fence flips the parity of a fenced example further down the file. Reproduced on `cb507a44` with a temporary test (deleted): REQ body `## What`, `## Review`, then a fenced example holding `## Hidden`; `--from` body ```` ``` ```` + `open`. Result: exit 0. `## Review` is hidden and `## Hidden` becomes a visible section, so the count still grows by one. A fix that does not depend on what comes after the insert point: refuse when the `--from` body leaves a fence or comment open, for example by checking that a sentinel heading appended after the body is visible (`VisibleSections(body + "\n\n## Sentinel\n")`). The design in the brief was marked as decided, so I did not change it. impact-user-visible → report only
- No shipped prose names the new `FRONTMATTER-FIELD-OWNED` code, and no action tells an agent to run `frontmatter set` on `status` or `id` (pre-dispatch checked this). If `work-reference.md` → Stamps are append-only ever lists the writer's refusal codes, add this one there. impact-negligible → report only

## Qualification

**Gate records:** `advance --diff-range 874ae4e6..7bb792bd` (first merge) returned the `qualify` gate `satisfied` (provenance `merged_range`) with no findings. After the review fix the cumulative range is `874ae4e6..d7cddc73`; `advance` was then in the review phase and refused gate input, so the diff hygiene was re-checked by hand: no print, TODO, FIXME or debug line added in the range.

**Requirement trace** (against `git diff 874ae4e6..d7cddc73`, 2 files, +66 -0, and the merged files):
- N1: `handleRequest` refuses `SECTION-WRITE-FAILED` before the request is resolved or read when the `--from` body leaves a fence or comment open (a `## Sentinel` heading after the body is not visible), so nothing is written. Because the insert point is a visible section start (or the end of the file), a body that closes what it opens cannot hide anything after it; this covers the REQ's case (later sections hidden mid-file), the coordinator's case (an open fence that hides `## Review` and reveals a fenced example heading), and review F1 (an append at the end that would hide the next section). The base single-copy re-check still refuses an append into a file that already ends inside an open fence. This is stricter than the REQ's "count grew by exactly one" wording, which D-11 and D-13 explain.
- F3: `handleFrontmatterSet` refuses `status` and `id` with `FRONTMATTER-FIELD-OWNED` right after argument parsing, before resolution, so the file is never read or written. The `status` evidence names `advance`, `unblock`, `finalize`; the `id` evidence says it is never rewritten. Other scalar fields keep the old path (`TestFrontmatterSetWritesQuotedScalars`, `TestFrontmatterSetStampIsAppendOnly` still pass).
- Red-Green Proof: RED case (1) is table rows 1 and 2, case (2) is `TestFrontmatterSetRefusesLifecycleOwnedFields/status`; both assert the refusal code and a byte-identical file (`assertRefusedUnchanged`).

**Extra cases (live CLI, scratch repo):** at the first merge, a valid append before `## Review` in a file that also holds a fenced example `## Review` succeeded, a CRLF request accepted an append, an open-comment body refused with the file byte-identical, and `frontmatter set REQ-701 status pending` refused with the file byte-identical. At `d7cddc73` the live CLI re-check is recorded under `## Testing`.

**Scope:** Route A, `write_set` names `request_writers.go` and `request_writers_test.go`; the range touches exactly those two files. No package-level identifier added (seam with REQ-693, which edits `handoff.go` in the same package).

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` twice: at the first merge `7bb792bd` (load 2.09, wall 121 s, exit 0, recorded through `advance --gate-arg ... --probe-file REQ-694-probe.sh`: `green-gate` and `run-blocked-check` satisfied) and again at `d7cddc73` after the review fix (load 2.98, wall 135 s, exit 0, recorded with `record-green-gate` because `advance` was in the review phase). GREEN probe `do-work/runs/work-2026-10-10-192201/REQ-694-probe.sh` exit 0 at both merges. Focused `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/corehelpers/` exit 0 on the builder branch.
**Result:** ✓ All passing (do-work-cli 915 Go tests, queue-kanban 423, no test file over the 30 s budget).

**Red-green validation:** (traces to `## Red-Green Proof`)
- `TestRequestAppendSectionRefusesBodyThatHidesLaterSections/an_open_fence_in_the_body` and `/an_open_comment_in_the_body` (RED case 1): ✗ before (builder RED run, `outcome=success findings=[]`) → ✓ after.
- `TestFrontmatterSetRefusesLifecycleOwnedFields/status` (RED case 2) and `/id`: ✗ before → ✓ after.
- `/an_open_fence_that_reveals_a_fenced_example_heading` (coordinator ruling, builder's repro): ✗ against the builder's count check at `cb507a44` → ✓ after.
- `/an_open_fence_in_a_body_appended_at_the_end` (review F1): ✗ at `87ffbb4e` (`outcome=success findings=[]`) → ✓ at `571b256c`.
- `/an_open_fence_already_at_the_end_of_the_file`: passes at base and after; pins that the base single-copy re-check still guards the request file's own open fence.
- Mutations at `571b256c` (restored, `cmp` equal): body sentinel deleted → rows 1, 2, 4, 5 fail; single-copy re-check deleted → row 3 fails.
- Live CLI at `d7cddc73` (scratch repo): an end-of-file open-fence body refuses with the file byte-identical; a body with a closed fence appends, then `## Review` appends after it, and an identical repeat is a no-op success.

**New tests added:**
- `TestFrontmatterSetRefusesLifecycleOwnedFields` (2 rows) and `TestRequestAppendSectionRefusesBodyThatHidesLaterSections` (5 rows) in `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go`.

**Heavy verification plan:**
- Range: 874ae4e6403e2c2e3669cd97e48de2ec3f1c1387..d7cddc738955b004eef7288567ceea6d5a2e0fec
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — both changed files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — both changed files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — both changed files matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 97%** | 2026-10-10T19:58:50Z (first review 92% at 7bb792bd; re-review of the fix delta 7bb792bd..d7cddc73 at 97%)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. Both guards work, and the body check covers every hidden-section case at its source. Full report: `do-work/runs/work-2026-10-10-192201/REQ-694-review.md`.

**Important findings:**
- F1 (first review): an append at the end of the file accepted a body that left a fence or comment open, hiding the next section written after it. Fixed on the builder branch (`571b256c`, D-13, test row 5), re-merged as `d7cddc73`, re-gated green, re-reviewed. Resolved.

**Minor findings:**
- F2: the refusal text misled for punctuation-run fences such as `...`. Fixed in `571b256c` (the text names the rule and the `--from` file). Resolved.
- F3: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:102` said "the open-fence case is still open". Fixed by the integrator in the finalization commit. — impact-negligible → report only
- F4: the hand-back's proposed CHANGELOG entry and lesson bullet described the replaced count check. Rewritten by the integrator for the release. — impact-negligible → report only
- N1: a prose line that starts with `...` is refused; correct for `advance`'s shared fence dialect, and any relaxation belongs in `requestmodel.VisibleSections`. — impact-negligible → report only

**Acceptance:** Pass — focused tests green, both mutations reproduced by the reviewer, live CLI confirms the F1 repro refuses with the file byte-identical and legitimate bodies (closed backtick/tilde fences, inline and multi-line comments, tables, CRLF) are accepted.
**Restatement sweep:** `SECTION-WRITE-FAILED` now also fires on the `--from` body; new `FRONTMATTER-FIELD-OWNED`; one stale restatement found (F3). Not the wave's last integration, so no wave-end sweep.
**Anti-bloat:** no new package-level names, helpers, options or files; one new refusal code (F3 of the REQ needs it); the final N1 change is one 3-line check plus a comment.
**Suggested testing:** 0 items
**Follow-ups created:** None (all remaining findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Checking the `--from` body on its own (a heading placed after it must stay visible) catches every hidden-section case at its source, because the insert point is always a visible section start where no fence or comment is open. Mutation runs per guard showed which row each check owns.
**What didn't:** Two whole-file checks after the insert. The net section count (the REQ's own wording) was fooled when an open fence hid `## Review` and revealed a fenced example heading below it. The coordinator's stricter "every earlier section stays visible" comparison fixed that but still passed an append at the end of the file, where nothing follows yet: the open fence then hides the next section written (review F1).
**Worth knowing:** A post-insert check sees only the sections that exist now. Validate the inserted text closes what it opens. `VisibleSections` treats any line starting with three or more of one punctuation character (such as `...`) as a fence, so the writer refuses such a prose line too; relaxing that belongs to the shared parser, not the writer.

## Orientation

Now `req append-section` refuses a section body that would hide later sections from `advance`, and `frontmatter set` refuses `status` and `id`; both live in the do-work-cli request writers (`internal/corehelpers/request_writers.go`, see `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`). Primes checked: both prime files' paths still exist.

## Heavy Verification Plan

- Base: 874ae4e6403e2c2e3669cd97e48de2ec3f1c1387
- Target: d7cddc738955b004eef7288567ceea6d5a2e0fec
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — both changed files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — both changed files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — both changed files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: d7cddc738955b004eef7288567ceea6d5a2e0fec
- Execution revision: d7cddc738955b004eef7288567ceea6d5a2e0fec (detached drain checkout, `QUEUE_KANBAN_BROWSER` set, 1-minute load 4.68 at start)
- do-work-cli-integrations: exit 0, executed (fingerprint_mismatch), 76 s
- staged-skills: exit 0, executed, 41 s
- updater: exit 0, executed, 71 s
- installer: exit 0, executed, 30 s

## Timing

Observed 2026-10-10T19:43:51Z to 2026-10-10T20:04:09Z: 20m 18s total, 19m 51s attributed across 5 events, 27s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 10m 50s | 1 |
| verification-gate | 8m 39s | 3 |
| handback-merge | 22s | 1 |

Slowest stage: handback-merge / hand-back merge, 22s, outcome success.
Slowest command: review / independent review and delta re-review, 10m 50s, exit 0, .

Notes: no builder-work event was recorded. The hand-back had landed before the integrator started (builder commit `cb507a44` at 2026-10-10T19:28:12Z, dispatch 2026-10-10T19:26:34Z), so per fan-out-reference "Landed hand-back" the wait is not measured. `record-timing-event` has no end-time option, so each event ends when it was recorded. The two repository gates measured 121 s and 135 s of wall time, and the heavy drain measured 218 s across its four lanes. The review event includes the integrator's F1 fix, re-merge and re-gate.
