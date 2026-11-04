---
id: UR-163
title: 'Add optional AI-Evaluation-Loop modes'
created_at: 2026-11-04T15:54:52Z
requests: [REQ-710]
word_count: 25
---
## Summary
Capture the optional AI-Evaluation-Loop proposal discussed in the tstack investigation. Preserve ordinary do-work runs.

## Extracted Requests
| Request | Intent |
| --- | --- |
| REQ-710 (Add optional AI-Evaluation-Loop modes) | Opt-in AI-Evaluation-Loop modes with a maximum number of rounds, traceable findings and truthful outcomes. |

## Batch Constraints
Capture only. Do not implement or install upstream skills during this invocation.

## Full Verbatim Input
> ```
> Add optional bounded review modes to do-work, preserving the default loop and using a finding ledger, verified model identities, current-diff approval, and explicit unresolved outcomes.
> ```

## Clarification — Optional Advisory Reviewer (2026-11-04)

The jester is optional. When enabled, its default model is the latest available Gemini; the named version is an example rather than a fixed model ID. This extends REQ-710 (Add optional AI-Evaluation-Loop modes).

User clarification:
> ```
> the jester is an optional advisory reviewer (latest gemini, like gemini 4 argon) is the default on it
> ```
