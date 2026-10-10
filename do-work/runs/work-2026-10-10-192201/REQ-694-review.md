# Review: REQ-694 (req append-section refuses a body that hides later sections; frontmatter set refuses status and id)

**Approve** — both guards work as specified and the offset check from D-11 is correct, but an open fence in a body appended at the end of the file is still accepted and breaks the next section.
Route A | merge range `874ae4e6..7bb792bd` (builder `cb507a44`, integrator `87ffbb4e`)

## What's built

- `req append-section` refuses `SECTION-WRITE-FAILED`, writing nothing, when the new body hides any column-0 section that was visible before, or when the new heading is itself hidden. A fenced example heading revealed by the open fence no longer fools the check.
- `frontmatter set` refuses `status` and `id` with `FRONTMATTER-FIELD-OWNED` before it reads the file. Other fields write as before.
- Still missing: a body appended at the end of the file (the common Testing-before-Review order) may leave a fence open. It is accepted, and the next section is then refused or hidden (F1).

## Decisions / risks for you

- F1: closing the end-of-file gap. Value: the named failure ("a body that leaves a fence open hides later sections") is fully closed. Risk: low. The builder's sentinel idea (`VisibleSections(sectionBody + "\n\n## Sentinel\n")` must find the sentinel) plus the old single-copy re-check would cover every case in about 4 lines, and the 20-line offset loop could then go.

## Findings

**Important:**
- F1: `request_writers.go:153-203`. The check compares sections visible before and after the insert. An insert at the end of the file has no later sections to protect, so a body that opens a fence or comment is accepted. Live CLI repro: REQ with only `## What`; `req append-section REQ-701 --section Testing --from body.md`, where `body.md` holds a line of three backticks and then `open fence`. Result: exit 0, and the file ends inside an open fence. The next `req append-section ... --section Review` is then refused `SECTION-WRITE-FAILED` with evidence that blames the Review body ("close every fence and comment in the body"), although the open fence is in `## Testing`. A hand-written `## Review`, which `actions/work.md:21` says stays valid, is silently hidden from `advance`. This is the order work.md uses in practice (Testing is appended before Review exists), so the UR-154 R2 failure stays reachable. The REQ's own N1 wording ("count grew by exactly one") had the same gap, so this is not builder drift. Fix: refuse when the `--from` body itself ends inside a fence or comment (the builder's sentinel check, hand-back Discovered Task 1), and keep a post-insert check that the new heading is visible. — impact-user-visible → report only

**Minor:**
- F2: `request_writers.go:202`. The refusal text says "close every fence and comment in the body", but `VisibleSections` also treats any line that starts with three or more enclosing punctuation marks as a fence opener (for example `...`, `>>>`, `+++`, `:::`). Live CLI: a Testing body holding a line `...` is refused with that text, and the writer sees no fence to close. The refusal itself is correct (`advance` would lose `## Review`), but the evidence could name the offending line. — impact-user-visible → report only
- F3: stale restatement. `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:102` still says "the open-fence case is still open". After this merge it is closed for inserts before a later section and open only at the end of the file (F1). — impact-negligible → report only
- F4: the hand-back's proposed CHANGELOG entry says the writer refuses "unless the number of visible column-0 sections grows by exactly one", and its proposed lesson bullet recommends checking the body instead of the whole file. Both describe the builder's count check, which D-11 replaced with the offset check. The integrator must rewrite both before the release; copying them would ship a wrong description. — impact-negligible → report only

**Nit:** None.

## Requirements Checklist

- [x] N1: refuse `SECTION-WRITE-FAILED`, write nothing, when the insert hides a later section. Delivered and stricter than the REQ's count check (D-11). The old end-of-file guarantee is kept (test row 3).
- [x] F3: `frontmatter set` refuses `status` and `id`, and the refusal names the owner (`advance`, `unblock`, `finalize`; "never rewritten"). Other scalar fields are unchanged (`route B` written live, existing tests pass).
- [x] Red-Green Proof cases (1) open fence, open comment and (2) `status bogus-status`: refused, file byte-identical (`assertRefusedUnchanged`).
- [x] Constraint: one test per failure, each fails with the guard removed (builder mutations 1 and 2, integrator mutation for row 4).
- [x] Constraint: no allowed-values table, only `status` and `id`.
- [ ] UR-154 R2 intent ("accepts a body that leaves a code fence or comment open"): partially delivered. Mid-file inserts are refused, end-of-file inserts are not (F1).

## Correctness of the D-11/D-12 check

