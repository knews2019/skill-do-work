## Review: REQ-686

**Approve**: both requested lines are present and true against the code, the CHANGELOG item is correctly dropped, and every changed line is a comment or doc text.
Route A | 93198e39 (range 577a4ac7..93198e39, 4 files, +8/-6)

### What's built
- `lessons-do-work-cli.md:121` now says what to do after a performed and verified `git revert` of the primary commit: delete the journal and its payloads folder by hand, because `--discard-journal` only accepts phase `prepared`.
- `heavy_commands.go:24` now says the three fast-stage commands serve `_dev/tests/maintainer-verify.sh` only and no action calls them.
- Two integrator addenda (comments only) fix earlier review findings: the HEAD-sentinel comment in `git_transaction.go:704-705` (REQ-685 review F1) and two stale comments in `timeline_scroll_browser_probe_test.go` (REQ-682 review F2).

### Decisions / risks for you
- None. F1 below is a one-phrase wording fix you may want in a later doc pass.

### Findings

**Important:**
- None.

**Minor:**
- F1. `lessons-do-work-cli.md:121` (wave-end sweep, element d): the bullet still says `--discard-journal` "is the only supported exit from a journal that refuses on every replay", and the new next sentence gives a second exit (hand delete after a verified revert) for exactly such a journal, since `FINALIZATION-PRIMARY-COMMIT` (`finalization_apply.go:114-115`) refuses on every replay once `PrimaryCommit` is set. A reader gets two statements that disagree. Fix: "the only supported exit for a journal still in phase `prepared`". `finalization_commands.go:121-127` ("the documented exit", "no longer the only way out") still reads true. impact-negligible → report only
- F2. Anti-bloat count. The diff adds nothing the REQ did not name except the two coordinator addenda: one comment pair in `gittransaction/git_transaction.go` and two comments in `timeline_scroll_browser_probe_test.go`. No helper, flag, option, file, section or test was added. The `git_transaction.go` edit is against the REQ's own constraint "This REQ must not edit `gittransaction/` files; REQ-685 owns them". REQ-685 is already archived, the coordinator ordered the edit, and `## Implementation Summary` and `## Qualification` record it, so it is traced, not drift. impact-negligible → report only

**Nit:**
- F3. `git_transaction.go:704-705`: the caller list is a closed enumeration. The three named files do test `CommitSHA != ""` and report the commit (`cleanup_apply.go:144`, `doctor_repair.go:173-174`, `publication_commands.go:234-235`), but the package's own `survivingChanges` (`transaction_findings.go:171`) also reads it, and `finalization_apply.go:142` tests it for its own transaction calls. "tells callers (for example ...)" would keep it true as callers change. impact-negligible → report only
- F4. `lessons-do-work-cli.md:121`: the literal path `.git/do-work-finalization/REQ-N.json` is right for the main checkout only. The code resolves the folder with `git rev-parse --git-path do-work-finalization` (`finalization_journal.go:102`), which in a linked worktree points under `.git/worktrees/<name>/`. The sentence also spells the id `REQ-N` where the sentence before it uses `REQ-NNN`. impact-negligible → report only
- F5. `heavy_commands.go:24`: the comment sits directly above `func Handlers()`, so Go doc tools show it as the doc comment of `Handlers`. The builder chose this to avoid a gofmt realignment of the const block (D-01). The text is true. impact-negligible → report only

### Requirements Checklist

