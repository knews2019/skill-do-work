---
id: REQ-202
title: Complete unsafe Remotion preview mutation detection
status: completed
domain: testing
created_at: 2026-08-15T18:45:11Z
claimed_at: 2026-08-15T20:33:47Z
completed_at: 2026-08-15T20:53:26Z
route: B
commit: 536fbd6
kb_status: pending
kb_entry:
user_request: UR-042
addendum_to: REQ-192
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: unsafe-remotion-preview-mutation-detection
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
write_set: [_dev/tests/contract-regressions.sh]
---

# Review Fix: Complete Unsafe Remotion Preview Mutation Detection

## What

Make the completed-work presentation regression detector recognize every executable fixed-port and platform-opener form while preserving safe foreground preview commands and non-executable prohibition prose.

This is one sweep because the missed forms share a single root cause: the unsafe-command extractor matches a narrow set of literal spellings instead of the complete prohibited executable command families.

## Context

Found during review of REQ-192. `_dev/tests/contract-regressions.sh` catches background/sleep, literal `open http://localhost:3000`, and render workflows, but mutation probes show that `remotion studio src/Root.tsx --port 3000` and `open "$REMOTION_PREVIEW_URL"` can pass.

## Instances

- [x] Fixed-port Remotion Studio flags, including separated and equals forms such as `--port 3000` and `--port=3000`.
- [x] Command-start platform opener forms whose target is a variable or other nonliteral expression, including `open "$url"`.

## Requirements

- Reject executable fixed-port Remotion Studio forms regardless of whether the flag uses a space or `=`.
- Reject executable command-start platform opener forms without depending on a literal localhost URL.
- Preserve the documented safe foreground `npm run preview` workflow.
- Do not treat explanatory prohibition prose as executable workflow content.
- Add replayable positive and negative mutation cases for every widened matcher family.

## Red-Green Proof

**RED prompt/case:** Feed the current unsafe-form detector `remotion studio src/Root.tsx --port 3000`, `remotion studio src/Root.tsx --port=3000`, and `open "$REMOTION_PREVIEW_URL"`; each prohibited executable form currently produces no match.
**Why RED now:** A future fixed-port or macOS-opener regression can pass the presentation contract suite even though the source-only video contract prohibits both workflows.
**GREEN when:** All executable fixed-port and platform-opener mutations fail, safe foreground preview examples pass, negative prose is ignored, and the focused and canonical suites remain green.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Isolate the executable-command detector, define positive and negative mutation cases for fixed-port Studio and platform-opener families, and preserve the safe foreground preview contract.
- [x] **[APPLY]:** Refactored the detector to evaluate supplied Markdown source, captured six semantic mutation failures, then widened only numeric fixed-port Studio commands and command-boundary lowercase `open` commands while retaining safe foreground/prohibition controls.
- [x] **[UNIFY]:** Passed focused syntax/contracts, implementation-summary qualification, scope drift, diff hygiene, and canonical maintainer verification; the change is ready for independent review before release.

## Triage

**Route: B — Medium.** The failure and owning test file are known, but the existing Markdown executable-segment parser needs focused exploration so the matcher expands without treating prohibition prose or safe foreground preview commands as unsafe.

## Exploration

