---
id: REQ-202
title: Complete unsafe Remotion preview mutation detection
status: pending
domain: testing
created_at: 2026-08-15T18:45:11Z
user_request: UR-042
addendum_to: REQ-192
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: unsafe-remotion-preview-mutation-detection
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
---

# Review Fix: Complete Unsafe Remotion Preview Mutation Detection

## What

Make the completed-work presentation regression detector recognize every executable fixed-port and platform-opener form while preserving safe foreground preview commands and non-executable prohibition prose.

This is one sweep because the missed forms share a single root cause: the unsafe-command extractor matches a narrow set of literal spellings instead of the complete prohibited executable command families.

## Context

Found during review of REQ-192. `_dev/tests/contract-regressions.sh` catches background/sleep, literal `open http://localhost:3000`, and render workflows, but mutation probes show that `remotion studio src/Root.tsx --port 3000` and `open "$REMOTION_PREVIEW_URL"` can pass.

## Instances

- [ ] Fixed-port Remotion Studio flags, including separated and equals forms such as `--port 3000` and `--port=3000`.
- [ ] Command-start platform opener forms whose target is a variable or other nonliteral expression, including `open "$url"`.

## Requirements

- Reject executable fixed-port Remotion Studio forms regardless of whether the flag uses a space or `=`.
- Reject executable command-start platform opener forms without depending on a literal localhost URL.
- Preserve the documented safe foreground `npm run preview` workflow.
- Do not treat explanatory prohibition prose as executable workflow content.
- Add replayable positive and negative mutation cases for every widened matcher family.

## Red-Green Proof

**RED prompt/case:** Feed the current unsafe-form detector `remotion studio src/Root.tsx --port 3000`, `remotion studio src/Root.tsx --port=3000`, and `open "$REMOTION_PREVIEW_URL"`; each prohibited executable form currently produces no match.
**Why RED now:** A future fixed-port or macOS-opener regression can pass the presentation contract suite even though the source-only video contract prohibits both workflows.
**GREEN when:** All executable fixed-port and platform-opener mutations fail, safe foreground preview examples pass, negative prose is ignored, and the focused and canonical suites remain green.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
