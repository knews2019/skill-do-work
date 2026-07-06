# Install Action

> **Part of the do-work skill.** Installs companion skills/tooling into the current project. Currently supports four targets: `ui-design` (frontend-design skill), `bowser` (Playwright CLI + Bowser skill for browser automation), `last30days` (engagement-ranked social-research engine, vendored project-scoped and keyless), and `just-kanban` (justfile recipes wiring `just run-kanban` to the shipped queue-kanban board).

Each target is idempotent — running it when the target is already present is a no-op. The action dispatches on the first argument; everything else (detect → install → verify → report) follows the same shape.

## When to Use

**Use when:**
- The project is about to start UI work or browser automation and the matching companion skill/tooling is missing.
- A `ui-review` pass flagged "frontend-design skill not installed" or "visual verification skipped" — install the missing piece.
- A `domain: ui-design` REQ is about to be built and the builder would benefit from skill-level design knowledge (`install ui-design`).
- The user asked for headed-browser workflows, screenshots, or visual verification (`install bowser`).
- The user asked for social research, trend scanning, or "what's the discourse on X" capabilities (`install last30days`).
- The user wants a standing `just run-kanban` shortcut so the board runs without invoking the agent (`install just-kanban`).

**Do NOT use when:**
- The target is already installed (Step 1 of the matching workflow detects this and exits).
- The project explicitly uses a different design system or browser-automation tool and adding the do-work default would conflict.
- The environment can't install global npm packages (for `bowser`) and the user hasn't consented to a local-only install.
- The user just wants to view the board once — that's `do-work board` (`actions/board.md`), no install needed.

## Input

`$ARGUMENTS` selects the install target:

- `ui-design` — Install Anthropic's `frontend-design` skill for production-grade UI design capabilities.
- `bowser` — Install Playwright CLI (global) plus the Bowser skill (project-scoped) for browser automation, screenshots, and visual UI verification.
- `last30days` — Vendor the engagement-ranked social-research engine (project-scoped, git-ignored, keyless).
- `just-kanban` — Append `just` recipes (`run-kanban`, `kanban-static`, `kanban-summary`) for the shipped queue-kanban board to the project's justfile.

If `$ARGUMENTS` is empty or doesn't match a known target, print the help block (target list + one-line blurb each) and stop.

## Install Manifest

Every target follows the same four-step shape (detect → install → verify → report). The per-target commands and blurbs live here:

| target | detect_cmd | install_cmd | verify_cmd | blurb |
|--------|------------|-------------|------------|-------|
| `ui-design` | `ls "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" 2>/dev/null` | `mkdir -p "$PROJECT_ROOT/.claude/skills/frontend-design" && curl -fsSL -o "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" https://raw.githubusercontent.com/anthropics/skills/main/skills/frontend-design/SKILL.md` | `test -s "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" && echo "Installed successfully" || echo "Installation failed"` | Anthropic's `frontend-design` Claude skill — production-grade UI design capabilities (typography, color, spacing, layout, component design, responsive/mobile-first, accessibility). |
| `bowser` | `playwright-cli --help >/dev/null 2>&1 && ls "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" 2>/dev/null` | (multi-step — see `bowser` workflow below) | (multi-step — see `bowser` workflow below) | Playwright CLI + Bowser skill — headed/headless browser sessions with Chromium, screenshots at any viewport, DOM snapshots, parallel named sessions, persistent profiles. |
| `last30days` | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | (multi-step — see `last30days` workflow below) | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | Engagement-ranked social-research engine — Reddit/HN/Polymarket/GitHub/YouTube keyless out of the box; X/TikTok/Instagram unlock only via user-global API keys. |
| `just-kanban` | (multi-step — see `just-kanban` workflow below) | (multi-step — see `just-kanban` workflow below; append-only) | (multi-step — see `just-kanban` workflow below) | Justfile recipes for the shipped queue-kanban board — `just run-kanban` serves the live board, `kanban-static`/`kanban-summary` cover the other modes; rebuilds the tool each run so `do-work update` refreshes take effect. |

