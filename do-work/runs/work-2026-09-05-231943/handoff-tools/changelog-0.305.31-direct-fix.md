## 0.305.31 — Finalization Refuses a Release Whose Implementation Ships Nothing (2026-09-06)

A release is a change to shipped files. Four entries below (0.305.17, 0.305.21, 0.305.22 and 0.305.23) bumped the version for changes under `_dev/tests` alone; the finalizer checked the release payload and never the implementation. Those four entries stand as history, and this entry closes the gap.

- `complete` with a release manifest now lists what the implementation changed (the supplied commit's first-parent diff, or the finalization commit's own paths) and refuses with `RELEASE-WITHOUT-SHIPPED-CHANGE` when nothing lies under a declared module root (`suite/modules.tsv`), `suite/` or `tools/` other than release metadata. A repository that declares no modules is a consumer project and is not guarded.
- Pinned by tests for a maintainer-only commit, metadata-only paths, a merge that brings in only maintainer files, primary-commit provenance, and a consumer repository.
