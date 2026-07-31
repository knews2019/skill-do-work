---
id: REQ-027
title: Make action-template sections earned, not mandatory
status: completed
created_at: 2026-07-27T07:34:50Z
claimed_at: 2026-07-27T07:40:27Z
completed_at: 2026-07-27T08:06:57Z
commit: f6cd577
user_request: UR-006
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-025, REQ-026, REQ-028, REQ-029, REQ-030, REQ-031]
batch: context-engineering-alignment
---

# Make Action-Template Sections Earned, Not Mandatory

## What

`CLAUDE.md` § Action File Conventions prescribes an action-file template whose Rules / Common Rationalizations / Red Flags / Verification Checklist sections are treated as near-mandatory. The result: 24–33 of the 43 action files carry all four, most filled with generic engineering advice that a capable model already follows. Two changes:

1. **Reword the template** so those four sections are included only when the file has content a capable model would otherwise get wrong — do-work machinery, a file/frontmatter contract, or a hard-won failure mode with a traceable origin. Make "generic engineering advice" an explicit *non*-reason to add a section. Add the positive rule: state intent, not a directive rule, when a capable model can infer the rest.
2. **Add a regrowth ratchet** to `_dev/tests/contract-regressions.sh`: fail a new action file whose Common Rationalizations rows contain no do-work-specific noun (REQ, UR, queue, frontmatter, pipeline, archive…). Same spirit as the existing router word budget.

While in `CLAUDE.md`: trim the Project Structure tree glosses and the Queue Path Convention section (both derivable from the repo). Keep every shell-trap gotcha and the Closed Enumerations Go Stale rule verbatim — those are hard-won and not inferable. **Leave Naming Conventions alone** — it is a deliberate cross-project preference, not drift.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, `crew-members/general.md`, `crew-members/coding-guardrails.md`, `crew-members/maintenance.md`, `CLAUDE.md` (current), and `_dev/tests/contract-regressions.sh` (current, for house style — the router-word-budget check and the self-citation check). Approach: (1) reword the Action File Conventions template in `CLAUDE.md` so Rules/Common Rationalizations/Red Flags/Verification Checklist are marked "include only if earned," give a concrete test (name the specific failure + where it happened), state the explicit non-reason (generic engineering advice), and state the omission rule (empty/generic is worse than absent); (2) trim Project Structure per-line glosses to what isn't obvious from the directory name, and compress Queue Path Convention to its one load-bearing sentence; leave every shell-trap gotcha, Closed Enumerations Go Stale, and Naming Conventions untouched; (3) add a regrowth-ratchet loop to `_dev/tests/contract-regressions.sh` that scans each `actions/*.md` file's `## Common Rationalizations` table rows for a do-work-specific noun, grandfathering the pre-REQ-027 action-file set (baked in as a checked-in baseline array) so the check doesn't retroactively fail the existing 24-file backlog — new files only. Validate red→green with a throwaway action file.
- [x] **[APPLY]:** Code written exactly as planned. Scope strictly limited to `CLAUDE.md` and `_dev/tests/contract-regressions.sh`.
- [x] **[UNIFY]:** `git diff --stat` reviewed both files — additive/reword only, no unrelated hunks. `git diff CLAUDE.md` confirms no hunk touches Naming Conventions or the shell-trap/Closed-Enumerations subsections. Ran `bash _dev/tests/contract-regressions.sh` clean on the real tree; wrote a throwaway `actions/zzz-throwaway-test.md` with a generic Common Rationalizations table → suite FAILed naming the file and the fix; edited the same throwaway file to reference REQ/frontmatter/queue → suite passed; deleted the throwaway file → suite still passes. No debug artifacts left in the diff (the throwaway file was removed, not committed as part of scope).

## Triage

**Route B** — outcome and locations were clear from the REQ (reword one prose section in `CLAUDE.md`, add one ratchet loop to a known test script); the open call was the ratchet's exact mechanics (grandfather baseline, noun list, row extraction), which the Builder Guidance explicitly deferred to the builder.

## Scope

**Files committed to before editing:**
- `CLAUDE.md` — reword Action File Conventions template (earned-sections test); trim Project Structure glosses; compress Queue Path Convention; leave shell-trap gotchas, Closed Enumerations Go Stale, and Naming Conventions untouched.
- `_dev/tests/contract-regressions.sh` — add the Common Rationalizations regrowth-ratchet check.

