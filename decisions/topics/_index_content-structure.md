# Topic: Content Structure

How the skill's **own** source files are organized so they stay readable, loadable within model context limits, and evolvable without breaking everything downstream.

## The Structure Principles

- **Domain knowledge lives in crew-members, not actions.** Actions describe *what to do*; crew-member files describe *how a specialist would do it*. Crew files load just-in-time based on REQ frontmatter so irrelevant rules don't flood the context.
- **No action file should exceed the single-read token budget.** When an action's content passes the ~10k-token limit that agents enforce on a single read, it splits into a primary file (flow + rules) and a companion reference file (templates + schemas). The two are cross-linked, not duplicated.

## ADRs in This Cluster

| ADR  | Title                                 | Status   |
|------|---------------------------------------|----------|
| [[../adr/0007-crew-member-jit-loading\|0007]]  | Crew members — JIT-loaded domain rules       | accepted |
| [[../adr/0009-companion-reference-files\|0009]] | Companion reference file pattern             | accepted |

## How They Relate

- [[../adr/0007-crew-member-jit-loading|0007]] is about *horizontal* structure — the right rules load based on the kind of REQ being processed.
- [[../adr/0009-companion-reference-files|0009]] is about *vertical* structure — a long action gets split in two so it stays loadable.

Both serve the same end: agents walk into each action with the smallest useful context.

## Template Enforcement

The consistent structure of every action file (Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist) is documented in `CLAUDE.md` under "Action File Conventions." That template isn't an ADR on its own — it's a convention enforced by reviewers and the `do work code-review` action. Treat it as the surface that these structural ADRs shape.
