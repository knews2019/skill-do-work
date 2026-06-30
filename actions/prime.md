# Prime Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to create or audit prime files. Prime files (`prime-*.md`) are AI context documents — semantic indexes that help an AI coder navigate a utility in minimum tokens. User-facing walkthrough: [`docs/prime-guide.md`](../docs/prime-guide.md).

## When to Use

**Use when:**
- A utility or directory is heavily used by builders and would benefit from a concise routing index (`prime create`).
- Prime files exist and may have drifted from the code — stale refs, broken links, missing primes for new utilities (`prime audit`).
- Onboarding a new directory of AI context docs and you want a health check across all of them.

**Do NOT use when:**
- The utility is tiny (a handful of files) — a prime file adds overhead without saving tokens.
- The user wants to *read* a prime file — that's just a file read, not this action.
- The user wants to build something from a prime file — use `do-work run` or `do-work capture-request:` and let the builder load the prime as context.

## Sub-Commands

The `prime` command accepts a sub-command as its first argument. If no sub-command is given, show the help menu.

| Sub-command | What it does |
|---|---|
| `create <path>` | Generate a new prime file for a utility via interactive Q&A |
| `audit` | Health-check all prime files (stale refs, missing primes, broken links) **and refresh their `## Stakes`** (Req/Value/Risk for decisions) |
| (none) | Show help menu |

---

## Help Menu (no sub-command)

When invoked with no sub-command (`do-work prime`), show:

```
prime — manage AI context documents (prime files)

  do-work prime create src/auth/    Generate a prime file via interactive Q&A
  do-work prime audit               Health-check all primes + refresh their Stakes (value/risk)
```

---

## Sub-Command: `create <path>`

Generate a prime file for the utility at `<path>` — a routing index that helps an AI coder navigate the utility in minimum tokens.

### Principles

- **Target: 15-30 lines** for the routing index. Every line must save the AI more tokens than it costs to read. (A `## Stakes` section, if present, is decision-context for the user — it sits **outside** this budget; see Step 4.)
- The AI has `Read`, `Grep`, `Glob`. Don't reproduce what tool calls discover.
- Only include what the AI CANNOT efficiently find: routing, traps, exclusions.
- NO: line numbers, code descriptions, DOM anchors, request flow diagrams, URL params, environment tables, external service catalogs. The AI will discover these via tool calls.
- Follow the **PRIME Files Philosophy** from `crew-members/general.md`: low noise, high value, pointers not copies, no volatile metrics.
- Save to `{utility}/prime-{short-name}.md`

### Workflow

#### Step 1: Scan

`Glob` and `Read` the utility directory. Identify:
- Entry points (the 2-4 files an AI should read first)
- Build system (if any)
- Generated/vendor/dead files the AI should skip

#### Step 2: Report

Show the user a 3-line summary: entry points, build system, files to skip. Then proceed to questions.

#### Step 3: Ask 3 Questions

Use your environment's ask-user prompt/tool to ask these questions (free-text answers):

**Q1: "Which files should an AI read first to understand this utility?"**
- The files that contain the core logic, not config or boilerplate
- This becomes the **Read first** section
- Files that drive the rest of the code so that the project becomes discoverable

**Q2: "What files or code sections should the AI NOT edit?"**
- Dead code kept for reference, generated output, vendor files, config that looks editable but isn't
- This becomes the **Do not edit** section

**Q3: "What traps would waste an AI's debugging time?"**
- Dev/prod differences, path resolution gotchas, naming that misleads, two-directory patterns
- This becomes the **Traps** section

#### Step 4: Generate

Combine auto-detected facts with user answers. Apply these rules:
- If a section would be empty, omit it entirely
- If there's no build step, omit **Must build**
- Every line must earn its place — "would an AI waste tokens without this?"
- **Spelunk the Stakes.** For the utility's *load-bearing* elements only (the few whose contract, if changed wrong, has real blast radius), read the code and record `Req:` (what it must do / why it exists), `Value:` (what it enables), `Risk:` (what breaks if changed wrong; reversibility). This becomes the `## Stakes` section — its purpose is to let the user make high-impact decisions confidently without re-deriving the context, and it feeds the Value/Risk that `do-work clarify` surfaces. Load-bearing only; if every element seems to qualify, the utility is too big — split it. No volatile metrics, pointers not copies. Omit the section entirely for a utility with no high-stakes elements.

#### Step 5: Write

Write to `{path}/prime-{short-name}.md` using this template:

```md
# Prime: {short-name}

{One line: what this utility is and where it lives.}

## Read first
- `{file}` — {why this one, max 8 words}
{2-4 files max}

## Do not edit
- `{file-or-pattern}` — {why}

## Must build
`{one-liner command}`

## Traps
- **{symptom}** — {cause and fix, one line}

## Stakes
- `{load-bearing element}`
  Req:   {what it must do / why it exists}
  Value: {what it enables}
  Risk:  {what breaks if changed wrong; reversibility}
```

(Include `## Stakes` only for load-bearing elements; omit for utilities with no high-stakes surface. It is the one section exempt from the 15-30 line budget.)

