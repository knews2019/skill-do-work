# Install Action

> **Part of the do-work-toolbox skill.** Installs four optional companion capabilities into the current project: `ui-design`, `bowser`, `last30days`, and `ideation-adhd`. Board recipes belong to `do-work-board`, memory setup belongs to `do-work-knowledge`, and suite self-update belongs to core `do-work`.

Each target is idempotent and follows detect → install → verify → report. The action dispatches on the first argument and never chains targets.

## When to Use

**Use when:**
- The project is about to start UI work or browser automation and the matching companion skill/tooling is missing.
- A `ui-review` pass flagged "frontend-design skill not installed" or "visual verification skipped" — install the missing piece.
- A `domain: ui-design` REQ is about to be built and the builder would benefit from skill-level design knowledge (`install ui-design`).
- The user asked for headed-browser workflows, screenshots, or visual verification (`install bowser`).
- The user asked for social research, trend scanning, or "what's the discourse on X" capabilities (`install last30days`).
- The user wants structured divergent ideation — deliberately unconventional candidate ideas for a named problem, beyond what `scan-ideas` surfaces from the repo (`install ideation-adhd`).

**Do NOT use when:**
- The project explicitly uses a different design system or browser-automation tool and adding the do-work default would conflict.
- The environment can't install global npm packages (for `bowser`) and the user hasn't consented to a local-only install.
- The user just wants to view the board once — that's `do-work-board board` (`../../do-work-board/actions/board.md`), no install needed.

## Input

`$ARGUMENTS` selects the install target:

- `ui-design` — Install Anthropic's `frontend-design` skill for production-grade UI design capabilities.
- `bowser` — Install Playwright CLI (global) plus the Bowser skill (project-scoped) for browser automation, screenshots, and visual UI verification.
- `last30days` — Vendor the engagement-ranked social-research engine (project-scoped, git-ignored, keyless).
- `ideation-adhd` — Install the upstream `adhd` skill (project-scoped) for parallel divergent ideation: isolated branches under distinct cognitive frames, scored, clustered, and deepened. (`adhd` is an accepted alias for this target.)

If `$ARGUMENTS` is empty or doesn't match a known target, print the help block (target list + one-line blurb each) and stop.

## Install Manifest

Every target follows the same four-step shape (detect → install → verify → report). The per-target commands and blurbs live here:

| target | detect_cmd | install_cmd | verify_cmd | blurb |
|--------|------------|-------------|------------|-------|
| `ui-design` | `test -s "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md"` | `mkdir -p "$PROJECT_ROOT/.claude/skills/frontend-design" && <skill-root>/../do-work/scripts/atomic-download.sh https://raw.githubusercontent.com/anthropics/skills/main/skills/frontend-design/SKILL.md "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md"` | `test -s "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" && echo "Installed successfully" || echo "Installation failed"` | Anthropic's `frontend-design` Claude skill — production-grade UI design capabilities (typography, color, spacing, layout, component design, responsive/mobile-first, accessibility). |
| `bowser` | `playwright-cli --help >/dev/null 2>&1 && test -s "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md"` | (multi-step — see `bowser` workflow below) | (multi-step — see `bowser` workflow below) | Playwright CLI + Bowser skill — headed/headless browser sessions with Chromium, screenshots at any viewport, DOM snapshots, parallel named sessions, persistent profiles. |
| `last30days` | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | (multi-step — see `last30days` workflow below) | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | Engagement-ranked social-research engine — Reddit/HN/Polymarket/GitHub/YouTube keyless out of the box; X/TikTok/Instagram unlock only via user-global API keys. |
| `ideation-adhd` | `test -s "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md"` | `mkdir -p "$PROJECT_ROOT/.claude/skills/adhd" && <skill-root>/../do-work/scripts/atomic-download.sh https://raw.githubusercontent.com/UditAkhourii/adhd/main/skills/adhd/SKILL.md "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md"` | `test -s "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md" && echo "Installed successfully" || echo "Installation failed"` | The `adhd` skill (MIT) — parallel divergent ideation: spawns isolated branches under distinct cognitive frames (regulator, biology, speedrunner, 10-year-old, $0 budget, …), scores on novelty/viability/fit, clusters, prunes traps, deepens the top survivors. Explicitly invoked (`/adhd`), never fires on its own. |