- The defect is isolated to `_dev/tests/contract-regressions.sh`; the shipped action and guide already prescribe safe foreground preview commands and prohibit the unsafe workflows.
- The current detector misses separated and equals `--port` flags on `remotion studio`, plus command-start `open` forms whose target is not a literal localhost URL.
- The lowest-risk shape is a source-text helper shared by live file checks and an in-memory mutation table, with family labels asserted for unsafe cases and zero findings required for safe commands and negative prose.
- No action, guide, prescribed-shell test, or new test file needs modification.

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` (modify) — add replayable source-text mutations and complete the fixed-port Studio and platform-opener matchers.

**Files I will NOT touch:** shipped presentation actions or guides, prescribed-shell suites, queue-board files, existing archive records, generated artifacts, or unrelated concurrent documentation edits. Shared release files are lifecycle-only edits after product verification.

**Acceptance criteria:**
- [x] Both `--port 3000` and `--port=3000` executable Studio forms fail, including another numeric literal proving the matcher is not hard-coded to one port.
- [x] Command-start or chained platform `open` forms fail even with variable and parameter-expansion targets.
- [x] Safe foreground `remotion studio src/Root.tsx` and `npm run preview` commands pass.
- [x] Negative explanatory prose and ordinary “Open …” prose remain ignored.
- [x] Focused and canonical suites pass.

## Decisions

### D-01: Recognize platform openers by shell command boundary and exact command casing

**Decision:** Match lowercase `open` only at a command start or after a shell chaining operator, while leaving ordinary capitalized “Open …” prose outside the executable family.

**Why:** The macOS executable is the lowercase command token. Command-boundary matching catches literal, variable, and parameter-expansion targets without turning ordinary documentation sentences into false positives.

### D-02: Cross physical lines only through an explicit shell continuation

**Decision:** Let the fixed-port matcher cross a newline only when the Studio command ends that physical line with `\`, and carry a prohibition lead-in only across its immediately following blank/list/indented-inline example lines.

**Why:** Review attempt 1 showed that physical-line matching missed real continued commands while treating each explanatory Markdown line independently. These two bounded continuations cover the executable and prose structures without scanning arbitrary later paragraphs as one command or one prohibition.

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — extracted source-text unsafe-video findings, added replayable unsafe/safe mutation tables, and completed indented/chained platform-opener plus same-line/continued numeric Studio-port detection while preserving multi-line prohibition examples.

**What was done:** The live presentation files and in-memory replay cases now use one detector seam. Unsafe mutations cover separated, equals, quoted, and shell-continued numeric port values plus direct, indented, and chained variable opener targets. Prohibition context now spans only its immediate Markdown example list, while controls retain foreground Studio/package preview commands, same-line and multi-line prohibition prose, and ordinary opener prose.

## Testing

**Tests run:** `bash -n _dev/tests/contract-regressions.sh`; `bash _dev/tests/contract-regressions.sh`; `git diff --check -- _dev/tests/contract-regressions.sh`; `skills/do-work/tools/checks/qualify.sh do-work/working/REQ-202-complete-unsafe-remotion-preview-mutation-detection.md`; `skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-202-complete-unsafe-remotion-preview-mutation-detection.md`; `bash _dev/tests/maintainer-verify.sh`

**Result:** ✓ Shell syntax, focused contract regressions, scope drift, diff hygiene, and canonical maintainer verification passed. The first qualification pass correctly reported the still-open UNIFY checkbox; after recording the completed verification, qualification passed.

**Review remediation attempt 1:** The focused syntax-and-contract command passed with exit 0 after closing both Important findings. Canonical maintainer verification was not rerun by the remediation builder; the orchestrator owns that later gate.

**Red-green validation:**
- ✗ RED — after adding the source-text replay table but before widening the matcher, `bash -n _dev/tests/contract-regressions.sh && bash _dev/tests/contract-regressions.sh` reported six semantic escapes: separated, equals, and two quoted numeric Studio-port mutations plus direct-variable and chained-parameter-expansion `open` mutations.
- ✓ GREEN — after the minimal matcher update, the same command passed all unsafe mutations, all safe/prohibition controls, staged-skill checks reached by the suite, and the complete contract-regression run with exit 0.
- ✗ REMEDIATION RED (review attempt 1) — after adding the review replay cases, the same command reported four semantic escapes for two-space/tab-indented `open` plus separated/equals continued ports, and rejected the safe multi-line prohibition list with three false-positive findings.
- ✓ REMEDIATION GREEN — after bounded indentation, shell-continuation, and prohibition-list handling, `bash -n _dev/tests/contract-regressions.sh && bash _dev/tests/contract-regressions.sh` completed with exit 0.

**New tests added:** `_dev/tests/contract-regressions.sh` — in-memory unsafe mutations for same-line/continued fixed ports and direct/indented/chained openers, plus safe cases for foreground preview, same-line and multi-line explanatory prohibition prose, and ordinary “Open …” prose.

## Review

**Overall: 97% — Approve (Pass)**

**Route:** B

### Summary

Remediation attempt 1 closes both prior Important findings at the shared detector seam. The matcher now recognizes indentation at a shell command boundary, crosses physical lines only through an explicit backslash continuation, and carries prohibition context across the immediately following Markdown example list without masking a later positive command paragraph. The replay table directly locks each remediated form.

### Findings

**Critical:** None.

**Important:** None.

**Minor:** None.

### Prior Findings

- **Indented opener / continued fixed-port finding: Closed.** Direct probes detect space- and tab-indented `open`, separated and equals backslash-continued `--port`, additional indentation and continuation variants, and `npx remotion studio` with a continued numeric port.
- **Multi-line prohibition finding: Closed.** Bullet, ordered-list, and indented-inline prohibition examples remain ignored; a later positive opener paragraph is still detected. Foreground Studio/package preview and ordinary capitalized “Open …” prose remain safe.

### Acceptance

**Result: Pass.** The exact requested mutation families and negative controls are closed, both live video surfaces still use the shared detector, and focused plus canonical verification pass.

| Dimension | Score |
|---|---:|
| Requirements | 100% |
| Code Quality | 94% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

No follow-up is required.

## Lessons Learned

**What worked:** Extracting a source-text detector made executable-safety rules directly mutation-testable without modifying shipped presentation instructions. Family-labeled assertions also prevented one unsafe pattern from accidentally masking another.

**What didn't:** The first matcher expansion overfit the exact one-line examples. Independent adversarial review exposed shell-significant indentation and continuation plus a multi-line documentation shape that the initial tests omitted.

**Worth knowing:** Safety matchers for Markdown command examples need two bounded grammars: shell continuation for executable content and local structural continuation for explanatory prohibition examples. Crossing arbitrary newlines in either direction creates false negatives or false positives.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Orientation

The completed-work presentation contract now rejects executable fixed-port Remotion Studio and macOS opener variants across same-line, indented, chained, and explicitly continued forms while preserving foreground preview commands and explanatory prohibition prose. The canonical detector and replay matrix live together in `_dev/tests/contract-regressions.sh`.
