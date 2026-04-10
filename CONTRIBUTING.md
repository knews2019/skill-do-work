# Contributing to do-work

## Adding a new action

1. **Create the action file** at `actions/{name}.md` following the template in `CLAUDE.md` — description blockquote, steps, and guardrail sections.

2. **Add a routing entry** in `SKILL.md` — insert a row in the priority table and add keywords to the verb reference. Pick a priority that avoids conflicts with existing routes.

3. **Add next-steps suggestions** in `next-steps.md` — what should the user do after this action completes?

4. **Create a docs guide** at `docs/{name}-guide.md` — a human-readable summary (1-3 pages) covering what it does, when to use it, and key concepts.

5. **Update CLAUDE.md** — add the new file to the Project Structure tree.

6. **Bump version and changelog** — see "Version bumping" below.

## Modifying an existing action

- Preserve section order: Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist
- Don't remove existing sections without good reason
- Test the action with a real invocation before committing
- Bump the version for any user-visible change

## Required elements

Every action file must have:

- **Description blockquote** starting with `> **Part of the do-work skill.**`
- **Steps** (numbered `### Step N:` headings)

## Encouraged elements

- **When to Use** — positive triggers and explicit exclusions with redirects
- **Common Rationalizations** — 3-column table: `If you're thinking... | STOP. Instead... | Because...`
- **Red Flags** — observable failure symptoms
- **Verification Checklist** — checkbox exit criteria (`- [ ]`)

## Style guide

- Cross-reference other actions by short name ("the work action", "do work clarify") — not by file path
- Use generalized agent language ("spawn a subagent", "use your environment's ask-user prompt")
- Each action file must work as a standalone prompt in a basic chat interface
- Design for the floor: the simplest agent that can read/write files and run shell commands must be able to follow the instructions

## Version bumping

Before every commit:

1. Bump the version in `actions/version.md` (line starting with `**Current version**:`). Semver: patch for fixes, minor for features, major for breaking changes.
2. Verify the new version is strictly greater than the latest entry in `CHANGELOG.md`.
3. Add a changelog entry at the top of `CHANGELOG.md` with a unique two-word codename:

```markdown
## X.Y.Z — The [Fun Two-Word Name] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

## Crew members

Domain-specific rules live in `crew-members/{domain}.md`. Each file should have:

- A `JIT_CONTEXT` comment documenting when it loads
- A Role Identity section describing the perspective the agent should adopt
- Practical rules and patterns, not abstract advice

If a rules file is missing, agents proceed without it — never block on a missing file.
