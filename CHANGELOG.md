# Changelog

What's new, what's better, what's different. Recent release entries are here, most recent first.

For the complete release history, read this file and the archives below from newest to oldest. Archives are tracked in git and excluded from the distribution tarball; these GitHub links also work from an installed copy.

- [0.121.1 through 0.303.6](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-09-05-up-to-v0.303.6.md)
- [0.110.0 through 0.121.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-07-13-up-to-v0.121.0.md)
- [0.65.0 through 0.109.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-07-07-up-to-v0.109.0.md)
- [0.50.0 through 0.64.1](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-04-13-up-to-v0.64.1.md)
- [0.1.0 through 0.49.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-04-07-up-to-v0.49.0.md)

## 0.305.113 — do-work run --coordinate Keeps the Main Session a Coordinator, With Run Rules, a Preflight and a Stall Check (2026-10-10)

About 13 sessions started with a pasted directive telling the main session to coordinate and not integrate. One flag, or the phrase "drive the queue", now does that, and the run rules that lived only in a memory note are in the skill. The rules make stalls and load-only failures visible and shorter; they do not prevent them.

- `do-work run --coordinate [--fan-out N] [REQ-NNN ...]`: the main session dispatches builders, writes briefs and reports. It never builds, merges or runs the span from the hand-back merge through Step 9; no integrator pushes; a blocking question becomes a `pending-answers` follow-up. Bare `--coordinate` implies `--fan-out`. Plain `--fan-out` is unchanged. Routed from "drive the queue", "use the main session as a coordinator" and `run --coordinate`, and shown in the argument hint.
- Run rules whenever the coordinator shape is used: each builder and integrator appends to `REQ-NNN-progress.log`; one full suite at a time under `full-gate.lock`; a red full gate reruns the failed suites alone first; three new `manifest.md` columns (takeover-to-finalization minutes, full gates run, stall restarts); a hand-back coverage check. Neither file is queue state or a liveness claim.
- A four-line preflight before the first spawn (board disk level, leftover processes, stale locks, optional `do-work/run-policy.md` copied into every brief) and a stall check every 15 to 20 minutes that reads `do-work status --watch` and the progress logs. An integrator silent about 20 minutes, and not inside a long command such as a full gate, is replaced by a fresh one through `advance REQ-NNN`, never `recover`.
- A handoff from a coordinated run resumes with `do-work run --coordinate --fan-out N` and hands leftover claims to integrators. The coordinator writes one on its own, in a writing gap and once per run, when the harness reports high context use. Teardown stops idle background agents and lists merged worktrees without removing them.
- `do-work status` no longer offers an `rm` line for a `full-gate.lock`. Stale restatements fixed at the end of the wave: two cleanup step pointers now say Step 9, and the ai-report Input sentence, guide and architecture-report row now name the `--kind` form for open work.

## 0.305.112 — do-work trace Shows How Much of a Spec Is Already Captured and Built Before Anything Is Captured Again (2026-10-10)

Asking "is this captured, and how much of it is built?" about a ticket page, a screenshot, pasted text or a UR used to mean a hand search, and shipped asks were sometimes captured twice. `do-work trace` answers it with one table and a dated verdict per ask, and captures only the gaps, and only when you say go.

- New action `do-work trace <url|image|text|UR-NNN>` (`actions/trace.md`, details in `actions/trace-reference.md`). It splits the source into asks A1, A2, ..., searches URs, REQs, credited commits, code and tests for each ask, and prints one row per ask: the matched REQ, the commit, the criteria with evidence (`N of M + GREEN yes|no`), the verdict `complete`, `partial`, `not started` or `unverified` with the date, the gap, and any open question.
- Read-only until go: the transcription stays in a temporary directory outside the repo, so `git status --porcelain` is unchanged after a no-go trace. A URL fetch that returns only a site name is never judged; trace asks for a paste or a screenshot instead.
- On go, trace makes one `capture-request:` call with only the `partial` and `not started` rows. A partial match on an archived REQ becomes an addendum REQ with the commit in `## Prior Implementation`, and an external ticket id starts the UR title and the transcription's first lines.
- New read-only board subcommand `queue-kanban request-commits [--repo-root DIR] REQ-NNN...` lists every commit the board credits to each REQ (a bracketed `[REQ-NNN]` subject token or a touched REQ file), with hash, date, the count of paths outside `do-work/`, and subject. It uses the same attribution helper as the board's activity lines, so the two cannot disagree.
- Routing: "is this captured", "already implemented", "how much of it is implemented", "didn't we have a request", `trace`, and "coverage" together with a spec, URL, image or UR now go to trace, from a row above the capture fallback. Also added to the argument hint and the help menu.

## 0.305.111 — do-work status Shows Each In-Flight REQ's Class, ETA and Remedy in One Table (2026-10-10)

Asking "is it stuck?" used to mean rebuilding the answer by hand from REQ files, the run manifest, builder branches and the run directory. `do-work status` now answers it in one table, from the same facts the board shows. It reports ages and never says a run or builder is dead.

- New action `do-work status [REQ] [--watch]` (`actions/status.md`, `docs/status-guide.md`), routed from "status", "status and ETA", "eta", "is it stuck", "stuck", "how do I unblock" and "why is REQ-N taking so long" or "why is REQ-N in". The bare word "blocked" still routes to clarify.
- New `do-work-cli run-status --board-facts FILE [--run DIR] [--req REQ-NNN] [--watch]`: one finding per claimed, blocked, waiting or earmarked REQ, in the same shape doctor uses. Each row gets one class from C1 to C8 (first match wins: C8, C3, C4, C5, C6, C7, C2, C1), an ETA from the frozen p50 estimate ("over estimate by N min" once past it, never negative), and a `run_status` block with ages, the run-manifest row, whether the hand-back landed, the builder branch, and run-directory files labelled "run-local, not interpreted".
- A claim older than the board's 3-hour threshold offers the read-only `git log --full-history` as its next command. `run-with-recovery` is named only in the explanation, with what it destroys and when it is safe.
- `queue-kanban open-work --format json` prints the board's open tickets with column, placement reason and last activity. The text digest is unchanged.
- `/loop 15m /do-work status --watch` prints a compact tick of about 10 lines for 8 open REQs and writes nothing.

## 0.305.110 — ai-report Writes a Decision Brief for Unfinished Work With --kind proposal or root-cause (2026-10-10)

Proposal and "why did this fail" reports were written from scratch each time, because `ai-report` accepted only completed work. Two opt-in kinds now give them a fixed, decision-first shape. Without `--kind`, the report works as before.

- `do-work-toolbox ai-report --kind proposal <topic|REQ-NNN|UR-NNN>` writes a brief that opens with the decision, then two to four options (O1..On) with benefit, risk and cost, a recommendation with its reason, an evidence ledger, limits, open questions (Q1..Qn) and one `do-work capture-request:` line per option. The action prints the capture lines and never runs them.
- The options always include the smallest-change option, such as deleting the mechanism or doing nothing extra, priced like the others.
- `--kind root-cause` puts "what happened, why, what to change" in place of the decision and options. It also accepts a failed or cancelled REQ.
- Both kinds accept a free topic, an open REQ or an open UR, and never present that work as shipped. The default kind keeps its completed-work check unchanged.
- Bundles are named `yyyy-mm-dd_hhmm_<kind>-<topic>` and carry `<meta name="ai-report-kind">`, so the report catalog can show the kind. Every mockup is labelled "MOCKUP — proposal" and kept in `generated/`, apart from real screenshots.
- Routing adds `proposal report`, `root cause report` and `options report`; the toolbox help and the ai-report guide list the new form.

## 0.305.109 — validate-feedback Can Capture, Verify and Run the Accepted Findings, and Core do-work Routes to It (2026-10-10)

Accepting review findings used to take two to five hand-typed prompts after the triage, and a review written for another repo could be triaged against the wrong code. One opt-in flag now carries the accepted findings into the queue, verifies them, and can start the run.

- `do-work-toolbox validate-feedback --capture [--run]` (or "then capture the accepted ones" / "capture and run" in the same request) asks one Accept / Park / Drop question per Discuss item, captures every accepted finding as one REQ under one new UR with its finding id and source, runs verify-requests on that UR, and prints one combined report. With `--run` and no Important, Minor or Ambiguous gap it continues into `do-work run UR-NNN`; with a gap it stops and prints the two commands to run after the fix. `--run` alone prints a usage line.
- With `--capture`, a wrong-repo check stops before any triage or write when more than half of the cited paths are not in this repo.
- The action reads a file path as input and records the path as each finding's source.
- Core `do-work` now routes `validate-feedback` and `triage feedback` to the toolbox action instead of capturing the text as a new request. `capture.md` accepts the validate-feedback hand-off as an explicit promotion.
- Without the flag the action stays read-only and its report is unchanged.

## 0.305.108 — capture-files Prints a Fill-In Example Manifest With Payload Templates, and the Capture Reference Shows the Fence Capture Accepts (2026-10-10)

A session writing a capture no longer has to read Go source to learn the manifest keys, and copying the reference UR example no longer fails the first dry run. One command prints a manifest and both payloads that pass `capture-files --dry-run` once the placeholders are filled.

- New `do-work-cli capture-files --example [--raw-input <file>]`. It prints the manifest, a UR template and a REQ template, each after a `==> <path> <==` line that names the file to write it to. With `--raw-input`, the UR template already holds the exact Full Verbatim Input block built from the file's bytes, `raw_input` names the file, and a header line says to leave that block as printed. It writes nothing and refuses under `--format json`.
- `actions/capture-reference.md`: the UR example now uses a bare three-backtick fence. The paragraph after it states the real rule once (one backtick longer than the longest backtick run, never shorter than three, no info string, every line prefixed `> `) and points at the new command.
- `actions/capture.md` Step 5 names the command as the way to get a valid manifest.
- The test that checks the published capture examples now runs the UR example through the raw-input check, so a wrong fence can no longer ship green.

## 0.305.107 — ai-report Revises a Report by Writing a New Linked Revision and Never Editing the Old One (2026-10-10)

You can now update an existing report without breaking the rule that a published report never changes. `ai-report revise <dir|latest> [what changed]` writes a new report next to the old one, and the catalog links the two.

- New form `ai-report revise <dir|latest> [what changed]`. It re-checks the old report's claims against the current repository, then writes a sibling bundle named `<date>_<time>_<same-slug>-rev<N>` through the normal collision-safe path. `latest` is the newest bundle in the catalog; if that bundle was already revised, the form follows the chain to the newest revision.
- The new report opens with a "rev-N (date)" block (Changed, Still to do, a link to the old report) and names the old folder in an `ai-report-supersedes` meta tag. The catalog is regenerated, so the old report shows "superseded by" the new one. The old bundle's bytes never change, and nothing is committed.
- A revise of a report that is not about completed work, such as a deploy guide, keeps that report's own sections. The render check reports the link back to the old report as broken, because it serves only the new folder; for a revise that one finding is expected.
- Report Design Rules: a report with more than six top-level sections or about 1,500 words gets a table of contents near the top.
- Every run that writes a report now ends with the report's path and a `file://` link to open it.
- The guide, the help list and the routing phrases (`revise the report`, `update the report`) name the new form.

## 0.305.106 — ai-report Render-Checks a Report Bundle With One Command at Phone and Wide Widths in Light and Dark (2026-10-10)

The report render check is now one shipped command instead of prose each session retyped. It catches a report that scrolls sideways on a phone, or links to a missing image, before a stakeholder opens it.

- A new toolbox CLI verb, `do-work-cli ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]`, serves the bundle on a free local port and renders `index.html` in a headless Chromium-family browser at 1440x1000 and 390x844, in light and dark.
- It fails on horizontal overflow (`scrollWidth > innerWidth`) and on any same-origin `href` or `src` that answers 400 or more. Other origins and `mailto:`, `data:` and `javascript:` links are never fetched.
- It writes `judge.json` and four full-page captures to `--out` or a fresh temporary directory outside the bundle, and stops its server and browser on every exit path. Verdicts: `pass` (exit 0), `fail` and `skipped` (exit 1), `error` (exit 2). With no browser found the verdict is `skipped`, never `pass`.
- `ai-report` Step 7, `stakeholder-report` Step 5 and `architecture-report` Step 5 call the command. `ai-report judge <dir>` is routed, listed in help and described in the ai-report guide. In `architecture-report`, the link back to the previous report is checked after publication, because the judge serves only the draft folder.

## 0.305.105 — ai-report Lists Every Report Bundle in One Catalog and Finds Reports on a Topic (2026-10-10)

Sessions in three repositories asked "is there an ai-report on X" and answered it by reading folder names. One built an index page by hand. `ai-report index` and `ai-report find <topic>` now answer it with one command, whatever folder naming style a project has used.

- A new toolbox CLI verb, `do-work-cli ai-report-index`, treats every directory under `ai-reports/` as a report bundle and skips loose files beside them. It picks the entry file (`index.html`, `index.md`, `README.md`, `prompt.md`, `report.md`, then the first `.html`, then the first `.md`) and records path, entry, title, date, kind, linked UR/REQ ids, verdict, `supersedes` and `superseded_by`.
- It writes `ai-reports/catalog.json` and a static `ai-reports/index.html`, grouped by kind and newest first, with superseded reports greyed and linked to the report that replaced them. A rerun over an unchanged folder writes the same bytes.
- Kind, verdict and the replaced report come from the optional `ai-report-kind`, `ai-report-verdict` and `ai-report-supersedes` meta tags. Without a kind tag, a known kind word in the folder name is used, else `unknown`.
- `--find <topic>` matches the topic against each folder name, title, kind, verdict and linked id, prints the matching bundle paths newest first, replaced reports included, and writes nothing. An id such as `REQ-412` also matches `REQ-0412`.
- A `catalog.json` or `index.html` that the command did not write is refused (`AI-REPORT-INDEX-HAND-MADE`) and never overwritten.
- The ai-report action, its guide, the toolbox router and the help list the two new forms.

