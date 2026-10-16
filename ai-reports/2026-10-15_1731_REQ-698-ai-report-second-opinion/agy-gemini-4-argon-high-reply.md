Filled the assigned staging file: [index-agy-gemini-4-argon-high.html](index-agy-gemini-4-argon-high.html) (second-opinion rebuild of [index.html](index.html); `index.html` and all repository files remain untouched).

### What Changed to Improve Clarity (5 Bullets)

- **Added a 3-sentence TL;DR, sticky navigation, and at-a-glance stat tiles:** Opened with a plain-language 3-sentence summary explaining what `ai-report second-opinion` does, why its staging/guard pipeline exists, and its current test/live status, paired with a sticky table of contents and 4 summary metric tiles.
- **Made the 6-stage pipeline interactive and separated the two failure gates:** Added a clickable SVG pipeline diagram (`report_second_opinion.go`) that reveals each stage's inputs, outputs, and exact line ranges on click/keyboard selection, and separated **Stages 1–4a (abort before publish)** from **Stage 4b → 5 (sibling page & reply published, but Stage 6 banner withheld if `judgeRender` does not pass)**.
- **Organized the codebase into 6 component cards with role filtering and a 4-mode tabbed walkthrough:** Replaced the flat file list with 6 structured component cards (covering purpose, inputs, outputs, repository path/lines, and status) plus interactive tabs for **Session**, **Targeted**, **Batch (`--all`)**, and **Plain-Terminal Fallback** modes.
- **Moved dense protocol and process-group mechanics behind native `<details>` disclosures:** Explained every acronym (`CLI`, `REQ`, `UR`, `AGY`, `PGID`, `ESRCH`, `inode`) on first use with hover tooltips and a full **Glossary**, while placing the 4-backend CLI flag/envelope table, the `2abd7939` POSIX supervisor race fix, and the 22-file accounting table inside expandable sections that work without JavaScript.
- **Added a visual verification chart and complete 6-row live proof matrix:** Built an inline SVG bar chart from recorded heavy-lane and review-score metrics (`80.63%` pre-remediation to `86.44%` post-remediation) next to the complete 6-row live verification table and a **Source Map** disclosing why repository paths outside the bundle are shown as plain text rather than unpushed remote URLs.

---

### Corrections to the Original Report ([index.html](index.html))

1. **C1 — Plumbing file count mismatch (`14 paths / groups` vs. 22+ listed files):**
   - **Original report said:** Section 2 headed its second column *"Plumbing That Came Along (14 paths / groups)"* while listing 10 bullet groups covering 22+ files.
   - **What is actually true:** The builder's `Implementation Summary` touched **20 files total** (6 core feature files + **14 non-feature plumbing, live-proof, and CLI prime/lesson files**). The original report kept the `14` count (`20 − 6`), then appended `do-work/lessons-index.md` and the release/archive files from commit `3645f3a1` without updating the heading.
   - **Proof:** `do-work/archive/UR-156/REQ-698-ai-report-second-opinion.md:255-275` and release commit `3645f3a1`.
2. **C2 — Incomplete batch table and conflation of the two AGY runs:**
   - **Original report said:** Section 4's live table omitted the batch run's `codex` (`excluded`) and `agy` (`failed`) rows, and attributed AGY's failure solely to placing a `<nav>` sidebar before the opening backlink.
   - **What is actually true:** AGY failed in two separate runs for different reasons: (a) during `--all --exclude codex`, AGY (`gemini-3.8-flash-low`) failed under the **old overly strict backlink parser** before commit `2f5cd87a` (Decision `D-09`) fixed false rejections on attribute order and inline tags, and its deleted stage was not retained (`REQ-698` explicitly notes: *"Do not call the original batch refusal proven model noncompliance"*); (b) in the later authorized single retry at `69a2a998`, the retained candidate placed a `Report Navigation` container/sidebar and full table of contents before `<main class="main-content">`.
   - **Proof:** `do-work/archive/UR-156/REQ-698-ai-report-second-opinion.md:146, 292, 305-316, 330-347` and commit `2f5cd87a`.
