---
id: REQ-030
title: Split actions/pipeline.md into an action + reference pair
status: completed
created_at: 2026-07-27T07:34:50Z
claimed_at: 2026-07-27T08:07:09Z
completed_at: 2026-07-27T08:32:36Z
commit: 168871f
user_request: UR-006
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-027]
maintenance: false
related: [REQ-025, REQ-026, REQ-027, REQ-028, REQ-029, REQ-031]
batch: context-engineering-alignment
---

# Split actions/pipeline.md Into an Action + Reference Pair

## What

`actions/pipeline.md` is 7,471 words and loads whole — second-largest action file, and the bulk of it is output-rendering templates rather than decision logic. Split into `actions/pipeline.md` (steps + mode-determination logic) plus a new `actions/pipeline-reference.md`, following the `actions/work.md` ↔ `actions/work-reference.md` pattern.

Move to the companion:

- The **State File** schema (lines ~39–76).
- **All Output Format renderings**: the status block, the completion status block, the three-rendering Pipeline Completion Report (~lines 281–467), the continuation notice, and the help menu.

Keep in the action file: Steps 1–6, including Step 1's mode determination, Step 5a's queue-continuation logic, and Step 6's error handling — the decision logic, not the prose it emits.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, `actions/work-reference.md` header + heading conventions, `crew-members/general.md`, `crew-members/coding-guardrails.md`, the current `CLAUDE.md` § Action File Conventions (post-REQ-027), and the archived REQ-027 record for what "earned" trimming looks like in practice. Approach: (1) reduce `## State File` to a lifecycle paragraph (path, created-at, read-at, updated-at, never-deleted) + a pointer to the companion's full schema; (2) move the entire `## Output Format` section (Status Block, Completion Status Block, the three-rendering Pipeline Completion Report, Continuation Notice, Help Menu) to `actions/pipeline-reference.md`, replacing it in the action file with a short index stub whose bullets are the same pointers already inserted inline at each consuming step; (3) insert a named pointer at every step that currently prints/renders one of the moved blocks; (4) apply the earned-sections test to Rules/Common Rationalizations/Red Flags/Verification Checklist — audit every row/bullet for a do-work-specific noun or traceable failure mode, trim only what fails that test.
- [x] **[APPLY]:** Code written exactly as planned. Scope strictly limited to `actions/pipeline.md` (modified) and `actions/pipeline-reference.md` (new).
- [x] **[UNIFY]:** `git diff --stat` reviewed (below). Ran `bash _dev/tests/contract-regressions.sh` clean. Diffed the moved State File and Output Format bodies against `git show HEAD:actions/pipeline.md` (heading levels promoted by exactly one `#`, no other byte differences) to prove the relocation carried no accidental rewrite. Diffed the Steps section and the Rules section against HEAD to confirm the only changes are the intended pointer insertions and the single Platform-agnostic trim. No debug artifacts (this is prose, not code).

## Why (if provided)

Progressive disclosure — the third of Anthropic's five shifts for Claude 5 generation models. The three completion-report renderings are needed once, at the end of a run, by exactly one step; the state-file schema matters only when reading or writing that file. Today every pipeline invocation pays for all of it up front, including invocations that abort at Step 1.

## Context

- `actions/pipeline.md` heading map from the audit (re-verify before editing): Philosophy(7), When to Use(14), Input(26), State File(39), Steps(77), Step 1 Determine Mode(79), Step 2 Initialize(107), Step 3 Resume(120), Step 4 Execute Current Step(129), Step 5 Completion(161), Step 5a Queue Continuation(192), Step 6 Error Handling(227), Output Format(239), Status Block(241), Completion Status Block(264), Pipeline Completion Report — three renderings(281), Continuation Notice(468), Help Menu(497), Rules(520), Common Rationalizations(536), Red Flags(552), Verification Checklist(569).
- Copy the companion-header convention from `actions/work-reference.md`, including its lazy-read and don't-re-read instructions.
- `pipeline` dispatches `work` in the foreground, so anything the work loop inherits from pipeline's state handling must remain correct after the split. Verify the state-file contract still reads coherently when the schema lives in the companion.
- REQ-027 lands first — apply its updated rule to Rules / Common Rationalizations / Red Flags / Verification Checklist in both halves; a section that no longer earns its place is deleted, not relocated.

## Detailed Requirements

