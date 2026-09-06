## 0.305.26 — Restore Single-Image Staging Claim to Adjacent to Target (2026-09-06)

Corrected the single-image generation command description in `skills/do-work-toolbox/actions/ai-report-reference.md` to reflect that `generateImage` stages its invocation-private file adjacent to the target output path (`filepath.Dir(outputPath)`), rather than in the system temporary directory.
