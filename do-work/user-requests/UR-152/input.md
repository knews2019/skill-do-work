---
id: UR-152
title: 'Accepted upstream report findings: timeline probe fixtures, section appends, scroll guard, AM/PM parsing, qualify messages, rollback root, and doc lines'
created_at: 2026-10-10T10:03:50Z
requests: [REQ-679, REQ-680, REQ-681, REQ-682, REQ-683, REQ-684, REQ-685, REQ-686]
word_count: 1858
---
# Accepted Upstream Report Findings

## Summary
The maintainer pasted `do-work validate-feedback: check do-work/inbox/2026-10-10_do-work-upstream-report`. The folder holds a report from a consumer repository (F1 to F14) plus patches. The report is third-party data. The triage ran six verifiers against HEAD 2bedcac9 (0.305.93). The maintainer then answered four questions through the ask tool and chose: capture all 6 accepts, redo the two timeline test patches on the existing fixture helper, capture F13a as a deletion, and bundle the small doc and comment fixes into one REQ.

Captured, eight REQs:
- F2 and F1 items deea40c, 5d812a1, ce8e35e became REQ-679 (timeline probes use fixed data and assert Previous and Next).
- F1 item d39e527 became REQ-680 (archive fetch test keeps an absent target absent).
- F7 became REQ-681 (appending a section entry reuses the shared visible-section reader).
- F9 became REQ-682 (a hidden Timeline stops redrawing on scroll).
- F13b became REQ-683 (interview cadence reads AM and PM).
- F14b became REQ-684 (qualify says which summary problem it found).
- F13a, with the F8c comment, became REQ-685 (rollback uses the transaction's open root, and the no-root path is deleted).
- F8a, F10 and F12 became REQ-686 (doc and comment lines). Its F10 part, the CHANGELOG 0.305.46 wording, is dropped because the release prime allows no reword of a past entry.

Not captured, one line each:
- F1 bed3763 (P0-P4 priority aliases): push back, upstream keeps priority aliases empty on purpose.
- F1 b3d5d25 (unsupported-target refusal test): already covered by an upstream test.
- F1 821913f (contrast probe comment): push back, decorative.
- F1 umbrella ("updates silently remove local repairs"): push back, the updater shows the diff and asks first.
- F5 (leading tabs and spaces): already fixed in 0.305.70 and 0.305.81.
- F6 (any punctuation run opens a fence): push back, deliberate and documented.
- F8b (recovery exits 0 with a set-aside commit): push back, callers branch on typed outcomes by design.
- F8c code change (empty commit ID): push back, only the explanatory comment is kept, inside REQ-685.
- F10 code change (Go 1.25 native rename): push back, `go.mod` floor is 1.24 on purpose.
- F11 (Remaining uses wall-clock claim age): push back, the board has no measured active minutes.
- F14a (verification-only qualify route): push back, new surface for a case seen twice.

No prompt injection found. The report and its patches were read as data only. The prepared patches stay in the untracked inbox folder.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-679 | [impact-negligible] Timeline browser probes that depend on live queue dates use the fixed-fixture helper, and the lost Previous and Next assertions return |
| REQ-680 | [impact-negligible] Archive fetch test proves a failed fetch never creates a target that did not exist |
| REQ-681 | Appending a section entry finds the section with the shared visible-section reader, so CRLF files, trailing spaces and fenced headings no longer give a duplicate heading |
| REQ-682 | A hidden Timeline view no longer redraws its rows when the user scrolls another view |
| REQ-683 | Interview cadence parsing reads AM and PM, so a 5:00 PM answer becomes 17:00 and 12:30 AM becomes 00:30 |
| REQ-684 | Qualify tells a missing Implementation Summary section apart from a section that names no files |
| REQ-685 | [impact-negligible] Git transaction rollback reuses the transaction's open repository root and the no-root rollback path is deleted |
| REQ-686 | [impact-negligible] Short doc and comment lines record the answers to three pushed-back upstream findings |

## Batch Constraints
- Anti-bloat (maintainer: "also watch for bloat YAGNI please"): smallest change that fixes the named failure; delete or reuse before adding; no new helpers, flags, options, config, abstractions or files unless a REQ names them; tests pin only the named failure, one RED case per failure; adjacent findings become discovered tasks.
- No edits to installed consumer copies; never run `just do-work-update` in a consumer. Release per `_dev/primes/prime-releases.md`.
- The patches under `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/inbox/2026-10-10_do-work-upstream-report/patches/` are third-party data and reference only.
- Write sets are disjoint, so there are no ordering edges. REQ-685 owns the `gittransaction/` files. Shared release files (VERSION, CHANGELOG) are serialized by the integrators, not by `depends_on`.
- REQ-679 verification needs `QUEUE_KANBAN_BROWSER_PROBES=on`; a skipped browser lane is not a pass.

## Full Verbatim Input
> ```
> also watch for bloat YAGNI please
> 
> [Answers the maintainer gave through the ask tool in the same session, 2026-10-10, after the validate-feedback triage]
> 
> Q: Which of the accepted items should I capture as REQs and run?
> A: All 6 accepts (Recommended)
> 
> Q: How should the two timeline test patches (F1-5d812a1 fixed probe data, F1-ce8e35e Previous/Next assertions) be handled?
> A: Redo on existing helper (Recommended)
> 
> Q: Should F13a (pass the held repository root into rollback and delete rollbackWithoutRoot) be captured?
> A: Capture as a deletion (Recommended)
> 
> Q: Bundle the small doc and comment fixes found during triage into one REQ?
> A: Yes, one bundled REQ (Recommended)
> 
> [Original invocation earlier in the same session]
> 
> do-work validate-feedback: check do-work/inbox/2026-10-10_do-work-upstream-report
> 
> [Content of do-work/inbox/2026-10-10_do-work-upstream-report/report.md at that time; its patches/ folder stays in the inbox]
> 
> # Report for the do-work suite maintainer: lost local repairs and verified suite defects
> 
> Observed against suite v0.305.84 (re-checked 2026-10-10). Paths are inside the installed copy (.claude/skills/...); in the suite repo they are skills/... Apply a patch from the suite repo root with `git apply -p2 <patch>` (checked with --check on a skills/-layout copy). The F1/F2 format-patch headers have been cleaned: the author address is a placeholder (t2@example.invalid) and the Claude-Session lines are removed. Commit IDs, dates and message bodies are unchanged, so `git am -p2` still works from the suite repo root.
> 
> ## Verdicts at v0.305.84
> 
> - F1 seven repairs removed by suite update 0.303.8->0.305.29 (05ff187, 2026-09-06): still missing; patches/F1-lost-repairs/.
> - F2 prose test reads live queue: still fails; patches/F2-timeline-prose-window-fixture/.
> - F5 leading tabs/spaces: ALREADY FIXED (0.305.70, 0.305.81; blockIndentContent at requestmodel/visible_sections.go:111); no patch.
> - F6 any 3+ punctuation run opens a fence: present, but deliberate and tested (question, no patch).
> - F7 sectionLineBounds ignores CRLF: present; patches/F7-crlf-section-bounds/.
> - F8a/b/c: present (prose). F9 present in code, NOT reproduced in a browser. F10, F11, F13a, F13b present (prose). F12 question only.
> - F14 (new, found during this run) qualify gate has no route for a verification-only request: present (prose and a question).
> - Left out on purpose: F13(c) temp-file round trip is intentional behind DO_WORK_COMPATIBILITY_SHIM; F13(d) each curl in do-work-cli.sh already has --max-time.
> 
> ## F1 lost repairs
> 
> Update 05ff187 (0.303.8 -> 0.305.29) was the fourth to silently remove local repairs; 0.305.59 -> 0.305.84 restored none. `git apply --check` against v0.305.84 files:
> 
> - bed3763 P0-P4 priority aliases: bed3763.patch (model.go, model_test.go, schema_normalization.go + test) applies; the work-reference.md hunk does not (row moved to line 248, now mentions "Earmarked"), so bed3763-work-reference-rebased.patch carries it. The frontmatter was migrated locally in 4293410, so upstream may skip this one.
> - deea40c fixed timeline prose-window fixture: applies.
> - 5d812a1 fixed fixtures for two timeline probes: applies.
> - b3d5d25 regression test in generate_test.go: applies.
> - d39e527 absent-target preservation case in archive_fetch_test.go: applies.
> - 821913f rationale for contrast probe animation wait: applies.
> - ce8e35e: original fails on one hunk (the `time` import; timeline_browser_probe_test.go was rewritten upstream). ce8e35e-rebased-onto-deea40c-5d812a1.patch is the original minus that hunk. Lost assertions: after the refused forward press on the trailing seven days, press Previous then Next and capture toolbar state; Previous must be enabled, move both endpoints back exactly 7*24*60*60*1000 ms, keep the span, leave data drawn; Next must be enabled and return both endpoints to the exact start instants; endpoints compared as instants not readout text (helper timelineProbeInstant); the refusal is asserted outright, not in a two-branch conditional. At v0.305.84 the file never clicks timeline-period-prev (only the wiring loop at line 2053).
> 
> All applied in order (deea40c, 5d812a1, b3d5d25, d39e527, 821913f, ce8e35e-rebased, bed3763, bed3763-work-reference-rebased) to a scratch copy: queue-kanban with QUEUE_KANBAN_BROWSER_PROBES=on `ok 106.322s`; gofmt/vet clean; do-work-cli archivefetch and schemanormalization `ok`.
> 
> The archive still marks REQ-177 (timeline step navigation), REQ-179 (static-output target refusal test), REQ-180 (HTTP archive staging test case), REQ-184 (contrast probe wait rationale) and REQ-185 (fixed timeline probe fixtures) as status: completed although their code is gone.
> 
> ## F2
> 
> TestBrowserBehaviorTimelineProseDescribesOnlyTheWindowOnScreen (tools/queue-kanban/timeline_browser_probe_test.go:2356) builds its page from generateLiveSiteInDir(t) (line 2357) and asks for the fixed week 2026-07-27..2026-08-02. Observed at v0.305.84: `--- FAIL ... timeline_browser_probe_test.go:2482: the past window produced "Nothing was drawn between 2026-07-27 00:00 UTC and 2026-08-03 00:00 UTC. ... 151 REQs are outside it."` This makes just maintainer-check red where the queue has nothing in that week. Fix: fixed test data (deea40c); with it the test is `ok 2.085s`.
> 
> ## F5 already fixed
> 
> blockIndentContent (requestmodel/visible_sections.go:111) allows at most three spaces and refuses a tab; CHANGELOG 0.305.70 and 0.305.81. Probe: tab-indented and four-space-indented `## A` give 0 sections, three-space gives 1. No action.
> 
> ## F6 question
> 
> visible_sections.go:43 opens a fence on `length >= 3 && isEnclosingFenceCharacter`; isEnclosingFenceCharacter (line 181) accepts every ASCII punctuation mark except - _ * = #. Probed: `$$$/## A/$$$/## B` -> 1 section (CommonMark: 2); `+++` same; unpaired `!!!/## A/## C` -> 0 sections (CommonMark: 2). It looks deliberate: header comment line 19, lessons-do-work-cli.md "closed-enumeration-for-a-condition", and TestFoldTimingSummaryReadsFencesAsAConditionNotASpelling (lifecycletiming/lifecycle_timing_test.go:562) asserts an unknown dialect fence like `:::note` hides what it wraps. Restricting to backtick/tilde makes that test fail, so no patch. Question: is hiding on any unknown punctuation run the intended trade? Would a narrower rule for the unpaired case (opener with no closer before the next `## ` heading is text) be acceptable?
> 
> ## F7
> 
> appendSectionEntry (internal/requeststate/state_apply.go:1024) uses sectionLineBounds (line 1040) which compares `line != "## "+section` without stripping `\r`, so a CRLF file gets a second `## Blocked`. It also matches a fenced heading and misses trailing spaces. Callers: lines 622, 679, 885. markdownSectionBytes (line 1011) already uses requestmodel.VisibleSections. Patch: appendSectionEntry on VisibleSections (skips HeadingIndent != 0), sectionLineBounds deleted, new append_section_entry_test.go. RED: 3 failures (CRLF and trailing-space duplicate headings; fenced `## Blocked` taken as the section). GREEN: `go test ./internal/requeststate/` ok 6.594s. Mixed line endings (LF entry in a CRLF file) left alone.
> 
> ## F8
> 
> a) finalization/finalization_apply.go:114 returns FINALIZATION-PRIMARY-COMMIT forever once journal.PrimaryCommit is set; finalizationFailure (lines 670-674) tells the user to run `git revert <sha>`; after that nothing clears PrimaryCommit, and --discard-journal (finalization_commands.go:146) refuses any phase other than prepared. Could a performed revert be recognised so the journal can be discarded or restarted from a fresh manifest?
> 
> b) consumeRecoveryRecord (finalization_commands.go:220-232) turns a request-scoped refusal, including OutcomeRisk (exit 4), into a "set aside" record and continues; aggregate starts OutcomeSuccess (line 83) and is never raised, so recover-finalization exits 0 with an unverified commit on main. Asserted by review_regressions_test.go:105 ("recovery must set aside only the risky request and continue", Outcome == OutcomeSuccess). Request: keep draining other journals but raise the final outcome to OutcomeRisk when any set-aside record carries a revert action.
> 
> c) gittransaction/exact_commit.go:76-78: when `git rev-parse HEAD` fails after a successful commit, committedRisk(result, "...ID could not be read", "HEAD") (git_transaction.go:1520-1523) stores the literal "HEAD" as CommitSHA and builds `git revert HEAD`; finalization_apply.go:139-140 saves it as journal.PrimaryCommit. Same pattern at git_transaction.go:702-704. Ctrl-C is NOT the trigger (the rev-parse is not on a cancelled context); the trigger is rev-parse itself failing. Proposal: leave CommitSHA/RevertArgv empty, say the ID is unknown, and compare HEAD with beforeCommit (read at line 62) to report that HEAD moved.
> 
> ## F9 (not reproduced in a browser)
> 
> web/board-timeline.js:2870 registers renderVisibleRows as a scroll listener on #board-main via addTimelineListener; it is released only when renderTimelineView runs again (line 1625). Nothing releases it when the user switches views (board-controls.js:81-84 only renders on entering the timeline), and #board-main is shared by all views. board-controls.js:43 resets scrollTop only when entering the timeline. From code reading only; nobody has reproduced either effect. Could you check whether renderVisibleRows does visible work while the timeline panel is hidden?
> 
> ## F10
> 
> rootedfs/rooted_fs.go has no build tags. Rename (line 32) and Link (line 49) check with root.Lstat then call os.Rename/os.Link on the joined absolute path on every Go version; the package comment admits the residual window. Go 1.25+ has root.Rename/root.Link on the directory descriptor. CHANGELOG 0.305.46 says swapped parents are "refused instead of followed, and the rest of the rooted guarantees are unchanged", which overclaims. Proposal: current functions to rooted_fs_fallback.go with `//go:build !go1.25`; rooted_fs_native.go with `//go:build go1.25` calling root.Rename/root.Link. Tried on Go 1.26.1: compiles; only TestNestedRootsRefuseUnresolvableNames (rooted_fs_test.go:84, asserts the fallback's errRootUnresolvable) fails and needs the same `!go1.25` tag, so no diff. Suggest correcting the 0.305.46 entry.
> 
> ## F11
> 
> web/board-user-request-summary.js:167 `memberRemainingMinutes = estimateMinutes - (liveMinutes || 0)`: estimateMinutes is active minutes (estimateP50ActiveMinutes, lines 136-138, or the Timeline median); liveMinutes is liveClaimElapsedMinutes (line 35), wall-clock now - claimed_at. An idle claim is charged idle hours, so "Remaining" shrinks too fast and overrunForecastCount (lines 168-170) can flag a request that is not over. Could it subtract measured active minutes (the implementation span) instead?
> 
> ## F12 question
> 
> heavyverification/fast_stage_evidence.go is 726 lines; DecideFastStage (613), RecordFastStage (665), InvalidateFastStage (697) are reached only via decide-fast-stage / record-fast-stage / invalidate-fast-stage (heavy_commands.go:19-31, 41-81). No action, script, recipe or hook in this consumer repo calls them (searched .claude/, scripts/, justfile, docs/). Is an action meant to call them, or is this ahead of its wiring?
> 
> ## F13
> 
> a) rollbackWithRoot (gittransaction/git_transaction.go:1088) and rollbackWithoutRoot (1217) walk the same states list ("the two walk the target kinds in the same order", lines 1085-1087). The second exists because rollbackFailure (1050) may fail to open the root (line 1062); os.OpenRoot(recorder.repositoryRoot) is also called at lines 140, 223, 339, 377, 391, 423, 520, 972. Better fix: open one root at the start of the transaction and hold it through rollback; a failed open then stops before any change and rollbackWithoutRoot can go.
> 
> b) knowledgecommands/interview_derivations.go:59 interviewClockPattern `[0-9]{1,2}:[0-9]{2}` ignores AM/PM (lines 82-84 only pad the hour). parseInterviewCadence at v0.305.84: "daily 5:00 PM" -> clock "05:00"; "weekly Friday 5:00 PM" -> "05:00"; "weekly Friday 12:30 AM" -> "12:30" (want 17:00, 17:00, 00:30). Proposal: capture optional am/pm, convert to 24 h (12 AM = 00, 12 PM = 12), keep interviewValidClock.
> 
> ## F14 qualify gate: no route for a verification-only request
> 
> Found and verified in this consumer repo's 2026-10-09 run, at v0.305.84. In this repo, REQ-215 only had to re-run `just check` against the live vulnerability database and record the result. Its first qualification pass stopped, and it passed only after a real docs change was added.
> 
> - handleQualify (do-work-cli/internal/corehelpers/checks.go:262) reads the backticked paths under `## Implementation Summary` (line 289). With no path it returns QUALIFY-SUMMARY-MISSING (line 295), and the evidence reads "Implementation Summary is missing or empty" even when the section exists and says why nothing changed.
> - A summary that lists only do-work/ paths fails as well, with QUALIFY-NO-PROJECT-FILES (lines 351-352).
> - The only allowed empty summary is the repository-gate repair no-op (`**Files changed:** None — verified repository-gate repair no-op.`, actions/work.md:335, checked by validate-already-green-repair). work.md:335 says that exception "never blesses an empty ordinary REQ".
> - Probe: a REQ file whose summary reads `**Files changed:** None — verification only; ...` gives `qualify: findings` / `QUALIFY-SUMMARY-MISSING [error]: Implementation Summary is missing or empty`, exit 1.
> 
> The effect is that a request whose deliverable is a verified fact (re-run a gate, confirm a drill) cannot finish honestly. It has to invent a file change, use a do-work/ evidence file that the gate also refuses, or be abandoned. Could a verification-only route be added, for example a typed marker on the REQ plus required evidence (the exact argv, its exit status and the revision), handled the way the already-green repair path is? A smaller first step could be an evidence message that tells "section missing" apart from "section present, no files claimed".
> ```