## 0.305.104 — do-work-cli Builds the Finalization Manifest From the Judged Inputs and Checks the Tree First (2026-10-10)

Sessions hand-built the finalization manifest about 171 times in 15 days and learned the required commit paths by reading refusals. The CLI now fills every mechanical field and runs the finalizer's own checks before anything is written. This closes the CLI-ergonomics batch (UR-145).

- `do-work-cli finalize --auto-manifest REQ-NNN --emit <path>` takes the judged inputs (`--transition`, `--terminal-status` or `--failure-type` with `--failure-error-file`, `--message-file`, `--provenance` with `--implementation-hash` for `supplied_commit`, optional `--release-manifest`, and `--extra-path` for each extra path). It fills the request path, both digests, `completed_at`, `release_at`, the writer label and the planner's required `commit_paths`.
- It refuses with exit 1 and writes nothing when a judged input or `--emit` is missing, the index has staged paths (it lists them), an unfinished journal exists (it names `recover-finalization`), or the release version is stale (`RELEASE-PREIMAGE-STALE`).
- The emitted file goes to `advance --finalization-manifest` unchanged. The command never finalizes by itself, so the `advance` phase gate always applies. Hand-built manifests stay valid.
- `actions/work.md` Steps 8 and 9 and the work-reference Commit procedure name the command and say which fields are judged and which are filled.

## 0.305.103 — do-work-cli Writes REQ Stamps and Sections and Refuses What the Schema Forbids (2026-10-10)

Agents stamped `review_at` and `integration_at` and appended `## Testing` with hand-written scripts, and nothing enforced the append-only or section-order rules: one upstream report counted 112 such scripts in 13 days. Two new commands now do these writes and refuse what the REQ schema forbids.

- `do-work-cli frontmatter set <file|REQ-N> <field> (<value> | --at now)` writes one scalar field on a REQ in `do-work/working/` or `do-work/queue/`, quoted the same way as every other frontmatter write. It refuses to overwrite any `*_at` stamp that already has a value, accepts a stamp only in `YYYY-MM-DDTHH:MM:SSZ` form, refuses list or nested fields, and prints the value it wrote.
- `do-work-cli req append-section REQ-N --section <name> --from <file>` inserts a lifecycle section in canonical order. Running it again with the same body changes nothing and exits 0. A different body, a second copy, a name outside the canonical list, or a `--from` file that contains its own `## ` heading is refused.
- `advance` and the new writer read one shared section order, so the writer cannot put a section where `advance` would refuse it.
- `actions/work.md` and `actions/work-reference.md` name the two commands where stamps and sections are written. Writing by hand stays valid.

## 0.305.102 — One do-work-cli Command Creates, Merges and Removes Builder Worktrees (2026-10-10)

Sessions typed `git worktree add`, `ln -s`, merge and cleanup commands by hand for every builder: one upstream report counted 85 hand-run `git worktree add` calls in 14 days. `do-work-cli worktree` now runs the steps `actions/fan-out-reference.md` already defines, and it refuses before any change when a step is unsafe. It never passes `--force` or `-D` to git.

- `worktree new REQ-NNN [--from <branch>]` names the branch and worktree `worktree-agent-REQ-NNN-<slug>` from the request filename, adds `-2`, `-3` when the name is taken without deleting anything, creates the worktree in the sibling `<repo>-worktrees/` directory, and symlinks each path listed in an optional `do-work/worktree-links` file (for example `node_modules`).
- `worktree status` lists each builder worktree with commits ahead and behind the integration branch, dirty or clean, and the age of its last commit.
- `worktree merge REQ-NNN [--name <operative_name>]` runs hand-back merge steps 1 to 4: it refuses a staged index, an empty hand-back and any builder commit under `do-work/`, then merges with `--no-ff`, commits `[REQ-NNN] merge builder branch <name>`, and reports `pre` and `merge_hash`. A conflict stops with the merge in progress and lists the conflicted paths.
- `merge` also refuses a builder commit that touches a linked path. A `node_modules/` ignore line does not match the symlink, so `git add -A` in the worktree can commit it, and merging that commit would delete the main tree's `node_modules/` directory. `new` now warns when a link is not ignored; the ignore line must be `node_modules` without the slash.
- `worktree cleanup REQ-NNN [--name <operative_name>]` refuses an unmerged branch or uncommitted work, then removes the links it created, the worktree and the branch, and prunes.
- The hand sequences in `fan-out-reference.md` stay valid for agents that cannot run the command.

## 0.305.101 — Lessons Record How to Clear a Stuck Finalization Journal After a Revert, and Comments Name Their Real Readers (2026-10-10)

An upstream report asked for code in three places. The triage answered each with a short written fact instead, so the next reader does not have to ask again. This also closes the upstream-report batch (UR-152).

- `lessons-do-work-cli.md` now says what to do after a performed and verified `git revert` of a finalization's primary commit: delete the stuck journal `.git/do-work-finalization/REQ-N.json` and its payloads folder by hand, because `--discard-journal` accepts only phase `prepared`.
- A comment above `Handlers()` in `heavy_commands.go` says the three fast-stage commands are called only by `_dev/tests/maintainer-verify.sh`.
- The comment in `git_transaction.go` about the literal commit ID "HEAD" now names the callers that read a non-empty commit ID (`cleanup_apply.go`, `doctor_repair.go`, `publication_commands.go`) instead of the wrong reader.
- Two comments in the Timeline scroll browser probe now say three probes share the helper and point at the test that measures the first child.
- The 0.305.46 wording fix was not made: release rules allow link edits to past entries, not rewording what they claim.

## 0.305.100 — A Hidden Timeline No Longer Redraws When You Scroll Another Board View (2026-10-10)

After you opened the Timeline once, every scroll in Board, Calendar, Activity or Testing rebuilt the Timeline's invisible rows, because all views share one scroll area. In Testing that was about one frame of work per scroll. The Timeline now ignores those scrolls while it is hidden.

- `renderVisibleRows` returns at once when the Timeline panel is hidden. A missing panel counts as visible, so the Node-lane test page still works.
- Arriving at the Timeline sends one scroll event, so a Timeline left scrolled down and returned to after a short view clamped the position still redraws its top rows. The guard alone would have left the old rows on screen.
- Two browser probes pin both behaviors: zero row changes over ten scrolls in Testing (2330 before the fix), and the top rows back after the clamped return.
- The Node-lane test page gained a `dispatchEvent` stub for the new arrival line.

## 0.305.99 — Transaction Rollback Reuses the Repository Folder It Already Opened, and the Second Rollback Path Is Gone (2026-10-10)

A git transaction opens the repository folder once and holds it until it ends, but rollback opened the same folder a second time and kept a near-copy of its own code for the case where that second open failed. That case cannot happen, because the first open already stops the transaction before any change. The copy and its two tests are deleted, about 200 lines in all. Rollback behavior does not change.

- `rollbackFailure` now receives the open folder from `ExecuteTransaction` as a required argument and uses the rooted rollback only. Removed: `rollbackWithoutRoot`, the `openRollbackRoot` test hook, the open-and-fallback branch, and the two tests that forced the fallback (`TestRollbackWithoutRootHandleUnstagesRestoresFromHeadAndReportsTheRest` and the "root unavailable" case of `TestCreationIntentPreservesForeignIndexBeforeTransactionStaging`).
- The empty-folder-handle crash fixed earlier (REQ-598) stays impossible: the argument is required and no code path passes nothing.
- The held folder follows a renamed repository folder while `git -C <path>` follows the path. Nothing renames the folder today.
- Two comments (`exact_commit.go` and `git_transaction.go`) now say the literal commit ID "HEAD" is deliberate, because a non-empty `PrimaryCommit` blocks rollback of a commit that landed. Code unchanged.
- The nil-root note in `_dev/tests/audit-lockins.sh` no longer names the deleted code; its check is unchanged.

## 0.305.98 — Timeline Browser Probes No Longer Depend on the Live Queue's Dates (2026-10-10)

Two Timeline browser probes read the checkout's real queue but typed a fixed week, so a checkout with no REQ dated 2026-07-27 to 2026-08-02 failed the strict browser lane for no product reason. They now build their page from a fixed fixture, and the Previous and Next buttons are pressed in a probe for the first time.

- `TestBrowserBehaviorTimelineProseDescribesOnlyTheWindowOnScreen` and `TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable` use the existing range-end fixture tree, which gained two rows: a short completed REQ-164 and a long open REQ-0003.
- The Now / Fit all probe presses Previous then Next on the trailing 7 days. Both window endpoints must move exactly 7 days and return to the start instants, compared in epoch milliseconds, with the span unchanged and data still drawn.
- The forward-arrow refusal on the trailing 7 days is asserted outright instead of in a two-branch conditional.
- The other nine live-queue probes were read and left alone; none depends on dates.

## 0.305.97 — Qualify Says Whether the Implementation Summary Is Missing or Names No Files (2026-10-10)

Qualify now tells the two Implementation Summary problems apart. Before, one message, "Implementation Summary is missing or empty", covered both. A builder whose summary said "None, verification only" was told the section was missing and went looking for the wrong fix.

- `QUALIFY-SUMMARY-MISSING` evidence reads "Implementation Summary section not found" when the section is absent.
- It reads "Implementation Summary lists no backticked file paths" when the section exists but names no files.
- The finding code, severity and fixability are unchanged, and so is the parse-error message.
- Two tests in `corehelpers/checks_test.go` pin the two messages.

## 0.305.96 — Interview Cadence Reads AM and PM, So 5:00 PM Is Stored as 17:00 (2026-10-10)

Interview cadence now reads AM and PM in a clock answer, so "daily 5:00 PM" is stored as 17:00 and "weekly Friday 12:30 AM" as 00:30. Before, the AM/PM marker was ignored: the first was stored as 05:00 and the second as 12:30, so the standing slot was 12 hours off.

- `parseInterviewCadence` accepts an optional `am` or `pm` after the `H:MM` clock, in any case and with or without a space, and converts to 24-hour time: 12 AM is 00, 12 PM stays 12, other PM hours add 12.
- An AM/PM marker on an hour above 12, such as "13:00 PM", is refused. `interviewValidClock` is still the last check.
- "5pm" without a colon still gives no time.
- Four rows were added to `TestInterviewStandingSlotsParseOnlyRecurringTiming`.

## 0.305.95 — Archive Fetch Test Proves a Failed Fetch Never Creates a Target That Did Not Exist (2026-10-10)

A failed archive update must not publish a half-fetched file at a path that held nothing before. Only the case with an existing target had a test, so removing the private staging step would have broken the absent-target guarantee without any test noticing.

- `TestTotalFailurePreservesTheTargetAndLeavesNoScratch` is now a table with a `pre-existing target` row and an `absent target` row.
- The absent row seeds no file, forces the same failing fetch, and checks that the target path still does not exist and that no scratch file is left beside it.
- Test only. `archivefetch` production code is unchanged.

## 0.305.94 — Section Appends No Longer Duplicate a Heading in CRLF Files, After Trailing Spaces or Inside a Code Fence (2026-10-10)

Blocking a request, cancelling it and recovering an interrupted claim each add a line to a section such as `## Blocked`. The old line scanner compared whole lines, so a CRLF file, a heading with trailing spaces or an example heading inside a code fence ended with a second `## Blocked` heading or an entry written in the wrong place.

- `appendSectionEntry` now finds its section with `requestmodel.VisibleSections`, the reader that `markdownSectionBytes` already uses. An indented heading is not taken as the section.
- `sectionLineBounds` is deleted, which leaves `state_apply.go` 14 lines shorter.
- New `append_section_entry_test.go` pins the CRLF, trailing-space and fenced-heading cases.
- A CRLF file still receives the new entry with LF line endings. Matching the file's line ending was left out on purpose.

## 0.305.93 — Toolbox Release Check Reports Each Delivery Stage and Will Not Call Stale or Empty Content Ready (2026-10-10)

A file that exists, decodes and carries the right version can still give the consumer nothing: a hit-mask image can be fully transparent, and a served copy can be an older build than the source. An old "deployed and verified" note proves nothing today. The new `do-work-toolbox release-check <target> [--brief <path>]` traces the content to the consumer and says whether it is ready.

- New read-only toolbox action `release-check` and a short user guide beside it in the toolbox docs. It changes no project file, creates no REQ, and deploys, rolls back or repairs nothing.
- Traces intended source, generated package, serving environment and consumer-visible behavior, and tests what the consumer does with the content. File presence, a clean decode or a matching identifier alone is never consumer evidence.
- Reports implementation, integration, deployment and live acceptance separately, using core review's stage definitions by citation. Each stage is verified, failed, unassessed or not applicable, and every piece of evidence is labeled current run or historical. Historical evidence never verifies a stage.
- The verdict is not ready when any stage failed, ready only when live acceptance is verified in this run, and unknown otherwise, with the gaps and the operator steps they call for.
- Listed in the toolbox router, both help menus and the README.

## 0.305.92 — Toolbox Journey QA Checks Whole User Journeys and Names Product, Test or Unresolved Failures (2026-10-10)

Zoom alone, pan alone and reset alone can each pass their check while zoom, then pan, then reset leaves the panel offset. A broken test or an emulated phone could also be reported as a product bug or a device check. The new `do-work-toolbox journey-qa <target> [--brief <path>]` reproduces the reported sequence in a browser, tries the combined transitions around it, and gives each result one class.

- New read-only toolbox action `journey-qa` and a short user guide beside it in the toolbox docs. It changes no project source and creates no REQ.
- Runs the reported sequence first, then each named transition alone and in the combined orders the requirements or brief name. It waits on observable state, never a fixed delay, and keeps screenshots outside timed spans.
- Each result is passed, product defect, test defect or unresolved, with revision, environment (browser, viewport, physical device or emulation), steps, expected, actual and evidence paths. Emulation is labeled emulation. No browser tool or no device means unresolved, never passed.
- Ends with the smallest justified repair, the checks still left, and a capture line the user runs.
- Listed in the toolbox router, both help menus and the README.