**Acceptance criteria restated from the REQ:**
- Template reworded so the four sections are earned via a concrete test, not a vibe; generic engineering advice is an explicit non-reason; omission (not empty/generic filler) is the stated fallback.
- Ratchet check added to `_dev/tests/contract-regressions.sh`; scans new action files' Common Rationalizations rows for a do-work-specific noun (illustrative list, marked as such per Closed Enumerations Go Stale).
- Check does not retroactively fail the existing tree (builder's call, recorded below: grandfather baseline of the 42 pre-REQ-027 action files).
- Project Structure tree kept; per-line glosses trimmed to non-obvious content. Queue Path Convention compressed to its one sentence.
- Shell-trap gotchas and Closed Enumerations Go Stale kept verbatim. Naming Conventions untouched (byte-identical).
- `bash _dev/tests/contract-regressions.sh` passes clean including the new check.
- Version bump + changelog entry — explicitly NOT this builder's job (orchestrator handles it per the task brief); not done here.
- Editing action files to conform — explicitly out of scope for this REQ.

## Why (if provided)

This is the root cause of the boilerplate the rest of the batch is cleaning up. Release 0.123.2 already deduped small action files once and the pattern grew back, because the template still asked for the sections. Trimming without changing the template buys one release of relief. Fixing the template plus adding a ratchet is the structural version — the same move the 2026-07-15 cleanup made when it paired the SKILL.md diet with the 2,650-word router budget.

## Context

- The four sections are genuinely valuable in some files: `actions/work.md`'s rationalizations encode real pipeline failure modes with traceable origins. The rule must distinguish those from "don't skip tests" filler, not ban the sections.
- `_dev/tests/contract-regressions.sh` is the repo's hard-ratchet suite (`_dev/` is `export-ignore`'d, so the check is maintainer-side only). Read the existing router-budget check for the house style — a clear FAIL message that names the fix, not just the violation.
- The new check must not retroactively fail the existing tree. Decide during implementation whether to scope it to files added after this REQ (e.g. compare against a checked-in baseline list) or to fix the offenders in the same commit; record the call and its rationale.
- Maintenance pass on the skill's own instructions — `crew-members/maintenance.md` loads via `maintenance: true`.

## Detailed Requirements

- **Files in scope:** `CLAUDE.md`, `_dev/tests/contract-regressions.sh`. Editing action files to conform is explicitly *out* of scope for this REQ — REQ-025/028/029/030/031 do that work under the new rule.
- **The reworded template must give a test, not a vibe.** A reader deciding whether to add a Common Rationalizations section needs a concrete question to answer ("can I name the specific failure this row prevents, and where it happened?"), not "include when useful."
- **Required vs. earned must be unambiguous.** The current text has Required / Common / Encouraged tiers; the reworded version must clearly place the four sections and state what happens when a file has nothing to put in one (omit it — an empty or generic section is worse than no section).
- **Ratchet check specifics:** for each action file's `## Common Rationalizations` table, scan the row text for at least one do-work-specific noun. Treat the noun list as illustrative and say so in a comment — per the Closed Enumerations rule, a hand-enumerated list goes stale as the vocabulary grows.
- **Trim, don't gut, `CLAUDE.md`.** The Project Structure tree stays; its per-line glosses shrink to what isn't obvious from the directory name. Queue Path Convention compresses to its one load-bearing sentence (`do-work/queue/`, not `do-work/` root).

## Constraints

- `bash _dev/tests/contract-regressions.sh` must pass clean *including the new check*.
- `CLAUDE.md` and `_dev/` are `export-ignore`'d, so nothing here ships — but any rule that a *shipped* file needs to know must be restated inline in that shipped file, never cited across the export boundary.
- Version bump + descriptive `CHANGELOG.md` entry.

## Dependencies

Gates REQ-029, REQ-030, REQ-031 — the new template determines what the split reference files are allowed to contain, so the splits must not run first.

## Builder Guidance

**Certainty: Firm on intent, mixed on the check's mechanics.** The wording change is decided. The ratchet's exact implementation (which files it scans, how it avoids failing the existing tree, how it tokenizes table rows) is the builder's call — prefer the simplest thing that catches a genuinely generic new table, and accept some false negatives over false positives that block honest work.

## Red-Green Proof

