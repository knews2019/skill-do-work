---
id: REQ-032
title: Write-set contract, parallel-dispatch gate, and serial-only resource classes
status: completed
commit: 97f65b9
completed_at: 2026-07-28T21:57:21Z
claimed_at: 2026-07-28T21:11:40Z
created_at: 2026-07-28T20:52:05Z
user_request: UR-007
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: [REQ-033, REQ-034]
batch: parallel-dispatch
write_set:
  - actions/work-reference.md
  - actions/work.md
  - actions/capture-reference.md
  - actions/board.md
  - docs/work-guide.md
  - tools/queue-kanban/model.go
  - tools/queue-kanban/generate.go
  - tools/queue-kanban/model_test.go
  - _dev/tests/contract-regressions.sh
  - CLAUDE.md
maintenance: false
---

# Write-set contract, parallel-dispatch gate, and serial-only resource classes

## What

Give the work pipeline a first-class notion of *which files a REQ will write* (`write_set`), and a dispatch rule that lets advanced harnesses run multiple ready REQs concurrently when their write-sets are pairwise disjoint. Certain resource classes are never parallelized regardless of disjointness.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Follow the `## Plan` section verbatim, using the `## Exploration` line anchors. Order: schema first (`actions/work-reference.md` frontmatter block + verbatim-list sentence + Scope Declaration Template note), then the capture seed (`actions/capture-reference.md`), then the pipeline rules (`actions/work.md`: Step 1 gate subsection after the Targeted-mode paragraph, Step 5.5 mirror, Step 6 builder bullet, one Rules bullet, one Common Rationalizations row), then the lock-step enumerations (`actions/board.md`, `CLAUDE.md`), then the user-facing mention (`docs/work-guide.md`), then the parser (`model.go` struct field + parse line; `generate.go` payload field + assignment), then tests (`model_test.go` fixture + `reflect.DeepEqual` assertion; two new ratchets in `_dev/tests/contract-regressions.sh`). No new abstractions: `WriteSet` reuses `coerceToStringList`, no alias resolution (no alias exists), no derived overlap computation (that is REQ-034). `prime_files` is empty — no primes to read. Verify with `go test ./...`, `bash _dev/tests/contract-regressions.sh`, and the REQ's GREEN grep. SKILL.md untouched (2-word budget headroom).
- [x] **[APPLY]:** Implemented as planned across exactly the ten declared files — no additions, no deletions, no new files. One in-flight correction: the Step 1 subsection was first inserted before the Targeted-mode paragraph and was moved to after it, per Exploration concern #1 (do not split the dependency-gating cluster).
- [x] **[UNIFY]:** `git diff --stat` shows 10 files, +60/-2, all declared. Reviewed each: `actions/work-reference.md` (schema line placed after `depends_on`, verbatim-list sentence below the Write-paths paragraph, Scope-template mirror note — Schema Read Contract still says "Nine fields", correct since no row was added), `actions/work.md` (gate subsection at Step 1 line 180 after Targeted-mode, Step 5.5 mirror, Step 6 builder bullet, Rules bullet, Rationalizations row), `actions/capture-reference.md` (field bullet + Populating paragraph), `actions/board.md` + `CLAUDE.md` (lock-step enumerations gain `write_set`), `docs/work-guide.md` (one user-facing bullet + stale count fix, D-02), `tools/queue-kanban/model.go` (struct field + parse line), `generate.go` (payload field + assignment), `model_test.go` (fixture + DeepEqual assertion), `_dev/tests/contract-regressions.sh` (three ratchets). Linters: `go vet ./...` clean, `gofmt -l` empty. No debug artifacts — no `console.log`/`fmt.Print`/`TODO`/commented-out code in the diff; no build outputs staged (scratch binary built and deleted outside the repo). No shipped file cites `CLAUDE.md`/`AGENTS.md`; SKILL.md untouched.

## Why (if provided)

Observed in a consumer project: builds serialize even when REQs touch disjoint files, while the verification phase (Playwright runs) is long and perfectly parallel. The runner had to invent single-writer discipline ad hoc; the skill should make it a contract.

## Detailed Requirements

