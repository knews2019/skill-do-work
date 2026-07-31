# Install Action

> **Part of the do-work skill.** Installs companion skills/tooling into the current project. Currently supports six targets: `ui-design` (frontend-design skill), `bowser` (Playwright CLI + Bowser skill for browser automation), `last30days` (engagement-ranked social-research engine, vendored project-scoped and keyless), `just-kanban` (justfile recipes wiring `just run-kanban` to the shipped queue-kanban board and `just run-do-work-update` to the project-local skill updater), `ideation-adhd` (parallel divergent-ideation skill), and `memory-module` (the ADR-017 memory engine: `memory/` scaffolding plus optional SessionStart/Stop hooks).

Each target is idempotent — running it when the target is already present and current is a no-op. One target goes further: `just-kanban` compares an already-present recipe block against the shipped version and offers a consent-gated upgrade when they diverge (Phase 1b of its workflow). The action dispatches on the first argument; everything else (detect → install → verify → report) follows the same shape.

## When to Use

**Use when:**
- The project is about to start UI work or browser automation and the matching companion skill/tooling is missing.
- A `ui-review` pass flagged "frontend-design skill not installed" or "visual verification skipped" — install the missing piece.
- A `domain: ui-design` REQ is about to be built and the builder would benefit from skill-level design knowledge (`install ui-design`).
- The user asked for headed-browser workflows, screenshots, or visual verification (`install bowser`).
- The user asked for social research, trend scanning, or "what's the discourse on X" capabilities (`install last30days`).
- The user wants a standing `just run-kanban` shortcut so the board runs without invoking the agent (`install just-kanban`).
- The user wants structured divergent ideation — deliberately unconventional candidate ideas for a named problem, beyond what `scan-ideas` surfaces from the repo (`install ideation-adhd`).
- The user wants session-persistent memory — `do-work memory` scaffolding and, on Claude Code, the auto-inject/auto-capture hooks (`install memory-module`).

**Do NOT use when:**
- The target is already installed and current (Phase 1 of the matching workflow detects this and exits; for `just-kanban`, an outdated recipe block gets a diff and a consent-gated upgrade instead of a plain exit).
- The project explicitly uses a different design system or browser-automation tool and adding the do-work default would conflict.
- The environment can't install global npm packages (for `bowser`) and the user hasn't consented to a local-only install.
- The user just wants to view the board once — that's `do-work board` (`actions/board.md`), no install needed.

## Input

`$ARGUMENTS` selects the install target:

- `ui-design` — Install Anthropic's `frontend-design` skill for production-grade UI design capabilities.
- `bowser` — Install Playwright CLI (global) plus the Bowser skill (project-scoped) for browser automation, screenshots, and visual UI verification.
- `last30days` — Vendor the engagement-ranked social-research engine (project-scoped, git-ignored, keyless).
- `just-kanban` — Append `just` recipes (`run-kanban`, `kanban-static`, `kanban-summary`, `run-do-work-update`) for the shipped queue-kanban board and project-local do-work updates to the project's justfile; if the recipes are already present but diverge from the shipped block, offer a consent-gated upgrade. (`run-do-work-update` is an accepted alias for this target.)
- `ideation-adhd` — Install the upstream `adhd` skill (project-scoped) for parallel divergent ideation: isolated branches under distinct cognitive frames, scored, clustered, and deepened. (`adhd` is an accepted alias for this target.)
- `memory-module` — Scaffold the `memory/` store (working-memory template, logs dir, usage ledger) and merge the memory SessionStart/Stop hook entries into `.claude/settings.json` — composing with, never clobbering, existing hooks.

If `$ARGUMENTS` is empty or doesn't match a known target, print the help block (target list + one-line blurb each) and stop.

## Install Manifest

Every target follows the same four-step shape (detect → install → verify → report). The per-target commands and blurbs live here:

| target | detect_cmd | install_cmd | verify_cmd | blurb |
|--------|------------|-------------|------------|-------|
| `ui-design` | `test -s "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md"` | `mkdir -p "$PROJECT_ROOT/.claude/skills/frontend-design" && curl -fsSL -o "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md.download" https://raw.githubusercontent.com/anthropics/skills/main/skills/frontend-design/SKILL.md && mv "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md.download" "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" || { rm -f "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md.download"; false; }` | `test -s "$PROJECT_ROOT/.claude/skills/frontend-design/SKILL.md" && echo "Installed successfully" || echo "Installation failed"` | Anthropic's `frontend-design` Claude skill — production-grade UI design capabilities (typography, color, spacing, layout, component design, responsive/mobile-first, accessibility). |
| `bowser` | `playwright-cli --help >/dev/null 2>&1 && test -s "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md"` | (multi-step — see `bowser` workflow below) | (multi-step — see `bowser` workflow below) | Playwright CLI + Bowser skill — headed/headless browser sessions with Chromium, screenshots at any viewport, DOM snapshots, parallel named sessions, persistent profiles. |
| `last30days` | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | (multi-step — see `last30days` workflow below) | (multi-step — see `last30days` workflow below; gates on the full guarantee set) | Engagement-ranked social-research engine — Reddit/HN/Polymarket/GitHub/YouTube keyless out of the box; X/TikTok/Instagram unlock only via user-global API keys. |
| `just-kanban` | (multi-step — see `just-kanban` workflow below) | (multi-step — see `just-kanban` workflow below; appends fresh, consent-gated upgrade when present-but-divergent) | (multi-step — see `just-kanban` workflow below) | Justfile recipes for the shipped queue-kanban board and project-local updater — `just run-kanban` serves the live board, `kanban-static`/`kanban-summary` cover the other modes, and `run-do-work-update` runs the guarded updater without an agent. |
| `memory-module` | (multi-step — see `memory-module` workflow below; gates on scaffolding + hook wiring) | (multi-step — see `memory-module` workflow below; repair mode when scaffolding exists but hooks are absent) | (multi-step — see `memory-module` workflow below) | Hermes-style working-memory + dated-logs engine with SessionStart/Stop hooks and layered recall — the experimental counterpart to `actions/bkb.md` (see `decisions/records/adr-017-run-a-parallel-memory-engine-experiment-with-usage-ledgers.md`). |
| `ideation-adhd` | `test -s "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md"` | `mkdir -p "$PROJECT_ROOT/.claude/skills/adhd" && curl -fsSL -o "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md.download" https://raw.githubusercontent.com/UditAkhourii/adhd/main/skills/adhd/SKILL.md && mv "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md.download" "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md" || { rm -f "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md.download"; false; }` | `test -s "$PROJECT_ROOT/.claude/skills/adhd/SKILL.md" && echo "Installed successfully" || echo "Installation failed"` | The `adhd` skill (MIT) — parallel divergent ideation: spawns isolated branches under distinct cognitive frames (regulator, biology, speedrunner, 10-year-old, $0 budget, …), scores on novelty/viability/fit, clusters, prunes traps, deepens the top survivors. Explicitly invoked (`/adhd`), never fires on its own. |

