# REQ-659 hand-back (frontmatter set and req append-section)

- Branch: `worktree-agent-REQ-659-frontmatter-set-append-section`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-659-frontmatter-set-append-section`
- Base commit: `bd56c4b0`
- Commits: `b16563a9` `[REQ-659] add frontmatter set and req append-section` (one commit)

## Write-boundary deviation (read first)

One file outside the brief's seven paths is changed: `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go`, one token. `TestEveryRemainingUtilityHasOneHandler` pins the exact `corehelpers.Handlers()` count (21). The brief's decided design (D-01: a new top-level `req` command with one `Handlers()` entry) makes it 22, so the test fails without the edit (`handlers=22 want 21`). I added `CommandRequest` to its expected list. This is an intentional cross-REQ test update per general.md → Cross-REQ Test-Break Rules, not a behaviour change. Recorded as D-11. If the integrator rejects it, the only alternative is to drop the top-level `req` command, which reopens D-01.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modified): `CommandRequest = "req"` constant, one `Handlers()` entry, `set` dispatch at the top of `handleFrontmatter`, usage string names both verbs. `get` is unchanged.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go` (new): `handleFrontmatterSet`, `handleRequest` (append-section), the REQ target resolver, and the shared refusal shape.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go` (new): the five Proof tests.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (modified): registry list gains `CommandRequest` (see deviation above).
- `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go` (modified): `CanonicalSectionOrder` and `SectionsAfter`.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified): the twelve literal later-section lists now call `requestmodel.SectionsAfter(...)`.
- `skills/do-work/actions/work.md` (modified): the Living Logs paragraph names the `req append-section` argv once. The stamp sentences at the old lines 201, 288, 309 and 375 name `frontmatter set` and cite **Stamps are append-only**. The Testing append at line 342 names `req append-section`.
- `skills/do-work/actions/work-reference.md` (modified): one sentence in **Stamps are append-only** names the `frontmatter set ... --at now` argv. The Timestamp-rule paragraph is untouched.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Reuse `ResolveTarget`, `ParseDocument`/`SetScalar`/`DocumentBytes`, `CanonicalTimestamp`, `VisibleSections`, `ReplaceBodySpan` and `atomicfile.ReplaceExisting`. Add no encoder. Write the five tests first, then the order list and helper, the two handlers, the advance list swap and the prose. Both lesson families shaped the plan. `lifecycle-section-evidence`: only column-0 visible headings count, which is the same rule as `advanceSections`. `alternate-writer-contract-drift`: the writer and `advance` share one order list.
- [x] **[APPLY]:** Done as planned. One change outside the plan: the registry test token (D-11).
- [x] **[UNIFY]:** `git diff bd56c4b0 --stat`:
  ```
   skills/do-work/actions/work-reference.md           |   2 +-
   skills/do-work/actions/work.md                     |  12 +-
   .../do-work-cli/internal/corehelpers/commands.go   |   8 +-
   .../internal/corehelpers/commands_test.go          |   2 +-
   .../internal/corehelpers/request_writers.go        | 246 ++++++++++++++++++++
   .../internal/corehelpers/request_writers_test.go   | 248 +++++++++++++++++++++
   .../internal/lifecycleadvance/advance_commands.go  |  24 +-
   .../internal/requestmodel/visible_sections.go      |  22 ++
   8 files changed, 542 insertions(+), 22 deletions(-)
  ```
  Checks (from the worktree root unless noted):
  | Check | Exit | Wall |
  |---|---|---|
  | RED `go test -run '^(TestFrontmatterSet\|TestRequestAppendSection\|TestRequestIDResolves)' -v ./internal/corehelpers/` | 1 | 1.0s |
  | GREEN same command | 0 | 1.0s |
  | Advance guard before the edit | 0 | 4.1s |
  | Advance guard after the edit | 0 | 3.1s |
  | `gofmt -l` on corehelpers, requestmodel, lifecycleadvance | 0, no output | <1s |
  | `go vet ./internal/corehelpers/ ./internal/requestmodel/ ./internal/lifecycleadvance/` | 0 | 0.4s (cached) |
  | `go test -count=1 ./internal/corehelpers/ ./internal/requestmodel/` (whole touched packages) | 0 | 9.2s |
  | `go test -count=1 -run Advance ./internal/lifecycleadvance/` | 0 | 21.6s |
  | `bash .../REQ-659-probe.sh` | 0 ("REQ-659 GREEN probe: pass") | 2.0s |
  | `git diff bd56c4b0 --check` | 0 | <1s |
  
  Every changed file was re-read in the diff. The diff holds no debug prints and no scratch files. Fixtures live in `t.TempDir()` or `mktemp -d`, and the smoke temp tree was removed.

## Proof record

- **RED** (constant added first, so the failure is a runtime result and not a compile error): `TestFrontmatterSet*` got `outcome=failure` with `HELPER-USAGE: usage: frontmatter get <file> <field> ...` from `commands.go`. `TestRequestAppendSection*` and the `req` part of `TestRequestIDResolves*` got `no handler registered for "req"`. All 5 tests FAIL.
- **GREEN:** `--- PASS:` for `TestFrontmatterSetStampIsAppendOnly`, `TestFrontmatterSetWritesQuotedScalars`, `TestRequestAppendSectionLandsInCanonicalOrder`, `TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses` and `TestRequestIDResolvesOneWorkingOrQueueFile`. No `--- SKIP`.
- **Advance guard:** `TestAdvanceCommandPhaseMatrix` and `TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates` passed before and after the edit. The advance tests were not edited.
- **Real-binary smoke** (`mktemp -d` tree with `git init`, launcher `skills/do-work/tools/do-work-cli.sh --repo-root <tmp> --format text`):
  - `frontmatter set REQ-701 review_at --at now` printed `2026-10-10T13:31:37Z`, exit 0.
  - The second call printed `FRONTMATTER-STAMP-EXISTS`, exit 1.
  - `req append-section REQ-701 --section Testing --from testing.md` placed Testing before Review, exit 0. A repeat exited 0 with no change.
- **Advance lists replaced (12):** four × `SectionsAfter("Triage")`, two × `SectionsAfter("Plan")`, and one each for Exploration, Scope, Implementation Summary, Qualification, Testing and Review. Each list was checked, by a script assertion, to be the exact tail of the canonical order before it was replaced.
- **Lists left literal:**
  - `:174`, the Route A set check (`Exploration`, `Scope`, `Pre-Flight`). The brief said to leave it, and it is a set, not a tail.
  - `hasSection(sections, "Orientation")` in the Lessons Learned branch. It is a single `hasSection`, not a `hasAnySection` list.

## Decisions

- **D-11 (DECIDE & STATE):** I edited `commands_test.go` (one token), outside the write boundary. The registry count test must name the new command, or the decided D-01 design leaves the package red. Reversible.
- **D-12 (DECIDE & STATE):** New text from `req append-section` uses LF. `requestmodel` keeps the line ending unexported (`lineEnding`), and the brief said not to add an accessor. On a CRLF file, the existing-body comparison sees `\r` and refuses a repeat as a conflict. It never duplicates a section.
- **D-13 (DECIDE & STATE):** Argument-shape errors exit 2 (outcome `failure`, `HELPER-USAGE`), matching `usageResult` across `corehelpers`. That covers a wrong argument count, an unknown option, a missing `--section` or `--from`, and an unreadable `--from` file. Every rule refusal exits 1 (outcome `refused`) with a typed code: `FRONTMATTER-AT-INVALID`, `FRONTMATTER-STAMP-NOT-CANONICAL`, `FRONTMATTER-STAMP-EXISTS`, `FRONTMATTER-FIELD-STRUCTURED`, `FRONTMATTER-SET-INVALID`, `SECTION-NOT-CANONICAL`, `SECTION-CONFLICT`, `SECTION-DUPLICATE`, `SECTION-WRITE-FAILED`, `REQUEST-NOT-ACTIVE`, `REQUEST-UNREADABLE`, `REQUEST-WRITE-FAILED`, plus `ResolveTarget`'s own codes passed through. A failed atomic write also reports `refused` (`REQUEST-WRITE-FAILED`), so there is one refusal shape. `ReplaceExisting` leaves the original bytes in place on error.
- **D-14 (DECIDE & STATE):** A `## <name>` section placed before an anchor is written as `## <name>\n\n<body>\n\n`. One placed at the end of the file gets one blank line before it, with no extra blank line if the file already ends with one. The `--from` body is trimmed of blank lines at both ends before it is written, which is the same trim used for the "same body" comparison. Without that shared trim, a no-op repeat would compare unequal.
- **D-15 (DECIDE & STATE):** After building the new body, `req append-section` re-reads `VisibleSections` and refuses unless exactly one visible column-0 `## <name>` results. An unclosed fence or comment at the end of the file would otherwise hide the new heading, and every retry would append another invisible copy. The check is 6 lines, earned by the `lifecycle-section-evidence` family. Not separately tested, because the brief caps the tests at five.
- **D-16 (DECIDE & STATE):** In the `frontmatter set` file-path form, the path is checked to sit under `do-work/working/` or `do-work/queue/` after `filepath.Clean`, so `..` cannot escape. No `ResolveTarget` identity check is applied to a path, because the brief's D-05 asks only for the location rule. Symlinks are refused by `ReplaceExisting` (regular files only).