- **RED now:** `CLAUDE.md`'s template lists Rules / Common Rationalizations / Red Flags / Verification Checklist with no test for when to omit them, and `bash _dev/tests/contract-regressions.sh` passes even if you add an action file whose entire Common Rationalizations table is generic advice ("If you're thinking 'I'll skip the tests' → STOP. Write the tests → Because untested code breaks").
- **GREEN when:** adding that same generic-table action file makes the suite FAIL with a message naming the fix; adding a file whose rows reference REQ/queue/frontmatter passes. `CLAUDE.md` states an omission test the reader can apply.
- **Validation:** write the throwaway generic action file, run the suite, confirm FAIL; delete it, confirm the suite passes on the real tree; paste both runs into the Implementation Summary. Confirm Naming Conventions is byte-identical (`git diff` shows no hunk in that section).

## Open Questions

- [~] Should the ratchet fail the existing tree's offenders, or only newly added files? → Deferred to the builder. Recommended: scope to new/changed files so the check lands green and the batch's own REQs clean up the backlog. Also: fix all offenders in this commit (larger diff, breaks the one-REQ-one-concern boundary), or warn-only for existing files.

## Full Context

See `do-work/user-requests/UR-006/input.md` for complete verbatim input.

## Implementation Summary

**What was done:**

1. **`CLAUDE.md` § Action File Conventions** — the template's four "earned" sections now each carry an inline `[Include only if earned — see below.]` marker instead of an example placeholder that implied every file needed one. Replaced the old Required/Common/Encouraged tier sentence with: a **Required/Common** split (Description blockquote + Steps; Input/Output Format/When to Use), an **Earned, not mandatory** paragraph naming the two things that earn a section (do-work machinery, or a hard-won failure mode with a traceable origin) and the explicit non-reason ("generic engineering advice a capable model already follows"), **the test** ("can I name the specific failure this row prevents, and where it happened?") with the omission rule stated directly ("a generic table is worse than no table... omit the section entirely; don't ship it empty or generic to satisfy the template"), the **state-intent-not-a-directive-rule** positive rule with a worked example, and a pointer to the new ratchet check. Section order sentence kept, reworded to match the new "when present" framing.
2. **`CLAUDE.md` § Project Structure** — trimmed 8 of 14 tree-line glosses to drop content restated elsewhere or inferable from the name (`next-steps.md`, `README.md`, `actions/`, `specs/`, `prompts/`, `hooks/`, `tools/`, `docs/`, `decisions/`, `CHANGELOG.md`); kept every gloss carrying non-obvious behavior (`SKILL.md`, `interviews/`, `crew-members/`, `AGENTS.md`).
3. **`CLAUDE.md` § Queue Path Convention** — compressed two sentences to the one load-bearing fact (`do-work/queue/`, not `do-work/` root).
4. **`CLAUDE.md`** — left every shell-trap gotcha, the Closed Enumerations Go Stale rule, and the entire Naming Conventions section untouched (verified via `git diff` — no hunk touches those spans).
5. **`_dev/tests/contract-regressions.sh`** — added a new check block: a checked-in baseline array of the 42 action files that existed at REQ-027 time (grandfathered, so the existing 24-file generic-table backlog doesn't fail the suite); for every `actions/*.md` file NOT in that baseline, if it has a `## Common Rationalizations` section, extract the table's data rows (strip the header row and the `---` separator row) and require at least one do-work-specific noun. Noun check is split into a case-sensitive pattern for `REQ`/`UR-` (to avoid `-i` colliding with ordinary words like "required"/"requires") and a case-insensitive pattern for the rest (queue, frontmatter, pipeline, archive, do-work, domain, blocked, kb/, prime, clarify, working/, crew-member, schema, status:) — commented as illustrative, not exhaustive, per the Closed Enumerations Go Stale rule. Failure message names the offending file and points at the fix (add a specific row, or omit the section) and at `CLAUDE.md`'s omission test (this script is `_dev/`-only, export-ignored, so citing `CLAUDE.md` here is allowed — it's not a shipped file).

**Decision recorded (Open Question in the REQ):** scoped the ratchet to *new* action files via a grandfather baseline, not a mass rewrite of the 24-file existing backlog — per the REQ's own recommendation and the one-REQ-one-concern boundary (REQ-025/028/029/030/031 clean up the backlog under the new rule; editing action files was explicitly out of this REQ's scope).