- [x] (a) One recovery sentence beside the `--discard-journal` sentence in `lessons-do-work-cli.md`. Delivered. It is true: `discardPreparedJournal` (`finalization_commands.go:146-148`) refuses any `journal.Phase != PhasePrepared`, and the journal and payloads paths match `journalLocations` (`finalization_journal.go:102-114`). No `actions/*.md` file mentions `FINALIZATION-PRIMARY-COMMIT` (grep), so the lessons bullet is the one home.
- [x] (c) One comment line at the fast-stage registration in `heavy_commands.go`. Delivered. It is true: outside `do-work/archive/`, `do-work/runs/` and historical or third-party records, the only callers of `decide-fast-stage`, `record-fast-stage` and `invalidate-fast-stage` are `_dev/tests/maintainer-verify.sh:160, 186, 199`. No action, script or hook under `skills/` calls them.
- [x] (b) CHANGELOG 0.305.46 wording not changed, with the reason recorded in `## Implementation Summary`. Delivered (CHANGELOG is not in the diff).
- [x] Constraint "documentation and comment lines only, no behavior change". Confirmed: every changed non-Markdown line in the range is a `//` comment line. `gofmt -l` on the three Go files is empty.
- [x] Addendum 1 (HEAD-sentinel comment): true. `committedRisk` (`git_transaction.go:1444-1447`) sets `CommitSHA` to "HEAD", and the three named callers report "committed in ..." or the commit evidence when it is non-empty.
- [x] Addendum 2 (probe comments): true. `buildTimelineScrollProbeSite` has three callers (`timeline_scroll_browser_probe_test.go:167, 768, 842`), and the first-visible-child measurement is asserted in `TestBrowserBehaviorTimelineViewHasOneScrollSurface` (`:261-263`).

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- Implementation: each new sentence and comment was checked against the code it describes (see the checklist).
- Integration: the integrator's full gate at `93198e39` exited 0 (`## Testing`, one run). I re-ran `gofmt -l` on the three Go files: clean. I did not run the repository gate (read-only reviewer brief).

### Restatement sweep (wave end, UR-152 run work-2026-10-10-100748)

This diff redefines the recovery story for a stuck finalization journal (element d) and the consumer list in the HEAD-sentinel comment. Inherited elements come from the `**Restatement sweep:**` lines of REQ-679 to REQ-685 in `do-work/archive/`. REQ-679 and REQ-680 recorded "nothing redefined", so every sibling has a recorded line and none is unread. Searched `skills/`, `_dev/`, `README.md` and `do-work/lessons-index.md`. CHANGELOG past entries, archives, run files, the inbox report and `ai-reports/` were treated as history.