#### Step 6: Post-creation checks

1. Show the user the generated file and ask if anything is missing.
2. Check whether the prime should be registered in CLAUDE.md:
   - **Utility-specific primes** (live in a utility root): NOT registered — discovered by convention via glob
   - **Cross-cutting primes** (shared docs not in a utility root): SHOULD be registered in CLAUDE.md if it has a prime registry section
3. Check if sibling primes exist in the same area. If so, add cross-links.
4. If the area now has 3+ primes, check whether an area index prime exists. If not, suggest creating one.

#### Report

After writing the file, output:

```
Prime created: {path}/prime-{short-name}.md ({line count} lines)
Sections: {list of included sections}
```

---

## Sub-Command: `audit`

Audit the repo's prime file system, **then refresh each prime's `## Stakes`.** Two jobs: (1) a **read-only health check** of the routing index — staleness, missing coverage, broken references — and (2) a **write**: spelunk each flagged prime's load-bearing elements and refresh its `## Stakes` (Req / Value / Risk) so the decision context the user relies on stays current. The routing index is never rewritten; `## Stakes` is. **Updating Stakes is a core purpose of `audit`, not an afterthought** — don't run this as a pure read-only pass.

**What audit writes:** the routing sections (Read first / Do not edit / Must build / Traps) are **read-only** — audit reports issues there and lets the user decide what to fix. The `## Stakes` section **is written** — audit spelunks and refreshes it (Step 6.5), because keeping the value/risk current is half of why audit exists. So "audit is read-only" applies *only* to the routing index; **audit does update the prime's Stakes** — never tell the user the audit is read-only across the board.

### Conventions

If CLAUDE.md has a section describing prime file conventions, read it to understand the project's specific rules. The general conventions are:
- **Utility-specific primes:** `<utility-dir>/prime-<name>.md` — discovered by convention (recursive glob), NOT registered in CLAUDE.md
- **Satellite docs:** `known-bugs-<name>.md`, `lessons-learned/<topic>.md` — live alongside the prime
- **Cross-cutting primes:** Registered in CLAUDE.md's prime registry section (if one exists) — only for shared docs that don't live in a utility root
- **Cross-linking:** Primes in the same area must cross-link to each other (not just operational dependencies)
- **Area indexes:** Areas with 3+ primes need one prime that lists all related primes as the entry point

### Step 1: Discover all prime files

```
glob **/prime-*.md
```

Build a table of every prime file with columns: path, utility it documents, last modified.

Skip directories that are clearly not source primes — build output (`dist/`, `build/`, `.next/`), dependencies (`node_modules/`, `vendor/`), and session/scratch artifacts (temp directories, `.cache/`).

### Step 2: Validate each prime file

For each prime file, check:

1. **Key files still exist** — read the prime, extract any file paths it references (e.g., "Read first: `src/index.js`"), verify those files still exist on disk via glob. Flag any that are missing.

2. **Utility directory still exists** — confirm the parent utility directory is populated (not empty/deleted).

3. **Internal links valid** — check any relative markdown links (`[text](path)`) in the prime resolve to real files.

4. **CLAUDE.md link correct** — if the prime has a link back to CLAUDE.md, verify the relative path depth is correct for its location.

5. **No absolute paths** — grep for `file:///` URLs in the prime. All links must be relative from the prime's directory. Flag any absolute paths as portability violations.

6. **Cross-links present** — primes in the same area (sharing a parent directory tree) should cross-link to each other, not just for operational dependencies. If two primes are siblings or cousins in the same area (e.g., multiple primes under the same utility root), flag missing cross-links.

7. **Area index exists** — if an area (parent directory tree) contains **3 or more primes**, one prime should serve as the **area index** that lists all related primes. An area index prime is identified by: (a) a filename matching `prime-*-index.md` (e.g., `prime-auth-index.md`), or (b) containing a `## Related Primes` or `## Index` section that lists other primes in the same area. Check if such an index prime exists and whether it lists all primes in that area. Flag areas with 3+ primes but no index prime.

### Step 3: Find utilities without primes

Identify directories that have source code but no prime file. A utility-sized directory typically has its own entry point, build config, or package manifest.

Look for directories containing source files (`.php`, `.js`, `.ts`, `.py`, `.go`, `.rs`, `.rb`, etc.) but no `prime-*.md`. Use the project's directory structure to identify utility-sized units — directories with their own `package.json`, `composer.json`, `Makefile`, `index.*` entry point, or similar markers of an independent unit.

Skip directories that are clearly not utility roots: `node_modules/`, `vendor/`, `dist/`, `build/`, `.next/`, `.git/`, test fixture directories.

If CLAUDE.md defines known utility locations or directory conventions, use those as the primary scan targets. Otherwise, scan the project root with reasonable depth limits (2-3 levels).

Focus on directories that represent distinct utilities or modules — not every subdirectory needs a prime. A directory is a good "missing prime" candidate if:
- It has 5+ source files
- It has its own entry point or build config
- An AI would need multiple tool calls to understand its structure

Report these as "missing prime" candidates.

### Step 4: Audit satellite docs

