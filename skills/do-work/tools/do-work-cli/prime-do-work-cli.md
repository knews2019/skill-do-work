# do-work-cli Prime

## Purpose

The standard-library Go module under this directory is the canonical implementation for deterministic do-work operations. Natural-language actions and retained shell paths delegate to it; `queue-kanban` remains a separate board binary.

## Read first

- `cmd/do-work-cli/main.go` is the command registration boundary. Shared packages must not register commands.
- `internal/resultmodel/` and `internal/commandruntime/` own stable command results and rendering.
- `internal/requestmodel/` owns lossless REQ/UR frontmatter documents and authorized field edits.
- `internal/lifecycleadvance/` owns queue selection-and-claim composition, working/archive lifecycle projection, request-bound execution of each mechanical evidence gate, request-bound composition of finalization, the public finalization-first recovery composition — which preserves the claim of a request held for heavy lanes instead of releasing it — and explicit session-end checkpoint mutation. Gate and finalization inputs stay structured and tokenized; judgment remains in the action.

For another command family, use [package routing](lessons-do-work-cli.md#package-routing).

See [package direction](lessons-do-work-cli.md#package-direction).

See [advance phase table](lessons-do-work-cli.md#advance-phase-table).

## Traps

- [family: guard-on-one-entry-path] Guarding one public route or its payload leaves other routes unprotected → test the governing condition at every entry into the shared authority.

- [family: shipped-module-test-self-containment] A maintainer-tree test can run in a consumer with a coincidentally named directory → keep those tests export-ignored and verify the installed module independently.

- [family: time-is-not-authorization] An earlier green does not authorize a later lane → recheck tree identity and expiry per decision, fingerprint effective inputs, and revoke success before retrying.

- [family: lifecycle-section-evidence] Quoted or hidden Markdown can look like workflow evidence → classify real visible sections and preserve route exceptions before advancing.

- [family: cross-action-exception-closure] A lifecycle exception accepted by one gate can remain impossible downstream → carry the same bound evidence through review, release, completion and recovery.

- [family: exact-basename-authority] A filesystem basename is authority before it is presentation: never trim or normalize it before deciding whether it reserves an ID. Canonical writers format the stored minimum-three-digit ID; compatibility readers accept only whole positive-numeric aliases, and documentation or fixtures that construct marker paths count as writers.
- [family: final-boundary-identity] Bind mutation authority to the filesystem object: inode, content digest, and complete mode metadata. Revalidate at the destructive syscall and refuse symlinks or special files as replacement targets. Exclusive creation establishes ownership; absence differs from replacement. Record every completed part of a move before fallible revalidation, so rollback retains both destination ownership and source removal. Recursive deletion must enumerate and recurse through the same directory-bound root; repository containment alone permits deleting unrelated internal targets.
- [family: silent-skip-reads-as-red] A lane that skips silently inside its own test → the runner records a red exit and the drain routes the REQ to remediation for work that never ran; every engine-gated lane announces `SKIP:` on its own output, and the runner keys `skipped` on that prefix **and** a zero exit status, never on a lane-name list — a lane that ran and exited non-zero is red whatever it printed, and its announcement rides along as second evidence on the red finding. The announcement is only trustworthy when it is the lane's own: a nested runner must never tee a fixture lane's output to the enclosing process, or the parent reads the child's `SKIP:` as its own.
- [family: closed-enumeration-for-a-condition] When a contract states a **condition**, a list of examples that happens to agree on those examples is a different contract. Test the condition's ingredients, taken from the governing spec's own definition, not the spellings today's cases use — and never replace one name list with a longer one. A partially-correct enumeration is the dangerous shape: the passing cases make the class look handled. Prefer an affirmative proof the repository itself declares (Git's index, a package marker) over inferring safety from the absence of known-bad names, and check both halves of a two-part proof separately — a differential that reddens on a rename proves nothing about the condition.
- [family: opaque-evidence-projection] A generic fallback or opaque aggregate can discard the exact blocker a caller must act on; derive output, specific typed records, and recovery argv from one observation set, and keep implementation dependencies in the Go standard library.
- [family: alternate-writer-contract-drift] Changing a stored-format contract in its canonical lifecycle owner without sweeping cleanup, recovery, and other alternate writers leaves repository behavior split; pair the canonical regression with equivalent seam coverage at every remaining writer.
- [family: observed-subset-is-not-semantic-completeness] A coherent subset of dirty recovery paths is not proof of the whole transaction; require the configured semantic member set and exact ownership evidence for each admitted byte before finalizing legacy state.
- [family: reaped-by-its-own-parent] A process is reaped by its parent and by nobody else, so terminating a tree leader-first orphans the descendants to init and a killed orphan stays a zombie — and a zombie still satisfies `kill(pid, 0)`. Signal descendants before their parents, one level at a time, and wait for each level to be reaped before climbing; a wait on our own child ends in milliseconds, while a wait on init is measured in seconds. A group signal may only be sent once the group is proved isolated; otherwise signal the bare pid.
- [family: interruptible-blocking-io] Reading confirmation synchronously and cancelling only the surrounding context → signal handling can wait forever for a caller that is still blocked on input; select cancellation against a buffered read result before writes, then leave post-write recovery with its transaction owner.
- [family: smoke-vs-characterization] Registration and smoke checks can stay green while replacement semantics diverge; compare status, ordered evidence, actions, paths, and effects at authority boundaries before retiring an implementation. Toolbox compatibility Markdown and typed audit JSON come from one result. Image backends receive prompts only as argv, run inside owned Unix process groups, and are reaped before scratch cleanup; platforms without provable descendant ownership fail closed. Batch output appears only as one verified `generated/` publication, while an all-failed batch intentionally succeeds with no output.
- [family: collision-fixture-identity] Reusing a target number in both fixture filenames and frontmatter → filename claims can make an incomplete frontmatter identity rule look correct; use unrelated filename numbers plus a malformed adjacent-value negative control.
- [family: projection-before-bounding] Applying a dispatch bound before projecting a frozen target ledger → newly discovered out-of-ledger records consume selector slots and starve retained work; observe canonically without the scheduling bound, project frozen membership, then bound dispatch.
- [family: publication-target-topology-classification] Treating a tracked path, a manifest's self-description, or a matching version VALUE as release ownership → dependency and generated trees can authorize their own metadata, and an independently versioned component that happens to sit on the application's version gets dragged into its bump or refuses it; seed ownership at an independent repository or declared maintainer root, then propagate only through proven relationships. Ownership is one definition in `internal/releaseownership`, asked before any value comparison, by both the plan-time planner and recovery-time discovery.
- [family: commit-preflight-before-side-effect] Passing a validation-only dry run before an owned filesystem claim → commit-only guards can refuse the real transaction after the command has already left durable residue; validate every commit precondition before the first side effect.

- [family: toolchain-behavior-contract] Verify the declared Go floor with that exact toolchain and check newer supported versions. Refusal fixtures must establish the refused condition, such as an unresolvable root name, rather than assume a library representation stays unchanged.

See [scoped command contracts](lessons-do-work-cli.md#scoped-command-contracts) before changing an individual command.

See [verify](lessons-do-work-cli.md#verify).

## Stakes

- Req: lossless shared models and atomic mutation primitives.
- Value: later commands receive one typed snapshot with stable evidence and exact paths.
- Risk: parser drift or unsafe publication can corrupt durable request history or allocate duplicate IDs.

## Lessons

See [`lessons-do-work-cli.md`](lessons-do-work-cli.md) before changing the shared repository evidence contracts named above.