- a. QUALIFY-SUMMARY-MISSING (REQ-684, qualify now says which summary problem it found). Terms: `QUALIFY-SUMMARY`, `missing or empty`, `summary section not found`, `lists no backticked`, `no claimed files`, `backticked file path`. Found: `checks.go:291-297` (new text), `commands.go:398` (keys on the code only), `checks_test.go:160, 176`, and `_dev/tests/contracts/core-checks.sh:833` ("shared principles file missing or empty", unrelated). No action, doc or crew file describes when the finding fires. Not stale.
- b. appendSectionEntry section lookup (REQ-681, now uses the shared visible-section reader). Terms: `appendSectionEntry`, `sectionLineBounds`, `duplicate heading`, `second ... heading`, `## Blocked`, `## Cancelled`, `In Progress (interrupted)`. Found: code and test only. Prose in `work-reference.md`, `abandon.md`, `stakeholder-answers.md`, `cleanup.md` and `fan-out-reference.md:168` names the sections but states no exact-line matching rule. Lessons lines 7 and 28 are history about other readers. Not stale.
- c. Rollback uses the open root (REQ-685). Terms: `reopens the root`, `no-handle rollback`, `REQ-598`, `rollbackWithoutRoot`, `no-root rollback`, `rollback root`, `rollback` in `lessons-do-work-cli.md` and `prime-do-work-cli.md`. Found: `_dev/tests/audit-lockins.sh:599-628`, which already describes the held root passed to `rollbackFailure`. The lessons file has no root wording (`:128` says only that gittransaction owns rollback). Not stale.
- d. "--discard-journal is the only supported exit" (REQ-686). Terms: `only supported exit`, `documented exit`, `discard-journal`, `by hand`, `FINALIZATION-PRIMARY-COMMIT`. Found: `lessons-do-work-cli.md:121` (stale against its own new next sentence, F1), `finalization_commands.go:121-127` (still true), the flag-parser errors at `:130, 197-211` (true). No action file restates it. One stale statement: F1.
- e. Interview cadence clock format (REQ-683, am/pm accepted, stored value still 24-hour). Terms: `HH:MM`, `needed_by`, `12-hour`, `24-hour`, `AM`/`PM`, `am/pm`, `standing_slots`. Found: `work-operating-model.md:104` (`needed_by` is a timing window) and `:405` (`time` = `HH:MM` or `""`), both still true. The other `HH:MM` hits are the Timestamp rule, memory log headings and interview log headings, which are unrelated formats. Not stale.
- f. Hidden Timeline scroll behavior (REQ-682). Terms: Timeline with scroll, listener, hidden, redraw, view switch in `skills/do-work-board/**/*.md`, `_dev/primes/prime-kanban-board.md` and `lessons-kanban-board.md`. Found: `lessons-kanban-board.md:41` (REQ-227, a view owns its listener teardown; still true, `board-timeline.js:159` removes listeners) and `:78` (REQ-682's own lesson). The code comment at `board-timeline.js:2140-2144` matches the new guard. No board doc describes scroll listener lifetime. Not stale.
- Sibling report-only findings were re-checked, not re-argued: REQ-685 F1 and REQ-682 F2 are fixed by this merge's addenda. The other report-only findings (REQ-679 F1-F3, REQ-680 nits, REQ-681 F1-F3, REQ-682 F1/F3/F4, REQ-683 F1-F2, REQ-684 nit) are untouched by this diff and still stand as recorded.

### Suggested Additional Testing

- Deployment: unassessed. After the release, check that the installed mirror of `lessons-do-work-cli.md` carries the new sentence (the release prime's mirror check covers this).

### Scores (on the record — not the headline)

**Overall: 97%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Both lines delivered and true; (b) correctly dropped |
| Code Quality | 95% | One sentence now disagrees with the one before it (F1) |
| Test Adequacy | N/A | Documentation and comment lines only |
| Scope | 95% | Two comment-only addenda outside the write set, both coordinator-ordered and recorded (F2) |
| Risk | None | No code behavior change |
| Acceptance | Pass | Implementation and integration stages |

### Follow-ups created
None (5 findings report only)

## Review

**Overall: 97%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1. `lessons-do-work-cli.md:121` still calls `--discard-journal` "the only supported exit from a journal that refuses on every replay" while the new next sentence gives a hand-delete exit for such a journal; reword to "for a journal still in phase `prepared`" — impact-negligible → report only
- F2. Anti-bloat: nothing added beyond the REQ except two coordinator-ordered comment-only addenda (`git_transaction.go`, `timeline_scroll_browser_probe_test.go`); the first is against the REQ's "must not edit gittransaction/" constraint but is recorded and traced — impact-negligible → report only
- F3 (nit). `git_transaction.go:704-705` lists three `CommitSHA` callers as if complete; `transaction_findings.go:171` and `finalization_apply.go:142` also read it — impact-negligible → report only
- F4 (nit). `lessons-do-work-cli.md:121` hard-codes `.git/do-work-finalization/`, true only in the main checkout (code uses `git rev-parse --git-path`), and spells `REQ-N` beside `REQ-NNN` — impact-negligible → report only
- F5 (nit). `heavy_commands.go:24` comment becomes the Go doc comment of `Handlers()` (placement chosen in D-01) — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: each new line checked against the code it describes; integrator gate exit 0 at 93198e39; gofmt clean.
**Restatement sweep:** redefined the stuck-journal recovery wording (`--discard-journal` exit plus hand delete after a verified revert) and the HEAD-sentinel comment's caller list; wave end swept inherited elements a-f (QUALIFY-SUMMARY-MISSING text, appendSectionEntry lookup, rollback open root and REQ-598 wording, discard-journal "only supported exit", interview HH:MM, hidden Timeline scroll); one stale statement found (F1); no sibling unread.
**Suggested testing:** 1 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*
