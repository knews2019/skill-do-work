# Install Action

> **Part of the do-work skill.** Installs companion skills/tooling into the current project. Currently supports three targets: `ui-design` (frontend-design skill), `bowser` (Playwright CLI + Bowser skill for browser automation), and `last30days` (engagement-ranked social-research engine, vendored project-scoped and keyless).

Each target is idempotent — running it when the target is already present is a no-op. The action dispatches on the first argument; everything else (detect → install → verify → report) follows the same shape.

## When to Use

**Use when:**
- The project is about to start UI work or browser automation and the matching companion skill/tooling is missing.
- A `ui-review` pass flagged "frontend-design skill not installed" or "visual verification skipped" — install the missing piece.
- A `domain: ui-design` REQ is about to be built and the builder would benefit from skill-level design knowledge (`install ui-design`).
- The user asked for headed-browser workflows, screenshots, or visual verification (`install bowser`).
- The user asked for social research, trend scanning, or "what's the discourse on X" capabilities (`install last30days`).

**Do NOT use when:**
- The target is already installed (Step 1 of the matching workflow detects this and exits).
- The project explicitly uses a different design system or browser-automation tool and adding the do-work default would conflict.
- The environment can't install global npm packages (for `bowser`) and the user hasn't consented to a local-only install.

## Input

`$ARGUMENTS` selects the install target:

- `ui-design` — Install Anthropic's `frontend-design` skill for production-grade UI design capabilities.
- `bowser` — Install Playwright CLI (global) plus the Bowser skill (project-scoped) for browser automation, screenshots, and visual UI verification.
- `last30days` — Vendor the engagement-ranked social-research engine (project-scoped, git-ignored, keyless).

If `$ARGUMENTS` is empty or doesn't match a known target, print the help block (target list + one-line blurb each) and stop.

## Install Manifest

Every target follows the same four-step shape (detect → install → verify → report). The per-target commands and blurbs live here:

| target | detect_cmd | install_cmd | verify_cmd | blurb |
|--------|------------|-------------|------------|-------|
| `ui-design` | `ls "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" 2>/dev/null` | `mkdir -p "$PROJECT_ROOT/.claude/skills/frontend-design" && curl -fsSL -o "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" https://raw.githubusercontent.com/anthropics/skills/main/skills/frontend-design/SKILL.md` | `test -s "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" && echo "Installed successfully" || echo "Installation failed"` | Anthropic's `frontend-design` Claude skill — production-grade UI design capabilities (typography, color, spacing, layout, component design, responsive/mobile-first, accessibility). |
| `bowser` | `playwright-cli --help >/dev/null 2>&1 && ls "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" 2>/dev/null` | (multi-step — see `bowser` workflow below) | (multi-step — see `bowser` workflow below) | Playwright CLI + Bowser skill — headed/headless browser sessions with Chromium, screenshots at any viewport, DOM snapshots, parallel named sessions, persistent profiles. |
| `last30days` | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | (multi-step — see `last30days` workflow below) | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | Engagement-ranked social-research engine — Reddit/HN/Polymarket/GitHub/YouTube keyless out of the box; X/TikTok/Instagram unlock only via user-global API keys. |

In every command above, resolve `PROJECT_ROOT` first:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
```

## Steps

### Step 1: Dispatch on `$ARGUMENTS`

- If `$ARGUMENTS` is empty, prints "help" / "?", or doesn't match a known target → print the help block (target list + blurb) and stop.
- If `$ARGUMENTS` matches a known target → proceed to the target-specific workflow below.

### Step 2: Run the target's workflow

Each workflow follows the same four-step shape. The `ui-design` workflow uses the manifest commands directly. The `bowser` and `last30days` workflows have multi-part installs and are spelled out below.

---

## Workflow: `ui-design`

#### Phase 1: Check if already installed

Resolve `PROJECT_ROOT`, then run the manifest's `detect_cmd`. If the file exists, report "already installed" and stop.

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

## Workflow: `bowser`

The `bowser` target installs two components: the global `playwright-cli` (plus a Chromium browser), and the project-scoped `playwright-bowser` skill.

#### Phase 1: Check if already installed

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
playwright-cli --help >/dev/null 2>&1 && echo "playwright-cli: installed" || echo "playwright-cli: not found"
ls "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" 2>/dev/null && echo "bowser skill: installed" || echo "bowser skill: not found"
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

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
mkdir -p "$PROJECT_ROOT/.claude/skills/playwright-bowser"
curl -fsSL -o "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" \
  https://raw.githubusercontent.com/disler/bowser/main/.claude/skills/playwright-bowser/SKILL.md
```

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