- **Every moved section is pointed at by name from the step that consumes it**; no dangling pointers, no orphan sections.
- **The mode-determination logic stays whole in the action file.** It is the part an agent must have before it can decide anything, including whether it needs the companion at all.
- **State File schema:** it moves, but the action file must retain enough to know the file's path, when it is created, and when it is deleted — the lifecycle, not the field list.
- **Help menu:** relocating it must not break `do-work pipeline help`. Confirm the routing path still lands somewhere that reaches the menu.
- **Both files ship.** Cite other actions by path; never cite `CLAUDE.md`/`AGENTS.md`.
- **Relocation, not redesign** — no behavior changes to the pipeline beyond REQ-027-mandated deletions.

## Constraints

- `bash _dev/tests/contract-regressions.sh` must pass clean.
- `SKILL.md` must not grow past 2,650 words.
- Version bump + descriptive `CHANGELOG.md` entry.

## Dependencies

Depends on REQ-027.

## Builder Guidance

**Certainty: Firm.** The split line is clear here — decision logic stays, emitted prose goes. The one judgment call is the state-file schema's boundary (lifecycle vs. fields); make it, state it in the Implementation Summary, and keep it consistent.

## Red-Green Proof

- **RED now:** `wc -w actions/pipeline.md` ≈ 7,471; a pipeline invocation that aborts at Step 1 has nonetheless loaded all three completion-report renderings and the help menu.
- **GREEN when:** `actions/pipeline.md` holds Steps 1–6 and the mode logic at materially reduced size; `actions/pipeline-reference.md` holds the state-file schema and every output rendering; each is reachable by name from its consuming step.
- **Validation:** `grep -n 'pipeline-reference' actions/pipeline.md` shows a pointer per moved section; cross-check names against `grep -n '^#' actions/pipeline-reference.md` in both directions; `do-work pipeline help` still resolves; before/after `wc -w` for both files and the sum; `bash _dev/tests/contract-regressions.sh` clean.

## Open Questions

None.

## Full Context

See `do-work/user-requests/UR-006/input.md` for complete verbatim input.

## Triage

