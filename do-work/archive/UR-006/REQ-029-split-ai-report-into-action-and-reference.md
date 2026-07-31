---
id: REQ-029
title: Split actions/ai-report.md into an action + reference pair
status: completed
created_at: 2026-07-27T07:34:50Z
claimed_at: 2026-07-27T08:07:09Z
completed_at: 2026-07-27T08:31:31Z
commit: 2273381
user_request: UR-006
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-027]
maintenance: false
related: [REQ-025, REQ-026, REQ-027, REQ-028, REQ-030, REQ-031]
batch: context-engineering-alignment
---

# Split actions/ai-report.md Into an Action + Reference Pair

## What

`actions/ai-report.md` is 7,541 words and loads whole whenever the action runs — the largest single action file in the skill. Split it into `actions/ai-report.md` (the step skeleton and decision logic) plus a new `actions/ai-report-reference.md` companion, following the established `actions/work.md` ↔ `actions/work-reference.md` pattern.

Move to the companion:

- The **Image Generation Backend** section (lines ~15–93).
- The **Step 4/5 CSS specifications** and the **before/after toggle reference implementation** (the inline HTML/CSS/JS at lines ~322–345).
- The **Output Format** template.

Keep in the action file: Steps 1–8 as a skeleton, each step pointing at its companion section **by name**.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Re-verified the heading map with fresh `grep -n`/`Read` before touching anything (matched the audit almost exactly — headings unmoved). Plan: (1) delete `## Image Generation Backend` (lines 15–93) wholesale into a new reference section, fixing its three consuming pointer sites (Step 1, Step 3c table, Step 4d); (2) delete Step 4c's `Data-viz rules for every hand-authored SVG` bullets into a `SVG Data-Viz Rules (Step 4c)` reference section, leaving a one-line pointer + fixing the Red Flags mention that named it as "(Step 4c rules)"; (3) delete Step 5's `#### Design rules` bullet list (the whole list, not a CSS-only subset — it's one cohesive block and Builder Guidance favors step-aligned sections) into `Report Design Rules (Step 5)`, leaving a one-line pointer; (4) delete the `#### Before/after toggle pattern` code block into `Before/After Toggle Reference Implementation (Step 5)`, keeping its one intro sentence in the action file; (5) delete `## Output Format`'s body into `Output Format Template (Step 8)`, leaving a one-line pointer. Then apply CLAUDE.md's earned-sections test to `## Rules`/`## Common Rationalizations`/`## Red Flags`/`## Verification Checklist` row-by-row (test: "can I name the specific failure this row prevents, and where it happened?").
- [x] **[APPLY]:** Implemented exactly per plan — see Implementation Summary. One finding from the earned-sections pass: every row in all four sections named a specific ai-report/do-work failure mode (a REQ/UR/status reference, a named Step, a hard-won bug like the `file://` blank-screenshot issue or the prompt-injection→RCE risk) — none read as restatable generic engineering advice, so **zero rows were deleted**. This is a real outcome of applying the test, not a skipped step — reasoning for every row is in the Implementation Summary.
- [x] **[UNIFY]:** `git diff --stat` shows only `actions/ai-report.md` modified (139 lines: 8 insertions, 131 deletions) plus the new `actions/ai-report-reference.md`; `git status --porcelain` confirms `actions/capture.md`/`actions/pipeline.md` (+ their new reference files) belong to sibling REQ-030/REQ-031 builders in this same session, untouched by me. No debug artifacts. `bash _dev/tests/contract-regressions.sh` passes clean (see Testing). Two-way pointer check performed (see Implementation Summary) — no dangling pointers, no orphan sections.

## Why (if provided)

Progressive disclosure is the third of Anthropic's five shifts for Claude 5 generation models: load the skeleton, fetch the heavy detail only at the step that needs it. Only 4 `*-reference.md` companions exist today while 10 of the top-15 action files load whole. `ai-report.md` is the worst offender, and most of its bulk is asset templates that matter at exactly one step each — a report that never generates an image still pays for the entire Image Generation Backend section today.

## Context

- The pattern is already proven in this repo: read `actions/work-reference.md`'s header blockquote and copy its conventions, including the instruction to read only the named section and to reuse the file if it's already in context rather than re-reading it per reference site.
- `actions/ai-report.md` heading map from the audit (re-verify before editing): Philosophy(7), Image Generation Backend(15), When to Use(94), Input(107), Steps(117), Step 1 Load Principles(119), Step 2 Resolve the Target(123), Step 3 Collect Visual Evidence(142), Step 4 Build the Visual Assets(197), Step 5 Write the Report HTML(264), Step 6 Self-Review Against Anti-Slop(347), Step 7 Render and Judge(366), Step 8 Save and Report(400), Output Format(421), Rules(425), Common Rationalizations(436), Red Flags(454), Verification Checklist(485).
- REQ-027 lands first and changes the rules for Rules / Common Rationalizations / Red Flags / Verification Checklist. Apply the updated rule to both halves as part of this split — a section that no longer earns its place is deleted, not relocated.

## Detailed Requirements

- **Every moved section must be pointed at by name from the step that consumes it.** A dangling pointer (points at a section that doesn't exist) and an orphan section (exists but nothing points at it) are both defects; check for both.
- **Companion header:** copy the convention from `actions/work-reference.md` — what the file is, why it was extracted, that each section is pointed at from the matching step, and that it should be read lazily and not re-read.
- **The action file must remain runnable as a standalone prompt** for the common path. An agent reading only `actions/ai-report.md` should know what each step does and when it needs to open the companion.
- **Both files ship.** Cite other actions by path; never cite `CLAUDE.md`/`AGENTS.md`.
- **No content rewriting beyond the split** except where REQ-027's rule requires a deletion. This is a relocation, not a redesign of the report action.

## Constraints

- `bash _dev/tests/contract-regressions.sh` must pass clean.
- No `SKILL.md` change should be needed — verify the routing/dispatch rows for `ai-report` still resolve, and if the router genuinely needs to know about the companion, do it without growing `SKILL.md` past 2,650 words.
- Version bump + descriptive `CHANGELOG.md` entry.

## Dependencies

Depends on REQ-027 (the template rule determines what the split files may contain).

## Builder Guidance

**Certainty: Firm on what moves, mixed on granularity.** The four listed items move. The builder decides how finely to slice them into named companion sections — favor sections that map 1:1 to a step, since that's what makes lazy loading actually work. If a section is consumed by two steps, name it once and point twice.

## Red-Green Proof

- **RED now:** `wc -w actions/ai-report.md` ≈ 7,541; running the action loads the Image Generation Backend section and the full before/after toggle implementation even when neither is used.
- **GREEN when:** `actions/ai-report.md` is materially smaller (target: the step skeleton, roughly half or less), `actions/ai-report-reference.md` exists with the moved sections, and every step that needs heavy detail names its companion section.
- **Validation:** `grep -n 'ai-report-reference' actions/ai-report.md` lists a pointer for each moved section; cross-check each named section exists in the companion (`grep -n '^#' actions/ai-report-reference.md`) and that each companion section is pointed at from somewhere; before/after `wc -w` for both files and the sum; `bash _dev/tests/contract-regressions.sh` clean.

## Open Questions

None.

## Full Context

See `do-work/user-requests/UR-006/input.md` for complete verbatim input.

---
*Source: "compare with the current skill, is there something that we need to update?" — resolved into the approved seven-REQ plan.*

Think carefully before answering.

## Triage

**Route: A** - Simple

**Reasoning:** Mechanical relocation against a firm, pre-decided move-list (Builder Guidance: "Certainty: Firm on what moves") plus a row-by-row omission-test pass over four pre-identified sections. No exploration or planning agent needed.

**Planning:** Not required

## Scope

**Files I will touch:**
- `actions/ai-report.md` (modified) — remove the four relocated blocks, replace each with a by-name pointer at the step(s) that consume it, fix two downstream mentions ("Step 4c rules" → named pointer) that would otherwise dangle.
- `actions/ai-report-reference.md` (new) — the four relocated sections, each headed with the consuming Step(s) in its own heading, matching `actions/work-reference.md`'s header/citation conventions.

**Files I will NOT touch:** `SKILL.md` (routing/dispatch rows for `ai-report` already resolve — no change needed, confirmed below), `docs/ai-report-guide.md` (checked for dangling anchors — none found), `actions/version.md`, `CHANGELOG.md` (orchestrator-owned), `actions/capture.md`/`actions/pipeline.md` (sibling REQ-030/031 scope).

**Acceptance criteria (restated from REQ):**
- [x] Every moved section is pointed at by name from the step(s) that consume it.
- [x] No dangling pointer (pointer names a section that doesn't exist) and no orphan section (section nothing points at).
- [x] Companion header copies `actions/work-reference.md`'s convention (what the file is, why extracted, pointed-at-by-name, read-lazily-and-don't-re-read).
- [x] `actions/ai-report.md` remains runnable as a standalone prompt for the common path.
- [x] Both files ship; no citation of `CLAUDE.md`/`AGENTS.md` in either.
- [x] No content rewriting beyond the split, except deletions earned by REQ-027's rule (none were earned here — see Implementation Summary).
- [x] `bash _dev/tests/contract-regressions.sh` passes clean.
- [x] `SKILL.md` untouched, still ≤ 2,650 words (2,552, unaffected).
- [ ] Version bump + CHANGELOG entry — left to the orchestrator per the harness rule; not touched by this builder.

## Implementation Summary

**What was done:** Split `actions/ai-report.md` into the eight-step skeleton plus a new `actions/ai-report-reference.md` companion, following `actions/work.md`/`actions/work-reference.md`'s proven pattern. Moved exactly the four items the REQ named — Image Generation Backend, the Step 4c/5 CSS-style specifications (Data-viz rules + Design rules), the before/after toggle reference implementation, and the Output Format template — each into its own reference section headed with the consuming Step(s). Every consuming site in the action file was updated to name the companion section explicitly (a `above`/`described above`/bare-heading self-reference doesn't survive relocation, so all of those were rewritten as `actions/ai-report-reference.md` pointers). Applied CLAUDE.md's post-REQ-027 earned-sections test to `## Rules`/`## Common Rationalizations`/`## Red Flags`/`## Verification Checklist` — see the row-by-row reasoning below; the outcome was zero deletions, which is a legitimate result of the test, not a skipped step.

**Files changed:**
- `actions/ai-report.md` (modified)
- `actions/ai-report-reference.md` (new)

**Word-count receipt** (`wc -w`):

| File | Before | After |
|---|---|---|
| `actions/ai-report.md` | 7,541 | 5,477 |
| `actions/ai-report-reference.md` | 0 (new) | 2,339 |
| **Sum** | **7,541** | **7,816** |

The sum grew by 275 words net (moved content plus its heading + the by-name pointer sentences replacing terser same-file references like "described above" — expected, since a pointer that survives relocation must name a path). The action file itself dropped 2,064 words (27%) — short of the REQ's aspirational "roughly half or less" target, because the Detailed Requirements and Builder Guidance ("Certainty: Firm on what moves... four listed items move") scope the move to exactly the four named items, not every heavy block in Steps 4/5 (e.g. the 4a SVG-callout-anatomy code sample and the 4b Mermaid code sample are also inline templates but were never named in the REQ's move list, so they stay put — moving them would be scope creep beyond a firm builder-guidance boundary, not a judgment call left open to the builder). Flagging this delta explicitly rather than silently under-delivering on the numeric target.