Works alongside do-work's `ui-review` action — when Playwright CLI is
detected, ui-review automatically runs visual verification. The `ai-report`
action also consumes Playwright when available for live screenshots;
without it, ai-report falls back to SVG + Mermaid diagrams.

To use directly: playwright-cli -s=my-session open https://example.com --persistent
```

---

## Workflow: `last30days`

The `last30days` target vendors the engagement-ranked social-research engine (https://github.com/mvanhorn/last30days-skill) into the project as a git-ignored, keyless drop. Reddit, Hacker News, Polymarket, GitHub, and YouTube work with no API keys; X/TikTok/Instagram unlock only via user-global keys — never via project files.

#### Phase 1: Check if already installed

Run the full guarantee check from Phase 3 (same commands — skill file, ignore rule, Python 3.12+). The install promises all three; detecting on the skill file alone would let a half-completed prior run masquerade as installed.

- **All checks pass** (the ignore rule counts as passing when the project isn't a git repo) → report "already installed" and stop.
- **Skill file present but the ignore rule failed** → a prior run half-completed. Proceed to Phase 2 in *repair mode*: skip the clone/copy and run only the ignore step — it's guarded, so re-running is safe. (A missing Python 3.12+ interpreter isn't repairable by this action — report it per Phase 3.)
- **Skill file missing** → run Phase 2 in full.

#### Phase 2: Vendor the skill

The upstream repo keeps the actual skill at `skills/last30days/` (self-contained — `SKILL.md`, `scripts/`, and supporting directories). Shallow-clone to a temp dir, copy only that subdirectory's contents, discard the clone — skipped in repair mode, since the skill file already exists:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
if [ ! -s "$PROJECT_ROOT/.claude/skills/last30days/SKILL.md" ]; then
  CLONE_DIR="$(mktemp -d)"
  git clone --depth 1 https://github.com/mvanhorn/last30days-skill "$CLONE_DIR" \
    && mkdir -p "$PROJECT_ROOT/.claude/skills/last30days" \
    && cp -R "$CLONE_DIR/skills/last30days/." "$PROJECT_ROOT/.claude/skills/last30days/"
  COPY_STATUS=$?
  rm -rf "$CLONE_DIR"
  [ "$COPY_STATUS" -eq 0 ] || echo "last30days: clone/copy FAILED"
fi
```

