---
id: REQ-659
title: 'frontmatter set and req append-section write REQ stamps and sections with the schema rules enforced'
status: completed
created_at: 2026-10-09T21:14:04Z
user_request: UR-145
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-658, REQ-660, REQ-661]
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go, skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go, skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go, skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md]
batch: cli-ergonomics
claimed_at: 2026-10-10T12:52:19Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-10-10T13:19:19Z
  basis:
    - Route B
    - 7-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
builder_handback_at: 2026-10-10T13:31:55Z
integration_at: 2026-10-10T14:42:07Z
review_at: 2026-10-10T14:54:26Z
remediation_at: 2026-10-10T14:54:27Z
re_review_at: 2026-10-10T15:02:18Z
kb_status: pending
commit: 8a33190caa8c3c4f6b41326b46f666ff668a5839
heavy_verified_at: 2026-10-10T15:02:57Z
heavy_verified_revision: 8a33190caa8c3c4f6b41326b46f666ff668a5839
completed_at: 2026-10-10T15:03:36Z
release_at: 2026-10-10T15:03:36Z
---
# frontmatter set and req append-section Write REQ Stamps and Sections With the Schema Rules Enforced
## What
`do-work-cli frontmatter set <file|REQ-N> <field> (<value> | --at now)` sits next to the existing `frontmatter get`. It refuses to overwrite an existing `*_at` field, formats the instant per the Timestamp rule, and keeps the Frontmatter Quoting contract. `do-work-cli req append-section REQ-N --section <name> --from <file>` inserts a REQ body section in canonical order, refuses a second copy, and is a no-op with exit 0 when the same bytes are already there.
## Why
Report item 2 (UR-145 input, "What happened"): 112 python heredocs that stamp `review_at:` or `integration_at:` into a REQ file, in about 24 sessions over 13 days. One example (2026-10-02, this checkout) asserted an anchor and ran `s.replace("\nintegration_at:","\nreview_at: 2026-10-02T18:47:04Z\nintegration_at:",1)`. Separately, `date -u +%Y-%m-%dT...` was used 356 times against 40 uses of `do-work-cli now`. The append-only and timestamp rules are ones a CLI can enforce and a heredoc cannot.
## Verified Facts (checked at 0.305.87)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:77-79`: `handleFrontmatter` accepts only `get` ("usage: frontmatter get <file> <field> [--normalize] [--in-set SET]").
- `skills/do-work/actions/work.md:309` (`integration_at: <now>`, "only if the field is absent") and `:375` (`review_at`, `remediation_at`, `re_review_at`, same absent-only rule). The report cites `:317` and `:383` at 0.305.84.
- `skills/do-work/actions/work-reference.md:59` (Timestamp rule: UTC `YYYY-MM-DDTHH:MM:SSZ`) and `:73` ("**Stamps are append-only.**", with four documented current-state exceptions: `status_changed_at`, `completed_at` on `failed` → `cancelled`, `blocked_at`, and `heavy_verified_at` / `heavy_verified_revision`).
- `skills/do-work/tools/replace-text-section.sh:10`: `replace-section` works on marker-delimited managed blocks, not REQ `## Testing` or `## Review` sections.
- `skills/do-work/actions/sample-archived-req.md`: the de facto section order (What, Red-Green Proof, AI Execution State, Triage, Plan, Exploration, Scope, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation). No command checks it. `advance` refuses an out-of-order write later ("later lifecycle evidence exists before ...", `tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:328`).
- The archived request that gave frontmatter its read surface (`do-work/archive/UR-021/REQ-112-frontmatter-read-cli-surface.md:51`) put a `set` verb out of scope. That rule counted the write surfaces of the queue-kanban tool (`_dev/primes/prime-kanban-board.md:14`, `skills/do-work-board/tools/queue-kanban/frontmatter_cli.go:29`). The report asks to revisit it because the write happens anyway, by hand, with no validation.
## Detailed Requirements
1. `frontmatter set <file|REQ-N> <field> <value>`, or `--at now` in place of `<value>` for `*_at` fields.
2. `set` refuses to overwrite an existing `*_at` field (append-only rule) and leaves the file bytes identical on refusal.
3. `set` formats an `--at now` instant per the Timestamp rule and writes every value per the Frontmatter Quoting contract (`actions/work-reference.md` → Request File Schema).
4. `req append-section REQ-N --section <name> --from <file>` inserts the section in canonical order, refuses a different body for a section that already exists, and exits 0 with no change when the same bytes are already there.
5. Both commands resolve `REQ-N` to its one file in `do-work/working/` or `do-work/queue/` and refuse on zero or several matches.
6. Change the action lines that describe a hand stamp or hand section write (`actions/work.md:309`, `:375` and any other stamp or section write site the builder finds) to name the commands, keeping the hand path valid.
7. Release per `_dev/primes/prime-releases.md`.
## Constraints
- No new frontmatter field and no new status.
- The queue-kanban tool's `frontmatter` command stays read-only; its write-surface count (`_dev/primes/prime-kanban-board.md:14`) is unchanged.
- This REQ is its own release; do not fold REQ-658, REQ-660 or REQ-661 into it.
## Assumptions (recorded at capture, no questions asked)
- The new write surface lives in do-work-cli (`corehelpers`), not in queue-kanban. REQ-112's out-of-scope line was about queue-kanban's counted write surfaces, so this REQ does not reopen that count. By capturing this item the maintainer accepts the report's request to revisit REQ-112's "no `set`" decision for do-work-cli.
- `set` refuses to overwrite every existing `*_at` field, including the four documented current-state exceptions. Those exceptions keep their owning commands (`advance`, `unblock` and the others); `set` is only for append-only stamps. The builder records in Decisions if a caller needs otherwise.
- `--at` accepts only `now`, and only on `*_at` fields. A non-`*_at` field given `--at` refuses.
- For a field that is not `*_at`, `set` adds it when absent and replaces it when present.
- A new field is appended at the end of the frontmatter block. If the schema already defines a field order the CLI uses elsewhere, use that order instead.
- The canonical section order is the `sample-archived-req.md` order. A section name outside that order refuses rather than being guessed into place. If the builder finds the order already encoded in Go (for example in `advance`), reuse that one list and do not add a second.
- Archived REQs are not targets. `review-work`'s archived `## Review` append is out of scope here.
## Dependencies
None. Independent of REQ-658, REQ-660 and REQ-661.
## Builder Guidance
Certainty is high on the rules being enforced (they already exist in prose). Latitude: whether `req append-section` is a new top-level `req` command or a sibling of `frontmatter`, the help text, and how a REQ id resolves to a file (reuse any existing resolver).
## Red-Green Proof
**RED prompt/case:** On a fixture REQ with no `review_at`, run `frontmatter set REQ-N review_at --at now`, then run it again.
**Why RED now:** `commands.go:77-79` rejects any verb other than `get`.
**GREEN when:** The first call adds `review_at` once in Timestamp-rule form. The second call exits non-zero and the file bytes are identical. `req append-section` run twice with the same file leaves one section. A different file for an existing section refuses. A `## Testing` appended to a REQ that already has `## Review` lands before `## Review`.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (19046 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers alternate stored-format writers and lifecycle section evidence; families `alternate-writer-contract-drift` and `lifecycle-section-evidence` fit a second writer of REQ frontmatter and body sections.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over budget; `slugged: partial`). Matching reason: family `alternate-writer-contract-drift` and its index row on alternate artifact writers.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: `prime_files` lists `prime-shell-commands.md` and the REQ changes prescribed command lines in `actions/work.md`.
## Full Context
See `do-work/user-requests/UR-145/input.md` for complete verbatim input (sections Request item 2, What happened, Where the behaviour lives today Item 2, Proposed direction 2, Acceptance check). No queued candidate shares this root cause (the queue held only REQ-654 to REQ-657, the ai-report modes, at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Reuse `ResolveTarget`, `ParseDocument`/`SetScalar`/`DocumentBytes`, `CanonicalTimestamp`, `VisibleSections`, `ReplaceBodySpan` and `atomicfile.ReplaceExisting`; add no encoder. Write the five tests first, then the order list and helper, the two handlers, the advance list swap and the prose. Lesson families `lifecycle-section-evidence` (only column-0 visible headings count, the same rule as `advanceSections`) and `alternate-writer-contract-drift` (the writer and `advance` share one order list) shaped the plan.
- [x] **[APPLY]:** (from the builder hand-back) Done as planned. One change outside the plan: the registry-count test token in `internal/corehelpers/commands_test.go` (D-11).
- [x] **[UNIFY]:** (from the builder hand-back) `git diff bd56c4b0 --stat`: 8 files, 542 insertions, 22 deletions. RED run of the five tests exit 1, GREEN exit 0 (1.0 s). Advance guard tests (`TestAdvanceCommandPhaseMatrix`, `TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates`) exit 0 before and after the edit. `gofmt -l` on corehelpers, requestmodel, lifecycleadvance: no output. `go vet` on the three packages: exit 0. `go test -count=1 ./internal/corehelpers/ ./internal/requestmodel/` exit 0 (9.2 s); `go test -count=1 -run Advance ./internal/lifecycleadvance/` exit 0 (21.6 s); `REQ-659-probe.sh` exit 0 (2.0 s); `git diff bd56c4b0 --check` exit 0. Every changed file re-read in the diff; no debug prints, no scratch files.
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 2: "`frontmatter set <file|REQ-N> <field> (<value> | --at now)` next to the existing `frontmatter get`, plus `req append-section REQ-N --section <name> --from <file>`."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (two verbs, their refusals, and the action lines to change), but where it goes needs discovery: the `corehelpers` dispatch, any existing REQ-id resolver, any section-order list already in Go, and the frontmatter writer that applies the Quoting contract. Exploration settles those; no architectural choice is open.

**Planning:** Not required

### Open Questions

No `## Open Questions` section exists, so Step 3.5 has nothing to resolve. The capture recorded its ambiguities as Assumptions; exploration decisions are recorded under `## Exploration` as D-01 onward.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Required-Lessons Consult

`do-work/lessons-index.md` consulted against the request text, the Scope below and `prime_files`. Added `_dev/primes/lessons-releases.md` (666 tokens, matches `prime_files` → `prime-releases.md`; read in full: its two families are about history links and manifest ownership, so the builder needs it only for the proposed changelog entry). The three over-budget satellites are listed under `## Required Lessons — Dropped for Budget`. The two `lessons-do-work-cli.md` families that matter were read by targeted grep: `lifecycle-section-evidence` (a section reader must ignore fenced and commented headings and count column-0 headings only) and `alternate-writer-contract-drift` (a second writer of a stored format must share the canonical writer's bounds, not re-derive them). Both shape D-02, D-04 and D-06 below.

## Exploration

Explored directly by the pre-dispatch agent at main `bd56c4b0` (no separate Explore agent: the code surface is three Go packages and two action files).

**Key files and what already exists**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:50-63`: `CommandFrontmatter` and the `Handlers()` registry; the CLI help lists commands from these maps, so a new command needs only a registry entry. `:77-146` `handleFrontmatter` accepts only `get`; `:148` `frontmatterUsageResult`. `corehelpers` already imports `requestmodel`, `requeststate` and `repositorymodel`. It must not import `lifecycleadvance` (that package imports `corehelpers`).
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go:120` `ParseDocument`, `:225` `DocumentBytes`, `:230` `BodyBytes`, `:252` `FieldValue`, `:367` `SetScalar` (appends a new field just before the closing fence, which matches the capture assumption), `:420` `CanonicalTimestamp` (the Timestamp rule), `:664` `encodeScalar`. `encodeScalar` already implements the Frontmatter Quoting contract: plain scalar when safe, single-quoted with doubled apostrophes otherwise, and a literal block with `|-` / `|` / `|+` chomping for text with an LF. `doctor/doctor_timestamps.go:180` is an existing caller of `SetScalar`.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go:25` `VisibleSections`: column-0 vs indented headings (`HeadingIndent`), fenced and commented regions hidden. `lifecycleadvance/advance_commands.go:346` `advanceSections` shows the house pattern: skip `HeadingIndent != 0`, refuse duplicates.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_targets.go:11` `ResolveTarget`: the existing REQ-id resolver over `snapshot.RequestsByID` (all trees), refusing zero matches (`REQUEST-NOT-FOUND`), several (`REQUEST-AMBIGUOUS`) and identity mismatch. `repositorymodel.DiscoverRepository` builds the snapshot.
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go:30` `ReplaceExisting`: the atomic writer `doctor` uses for REQ files.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:139-246`: the canonical section order is encoded only as twelve repeated literal "later sections" lists (Plan, Exploration, Scope, Pre-Flight, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation and their tails). There is no single list. `:174` (Route A must not contain Exploration, Scope, Pre-Flight) is a set check, not an order.
- Action prose: `skills/do-work/actions/work.md:21` (Request Files as Living Logs), `:201` (`planning_at`), `:288` (`dispatch_at`, `builder_handback_at`), `:309` (`integration_at`), `:342` (Testing append), `:375` (`review_at`, `remediation_at`, `re_review_at`). `skills/do-work/actions/work-reference.md:59-75`: the Timestamp rule paragraph says it is the only place in `actions/` that spells a command for obtaining an instant; `:73` is **Stamps are append-only**. `skills/do-work/actions/fan-out-reference.md:161` repeats the dispatch stamp rule.

**Decisions (recorded, no question asked)**
- **D-01**: `req append-section` is a new top-level `req` command in `corehelpers` with the one verb `append-section`, spelled as the REQ spells it. Reasoning: a `frontmatter` sibling verb would misname a body write. No other `req` verbs.
- **D-02**: The canonical section order becomes ONE exported ordered list in `requestmodel` (beside `VisibleSections`): What, Red-Green Proof, AI Execution State (P-A-U Loop), Triage, Plan, Exploration, Scope, Pre-Flight, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation, plus a small helper that returns the names after a given one. `advance`'s twelve literal tails are rewritten to read that helper, so the writer and the classifier cannot drift. Pre-Flight is included because `advance` orders it. Reasoning: the capture asked to reuse the order already in Go and not add a second; today's encoding is twelve inline copies, and a writer with its own list would be exactly the alternate-writer drift the lessons describe. Value: one list, and `append-section` can never place a section where `advance` refuses it. Risk: touches `advance_commands.go`; guarded by the unchanged `TestAdvanceCommandPhaseMatrix` and `TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates`. Reversible.
- **D-03**: The `--from` file holds the section body only; the command writes the `## <name>` heading. "Same bytes" means the existing section body equals the file body after trimming leading and trailing blank lines. A different body refuses and leaves the file bytes identical. Reasoning: the `--section` flag is then the only source of the heading, so there is no heading-mismatch case to handle.
- **D-04**: Insertion point: just before the heading of the first column-0 visible section whose canonical index is greater than the target's; with none, at the end of the file after one blank line. Sections outside the canonical list (Why, Constraints, Timing, Decisions and others) are never anchors. A target name outside the list refuses. A target already present twice refuses (the file is malformed, as `advance` says).
- **D-05**: `REQ-N` resolution reuses `requeststate.ResolveTarget`, then requires the tree section `working` or `queue`; an archived match refuses (archived REQs are not targets). A file path argument to `frontmatter set` is accepted only when it lies under `do-work/working/` or `do-work/queue/`. No zero-padding or case grammar is added beyond what `ResolveTarget` already does.
- **D-06**: `frontmatter set` writes through `ParseDocument` → `SetScalar` → `DocumentBytes` → `atomicfile.ReplaceExisting`. No new encoder, no new quoting code. A refusal never writes.
- **D-07**: A field whose name ends in `_at` refuses when it already holds a non-empty value, including the four documented current-state exceptions (capture assumption). `--at now` takes `requestmodel.CanonicalTimestamp(time.Now())`. An explicit value on an `_at` field is accepted only when it is already in canonical `YYYY-MM-DDTHH:MM:SSZ` form, because `work.md:288` stamps `dispatch_at` with an instant held earlier. `--at` with any word other than `now`, or on a field that does not end in `_at`, refuses.
- **D-08**: `frontmatter set` refuses to replace a field whose current value is a list or a nested block (for example `estimate:` or `depends_on:`), because `SetScalar` would collapse it to one scalar. Scalars only.
- **D-09**: Prose: the full launcher argv for each command is written once. The `frontmatter set` argv goes in `work-reference.md`'s **Stamps are append-only** paragraph, and the `req append-section` argv goes in `work.md` → Request Files as Living Logs. The stamp lines `work.md:201`, `:288`, `:309`, `:375` and the Testing append `:342` name the command and cite that paragraph, and each keeps the hand write valid. Reasoning: the Timestamp rule paragraph insists it is the only place that spells a command for an instant; a single argv site keeps that design. `fan-out-reference.md:161` is not edited: it cites the same rule, and REQ-660 and REQ-689 own edits in that file.

- **D-10**: On success, `frontmatter set` prints the value it wrote and nothing else in `--format text` (so a caller can hold the instant `--at now` chose, the way `now` prints one). `req append-section` prints nothing in text mode; both report a refusal as outcome `refused` (exit 1) with a typed finding, and the no-op as outcome `success` (exit 0) with no change.

<!-- D-XX counter: last used D-18. Next decision: D-19. -->

**Concerns**
- REQ-658 (finalize auto-manifest) also resolves "REQ-N to its one working file". Both should reuse `ResolveTarget`; neither should add a second resolver.
- REQ-658 may touch `advance_commands.go` near `:264` (the finalization hint). D-02 edits `:139-246`, so the hunks are separate.
- The run guide's table lists `skills/do-work-board/tools/queue-kanban/frontmatter_cli.go` for this REQ. That is wrong: the Constraints keep queue-kanban's `frontmatter` read-only, so the builder must not touch it.

*Generated by the pre-dispatch agent (direct exploration)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modify) — the set verb in the frontmatter dispatch, the new req command constant and registry entry
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go` (new) — frontmatter set and req append-section handlers, REQ target resolution
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go` (new) — the five focused tests named in the builder brief
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (modify) — the registry-count test names the new req command (added at integration: builder D-11, accepted as D-17)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go` (modify) — the one canonical section-order list and its later-sections helper
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modify) — later-section lists read from the shared order
- `skills/do-work/actions/work.md` (modify) — stamp and section write lines name the commands
- `skills/do-work/actions/work-reference.md` (modify) — Stamps are append-only names the enforcing writer

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/frontmatter_cli.go` and the rest of queue-kanban (its frontmatter command stays read-only and its write-surface count is unchanged), `skills/do-work/actions/fan-out-reference.md` (D-09), `skills/do-work/tools/replace-text-section.sh`, `skills/do-work/actions/sample-archived-req.md`, the release files (changelog, version mirrors: the integrator's).

**Acceptance criteria (restated from REQ):**
- [ ] `frontmatter set <file|REQ-N> <field> <value>` works, and `--at now` replaces the value for fields ending in `_at`.
- [ ] `set` refuses to overwrite an existing `_at` field and leaves the file bytes identical on refusal.
- [ ] `set` formats an `--at now` instant per the Timestamp rule and writes every value per the Frontmatter Quoting contract.
- [ ] `req append-section REQ-N --section <name> --from <file>` inserts in canonical order, refuses a different body for an existing section, and exits 0 with no change when the same bytes are already there.
- [ ] Both commands resolve `REQ-N` to its one file in `do-work/working/` or `do-work/queue/` and refuse on zero or several matches.
- [ ] The action lines that describe a hand stamp or hand section write (`actions/work.md:309`, `:375` and the others found) name the commands, and the hand path stays valid.
- [ ] Released per `_dev/primes/prime-releases.md` (integrator).

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)

**What was done:** Added `do-work-cli frontmatter set <file|REQ-N> <field> (<value> | --at now)` beside `frontmatter get`, and a new top-level `do-work-cli req append-section REQ-N --section <name> --from <file>`. `set` writes one scalar on a working or queue REQ through the existing parse, `SetScalar` and atomic-write path, refuses to overwrite a non-empty `*_at` field, accepts a stamp only in canonical form, and refuses list or nested fields. `append-section` inserts a lifecycle section in canonical order, is a no-op (exit 0) on the same body, and refuses a different body or a second copy. The canonical section order is now one list, `requestmodel.CanonicalSectionOrder`, and `advance` reads its twelve later-section lists from it through `SectionsAfter`. The registry-count test in `commands_test.go` names the new `req` command (D-11). `work.md` names `req append-section` in Request Files as Living Logs and at the Testing append, and names `frontmatter set` at the stamp lines; `work-reference.md` → **Stamps are append-only** holds the one `frontmatter set ... --at now` argv. Hand writes stay valid. After review (remediation re-merge `8a33190c`, builder-branch commit `fd8c934d` by the integrator): `req append-section` refuses a `--from` body that holds its own `## ` heading (`SECTION-BODY-HAS-HEADING`), one assertion pins it, and two `work-reference.md` sentences were corrected (the Timestamp rule is the only command site for an instant *to hold*; the four exceptions are kept by their owning command or the hand write at their defining site).

## Decisions

(from the builder hand-back, verbatim)

- **D-11 (DECIDE & STATE):** I edited `commands_test.go` (one token), outside the write boundary. The registry count test must name the new command, or the decided D-01 design leaves the package red. Reversible.
- **D-12 (DECIDE & STATE):** New text from `req append-section` uses LF. `requestmodel` keeps the line ending unexported (`lineEnding`), and the brief said not to add an accessor. On a CRLF file, the existing-body comparison sees `\r` and refuses a repeat as a conflict. It never duplicates a section.
- **D-13 (DECIDE & STATE):** Argument-shape errors exit 2 (outcome `failure`, `HELPER-USAGE`), matching `usageResult` across `corehelpers`. That covers a wrong argument count, an unknown option, a missing `--section` or `--from`, and an unreadable `--from` file. Every rule refusal exits 1 (outcome `refused`) with a typed code: `FRONTMATTER-AT-INVALID`, `FRONTMATTER-STAMP-NOT-CANONICAL`, `FRONTMATTER-STAMP-EXISTS`, `FRONTMATTER-FIELD-STRUCTURED`, `FRONTMATTER-SET-INVALID`, `SECTION-NOT-CANONICAL`, `SECTION-CONFLICT`, `SECTION-DUPLICATE`, `SECTION-WRITE-FAILED`, `REQUEST-NOT-ACTIVE`, `REQUEST-UNREADABLE`, `REQUEST-WRITE-FAILED`, plus `ResolveTarget`'s own codes passed through. A failed atomic write also reports `refused` (`REQUEST-WRITE-FAILED`), so there is one refusal shape. `ReplaceExisting` leaves the original bytes in place on error.
- **D-14 (DECIDE & STATE):** A `## <name>` section placed before an anchor is written as `## <name>\n\n<body>\n\n`. One placed at the end of the file gets one blank line before it, with no extra blank line if the file already ends with one. The `--from` body is trimmed of blank lines at both ends before it is written, which is the same trim used for the "same body" comparison. Without that shared trim, a no-op repeat would compare unequal.
- **D-15 (DECIDE & STATE):** After building the new body, `req append-section` re-reads `VisibleSections` and refuses unless exactly one visible column-0 `## <name>` results. An unclosed fence or comment at the end of the file would otherwise hide the new heading, and every retry would append another invisible copy. The check is 6 lines, earned by the `lifecycle-section-evidence` family. Not separately tested, because the brief caps the tests at five.
- **D-16 (DECIDE & STATE):** In the `frontmatter set` file-path form, the path is checked to sit under `do-work/working/` or `do-work/queue/` after `filepath.Clean`, so `..` cannot escape. No `ResolveTarget` identity check is applied to a path, because the brief's D-05 asks only for the location rule. Symlinks are refused by `ReplaceExisting` (regular files only).
- **D-17 (integrator, DECIDE & STATE):** D-11's out-of-boundary test edit is accepted as a necessary scope addition: the registry-count test pins the handler count, and the decided D-01 design adds one handler. The integrator stamped `builder_handback_at`, `integration_at` and appended this REQ's `## Implementation Summary` with the merged `frontmatter set` and `req append-section` commands (a same-body repeat exited 0 with no change; a second `integration_at` call refused `FRONTMATTER-STAMP-EXISTS`).
- **D-18 (integrator, after review, DECIDE & STATE):** Review F1 (a `--from` body holding its own `## ` heading either refused with a misleading reason or wrote a stray section `advance` later refuses; the shipped Testing template starts with `## Testing`, and `work.md` now points that template at the command), F4 (the Timestamp rule claimed to be the only command site for an instant) and F5 (the exceptions sentence said every exception has an owning command) were fixed on the builder branch (`fd8c934d`) and re-merged with the same `<pre>` (`8a33190c`), instead of being released as known defects: `req append-section` refuses `SECTION-BODY-HAS-HEADING` before any read of the REQ, one assertion pins it, and two sentences in `work-reference.md` were corrected. F2, F3, F6 and F7 stay report only. Also corrects D-12: a repeat on a CRLF file is a no-op (the inserted text is LF, so it compares equal), not a conflict; the file keeps mixed line endings (review F6).

## Discovered Tasks

(from the builder hand-back; the builder's `impact-minor` token is not in the impact vocabulary, so the integrator re-stamped each line by the two questions in `actions/review-work.md` Step 10)

- **impact-rule-change** The **Timestamp rule** paragraph in `actions/work-reference.md` says it is "the only place in `actions/` that spells a command for obtaining" an instant. `frontmatter set ... --at now` (now named in **Stamps are append-only**) also obtains one while writing it. The sentence could say "obtaining an instant to hold", or name the writer as a second sanctioned site. → report only
- **impact-negligible** REQ-658 (finalize auto-manifest) also resolves a REQ id to its working file. `resolveActiveRequest` in `corehelpers/request_writers.go` is unexported; if REQ-658 needs the same working-or-queue rule it should call `requeststate.ResolveTarget` plus a `TreeSection` check the same way, or a later REQ can export this one. → report only
- **impact-negligible** `actions/fan-out-reference.md:161` repeats the dispatch stamp rule without naming `frontmatter set` (deliberately not edited, per D-09; the hand write stays valid). → report only

## Qualification

**Gate records (`advance --diff-range 1660977c..a56cd9c4`):** `qualify` satisfied (success), `scope-drift` satisfied (success). Two `QUALIFY-NEW-FILE-UNWIRED` warnings, judged false positives: `request_writers.go` is a same-package Go file whose `handleFrontmatterSet` and `handleRequest` are called from `commands.go` (`handleFrontmatter` dispatch and the `Handlers()` map), and `request_writers_test.go` is a `_test.go` file the Go test runner finds by convention. No debug-artifact, P-A-U or output-primitive findings. The first run reported `SCOPE-UNDECLARED-TOUCH` for `internal/corehelpers/commands_test.go`; the integrator accepted that file as a necessary scope addition (builder D-11, integrator D-17) and declared it in `## Scope` and `write_set`, after which scope-drift passed.

**Scope:** declared `write_set` (8 paths after the D-17 addition) equals the touched set in `git diff 1660977c..a56cd9c4 --stat` (8 files, 542 insertions, 22 deletions). No `do-work/` path on the builder branch (queue guard printed nothing before the merge). The merge needed no conflict resolution; REQ-660's `result_model.go` and `fan-out-reference.md` edits are not touched by this diff.

**Requirement trace (read against the merged files):**
1. `frontmatter set <file|REQ-N> <field> <value>` and `--at now`: `commands.go` sends `set` to `handleFrontmatterSet` (`request_writers.go`) before the `get` path; `--at` is accepted only as `--at now` on a field ending in `_at` (`FRONTMATTER-AT-INVALID` otherwise), and `--at now` takes `requestmodel.CanonicalTimestamp(time.Now())`.
2. Append-only: a non-empty `*_at` field refuses `FRONTMATTER-STAMP-EXISTS` before any write; every refusal returns before `replaceRequestFile`, so the bytes stay identical. Checked live on this REQ: a second `frontmatter set REQ-659 integration_at --at now` refused and left the file unchanged.
3. Timestamp rule and quoting: an explicit `*_at` value must round-trip through `time.Parse(RFC3339)` to the canonical form (`FRONTMATTER-STAMP-NOT-CANONICAL`); values are written by `RequestDocument.SetScalar`, the existing Quoting-contract writer, and list or nested fields refuse (`FRONTMATTER-FIELD-STRUCTURED`).
4. `req append-section`: a new top-level `req` command (`CommandRequest`) refuses a non-canonical name, inserts before the first later canonical section from `requestmodel.SectionsAfter` (or at the end), returns success with no write on a same body (shared `trimBlankLines`), refuses `SECTION-CONFLICT` on a different body and `SECTION-DUPLICATE` on two copies, and re-checks that exactly one visible copy results. Checked live: this REQ's `## Implementation Summary` was written by it, and a repeat exited 0 with no change. `advance` now reads its twelve later-section lists from the same `CanonicalSectionOrder`; each replaced list equals the old literal tail (read in the diff).
5. Resolution: `resolveActiveRequest` maps `REQ-N` through `requeststate.ResolveTarget` (zero or several matches refuse with its own codes) and refuses a file outside `do-work/working/` or `do-work/queue/` (`REQUEST-NOT-ACTIVE`); a path argument is checked for the same location after cleaning.
6. Prose: `work.md` Request Files as Living Logs holds the one `req append-section` argv; the stamp sentences for `planning_at`, `dispatch_at`, `integration_at` and `review_at`/`remediation_at`/`re_review_at` name `frontmatter set` and cite **Stamps are append-only**; the Testing append names `req append-section`; `work-reference.md` → **Stamps are append-only** holds the one `frontmatter set ... --at now` argv. Each keeps the hand write valid.
7. Release: handled at finalization (patch bump, own changelog entry).

**Red-Green Proof trace:** `TestFrontmatterSetStampIsAppendOnly` (first call adds `review_at` once, second refuses with identical bytes), `TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses` (same file twice leaves one section; a different file refuses), `TestRequestAppendSectionLandsInCanonicalOrder` (`## Testing` lands before an existing `## Review`), plus `TestFrontmatterSetWritesQuotedScalars` and `TestRequestIDResolvesOneWorkingOrQueueFile` for requirements 3 and 5.

**Remediation re-merge (review F1, F4, F5), cumulative range `1660977c..8a33190c`:** `advance` no longer takes qualify input at this phase, so the integrator ran the same `qualify --request-path <P> --diff-range 1660977c..8a33190c` handler directly: success, only the two judged `QUALIFY-NEW-FILE-UNWIRED` warnings above. The touched set is still the 8 declared files (`git diff 1660977c..8a33190c --stat`: 554 insertions, 23 deletions); the delta touches `request_writers.go`, `request_writers_test.go` and `work-reference.md` only. The queue guard before the re-merge printed nothing.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `a56cd9c4` (machine quiet before launch: 1-minute load 2.22, no other gate running; load 6.60 at the end), then `advance REQ-659 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-659-probe.sh`.
**Result:** ✓ Repository gate passed on the first run (exit 0, gate wall 151 s, do-work-cli 893 tests in 69 s, slowest file `internal/finalization/finalization_recovery_test.go` 22.05 s under the 30 s limit; queue-kanban 421 tests, slowest 20.21 s). Probe exit 0 ("REQ-659 GREEN probe: pass"). Gate records `green-gate`, `scope-drift` and `run-blocked-check` satisfied.

**After the review-fix re-merge (`8a33190c`):** the same gate argv ran again (load 4.48 before launch, 5.11 at the end): exit 0, gate wall 149 s, 893 do-work-cli tests, slowest file 21.43 s. `advance` refuses gate input past this phase, so the green record was written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` at `8a33190c`, and `REQ-659-probe.sh` was run directly: exit 0. `go test -count=1 ./internal/corehelpers/` in the builder worktree: exit 0.

**Red-green validation:** (from the builder hand-back, traced to `## Red-Green Proof`)
- `TestFrontmatterSetStampIsAppendOnly`, `TestFrontmatterSetWritesQuotedScalars`, `TestRequestAppendSectionLandsInCanonicalOrder`, `TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses`, `TestRequestIDResolvesOneWorkingOrQueueFile` (`internal/corehelpers/request_writers_test.go`): ✗ before the handlers existed (exit 1; the `set` tests got `HELPER-USAGE: usage: frontmatter get <file> <field> ...`, the `req` tests got `no handler registered for "req"`; the constant was added first so RED is a runtime failure, not a compile error) → ✓ after (exit 0, no skips).
- `TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses`, new assertion for review F1 (a `--from` body holding `## Review`): ✗ before the fix (`request_writers_test.go:211: outcome=success findings=[], want refused (exit 1)`) → ✓ after (exit 0); the re-reviewer confirmed it fails again with the check disabled.
- Advance guard: `TestAdvanceCommandPhaseMatrix` and `TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates` passed before and after the twelve literal lists were replaced by `SectionsAfter`; the advance tests were not edited.

**New tests added:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go` (five tests through `Handlers()`)

**Existing tests updated:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (`TestEveryRemainingUtilityHasOneHandler` names `CommandRequest`; handler count 21 → 22)

**Heavy verification plan:**
- Range: 1660977cd29cf853a5f5abee540ee604a5f1a3b9..8a33190caa8c3c4f6b41326b46f666ff668a5839 (re-planned after the review-fix re-merge; same four lanes as the first plan at `a56cd9c4`)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files under `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files under `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files under `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 91%** | 2026-10-10T14:54:26Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 88% |
| Test Adequacy | 88% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token, the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. `req append-section` does not check the `--from` body for a column-0 `## ` heading (`internal/corehelpers/request_writers.go:109`, `:158-166`). The Testing Section Template's own `## Testing` line gets the misleading reason "would not be a visible section". A different canonical heading (for example `## Review`) is written unchecked. Fix: refuse any visible heading in the body. impact-user-visible → report only
- F2. The path-form containment check is lexical (`request_writers.go:178-180`). A symlinked directory under `do-work/working/` lets `frontmatter set` write outside the tree (reproduced), and an absolute path through the macOS `/tmp` alias falsely refuses. Fix: `filepath.EvalSymlinks` before `Rel`. impact-negligible → report only
- F3. `frontmatter set` replaces any non-`*_at` scalar with no check on the value (`status bogus-status` exit 0). Fix: refuse `status` and `id`, and name their owning commands. impact-rule-change → report only
- F4. The Timestamp rule says it is "the only place in `actions/` that spells a command for obtaining one" (`actions/work-reference.md:61`). `:73` now spells `frontmatter set ... --at now`. impact-rule-change → report only
- F5. `actions/work-reference.md:73` says the four exceptions "keep their owning commands", but `abandon.md:62` (`completed_at` re-stamp) and `clarify.md:138` and `:226` (`status_changed_at`) are hand writes. The refusal is safe. impact-negligible → report only
- F6 (nit). D-12 says a CRLF repeat refuses as a conflict. It is actually a no-op, and the file ends with mixed line endings. impact-negligible → report only
- F7 (nit, anti-bloat). Additions are low. `requestIDPattern` duplicates `commands.go:159`. The no-change write guard only fires on a same-value `set`. JSON `set` does not report the written instant (this matches `get`). impact-negligible → report only

**Acceptance:** Pass. UR-145 check 2 was reproduced on the real binary. Refusals leave the bytes identical. The 12 `advance` lists equal `SectionsAfter` tails exactly.
**Restatement sweep:** redefined (E1) how REQ `*_at` stamps are written (`frontmatter set`, append-only enforced), (E2) how lifecycle sections are written (`req append-section`), (E3) the section order as `requestmodel.CanonicalSectionOrder`. Stale: `work-reference.md:61` "only place … spells a command" (F4). `work-reference.md:73` "owning commands" vs `abandon.md:62`, `clarify.md:138`, `:226` (F5). Consistent: `work-reference.md:147`, `:155`, `:178`, `:182`, `:199`, `:209` (schema comments citing the rule). `fan-out-reference.md:163`, `:167` (absent-only `builder_handback_at` / `dispatch_at`, hand path valid, D-09). `work.md:21`, `:201`, `:288`, `:309`, `:342`, `:375`, `:409`, `:497`. `sample-archived-req.md` section order (matches. It omits Pre-Flight, which sits between Scope and Implementation Summary in the list and in `work-reference.md:494`). `review-work.md` Append to REQ File (`## Review` added at the end). `queue-kanban/frontmatter_cli.go:33` and `prime-do-kanban.md:3` (about queue-kanban, which stays read-only). Go: no other hard-coded section-order list. `advance_commands.go:174` and `:253` are a set check and a single check, not order tails.
**Suggested testing:** 4 items
**Follow-ups created:** None (7 findings report only)

Full report: `do-work/runs/work-2026-10-10-131527/REQ-659-review.md`.

*Reviewed by review-work action*

### Re-review (delta a56cd9c4..8a33190c) | 2026-10-10T15:02:18Z

F1, F4 and F5 were fixed on the builder branch and re-merged with the same `<pre>` (D-18).

- Overall: 94% | Acceptance: Pass (real binary at `8a33190c`: `## Testing` verbatim, a different `## Review`, a CRLF heading and a 2-space indented heading refuse `SECTION-BODY-HAS-HEADING` with bytes unchanged; `###` subheadings, fenced or commented `##` lines, 4-space code and `# ` headings are accepted; the new assertion fails with the check disabled)
- F1: closed. F4: closed. F5: closed.
- N1 `req append-section` accepts a body that leaves a fence or `<!--` comment open, which hides every later section (for example `## Review`) from `advance`; the visible-copy re-check counts only the inserted section (`request_writers.go:161-171`). Fix: refuse unless the column-0 visible section count grows by exactly one — impact-user-visible → report only

Full re-review report: `do-work/runs/work-2026-10-10-131527/REQ-659-rereview.md`.

## Lessons Learned

**What worked:** Moving `advance`'s twelve inline later-section lists into one `requestmodel.CanonicalSectionOrder` before the writer used it meant `req append-section` cannot place a section where `advance` would refuse it; a script compared each replaced list with its `SectionsAfter` tail, so the refactor is provably behaviour-identical. Reusing `ResolveTarget`, `SetScalar` and `atomicfile.ReplaceExisting` kept both commands to one new file with no new encoder. The integrator used both commands on this REQ (stamps and four sections); the append-only refusal and the same-body no-op worked live.
**What didn't:** The writer trusted the shape of the `--from` body. The shipped Testing template starts with `## Testing`, and `work.md` now points that template at the command, so a verbatim template got a misleading refusal, and a body holding a different `## ` heading wrote a stray section (review F1, fixed). The same blind spot remains for a body that leaves a fence or comment open (re-review N1, report only). A registry-count test is part of a new command's seam: the builder had to edit `commands_test.go` outside its write boundary (D-11).
**Worth knowing:** A section writer must validate the body it inserts as well as the file it inserts into: the body is untrusted Markdown that can carry headings, fences and comments, and each of them changes which sections `advance` sees. The qualifier reads a backticked bare word in `## Scope` (for example `req`) as a declared path, so write command names there without backticks.

## Orientation

Now a REQ stamp is written with `do-work-cli frontmatter set <file|REQ-N> <field> (<value> | --at now)` and a lifecycle section with `do-work-cli req append-section REQ-N --section <name> --from <file>`; both live in the do-work-cli `internal/corehelpers` package (`request_writers.go`), and the section order they share with `advance` is `requestmodel.CanonicalSectionOrder`. They are named in `actions/work.md` → Request Files as Living Logs and `actions/work-reference.md` → **Stamps are append-only**, where the hand write stays valid. [MAP CHANGED]: two new CLI commands and one shared section-order list. Prime spot-check: `prime-shell-commands.md`, `prime-action-files.md` and `prime-releases.md` name no path this change moved or removed; none is stale.

## Heavy Verification Plan

- Base: 1660977cd29cf853a5f5abee540ee604a5f1a3b9
- Target: 8a33190caa8c3c4f6b41326b46f666ff668a5839
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: 8a33190caa8c3c4f6b41326b46f666ff668a5839
- Execution revision: 8a33190caa8c3c4f6b41326b46f666ff668a5839 (detached checkout `.git/work-run-work-2026-10-10-131527/drain-head-REQ-659`, `QUEUE_KANBAN_BROWSER` set, removed afterwards)
- do-work-cli-integrations: exit 0, executed, 68 s
- staged-skills: exit 0, executed, 39 s
- updater: exit 0, executed, 67 s
- installer: exit 0, executed, 30 s
- Earlier drain at the first merge `a56cd9c4` (before the review fix): all four lanes exit 0, executed (80, 45, 70, 33 s).

## Timing

Observed 2026-10-10T14:41:46Z to 2026-10-10T15:02:19Z: 20m 33s total, 29m 09s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 15m 35s | 2 |
| verification-gate | 13m 09s | 4 |
| handback-merge | 25s | 1 |

Slowest stage: review / re-review of the review-fix delta, 8m 15s, outcome success.

Note: no builder-work event was recorded. The hand-back had landed before this integrator started (dispatch 13:24:46Z, builder commit 13:31:55Z, about 7 minutes of build), so per fan-out-reference.md → Landed hand-back the event is skipped rather than charging the builder with the wait. Review and verification overlap in wall time because the heavy drain ran while the first reviewer worked, so attributed time exceeds the observed span.
