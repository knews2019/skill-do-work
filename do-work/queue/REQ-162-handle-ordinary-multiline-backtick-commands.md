---
id: REQ-162
title: "Review fix: Handle ordinary multiline backtick commands"
status: pending-answers
domain: general
created_at: 2026-08-10T10:14:57Z
user_request: UR-031
addendum_to: REQ-159
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: just-multiline-literal-state
---

# Review Fix: Handle Ordinary Multiline Backtick Commands

## What

Extend reserved-recipe collision scanning to retain physical-line state for ordinary single-backtick Just command literals. Done means the broader multiline-literal scanner accepts reserved-looking command payload without weakening real definition detection, exact diagnostics, or pre-mutation preservation.

## Context

Discovered during REQ-159 and independently confirmed in its review. Just 1.46.0 accepts physical newlines inside ordinary single-backtick commands, but REQ-159 deliberately kept them line-local because its explicit Requirements named three other literal families and described ordinary backticks as one-line.

## Requirements

- Retain ordinary single-backtick command state across physical lines using Just's actual closing behavior.
- Keep same-line backtick commands, comments, indented recipe bodies, triple-backtick commands, and every existing string family from hiding real definitions.
- Add Just-parseable positive fixtures plus exact byte-preserving real-collision controls and keep paired helpers/full contracts green.

## Instances

- [ ] `tools/replace-text-section.sh`: an ordinary backtick command spanning physical lines exposes column-zero reserved-looking payload to `just_definition_names()`.
- [ ] `skills/do-work/tools/replace-text-section.sh`: apply the correction byte-identically to the shipped helper.
- [ ] `_dev/tests/contract-regressions.sh`: add positive and exact negative production-helper fixtures for the accepted multiline form.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-159: Just accepts ordinary single-backtick commands across physical lines, and the collision scanner still rejects safe command payload that resembles a reserved recipe. REQ-159 explicitly treated ordinary backticks as one-line, so the cascade-depth rule requires your consent before extending that boundary. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
