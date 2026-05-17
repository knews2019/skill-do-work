# Install Action

> **Part of the do-work skill.** Installs companion skills/tooling into the current project. Currently supports two targets: `ui-design` (frontend-design skill) and `bowser` (Playwright CLI + Bowser skill for browser automation).

Each target is idempotent — running it when the target is already present is a no-op. The action dispatches on the first argument; everything else (detect → install → verify → report) follows the same shape.

## When to Use

**Use when:**
- The project is about to start UI work or browser automation and the matching companion skill/tooling is missing.
- A `ui-review` pass flagged "frontend-design skill not installed" or "visual verification skipped" — install the missing piece.
- A `domain: ui-design` REQ is about to be built and the builder would benefit from skill-level design knowledge (`install ui-design`).
- The user asked for headed-browser workflows, screenshots, or visual verification (`install bowser`).

**Do NOT use when:**
- The target is already installed (Step 1 of the matching workflow detects this and exits).
- The project explicitly uses a different design system or browser-automation tool and adding the do-work default would conflict.
- The environment can't install global npm packages (for `bowser`) and the user hasn't consented to a local-only install.

## Input

`$ARGUMENTS` selects the install target:

- `ui-design` — Install Anthropic's `frontend-design` skill for production-grade UI design capabilities.
- `bowser` — Install Playwright CLI (global) plus the Bowser skill (project-scoped) for browser automation, screenshots, and visual UI verification.

If `$ARGUMENTS` is empty or doesn't match a known target, print the help block (target list + one-line blurb each) and stop.

## Install Manifest

Every target follows the same four-step shape (detect → install → verify → report). The per-target commands and blurbs live here:

| target | detect_cmd | install_cmd | verify_cmd | blurb |
|--------|------------|-------------|------------|-------|
| `ui-design` | `ls "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" 2>/dev/null` | `mkdir -p "$PROJECT_ROOT/.claude/skills/frontend-design" && curl -fsSL -o "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" https://raw.githubusercontent.com/anthropics/claude-code/main/skills/frontend-design/SKILL.md` | `test -s "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" && echo "Installed successfully" || echo "Installation failed"` | Anthropic's `frontend-design` Claude skill — production-grade UI design capabilities (typography, color, spacing, layout, component design, responsive/mobile-first, accessibility). |
| `bowser` | `playwright-cli --help >/dev/null 2>&1 && ls "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" 2>/dev/null` | (multi-step — see `bowser` workflow below) | (multi-step — see `bowser` workflow below) | Playwright CLI + Bowser skill — headed/headless browser sessions with Chromium, screenshots at any viewport, DOM snapshots, parallel named sessions, persistent profiles. |

In every command above, resolve `PROJECT_ROOT` first:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
```

## Steps

### Step 1: Dispatch on `$ARGUMENTS`

- If `$ARGUMENTS` is empty, prints "help" / "?", or doesn't match a known target → print the help block (target list + blurb) and stop.
- If `$ARGUMENTS` matches a known target → proceed to the target-specific workflow below.

### Step 2: Run the target's workflow

Each workflow follows the same four phases. The `ui-design` workflow uses the manifest commands directly. The `bowser` workflow has multi-part install/verify and is spelled out below.

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
  https://raw.githubusercontent.com/disler/bowser/main/SKILL.md
```

Fallback path if the primary URL 404s:

```bash
curl -fsSL -o "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" \
  https://raw.githubusercontent.com/disler/bowser/main/skills/playwright-bowser/SKILL.md
```

If both fail, report the error and direct the user to https://github.com/disler/bowser for manual install.

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
detected, ui-review automatically runs visual verification.

To use directly: playwright-cli -s=my-session open https://example.com --persistent
```

---

## Help Block (no/unknown target)

When `$ARGUMENTS` is empty or doesn't match a known target, print:

```
install — install companion skills/tooling into the current project

  do-work install ui-design   Anthropic's frontend-design skill for production-grade UI
  do-work install bowser      Playwright CLI + Bowser skill for browser automation
```

Then stop.

## Output Format

- **`ui-design`**: a short status line — "already installed", "installed successfully", or an error describing what failed and how to finish manually.
- **`bowser`**: a two-line status — one for `playwright-cli`, one for the Bowser skill. Each is either "OK" (installed and verified), "already installed" (detected in Phase 1), or an error with the exact command the user can re-run.
- **Unknown / missing target**: the help block above.

## Rules

- **Install to the project, not globally.** Skill files go under `<project-root>/.claude/skills/<skill-name>/` (`<project-root>` resolved via `git rev-parse --show-toplevel || pwd`). Do not write to `~/.claude/` or any global path.
- **CLI is global; skill is project-scoped.** For `bowser`, `playwright-cli` goes to the global npm prefix; the Bowser skill goes under `<project-root>/.claude/skills/playwright-bowser/`. Don't mix them.
- **Never overwrite an existing skill `SKILL.md`.** Phase 1 of each workflow is the gate. If the file is present, stop.
- **Only Chromium by default (bowser).** Other browsers bloat install time and aren't needed for ui-review's default flow.
- **Don't silently substitute a different skill or repo.** If the upstream URL fails, report the error — don't download a similarly-named skill from elsewhere.
- **One target per invocation.** If the user wants both, they run two separate commands. The action never chains targets.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The file is already there but looks stale — I'll overwrite it" | Report installed; let the user decide whether to re-download | Overwriting a user-customized skill file silently destroys their edits |
| "`curl` failed, I'll fetch it with `wget` from a mirror" | Report the failure and the URL; let the user install manually | Unknown mirrors risk installing a tampered skill file |
| "The user wanted UI help, I should install both targets while I'm here" | Stop after the requested target; mention the other as a next step if relevant | Each install target has a single, documented scope |
| "I'll install all three browsers to be safe (bowser)" | Install only Chromium; mention the manual command for the others | Firefox + WebKit roughly triple the install time and disk use, for a feature ui-review doesn't need |
| "npm install failed, I'll try yarn and pnpm and bun until something works (bowser)" | Try npm, then yarn; if both fail, stop and report | Quiet package-manager shopping leaves the user unsure what got installed |

## Red Flags

- The install command reported success but the verify step shows the file is empty — the URL may have changed; investigate before claiming success.
- `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists but has zero size — treat as not-installed and re-download (with user confirmation).
- You installed into `~/.claude/skills/` instead of the project — undo and re-install to the correct path.
- `git rev-parse --show-toplevel` fails (not in a git repo) and you installed into `pwd` — acceptable, but warn the user the path may drift if they `cd` elsewhere.
- (bowser) `playwright-cli --help` succeeds but `playwright-cli install` fails silently — browsers aren't actually installed; headless runs will error later.
- (bowser) You installed `playwright-cli` into a project-local `node_modules` instead of globally — the CLI won't be on PATH for other sessions.

## Verification Checklist

- [ ] Step 1 correctly dispatched on `$ARGUMENTS` (known target → workflow; unknown/empty → help block).
- [ ] Phase 1 detected an existing install and stopped, OR Phase 2+ ran the fetch/install commands.
- [ ] After the verify phase, `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists and is non-empty.
- [ ] (bowser only) `playwright-cli --help` runs without error and Chromium is installed.
- [ ] The report names the destination path so the user can verify location.
- [ ] No changes were made outside `<project-root>/.claude/skills/<skill-name>/` (and, for `bowser`, the global npm install).
