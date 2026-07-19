# Changelog

What's new, what's better, what's different. Most recent stuff on top.

> Older releases (0.65.x through the entry below the 20th here) live in [`CHANGELOG-archive.md`](./CHANGELOG-archive.md) — kept in the git repo but excluded from the distribution tarball. Tarball-installed copies (no local `.git`, no archive file on disk) can browse it at <https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-archive.md>. Pre-0.65 release notes live one hop further back: `CHANGELOG-2026-spring.md` and `CHANGELOG-pre-0.50.md` at commit [`bf15fe2`](https://github.com/knews2019/skill-do-work/tree/bf15fe2).

---

## 0.129.0 — Blocked-on-External-Condition Status (2026-07-19)

REQs can now wait on an external condition — LM Studio being up, a designer answering, credentials getting provisioned — instead of being mislabeled as "needs clarification" or dying as a failed environment error. A new `blocked` status names the condition and gets its own badge on the board.

- New `status: blocked` with a free-text `blocked_by` condition (plus optional `blocked_at` and a `blocked_check` shell probe). Distinct from `pending-answers` (a question for you) and `depends_on` (a wait on another REQ).
- `do-work run` re-probes each blocked REQ's `blocked_check` at scan time and auto-unblocks on exit 0 — the same "resolves dynamically" feel as dependency gating. The probe runs the repo-authored command safely (scratch file + `timeout`, fail-closed) and never halts the run.
- Mid-run, when a builder hits a missing external precondition before any edits land, the pipeline flips the REQ to `blocked` and moves on — instead of forcing a `failed` + follow-up cycle.
- `do-work capture` emits `blocked` when a task states it waits on something external; `do-work clarify` now also lets you confirm a blocked condition is met; `do-work abandon` / `roadmap` / `forensics` / `cleanup` all recognize it.
- The Kanban board shows blocked REQs in the *Needs input · Blocked* column with a distinct "blocked by: …" badge and drawer rows (condition, since, probe).

## 0.128.1 — Built-In Preferences Reference Doc (2026-07-18)

The operating nudges people paste at the start of every run — "keep writing lessons learned," "commit often," "I'm AFK, don't block on questions" — are almost all already the skill's defaults. A new reference doc maps each common nudge to where that behavior already lives, so you can stop re-typing them.

- New `docs/standing-preferences.md`: a table of common nudges → the built-in behavior and its home (lessons learned, discovered tasks, YAGNI, per-REQ atomic commits, background agents, non-blocking `pending-answers` questions).
- Calls out the two nudges that are deliberately *not* defaults — an unbounded queue drain (declined in ADR-006/014) and a backgrounded commit — so expectations match reality.
- New README Q&A entry pointing at the reference.

## 0.128.0 — Board Surfaces Completion Anomalies (2026-07-18)

A done REQ with no `completed_at` and no resolvable commit hash used to vanish from the live board — terminal, but with no instant to place it in Recently done. Those are bookkeeping bugs, and now the board shouts about them instead of hiding them.

- New always-visible "Completion anomalies" strip on the board (every view, immune to the recent-window toggle and filters); each card carries an `anomaly` badge, the reason, and the concrete fix, echoed in the detail drawer and the data-warnings banner.
- Three anomaly shapes detected: neither field present, a `completed_at` that doesn't parse (flagged even when the commit hash rescues the date), and a commit-hash field git can't resolve — the reason names the exact broken field.
- Anomalous tickets are never dated "now": no fabricated instant, no Recently-done membership, no mtime fallback — dated tickets keep the existing window behavior unchanged.
- Headless too: `queue-kanban summary` prints `completion anomalies : N` and lists the offending REQ ids.
- Prevention at the source: `actions/work.md`'s done/fail flips and `actions/work-reference.md`'s frontmatter template now make the `completed_at` + `commit` stamp an explicit hard rule on every terminal flip.

## 0.127.0 — Testing View Sorts Newest First and Filters by Date (2026-07-18)

With hundreds of finished REQs, the Ready-to-test column buried the work you just shipped at the bottom. Testing columns now read newest-first, and a date filter narrows them to a window.

- All four testing columns sort most-recent-first — by last testing activity, falling back to the REQ's completion instant; unknown dates sink to the bottom, ties break toward the higher REQ id.
- New date filter in the shared filter bar (visible only on the Testing view): Any date / Last 24 hours / Last 7 days / Last 30 days / Older than 30 days. It joins the existing search/domain/status filters and the Clear button, and never touches the Board or Calendar views.
- This repo's own queue data: four duplicate REQ ids (UR-003's doc-diet stream had reused REQ-015..018 from the earlier kanban stream) renumbered to REQ-021..024, with every frontmatter and audit-trail reference repointed — the board's duplicate-id warnings are gone.

