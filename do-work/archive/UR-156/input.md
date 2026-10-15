---
id: UR-156
title: 'ai-report second-opinion: another model rebuilds an existing report bundle'
created_at: 2026-10-15T11:02:38Z
requests: [REQ-698]
word_count: 1701
---
## Summary

Migrated from the consumer repository iwr2-imprint-wright (its UR-061 / REQ-287, cancelled there on 2026-10-15 in favour of this upstream capture). The user wants a second opinion on an AI-written report: `ai-report second-opinion` hands an existing `ai-reports/<bundle>/index.html` to another model, which rebuilds it as a clarity-edited sibling page `index-<cli>-<model-slug>.html` in the same bundle. Three modes: session (the calling AI coder's own CLI and model, the default), targeted (`--backend <cli> [--model <slug>]`), batch (`--all`, every installed backend except the caller). The original page is never touched by the model; after the checks pass, the launcher inserts one marker-wrapped banner per second opinion linking to it. The consumer will receive the feature through `just run-do-work-update`.

## Extracted Requests

| REQ | Title |
|-----|-------|
| REQ-698 | ai-report second-opinion: a second model rebuilds a report bundle as index-<cli>-<model-slug>.html (session, targeted and batch modes) |

## Batch Constraints

- The model that rebuilds the page never modifies `index.html`, existing assets, or anything outside the bundle folder; codex runs under its workspace-write sandbox with the bundle as working directory, every backend is followed by a `git status --porcelain` check that fails the run on any path outside the bundle.
- The output file name is derived by the launcher from the backend and model that actually ran, never chosen by the model or hard-coded.
- `index.html` changes only by the launcher's marker-wrapped banner blocks, inserted idempotently after the git status check and a judge `pass`.
- In session and targeted modes the launcher never substitutes another backend or model.

## Full Verbatim Input
> ```
> do-work validate-feedback:  Can we add the following skill so we extend the ai-report and can have a second opinion from another model that rebuilds a new html, like index-[model-name-version].html i.e. index-argon.html
> 
> # See bellow an attempt to build it
> 
> ROLE. You are a clarity editor for AI-generated technical reports. Your job is to make an existing report easier to understand and more accurate, not denser.
> 
> INPUT. A folder containing index.html and its assets (images, CSS, JS, data files), plus the code repository the report describes. The report folder lives inside it or is linked to it. Read index.html and every asset it references before you start.
> 
> USING THE REPOSITORY. You may explore the repository freely, including source code, READMEs, docs, configs, tests, and commit history, to fill gaps the report leaves unexplained, to confirm that the report's claims match the code, and to find concrete examples (file names, function names, real data flows) that make ideas easier to grasp. Look things up to clarify, not to pile on. Include a repository detail only if it helps the reader understand what is being built. Reference the files you used with their relative paths.
> 
> OUTPUT. Create one new file, index-argon.html, in the same folder. Never modify or delete index.html, any existing asset, or anything else in the repository. You may reuse existing assets through relative paths. Any new assets you need go in the same folder with an "argon-" prefix.
> 
> GOAL. After reading index-argon.html, a reader should clearly understand what is being built: what it is, why it exists, how its parts fit together, and what state it's in.
> 
> RULES.
> Rule 1, ground everything in sources. Every fact must come from the original report or the repository. Don't fill gaps from your own background knowledge. If something important can't be verified in either, mark it as an "Open question" instead of guessing.
> Rule 2, review the original report. Check its claims against the repository. When something is wrong, outdated, or misleading, use the correct version in the new report and record it in a "Corrections" section, showing what the original said, what's actually true, and the file that proves it. Also list anything you couldn't verify.
> Rule 3, don't make it dense. Use short paragraphs (no more than 3 to 4 sentences), one idea per section, and plain words. Explain every acronym or term the first time it appears. Put deep detail behind expandable sections.
> Rule 4, be visual. Wherever a picture explains better than prose, use one: architecture or flow diagrams (inline SVG or Mermaid) based on the real code structure, timelines, comparison tables, status cards, and charts built from real data.
> Rule 5, be interactive. Add things that help understanding, such as collapsible sections, tabs for alternative views, hover tooltips for terms, clickable diagram parts that reveal details and file paths, and a sticky table of contents. Keep it a single self-contained HTML file that runs offline, using vanilla JS or a small library loaded from a CDN, and degrade gracefully without JS.
> Rule 6, structure. TL;DR: what is being built, in 3 sentences. The big picture: one main diagram of the whole system. Components: one card per part, covering what it does, its inputs and outputs, where it lives in the repo, and its status. How it works: a step-by-step flow. Decisions and trade-offs. Status, risks, and open questions. Corrections to the original report. Glossary. Source map: links back to index.html sections and to repository files. Change this order only if the content clearly calls for it.
> Rule 7, keep it faithful. Reorganizing, simplifying, and visualizing are welcome. Change a claim only when the repository proves it's wrong, and say so when you do.
> Rule 8, accessibility and polish: responsive layout, readable font size, good contrast, alt text on images, and labels on diagrams.
> 
> PROCESS. Step a, read the full report and its assets. Step b, explore the repository enough to understand its structure and to check the report's key claims. Step c, write a short private outline: the core idea, the main parts, where readers are likely to get confused, and any errors you found. Step d, build index-argon.html. Step e, check it: every claim traces back to the report or the repo, all corrections are listed, nothing important is missing, all links resolve, and the page opens without errors.
> 
> FINAL REPLY. Confirm the file path, list in 3 to 5 bullets what you changed to improve clarity, then list the corrections you made and any open questions.
> 
> ---- Decisions taken in the validate-feedback session (2026-10-15), chosen by the user via the question tool:
> - Placement: project-owned skill folder .claude/skills/ai-report-second-opinion/ (survives do-work suite updates).
> - File name: index-<cli>-<model-slug>.html, written by the launcher from the resolved model (not a hard-coded nickname such as argon).
> - Backend order: codex with the workspace-write sandbox first, gemini as fallback.
> - Next step: capture one REQ.
> 
> ---- Capture request as sent to do-work:
> capture-request Add a project-owned skill `.claude/skills/ai-report-second-opinion/` that asks a second model to rebuild an existing ai-report bundle as a clarity-edited sibling page. Input: a bundle directory under `ai-reports/` holding `index.html`. Output: `index-<cli>-<model-slug>.html` in the same folder (e.g. `index-codex-gpt-5.html`), new assets prefixed `<cli>-<model-slug>-`; never modifies `index.html`, existing assets, or anything else in the repo. Pieces: (1) `prompts/clarity-editor.md` holding the user's clarity-editor prompt (role, rules 1-8, process a-e, final reply) with placeholders for absolute bundle path, absolute repo root, output file name and asset prefix; fix its CDN/offline contradiction to the suite design rule (CDN only Tailwind and Mermaid, rest inline). (2) A launcher script `scripts/run_second_opinion.sh` invoked via `bash`, not chmod +x: resolves the bundle, runs `codex exec -s workspace-write` with cwd set to the bundle dir so writes are fenced to that folder, reads the resolved model from the codex output to build the file name, falls back to `gemini -p --approval-mode auto_edit` when codex is absent, then runs `git status --porcelain` and fails if any path outside the bundle changed. (3) Render check: copy the new page into a scratch bundle as `index.html`, run `do-work-cli.sh ai-report-judge` on it, report the verdict, delete the scratch copy. (4) `SKILL.md` with argument `<bundle-dir|latest>`, the backend order, and a final reply that lists the file path, clarity changes, corrections and open questions. (5) Link the skill from `.claude/skills/make-ai-report-with-screenshot/SKILL.md` and from the do-work-toolbox ai-report guide only if that file is project-owned; otherwise note the upstream candidate in `docs/imprintwright-v2/PLAN.md`. Test on `ai-reports/2026-10-11_0935_imprintwright-onboarding/`.
> 
> ---- Update 1 (2026-10-15, sent in the consumer repo iwr2-imprint-wright after do-work verify-request):
> Then amend REQ-287 with a mutual-linking rule for the ai-report-second-opinion skill: the new index-<cli>-<model-slug>.html links back to index.html at the top, and index.html gets a small note or banner saying a second-opinion version is available, with a link to it. This conflicts with the rule that index.html stays byte-identical, so the model must never touch index.html. Only after the post-run git status check and the judge pass, the launcher inserts the banner (idempotently, skipping it if it's already there), and the constraint and GREEN criteria change to "byte-identical except the launcher's banner block". If that's cleaner, put the link on the ai-report catalog page (ai-reports/index.html) instead and leave index.html fully untouched. Pick one, say which, and update the REQ.
> 
> ---- Update 2 (2026-10-15, final wording; it superseded two earlier versions of the same message):
> Update REQ-287: replace the codex-then-gemini fallback chain with explicit modes for the ai-report-second-opinion skill, and set its official trigger phrase. Keep everything else in REQ-287 unchanged.
> 
> 0. Trigger phrase. The official trigger is "ai-report second-opinion", which parallels "do-work ai-report": ai-report creates a report, and ai-report second-opinion gives a second opinion on an existing one. Put it at the start of the SKILL.md frontmatter description so any AI coder (Claude Code, Codex, Antigravity/agy, Gemini) picks it up. Forms (bundle defaults to latest):
>    ai-report second-opinion [<bundle-dir>|latest]                                   (session mode)
>    ai-report second-opinion --backend <cli> [--model <slug>] [<bundle-dir>|latest]  (targeted mode)
>    ai-report second-opinion --all [<bundle-dir>|latest]                             (batch mode)
>    Natural phrasings map to the same forms: "generate a second opinion" is session mode, "Codex, generate a second opinion" is targeted mode, and "generate second opinions from all models" is batch mode.
> 
> 1. Session mode (default, no backend named). Inside an active AI coder session, the skill uses the model currently active in that session. Typed in an Antigravity session, it runs Antigravity's current model; in Codex, Codex's. The calling agent passes its own CLI and model slug (--backend <own-cli> --model <own-slug>); SKILL.md tells it to. Don't fall back to another backend here. If the session's model can't be determined, stop and say so. Only when invoked outside any AI coder session, such as from a plain terminal with no active model, does the skill use the configured backend order (codex, gemini, agy) and run the first one installed.
> 
> 2. Targeted mode. Naming one backend or model runs only that one. If it isn't installed or the model can't be resolved, stop and say so. Never substitute another backend.
> 
> 3. Batch mode (--all). Runs the clarity editor once on every installed backend (codex, gemini, agy) except the one driving the session (--exclude <own-cli>), producing one index-<cli>-<model-slug>.html per backend. Runs go one at a time, each with its own post-run git status check, judge check, banner insert and <cli>-<model-slug>-reply.md. A failed or missing backend is reported and skipped; the batch continues and ends with a per-backend summary (ran, file name, judge verdict, or the reason it was skipped).
> 
> In every mode, the launcher resolves the model slug from the backend that actually ran (per prime-agy-usage.md for agy and gemini) and writes the file name itself. Codex keeps the workspace-write sandbox. gemini and agy rely on the post-run git status check. Banner markers are per file name, so index.html can carry one banner per second opinion, each inserted idempotently. Update the launcher usage and the Red-Green Proof: "ai-report second-opinion" from a Claude Code session uses Claude Code's model, "--backend codex" proves one codex file, and "--all" from Claude Code proves one file per available non-calling backend, with index.html changed only by the marker-wrapped banners.
> 
> ---- Migration instruction (2026-10-15):
> ok, I need to you migrate UR/REQ to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
> because that is the upstream do-work repository, cancel it here, we will get it via just run-do-work-update  when it is done upstream
> ```