In every command above, resolve `PROJECT_ROOT` first:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
```

## Steps

### Step 1: Dispatch on `$ARGUMENTS`

- If `$ARGUMENTS` is empty, prints "help" / "?", or doesn't match a known target → print the help block (target list + blurb) and stop.
- If `$ARGUMENTS` matches a known target → proceed to the target-specific workflow below.

### Step 2: Run the target's workflow

Each workflow follows the same four-step shape. The `ui-design` workflow uses the manifest commands directly. The `bowser`, `last30days`, and `just-kanban` workflows have multi-part installs and are spelled out below.

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

## Workflow: `just-kanban`

The `just-kanban` target appends [`just`](https://github.com/casey/just) recipes for the shipped queue-kanban board (`tools/queue-kanban/`, normally run via `actions/board.md`) to the consuming project's justfile, so `just run-kanban` serves the live board — replacing a stale queue-kanban instance still holding the port, then opening your default browser at it — without going through the agent. The justfile is a **project-owned** file — `do-work update` never touches it — while the tool source the recipes point at is refreshed by every update; the recipes rebuild the binary on each run so those refreshes take effect automatically.

#### Phase 1: Check if already installed

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
JUSTFILE_PATH=""
for justfile_candidate in justfile Justfile .justfile; do
  [ -f "$PROJECT_ROOT/$justfile_candidate" ] && { JUSTFILE_PATH="$PROJECT_ROOT/$justfile_candidate"; break; }
done
[ -n "$JUSTFILE_PATH" ] && grep -qE '^run-kanban[ :]' "$JUSTFILE_PATH" \
  && echo "run-kanban recipe: present" || echo "run-kanban recipe: absent"
```

If a `run-kanban` recipe is present — even one that differs from the block below — report "already installed" and stop. Never replace an existing recipe; if the user wants the shipped version, they remove theirs first.

**Upgrading an already-installed project:** because installs are append-only, a project that ran `just-kanban` before the recipe gained the replace-stale-instance and auto-open behavior (or any future recipe change) will not pick it up automatically — Phase 1 only checks for *presence* of a `run-kanban` recipe, not which version. To pull in the new behavior: delete the `# --- do-work board recipes ... ---` block from the project's justfile, then re-run `do-work install just-kanban`.

#### Phase 2: Append the recipes

Resolve the two paths as a text operation first, then substitute the literal result into the block below (deriving values before they enter a file is the same discipline CLAUDE.md prescribes for shell quoting):

1. `<skill-root>` — the absolute directory containing `SKILL.md` (this action lives in its `actions/` subdir).
2. **Global-install gate:** if `<skill-root>` is not inside `PROJECT_ROOT`, stop and report — a project justfile must not point outside the project, and this skill's norms reject global installs anyway.
3. `<kanban-dir>` — the `PROJECT_ROOT`-relative path of `<skill-root>/tools/queue-kanban` (e.g. `.claude/skills/do-work/tools/queue-kanban`).

Pick the justfile the same way Phase 1 does (`justfile` / `Justfile` / `.justfile` at `PROJECT_ROOT`, first match); when none exists, create `PROJECT_ROOT/justfile`. Append the block with your file-editing capability (or a quoted heredoc — `<<'RECIPES'` quoting keeps every token literal), adding one blank line of separation when the file already has content. Substitute `<kanban-dir>`; keep the `{{…}}` tokens **verbatim** — that is `just`'s own interpolation, resolved when the recipe runs, not by this install:

```just
# --- do-work board recipes (installed by `do-work install just-kanban`) ---

# Serve the do-work queue as a live Kanban board, replacing a stale instance on the port and opening your browser (Ctrl-C to stop; reload the page to refresh)
run-kanban $port="8090":
    case "$port" in ''|*[!0-9]*) echo "queue-kanban: invalid port '$port' - must be digits only (for a LAN-exposed host:port bind, run the queue-kanban serve command directly)" >&2; exit 1;; esac
    if command -v lsof >/dev/null 2>&1; then PID="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"; if [ -n "$PID" ]; then COMM="$(ps -p "$PID" -o comm= 2>/dev/null)"; COMM="${COMM##*/}"; if [ "$COMM" = "queue-kanban" ]; then kill "$PID" 2>/dev/null; i=0; while kill -0 "$PID" 2>/dev/null && [ "$i" -lt 20 ]; do sleep 0.1; i=$((i+1)); done; else echo "queue-kanban: port $port is already in use by another process ($COMM, pid $PID) - refusing to kill it. Stop it manually, or run 'just run-kanban <port>' with a different port." >&2; exit 1; fi; fi; fi
    cd <kanban-dir> && go build -o queue-kanban . && ./queue-kanban serve --open --repo-root "{{justfile_directory()}}" --port "$port"

# Shareable static snapshot → build/queue-kanban-board/index.html
kanban-static:
    cd <kanban-dir> && go build -o queue-kanban . && ./queue-kanban generate --out "{{justfile_directory()}}/build/queue-kanban-board" --repo-root "{{justfile_directory()}}"

# Column counts in the terminal, no browser
kanban-summary:
    cd <kanban-dir> && go build -o queue-kanban . && ./queue-kanban summary --repo-root "{{justfile_directory()}}"
```