## 0.305.91 — Toolbox Source Audit Checks Each Cited Source Against Its Claim (2026-10-09)

A link that loads is not proof that the page supports a claim: error pages served with success, redirects to a home page and index pages all pass a naive check. The new `do-work-toolbox source-audit <report-or-url-list>` reads each cited source and says, claim by claim, whether it supports what the report says.

- New read-only toolbox action `source-audit` and a short user guide beside it in the toolbox docs. It never edits the report and writes no file.
- Prints a claim table with one judgment per claim (supported, contradicted, insufficient, unavailable), then the evidence for every source: requested and final URL, every retrieval attempt, what the page is, the quoted passage, publication date and authority.
- Error pages (also when served with status 200) and unrelated redirects count as unavailable; index pages count as insufficient. A publication date the page does not state stays "unknown".
- Replacement candidates are checked the same way and listed in their own section; they never change an original judgment.
- Listed in the toolbox router, both help menus and the README.

## 0.305.90 — Review Names the Delivery Stages It Checked and Lists the Rest as Unassessed (2026-10-09)

A review's "Acceptance: Pass" from local runs could be read as proof that the served build works. A review now says which delivery stages it checked and lists the applicable stages it did not check as unassessed.

- `actions/review-work.md` Step 7 defines four delivery stages once: implementation, integration, deployment and live acceptance. Which ones apply is a judgment, and many REQs have no deployment stage.
- The Acceptance result scores only the stages the review exercised. Its one-line summary names them, in the report and on the persisted `**Acceptance:**` line.
- Step 8 gains a first category, "Unassessed delivery stages": every applicable stage the review did not exercise, by stage name, with the check that would cover it.
- `docs/review-work-guide.md` Phase 3 says the same in user words. The persisted Review block, the scores, the verdict mapping and follow-up routing are unchanged.

## 0.305.89 — The Commit Prefix Rule Covers Every Run Mode and the Hand-Back Merge (2026-10-10)

A code review of 0.305.88 found that the board now trusts the `[REQ-NNN]` prefix completely, while the rule reached builders only through the fan-out brief and no step named the merge commit's message. The rule now has one home that every run mode reads.

- `actions/fan-out-reference.md` → *Naming* states the rule for the builder's commits and the hand-back merge: `[REQ-NNN] merge builder branch <operative_name>`. Merge step 3 spells that message. The brief row points at it instead of carrying a copy.
- The 0.305.88 note that the rule was already in the brief was wrong: the rule is new since 0.305.88. The code comment and the board lesson say so now, and the lesson's line count is corrected to the 58 net lines the commit removed.
- No code behaviour changed.

## 0.305.88 — Board Activity Credits Commits by Prefix, Path and Owned Branch Tip Only (2026-10-10)

The board no longer walks a merge's second-parent range to credit builder commits that carried no `[REQ-NNN]` prefix. That code existed for a case the builder brief forbids, had already shipped one defect (the inner merge of main fixed in 0.305.80), and in the week before this release credited nothing: every builder commit behind a prefixed merge carried its own prefix.

- A commit counts for a REQ when it touches that REQ's file or run artifact, or when its subject carries the bracketed `[REQ-NNN]` token. An in-flight builder still shows through its owned branch tip, so a claimed card's last-activity line is unchanged.
- Removed from `queue-kanban`: the `%P` parent field of the windowed log, the in-memory ancestry walk, and the direct-match snapshot. `activity_correlation.go` drops from 357 to 299 lines and its tests from 429 to 337.
- `actions/fan-out-reference.md` states the commit subject rule the brief carries, and the board prime's trap line names the new credit rule.

## 0.305.87 — Fan-Out Orchestration Moves to Its Own Reference File and the Two Run Files Shrink (2026-10-09)

The two files that drive `do-work run` had grown with dense fan-out prose, and a consumer review called them hard to follow. The fan-out detail now lives in one companion file with no size limit, and the run files keep each rule in a sentence or two with a pointer to it.

- New `actions/fan-out-reference.md` holds worktree dispatch (isolation ladder, naming, the hand-back merge sequence, merge range, cleanup), fan-out and auto-wave, delegated integration and the run-directory table, landed hand-backs and dispatch timing, mid-run message routing, and the wave-end restatement sweep.
- `actions/work.md` drops from 13,864 to 12,706 words and `actions/work-reference.md` from 22,815 to 17,046 words.
- Every citation of the moved sections in the core, board, knowledge and toolbox packages now points at the new file. No rule's condition changed.

## 0.305.86 — The Disk-Space Check Measures the Repo Root Only (2026-10-09)

The board's `low-disk-space` check now takes one reading: the disk that holds the repo root. The problem that created the check was the repo itself growing to about 20 GB, and no problem ever came from a builder worktree on another disk, so the per-worktree, per-disk code is gone.

- `queue-kanban verify` measures one path, the repo root. The thresholds stay the same: a warning below 10 GiB free and critical below 3 GiB.
- Removed: the measurement of each `worktree-agent-*` worktree, the merging of readings by disk, and the "disk-space probe for worktrees" skip line.
- The finding's text and the Testing page's disk line are unchanged.

## 0.305.85 — Routine Operator Work Stays Out of the Queue; a Tracked Operator Setup Waits Behind Its Code and Clarify Completes It (2026-10-09)

Deploys and publishes you do yourself no longer become REQs that block, re-block, and end in abandon. A special operator setup you ask to track now waits behind its code, shows up when that code is done, and moves to Done when you say you did it.

- Capture writes no REQ for a routine operator act after the code ships (deploy, publish, approve, verify on live hosts are examples). It captures only the AI-buildable parts and names the operator step in its summary. A request that holds nothing else writes no UR and no REQ.
- A special operator configuration you ask to track is captured `blocked` with `depends_on` on the AI-buildable REQs from the same request. The board shows it under Pending → Waiting until they finish, then under Needs input · Blocked.
- `do-work clarify` adds "4. Done — I did it myself" to the blocked prompt. It records your words under `## Operator receipts` and marks the REQ `completed` with `completed_at`, so the board shows it under Done and `cleanup` archives it. It does not run `unblock`.
- A run that finds only a routine operator act left completes the REQ and names the step, instead of flipping it to `blocked`.

## 0.305.84 — Board Earmarked Badge Tooltip Names Needs Input · Blocked as the Home for Operator-Gated Work (2026-10-08)

Hovering a card's `assigned` badge on the board now says where work that waits on you belongs. In a consumer run, two REQs that needed the operator's input sat under Pending → Earmarked, and nothing on the board said they were in the wrong place.

- The badge tooltip ends with a new sentence: "A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here."
- Text only. Columns, grouping, card order, and the board data are unchanged.

## 0.305.83 — Clarify Says When a REQ It Releases Is Still Earmarked and Offers to Clear It (2026-10-08)

When you unblock a REQ in `do-work clarify` and it also carries `assigned_to`, clarify now tells you the default run will still skip it. Before, clarify said the REQ was back in the queue, and in a consumer run two released REQs sat under Pending → Earmarked while every default run passed them by without a word.

- Clarify Step 5.5 names the session the REQ is still earmarked for and asks one question: clear the earmark now by removing the field by hand, or keep it and run the REQ by name with `do-work run REQ-NNN`. It never clears the field without your answer.
- The clarify report lists every released REQ that stays earmarked, with its session name and the run-by-name command.
- The Earmarking section of the user guide says clarify flags this case.

## 0.305.82 — The Run Action and the assigned_to Schema Line Say Operator-Gated Work Is Blocked, Never Earmarked (2026-10-08)

A run orchestrator could park work on you with a session earmark, because the rule against it lived only in capture. In a consumer run, an orchestrator wrote `assigned_to: 'user-interactive'` on two REQs that needed production access, so the default run skipped them and the board showed them as Earmarked instead of waiting on you.

- The `assigned_to` schema line in `actions/work-reference.md` now says the field is for another session or checkout only. Work that waits on you as operator is `status: blocked` with `blocked_by` naming you.
- The Error Handling table in `actions/work.md` has a new row for an orchestrator that wants to park a queued REQ on the operator: flip it to blocked, and never write `assigned_to`.
- The missing-precondition row pointed at a heading that does not exist. Both rows now point at `actions/work-reference.md` → Failure Classification (Step 8), Environment row.

## 0.305.81 — Recovery No Longer Deletes an Indented User Heading That Shares a Generated Section's Name (2026-10-08)

This fixes data loss in `recover --take-over` that the 0.305.70 fix (REQ-635) left open. A user's sample heading indented by one to three spaces, such as a two-space `## Plan` under a bullet, was removed with the text under it as if it were a generated section.

- Recovery: a heading indented by one to three spaces still ends the section above it, but it is never removed. Generated sections always start at column 0.
- Timing: the writer no longer overwrites an indented `## Timing` sample. It appends the real `## Timing` section instead.
- `advance`: an indented heading no longer counts as lifecycle evidence, and it no longer makes the request fail as a duplicate section.
- Three new do-work CLI tests pin these cases. The recovery and Timing tests cover both `\n` and `\r\n` line endings.

## 0.305.80 — Board Activity No Longer Credits Main's Commits to a REQ Through a Merge of Main Inside Its Builder Branch (2026-10-08)

The board counts a builder's commits as REQ activity through the REQ's hand-back merge. When the builder had merged main into its branch before handing back, main's commits were also counted as that REQ's activity, which made the drawer's "Largest gap between events" too small and could move a claimed card's last-activity time.

- Board: only a merge matched directly by its `[REQ-NNN]` subject or a REQ file path expands its merge range. A merge that only got the REQ id from another merge's range, such as a builder's merge of main, no longer expands its own range.
- A new queue-kanban test nests a merge of main inside a matched hand-back merge and checks that main's commits carry no REQ id. No git command or log format changed.

## 0.305.79 — The Board Shows Earmarked REQs in Their Own Pending Group Instead of Ready (2026-10-08)

Ready on the board means the run takes that REQ next, but the run's default scan skips a REQ with `assigned_to` set, and the board still listed it under Ready and counted it as ready. Now a pending REQ that is earmarked for a session and waits on nothing shows under Pending → Earmarked, and every ready count agrees with what the run will do.

- Board: the Pending column has a third group, Earmarked, after Ready and Waiting. An earmarked REQ that still waits on a dependency stays under Waiting. The `assigned` badge text is unchanged, and its tooltip no longer says the board never acts on the field.
- `queue-kanban summary` adds an `earmarked` line, and `ready to work` no longer counts earmarked REQs. The `open-work` digest reads `pending N (N ready, N waiting, N earmarked)`.
- The board data payload carries a `pendingEarmarked` list beside `pendingReady` and `pendingWaiting`.
- `actions/work-reference.md` (`assigned_to:`), `docs/board-guide.md`, `actions/board.md`, `docs/work-guide.md`, and the board prime describe the new group. The field still drives no scheduling, filter, or sort on the board.

## 0.305.78 — Capture Tells a Session Earmark Apart From Work That Needs You at the Keyboard (2026-10-07)

Before this, "leave this one for me" and "I have to be at the keyboard for this one" both read as an earmark. Capture wrote `assigned_to`, and nothing in the pipeline released the REQ. Now capture writes `assigned_to` only when another session or checkout will take the work. When you have to be present yourself, it captures the REQ as blocked on you, so it shows under Needs input · Blocked and `do-work clarify` releases it.

- `actions/capture.md`: the Earmark assessment states the two cases. Operator-present work is captured `blocked`, with `blocked_by` naming the person in your words and `blocked_at`.
- `actions/capture.md`: the External-condition assessment now lists the earmark among the look-alikes that are not blocked. The list is keyed on that condition instead of a fixed count of three.
- `docs/work-guide.md`: the Earmarking paragraph opens with "leave this one for that session" and says the same split in one sentence.

## 0.305.77 — The Wave-End Check Now Knows Which REQs Were Set Aside (2026-10-07)

The wave-end restatement check from 0.305.76 runs in the review of a wave's last successful integration. Two cases still skipped it with no record. A delegated integrator could not tell a set-aside member from one still building, because a set-aside REQ keeps its claim in `do-work/working/`. And when a member was set aside after the wave's last review had already run, no review got the fact.

- The integrator brief, `REQ-NNN-integrate.md`, now also names every wave member the coordinator has set aside.
- `actions/work.md` Step 7 decides "last successful integration" from the archive plus the set-aside list the orchestrator holds. Under delegated integration, that is the list in the integrator brief.
- When a member is set aside after the wave's last review already ran, the check did not run for that wave. The orchestrator now names that gap under HANDLED in the run's Decision Brief. No review is re-run.
- The Decision Brief's HANDLED block now also lists run-level calls the orchestrator states itself, so this note is not dropped when no REQ recorded a decision.

## 0.305.76 — The Last Review in a Parallel Wave Now Checks What Earlier REQs Redefined (2026-10-07)

When several REQs are built at once, a git merge only catches edits to the same lines. It cannot see a later REQ that restates a rule an earlier REQ just changed the meaning of. Each REQ's review ran before the next one merged, so no review saw both. Now the review of the last REQ to integrate in a wave checks for this.

- Every review now records one `**Restatement sweep:**` line in its `## Review` block. It lists what that REQ's own change redefined, or says nothing was redefined.
- When a REQ is the last successful integration of its wave, the orchestrator tells its reviewer so and passes the earlier members' REQ ids. The review then also sweeps every element those members recorded on their lines, read from their archived REQ files.
- The orchestrator decides "last" from the archive and from the members it set aside, with the run manifest naming the wave's members. A set-aside member does not skip the check.
- An earlier review with no such line is named as unread for the sweep. Its changes are never re-derived from its diff.
- Findings route as before: report only unless impact-critical.
- The Folder Structure in `actions/work-reference.md` now says that finished REQs of a still-open user request sit flat in `archive/` until the request closes.
- The user guide says the last review in a wave also checks for these meaning collisions.