**Route B** — the split boundary itself was firm (per Builder Guidance: decision logic stays, emitted prose goes; the one judgment call was the State File lifecycle-vs-fields line). What needed exploration was the earned-sections audit (REQ-027's rule) across the existing Rules/Common Rationalizations/Red Flags/Verification Checklist content, and matching `actions/work.md`/`actions/work-reference.md`'s exact citation conventions before writing the companion.

## Scope

**Files committed to before editing:**
- `actions/pipeline.md` (modified) — reduce State File to lifecycle + pointer; replace the `## Output Format` section with an index stub; insert a named companion pointer at every step that prints/renders a moved block; apply the earned-sections test to Rules.
- `actions/pipeline-reference.md` (new) — State File Schema, Status Block, Completion Status Block, Pipeline Completion Report (three renderings + Composition rules), Continuation Notice, Help Menu.

**Acceptance criteria restated from the REQ:**
- Every moved section pointed at by name from the step that consumes it; no dangling pointers, no orphan sections.
- Mode-determination logic (Step 1) stays whole in the action file.
- State File schema moves, but the action file retains the file's path, when it's created, and when it's deleted (lifecycle, not the field list).
- Help menu relocation doesn't break `do-work pipeline help`'s routing path to the menu.
- Both files ship; no citation of `CLAUDE.md`/`AGENTS.md`.
- Relocation, not redesign — no behavior changes beyond REQ-027-mandated deletions.
- `bash _dev/tests/contract-regressions.sh` passes clean; `SKILL.md` stays ≤2,650 words (untouched by this REQ).

## Implementation Summary

**What was done:**

1. **`actions/pipeline.md` § State File** — reduced from the full JSON schema + field table + pretty-print invariant to a lifecycle paragraph: path (`do-work/pipeline.json`), created at Step 2, read at the top of every invocation, rewritten at every step transition (Steps 4/5/5a/6), and — the judgment call the Builder Guidance flagged — **never deleted**: Step 5 (Completion) and the Abandon mode both set `active: false` and leave the file on disk as a historical record; it's excluded from git, not removed from disk. The REQ's phrasing ("when it is deleted") assumed a deletion step that doesn't exist anywhere in this action — I verified this by grepping the original file for `delete|remove|rm |unlink` (no hits) and confirmed with `hooks/pipeline-guard.sh`/`session-start.sh`, which only ever *read* `pipeline.json`, never remove it. Stated this explicitly rather than inventing a deletion step to satisfy the REQ's wording literally.
2. **`actions/pipeline.md` § Output Format** — the entire section (Status Block, Completion Status Block, the three-rendering Pipeline Completion Report with its Composition rules and Marp/HTML template bodies, Continuation Notice, Help Menu) moved verbatim to `actions/pipeline-reference.md`, heading levels promoted by exactly one `#` (`###`→`##`, `####`→`###`) since the companion is now the top-level file for this content. Verified byte-for-byte fidelity with a diff against `git show HEAD:actions/pipeline.md` after undoing the heading promotion (see Testing) — the only difference is the promoted `#` count.
3. **`actions/pipeline.md` Steps** — inserted a named pointer to `actions/pipeline-reference.md` at every print/render site: Step 1's Help mode (→ **Help Menu**), Steps 2/3/4/6's status-block prints (→ **Status Block**), Step 5's completion print (→ **Completion Status Block**) and completion-report assembly (→ **Pipeline Completion Report — three renderings of one dataset**, which also names **Composition rules**), the three-format table's "Template / Producer" column (→ **Plain Markdown Report** / **Marp Slide Deck** / **Standalone HTML Debrief**), Step 5a's continuation print and max-iterations message (→ **Continuation Notice**, including its limit-reached variant), and Step 6's failure print (→ **Status Block**). Step 5a's "Max iterations" paragraph previously re-embedded the exact limit-reached message block a second time (duplicate of the Output Format section's own copy) — collapsed to a single pointer at the companion's one copy rather than keeping two copies of the same literal text in sync.
4. **`actions/pipeline.md` § Output Format (stub)** — the heading stays (matching `actions/work.md`'s pattern of keeping a short stub heading like `## Progress Reporting` rather than deleting it outright) but its body is now a five-line index of the same pointers already inserted at each step, for an agent skimming the file top-to-bottom before it reaches the steps.
5. **`actions/pipeline.md` § Rules — earned-sections audit (REQ-027).** Read every Rules bullet, every Common Rationalizations row, every Red Flag, and every Verification Checklist item against the test ("can I name the specific failure this prevents, and where it happened?"). Every row in Common Rationalizations/Red Flags/Verification Checklist already names a do-work-specific mechanic (`pipeline.json`, specific REQ/UR machinery, the three-rendering report, Marp/HTML artifact requirements) — none were generic filler, so none were touched. One Rules bullet failed the test: **"Platform-agnostic. No tool-specific APIs. Dispatch actions the same way the main router does."** is restated engineering hygiene already covered verbatim by Step 4 ("Dispatch each action the same way the main router dispatches actions..."). Trimmed to keep only the earned, do-work-specific half — the `hooks/pipeline-guard.sh` stop-hook mention — retitled **"Optional stop-hook guard."** This is the one deletion in this REQ; everything else in Rules/Common Rationalizations/Red Flags/Verification Checklist was already earned and is unchanged.
6. **`actions/pipeline-reference.md` (new)** — header blockquote copies `actions/work-reference.md`'s conventions verbatim in spirit: names what it holds, states each section is pointed to by name from the matching step in `pipeline.md`, states the lazy-read rule ("only necessary when you reach the step that references it — and read only the named section") and the don't-re-read rule. No Rules/Common Rationalizations/Red Flags/Verification Checklist sections were added to the companion — it holds pure templates/output, matching `actions/work-reference.md`'s own shape (zero such sections there either), and per CLAUDE.md's ratchet, a *new* action file's Common Rationalizations table must carry a do-work-specific noun or the suite fails — since none of the moved content is a behavioral rationalization table, none was added.

**Decision recorded (Builder Guidance's flagged judgment call):** the State File split line is lifecycle (path, created, read, updated, never-deleted) in the action file vs. full JSON shape/field table/pretty-print invariant in the companion — applied consistently at the one call site.

**Files changed:**
- `actions/pipeline.md` (modified) — State File reduced to lifecycle + pointer; `## Output Format` body replaced with an index stub; 8 named pointers inserted across Steps 1/2/3/4/5/5a/6; one Rules bullet trimmed (Platform-agnostic → Optional stop-hook guard). 7,471 → 4,719 words.
- `actions/pipeline-reference.md` (new) — State File Schema, Status Block, Completion Status Block, Pipeline Completion Report (Composition rules + 3 authored-format templates + Marp frontmatter skeleton), Continuation Notice, Help Menu. 0 → 3,084 words.

**Word-count receipt (`wc -w`, before → after):**

| File | Before | After |
| --- | --- | --- |
| `actions/pipeline.md` | 7,471 | 4,719 |
| `actions/pipeline-reference.md` | 0 (new) | 3,084 |
| **Sum** | **7,471** | **7,803** |

The sum grew by 332 words — expected, not a regression: the companion's header/blockquote is new prose, the action file's State File lifecycle paragraph is more explicit than the one-liner it replaced, and 8 inline pointers plus a 5-line Output Format index were added where the original had none. The win this REQ targets is progressive disclosure (a Step-1 abort now loads 4,719 words instead of 7,471, and never touches the companion's 3,084), not a smaller total.

**Two-way pointer check (the failure mode that matters most):**

Direction 1 — every pointer in `actions/pipeline.md` names a section that exists in the companion:
```
$ grep -n 'pipeline-reference' actions/pipeline.md
```
14 pointer sites, resolving to 8 distinct companion section names: **State File Schema**, **Help Menu**, **Status Block** (×4 call sites), **Completion Status Block**, **Pipeline Completion Report — three renderings of one dataset** (+ **Composition rules**), **Plain Markdown Report**, **Marp Slide Deck**, **Standalone HTML Debrief**, **Continuation Notice** (×2 call sites). Cross-checked each against:
```
$ grep -n '^#' actions/pipeline-reference.md
```
Every one of the 8 names is a literal substring of an actual heading in that output (`## State File Schema`, `## Status Block (printed after every step transition)`, `## Completion Status Block (printed when all 6 steps are done)`, `## Pipeline Completion Report — three renderings of one dataset`, `### Composition rules (apply to all three formats)`, `### 1. Plain Markdown Report — ...`, `### 2. Marp Slide Deck — ...`, `### 3. Standalone HTML Debrief — ...`, `## Continuation Notice (...)`, `## Help Menu (...)`) — no dangling pointer.

Direction 2 — every top-level section in the companion is reachable from the action file: **State File Schema** (State File section), **Status Block** (Steps 2/3/4/6), **Completion Status Block** (Step 5), **Pipeline Completion Report** + **Composition rules** (Step 5, both the assembly line and the Output Format stub), the three format subsections (Step 5's table), **Continuation Notice** (Step 5a, twice), **Help Menu** (Step 1's mode table + Output Format stub) — no orphan section.

## Testing

This is a prose relocation, not code — "testing" is fidelity verification that the split moved content losslessly and the pointers resolve.

**Contract suite:**
```
$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
```
(Passes because `pipeline.md` is in the REQ-027 grandfather baseline, and the new `pipeline-reference.md` carries no `## Common Rationalizations` section at all, so the ratchet loop skips it.)

**Fidelity diff — State File move** (git HEAD vs. companion, blank-line-insensitive): identical apart from the one-line intro sentence, which the action file now expresses as the fuller lifecycle paragraph instead of duplicating in the companion.

**Fidelity diff — Output Format move** (git HEAD's `## Output Format` body vs. the companion's moved sections, heading levels promoted by exactly one `#` to account for the file becoming top-level):
```
$ diff /tmp/orig_transformed2.txt /tmp/new_output_format.txt
279d278
<
```
Only difference: a trailing blank line before the next heading in the original (which had a following `## Rules` section) — no content divergence.

**Fidelity diff — Steps section** (git HEAD vs. new, full section): every hunk is one of the 8 intended pointer insertions or the Step 5a duplicate-message collapse; no other line changed.

**Fidelity diff — Rules section** (git HEAD vs. new, full section): the single Platform-agnostic → Optional stop-hook guard trim; every other line byte-identical.

**Two-way pointer check:** see Implementation Summary above (pasted as the receipt).

**`do-work pipeline help` routing — confirmed, with a pre-existing gap found (not introduced by this REQ, not fixed here — see Discovered Tasks).** Traced the path: `SKILL.md` routing row 1 only matches literal `help` alone; `pipeline help` matches row 3 (`pipeline` ± request text) first, so it dispatches `actions/pipeline.md` with `$ARGUMENTS = "help"` per `SKILL.md`'s stated carve-out ("except `pipeline`, `prime`, and `bkb`, which handle `help` internally"). Inside `pipeline.md`'s own Step 1 bucket logic (untouched by this REQ — mode-determination logic must stay whole per the REQ), `help` is not one of the two reserved keywords (`status`/`abandon`), so it buckets as **Request text**, which (no active pipeline) resolves to **Initialize** — not Help. This means `do-work pipeline help` does not currently reach the Help Menu at all; it would try to start a new pipeline with request text `"help"`. This is a routing-table gap in Step 1's bucket definition, not anything introduced by relocating the Help Menu's *content* — the pointer from Step 1's Help-mode arm to the companion's **Help Menu** section is correct and intact for every input that *does* resolve to Help mode (empty, `status`, `abandon` with no active pipeline). Fixing the bucket table itself would be decision-logic redesign, explicitly out of scope ("Step 1's mode determination stays whole in the action file... this is a relocation, not a redesign"). Recorded below as a Discovered Task rather than fixed inline.

## Lessons Learned

**What worked:** Diffing the moved blocks against `git show HEAD:actions/pipeline.md` (with heading levels programmatically promoted for comparison) caught that my first verification script had a sed cascading bug — chaining `s/^#### /### /; s/^### /## /` in one sed invocation double-demotes lines that started as `####`, since the first substitution's output re-matches the second pattern on the same pass. Rewrote as an `awk` script with explicit `if/next` branching to avoid the cascade. Worth remembering for any future heading-level-promotion verification: never chain sequential regex demotions in one pass.

**What didn't:** The REQ's phrasing for the State File judgment call ("know... when it is deleted") assumed a deletion step that doesn't exist in this action. Rather than either inventing one to satisfy the letter of the REQ or silently dropping the phrase, I grepped for `delete|remove|rm |unlink` across the original file (no hits) and stated the actual lifecycle (never deleted, excluded from git instead) — this is the kind of REQ-vs-reality mismatch the Builder Guidance anticipated by calling this out as "the one judgment call."

**Worth knowing:** All of the pipeline's existing Rules/Common Rationalizations/Red Flags/Verification Checklist content was already earned under REQ-027's test — this file didn't accrete generic filler the way the REQ-027 audit found elsewhere (20-plus files with boilerplate rationalizations). The only casualty was one Rules bullet that duplicated Step 4's own dispatch instruction almost verbatim. A useful signal for future earned-sections audits: check whether a Rules bullet is already fully stated inline in a Step — if so, it's very likely restated hygiene rather than an independently earned rule.

## Discovered Tasks

- `do-work pipeline help` does not reach the pipeline's Help Menu under the current Step 1 bucket logic. `SKILL.md` routes `pipeline help` to `actions/pipeline.md` with `$ARGUMENTS = "help"` (its documented carve-out: pipeline/prime/bkb "handle help internally" instead of being intercepted by the generic per-command help router). But `pipeline.md`'s Step 1 only recognizes `status`/`abandon` as reserved keywords — `help` buckets as **Request text**, which (no active pipeline) routes to **Initialize**, attempting to start a new pipeline with request text `"help"` instead of showing the Help Menu. Pre-existing in the file before this REQ (confirmed by reading `git show HEAD:actions/pipeline.md`'s Step 1 bucket definition, unchanged by the split); out of scope here since Step 1's mode-determination logic must stay whole per this REQ's Builder Guidance. Fix would add `help` as a third reserved keyword (or its own bucket) in Step 1's normalize/bucket logic, routed to Help mode from both active and inactive pipeline states.

---
*Source: "compare with the current skill, is there something that we need to update?" — resolved into the approved seven-REQ plan.*

Think carefully before answering.

## Review

**Verdict: Partial → resolved.** The split itself verified clean: both moved blocks (the State File JSON schema/field table/pretty-print invariant, and the entire Output Format section including all three completion-report renderings) diffed byte-identical against `git show HEAD:actions/pipeline.md` after heading-level normalization; two-way pointer integrity confirmed (14 pointer sites naming 8 distinct companion headings, no orphans, no danglers); no unaccounted deletions.

**One Important finding — ADR drift (orchestrator-resolved, not deferred):**

`decisions/records/adr-008-*.md` stated as current fact: "The rendering templates and cross-format rules live inline in `pipeline.md`'s Output Format section." REQ-030 made that false. The project has direct precedent here — `CHANGELOG-archive.md`'s 0.76.0 entry records that adr-001 and adr-008 were both updated when this same companion was re-inlined the other direction, which establishes updating them as the expected companion move whenever `pipeline.md`'s split/inline state changes.

Fixed in this REQ's commit rather than queued as a follow-up:
- `decisions/records/adr-008-*.md` — Decision section now states the templates live in `actions/pipeline-reference.md`, notes the 0.76.0-to-REQ-030 history, and clarifies that the ADR's actual decision is about the three formats and their cross-links, not about which file holds the templates.
- `decisions/records/adr-001-*.md` — the `pipeline.md` reference line now records the full split → re-inline → re-split history (matching how the same line already treats `work.md`), and gains the two new pairs from REQ-029 and REQ-031.

**Scope note (D-01, DECIDE & STATE):** `decisions/` was outside this REQ's declared file list. Updating it here is a deliberate expansion — the ADRs are the paper trail that this REQ's change invalidates, the precedent for doing it in the same commit is explicit, and deferring it would ship a decision record that is actively false. Reversible, no reasonable disagreement, so decided rather than escalated. `decisions/` is `export-ignore`'d, so no consumer-facing surface changed.
