# Changelog

What's new, what's better, what's different. Most recent stuff on top.

> Covers **0.121.1 onward**. Older releases live in dated archives, each named for the range it holds: [0.110.0–0.121.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-07-13-up-to-v0.121.0.md) · [0.65.0–0.109.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-07-07-up-to-v0.109.0.md) · [0.50.0–0.64.1](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-04-13-up-to-v0.64.1.md) · [0.1.0–0.49.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-04-07-up-to-v0.49.0.md). They are tracked in git but export-ignored from the distribution tarball, so a tarball install browses them at <https://github.com/knews2019/skill-do-work/tree/main>.
>
> Keep this note short: `actions/version.md` reads only the first ~80 lines to find the newest 5 entries.

---

## 0.193.0 — Explicit Presentation Command Ownership (2026-08-15)

Completed-work presentation now has one discoverable command per outcome: detailed evidence in `ai-report`, cross-project portfolios in `present-work`, and explicit source-only walkthroughs in `present-video`.

- Aligns router aliases, argument hints, help, tutorials, README examples, and completion guidance on the same three-command map
- Makes archive and human-artifact guardrails condition-based so future readers and presentation producers inherit them without stale caller lists
- Adds durable ownership, evidence, portfolio-branch, source-video, safety-order, and shipped-inventory contracts while retiring ambiguous detail/video routes

## 0.192.1 — Inline PNG File Views (2026-08-15)

Captured screenshots opened from the live board no longer spill binary bytes across the browser. Byte-detected PNGs render inline while every other file keeps the inert text response that protects the board origin.

- Returns `image/png` only for PNG bytes, never from a filename extension alone
- Keeps HTML, SVG, misleading `.png` files, and ordinary documents on `text/plain` with the existing loopback, containment, size, and `nosniff` guards
- Adds end-to-end RED/GREEN coverage for a real encoded PNG, byte preservation, and the safe fallback

## 0.192.0 — Explicit Source-Only Video Walkthroughs (2026-08-15)

Completed work can now be expressed as a deliberate Remotion source walkthrough without making video a hidden side effect of reports, portfolios, or completion. The new action scales its four-scene story to verified evidence and leaves previewing entirely in the user's foreground terminal.

- Adds an explicit-only `present-video` action and guide for completed URs and REQs, with concise successful skips for trivial or genuinely non-visual work
- Defines a complete Problem → Solution → Architecture → Value React/TypeScript source tree with proportional frame math, module-level registration, and evidence-backed claims
- Uses a package-local foreground Studio preview command while prohibiting automatic installs, external assets, fixed ports, readiness sleeps, browser launchers, and rendered media

## 0.191.0 — Portfolio-Only Present Work (2026-08-15)

`present-work` now does one job: turn the successful archive into a cross-project portfolio. Item-specific calls stay non-writing and point directly to the detailed report and video commands instead of quietly generating competing artifacts.

- Restricts writing to `present-work all|portfolio`, with archive safety, terminal-success issue disclosure, and evidence-backed value language
- Makes bare and UR/REQ invocations non-writing guidance paths with exact `ai-report` and `present-video` replacements
- Refreshes one canonical portfolio summary and guides optional byte-identical, timestamped, no-clobber snapshot preservation without deleting prior artifacts

## 0.190.0 — Canonical Detailed Work Reports (2026-08-15)

Completed-work reporting now has one detailed HTML path for visual and non-visual changes. `ai-report` adapts its evidence to the work while shared archive rules keep every presentation honest and collision-safe.

- Adds one completed-work presentation reference for terminal-success resolution, safe archive reading, evidence provenance, merge/current-code inspection, and no-overwrite publication
- Preserves authentic screenshots, SVG annotations, real before/after comparison, generated-image provenance, responsive layout, and light/dark rendered verification for visual work
- Adds architecture, commit, current-code, test, and operational evidence for backend, refactor, infrastructure, and other work where UI captures are not expected

## 0.189.19 — Case-Exact Late Root Justfile Contracts (2026-08-15)

Four late aggregate assertions still opened `Justfile`, passing on common macOS filesystems but failing against the tracked lowercase root file on case-sensitive systems.

- Changes every remaining live root assertion input to exact tracked `justfile` casing
- Adds a filesystem-independent ratchet that requires all four exact assertion patterns and rejects wrong-case, missing, or duplicate inputs
- Preserves intentional prose and filename-variant fixtures while keeping the aggregate final marker and canonical maintainer gate green

## 0.189.18 — Modular Framework-Free Board Client (2026-08-15)

The queue board's single large browser source now has explicit responsibility boundaries while preserving the exact framework-free runtime. One private shell assembles eight ordered closure fragments for both static and live pages.

- Preserves the pre-change client byte-for-byte through a fixed Go manifest with no fragment requests, modules, globals, frameworks, or dependencies
- Mutation-locks authored inventory, execution order, exact boundaries, raw/canonical marker invariants, assembled syntax, and static/live byte parity
- Retargets source ownership and contract checks to the exact fragment owners while keeping the strict JavaScript behavior lane intact

## 0.189.17 — Structured Stray Request Verification (2026-08-15)

Misplaced request files could trigger a board warning yet disappear from the forensics report. Verify now forwards the board's canonical structured evidence while recognizing legitimate review follow-ups under closed User Requests.

- Emits one read-only, non-fixable verify finding for every retained off-section REQ path, regardless of status
- Keeps strays out of cards and ordinary probes, with direct tests rejecting warning parsing or a second filesystem walk
- Exempts only exact review-generated live members from the archived-UR anomaly and replaces reopen advice with a stays-closed remedy

## 0.189.16 — Closed-UR Review Follow-up Lifecycle (2026-08-15)

Standalone review could read a closed User Request and leave later work without a safe return path. Closed URs now stay stationary from review context through completed follow-up placement.

- Makes archived input context-only in every review mode while preserving same-UR queued follow-ups
- Returns successful review-generated REQs to an already-archived UR folder in place
- Mutation-locks the review marker, same-UR archived-folder existence, stationary-folder rule, and active-UR branch bypass

## 0.189.15 — Visible Unavailable Hotspot Evidence (2026-08-15)

Hotspot rankings could look complete while silently omitting tracked paths unavailable in the current worktree. Those paths now remain visible without distorting the numeric ranking.

- Separates measured hotspots from sorted unavailable churn-bearing paths
- Keeps numeric ordering and `topCount` behavior while rendering every unavailable path uncapped
- Shows known commit counts with `NOT-MEASURED` lines and scores, plus a visible incomplete-ranking warning

## 0.189.14 — Canonical Maintainer Verification Gate (2026-08-15)

Repository health no longer depends on remembering separate shell and Go commands. One local, export-ignored gate now owns the complete maintainer verification path and its failure contract.

- Pins the audited Go and ShellCheck versions, lints every tracked shell file at warning severity, and runs the aggregate once
- Vets and tests both Go modules, plus the strict board JavaScript lane whenever Node is available
- Proves exact-once execution and nonzero propagation with recursion-safe fixture shims; CLAUDE and Just remain thin delegates

## 0.189.13 — Single-Owner Baseline Suites (2026-08-15)

Two baseline child suites ran twice through identical paths, adding hand-back time without adding evidence. Each child now has one required owner while every distinct late aggregate check remains intact.

- Prescribed-shell behavior runs through the staged-skills contract instead of also running directly from the aggregate
- Shipped-package references remain aggregate-owned without a second mandatory standalone hand-back invocation
- Failure propagation, the late installer suite, and the aggregate's final pass marker remain covered

## 0.189.12 — Strict JavaScript Behavior Lane (2026-08-15)

The board's incident-sensitive JavaScript probes could all skip when Node was unavailable while the Go suite still reported success. A maintainer-only lane now distinguishes an intentional consumer skip from a false-green zero-probe verification run.

- Centralized all seven Node behavior probes behind one counted runner and stable test prefix
- Added a canonical strict entrypoint that fails after an otherwise green run when no JavaScript probe started, while ordinary package tests still skip cleanly without Node
- Replaced state source-shape claims with executable production behavior for recent-window, By-UR, empty-state, and confirmed testing transitions

## 0.189.11 — Listener-Anchored Live Board Authority (2026-08-15)

The live board trusted a matching Host and Origin even though both values came from the request. Every production route now sits behind a post-bind authority gate derived from the actual listener and assigned port.

- Concrete, loopback, LAN, IPv6, and wildcard binds have explicit normalized Host policies with no request-time DNS lookup
- Wildcard listeners accept only the connection's concrete numeric local destination, preserving intentional LAN access without trusting arbitrary DNS
- Testing writes retain their existing guards and now compare normalized HTTP Origin and request authorities

## 0.189.10 — Recoverable Static Board Publication (2026-08-15)

A failed static-board refresh could leave new card data beside old Markdown and HTML, producing a plausible but internally mixed bundle. The three public files now publish through one bounded all-or-recover operation.

- Builds and stages all three payloads privately before touching public targets
- Holds unique backups until publication succeeds, restoring exact pre-invocation bytes after a handled failure
- Preserves unrelated output entries and removes private staging/backup residue after success or completed rollback

## 0.189.9 — Public Vocabulary Parity (2026-08-15)

Public work aliases and testing-status aliases could drift from their runtime mirrors while every suite stayed green. The documented owners and executable readers now agree, with mutation-sensitive seam checks to keep them together.

- The work guide is the sole public alias inventory, exactly matched by the core router; README points to it instead of maintaining a third list
- The testing-status schema now includes both spaced aliases already accepted by the board normalizer
- Queue summaries surface dependency-cycle holds, and the schema gloss no longer freezes a stale field list

## 0.189.8 — Route-A Public Write Boundary (2026-08-15)

The new README boundary still implied every implementation had a declared Scope section, but simple Route A work deliberately does not. Public guidance now names both valid boundaries instead of flattening them into one.

- Routes B/C source writes are bounded by their declared `## Scope`
- Route A source writes are bounded by the focused REQ text

## 0.189.7 — Accurate Public Write Boundaries (2026-08-15)

The existing-project README claimed do-work only wrote beneath its queue directory, hiding both managed installation paths and request-authorized implementation edits. The public guidance now states the real boundaries before adopters trust the suite.

- Distinguishes reviewed four-skill/Just/settings install paths from durable `do-work/` queue state
- Makes project-source writes explicit: only invoked REQ implementation, bounded by that REQ's declared Scope

## 0.189.6 — Case-Exact Contract Suite Paths (2026-08-15)

Case-insensitive development filesystems hid that the late contract checks opened a nonexistent `Justfile` on Linux. The suite now names the tracked lowercase file exactly and reaches its final checks everywhere.

- Changed both Kanban shutdown-check inputs in `_dev/tests/contract-regressions.sh` from `Justfile` to `justfile`
- Full contract regression suite now completes at exit 0 after the repaired block

## 0.189.5 — Consolidated Archive Lesson Links (2026-08-15)

Closing UR-039 moved its completed requests into their durable UR archive, but two shipped lesson links kept pointing at the old root paths. The links now follow the consolidated archive so reference verification stays green.

- Updated the REQ-173 and REQ-174 lesson URLs in the do-work update prime to include `archive/UR-039/`
- Restored the shipped-package reference contract after the UR-039 lifecycle move

## 0.189.4 — Goldmark-Aligned Board Question Fences (2026-08-15)

Open Questions could lose their visual breaks after prose that merely looked like a code fence. The board preprocessor now follows Goldmark's marker-specific info rule, so guidance stays readable without changing real code blocks.

- Invalid backtick-info lookalikes remain ordinary prose, allowing `Recommended:` and `Also:` lines to receive their intended hard breaks
- Focused renderer tests preserve byte-verbatim handling for valid backtick and tilde fences

## 0.189.3 — Churn Exclude Fix and Calibration Fallback (2026-08-14)

Two PR-review findings on the new audit machinery: an excluded-but-still-tracked file that git flags as a copy source was treated as dead, silently handing its whole history to the surviving copy and inflating its hotspot score; and the manual inventory fallback gave the calibration gate nothing to derive FLAG = max(floor, p95) from.

- `audit-metrics` churn now judges copy-source aliveness against the unfiltered tracked set and applies excludes only to report output; lock-in test pins it (excluded-but-live source keeps its history)
- The reference's manual fallbacks gain nearest-rank distribution blocks (lines and words: count/median/p90/p95/max), so calibration completes without the Go toolchain

## 0.189.2 — Scope-Drift Check Parses Annotated Headers (2026-08-14)

The Step 5.5 scope-drift check silently skipped any REQ whose touch-list header carried a parenthetical — the exact silent self-disable its own guard exists to prevent. Found live during REQ-178's review.

- `tools/checks/scope-drift.sh` now tolerates trailing annotations before the colon at both match sites: an annotated header parses, or FAILs loudly when nothing parses — never SKIPs; absent Scope still SKIPs (Route A contract preserved)
- Three lock-in probes pin the contract in the regression suite

## 0.189.1 — Maintainability-Audit User Guide (2026-08-14)

The new audit shipped without a user-facing walkthrough of its loop. Now the guide exists and the stale trigger attribution is gone.

- New `docs/maintainability-audit-guide.md` in do-work-toolbox: run → calibrate → read → triage → capture → build → re-audit, plus lock-in limits and waivers, linked from the action's description
- code-review's guide no longer shows `audit codebase` as its own invocation (that phrase now runs the maintainability audit)

## 0.189.0 — Maintainability-Audit Action (2026-08-14)

The toolbox's findings family could produce reviews and receive triage, but nothing measured repo health repeatably. A new grounded, interactive, read-only audit closes the loop: measure with calibrated bands, judge only hotspots, emit refutable finding classes, track deltas across runs.

- New `maintainability-audit` action + reference companion in do-work-toolbox: grounding → calibration gate → metrics via the audit-metrics tool (manual fallbacks included) → hotspot-scoped judgment → root-cause classes → persistent `do-work/audits/` report whose Findings section pastes into `do-work-toolbox validate-feedback`
- Lock-in limits (pinned-at-worst regression ceilings) are proposals only — accepted ones land as lock-in tests through the normal capture flow
- The `audit codebase` trigger moves from code-review to the new action; code-review keeps `code-review` / `review codebase`

## 0.188.0 — Audit-Metrics Measurement Tool (2026-08-14)

The upcoming maintainability audit needs numbers, and hand-run wc/find/git pipelines are fragile and expensive per run. A vendored Go tool now answers everything deterministic — the audit action will paste its tables instead of re-deriving them.

- New `skills/do-work-toolbox/tools/audit-metrics/` module (zero dependencies, queue-kanban conventions): `inventory`, `folders`, `churn`, `hotspots` subcommands emitting pasteable markdown
- WATCH/FLAG bands come only from flags (strict-greater edges); exclude list is caller-supplied prefixes; output states what was excluded
- Churn is rename- AND copy-normalized (`-M -C --find-copies-harder` with dead-copy-source reassignment) — reproduces `git log --follow` across the skills/ restructure, where plain aggregation splits history onto dead paths; shallow clones are reported, never silently truncated
- 10 lock-in tests on real git and real `--depth 1` clone fixtures; in-dir `prime-audit-metrics.md` routing index

## 0.187.1 — Validated Runtime Boundaries (2026-08-13)

Timeouts now clean up the process trees they start, directory installs publish only complete payloads, and report images prove they came from the current invocation.

- Gives the stock-macOS blocked-check fallback an isolated process group with verified, escalating group cleanup while preserving GNU timeout behavior and ordinary exit statuses
- Publishes the complete `last30days` subtree as a validated adjacent rename transaction with restoration after copy, publication, or interruption failure
- Stages every report image per invocation, retains every parallel status, and keeps the full-host agentic backend behind the exact explicit opt-in
- Adds 22 runtime fixtures plus a durable linked lesson covering process ownership, directory transactions, artifact freshness, and authority boundaries

## 0.187.0 — Decision Revalidation Across Queued Work (2026-08-13)

Reversed decisions can now be checked once against every unfinished queued REQ, turning a changed premise into an evidence-backed reconciliation list without mutating queue state.

- Adds read-only `verify-requests --against` scans for superseded ADRs and answered builder-decision follow-ups, including repeated sources in one queue pass
- Ranks direct references and restatements as likely affected and defensible semantic dependencies as possibly affected, with quoted evidence and reconciliation commands
- Batches genuine clarify overrides, auto-scans queues up to 10,000 words, asks before larger automatic scans, and explicitly excludes claimed or archived work

## 0.186.37 — Parallel Writer Worktree Hand-Back (2026-08-13)

Overlapping parallel writers now get isolated before they edit, and completed work comes back through one verified integration path instead of being left on side branches.

- Makes explicitly overlapping declared write ownership the trigger for worktree-and-branch isolation across every package-local background-agent guide
- Requires serial merge attempts, conflict and semantic reconciliation, merged-state verification, and preservation of unsafe branches without force
- Adds the Step 6 reminder while leaving auto-wave's stronger one-worktree-per-builder contract and display-only `write_set` untouched

## 0.186.36 — Narrow Defensive Surface Ratchet (2026-08-12)

Incident-backed guidance can return without fighting a blanket section ban, while the generic text removed by REQ-168 stays deleted. The dated defensive-surface audit is now historical evidence instead of a registry every future shell file must extend.

- Replaces exhaustive shell/prose inventory enforcement with exact sentinels for the removed generic tables and arbitrary commit-size warning
- Allows future `Red Flags` and `Common Rationalizations` sections when a concrete incident earns them
- Freezes the audit after REQ-171 and removes 30 net lines of ongoing maintenance surface

## 0.186.35 — Maintainer Doc Split Into Core Rules Plus Prime Files (2026-08-12)

The root CLAUDE.md went from a 202-line everything-file to a short core in the maintainer's own vocabulary, with domain detail lazy-loaded from prime files. A session now reads the deep material only when working in that domain.

- New `_dev/primes/` (export-ignored) holds three primes: action-file conventions, prescribed-shell traps, and Kanban board tool rules — each moved intact from CLAUDE.md
- Core gains a plain-words glossary, personal coding preferences (YAGNI, focused tests, match ceremony to the task), a Verify section (exit code zero is the only proof), and push-back-then-continue communication rules
- The Kanban board three-write-surfaces sentence stays in CLAUDE.md (lock-in test unchanged); live pointers in `_dev/tests/contract-regressions.sh` and `.gitattributes` now cite the primes


## 0.186.34 — Prescribed Shell Becomes Executable (2026-08-11)

Reusable shell mechanics now ship as executable, fixture-tested scripts instead of copied multi-line Markdown. Callers keep local intent and policy while semantic traps have attributable runtime tests.

- Promotes 17 multi-line blocks into 11 scripts across core, knowledge, and toolbox while retaining 21 inline-residue and 2 Go-owned blocks
- Adds 11 named behavior cases for atomic publication, merge display, screenshot races, portable timeouts, protected inventory, literal exact deletion, memory hooks/recall, and toolbox installs
- Shrinks shipped Markdown by 303 nonblank shell-body lines and ratchets executable homes, package paths, defensive evidence, and the full 34-path scope


## 0.186.33 — Findings Close with Proof (2026-08-11)

Review and triage findings can no longer close on a bare patch or unrelated green tests. Finding-origin work must carry matching fail-before/pass-after evidence or delete the exact surface, while defensive fixes must justify the surface they add.

- Establishes one canonical Finding-Closure Ratchet independent of `tdd`, review score, and unrelated test results
- Adds the one-paragraph earned-defense question to the always-loaded simplicity guardrail and keeps triage behavior citation-sized
- Makes every shipped review-generated REQ template emit compatible Red-Green Proof and dynamically ratchets future producers


## 0.186.32 — Root Markdown Fence Info Validity (2026-08-11)

The shipped reference classifier now agrees with the pinned Markdown renderer when a root backtick-fence candidate contains a forbidden backtick in its info string. Invalid openers remain visible prose instead of masking later links.

- Shares one marker-aware info-string predicate across root, list, and paragraph-state classification
- Preserves valid tilde-fence behavior while consolidating the former list-only exception
- Adds Goldmark-differential root/list/tilde fixtures and keeps the full shipped-reference contract green


## 0.186.31 — BOM-Aware Just Collision Scans (2026-08-11)

Suite installation now recognizes a reserved Just recipe even when the file begins with a UTF-8 BOM and `just` is unavailable. Collision checks stay byte-preserving and still reject before confirmation or mutation.

- Removes exactly one leading BOM only from the first line's classification view
- Keeps both distributed managed-section helpers byte-identical without broader encoding normalization
- Replays the no-Just installer path and verifies byte, mode, settings, module, and Git-state preservation

## 0.186.30 — Unique Screenshot Installation Copies (2026-08-11)

Concurrent screenshot captures can no longer swap bytes between verification and publication. Each dispatch now verifies and installs its own adjacent temporary copy while retaining no-clobber and recovery semantics.

- Allocates the permanent-copy staging file with `mktemp`, then compares and hard-links that exact unique path
- Preserves the staged loser and existing permanent destination on collision while cleaning only the dispatch-owned temporary copy
- Coordinates two different staged byte sequences around copy, verification, and install to prove winner ownership, loser recovery, cleanup, and ordinary capture behavior

## 0.186.29 — Retry-Safe Screenshot Source Cleanup (2026-08-11)

A verified screenshot no longer becomes unusable just because its staged source could not be removed. Capture reports the leftover while keeping the installed asset valid and the no-clobber boundary strict.

- Makes staged-source deletion warned best-effort after byte verification and permanent installation
- Replays the exact failed-removal path, including preserved bytes, temporary-copy cleanup, and a strict later collision

## 0.186.28 — Feedback Prices Added Defense (2026-08-11)

External feedback can no longer smuggle speculative guards and fallbacks through a valid finding. Triage now verifies whether proposed defensive surface earns its long-term cost before recommending Accept.

- Applies the exact incident-and-surface-cost rubric to remedies that add guards, fallbacks, retries, validation, rules, or warning apparatus
- Routes unearned defense to Push back and unresolved real trade-offs to Discuss while leaving direct fixes, deletions, and simplifications unchanged
- Shows `Surface-cost: N/A / Earned / Flagged` per finding and pins scope, verdict impact, and output visibility in aggregate contracts

## 0.186.27 — Defensive Surface Earns Its Keep (2026-08-11)

Shipped defenses now have an incident-and-test ledger instead of being accepted on instinct. Generic warning apparatus that duplicated existing action contracts has been removed and ratcheted out.

- Audits all shipped shell sources and explicit Rules/Rationalizations/Red Flags/Warnings/recovery surfaces, with keep/delete evidence and behavior-changing candidates called out
- Removes 96 lines of decorative review guidance, duplicate rationale, generic inspect warnings, and an arbitrary commit-size heuristic without changing action steps or permissions
- Adds a focused completeness/deletion probe to the aggregate suite while preserving every tested data-loss, recovery, secret, parser, hook, and install defense

## 0.186.26 — Canonical Prescribed Shell Primitives (2026-08-11)

Shell safety rationale now has one shipped home instead of drifting across copied action prose. Callers remain independently runnable while future corrections have a single place to land.

- Adds a core guide covering eight recurring shell primitives plus a durable audit of canonical homes, former copies, and divergent variants
- Replaces repeated cross-package explanations with package-safe pointers while preserving caller commands, policy gates, and fixed semantics
- Ratchets the consolidation through focused pointer/stale-phrase checks, the shipped shell-fence harness, and the aggregate contract suite

## 0.186.25 — Fail-Soft SessionStart Hook (2026-08-11)

Session start keeps its cross-session status signal even when installed version metadata or the queue directory is missing. The hook is shorter, and fixture coverage now pins every fallback that previously died before the banner.

- Replaces the strict-shell `grep | sed` failure path with minimal empty-value defaults while preserving the anchored hook command and exact stdout contract
- Covers the happy path, missing version file, reformatted version label, and missing queue directory through the real copied hook
- Invokes the new SessionStart behavioral probe from `_dev/tests/contract-regressions.sh`

## 0.186.24 — Shipped Shell Lint Harness (2026-08-11)

Prescribed shell now fails in the contract suite instead of later in a consumer repository. A self-testing probe covers the full shipped shell surface while keeping diagnostics attributable and optional tooling optional.

