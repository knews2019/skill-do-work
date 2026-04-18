# Version Action

> **Part of the do-work skill.** Handles version reporting, update checks, and work recaps.

**Current version**: 0.69.5

**Upstream**: https://raw.githubusercontent.com/knews2019/skill-do-work/main/actions/version.md

## Responding to Version Requests

When user asks "what version", "version", "what's new", "release notes", "what's changed", "updates", or "history":

1. Report the version shown above
2. **Show last 5 skill releases**:
   - Read the first ~80 lines of `CHANGELOG.md` in the skill's root directory (same level as `SKILL.md`) — do NOT load the full file
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
2. **Preflight: confirm this is a project-local install, not a global one.** The update must overwrite the copy inside the current project — never a user-wide / global install.
   - Resolve the absolute path of the skill's root directory (where `SKILL.md` lives). Call this `<skill-root>`.
   - **Refuse to auto-update if `<skill-root>` is under any user-wide skills location**, including but not limited to:
     - `~/.claude/skills/...`
     - `~/.gemini/skills/...`
     - `~/.cursor/skills/...`
     - `~/.config/*/skills/...`
     - anything else under `$HOME` that isn't also inside the current project's git repo.
   - Resolve the current project's git root: `git -C <invocation-dir> rev-parse --show-toplevel`. If `<skill-root>` is **not** a descendant of that project root, stop and report:
     ```
     Skill is installed at <skill-root>, which is outside the current project (<project-root>).
     Refusing to update a global/shared install from here. Either:
       - cd into the project that owns <skill-root> and re-run, or
       - install the skill locally inside this project (e.g. <project-root>/.claude/skills/do-work/) and re-run.
     ```
     Do NOT proceed. Do NOT suggest the curl command.
   - Only continue once `<skill-root>` is confirmed to live inside the current project's git root.
3. **Check for local changes** to shipped skill files at `<skill-root>`:
   - **Scope the check to skill-owned files only.** Ignore `do-work/` (queue data, archives, deliverables) — those are generated at runtime and should never block an update.
   - If `<skill-root>` is a git repo, run `git -C <skill-root> status --porcelain -- SKILL.md actions/ crew-members/ prompts/ interviews/ specs/ docs/ decisions/ hooks/ CLAUDE.md AGENTS.md CHANGELOG.md README.md next-steps.md` (listing every shipped editable path) and check for uncommitted changes. Any dirty file in these paths will be clobbered by the tar extraction in step 5 if you proceed.
   - If it's **not** a git repo, check whether shipped skill files (actions/, crew-members/, prompts/, interviews/, specs/, docs/, decisions/, hooks/, SKILL.md, CLAUDE.md, next-steps.md, etc.) differ from a fresh install by looking for user-modified content (custom crew-members, edited action files, local prompt/template additions, ADR edits, etc.).
   - **If any shipped skill files are dirty / have local modifications**: Stop and warn the user. List the modified files and ask for explicit confirmation before proceeding. Do NOT auto-update.
   - **If clean**: Proceed to step 4 (pre-clean) then step 5 (extract).
4. **Pre-clean discoverable directories.** `prompts/` and `interviews/` are enumerated by `do work prompts list` and `do work interview list` (they glob `prompts/*.md` and `interviews/*.md`), so any upstream-removed file that stays on disk will still appear as a live workflow. The dirty check in step 3 has already confirmed these are clean, so removing the tracked `.md` files here is safe and the tar extraction will restore them fresh:
   ```
   find <skill-root>/prompts -maxdepth 1 -name '*.md' ! -name 'README.md' -delete
   find <skill-root>/interviews -maxdepth 1 -name '*.md' -delete
   ```
   Do NOT delete files in `prompts/` or `interviews/` subdirectories — only the top-level `.md` files are globbed. Do NOT touch `do-work/` or any other runtime directory.
5. **Run the update in place at `<skill-root>`** (the project-local path confirmed in step 2). `cd` there first so the extraction cannot land in a global directory by mistake:
   ```
   cd <skill-root> && curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'
   ```
   **Note:** tar extraction adds and overwrites files but does not delete files removed upstream. For non-discoverable directories (`actions/`, `crew-members/`, `specs/`, `docs/`, `decisions/`) leftovers are harmless — the skill only loads files it references by name. For `prompts/` and `interviews/`, the pre-clean step above is what prevents ghost entries; if you skipped it, run `do work prompts list` and `do work interview list` after updating and delete anything that looks obsolete. Never delete `do-work/` (runtime state).
6. **Verify**: Read `<skill-root>/actions/version.md` again and confirm the local version now matches the remote version.
7. **Report result**: `Updated to v{remote} at <skill-root>.`

Do NOT just print the curl command and ask the user to run it. You are the agent — run it yourself.

**If up to date** (local >= remote):

```
You're up to date (v{local})
```

**If fetch fails**:

```
Couldn't check for updates.
```

Attempt the update anyway using the curl command above (still respecting the preflight location check in step 2 and the dirty-tree check in step 3 — refuse if the install is global). If that also fails, report the error and provide the manual command as a fallback:

```
To manually update, cd into the **project-local** skill root (where SKILL.md lives inside *this* project — NOT ~/.claude/skills/, ~/.gemini/skills/, or any other global skills directory) and run:

cd <project-root>/path/to/skill-do-work
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'

Or visit: https://github.com/knews2019/skill-do-work
```

## Responding to Recap Requests

When user asks "recap":

1. **Archive source** (`do-work/archive/UR-*/`): Read as before — title from `input.md`, REQs from `REQ-*.md` files inside each UR folder.
2. **Active source** (`do-work/user-requests/UR-*/`): Read `input.md` for the title. For REQs, scan `do-work/queue/REQ-*.md` files whose `user_request:` frontmatter field matches the UR id (e.g., `user_request: UR-143`). Also check `do-work/working/` for claimed REQs belonging to the UR.
3. **Merge**: Combine both lists, deduplicate by UR id (archive version wins if both exist), sort by UR number descending, take top 5.
4. **Label each UR**:
   - No label if fully archived
   - `(pending)` if the UR has any pending REQs
   - `(completed, awaiting archive)` if all its REQs are completed/done but the UR isn't archived yet
5. **Format as a "Recent Work" section**:
   ```
   ## Recent Work

   UR-144 — Block-level improved translation for ZH pairs
     REQ-361 — Block-level improved translation
   UR-143 — Model selector thinking variants (completed, awaiting archive)
     REQ-360 — Model selector thinking variants
   UR-142 — Quality-Score-Driven Repair Loop (completed, awaiting archive)
     REQ-359 — Quality-Score-Driven Repair Loop
   UR-011 — Dark mode implementation
     REQ-043 — Theme store setup
     REQ-044 — Settings panel toggle
   ```
   One line per UR, one indented line per REQ. No descriptions, no scores, no file lists.
6. **If no archive exists AND no active URs found**: Print `No completed work yet.` and skip this section.