## 0.305.75 — Messages Sent During a Run Now Reach the Right REQ Without Stopping It (2026-10-07)

You can now keep talking to a running `do-work run`. Before this, a message sent during a run became an addendum REQ for the next run, and the builder working on that REQ never saw it. The new "Mid-Run Messages (any step)" rule in `actions/work.md` routes each message by what it changes, and the run never stops for it.

- A question about progress is answered from the REQ files, the run manifest and the hand-backs, not from memory of the conversation.
- A change to a REQ whose builder has not handed back yet is written into that REQ as one `## Addendum (mid-run)` section, in the user's own words. A later message for the same REQ adds an entry inside that section. Review now reads every `## Addendum` section as a requirement, so it judges the build against these words. Route C plan validation reads it too.
- Under delegated integration the coordinator keeps the words until its next writing gap and puts them in the integrator's brief. The integrator writes the section before the merge, so there is still only one writer in the project at a time.
- A change to work that already came back, to queued work, or a request for new work goes through capture as before.
- A note about what to show first goes into the run manifest, and the end-of-run Decision Brief reads it first and may reorder its sections by it.
- The user guide says you can keep talking while a run is in progress.

## 0.305.74 — Fan-Out Runs Can Hand Each REQ's Integration to One Agent at a Time (2026-10-07)

In `do-work run --fan-out`, the session that dispatches builders can now stay a coordinator. It hands each REQ's integration, from the hand-back merge through the release, to one agent at a time, so the user's conversation stays free during long integrations.

- A hand-back that has already landed is consumed, never built again. When the hand-back file named in the run manifest exists, Step 6 skips the builder and goes straight to the merge, whatever the manifest row's status says. This also covers a fresh session after a crash.
- The run manifest row now carries the builder's dispatch time, so a session that did not dispatch the builder can record its timing. The timing event is recorded only when the hand-back has just landed, because the recorder measures up to the current moment.
- The new "Delegated integration — the coordinator shape" paragraph in the work reference sets the rules. Integrators run one at a time in the same checkout. An integrator enters with `advance REQ-NNN` and never runs `recover`, run-with-recovery, `--assume-sole-authority` or `--take-over`. While it runs, the coordinator writes nothing in the project. The coordinator writes each integrator's brief, `REQ-NNN-integrate.md`, in the gap between integrators.
- An integrator stops after finalization and never runs the session checkpoint or the loop. The coordinator runs them after the last integrator returns.
- The time saved is still in the build phase only.

## 0.305.73 — Release Guard Reads the Module Declaration Once (2026-10-07)

Internal cleanup in the finalization release guard. Behaviour is unchanged.

- The shipped-change release guard now reads `suite/modules.tsv` once instead of twice. Before, it read the versioned module roots, then read all declared module sources again and merged the two lists, although the first list is only a filtered part of the second. Now it reads the declared sources once, treats the repository as a consumer project when no declared source has its own tracked `VERSION` (the same rule as before), and adds `suite` and `tools`. The list of shipped roots and the refusal message are the same for every input.

## 0.305.72 — Board Panel B Agrees With the Calibration Log at the 2h Gap Boundary (2026-10-07)

The Durations page's Panel B and the calibration log now make the same call on a REQ whose largest stamp gap is just over 2h, and the drawer row that shows the largest gap has a clearer name.

- Panel B now compares the largest gap between lifecycle stamps in whole minutes rounded down, the same way the calibration log records `max_stamp_gap_minutes`. Before, a 2h00m40s gap left a REQ out of the day medians while a re-fit reading the logged 120 kept it. A 2h01m gap is still excluded, and the rule itself is unchanged.
- The drawer row "Largest idle gap" is now "Largest gap between events". It measures only between recorded events, never the open stretch from the last event to now. We renamed the row instead of adding that open stretch, so the row means the same thing on claimed and done cards, and because the card's `last activity` line already shows the open stretch. The board guide, the calibration doc, and code comments use the new name.
- The Panel B code comment now says that a REQ with no phase stamps has one gap, its whole claim-to-completion span, so it is excluded when that span is over 2h even if the work was continuous.

## 0.305.71 — Board Activity Ignores Builder Branches With No Commits of Their Own (2026-10-07)

A claimed card's "last activity" no longer jumps to the moment a builder branch was created. A freshly cut `worktree-agent-REQ-NNN-*` branch points at the integration commit until the builder commits, and the board used to count that commit as the REQ's newest activity, which also split the card's real idle gap.

- A builder branch's tip counts as activity for its REQ only when the branch owns that commit, that is, when the tip is not reachable from the integration branch. A merged builder's commits still count through its merge.
- Each board response, and each static `generate`, now reads the worktree-agent worktrees and branches once. The VERIFY band and the activity line share that read, and it takes a fixed number of git commands however many builder branches exist, where it used to list both twice and run one more command per branch.
- The board guide's "last activity" row states the new rule.

## 0.305.70 — Recover Take-Over Keeps User Text Under an Indented Heading (2026-10-07)

Data-loss fix in `recover --take-over`: a request that showed an example `## Plan` heading indented by four spaces or a tab lost that example and the user text after it, because the section reader treated the indented line as a real heading and recovery deleted the "generated" section it seemed to start. The section reader now follows the CommonMark indent rule.

- A `## ` heading, and a code fence opener or closer, counts only after zero to three spaces. A tab in the indentation or a fourth space makes the line indented code, so recovery leaves it and the text below it alone, `advance` stops reporting a false duplicate section, and the Timing writer no longer overwrites an indented `## Timing` sample.
- A fence line whose info string holds `<!--` opens the fence and nothing else. Before, it also opened a comment that never closed, every later heading disappeared, and the Timing writer appended a second `## Timing` section.
- A `<!--` inside an inline code span on the same line is quoted text, not a comment. A request that quoted `<!--` in backticks used to hide all of its own later sections from `advance`.

## 0.305.69 — Commit Association Survives a REQ File It Cannot Parse (2026-10-06)

One archived request whose Implementation Summary bullet carried an odd number of backticks made `do-work commit` exit 2 with `PARSE-FAILED` at its association step, in every later commit, because archive records are immutable. One record's formatting no longer stops the commit.

- `protected-inventory associate` treats a request file whose Implementation Summary cannot be parsed as claiming no paths, keeps walking, and exits 0 with its owner rows unchanged. The skipped file is named by an `ASSOCIATION-SUMMARY-UNPARSED` warning, printed on stderr in text mode and carried as a finding in JSON.
- `PARSE-FAILED` is gone. The commit and inspect actions and the shell primitives guide say what a skipped record means instead of listing it as an exit-2 case.
- How a path may be written is unchanged: qualification and scope-drift still refuse an unmatched backtick in the request being finalized, where the author can fix it.

## 0.305.68 — Panel B and the Calibration Log Exclude by Largest Idle Gap (2026-10-05)

Panel B's day medians dropped every request whose claim-to-completion span passed four hours, as an assumed pause. A long run worked without a break was dropped too, while a short request that sat idle for hours was kept. The rule now looks at the largest gap between a request's lifecycle stamps.

- A request is excluded from the day medians when the largest gap between consecutive lifecycle stamps is over 2 hours (the whole span when it has no phase stamps). The verdict is `idle-gap`, and Panel B, the timeline forecast and the request-group summary say so. `reversed` keeps its meaning.
- `do-work/calibration-log.tsv` gains a `max_stamp_gap_minutes` column, so the next estimator re-fit can apply the same rule. The header decides the shape: a new log gets the six-column header, and an existing five-column log keeps writing five-column rows.
- `estimate-reference.md` → Calibration states the new rule, and the board reads it as its second reader.

## 0.305.67 — Board Cards Show Last Activity Instead of an Assumed-Pause Badge (2026-10-05)

Done cards on the board carried an `over 4h · assumed pause` badge whenever the claim-to-completion span passed four hours, even when the request was worked without a break. In-progress cards had no sign of a stalled request at all. The board now reads the evidence instead of guessing from one number.

- A claimed card shows `last activity`: the newest of the request's lifecycle stamps and the git commits linked to it, with a ticking stopwatch and the phase it fell in. A commit is linked when it touches the request's file or run artifacts, carries the `[REQ-NNN]` prefix, sits inside a merge of the request's builder branch, or is the tip of its live `worktree-agent-REQ-NNN-*` branch.
- The detail drawer of a claimed or recently finished request shows `Largest idle gap`, the longest stretch between two activity events and the phases on each side. It is an observation with no threshold.
- The assumed-pause badge is gone from done cards. Panel B's day-median exclusion is unchanged in this release.

## 0.305.66 — Recover's Takeover Message Says Where a Reset Claim Goes (2026-10-02)

Since 0.305.64, `recover` warns that `recover --take-over` resets a claimed request. The warning said the request goes back "as pending", but a request with an unanswered question goes back as `pending-answers`, and a blocked one stays `blocked`.

- The warning now says the takeover returns the request to the queue and strips its routing and orchestrator sections, without naming a single status. Nothing else changes.

## 0.305.65 — Board Intros Describe the Board by Its Page Switcher (2026-10-02)

The opening lines of the board guide and of the board tool's prime still described the board as a Kanban board plus a calendar and a testing track, three of its six pages.

- Both now describe it as a Kanban board plus the other pages in its page switcher, so they stay true when a page is added.

## 0.305.64 — Resuming a Handoff Keeps Its Claimed Requests (2026-10-02)

A session that resumed a handoff ran `recover`, followed the next step it suggested, `recover --take-over`, and lost its claimed requests: the takeover put each one back in the queue as pending and removed its written sections, even for work that was already merged. The suggested step is now safe.

- When `recover` finds a claimed request, its suggested next step is `advance REQ-NNN`, which changes nothing and names the request's next phase. The finding still says that `recover --take-over REQ-NNN` resets the claim, and that it is only for a claim no live session owns.
- A handoff now writes one `advance REQ-NNN` line per claimed request into its paste block and warns that a takeover resets a claim; the normal resume command never picks up claimed requests on its own.
- The work guide and the crash-recovery reference say the same.

## 0.305.63 — The Board Guide Covers Every Page, Its Links, and the Disk Line (2026-10-02)

The board guide still described the page switcher as three of its six pages, and never said that pages can be opened by link or that the Testing page shows free disk space. It now matches what shipped in 0.305.61 and 0.305.62.

- The guide names the switcher by what it covers: Board, Activity, Calendar, Timeline, Durations and Testing.
- A short passage lists each page and lens link (`#board`, `#board/by-ur`, `#board/urs-only`, `#activity`, `#calendar`, `#timeline`, `#durations`, `#testing`): clicking a page adds no Back steps, filters are not part of the link, and links work on a static snapshot opened from disk.
- The guide's Testing section and the board action's Testing step mention the free-disk-space line and its amber and red thresholds.

## 0.305.62 — The Testing Page Shows Free Disk Space (2026-10-02)

The low-disk-space check from 0.305.60 spoke up only below 10 GiB, so on a healthy machine the free space was shown nowhere. The Testing page now always shows it.

- One line in the Testing toolbar reads, for example, `disk: 61.3 GiB free of 177.5 GiB`, in the normal colour at 10 GiB or more, amber below 10 GiB and red below 3 GiB. A static snapshot adds "(at generation)", and a platform that cannot measure says "not measured".
- The figure comes from the same single measurement the VERIFY-band check already takes; the live board refreshes it on every load, and the warning findings themselves are unchanged.
- The board data carries the reading with the repo root shortened to `.`, so a shared snapshot still reveals no local paths.

## 0.305.61 — Every Board Page Has Its Own Link (2026-10-02)

The board has six pages and three Board lenses, but every link opened the Board page and the address bar never changed, so there was no way to point someone at the Timeline or the Testing page. Each page and lens now has its own address.

- Links such as `#timeline`, `#testing`, `#board/by-ur` and `#board/urs-only` open that page or lens directly, on the live board and on a static snapshot opened from disk. An unknown or empty fragment opens the Board page as before.
- Clicking a page or lens updates the address bar without adding history entries, so Back still leaves the board. Filters stay out of the link.
- The browser test harness now compares the probe page's address without its fragment, through one shared helper, so page clicks no longer break its page check.

## 0.305.60 — Low Disk Space Shows in the Board's VERIFY Band (2026-10-02)

A long fan-out run with browser QA can fill a disk in a few hours while nothing in the suite notices, and builders, gates and git then fail in confusing ways. The board now measures free space where the run operator already looks.

- `queue-kanban verify` gains a read-only `low-disk-space` probe: it measures the filesystem holding the repo root and each builder worktree, one finding per device, warning below 10 GiB free and critical below 3 GiB.
- The finding reaches the CLI report and the board's VERIFY band through the existing path, with the repo root reduced to `.` and a worktree named by its branch so static snapshots stay shareable.
- The probe never deletes anything, uses only the standard library (`syscall.Statfs` on unix, `GetDiskFreeSpaceExW` on Windows), and reports itself as skipped on platforms it cannot measure.
- The finalizer's release guard now counts every module declared in `suite/modules.tsv` as shipped, so a change to the board, knowledge or toolbox package alone is a release; only version ownership keeps the VERSION-carrying rule.

## 0.305.59 — Preserve Indented and Commented Section Boundaries (2026-09-24)

Timing replacement and claim recovery now preserve requirements under indented headings and headings containing inline HTML comments.

- Restore section boundaries without changing the original heading or body bytes.
- Cover both writers with regression tests for indentation, inline comments, and LF/CRLF line endings.

## 0.305.58 — Limit Collaborative Tone to Outgoing Drafts (2026-09-23)

Polite, collaborative wording now applies only to messages and stakeholder questions drafted for other people. Replies to the session user and all handoff documents retain their previous style.

- Allow optional alternative drafts with consistent facts and requested actions, without duplicating stakeholder records or question IDs.
- Preserve technical evidence, uncertainty, and warnings when wording outgoing requests; leave stored records and existing report bundles unchanged.