## Discovered Tasks

- The **Timestamp rule** paragraph in `actions/work-reference.md` says it is "the only place in `actions/` that spells a command for obtaining" an instant. `frontmatter set ... --at now` (now named in **Stamps are append-only**) also obtains one while writing it. A maintainer may want that sentence to say "obtaining an instant to hold", or to name the writer as a second sanctioned site. impact-minor → report only.
- REQ-658 (finalize auto-manifest) also resolves a REQ id to its working file. `resolveActiveRequest` in `corehelpers/request_writers.go` is unexported, and no shared resolver was exported, as the brief asked. If REQ-658 needs the same working-or-queue rule, it should call `requeststate.ResolveTarget` plus a `TreeSection` check the same way, or a later REQ can export this one. impact-minor → report only.
- `actions/fan-out-reference.md:161` repeats the dispatch stamp rule without naming `frontmatter set` (deliberately not edited, per D-09). impact-minor → report only.

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: only the `[family: lifecycle-section-evidence]` and `[family: alternate-writer-contract-drift]` bullets (grep -F).
- Primes: `prime-shell-commands.md`, `prime-action-files.md`, `prime-releases.md`, `tools/do-work-cli/prime-do-work-cli.md`. Crew: general, backend, coding-guardrails, shared-principles, communication-style, testing.

