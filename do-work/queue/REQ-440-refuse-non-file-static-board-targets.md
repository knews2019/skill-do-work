---
id: REQ-440
title: '[impact-critical] Refuse non-file static board output targets'
status: pending
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
related: [REQ-437, REQ-438, REQ-439, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
---

# Refuse Non-File Static Board Output Targets

## What

Refuse static-board publication when any generated output name already exists as a directory, symlink, or other non-regular object. Validate all three targets before the first backup rename so successful cleanup can never recursively delete user content hidden beneath an output filename.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Verbatim claim / severity:** `[P1] Refuse non-file static output targets before backing up.`
- **Evidence:** `generate.go:479-504` accepts any successful `Lstat`, renames the inode below the private directory, publishes a file, and later removes the private directory recursively.
- **Origin / earned by:** `803e4e77`/REQ-183 (Static Board Generation Can Publish a Mixed Three-File Bundle) introduced recovery backups but treated every existing object as replaceable. An isolated replay replaced an `index.html` directory and deleted its nested file while returning success.
- **Surface-cost:** Earned. The exact data-loss replay justifies one preflight over three fixed targets plus directory/symlink regressions; this is much cheaper than deleting arbitrary directory contents.

## Detailed Requirements

- Inspect all generated target paths before moving any existing target.
- Permit only absent targets and existing regular non-symlink files.
- Refuse directories, symlinks, and special files without mutating any target or leaving private publication residue.
- Preserve current regular-file publication and rollback behavior.

## Constraints

- Keep the existing three-output transactional publication shape.
- The validation is fixed-scope and must not become a configurable filesystem policy layer.

## Red-Green Proof

**RED prompt/case:** Replace `index.html` with a directory containing `kept.txt`, run static generation, and inspect the target and private scratch paths.
**Why RED now:** Publication moves the directory into its private backup and successful cleanup recursively deletes it.
**GREEN when:** Generation refuses before any rename, the directory and nested bytes remain unchanged, no private residue remains, and ordinary regular-file regeneration still succeeds.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. A preflight regular-file/no-symlink check before the first rename is the accepted remedy.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 19 from the validated external feedback.*