3. **C3 — Omitted post-generation verifier edit on the Codex live proof artifacts:**
   - **Original report said:** Sections 2 and 4 presented `index-codex-gpt-6-astra.html` and `codex-gpt-6-astra-reply.md` as direct model outputs from the targeted Codex run.
   - **What is actually true:** Codex's initial output included 22 source-map hyperlinks pointing to unpushed local commit `406ff3c4` and a worktree staging path in its reply. An authorized verifier edit (commit `9fa76500`, Decision `D-08`) stripped the 22 unpublished-commit URLs to plain text, normalized the reply link, appended a verifier disclosure note to `codex-gpt-6-astra-reply.md`, and re-ran `ai-report-judge`.
   - **Proof:** `do-work/archive/UR-156/REQ-698-ai-report-second-opinion.md:145, 277, 321-323` and `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-gpt-6-astra-reply.md:19`.
4. **C4 — Omitted 4th invocation mode (`terminal`) and `claude` exclusion from `--all` batch:**
   - **Original report said:** Described only 3 invocation modes (`session`, `targeted`, `batch`) and did not note that `claude` is excluded from `--all`.
   - **What is actually true:** `report_second_opinion.go` implements a 4th plain-terminal mode (`InvocationMode: "terminal"`) when neither `--backend` nor `--all` is passed, selecting the first installed binary in `codex → gemini → agy` order once without fallback on failure. Additionally, `claude` is valid only with explicit `--backend claude` and is not part of the `--all` batch slice (`[]string{"codex", "gemini", "agy"}`).
   - **Proof:** `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:52-69, 137-139` and `skills/do-work-toolbox/actions/ai-report.md:69`.
5. **C5 — AGY two-step catalog + zero-turn `/model` probe and 16m/16MiB generation bounds:**
   - **Original report said:** Stage 2 and Section 3 Row 4 described a uniform "90s read-only probe" without distinguishing AGY's zero-turn probe or the generation runner's limits.
   - **What is actually true:** When `--model` is omitted for `agy`, the probe first runs `agy models` and selects the newest `gemini-<ver>-flash-low|medium` (`secondOpinionFlashModel`), then runs the zero-turn slash command `agy --model <m> --print-timeout 30s --output-format json -p /model` instead of the 1-turn `"Reply OK only..."` prompt used by `codex`, `gemini`, and `claude`. Main generation runs under a 16-minute timeout (`--print-timeout 15m` for `agy`) and a 16 MiB per-stream (`stdout`/`stderr`) capture cap.
   - **Proof:** `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:329-354, 527-530, 555-582, 823-863`.
6. **C6 — Stage 4 publication precedes Stage 5 `judgeRender`:**
   - **Original report said:** The diagram's fail-closed note implied stages 2–4 abort before publication and stage 5 skips stage 6, without clarifying what is already written to disk before stage 5 runs.
   - **What is actually true:** `publish(pageName, page)` and `publish(replyName, []byte(reply))` execute at lines 446–457 **before** `secondOpinionJudge` runs at line 458 (and any model-created prefixed assets stay visible on refusal per lines 393–394). When `secondOpinionJudge` returns `fail`, `skipped`, or `error`, the sibling HTML and `-reply.md` remain published in the bundle (`RunState: "ran"`), while Stage 6 (banner insertion in `index.html`) is skipped.
   - **Proof:** `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:393-394, 446-465`.

**Unverified external references (outside repository boundary):**
- Temporary evidence logs and candidate captures under `/private/tmp/req698-evidence/*` and `/private/tmp/req698-review/*` cited in `REQ-698-ai-report-second-opinion.md`, as well as the external consumer repository `iwr2-imprint-wright`, lie outside `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and were not accessed. Test execution metrics are reported from committed verification records rather than re-running workflows.

---

### Open Questions

1. **Local Gemini CLI authentication (`exit 55 UNSUPPORTED_CLIENT`):** What local credential or client update is required so `gemini -p ... --output-format json --approval-mode auto_edit` succeeds during live `--all` batch runs?
2. **Live session-mode invocation from Claude Code:** Because `REQ-698` was built inside a Codex session where active model introspection was unavailable, Claude was verified via pinned `--backend claude --model claude-opus-5-5`. Has a live natural-language session-mode call (`"generate a second opinion"`) been exercised directly from an interactive Claude Code session?
3. **Consumer repository rollout (`iwr2-imprint-wright`):** Once the 15 unpushed local commits on `main` are pushed to `origin/main` and `just run-do-work-update` is run in `iwr2-imprint-wright`, does the original consumer target bundle (`ai-reports/2026-10-11_0935_imprintwright-onboarding/`) pass end-to-end?
