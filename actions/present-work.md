# Present Work Action

> **Part of the do-work skill.** Generates client-facing deliverables from completed work — briefs, architecture diagrams, value propositions, and video scripts. Turns the technical archive into something that educates and sells.

The code is done. Now communicate its value. This action reads the full history of completed work — what was requested, what was built, how it works, and what the Lessons Learned say — and produces artifacts that explain it to someone who doesn't read diffs.

## Philosophy

- **Educate first, sell second.** The client should understand what they got before hearing why it's valuable.
- **Technical accuracy in plain language.** Architecture and data flow matter, but describe them at the level of components and interactions, not code.
- **Honest value, not hype.** Say "enables X" not "will increase revenue by Y%." Don't fabricate metrics.
- **Pointers over prose.** Reference key files and code — the implementation is the source of truth. Don't rewrite it in paragraphs.
- **Proportional effort.** A config change gets a 2-paragraph brief. A multi-feature system gets architecture diagrams and a video script.

## When to Use

**Use when:**
- Work is **completed and archived**, and the user wants a client-facing or stakeholder-facing summary.
- A pipeline just finished and the user wants briefs/videos/diagrams as the closing step.
- The user asks for a portfolio summary across everything shipped to date (`do-work present all`).

**Do NOT use when:**
- Work is still in progress — the archive isn't populated, and briefs would be speculative.
- The user wants *technical review*, not a brief — use `do-work review work` or `do-work code-review`.
- The REQ has no user-visible output (infra-only, tooling) — a brief adds no value; skip.

## Two Modes

| Mode | Trigger | What it does |
|------|---------|-------------|
| **Detail** | `do-work present UR-003`, `do-work present REQ-005`, or `do-work present` (most recent) | Deep dive on specific completed work |
| **Portfolio** | `do-work present all` or `do-work present portfolio` | Summary of all completed work across the archive |

## Detail Mode Workflow

### Step 1: Find the Target

Same pattern as review work standalone mode:

1. **If user specifies a REQ** (e.g., "present REQ-005"): Find it in `do-work/archive/` or `do-work/archive/UR-NNN/`
2. **If user specifies a UR** (e.g., "present UR-003"): Find all completed REQs under that UR — present them as one deliverable
3. **If no target specified**: Find the most recently completed UR (or REQ if no UR). Check `do-work/archive/` for the highest UR/REQ number with `status: completed`

If the target has no completed REQs, report that there's nothing to present and exit.

### Step 2: Read the Full History

For each completed REQ, read the full file and extract:

- **What was requested** — the What/Detailed Requirements sections
- **Original input** — read the UR's `input.md` for the user's own words
- **Triage** — what route was chosen and why
- **Plan** — what was planned (Route C)
- **Exploration** — what was discovered about the codebase
- **Implementation Summary** — mandatory file manifest (files changed, what was done)
- **Testing** — what tests exist and pass
- **Review** — scores, findings, acceptance result
- **Lessons Learned** (if present) — what worked, what didn't, gotchas. Route A REQs may skip this section.

### Step 3: Read the Code