In every command above, resolve `PROJECT_ROOT` first:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
```

Every skill-file download follows the canonical [Atomic download publication](../../do-work/docs/prescribed-shell-primitives.md#atomic-download-publication) contract: land in a `SKILL.md.download` sibling, rename into place only on success, and remove the temporary file while preserving a failure status. The table keeps each executable command inline because every target must remain independently runnable.

## Steps

### Step 1: Dispatch on `$ARGUMENTS`

- If `$ARGUMENTS` is empty, prints "help" / "?", or doesn't match a known target → print the help block (target list + blurb) and stop.
- If `$ARGUMENTS` matches a known target → proceed to the target-specific workflow below.

### Step 2: Run the target's workflow

The `ui-design` and `ideation-adhd` workflows use the manifest commands directly. `bowser` and `last30days` use the multi-phase workflows below.

---

## Workflow: `ui-design`

#### Phase 1: Check if already installed

Resolve `PROJECT_ROOT`, then run the manifest's `detect_cmd`. If the file exists and is non-empty, report "already installed" and stop. (`test -s`, not `ls` — a zero-byte file from an interrupted download must read as absent so a re-run can repair it.)

#### Phase 2: Install the skill

Run the manifest's `install_cmd`. If `curl` is unavailable or the download fails, check the environment's plugin/skill registry (e.g., `/plugin install frontend-design`) and install from there.

#### Phase 3: Verify

Run the manifest's `verify_cmd`.

#### Phase 4: Report back

```
Installed: frontend-design skill

Gives Claude production-grade UI design capabilities:
- Professional visual aesthetics (typography, color, spacing, layout)
- Component design with proper states and variants
- Responsive, mobile-first implementations
- Accessibility-compliant interfaces

Works alongside do-work's `domain: ui-design` rules — the skill provides
implementation expertise while the domain rules provide workflow structure.
Requests tagged `domain: ui-design` benefit from both automatically.
```

---

## Workflow: `ideation-adhd`

A single self-contained `SKILL.md` — same shape as `ui-design`. The target is named `ideation-adhd` (the name says what it does), but the install **folder** stays `adhd` to match the upstream skill's own frontmatter `name:` field, so it auto-discovers as the `/adhd` slash command; do not rename the folder to match the target.

#### Phase 1: Check if already installed

Resolve `PROJECT_ROOT`, then run the manifest's `detect_cmd`. If the file exists and is non-empty, report "already installed" and stop. (`test -s`, not `ls` — a zero-byte file from an interrupted download must read as absent so a re-run can repair it.)

#### Phase 2: Install the skill

Run the manifest's `install_cmd`. If `curl` is unavailable or the download fails (offline, upstream restructured), report the error and direct the user to https://github.com/UditAkhourii/adhd for manual install — the file lives at `skills/adhd/SKILL.md` in that repo. Upstream also documents a global `npm install -g adhd-agent` CLI — do not run it; this skill's norms reject global installs, and the project-scoped `SKILL.md` copy is the only supported path.

#### Phase 3: Verify

Run the manifest's `verify_cmd`.

#### Phase 4: Report back

```
Installed: ideation-adhd — the adhd skill (parallel divergent ideation)

Gives the agent a structured divergence engine:
- Spawns isolated branches under distinct cognitive frames
  (regulator, biology, speedrunner, 10-year-old, $0 budget, ...)
- Scores on novelty/viability/fit, clusters by angle, prunes traps,
  and deepens the top 3 survivors
- Explicitly invoked via /adhd or "ADHD mode" — it never fires on its own

Complements do-work's ideation actions rather than replacing them:
`scan-ideas` finds grounded opportunities in THIS repo; adhd generates
deliberately unconventional candidates for a problem you name. Feed the
winners to `do-work capture-request:` to queue them.

Note: the skill leans on parallel subagents for its divergence phase — on
an agent without them, the frames run sequentially instead.
```

---

## Workflow: `bowser`

The `bowser` target installs two components: the global `playwright-cli` (plus a Chromium browser), and the project-scoped `playwright-bowser` skill.

#### Phase 1: Check if already installed

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
playwright-cli --help >/dev/null 2>&1 && echo "playwright-cli: installed" || echo "playwright-cli: not found"
test -s "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" && echo "bowser skill: installed" || echo "bowser skill: not found"
```

If both are present, report installed and stop.

#### Phase 2: Install Playwright CLI

```bash
npm install -g @anthropic-ai/playwright-cli@latest
```

If `npm` isn't available:

```bash
yarn global add @anthropic-ai/playwright-cli@latest
```

If neither package manager works, report the error and the install command so the user can run it manually.

#### Phase 3: Install Playwright browsers

```bash
playwright-cli install --with-deps chromium
```

Only Chromium is installed — sufficient for UI review. For Firefox/WebKit, the user can add `playwright-cli install firefox` or `webkit` later.