- Extracts Markdown `bash`/`sh` fences with valid indentation, narrowly neutralizes prose placeholders, and remaps Bash and ShellCheck diagnostics to their source file and line
- Lints complete shipped shell sources without snippet exclusions and degrades to `bash -n` with a note when ShellCheck is unavailable
- Proves the failure path with an indented malformed fixture and runs both negative and clean-tree checks from `_dev/tests/contract-regressions.sh`
- Makes the board action's prescribed repository-root directory change fail fast, closing the real warning exposed by the complete scan

## 0.186.23 — Collision-Safe Screenshots and Lossless Suite Guards (2026-08-11)

Screenshot capture now preserves every dispatch and every image without clobbering recovery evidence. Suite validation, configuration reconciliation, and preflight reporting also retain the exact file and path shapes they inspect.

- Allocates an exclusive staging directory per screenshot dispatch, installs ordinal asset paths without overwriting, and treats empty-directory cleanup as best-effort after verification
- Parses dirty Git paths through NUL-delimited porcelain so spaces, quotes, and glob characters survive preflight unchanged
- Rejects directory-shaped `SKILL.md` entries before installation and verifies regular files in installed-module fixtures
- Preserves unrelated empty Stop-hook wrappers in both jq and Python reconciliation paths
- Resolves the real Justfile directory-entry spelling before backup, write, and failed-install recovery

## 0.186.22 — Status-Colored Board Cards (2026-08-11)

Board cards now separate workflow states at a glance without tinting the whole surface. Mixed-status By-UR groups use the same restrained cues as the column view.

- Adds a 3px semantic status rail and softly tinted, written status pill to shared request cards
- Maps pending, claimed, blocked/failed, completed, cancelled, and invalid states to the existing amber, blue, red, green, and gray palette
- Preserves neutral card bodies, cancellation strike-through, hover/focus feedback, responsive wrapping, and light/dark themes

## 0.186.21 — Visible Failures and Resolvable Archived Attempts (2026-08-10)

Failed work now carries its recorded diagnosis into the board drawer, and legacy failures already consolidated inside closed UR folders have an explicit resolution path. Detail drawers also stop repeating a record's matching leading title without changing the exact Markdown copied from disk.

- Projects `error` and present-only normalized `error_type` values into failed REQ drawers, retaining invalid provenance and warnings without fabricating a missing type
- Removes only a matching leading H1 from rendered REQ and UR drawer bodies while keeping nonmatching headings and verbatim Copy payloads intact
- Lets an explicitly named failed REQ inside `archive/UR-NNN/` be confirmed and cancelled in place without bulk-discovering, moving, or reopening closed UR history

## 0.186.20 — Verified Screenshot Staging Cleanup (2026-08-10)

Screenshot captures dispatched to subagents no longer leave their hidden staging copies behind after the permanent UR assets are safe. Failed copies remain recoverable and visible in the capture report.

- Restores the modular dispatcher's `.pending-assets/` handoff and links it directly to capture Step 4 as the cleanup owner
- Copies through a temporary destination, verifies byte equality, and renames before removing the exact staged source and empty staging directory
- Preserves and reports staged screenshots on copy, verification, rename, or cleanup failure, with regression coverage for the lifecycle contract

## 0.186.19 — Exact Installer Recovery and Parser Guard Completion (2026-08-10)

Failed suite installs now restore both managed filesystem bytes and the exact pre-install Git index, including staged/unstaged distinctions. The same release closes the remaining delimiter, Markdown-fence, and version-file validation gaps found in the distribution guards.

- Snapshots the Git index before managed unstaging, arms recovery first, and restores the index atomically on failures and handled signals
- Compares managed bytes, ordinary and cached diffs, and porcelain status across staged-only, unstaged-only, partially staged, interrupted, cancelled, and successful installer fixtures
- Keeps embedded quote/backtick characters inside matching Just triple delimiters from terminating multiline defaults early
- Rejects root fence lookalikes with trailing text and version files with any bytes beyond their single newline-terminated semantic version

## 0.186.18 — Release, Parser, and Distribution Guard Fixes (2026-08-10)

Consumer releases no longer risk rewriting the installed do-work version, and the update/install path now rejects version drift before it can escape recovery. The same review pass closes the remaining Just/CommonMark parsing gaps and turns helper/help alignment into enforced contracts.

- Routes documented version/update phrases before generic `check`, removes the unsafe `next-version` accelerator, and excludes installed suite metadata from project version selection
- Validates root, core, and runtime action versions before and after installation, preserving all-or-recover behavior
- Detects reserved recipe headers with multiline defaults and honors CommonMark link adjacency and list-fence container boundaries
- Enforces byte identity across installer, replacer, and validator mirrors and documents `board [serve|static|summary|cli]`

## 0.186.17 — Core Help Lists Every Sibling Subcommand (2026-08-10)

`do-work help` used to name one sample command per extension package, which read like the whole menu and left the other sibling commands undiscoverable. It now names all 26, so knowing a command like `tidy-repo` is enough to find the package that owns it.

- Replaces the sampled `board`/`bkb`/`code-review` lines with the full subcommand list for each sibling package
- Points at `<package> help` for usage detail instead of restaging a sibling's menu in core
- Adds a help Rules bullet requiring the list to move in the same commit as any sibling command change

## 0.186.16 — UR Archive Lesson Link Repointing (2026-08-10)

Closing UR-031 moved its completed history into one self-contained archive folder, and the updater prime's durable lessons now follow it. Shipped-reference validation no longer depends on stale loose-archive paths.

- Repoints five canonical REQ-136/137/138/144/146 lesson URLs into `do-work/archive/UR-031/`
- Preserves each lesson anchor and the updater prime's existing inline-only lesson format
- Keeps the completed 27-REQ UR archive self-contained and the shipped Markdown reference contract green

## 0.186.15 — Complete Markdown Link and List-Fence Classification (2026-08-10)

The shipped-reference guard now reaches the last known Markdown visibility edges: even-parity relative links stay discoverable, escaped brackets inside labels no longer hide live links, and fenced code inside list items stays unpublished. Target policy remains unchanged.

- Extracts relative link targets behind zero, two, or four backslashes while retaining odd-parity masking
- Treats complete live links as one region so escaped opening brackets inside labels are not mistaken for independent escaped links
- Tracks backtick and tilde fences opened by bullet or one-to-nine-digit ordered list items, including attached info strings and nested indentation
- Preserves live post-fence and list-paragraph continuations plus root-level indented-code masking
- Adds exact RED/GREEN fixtures with normalized target order, source-length, newline-offset, full-contract, and distribution coverage

## 0.186.14 — Just Ordinary Multiline-Backtick State (2026-08-10)

Safe custom Justfiles no longer look like reserved-recipe collisions when an ordinary single-backtick command spans physical lines. Real definitions around the command still fail before mutation with the same exact sorted diagnostic.

- Persists ordinary backtick state through the existing raw active-delimiter path
- Closes on the next literal backtick without cooked-string escape parity and retains same-line close/reopen behavior
- Keeps exact triple-backticks longest-first and prevents comments, recipe bodies, or closed same-line commands from hiding definitions
- Adds Just-parseable acceptance, exact insertion-byte, sorted collision, and byte-preserving rejection fixtures
- Preserves paired-helper identity and every triple-string, ordinary-quote, installer, staged-distribution, and suite-manifest contract

## 0.186.13 — Escaped-Link and List-Paragraph Classification (2026-08-10)

The shipped-reference release guard no longer leaks first-party URLs through escaped closing delimiters or hides rendered four-column continuations inside ordinary list paragraphs. Publication and target-resolution policy remain unchanged.

- Applies odd/even escape parity to closing brackets and destination-opening parentheses in the shared rendered-region mask
- Preserves continuation links for nonempty bullet and one-to-nine-digit ordered list items
- Keeps empty markers, blanks, fences, and genuine indented code outside published-reference discovery
- Adds exact RED/GREEN production-helper fixtures with source-length and newline-offset invariants
- Records remaining relative-parity, label-content, and list-fence-info variants in consent-gated REQ-163

## 0.186.12 — Occurrence-Complete Retired Alias Matching (2026-08-10)

The test-only retired-command guard now evaluates exact source occurrences and overlapping historical install/setup candidates completely, so approved queue-board references cannot hide another retired invocation and invalid longer candidates cannot suppress a valid former command head.

- Applies queue-board branding and test-reference exemptions only to their exact source spans
- Continues eligible install/setup candidate evaluation after a longer historical trigger fails its right boundary
- Recognizes unknown former `install-<target>` routes through the test-only historical prefix without restoring runtime compatibility
- Preserves all 186 inventory identities, 585 direct-boundary negatives, current sibling routes, prime fingerprints, and repaired live surfaces
- Adds focused aggregate RED controls and passes a 34-case independent adversarial acceptance matrix plus full distribution contracts

## 0.186.11 — Just Ordinary-Quote and Triple-Backtick State (2026-08-10)

Safe custom Justfiles no longer look like reserved-recipe collisions merely because ordinary quoted values or triple-backtick commands span physical lines. Real definitions around every form still fail before mutation with the same exact diagnostic.

- Retains raw single-quote and cooked double-quote state with their actual closing and backslash-parity rules
- Recognizes exact triple-backtick command delimiters without letting comments, recipe bodies, or same-line backticks hide definitions
- Adds Just-parseable acceptance fixtures, exact sorted nearby-collision controls, and byte-preserving rejection checks
- Keeps root and shipped helpers byte-identical and preserves triple quotes, managed-span exclusion, reserved derivation, installer ordering, and full contracts
- Records Just's additional ordinary single-backtick multiline form in consent-gated REQ-162

## 0.186.10 — Shared Markdown Rendered-Region Classification (2026-08-10)

Shipped-reference checks now make one length-preserving Markdown visibility decision before any target-discovery path runs, closing false positives and false negatives across the approved indentation, code-span, and escaped-link cases.

- Classifies tab-expanded indented code by effective columns while preserving top-level four-space paragraph continuations
- Keeps links between escaped backticks live and retains ordinary/even-parity exact-run code-span masking
- Hides escaped inline-link and escaped reference-definition targets from both structural extraction and the bare first-party URL fallback
- Adds exact production-helper RED/GREEN fixtures while retaining topology, raw/blob, containment, changelog-identity, and distribution contracts
- Independent review captured the remaining delimiter and list-paragraph variants in consent-gated REQ-161

## 0.186.9 — Complete Retired Core Alias Inventory (2026-08-09)

Retired core command drift can no longer return under an alias merely because it was absent from a hand-picked sample. The distribution guard now reads one test-only historical inventory without republishing compatibility routes or banning ordinary prose.

- Reconstructs 186 concrete trigger rows from the deleted router, shim, and install-normalization contract: 117 direct aliases, 22 install targets across three families, and three bare historical heads
- Validates fixture shape, ownership, counts, install-family symmetry, exact boundaries, and every row on synthetic root and module surfaces
- Preserves current sibling commands, branding/noun phrases, generic pipeline prose, history/fixture exclusions, unique ownership, updater-prime fingerprints, and all REQ-153 live repairs
- Independent review isolated two occurrence-matching edge classes in consent-gated REQ-160; the historical vocabulary itself is complete

## 0.186.8 — Just Triple-String Collision Scanning (2026-08-09)

Valid raw and cooked triple-quoted Just variables can now contain recipe- or alias-shaped payload without blocking suite installation. Real reserved definitions around those values still fail before mutation with the existing exact diagnostic.

- Tracks triple-single and triple-double string state across physical lines, including raw closing, cooked escape parity, same-line close/reopen, and CRLF
- Keeps reserved-name derivation, managed-span exclusion, sorted reporting, and byte-preserving rejection unchanged
- Adds Just-parseable positive fixtures and nearby real-collision controls; adjacent multiline literal forms are captured as consent-gated REQ-159

## 0.186.7 — Exact Manual Stop-Hook Cleanup Path (2026-08-09)

Projects without jq or Python now get an unambiguous nested-object cleanup instruction that cannot be read as deleting custom neighbors. The manual result is checked against both automated reconciliation paths.

- Targets individual `hooks.Stop[*].hooks[*]` objects whose command contains the retired guard path
- Preserves same-wrapper custom hooks and removes an enclosing wrapper only when targeted cleanup leaves it empty
- Covers mixed and guard-only wrappers across jq, forced-Python, and byte-preserving manual fixtures

## 0.186.6 — Shipped Markdown Parser Boundary Hardening (2026-08-09)

Valid code examples, comments, and escaped destinations no longer trip the ordinary shipped-reference release checks. The guard now tests the same parser helpers its live four-module scan uses, so syntax fixes cannot pass only in a fixture-only path.

- Preserves source offsets while masking fenced/indented code, HTML comments, and exact-delimiter inline code spans
- Applies odd/even backslash parity to link and reference structure and normalizes standards-valid escaped destination punctuation
- Adds exact-target positive/negative fixtures while retaining source/install topology, first-party URL, path-containment, and changelog-identity checks

## 0.186.5 — Modular Command Restatement Sweep (2026-08-08)

Live board, knowledge, and toolbox guidance now names the skill that actually owns each command. Transition-era updater lessons have also been reduced to the permanent manifest, validation, marker, and recovery contracts.

- Replaces retired core command hints across managed recipes, board UI/source, knowledge hooks/guides, toolbox guides, and core diagnostics
- Aligns root and installed recipe guidance with the full-suite installer and keeps repaired lesson URLs intact
- Adds a narrow live-surface regression for sibling-owned commands and stale updater-transition fingerprints while preserving history and branding

## 0.186.4 — Tool-Independent Just Recipe Collision Guard (2026-08-08)

Installing the suite can no longer duplicate its managed Just recipes merely because the Just executable is absent. Reserved external recipes and aliases are rejected by name before confirmation or any client write.

- Derives the protected namespace from the shipped managed section instead of maintaining a second recipe list
- Excludes the replaceable managed span and accepts comments, variables, attributes, dependencies, bodies, and longer custom names
- Proves no-Just rejection preserves exact Just/settings bytes, modes, modules, and Git state before mutation

## 0.186.3 — Manual Pipeline-Guard Retirement Guidance (2026-08-08)

Projects without jq or Python now receive a targeted manual settings step instead of advice that would preserve the deleted pipeline guard. Custom Stop hooks and unrelated settings remain explicitly protected.

- Names the retired command substring and limits removal to nested Stop-hook content
- Preserves other hooks in the same Stop event and merges the current core hook fragment afterward
- Adds a mixed custom/retired settings fixture that proves manual mode never mutates the file

## 0.186.2 — Shipped Package Reference Guard (2026-08-08)

Links published by the four installed skills now resolve from both the source archive and the manifest-mapped client layout. Repository-only history stays browseable through canonical GitHub links, and installed changelog history stays current as a verified mirror.

- Points version discovery at the live modular action and repairs five updater-prime lesson links
- Replaces absent installed sidecar paths with tracked repository history URLs and synchronizes the installed changelog
- Adds a local, manifest-derived reference guard for relative source/install paths and first-party raw/blob targets

## 0.186.1 — Permanent Sibling Routing Contract (2026-08-08)

The expired moved-command follow-up now protects the permanent four-skill boundary instead of restoring retired aliases. Every extension action has one sibling owner/route, and the managed update recipe is validated as a core-only call.

- Table-drives all 23 active board, knowledge, and toolbox actions across exact ownership and sibling routing
- Guards each sibling router's pass-through/help behavior while keeping every extension action out of core
- Parses `run-do-work-update` as an isolated recipe and rejects board recipes or binaries in its body

## 0.186.0 — Modular Migration Shims Removed (2026-08-08)

The one-release compatibility window is closed. Core now owns only core commands, and current modular updates use one validated archive through the installed all-or-recover installer instead of carrying bridge, monolith, or legacy-configuration branches.

- Deletes the moved-command shim/routes and removes capability probing, root/monolith update handling, exact legacy Just migration, and old core memory-hook rewriting
- Preserves marker-managed configuration, current hooks, dirty/confirmation guards, installed-helper trust, byte verification, and exact recovery with focused RED/GREEN and full regression coverage
- Queues targeted follow-ups for tool-independent reserved-recipe collision rejection and the remaining retired-command restatement sweep found during independent review

## 0.185.0 — Stateful Pipeline Removed (2026-08-08)

End-to-end work no longer has a second state machine competing with the queue orchestrator. A copyable prompt now composes capture, verification, `do-work run`, and toolbox presentation while the tested work lifecycle stays in one place.

- Removes the `pipeline`/`full` routes, pipeline action/reference, `pipeline.json` reporting, Stop guard, and stale completion-report guidance
- Publishes the approved full-cycle prompt byte-for-byte in core help and the README, with testing and review explicitly owned by `do-work run`
- Migrates the retired guard safely through jq and Python settings paths while preserving custom Stop hooks; the no-JSON-tool manual fallback remains queued for correction

## 0.184.0 — Live Four-Skill Distribution (2026-08-08)

Existing bridge clients can now leave the monolith behind through the same recoverable transaction used by fresh installs. The live archive has one shared version and four focused sibling skills, while project data and unrelated configuration stay outside the managed plan.

- Makes core, board, knowledge, and toolbox the only shipped runtime sources and retains only the three root bootstrap utilities
- Publishes the tested one-archive bootstrap and makes direct plus managed-Just updates converge on validate-first, verify-complete, exact-recovery behavior
- Migrates owned Just recipes and enabled legacy memory hooks, preserves the feature-rich work pipeline for its scheduled follow-up, and starts the one-release moved-command compatibility window

## 0.183.23 — Empty-Quarantine Association Fix (2026-08-08)

Commit and unscoped inspect runs no longer lose every safe candidate when their run-level secret quarantine starts empty.

- Distinguishes the quarantine input by filename instead of the empty-file-unsafe `NR == FNR` idiom
- Preserves all safe M/A/D/XD candidates while retaining current and previously quarantined X paths
- Covers both bridge and staged modular action copies with empty and populated quarantine regressions

## 0.183.22 — Atomic Request Number Reservations (2026-08-08)

Concurrent captures now receive distinct REQ ids before either request file exists, without a global lock or recyclable stale claims.

- Reserves each returned number with an exclusively created marker under `do-work/.req-reservations/` and advances past both request records and prior markers
- Keeps abandoned markers as safe gaps, rejects path escapes and symlinked reservation stores, and roots marker writes inside the repository
- Makes capture stage the marker and proves sequential plus cross-process allocation, existing markers, unsafe paths, and the unchanged decimal output contract

## 0.183.21 — Cross-Platform Atomic Queue Writes (2026-08-08)

Queue-kanban now keeps its complete-file replacement guarantee on both Unix and Windows and refuses symlinked write targets without changing them.

- Uses same-directory rename on Unix and Windows `ReplaceFileW` instead of overstating the cross-platform guarantee of Go's generic `os.Rename`
- Rejects symlinks, special files, missing targets, and targets whose identity changes before replacement without leaving temporary artifacts
- Preserves modes, synced temporary writes, version read-back verification, normal bump behavior, and the Testing-view write guards

## 0.183.20 — Configuration-Aware Bridge Updater (2026-08-08)

Existing bridge clients can now complete the later modular cutover through the trusted updater they already have, including managed configuration migration.

- Keeps the installed bridge validator and full-suite installer as the trusted transaction engine; downloaded archive executables cannot redefine that boundary
- Reviews all four modules plus the owned Just section and known hook migrations, asks once, verifies installed bytes, and restores managed originals after failure
- Carries the full-suite installer inside the staged core for future updates while leaving the live monolithic export guards in place for this final bridge rollout

## 0.183.19 — Full-Suite Installer and Reconciler (2026-08-08)

The staged four-skill suite now has one copy-paste bootstrap and a recoverable client-configuration transaction, ready for live publication at cutover.

- Downloads one archive, validates the shared version and exact four-module manifest, asks once, verifies installed bytes, and restores exact managed originals on failure or interruption
- Reconciles complete, managed, custom, and legacy Justfiles through the board template and managed-section utility without changing exterior client bytes
- Composes core hooks with jq or Python, migrates only known legacy memory command paths, keeps fresh memory capture disabled, and leaves settings unchanged with an exact manual step when no JSON tool exists

## 0.183.18 — Staged Modular Toolbox Skill (2026-08-08)

Optional reviews, reporting, discovery, repository utilities, and companion installers now have an independently loadable staged package.

- Preserves all sixteen retained toolbox actions, their references, guides, and guardrails under `skills/do-work-toolbox`
- Gives toolbox a dedicated router and canonical command surface while resolving queue, board, and knowledge dependencies through explicit sibling paths
- Limits companion installation to ui-design, bowser, last30days, and ideation-adhd, with exact-route and suite-wide runtime-reference coverage

## 0.183.17 — Staged Modular Knowledge Skill (2026-08-08)

Knowledge retention now has an independently loadable staged package, while automatic memory capture remains off until the user explicitly enables it.

- Groups BKB, memory, dream, interviews, prompts, references, guides, assets, and privacy guardrails under `skills/do-work-knowledge`
- Gives explicit memory setup sole ownership of scaffolding, machine-local raw-store protection, hook composition, verification, and rollback
- Moves optional hook commands to deterministic knowledge paths and pins exact legacy migrations without adding memory hooks to core defaults

## 0.183.16 — Staged Modular Board Skill (2026-08-08)

Queue visualization now has an independently loadable staged package, including its full compiled server/CLI, embedded UI, and safe project recipes.

- Preserves live, static, summary, terminal, calendar, Testing, and done-window behavior in `skills/do-work-board`
- Moves the exact managed Just section to a board-owned template with bounded listener shutdown, foreign-process refusal, browser opening, and core updater paths
- Validates complete source inventory, router/runtime references, install paths, recipe safety, and the modular core version-file seam

## 0.183.15 — Staged Modular Core Skill (2026-08-08)

The feature-rich request lifecycle now has an independently loadable staged core package, while the current all-in-one distribution stays active for bridge clients.

- Curates capture, work, verify, review, queue maintenance, schemas, guardrails, specs, hooks, checks, and the suite updater under `skills/do-work`
- Gives core a reduced router and help menu with explicit board, knowledge, and toolbox sibling boundaries
- Adds package contracts for required contents, route/hook/reference resolution, manifest-declared siblings, maintainer-citation exclusions, and pre-cutover root activation

## 0.183.14 — Managed Justfile Recipe Section (2026-08-07)

Do-work now owns one explicit recipe section and can reconcile it without reformatting or overwriting the client's surrounding Justfile.

- Adds a byte-preserving, mode-preserving, same-directory atomic section replacement utility
- Creates fresh templates, appends to custom Justfiles, and migrates the exact legacy five-recipe block without duplication
- Rejects malformed markers and ambiguous interleaved legacy content unchanged, with idempotence and Just parse coverage

## 0.183.13 — Suite-Aware Bridge Updater (2026-08-07)

Existing clients can now receive the modular four-skill distribution safely after the rollout gate, while the current release continues shipping the monolithic layout.

- Makes agent-driven and Justfile updates use one capability-reporting engine for legacy and future suite archives
- Validates every module with the installed bridge validator, shows one reviewed diff and confirmation, and verifies the installed bytes
- Automatically restores only managed Git paths after failure while preserving application files, queue/KB data, Justfiles, and settings

## 0.183.12 — Four-Skill Suite Contract (2026-08-07)

The modular split now has one version, one exact manifest, and one validation boundary, while existing clients continue receiving the unchanged monolithic archive.

- Defines the four source-to-client mappings and all-or-recover update semantics in ADR-019
- Adds a reusable fail-closed manifest validator with malformed-path, duplicate, completeness, symlink, and missing-skill coverage
- Keeps `VERSION`, `suite/`, and `skills/` export-ignored until the bridge rollout gate is satisfied

## 0.183.11 — Warning-Clean Listener Fixture (2026-08-07)

The foreign-listener safety fixture is warning-clean again, keeping later modularization lint failures attributable to the work that introduced them.

- Makes ShellCheck-visible no-op reads of the output captured only for fixture execution and the port consumed indirectly through `eval`
- Preserves the listener refusal, foreign-process protection, and kill boundary unchanged

## 0.183.10 — Testing Date-Filter Empty State (2026-08-07)

Testing columns now say “No matches” when their date window hides existing cards, while ordinary Board empty states stay independent of that Testing-only filter.

- Passes Testing-visible filter state explicitly to the shared empty-copy helper
- Keeps Board, Calendar, and By-UR on the ordinary search/domain/status fallback
- Adds caller-level regression coverage for Board, hidden Board, and Testing columns

## 0.183.9 — Crash-Safe Version Allocation (2026-08-07)

`queue-kanban next-version` now replaces the version file atomically, so an interrupted write cannot leave the release marker truncated.

