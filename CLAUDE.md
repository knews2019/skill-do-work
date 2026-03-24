# Do-Work Skill Project

A task queue skill for agentic coding tools. Platform-agnostic — works with any agent that can read/write files and run shell commands.

## Project Structure

```
SKILL.md              # Entry point — routing logic, action dispatch
README.md             # Installation + quick usage guide
actions/              # Action files (each is a standalone prompt)
  capture.md          # Capture new requests → UR folders + REQ files
  work.md             # Process the queue — triage, plan, build, test, review
  verify-requests.md  # Quality-check captured REQs against original input
  review-work.md      # Post-work code review + acceptance testing
  present-work.md     # Client-facing deliverables (briefs, videos, diagrams)
  cleanup.md          # Archive consolidation
  commit.md           # Atomic git commits traced to REQs
  version.md          # Version reporting + update checks (current version lives here)
  quick-wins.md       # Scan for refactoring opportunities and low-hanging tests
  install-ui-design.md # Install the frontend-design skill for UI work
agent-rules/          # Domain-specific rules loaded by work action
CHANGELOG.md          # Release notes (newest on top)
```

## Before Every Commit

1. **Bump the version** in `actions/version.md` (line starting with `**Current version**:`). Use semver — patch for fixes, minor for features, major for breaking changes. When in doubt, patch.

2. **Add a changelog entry** at the top of `CHANGELOG.md` (below the header):

```markdown
## X.Y.Z — The [Fun Two-Word Name] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

Keep it brief, newest on top, lead with value not implementation. Every version gets an entry.

## Agent Compatibility

Action files must work with **any** agentic coding tool:

- Use generalized language ("spawn a subagent", "use your environment's ask-user prompt") — no tool-specific APIs in action files.
- Each action file should work as a standalone prompt pasted into a basic chat interface.
- Design for the floor: the simplest agent that can read/write files and run shell commands must be able to follow the instructions. Subagents and parallel execution are nice-to-haves.
