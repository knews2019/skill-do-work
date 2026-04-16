# Version Action

> **Part of the do-work skill.** Handles version reporting, update checks, and work recaps.

**Current version**: 0.67.4

**Upstream**: https://raw.githubusercontent.com/knews2019/skill-do-work/main/actions/version.md

## Responding to Version Requests

When user asks "what version", "version", "what's new", "release notes", "what's changed", "updates", or "history":

1.  **Report Current Version**: Print the "Current version" from this file.
2.  **Report Changes**: Print the top entry (the most recent release) from the `CHANGELOG.md` in the skill root. Do not summarize or paraphrase; the user wants to see the canonical log entry.
3.  **Propose Update Check**: If the user's project has an internet connection, offer to run `do work version check` to see if a newer version is available.

## Update Check Protocol

When user asks "check for updates", "is there a new version", or "update check":

1.  **Fetch Upstream Version**: Download the `actions/version.md` from the "Upstream" URL above.
2.  **Compare**: Extract the version string from the fetched file and compare it to the "Current version" here.
3.  **Report Outcome**:
    -   If version matches: "You are on the latest version."
    -   If upstream is newer: "A new version (<Upstream Version>) is available. Would you like to update?"
4.  **Execute Update** (Only on explicit confirmation):
    -   Run the "Auto-Update" protocol below.

## Auto-Update Protocol

When user confirms "yes, update":

1.  **Safety Check**: Run `git status --porcelain` in the skill directory (`~/.claude/skills/do-work/`).
    -   If there are uncommitted changes in `actions/`, `prompts/`, `interviews/`, `specs/`, `docs/`, `decisions/`, `hooks/`, `CLAUDE.md`, `AGENTS.md`, or `next-steps.md`, **ABORT** and report that the skill has local modifications that would be overwritten.
2.  **Pre-Clean**: Remove all `.md` files directly under `prompts/` and `interviews/` to ensure that upstream renames/removals don't leave stale local entries.
3.  **Download and Extract**:
    -   Download the latest release tarball from GitHub: `https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz`.
    -   Extract the tarball over the existing skill directory, overwriting all files.
4.  **Completion**: Report "Update complete. You are now on version <New Version>." and invite the user to read the new release notes with `do work version`.