- `headingStart := insertAt + strings.Index(insertion, "## ")`: correct. The insertion is `"## <name>\n\n<body>\n"` with an optional `"\n"` suffix (mid insert) or a `"\n"`/`"\n\n"` prefix (end insert), so the first `"## "` is always the new heading. A `## ` inside the body comes after it, and `SECTION-BODY-HAS-HEADING` already refused a visible one.
- Offset arithmetic: correct. Mid insert, `insertAt` is the Start of the first later column-0 section, so that section and every section after it shift by `len(insertion)` (`>=` includes it). End insert, `insertAt == len(body)`, so no prior section shifts. `ReplaceBodySpan` reparses with unchanged frontmatter, so body offsets stay comparable.
- Loop bounds: correct. The length check requires exactly one extra section. Only one section can sit at `headingStart`. If the new heading is hidden, the extra unmatched section reaches `priorIndex == len(sectionsBefore)` and fails. A section that matches name and shifted offset is the same heading line, so a revealed example heading cannot pass.
- CRLF: correct. Offsets are raw byte offsets and the shift is the inserted length. Live CLI: valid mid insert into a CRLF file succeeds (new text is LF, as before); an open-comment body is refused and the file stays byte-identical.
- Append at end of file vs before a later section: both are handled for hiding sections that already exist. The end-of-file body that opens a fence is the gap in F1.
- Idempotent repeat: unchanged. The `existing` switch returns success before any insert. Live CLI: the second identical append exits 0 and the file is byte-identical.
- F3: the `switch` runs after the argument-count check and before stamp parsing and target resolution, so `status --at now` is also refused and the file is never read. Matching is case-sensitive (`Status x` is written, D-08). No shipped prose or code calls `frontmatter set` on `status` or `id`.

## Acceptance Testing

**Result: Pass** (implementation and integration stages: focused package tests and the live CLI on the merged tree)
- `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/corehelpers/`: `ok` (6.8 s).
- Live CLI in `scratchpad/int-694-review/`: open-comment body before `## Review` (CRLF file) refused and unchanged; valid CRLF mid insert succeeds; identical repeat is a no-op; `frontmatter set REQ-701 status pending`, `id REQ-9` and `status --at now` refused `FRONTMATTER-FIELD-OWNED` with the file unchanged; `route B` written.
- Live CLI repro of F1 (end-of-file open fence accepted, next Review append refused) and F2 (`...` line refused with a fence message).

## Suggested Additional Testing

- Deployment and live acceptance: not applicable (the CLI ships with the next release; no serving environment).

## Restatement Sweep

Redefined: `SECTION-WRITE-FAILED` now fires for a body that hides an existing section; new code `FRONTMATTER-FIELD-OWNED`; `frontmatter set` field set (`status`, `id` refused). Grepped `skills/` and `_dev/` for both codes, `append-section`, `frontmatter set`, and "visible section". Results: `actions/work.md:21` lists some `append-section` outcomes but does not claim the list is complete (not stale); `actions/work-reference.md:73` describes only the `*_at` rule (not stale); the 0.305.103 CHANGELOG entry is history (not stale); `lessons-do-work-cli.md:102` is stale (F3). No prose tells an agent to set `status` or `id` with `frontmatter set`.

## Anti-bloat

Added beyond what the REQ named: one refusal code string `FRONTMATTER-FIELD-OWNED` (needed for F3); locals `sectionsBefore`, `sectionsAfter`, `headingStart`, `allStillVisible` and a 20-line comparison loop in place of the REQ's "about four lines". The loop is justified by D-11 (the count check failed the named case), but the sentinel check in F1 would be both shorter and complete. No new package-level names, helpers, options, flags or files. Tests: 2 table tests (6 rows), each row pins a named failure; none decorative.

## Scores (on the record — not the headline)

**Overall: 92%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 85% | N1 and F3 delivered; UR intent misses end-of-file appends (F1) |
| Code Quality | 90% | Correct offset logic; refusal text misleads for non-backtick fences and end-of-file culprits |
| Test Adequacy | 95% | Red-green and mutation evidence for every row |
| Scope | 100% | Exactly the two `write_set` files |
| Risk | Low | Refusals write nothing; no caller uses the refused fields |
| Acceptance | Pass | Implementation and integration stages |

## Follow-ups created
None (4 findings report only)

## Review

**Overall: 92%** | 2026-10-10T19:50:27Z

