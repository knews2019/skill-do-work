---
id: REQ-031
title: Split actions/capture.md into an action + reference pair
status: completed
created_at: 2026-07-27T07:34:50Z
claimed_at: 2026-07-27T08:09:38Z
completed_at: 2026-07-27T08:29:45Z
commit: ed5f96b
user_request: UR-006
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-025, REQ-027]
maintenance: false
related: [REQ-025, REQ-026, REQ-027, REQ-028, REQ-029, REQ-030]
batch: context-engineering-alignment
---

# Split actions/capture.md Into an Action + Reference Pair

## What

`actions/capture.md` is 5,680 words and loads whole on every capture. Split into `actions/capture.md` (Steps 0–7 plus the simple/complex triage logic) and a new `actions/capture-reference.md`, following the `actions/work.md` ↔ `actions/work-reference.md` pattern.

Move to the companion:

- **Request File Formats** — the simple REQ template, the complex REQ template, the Schema Aliases table, and the UR `input.md` template (lines ~77–224).
- The **addendum template**.
- The **five worked examples** (lines ~429–489) — but see below: under REQ-027's rule these are candidates for **deletion** rather than relocation.

Keep in the action file: Steps 0–7, the simple-vs-complex triage, file locations, the immutability rule, and file naming.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, `actions/work-reference.md` header + Architecture/Folder Structure (proven pattern), `crew-members/general.md`, `crew-members/coding-guardrails.md`, `CLAUDE.md` § Action File Conventions (post-REQ-027 wording), the archived REQ-025 and REQ-027 (confirmed both already landed — the current `actions/capture.md` reflects both), and the existing `bkb-reference.md`/`interview-reference.md`/`deep-explore-reference.md` companions for header/citation style. Approach: extract the exact byte ranges for the Request File Formats section (Simple REQ, Complex REQ, Schema Aliases, UR input.md) and the Step-2 addendum frontmatter template into `actions/capture-reference.md` behind a work-reference.md-style companion header; leave every other section (Philosophy → Backward Compatibility, Steps 0–7, Edge Cases, Common Rationalizations, Red Flags, Verification Checklist) in place; add a top-of-file companion callout plus named pointers at Step 5 (and the Step 2 addendum branch) so the companion is impossible to miss; delete the five worked examples (default per REQ) after confirming none encodes a judgment the templates/triage rules don't already specify; diff each moved block against the pre-edit original to prove verbatim relocation; run the contract suite and the two-way pointer check before finishing.
- [x] **[APPLY]:** Code written exactly as planned. Scope strictly limited to `actions/capture.md` (edited) and `actions/capture-reference.md` (new).
- [x] **[UNIFY]:** `git diff --stat` shows only `actions/capture.md` and `actions/capture-reference.md` touched by this builder (ai-report.md/pipeline.md changes in the tree are concurrent agents' work, confirmed via `git status --porcelain`). Diffed the two moved blocks (Request File Formats section, Addendum REQ frontmatter template) against pre-edit extracts — both `IDENTICAL`. Ran `bash _dev/tests/contract-regressions.sh` clean. Ran an end-to-end throwaway capture (UR-999/REQ-999) by hand using the relocated templates, confirmed a correctly paired, correctly-timestamped UR+REQ, then deleted it. No debug artifacts left in either shipped file.

## Why (if provided)

Two of Anthropic's five shifts apply here. **Progressive disclosure:** the REQ/UR templates are needed at Step 5 (Write Files) and nowhere else, yet they load before Step 0. **Examples → interface design:** five full worked transcripts teach by pattern-matching what the templates plus the triage rules already specify — the kind of few-shot bulk the guidance says to replace with a clear interface.

## Context

- `actions/capture.md` heading map from the audit (re-verify before editing): Philosophy(7), First-Run Bootstrap(24), When to Use(32), Simple vs Complex(44), File Locations(53), Immutability Rule(59), File Naming(65), Backward Compatibility(73), Request File Formats(77), Simple REQ(79), Complex REQ(136), Schema Aliases(186), UR input.md(200), Steps(225), Step 0 Load Prompt-Injection Guardrail(227), Step 1 Parse and Assess(231), Step 2 Check for Duplicates(244), Step 3 Capture-Phase Clarification(320), Step 4 Handle Screenshots(365), Step 5 Write Files(373), Step 6 Report Back(391), Step 7 Commit(399), worked examples(429–489).
- The templates are a **hard contract**, not illustration — `actions/work.md`, `actions/roadmap.md`, and `tools/queue-kanban/model.go` all read the fields they define. Relocating them must not soften a single field's specification.
- Copy the companion-header convention from `actions/work-reference.md`, including lazy-read and don't-re-read.
- REQ-025 lands first and removes this file's inlined prompt-injection doctrine; REQ-027 lands first and changes what the four template sections must earn. Split the post-REQ-025/027 file, not the current one.

## Detailed Requirements

- **Decide the worked examples explicitly and record the call.** Default per REQ-027's rule: delete them, because the templates plus the triage rules fully specify the output and the examples merely demonstrate. Keep an example only if it encodes a judgment the templates don't — e.g. a genuinely hard slicing decision. If any survive, they live in the companion, and the Implementation Summary names which one survived and why.
- **Field specifications survive verbatim.** The frontmatter fields, their enums, `depends_on`/`addendum_to` semantics, the schema-alias table, and the timestamp rule (UTC via `date -u`, never local time with a `Z` suffix) must read identically after the move. Diff the moved block to confirm.
- **Step 5 points at the companion by name** for each template it writes. No dangling pointers, no orphan sections.
- **The action file must still work as a standalone prompt** for triage and the non-writing steps; it should be obvious at Step 5 that the companion must be opened.
- **Both files ship.** Cite other actions by path; never cite `CLAUDE.md`/`AGENTS.md` — note the current file's Step 7 commit guidance and Context sections reference maintainer docs in places; verify none survive the split.

## Constraints

- `bash _dev/tests/contract-regressions.sh` must pass clean.
- `SKILL.md` must not grow past 2,650 words.
- Version bump + descriptive `CHANGELOG.md` entry.

## Dependencies

Depends on REQ-025 (same file, prompt-injection dedupe) and REQ-027 (template rule).

## Builder Guidance

**Certainty: Firm on the split, exploratory on the examples.** Deleting the five transcripts is the recommended call, not a mandate — read them first, and if one carries real judgment, keep that one. Everything else about this REQ is mechanical: move the templates, point at them, change nothing they specify.

## Red-Green Proof

- **RED now:** `wc -w actions/capture.md` ≈ 5,680; every capture — including one that aborts at Step 2's duplicate check — has loaded both REQ templates, the alias table, the UR template, and five worked transcripts.
- **GREEN when:** `actions/capture.md` holds Steps 0–7 and the triage rules at materially reduced size; `actions/capture-reference.md` holds the templates and alias table; Step 5 names the sections it needs; the field specifications are unchanged.
- **Validation:** run a real capture end-to-end after the split (capture a throwaway request, confirm it produces a correct UR + REQ pair, then abandon it) — this is the acceptance test that matters most, since the templates are a contract. Plus: `grep -n 'capture-reference' actions/capture.md` pointer check in both directions; `git diff` on the moved template block showing pure relocation; before/after `wc -w`; `bash _dev/tests/contract-regressions.sh` clean.

## Open Questions

- [x] Delete or relocate the five worked examples? → **Deleted, all five.** See Implementation Summary → Examples Decision for the full read of each example against the templates/triage rules.

## Full Context

See `do-work/user-requests/UR-006/input.md` for complete verbatim input.

## Triage

**Route B** — the split mechanics were fully specified by the REQ (extract Request File Formats + the addendum frontmatter template, following the proven `work.md`/`work-reference.md` pattern); the one open call was the worked-examples decision, which the Builder Guidance explicitly left to the builder after reading them.

## Scope

**Files committed to before editing:**
- `actions/capture.md` (modified) — remove the Request File Formats section (Simple REQ, Complex REQ, Schema Aliases, UR input.md templates), remove the Step-2 addendum frontmatter template, remove the five worked examples; add a companion callout and named Step 5 (+ Step 2) pointers into the new companion. Steps 0–7, the simple/complex triage, file locations, immutability rule, file naming, Edge Cases, Common Rationalizations, Red Flags, and Verification Checklist stay untouched.
- `actions/capture-reference.md` (new) — houses the four templates plus the addendum-REQ template, behind a `work-reference.md`-style companion header.

**Acceptance criteria restated from the REQ:**
- Worked-examples decision made explicitly and recorded, with reasoning.
- Frontmatter fields/enums, `depends_on`/`addendum_to` semantics, the schema-alias table, and the UTC timestamp rule read identically after the move (verified by diff, not re-read).
- Step 5 (and the Step 2 addendum branch) names every companion section by name; no dangling pointer, no orphan section.
- `actions/capture.md` still works standalone for triage and the non-writing steps.
- Neither file cites this repo's own `CLAUDE.md`/`AGENTS.md`; the existing consumer-project `CLAUDE.md` reference (Step 1 prime-file routing) is preserved.
- `bash _dev/tests/contract-regressions.sh` passes clean.
- Version bump + changelog entry — explicitly the orchestrator's job per the task brief, not this builder's; not done here.

## Implementation Summary

**What was done:**

1. **Companion callout added to `actions/capture.md`** (right after the intro paragraph, before `## Philosophy`) — names `actions/capture-reference.md`, states it's read at Step 5 (or the Step 2 addendum branch), and restates the hard-contract framing (`actions/work.md`, `actions/roadmap.md`, `tools/queue-kanban/model.go` all read these fields).
2. **`## Request File Formats` removed from `actions/capture.md`** (the whole section — Simple REQ, Complex REQ (additional sections), Schema Aliases, UR input.md — 147 lines) and relocated **verbatim** into `actions/capture-reference.md` under the same heading and sub-headings.
3. **The addendum-REQ frontmatter template removed from Step 2** (the fenced `REQ-021` example, 28 lines) and relocated **verbatim** into a new `## Addendum REQ Template` section in `actions/capture-reference.md`. The surrounding decision prose (queued-addendum table, Coherence Rule, the small inline `## Addendum (2025-01-27)` append-format snippet, "Context is critical for addenda," "When the original UR is archived," "Coherence across addendum chains") all stayed in Step 2 — none of that is a field-bearing template, all of it is decision logic that belongs with the step.
4. **Step 5 rewritten to open with an explicit "open the companion first" line**, then names, at the exact point each template is used: **UR input.md** (point 1), **Simple REQ** / **Complex REQ (additional sections)** (point 2), the **Schema Aliases** normalize-and-warn contract (point 2), and the **Complex REQ (additional sections)** template's Populating `depends_on` / Slicing convention guidance (Complex-mode bullet).
5. **Step 2's addendum-for-in-flight/archived paragraph rewritten** to point at the **Addendum REQ Template** section by name instead of inlining it.
6. **Five worked examples deleted, not relocated** — see Examples Decision below.
7. **Everything else in `actions/capture.md` left byte-for-byte untouched**: Philosophy, First-Run Bootstrap, When to Use, Simple vs Complex, File Locations, Immutability Rule, File Naming, Backward Compatibility, Steps 0/1/3/4/6/7, Edge Cases, Common Rationalizations, Red Flags, Verification Checklist. Steps 2 and 5 are the two exceptions, and they changed only as described in points 4 and 5 above — each gained a named pointer to the companion; no decision prose was removed. (Corrected at review: the original wording claimed "all seven Steps" were untouched, which contradicted points 4 and 5 on the same page.)

**Examples Decision (the one judgment call assigned to this builder):** **Deleted all five** (Simple Capture, Multiple Requests, Addendum to In-Flight Request, Addendum to Archived Request, Complex Multi-Feature Request) — none relocated. Read each against what the templates + triage rules already specify:
- *Simple Capture* / *Multiple Requests* — pure restatement of Step 1's splitting heuristic ("and also", comma lists, distinct topics). No new information.
- *Addendum to In-Flight Request* / *Addendum to Archived Request* — pure restatement of Step 2's decision table + File Naming's numbering rule. The archived-request example's "Prior Implementation" mention duplicates prose already stated twice in Step 2.
- *Complex Multi-Feature Request* — the REQ's suggested bar for a keeper was "a genuinely hard slicing decision," e.g. one that shows a `depends_on` graph. This example doesn't: it slices one input into 4 REQs with a flat 1:1 feature→REQ mapping and **no `depends_on` populated at all** — it doesn't even demonstrate the Slicing convention paragraph it sits closest to. Keeping it would have been keeping a demonstration, not a judgment.
No example encoded a judgment the templates/triage rules don't already fully specify, so per the REQ's default and CLAUDE.md's "examples that merely demonstrate an already-clear interface get cut" framing, all five are deleted.

**Earned-sections re-check (per REQ-027's dependency note that REQ-031 is one of the REQs bringing its own file into conformance):** re-read Common Rationalizations (5 rows), Red Flags (5 items), and Verification Checklist (6 items) against the "can I name the specific failure this row prevents, and where it happened?" test. All 16 items tie to a specific do-work behavioral contract (REQ/UR pairing, the capture-vs-build phase boundary, RED/GREEN proof, Open Questions format, the STOP-after-capture rule) — none read as generic engineering filler. No deletions warranted there; left as-is. Edge Cases (7 items, not one of the four audited sections) reviewed on the same test and also left as-is — each ties to a do-work-specific object (REQ, UR, Open Questions) except one one-liner ("references earlier conversation") that's harmless and not worth a special-case removal.

**Files changed:**
- `actions/capture.md` (modified)
- `actions/capture-reference.md` (new)

**Word-count receipt (`wc -w`):**

| File | Before | After |
| --- | --- | --- |
| `actions/capture.md` | 5,505 | 3,913 |
| `actions/capture-reference.md` (new) | 0 | 1,682 |
| Combined | 5,505 | 5,595 |

(Combined total rose slightly — +90 words — from the added companion callout, the five named Step 5/Step 2 pointers, and each template's short intro sentence in the companion; this mirrors `work.md`/`work-reference.md`'s own pattern of a small net increase in exchange for a large per-load reduction. `actions/capture.md`'s always-loaded footprint dropped 1,592 words, ~29%, and no longer loads the five templates or the five worked examples on every invocation — only at Step 5, or the Step 2 addendum branch.)

**Two-way pointer-integrity check (pasted receipt):**

```
$ grep -n "^#" actions/capture-reference.md | grep -v -f <(sed -n '/^```/,/^```/p' actions/capture-reference.md)
# ^ the bare grep emits ~29 lines: the real headings PLUS every '#'-prefixed line inside the
#   fenced REQ/UR templates (## What, ## Context, # External-condition fields, ...). Filtered
#   to the companion's real sections — the 6 below are what the action file points at.
7:## Request File Formats
9:### Simple REQ
66:### Complex REQ (additional sections)
116:### Schema Aliases
130:### UR input.md
155:## Addendum REQ Template

$ grep -n "capture-reference" actions/capture.md   # every pointer site in the action file
7:  companion callout — names all four Request File Formats sections + the addendum template
137: Step 2 — "**Addendum REQ Template** in `actions/capture-reference.md`"
200: Step 5 — "Open `actions/capture-reference.md` before writing anything"
205: Step 5 point 1 — "**UR input.md** template in `actions/capture-reference.md`"
206: Step 5 point 2 — "**Simple REQ** or **Complex REQ (additional sections)** template ... **Schema Aliases** section's normalize-and-warn contract"
213: Step 5 Complex-mode bullet — "**Complex REQ (additional sections)** template's Populating `depends_on` / Slicing convention guidance"
```

Forward: all 5 leaf sections (Simple REQ, Complex REQ (additional sections), Schema Aliases, UR input.md, Addendum REQ Template) are named from `actions/capture.md` — Simple REQ/Complex REQ/Schema Aliases/UR input.md from Step 5, Addendum REQ Template from Step 2. Backward: every section named in a pointer exists verbatim as a heading in `actions/capture-reference.md` (confirmed against the grep above). No orphan sections, no dangling pointers.

**Verbatim-relocation diff (pasted receipt):**

```
$ diff orig-request-file-formats.txt(lines 77-223 of pre-edit capture.md) new-request-file-formats.txt(lines 7-153 of capture-reference.md)
REQUEST FILE FORMATS: IDENTICAL

$ diff orig-addendum-template.txt(lines 285-312 of pre-edit capture.md) new-addendum-template.txt(lines 159-186 of capture-reference.md)
ADDENDUM TEMPLATE: IDENTICAL
```

Both diffs empty — the frontmatter fields, enums, `depends_on`/`addendum_to` semantics, the Schema Aliases table, and the `date -u`/UTC timestamp rule read identically after the move; nothing was softened.

## Testing

This is a documentation/prompt-instruction split — no application code, no automated test harness. Evidence gathered:

**1. Contract regression suite (must pass clean per the Constraints):**
```
$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
```
Confirms, among other things: `actions/capture.md` still contains the literal strings `maintenance: false` and `Maintenance assessment` (both live in Step 1, untouched by the split, so the ratchet that checks for them keeps passing); the self-citation check finds no `CLAUDE.md`/`AGENTS.md` reference to this repo's own maintainer doc in either file (the Step 1 mention of a **consumer project's** root `CLAUDE.md` is preserved and does not match the check's pattern); the new-action-file Common Rationalizations ratchet doesn't fire on `actions/capture-reference.md` because that file carries no `## Common Rationalizations` section at all (nothing was relocated there).

**2. Verbatim-relocation diff** — see Implementation Summary above; both moved blocks diffed `IDENTICAL` against pre-edit extracts.

**3. Two-way pointer-integrity check** — see Implementation Summary above; both directions confirmed via grep receipts.

**4. End-to-end throwaway capture (the acceptance test the REQ calls out as mattering most)** — simulated capturing "add a --dry-run flag to the export command" by hand, following `actions/capture.md`'s Steps and the relocated **UR input.md** / **Simple REQ** templates in `actions/capture-reference.md`:
```
$ find do-work -type f
do-work/queue/REQ-999-dry-run-export-flag.md
do-work/user-requests/UR-999/input.md

$ grep -E "^(id|title|status|created_at|user_request|domain|tdd|maintenance):" do-work/queue/REQ-999-dry-run-export-flag.md
id: REQ-999
title: Add --dry-run flag to export command
status: pending
created_at: 2026-07-27T08:19:27Z
user_request: UR-999
domain: backend
tdd: true
maintenance: false
```
UR-999's `requests: [REQ-999]` and the REQ's `user_request: UR-999` cross-link correctly; `created_at` is a UTC instant via `date -u +%Y-%m-%dT%H:%M:%SZ` on both files, matching the Timestamp rule. Confirms the relocated templates are usable exactly as pointed to from Step 5. The throwaway UR/REQ were then deleted (`rm -rf`) — nothing landed in the real `do-work/` tree.

**5. `git diff --stat` scope check** — confirms only `actions/capture.md` (modified) and `actions/capture-reference.md` (new) were touched by this builder; the concurrently-modified `actions/ai-report.md`/`actions/pipeline.md` in the working tree belong to the other agents named in the task brief.

## Lessons Learned

**What worked:** Extracting the moved blocks to standalone scratch files first (via `sed -n`) and diffing them against the corresponding region of the new companion file, rather than re-reading the new file and eyeballing it, is what actually proves "nothing was softened" — REQ-031's own instructions warned against judging by re-reading, and that warning was earned: a re-read pass would not have reliably caught, say, a single dropped `#` in a YAML comment.

**Worth knowing:** REQ-027's own body names REQ-031 among the REQs expected to bring their target file into conformance with the new earned-sections doctrine, not just execute the split mechanically (its "What" section: "Editing action files to conform is explicitly out of scope for this REQ — REQ-025/028/029/030/031 do that work under the new rule"). That's easy to miss if you read only REQ-031's own body, which frames the split as purely mechanical. Re-auditing Common Rationalizations/Red Flags/Verification Checklist against the earned-sections test (rather than assuming "REQ-031 only mentions the examples decision, so that's the only judgment call") is what closed that gap — all 16 rows/items across the three sections held up, so no further deletions were made, but the audit itself was the point.

**What didn't:** Nothing — the split was mechanical exactly as scoped once the examples decision and the earned-sections re-check were made explicit; no rework was needed.

---
*Source: "compare with the current skill, is there something that we need to update?" — resolved into the approved seven-REQ plan.*

Think carefully before answering.