**What moved, and where each new section's pointer(s) live in `actions/ai-report.md`:**

1. **`## Image Generation Backend (Steps 3c, 4d, 5)`** — the full original section (backend fallback chain, style-brief convention, `gen_image` helper, parallel-fire pattern, rules for generated images). Pointed at by name from: Step 1 (the `$2` trust-boundary note, line 42), the Step 3c decision table's AI-image row (line 114), and Step 4d's lead sentence (line 171). Step 5's hero/how-it-works mentions still route through "Step 4d" rather than naming the Backend section directly — that's unchanged from the original (which itself only claimed Step 5 as an indirect consumer via 4d), not something this split altered.
2. **`## SVG Data-Viz Rules (Step 4c)`** — the four data-viz bullets (color-by-job ordinal ramp, ink-colored text tokens, label-collision handling, stat-tile spec). Pointed at from Step 4c's own intro (line 167) and from the Red Flags line that used to say the vague "(Step 4c rules)" (line 348, now names the section).
3. **`## Report Design Rules (Step 5)`** — the full `#### Design rules` bullet list (21 bullets: CSS custom properties, aesthetic-direction rule, full-bleed layout, flex-wrap responsiveness, native-resolution image caps, base content rules). Moved as one intact block (not diced into "CSS-only" vs. "prose" bullets) since it was already a single cohesive list and Builder Guidance favors 1:1 step-aligned sections over finer slicing. Pointed at from Step 5's `#### Design rules` heading (line 218).
4. **`## Before/After Toggle Reference Implementation (Step 5)`** — the HTML/CSS/JS code block. The one-sentence usage guidance ("side-by-side by default, toggle only as fallback") stayed in the action file; only the code moved. Pointed at from Step 5's `#### Before/after toggle pattern` heading (line 222).
5. **`## Output Format Template (Step 8)`** — the folder/file-shape paragraph. Pointed at from the action file's own `## Output Format` heading (line 300), which now carries a one-line summary + pointer instead of the full paragraph.