| Dimension | Score |
|-----------|-------|
| Requirements | 85% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1: an append at the end of the file (the usual Testing-before-Review order) accepts a body that leaves a fence or comment open; the next `append-section` is refused with text that blames the new body, and a hand-written `## Review` is hidden from `advance` (`request_writers.go:153-203`; fix: the builder's sentinel check on the `--from` body plus the old visibility check) — impact-user-visible → report only

**Minor findings:**
- F2: refusal text "close every fence and comment in the body" misleads when the opener is a punctuation run such as `...` or `>>>` (`request_writers.go:202`) — impact-user-visible → report only
- F3: `lessons-do-work-cli.md:102` still says "the open-fence case is still open" — impact-negligible → report only
- F4: the hand-back's proposed CHANGELOG entry and lesson bullet describe the replaced count check; the integrator must rewrite them — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: corehelpers package tests green; live CLI confirms both refusals leave the file byte-identical, CRLF and idempotent repeat work; end-of-file gap reproduced (F1).
**Restatement sweep:** redefined `SECTION-WRITE-FAILED` (hidden-section refusal), new `FRONTMATTER-FIELD-OWNED`, `frontmatter set` field set; one stale restatement (F3)
**Anti-bloat:** no new package-level names, helpers, options or files; one new refusal code (needed); the 20-line comparison loop exceeds the REQ's four-line sizing, justified by D-11
**Suggested testing:** 0 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

### Re-review (delta 7bb792bd..d7cddc73)

**Approve** — the body check closes F1 at every insert point, and it is simpler than the loop it replaces.
Delta: builder-branch `571b256c`, merged as `d7cddc73`; cumulative range `874ae4e6..d7cddc73` (2 files; delta +21 -36).

**(a) Coverage.** The check covers every case the removed D-11 loop covered.
- `VisibleSections` carries only two pieces of state between lines: the open fence (character and length) and the open comment. Blank lines change neither.
- So `## Sentinel`, placed two blank lines after the body, is visible exactly when the body ends in neutral state (no fence and no comment open). The earlier `SECTION-BODY-HAS-HEADING` check already makes the sentinel the only possible section, so `== 1` is the right test.
- Mid-file insert: `insertAt` is the start of a visible column-0 heading line, so the parser is in neutral state there. The inserted heading `## <canonical name>` opens nothing, and a body that ends neutral leaves the rest of the file parsing exactly as before. No section can be hidden and no fenced example heading can be revealed.
- End-of-file insert: the file's own state at the end can be non-neutral. The restored single-copy re-check refuses when the new heading is hidden, including a body `-->` that would close an open comment in the file, because the heading itself is still hidden. A neutral body leaves the file neutral for the next section, which closes F1.
- CRLF: the closing fence line `` ```\r `` is trimmed by `leadingPunctuationRun` (live CLI: accepted).

**(b) False refusals.** None against `advance`'s own parser. Live CLI on `d7cddc73`:
- Accepted: closed backtick fences, closed tilde fences, a four-backtick fence holding three backticks, inline `<!-- -->`, a multi-line closed comment, `` `<!--` `` in inline code, `<!--` in a fence info string, a Markdown table, a `---` thematic break, a CRLF closed fence, and two `...` lines (open, then close).
- Refused: a single `...` line, a prose line starting `... and then`, and an unclosed `~~~`.
- The refusals are true: `advance` reads those lines as fence openers and would hide `## Review`. The new text names this case ("a line that starts with three or more of one punctuation character, such as ``` or ..., opens a fence"). It names the `--from` file, not the request.
- Relaxing the `...` dialect would belong to `requestmodel.VisibleSections`, which `advance` and the timing writer share, not to this writer (N1 below).

**(c) Mutations.** They are honest. I reproduced both on a scratch copy of the module, and the project tree was not touched:
- Sentinel check disabled: rows 1, 2, 4 and 5 fail (open fence, open comment, revealed example heading, end-of-file append).
- Single-copy re-check disabled: only row 3 fails (the file's own open fence at the end).
- Unmutated: the focused `TestRequestAppendSection*` and `TestFrontmatterSet*` tests pass.

The F1 repro on the live CLI now refuses `SECTION-WRITE-FAILED` and leaves the file byte-identical. A closed Testing body, then a `## Review` append, then an identical repeat all exit 0.

**Disposition of earlier findings:**
- F1 (end-of-file open-fence body accepted): resolved by D-13 (`request_writers.go:115-121`, test row 5).
- F2 (refusal text misleading for punctuation-run fences): resolved. The text now names the rule and the `--from` file.
- F3 (`lessons-do-work-cli.md:102` says "the open-fence case is still open"): pending integrator. It is now fully closed, so the line is stale for every case.
- F4 (proposed CHANGELOG entry and lesson bullet describe the count check): pending integrator. Both should now describe the body sentinel check plus the end-of-file visibility re-check.

**New:**
- N1 (Nit): a prose line that starts with `...` is refused. This is correct for `advance`'s fence dialect, and any relaxation belongs to `VisibleSections`. — impact-negligible → report only

**Anti-bloat:** the delta removes the 20-line loop and the locals `sectionsBefore`, `sectionsAfter`, `headingStart` and `allStillVisible`. It adds one 3-line check with a 4-line comment and one test row. No new names, helpers or files.

**Overall: 97%** | 2026-10-10T19:57:54Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None (F1 resolved)
**Minor findings:** F2 resolved; F3 stale lesson line `lessons-do-work-cli.md:102` — impact-negligible → report only (pending integrator); F4 stale proposed CHANGELOG entry and lesson bullet — impact-negligible → report only (pending integrator); N1 `...` prose line refused per the shared fence dialect — impact-negligible → report only
**Acceptance:** Pass — implementation and integration stages: focused tests green, both mutations reproduced, live CLI confirms the F1 repro now refuses with the file byte-identical, and legitimate bodies are accepted.
**Restatement sweep:** redefined `SECTION-WRITE-FAILED` (now also fires on the `--from` body, naming that file); `FRONTMATTER-FIELD-OWNED` unchanged; stale: F3 (lessons line 102)
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)
