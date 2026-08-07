---
id: UR-034
title: Preserve association candidates with empty quarantine
created_at: 2026-08-07T22:40:46Z
requests: [REQ-148]
word_count: 5
---

# Preserve Association Candidates with Empty Quarantine

## Summary

Capture the accepted validation finding that commit and unscoped inspect lose every REQ-association candidate when the run-level secret quarantine is empty.

## Referenced Conversation Context

The phrase "accepted issue" in the verbatim input refers to the preceding `do-work validate-feedback` finding about the two-file awk merge in `actions/commit.md` Step 3 and `actions/inspect.md` unscoped mode, plus their staged modular copies.

The merge loads quarantined paths from its first input and filters fresh `uncommitted-inventory.sh` rows from its second input. It identifies the first input with `NR == FNR`. When the quarantine file is empty, it contributes no records, so `NR` and `FNR` remain equal for every record in the inventory file. The first-file branch therefore consumes the complete inventory, not only its first row, and its `next` statement prevents every candidate from reaching the output.

The finding was reproduced with six legitimate `M`/`A`/`D` candidates and no current `X` rows: `associate-files.sh` received empty stdin and exited 1. The intended behavior is for an empty quarantine to emit every safe `M`, `A`, `D`, and `XD` path, while a populated quarantine continues excluding exact retained paths and every current `X` row. In commit mode, the failure makes all safe files skip automatic REQ association; in unscoped inspect mode, it omits them from the association pass.

The accepted remedy is to replace `NR == FNR` with the repository's existing portable `FILENAME == ARGV[1]` discriminator in every checked-in bridge or modular copy still present when the work runs. Keep the `associate-files.sh` interface and REQ-128's once-X-always-X safety behavior unchanged, do not add a micro-helper for this one-condition correction, and add regression coverage for empty and populated quarantine files with multiple safe candidates and an `X` row. The user confirmed this scope after the live reproduction and independent validation.

## Full Verbatim Input

```text
do-work capture-request for accepted issue
```

---
*Captured: 2026-08-07T22:40:46Z*