**Files changed:**
- `CLAUDE.md` (modified) — Action File Conventions template reworded (earned-sections test); Project Structure glosses trimmed; Queue Path Convention compressed. 2,598 → 2,859 words (net +261; the template rewording's concrete test/positive-rule prose outweighs the Project-Structure/Queue-Path trims — CLAUDE.md carries no word budget, unlike SKILL.md).
- `_dev/tests/contract-regressions.sh` (modified) — added the Common Rationalizations regrowth-ratchet check (baseline array, noun patterns, scan loop). 1,097 → 1,493 words.

**Word-count receipt (`wc -w`, before → after):**

| File | Before | After |
| --- | --- | --- |
| `CLAUDE.md` (touched) | 2,598 | 2,859 |
| `_dev/tests/contract-regressions.sh` (touched) | 1,097 | 1,493 |
| `SKILL.md` (always-read floor, untouched) | 2,557 | 2,557 |
| `next-steps.md` (always-read floor, untouched by this REQ — a concurrent agent is mid-edit on it per the task brief; snapshot only) | 1,741 | 383\* |
| `actions/work.md` (orchestrator load, untouched) | 11,983 | 11,983 |
| `actions/work-reference.md` (orchestrator load, untouched) | 7,426 | 7,426 |
| `crew-members/general.md` (orchestrator load, untouched) | 678 | 678 |
| `crew-members/coding-guardrails.md` (orchestrator load, untouched) | 846 | 846 |

\* `next-steps.md`'s "after" count reflects another agent's concurrent in-flight edit observed at the time this REQ ran its receipt — not a change made by this builder, and outside this REQ's declared file scope (`CLAUDE.md` and `_dev/tests/contract-regressions.sh` only).

## Testing

**Contract suite (real tree, before any throwaway file):**
```
$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
```

**RED — throwaway `actions/zzz-throwaway-test.md` with an all-generic Common Rationalizations table** (rows: "I'll skip the tests" / "This looks fine" — no do-work noun):
```
$ bash _dev/tests/contract-regressions.sh
FAIL: zzz-throwaway-test.md Common Rationalizations table has no do-work-specific noun (REQ, UR, queue, frontmatter, pipeline, archive, domain, blocked, kb/, prime, clarify, working/, crew-member, schema, status — illustrative list) in any row — every row reads as generic engineering advice a capable model already follows. Add rows naming a specific do-work failure mode, or omit the section entirely (see CLAUDE.md -> Action File Conventions for the omission test).
exit code: 1
```

**GREEN — same throwaway file, row rewritten to reference REQ/frontmatter/queue** ("I'll archive the REQ without checking status" / "Read the REQ frontmatter first" / "The queue depends on accurate status transitions"):
```
$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
exit code: 0
```

**Cleanup — throwaway file deleted, real tree re-verified:**
```
$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
exit code: 0
```

**Naming Conventions byte-identical check:** `git diff CLAUDE.md | sed -n '/Naming Conventions/,/^$/p'` produced no output — confirms no hunk touches that section.

## Lessons Learned

**What worked:** Grandfathering the existing action-file set via a checked-in baseline array (same pattern the script already uses for `hardened_check_scripts` and `shipped_citation_paths`) let the ratchet land green immediately without a same-commit mass rewrite, matching the REQ's "scope to new/changed files" recommendation and keeping this REQ's diff to exactly the two files in scope.

**What didn't:** A naive single case-insensitive noun pattern including `REQ` would have been useless — `-i` makes `REQ` match inside ordinary words like "required" and "requires," which show up constantly in generic advice, silently defeating the check's whole purpose. Had to split into a case-sensitive pattern for `REQ`/`UR-` and a case-insensitive pattern for the rest.

**Worth knowing:** Common Rationalizations tables don't have a fixed column width or exact header text across the 24 existing files, but every one starts data rows with `|` after a header row containing "If you're thinking" and a `---`-only separator row — extracting by stripping those two identifiable rows (rather than assuming a fixed line offset) is what makes the row-scan robust across the existing table-formatting variance. Also: under `set -euo pipefail`, an unguarded `var="$(cmd1 | cmd2 | grep ...)"` where the final grep matches nothing exits the whole script (pipefail treats "no match" as pipeline failure) — every multi-stage grep pipeline feeding the new check is terminated with `|| true`.

---
*Source: "compare with the current skill, is there something that we need to update?" — resolved into the approved seven-REQ plan.*

Think carefully before answering.
