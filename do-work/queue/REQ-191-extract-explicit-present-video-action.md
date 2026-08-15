---
id: REQ-191
title: Extract an explicit standalone present-video action
status: pending
created_at: 2026-08-15T09:10:53Z
user_request: UR-042
domain: frontend
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-189]
maintenance: false
effort_estimate: normal
related: [REQ-189, REQ-190, REQ-192]
batch: completed-work-presentation-consolidation
write_set: [skills/do-work-toolbox/actions/present-video.md, skills/do-work-toolbox/docs/present-video-guide.md]
---

# Extract an Explicit Standalone Present-Video Action

## What

Move the existing Remotion specification into `do-work-toolbox present-video [UR|REQ|most recent]` and a dedicated guide. Generate a valid, source-only animated walkthrough only after an explicit video request, with no automatic invocation and no MP4 render path.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Video source generation is currently embedded in `present-work` and can run as part of a broader presentation workflow. A standalone action makes the cost and intent explicit and allows `ai-report` and `present-work` to stay video-free.

## Context

REQ-189 creates `actions/completed-work-presentation-reference.md`; this action must consume that reference rather than restating its archive and evidence rules. REQ-192 owns public routing and command discovery after the action and guide exist.

## Detailed Requirements

- Add `do-work-toolbox present-video [UR|REQ|most recent]` as a standalone action with its own guide.
- Run it only through an explicit `present-video`, `remotion`, or `video walkthrough` request.
- Resolve a completed UR, completed REQ, or most-recent completed target through the shared completed-work presentation reference.
- Accept `completed` and `completed-with-issues`, surface known issues honestly, and reject cancelled or unfinished work.
- Load anti-slop and prompt-injection guidance before reading archived content.
- Produce `do-work/deliverables/<ID>-video/`.
- Generate valid Remotion source only and never render an MP4.
- Move and preserve the useful existing Remotion contract:
  - Problem → Solution → Architecture → Value structure;
  - proportional generation rules so output depth matches the work;
  - real archived content rather than invented feature claims;
  - a module-level `registerRoot` requirement;
  - a no-external-assets rule;
  - qualitative value rules with no fabricated metrics.
- Use `remotion studio src/Root.tsx` as the package preview command.
- Keep the preview in the foreground.
- Do not background the preview process.
- Do not sleep and assume the server is ready.
- Do not assume port 3000.
- Do not invoke macOS `open`.
- Do not add any MP4 render command or render an MP4.
- Skip trivial or genuinely non-visual work with a concise explanation instead of generating filler.
- Honor the shared no-overwrite behavior for an existing output directory.
- Do not duplicate terminal-success resolution, archive-reading safety, required archive fields, merge-aware commit inspection, current-code inspection, crew loading, or no-overwrite rules from the shared reference.

## Constraints

- Detailed report means `ai-report`; cross-project portfolio means `present-work`; animated walkthrough means `present-video`.
- `present-video` is never called automatically by `ai-report`, `present-work`, or a completion flow.
- Preserve all existing generated artifacts; do not migrate, overwrite, or delete prior briefs, `.single.html` files, reports, or video directories.
- Do not change UR/REQ schemas, archive formats, `review-work`, or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, or automatic video generation.
- Use only local/generated source and content allowed by the no-external-assets rule.

## Dependencies

Depends on REQ-189, which creates the shared completed-work presentation reference this action must use.

## Builder Guidance

Firm intent. Extract the proven Remotion content proportionally and delete its unsafe preview assumptions. A skipped trivial or non-visual target is a successful concise outcome, not a reason to manufacture scenes.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Request a video walkthrough for a completed visual REQ and inspect the generated package, then request one for trivial or genuinely non-visual work; also search the command for backgrounding, sleeps, fixed port 3000, macOS `open`, and MP4 rendering.
**Why RED now:** Remotion generation exists only inside `present-work`, its preview script backgrounds Studio, sleeps, assumes port 3000, and invokes `open`, and there is no explicit standalone routing boundary.
**GREEN when:** Explicit `present-video`, `remotion`, or `video walkthrough` routing can generate a valid `<ID>-video` Remotion source tree with module-level `registerRoot`, real archived content, no external assets, and a foreground `remotion studio src/Root.tsx` preview; no MP4 exists or can be produced by the action; trivial/non-visual work is skipped concisely; no other action invokes it automatically.
**Validation:** Inferred during capture from the supplied acceptance tests.

## Full Context

See `do-work/user-requests/UR-042/input.md` for the complete verbatim request and batch constraints.

---
*Source: attached “do-work capture-request: Consolidate completed-work presentation around ai-repo…” specification.*
