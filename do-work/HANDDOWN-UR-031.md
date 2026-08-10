# HANDDOWN — UR-031: Four-Skill Distribution Follow-ups

**Written:** 2026-08-10T12:04:18Z, at v0.186.14 after REQ-162 metadata commit `2a6e3a2`.
**Stop condition:** The approved REQ-161/REQ-162 serial run is complete. REQ-163 is the sole unresolved member and remains `pending-answers`; do not claim or implement it until the user clarifies it.

## Completed in this session

| REQ | Outcome | Review | Version | Implementation | Metadata |
|---|---|---:|---:|---|---|
| REQ-161 | Escaped-link delimiter and list-paragraph classification | 63%, Partial | 0.186.13 | `ad3f8bd` | `a06f27f` |
| REQ-162 | Ordinary multiline-backtick Just collision state | 100%, Pass | 0.186.14 | `aff7c9c` | `2a6e3a2` |

Each REQ completed its normal claim, fresh Plan/Explore/Builder/Review contexts, RED/GREEN, owner qualification and testing, Lessons/Orientation, archive, patch release, implementation commit, guarded `commit:` write-back, and metadata commit.

## What was delivered

### REQ-161

- The shipped-reference classifier now masks complete inline-link-shaped regions when a closing bracket or destination-opening parenthesis has odd escape parity, preventing hidden first-party URLs from reappearing through fallback discovery.
- Nonempty bullet and one-to-nine-digit ordered list items now establish paragraph state, preserving live four-column continuation links while retaining empty-marker, blank, fence, and genuine-code controls.
- Exact fixtures preserve source length/newline offsets and full publication/topology/path policy.
- Independent review passed the four approved cases but found adjacent relative-target parity, escaped-label-content, and list-fence-info-string gaps. Those are consolidated in consent-gated REQ-163.

### REQ-162

- Root and shipped `replace-text-section.sh` helpers now persist ordinary single-backtick command state across physical lines through the existing raw active-delimiter path.
- The next literal backtick closes even after a backslash, same-line close/reopen remains valid, and exact triple-backticks retain longest-first handling.
- Just-parseable positives, exact insertion bytes, exact four-name real-collision diagnostics, and byte-preserving rejection prove both acceptance and safety.
- Full review passed with no finding or follow-up.

## Live queue at handdown

| REQ | Status | Dependency / action |
|---|---|---|
| REQ-163 | `pending-answers` | Generation≥2 review sweep from REQ-161. One consent question covers the remaining inline-link parity/label-context and list-fence classification variants. Run `do-work clarify`; do not auto-claim it. |

`do-work/working/` contains no REQ. UR-031 remains open with 26 completed members and one unresolved follow-up; completed REQs remain loose in `do-work/archive/` until REQ-163 is terminally resolved.

## REQ-163 decision boundary

- Even two/four-backslash destination-opening-parenthesis forms can lose relative targets even though first-party URLs appear live through the bare-URL fallback.
- An escaped opening bracket inside an otherwise live label can cause the containing link to be falsely masked.
- Bullet and one-to-nine-digit ordered list items can miss backtick/tilde fences with attached info strings, including nested tilde fences, and scan their links as published Markdown.
- These are concrete, independently reproduced classifier gaps, but generation-depth policy requires explicit user consent before implementation. Publication, topology, target-resolution, raw/blob, and path policy remain out of scope.

## Verification evidence

- REQ-161 RED reproduced its four approved defects; GREEN passed focused shipped-reference, full contract, staged-skills, suite-manifest, Bash/ShellCheck, changelog identity, offset invariants, and protected-file checks. Review ran a 62-case actual-helper matrix and captured only the REQ-163 variants.
- REQ-162 RED reproduced exact false collisions `kanban-summary, run-kanban` and a literal-only `kanban-static` overreport while preserving bytes; GREEN passed exact insertion/diagnostic checks through both production helpers and Just 1.46.
- Post-v0.186.14 `bash _dev/tests/contract-regressions.sh` passes, including suite manifest, shipped-reference, record-hash/blanked-REQ guards, updater, staged-skills, and installer probes.
- Standalone installer/staged-skills/suite-manifest contracts, paired-helper and changelog mirrors, version mirrors, board-template Just parse, Bash syntax, warning-level ShellCheck, and `git diff --check` pass.
- Provenance read-back records exactly `ad3f8bd` and `aff7c9c` in archived REQ-161/162.

## Cleanup result

- Passes 0–6 completed. No terminal REQ is stranded in queue/working, no misplaced request tree/archive folder exists, no consumed run scratch exists, no orphan `worktree-agent-*` worktree/branch exists, and the blanked-REQ scanner found no blanked or unparseable REQ/UR file.
- UR membership was derived from `user_request:` across queue/working/archive: 26 completed and REQ-163 `pending-answers`. Cleanup correctly left UR-031 open and made no structural move.
- No cleanup move occurred, so no documentation link repoint or cleanup-only commit was required.

## Preserved approved files

- `decisions/records/adr-019-four-skill-suite-contract.md` — SHA-256 `2d5a54bc9435f8643f4d30e332c37426fbacb15442503eba338cb1f9ab11b282`.
- `do-work/user-requests/UR-031/input.md` — SHA-256 `ed156a18dc11f4f367e80d0e1cca8dbf676dffaae8030622df214e8070bab160`.

## Copy-paste restart prompt

```text
do-work clarify

Resume UR-031 from /Users/t2/Desktop/e1-experimental-repos/skill-do-work2. Read do-work/CHECKPOINT.md and do-work/HANDDOWN-UR-031.md first. REQ-163 is the sole unresolved member and remains pending-answers; clarify its consent question before any do-work run claims it. Preserve decisions/records/adr-019-four-skill-suite-contract.md and do-work/user-requests/UR-031/input.md byte-identically.
```