## Anti-bloat check

Named by the brief: `request_writers.go`, `request_writers_test.go`, `CanonicalSectionOrder`, `SectionsAfter`, `CommandRequest`, `handleRequest`, the `set` dispatch, and the five tests. Everything else added:

- `handleFrontmatterSet`: the `set` handler the brief puts in `request_writers.go`, named for the verb.
- `resolveActiveRequest`: the shared REQ-N or path resolution both commands need (D-05). It is one function so the two commands cannot drift.
- `readRequestDocument`: read plus parse, used by both handlers.
- `replaceRequestFile`: the atomic write plus the "no byte changed, no write" guard, used by both handlers.
- `writerRefusal`: the one refused-outcome shape (D-13), used by every refusal.
- `trimBlankLines`: the D-03 trim, used for both sides of the same-body comparison.
- `requestIDPattern`: package-level regexp so it is compiled once. It is used by the resolver and the `req` argv check.
- Test helpers `runRegisteredWriter`, `readFixture`, `writerOutput`, `assertRefusedUnchanged`, plus `writerFixtureRequest` and `canonicalInstantPattern`. They are shared across the five tests. `runRegisteredWriter` goes through `Handlers()`, so RED was a runtime failure, not a compile error.
- `commands_test.go` token: see D-11.

No flags, options, config, `--dry-run`, or other verbs were added.

## Proposed CHANGELOG entry

**REQ Stamps and Sections Get a Validating Writer**

Agents were stamping `review_at` and `integration_at` and appending `## Testing` with hand-written heredocs, and nothing checked the append-only or section-order rules. Two new commands now write these and refuse what the schema forbids.

- `do-work-cli frontmatter set <file|REQ-N> <field> (<value> | --at now)` writes one scalar on a working or queue REQ with the existing quoting rules. It refuses to overwrite any `*_at` stamp, accepts a stamp only in `YYYY-MM-DDTHH:MM:SSZ` form, and refuses list or nested fields.
- `do-work-cli req append-section REQ-N --section <name> --from <file>` inserts a lifecycle section in canonical order. A repeat with the same body is a no-op (exit 0), and a different body or a second copy is refused.
- The canonical section order is now one list that `advance` also reads, so the writer cannot place a section where `advance` would refuse it.
- `actions/work.md` and `actions/work-reference.md` name the two commands at the stamp and section write sites. Writing by hand stays valid.

## Proposed lesson bullet

For `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`:

`- [family: alternate-writer-contract-drift] [REQ-659: a new writer of an ordered structure must read the classifier's order, not copy it — advance's twelve inline later-section lists became one requestmodel.CanonicalSectionOrder before req append-section was allowed to insert by it, and a registry-count test is part of the writer's seam](../../../../do-work/archive/UR-145/REQ-659-frontmatter-set-and-req-append-section.md#lessons-learned)`

## Integration seams

- `skills/do-work/actions/work.md`: edits are sentence-local in the Living Logs paragraph and at the old lines 201, 288, 309, 342 and 375. No reflow. REQ-658 (near 467) and REQ-689 should merge cleanly unless they edit the same sentences.
- `skills/do-work/actions/work-reference.md`: one sentence appended at the end of the **Stamps are append-only** paragraph only.
- `internal/lifecycleadvance/advance_commands.go`: changes stay inside lines 139-246. REQ-658's finalization hint near 264 is untouched.
- `internal/corehelpers/commands.go`: one constant, one map entry, a three-line `set` dispatch, and the usage string.
- `internal/corehelpers/commands_test.go`: one token (D-11). If any sibling also adds a `corehelpers` command, this line will conflict.
- Test wall times: focused five 1.0s; whole `corehelpers` 8.8s; whole `requestmodel` 0.3s; `lifecycleadvance -run Advance` 21.4s; advance guard 3-4s. One run each, and no budget failures under load.
