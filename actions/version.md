# Version Action

> **Part of the do-work skill.** Handles version reporting and update checks.

**Current version**: 0.29.7

**Upstream**: https://raw.githubusercontent.com/knews2019/skill-do-work/main/actions/version.md

## Responding to Version Requests

When user asks "what version" or "version":
- Report the version shown above

## Responding to Update Checks

When user asks "check for updates", "update", or "is there a newer version":

1. **Fetch upstream**: Use your environment's web fetch capability to get the raw version.md from the upstream URL above
2. **Extract remote version**: Look for `**Current version**:` in the fetched content
3. **Compare versions**: Use semantic versioning comparison
4. **Report result** using the format below

### Report Format

**If update available** (remote > local):

1. **Tell the user**: `Update available: v{remote} (you have v{local}).`
2. **Check for local changes** in the skill's root directory (where SKILL.md lives):
   - If the directory is a git repo, run `git -C <skill-root> status --porcelain` and check for uncommitted changes to tracked files.
   - If it's **not** a git repo, check whether any files differ from a fresh install by looking for user-modified content (custom agent-rules, edited action files, etc.).
   - **If the working tree is dirty / has local modifications**: Stop and warn the user. List the modified files and ask for explicit confirmation before proceeding. Do NOT auto-update.
   - **If clean**: Proceed to step 3.
3. **Run the update** from the skill's root directory:
   ```
   curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'
   ```
   **Note:** tar extraction adds and overwrites files but does not delete files removed upstream. If the update changes significantly, stale files from older versions may remain. For a guaranteed clean update, delete the skill directory contents first (preserving `do-work/` queue data) and then extract.
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

## Responding to Changelog Requests

When user asks "changelog", "release notes", "what's new", "what's changed", "updates", or "history":

1. **Find the changelog**: Look for `CHANGELOG.md` in the skill's root directory (same level as `SKILL.md`)
2. **Read the file**: Load the full contents
3. **Reverse for terminal reading**: The changelog is written newest-on-top (conventional for file reading). For terminal output, reverse the version sections so the **most recent entries appear at the bottom** — right where the user's eyes are
   - Separate the header (everything before the first `## ` version heading) from the version entries
   - Split version entries at each `## ` heading (each heading + its body is one block)
   - Reverse the order of those blocks
   - Output: header first, then oldest-to-newest entries (so newest lands at the bottom)
4. **Print the result**: Output the reversed changelog directly — no file creation, just terminal output

### Why Reverse?

Changelogs are written newest-first so the file reads well. But in a terminal, the bottom of the output is where the user is looking. Reversing puts the latest changes at the bottom — no scrolling required.

### If No Changelog Exists

If `CHANGELOG.md` is not found in the skill root:

```
No changelog found for this skill.
```
