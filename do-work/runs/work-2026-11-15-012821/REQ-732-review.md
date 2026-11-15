## Review: REQ-732

**Approve** — the contract now accepts root `docs/run-prompts/`, still refuses every other new folder under root `docs/`, and the folder never ships.
Route A | merge `ce3cb5dc` (range `8ad193d8..ce3cb5dc`)

### What's built
- `_dev/tests/staged-skills-contract.sh` names `docs/run-prompts` in the export-ignore loop (line 75) and in the leftover-scan exclusions (line 107). `.gitattributes` line 57 export-ignores the folder.
- The six run prompts moved from `_dev/run-prompts/` to `docs/run-prompts/` as 100% renames. The old folder is gone from the main checkout.

### Decisions / risks for you
- One stale note outside the repository still tells agents the old rule (F2). It is the user's auto-memory, so the orchestrator or maintainer edits it, not a REQ.

### Findings

**Important:**
- F2 — `/Users/t2/.claude/projects/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/memory/run-prompts-not-in-root-docs.md` (whole file) and its index line in `MEMORY.md` of the same folder say "do not create `docs/run-prompts/`" and "write run prompts to `_dev/run-prompts/`". An agent that loads this memory acts on the old contract, which this REQ retired. The file is outside the repository and outside the write set, so the builder could not see it. — impact-rule-change → fix before archive (orchestrator rewrites or deletes the memory note and its index line; no repository change, no REQ)

**Minor:**
- F1 — `.gitattributes:19` is still 141 columns. Decision D-04 in the REQ says the builder reflowed the pre-existing long line. The reflow only moved the break: the old long line (about 118 columns) became a longer one. Comment-only, so behaviour is not affected, but the recorded decision does not match the diff. — impact-negligible → report only
- F3 — Non-archived queue records still name `_dev/run-prompts/` in their `*Source:*` lines: `do-work/queue/REQ-727-relation-map-second-opinion-preset.md:79`, `do-work/working/REQ-725-sort-local-memory-notes-by-destination.md:62`, `do-work/working/REQ-726-record-no-review-canvas-decision.md:67`, and REQ-732's own line 68. The requirement asks to repoint non-archived files. Triage decision D-01 chose to keep record text, and these lines quote the captured input, so the choice is reasonable. A builder that follows a Source line finds a missing folder, and `git log --follow` recovers the file. — impact-negligible → report only

**Nit:**
- F4 — `.gitattributes:57` puts the attribute one column right of its neighbours (`/docs/run-prompts` is 17 characters, the column is 16 + padding). Realigning every line for one entry would be churn. — impact-negligible → report only

### Requirements Checklist
- [x] Contract export-ignore assertion names `docs/run-prompts` — delivered (`_dev/tests/staged-skills-contract.sh:75`). `git check-attr` prints `docs/run-prompts: export-ignore: set`, and the assertion fails if the `.gitattributes` line is removed, because `check-attr` then prints `unspecified`.
- [x] Leftover-scan allows `docs/run-prompts/*` — delivered (line 107). The `-path "$repo_root/docs/run-prompts/*"` pattern does not match a sibling such as `docs/run-prompts-x/` or `docs/other/`, so other folders and loose files under root `docs/` still fail.
- [x] `.gitattributes` export-ignore line plus comment — delivered (lines 36-38, 57). The comment is three lines instead of the requested one, in the style of the neighbouring paragraphs. `git archive ce3cb5dc | tar -t` lists no run-prompts path.
- [x] Prompts moved with `git mv`, empty folder removed — delivered (six 100% renames, `ls _dev/run-prompts` reports no such folder).
- [x] Sweep of non-archived files — delivered outside `do-work/` (the only `_dev/run-prompts` mentions are inside two of the moved, historical prompts). Inside `do-work/` the Source lines keep their text per D-01 (F3).
- [x] Constraint: nothing under `skills/` or `suite/` changed — delivered (0 paths), so no release.
- [x] Constraint: contract keeps refusing other new files under root `docs/` — delivered (probe case 3 below).

### Acceptance Testing

**Result: Pass** (stages covered: GREEN probe re-run in this review, fast repository gate from the integrator; the `staged-skills` heavy lane is not yet recorded)
- `bash do-work/runs/work-2026-11-15-012821/REQ-732-probe.sh` → `REQ-732 probe PASS`, exit 0, at `ce3cb5dc`. The probe cuts the real contract's docs checks out by text markers and proves three facts: the real tree passes, a fixture with only `docs/run-prompts/probe.md` passes, and the same fixture plus `docs/other/probe.md` fails with the legacy-root message.
- The integrator's gate at the merge: `GATE-EXIT=0 WALL=125s`.
- `git archive ce3cb5dc | tar -t | grep -c run-prompts` → `0`.

### Suggested Additional Testing
- Confirm the `staged-skills` heavy lane passes in the drain; the REQ's `## Testing` points to a `## Heavy Verification Result` section that is not written yet.

### Scores (on the record — not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All delivered; F3 follows a stated triage decision |
| Code Quality | 92% | F1 long comment line and its inaccurate D-04 record |
| Test Adequacy | 90% | No repository test added; the contract itself is the test and the probe proves both directions |
| Scope | 98% | Write set matched exactly; a three-line comment where one line was asked |
| Risk | None | Comment, test list and file moves only |
| Acceptance | Pass | Probe and fast gate; heavy lane pending |

**Anti-bloat:** 0 unrequested helpers, options, files or tests. The probe is a run artifact, not a repository test.

### Follow-ups created
- None (4 findings report only or orchestrator-fixed; none is impact-critical)

## Review

**Overall: 95%** | 2026-11-15T02:39:27Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 98% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F2: the auto-memory note `run-prompts-not-in-root-docs.md` and its `MEMORY.md` index line (outside the repository) still say to use `_dev/run-prompts/` and never create `docs/run-prompts/` — impact-rule-change → fix before archive by the orchestrator (memory edit, no REQ)

**Minor findings:** F1 `.gitattributes:19` still 141 columns although D-04 claims the long line was reflowed — impact-negligible → report only; F3 non-archived queue `*Source:*` lines (REQ-727:79, REQ-725:62, REQ-726:67, REQ-732:68) still name `_dev/run-prompts/`, kept per D-01 — impact-negligible → report only; F4 (nit) `.gitattributes:57` attribute column one off — impact-negligible → report only
**Acceptance:** Pass — GREEN probe re-run at `ce3cb5dc` (allowed folder passes, `docs/other/` fails) plus the integrator's fast gate; `staged-skills` heavy lane not yet recorded
**Restatement sweep:** redefined the set of folders allowed under root `docs/` (kb, decisions, changelog, run-prompts) and the home of this repository's run prompts; grepped `docs/(kb|decisions|changelog)` enumerations and `_dev/run-prompts` outside and inside `do-work/`: `.gitattributes:17-18` and the contract comment at lines 71-74 agree, `skills/do-work-toolbox/actions/folder-structure.md:38` already lists `run-prompts/` under `docs/`, `CLAUDE.md` and `_dev/primes/` hold no restatement; stale: the out-of-repo auto-memory note (F2) and queue Source lines (F3)
**Anti-bloat:** 0 unrequested helpers, options, files or tests
**Suggested testing:** 1 items
**Follow-ups created:** None (3 findings report only, F2 is an orchestrator memory edit)

*Reviewed by review-work action*