## 0.305.57 — Softer, Collaborative Follow-Up Questions (2026-09-23)

Follow-up questions now ask for the next step politely while preserving technical details, constraints, and consequences. Core, knowledge, and toolbox question guidance share the same tone and example.

## 0.305.56 — Respect Server Retry Delays for Authenticated Downloads (2026-09-16)

Authenticated downloads now honor Retry-After seconds and HTTP dates within the existing retry budget. Cancellation, attempt limits, credential scoping, and private atomic publication remain intact; an excessive delay stops retries instead of retrying early.

## 0.305.55 — Restore Frontmatter Accessor Compatibility (2026-09-16)

Frontmatter reads now use the shared document parser and schema normalization. Lists, quoting, comments, aliases, and case-sensitive values retain their established meaning; membership checks return empty stdout and the correct predicate exit status.

## 0.305.54 — Allow Straightforward Route A Work to Finish Without Lessons (2026-09-16)

Lifecycle advancement now honors the documented Route A exception for optional Lessons Learned. Orientation remains required, and Routes B and C still require lessons before finalization.

## 0.305.53 — Find Archive Collisions in Nested Directories (2026-09-16)

The archive-collision compatibility command now finds requests inside grouped archives. It preserves sorted paths and collision exit status, reports traversal failures, and does not follow directory symlinks.

## 0.305.52 — Preserve Markdown Examples During Recovery (2026-09-16)

Recovery and lifecycle advancement now distinguish visible sections from fenced examples and hidden comments. Shared byte-preserving section discovery also keeps unclosed examples intact and preserves the timing writer’s fence behavior.

## 0.305.51 — Keep GitHub Tokens on Trusted HTTPS Destinations (2026-09-15)

Downloads no longer send environment GitHub tokens to arbitrary URLs or upstream mirrors.

- Authenticate only exact trusted GitHub HTTPS endpoints and strip authorization on redirects outside that boundary.
- Use the shared Go downloader for authenticated atomic downloads, keeping tokens out of curl arguments and preserving token precedence and private publication.
- Cover localhost leaks, HTTPS mirrors, lookalike hosts, ports, redirects, and interrupted-target repair with regression tests.

## 0.305.50 — Default AGY to Latest Gemini Flash with Explicit Escalation (2026-09-15)

AGY now defaults to the latest available Gemini Flash at Low or Medium effort, with explicit model verification on launches and continuations. The guide adds High, Boost, focused subagents and Teamwork escalation, worker-routing boundaries, and checks for partial or denied headless work.

## 0.305.49 — Preserve Interview Exports and Correct Publication Handoffs (2026-09-15)

Interview exports retain approved entries and scalar source names, and derive scheduling and stakeholder guidance from the supplied data.

- Historical JSON preserves stale entries and empty arrays. Active scheduling, trust, and scan derivations exclude stale inputs as declared by the template.
- Stakeholder tones, avoidance windows, and recurring dependency slots now populate when the supplied entries qualify; missing or ambiguous window times are omitted.
- Architecture report commits check staged changes before creating report directories.
- Queue selection retains the work-action handoff without advertising an uninstalled Just recipe.
- The CLI prime and lesson index record the failure cases and regression coverage.

## 0.305.48 — Document Rollback Identity and Toolchain Lessons (2026-09-15)

The CLI prime and lesson records now capture the dirty-move recording order, directory-bound recursive deletion, and portable nested-root refusal tests from 0.305.47. The earlier lesson no longer assumes every Go version gives nested roots an unresolvable name.

## 0.305.47 — Restore Dirty Moves and Confine Recursive Rollback (2026-09-15)

Failed publication restores dirty source contents even when another writer replaces a created file. Recursive rollback stays bound to the directory it opened, preserving unrelated files behind a replacement symlink.

- Move recording captures destination ownership and source removal before revalidation can fail.
- Recursive removal enumerates and deletes through the same directory-bound root.
- The nested-root refusal test uses an actually unresolvable directory name, including on Go 1.26.1.

## 0.305.46 — Go 1.24 Floor and Prebuilt Binary Fallback (2026-09-14)

Hosts pinned to Go 1.24, and hosts with no usable Go at all, can now install, update and run the suite.

- The `do-work-cli` module floor drops from Go 1.25.0 to 1.24.0. The five `os.Root` methods that only exist in 1.25 (`Rename`, `Link`, `MkdirAll`, `ReadFile`, `RemoveAll`) are backported in `internal/rootedfs`; a parent directory swapped under a publish is now refused instead of followed, and the rest of the rooted guarantees are unchanged.
- When Go is missing or too old, `tools/do-work-cli.sh` fetches the prebuilt binary for the installed version from the GitHub release, verifies it against `SHA256SUMS`, caches it under `$XDG_CACHE_HOME/do-work-cli`, and runs it. `DO_WORK_CLI_RELEASE_BASE` points the fetch at a mirror; the `release-binaries` workflow publishes Linux and macOS binaries for amd64 and arm64 on every version bump.
- The launcher tolerates Go 1.24's first uncached `go tool -n`, which prints a temporary build path instead of the cached one.

## 0.305.45 — Functional Manifest Release Changes (2026-09-13)

Dependency and configuration edits in package manifests can now authorize a release. Version-only manifest edits remain excluded, including TOML tables with trailing comments.

## 0.305.44 — Actionable Committed-Risk Diagnostics (2026-09-13)

Failed post-commit verification now preserves actionable committed-risk diagnostics. Recovery retains the exact revert command while setting that request aside and continuing unrelated work.

## 0.305.43 — Exact Git Paths in Release Checks (2026-09-13)

Release checks now accept valid shipped changes when filenames contain Unicode, quotes, whitespace or line breaks. Git paths retain their exact bytes across ordinary, root and merge commits.

## 0.305.42 — Publication Ownership and Finalization Corrections (2026-09-13)

Rollback preserves foreign replacement files and staged entries, and valid finalization commits no longer fail because of Git diff formatting or literal filenames.

- Bind publication rollback to the file identity returned by exclusive creation and limit index cleanup to attempted transaction staging.
- Use consistent diff prefixes and literal paths when comparing prepared and committed changes.
- Remove retired checkpoint summaries before or after live claims while preserving claims, notes, and legacy layouts.

## 0.305.41 — Restore Strict Consumer Test Entry Points (2026-09-13)

Existing consumer gates can again select the strict JavaScript and browser wrapper tests. Each wrapper enables its probes in a separate process and rejects a run that executes none, including when Node or a browser is unavailable.

- Preserve ordinary test runs without automatically launching heavy probes.
- Cover both legacy entry points with subprocess execution and missing-runtime regressions.

## 0.305.40 — Checkpoint Summary Cleanup (2026-09-13)

Checkpoint refresh removes obsolete generated session summaries while preserving live claims, legacy recovery evidence, and authored notes.

## 0.305.39 — Commit Evidence and Publication Recovery (2026-09-13)

Finalization preserves failed commit evidence and checks committed content against the complete prepared change, including new files. Publication failures keep every intended destination visible to rollback while preserving foreign files.

- Preserve stale reservations when committed Git authority is unavailable, including unborn repositories.
- Require actual shipped changes for releases and show image-batch failure diagnostics without changing compatibility stdout.
- Keep portable atomic publication tests active on Windows while skipping Unix special-mode assertions.

## 0.305.38 — AGY Handoff and Completion Checks (2026-09-09)

AGY guidance now covers the handoff and verification failures observed during REQ-624.

- Bound the remaining assignment and check fresh progress before interrupting or resuming.
- Review command findings independently of AGY success, with linked execution and verification evidence.

## 0.305.37 — Shorter Workflow Instructions (2026-09-09)

Work and review instructions now point to existing procedure owners instead of repeating them.

- Consolidate duplicate merge, resume, review and scope guidance while retaining command contracts, decision points and recovery evidence.
- Preserve core routing and verify simple, complex and interrupted workflow paths.

## 0.305.36 — Verified Prose Reconciliation (2026-09-09)

Workflow instructions and board guidance now match their current behavior.

- Correct verified drift in cancellation, clarification, containment, board search and copy guidance, and process-ownership documentation.
- Remove stale counts and redundant claims while preserving runtime behavior and existing contracts.

## 0.305.35 — Bundled AGY Usage Prime (2026-09-07)

Agents can reference portable Antigravity CLI guidance directly from the installed toolbox.

- Ship the AGY usage prime with model selection, headless execution, background monitoring, continuation and Boost guidance.
- Link it from core and toolbox skill entrypoints, keeping project-specific pipeline details out of the shared reference.

## 0.305.34 — Restore Unmerged Badge and Prime Discovery Fixes (2026-09-07)

Long blocked conditions now fit their cards, and lesson capture finds relevant primes even when a request's prime list is empty or stale.

- Adapt the historical badge fix to the split board client, retaining full conditions in tooltips.
- Share prime discovery across work and standalone review while preserving deferred archival writes and lesson satellites.

## 0.305.33 — Compact Prime Routing and Recurring-Lesson Traps (2026-09-06)

Prime readers load less repeated detail while keeping access to the original guidance and incident history.

- Move oversized action, shell and CLI reference material into paired lesson satellites and refresh their routing-index estimates.
- Add missing Stakes, correct the updater’s equal-version outcome, and promote eight recurring lesson families into prime traps.

## 0.305.32 — Release Planning Ignores Unrelated Missing Manifests (2026-09-06)

A root-only release can now proceed when an independent component's tracked project manifest has an unstaged deletion.

- Require manifest reads only after establishing root, suite, or workspace ownership. Missing owned declarations still refuse release planning.
- Regression coverage checks unrelated Node, Rust, and Python manifest deletions, plus missing root and nested-workspace declarations.

## 0.305.31 — Finalization Refuses a Release Whose Implementation Ships Nothing (2026-09-06)

A release is a change to shipped files. Four entries below (0.305.17, 0.305.21, 0.305.22 and 0.305.23) bumped the version for changes under `_dev/tests` alone; the finalizer checked the release payload and never the implementation. Those four entries stand as history, and this entry closes the gap.

- `complete` with a release manifest now lists what the implementation changed (the supplied commit's first-parent diff, or the finalization commit's own paths) and refuses with `RELEASE-WITHOUT-SHIPPED-CHANGE` when nothing lies under a declared module root (`suite/modules.tsv`), `suite/` or `tools/` other than release metadata. A repository that declares no modules is a consumer project and is not guarded.
- Pinned by tests for a maintainer-only commit, metadata-only paths, a merge that brings in only maintainer files, primary-commit provenance, and a consumer repository.

## 0.305.30 — Atomic Download Repairs an Interrupted Target Again, and the Inventory Launcher Forwards Only What It Honors (2026-09-06)

An audit of the queue-drain releases 0.305.14 through 0.305.25 found two shipped defects and three stale sentences; this release fixes them directly.

- `atomic-download` (0.305.15) refused every existing regular file, which closed the installer's documented repair: `install.md` reads a zero-byte `SKILL.md` from an interrupted download as absent so a re-run can replace it, and this command is what that re-run calls. A non-empty regular file is still refused with its bytes intact; a zero-byte one is replaced. Dry-run and live agree, and the prescribed-shell cases now prove both, plus the failed-transfer case that a pre-seeded target had made vacuous.
- `atomic-download` read the byte count it reports from the target after the rename, so a failed stat there reported a published file as unpublished. The count is read from the private file before the rename; a transfer that leaves no private file publishes nothing and says so.
- `scripts/protected-inventory.sh` (0.305.14) crashed under macOS `/bin/bash` 3.2 with `set -u` whenever no global flag was given (an empty array is an unbound variable there); the guard landed in an untagged commit and ships with this entry. The launcher's `--format` sifting was dead code, since the hard-coded `--format text` always won; it is removed, the guide no longer claims it, and a missing `--repo-root` value now says so on stderr instead of exiting 2 silently. The prescribed-shell case gains the `--repo-root` from-outside-the-repository proof the original request named first.
- `docs/prescribed-shell-primitives.md` describes the zero-byte repair and the before-rename byte count. The REQ-603, REQ-604 and REQ-605 records carry strikethrough corrections for the claims this audit found false, and `lessons-do-work-cli.md` records the lesson.

## 0.305.29 — Keep Cached and Fresh Test Results Consistent (2026-09-06)

Fast-stage reuse now matches the Git configuration used by fresh tests, and forced runs revoke earlier passes before execution.

- Isolate global and system Git configuration for both evidence decisions and test execution.
- Revoke old evidence before forced runs and stop if revocation fails, so failures and interruptions cannot leave a stale pass reusable.
- Make finalization fixtures select their initial branch explicitly and cover the stale-result failures with regression checks.

## 0.305.28 — Correct Test Reuse and Integrity Checks; Remove Unused Paths (2026-09-06)

Covered ignored files now invalidate test reuse even when their directory names have Git pathspec meaning, and fixture-integrity self-tests retain their failures.

- Treat coverage roots as literal Git paths, with a regression for a colon-prefixed directory.
- Preserve integrity assertion failures while clearing only expected mutation-probe increments.
- Remove unreachable per-file ShellCheck code and unused finalization tracked-path caching and wrapper code while retaining batching and committed-image caching.

## 0.305.27 — Timeline Disclosure Geometry and Markdown Version Mirrors (2026-09-06)

Keep Timeline rows visible when findings disclosures change height, and prevent releases from leaving owned Markdown version mirrors stale.

- Measure the chart's current vertical offset and repaint on disclosure toggles, including nested findings.
- Recognize `version.md` in owned release roots during planning and recovery while keeping independent components excluded.
- Add browser coverage for disclosure expansion and collapse, plus release-planning and recovery regressions for root and workspace mirrors.

## 0.305.26 — Restore Single-Image Staging Claim to Adjacent to Target (2026-09-06)

Corrected the single-image generation command description in `skills/do-work-toolbox/actions/ai-report-reference.md` to reflect that `generateImage` stages its invocation-private file adjacent to the target output path (`filepath.Dir(outputPath)`), rather than in the system temporary directory.