If `--with-deps` fails due to permissions:

```bash
npx playwright install chromium
```

#### Phase 4: Install the Bowser skill

Create `$PROJECT_ROOT/.claude/skills/playwright-bowser`, then invoke:

```bash
<skill-root>/../do-work/scripts/atomic-download.sh https://raw.githubusercontent.com/disler/bowser/main/.claude/skills/playwright-bowser/SKILL.md "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md"
```

The helper preserves a failed download status and never publishes a partial final file. The `ui-design` and `ideation-adhd` manifest rows follow the same atomic-download contract.

If the URL 404s (the repo may have restructured), report the error and direct the user to https://github.com/disler/bowser for manual install — the file lives somewhere under `.claude/skills/` in that repo.

#### Phase 5: Verify

```bash
playwright-cli --help >/dev/null 2>&1 && echo "playwright-cli: OK" || echo "playwright-cli: FAILED"
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
test -s "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" && echo "bowser skill: OK" || echo "bowser skill: FAILED"
```

#### Phase 6: Report back

```
Installed: Playwright CLI + Bowser skill

Gives agents browser automation capabilities:
- Headed/headless browser sessions with Chromium
- Screenshots at any viewport (mobile, tablet, desktop)
- DOM snapshots for accessibility and element inspection
- Parallel named sessions for independent browser tasks
- Persistent profiles (cookies, localStorage preserved)

Works alongside `do-work-toolbox ui-review` — when Playwright CLI is
detected, ui-review automatically runs visual verification. The `ai-report`
action also consumes Playwright when available for live screenshots;
without it, ai-report falls back to SVG + Mermaid diagrams.

To use directly: playwright-cli -s=my-session open https://example.com --persistent
```

---

## Workflow: `last30days`

