---
title: "Lessons from REQ-202: Complete unsafe Remotion preview mutation detection"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-202-complete-unsafe-remotion-preview-mutatio.md]
related:
  - page: concept-completed-work-presentation
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-202: Complete unsafe Remotion preview mutation detection

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Make the completed-work presentation regression detector recognize every executable fixed-port and platform-opener form while preserving safe foreground preview commands and non-executable prohibition prose.

This is one sweep because the missed forms share a single root cause: the unsafe-command extractor matches a narrow set of literal spellings instead of the complete prohibited executable command families.

## Solution summary

The live presentation files and in-memory replay cases now use one detector seam. Unsafe mutations cover separated, equals, quoted, and shell-continued numeric port values plus direct, indented, and chained variable opener targets. Prohibition context now spans only its immediate Markdown example list, while controls retain foreground Studio/package preview commands, same-line and multi-line prohibition prose, and ordinary opener prose.

## What worked

**What worked:** Extracting a source-text detector made executable-safety rules directly mutation-testable without modifying shipped presentation instructions. Family-labeled assertions also prevented one unsafe pattern from accidentally masking another.

**What didn't:** The first matcher expansion overfit the exact one-line examples. Independent adversarial review exposed shell-significant indentation and continuation plus a multi-line documentation shape that the initial tests omitted.

**Worth knowing:** Safety matchers for Markdown command examples need two bounded grammars: shell continuation for executable content and local structural continuation for explanatory prohibition examples. Crossing arbitrary newlines in either direction creates false negatives or false positives.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Back-reference

See `do-work/archive/UR-042/REQ-202-complete-unsafe-remotion-preview-mutation-detection.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `536fbd6`.