- Reuses the existing synced temporary-file and rename path
- Preserves the version file's permission bits and existing read-back verification
- Adds regression coverage for atomic replacement, normal bumps, and temporary-file cleanup

## 0.183.8 — Listener-Specific Kanban Shutdown Guard (2026-08-07)

The standing Kanban recipe now refuses to build or start a replacement server while any listener remains on the requested port after shutdown.

- Waits up to 320 one-tenth-second iterations for the old PID to stop listening on that port
- Re-queries the port and names the remaining PID and full command before refusing startup
- Preserves foreign-process protection and keeps the root and installer recipes regression-synchronized

## 0.183.7 — Metadata-Only Commit Hash Recovery (2026-08-07)

Commit-hash idempotency now distinguishes a stranded frontmatter-only metadata edit from unrelated archived-request changes before printing a staging instruction.

- Normalized HEAD and worktree content may differ only by the first frontmatter block's column-zero `commit:` field
- Body prose and fenced-example `commit:` lines remain significant and are rejected when changed
- Temp-file and producer failures remain visible under Bash 3.2 with `pipefail`

## 0.183.6 — Parseable Request Recovery Chains (2026-08-07)

Recovery now distinguishes parseable historical requests from merely non-empty blobs, so malformed metadata damage cannot be selected and reported as a successful repair.

- Historical source selection skips malformed and empty blobs until it finds valid frontmatter
- Recorded-hash traversal treats malformed blobs as part of the damage chain
- Restore validates recovered content before replacement, with a clean-rescan regression

## 0.183.5 — Copy-Aware Secret Inventory (2026-08-07)

Git copy records are now requested explicitly, preventing a destination copied from a secret-shaped source from degrading into a readable ordinary addition when repository rename detection is disabled.

- Secret-derived copy destinations are classified as `X`
- Secret renames remain `XD` plus `X`, while ordinary renames remain `M`
- Commit and inspect manual fallbacks use the same copy-aware porcelain command

## 0.183.4 — Secret Rename Quarantine Survives Re-inventory (2026-08-07)

Resetting a staged secret rename can no longer make its ordinary-looking destination readable on the next inventory pass. Deletion-only commits also accept the safe state Git already staged instead of failing while trying to stage it again.

- Inventory now forces rename detection and quarantines every ambiguous addition when a secret-shaped deletion has lost its rename provenance
- Commit and inspect retain every excluded path for the full action run, including their manual fallback paths
- Secret deletions verify cached name/status only, skipping `git add -u` when the exact deletion is already staged

## 0.183.3 — Residual Shell-Logic Candidates Have an Accurate Disposition (2026-08-07)

The census's remaining candidates no longer imply that every extraction is still unapproved. The record now says what happened without bringing its obsolete line-number table back to life.

- Candidate B is recorded as the separately approved and shipped REQ-121 work
- Candidates A and C remain explicit future decisions rather than accidental work from a queue run


## 0.183.2 — The Effort Chip Shows What Was Declared, Not Just What It Resolved To (2026-08-06)

0.180.1 taught the board's domain and route badges to carry raw provenance — what the REQ
actually declared, flagged when it isn't recognized. The effort chip that landed in parallel
now does the same, so a typo'd value can't be the one field that leaves no mark on the card.

- `effort_estimate` carries `originalEffortEstimate` and its unrecognized flag through to the board, matching how domain and route report themselves
- An unrecognized value now chips as `normal` with an `invalid` flag instead of rendering nothing — the resolved value alone would have hidden the typo everywhere except the warnings banner
- The drawer's Effort estimate row uses the shared declared-vs-normalized renderer

## 0.183.1 — Sweep Lookup, Gate Audit Trail, and Append Staging Hardened (2026-08-06)

Codex's review of the 0.181–0.183 set found four gaps between the new contracts and the
machinery around them — all closed before first real use.

- Sweep REQs carry a `sweep_key` root-cause slug, and the lookup matches on root cause (key first, then the What-section rule statement) instead of "append when one exists" — two unrelated sweeps under one UR can no longer swallow each other's instances
- The lookup also accepts `pending-answers` sweeps, so a generation-≥2 sweep awaiting clarify can't be duplicated by the next review
- The archived `## Review` template now persists each Important finding's `gate:` token and destination — the durable audit record the gate mandates, not just counts
- Both commit recipes (standalone review and pipeline Step 9) stage existing sweep REQs that were appended to, closing the path where new `## Instances` lines silently stayed uncommitted

## 0.183.0 — Same-Root-Cause Review Findings Consolidate Into One Sweep REQ (2026-08-06)

Fifteen facets of one root cause used to mean fifteen queued REQs to wade through. Findings
that share a cause now land in a single sweep REQ with a checklist of instances — approve one
decision, fix one class.

- Reviews route trivial and same-root-cause findings into a `sweep: true` REQ named for the root cause, appending to the existing pending sweep for that UR when one exists (found by grep, never by judging titles)
- Solving a sweep means the class cannot recur — the rule changes everywhere it applies, not N spots patched one at a time
- Only a standalone user-visible finding still earns its own REQ, stating in one line why it couldn't fold into a sweep (`actions/review-work.md` Step 10)

## 0.182.0 — Review Cascades Stop at Depth Two Behind a Consent Gate (2026-08-06)

Reviews of review-spawned REQs could mint fresh auto-worked REQs forever — the UR-489 chain
ran sixteen deep before a human noticed. The cascade now converges at depth two, with the
user as the only escalation path and nothing lost along the way.

- A review of a `review_generated: true` REQ creates its non-critical follow-ups as `pending-answers` — on the board with their effort chip, approved via `do-work clarify`, never auto-worked (`actions/review-work.md` Step 10 → Generation ≥ 2)
- Critical-grade findings (security, data loss, broken production paths) still auto-queue at any depth, prominently reported
- Failure-classification follow-ups and sweep appends are exempt — the stop governs REQ creation from review findings only

## 0.181.0 — Review Follow-Ups Get a Disposition Gate and a Trivial-Work Chip (2026-08-06)

One user request's review chain minted sixteen follow-up REQs over two days — fifteen of them
trivial facets of a single root cause — and nothing on the board said so until the user dug
through them by hand. Every automatic follow-up now declares its weight before it lands.

- Reviews record a `gate:` token (user-visible / rule-change / trivial) on each Important finding before any follow-up REQ is created — severity judgment is unchanged, the gate only routes (`actions/review-work.md` Step 10)
- New `effort_estimate: trivial | normal` frontmatter field, stamped from the gate on review and Discovered-Tasks follow-ups; capture may set it; absent reads as `normal`, so existing REQs need no migration
- The board chips `effort_estimate: trivial` cards and adds a drawer row — display-only, with domain-style normalize-and-warn handling (`tools/queue-kanban`)

## 0.180.1 — Hand-Back, Frontmatter, and Board-State Correctness (2026-08-06)

An external-feedback pass found real edge cases in the newest Git workflow, shell checks,
frontmatter reader, and By UR lens. The fixes fail closed, preserve raw schema evidence,
and give the board the right escape or refresh instead of a plausible stale answer.

- Hand-back bookkeeping now stages and rechecks an explicit path set before a plain commit;
  missing shell arguments and failed `git status` reads return diagnosed usage/tool errors.
- `frontmatter get` rejects empty or non-status membership gates and emits YAML sequences one
  item per line, with empty lists producing empty output instead of Go debug formatting.
- The board carries raw domain/route provenance, distinguishes scope-hidden search matches,
  keeps Testing-only filters out of By UR decisions, and invalidates By UR after test updates.
- Regression probes now isolate global Git ignores, clean up on setup failure, compare exact
  output, derive CLI subcommands structurally, and execute the critical JavaScript state cases.

## 0.180.0 — The Board's By UR Lens Stops Going Blank on a Shipped Queue (2026-08-06)

The moment a queue was fully shipped, the board's **By UR** lens with URs set to Active
showed nothing at all — while the Columns lens, same board, same moment, happily showed
the recently-done cards for those very requests. Since "everything shipped" is the normal
state after a run finishes, the lens was unusable exactly when you'd reach for it to see
what a session touched. Reported from a live 390-UR tree.

- **Active now means open work _or_ a REQ completed inside the RECENTLY DONE window.** It
  used to mean "holds a non-terminal REQ", which is unsatisfiable once everything ships.
  A UR that qualifies via the window shows all its REQ cards, not just the recent ones.
- **The RECENTLY DONE chips drive the By UR lens too.** They were a dead knob there —
  visible, repainting hidden columns, changing nothing on screen. Switching the window
  from one lens now updates the other instead of leaving it stale. Both lenses read the
  same window, so they can no longer disagree about what "recent" means.
- **The hidden-UR count finally shows up when every UR is hidden** — the one case it
  exists for, and the one case an early `return` used to skip. It stays quiet when a
  search, not the scope, emptied the list, where "switch to All" would be a false lead.
- **The empty-state copy tells you which escape to take** — widen the window or drop the
  scope — and names the window you actually selected instead of a hardcoded span.

## 0.179.0 — Naming Conventions Are Now a Loaded Guardrail (2026-08-06)

Every build now carries a naming rule you don't have to paste: no cryptic or single-word names for anything with reach, and names have to be findable with plain-text search. It was a per-project preference people re-typed into each `CLAUDE.md`; now it ships.

- `crew-members/coding-guardrails.md` gains § 5 Naming for Reach — the canonical home. Scopes the two-word minimum to names with *reach* (exported identifiers, struct fields, files, DB tables and columns, CLI flags, env vars, CSS classes), keeps a per-language form clause, and carves out idiomatic short locals behind a two-part test (conventional vocabulary **and** declaration-to-last-use fits one screen).
- Exempts single-word-by-design invocations — CLI subcommands and Make/just targets are already two words where you type them, and a conventions review must not "fix" them.
- States its precedence against § 3 Surgical Changes: the rule governs names you *introduce*, not names already in the host codebase. No drive-by renames.
- `actions/review-work.md`'s Principle Check gains a fifth (informational) item so review spot-checks new names, and explicitly excludes short locals and untouched pre-existing names from being findings.
- `docs/standing-preferences.md` maps the "no cryptic names / must be greppable" nudge to the new default, and notes that language-specific vocabulary still belongs in your own project's agent instructions.
- Every place that enumerated the guardrail set is updated and marked illustrative rather than closed — `actions/work.md` Step 6's loader gloss, the loader summary below it (now un-enumerated), the three spec templates, and the README. Per the Closed-Enumerations-Go-Stale rule, the guardrail file is now stated as authoritative so the next addition can't strand a caller list.

## 0.178.1 — The Secret Scan Catches `.envrc`, and the Scripts Run on Bash 3 (2026-08-06)

Two fixes from an external review of the scripts 0.178.0 shipped, both real.

- The `.env*` exclusion didn't actually cover `.env*`. `.env|.env.*` matches neither `.envrc` nor `.environment` — a suffix with no dot fails both branches — so a direnv file full of exported secrets was tagged as an ordinary new file and could be read and staged. Now `.env*|*.env`, with both callers advertising both forms
- `associate-files.sh` used `mapfile` and `declare -A`, which stock macOS bash 3.2 doesn't have. It died instead of degrading, and the documented fallback only covered a *missing* script. Ownership resolution moved into awk; the fallback now reads "missing or will not run"
- The new regression guard is a behavior probe, not a grep — the bug was a glob that looked right, so the test asserts the emitted tags instead

## 0.178.0 — The Secret-File Scan Is a Script, Not a Paragraph Twice (2026-08-06)

`commit` and `inspect` each carried a word-for-word copy of the same uncommitted-changes scan, including the paragraph explaining why the `-uall` flag is load-bearing. Both are now one shipped script, and so is the REQ-association pass they also duplicated.

- New `tools/checks/uncommitted-inventory.sh` — tags every uncommitted path M/A/D, and tags secret-shaped names `X` so they get reported without being read. Without `-uall`, a brand-new directory collapses to one `?? dir/` row and every file inside escapes the secret scan; that path is now closed by a tested script instead of a paragraph each caller had to get right
- New `tools/checks/associate-files.sh` — matches paths against archived and in-flight REQs, honoring the `status` aliases from the Schema Read Contract. `commit`'s copy knew a `status: done` REQ must still associate; `inspect`'s copy didn't. That's the drift two copies produce, and it's fixed in one place now
- Both callers keep their manual procedure as a documented fallback, and both scripts are pinned in the regression suite so a prose pointer at a missing script fails the build
- Candidate B of REQ-114, split out and shipped as REQ-121. Candidates A and C stay queued, with their grep counts refreshed — the census figures had gone stale exactly as that REQ predicted

## 0.177.0 — `cms` Is a Recognized Domain (2026-08-06)

Content-management work had nowhere to live in the domain vocabulary, so `domain: cms` normalized to `general` and warned on the board every time. It's now a first-class domain with its own crew rules.

- `cms` joins the `domain` enum in the Schema Read Contract, the board's normalizer table, capture's closed set, and both frontmatter samples — one addition, six definition sites, all in step
- `content-management` / `content_management` resolve to `cms`; `CMS` already worked, since the normalizer case-folds before the alias lookup
- New `crew-members/cms.md`, loaded just-in-time like every other domain file. Its opinions are the ones that don't fall out of `backend.md` or `frontend.md`: a content-model change is a data migration, content the CMS owns isn't yours to hand-edit, draft is a state rather than a copy, and editor-facing strings are a deliverable

## 0.176.6 — Shipped Comments Stop Pointing at a File Consumers Don't Have (2026-08-06)

Four comments in the shipped board tool cited this repo's maintainer doc, which is export-ignored — so the pointer dangles in every consumer install. The repo already has a check for exactly this; it had been failing, partly since 0.175.0, hidden in the same output as seven unrelated runner-dependent failures.

- Each rule is restated instead of dropped: two now point at `actions/board.md` and `actions/work-reference.md`, which ship; two state the obligation directly, having no shipped home to cite
- The board prime's route lesson keeps its content and gains the actionable half — a field joins the display-parsed list the moment the board starts parsing it
- The check's per-file allowlist was deliberately not widened: it exists for mentions of a *consumer's* CLAUDE.md, so using it here would have silenced the probe rather than satisfied it
- No behaviour change; the contract-regression suite is back to its true baseline

## 0.176.5 — An Off-Vocabulary Route Warns Like a Bad Domain Does (2026-08-06)

0.176.2 taught the board to normalize `route` and 0.176.3 gave `domain` a warning when its value is off-vocabulary. `route` got the first and not the second, so `route: z` displayed as `Z` with no footprint — the exact silence 0.176.3 removed, one field over. An external review caught the asymmetry.

- `route: z` still shows as `Z` and now raises the contract's warning; blanking it was never an option, since route has no documented default and an absent route would then look identical
- The two contract fields the board reads share one collector now, so the warning's wording still lives in exactly one function
- A lowercase `a` and an absent route both stay silent, each with its own test

## 0.176.4 — Reading a Timestamp With --normalize Stops Warning About It (2026-08-06)

`frontmatter get … created_at --normalize` printed a warning calling the timestamp "not recognized" on every single call — and the contract it cited says the opposite: a field with no canonical vocabulary is *outside* the contract and read verbatim. Now it's a clean no-op for those fields.

- A field with no contract row prints its value and nothing else, observably identical to the same `get` without the flag
- The gate is a lookup of the contract table, not a list of exempt field names — the exempt set is "whatever has no row", so a hand-written list would go stale the first time a row is added
- `--in-set` is deliberately not silent there: both set names are `status` sets, so a membership test on a timestamp is now a usage error instead of a "no" that reads as a real answer at a call site
- A typo'd `status` still warns exactly as before — 0.175.1's fix has a regression guard here so this couldn't quietly re-open it

## 0.176.3 — A Typo'd Domain Leaves a Footprint on the Board Again (2026-08-06)

`domain: quantum` rendered as a plain `general` badge with no warning anywhere, so a misspelled domain was *harder* to spot after 0.174.15 than before it — the value at least used to reach the card verbatim. The board now says so, in the same warnings banner it already uses for a typo'd testing status.

- The unrecognized flag is kept rather than discarded; the card still shows the contract's `general` default, and the board raises the contract's own warning line naming what was actually written
- The code comment defending the silence claimed the board had no warning channel for domain. It has one, the sibling field has used it since it shipped, and the frontend already renders it — the comment now says that instead
- A documented alias (`back-end`) and an absent domain stay silent, each with its own test: a channel that fires on ordinary REQs is a channel readers learn to ignore
- No frontend change was needed — anything appended to the board's warnings list already prints in `do-work board summary` and renders in the data-warnings banner

## 0.176.2 — The Board Now Normalizes Route, Which 0.174.15 Claimed and Didn't (2026-08-06)

An external review caught 0.174.15 overclaiming: its title said the board honored the Schema Read Contract for all nine fields, but only `domain` was ever wired. `route` kept reading verbatim, so a REQ written `route: a` showed a lowercase `a` badge — in the one other contract field the board actually displays. Now wired, with the correction on the record.

- `route` normalizes at the board's read site, through the same table REQ-111 added
- Deliberately `normalizeSchemaField`, not `resolveSchemaField`: route has no documented default, so resolving would blank an unrecognized letter and hide the REQ that needs re-triage
- **0.174.15's claim was too broad and stays on the record as written.** Only `domain` was wired then; five of the seven fields it named (`caveman`, `maintenance`, `tdd`, `error_type`, `kb_status`) the board still doesn't read at all, which is correct — they have no display role, and adding one to make an old title true would be backwards
- The maintainer doc's list of fields the board parses for display had never included `route` — which is why the field carried no keep-in-sync obligation and drifted in the first place

## 0.176.1 — The Just-Kanban Install Verifies All Five Recipes (2026-08-06)

`just-kanban`'s verify step checked two of its recipes and reported success for all of them. An absent recipe is not a syntax error, so `just --list` parsed happily over a partial append and the install claimed to have provided commands that weren't there.

- Phase 3 now greps every recipe header and names the missing ones
- Each pattern requires a space or colon after the name, so `run-kanban` can no longer be satisfied by `run-kanban-cli`

## 0.176.0 — Terminal Digest of What's In Flight (`just run-kanban-cli`) (2026-08-06)

Opening a board tab is more than "what am I working on?" deserves. `just run-kanban-cli` answers it in one screen: how many REQs are open, every claimed REQ with its title, and every needs-input/blocked REQ with the status that parked it there.

- New `queue-kanban open-work` subcommand (read-only) behind the recipe, also reachable as `do-work board cli`
- Blocked lines carry the `blocked_by` condition; an off-vocabulary status shows verbatim as `invalid:<value>` instead of hiding
- Finished work never appears — this is open work only; parse warnings show as a count pointing at `summary`
- `do-work install just-kanban` now installs five recipes, and the unknown-subcommand error is checked against the dispatch switch itself rather than a hand-written list

## 0.175.2 — The Commit Action Reads a REQ's Status Through the Tool (2026-08-05)

The new `frontmatter` command had no callers, so the hand-rolled reads it exists to replace were all still hand-rolled. `do-work commit` is the first one switched over, and it's the site where getting it wrong actually bit: testing for the literal `completed` drops every remediated-with-issues REQ, so its files never get associated.

- Step 3's terminal-success check now prefers `frontmatter get … --in-set terminal-success`, which normalizes aliases before testing
- The `awk` floor is spelled out as a working command, and building the tool to get the value is explicitly prohibited
- One site only — the other read sites stay candidates under REQ-114

## 0.175.1 — A Typo'd Status No Longer Passes Silently (2026-08-05)

`frontmatter get … status --normalize` printed a misspelled status straight through with no complaint, because the two fields whose aliases live in their own normalizers were being force-marked as recognized. That was the exact no-feedback path the command was added to replace.

- `status` and `testing_status` now warn like the other seven contract fields; their alias maps stay in one place and are unchanged
- A field the contract gives no default now says so, instead of claiming `Treating as ''`
- Aliases still resolve silently — `status: done` prints `completed` with nothing on stderr

## 0.175.0 — Read a REQ Field From the Command Line (2026-08-05)

The shipped frontmatter parser had no way to call it — none of the tool's subcommands took a file and a field — so every action that needed a REQ's `status` or `domain` hand-rolled its own `awk`. There's now a command for it, which means one tested implementation instead of ~95 copies.

- New `queue-kanban frontmatter get <file> <field>`, with `--normalize` to apply the Schema Read Contract and `--in-set terminal-success|terminal-resolved` for the finished-work check
- The value goes to stdout and every diagnostic to stderr, so `value=$(… )` captures cleanly even when a contract warning fires
- Read-only, and an accelerator only: the shell fallback stays documented and nothing builds the tool to read a field

## 0.174.15 — The Board Now Honors the Schema Read Contract for All Nine Fields (2026-08-05)

`domain: back-end` used to reach the board exactly as written, because only two of the contract's nine enum fields had a normalizer anywhere in the repo. All seven of the others now resolve their documented aliases through one table, so a muscle-memory spelling stops silently meaning something else.

- `domain`, `route`, `caveman`, `maintenance`, `tdd`, `error_type`, and `kb_status` now normalize per `actions/work-reference.md`'s Schema Read Contract
- An unrecognized value falls back to the documented default and reports itself unrecognized, rather than being silently remapped
- An absent field stays absent for the board — a domain-less card gains no badge and no filter entry

## 0.174.14 — Version Bumps Now Sync the Lockfile That Mirrors Them (2026-08-05)

Bumping a project's version file left its lockfile's copy of that same version behind, so the next build or install rewrote it and the tree read dirty for a change nobody made — one consumer repo drifted 8 patch versions and collected three ad-hoc "sync the lockfile" commits before anyone traced it. Step 9 now bumps the mirror in the same commit, by hand, so it works in a repo with no toolchain installed at all.

- **New "Lockfile mirror" note** in the Changelog Entry Procedure's source 1. The trigger is stated as a condition — the repo commits a lockfile recording this package's own version — with `package-lock.json`, `Cargo.lock` and `uv.lock` named as known instances rather than the boundary. `pnpm-lock.yaml`, `yarn.lock`, `poetry.lock` and `go.sum` carry no root version and need nothing, and the split follows the lock tool, not the manifest: `pyproject.toml` mirrors under uv and not under poetry.
- **The mirror is hand-edited, not delegated to the package manager.** `npm install --package-lock-only` runs the target repo's `preinstall`/`install`/`postinstall`/`prepare` hooks, restructures an old `lockfileVersion` 1 file wholesale, re-resolves dependencies the lockfile is behind on, and fabricates a `package-lock.json` in a pnpm or yarn repo — all at exit 0. `cargo generate-lockfile` drags unrelated dependencies and checksums forward. The hand edit writes the same bytes with none of that.
- Notes that the drift is cosmetic under npm (`npm ci` tolerates it) but hard-fails `cargo check --locked` and `uv lock --check` — so it isn't always just tidiness.
- **Workspace members mirror somewhere else.** Bumping a member's `package.json` leaves `packages["<member-path>"].version` in the root lockfile stale while both root-package sites correctly stay put — so reading only those two says "nothing to sync" and the drift survives. The table calls this out.
- **The Step 9 staging validation kept working.** Its exemption list only fires when nothing outside it is staged, so it gained the lockfile too — otherwise a fourth staged file class would have silently switched the "this commit contains no implementation" check off entirely.
- Every staging list names the lockfile, worktree dispatch mode included — its list says "stage **only**", so omitting the lockfile there would have contradicted the generic list directly above it.
- The `git add` block, the lifecycle-file list, and three restatements in `actions/work.md` now name the lockfile, including the Go accelerator note, which says plainly that it touches no manifest or lockfile.

## 0.174.13 — Recovered Trap's Evidence Corrected, Probe Rows Name the Right Status Set (2026-08-05)

An adversarial review of last release's recovered knowledge-base entry re-ran its verification steps and found one trap's evidence wrong on every point. Fixed, and the negative result is now recorded too — it is the more useful half.

- **The recovered stale-binary trap claimed a hazard wider than the real one.** It said nothing in the pipeline rebuilds the board binary before `verify` runs; in fact both shipped call sites — `actions/forensics.md` Check 14 and `actions/work.md` Step 9 — build in the same command, and did so when the trap was written. Its supporting evidence was wrong twice over: `go version -m` reports `vcs.time`, the *commit's* timestamp rather than the binary's mtime, and a byte-for-byte diff against a fresh build always differs once HEAD has moved, because Go stamps the revision and dirty flag into every build. The entry now scopes the hazard to a hand-run `verify` outside those blocks, and records both failed staleness tests so nobody re-derives them.
- **`actions/forensics.md`'s Check 14 rows now say "non-terminally-resolved"** instead of "non-terminal". The repo uses both senses — `actions/cleanup.md` calls `failed` terminal, while the Terminal-resolved set deliberately excludes it — so the looser word told readers a `failed` REQ was outside rows that in fact fire on it.

