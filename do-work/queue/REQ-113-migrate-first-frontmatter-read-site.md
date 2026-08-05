---
id: REQ-113
title: "Confirm: migrate the first prose read site onto queue-kanban frontmatter"
status: pending-answers
created_at: 2026-08-05T19:24:00Z
user_request: UR-021
addendum_to: REQ-112
domain: general
prime_files: []
tdd: false
depends_on: [REQ-112]
maintenance: false
builder_decided: true
---

# Confirm: Migrate the First Prose Read Site

## What

While building REQ-112, the builder had to decide whether to also change one action file to actually *use* the new command, or to ship the command with tests only. It chose tests only. This follow-up confirms whether that matches your intent.

The background, since this will read cold later: the repo now has a `queue-kanban frontmatter get <file> <field>` command that reads a field out of a REQ file. Roughly 95 places across `actions/` still read REQ fields by hand instead — usually a few lines of `awk` or `grep` written inline in the prose of an action. The new command exists so those places *can* stop hand-rolling it, but nothing has been switched over yet, so today the command has no user.

**Why this is your call rather than the builder's:** switching one of those places over means editing an action file, and every action file has to keep working for an agent that has no compiler available. So the real question isn't "does the command work" — it does — but "what should an action say to do when the compiled tool isn't there?" That's a judgment about how much the skill leans on optional tooling, which is a standing preference of yours, not a technical detail.

## What the Builder Chose

Ship the command with tests only. REQ-112's diff stayed entirely inside `tools/queue-kanban/`, so it could be reviewed as a self-contained tooling change.

## What Would Change

If you'd rather see a site migrated: one low-risk action would gain a line naming the command as the preferred way to read a field, plus a documented fallback for when the binary isn't built. The candidate the builder had in mind is `actions/commit.md` Step 3, which checks whether an archived REQ has a terminal-success status — a check the new `--in-set terminal-success` flag does exactly, and one of the five places where a documented Red Flag warns that hand-rolling it drops `completed-with-issues`.

## Open Questions

- [ ] Should one action file now be switched over to use the new command, as a worked example?
  Recommended: No — leave the command unused for now. It's available when any action wants it, and each switch-over can be reviewed on its own.
  Value: keeps tooling changes and prose changes in separate reviews, and avoids committing the whole skill to a pattern before one real use has been seen.
  Risk: a command with no users can turn out to fit no real call site, and the hand-rolled reads the census flagged all stay as they are until someone starts. Cheap to reverse — adopting it later is a one-line prose edit per action, with nothing to undo.
  Also: Yes, migrate `actions/commit.md` Step 3 as the worked example; or yes, but pick a different site.

## Full Context

See `do-work/user-requests/UR-021/input.md`. The originating decision is D-01 in `REQ-112`'s `## Decisions`.