- New optional frontmatter field `write_set:` (list of repo-relative paths/globs). Capture may seed it; the Pre-Flight Scope step ("Files I will touch") firms it up — the existing Scope prose and the frontmatter field must not drift from each other.
- **Contract:** a builder writes only inside its declared set. Discovering mid-build that it needs a file outside the set means **stop and report to the orchestrator** (which may extend the set if no concurrent REQ conflicts), never silently write. The report lands in the REQ's trail.
- **Dispatch gate** in `actions/work.md`: where the harness supports concurrent subagents, the orchestrator may dispatch multiple ready REQs (deps met, not blocked) simultaneously **iff** their write-sets are pairwise disjoint. Overlapping REQs serialize, or get an explicit partition directive at dispatch time.
- **Serial-only resource classes:** REQs touching migration/schema surfaces, dependency lockfiles, or generated bundles are never dispatched concurrently with each other, disjoint or not — semantic collisions (e.g., two migrations taking the same sequence number) pass through textual disjointness and git merges untouched. State the *condition* (ordered or generated resources whose correctness depends on global sequence); the class list is illustrative, per the Closed Enumerations Go Stale rule.
- Write-sets are a **scheduling input, not a safety guarantee** — phrase accordingly (REQ-033's worktree isolation is what makes overlap survivable; without it, overlap means serialize).
- `tools/queue-kanban/model.go` parses `write_set` for display in the same commit (parser lock-step). Display only — no column logic.
- Rejected alternative to record in the change: timed per-file lock files (TTL breaks over live slow agents — the same defect class as the 15s orchestrator-mutex break removed in 0.140.4).

## Constraints

- Design for the floor: floor agents ignore the gate and run the serial loop unchanged; absence of `write_set` means "unknown — treat as overlapping everything" (safe default).
- No new action file; changes go into `actions/work.md` / `actions/work-reference.md` / `actions/capture-reference.md` (schema). SKILL.md's word budget applies to any routing prose.
- `_dev/tests/contract-regressions.sh` may gain a ratchet asserting the gate text exists; keep the hook's existing checks green.

## Dependencies

Root of the batch — REQ-033 and REQ-034 build on this field and gate.

## Builder Guidance

Certainty: Firm on the contract and gate; latitude on exact field syntax (paths vs globs) and where the schema documentation lives (capture-reference vs work-reference Schema Read Contract) — follow existing schema-field precedent. `- [~]` Whether `write_set` is required or optional on new REQs → deferred to builder; recommended optional with the safe default above.

## Red-Green Proof

**RED prompt/case:** `grep -ri "write_set" actions/ tools/queue-kanban/model.go` returns nothing; `actions/work.md` contains no rule permitting concurrent dispatch of ready REQs — two ready REQs with disjoint file scopes are always processed one at a time.
**Why RED now:** The pipeline has no vocabulary for declared write surfaces, so an orchestrator has no sanctioned basis to parallelize.
**GREEN when:** The schema defines `write_set`; `actions/work.md` contains the disjointness dispatch gate and the serial-only class rule; `model.go` parses the field; a queue REQ carrying `write_set` renders without parser warnings on the board.
**Validation:** User confirmed (design agreed in conversation; "run the capture").

## Full Context

See `do-work/user-requests/UR-007/input.md` for complete verbatim input.

---
*Source: "can we do something to update this skill to encourage them to run in parallel? …" (UR-007)*

Think carefully before answering.

---

## Triage

**Route: C** - Complex

**Reasoning:** New pipeline feature spanning multiple systems — frontmatter schema (`actions/work-reference.md`), dispatch rule (`actions/work.md`), capture seeding (`actions/capture-reference.md`), and the Go board parser (`tools/queue-kanban/model.go`) in lock-step. Architectural contract work, not a localized change.

**Planning:** Required

## Plan

### 1. `actions/work-reference.md` — schema first (everything else cites it)

**Full Frontmatter block (§ Request File Schema, line ~97):** add `write_set: []` immediately after `depends_on:` in the "Set by capture action" group — it's list-valued and capture-seeded, so it belongs with its siblings. Annotation states: optional list of repo-relative paths/globs the REQ expects to write; **absent or empty ⇒ unknown ⇒ treat as overlapping every other REQ** (safe default); firmed by Step 5.5's Scope declaration; a **scheduling input, not a safety guarantee** — nothing enforces it at the filesystem, so overlap means serialize.

**Schema Read Contract:** no new row. The contract governs enum/boolean fields; `write_set` is a free-form list like `depends_on`/`related`/`prime_files`, which carry no rows either. Add one sentence under the table's write-paths paragraph noting list-valued path fields are read verbatim (no normalization, no alias) so a future reader doesn't look for a missing rule.

### 2. `actions/capture-reference.md`

Add `write_set: [...]` to **Additional frontmatter for complex requests** (beside `related`/`batch`), plus a short **Populating `write_set`** paragraph next to the existing `depends_on` one: seed only when the request names files or the slice is per-file; omit rather than guess — an invented set is worse than absence, because absence already means "overlaps everything." Note the field is refined by the work pipeline's Scope step, so capture's value is a hint, not a commitment.

### 3. `actions/work.md` — three edits

**a. Step 1, new subsection after the dependency-gating cluster (after "Targeted mode bypasses dependency gating"), before "Queue status summary":** `**Parallel dispatch (optional — advanced harnesses).**` Frame explicitly as an option over the default: floor agents ignore this and claim one REQ per loop, unchanged. Where the harness runs concurrent builders, the orchestrator MAY claim and dispatch several dependency-ready REQs at once **iff** their `write_set`s are pairwise disjoint. Any REQ missing/empty `write_set` overlaps everything ⇒ serialize it. Overlapping pairs either serialize or get an explicit partition directive at dispatch (each builder told which subset of the shared files it owns). Each concurrently dispatched REQ still runs Steps 2–9 in full, including the lock's `claimed_req` — concurrency is in dispatch, not in bookkeeping. One clause records the rejection: declared write-sets rather than timed per-file locks — a TTL expires over a live slow agent and hands the file to a second writer.

**b. Same subsection, serial-only classes.** State the *condition*: a REQ is serial-only when it writes an **ordered or generated resource whose correctness depends on global sequence**, so two textually disjoint edits still collide semantically (two migrations claiming one sequence number merge cleanly and are still wrong). Two serial-only REQs are never co-dispatched even with disjoint sets; a serial-only REQ may still run beside an unrelated non-serial one. Class list marked **illustrative, not exhaustive**: migration/schema surfaces, dependency lockfiles, generated bundles/codegen output, ordered seed or fixture data.

**c. Step 5.5 + Step 6 (the anti-drift and binding half).** Step 5.5: the Scope section's "Files I will touch" list is the authoritative firming of `write_set`; after writing `## Scope`, the orchestrator mirrors that list back into the frontmatter field so the two cannot drift (Scope is the source, frontmatter the mirror — one write direction only). Route A skips Scope, so its `write_set` stays as captured (or absent). Step 6 builder bullets: add one — **write only inside the declared `write_set`; discovering a needed file outside it is a stop-and-report to the orchestrator, never a silent write.** The orchestrator records the request and its resolution in the REQ trail (a `## Decisions` D-XX entry, since it is a scope judgment) and extends both the Scope list and `write_set` only when no concurrently dispatched REQ claims that file.

Also add one Rules bullet ("write-sets schedule work; they don't protect files") and one Common Rationalizations row: *"Their write-sets don't overlap, so I can run both migration REQs together"* → serial-only classes never co-dispatch → sequence collisions survive textual disjointness and git merge.

### 4. `tools/queue-kanban/model.go` + `generate.go` (lock-step, same commit)

`model.go`: add `WriteSet []string` to `RequestTicket` right after `Related` (same comment style — display only, never column logic), populated in `parseRequestTicket` with `WriteSet: coerceToStringList(fields["write_set"])`. No alias, no normalization, no derived annotation — the `UnmetDependencies`-style overlap computation is REQ-034's.

`generate.go`: add `WriteSet []string` with json tag `writeSet` to `generatedRequest` and the assignment in the payload builder. That is the minimum "display" requires — the field must reach the JSON the board renders from; `board.js`/`board.css` badge rendering is REQ-034 and needs no change here (the frontend ignores unknown keys).

### 5. Testing

- `cd tools/queue-kanban && go test ./...`; extend `TestParseRequestTicketNormalizesAndResolves`'s fixture with `write_set: [src/a.ts, src/b.ts]` and a `reflect.DeepEqual` assertion (cheapest correct home — it already covers list fields).
- `bash _dev/tests/contract-regressions.sh` must stay green. Optional new ratchet (recommended, cheap): assert `actions/work.md` contains the gate text and `tools/queue-kanban/model.go` contains `write_set` — that pairing is exactly the parser lock-step rule this REQ introduces.
- **SKILL.md: no change.** Verified 2648 words against a 2650 budget — two words of headroom, so any routing prose would fail the ratchet. The gate is a work.md-internal option and needs no routing row.
- Version bump in `actions/version.md` + `CHANGELOG.md` entry (title states what shipped), with a bullet recording the rejected timed-lock alternative — handled by the orchestrator at the Commit Phase.

### Architectural decisions

1. **Scope section is the source, frontmatter the mirror** — one-directional sync is the only shape that can't drift.
2. **No Schema Read Contract row** — the contract is for enum/boolean values; a path list has no canonical vocabulary to normalize against.
3. **Absence = overlaps everything** — makes the field purely additive; every existing REQ keeps today's serial behavior.
4. **Gate lives inside Step 1, not a new step** — the ready set is computed there, and a new step would read as mandatory in the Orchestrator Checklist.
5. **Payload field but no frontend** — satisfies "display" without claiming REQ-034's badge surface.
6. **Serial-only stated as condition + illustrative list** — per the Closed Enumerations Go Stale rule; the class list will grow.

### Not covered / flags

- Nothing enforces the write-set mechanically (no pre-write check). That is deliberate per the REQ ("scheduling input, not a safety guarantee") — real isolation is REQ-033's worktrees.
- The gate text stays harness-agnostic (says "concurrent builders", never "worktree") so REQ-033 can add its dispatch mode without rewriting this section.

**Plan validation (orchestrator):** All Detailed Requirements map to planned tasks; no orphan tasks (the optional ratchet traces to the REQ's Constraints). ⚠ Scope spans ~9 files — above the 3-task quality flag, inherent to a lock-step schema change; explorer confirmed each edit site. Exploration added three sites the plan missed: `actions/board.md`'s lock-step field enumeration, the Scope Declaration Template note in work-reference.md, and `docs/work-guide.md` staleness — folded into Scope below. No file conflicts (no other REQ in flight).

*Generated by Plan agent*

## Exploration

### 1. `actions/work-reference.md` — schema

- **Full Frontmatter block**: `## Request File Schema — Full Frontmatter` at **line 79**; `depends_on: []` is **line 97**, the last line of the `# Set by capture action` group (next line 99 is `# Set by work action when claimed`). Insert `write_set: []` at line 98. Style to imitate: single YAML line + long trailing `#` annotation (see `depends_on`/`maintenance` lines).
- **Schema Read Contract** (line 154): table is explicitly "Nine fields … enum-or-boolean-valued". `related`, `depends_on`, `prime_files`, `blocked_by` have **no rows** — `write_set` needs none. The sentence about verbatim list reads belongs after the **"Write paths are unaffected."** paragraph (line ~189, immediately below the table).
- **Scope Declaration Template (Step 5.5)** lives here at **line 480** (`**Files I will touch:**` bullets, lines 485–488). The mirror rule makes this template the natural place for a one-line note that the list is the source of `write_set`.

### 2. `actions/work.md`

- Step 1 `### Step 1: Find Next Request` line 115. `**Dependency-aware selection.**` is **line 164**, but four more paragraphs follow before `**Queue status summary:**` (line 180): roots (166), cycle detection (168), wave execution (170–176), `**Targeted mode bypasses dependency gating.**` (178). Insert the new subsection *after line 178* so it doesn't split the dependency-gating cluster.
- Step 5.5 at **line 311**; the `## Scope` write instruction is lines 317–319. Mirror instruction goes right after 319. Route A skip is line 315.
- Step 6 builder bullet list: lines ~365–381, all `- **Bold lead-in:** …` one-liners. Best sibling anchor: the `## Discovered Tasks` out-of-scope bullet — put the write-set bullet adjacent.
- `## Rules` line 660 (3 plain bullets, before "Common mistakes to avoid"); `## Common Rationalizations` table line 685 with 9 rows.
- `tools/checks/scope-drift.sh` reads the `## Scope` section — mirroring into frontmatter doesn't affect it.

### 3. `actions/capture-reference.md`

- Base frontmatter template: `depends_on: []` at **line 22**. `**Additional frontmatter for complex requests:**` at **line 103** (bullets 104–106: `related`, `batch`, `addendum_to`). `**Populating `depends_on`.**` at **line 110**; `**Slicing convention.**` at 112.
- `### Schema Aliases` table (lines ~118–125) enumerates canonical→alias→consumers. `write_set` has no alias — skip consciously.

### 4. `tools/queue-kanban`

- `model.go`: `Related []string // soft relations (not dependency edges)` **line 114** (struct), populated at **line 598** `Related: coerceToStringList(fields["related"]),` inside `parseRequestTicket` (line 542). `resolveDependsOn` line 736 (alias-only helper — not needed here). Insert `WriteSet` after line 114 and after 598.
- `generate.go`: `generatedRequest` struct line 95; `Related` json field **line 109**; assignment **line 267**. `.Related` appears **only** in generate.go — serve mode reuses this payload, so those two edits are the complete data path.
- Tests: `model_test.go` line 272 `TestParseRequestTicketNormalizesAndResolves`, raw-string fixture (275–290) already contains `related: [REQ-501]` with no assertion; `reflect.DeepEqual` + `t.Fatalf` style at line 311. `frontmatter_test.go` covers block-list recovery generically via `coerceToStringList` — no change needed.

### 5. `_dev/tests/contract-regressions.sh`

Pattern: `assert_contains "<repo-relative path>" '<ERE pattern>' '<failure message ending in a period>'` (helper line 8; examples lines 81–100). Router budget check lines 158–169: `router_word_budget=2650`, current `wc -w SKILL.md` = **2648** — two-word headroom confirmed; do not touch SKILL.md.

### 6. Other schema enumerations (Closed Enumerations rule)

- `actions/board.md` **line 115** — lock-step paragraph enumerating display-only parsed fields (`domain`, blocked fields). Adding a parsed field makes it stale — must gain `write_set`.
- `CLAUDE.md` line 158 — same enumeration (maintainer doc; update to keep lock-step doc accurate).
- `docs/work-guide.md` lines 86–89 — user-facing prose; a brief parallel-dispatch mention belongs there. `docs/capture-guide.md` line 54 — judgment.
- `actions/roadmap.md` line 60 alias enumeration — no `write_set` alias, skip. `actions/sample-archived-req.md` — optional field, correctly absent.

### 7. Concerns

1. Step 1 insertion point: after the Targeted-mode paragraph (line 178), before Queue status summary — not immediately after Dependency-aware selection.
2. `actions/board.md` § lock-step and the Scope Declaration Template are edit sites the plan missed — folded into Scope.
3. Everything else on disk matches the plan (no Schema Read Contract row needed; `.Related` single data path; word budget).

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `actions/work-reference.md` (modify) — `write_set` in Full Frontmatter block; verbatim-list sentence after the Write-paths paragraph; one-line source-of-`write_set` note in the Scope Declaration Template
- `actions/work.md` (modify) — Step 1 parallel-dispatch gate + serial-only classes subsection; Step 5.5 mirror instruction; Step 6 builder write-set bullet; one Rules bullet; one Common Rationalizations row
- `actions/capture-reference.md` (modify) — `write_set` bullet under Additional frontmatter; Populating `write_set` paragraph
- `actions/board.md` (modify) — lock-step enumeration gains `write_set` (display-only)
- `docs/work-guide.md` (modify) — brief user-facing mention of the optional parallel-dispatch gate
- `tools/queue-kanban/model.go` (modify) — `WriteSet` struct field + parse line
- `tools/queue-kanban/generate.go` (modify) — payload field + assignment
- `tools/queue-kanban/model_test.go` (modify) — fixture `write_set` + assertion
- `_dev/tests/contract-regressions.sh` (modify) — new ratchet asserting gate text + parser lock-step pairing
- `CLAUDE.md` (modify) — maintainer lock-step enumeration line

**Files I will NOT touch:** `SKILL.md` (2 words of budget headroom; no routing change needed), `next-steps.md`, `tools/queue-kanban/web/*` (badge rendering is REQ-034), `actions/capture.md` (slicing nudge is REQ-034), `actions/roadmap.md` / Schema Aliases table (no alias for `write_set`), `actions/sample-archived-req.md` (optional field, correctly absent).

**Acceptance criteria (restated from REQ):**
- [ ] Schema defines optional `write_set` (repo-relative paths/globs) with absence/empty = "overlaps everything" safe default
- [ ] Builder contract in work.md: write only inside the declared set; out-of-set need = stop and report to orchestrator (never silent write); report lands in the REQ trail
- [ ] Dispatch gate in work.md: concurrent dispatch of ready REQs iff write-sets pairwise disjoint; overlapping REQs serialize or get an explicit partition directive; floor agents unchanged
- [ ] Serial-only resource classes: condition stated (ordered/generated resources whose correctness depends on global sequence), class list illustrative
- [ ] Write-sets phrased as scheduling input, not a safety guarantee
- [ ] Scope section and `write_set` frontmatter cannot drift (one-directional mirror rule)
- [ ] `tools/queue-kanban/model.go` parses `write_set` for display in the same commit; no column logic
- [ ] Rejected timed-per-file-lock alternative recorded in the change
- [ ] `grep -ri write_set actions/ tools/queue-kanban/model.go` returns hits (RED→GREEN); go tests + contract-regressions.sh green

## Decisions

**D-01 — `write_set` is OPTIONAL on every REQ; absence means "overlaps everything." (DECIDE & STATE)**

The Builder Guidance deferred this (`- [~]` required vs. optional) with a recommendation, and I followed the recommendation. Optional is the only choice that makes the field purely additive: every REQ already in a consumer's queue keeps today's exact serial behavior with no migration, and no existing action has to start emitting a field it can't compute. Making it required would force capture to guess a write surface at the moment it has the least information, and a *guessed* set is strictly more dangerous than no set — absence fails safe (serialize), while a wrong set can green-light two builders onto the same file. The safe default is stated in three places so a reader hits it wherever they enter: the schema line (`actions/work-reference.md`), the dispatch gate's first bullet (`actions/work.md` Step 1), and capture's Populating paragraph.

**D-02 — Fixed the stale "three properties" count in `docs/work-guide.md`. (DECIDE & STATE)**

The lead-in to that bullet list already said "three properties" over four bullets before this REQ; adding the parallel-dispatch bullet made it five. Changed to "a few properties" rather than "five" — a count-free lead-in cannot go stale the next time a bullet is added, which is the same failure this one already demonstrated. Confined to the one file already in scope and to the exact paragraph this REQ extends.

**D-03 — The parser ratchet asserts `fields["write_set"]`, not the bare string. (DECIDE & STATE)**

A bare `write_set` grep against `model.go` passes on the doc comment alone, so deleting the parse line while leaving the comment would keep the ratchet green — a lock-step check that can't detect the break it exists for. Anchoring on the actual map read makes it bite. Verified by mutation: removing the parse line and rewording the two gate phrases fires all three new assertions; restoring turns them green again.

**D-04 — No `write_set` mention in `docs/capture-guide.md`. (SILENT/leaf, recorded because Exploration flagged it as a judgment call)**

That guide enumerates REQ *body sections* (Red-Green Proof, Open Questions, Detailed Requirements…), not optional frontmatter fields — it doesn't document `depends_on` as a field either. `write_set` is an agent-populated scheduling hint no user hand-writes, so it has no home there. Left untouched; it was also outside the declared scope.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — `write_set: []` schema line after `depends_on` (safe default, mirror direction, scheduling-not-safety framing); verbatim-list-read sentence under the Write-paths paragraph; Scope-template note naming "Files I will touch" as the source of `write_set`
- `actions/work.md` (modified) — Step 1 "Parallel dispatch (optional — advanced harnesses)" subsection: pairwise-disjoint gate, absent/empty set = overlaps everything, partition directive, per-REQ bookkeeping unchanged, rejected timed-lock rationale, serial-only resource classes (condition + illustrative list); Step 5.5 one-directional Scope→`write_set` mirror instruction; Step 6 builder write-boundary bullet (stop-and-report, never silent write); one Rules bullet; one Common Rationalizations row
- `actions/capture-reference.md` (modified) — `write_set` bullet under Additional frontmatter; "Populating `write_set`" paragraph (omit rather than guess; hint not commitment)
- `actions/board.md` (modified) — lock-step display-only field enumeration gains `write_set`
- `docs/work-guide.md` (modified) — user-facing parallel-dispatch bullet; count-free lead-in fix (D-02)
- `tools/queue-kanban/model.go` (modified) — `WriteSet []string` on `RequestTicket` + `coerceToStringList(fields["write_set"])` in `parseRequestTicket`
- `tools/queue-kanban/generate.go` (modified) — `WriteSet` payload field (`json:"writeSet"`) + assignment
- `tools/queue-kanban/model_test.go` (modified) — fixture gains `write_set: [src/a.ts, src/b.ts]` + `reflect.DeepEqual` assertion
- `_dev/tests/contract-regressions.sh` (modified) — three new ratchets: work.md gate text, work.md serial-only rule, model.go parse call (`fields["write_set"]`, D-03)
- `CLAUDE.md` (modified) — maintainer lock-step enumeration gains `write_set`

**What was done:** Added the optional `write_set` frontmatter field (repo-relative paths/globs; absence = overlaps everything) to the schema, capture guidance, and board parser/payload in lock-step, and gave `actions/work.md` the parallel-dispatch gate (concurrent dispatch iff pairwise-disjoint write-sets, serial-only resource classes stated as a condition), the builder stop-and-report write-boundary contract, and the one-directional Scope→`write_set` mirror rule. Three contract-regression ratchets pin the gate text and parser lock-step.

## Qualification

Passed — 10 files verified in `git diff --stat` (exact match to Scope; `tools/checks/scope-drift.sh` reports zero drift), all requirements traced to shipped edits (gate + serial-only condition in work.md Step 1, stop-and-report bullet in Step 6, mirror rule in Step 5.5 + Scope template, schema line + verbatim-list sentence in work-reference.md, capture seeding guidance, model.go/generate.go lock-step with complete payload path, rejected timed-lock rationale in the gate), P-A-U boxes confirmed against the diff (no debug artifacts). Mechanical checks via `tools/checks/qualify.sh`: OK. Orchestrator note: Scope file list mirrored into `write_set:` frontmatter per the Step 5.5 rule this REQ ships — first application of its own contract.

## Testing

**Tests run:** `cd tools/queue-kanban && go test ./...`; `bash _dev/tests/contract-regressions.sh` (both re-run by the orchestrator, independent of the builder's runs)
**Result:** ✓ All passing (go: ok, 1 package; contract regressions: passed)

**Red-green validation:**
- `grep -ri write_set actions/ tools/queue-kanban/model.go`: ✗ no hits before implementation (exit 1, captured in Red-Green Proof) → ✓ hits in actions/board.md, actions/capture-reference.md, actions/work.md, actions/work-reference.md, model.go
- New ratchets (work.md gate text, serial-only rule, model.go parse call): builder mutation-tested each — removed/reworded the guarded text, all three assertions fired, restored, re-verified green
- `TestParseRequestTicketNormalizesAndResolves`: fixture had `write_set` unasserted → now asserts `[src/a.ts src/b.ts]` via `reflect.DeepEqual`, passing

**New tests added:**
- `model_test.go`: `write_set` fixture lines + DeepEqual assertion
- `_dev/tests/contract-regressions.sh`: three ratchet assertions (gate text, serial-only rule, `fields["write_set"]` parse call)

*Verified by work action*

## Review

**Overall: 89%** | 2026-07-28T21:35:03Z

**Verdict: Approve with follow-ups** — the write-set contract, dispatch gate, serial-only rule, and parser lock-step all landed exactly as planned with zero scope drift and verified red→green; two spec-coherence holes in the new parallel path need a follow-up before any harness actually turns concurrency on.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 70% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 2 important, 4 minor, 2 nits

### Findings

**Important**

1. **`actions/work.md:184` — the gate prescribes bookkeeping the lock schema cannot hold.** The bullet says "Every concurrently dispatched REQ still runs Steps 2–9 in full, including the orchestrator lock's `claimed_req` bookkeeping … Concurrency is in the dispatch, not in the bookkeeping." But `claimed_req` is a **single string per session** (`actions/work-reference.md:223` — "names whatever REQ currently sits in `do-work/working/` under the holder's claim, or `null` between REQs"). With N concurrent REQs: Step 2 (`actions/work.md:238`) rewrites `claimed_req` to *this* REQ's id, erasing the sibling's claim; Step 8 (`actions/work.md:548`) clears it to `null` after its archive move while siblings are still building; and the heartbeat rule (`actions/work-reference.md:357`) — "`claimed_req` = whatever's currently in `do-work/working/` under this session's claim" — is undefined for more than one file. Downstream, Crash Recovery's per-file gate (`actions/work-reference.md:203`) and `actions/cleanup.md`'s live-claim gate exempt **only** the REQ named in `claimed_req`, so every other concurrently dispatched REQ reads as unclaimed and gets stripped-and-re-queued or swept-and-archived mid-build — the exact 2026-07-01 failure class the lock guard exists to prevent, and most likely precisely during the long verification phase that motivated UR-007. Fix is either a list-valued claim (or one `coexisting_sessions` entry per builder) or an explicit stated limitation in the gate.

2. **`actions/work.md:336` — the Step 5.5 mirror can silently widen a co-dispatched write-set with no conflict check.** The mirror instruction says to copy the Scope list into `write_set` "replacing whatever capture seeded" — but the Step 1 co-dispatch decision was made on the *capture-seeded* sets, which `actions/capture-reference.md:113` itself calls "a hint, never a commitment." Concretely: REQ-A seeds `[src/a.ts]`, REQ-B seeds `[src/b.ts]` → co-dispatched; each one's `## Scope` then declares `src/shared.ts` as well; both mirror, both write it, and no rule fires. The Step 6 bullet (`actions/work.md:395`) requires exactly the missing check for a *mid-build* extension ("extends both the Scope list and `write_set` only after confirming no concurrently dispatched REQ claims that file") — so the same field has two write paths, one guarded and one not. Same root cause erases the partition directive: `actions/work.md:182` hands a builder a narrowed declared set, and Step 5.5 then overwrites `write_set` from Scope. One clause fixes both: on a concurrent dispatch, the Step 5.5 mirror runs the Step 6 no-concurrent-claimant check and serializes/partitions when it fails.

**Minor**

3. **`actions/work-reference.md:206` — Crash Recovery deletes the mirror's source but keeps the mirror.** Recovery strips `## Scope` and removes `claimed_at`/`route`, but leaves `write_set` in place. A recovered REQ then carries a mirror with no source (contradicting "Scope is the source … never the reverse") and feeds the dispatch gate a set from an abandoned attempt. Clearing it would fall back to the safe default (absent ⇒ overlaps everything ⇒ serialize).

4. **`actions/board.md:115` — the enumeration insert attaches `write_set` to a blocked-badge predicate.** "`domain`, `write_set`, and the blocked fields … are parsed for display only — they drive the 'blocked by' badge and drawer rows on a `status: blocked` card" now asserts something untrue of `write_set` (and, pre-existing, of `domain`). The next sentence corrects it, so the damage is small, but a cold reader hits the wrong claim first.

5. **`actions/capture.md:213` — the complex-REQ field enumeration wasn't extended.** It still reads "Set `related` and `batch` fields across the batch; populate `depends_on` per …", so the only prompt-level trigger for seeding `write_set` lives in the reference file. `actions/capture.md` was deliberately out of scope (deferred to REQ-034), but this is the Closed-Enumerations-Go-Stale pattern the project ratchets against — worth folding into REQ-034.

6. **`prime_files: []` while the edited Go tool has a co-located prime.** `tools/queue-kanban/prime-do-kanban.md` carries `## Traps` and `## Lessons` for exactly the files this REQ modified (`model.go`, `generate.go`), so `crew-members/general.md`'s Lessons Discipline never ran for the tool. No harm here — the edit is two lines and matches the file's conventions — but it repeats for the next board-touching REQ.

**Nit**

7. **`_dev/tests/contract-regressions.sh:130` — the `serial-only` ratchet is a bare-word grep.** Unlike D-03's `fields\["write_set"\]` anchor, it passes on any rewrite that keeps the hyphenated word (e.g. downgrading the rule to "serial-only classes are a future idea"). Anchoring on "never co-dispatched" or "global sequence" would bite harder. Mutation-verified as working today.

8. **`docs/work-guide.md:87` — the derived serial-only list has no illustrative marker.** "Migrations, lockfiles, and generated bundles" appears without the condition or a hedge. The canonical home (`actions/work.md:188`) states both, which is what the rule requires, so this is a gloss rather than a violation — a one-word "like" keeps it from reading as exhaustive.

**Positive:** exact 10/10 file match to the declared Scope (`tools/checks/scope-drift.sh`: OK); all three new ratchets independently mutation-verified by this review (rewording the gate phrase, deleting the serial-only paragraph, and deleting the parse line each fire their assertion, and deleting the parse line also fails `go test`); the REQ mirrored its own Scope into `write_set` — the contract's first self-application — and it renders end-to-end on the board.

### Requirements Checklist

- [x] Optional `write_set` frontmatter (repo-relative paths/globs), absence/empty ⇒ overlaps everything — delivered (`actions/work-reference.md:98`), stated in three places per D-01
- [x] Builder contract: write only inside the set; out-of-set need is stop-and-report, never a silent write; resolution lands in the REQ trail as a `D-XX` entry — delivered (`actions/work.md:395`)
- [x] Dispatch gate: concurrent dispatch iff pairwise-disjoint write-sets; overlaps serialize or take an explicit partition directive; floor agents unchanged — delivered (`actions/work.md:180-186`); see Important 1–2 for the coherence gaps
- [x] Serial-only resource classes stated as a *condition* with an explicitly illustrative list — delivered (`actions/work.md:188`)
- [x] Phrased as scheduling input, not a safety guarantee — delivered (schema line, gate bullet, `actions/work.md:683` Rules bullet)
- [x] One-directional Scope → `write_set` mirror so the two cannot drift — delivered (`actions/work.md:336`, `actions/work-reference.md:501`); see Important 2
- [x] `tools/queue-kanban/model.go` parses `write_set` in the same commit, display only, no column logic — delivered and verified end-to-end into `board-data.js`
- [x] Rejected timed-per-file-lock alternative recorded in the change — delivered (`actions/work.md:186`, citing the 0.140.4 mutex-break defect class)
- [x] RED→GREEN: `grep -ri write_set actions/ tools/queue-kanban/model.go` returns hits; `go test ./...` and `bash _dev/tests/contract-regressions.sh` green — verified independently
- [x] UR-007 batch constraints: no new action file, SKILL.md untouched (2648/2650 words), parser lock-step in the same commit, timed locks not reintroduced, Closed Enumerations honored in the canonical home — all met

### Acceptance Testing

**Result: Pass**
- `cd tools/queue-kanban && go test ./...` → ok (1 package)
- `bash _dev/tests/contract-regressions.sh` → "Contract regression checks passed."
- `gofmt -l .` empty; `go vet ./...` clean
- GREEN grep returns hits across `actions/board.md`, `actions/capture-reference.md`, `actions/work.md`, `actions/work-reference.md`, `tools/queue-kanban/model.go`
- Mutation testing on a scratch copy: rewording `pairwise disjoint`, deleting the `Serial-only resource classes` paragraph, and deleting the `fields["write_set"]` parse line fired all three new ratchets; the deleted parse line also failed `TestParseRequestTicketNormalizesAndResolves`
- End-to-end board render: built the tool and ran `summary` + `generate` against the live tree — 34 REQs, 0 completion anomalies, no parser/data warnings, and REQ-032's own 10-path `write_set` appears verbatim in `board-data.js` as `"writeSet":[…]` (null for the 33 REQs without one, matching `related`'s existing shape)
- `bash tools/checks/scope-drift.sh` on the REQ → "OK: Implementation Summary matches the Scope declaration"

### Suggested Additional Testing

- Walk the concurrent path by hand once REQ-033 lands: two co-dispatched REQs, then inspect `do-work/orchestrator-lock.json` after the first reaches Step 8 — confirm whether the second is still visible as claimed (Important 1)
- Dry-run the Step 5.5 mirror against two REQs whose Scope sections both name a shared file, and check whether anything stops the second mirror (Important 2)
- Have a cold agent read only `actions/work.md` Step 1 and answer "may I dispatch these two REQs together?" for three cases (both sets present and disjoint, one set absent, both touching migrations) — the gate's readability is the deliverable here
- Verify a consumer install picks up the field: run `do-work update` into a scratch consumer and confirm `queue-kanban` there parses a hand-added `write_set` without warnings
- Confirm the Commit Phase bumps `actions/version.md` and adds the `CHANGELOG.md` entry with the rejected-timed-lock bullet (planned as orchestrator work, not present in this diff)

**Acceptance:** Pass — full suite, linters, mutation-verified ratchets, and an end-to-end board render all green; the field reaches the rendered payload.
**Suggested testing:** 5 items
**Follow-ups created:** None by this review — Important 1 and Important 2 are handed to the orchestrator for queuing (both `addendum_to: REQ-032`, domain `general`; Important 1 likely wants `depends_on: [REQ-033]` since the worktree mode is where concurrency actually turns on)

*Reviewed by review-work action*
**Adversarial verification (orchestrator addendum):** 13 serious findings (2 review + 11 contradiction-hunt) were each attacked by 2 independent refuters. Confirmed: the two Important findings above plus three Critical restatements of the same two root causes — (A) single-slot `claimed_req`/Crash-Recovery gate cannot represent N concurrent claims by one orchestrator; (B) dispatch-time disjointness is decided on capture seeds that Step 5.5 overwrites unvalidated. Refuted (killed): Step 10 exit racing live builders (cleanup's live-claim gate + status filters hold), shared-tree diff-gate corruption (pre-existing, REQ-033's declared surface), partition-vs-serial-only contradiction, recovery-keeps-`write_set` as anti-drift defeat (direction rule unviolated; kept as Minor 3), capture-hint wording, stop-and-report gaps. Split (1–1, judged Minor): absent-set builder-side gloss; Step 4 file-conflict non-blocker wording. Follow-ups: REQ-035 (cluster A), REQ-036 (cluster B, folds in the two split Minors).

## Lessons Learned

**What worked:** Line-anchored exploration before building — the explorer caught three stale-enumeration sites the plan missed (`actions/board.md` lock-step list, Scope-template note, work-guide) and corrected the Step 1 insertion point. Mutation-testing the new ratchets (delete/reword the guarded text, watch each assertion fire) proved they bite before shipping them.
**What didn't:** Writing new concurrency prose against the existing single-writer state machine from memory — the gate bullet promised `claimed_req` bookkeeping the lock schema literally cannot hold (single string). Concurrency text must be drafted with the lock schema and Crash Recovery gate open side-by-side; two review rounds of REQ-018 taught this and it still recurred.
**Worth knowing:** `write_set` is a scheduling input only — the dispatch gate is coherent today for the serial default; actually turning concurrency on waits for REQ-035 (multi-claim lock representation) and REQ-036 (dispatch-time re-validation of firmed sets). The parser lock-step is now ratcheted: `_dev/tests/contract-regressions.sh` fails if the gate text or `model.go`'s `write_set` parse disappears.

## Orientation

Now the queue schema can declare per-REQ write surfaces (`write_set`) and the work pipeline has a sanctioned parallel-dispatch option: advanced harnesses may co-dispatch dependency-ready REQs whose write-sets are pairwise disjoint, with serial-only resource classes and a builder stop-and-report write boundary. Lives in the work-pipeline docs (`actions/work.md` Step 1/5.5/6, `actions/work-reference.md` schema) and the queue-kanban parser/payload. **[MAP CHANGED]** — new scheduling contract in the work pipeline; two coherence follow-ups (REQ-035/036) gate real concurrent use. No prime files were listed; `tools/queue-kanban/prime-do-kanban.md` covers the two Go files touched (lesson link added at archive).