Four deliberate choices in these recipes:

- **`$port` is an exported parameter, validated before anything else runs**: `just` interpolates `{{…}}` tokens textually into each recipe line's shell source, so a raw `{{port}}` would let `just run-kanban '8090; echo PWNED'` inject arbitrary commands (CLAUDE.md's never-interpolate-raw-user-text rule applies to justfiles too). The `$` prefix hands the parameter to every recipe line as an environment variable instead — the shell reads `"$port"` as data, never code — and the first line rejects anything but digits before the kill-stale or build+serve lines can see it.
- **`go build` on every run** (`actions/board.md` Step 4's rule): `do-work update` overwrites the tool's source but leaves the previously compiled binary in place — a cached binary silently renders old logic. The incremental rebuild is near-instant when nothing changed, and the binary stays uncommittable via the tool's shipped `.gitignore`.
- **Each `cd … && …` chain stays on one logical line**: `just` runs every recipe line in a fresh shell, so a bare `cd` on its own line would not carry into the next — the same cross-shell state trap CLAUDE.md documents for prescribed action steps.
- **The kill-stale check is its own recipe line and needs no `cd`**: it only touches `lsof`/`ps`/`kill` against `"$port"`, so it doesn't need the `<kanban-dir>` context the build+serve line does. `just` aborts a recipe on the first line that exits non-zero, so a squatting non-`queue-kanban` process's `exit 1` here stops the recipe *before* the build+serve line ever runs — no build is attempted and nothing gets killed. It kills only a process whose own command name (verified via `ps -p PID -o comm=`) is `queue-kanban`; anything else is left running and named in the error. A missing `lsof` degrades gracefully — the check is skipped and the recipe proceeds straight to build+serve — rather than blocking the recipe on a tool that isn't guaranteed to exist.

#### Phase 3: Verify

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
JUSTFILE_PATH=""
for justfile_candidate in justfile Justfile .justfile; do
  [ -f "$PROJECT_ROOT/$justfile_candidate" ] && { JUSTFILE_PATH="$PROJECT_ROOT/$justfile_candidate"; break; }
done
[ -n "$JUSTFILE_PATH" ] && grep -qE '^run-kanban[ :]' "$JUSTFILE_PATH" && echo "recipes: OK" || echo "recipes: FAILED"
if command -v just >/dev/null; then
  just --justfile "$JUSTFILE_PATH" --list >/dev/null 2>&1 && echo "justfile parses: OK" || echo "justfile parses: FAILED"
else
  echo "just: MISSING (recipes are installed; install just to run them)"
fi
command -v go >/dev/null && echo "go toolchain: OK" || echo "go toolchain: MISSING (the board needs Go to build — see tools/queue-kanban/go.mod)"
```

Report "Installed successfully" only when `recipes: OK` and — whenever `just` is available to check — `justfile parses: OK` (a FAILED parse means the append corrupted the file: restore it and re-append). Missing `just` or `go` are **warnings, not failures**: the recipes are inert text until run, and `actions/board.md` already treats a missing Go toolchain as a graceful stop rather than a blocker.

#### Phase 4: Report back

```
Installed: just recipes for the queue-kanban board

Appended to: <project-root>/justfile

  just run-kanban          Live board at http://localhost:8090, opens in your browser (reload to refresh)
  just run-kanban 9000     Same, custom port
  just kanban-static       Snapshot → build/queue-kanban-board/index.html
  just kanban-summary      Column counts in the terminal

- `just run-kanban` replaces a stale queue-kanban instance already holding
  the port and opens your default browser at the board URL automatically —
  a non-queue-kanban process on the port is left alone and named in an error.
- Recipes rebuild the tool on every run, so `do-work update` refreshes take
  effect automatically (needs the Go toolchain — same requirement as
  `do-work board`).
- The justfile is project-owned: `do-work update` never touches it.
```

---

## Help Block (no/unknown target)

When `$ARGUMENTS` is empty or doesn't match a known target, print:

```
install — install companion skills/tooling into the current project

  do-work install ui-design   Anthropic's frontend-design skill for production-grade UI
  do-work install bowser      Playwright CLI + Bowser skill for browser automation
  do-work install last30days  Engagement-ranked social-research engine (vendored, keyless)
  do-work install just-kanban  Justfile recipes for the queue-kanban board (needs Go to run)
```

Then stop.

## Output Format

- **`ui-design`**: a short status line — "already installed", "installed successfully", or an error describing what failed and how to finish manually.
- **`bowser`**: a two-line status — one for `playwright-cli`, one for the Bowser skill. Each is either "OK" (installed and verified), "already installed" (detected in Phase 1), or an error with the exact command the user can re-run.
- **`last30days`**: a per-guarantee status (skill file, ignore rule, Python 3.12+) — "already installed" only when every guarantee holds; otherwise "installed successfully" with the destination path, or the FAILED line(s) and the exact command the user can re-run.
- **`just-kanban`**: a per-component status (recipes appended, justfile parses, `just`/`go` availability) — "already installed" when a `run-kanban` recipe already exists; missing toolchains are warnings, not failures.
- **Unknown / missing target**: the help block above.

## Rules

- **Install to the project, not globally.** Skill files go under `<project-root>/.claude/skills/<skill-name>/` (`<project-root>` resolved via `git rev-parse --show-toplevel || pwd`). Do not write to `~/.claude/` or any global path.
- **CLI is global; skill is project-scoped.** For `bowser`, `playwright-cli` goes to the global npm prefix; the Bowser skill goes under `<project-root>/.claude/skills/playwright-bowser/`. Don't mix them.
- **Never overwrite an existing skill `SKILL.md`.** Phase 1 of each workflow is the gate. If the file is present, stop.
- **Only Chromium by default (bowser).** Other browsers bloat install time and aren't needed for ui-review's default flow.
- **Don't silently substitute a different skill or repo.** If the upstream URL fails, report the error — don't download a similarly-named skill from elsewhere.
- **Keyless in the project (last30days).** This install writes no config file at all. If a project-local `.claude/last30days.env` ever exists, it must never contain API keys — real keys live only in the user-global `~/.config/last30days/.env`. Never write a secret into any file inside the repo.
- **The vendor drop must be ignored (last30days).** Phase 2 adds `**/.claude/skills/last30days/` to the enclosing repo's `.git/info/exclude` when it isn't already covered — machine-local, never the project's committable `.gitignore` — because ~15 MB of upstream Python must never become committable in the consuming repo.
- **Append-only in the justfile (just-kanban).** Never reorder, reformat, or replace existing justfile content; an existing `run-kanban` recipe — even a divergent one — means "already installed", not "overwrite". Create a `justfile` only when none of `justfile`/`Justfile`/`.justfile` exists at the project root.
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
| "Their existing `run-kanban` recipe is outdated — I'll replace it with the shipped block (just-kanban)" | Report already-installed and show the shipped block for manual comparison | The existing recipe may carry deliberate project-specific flags; replacing it destroys their edits |
| "The skill is installed globally — I'll hard-code its absolute path into the recipe (just-kanban)" | Stop at the Phase 2 global-install gate and report | A recipe pointing outside the project breaks on every other clone and machine, and the skill's norms reject global installs |

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
- (just-kanban) The justfile diff shows anything beyond one appended block — existing recipes were reordered or rewritten; restore the file and re-append.
- (just-kanban) The appended recipe contains an absolute path (especially into `$HOME`) — the skill-root resolution went wrong; recipes must use project-relative paths.

## Verification Checklist

- [ ] Step 1 correctly dispatched on `$ARGUMENTS` (known target → workflow; unknown/empty → help block).
- [ ] Phase 1 detected an existing install and stopped, OR Phase 2+ ran the fetch/install commands.
- [ ] After the verify phase, `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists and is non-empty (skill-file targets: `ui-design`, `bowser`, `last30days`).
- [ ] (bowser only) `playwright-cli --help` runs without error and Chromium is installed.
- [ ] (last30days only) a Python 3.12+ interpreter is on PATH, `git check-ignore` covers `.claude/skills/last30days/`, and no project file gained an API key.
- [ ] (just-kanban only) the justfile gained exactly one appended block, `run-kanban` greps present, `just --list` parses when `just` is available, and no existing recipe was modified.
- [ ] The report names the destination path so the user can verify location.
- [ ] No changes were made outside `<project-root>/.claude/skills/<skill-name>/` (plus, for `bowser`, the global npm install; for `last30days`, the machine-local `.git/info/exclude` entry; for `just-kanban`, the project justfile).
