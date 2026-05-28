# Changelog

What's new, what's better, what's different. Most recent stuff on top.

> Pre-0.65 release notes lived in `CHANGELOG-2026-spring.md` and `CHANGELOG-pre-0.50.md` through 0.75.x, then were removed in 0.76.0. Tarball-installed copies (no local `.git`) can browse both archives at commit [`bf15fe2`](https://github.com/knews2019/skill-do-work/tree/bf15fe2) on GitHub — they're preserved there.

---

## 0.82.0 — The Real Path (2026-05-28)

Closeout pass on five P2 findings from a fresh review of 0.81.0's `ai-report` action — three real bugs, one self-contradictory rule, and one design change to how the report stores its image binaries. One additional finding was rejected after history-check showed the reviewer was describing behavior that never existed.

- `actions/ai-report.md`: **image binaries are now external, not base64-inlined.** Each report owns a sibling `<report-stem>.assets/` folder next to its HTML — screenshots go in there with descriptive names (`before.png`, `after.png`, `live.png`) and the HTML references them via relative `src`. The HTML stays one file; it just has external image resources, like any normal webpage. Drops base64's ~33% per-image bloat, gives the user real binaries to inspect/re-process, and removes the awkward "self-contained but Tailwind+Mermaid via CDN" contradiction. The pair moves together — same portability story, smaller HTML. Updated across Philosophy, Step 3b, Step 4a, Step 5, Step 7, Output Format, Rules, Common Rationalizations, Red Flags, and Verification Checklist.
- `actions/pipeline.md` (line 349 markdown bullet + line 446 `.single.html` template): the relative path from `do-work/deliverables/` to project-root `ai-reports/` is `../../`, not `../../../` — the latter climbs one directory above the repo, so every generated deliverable link was broken. The markdown URL also pointed at the directory rather than the resolved filename. Both spots fixed; `grep -rn "\.\./\.\./\.\./ai-reports"` now returns zero.
- `actions/ai-report.md` Step 3a (asset search): `do-work/archive/UR-NNN/assets/` is now searched **first**. The action explicitly targets completed work, which `cleanup` has moved out of `do-work/user-requests/` and into `do-work/archive/` — so the live-only search was missing exactly the assets it most often needed and silently falling back to the weaker diagram path.
- `actions/ai-report.md` Step 3b (live screenshot): the dev-server probe walks `8080 → 5173 → 3000` but the screenshot URL was hardcoded to `http://localhost:8080/`, so a hit on 5173 or 3000 still tried to capture 8080 and failed. Now the probe captures the responding URL into `$DEV_URL` and the screenshot step uses that — and writes to the new `.assets/` layout.
- Reviewed and rejected: a fifth finding claimed the `forensics.md` Missing Qualifications check was "new" and would warn on REQs with `## Review acceptance`. Git history (`ab31114`, v0.39.0) shows the check is original to the action, and `## Review acceptance` isn't a heading anywhere in the codebase — `## Qualification` and `## Review` are designed as separate sections per `actions/sample-archived-req.md`. No change.

Plus four adjacent fixes from a second review pass on related surfaces:

- `actions/stray-check.md` Step 5 (line 73): the Input section advertised `default = report-only`, but the Step 5 skip condition only fired on the explicit `report` / `--report-only` token — so the default path still ran the fix prompt. Inverted to "skip unless `fix` / `--fix`", which is what the Input section already promised.
- `actions/stray-check.md` Step 2 + category 7: category 7 promised to flag empty directories, but the git-based inventory (`git ls-files` + `git ls-files --others --exclude-standard`) never emits directories — only files. Added a separate `find <scan-root> -type d -empty` filesystem pass in Step 2 that honors the skip-list, and clarified the category 7 detection wording to match. Outside-git mode already walked the filesystem, so it gets the same treatment for free.
- `actions/install.md` (line 35 + line 137): both upstream raw-GitHub URLs were 404. The `frontend-design` skill moved to a dedicated `anthropics/skills` repo (was `anthropics/claude-code`), and `disler/bowser` houses its SKILL.md under `.claude/skills/playwright-bowser/`, not `skills/playwright-bowser/`. Verified the new URLs with `curl -fsI` before swapping. Old fallback URL removed (it pointed at the same wrong path family).
- `next-steps.md`: added an **After ai-report** block. `SKILL.md` (line 324) requires next-step suggestions after every action, and ai-report was the only routable action without one. Suggests `slop-check` (validate against anti-slop), `present work` (complementary explainer), `inspect`, and `commit`.

## 0.81.0 — The Pixel Proof (2026-05-28)

New `ai-report` action: a single-file HTML report that anchors a completed UR/REQ in the literal pixels that changed. Where `present-work` explains the concept and `pipeline`'s `.single.html` debriefs the run, `ai-report` puts the screenshot front-and-centre with SVG callouts pointing at the delta — and falls back to SVG + Mermaid diagrams when bowser isn't available, so the report always ships.

- `actions/ai-report.md`: new standalone action. Inputs: `UR-NNN`, `REQ-NNN`, `most recent`, or empty. Outputs to `ai-reports/yyyy-mm-dd_hhmm_<slug>.html` (chronologically sortable filename per the external skill it was ported from). Pipeline: resolve target → collect before/after assets (UR `assets/`, `do-work/working/`, root `verify-*.png`, git-diff images) → optionally take a live screenshot via `playwright-cli` if a dev server responds → embed everything as base64 data URIs → write the report with hero / The Change / How It Works / What Changed / Verify It Yourself / Lessons sections.
- Bowser is **optional with graceful fallback**. If `playwright-cli` is missing or no dev server responds on common ports (8080, 5173, 3000), the report drops the live-screenshot pass and uses hand-coded SVG architecture + Mermaid data-flow diagrams instead. No install prompt, no block.
- Anti-slop applied inline. The action loads `crew-members/anti-slop.md` in Step 1 and runs the seven-principle self-check table at Step 6 before writing the file — no separate `slop-check` pass needed.
- `SKILL.md`: new priority-30 routing row (triggers: `ai-report`, `ai report`, `make-report`, `screenshot-report`, `visual report`, `proof of work`). Descriptive-content catch-all bumped to priority 31. Verb Reference, Action Dispatch, foreground-actions list, argument-hint frontmatter, top-of-file actions list, and help menu all updated.
- `actions/pipeline.md` Step 5: completion report now opportunistically links to `ai-reports/*{UR-NNN}*.html` when one exists — one bullet in the markdown rendering's "for the clueless-reader" Deliverables block, one tile in the `.single.html` "Related deliverables" card grid. Pure file-presence check; no dependency on the action having run.
- `actions/install.md`: bowser target now notes that `ai-report` is the second consumer (alongside `ui-review`) and falls back to diagrams when bowser isn't installed.
- `CLAUDE.md`: `actions/` listing gains `ai-report.md`.

## 0.80.2 — The Full Inventory (2026-05-27)

Two more correctness fixes from the same code review, plus the realization that one of them wasn't local. `stray-check` now sees junk inside brand-new untracked directories and stops letting its own skip-list hide committed artifacts — and the untracked-enumeration fix was applied everywhere the same pattern had been copy-pasted.

- `stray-check`: the untracked inventory now uses `git ls-files --others --exclude-standard` (lists files individually, honoring `.gitignore`) instead of `git status --porcelain`, which collapses a wholly-untracked directory into a single `?? dir/` row — junk like `tmp/debug.log` was invisible to every filename/size/content check.
- `stray-check`: the noise skip-list (`__pycache__/`, `dist/`, …) now applies only to untracked/ignored content; a *tracked* file inside those dirs still reaches the committed-artifact checks, which is the whole point of category 3.
- Same `git status --porcelain` → `--untracked-files=all` fix applied to `commit`, `inspect`, and `work`, where the identical pattern would have missed (or tried to "read") files inside a new untracked folder.
- Codified the two git-command traps in CLAUDE.md so future actions avoid them — and the rule that a prescribed-command bug found by review is rarely local: grep the primitive across every action.

## 0.80.1 — The Root Cause (2026-05-27)

Install-safety fix plus three correctness fixes from a code review. The big one: the shipped `.gitignore` ignored all of `do-work/`, and since the repo installs by extracting its files into your project root, that rule landed in end-user projects and blocked the do-work folder it's supposed to commit. An ignore rule's reach follows where it sits — a project-root rule over-reaches — so nothing `do-work/`-related ships into the root anymore.

- `.gitignore` now ships only `do-work/pipeline.json` (transient state); this source repo keeps its own `do-work/` untracked via local `.git/info/exclude`, which never ships.
- `stray-check`: the "tracked but should-be-gitignored" check now uses `git check-ignore --no-index` — plain `check-ignore` never reports already-tracked files, so the category was silently finding nothing.
- `prompts`: the `Runnable:` opt-out guard now parses the first token, so `Runnable: no — placeholder…` correctly refuses to run.
- `slop-check`: the report template no longer shows a `PASS` row with blank evidence, which had undercut its own "every row needs evidence" rule.

## 0.80.0 — The Lost and Found (2026-05-27)

New `stray-check` action: a repo-wide sweep for orphan and junk files that pollute where they don't belong — the whole-repo sibling to forensics, which only ever looked at do-work's own files. It reports first and touches nothing until you confirm.

- Detects stray temp/backup/OS files, committed build artifacts, tracked-but-should-be-gitignored files, committed secrets (critical), misplaced/duplicate/empty files, oversized binary blobs, AI scratch droppings, and best-effort dead code.
- Report-only by default; `fix` applies the safe, reversible fixes (delete untracked junk, `git rm --cached`, gitignore) only on explicit confirmation. Never `git add -A`, never auto-commits.
- Skips the entire `do-work/` tree and defers misplaced `do-work/` directories to cleanup. Routing carve-out keeps "clean up junk files" / "find orphan files" out of cleanup.

## 0.79.1 — The Dream Lane (2026-05-26)

Routing fix: `consolidate memory` / `clean up wiki` / `memory cleanup` now reach the dream action instead of being swallowed by cleanup.

- Scoped cleanup to archive-only; gave dream's memory/wiki/notes phrases precedence over cleanup's generic verbs.

## 0.79.0 — The Quiet Pass (2026-05-25)

A new `dream` action — a manual, explicit four-phase pass that consolidates a plain-text memory directory: lint mechanical rot, heal contradictions, prune near-duplicates, rebuild the index. Destructive by design, so it never auto-triggers; invoke it when memory has visibly decayed.

- `actions/dream.md`: new standalone action. Resolves a default memory dir (`./memory`, `./wiki`, `./kb/wiki`, `./knowledge-base/wiki`) or accepts an explicit path. Holds `.lock` for the duration. Phase 1 orients (read index, page frontmatter, recent log). Phase 2 runs seven deterministic checks expressed inline as prompt steps (pages missing from index, index dangling, broken `[[wiki-links]]`, orphan pages, stale frontmatter dates, relative-date occurrences, near-duplicate titles) — no script bundled, every check spelled out with the exact regex and worklist payload so the agent reproduces the deterministic behavior in-prompt. Phase 3 fixes mechanical issues first, then resolves contradictions (newest wins), pins relative dates to absolute, merges duplicates (repoints inbound links before deletion), prunes untrue. Phase 4 rebuilds the index ≤200 lines, bumps `last_updated`, appends a `[dream]` line to `log.md`, removes `.lock`.
- Auto-detects bkb wikis: `_master_index.md` is a first-class index alongside `MEMORY.md` and `index.md`, so `do-work dream` works against `kb/wiki/` without extra arguments. Coexists with `bkb lint`/`garden`/`defrag` — those are routine read-only or conservative hygiene; dream is the aggressive single-pass consolidation that merges, prunes, and rebuilds.
- `SKILL.md`: new priority-28 routing row for dream (triggers: `dream`, `consolidate memory`, `clean up wiki`, `lint and merge notes`, `memory cleanup`); descriptive-content catch-all moves to priority 29. Top-of-file actions list, Verb Reference, Action Dispatch table, foreground-actions list, argument-hint frontmatter, and help menu all updated.
- `next-steps.md`: new "After dream" block suggesting `commit`, `bkb lint`, or another dream pass.
- `CLAUDE.md`, `README.md`: `actions/` listing and Other-actions reference gain `dream`.

## 0.78.3 — The Dimension Pair (2026-05-25)

`code-review` Step 4 (Pattern & Architecture Review) now names two more dimensions reviewers were quietly missing. Folder cohesion catches files that don't belong in the folder they live in; cyclomatic complexity gets promoted from a quick-wins tie-breaker into a first-class architectural check, explicitly distinguished from Step 3's circular-dependency check so the two don't get conflated.

- `actions/code-review.md`: Step 4 dimension table gains two rows. **Folder cohesion / orphan files** — checks whether imports match folder domain, whether file shape matches siblings, and whether folders have become junk drawers; contrasted with Step 3's "structural consistency" angle. **Cyclomatic complexity** — branch counts, nested conditionals, sprawling switches, predicate chains; explicitly contrasted with Step 3's "Circular dependencies?" check to prevent the McCabe-vs-cyclic-deps mix-up. No thresholds named (consistent with the rest of the table). No changes to Step 9's report template — new findings flow into the existing Architecture table.
- `actions/quick-wins.md`: unchanged. Cyclomatic complexity stays in Step 5 as a risk-impact tie-breaker; this release adds coverage to code-review without removing anything from quick-wins. The two actions stay complementary.

## 0.78.2 — The Audit Sunset (2026-05-21)

Removed `DEAD_CODE.md` from the repo root. It was a point-in-time audit snapshot from 0.77.0, and every actionable finding it raised has since been closed out — so the report now describes a tree that no longer exists.

- `DEAD_CODE.md`: deleted. Its findings were all resolved in 0.77.0 (`performance.md` removed, action→guide cross-links added, placeholder-prompt opt-out marker added, broken ADR-012 link fixed, orphaned imported-spec annotated). The file was never wired into the skill — nothing routed to or loaded it — so removal is pure cleanup. Full report remains in git history at commit `73d4955`. The historical 0.77.0 changelog entry that mentions it is left intact as a record.

## 0.78.1 — The Review Trim (2026-05-21)

Codex review pass on 0.78.0. Two precision fixes — both surfaced unreachable or noisy behavior in the just-added slop-check routing and default-target resolution.

- `SKILL.md`: dropped `check slop`, `check draft`, `check artifact` from the priority-27 routing row and the Verb Reference. They collided with priority 5 verify (which already claims any `check ...` form), so users invoking them would have hit the verify route, not slop-check. Distinctive triggers (`slop-check`, `slop check`, `anti-slop`) stay; the Verb Reference now states the exclusion rule directly so future contributors don't re-add the trap.
- `actions/slop-check.md` (Step 2, point 3): "most recent" resolution now prefers authored artifacts. Globs `*.md` and `*.single.html`; skips `*.marp.html` (mechanical Marp-CLI exports of the `.marp.md` source) and `*-video/` directory contents (Remotion TSX source, not prose). Previous newest-by-mtime heuristic would frequently pick the mechanical Marp HTML right after a pipeline completion, flagging HTML scaffolding instead of the authored draft.
- `CHANGELOG.md` 0.78.0 entry: corrected the inaccurate trigger list (was claiming `check slop`/`check draft`/`check artifact` as distinctive — they weren't).

## 0.78.0 — The Slop Filter (2026-05-21)

A guardrail against AI slop — bloated, unverified, conclusion-buried artifacts that pass the cost of clarity onto the reader. Adds a new behavioral crew-member that auto-loads whenever an artifact is being generated for a human, plus a standalone `slop-check` action to grade any draft against the seven principles before it ships.

- `crew-members/anti-slop.md`: new always-on-during-artifact-generation crew-member. Seven principles in one frame — producer absorbs the cost of clarity, reader doesn't. Loaded by present-work (Step 4 drafting), review-work (Step 9 report), pipeline (Step 5 completion-report rendering), and kb-lessons-handoff (Step 2 source-document assembly). Boundaries explicitly exempt code output (karpathy.md territory), agent status updates (caveman.md / general.md), and commit messages.
- `actions/slop-check.md`: new read-only action that loads the crew-member and grades a target artifact against each of the seven principles. Inputs are flexible — file path, REQ/UR ID, "most recent" deliverable, or pasted text. Outputs a findings table (principle | status | evidence | fix) plus a top-line verdict (Clean / Borderline / Slop) and a single concrete top fix. Optional rewrite only on explicit user confirmation; preserves factual claims verbatim.
- `SKILL.md`: new priority-27 routing row for slop-check (distinctive triggers only — `slop-check`, `slop check`, `anti-slop`; any `check ...` form collides with verify priority 5 and is intentionally excluded). Verb Reference, Action Dispatch, foreground-actions list, argument-hint, top-of-file action listing, and help menu all updated. Descriptive-content catch-all moved to priority 28.
- `actions/present-work.md`, `actions/review-work.md`, `actions/pipeline.md`, `actions/kb-lessons-handoff.md`: each step that begins composing a human-facing artifact now loads `crew-members/anti-slop.md` explicitly — no behavioral change for any other step.
- `next-steps.md`: post-`present-work` suggestion now includes `do-work slop-check`; new "After slop-check" block points at re-checks, regeneration, and follow-up capture.
- `CLAUDE.md`: `actions/` listing gains `slop-check.md`; crew-members loading-behavior list gains `anti-slop.md` with its exact load conditions and boundaries.

## 0.77.0 — The Reach Audit (2026-05-19)

Closeout pass on a dead-code audit of the skill. Tightens the Schema Read Contract so the `domain` enum is honored consistently, removes a crew-member file that was reachable in letter but never in spirit, cross-links every action file to its user-facing guide, and adds a machine-readable opt-out marker for placeholder prompts.

- `actions/work.md`: Route C planning (Step 4) and review-work spawning (Step 9) now explicitly normalize `domain` per the Schema Read Contract — matches the narrative claim that "every read site honors a uniform normalize-and-warn contract." The per-field table's read-sites column for `domain` is updated to list all three load sites instead of only Step 6.
- `crew-members/performance.md`: removed. Unreachable for non-canonical domains under the tightened contract, and `performance` was never in the canonical enum (`frontend | backend | ui-design | general`). `CLAUDE.md`'s example list updated accordingly.
- `prompts/`: new optional `**Runnable:**` header key for placeholder/sidecar prompts. `prompts/weekly-signal-diff-personal.md` opts in with `Runnable: no`; the dispatcher in `actions/prompts.md` (Sub-Command `run`, step 3) refuses opt-out prompts with a contextual explanation from the prompt's first-line description. Absence of the key means runnable — the safe default.
- 18 action files now end their top-of-file blockquote with a `User-facing walkthrough:` link to the corresponding `docs/*-guide.md`. Previously only `capture`, `work`, and `interview` had docs links (and those weren't in the blockquote). Now uniform across `bkb`, `capture`, `cleanup`, `code-review`, `commit`, `forensics`, `inspect`, `interview`, `present-work`, `prime`, `prompts`, `quick-wins`, `review-work`, `roadmap`, `ui-review`, `verify-requests`, `version`, `work`.
- `decisions/records/adr-012-interview-v2-gap-closure.md`: removed the broken `References` bullet pointing at `decisions/imported-specs/2026-04-16_expand-skill-do-work-interview.md`. That file was intentionally deleted in `0.71.1` (commit `f7e4b61`); restoring it would re-open a closed decision.
- `decisions/imported-specs/2026-04-17_improve-weekly-diff-skill.md`: added a Status footer documenting that edits 1–3 from the spec landed in `prompts/weekly-signal-diff.md`. The spec is a candidate for a future ADR-013 if the maintainer wants the decision rationale in the ledger.
- `DEAD_CODE.md`: full audit report committed at the repo root with findings grouped by confidence. This release is the closeout of the items that were actionable; two other observations in the report (the `AGENTS.md` stub and `.vscode/tasks.json` portability) were independently addressed in 0.76.5.

## 0.76.5 — The Stale Wipe (2026-05-19)

Six janitorial fixes from a `quick-wins` self-scan — stale docs swept out, two shell hooks hardened, and an invariant documented so the non-jq fallback can't silently miscount.

- `CLAUDE.md` Project Structure: dropped the `_dev/` line — the directory was emptied in 0.75.0 and the entry was a dead pointer.
- `README.md` "fully clean update" path list now matches `actions/version.md`'s authoritative shipped-paths glob (was missing `prompts/`, `interviews/`, `specs/`, `docs/`, `decisions/`, `hooks/`, `CLAUDE.md`, `AGENTS.md`, `next-steps.md`).
- `.vscode/tasks.json` gained `linux` (`xdg-open`) and `windows` (`cmd /c start`) overrides for the "Open current HTML in browser" task; macOS behavior unchanged.
- `actions/pipeline.md`: documents the pretty-print invariant for `do-work/pipeline.json` — pipeline.md is the sole writer and the constraint protects `hooks/pipeline-guard.sh`'s line-oriented grep fallback from miscounting on compact JSON.
- `AGENTS.md`: replaced the newline-less `READ CLAUDE.md` stub with a one-line markdown link (`See [CLAUDE.md](CLAUDE.md).`) — clickable when rendered, POSIX-clean.
- `hooks/session-start.sh`: anchored the version-line `sed` so it strips only the `**Current version**:` prefix instead of greedily up to the last `: `. Same output today, robust to future colon-containing version lines.

## 0.76.4 — The Quiet Drain (2026-05-17)

Removes the `--halt-on-failure` flag from `do-work run`. The flag was redundant with the existing auto-follow-up pattern — `review-work` Step 10 already creates `pending` / `pending-answers` follow-ups for failed and completed-with-issues outcomes, and `do-work clarify` is the documented batch-triage path. The default loop is now the only loop: classify, archive, queue follow-ups, continue.

- `actions/work.md`: dropped the `--halt-on-failure` Input bullet and the halt-check branch at the top of Step 10. The four-section exit summary (completed/done, pending-answers, blocked-archive-collision, blocked-by-dependencies) and Session Checkpoint behavior are unchanged.
- `SKILL.md`: removed the flag from the priority-4 routing example and rewrote the work-action Notes cell to mention only `--wave N`.
- `docs/work-guide.md`: rewrote the third "What `run` does" bullet to state the loop-always-continues guarantee and point at `do-work clarify` for triage, instead of describing an opt-in halt.

## 0.76.3 — The Typo Guard (2026-05-17)

Extends `0.76.2`'s defensive `dependencies:` alias to every other field where a natural muscle-memory typo would have been silently swallowed. Pairs the read-only field-name alias pattern (when the YAML key is wrong) with a uniform normalize-and-warn contract (when the enum value is wrong), and closes a near-miss-keyword fall-through in the pipeline dispatcher.

- `actions/pipeline.md` Step 1: mode-selection table now normalizes `$ARGUMENTS` first (trim + lowercase), then guards against single-token near-misses of `status`/`abandon` (e.g., `stat`, `aban`) — they trigger a "Did you mean ...?" prompt instead of silently initializing a new pipeline with the typo as the request text. Same shape as the install-dispatch normalization that landed in 0.75.1.
- `actions/work.md` Schema Read Contract: new section documenting the uniform normalize-and-warn rule for seven enum-or-boolean fields (`domain`, `status`, `route`, `caveman`, `tdd`, `error_type`, `kb_status`) — each gets a per-field alias map (e.g., `back-end` → `backend`, `done` → `completed`) plus a documented default-on-unknown that emits the warning rather than silently dropping. `addendum_to` also gets the 0.76.2 field-name alias treatment (`amends`/`parent`/`amendment_to` recognized when canonical is absent). Step 6 crew-load and Step 8 upstream-walk/cycle-detection updated to honor the contract.
- `actions/capture.md`: new Schema Aliases section listing the five field-name aliases (`addendum_to`, `depends_on`, `batch`, `related`, `suggested_spec`); points downstream readers at `actions/work.md`'s contract for shared enum normalization. Capture validates non-canonical enum values during emission and prompts for correction — capture is the human-attention window for catching typos at the source.
- `actions/bkb.md`: new Schema Read Contract for wiki-page frontmatter — normalize-and-warn for `type`, `rel`, `confidence`; field-name alias for `topic_cluster` (`topic`/`topic_category`). Applies across `triage`, `ingest`, `lint`, `garden`, `defrag`, and `query` sub-commands.
- `actions/roadmap.md`: Ready and Blocked rubrics honor the `addendum_to` aliases; Step 1 inventory paragraph references the work.md contract so the same normalization applies to roadmap classification — REQs with non-canonical field values land in the right bucket instead of being orphaned.

## 0.76.2 — The Safety Alias (2026-05-17)

Belt-and-suspenders defensive read for the dependency-aware selection added in 0.74.0. Codex flagged a P1 on the 0.74.0 PR claiming `depends_on` was a rename of a legacy `dependencies:` frontmatter field and that pre-rename queues would silently bypass gating. The premise was wrong — `depends_on` was introduced fresh in 0.74.0 and no `dependencies:` frontmatter ever existed in the schema — but a downstream user typing `dependencies:` from Python/Node/Cargo muscle memory would today have it silently ignored. The alias makes the typo harmless.

- `actions/work.md`: schema, Step 1 dependency-aware selection, Step 1 cycle detection, Step 1 `--wave` depth calculation, and Step 8 upstream-failure short-circuit now all read `dependencies:` as a synonym for `depends_on` when only the alias is present. `depends_on` wins when both are present. Capture and Step 8 follow-up REQs continue to emit only the canonical `depends_on:` — the alias is read-only, never propagated.
- `actions/roadmap.md`: Ready and Blocked rubrics honor the same alias when classifying pending REQs.
- `docs/work-guide.md`: the dependency-aware ordering bullet names the alias so it surfaces in the user-facing doc, not just in the action spec.

## 0.76.1 — The Archive Pointer (2026-05-17)

Post-PR-review fixup for 0.76.0. Codex flagged that tarball-installed users lose access to pre-0.65 release notes once the archives are deleted — `.git` isn't always present, so "git history" isn't always a valid fallback. Restored discoverability without restoring the files.

- `CHANGELOG.md` header now points tarball users at commit [`bf15fe2`](https://github.com/knews2019/skill-do-work/tree/bf15fe2) on GitHub where both archive files are still readable.
- `actions/version.md` glob note refined to acknowledge the tarball gap explicitly.

## 0.76.0 — The Trim Pass (2026-05-17)

Three-way cleanup: stale CHANGELOG archives removed, `actions/pipeline-reference.md` re-inlined into `actions/pipeline.md`, and "Do NOT use when" routing bullets across seven action files collapsed to a single `SKILL.md` pointer. ~1,800 fewer lines on disk, same functionality.

- `CHANGELOG-2026-spring.md` and `CHANGELOG-pre-0.50.md` removed from the working tree. The split-archive scheme introduced in 0.75.0 went unmaintained; entries remain reachable via git history. `CHANGELOG.md` header pointer and `actions/version.md`'s dirty-file glob updated to match.
- `actions/pipeline-reference.md` re-inlined into `actions/pipeline.md`'s Output Format section as `#### Composition rules` plus the three numbered template subsections (Plain Markdown Report, Marp Slide Deck, Standalone HTML Debrief). The companion was always loaded with the main action, so the split was paying a cross-file tax for no gain — same proven pattern as 0.75.0's `work-reference.md` re-inline. Combined `pipeline.md` is 564 lines, well under the `work.md` baseline (1058) that's already shipping.
- `decisions/records/adr-001-modular-action-prompts-and-companion-references.md` and `adr-008-...` updated to acknowledge the re-inline as a counter-example. `decisions/topics/_index_pipeline-deliverables.md` sources list trimmed.
- "Do NOT use when" sections normalized in 7 action files where 2+ bullets were pure sibling-action routing (`bkb`, `deep-explore`, `interview`, `pipeline`, `quick-wins`, `scan-ideas`, `work`). Routing bullets replaced with one `SKILL.md` pointer; state and scope constraints preserved. The other 13 action files were intentionally left alone — their bullets carry non-routing guidance worth keeping.

## 0.75.1 — The Review Catch (2026-05-17)

Post-PR-review fixups for 0.75.0. Codex caught two real regressions and they're now addressed.

- `crew-members/general.md`: extended with four sections (Lessons Discipline, Test-Writing Posture, Cross-REQ Test-Break Rules, Discovered-Tasks Contract). The 0.75.0 trim of `actions/work.md` Step 6 replaced inline builder rules with a pointer claiming `general.md` carried them — but `general.md` only had PRIME-file philosophy. The pointer claim is now true; enforcement restored for every REQ.
- `SKILL.md` install dispatch note: explicit normalization rules for the trigger aliases. Hyphenated forms (`install-ui-design`, `install-bowser`) now strip the `install-` prefix before target extraction; `setup`-prefixed forms (`setup bowser`, `setup ui design`) strip the leading `setup`. Previously these aliases would fall through to the help block instead of installing.

## 0.75.0 — The Lighter Pack (2026-05-17)

Cross-cutting cleanup pass: seven simplifications dispatched as parallel background agents, all touching disjoint file sets. Smaller surface, same functionality.

- `CHANGELOG.md` split into a 501-line live file plus two archives (`CHANGELOG-2026-spring.md` for 0.50.0–0.64.x, `CHANGELOG-pre-0.50.md` for 0.1.0–0.49.x). All 248 historical entries preserved, just relocated. `actions/version.md` Step 1 unchanged — pattern-based parse still picks the newest 5 from the live file.
- `actions/install-ui-design.md` + `actions/install-bowser.md` folded into one parameterized `actions/install.md` with a per-target manifest table. Users still type `install ui-design` / `install bowser` — both trigger phrases remain registered keywords.
- `actions/work-reference.md` re-inlined into `actions/work.md`. The companion was always loaded with the main action, so the split was paying a cross-file tax for no gain. `decisions/records/adr-001-modular-action-prompts-and-companion-references.md` updated to acknowledge the re-inlining as a counter-example.
- `actions/work.md` Step 6 "all routes include these instructions" block compressed from 36 lines to 10 bullet pointers — the rules live in `crew-members/general.md` and `crew-members/karpathy.md`, so the agent only needs pointers, not restatement.
- `actions/capture.md`: removed inline "Next steps" block that duplicated `next-steps.md`, plus the second copy of the "you can run verify requests" suggestion.
- "Do NOT use when" sibling-action redirects trimmed across 7 review/diagnostic action files (forensics, roadmap, inspect, code-review, ui-review, review-work, verify-requests). `SKILL.md`'s routing table is the authoritative dispatcher; per-file restatement was just drift risk.
- `README.md`: 427 → 153 lines. Per-action sections that duplicated `SKILL.md`'s help menu collapsed into one "Other actions" pointer. Install + capture + run + pipeline workflows kept full prose since those are the headline workflows.
- `next-steps.md`: 6 BKB maintenance sub-commands (close, rollup, defrag, garden, crew, resolve) merged into one "After bkb (maintenance subcommands)" block; install blocks merged into a single "After install".
- `_dev/code-review-20-commits.md` deleted — marked RESOLVED back in April 2026 and retained "as historical artifact"; no longer needed.

## 0.74.0 — The Linked Run (2026-05-17)

Adds dependency-aware execution to `do-work run`. REQs can now declare `depends_on` in frontmatter; the work loop honors it for selection order, surfaces upstream failures during classification (so cascading failures aren't misdiagnosed as fresh code bugs), and supports `--halt-on-failure` and `--wave N` flags for foundation-phase work where late-stage REQs depend on early-stage ones being correct. Also folds in the Codex P2 finding on `actions/roadmap.md` from the 0.73.5 PR.

- `actions/work.md` Request File Schema: new optional `depends_on: []` field, semantically distinct from `addendum_to` ("requires that REQ to be done first" vs. "amends that REQ"); a REQ can carry both.
- `actions/work.md` Step 1: honors `depends_on` for selection (a REQ is dependency-ready when all its `depends_on` resolve to completed/completed-with-issues); new `blocked-by-dependencies` section in the composed exit summary; new `blocked-dependency-cycle` held status for cycles in the `depends_on` graph.
- `actions/work.md` Step 1: optional `--wave N` flag filters the scan to REQs at dependency depth N for checkpointed wave-by-wave execution. Mutually exclusive with targeted REQ IDs.
- `actions/work.md` Step 8: upstream-failure short-circuit during failure classification — if any `addendum_to` or `depends_on` ancestor is `failed`, classify as `spec` with the original error wrapped in an upstream pointer. Cascades now show up in the follow-up REQ's error message instead of presenting as fresh code bugs in the wrong domain.
- `actions/work.md` Step 10: optional `--halt-on-failure` halts the loop after a failed or completed-with-issues REQ; default behavior unchanged.
- `actions/capture.md`: documents emitting `depends_on` and the slicing convention; the optional `## Dependencies` prose section remains for human readers, but frontmatter is the source of truth for tooling.
- `actions/roadmap.md`: Ready and Blocked rubrics both reference `depends_on` and `blocked-dependency-cycle` — resolves the fall-through bucket the Codex P2 review flagged on the 0.73.5 PR (Ready required dependencies-archived while Blocked no longer included dependency-list checks).
- `SKILL.md`: work-action routing recognizes `--halt-on-failure` and `--wave N` flags; strips them before extracting REQ IDs.
- `docs/work-guide.md`: rewrites the "What `run` does NOT do" section to reflect the new opt-ins.

## 0.73.5 — The Honest Run (2026-05-17)

Documents what `do-work run` does (and does not) do across a bulk queue, and removes a roadmap rule that referenced a frontmatter field the REQ schema never defined.

- `docs/work-guide.md`: new "What `run` does NOT do" section — no dependency ordering, no mid-run pause, no halt on failure. Surfaces three properties first-time users routinely hit by surprise.
- `actions/work.md` Step 1: now states queue order is purely numeric and points readers to `do-work roadmap` before bulk runs, right where the order is established.
- `actions/roadmap.md` Step 2: Blocked classification no longer references a non-existent `dependencies` frontmatter field; `addendum_to` and external-dependency-in-prose remain — both are real backings.

## 0.73.4 — The Fresh Read (2026-05-14)

Fixed a spec bug in the interview `status` sub-command that could report stale data after an in-memory migration.

- `actions/interview.md` `<template> status`: the step ran the Session-Load Protocol in dry-run mode (in-memory migration only) but then told the agent to re-read `session.json` from disk — discarding the migrated shape and rendering status from stale pre-migration data. It now renders directly from the in-memory session object the protocol hands back.

## 0.73.3 — The Downgrade Guard (2026-05-12)

Four bug fixes from a code review of the 0.72.x → 0.73.x cluster. Two were real spec bugs in the Session-Load Protocol (silent template downgrade, ambiguous CHANGELOG noise for stamp-only refreshes); two were defects in the roadmap action's `find` examples (literal `HHMMSS-` placeholder, fragile `-o` precedence).

- `actions/interview-reference.md` Session-Load Protocol Step 3: split the same-major branch by direction. Session-older-than-template stamps forward as before; session-newer-than-template (template was rolled back via `git checkout`) is now a no-op read in both modes — the old wording would have downgraded the stamp and lost the record of which template version generated the session.
- `actions/interview-reference.md` Session-Load Protocol Step 3: explicitly carved out the CHANGELOG.md append rule. Stamp-only refreshes (same-major minor/patch bumps) skip the CHANGELOG entirely — only cross-major migrations in Step 4c append. Previous wording referenced "4c's persist rules" which would have logged every minor template bump as `auto-migrated session: X → Y` even though no migration ran.
- `actions/roadmap.md` Step 3: rewrote the `find` examples to use the actual `[0-9][0-9][0-9][0-9][0-9][0-9]-<kb_entry>` glob instead of the literal `HHMMSS-` placeholder. An agent following the code block literally would have searched for files named `HHMMSS-foo.md` and found none; the clarifying paragraph below was a workaround, not a fix.
- `actions/roadmap.md` Step 3: wrapped the `-name A -o -name B` predicates in `\( … \)` parentheses and added explicit `-print`. Without them, appending any predicate to the command makes `-o` bind lower than the implicit action and silently drops half the matches.

## 0.73.2 — The Dry Verbs (2026-05-11)

Replaced the "drain" metaphor in queue-processing docs with clearer verbs (work through / process / clear). User feedback flagged "drain" as reading wet/unnatural for a task queue.

- Swept SKILL.md, next-steps.md, and action files (work, pipeline, roadmap, kb-lessons-handoff) to swap "drain"/"draining"/"drains" for context-fit alternatives.
- Renamed ADR-006 from `pipeline-drains-follow-up-work` to `pipeline-processes-follow-up-work` and updated every wikilink reference across decisions/.
- No behavior change — pure docs/prompts polish.

## 0.73.1 — The Convention Match (2026-05-09)

The two editorial polish items from the 0.73.0 review pass. Both align the just-extracted Session-Load Protocol references with the conventions the rest of the file already uses for extracted heavy content.

- `actions/interview.md`: dropped the standalone `## Session-Load Protocol` heading and folded the stub into "Locating the Session" as a follow-on paragraph. The other heavy content extracted to interview-reference.md (Template File Format, Canonical Entry Contract, Checkpoint File Format, Re-run Modes, Versioning Scheme, Ingest File Mapping, Export Schemas, Mid-layer Recovery) doesn't get parallel `##` headings in the action file — they're referenced inline. The new paragraph names the protocol, lists its two modes, and points to the reference, all inside the existing session-location section.
- `actions/interview.md`: the per-subcommand pointers no longer repeat the protocol's location on every invocation. The first mention in "Locating the Session" establishes where to find the spec; the per-subcommand calls just say "Run the **Session-Load Protocol** in **persist** mode" without re-pointing. Drops a doubled "Session-Load Protocol" phrase that read awkwardly.

## 0.73.0 — The Protocol Move (2026-05-09)

The Session-Load Protocol grew to ~50 lines of dense spec across seven patches in 0.72.x. That's the kind of heavy content `actions/interview-reference.md` is for per ADR-001, so the protocol moved there.

- `actions/interview-reference.md`: gained a "Session-Load Protocol" section between the `session.json` schema and the Checkpoint File Format. Full spec — mode-selection table, all four steps with substeps, version placeholder conventions, multi-major chain rules, atomic write semantics, concrete dry-run rendering — lives here now.
- `actions/interview.md`: collapsed the Session-Load Protocol section to a six-line stub that names the protocol, lists its two modes (persist vs dry-run), and points readers to the reference. Per-subcommand pointers now read `(spec in actions/interview-reference.md)` instead of "see top of this file." The action stays a short entry-point document; the heavy specification stays in the companion file.
- Top-of-file architecture summary in `actions/interview.md` now lists "Session-Load Protocol" alongside the other heavy-content items extracted per ADR-001.

No behavior change — only the location and discoverability of the spec.

## 0.72.7 — The Semver Honor (2026-05-09)

Two real correctness bugs from the 0.72.6 review pass plus an editorial cleanup. Both bugs would have fired on every minor or patch template bump, so worth the patch.

- `actions/interview.md` Session-Load Protocol Step 4b: dropped the stray leading `v` from the bump expression. The previous text said "bump the in-memory `template_version` to `v<old-major+1>.0.0`" — but the `template_version` field is a bare semver string per `actions/interview-reference.md`'s schema. Implementations writing `v2.0.0` into the field would break version comparison on the next protocol run.
- `actions/interview.md` Session-Load Protocol Step 3: same-major older versions (e.g., session `2.3.0` against template `2.5.0`) now short-circuit to a stamp-only path instead of falling through to Step 4. Semver minor/patch bumps are non-breaking by contract, so applying a `Migration from v2.x` section to a same-major upgrade would corrupt rather than migrate. Step 4 now triggers only on cross-major drift, which is what the chain logic was designed for.
- `actions/interview.md` Session-Load Protocol Step 4 header: dropped the backslash-escaped angle brackets (`v\<major\>.x` → `v<major>.x`) so the placeholder convention matches the rest of the section.

## 0.72.6 — The Spec Sharpening (2026-05-09)

The four review carryovers from the 0.72.4/0.72.5 review pass. The standout was a real spec ambiguity in the migration protocol — `<old>` was used to mean both "full version string" (for messages) and "major-version component" (for section lookup). Two implementations following the spec literally would diverge.

- `actions/interview.md` Session-Load Protocol Step 4: introduced explicit placeholder conventions — `<old>` and `<new>` for full version strings (used in user-facing messages), `<old-major>` and `<new-major>` for major-version components (used to look up `## Migration from v<major>.x` sections). Section-lookup now unambiguously matches `v1.x` for any session at `1.0.0`/`1.4.7`/etc.
- `actions/interview.md` Session-Load Protocol Step 4a/4b: spec'd multi-major-version migration chains. A session at `1.x` against a `3.x` template now requires `Migration from v1.x` AND `Migration from v2.x` to both exist; the protocol applies them in order, advancing the in-memory `template_version` by one major per pass. Authors who want to skip a major must write a passthrough section rather than omitting it.
- `actions/roadmap.md` Step 3: added a durability caveat to the within-branch tie resolution — `processed`'s `YYYY-MM-DD/` lexicographic sort survives `git clone` and archive restores; `capture` and `inbox`'s mtime-based resolution does not. Readers should treat the roadmap as a snapshot of the current filesystem, not a stable identifier across machines.
- `actions/interview.md` Session-Load Protocol Step 4c: replaced the placeholder dry-run example block with a concrete rendering — full status output (Interview status header, layer table, Review/Previous version lines) followed by the `⚠` staleness notice with real version strings (`1.0.0` → `2.0.0`). Implementations can now diff against a real format instead of reconstructing it.

## 0.72.5 — The Polish Bundle (2026-05-09)

The four P3 carryovers from the 0.72.2/0.72.3 self-review pass.

- `actions/roadmap.md` Step 3: added a bash globstar caveat. `**` is opt-in in bash (`shopt -s globstar`, off by default) but on by default in zsh — readers running in bash without globstar get either no match or a literal-`**` match. Recommended `find` instead in that case.
- `actions/roadmap.md` Step 3: spec'd within-branch tie resolution. When the same `kb_entry` matches multiple files in one branch (re-ingest, `HHMMSS-` collision sibling), the most recent wins — lexicographic sort on `processed/YYYY-MM-DD/`, mtime elsewhere.
- `actions/roadmap.md` Output Format header: aligned the `**Lessons:**` line's last label with its section header — `[N missing]` is now `[N file not found]`, matching `## Lessons File Not Found` so totals roll up to readable section names.
- `actions/interview.md` Session-Load Protocol Step 4c: spec'd the dry-run staleness notice's placement in the `status` output — blank-line separator, then the `⚠` line, no trailing blank — and showed the exact format inline.

## 0.72.4 — The Precondition Fix (2026-05-09)

The Session-Load Protocol's "no migration path documented" branch was sequenced after the migration-apply step instead of before it. If a template lacked a `Migration from vX.x` section, the protocol would attempt to apply zero steps and silently bump `template_version`, corrupting the session shape, instead of bailing with the documented error message. Restructured so the precondition check runs first.

- `actions/interview.md` Session-Load Protocol Step 4: split into 4a (verify a migration path exists; abort with the documented error if not), 4b (apply migration steps), 4c (persist or report). The previous Step 5 is now Step 4a — it gates Step 4b instead of being a never-fires fallback after it.

## 0.72.3 — The Lesson Roll-Up (2026-05-09)

The two P3 carryovers from 0.72.2's self-review. The roadmap report now surfaces lesson workload at the same altitude as REQ workload.

- `actions/roadmap.md` Output Format header: added a `**Lessons:**` totals line next to the existing `**Totals:**` and `**TDD posture (pending):**` lines, rolling up all five lesson buckets (awaiting triage / awaiting ingest / processed / pending handoff / missing) at a glance.
- `actions/roadmap.md` Suggested Next Steps: added template lines for the four actionable lesson buckets (`bkb triage` + `bkb ingest`, `bkb ingest`, investigate File Not Found, re-run handoff). The list is filtered — items only emit when their bucket has at least one REQ — so the rendered output stays compact when there's nothing to do.

## 0.72.2 — The Read-Only Honor (2026-05-09)

A self-review of 0.72.1 caught four real issues in the just-shipped code: a read-only subcommand had been quietly turned into a mutator, the migration write had no error path, the work loop's exit semantics weren't explicit, and the KB lookup ignored `bkb`'s collision-prefix rule. All four fixed.

- `actions/interview.md`: split the Session-Load Protocol into **persist** and **dry-run** modes. `status` now uses dry-run — migration happens in-memory only, no `session.json` write, no `CHANGELOG.md` append, and the output gets a one-line staleness notice instead. Mutating subcommands (`<template>` resume, `review`, `export`, `ingest`) keep the persist mode but now use atomic write-then-rename and abort the calling subcommand on write failure rather than silently leaving an inconsistent on-disk state. Dropped the misleading `versions` reference from the protocol's enumeration.
- `actions/work.md`: Step 1's composed exit path now states explicitly "After rendering all applicable sections, exit the work loop" so an agent reading strictly doesn't fall through to Step 2.0 after rendering. The pending-answers section dropped the `[N] open questions` count — Step 1 only reads frontmatter, so the count would have required reaching into REQ bodies. The count belongs to `do-work clarify`, where it lives now.
- `actions/roadmap.md` Step 3: the recursive `kb_entry` lookup now also matches `HHMMSS-<kb_entry>` (bkb's collision-prefix rule from `bkb.md` Step 6 of ingest), so collision-renamed files surface in the right bucket instead of dropping into "File Not Found." Added an explicit resolution rule for multi-branch matches: later in the pipeline wins (`processed` > `capture` > `inbox`), so a single REQ never appears in two lesson sections.

## 0.72.1 — The Follow-On Four (2026-05-09)

A second-round review caught four follow-on bugs from 0.72.0. Two were narrow scoping mistakes (migration check only on one entry point, exit branches that excluded mixed cases), two were paths I didn't follow deep enough into bkb's directory layout. All four fixed.

- `actions/interview.md`: hoisted the v1→v2 migration into a shared **Session-Load Protocol** at the top of the action. Every subcommand that reads `session.json` (`<template>` resume, `status`, `review`, `export`, `ingest`) now invokes it before any other read, so an updated v1.x session can't bypass migration by entering through a non-bare subcommand. Step 2 now references the protocol instead of duplicating it.
- `actions/work.md`: Step 1's exit paths and Step 10's loop-or-exit are now **composed from sections** (completed/done, pending-answers, blocked-archive-collision) rather than disjoint "only X" branches. A queue with pending-answers + blocked-archive-collision (the gap the reviewer found) now renders both sections in one report instead of falling through into "no REQs at all."
- `actions/roadmap.md` Step 3: `kb_entry` lookup is now **recursive** under each `<kb>/raw/` branch. The previous top-level glob missed `raw/capture/<type>/` (triage's type subdirs) and `raw/processed/YYYY-MM-DD/` (ingest's date subdirs) — exactly the cases the new buckets were added to handle. Spelled out as `find <kb>/raw/<branch> -name <kb_entry>` with equivalent recursive globs.
- `actions/roadmap.md` Output Format: added rendering sections for the new buckets — **Lessons Promoted (Awaiting Ingest)**, **Lessons Processed (Terminal)**, **Lessons File Not Found** — so a `kb_entry` in `raw/capture/` or `raw/processed/` lands in the right section with the right next-step suggestion (`bkb ingest` for capture, no action for processed, investigate for missing) instead of falling back to the awaiting-triage section.

## 0.72.0 — The Five Patches (2026-05-09)

A review pass turned up five issues across capture, work, roadmap, the interview action, and the prompt-library README. All accepted, all fixed in one batch. The schema addition to interview sessions (a `template_version` field) is what bumps this to a minor.

- `actions/capture.md`: tightened the new TDD-on heuristic. `tdd: true` now requires that a *runnable* failing test can realistically be written first in the project's existing harness — not just a describable RED case. Manual/prompt-only proofs go in `## Red-Green Proof` with `tdd: false` instead, so capture stops creating REQs the work loop's mandatory test-first gate can't complete.
- `actions/work.md`: queue summary, Step 1 exit branches, and Step 10 loop-or-exit now account for `blocked-archive-collision`. Held duplicates are listed with their archived twin and a recovery instruction instead of disappearing into the silence between "no pending" and "no REQs at all."
- `actions/roadmap.md` Step 3: `kb_status: promoted` is a one-way stamp — the file moves through `raw/inbox/` → `raw/capture/` → `raw/processed/`. Roadmap now globs `kb_entry` across all three locations and buckets accordingly (awaiting triage / mid-pipeline / processed / not-found), so already-processed lessons stop showing up as actionable.
- `actions/interview-reference.md` + `actions/interview.md` + `interviews/work-operating-model.md`: added `template_version` to the `session.json` schema, the new-session write path (Step 1), all three re-run modes (`fresh`, `version`, plus the `update` shape via reference), and a new Step 2 migration check that auto-runs the template's documented "Migration from vX.x" steps. The work-operating-model migration text is now actionable instead of pointing at a phantom field.
- `prompts/README.md`: documented the exact-alias resolution tier the dispatcher already supports, so users can actually invoke aliases like `adr` / `adr-log` / `decisions` from the README's instructions.

## 0.71.2 — The TDD Default (2026-05-09)

Capture now defaults `tdd: true` instead of `tdd: false`. Most behavior-changing work benefits from a RED/GREEN cycle, so the bar is now "turn it off when it doesn't fit" rather than "turn it on when it clearly applies."

- `actions/capture.md` Step 1 TDD assessment: flipped default to true and rewrote the heuristic. Lists the narrow set where `tdd: false` is reasonable (pure styling/layout, copy/content, config bumps, doc-only, explicit throwaway spikes, no definable RED state).
- `actions/capture.md` Simple REQ frontmatter: `tdd: true` with a comment that flipping it off needs a real reason.

## 0.71.1 — The Deferred Link (2026-05-07)

The work action used to write prime-file lessons links from Step 7.5 — before Step 8 actually moved the REQ to its archive path, so the existence-verify either failed or the agent linked to the transient `working/` location. And nothing stopped a duplicate queue file from being silently re-processed when its twin was already archived. Both fixed.

- `actions/work.md` Step 7.5: prime-link writes are now COLLECTED as pending operations; the actual append + existence-verify happens in Step 8 substep 7, after the archive move.
- `actions/work.md` Step 8: new substep 7 walks the pending prime-link writes against the actual archived path.
- `actions/work.md` Step 2.0 (new): pre-claim glob check against `do-work/archive/**/REQ-NNN-*.md` AND `do-work/archive/**/REQ-NNN.md`. Bails cleanly with a clear message if the REQ id is already archived, and sets `status: blocked-archive-collision` on the duplicate to prevent livelock. Minimal scope (single-orchestrator); no post-move verify or pre-commit collision guard added.

## 0.71.0 — The Sweep (2026-05-07)

A pass through review findings: stale references, drifting pointers, parallel actions that resolved paths differently, a missing guide, and a template that mixed mechanical handlebars with natural-language directives. Plus a real semver fix on the work-operating-model template.

- `interviews/work-operating-model.md`: bumped to **2.0.0** (breaking) — `details.interruptions` is now `list[{source, priority}]` and `details.time_windows` requires a `days` array. Added a **Migration from v1.x** section with hand-migration steps for in-flight `session.json` files. Previously these schema changes shipped under a 1.1.0 minor bump.
- `interviews/work-operating-model.md`: declared the Export Templates dialect explicitly (handlebars-style with `where`/`sorted by` extensions and explicit `{{derived.<name>}}` slots). Replaced every natural-language directive embedded in `{{ … }}` (synthesis paragraphs, "items appearing in 2+", `{{#for each}}`, etc.) with named derived fields and per-template **Synthesized fields** blocks that say how each is computed.
- `actions/install-bowser.md` + `actions/install-ui-design.md`: both install actions now resolve project root the same way (`git rev-parse --show-toplevel || pwd`). Fixed an internal bowser inconsistency where Step 1 was cwd-relative while Step 4 was project-root-relative — these would mismatch when invoked from a subdirectory.
- `decisions/topics/_index_skill-architecture.md`: `sources:` frontmatter pointed at the deleted `actions/build-knowledge-base.md`; replaced with `actions/bkb.md` and `actions/bkb-reference.md`.
- `actions/interview.md`: Step 2 of `export` no longer claims `interview-reference.md` has a per-export schema list — the parenthetical now points readers at the template file's `## Export Templates` section directly, matching the reference file's actual content.
- `docs/forensics-guide.md`: added a "sister action" pointer to roadmap so the broken-vs-intended split is discoverable from either side.
- `docs/roadmap-guide.md`: new — every other first-class action had a guide except roadmap.
- `actions/clarify.md`: promoted "always show the builder's recommended choice" from a Rules-line into Step 3, where the verification checklist already asserted it.
- `decisions/imported-specs/2026-04-16_expand-skill-do-work-interview.md`: added a one-line footer noting that `actions/build-knowledge-base.md` was later split into `actions/bkb.md` + `actions/bkb-reference.md` so anyone reading the imported spec doesn't chase a deleted file.

---

## 0.70.5 — The Two Buckets (2026-05-07)

Two review findings fixed: roadmap's `kb_status: pending` recovery instruction was wrong (it pointed at `bkb triage`, but pending means no file was ever staged), and prompt aliases declared in headers (`dca`, `clg`, `cg`, `adr`, etc.) were unreachable because the dispatcher only resolved by filename.

- `actions/roadmap.md`: Step 3 rollup now splits `kb_status: promoted` (file staged → `bkb triage` + `bkb ingest`) from `kb_status: pending` (nothing staged → re-run handoff via `do-work review REQ-NNN`, possibly after `bkb init`). Output Format replaces the single "Lessons Awaiting Promotion" section with two distinct sections.
- `next-steps.md`: "After roadmap" block now suggests `bkb triage` only when promoted lessons exist, and points at `do-work review REQ-NNN` for pending lessons.
- `actions/prompts.md`: Resolution rules now include alias matching (priority 2, between exact filename and prefix). Aliases parsed from each prompt's `**Aliases:**` header line. Cross-file alias collisions are surfaced rather than silently picking one. `list` output gains an ALIASES column and warns on collisions.

---

## 0.70.4 — The Composed Key (2026-05-06)

Bare `status` and space-form `queue status` removed from the roadmap route — they caused first-match conflicts with any `<action> status` sub-command (interview, bkb, etc.). Use `do-work roadmap` or `do-work queue-status` (hyphenated) instead.

- SKILL.md routing table row 17: removed `do-work status` and `do-work queue status` examples
- Verb Reference roadmap entry: removed `status` and `queue status` triggers, replaced single-action exceptions with a general `"<action> status" → that action` rule

---

## 0.70.3 — The Wired Roadmap (2026-05-06)

`roadmap` was drafted but unrouted — the dispatch table didn't list it, the help menu didn't mention it, and no other action ever suggested it. Now it's wired end-to-end so users can actually find and run it.

- `SKILL.md`: Added roadmap to the Actions list, frontmatter argument-hint, routing priority table (priority 17, triggered by `roadmap`, `queue-status`, `status`, `where are we`, `what's left`, `what's feasible`, `what should I work on next`), Verb Reference, Help Menu, Action Dispatch table, and the foreground subagent list.
- `next-steps.md`: New `After roadmap:` block. Surfaced `do-work roadmap` as a follow-up after `forensics`, `verify requests`, and `work` so users discover it in flow.
- `README.md`: New section 18 covering the roadmap action with example invocations; sections 19-24 renumbered.
- `actions/forensics.md`: Added a "Do NOT use when" pointer to roadmap to clarify the broken-vs-intended split between the two read-only surveys.

## 0.70.2 — The TDD Telltale (2026-05-06)

The roadmap action now reads `tdd` posture per REQ and flags pending items where TDD is off but the behavior is testable — so reviewers can decide to flip it on before pickup. Also picks up `queue-status` as an explicit trigger phrase.

- `actions/roadmap.md`: New Step 2.5 classifies pending REQs as TDD on / eligible / not applicable, with evidence (frontmatter, `## Red-Green Proof`, domain, input/output examples). Output Format adds a `TDD Eligible` section and per-row `tdd:` annotations across Ready / Needs Clarification / In Progress / Recently Completed.
- `actions/roadmap.md`: Added `queue-status` and `queue status` to the When-to-Use trigger phrases. Also added a rationalization, two red flags, and a verification-checklist item for TDD reporting.

## 0.70.1 — The Lookahead Lens (2026-05-05)

Drafted a new `roadmap` action — a read-only survey of the do-work queue that classifies pending REQs as ready / needs-clarification / blocked / stale and rolls up in-progress and recently-completed work. Sits alongside `forensics` (which finds *broken* state) by reporting *intended* state and feasibility instead. Not yet wired into SKILL.md routing.

- New `actions/roadmap.md` covering pending-feasibility classification, in-progress reporting, completed-work roll-up by UR, and a "Suggested Next Steps" punch list.
- Explicit boundaries against `forensics`, `scan-ideas`, `clarify`, and `inspect` so routing stays clean once the action is registered.

## 0.70.0 — The Karpathy Echo (2026-05-04)

Karpathy guardrails were already auto-loaded at implementation, but the principles were invisible elsewhere — specs didn't cite them, entry-point docs didn't name them, and the upstream's verifiable-goals examples never made it into our adaptation. This release surfaces the four principles across the funnel without spamming citations, and backfills the dropped content.

- `crew-members/karpathy.md`: Backfilled upstream's transformation examples and multi-step plan template under Goal-Driven Execution.
- `specs/api-endpoint.md`, `specs/ui-component.md`, `specs/refactor.md`, `specs/bug-fix.md`: Added one-line guardrail citations to each Quality Standards section.
- `actions/review-work.md`: Added principle→dimension cross-reference table after the Karpathy Principle Check.
- `actions/capture.md`, `actions/clarify.md`: Connected the `- [~]` open-question convention to the "Think Before Coding" guardrail.
- `SKILL.md`, `README.md`: Named the four principles in the entry-point docs.

## 0.69.17 — The Thin Week Allowance (2026-04-27)

Resolves an internal contradiction in the `weekly-signal-diff` Verification checklist for the new "Top of mind this week" section. The section spec allows fewer bullets when the week is thin ("give fewer bullets rather than padding"), but the checklist required 3–5 bullets — so a compliant thin-week output would fail self-check or get padded with filler. Codex flagged it on PR #96.

- `prompts/weekly-signal-diff.md`: Verification checklist for "Top of mind this week" now enforces only the upper bound (≤5 bullets, ≤150 words) and explicitly permits fewer bullets in thin weeks.

## 0.69.16 — The Archetype Bullet (2026-04-27)

Adds a per-shift "For client archetypes" bullet to every headline structural shift in the `weekly-signal-diff` digest. Naming the archetype and a one-line outreach angle inside the shift itself — kept visually separate from "Why it matters to this user" — turns each shift into a scannable outreach prompt instead of a synthesis the operator has to redo at the desk.

- `prompts/weekly-signal-diff.md`: New `**For client archetypes**` bullet inserted after `**Why it matters to this user**` in the headline-shift template. Optional per-shift; "No direct client angle" is the explicit empty form. New Common Rationalizations row blocks the "obvious from context" shortcut.

## 0.69.15 — The Action Split (2026-04-27)

Promotes Actions from an optional tail section to a mandatory block at the head of the `weekly-signal-diff` digest, and splits it into two groups: operator-facing captures and proactive client-outreach angles. Pushes the digest's value outward toward the operator's clients, not just inward toward the operator's own backlog.

- `prompts/weekly-signal-diff.md`: Removed `### Actions (optional)` from the bottom of Phase 7. Added `### Actions this week` between "Top of mind this week" and "Coverage note" with two mandatory groups. Empty groups must be stated explicitly — silence isn't allowed. Matching Rule and Verification checklist entry added.

## 0.69.14 — The Top Of Mind (2026-04-27)

Adds a mandatory "Top of mind this week" subsection at the head of the `weekly-signal-diff` digest. Forces the agent to lead with the 3–5 things the operator should hold in working memory — synthesis, not detail — so the rest of the digest reads as support material for mid-week re-reading. Hard cap of 5 bullets / 150 words; thin weeks shrink the bullet count rather than padding.

- `prompts/weekly-signal-diff.md`: New `### Top of mind this week` block in Phase 7, placed before `### Coverage note`. Matching Rule and Verification checklist entry added so the cap is enforceable, not advisory.

## 0.69.13 — The Symmetry Patch (2026-04-23)

Closes five findings from a contradictions-and-gaps sweep of the repo. Main move: the bkb action's filename now matches its trigger word, so every action follows the same naming rule.

- Renamed `actions/build-knowledge-base.md` → `actions/bkb.md`; updated live cross-refs in `SKILL.md`, `CLAUDE.md`, `bkb-reference.md`, `kb-lessons-handoff.md`, and `prompts/architecture-decisions-log_create-or-expand.md` (historical references in `decisions/` and CHANGELOG preserved).
- `CLAUDE.md` Project Structure now lists `_dev/` and `decisions/` — both tracked directories that were absent from the tree.
- `CLAUDE.md` Agent Rules now documents `interviewer.md` (loaded by the interview action across all sub-commands).
- `CLAUDE.md` docs exemption now covers `kb-lessons-handoff` explicitly as a reference-only action invoked by other actions.
- `actions/work.md` Request File Schema now includes the `caveman` frontmatter field with its intensity values (`lite` | `full` | `ultra`).

## 0.69.12 — The Dark Code Kit (2026-04-23)

Captures a three-prompt kit for fighting "dark code" — code that was never understood by anyone at any point in its lifecycle. Shared `dark-code-kit_` prefix groups them as sibling tools in the library.

- `prompts/dark-code-kit_audit.md`: four-group interview (architecture, AI tool usage, team/ownership, deployment) that produces a hotspot map across structural and velocity dimensions, with severity ratings, ownership gaps, and a prioritized action plan.
- `prompts/dark-code-kit_context-layer-generator.md`: per-module interview walking through structural → semantic → philosophical context, emitting a module manifest, behavioral contracts, and a decision log that make the module self-describing.
- `prompts/dark-code-kit_comprehension-gate.md`: senior-engineer-style PR review across seven dimensions (credentials, cross-service side effects, blast radius, state, tokens, assumptions, comprehension) with CLEAR / REVIEW REQUIRED / HOLD verdicts.
- `prompts/README.md`: three entries added to the Available prompts table.

## 0.69.11 — The Ingest Correction (2026-04-23)

Fixes the kb-lessons handoff's user-facing messages: both the no-KB fallback and the promoted confirmation told users to run `bkb triage` alone, but triage only sorts inbox files — compilation into the wiki happens in `bkb ingest`. Following the old messages left lessons stuck in `capture/notes/` and invisible in the wiki.

- `actions/kb-lessons-handoff.md`: no-KB fallback now documents the full re-promotion path — `bkb init` → re-run handoff (e.g. `do-work review REQ-NNN`) → `bkb triage` → `bkb ingest`. Previously it stopped at triage and also glossed over the fact that the handoff set `kb_status: pending` without dropping the file, so even a correct triage+ingest pair would have found an empty inbox.
- `actions/kb-lessons-handoff.md`: "Promoted to …" confirmation now instructs `bkb triage` then `bkb ingest`. Previously users on the happy path were told `bkb triage` was the last step, leaving the lesson sorted but uncompiled.

## 0.69.10 — The Gap Patrol (2026-04-23)

Audit-driven cleanup of the three recent handoff commits (0.69.7–0.69.9). Fills in the spots where the new `kb_status`/`kb_entry` fields and the handoff flow weren't yet mentioned in sibling docs. Nothing behavioral — just the cross-references finally catching up with the feature.

- `actions/work.md`: `## Request File Schema` now documents the two optional `kb_status` and `kb_entry` frontmatter fields alongside the existing ones. Previously only `sample-archived-req.md` mentioned them, so agents reading the schema block thought they were non-standard.
- `next-steps.md`: "After work" and "After review work" blocks now suggest `do-work bkb triage` as a follow-up when lessons were promoted, and `do-work bkb init` when the handoff deferred because no `kb/` existed.
- `actions/build-knowledge-base.md`: `triage` classification table now recognizes `.md` files with `source_type: req_lesson` frontmatter (written by the kb-lessons handoff). They route to `capture/notes/` — no new capture subdir needed — with a note that the `domain` field is a reliable topic hint and `req_path` is a back-reference to the originating REQ.

## 0.69.9 — The Handoff Cleanup (2026-04-23)

Two bot-reviewer findings against the kb-lessons handoff, both legitimate and both fixed. Metadata now populates correctly in pipeline mode, and the `declined` vs `skipped` statuses are actually reachable as designed.

- `actions/kb-lessons-handoff.md`: `date` now falls back to today's date when `completed_at` isn't set yet — the handoff runs at Step 7.5 (pipeline mode), before Step 8 writes `completed_at`, so the old "source from `completed_at`" rule produced empty dates on every pipeline run.
- `actions/kb-lessons-handoff.md`: user's explicit "Skip" choice in Step 3/4 now records `kb_status: declined` instead of `skipped`, matching Step 5's semantics (`declined` = active refusal, `skipped` = silent auto-skip when trigger conditions aren't met). Previously `declined` was effectively unreachable.

## 0.69.8 — The Homegrown Handoff (2026-04-23)

Replaces the compound-engineering integration from 0.69.7 with a zero-dependency version that uses do-work's own knowledge base (`kb/`). After a REQ's review passes and Lessons Learned are captured, do-work drops a structured source document into `kb/raw/inbox/` so the existing `bkb triage` → `bkb ingest` pipeline compiles it into the wiki. Same consent-driven shape as before, just no external plugin required.

- `actions/kb-lessons-handoff.md`: New handoff reference. Writes to `<kb>/raw/inbox/REQ-NNN-<slug>.md`, defers to `kb_status: pending` if no `kb/` exists (never auto-inits), and stops at the drop — triage and ingest stay in the bkb action's lane.
- `actions/review-work.md`, `actions/work.md`: Step 9.5 / Step 7.5 now call the kb-lessons handoff instead of the CE one. Unattended runs default to `kb_status: pending`.
- `actions/sample-archived-req.md`: Frontmatter fields renamed — `ce_compound_status` → `kb_status`, `ce_solution_path` → `kb_entry` (filename only, survives bkb's moves through `capture/` and `processed/`).
- `CLAUDE.md`: "Compound-engineering Integration" section replaced with a shorter "Lessons → Knowledge Base Handoff" section that documents the in-skill flow only.
- Removed: `actions/ce-compound-handoff.md`, `docs/ce-integration-guide.md` — both were CE-specific and no longer apply.

## 0.69.7 — The Compound Handoff (2026-04-23)

First integration point with the [compound-engineering plugin](https://github.com/EveryInc/compound-engineering-plugin). After a REQ's review passes and Lessons Learned are captured, do-work now offers to promote those lessons into CE's `docs/solutions/` knowledge base via the `ce-compound` skill. The handoff asks before dispatching, degrades to a saved prompt if CE isn't installed, and never blocks archival.

- `actions/ce-compound-handoff.md`: New reference file describing the handoff payload shape, user consent flow, and REQ frontmatter updates. Both review-work Step 9.5 (standalone) and work Step 7.5 (pipeline) dispatch into this single reference.
- `actions/review-work.md`: Step 9.5 now runs the compound handoff after lesson capture in standalone mode.
- `actions/work.md`: Step 7.5 now runs the compound handoff after lesson capture in pipeline mode. Unattended runs default to `ce_compound_status: pending` — no auto-promotion.
- `actions/sample-archived-req.md`: Sample frontmatter now shows the two new optional fields (`ce_compound_status`, `ce_solution_path`) so REQ authors know the schema.
- `CLAUDE.md`: New "Compound-engineering Integration" section documents the augmentation model, the three CE artifact paths, and the current integration point.
- `docs/ce-integration-guide.md`: New user-facing guide covering install, the handoff flow with sample payload, troubleshooting, roadmap for future integration points (reviewer agents, ce-plan, ce-brainstorm), and design principles for contributors wiring up the next seam.

## 0.69.6 — The Audit Ratchet (2026-04-22)

Close the contradictions and gaps found in a self-audit of the skill: a broken link, a missing `next-steps.md` entry, an out-of-date README, a missing docs guide, two action files that didn't follow the template, and a wave of missing `When to Use` / `Red Flags` / `Verification Checklist` sections across core actions. Nothing behavioral — just the docs finally matching the conventions CLAUDE.md claims.

- `decisions/records/adr-012-interview-v2-gap-closure.md`: Fixed broken link to v1 spec — the filename is date-prefixed (`2026-04-16_expand-skill-do-work-interview.md`).
- `next-steps.md`: New `After interview` blocks covering session-in-progress, all-layers-complete, export, and list — previously absent despite `interview` being a first-class action.
- `README.md`: `bkb` usage list now includes `defrag`, `garden`, and `crew [action]` — all three were already in the action file and `next-steps.md`, just missing from the README overview.
- `CLAUDE.md`: `prompts/` tree entry now points at `prompts/README.md` as the authoritative index instead of listing one outdated prompt.
- `docs/prompts-guide.md`: New guide for the prompts dispatcher — sub-commands, name resolution, safety model, and how to add a new prompt.
- `actions/install-ui-design.md`, `actions/install-bowser.md`: Restructured to follow the CLAUDE.md action template (When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist).
- `actions/capture.md`, `actions/clarify.md`, `actions/work.md`, `actions/pipeline.md`, `actions/ui-review.md`, `actions/prompts.md`, `actions/present-work.md`, `actions/prime.md`, `actions/version.md`, `actions/tutorial.md`, `actions/build-knowledge-base.md`, `actions/forensics.md`, `actions/deep-explore.md`, `actions/scan-ideas.md`: Added missing `When to Use`, `Red Flags`, and/or `Verification Checklist` sections per CLAUDE.md's action-template spec. All 14 core actions now carry the full template.

## 0.69.5 — The Hyphen Hustle (2026-04-22)

Every `do work` command invocation is now written `do-work` across docs, actions, crew rules, and the session-start hook. Matches the skill's actual name and makes it unambiguous to agents that it's a real command, not a verb phrase.

- All `*.md` files and `hooks/session-start.sh`: `do work <action>` → `do-work <action>`, including README examples, SKILL.md routing tables, action files, docs, CHANGELOG prose, crew rules, prompts, and decision records.
- No behavior change — natural-language triggering still works; the skill's name has always been `do-work`, so hyphenated references stay consistent with the skill manifest.

## 0.69.4 — The Review Ratchet (2026-04-17)

Follow-up to 0.68.2: fixes three defects from code review on the interview v2 gap-closure patch. One was a JSON rendering bug, one was a reference to a session field that doesn't exist, and one was a stale-entry leak into agent rules that violated ADR-012's own promise. ADR-012 gets a "Post-merge corrections" section documenting each.

- `interviews/work-operating-model.md`: `operating-model.json` template now uses `{{json_entries <layer>}}` instead of `[ "{{canonical_entries}}" ]` — emits a proper JSON array of entry objects instead of a single-element array of strings.
- `interviews/work-operating-model.md`: All `{{session.completed_at}}` references changed to `{{session.last_exported_at}}` (the field that actually exists on `session.json`).
- `actions/interview.md`: `export` sub-command reordered to stamp `last_exported_at` in-memory **before** rendering (step 2), then persist `session.json` after artifacts are on disk (step 4). Prevents templates from substituting a null timestamp on first export.
- `interviews/work-operating-model.md`: SOUL.md and HEARTBEAT.md templates now filter `where status != "stale"` on every entry-iterating block. USER.md's active sections do the same, plus a new "Stale or deprecated" section labels stale entries at the bottom (narrative context preserved, but they no longer appear as active rules).
- `actions/interview-reference.md`: Ingest frontmatter `created:` fields follow the template fix (`last_exported_at` in place of `completed_at`).
- `decisions/records/adr-012-interview-v2-gap-closure.md`: "Post-merge corrections" section added under Consequences.

## 0.69.3 — The Honored Flag (2026-04-17)

Fixes an inconsistency in the eval-harness prompt flagged in code review: `--tasks <n>` was documented but the interview and output flow were hard-coded to exactly three test cases. The prompt now resolves N from the flag up front (default 3, clamped to 1–7) and uses N everywhere — task inventory, priority selection, case count, verification.

- `prompts/prompt-kit-step5-eval-harness.md`: new Step 0 resolves and clamps N; Steps 1, 2, 3, 5 reference N instead of literal 3; Rules and Verification Checklist enforce the contract; Red Flags call out suite-size drift; template placeholder for the per-case index changed from `[N]` to `[#]` to avoid visual collision with the count variable.

## 0.69.2 — The Topical Shelving (2026-04-17)

Regroups the five AI-industry analytical prompts by the discipline they're drawn from — business, economics, or tech — dropping the redundant `ai-` umbrella (the whole library is AI-oriented). One of the tech prompts gains an `infrastructure` sub-prefix to mark it as an infra decision rather than an architecture one.

- `prompts/ai-vendor-strategic-sort.md` → `prompts/business-vendor-strategic-sort.md`
- `prompts/inference-economics-stress-test.md` → `prompts/economics-inference-stress-test.md`
- `prompts/saas-repricing-exposure.md` → `prompts/economics-saas-repricing-exposure.md`
- `prompts/compute-geography-risk.md` → `prompts/tech-infrastructure-compute-geography-risk.md`
- `prompts/inference-architecture-decision.md` → `prompts/tech-inference-architecture-decision.md`
- `prompts/README.md` index rows updated to match. Historical references in `CHANGELOG.md` left as-is.

## 0.69.1 — The Spelled-Out Name (2026-04-17)

Renames the ADR-log prompt so its filename actually says what it does. Establishes a `[noun]_[action]` convention (underscore between the subject and the verb phrase) that leaves room for sibling actions on the same noun later.

- `prompts/adr-log.md` → `prompts/architecture-decisions-log_create-or-expand.md`: renamed; H1 and aliases updated inside the file (`adr`, `adr-log`, `decisions`, `architecture-decisions` all still work as documentation hints — the dispatcher resolves via prefix match against the new filename).
- Cross-references updated in `SKILL.md`, `CLAUDE.md`, `README.md`, `actions/prompts.md`, and `prompts/README.md`. Historical references in `CHANGELOG.md` left as-is.

## 0.69.0 — The Seven Steps (2026-04-17)

Extracts the Prompt Kit article's progression into the library as seven numbered prompts. One pre-flight pen-and-paper exercise plus six runnable disciplines — diagnostic, context doc, spec engineer, intent framework, eval harness, constraints — all `step[n]`-prefixed so they sort in workflow order.

- `prompts/prompt-kit-step0-pen-and-paper-exercises-to-prepare-prompt.md`: handoff prompt that tells the user to step away from the screen and work the seven questions offline, then structures the returning notes into a PRE-FLIGHT BRIEF.
- `prompts/prompt-kit-step1-four-discipline-diagnostic.md`: scored audit across Prompt Craft, Context, Intent, Specification — with a 4-month personalized roadmap.
- `prompts/prompt-kit-step2-personal-context-doc.md`: seven-domain interview producing the user's "CLAUDE.md for everything."
- `prompts/prompt-kit-step3-spec-engineer.md`: collaborative spec builder for real projects — acceptance criteria, constraint architecture, task decomposition, definition of done.
- `prompts/prompt-kit-step4-intent-and-delegation-framework.md`: extracts implicit decision rules into a deployable framework, with a Klarna Test self-check.
- `prompts/prompt-kit-step5-eval-harness.md`: Lütke-pattern test suite over the user's actual recurring tasks.
- `prompts/prompt-kit-step6-constraint-architecture.md`: pre-delegation Must Do / Must Not / Prefer / Escalate document tied to the user's stated failure modes.
- `prompts/README.md`: index updated with all seven new entries.

## 0.68.2 — The Paved Cowpath (2026-04-17)

Closes five v1 gaps in the `interview` action per the v2 imported spec — export templates move into the template file as mechanical render templates, `update` goes entry-level, mid-layer quits become recoverable, and `ingest` lands 10 files in `kb/raw/inbox/` instead of inventing its own frontmatter shape. Surgical patches, not a rewrite. Recorded as ADR-012.

- `interviews/work-operating-model.md`: New `## Export Templates` section with verbatim handlebars-style render templates for `USER.md`, `SOUL.md`, `HEARTBEAT.md`, `operating-model.json`, and `schedule-recommendations.json`. An implementation can now render exports mechanically against the approved session — different runs produce the same file shape.
- `actions/interview-reference.md`: `## Export Schemas` trimmed to framework-level invariants only (narrative tone, source-confidence filtering, cadence, traceability). Template-specific rendering now lives in the template.
- `actions/interview-reference.md`: `update` re-run mode rewritten to walk entries individually — `[confirm / edit / mark-stale / delete / skip]` per entry. Explicitly overrides v1's "do not invent a per-entry patch path." CHANGELOG format for update runs is now `N confirmed, N edited, N marked stale, N deleted, N added`.
- `actions/interview-reference.md`: New `### Mid-layer recovery` section. On resume, the action checks for `.draft-<layer-id>.md` written opportunistically during the interview and offers pick-up vs. start-over.
- `actions/interview-reference.md`: `## Ingest Frontmatter` rewritten as `## Ingest File Mapping`. Specifies 5 export files + 5 layer summaries = 10 files per run for `work-operating-model`, plus a manifest row per file in `kb/raw/_inbox_queue.md`. Frontmatter aligns with BKB's canonical schema (`sources:` list, `related:` with `rel`, `type: source-summary` for exports, `type: concept` for layer summaries).
- `actions/interview.md`: New draft-checkpoint step in the layer interview workflow. Subsequent steps renumbered. `ingest` sub-command body rewritten to reference the new File Mapping section in the reference.
- `decisions/records/adr-012-interview-v2-gap-closure.md`: New ADR documenting the five patches. Extends ADR-011. Crew placement audit confirmed `crew-members/interviewer.md` stays put — the directory is a generic persona pool, not `work`-scoped.
- `decisions/_master_index.md` + `decisions/topics/_index_skill-architecture.md`: Bumped to list ADR-012.

## 0.68.1 — The Rename Tag (2026-04-16)

Renames the Weekly Structural Diff prompt so "original" is explicit in the filename — clears the way for variant versions of the same framework to coexist in the library.

- `prompts/weekly-structural-diff.md` → `prompts/weekly-structural-diff-original.md`: renamed; index entry in `prompts/README.md` updated to match. Invoke with `do-work prompts run weekly-structural-diff-original` (prefix match `weekly-structural-diff` still resolves unambiguously while it's the only variant).

## 0.68.0 — The Promptkit Drop (2026-04-16)

Six new reusable prompts ingested from the Prompt Kit article on the 2026 capability-phase → economics-phase transition. They turn the article's analytical framework into runnable tools for tracking AI news, stress-testing product economics, mapping infrastructure risk, pricing SaaS seat compression, sorting vendors, and designing inference architectures.

- `prompts/weekly-structural-diff.md`: Signal/noise sort for AI news across five altitudes (physics, monetization, geography, business models, geopolitics), with a "what didn't change" calibration and prioritized takeaways.
- `prompts/inference-economics-stress-test.md`: Sora-style economics stress test — sustainability ratio, three-scenario pressure test, emoji verdict (🟢/🟡/🟠/🔴), and a concrete "what would fix it" plan. Benefits from a thinking-capable model.
- `prompts/compute-geography-risk.md`: Location-by-location risk matrix (power/grid, permitting/politics, geopolitics, data residency) with a deployment strategy and contingency playbook per location.
- `prompts/saas-repricing-exposure.md`: Seat compression estimate, "The Clock" (months until compression shows in reported numbers), transition readiness score, and an Atlassian benchmark.
- `prompts/ai-vendor-strategic-sort.md`: Vendor assessment matrix across five structural-sustainability dimensions, one tripwire event per vendor, and a portfolio concentration score.
- `prompts/inference-architecture-decision.md`: API vs. self-hosted vs. hybrid comparison, model selection matrix, Sora test, and a Now / 3× / 10× migration path with triggers.
- `prompts/README.md`: Index table extended with the six new prompts.

## 0.67.5 — The Weekly Witness (2026-04-17)

New prompt in the library: `weekly-signal-diff` — a weekly structural diff of AI-industry news, personalized via BKB. Ships with a 10-lane core starter universe and auto-loads a personal sidecar at `prompts/weekly-signal-diff-personal.md` when present for user-specific lanes. Every loaded lane gets full coverage every week — no lane is ever compressed or dropped.

- `prompts/weekly-signal-diff.md`: New prompt. Produces both an inline digest and a durable deliverable at `do-work/deliverables/weekly-signal-diff/<week-ending>.md` staged for BKB ingest. Idempotent per week-ending date (appends timestamped revisions rather than overwriting). Supports `--week-ending`, `--source-packet`, `--topic`, `--dry-run`, `--no-ingest`. Aliases: `wsd`, `signal-diff`.
- `prompts/weekly-signal-diff-personal.md`: New placeholder template. Ships with no real lanes — users copy it anywhere in their project (project root, `.claude/`, `do-work/`, etc.) and fill in real lanes. At Phase 3 the main prompt searches the user's project and loads whatever project-local copy it finds; the shipped placeholder is only a template, never treated as a source of real lanes. Library prompt and shipped placeholder stay generic; personal content lives exclusively in the user's project.
- `prompts/README.md`: New rows for `weekly-signal-diff` and `weekly-signal-diff-personal` in the Available prompts table.
- `decisions/imported-specs/2026-04-16_weekly-signal-diff-authoring-prompt.md`, `decisions/imported-specs/2026-04-17_starter-universe.md`: Spec updates — demotion language removed, 3–7 shift cap removed, forbidden-memory-layer name-drops stripped, personal sidecar pattern documented.

## 0.67.4 — The Gap Sealer (2026-04-16)

Folds in the legitimate improvements from a parallel branch that landed alongside 0.67.2/0.67.3. The earlier "Unified Trunk" merge tried to combine both lines but truncated `CHANGELOG.md` and rewrote `actions/version.md` losing the global-install guard and the recap section — that merge was reverted and only the load-bearing changes were re-applied here.

- `actions/version.md`: Widened the auto-update dirty check to scope every shipped editable path (`prompts/`, `interviews/`, `specs/`, `docs/`, `decisions/`, `hooks/`, `CLAUDE.md`, `AGENTS.md`, `next-steps.md`) — anything tar would clobber. Anything dirty in those paths now blocks the update.
- `actions/version.md`: New pre-clean step (4) for `prompts/` and `interviews/` — top-level `.md` files are deleted before extraction so upstream-removed entries don't linger as ghost workflows in `do-work prompts list` / `do-work interview list`. Subsequent steps renumbered 4→5, 5→6, 6→7.
- `actions/interview-reference.md`: `update` re-run mode now tracks an in-memory `any_edits` flag. If any layer's approval committed a non-zero diff, the export gate state (`review_completed_at`, `review_runs`) is cleared on completion — the user must re-run `review` before the next `export`. Pure re-confirms leave the gate untouched.
- `actions/interview-reference.md`: `fresh` and `version` empty session shapes now include `last_activity_at: <now>` so the freshness preflight has something to compare against on the very first export.
- `actions/interview.md`: Exports gate rule documents that `update` clears the review state when edits are committed.
- `interviews/work-operating-model.md`: Layer 1 schema fix — `time_windows` entries gain a required `days` field (weekday abbreviations) so `schedule-recommendations.json` can emit `days` without inventing data; `interruptions` is now a list of `{source, priority}` objects (priority drawn from `low`/`medium`/`high`) so `HEARTBEAT.md`'s "What to ignore" section has a real signal to filter on. Template version bumped 1.0.0 → 1.1.0.

## 0.67.3 — The Right Shelf (2026-04-16)

Moves the 0.67.2 export freshness stamp out of `exports/` and into `session.json.last_exported_at`. The sidecar-file approach would have been picked up by `ingest`'s "for each file" loop and polluted `kb/raw/inbox/` with bogus timestamp documents. Caught in review; the field-on-session.json approach was always the right one.

- `actions/interview.md`: `export` preflight and stamp-write now read/write `session.json.last_exported_at` instead of a sidecar file. Empty session shape gains the new field.
- `actions/interview-reference.md`: `session.json` schema gains `last_exported_at`. Status Vocabulary row updated with a note explaining why the stamp lives on the session, not in `exports/`. `fresh` re-run mode writes `last_exported_at: null` in the new empty session.

## 0.67.2 — The Status Ledger (2026-04-16)

Interview recipe gains a stale-export warning and a consolidated status vocabulary — small operational patches for when an operating model gets re-run in anger. Addresses gaps surfaced by a recent design review of the `work-operating-model` activation path.

- `actions/interview.md`: `export` sub-command now stamps `exports/.exported_at` after each run and does a freshness preflight on the next run — if `session.json.last_activity_at` is newer than the stamp, the user hears about it before exports are regenerated.
- `actions/interview-reference.md`: New Status Vocabulary table consolidates the four independent status fields (session `status`, layer `approved`, entry `status`, export freshness stamp) into a single reference. Explicitly notes that prior runs are archived directories, not `superseded` flags.
- `actions/interview-reference.md`: `update` re-run mode now documents the "empty a layer" path (user can nuke a layer; same approval gate applies, empty layer still counts as approved) and calls out that per-entry edit friction is intentional — the approval gate is the whole point.

## 0.67.1 — The Settled Tenant (2026-04-16)

Interview action now works the moment the skill is installed into a project, and session state lives in `do-work/` alongside the rest of the per-repo workspace — tracked in git like URs and REQs.

- Templates resolve from `<skill-root>/interviews/` (the `interviews/` directory inside the skill bundle), not the user's project root. Fixes `do-work interview list` and `do-work interview <template>` finding nothing when the skill ships from `~/.claude/skills/do-work/`.
- Session state moved from `./interview/<template>/` to `./do-work/interview/<template>/`. It joins `queue/`, `user-requests/`, `archive/`, and `working/` under the canonical workspace and is tracked in git — the elicited operating model is durable per-repo knowledge, not transient orchestration state.
- Removed the stale `interview/` entry from the skill repo's own `.gitignore` so the skill no longer models the wrong behaviour.

## 0.67.0 — The Open Ear (2026-04-16)

New `interview` action — a generalized elicitation framework that runs prescriptive templates to turn tacit work knowledge into agent-ready operating artifacts. First template `work-operating-model` walks the five-layer Work Operating Model (Nate B. Jones and Jonathan Edwards) across ~45 focused minutes and produces `USER.md` / `SOUL.md` / `HEARTBEAT.md` plus machine-readable exports. Session state is resumable, cross-layer contradictions get surfaced explicitly, and exports flow into BKB via `ingest` for querying.

- `actions/interview.md`: New sub-command dispatcher — `list`, `<template>`, `<template> status`, `<template> review`, `<template> export`, `<template> ingest`, `<template> reset`, `<template> versions`. Session state lives at `./interview/<template>/session.json` and persists across sessions per ADR-005. Export gates on all layers approved + at least one review pass complete. Re-run modes (`fresh`, `update`, `version`) archive prior runs as immutable `versions/v<N>-<date>/` directories.
- `actions/interview-reference.md`: Companion per ADR-001 holding the heavy content — template file format, canonical 11-field entry contract, `session.json` schema (including `review_completed_at` + `review_runs` gate fields), checkpoint format, per-export schemas for the five `work-operating-model` artifacts, re-run mode specifications, versioning scheme, and ingest frontmatter shape.
- `interviews/work-operating-model.md`: First template. Five layers — operating rhythms, recurring decisions, dependencies, institutional knowledge, friction — each with concrete prompt patterns and layer-specific `details` shape. Declares four named cross-layer contradiction checks the `review` sub-command surfaces.
- `crew-members/interviewer.md`: New persona loaded during every interview sub-command. Concrete-before-abstract, one-question-at-a-time, checkpoint-gated, honest-confidence standards. Never invents fields the user didn't provide.
- `docs/interview-guide.md`: Onboarding guide — when to run (45-minute focused session), the five export files, re-run cadence (quarterly), BKB integration flow, multi-repo context separation, and troubleshooting.
- `decisions/records/adr-011-interview-framework-with-prescriptive-templates.md`: New ADR documenting the prescriptive-not-minimal template shape, single-instance-per-repo design, and local-files-only constraint. Depends on ADR-001, ADR-002, ADR-005; complements ADR-010.
- `SKILL.md`: Registered in action list, routing table (priority 19), Verb Reference, Action Dispatch table, bare-invocation help menu (new Interviews block), and foreground-dispatch list. Frontmatter `argument-hint` updated.
- `README.md`: New numbered scenario "19. Run a structured interview"; renumbered later scenarios 20→21, 21→22, 22→23.
- `decisions/_master_index.md`, `decisions/_progress.md`, `decisions/topics/_index_skill-architecture.md`: ADR-011 added to the index and topic cluster; progress tracker bumped to `Next ADR number: ADR-012`.
- `.gitignore`: New `interview/` line so per-repo session state isn't accidentally committed. Templates under `interviews/` remain tracked.

## 0.66.1 — The Local Landlord (2026-04-16)

The `do-work update` flow now refuses to overwrite a global/shared install. If `SKILL.md` lives under `~/.claude/skills/`, `~/.gemini/skills/`, or anywhere else outside the current project's git root, the update stops and tells the user to either `cd` into the owning project or install the skill locally — no more silent updates to a user-wide copy.

- `actions/version.md`: Added an explicit preflight location check as step 2 of the update flow that resolves the skill root, compares it to `git rev-parse --show-toplevel`, and refuses to proceed if the skill sits under a user-wide skills directory. Renumbered the dirty-tree / run / verify / report steps accordingly. The curl command is now prefixed with `cd <skill-root> &&` so extraction can't land in a global directory by mistake. The fetch-failed fallback message was rewritten to call out the global paths by name.

## 0.66.0 — The Four Corners (2026-04-16)

Deliverables now follow an unambiguous naming convention: `.marp.md` for LLM-authored Marp source, `.marp.html` for the marp-cli export of that source, and `.single.html` for LLM-authored standalone HTML (explainer or debrief). The pipeline now ships four files per completion — three LLM renderings plus the mechanical Marp HTML export — so a stakeholder without marp-cli can still view the deck.

- `actions/pipeline.md`: Step 5 table expanded to four rows — `.md`, `.marp.md`, `.marp.html`, `.single.html` — with the Marp HTML export marked as mechanically produced by `npx @marp-team/marp-cli ... --html`. Narrative, rationalizations, red flags, and the verification checklist updated to distinguish the three LLM renderings from the fourth mechanical export, and to scope the Tailwind/Mermaid CDN constraint specifically to `.single.html`.
- `actions/pipeline-reference.md`: Section 3 heading + filename renamed to `.single.html`. Section 2 (Marp Slide Deck) now calls out the `.marp.html` export with the exact command. Sibling-link lists, preview commands, and the HTML Related-deliverables card grid updated to link both `.marp.html` and `.single.html` where relevant.
- `actions/present-work.md`: Interactive explainer renamed to `{UR-NNN}-interactive-explainer.single.html` with a note explaining the `.single.` vs `.marp.` distinction. Client-brief "Related Reading" footer and the terminal summary updated to the new filenames, and the Keep-exploring footer now links both pipeline summary formats.

## 0.65.2 — The Dry-Run Reprieve (2026-04-15)

Fixes two review findings on the `adr-log` prompt. Phase 0 no longer hard-blocks every run on `main`/`master` — `--dry-run` now skips the tree/branch blockers entirely (they're zero-risk in a read-only run), and non-dry-run invocations on `main` pause and ask for authorization instead of refusing outright. README's description of the prompt's source model was stale; it now accurately reflects the layered spine (`implementation-history.md` primary, `lessons-learned/` secondary, code verification, `CHANGELOG.md` fallback).

- `prompts/adr-log.md`: Rewrote Phase 0 to parse flags first, skip dirty-tree / branch-name blockers under `--dry-run`, and prompt for authorization on `main`/`master` (with three accepted responses: yes / feature-branch-name / no). Authorization persists across resume via `authorized_main_branch: true` in `_progress.md`. Updated the "Never push to main/master" and "`--dry-run` means read-only" rules to match. Added two new Common Rationalization rows (--no-push on main is still a write; dry-run can't skip source verification).
- `README.md`: Replaced the stale "mines CHANGELOG.md for load-bearing decisions" description in scenario 19 with the current layered source model (implementation-history primary, lessons-learned secondary, code verification, CHANGELOG fallback).

## 0.65.1 — The Layered Spine (2026-04-15)

Rewrote `prompts/adr-log.md` to merge the better ideas from the user's own ADR-extraction prompt with the safety envelope from the first draft. Same prompt, much sharper — layered source mining with `implementation-history.md` as the primary spine, REQ/UR-keyed idempotency instead of fuzzy CHANGELOG-version matching, proper YAML `related: [{page, rel}]` relationships, per-cluster `topics/_index_*.md` wiki pages, and a completion report that forecasts remaining work sized S/M/L.

- `prompts/adr-log.md`: Replaced the mining spine (`CHANGELOG.md` → `implementation-history.md` primary, `lessons-learned/` secondary, current code for verification, `CHANGELOG.md` as portable fallback). Replaced the frontmatter schema (now `req:`, `ur:`, `sources:`, `related: [{page, rel}]`, `confidence`). Moved ADR files into `decisions/records/` and clusters into `decisions/topics/_index_<cluster>.md` as first-class wiki pages. Added explicit supersession workflow that flips the old ADR's `status` and adds the inverse `rel: superseded-by` to its `related` list in the same commit. Commit messages now follow `docs(adr): …` conventional shape. Added a completion-report section with a remaining-candidates forecast (sized S/M/L per UR). Kept the pre-flight safety checks, `--dry-run` / `--no-push` / `--batch-size` / `--from` flags, "infer alternatives if absent and mark `(inferred)`" guidance, and the Common Rationalizations / Red Flags / Verification Checklist guardrails.
- `prompts/README.md`: Updated the `adr-log` description to reflect the layered source model and REQ/UR-based idempotency.

## 0.65.0 — The Prompt Shelf (2026-04-15)

New `prompts` action — a dispatcher over a growing library of reusable, battle-tested prompts for recurring jobs the skill doesn't have a first-class action for. Seeded with `adr-log`, a create-or-update prompt that builds a project-wide Architecture Decision Record log at `decisions/` (BKB wiki pattern) by mining `CHANGELOG.md` for load-bearing decisions. Idempotent, resumable, supersession-aware.

- `actions/prompts.md`: New sub-command dispatcher (`list`, `show <name>`, `run <name>`, shorthand `<name>`) that resolves prompt names against `prompts/*.md` by exact match or unambiguous prefix. `show` is strictly read-only; `run` adopts the body below the `---` separator as operational instructions.
- `prompts/README.md`: Library index explaining the prompt file shape (title + blockquote + metadata + `---` + body) and how to add new entries.
- `prompts/adr-log.md`: First library entry. Detects create-vs-update mode via `decisions/_master_index.md`, resumes from `_progress.md` mid-run, allocates sequential `ADR-NNNN` numbers without reuse, handles supersession (sets `status: superseded` + `superseded_by` on the old ADR, never deletes), de-duplicates on re-run via a `source:` frontmatter field, and commits+pushes in batches (scaffolding → mining → ADRs in groups of 3 → final reconciliation).
- `SKILL.md`: New priority-19 routing row for `prompts` / `prompt`, new Verb Reference entry, new Action Dispatch entry, new "Prompt library:" block in the bare-invocation help menu, and `prompts` added to the foreground-dispatch list.
- `next-steps.md`: Three new post-action sections (`prompts list`, `prompts show`, `prompts run`).
- `README.md`: New numbered scenario "19. Run a saved prompt"; renumbered later scenarios 19→20, 20→21, 21→22.
- `CLAUDE.md`: Registered `actions/prompts.md` and the `prompts/` directory in the Project Structure tree.

