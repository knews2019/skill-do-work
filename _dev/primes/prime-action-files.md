# Prime: Action Files

> Read this before adding or modifying any action file under `skills/do-work*/actions/`.
> Prime files are low noise, high value: they point at the code that is the source of
> truth, carry no volatile metrics, and every new lesson in this domain gets linked from
> here so future sessions find it. This file is maintainer-side (`_dev/` is
> export-ignored) — nothing shipped may cite it.

## Every New Action Must Justify Its Package

Before an action file is added, state — in its description blockquote or an accompanying ADR — why it belongs in core, board, knowledge, or toolbox: what package machinery it needs or which existing action it completes. Reviewers reject additions without this justification. Every new action also updates the owning `skills/do-work*/SKILL.md`, help, and any next-step surface; router budgets are enforced by `_dev/tests/contract-regressions.sh`.

## Template

See [the action template and earned-section rules](lessons-action-files.md#template).

## Accepted Variants

- **Sub-command dispatchers** (`prime.md`, `bkb.md`) — Use a Sub-Commands table instead of flat steps. Each sub-command has its own workflow section.
- **Multi-mode actions** (`present-work.md`, `review-work.md`, `tutorial.md`) — Use a Modes table, then separate workflow sections per mode. A single `Step 1: Mode Selection` dispatcher at the top is acceptable.
- **State-based actions** (`version.md`, `pipeline.md`) — Response sections keyed by input type instead of sequential steps.
- **Checklist-based diagnostics** (`forensics.md`) — Use a `## Checks` section with independently-runnable items instead of ordered `## Steps`. Each check is a diagnostic probe, not a sequential step.

## Cross-Referencing

Cross-reference same-package actions by their local path (for example `actions/work.md`). Cross a package boundary with the **literal relative path from the citing file's own directory** — the spelling a reader could paste into a terminal there and have resolve, in both the source tree and an install. The depth is per-file, not a fixed prefix: from a package-root `SKILL.md` the sibling is `../do-work-board/...`, from `actions/` or `crew-members/` it is `../../do-work-board/...`, and from `tools/queue-kanban/` it is `../../../do-work/...`. Never write `../<package>/...` as shorthand for "up to the skills folder" — that skill-root-relative reading was retired by REQ-249, and from anywhere below a package root it points at nothing. **A fenced block's payload is exempt; its annotations are not.** The exemption's reason is that the payload's text lands in some other file — a consumer's REQ, a report — so it is content rather than a citation from here. A `#` or `<!--` annotation beside that payload lands nowhere: it is documentation addressed to the reader of *this* file, so its paths are citations and are checked. Keying the exemption to the fence character instead of to that reason is what shielded eight wrong-depth citations inside one schema block.

**The citation class is drawn by what a citation is, not by the punctuation around it** (REQ-269). A token is a cross-package citation when its first path segment — after any `../` lead — names a sibling package directory, whether or not it is backticked, and whether or not it leads with `../`. Two conditions take a token back out of the class, and both are about meaning rather than spelling: a path rooted somewhere the reader is told to resolve from instead (`<skill-root>/…`, `<suite-root>/…`, `.claude/skills/…`) is not a path from here at all, and a **bare** `do-work/…` token is the consuming project's own queue state, because the core package's directory name is also that root. Only `do-work` carries that collision; a bare `do-work-board/…`, `do-work-knowledge/…` or `do-work-toolbox/…` token has no consumer-state reading and is always a citation. `_dev/tests/shipped-package-reference-contract.sh` enforces exactly this condition in both topologies (alongside its Markdown-link checks), and its `run_citation_surface_fixtures` pins which tokens the check ever sees.

Shipped files must never cite this repo's own `CLAUDE.md` or `AGENTS.md` — both are export-ignored maintainer instructions, so the citation dangles in every consumer install. `_dev/tests/contract-regressions.sh` enforces this across the shipped `skills/` tree.

## Descriptions Are Triggers, Not Summaries

A skill or action description is loaded whether or not the thing gets used, so it should carry the words that tell the router when to pull it in — not an explanation so complete the file never gets read.

**Bad:** "Monitors a pull request through review and CI, rebasing as needed and addressing reviewer comments until checks pass."

**Good:** "Use when the user asks to monitor, watch, or babysit a PR."

## Traps

- [family: alternate-writer-contract-drift] Changing an emitted artifact or routing rule only at its primary writer leaves alternate modes and action-bearing readers silently following the old contract; grep every writer and downstream reader against the recorded scope before declaring it shipped — including sibling-package actions, board labels, and prime sentences that restate the rule (REQ-566 found six such restatements outside the three files that defined the contract).
- [family: budgeted-context-routing] Routing lessons only when a REQ is captured → old REQs and later serial siblings miss lessons written before their claim; re-run the same bounded projection at claim for every route, counting captured entries first.

## Agent Compatibility

See [`lessons-action-files.md#agent-compatibility`](lessons-action-files.md#agent-compatibility).

## Stakes

- Action routing and authored artifacts
  Req: keep owning skill routers, emitted fields and downstream readers aligned.
  Value: the requested action is discoverable and its output remains usable across agents.
  Risk: a local prose change can silently bypass a workflow or strand records; follow the Cross-Referencing contract before shipping.

## Lessons

See [`lessons-action-files.md`](lessons-action-files.md) — read it before changing what **Read first** or **Traps** name above.
