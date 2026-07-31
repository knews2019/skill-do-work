---
id: UR-006
title: Align do-work with Claude 5 context-engineering guidance
created_at: 2026-07-27T07:34:50Z
requests: [REQ-025, REQ-026, REQ-027, REQ-028, REQ-029, REQ-030, REQ-031]
word_count: 31
---

# Align do-work with Claude 5 Context-Engineering Guidance

## Summary

The user asked the session to read Anthropic's post *"The New Rules of Context Engineering for Claude 5 Generation Models"* and then compare it against this skill to decide what needs updating. An audit of the tree against the post's five shifts (rules → judgment, examples → interface design, upfront context → progressive disclosure, repetition → focused descriptions, manual memory → auto-memory) produced a seven-REQ plan, approved by the user.

The skill already ran a bloat cleanup on 2026-07-15 (`decisions/audits/2026-07-15-harness-bloat-audit-phase1-2.md`) that cut `SKILL.md` from 5,507 → 2,396 words and installed a 2,650-word router ratchet. This UR covers the four areas that pass did **not** address.

## Findings the REQs Address

| Shift | Finding |
| --- | --- |
| Rules → judgment | Six domain crew files (~6,737 w) are largely generic engineering advice, with real overlap between `frontend.md`↔`security.md` and `backend.md`↔`security.md` |
| Examples → interfaces | `next-steps.md` (1,741 w, read after *every* action) is 40 hard-coded worked examples instead of a stated intent; `capture.md` carries 5 full transcripts; `ai-report.md` pastes raw HTML/CSS/JS inline |
| Progressive disclosure | Only 4 `*-reference.md` companions exist; 10 of the top-15 action files load whole — `ai-report.md` (7,541 w), `pipeline.md` (7,471 w), `capture.md` (5,680 w) worst |
| Repetition | The `git add -A` / `--no-verify` guard is restated across 7 files (5× inside `stray-check.md` alone); prompt-injection doctrine is inlined in `capture.md`, `bkb.md`, `validate-feedback.md` despite `crew-members/prompt-injection.md` existing to be JIT-loaded; "read-only" is restated 3–4× per file in `forensics.md`, `quick-wins.md`, `code-review.md` |

**Root cause of the boilerplate:** the maintainer-side action-file template prescribes Rules / Common Rationalizations / Red Flags / Verification Checklist, so 24–33 of 43 action files carry all four whether or not each has real content. Release 0.123.2 already deduped small actions once and the pattern came back — hence REQ-027 makes the fix structural rather than another one-off trim.

The fifth shift (manual memory → auto-memory) is deliberately **not** pursued: do-work's REQ/UR trail is the product, not a context tax.

## Extracted Requests

| REQ | Title | Kind |
| --- | --- | --- |
| REQ-025 | State each shipped guard once | maintenance |
| REQ-026 | next-steps.md — intent over enumeration | maintenance |
| REQ-027 | Make action-template sections earned, not mandatory | maintenance |
| REQ-028 | Trim domain crew files to opinions | maintenance |
| REQ-029 | Split `actions/ai-report.md` into an action + reference pair | implementation |
| REQ-030 | Split `actions/pipeline.md` into an action + reference pair | implementation |
| REQ-031 | Split `actions/capture.md` into an action + reference pair | implementation |

## Capture-Time Decisions (user-selected, 2026-07-26)

1. **Scope:** full — dedupe repetition + fix the action-file template + trim crew files + split the three largest monolith action files.
2. **Crew files:** trim to opinions only — cut lines a capable model already follows; keep lines that are a deliberate project preference or tied to do-work machinery. The user accepted the stated risk that this removes a safety net when running the skill on a weaker model.
3. **Execution:** dogfood — capture as REQs, then run the work pipeline over them.

## Batch Constraints

- **Sequence:** REQ-027 (the template change) gates REQ-029/030/031, because it determines what the split files are allowed to contain. REQ-025 gates REQ-031 (both touch `capture.md`).
- **One REQ per commit**, each with its own `actions/version.md` bump and descriptive `CHANGELOG.md` entry (no codenames, no duplicate version numbers or titles).
- **`bash _dev/tests/contract-regressions.sh` must pass clean before every commit.** It gates: `SKILL.md` ≤ 2,650 words; the Action Dispatch `work` row still passing `$ARGUMENTS`; `tools/checks/{archive-collision,preflight,scope-drift,qualify}.sh` existing, executable, and referenced by basename from `actions/work.md`; no `ultracode|fable` in the active runtime docs; `.gitattributes` keeping `/CLAUDE.md` and `/AGENTS.md` as `export-ignore`; and shipped files never citing the maintainer docs.
- **Word-count receipt** in every REQ's Implementation Summary: before/after `wc -w` for each touched file, plus the always-read floor (`SKILL.md` + `next-steps.md`) and the `do-work run` orchestrator load (`SKILL.md` + `actions/work.md` + `actions/work-reference.md` + `crew-members/general.md` + `crew-members/coding-guardrails.md` + `next-steps.md`, currently **25,231 w**).

## Out of Scope (do not silently absorb)

- **Relocation Plans A/B/C** (`decisions/audits/2026-07-15-relocation-extraction-plans.md`, ~47k words). Still unexecuted and still valid, but a distribution decision independent of this guidance.
- **`docs/`** (13.7k w, 23 files) — restates `actions/` content but is human-facing only and never agent-loaded, so it costs zero context.
- **The maintainer doc's Naming Conventions section** — a deliberate cross-project preference, not drift.
- **`SKILL.md` routing** — already at 2,557/2,650 words. This work must not grow it.

## Full Verbatim Input

> can you read: https://claude.com/blog/the-new-rules-of-context-engineering-for-claude-5-generation-models ?

> compare with the current skill, is there something that we need to update?

Followed by three option selections: scope = "Add monolith splits"; crew files = "Trim to opinions only"; execution = "Capture as REQs, then run". The resulting plan was approved unchanged.

---
*Captured: 2026-07-27T07:34:50Z*