## 0.305.25 — Batch Repeated Git Reads in Finalization and Request State (2026-09-06)

Profiled recovery and state planning to batch repeated Git reads and memoize historical commits during finalization discovery.

- Batched `existingUntrackedPaths` queries in `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` using a single `git --literal-pathspecs ls-files -z -- <paths...>` command, eliminating over 1,180 per-path subprocesses.
- Batched `existingDirtyTrackedPaths` queries in `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` using a single `git status --porcelain=v1 -z --untracked-files=no -- <paths...>` command, eliminating over 660 per-path subprocesses.
- Introduced `discoverySession` in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` to memoize immutable HEAD file images and tracked release paths across discovery phases, keyed to the operation entry head commit to guarantee freshness across write operations.
- Reduced redundant `git show` subprocess calls from 645 to 274 (-57.5%) and overall Git subprocesses in `TestRecoverFinalization` from 2,950 to 2,446, reducing recovery wall time to 18.13s.
- Added unit tests verifying batched path classification, session memoization, negative lookups for absent files, and cache invalidation across commits.

## 0.305.24 — Reduce Fast-Stage Evidence Computation Cost (2026-09-06)

Profiled and simplified fast-stage evidence computation in `heavy-runtime-fingerprint.py` and `do-work-cli` without weakening invalidation.

- Added in-memory caching of resolved binary seals in `_dev/tests/heavy-runtime-fingerprint.py`, eliminating redundant disk reads and SHA-256 hashing of identical executables.
- Consolidated working tree enumeration in `fast_stage_evidence.go` into a single `git ls-files` subprocess with combined `--cached` and `--others` flags.
- Scoped ignored file queries to stage coverage directory roots using `fastStageCoverageRoots`, eliminating full-tree ignored-file traversals while strictly preserving `laneCoversPath` checks.
- Preallocated internal slices and maps for stage seals to avoid dynamic slice reallocations.
- Added comprehensive unit tests and end-to-end behavior probes verifying that ignored files under stage coverage force execution while ignored files outside coverage preserve reuse.

## 0.305.23 — Remove Redundant Go Test Discovery using Native -skip (2026-09-06)

Replaced pre-execution Go test discovery in `run-go-tests-with-budget.sh` with Go 1.20+ native `-skip`, eliminating subprocess startup while strictly preserving test selection and safety guards.

- Replaces `go test -list` and Python regex assembly with pure-bash `-skip "^($skip_pattern)"`, escaping regex metacharacters with zero child processes.
- Preserves exact 100% test selection equivalence across fast test suites (verified on `queue-kanban` 402 fast tests).
- Retains empty-selection refusal (`no fast Go tests remain after applying the heavy prefixes`) in the post-test results processor with zero discovery overhead.
- Added comprehensive behavior probes in `_dev/tests/run-go-tests-with-budget-behavior.sh` for heavy exclusion, empty-selection refusal, and metacharacter escaping.

## 0.305.22 — Batch Shell Audits while Preserving Diagnostics (2026-09-06)

Batched compatible ShellCheck invocations and Markdown pattern audits across the shell verification suite, eliminating repeated process startup while preserving fragment syntax isolation and exact line attribution.

- `_dev/tests/action-shell-blocks.sh` batches ShellCheck across 74 extracted Markdown fences and 33 shipped shell files, writing `.meta` companion files to map GCC-formatted diagnostics back to exact Markdown paths and line numbers, eliminating 106 ShellCheck Haskell process spawns (-33.5% wall time).
- `_dev/tests/prescribed-shell-canonicalization.sh` scans 165 Markdown files using multi-pattern `grep -F -f` passes, eliminating over 2,600 child grep process launches (-84.1% wall time).
- Verifies defect detection across fragment syntax errors, quiet-grep pipeline violations, ShellCheck warnings, and canonicalization guide restatements.

## 0.305.21 — Use Fast POSIX Checksums for Fixture Integrity Verification (2026-09-06)

Optimized shared fixture integrity checking in `_dev/tests/session-start-hook-behavior.sh` by replacing Perl-based `shasum` with native POSIX `cksum`.

- Eliminates 20 Perl interpreter startup invocations across scenario runs, cutting CPU usage by 11.5% (-0.13s).
- Adds explicit mutation probes in `_dev/tests/session-start-hook-behavior.sh` verifying that byte changes, added files, and deleted files fail closed while preserving expected `actions/version.md` rewrites.

## 0.305.20 — Copy Prepared Recovery States in Finalization Tests (2026-09-06)

Extended fixture reuse in `internal/finalization` tests to copy prepared semantic legacy and planned recovery baseline states, eliminating repetitive Git seed commits and repository construction.

- `internal/finalization/finalization_commands_test.go` defines thread-safe singletons and lazy initializers for semantic legacy and planned finalization template repositories, cleaning them up in `TestMain`.
- `internal/finalization/finalization_recovery_test.go` updates `seedSemanticLegacyTail` and `seedPlannedFinalization` to instantiate isolated copies via `os.CopyFS`, eliminating 69 Git subprocess executions and reducing cold wall time by 13.4% (-6.45s) and CPU by 14.0% (-5.38s).
- Adds `TestPreparedRecoveryTemplateIsolation` proving complete filesystem and Git history isolation across copies and underlying templates, with mutation tests confirming robust defect rejection.

## 0.305.19 — Run Inventory Data Matrix In-Process Without Subprocesses (2026-09-06)

Decoupled porcelain byte parsing from Git acquisition in `internal/corehelpers/inventory.go`, allowing synthetic inventory test matrices to run completely in-process without spawning Git and CLI subprocesses.

- `internal/corehelpers/inventory.go` extracts `parseInventoryBytes` to cleanly separate byte stream parsing from `gitOutput`.
- `internal/corehelpers/inventory_test.go` executes the 45-case porcelain status matrix and 10-case secret origin/ambiguity matrix in-process, eliminating 56 `do-work-cli` subprocess executions and reducing test wall time by 82.5% (-7.31s) and child CPU by 61.2% (-1.01s).
- Adds explicit contract tests for malformed short porcelain records, missing rename origins, record ordering, metadata exclusions, and cross-row secret promotion while retaining end-to-end launcher and real-Git coverage.

## 0.305.18 — Build Integration Test CLI Once per Test Binary in Suiteinstall (2026-09-06)

Integration tests in `suiteinstall` now compile the `do-work-cli` executable once per test binary using `sync.Once` and clean up the temporary directory in `TestMain`, eliminating redundant `go build` invocations across heavy signal handling tests.

- `internal/suiteinstall/suite_commands_test.go` builds the test CLI binary on demand in an isolated temporary directory via `sync.Once` and deletes it upon process completion in `TestMain`.
- Retains signal interruption, recovery, and post-verification exit status assertions while reducing wall time by ~1.15s and total CPU by ~1.00s.

## 0.305.17 — Establish Test Efficiency Baseline with Descendant CPU and Work Counts (2026-09-06)

Introduced opt-in test efficiency instrumentation and reproducible baseline measurement tooling to record process descendant CPU times, toolchain subprocess counts, and per-file Go test duration attribution across representative verification selections.

- `_dev/tests/test-duration-log.sh` provides opt-in measurement functions capturing child CPU via POSIX getrusage and toolchain subprocess executions via PATH shims.
- `_dev/tests/test-efficiency-baseline.sh` orchestrates multi-run benchmarking across inventory, finalization, session-start, shell-audit, and heavy CLI build cases, emitting a structured evidence table.
- `_dev/tests/test-efficiency-baseline-behavior.sh` verifies opt-in isolation, exit code preservation, and Go test event stream parsing.
- Repointed archived lesson satellite links in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `_dev/primes/lessons-action-files.md`, and `_dev/primes/lessons-shell-commands.md`.

## 0.305.16 — List Merge Commit Paths with Diff-Tree -m in Finalization Range Check (2026-09-06)

`matchingHeadCommit` in the finalization subsystem now passes `-m` to `git diff-tree` when verifying candidate commits in the `PreparedHead..HEAD` range, preventing merge commits from emitting empty path listings and bypassing exactness path constraints.

- `internal/finalization/finalization_apply.go` invokes `git diff-tree` with `-m` so candidate merge commits enumerate paths modified across each parent against `EffectiveCommitPaths`.
- `internal/finalization/finalization_apply_test.go` verifies rejection of merge commits containing foreign paths and acceptance of clean merges.

## 0.305.15 — Unify Atomic-Download Occupancy and Handle Post-Publish Stat Error (2026-09-06)

The `atomic-download` core helper now enforces the same occupancy rule in live execution as it does during dry-run, refusing an existing regular file with exit 2 (`DOWNLOAD-TARGET-OCCUPIED`) rather than silently replacing it, and checks the error on post-rename `os.Stat` to prevent nil pointer panics under concurrent races.

- `internal/corehelpers/commands.go` checks destination occupancy before branching on dry-run, ensuring both dry-run and live runs refuse directories with exit 1 (`target is a directory`) and existing regular files with exit 2 (`target already exists`).
- `internal/corehelpers/commands.go` handles errors on post-rename stat, returning typed failure finding `DOWNLOAD-FAILED` instead of panicking on nil dereference.
- `prescribed-shell-primitives.md` documents the unified occupancy rule refusing existing destinations before fetching and the error-checked stat.

## 0.305.14 — Protected Inventory Launcher Passes Global Flags and Preserves Diagnostic Text (2026-09-06)

The `protected-inventory` compatibility launcher now forwards global flags (`--repo-root`, `--format`) ahead of the command token so it can run from any working directory, and the compatibility shim preserves diagnostic text and findings on failures instead of discarding them.

- `protected-inventory.sh` parses and forwards global options before the command verb, allowing `--repo-root` to be passed before or after the mode (`start` or `associate`).
- `internal/corehelpers/inventory.go` compatibility shim preserves prepared diagnostic text (`NO-DO-WORK-DIR`, `PARSE-FAILED`) and walk error findings on failure, ensuring exit 2 causes are visible.
- `commit.md` and `inspect.md` document that exit 2 with a finding or diagnostic text is an error to report rather than a skip condition, clarify exit 2 causes (such as non-git repository, status failure, or quarantine write failure), and note that `associate` exits 2 as not-started after `start --dry-run`.
- `commit.md` clarifies that re-running inventory uses `start` (which replaces quarantine rather than appending, whereas `associate` unions).
- `prescribed-shell-primitives.md` documents global flag forwarding and `--repo-root` support in the launcher.

## 0.305.13 — Correct Stale Mechanism and Script Claims Across Shipped Callers (2026-09-06)

Documentation and action files across toolbox, board, and core now accurately describe underlying Go command implementations, staging paths, error conditions, and conflict resolution semantics.

- `ai-report-reference.md` documents that single-image and batch image generation stage in the system temporary directory rather than adjacent to the target, and removes nonexistent retained script fallback references.
- `install.md` documents temporary skill downloads landing in `SKILL.md.download.<random>`, specifies that local Git ignores are managed inline following the canonical contract rather than calling an external helper, and removes nonexistent compatibility script references.
- `present-work.md` clarifies that portfolio summary publication reads the source once and writes both outputs from that buffer, and removes nonexistent compatibility script references.
- `board.md` clarifies that adding a local git exclude on non-git projects exits zero with a `GIT-EXCLUDE-NOT-A-REPOSITORY` warning finding.
- `prescribed-shell-primitives.md` accurately describes inventory conflict resolution: a parseable `completed_at` beats an unparseable or missing timestamp regardless of walk order, and ties fall back to `working/` before `archive/`.
- `architecture-report.md` removes nonexistent compatibility script references.
- `work-reference.md` clarifies that `run-timed-command` connects both child streams directly to the CLI's stderr handle.

## 0.305.12 — Rollback Decides the Handle Once and Closes Nil-Handle Panic (2026-09-06)

Transaction rollback now settles the rooted filesystem handle once upon entering `rollbackFailure`. If `os.OpenRoot` fails, the transaction immediately executes `rollbackWithoutRoot` to perform Git unstaging, tracked restoration from HEAD, and preimage restoration by pathname while safely recording unavailable root errors for rooted mutations. The missing guard in `quarantineAndRollbackPrivate` that caused a nil-pointer dereference panic when rolling back identity-recorded private untracked files is eliminated, and all eight downstream nil checks across rollback loops 1–3 were removed and pinned at zero.

- `rollbackFailure` cleanly partitions into `rollbackWithRoot` and `rollbackWithoutRoot`, ensuring that rooted operations only execute under a valid handle.
- All eight defensive nil-root checks in `git_transaction.go` have been removed, and `_dev/tests/audit-lockins.sh` Finding 3 now pins the count at zero.
- Added `TestRollbackWithoutRootHandleUnstagesRestoresFromHeadAndReportsTheRest`, providing the package's first unit test verifying rollback behavior when the root handle cannot be opened.

## 0.305.11 — The Rest of the Shell Guide Accurately Describes Shipped Primitives, and Inspect Associates Files (2026-09-06)

Sixteen stale claims across five sections of the prescribed shell primitives guide, plus the association prose in commit and inspect, now truthfully describe what the Go implementations and syscalls do. The two prescribed blocks in inspect that called the protected inventory wrapper with positional arguments now pass --quarantine-name, so inspect associates files with requests rather than failing with an unknown option error.

- inspect.md now passes `--quarantine-name do-work-inspect-secret-quarantine` to `protected-inventory.sh` rather than positional arguments, enabling the action to associate modified files with requests.
- Both `commit` and `inspect` now accurately document that the protected inventory wrapper takes no `--repo-root` argument, uses the current directory as the repository root, drops untracked hidden files under `do-work/` as metadata, and prints no placeholder rows for unassociated candidates.
- The prescribed shell guide's Verified Exact Publication section accurately distinguishes syscall semantics (`rename(2)` silently replacing a regular file while refusing a directory, `link(2)` refusing all existing targets) from shell `mv` and `ln` nesting.
- Lifecycle timing claims now truthfully state that elapsed seconds are derived by the command, that child processes inherit the CLI's stderr handle directly, and that redaction applies to the executable name rather than caller parameters.
- Commit file listing guidance clarifies quoting rules and the necessity of `--root` for the initial commit, and documents that `publish-portfolio-summary` and report image batch generation run through Go commands rather than non-existent retained shell scripts.

## 0.305.10 — Every Shipped Shell Block Is Now Checked for the Pipe Shape That Hides a Failed Command (2026-09-06)

0.305.5 removed `producer | grep -q PATTERN` from the repository's own scripts because, under the shell setting they all run with, the producer's death reads as the pattern's absence. One shipped block still carried it: the probe in the memory guidance that asks a local model server whether an embedding model is pulled. It could not misfire in practice, because that listing is far smaller than the size where the shape starts to fail, but it is the block agents copy, and copied guidance is where the shape spreads from.

- The probe now captures the listing first and searches the captured text. A missing or stopped server is the "no local model" answer the surrounding guidance already gives, and the block says so where the reader sees it.
- The maintainer check that already syntax-checks every shell block inside the shipped guidance (74 blocks across 32 files) now scans each one for the shape and fails at the block's own line in the Markdown file. Until now that guard read only shell scripts, so a block in guidance could ship carrying what a script could not.
- The scanner has one home, shared by both checks, so the two cannot drift apart. The check's own fixture still carries the nineteen forbidden spellings and seven safe ones, and a one-shape wiring test runs on every ordinary invocation, so removing the scan from the guidance check fails loudly rather than passing quietly.
- The maintainer's shell guide now describes the trap in the place an author reads before writing shell: the condition, why it is wrong in both directions, the measured window, the fix, and the reader shapes the guard cannot see.

## 0.305.9 — A Maintainer-Only Test No Longer Ships, and a Skewed Claim Keeps Its Progress Figures Ticking (2026-09-06)

Two fixes from an external review of the open pull request, applied directly.

- One test in the shipped `do-work-cli` module read this repository's own maintainer tree and skipped only when a directory six levels above the package was missing. In an installed copy that directory is `<project>/.claude/_dev/tests`; a project that happens to have it ran the test and failed on a manifest that never ships. The test now lives in the export-ignored file with the other maintainer-tree tests, so an installed module's `go test ./...` cannot meet it.
- On the board, a request claimed with a timestamp ahead of the viewer's clock is shown as clock skew rather than counted. Its user-request summary was never subscribed to the one-second refresh, so when the clock caught up the card's stopwatch moved while the header and drawer figures stayed on clock skew until something else redrew them. The summary now subscribes while a claim is skewed and recovers on the next tick.

## 0.305.8 — Archived Requests Stay Archived When the Project Sits Under a Directory Named working (2026-09-06)

The commit action asks which in-flight request owns each changed file. That answer came from the wrong place: the tool decided whether a request was in flight by looking for `/working/` anywhere in the request file's absolute path, so a project checked out beneath any directory called `working` made every archived request look active. Archived requests that had been cancelled or blocked could then claim files they never owned.

- In-flight-ness now comes from which directory is being walked, `do-work/working` or `do-work/archive`, and never from the path's spelling. The rule for an ordinary checkout is unchanged: a working request counts whatever its status says, an archived one only when it finished successfully.
- A request file stored under `do-work/archive/working/` was caught by the same mistake and is now treated as what it is, an archived request.
- A test builds its checkout under a `do-work/working/` directory, so any future fix that reads the path instead of the walked root fails it.
- When two requests claim one file, the later completion time wins; on an equal or missing time the request found first stands, working before archive. That order is now stated in the code where it is decided.

## 0.305.7 — One Redundant Rollback Check Removed, Eight Load-Bearing Ones Kept and Pinned (2026-09-06)

A maintenance audit counted nine places where the transaction rollback tests a filesystem handle for nil and asked for eight of them to go. Tracing every path that can reach them showed the opposite: eight are the only thing between a failed handle and a crash in the middle of restoring the tree.

- When the rollback cannot open its rooted handle it records that and keeps going, so a failed transaction still unstages its paths and still reports what it could not restore. That one missing handle then reaches eleven places, and each decides for itself what to do without it. Removing any one of seven of those checks turns a reported incomplete rollback into a crash mid-rollback.
- Exactly one check re-tested a handle every caller had already settled. That one is gone, and its precondition is stated where the callers can read it.
- A maintainer check pins the count at eight in both directions and counts the checks themselves, not any line that happens to mention them, so a deleted check cannot be paid for with a comment.
- Found on the way and not yet fixed: one place the rollback hands that possibly-missing handle onward with no check at all, so a transaction with a recorded private untracked target can crash during rollback when the handle cannot be opened. That is queued as its own fix with the package's first test for the no-handle path.
- This change was already present in 0.305.5 and 0.305.6, which were released on top of it while it was held for review; it changes no behaviour a user sees.

## 0.305.6 — The Manual Fallbacks in the Shell Guide Now Give the Tool's Answer (2026-09-06)

Two procedures in the shell guide tell you what to do by hand when the inventory tool will not run. Both gave different answers from the tool. That matters more than an ordinary documentation slip, because a person only reaches for them when the tool is already unavailable.

- The pattern list for secret-shaped filenames was narrower than the code's, so a by-hand inventory left `credential.json` unquarantined where the tool excludes it. A maintainer check now derives the expected patterns from the code itself and fails if the guide does not name every one — it catches a pattern added to the code, a pattern dropped from the guide, and the code's own function being renamed away.
- The by-hand classification said a renamed file's new path is always "modified". The code checks for deletion first, so a rename whose new path is then deleted is a deletion — and the guide told you to run a diff on a file that is no longer there.
- The by-hand association filtered out every in-flight request, which is the exact failure the paragraph above it exists to prevent. It also used a shell glob that skips the request files sitting directly in the archive directory, failed outright on a project that has never archived anything, and could name a different owner from the tool when two requests tie, because the file search returned a different order.
- The rule for resolving that tie said an archived request outranks an in-flight one. Nothing in the code compares those; the latest completion time wins, and on a tie the in-flight one is what stands.
- The excluded-file tag now has a reading rule of its own. An ordinary new file is marked excluded whenever a secret is present in the same inventory, and until now the closest rule told you to read it.
- Two descriptions of publication were also corrected: only one of the two commands verifies what it wrote, and the other reports a byte count without checking anything.

## 0.305.5 — Maintainer Checks Stop Reporting a Passed Check When Their Input Died (2026-09-06)

A shell check written as `producer | grep -q PATTERN` reports the producer's death as the pattern's absence. Under the shell setting these checks all run with, that makes the check report the wrong answer — and it is wrong in both directions, so a check written to fail when it finds something can quietly pass instead. It is invisible below roughly 36 KB of output from the producer and certain above about 200 KB.

- 129 such checks across 22 files are now 0. Ninety-nine were replaying text already captured; the other thirty ran a real command, and each of those now checks whether that command succeeded rather than trusting what it printed.
- **Two of the 129 were not waiting to break — they were broken.** A check that the old runtime is gone reported it gone whenever its file search failed for any reason, which is what happens when one directory in the tree cannot be read. Both now report the failure.
- The guard that forbids the shape moved out of the one probe file it lived in and now covers every shell file the repository tracks, running on an ordinary commit rather than only in the heavy tier. It costs a fifth of a second.
- The guard's own test carries nineteen ways of writing the forbidden shape and seven safe shapes it must leave alone, each named, so removing part of the guard reports which spelling it stopped catching rather than that a number moved. An earlier version of the guard could be walked past five different ways, including by a blank line.
- The guard's header states what it cannot see — a reader that is not `grep`, a pipeline assembled while the script runs, and three shapes it can flag by mistake, all of which fail loudly rather than passing quietly.

## 0.305.4 — Six Helper Names Defined Fifteen Times Now Have One Home Each (2026-09-06)

Six small helpers were copied across the tool's internal packages, and three of the copies disagreed with each other. That is how a version comparison ended up returning opposite answers in two places and a safety check ended up switched off without anyone noticing.

- One shared package now holds the four helpers more than one package needs. It imports nothing else inside the tool, so any package can use it and no copy has a reason to come back.
- The two copies of the version comparison returned **opposite** results for the same pair. One place asked "is the new version greater" and the other asked "is the old one older", and each had written its check to match its own copy. There is one comparison now, and it reports separately whether the version could be read at all, so nothing can mistake "not a version number" for "the same version".
- One copy was not a helper at all. It ran a path validator and threw the error away, returning nothing for the whole list whenever any single path was rejected — which quietly switched off the check that a finalization lists every file it plans to commit. That check is back on and has a test that fails if it is switched off again.
- Two path resolvers shared a name and disagreed about whether a missing file is an error. They keep their behaviour and now have two names, because the one that treats absence as an error uses that error as its existence check.
- A version string a project writes by hand — `1.09.0`, `1.0.` or `1.0.x` — is now refused with a clear message where the looser comparison used to guess at an order. Measured across 441 version pairs: 242 change, and every one of them changes to a refusal rather than to a different silent answer.
- A maintainer check pins the count of these definitions in both directions, so a seventh copy fails the build and so does losing one.

## 0.305.3 — The Shell Guide Says What Each Command Actually Does (2026-09-06)

The guide's table of shipped executable homes described work that moved into Go years ago, or credited a command with a step something else performs. All fourteen rows were checked against the code that implements them; three were wrong.

- `run-blocked-check` was described as selecting a GNU `timeout` binary and building a Bash process group. It does neither: there is no `timeout` lookup anywhere in the tool, and the process group comes from the Go runtime. Its entry now names what it owns — the probe's process-group timeout and kill escalation, first-hand launch and timeout evidence rather than facts guessed from an exit code, a bounded diagnostic identity, and the baseline comparison a focused test run needs.
- `install-memory-hooks` was credited with verification and rollback. Both are still steps a person follows in the memory actions; the command stops at a backup and a rename. Its entry now names the per-event gating and the settings merge that keeps your own keys in order.
- `record-timing-event` was credited with the folded per-request summary, which a different command produces — and the prose lower on the same page already said so, so the table had been contradicting its own document.
- Two descriptions of how a report's images are published were also wrong. Publication is not a single rename: the command claims the output directory outright, writes each finished image into it, and deletes the whole directory if any write fails, so a half-published directory never survives. And the staging directory is a temporary one, not a directory beside the output.
- A maintainer check now fails when an entry in that table credits a command with shell machinery, since every entry in it is a Go command. Its own comment states the case it cannot catch: an entry that claims another command's work reads exactly like a true one.

## 0.305.2 — The Shell Guide's Table Points at the Command, Not at Its Launchers (2026-09-06)

The "Shipped executable homes" table sent readers to nine shell scripts for mechanics that live in Go. Each of those scripts is a few lines that hand the work straight to a `do-work-cli` subcommand, so the table was routing people to the wrong file.

- All nine rows now name the `do-work-cli` subcommand that owns the mechanic. Each subcommand was read from its launcher's own `exec` line, not guessed from the file name.
- The sentence claiming that a six-line launcher orchestrates two other check scripts is gone. It was false twice over: that launcher hands everything to the Go command, and the two scripts it was said to drive are themselves launchers over the same command that nothing in the Go tree ever starts.
- The guide now explains the difference between the two, because it matters. The table says where a mechanic is owned; the launcher is what the guide's own instructions and the shipped actions invoke. They are not interchangeable — six launchers translate a legacy positional call into the flags the subcommand needs, and the protected-inventory launcher additionally selects the tag-and-path output that two actions parse one row per file from. Change the command's flags or its output and you have changed the launcher's contract.
- A maintainer check keeps both shapes from returning, and it asks what a row means rather than how it is spelled: it finds the table by its own header row, and it opens each script a row names to decide whether that script is a launcher. A row naming a file that does not exist fails, and a table with no rows fails. The first version of this check could be walked past seven different ways, including by the table's own formatting.

## 0.305.1 — The Fast Gate's Skip Decision Reads the Queue It Builds From (2026-09-06)

The stage-reuse seal added in 0.305.0 inherited a rule that ignored the whole `do-work/` tree. One of the two stages builds the Kanban board from that tree, so editing a request could leave the gate reporting a pass from a run that never saw the change.

- The board stage now declares `do-work` as an input it reads, and the tool stage declares the single archived file it reads there. Neither is assumed: the board's own prune rules and its file-mention list were walked to confirm nothing else is read.
- A separate exclusion list keeps four churn paths from invalidating a stage that never reads their bytes. The load-bearing one is the gate's own test-duration log: the stage appends to it while running, so a seal covering it could never match again and the stage would never be skipped.
- The two tests that pinned the old behaviour were rewritten rather than deleted, and each names the failure it now catches. The "queue changed" case became two — one proving the stage that reads the tree runs, one proving the stage that does not still skips — so the fix cannot become "make every change invalidate everything", which would have removed the saving entirely.
- After a review found the shipped input list itself was read by no test, restoring it to its pre-fix content now fails with three named messages, deleting either stage's rule fails with the message naming that stage, and broadening an exclusion until it swallows a real input fails too. All four of those used to pass silently.
- A typo in an exclusion's kind used to decode cleanly, match nothing, and turn skipping off for that stage forever. It is now rejected when the file is read.

## 0.305.0 — The Fast Gate Skips a Stage Whose Inputs Have Not Moved (2026-09-06)

Running the fast gate twice in a row used to re-run everything, even when nothing a stage reads had changed since the last green result. It now records what each stage read and skips a stage whose complete inputs are unchanged, printing one line per stage saying what it did and why.

- The seal is over the working tree, not over committed content. The fast gate exists to run on a tree with uncommitted work, so sealing commits would report a false pass the moment anything was uncommitted.
- Records are kept in their own key space, keyed by stage and by working-tree root, so two sibling worktrees cannot invalidate each other's evidence.
- Three commands expose the decision, the recording and the invalidation, and the gate wraps its two Go test stages in them. Anything the engine cannot determine forces the stage to run — an unreadable or missing decision never reads as a skip.
- Separately, the SessionStart hook probe stopped copying the tool's module once per scenario and now builds one shared tree. The cost was never the copying: each copy created a new absolute path, which the Go toolchain treats as a different build and re-links from scratch. Each scenario still keeps its own writable state, the banner input is now a required argument so forgetting it is an error rather than a silent pass, and a new check proves the shared tree is unchanged after every case.

## 0.304.7 — The Timeline Scrolls With the Board Instead of Inside Its Own Box (2026-09-06)

The Timeline's rows used to scroll inside a box about half the window tall, so the page had two scrollbars and the reader had to find the right one. The rows now scroll with the board, and the time axis stays pinned to the top edge the way the Activity view's column header does.

- The chart no longer has a height cap or its own scroll area. Keeping either would have left it a scroll container while measuring as though it were not.
- The time axis has a sticky rule of its own, painted on the base background, and carries the one-pixel separator the rows box used to own, so the line stays on screen instead of scrolling away with the content.
- The board's top padding now sits on whichever child the reader sees first, decided by that condition rather than by a list of the optional strips that can come before the view.
- Everything that reads or writes a scroll position moved into the board's coordinate space: the top-visible-row anchor and its restore, the virtualized visible range, scrolling a focused row into view, and jump-to-open-work. The two geometry numbers those need are measured once per render and refreshed beside the existing width cache, because the visible-row render is the scroll listener and extra layout reads would land on every frame of a drag.
- Only the scroll listener moved. Wheel zoom, drag-pan, pointer capture, keyboard, focus, hover and the width observer all stay on the chart.
- Entering the Timeline now resets the board's scroll position, which nothing did before.
- Two defects found by measuring rather than reasoning are fixed: removing the chart's focus ring handed the browser's own default ring the full-height box, caught by a screenshot; and clamping the board-relative scroll at zero made the anchor's write disagree with its read, so pressing a window chip while the board sat above the chart jumped it down by the whole offset.

## 0.304.6 — Three Lifecycle Behaviours Are Now Held By Tests That Name What They Catch (2026-09-06)

Three things the lifecycle gate does could be deleted from the code and every test would still pass. Each one now has a test that fails when it is deleted, and each test says in its own comment which deletion it catches.

- When a finding's suggested command belongs to a helper rather than to the current step, the gate rewrites it to point at the step you are actually on. That rewrite is now checked at both places it happens, by reading the rewritten command itself rather than a summary line beside it.
- Three negative controls prove the rewrite is selective rather than blanket: a neighbouring finding whose command must stay empty, a git command whose own wording must survive untouched, and one finding where one field is rewritten and the other deliberately is not.
- The guard that decides whether a focused test run counts as evidence is pinned by a nine-row table keyed on the facts the code actually reads — did it launch, did it time out, what was the baseline — rather than on the exit numbers that happen to produce those facts today. Four of the nine rows prove the guard still admits a valid run instead of refusing everything.
- An interrupted focused test is pinned by a case that signals the tool from inside the probe, paired with an ordinary failing run at the same exit status, so the only difference between them is whether the interruption was observed.
- No behaviour changed. Each of the three behaviours was temporarily broken to watch its new test fail, then restored and checked byte-for-byte.

## 0.304.5 — The Heavy Test Tier Stops Reporting a Failed Run as a Passed One (2026-09-06)

The heavy verification tier could tell you a lane passed when it had failed, and could tell you a run succeeded when it had verified nothing. Three separate routes to a false green are closed, each with a test that fails when the fix is taken back out.

- A lane that ran and exited with an error is now reported red, whatever it printed while running. It used to be enough for the lane to print a line starting with `SKIP:` — including a line printed by a test's own fixture — for the whole lane to be recorded as skipped, which reads as success. The lane's own announcement is kept and shown as extra evidence on the failure.
- A run that names no lane at all is refused instead of returning success. Before, a caller inside the tool that forgot to list lanes got a clean verdict for verifying nothing.
- A lane timeout of zero now means "use the default" rather than "expire immediately". The thirty-minute default only existed on the command-line path, so a caller inside the tool that left the field empty had its lanes killed mid-run — which then surfaced as a failing test that looked like bad luck.
- A test fixture's output can no longer reach the process running it. The check for that swaps the real output stream rather than the file descriptor, so it fails if the fix is reverted; the earlier evidence for the same fix could not tell the difference.
- The maintainer check that forbids a particular fragile shell pipeline was rewritten. The old pattern missed five ordinary ways of writing the same thing; the new one matches what makes the pipeline fragile instead of one spelling of it, and states in its own source the two shapes no source scan can catch.
- Two archive checks were reading a file listing while throwing away the error that said the archive could not be read. Both now check readability first, so a truncated archive fails instead of passing.

## 0.304.4 — The User Request Progress Figures Say When They Are Guessing (2026-09-06)

An independent review of 0.304.0 found the new per-user-request progress strip printing a confident number in two cases where it did not have one, and found the browser probe that was supposed to prove the layout measuring a page two of its three widths never had.

- A user request whose members have all run past their estimates now reads `~0 min (2 over estimate)` instead of a bare `~0 min`. On a real board this was 4 of the 5 user requests with a live claim, one of them 9.4 hours against a saved 20-minute estimate, and the reader could not tell "almost done" from "every member has blown its estimate".
- A claim timestamp the board rejects is now counted as unknown remaining time rather than as zero time spent. Before, the member was quietly charged its full estimate and the figure rendered with no qualifier, even though the same rejected stamp already showed a warning on the Active figure.
- The browser probe measures the real layout again. It reuses one page across three widths, and two things leaked between measurements: a detail drawer left open, which is a grid column rather than an overlay at these widths, and the probe's own result block, which lands in the same grid and can be a 32,000px-wide line. The group column measured 273 / 259 / 579 CSS px; it now measures 273 / 697 / 1209 again, and the probe fails if the measured box does not match the width its own case claims.
- Three claims that were written as if a test enforced them are now enforced by tests: the order of the two refresh passes inside a tick, and the absence of `try`/`catch` and of a completion timestamp read in the summary code path.
- A misspelled field name is corrected so a plain-text search finds every use of it.

## 0.304.3 — The Debug-Artifact Rule Is Stated Once, Where It Is Enforced (2026-09-06)

Three shipped action files repeated the same debug-artifact and P-A-U honesty rule seven times between them. The rule is enforced in code, so every prose copy was a restatement that could drift away from what actually runs.

- Seven mentions become two. `work.md`, `review-work.md` and `work-reference.md` now point at the finding codes `QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED` and `QUALIFY-UNIFY-DISARMED`, which is where the rule is decided.
- Two mentions are kept on purpose: `review-work.md`'s standalone-review hygiene bullet is a read the canonical `qualify` never makes, and the P-A-U template payload is byte-identical across four shipped files.
- The request claimed nine sites and the tree has seven. The commit its reproduction line names is not in this clone, so the captured failure could not be replayed and no lines were invented to reach nine. One anchor had moved eight lines earlier in the same run, so every edit was located by text rather than by line number.
- A maintainer lock-in pins the count of remaining mentions with a floor as well as a ceiling, so losing one of the two protected mentions is no longer silently green. It counts matches rather than lines — the earlier version failed on a pure reflow — reads its scanner's exit status without a pipeline, and reports path, line and matched text for each site.
- A companion assertion checks every `QUALIFY-*` code named in the action files against the code that defines it, so the enumeration goes stale loudly instead of quietly.

## 0.304.2 — The Commit and Inspect Actions Stop Repeating the Same Shell Rules (2026-09-06)

Two shipped actions carried the same block of file-inventory rules word for word, so a correction to one left the other behind. That block now lives once, in the prescribed-shell guide, and both actions point at it.

- The inventory tag legend, the secret-shaped matching sentence, the four file-reading bullets, the association semantics and both manual fallbacks are stated once in a new section of `prescribed-shell-primitives.md`. `commit.md` and `inspect.md` reference it instead of restating it.
- What each action *does* with an excluded row stays where it was — a writing action lets only the deletion proceed, a read-only one inspects the path and the deletion state. That is caller policy the guide's own charter excludes, so collapsing it would have been a behaviour change dressed as deduplication.
- The by-hand association fallback now names the two directories it tells you to glob, checked against the roots the code really walks, and the quarantine sentence names the `git rev-parse --git-path` call that produces the failure, confirmed by running both modes outside a repository.
- A maintainer lock-in now asserts each moved passage is absent from each action on its own, so a return to either file is caught by itself. It was proven in both directions: the legend pasted back into one action gives four failures naming that file, the whole pre-move file restored gives eight. The earlier line-count version could not see a one-sided return at all, and a single deleted line could push its similarity score past the ceiling — 660 fuzz runs now trip nothing.
- The lock-in reads its scanner's exit status instead of a piped total, so a scanner that cannot run fails instead of reading clean.

## 0.304.1 — The Archive Audit and Report Publishing Stop Launching Coreutils (2026-09-06)

Two places in the tool started a `find` or a `cp` process to do work it already does in Go. Both now do it directly, so they behave the same on a machine where those programs are missing, older, or a different implementation.

- The archive timestamp audit checks that it can read the archive with its own walk instead of launching `find`, and still stops with a clear message naming the path it could not read.
- Publishing an architecture report stages the draft in memory instead of launching `cp`; the published report carries the same bytes as before.
- Two maintainer test cases used to drive those failures by putting a fake program on the path, which stopped working the moment the code stopped launching one. Both now provoke the same failure inside the process, and each names the failure it expects, so neither can quietly pass on something else going wrong.

## 0.304.0 — User Request Groups Fold, and Report Their Own Progress (2026-09-06)

A user request with dozens of requests under it filled the board with cards and still told you nothing about how far along it was. Its groups now fold, and its header answers the question the card wall never did.

- Every **By UR** group header folds its own card grid, and starts open. More than one can be folded at a time, the control is a real button that announces whether it is open, and **Details** beside it still opens the drawer. The **URs only** reading is unchanged, including its collapsed default.
- The By UR header and the user request's drawer both show the same five figures for the whole request: how many requests it groups, active time already spent, an approximate remaining time, and the successful and resolved percentages with their counts. Both read the request's complete membership, so filters change which cards you see and never move the numbers.
- Active time ticks with the rest of the board while the page is open, so the header can never drift from the stopwatch on a claimed card below it.
- Missing evidence is stated, never counted as zero. A refused span, a member whose work ended with nothing measurable, an unfinished member nobody has estimated, and a claim stamped ahead of your clock each get their own qualifier, and a request with no members reads `unavailable` instead of dividing by zero.
- The drawer's grouped REQ id list starts open, is height-capped so it can no longer push `input.md` and the body out of the panel, and folds away entirely with one click.
- The board now reads each request's saved `estimate.p50_active_minutes` for the remaining-time figure, falling back to the Timeline's median only while the Timeline has enough history to call it confident. The Timeline's own forecasting is unchanged.

## 0.303.10 — The Top Bar's Controls Stay On Screen at Narrower Widths (2026-09-05)

Making the identity one unwrapping line (0.299.0) also gave it a minimum width, and on a window between roughly 760px and 1000px wide that pushed the filters and view buttons past the right edge. The page hides horizontal overflow, so there was no scrollbar to reach them — measured at 800px, the bar's contents wanted 905px and the last 105px were simply gone.

- The bar now stacks the identity above the controls below 1000px, where it used to stack only below 760px, so both halves stay reachable across that whole band.
- A project directory name longer than the line can carry truncates with an ellipsis instead of pushing the controls off-screen, and the controls take whatever width the identity does not want.
- Above 1000px the bar is the same single line 0.299.0 delivered; a browser probe measures both widths, so neither can be fixed at the other's expense.

## 0.303.9 — A Blocked Request No Longer Strands the Ones Behind It (2026-09-05)

The last two findings of the same review that produced 0.303.5. Both are cases where one item's problem was silently applied to everything after it.

- A frozen queue continuation used to stop at the first member it could not claim, so a request awaiting answers hid every sibling behind it — including one the just-claimed prerequisite had made ready. Because the continuation is re-emitted unchanged, replaying it stopped in the same place and claimed nothing. The pass now skips that member, leaves it in the queue, still reports why it was excluded, and carries on down the list.
- An ordinary `### ` sub-heading in a request body was read as an opening fence, which hid the request's own `## Timing` section from replacement: the fold left the stale summary and appended a second heading. Headings enclose nothing, so they no longer open one.
- Both lessons are recorded in `lessons-do-work-cli.md`.

## 0.303.8 — A Verify Finding Names Five Paths, Then Counts the Rest (2026-09-05)

One long list could take the board over. A builder worktree whose `do-work/` was untracked put roughly 700 paths into a single finding, and every surface that prints a finding printed all of them: the board's findings strip filled the page, the terminal report put them on one line, and the shareable snapshot carried the lot.

- Three `queue-kanban verify` findings now name five entries and then say how many they did not name — uncommitted queue state in a builder worktree, queue state committed on a builder branch, and the live members of an archived user request.
- A list of five entries or fewer reads exactly as before, so a finding about a single file is unchanged.

## 0.303.7 — The Verify Findings Strip Is One Line Until You Open It (2026-09-05)

This morning's rows (0.303.2) were still too tall: every remedy printed in full, all the time, inside a card. A remedy is what you read after deciding to act, so it now sits behind a click, and the strip is one line until you open it.

- Closed, the strip is a single line: a small warning mark, the counts, and every subject with a coloured dot for its weight (amber: read this, green: cleanup can fix, grey: a probe that never ran), with a Show button at the right.
- Open, each finding is one row: dot, subject, category in plain words, the detail clipped at the line's end, and a chevron that opens that finding's "What to do" under the row.
- The board remembers whether you left the strip open or closed.