In every command above, resolve `PROJECT_ROOT` first:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
```

Every skill-file download lands in a `SKILL.md.download` temp name and is renamed into place only on success, with the temp removed on failure. `curl -o` writes the final path incrementally, so a mid-transfer failure would otherwise leave a non-empty partial file that `test -s` cannot distinguish from a complete install — the same trap the zero-byte detect fix closed, one failure mode further along. (`curl --remove-on-error` also covers this but only exists in curl ≥ 7.83; the rename works everywhere.)

## Steps

### Step 1: Dispatch on `$ARGUMENTS`

- If `$ARGUMENTS` is empty, prints "help" / "?", or doesn't match a known target → print the help block (target list + blurb) and stop.
- If `$ARGUMENTS` matches a known target → proceed to the target-specific workflow below.

### Step 2: Run the target's workflow

Each workflow follows the same four-step shape. The `ui-design` and `ideation-adhd` workflows use the manifest commands directly. The `bowser`, `last30days`, `just-kanban`, and `memory-module` workflows have multi-part installs and are spelled out below.

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

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
mkdir -p "$PROJECT_ROOT/.claude/skills/playwright-bowser"
curl -fsSL -o "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md.download" \
  https://raw.githubusercontent.com/disler/bowser/main/.claude/skills/playwright-bowser/SKILL.md \
  && mv "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md.download" \
        "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" \
  || { rm -f "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md.download"; false; }
```

The trailing `false` is load-bearing: `rm -f` on an absent path exits 0, so a plain `|| rm -f …` would hand the block a success status after a failed download and the phase would report installed. Same form as the `ui-design` and `ideation-adhd` manifest rows.

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

Then make the ignore claim true — the vendored engine is ~15 MB of upstream Python that must never become committable in the consuming repo. Add it to the enclosing repo's `.git/info/exclude` (machine-local — never committed, never shipped); do **not** touch the project's committable `.gitignore`. This is the exact snippet from `crew-members/background-agents.md` → **Local-ignore snippet (for genuinely-transient paths)**, substituting this path (see that file for why, including the linked-worktree caveat):

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