The `last30days` target vendors the engagement-ranked social-research engine (https://github.com/mvanhorn/last30days-skill) into the project as a git-ignored, keyless drop. Reddit, Hacker News, Polymarket, GitHub, and YouTube work with no API keys; X/TikTok/Instagram unlock only via user-global keys — never via project files.

#### Phase 1: Check if already installed

Run the shipped full-guarantee check (complete runtime payload, ignore rule, Python 3.12+). The runtime payload requires, at minimum, non-empty `SKILL.md` and `<project-root>/.claude/skills/last30days/scripts/last30days.py`; the installer and checker share that predicate so a sentinel-only or half-copied tree cannot masquerade as installed.

```bash
<skill-root>/../do-work/tools/do-work-cli.sh --repo-root <project-root> --format json install-last30days check <project-root>
```

- **All checks pass** (the ignore rule counts as passing when the project isn't a git repo) → report "already installed" and stop.
- **Runtime payload complete but the ignore rule failed** → proceed to Phase 2 in *repair mode*: skip publication and run only the guarded ignore step. (A missing Python 3.12+ interpreter isn't repairable by this action — report it per Phase 3.)
- **Runtime payload incomplete** → run Phase 2 in full, even when `SKILL.md` exists.

#### Phase 2: Vendor the skill

The upstream repo keeps the actual skill at `skills/last30days/` (self-contained — `SKILL.md`, a `scripts` directory, and supporting directories). The installer shallow-clones to temporary storage, rejects an incomplete source, copies the full subtree into a private staging directory adjacent to the destination, and validates the staged payload before a same-filesystem rename publishes it. When replacing an incomplete destination, it holds the prior tree in a private adjacent backup until the validated replacement is live; a publication failure restores the prior tree byte-for-byte. A complete destination keeps the existing no-op/ignore-repair behavior.

```bash
<skill-root>/../do-work/tools/do-work-cli.sh --repo-root <project-root> --format json install-last30days install <project-root>
```

If the block prints FAILED (offline, incomplete upstream payload, copy failure, or publication failure), **stop here** and report the error. The final path remains either the prior incomplete tree or a fully validated replacement; the installer never merges files from different upstream versions and cleans private staging/backup paths on exit.

The installer then makes the ignore claim true with the canonical [Local Git ignore](../../do-work/docs/prescribed-shell-primitives.md#local-git-ignore) helper and rejects an already-tracked vendored copy. It never touches the project's committable `.gitignore`.

Two hard constraints on this phase:

- **Write no config file — anywhere.** The engine reads API keys and settings from the user-global `~/.config/last30days/.env`, which the user manages themselves. Upstream also supports a project-local `.claude/last30days.env`, but it is trust-gated: the engine only reads it when `LAST30DAYS_TRUST_PROJECT_CONFIG` is already set in the environment or the user-global config — writing the trust flag *inside* the project file it gates is circular and does nothing. If the user wants project-local overrides, they create that file (keyless — never an API key in any repo file) and set the trust flag themselves.
- **Reject the global install paths.** Upstream documents `npx skills add … -g` and `/plugin marketplace add` — both write to `~/.claude`, which this skill's norms avoid. The vendored project copy above is the only supported install.

#### Phase 3: Verify

Check every guarantee the workflow promises, one line per component (this is also the Phase 1 detect check):

```bash
<skill-root>/../do-work/tools/do-work-cli.sh --repo-root <project-root> --format json install-last30days check <project-root>
```

Missing, failed, or malformed canonical tooling stops actionably; do not fall back to the compatibility script, a fresh clone, or direct filesystem mutation.

Report "Installed successfully" only when no line prints FAILED. The engine resolves a Python 3.12+ interpreter at run time (upstream keeps it in `LAST30DAYS_PYTHON`), so no qualifying interpreter is a real failure, not a warning — report the install as failed and name Python 3.12+ as the missing piece. A FAILED ignore line means the vendored ~15 MB is committable in the consuming repo — that's a broken install even though the engine itself would run.

#### Phase 4: Report back

```
Installed: last30days skill (vendored)

Destination: <project-root>/.claude/skills/last30days/
  (ignored via .git/info/exclude — machine-local; your .gitignore is untouched)

- Auto-discovers as the /last30days slash command.
- Reddit, Hacker News, Polymarket, GitHub, and YouTube work with no API keys.
- X/TikTok/Instagram need keys in the user-global ~/.config/last30days/.env
  — never in project files.
- Research memory defaults to ~/Documents/Last30Days/ (outside this repo).
  To relocate it, set LAST30DAYS_MEMORY_DIR in your environment or the
  user-global config — if you point it inside this repo, add an ignore
  rule for that path too.

Usage doctrine — when it's appropriate to invoke and what NOT to use it
for — is this project's responsibility to document. Add it wherever the
project keeps its action-usage docs.
```

---

## Help Block (no/unknown target)

Print this menu and stop:

```text
do-work-toolbox install ui-design      Frontend-design skill for production-grade UI
do-work-toolbox install bowser         Playwright CLI + Bowser browser skill
do-work-toolbox install last30days     Project-scoped social-research engine
do-work-toolbox install ideation-adhd  Parallel cognitive-frame ideation skill
```
## Output Format

- **`ui-design`**, **`ideation-adhd`**: a short status line — "already installed", "installed successfully", or an error describing what failed and how to finish manually.
- **`bowser`**: a two-line status — one for `playwright-cli`, one for the Bowser skill. Each is either "OK" (installed and verified), "already installed" (detected in Phase 1), or an error with the exact command the user can re-run.
- **`last30days`**: a per-guarantee status (skill file, ignore rule, Python 3.12+) — "already installed" only when every guarantee holds; otherwise "installed successfully" with the destination path, or the FAILED line(s) and the exact command the user can re-run.
- **Unknown / missing target**: the help block above.

## Rules

- **Install to the project, not globally.** Skill files go under `<project-root>/.claude/skills/<skill-name>/` (`<project-root>` resolved via `git rev-parse --show-toplevel || pwd`). Do not write to `~/.claude/` or any global path.
- **CLI is global; skill is project-scoped.** For `bowser`, `playwright-cli` goes to the global npm prefix; the Bowser skill goes under `<project-root>/.claude/skills/playwright-bowser/`. Don't mix them.
- **Never overwrite an existing non-empty skill `SKILL.md`.** Phase 1 of each workflow is the gate. If the file is present and non-empty, stop. (A zero-byte file is a failed download — reinstalling over it is repair, not overwrite, which is why the detect commands use `test -s` rather than `ls`.)
- **Only Chromium by default (bowser).** Other browsers bloat install time and aren't needed for ui-review's default flow.
- **Don't silently substitute a different skill or repo.** If the upstream URL fails, report the error — don't download a similarly-named skill from elsewhere.
- **Keyless in the project (last30days).** This install writes no config file at all. If a project-local `.claude/last30days.env` ever exists, it must never contain API keys — real keys live only in the user-global `~/.config/last30days/.env`. Never write a secret into any file inside the repo.
- **The vendor drop must be ignored (last30days).** Phase 2 adds `**/.claude/skills/last30days/` to the enclosing repo's `.git/info/exclude` when it isn't already covered — machine-local, never the project's committable `.gitignore` — because ~15 MB of upstream Python must never become committable in the consuming repo.
- **One target per invocation.** If the user wants both, they run two separate commands. The action never chains targets.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The file is already there but looks stale — I'll overwrite it" | Report installed; let the user decide whether to re-download | Overwriting a user-customized skill file silently destroys their edits |
| "`curl` failed, I'll fetch it with `wget` from a mirror" | Report the failure and the URL; let the user install manually | Unknown mirrors risk installing a tampered skill file |
| "The user wanted UI help, I should install the other targets while I'm here" | Stop after the requested target; mention the others as next steps if relevant | Each install target has a single, documented scope |
| "I'll install all three browsers to be safe (bowser)" | Install only Chromium; mention the manual command for the others | Firefox + WebKit roughly triple the install time and disk use, for a feature ui-review doesn't need |
| "npm install failed, I'll try yarn and pnpm and bun until something works (bowser)" | Try npm, then yarn; if both fail, stop and report | Quiet package-manager shopping leaves the user unsure what got installed |
| "Upstream's README says `npx skills add … -g` — I'll just follow upstream (last30days)" | Vendor into `$PROJECT_ROOT/.claude/skills/last30days/` per the workflow | Both `-g` and `/plugin marketplace add` write to `~/.claude`, which this skill never touches |
| "The user gave me an X API key — I'll put it in `.claude/last30days.env` so it works right away" | Direct them to `~/.config/last30days/.env`; never write a key into a project file | A key in a repo file leaks on the next commit |
| "I'll write `LAST30DAYS_TRUST_PROJECT_CONFIG=1` into `.claude/last30days.env` so project config just works (last30days)" | Write no config file; tell the user the trust flag must come from their environment or user-global config | The flag gates whether the engine reads the project file — setting it *inside* that file is circular and inert |

## Red Flags

- The install command reported success but the verify step shows the file is empty — the URL may have changed; investigate before claiming success.
- `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists but has zero size and Phase 1 still reported "already installed" — the detect command regressed to a bare existence check (`ls`); it must be `test -s` so a re-run repairs the failed download.
- An install command writes `SKILL.md` directly instead of following the canonical [Atomic download publication](../../do-work/docs/prescribed-shell-primitives.md#atomic-download-publication) shape — restore the temporary sibling, rename-on-success, and failure-preserving cleanup.
- A stray `SKILL.md.download` file sits next to an installed skill — a prior run's failure cleanup didn't run; delete it (the rename-on-success flow never leaves one behind).
- You installed into `~/.claude/skills/` instead of the project — undo and re-install to the correct path.
- `git rev-parse --show-toplevel` fails (not in a git repo) and you installed into `pwd` — acceptable, but warn the user the path may drift if they `cd` elsewhere.
- (bowser) `playwright-cli --help` succeeds but `playwright-cli install` fails silently — browsers aren't actually installed; headless runs will error later.
- (bowser) You installed `playwright-cli` into a project-local `node_modules` instead of globally — the CLI won't be on PATH for other sessions.
- (last30days) `git check-ignore -q .claude/skills/last30days/SKILL.md` exits non-zero in a git repo — the exclude entry was skipped or mismatched; fix it before anything gets committed. (Don't eyeball `git status` for this: a wholly-untracked `.claude/` collapses to a single `?? .claude/` row that hides the path either way.)
- (last30days) `SKILL.md` exists but `<project-root>/.claude/skills/last30days/scripts/last30days.py` is absent or empty — the runtime payload is incomplete; re-run the transactional install instead of treating the sentinel as healthy.
- (last30days) A project file (e.g. `.claude/last30days.env`) contains anything that looks like a credential — remove it and move the key to the user-global `~/.config/last30days/.env`.
- (last30days) Verify found no Python 3.12+ interpreter — the engine can't run; treat it as a failed install, not a soft warning.

## Verification Checklist

- [ ] Step 1 correctly dispatched on `$ARGUMENTS` (known target → workflow; unknown/empty → help block).
- [ ] Phase 1 detected an existing install and stopped, OR Phase 2+ ran the fetch/install commands.
- [ ] After the verify phase, `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists and is non-empty (skill-file targets: `ui-design`, `bowser`, `last30days`, `ideation-adhd`).
- [ ] (bowser only) `playwright-cli --help` runs without error and Chromium is installed.
- [ ] (last30days only) `<project-root>/.claude/skills/last30days/scripts/last30days.py` is non-empty, a Python 3.12+ interpreter is on PATH, `git check-ignore` covers `.claude/skills/last30days/`, and no project file gained an API key.
- [ ] The report names the destination path so the user can verify location.