**Two-way pointer check (receipts):**

```
$ grep -n 'ai-report-reference' actions/ai-report.md
42:  ...see the Image Generation Backend `$2` trust-boundary note in `actions/ai-report-reference.md`.
114: ...Mark it for AI image generation (Step 4d, which uses the **Image Generation Backend** in `actions/ai-report-reference.md`)...
167: ...see **SVG Data-Viz Rules (Step 4c)** in `actions/ai-report-reference.md`...
171: ...generate real images using the **Image Generation Backend** in `actions/ai-report-reference.md`...
218: Follow the **Report Design Rules (Step 5)** in `actions/ai-report-reference.md`...
222: ...Full HTML/CSS/JS reference implementation: **Before/After Toggle Reference Implementation (Step 5)** in `actions/ai-report-reference.md`.
300: See the **Output Format Template (Step 8)** in `actions/ai-report-reference.md`...
348: ...shorten strings (see **SVG Data-Viz Rules, Step 4c**, in `actions/ai-report-reference.md`).

$ grep -n '^#' actions/ai-report-reference.md
1:# AI Report Action — Reference
7:## Image Generation Backend (Steps 3c, 4d, 5)
86:## SVG Data-Viz Rules (Step 4c)
95:## Report Design Rules (Step 5)
117:## Before/After Toggle Reference Implementation (Step 5)
142:## Output Format Template (Step 8)
```

