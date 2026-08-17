# Present Video Action

> **Part of the do-work-toolbox skill.** Generates a source-only animated walkthrough for one completed UR or REQ when the user explicitly asks for `present-video`, Remotion, or a video walkthrough. It belongs in toolbox because it completes the package's presentation set beside `ai-report` and `present-work` while sharing their completed-work evidence machinery. User-facing walkthrough: [`docs/present-video-guide.md`](../docs/present-video-guide.md).

`present-video` is a direct, opt-in action. Never invoke it from `ai-report`, `present-work`, a completion flow, or another automatic workflow.

## Philosophy

- **Meaning before motion.** Animation reveals the shipped problem, solution, architecture, and value; it does not decorate weak evidence.
- **Source is the deliverable.** Generate a complete Remotion project that a user may preview manually, never rendered media.
- **Proportional depth.** Match scene detail and duration to the evidence instead of stretching small work into filler.

## When to Use

**Use when:**

- The user explicitly requests `present-video`, Remotion, or a video walkthrough for completed work.
- The completed work has visible behavior, a meaningful user flow, or evidence-backed architecture or data flow that benefits from animation.
- The target is substantial enough to support a concise Problem → Solution → Architecture → Value narrative.

**Do NOT use when:**

- The user wants detailed stakeholder HTML; use `do-work-toolbox ai-report <target>`.
- The user wants a cross-project portfolio; use `do-work-toolbox present-work all|portfolio`.
- The work is trivial or genuinely non-visual; return the successful skip described below instead of manufacturing scenes.

## Input

`$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank. Naming this action is the explicit authority to run it; do not infer or queue a video from completion state or from another presentation action.

## Steps

### Step 1: Resolve and Read the Completed Work

Read and follow [`completed-work-presentation-reference.md`](completed-work-presentation-reference.md) in full **before opening archived user content**. It is the sole contract for safety load order, target resolution, archive fields, missing evidence, merge-aware commit and current-code inspection, evidence honesty, terminal-success qualifications, and no-overwrite publication. Do not recreate those rules here.

Build the reference's provenance ledger for the selected work before judging eligibility or drafting scene content.

### Step 2: Apply the Walkthrough Eligibility Gate

Decide whether animation would materially clarify the completed work. Eligible evidence includes a shipped UI, an observable user journey, or a substantial architecture or data-flow change that can be explained accurately with animated components and connections. UI screenshots are not required.

Skip when the work is trivial or genuinely non-visual and animation would add only generic text reveals or invented imagery. A skip is a successful result: create no output path and state the canonical target ID plus one concise evidence-based reason.

Apply this gate before creating a directory or file.

### Step 3: Design a Proportional Evidence-Backed Narrative

Map the provenance ledger into exactly four ordered scenes:

| Scene | Purpose |
|---|---|
| **Problem** | Explain the recorded pain point or prior limitation. Do not invent a visual or factual before state. |
| **Solution** | Walk through the delivered behavior or flow and emphasize the central change. |
| **Architecture** | Show only verified components, boundaries, and data flow from the implementation and current code. |
| **Value** | End with qualitative stakeholder value and any recorded qualification or issue; never fabricate metrics. |

Choose each scene's frames from the amount of useful evidence. The solution may be longest when the delivered flow warrants it, but no scene receives filler to reach a preset runtime. Define named frame constants, derive each sequence offset from preceding durations, and derive the composition total from all four durations so the values cannot drift.

Keep on-screen copy concise and readable. Use progressive transforms and opacity to reveal meaning; do not animate layout properties merely for decoration.

### Step 4: Define the Complete Source Package

Generate this complete tree beneath the final output directory:

```text
package.json
tsconfig.json
src/
  Root.tsx
  Video.tsx
  styles.ts
  scenes/
    ProblemScene.tsx
    SolutionScene.tsx
    ArchitectureScene.tsx
    ValueScene.tsx
```

Every file must contain valid, complete JSON, TypeScript, or TSX—never placeholders, ellipses, prose tokens, or pseudocode. Use compatible local React, React DOM, Remotion CLI, Remotion, TypeScript, and React type dependencies. Keep the package private. The only package script is:

```json
{
  "scripts": {
    "preview": "remotion studio src/Root.tsx"
  }
}
```

Use a strict no-emit TypeScript configuration with React JSX support and include `src/`. The action creates no lockfile, installs no package, and launches no process.

`src/Root.tsx` must define a valid `<Composition>` using the same FPS, dimensions, and total frame constant as the sequences, then call `registerRoot(RemotionRoot)` at module scope. Exporting a root component without the module-level call is invalid because Studio will not register the composition.

`src/Video.tsx` must sequence `ProblemScene`, `SolutionScene`, `ArchitectureScene`, and `ValueScene` in that order using internally consistent named duration and offset constants. Each scene is a complete React component using Remotion primitives.

Keep the shared visual palette and typography constants in `src/styles.ts`. Build every visual from React elements, inline styles, shared style constants, system fonts, and inline SVG markup. Do not use network requests, CDNs, CSS frameworks, imported CSS, or imported image, video, audio, or font assets. Do not depend on a client brief or another generated artifact: scene claims come directly from the provenance ledger.

### Step 5: Publish the Source-Only Project

Use the resolved target's canonical archive ID for the preferred directory:

```text
do-work/deliverables/<canonical-ID>-video/
```

Use that directory as this consumer's preferred path and apply the shared reference's **Collision-Safe Publication** section to the entire source project before writing any file. Report the path that contract resolved.

Write only the source tree from Step 4. Do not install dependencies, start Studio, add a media-rendering script, invoke a renderer, or create MP4 or other rendered media.

### Step 6: Verify and Report the Result

Read every generated file and verify that:

- `package.json` parses, is private, has compatible dependencies, and has only the exact foreground preview script;
- `tsconfig.json` parses and covers the source tree;
- imports resolve within the declared tree and the TSX contains no placeholders;
- `registerRoot` runs at module scope;
- four sequence offsets, durations, and the composition total agree;
- every scene claim traces to the provenance ledger and recorded issues remain visible;
- no external asset, installation, launched process, render path, lockfile, or rendered media exists.

Print a compact result containing the output path, canonical target and status, total frames and FPS, a one-line scene summary, recorded issues, and the manual preview command `npm run preview`. Label preview as not run.

## Output Format

An eligible target creates one fresh, self-contained Remotion source project using `do-work/deliverables/<canonical-ID>-video/` as its preferred path. A successful skip creates nothing and explains why animation would not improve the evidence.

## Rules

- This action is explicit-only and source-only. It never runs as a completion side effect and never renders media.
- The shared completed-work presentation reference owns archive ingestion and publication mechanics; this action owns only walkthrough eligibility, narrative, and Remotion source shape.

## Verification Checklist

- [ ] Invocation was an explicit request for this action; no other action or completion flow delegated automatically.
- [ ] The shared reference was loaded before archived content and its target, evidence, issue, and **Collision-Safe Publication** contracts were followed.
- [ ] Trivial or genuinely non-visual work skipped before output creation; eligible architecture or data flow was not rejected merely for lacking UI.
- [ ] The project contains the complete declared tree, four ordered scenes, proportional consistent frame math, and module-level `registerRoot`.
- [ ] Claims and qualitative value are evidence-backed, with no invented before state or metric.
- [ ] Visuals use only React/CSS-in-JS/system fonts/inline SVG and no external or imported media assets.
- [ ] The only package script is exactly `remotion studio src/Root.tsx`; the action did not install, launch, render, or create media.