## 0.174.12 — Assignment Probe Stops Double-Reporting Finished Work (2026-08-05)

The double-fire fixed in 0.174.11 was hiding one probe over, so it got the same carve-out. Plus a checklist line that was quietly too narrow, and a note explaining a test suite that only goes red for root.

- **`assigned-elsewhere-claimed-here` now skips terminally-resolved REQs.** A finished REQ stranded in `working/` while still carrying `assigned_to` tripped this probe *and* the stranded-finished one, and this one's remedy told you to clear or release a claim on work that was already done. The stranded-finished probe owns that state alone now — the same carve-out, keyed on the same shared predicate, as the archived-UR fix last release. `actions/forensics.md`'s Check 14 row is qualified to match.
- **`actions/work.md`'s verification checklist covers every claim Step 1 deliberately leaves alone**, not just a reported foreign one. A label-less claim is left intact too, and the old wording quietly excluded it — the last instance of a vocabulary drift that had been surfacing one file at a time.
- `_dev/tests/update-script-behavior.sh` says up front that it needs a non-root runner. Its failure injection drops write permission on a directory, which root ignores, so under root seven probes fail on messages the updater never had reason to print — a property of the runner, not the code.
- Recovered the UR-018 session's still-accurate operational traps into the knowledge base before they were only in git history — exact flag orders, the stale-binary trap that makes `verify` lie, and the contract-suite phrases to reword around. Repo-internal; `kb/` does not ship.

## 0.174.11 — Claim-Sync Timing, Verify Probe Overlap, and Addendum Earmarks (2026-08-05)

An external review of the 0.174.x series turned up six things the docs and the board's verifier were quietly getting wrong — mostly promises a bit broader than the code behind them, plus two probes that both fired on the same state. All small, all now honest.

- **Claims travel when bookkeeping commits, not when the claim happens.** The Execution Model said claims reach other checkouts "by ordinary git sync" without saying that nothing commits a claim until a checkpoint, a hand-back, or the release tail — so the duplicate window was real but undocumented. Crash recovery's detection claim is qualified the same way.
- **`queue-kanban verify` no longer double-reports a stranded finished REQ.** A terminally-resolved member of an archived UR tripped both `ur-archived-with-live-member` and the stranded-finished probe, and the former's remedy told you to run or abandon work that was already done. The stranded-finished probe owns that state now.
- **Capture carries the earmark through both addendum paths.** Step 1's `assigned_to` assessment reached a fresh REQ but not an addendum: appending to a queued REQ touched only the body, and the Addendum REQ Template had no `assigned_to` line to fill in.
- **The exit summary has a headline for "ready, but assigned elsewhere."** When every dependency-ready REQ was earmarked for another session, neither documented headline fit. The reference and `actions/work.md` now say *claimable* where they meant it, and that sixth section finally has a lead that matches it.
- **The `write_set` fan-out bullet no longer contradicts itself** — its bold lead called the field "not an input to either" pick while the next sentence called it advisory input to the manual one.
- **Fixed the archive paths the UR-018 consolidation left behind** in `actions/work-reference.md` and ADR-018: REQ-095, REQ-100, `input.md`, and the approved plan all live under `do-work/archive/UR-018/` now.

## 0.174.10 — Merge-Time Queue Guard, Clean-Index Hand-Back, Synthesized-UR Copy Fix (2026-08-05)

Three accepted findings from an external review of the worktree-dispatch pipeline and the board. The common thread on the first two: the mechanical checks all excluded `do-work/`, so a builder's committed queue edits could integrate with only human review watching.