## 0.126.1 — Testing View Review Fixes: Status Gate, Duplicate Keys, Symlink Guard (2026-07-17)

Four PR-review catches on the new Testing view (thanks, Codex review on #119) — all hardening the write path before it ships.

- The status API now rejects non-`clear` transitions on unfinished REQs (409): only terminal-success REQs — or REQs already carrying a testing record, so a returned-then-requeued REQ can restart testing — accept testing writes. A stale browser tab can no longer stamp `in-testing` onto a pending REQ.
- The frontmatter upsert consumes **every** occurrence of a duplicated key, not just the first — the YAML reader keeps the last occurrence, so a first-only edit could look successful yet read back unchanged.
- Testing writes refuse symlinked targets: the REQ file must be a regular file whose parent resolves inside `do-work/`, and `testers.md` gets the same guard — a hostile checkout can't redirect a write outside the tree.
- The Testing view keeps REQs with an *invalid* `testing_status` visible even after their pipeline status leaves terminal-success (the record and its invalid flag no longer vanish on requeue).

## 0.126.0 — Board Testing View: Track Who Tested Which REQ (2026-07-17)

With thousands of REQs, "done" told you nothing about whether anyone actually tested it. The kanban board now has a Testing view (next to Board / Calendar) where a tester picks their profile, selects a finished REQ, and marks it in-testing, tested, or returned with feedback — and the record lives in the markdown itself, so git is the audit trail.

- New Testing view in `do-work board` serve mode: four columns (Ready to test → In testing → Returned with feedback → Tested) over every terminal-success REQ, with per-card actions and an inline feedback form.
- The markdown files are the database: actions write `testing_status` / `tested_by` / `testing_updated_at` / `testing_feedback` placeholder frontmatter into the REQ file via new loopback-only `/api/testing/*` endpoints (surgical line-level upsert — everything else in the file stays byte-identical). No locking by design — changes land in the working tree and commit like any other edit.
- Tester profiles are add-or-select in the view's toolbar, stored as plain bullets in `do-work/testers.md` (created on first use, hand-editable).
- The main Board view shows a `testing` badge on any card carrying a record, and the detail drawer lists the testing meta, so testing state is visible without switching views. Static snapshots render the view read-only.
- Schema Read Contract gains the `testing_status` vocabulary (normalize-and-warn like every other enum — an off-vocabulary value renders as not-tested with an invalid flag and a data warning).

## 0.125.2 — Crew Member Renamed: karpathy.md → coding-guardrails.md (2026-07-16)

The always-loaded implementation crew member is now named for what it does, not for a person — Andrej Karpathy is more than four coding rules. The source attribution inside the file stays.

- `crew-members/karpathy.md` → `crew-members/coding-guardrails.md`; H1 retitled to "Coding Guardrails Crew Member".
- All live references updated (SKILL.md, CLAUDE.md, README, actions, specs, sibling crew files); review-work's audit heading is now "Coding-Guardrails Principle Check".
- Historical records (ADRs, archives) left as written; ADR-003 gained a one-line rename pointer.

## 0.125.1 — Reservation Review Fixes: UR Closure, Release Routing, Roadmap Section (2026-07-16)

Three PR-review catches on the new reservation feature (thanks, Codex review on #118).

- Step 8's UR-finalization check now holds a UR open for **any non-terminal** sibling — a reserved REQ no longer lets its UR archive out from under it.
- `do-work release REQ-042` now actually releases: the router passes `release <rest>` for the `release`/`unreserve` triggers so the reserve action enters release mode instead of trying to reserve the bare ID.
- The roadmap report gained the promised `## Reserved (Other Sessions)` section (with the stale-reservation recategorize hint), a reserved total, and a matching next-step line.

## 0.125.0 — REQ Reservations for Other Worktrees and Cloud Sessions (2026-07-16)

You can now reserve pending REQs for a different worktree or cloud session (`do-work reserve REQ-042 for cloud-alpha`) so the local work loop walks past them. Unlike a claim, a reservation stays in `do-work/queue/` — crash recovery can't steal it — and it travels to sibling checkouts via a normal git sync.

- New `reserved` status in the Schema Read Contract, with `reserved_for` (owner label) and `reserved_at` frontmatter; new `actions/reserve.md` (reserve / release / list).
- The default queue scan skips reserved REQs; targeted `do-work run REQ-NNN` claims them — that's how the owning session picks up its slice.
- Reservations older than 24 hours are flagged as stale everywhere they render (work-loop queue summary, exit summary, forensics, roadmap, board) with a recategorize suggestion — release, claim here, or leave it. Never auto-released.
- The Kanban board shows reserved REQs grayed out in the Claimed column with a "reserved for" badge and a stale marker.
- Intent and contract recorded in `actions/prime-req-reservation.md`.

## 0.124.4 — Qualify and Scope-Drift Checks Tightened Against False Passes (2026-07-15)

Second Codex review round on #117 caught three ways the new checks could be fooled; all three are closed.

- `qualify.sh` no longer counts the previous commit's diff as current work — a no-op builder can't pass on the back of the last REQ's changes.
- `(deleted)` summary entries now need deletion evidence in the working/staged diff, not just disk absence — a typo'd path no longer qualifies.
- `scope-drift.sh` reads only the "Files I will touch" list, so documenting out-of-scope files in "Files I will NOT touch" no longer reports false drift.

## 0.124.3 — Portable Check Scripts and Stale-Baseline Cleanup (2026-07-15)

Two PR-review fixes to the new tools/checks/ scripts (thanks, Codex review on #117).

- Replaced GNU-only `grep -P` extraction (and `\s` ERE classes) with portable `sed`/`grep -E [[:space:]]` — the checks now run on BSD/macOS grep, matching the skill's any-environment contract.
- `preflight.sh` deletes a stale `baseline-failures.txt` when the baseline passes, so Step 6.5 can never misclassify a new regression as pre-existing.

## 0.124.2 — Regrowth Ratchets: Router Word Budget and Sibling-Skill Gate (2026-07-15)

Two guards so the bloat this cleanup removed can't quietly come back.

- Contract tests now fail any commit that pushes SKILL.md past 2,650 words (post-diet count + ~10% headroom); the prescribed fix is a merge or lazy-load, never a bigger budget.
- CLAUDE.md: every NEW action must state why it belongs inside do-work rather than a sibling skill — reviewers reject additions without the justification.

## 0.124.1 — Extraction Plans for the Three Relocatable Subsystems (2026-07-15)

Plan-only release: grep-verified extraction plans for the prompt library, the interview framework, and bkb+dream now live in `decisions/audits/2026-07-15-relocation-extraction-plans.md` (maintainer docs, not shipped). No files moved; nothing changes for consumers in this release.

- Each plan names the target sibling repo, the full manifest with word counts, every inbound-reference seam to cut, and a migration note for git-clone and tarball installs.
- Recommended sequence: prompts → interview → bkb+dream (~47k words would leave the shipped skill if all three run).

## 0.124.0 — Mechanical Work-Loop Checks Ship as Scripts (2026-07-15)

Four parts of the work loop that were pure shell-logic-in-prose are now shipped executables under `tools/checks/`, so they run the same way every time instead of being re-derived from paragraphs. Judgment stays in the prose; mechanics move to code.

- `archive-collision.sh` (Step 2.0, full), `preflight.sh` (Step 5.75, full — also records a machine-readable test baseline for Step 6.5), `scope-drift.sh` (Step 5.5's review-time comparison), `qualify.sh` (Step 6.3's items 1/4/5 + the only-do-work-paths rule).
- work.md steps shrink to pointers + the judgment that remains; every pointer has a script-missing fallback.
- Contract tests now assert the pointers and scripts stay in sync.

## 0.123.2 — Small Actions State Each Guard Once (2026-07-15)

Four action files said the same rules two to seven times over (commit.md stated the .env exclusion in seven places). The guard content survives — stated once, in the section that owns it.

- note.md and scan-ideas.md: Common Rationalizations / Red Flags / Verification Checklist removed — every row mapped 1:1 onto the files' own Rules (mappings recorded in REQ-023, renumbered from REQ-017 in 0.127.0).
- commit.md: step-recap Checklist and "Common mistakes" blocks removed; generic git-advice rationalization rows dropped; the REQ-traceability rows and the hard-won terminal-status Red Flag stay.
- quick-wins.md: two generic rationalization rows dropped; the scan-breadth and dynamic-reference rows stay.

## 0.123.1 — Changelog Trimmed to the Newest 20 Entries (2026-07-15)

The live changelog was 162 entries (~24k words of shipped payload) while the version action only ever reads the newest five. Older entries moved verbatim to `CHANGELOG-archive.md`, which stays in the git repo but is export-ignored from the distribution tarball.

- Live file keeps the newest 20 entries; everything older is in the archive.
- Tarball installs (no `.git`, no archive file) can browse the archive on GitHub — link in the header.
- `actions/version.md`'s "last 5 releases" read is unaffected (first ~80 lines).

## 0.123.0 — Router Diet: One Routing Table, Help Menu Loads Lazily (2026-07-15)

SKILL.md dropped from ~5,500 to ~2,400 words with zero routing changes. The router used to enumerate the action set five times; now the priority table (with the old Verb Reference's disambiguation folded into its Notes column) and the Action Dispatch table are the only two, and the help menu lives in its own action file that loads only when you actually ask for help.

- Actions bullet list deleted — each action file's own blockquote already carries its description.
- Verb Reference merged into the routing table; every trigger verb and precedence rule preserved.
- Help menu + per-command help moved to `actions/help.md` (new `help` dispatch row).
- Every invocation now loads ~3,100 fewer words of router text before your content is touched.

## 0.122.0 — AI-Report Render-Judge Pass and SVG Design Rules (2026-07-14)

The ai-report action now looks at its own output before shipping: when browser automation is available it serves the report over HTTP, takes full-page light+dark screenshots, and judges them against an explicit layout rubric — catching the dead-gutter columns, SVG label collisions, and buried-lede layouts that read fine in source and broke on screen.

- New mandatory Step 7 "Render and Judge": HTTP serve (never `file://` — it screenshots blank in headless Chrome), full-page light AND dark captures (dark via browser color-scheme emulation), fix-and-re-render loop with two passes minimum when any SVG has text labels; graceful footer disclosure when browser automation is absent
- Six-dimension judge rubric applied to the screenshot, not the source: width usage, table shape, diagram informativeness, emphasis hierarchy, theme robustness, SVG label collisions/clipping
- Data-viz rules for hand-authored SVGs: single-hue ordinal ramps for ordered data, ink-colored labels with identity swatches, above/below label lanes with edge-aware text anchors, stat-tile typography
- Reports commit to one coherent aesthetic direction per report via characterful system font stacks (CDN allowlist unchanged: Tailwind + Mermaid only)
- Matching Red Flags, Common Rationalizations, and Verification Checklist entries; user guide updated to match

## 0.121.1 — Recoverable Runs and Leaner Board Loading (2026-07-13)

Fan-out runs now distinguish “assembled” from “delivered,” so an interrupted review or exploration can resume without cleanup deleting its only result. The board also keeps exact-copy Markdown out of the initial payload until someone actually presses Copy.

- Added `in-progress` → `synthesized` → `consumed` run states, persisted code-review reports, a root deep-explore manifest, and consumed-only cleanup with explicit staging for deleted run paths
- Aligned cleanup's five-pass documentation, changelog-title examples, and prime's interactive questions with their canonical contracts
- Moved raw REQ/UR Markdown into lazy `board-markdown.js`; the current tree's initial `board-data.js` is 43% smaller while generated and live boards still copy exact source text

## 0.121.0 — Tidy-Repo Rename and Safer Layout Planning (2026-07-13)

`file-reorg` is now `tidy-repo`: a clearer name for the same reference-safe repository-layout job, with the old command retained as a compatibility alias. The workflow is tighter about what belongs in a layout pass and more careful around real-world repositories that already have local changes, generators, or platform-sensitive paths.

- Renamed `actions/file-reorg.md` to `actions/tidy-repo.md` and promoted `do-work tidy-repo [path] [plan]` across routing, help, dispatch, README, and next-step guidance
- Added an explicit target-design step, dirty-path overlap handling, generated-source mapping, case-only rename handling, and post-move diff verification
- Made README/CLAUDE edits conditional on actual layout drift; unrelated link fixes, boilerplate rewrites, and permanent link-checker creation are follow-up work instead of mandatory side effects
- Preserved `do-work file-reorg` as a legacy alias so existing prompts keep working

## 0.120.0 — Run Dirs Are Committed, Then Cleaned Up on Consumption (2026-07-13)

Fan-out run directories (`do-work/runs/`) are no longer gitignored transient scratch — they're now committable, so a review or exploration is visible and doesn't get silently lost mid-run. In exchange, the run dir gets deleted the moment its findings are consumed (synthesized and promoted to a report, REQs, or deliverables), which keeps `do-work/runs/` from growing without bound. That whole create → inspect → promote → delete lifecycle is now part of the job, not an afterthought.

- `.gitignore` no longer excludes `do-work/runs/` (`do-work/pipeline.json` stays excluded — it's live state, not work).
- `crew-members/background-agents.md` is the canonical lifecycle: run dirs are committable (step 1) and deleted once consumed (new step 5). The old `.git/info/exclude` append for run dirs is gone.
- `code-review` and `deep-explore` now delete their run/session directory as the final step, after promoting anything worth keeping into `do-work/deliverables/`.
- `cleanup` gains a safety-net pass that sweeps abandoned `Status: complete` run dirs (and leaves incomplete, possibly-resumable ones alone).
- The shared local-ignore snippet still used by `pipeline.json`, the vendored `last30days` engine, and build artifacts moved to a dedicated section in `background-agents.md`; its former callers point there.

## 0.119.0 — Board Drawer Copy Button (2026-07-11)

The Kanban board's ticket drawer gets a Copy button next to Close: one click puts the open REQ's (or UR's) raw Markdown on the clipboard, ready to paste into chat, email, or another ticket without losing headings, checkboxes, or links.

- The data island now ships `bodyMarkdown` beside the pre-rendered `bodyHtml`, so the copy is the ticket's source text, not scraped HTML.
- Transient feedback ("Copied ✓" / "Copy failed") resets on every drawer open; a hidden-textarea fallback covers contexts where the async Clipboard API is missing or denied (file://, plain http).

## 0.118.0 — Cleanup Repoints Doc Links to Moved Files (2026-07-11)

Cleanup's consolidation passes move REQ files around the archive, which used to silently break any doc that linked to them (one consumer repo hit 39 broken prime-doc links). Cleanup now records every move's old → new path and rewrites the referring links itself.

- New `Repoint Documentation Links` step in `actions/cleanup.md`: after all passes, filename-grep tracked markdown outside `do-work/` for each moved file and rewrite link targets from the per-move mapping — preserving `#anchors`, skipping bare prose mentions, tracked files only by design.
- Summary gains a `Repointed: N doc links in M files` line (`Repointed: none` when nothing referenced the moved files, so the step visibly ran).
- The cleanup commit stages the rewritten docs alongside the moves they repair; `docs/cleanup-guide.md` documents the behavior.

## 0.117.1 — Retroactive Descriptive Changelog Titles (2026-07-11)

The descriptive-title convention from 0.117.0 now applies to the whole file: all 152 pre-0.117.0 codename headings ("The Red Pen", "The Court Scribe", …) were rewritten to say what each release delivered. Bodies are untouched — only the heading titles changed.

- Every `## X.Y.Z — The [Codename] (date)` heading from 0.65.0 through 0.115.0 replaced with a short descriptive title derived from that entry's own body.
- Verified no duplicate titles across the file and no codename headings remain.
- CLAUDE.md's "leave pre-0.117.0 entries as-is" note removed — it no longer applies.

## 0.117.0 — Board View Filters (2026-07-11)

The board's By-UR lens rendered the entire archive — after months of history it was an archive dump, not a work view. Every view now filters: a shared search + domain/status bar in the topbar, and an Active/All toggle that hides fully resolved URs by default.

- By-UR lens defaults to Active (URs with at least one unresolved REQ); a footer note counts the hidden resolved URs, and All brings them back.
- Shared filter bar applies to whichever view is active: search matches REQ/UR ids and titles, domain and status selects populate from the data. Column and UR counts read "shown / total" while filtering; the calendar hides days with no matches.
- A search hit on a UR header keeps its whole group visible (domain/status still filter the cards inside).

## 0.116.1 — Clear Questions in Review-Work Follow-Ups (2026-07-11)

0.116.0 required cold-reader question authoring in work.md's follow-ups but missed the copy-paste sibling: review-work's ambiguous-requirements follow-ups emit the same `Recommended:`/`Also:` template. A grep for every `pending-answers` authoring site found this one remaining gap.

- `actions/review-work.md` ambiguous-requirements follow-ups now load `crew-members/clear-questions.md` and author Open Questions for a cold reader (gloss shorthand, state why the decision is the user's — Principle 7), matching work.md Step 8.

## 0.116.0 — Escalated Questions Explain Themselves (2026-07-11)

Escalated questions were reaching the user written in builder shorthand — technically asked, practically unanswerable. Now clarity is enforced at both ends: builders author Open Questions for a cold reader, and clarify rewrites what slips through.

- `actions/clarify.md` Step 3 now loads `crew-members/clear-questions.md` and rewrites stored question text to its contract instead of rendering it verbatim.
- New clear-questions Principle 7: an escalated question must say why the decision is the user's — the rule that forced the escalation and what silently deciding would have cost.
- `actions/work.md` Step 8 and the follow-up template in `actions/work-reference.md` require Open Questions destined for clarify to meet the contract at authoring time.

## 0.115.0 — Board Flags Invalid REQ Statuses (2026-07-10)

The Kanban board now marks a REQ whose `status:` is outside the schema vocabulary as *invalid* — red status, an INVALID pill on the card, and a drawer note telling you exactly how to fix it — instead of letting it blend in with normal blocked tickets. Came out of triaging review feedback: the live-tree bucketing test contradicted the board's own deliberate catch-all and would have failed on any off-vocabulary status.

- `bucketColumns` flags off-vocabulary tickets (`StatusUnrecognized`), and its warning now carries the fix prompt (edit `status:` per the Schema Read Contract, or run `do-work forensics`).
- New forensics check 11 sweeps queue/working/archive for unrecognized statuses — the mechanical fix path the board's warning points at.
- `TestLiveTreeColumnBucketingMatchesStatus` now asserts the real invariant (unrecognized statuses legitimately live in Needs-input *when flagged*), plus a seeded synthetic regression test so the live queue can't mask it.

## 0.114.0 — Retire the Weekly-Signal-Diff Prompt (2026-07-10)

Retired the `weekly-signal-diff` prompt from the library. It graduated into the consumer project's own `wsd-skill` (as `daily-signal-diff`, driven by the `wsd-full` / `wsd-go` / `wsd-refresh` family) months ago — the shipped copy was a stale duplicate that every `do-work update` kept reinstalling.

- Removed `prompts/weekly-signal-diff.md` and `prompts/weekly-signal-diff-personal.md`; dropped their rows from `prompts/README.md`.
- `decisions/imported-specs/2026-04-17_improve-weekly-diff-skill.md` gained a Status footer recording the removal; changelog history stays as-is.
- The `**Runnable:**` header key in `actions/prompts.md` is generic and remains — it just no longer has a shipped opt-out example.

## 0.113.2 — Drawer Formatting for Questions and Prose (2026-07-10)

The drawer was mashing a REQ's Open Questions into one run-on paragraph and stretching prose across the whole panel. Both readable now.

- `Recommended:` / `Also:` / `Value:` / `Risk:` / `→` continuation lines render on their own lines instead of lazily merging into the question sentence (fenced code blocks stay verbatim).
- Markdown body text caps at ~90 characters per line, so a wide drawer no longer means 200-character lines.

## 0.113.1 — Notes Strip Parses Only Bullet Lines (2026-07-10)

The Notes strip was reading a real `notes.md` as eighteen notes when it held two. Only bullet lines are notes now.

- The `#` heading, the prose preamble, and horizontal rules are skipped instead of rendered as notes.
- `<!-- ... -->` comment blocks are stripped **before** the bullet test — that's where pruned entries get parked, and their bullets were resurfacing on the board.
- `do-work roadmap` and `do-work note` carry the same rule, so every reader of `notes.md` agrees on what a note is.

## 0.113.0 — Board Dependency Graph: Ready vs Waiting (2026-07-10)

The board finally draws the dependency graph it was already parsing. Pending now separates what you can pick up right now from what's still waiting on an upstream REQ, and every card tells you how much is waiting on *it*.

- **Ready vs. Waiting.** The Pending column splits in two. When nothing is waiting, it stays a flat list — no new headers for a queue without dependencies.
- **Unblocks N.** A card carrying that badge is the one to work on: N unresolved REQs are waiting for it. The full list is in the detail drawer.
- **Dangling dependencies are now loud.** A `depends_on` pointing at a REQ that isn't in the tree fails closed (the dependent stays waiting, never quietly ready) and raises a data warning — it can never self-resolve.
- Dependency chips show met (struck through) vs. unmet (amber), and the drawer lists each dependency with the status that decides it. `cancelled` never satisfies gating, matching the work loop.
- `do-work board summary` now prints the ready / waiting breakdown.

## 0.112.0 — Notes Strip on the Kanban Board (2026-07-10)

Your `do-work note` hints now show up on the Kanban board, not just in `do-work roadmap`. They sit in a collapsible Notes strip above the columns, so the thing you told yourself to check next is visible while you're staring at the queue.

- `do-work board` reads `do-work/notes.md` and renders each line with its date, in append order.
- The strip stays visible in the calendar view too, and disappears entirely when there are no notes.
- Notes render as plain text, never Markdown — they're hints, not tickets, so they get no column, no calendar entry, and no detail drawer.
- Serve mode watches `notes.md`, so appending a note and reloading the page shows it.

## 0.111.0 — Versioned Changelog Entries in Target Repos (2026-07-09)

Changelog entries in unversioned repos came out keyed by date alone, so nothing told you whether an entry was a typo fix or a rewrite. Every entry now carries a version and a date, and the number is earned — bumped by what the change actually did to people using the code.

- Entry key is always `## X.Y.Z — The [Codename] (YYYY-MM-DD)`
- Version source resolves in order: a version file in the repo (bumped and staged with the REQ commit), release tags (read, never created — a tag is a human's release call), or the changelog's own counter seeded at `0.1.0` for repos with no version at all
- Bump size reads the delivered change: breaking a consumer is major, a new user-invocable capability is minor, everything else is patch. Ties break downward; below `1.0.0` a breaking change bumps the minor, so a seeded repo never silently promotes itself to a `1.0.0` release
- Fixes a duplicate-header bug on the versioned path, which reused the repo's current version for every entry instead of bumping it
- Guards added for disagreeing version files (leave them alone, fall back to the counter, report it) and for out-of-band releases (bump from whichever source is higher)
- The commit's "did we actually stage an implementation?" check now knows the version file is bookkeeping, not implementation — so a lone version bump can't masquerade as delivered work

## 0.110.0 — Work Pipeline Writes Target-Repo Changelogs (2026-07-07)

This changelog was the only one do-work ever kept — every target repo's history lived in commit messages nobody rereads. Now the work pipeline writes a changelog entry in every repo it works in, by default, in the house voice (picked from a six-voice side-by-side style lab over four real entries).

- New **Changelog Entry Procedure (Step 9)** in `actions/work-reference.md`: house-style contract (value-first lead + technical bullets), `## YYYY-MM-DD — The [Codename]` keys for unversioned repos, the repo's own version when it has one — never invented
- Bootstrap when `CHANGELOG.md` is missing; an existing changelog in a different format wins over the house voice
- Successful REQs only — failed and cancelled work gets no entry; `CHANGELOG.md` joins the explicit staging list and doesn't count as implementation in the commit validation check
- Wired into `actions/work.md`'s Commit Phase; entries load `crew-members/anti-slop.md` like any human-facing artifact