If the block prints FAILED (offline, upstream repo moved), **stop here** — report the error and skip the ignore step below; a failed install must not leave stray side effects in the consuming repo. The `cp -R …/. ` form copies the *contents* into the destination, so re-running over a broken partial directory merges cleanly instead of nesting a second `last30days/` inside (Phase 1's skill-file gate keeps healthy installs from ever reaching this block).

Then make the ignore claim true — the vendored engine is ~15 MB of upstream Python that must never become committable in the consuming repo. Add it to the enclosing repo's `.git/info/exclude` (machine-local — never committed, never shipped); do **not** touch the project's committable `.gitignore`. This is the exact snippet from `crew-members/background-agents.md` step 1, substituting this path (see that file for why, including the linked-worktree caveat):

```bash
exclude=$(git rev-parse --git-path info/exclude 2>/dev/null) || exclude=""
if [ -n "$exclude" ]; then
  git check-ignore -q .claude/skills/last30days/SKILL.md 2>/dev/null \
    || echo '**/.claude/skills/last30days/' >> "$exclude"
fi
```

Two hard constraints on this phase:

- **Write no config file — anywhere.** The engine reads API keys and settings from the user-global `~/.config/last30days/.env`, which the user manages themselves. Upstream also supports a project-local `.claude/last30days.env`, but it is trust-gated: the engine only reads it when `LAST30DAYS_TRUST_PROJECT_CONFIG` is already set in the environment or the user-global config — writing the trust flag *inside* the project file it gates is circular and does nothing. If the user wants project-local overrides, they create that file (keyless — never an API key in any repo file) and set the trust flag themselves.
- **Reject the global install paths.** Upstream documents `npx skills add … -g` and `/plugin marketplace add` — both write to `~/.claude`, which this skill's norms avoid. The vendored project copy above is the only supported install.

#### Phase 3: Verify

Check every guarantee the workflow promises, one line per component (this is also the Phase 1 detect check):

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
test -s "$PROJECT_ROOT/.claude/skills/last30days/SKILL.md" && echo "skill file: OK" || echo "skill file: FAILED"
if git -C "$PROJECT_ROOT" rev-parse --git-dir >/dev/null 2>&1; then
  git -C "$PROJECT_ROOT" check-ignore -q .claude/skills/last30days/SKILL.md && echo "ignore rule: OK" || echo "ignore rule: FAILED"
else
  echo "ignore rule: n/a (not a git repo)"
fi
FOUND_PYTHON=""
for python_candidate in python3.13 python3.12 python3 python; do
  command -v "$python_candidate" >/dev/null 2>&1 \
    && "$python_candidate" -c 'import sys; raise SystemExit(0 if sys.version_info >= (3, 12) else 1)' 2>/dev/null \
    && { FOUND_PYTHON="$python_candidate"; break; }
done
[ -n "$FOUND_PYTHON" ] && echo "python 3.12+: OK ($FOUND_PYTHON)" || echo "python 3.12+: FAILED"
```

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

When `$ARGUMENTS` is empty or doesn't match a known target, print:

```
install — install companion skills/tooling into the current project

  do-work install ui-design   Anthropic's frontend-design skill for production-grade UI
  do-work install bowser      Playwright CLI + Bowser skill for browser automation
  do-work install last30days  Engagement-ranked social-research engine (vendored, keyless)
```

Then stop.

## Output Format

- **`ui-design`**: a short status line — "already installed", "installed successfully", or an error describing what failed and how to finish manually.
- **`bowser`**: a two-line status — one for `playwright-cli`, one for the Bowser skill. Each is either "OK" (installed and verified), "already installed" (detected in Phase 1), or an error with the exact command the user can re-run.
- **`last30days`**: a per-guarantee status (skill file, ignore rule, Python 3.12+) — "already installed" only when every guarantee holds; otherwise "installed successfully" with the destination path, or the FAILED line(s) and the exact command the user can re-run.
- **Unknown / missing target**: the help block above.

## Rules

- **Install to the project, not globally.** Skill files go under `<project-root>/.claude/skills/<skill-name>/` (`<project-root>` resolved via `git rev-parse --show-toplevel || pwd`). Do not write to `~/.claude/` or any global path.
- **CLI is global; skill is project-scoped.** For `bowser`, `playwright-cli` goes to the global npm prefix; the Bowser skill goes under `<project-root>/.claude/skills/playwright-bowser/`. Don't mix them.
- **Never overwrite an existing skill `SKILL.md`.** Phase 1 of each workflow is the gate. If the file is present, stop.
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
- `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists but has zero size — treat as not-installed and re-download (with user confirmation).
- You installed into `~/.claude/skills/` instead of the project — undo and re-install to the correct path.
- `git rev-parse --show-toplevel` fails (not in a git repo) and you installed into `pwd` — acceptable, but warn the user the path may drift if they `cd` elsewhere.
- (bowser) `playwright-cli --help` succeeds but `playwright-cli install` fails silently — browsers aren't actually installed; headless runs will error later.
- (bowser) You installed `playwright-cli` into a project-local `node_modules` instead of globally — the CLI won't be on PATH for other sessions.
- (last30days) `git check-ignore -q .claude/skills/last30days/SKILL.md` exits non-zero in a git repo — the exclude entry was skipped or mismatched; fix it before anything gets committed. (Don't eyeball `git status` for this: a wholly-untracked `.claude/` collapses to a single `?? .claude/` row that hides the path either way.)
- (last30days) A project file (e.g. `.claude/last30days.env`) contains anything that looks like a credential — remove it and move the key to the user-global `~/.config/last30days/.env`.
- (last30days) Verify found no Python 3.12+ interpreter — the engine can't run; treat it as a failed install, not a soft warning.

## Verification Checklist

- [ ] Step 1 correctly dispatched on `$ARGUMENTS` (known target → workflow; unknown/empty → help block).
- [ ] Phase 1 detected an existing install and stopped, OR Phase 2+ ran the fetch/install commands.
- [ ] After the verify phase, `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists and is non-empty.
- [ ] (bowser only) `playwright-cli --help` runs without error and Chromium is installed.
- [ ] (last30days only) a Python 3.12+ interpreter is on PATH, `git check-ignore` covers `.claude/skills/last30days/`, and no project file gained an API key.
- [ ] The report names the destination path so the user can verify location.
- [ ] No changes were made outside `<project-root>/.claude/skills/<skill-name>/` (plus, for `bowser`, the global npm install; for `last30days`, the machine-local `.git/info/exclude` entry).