- Hand-back step 2 now opens with a pre-merge guard: `git diff <pre>...<operative_name> -- do-work/` must print nothing — run before the merge because the three-dot diff (and verify's probe, which reads the same one) goes blind after it (`actions/work-reference.md`)
- Hand-back step 0 now requires a clean index: bookkeeping commits path-limited, anything else staged stops and surfaces to the user — a plain commit was silently sweeping unrelated staged files below `<pre>`, outside every evidence step's range
- The board's Copy button no longer copies an empty string for a synthesized UR card (a REQ pointing at a UR with no `input.md`) — those URs are omitted from the markdown map so the rendered-text fallback fires (`tools/queue-kanban/generate.go`)

## 0.174.9 — Session-Start Note Uses the Canonical Recovery Terms (2026-08-05)

Last stale restatement from the 0.174.7 vocabulary change (or so we thought — review found one more, see below). work.md's Step 10 session-start note now states the own-label condition and defers the recovery case list to Crash Recovery, instead of hand-enumerating the cases and calling a label-less entry a "foreign claim" (canonically a *claim of unknown origin* since 0.174.7).

- `actions/work.md` Step 10 session-start note, item 2: condition stated once, cases deferred to `actions/work-reference.md` → Crash Recovery (Step 1)
- Closes UR-018 (parallel building across checkouts) — REQ-109 was its last live member
- Known remainder: the review's restatement sweep flagged `actions/work.md`'s Verification Checklist ("a reported foreign claim") as the same drift class — minor, not yet fixed

## 0.174.8 — Label-Less Drop Carried Through the Satellite Rules (2026-08-05)

Follow-through on 0.174.7: the rules that orbit crash recovery now agree with it. The In-Progress Record paragraph states the own-label condition instead of restating the case list (one enumeration, at its canonical home), and a label-less checkpoint entry finally has a documented exit — it leaves with the REQ when a human reclaims it, in both the removal rule and forensics Check 1's manual reset.

- `actions/work-reference.md` In-Progress Record: condition stated once, case list deferred to Crash Recovery; removal rule covers the label-less reclaim
- `actions/forensics.md` Check 1: the manual reset sends a label-less entry with the reclaimed REQ
- Decision records: ADR-018 `updated:` bumped, decision log notes the edge closed

## 0.174.7 — Label-Less Checkpoint Entries Are Never Auto-Recovered (2026-08-05)

Crash recovery no longer treats a locally-modified `CHECKPOINT.md` as proof that this checkout wrote a label-less claim entry — under claim-anywhere, every concurrent claim dirties that file with a merge resolution, so the inference could classify another checkout's live claim as an own crash and strip it (REQ-095 F-06/F-07 reproduced exactly that). A label-less entry is now a claim of unknown origin, always report-only.

- `actions/work-reference.md` Crash Recovery: the label-less bullet drops the authorship heuristic; reclaiming a genuinely-own pre-0.170.0 entry stays a human act (takeover ladder or `actions/forensics.md` Check 1)
- Contract suite pins the new rule both ways: report-only classification must stay stated, and the retired "locally modified ⇒ mine" inference must stay gone
- ADR-018's Consequences updated: the known-unsound edge it recorded is now resolved

## 0.174.6 — Assigned-Badge Comment No Longer Over-Claims (2026-08-05)

The board's assigned-badge comment said "nothing here folds, trims, or rewrites the value" three lines above a call that truncates the badge text. Behavior was always fine — the tooltip and drawer carry the full value — the comment just conflated value normalization with display truncation. It now states both halves accurately.

- `tools/queue-kanban/web/board.js` comment-only reword; no executable change, Go tests pass.

## 0.174.5 — Auto-Wave Honors Targeting-Token Provenance (2026-08-05)

Reading only work-reference's Auto-wave list, an explicitly-named but dependency-blocked REQ looked excluded from the wave — while work.md said it runs. The two files now give the same answer.

- Auto-wave condition 2 (dependency-ready) carries the serial scan's provenance carve-out: a named `REQ-NNN` enters the wave regardless of `depends_on`; a UR-expanded member stays gated, scoped to the UR's set.
- A new paragraph states that targeting tokens scope the wave's candidate pool, per-token provenance survives, and there is no separate readiness predicate for waves — the rule work.md already carried.
- Doc-only alignment; no behavior change to the default (untargeted) wave computation.

## 0.174.4 — Capture Now Seeds `assigned_to` Earmarks (2026-08-05)

The schema said "seeded by capture when the user earmarks work" — but the capture action never mentioned the field, so nothing could ever seed it. A downstream sync review caught the gap; capture now carries the instruction.

- Capture Step 1 gains an **Earmark assessment**: "leave this for cloud-alpha" in the request text seeds `assigned_to: "cloud-alpha"` — verbatim, YAML-quoted, never invented when the user names no session.
- The Simple REQ template documents the optional field alongside the other conditional frontmatter.
- Scope correction from triage: the seeding claim only ever lived in the schema line, not the 0.172.0 changelog entry — so the docs and the behavior now agree everywhere.

## 0.174.3 — Verify's UR Probe Works Under a Relative `--repo-root` (2026-08-04)

Two review findings on the 0.174.x work, both real. One was a probe that silently found nothing; the other was a contract sentence claiming something the code has never done.

- `ur-archived-with-live-member` required a leading `/do-work/archive/` in the UR's path, so `do-work verify --repo-root .` recognized zero archived URs and reported nothing — the quiet kind of wrong. It now anchors on the trailing separator and works under both absolute and relative roots, while still refusing a directory merely named `archive`. Every existing test built its fixture under an absolute temp dir, which is why none of them caught it.
- The `assigned_to` schema text said readers must not trim the value. They always have: every field in the verbatim-read class, `write_set` and `prime_files` included, goes through the same scalar coercion. Corrected the contract rather than the parser — padding survives only explicit YAML quoting, means nothing in a name, and preserving it would make `" cloud-alpha "` a different session from `"cloud-alpha"`. Verbatim means no alias map, no case folding, no path canonicalization.
- Three regression tests: the relative-root probe, the lookalike-directory rejections, and case-preserved-padding-trimmed.

## 0.174.2 — Multi-Checkout Guide and the Session-Ownership Decision Record (2026-08-04)

The claim-anywhere model now has documentation and a decision record. The guide gets a walkthrough for working one queue from several checkouts; ADR-018 records why the contract changed, including the parts that turned out to cost something.

- New `docs/work-guide.md` section: claiming from a workspace, clone or cloud session; earmarking with `assigned_to` and the two ways to override it; the one-releaser rule; both conflict shapes and how to resolve them without deleting a real claim; `--fan-out`; and the caveat that all of it presumes a committed `do-work/`.
- ADR-018 is the first decision record this skill has for session ownership at all — the 0.161.0 exclusive-session decision it partially reverses was never written down as one. Six decisions, seven rejected alternatives (including the two that were nearly taken), and the costs the acceptance runs found rather than the ones the plan predicted.
- The ADR indexes were three counts behind. Those hand-maintained numbers are gone, replaced by a pointer to read `records/` — a count in a file nobody recounts is a stale count.

## 0.174.1 — Wall-Clock Concurrency Proven for the First Time (2026-08-04)

Fan-out has been in the contract for a while, and until now nobody had watched two builders actually run at the same time — the one recorded attempt logged Partial. This one measured a 4.109-second window where both were genuinely running.

- Two builders in two git worktrees on two branches, dispatched by the automatic ready-set computation and launched before the owner waited on either.
- The computation was shown *excluding* as well as including: a REQ whose dependency was still pending stayed out, and one earmarked with `assigned_to` was left for its own session.
- Serial integration produced two different `<pre>` values, because the first REQ's archive commit moved the tip between merges. That is why the contract says capture it once per REQ — and writing the once-per-wave version by mistake is what proved the rule earned.
- Wave 2 recomputed and picked up the REQ whose dependency had just landed, which a carried-forward remainder list would have missed.
- The negative case ran too: two REQs declaring the same `write_set` both entered the wave, and the second merge refused with a content conflict. A computed set claims its REQs are runnable, never that they don't collide.
- Honest gap, recorded: the builders were scripts. The mechanism is proven under real concurrency; agent behavior under concurrency still isn't.

## 0.174.0 — `do-work run --fan-out` Computes the Wave Itself (2026-08-04)

The pipeline used to say flatly that nothing computes the set — you picked which REQs ran together, by hand. That was a deliberate design choice, and it's now a deliberate reversal: pass `--fan-out` and the loop works out the ready set and dispatches builders with no confirmation prompt.

- Ready means pending, dependencies met, unclaimed, and not earmarked for another session. Same predicate the serial scan already uses — no second definition of readiness.
- Waves are bounded, never an unbounded fan-out over the whole queue. `--fan-out N` sets it; bare `--fan-out` uses the harness limit, or two when that's unknown.
- Each wave is recomputed after the previous one integrates, so a REQ whose dependency just landed joins the next one.
- `--fan-out` and `--wave N` compose: `--wave` picks which REQs, `--fan-out` picks how many at once.
- `write_set` is still not a scheduling input, and now explicitly is not read by the wave at all. A computed set claims its REQs are runnable — not that they don't overlap. The merge is still what proves that.
- The default is unchanged and still serial, so the simplest agent reads the same instructions it always did. No worktree support means the flag quietly does nothing.

## 0.173.0 — Verify Catches Assignment Drift and Half-Closed URs (2026-08-04)

Two new read-only probes for drift that only became reachable once several checkouts share a queue. Both report and route; neither repairs, and neither is marked fixable, because both resolutions are yours to make.

- `assigned-elsewhere-claimed-here`: a REQ sitting in `working/` while still earmarked for another session. Its marker is now telling every other checkout to skip a REQ you're already building.
- `ur-archived-with-live-member`: a UR archived while one of its members is still in `queue/` or `working/` — the live REQ is orphaned from the `input.md` that explains why it exists.
- The UR probe scans `user_request:` frontmatter, never the UR's own `requests:` array, which is capture-time-only and misses follow-up REQs. There's a test that fails if anyone changes that.
- Both show up in `do-work forensics` Check 14, whose probe table now says out loud that it's a snapshot rather than the authority — the tool's own output is.

## 0.172.0 — Earmark a REQ for Another Session with `assigned_to` (2026-08-04)

Now that any checkout can claim, there needs to be a way to say "leave this one, it's mine." That's `assigned_to` — one optional frontmatter field naming a session. Another session's default run skips it and tells you why; naming it explicitly takes it anyway and clears the marker. No verb, no status, no clock.

- `assigned_to: "cloud-alpha"` is read verbatim — session names have no canonical spelling, so nothing folds case or trims it.
- The default work scan skips an assigned REQ and lists it in the exit summary as *assigned to <name>*. It's a courtesy, not a lock: nothing checks whether that session is actually running.
- `do-work run REQ-NNN` overrides the skip and clears the field as part of the claim. Reaching it through `do-work run UR-NNN` does not — naming a batch is a weaker signal than naming a member.
- The board shows an `assigned` badge and a drawer row. It never moves, hides, or reorders a card on it.
- Nothing else came back with it: no `assigned_at`, no staleness threshold, no auto-release, and the reserve verb and `reserved` status stay dead.

## 0.171.0 — Claim From Any Checkout, Release From One (2026-08-04)

The ownership contract used to say one queue owner per checkout, and cross-session ownership was out of bounds. It now says the opposite about *claiming*: point as many checkouts at a queue as you like — a worktree, a second workspace, a clone, a cloud session — and each may capture, claim and build. What stays single is the release tail.

- `Execution Model — Exclusive Session` is now `Execution Model — Claim Anywhere, One Releaser`. Any checkout claims; exactly one merges, bumps the version, writes the changelog, archives, and closes URs.
- Two releasers against one queue, and two sessions in one working tree, both stay unspecified — nothing prevents them, and `do-work forensics` / `do-work cleanup` are the repair path.
- Cross-checkout collisions are ordinary merge conflicts fixed when branches meet. The duplicate-id probe catches colliding captures; the checkpoint writer label catches double claims.
- A builder tree no longer has to be a spawned worktree — a workspace, clone, or remote sandbox works too. A remote builder's hand-back travels on its branch, and a non-releaser checkout treats its synced queue as a snapshot rather than the truth.
- New guidance for the conflict you *will* hit: `CHECKPOINT.md` collides on every concurrent claim, even disjoint ones. Keep every entry from both sides — dropping either one loses a real claim.
- New red flag: a second checkout running the release tail. Claiming and building elsewhere is now in contract; releasing twice is not.

## 0.170.2 — Two-Clone Acceptance Run Confirms the Writer Label (2026-08-04)

The checkpoint writer label from 0.170.0 now has a real experiment behind it instead of an argument. Two clones sharing a bare origin reproduced the old poisoning as a silent fast-forward that erased a live claim, then the shipped rule left the same claim byte-for-byte untouched. The run also found a hole nobody had reasoned their way to.

- Reproduced the pre-0.170.0 strip deterministically: a routine sync exits 0, fast-forwards, and deletes another checkout's live claim along with its Triage and Scope — no conflict, no warning.
- Confirmed the shipped rule reports `claim held by <writer>, not touched` and writes nothing, verified by matching sha256 before and after.
- Captured what a cross-checkout claim conflict actually looks like: plain content conflicts, never a rename conflict, plus an `add/add` conflict on `CHECKPOINT.md` for *every* concurrent claim — including two that touch nothing in common.
- Found that with byte-identical claim writes the REQ file merges clean, so the writer label is the only thing that surfaces the double claim at all.
- Filed REQ-104: the label-less recovery rule reads "checkpoint is locally modified" as proof this checkout wrote it, which stops being true once resolving a checkpoint conflict is routine.

## 0.170.1 — Step 10 Preserves Label-Less Checkpoint Entries Too (2026-08-04)

Follow-up to 0.170.0's writer label: the work action's session-end rewrite and session-start delete were scoped to entries "carrying another checkout's label", which quietly excluded label-less legacy entries — recovery refused to touch them, then Step 10 deleted them anyway, sending the REQ into the takeover ladder next run. Both clauses now preserve every entry this checkout did not write, and two new contract assertions pin both destruction paths.

## 0.170.0 — Checkpoint Writer Label Scopes Crash Recovery to the Claiming Checkout (2026-08-04)

In-progress checkpoint entries now record which checkout wrote them (`writer: <hostname>:<checkout-path>`), and crash recovery only ever strips its own — another checkout's synced live claim is reported as `claim held by <writer>, not touched` instead of being read as a local crash and re-queued. First piece of the UR-018 multi-checkout work: the old behavior replayed the 2026-07-01 collision deterministically on any tree that commits `do-work/`.

- Crash recovery classifies four claim origins: own-label (recover), foreign-label (report, never in the takeover ladder), label-less legacy (own only if the checkpoint is locally modified), unnamed/no-checkpoint (unchanged 3-hour ladder)
- Step 10's session-end rewrite and the session-start delete preserve entries this checkout didn't write
- The no-liveness tripwire still bans refresh intervals, staleness checks, and liveness probes by name — the static label is the one carve-out, pinned in the same paragraph by a new contract assertion (3 new assertions total)
- `queue-kanban verify`'s ghost-REQ remedy no longer suggests deleting a checkpoint that may hold another checkout's live entries

## 0.169.9 — Verify No Longer Reads a Quoted Shell Glob as a Ghost REQ (2026-08-04)

`queue-kanban verify` exited 1 on this very repo: the checkpoint ghost-REQ scanner's `REQ-\d+`
pattern matched the `REQ-0` prefix of the shell glob `REQ-0[0-9][0-9]-*.md` quoted in
CHECKPOINT.md's session notes, and reported a REQ that never existed.

- The mention scan now skips a digit run that continues straight into `[` — that's a glob
  character class, not an id. Only `[` is treated as a continuation: `*` and `?` legitimately
  follow real ids in prose (`**REQ-093**` emphasis, a sentence ending in a question mark), so
  skipping on those would hide genuine ghosts.
- The scan lives in its own tested helper (`checkpointMentionedRequestIds`), pinned against the
  exact line that produced the false positive.

## 0.169.8 — Update Flow Deletes Orphaned Reserve Files; Verify Keeps Space Paths Whole (2026-08-04)

Two follow-through fixes from an external review of a consumer install. Tar extraction never deletes
what upstream removed, so pre-0.161.0 installs still carried the two reservation-workflow files the
exclusive-session cleanup deleted — and one of them is a `prime-*.md`, which the prime audit's
recursive glob keeps rediscovering as a live doc.

- `do-work update` Step 5 now removes the two orphaned reservation-workflow files alongside the
  stale maintainer docs, and the "leftovers in `actions/` are harmless" note gains the prime-glob
  exception that made it wrong.
- The contract suite's reservation-token sweep grew a scoped exemption: `actions/version.md` may
  name the removed files only on `rm -f` lines — any other mention still fails (probed both ways).
- `queue-kanban verify`'s dirty-worktree probe parsed porcelain output by last whitespace field,
  truncating any path containing spaces in the report line; it now slices at porcelain's fixed
  3-byte prefix and keeps a rename's destination side (unit-tested).

## 0.169.7 — Board Copy Round-Trips CRLF Files Byte-for-Byte (2026-08-04)

The drawer's Copy payload was corrupted for Windows-authored (CRLF) REQ/UR files: the frontmatter
fence was measured by subtracting the parsed body's length from the file's, but the parser had
CRLF-normalized that body, so the fence kept stolen bytes from the body's start and the body pasted
back with Unix endings. Ironic, given the exact-suffix approach exists precisely to preserve CRLF.

- `splitFrontmatter` now scans the raw bytes and reports the true body-start offset; callers slice
  the fence at that offset instead of inferring it by arithmetic.
- `BodyMarkdown` is now genuinely verbatim (its comment always claimed it was) — CRLF endings
  survive into the Copy payload, and `FrontmatterMarkdown + BodyMarkdown` equals the file on disk
  byte-for-byte. YAML parsing still sees normalized text.
- New tests pin the round-trip invariant at both the splitter and the ticket-parse level.

## 0.169.6 — Every Historical Changelog Is Back on Disk, Named by Range (2026-08-04)

192 entries of release history hadn't been in the working tree since 0.76.0 — readable only via
`git show bf15fe2:…` or a link to a four-month-old commit. All of it is back, and every archive now
says in its own filename which releases it holds.

- The two deleted archives are restored verbatim as `CHANGELOG-2026-04-07-up-to-v0.49.0.md` and
  `CHANGELOG-2026-04-13-up-to-v0.64.1.md`; `CHANGELOG-archive.md` is renamed to
  `CHANGELOG-2026-07-07-up-to-v0.109.0.md`. The old name said when it was cut, not what was in it.
- A new slice, `CHANGELOG-2026-07-13-up-to-v0.121.0.md`, takes 0.110.0–0.121.0 out of the live file.
  All 478 prior entries are now on disk in a tracked file, so the `bf15fe2` recovery pointer is gone
  from all three places that carried it.
- The live header lists ranges, not counts — it used to say "the entry below the 20th here," true
  when 0.123.1 wrote it and off by 122 by the time anyone read it again.
- `.gitattributes` export-ignores the glob `/CHANGELOG-20*.md` rather than one hardcoded name, so the
  next cut can't silently start shipping 20k words.

## 0.169.5 — Current-REQ Relevance Covers Dirty Trees and Live Sibling Sessions (2026-08-04)

Two things a run can trip over that aren't its problem: files that were already dirty when it
started, and another session working in the same checkout. Both now route through one rule, and
the second one says the quiet part out loud — don't go looking, and don't make the user adjudicate it.

- **Pre-flight's dirty-file warning says what's actually at risk.** It claimed unrelated files
  "may be swept into the commit," which contradicts Step 9's explicit-pathspec staging and made the
  check look removable. The real exposure is that qualification and review read the repository-wide
  working/staged diff, so pre-existing changes contaminate *evidence*. Broad detection stays;
  `tools/checks/preflight.sh`, `actions/work.md` Step 5.75, and the Pre-Flight template now point at
  Current-REQ relevance instead of at a staging risk that doesn't exist.
- **Current-REQ relevance now covers session state too.** A run that notices a commit landing
  mid-flight, or another live process against the checkout, treats it exactly like an unexpected
  diff: preserve, exclude, continue. Added explicitly — **never probe for a concurrent session, and
  never ask the user to arbitrate one.** Exclusivity is the user's guarantee; a prompt asking them to
  confirm it rebuilds in conversation the coordination 0.161.0 deleted.
- A regression assertion pins the new clause, so the concurrency check can't be re-derived from
  first principles and quietly reappear as a four-option prompt that stalls the loop.

## 0.169.4 — Citation Guard Flags Mentions Instead of Idioms (2026-08-04)

The check that stops shipped files pointing at the maintainer doc was matching a hand-written list of
citation phrasings, and it turned out to catch none of the eight real ones in the tree — including
six comments in the board's `verify` code that had been shipping for months. It now flags any mention
and keeps a short per-file allowlist for the files that legitimately mean *your* `CLAUDE.md`.

- `_dev/tests/contract-regressions.sh` inverted: any `CLAUDE.md`/`AGENTS.md` mention in a shipped
  path fails, with a 14-file per-file allowlist carrying the judgement
- Six comments in `tools/queue-kanban/verify.go` and `verify_test.go` now state the release rule
  they describe instead of citing the maintainer doc
- `prompts/prompt-kit-step6-constraint-architecture.md` drops an attribution to a "CLAUDE.md
  standard" that was never in the file
- Measured, not assumed: old pattern caught 0 of 8, the inverted check catches 8 of 8, and a
  bare-prose mention (which a wider idiom list would still have missed) now fails

## 0.169.3 — Memory Recall's Shell-State Note No Longer Cites the Maintainer Doc (2026-08-04)

The memory action's recall procedure explained a rule by pointing at this repo's `CLAUDE.md`, which
is excluded from the tarball — so in every install that pointer went nowhere. It now states the rule
where you're standing, matching the nine other shipped files that already do.

- `actions/memory-reference.md` states the shell-state rule inline instead of citing `CLAUDE.md`
- Sweeping every shipped path for the same defect turned up six more sites in
  `tools/queue-kanban/verify.go` and `verify_test.go`, plus the fact that the contract suite's
  citation guard matches none of them — queued for your decision as REQ-093, not fixed here

## 0.169.2 — The Work Action Says Who Runs A Fan-Out (2026-08-04)

Worktree fan-out was described in one document and absent from the other, so `--wave N` read like a
concurrency switch while every step of the work action processes one request at a time. The two now
agree: running several builders is a procedure you or an advanced harness drive, not something
`do-work run` does.

- `actions/work.md` says plainly that it works one request at a time by design, and points at the
  procedure that does the other thing
- `--wave N` keeps doing what it always did — pick a batch of mutually independent requests — and now
  says so instead of implying parallelism
- The decision favours the floor: the work action stays followable by the simplest agent that can read
  files and run shell commands

## 0.169.1 — Hand-Back Merge Survives The Owner's Own Claim (2026-08-04)

Worktree dispatch told the orchestrator to merge a builder's branch, but the claim it had just made was
still sitting uncommitted — so on any project that keeps `do-work/` in git, the merge refused before it
started. The sequence now says to settle that first.

- New step 0 in the hand-back sequence: commit the claim moves, the checkpoint and the run directory
  before capturing the merge range's lower bound, so the bookkeeping stays outside the range
- Says plainly that it's a no-op to skip where `do-work/` isn't tracked, which is the common install
- Stated in both places the sequence is written, including the one an orchestrator actually follows

## 0.169.0 — Board Copy Carries the Ticket's Frontmatter (2026-08-04)

The drawer shows a ticket's status, domain, user request and timestamps in labelled rows — and **Copy**
dropped every one of them, handing you an anonymous body. It now copies the ticket file exactly as it
sits on disk, fence and all, so a paste can be saved straight back as a valid REQ or UR.

- Verbatim means verbatim: the original key order, comments and spacing survive, because the fence is
  taken from the file's own bytes rather than rebuilt from the parsed fields
- Works on user requests too, and on a live `do-work board` as well as a generated one
- Three drawer rows still can't come along — Tree, Overlapping write sets and Unblocks are worked out
  while the board is built, not stored in the file — and the board guide now says so
- Older boards missing their source-text bundle still fall back to the rendered text under a
  `# REQ-042: <title>` heading; no frontmatter is ever fabricated from what's on screen

## 0.168.6 — Fan-Out Dispatch Actually Run For The First Time (2026-08-04)

Worktree fan-out shipped three versions ago with its live acceptance test never run — two checkpoints
carried it as deferred, so the feature was a claim rather than a demonstrated capability. It has now
been run end to end, and the record is in the archive rather than in a promise.

- Two real requests built in two worktrees on two branches, integrated one at a time: both merged
  cleanly, each got its own changelog entry with an increasing version, and both branches came off with
  a plain `git branch -d` — no force needed
- A deliberately overlapping pair was confirmed to fail at the merge rather than combine silently
- Two defects found and queued: the hand-back merge collides with the owner's own uncommitted claim
  where `do-work/` is tracked, and nothing in the work action actually drives a wave
- Acceptance recorded as **Partial**, on purpose — the builders never overlapped in time, so genuine
  concurrency is still unproven and the record says which properties that leaves open

## 0.168.5 — Board And Verify Stop Handing Windows A Dead Command (2026-08-04)

The board's tooltips, its data warnings and `verify`'s remedy line all told you to fix a bad timestamp
with `date -u +%Y-%m-%dT%H:%M:%SZ` — a command that doesn't exist on Windows. The instruction layer was
fixed a release ago; these four display strings had been left carrying the old one.

- The tooltips and warnings now name the shape (`YYYY-MM-DDTHH:MM:SSZ`) and point at the Timestamp rule
- `queue-kanban verify`'s remedy keeps a runnable command, but it's `queue-kanban now` — same output on
  every platform, and you're already looking at that binary's output when you read it
- A fourth site turned up during the fix: the board's server-side future-timestamp warning had the same
  problem

## 0.168.4 — Checkpoint Bookkeeping Stated At Every Mover (2026-08-04)

A completed request could leave a stale entry behind in the session checkpoint, so the next run would
report a contradiction about a request that finished perfectly normally — and a warning that fires on
the happy path is the fastest way to teach you to ignore it. The two procedures that move a request out
of `working/` from outside the pipeline now say the entry goes with it.

- `do-work cleanup` Pass 0 and `do-work forensics` Check 1 each state the rule, citing its canonical
  home rather than restating it
- The work guide no longer says the checkpoint is written only at session end — it's written as each
  request is claimed, which is what lets a crashed run pick its own work back up

## 0.168.3 — Verify Catches A Builder That Committed Its Queue Edits (2026-08-04)

The check that stops a builder from writing queue state only looked at uncommitted changes, so it
caught a builder interrupted mid-write and missed one that finished and committed — the likelier shape,
since builders commit their work by design. It now also compares the builder's branch against the
integration branch, so the edits show up whether they were committed or not.

- New `worktree-committed-queue-state` finding, reported separately from the uncommitted one because
  the fix differs: drop the commits before merging, versus discard loose working-tree edits
- Uses a merge-base comparison, so a worktree that is merely *behind* the main tree stays silent —
  that legitimate case now has a regression test rather than a narrowed probe protecting it
- Neither state is marked fixable; discarding a builder's work is never mechanical
- `do-work forensics` Check 14's probe table now lists the two worktree probes separately

## 0.168.2 — Verify Marks Only Merged Worktrees As Fixable (2026-08-03)

`do-work forensics` was reporting every leftover builder worktree as something `do-work cleanup` would
mechanically resolve — including builders still mid-run and branches holding work that exists nowhere
else. Sent to cleanup, those cases stop and ask instead, so the count was pointing you at a command
that could not help. Verify now checks whether the branch is actually merged and says which case it
found.

- Merged residue reports as `merged-worktree-leftover [fixable]` and is the only kind counted toward
  `N fixable: run do-work cleanup` — matching what cleanup Pass 5 will actually do
- Unmerged work reports as `unmerged-worktree-leftover`, not fixable, and its remedy says out loud that
  verify cannot tell a live builder from a dead one
- A worktree whose branch is gone gets its own `worktree-merge-state-undetermined` state, so an
  unanswerable question is never reported as an answer
- `verify_test.go` gained a real git-repo fixture; every worktree probe used to skip silently in tests

## 0.168.1 — The Fan-Out Hand-Back File Has Somewhere Legal To Go (2026-08-03)

Fan-out dispatch made a per-builder `REQ-NNN-handback.md` mandatory, and the worktree rules forbade builders from writing the main tree — where `do-work/` lives. There was no location satisfying both, so an agent that got there had three moves and two of them corrupted the run quietly: write the main tree anyway, or write the worktree's copy where the file lands in the builder's branch and the orchestrator reads nothing.

- Sole-integrator now names exactly one path a builder may write — its own hand-back file, by absolute main-tree path, never staged or committed. Everything else stays forbidden, by name: a sibling's hand-back, `manifest.md`, anything else under `do-work/runs/`
- The "what lives under `do-work/`" rule is a condition now instead of a three-item list written before `do-work/runs/` existed — which is what let a reader argue the run directory was out of scope
- The path-resolution trap already documented for the inbound brief now says it applies in the return direction too, where it's worse: the write succeeds, into the wrong tree, silently
- Four contract checks pin the carve-out, so a maintenance pass can't delete it as redundant and restore the contradiction

## 0.168.0 — next-version Writes The Repo You Point It At (2026-08-03)

`queue-kanban next-version` takes the bump size as a positional and `--repo-root` as a flag, and Go's flag parser stops at the first non-flag argument — so every flag placed after the bump size was thrown away. The invocation the skill itself prescribed put `--repo-root` last, which meant it bumped whatever repo you happened to be standing in, exited 0, and printed a plausible version number. This is the tool's only write outside `do-work/`, and it was writing the wrong file.

- Flags now work on either side of the bump size, in any order
- Leftover tokens are an error on **every** subcommand instead of being ignored — that silence is what hid this, and the same shape was live on the flags-only subcommands too (`verify stray --repo-root X` was quietly reading the calling repo)
- The argument handling moved into a function without an `os.Exit` in it, so it can finally be tested — 11 table cases now cover every order, and they were watched failing against the old code first
- The prescribed invocation in `actions/work.md` puts the flags first, and a contract check pins that order rather than just the subcommand's name

## 0.167.4 — Capture Stops Stamping A Stray Instruction Into Every REQ (2026-08-03)

The Simple REQ template ended with a bare `Think carefully before answering.` *inside* the fence, so every request `do-work capture` wrote inherited it — 25 archived REQs carry the line. A builder spotted it once, correctly treated it as data rather than a directive, and logged the sighting in its own REQ. Nobody looked at the template, so it kept emitting for another 25 captures.

- The line is gone from the template; new captures end at `*Source: …*`
- A contract check pins it and two sibling phrasings, so it can't come back quietly
- The 25 archived REQs keep it on purpose — the archive is immutable, and the artifact is inert once it's history
- Swept `specs/`, `prompts/` and `interviews/` for the same shape and found nothing: their instructional prose sits outside the fenced bodies, which is the distinction that matters

## 0.167.3 — The write_set Guards Now Catch The Premise, Not One Phrasing (2026-08-03)

A guard that only matches the wording you already deleted is worse than no guard, because the green suite reads as coverage. The two checks protecting `write_set`'s display-only rule were in exactly that state — they pinned "one REQ at a time" and nothing else, so the retired premise could walk back in under half a dozen other phrasings, including the one the project's own lesson file calls the more dangerous of the two.

- The count-based pattern is defined once now (it was spelled out three times, which is why its gap was triplicated) and covers the class: `one`/`a single`/`only one`, `REQ`/`builder`/`coder`/`agent`, `at a time`/`at once`/`concurrently`, and the "is ever building" shape
- A second guard pins the weak fingerprint — justifying `write_set` by naming the exclusive-session model — across prose by line and across `model.go`/`board.js` by file
- Every one of them was demonstrated failing and then reverted; the old suite exits 0 on the same injected lines, which is the gap observed rather than argued
- `do-work cleanup`'s Pass 0 no longer argues its own safety from that premise either — it's safe because it sweeps only terminal statuses, which holds at any builder count
- Two other sites that name the exclusive-session model were checked and deliberately left alone: they cite it for what it actually says, and sweeping a phrase mechanically would delete true statements

## 0.167.2 — A Windows Timestamp Command That Actually Runs (2026-08-03)

0.167.0's notes said it fixed a real gap on Windows `cmd`. It didn't — the PowerShell form it shipped needs PowerShell 7 (stock Windows has 5.1), and it named a cmdlet as the remedy for `cmd`, which can't invoke one. Both were readable defects, not runtime surprises. This release is what makes that claim true, and it also does the subtraction that makes the fix reachable at all: eleven action files each spelled `date -u` inline, so an agent following one of them never got as far as the rule.

- The Windows form is now `(Get-Date).ToUniversalTime().ToString("yyyy-MM-dd\THH:mm:ss\Z")`, which runs on Windows PowerShell 5.1; `-AsUTC` is named only as the 7+ shorthand it is
- `cmd` gets a real entry point: `powershell -NoProfile -Command "…"` — `-NoProfile` so a profile banner can't land in the captured stamp
- Eleven inline copies of the command across seven action files are gone; the Timestamp rule is the one place that spells it, and a new check names any file that reintroduces a copy
- The rule was restructured rather than extended — a lead, three numbered sources, and two labelled trailing paragraphs, one of which finally covers the UTC *date-only* shape `memory/logs/` needs
- The Windows path is reasoned, not executed: nothing here can run PowerShell, and the review says so instead of claiming a test

## 0.167.1 — Crash Recovery Can Actually Recover A Crash (2026-08-03)

**0.164.0 quietly broke automatic crash recovery, and this fixes it.** That release made recovery check the checkpoint's in-progress record before touching a claimed REQ — a good gate, protecting work that exists nowhere but that file. But the record was only ever written at *session end*, which a hard crash never reaches. So the common case had no record, every crashed REQ was classified as someone else's claim, and it sat in `do-work/working/` untouched for good, one warning line per run, until someone ran `do-work forensics` and reset it by hand. The 0.164.0 notes described the gate; they didn't say the automatic path had stopped working. It had.

- The in-progress record is now written **at claim time** (Step 2), so a run that dies mid-REQ leaves the record its successor needs
- It's a list, one entry per claimed REQ — a crash during fan-out recovers every claim, not just the newest
- The entry is dropped whenever the REQ leaves `working/`, so a normal completion stops generating a phantom stranded-claim warning
- Explicitly **not** a lock: no exclusivity, no coordination, no second reader — the exclusive-session model is unchanged, and 0.164.0's protections (three-hour threshold gates the *offer* only, a human authorizes every takeover, an absent checkpoint stays ambiguous) are all intact
- A retired premise had left a third copy of itself in `actions/work.md` that the guard couldn't see; the two narrow guards are now one sweep over both pipeline files

## 0.167.0 — The Go Utility Can Now Produce Timestamps (2026-08-03)

`queue-kanban now` prints the exact UTC stamp every `*_at` field wants, and the Timestamp rule prefers it whenever the binary is already built. It joins `next-req` and `next-version` — the other two things the pipeline has to get right and can't eyeball. `date -u` stays the documented fallback, so no action ever needs a compiler.

- Fixes a real gap on Windows `cmd`, where the prescribed `date -u +FORMAT` doesn't exist at all — a PowerShell fallback is now named too
- The rule never tells you to *build* the tool for a stamp; a compile per timestamp would be worse than the shell command it replaces
- Amended in one place. The eleven sites that cite the rule were deliberately left alone
- Read-only, and the one subcommand that needs no `--repo-root` — it reads a clock, not your queue
- Writer now lives next to the readers that parse it, with a round-trip test through the board's own parser

## 0.166.2 — The Overlaps Badge Explains Itself Correctly Again (2026-08-03)

Eleven places told you the board's file-overlap badge was harmless *because only one REQ runs at a time* — a reason that stopped being true when fan-out dispatch landed. The badge's behavior never changed; the explanation did. It is advisory input to your pick, and the merge is what proves two builders didn't collide.

- Fixed in the `write_set` schema itself, both board-action sites, the board and work guides, the tool's prime, four Go comments, and the badge tooltip you actually read in the browser
- Each site now points at the Fan-Out Dispatch contract instead of restating it
- New contract assertion, in two parts: a line sweep for prose plus a file-level guard for the Go and JS sources, whose comments wrap and so slip past a line-based check
- No behavior change anywhere: no Go logic, no schema field, no board column touched

## 0.166.1 — Crash Recovery Records When It Reset a REQ (2026-08-03)

A REQ the pipeline re-queues after a crash now carries the instant it was recovered, so the board's waiting-time figure counts from the reset rather than from the day the request was written. The by-hand reset in forensics already did this; the automatic one had never done it.

- Crash Recovery substep 1 stamps `status_changed_at` on the reset, on both flip targets (`pending` and `pending-answers`), with the preserved-`blocked` case explicitly excluded — its `blocked_at` is intact
- The same substep removes `claimed_at`, which is why the stamp matters: without it a recovered REQ has no trace of the reset at all
- New contract-suite assertion pins the prescription, so a later edit can't quietly drop it

## 0.166.0 — Several Builders at Once, Under One Queue Owner (2026-08-03)

You can now hand several REQs to several builders at the same time — each in its own git worktree — without any of the locks, heartbeats, or liveness checks that got deleted at 0.161.0. The trick is that none of that was ever needed: worktree dispatch was already written per REQ, so only two sentences were capping the builder count. The boundary that matters is *who owns the queue*, not how many builds run.

- The operating invariant is now **one queue owner per checkout** — the single session that claims, flips status, and archives. Builders aren't owners, so any number can build at once
- Two queue owners on one checkout stays out of the contract, stated more explicitly than before — nothing that was removed at 0.161.0 comes back
- New Fan-Out Dispatch section: you pick which REQs run together, `write_set` overlaps are advisory input to that pick and never a gate, and the merge refusing is the only real proof two builders didn't collide
- It says plainly where the merge gate stops helping: git detects conflicts by line proximity, not meaning, so two REQs each appending to a shared registry merge cleanly and can still be jointly wrong
- Named serial-only list: queue transitions, REQ id allocation, and `actions/version.md` + `CHANGELOG.md` — one entry per REQ, written by the owner at merge time
- Integration stays serial. The saving is in the build phase, and the docs say so instead of promising more
- A worktree per builder is mandatory, with the reason kept in the prose: sharing one tree means every test run and review diff reads the other builder's unfinished edits
- The version-bump-and-changelog ritual in `CLAUDE.md` is now scoped to the integrating commit, so a builder in a worktree skips it instead of racing its siblings

## 0.165.0 — REQ and Version Numbers Get Allocated, Release Rules Get Checked (2026-08-03)

Two of the rules you're asked to hand-check on every single commit — the version must beat the newest changelog entry, and the entry title must not already be in use — are documented as having already been gotten wrong. They're now machine-checked, and the board tool can hand you the next REQ number and version instead of you eyeballing a file listing for it.

- New `queue-kanban next-req` prints the next free REQ number, running capture's own scan (gaps tolerated, `archive/UR-*/` included)
- New `queue-kanban next-version <patch|minor|major>` bumps the version line, reads it back to confirm the write landed, and prints the number — the bump size is always yours to name, never inferred
- New `queue-kanban verify` reports eight invariants read-only: version/changelog agreement, duplicate or reused entry versions and titles, duplicate REQ numbers, a checkpoint naming a REQ that no longer exists, untrustworthy `claimed_at` stamps, finished REQs stranded outside the archive, and `worktree-agent-*` leftovers
- `verify` marks the findings `do-work cleanup` can actually fix and points you there; it never repairs anything itself, and it never writes `CHANGELOG.md` — that stays a human write
- Wired in as optional accelerators: capture (REQ numbers), the Step 9 commit ritual (version bump), and forensics as Check 14. Each falls back to the existing manual procedure when the Go toolchain is absent
- A probe that can't run (no git, no version line, a changelog in another convention) is reported as skipped rather than passing quietly

## 0.164.0 — Crash Recovery Asks Before Discarding a Claimed REQ (2026-08-03)

Crash and restart thirty seconds later, and recovery used to throw away the plan, exploration, and scope the interrupted run had just finished writing — nothing is committed before Step 9, so those sections were gone for good. Recovery is now conditional: it only resets what the checkpoint records as this session's own interrupted work, and everything else is left exactly as it is.

- A claimed REQ in `do-work/working/` that `do-work/CHECKPOINT.md` doesn't name as in-progress is reported and left byte-identical — no frontmatter reset, no section stripping, no move
- Takeover is offered only once a claim is over three hours old, and a human always authorizes it; the threshold reports, it never decides
- An unparseable, future-dated, or missing `claimed_at` makes a REQ immediately eligible for that offer, so a corrupt stamp can't protect a dead claim forever
- No checkpoint at all is treated as ambiguous, not as permission — which is the common case, since a hard crash never gets to write one
- Unattended runs never stall on the prompt: they report the claim and continue to the next REQ
- Step 1 now reads the checkpoint *before* recovery and says why — it's recovery's input, not just resume context

## 0.163.3 — Board Copy Includes the REQ Id and Title (2026-08-02)

The drawer's Copy button used to hand you the body text alone — frontmatter holds the id and title, so a paste landed somewhere else as anonymous prose. Every copy now leads with a `# REQ-236: <title>` heading.

- Applies to URs too, and to the rendered-text fallback path used when a generated bundle is missing `board-markdown.js`.
- When the body already opens with an H1 restating the title (the usual REQ shape), that duplicate line is dropped rather than repeated under the new heading.

## 0.163.2 — Audit Sweep: Last Stale Release-Flow Comment (2026-08-01)

A post-removal audit found one straggler: a `model.go` comment still listed the removed release action among `status_changed_at`'s writers. Comment corrected; the audit found nothing else — numbering, routing, board fields, and docs are all consistent with the reservation removal.

## 0.163.1 — Drop the Last Parallel-Actions Sentence and Fix Roadmap Numbering (2026-08-01)

Re-validation caught two leftovers. The Execution Model no longer reasons about parallel execution at all — the boundary now ends at "behavior is unspecified" instead of enumerating which actions may run alongside a build (the session assumes it is alone; that's the whole rule). Also renumbered roadmap's Suggested Next Steps template, which skipped from 7 to 9 after the reservation line was deleted.

## 0.163.0 — Remove the Reserve Action and the Read-Only Parallel Carveout (2026-08-01)

Validation of the exclusive-session work flagged two contradictions it had left standing, plus a router-budget overrun. All three are fixed.

- **The reserve action is gone.** `do-work reserve`/`release`, the `reserved` status, and the `reserved_for`/`reserved_at` fields existed solely to allocate REQs to *another* worktree/cloud session — cross-session ownership, which the exclusive-session contract declares unsupported. Removed from the router, work pipeline, roadmap, forensics, cleanup, abandon, help, the board tool (parser, column logic, badges, drawer row, CSS), and the docs. A leftover `status: reserved` REQ is caught by the existing unrecognized-status guards — the pipeline skips it with a warning and the board shows it under Needs input / Blocked with a fix hint; edit its status back to `pending`.
- **The parallel carveout no longer blesses queue writers.** The Execution Model note had listed `board` and the reviews as read-only — but board's Testing view writes REQ frontmatter and `testers.md`, and the reviews append REQ metadata. The rule now permits only actions that claim nothing and write no queue or REQ state, still within the ≤200-word replacement cap.
- **SKILL.md is back under its 2,650-word budget** (2,588) via the removed reserve routing/dispatch rows. The contract suite now ratchets the removal: reservation vocabulary is a forbidden token across shipped files, and `actions/reserve.md` must stay deleted.

## 0.162.4 — Reconcile Reservation, Maintainer-Contract, and Token-Canonicalization Docs (2026-08-01)

Codex re-review of the final tree caught three consistency gaps the UR/exclusive-session work left behind:

- **Reservation override.** `actions/reserve.md` and `actions/prime-req-reservation.md` still said "only the default full-queue scan honors reservations" — but a `do-work run UR-NNN` scoped run now honors them too. Corrected to "only explicit per-REQ naming overrides a reservation".
- **Maintainer contract.** `CLAUDE.md`'s board lock-step rule still told maintainers "which REQs may co-dispatch stays with the dispatch gate" — machinery REQ-069 deleted. Reworded to the exclusive-session, display-only framing.
- **Token canonicalization.** The Target ID Resolution contract advertised case-insensitive `req-42` but callers glob/compare against zero-padded upper-case stored ids. Added a canonicalize-before-lookup step (uppercase prefix, numeric-value digit match) so `req-42` actually resolves to `REQ-042`.

## 0.162.3 — Sweep the Last "Overlaps Everything" Mis-Framings (2026-08-01)

Verification found two more spots that inherited the deleted dispatch gate's "absent write_set ⇒ overlaps everything" wording — `actions/capture-reference.md` and a malformed-glob parenthetical in `actions/work-reference.md`. Both corrected so all five board-story files agree: an absent/empty `write_set` gets **no** overlaps badge (unknown, not conflict), and the docs make no false universal claim about malformed globs.

## 0.162.2 — Finish the Board Comment Sweep (2026-08-01)

Re-verification caught one `tools/queue-kanban/model.go` comment (`isWriteSetOverlapCandidateStatus`) still using the removed "a dispatcher could put these in flight together" framing while its four siblings had been refreshed. Reworded it to the exclusive-session, contention-heads-up wording so the file tells one story. Comment-only; the Go tests still pass.

## 0.162.1 — Correct the Board Overlap Semantics After the Exclusive-Session Cut (2026-08-01)

An adversarial verification pass over the batch caught a factual error introduced by the concurrency-machinery removal, plus stale comments. Fixes:

- `actions/work-reference.md` said an absent/empty `write_set` makes the board's overlaps badge "read as overlaps everything." The board does the opposite — it shows **no** badge for an empty set (absence = unknown, not conflict). Corrected the schema note and the crash-recovery rationale to match `actions/board.md`, `tools/queue-kanban/prime-do-kanban.md`, and the board code.
- `actions/work.md` `--wave` restatement now says "any targeting token (`REQ-`/`UR-`)" instead of "targeted REQ IDs", matching the authoritative Input guard.
- Refreshed `tools/queue-kanban/model.go` and `web/board.js` **comments/tooltip** that still described a "dispatch gate" / "co-dispatch" decision the exclusive-session model removed (comment-only; the board's overlap display is unchanged and its Go tests still pass).

## 0.162.0 — REQ IDs Accepted by `do-work roadmap` (2026-08-01)

Closes the inverse of the asymmetry this batch fixed: `roadmap` took `UR-NNN` as a scope token but not `REQ-NNN`, so `do-work roadmap REQ-067` silently returned a whole-queue survey. Now every id-taking action in the skill accepts both prefixes.

- `do-work roadmap REQ-NNN` scopes to a single REQ — its status, dependency position, feasibility read, and its UR siblings for context. Deliberately thin (per-REQ detail stays with `do-work inspect`).
- Multiple id tokens resolve to their union; the soft unrecognized-argument fallback stays, so a genuinely unknown token still yields the full survey with a note.
- Cites the shared Target ID Resolution contract for token shapes rather than restating them.

## 0.161.1 — UR Token Precedence and Reserve Mode-Table Fixes (2026-08-01)

Two corrections from PR review of the UR-id work.

- `do-work reserve`/`release` mode table now lets either `REQ-NNN` or `UR-NNN` lead — a UR-only `do-work reserve UR-011 for cloud-alpha` or `do-work release UR-011` no longer misses every row (or gets read as a free-text label).
- `do-work run` now defines precedence for a REQ reached both by explicit name and by UR expansion (`do-work run UR-011 REQ-068`): explicit naming wins, so the deduped target takes the named branch (bypasses `depends_on`, claims a reservation).

## 0.161.0 — Exclusive-Session Model Replaces the Concurrency Machinery (2026-08-01)

The work pipeline now states plainly what it always assumed: one `do-work` session, one active REQ, one coder context. The ~6,500 words of orchestrator-lock, parallel-dispatch, and co-dispatch-re-validation machinery that existed to detect and recover unsupported concurrent runs are gone, replaced by a short operating rule. Behavior for the normal single-session run is unchanged; the pipeline is just far smaller and easier to follow.

- Removed: the `Concurrent-Orchestrator Lock Guard` section, the Step 1 parallel-dispatch gate and serial-only rule, the Step 3 / Step 5.5 co-dispatch re-validations, the crash-recovery concurrency gate, cleanup's Pass 0 live-claim gate, and every orchestrator-lock heartbeat/claim touchpoint.
- Added a `## Execution Model — Exclusive Session` rule: unexpected repo state matters only when it blocks the active REQ; the coder stops after three consecutive fix attempts; and **read-only actions (roadmap, board, inspect, forensics, recap, reviews) may run in parallel** — the boundary governs writers only.
- Kept: crash recovery of an interrupted single session, single-builder worktree isolation, and the `write_set` field (now display-only, feeding the board's overlaps badge; the board tool is untouched).

## 0.160.0 — UR IDs Accepted by `do-work abandon` and `do-work reserve` (2026-08-01)

The last two REQ-only actions now take a UR. `do-work abandon UR-011` cancels the UR's cancellable members (including any `failed` one still holding it open) behind a single itemized confirmation; `do-work reserve UR-011 for cloud-alpha` reserves its pending members; `do-work release UR-011` returns them. Both cite REQ-067's Target ID Resolution contract instead of re-deriving the rule.

- Each action resolves tokens in front of its existing per-target loop, so every per-REQ gate (status refusals, the `claimed`/`reserved` extra confirmations, reserve's pending-only capture, release's reserved-only touch) applies unchanged to expanded members.
- Bulk cancel is protected by one prompt that itemizes every resolved target with a total count — no `--yes` or per-member bypass.
- `release <token>` precedence is stated explicitly: a `UR-` token resolves as an id, so a reservation *label* literally named `UR-011` is released by naming its REQ ids, not by `release UR-011`.

## 0.159.0 — UR IDs Accepted by `do-work run` (2026-08-01)

`do-work run UR-011` now works: a UR argument expands to its member REQs and runs them in dependency order, instead of being rejected as an unrecognized token. A new shared Target ID Resolution contract defines the token grammar once so the other id-taking actions can cite it rather than each restating it.

- New `Target ID Resolution` contract in `actions/work-reference.md` — `REQ-`/`UR-` token shapes (case-insensitive), UR→REQ expansion by scanning `user_request:` frontmatter (never the `requests:` array), and an empty resolution that stops the action instead of falling through to a full-queue run.
- `do-work run` accepts `UR-NNN` alongside `REQ-NNN`; an explicitly-named REQ bypasses `depends_on`, a UR-expanded member does not. `--wave` stays mutually exclusive with any targeting token.
- The unrecognized-argument guard is unchanged — a bad token like `REG-042` still errors, never a silent full-queue run.

## 0.158.1 — Pre-Flight Tells a Failed Test Run From a Test Command That Never Ran (2026-08-01)

The pre-flight check records a test baseline so the reviewer can tell a pre-existing failure from a regression the builder just introduced. It was recording one even when the test command never launched — a typo'd command exits 127, which looked exactly like a red suite, so a session could start with a baseline that described no test run at all and later excuse a real regression as "already failing before we started". The whole point of the check is attribution, and that was the one way it could get attribution wrong.

- A command that can't be launched (exit 126/127) now reports `WARN: could not run the test command — no baseline recorded`, writes no failures file, clears any stale one, and records `"launched": false`. The review step reads that field and refuses to compare rather than comparing against a fiction.
- `preflight.sh "npx jest --silent"` works now. Passing the test command as one quoted string used to die with an opaque `command not found` — which is how the bug above got triggered in the first place. The quoted form runs via `sh -c`, so `"cd app && npm test"` works too.
- The baseline path is written once instead of spelled two different ways, and the script's header documents that it runs from the project root like its siblings.

## 0.158.0 — An Answer You Give Mid-Run Is Written Into the REQ (2026-08-01)

The question rulebook governed how to ask you something and said nothing about where your answer had to land — so a long run could stop, ask two decisions that were genuinely yours, get detailed answers, and write none of it down. The next builder started fresh, found one question still unanswered, and re-decided it off the stored recommendation whose reasoning you had just rejected. Same outcome, wrong reason, no record.

- New `crew-members/clear-questions.md` Principle 8: an answer obtained interactively outlives the transcript or it may as well not exist — write it into the durable record before acting on it, and capture any new work it implies as its own REQ. Scoped to decisions the work is later read from; a plain proceed/abort gate needs nothing extra.
- `actions/work.md` Step 3.5 gained the missing third branch. It had two — the builder decides (`- [~]` plus a D-XX), or the question waits for `do-work clarify`. Now: a decision escalated mid-run and answered is written in by the orchestrator before dispatch as `- [x]`, never `- [~]` and never with a D-XX, because it is your decision and not the builder's guess.
- `actions/clarify.md` Step 4's `- [x] question → answer` form is now a named entry point, canonical for any caller that obtains an answer — cited by name so a renumbered step can't dangle. Step 5's status flip stays clarify-only; a REQ already in flight has no `pending-answers` to leave.
- Three assertions in `_dev/tests/contract-regressions.sh` pin the principle, the branch, and the named format so a later cleanup pass can't drop one half and leave the others pointing at nothing.

Mid-run questions are deliberately still allowed — asking was the right call. The bug was losing the answer.

## 0.157.1 — Board Pills and Drawer Fields Wrap Instead of Overflowing (2026-07-31)

Long badge values on the Kanban board — most visibly a card's "blocked by \<condition\>" pill — used to run off the card edge because badges were locked to a single line. Pills now wrap onto extra lines, and the drawer's metadata values break long unbroken tokens (write_set paths, blocked_check commands) instead of widening past the drawer.

- `.badge` drops `white-space: nowrap`: wraps within the card, keeps the small-caps label ("BLOCKED BY", "NEEDS") on one piece.
- Drawer `detail-meta` values get `min-width: 0` + `overflow-wrap: anywhere` so a single long path can't stretch the grid.
- Testing-view meta chips wrap the same way.

## 0.157.0 — Source Repo Now Tracks Its Own Queue and Knowledge Base (2026-07-31)

This repo had been keeping its own `do-work/` and `kb/` untracked via a local `.git/info/exclude` entry, on the theory that committing them would leak the maintainer's queue into consumer installs. That theory was already false — `.gitattributes` `export-ignore` plus the tar `--exclude` flags do that job, and the tarball has never contained either directory. Meanwhile the blanket ignore was costing real safety: several of the skill's own guards only work on tracked files, so the repo that dogfoods do-work was running with them silently disabled.

- `do-work/` and `kb/` are now tracked here, matching the Trail of Intent the skill tells consumers to commit. That re-arms `record-commit-hash.sh`'s HEAD size-floor and numstat data-loss guards, and makes `do-work cleanup` Pass 6 blanked-REQ recovery functional instead of a permanent no-op. It also ends a hybrid state where 25 archived REQs were tracked and 54 weren't.
- `/kb export-ignore` added to `.gitattributes`, and `--exclude='kb'` added to all six install/update tar invocations (README, `actions/version.md` ×3, `tools/do-work-update.sh` ×2). Consumers' own `kb/` was already safe — extraction never deleted it — this keeps *upstream's* KB from landing in their skill directory.
- New contract regressions assert `/do-work export-ignore`, `/kb export-ignore`, and the `kb` tar exclude. Those lines are now the only barrier to shipping the maintainer's archive, so they get a ratchet rather than a comment.
- Only genuinely-transient runtime state stays locally excluded: the orchestrator lock and its mutex files, plus preflight's `do-work/working/baseline.json` and `baseline-failures.txt`.
- Caught on the way in: archived `REQ-034` had 833 KB of raw Explore-agent session JSONL pasted into its `## Exploration` section where the agent's summary belonged. Verbatim session capture — UUIDs, prompt text, local paths — has no business in permanent history, so the block is replaced by a note explaining the removal; the REQ's Scope/Verification/Decisions sections are untouched (844 KB → 32 KB).
- Fixed a stale claim in `actions/note.md` that `do-work/pipeline.json` is kept out of git "via the shipped `.gitignore`" — `.gitignore` is itself export-ignored, so no `.gitignore` ever ships. `actions/pipeline.md` Step 4's `.git/info/exclude` entry is the real mechanism.

## 0.156.1 — Verify Pins Removed Lines to the Parent's Frontmatter, Updater Warns About Uncommitted Edits (2026-07-31)

Two follow-ups from an external review of 0.155–0.156. Both were narrow, and both had a way of losing content quietly.

- `record-commit-hash.sh --verify` used to accept *any* removed line starting with `commit:` as "the old field being replaced". Archived REQs that quote the schema have `commit:` lines in their **body**, so a hook that deleted one while the write-back inserted the frontmatter line netted +1/−1 and passed — with a message claiming the patch was a single line. The removal side is now measured against the parent commit's actual frontmatter: zero removals on an insert, exactly that line on a replace.
- The updater's recovery instructions read like a full undo. They aren't: git restores what was **committed**, so an uncommitted edit to a shipped file dies at the extraction now that no rollback copy is kept. It now says that before the confirmation prompt — with the files named and `git stash` suggested — and repeats in the recovery output that the printed `git checkout` won't bring those edits back. The `git clean` line also flags that a root install's shipped paths hold project-owned files.

## 0.156.0 — Update Script Keeps No Rollback Copy (2026-07-31)

`just run-do-work-update` no longer duplicates your whole install before extracting. Version control is the undo, and copying a tracked tree on every run buys nothing git does not already hold — it just left a `do-work.preupdate-<timestamp>.bak` sitting in the project as untracked noise.

- The `cp -R` snapshot and the automatic restore are gone. The prompt now reads "Files are overwritten in place and no rollback copy is kept."
- A failure inside the destructive region reports the partial install and prints recovery commands you can paste — `git checkout --` scoped to the shipped paths git actually tracks, then `git clean -nd --` to review what the extraction added before deleting it.
- If the install is **not** tracked in git (a project that gitignores `.claude/`, or no repo at all), you get told that before the confirmation prompt, because there nothing can be restored.
- New `_dev/tests/update-script-behavior.sh` runs the real updater against a synthetic install with a stubbed upstream fetch: the happy path leaves no `.bak`, a mid-extraction failure reports instead of restoring, and a declined update changes nothing. A contract check fails the build if a `cp -R` snapshot is ever reintroduced.

## 0.155.0 — Commit-Hash Verify Inspects the Committed Patch, Partial Restores Fail Loudly, Updater Flags Stale Files (2026-07-31)

Three fixes from an external review of the data-loss guards shipped in 0.151–0.153. The headline: `--verify` was making a promise it couldn't keep against the commonest kind of pre-commit hook.

- `record-commit-hash.sh --verify` now asserts the **committed patch** — HEAD introduced exactly one new line for the REQ, and that line is the `commit:` field. It used to compare the committed blob against the worktree, which proves nothing when a hook rewrites the file and re-stages it: both sides move together, so a body could be silently gutted while the sizes agreed. Where the patch can't be isolated (root commit, merge HEAD, file added by that same commit, or `--verify` run too late) it now says so and labels the weaker guarantee instead of reporting a clean pass.
- `blanked-req-scan.sh --restore` no longer reports a partial repair as a repair. Content restored but its recorded hash rejected is its own outcome now: exit 1, a `FAIL:` line, and the write-back's own diagnosis passed through rather than swallowed. Previously that file counted as restored and `do-work cleanup` Pass 6 was told everything was fixed — then committed it with provenance pointing at nothing.
- `do-work update` lists files sitting in your install that upstream no longer ships. The extraction only overwrites, so a deleted action or check used to survive downstream forever while the post-update audit reported clean — the old filter dropped "upstream removed this" and "that's your `.DS_Store`" on the same line. Reported, never auto-deleted: your own file in a shipped directory looks identical from there. (The `prompts/`/`interviews/` pre-clean already covered those two directories; this covers the rest.)
- Four new probes in the guard fixture, including a re-staging hook that rewrites the body without changing what a blob read-back can see.

## 0.154.0 — Project Justfile Left Alone, Failed Updates Roll Back, Lock Race Closed (2026-07-31)

Three follow-ups from a review of the update script and the orchestrator lock. The headline: a failed `do-work update` now restores itself instead of leaving you a partial install and a path to fix by hand.

- `do-work update` no longer overwrites a project-owned justfile when the skill lives at the project root. It records which of `justfile`/`Justfile`/`.justfile` the project actually uses — by real directory entry, since a case-insensitive filesystem makes `[ -f justfile ]` match a `Justfile` — and restores that exact name and content across the extraction. A nested install's justfile is the skill's own and still gets refreshed.
- Any failure inside the destructive region now rolls the install back automatically: shipped paths are restored from the rollback copy, files the update added are cleared, and the copy is kept as the audit trail. Previously it printed the backup path and left the partial install in place.
- The lock mutex's remaining lost-update window is closed, not just narrowed. The staged lock image now lives *inside* the mutex directory, so an evicting `rm -rf` destroys it and the evicted owner's publishing rename can only fail — and that rename now fails closed (exit 3, re-acquire) instead of swallowing the error.

## 0.153.1 — Guard Fixture Lints Clean (2026-07-31)

Housekeeping in the commit-hash test fixture, plus one assertion that was written but never wired up.

- The scan probes now assert the *recoverable* byte count reported from the pre-blanking commit, not just that a file was found — the number an operator actually decides on. It was computed and then dropped on the floor.
- `_dev/tests/record-commit-hash-guards.sh` is shellcheck-clean at default severity.

## 0.153.0 — Cleanup Restores Blanked Archived REQs (2026-07-31)

Forensics could tell you an archived REQ's content had been destroyed; now cleanup can put it back. The recovery used to be hand-rolled git archaeology, one file at a time, against a deadline — the lost content only survives until `git gc` collects it.

- New `### Pass 6` in `do-work cleanup`: shows a dry run first (which file, which commit, how many bytes, which hash goes back), asks before writing, and restores only what you approve. Unattended runs report and stop.
- `tools/checks/blanked-req-scan.sh --restore` does the work — temp file plus atomic rename, refuses to write recovered content that's itself empty, and re-applies `commit:` by calling `record-commit-hash.sh` rather than hand-editing frontmatter, so the guards come along.
- `--dry-run` writes nothing and keeps the finding exit code; a completed repair exits 0, because a fixed thing is not a finding.
- A file git has no non-empty version of is reported as a permanent loss, never silently skipped.
- Seven restore probes added to the git fixture suite, including the full six-file incident reproduction and a byte-identity assertion against the pre-blanking blob.

## 0.152.0 — Forensics Detects Blanked REQ Files (2026-07-30)

`do-work forensics` can now tell "this REQ's content was destroyed" apart from "this REQ has a typo in its status." Previously a 0-byte archived REQ showed up as an `unrecognized status ''` warning whose suggested fix — edit the status field — would have written over a file that needed recovering first.

- New check 13 flags any REQ/UR file that is 0 bytes or has lost its frontmatter, as **Critical**, with the recovery commit and the recorded implementation hash already resolved from git history.
- New `tools/checks/blanked-req-scan.sh` does the scanning and the history walk. Read-only, so forensics keeps its never-modifies-anything contract; `--porcelain` emits machine-readable records.
- Check 11 now skips files with no parseable frontmatter, so a destroyed file is reported once, with the remedy that fits.
- A file with no non-empty version anywhere in history is reported as unrecoverable rather than silently passed over.

## 0.151.0 — Guarded Commit-Hash Write-Back (2026-07-30)

The Step 9 "record commit hash" step used to be prose — write the hash into the archived REQ's `commit:` field, then commit. In a repo using do-work, that free-form edit truncated six archived REQ files to 0 bytes, destroying 9 KB to 26 KB of decision trail each, with commit messages that claimed success. It's a script now, and the guards run before anything is staged.

- New `tools/checks/record-commit-hash.sh`: edits only the `commit:` line inside the frontmatter block, and refuses to write unless the rewrite changed exactly that one line. A `commit:` quoted in body prose is structurally unreachable.
- Guards for the real failure shapes — an already-blanked or truncated REQ (checked against its size in `HEAD`), an unterminated frontmatter block, duplicate `commit:` keys, CRLF, a hash the repo can't resolve, and the literal `<hash>` placeholder pasted straight out of the docs.
- Running it twice is a no-op, but an edit that never got committed (a pre-commit hook rejected it) is detected and reported for committing rather than silently skipped.
- A `--verify` mode reads the committed blob back, which is the only way to catch a content-mutating pre-commit hook rewriting the file after every other guard passed.
- Works the same on the worktree-dispatch path, where the hash is the `--no-ff` merge commit; degrades cleanly outside git or where `do-work/` is git-excluded.
- First behavioral test fixture in the repo: `_dev/tests/record-commit-hash-guards.sh` builds a throwaway git repo and asserts each guard actually fires.

## 0.150.15 — Update Script: Guard Project Docs, Clean Non-Interactive Cancel, Rollback Pointer on Failure (2026-07-30)

Four fixes to `do-work update`'s helper script after a code review of the new just shortcut. The headline: running it where the skill *is* the project root (a dev repo or direct clone) no longer risks deleting your project's own `CLAUDE.md`/`AGENTS.md`, and a piped or CI invocation now cancels cleanly instead of dying mid-prompt.

- Non-interactive stdin (piped, CI, `</dev/null`) now defaults to No and exits 0 — previously the bare `read` returned non-zero at EOF and `set -e` aborted before the cancel branch could run.
- The stale-vendored-doc cleanup (`rm` of `CLAUDE.md`/`AGENTS.md`) runs only for a nested install; when `skill_root == project_root` those are the project's own instruction files and are left untouched.
- `justfile` joined `shipped_paths`, so its overwrite now shows in the pre-confirmation diff and the uncommitted-changes warning instead of happening silently off-list.
- A failure mid-update (e.g. ENOSPC during extraction — the `cp -R` backup just doubled the on-disk size) now always prints the rollback-copy path, closing the one path that left a half-updated install with no pointer to the backup.

## 0.150.14 — Update Check No Longer False-Fails on Stray Local Files (2026-07-30)

`do-work update` could report a clean, successful update as broken — and steer you toward a needless rollback — whenever a Finder `.DS_Store`, an editor swapfile, or any other local-only file sat under a shipped path. The post-update integrity check now ignores install-only extras while still catching genuine extraction failures.

- Post-update check drops `--new-file` and filters install-only "extras" with a fixed-string match, so a stray `.DS_Store` or `*.swp` no longer trips the false `exit 1`.
- The filter is metacharacter-proof (`grep -vF`), so a `[`, `+`, or `.bak` in the install path can't defeat it.
- Added an explicit missing-file check so a wholly-absent shipped file or directory is still caught as a real failure — a gap that dropping `--new-file` would otherwise have opened.

## 0.150.13 — Project-Local Just Update Shortcut (2026-07-29)

Projects that install the do-work just recipes can now update their local skill without spending an agent turn.

- `do-work install just-kanban` now adds `just run-do-work-update`, which checks the upstream version, shows the installed-versus-upstream diff, asks before overwriting, creates a rollback copy, and preserves the runtime `do-work/` directory.
- Existing justfile recipes detect the new command as drift and offer the same consent-gated upgrade as other shipped recipe changes.

## 0.150.12 — Board No Longer Flags assets/ Deliverable Copies as Duplicate REQs (2026-07-29)

The Kanban board (`do-work board`) walked every `REQ-*.md` file, including deliverable copies parked under a UR's `assets/` folder. Those attachments have no frontmatter `id`, so their id fell back to the filename and collided with the real ticket — producing a spurious "duplicate REQ id" data warning (or a phantom card, for a uniquely-named asset).

- The walk now prunes any `assets/` folder at any depth, alongside the existing `deliverables/` and `runs/` exclusions.

## 0.150.11 — Abandon Resolves a Failed REQ So Its UR Can Close (2026-07-29)

A `failed` REQ had no way out — nothing in the skill could move it off `failed`, so any User Request holding one stayed open forever (the gap 0.150.10 uncovered). Now `do-work abandon` cancels an already-archived failed REQ in place, flipping it to `cancelled` so its UR can close, while keeping the failure record intact.

- `do-work abandon` accepts a `failed` REQ (at `do-work/archive/` root or `legacy/`) and cancels it in place — no move, and the failure signal (`error`/`error_type` plus a `## Cancelled` "Previously: failed" note) is preserved.
- **The one thing to know:** a completed follow-up REQ recovers the work but never flips the original out of `failed` — cancelling is the only transition, needed whether or not a follow-up ran. This corrects the 0.150.10 note that framed a follow-up as resolving it.
- `cleanup` Pass 1 and `forensics` Check 6 now point you at `do-work abandon` when a failure is holding a UR open; the failure-resolution rule lives canonically in work-reference.md's Terminal-resolved statement, with the three closure readers deferring to it (no restated copies to drift).
- No board or schema change — `cancelled` was already terminally-resolved; the change is purely which inputs abandon accepts and how the resolution is documented.

## 0.150.10 — UR Closure Keys on the Terminal-Resolved Set Everywhere (2026-07-29)

work.md Step 8's archive table was the last reader still counting `failed` as closing a User Request — cleanup Pass 1 and forensics Check 4 already keyed on work-reference.md's terminal-resolved set, so the two halves of the pipeline disagreed on whether a failed REQ holds its UR open. Now all three readers cite the one canonical set: a `failed` REQ keeps its UR open until a follow-up resolves it.

- Step 8's row cites the set instead of restating it; the canonical paragraph's caller list is now marked illustrative and includes forensics Check 4.
- docs/cleanup-guide.md and docs/forensics-guide.md stop describing closure as "all REQs archived" (a failed REQ is archived but not resolved).
- Discovered en route: nothing in the skill can resolve a `failed` REQ at all — queued as a follow-up decision (REQ-060, pending-answers).

## 0.150.9 — Blocked-Flip Guard Judges Worktree Builders by Their Branch (2026-07-29)

Step 8's blocked-vs-failed call used `git diff` on the main tree to ask "did edits land this attempt?" — but a worktree builder commits on its own branch, so the main tree always reads clean and real work got wrongly parked as `blocked`. The guard now reads the builder's branch in worktree mode.

- Case order: a completed hand-back merge (`<merge_hash>` held) proves edits landed; otherwise a quoted `git rev-parse --verify -q '<operative_name>'` existence probe (a missing branch = genuine before-any-work, and `rev-list` on it would exit fatal instead of printing a count); only then `git rev-list --count HEAD..<operative_name>` decides.
- Judged from git state, never the builder's handed-back manifest; serial-mode behavior untouched.

## 0.150.8 — Merge-Aware Diff Reads for Worktree-Merged REQs (2026-07-29)

Every consumer that reads a REQ's `commit:` hash as a diff source now detects the worktree `--no-ff` merge case and diffs against the first parent — plain `git show` on a merge prints a near-empty combined diff, so standalone reviews and receipts of worktree-merged work silently saw nothing.

- Shared idiom at all sites: detect via `git rev-parse --verify -q '<sha>^2'` (quoted — `^` is special in zsh/cmd.exe), then `git show --first-parent -m <sha>`; ordinary serial commits unchanged.
- Covered: review-work.md Get-the-Diff + Two-Modes table, present-work.md (three sites incl. the interactive explainer's receipt), pipeline.md's Completion-Report bullet + pipeline-reference.md's rendering template, ai-report.md's Verify-It-Yourself spec.
- The maintainer shell-trap catalog gains the merge-commit/empty-combined-diff trap.

## 0.150.7 — Lessons-Capture Honors a Prime's Inline-Only Marker (2026-07-29)

The pipeline's Lessons-capture step now inlines a lesson into a prime file that declares itself inline-only, instead of appending an archive link that would be dead in every consumer install.

- Both write sites (`actions/work.md` Step 8, `actions/review-work.md` standalone twin) branch on the prime's `## Lessons` marker comment ("inlined, not linked").
- Keyed off the marker condition in the prime's header, never a hand-list of primes (Closed Enumerations rule).
- The normal (non-marked) link path is unchanged.

## 0.150.6 — Forensics Check 4 Keys UR Closure on user_request Scan (2026-07-29)

`forensics` Check 4 (Orphaned URs) was keying UR closure on the capture-time `requests:` array — the same stale-list bug REQ-048 already fixed in `cleanup` Pass 1. That let it false-positive on UR-007 today, warning to archive a UR that still has six pending follow-up REQs.

- Check 4 now derives UR membership by scanning `user_request` across `queue/`, `working/`, `archive/` root, and `archive/UR-NNN/`, gating on `work-reference.md`'s Terminal-resolved status set by pointer (mirroring the REQ-048 fix), instead of testing whether the `requests:` array's ids all live under `archive/`.
- The `requests:` array is no longer read as the closure predicate; the live UR-007 false positive is gone.

## 0.150.5 — Route A Keeps Its Capture-Seeded write_set (2026-07-29)

Small doc-accuracy fix found during REQ-045: `capture-reference.md` said the pipeline's Scope step "firms up and overwrites" a REQ's `write_set` — true only for Routes B and C, since a Route A REQ never runs that step.

- Notes that a Route A REQ keeps its capture-seeded `write_set` for the whole run, and that value is what `work.md` Step 3 re-validates for disjointness when co-dispatched (per REQ-045).

## 0.150.4 — Board User Guide (2026-07-29)

The board's features were only documented in the agent-facing action file — no linkable human tour existed. New `docs/board-guide.md` covers what a user actually sees: modes, columns, badges, and the Testing view.

- Covers serve/static/summary modes, the four board columns, the Notes and Completion-anomalies strips, the toolbar, the card drawer, and the Testing view's columns and per-card actions.
- The `overlaps` badge gets its own subsection on the four ways it can under-report (no `write_set` declared, `*`/`**` glob quirks, identical malformed patterns, directory entries never badging files inside them) so it isn't misread as a safety guarantee.
- Keeps the human-tour/agent-contract boundary: feature facts sourced from `actions/board.md`, but none of its build/dispatch internals are duplicated here.

## 0.150.3 — Display-Only Overlap-Annotation Invariant Ratchet (2026-07-29)

The board's write-set overlap badge is display-only by design — it must never affect column placement — but that invariant was only protected by one Go test plus prose. `contract-regressions.sh` now pins it on both sides.

- Ratchet asserts `annotateWriteSetOverlap` runs after `bucketColumns` in `model.go`'s `buildBoard`, and that `board.md`'s Rules block keeps its "display-only, never column logic" wording.
- Anchors are call-site and heading-scoped (not file-wide greps), so a hoisted call or a relocated doc claim fails loud with a message naming the file and the fix.
- Red-green verified across five mutation scenarios in a sandboxed clone; live tree stays green.

## 0.150.2 — Board Badge Render-Path Test Assertion (2026-07-29)

The overlap badge's frontend render path had zero test coverage — only the Go-side annotation logic was tested. `generate_test.go` now asserts the badge's render tokens actually make it into the generated board HTML, so a regression in that path fails loudly instead of shipping silent.

- New test anchors on the inlined `web/board.js`/`web/board.css` tokens (`badge-write-overlap`, `writeSetOverlaps`, the drawer row, the CSS rule) rather than rendered DOM, so it holds regardless of live queue contents.
- Red-green verified: mutating the anchor token fails the test; restoring it passes.

## 0.150.1 — Doc-Accuracy Fixes: Legacy-Suppression Comment and Board Glob Miss-Classes (2026-07-29)

Three documentation inaccuracies from the deep review, all comment/prose with no behavior change.

- The memory hook's legacy-suppression comment claimed it "self-clears as soon as the next capture is written" — it actually suppresses to end-of-file and clears at the next UTC day's fresh log; the comment now matches the awk.
- The board's "a malformed glob pattern matches nothing" overstated: `writeSetPatternsIntersect` short-circuits on literal equality first, so two REQs declaring the identical malformed pattern still badge. Aligned across `model.go` (source), `board.md`, `work-reference.md`, and the board prime.
- Added the directory-entry case (`actions/` never badges `actions/board.md`) to the board's illustrative miss-class list, kept explicitly illustrative.

## 0.150.0 — Review Restatement Sweep (2026-07-29)

Six concurrency-spec REQs passed adversarial review at 86 to 98 percent, yet a later pass found every top defect was the same class: a contract changed in its canonical home while a restatement elsewhere kept the old meaning. The review step now forces a sweep for exactly that.

- `review-work.md` Step 6 gains a required **Restatement Sweep**: when a diff redefines something other text restates (a token, a field's semantics, what a hash holds, a command's output shape), grep every other statement/consumer and flag stale ones as findings — including in files outside the REQ's Scope (routed to follow-ups, not scored as builder scope drift). `work.md` Step 7 cross-references it.
- Trigger is condition-based, not a token list (Closed Enumerations Go Stale); a proportionality guard skips diffs that redefine nothing. Inherited by both pipeline and standalone `do-work review`.

## 0.149.3 — Cleanup Keys UR Closure on user_request, Not the Stale requests Array (2026-07-29)

`cleanup` Pass 1 decided a UR was done by reading its capture-time `requests:` array — but review-spawned and addendum follow-ups carry `user_request:` without ever being added to that array, so a UR with pending follow-ups could be archived out from under them. Pass 1 now uses the same predicate `work` Step 8 does.

- `cleanup` Pass 1 derives UR membership by scanning `user_request:` across `queue/`, `working/`, `archive/` root, and `archive/UR-NNN/`, gating on the terminal-resolved set (with `failed` holding the UR open); the `requests:` array is now a report-only cross-check.
- `capture` documents `requests:` as the capture-time record only — never the closure predicate. Two more readers of the old predicate were found and queued as follow-ups (forensics Check 4; a `failed`-status contradiction in work Step 8).

## 0.149.2 — Lock Mutex Re-Verifies Ownership Before Publishing (2026-07-29)

The serialized-lock mutex could evict a slow-but-live owner on the one-minute age check, and that owner's already-staged lock write would still land — clobbering its successor and losing a claim, the exact failure the mutex exists to prevent.

- The prescribed block now re-checks the mutex owner token immediately before the publishing `mv`; on a mismatch it discards the staged temp file, writes nothing, and re-acquires (exit 3) — narrowing the lost-update window from model-round-trip scale to the instant before the rename.
- The mtime-reclaim comment now says the age check proves age, not death, and points at the re-check as what makes eviction safe. The one-minute bound and the fixed-mtime property are unchanged; serial and single-session runs behave identically.

## 0.149.1 — Worktree Cleanup Uses the Recorded Operative Name (2026-07-29)

REQ-038 taught a crash-recovered worktree REQ to re-dispatch under a fresh name variant, but every later step still re-derived the original slug-based name — so after a collision the merge and cleanup targeted the *leftover*, not the builder's actual worktree.

- The name `git worktree add` actually succeeded with is now the REQ's held **operative name**, used by the hand-back merge, Step 8 cleanup, and the own-session crash-sweep exemption — never re-derived from the slug. (The merge site was the sharper bug: it would have integrated the wrong branch.)
- One variant scheme (incrementing numeric token) replaces the free counter-or-timestamp choice; the `worktree-agent-REQ-NNN-` prefix invariant is preserved. No-collision and serial behavior are unchanged.

## 0.149.0 — Dispatch Re-Validation: Full Route Coverage and the Serialization Loser (2026-07-29)

REQ-036's write-set re-validation was written against the Route B/C pipeline only, leaving a co-dispatched Route A REQ building under an unvalidated hint. This states one covering invariant and fixes three coherence gaps around it.

- Every co-dispatched REQ now gets exactly one post-dispatch re-validation, and its route picks which: Routes B/C at Step 5.5, Route A at Step 3 (serialize-only — it has no `## Scope` to hold a partition). The three previously-contradicting sentences now agree.
- The serialization "loser" is defined (the REQ at the re-check is held, never a dispatched sibling mid-build), with a deadlock guard for the two-discoverer case; dispatch-time partitions are persisted into `write_set` so a sibling's re-check compares against the real subset.
- The absent-`write_set` gloss is reworded condition-first and reconciled with REQ-044's conditional recovery clear; a new contract-regression ratchet pins the whole contract.

## 0.148.0 — Lock Claim Coherence: Dispatch Record as the Recompute Source (2026-07-29)

Four coherence defects survived REQ-035's move to a canonical `claimed_reqs` list. The fix names a single source of truth for the recompute — the orchestrator's in-memory dispatch record — instead of a `working/` listing.

- The heartbeat recompute reads the session's dispatch record, never a directory listing and never the lock's own previous `claimed_reqs`. Step 2's claim-before-move append enters the id into that record before the file moves, so the refresh carrying the claim no longer erases it.
- A known-dead builder's id leaves the record (so its REQ is reclaimable that session), explicitly distinguished from an ordinary failed build that keeps its claim through remediation and Step 8.
- Crash Recovery clears `write_set` only for Scope-mirrored sets and preserves capture-seeded / Route-A sets; a stale proceed-anyway gate restatement is aligned and pinned by a new block-scoped contract-regression ratchet.

## 0.147.0 — Worktree Merge Range: Fail-Loud Validation and Seam-in-Range Merge (2026-07-29)

The worktree-dispatch merge range `<pre>..<merge_hash>` had six confirmed defects in how it was produced, validated, and restated — most dangerously, `qualify.sh` printed OK on a broken range instead of failing. All six are fixed as one coherence contract.

- `qualify.sh` now hard-FAILs (exit 1, naming the range) on an unresolvable `DO_WORK_DIFF_RANGE` instead of reading an empty diff and passing vacuously; serial mode is byte-unchanged.
- The integration seam is folded into the merge commit (`git merge --no-ff --no-commit` → apply seam → commit), so it provably lands inside the range and `commit:` still records the merge.
- Remediation re-merges get a defined cumulative range (`<pre1>..<merge_hash2>`); Step 6 gains an imperative orchestrator-side hand-back merge instruction; the hash-writeback block gains a worktree carve-out so it can't record the changelog commit's hash.

## 0.146.3 — Board Overlap Badge Uses OS-Independent Glob Matching (2026-07-29)

The Kanban board's write-set overlap badge matched globs with `filepath.Match`, whose separator is `\` on Windows — so `*` could wrongly cross `/` and the badge would misjudge contention off-platform. It now uses `path.Match` (correct for slash-separated repo-relative paths), and the glob dialect is finally written down where readers meet the field.

- `writeSetPatternsIntersect` uses `path.Match`; the doc comment, `actions/board.md`, the board prime, and the `write_set` schema line all state the dialect: `*` never crosses `/`, `**` is not recursive, malformed patterns match nothing on the board — the dispatch gate still treats an unexpandable/overlapping glob as overlapping, so a board false-negative never loosens it.
- New tests pin the slash boundary (with a same-segment positive control) and malformed-pattern behavior; the badge stays display-only (no schema or column-placement change).

## 0.146.2 — Note Worktree Isolation in the Harness-Tier Guide (2026-07-29)

Someone sizing a harness against `background-agents.md`'s three fan-out rungs had no way to learn that per-builder git-worktree isolation exists or where it's documented. A short cross-reference now closes that gap.

- Added a "worktree isolation is a separate axis" note after the harness rungs, pointing at `actions/work-reference.md` → Worktree Dispatch Mode; the rungs and the file's load contract are unchanged.

## 0.146.1 — Worktree Name-Collision Handling on Re-Dispatch (2026-07-29)

Worktree dispatch mode names each builder's worktree and branch deterministically from the REQ id, and the crash sweep reports (never deletes) an unmerged leftover — so a crash-recovered REQ would re-dispatch straight into the name its own leftover still holds and fail to start. Re-dispatch now sidesteps the occupied name instead of deadlocking on it.

- On a name collision at creation, dispatch under a fresh unique variant (an incrementing `-2`/`-3` or short timestamp token), keeping the `worktree-agent-REQ-NNN-` prefix so sweeps still correlate both names to the REQ.
- The crash sweep now states that a reported unmerged leftover doesn't block re-dispatch — the collision variant covers it, and the two coexist until cleanup Pass 5 resolves the leftover.

## 0.146.0 — Worktree Merge Placement and Evidence Re-Pointing (2026-07-29)

Worktree dispatch mode said who merges and how, but never *when* in the pipeline — and after a merge the main tree is clean, so the qualify check, the review step, and the commit step all read an empty diff and quietly passed nothing. The merge now has a fixed place in the sequence and a defined range, and every evidence step reads that range instead of the post-merge-clean tree.

- The orchestrator merges each builder branch at hand-back (end of Step 6, before the Implementation Summary), and captures the range `<pre>..<merge_hash>` around it — stated once in Worktree Dispatch Mode and consumed by qualify (Step 6.3), review (Step 7), post-merge verification (Step 8), and Step 9's validation.
- `tools/checks/qualify.sh` gained an optional `DO_WORK_DIFF_RANGE` env var; unset, it reads the working+staged diff exactly as before, so serial runs are byte-for-byte unchanged.
- Step 9 is reconciled for merged work: it stages only the changelog/version/metadata (the implementation is already in the merge commit) and records the merge commit's hash.
- Fixed a latent gap where `work.md` and `work-reference.md` disagreed on the post-merge verification default (now both say per-merge whenever more than one REQ is in flight).

## 0.145.0 — Re-Validate Write-Set Disjointness When Scope Firms It (2026-07-29)

The parallel-dispatch gate decided co-dispatch on capture's write-set guess, but Step 5.5 then rewrote that field from each REQ's real scope with no second look — so two REQs seeded as disjoint could both quietly claim the same file once their scopes firmed. Step 5.5 now re-checks disjointness before it commits the field, and a dispatch-time partition directive survives the mirror instead of being erased.

- New Step 5.5 re-validation runs only under co-dispatch: it re-checks the firmed scope against every other in-flight REQ's current `write_set` and serializes or partitions the loser before its builder starts — the same check Step 6 already runs before a mid-build write-set extension.
- The Step 1 gate and Step 5.5 now agree that both steps enforce, at different times (the gate on capture's hint, Step 5.5 after firming); Step 4's plan-validation flag is documented as a warning, not the enforcement point.
- The Step 6 write-boundary bullet clarifies that an absent `write_set` handed to a builder means "dispatched serially, full-scope freedom," never "write nothing."
- Serial (floor) runs are completely unaffected — every new clause is gated to the parallel-dispatch path.

## 0.144.0 — Concurrent Claim Tracking in the Orchestrator Lock (2026-07-29)

The orchestrator lock could only name one in-flight REQ per session, so the moment a single orchestrator dispatched more than one builder at once, Crash Recovery and cleanup mistook the siblings for abandoned crash artifacts and re-queued them mid-build. The lock now tracks every concurrent claim, and the recovery and cleanup gates honor the whole set — the session's own claims included — so parallel dispatch is finally safe inside the skill's own protocol.

- New canonical `claimed_reqs` list on the holder and each coexisting-session entry; the old `claimed_req` stays as a derived legacy mirror (`claimed_reqs[0]`), so older readers and the serial default are completely untouched.
- Crash Recovery now gates on freshness alone and skips any file in a live claim set — the session's own included — so a Step 10 → Step 1 loop no longer strips and re-queues a still-building sibling.
- Per-merge post-merge verification becomes the default whenever more than one REQ is in flight, and `cleanup.md`'s Pass 0 live-claim gate reads the whole claim list.
- Contract-regression ratchets pin both the field's presence and the same-story gate phrasing across `work.md`, `work-reference.md`, and `cleanup.md`.

## 0.143.0 — Capture Slicing Nudge and Board Write-Set Overlap Badges (2026-07-29)

Two upstream levers for parallel-friendly queues: capture's slicing convention now prefers boundaries that give each REQ its own files (declaring unavoidable overlap in `write_set`), and the Kanban board shows an `overlaps` badge on pending/claimed cards whose declared write-sets could touch the same files.

- The overlap annotation computes in Go after column bucketing — structurally display-only (badge + drawer rows, never column logic; co-dispatch decisions stay with the work pipeline's gate) — and is glob-aware in both directions.
- Drawer gains "Write set" and linked "Overlapping write sets" rows; no badge on a REQ without a declared `write_set` (unknown, not safe — the gate's serialize reading is documented alongside).
- The three stale "no overlap computation on the board" claims (board action, maintainer doc, board prime) updated in lock-step.

## 0.142.0 — Worktree Dispatch Mode with Defined Cleanup Ownership (2026-07-29)

The work pipeline now documents running builders in orchestrator-created git worktrees: each builder commits on its own `worktree-agent-REQ-NNN-*` branch, the orchestrator stays the sole writer of the main tree and merges in dependency order, and nothing archives until the merged state re-passes the REQ's checks. Every leftover now has an owner.

- Happy path: the archive step removes the worktree and deletes the branch with `git branch -d` from the integration branch — a refusal is the signal a merge was skipped or lost, so never `-D`/`--force`.
- Crash path: Crash Recovery sweeps `worktree-agent-*` leftovers — merged ones removed mechanically, unmerged ones only reported; discarding unmerged work belongs to cleanup's new consent-gated Pass 5 (its first interactive pass — six passes now, mirrored in the cleanup guide).
- `do-work/` state stays in the main tree only; builders get their brief in the dispatch prompt (and treat any committed `do-work/` snapshot in a worktree as absent).
- Ships honestly single-builder: co-dispatching several worktree builders waits on the lock's multi-claim work (queued follow-up), and four contract-regression ratchets pin the naming, the `-d` assertion, the post-merge gate, and the consent-gated pass.

## 0.141.0 — Write-Set Declarations and a Parallel-Dispatch Gate (2026-07-29)

The queue schema gains an optional `write_set` field (repo-relative paths/globs a REQ expects to write; absent means it overlaps everything), and the work pipeline gains an opt-in dispatch gate: advanced harnesses may co-dispatch dependency-ready REQs whose write-sets are pairwise disjoint. The serial default is untouched.

- Serial-only resource classes: REQs writing ordered/generated resources (migrations, lockfiles, generated bundles — illustrative list) never co-dispatch, disjoint or not.
- Builders get a write boundary: out-of-set needs are a stop-and-report to the orchestrator, never a silent write; the Scope declaration one-directionally mirrors into `write_set` so the two can't drift.
- Chosen over timed per-file locks — a TTL expires over a live slow agent and hands the file to a second writer (the 0.140.4 mutex-break defect class).
- `tools/queue-kanban` parses `write_set` into the board payload (display only), and three new contract-regression ratchets pin the gate text and the parser lock-step.

## 0.140.4 — Owner-Checked Lock Mutex, Atomic Capture Appends (2026-07-28)

Two accepted findings from an external concurrency review. The lock mutex could be forcibly broken after 15 seconds even when its owner was live mid-write — dangerous now that the critical section legitimately spans a model round-trip — and concurrent Stop-hook captures could interleave their writes in the shared daily log.

- `actions/work-reference.md`: removed the 15-second attempt-count mutex break — the one-minute mtime check (a verified stale-owner bound) is now the only reclaim path. The winner records an owner token in the mutex so release can't delete a successor's mutex, and an `mkdir` failure with no contender present reports and stops instead of spinning forever.
- `hooks/memory-stop-capture.sh` (+ spec in `actions/memory-reference.md`): each capture section is composed first and appended in a single `printf` — one atomic `O_APPEND` write — so near-simultaneous stops from sessions sharing the daily log can no longer garble section structure. No lock: `flock` doesn't exist on macOS and the hook must never block session end.

## 0.140.3 — Loopback-Only Board Writes, In-Mutex Stale Revalidation, Redact-Before-Truncate (2026-07-28)

Five accepted findings from an external review, all verified against the code before fixing. The two that mattered most: a LAN-exposed board's testing endpoints accepted writes from any machine (the Origin check only fires when a browser sends one), and a stale-lock takeover judged staleness on a pre-mutex read, so a holder that heartbeated in the gap could be overwritten and its in-flight REQ re-queued.

- Kanban testing writes (`/api/testing/profile`, `/api/testing/status`) now require a loopback peer, same as `/file` — a network-exposed board is read-only.
- The stale-lock takeover re-confirms the holder's identity and recomputes its heartbeat age from the fresh read inside the mutex before overwriting; the user-gated take-over keeps the identity check, and coexisting-session prune ages come from the same fresh read. Wording now ratcheted by `_dev/tests/contract-regressions.sh`.
- The memory Stop hook redacts credentials (and judges the private-key drop) on the full extracted messages *before* truncation — a byte-budget cut can no longer sever a token into an unmatched, persisted fragment. Spec in `actions/memory-reference.md` reordered to match; regression probe reproduces the straddling-token case.
- Stray (misplaced) REQ files now feed the live board's mtime fingerprint, so their warning appears and clears without waiting for an unrelated file change.
- Testing updates preserve the REQ file's existing permission bits across the atomic rewrite instead of forcing 0644.

The remaining four accepted findings from the same external triage 0.140.1 started on. The big one: a live Route C build could go 45+ minutes between lock touchpoints and get taken over — and re-queued — mid-build.

- Heartbeat refreshes are no longer phase-boundaries-only: refresh before dispatching and when each long-running agent returns (explore, plan, build, review), plus at any pause once 15 minutes have passed since the last lock write (`actions/work-reference.md`, `actions/work.md` Steps 4–6). The 45-minute stale threshold's rationale now states the schedule it depends on.
- Cleanup's Pass 0 gained a live-claim gate: any REQ freshly claimed in the orchestrator lock by another session is exempt from the queue and `working/` sweeps — a coexisting session flips its REQ terminal *before* moving it, and that window used to look like an abandoned sweep candidate (`actions/cleanup.md`).
- Next-step suggestions are sourced from SKILL.md's routing triggers, not Action Dispatch names — `do-work capture-requests` didn't route; the capture form is `do-work capture-request: <text>` (`next-steps.md`).
- Lock acquisition and pipeline-state setup now run the `git ls-files` tracked-path check behind their ignore appends, reporting the `git rm --cached` remedy when an earlier session committed the file — an ignore rule can't rescue an already-indexed path (`actions/work-reference.md`, `actions/pipeline.md`; the memory-module installer already had this check since 0.139.2).

## 0.140.1 — Claim-Before-Move Closes the Lock Races; Checkpoints and Temp Files Stop Lingering (2026-07-28)

Four accepted findings from external feedback triage, all in the concurrent-orchestrator machinery. The claim/recovery race and the acquisition race could each let two sessions fight over one REQ; the other two were slow leaks.

- Step 2 now claims the REQ in the orchestrator lock *before* moving it into `working/`, and the Crash Recovery gate lists `working/` *before* reading the lock — the paired ordering means a file can never be observed unclaimed by a live scan (`actions/work.md`, `actions/work-reference.md`). Archive-time gets the explicit mirror order: move out first, then clear `claimed_req`.
- Lock acquisition re-validates existence inside the serialized mutex — a session that loses the empty-queue race now falls into the existing-lock decision tree instead of silently overwriting the winner's holder slot (`actions/work-reference.md`).
- The local-ignore prescription for the lock now uses the glob `do-work/orchestrator-lock.json.*`, covering orphaned PID-suffixed `.tmp` files as well as the mutex directory.
- Checkpoint deletion is scoped to the files a session is allowed to recover, so a coexisting session's live claim in `working/` no longer keeps a stale `CHECKPOINT.md` alive forever (`actions/work.md`).

## 0.140.0 — UTF-8-Safe Capture Truncation; Forget Now Scrubs the Daily Logs Too (2026-07-28)

Two accepted findings from external feedback triage. The stop-capture hook could tear a multi-byte character in half, and "forget" only forgot half the store.

- The stop-capture hook's byte-budget cuts now pipe through `iconv -c -f UTF-8 -t UTF-8` (plain cut when iconv is absent). A raw `head -c` cut lands mid-character routinely on CJK text (~3 bytes/char), persisting invalid bytes into the log and the dedup hash — and on macOS the torn sequence made BSD sed fail in the redaction pipeline, silently dropping the whole capture.
- `memory forget` is now an explicit, confirmation-gated sub-command (with a `do-work forget` alias) instead of one clause inside `remember`. It removes the working-memory bullet AND redacts matching daily-log lines in place with a `[forgotten — …]` marker, since recall searches the logs too. It is the one named exception to the logs-are-append-only rule, scoped to explicit user invocation — automatic writers still only append, capture-body lines keep their `> ` quoting, and heading lines (the dedup key) are never touched.

## 0.139.4 — Capture Boundaries Can't Be Spoofed; Tracked Logs Caught Before the Install Fast Path (2026-07-28)

Four follow-up fixes, all found by reviewing the previous two rounds' fixes. Two of them were holes the earlier fixes left open rather than new bugs.

- Capture sections now open with a sentinel and quote every body line, and the session-start filter ends a section only at an unquoted heading. Heading grammar alone was still spoofable: raw capture text can contain `## 12:34 UTC note`, which ended the section and injected the rest. Legacy sections written before this release have unquoted bodies where no boundary is trustworthy, so they're suppressed to end-of-file — that hides curated entries written after them that day, and self-clears with the next capture.
- The memory-module install detects a tracked raw store in Phase 1 instead of Phase 4. A fully-wired install with a committed log plus an ignore rule took the "already installed" early return, so the check added in 0.139.2 never ran in the one scenario it was written for.
- A `memory` payload that arrives with no sub-command now falls back to `recall`, never `remember`. Real recall queries are usually noun phrases, so the previous sentence-shape test would have classified most reads as writes and mutated the store on a request to read it.
- `do-work recall` with an empty query — including the `what do you remember` phrasing — is now defined: it presents working memory plus recent curated log entries instead of searching for nothing.

## 0.139.3 — Orchestrator Lock Updates Are Serialized; Long Prompts No Longer Eat the Agent Reply (2026-07-28)

Three fixes from a second review pass. The lock guard turns out to have had a lost-update hole since coexisting sessions were added — measured at 19 of 20 concurrent writes discarded.

- Every write to `do-work/orchestrator-lock.json` now runs inside a `mkdir` mutex and lands via temp-file-plus-rename. Multiple sessions write that file (a coexisting session refreshing its entry, the holder's prune, a take-over), so "each writer touches only its own fields" was never enough: a lost `claimed_req` makes Crash Recovery re-queue a REQ that's actively being built. The mutex self-heals after a minute and never blocks the pipeline.
- The Stop hook budgets the user and assistant sides of a capture separately. Truncating the combined string meant a long prompt silently dropped the entire assistant reply — the half holding the decisions. Each side is guaranteed half the budget and a short side yields its slack; cuts are marked `[truncated]`. This was reachable only because 0.139.2 started populating the user side.
- `do-work run` holds the orchestrator lock through `cleanup` and its commit instead of releasing first. Cleanup sweeps `do-work/queue/` and `do-work/working/`, so releasing early let a departing session sweep and commit an arriving session's just-claimed REQ.

## 0.139.2 — Raw Session Captures Stay Out of Session Start; Tool-Using Sessions Capture the Human Prompt (2026-07-28)

Four review fixes to the memory engine, its installer, and the router. The important one: a Markdown heading inside a captured exchange could end the capture section early, so the rest of the raw transcript was injected as curated memory at the next session start — exactly what that filter exists to prevent.

- Capture bodies are now blockquoted at write time and the session-start filter only treats `## HH:MM UTC …` as a section boundary, so a captured `## Findings` can't end the section. The reader-side rule is what protects logs already on disk.
- The Stop hook now finds the real final exchange. Claude Code records tool results as `type: "user"` entries, so any session whose last turn used a tool was storing an empty `User:` side — 7 of 8 recent transcripts on the author's machine. Text is pulled only from `text` blocks, `isMeta` entries are skipped, and blank entries are dropped; the same fix repairs assistant turns that end in a tool call.
- `do-work remember X` and `do-work recall X` keep their verb as the sub-command instead of arriving as a bare payload, and `what do you remember` maps to `recall`.
- The memory-module install verifies with `git ls-files`, not just `git check-ignore`. Ignore rules don't apply to tracked files, so a repair install over a committed log reported "logs ignored: OK" while the Stop hook kept appending to it. The canonical local-ignore snippet in `crew-members/background-agents.md` now carries the same caveat so the next copy-paste doesn't reintroduce it.

## 0.139.1 — Bootstrap Sentinel Stays Machine-Local; "Remember To" Routes to Capture (2026-07-27)

Three review fixes to the memory engine and the installer. The big one: the bootstrap sentinel was committable, so one machine's import would have silently blocked every other clone from importing its own history.

- `do-work install memory-module` now adds `memory/.bootstrap-imported` to `.git/info/exclude` alongside `memory/logs/` and the usage ledger, and verifies it. `memory bootstrap` refuses to re-run when the sentinel exists — committed, that refusal would have followed the repo to every other machine.
- Routing: `remember to fix X` is queued work, not a fact. SKILL.md row 37 now sends task-shaped `remember` phrasings to capture; `actions/memory.md` already documented the boundary, but the router decides the route before the action file is ever read.
- The `bowser` skill download's cleanup returns failure (`|| { rm -f …; false; }`) instead of reporting a failed download as a success — matching the `ui-design` and `ideation-adhd` install rows.

## 0.139.0 — Parallel Memory Engine: memory module, install target, and value auditor (2026-07-25)

A second, capture-first memory engine now runs alongside bkb so real usage data — not theory — decides which one earns its keep (ADR-017). Both engines log usage, and a new auditor renders the head-to-head verdict.

- New `do-work memory` action (`remember` / `recall` / `status` / `bootstrap` / `audit`): a 2,500-char capped `memory/working-memory.md`, dated daily logs, and layered recall (lexical always; semantic only when an embedding backend is detected) with cited sources. Companion schemas in `actions/memory-reference.md`.
- New `do-work install memory-module` target scaffolds `memory/` and merges optional SessionStart/Stop hooks into `.claude/settings.json` — composing with existing hook entries (backup + parse-verify), never clobbering them. The Stop hook appends a hash-deduplicated capture of each session's final exchange and never blocks a session end.
- New `do-work memory audit` (`actions/memory-value.md`): read-only, engine-agnostic value audit of bkb and the memory engine — git/log history probes, usage-ledger stats, hit-cited rate as the verdict signal, with an explicit fairness rule for bkb's pre-instrumentation era.
- bkb `query`/`ingest` now append best-effort usage-ledger events so the comparison sees both engines.
- Review hardening (PR #122): the Stop hook redacts credential-shaped text before persisting captures; session start injects only curated memory (raw captures stay behind `memory recall`, which loads the prompt-injection guardrail); the hook merge checks and appends each settings entry independently; the installer repairs any missing scaffold component; and the auditor locates the KB via bkb's locating contract instead of assuming `kb/`.
- Raw captures and the usage ledger are machine-local: the installer adds `memory/logs/` and `memory/usage-ledger.jsonl` to the repo's `.git/info/exclude` (never your committable `.gitignore`), so only the curated `working-memory.md` is shareable. Redaction stays as a second line of defense.

## 0.138.3 — Skill Downloads Are Atomic: Temp File, Then Rename (2026-07-25)

The 0.138.1 fix caught zero-byte downloads, but `curl -o` writes the final path incrementally — a connection dropped mid-transfer left a non-empty partial `SKILL.md` that `test -s` read as a complete install, unrepairable by re-running. Downloads now land in a `SKILL.md.download` temp name and only rename into place on success.

- Applies to all three curl-based skill installs (`ui-design`, `ideation-adhd`, `bowser`); the temp file is removed on failure, so nothing is left behind either way.
- Chosen over `curl --remove-on-error`, which needs curl ≥ 7.83 — the rename works on any curl.

## 0.138.2 — Install Target Renamed to ideation-adhd (2026-07-25)

The 0.138.0 target ships as `install ideation-adhd` — the name now says what it does (the "adhd" is the upstream skill's metaphor for its branching style, not the substance). `install adhd` and the `adhd-mode` spellings still work as aliases.

- The install **folder** stays `.claude/skills/adhd/` — it must match the upstream frontmatter `name:` field so `/adhd` auto-discovers.

## 0.138.1 — Install Detect Treats a Zero-Byte Skill File as Absent (2026-07-25)

An interrupted download could leave a zero-byte `SKILL.md` that the `ls`-based detect read as "already installed", making the failed install unrepairable by re-running. Review caught it on the new `adhd` target; the same copy-pasted primitive was fixed in `ui-design` and `bowser` too.

- Detect commands for the single-file skill targets now use `test -s` (non-empty), so a re-run over a failed download repairs it instead of stopping at Phase 1.
- The never-overwrite rule is scoped to non-empty files: reinstalling over a zero-byte file is repair, not overwrite.

## 0.138.0 — New Install Target: adhd Divergent-Ideation Skill (2026-07-25)

`do-work install adhd` vendors the [adhd skill](https://github.com/UditAkhourii/adhd) (MIT) into the project — parallel divergent ideation that branches a named problem across distinct cognitive frames, then scores, clusters, and deepens the top candidates. Complements `scan-ideas` (repo-grounded) with deliberately unconventional exploration; feed the winners to `capture-request:`.

- Single self-contained `SKILL.md` installed project-scoped to `.claude/skills/adhd/` — folder name matches upstream so `/adhd` auto-discovers; no global npm install.
- Same manifest-driven detect → install → verify → report shape as `ui-design`; idempotent, never overwrites an existing copy.
- Routing accepts `install adhd` (also the `install adhd-mode` / `install adhd mode` / `setup adhd` spellings — the target normalizes after the install verb; bare `adhd` without the verb is not a route).

## 0.137.0 — Clarify Opens Each Question With a Plain-Language Story (2026-07-24)

Answering pending questions used to mean remembering what REQ-025 was about — often days after the work happened. Now every question arrives with a short story above it: what you asked for, what the builder ran into, and why the call is yours. The decision block underneath is unchanged.

- Questions are presented in three layers — a 1–4 sentence story, then the existing `Decision / Value / Risk / Also` block, then the builder's original wording and file paths *only if you ask*.
- Layer one is written to be read aloud: no file paths, no bare identifiers, no CamelCase, one idea per sentence. Any technical term used lower down gets paraphrased in the story first.
- Blocked REQs waiting on an external condition now get a one-line "what it was for" too — those are the ones you've had the longest to forget.
- New red flags catch the failure mode this invites: a story that just restates the question is padding, not context.

## 0.136.1 — Board Flags REQ Files Found Outside the Scanned Sections (2026-07-24)

A REQ that lands somewhere other than `queue/`, `working/`, or `archive/` — say a work agent that archived to `do-work/user-requests/UR-NNN/` instead of `do-work/archive/` — used to vanish from the board with no trace. Now the walk catches it and raises a data warning instead of silently dropping it.

- The board now emits a warning naming the misplaced REQ, its exact path, and how to fix it (move into `archive/` or `queue/`) — shown in the web warnings banner and the `board summary` output.
- A stray REQ is only flagged, never rendered as a card, so its off-vocab location can't masquerade as a real column placement.

## 0.136.0 — Maintainer Docs No Longer Ship to Consumer Installs (2026-07-23)

The repo's own `CLAUDE.md` and `AGENTS.md` were landing in every consumer install, where Claude Code auto-loads the nested `CLAUDE.md` on every skill-file read — a ~2.5k-word context tax whose commit protocol (bump the version, add a changelog entry) is actively wrong advice inside someone else's project. They're maintainer docs, not skill content, so they no longer ship.

- Both files are now `export-ignore`'d, and `do-work update` deletes the stale copies that installs ≤0.135.x left behind (tar extraction never removes files dropped upstream).
- Every shipped file's citation of the skill's own CLAUDE.md was reworded to be self-contained or point at a shipped home (e.g. `actions/kb-lessons-handoff.md` for the KB handoff contract) — 14 sites across actions, crew-members, and hooks.
- New contract-regression checks keep it that way: the export-ignore lines must exist, and shipped files must not cite the unshipped docs.

## 0.135.0 — Board Drawer Links Every REQ/UR Mention, URL, and File Path (2026-07-23)

The detail drawer's cross-references are now real, obviously-styled links instead of plain text or button-shaped chips. File paths get existence-checked at build time, so a stale reference is visible at a glance.

- Every REQ/UR id in the drawer is a link: the UR drawer's "REQ ids" row, the REQ drawer's "User request" / "Depends on" / "Unblocks" / "Blocked by" rows, and any `REQ-…`/`UR-…` mention inside a rendered body (only ids actually on the board — unknown mentions stay text). Short mentions resolve compound card ids (`REQ-031` → `UR-002-REQ-031`).
- All links are visibly links: accent color + underline (the old "User request" chip looked like a badge).
- File paths in code spans are checked against the repo at board-build time: existing files render as blue links that open read-only via the live server's new `GET /file` endpoint (loopback-only, repo-contained, always text/plain); missing files render red with a "Not found in this repository" tooltip — in static snapshots too, where the existence verdict is baked into the data.
- URLs in code spans become clickable, and every http(s) link in a body opens in a new tab instead of navigating the board away.

## 0.134.0 — Pending-Card Timer Tracks the Last Transition, Not Capture Time (2026-07-23)

A REQ answered via `do-work clarify` went back to `pending` but its board card kept counting from capture time ("queued … · 28m" seconds after the flip). Pending-tier cards now time from the last state change.

- New optional `status_changed_at` frontmatter field, stamped on any status flip that has no dedicated `*_at` stamp of its own — clarify's answered/unblock/discovered-task flips, work Step 1's probe unblock, reserve release, manual resets. Schema + board parser updated in lock-step.
- The board's state timer for pending-tier cards resolves: `status_changed_at` → the later of `created_at` / file mtime (only when mtime beats capture by >5min, so untouched cards still read "queued") — verb "updated" whenever the instant isn't capture time.
- Dedicated stamps stay authoritative: mtime never drives the claim/blocked/reserved stopwatches (the pipeline edits the file all through a claim), and mtime remains banned from completion dating.
- `status_changed_at` is covered by the 0.133.0 future-stamp guard.

## 0.133.0 — UTC Stamping Rule + Board Defense Against Future-Dated Timestamps (2026-07-22)

A session stamped `claimed_at` with local wall-clock time plus a `Z` suffix, and the board's stopwatch froze at "0s" until the wall clock caught up. Now the instructions say exactly how to stamp, and the board calls out bad stamps instead of rendering them silently.

- New Timestamp rule in `actions/work-reference.md` (Full Frontmatter): every `*_at` field is the current UTC instant from `date -u +%Y-%m-%dT%H:%M:%SZ` — never local time with a `Z` suffix. Every write site that said just `<timestamp>` / `<now>` (work claim + terminal flips + blocked flip, capture templates, clarify, reserve, abandon) now says so or points at the rule.
- queue-kanban flags any frontmatter timestamp parsing later than board generation time + 2min skew: a "⚠ future stamp" card badge naming the field(s) plus a data warning (banner in serve/static, listed in summary) with the fix.
- The stopwatch renders an honest "⚠ clock skew" marker (with a tooltip explaining the likely cause) instead of a dead-looking "0s" when its instant is beyond the skew allowance, and recovers to normal ticking once the clock catches up.
- New forensics check 12 sweeps all REQ frontmatter for future-dated timestamps.

## 0.132.0 — Every Non-Terminal Card Gets a Live State Timer (2026-07-22)

The claimed-card stopwatch from 0.131.0 now covers every state: any card that isn't done shows when it entered its current state plus a ticking elapsed timer — "queued … · 3d 04h", "blocked … · 2h 15m", "reserved … · 12m 30s" — so you can see at a glance where time is going before a task lands in Done.

- Pending / pending-answers / failed cards count from `created_at` (labeled "queued" — time since capture; those states write no transition instant of their own), blocked cards from `blocked_at`, reserved from `reserved_at`, claimed from `claimed_at` as before.
- Durations past a day render as "3d 04h" instead of a wall of hours.
- Drawer parity: "Blocked since" and a new "Reserved" row tick while that hold is live, degrading to the plain instant when the field is a stale leftover.

---

## 0.131.1 — Kill Boards by Executable, Bound Probes Everywhere, Guard a Symlinked do-work (2026-07-21)

Four review follow-ups, each pinned by a new contract-regression check.

- `just run-kanban` (and the installed recipe) now identifies a stale board by the listener's **executable path** via `lsof -d txt`, not by grepping argv — a process that merely mentions "queue-kanban" in an argument can no longer be killed by mistake. No `lsof` info still means refuse-to-kill.
- The work pipeline's `blocked_check` probe keeps its 30-second bound even where GNU `timeout`/`gtimeout` don't exist (stock macOS): a background-and-poll fallback kills a hung probe and reports exit 124 instead of running unbounded.
- Roadmap's **Ready** bucket now requires normalized `status: pending`, so `blocked`/`reserved`/`claimed` REQs can't classify as Ready.
- The board's testing writes (REQ placeholders *and* `testers.md`) verify that `do-work/` itself resolves inside the repository — a symlinked `do-work/` pointing elsewhere is refused. If you deliberately symlink `do-work/` outside your repo, the board's Testing view will decline to write there (the work pipeline is unaffected).

## 0.131.0 — Claimed Cards Show Claim Time and a Ticking Duration (2026-07-21)

A claimed card used to sit in its column with no time at all. It now shows when it was claimed plus a live stopwatch — "claimed Jul 21, 19:44 UTC · 4m 01s" — so you can see at a glance how long the current REQ has been in flight.

- The duration ticks every second (s → "Xm YYs" → "Xh YYm") with tabular digits so the line doesn't jitter.
- The detail drawer gets a matching "Claimed" row: ticking stopwatch while the claim is live, plain instant + relative label if a stale `claimed_at` lingers on any other status.

## 0.130.4 — Testing View: Safer Writes, Sturdier Feedback Form (2026-07-21)

Hardening for the board's one write surface and the client flows that drive it.

- Server: all testing writes are serialized behind a mutex (concurrent add-tester or status posts raced their read-modify-write cycles); REQ frontmatter updates land via atomic temp-file-and-rename so a crash can never leave a zero-byte REQ; the testers file opens with O_APPEND.
- Client: typed feedback survives testing-view re-renders and failed posts (the form now closes only on server confirmation); Clear no longer requires a tester profile (it only removes fields); compound ids like UR-002-REQ-031 sort by their REQ number; the Recently-done and testing date windows anchor to the wall clock instead of page-generation time, so a long-open tab keeps meaning "last 24 hours".

## 0.130.3 — Board Names the Real Cause When git Is Unavailable (2026-07-21)

Running the board without a `git` binary (or outside a repo) used to produce a per-ticket anomaly blaming each commit hash. Now a one-time probe logs a single clear line ("git binary not found on PATH"), skips the doomed per-ticket subprocesses, and the anomaly reason says the hash could not be dated rather than asserting it is invalid.

## 0.130.2 — Blocked-REQ Probe Works on Stock macOS (2026-07-21)

Two fixes to the blocked-status machinery in `actions/work.md`.

- The `blocked_check` probe no longer assumes GNU `timeout` exists: it resolves `timeout` → `gtimeout` → unwrapped, so stock macOS (which ships neither) probes the condition instead of failing on exit 127 and wrongly reporting "probe failed".
- The mid-run blocked flip's "no edits landed" test now excludes `do-work/` from its porcelain/diff check — the REQ's own bookkeeping is always dirty mid-run, so the unscoped check could never read clean and silently defeated the flip.

## 0.130.1 — Check Scripts Stop Misreading Diffs and Scope Lists (2026-07-21)

Two bugfix rounds for the shipped verification scripts from the work-pipeline hardening.

- `tools/checks/qualify.sh`: `grep -q` on a piped `git diff` could SIGPIPE the pipeline and mark genuinely-changed files as absent (false WARNs); the diff file list is now computed once. The debug-artifact grep now excludes `do-work/` at the pathspec level, so REQ prose merely *mentioning* console.log/TODO no longer FAILs clean implementations.
- `tools/checks/scope-drift.sh`: inline `**Files I will touch:** \`a\`, \`b\`` lists now parse (previously only bullet lists did — an inline list silently turned the whole check into a SKIP); a touch-list header with zero parseable paths is now a FAIL instead of a silent SKIP; drift path lists print one per line unsplit.

## 0.130.0 — Relative Times Next to Every Board Timestamp (2026-07-21)

Every timestamp on the Kanban board now carries a live relative label — "done Jul 21, 16:24 UTC · 6min ago" — so you can tell at a glance how fresh a card is without doing UTC math.

- Covers card done/cancelled lines, the testing view's tester chips, the detail drawer's Created / Blocked since / Completed / Testing updated rows, the reserved/blocked badge tooltips, and the "Generated" page header.
- Labels tick every second client-side (s → min → h → d buckets), so a tab left open stays accurate; tooltips get a render-time snapshot.

## 0.129.1 — Kill-Stale Guard Matches Cross-Repo Kanban Binaries (2026-07-21)

`just run-kanban` can now restart the board when the port is held by a queue-kanban instance started from *another* repo. Other projects' recipes build the same tool under different names (e.g. `build/go-bin-queue-kanban`), and the old guard only killed a process named exactly `queue-kanban` — so the recipe refused and failed instead of reclaiming the port.

- The kill-stale check now substring-matches `queue-kanban` against the listener's full command line (`ps -o args=`, which unlike `comm` isn't truncated on Linux) and echoes what it stopped.
- Unrelated listeners are still left alone and named in the error, exactly as before.
- Applies to both the shipped `just-kanban` install template (`actions/install.md`) and this repo's own justfile; already-installed projects get the fix as a consent-gated upgrade offer on re-running `do-work install just-kanban`.

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
