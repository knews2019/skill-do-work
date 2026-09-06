# Prime: Prescribed Shell Commands

> Read this before writing or reviewing any shell block prescribed inside an action file,
> hook, or tool. Prime files are low noise, high value: every trap below was earned in a
> real debugging session, and new lessons in this domain get added here so future sessions
> find them. This file is maintainer-side (`_dev/` is export-ignored) — nothing shipped
> may cite it.

## Prescribed Shell Commands Must Surface What the Steps Consume

See [`lessons-shell-commands.md#prescribed-shell-commands-must-surface-what-the-steps-consume`](lessons-shell-commands.md#prescribed-shell-commands-must-surface-what-the-steps-consume).

## Unchecked Exit Status Reads as Content

See [`lessons-shell-commands.md#unchecked-exit-status-reads-as-content`](lessons-shell-commands.md#unchecked-exit-status-reads-as-content).

## A Writer's SIGPIPE Death Reads as the Reader's Verdict

See [`lessons-shell-commands.md#a-writers-sigpipe-death-reads-as-the-readers-verdict`](lessons-shell-commands.md#a-writers-sigpipe-death-reads-as-the-readers-verdict).

## Closed Enumerations Go Stale

When a rule applies "whenever X happens" (load a guardrail, honor an enum, keep a guide in sync), state the trigger _condition_ in the rule's canonical home and mark any caller/value list as illustrative, not exhaustive. Hand-enumerated lists silently go stale the moment the set grows — one review traced four independent defects to this pattern (capture's stale domain enum, prompt-injection's five-caller list, the docs-exemption list, security.md's loader claims). When extending a set, grep for every other enumeration of it and update or generalize each one.

## Every Flag on a Shipped Script Needs a Non-Test Caller

See [`lessons-shell-commands.md#every-flag-on-a-shipped-script-needs-a-non-test-caller`](lessons-shell-commands.md#every-flag-on-a-shipped-script-needs-a-non-test-caller).

## Traps

- [family: assertion-passes-when-feature-is-dead] A green assertion can pass when its feature never runs → prove both execution/reuse and invalidation, including the manifest inputs that select the path.

## Stakes

- Prescribed command results and mutation boundaries
  Req: distinguish failed observation from empty results and bind writes to current, owned targets.
  Value: workflow decisions consume evidence the command actually established.
  Risk: swallowed failures can authorize destructive edits; inspect the relocated shell traps before changing a prescribed primitive.

## Lessons

See [`lessons-shell-commands.md`](lessons-shell-commands.md) — read it before changing what **Read first** or **Traps** name above.