The `just-kanban` target appends [`just`](https://github.com/casey/just) recipes for the shipped queue-kanban board (`tools/queue-kanban/`, normally run via `actions/board.md`) and a project-local do-work updater to the consuming project's justfile. `just run-kanban` serves the live board — replacing a stale queue-kanban instance still holding the port, then opening your default browser at it — while `just run-do-work-update` runs the shipped updater without consuming an agent turn. The updater still reviews the installed-versus-upstream diff, asks before overwriting, and excludes the runtime `do-work/` directory — but it keeps **no** rollback copy: version control is the undo, and duplicating a tracked tree on every run buys nothing git does not already hold. A failure partway therefore leaves a partial install and prints the exact `git checkout` / `git clean` commands for the shipped paths git actually tracks; when the install is *not* tracked in git the updater says so before the confirmation prompt, because there nothing can be restored at all. It also lists, before the prompt, any file sitting in your install that upstream no longer ships — the extraction overwrites but never deletes, so a removed action or check would otherwise stay readable downstream forever. Those are reported, never removed: at that level a file your project put in a shipped directory looks identical to a stale one. The justfile is a **project-owned** file — `do-work update` never touches it, and `tools/do-work-update.sh` enforces that by name: when the skill is installed at the project root (the recipe's fallback branch, where the tarball's own `justfile` would land on top of the project's), it records which of `justfile`/`Justfile`/`.justfile` the project uses, keeps it out of the reviewed shipped set, and restores that exact name and content across the extraction. A nested install's `justfile` is the skill's own copy and *is* refreshed like any shipped file. Meanwhile the tool source the recipes point at is refreshed by every update; the recipes rebuild the binary on each run so those refreshes take effect automatically. Changes to the recipe *text* itself don't propagate that way, which is why re-running this install on an already-installed project compares and offers an upgrade (Phase 1b) instead of stopping at "already installed".

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

If no `run-kanban` recipe is present, skip to Phase 2 (fresh append). If one is present, do **not** stop at "already installed" — presence alone doesn't mean current. Continue to Phase 1b.

#### Phase 1b: Compare and upgrade (recipe already present)

A project that installed before a recipe change (auto-open, stale-instance replacement, port validation, or any future one) still carries the old text — the justfile is project-owned, so nothing else ever refreshes it. Compare before deciding:

1. Render the current shipped block: resolve `<kanban-dir>` per Phase 2's path-resolution steps and substitute it into Phase 2's recipe block.
2. Extract the installed do-work recipes from the justfile: for each of `run-kanban`, `kanban-static`, `kanban-summary`, `run-do-work-update`, the span from the comment line(s) immediately above its `name …:` header line through its last indented line. A recipe that is missing entirely counts as divergent.
3. Compare each extracted recipe against its rendered shipped version, ignoring trailing whitespace.

Then branch on the result:

- **All four identical** → report "already installed (current)" and stop.
- **Any divergence** → show the user a unified diff (installed vs. shipped) and ask with your environment's ask-user prompt whether to upgrade. This consent gate is the load-bearing safety step, not a formality: the divergence may be an older shipped version, but it may equally be deliberate project-specific edits — only the user can tell those apart, and the diff is what lets them.
  - **User accepts** → replace each divergent recipe span in place with its shipped version, and append any of the four that were missing (inside the `# --- do-work board recipes … ---` block when that header comment exists, else at the end of the file). Touch nothing else in the justfile — no reordering, no reformatting of the user's own recipes or variables. Then continue to Phase 3 (verify).
  - **User declines** → stop and report "kept existing recipes", naming the shipped behaviors their version lacks so the choice is informed, not final by accident.

#### Phase 2: Append the recipes

Resolve the two paths as a text operation first, then substitute the literal result into the block below (derive values *before* they enter a file, the same discipline as never interpolating raw text inside shell quoting):

1. `<skill-root>` — the absolute directory containing `SKILL.md` (this action lives in its `actions/` subdir).
2. **Global-install gate:** if `<skill-root>` is not inside `PROJECT_ROOT`, stop and report — a project justfile must not point outside the project, and this skill's norms reject global installs anyway.
3. `<kanban-dir>` — the `PROJECT_ROOT`-relative path of `<skill-root>/tools/queue-kanban` (e.g. `.claude/skills/do-work/tools/queue-kanban`).

Pick the justfile the same way Phase 1 does (`justfile` / `Justfile` / `.justfile` at `PROJECT_ROOT`, first match); when none exists, create `PROJECT_ROOT/justfile`. Append the block with your file-editing capability (or a quoted heredoc — `<<'RECIPES'` quoting keeps every token literal), adding one blank line of separation when the file already has content. Substitute `<kanban-dir>`; keep the `{{…}}` tokens **verbatim** — that is `just`'s own interpolation, resolved when the recipe runs, not by this install:

```just
# --- do-work board recipes (installed by `do-work install just-kanban`) ---

# Serve the do-work queue as a live Kanban board, replacing a stale instance on the port and opening your browser (Ctrl-C to stop; reload the page to refresh)
run-kanban $port="8090":
    case "$port" in ''|*[!0-9]*) echo "queue-kanban: invalid port '$port' - must be digits only (for a LAN-exposed host:port bind, run the queue-kanban serve command directly)" >&2; exit 1;; esac
    if command -v lsof >/dev/null 2>&1; then listener_pid="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"; if [ -n "$listener_pid" ]; then listener_executable="$(lsof -a -p "$listener_pid" -d txt -Fn 2>/dev/null | sed -n 's/^n//p' | head -n1)"; listener_executable_name="${listener_executable##*/}"; listener_command="$(ps -p "$listener_pid" -o args= 2>/dev/null)"; case "$listener_executable_name" in *queue-kanban*) echo "queue-kanban: stopping previous session on :$port (pid $listener_pid): $listener_command"; kill "$listener_pid" 2>/dev/null; wait_count=0; while kill -0 "$listener_pid" 2>/dev/null && [ "$wait_count" -lt 20 ]; do sleep 0.1; wait_count=$((wait_count+1)); done;; *) echo "queue-kanban: port $port is already in use by another process ($listener_command, pid $listener_pid) - refusing to kill it. Stop it manually, or run 'just run-kanban <port>' with a different port." >&2; exit 1;; esac; fi; fi
    cd <kanban-dir> && go build -o queue-kanban . && ./queue-kanban serve --open --repo-root "{{justfile_directory()}}" --port "$port"

# Shareable static snapshot → build/queue-kanban-board/index.html (locally git-excluded so it never dirties git status)
kanban-static:
    cd <kanban-dir> && go build -o queue-kanban . && ./queue-kanban generate --out "{{justfile_directory()}}/build/queue-kanban-board" --repo-root "{{justfile_directory()}}"
    cd "{{justfile_directory()}}" && if git rev-parse --git-dir >/dev/null 2>&1 && ! git check-ignore -q build/queue-kanban-board/index.html; then exclude_file="$(git rev-parse --git-path info/exclude)"; mkdir -p "$(dirname "$exclude_file")"; echo '/build/queue-kanban-board/' >> "$exclude_file"; echo "kanban-static: added /build/queue-kanban-board/ to .git/info/exclude (local-only ignore)"; fi

# Column counts in the terminal, no browser
kanban-summary:
    cd <kanban-dir> && go build -o queue-kanban . && ./queue-kanban summary --repo-root "{{justfile_directory()}}"

# Update the project-local do-work skill without an agent (reviews differences, then overwrites in place — git is the undo)
run-do-work-update:
    project_root="{{justfile_directory()}}"; skill_root="$project_root/.claude/skills/do-work"; [ -f "$skill_root/SKILL.md" ] || skill_root="$project_root"; bash "$skill_root/tools/do-work-update.sh" --project-root "$project_root"
```

Six deliberate choices in these recipes:

- **`$port` is an exported parameter, validated before anything else runs**: `just` interpolates `{{…}}` tokens textually into each recipe line's shell source, so a raw `{{port}}` would let `just run-kanban '8090; echo PWNED'` inject arbitrary commands (never interpolate raw user text into shell source — the rule applies to justfiles too). The `$` prefix hands the parameter to every recipe line as an environment variable instead — the shell reads `"$port"` as data, never code — and the first line rejects anything but digits before the kill-stale or build+serve lines can see it.
- **`go build` on every run** (`actions/board.md` Step 4's rule): `do-work update` overwrites the tool's source but leaves the previously compiled binary in place — a cached binary silently renders old logic. The incremental rebuild is near-instant when nothing changed, and the binary stays uncommittable via the tool's shipped `.gitignore`.
- **Each `cd … && …` chain stays on one logical line**: `just` runs every recipe line in a fresh shell, so a bare `cd` on its own line would not carry into the next — the same trap as prescribed action steps, where shell state never survives between command blocks.
- **The kill-stale check is its own recipe line and needs no `cd`**: it only touches `lsof`/`ps`/`kill` against `"$port"`, so it doesn't need the `<kanban-dir>` context the build+serve line does. `just` aborts a recipe on the first line that exits non-zero, so a squatting non-`queue-kanban` process's `exit 1` here stops the recipe *before* the build+serve line ever runs — no build is attempted and nothing gets killed. It kills only a listener whose full command line (verified via `ps -p PID -o args=` — `args`, unlike `comm`, is never truncated on Linux) contains `queue-kanban`, which covers the same binary built under another name or path by a different repo's recipes (e.g. `build/go-bin-queue-kanban`) — that's what lets the recipe reclaim the port from a board started in another folder; anything else is left running and named in the error. A missing `lsof` degrades gracefully — the check is skipped and the recipe proceeds straight to build+serve — rather than blocking the recipe on a tool that isn't guaranteed to exist.
- **`kanban-static` excludes its own output locally instead of dirtying the project**: the snapshot lands at `build/queue-kanban-board/`, which would otherwise sit in `git status` as untracked noise forever. The second recipe line appends a root-anchored pattern to `.git/info/exclude` — git's local-only ignore list — so no tracked file (like the project's `.gitignore`) is modified by a viewer command. The `check-ignore` test makes the append idempotent, it runs from `{{justfile_directory()}}` so the root-anchored pattern and the cwd-relative check can't mismatch (an interior-slash ignore pattern is root-anchored while `check-ignore` tests cwd-relative paths), and the exclude path comes from `git rev-parse --git-path` (never assembled from `--show-toplevel`), keeping it worktree-safe. In a non-git project the guard skips silently.
- **`run-do-work-update` resolves the installed skill relative to the justfile**: it prefers `.claude/skills/do-work/` for consumers and falls back to the project root for this repository's own recipe. The updater derives and verifies the skill root again, so a copied recipe cannot point an update outside the project. It is intentionally interactive at the overwrite gate: the direct command does not spend an agent turn, but it still shows the installed-versus-upstream diff and requires `y` before replacing skill files.

#### Phase 3: Verify

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
JUSTFILE_PATH=""
for justfile_candidate in justfile Justfile .justfile; do
  [ -f "$PROJECT_ROOT/$justfile_candidate" ] && { JUSTFILE_PATH="$PROJECT_ROOT/$justfile_candidate"; break; }
done
[ -n "$JUSTFILE_PATH" ] && grep -qE '^run-kanban[ :]' "$JUSTFILE_PATH" && grep -qE '^run-do-work-update[ :]' "$JUSTFILE_PATH" && echo "recipes: OK" || echo "recipes: FAILED"
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
  just run-do-work-update  Update the project-local skill (shows diff, asks before overwrite)

- `just run-kanban` replaces a stale queue-kanban instance already holding
  the port and opens your default browser at the board URL automatically —
  a non-queue-kanban process on the port is left alone and named in an error.
- Recipes rebuild the tool on every run, so `do-work update` refreshes take
  effect automatically (needs the Go toolchain — same requirement as
  `do-work board`).
- The justfile is project-owned: `do-work update` never touches it. When the
  shipped recipes change, re-run `do-work install just-kanban` — it diffs
  your installed block against the shipped one and offers an upgrade.
```

For the Phase 1b upgrade path, the first two lines read `Upgraded: just recipes for the queue-kanban board` / `Updated in: <project-root>/justfile` instead — the rest is identical.

---

## Workflow: `memory-module`

The `memory-module` target sets up the ADR-017 memory engine in the consuming project: the `memory/` store that `actions/memory.md` operates on, plus (on Claude Code) the two optional hooks that make it automatic — `hooks/memory-session-start.sh` injects the frozen snapshot at session start, `hooks/memory-stop-capture.sh` appends a deduplicated capture of each session's final exchange. Every `do-work memory` sub-command works without the hooks; they are the enhancement, not the dependency.

#### Phase 1: Check if already installed

Check the full guarantee set — scaffolding AND hook wiring (detecting on the scaffolding alone would let a half-completed prior run masquerade as installed):

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
# Scaffolding is "present" only when ALL components exist — a partial prior run
# (say, working-memory.md without logs/ or the ledger, or the store present but no
# longer locally ignored) must route back through Phase 2, which is idempotent and
# only fills gaps. Outside a git repo the ignore clause is vacuously satisfied.
ignore_ok=1
tracked_raw_store=""
if git -C "$PROJECT_ROOT" rev-parse --git-dir >/dev/null 2>&1; then
  git check-ignore -q "$PROJECT_ROOT/memory/logs/.probe" 2>/dev/null || ignore_ok=0
  git check-ignore -q "$PROJECT_ROOT/memory/usage-ledger.jsonl" 2>/dev/null || ignore_ok=0
  git check-ignore -q "$PROJECT_ROOT/memory/.bootstrap-imported" 2>/dev/null || ignore_ok=0
  # A tracked raw store is the ONE state that a fully-wired install can be in while
  # still leaking verbatim captures, and an ignore rule cannot repair it. It must be
  # detected HERE, not at Phase 4 verify: with the rule present, `ignore_ok` stays 1
  # and every other probe passes, so the "already installed" early return below would
  # exit before any later check ever runs — the exact repair case this is meant to catch.
  tracked_raw_store="$(git -C "$PROJECT_ROOT" ls-files -- memory/logs memory/usage-ledger.jsonl memory/.bootstrap-imported)"
fi
if [ -n "$tracked_raw_store" ]; then
  echo "raw store: TRACKED — repair required, not installable as-is:"
  printf '%s\n' "$tracked_raw_store" | sed 's/^/    /'
else
  echo "raw store: untracked"
fi
if test -s "$PROJECT_ROOT/memory/working-memory.md" \
   && test -d "$PROJECT_ROOT/memory/logs" \
   && test -f "$PROJECT_ROOT/memory/usage-ledger.jsonl" \
   && [ "$ignore_ok" -eq 1 ]; then
  echo "scaffolding: present"
else
  echo "scaffolding: absent or partial"
fi
SETTINGS_FILE="$PROJECT_ROOT/.claude/settings.json"
if [ -f "$SETTINGS_FILE" ] && grep -q 'memory-session-start.sh' "$SETTINGS_FILE" && grep -q 'memory-stop-capture.sh' "$SETTINGS_FILE"; then
  echo "hooks: wired"
else
  echo "hooks: absent or partial"
fi
```

- **`raw store: TRACKED` (checked FIRST, whatever the other two lines say)** → stop before Phase 2. Report the named paths, tell the user to `git rm --cached <path>` and commit the removal, and say plainly that until then the Stop hook keeps appending verbatim captures to a file that gets committed. Never untrack on their behalf — removing a path from the index is the user's call. This outranks "already installed" precisely because a fully-wired install is the state it hides in.
- **Both present** → report "already installed" and stop.
- **Scaffolding present, hooks absent/partial** → *repair mode*: skip Phase 2, run Phase 3 only. (Hooks wired but scaffolding absent/partial is also repair: run Phase 2 only.)
- **Both absent/partial** → run Phases 2–3 in full.

Phases 2 and 3 are both safe to re-run over partial state: Phase 2 only creates what's missing (never overwrites `working-memory.md`), and Phase 3 gates each hook entry independently — so "partial" never needs special-casing beyond routing into the phase.

#### Phase 2: Scaffold the memory store

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
mkdir -p "$PROJECT_ROOT/memory/logs"
touch "$PROJECT_ROOT/memory/usage-ledger.jsonl"
```

Then create `$PROJECT_ROOT/memory/working-memory.md` from the template in `actions/memory-reference.md` (with today's date in the `updated:` frontmatter) — **only if the file is absent or empty. Never overwrite an existing `working-memory.md`**: it is the user's standing memory, the same gate as the never-overwrite-an-existing-`SKILL.md` rule.

Finally, make the raw store machine-local. `memory/logs/` holds verbatim hook captures, `memory/usage-ledger.jsonl` is per-machine instrumentation, and `memory/.bootstrap-imported` records that *this machine's* session transcripts were imported — none should ever become committable. Redaction (`actions/memory-reference.md` → Capture Redaction Spec) is best-effort pattern matching, so version control must not be the place a missed credential lands. The sentinel is machine-local for a different reason: `do-work memory bootstrap` refuses to run when it exists, so a committed sentinel from one machine permanently blocks every other clone from importing its own local history — a silent failure that reads as "already imported". Only the curated `working-memory.md` stays committed and shareable. Add the three paths to the enclosing repo's `.git/info/exclude` (machine-local — never committed, never shipped); do **not** touch the project's committable `.gitignore`. This is the snippet from `crew-members/background-agents.md` → **Local-ignore snippet (for genuinely-transient paths)**, substituting these paths (see that file for why, including the linked-worktree caveat):

```bash
exclude=$(git rev-parse --git-path info/exclude 2>/dev/null) || exclude=""
if [ -n "$exclude" ]; then
  git check-ignore -q "$PROJECT_ROOT/memory/logs/.probe" 2>/dev/null \
    || echo '**/memory/logs/' >> "$exclude"
  git check-ignore -q "$PROJECT_ROOT/memory/usage-ledger.jsonl" 2>/dev/null \
    || echo '**/memory/usage-ledger.jsonl' >> "$exclude"
  git check-ignore -q "$PROJECT_ROOT/memory/.bootstrap-imported" 2>/dev/null \
    || echo '**/memory/.bootstrap-imported' >> "$exclude"
fi
```

Each path is gated **independently** — the same per-entry discipline Phase 3 uses for the hooks, so a partial prior state (one pattern present, one missing) repairs instead of skipping both. `git check-ignore -q` already succeeds when any ignore source covers the path, so the appends never duplicate on a re-run, and the guard is a clean no-op outside a git repo. (Ignored ≠ untracked: the Phase 1 fast-path gate and Phase 4's verification both run the `git ls-files` tracked-store check with the `git rm --cached` remedy — this block only makes the paths ignorable.)

#### Phase 3: Merge the hook entries

This phase only applies on Claude Code (a `.claude/` directory convention); on other platforms report `hooks: n/a (not Claude Code)` and continue to Phase 4 — the actions work hook-less.

Follow the **Hook Install Internals** in `actions/memory-reference.md` exactly: back up `.claude/settings.json` to `settings.json.pre-memory-module`, then jq-**append** the two entries from `hooks/memory-hooks.json` (resolve `<skill-root>` as the directory containing `SKILL.md`) into the `hooks.SessionStart` / `hooks.Stop` arrays. The dedup gate is **per entry** — one grep for `memory-session-start.sh`, one for `memory-stop-capture.sh`, each appending only its own missing entry (a partial/manual prior install may hold one hook but not the other). Compose, never clobber: append with `+`, never assign a whole new array — the consumer's existing hooks (including do-work's own `session-start.sh` / `pipeline-guard.sh`, if installed) must survive untouched.

If `jq` is unavailable: do NOT attempt a sed/awk merge. Print the two entries from `hooks/memory-hooks.json` with "merge these manually into `.claude/settings.json`" and record `hooks: MANUAL STEP` — a warning, not a failure.

#### Phase 4: Verify

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
test -d "$PROJECT_ROOT/memory/logs" && echo "logs dir: OK" || echo "logs dir: FAILED"
test -s "$PROJECT_ROOT/memory/working-memory.md" && echo "working memory: OK" || echo "working memory: FAILED"
test -f "$PROJECT_ROOT/memory/usage-ledger.jsonl" && echo "ledger: OK" || echo "ledger: FAILED"
if git -C "$PROJECT_ROOT" rev-parse --git-dir >/dev/null 2>&1; then
  # Tracked beats ignored. Ignore rules do not apply to files already in the index,
  # so on a repair install `check-ignore` prints OK while the Stop hook keeps
  # appending verbatim captures to a TRACKED log. Test tracking first and let it
  # override — an OK here must mean "cannot be committed", not "a pattern exists".
  tracked_raw_store="$(git -C "$PROJECT_ROOT" ls-files -- memory/logs memory/usage-ledger.jsonl memory/.bootstrap-imported)"
  if [ -n "$tracked_raw_store" ]; then
    echo "raw store untracked: FAILED — these are in the index, ignore rules cannot help:"
    printf '%s\n' "$tracked_raw_store" | sed 's/^/    /'
  else
    echo "raw store untracked: OK"
  fi
  git check-ignore -q "$PROJECT_ROOT/memory/logs/.probe" 2>/dev/null && echo "logs ignored: OK" || echo "logs ignored: FAILED"
  git check-ignore -q "$PROJECT_ROOT/memory/usage-ledger.jsonl" 2>/dev/null && echo "ledger ignored: OK" || echo "ledger ignored: FAILED"
  git check-ignore -q "$PROJECT_ROOT/memory/.bootstrap-imported" 2>/dev/null && echo "sentinel ignored: OK" || echo "sentinel ignored: FAILED"
fi
SETTINGS_FILE="$PROJECT_ROOT/.claude/settings.json"
if [ -f "$SETTINGS_FILE" ]; then
  if command -v jq >/dev/null 2>&1; then jq . "$SETTINGS_FILE" >/dev/null 2>&1 && echo "settings parses: OK" || echo "settings parses: FAILED"; else python3 -m json.tool "$SETTINGS_FILE" >/dev/null 2>&1 && echo "settings parses: OK" || echo "settings parses: UNCHECKED (no jq/python3)"; fi
  grep -q 'memory-session-start.sh' "$SETTINGS_FILE" && grep -q 'memory-stop-capture.sh' "$SETTINGS_FILE" && echo "memory hooks: OK" || echo "memory hooks: ABSENT (manual step pending?)"
  grep -c 'session-start.sh\|pipeline-guard.sh' "$SETTINGS_FILE" >/dev/null 2>&1 || true
fi
```

Three hard gates beyond the OK lines:

- **`settings parses: FAILED` after a merge = broken install.** Restore from `settings.json.pre-memory-module` immediately and report the failure. On success, remove the backup.
- **Pre-existing hooks must still be present.** Every hook entry that existed in the backup must still appear in the merged file (compare the backup's entry set against the merged file — a grep per pre-existing script filename is enough). Any entry lost → restore from backup and report.
- **`raw store untracked: FAILED` = the install is not safe to use.** The exclude entries added in Phase 2 cannot un-track a file, so a raw log or the sentinel that is already in the index stays committable no matter what `check-ignore` says. Do **not** report success and do **not** untrack anything yourself — removing a path from the index is the user's call. Report the named paths and the one-line remedy (`git rm --cached <path>`, then commit the removal), and say plainly that until they run it the Stop hook will keep appending verbatim captures to a tracked file. If the sentinel is the tracked path, add that every other clone will refuse to bootstrap its own history until it is removed.

#### Phase 5: Report back

```
Installed: memory-module (ADR-017 memory engine)

Destination: <project-root>/memory/
  working-memory.md   Standing memory, 2,500-char hard cap (curate via `do-work memory remember`)
  logs/               Dated daily logs (auto-appended by the Stop hook) — machine-local
  usage-ledger.jsonl  Usage instrumentation read by `do-work memory audit` — machine-local
  .bootstrap-imported Sentinel written by `do-work memory bootstrap` — machine-local

Only working-memory.md is committable. logs/, usage-ledger.jsonl, and
.bootstrap-imported were added to .git/info/exclude (machine-local; your .gitignore
is untouched) so verbatim captures never reach version control and each machine can
still bootstrap its own session history.

Hooks (Claude Code only, merged into .claude/settings.json):
  memory-session-start.sh  Injects the frozen snapshot at session start — writes surface NEXT session
  memory-stop-capture.sh   Captures each session's final exchange into today's log (hash-deduplicated)

- Every `do-work memory` sub-command works without the hooks — they are optional.
- Consider `do-work memory bootstrap` for a one-time import of past session history.
- This engine runs in parallel with bkb during the ADR-017 experiment;
  `do-work memory audit` renders the head-to-head.
```

If Phase 3 ended in `MANUAL STEP`, the Hooks section instead prints the two JSON entries and the manual-merge instruction.

---

## Help Block (no/unknown target)

When `$ARGUMENTS` is empty or doesn't match a known target, print:

```
install — install companion skills/tooling into the current project

  do-work install ui-design   Anthropic's frontend-design skill for production-grade UI
  do-work install bowser      Playwright CLI + Bowser skill for browser automation
  do-work install last30days  Engagement-ranked social-research engine (vendored, keyless)
  do-work install just-kanban  Justfile recipes for the board + updater (board needs Go)
  do-work install ideation-adhd  Parallel divergent-ideation skill (/adhd — cognitive-frame branching)
  do-work install memory-module  memory/ store + SessionStart/Stop hooks (ADR-017 memory engine)
```

Then stop.

## Output Format

- **`ui-design`**, **`ideation-adhd`**: a short status line — "already installed", "installed successfully", or an error describing what failed and how to finish manually.
- **`bowser`**: a two-line status — one for `playwright-cli`, one for the Bowser skill. Each is either "OK" (installed and verified), "already installed" (detected in Phase 1), or an error with the exact command the user can re-run.
- **`last30days`**: a per-guarantee status (skill file, ignore rule, Python 3.12+) — "already installed" only when every guarantee holds; otherwise "installed successfully" with the destination path, or the FAILED line(s) and the exact command the user can re-run.
- **`just-kanban`**: a per-component status (recipes appended or upgraded, justfile parses, `just`/`go` availability) — "already installed (current)" only when the installed recipes match the shipped block; a divergent block gets a diff and a consent-gated upgrade (Phase 1b), reported as "Upgraded" or "kept existing recipes". Missing toolchains are warnings, not failures.
- **`memory-module`**: a per-guarantee status (logs dir, working memory, ledger, local-ignore entries, settings parses, memory hooks) — "already installed" only when scaffolding AND hook wiring hold; hooks may report `MANUAL STEP` (no jq) or `n/a` (not Claude Code) as warnings, not failures.
- **Unknown / missing target**: the help block above.

## Rules

- **Install to the project, not globally.** Skill files go under `<project-root>/.claude/skills/<skill-name>/` (`<project-root>` resolved via `git rev-parse --show-toplevel || pwd`). Do not write to `~/.claude/` or any global path.
- **CLI is global; skill is project-scoped.** For `bowser`, `playwright-cli` goes to the global npm prefix; the Bowser skill goes under `<project-root>/.claude/skills/playwright-bowser/`. Don't mix them.
- **Never overwrite an existing non-empty skill `SKILL.md`.** Phase 1 of each workflow is the gate. If the file is present and non-empty, stop. (A zero-byte file is a failed download — reinstalling over it is repair, not overwrite, which is why the detect commands use `test -s` rather than `ls`.)
- **Only Chromium by default (bowser).** Other browsers bloat install time and aren't needed for ui-review's default flow.
- **Don't silently substitute a different skill or repo.** If the upstream URL fails, report the error — don't download a similarly-named skill from elsewhere.
- **Keyless in the project (last30days).** This install writes no config file at all. If a project-local `.claude/last30days.env` ever exists, it must never contain API keys — real keys live only in the user-global `~/.config/last30days/.env`. Never write a secret into any file inside the repo.
- **The vendor drop must be ignored (last30days).** Phase 2 adds `**/.claude/skills/last30days/` to the enclosing repo's `.git/info/exclude` when it isn't already covered — machine-local, never the project's committable `.gitignore` — because ~15 MB of upstream Python must never become committable in the consuming repo.
- **Touch only the four do-work recipes in the justfile (just-kanban).** Never reorder, reformat, or modify justfile content outside `run-kanban`/`kanban-static`/`kanban-summary`/`run-do-work-update`. A divergent installed recipe is replaced only after the user has seen the diff and accepted (Phase 1b) — never silently, and never on the assumption that different means stale. Create a `justfile` only when none of `justfile`/`Justfile`/`.justfile` exists at the project root.
- **Compose hook entries (memory-module).** Append to the `hooks.SessionStart`/`hooks.Stop` arrays in `.claude/settings.json`; never replace, reorder, or rewrite existing entries. Back up to `settings.json.pre-memory-module` before the merge; a post-merge parse failure or any lost pre-existing entry → restore the backup and report. Never overwrite an existing `memory/working-memory.md`.
- **Keep the raw memory store out of version control (memory-module).** `memory/logs/`, `memory/usage-ledger.jsonl`, and `memory/.bootstrap-imported` go in the enclosing repo's `.git/info/exclude` — machine-local, never the project's committable `.gitignore`. Only the curated `working-memory.md` is committable. The sentinel belongs there because `memory bootstrap` refuses to run when it exists: committed, one machine's import permanently blocks every other clone from importing its own history. **An ignore rule is not proof — verify with `git ls-files` too.** Ignore rules have no effect on tracked files, so on a repair install `check-ignore` reports OK over a raw log that is already in the index and still being appended to. Report tracked paths and hand the user the `git rm --cached` remedy; never untrack a file on their behalf.
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
| "Their existing `run-kanban` recipe differs from the shipped one — I'll just swap it out (just-kanban)" | Run Phase 1b: show the unified diff and get explicit consent before replacing anything | The divergence may be deliberate project-specific edits, not staleness — only the user can tell, and a silent swap destroys their edits |
| "The diff is obviously just an older shipped version — asking is pointless (just-kanban)" | Ask anyway; the consent prompt is one question and the diff makes it fast | "Obviously stale" is exactly how a deliberate one-flag customization gets flattened; there is no reliable way to distinguish old-shipped from hand-edited |
| "The skill is installed globally — I'll hard-code its absolute path into the recipe (just-kanban)" | Stop at the Phase 2 global-install gate and report | A recipe pointing outside the project breaks on every other clone and machine, and the skill's norms reject global installs |
| "settings.json already has hooks — I'll rewrite the whole hooks object cleanly (memory-module)" | Append via the jq recipe in `actions/memory-reference.md` | Assigning a fresh hooks object clobbers the user's other hooks — including do-work's own session-start and pipeline-guard |
| "No jq here — a quick sed insert will do (memory-module)" | Print the manual-merge snippet and report `hooks: MANUAL STEP` | Text-patching JSON corrupts settings.json on the first nested bracket it didn't expect |

## Red Flags

- The install command reported success but the verify step shows the file is empty — the URL may have changed; investigate before claiming success.
- `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists but has zero size and Phase 1 still reported "already installed" — the detect command regressed to a bare existence check (`ls`); it must be `test -s` so a re-run repairs the failed download.
- An install command writes `SKILL.md` directly instead of downloading to `SKILL.md.download` and renaming on success — a mid-transfer failure would leave a non-empty partial file the detect reads as installed; restore the temp-then-rename shape.
- A stray `SKILL.md.download` file sits next to an installed skill — a prior run's failure cleanup didn't run; delete it (the rename-on-success flow never leaves one behind).
- You installed into `~/.claude/skills/` instead of the project — undo and re-install to the correct path.
- `git rev-parse --show-toplevel` fails (not in a git repo) and you installed into `pwd` — acceptable, but warn the user the path may drift if they `cd` elsewhere.
- (bowser) `playwright-cli --help` succeeds but `playwright-cli install` fails silently — browsers aren't actually installed; headless runs will error later.
- (bowser) You installed `playwright-cli` into a project-local `node_modules` instead of globally — the CLI won't be on PATH for other sessions.
- (last30days) `git check-ignore -q .claude/skills/last30days/SKILL.md` exits non-zero in a git repo — the exclude entry was skipped or mismatched; fix it before anything gets committed. (Don't eyeball `git status` for this: a wholly-untracked `.claude/` collapses to a single `?? .claude/` row that hides the path either way.)
- (last30days) A project file (e.g. `.claude/last30days.env`) contains anything that looks like a credential — remove it and move the key to the user-global `~/.config/last30days/.env`.
- (last30days) Verify found no Python 3.12+ interpreter — the engine can't run; treat it as a failed install, not a soft warning.
- (just-kanban) The justfile diff shows anything beyond one appended block — existing recipes were reordered or rewritten; restore the file and re-append.
- (just-kanban) The appended recipe contains an absolute path (especially into `$HOME`) — the skill-root resolution went wrong; recipes must use project-relative paths.
- (memory-module) The merged `.claude/settings.json` has fewer hook entries than the `settings.json.pre-memory-module` backup — an existing hook was clobbered; restore the backup.
- (memory-module) `memory/working-memory.md` changed content during install — the never-overwrite gate was bypassed.
- (memory-module) `git check-ignore` fails for `memory/logs/`, `memory/usage-ledger.jsonl`, or `memory/.bootstrap-imported` after a reported success — the local-ignore step was skipped, and verbatim captures are one `git add -A` from being committed. (Don't eyeball `git status` for this: a wholly-untracked `memory/` collapses to a single `?? memory/` row, and `.bootstrap-imported` is a dotfile that never shows in it at all.)
- (memory-module) Anything under `memory/logs/`, or `memory/usage-ledger.jsonl`, or `memory/.bootstrap-imported` appears in `git ls-files` — the raw store got committed. Ignore rules don't apply to tracked files, so the Phase 2 exclude entries can't fix it and `check-ignore` will still report OK: a tracked log keeps collecting verbatim captures on every `git add -A`, and a tracked sentinel makes every other clone refuse to bootstrap its own session history. Tell the user to `git rm --cached <path>` and commit the removal.
- (memory-module) The project's `.gitignore` gained a `memory/` entry — the ignore belongs in `.git/info/exclude`; a committable ignore rule was written where a machine-local one was specified.
- (memory-module) A `settings.json.pre-memory-module` backup file left behind after a reported success — cleanup was skipped or the merge silently failed.

## Verification Checklist

- [ ] Step 1 correctly dispatched on `$ARGUMENTS` (known target → workflow; unknown/empty → help block).
- [ ] Phase 1 detected an existing install and stopped, OR Phase 2+ ran the fetch/install commands.
- [ ] After the verify phase, `<project-root>/.claude/skills/<skill-name>/SKILL.md` exists and is non-empty (skill-file targets: `ui-design`, `bowser`, `last30days`, `ideation-adhd`).
- [ ] (bowser only) `playwright-cli --help` runs without error and Chromium is installed.
- [ ] (last30days only) a Python 3.12+ interpreter is on PATH, `git check-ignore` covers `.claude/skills/last30days/`, and no project file gained an API key.
- [ ] (just-kanban only) the justfile gained exactly one appended block, `run-kanban` and `run-do-work-update` greps present, `just --list` parses when `just` is available, and no existing recipe was modified.
- [ ] (memory-module only) `memory/logs/`, a non-empty `working-memory.md`, and `usage-ledger.jsonl` exist; `git check-ignore` covers `memory/logs/`, `memory/usage-ledger.jsonl`, and `memory/.bootstrap-imported`, **and `git ls-files` returns nothing for those three paths** (ignored ≠ untracked), while the project's `.gitignore` is unmodified; `.claude/settings.json` parses and every pre-existing hook entry survived; the backup was removed on success.
- [ ] The report names the destination path so the user can verify location.
- [ ] No changes were made outside `<project-root>/.claude/skills/<skill-name>/` (plus, for `bowser`, the global npm install; for `last30days`, the machine-local `.git/info/exclude` entry; for `just-kanban`, the project justfile; for `memory-module`, `<project-root>/memory/`, `.claude/settings.json`, and the machine-local `.git/info/exclude` entries).
