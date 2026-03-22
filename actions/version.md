# Version Action

> **Part of the do-work skill.** Handles version reporting and update checks.

**Current version**: 0.30.4

**Upstream**: https://raw.githubusercontent.com/knews2019/skill-do-work/main/actions/version.md

## Responding to Version Requests

When user asks "what version" or "version":

1. Report the version shown above
2. **Show last 5 skill releases**:
   - Read the first ~80 lines of `CHANGELOG.md` (do NOT load the full file)
   - Extract the 5 most recent version entries (split at `## ` headings, take first 5 blocks)
   - Reverse so newest is at the bottom (right where the user's eyes are)
   - Print them after the version number

## Responding to Update Checks

When user asks "check for updates", "update", or "is there a newer version":

1. **Fetch upstream**: Use your environment's web fetch capability to get the raw version.md from the upstream URL above
2. **Extract remote version**: Look for `**Current version**:` in the fetched content
3. **Compare versions**: Use semantic versioning comparison
4. **Report result** using the format below

### Report Format

**If update available** (remote > local):

1. **Tell the user**: `Update available: v{remote} (you have v{local}).`
2. **Check for local changes** to shipped skill files (where SKILL.md lives):
   - **Scope the check to skill-owned files only.** Ignore `do-work/` (queue data, archives, deliverables) — those are generated at runtime and should never block an update.
   - If the directory is a git repo, run `git -C <skill-root> status --porcelain -- SKILL.md actions/ agent-rules/ CHANGELOG.md README.md` (listing only shipped paths) and check for uncommitted changes.
   - If it's **not** a git repo, check whether shipped skill files (actions/, agent-rules/, SKILL.md, etc.) differ from a fresh install by looking for user-modified content (custom agent-rules, edited action files, etc.).
   - **If any shipped skill files are dirty / have local modifications**: Stop and warn the user. List the modified files and ask for explicit confirmation before proceeding. Do NOT auto-update.
   - **If clean**: Proceed to step 3.
3. **Run the update** from the skill's root directory:
   ```
   curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'
   ```
   **Note:** tar extraction adds and overwrites files but does not delete files removed upstream. Stale files from older versions may remain. This is generally harmless — the skill only loads files it references. If you need a fully clean update, delete only the known skill paths (`actions/`, `agent-rules/`, `SKILL.md`, `CHANGELOG.md`, `README.md`) before extracting — never delete `do-work/` or other project files.
4. **Verify**: Read `actions/version.md` again and confirm the local version now matches the remote version.
5. **Report result**: `Updated to v{remote}.`

Do NOT just print the curl command and ask the user to run it. You are the agent — run it yourself.

**If up to date** (local >= remote):

```
You're up to date (v{local})
```

**If fetch fails**:

```
Couldn't check for updates.
```

Attempt the update anyway using the curl command above (still respecting the dirty-tree check in step 2). If that also fails, report the error and provide the manual command as a fallback:

```
To manually update, run this from the skill's root directory (where SKILL.md lives):
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'

Or visit: https://github.com/knews2019/skill-do-work
```

## Responding to Recap Requests

When user asks "recap":

1. **Find the archive**: Look for `do-work/archive/` in the project root
2. **Find the 5 highest-numbered UR folders** (e.g., `UR-012`, `UR-011`, etc.)
3. **For each UR**:
   - Read `input.md` for the title
   - List REQ files and extract titles from frontmatter or filename slug
4. **Format as a "Recent Work" section**:
   ```
   ## Recent Work

   UR-012 — [title from input.md]
     REQ-045 — [title]
     REQ-046 — [title]
   UR-011 — [title from input.md]
     REQ-043 — [title]
   ```
   One line per UR, one indented line per REQ. No descriptions, no scores, no file lists.
5. **If no archive exists** (`do-work/archive/` not found or empty): Print `No completed work yet.` and skip this section.
