# Install UI Design Action

> **Part of the do-work skill.** Installs the `frontend-design` Claude skill into the current project — specialized design knowledge that improves visual quality, layout, and interaction patterns for UI work.

Once installed, the skill is available to all agents in the project, including do-work builders processing `domain: ui-design` requests.

## When to Use

**Use when:**
- The project is about to start UI work and no frontend-design skill is present under `.claude/skills/`.
- A `domain: ui-design` REQ is about to be built and the builder needs design knowledge beyond what's in `crew-members/ui-design.md`.
- A ui-review pass flagged the absence of the skill as a cause of low-quality output.

**Do NOT use when:**
- The skill is already installed (Step 1 detects this and exits).
- The user is asking for browser automation or visual verification — that's `install-bowser`.
- The project explicitly uses a different design system/skill and adding frontend-design would conflict.

## Input

No arguments. The action is idempotent — running it when the skill is already present is a no-op.

## Steps

### Step 1: Check if already installed

```bash
ls .claude/skills/frontend-design/SKILL.md 2>/dev/null
```

If the file exists, report "already installed" and stop.

### Step 2: Install the skill

```bash
mkdir -p .claude/skills/frontend-design
curl -fsSL -o .claude/skills/frontend-design/SKILL.md \
  https://raw.githubusercontent.com/anthropics/claude-code/main/skills/frontend-design/SKILL.md
```

If `curl` is unavailable or the download fails, check the environment's plugin/skill registry (e.g., `/plugin install frontend-design`) and install from there.

### Step 3: Verify installation

```bash
test -s .claude/skills/frontend-design/SKILL.md && echo "Installed successfully" || echo "Installation failed"
```

### Step 4: Report back

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

## Output Format

A short status line: either "already installed", "installed successfully", or an error describing what failed and how to finish manually.

## Rules

- **Install to the project, not globally.** The skill file goes under `.claude/skills/frontend-design/` in the current project. Do not write to `~/.claude/` or any global path.
- **Never overwrite an existing `SKILL.md`.** Step 1 is the gate. If the file is present, stop.
- **Don't silently substitute a different skill.** If the upstream URL fails, report the error — don't download a similarly-named skill from elsewhere.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The file is already there but looks stale — I'll overwrite it" | Report installed; let the user decide whether to re-download | Overwriting a user-customized skill file silently destroys their edits |
| "`curl` failed, I'll fetch it with `wget` from a mirror" | Report the failure and the URL; let the user install manually | Unknown mirrors risk installing a tampered skill file |
| "The user wanted UI help, I should install bowser too while I'm here" | Stop after the frontend-design install; mention `install-bowser` as a next step if relevant | Each install action has a single, documented scope |

## Red Flags

- `curl` reported success but `test -s` says the file is empty — the URL may have changed; investigate before claiming success.
- `.claude/skills/frontend-design/SKILL.md` exists but has zero size — treat as not-installed and re-download (with user confirmation).
- You installed into `~/.claude/skills/` instead of the project — undo and re-install to the correct path.

## Verification Checklist

- [ ] Step 1 detected existing install and stopped, OR Step 2 ran fetch commands.
- [ ] `.claude/skills/frontend-design/SKILL.md` exists and is non-empty after Step 3.
- [ ] The report names the destination path so the user can verify location.
- [ ] No changes were made outside `.claude/skills/frontend-design/`.
