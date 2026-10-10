---
id: REQ-659
title: 'frontmatter set and req append-section write REQ stamps and sections with the schema rules enforced'
status: claimed
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
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go, skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go, skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md]
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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

<!-- D-XX counter: last used D-10. Next decision: D-11. -->

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