```
glob **/known-bugs-*.md
glob **/lessons-learned/**/*.md
```

For each satellite doc, verify its parent directory also has a prime file. Flag orphaned satellites (satellite exists but no prime in the same utility).

### Step 5: Verify CLAUDE.md registry (if applicable)

Read CLAUDE.md and check whether it has a prime file registry section. If it does:
1. Every path listed in the registry points to a real file
2. No utility-specific primes have been accidentally registered (they should be discovered by convention)
3. The convention description is still accurate

If CLAUDE.md has no prime registry section, skip this step.

### Step 6: Content freshness spot-check

For each prime, do a quick sanity check:
- If it references a dev server command, does that script still exist?
- If it references specific config files, do they still exist?
- If it has a "Do not edit" section with vendor files, are those files still present?

Don't read every line of source code — just verify the pointers are valid.

### Step 6.5: Refresh Stakes — audit WRITES here (don't skip or downgrade to a flag)

Spelunk and **write** `## Stakes` so the user can make high-impact decisions confidently. This is half of why audit exists — the producer absorbs the cost of clarity so the reader doesn't (`crew-members/anti-slop.md` § 1). **Actually edit the file**; do not merely report that Stakes is missing or stale.

Scope the write to the primes that need it (don't re-spelunk a prime whose Stakes already matches the code):
- **Add when missing.** If a prime documents a load-bearing utility but has no `## Stakes`, spelunk its load-bearing elements and **add** the section — `Req:` (what it must do / why it exists), `Value:` (what it enables), `Risk:` (what breaks if changed wrong; reversibility).
- **Refresh when stale.** If an existing Stakes entry points at a removed file or a requirement that no longer holds, **rewrite** it from the current code.
- **Leave current Stakes alone.** A prime whose Stakes still matches the code needs no write — note it as current.
- **Load-bearing only; no volatile metrics; pointers not copies.** Keep `## Stakes` outside the 15-30 line routing-index budget. If every element seems to qualify, the utility is too big — split it.

`## Stakes` is the only section audit writes — the routing sections stay read-only. **Report what you wrote** so the write is visible: end the audit with `Stakes: added M, refreshed N, current K`.

### Output Format

Report findings as a structured checklist:

```markdown
## Prime Audit Report — YYYY-MM-DD

### Summary
- Total primes found: N
- Healthy: N
- Issues found: N
- Utilities missing primes: N
- Stakes: added M, refreshed N, current K   ← the write audit performs

### Issues

#### Stale references
- [ ] `path/prime-foo.md` references `src/old-file.js` which no longer exists
- [ ] ...

#### Stale or missing Stakes
- [ ] `path/prime-foo.md` Stakes for `src/x.ts` references a removed contract — refreshed
- [ ] `path/prime-bar.md` documents a load-bearing utility but has no Stakes — added
- [ ] ...

#### Missing primes
- [ ] `web/.../some-utility/` has source files but no prime
- [ ] ...

#### Broken links
- [ ] `path/prime-foo.md` has broken link to `../../CLAUDE.md` (wrong depth)
- [ ] ...

#### Absolute paths (portability)
- [ ] `path/prime-foo.md` uses `file:///Users/...` absolute URL — must be relative
- [ ] ...

#### Missing area indexes
- [ ] `web/checkout-eet/` has N primes but no index prime listing them all
- [ ] ...

#### Orphaned satellites
- [ ] `path/known-bugs-foo.md` exists but no prime in that directory
- [ ] ...

#### CLAUDE.md registry
- [ ] All registered paths valid: YES/NO
- [ ] No utility-specific primes registered: YES/NO

### Recommendations
[Actionable next steps — which primes to update, which to create, which links to fix]
```

Be concise. Only flag actual issues. "Everything looks fine" for a prime is not worth reporting individually.

## Red Flags

- A newly created prime file's *routing index* (everything but `## Stakes`) is longer than 30 lines — it's drifting into documentation; tighten it or split the utility.
- `audit` reports no issues across the whole repo, but several primes haven't been touched in months — they're likely stale; spot-check before trusting the clean result.
- `create <path>` was run on a path that already has a prime — avoid silent overwrite; ask before replacing.
- A prime file lists line numbers or reproduces code — violates the "pointers over copies" principle; rewrite.
- Audit flags missing primes for paths the user doesn't care about (experimental, deprecated) — add an exclusion rather than forcing creation.

## Verification Checklist

- [ ] Prime files' routing index is 15–30 lines (create mode); a `## Stakes` section, if present, is excluded from that budget and scoped to load-bearing elements.
- [ ] `audit` actually wrote `## Stakes` where missing/stale (added or rewrote it, not just flagged it) and reported the `added/refreshed/current` counts — it did not describe itself as read-only across the board.
- [ ] No line numbers, no reproduced code, no volatile metrics in the generated prime.
- [ ] `audit` output names each issue by file path and type (stale ref / broken link / missing prime).
- [ ] "Everything looks fine" primes are omitted from the audit report (only issues are listed).
- [ ] Generated primes follow the PRIME Files Philosophy from `crew-members/general.md`.