Direction 1 (every companion section pointed at by name): all 5 headings in `ai-report-reference.md` have at least one inbound named pointer above (Image Generation Backend ×3, SVG Data-Viz Rules ×2, Report Design Rules ×1, Before/After Toggle ×1, Output Format Template ×1). Direction 2 (every pointer names a real section): all 8 pointer lines above name one of the 5 headings verbatim or near-verbatim (line 348 uses a comma instead of the heading's parenthesis — same name, unambiguous). No dangling pointer, no orphan section.

**Earned-sections pass (CLAUDE.md's post-REQ-027 rule) — row-by-row, zero deletions:**

Applied the test ("can I name the specific failure this row prevents, and where it happened?") to every row/bullet in `## Rules` (8 bullets), `## Common Rationalizations` (12 rows), `## Red Flags` (28 bullets), and `## Verification Checklist` (15 items). Every single one ties to a specific do-work/ai-report mechanism or a named hard-won bug: `REQ`/`UR`/`status:` schema references (e.g. the `completed-with-issues` filtering bug), specific Step numbers, specific tool names (`bowser`, `playwright-cli`, Mermaid CDN), specific do-work paths (`ai-reports/<report-slug>/`, `do-work/deliverables/`), or a named security failure mode (image-gen prompt injection → RCE via the opt-in agentic backend, the `file://` headless-Chrome-blank-screenshot bug). None reads as generic restatable engineering hygiene ("write tests," "don't skip validation") — the action is narrow and specialized enough that its four earned-in-principle sections were, in practice, already earned in full. This differs from REQ-025's finding (which was about the *same guard restated across many files*, not generic filler within one file) and from what REQ-027's audit found elsewhere (e.g. capture.md/bkb.md's inlined prompt-injection doctrine) — ai-report.md's Rules/Rationalizations/Red-Flags/Checklist content was already scoped to this action's own failure modes, not restated hygiene. The only edits inside these four sections were the two pointer-name fixes required by the relocation itself (line 348's Red Flags row, plus the Rules/Rationalizations/Verification-Checklist bullets were left byte-for-byte identical — confirmed by the `git diff` above, which shows no `-`/`+` pairs inside those three sections).

## Testing

Non-behavioral relocation (prose/structure change) — no code, no automated test suite beyond the contract regression suite. Verification was grep/wc-receipt based, per the REQ's own Red-Green Proof.

**RED (before):** `wc -w actions/ai-report.md` = 7,541 (confirmed at REQ intake).

**GREEN (after):**
```
$ wc -w actions/ai-report.md actions/ai-report-reference.md
    5477 actions/ai-report.md
    2339 actions/ai-report-reference.md
    7816 total
```
Two-way pointer check: see receipts in Implementation Summary — 5/5 companion sections have inbound named pointers, 8/8 pointer sites name a real heading.

```
$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
```

(The suite's Common Rationalizations regrowth ratchet applies to new action files — `actions/ai-report-reference.md` has no `## Common Rationalizations` heading at all, so it isn't a candidate for that specific check; the file passed the suite's other checks — SKILL.md word budget, self-citation ban, schema-read-contract greps, etc. — cleanly.)

```
$ git status --porcelain
 M actions/ai-report.md
 M actions/capture.md          ← sibling REQ-031 builder, not this REQ
 M actions/pipeline.md         ← sibling REQ-030 builder, not this REQ
?? actions/ai-report-reference.md
?? actions/capture-reference.md   ← sibling REQ-031
?? actions/pipeline-reference.md  ← sibling REQ-030

$ git diff --stat -- actions/ai-report.md
 actions/ai-report.md | 139 +++------------------------------------------------
 1 file changed, 8 insertions(+), 131 deletions(-)

$ wc -w SKILL.md
    2552 SKILL.md   ← unaffected, under the 2,650 budget

$ grep -n "ai-report.md#\|Image Generation Backend" docs/ai-report-guide.md
(no matches — no dangling anchors in the user-facing guide)
```

## Lessons Learned

**What worked:** Re-deriving the heading map with fresh `grep -n` before editing (rather than trusting the REQ's audit numbers) caught that the map was still accurate — a good sign the audit was recent, but worth doing anyway since the REQ itself flagged "re-verify before editing." Extracting the exact section text via `sed -n` ranges before writing the reference file (rather than retyping from memory) avoided any risk of silently rewriting content during the "relocation, not redesign."

**What didn't:** The REQ's Red-Green "roughly half or less" word-count target wasn't met (27% reduction, not ~50%) — see the Word-count receipt note above for why (the firm four-item move-list is smaller than the target implies). Worth a note back to whoever wrote the REQ's Red-Green Proof: for a file this specialized, "half" may have been calibrated against a rougher average across all the batch's target files rather than this file's actual heavy-content footprint.

**Worth knowing:** Applying CLAUDE.md's earned-sections test to `ai-report.md` produced a clean pass with zero deletions — a useful data point for future maintenance passes that a narrow, single-purpose action file (as opposed to a broad orchestrator like `work.md`, or files with copy-pasted cross-file guards like the REQ-025 targets) may already satisfy the earned-sections bar throughout, precisely because every rule in it was written against this action's own specific failure modes rather than generic hygiene. Don't assume every action needs trimming just because the rule now exists — confirm row-by-row.
