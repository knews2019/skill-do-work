# Install Bowser Action

> **Part of the do-work skill.** Installs Playwright CLI (`playwright-cli`) plus the Bowser skill, giving agents headed-browser automation, screenshots, and visual UI verification.

Once installed, `do-work ui-review` can use Playwright CLI for viewport screenshots, accessibility checks on rendered pages, and layout issues static analysis misses.

## When to Use

**Use when:**
- The project needs browser automation (screenshots, DOM snapshots, interaction tests) and neither `playwright-cli` nor the Bowser skill is present.
- A ui-review pass flagged "visual verification skipped" as a reason for incomplete findings.
- The user asked for headed-browser workflows, parallel named sessions, or persistent profiles.

**Do NOT use when:**
- `playwright-cli --help` already succeeds AND `.claude/skills/playwright-bowser/SKILL.md` exists (Step 1 detects this and exits).
- The user is asking for design-quality help — that's `install-ui-design`.
- The environment can't install global npm packages and the user hasn't consented to a local-only install.

## Input

No arguments. The action is idempotent — re-running is a no-op when both components are present.

## Steps

### Step 1: Check if already installed

```bash
playwright-cli --help >/dev/null 2>&1 && echo "playwright-cli: installed" || echo "playwright-cli: not found"
ls .claude/skills/playwright-bowser/SKILL.md 2>/dev/null && echo "bowser skill: installed" || echo "bowser skill: not found"
```

If both are present, report installed and stop.

### Step 2: Install Playwright CLI

```bash
npm install -g @anthropic-ai/playwright-cli@latest
```

If `npm` isn't available:

```bash
yarn global add @anthropic-ai/playwright-cli@latest
```

If neither package manager works, report the error and the install command so the user can run it manually.

### Step 3: Install Playwright browsers

```bash
playwright-cli install --with-deps chromium
```

Only Chromium is installed — sufficient for UI review. For Firefox/WebKit, the user can add `playwright-cli install firefox` or `webkit` later.

If `--with-deps` fails due to permissions:

```bash
npx playwright install chromium
```

### Step 4: Install the Bowser skill

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

### Step 5: Verify installation

```bash
playwright-cli --help >/dev/null 2>&1 && echo "playwright-cli: OK" || echo "playwright-cli: FAILED"
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
test -s "$PROJECT_ROOT/.claude/skills/playwright-bowser/SKILL.md" && echo "bowser skill: OK" || echo "bowser skill: FAILED"
```

### Step 6: Report back

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

## Output Format

A two-line status: one for `playwright-cli`, one for the Bowser skill. Each is either "OK" (installed and verified), "already installed" (detected in Step 1), or an error with the exact command the user can re-run.

## Rules

- **CLI is global; skill is project-scoped.** `playwright-cli` goes to the global npm prefix; the Bowser skill goes under `<project-root>/.claude/skills/playwright-bowser/`. Don't mix them.
- **Only Chromium by default.** Other browsers bloat install time and aren't needed for ui-review's default flow.
- **Never overwrite an existing Bowser skill file.** Step 1 is the gate.
- **Don't silently fall back to a different Bowser repo.** If the two documented URLs fail, stop and tell the user.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll install all three browsers to be safe" | Install only Chromium; mention the manual command for the others | Firefox + WebKit roughly triple the install time and disk use, for a feature ui-review doesn't need |
| "The Bowser skill file is old — I'll overwrite it with the latest" | Report installed; let the user decide whether to refresh | User may have customized the skill for their workflow |
| "npm install failed, I'll try yarn and pnpm and bun until something works" | Try npm, then yarn; if both fail, stop and report | Quiet package-manager shopping leaves the user unsure what got installed |
| "I'll also install the UI design skill since they're both UI-related" | Stop after bowser; mention `install-ui-design` as a next step if relevant | Each install action has a single documented scope |

## Red Flags

- `playwright-cli --help` succeeds but `playwright-cli install` fails silently — browsers aren't actually installed; headless runs will error later.
- The Bowser `SKILL.md` file exists but is empty — a prior download was interrupted; re-fetch.
- `git rev-parse --show-toplevel` fails (not in a git repo) and you installed the skill into `pwd` — acceptable, but warn the user the path may drift if they `cd` elsewhere.
- You installed `playwright-cli` into a project-local `node_modules` instead of globally — the CLI won't be on PATH for other sessions.

## Verification Checklist

- [ ] `playwright-cli --help` runs without error after Step 2.
- [ ] Chromium is installed (Step 3 reported success, or `playwright-cli install chromium` returns exit 0).
- [ ] `<project-root>/.claude/skills/playwright-bowser/SKILL.md` exists and is non-empty.
- [ ] Step 5 printed "OK" for both components.
- [ ] The report includes a usage example the user can paste directly.