Use `git show <commit>` (from the REQ's `commit` frontmatter) to get the diff. Then read the actual created/modified files to understand:

- **Architecture** — what components exist, their roles, how they connect
- **Data flow** — how information enters, moves through, and exits the system
- **Patterns used** — frameworks, libraries, conventions followed
- **Scale** — how many files, how much new code, how much modified

Don't just summarize the diff — understand the system that was built.

### Step 4: Generate Artifacts

Based on the work scope, generate the appropriate deliverables:

#### 4a: Client Brief (always generated)

Write to `do-work/deliverables/UR-NNN-client-brief.md` (or `REQ-NNN-client-brief.md` for standalone REQs):

```markdown
# [Feature/Project Name]

## What We Built

[1-2 paragraphs. Plain language. What it does from the user's perspective.
No code, no jargon — what problem it solves and how the user interacts with it.]

## How It Works

### Architecture

[ASCII diagram: components, their roles, how they connect.
Use boxes, arrows, labels. Keep it to one screen.]

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  User Input │────►│  Processing  │────►│   Output    │
│  (Form/API) │     │  (Service)   │     │  (UI/File)  │
└─────────────┘     └──────┬───────┘     └─────────────┘
                           │
                    ┌──────▼───────┐
                    │   Storage    │
                    │  (Database)  │
                    └──────────────┘
```

### Data Flow

[Numbered walkthrough showing how data moves through the system:]
1. User submits [input] via [interface]
2. [Component] validates and transforms the data
3. [Service] processes it by [doing what]
4. Result is stored in [where] and returned to [whom]

### Key Design Decisions

- [Why this approach over alternatives — what made it the right call]
- [Trade-offs made and why they're acceptable]
- [What the Lessons Learned revealed about the approach]

## Why This Works

[1 paragraph. Business-level justification. Why this solution solves the
problem well, not just technically but for the business.]

## Value Delivered

### Immediate Impact
- [What the client gets right now — specific capabilities]
- [Problems solved — pain points eliminated]

### Revenue & Growth Opportunities
- [How this enables the client to make more money]
- [New capabilities unlocked — what's now possible that wasn't before]
- [Efficiency gains — time saved, errors prevented, processes streamlined]

### Competitive Advantage
- [What this gives the client that their competitors likely don't have]

## Key Files

[Pointers to the most important files in the implementation — pull from
the REQ's ## Implementation Summary. These are the source of truth.]

- `src/components/Feature.tsx` — main component
- `src/services/feature-service.ts` — business logic
- `tests/feature.spec.ts` — test coverage

## How to Verify

[Concrete steps someone can follow to check that the work is real and correct.
Use copy-pasteable commands — no abstractions. Pull commits from each REQ's
frontmatter.]

1. **See the change in the codebase:**
   ```
   git show {commit-sha}   # {one-line: what this commit did}
   ```
2. **Run the tests that prove it works:**
   ```
   {project's test command — pull from prime files or REQ Testing section}
   ```
3. **Try it for yourself:**
   - [User-facing steps: open page X, click Y, expect Z]
   - [Or: call endpoint with curl, check response]

## What's Next

- [Follow-up opportunities — natural extensions of this work]
- [Phase 2 ideas — features that build on what was just shipped]
- [Quick wins — small additions that compound the value]

## Related Reading

Cross-link sibling deliverables so the reader can drill in or zoom out. Render as a markdown list of relative links; include only files that actually exist in `do-work/deliverables/` — check the folder before writing this section.

- [`{UR-NNN}-interactive-explainer.single.html`](./{UR-NNN}-interactive-explainer.single.html) — click-through Before/After demo in any browser
- [`{UR-NNN}-video/`](./{UR-NNN}-video/) — animated walkthrough (run `npm install && npm run preview`)
- [`{UR-NNN}-pipeline-summary.md`](./{UR-NNN}-pipeline-summary.md) — developer-facing debrief: commits, test deltas, how to verify *(exists if this work was produced by a pipeline run)*
- [`{UR-NNN}-pipeline-summary.single.html`](./{UR-NNN}-pipeline-summary.single.html) — same debrief, browser-readable *(same condition)*
- [`{UR-NNN}-pipeline-summary.marp.html`](./{UR-NNN}-pipeline-summary.marp.html) — stakeholder-facing Marp slide deck, exported to HTML *(same condition)*
```

#### 4b: Remotion Video (when the feature is user-facing/demo-able)

Generate a Remotion video project for browser preview. **No mp4 rendering** — the deliverable is the Remotion source code itself, previewed via `npx remotion studio`.

Write the video project to `do-work/deliverables/UR-NNN-video/` (or `REQ-NNN-video/` for standalone REQs):

**Project structure:**

```
do-work/deliverables/UR-NNN-video/
├── package.json
├── src/
│   ├── Root.tsx              # Remotion <Composition> entry point
│   ├── Video.tsx             # Main composition — sequences all scenes
│   ├── scenes/
│   │   ├── ProblemScene.tsx   # Scene 1: The pain point
│   │   ├── SolutionScene.tsx  # Scene 2: The feature demo
│   │   ├── ArchScene.tsx      # Scene 3: Under the hood
│   │   └── ValueScene.tsx     # Scene 4: Business value + CTA
│   └── styles.ts             # Shared colors, fonts, layout constants
└── tsconfig.json
```

**`package.json`:**

```json
{
  "name": "[ur-or-req-id]-video",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "preview": "npx remotion studio src/Root.tsx --no-open & sleep 3 && open http://localhost:3000"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "remotion": "^4.0.0",
    "@remotion/cli": "^4.0.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/react": "^18.2.0"
  }
}
```

**`src/Root.tsx`:**

```tsx
import {registerRoot, Composition} from 'remotion';
import {Video} from './Video';

const RemotionRoot: React.FC = () => {
  return (
    <Composition
      id="[FeatureName]"
      component={Video}
      durationInFrames={[total frames — 30fps * seconds]}
      fps={30}
      width={1920}
      height={1080}
    />
  );
};

registerRoot(RemotionRoot);
```

**Important:** The `registerRoot()` call is required — without it Remotion Studio will hang on "Waiting for registerRoot() to get called." Do NOT use `export const RemotionRoot` — call `registerRoot()` at module level instead.

**`src/Video.tsx`:**

Compose all scenes using `<Sequence>` with appropriate frame offsets:

```tsx
import {AbsoluteFill, Sequence} from 'remotion';
import {ProblemScene} from './scenes/ProblemScene';
import {SolutionScene} from './scenes/SolutionScene';
import {ArchScene} from './scenes/ArchScene';
import {ValueScene} from './scenes/ValueScene';

export const Video: React.FC = () => {
  return (
    <AbsoluteFill style={{backgroundColor: '#0f172a'}}>
      <Sequence from={0} durationInFrames={[scene1 frames]}>
        <ProblemScene />
      </Sequence>
      <Sequence from={[scene1 end]} durationInFrames={[scene2 frames]}>
        <SolutionScene />
      </Sequence>
      <Sequence from={[scene2 end]} durationInFrames={[scene3 frames]}>
        <ArchScene />
      </Sequence>
      <Sequence from={[scene3 end]} durationInFrames={[scene4 frames]}>
        <ValueScene />
      </Sequence>
    </AbsoluteFill>
  );
};
```

**Scene components — guidelines:**

Each scene is a React component using Remotion primitives. Build visuals with:

- `<AbsoluteFill>` for full-frame layouts
- `useCurrentFrame()` and `interpolate()` for animations (fade in/out, slide, scale)
- `<Sequence>` for staggered element timing within a scene
- `spring()` for natural motion on key elements
- Inline styles or the shared `styles.ts` constants — no external CSS frameworks
- Text, boxes, arrows, and simple shapes to illustrate concepts — no external images unless the project has them

**Scene content mapping:**

| Scene | Duration | Content |
|-------|----------|---------|
| **ProblemScene** | ~10-15s (300-450 frames) | Animate the pain point — what was missing, broken, or manual. Use text reveals, fading lists, or a "before" state visualization. |
| **SolutionScene** | ~20-30s (600-900 frames) | Walk through the feature. Animate a simulated UI or step-by-step flow. Highlight the "aha" moment with emphasis animation. This is the longest scene. |
| **ArchScene** | ~15-20s (450-600 frames) | Animate the architecture diagram from the client brief — boxes appearing, arrows drawing between them, labels fading in. Show data flow as animated paths. |
| **ValueScene** | ~10-15s (300-450 frames) | Before/after comparison or animated capability list. End with the project name and a call to action. |

**Scene component example (ProblemScene):**

```tsx
import {AbsoluteFill, useCurrentFrame, interpolate, Sequence} from 'remotion';

export const ProblemScene: React.FC = () => {
  const frame = useCurrentFrame();
  const titleOpacity = interpolate(frame, [0, 30], [0, 1], {
    extrapolateRight: 'clamp',
  });

  return (
    <AbsoluteFill
      style={{
        backgroundColor: '#1e293b',
        justifyContent: 'center',
        alignItems: 'center',
        padding: 80,
      }}
    >
      <h1
        style={{
          color: '#f8fafc',
          fontSize: 72,
          fontFamily: 'system-ui, sans-serif',
          fontWeight: 700,
          opacity: titleOpacity,
          textAlign: 'center',
        }}
      >
        [The problem statement — concise, relatable]
      </h1>
      <Sequence from={45}>
        {/* Staggered bullet points, animated elements, etc. */}
      </Sequence>
    </AbsoluteFill>
  );
};
```

**`src/styles.ts`:**

```ts
export const COLORS = {
  bg: '#0f172a',
  surface: '#1e293b',
  primary: '#3b82f6',
  accent: '#10b981',
  text: '#f8fafc',
  muted: '#94a3b8',
};

export const FONTS = {
  heading: 'system-ui, -apple-system, sans-serif',
  body: 'system-ui, -apple-system, sans-serif',
  mono: 'ui-monospace, monospace',
};
```

**Key rules for generating scenes:**

- **Use real content from the completed work.** Scene text comes from the REQ's What/Implementation Summary sections and the client brief — not placeholder lorem ipsum.
- **Animate meaningfully.** Every animation should reveal information progressively. Don't animate for the sake of it.
- **Keep it readable.** Large fonts (48px+ for body, 72px+ for headings), high contrast, minimal text per frame.
- **No external assets.** Everything is built with React components, styled divs, and SVG. No image imports, no video embeds.
- **Match the architecture diagram.** The ArchScene should visualize the same components and data flow from the client brief's architecture section, animated.

**When to generate a Remotion video:**
- The feature has visible UI or user-facing output
- There's a clear before/after to demonstrate
- The work is substantial enough to warrant a walkthrough (Route B or C)

**When to skip:**
- Backend-only changes, config tweaks, refactors, infrastructure
- Bug fixes that aren't visually interesting
- Route A changes (too small)

#### 4c: Interactive Explainer (Single-File HTML)

Always generate a self-contained HTML file at `do-work/deliverables/UR-NNN-interactive-explainer.single.html` (or `REQ-NNN-interactive-explainer.single.html`). The `.single.html` suffix marks this as LLM-authored standalone HTML — distinct from any `.marp.html` file, which is the mechanical output of `marp-cli` converting a `.marp.md` source.

**Guidelines for the Interactive Explainer:**

- **Zero build steps:** It MUST be a single `.html` file. Use standard HTML5, Tailwind CSS via CDN (`<script src="https://cdn.tailwindcss.com"></script>`), and Vanilla JavaScript. No React, Vite, or npm installs required.
- **Goal:** Visually explain the problem (Before) and the solution (After) to a non-technical stakeholder.
- **Interactivity:** You MUST include interactive elements. Examples: A 'Before / After' visual toggle slider, or a clickable 'Step-by-Step' data flow diagram where clicking 'Next' highlights different parts of the architecture.
- **Theme:** Light theme by default. Use `prefers-color-scheme: dark` media query to support OS-level dark mode. Define CSS custom properties (e.g., `--bg`, `--surface`, `--text`) at `:root` for light values and override them inside `@media (prefers-color-scheme: dark)`. Light palette: white/slate-50 backgrounds, slate-800/900 text, blue-600 accents. Dark palette: slate-900 backgrounds, slate-100 text, blue-400 accents.
- **Design:** Make it look modern and highly polished — large typography, soft shadows, generous whitespace. Include tooltips or sidebars that explain technical decisions in plain English.
- **Content:** Pull real context from the REQ files. Include a 'The Problem', 'The Interactive Demo', and a 'Value Delivered' section. Also include a collapsible (or small-print, bottom-of-page) 'For the developer' section listing the commit SHAs from each REQ's frontmatter, a copy-pasteable `git show <sha>` block, and the project's test command. A dev who landed on the explainer shouldn't have to leave to verify it's real — the explainer serves both the "no clue" reader and the "show me the receipts" reader in one file.
- **Navigation footer:** End the page with a "Keep exploring" section — a responsive card grid of `<a>` links to sibling deliverables that exist in the same folder. Always link the client brief (`./{UR-NNN}-client-brief.md`) with a note that GitHub/VS Code renders markdown natively. When the pipeline ran, also link `./{UR-NNN}-pipeline-summary.single.html` (developer debrief), `./{UR-NNN}-pipeline-summary.marp.html` (stakeholder deck), and `./{UR-NNN}-video/` (walkthrough). Check the `do-work/deliverables/` folder before rendering — only include tiles for files that actually exist. This is the reader's escape hatch from the explainer's breadth-first view into deeper, audience-specific context.

#### 4d: Portfolio artifacts (portfolio mode only — see below)

### Step 5: Save and Present

1. Create `do-work/deliverables/` directory if it doesn't exist
2. Write all generated artifacts to `do-work/deliverables/`
3. Present a summary to the user:

```
Generated deliverables for UR-003:

  do-work/deliverables/UR-003-client-brief.md                     Client brief with architecture + value prop
  do-work/deliverables/UR-003-video/                              Remotion video (4 scenes, ~90s)
  do-work/deliverables/UR-003-interactive-explainer.single.html   Interactive explainer (LLM-authored single-file HTML)

Each artifact links to its siblings — open any one as an entry point.

Key value points:
  - [Top 1-2 value propositions from the brief]

To preview the video:
  cd do-work/deliverables/UR-003-video && npm install && npm run preview

To view the interactive explainer:
  Simply double-click the HTML file to open the interactive presentation directly in any browser.
```

## Portfolio Mode Workflow

### Step 1: Scan the Archive

List all UR folders and completed REQs in `do-work/archive/`:
- Read each UR's `input.md` for the title and request list
- Read each completed REQ's frontmatter (title, route, commit, completed_at) and Review section (overall score)
- Check `do-work/archive/legacy/` for standalone completed REQs

### Step 2: Build the Overview

For each UR/REQ, extract a one-line summary from the What section and the review score.

### Step 3: Generate Portfolio Summary

Write to `do-work/deliverables/portfolio-summary.md`:

```markdown
# Work Portfolio

**Total completed:** [N] requests across [M] user requests
**Period:** [earliest created_at] — [latest completed_at]

---

## Completed Work

### UR-001: [Title from input.md]

[1-2 sentences on what was delivered — synthesized from the REQs' What sections]

| REQ | Title | Route | Review | Commit |
|-----|-------|-------|--------|--------|
| REQ-010 | [title] | C | 92% | abc1234 |
| REQ-011 | [title] | B | 88% | def5678 |

---

### UR-003: [Title]

[1-2 sentences]

| REQ | Title | Route | Review | Commit |
|-----|-------|-------|--------|--------|
| REQ-020 | [title] | A | 95% | ghi9012 |

---

## Cumulative Value Proposition

### What Was Built
[High-level summary across all completed work — the big picture.
What capabilities does the system now have?]

### Total Value Delivered
- [Aggregated business impact across all work]
- [Capabilities added — the full list]
- [Technical quality — average review scores, test coverage]

### Growth Opportunities
- [Cross-cutting opportunities that only become visible when viewing all work together]
- [Natural next phases based on what's been built]
- [Quick wins that compound existing value]

### Lessons Learned (Cross-Project)
- [Patterns that worked well across multiple REQs]
- [Common pitfalls encountered]
- [Architectural decisions that should inform future work]
```

### Step 4: Save and Present

Same as detail mode — save to `do-work/deliverables/` and summarize.

## Calibrating Depth

| Work scope | Artifact depth |
|-----------|---------------|
| Single Route A REQ (bug fix, config) | Minimal brief — 2-3 paragraphs, skip architecture diagram, skip video script, minimal value prop. |
| Single Route B REQ (clear feature) | Standard brief with architecture section. Video script if user-facing. |
| Route C or multi-REQ UR | Full brief with detailed architecture, data flow, value prop, design decisions. Video script for demo-able features. |
| Portfolio (all work) | Portfolio summary with per-UR breakdowns and cumulative value proposition. |

## What NOT to Do

- Don't include code snippets — this is client-facing, not developer-facing. Point to files instead.
- Don't oversell — be honest about what was built and its limitations
- Don't fabricate metrics — quantify where data exists, qualify where it doesn't
- Don't include internal details (review scores, triage routes, test commands) in client briefs — those are for the portfolio summary
- Don't generate a video script for non-visual changes
- Don't write walls of text when pointers to code would be more accurate and durable
- Don't regenerate deliverables that already exist without the user asking — check `do-work/deliverables/` first and offer to update

## Red Flags

- The target UR/REQ isn't in `do-work/archive/` yet — don't present work that isn't done; stop and tell the user.
- The brief quotes metrics that don't appear anywhere in the REQ, implementation summary, or commits — fabricated value prop; remove or replace with qualitative framing.
- The brief copies code snippets instead of pointing at files — client-facing docs shouldn't show diffs.
- You're about to overwrite an existing deliverable at `do-work/deliverables/` — pause and confirm with the user whether to update in-place or version it.
- Portfolio mode produced a brief for every archive entry including internal-only cleanup REQs — filter those out; portfolios are about shipped user value.

## Verification Checklist

- [ ] Target UR/REQ was found in `do-work/archive/` (not `working/`, not `queue/`).
- [ ] Deliverables saved under `do-work/deliverables/{UR-NNN}/` (or portfolio path for all-mode).
- [ ] No code snippets in the brief — only file path references.
- [ ] Value prop is qualitative unless a specific metric exists in the archive record.
- [ ] Architecture diagram was generated only when the work was non-trivial (Route B or C).
- [ ] Depth matched the work scope per the "Calibrating Depth" table.
- [ ] If a prior deliverable existed, it was updated in place or versioned — not silently duplicated.
