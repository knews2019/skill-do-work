# REQ-230 Hand-Back — Point caller docs at the canonical publication rationale

**Branch:** `worktree-agent-REQ-230-point-caller-docs`
**Implementation commit:** `6f65763` — `[REQ-230] Point present-work at the canonical publication rationale`

## File Manifest

- **Modified** `skills/do-work-toolbox/actions/present-work.md` — replaced the restated `ln`/`mv` container-not-a-collision rationale in Step 6 with a pointer to the canonical [Verified exact publication](../../do-work/docs/prescribed-shell-primitives.md#verified-exact-publication) section, keeping the action's own policy sentence (snapshot candidate advances a numeric suffix, canonical path fails closed, occupying directory untouched and nothing nested inside it).
- **Modified** `_dev/tests/prescribed-shell-canonicalization.sh` — added `'container rather than a collision'` to the stale-pattern list, so any future restatement in shipped markdown outside the canonical guide fails the suite.

No other file was created, modified, or deleted. `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, and `skills/do-work/CHANGELOG.md` were left untouched, as instructed.

## Red-Green Proof

### RED (pattern added first, restatement still present)

```
$ bash _dev/tests/prescribed-shell-canonicalization.sh
FAIL: skills/do-work-toolbox/actions/present-work.md restates canonical prescribed-shell rationale <container rather than a collision>; keep local intent and point at the guide.
EXIT=1
```

Exactly one failure, naming exactly the target file — confirming the pattern is neither inert nor over-broad across the rest of the shipped markdown tree.

### GREEN (restatement converted to a pointer, pattern still in the list)

```
$ bash _dev/tests/prescribed-shell-canonicalization.sh
Prescribed shell primitive canonicalization checks passed.
EXIT=0
```

## Verify

`bash _dev/tests/maintainer-verify.sh` from the worktree root — **exit code 0**, run unpiped.

Tail of the run:

```
maintainer-verify: aggregate contract suite
...
Prescribed shell primitive canonicalization checks passed.
...
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.528s
Maintainer verification passed.
```

## Integration Seams

None. Both edits are self-contained within the write set; no shared file needs a line from this REQ.

## Decisions Not Dictated by the REQ

1. **Pointer wording mirrors the canonical guide's own sentence.** The replacement reads "Both publication steps make the canonical [Verified exact publication](...) check. This helper's answers to it are that …", which is the same construction the canonical `## Portfolio summary publication` section already uses ("Each publication makes the [Verified exact publication](#verified-exact-publication) check, and this helper's answers to it are that …"). Reason: the two now say the same thing in the same words, so a future correction to either reads as a correction to the other rather than as a new divergence.
2. **Pattern placed next to the `curl -o` entry** in the stale-pattern list rather than appended at the end. Reason: that groups the two publication-side rationales together; the list is order-independent, so this is legibility only.
3. **Kept the local policy sentence verbatim.** The suffix-advance / fail-closed / nothing-nested clause is the action's policy, which the canonical guide explicitly delegates back to each helper's own section ("What a helper does *about* a nested write is its own policy and stays in its own section"). Removing it would have deleted policy, not a restatement.
4. **Did not widen the pattern beyond the REQ's phrase.** A broader pattern (e.g. `ln` + `mv` together) would have caught the same instance but risks false positives on legitimate helper-policy prose. The REQ named the phrase; the RED run proves it is sufficient.

## Discovered Tasks

- `skills/do-work-toolbox/actions/present-work.md`, the bullet immediately above the one this REQ fixed, restates the independent-bytes rationale ("A snapshot that shared storage with the canonical file would silently follow every later in-place edit of it…"). The canonical `## Portfolio summary publication` section already carries the same reasoning ("a snapshot linked to the canonical file would follow every later in-place edit of it"). It is the same class as this REQ — a caller-side copy of a paragraph that now has a single home — and it is not covered by any stale pattern, so a future correction to the canonical version would leave this copy teaching the old one. Not fixed inline: it is outside this REQ's named instance, and the fix is the same shape (convert to pointer + add the phrase to the stale-pattern list), so it wants its own RED/GREEN.
